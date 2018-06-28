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
package org.apache.plc4x.java.modbus.netty;

import com.digitalpetri.modbus.ModbusPdu;
import com.digitalpetri.modbus.codec.ModbusTcpPayload;
import com.digitalpetri.modbus.requests.*;
import com.digitalpetri.modbus.responses.*;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.*;
import org.apache.plc4x.java.api.types.ResponseCode;
import org.apache.plc4x.java.modbus.model.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4XModbusProtocol extends MessageToMessageCodec<ModbusTcpPayload, PlcRequestContainer<PlcRequest, PlcResponse>> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XModbusProtocol.class);

    public final AtomicInteger transactionId = new AtomicInteger();

    private final ConcurrentMap<Short, PlcRequestContainer<PlcRequest, PlcResponse>> requestsMap = new ConcurrentHashMap<>();

    public Plc4XModbusProtocol() {
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer<PlcRequest, PlcResponse> msg, List<Object> out) throws Exception {
        LOGGER.trace("(<--OUT): {}, {}, {}", ctx, msg, out);
        // Reset transactionId on overflow
        transactionId.compareAndSet(Short.MAX_VALUE + 1, 0);
        PlcRequest request = msg.getRequest();
        if (request instanceof PlcReadRequest) {
            encodeReadRequest(msg, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(msg, out);
        }
    }

    private void encodeWriteRequest(PlcRequestContainer<PlcRequest, PlcResponse> msg, List<Object> out) throws PlcException {
        PlcWriteRequest request = (PlcWriteRequest) msg.getRequest();
        // TODO: support multiple requests
        WriteRequestItem<?> writeRequestItem = request.getRequestItem().get();
        int quantity = writeRequestItem.getSize();
        short unitId = 0;

        ModbusAddress address = (ModbusAddress) writeRequestItem.getAddress();
        ModbusPdu modbusRequest;
        if (address instanceof RegisterAddress) {
            RegisterAddress registerAddress = (RegisterAddress) address;
            if (quantity > 1) {
                // TODO: ensure we have a least quantity * 2 = N bytes
                byte[] bytesToWrite = (byte[]) writeRequestItem.getValues().get(0);
                modbusRequest = new WriteMultipleRegistersRequest(registerAddress.getAddress(), quantity, bytesToWrite);
            } else {
                int intToWrite = (Integer) writeRequestItem.getValues().get(0);
                modbusRequest = new WriteSingleRegisterRequest(registerAddress.getAddress(), intToWrite);
            }
        } else if (address instanceof CoilAddress) {
            CoilAddress coilAddress = (CoilAddress) address;
            if (quantity > 1) {
                // TODO: ensure we have a least (quantity + 7) / 8 = N bytes
                byte[] bytesToWrite = (byte[]) writeRequestItem.getValues().get(0);
                modbusRequest = new WriteMultipleCoilsRequest(coilAddress.getAddress(), quantity, bytesToWrite);
            } else {
                boolean booleanToWrite = (Boolean) writeRequestItem.getValues().get(0);
                modbusRequest = new WriteSingleCoilRequest(coilAddress.getAddress(), booleanToWrite);
            }
        } else {
            throw new PlcProtocolException("Unsupported address type" + address.getClass());
        }
        short transactionId = (short) this.transactionId.getAndIncrement();
        requestsMap.put(transactionId, msg);
        out.add(new ModbusTcpPayload(transactionId, unitId, modbusRequest));
    }

    private void encodeReadRequest(PlcRequestContainer<PlcRequest, PlcResponse> msg, List<Object> out) throws PlcException {
        PlcReadRequest request = (PlcReadRequest) msg.getRequest();
        // TODO: support multiple requests
        ReadRequestItem<?> readRequestItem = request.getRequestItem().get();
        int quantity = readRequestItem.getSize();
        // TODO: the unit the should be used for multiple Requests
        short unitId = 0;

        ModbusAddress address = (ModbusAddress) readRequestItem.getAddress();
        ModbusPdu modbusRequest;
        if (address instanceof CoilAddress) {
            CoilAddress coilAddress = (CoilAddress) address;
            modbusRequest = new ReadCoilsRequest(coilAddress.getAddress(), quantity);
        } else if (address instanceof RegisterAddress) {
            RegisterAddress registerAddress = (RegisterAddress) address;
            modbusRequest = new ReadHoldingRegistersRequest(registerAddress.getAddress(), quantity);
        } else if (address instanceof ReadDiscreteInputsModbusAddress) {
            ReadDiscreteInputsModbusAddress readDiscreteInputsModbusAddress = (ReadDiscreteInputsModbusAddress) address;
            modbusRequest = new ReadDiscreteInputsRequest(readDiscreteInputsModbusAddress.getAddress(), quantity);
        } else if (address instanceof ReadHoldingRegistersModbusAddress) {
            ReadHoldingRegistersModbusAddress readHoldingRegistersModbusAddress = (ReadHoldingRegistersModbusAddress) address;
            modbusRequest = new ReadHoldingRegistersRequest(readHoldingRegistersModbusAddress.getAddress(), quantity);
        } else if (address instanceof ReadInputRegistersModbusAddress) {
            ReadInputRegistersModbusAddress readInputRegistersModbusAddress = (ReadInputRegistersModbusAddress) address;
            modbusRequest = new ReadInputRegistersRequest(readInputRegistersModbusAddress.getAddress(), quantity);
        } else {
            throw new PlcProtocolException("Unsupported address type" + address.getClass());
        }
        short transactionId = (short) this.transactionId.getAndIncrement();
        requestsMap.put(transactionId, msg);
        out.add(new ModbusTcpPayload(transactionId, unitId, modbusRequest));
    }

    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext ctx, ModbusTcpPayload msg, List<Object> out) throws Exception {
        LOGGER.trace("(-->IN): {}, {}, {}", ctx, msg, out);
        // TODO: implement me
        short transactionId = msg.getTransactionId();
        PlcRequestContainer<PlcRequest, PlcResponse> plcRequestContainer = requestsMap.get(transactionId);
        if (plcRequestContainer == null) {
            throw new PlcProtocolException("Unrelated payload received" + msg);
        }

        // TODO: only single Item supported for now
        PlcRequest request = plcRequestContainer.getRequest();
        RequestItem requestItem = (RequestItem) request.getRequestItem().get();

        ModbusPdu modbusPdu = msg.getModbusPdu();
        short unitId = msg.getUnitId();

        if (modbusPdu instanceof WriteMultipleCoilsResponse) {
            // TODO: finish implementation
            WriteMultipleCoilsResponse writeMultipleCoilsResponse = (WriteMultipleCoilsResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, quantity:{}", writeMultipleCoilsResponse, writeMultipleCoilsResponse.getAddress(), writeMultipleCoilsResponse.getQuantity());
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem<? extends Object>) requestItem, ResponseCode.OK)));
        } else if (modbusPdu instanceof WriteMultipleRegistersResponse) {
            // TODO: finish implementation
            WriteMultipleRegistersResponse writeMultipleRegistersResponse = (WriteMultipleRegistersResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, quantity:{}", writeMultipleRegistersResponse, writeMultipleRegistersResponse.getAddress(), writeMultipleRegistersResponse.getQuantity());
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem<? extends Object>) requestItem, ResponseCode.OK)));
        } else if (modbusPdu instanceof WriteSingleCoilResponse) {
            // TODO: finish implementation
            WriteSingleCoilResponse writeSingleCoilResponse = (WriteSingleCoilResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, value:{}", writeSingleCoilResponse, writeSingleCoilResponse.getAddress(), writeSingleCoilResponse.getValue());
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem<? extends Object>) requestItem, ResponseCode.OK)));
        } else if (modbusPdu instanceof WriteSingleRegisterResponse) {
            // TODO: finish implementation
            WriteSingleRegisterResponse writeSingleRegisterResponse = (WriteSingleRegisterResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, value:{}", writeSingleRegisterResponse, writeSingleRegisterResponse.getAddress(), writeSingleRegisterResponse.getValue());
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem<? extends Object>) requestItem, ResponseCode.OK)));
        } else if (modbusPdu instanceof ReadCoilsResponse) {
            // TODO: finish implementation
            ReadCoilsResponse readCoilsResponse = (ReadCoilsResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readCoilsResponse);
            ByteBuf byteBuf = readCoilsResponse.getCoilStatus();
            byte[] bytes = new byte[byteBuf.readableBytes()];
            byteBuf.readBytes(bytes);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem<? extends Object>) requestItem, ResponseCode.OK, (Object) bytes)));
        } else if (modbusPdu instanceof ReadDiscreteInputsResponse) {
            // TODO: finish implementation
            ReadDiscreteInputsResponse readDiscreteInputsResponse = (ReadDiscreteInputsResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readDiscreteInputsResponse);
            ByteBuf byteBuf = readDiscreteInputsResponse.getInputStatus();
            byte[] bytes = new byte[byteBuf.readableBytes()];
            byteBuf.readBytes(bytes);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem<? extends Object>) requestItem, ResponseCode.OK, (Object) bytes)));
        } else if (modbusPdu instanceof ReadHoldingRegistersResponse) {
            // TODO: finish implementation
            ReadHoldingRegistersResponse readHoldingRegistersResponse = (ReadHoldingRegistersResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readHoldingRegistersResponse);
            ByteBuf byteBuf = readHoldingRegistersResponse.getRegisters();
            byte[] bytes = new byte[byteBuf.readableBytes()];
            byteBuf.readBytes(bytes);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem<? extends Object>) requestItem, ResponseCode.OK, (Object) bytes)));
        } else if (modbusPdu instanceof ReadInputRegistersResponse) {
            // TODO: finish implementation
            ReadInputRegistersResponse readInputRegistersResponse = (ReadInputRegistersResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readInputRegistersResponse);
            ByteBuf byteBuf = readInputRegistersResponse.getRegisters();
            byte[] bytes = new byte[byteBuf.readableBytes()];
            byteBuf.readBytes(bytes);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem<? extends Object>) requestItem, ResponseCode.OK, (Object) bytes)));
        } else if (modbusPdu instanceof ExceptionResponse) {
            ExceptionResponse exceptionResponse = (ExceptionResponse) modbusPdu;
            throw new PlcProtocolException("Error received " + exceptionResponse.getExceptionCode());
        } else {
            throw new PlcProtocolException("Unsupported messageTyp type" + modbusPdu.getClass());
        }
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        LOGGER.trace("(-->ERR): {}", ctx, cause);
        super.exceptionCaught(ctx, cause);
    }

    ////////////////////////////////////////////////////////////////////////////////
    // Encoding helpers.
    ////////////////////////////////////////////////////////////////////////////////


    ////////////////////////////////////////////////////////////////////////////////
    // Decoding helpers.
    ////////////////////////////////////////////////////////////////////////////////


}
