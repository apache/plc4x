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

import org.apache.plc4x.java.api.types.PlcValueType;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcNull - Null/empty value representation
 */
class PlcNullTest {

    @Test
    void testGetPlcValueType() {
        PlcNull value = new PlcNull();
        assertEquals(PlcValueType.NULL, value.getPlcValueType());
    }

    @Test
    void testGetObject() {
        PlcNull value = new PlcNull();
        assertNull(value.getObject());
    }

    @Test
    void testIsSimple() {
        PlcNull value = new PlcNull();
        assertTrue(value.isSimple());
    }

    @Test
    void testIsNullable() {
        PlcNull value = new PlcNull();
        assertTrue(value.isNullable());
    }

    @Test
    void testIsNull() {
        PlcNull value = new PlcNull();
        assertTrue(value.isNull());
    }

    @Test
    void testIs() {
        PlcNull value = new PlcNull();
        assertFalse(value.is(String.class));
        assertFalse(value.is(Integer.class));
        assertFalse(value.is(Object.class));
    }

    @Test
    void testIsConvertibleTo() {
        PlcNull value = new PlcNull();
        assertFalse(value.isConvertibleTo(String.class));
        assertFalse(value.isConvertibleTo(Integer.class));
    }

    @Test
    void testGet() {
        PlcNull value = new PlcNull();
        assertNull(value.get(String.class));
        assertNull(value.get(Integer.class));
    }

    @Test
    void testBooleanMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isBoolean());
        assertFalse(value.getBoolean());
    }

    @Test
    void testByteMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isByte());
        assertEquals(0, value.getByte());
    }

    @Test
    void testShortMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isShort());
        assertEquals(0, value.getShort());
    }

    @Test
    void testIntegerMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isInteger());
        assertEquals(0, value.getInteger());
        assertEquals(0, value.getInt());
    }

    @Test
    void testLongMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isLong());
        assertEquals(0L, value.getLong());
    }

    @Test
    void testBigIntegerMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isBigInteger());
        assertNull(value.getBigInteger());
    }

    @Test
    void testFloatMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isFloat());
        assertEquals(0.0f, value.getFloat());
    }

    @Test
    void testDoubleMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isDouble());
        assertEquals(0.0, value.getDouble());
    }

    @Test
    void testBigDecimalMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isBigDecimal());
        assertNull(value.getBigDecimal());
    }

    @Test
    void testStringMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isString());
        assertNull(value.getString());
    }

    @Test
    void testDurationMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isDuration());
        assertNull(value.getDuration());
    }

    @Test
    void testTimeMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isTime());
        assertNull(value.getTime());
    }

    @Test
    void testDateMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isDate());
        assertNull(value.getDate());
    }

    @Test
    void testDateTimeMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isDateTime());
        assertNull(value.getDateTime());
    }

    @Test
    void testGetRaw() {
        PlcNull value = new PlcNull();
        assertArrayEquals(new byte[0], value.getRaw());
    }

    @Test
    void testListMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isList());
        assertEquals(0, value.getLength());
        assertNull(value.getIndex(0));
        assertNull(value.getList());
    }

    @Test
    void testStructMethods() {
        PlcNull value = new PlcNull();
        assertFalse(value.isStruct());
        assertNull(value.getKeys());
        assertFalse(value.hasKey("test"));
        assertNull(value.getValue("test"));
        assertNull(value.getStruct());
    }

    @Test
    void testMetaDataMethods() {
        PlcNull value = new PlcNull();
        assertNotNull(value.getMetaDataNames());
        assertTrue(value.getMetaDataNames().isEmpty());
        assertFalse(value.hasMetaData("test"));
        assertNull(value.getMetaData("test"));
    }

    @Test
    void testEquals() {
        PlcNull value1 = new PlcNull();
        PlcNull value2 = new PlcNull();
        PlcBOOL otherType = new PlcBOOL(true);

        // Same instance
        assertEquals(value1, value1);

        // Different instances, same type
        assertEquals(value1, value2);
        assertEquals(value2, value1);

        // Different type
        assertNotEquals(value1, otherType);

        // Null
        assertNotEquals(value1, null);
    }

    @Test
    void testHashCode() {
        PlcNull value1 = new PlcNull();
        PlcNull value2 = new PlcNull();

        // All PlcNull instances should have same hash code
        assertEquals(value1.hashCode(), value2.hashCode());
        assertEquals(0, value1.hashCode());
    }
}
