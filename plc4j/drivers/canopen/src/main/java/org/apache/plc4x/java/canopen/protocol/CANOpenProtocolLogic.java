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
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.can.adapter.Plc4xCANProtocolBase;
import org.apache.plc4x.java.canopen.transport.CANOpenAbortException;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.api.conversation.canopen.CANConversation;
import org.apache.plc4x.java.canopen.api.conversation.canopen.SDODownloadConversation;
import org.apache.plc4x.java.canopen.api.conversation.canopen.SDOUploadConversation;
import org.apache.plc4x.java.canopen.configuration.CANOpenConfiguration;
import org.apache.plc4x.java.canopen.context.CANOpenDriverContext;
import org.apache.plc4x.java.canopen.field.CANOpenField;
import org.apache.plc4x.java.canopen.field.CANOpenHeartbeatField;
import org.apache.plc4x.java.canopen.field.CANOpenNMTField;
import org.apache.plc4x.java.canopen.field.CANOpenPDOField;
import org.apache.plc4x.java.canopen.field.CANOpenSDOField;
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
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
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
    private Logger logger = LoggerFactory.getLogger(CANOpenProtocolLogic.class);

    private CANOpenConfiguration configuration;
    private RequestTransactionManager tm;
    private Timer heartbeat;
    private CANOpenDriverContext canContext;
    private CANConversation conversation;

    private Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

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
        if (writeRequest.getFieldNames().size() != 1) {
            response.completeExceptionally(new IllegalArgumentException("You can write only one field at the time"));
            return response;
        }

        PlcField field = writeRequest.getFields().get(0);
        if (!(field instanceof CANOpenField)) {
            response.completeExceptionally(new IllegalArgumentException("Only CANOpenField instances are supported"));
            return response;
        }

        if (field instanceof CANOpenSDOField) {
            writeInternally((DefaultPlcWriteRequest) writeRequest, (CANOpenSDOField) field, response);
            return response;
        }
        if (field instanceof CANOpenPDOField) {
            writeInternally((DefaultPlcWriteRequest) writeRequest, (CANOpenPDOField) field, response);
            return response;
        }

        response.completeExceptionally(new IllegalArgumentException("Only CANOpenSDOField instances are supported"));
        return response;
    }

    private void writeInternally(DefaultPlcWriteRequest writeRequest, CANOpenSDOField field, CompletableFuture<PlcWriteResponse> response) {
        final RequestTransactionManager.RequestTransaction transaction = tm.startRequest();

        String fieldName = writeRequest.getFieldNames().iterator().next();

        CompletableFuture<PlcResponseCode> callback = new CompletableFuture<>();
        callback.whenComplete((code, error) -> {
            if (error != null) {
                if (error instanceof CANOpenAbortException) {
                    response.complete(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(fieldName, PlcResponseCode.REMOTE_ERROR)));
                } else {
                    response.complete(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(fieldName, PlcResponseCode.INTERNAL_ERROR)));
                }
                transaction.endRequest();
                return;
            }
            response.complete(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(fieldName, code)));
            transaction.endRequest();
        });

        PlcValue writeValue = writeRequest.getPlcValues().get(0);
        SDODownloadConversation download = new SDODownloadConversation(conversation, field.getNodeId(), field.getAnswerNodeId(),
            new IndexAddress(field.getIndex(), field.getSubIndex()), writeValue, field.getCanOpenDataType());
        transaction.submit(() -> download.execute(callback));
    }

    private void writeInternally(DefaultPlcWriteRequest writeRequest, CANOpenPDOField field, CompletableFuture<PlcWriteResponse> response) {
        PlcValue writeValue = writeRequest.getPlcValues().get(0);

        try {
            String fieldName = writeRequest.getFieldNames().iterator().next();

            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(DataItem.getLengthInBytes(writeValue, field.getCanOpenDataType(), writeValue.getLength()), ByteOrder.LITTLE_ENDIAN);
            DataItem.staticSerialize(writeBuffer, writeValue, field.getCanOpenDataType(), writeValue.getLength(), ByteOrder.LITTLE_ENDIAN);
            final CANOpenPDOPayload payload = new CANOpenPDOPayload(new CANOpenPDO(writeBuffer.getData()));
            context.sendToWire(new CANOpenFrame((short) field.getNodeId(), field.getService(), payload));
            response.complete(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(fieldName, PlcResponseCode.OK)));
        } catch (Exception e) {
            response.completeExceptionally(e);
        }
    }

    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> response = new CompletableFuture<>();
        if (readRequest.getFieldNames().size() != 1) {
            response.completeExceptionally(new IllegalArgumentException("SDO requires single field to be read"));
            return response;
        }

        PlcField field = readRequest.getFields().get(0);
        if (!(field instanceof CANOpenField)) {
            response.completeExceptionally(new IllegalArgumentException("Only CANOpenField instances are supported"));
            return response;
        }

        if (!(field instanceof CANOpenSDOField)) {
            response.completeExceptionally(new IllegalArgumentException("Only CANOpenSDOField instances are supported"));
            return response;
        }

        readInternally(readRequest, (CANOpenSDOField) field, response);
        return response;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest request) {
        DefaultPlcSubscriptionRequest rq = (DefaultPlcSubscriptionRequest) request;

        Map<String, ResponseItem<PlcSubscriptionHandle>> answers = new LinkedHashMap<>();
        DefaultPlcSubscriptionResponse response = new DefaultPlcSubscriptionResponse(rq, answers);

        for (String key : rq.getFieldNames()) {
            DefaultPlcSubscriptionField subscription = (DefaultPlcSubscriptionField) rq.getField(key);
            if (subscription.getPlcSubscriptionType() != PlcSubscriptionType.EVENT) {
                answers.put(key, new ResponseItem<>(PlcResponseCode.UNSUPPORTED, null));
            } else if ((subscription.getPlcField() instanceof CANOpenPDOField)) {
                answers.put(key, new ResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, (CANOpenPDOField) subscription.getPlcField())
                ));
            } else if ((subscription.getPlcField() instanceof CANOpenNMTField)) {
                answers.put(key, new ResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, (CANOpenNMTField) subscription.getPlcField())
                ));
            } else if ((subscription.getPlcField() instanceof CANOpenHeartbeatField)) {
                answers.put(key, new ResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, (CANOpenHeartbeatField) subscription.getPlcField())
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

    private void readInternally(PlcReadRequest readRequest, CANOpenSDOField field, CompletableFuture<PlcReadResponse> response) {
        String fieldName = readRequest.getFieldNames().iterator().next();

        final RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        CompletableFuture<PlcValue> callback = new CompletableFuture<>();
        callback.whenComplete((value, error) -> {
            if (error != null) {
                Map<String, ResponseItem<PlcValue>> fields = new HashMap<>();
                if (error instanceof CANOpenAbortException) {
                    fields.put(fieldName, new ResponseItem<>(PlcResponseCode.REMOTE_ERROR, new PlcLINT(((CANOpenAbortException) error).getAbortCode())));
                } else {
                    fields.put(fieldName, new ResponseItem<>(PlcResponseCode.REMOTE_ERROR, null));
                }
                response.complete(new DefaultPlcReadResponse(readRequest, fields));
                transaction.endRequest();

                return;
            }

            Map<String, ResponseItem<PlcValue>> fields = new HashMap<>();
            fields.put(fieldName, new ResponseItem<>(PlcResponseCode.OK, value));
            response.complete(new DefaultPlcReadResponse(readRequest, fields));
            transaction.endRequest();
        });

        SDOUploadConversation upload = new SDOUploadConversation(conversation, field.getNodeId(), field.getAnswerNodeId(), new IndexAddress(field.getIndex(), field.getSubIndex()), field.getCanOpenDataType());
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

                        CANOpenPDOField field = (CANOpenPDOField) handle.getField();
                        byte[] data = ((CANOpenPDOPayload) payload).getPdo().getData();
                        try {
                            PlcValue value = DataItem.staticParse(new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN), field.getCanOpenDataType(), data.length);
                            DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(
                                Instant.now(),
                                Collections.singletonMap(
                                    handle.getName(),
                                    new ResponseItem<>(PlcResponseCode.OK, value)
                                )
                            );
                            consumer.accept(event);
                        } catch (ParseException e) {
                            logger.warn("Could not parse data to desired type: {}", field.getCanOpenDataType(), e);
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
                        Map<String, PlcValue> fields = new HashMap<>();
                        fields.put("state", new PlcUSINT(state.getValue()));
                        fields.put("node", new PlcUSINT(nodeId));
                        PlcStruct struct = new PlcStruct(fields);
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
                        Map<String, PlcValue> fields = new HashMap<>();
                        fields.put("state", new PlcUSINT(state.getValue()));
                        fields.put("node", new PlcUSINT(nodeId));
                        PlcStruct struct = new PlcStruct(fields);
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
