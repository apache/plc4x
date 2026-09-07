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

public class PlcCHARTest {

    @Test
    void testConstructorWithBoolean() {
        PlcCHAR plcCHAR = new PlcCHAR(true);
        assertEquals('T', plcCHAR.getString().charAt(0));
        assertFalse(plcCHAR.isNullable());

        plcCHAR = new PlcCHAR(false);
        assertEquals('F', plcCHAR.getString().charAt(0));
        assertFalse(plcCHAR.isNullable());
    }

    @Test
    void testConstructorWithCharacter() {
        PlcCHAR plcCHAR = new PlcCHAR('A');
        assertEquals("A", plcCHAR.getString());
        assertFalse(plcCHAR.isNullable());
    }

    @Test
    void testConstructorWithValidByte() {
        PlcCHAR plcCHAR = new PlcCHAR((byte) 65);
        assertEquals("A", plcCHAR.getString());
        assertFalse(plcCHAR.isNullable());
    }

    @Test
    void testConstructorWithValidShort() {
        PlcCHAR plcCHAR = new PlcCHAR((short) 66);
        assertEquals("B", plcCHAR.getString());
        assertFalse(plcCHAR.isNullable());
    }

    @Test
    void testConstructorWithValidInteger() {
        PlcCHAR plcCHAR = new PlcCHAR(67);
        assertEquals("C", plcCHAR.getString());
        assertFalse(plcCHAR.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeInteger() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcCHAR(-1));
        assertThrows(PlcInvalidTagException.class, () -> new PlcCHAR(256));
    }

    @Test
    void testConstructorWithValidString() {
        PlcCHAR plcCHAR = new PlcCHAR("X");
        assertEquals("X", plcCHAR.getString());
        assertFalse(plcCHAR.isNullable());
    }

    @Test
    void testConstructorWithStringWithWhitespace() {
        PlcCHAR plcCHAR = new PlcCHAR("  Y  ");
        assertEquals("Y", plcCHAR.getString());
        assertFalse(plcCHAR.isNullable());
    }

    @Test
    void testConstructorWithSpace() {
        PlcCHAR plcCHAR = new PlcCHAR("   ");
        assertEquals(" ", plcCHAR.getString());
        assertFalse(plcCHAR.isNullable());
    }

    @Test
    void testGetPlcValueType() {
        PlcCHAR plcCHAR = new PlcCHAR('Z');
        assertEquals(PlcValueType.CHAR, plcCHAR.getPlcValueType());
    }

    @Test
    void testIsBoolean() {
        PlcCHAR plcCHAR = new PlcCHAR('A');
        assertTrue(plcCHAR.isBoolean());
    }

    @Test
    void testGetBoolean() {
        PlcCHAR plcCHAR = new PlcCHAR('A');
        assertTrue(plcCHAR.getBoolean());

        PlcCHAR plcCHAR2 = new PlcCHAR((char) 0);
        assertFalse(plcCHAR2.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcCHAR plcCHAR = new PlcCHAR('A');
        assertEquals(65, plcCHAR.getByte());
    }

    @Test
    void testGetShort() {
        PlcCHAR plcCHAR = new PlcCHAR('B');
        assertEquals(66, plcCHAR.getShort());
    }

    @Test
    void testGetInteger() {
        PlcCHAR plcCHAR = new PlcCHAR('C');
        assertEquals(67, plcCHAR.getInteger());
    }

    @Test
    void testGetLong() {
        PlcCHAR plcCHAR = new PlcCHAR('D');
        assertEquals(68L, plcCHAR.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcCHAR plcCHAR = new PlcCHAR('E');
        assertEquals(BigInteger.valueOf(69), plcCHAR.getBigInteger());
    }

    @Test
    void testGetString() {
        PlcCHAR plcCHAR = new PlcCHAR('Z');
        assertEquals("Z", plcCHAR.getString());
    }

    @Test
    void testToString() {
        PlcCHAR plcCHAR = new PlcCHAR('M');
        assertEquals("M", plcCHAR.toString());
    }

    @Test
    void testGetObject() {
        PlcCHAR plcCHAR = new PlcCHAR('Q');
        assertEquals("Q", plcCHAR.getObject());
    }

    @Test
    void testGetBytes() {
        PlcCHAR plcCHAR = new PlcCHAR('A');
        byte[] bytes = plcCHAR.getBytes();
        assertEquals(1, bytes.length);
        assertEquals(65, bytes[0]);
    }

    @Test
    void testOfWithPlcCHAR() {
        PlcCHAR original = new PlcCHAR('A');
        PlcCHAR copy = PlcCHAR.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithCharacter() {
        PlcCHAR plcCHAR = PlcCHAR.of('X');
        assertEquals("X", plcCHAR.getString());
    }

    @Test
    void testOfWithString() {
        PlcCHAR plcCHAR = PlcCHAR.of("Y");
        assertEquals("Y", plcCHAR.getString());
    }

    @Test
    void testOfWithBoolean() {
        PlcCHAR plcCHAR = PlcCHAR.of(true);
        assertEquals('T', plcCHAR.getString().charAt(0));
    }

    @Test
    void testOfWithByte() {
        PlcCHAR plcCHAR = PlcCHAR.of((byte) 65);
        assertEquals("A", plcCHAR.getString());
    }

    @Test
    void testOfWithInteger() {
        PlcCHAR plcCHAR = PlcCHAR.of(66);
        assertEquals("B", plcCHAR.getString());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testLongConstructor() {
        PlcCHAR value = new PlcCHAR(100L);
        assertEquals(100, value.getInteger());
    }

    @Test
    void testLongConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((Long) null));
    }

    @Test
    void testFloatConstructor() {
        PlcCHAR value = new PlcCHAR(123.0f);
        assertEquals(123, value.getInteger());
    }

    @Test
    void testFloatConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((Float) null));
    }

    @Test
    void testDoubleConstructor() {
        PlcCHAR value = new PlcCHAR(200.0);
        assertEquals(200, value.getInteger());
    }

    @Test
    void testDoubleConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((Double) null));
    }

    @Test
    void testBigIntegerConstructor() {
        PlcCHAR value = new PlcCHAR(BigInteger.valueOf(150));
        assertEquals(150, value.getInteger());
    }

    @Test
    void testBigIntegerConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((BigInteger) null));
    }

    @Test
    void testBigDecimalConstructor() {
        PlcCHAR value = new PlcCHAR(BigDecimal.valueOf(100));
        assertEquals(100, value.getInteger());
    }

    @Test
    void testBigDecimalConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((BigDecimal) null));
    }

    // ========== Additional Getter Tests ==========

    @Test
    void testGetFloat() {
        PlcCHAR value = new PlcCHAR(100);
        assertEquals(100.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcCHAR value = new PlcCHAR(150);
        assertEquals(150.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcCHAR value = new PlcCHAR(123);
        assertEquals(0, BigDecimal.valueOf(123).compareTo(value.getBigDecimal()));
    }

    // ========== Additional Is* Methods Tests ==========

    @Test
    void testIsLong() {
        PlcCHAR value = new PlcCHAR(200);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcCHAR value = new PlcCHAR(250);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcCHAR value = new PlcCHAR(100);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcCHAR value = new PlcCHAR(150);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcCHAR value = new PlcCHAR(123);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsByte() {
        PlcCHAR value = new PlcCHAR(100);
        assertTrue(value.isByte());
    }

    @Test
    void testIsShort() {
        PlcCHAR value = new PlcCHAR(200);
        assertTrue(value.isShort());
    }

    @Test
    void testIsInteger() {
        PlcCHAR value = new PlcCHAR(250);
        assertTrue(value.isInteger());
    }

    // ========== Edge Case Tests ==========

    @Test
    void testIntegerConstructorOutOfRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcCHAR(300));
    }

    @Test
    void testLongConstructorOutOfRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcCHAR(300L));
    }

    @Test
    void testByteConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((Byte) null));
    }

    @Test
    void testShortConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((Short) null));
    }

    @Test
    void testIntegerConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((Integer) null));
    }

    @Test
    void testMinValue() {
        PlcCHAR value = new PlcCHAR(0);
        assertEquals(0, value.getInteger());
    }

    @Test
    void testMaxValue() {
        PlcCHAR value = new PlcCHAR(255);
        assertEquals(255, value.getInteger());
    }

    @Test
    void testNegativeValue() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcCHAR(-1));
    }

    // ========== Additional Tests for Uncovered Code Paths ==========

    @Test
    void testCharacterConstructorOutOfRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcCHAR((char) 256));
    }

    @Test
    void testCharacterConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcCHAR((Character) null));
    }

    @Test
    void testBigDecimalWithScale() {
        // BigDecimal with scale > 0 should throw exception
        assertThrows(PlcInvalidTagException.class, () -> new PlcCHAR(BigDecimal.valueOf(100.5)));
    }

    @Test
    void testBigIntegerIsNullable() {
        PlcCHAR value = new PlcCHAR(BigInteger.valueOf(100));
        assertTrue(value.isNullable());
    }

    @Test
    void testBigDecimalIsNullable() {
        PlcCHAR value = new PlcCHAR(BigDecimal.valueOf(100));
        assertTrue(value.isNullable());
    }

    @Test
    void testStringConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcCHAR((String) null));
    }

    @Test
    void testIsString() {
        PlcCHAR value = new PlcCHAR(100);
        assertTrue(value.isString());
    }

    @Test
    void testGetRaw() {
        PlcCHAR value = new PlcCHAR('A');
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    @Test
    void testIsByteReturnsFalse() {
        // Value > Byte.MAX_VALUE (127) should return false for isByte()
        PlcCHAR value = new PlcCHAR(200);
        assertFalse(value.isByte());
    }

    @Test
    void testIsByteReturnsTrue() {
        // Value <= Byte.MAX_VALUE (127) should return true for isByte()
        PlcCHAR value = new PlcCHAR(100);
        assertTrue(value.isByte());
    }
}
