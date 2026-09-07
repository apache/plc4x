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

public class PlcDATE_AND_LTIMETest {

    @Test
    void testConstructorWithLocalDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(dateTime);
        assertEquals(dateTime, plcDATE_AND_LTIME.getDateTime());
        assertFalse(plcDATE_AND_LTIME.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(5000000000L);
        assertNotNull(plcDATE_AND_LTIME.getDateTime());
        assertFalse(plcDATE_AND_LTIME.isNullable());
    }

    @Test
    void testOfNanosecondsSinceEpoch() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = PlcDATE_AND_LTIME.ofNanosecondsSinceEpoch(BigInteger.valueOf(86400L * 1000_000_000));
        assertNotNull(plcDATE_AND_LTIME.getDateTime());
    }

    @Test
    void testOfSegments() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = PlcDATE_AND_LTIME.ofSegments(2025, 1, 15, 12, 30, 45, 500000000);
        assertEquals(2025, plcDATE_AND_LTIME.getDateTime().getYear());
        assertEquals(1, plcDATE_AND_LTIME.getDateTime().getMonthValue());
        assertEquals(15, plcDATE_AND_LTIME.getDateTime().getDayOfMonth());
    }

    @Test
    void testGetPlcValueType() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(LocalDateTime.now());
        assertEquals(PlcValueType.DATE_AND_LTIME, plcDATE_AND_LTIME.getPlcValueType());
    }

    @Test
    void testIsDateTime() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(LocalDateTime.now());
        assertTrue(plcDATE_AND_LTIME.isDateTime());
    }

    @Test
    void testGetDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(dateTime);
        assertEquals(dateTime, plcDATE_AND_LTIME.getDateTime());
    }

    @Test
    void testIsDate() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(LocalDateTime.now());
        assertTrue(plcDATE_AND_LTIME.isDate());
    }

    @Test
    void testGetDate() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(dateTime);
        assertEquals(LocalDate.of(2025, 1, 15), plcDATE_AND_LTIME.getDate());
    }

    @Test
    void testIsTime() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(LocalDateTime.now());
        assertTrue(plcDATE_AND_LTIME.isTime());
    }

    @Test
    void testGetTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(dateTime);
        assertEquals(LocalTime.of(12, 30, 45), plcDATE_AND_LTIME.getTime());
    }

    @Test
    void testGetNanosecondsSinceEpoch() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(5000000000L);
        assertNotNull(plcDATE_AND_LTIME.getNanosecondsSinceEpoch());
    }

    @Test
    void testGetLong() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(5000000000L);
        assertTrue(plcDATE_AND_LTIME.isLong());
        assertNotNull(plcDATE_AND_LTIME.getLong());
    }

    @Test
    void testIsString() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(LocalDateTime.now());
        assertTrue(plcDATE_AND_LTIME.isString());
    }

    @Test
    void testGetString() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = new PlcDATE_AND_LTIME(dateTime);
        assertNotNull(plcDATE_AND_LTIME.getString());
    }

    @Test
    void testOfWithPlcDATE_AND_LTIME() {
        PlcDATE_AND_LTIME original = new PlcDATE_AND_LTIME(LocalDateTime.now());
        PlcDATE_AND_LTIME copy = PlcDATE_AND_LTIME.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithLocalDateTime() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = PlcDATE_AND_LTIME.of(dateTime);
        assertEquals(dateTime, plcDATE_AND_LTIME.getDateTime());
    }

    @Test
    void testOfWithLong() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = PlcDATE_AND_LTIME.of(5000000000L);
        assertNotNull(plcDATE_AND_LTIME.getDateTime());
    }

    @Test
    void testOfWithBigInteger() {
        PlcDATE_AND_LTIME plcDATE_AND_LTIME = PlcDATE_AND_LTIME.of(BigInteger.valueOf(5000000000L));
        assertNotNull(plcDATE_AND_LTIME.getDateTime());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcDATE_AND_LTIME value = new PlcDATE_AND_LTIME((byte) 100);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcDATE_AND_LTIME value = new PlcDATE_AND_LTIME((short) 1000);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcDATE_AND_LTIME value = new PlcDATE_AND_LTIME(86400);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcDATE_AND_LTIME value = new PlcDATE_AND_LTIME(86400.0f);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcDATE_AND_LTIME value = new PlcDATE_AND_LTIME(86400.0);
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    @Test
    void testConstructorWithBigDecimal() {
        PlcDATE_AND_LTIME value = new PlcDATE_AND_LTIME(BigDecimal.valueOf(86400));
        assertNotNull(value.getDateTime());
        assertFalse(value.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcDATE_AND_LTIME value = PlcDATE_AND_LTIME.of((byte) 100);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithShort() {
        PlcDATE_AND_LTIME value = PlcDATE_AND_LTIME.of((short) 1000);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithInteger() {
        PlcDATE_AND_LTIME value = PlcDATE_AND_LTIME.of(86400);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithFloat() {
        PlcDATE_AND_LTIME value = PlcDATE_AND_LTIME.of(86400.0f);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithDouble() {
        PlcDATE_AND_LTIME value = PlcDATE_AND_LTIME.of(86400.0);
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcDATE_AND_LTIME value = PlcDATE_AND_LTIME.of(BigDecimal.valueOf(86400));
        assertNotNull(value.getDateTime());
    }

    @Test
    void testOfWithString() {
        PlcDATE_AND_LTIME value = PlcDATE_AND_LTIME.of("2025-01-15T12:30:45");
        assertEquals(LocalDateTime.of(2025, 1, 15, 12, 30, 45), value.getDateTime());
    }

    @Test
    void testToString() {
        LocalDateTime dateTime = LocalDateTime.of(2025, 1, 15, 12, 30, 45);
        PlcDATE_AND_LTIME value = new PlcDATE_AND_LTIME(dateTime);
        assertNotNull(value.toString());
    }
}
