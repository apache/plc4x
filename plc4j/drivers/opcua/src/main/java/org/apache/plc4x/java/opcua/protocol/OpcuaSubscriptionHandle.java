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

import static java.util.concurrent.Executors.newSingleThreadExecutor;
import static java.util.concurrent.Executors.newSingleThreadScheduledExecutor;

import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.context.Conversation;
import org.apache.plc4x.java.opcua.tag.OpcuaTag;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager.RequestTransaction;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Consumer;

public class OpcuaSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaSubscriptionHandle.class);

    private final static ScheduledExecutorService EXECUTOR = newSingleThreadScheduledExecutor(runnable -> new Thread(runnable, "plc4x-opcua-subscription-scheduler"));

    private final Set<Consumer<PlcSubscriptionEvent>> consumers;
    private final List<String> tagNames;
    private final Conversation conversation;
    private final PlcSubscriptionRequest subscriptionRequest;
    private final OpcuaProtocolLogic plcSubscriber;
    private final Long subscriptionId;
    private final long cycleTime;
    private final long revisedCycleTime;

    private final AtomicLong clientHandles = new AtomicLong(1L);
    private final RequestTransactionManager tm;
    private ScheduledFuture<?> publishTask;

    public OpcuaSubscriptionHandle(OpcuaProtocolLogic plcSubscriber, RequestTransactionManager tm,
        Conversation conversation, PlcSubscriptionRequest subscriptionRequest, Long subscriptionId, long cycleTime) {
        super(plcSubscriber);
        this.tm = tm;
        this.consumers = new HashSet<>();
        this.subscriptionRequest = subscriptionRequest;
        this.tagNames = new ArrayList<>(subscriptionRequest.getTagNames());
        this.conversation = conversation;
        this.subscriptionId = subscriptionId;
        this.plcSubscriber = plcSubscriber;
        this.cycleTime = cycleTime;
        this.revisedCycleTime = cycleTime;
    }

    public CompletableFuture<OpcuaSubscriptionHandle> onSubscribeCreateMonitoredItemsRequest() {
        List<ExtensionObjectDefinition> requestList = new ArrayList<>(this.tagNames.size());
        for (String tagName : this.tagNames) {
            final DefaultPlcSubscriptionTag tagDefaultPlcSubscription = (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(tagName);

            NodeId idNode = OpcuaProtocolLogic.generateNodeId((OpcuaTag) tagDefaultPlcSubscription.getTag());

            ReadValueId readValueId = new ReadValueId(
                idNode,
                0xD,
                OpcuaProtocolLogic.NULL_STRING,
                new QualifiedName(0, OpcuaProtocolLogic.NULL_STRING));

            MonitoringMode monitoringMode;
            switch (tagDefaultPlcSubscription.getPlcSubscriptionType()) {
                case CYCLIC:
                    monitoringMode = MonitoringMode.monitoringModeSampling;
                    break;
                case CHANGE_OF_STATE:
                    monitoringMode = MonitoringMode.monitoringModeReporting;
                    break;
                case EVENT:
                    monitoringMode = MonitoringMode.monitoringModeReporting;
                    break;
                default:
                    monitoringMode = MonitoringMode.monitoringModeReporting;
            }

            long clientHandle = clientHandles.getAndIncrement();

            MonitoringParameters parameters = new MonitoringParameters(
                clientHandle,
                (double) cycleTime,     // sampling interval
                OpcuaProtocolLogic.NULL_EXTENSION_OBJECT,       // filter, null means use default
                1L,   // queue size
                true        // discard oldest
            );

            MonitoredItemCreateRequest request = new MonitoredItemCreateRequest(
                readValueId, monitoringMode, parameters);

            requestList.add(request);
        }

        RequestHeader requestHeader = conversation.createRequestHeader();
        CreateMonitoredItemsRequest createMonitoredItemsRequest = new CreateMonitoredItemsRequest(
            requestHeader,
            subscriptionId,
            TimestampsToReturn.timestampsToReturnBoth,
            requestList.size(),
            requestList
        );

        return conversation.submit(createMonitoredItemsRequest, CreateMonitoredItemsResponse.class)
            .whenComplete((response, error) -> {
                if (error instanceof TimeoutException) {
                    LOGGER.info("Timeout while sending the Create Monitored Item Subscription Message", error);
                } else if (error != null) {
                    LOGGER.info("Error while sending the Create Monitored Item Subscription Message", error);
                }
            }).thenApply(responseMessage -> {
                MonitoredItemCreateResult[] array = responseMessage.getResults().toArray(new MonitoredItemCreateResult[0]);
                for (int index = 0, arrayLength = array.length; index < arrayLength; index++) {
                    MonitoredItemCreateResult result = array[index];
                    if (OpcuaStatusCode.enumForValue(result.getStatusCode().getStatusCode()) != OpcuaStatusCode.Good) {
                        LOGGER.error("Invalid Tag {}, subscription created without this tag", tagNames.get(index));
                    } else {
                        LOGGER.debug("Tag {} was added to the subscription", tagNames.get(index));
                    }
                }

                LOGGER.trace("Scheduling publish event for subscription {}", subscriptionId);
                publishTask = EXECUTOR.scheduleAtFixedRate(this::sendPublishRequest, revisedCycleTime / 2, revisedCycleTime, TimeUnit.MILLISECONDS);
                return this;
            });
    }

    /**
     * Main subscriber loop. For subscription, we still need to send a request the server on every cycle.
     * Which includes a request for an update of the previously agreed upon list of tags.
     * The server will respond at most once every cycle.
     *
     * @return
     */
    private void sendPublishRequest() {
        List<ExtensionObjectDefinition> outstandingAcknowledgements = new LinkedList<>();
        List<Long> outstandingRequests = new LinkedList<>();

        //If we are waiting on a response and haven't received one, just wait until we do. A keep alive will be sent out eventually
        if (outstandingRequests.size() <= 1) {
            RequestHeader requestHeader = conversation.createRequestHeader(this.revisedCycleTime * 10);

            //Make a copy of the outstanding requests, so it isn't modified while we are putting the ack list together.
            List<ExtensionObjectDefinition> acks = new ArrayList<>(outstandingAcknowledgements);
            // do not send -1 when requesting publish, the -1 value indicates NULL value
            // which might result in corruption of subscription for some servers
            int ackLength = acks.size();
            outstandingAcknowledgements.removeAll(acks);

            PublishRequest publishRequest = new PublishRequest(requestHeader, ackLength, acks);
            // we work in external thread - we need to coordinate access to conversation pipeline
            RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> {
                //  Create Consumer for the response message, error and timeout to be sent to the Secure Channel
                conversation.submit(publishRequest, PublishResponse.class).thenAccept(responseMessage -> {
                    outstandingRequests.remove(((ResponseHeader) responseMessage.getResponseHeader()).getRequestHandle());

                    for (long availableSequenceNumber : responseMessage.getAvailableSequenceNumbers()) {
                        outstandingAcknowledgements.add(new SubscriptionAcknowledgement(this.subscriptionId, availableSequenceNumber));
                    }

                    for (ExtensionObject notificationMessage : ((NotificationMessage) responseMessage.getNotificationMessage()).getNotificationData()) {
                        ExtensionObjectDefinition notification = notificationMessage.getBody();
                        if (notification instanceof DataChangeNotification) {
                            LOGGER.trace("Found a Data Change notification");
                            List<ExtensionObjectDefinition> items = ((DataChangeNotification) notification).getMonitoredItems();
                            onSubscriptionValue(items.toArray(new MonitoredItemNotification[0]));
                        } else {
                            LOGGER.warn("Unsupported Notification type");
                        }
                    }
                }).whenComplete((result, error) -> {
                    if (error != null) {
                        LOGGER.warn("Publish request of subscription {} resulted in error reported by server", subscriptionId, error);
                        transaction.failRequest(error);
                    } else {
                        LOGGER.trace("Completed publish request for subscription {}", subscriptionId);
                        transaction.endRequest();
                    }
                });
                outstandingRequests.add(requestHeader.getRequestHandle());
            });
        }
    }


    /**
     * Stop the subscriber either on disconnect or on error
     *
     * @return
     */
    public void stopSubscriber() {
        RequestHeader requestHeader = conversation.createRequestHeader(this.revisedCycleTime * 10);
        List<Long> subscriptions = Collections.singletonList(subscriptionId);
        DeleteSubscriptionsRequest deleteSubscriptionRequest = new DeleteSubscriptionsRequest(requestHeader,
            1,
            subscriptions
        );

        // subscription suspend can be invoked from multiple places, hence we manage transaction side of it
        RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> {
            //  Create Consumer for the response message, error and timeout to be sent to the Secure Channel
            conversation.submit(deleteSubscriptionRequest, DeleteSubscriptionsResponse.class)
                .thenAccept(responseMessage -> publishTask.cancel(true))
                .whenComplete((result, error) -> {
                    if (error != null) {
                        LOGGER.error("Deletion of subscription resulted in error", error);
                        transaction.failRequest(error);
                    } else {
                        transaction.endRequest();
                    }
                    plcSubscriber.removeSubscription(subscriptionId);
                });
        });
    }

    /**
     * Receive the returned values from the OPCUA server and format it so that it can be received by the PLC4X client.
     *
     * @param values - array of data values to be sent to the client.
     */
    private void onSubscriptionValue(MonitoredItemNotification[] values) {
        LinkedHashSet<String> tagNameList = new LinkedHashSet<>();
        List<DataValue> dataValues = new ArrayList<>(values.length);
        for (MonitoredItemNotification value : values) {
            tagNameList.add(tagNames.get((int) value.getClientHandle() - 1));
            dataValues.add(value.getValue());
        }
        Map<String, ResponseItem<PlcValue>> tags = plcSubscriber.readResponse(tagNameList, dataValues);
        final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(), tags);

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
        LOGGER.info("Registering a new OPCUA subscription consumer");
        consumers.add(consumer);
        return new DefaultPlcConsumerRegistration(plcSubscriber, consumer, this);
    }

    public Long getSubscriptionId() {
        return subscriptionId;
    }

}
