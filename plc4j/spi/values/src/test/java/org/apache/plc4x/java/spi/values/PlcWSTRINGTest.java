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

public class PlcWSTRINGTest {

    @Test
    void testConstructorWithString() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("Hello World");
        assertEquals("Hello World", plcWSTRING.getString());
        assertFalse(plcWSTRING.isNullable());
    }

    @Test
    void testGetPlcValueType() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("test");
        assertEquals(PlcValueType.WSTRING, plcWSTRING.getPlcValueType());
    }

    @Test
    void testIsString() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("test");
        assertTrue(plcWSTRING.isString());
    }

    @Test
    void testGetString() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("test string");
        assertEquals("test string", plcWSTRING.getString());
    }

    @Test
    void testIsBooleanValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("true");
        assertTrue(plcWSTRING.isBoolean());
    }

    @Test
    void testGetBoolean() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("true");
        assertTrue(plcWSTRING.getBoolean());

        PlcWSTRING plcWSTRING2 = new PlcWSTRING("false");
        assertFalse(plcWSTRING2.getBoolean());
    }

    @Test
    void testIsByteValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("123");
        assertTrue(plcWSTRING.isByte());
    }

    @Test
    void testGetByte() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("42");
        assertEquals(42, plcWSTRING.getByte());
    }

    @Test
    void testIsShortValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("12345");
        assertTrue(plcWSTRING.isShort());
    }

    @Test
    void testGetShort() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("1000");
        assertEquals(1000, plcWSTRING.getShort());
    }

    @Test
    void testIsIntegerValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("123456");
        assertTrue(plcWSTRING.isInteger());
    }

    @Test
    void testGetInteger() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("100000");
        assertEquals(100000, plcWSTRING.getInteger());
    }

    @Test
    void testIsLongValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("9223372036854775807");
        assertTrue(plcWSTRING.isLong());
    }

    @Test
    void testGetLong() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("9000000000");
        assertEquals(9000000000L, plcWSTRING.getLong());
    }

    @Test
    void testIsBigIntegerValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("12345678901234567890");
        assertTrue(plcWSTRING.isBigInteger());
    }

    @Test
    void testGetBigInteger() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("12345678901234567890");
        assertEquals(new BigInteger("12345678901234567890"), plcWSTRING.getBigInteger());
    }

    @Test
    void testIsFloatValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("123.45");
        assertTrue(plcWSTRING.isFloat());
    }

    @Test
    void testGetFloat() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("123.45");
        assertEquals(123.45f, plcWSTRING.getFloat(), 0.001);
    }

    @Test
    void testIsDoubleValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("123.456789");
        assertTrue(plcWSTRING.isDouble());
    }

    @Test
    void testGetDouble() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("123.456789");
        assertEquals(123.456789, plcWSTRING.getDouble(), 0.000001);
    }

    @Test
    void testIsBigDecimalValid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("123.456789012345");
        assertTrue(plcWSTRING.isBigDecimal());
    }

    @Test
    void testGetBigDecimal() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("123.456");
        assertEquals(new BigDecimal("123.456"), plcWSTRING.getBigDecimal());
    }

    @Test
    void testGetLength() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("Hello");
        assertEquals(5, plcWSTRING.getLength());
    }

    @Test
    void testToString() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("test");
        assertEquals("test", plcWSTRING.toString());
    }

    @Test
    void testOfWithPlcWSTRING() {
        PlcWSTRING original = new PlcWSTRING("test");
        PlcWSTRING copy = PlcWSTRING.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithString() {
        PlcWSTRING plcWSTRING = PlcWSTRING.of("test");
        assertEquals("test", plcWSTRING.getString());
    }

    @Test
    void testOfWithObject() {
        PlcWSTRING plcWSTRING = PlcWSTRING.of(123);
        assertEquals("123", plcWSTRING.getString());
    }

    @Test
    void testIsBooleanInvalid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("not a boolean");
        assertTrue(plcWSTRING.isBoolean()); // Boolean.parseBoolean always succeeds
    }

    @Test
    void testIsByteInvalid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("not a number");
        assertFalse(plcWSTRING.isByte());
    }

    @Test
    void testIsIntegerInvalid() {
        PlcWSTRING plcWSTRING = new PlcWSTRING("not a number");
        assertFalse(plcWSTRING.isInteger());
    }
}
