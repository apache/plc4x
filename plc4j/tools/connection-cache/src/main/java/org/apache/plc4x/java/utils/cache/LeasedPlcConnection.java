/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogProvider;
import org.apache.plc4x.java.api.EventPlcConnection;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.listener.EventListener;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicReference;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

/**
 * A wrapper around a PlcConnection that represents a "leased" connection from the cache.
 * <p>
 * This class implements AutoCloseable, so it should be used in a try-with-resources statement.
 * When closed, it returns the connection to the cache instead of actually closing it.
 * <p>
 * Thread Safety: This class is thread-safe using ReentrantLock for better performance
 * under contention compared to synchronized blocks. Multiple threads can call methods on
 * the same leased connection, but only one thread should call close().
 */
public class LeasedPlcConnection implements PlcConnection, EventPlcConnection, AuditLogProvider {

    private static final Logger LOGGER = LoggerFactory.getLogger(LeasedPlcConnection.class);

    private final String connectionName;
    private final AtomicReference<PlcConnection> delegate;
    private final ConnectionStateTracker stateTracker;
    private final ReturnConnectionCallback returnConnectionCallback;
    private final long leaseTime;
    private final long leaseId;
    private final ReentrantLock lock;

    private boolean invalidate;

    /**
     * Creates a new leased connection with state tracking and lease identity support.
     *
     * @param connectionName The connection identifier
     * @param delegate The underlying connection
     * @param stateTracker The state tracker for subscription/listener recovery (may be null)
     * @param returnConnectionCallback Callback to return the connection to the cache (includes lease identity)
     * @param leaseId Unique monotonically increasing identifier for this lease, used to prevent
     *                stale invalidation callbacks from destroying active connections
     */
    public LeasedPlcConnection(String connectionName, PlcConnection delegate,
                               ConnectionStateTracker stateTracker,
                               ReturnConnectionCallback returnConnectionCallback,
                               long leaseId) {
        this.connectionName = connectionName;
        this.delegate = new AtomicReference<>(delegate);
        this.stateTracker = stateTracker;
        this.returnConnectionCallback = returnConnectionCallback;
        this.leaseTime = System.currentTimeMillis();
        this.leaseId = leaseId;
        this.lock = new ReentrantLock();
        this.invalidate = false;
    }

    /**
     * Creates a new leased connection with state tracking support.
     * <p>
     * Uses the new ReturnConnectionCallback with a default leaseId of 0.
     *
     * @param connectionName The connection identifier
     * @param delegate The underlying connection
     * @param stateTracker The state tracker for subscription/listener recovery (may be null)
     * @param returnConnectionCallback Callback to return the connection to the cache (includes lease identity)
     */
    public LeasedPlcConnection(String connectionName, PlcConnection delegate,
                               ConnectionStateTracker stateTracker,
                               ReturnConnectionCallback returnConnectionCallback) {
        this(connectionName, delegate, stateTracker, returnConnectionCallback, 0L);
    }

    /**
     * Creates a new leased connection without state tracking.
     * <p>
     * This constructor is maintained for backward compatibility with existing tests.
     * Adapts the BiConsumer callback to the new ReturnConnectionCallback interface
     * by ignoring the leaseId parameter.
     *
     * @param connectionName The connection identifier
     * @param delegate The underlying connection
     * @param returnConnectionCallback Legacy callback to return the connection to the cache
     */
    public LeasedPlcConnection(String connectionName, PlcConnection delegate,
                               BiConsumer<PlcConnection, Boolean> returnConnectionCallback) {
        this(connectionName, delegate, null,
            (conn, invalidated, lid) -> returnConnectionCallback.accept(conn, invalidated), 0L);
    }

    @Override
    public AuditLog getAuditLog() {
        if (delegate.get() == null) {
            return null;
        }

        PlcConnection connection = delegate.get();
        if(connection instanceof AuditLogProvider auditLogProvider) {
            return auditLogProvider.getAuditLog();
        }

        return null;
    }

    @Override
    public void connect() {
        // The cache already established the connection
        throw new UnsupportedOperationException("Cannot call connect() on a leased connection");
    }

    /**
     * Is also closed from the ConnectionContainer if it's been used for too long.
     */
    @Override
    public void close() {
        lock.lock();
        try {
            // In this case the connection was already closed (possibly by the timer)
            PlcConnection connection = delegate.get();
            if(connection == null) {
                return;
            }

            // Clear the connection reference.
            delegate.set(null);

            // Tell the connection container that the connection is free to be reused.
            // Pass the leaseId so the container can verify this callback comes from
            // the current active lease and ignore stale callbacks from old leases.
            try {
                returnConnectionCallback.accept(connection, invalidate, leaseId);
            } catch (Exception e) {
                LOGGER.warn("Exception in return connection callback for {}: {}", connectionName, e.getMessage());
            }
        } finally {
            lock.unlock();
        }
    }

    /**
     * Handle an exception that occurred while using this connection.
     * Invalidates the connection and removes it from the cache.
     *
     * @param operation The operation that failed (for logging)
     * @param exception The exception that occurred
     */
    private void handleException(String operation, Throwable exception) {
        lock.lock();
        try {
            // Only invalidate once
            if (!invalidate) {
                if (LOGGER.isTraceEnabled()) {
                    LOGGER.warn("Connection error during {}, marking connection for invalidation '{}'", operation, connectionName, exception);
                } else {
                    LOGGER.warn("Connection error during {}, marking connection for invalidation '{}': {}", operation, connectionName, exception.getMessage());
                }
                invalidate = true;
            }
        } finally {
            lock.unlock();
        }
    }

    @Override
    public boolean isConnected() {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return plcConnection.isConnected();
    }


    @Override
    public PlcConnectionMetadata getMetadata() {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return plcConnection.getMetadata();
    }

    @Override
    public Optional<PlcTag> parseTagAddress(String tagAddress) {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return plcConnection.parseTagAddress(tagAddress);
    }

    @Override
    public Optional<PlcValue> parseTagValue(PlcTag tag, Object... values) {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return plcConnection.parseTagValue(tag, values);
    }

    @Override
    public CompletableFuture<? extends PlcPingResponse> ping() {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        try {
            CompletableFuture<? extends PlcPingResponse> future = plcConnection.ping();
            // Wrap the future to catch exceptions during execution
            return future.whenComplete((result, throwable) -> {
                if (throwable != null) {
                    handleException("ping", throwable);
                }
            });
        } catch (Exception e) {
            handleException("ping", e);
            throw e;
        }
    }

    @Override
    public PlcBrowseRequest.Builder browseRequestBuilder() {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        try {
            PlcBrowseRequest.Builder builder = plcConnection.browseRequestBuilder();
            // Wrap the builder to intercept execute() calls
            return new WrappedBrowseRequestBuilder(builder, this::handleException);
        } catch (Exception e) {
            handleException("browseRequestBuilder", e);
            throw e;
        }
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        try {
            PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
            // Wrap the builder to intercept execute() calls
            return new WrappedReadRequestBuilder(builder, this::handleException);
        } catch (Exception e) {
            handleException("readRequestBuilder", e);
            throw e;
        }
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        try {
            PlcWriteRequest.Builder builder = plcConnection.writeRequestBuilder();
            // Wrap the builder to intercept execute() calls
            return new WrappedWriteRequestBuilder(builder, this::handleException);
        } catch (Exception e) {
            handleException("writeRequestBuilder", e);
            throw e;
        }
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        try {
            PlcSubscriptionRequest.Builder builder = plcConnection.subscriptionRequestBuilder();
            // Wrap the builder to intercept execute() calls and track subscriptions
            return new WrappedPlcSubscriptionRequestBuilder(builder, this::handleException, stateTracker);
        } catch (Exception e) {
            handleException("subscriptionRequestBuilder", e);
            throw e;
        }
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        PlcConnection plcConnection = delegate.get();
        if (plcConnection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        try {
            PlcUnsubscriptionRequest.Builder builder = plcConnection.unsubscriptionRequestBuilder();
            // Wrap the builder to intercept execute() calls and track unsubscriptions
            return new WrappedPlcUnsubscriptionRequestBuilder(builder, this::handleException, stateTracker);
        } catch (Exception e) {
            handleException("unsubscriptionRequestBuilder", e);
            throw e;
        }
    }

    /**
     * Get the time this lease was created.
     *
     * @return The time this lease was created, in milliseconds since epoch.
     */
    public long getLeaseTime() {
        return leaseTime;
    }

    /**
     * Get the unique identifier for this lease.
     * <p>
     * The lease ID is a monotonically increasing counter assigned by the
     * ConnectionContainer. It is used to prevent stale invalidation callbacks
     * from old leases from destroying active connections.
     *
     * @return The lease ID (immutable, set at construction time)
     */
    public long getLeaseId() {
        return leaseId;
    }

    /**
     * Get the underlying connection (for testing purposes only).
     *
     * @return The wrapped connection
     */
    AtomicReference<PlcConnection> getDelegate() {
        return delegate;
    }

    /**
     * Check if this connection has been returned to the cache.
     *
     * @return true if returned, false otherwise
     */
    boolean isReturned() {
        return delegate.get() == null;
    }

    /**
     * Callback interface for returning a leased connection to the cache.
     * <p>
     * Extends the original BiConsumer pattern to include a lease identity parameter,
     * allowing the container to verify that the callback comes from the current
     * active lease and ignore stale callbacks from previous leases.
     */
    @FunctionalInterface
    public interface ReturnConnectionCallback {
        /**
         * Called when a leased connection is being returned to the cache.
         *
         * @param connection The underlying connection being returned
         * @param invalidated true if the connection should be invalidated (error occurred),
         *                    false for normal return
         * @param leaseId The lease identity of the returning connection, used by the
         *                container to detect and ignore stale callbacks
         */
        void accept(PlcConnection connection, boolean invalidated, long leaseId);
    }

    /**
     * Functional interface for exception handling callback.
     */
    @FunctionalInterface
    private interface ExceptionHandler {
        void handle(String operation, Throwable exception);
    }

    /**
     * Wrapper for PlcBrowseRequest.Builder that intercepts execute() to handle exceptions.
     */
    private static class WrappedBrowseRequestBuilder implements PlcBrowseRequest.Builder {
        private final PlcBrowseRequest.Builder delegate;
        private final ExceptionHandler exceptionHandler;

        WrappedBrowseRequestBuilder(PlcBrowseRequest.Builder delegate, ExceptionHandler exceptionHandler) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
        }

        @Override
        public PlcBrowseRequest.Builder addQuery(String tagName, String query) {
            delegate.addQuery(tagName, query);
            return this;
        }

        @Override
        public PlcBrowseRequest build() {
            PlcBrowseRequest request = delegate.build();
            // Wrap the request to intercept execute()
            return new WrappedBrowseRequest(request, exceptionHandler);
        }
    }

    /**
     * Wrapper for PlcBrowseRequest that intercepts execute() to handle exceptions.
     */
    private static class WrappedBrowseRequest implements PlcBrowseRequest {
        private final PlcBrowseRequest delegate;
        private final ExceptionHandler exceptionHandler;

        WrappedBrowseRequest(PlcBrowseRequest delegate, ExceptionHandler exceptionHandler) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
        }

        @Override
        public LinkedHashSet<String> getQueryNames() {
            return delegate.getQueryNames();
        }

        @Override
        public PlcQuery getQuery(String tagName) {
            return delegate.getQuery(tagName);
        }

        @Override
        public CompletableFuture<? extends PlcBrowseResponse> execute() {
            try {
                CompletableFuture<? extends PlcBrowseResponse> future = delegate.execute();
                // Wrap the future to catch exceptions during execution
                return future.whenComplete((result, throwable) -> {
                    if (throwable != null) {
                        exceptionHandler.handle("browse.execute", throwable);
                    }
                });
            } catch (Exception e) {
                exceptionHandler.handle("browse.execute", e);
                throw e;
            }
        }

        @Override
        public CompletableFuture<? extends PlcBrowseResponse> executeWithInterceptor(PlcBrowseRequestInterceptor plcBrowseRequestInterceptor) {
            try {
                CompletableFuture<? extends PlcBrowseResponse> future = delegate.executeWithInterceptor(plcBrowseRequestInterceptor);
                // Wrap the future to catch exceptions during execution
                return future.whenComplete((result, throwable) -> {
                    if (throwable != null) {
                        exceptionHandler.handle("browse.executeWithInterceptor", throwable);
                    }
                });
            } catch (Exception e) {
                exceptionHandler.handle("browse.executeWithInterceptor", e);
                throw e;
            }
        }
    }

    /**
     * Wrapper for PlcReadRequest.Builder that intercepts execute() to handle exceptions.
     */
    private static class WrappedReadRequestBuilder implements PlcReadRequest.Builder {
        private final PlcReadRequest.Builder delegate;
        private final ExceptionHandler exceptionHandler;

        WrappedReadRequestBuilder(PlcReadRequest.Builder delegate, ExceptionHandler exceptionHandler) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
        }

        @Override
        public PlcReadRequest.Builder addTagAddress(String tagName, String tagAddress) {
            delegate.addTagAddress(tagName, tagAddress);
            return this;
        }

        @Override
        public PlcReadRequest.Builder addTag(String tagName, PlcTag tag) {
            delegate.addTag(tagName, tag);
            return this;
        }

        @Override
        public PlcReadRequest build() {
            PlcReadRequest request = delegate.build();
            // Wrap the request to intercept execute()
            return new WrappedReadRequest(request, exceptionHandler);
        }
    }

    /**
     * Wrapper for PlcReadRequest that intercepts execute() to handle exceptions.
     */
    private static class WrappedReadRequest implements PlcReadRequest {
        private final PlcReadRequest delegate;
        private final ExceptionHandler exceptionHandler;

        WrappedReadRequest(PlcReadRequest delegate, ExceptionHandler exceptionHandler) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
        }

        @Override
        public int getNumberOfTags() {
            return delegate.getNumberOfTags();
        }

        @Override
        public LinkedHashSet<String> getTagNames() {
            return delegate.getTagNames();
        }

        @Override
        public PlcResponseCode getTagResponseCode(String tagName) {
            return delegate.getTagResponseCode(tagName);
        }

        @Override
        public PlcTag getTag(String tagName) {
            return delegate.getTag(tagName);
        }

        @Override
        public List<PlcTag> getTags() {
            return delegate.getTags();
        }

        @Override
        public CompletableFuture<? extends PlcReadResponse> execute() {
            try {
                CompletableFuture<? extends PlcReadResponse> future = delegate.execute();
                // Wrap the future to catch exceptions during execution
                return future.whenComplete((result, throwable) -> {
                    if (throwable != null) {
                        exceptionHandler.handle("read.execute", throwable);
                    }
                });
            } catch (Exception e) {
                exceptionHandler.handle("read.execute", e);
                throw e;
            }
        }
    }

    /**
     * Wrapper for PlcWriteRequest.Builder that intercepts execute() to handle exceptions.
     */
    private static class WrappedWriteRequestBuilder implements PlcWriteRequest.Builder {
        private final PlcWriteRequest.Builder delegate;
        private final ExceptionHandler exceptionHandler;

        WrappedWriteRequestBuilder(PlcWriteRequest.Builder delegate, ExceptionHandler exceptionHandler) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
        }

        @Override
        public PlcWriteRequest.Builder addTagAddress(String tagName, String tagAddress, Object... values) {
            delegate.addTagAddress(tagName, tagAddress, values);
            return this;
        }

        @Override
        public PlcWriteRequest.Builder addTag(String tagName, PlcTag tag, Object... values) {
            delegate.addTag(tagName, tag, values);
            return this;
        }

        @Override
        public PlcWriteRequest build() {
            PlcWriteRequest request = delegate.build();
            // Wrap the request to intercept execute()
            return new WrappedWriteRequest(request, exceptionHandler);
        }
    }

    /**
     * Wrapper for PlcWriteRequest that intercepts execute() to handle exceptions.
     */
    private static class WrappedWriteRequest implements PlcWriteRequest {
        private final PlcWriteRequest delegate;
        private final ExceptionHandler exceptionHandler;

        WrappedWriteRequest(PlcWriteRequest delegate, ExceptionHandler exceptionHandler) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
        }

        @Override
        public int getNumberOfTags() {
            return delegate.getNumberOfTags();
        }

        @Override
        public LinkedHashSet<String> getTagNames() {
            return delegate.getTagNames();
        }

        @Override
        public PlcResponseCode getTagResponseCode(String tagName) {
            return delegate.getTagResponseCode(tagName);
        }

        @Override
        public PlcTag getTag(String tagName) {
            return delegate.getTag(tagName);
        }

        @Override
        public List<PlcTag> getTags() {
            return delegate.getTags();
        }

        @Override
        public int getNumberOfValues(String tagName) {
            return delegate.getNumberOfValues(tagName);
        }

        @Override
        public PlcValue getPlcValue(String tagName) {
            return delegate.getPlcValue(tagName);
        }

        @Override
        public CompletableFuture<? extends PlcWriteResponse> execute() {
            try {
                CompletableFuture<? extends PlcWriteResponse> future = delegate.execute();
                // Wrap the future to catch exceptions during execution
                return future.whenComplete((result, throwable) -> {
                    if (throwable != null) {
                        exceptionHandler.handle("write.execute", throwable);
                    }
                });
            } catch (Exception e) {
                exceptionHandler.handle("write.execute", e);
                throw e;
            }
        }
    }

    private static class WrappedPlcSubscriptionRequestBuilder implements PlcSubscriptionRequest.Builder {
        private final PlcSubscriptionRequest.Builder delegate;
        private final ExceptionHandler exceptionHandler;
        private final ConnectionStateTracker stateTracker;

        WrappedPlcSubscriptionRequestBuilder(PlcSubscriptionRequest.Builder delegate, ExceptionHandler exceptionHandler, ConnectionStateTracker stateTracker) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
            this.stateTracker = stateTracker;
        }

        @Override
        public PlcSubscriptionRequest.Builder setConsumer(Consumer<PlcSubscriptionEvent> consumer) {
            return delegate.setConsumer(consumer);
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval) {
            return delegate.addCyclicTagAddress(name, tagAddress, pollingInterval);
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval, Consumer<PlcSubscriptionEvent> consumer) {
            return delegate.addCyclicTagAddress(name, tagAddress, pollingInterval, consumer);
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval) {
            return delegate.addCyclicTag(name, tag, pollingInterval);
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval, Consumer<PlcSubscriptionEvent> consumer) {
            return delegate.addCyclicTag(name, tag, pollingInterval, consumer);
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress) {
            return delegate.addChangeOfStateTagAddress(name, tagAddress);
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer) {
            return delegate.addChangeOfStateTagAddress(name, tagAddress, consumer);
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, Duration duration) {
            return delegate.addChangeOfStateTagAddress(name, tagAddress, duration);
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer, Duration duration) {
            return delegate.addChangeOfStateTagAddress(name, tagAddress, consumer, duration);
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag) {
            return delegate.addChangeOfStateTag(name, tag);
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer) {
            return delegate.addChangeOfStateTag(name, tag, consumer);
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag, Duration duration) {
            return delegate.addChangeOfStateTag(name, tag, duration);
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer, Duration duration) {
            return delegate.addChangeOfStateTag(name, tag, consumer, duration);
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress) {
            return delegate.addEventTagAddress(name, tagAddress);
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer) {
            return delegate.addEventTagAddress(name, tagAddress, consumer);
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag) {
            return delegate.addEventTag(name, tag);
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer) {
            return delegate.addEventTag(name, tag, consumer);
        }

        @Override
        public PlcSubscriptionRequest build() {
            PlcSubscriptionRequest request = delegate.build();
            // Wrap the request to intercept execute() and track subscriptions
            return new WrappedSubscriptionRequest(request, exceptionHandler, stateTracker);
        }
    }

    private static class WrappedSubscriptionRequest implements PlcSubscriptionRequest {
        private final PlcSubscriptionRequest delegate;
        private final ExceptionHandler exceptionHandler;
        private final ConnectionStateTracker stateTracker;

        public WrappedSubscriptionRequest(PlcSubscriptionRequest delegate, ExceptionHandler exceptionHandler, ConnectionStateTracker stateTracker) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
            this.stateTracker = stateTracker;
        }

        @Override
        public int getNumberOfTags() {
            return delegate.getNumberOfTags();
        }

        @Override
        public LinkedHashSet<String> getTagNames() {
            return delegate.getTagNames();
        }

        @Override
        public PlcResponseCode getTagResponseCode(String tagName) {
            return delegate.getTagResponseCode(tagName);
        }

        @Override
        public PlcSubscriptionTag getTag(String name) {
            return delegate.getTag(name);
        }

        @Override
        public List<PlcSubscriptionTag> getTags() {
            return delegate.getTags();
        }

        @Override
        public Consumer<PlcSubscriptionEvent> getConsumer() {
            return delegate.getConsumer();
        }

        @Override
        public Consumer<PlcSubscriptionEvent> getTagConsumer(String name) {
            return delegate.getTagConsumer(name);
        }

        @Override
        public CompletableFuture<? extends PlcSubscriptionResponse> execute() {
            try {
                return delegate.execute()
                    .thenApply(result -> {
                        if (result != null && stateTracker != null) {
                            // Record the RAW response (its real handles) so the subscription can be
                            // restored across a cache-managed reconnection.
                            stateTracker.recordSubscription(delegate, delegate.getConsumer(), result);
                        }
                        // Hand the client a response whose handles reject a post-subscribe
                        // register(...) (see WrappedSubscriptionHandle) — such a consumer could not be
                        // restored across a reconnection. Consumers must be set at subscribe time.
                        return result != null ? new WrappedSubscriptionResponse(result) : null;
                    })
                    .whenComplete((result, throwable) -> {
                        if (throwable != null) {
                            exceptionHandler.handle("subscribe.execute", throwable);
                        }
                    });
            } catch (Exception e) {
                exceptionHandler.handle("subscribe.execute", e);
                throw e;
            }
        }
    }

    /**
     * Wraps a subscription response so that the handles handed to the client are
     * {@link WrappedSubscriptionHandle}s that reject post-subscribe {@code register(...)}. Everything
     * else delegates to the real response. The cache itself keeps using the raw response/handles
     * (recorded in the state tracker) for restoration and unsubscription.
     */
    private static final class WrappedSubscriptionResponse implements PlcSubscriptionResponse {
        private final PlcSubscriptionResponse delegate;

        WrappedSubscriptionResponse(PlcSubscriptionResponse delegate) {
            this.delegate = delegate;
        }

        @Override
        public PlcSubscriptionRequest getRequest() {
            return delegate.getRequest();
        }

        @Override
        public Collection<String> getTagNames() {
            return delegate.getTagNames();
        }

        @Override
        public PlcSubscriptionTag getTag(String name) {
            return delegate.getTag(name);
        }

        @Override
        public PlcResponseCode getResponseCode(String name) {
            return delegate.getResponseCode(name);
        }

        @Override
        public PlcSubscriptionHandle getSubscriptionHandle(String name) {
            PlcSubscriptionHandle handle = delegate.getSubscriptionHandle(name);
            return handle != null ? new WrappedSubscriptionHandle(handle) : null;
        }

        @Override
        public Collection<PlcSubscriptionHandle> getSubscriptionHandles() {
            List<PlcSubscriptionHandle> wrapped = new ArrayList<>();
            for (PlcSubscriptionHandle handle : delegate.getSubscriptionHandles()) {
                wrapped.add(new WrappedSubscriptionHandle(handle));
            }
            return wrapped;
        }
    }

    /**
     * A subscription handle handed to a client of the connection cache. It delegates to the real
     * handle for everything except {@link #register(Consumer)}, which is rejected: a consumer attached
     * after subscribing is detached runtime state that the cache cannot reconstruct when it transparently
     * recreates the connection, so it would be silently lost on the first reconnect. Consumers must be
     * supplied at subscribe time via {@code setConsumer(...)} or {@code addXxxTag(..., consumer)}.
     * <p>
     * {@link #unwrap()} exposes the real handle so the cache's own unsubscription / restoration paths
     * keep operating on the underlying connection's handles.
     */
    private static final class WrappedSubscriptionHandle implements PlcSubscriptionHandle {
        private final PlcSubscriptionHandle delegate;

        WrappedSubscriptionHandle(PlcSubscriptionHandle delegate) {
            this.delegate = delegate;
        }

        PlcSubscriptionHandle unwrap() {
            return delegate;
        }

        @Override
        public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
            throw new UnsupportedOperationException(
                "Registering a consumer on a subscription handle after subscribing is not supported by the "
                    + "connection cache: such a consumer cannot be restored when the cache transparently "
                    + "recreates the connection, so it would be silently lost on a reconnect. Provide the "
                    + "consumer at subscribe time via setConsumer(...) or addXxxTag(..., consumer) instead.");
        }
    }

    // ========== EventPlcConnection Implementation ==========

    @Override
    public void addEventListener(EventListener listener) {
        PlcConnection conn = delegate.get();
        if (conn == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        if (conn instanceof EventPlcConnection eventConnection) {
            eventConnection.addEventListener(listener);
            // Track the listener for recovery after reconnection
            if (stateTracker != null) {
                stateTracker.addEventListener(listener);
            }
        } else {
            LOGGER.warn("Underlying connection {} does not support events", connectionName);
        }
    }

    @Override
    public void removeEventListener(EventListener listener) {
        PlcConnection conn = delegate.get();
        if (conn == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        if (conn instanceof EventPlcConnection eventConnection) {
            eventConnection.removeEventListener(listener);
            // Stop tracking the listener
            if (stateTracker != null) {
                stateTracker.removeEventListener(listener);
            }
        } else {
            LOGGER.warn("Underlying connection {} does not support events", connectionName);
        }
    }

    private static class WrappedPlcUnsubscriptionRequestBuilder implements PlcUnsubscriptionRequest.Builder {
        private final PlcUnsubscriptionRequest.Builder delegate;
        private final ExceptionHandler exceptionHandler;
        private final ConnectionStateTracker stateTracker;

        public WrappedPlcUnsubscriptionRequestBuilder(PlcUnsubscriptionRequest.Builder delegate, ExceptionHandler exceptionHandler, ConnectionStateTracker stateTracker) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
            this.stateTracker = stateTracker;
        }

        // Handles handed to the client are WrappedSubscriptionHandles; unwrap to the real handle, then
        // translate it through the state tracker so that a handle obtained BEFORE a cache-managed
        // reconnection maps to the connection's current handle for that subscription. Untracked handles
        // pass through unchanged so the underlying connection always receives a real handle.
        private PlcSubscriptionHandle current(PlcSubscriptionHandle handle) {
            PlcSubscriptionHandle raw = (handle instanceof WrappedSubscriptionHandle wrapped) ? wrapped.unwrap() : handle;
            return stateTracker != null ? stateTracker.getCurrentHandle(raw) : raw;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle plcSubscriptionHandle) {
            return delegate.addHandles(current(plcSubscriptionHandle));
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle plcSubscriptionHandle1, PlcSubscriptionHandle... plcSubscriptionHandles) {
            PlcSubscriptionHandle[] mapped = new PlcSubscriptionHandle[plcSubscriptionHandles.length];
            for (int i = 0; i < plcSubscriptionHandles.length; i++) {
                mapped[i] = current(plcSubscriptionHandles[i]);
            }
            return delegate.addHandles(current(plcSubscriptionHandle1), mapped);
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(Collection<PlcSubscriptionHandle> plcSubscriptionHandles) {
            List<PlcSubscriptionHandle> mapped = new ArrayList<>(plcSubscriptionHandles.size());
            for (PlcSubscriptionHandle handle : plcSubscriptionHandles) {
                mapped.add(current(handle));
            }
            return delegate.addHandles(mapped);
        }

        @Override
        public PlcUnsubscriptionRequest build() {
            PlcUnsubscriptionRequest request = delegate.build();
            // Wrap the request to intercept execute() and track unsubscriptions
            return new WrappedUnsubscriptionRequest(request, exceptionHandler, stateTracker);
        }
    }

    private static class WrappedUnsubscriptionRequest implements PlcUnsubscriptionRequest {
        private final PlcUnsubscriptionRequest delegate;
        private final ExceptionHandler exceptionHandler;
        private final ConnectionStateTracker stateTracker;

        public WrappedUnsubscriptionRequest(PlcUnsubscriptionRequest delegate, ExceptionHandler exceptionHandler, ConnectionStateTracker stateTracker) {
            this.delegate = delegate;
            this.exceptionHandler = exceptionHandler;
            this.stateTracker = stateTracker;
        }

        @Override
        public List<PlcSubscriptionHandle> getSubscriptionHandles() {
            return delegate.getSubscriptionHandles();
        }

        @Override
        public CompletableFuture<PlcUnsubscriptionResponse> execute() {
            // Remove tracked subscriptions before executing
            if (stateTracker != null) {
                List<PlcSubscriptionHandle> handles = delegate.getSubscriptionHandles();
                if (handles != null) {
                    stateTracker.removeSubscription(handles);
                }
            }
            try {
                CompletableFuture<PlcUnsubscriptionResponse> future = delegate.execute();
                // Wrap the future to catch exceptions during execution
                return future.whenComplete((result, throwable) -> {
                    if (throwable != null) {
                        exceptionHandler.handle("unsubscribe.execute", throwable);
                    }
                });
            } catch (Exception e) {
                exceptionHandler.handle("unsubscribe.execute", e);
                throw e;
            }
        }
    }

}
