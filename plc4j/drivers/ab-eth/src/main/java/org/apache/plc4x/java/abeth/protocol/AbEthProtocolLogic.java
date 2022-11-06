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

import org.apache.plc4x.java.abeth.configuration.AbEthConfiguration;
import org.apache.plc4x.java.abeth.field.AbEthField;
import org.apache.plc4x.java.abeth.readwrite.*;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

public class AbEthProtocolLogic extends Plc4xProtocolBase<CIPEncapsulationPacket> implements HasConfiguration<AbEthConfiguration> {

    private static final Logger logger = LoggerFactory.getLogger(AbEthProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private static final List<Short> emptySenderContext = Arrays.asList((short) 0x00 ,(short) 0x00 ,(short) 0x00,
        (short) 0x00,(short) 0x00,(short) 0x00, (short) 0x00,(short) 0x00);

    private AbEthConfiguration configuration;

    private final AtomicInteger transactionCounterGenerator = new AtomicInteger(10);
    private RequestTransactionManager tm;
    private long sessionHandle;

    @Override
    public void setConfiguration(AbEthConfiguration configuration) {
        this.configuration = configuration;
        // Set the transaction manager to allow only one message at a time.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void onConnect(ConversationContext<CIPEncapsulationPacket> context) {
        logger.debug("Sending COTP Connection Request");
        CIPEncapsulationConnectionRequest connectionRequest =
            new CIPEncapsulationConnectionRequest(0L, 0L, emptySenderContext, 0L);
        context.sendRequest(connectionRequest)
            .expectResponse(CIPEncapsulationPacket.class, REQUEST_TIMEOUT)
            .check(p -> p instanceof CIPEncapsulationConnectionResponse)
            .unwrap(p -> (CIPEncapsulationConnectionResponse) p)
            .handle(cipEncapsulationConnectionResponse -> {
                sessionHandle = cipEncapsulationConnectionResponse.getSessionHandle();
                // Send an event that connection setup is complete.
                context.fireConnected();
            });
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        // TODO: Warning ... we are senging one request per field ... the result has to be merged back together ...
        for (String fieldName : readRequest.getFieldNames()) {
            PlcField field = readRequest.getField(fieldName);
            if (!(field instanceof AbEthField)) {
                logger.error("The field should have been of type AbEthField");
            }
            AbEthField abEthField = (AbEthField) field;

            DF1RequestProtectedTypedLogicalRead logicalRead = new DF1RequestProtectedTypedLogicalRead(
                abEthField.getByteSize(), abEthField.getFileNumber(), abEthField.getFileType().getTypeCode(),
                abEthField.getElementNumber(), (short) 0); // Subelementnumber default to zero

            final int transactionCounter = transactionCounterGenerator.incrementAndGet();
            // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
            if(transactionCounterGenerator.get() == 0xFFFF) {
                transactionCounterGenerator.set(1);
            }
// origin/sender: constant = 5
            DF1RequestMessage requestMessage = new DF1CommandRequestMessage(
                (short) configuration.getStation(), (short) 5, (short) 0,
                transactionCounter, logicalRead);
            CIPEncapsulationReadRequest read = new CIPEncapsulationReadRequest(
                sessionHandle, 0, emptySenderContext, 0, requestMessage);

            CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> context.sendRequest(read)
                .expectResponse(CIPEncapsulationPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CIPEncapsulationReadResponse)
                .unwrap(p -> (CIPEncapsulationReadResponse) p)
                .check(p -> p.getResponse().getTransactionCounter() == transactionCounter)
                .handle(p -> {
                    PlcResponse response = decodeReadResponse(p, readRequest);

                    // TODO: Not sure how to merge things back together ...

                    //future.complete(response);
                    // Finish the request-transaction.
                    transaction.endRequest();
//                    future.complete(((PlcReadResponse) decodeReadResponse(p, ((InternalPlcReadRequest) readRequest))));
                }));

            // TODO: This aborts reading other fields after sending the first fields request ... refactor.
            return future;
        }
        // TODO: Should return an aggregated future ....
        return null;
    }

    @Override
    public void close(ConversationContext<CIPEncapsulationPacket> context) {

    }

    private PlcResponse decodeReadResponse(
        CIPEncapsulationReadResponse plcReadResponse, PlcReadRequest plcReadRequest) {
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
                                    plcValue = PlcValueHandler.of(data);
                                }
                            }
                            break;
                        case WORD:
                            if(plcReadResponse.getResponse() instanceof DF1CommandResponseMessageProtectedTypedLogicalRead) {
                                DF1CommandResponseMessageProtectedTypedLogicalRead df1PTLR = (DF1CommandResponseMessageProtectedTypedLogicalRead) plcReadResponse.getResponse();
                                List<Short> data = df1PTLR.getData();
                                if (((data.get(1)>> 7) & 1) == 0)  {
                                    plcValue = PlcValueHandler.of((data.get(1) << 8) + data.get(0));  // positive number
                                } else {
                                    plcValue = PlcValueHandler.of((((~data.get(1) & 0b01111111) << 8) + (~(data.get(0)-1) & 0b11111111))  * -1);  // negative number
                                }
                            }
                            break;
                        case DWORD:
                            if(plcReadResponse.getResponse() instanceof DF1CommandResponseMessageProtectedTypedLogicalRead) {
                                DF1CommandResponseMessageProtectedTypedLogicalRead df1PTLR = (DF1CommandResponseMessageProtectedTypedLogicalRead) plcReadResponse.getResponse();
                                List<Short> data = df1PTLR.getData();
                                if (((data.get(3)>> 7) & 1) == 0)  {
                                    plcValue = PlcValueHandler.of((data.get(3) << 24) + (data.get(2) << 16) + (data.get(1) << 8) + data.get(0));  // positive number
                                } else {
                                    plcValue = PlcValueHandler.of((((~data.get(3) & 0b01111111) << 24) + ((~(data.get(2)-1) & 0b11111111) << 16)+ ((~(data.get(1)-1) & 0b11111111) << 8) + (~(data.get(0)-1) & 0b11111111))  * -1);  // negative number
                                }
                            }
                            break;
                        case SINGLEBIT:
                            if(plcReadResponse.getResponse() instanceof DF1CommandResponseMessageProtectedTypedLogicalRead) {
                                DF1CommandResponseMessageProtectedTypedLogicalRead df1PTLR = (DF1CommandResponseMessageProtectedTypedLogicalRead) plcReadResponse.getResponse();
                                List<Short> data = df1PTLR.getData();
                                if (field.getBitNumber() < 8) {
                                    plcValue = PlcValueHandler.of((data.get(0) & (1 <<  field.getBitNumber())) != 0);         // read from first byte
                                } else {
                                    plcValue = PlcValueHandler.of((data.get(1) & (1 << (field.getBitNumber() - 8) )) != 0);   // read from second byte
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

        // TODO: Double check if it's really a InternalPlcReadRequest ...
        return new DefaultPlcReadResponse(plcReadRequest, values);
    }

    private PlcResponseCode decodeResponseCode(short status) {
        if(status == 0) {
            return PlcResponseCode.OK;
        }
        return PlcResponseCode.NOT_FOUND;
    }

}
