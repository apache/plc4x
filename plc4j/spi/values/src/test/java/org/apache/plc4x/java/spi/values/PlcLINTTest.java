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
 * Test class for PlcLINT - Signed 64-bit integer (-9223372036854775808 to 9223372036854775807)
 */
class PlcLINTTest {

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcLINT value = new PlcLINT(Long.MIN_VALUE);
        assertEquals(Long.MIN_VALUE, value.getLong());
        assertEquals(PlcValueType.LINT, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcLINT value = new PlcLINT(Long.MAX_VALUE);
        assertEquals(Long.MAX_VALUE, value.getLong());
    }

    @Test
    void testMinValueMinusOne() {
        BigInteger belowMin = BigInteger.valueOf(Long.MIN_VALUE).subtract(BigInteger.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT(belowMin));
    }

    @Test
    void testMaxValuePlusOne() {
        BigInteger aboveMax = BigInteger.valueOf(Long.MAX_VALUE).add(BigInteger.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT(aboveMax));
    }

    @Test
    void testZero() {
        PlcLINT value = new PlcLINT(0L);
        assertEquals(0L, value.getLong());
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcLINT trueValue = new PlcLINT(true);
        assertEquals(1L, trueValue.getLong());

        PlcLINT falseValue = new PlcLINT(false);
        assertEquals(0L, falseValue.getLong());
    }

    @Test
    void testByteConstructor() {
        PlcLINT value = new PlcLINT((byte) -100);
        assertEquals(-100L, value.getLong());
    }

    @Test
    void testShortConstructor() {
        PlcLINT value = new PlcLINT((short) -30000);
        assertEquals(-30000L, value.getLong());
    }

    @Test
    void testIntegerConstructor() {
        PlcLINT value = new PlcLINT(-1000000000);
        assertEquals(-1000000000L, value.getLong());
    }

    @Test
    void testLongConstructor() {
        PlcLINT value = new PlcLINT(5000000000000L);
        assertEquals(5000000000000L, value.getLong());
    }

    @Test
    void testFloatConstructorValid() {
        PlcLINT value = new PlcLINT(1000000.5f);
        assertEquals(1000000L, value.getLong());
    }

    @Test
    void testDoubleConstructorValid() {
        PlcLINT value = new PlcLINT(-5000000000000.9);
        assertEquals(-5000000000000L, value.getLong());
    }

    @Test
    void testBigIntegerConstructorValid() {
        PlcLINT value = new PlcLINT(BigInteger.valueOf(Long.MAX_VALUE));
        assertEquals(Long.MAX_VALUE, value.getLong());
    }

    @Test
    void testBigIntegerConstructorBelowRange() {
        BigInteger belowMin = BigInteger.valueOf(Long.MIN_VALUE).subtract(BigInteger.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT(belowMin));
    }

    @Test
    void testBigIntegerConstructorAboveRange() {
        BigInteger aboveMax = BigInteger.valueOf(Long.MAX_VALUE).add(BigInteger.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT(aboveMax));
    }

    @Test
    void testBigDecimalConstructorValid() {
        PlcLINT value = new PlcLINT(BigDecimal.valueOf(5000000000000L));
        assertEquals(5000000000000L, value.getLong());
    }

    @Test
    void testBigDecimalConstructorBelowRange() {
        BigDecimal belowMin = new BigDecimal(Long.MIN_VALUE).subtract(BigDecimal.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT(belowMin));
    }

    @Test
    void testBigDecimalConstructorAboveRange() {
        BigDecimal aboveMax = new BigDecimal(Long.MAX_VALUE).add(BigDecimal.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT(aboveMax));
    }

    @Test
    void testStringConstructorValid() {
        PlcLINT value = new PlcLINT("-5000000000000");
        assertEquals(-5000000000000L, value.getLong());
    }

    @Test
    void testStringConstructorMinValue() {
        PlcLINT value = new PlcLINT("-9223372036854775808");
        assertEquals(Long.MIN_VALUE, value.getLong());
    }

    @Test
    void testStringConstructorMaxValue() {
        PlcLINT value = new PlcLINT("9223372036854775807");
        assertEquals(Long.MAX_VALUE, value.getLong());
    }

    @Test
    void testStringConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT("-9223372036854775809"));
    }

    @Test
    void testStringConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT("9223372036854775808"));
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLINT("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcLINT value = new PlcLINT("  -5000000000000  ");
        assertEquals(-5000000000000L, value.getLong());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcLINT() {
        PlcLINT original = new PlcLINT(-5000000000000L);
        PlcLINT copy = PlcLINT.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithLong() {
        PlcLINT value = PlcLINT.of(5000000000000L);
        assertEquals(5000000000000L, value.getLong());
    }

    @Test
    void testOfMethodWithString() {
        PlcLINT value = PlcLINT.of("-7000000000000");
        assertEquals(-7000000000000L, value.getLong());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcLINT zero = new PlcLINT(0L);
        assertFalse(zero.getBoolean());

        PlcLINT positive = new PlcLINT(1L);
        assertTrue(positive.getBoolean());

        PlcLINT negative = new PlcLINT(-1L);
        assertTrue(negative.getBoolean());

        PlcLINT large = new PlcLINT(Long.MAX_VALUE);
        assertTrue(large.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcLINT value = new PlcLINT(100L);
        assertEquals(100, value.getByte());
    }

    @Test
    void testGetShort() {
        PlcLINT value = new PlcLINT(30000L);
        assertEquals(30000, value.getShort());
    }

    @Test
    void testGetInteger() {
        PlcLINT value = new PlcLINT(2000000000L);
        assertEquals(2000000000, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcLINT value = new PlcLINT(Long.MIN_VALUE);
        assertEquals(Long.MIN_VALUE, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcLINT value = new PlcLINT(Long.MAX_VALUE);
        assertEquals(BigInteger.valueOf(Long.MAX_VALUE), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcLINT value = new PlcLINT(-1000000L);
        assertEquals(-1000000.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcLINT value = new PlcLINT(5000000000000L);
        assertEquals(5000000000000.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcLINT value = new PlcLINT(-5000000000000L);
        assertEquals(0, new BigDecimal(-5000000000000L).compareTo(value.getBigDecimal()));
    }

    @Test
    void testGetString() {
        PlcLINT value = new PlcLINT(-5123456789012L);
        assertEquals("-5123456789012", value.getString());
    }

    @Test
    void testToString() {
        PlcLINT value = new PlcLINT(Long.MAX_VALUE);
        assertEquals("9223372036854775807", value.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcLINT value = new PlcLINT(1L);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcLINT smallValue = new PlcLINT(100L);
        assertTrue(smallValue.isByte());

        PlcLINT largeValue = new PlcLINT(1000L);
        assertFalse(largeValue.isByte());
    }

    @Test
    void testIsShort() {
        PlcLINT smallValue = new PlcLINT(30000L);
        assertTrue(smallValue.isShort());

        PlcLINT largeValue = new PlcLINT(40000L);
        assertFalse(largeValue.isShort());
    }

    @Test
    void testIsInteger() {
        PlcLINT smallValue = new PlcLINT(2000000000L);
        assertTrue(smallValue.isInteger());

        PlcLINT largeValue = new PlcLINT(3000000000L);
        assertFalse(largeValue.isInteger());
    }

    @Test
    void testIsLong() {
        PlcLINT value = new PlcLINT(Long.MAX_VALUE);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcLINT value = new PlcLINT(-5000000000000L);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcLINT value = new PlcLINT(1000000L);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcLINT value = new PlcLINT(-5000000000000L);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcLINT value = new PlcLINT(Long.MAX_VALUE);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcLINT value = new PlcLINT(-5123456789012L);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesMinValue() {
        PlcLINT value = new PlcLINT(Long.MIN_VALUE);
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        assertEquals((byte) 0x80, bytes[0]);
        for (int i = 1; i < 8; i++) {
            assertEquals((byte) 0x00, bytes[i]);
        }
    }

    @Test
    void testGetBytesMaxValue() {
        PlcLINT value = new PlcLINT(Long.MAX_VALUE);
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        assertEquals((byte) 0x7F, bytes[0]);
        for (int i = 1; i < 8; i++) {
            assertEquals((byte) 0xFF, bytes[i]);
        }
    }

    @Test
    void testGetBytesZero() {
        PlcLINT value = new PlcLINT(0L);
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        for (int i = 0; i < 8; i++) {
            assertEquals(0, bytes[i]);
        }
    }

    @Test
    void testGetBytesNegativeValue() {
        PlcLINT value = new PlcLINT(-1L);
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        for (int i = 0; i < 8; i++) {
            assertEquals((byte) 0xFF, bytes[i]);
        }
    }

    @Test
    void testGetRaw() {
        PlcLINT value = new PlcLINT(-5000000000000L);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    @Test
    void testByteOrderInSerialization() {
        // Test that byte order is big-endian
        PlcLINT value = new PlcLINT(0x0102030405060708L);
        byte[] bytes = value.getBytes();
        assertEquals((byte) 0x01, bytes[0]); // High byte first
        assertEquals((byte) 0x02, bytes[1]);
        assertEquals((byte) 0x03, bytes[2]);
        assertEquals((byte) 0x04, bytes[3]);
        assertEquals((byte) 0x05, bytes[4]);
        assertEquals((byte) 0x06, bytes[5]);
        assertEquals((byte) 0x07, bytes[6]);
        assertEquals((byte) 0x08, bytes[7]); // Low byte last
    }

    // ========== Edge Cases ==========

    @Test
    void testNegativeOneValue() {
        PlcLINT value = new PlcLINT(-1L);
        assertEquals(-1L, value.getLong());
    }

    @Test
    void testPositiveOneValue() {
        PlcLINT value = new PlcLINT(1L);
        assertEquals(1L, value.getLong());
    }

    @Test
    void testValueJustBelowMax() {
        PlcLINT value = new PlcLINT(Long.MAX_VALUE - 1);
        assertEquals(9223372036854775806L, value.getLong());
    }

    @Test
    void testValueJustAboveMin() {
        PlcLINT value = new PlcLINT(Long.MIN_VALUE + 1);
        assertEquals(-9223372036854775807L, value.getLong());
    }

    @Test
    void testTrillionValue() {
        PlcLINT value = new PlcLINT(1000000000000L);
        assertEquals(1000000000000L, value.getLong());
    }

    @Test
    void testNegativeTrillionValue() {
        PlcLINT value = new PlcLINT(-1000000000000L);
        assertEquals(-1000000000000L, value.getLong());
    }

    @Test
    void testIntegerMaxAsLong() {
        PlcLINT value = new PlcLINT((long) Integer.MAX_VALUE);
        assertEquals(2147483647L, value.getLong());
    }

    @Test
    void testBeyondIntegerMax() {
        PlcLINT value = new PlcLINT(5000000000L);
        assertEquals(5000000000L, value.getLong());
        assertTrue(value.getLong() > Integer.MAX_VALUE);
    }
}
