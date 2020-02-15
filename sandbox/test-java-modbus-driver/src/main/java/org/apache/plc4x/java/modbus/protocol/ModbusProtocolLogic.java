/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.modbus.protocol;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcBoolean;
import org.apache.plc4x.java.api.value.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.config.ModbusConfiguration;
import org.apache.plc4x.java.modbus.field.ModbusFieldCoil;
import org.apache.plc4x.java.modbus.field.ModbusFieldDiscreteInput;
import org.apache.plc4x.java.modbus.field.ModbusFieldHoldingRegister;
import org.apache.plc4x.java.modbus.field.ModbusFieldInputRegister;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.modbus.readwrite.io.DataItemIO;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.optimizer.RequestTransactionManager;

import java.time.Duration;
import java.util.ArrayList;
import java.util.BitSet;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

public class ModbusProtocolLogic extends Plc4xProtocolBase<ModbusTcpADU> implements HasConfiguration<ModbusConfiguration> {

    private Duration requestTimeout;
    private short unitIdentifier;
    private RequestTransactionManager tm;
    private AtomicInteger transactionIdentifierGenerator = new AtomicInteger(10);

    @Override
    public void setConfiguration(ModbusConfiguration configuration) {
        this.requestTimeout = Duration.ofMillis(configuration.getRequestTimeout());
        this.unitIdentifier = (short) configuration.getUnitIdentifier();
        this.tm = new RequestTransactionManager(1);
        this.transactionIdentifierGenerator = new AtomicInteger(10);
    }

    @Override
    public void close(ConversationContext<ModbusTcpADU> context) {
        // Nothing to do here ...
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
        // TODO: Break down multiple items into sub-futures that are acknowledged at the end
        for (String fieldName : request.getFieldNames()) {
            PlcField field = request.getField(fieldName);
            ModbusPDU requestPdu = getRequestPdu(field);
            int transactionIdentifier = transactionIdentifierGenerator.getAndIncrement();
            ModbusTcpADU modbusTcpADU = new ModbusTcpADU(transactionIdentifier, unitIdentifier, requestPdu);
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> context.sendRequest(modbusTcpADU)
                .expectResponse(ModbusTcpADU.class, requestTimeout)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p.getTransactionIdentifier() == transactionIdentifier)
                .unwrap(ModbusTcpADU::getPdu)
                .handle(responsePdu -> {
                    // Try to decode the response data based on the corresponding request.
                    PlcValue plcValue = null;
                    PlcResponseCode responseCode;
                    try {
                        plcValue = toPlcValue(requestPdu, responsePdu);
                        responseCode = PlcResponseCode.OK;
                    } catch (ParseException e) {
                        // Add an error response code ...
                        responseCode = PlcResponseCode.INTERNAL_ERROR;
                    }

                    // Prepare the response.
                    PlcReadResponse response = new DefaultPlcReadResponse(request,
                        Collections.singletonMap(fieldName, Pair.of(responseCode, plcValue)));

                    // Pass the response back to the application.
                    future.complete(response);

                    // Finish the request-transaction.
                    transaction.endRequest();
            }));
        }
        // TODO: Merge all the sub-futures to one big response.
        return future;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        // 1. Sort all items by type:
        //      - DiscreteInput     (read-only)     --> Error
        //      - Coil              (read-write)    --> ModbusPduWriteSingleCoilRequest / ModbusPduWriteMultipleCoilsRequest
        //      - InputRegister     (read-only)     --> Error
        //      - HoldingRegister   (read-write)    --> ModbusPduWriteSingleRegisterRequest / ModbusPduWriteMultipleRegistersRequest
        //      - FifoQueue         (read-only)     --> Error
        //      - FileRecord        (read-write)    --> ModbusPduWriteFileRecordRequest
        // 2. Split up into multiple sub-requests
        return super.write(writeRequest);
    }

    private ModbusPDU getRequestPdu(PlcField field) {
        if(field instanceof ModbusFieldDiscreteInput) {
            ModbusFieldDiscreteInput discreteInput = (ModbusFieldDiscreteInput) field;
            return new ModbusPDUReadDiscreteInputsRequest(discreteInput.getAddress(), discreteInput.getQuantity());
        } else if(field instanceof ModbusFieldCoil) {
            ModbusFieldCoil coil = (ModbusFieldCoil) field;
            return new ModbusPDUReadCoilsRequest(coil.getAddress(), coil.getQuantity());
        } else if(field instanceof ModbusFieldInputRegister) {
            ModbusFieldInputRegister inputRegister = (ModbusFieldInputRegister) field;
            return new ModbusPDUReadInputRegistersRequest(inputRegister.getAddress(), inputRegister.getQuantity());
        } else if(field instanceof ModbusFieldHoldingRegister) {
            ModbusFieldHoldingRegister holdingRegister = (ModbusFieldHoldingRegister) field;
            return new ModbusPDUReadHoldingRegistersRequest(holdingRegister.getAddress(), holdingRegister.getQuantity());
        }
        throw new PlcRuntimeException("Unsupported field type " + field.getClass().getName());
    }

    private PlcValue toPlcValue(ModbusPDU request, ModbusPDU response) throws ParseException {
        if (request instanceof ModbusPDUReadDiscreteInputsRequest) {
            if (!(response instanceof ModbusPDUReadDiscreteInputsResponse)) {
                throw new PlcRuntimeException("Unexpected response type ModbusPDUReadDiscreteInputsResponse");
            }
            ModbusPDUReadDiscreteInputsRequest req = (ModbusPDUReadDiscreteInputsRequest) request;
            ModbusPDUReadDiscreteInputsResponse resp = (ModbusPDUReadDiscreteInputsResponse) response;
            return readBooleanList(req.getQuantity(), resp.getValue());
        } else if (request instanceof ModbusPDUReadCoilsRequest) {
            if (!(response instanceof ModbusPDUReadCoilsResponse)) {
                throw new PlcRuntimeException("Unexpected response type ModbusPDUReadCoilsResponse");
            }
            ModbusPDUReadCoilsRequest req = (ModbusPDUReadCoilsRequest) request;
            ModbusPDUReadCoilsResponse resp = (ModbusPDUReadCoilsResponse) response;
            return readBooleanList(req.getQuantity(), resp.getValue());
        } else if (request instanceof ModbusPDUReadInputRegistersRequest) {
            if (!(response instanceof ModbusPDUReadInputRegistersResponse)) {
                throw new PlcRuntimeException("Unexpected response type ModbusPDUReadInputRegistersResponse");
            }
            ModbusPDUReadInputRegistersRequest req = (ModbusPDUReadInputRegistersRequest) request;
            ModbusPDUReadInputRegistersResponse resp = (ModbusPDUReadInputRegistersResponse) response;
            ReadBuffer io = new ReadBuffer(resp.getValue());
            return DataItemIO.staticParse(io, (short) 2, (short) req.getQuantity());
        } else if (request instanceof ModbusPDUReadHoldingRegistersRequest) {
            if (!(response instanceof ModbusPDUReadHoldingRegistersResponse)) {
                throw new PlcRuntimeException("Unexpected response type ModbusPDUReadHoldingRegistersResponse");
            }
            ModbusPDUReadHoldingRegistersRequest req = (ModbusPDUReadHoldingRegistersRequest) request;
            ModbusPDUReadHoldingRegistersResponse resp = (ModbusPDUReadHoldingRegistersResponse) response;
            ReadBuffer io = new ReadBuffer(resp.getValue());
            return DataItemIO.staticParse(io, (short) 2, (short) req.getQuantity());
        }
        return null;
    }

    private PlcValue readBooleanList(int count, byte[] data) throws ParseException {
        ReadBuffer io = new ReadBuffer(data);
        if(count == 1) {
            return DataItemIO.staticParse(io, (short) 1, (short) 1);
        }
        // Make sure we read in all the bytes. Unfortunately when requesting 9 bytes
        // they are ordered like this: 8 7 6 5 4 3 2 1 | 0 0 0 0 0 0 0 9
        // Luckily it turns out that this is exactly how BitSet parses byte[]
        BitSet bits = BitSet.valueOf(data);
        List<PlcBoolean> result = new ArrayList<>(count);
        for(int i = 0; i < count; i++) {
            result.add(new PlcBoolean(bits.get(i)));
        }
        return new PlcList(result);
    }

}
