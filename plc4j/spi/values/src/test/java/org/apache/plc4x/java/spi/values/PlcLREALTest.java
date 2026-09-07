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

import java.math.BigDecimal;
import java.math.BigInteger;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcLREAL - 64-bit floating point (IEEE 754 double precision)
 */
class PlcLREALTest {

    private static final double DELTA = 0.0000001;

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcLREAL value = new PlcLREAL(-Double.MAX_VALUE);
        assertEquals(-Double.MAX_VALUE, value.getDouble(), DELTA);
        assertEquals(PlcValueType.LREAL, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcLREAL value = new PlcLREAL(Double.MAX_VALUE);
        assertEquals(Double.MAX_VALUE, value.getDouble(), DELTA);
    }

    @Test
    void testZero() {
        PlcLREAL value = new PlcLREAL(0.0);
        assertEquals(0.0, value.getDouble(), DELTA);
    }

    @Test
    void testPositiveInfinity() {
        PlcLREAL value = new PlcLREAL(Double.POSITIVE_INFINITY);
        assertEquals(Double.POSITIVE_INFINITY, value.getDouble());
        assertTrue(Double.isInfinite(value.getDouble()));
    }

    @Test
    void testNegativeInfinity() {
        PlcLREAL value = new PlcLREAL(Double.NEGATIVE_INFINITY);
        assertEquals(Double.NEGATIVE_INFINITY, value.getDouble());
        assertTrue(Double.isInfinite(value.getDouble()));
    }

    @Test
    void testNaN() {
        PlcLREAL value = new PlcLREAL(Double.NaN);
        assertTrue(Double.isNaN(value.getDouble()));
    }

    @Test
    void testMinPositiveValue() {
        PlcLREAL value = new PlcLREAL(Double.MIN_VALUE);
        assertEquals(Double.MIN_VALUE, value.getDouble(), DELTA);
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcLREAL trueValue = new PlcLREAL(true);
        assertEquals(1.0, trueValue.getDouble(), DELTA);

        PlcLREAL falseValue = new PlcLREAL(false);
        assertEquals(0.0, falseValue.getDouble(), DELTA);
    }

    @Test
    void testByteConstructor() {
        PlcLREAL value = new PlcLREAL((byte) 100);
        assertEquals(100.0, value.getDouble(), DELTA);
    }

    @Test
    void testShortConstructor() {
        PlcLREAL value = new PlcLREAL((short) -30000);
        assertEquals(-30000.0, value.getDouble(), DELTA);
    }

    @Test
    void testIntegerConstructor() {
        PlcLREAL value = new PlcLREAL(1000000000);
        assertEquals(1000000000.0, value.getDouble(), DELTA);
    }

    @Test
    void testLongConstructor() {
        PlcLREAL value = new PlcLREAL(9000000000000L);
        assertEquals(9000000000000.0, value.getDouble(), DELTA);
    }

    @Test
    void testFloatConstructor() {
        PlcLREAL value = new PlcLREAL(3.14159f);
        assertEquals(3.14159, value.getDouble(), 0.00001);
    }

    @Test
    void testDoubleConstructor() {
        PlcLREAL value = new PlcLREAL(3.141592653589793);
        assertEquals(3.141592653589793, value.getDouble(), DELTA);
    }

    @Test
    void testBigIntegerConstructor() {
        PlcLREAL value = new PlcLREAL(BigInteger.valueOf(123456789012345L));
        assertEquals(123456789012345.0, value.getDouble(), DELTA);
    }

    @Test
    void testBigDecimalConstructor() {
        PlcLREAL value = new PlcLREAL(new BigDecimal("123.456789012345"));
        assertEquals(123.456789012345, value.getDouble(), DELTA);
    }

    @Test
    void testStringConstructorValid() {
        PlcLREAL value = new PlcLREAL("3.141592653589793");
        assertEquals(3.141592653589793, value.getDouble(), DELTA);
    }

    @Test
    void testStringConstructorNegative() {
        PlcLREAL value = new PlcLREAL("-123.456789");
        assertEquals(-123.456789, value.getDouble(), DELTA);
    }

    @Test
    void testStringConstructorScientific() {
        PlcLREAL value = new PlcLREAL("1.23e15");
        assertEquals(1.23e15, value.getDouble(), DELTA);
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcLREAL value = new PlcLREAL("  -123.456789  ");
        assertEquals(-123.456789, value.getDouble(), DELTA);
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcLREAL() {
        PlcLREAL original = new PlcLREAL(3.141592653589793);
        PlcLREAL copy = PlcLREAL.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithDouble() {
        PlcLREAL value = PlcLREAL.of(2.718281828459045);
        assertEquals(2.718281828459045, value.getDouble(), DELTA);
    }

    @Test
    void testOfMethodWithString() {
        PlcLREAL value = PlcLREAL.of("1.414213562373095");
        assertEquals(1.414213562373095, value.getDouble(), DELTA);
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcLREAL zero = new PlcLREAL(0.0);
        assertFalse(zero.getBoolean());

        PlcLREAL positive = new PlcLREAL(0.1);
        assertTrue(positive.getBoolean());

        PlcLREAL negative = new PlcLREAL(-0.1);
        assertTrue(negative.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcLREAL value = new PlcLREAL(100.7);
        assertEquals(100, value.getByte());
    }

    @Test
    void testGetShort() {
        PlcLREAL value = new PlcLREAL(30000.5);
        assertEquals(30000, value.getShort());
    }

    @Test
    void testGetInteger() {
        PlcLREAL value = new PlcLREAL(1000000000.9);
        assertEquals(1000000000, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcLREAL value = new PlcLREAL(9000000000000.0);
        assertEquals(9000000000000L, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcLREAL value = new PlcLREAL(123456789012345.0);
        assertEquals(BigInteger.valueOf(123456789012345L), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcLREAL value = new PlcLREAL(3.14159);
        assertEquals(3.14159f, value.getFloat(), 0.00001f);
    }

    @Test
    void testGetDouble() {
        PlcLREAL value = new PlcLREAL(3.141592653589793);
        assertEquals(3.141592653589793, value.getDouble(), DELTA);
    }

    @Test
    void testGetBigDecimal() {
        PlcLREAL value = new PlcLREAL(123.456789012345);
        BigDecimal result = value.getBigDecimal();
        assertTrue(Math.abs(result.doubleValue() - 123.456789012345) < 0.0001);
    }

    @Test
    void testGetString() {
        PlcLREAL value = new PlcLREAL(123.456789);
        assertTrue(value.getString().startsWith("123.456"));
    }

    @Test
    void testToString() {
        PlcLREAL value = new PlcLREAL(3.14159);
        assertTrue(value.toString().startsWith("3.14"));
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcLREAL value = new PlcLREAL(1.0);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcLREAL smallValue = new PlcLREAL(100.0);
        assertTrue(smallValue.isByte());

        PlcLREAL largeValue = new PlcLREAL(300.0);
        assertFalse(largeValue.isByte());
    }

    @Test
    void testIsShort() {
        PlcLREAL smallValue = new PlcLREAL(30000.0);
        assertTrue(smallValue.isShort());

        PlcLREAL largeValue = new PlcLREAL(40000.0);
        assertFalse(largeValue.isShort());
    }

    @Test
    void testIsInteger() {
        PlcLREAL smallValue = new PlcLREAL(2000000000.0);
        assertTrue(smallValue.isInteger());

        PlcLREAL largeValue = new PlcLREAL(3000000000.0);
        assertFalse(largeValue.isInteger());
    }

    @Test
    void testIsLong() {
        PlcLREAL smallValue = new PlcLREAL(9000000000000.0);
        assertTrue(smallValue.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcLREAL value = new PlcLREAL(123456789012345.0);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcLREAL value = new PlcLREAL(3.14159);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcLREAL value = new PlcLREAL(123.456789);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcLREAL value = new PlcLREAL(123.456789);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcLREAL value = new PlcLREAL(123.456789);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesZero() {
        PlcLREAL value = new PlcLREAL(0.0);
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        for (int i = 0; i < 8; i++) {
            assertEquals(0, bytes[i]);
        }
    }

    @Test
    void testGetBytesPositive() {
        PlcLREAL value = new PlcLREAL(1.0);
        byte[] bytes = value.getBytes();
        assertEquals(8, bytes.length);
        // IEEE 754: 1.0 = 0x3FF0000000000000
        assertEquals((byte) 0x3F, bytes[0]);
        assertEquals((byte) 0xF0, bytes[1]);
        for (int i = 2; i < 8; i++) {
            assertEquals((byte) 0x00, bytes[i]);
        }
    }

    @Test
    void testGetRaw() {
        PlcLREAL value = new PlcLREAL(3.141592653589793);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    // ========== Edge Cases ==========

    @Test
    void testVerySmallPositiveValue() {
        PlcLREAL value = new PlcLREAL(0.0000000001);
        assertEquals(0.0000000001, value.getDouble(), 0.00000000001);
    }

    @Test
    void testVerySmallNegativeValue() {
        PlcLREAL value = new PlcLREAL(-0.0000000001);
        assertEquals(-0.0000000001, value.getDouble(), 0.00000000001);
    }

    @Test
    void testHighPrecision() {
        // Double has ~15-17 decimal digits of precision
        PlcLREAL value = new PlcLREAL(1.234567890123456);
        assertEquals(1.234567890123456, value.getDouble(), DELTA);
    }

    @Test
    void testNegativeZero() {
        PlcLREAL value = new PlcLREAL(-0.0);
        assertEquals(-0.0, value.getDouble());
    }

    @Test
    void testLargeInteger() {
        PlcLREAL value = new PlcLREAL(9007199254740992.0); // 2^53 - exact representation
        assertEquals(9007199254740992.0, value.getDouble(), DELTA);
    }

    @Test
    void testVeryLargeValue() {
        PlcLREAL value = new PlcLREAL(1.0e308);
        assertEquals(1.0e308, value.getDouble(), DELTA);
    }

    @Test
    void testVerySmallValue() {
        PlcLREAL value = new PlcLREAL(1.0e-308);
        assertEquals(1.0e-308, value.getDouble(), DELTA);
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfMethodWithByte() {
        PlcLREAL value = PlcLREAL.of((byte) 100);
        assertEquals(100.0, value.getDouble(), DELTA);
    }

    @Test
    void testOfMethodWithShort() {
        PlcLREAL value = PlcLREAL.of((short) 30000);
        assertEquals(30000.0, value.getDouble(), DELTA);
    }

    @Test
    void testOfMethodWithInteger() {
        PlcLREAL value = PlcLREAL.of(1000000000);
        assertEquals(1000000000.0, value.getDouble(), DELTA);
    }

    @Test
    void testOfMethodWithLong() {
        PlcLREAL value = PlcLREAL.of(9000000000000L);
        assertEquals(9000000000000.0, value.getDouble(), DELTA);
    }

    @Test
    void testOfMethodWithFloat() {
        PlcLREAL value = PlcLREAL.of(3.14159f);
        assertEquals(3.14159, value.getDouble(), 0.00001);
    }

    @Test
    void testOfMethodWithBigInteger() {
        PlcLREAL value = PlcLREAL.of(BigInteger.valueOf(123456789012345L));
        assertEquals(123456789012345.0, value.getDouble(), DELTA);
    }

    @Test
    void testOfMethodWithBigDecimal() {
        PlcLREAL value = PlcLREAL.of(new BigDecimal("123.456789012345"));
        assertEquals(123.456789012345, value.getDouble(), DELTA);
    }

    @Test
    void testOfMethodWithBoolean() {
        PlcLREAL trueValue = PlcLREAL.of(true);
        assertEquals(1.0, trueValue.getDouble(), DELTA);

        PlcLREAL falseValue = PlcLREAL.of(false);
        assertEquals(0.0, falseValue.getDouble(), DELTA);
    }

    // ========== Additional Edge Case Tests ==========

    @Test
    void testNullableStatus() {
        PlcLREAL value = new PlcLREAL(123.456789012345);
        assertFalse(value.isNullable());
    }

    @Test
    void testGetLength() {
        PlcLREAL value = new PlcLREAL(123.456789012345);
        assertEquals(1, value.getLength());
    }

    @Test
    void testGetIndex() {
        PlcLREAL value = new PlcLREAL(123.456789012345);
        assertEquals(value, value.getIndex(0));
    }

    @Test
    void testIsSimpleType() {
        PlcLREAL value = new PlcLREAL(123.456789012345);
        assertTrue(value.isSimple());
    }

    @Test
    void testGetObject() {
        PlcLREAL value = new PlcLREAL(123.456789012345);
        Object obj = value.getObject();
        assertNotNull(obj);
        assertEquals(123.456789012345, ((Double) obj).doubleValue(), DELTA);
    }

    @Test
    void testSerializeToString() {
        PlcLREAL value = new PlcLREAL(123.456789012345);
        String serialized = value.toString();
        assertNotNull(serialized);
        assertTrue(serialized.contains("123"));
    }
}
