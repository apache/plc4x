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
package org.apache.plc4x.java.abeth.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.abeth.field.AbEthField;
import org.apache.plc4x.java.abeth.readwrite.*;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.PlcMessageToMessageCodec;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.PlcRequestContainer;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;

@Deprecated
public class Plc4xAbEthProtocol extends PlcMessageToMessageCodec<CIPEncapsulationPacket, PlcRequestContainer> {

    private static final Logger logger = LoggerFactory.getLogger(Plc4xAbEthProtocol.class);

    private static final AtomicInteger transactionCounterGenerator = new AtomicInteger(10);
    private static final List<Short> emptySenderContext = Arrays.asList((short) 0x00 ,(short) 0x00 ,(short) 0x00,
            (short) 0x00,(short) 0x00,(short) 0x00, (short) 0x00,(short) 0x00);

    private long sessionHandle;
    private Map<Integer, PlcRequestContainer> requests;
    private int station;

    public Plc4xAbEthProtocol(int station) {
        logger.trace("Created new instance of PLC4X-AB-ETH Protocol");
        this.requests = new HashMap<>();
        this.station = station;
    }

    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        logger.trace("Registered user event {}", evt);
        // If the connection has just been established, start setting up the connection
        // by sending a connection request to the plc.
        if (evt instanceof ConnectEvent) {
            logger.debug("AB-ETH Sending Connection Request");
            // Open the session on ISO Transport Protocol first.
            CIPEncapsulationConnectionRequest connectionRequest = new CIPEncapsulationConnectionRequest(0L, 0L,
                emptySenderContext, 0L);
            ctx.channel().writeAndFlush(connectionRequest);
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
        logger.trace("Encoding {}", msg);
        PlcRequest request = msg.getRequest();

        // reset counter since two byte values are possible in DF1
        if (transactionCounterGenerator.get() > 65000) {
            transactionCounterGenerator.set(10);
        }

        if (request instanceof PlcReadRequest) {
            PlcReadRequest readRequest = (PlcReadRequest) msg.getRequest();

            for (String fieldName : readRequest.getFieldNames()) {
                PlcField field = readRequest.getField(fieldName);
                if (!(field instanceof AbEthField)) {
                    throw new PlcProtocolException("The field should have been of type AbEthField");
                }
                AbEthField abEthField = (AbEthField) field;

                DF1RequestProtectedTypedLogicalRead logicalRead = new DF1RequestProtectedTypedLogicalRead(
                    abEthField.getByteSize(), abEthField.getFileNumber(), abEthField.getFileType().getTypeCode(),
                    abEthField.getElementNumber(), (short) 0); // Subelementnumber default to zero
                // origin/sender: constant = 5
                DF1RequestMessage requestMessage = new DF1CommandRequestMessage(
                    (short) station, (short) 5, (short) 0, transactionCounterGenerator.incrementAndGet(), logicalRead);
                CIPEncapsulationReadRequest read = new CIPEncapsulationReadRequest(
                    sessionHandle, 0, emptySenderContext, 0, requestMessage);

                requests.put(requestMessage.getTransactionCounter(), msg);

                out.add(read);
            }
        } else {
            ctx.fireExceptionCaught(
                new PlcProtocolException("Unsupported request type " + request.getClass().getName()));
        }
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, CIPEncapsulationPacket packet, List<Object> out) throws Exception {
        logger.trace("Received {}, decoding...", packet);
        if(packet instanceof CIPEncapsulationConnectionResponse) {
            CIPEncapsulationConnectionResponse connectionResponse = (CIPEncapsulationConnectionResponse) packet;
            // Save the session handle
            sessionHandle = connectionResponse.getSessionHandle();

            // Tell Netty we're finished connecting
            ctx.channel().pipeline().fireUserEventTriggered(new ConnectedEvent());
        } else {
            // We're currently just expecting responses.
            if (!(packet instanceof CIPEncapsulationReadResponse)) {
                return;
            }
            CIPEncapsulationReadResponse cipResponse = (CIPEncapsulationReadResponse) packet;
            int transactionCounter = cipResponse.getResponse().getTransactionCounter();
            if(!requests.containsKey(transactionCounter)) {
                ctx.fireExceptionCaught(
                    new PlcProtocolException(
                        "Couldn't find request for response with transaction counter " + transactionCounter));
                return;
            }

            PlcRequestContainer requestContainer = requests.remove(transactionCounter);
            PlcRequest request = requestContainer.getRequest();
            PlcResponse response = null;
            if (request instanceof PlcReadRequest) {
                response = decodeReadResponse(cipResponse, requestContainer);
            } else {
                ctx.fireExceptionCaught(
                    new PlcProtocolException("Unsupported request type " + request.getClass().getName()));
            }

            // Confirm the response being handled.
            if (response != null) {
                requestContainer.getResponseFuture().complete(response);
            }
        }
    }

    private PlcResponse decodeReadResponse(
        CIPEncapsulationReadResponse plcReadResponse, PlcRequestContainer requestContainer) {

        PlcReadRequest plcReadRequest = (PlcReadRequest) requestContainer.getRequest();

        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
        for (String fieldName : plcReadRequest.getFieldNames()) {
            AbEthField field = (AbEthField) plcReadRequest.getField(fieldName);
            PlcResponseCode responseCode = decodeResponseCode(plcReadResponse.getResponse().getStatus());

            PlcValue plcValue = null;
            if (responseCode == PlcResponseCode.OK) {
                try {
                    switch (field.getFileType()) {
                        case INTEGER: // output as single bytes
                            if(plcReadResponse.getResponse() instanceof DF1CommandResponseMessageProtectedTypedLogicalRead) {
                                DF1CommandResponseMessageProtectedTypedLogicalRead df1PTLR = (DF1CommandResponseMessageProtectedTypedLogicalRead) plcReadResponse.getResponse();
                                List<Short> data = df1PTLR.getData();
                                if(data.size() == 1) {
                                    plcValue = new PlcINT(data.get(0));
                                } else {
                                    plcValue = IEC61131ValueHandler.of(data);
                                }
                            }
                            break;
                        case WORD:
                            if(plcReadResponse.getResponse() instanceof DF1CommandResponseMessageProtectedTypedLogicalRead) {
                                DF1CommandResponseMessageProtectedTypedLogicalRead df1PTLR = (DF1CommandResponseMessageProtectedTypedLogicalRead) plcReadResponse.getResponse();
                                List<Short> data = df1PTLR.getData();
                                if (((data.get(1)>> 7) & 1) == 0)  {
                                    plcValue = IEC61131ValueHandler.of((data.get(1) << 8) + data.get(0));  // positive number
                                } else {
                                    plcValue = IEC61131ValueHandler.of((((~data.get(1) & 0b01111111) << 8) + (~(data.get(0)-1) & 0b11111111))  * -1);  // negative number
                                }
                            }
                            break;
                        case DWORD:
                            if(plcReadResponse.getResponse() instanceof DF1CommandResponseMessageProtectedTypedLogicalRead) {
                                DF1CommandResponseMessageProtectedTypedLogicalRead df1PTLR = (DF1CommandResponseMessageProtectedTypedLogicalRead) plcReadResponse.getResponse();
                                List<Short> data = df1PTLR.getData();
                                if (((data.get(3)>> 7) & 1) == 0)  {
                                    plcValue = IEC61131ValueHandler.of((data.get(3) << 24) + (data.get(2) << 16) + (data.get(1) << 8) + data.get(0));  // positive number
                                } else {
                                    plcValue = IEC61131ValueHandler.of((((~data.get(3) & 0b01111111) << 24) + ((~(data.get(2)-1) & 0b11111111) << 16)+ ((~(data.get(1)-1) & 0b11111111) << 8) + (~(data.get(0)-1) & 0b11111111))  * -1);  // negative number
                                }
                            }
                            break;
                        case SINGLEBIT:
                            if(plcReadResponse.getResponse() instanceof DF1CommandResponseMessageProtectedTypedLogicalRead) {
                                DF1CommandResponseMessageProtectedTypedLogicalRead df1PTLR = (DF1CommandResponseMessageProtectedTypedLogicalRead) plcReadResponse.getResponse();
                                List<Short> data = df1PTLR.getData();
                                if (field.getBitNumber() < 8) {
                                    plcValue = IEC61131ValueHandler.of((data.get(0) & (1 <<  field.getBitNumber())) != 0);         // read from first byte
                                } else {
                                    plcValue = IEC61131ValueHandler.of((data.get(1) & (1 << (field.getBitNumber() - 8) )) != 0);   // read from second byte
                                }
                            }
                            break;
                        default:
                            logger.warn("Problem during decoding of field {}: Decoding of file type not implemented; " +
                                "FieldInformation: {}", fieldName, field);
                    }
                }
                catch (Exception e) {
                    logger.warn("Some other error occurred casting field {}, FieldInformation: {}",fieldName, field,e);
                }
            }
            ResponseItem<PlcValue> result = new ResponseItem<>(responseCode, plcValue);
            values.put(fieldName, result);
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);
    }

    private PlcResponseCode decodeResponseCode(short status) {
        if(status == 0) {
            return PlcResponseCode.OK;
        }
        return PlcResponseCode.NOT_FOUND;
    }

    private PlcValue decodeReadResponseUnsignedBytePlcValue(AbEthField field, ByteBuf data) {
        Short[] shorts = null;//readAllValues(Short.class, field, i -> data.readUnsignedByte());
        return new PlcSINT(1/*shorts*/);
    }

}
