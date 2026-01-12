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

public class PlcSTRINGTest {

    @Test
    void testConstructorWithString() {
        PlcSTRING plcSTRING = new PlcSTRING("Hello World");
        assertEquals("Hello World", plcSTRING.getString());
        assertFalse(plcSTRING.isNullable());
    }

    @Test
    void testGetPlcValueType() {
        PlcSTRING plcSTRING = new PlcSTRING("test");
        assertEquals(PlcValueType.STRING, plcSTRING.getPlcValueType());
    }

    @Test
    void testIsString() {
        PlcSTRING plcSTRING = new PlcSTRING("test");
        assertTrue(plcSTRING.isString());
    }

    @Test
    void testGetString() {
        PlcSTRING plcSTRING = new PlcSTRING("test string");
        assertEquals("test string", plcSTRING.getString());
    }

    @Test
    void testIsBooleanValid() {
        PlcSTRING plcSTRING = new PlcSTRING("true");
        assertTrue(plcSTRING.isBoolean());
    }

    @Test
    void testGetBoolean() {
        PlcSTRING plcSTRING = new PlcSTRING("true");
        assertTrue(plcSTRING.getBoolean());

        PlcSTRING plcSTRING2 = new PlcSTRING("false");
        assertFalse(plcSTRING2.getBoolean());
    }

    @Test
    void testIsByteValid() {
        PlcSTRING plcSTRING = new PlcSTRING("123");
        assertTrue(plcSTRING.isByte());
    }

    @Test
    void testGetByte() {
        PlcSTRING plcSTRING = new PlcSTRING("42");
        assertEquals(42, plcSTRING.getByte());
    }

    @Test
    void testIsShortValid() {
        PlcSTRING plcSTRING = new PlcSTRING("12345");
        assertTrue(plcSTRING.isShort());
    }

    @Test
    void testGetShort() {
        PlcSTRING plcSTRING = new PlcSTRING("1000");
        assertEquals(1000, plcSTRING.getShort());
    }

    @Test
    void testIsIntegerValid() {
        PlcSTRING plcSTRING = new PlcSTRING("123456");
        assertTrue(plcSTRING.isInteger());
    }

    @Test
    void testGetInteger() {
        PlcSTRING plcSTRING = new PlcSTRING("100000");
        assertEquals(100000, plcSTRING.getInteger());
    }

    @Test
    void testIsLongValid() {
        PlcSTRING plcSTRING = new PlcSTRING("9223372036854775807");
        assertTrue(plcSTRING.isLong());
    }

    @Test
    void testGetLong() {
        PlcSTRING plcSTRING = new PlcSTRING("9000000000");
        assertEquals(9000000000L, plcSTRING.getLong());
    }

    @Test
    void testIsBigIntegerValid() {
        PlcSTRING plcSTRING = new PlcSTRING("12345678901234567890");
        assertTrue(plcSTRING.isBigInteger());
    }

    @Test
    void testGetBigInteger() {
        PlcSTRING plcSTRING = new PlcSTRING("12345678901234567890");
        assertEquals(new BigInteger("12345678901234567890"), plcSTRING.getBigInteger());
    }

    @Test
    void testIsFloatValid() {
        PlcSTRING plcSTRING = new PlcSTRING("123.45");
        assertTrue(plcSTRING.isFloat());
    }

    @Test
    void testGetFloat() {
        PlcSTRING plcSTRING = new PlcSTRING("123.45");
        assertEquals(123.45f, plcSTRING.getFloat(), 0.001);
    }

    @Test
    void testIsDoubleValid() {
        PlcSTRING plcSTRING = new PlcSTRING("123.456789");
        assertTrue(plcSTRING.isDouble());
    }

    @Test
    void testGetDouble() {
        PlcSTRING plcSTRING = new PlcSTRING("123.456789");
        assertEquals(123.456789, plcSTRING.getDouble(), 0.000001);
    }

    @Test
    void testIsBigDecimalValid() {
        PlcSTRING plcSTRING = new PlcSTRING("123.456789012345");
        assertTrue(plcSTRING.isBigDecimal());
    }

    @Test
    void testGetBigDecimal() {
        PlcSTRING plcSTRING = new PlcSTRING("123.456");
        assertEquals(new BigDecimal("123.456"), plcSTRING.getBigDecimal());
    }

    @Test
    void testGetLength() {
        PlcSTRING plcSTRING = new PlcSTRING("Hello");
        assertEquals(5, plcSTRING.getLength());
    }

    @Test
    void testToString() {
        PlcSTRING plcSTRING = new PlcSTRING("test");
        assertEquals("test", plcSTRING.toString());
    }

    @Test
    void testOfWithPlcSTRING() {
        PlcSTRING original = new PlcSTRING("test");
        PlcSTRING copy = PlcSTRING.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithString() {
        PlcSTRING plcSTRING = PlcSTRING.of("test");
        assertEquals("test", plcSTRING.getString());
    }

    @Test
    void testOfWithObject() {
        PlcSTRING plcSTRING = PlcSTRING.of(123);
        assertEquals("123", plcSTRING.getString());
    }

    @Test
    void testIsBooleanInvalid() {
        PlcSTRING plcSTRING = new PlcSTRING("not a boolean");
        assertTrue(plcSTRING.isBoolean()); // Boolean.parseBoolean always succeeds
    }

    @Test
    void testIsByteInvalid() {
        PlcSTRING plcSTRING = new PlcSTRING("not a number");
        assertFalse(plcSTRING.isByte());
    }

    @Test
    void testIsIntegerInvalid() {
        PlcSTRING plcSTRING = new PlcSTRING("not a number");
        assertFalse(plcSTRING.isInteger());
    }
}
