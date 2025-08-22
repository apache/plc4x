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
package org.apache.plc4x.java.s7light.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.tag.*;
import org.apache.plc4x.java.s7light.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.s7light.readwrite.tag.S7PlcTagHandler;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.*;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

/**
 * The S7 Protocol states that there can not be more then {min(maxAmqCaller, maxAmqCallee} "ongoing" requests.
 * So we need to limit those.
 * Thus, each request goes to a Work Queue and this Queue ensures, that only 3 are open at the same time.
 */
public class S7ProtocolLogic extends Plc4xProtocolBase<TPKTPacket> {

    private static final Logger logger = LoggerFactory.getLogger(S7ProtocolLogic.class);

    private final AtomicInteger tpduGenerator = new AtomicInteger(10);
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
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new S7PlcTagHandler();
    }

    @Override
    public void close(ConversationContext<TPKTPacket> context) {
        tm.shutdown();
    }

    @Override
    public void onConnect(ConversationContext<TPKTPacket> context) {
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
            .expectResponse(TPKTPacket.class, s7DriverContext.getReadTimeoutDuration())
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
                    .expectResponse(TPKTPacket.class, s7DriverContext.getReadTimeoutDuration())
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
                            .expectResponse(TPKTPacket.class, s7DriverContext.getReadTimeoutDuration())
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

    @Override
    public void onDisconnect(ConversationContext<TPKTPacket> context) {
        // 2. Performs the shutdown of the transaction executor.
        tm.shutdown();
        // 4. Executes the closing of the main channel.
        context.getChannel().close();
   }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        CompletableFuture<S7Message> responseFuture;
        // If the request contains at least one var-length string field, we need to get the real length first.
        if (request.getTagNames().stream().anyMatch(t -> request.getTag(t) instanceof S7StringVarLengthTag)) {
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
        //Pointers
        S7Message[] responseMessage = new S7Message[1];
        PlcReadRequest[] plcReadRequest = new PlcReadRequest[1];
        
        responseFuture.whenComplete((s7Message, throwable) -> {
            if (throwable != null) {
                clientFuture.completeExceptionally(new PlcProtocolException("Error reading", throwable));
            } else {
                try {
                    responseMessage[0] = s7Message;
                    plcReadRequest[0] = readRequest;
                    PlcReadResponse response = (PlcReadResponse) decodeReadResponse(responseMessage[0], plcReadRequest[0]);
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
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        CompletableFuture<S7Message> responseFuture = new CompletableFuture<>();

        // If the list of tags contains at least one STRING/WSTRING element,
        // we need to check the sizes of the string fields in a first request.
        if (request.getTagNames().stream().anyMatch(t -> request.getTag(t) instanceof S7StringVarLengthTag)) {
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
            .expectResponse(TPKTPacket.class, s7DriverContext.getReadTimeoutDuration())
            .unwrap(TPKTPacket::getPayload)
            .only(COTPPacketData.class)
            .check(p -> p.getPayload() != null)
            .unwrap(COTPPacket::getPayload)
            .check(p -> p.getTpduReference() == tpduId)
            .handle(p -> {
                // Finish the request-transaction.
                transaction.endRequest();

                try {
                    future.complete(p);
                } catch (Exception e) {
                    logger.warn("Error sending 'write' message: '{}'", e.getMessage(), e);
                }
            }));

        return future;
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

        S7ParameterUserDataItemCPUFunctions parameterItem;
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
            parameterItem = (S7ParameterUserDataItemCPUFunctions) parameters.getItems().get(0);
            errorClass = 0;
            errorCode = parameterItem.getErrorCode().shortValue();
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

        // In all other cases all went well.
        S7PayloadReadVarResponse payload = (S7PayloadReadVarResponse) responseMessage.getPayload();

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as the PLC does not return this information.
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
            } else if((tag.getDataType() == TransportSize.BOOL) && (tag.getNumberOfElements() > 1)) {
                if(!(plcValue instanceof PlcList)) {
                    throw new PlcRuntimeException(String.format("Expected a PlcList with %d PlcBOOL elements", tag.getNumberOfElements()));
                }
                PlcList plcList = (PlcList) plcValue;
                int numBytes = (tag.getNumberOfElements() + 7) / 8;
                byteBuffer = ByteBuffer.allocate(numBytes);
                for (int i = 0; i < tag.getNumberOfElements(); i++) {
                    if(!(plcList.getIndex(i) instanceof PlcBOOL)) {
                        throw new PlcRuntimeException(String.format("Expected a PlcList with %d PlcBOOL elements", tag.getNumberOfElements()));
                    }
                    PlcBOOL plcBOOL = (PlcBOOL) plcList.getIndex(i);
                    if(plcBOOL.getBoolean()) {
                        int curByte = i / 8;
                        int curBit = i % 8;
                        byteBuffer.put(curByte, (byte) (1 << curBit | byteBuffer.get(curByte)));
                    }
                }
                transportSize = DataTransportSize.BYTE_WORD_DWORD;
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
     * Currently, we only support the S7 Any type of addresses. This helper simply converts the S7Tag
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
        } else if ((transportSize == TransportSize.BOOL) && (s7Tag.getNumberOfElements() > 1)) {
            numElements = (s7Tag.getNumberOfElements() + 7) / 8;
            transportSize = TransportSize.BYTE;
        }
        if (transportSize.getCode() == 0x00) {
            numElements = numElements * transportSize.getSizeInBytes();
            transportSize = TransportSize.BYTE;
        }
        return new S7AddressAny(transportSize, numElements, s7Tag.getBlockNumber(),
            s7Tag.getMemoryArea(), s7Tag.getByteOffset(), s7Tag.getBitOffset());
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
