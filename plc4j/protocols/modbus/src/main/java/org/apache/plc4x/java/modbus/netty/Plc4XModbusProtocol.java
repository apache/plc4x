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
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.*;
import org.apache.plc4x.java.api.types.ResponseCode;
import org.apache.plc4x.java.modbus.model.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.LinkedList;
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
        WriteRequestItem<?> writeRequestItem = request.getRequestItem().orElseThrow(() -> new PlcProtocolException("Only single message supported for now"));
        // TODO: check if we can map like this. Implication is that we can only work with int, short, byte and boolean
        // TODO: for higher datatypes float, double etc we might need to split the bytes into chunks
        int quantity = writeRequestItem.getSize();
        short unitId = 0;
        checkSupportedDataType(writeRequestItem.getValues());

        ModbusAddress address = (ModbusAddress) writeRequestItem.getAddress();
        ModbusPdu modbusRequest;
        if (address instanceof RegisterAddress) {
            RegisterAddress registerAddress = (RegisterAddress) address;
            if (quantity > 1) {
                byte[] bytesToWrite = produceRegisterValue(writeRequestItem.getValues());
                modbusRequest = new WriteMultipleRegistersRequest(registerAddress.getAddress(), quantity, bytesToWrite);
            } else {
                byte[] register = produceRegisterValue(writeRequestItem.getValues());
                int intToWrite = register[0] << 8 | register[1];
                modbusRequest = new WriteSingleRegisterRequest(registerAddress.getAddress(), intToWrite);
            }
        } else if (address instanceof CoilAddress) {
            CoilAddress coilAddress = (CoilAddress) address;
            if (quantity > 1) {
                byte[] bytesToWrite = produceCoilValue(writeRequestItem.getValues());
                modbusRequest = new WriteMultipleCoilsRequest(coilAddress.getAddress(), quantity, bytesToWrite);
            } else {
                byte[] coil = produceCoilValue(writeRequestItem.getValues());
                boolean booleanToWrite = (coil[0] >> 8) == 1;
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
        ReadRequestItem<?> readRequestItem = request.getRequestItem().orElseThrow(() -> new PlcProtocolException("Only single message supported for now"));
        // TODO: check if we can map like this. Implication is that we can only work with int, short, byte and boolean
        // TODO: for higher datatypes float, double etc we might need to split the bytes into chunks
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
        LOGGER.debug("{}: transactionId: {}, unitId: {}, modbusPdu:{}", msg, msg.getTransactionId(), msg.getUnitId(), msg.getModbusPdu());
        // TODO: implement me
        short transactionId = msg.getTransactionId();
        PlcRequestContainer<PlcRequest, PlcResponse> plcRequestContainer = requestsMap.get(transactionId);
        if (plcRequestContainer == null) {
            throw new PlcProtocolException("Unrelated payload received" + msg);
        }

        // TODO: only single Item supported for now
        PlcRequest<?> request = plcRequestContainer.getRequest();
        RequestItem requestItem = request.getRequestItem().orElseThrow(() -> new PlcProtocolException("Only single message supported for now"));
        Class datatype = requestItem.getDatatype();

        ModbusPdu modbusPdu = msg.getModbusPdu();
        short unitId = msg.getUnitId();

        if (modbusPdu instanceof WriteMultipleCoilsResponse) {
            // TODO: finish implementation
            WriteMultipleCoilsResponse writeMultipleCoilsResponse = (WriteMultipleCoilsResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, quantity:{}", writeMultipleCoilsResponse, writeMultipleCoilsResponse.getAddress(), writeMultipleCoilsResponse.getQuantity());
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem) requestItem, ResponseCode.OK)));
        } else if (modbusPdu instanceof WriteMultipleRegistersResponse) {
            // TODO: finish implementation
            WriteMultipleRegistersResponse writeMultipleRegistersResponse = (WriteMultipleRegistersResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, quantity:{}", writeMultipleRegistersResponse, writeMultipleRegistersResponse.getAddress(), writeMultipleRegistersResponse.getQuantity());
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem) requestItem, ResponseCode.OK)));
        } else if (modbusPdu instanceof WriteSingleCoilResponse) {
            // TODO: finish implementation
            WriteSingleCoilResponse writeSingleCoilResponse = (WriteSingleCoilResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, value:{}", writeSingleCoilResponse, writeSingleCoilResponse.getAddress(), writeSingleCoilResponse.getValue());
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem) requestItem, ResponseCode.OK)));
        } else if (modbusPdu instanceof WriteSingleRegisterResponse) {
            // TODO: finish implementation
            WriteSingleRegisterResponse writeSingleRegisterResponse = (WriteSingleRegisterResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, value:{}", writeSingleRegisterResponse, writeSingleRegisterResponse.getAddress(), writeSingleRegisterResponse.getValue());
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem) requestItem, ResponseCode.OK)));
        } else if (modbusPdu instanceof ReadCoilsResponse) {
            // TODO: finish implementation
            ReadCoilsResponse readCoilsResponse = (ReadCoilsResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readCoilsResponse);
            ByteBuf byteBuf = readCoilsResponse.getCoilStatus();
            List data = produceCoilValueList(requestItem, datatype, byteBuf);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem) requestItem, ResponseCode.OK, data)));
        } else if (modbusPdu instanceof ReadDiscreteInputsResponse) {
            // TODO: finish implementation
            ReadDiscreteInputsResponse readDiscreteInputsResponse = (ReadDiscreteInputsResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readDiscreteInputsResponse);
            ByteBuf byteBuf = readDiscreteInputsResponse.getInputStatus();
            List data = produceCoilValueList(requestItem, datatype, byteBuf);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem) requestItem, ResponseCode.OK, data)));
        } else if (modbusPdu instanceof ReadHoldingRegistersResponse) {
            // TODO: finish implementation
            ReadHoldingRegistersResponse readHoldingRegistersResponse = (ReadHoldingRegistersResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readHoldingRegistersResponse);
            ByteBuf byteBuf = readHoldingRegistersResponse.getRegisters();
            // TODO: use register method
            List data = produceRegisterValueList(requestItem, datatype, byteBuf);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem) requestItem, ResponseCode.OK, data)));
        } else if (modbusPdu instanceof ReadInputRegistersResponse) {
            // TODO: finish implementation
            ReadInputRegistersResponse readInputRegistersResponse = (ReadInputRegistersResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readInputRegistersResponse);
            ByteBuf byteBuf = readInputRegistersResponse.getRegisters();
            // TODO: use register method
            List data = produceRegisterValueList(requestItem, datatype, byteBuf);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem) requestItem, ResponseCode.OK, data)));
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

    private void checkSupportedDataType(List<?> values) {
        if (values == null || values.size() == 0) {
            return;
        }
        for (Object value : values) {
            if (
                !(value instanceof Boolean)
                    && !(value instanceof Byte)
                    && !(value instanceof byte[])
                    && !(value instanceof Short)
                    && !(value instanceof Integer)
                ) {
                throw new PlcRuntimeException("Unsupported datatype detected " + value.getClass());
            }
        }
    }

    ////////////////////////////////////////////////////////////////////////////////
    // Encoding helpers.
    ////////////////////////////////////////////////////////////////////////////////
    private byte[] produceCoilValue(List<?> values) throws PlcProtocolException {
        List<Byte> coils = new LinkedList<>();
        Byte actualCoil = 0;
        int i = 7;
        for (Object value : values) {
            boolean coilSet = false;
            if (value.getClass() == Boolean.class) {
                coilSet = (Boolean) value;
            } else if (value.getClass() == Byte.class) {
                if ((byte) value > 1) {
                    throw new PlcProtocolException("Value to high to fit into Byte: " + value);
                }
                coilSet = (byte) value == 1;
            } else if (value.getClass() == byte[].class) {
                byte[] bytes = (byte[]) value;
                if (bytes.length != 1) {
                    throw new PlcProtocolException("Only exactly one byte is allowed: " + bytes.length);
                }
                byte byteValue = bytes[1];
                if (byteValue > 1) {
                    throw new PlcProtocolException("Value to high to fit into Byte: " + value);
                }
                coilSet = byteValue == 1;
            } else if (value.getClass() == Short.class) {
                if ((short) value > 1) {
                    throw new PlcProtocolException("Value to high to fit into Byte: " + value);
                }
                coilSet = (Short) value == 1;
            } else if (value.getClass() == Integer.class) {
                if ((int) value > 1) {
                    throw new PlcProtocolException("Value to high to fit into Byte: " + value);
                }
                coilSet = (int) value == 1;
            }
            byte coilToSet = (coilSet ? (byte) 1 : (byte) 0);
            actualCoil = (byte) (actualCoil | coilToSet << i);
            i--;
            if (i < 0) {
                coils.add(actualCoil);
                actualCoil = 0;
                i = 8;
            }
        }
        // TODO: ensure we have a least (quantity + 7) / 8 = N bytes
        return ArrayUtils.toPrimitive(coils.toArray(new Byte[0]));
    }

    private byte[] produceRegisterValue(List<?> values) throws PlcProtocolException {
        ByteBuf buffer = Unpooled.buffer();
        for (Object value : values) {
            if (value.getClass() == Boolean.class) {
                buffer.writeByte(0);
                buffer.writeByte((boolean) value ? 1 : 0);
            } else if (value.getClass() == Byte.class) {
                buffer.writeByte(0);
                buffer.writeByte((byte) value);
            } else if (value.getClass() == byte[].class) {
                byte[] bytes = (byte[]) value;
                if (bytes.length != 2) {
                    throw new PlcProtocolException("Only exactly two bytes are allowed: " + bytes.length);
                }
                buffer.writeBytes(bytes);
            } else if (value.getClass() == Short.class) {
                buffer.writeShort((int) value);
            } else if (value.getClass() == Integer.class) {
                if ((int) value > Integer.MAX_VALUE) {
                    throw new PlcProtocolException("Value to high to fit into register: " + value);
                }
                buffer.writeShort((int) value);
            }
        }
        // TODO: ensure we have a least quantity * 2 = N bytes
        byte[] result = new byte[buffer.writerIndex()];
        buffer.readBytes(result);
        return result;
    }

    ////////////////////////////////////////////////////////////////////////////////
    // Decoding helpers.
    ////////////////////////////////////////////////////////////////////////////////

    @SuppressWarnings("unchecked")
    private List produceCoilValueList(RequestItem requestItem, Class datatype, ByteBuf byteBuf) {
        ReadRequestItem readRequestItem = (ReadRequestItem) requestItem;
        byte[] bytes = new byte[byteBuf.readableBytes()];
        byteBuf.readBytes(bytes);
        List data = new LinkedList();
        for (int i = 0, j = 0; i < readRequestItem.getSize(); i++) {
            if (i != 0 && i % 8 == 0) {
                // Every 8 Coils we need to increase the access
                j++;
            }
            Boolean coilValue = (1 << i & bytes[j]) == 1;
            if (datatype == Boolean.class) {
                data.add(coilValue);
            } else if (datatype == Byte.class) {
                data.add((byte) (coilValue ? 1 : 0));
            } else if (datatype == byte[].class) {
                data.add(new byte[]{(byte) (coilValue ? 1 : 0)});
            } else if (datatype == Short.class) {
                data.add((short) (coilValue ? 1 : 0));
            } else if (datatype == Integer.class) {
                data.add(coilValue ? 1 : 0);
            }
        }
        return data;
    }

    @SuppressWarnings("unchecked")
    private List produceRegisterValueList(RequestItem requestItem, Class datatype, ByteBuf byteBuf) throws PlcProtocolException {
        ReadRequestItem readRequestItem = (ReadRequestItem) requestItem;
        int readableBytes = byteBuf.readableBytes();
        if (readableBytes % 2 != 0) {
            throw new PlcProtocolException("Readables bytes should even: " + readableBytes);
        }
        List data = new LinkedList();
        for (int i = 0; i < readRequestItem.getSize(); i++) {
            byte[] register = new byte[2];
            byteBuf.readBytes(register);
            int intValue = register[0] << 8 | register[1];
            if (datatype == Boolean.class) {
                data.add(intValue == 1);
            } else if (datatype == Byte.class) {
                if (intValue > Byte.MAX_VALUE) {
                    throw new PlcProtocolException("Value to high to fit into Byte: " + intValue);
                }
                data.add((byte) intValue);
            } else if (datatype == byte[].class) {
                data.add(register);
            } else if (datatype == Short.class) {
                if (intValue > Short.MAX_VALUE) {
                    throw new PlcProtocolException("Value to high to fit into Short: " + intValue);
                }
                data.add((short) intValue);
            } else if (datatype == Integer.class) {
                data.add(intValue);
            }
        }
        return data;
    }
}
