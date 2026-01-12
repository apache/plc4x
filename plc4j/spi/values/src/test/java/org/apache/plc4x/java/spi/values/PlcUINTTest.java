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
 * Test class for PlcUINT - Unsigned 16-bit integer (0 to 65535)
 */
class PlcUINTTest {

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcUINT value = new PlcUINT(0);
        assertEquals(0, value.getInteger());
        assertEquals(PlcValueType.UINT, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcUINT value = new PlcUINT(65535);
        assertEquals(65535, value.getInteger());
    }

    @Test
    void testMinValueMinusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(-1));
    }

    @Test
    void testMaxValuePlusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(65536));
    }

    @Test
    void testMinValueFromShort() {
        PlcUINT value = new PlcUINT((short) 0);
        assertEquals(0, value.getInteger());
    }

    @Test
    void testMaxValueFromLong() {
        PlcUINT value = new PlcUINT(65535L);
        assertEquals(65535, value.getInteger());
    }

    @Test
    void testBelowRangeFromLong() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(-1L));
    }

    @Test
    void testAboveRangeFromLong() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(65536L));
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcUINT trueValue = new PlcUINT(true);
        assertEquals(1, trueValue.getInteger());

        PlcUINT falseValue = new PlcUINT(false);
        assertEquals(0, falseValue.getInteger());
    }

    @Test
    void testByteConstructor() {
        PlcUINT value = new PlcUINT((byte) 127);
        assertEquals(127, value.getInteger());
    }

    @Test
    void testShortConstructorValid() {
        PlcUINT value = new PlcUINT((short) 30000);
        assertEquals(30000, value.getInteger());
    }

    @Test
    void testShortConstructorNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT((short) -1));
    }

    @Test
    void testFloatConstructorValid() {
        PlcUINT value = new PlcUINT(32000.5f);
        assertEquals(32000, value.getInteger());
    }

    @Test
    void testFloatConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(70000.0f));
    }

    @Test
    void testDoubleConstructorValid() {
        PlcUINT value = new PlcUINT(50000.9);
        assertEquals(50000, value.getInteger());
    }

    @Test
    void testDoubleConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(-1.0));
    }

    @Test
    void testBigIntegerConstructorValid() {
        PlcUINT value = new PlcUINT(BigInteger.valueOf(65535));
        assertEquals(65535, value.getInteger());
    }

    @Test
    void testBigIntegerConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(BigInteger.valueOf(-1)));
    }

    @Test
    void testBigIntegerConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(BigInteger.valueOf(65536)));
    }

    @Test
    void testBigDecimalConstructorValid() {
        PlcUINT value = new PlcUINT(BigDecimal.valueOf(40000));
        assertEquals(40000, value.getInteger());
    }

    @Test
    void testBigDecimalConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(BigDecimal.valueOf(-1)));
    }

    @Test
    void testBigDecimalConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT(BigDecimal.valueOf(65536)));
    }

    @Test
    void testStringConstructorValid() {
        PlcUINT value = new PlcUINT("30000");
        assertEquals(30000, value.getInteger());
    }

    @Test
    void testStringConstructorMinValue() {
        PlcUINT value = new PlcUINT("0");
        assertEquals(0, value.getInteger());
    }

    @Test
    void testStringConstructorMaxValue() {
        PlcUINT value = new PlcUINT("65535");
        assertEquals(65535, value.getInteger());
    }

    @Test
    void testStringConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT("-1"));
    }

    @Test
    void testStringConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT("65536"));
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcUINT("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcUINT value = new PlcUINT("  32000  ");
        assertEquals(32000, value.getInteger());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcUINT() {
        PlcUINT original = new PlcUINT(40000);
        PlcUINT copy = PlcUINT.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithInteger() {
        PlcUINT value = PlcUINT.of(50000);
        assertEquals(50000, value.getInteger());
    }

    @Test
    void testOfMethodWithString() {
        PlcUINT value = PlcUINT.of("60000");
        assertEquals(60000, value.getInteger());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcUINT zero = new PlcUINT(0);
        assertFalse(zero.getBoolean());

        PlcUINT nonZero = new PlcUINT(1);
        assertTrue(nonZero.getBoolean());

        PlcUINT large = new PlcUINT(65535);
        assertTrue(large.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcUINT value = new PlcUINT(127);
        assertEquals(127, value.getByte());
    }

    @Test
    void testGetShort() {
        PlcUINT value = new PlcUINT(32000);
        assertEquals(32000, value.getShort());
    }

    @Test
    void testGetInteger() {
        PlcUINT value = new PlcUINT(65535);
        assertEquals(65535, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcUINT value = new PlcUINT(50000);
        assertEquals(50000L, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcUINT value = new PlcUINT(65535);
        assertEquals(BigInteger.valueOf(65535), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcUINT value = new PlcUINT(30000);
        assertEquals(30000.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcUINT value = new PlcUINT(40000);
        assertEquals(40000.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcUINT value = new PlcUINT(50000);
        assertEquals(BigDecimal.valueOf(50000.0f), value.getBigDecimal());
    }

    @Test
    void testGetString() {
        PlcUINT value = new PlcUINT(12345);
        assertEquals("12345", value.getString());
    }

    @Test
    void testToString() {
        PlcUINT value = new PlcUINT(65535);
        assertEquals("65535", value.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcUINT value = new PlcUINT(1);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcUINT smallValue = new PlcUINT(100);
        assertTrue(smallValue.isByte());

        PlcUINT largeValue = new PlcUINT(300);
        assertFalse(largeValue.isByte());
    }

    @Test
    void testIsShort() {
        PlcUINT value = new PlcUINT(30000);
        assertTrue(value.isShort());
    }

    @Test
    void testIsInteger() {
        PlcUINT value = new PlcUINT(65535);
        assertTrue(value.isInteger());
    }

    @Test
    void testIsLong() {
        PlcUINT value = new PlcUINT(50000);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcUINT value = new PlcUINT(40000);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcUINT value = new PlcUINT(32000);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcUINT value = new PlcUINT(20000);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcUINT value = new PlcUINT(65535);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcUINT value = new PlcUINT(12345);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesMinValue() {
        PlcUINT value = new PlcUINT(0);
        byte[] bytes = value.getBytes();
        assertEquals(2, bytes.length);
        assertEquals(0, bytes[0]);
        assertEquals(0, bytes[1]);
    }

    @Test
    void testGetBytesMaxValue() {
        PlcUINT value = new PlcUINT(65535);
        byte[] bytes = value.getBytes();
        assertEquals(2, bytes.length);
        assertEquals((byte) 0xFF, bytes[0]);
        assertEquals((byte) 0xFF, bytes[1]);
    }

    @Test
    void testGetBytesMidValue() {
        PlcUINT value = new PlcUINT(0x1234);
        byte[] bytes = value.getBytes();
        assertEquals(2, bytes.length);
        assertEquals((byte) 0x12, bytes[0]);
        assertEquals((byte) 0x34, bytes[1]);
    }

    @Test
    void testGetRaw() {
        PlcUINT value = new PlcUINT(30000);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    // ========== Edge Cases ==========

    @Test
    void testMidRangeValue() {
        PlcUINT value = new PlcUINT(32767);
        assertEquals(32767, value.getInteger());
    }

    @Test
    void testPowerOfTwo() {
        PlcUINT value = new PlcUINT(32768);
        assertEquals(32768, value.getInteger());
    }

    @Test
    void testValueJustBelowMax() {
        PlcUINT value = new PlcUINT(65534);
        assertEquals(65534, value.getInteger());
    }

    @Test
    void testValueJustAboveMin() {
        PlcUINT value = new PlcUINT(1);
        assertEquals(1, value.getInteger());
    }

    @Test
    void testByteOrderInSerialization() {
        // Test that byte order is big-endian
        PlcUINT value = new PlcUINT(0x0102);
        byte[] bytes = value.getBytes();
        assertEquals((byte) 0x01, bytes[0]); // High byte first
        assertEquals((byte) 0x02, bytes[1]); // Low byte second
    }
}
