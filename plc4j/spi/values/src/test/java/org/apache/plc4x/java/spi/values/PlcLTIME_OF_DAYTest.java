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

public class PlcLTIME_OF_DAYTest {

    @Test
    void testConstructorWithLocalTime() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(time);
        assertEquals(time, plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(5000L);
        assertNotNull(plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    @Test
    void testOfNanosecondsSinceMidnight() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.ofNanosecondsSinceMidnight(BigInteger.valueOf(5000000000L));
        assertNotNull(plcLTIME_OF_DAY.getNanosecondsSinceMidnight());
    }

    @Test
    void testGetPlcValueType() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.now());
        assertEquals(PlcValueType.LTIME_OF_DAY, plcLTIME_OF_DAY.getPlcValueType());
    }

    @Test
    void testIsTime() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.now());
        assertTrue(plcLTIME_OF_DAY.isTime());
    }

    @Test
    void testGetTime() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(time);
        assertEquals(time, plcLTIME_OF_DAY.getTime());
    }

    @Test
    void testGetNanosecondsSinceMidnight() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(5000L);
        assertNotNull(plcLTIME_OF_DAY.getNanosecondsSinceMidnight());
    }

    @Test
    void testGetLong() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(5000L);
        assertTrue(plcLTIME_OF_DAY.isLong());
        assertNotNull(plcLTIME_OF_DAY.getLong());
    }

    @Test
    void testIsString() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.now());
        assertTrue(plcLTIME_OF_DAY.isString());
    }

    @Test
    void testGetString() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(time);
        assertNotNull(plcLTIME_OF_DAY.getString());
    }

    @Test
    void testOfWithPlcLTIME_OF_DAY() {
        PlcLTIME_OF_DAY original = new PlcLTIME_OF_DAY(LocalTime.now());
        PlcLTIME_OF_DAY copy = PlcLTIME_OF_DAY.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithLocalTime() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of(time);
        assertEquals(time, plcLTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithLong() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of(5000L);
        assertNotNull(plcLTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithBigInteger() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of(BigInteger.valueOf(5000));
        assertNotNull(plcLTIME_OF_DAY.getTime());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY((byte) 100);
        assertNotNull(plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY((short) 5000);
        assertNotNull(plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(1000000);
        assertNotNull(plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(3600.0f);
        assertNotNull(plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(7200.0);
        assertNotNull(plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithBigInteger() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(BigInteger.valueOf(1000000));
        assertNotNull(plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    @Test
    void testConstructorWithBigDecimal() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(BigDecimal.valueOf(1000000));
        assertNotNull(plcLTIME_OF_DAY.getTime());
        assertFalse(plcLTIME_OF_DAY.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of((byte) 100);
        assertNotNull(plcLTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithShort() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of((short) 5000);
        assertNotNull(plcLTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithInteger() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of(1000000);
        assertNotNull(plcLTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithFloat() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of(3600.0f);
        assertNotNull(plcLTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithDouble() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of(7200.0);
        assertNotNull(plcLTIME_OF_DAY.getTime());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.of(BigDecimal.valueOf(1000000));
        assertNotNull(plcLTIME_OF_DAY.getTime());
    }


    // ========== Edge Case Tests ==========

    @Test
    void testMidnight() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.MIDNIGHT);
        assertEquals(BigInteger.ZERO, plcLTIME_OF_DAY.getNanosecondsSinceMidnight());
    }

    @Test
    void testNoon() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.NOON);
        assertEquals(BigInteger.valueOf(12L * 60 * 60 * 1000000000L), plcLTIME_OF_DAY.getNanosecondsSinceMidnight());
    }

    @Test
    void testEndOfDay() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.of(23, 59, 59));
        assertTrue(plcLTIME_OF_DAY.getNanosecondsSinceMidnight().compareTo(BigInteger.ZERO) > 0);
    }

    @Test
    void testHourMinuteSecond() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.of(10, 30, 45));
        assertEquals(BigInteger.valueOf((10 * 60 * 60 + 30 * 60 + 45) * 1000000000L), plcLTIME_OF_DAY.getNanosecondsSinceMidnight());
    }

    @Test
    void testNanosecondsSinceMidnight() {
        long nanos = 5000000000L;
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = PlcLTIME_OF_DAY.ofNanosecondsSinceMidnight(BigInteger.valueOf(nanos));
        assertEquals(BigInteger.valueOf(nanos), plcLTIME_OF_DAY.getNanosecondsSinceMidnight());
    }

    @Test
    void testNanosecondPrecision() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.of(0, 0, 0, 123456789));
        assertEquals(BigInteger.valueOf(123456789), plcLTIME_OF_DAY.getNanosecondsSinceMidnight());
    }

    // ========== Metadata Tests ==========

    @Test
    void testGetLength() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.of(12, 30, 45));
        assertEquals(1, plcLTIME_OF_DAY.getLength());
    }

    @Test
    void testGetIndex() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.of(12, 30, 45));
        assertEquals(plcLTIME_OF_DAY, plcLTIME_OF_DAY.getIndex(0));
    }

    @Test
    void testIsSimple() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.of(12, 30, 45));
        assertTrue(plcLTIME_OF_DAY.isSimple());
    }

    @Test
    void testGetObject() {
        LocalTime time = LocalTime.of(12, 30, 45);
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(time);
        Object obj = plcLTIME_OF_DAY.getObject();
        assertNotNull(obj);
        assertTrue(obj instanceof LocalTime);
        assertEquals(time, obj);
    }

    @Test
    void testToString() {
        PlcLTIME_OF_DAY plcLTIME_OF_DAY = new PlcLTIME_OF_DAY(LocalTime.of(12, 30, 45));
        String str = plcLTIME_OF_DAY.toString();
        assertNotNull(str);
        assertTrue(str.contains("12") || str.contains("30"));
    }
}
