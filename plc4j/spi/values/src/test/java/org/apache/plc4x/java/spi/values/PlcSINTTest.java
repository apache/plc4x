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
 * Test class for PlcSINT - Signed 8-bit integer (-128 to 127)
 */
class PlcSINTTest {

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcSINT value = new PlcSINT((byte) -128);
        assertEquals(-128, value.getByte());
        assertEquals(PlcValueType.SINT, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcSINT value = new PlcSINT((byte) 127);
        assertEquals(127, value.getByte());
    }

    @Test
    void testMinValueMinusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT((short) -129));
    }

    @Test
    void testMaxValuePlusOne() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT((short) 128));
    }

    @Test
    void testZero() {
        PlcSINT value = new PlcSINT((byte) 0);
        assertEquals(0, value.getByte());
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcSINT trueValue = new PlcSINT(true);
        assertEquals(1, trueValue.getByte());

        PlcSINT falseValue = new PlcSINT(false);
        assertEquals(0, falseValue.getByte());
    }

    @Test
    void testByteConstructor() {
        PlcSINT value = new PlcSINT((byte) -50);
        assertEquals(-50, value.getByte());
    }

    @Test
    void testShortConstructorValid() {
        PlcSINT value = new PlcSINT((short) 100);
        assertEquals(100, value.getByte());
    }

    @Test
    void testShortConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT((short) -200));
    }

    @Test
    void testShortConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT((short) 200));
    }

    @Test
    void testIntegerConstructorValid() {
        PlcSINT value = new PlcSINT(-100);
        assertEquals(-100, value.getByte());
    }

    @Test
    void testIntegerConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT(300));
    }

    @Test
    void testLongConstructorValid() {
        PlcSINT value = new PlcSINT(50L);
        assertEquals(50, value.getByte());
    }

    @Test
    void testLongConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT(-200L));
    }

    @Test
    void testFloatConstructorValid() {
        PlcSINT value = new PlcSINT(75.5f);
        assertEquals(75, value.getByte());
    }

    @Test
    void testFloatConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT(150.0f));
    }

    @Test
    void testDoubleConstructorValid() {
        PlcSINT value = new PlcSINT(-75.9);
        assertEquals(-75, value.getByte());
    }

    @Test
    void testDoubleConstructorInvalid() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT(-200.0));
    }

    @Test
    void testBigIntegerConstructorValid() {
        PlcSINT value = new PlcSINT(BigInteger.valueOf(-128));
        assertEquals(-128, value.getByte());
    }

    @Test
    void testBigIntegerConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT(BigInteger.valueOf(-129)));
    }

    @Test
    void testBigIntegerConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT(BigInteger.valueOf(128)));
    }

    @Test
    void testBigDecimalConstructorValid() {
        PlcSINT value = new PlcSINT(BigDecimal.valueOf(100));
        assertEquals(100, value.getByte());
    }

    @Test
    void testBigDecimalConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT(BigDecimal.valueOf(-150)));
    }

    @Test
    void testBigDecimalConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT(BigDecimal.valueOf(200)));
    }

    @Test
    void testStringConstructorValid() {
        PlcSINT value = new PlcSINT("-50");
        assertEquals(-50, value.getByte());
    }

    @Test
    void testStringConstructorMinValue() {
        PlcSINT value = new PlcSINT("-128");
        assertEquals(-128, value.getByte());
    }

    @Test
    void testStringConstructorMaxValue() {
        PlcSINT value = new PlcSINT("127");
        assertEquals(127, value.getByte());
    }

    @Test
    void testStringConstructorBelowRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT("-129"));
    }

    @Test
    void testStringConstructorAboveRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT("128"));
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcSINT("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcSINT value = new PlcSINT("  -75  ");
        assertEquals(-75, value.getByte());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcSINT() {
        PlcSINT original = new PlcSINT((byte) -50);
        PlcSINT copy = PlcSINT.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithInteger() {
        PlcSINT value = PlcSINT.of(100);
        assertEquals(100, value.getByte());
    }

    @Test
    void testOfMethodWithString() {
        PlcSINT value = PlcSINT.of("-100");
        assertEquals(-100, value.getByte());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcSINT zero = new PlcSINT((byte) 0);
        assertFalse(zero.getBoolean());

        PlcSINT positive = new PlcSINT((byte) 1);
        assertTrue(positive.getBoolean());

        PlcSINT negative = new PlcSINT((byte) -1);
        assertTrue(negative.getBoolean());

        PlcSINT max = new PlcSINT((byte) 127);
        assertTrue(max.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcSINT value = new PlcSINT((byte) -128);
        assertEquals(-128, value.getByte());
    }

    @Test
    void testGetShort() {
        PlcSINT value = new PlcSINT((byte) -50);
        assertEquals(-50, value.getShort());
    }

    @Test
    void testGetInteger() {
        PlcSINT value = new PlcSINT((byte) 100);
        assertEquals(100, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcSINT value = new PlcSINT((byte) -75);
        assertEquals(-75L, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcSINT value = new PlcSINT((byte) 127);
        assertEquals(BigInteger.valueOf(127), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcSINT value = new PlcSINT((byte) -50);
        assertEquals(-50.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcSINT value = new PlcSINT((byte) 100);
        assertEquals(100.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcSINT value = new PlcSINT((byte) -75);
        assertEquals(BigDecimal.valueOf(-75.0f), value.getBigDecimal());
    }

    @Test
    void testGetString() {
        PlcSINT value = new PlcSINT((byte) -100);
        assertEquals("-100", value.getString());
    }

    @Test
    void testToString() {
        PlcSINT value = new PlcSINT((byte) 127);
        assertEquals("127", value.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcSINT value = new PlcSINT((byte) 1);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcSINT value = new PlcSINT((byte) -100);
        assertTrue(value.isByte());
    }

    @Test
    void testIsShort() {
        PlcSINT value = new PlcSINT((byte) 50);
        assertTrue(value.isShort());
    }

    @Test
    void testIsInteger() {
        PlcSINT value = new PlcSINT((byte) -128);
        assertTrue(value.isInteger());
    }

    @Test
    void testIsLong() {
        PlcSINT value = new PlcSINT((byte) 100);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcSINT value = new PlcSINT((byte) -75);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcSINT value = new PlcSINT((byte) 50);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcSINT value = new PlcSINT((byte) -50);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcSINT value = new PlcSINT((byte) 127);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcSINT value = new PlcSINT((byte) -100);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesMinValue() {
        PlcSINT value = new PlcSINT((byte) -128);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals((byte) 0x80, bytes[0]);
    }

    @Test
    void testGetBytesMaxValue() {
        PlcSINT value = new PlcSINT((byte) 127);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals((byte) 0x7F, bytes[0]);
    }

    @Test
    void testGetBytesZero() {
        PlcSINT value = new PlcSINT((byte) 0);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals(0, bytes[0]);
    }

    @Test
    void testGetBytesNegativeValue() {
        PlcSINT value = new PlcSINT((byte) -1);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals((byte) 0xFF, bytes[0]);
    }

    @Test
    void testGetRaw() {
        PlcSINT value = new PlcSINT((byte) -50);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    // ========== Edge Cases ==========

    @Test
    void testNegativeOneValue() {
        PlcSINT value = new PlcSINT((byte) -1);
        assertEquals(-1, value.getByte());
    }

    @Test
    void testPositiveOneValue() {
        PlcSINT value = new PlcSINT((byte) 1);
        assertEquals(1, value.getByte());
    }

    @Test
    void testValueJustBelowMax() {
        PlcSINT value = new PlcSINT((byte) 126);
        assertEquals(126, value.getByte());
    }

    @Test
    void testValueJustAboveMin() {
        PlcSINT value = new PlcSINT((byte) -127);
        assertEquals(-127, value.getByte());
    }

    @Test
    void testNegativeValues() {
        for (int i = -128; i <= -1; i++) {
            PlcSINT value = new PlcSINT((byte) i);
            assertEquals(i, value.getByte());
        }
    }

    @Test
    void testPositiveValues() {
        for (int i = 0; i <= 127; i++) {
            PlcSINT value = new PlcSINT((byte) i);
            assertEquals(i, value.getByte());
        }
    }
}
