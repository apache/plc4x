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
package org.apache.plc4x.java.opcua.protocol;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.context.SecureChannel;
import org.apache.plc4x.java.opcua.field.OpcuaField;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.readwrite.io.ExtensionObjectIO;
import org.apache.plc4x.java.opcua.readwrite.types.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

public class OpcuaSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaSubscriptionHandle.class);

    private Set<Consumer<PlcSubscriptionEvent>> consumers;
    private List<String> fieldNames;
    private SecureChannel channel;
    private PlcSubscriptionRequest subscriptionRequest;
    private AtomicBoolean destroy = new AtomicBoolean(false);
    private OpcuaProtocolLogic plcSubscriber;
    private Long subscriptionId;
    private long cycleTime;
    private long revisedCycleTime;
    private boolean complete = false;

    private final AtomicLong clientHandles = new AtomicLong(1L);

    private ConversationContext context;

    public OpcuaSubscriptionHandle(ConversationContext<OpcuaAPU> context, OpcuaProtocolLogic plcSubscriber, SecureChannel channel, PlcSubscriptionRequest subscriptionRequest, Long subscriptionId, long cycleTime) {
        super(plcSubscriber);
        this.consumers = new HashSet<>();
        this.subscriptionRequest = subscriptionRequest;
        this.fieldNames = new ArrayList<>( subscriptionRequest.getFieldNames() );
        this.channel = channel;
        this.subscriptionId = subscriptionId;
        this.plcSubscriber = plcSubscriber;
        this.cycleTime = cycleTime;
        this.revisedCycleTime = cycleTime;
        this.context = context;
        try {
            onSubscribeCreateMonitoredItemsRequest().get();
        } catch (Exception e) {
            LOGGER.info("Unable to serialize the Create Monitored Item Subscription Message");
            e.printStackTrace();
            plcSubscriber.onDisconnect(context);
        }
        startSubscriber();
    }

    private CompletableFuture<CreateMonitoredItemsResponse> onSubscribeCreateMonitoredItemsRequest()  {
        MonitoredItemCreateRequest[] requestList = new MonitoredItemCreateRequest[this.fieldNames.size()];
        for (int i = 0; i <  this.fieldNames.size(); i++) {
            final DefaultPlcSubscriptionField fieldDefaultPlcSubscription = (DefaultPlcSubscriptionField) subscriptionRequest.getField(fieldNames.get(i));

            NodeId idNode = generateNodeId((OpcuaField) fieldDefaultPlcSubscription.getPlcField());

            ReadValueId readValueId = new ReadValueId(
                idNode,
                0xD,
                OpcuaProtocolLogic.NULL_STRING,
                new QualifiedName(0, OpcuaProtocolLogic.NULL_STRING));

            MonitoringMode monitoringMode;
            switch (fieldDefaultPlcSubscription.getPlcSubscriptionType()) {
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

            requestList[i] = request;
        }

        CompletableFuture<CreateMonitoredItemsResponse> future = new CompletableFuture<>();

        RequestHeader requestHeader = new RequestHeader(channel.getAuthenticationToken(),
            SecureChannel.getCurrentDateTime(),
            channel.getRequestHandle(),
            0L,
            OpcuaProtocolLogic.NULL_STRING,
            SecureChannel.REQUEST_TIMEOUT_LONG,
            OpcuaProtocolLogic.NULL_EXTENSION_OBJECT);

        CreateMonitoredItemsRequest createMonitoredItemsRequest = new CreateMonitoredItemsRequest(
            requestHeader,
            subscriptionId,
            TimestampsToReturn.timestampsToReturnNeither,
            requestList.length,
            requestList
        );

        ExpandedNodeId expandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(createMonitoredItemsRequest.getIdentifier())),
            null,
            null);

        ExtensionObject extObject = new ExtensionObject(
            expandedNodeId,
            null,
            createMonitoredItemsRequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            ExtensionObjectIO.staticSerialize(buffer, extObject);

            Consumer<byte[]> consumer = opcuaResponse -> {
                CreateMonitoredItemsResponse responseMessage = null;
                try {
                    ExtensionObjectDefinition unknownExtensionObject = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, ByteOrder.LITTLE_ENDIAN), false).getBody();
                    if (unknownExtensionObject instanceof CreateMonitoredItemsResponse) {
                        responseMessage = (CreateMonitoredItemsResponse) unknownExtensionObject;
                    } else {
                        ServiceFault serviceFault = (ServiceFault) unknownExtensionObject;
                        ResponseHeader header = (ResponseHeader) serviceFault.getResponseHeader();
                        LOGGER.error("Subscription ServiceFault returned from server with error code,  '{}'", header.getServiceResult().toString());
                        plcSubscriber.onDisconnect(context);
                    }
                } catch (ParseException e) {
                    LOGGER.error("Unable to parse the returned Subscription response");
                    e.printStackTrace();
                    plcSubscriber.onDisconnect(context);
                }
                for (MonitoredItemCreateResult result : Arrays.stream(responseMessage.getResults()).toArray(MonitoredItemCreateResult[]::new)) {
                    if (OpcuaStatusCode.enumForValue(result.getStatusCode().getStatusCode()) != OpcuaStatusCode.Good) {
                        LOGGER.error("Invalid Field {}, subscription created without this field", fieldNames.get((int) result.getMonitoredItemId()));
                    } else {
                        LOGGER.debug("Field {} was added to the subscription", fieldNames.get((int) result.getMonitoredItemId() - 1));
                    }
                }
                future.complete(responseMessage);
            };

            Consumer<TimeoutException> timeout = e -> {
                LOGGER.info("Timeout while sending the Create Monitored Item Subscription Message");
                e.printStackTrace();
                plcSubscriber.onDisconnect(context);
            };

            BiConsumer<OpcuaAPU, Throwable> error = (message, e) -> {
                LOGGER.info("Error while sending the Create Monitored Item Subscription Message");
                e.printStackTrace();
                plcSubscriber.onDisconnect(context);
            };

            channel.submit(context, timeout, error, consumer, buffer);

        } catch (ParseException e) {
            LOGGER.info("Unable to serialize the Create Monitored Item Subscription Message");
            e.printStackTrace();
            plcSubscriber.onDisconnect(context);
        }
        return future;
    }

    private void sleep(long length) {
        try {
            Thread.sleep(length);
        } catch (InterruptedException e) {
            LOGGER.trace("Interrupted Exception");
        }
    }

    /**
     * Main subscriber loop. For subscription we still need to send a request the server on every cycle.
     * Which includes a request for an update of the previsouly agreed upon list of tags.
     * The server will respond at most once every cycle.
     * @return
     */
    public void startSubscriber() {
        LOGGER.trace("Starting Subscription");
        CompletableFuture.supplyAsync(() -> {
            try {
                LinkedList<SubscriptionAcknowledgement> outstandingAcknowledgements = new LinkedList<>();
                LinkedList<Long> outstandingRequests = new LinkedList<>();
                while (!this.destroy.get()) {
                    long requestHandle = channel.getRequestHandle();

                    //If we are waiting on a response and haven't received one, just wait until we do. A keep alive will be sent out eventually
                    if (outstandingRequests.size() <= 1) {
                        RequestHeader requestHeader = new RequestHeader(channel.getAuthenticationToken(),
                            SecureChannel.getCurrentDateTime(),
                            requestHandle,
                            0L,
                            OpcuaProtocolLogic.NULL_STRING,
                            this.revisedCycleTime * 10,
                            OpcuaProtocolLogic.NULL_EXTENSION_OBJECT);

                        //Make a copy of the outstanding requests so it isn't modified while we are putting the ack list together.
                        LinkedList<Long> outstandingAcknowledgementsSnapshot = (LinkedList<Long>) outstandingAcknowledgements.clone();
                        SubscriptionAcknowledgement[] acks = new SubscriptionAcknowledgement[outstandingAcknowledgementsSnapshot.size()];
                        ;
                        outstandingAcknowledgementsSnapshot.toArray(acks);
                        int ackLength = outstandingAcknowledgementsSnapshot.size() == 0 ? -1 : outstandingAcknowledgementsSnapshot.size();
                        outstandingAcknowledgements.removeAll(outstandingAcknowledgementsSnapshot);

                        PublishRequest publishRequest = new PublishRequest(
                            requestHeader,
                            ackLength,
                            acks
                        );

                        ExpandedNodeId extExpandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
                            false,            //Server Index Specified
                            new NodeIdFourByte((short) 0, Integer.valueOf(publishRequest.getIdentifier())),
                            null,
                            null);

                        ExtensionObject extObject = new ExtensionObject(
                            extExpandedNodeId,
                            null,
                            publishRequest);

                        try {
                            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
                            ExtensionObjectIO.staticSerialize(buffer, extObject);

                            /*  Create Consumer for the response message, error and timeout to be sent to the Secure Channel */
                            Consumer<byte[]> consumer = opcuaResponse -> {
                                PublishResponse responseMessage = null;
                                ServiceFault serviceFault = null;
                                try {
                                    ExtensionObjectDefinition unknownExtensionObject = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, ByteOrder.LITTLE_ENDIAN), false).getBody();
                                    if (unknownExtensionObject instanceof PublishResponse) {
                                        responseMessage = (PublishResponse) unknownExtensionObject;
                                    } else {
                                        serviceFault = (ServiceFault) unknownExtensionObject;
                                        ResponseHeader header = (ResponseHeader) serviceFault.getResponseHeader();
                                        LOGGER.error("Subscription ServiceFault returned from server with error code,  '{}'", header.getServiceResult().toString());
                                        //plcSubscriber.onDisconnect(context);
                                    }
                                } catch (ParseException e) {
                                    LOGGER.error("Unable to parse the returned Subscription response");
                                    e.printStackTrace();
                                    plcSubscriber.onDisconnect(context);
                                }
                                if (serviceFault == null) {
                                    outstandingRequests.remove(((ResponseHeader) responseMessage.getResponseHeader()).getRequestHandle());

                                    for (long availableSequenceNumber : responseMessage.getAvailableSequenceNumbers()) {
                                        outstandingAcknowledgements.add(new SubscriptionAcknowledgement(this.subscriptionId, availableSequenceNumber));
                                    }

                                    for (ExtensionObject notificationMessage : ((NotificationMessage) responseMessage.getNotificationMessage()).getNotificationData()) {
                                        ExtensionObjectDefinition notification = notificationMessage.getBody();
                                        if (notification instanceof DataChangeNotification) {
                                            LOGGER.trace("Found a Data Change notification");
                                            ExtensionObjectDefinition[] items = ((DataChangeNotification) notification).getMonitoredItems();
                                            onSubscriptionValue(Arrays.stream(items).toArray(MonitoredItemNotification[]::new));
                                        } else {
                                            LOGGER.warn("Unsupported Notification type");
                                        }
                                    }
                                }
                            };

                            Consumer<TimeoutException> timeout = e -> {
                                LOGGER.error("Timeout while waiting for subscription response");
                                e.printStackTrace();
                                plcSubscriber.onDisconnect(context);
                            };

                            BiConsumer<OpcuaAPU, Throwable> error = (message, e) -> {
                                LOGGER.error("Error while waiting for subscription response");
                                e.printStackTrace();
                                plcSubscriber.onDisconnect(context);
                            };

                            outstandingRequests.add(requestHandle);
                            channel.submit(context, timeout, error, consumer, buffer);

                        } catch (ParseException e) {
                            LOGGER.warn("Unable to serialize subscription request");
                            e.printStackTrace();
                        }
                    }
                    /* Put the subscriber loop to sleep for the rest of the cycle. */
                    sleep(this.revisedCycleTime);
                }
                //Wait for any outstanding responses to arrive, using the request timeout length
                //sleep(this.revisedCycleTime * 10);
                complete = true;
            } catch (Exception e) {
                LOGGER.error("Failed :(");
                e.printStackTrace();
            }
            return null;
        });
        return;
    }


    /**
     * Stop the subscriber either on disconnect or on error
     * @return
     */
    public void stopSubscriber() {
        this.destroy.set(true);

        long requestHandle = channel.getRequestHandle();

        RequestHeader requestHeader = new RequestHeader(channel.getAuthenticationToken(),
            SecureChannel.getCurrentDateTime(),
            requestHandle,
            0L,
            OpcuaProtocolLogic.NULL_STRING,
            this.revisedCycleTime * 10,
            OpcuaProtocolLogic.NULL_EXTENSION_OBJECT);

        long[] subscriptions = new long[1];
        subscriptions[0] = subscriptionId;
        DeleteSubscriptionsRequest deleteSubscriptionrequest = new DeleteSubscriptionsRequest(requestHeader,
            1,
            subscriptions
            );

        ExpandedNodeId extExpandedNodeId = new ExpandedNodeId(false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte((short) 0, Integer.valueOf(deleteSubscriptionrequest.getIdentifier())),
            null,
            null);

        ExtensionObject extObject = new ExtensionObject(
            extExpandedNodeId,
            null,
            deleteSubscriptionrequest);

        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(extObject.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            ExtensionObjectIO.staticSerialize(buffer, extObject);

            //  Create Consumer for the response message, error and timeout to be sent to the Secure Channel
            Consumer<byte[]> consumer = opcuaResponse -> {
                DeleteSubscriptionsResponse responseMessage = null;
                try {
                    ExtensionObjectDefinition unknownExtensionObject = ExtensionObjectIO.staticParse(new ReadBufferByteBased(opcuaResponse, ByteOrder.LITTLE_ENDIAN), false).getBody();
                    if (unknownExtensionObject instanceof DeleteSubscriptionsResponse) {
                        responseMessage = (DeleteSubscriptionsResponse) unknownExtensionObject;
                    } else {
                        ServiceFault serviceFault = (ServiceFault) unknownExtensionObject;
                        ResponseHeader header = (ResponseHeader) serviceFault.getResponseHeader();
                        LOGGER.error("Fault when deleteing Subscription ServiceFault return from server with error code,  '{}'", header.getServiceResult().toString());
                    }
                } catch (ParseException e) {
                    LOGGER.error("Unable to parse the returned Delete Subscriptions Response");
                    e.printStackTrace();
                }
            };

            Consumer<TimeoutException> timeout = e -> {
                LOGGER.error("Timeout while waiting for delete subscription response");
                e.printStackTrace();
                plcSubscriber.onDisconnect(context);
            };

            BiConsumer<OpcuaAPU, Throwable> error = (message, e) -> {
                LOGGER.error("Error while waiting for delete subscription response");
                e.printStackTrace();
                plcSubscriber.onDisconnect(context);
            };

            channel.submit(context, timeout, error, consumer, buffer);
        } catch (ParseException e) {
            LOGGER.warn("Unable to serialize subscription request");
            e.printStackTrace();
        }

        sleep(500);
        plcSubscriber.removeSubscription(subscriptionId);
    }

    /**
     * Receive the returned values from the OPCUA server and format it so that it can be received by the PLC4X client.
     * @param values - array of data values to be sent to the client.
     */
    private void onSubscriptionValue(MonitoredItemNotification[] values) {
        LinkedHashSet<String> fieldList = new LinkedHashSet<>();
        DataValue[] dataValues = new DataValue[values.length];
        int i = 0;
        for (MonitoredItemNotification value : values) {
            fieldList.add(fieldNames.get((int) value.getClientHandle() - 1));
            dataValues[i] = value.getValue();
            i++;
        }
        Map<String, ResponseItem<PlcValue>> fields = plcSubscriber.readResponse(fieldList, dataValues);
        final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(), fields);

        consumers.forEach(plcSubscriptionEventConsumer -> {
            plcSubscriptionEventConsumer.accept(event);
        });
    }

    /**
     * Registers a new Consumer, this allows multiple PLC4X consumers to use the same subscription.
     * @param consumer - Consumer to be used to send any returned values.
     * @return PlcConsumerRegistration - return the important information back to the client.
     */
    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        LOGGER.info("Registering a new OPCUA subscription consumer");
        consumers.add(consumer);
        return new DefaultPlcConsumerRegistration(plcSubscriber, consumer, this);
    }

    /**
     * Given an PLC4X OpcuaField generate the OPC UA Node Id
     * @param field - The PLC4X OpcuaField, this is the field generated from the OpcuaField class from the parsed field string.
     * @return NodeId - Returns an OPC UA Node Id which can be sent over the wire.
     */
    private NodeId generateNodeId(OpcuaField field) {
        NodeId nodeId = null;
        if (field.getIdentifierType() == OpcuaIdentifierType.BINARY_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdTwoByte(Short.valueOf(field.getIdentifier())));
        } else if (field.getIdentifierType() == OpcuaIdentifierType.NUMBER_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdNumeric((short) field.getNamespace(), Long.valueOf(field.getIdentifier())));
        } else if (field.getIdentifierType() == OpcuaIdentifierType.GUID_IDENTIFIER) {
            UUID guid = UUID.fromString(field.getIdentifier());
            byte[] guidBytes = new byte[16];
            System.arraycopy(guid.getMostSignificantBits(), 0, guidBytes, 0, 8);
            System.arraycopy(guid.getLeastSignificantBits(), 0, guidBytes, 8, 8);
            nodeId = new NodeId(new NodeIdGuid((short) field.getNamespace(), guidBytes));
        } else if (field.getIdentifierType() == OpcuaIdentifierType.STRING_IDENTIFIER) {
            nodeId = new NodeId(new NodeIdString((short) field.getNamespace(), new PascalString(field.getIdentifier())));
        }
        return nodeId;
    }



}
