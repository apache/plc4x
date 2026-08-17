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
package org.apache.plc4x.java.modbus.tcp;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.tag.*;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.modbus.tcp.config.ModbusTcpConfiguration;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class ModbusTcpConnectionTest {

    private ModbusTcpConnection connection;

    @BeforeEach
    void setUp() {
        ModbusTcpConfiguration config = new ModbusTcpConfiguration();
        config.setRequestTimeout(5000);
        config.setDefaultUnitIdentifier(1);
        config.setPingAddress("4x00001:BOOL");
        config.setDefaultPayloadByteOrder(ModbusByteOrder.BIG_ENDIAN);
        config.setMaxCoilsPerRequest(2000);
        config.setMaxRegistersPerRequest(125);

        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        AuditLog auditLog = mock(AuditLog.class);
        when(auditLog.isEnabled()).thenReturn(false);

        connection = new ModbusTcpConnection(config, transportInstance, auditLog);
    }

    @Test
    void testIsConnected_beforeConnect() {
        assertFalse(connection.isConnected());
    }

    @Test
    void testByteSwap_evenLength() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("byteSwap", byte[].class);
        method.setAccessible(true);
        byte[] input = new byte[]{0x01, 0x02, 0x03, 0x04};
        byte[] result = (byte[]) method.invoke(null, (Object) input);
        assertArrayEquals(new byte[]{0x02, 0x01, 0x04, 0x03}, result);
    }

    @Test
    void testByteSwap_oddLength() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("byteSwap", byte[].class);
        method.setAccessible(true);
        byte[] input = new byte[]{0x01, 0x02, 0x03};
        byte[] result = (byte[]) method.invoke(null, (Object) input);
        assertArrayEquals(new byte[]{0x02, 0x01, 0x03}, result);
    }

    @Test
    void testReverseBitsOfByte() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("reverseBitsOfByte", byte.class);
        method.setAccessible(true);
        // 0b00000001 -> 0b10000000
        assertEquals((byte) 0x80, method.invoke(null, (byte) 0x01));
        // 0b10101010 -> 0b01010101
        assertEquals((byte) 0x55, method.invoke(null, (byte) 0xAA));
        // 0b00000000 -> 0b00000000
        assertEquals((byte) 0x00, method.invoke(null, (byte) 0x00));
    }

    @Test
    void testGetUnitId_withModbusTag() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getUnitId", PlcTag.class);
        method.setAccessible(true);

        // Tag without explicit unit id — should use default (1)
        ModbusTag tag = ModbusTag.of("4x00001:INT");
        short unitId = (short) method.invoke(connection, tag);
        assertEquals(1, unitId);
    }

    @Test
    void testGetUnitId_withUnitIdInTag() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getUnitId", PlcTag.class);
        method.setAccessible(true);

        // Tag with explicit unit id
        ModbusTag tag = ModbusTag.of("4x00001:INT{unit-id: 5}");
        short unitId = (short) method.invoke(connection, tag);
        assertEquals(5, unitId);
    }

    @Test
    void testGetReadRequestPdu_discreteInput() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("1x00001:BOOL");
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag);
        assertInstanceOf(ModbusPDUReadDiscreteInputsRequest.class, pdu);
    }

    @Test
    void testGetReadRequestPdu_coil() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("0x00001:BOOL");
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag);
        assertInstanceOf(ModbusPDUReadCoilsRequest.class, pdu);
    }

    @Test
    void testGetReadRequestPdu_inputRegister() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("3x00001:INT");
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag);
        assertInstanceOf(ModbusPDUReadInputRegistersRequest.class, pdu);
    }

    @Test
    void testGetReadRequestPdu_holdingRegister() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag);
        assertInstanceOf(ModbusPDUReadHoldingRegistersRequest.class, pdu);
    }

    @Test
    void testGetReadRequestPdu_extendedRegister() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("6x00001:INT");
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag);
        assertInstanceOf(ModbusPDUReadFileRecordRequest.class, pdu);
    }

    @Test
    void testGetErrorCode() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getErrorCode", ModbusPDUError.class);
        method.setAccessible(true);

        assertEquals(PlcResponseCode.UNSUPPORTED, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.ILLEGAL_FUNCTION)));
        assertEquals(PlcResponseCode.INVALID_ADDRESS, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.ILLEGAL_DATA_ADDRESS)));
        assertEquals(PlcResponseCode.INVALID_DATA, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.ILLEGAL_DATA_VALUE)));
        assertEquals(PlcResponseCode.REMOTE_ERROR, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.SLAVE_DEVICE_FAILURE)));
        assertEquals(PlcResponseCode.OK, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.ACKNOWLEDGE)));
        assertEquals(PlcResponseCode.REMOTE_BUSY, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.SLAVE_DEVICE_BUSY)));
        assertEquals(PlcResponseCode.REMOTE_ERROR, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.NEGATIVE_ACKNOWLEDGE)));
        assertEquals(PlcResponseCode.INTERNAL_ERROR, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.MEMORY_PARITY_ERROR)));
        assertEquals(PlcResponseCode.INTERNAL_ERROR, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.GATEWAY_PATH_UNAVAILABLE)));
        assertEquals(PlcResponseCode.REMOTE_ERROR, method.invoke(connection,
            new ModbusPDUError(ModbusErrorCode.GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND)));
    }

    @Test
    void testGetEffectiveByteOrder_default() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getEffectiveByteOrder", ModbusTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        ModbusByteOrder result = (ModbusByteOrder) method.invoke(connection, tag);
        assertEquals(ModbusByteOrder.BIG_ENDIAN, result);
    }

    @Test
    void testExtractResponseData_coils() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("extractResponseData", ModbusPDU.class, ModbusPDU.class);
        method.setAccessible(true);

        ModbusPDUReadCoilsRequest request = new ModbusPDUReadCoilsRequest(0, 8);
        ModbusPDUReadCoilsResponse response = new ModbusPDUReadCoilsResponse(new byte[]{(byte) 0xFF});
        byte[] data = (byte[]) method.invoke(connection, request, response);
        assertNotNull(data);
        assertEquals(1, data.length);
        assertEquals((byte) 0xFF, data[0]);
    }

    @Test
    void testExtractResponseData_discreteInputs() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("extractResponseData", ModbusPDU.class, ModbusPDU.class);
        method.setAccessible(true);

        ModbusPDUReadDiscreteInputsRequest request = new ModbusPDUReadDiscreteInputsRequest(0, 8);
        ModbusPDUReadDiscreteInputsResponse response = new ModbusPDUReadDiscreteInputsResponse(new byte[]{0x55});
        byte[] data = (byte[]) method.invoke(connection, request, response);
        assertNotNull(data);
        assertEquals(0x55, data[0]);
    }

    @Test
    void testExtractResponseData_holdingRegisters() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("extractResponseData", ModbusPDU.class, ModbusPDU.class);
        method.setAccessible(true);

        ModbusPDUReadHoldingRegistersRequest request = new ModbusPDUReadHoldingRegistersRequest(0, 1);
        ModbusPDUReadHoldingRegistersResponse response = new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A});
        byte[] data = (byte[]) method.invoke(connection, request, response);
        assertNotNull(data);
        assertEquals(2, data.length);
    }

    @Test
    void testExtractResponseData_inputRegisters() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("extractResponseData", ModbusPDU.class, ModbusPDU.class);
        method.setAccessible(true);

        ModbusPDUReadInputRegistersRequest request = new ModbusPDUReadInputRegistersRequest(0, 1);
        ModbusPDUReadInputRegistersResponse response = new ModbusPDUReadInputRegistersResponse(new byte[]{0x00, 0x2A});
        byte[] data = (byte[]) method.invoke(connection, request, response);
        assertNotNull(data);
        assertEquals(2, data.length);
    }

    @Test
    void testExtractResponseData_mismatchedTypes() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("extractResponseData", ModbusPDU.class, ModbusPDU.class);
        method.setAccessible(true);

        ModbusPDUReadCoilsRequest request = new ModbusPDUReadCoilsRequest(0, 8);
        ModbusPDUReadHoldingRegistersResponse response = new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A});
        byte[] data = (byte[]) method.invoke(connection, request, response);
        assertNull(data);
    }

    @Test
    void testGetWriteRequestPdu_coilSingleBool() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getWriteRequestPdu", PlcTag.class, PlcValue.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("0x00001:BOOL");
        PlcValue value = new PlcBOOL(true);
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag, value);
        assertInstanceOf(ModbusPDUWriteMultipleCoilsRequest.class, pdu);
    }

    @Test
    void testGetWriteRequestPdu_holdingRegisterSingleValue() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getWriteRequestPdu", PlcTag.class, PlcValue.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        PlcValue value = new PlcINT(42);
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag, value);
        assertInstanceOf(ModbusPDUWriteSingleRegisterRequest.class, pdu);
    }

    @Test
    void testGetWriteRequestPdu_extendedRegister() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getWriteRequestPdu", PlcTag.class, PlcValue.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("6x00001:INT");
        PlcValue value = new PlcINT(42);
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag, value);
        assertInstanceOf(ModbusPDUWriteFileRecordRequest.class, pdu);
    }

    @Test
    void testFromPlcValueCoil_singleBool() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("fromPlcValueCoil", PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        byte[] result = (byte[]) method.invoke(connection, new PlcBOOL(true), ModbusByteOrder.BIG_ENDIAN);
        assertEquals(1, result.length);
        assertEquals(1, result[0]);

        result = (byte[]) method.invoke(connection, new PlcBOOL(false), ModbusByteOrder.BIG_ENDIAN);
        assertEquals(1, result.length);
        assertEquals(0, result[0]);
    }

    @Test
    void testFromPlcValueCoil_boolList() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("fromPlcValueCoil", PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        PlcList list = new PlcList(Arrays.asList(new PlcBOOL(true), new PlcBOOL(false), new PlcBOOL(true)));
        byte[] result = (byte[]) method.invoke(connection, list, ModbusByteOrder.BIG_ENDIAN);
        assertNotNull(result);
        assertTrue(result.length > 0);
    }

    @Test
    void testFromPlcValueCoil_invalidValue() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("fromPlcValueCoil", PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        try {
            method.invoke(connection, new PlcINT(42), ModbusByteOrder.BIG_ENDIAN);
            fail("Expected PlcRuntimeException");
        } catch (java.lang.reflect.InvocationTargetException e) {
            assertInstanceOf(PlcRuntimeException.class, e.getCause());
        }
    }

    @Test
    void testCreateWriteBuffer_bigEndian() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("createWriteBuffer", int.class, ModbusByteOrder.class);
        method.setAccessible(true);
        assertNotNull(method.invoke(connection, 10, ModbusByteOrder.BIG_ENDIAN));
    }

    @Test
    void testCreateWriteBuffer_littleEndian() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("createWriteBuffer", int.class, ModbusByteOrder.class);
        method.setAccessible(true);
        assertNotNull(method.invoke(connection, 10, ModbusByteOrder.LITTLE_ENDIAN));
    }

    @Test
    void testCreateWriteBuffer_byteSwap() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("createWriteBuffer", int.class, ModbusByteOrder.class);
        method.setAccessible(true);
        assertNotNull(method.invoke(connection, 10, ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP));
        assertNotNull(method.invoke(connection, 10, ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP));
    }

    @Test
    void testNextTransactionId() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("nextTransactionId");
        method.setAccessible(true);

        int first = (int) method.invoke(connection);
        int second = (int) method.invoke(connection);
        assertEquals(first + 1, second);
    }

    @Test
    void testNextTransactionId_wraparound() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("nextTransactionId");
        method.setAccessible(true);

        // Set the generator close to wraparound
        Field generatorField = ModbusTcpConnection.class.getDeclaredField("transactionIdentifierGenerator");
        generatorField.setAccessible(true);
        AtomicInteger generator = (AtomicInteger) generatorField.get(connection);
        generator.set(0xFFFE);

        int id = (int) method.invoke(connection);
        assertEquals(0xFFFE, id);
        // After this call, the generator should have wrapped
        int nextId = (int) method.invoke(connection);
        assertTrue(nextId > 0 && nextId < 0xFFFF);
    }

    @Test
    void testHandleIncomingMessage() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("handleIncomingMessage", ModbusTcpADU.class);
        method.setAccessible(true);

        // Set up a pending request
        Field pendingField = ModbusTcpConnection.class.getDeclaredField("pendingRequests");
        pendingField.setAccessible(true);
        @SuppressWarnings("unchecked")
        Map<Integer, CompletableFuture<ModbusTcpADU>> pendingRequests =
            (Map<Integer, CompletableFuture<ModbusTcpADU>>) pendingField.get(connection);

        CompletableFuture<ModbusTcpADU> future = new CompletableFuture<>();
        pendingRequests.put(42, future);

        // Create a response ADU
        ModbusTcpADU responseAdu = new ModbusTcpADU(42, (short) 1,
            new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A}));

        method.invoke(connection, responseAdu);
        assertTrue(future.isDone());
        assertEquals(responseAdu, future.get());
    }

    @Test
    void testHandleIncomingMessage_unknownTransaction() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("handleIncomingMessage", ModbusTcpADU.class);
        method.setAccessible(true);

        // No pending request for this transaction ID - should just log a warning
        ModbusTcpADU responseAdu = new ModbusTcpADU(999, (short) 1,
            new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A}));

        // Should not throw
        method.invoke(connection, responseAdu);
    }

    @Test
    void testClose_clearsPendingRequests() throws Exception {
        // Set up pending requests
        Field pendingField = ModbusTcpConnection.class.getDeclaredField("pendingRequests");
        pendingField.setAccessible(true);
        @SuppressWarnings("unchecked")
        Map<Integer, CompletableFuture<ModbusTcpADU>> pendingRequests =
            (Map<Integer, CompletableFuture<ModbusTcpADU>>) pendingField.get(connection);

        CompletableFuture<ModbusTcpADU> future = new CompletableFuture<>();
        pendingRequests.put(1, future);

        connection.close();

        assertTrue(future.isCompletedExceptionally());
        assertTrue(pendingRequests.isEmpty());
    }

    @Test
    void testOnTransportDisconnected() throws Exception {
        // Set up pending requests
        Field pendingField = ModbusTcpConnection.class.getDeclaredField("pendingRequests");
        pendingField.setAccessible(true);
        @SuppressWarnings("unchecked")
        Map<Integer, CompletableFuture<ModbusTcpADU>> pendingRequests =
            (Map<Integer, CompletableFuture<ModbusTcpADU>>) pendingField.get(connection);

        CompletableFuture<ModbusTcpADU> future = new CompletableFuture<>();
        pendingRequests.put(1, future);

        Method method = ModbusTcpConnection.class.getDeclaredMethod("onTransportDisconnected", Throwable.class);
        method.setAccessible(true);
        method.invoke(connection, new RuntimeException("test disconnect"));

        assertTrue(future.isCompletedExceptionally());
        assertTrue(pendingRequests.isEmpty());
    }

    @Test
    void testOnTransportDisconnected_nullCause() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("onTransportDisconnected", Throwable.class);
        method.setAccessible(true);
        // Should not throw with null cause
        method.invoke(connection, (Throwable) null);
    }

    @Test
    void testFromPlcValue_holdingRegister() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        PlcValue value = new PlcINT(42);
        byte[] result = (byte[]) method.invoke(connection, tag, value, ModbusByteOrder.BIG_ENDIAN);
        assertNotNull(result);
        assertEquals(2, result.length);
    }

    @Test
    void testFromPlcValue_withByteSwap() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        PlcValue value = new PlcINT(42);
        byte[] bigEndian = (byte[]) method.invoke(connection, tag, value, ModbusByteOrder.BIG_ENDIAN);
        byte[] byteSwap = (byte[]) method.invoke(connection, tag, value, ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP);
        // Byte-swapped result should have bytes in swapped order
        assertNotNull(byteSwap);
        assertArrayEquals(new byte[]{bigEndian[1], bigEndian[0]}, byteSwap);
    }

    @Test
    void testFromPlcValue_littleEndian() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        PlcValue value = new PlcINT(42);
        byte[] result = (byte[]) method.invoke(connection, tag, value, ModbusByteOrder.LITTLE_ENDIAN);
        assertNotNull(result);
        assertEquals(2, result.length);
    }

    @Test
    void testToPlcValue() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("toPlcValue",
            ModbusPDU.class, ModbusPDU.class, ModbusDataType.class, int.class, ModbusByteOrder.class, int.class);
        method.setAccessible(true);

        ModbusPDUReadHoldingRegistersRequest request = new ModbusPDUReadHoldingRegistersRequest(0, 1);
        ModbusPDUReadHoldingRegistersResponse response = new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A});
        PlcValue result = (PlcValue) method.invoke(connection, request, response, ModbusDataType.INT, 1, ModbusByteOrder.BIG_ENDIAN, 1);
        assertNotNull(result);
        assertEquals(42, result.getInteger());
    }

    @Test
    void testToPlcValue_withByteSwap() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("toPlcValue",
            ModbusPDU.class, ModbusPDU.class, ModbusDataType.class, int.class, ModbusByteOrder.class, int.class);
        method.setAccessible(true);

        ModbusPDUReadHoldingRegistersRequest request = new ModbusPDUReadHoldingRegistersRequest(0, 1);
        // Byte-swapped data for value 42 (0x002A) would be stored as 0x2A00
        ModbusPDUReadHoldingRegistersResponse response = new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x2A, 0x00});
        PlcValue result = (PlcValue) method.invoke(connection, request, response, ModbusDataType.INT, 1, ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP, 1);
        assertNotNull(result);
        assertEquals(42, result.getInteger());
    }

    @Test
    void testGetReadRequestPdu_extendedRegister_crossBoundary() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        // Create a tag that crosses the 10000 boundary for extended registers
        ModbusTag tag = new ModbusTagExtendedRegister(9999, 2, ModbusDataType.INT, Collections.emptyMap());
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag);
        assertInstanceOf(ModbusPDUReadFileRecordRequest.class, pdu);
        ModbusPDUReadFileRecordRequest fileRecordRequest = (ModbusPDUReadFileRecordRequest) pdu;
        // Should produce 2 groups since the read crosses the boundary
        assertEquals(2, fileRecordRequest.getItems().size());
    }

    @Test
    void testGetWriteRequestPdu_holdingRegister_multipleWords() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getWriteRequestPdu", PlcTag.class, PlcValue.class);
        method.setAccessible(true);

        // DINT = 4 bytes = 2 words, should use WriteMultipleHoldingRegisters
        ModbusTag tag = ModbusTag.of("4x00001:DINT");
        PlcValue value = new org.apache.plc4x.java.spi.values.PlcDINT(100000);
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag, value);
        assertInstanceOf(ModbusPDUWriteMultipleHoldingRegistersRequest.class, pdu);
    }

    @Test
    void testGetWriteRequestPdu_unsupportedType() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getWriteRequestPdu", PlcTag.class, PlcValue.class);
        method.setAccessible(true);

        // Discrete input is read-only
        ModbusTag tag = ModbusTag.of("1x00001:BOOL");
        try {
            method.invoke(connection, tag, new PlcBOOL(true));
            fail("Expected exception for unsupported write type");
        } catch (java.lang.reflect.InvocationTargetException e) {
            assertInstanceOf(PlcRuntimeException.class, e.getCause());
        }
    }

    @Test
    void testFromPlcValue_littleEndianByteSwap() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        byte[] result = (byte[]) method.invoke(connection, tag, new PlcINT(42), ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP);
        assertNotNull(result);
    }

    @Test
    void testToPlcValue_littleEndian() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("toPlcValue",
            ModbusPDU.class, ModbusPDU.class, ModbusDataType.class, int.class, ModbusByteOrder.class, int.class);
        method.setAccessible(true);

        ModbusPDUReadHoldingRegistersRequest request = new ModbusPDUReadHoldingRegistersRequest(0, 1);
        ModbusPDUReadHoldingRegistersResponse response = new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A});
        PlcValue result = (PlcValue) method.invoke(connection, request, response, ModbusDataType.INT, 1, ModbusByteOrder.LITTLE_ENDIAN, 1);
        assertNotNull(result);
    }

    @Test
    void testToPlcValue_littleEndianByteSwap() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("toPlcValue",
            ModbusPDU.class, ModbusPDU.class, ModbusDataType.class, int.class, ModbusByteOrder.class, int.class);
        method.setAccessible(true);

        ModbusPDUReadHoldingRegistersRequest request = new ModbusPDUReadHoldingRegistersRequest(0, 1);
        ModbusPDUReadHoldingRegistersResponse response = new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A});
        PlcValue result = (PlcValue) method.invoke(connection, request, response, ModbusDataType.INT, 1, ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP, 1);
        assertNotNull(result);
    }

    @Test
    void testGetReadRequestPdu_unsupportedType() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        // Use a mock tag that doesn't match any known type
        PlcTag unknownTag = mock(PlcTag.class);
        try {
            method.invoke(connection, unknownTag);
            fail("Expected exception for unsupported read type");
        } catch (java.lang.reflect.InvocationTargetException e) {
            assertInstanceOf(PlcRuntimeException.class, e.getCause());
        }
    }

    @Test
    void testGetEffectiveByteOrder_withTagOverride() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("getEffectiveByteOrder", ModbusTag.class);
        method.setAccessible(true);

        // Tag with explicit byte order
        ModbusTag tag = ModbusTag.of("4x00001:INT{byte-order: 'LITTLE_ENDIAN'}");
        ModbusByteOrder result = (ModbusByteOrder) method.invoke(connection, tag);
        assertEquals(ModbusByteOrder.LITTLE_ENDIAN, result);
    }

    @Test
    void testExtractResponseData_fileRecord() throws Exception {
        Method method = ModbusTcpConnection.class.getDeclaredMethod("extractResponseData", ModbusPDU.class, ModbusPDU.class);
        method.setAccessible(true);

        ModbusPDUReadFileRecordRequest request = new ModbusPDUReadFileRecordRequest(
            Collections.singletonList(new ModbusPDUReadFileRecordRequestItem((short) 6, 1, 0, 1)));
        ModbusPDUReadFileRecordResponse response = new ModbusPDUReadFileRecordResponse(
            Collections.singletonList(new ModbusPDUReadFileRecordResponseItem((short) 6, new byte[]{0x00, 0x2A})));
        byte[] data = (byte[]) method.invoke(connection, request, response);
        assertNotNull(data);
        assertEquals(2, data.length);
    }
}
