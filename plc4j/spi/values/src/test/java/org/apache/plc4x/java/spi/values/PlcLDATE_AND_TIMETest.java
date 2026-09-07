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
import java.time.LocalDateTime;
import java.time.LocalTime;

import static org.junit.jupiter.api.Assertions.*;

public class PlcLDATE_AND_TIMETest {

    @Test
    void testConstructorWithLocalDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(dateTime);
        assertEquals(dateTime, plcLDATE_AND_TIME.getDateTime());
        assertFalse(plcLDATE_AND_TIME.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(86400000L);
        assertNotNull(plcLDATE_AND_TIME.getDateTime());
        assertFalse(plcLDATE_AND_TIME.isNullable());
    }

    @Test
    void testOfNanosecondsSinceEpoch() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = PlcLDATE_AND_TIME.ofNanosecondsSinceEpoch(BigInteger.valueOf(86400L * 1000_000_000));
        assertNotNull(plcLDATE_AND_TIME.getDateTime());
    }

    @Test
    void testGetPlcValueType() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(LocalDateTime.now());
        assertEquals(PlcValueType.DATE_AND_TIME, plcLDATE_AND_TIME.getPlcValueType());
    }

    @Test
    void testIsDateTime() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(LocalDateTime.now());
        assertTrue(plcLDATE_AND_TIME.isDateTime());
    }

    @Test
    void testGetDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(dateTime);
        assertEquals(dateTime, plcLDATE_AND_TIME.getDateTime());
    }

    @Test
    void testIsDate() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(LocalDateTime.now());
        assertTrue(plcLDATE_AND_TIME.isDate());
    }

    @Test
    void testGetDate() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(dateTime);
        assertEquals(LocalDate.of(2025, 1, 15), plcLDATE_AND_TIME.getDate());
    }

    @Test
    void testIsTime() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(LocalDateTime.now());
        assertTrue(plcLDATE_AND_TIME.isTime());
    }

    @Test
    void testGetTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(dateTime);
        assertEquals(LocalTime.of(12, 30, 45), plcLDATE_AND_TIME.getTime());
    }

    @Test
    void testGetNanosecondsSinceEpoch() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(86400000L);
        assertNotNull(plcLDATE_AND_TIME.getNanosecondsSinceEpoch());
    }

    @Test
    void testGetLong() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(86400000L);
        assertTrue(plcLDATE_AND_TIME.isLong());
        assertNotNull(plcLDATE_AND_TIME.getLong());
    }

    @Test
    void testIsString() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(LocalDateTime.now());
        assertTrue(plcLDATE_AND_TIME.isString());
    }

    @Test
    void testGetString() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = new PlcLDATE_AND_TIME(dateTime);
        assertNotNull(plcLDATE_AND_TIME.getString());
    }

    @Test
    void testOfWithPlcLDATE_AND_TIME() {
        PlcLDATE_AND_TIME original = new PlcLDATE_AND_TIME(LocalDateTime.now());
        PlcLDATE_AND_TIME copy = PlcLDATE_AND_TIME.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithLocalDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = PlcLDATE_AND_TIME.of(dateTime);
        assertEquals(dateTime, plcLDATE_AND_TIME.getDateTime());
    }

    @Test
    void testOfWithLong() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = PlcLDATE_AND_TIME.of(86400000L);
        assertNotNull(plcLDATE_AND_TIME.getDateTime());
    }

    @Test
    void testOfWithBigInteger() {
        PlcLDATE_AND_TIME plcLDATE_AND_TIME = PlcLDATE_AND_TIME.of(BigInteger.valueOf(86400000));
        assertNotNull(plcLDATE_AND_TIME.getDateTime());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcLDATE_AND_TIME value = new PlcLDATE_AND_TIME((byte) 100);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcLDATE_AND_TIME value = new PlcLDATE_AND_TIME((short) 1000);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcLDATE_AND_TIME value = new PlcLDATE_AND_TIME(86400);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcLDATE_AND_TIME value = new PlcLDATE_AND_TIME(86400.0f);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcLDATE_AND_TIME value = new PlcLDATE_AND_TIME(86400.0);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithBigDecimal() {
        PlcLDATE_AND_TIME value = new PlcLDATE_AND_TIME(BigDecimal.valueOf(86400));
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcLDATE_AND_TIME value = PlcLDATE_AND_TIME.of((byte) 100);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithShort() {
        PlcLDATE_AND_TIME value = PlcLDATE_AND_TIME.of((short) 1000);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithInteger() {
        PlcLDATE_AND_TIME value = PlcLDATE_AND_TIME.of(86400);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithFloat() {
        PlcLDATE_AND_TIME value = PlcLDATE_AND_TIME.of(86400.0f);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithDouble() {
        PlcLDATE_AND_TIME value = PlcLDATE_AND_TIME.of(86400.0);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcLDATE_AND_TIME value = PlcLDATE_AND_TIME.of(BigDecimal.valueOf(86400));
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithString() {
        PlcLDATE_AND_TIME value = PlcLDATE_AND_TIME.of("2025-01-15T12:30:45");
        assertEquals(LocalDateTime.of(2025, 1, 15, 12, 30, 45), value.getDateTime());
    }

    @Test
    void testToString() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcLDATE_AND_TIME value = new PlcLDATE_AND_TIME(dateTime);
        assertNotNull(value.toString());
    }
}
