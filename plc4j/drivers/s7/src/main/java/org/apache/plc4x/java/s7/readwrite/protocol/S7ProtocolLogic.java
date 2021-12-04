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
package org.apache.plc4x.java.s7.readwrite.protocol;

import org.apache.plc4x.java.s7.readwrite.utils.S7PlcSubscriptionHandle;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import io.netty.util.ByteProcessor;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.s7.readwrite.field.S7StringField;
import org.apache.plc4x.java.s7.readwrite.io.DataItemIO;
import org.apache.plc4x.java.s7.readwrite.types.*;
import org.apache.plc4x.java.s7.readwrite.field.S7Field;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.time.Duration;
import java.util.ArrayList;
import java.util.BitSet;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Queue;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;
import java.util.concurrent.atomic.AtomicReferenceArray;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.logging.Level;
import java.util.stream.IntStream;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.s7.events.S7CyclicEvent;
import org.apache.plc4x.java.s7.readwrite.field.S7AckField;
import org.apache.plc4x.java.s7.readwrite.field.S7SzlField;
import org.apache.plc4x.java.s7.readwrite.field.S7SubscriptionField;
import org.apache.plc4x.java.s7.readwrite.io.S7PayloadUserDataItemIO;
import org.apache.plc4x.java.s7.utils.S7ParamErrorCode;
import org.apache.plc4x.java.spi.ConversationContext.ContextHandler;
import org.apache.plc4x.java.spi.ConversationContext.SendRequestContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcUnsubscriptionRequest;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcSINT;
import org.apache.plc4x.java.spi.values.PlcUINT;

/**
 * The S7 Protocol states that there can not be more then {min(maxAmqCaller, maxAmqCallee} "ongoing" requests.
 * So we need to limit those.
 * Thus, each request goes to a Work Queue and this Queue ensures, that only 3 are open at the same time.
 * 
 * TODO: Evaluate the implementation of the "Visitor"  design pattern in the mspec code generator.
 * 
 */
public class S7ProtocolLogic extends Plc4xProtocolBase<TPKTPacket> {

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private final Logger logger = LoggerFactory.getLogger(S7ProtocolLogic.class);
    private final AtomicInteger tpduGenerator = new AtomicInteger(10);

    /*
     * Take into account that the size of this buffer depends on the final device.
     * S7-300 goes from 20 to 300 and for S7-400 it goes from 300 to 10000.
     * Depending on the configuration of the alarm system, a large number of
     * them should be expected when starting the connection.
     * (Examples of this are PCS7 and Braumat).
     * Alarm filtering, ack, etc. must be performed by the client application.
    */
    private final BlockingQueue eventqueue = new ArrayBlockingQueue<>(1024);
    
    private final S7ProtocolEventLogic EventLogic = new S7ProtocolEventLogic(eventqueue);
    private final S7PlcSubscriptionHandle modeHandle = new S7PlcSubscriptionHandle(EventType.MODE,EventLogic);
    private final S7PlcSubscriptionHandle sysHandle = new S7PlcSubscriptionHandle(EventType.SYS,EventLogic);
    private final S7PlcSubscriptionHandle usrHandle = new S7PlcSubscriptionHandle(EventType.USR,EventLogic);
    private final S7PlcSubscriptionHandle almHandle = new S7PlcSubscriptionHandle(EventType.ALM,EventLogic);
    private final S7PlcSubscriptionHandle cycHandle = new S7PlcSubscriptionHandle(EventType.CYC,EventLogic);    

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
        EventLogic.start();
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

        // Only the TCP transport supports login.
        logger.info("S7 Driver running in ACTIVE mode.");
        logger.debug("Sending COTP Connection Request");
        // Open the session on ISO Transport Protocol first.
        TPKTPacket packet = new TPKTPacket(createCOTPConnectionRequest(
            s7DriverContext.getCalledTsapId(), s7DriverContext.getCallingTsapId(), s7DriverContext.getCotpTpduSize()));

        context.sendRequest(packet)
            .onTimeout(e -> {
                logger.warn("Timeout during Connection establishing, closing channel...");
                context.getChannel().close();
            })
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .check(p -> p.getPayload() instanceof COTPPacketConnectionResponse)
            .unwrap(p -> (COTPPacketConnectionResponse) p.getPayload())
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
                            if (s7DriverContext.getControllerType() != S7ControllerType.ANY) {
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
                            .check(p -> p.getPayload() instanceof COTPPacketData)
                            .unwrap(p -> ((COTPPacketData) p.getPayload()))
                            .check(p -> p.getPayload() instanceof S7MessageUserData)
                            .unwrap(p -> ((S7MessageUserData) p.getPayload()))
                            .check(p -> p.getPayload() instanceof S7PayloadUserData)
                            .handle(messageUserData -> {
                                logger.debug("Got S7 Identification Response");
                                S7PayloadUserData payloadUserData = (S7PayloadUserData) messageUserData.getPayload();
                                extractControllerTypeAndFireConnected(context, payloadUserData);
                            });
                    });
            });
    }

    /*
    * Support read varibles from PLC (S7Field) and system status list (S7SzlField).
    * 
    */
    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<S7VarRequestParameterItem> requestItems = new ArrayList<>(request.getNumberOfFields());
     
        if (request.getFields().get(0) instanceof S7SzlField) {
            S7SzlField szlfield = (S7SzlField) request.getFields().get(0);
            final S7MessageUserData s7SslMessageRequest = new S7MessageUserData(1, new S7ParameterUserData(new S7ParameterUserDataItem[]{
                new S7ParameterUserDataItemCPUFunctions((short) 0x11, (byte) 0x4, (byte) 0x4, (short) 0x01, (short) 0x00, null, null, null)
            }), new S7PayloadUserData(new S7PayloadUserDataItem[]{
                new S7PayloadUserDataItemCpuFunctionReadSzlRequest(DataTransportErrorCode.OK, 
                        DataTransportSize.OCTET_STRING, 
                        0x04,
                        new SzlId(SzlModuleTypeClass.enumForValue((byte) ((szlfield.getSzlId() & 0xf000) >> 12)), 
                                (byte) ((szlfield.getSzlId() & 0x0f00) >> 8), 
                                SzlSublist.enumForValue((short) (szlfield.getSzlId() & 0x00ff))), 
                        szlfield.getIndex())
            }));         
            return toPlcReadResponse(readRequest, readInternal(s7SslMessageRequest));   
            
        } else if (request.getFields().get(0) instanceof S7AckField) {
            List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfFields());
            List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfFields());
            encodeAlarmAckRequest(request, parameterItems, payloadItems);
            final S7MessageUserData s7MessageRequest = new S7MessageUserData(-1, 
                new S7ParameterUserData(parameterItems.toArray(new S7ParameterUserDataItem[0])),
                new S7PayloadUserData(payloadItems.toArray(new S7PayloadUserDataItemCpuFunctionAlarmAckRequest[0])));                    
            return toPlcReadResponse(readRequest, readInternal(s7MessageRequest ));               
        }
    
     
        for (PlcField field : request.getFields()) {
            requestItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(field)));
        }        

        // Create a read request template.
        // tpuId will be inserted before sending in #readInternal so we insert -1 as dummy here
        final S7MessageRequest s7MessageRequest = new S7MessageRequest(-1,
            new S7ParameterReadVarRequest(requestItems.toArray(new S7VarRequestParameterItem[0])),
            null);
        
        // Just send a single response and chain it as Response
        return toPlcReadResponse(readRequest, readInternal(s7MessageRequest));
    }

    /** Maps the S7ReadResponse of a PlcReadRequest to a PlcReadResponse 
     * 
     * 
     * @param readRequest
     * @param response
     * @return 
     */
    private CompletableFuture<PlcReadResponse> toPlcReadResponse(PlcReadRequest readRequest, CompletableFuture<S7Message> response) {
        CompletableFuture<PlcReadResponse>  future = new CompletableFuture<>();        
        PlcReadResponse plcitems;    
        try {
            plcitems = (PlcReadResponse) decodeReadResponse(response.get(), readRequest);
            future.complete(plcitems);
        } catch (Exception ex) {
            logger.info(ex.toString());
        }        
        return future;
    }

    /**
     * Sends one Read over the Wire and internally returns the Response
     * Do sending of normally sized single-message request.
     *
     * Assumes that the {@link S7MessageRequest} and its expected {@link S7MessageResponseData}
     * and does not further check that!
     */
    private CompletableFuture<S7Message> readInternal(S7MessageRequest request) {
        CompletableFuture<S7Message> future = new CompletableFuture<>();
        int tpduId = tpduGenerator.getAndIncrement();
        // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
        if(tpduGenerator.get() == 0xFFFF) {
            tpduGenerator.set(1);
        }

        // Create a new Request with correct tpuId (is not known before)
        S7MessageRequest s7MessageRequest = new S7MessageRequest(tpduId, request.getParameter(), request.getPayload());

        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null, s7MessageRequest, true, (short) tpduId));
        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(tpktPacket)
            .onTimeout(new TransactionErrorCallback<>(future, transaction))
            .onError(new TransactionErrorCallback<>(future, transaction))
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .check(p -> p.getPayload() instanceof COTPPacketData)
            .unwrap(p -> (COTPPacketData) p.getPayload())
            .check(p -> p.getPayload()  != null)
            .unwrap(COTPPacket::getPayload)
            .check(p -> p.getTpduReference() == tpduId)
            .handle(p -> {
                future.complete(p);
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }
    
    private CompletableFuture<S7Message> readInternal(S7MessageUserData request) {
        CompletableFuture<S7Message> future = new CompletableFuture<>();
        int tpduId = tpduGenerator.getAndIncrement();
        // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
        if(tpduGenerator.get() == 0xFFFF) {
            tpduGenerator.set(1);
        }

        // Create a new Request with correct tpuId (is not known before)
        S7MessageUserData s7MessageRequest = new S7MessageUserData(tpduId, request.getParameter(), request.getPayload());

        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null, s7MessageRequest, true, (short) tpduId));
        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(tpktPacket)
            .onTimeout(new TransactionErrorCallback<>(future, transaction))
            .onError(new TransactionErrorCallback<>(future, transaction))
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .check(p -> p.getPayload() instanceof COTPPacketData)
            .unwrap(p -> (COTPPacketData) p.getPayload())
            .check(p -> p.getPayload()  != null)
            .unwrap(COTPPacket::getPayload)
            .check(p -> p.getTpduReference() == tpduId)
            .handle(p -> {
                future.complete(p);
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }    
    
    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<S7VarRequestParameterItem> parameterItems = new ArrayList<>(request.getNumberOfFields());
        List<S7VarPayloadDataItem> payloadItems = new ArrayList<>(request.getNumberOfFields());
        for (String fieldName : request.getFieldNames()) {
            final S7Field field = (S7Field) request.getField(fieldName);
            final PlcValue plcValue = request.getPlcValue(fieldName);
            parameterItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(field)));
            payloadItems.add(serializePlcValue(field, plcValue));
        }
        final int tpduId = tpduGenerator.getAndIncrement();
        // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
        if(tpduGenerator.get() == 0xFFFF) {
            tpduGenerator.set(1);
        }

        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
            new S7MessageRequest(tpduId,
                new S7ParameterWriteVarRequest(parameterItems.toArray(new S7VarRequestParameterItem[0])),
                new S7PayloadWriteVarRequest(payloadItems.toArray(new S7VarPayloadDataItem[0]))),
            true, (short) tpduId));

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(tpktPacket)
            .onTimeout(new TransactionErrorCallback<>(future, transaction))
            .onError(new TransactionErrorCallback<>(future, transaction))
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .check(p -> p.getPayload() instanceof COTPPacketData)
            .unwrap(p -> ((COTPPacketData) p.getPayload()))
            .unwrap(COTPPacket::getPayload)
            .check(p -> p.getTpduReference() == tpduId)
            .handle(p -> {
                try {
                    future.complete(((PlcWriteResponse) decodeWriteResponse(p, writeRequest)));
                } catch (PlcProtocolException e) {
                    logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {

        CompletableFuture<PlcSubscriptionResponse> response =  new CompletableFuture<>();
        HashMap<String, PlcSubscriptionResponse> valuesResponse = new HashMap<>();


        CompletableFuture<S7Message> future = new CompletableFuture<>();

        HashMap<String, CompletableFuture<S7Message>> futures = new HashMap<>();
        subscriptionRequest.getFieldNames().forEach(fieldname -> futures.put(fieldname,  new CompletableFuture<>()));


        DefaultPlcSubscriptionRequest request = (DefaultPlcSubscriptionRequest) subscriptionRequest;

        List<S7ParameterUserDataItem> parameterItems = new ArrayList<>(request.getNumberOfFields());
        List<S7PayloadUserDataItem> payloadItems = new ArrayList<>(request.getNumberOfFields());

        CompletableFuture<Void> maintask;
        maintask = CompletableFuture.
                allOf(futures.values().toArray(new CompletableFuture[0]));
        
        Thread t1 = new Thread(()->{

            S7ParameterUserData txParameters = null;
            S7PayloadUserData txPayloads = null;

            for (String fieldName : request.getFieldNames()) {
                final DefaultPlcSubscriptionField sf = (DefaultPlcSubscriptionField) request.getField(fieldName);
                final S7SubscriptionField  field =  (S7SubscriptionField) sf.getPlcField();

                final int tpduId = tpduGenerator.getAndIncrement();
                // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
                if(tpduGenerator.get() == 0xFFFF) {
                    tpduGenerator.set(1);
                }

                switch(field.getFieldtype()){
                    case EVENT_SUBSCRIPTION: {
                        parameterItems.clear();
                        payloadItems.clear();
                        encodeEventSubcriptionRequest(request, parameterItems, payloadItems);
                        txParameters = new S7ParameterUserData(parameterItems.toArray(new S7ParameterUserDataItem[0]));
                        txPayloads = new S7PayloadUserData(payloadItems.toArray(new S7PayloadUserDataItemCpuFunctionMsgSubscription[0]));

                        //TODO: Diferent message type in the Payload
                        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
                                new S7MessageUserData(tpduId,                        
                                        // new S7ParameterUserData(parameterItems.toArray(new S7ParameterUserDataItem[0])),
                                        // new S7PayloadUserData(payloadItems.toArray(new S7PayloadUserDataItemCpuFunctionMsgSubscription[0]))),
                                        txParameters,
                                        txPayloads),
                                //new S7PayloadUserData(payloadItems.toArray(new S7PayloadUserDataItemCpuFunctionAlarmAckRequest[0]))),
                                true, (short) tpduId));

                        // Start a new request-transaction (Is ended in the response-handler)
                        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
                        transaction.submit(() -> context.sendRequest(tpktPacket)
                                .onTimeout(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                .onError(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                                .check(p -> p.getPayload() instanceof COTPPacketData)
                                .unwrap(p -> ((COTPPacketData) p.getPayload()))
                                .unwrap(COTPPacket::getPayload)
                                .check(p -> p.getTpduReference() == tpduId)
                                .handle(p -> {
                                    try {
                                        futures.get(fieldName).complete(p);
                                    } catch (Exception e) {
                                        logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                                    }
                                    // Finish the request-transaction.
                                    transaction.endRequest();
                                }));



                        //return toPlcSubscriptionResponse(subscriptionRequest, future);
                    }
                    break;
                    case EVENT_UNSUBSCRIPTION:
                        //encodeEventUnSubcriptionRequest(msg, out);
                    break;
                    case ALARM_ACK:
                        //encodeAlarmAckRequest(request, parameterItems, payloadItems);
                    break;
                    case ALARM_QUERY:{
                        parameterItems.clear();
                        payloadItems.clear();
                        encodeAlarmQueryRequest(request, parameterItems, payloadItems);
                        txParameters = new S7ParameterUserData(parameterItems.toArray(new S7ParameterUserDataItem[0]));
                        txPayloads = new S7PayloadUserData(payloadItems.toArray(new S7PayloadUserDataItemCpuFunctionAlarmQueryRequest[0]));


                        //TODO: Diferent message type in the Payload
                        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
                                new S7MessageUserData(tpduId,
                                        txParameters,
                                        txPayloads),
                                true, (short) tpduId));

                        // Start a new request-transaction (Is ended in the response-handler)

                        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
                        transaction.submit(() -> context.sendRequest(tpktPacket)
                                .onTimeout(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                .onError(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                                .check(p -> p.getPayload() instanceof COTPPacketData)
                                .unwrap(p -> ((COTPPacketData) p.getPayload()))
                                .unwrap(COTPPacket::getPayload)
                                .check(p -> p.getTpduReference() == tpduId)
                                .handle(p -> {
                                    try {
                                        futures.get(fieldName).complete(p);
                                    } catch (Exception e) {
                                        logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                                    }
                                    // Finish the request-transaction.
                                    transaction.endRequest();
                                }));

                    }
                    break;
                    case CYCLIC_SUBSCRIPTION:{

                        parameterItems.clear();
                        payloadItems.clear();
                        encodeCycledS7ANYSubscriptionRequest(request.getField(fieldName), parameterItems, payloadItems);
                        if (!parameterItems.isEmpty() && !payloadItems.isEmpty()){
                            txParameters = new S7ParameterUserData(parameterItems.toArray(new S7ParameterUserDataItem[0]));
                            txPayloads = new S7PayloadUserData(payloadItems.toArray(new S7PayloadUserDataItemCyclicServicesSubscribeRequest[0]));

                            //TODO: Diferent message type in the Payload
                            TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
                                    new S7MessageUserData(tpduId,
                                            txParameters,
                                            txPayloads),
                                    true, (short) tpduId));


                            // Start a new request-transaction (Is ended in the response-handler)

                            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
                            transaction.submit(() -> context.sendRequest(tpktPacket)
                                    .onTimeout(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                    .onError(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                    .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                                    .check(p -> p.getPayload() instanceof COTPPacketData)
                                    .unwrap(p -> ((COTPPacketData) p.getPayload()))
                                    .unwrap(COTPPacket::getPayload)
                                    .check(p -> p.getTpduReference() == tpduId)
                                    .handle(p -> {
                                        try {
                                            futures.get(fieldName).complete(p);
                                        } catch (Exception e) {
                                            logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                                        }
                                        // Finish the request-transaction.
                                        transaction.endRequest();
                                    }));
                        }
                    }
                    break;
                    case CYCLIC_DB_SUBSCRIPTION:{

                        parameterItems.clear();
                        payloadItems.clear();
                        encodeCycledDBREADSubscriptionRequest(request.getField(fieldName), parameterItems, payloadItems);
                        if (!parameterItems.isEmpty() && !payloadItems.isEmpty()){
                            txParameters = new S7ParameterUserData(parameterItems.toArray(new S7ParameterUserDataItem[0]));
                            txPayloads = new S7PayloadUserData(payloadItems.toArray(new S7PayloadUserDataItemCyclicServicesSubscribeRequest[0]));

                            //TODO: Diferent message type in the Payload
                            TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
                                    new S7MessageUserData(tpduId,
                                            txParameters,
                                            txPayloads),
                                    true, (short) tpduId));


                            // Start a new request-transaction (Is ended in the response-handler)

                            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
                            transaction.submit(() -> context.sendRequest(tpktPacket)
                                    .onTimeout(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                    .onError(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                    .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                                    .check(p -> p.getPayload() instanceof COTPPacketData)
                                    .unwrap(p -> ((COTPPacketData) p.getPayload()))
                                    .unwrap(COTPPacket::getPayload)
                                    .check(p -> p.getTpduReference() == tpduId)
                                    .handle(p -> {
                                        try {
                                            futures.get(fieldName).complete(p);
                                        } catch (Exception e) {
                                            logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                                        }
                                        // Finish the request-transaction.
                                        transaction.endRequest();
                                    }));
                        }
                    }
                
                    break;
                    case CYCLIC_UNSUBSCRIPTION:{
                        System.out.println("CYCLIC_UNSUBSCRIPTION");
                        encodeCycledUnSubscriptionRequest(request.getField(fieldName), parameterItems, payloadItems);
                        txParameters = new S7ParameterUserData(parameterItems.toArray(new S7ParameterUserDataItem[0]));
                        txPayloads = new S7PayloadUserData(payloadItems.toArray(new S7PayloadUserDataItemCyclicServicesUnsubscribeRequest[0]));                        
                        
                        //TODO: Diferent message type in the Payload
                        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
                                new S7MessageUserData(tpduId,
                                        txParameters,
                                        txPayloads),
                                true, (short) tpduId)); 
                        
                        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
                        transaction.submit(() -> context.sendRequest(tpktPacket)
                                .onTimeout(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                .onError(new TransactionErrorCallback<>(futures.get(fieldName), transaction))
                                .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                                .check(p -> p.getPayload() instanceof COTPPacketData)
                                .unwrap(p -> ((COTPPacketData) p.getPayload()))
                                .unwrap(COTPPacket::getPayload)
                                .check(p -> p.getTpduReference() == tpduId)
                                .handle(p -> {
                                    try {
                                        futures.get(fieldName).complete(p);
                                    } catch (Exception e) {
                                        logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                                    }
                                    // Finish the request-transaction.
                                    transaction.endRequest();
                                }));                        
                        
                    }
                    break;
                    default:
                }

                try {
                    valuesResponse.put(fieldName, decodeEventSubcriptionRequest(fieldName, subscriptionRequest, futures.get(fieldName).get()));
                } catch (Exception ex) {
                    logger.warn(ex.toString());
                }
            } //Next


            try {
                maintask.get();              
                HashMap<String, ResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
                
                valuesResponse.forEach((s,p)->{
                    if (p != null)
                    values.putAll(((DefaultPlcSubscriptionResponse)p).getValues());                    
                });          
                
                response.complete(new DefaultPlcSubscriptionResponse(subscriptionRequest, values));                

            } catch (Exception ex) {
                java.util.logging.Logger.getLogger(S7ProtocolLogic.class.getName()).log(Level.SEVERE, null, ex);
            }              
        });
        
        t1.start();

        return response;        
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<>();
        DefaultPlcUnsubscriptionRequest request = (DefaultPlcUnsubscriptionRequest) unsubscriptionRequest;

        return future;
    }
    
    private CompletableFuture<PlcSubscriptionResponse>  toPlcSubscriptionResponse(PlcSubscriptionRequest subscriptionRequest, CompletableFuture<S7Message> response) {
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();      
        PlcSubscriptionResponse plcresponse;    
        try {
            plcresponse = decodeEventSubcriptionRequest(null, subscriptionRequest, response.get());
            future.complete(plcresponse);
        } catch (Exception ex) {
            logger.warn(ex.toString());
        }        
        return future;
    }    
    
    private void encodeEventSubcriptionRequest(DefaultPlcSubscriptionRequest request,
                                                List<S7ParameterUserDataItem> parameterItems,
                                                List<S7PayloadUserDataItem> payloadItems){
        byte subsevent = 0;
        for (String fieldName : request.getFieldNames()) {
            if (request.getField(fieldName) instanceof DefaultPlcSubscriptionField){
                PlcField event = ((DefaultPlcSubscriptionField) request.getField(fieldName)).getPlcField();
                if (event instanceof S7SubscriptionField) {
                    if (((S7SubscriptionField) event).getFieldtype() == S7SubscriptionFieldType.EVENT_SUBSCRIPTION)
                    subsevent = (byte) (subsevent | ((S7SubscriptionField) event).getEventtype().getValue());
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
        parameterItems.clear();
        parameterItems.add(parameter);


        S7PayloadUserDataItemCpuFunctionMsgSubscription payload = null;

        if (subsevent > 0) {
            payload = new S7PayloadUserDataItemCpuFunctionMsgSubscription(
                            DataTransportErrorCode.OK,
                            DataTransportSize.OCTET_STRING,
                            0x0a,
                            subsevent,
                            "HmiRtm  ",
                            null,
                            null);
        } else {
            //TODO: Check for ALARM_S (S7300) and ALARM_8 (S7400), maybe we need verify the CPU
            AlarmStateType alarmtype;
            if (s7DriverContext.getControllerType() == S7ControllerType.S7_400){
                alarmtype =AlarmStateType.ALARM_INITIATE;
            } else {
                alarmtype =AlarmStateType.ALARM_S_INITIATE;
            }
            payload = new S7PayloadUserDataItemCpuFunctionMsgSubscription(
                            DataTransportErrorCode.OK,
                            DataTransportSize.OCTET_STRING,
                            0x0c,
                            subsevent,
                            "HmiRtm  ",
                            alarmtype,
                            (short) 0x00);
        }
        payloadItems.clear();
        payloadItems.add(payload);

    }

    private void encodeEventUnSubcriptionRequest(DefaultPlcSubscriptionRequest request,
                                                List<S7ParameterUserDataItem> parameterItems,
                                                 List<S7PayloadUserDataItem> payloadItems){

    }
        
    private void encodeAlarmAckRequest(DefaultPlcReadRequest request,
                                                List<S7ParameterUserDataItem> parameterItems,
                                                 List<S7PayloadUserDataItem> payloadItems){
        
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
        parameterItems.clear();
        parameterItems.add(parameter);
        
        AlarmMessageObjectAckType[] messageObjects = null;
        BitSet bs = new BitSet();
        for (String fieldName : request.getFieldNames()) {
            if (request.getField(fieldName) instanceof S7AckField) {
                PlcField field = request.getField(fieldName);
                if (field instanceof S7AckField) {
                    ArrayList<Integer> arrAlarmIds = ((S7AckField) field).getAlarmIds();
                    ArrayList<Integer> arrAlarmSigs = ((S7AckField) field).getAlarmSigs();
                    messageObjects = new AlarmMessageObjectAckType[arrAlarmIds.size()];
                    int j = 0;
                    for (int i = 0; i < arrAlarmIds.size(); i++) {
                            bs = BitSet.valueOf(new byte[]{arrAlarmSigs.get(i).byteValue()});
                           
                            AlarmMessageObjectAckType messageObject = new AlarmMessageObjectAckType(
                                                                              SyntaxIdType.ALARM_ACKSET,
                                                                              (short) 0,
                                                                              arrAlarmIds.get(i),
                                                                              new State(bs.get(7),bs.get(6),bs.get(5),bs.get(4),bs.get(3),bs.get(2),bs.get(1),bs.get(0)),
                                                                              new State(bs.get(7),bs.get(6),bs.get(5),bs.get(4),bs.get(3),bs.get(2),bs.get(1),bs.get(0))                                
                                                                        );
                            messageObjects[i] = messageObject;
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

        payloadItems.clear();
        payloadItems.add(payload);        
        
    }    
    
    private void encodeAlarmQueryRequest(DefaultPlcSubscriptionRequest request,
                                                List<S7ParameterUserDataItem> parameterItems,
                                                 List<S7PayloadUserDataItem> payloadItems){
        
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
                                
        parameterItems.clear();
        parameterItems.add(parameter);
        

        //TODO: Chequear el tipo dfe larma.
        S7PayloadUserDataItemCpuFunctionAlarmQueryRequest payload =
            new S7PayloadUserDataItemCpuFunctionAlarmQueryRequest(
                    DataTransportErrorCode.OK,
                    DataTransportSize.OCTET_STRING,
                    0x0c,                    
                    SyntaxIdType.ALARM_QUERYREQSET,
                    QueryType.ALARM_8P,
                    AlarmType.ALARM_8);
        
        payloadItems.clear();
        payloadItems.add(payload);           
        
    }

    private void encodeCycledS7ANYSubscriptionRequest(PlcSubscriptionField plcfield,
                                                List<S7ParameterUserDataItem> parameterItems,
                                                 List<S7PayloadUserDataItem> payloadItems){
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
                                
        parameterItems.clear();
        parameterItems.add(parameter);
        
        //TODO: Chequear la asignacionde tipo
        
        PlcField field = ((DefaultPlcSubscriptionField) plcfield).getPlcField(); 
        S7SubscriptionField s7field = (S7SubscriptionField) field;

        CycServiceItemAnyType[] items = new CycServiceItemAnyType[s7field.getS7field().length];
        
        int i=0;
        for (S7Field userfield:s7field.getS7field()){
            items[i] = new CycServiceItemAnyType(
                    (short) 0x0a,
                    (short) 0x10, //S7ANY
                    userfield.getDataType(),
                    userfield.getNumberOfElements(),                    
                    userfield.getBlockNumber(),
                    userfield.getMemoryArea(),
                    (long) (((userfield.getByteOffset() << 3) | (userfield.getBitOffset() & 0x0007)))                    
            );
            i++;
        }
        
        //Length in bytes
        //4+12*items.legth
        int lengthInBytes = 4 + items.length * 12;

        S7PayloadUserDataItemCyclicServicesSubscribeRequest payload =
                new S7PayloadUserDataItemCyclicServicesSubscribeRequest(
                        DataTransportErrorCode.OK,
                        DataTransportSize.OCTET_STRING,
                        lengthInBytes,
                        items.length,
                        s7field.getTimeBase(),
                        s7field.getMultiplier(),
                        items
                );
        
        payloadItems.clear();
        payloadItems.add(payload);
    }
    

    private void encodeCycledDBREADSubscriptionRequest(PlcSubscriptionField plcfield,
                                                List<S7ParameterUserDataItem> parameterItems,
                                                 List<S7PayloadUserDataItem> payloadItems){
        S7ParameterUserDataItemCPUFunctions parameter = new S7ParameterUserDataItemCPUFunctions(
                                                                (short) 0x11,   //Method
                                                                (byte) 0x04,    //FunctionType
                                                                (byte) 0x02,    //FunctionGroup
                                                                (short) 0x05,   //SubFunction
                                                                (short) 0x00,   //SequenceNumber
                                                                null,   //DataUnitReferenceNumber
                                                                null,   //LastDataUnit
                                                                null         //errorCode
                                                    );
                                
        parameterItems.clear();
        parameterItems.add(parameter);
        
        //TODO: Chequear la asignacionde tipo
        
        PlcField field = ((DefaultPlcSubscriptionField) plcfield).getPlcField(); 
        S7SubscriptionField s7field = (S7SubscriptionField) field;

        double cycitemcount = Math.ceil(((double)s7field.getS7field().length / (double) 50));
        
        CycServiceItemDbReadType[] cycitems = new CycServiceItemDbReadType[(int) cycitemcount];
        
        ArrayList<SubItem> subitems = new ArrayList(); 
        for (S7Field userfield:s7field.getS7field()){
            subitems.add(new SubItem((short) userfield.getNumberOfElements(),
                            userfield.getBlockNumber(),
                            userfield.getByteOffset()));
        }        
        
        int initpos = 0;
        int endpos = (subitems.size() < 50)?subitems.size():49;
        int j=0;   
        int lengthInBytes = 4;        
        do{
            SubItem[] arraysubitems = subitems.subList(initpos, endpos)
                    .toArray(new SubItem[0]);
                    
            cycitems[j] = new CycServiceItemDbReadType(                                        
                    (short) (arraysubitems.length*5+2),
                    (short) 0xb0,
                    (short) arraysubitems.length,
                    arraysubitems);
            lengthInBytes += 4 + arraysubitems.length * 5;
            initpos = endpos++;
            endpos = ((initpos+49) < subitems.size())?(initpos+49):subitems.size();
            j++;
        } while (j < cycitems.length);

        S7PayloadUserDataItemCyclicServicesSubscribeRequest payload =
                new S7PayloadUserDataItemCyclicServicesSubscribeRequest(
                        DataTransportErrorCode.OK,
                        DataTransportSize.OCTET_STRING,
                        lengthInBytes,
                        cycitems.length,
                        s7field.getTimeBase(),
                        s7field.getMultiplier(),
                        cycitems
                );
        
        payloadItems.clear();
        payloadItems.add(payload);
    }
  
    

    private void encodeCycledUnSubscriptionRequest(PlcSubscriptionField  plcfield,
                                                List<S7ParameterUserDataItem> parameterItems,
                                                 List<S7PayloadUserDataItem> payloadItems){
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
        
        PlcField field = ((DefaultPlcSubscriptionField) plcfield).getPlcField(); 
        S7SubscriptionField s7field = (S7SubscriptionField) field;
        

        //TODO:Check CPU type
        S7PayloadUserDataItemCyclicServicesUnsubscribeRequest payload =
                new S7PayloadUserDataItemCyclicServicesUnsubscribeRequest (
                        DataTransportErrorCode.OK,
                        DataTransportSize.OCTET_STRING,
                        0x02,
                        (short) 0x05,
                        s7field.getAckalarms().get(0).byteValue()
                );        
        
        payloadItems.clear();
        payloadItems.add(payload);        
    }    
    

    private PlcSubscriptionResponse decodeEventSubcriptionRequest(String strField, 
                                    PlcSubscriptionRequest plcSubscriptionRequest,
                                    S7Message responseMessage
                                ) throws PlcProtocolException
    {
        //logger.info("decodeEventSubcriptionRequest: " + responseMessage);
        Map<String, ResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
        short errorClass = 0;
        short errorCode = 0;
        if(responseMessage instanceof S7MessageUserData) {
            S7MessageUserData messageUserData = (S7MessageUserData) responseMessage;
            S7PayloadUserData payload = (S7PayloadUserData) messageUserData.getPayload();
            //errorClass = payload.getItems()[0].
           // errorCode = messageUserData.getParameter().
        } else if(responseMessage instanceof S7MessageResponse) {
            S7MessageResponse messageResponse = (S7MessageResponse) responseMessage;
            errorClass = messageResponse.getErrorClass();
            errorCode = messageResponse.getErrorCode();
        } else {
            throw new PlcProtocolException("Unsupported message type " + responseMessage.getClass().getName());
        }

        // If the result contains any form of non-null error code, handle this instead.
        if((errorClass != 0) || (errorCode != 0)) {
            // This is usually the case if PUT/GET wasn't enabled on the PLC
            if((errorClass == 129) && (errorCode == 4)) {
                logger.warn("Got an error response from the PLC. This particular response code usually indicates " +
                    "that PUT/GET is not enabled on the PLC.");
                /*
                for (String fieldName : plcSubscriptionRequest.getFieldNames()) {
                    values.put(fieldName, null);
                }*/
                logger.info("Request field: " + strField + ": " + S7ParamErrorCode.valueOf(errorCode)+ " " + S7ParamErrorCode.valueOf(errorCode).getEvent());           
                values.put(strField,null);
                return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
            } else {
                logger.warn("Got an unknown error response from the PLC. Error Class: {}, Error Code {}. " +
                    "We probably need to implement explicit handling for this, so please file a bug-report " +
                    "on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump " +
                    "containing a capture of the communication.",
                    errorClass, errorCode);
                /*
                for (String fieldName : plcSubscriptionRequest.getFieldNames()) {
                    values.put(fieldName, null);
                }
                */
                logger.info("Request field: " + strField + ": " + S7ParamErrorCode.valueOf(errorCode)+ " " + S7ParamErrorCode.valueOf(errorCode).getEvent());                
                values.put(strField,null);
                return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest, values);
            }
        }
        
        S7ParameterUserData parameter = (S7ParameterUserData) responseMessage.getParameter();
        S7ParameterUserDataItem[] parameters = parameter.getItems();
        S7ParameterUserDataItemCPUFunctions itemparameter = (S7ParameterUserDataItemCPUFunctions) parameters[0];        
        errorCode = itemparameter.getErrorCode().shortValue();
        
        
        // In all other cases all went well.
        S7PayloadUserData payload = (S7PayloadUserData) responseMessage.getPayload();

        S7PayloadUserDataItem[] payloadItems = payload.getItems();

        //Only one item for any number of subscription (4)
        if (payloadItems.length == 0) {
            throw new PlcProtocolException(
                    "The number of requested items doesn't match the number of returned items");
        }

        boolean responseOk = false;
        if (payloadItems[0] instanceof S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse) {
                S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse item =
                (S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse)
                payloadItems[0];
                if ((item.getReturnCode() == DataTransportErrorCode.OK) &&
                    (item.getTransportSize() == DataTransportSize.OCTET_STRING)) {
                    responseOk = true;
                }
        } else if (payloadItems[0] instanceof S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse) {
                 S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse item =
                (S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse)
                payloadItems[0];
                if ((item.getReturnCode() == DataTransportErrorCode.OK) &&
                    (item.getTransportSize() == DataTransportSize.OCTET_STRING)) {
                    responseOk = true;
                }
        } else if (payloadItems[0] instanceof S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse) {
                 S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse item =
                (S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse)
                payloadItems[0];
                if ((item.getReturnCode() == DataTransportErrorCode.OK) &&
                    (item.getTransportSize() == DataTransportSize.OCTET_STRING)) {
                    responseOk = true;
                }
        }  else if (payloadItems[0] instanceof S7PayloadUserDataItemCpuFunctionAlarmAckResponse) {
            S7PayloadUserDataItemCpuFunctionAlarmAckResponse items =
            (S7PayloadUserDataItemCpuFunctionAlarmAckResponse) 
            payloadItems[0];
            String fieldName = (String) plcSubscriptionRequest.getFieldNames().toArray()[0];

            values.put(strField, new ResponseItem(PlcResponseCode.OK, null));            
            for (short s:items.getMessageObjects()){

                if (s == 0x0000) {
                     values.put(Integer.toHexString(s), new ResponseItem(PlcResponseCode.OK, null));
                } else if (s == 0x000a) {
                     values.put(Integer.toHexString(s), new ResponseItem(PlcResponseCode.NOT_FOUND, null));
                }
            }
           return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest,values);    
           
        }   else if (payloadItems[0] instanceof S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse) {
            
            S7PayloadUserDataItemCpuFunctionAlarmAckResponse items =
            (S7PayloadUserDataItemCpuFunctionAlarmAckResponse) 
            payloadItems[0];
            //String fieldName = (String) S7PayloadUserDataItemCyclicServicesPush .getFieldNames().toArray()[0];
            logger.warn("Request field: " + strField + ": " + S7ParamErrorCode.valueOf(errorCode)+ " " + S7ParamErrorCode.valueOf(errorCode).getEvent());          
            values.put(strField, new ResponseItem(PlcResponseCode.NOT_FOUND, null));              
           return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest,values);  
           
        }   else if (payloadItems[0] instanceof S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) {

            S7PayloadUserDataItemCpuFunctionAlarmQueryResponse items =
            (S7PayloadUserDataItemCpuFunctionAlarmQueryResponse)  payloadItems[0];
          
            ByteBuf buffer = Unpooled.directBuffer(items.getItems().length*2);
            ByteBuf rxbuffer = Unpooled.directBuffer(items.getItems().length*2);
            buffer.writeBytes(items.getItems());
            
            int numberOfItems = 1;
                   
            if (itemparameter.getLastDataUnit() == 1){

                short loop = 0xff;          
                CompletableFuture<S7MessageUserData> loopfuture = null;
                S7MessageUserData msg = null; 
                S7ParameterUserDataItemCPUFunctions loopparameter = null;
                S7PayloadUserDataItemCpuFunctionAlarmQueryResponse looppayload = null;        
                               
                do {
                    loopfuture = reassembledAlarmEvents(itemparameter.getSequenceNumber());
                    
                    try {
                      
                        msg = loopfuture.get();
                        if (msg != null) {
                            loopparameter =  (S7ParameterUserDataItemCPUFunctions) ((S7ParameterUserData) msg.getParameter()).getItems()[0];
                            looppayload = (S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) ((S7PayloadUserData) msg.getPayload()).getItems()[0];

                            buffer.writeBytes(looppayload.getItems());

                            loop = loopparameter.getLastDataUnit();

                        } else loop = 0x00;
                     
                    } catch (Exception ex) {
                        logger.warn(ex.toString());
                    }
                } while (loop > 0x00);

                rxbuffer.writeByte(looppayload .getReturnCode().getValue());
                rxbuffer.writeByte(looppayload .getTransportSize().getValue());
                rxbuffer.writeShort(looppayload.getDataLength());
                rxbuffer.writeBytes(buffer);

            } else {
               rxbuffer.writeByte(payloadItems[0].getReturnCode().getValue());
               rxbuffer.writeByte(payloadItems[0].getTransportSize().getValue());
               rxbuffer.writeShort(payloadItems[0].getDataLength());               
               rxbuffer.writeBytes(buffer);                
            }
          
            ReadBuffer readBuffer = new ReadBufferByteBased(ByteBufUtil.getBytes(rxbuffer));

            try {
                
                short cpuSubfunction = 0x00;
                
                if (s7DriverContext.getControllerType() == S7ControllerType.S7_300) {
                    cpuSubfunction = 0x13; 
                } else cpuSubfunction = 0xf0 ;

                S7PayloadUserDataItem payloaditem = 
                        S7PayloadUserDataItemIO.staticParse(readBuffer, 
                        (byte) 0x04,
                        (byte) 0x00, 
                        cpuSubfunction);

                eventqueue.add(payloaditem);
            } catch (Exception ex) {                
                logger.info(ex.toString());
            }                  
            
            PlcResponseCode rescode = (items.getReturnCode() == DataTransportErrorCode.OK) ? PlcResponseCode.OK : PlcResponseCode.INTERNAL_ERROR;
            values.put(strField, new ResponseItem(rescode, null));
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest,values); 
            
        }  else if (payloadItems[0] instanceof S7PayloadUserDataItemCyclicServicesSubscribeResponse) {
            //S7ParameterUserData parameter = (S7ParameterUserData) responseMessage.getParameter();  
            //logger.info("Aqui debe responder a Cyclic transfer");
            S7ParameterUserDataItemCPUFunctions msgparameter = (S7ParameterUserDataItemCPUFunctions)
                    parameter.getItems()[0];
            
                S7CyclicEvent cycevent = new S7CyclicEvent(plcSubscriptionRequest,
                        msgparameter.getSequenceNumber(), 
                        (S7PayloadUserDataItemCyclicServicesSubscribeResponse) payloadItems[0]); 
            eventqueue.add(cycevent);                  
            values.put(strField, new ResponseItem(PlcResponseCode.OK, cycHandle));  
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest,values);    
            
        } else if (payloadItems[0] instanceof S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse) {  
            //logger.info("Aqui debe responder a Cyclic transfer Change Driven");
            S7ParameterUserDataItemCPUFunctions msgparameter = (S7ParameterUserDataItemCPUFunctions)
                    parameter.getItems()[0];
            
                S7CyclicEvent cycevent = new S7CyclicEvent(plcSubscriptionRequest,
                        msgparameter.getSequenceNumber(), 
                        (S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse) payloadItems[0]); 
            eventqueue.add(cycevent);                  
            values.put(strField, new ResponseItem(PlcResponseCode.OK, cycHandle));  
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest,values);  
            
        }  else if (payloadItems[0] instanceof S7PayloadUserDataItemCyclicServicesErrorResponse) {
            //S7ParameterUserData parameter = (S7ParameterUserData) responseMessage.getParameter();
            //S7ParameterUserDataItem[] parameters = parameter.getItems();
            //S7ParameterUserDataItemCPUFunctions itemparameter = (S7ParameterUserDataItemCPUFunctions) parameters[0];
            //errorCode = itemparameter.getErrorCode().shortValue();
            logger.warn("Request field: " + strField + ": " + S7ParamErrorCode.valueOf(errorCode)+ " " + S7ParamErrorCode.valueOf(errorCode).getEvent());
            if (errorCode == 0x8104){
                values.put(strField, new ResponseItem(PlcResponseCode.UNSUPPORTED, null));              
            } else 
                values.put(strField, new ResponseItem(PlcResponseCode.INTERNAL_ERROR, null));
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest,values);                          
        } else if (payloadItems[0] instanceof S7PayloadUserDataItemCyclicServicesUnsubscribeResponse) {
            values.put(strField, new ResponseItem(PlcResponseCode.OK, null));  
            return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest,values);             
        }
        
        if (responseOk) {
            DefaultPlcSubscriptionField dfield = (DefaultPlcSubscriptionField) plcSubscriptionRequest.getField(strField);
            S7SubscriptionField field = (S7SubscriptionField) dfield.getPlcField();
            switch(field.getEventtype()) {
                case MODE: values.put(strField, new ResponseItem(PlcResponseCode.OK, modeHandle));
                break;
                case SYS: values.put(strField, new ResponseItem(PlcResponseCode.OK, sysHandle));
                break;
                case USR: values.put(strField, new ResponseItem(PlcResponseCode.OK, usrHandle));
                break;
                case ALM: values.put(strField, new ResponseItem(PlcResponseCode.OK, almHandle));
                break;
            }
            
           return new DefaultPlcSubscriptionResponse(plcSubscriptionRequest,values);
        }

        return null;
    }



    /**
     * This method is only called when there is no Response Handler.
     *
     */
    @Override
    protected void decode(ConversationContext<TPKTPacket> context, TPKTPacket msg) throws Exception {
        //logger.info(msg.toString());
        S7Message s7msg = msg.getPayload().getPayload();
        S7Parameter parameter = s7msg.getParameter();
        if (parameter instanceof S7ParameterModeTransition){
            eventqueue.add(parameter);
        } else
        if (parameter instanceof S7ParameterUserData) {
            S7ParameterUserData parameterud = (S7ParameterUserData) parameter;
            S7ParameterUserDataItem[] parameterudis = parameterud.getItems();
            for (S7ParameterUserDataItem parameterudi:parameterudis){
                if (parameterudi instanceof S7ParameterUserDataItemCPUFunctions) {
                    S7ParameterUserDataItemCPUFunctions myparameter = (S7ParameterUserDataItemCPUFunctions) parameterudi;
                    //TODO: Check from mspec. We can try using "instanceof"
                    if ((myparameter.getCpuFunctionGroup() == 0x04) && (myparameter.getCpuFunctionType() == 0x00) && (myparameter.getCpuSubfunction() == 0x03)) {
                        S7PayloadUserData payload = (S7PayloadUserData) s7msg.getPayload();
                        S7PayloadUserDataItem[] items = payload.getItems();
                        for (S7PayloadUserDataItem item:items){
                            if (item instanceof S7PayloadDiagnosticMessage){
                                eventqueue.add(item);
                            }
                        }
                    } else if ((myparameter.getCpuFunctionGroup() == 0x04) && 
                              (myparameter.getCpuFunctionType() == 0x00) &&
                              ((myparameter.getCpuSubfunction() == 0x05) ||
                               (myparameter.getCpuSubfunction() == 0x06) ||
                               (myparameter.getCpuSubfunction() == 0x0c) ||
                               (myparameter.getCpuSubfunction() == 0x11) ||
                               (myparameter.getCpuSubfunction() == 0x12) ||
                               (myparameter.getCpuSubfunction() == 0x13) ||
                               (myparameter.getCpuSubfunction() == 0x16))) {
                        S7PayloadUserData payload = (S7PayloadUserData) s7msg.getPayload();
                        S7PayloadUserDataItem[] items = payload.getItems();
                        for (Object item:items){
                            eventqueue.add(item);
                        }
                    } else if ((myparameter.getCpuFunctionGroup() == 0x02) && (myparameter.getCpuFunctionType() == 0x00) && (myparameter.getCpuSubfunction() == 0x01)){
                            S7ParameterUserDataItemCPUFunctions parameteritem = 
                                    (S7ParameterUserDataItemCPUFunctions)
                                    ((S7ParameterUserData) parameter).getItems()[0];
                            
                            S7PayloadUserData payload = (S7PayloadUserData) s7msg.getPayload();
                            
                            S7PayloadUserDataItemCyclicServicesPush payloaditem = 
                                    (S7PayloadUserDataItemCyclicServicesPush)
                                    payload.getItems()[0];
                            
                            S7CyclicEvent cycevent = new S7CyclicEvent(null,
                                    parameteritem.getSequenceNumber(),
                                    payloaditem);
                            eventqueue.add(cycevent);
                    
                    } else if  ((myparameter.getCpuFunctionGroup() == 0x02) && (myparameter.getCpuFunctionType() == 0x00) && (myparameter.getCpuSubfunction() == 0x05)) {   
                            S7ParameterUserDataItemCPUFunctions parameteritem = 
                                    (S7ParameterUserDataItemCPUFunctions)
                                    ((S7ParameterUserData) parameter).getItems()[0];
                            
                            S7PayloadUserData payload = (S7PayloadUserData) s7msg.getPayload();
                            
                            S7PayloadUserDataItemCyclicServicesChangeDrivenPush payloaditem = 
                                    (S7PayloadUserDataItemCyclicServicesChangeDrivenPush)
                                    payload.getItems()[0];
                            
                            S7CyclicEvent cycevent = new S7CyclicEvent(null,
                                    parameteritem.getSequenceNumber(),
                                    payloaditem);                           
                            eventqueue.add(cycevent);                        
                        
                    } else if ((myparameter.getCpuFunctionType() == 0x08) && (myparameter.getCpuSubfunction() == 0x01)) {

                    } else if ((myparameter.getCpuFunctionType() == 0x08) && (myparameter.getCpuSubfunction() == 0x04)) {

                    }
                }
            }
        }
    }

    @Override
    public void close(ConversationContext<TPKTPacket> context) {
        // TODO Implement Closing on Protocol Level
        EventLogic.stop();
    }

    private void extractControllerTypeAndFireConnected(ConversationContext<TPKTPacket> context, S7PayloadUserData payloadUserData) {
        for (S7PayloadUserDataItem item : payloadUserData.getItems()) {
            if (!(item instanceof S7PayloadUserDataItemCpuFunctionReadSzlResponse)) {
                continue;
            }
            S7PayloadUserDataItemCpuFunctionReadSzlResponse readSzlResponseItem =
                (S7PayloadUserDataItemCpuFunctionReadSzlResponse) item;
            
            //TODO: Extrac from byte array
            /*
            for (SzlDataTreeItem readSzlResponseItemItem : readSzlResponseItem.getItems()) {
                if (readSzlResponseItemItem.getItemIndex() != 0x0001) {
                    continue;
                }
                final String articleNumber = new String(readSzlResponseItemItem.getMlfb());
                s7DriverContext.setControllerType(decodeControllerType(articleNumber));

                // Send an event that connection setup is complete.
                context.fireConnected();
            }
            */
            
            ByteBuf szlitem = Unpooled.wrappedBuffer(readSzlResponseItem.getItems());
            String articleNumber = szlitem.toString(2, 20, Charset.defaultCharset());
            s7DriverContext.setControllerType(decodeControllerType(articleNumber));            
            context.fireConnected();            
        }
    }

    private TPKTPacket createIdentifyRemoteMessage() {
        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(1, new S7ParameterUserData(new S7ParameterUserDataItem[]{
                new S7ParameterUserDataItemCPUFunctions((short) 0x11, (byte) 0x4, (byte) 0x4, (short) 0x01, (short) 0x00, null, null, null)
            }), new S7PayloadUserData(new S7PayloadUserDataItem[]{
                new S7PayloadUserDataItemCpuFunctionReadSzlRequest(DataTransportErrorCode.OK, 
                        DataTransportSize.OCTET_STRING, 
                        0x04,
                        new SzlId(SzlModuleTypeClass.CPU, (byte) 0x00, 
                        SzlSublist.MODULE_IDENTIFICATION), 
                        0x0000)
            }));
        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (short) 2);
        return new TPKTPacket(cotpPacketData);
    }

    /*
    * 
    */
    private TPKTPacket createSzlReassembledRequest(int tpduId, short sequenceNumber) {
        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(tpduId, new S7ParameterUserData(new S7ParameterUserDataItem[]{
                new S7ParameterUserDataItemCPUFunctions((short) 0x12, (byte) 0x4, (byte) 0x4, (short) 0x01, sequenceNumber, (short) 0x00, (short) 0x00, 0)
            }), new S7PayloadUserData(new S7PayloadUserDataItem[]{
                new S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest(
                        DataTransportErrorCode.NOT_FOUND, 
                        DataTransportSize.NULL, 
                        0x00)
            }));
        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (short) 2);
        return new TPKTPacket(cotpPacketData);
    } 
    
    //TODO: S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest to S7PayloadUserDataItemCpuFunctionAlarmQueryNoDataRequest 
    private TPKTPacket createAlarmQueryReassembledRequest(int tpduId, short sequenceNumber) {
        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(tpduId, new S7ParameterUserData(new S7ParameterUserDataItem[]{
                new S7ParameterUserDataItemCPUFunctions((short) 0x12, (byte) 0x4, (byte) 0x4, (short) 0x13, sequenceNumber, (short) 0x00, (short) 0x00, 0)
            }), new S7PayloadUserData(new S7PayloadUserDataItem[]{
                new S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest(
                        DataTransportErrorCode.NOT_FOUND, 
                        DataTransportSize.NULL, 
                        0x00)
            }));
        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (short) 2);
        return new TPKTPacket(cotpPacketData);
    }        
    
    private TPKTPacket createS7ConnectionRequest(COTPPacketConnectionResponse cotpPacketConnectionResponse) {
        for (COTPParameter parameter : cotpPacketConnectionResponse.getParameters()) {
            if (parameter instanceof COTPParameterCalledTsap) {
                COTPParameterCalledTsap cotpParameterCalledTsap = (COTPParameterCalledTsap) parameter;
                s7DriverContext.setCalledTsapId(cotpParameterCalledTsap.getTsapId());
            } else if (parameter instanceof COTPParameterCallingTsap) {
                COTPParameterCallingTsap cotpParameterCallingTsap = (COTPParameterCallingTsap) parameter;
                if(cotpParameterCallingTsap.getTsapId() != s7DriverContext.getCallingTsapId()) {
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
        COTPPacketData cotpPacketData = new COTPPacketData(null, s7Message, true, (short) 1);
        return new TPKTPacket(cotpPacketData);
    }

    private COTPPacketConnectionRequest createCOTPConnectionRequest(int calledTsapId, int callingTsapId, COTPTpduSize cotpTpduSize) {
        return new COTPPacketConnectionRequest(
            new COTPParameter[]{
                new COTPParameterCalledTsap(calledTsapId),
                new COTPParameterCallingTsap(callingTsapId),
                new COTPParameterTpduSize(cotpTpduSize)
            }, null, (short) 0x0000, (short) 0x000F, COTPProtocolClass.CLASS_0);
    }
    
    private PlcResponse decodeReadResponse(S7Message responseMessage, PlcReadRequest plcReadRequest) throws PlcProtocolException {
        //logger.info("decodeReadResponse: " + responseMessage);
        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
        short errorClass;
        short errorCode;
        S7ParameterUserDataItemCPUFunctions parameteritem = null;
        if(responseMessage instanceof S7MessageResponseData) {
            S7MessageResponseData messageResponseData = (S7MessageResponseData) responseMessage;
            errorClass = messageResponseData.getErrorClass();
            errorCode = messageResponseData.getErrorCode();
        } else if(responseMessage instanceof S7MessageResponse) {
            S7MessageResponse messageResponse = (S7MessageResponse) responseMessage;
            errorClass = messageResponse.getErrorClass();
            errorCode = messageResponse.getErrorCode();
        } else if(responseMessage instanceof S7MessageUserData) {
            S7MessageUserData messageResponse = (S7MessageUserData) responseMessage;
            S7ParameterUserData parameters = (S7ParameterUserData) messageResponse.getParameter();
            parameteritem = (S7ParameterUserDataItemCPUFunctions) parameters.getItems()[0];
            errorClass = 0;
            errorCode = parameteritem.getErrorCode().shortValue();
        } else {
            throw new PlcProtocolException("Unsupported message type " + responseMessage.getClass().getName());
        }
        // If the result contains any form of non-null error code, handle this instead.
        if((errorClass != 0) || (errorCode != 0)) {
            // This is usually the case if PUT/GET wasn't enabled on the PLC
            if((errorClass == 129) && (errorCode == 4)) {
                logger.warn("Got an error response from the PLC. This particular response code usually indicates " +
                    "that PUT/GET is not enabled on the PLC.");
                for (String fieldName : plcReadRequest.getFieldNames()) {
                    ResponseItem<PlcValue> result = new ResponseItem<>(PlcResponseCode.ACCESS_DENIED, new PlcNull());
                    values.put(fieldName, result);
                }
                return new DefaultPlcReadResponse(plcReadRequest, values);
            } else {
                logger.warn("Got an unknown error response from the PLC. Error Class: {}, Error Code {}. " +
                    "We probably need to implement explicit handling for this, so please file a bug-report " +
                    "on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump " +
                    "containing a capture of the communication.",
                    errorClass, errorCode);
                for (String fieldName : plcReadRequest.getFieldNames()) {
                    ResponseItem<PlcValue> result = new ResponseItem<>(PlcResponseCode.INTERNAL_ERROR,  new PlcNull());
                    values.put(fieldName, result);
                }
                return new DefaultPlcReadResponse(plcReadRequest, values);
            }
        }

        //TODO: Reemsambling message.
        if (responseMessage instanceof S7MessageUserData) {
            
            S7PayloadUserData payload = (S7PayloadUserData) responseMessage.getPayload();             
            if (plcReadRequest.getNumberOfFields() != payload.getItems().length) {
                throw new PlcProtocolException(
                    "The number of requested items doesn't match the number of returned items");
            } 
            
            S7PayloadUserDataItem[] payloadItems = payload.getItems();             

            PlcResponseCode responseCode = PlcResponseCode.INTERNAL_ERROR;
            PlcValue plcValue = null;            
            int index = 0;
            for (String fieldName : plcReadRequest.getFieldNames()) {
                
                if (plcReadRequest.getField(fieldName) instanceof S7SzlField ) {

                    S7SzlField field = (S7SzlField) plcReadRequest.getField(fieldName);
                    S7PayloadUserDataItemCpuFunctionReadSzlResponse payloadItem = (S7PayloadUserDataItemCpuFunctionReadSzlResponse) payloadItems[index];               
                    responseCode = decodeResponseCode(payloadItem.getReturnCode());
                    
                    if (responseCode == PlcResponseCode.OK) {
                        try {
                            LinkedList plcvalues = null;
                            byte[] data = payloadItem.getItems();

                            plcvalues = new LinkedList<PlcSINT>();
                            for (byte b:data ) plcvalues.add(new PlcSINT(b));

                            if (parameteritem.getLastDataUnit() == 1) {

                                //TODO: reemsamble multiples mensajes
                                CompletableFuture<S7MessageUserData> future  = reassembledMessage(parameteritem.getSequenceNumber(), plcvalues);
                                S7MessageUserData msg = future.get();
                                if (msg != null) {
                                    S7ParameterUserData parameter = (S7ParameterUserData) msg.getParameter();
                                    S7PayloadUserData payload2 = (S7PayloadUserData) msg.getPayload();
                                    S7PayloadUserDataItemCpuFunctionReadSzlResponse payloadItem2 = (S7PayloadUserDataItemCpuFunctionReadSzlResponse) payload2.getItems()[0];
                                    for (byte b:payloadItem2.getItems()) plcvalues.add(new PlcSINT(b));
                                }
                               plcValue = new PlcList(plcvalues); 
                            } else {
                               plcValue = new PlcList(plcvalues);                            
                            }
                        } catch(Exception e) {
                            throw new PlcProtocolException("Error decoding PlcValue", e);
                        }

                    }
                } if (plcReadRequest.getField(fieldName) instanceof S7AckField ) {
                    S7AckField field = (S7AckField) plcReadRequest.getField(fieldName);
                    S7PayloadUserDataItemCpuFunctionAlarmAckResponse payloadItem = 
                        (S7PayloadUserDataItemCpuFunctionAlarmAckResponse) payloadItems[index];               
                    responseCode = decodeResponseCode(payloadItem.getReturnCode());  
                    short[] data = payloadItem.getMessageObjects();
                    LinkedList plcvalues = new LinkedList<PlcSINT>();
                    for (short b:data ) plcvalues.add(new PlcSINT((byte) b));  
                    plcValue = new PlcList(plcvalues);                       
                }
                
                ResponseItem<PlcValue> result = new ResponseItem<>(responseCode, plcValue);
                values.put(fieldName, result);
                index++;
            }

            return new DefaultPlcReadResponse(plcReadRequest, values);   
 
        }
        
        // In all other cases all went well.
        S7PayloadReadVarResponse payload = (S7PayloadReadVarResponse) responseMessage.getPayload();


        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcReadRequest.getNumberOfFields() != payload.getItems().length) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }


        S7VarPayloadDataItem[] payloadItems = payload.getItems();            

        int index = 0;
        for (String fieldName : plcReadRequest.getFieldNames()) {
            S7Field field = (S7Field) plcReadRequest.getField(fieldName);
           
            S7VarPayloadDataItem payloadItem = payloadItems[index];
     
            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            PlcValue plcValue = null;
            ByteBuf data = Unpooled.wrappedBuffer(payloadItem.getData());
            if (responseCode == PlcResponseCode.OK) {
                try {
                    plcValue = parsePlcValue(field, data);
                } catch(Exception e) {
                    throw new PlcProtocolException("Error decoding PlcValue", e);
                }
            }
            ResponseItem<PlcValue> result = new ResponseItem<>(responseCode, plcValue);
            values.put(fieldName, result);
            index++;
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);

    }
    
    private CompletableFuture<S7MessageUserData> reassembledMessage(short sequenceNumber, LinkedList<PlcSINT> plcValues){
        
        CompletableFuture<S7MessageUserData> future = new CompletableFuture<>();
        
        //TODO: PDU id is the same, we need check.
        int tpduId = tpduGenerator.getAndIncrement();
        // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
        if(tpduGenerator.get() == 0xFFFF) {
            tpduGenerator.set(1);
        }        
                
        TPKTPacket request = createSzlReassembledRequest(tpduId, sequenceNumber);

        context.sendRequest(request)
        .onTimeout(e -> {
            logger.warn("Timeout during Connection establishing, closing channel...");
            context.getChannel().close();
        })
        .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
        .check(p -> p.getPayload() instanceof COTPPacketData)
        .unwrap(p -> ((COTPPacketData) p.getPayload()))
        .check(p -> p.getPayload() instanceof S7MessageUserData)
        .unwrap(p -> ((S7MessageUserData) p.getPayload()))
        .check(p -> p.getPayload() instanceof S7PayloadUserData)
        .handle(messageUserData -> {
            future.complete(messageUserData);           
        });

        return future;
    }
    
    private CompletableFuture<S7MessageUserData>  reassembledAlarmEvents(short sequenceNumber) {
        CompletableFuture<S7MessageUserData> future = new CompletableFuture<>();
        
        //TODO: PDU id is the same, we need check.
        int tpduId = tpduGenerator.getAndIncrement();
        // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
        if(tpduGenerator.get() == 0xFFFF) {
            tpduGenerator.set(1);
        }        
                
        TPKTPacket request = createAlarmQueryReassembledRequest(tpduId, sequenceNumber);

        context.sendRequest(request)
        .onTimeout(e -> {
            logger.warn("Timeout during Connection establishing, closing channel...");
            context.getChannel().close();
        })
        .expectResponse(TPKTPacket.class, Duration.ofMillis(1000))
        .check(p -> p.getPayload() instanceof COTPPacketData)
        .unwrap(p -> ((COTPPacketData) p.getPayload()))
        .check(p -> p.getPayload() instanceof S7MessageUserData)
        .unwrap(p -> ((S7MessageUserData) p.getPayload()))
        .check(p -> p.getPayload() instanceof S7PayloadUserData)
        .handle(messageUserData -> {
            future.complete(messageUserData);           
        });

        return future;       
    }
        
    private PlcResponse decodeWriteResponse(S7Message responseMessage, PlcWriteRequest plcWriteRequest) throws PlcProtocolException {
        Map<String, PlcResponseCode> responses = new HashMap<>();
        short errorClass;
        short errorCode;
        if(responseMessage instanceof S7MessageResponseData) {
            S7MessageResponseData messageResponseData = (S7MessageResponseData) responseMessage;
            errorClass = messageResponseData.getErrorClass();
            errorCode = messageResponseData.getErrorCode();
        } else if(responseMessage instanceof S7MessageResponse) {
            S7MessageResponse messageResponse = (S7MessageResponse) responseMessage;
            errorClass = messageResponse.getErrorClass();
            errorCode = messageResponse.getErrorCode();
        } else {
            throw new PlcProtocolException("Unsupported message type " + responseMessage.getClass().getName());
        }
        // If the result contains any form of non-null error code, handle this instead.
        if((errorClass != 0) || (errorCode != 0)) {
            // This is usually the case if PUT/GET wasn't enabled on the PLC
            if((errorClass == 129) && (errorCode == 4)) {
                logger.warn("Got an error response from the PLC. This particular response code usually indicates " +
                    "that PUT/GET is not enabled on the PLC.");
                for (String fieldName : plcWriteRequest.getFieldNames()) {
                    responses.put(fieldName, PlcResponseCode.ACCESS_DENIED);
                }
                return new DefaultPlcWriteResponse(plcWriteRequest, responses);
            } else {
                logger.warn("Got an unknown error response from the PLC. Error Class: {}, Error Code {}. " +
                        "We probably need to implement explicit handling for this, so please file a bug-report " +
                        "on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump " +
                        "containing a capture of the communication.",
                    errorClass, errorCode);
                for (String fieldName : plcWriteRequest.getFieldNames()) {
                    responses.put(fieldName, PlcResponseCode.INTERNAL_ERROR);
                }
                return new DefaultPlcWriteResponse(plcWriteRequest, responses);
            }
        }
        
        // In all other cases all went well.
        S7PayloadWriteVarResponse payload = (S7PayloadWriteVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcWriteRequest.getNumberOfFields() != payload.getItems().length) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        S7VarPayloadStatusItem[] payloadItems = payload.getItems();
        int index = 0;
        for (String fieldName : plcWriteRequest.getFieldNames()) {
            S7VarPayloadStatusItem payloadItem = payloadItems[index];

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            responses.put(fieldName, responseCode);
            index++;
        }

        return new DefaultPlcWriteResponse(plcWriteRequest, responses);
    }

    private S7VarPayloadDataItem serializePlcValue(S7Field field, PlcValue plcValue) {
        try {
            DataTransportSize transportSize = field.getDataType().getDataTransportSize();
            int stringLength = (field instanceof S7StringField) ? ((S7StringField) field).getStringLength() : 254;
            ByteBuffer byteBuffer = null;
            for(int i = 0; i < field.getNumberOfElements(); i++) {
                WriteBufferByteBased writeBuffer = DataItemIO.staticSerialize(plcValue.getIndex(i),
                    field.getDataType().getDataProtocolId(), stringLength);
                if(writeBuffer != null) {
                    // Allocate enough space for all items.
                    if(byteBuffer == null) {
                        byteBuffer = ByteBuffer.allocate(writeBuffer.getData().length * field.getNumberOfElements());
                    }
                    byteBuffer.put(writeBuffer.getData());
                }
            }
            if(byteBuffer != null) {
                byte[] data = byteBuffer.array();
                return new S7VarPayloadDataItem(DataTransportErrorCode.OK, transportSize, data);
            }
        } catch (ParseException e) {
            logger.warn("Error serializing field item of type: '{}'", field.getDataType().name(), e);
        }
        return null;
    }

    private PlcValue parsePlcValue(S7Field field, ByteBuf data) {
        ReadBuffer readBuffer = new ReadBufferByteBased(data.array());
        try {
            int stringLength = (field instanceof S7StringField) ? ((S7StringField) field).getStringLength() : 254;
            if (field.getNumberOfElements() == 1) {
                return DataItemIO.staticParse(readBuffer, field.getDataType().getDataProtocolId(),
                    stringLength);
            } else {
                // Fetch all
                final PlcValue[] resultItems = IntStream.range(0, field.getNumberOfElements()).mapToObj(i -> {
                    try {
                        return DataItemIO.staticParse(readBuffer, field.getDataType().getDataProtocolId(),
                            stringLength);
                    } catch (ParseException e) {
                        logger.warn("Error parsing field item of type: '{}' (at position {}})", field.getDataType().name(), i, e);
                    }
                    return null;
                }).toArray(PlcValue[]::new);
                return IEC61131ValueHandler.of(resultItems);
            }
        } catch (ParseException e) {
            logger.warn("Error parsing field item of type: '{}'", field.getDataType().name(), e);
        }
        return null;
    }

    /**
     * Helper to convert the return codes returned from the S7 into one of our standard
     * PLC4X return codes
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
                return PlcResponseCode.NOT_FOUND;
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
    private S7ControllerType decodeControllerType(String articleNumber) {
        if (!articleNumber.startsWith("6ES7 ")) {
            return S7ControllerType.ANY;
        }
        String model = articleNumber.substring(articleNumber.indexOf(' ') + 1, articleNumber.indexOf(' ') + 2);
        switch (model) {
            case "2":
                return S7ControllerType.S7_1200;
            case "5":
                return S7ControllerType.S7_1500;
            case "3":
                return S7ControllerType.S7_300;
            case "4":
                return S7ControllerType.S7_400;
            default:
                if (logger.isInfoEnabled()) {
                    logger.info("Looking up unknown article number {}", articleNumber);
                }
                return S7ControllerType.ANY;
        }
    }

    /**
     * Currently we only support the S7 Any type of addresses. This helper simply converts the S7Field
     * from PLC4X into S7Address objects.
     * @param field S7Field instance we need to convert into an S7Address
     * @return the S7Address
     */
    protected S7Address encodeS7Address(PlcField field) {
        if (!(field instanceof S7Field)) {
            throw new PlcRuntimeException("Unsupported address type " + field.getClass().getName());
        }
        S7Field s7Field = (S7Field) field;
        TransportSize transportSize = s7Field.getDataType();
        int numElements = s7Field.getNumberOfElements();
        // For these date-types we have to convert the requests to simple byte-array requests
        // As otherwise the S7 will deny them with "Data type not supported" replies.
        if((transportSize == TransportSize.TIME) /*|| (transportSize == TransportSize.S7_S5TIME)*/ ||
            (transportSize == TransportSize.LTIME) || (transportSize == TransportSize.DATE) ||
            (transportSize == TransportSize.TIME_OF_DAY) || (transportSize == TransportSize.DATE_AND_TIME)) {
            numElements = numElements * transportSize.getSizeInBytes();
            transportSize = TransportSize.BYTE;
        }
        if(transportSize == TransportSize.STRING) {
            transportSize = TransportSize.CHAR;
            int stringLength = (s7Field instanceof S7StringField) ? ((S7StringField) s7Field).getStringLength() : 254;
            numElements = numElements * (stringLength + 2);
        } else if(transportSize == TransportSize.WSTRING) {
            transportSize = TransportSize.CHAR;
            int stringLength = (s7Field instanceof S7StringField) ? ((S7StringField) s7Field).getStringLength() : 254;
            numElements = numElements * (stringLength + 2) * 2;
        }
        return new S7AddressAny(transportSize, numElements, s7Field.getBlockNumber(),
            s7Field.getMemoryArea(), s7Field.getByteOffset(), s7Field.getBitOffset());
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
            transaction.endRequest();
            future.completeExceptionally(e);
        }

        @Override
        public void accept(TPKTPacket tpktPacket, E e) {
            transaction.endRequest();
            future.completeExceptionally(e);
        }
    }

}
