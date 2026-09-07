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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collections;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * An event-driven system for fetching PLC tag data and delivering it to consumers.
 * <p>
 * The EventPump manages multiple {@link TagBatch}es, each of which:
 * - Connects to one or more PLCs
 * - Defines a set of tags to fetch
 * - Has a trigger that determines when to fetch
 * - Delivers results to a listener
 * <p>
 * One EventPump instance can manage batches from multiple PLCs, and multiple
 * batches from the same PLC with different triggers.
 * <p>
 * Features:
 * - Multiple batches per PLC
 * - Multiple PLCs per pump
 * - Different triggers per batch (timer, subscription, etc.)
 * - Thread-safe batch management
 * - Automatic lifecycle management
 * <p>
 * Example usage:
 * <pre>
 * EventPump pump = new EventPump();
 *
 * // Add a batch that fetches every 5 seconds
 * TagBatch batch1 = new TagBatch("batch1", connection, tags, new TimerTrigger(5, TimeUnit.SECONDS));
 * batch1.setListener((batch, response) -> {
 *     // Process results
 * });
 * pump.addBatch(batch1);
 *
 * // Add a batch triggered by a tag change
 * TagBatch batch2 = new TagBatch("batch2", connection, tags2, new SubscriptionTrigger(connection, "trigger", "MAIN.trigger"));
 * batch2.setListener((batch, response) -> {
 *     // Process results
 * });
 * pump.addBatch(batch2);
 *
 * // Start all batches
 * pump.startAll();
 *
 * // Stop when done
 * pump.close();
 * </pre>
 */
public class EventPump implements AutoCloseable {

    private static final Logger LOGGER = LoggerFactory.getLogger(EventPump.class);

    // Thread-safe map of batch ID -> TagBatch
    private final ConcurrentHashMap<String, TagBatch> batches = new ConcurrentHashMap<>();

    private volatile boolean closed = false;

    /**
     * Create a new EventPump.
     */
    public EventPump() {
        LOGGER.debug("Created EventPump");
    }

    /**
     * Add a batch to the pump.
     * The batch will not start automatically - call {@link #startBatch(String)} or {@link #startAll()}.
     *
     * @param batch The batch to add
     * @throws IllegalArgumentException if a batch with the same ID already exists
     * @throws IllegalStateException if the pump is closed
     */
    public void addBatch(TagBatch batch) {
        if (batch == null) {
            throw new IllegalArgumentException("Batch cannot be null");
        }
        if (closed) {
            throw new IllegalStateException("EventPump is closed");
        }

        String batchId = batch.getBatchId();

        TagBatch existing = batches.putIfAbsent(batchId, batch);
        if (existing != null) {
            throw new IllegalArgumentException("Batch with ID '" + batchId + "' already exists");
        }

        LOGGER.info("Added batch '{}' to EventPump", batchId);
    }

    /**
     * Remove a batch from the pump.
     * The batch will be stopped and closed.
     *
     * @param batchId The ID of the batch to remove
     * @return The removed batch, or null if not found
     */
    public TagBatch removeBatch(String batchId) {
        if (batchId == null) {
            throw new IllegalArgumentException("Batch ID cannot be null");
        }

        TagBatch batch = batches.remove(batchId);
        if (batch != null) {
            try {
                batch.close();
                LOGGER.info("Removed batch '{}' from EventPump", batchId);
            } catch (Exception e) {
                if (LOGGER.isTraceEnabled()) {
                    LOGGER.error("Error closing removed batch '{}'", batchId, e);
                } else {
                    LOGGER.error("Error closing removed batch '{}': {}", batchId, e.getMessage());
                }
            }
        }

        return batch;
    }

    /**
     * Get a batch by ID.
     *
     * @param batchId The batch ID
     * @return The batch, or null if not found
     */
    public TagBatch getBatch(String batchId) {
        return batches.get(batchId);
    }

    /**
     * Get all batches.
     *
     * @return An unmodifiable map of batch ID -> TagBatch
     */
    public Map<String, TagBatch> getAllBatches() {
        return Collections.unmodifiableMap(batches);
    }

    /**
     * Get all batch IDs.
     *
     * @return An unmodifiable set of batch IDs
     */
    public java.util.Set<String> getBatchIds() {
        return Collections.unmodifiableSet(batches.keySet());
    }

    /**
     * Start a specific batch.
     *
     * @param batchId The ID of the batch to start
     * @throws IllegalArgumentException if the batch doesn't exist
     * @throws IllegalStateException if the pump is closed
     */
    public void startBatch(String batchId) {
        if (closed) {
            throw new IllegalStateException("EventPump is closed");
        }

        TagBatch batch = batches.get(batchId);
        if (batch == null) {
            throw new IllegalArgumentException("Batch '" + batchId + "' not found");
        }

        batch.start();
        LOGGER.info("Started batch '{}'", batchId);
    }

    /**
     * Stop a specific batch.
     *
     * @param batchId The ID of the batch to stop
     * @throws IllegalArgumentException if the batch doesn't exist
     */
    public void stopBatch(String batchId) {
        TagBatch batch = batches.get(batchId);
        if (batch == null) {
            throw new IllegalArgumentException("Batch '" + batchId + "' not found");
        }

        batch.stop();
        LOGGER.info("Stopped batch '{}'", batchId);
    }

    /**
     * Start all batches.
     *
     * @throws IllegalStateException if the pump is closed
     */
    public void startAll() {
        if (closed) {
            throw new IllegalStateException("EventPump is closed");
        }

        batches.values().forEach(batch -> {
            if (!batch.isStarted()) {
                try {
                    batch.start();
                } catch (Exception e) {
                    if (LOGGER.isTraceEnabled()) {
                        LOGGER.error("Error starting batch '{}'", batch.getBatchId(), e);
                    } else {
                        LOGGER.error("Error starting batch '{}': {}", batch.getBatchId(), e.getMessage());
                    }
                }
            }
        });

        LOGGER.info("Started all batches ({} total)", batches.size());
    }

    /**
     * Stop all batches.
     */
    public void stopAll() {
        batches.values().forEach(batch -> {
            if (batch.isStarted()) {
                try {
                    batch.stop();
                } catch (Exception e) {
                    if (LOGGER.isTraceEnabled()) {
                        LOGGER.error("Error stopping batch '{}'", batch.getBatchId(), e);
                    } else {
                        LOGGER.error("Error stopping batch '{}': {}", batch.getBatchId(), e.getMessage());
                    }
                }
            }
        });

        LOGGER.info("Stopped all batches ({} total)", batches.size());
    }

    /**
     * Get the number of batches in the pump.
     *
     * @return The number of batches
     */
    public int getBatchCount() {
        return batches.size();
    }

    /**
     * Get the number of started batches.
     *
     * @return The number of started batches
     */
    public int getStartedBatchCount() {
        return (int) batches.values().stream()
            .filter(TagBatch::isStarted)
            .count();
    }

    /**
     * Check if the pump is closed.
     *
     * @return true if closed, false otherwise
     */
    public boolean isClosed() {
        return closed;
    }

    @Override
    public void close() {
        if (closed) {
            return;
        }

        LOGGER.info("Closing EventPump with {} batches", batches.size());

        // Stop and close all batches
        batches.values().forEach(batch -> {
            try {
                batch.close();
            } catch (Exception e) {
                if (LOGGER.isTraceEnabled()) {
                    LOGGER.error("Error closing batch '{}'", batch.getBatchId(), e);
                } else {
                    LOGGER.error("Error closing batch '{}': {}", batch.getBatchId(), e.getMessage());
                }
            }
        });

        batches.clear();
        closed = true;

        LOGGER.info("EventPump closed");
    }
}
