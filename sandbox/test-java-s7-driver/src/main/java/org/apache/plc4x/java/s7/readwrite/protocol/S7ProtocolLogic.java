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

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.vavr.control.Either;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.COTPPacket;
import org.apache.plc4x.java.s7.readwrite.COTPPacketConnectionRequest;
import org.apache.plc4x.java.s7.readwrite.COTPPacketConnectionResponse;
import org.apache.plc4x.java.s7.readwrite.COTPPacketData;
import org.apache.plc4x.java.s7.readwrite.COTPParameter;
import org.apache.plc4x.java.s7.readwrite.COTPParameterCalledTsap;
import org.apache.plc4x.java.s7.readwrite.COTPParameterCallingTsap;
import org.apache.plc4x.java.s7.readwrite.COTPParameterTpduSize;
import org.apache.plc4x.java.s7.readwrite.S7Address;
import org.apache.plc4x.java.s7.readwrite.S7AddressAny;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageRequest;
import org.apache.plc4x.java.s7.readwrite.S7MessageResponse;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterReadVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7ParameterReadVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7ParameterSetupCommunication;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7ParameterWriteVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7ParameterWriteVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadReadVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadReadVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadSetupCommunication;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7VarPayloadDataItem;
import org.apache.plc4x.java.s7.readwrite.S7VarPayloadStatusItem;
import org.apache.plc4x.java.s7.readwrite.S7VarRequestParameterItem;
import org.apache.plc4x.java.s7.readwrite.S7VarRequestParameterItemAddress;
import org.apache.plc4x.java.s7.readwrite.SzlDataTreeItem;
import org.apache.plc4x.java.s7.readwrite.SzlId;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.s7.readwrite.configuration.S7Configuration;
import org.apache.plc4x.java.s7.readwrite.io.DataItemIO;
import org.apache.plc4x.java.s7.readwrite.optimizer.S7MessageProcessor;
import org.apache.plc4x.java.s7.readwrite.types.COTPProtocolClass;
import org.apache.plc4x.java.s7.readwrite.types.COTPTpduSize;
import org.apache.plc4x.java.s7.readwrite.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.types.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.types.DeviceGroup;
import org.apache.plc4x.java.s7.readwrite.types.S7ControllerType;
import org.apache.plc4x.java.s7.readwrite.types.SzlModuleTypeClass;
import org.apache.plc4x.java.s7.readwrite.types.SzlSublist;
import org.apache.plc4x.java.s7.readwrite.field.S7Field;
import org.apache.plc4x.java.s7.readwrite.utils.S7TsapIdEncoder;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.spi.messages.InternalPlcWriteRequest;
import org.apache.plc4x.java.spi.optimizer.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

/**
 * The S7 Protocol states that there can not be more then {min(maxAmqCaller, maxAmqCallee} "ongoing" requests.
 * So we need to limit those.
 * Thus, each request goes to a Work Queue and this Queue ensures, that only 3 are open at the same time.
 */
public class S7ProtocolLogic extends Plc4xProtocolBase<TPKTPacket> implements HasConfiguration<S7Configuration> {

    private static final Logger logger = LoggerFactory.getLogger(S7ProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private int callingTsapId;
    private int calledTsapId;
    private COTPTpduSize cotpTpduSize;
    private int pduSize;
    private int maxAmqCaller;
    private int maxAmqCallee;
    private S7ControllerType controllerType;

    private static final AtomicInteger tpduGenerator = new AtomicInteger(10);
    private RequestTransactionManager tm;

    private S7MessageProcessor processor = null;

    @Override
    public void setConfiguration(S7Configuration configuration) {
        this.callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OTHERS,
            configuration.localRack, configuration.localSlot);
        this.calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC,
            configuration.remoteRack, configuration.remoteSlot);

        this.controllerType = configuration.controllerType == null ? S7ControllerType.ANY : S7ControllerType.valueOf(configuration.controllerType);
        // The Siemens LOGO device seems to only work with very limited settings,
        // so we're overriding some of the defaults.
        if (this.controllerType == S7ControllerType.LOGO && configuration.pduSize == 1024) {
            configuration.pduSize = 480;
        }

        // Initialize the parameters with initial version (Will be updated during the login process)
        this.cotpTpduSize = getNearestMatchingTpduSize((short) configuration.getPduSize());
        // The PDU size is theoretically not bound by the COTP TPDU size, however having a larger
        // PDU size would make the code extremely complex. But even if the protocol would allow this
        // I have never seen this happen in reality. Making is smaller would unnecessarily limit the
        // size, so we're setting it to the maximum that can be included.
        this.pduSize = cotpTpduSize.getSizeInBytes() - 16;
        this.maxAmqCaller = configuration.maxAmqCaller;
        this.maxAmqCallee = configuration.maxAmqCallee;

        // Initialize Transaction Manager.
        // Until the number of concurrent requests is successfully negotiated we set it to a
        // maximum of only one request being able to be sent at a time. During the login process
        // No concurrent requests can be sent anyway. It will be updated when receiving the
        // S7ParameterSetupCommunication response.
        this.tm = new RequestTransactionManager(1);

        // REMARK: Perhaps make this configurable.
        // TODO: Comment in to enable ...
        //this.processor = new DefaultS7MessageProcessor(tpduGenerator);
    }

    @Override
    public void onConnect(ConversationContext<TPKTPacket> context) {
        logger.debug("Sending COTP Connection Request");
        // Open the session on ISO Transport Protocol first.
        TPKTPacket packet = new TPKTPacket(createCOTPConnectionRequest(calledTsapId, callingTsapId, cotpTpduSize));

        context.sendRequest(packet)
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .check(p -> p.getPayload() instanceof COTPPacketConnectionResponse)
            .unwrap(p -> (COTPPacketConnectionResponse) p.getPayload())
            .handle(cotpPacketConnectionResponse -> {
                logger.debug("Got COTP Connection Response");
                logger.debug("Sending S7 Connection Request");
                context.sendRequest(createS7ConnectionRequest(cotpPacketConnectionResponse))
                    .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
                    .unwrap(TPKTPacket::getPayload)
                    .only(COTPPacketData.class)
                    .unwrap(COTPPacket::getPayload)
                    .only(S7MessageResponse.class)
                    .unwrap(S7Message::getParameter)
                    .only(S7ParameterSetupCommunication.class)
                    .handle(setupCommunication -> {
                        logger.debug("Got S7 Connection Response");
                        // Save some data from the response.
                        maxAmqCaller = setupCommunication.getMaxAmqCaller();
                        maxAmqCallee = setupCommunication.getMaxAmqCallee();
                        pduSize = setupCommunication.getPduLength();

                        // Update the number of concurrent requests to the negotiated number.
                        // I have never seen anything else than equal values for caller and
                        // callee, but if they were different, we're only limiting the outgoing
                        // requests.
                        tm.setNumberOfConcurrentRequests(maxAmqCaller);

                        // If the controller type is explicitly set, were finished with the login
                        // process. If it's set to ANY, we have to query the serial number information
                        // in order to detect the type of PLC.
                        if (controllerType != S7ControllerType.ANY) {
                            // Send an event that connection setup is complete.
                            context.fireConnected();
                            return;
                        }

                        // Prepare a message to request the remote to identify itself.
                        logger.debug("Sending S7 Identification Request");
                        TPKTPacket tpktPacket = createIdentifyRemoteMessage();
                        context.sendRequest(tpktPacket)
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

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<S7VarRequestParameterItem> requestItems = new ArrayList<>(request.getNumberOfFields());
        for (PlcField field : request.getFields()) {
            requestItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(field)));
        }

        // Create a read request template.
        // tpuId will be inserted before sending in #readInternal so we insert -1 as dummy here
        final S7MessageRequest s7MessageRequest = new S7MessageRequest(-1,
            new S7ParameterReadVarRequest(requestItems.toArray(new S7VarRequestParameterItem[0])),
            new S7PayloadReadVarRequest());

        // If no processor is provided the request is output as-is without any modification.
        if (processor == null) {
            // Just send a single response and chain it as Response
            return toPlcReadResponse((InternalPlcReadRequest) readRequest, readInternal(s7MessageRequest));
        }

        try {
            // Do the preprocessing and eventually split up into multiple requests.
            final Collection<S7MessageRequest> s7MessageRequests = processor.processRequest(s7MessageRequest, pduSize);

            // Only if more than one sub-request is returned, do something special ...
            // otherwise just do the normal sending.
            if (s7MessageRequests.size() == 1) {
                return toPlcReadResponse((InternalPlcReadRequest) readRequest, readInternal(s7MessageRequest));
            }

            /////////////////////////////////////////////////////////////////
            //  Here we are in the case that we have multiple requests,
            //  so we need splitting
            /////////////////////////////////////////////////////////////////
            ParentFuture multiRequestFuture = new ParentFuture(
                result -> {
                    try {
                        final S7MessageResponse s7MessageResponse =
                            processor.processResponse(s7MessageRequest, result);
                        final PlcReadResponse plcReadResponse = (PlcReadResponse)
                            decodeReadResponse(s7MessageResponse, ((InternalPlcReadRequest) readRequest));
                        future.complete(plcReadResponse);
                    } catch (PlcException e) {
                        logger.error("Something went wrong", e);
                    }
                });
            for (S7MessageRequest messageRequest : s7MessageRequests) {
                ChildFuture childFuture = new ChildFuture(messageRequest);
                multiRequestFuture.addChildFuture(childFuture);

                readInternal(messageRequest)
                    // Forward everything to the remaining future
                    .handle((res, ex) -> res != null ? childFuture.complete(res) : childFuture.completeExceptionally(ex));
            }
            return multiRequestFuture;
        } catch (PlcException e) {
            logger.error("Something went wrong", e);
        }
        return null;
    }

    /** Maps the S7ReadResponse of a PlcReadRequest to a PlcReadRespoonse */
    private CompletableFuture<PlcReadResponse> toPlcReadResponse(InternalPlcReadRequest readRequest, CompletableFuture<S7MessageResponse> response) {
        return response
            .thenApply(p -> {
                try {
                    return ((PlcReadResponse) decodeReadResponse(p, readRequest));
                } catch (PlcProtocolException e) {
                    throw new PlcRuntimeException("Unable to decode Response", e);
                }
            });
    }

    /**
     * Sends one Read over the Wire and internally returns the Response
     * Do sending of normally sized single-message request.
     *
     * Assumes that the {@link S7MessageRequest} and its expected {@link S7MessageResponse}
     * and does not further check that!
     */
    private CompletableFuture<S7MessageResponse> readInternal(S7MessageRequest request) {
        CompletableFuture<S7MessageResponse> future = new CompletableFuture<>();
        int tpduId = tpduGenerator.getAndIncrement();

        // Create a new Request with correct tpuId (is not known before)
        S7MessageRequest s7MessageRequest = new S7MessageRequest(tpduId, request.getParameter(), request.getPayload());

        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null, s7MessageRequest, true, (short) tpduId));
        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(tpktPacket)
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p.getPayload() instanceof COTPPacketData)
            .unwrap(p -> (COTPPacketData) p.getPayload())
            .check(p -> p.getPayload() instanceof S7MessageResponse)
            .unwrap(p -> (S7MessageResponse) p.getPayload())
            .check(p -> p.getTpduReference() == tpduId)
            .check(p -> p.getParameter() instanceof S7ParameterReadVarResponse)
            .handle(p -> {
                future.complete(p);
                // Finish the request-transaction.
                transaction.endRequest();
//                    future.complete(((PlcReadResponse) decodeReadResponse(p, ((InternalPlcReadRequest) readRequest))));
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
        TPKTPacket tpktPacket = new TPKTPacket(new COTPPacketData(null,
            new S7MessageRequest(tpduId,
                new S7ParameterWriteVarRequest(parameterItems.toArray(new S7VarRequestParameterItem[0])),
                new S7PayloadWriteVarRequest(payloadItems.toArray(new S7VarPayloadDataItem[0]))),
            true, (short) tpduId));

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(tpktPacket)
            .expectResponse(TPKTPacket.class, REQUEST_TIMEOUT)
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p.getPayload() instanceof COTPPacketData)
            .unwrap(p -> ((COTPPacketData) p.getPayload()))
            .check(p -> p.getPayload() instanceof S7MessageResponse)
            .unwrap(p -> ((S7MessageResponse) p.getPayload()))
            .check(p -> p.getTpduReference() == tpduId)
            .check(p -> p.getParameter() instanceof S7ParameterWriteVarResponse)
            .handle(p -> {
                try {
                    future.complete(((PlcWriteResponse) decodeWriteResponse(p, ((InternalPlcWriteRequest) writeRequest))));
                } catch (PlcProtocolException e) {
                    logger.warn(String.format("Error sending 'write' message: '%s'", e.getMessage()), e);
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    @Override
    public void close(ConversationContext<TPKTPacket> context) {
        // TODO Implement Closing on Protocol Level
    }

    private void extractControllerTypeAndFireConnected(ConversationContext<TPKTPacket> context, S7PayloadUserData payloadUserData) {
        for (S7PayloadUserDataItem item : payloadUserData.getItems()) {
            if (!(item instanceof S7PayloadUserDataItemCpuFunctionReadSzlResponse)) {
                continue;
            }
            S7PayloadUserDataItemCpuFunctionReadSzlResponse readSzlResponseItem =
                (S7PayloadUserDataItemCpuFunctionReadSzlResponse) item;
            for (SzlDataTreeItem readSzlResponseItemItem : readSzlResponseItem.getItems()) {
                if (readSzlResponseItemItem.getItemIndex() != 0x0001) {
                    continue;
                }
                final String articleNumber = new String(readSzlResponseItemItem.getMlfb());
                controllerType = decodeControllerType(articleNumber);

                // Send an event that connection setup is complete.
                context.fireConnected();
            }
        }
    }

    private TPKTPacket createIdentifyRemoteMessage() {
        S7MessageUserData identifyRemoteMessage = new S7MessageUserData(1, new S7ParameterUserData(new S7ParameterUserDataItem[]{
            new S7ParameterUserDataItemCPUFunctions((short) 0x11, (byte) 0x4, (byte) 0x4, (short) 0x01, (short) 0x00, null, null, null)
        }), new S7PayloadUserData(new S7PayloadUserDataItem[]{
            new S7PayloadUserDataItemCpuFunctionReadSzlRequest(DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, new SzlId(SzlModuleTypeClass.CPU, (byte) 0x00, SzlSublist.MODULE_IDENTIFICATION), 0x0000)
        }));
        COTPPacketData cotpPacketData = new COTPPacketData(null, identifyRemoteMessage, true, (short) 2);
        return new TPKTPacket(cotpPacketData);
    }

    private TPKTPacket createS7ConnectionRequest(COTPPacketConnectionResponse cotpPacketConnectionResponse) {
        for (COTPParameter parameter : cotpPacketConnectionResponse.getParameters()) {
            if (parameter instanceof COTPParameterCalledTsap) {
                COTPParameterCalledTsap cotpParameterCalledTsap = (COTPParameterCalledTsap) parameter;
                calledTsapId = cotpParameterCalledTsap.getTsapId();
            } else if (parameter instanceof COTPParameterCallingTsap) {
                COTPParameterCallingTsap cotpParameterCallingTsap = (COTPParameterCallingTsap) parameter;
                if(cotpParameterCallingTsap.getTsapId() != callingTsapId) {
                    callingTsapId = cotpParameterCallingTsap.getTsapId();
                    logger.warn(String.format("Switching calling TSAP id to '%s'", callingTsapId));
                }
            } else if (parameter instanceof COTPParameterTpduSize) {
                COTPParameterTpduSize cotpParameterTpduSize = (COTPParameterTpduSize) parameter;
                cotpTpduSize = cotpParameterTpduSize.getTpduSize();
            } else {
                logger.warn(String.format("Got unknown parameter type '%s'", parameter.getClass().getName()));
            }
        }

        // Send an S7 login message.
        S7ParameterSetupCommunication s7ParameterSetupCommunication =
            new S7ParameterSetupCommunication(maxAmqCaller, maxAmqCallee, pduSize);
        S7Message s7Message = new S7MessageRequest(0, s7ParameterSetupCommunication,
            new S7PayloadSetupCommunication());
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

    private PlcResponse decodeReadResponse(S7MessageResponse responseMessage, InternalPlcReadRequest plcReadRequest) throws PlcProtocolException {
        S7PayloadReadVarResponse payload = (S7PayloadReadVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcReadRequest.getNumberOfFields() != payload.getItems().length) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        Map<String, Pair<PlcResponseCode, PlcValue>> values = new HashMap<>();
        S7VarPayloadDataItem[] payloadItems = payload.getItems();
        int index = 0;
        for (String fieldName : plcReadRequest.getFieldNames()) {
            S7Field field = (S7Field) plcReadRequest.getField(fieldName);
            S7VarPayloadDataItem payloadItem = payloadItems[index];

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            PlcValue plcValue = null;
            ByteBuf data = Unpooled.wrappedBuffer(payloadItem.getData());
            if (responseCode == PlcResponseCode.OK) {
                plcValue = parsePlcValue(field, data);
            }
            Pair<PlcResponseCode, PlcValue> result = new ImmutablePair<>(responseCode, plcValue);
            values.put(fieldName, result);
            index++;
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);
    }

    private PlcResponse decodeWriteResponse(S7MessageResponse responseMessage, InternalPlcWriteRequest plcWriteRequest) throws PlcProtocolException {
        S7PayloadWriteVarResponse payload = (S7PayloadWriteVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcWriteRequest.getNumberOfFields() != payload.getItems().length) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        Map<String, PlcResponseCode> responses = new HashMap<>();
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
            DataTransportSize transportSize = (field.getDataType().getDataProtocolId() == 1) ?
                DataTransportSize.BIT : DataTransportSize.BYTE_WORD_DWORD;
            WriteBuffer writeBuffer = DataItemIO.serialize(plcValue, field.getDataType().getDataProtocolId());
            if(writeBuffer != null) {
                byte[] data = writeBuffer.getData();
                return new S7VarPayloadDataItem(DataTransportErrorCode.OK, transportSize, data.length, data);
            }
        } catch (ParseException e) {
            logger.warn(String.format("Error serializing field item of type: '%s'", field.getDataType().name()), e);
        }
        return null;
    }

    private PlcValue parsePlcValue(S7Field field, ByteBuf data) {
        ReadBuffer readBuffer = new ReadBuffer(data.array());
        try {
            return DataItemIO.parse(readBuffer, field.getDataType().getDataProtocolId());
        } catch (ParseException e) {
            logger.warn(String.format("Error parsing field item of type: '%s'", field.getDataType().name()), e);
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
                    logger.info(String.format("Looking up unknown article number %s", articleNumber));
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
            throw new RuntimeException("Unsupported address type " + field.getClass().getName());
        }
        S7Field s7Field = (S7Field) field;
        return new S7AddressAny(s7Field.getDataType(), s7Field.getNumElements(), s7Field.getBlockNumber(),
            s7Field.getMemoryArea(), s7Field.getByteOffset(), s7Field.getBitOffset());
    }

    /**
     * Iterate over all values until one is found that the given tpdu size will fit.
     *
     * @param tpduSizeParameter requested tpdu size.
     * @return smallest {@link COTPTpduSize} which will fit a given size of tpdu.
     */
    protected COTPTpduSize getNearestMatchingTpduSize(short tpduSizeParameter) {
        for (COTPTpduSize value : COTPTpduSize.values()) {
            if (value.getSizeInBytes() >= tpduSizeParameter) {
                return value;
            }
        }
        return null;
    }

    private static class ChildFuture extends CompletableFuture<S7MessageResponse> {

        private final S7MessageRequest request;

        private Throwable exception;

        public ChildFuture(S7MessageRequest request) {
            this.request = request;
        }

        @Override
        public boolean completeExceptionally(Throwable exception) {
            this.exception = exception;
            return super.completeExceptionally(exception);
        }

        public S7MessageRequest getRequest() {
            return request;
        }

        public Throwable getException() {
            return exception;
        }

    }

    private static class ParentFuture extends CompletableFuture<PlcReadResponse> {

        private final Consumer<Map<S7MessageRequest, Either<S7MessageResponse, Throwable>>> action;
        private final List<ChildFuture> childFutures;

        public ParentFuture(Consumer<Map<S7MessageRequest, Either<S7MessageResponse, Throwable>>> action) {
            this.action = action;
            this.childFutures = new LinkedList<>();
        }

        public void addChildFuture(ChildFuture future) {
            future.whenComplete((value, throwable) -> {
                // If all futures are done (None are not done), do the completion action.
                if(childFutures.stream().allMatch(CompletableFuture::isDone)) {
                    Map<S7MessageRequest, Either<S7MessageResponse, Throwable>> result = new HashMap<>();
                    childFutures.forEach(childFuture -> {
                        Either<S7MessageResponse, Throwable> subResult = null;
                        if (childFuture.isCancelled()) {
                            subResult = Either.right(new Exception("Cancelled"));
                        } else if (childFuture.isCompletedExceptionally()) {
                            // Get the original exception and pass it along.
                            subResult = Either.right(childFuture.getException());
                        } else {
                            try {
                                subResult = Either.left(childFuture.get());
                            } catch (InterruptedException | ExecutionException e) {
                                e.printStackTrace();
                            }
                        }
                        result.put(childFuture.getRequest(), subResult);
                    });

                    // Call the result handler and pass in the list of results.
                    action.accept(result);
                }
            });
            childFutures.add(future);
        }

    }

}
