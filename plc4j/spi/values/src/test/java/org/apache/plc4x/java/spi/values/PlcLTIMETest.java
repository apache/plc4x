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

public class PlcLTIMETest {

    @Test
    void testConstructorWithDuration() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofNanos(5000000));
        assertEquals(5000000, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    @Test
    void testConstructorWithLong() {
        PlcLTIME plcLTIME = new PlcLTIME(3000000L);
        assertEquals(3000000, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    @Test
    void testConstructorWithInteger() {
        PlcLTIME plcLTIME = new PlcLTIME(2000000);
        assertEquals(2000000, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    @Test
    void testOfNanoseconds() {
        PlcLTIME plcLTIME = PlcLTIME.ofNanoseconds(1500000);
        assertEquals(1500000, plcLTIME.getNanoseconds());
    }

    @Test
    void testOfNanosecondsWithBigInteger() {
        PlcLTIME plcLTIME = PlcLTIME.ofNanoseconds(BigInteger.valueOf(2500000));
        assertEquals(2500000, plcLTIME.getNanoseconds());
    }

    @Test
    void testGetPlcValueType() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofNanos(1000000));
        assertEquals(PlcValueType.LTIME, plcLTIME.getPlcValueType());
    }

    @Test
    void testIsDuration() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofNanos(1000000));
        assertTrue(plcLTIME.isDuration());
    }

    @Test
    void testGetDuration() {
        Duration duration = Duration.ofNanos(4000000);
        PlcLTIME plcLTIME = new PlcLTIME(duration);
        assertEquals(duration, plcLTIME.getDuration());
    }

    @Test
    void testGetLong() {
        PlcLTIME plcLTIME = new PlcLTIME(5000000L);
        assertEquals(5000000, plcLTIME.getLong());
    }

    @Test
    void testGetBigInteger() {
        PlcLTIME plcLTIME = new PlcLTIME(6000000L);
        assertEquals(BigInteger.valueOf(6000000), plcLTIME.getBigInteger());
    }

    @Test
    void testIsString() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofNanos(1000000));
        assertTrue(plcLTIME.isString());
    }

    @Test
    void testGetString() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofSeconds(60));
        assertNotNull(plcLTIME.getString());
    }

    @Test
    void testOfWithPlcLTIME() {
        PlcLTIME original = new PlcLTIME(Duration.ofNanos(1000000));
        PlcLTIME copy = PlcLTIME.of(original);
        assertSame(original, copy);
    }

    @Test
    void testOfWithDuration() {
        Duration duration = Duration.ofNanos(2000000);
        PlcLTIME plcLTIME = PlcLTIME.of(duration);
        assertEquals(2000000, plcLTIME.getNanoseconds());
    }

    @Test
    void testOfWithLong() {
        PlcLTIME plcLTIME = PlcLTIME.of(3000000L);
        assertEquals(3000000, plcLTIME.getNanoseconds());
    }

    @Test
    void testOfWithBigInteger() {
        PlcLTIME plcLTIME = PlcLTIME.of(BigInteger.valueOf(4000000));
        assertEquals(4000000, plcLTIME.getNanoseconds());
    }

    @Test
    void testOfWithBigDecimal() {
        PlcLTIME plcLTIME = PlcLTIME.of(BigDecimal.valueOf(5000000));
        assertEquals(5000000, plcLTIME.getNanoseconds());
    }

    // ========== Additional Constructor Tests ==========

    @Test
    void testConstructorWithByte() {
        PlcLTIME plcLTIME = new PlcLTIME((byte) 100);
        assertEquals(100, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    @Test
    void testConstructorWithShort() {
        PlcLTIME plcLTIME = new PlcLTIME((short) 5000);
        assertEquals(5000, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    @Test
    void testConstructorWithFloat() {
        PlcLTIME plcLTIME = new PlcLTIME(3000000.0f);
        assertEquals(3000000, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    @Test
    void testConstructorWithDouble() {
        PlcLTIME plcLTIME = new PlcLTIME(4000000.0);
        assertEquals(4000000, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    @Test
    void testConstructorWithBigInteger() {
        PlcLTIME plcLTIME = new PlcLTIME(BigInteger.valueOf(6000000));
        assertEquals(6000000, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    @Test
    void testConstructorWithBigDecimal() {
        PlcLTIME plcLTIME = new PlcLTIME(BigDecimal.valueOf(7000000));
        assertEquals(7000000, plcLTIME.getNanoseconds());
        assertFalse(plcLTIME.isNullable());
    }

    // ========== Additional Of Method Tests ==========

    @Test
    void testOfWithByte() {
        PlcLTIME plcLTIME = PlcLTIME.of((byte) 100);
        assertEquals(100, plcLTIME.getNanoseconds());
    }

    @Test
    void testOfWithShort() {
        PlcLTIME plcLTIME = PlcLTIME.of((short) 5000);
        assertEquals(5000, plcLTIME.getNanoseconds());
    }

    @Test
    void testOfWithInteger() {
        PlcLTIME plcLTIME = PlcLTIME.of(3000000);
        assertEquals(3000000, plcLTIME.getNanoseconds());
    }

    @Test
    void testOfWithFloat() {
        PlcLTIME plcLTIME = PlcLTIME.of(4000000.0f);
        assertEquals(4000000, plcLTIME.getNanoseconds());
    }

    @Test
    void testOfWithDouble() {
        PlcLTIME plcLTIME = PlcLTIME.of(5000000.0);
        assertEquals(5000000, plcLTIME.getNanoseconds());
    }

    // ========== Is-Type Tests ==========

    @Test
    void testIsLong() {
        PlcLTIME value = new PlcLTIME(5000000000L);
        assertTrue(value.isLong());
    }

    // ========== Edge Case Tests ==========

    @Test
    void testZeroDuration() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ZERO);
        assertEquals(0, plcLTIME.getNanoseconds());
    }

    @Test
    void testNegativeDuration() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofNanos(-1000000));
        assertEquals(-1000000, plcLTIME.getNanoseconds());
    }

    @Test
    void testVeryLargeDuration() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofDays(365));
        assertEquals(365L * 24 * 60 * 60 * 1000000000L, plcLTIME.getNanoseconds());
    }

    @Test
    void testHoursMinutesSeconds() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofHours(1).plusMinutes(30).plusSeconds(45));
        assertEquals((1 * 60 * 60 + 30 * 60 + 45) * 1000000000L, plcLTIME.getNanoseconds());
    }

    @Test
    void testNanosecondPrecision() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofNanos(1234567));
        assertEquals(1234567, plcLTIME.getNanoseconds());
    }

    @Test
    void testMillisToNanos() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofMillis(1));
        assertEquals(1000000, plcLTIME.getNanoseconds());
    }

    // ========== Metadata Tests ==========

    @Test
    void testGetLength() {
        PlcLTIME plcLTIME = new PlcLTIME(5000000);
        assertEquals(1, plcLTIME.getLength());
    }

    @Test
    void testGetIndex() {
        PlcLTIME plcLTIME = new PlcLTIME(5000000);
        assertEquals(plcLTIME, plcLTIME.getIndex(0));
    }

    @Test
    void testIsSimple() {
        PlcLTIME plcLTIME = new PlcLTIME(5000000);
        assertTrue(plcLTIME.isSimple());
    }

    @Test
    void testGetObject() {
        PlcLTIME plcLTIME = new PlcLTIME(5000000);
        Object obj = plcLTIME.getObject();
        assertNotNull(obj);
        assertTrue(obj instanceof Duration);
        assertEquals(Duration.ofNanos(5000000), obj);
    }

    @Test
    void testToString() {
        PlcLTIME plcLTIME = new PlcLTIME(Duration.ofSeconds(5));
        String str = plcLTIME.toString();
        assertNotNull(str);
        assertTrue(str.contains("PT") || str.contains("5"));
    }
}
