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
package org.apache.plc4x.java.opcua.protocol;

import static java.util.concurrent.Executors.newSingleThreadScheduledExecutor;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.OpcuaConnection;
import org.apache.plc4x.java.opcua.context.Conversation;
import org.apache.plc4x.java.opcua.tag.OpcuaTag;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Consumer;

public class OpcuaSubscriptionHandle implements PlcSubscriptionHandle {

    private final static ScheduledExecutorService EXECUTOR = newSingleThreadScheduledExecutor(runnable -> new Thread(runnable, "plc4x-opcua-subscription-scheduler"));

    private final Logger logger = LoggerFactory.getLogger(OpcuaSubscriptionHandle.class);
    private final Set<Consumer<PlcSubscriptionEvent>> consumers;
    private final Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers;
    private final List<String> tagNames;
    private final Conversation conversation;
    private final PlcSubscriptionRequest subscriptionRequest;
    private final OpcuaConnection plcSubscriber;
    private final Long subscriptionId;
    private final long cycleTime;
    private final long revisedCycleTime;
    private final long queueSize;

    private final AtomicLong clientHandles = new AtomicLong(1L);

    private final List<SubscriptionAcknowledgement> outstandingAcknowledgements = new CopyOnWriteArrayList<>();
    private ScheduledFuture<?> publishTask;

    /** Most recent value per tag, used to re-report CYCLIC tags on their own schedule. */
    private final Map<String, PlcResponseItem<PlcValue>> lastValues = new ConcurrentHashMap<>();
    private final List<ScheduledFuture<?>> cyclicTasks = new CopyOnWriteArrayList<>();

    public OpcuaSubscriptionHandle(OpcuaConnection plcSubscriber,
        Conversation conversation, PlcSubscriptionRequest subscriptionRequest, Long subscriptionId, long cycleTime) {
        this(plcSubscriber, conversation, subscriptionRequest, subscriptionId, cycleTime, cycleTime, 1L);
    }

    /**
     * @param cycleTime        the publishing interval that was requested
     * @param revisedCycleTime the publishing interval the server granted; the publish request
     *                         cadence and its timeouts are derived from this one
     */
    public OpcuaSubscriptionHandle(OpcuaConnection plcSubscriber,
        Conversation conversation, PlcSubscriptionRequest subscriptionRequest, Long subscriptionId,
        long cycleTime, long revisedCycleTime, long queueSize) {
        this.consumers = new HashSet<>();
        this.tagConsumers = new HashMap<>();
        this.subscriptionRequest = subscriptionRequest;
        this.tagNames = new ArrayList<>(subscriptionRequest.getTagNames());
        this.conversation = conversation;
        this.subscriptionId = subscriptionId;
        this.plcSubscriber = plcSubscriber;
        this.cycleTime = cycleTime;
        this.revisedCycleTime = revisedCycleTime;
        this.queueSize = queueSize;
    }

    public CompletableFuture<OpcuaSubscriptionHandle> onSubscribeCreateMonitoredItemsRequest() {
        List<MonitoredItemCreateRequest> requestList = new ArrayList<>(this.tagNames.size());
        for (String tagName : this.tagNames) {
            final DefaultPlcSubscriptionTag tagDefaultPlcSubscription = (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(tagName);

            OpcuaTag opcTag = (OpcuaTag) tagDefaultPlcSubscription.getTag();
            NodeId idNode = OpcuaConnection.generateNodeId(opcTag);

            ReadValueId readValueId = new ReadValueId(
                idNode,
                opcTag.getAttributeId().getValue(),
                OpcuaConnection.NULL_STRING,
                new QualifiedName(0, OpcuaConnection.NULL_STRING));

            MonitoringMode monitoringMode = MonitoringMode.monitoringModeReporting;
            ExtensionObject eventFilter = OpcuaConnection.NULL_EXTENSION_OBJECT;
            if (tagDefaultPlcSubscription.getPlcSubscriptionType() == PlcSubscriptionType.EVENT) {
                NodeId nodeId = new NodeId(new NodeIdFourByte((short) 0, OpcuaNodeIdServicesObjectType.BaseEventType.getValue()));
                List<SimpleAttributeOperand> filterOperand = new ArrayList<>();
                Map<String, String> config = opcTag.getConfig();
                for (Map.Entry<String, String> entry : config.entrySet()) {
                    filterOperand.add(new SimpleAttributeOperand(nodeId,
                        List.of(new QualifiedName(0, new PascalString(entry.getKey()))),
                        AttributeId.Value.getValue(),
                        OpcuaConnection.NULL_STRING
                    ));
                }

                EventFilter filterPayload = new EventFilter(filterOperand, new ContentFilter(null));
                ExpandedNodeId expandedNodeId = new ExpandedNodeId(false, false,
                    new NodeIdFourByte((short) 0, filterPayload.getExtensionId()),
                    null, null
                );
                eventFilter = new BinaryExtensionObjectWithMask(
                    expandedNodeId,
                    new ExtensionObjectEncodingMask(false, false, true),
                    filterPayload
                );
                readValueId = new ReadValueId(
                    idNode,
                    AttributeId.EventNotifier.getValue(),
                    OpcuaConnection.NULL_STRING,
                    new QualifiedName(0, OpcuaConnection.NULL_STRING));
            }

            long clientHandle = clientHandles.getAndIncrement();
            // Each monitored item is sampled at the rate its own tag asked for. Using the
            // subscription's cycle time for every item would silently give all tags the rate of
            // whichever tag happened to come first - see the discussion in GH-1896.
            // When a queue is configured, request the server's fastest sampling (0.0)
            // for change-of-state tags so intermediate values accumulate in the queue
            // between publishes. Otherwise keep the existing per-tag behavior (sample at the tag's
            // own rate, else the publish cycle)
            double samplingInterval;

            if (queueSize > 1
                && tagDefaultPlcSubscription.getPlcSubscriptionType() == PlcSubscriptionType.CHANGE_OF_STATE
                && tagDefaultPlcSubscription.getDuration().isEmpty()) {
                samplingInterval = 0.0;
            } else {
                samplingInterval = tagDefaultPlcSubscription.getDuration()
                    .map(Duration::toMillis).map(Long::doubleValue).orElse((double) cycleTime);
            }
            // Only change-of-state items get the configured queue depth. Cyclic items never sample
            // faster than they publish, so a deeper queue would sit unused; event items are fanned
            // out through a name-keyed map in onEventNotification that cannot hold two values for the
            // same tag, so a deeper queue there would silently drop the extra notifications.
            long itemQueueSize =
                tagDefaultPlcSubscription.getPlcSubscriptionType() == PlcSubscriptionType.CHANGE_OF_STATE
                    ? queueSize : 1L;
            MonitoringParameters parameters = new MonitoringParameters(
                clientHandle,
                samplingInterval,
                eventFilter,       // filter, null means use default
                itemQueueSize,   // queue size
                true        // discard oldest
            );

            MonitoredItemCreateRequest request = new MonitoredItemCreateRequest(readValueId, monitoringMode, parameters);

            requestList.add(request);
        }

        RequestHeader requestHeader = conversation.createRequestHeader();
        CreateMonitoredItemsRequest createMonitoredItemsRequest = new CreateMonitoredItemsRequest(
            requestHeader,
            subscriptionId,
            TimestampsToReturn.timestampsToReturnBoth,
            requestList
        );

        return conversation.submit(createMonitoredItemsRequest, CreateMonitoredItemsResponse.class)
            .whenComplete((response, error) -> {
                if (error instanceof TimeoutException) {
                    logger.info("Timeout while sending the Create Monitored Item Subscription Message", error);
                } else if (error != null) {
                    logger.info("Error while sending the Create Monitored Item Subscription Message", error);
                }
            }).thenApply(responseMessage -> {
                MonitoredItemCreateResult[] array = responseMessage.getResults().toArray(MonitoredItemCreateResult[]::new);
                for (int index = 0, arrayLength = array.length; index < arrayLength; index++) {
                    MonitoredItemCreateResult result = array[index];
                    if (OpcuaStatusCode.enumForValue(result.getStatusCode().getStatusCode()) != OpcuaStatusCode.Good) {
                        logger.error("Invalid Tag {}, subscription created without this tag", tagNames.get(index));
                    } else {
                        logger.debug("Tag {} was added to the subscription", tagNames.get(index));
                    }
                }

                logger.trace("Scheduling publish event for subscription {}", subscriptionId);
                publishTask = EXECUTOR.scheduleAtFixedRate(this::sendPublishRequest, revisedCycleTime / 2, revisedCycleTime, TimeUnit.MILLISECONDS);
                startCyclicEmitters();
                return this;
            });
    }

    /**
     * Main subscriber loop. For subscription, we still need to send a request the server on every cycle.
     * Which includes a request for an update of the previously agreed upon list of tags.
     * The server will respond at most once every cycle.
     */
    private void sendPublishRequest() {
        RequestHeader requestHeader = conversation.createRequestHeader(this.revisedCycleTime * 10);

        //Make a copy of the outstanding requests, so it isn't modified while we are putting the ack list together.
        List<SubscriptionAcknowledgement> acks = new ArrayList<>(outstandingAcknowledgements);
        // do not send -1 when requesting publish, the -1 value indicates NULL value
        // which might result in corruption of subscription for some servers
        int ackLength = acks.size();
        outstandingAcknowledgements.removeAll(acks);

        PublishRequest publishRequest = new PublishRequest(requestHeader, acks);
        logger.trace("Sent publish request with {} acks", ackLength);
        //  Create Consumer for the response message, error and timeout to be sent to the Secure Channel
        conversation.submit(publishRequest, PublishResponse.class).thenAccept(responseMessage -> {
            for (long availableSequenceNumber : responseMessage.getAvailableSequenceNumbers()) {
                outstandingAcknowledgements.add(new SubscriptionAcknowledgement(this.subscriptionId, availableSequenceNumber));
            }

            NotificationMessage message = responseMessage.getNotificationMessage();
            if (message.getNotificationData() != null) {
                for (ExtensionObject notificationMessage : message.getNotificationData()) {
                    ExtensionObjectDefinition notification = notificationMessage.getBody();
                    if (notification instanceof DataChangeNotification data) {
                        logger.trace("Found a Data Change Notification");
                        if (!data.getMonitoredItems().isEmpty()) {
                            onMonitoredValue(data.getMonitoredItems());
                        }
                    } else if (notification instanceof EventNotificationList data) {
                        logger.trace("Found a Event Notification");
                        if (!data.getEvents().isEmpty()) {
                            onEventNotification(data.getEvents());
                        }
                    } else {
                        logger.warn("Unsupported Notification type {}", notification.getClass().getName());
                    }
                }
            }
        }).whenComplete((result, error) -> {
            if (error != null) {
                logger.warn("Publish request of subscription {} resulted in error reported by server", subscriptionId, error);
            } else {
                logger.trace("Completed publish request for subscription {}", subscriptionId);
            }
        });
    }


    /**
     * Stop the subscriber either on disconnect or on error
     */
    public void stopSubscriber() {
        RequestHeader requestHeader = conversation.createRequestHeader(this.revisedCycleTime * 10);
        List<Long> subscriptions = Collections.singletonList(subscriptionId);
        DeleteSubscriptionsRequest deleteSubscriptionRequest = new DeleteSubscriptionsRequest(requestHeader, subscriptions);

        //  Create Consumer for the response message, error and timeout to be sent to the Secure Channel
        conversation.submit(deleteSubscriptionRequest, DeleteSubscriptionsResponse.class)
            .whenComplete((result, error) -> {
                if (error != null) {
                    logger.error("Deletion of subscription resulted in error", error);
                }
                // Stop our own scheduled work regardless of how the server answered - a failed
                // delete must not leave the publish loop and the cyclic emitters running.
                cancelScheduledTasks();
                plcSubscriber.removeSubscription(subscriptionId);
            });
    }

    private void cancelScheduledTasks() {
        if (publishTask != null) {
            publishTask.cancel(true);
            publishTask = null;
        }
        cyclicTasks.forEach(task -> task.cancel(true));
        cyclicTasks.clear();
    }

    /**
     * Receive the returned values from the OPCUA server and format it so that it can be received by the PLC4X client.
     *
     * @param values - array of data values to be sent to the client.
     */
    private void onMonitoredValue(List<MonitoredItemNotification> values) {
        long receiveTs = System.currentTimeMillis();

        List<DataValue> dataValues = new ArrayList<>(values.size());
        Map<String, PlcTag> tagMap = new LinkedHashMap<>();
        for (MonitoredItemNotification value : values) {
            String tagName = tagNames.get((int) value.getClientHandle() - 1);
            PlcTag tag = subscriptionRequest.getTag(tagName).getTag();

            // Per-tag consumers already receive each reading individually
            Consumer<PlcSubscriptionEvent> tagConsumer = tagConsumers.get(tagName);
            if (tagConsumer != null) {
                Map<String, PlcResponseItem<PlcValue>> mappedResponse = plcSubscriber.readResponse(Map.of(tagName, tag), List.of(value.getValue()));
                PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.ofEpochMilli(receiveTs), mappedResponse);
                tagConsumer.accept(event);
            }

            // A single event is name-keyed and cannot hold two values for one tag, so flush the
            // batch before a duplicate would overwrite it - this keeps every queued value.
            if (queueSize > 1 && tagMap.containsKey(tagName)) {
                emitMonitoredBatch(tagMap, dataValues, receiveTs);
                tagMap = new LinkedHashMap<>();
                dataValues = new ArrayList<>();
            }
            tagMap.put(tagName, tag);
            dataValues.add(value.getValue());
        }
        // Flush events
        emitMonitoredBatch(tagMap, dataValues, receiveTs);
    }

    private void emitMonitoredBatch(Map<String, PlcTag> tagMap, List<DataValue> dataValues, long receiveTs) {
        if (tagMap.isEmpty()) {
            return;
        }
        Map<String, PlcResponseItem<PlcValue>> mappedResponse = plcSubscriber.readResponse(tagMap, dataValues);
        // Remember the values so cyclic tags can be reported again on their own schedule even
        // though the server only notifies us when something actually changes - see GH-1102.
        lastValues.putAll(mappedResponse);
        PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.ofEpochMilli(receiveTs), mappedResponse);
        consumers.forEach(plcSubscriptionEventConsumer -> plcSubscriptionEventConsumer.accept(event));
    }

    /**
     * Starts the emitters for CYCLIC tags.
     * <p>
     * OPC UA has no notion of a polling interval: a monitored item reports its initial value and
     * then only reports again when the value changes. A CYCLIC subscription however promises an
     * event every interval, so for those tags we re-report the most recently received value on
     * the requested schedule (see GH-1102). Tags sharing an interval are reported together.
     */
    private void startCyclicEmitters() {
        Map<Long, List<String>> tagsByInterval = new LinkedHashMap<>();
        for (String tagName : tagNames) {
            PlcSubscriptionTag tag = subscriptionRequest.getTag(tagName);
            if (tag.getPlcSubscriptionType() != PlcSubscriptionType.CYCLIC) {
                continue;
            }
            long interval = tag.getDuration().map(Duration::toMillis).orElse(cycleTime);
            tagsByInterval.computeIfAbsent(interval, k -> new ArrayList<>()).add(tagName);
        }

        for (Map.Entry<Long, List<String>> entry : tagsByInterval.entrySet()) {
            long interval = entry.getKey();
            List<String> names = entry.getValue();
            logger.debug("Reporting cyclic tags {} every {}ms", names, interval);
            cyclicTasks.add(EXECUTOR.scheduleAtFixedRate(
                () -> emitCyclicValues(names), interval, interval, TimeUnit.MILLISECONDS));
        }
    }

    /** Reports the last known value of the given tags to the registered consumers. */
    private void emitCyclicValues(List<String> names) {
        try {
            Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
            for (String tagName : names) {
                PlcResponseItem<PlcValue> value = lastValues.get(tagName);
                if (value != null) {
                    values.put(tagName, value);
                }
            }
            if (values.isEmpty()) {
                // Nothing received from the server yet - nothing to report.
                return;
            }
            PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(), values);
            for (String tagName : values.keySet()) {
                Consumer<PlcSubscriptionEvent> tagConsumer = tagConsumers.get(tagName);
                if (tagConsumer != null) {
                    tagConsumer.accept(new DefaultPlcSubscriptionEvent(Instant.now(),
                        Map.of(tagName, values.get(tagName))));
                }
            }
            consumers.forEach(consumer -> consumer.accept(event));
        } catch (Exception e) {
            // A failing consumer must not kill the scheduled task for good.
            logger.error("Error while reporting cyclic values for {}", names, e);
        }
    }

    private void onEventNotification(List<EventFieldList> events) {
        long receiveTs = System.currentTimeMillis();

        Map<String, PlcResponseItem<PlcValue>> tagValues = new LinkedHashMap<>();
        for (EventFieldList event : events) {
            String tagName = tagNames.get((int) event.getClientHandle() - 1);
            OpcuaTag tag = (OpcuaTag) subscriptionRequest.getTag(tagName).getTag();

            Iterator<String> fieldNames = tag.getConfig().keySet().iterator();
            Map<String, PlcValue> mapping = new LinkedHashMap<>();
            for (Variant variant : event.getEventFields()) {
                if (fieldNames.hasNext()) {
                    String fieldName = fieldNames.next();
                    PlcValue plcValue = OpcuaConnection.variantToPlcValue(tag, variant);
                    if (plcValue == null) {
                        // Unsupported variant type: keep the field in the struct as an empty
                        // value instead of putting a raw null into it.
                        logger.error("Event field '{}' has unsupported variant type {}", fieldName,
                            variant.getClass().getSimpleName());
                        plcValue = new PlcNull();
                    }
                    mapping.put(fieldName, plcValue);
                    tagValues.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcStruct(mapping)));
                } else {
                    logger.error("Could not map event notification response, subscription received more data than expected");
                    tagValues.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, new PlcNull()));
                }
            }
        }

        PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.ofEpochMilli(receiveTs), tagValues);
        consumers.forEach(plcSubscriptionEventConsumer -> plcSubscriptionEventConsumer.accept(event));
    }

    /**
     * Registers a new Consumer, this allows multiple PLC4X consumers to use the same subscription.
     *
     * @param consumer - Consumer to be used to send any returned values.
     * @return PlcConsumerRegistration - return the important information back to the client.
     */
    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        logger.info("Registering a new OPCUA subscription consumer");
        consumers.add(consumer);
        return new DefaultPlcConsumerRegistration(plcSubscriber, consumer, this);
    }

    public PlcConsumerRegistration registerTagConsumer(String tagName, Consumer<PlcSubscriptionEvent> consumer) {
        logger.info("Registering a new OPCUA subscription consumer for tag with name {}", tagName);
        tagConsumers.put(tagName, consumer);
        return new DefaultPlcConsumerRegistration(plcSubscriber, consumer, this);
    }

    public Long getSubscriptionId() {
        return subscriptionId;
    }

}
