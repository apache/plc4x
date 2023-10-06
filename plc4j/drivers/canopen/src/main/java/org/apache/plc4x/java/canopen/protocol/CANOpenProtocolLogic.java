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
package org.apache.plc4x.java.canopen.protocol;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.can.adapter.Plc4xCANProtocolBase;
import org.apache.plc4x.java.canopen.tag.*;
import org.apache.plc4x.java.canopen.transport.CANOpenAbortException;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.api.conversation.canopen.CANConversation;
import org.apache.plc4x.java.canopen.api.conversation.canopen.SDODownloadConversation;
import org.apache.plc4x.java.canopen.api.conversation.canopen.SDOUploadConversation;
import org.apache.plc4x.java.canopen.configuration.CANOpenConfiguration;
import org.apache.plc4x.java.canopen.context.CANOpenDriverContext;
import org.apache.plc4x.java.canopen.conversation.CANTransportConversation;
import org.apache.plc4x.java.canopen.readwrite.CANOpenHeartbeatPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenNetworkPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPDO;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPDOPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.IndexAddress;
import org.apache.plc4x.java.canopen.readwrite.DataItem;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.canopen.readwrite.NMTState;
import org.apache.plc4x.java.canopen.readwrite.NMTStateRequest;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.spi.values.PlcUSINT;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.Instant;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

public class CANOpenProtocolLogic extends Plc4xCANProtocolBase<CANOpenFrame>
    implements HasConfiguration<CANOpenConfiguration>, PlcSubscriber {

    private static final Duration REQUEST_TIMEOUT = Duration.ofSeconds(10L);
    private final Logger logger = LoggerFactory.getLogger(CANOpenProtocolLogic.class);

    private CANOpenConfiguration configuration;
    private RequestTransactionManager tm;
    private Timer heartbeat;
    private CANOpenDriverContext canContext;
    private CANConversation conversation;

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    @Override
    public void setConfiguration(CANOpenConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);
        this.canContext = (CANOpenDriverContext) driverContext;

        // Initialize Transaction Manager.
        // Until the number of concurrent requests is successfully negotiated we set it to a
        // maximum of only one request being able to be sent at a time. During the login process
        // No concurrent requests can be sent anyway. It will be updated when receiving the
        // S7ParameterSetupCommunication response.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void onConnect(ConversationContext<CANOpenFrame> context) {
        try {
            if (configuration.isHeartbeat()) {
                context.sendToWire(createFrame(new CANOpenHeartbeatPayload(NMTState.BOOTED_UP)));

                this.heartbeat = new Timer("can-heartbeat");
                this.heartbeat.scheduleAtFixedRate(new TimerTask() {
                    @Override
                    public void run() {
                        try {
                            context.sendToWire(createFrame(new CANOpenHeartbeatPayload(NMTState.OPERATIONAL)));
                        } catch (ParseException e) {
                            throw new PlcRuntimeException(e);
                        }
                    }
                }, 10000, 10000);
            }
            context.fireConnected();
        } catch (ParseException e) {
            throw new PlcRuntimeException(e);
        }
    }

    @Override
    public void setContext(ConversationContext<CANOpenFrame> context) {
        super.setContext(context);
        this.conversation = new CANTransportConversation(configuration.getNodeId(), context, configuration.getRequestTimeout());
    }

    private CANOpenFrame createFrame(CANOpenHeartbeatPayload state) throws ParseException {
        return new CANOpenFrame((short) configuration.getNodeId(), CANOpenService.HEARTBEAT, state);
    }

    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> response = new CompletableFuture<>();
        if (writeRequest.getTagNames().size() != 1) {
            response.completeExceptionally(new IllegalArgumentException("You can write only one tag at the time"));
            return response;
        }

        PlcTag tag = writeRequest.getTags().get(0);
        if (!(tag instanceof CANOpenTag)) {
            response.completeExceptionally(new IllegalArgumentException("Only CANOpenTag instances are supported"));
            return response;
        }

        if (tag instanceof CANOpenSDOTag) {
            writeInternally((DefaultPlcWriteRequest) writeRequest, (CANOpenSDOTag) tag, response);
            return response;
        }
        if (tag instanceof CANOpenPDOTag) {
            writeInternally((DefaultPlcWriteRequest) writeRequest, (CANOpenPDOTag) tag, response);
            return response;
        }

        response.completeExceptionally(new IllegalArgumentException("Only CANOpenSDOTag instances are supported"));
        return response;
    }

    private void writeInternally(DefaultPlcWriteRequest writeRequest, CANOpenSDOTag tag, CompletableFuture<PlcWriteResponse> response) {
        final RequestTransactionManager.RequestTransaction transaction = tm.startRequest();

        String tagName = writeRequest.getTagNames().iterator().next();

        CompletableFuture<PlcResponseCode> callback = new CompletableFuture<>();
        callback.whenComplete((code, error) -> {
            if (error != null) {
                if (error instanceof CANOpenAbortException) {
                    response.complete(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(tagName, PlcResponseCode.REMOTE_ERROR)));
                } else {
                    response.complete(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(tagName, PlcResponseCode.INTERNAL_ERROR)));
                }
                transaction.endRequest();
                return;
            }
            response.complete(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(tagName, code)));
            transaction.endRequest();
        });

        PlcValue writeValue = writeRequest.getPlcValues().get(0);
        SDODownloadConversation download = new SDODownloadConversation(conversation, tag.getNodeId(), tag.getAnswerNodeId(),
            new IndexAddress(tag.getIndex(), tag.getSubIndex()), writeValue, tag.getCanOpenDataType());
        transaction.submit(() -> download.execute(callback));
    }

    private void writeInternally(DefaultPlcWriteRequest writeRequest, CANOpenPDOTag tag, CompletableFuture<PlcWriteResponse> response) {
        PlcValue writeValue = writeRequest.getPlcValues().get(0);

        try {
            String tagName = writeRequest.getTagNames().iterator().next();

            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(DataItem.getLengthInBytes(writeValue, tag.getCanOpenDataType(), writeValue.getLength()), ByteOrder.LITTLE_ENDIAN);
            DataItem.staticSerialize(writeBuffer, writeValue, tag.getCanOpenDataType(), writeValue.getLength(), ByteOrder.LITTLE_ENDIAN);
            final CANOpenPDOPayload payload = new CANOpenPDOPayload(new CANOpenPDO(writeBuffer.getData()));
            context.sendToWire(new CANOpenFrame((short) tag.getNodeId(), tag.getService(), payload));
            response.complete(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(tagName, PlcResponseCode.OK)));
        } catch (Exception e) {
            response.completeExceptionally(e);
        }
    }

    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> response = new CompletableFuture<>();
        if (readRequest.getTagNames().size() != 1) {
            response.completeExceptionally(new IllegalArgumentException("SDO requires single tag to be read"));
            return response;
        }

        PlcTag tag = readRequest.getTags().get(0);
        if (!(tag instanceof CANOpenTag)) {
            response.completeExceptionally(new IllegalArgumentException("Only CANOpenTag instances are supported"));
            return response;
        }

        if (!(tag instanceof CANOpenSDOTag)) {
            response.completeExceptionally(new IllegalArgumentException("Only CANOpenSDOTag instances are supported"));
            return response;
        }

        readInternally(readRequest, (CANOpenSDOTag) tag, response);
        return response;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest request) {
        DefaultPlcSubscriptionRequest rq = (DefaultPlcSubscriptionRequest) request;

        Map<String, ResponseItem<PlcSubscriptionHandle>> answers = new LinkedHashMap<>();
        DefaultPlcSubscriptionResponse response = new DefaultPlcSubscriptionResponse(rq, answers);

        for (String key : rq.getTagNames()) {
            DefaultPlcSubscriptionTag subscription = (DefaultPlcSubscriptionTag) rq.getTag(key);
            if (subscription.getPlcSubscriptionType() != PlcSubscriptionType.EVENT) {
                answers.put(key, new ResponseItem<>(PlcResponseCode.UNSUPPORTED, null));
            } else if ((subscription.getTag() instanceof CANOpenPDOTag)) {
                answers.put(key, new ResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, (CANOpenPDOTag) subscription.getTag())
                ));
            } else if ((subscription.getTag() instanceof CANOpenNMTTag)) {
                answers.put(key, new ResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, (CANOpenNMTTag) subscription.getTag())
                ));
            } else if ((subscription.getTag() instanceof CANOpenHeartbeatTag)) {
                answers.put(key, new ResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, (CANOpenHeartbeatTag) subscription.getTag())
                ));
            } else {
                answers.put(key, new ResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            }
        }

        return CompletableFuture.completedFuture(response);
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest request) {
        List<PlcSubscriptionHandle> handles = request.getSubscriptionHandles();

        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            entry.getKey().getSubscriptionHandles().removeAll(handles);
        }

        return CompletableFuture.completedFuture(new DefaultPlcUnsubscriptionResponse(request));
    }

    private void readInternally(PlcReadRequest readRequest, CANOpenSDOTag tag, CompletableFuture<PlcReadResponse> response) {
        String tagName = readRequest.getTagNames().iterator().next();

        final RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        CompletableFuture<PlcValue> callback = new CompletableFuture<>();
        callback.whenComplete((value, error) -> {
            if (error != null) {
                Map<String, ResponseItem<PlcValue>> tags = new HashMap<>();
                if (error instanceof CANOpenAbortException) {
                    tags.put(tagName, new ResponseItem<>(PlcResponseCode.REMOTE_ERROR, new PlcLINT(((CANOpenAbortException) error).getAbortCode())));
                } else {
                    tags.put(tagName, new ResponseItem<>(PlcResponseCode.REMOTE_ERROR, null));
                }
                response.complete(new DefaultPlcReadResponse(readRequest, tags));
                transaction.endRequest();

                return;
            }

            Map<String, ResponseItem<PlcValue>> tags = new HashMap<>();
            tags.put(tagName, new ResponseItem<>(PlcResponseCode.OK, value));
            response.complete(new DefaultPlcReadResponse(readRequest, tags));
            transaction.endRequest();
        });

        SDOUploadConversation upload = new SDOUploadConversation(conversation, tag.getNodeId(), tag.getAnswerNodeId(), new IndexAddress(tag.getIndex(), tag.getSubIndex()), tag.getCanOpenDataType());
        transaction.submit(() -> upload.execute(callback));
    }

    @Override
    public void decode(ConversationContext<CANOpenFrame> context, CANOpenFrame msg) throws Exception {
        int nodeId = msg.getNodeId();
        CANOpenService service = msg.getService();
        CANOpenPayload payload = msg.getPayload();

        if (service != null && nodeId != this.configuration.getNodeId()) {
            if (service.getPdo() && payload instanceof CANOpenPDOPayload) {
                publishEvent(service, nodeId, payload);
            } else if (service == CANOpenService.HEARTBEAT && payload instanceof CANOpenHeartbeatPayload) {
                publishEvent(service, nodeId, payload);
            } else {
                logger.debug("Decoded CANOpen {} from {}, message {}", service, nodeId, payload);
            }
        }

//        int identifier = msg.getIdentifier();
//        CANOpenService service = CANOpenService.valueOf((byte) (identifier >> 7));
//        if (service != null) {
//            ReadBuffer buffer = new ReadBuffer(msg.getData());
//            CANOpenPayload payload = CANOpenPayloadIO.staticParse(buffer, service);
//
//
//        }
    }

    private void publishEvent(CANOpenService service, int nodeId, CANOpenPayload payload) {
        DefaultPlcSubscriptionHandle dispatchedHandle = null;
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            DefaultPlcConsumerRegistration registration = entry.getKey();
            Consumer<PlcSubscriptionEvent> consumer = entry.getValue();

            for (PlcSubscriptionHandle handler : registration.getSubscriptionHandles()) {
                CANOpenSubscriptionHandle handle = (CANOpenSubscriptionHandle) handler;
                if (payload instanceof CANOpenPDOPayload) {

                    if (handle.matches(service, nodeId)) {
                        logger.trace("Dispatching notification {} for node {} to {}", service, nodeId, handle);
                        dispatchedHandle = handle;

                        CANOpenPDOTag tag = (CANOpenPDOTag) handle.getTag();
                        byte[] data = ((CANOpenPDOPayload) payload).getPdo().getData();
                        try {
                            PlcValue value = DataItem.staticParse(new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN), tag.getCanOpenDataType(), data.length);
                            DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(
                                Instant.now(),
                                Collections.singletonMap(
                                    handle.getName(),
                                    new ResponseItem<>(PlcResponseCode.OK, value)
                                )
                            );
                            consumer.accept(event);
                        } catch (ParseException e) {
                            logger.warn("Could not parse data to desired type: {}", tag.getCanOpenDataType(), e);
                            DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(
                                Instant.now(),
                                Collections.singletonMap(
                                    handle.getName(),
                                    new ResponseItem<>(PlcResponseCode.INVALID_DATA, new PlcNull())
                                )
                            );
                            consumer.accept(event);
                        }
                    }
                } else if (payload instanceof CANOpenHeartbeatPayload) {
                    if (handle.matches(service, nodeId)) {
                        logger.trace("Dispatching notification {} for node {} to {}", service, nodeId, handle);
                        dispatchedHandle = handle;

                        final NMTState state = ((CANOpenHeartbeatPayload) payload).getState();
                        Map<String, PlcValue> tags = new HashMap<>();
                        tags.put("state", new PlcUSINT(state.getValue()));
                        tags.put("node", new PlcUSINT(nodeId));
                        PlcStruct struct = new PlcStruct(tags);
                        DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(
                            Instant.now(),
                            Collections.singletonMap(
                                handle.getName(),
                                new ResponseItem<>(PlcResponseCode.OK, struct)
                            )
                        );
                        consumer.accept(event);
                    }
                } else if (payload instanceof CANOpenNetworkPayload) {
                    if (handle.matches(service, nodeId)) {
                        logger.trace("Dispatching notification {} for node {} to {}", service, nodeId, handle);
                        dispatchedHandle = handle;

                        final NMTStateRequest state = ((CANOpenNetworkPayload) payload).getRequest();
                        Map<String, PlcValue> tags = new HashMap<>();
                        tags.put("state", new PlcUSINT(state.getValue()));
                        tags.put("node", new PlcUSINT(nodeId));
                        PlcStruct struct = new PlcStruct(tags);
                        DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(
                            Instant.now(),
                            Collections.singletonMap(
                                handle.getName(),
                                new ResponseItem<>(PlcResponseCode.OK, struct)
                            )
                        );
                        consumer.accept(event);
                    }
                }
            }
        }

        if (dispatchedHandle == null) {
            logger.trace("Could not find subscription matching {} and node {}", service, nodeId);
        }
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        final DefaultPlcConsumerRegistration consumerRegistration = new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new DefaultPlcSubscriptionHandle[0]));
        consumers.put(consumerRegistration, consumer);
        return consumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        consumers.remove(registration);
    }

    @Override
    public void close(ConversationContext<CANOpenFrame> context) {

    }

    @Override
    public void onDisconnect(ConversationContext<CANOpenFrame> context) {
        if (this.heartbeat != null) {
            this.heartbeat.cancel();
            this.heartbeat = null;
        }
    }

}
