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

public class PlcLDATETest {

    @Test
    void testConstructorWithLocalDate() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcLDATE plcLDATE = new PlcLDATE(date);
        assertEquals(date, plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcLDATE plcLDATE = new PlcLDATE(86400L);
        assertNotNull(plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    @Test
    void testOfNanosecondsSinceEpoch() {
        PlcLDATE plcLDATE = PlcLDATE.ofNanosecondsSinceEpoch(BigInteger.valueOf(86400L * 1000_000_000));
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testGetPlcValueType() {
        PlcLDATE plcLDATE = new PlcLDATE(LocalDate.now());
        assertEquals(PlcValueType.LDATE, plcLDATE.getPlcValueType());
    }

    @Test
    void testIsDate() {
        PlcLDATE plcLDATE = new PlcLDATE(LocalDate.now());
        assertTrue(plcLDATE.isDate());
    }

    @Test
    void testGetDate() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcLDATE plcLDATE = new PlcLDATE(date);
        assertEquals(date, plcLDATE.getDate());
    }

    @Test
    void testGetDateTime() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcLDATE plcLDATE = new PlcLDATE(date);
        assertEquals(date.atStartOfDay(), plcLDATE.getDateTime());
    }

    @Test
    void testGetNanosecondsSinceEpoch() {
        PlcLDATE plcLDATE = new PlcLDATE(86400L);
        assertNotNull(plcLDATE.getNanosecondsSinceEpoch());
    }

    @Test
    void testIsString() {
        PlcLDATE plcLDATE = new PlcLDATE(LocalDate.now());
        assertTrue(plcLDATE.isString());
    }

    @Test
    void testGetString() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcLDATE plcLDATE = new PlcLDATE(date);
        assertEquals("2025-01-15", plcLDATE.getString());
    }

    @Test
    void testOfWithPlcLDATE() {
        PlcLDATE original = new PlcLDATE(LocalDate.now());
        PlcLDATE copy = PlcLDATE.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithLocalDate() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcLDATE plcLDATE = PlcLDATE.of(date);
        assertEquals(date, plcLDATE.getDate());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcLDATE plcLDATE = new PlcLDATE((byte) 100);
        assertNotNull(plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcLDATE plcLDATE = new PlcLDATE((short) 86400);
        assertNotNull(plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcLDATE plcLDATE = new PlcLDATE(86400);
        assertNotNull(plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcLDATE plcLDATE = new PlcLDATE(86400.0f);
        assertNotNull(plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcLDATE plcLDATE = new PlcLDATE(86400.0);
        assertNotNull(plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    @Test
    void testConstructorWithBigInteger() {
        PlcLDATE plcLDATE = new PlcLDATE(BigInteger.valueOf(86400));
        assertNotNull(plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    @Test
    void testConstructorWithBigDecimal() {
        PlcLDATE plcLDATE = new PlcLDATE(BigDecimal.valueOf(86400));
        assertNotNull(plcLDATE.getDate());
        assertFalse(plcLDATE.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcLDATE plcLDATE = PlcLDATE.of((byte) 100);
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testOfWithShort() {
        PlcLDATE plcLDATE = PlcLDATE.of((short) 86400);
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testOfWithInteger() {
        PlcLDATE plcLDATE = PlcLDATE.of(86400);
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testOfWithLong() {
        PlcLDATE plcLDATE = PlcLDATE.of(86400L);
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testOfWithFloat() {
        PlcLDATE plcLDATE = PlcLDATE.of(86400.0f);
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testOfWithDouble() {
        PlcLDATE plcLDATE = PlcLDATE.of(86400.0);
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testOfWithBigInteger() {
        PlcLDATE plcLDATE = PlcLDATE.of(BigInteger.valueOf(86400));
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcLDATE plcLDATE = PlcLDATE.of(BigDecimal.valueOf(86400));
        assertNotNull(plcLDATE.getDate());
    }

    @Test
    void testOfWithString() {
        PlcLDATE plcLDATE = PlcLDATE.of("2025-01-15");
        assertEquals(LocalDate.of(2025, 1, 15), plcLDATE.getDate());
    }

    // ========== Additional Getter Tests ==========

    @Test
    void testIsLong() {
        PlcLDATE plcLDATE = new PlcLDATE(LocalDate.now());
        assertTrue(plcLDATE.isLong());
    }

    @Test
    void testGetLong() {
        LocalDate date = LocalDate.of(1970, 1, 2);
        PlcLDATE plcLDATE = new PlcLDATE(date);
        assertEquals(86400L, plcLDATE.getLong());
    }

    @Test
    void testToString() {
        LocalDate date = LocalDate.of(2025, 1, 15);
        PlcLDATE plcLDATE = new PlcLDATE(date);
        assertEquals("2025-01-15", plcLDATE.toString());
    }

    @Test
    void testGetNanosecondsSinceEpochSpecificValue() {
        LocalDate date = LocalDate.of(1970, 1, 2);
        PlcLDATE plcLDATE = new PlcLDATE(date);
        BigInteger nanos = plcLDATE.getNanosecondsSinceEpoch();
        assertEquals(BigInteger.valueOf(86400L * 1000_000_000L), nanos);
    }

    @Test
    void testEpochDate() {
        LocalDate epoch = LocalDate.of(1970, 1, 1);
        PlcLDATE plcLDATE = new PlcLDATE(epoch);
        assertEquals(0L, plcLDATE.getLong());
        assertEquals(epoch, plcLDATE.getDate());
    }
}
