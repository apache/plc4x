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
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.items.PlcReadResponseItem;
import org.apache.plc4x.java.api.messages.items.PlcWriteResponseItem;
import org.apache.plc4x.java.api.messages.items.ResponseItem;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.base.protocol.Plc4XSupportedDataTypes;
import org.apache.plc4x.java.modbus.model.*;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.math.BigInteger;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.apache.plc4x.java.base.protocol.Plc4XSupportedDataTypes.defaultAssert;
import static org.apache.plc4x.java.base.protocol.Plc4XSupportedDataTypes.streamOfBigEndianDataTypePairs;
import static org.apache.plc4x.java.base.util.Assert.assertByteEquals;
import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThat;
import static org.junit.Assume.assumeThat;

@RunWith(Parameterized.class)
public class Plc4XModbusProtocolTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XModbusProtocolTest.class);

    private Plc4XModbusProtocol SUT;

    @Rule
    public ExpectedException expectedException = ExpectedException.none();
    // TODO: implement these types
    private List<String> notYetSupportedDataType = Stream.of(
        GregorianCalendar.class,
        String.class
    ).map(Class::getSimpleName).collect(Collectors.toList());

    @Parameterized.Parameter
    public String payloadClazzName;

    @Parameterized.Parameter(1)
    public PlcRequestContainer<PlcRequest, PlcResponse> plcRequestContainer;

    @Parameterized.Parameter(2)
    public CompletableFuture completableFuture;

    @Parameterized.Parameter(3)
    public String plcRequestContainerClassName;

    @Parameterized.Parameter(4)
    public ModbusTcpPayload modbusTcpPayload;

    @Parameterized.Parameter(5)
    public String modbusPduName;

    @Parameterized.Parameters(name = "{index} Type:{0} {3} {5}")
    public static Collection<Object[]> data() {
        return streamOfBigEndianDataTypePairs()
            // Float and Double getting truncated in modbus.
            .map(dataTypePair -> {
                if (dataTypePair.getDataTypeClass() == Float.class) {
                    return Plc4XSupportedDataTypes.DataTypePair.of(1f, new byte[]{0x0, 0x1});
                } else if (dataTypePair.getDataTypeClass() == Double.class) {
                    return Plc4XSupportedDataTypes.DataTypePair.of(1d, new byte[]{0x0, 0x1});
                } else {
                    return dataTypePair;
                }
            })
            .map(dataTypePair -> Stream.of(
                producePair(dataTypePair.getDataTypeClass(), st(() -> CoilModbusField.of("coil:1")), new ReadCoilsResponse(Unpooled.wrappedBuffer(new byte[]{(byte) 0x1}))),
                producePair(st(() -> CoilModbusField.of("coil:1")), new WriteSingleCoilResponse(1, 1), mapDataTypePairForCoil(dataTypePair).getValue()),
                /* Read request no supported on maskwrite so how to handle?
                producePair(pair.getDataTypeClass(), MaskWriteRegisterModbusField.of("maskwrite:1/1/2"), new MaskWriteRegisterResponse(1, 1, 2)),
                */
                producePair(st(() -> MaskWriteRegisterModbusField.of("maskwrite:1/1/2")), new MaskWriteRegisterResponse(1, 1, 2), mapDataTypePairForRegister(dataTypePair).getValue()),
                producePair(dataTypePair.getDataTypeClass(), st(() -> ReadDiscreteInputsModbusField.of("readdiscreteinputs:1")), new ReadDiscreteInputsResponse(Unpooled.wrappedBuffer(new byte[]{(byte) 0x01}))),
                producePair(dataTypePair.getDataTypeClass(), st(() -> ReadHoldingRegistersModbusField.of("readholdingregisters:1")), new ReadHoldingRegistersResponse(Unpooled.wrappedBuffer(cutRegister(mapDataTypePairForRegister(dataTypePair).getByteRepresentation())))),
                producePair(dataTypePair.getDataTypeClass(), st(() -> ReadInputRegistersModbusField.of("readinputregisters:1")), new ReadInputRegistersResponse(Unpooled.wrappedBuffer(cutRegister(mapDataTypePairForRegister(dataTypePair).getByteRepresentation())))),
                producePair(st(() -> CoilModbusField.of("coil:1")), new WriteMultipleCoilsResponse(1, 3), mapDataTypePairForCoil(dataTypePair).getValue(), mapDataTypePairForCoil(dataTypePair).getValue(), mapDataTypePairForCoil(dataTypePair).getValue()),
                producePair(st(() -> RegisterModbusField.of("register:1")), new WriteMultipleRegistersResponse(1, 3), mapDataTypePairForRegister(dataTypePair).getValue(), mapDataTypePairForRegister(dataTypePair).getValue(), mapDataTypePairForRegister(dataTypePair).getValue()),
                producePair(st(() -> RegisterModbusField.of("register:1")), new WriteSingleRegisterResponse(1, cutRegister(mapDataTypePairForRegister(dataTypePair).getByteRepresentation())[0]), mapDataTypePairForRegister(dataTypePair).getValue())
            ))
            .flatMap(stream -> stream)
            .map(pair -> new Object[]{pair.left.getRequest().getRequestItem().orElseThrow(IllegalStateException::new).getDatatype().getSimpleName(), pair.left, pair.left.getResponseFuture(), pair.left.getRequest().getClass().getSimpleName(), pair.right, pair.right.getModbusPdu().getClass().getSimpleName()}).collect(Collectors.toList());
    }

    @FunctionalInterface
    public interface FieldSupplier<T> {
        T getField() throws PlcInvalidFieldException;
    }

    public static <T> T st(FieldSupplier<T> supplier) {
        try {
            return supplier.getField();
        } catch (PlcInvalidFieldException e) {
            throw new RuntimeException(e);
        }
    }

    private static ImmutablePair<PlcRequestContainer<PlcReadRequest, PlcResponse>, ModbusTcpPayload> producePair(Class type, PlcField field, ModbusPdu modbusPdu) {
        return ImmutablePair.of(
            new PlcRequestContainer<>(
                PlcReadRequest
                    .builder()
                    .addItem(type, field)
                    .build(), new CompletableFuture<>()),
            new ModbusTcpPayload((short) 0, (short) 0, modbusPdu)
        );
    }

    @SuppressWarnings("unchecked")
    private static ImmutablePair<PlcRequestContainer<PlcWriteRequest, PlcResponse>, ModbusTcpPayload> producePair(PlcField field, ModbusPdu modbusPdu, Object... values) {
        if (values.length == 0) {
            throw new IllegalArgumentException("At least one value ist required");
        }
        if (values.length == 1) {
            return ImmutablePair.of(
                new PlcRequestContainer<>(
                    PlcWriteRequest
                        .builder()
                        .addItem(field, values[0])
                        .build(), new CompletableFuture<>()),
                new ModbusTcpPayload((short) 0, (short) 0, modbusPdu)
            );
        } else {
            return ImmutablePair.of(
                new PlcRequestContainer<>(
                    PlcWriteRequest
                        .builder()
                        .addItem((Class<Object>) values[0].getClass(), field, values)
                        .build(), new CompletableFuture<>()),
                new ModbusTcpPayload((short) 0, (short) 0, modbusPdu)
            );
        }
    }

    private static byte[] cutRegister(byte[] right) {
        return new byte[]{right.length > 1 ? right[right.length - 2] : 0x0, right[right.length - 1]};
    }

    @Before
    public void setUp() {
        SUT = new Plc4XModbusProtocol();
    }

    @Test
    public void encode() throws Exception {
        assumeThat(payloadClazzName + " not yet implemented", notYetSupportedDataType, not(hasItem(payloadClazzName)));
        ArrayList<Object> out = new ArrayList<>();
        SUT.encode(null, plcRequestContainer, out);
        assertThat(out, hasSize(1));
        assertThat(out.get(0), instanceOf(ModbusTcpPayload.class));
        ModbusTcpPayload modbusTcpPayload = (ModbusTcpPayload) out.get(0);
        ModbusPdu modbusPdu = modbusTcpPayload.getModbusPdu();
        if (modbusPdu instanceof MaskWriteRegisterRequest) {
            MaskWriteRegisterRequest maskWriteRegisterRequest = (MaskWriteRegisterRequest) modbusPdu;
            int address = maskWriteRegisterRequest.getAddress();
            assertThat(address, equalTo(1));
            int andMask = maskWriteRegisterRequest.getAndMask();
            int orMask = maskWriteRegisterRequest.getOrMask();
            if (payloadClazzName.equals(Boolean.class.getSimpleName())) {
                assertThat(andMask, equalTo(1));
                assertThat(orMask, equalTo(2));
            } else if (payloadClazzName.equals(Byte.class.getSimpleName())) {
                assertThat(andMask, equalTo(1));
                assertThat(orMask, equalTo(2));
            } else if (payloadClazzName.equals(Short.class.getSimpleName())) {
                assertThat(andMask, equalTo(1));
                assertThat(orMask, equalTo(2));
            } else if (payloadClazzName.equals(GregorianCalendar.class.getSimpleName())) {
                assertThat(andMask, equalTo(1));
                assertThat(orMask, equalTo(2));
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertThat(andMask, equalTo(1));
                assertThat(orMask, equalTo(2));
            } else if (payloadClazzName.equals(Integer.class.getSimpleName())) {
                assertThat(andMask, equalTo(1));
                assertThat(orMask, equalTo(2));
            } else if (payloadClazzName.equals(BigInteger.class.getSimpleName())) {
                assertThat(andMask, equalTo(1));
                assertThat(orMask, equalTo(2));
            } else if (payloadClazzName.equals(String.class.getSimpleName())) {
                assertThat(andMask, equalTo(1));
                assertThat(orMask, equalTo(2));
            }
        } else if (modbusPdu instanceof ReadCoilsRequest) {
            ReadCoilsRequest readCoilsRequest = (ReadCoilsRequest) modbusPdu;
            int address = readCoilsRequest.getAddress();
            int quantity = readCoilsRequest.getQuantity();
            assertThat(address, equalTo(1));
            assertThat(quantity, equalTo(1));
        } else if (modbusPdu instanceof ReadDiscreteInputsRequest) {
            ReadDiscreteInputsRequest readCoilsRequest = (ReadDiscreteInputsRequest) modbusPdu;
            int address = readCoilsRequest.getAddress();
            int quantity = readCoilsRequest.getQuantity();
            assertThat(address, equalTo(1));
            assertThat(quantity, equalTo(1));
        } else if (modbusPdu instanceof ReadHoldingRegistersRequest) {
            ReadHoldingRegistersRequest readCoilsRequest = (ReadHoldingRegistersRequest) modbusPdu;
            int address = readCoilsRequest.getAddress();
            int quantity = readCoilsRequest.getQuantity();
            assertThat(address, equalTo(1));
            assertThat(quantity, equalTo(1));
        } else if (modbusPdu instanceof ReadInputRegistersRequest) {
            ReadInputRegistersRequest readCoilsRequest = (ReadInputRegistersRequest) modbusPdu;
            int address = readCoilsRequest.getAddress();
            int quantity = readCoilsRequest.getQuantity();
            assertThat(address, equalTo(1));
            assertThat(quantity, equalTo(1));
        } else if (modbusPdu instanceof WriteMultipleCoilsRequest) {
            WriteMultipleCoilsRequest writeMultipleCoilsRequest = (WriteMultipleCoilsRequest) modbusPdu;
            int address = writeMultipleCoilsRequest.getAddress();
            assertThat(address, equalTo(1));
            ByteBuf value = writeMultipleCoilsRequest.getValues();
            byte[] bytes = new byte[value.readableBytes()];
            value.readBytes(value);
            if (payloadClazzName.equals(Boolean.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0}, bytes);
            } else if (payloadClazzName.equals(Byte.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0}, bytes);
            } else if (payloadClazzName.equals(Short.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0}, bytes);
            } else if (payloadClazzName.equals(Calendar.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0}, bytes);
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0}, bytes);
            } else if (payloadClazzName.equals(Integer.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0}, bytes);
            } else if (payloadClazzName.equals(BigInteger.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0}, bytes);
            } else if (payloadClazzName.equals(String.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0}, bytes);
            }
        } else if (modbusPdu instanceof WriteMultipleRegistersRequest) {
            WriteMultipleRegistersRequest writeMultipleRegistersRequest = (WriteMultipleRegistersRequest) modbusPdu;
            int address = writeMultipleRegistersRequest.getAddress();
            assertThat(address, equalTo(1));
            ByteBuf value = writeMultipleRegistersRequest.getValues();
            byte[] bytes = new byte[value.readableBytes()];
            value.readBytes(value);
            if (payloadClazzName.equals(Boolean.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bytes);
            } else if (payloadClazzName.equals(Byte.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bytes);
            } else if (payloadClazzName.equals(Short.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bytes);
            } else if (payloadClazzName.equals(Calendar.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bytes);
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bytes);
            } else if (payloadClazzName.equals(Integer.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bytes);
            } else if (payloadClazzName.equals(BigInteger.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bytes);
            } else if (payloadClazzName.equals(String.class.getSimpleName())) {
                assertByteEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, bytes);
            }
        } else if (modbusPdu instanceof WriteSingleCoilRequest) {
            WriteSingleCoilRequest writeSingleCoilRequest = (WriteSingleCoilRequest) modbusPdu;
            int address = writeSingleCoilRequest.getAddress();
            assertThat(address, equalTo(1));
            int value = writeSingleCoilRequest.getValue();
            boolean coilValue = value == 0xFF00;
            if (payloadClazzName.equals(Boolean.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            } else if (payloadClazzName.equals(Byte.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            } else if (payloadClazzName.equals(Short.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            } else if (payloadClazzName.equals(Calendar.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            } else if (payloadClazzName.equals(Double.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            } else if (payloadClazzName.equals(Integer.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            } else if (payloadClazzName.equals(BigInteger.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            } else if (payloadClazzName.equals(String.class.getSimpleName())) {
                assertThat(coilValue, equalTo(true));
            }
        } else if (modbusPdu instanceof WriteSingleRegisterRequest) {
            WriteSingleRegisterRequest writeSingleRegisterRequest = (WriteSingleRegisterRequest) modbusPdu;
            int address = writeSingleRegisterRequest.getAddress();
            assertThat(address, equalTo(1));
            int value = writeSingleRegisterRequest.getValue();
            if (payloadClazzName.equals(Boolean.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            } else if (payloadClazzName.equals(Byte.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            } else if (payloadClazzName.equals(Short.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            } else if (payloadClazzName.equals(Calendar.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            } else if (payloadClazzName.equals(Float.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            } else if (payloadClazzName.equals(Double.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            } else if (payloadClazzName.equals(Integer.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            } else if (payloadClazzName.equals(BigInteger.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            } else if (payloadClazzName.equals(String.class.getSimpleName())) {
                assertThat(value, equalTo(1));
            }
        }
    }

    @Test
    public void decode() throws Exception {
        assumeThat(payloadClazzName + " not yet implemented", notYetSupportedDataType, not(hasItem(payloadClazzName)));
        ArrayList<Object> in = new ArrayList<>();
        SUT.encode(null, plcRequestContainer, in);
        assertThat(in, hasSize(1));
        syncInvoiceId();
        ArrayList<Object> out = new ArrayList<>();
        SUT.decode(null, modbusTcpPayload, out);
        assertThat(out, hasSize(0));
        LOGGER.info("PlcRequestContainer {}", plcRequestContainer);
        PlcResponse plcResponse = plcRequestContainer.getResponseFuture().get();
        ResponseItem responseItem = (ResponseItem) plcResponse.getResponseItem().get();
        LOGGER.info("ResponseItem {}", responseItem);
        ModbusPdu modbusPdu = modbusTcpPayload.getModbusPdu();
        if (modbusPdu instanceof MaskWriteRegisterResponse) {
            PlcWriteResponseItem writeResponseItem = (PlcWriteResponseItem) responseItem;
            assertEquals(PlcResponseCode.OK, writeResponseItem.getResponseCode());
        } else if (modbusPdu instanceof ReadCoilsResponse) {
            PlcReadResponseItem readResponseItem = (PlcReadResponseItem) responseItem;
            Object value = readResponseItem.getValues().get(0);
            defaultAssert(value, Plc4XModbusProtocolTest::mapDataTypePairForCoil);
        } else if (modbusPdu instanceof ReadDiscreteInputsResponse) {
            PlcReadResponseItem readResponseItem = (PlcReadResponseItem) responseItem;
            Object value = readResponseItem.getValues().get(0);
            defaultAssert(value, Plc4XModbusProtocolTest::mapDataTypePairForCoil);
        } else if (modbusPdu instanceof ReadHoldingRegistersResponse) {
            PlcReadResponseItem readResponseItem = (PlcReadResponseItem) responseItem;
            Object value = readResponseItem.getValues().get(0);
            defaultAssert(value, Plc4XModbusProtocolTest::mapDataTypePairForRegister);
        } else if (modbusPdu instanceof ReadInputRegistersResponse) {
            PlcReadResponseItem readResponseItem = (PlcReadResponseItem) responseItem;
            Object value = readResponseItem.getValues().get(0);
            defaultAssert(value, Plc4XModbusProtocolTest::mapDataTypePairForRegister);
        } else if (modbusPdu instanceof WriteMultipleCoilsResponse) {
            PlcWriteResponseItem writeResponseItem = (PlcWriteResponseItem) responseItem;
            assertEquals(PlcResponseCode.OK, writeResponseItem.getResponseCode());
        } else if (modbusPdu instanceof WriteMultipleRegistersResponse) {
            PlcWriteResponseItem writeResponseItem = (PlcWriteResponseItem) responseItem;
            assertEquals(PlcResponseCode.OK, writeResponseItem.getResponseCode());
        } else if (modbusPdu instanceof WriteSingleCoilResponse) {
            PlcWriteResponseItem writeResponseItem = (PlcWriteResponseItem) responseItem;
            assertEquals(PlcResponseCode.OK, writeResponseItem.getResponseCode());
        } else if (modbusPdu instanceof WriteSingleRegisterResponse) {
            PlcWriteResponseItem writeResponseItem = (PlcWriteResponseItem) responseItem;
            assertEquals(PlcResponseCode.OK, writeResponseItem.getResponseCode());
        }
    }

    private static Plc4XSupportedDataTypes.DataTypePair mapDataTypePairForCoil(Plc4XSupportedDataTypes.DataTypePair dataTypePair) {
        return mapDataTypePairForRegisterOrCoil(dataTypePair, 1);
    }

    private static Plc4XSupportedDataTypes.DataTypePair mapDataTypePairForRegister(Plc4XSupportedDataTypes.DataTypePair dataTypePair) {
        return mapDataTypePairForRegisterOrCoil(dataTypePair, 2);
    }

    private static Plc4XSupportedDataTypes.DataTypePair mapDataTypePairForRegisterOrCoil(Plc4XSupportedDataTypes.DataTypePair dataTypePair, int size) {
        byte[] byteRepresentation = dataTypePair.getByteRepresentation();
        if (dataTypePair.getDataTypeClass() == byte[].class) {
            byte[] value = (byte[]) dataTypePair.getValue();
            byte[] mappedValue = Arrays.copyOfRange(value, 0, size);
            byte[] mappedByteRepresentation = Arrays.copyOfRange(byteRepresentation, 0, size);
            return Plc4XSupportedDataTypes.DataTypePair.of(mappedValue, mappedByteRepresentation);
        } else if (dataTypePair.getDataTypeClass() == Byte[].class) {
            Byte[] value = (Byte[]) dataTypePair.getValue();
            Byte[] mappedValue = Arrays.copyOfRange(value, 0, size);
            byte[] mappedByteRepresentation = Arrays.copyOfRange(byteRepresentation, 0, size);
            return Plc4XSupportedDataTypes.DataTypePair.of(mappedValue, mappedByteRepresentation);
        } else {
            return dataTypePair;
        }
    }

    private void syncInvoiceId() throws Exception {
        Field transactionId = SUT.getClass().getDeclaredField("transactionId");
        transactionId.setAccessible(true);
        AtomicInteger correlationBuilder = (AtomicInteger) transactionId.get(SUT);

        Field invokeIdField = ModbusTcpPayload.class.getDeclaredField("transactionId");
        Field modifiersField = Field.class.getDeclaredField("modifiers");
        modifiersField.setAccessible(true);
        modifiersField.setInt(invokeIdField, invokeIdField.getModifiers() & ~Modifier.FINAL);
        invokeIdField.setAccessible(true);
        invokeIdField.set(modbusTcpPayload, (short) (correlationBuilder.get() - 1));
    }
}