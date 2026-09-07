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
 * Test class for PlcUDINT - Unsigned 32-bit integer (0 to 4294967295)
 */
class PlcUDINTTest {

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcUDINT value = new PlcUDINT(0L);
        assertEquals(0L, value.getLong());
        assertEquals(PlcValueType.UDINT, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcUDINT value = new PlcUDINT(4294967295L);
        assertEquals(4294967295L, value.getLong());
    }

    @Test
    void testMinValueMinusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(-1L));
    }

    @Test
    void testMaxValuePlusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(4294967296L));
    }

    @Test
    void testMinValueFromInteger() {
        PlcUDINT value = new PlcUDINT(0);
        assertEquals(0L, value.getLong());
    }

    @Test
    void testMaxValueFromBigInteger() {
        PlcUDINT value = new PlcUDINT(BigInteger.valueOf(4294967295L));
        assertEquals(4294967295L, value.getLong());
    }

    @Test
    void testBelowRangeFromBigInteger() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(BigInteger.valueOf(-1)));
    }

    @Test
    void testAboveRangeFromBigInteger() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(BigInteger.valueOf(4294967296L)));
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcUDINT trueValue = new PlcUDINT(true);
        assertEquals(1L, trueValue.getLong());

        PlcUDINT falseValue = new PlcUDINT(false);
        assertEquals(0L, falseValue.getLong());
    }

    @Test
    void testByteConstructor() {
        PlcUDINT value = new PlcUDINT((byte) 127);
        assertEquals(127L, value.getLong());
    }

    @Test
    void testByteConstructorNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT((byte) -1));
    }

    @Test
    void testShortConstructorValid() {
        PlcUDINT value = new PlcUDINT((short) 30000);
        assertEquals(30000L, value.getLong());
    }

    @Test
    void testShortConstructorNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT((short) -1));
    }

    @Test
    void testIntegerConstructorValid() {
        PlcUDINT value = new PlcUDINT(2000000000);
        assertEquals(2000000000L, value.getLong());
    }

    @Test
    void testIntegerConstructorNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(-1));
    }

    @Test
    void testFloatConstructorValid() {
        PlcUDINT value = new PlcUDINT(1000000.5f);
        assertEquals(1000000L, value.getLong());
    }

    @Test
    void testFloatConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(5000000000.0f));
    }

    @Test
    void testDoubleConstructorValid() {
        PlcUDINT value = new PlcUDINT(3000000000.9);
        assertEquals(3000000000L, value.getLong());
    }

    @Test
    void testDoubleConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(-1.0));
    }

    @Test
    void testBigDecimalConstructorValid() {
        PlcUDINT value = new PlcUDINT(BigDecimal.valueOf(2000000000L));
        assertEquals(2000000000L, value.getLong());
    }

    @Test
    void testBigDecimalConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(BigDecimal.valueOf(-1)));
    }

    @Test
    void testBigDecimalConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT(BigDecimal.valueOf(4294967296L)));
    }

    @Test
    void testStringConstructorValid() {
        PlcUDINT value = new PlcUDINT("3000000000");
        assertEquals(3000000000L, value.getLong());
    }

    @Test
    void testStringConstructorMinValue() {
        PlcUDINT value = new PlcUDINT("0");
        assertEquals(0L, value.getLong());
    }

    @Test
    void testStringConstructorMaxValue() {
        PlcUDINT value = new PlcUDINT("4294967295");
        assertEquals(4294967295L, value.getLong());
    }

    @Test
    void testStringConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT("-1"));
    }

    @Test
    void testStringConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT("4294967296"));
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUDINT("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcUDINT value = new PlcUDINT("  2000000000  ");
        assertEquals(2000000000L, value.getLong());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcUDINT() {
        PlcUDINT original = new PlcUDINT(3000000000L);
        PlcUDINT copy = PlcUDINT.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithLong() {
        PlcUDINT value = PlcUDINT.of(4000000000L);
        assertEquals(4000000000L, value.getLong());
    }

    @Test
    void testOfMethodWithString() {
        PlcUDINT value = PlcUDINT.of("3500000000");
        assertEquals(3500000000L, value.getLong());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcUDINT zero = new PlcUDINT(0L);
        assertFalse(zero.getBoolean());

        PlcUDINT nonZero = new PlcUDINT(1L);
        assertTrue(nonZero.getBoolean());

        PlcUDINT large = new PlcUDINT(4294967295L);
        assertTrue(large.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcUDINT value = new PlcUDINT(127L);
        assertEquals(127, value.getByte());
    }

    @Test
    void testIsByte() {
        PlcUDINT smallValue = new PlcUDINT(100L);
        assertTrue(smallValue.isByte());

        PlcUDINT largeValue = new PlcUDINT(300L);
        assertFalse(largeValue.isByte());
    }

    @Test
    void testGetShort() {
        PlcUDINT value = new PlcUDINT(30000L);
        assertEquals(30000, value.getShort());
    }

    @Test
    void testIsShort() {
        PlcUDINT smallValue = new PlcUDINT(30000L);
        assertTrue(smallValue.isShort());

        PlcUDINT largeValue = new PlcUDINT(40000L);
        assertFalse(largeValue.isShort());
    }

    @Test
    void testGetInteger() {
        PlcUDINT value = new PlcUDINT(2000000000L);
        assertEquals(2000000000, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcUDINT value = new PlcUDINT(4294967295L);
        assertEquals(4294967295L, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcUDINT value = new PlcUDINT(4294967295L);
        assertEquals(BigInteger.valueOf(4294967295L), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcUDINT value = new PlcUDINT(1000000L);
        assertEquals(1000000.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcUDINT value = new PlcUDINT(3000000000L);
        assertEquals(3000000000.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcUDINT value = new PlcUDINT(2000000000L);
        assertEquals(BigDecimal.valueOf(2000000000.0), value.getBigDecimal());
    }

    @Test
    void testGetString() {
        PlcUDINT value = new PlcUDINT(123456789L);
        assertEquals("123456789", value.getString());
    }

    @Test
    void testToString() {
        PlcUDINT value = new PlcUDINT(4294967295L);
        assertEquals("4294967295", value.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcUDINT value = new PlcUDINT(1L);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsInteger() {
        PlcUDINT value = new PlcUDINT(2000000000L);
        assertTrue(value.isInteger());
    }

    @Test
    void testIsLong() {
        PlcUDINT value = new PlcUDINT(4294967295L);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcUDINT value = new PlcUDINT(3000000000L);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcUDINT value = new PlcUDINT(1000000L);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcUDINT value = new PlcUDINT(2000000000L);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcUDINT value = new PlcUDINT(4294967295L);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcUDINT value = new PlcUDINT(123456789L);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesMinValue() {
        PlcUDINT value = new PlcUDINT(0L);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        assertEquals(0, bytes[0]);
        assertEquals(0, bytes[1]);
        assertEquals(0, bytes[2]);
        assertEquals(0, bytes[3]);
    }

    @Test
    void testGetBytesMaxValue() {
        PlcUDINT value = new PlcUDINT(4294967295L);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        assertEquals((byte) 0xFF, bytes[0]);
        assertEquals((byte) 0xFF, bytes[1]);
        assertEquals((byte) 0xFF, bytes[2]);
        assertEquals((byte) 0xFF, bytes[3]);
    }

    @Test
    void testGetBytesMidValue() {
        PlcUDINT value = new PlcUDINT(0x12345678L);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        assertEquals((byte) 0x12, bytes[0]);
        assertEquals((byte) 0x34, bytes[1]);
        assertEquals((byte) 0x56, bytes[2]);
        assertEquals((byte) 0x78, bytes[3]);
    }

    @Test
    void testGetRaw() {
        PlcUDINT value = new PlcUDINT(3000000000L);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    // ========== Edge Cases ==========

    @Test
    void testMidRangeValue() {
        PlcUDINT value = new PlcUDINT(2147483647L); // Integer.MAX_VALUE
        assertEquals(2147483647L, value.getLong());
    }

    @Test
    void testPowerOfTwo() {
        PlcUDINT value = new PlcUDINT(2147483648L); // Integer.MAX_VALUE + 1
        assertEquals(2147483648L, value.getLong());
    }

    @Test
    void testValueJustBelowMax() {
        PlcUDINT value = new PlcUDINT(4294967294L);
        assertEquals(4294967294L, value.getLong());
    }

    @Test
    void testValueJustAboveMin() {
        PlcUDINT value = new PlcUDINT(1L);
        assertEquals(1L, value.getLong());
    }

    @Test
    void testByteOrderInSerialization() {
        // Test that byte order is big-endian
        PlcUDINT value = new PlcUDINT(0x01020304L);
        byte[] bytes = value.getBytes();
        assertEquals((byte) 0x01, bytes[0]); // High byte first
        assertEquals((byte) 0x02, bytes[1]);
        assertEquals((byte) 0x03, bytes[2]);
        assertEquals((byte) 0x04, bytes[3]); // Low byte last
    }

    @Test
    void testUpperHalfRange() {
        // Test values in the upper half (> Integer.MAX_VALUE)
        PlcUDINT value = new PlcUDINT(3000000000L);
        assertEquals(3000000000L, value.getLong());
        assertTrue(value.getLong() > Integer.MAX_VALUE);
    }
}
