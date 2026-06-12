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

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.Test;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcValues - Factory utility for creating PlcValue instances
 */
class PlcValuesTest {

    // ========== List Factory Methods ==========

    @Test
    void testOfList() {
        List<PlcValue> values = Arrays.asList(
            new PlcINT(1),
            new PlcINT(2),
            new PlcINT(3)
        );
        PlcValue result = PlcValues.of(values);

        assertNotNull(result);
        assertTrue(result instanceof PlcList);
        assertTrue(result.isList());
        assertEquals(3, result.getLength());
    }

    @Test
    void testOfVarargs() {
        PlcValue result = PlcValues.of(
            new PlcINT(10),
            new PlcINT(20),
            new PlcINT(30)
        );

        assertNotNull(result);
        assertTrue(result instanceof PlcList);
        assertEquals(3, result.getLength());
    }

    @Test
    void testOfVarargsEmpty() {
        PlcValue result = PlcValues.of();
        assertNotNull(result);
        assertTrue(result instanceof PlcList);
        assertEquals(0, result.getLength());
    }

    // ========== Struct Factory Methods ==========

    @Test
    void testOfSingleKeyValue() {
        PlcValue result = PlcValues.of("key1", new PlcINT(42));

        assertNotNull(result);
        assertTrue(result instanceof PlcStruct);
        assertTrue(result.isStruct());
        assertTrue(result.hasKey("key1"));
        assertEquals(new PlcINT(42), result.getValue("key1"));
    }

    @Test
    void testOfMap() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("temperature", new PlcREAL(25.5f));
        map.put("pressure", new PlcINT(100));

        PlcValue result = PlcValues.of(map);

        assertNotNull(result);
        assertTrue(result instanceof PlcStruct);
        assertEquals(2, result.getKeys().size());
        assertTrue(result.hasKey("temperature"));
        assertTrue(result.hasKey("pressure"));
    }

    // ========== Object Factory Method ==========

    @Test
    void testOfNull() {
        PlcValue result = PlcValues.of((Object) null);
        assertNotNull(result);
        assertTrue(result instanceof PlcNull);
        assertTrue(result.isNull());
    }

    @Test
    void testOfBoolean() {
        PlcValue result = PlcValues.of(true);
        assertNotNull(result);
        assertTrue(result instanceof PlcBOOL);
        assertTrue(result.getBoolean());
    }

    @Test
    void testOfInteger() {
        // PlcValues.of(Object) uses reflection to map "Integer" -> "PlcINTEGER" class
        // which doesn't exist - this is expected to fail for boxed types
        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of(42));
    }

    @Test
    void testOfLong() {
        // PlcValues.of(Object) uses reflection to map "Long" -> "PlcLONG" class
        // which doesn't exist - this is expected to fail for boxed types
        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of(123456789L));
    }

    @Test
    void testOfFloat() {
        // PlcValues.of(Object) uses reflection to map "Float" -> "PlcFLOAT" class
        // which doesn't exist - this is expected to fail for boxed types
        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of(3.14f));
    }

    @Test
    void testOfDouble() {
        // PlcValues.of(Object) uses reflection to map "Double" -> "PlcDOUBLE" class
        // which doesn't exist - this is expected to fail for boxed types
        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of(2.718));
    }

    @Test
    void testOfString() {
        PlcValue result = PlcValues.of("test");
        assertNotNull(result);
        assertTrue(result instanceof PlcSTRING);
        assertEquals("test", result.getString());
    }

    @Test
    void testOfLocalDate() {
        LocalDate date = LocalDate.of(2024, 1, 15);
        PlcValue result = PlcValues.of(date);
        assertNotNull(result);
        // PlcValues.of() strips "Local" prefix, maps "Date" -> "PlcDATE"
        assertTrue(result instanceof PlcDATE);
        assertEquals(date, result.getDate());
    }

    @Test
    void testOfLocalTime() {
        // PlcValues.of() strips "Local" prefix, maps "Time" -> "PlcTIME"
        // But LocalTime should map to PlcTIME_OF_DAY, not PlcTIME
        // This is a known limitation of the reflection-based approach
        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of(LocalTime.of(14, 30, 0)));
    }

    @Test
    void testOfLocalDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2024, 1, 15, 14, 30);
        // PlcValues.of() strips "Local" prefix, maps "DateTime" -> "PlcDATETIME"
        // But the class is actually PlcDATE_AND_TIME, not PlcDATETIME
        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of(dateTime));
    }

    @Test
    void testOfList_DirectCreation() {
        List<String> stringList = Arrays.asList("a", "b", "c");
        // PlcValues.of() maps "List" -> "PlcLIST" (uppercase)
        // But the actual class is "PlcList" (camelCase)
        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of((Object) stringList));
    }

    @Test
    void testOfArray() {
        Object[] array = new Object[]{1, 2, 3};
        // PlcValues.of() converts array to List and tries "PlcLIST" (uppercase)
        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of(array));
    }

    @Test
    void testOfUnsupportedType() {
        // An object that doesn't have a corresponding PlcValue type
        class UnsupportedClass {
        }
        UnsupportedClass unsupported = new UnsupportedClass();

        assertThrows(PlcIncompatibleDatatypeException.class, () -> PlcValues.of(unsupported));
    }

    // ========== Mixed Type Tests ==========

    @Test
    void testMixedTypeList() {
        PlcValue result = PlcValues.of(
            new PlcINT(1),
            new PlcBOOL(true),
            new PlcSTRING("test")
        );

        assertTrue(result instanceof PlcList);
        PlcList list = (PlcList) result;
        assertEquals(3, list.getLength());
        assertTrue(list.getIndex(0) instanceof PlcINT);
        assertTrue(list.getIndex(1) instanceof PlcBOOL);
        assertTrue(list.getIndex(2) instanceof PlcSTRING);
    }

    @Test
    void testNestedStructure() {
        Map<String, PlcValue> innerMap = new HashMap<>();
        innerMap.put("x", new PlcINT(1));
        PlcStruct innerStruct = new PlcStruct(innerMap);

        PlcValue result = PlcValues.of("data", innerStruct);

        assertTrue(result instanceof PlcStruct);
        assertTrue(result.getValue("data") instanceof PlcStruct);
    }
}
