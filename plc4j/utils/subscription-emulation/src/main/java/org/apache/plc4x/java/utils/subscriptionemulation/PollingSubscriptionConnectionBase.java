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

package org.apache.plc4x.java.utils.subscriptionemulation;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionFactory;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.model.DefaultPlcSubscriptionHandle;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.tools.eventpump.EventPump;
import org.apache.plc4x.java.tools.eventpump.TagBatch;
import org.apache.plc4x.java.tools.eventpump.triggers.TimerTrigger;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Consumer;

/**
 * Extended base class for PLC connections that emulates subscriptions using polling.
 * <p>
 * This class is designed for drivers that don't natively support subscriptions (like Modbus).
 * It uses the EventPump tool to poll tags at regular intervals and emulates:
 * <ul>
 *   <li>CYCLIC subscriptions - polls at the specified interval</li>
 *   <li>CHANGE_OF_STATE subscriptions - polls and detects value changes</li>
 *   <li>EVENT subscriptions - not supported (throws UnsupportedOperationException)</li>
 * </ul>
 * <p>
 * To use this class, simply extend it instead of ConnectionBase and implement PlcReader:
 * <pre>
 * public class MyConnection extends PollingSubscriptionConnectionBase&lt;MyConfig&gt; implements PlcReader {
 *     // Your driver implementation
 * }
 * </pre>
 *
 * @param <C> the configuration type
 */
public abstract class PollingSubscriptionConnectionBase<C extends Configuration> extends ConnectionBase<C> {

    private static final Logger LOGGER = LoggerFactory.getLogger(PollingSubscriptionConnectionBase.class);

    private final EventPump eventPump = new EventPump();
    private final AtomicLong subscriptionIdGenerator = new AtomicLong(1);
    private final Map<Long, SubscriptionInfo> subscriptions = new ConcurrentHashMap<>();
    private final Map<PlcConsumerRegistration, ConsumerRegistrationInfo> consumerRegistrations = new ConcurrentHashMap<>();
    private final PlcConnectionFactory singleConnectionManager;

    public PollingSubscriptionConnectionBase(C configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
        // Create a connection manager that always returns this connection
        this.singleConnectionManager = new SingleConnectionManager(this);
    }

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        LinkedHashMap<String, PlcResponseItem<PlcSubscriptionHandle>> responseItems = new LinkedHashMap<>();

        try {
            // Group tags by subscription type and polling interval
            Map<SubscriptionKey, Map<String, PlcSubscriptionTag>> tagGroups = groupTagsBySubscriptionType(subscriptionRequest);

            // Create a batch for each group
            for (Map.Entry<SubscriptionKey, Map<String, PlcSubscriptionTag>> entry : tagGroups.entrySet()) {
                SubscriptionKey key = entry.getKey();
                Map<String, PlcSubscriptionTag> groupTags = entry.getValue();

                Long subscriptionId = subscriptionIdGenerator.getAndIncrement();
                String batchId = "subscription-" + subscriptionId;

                // Create tag batch with appropriate trigger
                TimerTrigger trigger = createTrigger(key.pollingInterval);

                // Convert tags into an address map
                Map<String, String> tagAddresses = new HashMap<>();
                groupTags.forEach((name, tag) -> tagAddresses.put(name, tag.getAddressString()));

                // Create subscription info for value change detection
                SubscriptionInfo subscriptionInfo = new SubscriptionInfo(
                    subscriptionId,
                    batchId,
                    key.subscriptionType,
                    groupTags,
                    subscriptionRequest.getConsumer(),
                    new ConcurrentHashMap<>() // tag-specific consumers
                );

                // Add tag-specific consumers
                groupTags.keySet().forEach(tagName -> {
                    Consumer<PlcSubscriptionEvent> tagConsumer = subscriptionRequest.getTagConsumer(tagName);
                    if (tagConsumer != null) {
                        subscriptionInfo.tagConsumers.put(tagName, tagConsumer);
                    }
                });

                TagBatch batch = TagBatch.builder()
                    .withBatchId(batchId)
                    .withConnectionFactory(singleConnectionManager)
                    .withConnectionString("internal://subscription")
                    .addTagAddresses(tagAddresses)
                    .withTrigger(trigger)
                    .withListener((b, response) -> handleSubscriptionEvent(subscriptionInfo, response))
                    .build();

                eventPump.addBatch(batch);
                eventPump.startBatch(batchId);

                subscriptions.put(subscriptionId, subscriptionInfo);

                // Create subscription handles for response
                for (String tagName : groupTags.keySet()) {
                    PlcSubscriptionHandle handle = new PollingSubscriptionHandle(this, subscriptionId, tagName);
                    responseItems.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
                }

                LOGGER.debug("Created {} subscription batch '{}' with {} tags, polling interval: {}ms",
                    key.subscriptionType, batchId, groupTags.size(), key.pollingInterval);
            }

            return CompletableFuture.completedFuture(
                new DefaultPlcSubscriptionResponse(subscriptionRequest, responseItems));

        } catch (Exception e) {
            LOGGER.error("Error creating subscription", e);
            return CompletableFuture.failedFuture(e);
        }
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        Set<Long> subscriptionsToRemove = new HashSet<>();

        for (PlcSubscriptionHandle handle : unsubscriptionRequest.getSubscriptionHandles()) {
            if (handle instanceof PollingSubscriptionHandle pollingHandle) {
                Long subscriptionId = pollingHandle.getSubscriptionId();
                SubscriptionInfo info = subscriptions.get(subscriptionId);

                if (info != null) {
                    subscriptionsToRemove.add(subscriptionId);
                } else {
                    LOGGER.warn("Subscription not found for handle: {}", handle);
                }
            } else {
                LOGGER.warn("Invalid handle type: {}", handle.getClass());
            }
        }

        // Stop and remove all identified subscriptions
        for (Long subscriptionId : subscriptionsToRemove) {
            SubscriptionInfo info = subscriptions.remove(subscriptionId);
            if (info != null) {
                eventPump.stopBatch(info.batchId);
                eventPump.removeBatch(info.batchId);
                LOGGER.debug("Unsubscribed from batch '{}'", info.batchId);
            }
        }

        return CompletableFuture.completedFuture(
            new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest));
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        List<PlcSubscriptionHandle> handleList = new ArrayList<>(handles);
        int consumerId = (int) subscriptionIdGenerator.getAndIncrement();
        PlcConsumerRegistration registration = new PlcConsumerRegistration() {
            @Override
            public void unregister() {
                onUnregisterConsumer(this);
            }

            @Override
            public List<PlcSubscriptionHandle> getSubscriptionHandles() {
                return Collections.unmodifiableList(handleList);
            }

            @Override
            public Integer getConsumerId() {
                return consumerId;
            }
        };
        ConsumerRegistrationInfo info = new ConsumerRegistrationInfo(consumer, new HashSet<>(handleList));
        consumerRegistrations.put(registration, info);

        LOGGER.debug("Registered consumer for {} subscription handles", handles.size());
        return registration;
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        ConsumerRegistrationInfo removed = consumerRegistrations.remove(registration);
        if (removed != null) {
            LOGGER.debug("Unregistered consumer");
        }
    }

    @Override
    public void close() throws Exception {
        // Stop all subscriptions
        eventPump.close();
        subscriptions.clear();
        consumerRegistrations.clear();

        super.close();
        LOGGER.info("Closed polling subscription connection, all subscriptions stopped");
    }

    /**
     * Handles a subscription event from the event pump.
     * For CHANGE_OF_STATE subscriptions, only fires if values have changed.
     */
    private void handleSubscriptionEvent(SubscriptionInfo info, PlcReadResponse response) {
        try {
            boolean shouldFire = false;
            LinkedHashMap<String, PlcResponseItem<PlcValue>> eventValues = new LinkedHashMap<>();

            for (String tagName : response.getTagNames()) {
                PlcResponseCode code = response.getResponseCode(tagName);
                if (code == PlcResponseCode.OK) {
                    PlcValue currentValue = response.getPlcValue(tagName);
                    PlcValue previousValue = info.previousValues.get(tagName);

                    if (info.subscriptionType == PlcSubscriptionType.CYCLIC) {
                        // For cyclic subscriptions, always fire
                        shouldFire = true;
                        eventValues.put(tagName, new DefaultPlcResponseItem<>(code, currentValue));
                        info.previousValues.put(tagName, currentValue);

                    } else if (info.subscriptionType == PlcSubscriptionType.CHANGE_OF_STATE) {
                        // For change-of-state, only fire if value changed
                        if (previousValue == null || !valuesEqual(previousValue, currentValue)) {
                            shouldFire = true;
                            eventValues.put(tagName, new DefaultPlcResponseItem<>(code, currentValue));
                            info.previousValues.put(tagName, currentValue);
                            LOGGER.trace("Value changed for tag '{}': {} -> {}", tagName, previousValue, currentValue);
                        }
                    }
                } else {
                    // Include error responses
                    eventValues.put(tagName, new DefaultPlcResponseItem<>(code, null));
                }
            }

            if (shouldFire && !eventValues.isEmpty()) {
                PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(
                    Instant.now(),
                    eventValues
                );

                // Fire to request-level consumer
                if (info.requestConsumer != null) {
                    info.requestConsumer.accept(event);
                }

                // Fire to tag-specific consumers
                for (String tagName : eventValues.keySet()) {
                    Consumer<PlcSubscriptionEvent> tagConsumer = info.tagConsumers.get(tagName);
                    if (tagConsumer != null) {
                        tagConsumer.accept(event);
                    }
                }

                // Fire to registered consumers
                for (ConsumerRegistrationInfo regInfo : consumerRegistrations.values()) {
                    // Check if this registration includes any of our subscription handles
                    for (PlcSubscriptionHandle handle : regInfo.handles) {
                        if (handle instanceof PollingSubscriptionHandle pollingHandle) {
                            if (pollingHandle.getSubscriptionId().equals(info.subscriptionId)) {
                                regInfo.consumer.accept(event);
                                break;
                            }
                        }
                    }
                }
            }

        } catch (Exception e) {
            LOGGER.error("Error handling subscription event for batch '{}'", info.batchId, e);
        }
    }

    /**
     * Groups subscription tags by type and polling interval.
     */
    private Map<SubscriptionKey, Map<String, PlcSubscriptionTag>> groupTagsBySubscriptionType(PlcSubscriptionRequest request) {
        Map<SubscriptionKey, Map<String, PlcSubscriptionTag>> groups = new LinkedHashMap<>();

        for (String tagName : request.getTagNames()) {
            PlcSubscriptionTag tag = request.getTag(tagName);
            PlcSubscriptionType type = tag.getPlcSubscriptionType();

            // EVENT subscriptions are not supported
            if (type == PlcSubscriptionType.EVENT) {
                throw new UnsupportedOperationException(
                    "EVENT subscriptions are not supported with polling-based subscriptions. " +
                    "Use CYCLIC or CHANGE_OF_STATE instead.");
            }

            long pollingInterval;
            if (type == PlcSubscriptionType.CYCLIC && tag.getDuration().isPresent()) {
                pollingInterval = tag.getDuration().get().toMillis();
            } else {
                // Default polling interval for CHANGE_OF_STATE or if not specified
                pollingInterval = getDefaultPollingInterval();
            }

            SubscriptionKey key = new SubscriptionKey(type, pollingInterval);
            groups.computeIfAbsent(key, k -> new LinkedHashMap<>()).put(tagName, tag);
        }

        return groups;
    }

    /**
     * Creates a trigger based on the polling interval.
     */
    private TimerTrigger createTrigger(long pollingIntervalMs) {
        return new TimerTrigger(pollingIntervalMs, TimeUnit.MILLISECONDS);
    }

    /**
     * Compares two PlcValues for equality.
     * Override this method if you need custom value comparison logic.
     */
    protected boolean valuesEqual(PlcValue v1, PlcValue v2) {
        if (v1 == v2) return true;
        if (v1 == null || v2 == null) return false;

        // Compare as objects
        return Objects.equals(v1.getObject(), v2.getObject());
    }

    /**
     * Get the default polling interval for CHANGE_OF_STATE subscriptions.
     * Override this method to provide a driver-specific default.
     *
     * @return Default polling interval in milliseconds (default: 1000ms)
     */
    protected long getDefaultPollingInterval() {
        return 1000; // 1 second default
    }

    /**
     * Key for grouping subscriptions by type and interval.
     */
    private static class SubscriptionKey {
        final PlcSubscriptionType subscriptionType;
        final long pollingInterval;

        SubscriptionKey(PlcSubscriptionType subscriptionType, long pollingInterval) {
            this.subscriptionType = subscriptionType;
            this.pollingInterval = pollingInterval;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (!(o instanceof SubscriptionKey that)) return false;
            return pollingInterval == that.pollingInterval && subscriptionType == that.subscriptionType;
        }

        @Override
        public int hashCode() {
            return Objects.hash(subscriptionType, pollingInterval);
        }
    }

    /**
     * Information about an active subscription.
     */
    private static class SubscriptionInfo {
        final Long subscriptionId;
        final String batchId;
        final PlcSubscriptionType subscriptionType;
        final Map<String, PlcSubscriptionTag> tags;
        final Consumer<PlcSubscriptionEvent> requestConsumer;
        final Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers;
        final Map<String, PlcValue> previousValues = new ConcurrentHashMap<>();

        SubscriptionInfo(Long subscriptionId,
                        String batchId,
                        PlcSubscriptionType subscriptionType,
                        Map<String, PlcSubscriptionTag> tags,
                        Consumer<PlcSubscriptionEvent> requestConsumer,
                        Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers) {
            this.subscriptionId = subscriptionId;
            this.batchId = batchId;
            this.subscriptionType = subscriptionType;
            this.tags = tags;
            this.requestConsumer = requestConsumer;
            this.tagConsumers = tagConsumers;
        }
    }

    /**
     * Information about a consumer registration.
     */
    private static class ConsumerRegistrationInfo {
        final Consumer<PlcSubscriptionEvent> consumer;
        final Set<PlcSubscriptionHandle> handles;

        ConsumerRegistrationInfo(Consumer<PlcSubscriptionEvent> consumer, Set<PlcSubscriptionHandle> handles) {
            this.consumer = consumer;
            this.handles = handles;
        }
    }

    /**
     * Custom subscription handle for polling-based subscriptions that stores the subscription ID and tag name.
     */
    private static class PollingSubscriptionHandle extends DefaultPlcSubscriptionHandle {
        private final Long subscriptionId;
        private final String tagName;

        PollingSubscriptionHandle(PlcSubscriber subscriber, Long subscriptionId, String tagName) {
            super(subscriber);
            this.subscriptionId = subscriptionId;
            this.tagName = tagName;
        }

        public Long getSubscriptionId() {
            return subscriptionId;
        }

        public String getTagName() {
            return tagName;
        }

        @Override
        public boolean equals(Object obj) {
            if (!(obj instanceof PollingSubscriptionHandle other)) {
                return false;
            }
            return Objects.equals(subscriptionId, other.subscriptionId) &&
                   Objects.equals(tagName, other.tagName);
        }

        @Override
        public int hashCode() {
            return Objects.hash(subscriptionId, tagName);
        }

        @Override
        public String toString() {
            return "PollingSubscriptionHandle{" +
                "subscriptionId=" + subscriptionId +
                ", tagName='" + tagName + '\'' +
                '}';
        }
    }

    /**
     * Simple connection manager that always returns the same connection wrapped in a non-closing proxy.
     * This is used internally to adapt TagBatch's new API to work with a single connection.
     * The wrapper prevents TagBatch from closing the connection when using try-with-resources.
     */
    private static class SingleConnectionManager implements PlcConnectionFactory {
        private final PlcConnection wrappedConnection;

        SingleConnectionManager(PlcConnection connection) {
            // Wrap the connection to prevent it from being closed by TagBatch
            this.wrappedConnection = new NonClosingConnectionWrapper(connection);
        }

        @Override
        public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
            return wrappedConnection;
        }

        @Override
        public PlcConnection getConnection(String connectionString, PlcAuthentication authentication) throws PlcConnectionException {
            return wrappedConnection;
        }
    }

    /**
     * Wrapper that delegates all calls to the underlying connection except close().
     * This prevents TagBatch from closing the connection when used in try-with-resources.
     */
    private static class NonClosingConnectionWrapper implements PlcConnection {
        private final PlcConnection delegate;

        NonClosingConnectionWrapper(PlcConnection delegate) {
            this.delegate = delegate;
        }

        @Override
        public void connect() {
            // Already connected
        }

        @Override
        public boolean isConnected() {
            return delegate.isConnected();
        }

        @Override
        public void close() {
            // Intentionally do nothing - we don't want TagBatch to close the underlying connection
        }

        @Override
        public Optional<PlcTag> parseTagAddress(String tagAddress) {
            return delegate.parseTagAddress(tagAddress);
        }

        @Override
        public Optional<PlcValue> parseTagValue(PlcTag tag, Object... values) {
            return delegate.parseTagValue(tag, values);
        }

        @Override
        public CompletableFuture<? extends org.apache.plc4x.java.api.messages.PlcPingResponse> ping() {
            return delegate.ping();
        }

        @Override
        public org.apache.plc4x.java.api.messages.PlcReadRequest.Builder readRequestBuilder() {
            return delegate.readRequestBuilder();
        }

        @Override
        public org.apache.plc4x.java.api.messages.PlcWriteRequest.Builder writeRequestBuilder() {
            return delegate.writeRequestBuilder();
        }

        @Override
        public org.apache.plc4x.java.api.messages.PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
            return delegate.subscriptionRequestBuilder();
        }

        @Override
        public org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
            return delegate.unsubscriptionRequestBuilder();
        }

        @Override
        public org.apache.plc4x.java.api.messages.PlcBrowseRequest.Builder browseRequestBuilder() {
            return delegate.browseRequestBuilder();
        }

        @Override
        public PlcConnectionMetadata getMetadata() {
            return delegate.getMetadata();
        }
    }
}
