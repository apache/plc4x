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
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultBooleanFieldItem;
import org.apache.plc4x.java.modbus.messages.items.DefaultModbusByteArrayFieldItem;
import org.apache.plc4x.java.modbus.model.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicInteger;


public class Plc4XModbusProtocol extends MessageToMessageCodec<ModbusTcpPayload, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XModbusProtocol.class);

    public final AtomicInteger transactionId = new AtomicInteger();

    private final ConcurrentMap<Short, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>> requestsMap = new ConcurrentHashMap<>();

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) throws Exception {
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

    private void encodeWriteRequest(PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) throws PlcException {
        InternalPlcWriteRequest request = (InternalPlcWriteRequest) msg.getRequest();

        // TODO: support multiple requests
        if (request.getFieldNames().size() != 1) {
            throw new PlcNotImplementedException("Only single message supported for now");
        }
        // TODO: check if we can map like this. Implication is that we can only work with int, short, byte and boolean
        // TODO: for higher data types float, double etc we might need to split the bytes into chunks
        String fieldName = request.getFieldNames().iterator().next();
        int quantity = request.getNumberOfValues(fieldName);
        ModbusField field = (ModbusField) request.getField(fieldName);
        if (quantity != field.getQuantity()) {
            LOGGER.warn("Supplied number of values [{}] don't match t the addressed quantity of [{}]", field.getQuantity(), quantity);
        }

        short unitId = 0;

        /*
         * It seems that in Modbus, there are only two types of resources, that can be accessed:
         * - Register: 2 byte value
         * - Coil: 1 bit value
         *
         * Registers:
         * When writing a bit, byte or char (a one byte or less data types) into a register the value is filled up to
         * fit the 2 bytes.
         * When writing a data type that has more than 2 bytes, subsequent registers are written to automatically.
         *
         * Coils:
         * When transferring data from/to a coil, 8 coil values can be transferred in a single byte.
         * Naturally a coil is a boolean data type, however, similar to the registers bigger data types can be
         * transported, by setting multiple subsequent coils: 32bit int -> 32 coils
         *
         * In all cases where we are accessing more than the natural size of the datatype, we have to keep this in mind
         * when addressing them. So if reading 2 32bit integers, this is split up into four registers. So for the second
         * int we have to increment the address accordingly.
         */
        ModbusPdu modbusRequest;
        if (field instanceof RegisterModbusField) {
            RegisterModbusField registerModbusField = (RegisterModbusField) field;
            if (quantity > 1) {
                byte[] bytesToWrite = produceRegisterValue(Arrays.asList(request.getFieldItem(fieldName).getValues()));
                // A register is a 16 bit (2 byte) value ... so every value needs 2 byte.
                int requiredLength = 2 * quantity;
                if (bytesToWrite.length != requiredLength) {
                    throw new PlcProtocolException("Invalid register values created. Should be at least quantity * 2 = N bytes. Was " + bytesToWrite.length + ", expected " + requiredLength);
                }
                modbusRequest = new WriteMultipleRegistersRequest(registerModbusField.getAddress(), quantity, bytesToWrite);
            } else {
                byte[] register = produceRegisterValue(Arrays.asList(request.getFieldItem(fieldName).getValues()));
                if ((register == null) || (register.length != 2)) {
                    throw new PlcProtocolException("Invalid register values created. Should be 2 bytes. Was " +
                        ((register != null) ? register.length : 0));
                }
                // Reconvert the two bytes back to an int.
                int intToWrite = register[0] << 8 | register[1] & 0xff;
                modbusRequest = new WriteSingleRegisterRequest(registerModbusField.getAddress(), intToWrite);
            }
        } else if (field instanceof CoilModbusField) {
            CoilModbusField coilModbusField = (CoilModbusField) field;
            if (quantity > 1) {
                byte[] bytesToWrite = produceCoilValues(Arrays.asList(request.getFieldItem(fieldName).getValues()));
                // As each coil value represents a bit, the number of bytes needed
                // equals "ceil(quantity/8)" (a 3 bit shift is a division by 8 ... the +1 is the "ceil")
                int requiredLength = (quantity >> 3) + 1;
                if (bytesToWrite.length != requiredLength) {
                    throw new PlcProtocolException(
                        "Invalid coil values created. Should be big enough to transport N bits. Was " +
                            bytesToWrite.length + ", expected " + requiredLength);
                }
                modbusRequest = new WriteMultipleCoilsRequest(coilModbusField.getAddress(), quantity, bytesToWrite);
            } else {
                boolean booleanToWrite = produceCoilValue(Arrays.asList(request.getFieldItem(fieldName).getValues()));
                modbusRequest = new WriteSingleCoilRequest(coilModbusField.getAddress(), booleanToWrite);
            }
        } else if (field instanceof MaskWriteRegisterModbusField) {
            MaskWriteRegisterModbusField maskWriteRegisterModbusField = (MaskWriteRegisterModbusField) field;
            if (quantity > 1) {
                throw new PlcProtocolException("Mask write request can only write one value");
            } else {
                // TODO: this should be better part of the payload not the addressing.
                int andMask = maskWriteRegisterModbusField.getAndMask();
                int orMask = maskWriteRegisterModbusField.getOrMask();
                modbusRequest = new MaskWriteRegisterRequest(maskWriteRegisterModbusField.getAddress(), andMask, orMask);
            }
        } else {
            throw new PlcProtocolException("Unsupported field type " + field.getClass() + " for a write request.");
        }
        short transactionId = (short) this.transactionId.getAndIncrement();
        requestsMap.put(transactionId, msg);
        out.add(new ModbusTcpPayload(transactionId, unitId, modbusRequest));
    }

    private void encodeReadRequest(PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) throws PlcException {
        PlcReadRequest request = (PlcReadRequest) msg.getRequest();
        // TODO: support multiple requests
        if (request.getFieldNames().size() != 1) {
            throw new PlcNotImplementedException("Only single message supported for now");
        }
        // TODO: check if we can map like this. Implication is that we can only work with int, short, byte and boolean
        // TODO: for higher data types float, double etc we might need to split the bytes into chunks
        String fieldName = request.getFieldNames().iterator().next();

        ModbusField field = (ModbusField) request.getField(fieldName);
        int quantity = field.getQuantity();
        // TODO: the unit the should be used for multiple Requests
        short unitId = 0;

        ModbusPdu modbusRequest;
        if (field instanceof CoilModbusField) {
            CoilModbusField coilModbusField = (CoilModbusField) field;
            modbusRequest = new ReadCoilsRequest(coilModbusField.getAddress(), quantity);
        } else if (field instanceof RegisterModbusField) {
            RegisterModbusField registerModbusField = (RegisterModbusField) field;
            modbusRequest = new ReadHoldingRegistersRequest(registerModbusField.getAddress(), quantity);
        } else if (field instanceof ReadDiscreteInputsModbusField) {
            ReadDiscreteInputsModbusField readDiscreteInputsModbusField = (ReadDiscreteInputsModbusField) field;
            modbusRequest = new ReadDiscreteInputsRequest(readDiscreteInputsModbusField.getAddress(), quantity);
        } else if (field instanceof ReadHoldingRegistersModbusField) {
            ReadHoldingRegistersModbusField readHoldingRegistersModbusField = (ReadHoldingRegistersModbusField) field;
            modbusRequest = new ReadHoldingRegistersRequest(readHoldingRegistersModbusField.getAddress(), quantity);
        } else if (field instanceof ReadInputRegistersModbusField) {
            ReadInputRegistersModbusField readInputRegistersModbusField = (ReadInputRegistersModbusField) field;
            modbusRequest = new ReadInputRegistersRequest(readInputRegistersModbusField.getAddress(), quantity);
        } else {
            throw new PlcProtocolException("Unsupported field type " + field.getClass() + " for a read request.");
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
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> plcRequestContainer = requestsMap.get(transactionId);
        if (plcRequestContainer == null) {
            throw new PlcProtocolException("Unrelated payload received. [transactionId: " + msg.getTransactionId() + ", unitId: " + msg.getUnitId() + ", modbusPdu: " + msg.getModbusPdu() + "]");
        }

        // TODO: only single Item supported for now
        InternalPlcFieldRequest request = (InternalPlcFieldRequest) plcRequestContainer.getRequest();
        // TODO: support multiple requests (Shouldn't be needed as the request wouldn't have been sent)
        if (request.getFieldNames().size() != 1) {
            throw new PlcNotImplementedException("Only single message supported for now");
        }
        String fieldName = request.getFieldNames().iterator().next();
        ModbusField field = (ModbusField) request.getField(fieldName);

        ModbusPdu modbusPdu = msg.getModbusPdu();
        //short unitId = msg.getUnitId();

        if (modbusPdu instanceof WriteMultipleCoilsResponse) {
            // TODO: finish implementation
            WriteMultipleCoilsResponse writeMultipleCoilsResponse = (WriteMultipleCoilsResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, quantity:{}", writeMultipleCoilsResponse, writeMultipleCoilsResponse.getAddress(), writeMultipleCoilsResponse.getQuantity());
            Map<String, PlcResponseCode> responseValues = new HashMap<>();
            responseValues.put(fieldName, PlcResponseCode.OK);
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcWriteResponse((InternalPlcWriteRequest) request, responseValues));
        } else if (modbusPdu instanceof WriteMultipleRegistersResponse) {
            // TODO: finish implementation
            WriteMultipleRegistersResponse writeMultipleRegistersResponse = (WriteMultipleRegistersResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, quantity:{}", writeMultipleRegistersResponse, writeMultipleRegistersResponse.getAddress(), writeMultipleRegistersResponse.getQuantity());
            Map<String, PlcResponseCode> responseValues = new HashMap<>();
            responseValues.put(fieldName, PlcResponseCode.OK);
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcWriteResponse((InternalPlcWriteRequest) request, responseValues));
        } else if (modbusPdu instanceof WriteSingleCoilResponse) {
            // TODO: finish implementation
            WriteSingleCoilResponse writeSingleCoilResponse = (WriteSingleCoilResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, value:{}", writeSingleCoilResponse, writeSingleCoilResponse.getAddress(), writeSingleCoilResponse.getValue());
            Map<String, PlcResponseCode> responseValues = new HashMap<>();
            responseValues.put(fieldName, PlcResponseCode.OK);
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcWriteResponse((InternalPlcWriteRequest) request, responseValues));
        } else if (modbusPdu instanceof WriteSingleRegisterResponse) {
            // TODO: finish implementation
            WriteSingleRegisterResponse writeSingleRegisterResponse = (WriteSingleRegisterResponse) modbusPdu;
            LOGGER.debug("{}: address:{}, value:{}", writeSingleRegisterResponse, writeSingleRegisterResponse.getAddress(), writeSingleRegisterResponse.getValue());
            Map<String, PlcResponseCode> responseValues = new HashMap<>();
            responseValues.put(fieldName, PlcResponseCode.OK);
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcWriteResponse((InternalPlcWriteRequest) request, responseValues));
        } else if (modbusPdu instanceof ReadCoilsResponse) {
            // TODO: finish implementation
            ReadCoilsResponse readCoilsResponse = (ReadCoilsResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readCoilsResponse);
            ByteBuf byteBuf = readCoilsResponse.getCoilStatus();
            DefaultBooleanFieldItem data = produceCoilValueList(byteBuf, field.getQuantity());
            Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> responseValues = new HashMap<>();
            responseValues.put(fieldName, new ImmutablePair<>(PlcResponseCode.OK, data));
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcReadResponse((InternalPlcReadRequest) request, responseValues));
        } else if (modbusPdu instanceof ReadDiscreteInputsResponse) {
            // TODO: finish implementation
            ReadDiscreteInputsResponse readDiscreteInputsResponse = (ReadDiscreteInputsResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readDiscreteInputsResponse);
            ByteBuf byteBuf = readDiscreteInputsResponse.getInputStatus();
            DefaultBooleanFieldItem data = produceCoilValueList(byteBuf, field.getQuantity());
            Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> responseValues = new HashMap<>();
            responseValues.put(fieldName, new ImmutablePair<>(PlcResponseCode.OK, data));
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcReadResponse((InternalPlcReadRequest) request, responseValues));
        } else if (modbusPdu instanceof ReadHoldingRegistersResponse) {
            // TODO: finish implementation
            ReadHoldingRegistersResponse readHoldingRegistersResponse = (ReadHoldingRegistersResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readHoldingRegistersResponse);
            ByteBuf byteBuf = readHoldingRegistersResponse.getRegisters();
            // TODO: use register method
            DefaultModbusByteArrayFieldItem data = produceRegisterValueList(byteBuf, field.getQuantity());
            Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> responseValues = new HashMap<>();
            responseValues.put(fieldName, new ImmutablePair<>(PlcResponseCode.OK, data));
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcReadResponse((InternalPlcReadRequest) request, responseValues));
        } else if (modbusPdu instanceof ReadInputRegistersResponse) {
            // TODO: finish implementation
            ReadInputRegistersResponse readInputRegistersResponse = (ReadInputRegistersResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", readInputRegistersResponse);
            ByteBuf byteBuf = readInputRegistersResponse.getRegisters();
            // TODO: use register method
            DefaultModbusByteArrayFieldItem data = produceRegisterValueList(byteBuf, field.getQuantity());
            Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> responseValues = new HashMap<>();
            responseValues.put(fieldName, new ImmutablePair<>(PlcResponseCode.OK, data));
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcReadResponse((InternalPlcReadRequest) request, responseValues));
        } else if (modbusPdu instanceof MaskWriteRegisterResponse) {
            // TODO: finish implementation
            MaskWriteRegisterResponse maskWriteRegisterResponse = (MaskWriteRegisterResponse) modbusPdu;
            LOGGER.debug("{}: Nothing", maskWriteRegisterResponse);
            Map<String, PlcResponseCode> responseValues = new HashMap<>();
            responseValues.put(fieldName, PlcResponseCode.OK);
            plcRequestContainer.getResponseFuture().complete(new DefaultPlcWriteResponse((InternalPlcWriteRequest) request, responseValues));
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
    private DefaultBooleanFieldItem produceCoilValueList(ByteBuf byteBuf, int expectedQuantity) {
        if (byteBuf.readableBytes() < expectedQuantity / 8) {
            LOGGER.warn("Expected to read {} coils but only max of {} can be supplied", expectedQuantity, byteBuf.readableBytes() * 8);
        }
        byte[] bytes = new byte[byteBuf.readableBytes()];
        if (bytes.length < 1) {
            return new DefaultBooleanFieldItem();
        }
        byteBuf.readBytes(bytes);
        List<Boolean> data = new LinkedList<>();
        int bitIndex = 0;
        int coilIndex = 0;
        while (coilIndex < bytes.length && data.size() < expectedQuantity) {
            if (bitIndex > 7) {
                // Every 8 Coils we need to increase the access
                coilIndex++;
                bitIndex = 0;
                if (coilIndex >= bytes.length) {
                    break;
                }
            }
            boolean coilSet = (bytes[coilIndex] & 0xff & (1L << bitIndex)) != 0;
            data.add(coilSet);
            bitIndex++;
        }
        return new DefaultBooleanFieldItem(data.toArray(new Boolean[0]));
    }

    private DefaultModbusByteArrayFieldItem produceRegisterValueList(ByteBuf byteBuf, int expectedQuantity) throws PlcProtocolException {
        int readableBytes = byteBuf.readableBytes();
        if (readableBytes % 2 != 0) {
            throw new PlcProtocolException("Readables bytes should even: " + readableBytes);
        }
        List<Byte[]> data = new LinkedList<>();
        while (byteBuf.readableBytes() > 0) {
            byte[] register = new byte[2];
            byteBuf.readBytes(register);
            data.add(ArrayUtils.toObject(register));
        }
        return new DefaultModbusByteArrayFieldItem(data.toArray(new Byte[0][0]));
    }
}
