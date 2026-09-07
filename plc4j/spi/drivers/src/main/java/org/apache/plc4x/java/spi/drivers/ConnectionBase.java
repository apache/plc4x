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
package org.apache.plc4x.java.spi.drivers;

import org.apache.plc4x.java.api.EventPlcConnection;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.listener.EventListener;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcConnectionStateChangedEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.functions.*;
import org.apache.plc4x.java.spi.drivers.messages.*;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.drivers.throttle.RequestThrottle;
import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Supplier;

/**
 * Base class for PLC connections that provides adaptive I/O handling.
 * Automatically detects and uses AsyncTransportInstance for event-driven I/O,
 * or falls back to polling with a daemon thread for non-async transports.
 *
 * @param <C> the configuration type
 */
public abstract class ConnectionBase<C extends Configuration> implements PlcConnection, EventPlcConnection, AuditLogProvider, PlcReader, PlcWriter, PlcSubscriber, PlcBrowser, PlcPinger {

    private static final Logger LOGGER = LoggerFactory.getLogger(ConnectionBase.class);

    protected final C configuration;
    protected final TransportInstance<?> transportInstance;
    protected final AuditLog auditLog;
    protected final RequestThrottle requestThrottle;
    private volatile boolean pollingRunning = false;

    // Connection identity metadata (set by DriverBase after construction)
    private String protocolCode = "unknown";
    private String protocolName = "unknown";
    private String transportCode = "unknown";
    private String transportName = "unknown";
    // Authentication passed to PlcDriverManager.getConnection(url, authentication), set by
    // DriverBase after construction. Null when the caller didn't supply one (credentials may
    // then still come from the connection-string parameters).
    private PlcAuthentication authentication;

    // Event listeners for connection state changes
    private final List<EventListener> eventListeners = new ArrayList<>();

    public ConnectionBase(C configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        this.configuration = configuration;
        this.transportInstance = transportInstance;
        this.auditLog = auditLog;

        int maxConcurrent = getMaxConcurrentRequests();
        this.requestThrottle = new RequestThrottle(maxConcurrent);
    }

    protected C getConfiguration() {
        return configuration;
    }

    protected int getMaxConcurrentRequests() {
        return 1;
    }

    protected void setMaxConcurrentRequests(int newMax) {
        requestThrottle.adjustMaxConcurrentRequests(newMax);
    }

    protected <R> CompletableFuture<R> executeThrottled(Supplier<CompletableFuture<R>> requestSupplier) {
        return requestThrottle.execute(requestSupplier);
    }

    protected boolean isAsyncTransport() {
        return transportInstance instanceof AsyncTransportInstance;
    }

    public TransportInstance<?> getTransportInstance() {
        return transportInstance;
    }

    @Override
    public final void connect() throws PlcConnectionException {
        onConnect();
    }

    protected void onConnect() throws PlcConnectionException {
        // Default implementation does nothing; drivers override this.
    }

    protected void startReceiving(Runnable dataHandler) {
        if (transportInstance instanceof AsyncTransportInstance<?> asyncTransport) {
            asyncTransport.registerDataListener(dataHandler);
            asyncTransport.registerDisconnectListener(this::onTransportDisconnected);
            LOGGER.info("Started async event-driven receive mode");
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.SYSTEM, "Started async event-driven receive mode");
            }
        } else {
            startPollingReceiveLoop(dataHandler);
            LOGGER.info("Started polling receive mode with daemon thread");
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.SYSTEM, "Started polling receive mode");
            }
        }
    }

    protected void stopReceiving() {
        if (transportInstance instanceof AsyncTransportInstance<?> asyncTransport) {
            asyncTransport.removeDataListener();
            asyncTransport.removeDisconnectListener();
        } else {
            pollingRunning = false;
        }
    }

    protected void onTransportDisconnected(Throwable cause) {
        if (cause != null) {
            LOGGER.error("Transport disconnected unexpectedly", cause);
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.ERROR, "Transport disconnected: " + cause.getMessage());
            }
        } else {
            LOGGER.info("Transport disconnected gracefully");
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.SYSTEM, "Transport disconnected");
            }
        }
    }

    private void startPollingReceiveLoop(Runnable dataHandler) {
        pollingRunning = true;
        Thread receiverThread = new Thread(() -> {
            LOGGER.debug("Polling receive loop started");
            while (pollingRunning && !Thread.currentThread().isInterrupted()) {
                try {
                    dataHandler.run();
                    Thread.onSpinWait();
                } catch (Exception e) {
                    if (pollingRunning) {
                        LOGGER.error("Error in receive loop", e);
                    }
                }
            }
            LOGGER.debug("Polling receive loop stopped");
        }, getClass().getSimpleName() + "-Receiver");
        receiverThread.setDaemon(true);
        receiverThread.start();
    }

    @Override
    public boolean isConnected() {
        return false;
    }

    @Override
    public void close() throws Exception {
        stopReceiving();
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CLOSE, "Connection closed");
        }
    }

    protected abstract PlcTagHandler getTagHandler();
    protected abstract PlcValueHandler getValueHandler();

    @Override
    public Optional<PlcTag> parseTagAddress(String tagAddress) {
        try {
            return Optional.of(getTagHandler().parseTag(tagAddress));
        } catch (PlcInvalidTagException e) {
            return Optional.empty();
        }
    }

    @Override
    public Optional<PlcValue> parseTagValue(PlcTag plcTag, Object... objects) {
        try {
            return Optional.of(getValueHandler().newPlcValue(plcTag, objects));
        } catch (PlcInvalidTagException e) {
            return Optional.empty();
        }
    }

    void setConnectionInfo(String protocolCode, String protocolName, String transportCode, String transportName) {
        this.protocolCode = protocolCode;
        this.protocolName = protocolName;
        this.transportCode = transportCode;
        this.transportName = transportName;
    }

    void setAuthentication(PlcAuthentication authentication) {
        this.authentication = authentication;
    }

    /**
     * The authentication supplied to {@code getConnection(url, authentication)}, or {@code null}
     * if none was given. Drivers should prefer this over connection-string credentials.
     */
    protected PlcAuthentication getAuthentication() {
        return authentication;
    }

    /**
     * Derived from the operations the driver actually implements.
     * <p>
     * Every {@code on*} hook below has a default that fails the request with "... not supported",
     * so a driver only supports an operation if it overrides the corresponding hook. Reporting
     * a fixed {@code true} here - as this class used to - told every caller that every driver
     * supports everything, which is wrong for most of them: a driver that connects but implements
     * nothing claimed full support, and callers that check the metadata before deciding what to do
     * were sent into requests that could only fail.
     * <p>
     * The lookup walks the class hierarchy up to (but excluding) {@link ConnectionBase}, so a hook
     * inherited from an intermediate base counts as implemented - {@code PollingSubscriptionConnectionBase}
     * for instance provides {@code onSubscribe} for drivers whose protocol has no native subscriptions.
     * <p>
     * What this can tell is what the <em>driver</em> implements, which is a property of the class.
     * Where the answer depends on the device at the other end - a controller that does not offer an
     * operation the driver otherwise speaks - only the driver can know, so it should override this
     * method and answer from what it learned during connect. This derivation is the default for
     * drivers that have no such distinction to make, and the result is deliberately computed per
     * call so an override is free to vary per connection.
     */
    @Override
    public PlcConnectionMetadata getMetadata() {
        return deriveMetadata(getClass());
    }

    private static PlcConnectionMetadata deriveMetadata(Class<?> connectionClass) {
        boolean read = overridesHook(connectionClass, "onRead", PlcReadRequest.class);
        boolean write = overridesHook(connectionClass, "onWrite", PlcWriteRequest.class);
        boolean subscribe = overridesHook(connectionClass, "onSubscribe", PlcSubscriptionRequest.class);
        boolean browse = overridesHook(connectionClass, "onBrowse", PlcBrowseRequest.class)
            || overridesHook(connectionClass, "onBrowseWithInterceptor", PlcBrowseRequest.class, PlcBrowseRequestInterceptor.class);
        LOGGER.debug("Derived metadata for {}: read={}, write={}, subscribe={}, browse={}",
            connectionClass.getSimpleName(), read, write, subscribe, browse);
        return new PlcConnectionMetadata() {
            @Override
            public boolean isReadSupported() {
                return read;
            }

            @Override
            public boolean isWriteSupported() {
                return write;
            }

            @Override
            public boolean isSubscribeSupported() {
                return subscribe;
            }

            @Override
            public boolean isBrowseSupported() {
                return browse;
            }
        };
    }

    /**
     * Whether {@code connectionClass} or any class between it and {@link ConnectionBase} declares
     * the given hook - i.e. whether the default "not supported" implementation has been replaced.
     */
    private static boolean overridesHook(Class<?> connectionClass, String name, Class<?>... parameterTypes) {
        for (Class<?> current = connectionClass; (current != null) && (current != ConnectionBase.class); current = current.getSuperclass()) {
            try {
                current.getDeclaredMethod(name, parameterTypes);
                return true;
            } catch (NoSuchMethodException e) {
                // Not declared here - keep walking up towards ConnectionBase.
            }
        }
        return false;
    }

    @Override
    public String getProtocolCode() {
        return protocolCode;
    }

    @Override
    public String getProtocolName() {
        return protocolName;
    }

    @Override
    public String getTransportCode() {
        return transportCode;
    }

    @Override
    public String getTransportName() {
        return transportName;
    }

    @Override
    public CompletableFuture<? extends PlcPingResponse> ping() {
        return ping(new DefaultPlcPingRequest(this));
    }

    public void handleUnexpectedIncomingMessage(Message message) {
        LOGGER.warn("Received unexpected message: {}", message);
    }

    @Override
    public PlcBrowseRequest.Builder browseRequestBuilder() {
        return new DefaultPlcBrowseRequest.Builder(this, getTagHandler());
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, getTagHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return new DefaultPlcWriteRequest.Builder(this, getTagHandler(), getValueHandler());
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        return new DefaultPlcSubscriptionRequest.Builder(this, getTagHandler());
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Template methods with audit logging
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public final CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.API_REQUEST, "Read request: " + readRequest.getTagNames());
        }
        return onRead(readRequest);
    }

    @Override
    public final CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.API_REQUEST, "Write request: " + writeRequest.getTagNames());
        }
        return onWrite(writeRequest);
    }

    @Override
    public final CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.API_REQUEST, "Subscribe request");
        }
        return onSubscribe(subscriptionRequest);
    }

    @Override
    public final CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.API_REQUEST, "Browse request");
        }
        return onBrowse(browseRequest);
    }

    @Override
    public final CompletableFuture<PlcBrowseResponse> browseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        return onBrowseWithInterceptor(browseRequest, interceptor);
    }

    @Override
    public final CompletableFuture<PlcPingResponse> ping(PlcPingRequest pingRequest) {
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.API_REQUEST, "Ping request");
        }
        return onPing(pingRequest);
    }

    @Override
    public final CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return onUnsubscribe(unsubscriptionRequest);
    }

    @Override
    public final PlcConsumerRegistration registerConsumer(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        return onRegisterConsumer(consumer, handles);
    }

    @Override
    public final void unregisterConsumer(PlcConsumerRegistration registration) {
        onUnregisterConsumer(registration);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Protected on*() hooks - drivers override these
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        return CompletableFuture.failedFuture(new PlcRuntimeException("Read not supported"));
    }

    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        return CompletableFuture.failedFuture(new PlcRuntimeException("Write not supported"));
    }

    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        return CompletableFuture.failedFuture(new PlcRuntimeException("Subscribe not supported"));
    }

    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return CompletableFuture.failedFuture(new PlcRuntimeException("Unsubscribe not supported"));
    }

    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        throw new PlcRuntimeException("Register consumer not supported");
    }

    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        throw new PlcRuntimeException("Unregister consumer not supported");
    }

    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        return CompletableFuture.failedFuture(new PlcRuntimeException("Browse not supported"));
    }

    protected CompletableFuture<PlcBrowseResponse> onBrowseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        return CompletableFuture.failedFuture(new PlcRuntimeException("Browse not supported"));
    }

    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        return CompletableFuture.failedFuture(new PlcRuntimeException("Ping not supported"));
    }

    @Override
    public AuditLog getAuditLog() {
        return auditLog;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // EventPlcConnection implementation
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void addEventListener(EventListener listener) {
        if (listener != null) {
            synchronized (eventListeners) {
                eventListeners.add(listener);
            }
        }
    }

    @Override
    public void removeEventListener(EventListener listener) {
        if (listener != null) {
            synchronized (eventListeners) {
                eventListeners.remove(listener);
            }
        }
    }

    protected void fireConnectionStateChanged(ConnectionStateChangeType changeType, String details) {
        PlcConnectionStateChangedEvent event = new PlcConnectionStateChangedEvent(changeType, details);
        LOGGER.debug("Firing connection state changed event: type={}, details={}", changeType, details);

        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.SYSTEM, "Connection state changed: " + changeType +
                (details != null ? " - " + details : ""));
        }

        List<EventListener> listenersCopy;
        synchronized (eventListeners) {
            listenersCopy = new ArrayList<>(eventListeners);
        }

        for (EventListener listener : listenersCopy) {
            if (listener instanceof ConnectionStateListener connectionStateListener) {
                try {
                    connectionStateListener.onConnectionStateChanged(event);
                } catch (Exception e) {
                    LOGGER.warn("Error notifying connection state listener", e);
                }
            }
        }
    }

}
