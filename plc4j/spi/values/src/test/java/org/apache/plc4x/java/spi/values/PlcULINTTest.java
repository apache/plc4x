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
 * Test class for PlcULINT - Unsigned 64-bit integer (0 to 2^64-1)
 */
class PlcULINTTest {

    private static final BigInteger MAX_ULINT = BigInteger.ONE.shiftLeft(64).subtract(BigInteger.ONE);

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcULINT value = new PlcULINT(BigInteger.ZERO);
        assertEquals(BigInteger.ZERO, value.getBigInteger());
        assertEquals(PlcValueType.ULINT, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcULINT value = new PlcULINT(MAX_ULINT);
        assertEquals(MAX_ULINT, value.getBigInteger());
    }

    @Test
    void testMinValueMinusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(BigInteger.valueOf(-1)));
    }

    @Test
    void testMaxValuePlusOne() {
        BigInteger overMax = MAX_ULINT.add(BigInteger.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(overMax));
    }

    @Test
    void testMinValueFromLong() {
        PlcULINT value = new PlcULINT(0L);
        assertEquals(BigInteger.ZERO, value.getBigInteger());
    }

    @Test
    void testBelowRangeFromLong() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(-1L));
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcULINT trueValue = new PlcULINT(true);
        assertEquals(BigInteger.ONE, trueValue.getBigInteger());

        PlcULINT falseValue = new PlcULINT(false);
        assertEquals(BigInteger.ZERO, falseValue.getBigInteger());
    }

    @Test
    void testByteConstructor() {
        PlcULINT value = new PlcULINT((byte) 127);
        assertEquals(BigInteger.valueOf(127), value.getBigInteger());
    }

    @Test
    void testByteConstructorNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT((byte) -1));
    }

    @Test
    void testShortConstructorValid() {
        PlcULINT value = new PlcULINT((short) 30000);
        assertEquals(BigInteger.valueOf(30000), value.getBigInteger());
    }

    @Test
    void testShortConstructorNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT((short) -1));
    }

    @Test
    void testIntegerConstructorValid() {
        PlcULINT value = new PlcULINT(2000000000);
        assertEquals(BigInteger.valueOf(2000000000), value.getBigInteger());
    }

    @Test
    void testIntegerConstructorNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(-1));
    }

    @Test
    void testLongConstructorValid() {
        PlcULINT value = new PlcULINT(Long.MAX_VALUE);
        assertEquals(BigInteger.valueOf(Long.MAX_VALUE), value.getBigInteger());
    }

    @Test
    void testLongConstructorNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(-1L));
    }

    @Test
    void testFloatConstructorValid() {
        PlcULINT value = new PlcULINT(1000000.5f);
        assertEquals(BigInteger.valueOf(1000000), value.getBigInteger());
    }

    @Test
    void testFloatConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(-100.0f));
    }

    @Test
    void testDoubleConstructorValid() {
        PlcULINT value = new PlcULINT(3000000000.9);
        assertEquals(BigInteger.valueOf(3000000000L), value.getBigInteger());
    }

    @Test
    void testDoubleConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(-1.0));
    }

    @Test
    void testBigIntegerConstructorValid() {
        BigInteger large = BigInteger.valueOf(Long.MAX_VALUE).add(BigInteger.valueOf(1000));
        PlcULINT value = new PlcULINT(large);
        assertEquals(large, value.getBigInteger());
    }

    @Test
    void testBigIntegerConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(BigInteger.valueOf(-1)));
    }

    @Test
    void testBigIntegerConstructorAboveRange() {
        BigInteger tooLarge = MAX_ULINT.add(BigInteger.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(tooLarge));
    }

    @Test
    void testBigDecimalConstructorValid() {
        PlcULINT value = new PlcULINT(BigDecimal.valueOf(2000000000L));
        assertEquals(BigInteger.valueOf(2000000000L), value.getBigInteger());
    }

    @Test
    void testBigDecimalConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(BigDecimal.valueOf(-1)));
    }

    @Test
    void testBigDecimalConstructorAboveRange() {
        BigDecimal tooLarge = new BigDecimal(MAX_ULINT).add(BigDecimal.ONE);
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(tooLarge));
    }

    @Test
    void testStringConstructorValid() {
        PlcULINT value = new PlcULINT("3000000000");
        assertEquals(BigInteger.valueOf(3000000000L), value.getBigInteger());
    }

    @Test
    void testStringConstructorMinValue() {
        PlcULINT value = new PlcULINT("0");
        assertEquals(BigInteger.ZERO, value.getBigInteger());
    }

    @Test
    void testStringConstructorMaxValue() {
        PlcULINT value = new PlcULINT(MAX_ULINT.toString());
        assertEquals(MAX_ULINT, value.getBigInteger());
    }

    @Test
    void testStringConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT("-1"));
    }

    @Test
    void testStringConstructorAboveRange() {
        String tooLarge = MAX_ULINT.add(BigInteger.ONE).toString();
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT(tooLarge));
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcULINT("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcULINT value = new PlcULINT("  2000000000  ");
        assertEquals(BigInteger.valueOf(2000000000L), value.getBigInteger());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcULINT() {
        PlcULINT original = new PlcULINT(BigInteger.valueOf(3000000000L));
        PlcULINT copy = PlcULINT.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithLong() {
        PlcULINT value = PlcULINT.of(4000000000L);
        assertEquals(BigInteger.valueOf(4000000000L), value.getBigInteger());
    }

    @Test
    void testOfMethodWithString() {
        PlcULINT value = PlcULINT.of("3500000000");
        assertEquals(BigInteger.valueOf(3500000000L), value.getBigInteger());
    }

    @Test
    void testOfMethodWithBigInteger() {
        BigInteger large = BigInteger.valueOf(Long.MAX_VALUE).add(BigInteger.valueOf(1000));
        PlcULINT value = PlcULINT.of(large);
        assertEquals(large, value.getBigInteger());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcULINT zero = new PlcULINT(BigInteger.ZERO);
        assertFalse(zero.getBoolean());

        PlcULINT nonZero = new PlcULINT(BigInteger.ONE);
        assertTrue(nonZero.getBoolean());

        PlcULINT large = new PlcULINT(MAX_ULINT);
        assertTrue(large.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(127));
        assertEquals(127, value.getByte());
    }

    @Test
    void testIsByte() {
        PlcULINT smallValue = new PlcULINT(BigInteger.valueOf(100));
        assertTrue(smallValue.isByte());

        PlcULINT largeValue = new PlcULINT(BigInteger.valueOf(300));
        assertFalse(largeValue.isByte());
    }

    @Test
    void testGetShort() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(30000));
        assertEquals(30000, value.getShort());
    }

    @Test
    void testIsShort() {
        PlcULINT smallValue = new PlcULINT(BigInteger.valueOf(30000));
        assertTrue(smallValue.isShort());

        PlcULINT largeValue = new PlcULINT(BigInteger.valueOf(40000));
        assertFalse(largeValue.isShort());
    }

    @Test
    void testGetInteger() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(2000000000));
        assertEquals(2000000000, value.getInteger());
    }

    @Test
    void testIsInteger() {
        PlcULINT smallValue = new PlcULINT(BigInteger.valueOf(2000000000));
        assertTrue(smallValue.isInteger());

        PlcULINT largeValue = new PlcULINT(BigInteger.valueOf(3000000000L));
        assertFalse(largeValue.isInteger());
    }

    @Test
    void testGetLong() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(Long.MAX_VALUE));
        assertEquals(Long.MAX_VALUE, value.getLong());
    }

    @Test
    void testIsLong() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(Long.MAX_VALUE));
        assertTrue(value.isLong());
    }

    @Test
    void testGetBigInteger() {
        BigInteger large = BigInteger.valueOf(Long.MAX_VALUE).add(BigInteger.valueOf(1000));
        PlcULINT value = new PlcULINT(large);
        assertEquals(large, value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(1000000));
        assertEquals(1000000.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(3000000000L));
        assertEquals(3000000000.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        BigInteger large = BigInteger.valueOf(Long.MAX_VALUE).add(BigInteger.valueOf(1000));
        PlcULINT value = new PlcULINT(large);
        assertEquals(new BigDecimal(large), value.getBigDecimal());
    }

    @Test
    void testGetString() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(123456789));
        assertEquals("123456789", value.getString());
    }

    @Test
    void testToString() {
        PlcULINT value = new PlcULINT(MAX_ULINT);
        assertEquals(MAX_ULINT.toString(), value.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcULINT value = new PlcULINT(BigInteger.ONE);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsBigInteger() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(3000000000L));
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(1000000));
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(2000000000));
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcULINT value = new PlcULINT(MAX_ULINT);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(123456789));
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesMinValue() {
        PlcULINT value = new PlcULINT(BigInteger.ZERO);
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        for (int i = 0; i < 8; i++) {
            assertEquals(0, bytes[i], "Byte at position " + i + " should be 0");
        }
    }

    @Test
    void testGetBytesMaxValue() {
        PlcULINT value = new PlcULINT(MAX_ULINT);
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        for (int i = 0; i < 8; i++) {
            assertEquals((byte) 0xFF, bytes[i], "Byte at position " + i + " should be 0xFF");
        }
    }

    @Test
    void testGetBytesMidValue() {
        PlcULINT value = new PlcULINT(new BigInteger("1234567890ABCDEF", 16));
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        assertEquals((byte) 0x12, bytes[0]);
        assertEquals((byte) 0x34, bytes[1]);
        assertEquals((byte) 0x56, bytes[2]);
        assertEquals((byte) 0x78, bytes[3]);
        assertEquals((byte) 0x90, bytes[4]);
        assertEquals((byte) 0xAB, bytes[5]);
        assertEquals((byte) 0xCD, bytes[6]);
        assertEquals((byte) 0xEF, bytes[7]);
    }

    @Test
    void testGetRaw() {
        PlcULINT value = new PlcULINT(BigInteger.valueOf(3000000000L));
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    // ========== Edge Cases ==========

    @Test
    void testLongMaxValue() {
        PlcULINT value = new PlcULINT(Long.MAX_VALUE);
        assertEquals(BigInteger.valueOf(Long.MAX_VALUE), value.getBigInteger());
    }

    @Test
    void testLongMaxValuePlusOne() {
        BigInteger longMaxPlusOne = BigInteger.valueOf(Long.MAX_VALUE).add(BigInteger.ONE);
        PlcULINT value = new PlcULINT(longMaxPlusOne);
        assertEquals(longMaxPlusOne, value.getBigInteger());
    }

    @Test
    void testValueJustBelowMax() {
        BigInteger justBelowMax = MAX_ULINT.subtract(BigInteger.ONE);
        PlcULINT value = new PlcULINT(justBelowMax);
        assertEquals(justBelowMax, value.getBigInteger());
    }

    @Test
    void testValueJustAboveMin() {
        PlcULINT value = new PlcULINT(BigInteger.ONE);
        assertEquals(BigInteger.ONE, value.getBigInteger());
    }

    @Test
    void testByteOrderInSerialization() {
        // Test that byte order is big-endian
        PlcULINT value = new PlcULINT(new BigInteger("0102030405060708", 16));
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

    @Test
    void testUpperRange() {
        // Test values significantly above Long.MAX_VALUE
        BigInteger large = BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.TWO);
        PlcULINT value = new PlcULINT(large);
        assertEquals(large, value.getBigInteger());
        assertTrue(value.getBigInteger().compareTo(BigInteger.valueOf(Long.MAX_VALUE)) > 0);
    }

    @Test
    void testMaxValueCalculation() {
        // Verify that MAX_ULINT is actually 2^64 - 1
        BigInteger expected = BigInteger.ONE.shiftLeft(64).subtract(BigInteger.ONE);
        assertEquals(expected, MAX_ULINT);
        assertEquals("18446744073709551615", MAX_ULINT.toString());
    }

    @Test
    void testPowerOfTwo() {
        // Test 2^63
        BigInteger powerOf63 = BigInteger.ONE.shiftLeft(63);
        PlcULINT value = new PlcULINT(powerOf63);
        assertEquals(powerOf63, value.getBigInteger());
    }
}
