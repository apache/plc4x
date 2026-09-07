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
 * Test class for PlcINT - Signed 16-bit integer (-32768 to 32767)
 */
class PlcINTTest {

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcINT value = new PlcINT((short) -32768);
        assertEquals(-32768, value.getShort());
        assertEquals(PlcValueType.INT, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcINT value = new PlcINT((short) 32767);
        assertEquals(32767, value.getShort());
    }

    @Test
    void testMinValueMinusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(-32769));
    }

    @Test
    void testMaxValuePlusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(32768));
    }

    @Test
    void testZero() {
        PlcINT value = new PlcINT((short) 0);
        assertEquals(0, value.getShort());
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcINT trueValue = new PlcINT(true);
        assertEquals(1, trueValue.getShort());

        PlcINT falseValue = new PlcINT(false);
        assertEquals(0, falseValue.getShort());
    }

    @Test
    void testByteConstructor() {
        PlcINT value = new PlcINT((byte) -50);
        assertEquals(-50, value.getShort());
    }

    @Test
    void testShortConstructor() {
        PlcINT value = new PlcINT((short) -30000);
        assertEquals(-30000, value.getShort());
    }

    @Test
    void testIntegerConstructorValid() {
        PlcINT value = new PlcINT(20000);
        assertEquals(20000, value.getShort());
    }

    @Test
    void testIntegerConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(-40000));
    }

    @Test
    void testIntegerConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(40000));
    }

    @Test
    void testLongConstructorValid() {
        PlcINT value = new PlcINT(-15000L);
        assertEquals(-15000, value.getShort());
    }

    @Test
    void testLongConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(50000L));
    }

    @Test
    void testFloatConstructorValid() {
        PlcINT value = new PlcINT(10000.5f);
        assertEquals(10000, value.getShort());
    }

    @Test
    void testFloatConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(-40000.0f));
    }

    @Test
    void testDoubleConstructorValid() {
        PlcINT value = new PlcINT(-20000.9);
        assertEquals(-20000, value.getShort());
    }

    @Test
    void testDoubleConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(35000.0));
    }

    @Test
    void testBigIntegerConstructorValid() {
        PlcINT value = new PlcINT(BigInteger.valueOf(-32768));
        assertEquals(-32768, value.getShort());
    }

    @Test
    void testBigIntegerConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(BigInteger.valueOf(-32769)));
    }

    @Test
    void testBigIntegerConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(BigInteger.valueOf(32768)));
    }

    @Test
    void testBigDecimalConstructorValid() {
        PlcINT value = new PlcINT(BigDecimal.valueOf(15000));
        assertEquals(15000, value.getShort());
    }

    @Test
    void testBigDecimalConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(BigDecimal.valueOf(-33000)));
    }

    @Test
    void testBigDecimalConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT(BigDecimal.valueOf(35000)));
    }

    @Test
    void testStringConstructorValid() {
        PlcINT value = new PlcINT("-15000");
        assertEquals(-15000, value.getShort());
    }

    @Test
    void testStringConstructorMinValue() {
        PlcINT value = new PlcINT("-32768");
        assertEquals(-32768, value.getShort());
    }

    @Test
    void testStringConstructorMaxValue() {
        PlcINT value = new PlcINT("32767");
        assertEquals(32767, value.getShort());
    }

    @Test
    void testStringConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT("-32769"));
    }

    @Test
    void testStringConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT("32768"));
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcINT("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcINT value = new PlcINT("  -10000  ");
        assertEquals(-10000, value.getShort());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcINT() {
        PlcINT original = new PlcINT((short) -20000);
        PlcINT copy = PlcINT.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithInteger() {
        PlcINT value = PlcINT.of(15000);
        assertEquals(15000, value.getShort());
    }

    @Test
    void testOfMethodWithString() {
        PlcINT value = PlcINT.of("-25000");
        assertEquals(-25000, value.getShort());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcINT zero = new PlcINT((short) 0);
        assertFalse(zero.getBoolean());

        PlcINT positive = new PlcINT((short) 1);
        assertTrue(positive.getBoolean());

        PlcINT negative = new PlcINT((short) -1);
        assertTrue(negative.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcINT value = new PlcINT((short) 100);
        assertEquals(100, value.getByte());
    }

    @Test
    void testGetShort() {
        PlcINT value = new PlcINT((short) -32768);
        assertEquals(-32768, value.getShort());
    }

    @Test
    void testGetInteger() {
        PlcINT value = new PlcINT((short) 20000);
        assertEquals(20000, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcINT value = new PlcINT((short) -15000);
        assertEquals(-15000L, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcINT value = new PlcINT((short) 32767);
        assertEquals(BigInteger.valueOf(32767), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcINT value = new PlcINT((short) -10000);
        assertEquals(-10000.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcINT value = new PlcINT((short) 25000);
        assertEquals(25000.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcINT value = new PlcINT((short) -20000);
        assertEquals(BigDecimal.valueOf(-20000.0f), value.getBigDecimal());
    }

    @Test
    void testGetString() {
        PlcINT value = new PlcINT((short) -15000);
        assertEquals("-15000", value.getString());
    }

    @Test
    void testToString() {
        PlcINT value = new PlcINT((short) 32767);
        assertEquals("32767", value.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcINT value = new PlcINT((short) 1);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcINT smallValue = new PlcINT((short) 100);
        assertTrue(smallValue.isByte());

        PlcINT largeValue = new PlcINT((short) 300);
        assertFalse(largeValue.isByte());
    }

    @Test
    void testIsShort() {
        PlcINT value = new PlcINT((short) -30000);
        assertTrue(value.isShort());
    }

    @Test
    void testIsInteger() {
        PlcINT value = new PlcINT((short) 20000);
        assertTrue(value.isInteger());
    }

    @Test
    void testIsLong() {
        PlcINT value = new PlcINT((short) -25000);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcINT value = new PlcINT((short) 32767);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcINT value = new PlcINT((short) 15000);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcINT value = new PlcINT((short) -10000);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcINT value = new PlcINT((short) 32767);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcINT value = new PlcINT((short) -15000);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesMinValue() {
        PlcINT value = new PlcINT((short) -32768);
        byte[] bytes = value.getBytes();
        assertEquals(2, bytes.length);
        assertEquals((byte) 0x80, bytes[0]);
        assertEquals((byte) 0x00, bytes[1]);
    }

    @Test
    void testGetBytesMaxValue() {
        PlcINT value = new PlcINT((short) 32767);
        byte[] bytes = value.getBytes();
        assertEquals(2, bytes.length);
        assertEquals((byte) 0x7F, bytes[0]);
        assertEquals((byte) 0xFF, bytes[1]);
    }

    @Test
    void testGetBytesZero() {
        PlcINT value = new PlcINT((short) 0);
        byte[] bytes = value.getBytes();
        assertEquals(2, bytes.length);
        assertEquals(0, bytes[0]);
        assertEquals(0, bytes[1]);
    }

    @Test
    void testGetBytesNegativeValue() {
        PlcINT value = new PlcINT((short) -1);
        byte[] bytes = value.getBytes();
        assertEquals(2, bytes.length);
        assertEquals((byte) 0xFF, bytes[0]);
        assertEquals((byte) 0xFF, bytes[1]);
    }

    @Test
    void testGetRaw() {
        PlcINT value = new PlcINT((short) -15000);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    @Test
    void testByteOrderInSerialization() {
        // Test that byte order is big-endian
        PlcINT value = new PlcINT((short) 0x0102);
        byte[] bytes = value.getBytes();
        assertEquals((byte) 0x01, bytes[0]); // High byte first
        assertEquals((byte) 0x02, bytes[1]); // Low byte second
    }

    // ========== Edge Cases ==========

    @Test
    void testNegativeOneValue() {
        PlcINT value = new PlcINT((short) -1);
        assertEquals(-1, value.getShort());
    }

    @Test
    void testPositiveOneValue() {
        PlcINT value = new PlcINT((short) 1);
        assertEquals(1, value.getShort());
    }

    @Test
    void testValueJustBelowMax() {
        PlcINT value = new PlcINT((short) 32766);
        assertEquals(32766, value.getShort());
    }

    @Test
    void testValueJustAboveMin() {
        PlcINT value = new PlcINT((short) -32767);
        assertEquals(-32767, value.getShort());
    }
}
