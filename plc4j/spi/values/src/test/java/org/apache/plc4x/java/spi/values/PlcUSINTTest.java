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

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.math.BigInteger;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcUSINT - Unsigned 8-bit integer (0 to 255)
 */
class PlcUSINTTest {

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcUSINT value = new PlcUSINT((short) 0);
        assertEquals(0, value.getShort());
        assertEquals(PlcValueType.USINT, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcUSINT value = new PlcUSINT((short) 255);
        assertEquals(255, value.getShort());
    }

    @Test
    void testMinValueMinusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT((short) -1));
    }

    @Test
    void testMaxValuePlusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT((short) 256));
    }

    @Test
    void testMinValueFromInteger() {
        PlcUSINT value = new PlcUSINT(0);
        assertEquals(0, value.getShort());
    }

    @Test
    void testMaxValueFromInteger() {
        PlcUSINT value = new PlcUSINT(255);
        assertEquals(255, value.getShort());
    }

    @Test
    void testBelowRangeFromInteger() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(-1));
    }

    @Test
    void testAboveRangeFromInteger() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(256));
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcUSINT trueValue = new PlcUSINT(true);
        assertEquals(1, trueValue.getShort());

        PlcUSINT falseValue = new PlcUSINT(false);
        assertEquals(0, falseValue.getShort());
    }

    @Test
    void testByteConstructor() {
        PlcUSINT value = new PlcUSINT((byte) 127);
        assertEquals(127, value.getShort());
    }

    @Test
    void testByteConstructorNegative() {
        // Negative byte values should throw
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT((byte) -1));
    }

    @Test
    void testLongConstructorValid() {
        PlcUSINT value = new PlcUSINT(200L);
        assertEquals(200, value.getShort());
    }

    @Test
    void testLongConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(300L));
    }

    @Test
    void testFloatConstructorValid() {
        PlcUSINT value = new PlcUSINT(100.5f);
        assertEquals(100, value.getShort());
    }

    @Test
    void testFloatConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(256.0f));
    }

    @Test
    void testDoubleConstructorValid() {
        PlcUSINT value = new PlcUSINT(150.9);
        assertEquals(150, value.getShort());
    }

    @Test
    void testDoubleConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(-1.0));
    }

    @Test
    void testBigIntegerConstructorValid() {
        PlcUSINT value = new PlcUSINT(BigInteger.valueOf(255));
        assertEquals(255, value.getShort());
    }

    @Test
    void testBigIntegerConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(BigInteger.valueOf(-1)));
    }

    @Test
    void testBigIntegerConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(BigInteger.valueOf(256)));
    }

    @Test
    void testBigDecimalConstructorValid() {
        PlcUSINT value = new PlcUSINT(BigDecimal.valueOf(200));
        assertEquals(200, value.getShort());
    }

    @Test
    void testBigDecimalConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(BigDecimal.valueOf(-1)));
    }

    @Test
    void testBigDecimalConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT(BigDecimal.valueOf(256)));
    }

    @Test
    void testStringConstructorValid() {
        PlcUSINT value = new PlcUSINT("100");
        assertEquals(100, value.getShort());
    }

    @Test
    void testStringConstructorMinValue() {
        PlcUSINT value = new PlcUSINT("0");
        assertEquals(0, value.getShort());
    }

    @Test
    void testStringConstructorMaxValue() {
        PlcUSINT value = new PlcUSINT("255");
        assertEquals(255, value.getShort());
    }

    @Test
    void testStringConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT("-1"));
    }

    @Test
    void testStringConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT("256"));
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUSINT("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcUSINT value = new PlcUSINT("  100  ");
        assertEquals(100, value.getShort());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcUSINT() {
        PlcUSINT original = new PlcUSINT((short) 100);
        PlcUSINT copy = PlcUSINT.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithInteger() {
        PlcUSINT value = PlcUSINT.of(200);
        assertEquals(200, value.getShort());
    }

    @Test
    void testOfMethodWithString() {
        PlcUSINT value = PlcUSINT.of("150");
        assertEquals(150, value.getShort());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcUSINT zero = new PlcUSINT((short) 0);
        assertFalse(zero.getBoolean());

        PlcUSINT nonZero = new PlcUSINT((short) 1);
        assertTrue(nonZero.getBoolean());

        PlcUSINT large = new PlcUSINT((short) 255);
        assertTrue(large.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcUSINT value = new PlcUSINT((short) 127);
        assertEquals(127, value.getByte());
    }

    @Test
    void testGetShort() {
        PlcUSINT value = new PlcUSINT((short) 200);
        assertEquals(200, value.getShort());
    }

    @Test
    void testGetInteger() {
        PlcUSINT value = new PlcUSINT((short) 255);
        assertEquals(255, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcUSINT value = new PlcUSINT((short) 100);
        assertEquals(100L, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcUSINT value = new PlcUSINT((short) 255);
        assertEquals(BigInteger.valueOf(255), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcUSINT value = new PlcUSINT((short) 100);
        assertEquals(100.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcUSINT value = new PlcUSINT((short) 200);
        assertEquals(200.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcUSINT value = new PlcUSINT((short) 150);
        assertEquals(BigDecimal.valueOf(150.0f), value.getBigDecimal());
    }

    @Test
    void testGetString() {
        PlcUSINT value = new PlcUSINT((short) 123);
        assertEquals("123", value.getString());
    }

    @Test
    void testToString() {
        PlcUSINT value = new PlcUSINT((short) 255);
        assertEquals("255", value.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcUSINT value = new PlcUSINT((short) 1);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcUSINT value = new PlcUSINT((short) 100);
        assertTrue(value.isByte());
    }

    @Test
    void testIsShort() {
        PlcUSINT value = new PlcUSINT((short) 200);
        assertTrue(value.isShort());
    }

    @Test
    void testIsInteger() {
        PlcUSINT value = new PlcUSINT((short) 255);
        assertTrue(value.isInteger());
    }

    @Test
    void testIsLong() {
        PlcUSINT value = new PlcUSINT((short) 100);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcUSINT value = new PlcUSINT((short) 200);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcUSINT value = new PlcUSINT((short) 150);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcUSINT value = new PlcUSINT((short) 100);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcUSINT value = new PlcUSINT((short) 255);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcUSINT value = new PlcUSINT((short) 123);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesMinValue() {
        PlcUSINT value = new PlcUSINT((short) 0);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals(0, bytes[0]);
    }

    @Test
    void testGetBytesMaxValue() {
        PlcUSINT value = new PlcUSINT((short) 255);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals((byte) 0xFF, bytes[0]);
    }

    @Test
    void testGetBytesMidValue() {
        PlcUSINT value = new PlcUSINT((short) 128);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals((byte) 0x80, bytes[0]);
    }

    @Test
    void testGetRaw() {
        PlcUSINT value = new PlcUSINT((short) 100);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    // ========== Edge Cases ==========

    @Test
    void testMidRangeValue() {
        PlcUSINT value = new PlcUSINT((short) 127);
        assertEquals(127, value.getShort());
    }

    @Test
    void testPowerOfTwo() {
        PlcUSINT value = new PlcUSINT((short) 128);
        assertEquals(128, value.getShort());
    }

    @Test
    void testValueJustBelowMax() {
        PlcUSINT value = new PlcUSINT((short) 254);
        assertEquals(254, value.getShort());
    }

    @Test
    void testValueJustAboveMin() {
        PlcUSINT value = new PlcUSINT((short) 1);
        assertEquals(1, value.getShort());
    }
}
