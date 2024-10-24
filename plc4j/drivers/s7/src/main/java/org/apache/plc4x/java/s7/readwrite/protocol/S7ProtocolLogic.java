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
package org.apache.plc4x.java.s7.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.events.S7CyclicEvent;
import org.apache.plc4x.java.s7.events.S7Event;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.s7.readwrite.tag.*;
import org.apache.plc4x.java.s7.readwrite.types.S7SubscriptionType;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcSubscriptionHandle;
import org.apache.plc4x.java.s7.utils.S7ParamErrorCode;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagValueItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcTagItem;
import org.apache.plc4x.java.spi.messages.utils.PlcTagValueItem;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.time.Duration;
import java.time.LocalDateTime;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.s7.events.S7AlarmEvent;
import org.apache.plc4x.java.s7.events.S7ModeEvent;
import org.apache.plc4x.java.s7.events.S7SysEvent;
import org.apache.plc4x.java.s7.events.S7UserEvent;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcSubscriptionRequest;

/**
 * The S7 Protocol states that there can not be more then {min(maxAmqCaller, maxAmqCallee} "ongoing" requests.
 * So we need to limit those.
 * Thus, each request goes to a Work Queue and this Queue ensures, that only 3 are open at the same time.
 */
public class S7ProtocolLogic extends Plc4xProtocolBase<TPKTPacket> {

    private static final Logger logger = LoggerFactory.getLogger(S7ProtocolLogic.class);

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);
    private final AtomicInteger tpduGenerator = new AtomicInteger(10);

    /*
     * Task group for managing connection redundancy.
     */
    private final ExecutorService clientExecutorService = Executors.newFixedThreadPool(4, new BasicThreadFactory.Builder()
        .namingPattern("plc4x-app-thread-%d")
        .daemon(true)
        .priority(Thread.MAX_PRIORITY)
        .build());

    /*
     * Take into account that the size of this buffer depends on the final device.
     * S7-300 goes from 20 to 300 and for S7-400 it goes from 300 to 10000.
     * Depending on the configuration of the alarm system, a large number of
     * them should be expected when starting the connection.
     * (Examples of this are PCS7 and Braumat).
     * Alarm filtering, ack, etc. must be performed by the client application.
     */
    private final BlockingQueue<S7Event> eventQueue = new ArrayBlockingQueue<>(1024);
    private final S7ProtocolEventLogic eventLogic = new S7ProtocolEventLogic(eventQueue);
    private final S7PlcSubscriptionHandle modeHandle = new S7PlcSubscriptionHandle(EventType.MODE, eventLogic);
    private final S7PlcSubscriptionHandle sysHandle = new S7PlcSubscriptionHandle(EventType.SYS, eventLogic);
    private final S7PlcSubscriptionHandle usrHandle = new S7PlcSubscriptionHandle(EventType.USR, eventLogic);
    private final S7PlcSubscriptionHandle almHandle = new S7PlcSubscriptionHandle(EventType.ALM, eventLogic);
    //private final S7PlcSubscriptionHandle cycHandle = new S7PlcSubscriptionHandle(EventType.CYC, EventLogic);    

    /*
     * This array stores the cyclic subscription requests between the driver
     * and the PLC. The purpose is to document the tags associated with the
     * request. Each subscription uses a 'JobID' that is managed by the PLC and
     * obtained from the response to the request. In the following,
     * the values sent PUSH from the PLC to the driver refer to this JobID.
     */
    private final Map<Short, PlcSubscriptionRequest> cycRequests = new HashMap<>();
    
    /*
    * This data structure stores the last value associated with a cyclic 
    * subscription request. In each event received, the values of the internal 
    * PlcValue are compared and if any of them are different, the new value is 
    * transferred to the event stack and the value is updated in this HashMap.
    */
    private final Map<Short, S7CyclicEvent> cycChangeValueEvents = new HashMap<>();

    private S7DriverContext s7DriverContext;
    private RequestTransactionManager tm;

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);
        this.s7DriverContext = (S7DriverContext) driverContext;

        // Initialize Transaction Manager.
        // Until the number of concurrent requests is successfully negotiated we set it to a
        // maximum of only one request being able to be sent at a time. During the login process
        // No concurrent requests can be sent anyway. It will be updated when receiving the
        // S7ParameterSetupCommunication response.
        this.tm = new RequestTransactionManager(1);
        eventLogic.start();
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new S7PlcTagHandler();
    }

    @Override
    public void close(ConversationContext<TPKTPacket> context) {
        // TODO: Find out how to close this prior to Java 19
        //clientExecutorService.close();
        tm.shutdown();
        eventLogic.stop();
        // TODO Implement Closing on Protocol Level
    }

    @Override
    public void onConnect(ConversationContext<TPKTPacket> context) {
        if (context.isPassive()) {
            logger.info("S7 Driver running in PASSIVE mode.");
            s7DriverContext.setPassiveMode(true);
            // No login required, just confirm that we're connected.
            context.fireConnected();
            return;
        }

        //Set feature for all handlers in the pipeline from
        //the driver configuration.
        setChannelFeatures();

        // Only the TCP transport supports login.
        logger.info("S7 Driver running in ACTIVE mode.");
        logger.debug("Sending COTP Connection Request");
        // Open the session on ISO Transport Protocol first.
        TPKTPacket packet = new TPKTPacket(createCOTPConnectionRequest(
            s7DriverContext.getCalledTsapId(), s7DriverContext.getCallingTsapId(), s7DriverContext.getCotpTpduSize()));

        context.sendRequest(packet)
            .onTimeout(e -> {
                logger.info("Timeout during Connection establishing, closing channel...");
                // TODO: We're saying that we're closing the channel, but not closing the channel ... sure, this is what we want?
                //context.getChannel().close();
            })
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .unwrap(TPKTPacket::getPayload)
            .only(COTPPacketConnectionResponse.class)
            .handle(cotpPacketConnectionResponse -> {
                logger.debug("Got COTP Connection Response");
                logger.debug("Sending S7 Connection Request");
                context.sendRequest(createS7ConnectionRequest(cotpPacketConnectionResponse))
                    .onTimeout(e -> {
                        logger.warn("Timeout during Connection establishing, closing channel...");
                        context.getChannel().close();
                    })
                    .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                    .unwrap(TPKTPacket::getPayload)
                    .only(COTPPacketData.class)
                    .unwrap(COTPPacket::getPayload)
                    .only(S7MessageResponseData.class)
                    .unwrap(S7Message::getParameter)
                    .only(S7ParameterSetupCommunication.class)
                    .handle(setupCommunication -> {
                        logger.debug("Got S7 Connection Response");
                        // Save some data from the response.
                        s7DriverContext.setMaxAmqCaller(setupCommunication.getMaxAmqCaller());
                        s7DriverContext.setMaxAmqCallee(setupCommunication.getMaxAmqCallee());
                        s7DriverContext.setPduSize(setupCommunication.getPduLength());

                        // Update the number of concurrent requests to the negotiated number.
                        // I have never seen anything else than equal values for caller and
                        // callee, but if they were different, we're only limiting the outgoing
                        // requests.
                        tm.setNumberOfConcurrentRequests(s7DriverContext.getMaxAmqCallee());

                        // If the controller type is explicitly set, were finished with the login
                        // process. If it's set to ANY, we have to query the serial number information
                        // in order to detect the type of PLC.
                        if (s7DriverContext.getControllerType() != ControllerType.ANY) {
                            // Send an event that connection setup is complete.
                            context.fireConnected();
                            return;
                        }

                        // Prepare a message to request the remote to identify itself.
                        logger.debug("Sending S7 Identification Request");
                        TPKTPacket tpktPacket = createIdentifyRemoteMessage();
                        context.sendRequest(tpktPacket)
                            .onTimeout(e -> {
                                logger.warn("Timeout during Connection establishing, closing channel...");
                                context.getChannel().close();
                            })
                            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                            .unwrap(TPKTPacket::getPayload)
                            .only(COTPPacketData.class)
                            .unwrap(COTPPacketData::getPayload)
                            .only(S7MessageUserData.class)
                            .unwrap(S7MessageUserData::getPayload)
                            .only(S7PayloadUserData.class)
                            .handle(payloadUserData -> {
                                logger.debug("Got S7 Identification Response");
                                extractControllerTypeAndFireConnected(context, payloadUserData);
                            });
                    });
            });
    }

    /*
     * It performs the sequential and safe shutdown of the driver.
     * Completion of pending requests, executors and associated tasks.
     */
    @Override
    public void onDisconnect(ConversationContext<TPKTPacket> context) {
        // 1. Here we shut down the local task executor.
        clientExecutorService.shutdown();
        // 2. Performs the shutdown of the transaction executor.
        tm.shutdown();
        // 3. Finish the execution of the tasks for the handling of Events.
        eventLogic.stop();
        // 4. Executes the closing of the main channel.
        context.getChannel().close();
        // 5. Here is the stop of any task or state machine that is added.
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        // If we're not connected, just abort with an error.
        if (!isConnected()) {
            CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
            future.completeExceptionally(new PlcRuntimeException("Disconnected"));
            return future;
        }

        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        CompletableFuture<S7Message> responseFuture;
        if (request.getTagNames().stream().anyMatch(t -> request.getTag(t) instanceof S7SzlTag)) {
            // TODO: Is it correct, that there can only be one szl tag?
            S7SzlTag szlTag = (S7SzlTag) request.getTags().get(0);
            S7Message s7Message = new S7MessageUserData(getTpduId(),
                new S7ParameterUserData(List.of(
                    new S7ParameterUserDataItemCPUFunctions(
                        (short) 0x11, (byte) 0x4, (byte) 0x4, (short) 0x01, (short) 0x00,
                        null, null, null)
                )),
                new S7PayloadUserData(List.of(
                    new S7PayloadUserDataItemCpuFunctionReadSzlRequest(DataTransportErrorCode.OK,
                        DataTransportSize.OCTET_STRING,
                        0x04,
                        new SzlId(SzlModuleTypeClass.enumForValue((byte) ((szlTag.getSzlId() & 0xf000) >> 12)),
                            (byte) ((szlTag.getSzlId() & 0x0f00) >> 8),
                            SzlSublist.enumForValue((short) (szlTag.getSzlId() & 0x00ff))),
                        szlTag.getIndex())
                )));
            responseFuture = sendInternal(s7Message);
        } else if (request.getTagNames().stream().anyMatch(t -> request.getTag(t) instanceof S7AckTag)) {
            responseFuture = performAlarmAckRequest(request);
        } else if (request.getTagNames().stream().anyMatch(t -> request.getTag(t) instanceof S7ClkTag)) {
            responseFuture = performClkRequest(request);
        }

        // If the request contains at least one var-length string field, we need to get the real length first.
        else if (request.getTagNames().stream().anyMatch(t -> request.getTag(t) instanceof S7StringVarLengthTag)) {
            responseFuture = performVarLengthStringReadRequest(request);
        }

        // This is a "normal" read request.
        else {
            responseFuture = performOrdinaryReadRequest(request);
        }

        // Just send a single response and chain it as Response
        return toPlcReadResponse(readRequest, responseFuture);
    }

    /**
     * Maps the S7ReadResponse of a PlcReadRequest to a PlcReadResponse
     */
    private CompletableFuture<PlcReadResponse> toPlcReadResponse(PlcReadRequest readRequest, CompletableFuture<S7Message> responseFuture) {
        CompletableFuture<PlcReadResponse> clientFuture = new CompletableFuture<>();

        responseFuture.whenComplete((s7Message, throwable) -> {
            if (throwable != null) {
                clientFuture.completeExceptionally(new PlcProtocolException("Error reading", throwable));
            } else {
                try {
                    PlcReadResponse response = (PlcReadResponse) decodeReadResponse(s7Message, readRequest);
                    clientFuture.complete(response);
                } catch (Exception ex) {
                    logger.info(ex.toString());
                }
            }
        });

        return clientFuture;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        // If we're not connected, just abort with an error.
        if (!isConnected()) {
            CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
            future.completeExceptionally(new PlcRuntimeException("Disconnected"));
            return future;
        }

        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        CompletableFuture<S7Message> responseFuture = new CompletableFuture<>();
        // TODO: Write one or two lines on what happens here ... to me it looks as if there's at least on S7ClkTag, then all is handled by the writeClk method, but what happens if a request would contain mixed tag types?
        if (request.getTagNames().stream().anyMatch(t -> request.getTag(t) instanceof S7ClkTag)) {
            responseFuture = performClkSetRequest((DefaultPlcWriteRequest) writeRequest);
        }

        // If the list of tags contains at least one STRING/WSTRING element,
        // we need to check the sizes of the string fields in a first request.
        else if (request.getTagNames().stream().anyMatch(t -> request.getTag(t) instanceof S7StringVarLengthTag)) {
            responseFuture = performVarLengthStringWriteRequest((DefaultPlcWriteRequest) writeRequest);
        }

        // This is a request only containing ordinary tags
        else {
            responseFuture = performOrdinaryWriteRequest(request);
        }

        return toPlcWriteResponse(writeRequest, responseFuture);
    }

    private CompletableFuture<PlcWriteResponse> toPlcWriteResponse(PlcWriteRequest writeRequest, CompletableFuture<S7Message> responseFuture) {
        CompletableFuture<PlcWriteResponse> clientFuture = new CompletableFuture<>();

        responseFuture.whenComplete((s7Message, throwable) -> {
            if (throwable != null) {
                clientFuture.completeExceptionally(new PlcProtocolException("Error writing", throwable));
            } else {
                try {
                    PlcWriteResponse response = (PlcWriteResponse) decodeWriteResponse(s7Message, writeRequest);
                    clientFuture.complete(response);
                } catch (Exception ex) {
                    logger.info(ex.toString());
                }
            }
        });

        return clientFuture;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        if (!isConnected()) {
            CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
            future.completeExceptionally(new PlcRuntimeException("Disconnected"));
            return future;
        }
        if (!isFeatureSupported()) {
            CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
            future.completeExceptionally(new PlcRuntimeException("Not Supported"));
            return future;
        }

        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        HashMap<String, PlcSubscriptionResponse> valuesResponse = new HashMap<>();
        HashMap<String, CompletableFuture<S7Message>> futures = new HashMap<>();

        //Initialize multiple requests.
        CompletableFuture<PlcSubscriptionResponse> response = new CompletableFuture<>();
        subscriptionRequest.getTagNames().forEach(fieldName -> futures.put(fieldName, new CompletableFuture<>()));

        futures.put("DATA_", new CompletableFuture<>());

        S7PlcSubscriptionRequest request = (S7PlcSubscriptionRequest) subscriptionRequest;

        int tpduId = getTpduId();

        //The main task that runs the subscriptions.
        Thread t1 = new Thread(() -> {
            final DefaultPlcSubscriptionTag sf = (DefaultPlcSubscriptionTag) request.getTags().get(0);
            final S7SubscriptionTag tag = (S7SubscriptionTag) sf.getTag();

            S7Message s7Message = null;
            switch (tag.getTagType()) {
                case EVENT_SUBSCRIPTION:
                    s7Message = encodeEventSubscriptionRequest(request, tpduId);
                    break;
                case EVENT_UNSUBSCRIPTION:
                    //encodeEventUnSubscriptionRequest(request, parameterItems, payloadItems);
                    break;
                case ALARM_ACK:
                    //encodeAlarmAckRequest(request, parameterItems, payloadItems);
                    break;
                case ALARM_QUERY:
                    s7Message = encodeAlarmQueryRequest(request, tpduId);
                    break;
                case CYCLIC_SUBSCRIPTION:
                    s7Message = encodeCycledS7ANYSubscriptionRequest(request, tpduId);
                    break;
                case CYCLIC_DB_SUBSCRIPTION:
                    s7Message = encodeCycledDBREADSubscriptionRequest(request, tpduId);
                    break;
                case CYCLIC_UNSUBSCRIPTION:
                    //encodeCycledUnSubscriptionRequest(request, parameterItems, payloadItems);
                    break;
                default:
            }
            if (s7Message == null) {
                throw new PlcInvalidTagException("Unsupported tag of type: " + tag.getTagType());
            }

            TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null, s7Message, true, (byte) tpduId));

            // Start a new request-transaction (Is ended in the response-handler)
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> conversationContext.sendRequest(tpktPacket)
                .onTimeout(new TransactionErrorCallback<>(future, transaction))
                .onError(new TransactionErrorCallback<>(future, transaction))
                .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                .unwrap(TPKTPacket::getPayload)
                .only(COTPPacketData.class)
                .unwrap(COTPPacket::getPayload)
                .check(p -> p.getTpduReference() == tpduId)
                .handle(p -> {
                    try {
                        //future.complete(decodeEventSubscriptionRequest(tagName, p, subscriptionRequest));
                        futures.get("DATA_").complete(p);
                    } catch (Exception e) {
                        logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                    }
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));

            try {
                S7Message responseMessage = futures.get("DATA_").get();
                S7ParameterUserData parameter = (S7ParameterUserData) responseMessage.getParameter();
                S7ParameterUserDataItemCPUFunctions msgParameter = (S7ParameterUserDataItemCPUFunctions) parameter.getItems().get(0);

                valuesResponse.put(Integer.toString(msgParameter.getSequenceNumber()),
                    decodeEventSubscriptionResponse(Integer.toString(msgParameter.getSequenceNumber()), subscriptionRequest, futures.get("DATA_").get()));
            } catch (Exception ex) {
                logger.warn(ex.toString());
            }

            try {
                HashMap<String, PlcResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
                valuesResponse.forEach((s, p) -> {
                    if (p != null)
                        values.putAll(((DefaultPlcSubscriptionResponse) p).getValues());
                });
                response.complete(new DefaultPlcSubscriptionResponse(subscriptionRequest, values));
            } catch (Exception ex) {
                logger.warn(ex.getMessage());
            }
        });

        t1.start();

        return response;
    }


    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        if (!isConnected()) {
            CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<PlcUnsubscriptionResponse>();
            future.completeExceptionally(new PlcRuntimeException("Disconnected"));
            return future;
        }
        if (!isFeatureSupported()) {
            CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<PlcUnsubscriptionResponse>();
            future.completeExceptionally(new PlcRuntimeException("Not Supported"));
            return future;
        }
        CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<>();
        DefaultPlcUnsubscriptionRequest request = (DefaultPlcUnsubscriptionRequest) unsubscriptionRequest;

        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>();
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>();

        encodeCycledUnSubscriptionRequest(request, parameterItems, payloadItems);

        int tpduId = getTpduId();

        TPKTPacket tpktPacket = new TPKTPacket(
            new COTPPacketData(null,
                new S7MessageUserData(tpduId,
                    new S7ParameterUserData(parameterItems),
                    new S7PayloadUserData(payloadItems)),
                true,
                (byte) tpduId)
        );

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(tpktPacket)
            .onTimeout(new TransactionErrorCallback<>(future, transaction))
            .onError(new TransactionErrorCallback<>(future, transaction))
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .unwrap(TPKTPacket::getPayload)
            .only(COTPPacketData.class)
            .unwrap(COTPPacket::getPayload)
            .check(p -> p.getTpduReference() == tpduId)
            .handle(p -> {
                try {
                    future.complete(null);
                } catch (Exception e) {
                    logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));

        return future;
    }

    private S7Message encodeEventSubscriptionRequest(S7PlcSubscriptionRequest request, int tpduId) {
        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfTags());
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfTags());

        byte subsEvent = 0x00;
        for (String tagName : request.getTagNames()) {
            if (request.getTag(tagName) instanceof DefaultPlcSubscriptionTag) {
                PlcTag event = ((DefaultPlcSubscriptionTag) request.getTag(tagName)).getTag();
                if (event instanceof S7SubscriptionTag) {
                    if (((S7SubscriptionTag) event).getTagType() == S7SubscriptionType.EVENT_SUBSCRIPTION)
                        subsEvent = (byte) (subsEvent | ((S7SubscriptionTag) event).getEventType().getValue());
                }
            }
        }

        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,   //Method
            (byte) 0x04,    //FunctionType
            (byte) 0x04,    //FunctionGroup
            (short) 0x02,   //SubFunction
            (short) 0x00,   //SequenceNumber
            null,   //DataUnitReferenceNumber
            null,   //LastDataUnit
            null         //errorCode
        );
        parameterItems.add(parameter);

        S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest payload;
        if (subsEvent > 0) {
            payload = new S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest(
                DataTransportErrorCode.OK,
                DataTransportSize.OCTET_STRING,
                0x0a,
                subsEvent,
                "HmiRtm  ",
                null,
                null);
        } else {
            //TODO: Check for ALARM_S (S7300) and ALARM_8 (S7400), maybe we need verify the CPU
            AlarmStateType alarmType;
            if (s7DriverContext.getControllerType() == ControllerType.S7_400) {
                alarmType = AlarmStateType.ALARM_INITIATE;
            } else {
                alarmType = AlarmStateType.ALARM_S_INITIATE;
            }

            short auxSubsEvent = (short) (subsEvent & 0xFF);

            payload = new S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest(
                DataTransportErrorCode.OK,
                DataTransportSize.OCTET_STRING,
                (short) 0x0c,
                auxSubsEvent,
                "HmiRtm  ",
                alarmType,
                (short) 0x00);
        }
        payloadItems.add(payload);

        return new S7MessageUserData(tpduId,
            new S7ParameterUserData(parameterItems),
            new S7PayloadUserData(payloadItems));
    }

    private PlcSubscriptionResponse decodeEventSubscriptionResponse(String strTagName,
                                                                    PlcSubscriptionRequest plcSubscriptionRequest,
                                                                    S7Message responseMessage)
        throws PlcProtocolException {

        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
        short errorClass = 0;
        short errorCode = 0;
        if (responseMessage instanceof S7MessageUserData) {
            S7MessageUserData messageUserData = (S7MessageUserData) responseMessage;
            S7Parameter s7Parameter = messageUserData.getParameter();
            if(s7Parameter instanceof S7ParameterUserData) {
                S7ParameterUserData s7ParameterUserData = (S7ParameterUserData) s7Parameter;
                S7ParameterUserDataItem s7ParameterUserDataItem = s7ParameterUserData.getItems().get(0);
                if(s7ParameterUserDataItem instanceof S7ParameterUserDataItemCPUFunctions) {
                    S7ParameterUserDataItemCPUFunctions s7ParameterUserDataItemCPUFunctions = (S7ParameterUserDataItemCPUFunctions) s7ParameterUserDataItem;
                    errorCode = s7ParameterUserDataItemCPUFunctions.getErrorCode().shortValue();
                }
            }
            //S7PayloadUserData payload = (S7PayloadUserData) messageUserData.getPayload();
            //errorClass = payload.getItems()[0]
            //errorCode = messageUserData.getParameter().
        } else if (responseMessage instanceof S7MessageResponse) {
            S7MessageResponse messageResponse = (S7MessageResponse) responseMessage;
            errorClass = messageResponse.getErrorClass();
            errorCode = messageResponse.getErrorCode();
        } else {
            throw new PlcProtocolException("Unsupported message type " + responseMessage.getClass().getName());
        }

        // If the result contains any form of non-null error code, handle this instead.
        if ((errorClass != 0) || (errorCode != 0)) {
            // This is usually the case if PUT/GET wasn't enabled on the PLC
            if ((errorClass == 129) && (errorCode == 4)) {
                logger.warn("Got an error response from the PLC. This particular response code usually indicates " +
                    "that PUT/GET is not enabled on the PLC.");
                for (String tagName : plcSubscriptionRequest.getTagNames()) {
                    values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null));
                }
                return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
            }
            // This seems to be the case if we're trying to do a subscription on a device that doesn't support that.
            else if((errorClass == 0) && (errorCode == (short) 0x8104)) {
                logger.warn("Got an error response from the PLC. Error Class: {}, Error Code {}. This particular response code usually indicates " +
                    "that a given service is not implemented on the PLC. Most probably you tried to subscribe to " +
                    "data on a PLC that doesn't support subscriptions (S7-1200 or S7-1500)",
                    errorClass, errorCode);
                for (String tagName : plcSubscriptionRequest.getTagNames()) {
                    values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.UNSUPPORTED, null));
                }
                return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
            } else {
                logger.warn("Got an unknown error response from the PLC. Error Class: {}, Error Code {}. " +
                        "We probably need to implement explicit handling for this, so please file a bug-report " +
                        "on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump " +
                        "containing a capture of the communication.",
                    errorClass, errorCode);
                for (String tagName : plcSubscriptionRequest.getTagNames()) {
                    values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
                return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
            }
        }

        S7ParameterUserData parameter = (S7ParameterUserData) responseMessage.getParameter();
        List<S7ParameterUserDataItem> parameters = parameter.getItems();
        S7ParameterUserDataItemCPUFunctions itemparameter = (S7ParameterUserDataItemCPUFunctions) parameters.get(0);
        errorCode = itemparameter.getErrorCode().shortValue();

        // In all other cases all went well.
        S7PayloadUserData payload = (S7PayloadUserData) responseMessage.getPayload();

        List<S7PayloadUserDataItem> payloadItems = payload.getItems();

        //Only one item for any number of subscription (4)
        if (payloadItems.isEmpty()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        boolean responseOk = false;
        if (payloadItems.get(0) instanceof S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse) {
            S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse item =
                (S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse) payloadItems.get(0);
            if ((item.getReturnCode() == DataTransportErrorCode.OK) &&
                (item.getTransportSize() == DataTransportSize.OCTET_STRING)) {
                responseOk = true;
            }
        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse) {
            S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse item =
                (S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse) payloadItems.get(0);
            if ((item.getReturnCode() == DataTransportErrorCode.OK) &&
                (item.getTransportSize() == DataTransportSize.OCTET_STRING)) {
                responseOk = true;
            }
        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse) {
            S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse item =
                (S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse) payloadItems.get(0);
            if ((item.getReturnCode() == DataTransportErrorCode.OK) &&
                (item.getTransportSize() == DataTransportSize.OCTET_STRING)) {
                responseOk = true;
            }
        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCpuFunctionAlarmAckResponse) {
            S7PayloadUserDataItemCpuFunctionAlarmAckResponse items =
                (S7PayloadUserDataItemCpuFunctionAlarmAckResponse)
                    payloadItems.get(0);
            //String tagName = (String) plcSubscriptionRequest.getTagNames().toArray()[0];
            //TODO: Chequear si tagName es el correcto           
            //logger.info("strTagName: " + strTagName);
            values.put(strTagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, null));
            for (short s : items.getMessageObjects()) {
                if (s == 0x0000) {
                    values.put(Integer.toHexString(s), new DefaultPlcResponseItem<>(PlcResponseCode.OK, null));
                } else if (s == 0x000a) {
                    values.put(Integer.toHexString(s), new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null));
                }
            }

            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse) {
            // TODO: Items are ignored?
            //S7PayloadUserDataItemCpuFunctionAlarmAckResponse items =
            //    (S7PayloadUserDataItemCpuFunctionAlarmAckResponse)
            //        payloadItems.get(0);
            //String fieldName = (String) S7PayloadUserDataItemCyclicServicesPush .getFieldNames().toArray()[0];
            //logger.warn("Request field: " + strTagName + ": " + S7ParamErrorCode.valueOf(errorCode) + " " + S7ParamErrorCode.valueOf(errorCode).getEvent());
            values.put(strTagName, new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null));
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);

        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) {

            S7PayloadUserDataItemCpuFunctionAlarmQueryResponse items =
                (S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) payloadItems.get(0);

            ByteBuf buffer = Unpooled.directBuffer(items.getItems().length * 2);
            ByteBuf rxBuffer = Unpooled.directBuffer(items.getItems().length * 2);
            buffer.writeBytes(items.getItems());

            if (itemparameter.getLastDataUnit() == 1) {
                short loop = 0xff;
                CompletableFuture<S7MessageUserData> loopFuture;
                S7MessageUserData msg;
                S7ParameterUserDataItemCPUFunctions loopParameter;
                S7PayloadUserDataItemCpuFunctionAlarmQueryResponse loopPayload = null;

                do {
                    loopFuture = reassembledAlarmEvents(itemparameter.getSequenceNumber());

                    try {
                        msg = loopFuture.get();
                        if (msg != null) {
                            loopParameter = (S7ParameterUserDataItemCPUFunctions) ((S7ParameterUserData) msg.getParameter()).getItems().get(0);
                            loopPayload = (S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) ((S7PayloadUserData) msg.getPayload()).getItems().get(0);
                            buffer.writeBytes(loopPayload.getItems());
                            loop = loopParameter.getLastDataUnit();
                        } else {
                            loop = 0x00;
                        }
                    } catch (Exception ex) {
                        logger.warn(ex.toString());
                    }
                } while (loop > 0x00);

                rxBuffer.writeByte(loopPayload.getReturnCode().getValue());
                rxBuffer.writeByte(loopPayload.getTransportSize().getValue());
                rxBuffer.writeShort(loopPayload.getDataLength());
                rxBuffer.writeBytes(buffer);

            } else {
                rxBuffer.writeByte(payloadItems.get(0).getReturnCode().getValue());
                rxBuffer.writeByte(payloadItems.get(0).getTransportSize().getValue());
                rxBuffer.writeShort(payloadItems.get(0).getDataLength());
                rxBuffer.writeBytes(buffer);
            }

            ReadBuffer readBuffer = new ReadBufferByteBased(ByteBufUtil.getBytes(rxBuffer));

            try {
                short cpuSubFunction;
                if (s7DriverContext.getControllerType() == ControllerType.S7_300) {
                    cpuSubFunction = 0x13;
                } else {
                    cpuSubFunction = 0xf0;
                }

                S7PayloadUserDataItem payloadItem =
                    S7PayloadUserDataItem.staticParse(readBuffer,
                        (byte) 0x04,
                        (byte) 0x00,
                        cpuSubFunction);

                // TODO: The eventQueue is only drained in the S7ProtocolEventLogic.ObjectProcessor and here only messages of type S7Event are processed, so S7PayloadUserDataItem elements will just be ignored.
                //eventQueue.add(payloadItem);
            } catch (Exception ex) {
                logger.info(ex.toString());
            }

            PlcResponseCode resCode = (items.getReturnCode() == DataTransportErrorCode.OK) ? PlcResponseCode.OK : PlcResponseCode.INTERNAL_ERROR;
            values.put(strTagName, new DefaultPlcResponseItem<>(resCode, null));
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);

        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCyclicServicesSubscribeResponse) {
            //S7ParameterUserData parameter = (S7ParameterUserData) responseMessage.getParameter();  
            //logger.info("Aqui debe responder a Cyclic transfer");
            S7ParameterUserDataItemCPUFunctions msgParameter = (S7ParameterUserDataItemCPUFunctions)
                parameter.getItems().get(0);

            cycRequests.put(msgParameter.getSequenceNumber(), plcSubscriptionRequest);

            S7CyclicEvent cycEvent = new S7CyclicEvent(plcSubscriptionRequest,
                msgParameter.getSequenceNumber(),
                (S7PayloadUserDataItemCyclicServicesSubscribeResponse) payloadItems.get(0));
            
            if (plcSubscriptionRequest.getTags().get(0).getPlcSubscriptionType() == PlcSubscriptionType.CHANGE_OF_STATE) {
                cycChangeValueEvents.put(msgParameter.getSequenceNumber(), cycEvent);
            }

            eventQueue.add(cycEvent);

            S7PlcSubscriptionHandle cycHandle = new S7PlcSubscriptionHandle(strTagName, EventType.CYC, eventLogic);

            PlcResponseItem<PlcSubscriptionHandle> response = new DefaultPlcResponseItem<>(PlcResponseCode.OK, cycHandle);
            plcSubscriptionRequest.getTagNames().forEach(s -> values.put(s, response));

            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);

        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse) {
            //logger.info("Aqui debe responder a Cyclic transfer Change Driven");
            S7ParameterUserDataItemCPUFunctions msgParameter = (S7ParameterUserDataItemCPUFunctions)
                parameter.getItems().get(0);

            cycRequests.put(msgParameter.getSequenceNumber(), plcSubscriptionRequest);

            S7CyclicEvent cycEvent = new S7CyclicEvent(plcSubscriptionRequest,
                msgParameter.getSequenceNumber(),
                (S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse) payloadItems.get(0));

            eventQueue.add(cycEvent);

            S7PlcSubscriptionHandle cycHandle = new S7PlcSubscriptionHandle(strTagName, EventType.CYC, eventLogic);
            values.put(strTagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, cycHandle));
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);

        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCyclicServicesErrorResponse) {
            //S7ParameterUserData parameter = (S7ParameterUserData) responseMessage.getParameter();
            //S7ParameterUserDataItem[] parameters = parameter.getItems();
            //S7ParameterUserDataItemCPUFunctions itemparameter = (S7ParameterUserDataItemCPUFunctions) parameters[0];
            //errorCode = itemparameter.getErrorCode().shortValue();
            logger.warn("Request field: " + strTagName + ": " + S7ParamErrorCode.valueOf(errorCode) + " " + S7ParamErrorCode.valueOf(errorCode).getEvent());
            // TODO: This is always false because of commented out code above.
            /*if (errorCode == 0x8104) {
                values.put(strTagName, new ResponseItem(PlcResponseCode.UNSUPPORTED, null));
            } else {*/
            values.put(strTagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
            // }
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
        } else if (payloadItems.get(0) instanceof S7PayloadUserDataItemCyclicServicesUnsubscribeResponse) {
            values.put(strTagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, null));
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
        }

        if (responseOk) {
            for (String tagName : plcSubscriptionRequest.getTagNames()) {
                DefaultPlcSubscriptionTag dTag = (DefaultPlcSubscriptionTag) plcSubscriptionRequest.getTag(tagName);
                S7SubscriptionTag tag = (S7SubscriptionTag) dTag.getTag();
                switch (tag.getEventType()) {
                    case MODE:
                        values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, modeHandle));
                        break;
                    case SYS:
                        values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, sysHandle));
                        break;
                    case USR:
                        values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, usrHandle));
                        break;
                    case ALM:
                        values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, almHandle));
                        break;
                }

            }
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
        }

        return null;
    }

    private CompletableFuture<S7Message> performAlarmAckRequest(DefaultPlcReadRequest request) {
        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfTags());
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfTags());

        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,   //Method
            (byte) 0x04,    //FunctionType
            (byte) 0x04,    //FunctionGroup
            (short) 0x0b,   //SubFunction
            (short) 0x00,   //SequenceNumber
            null,   //DataUnitReferenceNumber
            null,   //LastDataUnit
            null         //errorCode
        );
        parameterItems.add(parameter);

        ArrayList<AlarmMessageObjectAckType> messageObjects = null;
        BitSet bs;
        for (String fieldName : request.getTagNames()) {
            if (request.getTag(fieldName) instanceof S7AckTag) {
                PlcTag field = request.getTag(fieldName);
                if (field instanceof S7AckTag) {
                    ArrayList<Integer> arrAlarmIds = ((S7AckTag) field).getAlarmIds();
                    ArrayList<Integer> arrAlarmSigs = ((S7AckTag) field).getAlarmSigs();
                    messageObjects = new ArrayList<>();
                    for (int i = 0; i < arrAlarmIds.size(); i++) {
                        bs = BitSet.valueOf(new byte[]{arrAlarmSigs.get(i).byteValue()});

                        AlarmMessageObjectAckType messageObject = new AlarmMessageObjectAckType(
                            SyntaxIdType.ALARM_ACKSET,
                            (short) 0,
                            arrAlarmIds.get(i),
                            new State(bs.get(7), bs.get(6), bs.get(5), bs.get(4), bs.get(3), bs.get(2), bs.get(1), bs.get(0)),
                            new State(bs.get(7), bs.get(6), bs.get(5), bs.get(4), bs.get(3), bs.get(2), bs.get(1), bs.get(0))
                        );
                        messageObjects.add(messageObject);
                    }
                }
            }
        }

        S7PayloadUserDataItemCpuFunctionAlarmAckRequest payload =
            new S7PayloadUserDataItemCpuFunctionAlarmAckRequest(
                DataTransportErrorCode.OK,
                DataTransportSize.OCTET_STRING,
                0x0c,
                messageObjects);
        payloadItems.add(payload);

        return sendInternal(new S7MessageUserData(getTpduId(),
            new S7ParameterUserData(parameterItems),
            new S7PayloadUserData(payloadItems)));
    }

    private S7Message encodeAlarmQueryRequest(S7PlcSubscriptionRequest request, int tpduId) {
        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfTags());
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfTags());

        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,   //Method
            (byte) 0x04,    //FunctionType
            (byte) 0x04,    //FunctionGroup
            (short) 0x13,   //SubFunction
            (short) 0x00,   //SequenceNumber
            null,   //DataUnitReferenceNumber
            null,   //LastDataUnit
            null         //errorCode
        );
        parameterItems.add(parameter);

        //TODO: Check the type of alarm
        S7PayloadUserDataItemCpuFunctionAlarmQueryRequest payload =
            new S7PayloadUserDataItemCpuFunctionAlarmQueryRequest(
                DataTransportErrorCode.OK,
                DataTransportSize.OCTET_STRING,
                0x0c,
                SyntaxIdType.ALARM_QUERYREQSET,
                QueryType.ALARM_8P,
                AlarmType.ALARM_8);
        payloadItems.add(payload);

        return new S7MessageUserData(tpduId,
            new S7ParameterUserData(parameterItems),
            new S7PayloadUserData(payloadItems));
    }

    private S7Message encodeCycledS7ANYSubscriptionRequest(S7PlcSubscriptionRequest request, int tpduId) {
        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfTags());
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfTags());

        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,   //Method
            (byte) 0x04,    //FunctionType
            (byte) 0x02,    //FunctionGroup
            (short) 0x01,   //SubFunction
            (short) 0x00,   //SequenceNumber
            null,   //DataUnitReferenceNumber
            null,   //LastDataUnit
            null         //errorCode
        );
        parameterItems.add(parameter);

        //TODO: Chequear la asignacionde tipo

        //PlcTag tag = ((DefaultPlcSubscriptionTag) plctag).getTag();
        //S7SubscriptionTag s7tag = (S7SubscriptionTag) tag;
//
        List<CycServiceItemType> items = new ArrayList<>();

        request.getTags().forEach(tag -> {

//            PlcTag plctag = ((DefaultPlcSubscriptionTag) tag).getTag();
            S7SubscriptionTag s7tag = (S7SubscriptionTag) ((DefaultPlcSubscriptionTag) tag).getTag();
//            
            for (S7Tag userField : s7tag.getS7Tags()) {
                items.add(new CycServiceItemAnyType(
                    (short) 0x0a,
                    (short) 0x10, //S7ANY
                    userField.getDataType(),
                    userField.getNumberOfElements(),
                    userField.getBlockNumber(),
                    userField.getMemoryArea(),
                    (((userField.getByteOffset() << 3) | (userField.getBitOffset() & 0x0007)))
                ));
            }
        });

        S7SubscriptionTag s7tag_base = (S7SubscriptionTag) ((DefaultPlcSubscriptionTag) (request.getTags().get(0))).getTag();

//        
//        int i=0;
//        for (S7Tag userfield:s7tag.getS7Tags()) {
//            items.add(new CycServiceItemAnyType(
//                    (short) 0x0a,
//                    (short) 0x10, //S7ANY
//                    userfield.getDataType(),
//                    userfield.getNumberOfElements(),                    
//                    userfield.getBlockNumber(),
//                    userfield.getMemoryArea(),
//                    (long) (((userfield.getByteOffset() << 3) | (userfield.getBitOffset() & 0x0007)))                    
//            ));
//            i++;
//        }
//        
//        //Length in bytes
//        //4+12*items.legth
        int lengthInBytes = 4 + items.size() * 12;
//
        S7PayloadUserDataItemCyclicServicesSubscribeRequest payload =
            new S7PayloadUserDataItemCyclicServicesSubscribeRequest(
                DataTransportErrorCode.OK,
                DataTransportSize.OCTET_STRING,
                lengthInBytes,
                items.size(),
                s7tag_base.getTimeBase(),
                s7tag_base.getMultiplier(),
                items
            );
        payloadItems.add(payload);

        return new S7MessageUserData(tpduId,
            new S7ParameterUserData(parameterItems),
            new S7PayloadUserData(payloadItems));
    }


    private S7Message encodeCycledDBREADSubscriptionRequest(S7PlcSubscriptionRequest request, int tpduId) {
        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfTags());
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfTags());

        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,   //Method
            (byte) 0x04,    //FunctionType
            (byte) 0x02,    //FunctionGroup
            (short) 0x01,   //SubFunction
            (short) 0x00,   //SequenceNumber
            null,   //DataUnitReferenceNumber
            null,   //LastDataUnit
            null         //errorCode
        );
        parameterItems.add(parameter);

        //TODO: Chequear la asignacionde tipo

        PlcSubscriptionTag plctag = request.getTags().get(0);

        PlcTag tag = ((DefaultPlcSubscriptionTag) plctag).getTag();
        S7SubscriptionTag s7tag = (S7SubscriptionTag) tag;

        //CycServiceItemDbReadType[] cycitems = new CycServiceItemDbReadType[(int) cycitemcount];
        List<CycServiceItemType> cycItems = new ArrayList<>();
        ArrayList<SubItem> subItems = new ArrayList<>();
        for (S7Tag userTag : s7tag.getS7Tags()) {
            subItems.add(new SubItem((short) userTag.getNumberOfElements(),
                userTag.getBlockNumber(),
                userTag.getByteOffset()));
        }

        int initPos = 0;
        int endPos = (subItems.size() < 50) ? subItems.size() : 49;
        int j = 0;
        int lengthInBytes = 4;
        do {
            List<SubItem> arraySubItems = subItems.subList(initPos, endPos);

            cycItems.add(j, new CycServiceItemDbReadType(
                (short) (arraySubItems.size() * 5 + 2),
                (short) 0xb0,
                (short) arraySubItems.size(),
                arraySubItems));

            lengthInBytes += 4 + arraySubItems.size() * 5;
            initPos = endPos + 1;
            endPos = Math.min((initPos + 49), subItems.size());
            j++;
            // TODO: j is always equal to cycItems.size ... to technically this loop doesn't make any sense and can be replaced by a single execution.
        } while (j < cycItems.size());

        S7PayloadUserDataItemCyclicServicesSubscribeRequest payload =
            new S7PayloadUserDataItemCyclicServicesSubscribeRequest(
                DataTransportErrorCode.OK,
                DataTransportSize.OCTET_STRING,
                lengthInBytes,
                cycItems.size(),
                s7tag.getTimeBase(),
                s7tag.getMultiplier(),
                cycItems
            );
        payloadItems.add(payload);

        return new S7MessageUserData(tpduId,
            new S7ParameterUserData(parameterItems),
            new S7PayloadUserData(payloadItems));
    }


    private void encodeCycledUnSubscriptionRequest(DefaultPlcUnsubscriptionRequest request,
                                                   List<S7ParameterUserDataItem> parameterItems,
                                                   List<S7PayloadUserDataItem> payloadItems) {

        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,   //Method
            (byte) 0x04,    //FunctionType
            (byte) 0x02,    //FunctionGroup
            (short) 0x04,   //SubFunction
            (short) 0x00,   //SequenceNumber
            null,   //DataUnitReferenceNumber
            null,   //LastDataUnit
            null         //errorCode
        );

        parameterItems.clear();
        parameterItems.add(parameter);


        List<PlcSubscriptionHandle> handles = request.getSubscriptionHandles();

        //PlcField field = ((DefaultPlcSubscriptionField) plcfield).getPlcField(); 
        //S7SubscriptionField s7field = (S7SubscriptionField) field;

        payloadItems.clear();
        //TODO:Check CPU type
        handles.forEach(h -> {
            S7PayloadUserDataItemCyclicServicesUnsubscribeRequest payload =
                new S7PayloadUserDataItemCyclicServicesUnsubscribeRequest(
                    DataTransportErrorCode.OK,
                    DataTransportSize.OCTET_STRING,
                    0x02,
                    (short) 0x01,
                    Short.parseShort(((S7PlcSubscriptionHandle) h).getEventId()));

            payloadItems.add(payload);
        });

//        s7tag.getAckAlarms().forEach(jobid -> {
//                S7PayloadUserDataItemCyclicServicesUnsubscribeRequest payload =
//                        new S7PayloadUserDataItemCyclicServicesUnsubscribeRequest (
//                                DataTransportErrorCode.OK,
//                                DataTransportSize.OCTET_STRING,
//                                0x02,
//                                (short) 0x05,
//                                jobid.byteValue());
//
//                payloadItems.add(payload);
//            });
    }

    /*
     *
     */
    private CompletableFuture<S7Message> performClkRequest(DefaultPlcReadRequest request) {
        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfTags());
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfTags());

        final S7ClkTag tag = (S7ClkTag) request.getTags().get(0);
        int subFunction = tag.getAddressString().equals("CLK") ? 1 : 3;

        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,   //Method
            (byte) 0x04,    //FunctionType
            (byte) 0x07,    //FunctionGroup
            (short) subFunction,   //SubFunction
            (short) 0x00,   //SequenceNumber
            null,   //DataUnitReferenceNumber
            null,   //LastDataUnit
            null         //errorCode
        );
        parameterItems.add(parameter);

        S7PayloadUserDataItemClkRequest payload;
        payload = new S7PayloadUserDataItemClkRequest(
            DataTransportErrorCode.NOT_FOUND,
            DataTransportSize.NULL,
            0x00);
        payloadItems.add(payload);

        return sendInternal(new S7MessageUserData(getTpduId(),
            new S7ParameterUserData(parameterItems),
            new S7PayloadUserData(payloadItems)));
    }

    /*
     *
     */
    private CompletableFuture<S7Message> performClkSetRequest(DefaultPlcWriteRequest request) {
        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfTags());
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfTags());

        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,   //Method
            (byte) 0x04,    //FunctionType
            (byte) 0x07,    //FunctionGroup
            (short) 0x04,   //SubFunction
            (short) 0x00,   //SequenceNumber
            null,   //DataUnitReferenceNumber
            null,   //LastDataUnit
            null         //errorCode
        );
        parameterItems.add(parameter);

        S7ClkTag tag = (S7ClkTag) request.getTags().get(0);

        S7PayloadUserDataItemClkSetRequest payload;
        payload = new S7PayloadUserDataItemClkSetRequest(
            DataTransportErrorCode.OK,
            DataTransportSize.OCTET_STRING,
            0x0A,
            tag.getDateAndTime());
        payloadItems.add(payload);

        return sendInternal(
            new S7MessageUserData(getTpduId(),
                new S7ParameterUserData(parameterItems),
                new S7PayloadUserData(payloadItems)));
    }


    /*
     * Encoding of STRING types ... for WSTRING types the "Maximum length" and
     * "Current length" are both 16 bit unsigned integer values:
     *
     *            +-------------------+
     * Byte n     | Maximum length    | (intMaxChars)
     *            +-------------------+
     * Byte n+1   | Current Length    | (intActualChars)
     *            +-------------------+
     * Byte n+2   | 1st character     | \         \
     *            +-------------------+  |         |
     * Byte n+3   | 2st character     |  | Current |
     *            +-------------------+   >        |
     * Byte ...   | ...               |  | length  |  Maximum
     *            +-------------------+  |          >
     * Byte n+m+1 | mth character     | /          |  length
     *            +-------------------+            |
     * Byte ...   | ...               |            |
     *            +-------------------+            |
     * Byte ...   | ...               |           /
     *            +-------------------+
     *
     * Reading text strings in two steps:
     * 1. For the operation on texts, you must first evaluate the
     *    space created on DB of type STRING or WSTRING.
     * 2. The first two bytes have the maximum number of characters (bytes)
     *    available to store text strings (intMaxChars) and the number of
     *    characters available (intActualChars).
     * 3. In the specific case of reading, only the characters defined
     *    by "intActualChars" are recovered.
     */
    private CompletableFuture<S7Message> performVarLengthStringReadRequest(DefaultPlcReadRequest request) {
        CompletableFuture<S7Message> future = new CompletableFuture<>();

        // Replace the var-length string fields with requests to read the
        // length information instead of the string content.
        int numVarLengthStrings = 0;
        LinkedHashMap<String, PlcTagItem<PlcTag>> updatedRequestItems = new LinkedHashMap<>(request.getNumberOfTags());
        for (String tagName : request.getTagNames()) {
            PlcTagItem<PlcTag> plcTagItem = request.getTagItem(tagName);
            if(plcTagItem.getTag() instanceof S7StringVarLengthTag) {
                S7Tag s7Tag = (S7Tag) plcTagItem.getTag();
                TransportSize dataType = s7Tag.getDataType();
                if(dataType == TransportSize.STRING) {
                    updatedRequestItems.put(tagName, new DefaultPlcTagItem<>(new S7Tag(TransportSize.BYTE, s7Tag.getMemoryArea(), s7Tag.getBlockNumber(), s7Tag.getByteOffset(), s7Tag.getBitOffset(), 2)));
                    numVarLengthStrings++;
                } else if(dataType == TransportSize.WSTRING) {
                    updatedRequestItems.put(tagName, new DefaultPlcTagItem<>(new S7Tag(TransportSize.BYTE, s7Tag.getMemoryArea(), s7Tag.getBlockNumber(), s7Tag.getByteOffset(), s7Tag.getBitOffset(), 4)));
                    numVarLengthStrings++;
                }
            } else {
                updatedRequestItems.put(tagName, plcTagItem);
            }
        }

        CompletableFuture<S7Message> s7MessageCompletableFuture = performOrdinaryReadRequest(new DefaultPlcReadRequest(request.getReader(), updatedRequestItems));
        int finalNumVarLengthStrings = numVarLengthStrings;
        s7MessageCompletableFuture.whenComplete((s7Message, throwable1) -> {
            if (throwable1 != null) {
                future.completeExceptionally(throwable1);
                return;
            }
            // Collect the responses for the var-length strings and read them separately.
            LinkedHashMap<String, PlcTagItem<PlcTag>> varLengthStringTags = new LinkedHashMap<>(finalNumVarLengthStrings);
            int curItem = 0;
            for (String tagName : request.getTagNames()) {
                S7Tag s7tag = (S7Tag) request.getTag(tagName);
                if(s7tag instanceof S7StringVarLengthTag) {
                    S7VarPayloadDataItem s7VarPayloadDataItem = ((S7PayloadReadVarResponse) s7Message.getPayload()).getItems().get(curItem);
                    // Simply ignore processing var-length strings that are not ok
                    if(s7VarPayloadDataItem.getReturnCode() == DataTransportErrorCode.OK) {
                        ReadBuffer rb = new ReadBufferByteBased(s7VarPayloadDataItem.getData());
                        try {
                            if (s7tag.getDataType() == TransportSize.STRING) {
                                rb.readShort(8);
                                int stringLength = rb.readShort(8);
                                varLengthStringTags.put(tagName, new DefaultPlcTagItem<>(new S7StringFixedLengthTag(TransportSize.STRING, s7tag.getMemoryArea(), s7tag.getBlockNumber(), s7tag.getByteOffset(), s7tag.getBitOffset(), 1, stringLength)));
                            } else if (s7tag.getDataType() == TransportSize.WSTRING) {
                                rb.readInt(16);
                                int stringLength = rb.readInt(16);
                                varLengthStringTags.put(tagName, new DefaultPlcTagItem<>(new S7StringFixedLengthTag(TransportSize.WSTRING, s7tag.getMemoryArea(), s7tag.getBlockNumber(), s7tag.getByteOffset(), s7tag.getBitOffset(), 1, stringLength)));
                            }
                        } catch (Exception e) {
                            logger.warn("Error parsing string size for tag {}", tagName, e);
                        }
                    }
                }
                curItem++;
            }
            // TODO: Technically we would need to let this go through the optimizer in order to split things up.
            //  For this we need access to the PlcReader instance of this driver
            CompletableFuture<S7Message> readStringsCompletableFuture = performOrdinaryReadRequest(new DefaultPlcReadRequest(request.getReader(), varLengthStringTags));
            readStringsCompletableFuture.whenComplete((s7StringMessage, throwable2) -> {
                // Build a new S7Message that replaces the var-length string items of the previous request with the responses of this response.
                int curInitialItem = 0;
                int curVarLengthStringItem = 0;
                List<S7VarPayloadDataItem> varLengthStringItems = new ArrayList<>(request.getNumberOfTags());
                for (String tagName : request.getTagNames()) {
                    S7Tag s7tag = (S7Tag) request.getTag(tagName);
                    S7VarPayloadDataItem curResultItem = ((S7PayloadReadVarResponse) s7Message.getPayload()).getItems().get(curInitialItem);
                    if(s7tag instanceof S7StringVarLengthTag) {
                        if(curResultItem.getReturnCode() == DataTransportErrorCode.OK) {
                            curResultItem = ((S7PayloadReadVarResponse) s7StringMessage.getPayload()).getItems().get(curVarLengthStringItem);
                            curVarLengthStringItem++;
                        }
                    }
                    varLengthStringItems.add(curResultItem);
                    curInitialItem++;
                }
                future.complete(new S7MessageResponse(s7Message.getTpduReference(), s7Message.getParameter(), new S7PayloadReadVarResponse(varLengthStringItems), (short) 0, (short) 0));
            });
        });
        return future;
    }

    private CompletableFuture<S7Message> performOrdinaryReadRequest(DefaultPlcReadRequest request) {
        // Convert each tag in the request into a corresponding item used in the S7 protocol.
        List<S7VarRequestParameterItem> requestItems = new ArrayList<>(request.getNumberOfTags());
        for (PlcTag tag : request.getTags()) {
            requestItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(tag)));
        }

        // Create a read request template.
        // tpuId will be inserted before sending in #readInternal, so we insert -1 as dummy here
        S7Message requestMessage = new S7MessageRequest(getTpduId(),
            new S7ParameterReadVarRequest(requestItems),
            null);

        return sendInternal(requestMessage);
    }


    /*
     * Encoding of STRING types ... for WSTRING types the "Maximum length" and
     * "Current length" are both 16 bit unsigned integer values:
     *
     *            +-------------------+
     * Byte n     | Maximum length    | (intMaxChars)
     *            +-------------------+
     * Byte n+1   | Current Length    | (intActualChars)
     *            +-------------------+
     * Byte n+2   | 1st character     | \         \
     *            +-------------------+  |         |
     * Byte n+3   | 2st character     |  | Current |
     *            +-------------------+   >        |
     * Byte ...   | ...               |  | length  |  Maximum
     *            +-------------------+  |          >
     * Byte n+m+1 | mth character     | /          |  length
     *            +-------------------+            |
     * Byte ...   | ...               |            |
     *            +-------------------+            |
     * Byte ...   | ...               |           /
     *            +-------------------+
     *
     * Reading text strings in two steps:
     * 1. For the operation on texts, you must first evaluate the
     *    space created on DB of type STRING or WSTRING.
     * 2. The first two bytes have the maximum number of characters (bytes)
     *    available to store text strings (intMaxChars) and the number of
     *    characters available (intActualChars).
     * 3. In the specific case of write string, only the max characters defined
     *    by "intMaxChars" are written.
     * TODO: Maximum waiting time managed by system variables.
     */
    private CompletableFuture<S7Message> performVarLengthStringWriteRequest(DefaultPlcWriteRequest request) {
        CompletableFuture<S7Message> future = new CompletableFuture<>();

        // Resolve the lengths of all var-length string fields in the request.
        CompletableFuture<Map<S7StringVarLengthTag, StringSizes>> stringSizesFuture = getStringSizes(request);
        stringSizesFuture.whenComplete((s7StringVarLengthTagStringSizesMap, throwable) -> {
            if (throwable != null) {
                future.completeExceptionally(new PlcProtocolException("Error resolving string sizes", throwable));
            } else {
                // Create an alternative list of request items, where all var-length string tags are replaced with
                // fixed-length string tags using the string length returned by the previous request.
                LinkedHashMap<String, PlcTagValueItem<PlcTag>> updatedRequestItems = new LinkedHashMap<>(request.getNumberOfTags());
                for (String tagName : request.getTagNames()) {
                    PlcTag tag = request.getTag(tagName);
                    PlcValue value = request.getPlcValue(tagName);
                    if (tag instanceof S7StringVarLengthTag) {
                        S7StringVarLengthTag varLengthTag = (S7StringVarLengthTag) tag;
                        int stringLength = s7StringVarLengthTagStringSizesMap.get(varLengthTag).getCurLength();
                        S7StringFixedLengthTag newTag = new S7StringFixedLengthTag(varLengthTag.getDataType(), varLengthTag.getMemoryArea(),
                            varLengthTag.getBlockNumber(), varLengthTag.getByteOffset(), varLengthTag.getBitOffset(),
                            varLengthTag.getNumberOfElements(), stringLength);
                        updatedRequestItems.put(tagName, new DefaultPlcTagValueItem<>(newTag, value));
                    } else {
                        updatedRequestItems.put(tagName, new DefaultPlcTagValueItem<>(tag, value));
                    }
                }

                // Use the normal functionality to execute the read request.
                // TODO: Here technically the request object in the response will not match the original request.
                CompletableFuture<S7Message> s7MessageCompletableFuture = performOrdinaryWriteRequest(
                    new DefaultPlcWriteRequest(request.getWriter(), updatedRequestItems));
                s7MessageCompletableFuture.whenComplete((s7Message, throwable1) -> {
                    if (throwable1 != null) {
                        future.completeExceptionally(throwable1);
                    } else {
                        future.complete(s7Message);
                    }
                });
            }
        });

        return future;
    }

    private CompletableFuture<S7Message> performOrdinaryWriteRequest(DefaultPlcWriteRequest request) {
        List<S7VarRequestParameterItem> parameterItems = new ArrayList<>(request.getNumberOfTags());
        List<S7VarPayloadDataItem> payloadItems = new ArrayList<>(request.getNumberOfTags());

        for (String tagName : request.getTagNames()) {
            final S7Tag tag = (S7Tag) request.getTag(tagName);
            final PlcValue plcValue = request.getPlcValue(tagName);
            parameterItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(tag)));
            payloadItems.add(serializePlcValue(tag, plcValue));
        }

        return sendInternal(
            new S7MessageRequest(getTpduId(),
                new S7ParameterWriteVarRequest(parameterItems),
                new S7PayloadWriteVarRequest(payloadItems)
            ));
    }

    /**
     * Sends one Read over the Wire and internally returns the Response
     * Do sending of normally sized single-message request.
     * <p>
     * Assumes that the {@link S7MessageRequest} and its expected {@link S7MessageResponseData}
     * and does not further check that!
     */
    private CompletableFuture<S7Message> sendInternal(S7Message request) {
        CompletableFuture<S7Message> future = new CompletableFuture<>();

        // Get the tpduId from the S7 message.
        int tpduId = request.getTpduReference();
        //The "COTP - TPDU nnumber" field must always be zero.
        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null, request, true, (byte) 0));

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        // Send the request.
        transaction.submit(() -> conversationContext.sendRequest(tpktPacket)
            .onTimeout(new TransactionErrorCallback<>(future, transaction))
            .onError(new TransactionErrorCallback<>(future, transaction))
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .unwrap(TPKTPacket::getPayload)
            .only(COTPPacketData.class)
            .check(p -> p.getPayload() != null)
            .unwrap(COTPPacket::getPayload)
            .check(p -> p.getTpduReference() == tpduId)
            .handle(p -> {
                future.complete(p);
                // Finish the request-transaction.
                transaction.endRequest();
            }));

        return future;
    }

    /**
     * DECODE:
     * This method is called when there is no handler for the message. 
     * By default it must correspond to asynchronous events, which if so, 
     * must be transferred to the event queue.
     * The event's own information is encapsulated in the parameters and payload
     * field. From this it is abstracted to the corresponding event model.
     * 01. S7ModeEvent:
     * 02. S7UserEvent:
     * 03. S7SysEvent:
     * 04. S7AlarmEvent
     * 05. S7CyclicEvent:
     * 06. S7CyclicEvent:
     * TODO: Use mspec to generate types that allow better interpretation of
     * the code using "instanceof".
     */
    @Override
    protected void decode(ConversationContext<TPKTPacket> context, TPKTPacket msg) throws Exception {
        
        final S7Message s7msg = msg.getPayload().getPayload();
        final S7Parameter parameter = s7msg.getParameter();
        final S7PayloadUserData payload = (S7PayloadUserData) s7msg.getPayload();       
        
        if (parameter instanceof S7ParameterModeTransition) {  //(01)  
            
            S7ModeEvent modeEvent = new S7ModeEvent((S7ParameterModeTransition) parameter);
            eventQueue.add(modeEvent);
            
        } else if (parameter instanceof S7ParameterUserData) {
            
            S7ParameterUserData parameterUD = (S7ParameterUserData) parameter;
            List<S7ParameterUserDataItem> parameterUDItems = parameterUD.getItems();
            
            for (S7ParameterUserDataItem parameterUDItem : parameterUDItems) {
                
                if (parameterUDItem instanceof S7ParameterUserDataItemCPUFunctions) {
                    
                    S7ParameterUserDataItemCPUFunctions myParameter = (S7ParameterUserDataItemCPUFunctions) parameterUDItem;
                    
                    if ((myParameter.getCpuFunctionType() == 0x00) && (myParameter.getCpuSubfunction() == 0x03)) { //(02)                      
                        
                        payload.getItems().forEach(item ->{
                            if (item instanceof S7PayloadDiagnosticMessage) {
                                final S7PayloadDiagnosticMessage pload = (S7PayloadDiagnosticMessage) item; 
                                if ((pload.getEventId() >= 0x0A000) & (pload.getEventId() <= 0x0BFFF)) {
                                    S7UserEvent userEvent = S7UserEvent.of(pload);
                                    eventQueue.add(userEvent);                                
                                } else {
                                    S7SysEvent sysEvent = S7SysEvent.of(pload);
                                    eventQueue.add(sysEvent);                                        
                                }
                            } 
                        });

                    } else if ((myParameter.getCpuFunctionType() == 0x00) &&
                        ((myParameter.getCpuSubfunction() == 0x05) ||
                            (myParameter.getCpuSubfunction() == 0x06) ||
                            (myParameter.getCpuSubfunction() == 0x0c) ||
                            (myParameter.getCpuSubfunction() == 0x11) ||
                            (myParameter.getCpuSubfunction() == 0x12) ||
                            (myParameter.getCpuSubfunction() == 0x13) ||
                            (myParameter.getCpuSubfunction() == 0x16))) { //(04)
                        
                        payload.getItems().forEach(item ->{
                            S7AlarmEvent alrmEvent = S7AlarmEvent.of(item);
                            eventQueue.add(alrmEvent);                                
                        });
                        
                                                
                    } else if ((myParameter.getCpuFunctionType() == 0x00) && (myParameter.getCpuSubfunction() == 0x13)) {
                        //TODO: Requires reverse engineering.
                    } else if (((myParameter.getCpuFunctionGroup() == 0x02) && (myParameter.getCpuFunctionType() == 0x00) && (myParameter.getCpuSubfunction() == 0x01))) { //(05)

                        S7ParameterUserDataItemCPUFunctions parameterItem =
                            (S7ParameterUserDataItemCPUFunctions)
                                ((S7ParameterUserData) parameter).getItems().get(0);

                        S7PayloadUserDataItemCyclicServicesPush payloadItem =
                            (S7PayloadUserDataItemCyclicServicesPush)
                                payload.getItems().get(0);

                        S7CyclicEvent cycEvent = new S7CyclicEvent(cycRequests.get(parameterItem.getSequenceNumber()),
                            parameterItem.getSequenceNumber(),
                            payloadItem);
                        
                        if (cycChangeValueEvents.containsKey(parameterItem.getSequenceNumber())){
                            S7CyclicEvent lastCycEvent = cycChangeValueEvents.get(parameterItem.getSequenceNumber());
                            if (!cycEvent.equals(lastCycEvent)) {
                                cycChangeValueEvents.replace(parameterItem.getSequenceNumber(), cycEvent);
                                eventQueue.add(cycEvent);                                
                            }
                            
                        } else {
                            eventQueue.add(cycEvent);
                        }

                    } else if ((myParameter.getCpuFunctionGroup() == 0x02) && (myParameter.getCpuFunctionType() == 0x00) && (myParameter.getCpuSubfunction() == 0x05)) { //(06)
                        
                        S7ParameterUserDataItemCPUFunctions parameterItem =
                            (S7ParameterUserDataItemCPUFunctions)
                                ((S7ParameterUserData) parameter).getItems().get(0);

                        S7PayloadUserDataItemCyclicServicesChangeDrivenPush payloadItem =
                            (S7PayloadUserDataItemCyclicServicesChangeDrivenPush)
                                payload.getItems().get(0);

                        S7CyclicEvent cycEvent = new S7CyclicEvent(null,
                            parameterItem.getSequenceNumber(),
                            payloadItem);
                        eventQueue.add(cycEvent);

                    } else if ((myParameter.getCpuFunctionType() == 0x08) && (myParameter.getCpuSubfunction() == 0x01)) {
                        //TODO: Requires reverse engineering.
                    } else if ((myParameter.getCpuFunctionType() == 0x08) && (myParameter.getCpuSubfunction() == 0x04)) {
                        //TODO: Requires reverse engineering.
                    }
                }
            }
        }
    }

    private void extractControllerTypeAndFireConnected(ConversationContext<TPKTPacket> context, S7PayloadUserData payloadUserData) {
        for (S7PayloadUserDataItem item : payloadUserData.getItems()) {
            if (!(item instanceof S7PayloadUserDataItemCpuFunctionReadSzlResponse)) {
                continue;
            }
            S7PayloadUserDataItemCpuFunctionReadSzlResponse readSzlResponseItem =
                (S7PayloadUserDataItemCpuFunctionReadSzlResponse) item;

//            for (SzlDataTreeItem readSzlResponseItemItem : readSzlResponseItem.getItems()) {
//                if (readSzlResponseItemItem.getItemIndex() != 0x0001) {
//                    continue;
//                }
//                final String articleNumber = new String(readSzlResponseItemItem.getMlfb());
//                s7DriverContext.setControllerType(decodeControllerType(articleNumber));
//
//                // Send an event that connection setup is complete.
//                context.fireConnected();
//            }
            ByteBuf szlItem = Unpooled.wrappedBuffer(readSzlResponseItem.getItems());
            String articleNumber = szlItem.toString(2, 20, Charset.defaultCharset());
            s7DriverContext.setControllerType(decodeControllerType(articleNumber));
            context.fireConnected();
        }
    }

    private TPKTPacket createIdentifyRemoteMessage() {
        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(1, new S7ParameterUserData(Collections.singletonList(
            new S7ParameterUserDataItemCPUFunctions((short) 0x11, (byte) 0x4, (byte) 0x4, (short) 0x01, (short) 0x00, null, null, null)
        )), new S7PayloadUserData(Collections.singletonList(
            new S7PayloadUserDataItemCpuFunctionReadSzlRequest(DataTransportErrorCode.OK,
                DataTransportSize.OCTET_STRING,
                0x0C,
                new SzlId(SzlModuleTypeClass.CPU,
                    (byte) 0x00,
                    SzlSublist.MODULE_IDENTIFICATION),
                0x0000)
        )));
        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (byte) 2);
        return new TPKTPacket(cotpPacketData);
    }

    private TPKTPacket createS7ConnectionRequest(COTPPacketConnectionResponse cotpPacketConnectionResponse) {
        for (COTPParameter parameter : cotpPacketConnectionResponse.getParameters()) {
            if (parameter instanceof COTPParameterCalledTsap) {
                COTPParameterCalledTsap cotpParameterCalledTsap = (COTPParameterCalledTsap) parameter;
                s7DriverContext.setCalledTsapId(cotpParameterCalledTsap.getTsapId());
            } else if (parameter instanceof COTPParameterCallingTsap) {
                COTPParameterCallingTsap cotpParameterCallingTsap = (COTPParameterCallingTsap) parameter;
                if (cotpParameterCallingTsap.getTsapId() != s7DriverContext.getCallingTsapId()) {
                    s7DriverContext.setCallingTsapId(cotpParameterCallingTsap.getTsapId());
                    logger.warn("Switching calling TSAP id to '{}'", s7DriverContext.getCallingTsapId());
                }
            } else if (parameter instanceof COTPParameterTpduSize) {
                COTPParameterTpduSize cotpParameterTpduSize = (COTPParameterTpduSize) parameter;
                s7DriverContext.setCotpTpduSize(cotpParameterTpduSize.getTpduSize());
            } else {
                logger.warn("Got unknown parameter type '{}'", parameter.getClass().getName());
            }
        }

        // Send an S7 login message.
        S7ParameterSetupCommunication s7ParameterSetupCommunication =
            new S7ParameterSetupCommunication(
                s7DriverContext.getMaxAmqCaller(), s7DriverContext.getMaxAmqCallee(), s7DriverContext.getPduSize());
        S7Message s7Message = new S7MessageRequest(0, s7ParameterSetupCommunication,
            null);
        COTPPacketData cotpPacketData = new COTPPacketData(null, s7Message, true, (byte) 0);
        return new TPKTPacket(cotpPacketData);
    }

    private COTPPacketConnectionRequest createCOTPConnectionRequest(int calledTsapId, int callingTsapId, COTPTpduSize cotpTpduSize) {
        return new COTPPacketConnectionRequest(
            Arrays.asList(
                new COTPParameterCallingTsap(callingTsapId),
                new COTPParameterCalledTsap(calledTsapId),
                new COTPParameterTpduSize(cotpTpduSize)
            ), null, (short) 0x0000, (short) 0x000F, COTPProtocolClass.CLASS_0);
    }

    private PlcResponse decodeReadResponse(S7Message responseMessage, PlcReadRequest plcReadRequest) throws PlcProtocolException {
        Map<String, PlcResponseItem<PlcValue>> values = new HashMap<>();
        short errorClass;
        short errorCode;

        S7ParameterUserDataItemCPUFunctions parameteritem = null;
        if (responseMessage instanceof S7MessageResponseData) {
            S7MessageResponseData messageResponseData = (S7MessageResponseData) responseMessage;
            errorClass = messageResponseData.getErrorClass();
            errorCode = messageResponseData.getErrorCode();
        } else if (responseMessage instanceof S7MessageResponse) {
            S7MessageResponse messageResponse = (S7MessageResponse) responseMessage;
            errorClass = messageResponse.getErrorClass();
            errorCode = messageResponse.getErrorCode();
        } else if (responseMessage instanceof S7MessageUserData) {
            S7MessageUserData messageResponse = (S7MessageUserData) responseMessage;
            S7ParameterUserData parameters = (S7ParameterUserData) messageResponse.getParameter();
            parameteritem = (S7ParameterUserDataItemCPUFunctions) parameters.getItems().get(0);
            errorClass = 0;
            errorCode = parameteritem.getErrorCode().shortValue();
        } else {
            throw new PlcProtocolException("Unsupported message type " + responseMessage.getClass().getName());
        }
        // If the result contains any form of non-null error code, handle this instead.
        if ((errorClass != 0) || (errorCode != 0)) {
            // This is usually the case if PUT/GET wasn't enabled on the PLC
            if ((errorClass == 129) && (errorCode == 4)) {
                logger.warn("Got an error response from the PLC. This particular response code usually indicates " +
                    "that PUT/GET is not enabled on the PLC.");
                for (String tagName : plcReadRequest.getTagNames()) {
                    PlcResponseItem<PlcValue> result = new DefaultPlcResponseItem<>(PlcResponseCode.ACCESS_DENIED, new PlcNull());
                    values.put(tagName, result);
                }
                return new DefaultPlcReadResponse(plcReadRequest, values);
            } else if ((errorClass == 0x85) && (errorCode == 0)) {
                logger.warn("Got an error response from the PLC. This particular response code usually indicates " +
                    "that we sent a too large packet or would be receiving a too large one. " +
                    "Please report this, as this is most probably a bug.");
                for (String tagName : plcReadRequest.getTagNames()) {
                    PlcResponseItem<PlcValue> result = new DefaultPlcResponseItem<>(PlcResponseCode.ACCESS_DENIED, new PlcNull());
                    values.put(tagName, result);
                }
                return new DefaultPlcReadResponse(plcReadRequest, values);
            } else {
                logger.warn("Got an unknown error response from the PLC. Error Class: {}, Error Code {}. " +
                        "We probably need to implement explicit handling for this, so please file a bug-report " +
                        "on https://github.com/apache/plc4x/issues and ideally attach a WireShark dump " +
                        "containing a capture of the communication.",
                    errorClass, errorCode);
                for (String tagName : plcReadRequest.getTagNames()) {
                    PlcResponseItem<PlcValue> result = new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, new PlcNull());
                    values.put(tagName, result);
                }
                return new DefaultPlcReadResponse(plcReadRequest, values);
            }
        }

        if (responseMessage instanceof S7MessageUserData) {
            S7PayloadUserData payload = (S7PayloadUserData) responseMessage.getPayload();
            if (plcReadRequest.getNumberOfTags() != payload.getItems().size()) {
                throw new PlcProtocolException(
                    "The number of requested items doesn't match the number of returned items");
            }

            List<S7PayloadUserDataItem> payloadItems = payload.getItems();

            PlcResponseCode responseCode = PlcResponseCode.INTERNAL_ERROR;
            PlcValue plcValue = null;
            int index = 0;
            for (String tagName : plcReadRequest.getTagNames()) {

                if (plcReadRequest.getTag(tagName) instanceof S7SzlTag) {

                    S7PayloadUserDataItemCpuFunctionReadSzlResponse payloadItem = (S7PayloadUserDataItemCpuFunctionReadSzlResponse) payloadItems.get(index);
                    responseCode = decodeResponseCode(payloadItem.getReturnCode());

                    if (responseCode == PlcResponseCode.OK) {
                        try {
                            List<PlcValue> plcValues;
                            byte[] data = payloadItem.getItems();

                            plcValues = new LinkedList<>();
                            for (byte b : data) {
                                plcValues.add(new PlcSINT(b));
                            }

                            if (parameteritem.getLastDataUnit() == 1) {
                                CompletableFuture<S7MessageUserData> nextFuture;
                                S7ParameterUserData nextParameter;
                                S7PayloadUserData nextPayload;
                                S7PayloadUserDataItemCpuFunctionReadSzlResponse nextPayloadItem;

                                while (parameteritem.getLastDataUnit() == 1) {
                                    //TODO: Just wait for one answer!. Pending for other packages for rearm.
                                    nextFuture = reassembledMessage(parameteritem.getSequenceNumber(), plcValues);

                                    S7MessageUserData msg;

                                    msg = nextFuture.get();
                                    if (msg != null) {
                                        nextParameter = (S7ParameterUserData) msg.getParameter();
                                        parameteritem = (S7ParameterUserDataItemCPUFunctions) nextParameter.getItems().get(0);
                                        nextPayload = (S7PayloadUserData) msg.getPayload();
                                        nextPayloadItem = (S7PayloadUserDataItemCpuFunctionReadSzlResponse) nextPayload.getItems().get(0);
                                        for (byte b : nextPayloadItem.getItems()) {
                                            plcValues.add(new PlcSINT(b));
                                        }
                                    }

                                    plcValue = new PlcList(plcValues);
                                }
                            } else {
                                plcValue = new PlcList(plcValues);
                            }
                        } catch (Exception e) {
                            throw new PlcProtocolException("Error decoding PlcValue", e);
                        }

                    }

                }
                if (plcReadRequest.getTag(tagName) instanceof S7AckTag) {
                    S7PayloadUserDataItemCpuFunctionAlarmAckResponse payloadItem =
                        (S7PayloadUserDataItemCpuFunctionAlarmAckResponse) payloadItems.get(index);
                    responseCode = decodeResponseCode(payloadItem.getReturnCode());
                    List<Short> data = payloadItem.getMessageObjects();
                    List<PlcValue> plcValues = new LinkedList<>();
                    for (short b : data) {
                        plcValues.add(new PlcSINT((byte) b));
                    }
                    plcValue = new PlcList(plcValues);
                }
                if (plcReadRequest.getTag(tagName) instanceof S7ClkTag) {
                    DateAndTime dt;
                    if (payloadItems.get(index) instanceof S7PayloadUserDataItemClkResponse) {
                        final S7PayloadUserDataItemClkResponse payloadItem =
                            (S7PayloadUserDataItemClkResponse) payloadItems.get(index);
                        responseCode = decodeResponseCode(payloadItem.getReturnCode());
                        dt = payloadItem.getTimeStamp();
                    } else if (payloadItems.get(index) instanceof S7PayloadUserDataItemClkFResponse) {
                        final S7PayloadUserDataItemClkFResponse payloadItem =
                            (S7PayloadUserDataItemClkFResponse) payloadItems.get(index);
                        responseCode = decodeResponseCode(payloadItem.getReturnCode());
                        dt = payloadItem.getTimeStamp();
                    } else {
                        throw new PlcRuntimeException("unknown date-time type.");
                    }

                    List<PlcValue> plcValues = new LinkedList<>();
                    plcValues.add(PlcDATE_AND_LTIME.of(LocalDateTime.of(
                        dt.getYear() + 2000,
                        dt.getMonth(),
                        dt.getDay(),
                        dt.getHour(),
                        dt.getMinutes(),
                        dt.getSeconds(),
                        dt.getMsec() * 1000000)));
                    plcValue = new PlcList(plcValues);
                }

                PlcResponseItem<PlcValue> result = new DefaultPlcResponseItem<>(responseCode, plcValue);
                values.put(tagName, result);
                index++;
            }

            return new DefaultPlcReadResponse(plcReadRequest, values);
        }

        // In all other cases all went well.
        S7PayloadReadVarResponse payload = (S7PayloadReadVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcReadRequest.getNumberOfTags() != payload.getItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        List<S7VarPayloadDataItem> payloadItems = payload.getItems();
        int index = 0;
        PlcResponseCode responseCode;
        PlcValue plcValue;
        for (String tagName : plcReadRequest.getTagNames()) {
            S7Tag tag = (S7Tag) plcReadRequest.getTag(tagName);
            S7VarPayloadDataItem payloadItem = payloadItems.get(index);

            responseCode = decodeResponseCode(payloadItem.getReturnCode());
            plcValue = null;

            if (responseCode == PlcResponseCode.OK) {
                try {
                    plcValue = parsePlcValue(tag, payloadItem.getData());
                } catch (Exception e) {
                    throw new PlcProtocolException("Error decoding PlcValue", e);
                }
            }

            PlcResponseItem<PlcValue> result = new DefaultPlcResponseItem<>(responseCode, plcValue);
            values.put(tagName, result);
            index++;
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);
    }

    private PlcResponse decodeWriteResponse(S7Message responseMessage, PlcWriteRequest plcWriteRequest) throws PlcProtocolException {
        Map<String, PlcResponseCode> responses = new HashMap<>();
        short errorClass;
        short errorCode;

        if (responseMessage instanceof S7MessageResponseData) {
            S7MessageResponseData messageResponseData = (S7MessageResponseData) responseMessage;
            errorClass = messageResponseData.getErrorClass();
            errorCode = messageResponseData.getErrorCode();
        } else if (responseMessage instanceof S7MessageResponse) {
            S7MessageResponse messageResponse = (S7MessageResponse) responseMessage;
            errorClass = messageResponse.getErrorClass();
            errorCode = messageResponse.getErrorCode();
        } else if (responseMessage instanceof S7MessageUserData) {
            String tagName = (String) plcWriteRequest.getTagNames().toArray()[0];
            responses.put(tagName, PlcResponseCode.OK);
            return new DefaultPlcWriteResponse(plcWriteRequest, responses);
        } else {
            throw new PlcProtocolException("Unsupported message type " + responseMessage.getClass().getName());
        }
        // If the result contains any form of non-null error code, handle this instead.
        if ((errorClass != 0) || (errorCode != 0)) {
            // This is usually the case if PUT/GET wasn't enabled on the PLC
            if ((errorClass == 129) && (errorCode == 4)) {
                logger.warn("Got an error response from the PLC. This particular response code usually indicates " +
                    "that PUT/GET is not enabled on the PLC.");
                for (String tagName : plcWriteRequest.getTagNames()) {
                    responses.put(tagName, PlcResponseCode.ACCESS_DENIED);
                }
                return new DefaultPlcWriteResponse(plcWriteRequest, responses);
            } else if ((errorClass == 0x85) && (errorCode == 0)) {
                logger.warn("Got an error response from the PLC. This particular response code usually indicates " +
                    "that we sent a too large packet or would be receiving a too large one. " +
                    "Please report this, as this is most probably a bug.");
                for (String tagName : plcWriteRequest.getTagNames()) {
                    responses.put(tagName, PlcResponseCode.INTERNAL_ERROR);
                }
                return new DefaultPlcWriteResponse(plcWriteRequest, responses);
            } else {
                logger.warn("Got an unknown error response from the PLC. Error Class: {}, Error Code {}. " +
                        "We probably need to implement explicit handling for this, so please file a bug-report " +
                        "on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump " +
                        "containing a capture of the communication.",
                    errorClass, errorCode);
                for (String tagName : plcWriteRequest.getTagNames()) {
                    responses.put(tagName, PlcResponseCode.INTERNAL_ERROR);
                }
                return new DefaultPlcWriteResponse(plcWriteRequest, responses);
            }
        }

        // In all other cases all went well.
        S7PayloadWriteVarResponse payload = (S7PayloadWriteVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcWriteRequest.getNumberOfTags() != payload.getItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        List<S7VarPayloadStatusItem> payloadItems = payload.getItems();
        int index = 0;
        for (String tagName : plcWriteRequest.getTagNames()) {
            S7VarPayloadStatusItem payloadItem = payloadItems.get(index);

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            responses.put(tagName, responseCode);
            index++;
        }

        return new DefaultPlcWriteResponse(plcWriteRequest, responses);
    }

    private S7VarPayloadDataItem serializePlcValue(S7Tag tag, PlcValue plcValue) {
        try {
            DataTransportSize transportSize = tag.getDataType().getDataTransportSize();
            int stringLength = (tag instanceof S7StringFixedLengthTag) ? ((S7StringFixedLengthTag) tag).getStringLength() : 254;
            ByteBuffer byteBuffer = null;
            if((tag.getDataType() == TransportSize.BYTE) && (tag.getNumberOfElements() > 1)) {
                byteBuffer = ByteBuffer.allocate(tag.getNumberOfElements());
                byteBuffer.put(plcValue.getRaw());
            } else {
                for (int i = 0; i < tag.getNumberOfElements(); i++) {
                    int lengthInBits = DataItem.getLengthInBits(plcValue.getIndex(i), tag.getDataType().getDataProtocolId(), s7DriverContext.getControllerType(), stringLength);
                    // Cap the length of the string with the maximum allowed size.
                    if (tag.getDataType() == TransportSize.STRING) {
                        lengthInBits = Math.min(lengthInBits, (stringLength * 8) + 16);
                    } else if (tag.getDataType() == TransportSize.WSTRING) {
                        lengthInBits = Math.min(lengthInBits, (stringLength * 16) + 32);
                    } else if (tag.getDataType() == TransportSize.S5TIME) {
                        lengthInBits = lengthInBits * 8;
                    }
                    final WriteBufferByteBased writeBuffer = new WriteBufferByteBased((int) Math.ceil(((float) lengthInBits) / 8.0f));
                    DataItem.staticSerialize(writeBuffer, plcValue.getIndex(i), tag.getDataType().getDataProtocolId(), s7DriverContext.getControllerType(), stringLength);
                    // Allocate enough space for all items.
                    if (byteBuffer == null) {
                        // TODO: This logic will cause problems when reading arrays of strings.
                        byteBuffer = ByteBuffer.allocate(writeBuffer.getBytes().length * tag.getNumberOfElements());
                    }
                    byteBuffer.put(writeBuffer.getBytes());
                }
            }
            if (byteBuffer != null) {
                byte[] data = byteBuffer.array();
                return new S7VarPayloadDataItem(DataTransportErrorCode.OK, transportSize, data);
            }
        } catch (SerializationException e) {
            logger.warn("Error serializing tag item of type: '{}'", tag.getDataType().name(), e);
        }
        return null;
    }

    private PlcValue parsePlcValue(S7Tag tag, byte[] data) {
        ReadBuffer readBuffer = new ReadBufferByteBased(data);
        try {
            int stringLength = (tag instanceof S7StringFixedLengthTag) ? ((S7StringFixedLengthTag) tag).getStringLength() : 254;
            if (tag.getNumberOfElements() == 1) {
                // TODO: Pass the type of plc into the parse function ...
                return DataItem.staticParse(readBuffer, tag.getDataType().getDataProtocolId(),
                    s7DriverContext.getControllerType(), stringLength);
            } else {
                // In case of reading an array of bytes, make use of our simpler PlcRawByteArray as the user is
                // probably expecting to process the read raw data.
                if(tag.getDataType() == TransportSize.BYTE) {
                    return new PlcRawByteArray(data);
                } else {
                    // Fetch all
                    final PlcValue[] resultItems = IntStream.range(0, tag.getNumberOfElements()).mapToObj(i -> {
                        try {
                            return DataItem.staticParse(readBuffer, tag.getDataType().getDataProtocolId(),
                                s7DriverContext.getControllerType(), stringLength);
                        } catch (ParseException e) {
                            logger.warn("Error parsing tag item of type: '{}' (at position {}})", tag.getDataType().name(), i, e);
                        }
                        return null;
                    }).toArray(PlcValue[]::new);
                    return DefaultPlcValueHandler.of(tag, resultItems);
                }
            }
        } catch (ParseException e) {
            logger.warn("Error parsing tag item of type: '{}'", tag.getDataType().name(), e);
        }
        return null;
    }

    /**
     * Helper to convert the return codes returned from the S7 into one of our standard
     * PLC4X return codes
     *
     * @param dataTransportErrorCode S7 return code
     * @return PLC4X return code.
     */
    private PlcResponseCode decodeResponseCode(DataTransportErrorCode dataTransportErrorCode) {
        if (dataTransportErrorCode == null) {
            return PlcResponseCode.INTERNAL_ERROR;
        }
        switch (dataTransportErrorCode) {
            case OK:
                return PlcResponseCode.OK;
            case NOT_FOUND:
                // It seems the S7 devices return NOT_FOUND if for example we try to access a DB number
                // which doesn't exist. In other protocols we all map that to invalid address.
            case INVALID_ADDRESS:
                return PlcResponseCode.INVALID_ADDRESS;
            case DATA_TYPE_NOT_SUPPORTED:
                return PlcResponseCode.INVALID_DATATYPE;
            default:
                return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    /**
     * Little helper method to parse Siemens article numbers and extract the type of controller.
     *
     * @param articleNumber article number string.
     * @return type of controller.
     */
    private ControllerType decodeControllerType(String articleNumber) {
        if (!articleNumber.startsWith("6ES7 ")) {
            return ControllerType.ANY;
        }
        String model = articleNumber.substring(articleNumber.indexOf(' ') + 1, articleNumber.indexOf(' ') + 2);
        switch (model) {
            case "2":
                return ControllerType.S7_1200;
            case "5":
                return ControllerType.S7_1500;
            case "3":
                return ControllerType.S7_300;
            case "4":
                return ControllerType.S7_400;
            default:
                if (logger.isInfoEnabled()) {
                    logger.info("Looking up unknown article number {}", articleNumber);
                }
                return ControllerType.ANY;
        }
    }

    /**
     * Currently we only support the S7 Any type of addresses. This helper simply converts the S7Tag
     * from PLC4X into S7Address objects.
     *
     * @param tag S7Tag instance we need to convert into an S7Address
     * @return the S7Address
     */
    protected S7Address encodeS7Address(PlcTag tag) {
        if (!(tag instanceof S7Tag)) {
            throw new PlcRuntimeException("Unsupported address type " + tag.getClass().getName());
        }
        S7Tag s7Tag = (S7Tag) tag;
        TransportSize transportSize = s7Tag.getDataType();
        int numElements = s7Tag.getNumberOfElements();
        // For these date-types we have to convert the requests to simple byte-array requests
        // As otherwise the S7 will deny them with "Data type not supported" replies.
        if (transportSize == TransportSize.STRING) {
            transportSize = TransportSize.CHAR;
            int stringLength = (s7Tag instanceof S7StringFixedLengthTag) ? ((S7StringFixedLengthTag) s7Tag).getStringLength() : 254;
            numElements = numElements * (stringLength + 2);
        } else if (transportSize == TransportSize.WSTRING) {
            transportSize = TransportSize.CHAR;
            int stringLength = (s7Tag instanceof S7StringFixedLengthTag) ? ((S7StringFixedLengthTag) s7Tag).getStringLength() : 254;
            numElements = numElements * (stringLength + 2) * 2;
        }
        if (transportSize.getCode() == 0x00) {
            numElements = numElements * transportSize.getSizeInBytes();
            transportSize = TransportSize.BYTE;
        }
        return new S7AddressAny(transportSize, numElements, s7Tag.getBlockNumber(),
            s7Tag.getMemoryArea(), s7Tag.getByteOffset(), s7Tag.getBitOffset());
    }

    private boolean isConnected() {
        return conversationContext.getChannel().attr(S7HMuxImpl.IS_CONNECTED).get();
        //return true;
    }

    private boolean isPrimaryChannel() {
        return conversationContext.getChannel().attr(S7HMuxImpl.IS_PRIMARY).get() == null || conversationContext.getChannel().attr(S7HMuxImpl.IS_PRIMARY).get();
    }

    private void setChannelFeatures() {
        conversationContext.getChannel().attr(S7HMuxImpl.READ_TIME_OUT).set(s7DriverContext.getReadTimeout());
        conversationContext.getChannel().attr(S7HMuxImpl.IS_PING_ACTIVE).set(s7DriverContext.getPing());
        conversationContext.getChannel().attr(S7HMuxImpl.PING_TIME).set(s7DriverContext.getPingTime());
        conversationContext.getChannel().attr(S7HMuxImpl.RETRY_TIME).set(s7DriverContext.getRetryTime());
    }


    private boolean isFeatureSupported() {
        return (s7DriverContext.getControllerType() == ControllerType.S7_300) ||
            (s7DriverContext.getControllerType() == ControllerType.S7_400);
    }

    private CompletableFuture<S7MessageUserData> reassembledMessage(short sequenceNumber, List<PlcValue> plcValues) {

        CompletableFuture<S7MessageUserData> future = new CompletableFuture<>();

        //TODO: We need to verify that the returned tpdu id is the same in the response.
        int tpduId = getTpduId();

        TPKTPacket request = createSzlReassembledRequest(tpduId, sequenceNumber);

        conversationContext.sendRequest(request)
            .onTimeout(e -> {
                logger.warn("Timeout during Connection establishing, closing channel...");
                //context.getChannel().close();
            })
            .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
            .unwrap(TPKTPacket::getPayload)
            .only(COTPPacketData.class)
            .unwrap(COTPPacketData::getPayload)
            .only(S7MessageUserData.class)
            .check(p -> p.getPayload() instanceof S7PayloadUserData)
            .handle(future::complete);

        return future;
    }

    /*
     *
     */
    private TPKTPacket createSzlReassembledRequest(int tpduId, short sequenceNumber) {
        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(tpduId, new S7ParameterUserData(List.of(
            new S7ParameterUserDataItemCPUFunctions((short) 0x12, (byte) 0x4, (byte) 0x4, (short) 0x01, sequenceNumber, (short) 0x00, (short) 0x00, 0)
        )), new S7PayloadUserData(List.of(
            new S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest(
                DataTransportErrorCode.NOT_FOUND,
                DataTransportSize.NULL,
                0x00)
        )));
        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (byte) 2);
        return new TPKTPacket(cotpPacketData);
    }

    private CompletableFuture<S7MessageUserData> reassembledAlarmEvents(short sequenceNumber) {
        CompletableFuture<S7MessageUserData> future = new CompletableFuture<>();

        //TODO: We need to verify that the returned tpdu id is the same in the response.
        int tpduId = getTpduId();

        TPKTPacket request = createAlarmQueryReassembledRequest(tpduId, sequenceNumber);

        conversationContext.sendRequest(request)
            .onTimeout(e -> {
                logger.warn("Timeout during Connection establishing, closing channel...");
                //context.getChannel().close();
            })
            .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
            .unwrap(TPKTPacket::getPayload)
            .only(COTPPacketData.class)
            .unwrap(COTPPacketData::getPayload)
            .only(S7MessageUserData.class)
            .check(p -> p.getPayload() instanceof S7PayloadUserData)
            .handle(future::complete);

        return future;
    }

    //TODO: S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest to S7PayloadUserDataItemCpuFunctionAlarmQueryNoDataRequest 
    private TPKTPacket createAlarmQueryReassembledRequest(int tpduId, short sequenceNumber) {
        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(tpduId, new S7ParameterUserData(List.of(
            new S7ParameterUserDataItemCPUFunctions((short) 0x12, (byte) 0x4, (byte) 0x4, (short) 0x13, sequenceNumber, (short) 0x00, (short) 0x00, 0)
        )), new S7PayloadUserData(List.of(
            new S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest(
                DataTransportErrorCode.NOT_FOUND,
                DataTransportSize.NULL,
                0x00)
        )));
        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (byte) 2);
        return new TPKTPacket(cotpPacketData);
    }

    private int getTpduId() {
        int tpduId = tpduGenerator.getAndIncrement();
        // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
        if (tpduGenerator.get() == 0xFFFF) {
            tpduGenerator.set(1);
        }
        return tpduId;
    }


    /**
     * A generic purpose error handler which terminates transaction and calls back given future with error message.
     */
    static class TransactionErrorCallback<T, E extends Throwable> implements Consumer<TimeoutException>, BiConsumer<TPKTPacket, E> {

        private final CompletableFuture<T> future;
        private final RequestTransactionManager.RequestTransaction transaction;

        TransactionErrorCallback(CompletableFuture<T> future, RequestTransactionManager.RequestTransaction transaction) {
            this.future = future;
            this.transaction = transaction;
        }

        @Override
        public void accept(TimeoutException e) {
            try {
                transaction.endRequest();
            } catch (Exception ex) {
                logger.info(ex.getMessage());
            }
            future.completeExceptionally(e);
        }

        @Override
        public void accept(TPKTPacket tpktPacket, E e) {
            try {
                transaction.endRequest();
            } catch (Exception ex) {
                logger.info(ex.getMessage());
            }
            future.completeExceptionally(e);
        }
    }

    protected CompletableFuture<Map<S7StringVarLengthTag, StringSizes>> getStringSizes(PlcTagRequest request) {
        CompletableFuture<Map<S7StringVarLengthTag, StringSizes>> future = new CompletableFuture<>();

        // Build a request to read the length information for every var-length string in the request.
        List<S7StringVarLengthTag> varLengthStringTags = request.getTags().stream()
            .filter(plcTag -> plcTag instanceof S7StringVarLengthTag)
            .map(plcTag -> (S7StringVarLengthTag) plcTag)
            .collect(Collectors.toList());
        List<S7VarRequestParameterItem> stringFields = new ArrayList<>(varLengthStringTags.size());
        for (S7StringVarLengthTag varLengthStringTag : varLengthStringTags) {
            // For STRING, the header is 2 bytes (first byte contains the max length and the second the actual length)
            if (varLengthStringTag.getDataType() == TransportSize.STRING) {
                stringFields.add(new S7VarRequestParameterItemAddress(
                    new S7AddressAny(
                        TransportSize.BYTE,
                        2,
                        varLengthStringTag.getBlockNumber(),
                        MemoryArea.DATA_BLOCKS,
                        varLengthStringTag.getByteOffset(),
                        varLengthStringTag.getBitOffset()
                    )));
            }
            // For WSTRING, the header is 4 bytes (first word contains the max length and the second the actual length)
            else if (varLengthStringTag.getDataType() == TransportSize.WSTRING) {
                stringFields.add(new S7VarRequestParameterItemAddress(
                    new S7AddressAny(
                        TransportSize.BYTE,
                        4,
                        varLengthStringTag.getBlockNumber(),
                        MemoryArea.DATA_BLOCKS,
                        varLengthStringTag.getByteOffset(),
                        varLengthStringTag.getBitOffset()
                    )));
            } else {
                throw new PlcInvalidTagException("Only STRING and WSTRING allowed here.");
            }
        }
        final S7MessageRequest readRequest = new S7MessageRequest(
            getTpduId(), new S7ParameterReadVarRequest(stringFields), null);

        // Read the max length and actual size for each of the var-length strings.
        CompletableFuture<S7Message> resolveSizesRequestFuture = sendInternal(readRequest);
        resolveSizesRequestFuture.whenComplete((s7Message, throwable) -> {
            if (throwable != null) {
                future.completeExceptionally(new PlcProtocolException("Error resolving string sizes", throwable));
                return;
            }

            Map<S7StringVarLengthTag, StringSizes> stringLengths = new HashMap<>(varLengthStringTags.size());
            S7PayloadReadVarResponse getLengthsResponse = (S7PayloadReadVarResponse) s7Message.getPayload();
            int curItemIndex = 0;
            for (S7StringVarLengthTag varLengthStringTag : varLengthStringTags) {
                S7VarPayloadDataItem s7VarPayloadDataItem = getLengthsResponse.getItems().get(curItemIndex);
                // Something went wrong.
                // We simply don't save the length information, and then we'll treat it as INVALID_ADDRESS later on
                // in the code.
                if(s7VarPayloadDataItem.getReturnCode() != DataTransportErrorCode.OK) {
                    continue;
                }
                ReadBufferByteBased readBuffer = new ReadBufferByteBased(s7VarPayloadDataItem.getData());
                try {
                    if (varLengthStringTag.getDataType() == TransportSize.STRING) {
                        int maxChars = readBuffer.readUnsignedInt("maxLength", 8);
                        int actualChars = readBuffer.readUnsignedInt("maxLength", 8);
                        stringLengths.put(varLengthStringTag, new StringSizes(maxChars, actualChars));
                    } else if (varLengthStringTag.getDataType() == TransportSize.WSTRING) {
                        int maxChars = readBuffer.readUnsignedInt("maxLength", 16);
                        int actualChars = readBuffer.readUnsignedInt("maxLength", 16);
                        stringLengths.put(varLengthStringTag, new StringSizes(maxChars, actualChars));
                    } else {
                        throw new PlcInvalidTagException("Only STRING and WSTRING allowed here.");
                    }
                } catch (ParseException e) {
                    throw new PlcInvalidTagException("Error parsing var-length string actual lengths.");
                }
            }

            future.complete(stringLengths);
        });

        return future;
    }

    public static class StringSizes {

        private final int maxLength;
        private final int curLength;

        public StringSizes(int maxLength, int curLength) {
            this.maxLength = maxLength;
            this.curLength = curLength;
        }

        public int getMaxLength() {
            return maxLength;
        }

        public int getCurLength() {
            return curLength;
        }

    }

}
