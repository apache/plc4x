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

public class PlcLWORDTest {

    @Test
    void testConstructorWithBoolean() {
        PlcLWORD plcLWORD = new PlcLWORD(true);
        assertEquals(BigInteger.ONE, plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());

        plcLWORD = new PlcLWORD(false);
        assertEquals(BigInteger.ZERO, plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithValidByte() {
        PlcLWORD plcLWORD = new PlcLWORD((byte) 42);
        assertEquals(BigInteger.valueOf(42), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithValidShort() {
        PlcLWORD plcLWORD = new PlcLWORD((short) 1000);
        assertEquals(BigInteger.valueOf(1000), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithValidInteger() {
        PlcLWORD plcLWORD = new PlcLWORD(1000000);
        assertEquals(BigInteger.valueOf(1000000), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithValidLong() {
        PlcLWORD plcLWORD = new PlcLWORD(9000000000000000000L);
        assertEquals(BigInteger.valueOf(9000000000000000000L), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithMinValue() {
        PlcLWORD plcLWORD = new PlcLWORD(PlcLWORD.MIN_VALUE);
        assertEquals(BigInteger.ZERO, plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithMaxValue() {
        PlcLWORD plcLWORD = new PlcLWORD(PlcLWORD.MAX_VALUE);
        assertEquals(new BigInteger("18446744073709551615"), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeBigIntegerNegative() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD(BigInteger.valueOf(-1)));
    }

    @Test
    void testConstructorWithOutOfRangeBigIntegerTooLarge() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD(new BigInteger("18446744073709551616")));
    }

    @Test
    void testConstructorWithValidFloat() {
        PlcLWORD plcLWORD = new PlcLWORD(1000.5f);
        assertEquals(BigInteger.valueOf(1000), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeFloat() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD(-1.0f));
    }

    @Test
    void testConstructorWithValidDouble() {
        PlcLWORD plcLWORD = new PlcLWORD(2000.7);
        assertEquals(BigInteger.valueOf(2000), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeDouble() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD(-1.0));
    }

    @Test
    void testConstructorWithValidBigInteger() {
        PlcLWORD plcLWORD = new PlcLWORD(new BigInteger("10000000000000000000"));
        assertEquals(new BigInteger("10000000000000000000"), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithValidBigDecimal() {
        PlcLWORD plcLWORD = new PlcLWORD(new BigDecimal("10000000000000000000.5"));
        assertEquals(new BigInteger("10000000000000000000"), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeBigDecimal() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD(BigDecimal.valueOf(-1)));
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD(new BigDecimal("18446744073709551616")));
    }

    @Test
    void testConstructorWithValidString() {
        PlcLWORD plcLWORD = new PlcLWORD("123456789012345");
        assertEquals(new BigInteger("123456789012345"), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithValidStringWithWhitespace() {
        PlcLWORD plcLWORD = new PlcLWORD("  999999999999  ");
        assertEquals(new BigInteger("999999999999"), plcLWORD.getBigInteger());
        assertFalse(plcLWORD.isNullable());
    }

    @Test
    void testConstructorWithOutOfRangeString() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD("-1"));
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD("18446744073709551616"));
    }

    @Test
    void testConstructorWithInvalidString() {
        assertThrows(PlcInvalidTagException.class, () -> new PlcLWORD("not a number"));
    }

    @Test
    void testGetPlcValueType() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(1000));
        assertEquals(PlcValueType.LWORD, plcLWORD.getPlcValueType());
    }

    @Test
    void testIsBoolean() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.ONE);
        assertTrue(plcLWORD.isBoolean());
    }

    @Test
    void testGetBooleanTrue() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.ONE);
        assertTrue(plcLWORD.getBoolean());
    }

    @Test
    void testGetBooleanFalse() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.ZERO);
        assertFalse(plcLWORD.getBoolean());
    }

    @Test
    void testGetBooleanArray() {
        PlcLWORD plcLWORD = new PlcLWORD(new BigInteger("170"));  // 0b10101010 in lower byte
        boolean[] booleanArray = plcLWORD.getBooleanArray();
        assertEquals(64, booleanArray.length);
        // BigInteger.toByteArray() includes a sign byte, so bits are at positions 8-15
        assertFalse(booleanArray[8]);
        assertTrue(booleanArray[9]);
        assertFalse(booleanArray[10]);
        assertTrue(booleanArray[11]);
    }

    @Test
    void testGetByte() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(100));
        assertEquals(100, plcLWORD.getByte());
    }

    @Test
    void testGetShort() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(30000));
        assertEquals(30000, plcLWORD.getShort());
    }

    @Test
    void testGetInteger() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(2000000000));
        assertEquals(2000000000, plcLWORD.getInteger());
    }

    @Test
    void testGetLong() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(9000000000000000000L));
        assertEquals(9000000000000000000L, plcLWORD.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcLWORD plcLWORD = new PlcLWORD(new BigInteger("10000000000000000000"));
        assertEquals(new BigInteger("10000000000000000000"), plcLWORD.getBigInteger());
    }

    @Test
    void testGetFloat() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(123456));
        assertEquals(123456.0f, plcLWORD.getFloat());
    }

    @Test
    void testGetDouble() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(7890123));
        assertEquals(7890123.0, plcLWORD.getDouble());
    }

    @Test
    void testGetBigDecimal() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(999999));
        assertEquals(new BigDecimal("999999"), plcLWORD.getBigDecimal());
    }

    @Test
    void testIsString() {
        PlcLWORD plcLWORD = new PlcLWORD(BigInteger.valueOf(123));
        assertTrue(plcLWORD.isString());
    }

    @Test
    void testGetString() {
        PlcLWORD plcLWORD = new PlcLWORD(new BigInteger("87654321098765432"));
        assertEquals("87654321098765432", plcLWORD.getString());
    }

    @Test
    void testToString() {
        PlcLWORD plcLWORD = new PlcLWORD(new BigInteger("11111111111111111"));
        assertEquals("11111111111111111", plcLWORD.toString());
    }

    @Test
    void testGetBytes() {
        PlcLWORD plcLWORD = new PlcLWORD(new BigInteger("1234567890123456", 16));
        byte[] bytes = plcLWORD.getBytes();
        assertEquals(8, bytes.length);
        assertEquals((byte) 0x12, bytes[0]);
        assertEquals((byte) 0x34, bytes[1]);
        assertEquals((byte) 0x56, bytes[2]);
        assertEquals((byte) 0x78, bytes[3]);
        assertEquals((byte) 0x90, bytes[4]);
        assertEquals((byte) 0x12, bytes[5]);
        assertEquals((byte) 0x34, bytes[6]);
        assertEquals((byte) 0x56, bytes[7]);
    }

    @Test
    void testGetRaw() {
        PlcLWORD plcLWORD = new PlcLWORD(new BigInteger("ABCDEF0123456789", 16));
        byte[] raw = plcLWORD.getRaw();
        assertEquals(8, raw.length);
        assertEquals((byte) 0xAB, raw[0]);
        assertEquals((byte) 0xCD, raw[1]);
        assertEquals((byte) 0xEF, raw[2]);
        assertEquals((byte) 0x01, raw[3]);
        assertEquals((byte) 0x23, raw[4]);
        assertEquals((byte) 0x45, raw[5]);
        assertEquals((byte) 0x67, raw[6]);
        assertEquals((byte) 0x89, raw[7]);
    }

    @Test
    void testOfWithPlcLWORD() {
        PlcLWORD original = new PlcLWORD(BigInteger.valueOf(1000000));
        PlcLWORD copy = PlcLWORD.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithBoolean() {
        PlcLWORD plcLWORD = PlcLWORD.of(true);
        assertEquals(BigInteger.ONE, plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithByte() {
        PlcLWORD plcLWORD = PlcLWORD.of((byte) 100);
        assertEquals(BigInteger.valueOf(100), plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithShort() {
        PlcLWORD plcLWORD = PlcLWORD.of((short) 2000);
        assertEquals(BigInteger.valueOf(2000), plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithInteger() {
        PlcLWORD plcLWORD = PlcLWORD.of(1000000);
        assertEquals(BigInteger.valueOf(1000000), plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithLong() {
        PlcLWORD plcLWORD = PlcLWORD.of(9000000000000000000L);
        assertEquals(BigInteger.valueOf(9000000000000000000L), plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithFloat() {
        PlcLWORD plcLWORD = PlcLWORD.of(500.5f);
        assertEquals(BigInteger.valueOf(500), plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithDouble() {
        PlcLWORD plcLWORD = PlcLWORD.of(600.7);
        assertEquals(BigInteger.valueOf(600), plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithBigInteger() {
        PlcLWORD plcLWORD = PlcLWORD.of(new BigInteger("10000000000000000000"));
        assertEquals(new BigInteger("10000000000000000000"), plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcLWORD plcLWORD = PlcLWORD.of(new BigDecimal("10000000000000000000"));
        assertEquals(new BigInteger("10000000000000000000"), plcLWORD.getBigInteger());
    }

    @Test
    void testOfWithString() {
        PlcLWORD plcLWORD = PlcLWORD.of("123456789012345");
        assertEquals(new BigInteger("123456789012345"), plcLWORD.getBigInteger());
    }
}
