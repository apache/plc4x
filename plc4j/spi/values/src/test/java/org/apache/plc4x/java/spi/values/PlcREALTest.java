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
 * Test class for PlcREAL - 32-bit floating point (IEEE 754 single precision)
 */
class PlcREALTest {

    private static final float DELTA = 0.0001f;

    // ========== Boundary Value Tests ==========

    @Test
    void testMinValue() {
        PlcREAL value = new PlcREAL(-Float.MAX_VALUE);
        assertEquals(-Float.MAX_VALUE, value.getFloat(), DELTA);
        assertEquals(PlcValueType.REAL, value.getPlcValueType());
    }

    @Test
    void testMaxValue() {
        PlcREAL value = new PlcREAL(Float.MAX_VALUE);
        assertEquals(Float.MAX_VALUE, value.getFloat(), DELTA);
    }

    @Test
    void testZero() {
        PlcREAL value = new PlcREAL(0.0f);
        assertEquals(0.0f, value.getFloat(), DELTA);
    }

    @Test
    void testPositiveInfinity() {
        PlcREAL value = new PlcREAL(Float.POSITIVE_INFINITY);
        assertEquals(Float.POSITIVE_INFINITY, value.getFloat());
        assertTrue(Float.isInfinite(value.getFloat()));
    }

    @Test
    void testNegativeInfinity() {
        PlcREAL value = new PlcREAL(Float.NEGATIVE_INFINITY);
        assertEquals(Float.NEGATIVE_INFINITY, value.getFloat());
        assertTrue(Float.isInfinite(value.getFloat()));
    }

    @Test
    void testNaN() {
        PlcREAL value = new PlcREAL(Float.NaN);
        assertTrue(Float.isNaN(value.getFloat()));
    }

    @Test
    void testMinPositiveValue() {
        PlcREAL value = new PlcREAL(Float.MIN_VALUE);
        assertEquals(Float.MIN_VALUE, value.getFloat(), DELTA);
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcREAL trueValue = new PlcREAL(true);
        assertEquals(1.0f, trueValue.getFloat(), DELTA);

        PlcREAL falseValue = new PlcREAL(false);
        assertEquals(0.0f, falseValue.getFloat(), DELTA);
    }

    @Test
    void testByteConstructor() {
        PlcREAL value = new PlcREAL((byte) 100);
        assertEquals(100.0f, value.getFloat(), DELTA);
    }

    @Test
    void testShortConstructor() {
        PlcREAL value = new PlcREAL((short) -30000);
        assertEquals(-30000.0f, value.getFloat(), DELTA);
    }

    @Test
    void testIntegerConstructor() {
        PlcREAL value = new PlcREAL(1000000);
        assertEquals(1000000.0f, value.getFloat(), DELTA);
    }

    @Test
    void testLongConstructor() {
        PlcREAL value = new PlcREAL(5000000000L);
        assertEquals(5000000000.0f, value.getFloat(), DELTA);
    }

    @Test
    void testFloatConstructor() {
        PlcREAL value = new PlcREAL(3.14159f);
        assertEquals(3.14159f, value.getFloat(), DELTA);
    }

    @Test
    void testDoubleConstructorValid() {
        PlcREAL value = new PlcREAL(123.456);
        assertEquals(123.456f, value.getFloat(), DELTA);
    }

    @Test
    void testDoubleConstructorInfinity() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcREAL(Double.MAX_VALUE));
    }

    @Test
    void testBigIntegerConstructor() {
        PlcREAL value = new PlcREAL(BigInteger.valueOf(1234567));
        assertEquals(1234567.0f, value.getFloat(), DELTA);
    }

    @Test
    void testBigDecimalConstructor() {
        PlcREAL value = new PlcREAL(new BigDecimal("123.456"));
        assertEquals(123.456f, value.getFloat(), DELTA);
    }

    @Test
    void testStringConstructorValid() {
        PlcREAL value = new PlcREAL("3.14159");
        assertEquals(3.14159f, value.getFloat(), DELTA);
    }

    @Test
    void testStringConstructorNegative() {
        PlcREAL value = new PlcREAL("-123.456");
        assertEquals(-123.456f, value.getFloat(), DELTA);
    }

    @Test
    void testStringConstructorScientific() {
        PlcREAL value = new PlcREAL("1.23e5");
        assertEquals(123000.0f, value.getFloat(), DELTA);
    }

    @Test
    void testStringConstructorInvalidFormat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcREAL("not a number"));
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcREAL value = new PlcREAL("  -123.456  ");
        assertEquals(-123.456f, value.getFloat(), DELTA);
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcREAL() {
        PlcREAL original = new PlcREAL(3.14159f);
        PlcREAL copy = PlcREAL.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithFloat() {
        PlcREAL value = PlcREAL.of(2.71828f);
        assertEquals(2.71828f, value.getFloat(), DELTA);
    }

    @Test
    void testOfMethodWithString() {
        PlcREAL value = PlcREAL.of("1.41421");
        assertEquals(1.41421f, value.getFloat(), DELTA);
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcREAL zero = new PlcREAL(0.0f);
        assertFalse(zero.getBoolean());

        PlcREAL positive = new PlcREAL(0.1f);
        assertTrue(positive.getBoolean());

        PlcREAL negative = new PlcREAL(-0.1f);
        assertTrue(negative.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcREAL value = new PlcREAL(100.7f);
        assertEquals(100, value.getByte());
    }

    @Test
    void testGetShort() {
        PlcREAL value = new PlcREAL(30000.5f);
        assertEquals(30000, value.getShort());
    }

    @Test
    void testGetInteger() {
        PlcREAL value = new PlcREAL(1000000.9f);
        assertEquals(1000000, value.getInteger());
    }

    @Test
    void testGetLong() {
        PlcREAL value = new PlcREAL(5000000.0f);
        assertEquals(5000000L, value.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcREAL value = new PlcREAL(1234567.0f);
        assertEquals(BigInteger.valueOf(1234567), value.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcREAL value = new PlcREAL(3.14159f);
        assertEquals(3.14159f, value.getFloat(), DELTA);
    }

    @Test
    void testGetDouble() {
        PlcREAL value = new PlcREAL(123.456f);
        assertEquals(123.456, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcREAL value = new PlcREAL(123.456f);
        BigDecimal result = value.getBigDecimal();
        assertTrue(result.doubleValue() > 123.4 && result.doubleValue() < 123.5);
    }

    @Test
    void testGetString() {
        PlcREAL value = new PlcREAL(123.456f);
        assertTrue(value.getString().startsWith("123.456"));
    }

    @Test
    void testToString() {
        PlcREAL value = new PlcREAL(3.14f);
        assertTrue(value.toString().startsWith("3.14"));
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcREAL value = new PlcREAL(1.0f);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcREAL smallValue = new PlcREAL(100.0f);
        assertTrue(smallValue.isByte());

        PlcREAL largeValue = new PlcREAL(300.0f);
        assertFalse(largeValue.isByte());
    }

    @Test
    void testIsShort() {
        PlcREAL smallValue = new PlcREAL(30000.0f);
        assertTrue(smallValue.isShort());

        PlcREAL largeValue = new PlcREAL(40000.0f);
        assertFalse(largeValue.isShort());
    }

    @Test
    void testIsInteger() {
        PlcREAL smallValue = new PlcREAL(2000000000.0f);
        assertTrue(smallValue.isInteger());

        PlcREAL largeValue = new PlcREAL(3000000000.0f);
        assertFalse(largeValue.isInteger());
    }

    @Test
    void testIsLong() {
        PlcREAL value = new PlcREAL(5000000.0f);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcREAL value = new PlcREAL(1234567.0f);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcREAL value = new PlcREAL(3.14159f);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcREAL value = new PlcREAL(123.456f);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcREAL value = new PlcREAL(123.456f);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsString() {
        PlcREAL value = new PlcREAL(123.456f);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesZero() {
        PlcREAL value = new PlcREAL(0.0f);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        for (int i = 0; i < 4; i++) {
            assertEquals(0, bytes[i]);
        }
    }

    @Test
    void testGetBytesPositive() {
        PlcREAL value = new PlcREAL(1.0f);
        byte[] bytes = value.getBytes();
        assertEquals(4, bytes.length);
        // IEEE 754: 1.0f = 0x3F800000
        assertEquals((byte) 0x3F, bytes[0]);
        assertEquals((byte) 0x80, bytes[1]);
        assertEquals((byte) 0x00, bytes[2]);
        assertEquals((byte) 0x00, bytes[3]);
    }

    @Test
    void testGetRaw() {
        PlcREAL value = new PlcREAL(3.14159f);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    // ========== Edge Cases ==========

    @Test
    void testVerySmallPositiveValue() {
        PlcREAL value = new PlcREAL(0.000001f);
        assertEquals(0.000001f, value.getFloat(), 0.0000001f);
    }

    @Test
    void testVerySmallNegativeValue() {
        PlcREAL value = new PlcREAL(-0.000001f);
        assertEquals(-0.000001f, value.getFloat(), 0.0000001f);
    }

    @Test
    void testPrecisionLoss() {
        // Float has ~7 decimal digits of precision
        PlcREAL value = new PlcREAL(1.23456789f);
        float result = value.getFloat();
        assertTrue(Math.abs(result - 1.23456789f) < 0.00001f);
    }

    @Test
    void testNegativeZero() {
        PlcREAL value = new PlcREAL(-0.0f);
        assertEquals(-0.0f, value.getFloat());
    }

    @Test
    void testLargeInteger() {
        PlcREAL value = new PlcREAL(16777216.0f); // 2^24 - exact representation
        assertEquals(16777216.0f, value.getFloat(), DELTA);
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfMethodWithByte() {
        PlcREAL value = PlcREAL.of((byte) 100);
        assertEquals(100.0f, value.getFloat(), DELTA);
    }

    @Test
    void testOfMethodWithShort() {
        PlcREAL value = PlcREAL.of((short) 30000);
        assertEquals(30000.0f, value.getFloat(), DELTA);
    }

    @Test
    void testOfMethodWithInteger() {
        PlcREAL value = PlcREAL.of(1000000);
        assertEquals(1000000.0f, value.getFloat(), DELTA);
    }

    @Test
    void testOfMethodWithLong() {
        PlcREAL value = PlcREAL.of(5000000L);
        assertEquals(5000000.0f, value.getFloat(), DELTA);
    }

    @Test
    void testOfMethodWithDouble() {
        PlcREAL value = PlcREAL.of(123.456);
        assertEquals(123.456f, value.getFloat(), DELTA);
    }

    @Test
    void testOfMethodWithBigInteger() {
        PlcREAL value = PlcREAL.of(BigInteger.valueOf(1000000));
        assertEquals(1000000.0f, value.getFloat(), DELTA);
    }

    @Test
    void testOfMethodWithBigDecimal() {
        PlcREAL value = PlcREAL.of(new BigDecimal("123.456"));
        assertEquals(123.456f, value.getFloat(), DELTA);
    }

    @Test
    void testOfMethodWithBoolean() {
        PlcREAL trueValue = PlcREAL.of(true);
        assertEquals(1.0f, trueValue.getFloat(), DELTA);

        PlcREAL falseValue = PlcREAL.of(false);
        assertEquals(0.0f, falseValue.getFloat(), DELTA);
    }

    // ========== Additional Edge Case Tests ==========

    @Test
    void testNullableStatus() {
        PlcREAL value = new PlcREAL(123.456f);
        assertFalse(value.isNullable());
    }

    @Test
    void testGetLength() {
        PlcREAL value = new PlcREAL(123.456f);
        assertEquals(1, value.getLength());
    }

    @Test
    void testGetIndex() {
        PlcREAL value = new PlcREAL(123.456f);
        assertEquals(value, value.getIndex(0));
    }

    @Test
    void testIsSimpleType() {
        PlcREAL value = new PlcREAL(123.456f);
        assertTrue(value.isSimple());
    }

    @Test
    void testGetObject() {
        PlcREAL value = new PlcREAL(123.456f);
        Object obj = value.getObject();
        assertNotNull(obj);
        assertEquals(123.456f, ((Float) obj).floatValue(), DELTA);
    }

    @Test
    void testSerializeToString() {
        PlcREAL value = new PlcREAL(123.456f);
        String serialized = value.toString();
        assertNotNull(serialized);
        assertTrue(serialized.contains("123"));
    }
}
