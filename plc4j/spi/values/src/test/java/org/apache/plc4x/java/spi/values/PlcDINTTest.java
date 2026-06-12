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
 * Test class for PlcDINT - Signed 32-bit integer (-2147483648 to 2147483647)
 */
class PlcDINTTest {

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcDINT value = new PlcDINT(Integer.MIN_VALUE);
        assertEquals(Integer.MIN_VALUE, value.getInteger());
        assertEquals(PlcValueType.DINT, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcDINT value = new PlcDINT(Integer.MAX_VALUE);
        assertEquals(Integer.MAX_VALUE, value.getInteger());
    }

    @Test
    void testMinValueMinusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(-2147483649L));
    }

    @Test
    void testMaxValuePlusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(2147483648L));
    }

    @Test
    void testZero() {
        PlcDINT value = new PlcDINT(0);
        assertEquals(0, value.getInteger());
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcDINT trueValue = new PlcDINT(true);
        assertEquals(1, trueValue.getInteger());

        PlcDINT falseValue = new PlcDINT(false);
        assertEquals(0, falseValue.getInteger());
    }

    @Test
    void testByteConstructor() {
        PlcDINT value = new PlcDINT((byte) -100);
        assertEquals(-100, value.getInteger());
    }

    @Test
    void testShortConstructor() {
        PlcDINT value = new PlcDINT((short) -30000);
        assertEquals(-30000, value.getInteger());
    }

    @Test
    void testIntegerConstructor() {
        PlcDINT value = new PlcDINT(-1000000000);
        assertEquals(-1000000000, value.getInteger());
    }

    @Test
    void testLongConstructorValid() {
        PlcDINT value = new PlcDINT(2000000000L);
        assertEquals(2000000000, value.getInteger());
    }

    @Test
    void testLongConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(-3000000000L));
    }

    @Test
    void testLongConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(3000000000L));
    }

    @Test
    void testFloatConstructorValid() {
        PlcDINT value = new PlcDINT(1000000.5f);
        assertEquals(1000000, value.getInteger());
    }

    @Test
    void testFloatConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(3000000000.0f));
    }

    @Test
    void testDoubleConstructorValid() {
        PlcDINT value = new PlcDINT(-2000000000.9);
        assertEquals(-2000000000, value.getInteger());
    }

    @Test
    void testDoubleConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(-3000000000.0));
    }

    @Test
    void testBigIntegerConstructorValid() {
        PlcDINT value = new PlcDINT(BigInteger.valueOf(Integer.MAX_VALUE));
        assertEquals(Integer.MAX_VALUE, value.getInteger());
    }

    @Test
    void testBigIntegerConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(BigInteger.valueOf(-3000000000L)));
    }

    @Test
    void testBigIntegerConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(BigInteger.valueOf(3000000000L)));
    }

    @Test
    void testBigDecimalConstructorValid() {
        PlcDINT value = new PlcDINT(BigDecimal.valueOf(1500000000));
        assertEquals(1500000000, value.getInteger());
    }

    @Test
    void testBigDecimalConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(BigDecimal.valueOf(-2500000000L)));
    }

    @Test
    void testBigDecimalConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT(BigDecimal.valueOf(2500000000L)));
    }

    @Test
    void testStringConstructorValid() {
        PlcDINT value = new PlcDINT("-1500000000");
        assertEquals(-1500000000, value.getInteger());
    }

    @Test
    void testStringConstructorMinValue() {
        PlcDINT value = new PlcDINT("-2147483648");
        assertEquals(Integer.MIN_VALUE, value.getInteger());
    }

    @Test
    void testStringConstructorMaxValue() {
        PlcDINT value = new PlcDINT("2147483647");
        assertEquals(Integer.MAX_VALUE, value.getInteger());
    }

    @Test
    void testStringConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT("-2147483649"));
    }

    @Test
    void testStringConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT("2147483648"));
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDINT("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcDINT value = new PlcDINT("  -1000000000  ");
        assertEquals(-1000000000, value.getInteger());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcDINT() {
        PlcDINT original = new PlcDINT(-1500000000);
        PlcDINT copy = PlcDINT.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithInteger() {
        PlcDINT value = PlcDINT.of(1000000000);
        assertEquals(1000000000, value.getInteger());
    }

    @Test
    void testOfMethodWithString() {
        PlcDINT value = PlcDINT.of("-2000000000");
        assertEquals(-2000000000, value.getInteger());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcDINT zero = new PlcDINT(0);
        assertFalse(zero.getBoolean());

        PlcDINT positive = new PlcDINT(1);
        assertTrue(positive.getBoolean());

        PlcDINT negative = new PlcDINT(-1);
        assertTrue(negative.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcDINT value = new PlcDINT(100);
        assertEquals(100, value.getByte());
    }

    @Test
    void testGetShort() {
        PlcDINT value = new PlcDINT(30000);
        assertEquals(30000, value.getShort());
    }

    @Test
    void testGetInteger() {
        PlcDINT value = new PlcDINT(Integer.MIN_VALUE);
        assertEquals(Integer.MIN_VALUE, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcDINT value = new PlcDINT(-1500000000);
        assertEquals(-1500000000L, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcDINT value = new PlcDINT(Integer.MAX_VALUE);
        assertEquals(BigInteger.valueOf(Integer.MAX_VALUE), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcDINT value = new PlcDINT(-1000000);
        assertEquals(-1000000.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcDINT value = new PlcDINT(2000000000);
        assertEquals(2000000000.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcDINT value = new PlcDINT(-1500000000);
        assertEquals(BigDecimal.valueOf(-1500000000.0f), value.getBigDecimal());
    }

    @Test
    void testGetString() {
        PlcDINT value = new PlcDINT(-1234567890);
        assertEquals("-1234567890", value.getString());
    }

    @Test
    void testToString() {
        PlcDINT value = new PlcDINT(Integer.MAX_VALUE);
        assertEquals("2147483647", value.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcDINT value = new PlcDINT(1);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcDINT smallValue = new PlcDINT(100);
        assertTrue(smallValue.isByte());

        PlcDINT largeValue = new PlcDINT(1000);
        assertFalse(largeValue.isByte());
    }

    @Test
    void testIsShort() {
        PlcDINT smallValue = new PlcDINT(30000);
        assertTrue(smallValue.isShort());

        PlcDINT largeValue = new PlcDINT(40000);
        assertFalse(largeValue.isShort());
    }

    @Test
    void testIsInteger() {
        PlcDINT value = new PlcDINT(Integer.MAX_VALUE);
        assertTrue(value.isInteger());
    }

    @Test
    void testIsLong() {
        PlcDINT value = new PlcDINT(-2000000000);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcDINT value = new PlcDINT(Integer.MAX_VALUE);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcDINT value = new PlcDINT(1000000);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcDINT value = new PlcDINT(-1500000000);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcDINT value = new PlcDINT(Integer.MAX_VALUE);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcDINT value = new PlcDINT(-1234567890);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesMinValue() {
        PlcDINT value = new PlcDINT(Integer.MIN_VALUE);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        assertEquals((byte) 0x80, bytes[0]);
        assertEquals((byte) 0x00, bytes[1]);
        assertEquals((byte) 0x00, bytes[2]);
        assertEquals((byte) 0x00, bytes[3]);
    }

    @Test
    void testGetBytesMaxValue() {
        PlcDINT value = new PlcDINT(Integer.MAX_VALUE);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        assertEquals((byte) 0x7F, bytes[0]);
        assertEquals((byte) 0xFF, bytes[1]);
        assertEquals((byte) 0xFF, bytes[2]);
        assertEquals((byte) 0xFF, bytes[3]);
    }

    @Test
    void testGetBytesZero() {
        PlcDINT value = new PlcDINT(0);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        for (int i = 0; i < 4; i++) {
            assertEquals(0, bytes[i]);
        }
    }

    @Test
    void testGetBytesNegativeValue() {
        PlcDINT value = new PlcDINT(-1);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        for (int i = 0; i < 4; i++) {
            assertEquals((byte) 0xFF, bytes[i]);
        }
    }

    @Test
    void testGetRaw() {
        PlcDINT value = new PlcDINT(-1500000000);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    @Test
    void testByteOrderInSerialization() {
        // Test that byte order is big-endian
        PlcDINT value = new PlcDINT(0x01020304);
        byte[] bytes = value.getBytes();
        assertEquals((byte) 0x01, bytes[0]); // High byte first
        assertEquals((byte) 0x02, bytes[1]);
        assertEquals((byte) 0x03, bytes[2]);
        assertEquals((byte) 0x04, bytes[3]); // Low byte last
    }

    // ========== Edge Cases ==========

    @Test
    void testNegativeOneValue() {
        PlcDINT value = new PlcDINT(-1);
        assertEquals(-1, value.getInteger());
    }

    @Test
    void testPositiveOneValue() {
        PlcDINT value = new PlcDINT(1);
        assertEquals(1, value.getInteger());
    }

    @Test
    void testValueJustBelowMax() {
        PlcDINT value = new PlcDINT(Integer.MAX_VALUE - 1);
        assertEquals(2147483646, value.getInteger());
    }

    @Test
    void testValueJustAboveMin() {
        PlcDINT value = new PlcDINT(Integer.MIN_VALUE + 1);
        assertEquals(-2147483647, value.getInteger());
    }

    @Test
    void testBillionValue() {
        PlcDINT value = new PlcDINT(1000000000);
        assertEquals(1000000000, value.getInteger());
    }

    @Test
    void testNegativeBillionValue() {
        PlcDINT value = new PlcDINT(-1000000000);
        assertEquals(-1000000000, value.getInteger());
    }
}
