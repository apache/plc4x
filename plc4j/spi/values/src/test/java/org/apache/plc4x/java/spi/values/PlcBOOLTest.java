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
 * Test class for PlcBOOL - Boolean value (true/false)
 */
class PlcBOOLTest {

    // ========== Basic Value Tests ==========

    @Test
    void testTrueValue() {
        PlcBOOL value = new PlcBOOL(true);
        assertTrue(value.getBoolean());
        assertEquals(PlcValueType.BOOL, value.getPlcValueType());
    }

    @Test
    void testFalseValue() {
        PlcBOOL value = new PlcBOOL(false);
        assertFalse(value.getBoolean());
    }

    @Test
    void testNullValue() {
        PlcBOOL value = new PlcBOOL((Boolean) null);
        assertNull(value.value); // Direct field access since getBoolean may handle null differently
        assertTrue(value.isNullable());
    }

    // ========== Constructor Tests with Different Types ==========

    @Test
    void testBooleanConstructor() {
        PlcBOOL trueValue = new PlcBOOL(true);
        assertTrue(trueValue.getBoolean());

        PlcBOOL falseValue = new PlcBOOL(false);
        assertFalse(falseValue.getBoolean());
    }

    @Test
    void testByteConstructorZero() {
        PlcBOOL value = new PlcBOOL((byte) 0);
        assertFalse(value.getBoolean());
    }

    @Test
    void testByteConstructorNonZero() {
        PlcBOOL value = new PlcBOOL((byte) 1);
        assertTrue(value.getBoolean());

        PlcBOOL value2 = new PlcBOOL((byte) -1);
        assertTrue(value2.getBoolean());

        PlcBOOL value3 = new PlcBOOL((byte) 100);
        assertTrue(value3.getBoolean());
    }

    @Test
    void testByteConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcBOOL((Byte) null));
    }

    @Test
    void testShortConstructorZero() {
        PlcBOOL value = new PlcBOOL((short) 0);
        assertFalse(value.getBoolean());
    }

    @Test
    void testShortConstructorNonZero() {
        PlcBOOL value = new PlcBOOL((short) 1);
        assertTrue(value.getBoolean());

        PlcBOOL value2 = new PlcBOOL((short) -1000);
        assertTrue(value2.getBoolean());
    }

    @Test
    void testShortConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcBOOL((Short) null));
    }

    @Test
    void testIntegerConstructorZero() {
        PlcBOOL value = new PlcBOOL(0);
        assertFalse(value.getBoolean());
    }

    @Test
    void testIntegerConstructorNonZero() {
        PlcBOOL value = new PlcBOOL(1);
        assertTrue(value.getBoolean());

        PlcBOOL value2 = new PlcBOOL(-1000000);
        assertTrue(value2.getBoolean());
    }

    @Test
    void testIntegerConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcBOOL((Integer) null));
    }

    @Test
    void testLongConstructorZero() {
        PlcBOOL value = new PlcBOOL(0L);
        assertFalse(value.getBoolean());
    }

    @Test
    void testLongConstructorNonZero() {
        PlcBOOL value = new PlcBOOL(1L);
        assertTrue(value.getBoolean());

        PlcBOOL value2 = new PlcBOOL(5000000000L);
        assertTrue(value2.getBoolean());
    }

    @Test
    void testLongConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcBOOL((Long) null));
    }

    @Test
    void testFloatConstructorZero() {
        PlcBOOL value = new PlcBOOL(0.0f);
        assertFalse(value.getBoolean());
    }

    @Test
    void testFloatConstructorNonZero() {
        PlcBOOL value = new PlcBOOL(0.1f);
        assertTrue(value.getBoolean());

        PlcBOOL value2 = new PlcBOOL(-0.1f);
        assertTrue(value2.getBoolean());
    }

    @Test
    void testFloatConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcBOOL((Float) null));
    }

    @Test
    void testDoubleConstructorZero() {
        PlcBOOL value = new PlcBOOL(0.0);
        assertFalse(value.getBoolean());
    }

    @Test
    void testDoubleConstructorNonZero() {
        PlcBOOL value = new PlcBOOL(0.1);
        assertTrue(value.getBoolean());

        PlcBOOL value2 = new PlcBOOL(-0.1);
        assertTrue(value2.getBoolean());
    }

    @Test
    void testDoubleConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcBOOL((Double) null));
    }

    @Test
    void testBigIntegerConstructorZero() {
        PlcBOOL value = new PlcBOOL(BigInteger.ZERO);
        assertFalse(value.getBoolean());
    }

    @Test
    void testBigIntegerConstructorNonZero() {
        PlcBOOL value = new PlcBOOL(BigInteger.ONE);
        assertTrue(value.getBoolean());

        PlcBOOL value2 = new PlcBOOL(BigInteger.valueOf(-1000));
        assertTrue(value2.getBoolean());
    }

    @Test
    void testBigIntegerConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcBOOL((BigInteger) null));
    }

    @Test
    void testBigDecimalConstructorZero() {
        PlcBOOL value = new PlcBOOL(BigDecimal.ZERO);
        assertFalse(value.getBoolean());
    }

    @Test
    void testBigDecimalConstructorNonZero() {
        PlcBOOL value = new PlcBOOL(BigDecimal.ONE);
        assertTrue(value.getBoolean());

        PlcBOOL value2 = new PlcBOOL(new BigDecimal("0.0001"));
        assertTrue(value2.getBoolean());
    }

    @Test
    void testBigDecimalConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcBOOL((BigDecimal) null));
    }

    @Test
    void testStringConstructorTrue() {
        PlcBOOL value1 = new PlcBOOL("true");
        assertTrue(value1.getBoolean());

        PlcBOOL value2 = new PlcBOOL("TRUE");
        assertTrue(value2.getBoolean());

        PlcBOOL value3 = new PlcBOOL("True");
        assertTrue(value3.getBoolean());
    }

    @Test
    void testStringConstructorFalse() {
        PlcBOOL value1 = new PlcBOOL("false");
        assertFalse(value1.getBoolean());

        PlcBOOL value2 = new PlcBOOL("FALSE");
        assertFalse(value2.getBoolean());

        PlcBOOL value3 = new PlcBOOL("False");
        assertFalse(value3.getBoolean());
    }

    @Test
    void testStringConstructorNumericTrue() {
        PlcBOOL value1 = new PlcBOOL("1");
        assertTrue(value1.getBoolean());

        PlcBOOL value2 = new PlcBOOL("100");
        assertTrue(value2.getBoolean());

        // -1 is NOT > 0, so it should be false based on implementation
        PlcBOOL value3 = new PlcBOOL("-1");
        assertFalse(value3.getBoolean());
    }

    @Test
    void testStringConstructorNumericFalse() {
        PlcBOOL value = new PlcBOOL("0");
        assertFalse(value.getBoolean());
    }

    @Test
    void testStringConstructorWithWhitespace() {
        PlcBOOL value = new PlcBOOL("  true  ");
        assertTrue(value.getBoolean());
    }

    // ========== Factory Method Tests ==========

    @Test
    void testOfMethodWithPlcBOOL() {
        PlcBOOL original = new PlcBOOL(true);
        PlcBOOL copy = PlcBOOL.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfMethodWithBoolean() {
        PlcBOOL value = PlcBOOL.of(true);
        assertTrue(value.getBoolean());
    }

    @Test
    void testOfMethodWithInteger() {
        PlcBOOL value = PlcBOOL.of(1);
        assertTrue(value.getBoolean());
    }

    @Test
    void testOfMethodWithString() {
        PlcBOOL value = PlcBOOL.of("true");
        assertTrue(value.getBoolean());
    }

    // ========== Getter Tests ==========

    @Test
    void testGetBoolean() {
        PlcBOOL trueValue = new PlcBOOL(true);
        assertTrue(trueValue.getBoolean());

        PlcBOOL falseValue = new PlcBOOL(false);
        assertFalse(falseValue.getBoolean());
    }

    @Test
    void testGetByteTrue() {
        PlcBOOL value = new PlcBOOL(true);
        assertEquals(1, value.getByte());
    }

    @Test
    void testGetByteFalse() {
        PlcBOOL value = new PlcBOOL(false);
        assertEquals(0, value.getByte());
    }

    @Test
    void testGetShortTrue() {
        PlcBOOL value = new PlcBOOL(true);
        assertEquals(1, value.getShort());
    }

    @Test
    void testGetShortFalse() {
        PlcBOOL value = new PlcBOOL(false);
        assertEquals(0, value.getShort());
    }

    @Test
    void testGetIntegerTrue() {
        PlcBOOL value = new PlcBOOL(true);
        assertEquals(1, value.getInteger());
    }

    @Test
    void testGetIntegerFalse() {
        PlcBOOL value = new PlcBOOL(false);
        assertEquals(0, value.getInteger());
    }

    @Test
    void testGetLongTrue() {
        PlcBOOL value = new PlcBOOL(true);
        assertEquals(1L, value.getLong());
    }

    @Test
    void testGetLongFalse() {
        PlcBOOL value = new PlcBOOL(false);
        assertEquals(0L, value.getLong());
    }

    @Test
    void testGetBigIntegerTrue() {
        PlcBOOL value = new PlcBOOL(true);
        assertEquals(BigInteger.ONE, value.getBigInteger());
    }

    @Test
    void testGetBigIntegerFalse() {
        PlcBOOL value = new PlcBOOL(false);
        assertEquals(BigInteger.ZERO, value.getBigInteger());
    }

    @Test
    void testGetStringTrue() {
        PlcBOOL value = new PlcBOOL(true);
        assertEquals("true", value.getString());
    }

    @Test
    void testGetStringFalse() {
        PlcBOOL value = new PlcBOOL(false);
        assertEquals("false", value.getString());
    }

    @Test
    void testToString() {
        PlcBOOL trueValue = new PlcBOOL(true);
        assertEquals("true", trueValue.toString());

        PlcBOOL falseValue = new PlcBOOL(false);
        assertEquals("false", falseValue.toString());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsBoolean() {
        PlcBOOL value = new PlcBOOL(true);
        assertTrue(value.isBoolean());
    }

    @Test
    void testIsByte() {
        PlcBOOL value = new PlcBOOL(true);
        assertTrue(value.isByte());
    }

    @Test
    void testIsShort() {
        PlcBOOL value = new PlcBOOL(false);
        assertTrue(value.isShort());
    }

    @Test
    void testIsInteger() {
        PlcBOOL value = new PlcBOOL(true);
        assertTrue(value.isInteger());
    }

    @Test
    void testIsLong() {
        PlcBOOL value = new PlcBOOL(false);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcBOOL value = new PlcBOOL(true);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsString() {
        PlcBOOL value = new PlcBOOL(false);
        assertTrue(value.isString());
    }

    // ========== Byte Serialization Tests ==========

    @Test
    void testGetBytesTrue() {
        PlcBOOL value = new PlcBOOL(true);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals(1, bytes[0]);
    }

    @Test
    void testGetBytesFalse() {
        PlcBOOL value = new PlcBOOL(false);
        byte[] bytes = value.getBytes();
        assertEquals(1, bytes.length);
        assertEquals(0, bytes[0]);
    }

    @Test
    void testGetRaw() {
        PlcBOOL value = new PlcBOOL(true);
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }
}
