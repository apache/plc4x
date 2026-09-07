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

public class PlcDWORDTest {

    @Test
    void testConstructorWithBoolean() {
        PlcDWORD plcDWORD = new PlcDWORD(true);
        assertEquals(1L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());

        plcDWORD = new PlcDWORD(false);
        assertEquals(0L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithValidByte() {
        PlcDWORD plcDWORD = new PlcDWORD((byte) 42);
        assertEquals(42L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithValidShort() {
        PlcDWORD plcDWORD = new PlcDWORD((short) 1000);
        assertEquals(1000L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithValidInteger() {
        PlcDWORD plcDWORD = new PlcDWORD(1000000);
        assertEquals(1000000L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithValidLong() {
        PlcDWORD plcDWORD = new PlcDWORD(3000000000L);
        assertEquals(3000000000L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithMinValue() {
        PlcDWORD plcDWORD = new PlcDWORD(PlcDWORD.MIN_VALUE);
        assertEquals(0L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithMaxValue() {
        PlcDWORD plcDWORD = new PlcDWORD(PlcDWORD.MAX_VALUE);
        assertEquals(4294967295L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeLongNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(-1L));
    }

    @Test
    void testConstructorWithOutOfRangeLongTooLarge() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(4294967296L));
    }

    @Test
    void testConstructorWithValidFloat() {
        PlcDWORD plcDWORD = new PlcDWORD(1000.5f);
        assertEquals(1000L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeFloat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(-1.0f));
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(5000000000.0f));
    }

    @Test
    void testConstructorWithValidDouble() {
        PlcDWORD plcDWORD = new PlcDWORD(2000.7);
        assertEquals(2000L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeDouble() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(-1.0));
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(4294967296.0));
    }

    @Test
    void testConstructorWithValidBigInteger() {
        PlcDWORD plcDWORD = new PlcDWORD(BigInteger.valueOf(2000000000));
        assertEquals(2000000000L, plcDWORD.getLong());
        assertTrue(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeBigInteger() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(BigInteger.valueOf(-1)));
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(BigInteger.valueOf(4294967296L)));
    }

    @Test
    void testConstructorWithValidBigDecimal() {
        PlcDWORD plcDWORD = new PlcDWORD(BigDecimal.valueOf(3000000000.5));
        assertEquals(3000000000L, plcDWORD.getLong());
        assertTrue(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeBigDecimal() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(BigDecimal.valueOf(-1)));
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD(BigDecimal.valueOf(4294967296L)));
    }

    @Test
    void testConstructorWithValidString() {
        PlcDWORD plcDWORD = new PlcDWORD("12345678");
        assertEquals(12345678L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithValidStringWithWhitespace() {
        PlcDWORD plcDWORD = new PlcDWORD("  999999  ");
        assertEquals(999999L, plcDWORD.getLong());
        assertFalse(plcDWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeString() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD("-1"));
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD("4294967296"));
    }

    @Test
    void testConstructorWithInvalidString() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcDWORD("not a number"));
    }

    @Test
    void testGetPlcValueType() {
        PlcDWORD plcDWORD = new PlcDWORD(1000L);
        assertEquals(PlcValueType.DWORD, plcDWORD.getPlcValueType());
    }

    @Test
    void testIsBoolean() {
        PlcDWORD plcDWORD = new PlcDWORD(1L);
        assertTrue(plcDWORD.isBoolean());
    }

    @Test
    void testGetBooleanTrue() {
        PlcDWORD plcDWORD = new PlcDWORD(1L);
        assertTrue(plcDWORD.getBoolean());
    }

    @Test
    void testGetBooleanFalse() {
        PlcDWORD plcDWORD = new PlcDWORD(0L);
        assertFalse(plcDWORD.getBoolean());
    }

    @Test
    void testGetBooleanArray() {
        PlcDWORD plcDWORD = new PlcDWORD(0b10101010101010101010101010101010L);
        boolean[] booleanArray = plcDWORD.getBooleanArray();
        assertEquals(32, booleanArray.length);
        assertFalse(booleanArray[0]);
        assertTrue(booleanArray[1]);
        assertFalse(booleanArray[2]);
        assertTrue(booleanArray[3]);
    }

    @Test
    void testGetByte() {
        PlcDWORD plcDWORD = new PlcDWORD(100L);
        assertEquals(100, plcDWORD.getByte());
    }

    @Test
    void testGetShort() {
        PlcDWORD plcDWORD = new PlcDWORD(30000L);
        assertEquals(30000, plcDWORD.getShort());
    }

    @Test
    void testGetInteger() {
        PlcDWORD plcDWORD = new PlcDWORD(2000000000L);
        assertEquals(2000000000, plcDWORD.getInteger());
    }

    @Test
    void testGetLong() {
        PlcDWORD plcDWORD = new PlcDWORD(4000000000L);
        assertEquals(4000000000L, plcDWORD.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcDWORD plcDWORD = new PlcDWORD(3000000000L);
        assertEquals(BigInteger.valueOf(3000000000L), plcDWORD.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcDWORD plcDWORD = new PlcDWORD(123456L);
        assertEquals(123456.0f, plcDWORD.getFloat());
    }

    @Test
    void testGetDouble() {
        PlcDWORD plcDWORD = new PlcDWORD(7890123L);
        assertEquals(7890123.0, plcDWORD.getDouble());
    }

    @Test
    void testGetBigDecimal() {
        PlcDWORD plcDWORD = new PlcDWORD(999999L);
        assertEquals(BigDecimal.valueOf(999999.0f), plcDWORD.getBigDecimal());
    }

    @Test
    void testIsString() {
        PlcDWORD plcDWORD = new PlcDWORD(123L);
        assertTrue(plcDWORD.isString());
    }

    @Test
    void testGetString() {
        PlcDWORD plcDWORD = new PlcDWORD(87654321L);
        assertEquals("87654321", plcDWORD.getString());
    }

    @Test
    void testToString() {
        PlcDWORD plcDWORD = new PlcDWORD(11111111L);
        assertEquals("11111111", plcDWORD.toString());
    }

    @Test
    void testGetBytes() {
        PlcDWORD plcDWORD = new PlcDWORD(0x12345678L);
        byte[] bytes = plcDWORD.getBytes();
        assertEquals(4, bytes.length);
        assertEquals(0x12, bytes[0]);
        assertEquals(0x34, bytes[1]);
        assertEquals(0x56, bytes[2]);
        assertEquals(0x78, bytes[3]);
    }

    @Test
    void testGetRaw() {
        PlcDWORD plcDWORD = new PlcDWORD(0xABCDEF01L);
        byte[] raw = plcDWORD.getRaw();
        assertEquals(4, raw.length);
        assertEquals((byte) 0xAB, raw[0]);
        assertEquals((byte) 0xCD, raw[1]);
        assertEquals((byte) 0xEF, raw[2]);
        assertEquals((byte) 0x01, raw[3]);
    }

    @Test
    void testOfWithPlcDWORD() {
        PlcDWORD original = new PlcDWORD(1000000L);
        PlcDWORD copy = PlcDWORD.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithBoolean() {
        PlcDWORD plcDWORD = PlcDWORD.of(true);
        assertEquals(1L, plcDWORD.getLong());
    }

    @Test
    void testOfWithByte() {
        PlcDWORD plcDWORD = PlcDWORD.of((byte) 100);
        assertEquals(100L, plcDWORD.getLong());
    }

    @Test
    void testOfWithShort() {
        PlcDWORD plcDWORD = PlcDWORD.of((short) 2000);
        assertEquals(2000L, plcDWORD.getLong());
    }

    @Test
    void testOfWithInteger() {
        PlcDWORD plcDWORD = PlcDWORD.of(1000000);
        assertEquals(1000000L, plcDWORD.getLong());
    }

    @Test
    void testOfWithLong() {
        PlcDWORD plcDWORD = PlcDWORD.of(3000000000L);
        assertEquals(3000000000L, plcDWORD.getLong());
    }

    @Test
    void testOfWithFloat() {
        PlcDWORD plcDWORD = PlcDWORD.of(500.5f);
        assertEquals(500L, plcDWORD.getLong());
    }

    @Test
    void testOfWithDouble() {
        PlcDWORD plcDWORD = PlcDWORD.of(600.7);
        assertEquals(600L, plcDWORD.getLong());
    }

    @Test
    void testOfWithBigInteger() {
        PlcDWORD plcDWORD = PlcDWORD.of(BigInteger.valueOf(2000000000));
        assertEquals(2000000000L, plcDWORD.getLong());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcDWORD plcDWORD = PlcDWORD.of(BigDecimal.valueOf(3000000000L));
        assertEquals(3000000000L, plcDWORD.getLong());
    }

    @Test
    void testOfWithString() {
        PlcDWORD plcDWORD = PlcDWORD.of("123456789");
        assertEquals(123456789L, plcDWORD.getLong());
    }
}
