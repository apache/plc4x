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
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.*;
import org.apache.plc4x.java.api.types.ResponseCode;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.modbus.model.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicInteger;


public class Plc4XModbusProtocol extends MessageToMessageCodec<ModbusTcpPayload, PlcRequestContainer<PlcRequest, PlcResponse>> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XModbusProtocol.class);

    public final AtomicInteger transactionId = new AtomicInteger();

    private final ConcurrentMap<Short, PlcRequestContainer<PlcRequest, PlcResponse>> requestsMap = new ConcurrentHashMap<>();

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
        WriteRequestItem<?> writeRequestItem = request.getRequestItem().orElseThrow(() -> new PlcNotImplementedException("Only single message supported for now"));
        // TODO: check if we can map like this. Implication is that we can only work with int, short, byte and boolean
        // TODO: for higher datatypes float, double etc we might need to split the bytes into chunks
        int quantity = writeRequestItem.getSize();
        short unitId = 0;

        ModbusAddress address = (ModbusAddress) writeRequestItem.getAddress();
        ModbusPdu modbusRequest;
        if (address instanceof RegisterModbusAddress) {
            RegisterModbusAddress registerModbusAddress = (RegisterModbusAddress) address;
            if (quantity > 1) {
                byte[] bytesToWrite = produceRegisterValue(writeRequestItem.getValues());
                int requiredLength = 2 * quantity;
                if (bytesToWrite.length != requiredLength) {
                    throw new PlcProtocolException("Invalid register values created. Should be at least quantity * 2 = N bytes. Was " + bytesToWrite.length + ", expected " + requiredLength);
                }
                modbusRequest = new WriteMultipleRegistersRequest(registerModbusAddress.getAddress(), quantity, bytesToWrite);
            } else {
                byte[] register = produceRegisterValue(writeRequestItem.getValues());
                int intToWrite = register[0] << 8 | register[1] & 0xff;
                modbusRequest = new WriteSingleRegisterRequest(registerModbusAddress.getAddress(), intToWrite);
            }
        } else if (address instanceof CoilModbusAddress) {
            CoilModbusAddress coilModbusAddress = (CoilModbusAddress) address;
            if (quantity > 1) {
                byte[] bytesToWrite = produceCoilValues(writeRequestItem.getValues());
                int requiredLength = (quantity + 7) / 8;
                if (bytesToWrite.length != requiredLength) {
                    throw new PlcProtocolException("Invalid coil values created. Should be at least (quantity + 7) / 8 = N bytes. Was " + bytesToWrite.length + ", expected " + requiredLength);
                }
                modbusRequest = new WriteMultipleCoilsRequest(coilModbusAddress.getAddress(), quantity, bytesToWrite);
            } else {
                boolean booleanToWrite = produceCoilValue(writeRequestItem.getValues());
                modbusRequest = new WriteSingleCoilRequest(coilModbusAddress.getAddress(), booleanToWrite);
            }
        } else if (address instanceof MaskWriteRegisterModbusAddress) {
            MaskWriteRegisterModbusAddress maskWriteRegisterModbusAddress = (MaskWriteRegisterModbusAddress) address;
            if (quantity > 1) {
                throw new PlcProtocolException("Mask write request can only write one value");
            } else {
                // TODO: this should be better part of the payload not the addressing.
                int andMask = maskWriteRegisterModbusAddress.getAndMask();
                int orMask = maskWriteRegisterModbusAddress.getOrMask();
                modbusRequest = new MaskWriteRegisterRequest(maskWriteRegisterModbusAddress.getAddress(), andMask, orMask);
            }
        } else {
            throw new PlcProtocolException("Unsupported address type " + address.getClass() + " for a write request.");
        }
        short transactionId = (short) this.transactionId.getAndIncrement();
        requestsMap.put(transactionId, msg);
        out.add(new ModbusTcpPayload(transactionId, unitId, modbusRequest));
    }

    private void encodeReadRequest(PlcRequestContainer<PlcRequest, PlcResponse> msg, List<Object> out) throws PlcException {
        PlcReadRequest request = (PlcReadRequest) msg.getRequest();
        // TODO: support multiple requests
        ReadRequestItem<?> readRequestItem = request.getRequestItem().orElseThrow(() -> new PlcNotImplementedException("Only single message supported for now"));
        // TODO: check if we can map like this. Implication is that we can only work with int, short, byte and boolean
        // TODO: for higher datatypes float, double etc we might need to split the bytes into chunks
        int quantity = readRequestItem.getSize();
        // TODO: the unit the should be used for multiple Requests
        short unitId = 0;

        ModbusAddress address = (ModbusAddress) readRequestItem.getAddress();
        ModbusPdu modbusRequest;
        if (address instanceof CoilModbusAddress) {
            CoilModbusAddress coilModbusAddress = (CoilModbusAddress) address;
            modbusRequest = new ReadCoilsRequest(coilModbusAddress.getAddress(), quantity);
        } else if (address instanceof RegisterModbusAddress) {
            RegisterModbusAddress registerModbusAddress = (RegisterModbusAddress) address;
            modbusRequest = new ReadHoldingRegistersRequest(registerModbusAddress.getAddress(), quantity);
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
            throw new PlcProtocolException("Unsupported address type " + address.getClass() + " for a read request.");
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
            throw new PlcProtocolException("Unrelated payload received. [transactionId: " + msg.getTransactionId() + ", unitId: " + msg.getUnitId() + ", modbusPdu: " + msg.getModbusPdu() + "]");
        }

        // TODO: only single Item supported for now
        PlcRequest<?> request = plcRequestContainer.getRequest();
        RequestItem requestItem = request.getRequestItem().orElseThrow(() -> new PlcNotImplementedException("Only single message supported for now"));
        Class<?> dataType = requestItem.getDatatype();

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
            List<?> data = produceCoilValueList(requestItem, dataType, byteBuf);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem) requestItem, ResponseCode.OK, data)));
        } else if (modbusPdu instanceof ReadDiscreteInputsResponse) {
            // TODO: finish implementation
            ReadDiscreteInputsResponse readDiscreteInputsResponse = (ReadDiscreteInputsResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readDiscreteInputsResponse);
            ByteBuf byteBuf = readDiscreteInputsResponse.getInputStatus();
            List<?> data = produceCoilValueList(requestItem, dataType, byteBuf);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem) requestItem, ResponseCode.OK, data)));
        } else if (modbusPdu instanceof ReadHoldingRegistersResponse) {
            // TODO: finish implementation
            ReadHoldingRegistersResponse readHoldingRegistersResponse = (ReadHoldingRegistersResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readHoldingRegistersResponse);
            ByteBuf byteBuf = readHoldingRegistersResponse.getRegisters();
            // TODO: use register method
            List<?> data = produceRegisterValueList(requestItem, dataType, byteBuf);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem) requestItem, ResponseCode.OK, data)));
        } else if (modbusPdu instanceof ReadInputRegistersResponse) {
            // TODO: finish implementation
            ReadInputRegistersResponse readInputRegistersResponse = (ReadInputRegistersResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readInputRegistersResponse);
            ByteBuf byteBuf = readInputRegistersResponse.getRegisters();
            // TODO: use register method
            List<?> data = produceRegisterValueList(requestItem, dataType, byteBuf);
            plcRequestContainer.getResponseFuture().complete(new PlcReadResponse((PlcReadRequest) request, new ReadResponseItem((ReadRequestItem) requestItem, ResponseCode.OK, data)));
        } else if (modbusPdu instanceof MaskWriteRegisterResponse) {
            // TODO: finish implementation
            MaskWriteRegisterResponse maskWriteRegisterResponse = (MaskWriteRegisterResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", maskWriteRegisterResponse);
            plcRequestContainer.getResponseFuture().complete(new PlcWriteResponse((PlcWriteRequest) request, new WriteResponseItem<>((WriteRequestItem) requestItem, ResponseCode.OK)));
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

    private boolean produceCoilValue(List<?> values) throws PlcProtocolException {
        if (values.size() != 1) {
            throw new PlcProtocolException("Only one value allowed");
        }
        byte multiCoil = produceCoilValues(values)[0];
        return multiCoil != 0;
    }

    private byte[] produceCoilValues(List<?> values) throws PlcProtocolException {
        List<Byte> coils = new LinkedList<>();
        byte actualCoil = 0;
        int i = 7;
        for (Object value : values) {
            final boolean coilSet;
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
                    throw new PlcProtocolException("Exactly one byte is allowed: " + bytes.length);
                }
                byte byteValue = bytes[0];
                if (byteValue > 1) {
                    throw new PlcProtocolException("Value to high to fit into Byte: " + value);
                }
                coilSet = byteValue == 1;
            } else if (value.getClass() == Byte[].class) {
                Byte[] bytes = (Byte[]) value;
                if (bytes.length != 1) {
                    throw new PlcProtocolException("Exactly one byte is allowed: " + bytes.length);
                }
                byte byteValue = bytes[0];
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
            } else if (value.getClass() == BigInteger.class) {
                coilSet = value.equals(BigInteger.ONE);
            } else if (value.getClass() == Float.class) {
                coilSet = value.equals(1.0f);
            } else if (value.getClass() == Double.class) {
                coilSet = value.equals(1.0d);
            } else {
                throw new PlcUnsupportedDataTypeException(value.getClass());
            }
            byte coilToSet = coilSet ? (byte) 1 : (byte) 0;
            actualCoil = (byte) (actualCoil & 0xff | coilToSet << i);
            i--;
            if (i < 0) {
                coils.add(actualCoil);
                actualCoil = 0;
                i = 8;
            }
        }
        if (coils.isEmpty()) {
            // We only have one coil
            return new byte[]{actualCoil};
        }
        return ArrayUtils.toPrimitive(coils.toArray(new Byte[0]));
    }

    private byte[] produceRegisterValue(List<?> values) throws PlcProtocolException {
        ByteBuf buffer = Unpooled.buffer();
        long upperRegisterValue = 0xFFFFL;
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
                    throw new PlcProtocolException("Exactly two bytes are allowed: " + bytes.length);
                }
                buffer.writeBytes(bytes);
            } else if (value.getClass() == Byte[].class) {
                Byte[] bytes = (Byte[]) value;
                if (bytes.length != 2) {
                    throw new PlcProtocolException("Exactly two bytes are allowed: " + bytes.length);
                }
                buffer.writeBytes(ArrayUtils.toPrimitive(bytes));
            } else if (value.getClass() == Short.class) {
                if ((short) value < 0) {
                    throw new PlcProtocolException("Only positive values are supported for Short: " + value);
                }
                buffer.writeShort((short) value);
            } else if (value.getClass() == Integer.class) {
                if ((int) value > upperRegisterValue) {
                    throw new PlcProtocolException("Value to high to fit into register for Integer: " + value);
                }
                if ((int) value < 0) {
                    throw new PlcProtocolException("Only positive values are supported for Integer: " + value);
                }
                buffer.writeShort((int) value);
            } else if (value.getClass() == BigInteger.class) {
                if (((BigInteger) value).compareTo(BigInteger.ZERO) < 0) {
                    throw new PlcProtocolException("Only positive values are supported for BigInteger: " + value);
                }
                if (((BigInteger) value).compareTo(BigInteger.valueOf(0XFFFF_FFFFL)) > 0) {
                    throw new PlcProtocolException("Value to high to fit into register for BigInteger: " + value);
                }
                // TODO: for now we can't support big values as we only write one register at once
                if (((BigInteger) value).compareTo(BigInteger.valueOf(upperRegisterValue)) > 0) {
                    throw new PlcProtocolException("Value to high to fit into register for BigInteger: " + value);
                }
                // TODO: Register has 2 bytes so we trim to 2 instead of 4 like the second if above
                int maxBytes = 2;
                byte[] bigIntegerBytes = ((BigInteger) value).toByteArray();
                byte[] bytes = new byte[maxBytes];
                int lengthToCopy = Math.min(bigIntegerBytes.length, maxBytes);
                int srcPosition = Math.max(bigIntegerBytes.length - maxBytes, 0);
                int destPosition = maxBytes - lengthToCopy;
                System.arraycopy(bigIntegerBytes, srcPosition, bytes, destPosition, lengthToCopy);

                // TODO: check if this is a good representation.
                // TODO: can a big integer span multiple registers?
                buffer.writeBytes(bytes);
            } else if (value.getClass() == Float.class) {
                if (((float) value) < 0) {
                    throw new PlcProtocolException("Only positive values are supported for Float: " + value);
                }
                if (((float) value) > upperRegisterValue) {
                    throw new PlcProtocolException("Value to high to fit into register for Float: " + value);
                }
                buffer.writeShort(Math.round((float) value));
            } else if (value.getClass() == Double.class) {
                if (((double) value) < 0) {
                    throw new PlcProtocolException("Only positive values are supported for Double: " + value);
                }
                if (((double) value) > upperRegisterValue) {
                    throw new PlcProtocolException("Value to high to fit into register for Double: " + value);
                }
                buffer.writeShort((int) Math.round((double) value));
            } else {
                throw new PlcUnsupportedDataTypeException(value.getClass());
            }
        }
        byte[] result = new byte[buffer.writerIndex()];
        buffer.readBytes(result);
        return result;
    }

    ////////////////////////////////////////////////////////////////////////////////
    // Decoding helpers.
    ////////////////////////////////////////////////////////////////////////////////
    private <T> List<T> produceCoilValueList(RequestItem requestItem, Class<T> dataType, ByteBuf byteBuf) {
        ReadRequestItem readRequestItem = (ReadRequestItem) requestItem;
        byte[] bytes = new byte[byteBuf.readableBytes()];
        if (bytes.length < 1) {
            return Collections.emptyList();
        }
        byteBuf.readBytes(bytes);
        List<T> data = new LinkedList<>();
        int bitIndex = 0;
        int coilIndex = 0;
        while (data.size() < readRequestItem.getSize()) {
            if (bitIndex > 7) {
                // Every 8 Coils we need to increase the access
                coilIndex++;
                bitIndex = 0;
            }
            boolean coilSet = (bytes[coilIndex] & 0xff & (1L << bitIndex)) != 0;
            byte coilFlag = coilSet ? (byte) 1 : (byte) 0;
            if (dataType == Boolean.class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Boolean.valueOf(coilSet);
                data.add(itemToBeAdded);
            } else if (dataType == Byte.class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Byte.valueOf(coilFlag);
                data.add(itemToBeAdded);
            } else if (dataType == byte[].class) {
                data.add((T) new byte[]{coilFlag});
            } else if (dataType == Byte[].class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) new Byte[]{coilFlag};
                data.add(itemToBeAdded);
            } else if (dataType == Short.class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Short.valueOf(coilFlag);
                data.add(itemToBeAdded);
            } else if (dataType == Integer.class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Integer.valueOf(coilFlag);
                data.add(itemToBeAdded);
            } else if (dataType == BigInteger.class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) BigInteger.valueOf(coilFlag);
                data.add(itemToBeAdded);
            } else if (dataType == Float.class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Float.valueOf(coilFlag);
                data.add(itemToBeAdded);
            } else if (dataType == Double.class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Double.valueOf(coilFlag);
                data.add(itemToBeAdded);
            } else {
                throw new PlcUnsupportedDataTypeException(dataType);
            }
            bitIndex++;
        }
        return data;
    }

    private <T> List<T> produceRegisterValueList(RequestItem requestItem, Class<T> dataType, ByteBuf byteBuf) throws PlcProtocolException {
        ReadRequestItem readRequestItem = (ReadRequestItem) requestItem;
        int readableBytes = byteBuf.readableBytes();
        if (readableBytes % 2 != 0) {
            throw new PlcProtocolException("Readables bytes should even: " + readableBytes);
        }
        List<T> data = new LinkedList<>();
        for (int i = 0; i < readRequestItem.getSize(); i++) {
            byte[] register = new byte[2];
            byteBuf.readBytes(register);
            int intValue = register[0] << 8 | register[1] & 0xff;
            if (dataType == Boolean.class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Boolean.valueOf(intValue == 1);
                data.add(itemToBeAdded);
            } else if (dataType == Byte.class) {
                if (intValue > Byte.MAX_VALUE) {
                    throw new PlcProtocolException("Value to high to fit into Byte: " + intValue);
                }
                @SuppressWarnings("unchecked")
                T itemToBeADded = (T) Byte.valueOf((byte) intValue);
                data.add(itemToBeADded);
            } else if (dataType == byte[].class) {
                T itemToBeAdded = (T) register;
                data.add(itemToBeAdded);
            } else if (dataType == Byte[].class) {
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) ArrayUtils.toObject(register);
                data.add(itemToBeAdded);
            } else if (dataType == Short.class) {
                if (intValue > Short.MAX_VALUE) {
                    throw new PlcProtocolException("Value to high to fit into Short: " + intValue);
                }
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Short.valueOf((short) intValue);
                data.add(itemToBeAdded);
            } else if (dataType == Integer.class) {
                if (intValue < 0) {
                    throw new PlcProtocolException("Integer underflow: " + intValue);
                }
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) Integer.valueOf(intValue);
                data.add(itemToBeAdded);
            } else if (dataType == BigInteger.class) {
                if (intValue < 0) {
                    throw new PlcProtocolException("BigInteger underflow: " + intValue);
                }
                // TODO: can a big integer span multiple registers?
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) new BigInteger(register);
                data.add(itemToBeAdded);
            } else if (dataType == Float.class) {
                if (intValue < 0) {
                    throw new PlcProtocolException("BigInteger underflow: " + intValue);
                }
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) new Float(intValue);
                data.add(itemToBeAdded);
            } else if (dataType == Double.class) {
                if (intValue < 0) {
                    throw new PlcProtocolException("BigInteger underflow: " + intValue);
                }
                @SuppressWarnings("unchecked")
                T itemToBeAdded = (T) new Double(intValue);
                data.add(itemToBeAdded);
            } else {
                throw new PlcUnsupportedDataTypeException(dataType);
            }
        }
        return data;
    }
}
