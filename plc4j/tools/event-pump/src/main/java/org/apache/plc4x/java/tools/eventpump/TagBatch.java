/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.tools.eventpump;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionFactory;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.tools.eventpump.transform.TransformException;
import org.apache.plc4x.java.tools.eventpump.transform.ValueTransformer;
import org.apache.plc4x.java.tools.eventpump.transform.ValueTransformerRegistry;
import org.apache.plc4x.java.tools.eventpump.triggers.Trigger;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * Represents a batch of tags to be fetched from a PLC.
 * <p>
 * A batch consists of:
 * - A unique batch ID
 * - A PLC connection
 * - A map of tag names to tag addresses
 * - A trigger that determines when to fetch the tags
 * - A listener that is notified when tags are fetched
 * <p>
 * When the trigger fires, the batch automatically reads all configured tags
 * from the PLC and notifies the listener with the results.
 * <p>
 * Example using the Builder pattern:
 * <pre>
 * TagBatch batch = TagBatch.builder()
 *     .withBatchId("batch1")
 *     .withConnectionFactory(connectionFactory)
 *     .withConnectionString("ads://192.168.1.1:851")
 *     .addTag("temperature", "MAIN.temperature")
 *     .addTag("pressure", "MAIN.pressure")
 *     .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
 *     .withListener((b, response) -> {
 *         System.out.println("Temperature: " + response.getInt("temperature"));
 *         System.out.println("Pressure: " + response.getInt("pressure"));
 *     })
 *     .build();
 *
 * batch.start();
 * </pre>
 */
public class TagBatch implements AutoCloseable {

    private static final Logger LOGGER = LoggerFactory.getLogger(TagBatch.class);

    private final String batchId;
    private final PlcConnectionFactory connectionFactory;
    private final String connectionString;
    private final Map<String, String> tags; // tagName -> tagAddress
    private final Map<String, String> transforms; // tagName -> transformation expression
    private final Trigger trigger;
    private final ValueTransformer valueTransformer;

    // Exponential backoff state for handling consecutive fetch failures.
    // Prevents log flooding and resource waste during network outages.
    private static final long DEFAULT_MAX_BACKOFF_MS = 60_000;
    private static final long DEFAULT_INITIAL_BACKOFF_MS = 1_000;
    private final long maxBackoffMs;
    private final long initialBackoffMs;
    // Written from whichever thread completes a fetch, read from the trigger
    // thread in fetchTags() — must be safely published.
    private volatile int consecutiveFailures;
    private volatile long nextAllowedFetchTimeMs;
    private volatile long currentBackoffMs;

    // Upper bound on a single fetch cycle. This is NOT a replacement for the driver's
    // own request timeout (configure that in the connection string, e.g.
    // "?request-timeout=10000") — it is deliberately set an order of magnitude higher,
    // and exists only so that a driver which never completes its read future cannot
    // wedge this batch forever: fetchInProgress would stay true and every later trigger
    // would be skipped, silently killing the batch with no recovery.
    private static final long DEFAULT_FETCH_TIMEOUT_MS = 300_000;
    private final long fetchTimeoutMs;

    // Guard against overlapping fetch cycles. If a previous fetch hasn't completed
    // when the next trigger fires, the new fetch is skipped to prevent queue pile-up.
    private final AtomicBoolean fetchInProgress = new AtomicBoolean(false);
    private volatile long lastFetchDurationMs = -1;
    private volatile long skippedFetchCount = 0;

    private volatile TagBatchListener listener;
    private volatile boolean started = false;
    private volatile boolean closed = false;

    // Cache for read request optimization — reused across poll cycles as
    // long as the tag set is unchanged. Re-enabling this cache is the
    // single biggest steady-state polling-CPU optimization for the S7+
    // driver (and others that stamp resolution state on tag instances):
    // a fresh PlcReadRequest every cycle would re-parse every address
    // string and allocate fresh PlcTag instances, defeating any tag-side
    // resolution caches. Invalidated whenever the underlying tag set is
    // mutated via {@link #addTag(String, String)} / {@link #addTags(Map)} /
    // {@link #removeTag(String)} / {@link #removeTags} / {@link #clearTags()}.
    private volatile Map<String, String> cachedTagSnapshot = null;
    private volatile PlcReadRequest cachedReadRequest = null;

    /**
     * Private constructor - use Builder to create instances.
     *
     * @param batchId A unique identifier for this batch
     * @param connectionFactory The connection manager to use for leasing connections
     * @param connectionString The connection string to use for getting connections
     * @param tags A map of tag names to tag addresses
     * @param transforms A map of tag names to transformation expressions
     * @param trigger The trigger that determines when to fetch tags
     * @param listener The listener to be notified when tags are fetched
     * @param valueTransformer The value transformer to use (or null for default from registry)
     * @param transformerRegistry The transformer registry (or null to create default)
     */
    private TagBatch(String batchId, PlcConnectionFactory connectionFactory, String connectionString,
                     Map<String, String> tags, Map<String, String> transforms, Trigger trigger,
                     TagBatchListener listener, ValueTransformer valueTransformer,
                     ValueTransformerRegistry transformerRegistry,
                     long maxBackoffMs, long initialBackoffMs, long fetchTimeoutMs) {
        if (batchId == null || batchId.trim().isEmpty()) {
            throw new IllegalArgumentException("Batch ID cannot be null or empty");
        }
        if (connectionFactory == null) {
            throw new IllegalArgumentException("Connection manager cannot be null");
        }
        if (connectionString == null || connectionString.trim().isEmpty()) {
            throw new IllegalArgumentException("Connection string cannot be null or empty");
        }
        if (tags == null || tags.isEmpty()) {
            tags = new LinkedHashMap<>();
        }
        if (trigger == null) {
            throw new IllegalArgumentException("Trigger cannot be null");
        }

        this.batchId = batchId;
        this.connectionFactory = connectionFactory;
        this.connectionString = connectionString;
        this.tags = new LinkedHashMap<>(tags); // Preserve order
        this.transforms = transforms != null ? new LinkedHashMap<>(transforms) : new LinkedHashMap<>();
        this.trigger = trigger;
        this.listener = listener;

        // Determine the value transformer to use
        if (valueTransformer != null) {
            // Use explicitly provided transformer
            this.valueTransformer = valueTransformer;
        } else {
            // Use default from registry (create default registry if not provided)
            ValueTransformerRegistry registry = transformerRegistry != null ?
                transformerRegistry : ValueTransformerRegistry.createDefault();
            this.valueTransformer = registry.getDefault();
        }

        // Initialize backoff state
        this.maxBackoffMs = maxBackoffMs;
        this.initialBackoffMs = initialBackoffMs;
        this.consecutiveFailures = 0;
        this.nextAllowedFetchTimeMs = 0;
        this.currentBackoffMs = 0;
        this.fetchTimeoutMs = fetchTimeoutMs;

        LOGGER.debug("Created TagBatch '{}' with {} tags ({} with transformations)",
            batchId, tags.size(), this.transforms.size());
    }

    /**
     * Set the listener to be notified when tags are fetched.
     *
     * @param listener The listener
     */
    public void setListener(TagBatchListener listener) {
        this.listener = listener;
    }

    /**
     * Create a new Builder for constructing TagBatch instances.
     *
     * @return A new Builder
     */
    public static Builder builder() {
        return new Builder();
    }

    /**
     * Builder for creating TagBatch instances using a fluent API.
     */
    public static class Builder {
        private String batchId;
        private PlcConnectionFactory connectionFactory;
        private String connectionString;
        private final Map<String, String> tagAddresses = new LinkedHashMap<>();
        private final Map<String, String> transforms = new LinkedHashMap<>();
        private Trigger trigger;
        private TagBatchListener listener;
        private ValueTransformer valueTransformer;
        private ValueTransformerRegistry transformerRegistry;
        private long maxBackoffMs = DEFAULT_MAX_BACKOFF_MS;
        private long initialBackoffMs = DEFAULT_INITIAL_BACKOFF_MS;
        private long fetchTimeoutMs = DEFAULT_FETCH_TIMEOUT_MS;

        /**
         * Set the upper bound for a single fetch cycle.
         * <p>
         * This is a watchdog, not a request timeout: it only exists so a driver that never
         * completes its read future cannot wedge the batch permanently. Configure the actual
         * request timeout on the connection string instead (for example
         * {@code "opcua:tcp://host:4840?request-timeout=10000"}), and leave this comfortably
         * above it — the default is 5 minutes.
         *
         * @param timeout The timeout, or a value &lt;= 0 to disable the watchdog entirely
         * @param unit The time unit of the timeout
         * @return This builder
         */
        public Builder withFetchTimeout(long timeout, TimeUnit unit) {
            this.fetchTimeoutMs = timeout <= 0 ? 0 : unit.toMillis(timeout);
            return this;
        }

        /**
         * Set the batch ID.
         *
         * @param batchId The batch ID
         * @return This builder
         */
        public Builder withBatchId(String batchId) {
            this.batchId = batchId;
            return this;
        }

        /**
         * Set the connection manager.
         *
         * @param connectionFactory The connection manager to use for leasing connections
         * @return This builder
         */
        public Builder withConnectionFactory(PlcConnectionFactory connectionFactory) {
            this.connectionFactory = connectionFactory;
            return this;
        }

        /**
         * Set the connection string.
         *
         * @param connectionString The connection string to use when getting connections
         * @return This builder
         */
        public Builder withConnectionString(String connectionString) {
            this.connectionString = connectionString;
            return this;
        }

        /**
         * Add a single tag address to the batch.
         *
         * @param tagName The name of the tag
         * @param tagAddress The address of the tag
         * @return This builder
         */
        public Builder addTagAddress(String tagName, String tagAddress) {
            this.tagAddresses.put(tagName, tagAddress);
            return this;
        }

        /**
         * Add multiple tag addresses to the batch.
         *
         * @param tags A map of tag names to tag addresses
         * @return This builder
         */
        public Builder addTagAddresses(Map<String, String> tags) {
            this.tagAddresses.putAll(tags);
            return this;
        }

        /**
         * Set the trigger for this batch.
         *
         * @param trigger The trigger
         * @return This builder
         */
        public Builder withTrigger(Trigger trigger) {
            this.trigger = trigger;
            return this;
        }

        /**
         * Set the listener for this batch.
         *
         * @param listener The listener
         * @return This builder
         */
        public Builder withListener(TagBatchListener listener) {
            this.listener = listener;
            return this;
        }

        /**
         * Add a transformation expression for a specific tag.
         *
         * @param tagName The tag name
         * @param expression The transformation expression
         * @return This builder
         */
        public Builder addTransform(String tagName, String expression) {
            this.transforms.put(tagName, expression);
            return this;
        }

        /**
         * Add multiple transformation expressions.
         *
         * @param transforms A map of tag names to transformation expressions
         * @return This builder
         */
        public Builder addTransforms(Map<String, String> transforms) {
            this.transforms.putAll(transforms);
            return this;
        }

        /**
         * Set a custom value transformer.
         * If not set, the default transformer from the registry will be used.
         *
         * @param valueTransformer The value transformer
         * @return This builder
         */
        public Builder withValueTransformer(ValueTransformer valueTransformer) {
            this.valueTransformer = valueTransformer;
            return this;
        }

        /**
         * Set the transformer registry to use.
         * If not set, a default registry will be created with standard transformers.
         *
         * @param transformerRegistry The transformer registry
         * @return This builder
         */
        public Builder withTransformerRegistry(ValueTransformerRegistry transformerRegistry) {
            this.transformerRegistry = transformerRegistry;
            return this;
        }

        /**
         * Set the maximum backoff interval for exponential backoff during consecutive fetch failures.
         * Default is 60 seconds.
         *
         * @param maxBackoffMs The maximum backoff interval in milliseconds
         * @return This builder
         */
        public Builder withMaxBackoffMs(long maxBackoffMs) {
            this.maxBackoffMs = maxBackoffMs;
            return this;
        }

        /**
         * Set the initial backoff interval for exponential backoff during consecutive fetch failures.
         * Default is 1 second. The backoff doubles on each consecutive failure.
         *
         * @param initialBackoffMs The initial backoff interval in milliseconds
         * @return This builder
         */
        public Builder withInitialBackoffMs(long initialBackoffMs) {
            this.initialBackoffMs = initialBackoffMs;
            return this;
        }

        /**
         * Build the TagBatch instance.
         * Validates that all required fields are set.
         * Note: The listener is optional and can be set later using setListener().
         *
         * @return A new TagBatch instance
         * @throws IllegalStateException if any required field is not set
         */
        public TagBatch build() {
            if (batchId == null || batchId.trim().isEmpty()) {
                throw new IllegalStateException("Batch ID must be set");
            }
            if (connectionFactory == null) {
                throw new IllegalStateException("Connection manager must be set");
            }
            if (connectionString == null || connectionString.trim().isEmpty()) {
                throw new IllegalStateException("Connection string must be set");
            }
            if (trigger == null) {
                throw new IllegalStateException("Trigger must be set");
            }
            // Note: listener is optional and can be null

            return new TagBatch(batchId, connectionFactory, connectionString, tagAddresses, transforms,
                trigger, listener, valueTransformer, transformerRegistry, maxBackoffMs, initialBackoffMs,
                fetchTimeoutMs);
        }
    }

    /**
     * Start the batch.
     * This starts the trigger, which will cause tag fetching to begin.
     *
     * @throws IllegalStateException if the batch is already started or no listener is set
     */
    public void start() {
        if (closed) {
            throw new IllegalStateException("Batch is closed");
        }
        if (started) {
            throw new IllegalStateException("Batch is already started");
        }
        if (listener == null) {
            throw new IllegalStateException("Listener must be set before starting");
        }

        // Start the trigger with a callback to fetch tags
        trigger.start(t -> fetchTags());
        started = true;

        LOGGER.info("Started TagBatch '{}'", batchId);
    }

    /**
     * Stop the batch.
     * This stops the trigger, preventing further tag fetching.
     */
    public void stop() {
        if (!started) {
            return;
        }

        trigger.stop();
        started = false;

        LOGGER.info("Stopped TagBatch '{}'", batchId);
    }

    /**
     * Manually fetch tags (regardless of trigger state).
     * This can be used to force an immediate read.
     * <p>
     * Each call to fetchTags() will lease a fresh connection from the connection manager,
     * perform the read, and then return the connection (via try-with-resources).
     * <p>
     * To optimize performance, this method caches the PlcReadRequest and reuses it
     * when the batch configuration (tags) hasn't changed between runs.
     *
     * @return A future that completes when the fetch is done
     */
    public CompletableFuture<Void> fetchTags() {
        if (closed) {
            LOGGER.warn("Batch '{}' is closed, cannot fetch tags", batchId);
            return CompletableFuture.failedFuture(new IllegalStateException("Batch is closed"));
        }
        // If we don't have any tags, we can skip the read
        if (tags.isEmpty()) {
            LOGGER.debug("Batch '{}' has no tags, skipping fetch", batchId);
            return CompletableFuture.completedFuture(null);
        }

        // Skip this trigger if we're within the backoff window due to consecutive failures
        if (System.currentTimeMillis() < nextAllowedFetchTimeMs) {
            LOGGER.trace("Batch '{}' skipping fetch due to backoff (next attempt at {})",
                batchId, nextAllowedFetchTimeMs);
            return CompletableFuture.completedFuture(null);
        }

        // Skip this trigger if the previous fetch hasn't completed yet.
        // This prevents queue pile-up when the PLC response time exceeds the polling interval.
        if (!fetchInProgress.compareAndSet(false, true)) {
            skippedFetchCount++;
            LOGGER.warn("Batch '{}' skipping fetch — previous fetch still in progress " +
                "(last fetch took {}ms, skipped {} consecutive triggers)",
                batchId, lastFetchDurationMs, skippedFetchCount);
            if (listener != null) {
                listener.onFetchSkipped(this, lastFetchDurationMs, skippedFetchCount);
            }
            return CompletableFuture.completedFuture(null);
        }
        skippedFetchCount = 0;
        long fetchStartTimeMs = System.currentTimeMillis();

        LOGGER.debug("Fetching tags for batch '{}' (started={})", batchId, started);

        try {
            // Create a snapshot of tags to avoid ConcurrentModificationException
            Map<String, String> tagSnapshot;
            synchronized (this) {
                tagSnapshot = new LinkedHashMap<>(tags);
            }

            // Lease a connection from the manager for this read operation
            LOGGER.debug("Batch '{}' - Leasing connection for fetch", batchId);
            PlcConnection connection = connectionFactory.getConnection(connectionString);
            LOGGER.debug("Batch '{}' - Connection leased successfully", batchId);

            try {
                PlcReadRequest request;

                // Check if we can reuse the cached request. Reusing the
                // PlcReadRequest reuses the underlying PlcTag instances,
                // so driver-side resolution caches (e.g. S7+'s
                // S7CommPlusResolutionState stamped on S7Tag) survive
                // across poll cycles. Without this cache, every cycle
                // rebuilds tags from address strings, which defeats
                // those resolution caches and re-parses every address.
                Map<String, String> cachedSnap = cachedTagSnapshot;
                PlcReadRequest cachedReq = cachedReadRequest;
                if (cachedReq != null && tagSnapshot.equals(cachedSnap)) {
                    request = cachedReq;
                    LOGGER.trace("Reusing cached read request for batch '{}'", batchId);
                } else {
                    PlcReadRequest.Builder builder = connection.readRequestBuilder();
                    for (Map.Entry<String, String> entry : tagSnapshot.entrySet()) {
                        builder.addTagAddress(entry.getKey(), entry.getValue());
                    }
                    request = builder.build();
                    cachedReadRequest = request;
                    cachedTagSnapshot = tagSnapshot;
                    LOGGER.trace("Built and cached new read request for batch '{}'", batchId);
                }

                // Execute read and ensure the connection is returned AFTER the future completes
                LOGGER.debug("Batch '{}' - Executing read request", batchId);
                CompletableFuture<? extends PlcReadResponse> readFuture = request.execute();

                // Bound the cycle so a driver that never completes its future cannot wedge
                // this batch forever. This does not cancel or shorten the driver's own
                // request timeout — it fires only long after that should have.
                CompletableFuture<? extends PlcReadResponse> guardedFuture = fetchTimeoutMs > 0
                    ? readFuture.orTimeout(fetchTimeoutMs, TimeUnit.MILLISECONDS)
                    : readFuture;

                // Process the response and close the connection afterward
                return guardedFuture
                    .whenComplete((response, throwable) -> {
                        try {
                            if (throwable != null) {
                                // Handle error case — enter or increase backoff
                                enterBackoff();
                                LOGGER.error("Error fetching tags for batch '{}': {}", batchId, throwable.getMessage());
                                if (listener != null) {
                                    listener.onError(this, throwable);
                                }
                            } else {
                                // Handle success case
                                try {
                                    LOGGER.debug("Batch '{}' - Read request completed successfully", batchId);
                                    // Apply transformations if any
                                    PlcReadResponse finalResponse = applyTransformations(response);

                                    // Notify listener
                                    if (listener != null) {
                                        listener.onTagsFetched(this, finalResponse);
                                    }
                                    LOGGER.trace("Successfully fetched tags for batch '{}'", batchId);
                                    // Reset backoff only after listener processed successfully
                                    resetBackoff();
                                } catch (Exception e) {
                                    // Listener failed to process the response — enter backoff
                                    // to prevent rapid-fire retries
                                    enterBackoff();
                                    if (LOGGER.isTraceEnabled()) {
                                        LOGGER.error("Error in tag batch listener for batch '{}'", batchId, e);
                                    } else {
                                        LOGGER.error("Error in tag batch listener for batch '{}': {}", batchId, e.getMessage());
                                    }
                                    if (listener != null) {
                                        listener.onError(this, e);
                                    }
                                }
                            }
                        } finally {
                            // Track fetch duration and release the in-progress guard
                            lastFetchDurationMs = System.currentTimeMillis() - fetchStartTimeMs;
                            fetchInProgress.set(false);
                            LOGGER.debug("Batch '{}' - Fetch cycle completed in {}ms", batchId, lastFetchDurationMs);

                            // ALWAYS close connection after processing (success or error)
                            LOGGER.debug("Batch '{}' - Returning connection to pool", batchId);
                            try {
                                connection.close();
                                LOGGER.debug("Batch '{}' - Connection returned successfully", batchId);
                            } catch (Exception e) {
                                if (LOGGER.isTraceEnabled()) {
                                    LOGGER.error("Error returning connection to pool for batch '{}'", batchId, e);
                                } else {
                                    LOGGER.error("Error returning connection to pool for batch '{}': {}", batchId, e.getMessage());
                                }
                            }
                        }
                    })
                    .thenApply(response -> null); // Convert to CompletableFuture<Void>
            } catch (Exception e) {
                // If we fail to build the request, close the connection immediately
                lastFetchDurationMs = System.currentTimeMillis() - fetchStartTimeMs;
                fetchInProgress.set(false);
                try {
                    connection.close();
                } catch (Exception closeEx) {
                    if (LOGGER.isTraceEnabled()) {
                        LOGGER.error("Error closing connection after request build failure for batch '{}'", batchId, closeEx);
                    } else {
                        LOGGER.error("Error closing connection after request build failure for batch '{}': {}", batchId, closeEx.getMessage());
                    }
                }
                throw e;
            }

        } catch (Exception e) {
            // Connection or request build failure — enter or increase backoff
            lastFetchDurationMs = System.currentTimeMillis() - fetchStartTimeMs;
            fetchInProgress.set(false);
            enterBackoff();
            if (LOGGER.isTraceEnabled()) {
                LOGGER.error("Error building read request for batch '{}'", batchId, e);
            } else {
                LOGGER.error("Error building read request for batch '{}': {}", batchId, e.getMessage());
            }
            if (listener != null) {
                listener.onError(this, e);
            }
            return CompletableFuture.failedFuture(e);
        }
    }

    /**
     * Check if the batch is started.
     *
     * @return true if started, false otherwise
     */
    public boolean isStarted() {
        return started;
    }

    /**
     * Get the batch ID.
     *
     * @return The batch ID
     */
    public String getBatchId() {
        return batchId;
    }

    /**
     * Get the connection manager being used.
     *
     * @return The connection manager
     */
    public PlcConnectionFactory getConnectionFactory() {
        return connectionFactory;
    }

    /**
     * Get the connection string being used.
     *
     * @return The connection string
     */
    public String getConnectionString() {
        return connectionString;
    }

    /**
     * Get the tags being fetched.
     * <p>
     * Returns a snapshot: the tag map is mutable at runtime, so handing out a view of the
     * live map would let a caller iterating it collide with a concurrent
     * {@link #addTag(String, String)} and see a ConcurrentModificationException.
     *
     * @return An unmodifiable snapshot of tag names to addresses
     */
    public synchronized Map<String, String> getTags() {
        return Collections.unmodifiableMap(new LinkedHashMap<>(tags));
    }

    /**
     * Get the names of all tags in the batch.
     *
     * @return An unmodifiable snapshot of the tag names
     */
    public synchronized Set<String> getTagNames() {
        return Collections.unmodifiableSet(new LinkedHashSet<>(tags.keySet()));
    }

    /**
     * Invalidate the cached read request — the next fetch cycle will
     * rebuild from the current tag set. Called whenever the underlying
     * tag map is mutated.
     */
    private void invalidateCache() {
        cachedReadRequest = null;
        cachedTagSnapshot = null;
    }

    /**
     * Add a tag to the batch.
     * Note: If the batch is currently started, the tag will be included in the next fetch cycle.
     *
     * @param tagName The name of the tag
     * @param tagAddress The address of the tag
     * @throws IllegalArgumentException if tagName or tagAddress is null or empty
     */
    public synchronized void addTag(String tagName, String tagAddress) {
        if (tagName == null || tagName.trim().isEmpty()) {
            throw new IllegalArgumentException("Tag name cannot be null or empty");
        }
        if (tagAddress == null || tagAddress.trim().isEmpty()) {
            throw new IllegalArgumentException("Tag address cannot be null or empty");
        }
        if (closed) {
            throw new IllegalStateException("Batch is closed");
        }

        tags.put(tagName, tagAddress);
        invalidateCache();
        LOGGER.debug("Added tag '{}' to batch '{}'", tagName, batchId);
    }

    /**
     * Add multiple tags to the batch.
     * Note: If the batch is currently started, the tags will be included in the next fetch cycle.
     *
     * @param tagsToAdd A map of tag names to tag addresses
     * @throws IllegalArgumentException if tagsToAdd is null or empty
     */
    public synchronized void addTags(Map<String, String> tagsToAdd) {
        if (tagsToAdd == null || tagsToAdd.isEmpty()) {
            throw new IllegalArgumentException("Tags map cannot be null or empty");
        }
        if (closed) {
            throw new IllegalStateException("Batch is closed");
        }

        tags.putAll(tagsToAdd);
        invalidateCache();
        LOGGER.debug("Added {} tags to batch '{}'", tagsToAdd.size(), batchId);
    }

    /**
     * Remove a tag from the batch.
     * Note: If the batch is currently started, the tag will be excluded from the next fetch cycle.
     *
     * @param tagName The name of the tag to remove
     * @return true if the tag was removed, false if it didn't exist
     */
    public synchronized boolean removeTag(String tagName) {
        if (closed) {
            throw new IllegalStateException("Batch is closed");
        }

        boolean removed = tags.remove(tagName) != null;
        if (removed) {
            invalidateCache();
            LOGGER.debug("Removed tag '{}' from batch '{}'", tagName, batchId);
        }
        return removed;
    }

    /**
     * Remove multiple tags from the batch.
     * Note: If the batch is currently started, the tags will be excluded from the next fetch cycle.
     *
     * @param tagNames Collection of tag names to remove
     * @return The number of tags that were actually removed
     */
    public synchronized int removeTags(Collection<String> tagNames) {
        if (closed) {
            throw new IllegalStateException("Batch is closed");
        }
        if (tagNames == null) {
            return 0;
        }

        int count = 0;
        for (String tagName : tagNames) {
            if (tags.remove(tagName) != null) {
                count++;
            }
        }

        if (count > 0) {
            invalidateCache();
            LOGGER.debug("Removed {} tags from batch '{}'", count, batchId);
        }
        return count;
    }

    /**
     * Remove all tags from the batch.
     * Warning: After calling this method, the batch will have no tags and fetchTags() will do nothing
     * until tags are added again.
     */
    public synchronized void clearTags() {
        if (closed) {
            throw new IllegalStateException("Batch is closed");
        }

        int count = tags.size();
        tags.clear();
        invalidateCache();
        LOGGER.debug("Cleared all {} tags from batch '{}'", count, batchId);
    }

    /**
     * Check if the batch contains a specific tag.
     *
     * @param tagName The tag name to check
     * @return true if the tag exists in the batch, false otherwise
     */
    public synchronized boolean hasTag(String tagName) {
        return tags.containsKey(tagName);
    }

    /**
     * Get the number of tags in the batch.
     *
     * @return The number of tags
     */
    public synchronized int getTagCount() {
        return tags.size();
    }

    /**
     * Get the trigger.
     *
     * @return The trigger
     */
    public Trigger getTrigger() {
        return trigger;
    }

    /**
     * Apply transformations to the read response.
     *
     * @param response The original response
     * @return The transformed response (or original if no transformations)
     */
    private PlcReadResponse applyTransformations(PlcReadResponse response) {
        if (transforms.isEmpty()) {
            return response;
        }

        Map<String, PlcValue> transformedValues = new HashMap<>();
        boolean hasTransformations = false;

        for (String tagName : response.getTagNames()) {
            String expression = transforms.get(tagName);
            if (expression != null && !expression.trim().isEmpty()) {
                try {
                    // Get the original value
                    PlcValue originalValue = response.getPlcValue(tagName);

                    // Build context with the original value available as "value"
                    Map<String, PlcValue> context = new HashMap<>();
                    context.put("value", originalValue);

                    // Also add all tag values by name for cross-tag expressions
                    for (String otherTagName : response.getTagNames()) {
                        context.put(otherTagName, response.getPlcValue(otherTagName));
                    }

                    // Apply transformation
                    PlcValue transformedValue = valueTransformer.transform(expression, context);
                    transformedValues.put(tagName, transformedValue);
                    hasTransformations = true;

                    LOGGER.trace("Transformed tag '{}' from {} to {} using expression: {}",
                        tagName, originalValue, transformedValue, expression);

                } catch (TransformException e) {
                    LOGGER.error("Failed to transform tag '{}' with expression '{}': {}",
                        tagName, expression, e.getMessage());
                    // Keep original value on transformation error
                }
            }
        }

        if (hasTransformations) {
            return new TransformedPlcReadResponse(response, transformedValues);
        } else {
            return response;
        }
    }

    /**
     * Enter or increase backoff after a failed fetch attempt.
     * Uses exponential backoff: the window doubles per consecutive failure, capped at
     * maxBackoffMs.
     * <p>
     * The window is derived by doubling the previous one rather than evaluating
     * {@code initialBackoffMs * (1L << (failures - 1))}. That expression is unsafe for a
     * long outage: the multiplication overflows to a negative value from the 55th
     * consecutive failure (with the 1s/60s defaults), which puts the next allowed fetch
     * time in the past and disables the backoff altogether, and the shift itself wraps at
     * 64 (Java masks long shift distances to {@code & 63}), restarting the ramp at the
     * initial window. Both turn a long outage back into full-rate hammering of a PLC that
     * is already in trouble. Doubling a value that is capped at maxBackoffMs every step
     * cannot overflow.
     */
    private void enterBackoff() {
        consecutiveFailures++;
        long previous = currentBackoffMs;
        long next = (previous <= 0) ? initialBackoffMs : previous * 2;
        currentBackoffMs = Math.min(next, maxBackoffMs);
        nextAllowedFetchTimeMs = System.currentTimeMillis() + currentBackoffMs;
        LOGGER.info("Batch '{}' entering backoff: next attempt in {}ms (failure {})",
            batchId, currentBackoffMs, consecutiveFailures);
    }

    /**
     * Reset backoff after a successful fetch. Resumes normal polling immediately.
     */
    private void resetBackoff() {
        if (consecutiveFailures > 0) {
            LOGGER.info("Batch '{}' backoff reset after {} failures", batchId, consecutiveFailures);
            consecutiveFailures = 0;
            nextAllowedFetchTimeMs = 0;
            currentBackoffMs = 0;
        }
    }

    /**
     * Get the current number of consecutive failures (for testing).
     *
     * @return The number of consecutive failures
     */
    int getConsecutiveFailures() {
        return consecutiveFailures;
    }

    /**
     * Get the length of the current backoff window in milliseconds (for testing).
     * <p>
     * Unlike {@link #getNextAllowedFetchTimeMs()} this doesn't decay as time passes, so a test
     * can assert on the window that was computed rather than on whatever is left of it.
     *
     * @return The current backoff window, or 0 if the batch isn't in backoff
     */
    long getCurrentBackoffMs() {
        return currentBackoffMs;
    }

    /**
     * Get the next allowed fetch time in millis (for testing).
     *
     * @return The next allowed fetch time
     */
    long getNextAllowedFetchTimeMs() {
        return nextAllowedFetchTimeMs;
    }

    @Override
    public void close() {
        if (closed) {
            return;
        }

        stop();
        trigger.close();
        closed = true;

        LOGGER.debug("Closed TagBatch '{}'", batchId);
    }

    /**
     * Get the duration of the last fetch cycle in milliseconds.
     *
     * @return The duration in milliseconds, or -1 if no fetch has completed yet
     */
    public long getLastFetchDurationMs() {
        return lastFetchDurationMs;
    }

    /**
     * Listener interface for tag batch events.
     */
    public interface TagBatchListener {
        /**
         * Called when tags are successfully fetched.
         *
         * @param batch The batch that was fetched
         * @param response The read response containing tag values
         */
        void onTagsFetched(TagBatch batch, PlcReadResponse response);

        /**
         * Called when an error occurs during tag fetching.
         *
         * @param batch The batch that encountered the error
         * @param error The error that occurred
         */
        default void onError(TagBatch batch, Throwable error) {
            // Default implementation: log the error (stacktrace only at trace level)
            Logger logger = LoggerFactory.getLogger(TagBatchListener.class);
            if (logger.isTraceEnabled()) {
                logger.error("Error in tag batch '{}': {}", batch.getBatchId(), error.getMessage(), error);
            } else {
                logger.error("Error in tag batch '{}': {}", batch.getBatchId(), error.getMessage());
            }
        }

        /**
         * Called when a fetch trigger is skipped because the previous fetch is still in progress.
         * This indicates that the polling interval is shorter than the PLC response time.
         *
         * @param batch The batch that skipped the fetch
         * @param lastFetchDurationMs The duration of the last completed fetch in milliseconds
         * @param consecutiveSkips The number of consecutive triggers that have been skipped
         */
        default void onFetchSkipped(TagBatch batch, long lastFetchDurationMs, long consecutiveSkips) {
            // Default implementation: already logged at WARN level in TagBatch.fetchTags()
        }
    }
}
