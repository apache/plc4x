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
import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;

public class PlcTIMETest {

    @Test
    void testConstructorWithDuration() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofMillis(5000));
        assertEquals(5000, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcTIME plcTIME = new PlcTIME(3000L);
        assertEquals(3000, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcTIME plcTIME = new PlcTIME(2000);
        assertEquals(2000, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    @Test
    void testOfMilliseconds() {
        PlcTIME plcTIME = PlcTIME.ofMilliseconds(1500);
        assertEquals(1500, plcTIME.getMilliseconds());
    }

    @Test
    void testGetPlcValueType() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofMillis(1000));
        assertEquals(PlcValueType.TIME, plcTIME.getPlcValueType());
    }

    @Test
    void testIsDuration() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofMillis(1000));
        assertTrue(plcTIME.isDuration());
    }

    @Test
    void testGetDuration() {
        Duration duration = Duration.ofMillis(4000);
        PlcTIME plcTIME = new PlcTIME(duration);
        assertEquals(duration, plcTIME.getDuration());
    }

    @Test
    void testGetLong() {
        PlcTIME plcTIME = new PlcTIME(5000L);
        assertEquals(5000, plcTIME.getLong());
    }

    @Test
    void testGetInteger() {
        PlcTIME plcTIME = new PlcTIME(3000);
        assertEquals(3000, plcTIME.getInteger());
    }

    @Test
    void testIsString() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofMillis(1000));
        assertTrue(plcTIME.isString());
    }

    @Test
    void testGetString() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofSeconds(60));
        assertNotNull(plcTIME.getString());
    }

    @Test
    void testOfWithPlcTIME() {
        PlcTIME original = new PlcTIME(Duration.ofMillis(1000));
        PlcTIME copy = PlcTIME.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithDuration() {
        Duration duration = Duration.ofMillis(2000);
        PlcTIME plcTIME = PlcTIME.of(duration);
        assertEquals(2000, plcTIME.getMilliseconds());
    }

    @Test
    void testOfWithLong() {
        PlcTIME plcTIME = PlcTIME.of(3000L);
        assertEquals(3000, plcTIME.getMilliseconds());
    }

    @Test
    void testOfWithBigInteger() {
        PlcTIME plcTIME = PlcTIME.of(BigInteger.valueOf(4000));
        assertEquals(4000, plcTIME.getMilliseconds());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcTIME plcTIME = PlcTIME.of(BigDecimal.valueOf(5000));
        assertEquals(5000, plcTIME.getMilliseconds());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcTIME plcTIME = new PlcTIME((byte) 100);
        assertEquals(100, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcTIME plcTIME = new PlcTIME((short) 5000);
        assertEquals(5000, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcTIME plcTIME = new PlcTIME(3000.0f);
        assertEquals(3000, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcTIME plcTIME = new PlcTIME(4000.0);
        assertEquals(4000, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    @Test
    void testConstructorWithBigInteger() {
        PlcTIME plcTIME = new PlcTIME(BigInteger.valueOf(6000));
        assertEquals(6000, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    @Test
    void testConstructorWithBigDecimal() {
        PlcTIME plcTIME = new PlcTIME(BigDecimal.valueOf(7000));
        assertEquals(7000, plcTIME.getMilliseconds());
        assertFalse(plcTIME.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcTIME plcTIME = PlcTIME.of((byte) 100);
        assertEquals(100, plcTIME.getMilliseconds());
    }

    @Test
    void testOfWithShort() {
        PlcTIME plcTIME = PlcTIME.of((short) 5000);
        assertEquals(5000, plcTIME.getMilliseconds());
    }

    @Test
    void testOfWithInteger() {
        PlcTIME plcTIME = PlcTIME.of(3000);
        assertEquals(3000, plcTIME.getMilliseconds());
    }

    @Test
    void testOfWithFloat() {
        PlcTIME plcTIME = PlcTIME.of(4000.0f);
        assertEquals(4000, plcTIME.getMilliseconds());
    }

    @Test
    void testOfWithDouble() {
        PlcTIME plcTIME = PlcTIME.of(5000.0);
        assertEquals(5000, plcTIME.getMilliseconds());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsLong() {
        PlcTIME value = new PlcTIME(5000000000L);
        assertTrue(value.isLong());
    }

    // ========== Edge Case Tests ==========

    @Test
    void testZeroDuration() {
        PlcTIME plcTIME = new PlcTIME(Duration.ZERO);
        assertEquals(0, plcTIME.getMilliseconds());
    }

    @Test
    void testNegativeDuration() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofMillis(-1000));
        assertEquals(-1000, plcTIME.getMilliseconds());
    }

    @Test
    void testVeryLargeDuration() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofDays(365));
        assertEquals(365L * 24 * 60 * 60 * 1000, plcTIME.getMilliseconds());
    }

    @Test
    void testHoursMinutesSeconds() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofHours(1).plusMinutes(30).plusSeconds(45));
        assertEquals((1 * 60 * 60 + 30 * 60 + 45) * 1000, plcTIME.getMilliseconds());
    }

    @Test
    void testMillisecondPrecision() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofMillis(1234));
        assertEquals(1234, plcTIME.getMilliseconds());
    }

    // ========== Metadata Tests ==========

    @Test
    void testGetLength() {
        PlcTIME plcTIME = new PlcTIME(5000);
        assertEquals(1, plcTIME.getLength());
    }

    @Test
    void testGetIndex() {
        PlcTIME plcTIME = new PlcTIME(5000);
        assertEquals(plcTIME, plcTIME.getIndex(0));
    }

    @Test
    void testIsSimple() {
        PlcTIME plcTIME = new PlcTIME(5000);
        assertTrue(plcTIME.isSimple());
    }

    @Test
    void testGetObject() {
        PlcTIME plcTIME = new PlcTIME(5000);
        Object obj = plcTIME.getObject();
        assertNotNull(obj);
        assertTrue(obj instanceof Duration);
        assertEquals(Duration.ofMillis(5000), obj);
    }

    @Test
    void testToString() {
        PlcTIME plcTIME = new PlcTIME(Duration.ofSeconds(5));
        String str = plcTIME.toString();
        assertNotNull(str);
        assertTrue(str.contains("PT") || str.contains("5"));
    }
}
