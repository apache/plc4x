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

public class PlcWCHARTest {

    @Test
    void testConstructorWithBoolean() {
        PlcWCHAR plcWCHAR = new PlcWCHAR(true);
        assertEquals('T', plcWCHAR.getString().charAt(0));
        assertFalse(plcWCHAR.isNullable());

        plcWCHAR = new PlcWCHAR(false);
        assertEquals('F', plcWCHAR.getString().charAt(0));
        assertFalse(plcWCHAR.isNullable());
    }

    @Test
    void testConstructorWithCharacter() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('A');
        assertEquals("A", plcWCHAR.getString());
        assertFalse(plcWCHAR.isNullable());
    }

    @Test
    void testConstructorWithValidByte() {
        PlcWCHAR plcWCHAR = new PlcWCHAR((byte) 65);
        assertEquals("A", plcWCHAR.getString());
        assertFalse(plcWCHAR.isNullable());
    }

    @Test
    void testConstructorWithValidShort() {
        PlcWCHAR plcWCHAR = new PlcWCHAR((short) 12354);
        assertEquals(12354, plcWCHAR.getInteger());
        assertFalse(plcWCHAR.isNullable());
    }

    @Test
    void testConstructorWithValidInteger() {
        PlcWCHAR plcWCHAR = new PlcWCHAR(30000);
        assertEquals(30000, plcWCHAR.getInteger());
        assertFalse(plcWCHAR.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeInteger() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWCHAR(-1));
        assertThrows(PlcInvalidTagException.class, () -> new PlcWCHAR(65536));
    }

    @Test
    void testConstructorWithValidString() {
        PlcWCHAR plcWCHAR = new PlcWCHAR("X");
        assertEquals("X", plcWCHAR.getString());
        assertFalse(plcWCHAR.isNullable());
    }

    @Test
    void testConstructorWithStringWithWhitespace() {
        PlcWCHAR plcWCHAR = new PlcWCHAR("  Y  ");
        assertEquals("Y", plcWCHAR.getString());
        assertFalse(plcWCHAR.isNullable());
    }

    @Test
    void testConstructorWithSpace() {
        PlcWCHAR plcWCHAR = new PlcWCHAR("   ");
        assertEquals(" ", plcWCHAR.getString());
        assertFalse(plcWCHAR.isNullable());
    }

    @Test
    void testGetPlcValueType() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('Z');
        assertEquals(PlcValueType.WCHAR, plcWCHAR.getPlcValueType());
    }

    @Test
    void testIsBoolean() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('A');
        assertTrue(plcWCHAR.isBoolean());
    }

    @Test
    void testGetBoolean() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('A');
        assertTrue(plcWCHAR.getBoolean());

        PlcWCHAR plcWCHAR2 = new PlcWCHAR((char) 0);
        assertFalse(plcWCHAR2.getBoolean());
    }

    @Test
    void testGetByte() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('A');
        assertEquals(65, plcWCHAR.getByte());
    }

    @Test
    void testGetShort() {
        PlcWCHAR plcWCHAR = new PlcWCHAR((char) 12345);
        assertEquals(12345, plcWCHAR.getShort());
    }

    @Test
    void testGetInteger() {
        PlcWCHAR plcWCHAR = new PlcWCHAR((char) 30000);
        assertEquals(30000, plcWCHAR.getInteger());
    }

    @Test
    void testGetLong() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('D');
        assertEquals(68L, plcWCHAR.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('E');
        assertEquals(BigInteger.valueOf(69), plcWCHAR.getBigInteger());
    }

    @Test
    void testGetString() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('Z');
        assertEquals("Z", plcWCHAR.getString());
    }

    @Test
    void testToString() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('M');
        assertEquals("M", plcWCHAR.toString());
    }

    @Test
    void testGetObject() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('Q');
        assertEquals("Q", plcWCHAR.getObject());
    }

    @Test
    void testGetBytes() {
        PlcWCHAR plcWCHAR = new PlcWCHAR((char) 0x4142);
        byte[] bytes = plcWCHAR.getBytes();
        assertEquals(2, bytes.length);
        assertEquals(0x41, bytes[0]);
        assertEquals(0x42, bytes[1]);
    }

    @Test
    void testOfWithPlcWCHAR() {
        PlcWCHAR original = new PlcWCHAR('A');
        PlcWCHAR copy = PlcWCHAR.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithString() {
        PlcWCHAR plcWCHAR = PlcWCHAR.of("Y");
        assertEquals("Y", plcWCHAR.getString());
    }

    @Test
    void testOfWithBoolean() {
        PlcWCHAR plcWCHAR = PlcWCHAR.of(true);
        assertEquals('T', plcWCHAR.getString().charAt(0));
    }

    @Test
    void testOfWithByte() {
        PlcWCHAR plcWCHAR = PlcWCHAR.of((byte) 65);
        assertEquals("A", plcWCHAR.getString());
    }

    @Test
    void testOfWithInteger() {
        PlcWCHAR plcWCHAR = PlcWCHAR.of(30000);
        assertEquals(30000, plcWCHAR.getInteger());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testLongConstructor() {
        PlcWCHAR value = new PlcWCHAR(30000L);
        assertEquals(30000, value.getInteger());
    }

    @Test
    void testLongConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((Long) null));
    }

    @Test
    void testFloatConstructor() {
        PlcWCHAR value = new PlcWCHAR(12345.0f);
        assertEquals(12345, value.getInteger());
    }

    @Test
    void testFloatConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((Float) null));
    }

    @Test
    void testDoubleConstructor() {
        PlcWCHAR value = new PlcWCHAR(20000.0);
        assertEquals(20000, value.getInteger());
    }

    @Test
    void testDoubleConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((Double) null));
    }

    @Test
    void testBigIntegerConstructor() {
        PlcWCHAR value = new PlcWCHAR(BigInteger.valueOf(15000));
        assertEquals(15000, value.getInteger());
    }

    @Test
    void testBigIntegerConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((BigInteger) null));
    }

    @Test
    void testBigDecimalConstructor() {
        PlcWCHAR value = new PlcWCHAR(BigDecimal.valueOf(10000));
        assertEquals(10000, value.getInteger());
    }

    @Test
    void testBigDecimalConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((BigDecimal) null));
    }

    // ========== Additional Getter Tests ==========

    @Test
    void testGetFloat() {
        PlcWCHAR value = new PlcWCHAR(12000);
        assertEquals(12000.0f, value.getFloat(), 0.001f);
    }

    @Test
    void testGetDouble() {
        PlcWCHAR value = new PlcWCHAR(15000);
        assertEquals(15000.0, value.getDouble(), 0.001);
    }

    @Test
    void testGetBigDecimal() {
        PlcWCHAR value = new PlcWCHAR(20000);
        assertEquals(0, BigDecimal.valueOf(20000).compareTo(value.getBigDecimal()));
    }

    // ========== Additional Is* Methods Tests ==========

    @Test
    void testIsLong() {
        PlcWCHAR value = new PlcWCHAR(25000);
        assertTrue(value.isLong());
    }

    @Test
    void testIsBigInteger() {
        PlcWCHAR value = new PlcWCHAR(30000);
        assertTrue(value.isBigInteger());
    }

    @Test
    void testIsFloat() {
        PlcWCHAR value = new PlcWCHAR(12000);
        assertTrue(value.isFloat());
    }

    @Test
    void testIsDouble() {
        PlcWCHAR value = new PlcWCHAR(15000);
        assertTrue(value.isDouble());
    }

    @Test
    void testIsBigDecimal() {
        PlcWCHAR value = new PlcWCHAR(20000);
        assertTrue(value.isBigDecimal());
    }

    @Test
    void testIsByte() {
        PlcWCHAR value = new PlcWCHAR(100);
        assertTrue(value.isByte());
    }

    @Test
    void testIsShort() {
        PlcWCHAR value = new PlcWCHAR(1000);
        assertTrue(value.isShort());
    }

    @Test
    void testIsInteger() {
        PlcWCHAR value = new PlcWCHAR(30000);
        assertTrue(value.isInteger());
    }

    // ========== Edge Case Tests ==========

    @Test
    void testIntegerConstructorOutOfRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWCHAR(70000));
    }

    @Test
    void testLongConstructorOutOfRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWCHAR(70000L));
    }

    @Test
    void testByteConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((Byte) null));
    }

    @Test
    void testShortConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((Short) null));
    }

    @Test
    void testIntegerConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((Integer) null));
    }

    @Test
    void testMinValue() {
        PlcWCHAR value = new PlcWCHAR(0);
        assertEquals(0, value.getInteger());
    }

    @Test
    void testMaxValue() {
        PlcWCHAR value = new PlcWCHAR(65535);
        assertEquals(65535, value.getInteger());
    }

    @Test
    void testNegativeValue() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWCHAR(-1));
    }

    // ========== Additional Tests for Uncovered Code Paths ==========

    @Test
    void testCharacterConstructor() {
        PlcWCHAR plcWCHAR = new PlcWCHAR('A');
        assertEquals("A", plcWCHAR.getString());
    }

    @Test
    void testCharacterConstructorNull() {
        assertThrows(NullPointerException.class, () -> new PlcWCHAR((Character) null));
    }

    @Test
    void testShortConstructorOutOfRange() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWCHAR((short) -1));
    }

    @Test
    void testBigDecimalWithScale() {
        // BigDecimal with scale > 0 should throw exception
        assertThrows(PlcInvalidTagException.class, () -> new PlcWCHAR(BigDecimal.valueOf(30000.5)));
    }

    @Test
    void testBigIntegerIsNullable() {
        PlcWCHAR value = new PlcWCHAR(BigInteger.valueOf(30000));
        assertTrue(value.isNullable());
    }

    @Test
    void testBigDecimalIsNullable() {
        PlcWCHAR value = new PlcWCHAR(BigDecimal.valueOf(30000));
        assertTrue(value.isNullable());
    }

    @Test
    void testStringConstructorNull() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWCHAR((String) null));
    }

    @Test
    void testIsString() {
        PlcWCHAR value = new PlcWCHAR(30000);
        assertTrue(value.isString());
    }

    @Test
    void testGetRaw() {
        PlcWCHAR value = new PlcWCHAR('A');
        byte[] raw = value.getRaw();
        byte[] bytes = value.getBytes();
        assertArrayEquals(bytes, raw);
    }

    @Test
    void testIsByteReturnsFalse() {
        // Value > Byte.MAX_VALUE (127) should return false for isByte()
        PlcWCHAR value = new PlcWCHAR(200);
        assertFalse(value.isByte());
    }

    @Test
    void testIsByteReturnsTrue() {
        // Value <= Byte.MAX_VALUE (127) should return true for isByte()
        PlcWCHAR value = new PlcWCHAR(100);
        assertTrue(value.isByte());
    }

    @Test
    void testIsShortReturnsFalse() {
        // Value > Short.MAX_VALUE (32767) should return false for isShort()
        PlcWCHAR value = new PlcWCHAR(40000);
        assertFalse(value.isShort());
    }

    @Test
    void testIsShortReturnsTrue() {
        // Value <= Short.MAX_VALUE (32767) should return true for isShort()
        PlcWCHAR value = new PlcWCHAR(30000);
        assertTrue(value.isShort());
    }
}
