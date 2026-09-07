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

public class PlcDATE_AND_TIMETest {

    @Test
    void testConstructorWithLocalDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(dateTime);
        assertEquals(dateTime, plcDATE_AND_TIME.getDateTime());
        assertFalse(plcDATE_AND_TIME.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(86400L);
        assertNotNull(plcDATE_AND_TIME.getDateTime());
        assertFalse(plcDATE_AND_TIME.isNullable());
    }

    @Test
    void testConstructorWithSegments() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(2025, 1, 15, 12, 30, 45, 0);
        assertEquals(2025, plcDATE_AND_TIME.getYear());
        assertEquals(1, plcDATE_AND_TIME.getMonth());
        assertEquals(15, plcDATE_AND_TIME.getDay());
        assertEquals(12, plcDATE_AND_TIME.getHour());
        assertEquals(30, plcDATE_AND_TIME.getMinutes());
        assertEquals(45, plcDATE_AND_TIME.getSeconds());
    }

    @Test
    void testOfSecondsSinceEpoch() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = PlcDATE_AND_TIME.ofSecondsSinceEpoch(86400);
        assertNotNull(plcDATE_AND_TIME.getDateTime());
    }

    @Test
    void testOfSegments() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = PlcDATE_AND_TIME.ofSegments(2025, 1, 15, 12, 30, 45, 0);
        assertEquals(2025, plcDATE_AND_TIME.getYear());
        assertEquals(1, plcDATE_AND_TIME.getMonth());
        assertEquals(15, plcDATE_AND_TIME.getDay());
    }

    @Test
    void testGetPlcValueType() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(LocalDateTime.now());
        assertEquals(PlcValueType.DATE_AND_TIME, plcDATE_AND_TIME.getPlcValueType());
    }

    @Test
    void testIsDateTime() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(LocalDateTime.now());
        assertTrue(plcDATE_AND_TIME.isDateTime());
    }

    @Test
    void testGetDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(dateTime);
        assertEquals(dateTime, plcDATE_AND_TIME.getDateTime());
    }

    @Test
    void testIsDate() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(LocalDateTime.now());
        assertTrue(plcDATE_AND_TIME.isDate());
    }

    @Test
    void testGetDate() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(dateTime);
        assertEquals(LocalDate.of(2025, 1, 15), plcDATE_AND_TIME.getDate());
    }

    @Test
    void testIsTime() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(LocalDateTime.now());
        assertTrue(plcDATE_AND_TIME.isTime());
    }

    @Test
    void testGetTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(dateTime);
        assertEquals(LocalTime.of(12, 30, 45), plcDATE_AND_TIME.getTime());
    }

    @Test
    void testGetSecondsSinceEpoch() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(86400L);
        assertNotNull(plcDATE_AND_TIME.getSecondsSinceEpoch());
    }

    @Test
    void testGetLong() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(86400L);
        assertTrue(plcDATE_AND_TIME.isLong());
        assertNotNull(plcDATE_AND_TIME.getLong());
    }

    @Test
    void testIsString() {
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(LocalDateTime.now());
        assertTrue(plcDATE_AND_TIME.isString());
    }

    @Test
    void testGetString() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_TIME plcDATE_AND_TIME = new PlcDATE_AND_TIME(dateTime);
        assertNotNull(plcDATE_AND_TIME.getString());
    }

    @Test
    void testOfWithPlcDATE_AND_TIME() {
        PlcDATE_AND_TIME original = new PlcDATE_AND_TIME(LocalDateTime.now());
        PlcDATE_AND_TIME copy = PlcDATE_AND_TIME.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithLocalDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_TIME plcDATE_AND_TIME = PlcDATE_AND_TIME.of(dateTime);
        assertEquals(dateTime, plcDATE_AND_TIME.getDateTime());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcDATE_AND_TIME value = new PlcDATE_AND_TIME((byte) 100);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcDATE_AND_TIME value = new PlcDATE_AND_TIME((short) 86400);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcDATE_AND_TIME value = new PlcDATE_AND_TIME(86400);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcDATE_AND_TIME value = new PlcDATE_AND_TIME(86400.0f);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcDATE_AND_TIME value = new PlcDATE_AND_TIME(86400.0);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithBigInteger() {
        PlcDATE_AND_TIME value = new PlcDATE_AND_TIME(BigInteger.valueOf(86400));
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithBigDecimal() {
        PlcDATE_AND_TIME value = new PlcDATE_AND_TIME(BigDecimal.valueOf(86400));
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of((byte) 100);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithShort() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of((short) 86400);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithInteger() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of(86400);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithLong() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of(86400L);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithFloat() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of(86400.0f);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithDouble() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of(86400.0);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithBigInteger() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of(BigInteger.valueOf(86400));
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of(BigDecimal.valueOf(86400));
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithString() {
        PlcDATE_AND_TIME value = PlcDATE_AND_TIME.of("2025-01-15T12:30:45");
        assertEquals(LocalDateTime.of(2025, 1, 15, 12, 30, 45), value.getDateTime());
    }

    @Test
    void testToString() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_TIME value = new PlcDATE_AND_TIME(dateTime);
        assertNotNull(value.toString());
    }
}
