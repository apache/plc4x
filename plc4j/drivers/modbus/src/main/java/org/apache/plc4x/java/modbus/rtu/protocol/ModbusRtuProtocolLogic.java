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
package org.apache.plc4x.java.modbus.rtu.protocol;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.base.protocol.ModbusProtocolLogic;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHandler;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.modbus.rtu.config.ModbusRtuConfiguration;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;

import java.time.Duration;
import java.util.Collections;
import java.util.concurrent.CompletableFuture;

public class ModbusRtuProtocolLogic extends ModbusProtocolLogic<ModbusRtuADU> implements HasConfiguration<ModbusRtuConfiguration> {

    public ModbusRtuProtocolLogic() {
        super(DriverType.MODBUS_RTU);
    }

    @Override
    public void setConfiguration(ModbusRtuConfiguration configuration) {
        this.requestTimeout = Duration.ofMillis(configuration.getRequestTimeout());
        this.unitIdentifier = (short) configuration.getDefaultUnitIdentifier();
        this.defaultPayloadByteOrder = configuration.getDefaultPayloadByteOrder();
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new ModbusTagHandler();
    }

    @Override
    public void close(ConversationContext<ModbusRtuADU> context) {
        tm.shutdown();
    }

    @Override
    public CompletableFuture<PlcPingResponse> ping(PlcPingRequest pingRequest) {
        CompletableFuture<PlcPingResponse> future = new CompletableFuture<>();

        // As it seems that even, if Modbus defines a DeviceIdentificationRequest, no device actually implements this.
        // So we fall back to a request, that most certainly is implemented by any device. Even if the device doesn't
        // have any holding-register:1, it should still gracefully respond.
        ModbusPDU readRequestPdu = getReadRequestPdu(pingAddress);
        final short unitId = getUnitId(pingAddress);
        ModbusRtuADU modbusRtuADU = new ModbusRtuADU(unitId, readRequestPdu);

        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(modbusRtuADU)
            .expectResponse(ModbusRtuADU.class, requestTimeout)
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .unwrap(ModbusRtuADU::getPdu)
            .handle(responsePdu -> {
                transaction.endRequest();
                // We really don't care about what we got back. As long as it's a Modbus PDU, we're ok.
                future.complete(new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK));
            }));
        return future;
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;

        // 1. Sort all items by type:
        //      - DiscreteInput     (read-only)     --> ModbusPduReadDiscreteInputsRequest
        //      - Coil              (read-write)    --> ModbusPduReadCoilsRequest
        //      - InputRegister     (read-only)     --> ModbusPduReadInputRegistersRequest
        //      - HoldingRegister   (read-write)    --> ModbusPduReadHoldingRegistersRequest
        //      - FifoQueue         (read-only)     --> ModbusPduReadFifoQueueRequest
        //      - FileRecord        (read-write)    --> ModbusPduReadFileRecordRequest
        // 2. Split up into multiple sub-requests

        // Example for sending a request ...
        if (request.getTagNames().size() == 1) {
            String tagName = request.getTagNames().iterator().next();
            ModbusTag tag = (ModbusTag) request.getTag(tagName);
            final ModbusPDU requestPdu = getReadRequestPdu(tag);
            final short unitId = getUnitId(tag);

            ModbusRtuADU modbusRtuADU = new ModbusRtuADU(unitId, requestPdu);
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> conversationContext.sendRequest(modbusRtuADU)
                .expectResponse(ModbusRtuADU.class, requestTimeout)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .unwrap(ModbusRtuADU::getPdu)
                .handle(responsePdu -> {
                    // Try to decode the response data based on the corresponding request.
                    PlcValue plcValue = null;
                    PlcResponseCode responseCode;
                    // Check if the response was an error response.
                    if (responsePdu instanceof ModbusPDUError) {
                        ModbusPDUError errorResponse = (ModbusPDUError) responsePdu;
                        responseCode = getErrorCode(errorResponse);
                    } else {
                        try {
                            ModbusByteOrder byteOrder = defaultPayloadByteOrder;
                            if(tag.getByteOrder() != null) {
                                byteOrder = tag.getByteOrder();
                            }
                            plcValue = toPlcValue(requestPdu, responsePdu, tag.getDataType(), byteOrder);
                            responseCode = PlcResponseCode.OK;
                        } catch (ParseException e) {
                            // Add an error response code ...
                            responseCode = PlcResponseCode.INTERNAL_ERROR;
                        }
                    }

                    // Prepare the response.
                    PlcReadResponse response = new DefaultPlcReadResponse(request,
                        Collections.singletonMap(tagName, new DefaultPlcResponseItem<>(responseCode, plcValue)));

                    // Pass the response back to the application.
                    future.complete(response);

                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        } else {
            future.completeExceptionally(new PlcRuntimeException("Modbus only supports single filed requests"));
        }
        return future;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;

        // 1. Sort all items by type:
        //      - DiscreteInput     (read-only)     --> Error
        //      - Coil              (read-write)    --> ModbusPduWriteSingleCoilRequest / ModbusPduWriteMultipleCoilsRequest
        //      - InputRegister     (read-only)     --> Error
        //      - HoldingRegister   (read-write)    --> ModbusPduWriteSingleRegisterRequest / ModbusPduWriteMultipleRegistersRequest
        //      - FifoQueue         (read-only)     --> Error
        //      - FileRecord        (read-write)    --> ModbusPduWriteFileRecordRequest
        // 2. Split up into multiple sub-requests
        if (request.getTagNames().size() == 1) {
            String tagName = request.getTagNames().iterator().next();
            PlcTag tag = request.getTag(tagName);
            final ModbusPDU requestPdu = getWriteRequestPdu(tag, writeRequest.getPlcValue(tagName));
            final short unitId = getUnitId(tag);
            ModbusRtuADU modbusRtuADU = new ModbusRtuADU(unitId, requestPdu);
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> conversationContext.sendRequest(modbusRtuADU)
                .expectResponse(ModbusRtuADU.class, requestTimeout)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .unwrap(ModbusRtuADU::getPdu)
                .handle(responsePdu -> {
                    // Try to decode the response data based on the corresponding request.
                    PlcResponseCode responseCode;

                    // Check if the response was an error response.
                    if (responsePdu instanceof ModbusPDUError) {
                        ModbusPDUError errorResponse = (ModbusPDUError) responsePdu;
                        responseCode = getErrorCode(errorResponse);
                    } else {
                        responseCode = PlcResponseCode.OK;
                        // TODO: Check the correct number of elements were written.
                        if (responsePdu instanceof ModbusPDUWriteSingleCoilResponse) {
                            ModbusPDUWriteSingleCoilResponse response = (ModbusPDUWriteSingleCoilResponse) responsePdu;
                            ModbusPDUWriteSingleCoilRequest requestSingleCoil = (ModbusPDUWriteSingleCoilRequest) requestPdu;
                            if (!((response.getValue() == requestSingleCoil.getValue()) && (response.getAddress() == requestSingleCoil.getAddress()))) {
                                responseCode = PlcResponseCode.REMOTE_ERROR;
                            }
                        }
                    }

                    // Prepare the response.
                    PlcWriteResponse response = new DefaultPlcWriteResponse(request,
                        Collections.singletonMap(tagName, responseCode));

                    // Pass the response back to the application.
                    future.complete(response);

                    // Finish the request-transaction.
                    transaction.endRequest();
                }));

        } else {
            future.completeExceptionally(new PlcRuntimeException("Modbus only supports single filed requests"));
        }
        return future;
    }

    @Override
    protected void decode(ConversationContext<ModbusRtuADU> context, ModbusRtuADU msg) throws Exception {
        System.out.println(msg);
    }

}
