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
import java.time.LocalDate;

import static org.junit.jupiter.api.Assertions.*;

public class PlcDATETest {

    @Test
    void testConstructorWithLocalDate() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcDATE plcDATE = new PlcDATE(date);
        assertEquals(date, plcDATE.getDate());
        assertFalse(plcDATE.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcDATE plcDATE = new PlcDATE(100);
        assertNotNull(plcDATE.getDate());
        assertFalse(plcDATE.isNullable());
    }

    @Test
    void testOfDaysSinceEpoch() {
        PlcDATE plcDATE = PlcDATE.ofDaysSinceEpoch(100);
        assertEquals(100, plcDATE.getDaysSinceEpoch());
    }

    @Test
    void testOfSecondsSinceEpoch() {
        PlcDATE plcDATE = PlcDATE.ofSecondsSinceEpoch(86400 * 100);
        assertEquals(100, plcDATE.getDaysSinceEpoch());
    }

    @Test
    void testGetPlcValueType() {
        PlcDATE plcDATE = new PlcDATE(LocalDate.now());
        assertEquals(PlcValueType.DATE, plcDATE.getPlcValueType());
    }

    @Test
    void testIsDate() {
        PlcDATE plcDATE = new PlcDATE(LocalDate.now());
        assertTrue(plcDATE.isDate());
    }

    @Test
    void testGetDate() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcDATE plcDATE = new PlcDATE(date);
        assertEquals(date, plcDATE.getDate());
    }

    @Test
    void testGetDateTime() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcDATE plcDATE = new PlcDATE(date);
        assertEquals(date.atStartOfDay(), plcDATE.getDateTime());
    }

    @Test
    void testGetLong() {
        PlcDATE plcDATE = new PlcDATE(100);
        assertTrue(plcDATE.isLong());
        assertNotNull(plcDATE.getLong());
    }

    @Test
    void testIsString() {
        PlcDATE plcDATE = new PlcDATE(LocalDate.now());
        assertTrue(plcDATE.isString());
    }

    @Test
    void testGetString() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcDATE plcDATE = new PlcDATE(date);
        assertEquals("2025-01-15", plcDATE.getString());
    }

    @Test
    void testOfWithPlcDATE() {
        PlcDATE original = new PlcDATE(LocalDate.now());
        PlcDATE copy = PlcDATE.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithLocalDate() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcDATE plcDATE = PlcDATE.of(date);
        assertEquals(date, plcDATE.getDate());
    }

    @Test
    void testOfWithInteger() {
        PlcDATE plcDATE = PlcDATE.of(100);
        assertEquals(100, plcDATE.getDaysSinceEpoch());
    }

    @Test
    void testOfWithBigInteger() {
        PlcDATE plcDATE = PlcDATE.of(BigInteger.valueOf(100));
        assertEquals(100, plcDATE.getDaysSinceEpoch());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcDATE plcDATE = PlcDATE.of(BigDecimal.valueOf(100));
        assertEquals(100, plcDATE.getDaysSinceEpoch());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcDATE value = new PlcDATE((byte) 10);
        assertNotNull(value.getDate());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcDATE value = new PlcDATE((short) 100);
        assertNotNull(value.getDate());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcDATE value = new PlcDATE(100L);
        assertNotNull(value.getDate());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcDATE value = new PlcDATE(100.0f);
        assertNotNull(value.getDate());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcDATE value = new PlcDATE(100.0);
        assertNotNull(value.getDate());
        assertFalse(value.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcDATE value = PlcDATE.of((byte) 10);
        assertNotNull(value.getDate());
    }

    @Test
    void testOfWithShort() {
        PlcDATE value = PlcDATE.of((short) 100);
        assertNotNull(value.getDate());
    }

    @Test
    void testOfWithLong() {
        PlcDATE value = PlcDATE.of(100L);
        assertNotNull(value.getDate());
    }

    @Test
    void testOfWithFloat() {
        PlcDATE value = PlcDATE.of(100.0f);
        assertNotNull(value.getDate());
    }

    @Test
    void testOfWithDouble() {
        PlcDATE value = PlcDATE.of(100.0);
        assertNotNull(value.getDate());
    }

    @Test
    void testOfWithString() {
        PlcDATE value = PlcDATE.of("2025-01-15");
        assertEquals(LocalDate.of(2025, 1, 15), value.getDate());
    }

    @Test
    void testToString() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcDATE value = new PlcDATE(date);
        assertNotNull(value.toString());
    }
}
