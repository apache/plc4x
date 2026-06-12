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

public class PlcWORDTest {

    @Test
    void testConstructorWithBoolean() {
        PlcWORD plcWORD = new PlcWORD(true);
        assertEquals(1, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());

        plcWORD = new PlcWORD(false);
        assertEquals(0, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithValidByte() {
        PlcWORD plcWORD = new PlcWORD((byte) 42);
        assertEquals(42, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithValidShort() {
        PlcWORD plcWORD = new PlcWORD((short) 1000);
        assertEquals(1000, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithValidInteger() {
        PlcWORD plcWORD = new PlcWORD(32000);
        assertEquals(32000, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithMinValue() {
        PlcWORD plcWORD = new PlcWORD(PlcWORD.MIN_VALUE);
        assertEquals(0, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithMaxValue() {
        PlcWORD plcWORD = new PlcWORD(PlcWORD.MAX_VALUE);
        assertEquals(65535, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeIntegerNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(-1));
    }

    @Test
    void testConstructorWithOutOfRangeIntegerTooLarge() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(65536));
    }

    @Test
    void testConstructorWithValidLong() {
        PlcWORD plcWORD = new PlcWORD(30000L);
        assertEquals(30000, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeLong() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(-1L));
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(65536L));
    }

    @Test
    void testConstructorWithValidFloat() {
        PlcWORD plcWORD = new PlcWORD(100.5f);
        assertEquals(100, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeFloat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(-1.0f));
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(65536.0f));
    }

    @Test
    void testConstructorWithValidDouble() {
        PlcWORD plcWORD = new PlcWORD(200.7);
        assertEquals(200, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeDouble() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(-1.0));
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(65536.0));
    }

    @Test
    void testConstructorWithValidBigInteger() {
        PlcWORD plcWORD = new PlcWORD(BigInteger.valueOf(50000));
        assertEquals(50000, plcWORD.getInteger());
        assertTrue(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeBigInteger() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(BigInteger.valueOf(-1)));
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(BigInteger.valueOf(65536)));
    }

    @Test
    void testConstructorWithValidBigDecimal() {
        PlcWORD plcWORD = new PlcWORD(BigDecimal.valueOf(40000.5));
        assertEquals(40000, plcWORD.getInteger());
        assertTrue(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeBigDecimal() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(BigDecimal.valueOf(-1)));
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD(BigDecimal.valueOf(65536)));
    }

    @Test
    void testConstructorWithValidString() {
        PlcWORD plcWORD = new PlcWORD("12345");
        assertEquals(12345, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithValidStringWithWhitespace() {
        PlcWORD plcWORD = new PlcWORD("  999  ");
        assertEquals(999, plcWORD.getInteger());
        assertFalse(plcWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeString() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD("-1"));
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD("65536"));
    }

    @Test
    void testConstructorWithInvalidString() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcWORD("not a number"));
    }

    @Test
    void testGetPlcValueType() {
        PlcWORD plcWORD = new PlcWORD(100);
        assertEquals(PlcValueType.WORD, plcWORD.getPlcValueType());
    }

    @Test
    void testIsBoolean() {
        PlcWORD plcWORD = new PlcWORD(1);
        assertTrue(plcWORD.isBoolean());
    }

    @Test
    void testGetBooleanTrue() {
        PlcWORD plcWORD = new PlcWORD(1);
        assertTrue(plcWORD.getBoolean());
    }

    @Test
    void testGetBooleanFalse() {
        PlcWORD plcWORD = new PlcWORD(0);
        assertFalse(plcWORD.getBoolean());
    }

    @Test
    void testGetBooleanArray() {
        PlcWORD plcWORD = new PlcWORD(0b1010101010101010);
        boolean[] booleanArray = plcWORD.getBooleanArray();
        assertEquals(16, booleanArray.length);
        assertFalse(booleanArray[0]);
        assertTrue(booleanArray[1]);
        assertFalse(booleanArray[2]);
        assertTrue(booleanArray[3]);
    }

    @Test
    void testGetByte() {
        PlcWORD plcWORD = new PlcWORD(100);
        assertEquals(100, plcWORD.getByte());
    }

    @Test
    void testGetShort() {
        PlcWORD plcWORD = new PlcWORD(30000);
        assertEquals(30000, plcWORD.getShort());
    }

    @Test
    void testGetInteger() {
        PlcWORD plcWORD = new PlcWORD(50000);
        assertEquals(50000, plcWORD.getInteger());
    }

    @Test
    void testGetLong() {
        PlcWORD plcWORD = new PlcWORD(40000);
        assertEquals(40000L, plcWORD.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcWORD plcWORD = new PlcWORD(60000);
        assertEquals(BigInteger.valueOf(60000), plcWORD.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcWORD plcWORD = new PlcWORD(1234);
        assertEquals(1234.0f, plcWORD.getFloat());
    }

    @Test
    void testGetDouble() {
        PlcWORD plcWORD = new PlcWORD(5678);
        assertEquals(5678.0, plcWORD.getDouble());
    }

    @Test
    void testGetBigDecimal() {
        PlcWORD plcWORD = new PlcWORD(9999);
        assertEquals(BigDecimal.valueOf(9999.0f), plcWORD.getBigDecimal());
    }

    @Test
    void testIsString() {
        PlcWORD plcWORD = new PlcWORD(123);
        assertTrue(plcWORD.isString());
    }

    @Test
    void testGetString() {
        PlcWORD plcWORD = new PlcWORD(54321);
        assertEquals("54321", plcWORD.getString());
    }

    @Test
    void testToString() {
        PlcWORD plcWORD = new PlcWORD(11111);
        assertEquals("11111", plcWORD.toString());
    }

    @Test
    void testGetBytes() {
        PlcWORD plcWORD = new PlcWORD(0x1234);
        byte[] bytes = plcWORD.getBytes();
        assertEquals(2, bytes.length);
        assertEquals(0x12, bytes[0]);
        assertEquals(0x34, bytes[1]);
    }

    @Test
    void testGetRaw() {
        PlcWORD plcWORD = new PlcWORD(0xABCD);
        byte[] raw = plcWORD.getRaw();
        assertEquals(2, raw.length);
        assertEquals((byte) 0xAB, raw[0]);
        assertEquals((byte) 0xCD, raw[1]);
    }

    @Test
    void testOfWithPlcWORD() {
        PlcWORD original = new PlcWORD(1000);
        PlcWORD copy = PlcWORD.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithBoolean() {
        PlcWORD plcWORD = PlcWORD.of(true);
        assertEquals(1, plcWORD.getInteger());
    }

    @Test
    void testOfWithByte() {
        PlcWORD plcWORD = PlcWORD.of((byte) 100);
        assertEquals(100, plcWORD.getInteger());
    }

    @Test
    void testOfWithShort() {
        PlcWORD plcWORD = PlcWORD.of((short) 2000);
        assertEquals(2000, plcWORD.getInteger());
    }

    @Test
    void testOfWithInteger() {
        PlcWORD plcWORD = PlcWORD.of(30000);
        assertEquals(30000, plcWORD.getInteger());
    }

    @Test
    void testOfWithLong() {
        PlcWORD plcWORD = PlcWORD.of(40000L);
        assertEquals(40000, plcWORD.getInteger());
    }

    @Test
    void testOfWithFloat() {
        PlcWORD plcWORD = PlcWORD.of(500.5f);
        assertEquals(500, plcWORD.getInteger());
    }

    @Test
    void testOfWithDouble() {
        PlcWORD plcWORD = PlcWORD.of(600.7);
        assertEquals(600, plcWORD.getInteger());
    }

    @Test
    void testOfWithBigInteger() {
        PlcWORD plcWORD = PlcWORD.of(BigInteger.valueOf(50000));
        assertEquals(50000, plcWORD.getInteger());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcWORD plcWORD = PlcWORD.of(BigDecimal.valueOf(55000));
        assertEquals(55000, plcWORD.getInteger());
    }

    @Test
    void testOfWithString() {
        PlcWORD plcWORD = PlcWORD.of("12345");
        assertEquals(12345, plcWORD.getInteger());
    }
}