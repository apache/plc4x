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
package org.apache.plc4x.java.modbus.rtu;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.tag.*;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.modbus.rtu.config.ModbusRtuConfiguration;
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
import java.util.Map;
import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class ModbusRtuConnectionTest {

    private ModbusRtuConnection connection;

    @BeforeEach
    void setUp() {
        ModbusRtuConfiguration config = new ModbusRtuConfiguration();
        config.setRequestTimeout(5000);
        config.setDefaultUnitIdentifier(1);
        config.setPingAddress("4x00001:BOOL");
        config.setDefaultPayloadByteOrder(ModbusByteOrder.BIG_ENDIAN);
        config.setMaxCoilsPerRequest(2000);
        config.setMaxRegistersPerRequest(125);

        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        AuditLog auditLog = mock(AuditLog.class);
        when(auditLog.isEnabled()).thenReturn(false);

        connection = new ModbusRtuConnection(config, transportInstance, auditLog);
    }

    @Test
    void testIsConnected_beforeConnect() {
        assertFalse(connection.isConnected());
    }

    @Test
    void testByteSwap() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("byteSwap", byte[].class);
        method.setAccessible(true);
        byte[] input = new byte[]{0x01, 0x02, 0x03, 0x04};
        byte[] result = (byte[]) method.invoke(null, (Object) input);
        assertArrayEquals(new byte[]{0x02, 0x01, 0x04, 0x03}, result);
    }

    @Test
    void testByteSwap_oddLength() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("byteSwap", byte[].class);
        method.setAccessible(true);
        byte[] input = new byte[]{0x01, 0x02, 0x03};
        byte[] result = (byte[]) method.invoke(null, (Object) input);
        assertArrayEquals(new byte[]{0x02, 0x01, 0x03}, result);
    }

    @Test
    void testReverseBitsOfByte() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("reverseBitsOfByte", byte.class);
        method.setAccessible(true);
        assertEquals((byte) 0x80, method.invoke(null, (byte) 0x01));
        assertEquals((byte) 0x55, method.invoke(null, (byte) 0xAA));
    }

    @Test
    void testGetUnitId() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getUnitId", PlcTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        assertEquals((short) 1, method.invoke(connection, tag));

        ModbusTag tagWithUnit = ModbusTag.of("4x00001:INT{unit-id: 7}");
        assertEquals((short) 7, method.invoke(connection, tagWithUnit));
    }

    @Test
    void testGetReadRequestPdu_allTypes() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        assertInstanceOf(ModbusPDUReadDiscreteInputsRequest.class, method.invoke(connection, ModbusTag.of("1x00001:BOOL")));
        assertInstanceOf(ModbusPDUReadCoilsRequest.class, method.invoke(connection, ModbusTag.of("0x00001:BOOL")));
        assertInstanceOf(ModbusPDUReadInputRegistersRequest.class, method.invoke(connection, ModbusTag.of("3x00001:INT")));
        assertInstanceOf(ModbusPDUReadHoldingRegistersRequest.class, method.invoke(connection, ModbusTag.of("4x00001:INT")));
        assertInstanceOf(ModbusPDUReadFileRecordRequest.class, method.invoke(connection, ModbusTag.of("6x00001:INT")));
    }

    @Test
    void testGetErrorCode() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getErrorCode", ModbusPDUError.class);
        method.setAccessible(true);

        assertEquals(PlcResponseCode.UNSUPPORTED, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.ILLEGAL_FUNCTION)));
        assertEquals(PlcResponseCode.INVALID_ADDRESS, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.ILLEGAL_DATA_ADDRESS)));
        assertEquals(PlcResponseCode.INVALID_DATA, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.ILLEGAL_DATA_VALUE)));
        assertEquals(PlcResponseCode.REMOTE_ERROR, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.SLAVE_DEVICE_FAILURE)));
        assertEquals(PlcResponseCode.OK, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.ACKNOWLEDGE)));
        assertEquals(PlcResponseCode.REMOTE_BUSY, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.SLAVE_DEVICE_BUSY)));
        assertEquals(PlcResponseCode.REMOTE_ERROR, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.NEGATIVE_ACKNOWLEDGE)));
        assertEquals(PlcResponseCode.INTERNAL_ERROR, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.MEMORY_PARITY_ERROR)));
        assertEquals(PlcResponseCode.INTERNAL_ERROR, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.GATEWAY_PATH_UNAVAILABLE)));
        assertEquals(PlcResponseCode.REMOTE_ERROR, method.invoke(connection, new ModbusPDUError(ModbusErrorCode.GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND)));
    }

    @Test
    void testExtractResponseData() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("extractResponseData", ModbusPDU.class, ModbusPDU.class);
        method.setAccessible(true);

        // Coils
        assertNotNull(method.invoke(connection, new ModbusPDUReadCoilsRequest(0, 8), new ModbusPDUReadCoilsResponse(new byte[]{0x01})));
        // Discrete inputs
        assertNotNull(method.invoke(connection, new ModbusPDUReadDiscreteInputsRequest(0, 8), new ModbusPDUReadDiscreteInputsResponse(new byte[]{0x01})));
        // Holding registers
        assertNotNull(method.invoke(connection, new ModbusPDUReadHoldingRegistersRequest(0, 1), new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x01})));
        // Input registers
        assertNotNull(method.invoke(connection, new ModbusPDUReadInputRegistersRequest(0, 1), new ModbusPDUReadInputRegistersResponse(new byte[]{0x00, 0x01})));
        // Mismatched
        assertNull(method.invoke(connection, new ModbusPDUReadCoilsRequest(0, 8), new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x01})));
    }

    @Test
    void testGetWriteRequestPdu() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getWriteRequestPdu", PlcTag.class, PlcValue.class);
        method.setAccessible(true);

        // Coil write
        assertInstanceOf(ModbusPDUWriteMultipleCoilsRequest.class, method.invoke(connection, ModbusTag.of("0x00001:BOOL"), new PlcBOOL(true)));
        // Holding register single value
        assertInstanceOf(ModbusPDUWriteSingleRegisterRequest.class, method.invoke(connection, ModbusTag.of("4x00001:INT"), new PlcINT(42)));
        // Extended register
        assertInstanceOf(ModbusPDUWriteFileRecordRequest.class, method.invoke(connection, ModbusTag.of("6x00001:INT"), new PlcINT(42)));
    }

    @Test
    void testFromPlcValueCoil() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("fromPlcValueCoil", PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        byte[] result = (byte[]) method.invoke(connection, new PlcBOOL(true), ModbusByteOrder.BIG_ENDIAN);
        assertEquals(1, result[0]);

        result = (byte[]) method.invoke(connection, new PlcBOOL(false), ModbusByteOrder.BIG_ENDIAN);
        assertEquals(0, result[0]);

        // List
        PlcList list = new PlcList(Arrays.asList(new PlcBOOL(true), new PlcBOOL(false)));
        assertNotNull(method.invoke(connection, list, ModbusByteOrder.BIG_ENDIAN));
    }

    @Test
    void testCreateWriteBuffer() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("createWriteBuffer", int.class, ModbusByteOrder.class);
        method.setAccessible(true);
        assertNotNull(method.invoke(connection, 10, ModbusByteOrder.BIG_ENDIAN));
        assertNotNull(method.invoke(connection, 10, ModbusByteOrder.LITTLE_ENDIAN));
        assertNotNull(method.invoke(connection, 10, ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP));
    }

    @Test
    void testGetEffectiveByteOrder() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getEffectiveByteOrder", ModbusTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        assertEquals(ModbusByteOrder.BIG_ENDIAN, method.invoke(connection, tag));
    }

    @Test
    void testGetEffectiveByteOrder_withTagOverride() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getEffectiveByteOrder", ModbusTag.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT{byte-order: 'LITTLE_ENDIAN'}");
        assertEquals(ModbusByteOrder.LITTLE_ENDIAN, method.invoke(connection, tag));
    }

    @Test
    void testHandleIncomingMessage() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("handleIncomingMessage", ModbusRtuADU.class);
        method.setAccessible(true);

        Field pendingField = ModbusRtuConnection.class.getDeclaredField("pendingRequests");
        pendingField.setAccessible(true);
        @SuppressWarnings("unchecked")
        Map<Short, CompletableFuture<ModbusRtuADU>> pendingRequests =
            (Map<Short, CompletableFuture<ModbusRtuADU>>) pendingField.get(connection);

        CompletableFuture<ModbusRtuADU> future = new CompletableFuture<>();
        pendingRequests.put((short) 1, future);

        ModbusRtuADU responseAdu = new ModbusRtuADU((short) 1,
            new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A}));

        method.invoke(connection, responseAdu);
        assertTrue(future.isDone());
    }

    @Test
    void testHandleIncomingMessage_unknownAddress() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("handleIncomingMessage", ModbusRtuADU.class);
        method.setAccessible(true);

        ModbusRtuADU responseAdu = new ModbusRtuADU((short) 99,
            new ModbusPDUReadHoldingRegistersResponse(new byte[]{0x00, 0x2A}));
        // Should not throw
        method.invoke(connection, responseAdu);
    }

    @Test
    void testClose_clearsPendingRequests() throws Exception {
        Field pendingField = ModbusRtuConnection.class.getDeclaredField("pendingRequests");
        pendingField.setAccessible(true);
        @SuppressWarnings("unchecked")
        Map<Short, CompletableFuture<ModbusRtuADU>> pendingRequests =
            (Map<Short, CompletableFuture<ModbusRtuADU>>) pendingField.get(connection);

        CompletableFuture<ModbusRtuADU> future = new CompletableFuture<>();
        pendingRequests.put((short) 1, future);

        connection.close();

        assertTrue(future.isCompletedExceptionally());
        assertTrue(pendingRequests.isEmpty());
    }

    @Test
    void testOnTransportDisconnected() throws Exception {
        Field pendingField = ModbusRtuConnection.class.getDeclaredField("pendingRequests");
        pendingField.setAccessible(true);
        @SuppressWarnings("unchecked")
        Map<Short, CompletableFuture<ModbusRtuADU>> pendingRequests =
            (Map<Short, CompletableFuture<ModbusRtuADU>>) pendingField.get(connection);

        CompletableFuture<ModbusRtuADU> future = new CompletableFuture<>();
        pendingRequests.put((short) 1, future);

        Method method = ModbusRtuConnection.class.getDeclaredMethod("onTransportDisconnected", Throwable.class);
        method.setAccessible(true);
        method.invoke(connection, new RuntimeException("test disconnect"));

        assertTrue(future.isCompletedExceptionally());
    }

    @Test
    void testOnTransportDisconnected_nullCause() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("onTransportDisconnected", Throwable.class);
        method.setAccessible(true);
        method.invoke(connection, (Throwable) null);
    }

    @Test
    void testFromPlcValue_holdingRegister() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        byte[] result = (byte[]) method.invoke(connection, tag, new PlcINT(42), ModbusByteOrder.BIG_ENDIAN);
        assertNotNull(result);
        assertEquals(2, result.length);
    }

    @Test
    void testFromPlcValue_withByteSwap() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        byte[] result = (byte[]) method.invoke(connection, tag, new PlcINT(42), ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP);
        assertNotNull(result);
    }

    @Test
    void testFromPlcValue_littleEndian() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        byte[] result = (byte[]) method.invoke(connection, tag, new PlcINT(42), ModbusByteOrder.LITTLE_ENDIAN);
        assertNotNull(result);
    }

    @Test
    void testExtractResponseData_fileRecord() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("extractResponseData", ModbusPDU.class, ModbusPDU.class);
        method.setAccessible(true);

        ModbusPDUReadFileRecordRequest request = new ModbusPDUReadFileRecordRequest(
            Collections.singletonList(new ModbusPDUReadFileRecordRequestItem((short) 6, 1, 0, 1)));
        ModbusPDUReadFileRecordResponse response = new ModbusPDUReadFileRecordResponse(
            Collections.singletonList(new ModbusPDUReadFileRecordResponseItem((short) 6, new byte[]{0x00, 0x2A})));
        byte[] data = (byte[]) method.invoke(connection, request, response);
        assertNotNull(data);
    }

    @Test
    void testGetReadRequestPdu_extendedRegister_crossBoundary() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getReadRequestPdu", PlcTag.class);
        method.setAccessible(true);

        ModbusTag tag = new ModbusTagExtendedRegister(9999, 2, ModbusDataType.INT, Collections.emptyMap());
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag);
        assertInstanceOf(ModbusPDUReadFileRecordRequest.class, pdu);
        assertEquals(2, ((ModbusPDUReadFileRecordRequest) pdu).getItems().size());
    }

    @Test
    void testGetWriteRequestPdu_holdingRegister_multipleWords() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getWriteRequestPdu", PlcTag.class, PlcValue.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:DINT");
        ModbusPDU pdu = (ModbusPDU) method.invoke(connection, tag, new org.apache.plc4x.java.spi.values.PlcDINT(100000));
        assertInstanceOf(ModbusPDUWriteMultipleHoldingRegistersRequest.class, pdu);
    }

    @Test
    void testFromPlcValue_boolData_reverseBits() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        // BOOL data type on a holding register triggers bit reversal
        ModbusTag tag = ModbusTag.of("4x00001:BOOL");
        byte[] result = (byte[]) method.invoke(connection, tag, new PlcBOOL(true), ModbusByteOrder.BIG_ENDIAN);
        assertNotNull(result);
    }

    @Test
    void testFromPlcValueCoil_withByteSwap() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("fromPlcValueCoil", PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        PlcList list = new PlcList(Arrays.asList(new PlcBOOL(true), new PlcBOOL(false)));
        byte[] result = (byte[]) method.invoke(connection, list, ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP);
        assertNotNull(result);
    }

    @Test
    void testFromPlcValue_littleEndianByteSwap() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("fromPlcValue", PlcTag.class, PlcValue.class, ModbusByteOrder.class);
        method.setAccessible(true);

        ModbusTag tag = ModbusTag.of("4x00001:INT");
        byte[] result = (byte[]) method.invoke(connection, tag, new PlcINT(42), ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP);
        assertNotNull(result);
    }

    @Test
    void testGetWriteRequestPdu_unsupportedType() throws Exception {
        Method method = ModbusRtuConnection.class.getDeclaredMethod("getWriteRequestPdu", PlcTag.class, PlcValue.class);
        method.setAccessible(true);

        try {
            method.invoke(connection, ModbusTag.of("1x00001:BOOL"), new PlcBOOL(true));
            fail("Expected exception");
        } catch (java.lang.reflect.InvocationTargetException e) {
            assertInstanceOf(PlcRuntimeException.class, e.getCause());
        }
    }
}
