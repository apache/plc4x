/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.values;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Test class for DefaultPlcValueHandler
 */
class DefaultPlcValueHandlerTest {

    private final DefaultPlcValueHandler handler = new DefaultPlcValueHandler();

    // ========== Basic Single Value Tests ==========

    @Test
    void testNewPlcValue_Boolean() {
        PlcTag tag = createMockTag(PlcValueType.BOOL, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, true);

        assertNotNull(result);
        assertTrue(result instanceof PlcBOOL);
        assertTrue(result.getBoolean());
    }

    @Test
    void testNewPlcValue_Integer() {
        PlcTag tag = createMockTag(PlcValueType.DINT, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, 42);

        assertNotNull(result);
        assertTrue(result instanceof PlcDINT);
        assertEquals(42, result.getInteger());
    }

    @Test
    void testNewPlcValue_String() {
        PlcTag tag = createMockTag(PlcValueType.STRING, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, "test");

        assertNotNull(result);
        assertTrue(result instanceof PlcSTRING);
        assertEquals("test", result.getString());
    }

    // ========== Array Value Tests ==========

    @Test
    void testNewPlcValue_IntArray() {
        ArrayInfo arrayInfo = createMockArrayInfo(3);
        PlcTag tag = createMockTag(PlcValueType.INT, Collections.singletonList(arrayInfo));

        PlcValue result = handler.newPlcValue(tag, new Object[]{(short) 1, (short) 2, (short) 3});

        assertNotNull(result);
        assertTrue(result instanceof PlcList);
        assertEquals(3, result.getLength());
    }

    @Test
    void testNewPlcValue_ArrayFromList() {
        ArrayInfo arrayInfo = createMockArrayInfo(2);
        PlcTag tag = createMockTag(PlcValueType.DINT, Collections.singletonList(arrayInfo));

        List<Integer> values = Arrays.asList(10, 20);
        PlcValue result = handler.newPlcValue(tag, values);

        assertNotNull(result);
        assertTrue(result instanceof PlcList);
        assertEquals(2, result.getLength());
    }

    // ========== Error Handling Tests ==========

    @Test
    void testNewPlcValue_NoValueProvided() {
        PlcTag tag = createMockTag(PlcValueType.INT, Collections.emptyList());

        assertThrows(PlcRuntimeException.class, () -> handler.newPlcValue(tag, new Object[0]));
    }

    @Test
    void testNewPlcValue_MultipleValuesForNonArray() {
        PlcTag tag = createMockTag(PlcValueType.INT, Collections.emptyList());

        assertThrows(PlcRuntimeException.class, () ->
            handler.newPlcValue(tag, new Object[]{1, 2}));
    }

    @Test
    void testNewPlcValue_EmptyCollectionForNonArray() {
        PlcTag tag = createMockTag(PlcValueType.INT, Collections.emptyList());

        assertThrows(PlcRuntimeException.class, () ->
            handler.newPlcValue(tag, new ArrayList<>()));
    }

    @Test
    void testNewPlcValue_CollectionWithMultipleItemsForNonArray() {
        PlcTag tag = createMockTag(PlcValueType.INT, Collections.emptyList());

        assertThrows(PlcRuntimeException.class, () ->
            handler.newPlcValue(tag, Arrays.asList(1, 2)));
    }

    @Test
    void testNewPlcValue_ArraySizeMismatch() {
        ArrayInfo arrayInfo = createMockArrayInfo(3);
        PlcTag tag = createMockTag(PlcValueType.INT, Collections.singletonList(arrayInfo));

        assertThrows(PlcRuntimeException.class, () ->
            handler.newPlcValue(tag, new Object[]{1, 2})); // Expected 3, got 2
    }

    // ========== All PlcValueType Tests ==========

    @Test
    void testNewPlcValue_AllNumericTypes() {
        testValueType(PlcValueType.BYTE, (byte) 1, PlcBYTE.class);
        testValueType(PlcValueType.SINT, (byte) 1, PlcSINT.class);
        testValueType(PlcValueType.USINT, (short) 1, PlcUSINT.class);
        testValueType(PlcValueType.INT, (short) 1, PlcINT.class);
        testValueType(PlcValueType.UINT, 1, PlcUINT.class);
        testValueType(PlcValueType.WORD, 1, PlcWORD.class);
        testValueType(PlcValueType.DINT, 1, PlcDINT.class);
        testValueType(PlcValueType.UDINT, 1L, PlcUDINT.class);
        testValueType(PlcValueType.DWORD, 1L, PlcDWORD.class);
        testValueType(PlcValueType.LINT, 1L, PlcLINT.class);
        testValueType(PlcValueType.ULINT, BigInteger.ONE, PlcULINT.class);
        testValueType(PlcValueType.LWORD, BigInteger.ONE, PlcLWORD.class);
        testValueType(PlcValueType.REAL, 1.0f, PlcREAL.class);
        testValueType(PlcValueType.LREAL, 1.0, PlcLREAL.class);
    }

    @Test
    void testNewPlcValue_CharacterTypes() {
        testValueType(PlcValueType.CHAR, "A", PlcCHAR.class);
        testValueType(PlcValueType.WCHAR, "B", PlcWCHAR.class);
        testValueType(PlcValueType.STRING, "test", PlcSTRING.class);
        testValueType(PlcValueType.WSTRING, "test", PlcWSTRING.class);
    }

    @Test
    void testNewPlcValue_TimeTypes() {
        testValueType(PlcValueType.TIME, Duration.ofSeconds(10), PlcTIME.class);
        testValueType(PlcValueType.LTIME, Duration.ofSeconds(10), PlcLTIME.class);
        testValueType(PlcValueType.DATE, LocalDate.now(), PlcDATE.class);
        testValueType(PlcValueType.LDATE, LocalDate.now(), PlcLDATE.class);
        testValueType(PlcValueType.TIME_OF_DAY, LocalTime.now(), PlcTIME_OF_DAY.class);
        testValueType(PlcValueType.LTIME_OF_DAY, LocalTime.now(), PlcLTIME_OF_DAY.class);
        testValueType(PlcValueType.DATE_AND_TIME, LocalDateTime.now(), PlcDATE_AND_TIME.class);
        testValueType(PlcValueType.DATE_AND_LTIME, LocalDateTime.now(), PlcDATE_AND_LTIME.class);
        testValueType(PlcValueType.LDATE_AND_TIME, LocalDateTime.now(), PlcLDATE_AND_TIME.class);
    }

    @Test
    void testNewPlcValue_RawByteArray() {
        testValueType(PlcValueType.RAW_BYTE_ARRAY, new byte[]{1, 2, 3}, PlcRawByteArray.class);
    }

    @Test
    void testNewPlcValue_Null() {
        PlcTag tag = createMockTag(PlcValueType.NULL, Collections.emptyList());
        // Pass a PlcNull instance rather than Java null
        PlcValue result = handler.newPlcValue(tag, new PlcNull());

        assertNotNull(result);
        assertTrue(result instanceof PlcNull);
    }

    // ========== Static Method Tests ==========

    @Test
    void testStaticOf_SingleValue() {
        PlcTag tag = createMockTag(PlcValueType.INT, Collections.emptyList());
        PlcValue result = DefaultPlcValueHandler.of(tag, (short) 42);

        assertNotNull(result);
        assertTrue(result instanceof PlcINT);
    }

    @Test
    void testStaticOf_ArrayValues() {
        ArrayInfo arrayInfo = createMockArrayInfo(2);
        PlcTag tag = createMockTag(PlcValueType.INT, Collections.singletonList(arrayInfo));

        PlcValue result = DefaultPlcValueHandler.of(tag, new Object[]{(short) 1, (short) 2});

        assertNotNull(result);
        assertTrue(result instanceof PlcList);
    }

    // ========== Legacy Type Inference Tests (null PlcValueType) ==========

    @Test
    void testNewPlcValue_NullType_InfersBoolean() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, true);

        assertNotNull(result);
        assertTrue(result instanceof PlcBOOL);
    }

    @Test
    void testNewPlcValue_NullType_InfersByte() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, (byte) 10);

        assertNotNull(result);
        assertTrue(result instanceof PlcSINT);
    }

    @Test
    void testNewPlcValue_NullType_InfersShort() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, (short) 100);

        assertNotNull(result);
        assertTrue(result instanceof PlcINT);
    }

    @Test
    void testNewPlcValue_NullType_InfersInteger() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, 1000);

        assertNotNull(result);
        assertTrue(result instanceof PlcDINT);
    }

    @Test
    void testNewPlcValue_NullType_InfersLong() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, 10000L);

        assertNotNull(result);
        assertTrue(result instanceof PlcLINT);
    }

    @Test
    void testNewPlcValue_NullType_InfersBigInteger() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, BigInteger.valueOf(12345));

        assertNotNull(result);
        assertTrue(result instanceof PlcULINT);
    }

    @Test
    void testNewPlcValue_NullType_InfersFloat() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, 3.14f);

        assertNotNull(result);
        assertTrue(result instanceof PlcREAL);
    }

    @Test
    void testNewPlcValue_NullType_InfersDouble() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, 2.718);

        assertNotNull(result);
        assertTrue(result instanceof PlcLREAL);
    }

    @Test
    void testNewPlcValue_NullType_InfersString() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, "test");

        assertNotNull(result);
        assertTrue(result instanceof PlcSTRING);
    }

    @Test
    void testNewPlcValue_NullType_InfersByteArray() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, new byte[]{1, 2, 3});

        assertNotNull(result);
        assertTrue(result instanceof PlcRawByteArray);
    }

    @Test
    void testNewPlcValue_NullType_InfersDuration() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, Duration.ofSeconds(10));

        assertNotNull(result);
        assertTrue(result instanceof PlcTIME);
    }

    @Test
    void testNewPlcValue_NullType_InfersLocalTime() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, LocalTime.of(14, 30));

        assertNotNull(result);
        assertTrue(result instanceof PlcTIME_OF_DAY);
    }

    @Test
    void testNewPlcValue_NullType_InfersLocalDate() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, LocalDate.of(2024, 1, 15));

        assertNotNull(result);
        assertTrue(result instanceof PlcDATE);
    }

    @Test
    void testNewPlcValue_NullType_InfersLocalDateTime() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, LocalDateTime.of(2024, 1, 15, 14, 30));

        assertNotNull(result);
        assertTrue(result instanceof PlcDATE_AND_TIME);
    }

    @Test
    void testNewPlcValue_NullType_InfersPlcValue() {
        PlcTag tag = createMockTag(null, Collections.emptyList());
        PlcINT originalValue = new PlcINT(42);
        PlcValue result = handler.newPlcValue(tag, originalValue);

        assertSame(originalValue, result);
    }

    // ========== Helper Methods ==========

    private void testValueType(PlcValueType type, Object value, Class<? extends PlcValue> expectedClass) {
        PlcTag tag = createMockTag(type, Collections.emptyList());
        PlcValue result = handler.newPlcValue(tag, value);

        assertNotNull(result);
        assertTrue(expectedClass.isInstance(result),
            "Expected " + expectedClass.getSimpleName() + " but got " + result.getClass().getSimpleName());
    }

    private PlcTag createMockTag(PlcValueType type, List<ArrayInfo> arrayInfo) {
        PlcTag tag = mock(PlcTag.class);
        when(tag.getPlcValueType()).thenReturn(type);
        when(tag.getArrayInfo()).thenReturn(arrayInfo);
        return tag;
    }

    private ArrayInfo createMockArrayInfo(int size) {
        ArrayInfo arrayInfo = mock(ArrayInfo.class);
        when(arrayInfo.getSize()).thenReturn(size);
        return arrayInfo;
    }
}
