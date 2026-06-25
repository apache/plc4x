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

import org.apache.plc4x.java.api.EventPlcConnection;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.listener.EventListener;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.Consumer;

/**
 * Tracks the state that needs to be restored when a cached connection is recreated.
 * <p>
 * This includes:
 * <ul>
 *   <li>Active subscriptions (with their requests, consumers, and handle mappings)</li>
 *   <li>Registered event listeners</li>
 * </ul>
 * <p>
 * When a connection is invalidated and a new one is created, this tracker
 * re-establishes all subscriptions and re-registers all event listeners
 * on the new connection, transparent to the client.
 * <p>
 * Thread Safety: This class uses concurrent collections and is safe for use
 * from multiple threads.
 */
public class ConnectionStateTracker {

    private static final Logger LOGGER = LoggerFactory.getLogger(ConnectionStateTracker.class);

    /** Bound on each per-subscription re-execute during restoration, so a slow controller can't hang reconnect. */
    private static final long RESUBSCRIBE_TIMEOUT_MS = 5_000;

    private final String connectionString;

    // Active subscriptions indexed by a unique ID
    private final Map<String, SubscriptionRecord> subscriptions;

    // Registered event listeners
    private final List<EventListener> eventListeners;

    // Counter for generating unique subscription IDs
    private int subscriptionCounter = 0;

    public ConnectionStateTracker(String connectionString) {
        this.connectionString = connectionString;
        this.subscriptions = new ConcurrentHashMap<>();
        this.eventListeners = new CopyOnWriteArrayList<>();
    }

    /**
     * Records a new subscription for potential recovery.
     *
     * @param request  The subscription request
     * @param consumer The event consumer
     * @param response The subscription response containing handles
     * @return A unique ID for this subscription record
     */
    public String recordSubscription(PlcSubscriptionRequest request,
                                     Consumer<PlcSubscriptionEvent> consumer,
                                     PlcSubscriptionResponse response) {
        String subscriptionId = generateSubscriptionId();
        SubscriptionRecord record = new SubscriptionRecord(request, consumer, response);
        subscriptions.put(subscriptionId, record);
        LOGGER.debug("[{}] Recorded subscription {}: {} tags",
            connectionString, subscriptionId, request.getTagNames().size());
        return subscriptionId;
    }

    /**
     * Removes a subscription record when unsubscribed.
     *
     * @param handles The handles being unsubscribed
     */
    public void removeSubscription(List<PlcSubscriptionHandle> handles) {
        if (handles == null || handles.isEmpty()) {
            return;
        }

        // Find and remove records that contain any of the given handles
        List<String> toRemove = new ArrayList<>();
        for (Map.Entry<String, SubscriptionRecord> entry : subscriptions.entrySet()) {
            for (PlcSubscriptionHandle handle : handles) {
                if (entry.getValue().containsHandle(handle)) {
                    toRemove.add(entry.getKey());
                    break;
                }
            }
        }

        for (String id : toRemove) {
            subscriptions.remove(id);
            LOGGER.debug("[{}] Removed subscription record: {}", connectionString, id);
        }
    }

    /**
     * Gets the current handle for an original handle.
     * <p>
     * After a connection is recreated, clients may still hold references to
     * old handles. This method returns the current valid handle.
     *
     * @param originalHandle The handle the client has
     * @return The current valid handle
     */
    public PlcSubscriptionHandle getCurrentHandle(PlcSubscriptionHandle originalHandle) {
        for (SubscriptionRecord record : subscriptions.values()) {
            if (record.containsHandle(originalHandle)) {
                return record.getCurrentHandle(originalHandle);
            }
        }
        return originalHandle;
    }

    /**
     * Registers an event listener for tracking.
     *
     * @param listener The listener to track
     */
    public void addEventListener(EventListener listener) {
        if (listener != null && !eventListeners.contains(listener)) {
            eventListeners.add(listener);
            LOGGER.debug("[{}] Tracking event listener: {}", connectionString, listener.getClass().getSimpleName());
        }
    }

    /**
     * Removes an event listener from tracking.
     *
     * @param listener The listener to remove
     */
    public void removeEventListener(EventListener listener) {
        if (listener != null) {
            eventListeners.remove(listener);
            LOGGER.debug("[{}] Stopped tracking event listener: {}", connectionString, listener.getClass().getSimpleName());
        }
    }

    /**
     * Restores all tracked state on a new connection.
     * <p>
     * This method is called after a connection has been recreated following
     * invalidation. It re-registers all event listeners and re-establishes
     * all subscriptions on the new connection.
     *
     * @param newConnection The newly created connection
     */
    public void restoreState(PlcConnection newConnection) {
        if (newConnection == null) {
            LOGGER.warn("[{}] Cannot restore state: new connection is null", connectionString);
            return;
        }

        LOGGER.info("[{}] Restoring connection state: {} listeners, {} subscriptions",
            connectionString, eventListeners.size(), subscriptions.size());

        // Re-register event listeners
        restoreEventListeners(newConnection);

        // Re-establish subscriptions
        restoreSubscriptions(newConnection);
    }

    /**
     * Re-registers all tracked event listeners on the new connection.
     */
    private void restoreEventListeners(PlcConnection connection) {
        if (eventListeners.isEmpty()) {
            return;
        }

        if (!(connection instanceof EventPlcConnection eventConnection)) {
            LOGGER.warn("[{}] New connection does not support events, {} listeners will not be restored",
                connectionString, eventListeners.size());
            return;
        }

        for (EventListener listener : eventListeners) {
            try {
                eventConnection.addEventListener(listener);
                LOGGER.debug("[{}] Restored event listener: {}",
                    connectionString, listener.getClass().getSimpleName());
            } catch (Exception e) {
                if (LOGGER.isTraceEnabled()) {
                    LOGGER.error("[{}] Failed to restore event listener: {}",
                        connectionString, listener.getClass().getSimpleName(), e);
                } else {
                    LOGGER.error("[{}] Failed to restore event listener: {}",
                        connectionString, listener.getClass().getSimpleName());
                }
            }
        }

        LOGGER.info("[{}] Restored {} event listeners", connectionString, eventListeners.size());
    }

    /**
     * Re-establishes all tracked subscriptions on the new connection.
     * <p>
     * The original {@link PlcSubscriptionRequest} cannot be re-executed directly — it is bound to the
     * (now dead) connection that built it. Instead, each subscription is rebuilt tag-by-tag (type,
     * underlying tag, interval, consumer) against the new connection's builder and re-executed, and the
     * record's handle mapping is re-pointed so a client still holding a pre-reconnect handle keeps
     * working (e.g. for unsubscription). A subscription that fails to restore is left recorded so a
     * later reconnection can retry it; the failure is logged.
     */
    private void restoreSubscriptions(PlcConnection connection) {
        if (subscriptions.isEmpty()) {
            return;
        }
        int restored = 0;
        int failed = 0;
        for (Map.Entry<String, SubscriptionRecord> entry : subscriptions.entrySet()) {
            SubscriptionRecord record = entry.getValue();
            try {
                PlcSubscriptionResponse newResponse = resubscribe(connection, record);
                record.updateAfterReconnection(newResponse);
                restored++;
                LOGGER.debug("[{}] Restored subscription {} ({} tags) on the new connection",
                    connectionString, entry.getKey(), newResponse.getTagNames().size());
            } catch (Exception e) {
                failed++;
                // Keep the record so a later reconnection can retry; getCurrentHandle still maps.
                if (LOGGER.isTraceEnabled()) {
                    LOGGER.error("[{}] Failed to restore subscription {}", connectionString, entry.getKey(), e);
                } else {
                    LOGGER.error("[{}] Failed to restore subscription {}: {}",
                        connectionString, entry.getKey(), e.getMessage());
                }
            }
        }
        LOGGER.info("[{}] Subscription restoration complete: {} restored, {} failed",
            connectionString, restored, failed);
    }

    /**
     * Re-establishes one recorded subscription on {@code connection} by replaying its tags against a
     * fresh subscription request builder (type, underlying tag, polling / min interval, and per-tag or
     * request-level consumer), then executing it (bounded by {@link #RESUBSCRIBE_TIMEOUT_MS}).
     */
    private PlcSubscriptionResponse resubscribe(PlcConnection connection, SubscriptionRecord record)
            throws InterruptedException, ExecutionException, TimeoutException {
        PlcSubscriptionRequest original = record.getOriginalRequest();
        PlcSubscriptionRequest.Builder builder = connection.subscriptionRequestBuilder();

        // Re-apply the request-level consumer (the common setConsumer(...) pattern); a per-tag consumer
        // below takes precedence for its own tag.
        if (record.getConsumer() != null) {
            builder.setConsumer(record.getConsumer());
        }

        for (String name : original.getTagNames()) {
            PlcSubscriptionTag tag = original.getTag(name);
            PlcTag inner = tag.getTag();
            Consumer<PlcSubscriptionEvent> tagConsumer = original.getTagConsumer(name);
            switch (tag.getPlcSubscriptionType()) {
                case CYCLIC -> {
                    Duration interval = tag.getDuration().orElse(Duration.ofSeconds(1));
                    if (tagConsumer != null) {
                        builder.addCyclicTag(name, inner, interval, tagConsumer);
                    } else {
                        builder.addCyclicTag(name, inner, interval);
                    }
                }
                case CHANGE_OF_STATE -> {
                    if (tag.getDuration().isPresent()) {
                        Duration minInterval = tag.getDuration().get();
                        if (tagConsumer != null) {
                            builder.addChangeOfStateTag(name, inner, tagConsumer, minInterval);
                        } else {
                            builder.addChangeOfStateTag(name, inner, minInterval);
                        }
                    } else if (tagConsumer != null) {
                        builder.addChangeOfStateTag(name, inner, tagConsumer);
                    } else {
                        builder.addChangeOfStateTag(name, inner);
                    }
                }
                case EVENT -> {
                    if (tagConsumer != null) {
                        builder.addEventTag(name, inner, tagConsumer);
                    } else {
                        builder.addEventTag(name, inner);
                    }
                }
            }
        }

        return builder.build().execute().get(RESUBSCRIBE_TIMEOUT_MS, TimeUnit.MILLISECONDS);
    }

    /**
     * Gets the number of tracked subscriptions.
     */
    public int getSubscriptionCount() {
        return subscriptions.size();
    }

    /**
     * Gets the number of tracked event listeners.
     */
    public int getEventListenerCount() {
        return eventListeners.size();
    }

    /**
     * Checks if there is any state to restore.
     */
    public boolean hasStateToRestore() {
        return !subscriptions.isEmpty() || !eventListeners.isEmpty();
    }

    /**
     * Clears all tracked state.
     * Called when the connection is permanently closed.
     */
    public void clear() {
        subscriptions.clear();
        eventListeners.clear();
        LOGGER.debug("[{}] Cleared all tracked state", connectionString);
    }

    private synchronized String generateSubscriptionId() {
        return "sub-" + (++subscriptionCounter);
    }
}
