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
import java.time.LocalTime;

import static org.junit.jupiter.api.Assertions.*;

public class PlcTIME_OF_DAYTest {

    @Test
    void testConstructorWithLocalTime() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(time);
        assertEquals(time, plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(3600L);
        assertNotNull(plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    @Test
    void testOfMillisecondsSinceMidnight() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.ofMillisecondsSinceMidnight(5000);
        assertEquals(5000, plcTIME_OF_DAY.getMillisecondsSinceMidnight());
    }

    @Test
    void testGetPlcValueType() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.now());
        assertEquals(PlcValueType.TIME_OF_DAY, plcTIME_OF_DAY.getPlcValueType());
    }

    @Test
    void testIsTime() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.now());
        assertTrue(plcTIME_OF_DAY.isTime());
    }

    @Test
    void testGetTime() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(time);
        assertEquals(time, plcTIME_OF_DAY.getTime());
    }

    @Test
    void testGetLong() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(3600);
        assertTrue(plcTIME_OF_DAY.isLong());
        assertNotNull(plcTIME_OF_DAY.getLong());
    }

    @Test
    void testIsString() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.now());
        assertTrue(plcTIME_OF_DAY.isString());
    }

    @Test
    void testGetString() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(time);
        assertNotNull(plcTIME_OF_DAY.getString());
    }

    @Test
    void testOfWithPlcTIME_OF_DAY() {
        PlcTIME_OF_DAY original = new PlcTIME_OF_DAY(LocalTime.now());
        PlcTIME_OF_DAY copy = PlcTIME_OF_DAY.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithLocalTime() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of(time);
        assertEquals(time, plcTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithLong() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of(3600L);
        assertNotNull(plcTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithBigInteger() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of(BigInteger.valueOf(3600));
        assertNotNull(plcTIME_OF_DAY.getTime());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY((byte) 100);
        assertNotNull(plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY((short) 5000);
        assertNotNull(plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(10000);
        assertNotNull(plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(3600.0f);
        assertNotNull(plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(7200.0);
        assertNotNull(plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithBigInteger() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(BigInteger.valueOf(10000));
        assertNotNull(plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithBigDecimal() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(BigDecimal.valueOf(10000));
        assertNotNull(plcTIME_OF_DAY.getTime());
        assertFalse(plcTIME_OF_DAY.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of((byte) 100);
        assertNotNull(plcTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithShort() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of((short) 5000);
        assertNotNull(plcTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithInteger() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of(10000);
        assertNotNull(plcTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithFloat() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of(3600.0f);
        assertNotNull(plcTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithDouble() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of(7200.0);
        assertNotNull(plcTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.of(BigDecimal.valueOf(10000));
        assertNotNull(plcTIME_OF_DAY.getTime());
    }


    // ========== Edge Case Tests ==========

    @Test
    void testMidnight() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.MIDNIGHT);
        assertEquals(0, plcTIME_OF_DAY.getMillisecondsSinceMidnight());
    }

    @Test
    void testNoon() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.NOON);
        assertEquals(12 * 60 * 60 * 1000, plcTIME_OF_DAY.getMillisecondsSinceMidnight());
    }

    @Test
    void testEndOfDay() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.of(23, 59, 59));
        assertTrue(plcTIME_OF_DAY.getMillisecondsSinceMidnight() > 0);
    }

    @Test
    void testHourMinuteSecond() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.of(10, 30, 45));
        assertEquals((10 * 60 * 60 + 30 * 60 + 45) * 1000, plcTIME_OF_DAY.getMillisecondsSinceMidnight());
    }

    @Test
    void testMillisecondsSinceMidnight() {
        int millis = 5000;
        PlcTIME_OF_DAY plcTIME_OF_DAY = PlcTIME_OF_DAY.ofMillisecondsSinceMidnight(millis);
        assertEquals(millis, plcTIME_OF_DAY.getMillisecondsSinceMidnight());
    }

    // ========== Metadata Tests ==========

    @Test
    void testGetLength() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.of(12, 30, 45));
        assertEquals(1, plcTIME_OF_DAY.getLength());
    }

    @Test
    void testGetIndex() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.of(12, 30, 45));
        assertEquals(plcTIME_OF_DAY, plcTIME_OF_DAY.getIndex(0));
    }

    @Test
    void testIsSimple() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.of(12, 30, 45));
        assertTrue(plcTIME_OF_DAY.isSimple());
    }

    @Test
    void testGetObject() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(time);
        Object obj = plcTIME_OF_DAY.getObject();
        assertNotNull(obj);
        assertTrue(obj instanceof LocalTime);
        assertEquals(time, obj);
    }

    @Test
    void testToString() {
        PlcTIME_OF_DAY plcTIME_OF_DAY = new PlcTIME_OF_DAY(LocalTime.of(12, 30, 45));
        String str = plcTIME_OF_DAY.toString();
        assertNotNull(str);
        assertTrue(str.contains("12") || str.contains("30"));
    }
}
