/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.iec608705104;

import org.apache.plc4x.java.iec608705104.readwrite.SevenOctetBinaryTime;
import org.apache.plc4x.java.iec608705104.readwrite.ThreeOctetBinaryTime;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

/**
 * The station fills in these fields, so it can fill in a month of zero. LocalDateTime answers an
 * impossible date with an exception, and that exception used to travel out through the frame it
 * arrived in, taking every information object behind it in the same ASDU with it.
 */
class Iec60870TimeRangeTest {

    private static SevenOctetBinaryTime cp56(int ms, int minutes, int hour, int day, int month, int year) {
        return new SevenOctetBinaryTime(ms, false, false, (byte) minutes, false, (byte) hour,
            (byte) 1, (byte) day, (byte) month, (byte) year);
    }

    private static ThreeOctetBinaryTime cp24(int ms, int minutes) {
        return new ThreeOctetBinaryTime(ms, false, (byte) minutes);
    }

    @Test
    void aCp56TimeThatCannotBeATimeIsReportedAsNoTime() {
        // The month of zero the report named.
        assertNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(0, 0, 0, 1, 0, 24)));
        assertNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(0, 0, 0, 1, 13, 24)));
        // And the other fields, each of which the station picks just as freely.
        assertNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(0, 0, 0, 0, 6, 24)));
        assertNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(0, 0, 25, 1, 6, 24)));
        assertNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(0, 61, 0, 1, 6, 24)));
        assertNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(60_001, 0, 0, 1, 6, 24)));
    }

    @Test
    void aCp56TimeThatCouldBeATimeIsRead() {
        assertNotNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(1500, 30, 12, 15, 6, 24)));
        // The edges of every field.
        assertNotNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(0, 0, 0, 1, 1, 0)));
        assertNotNull(Iec60870Connection.convertCp56Time2aToCalendar(cp56(59_999, 59, 23, 31, 12, 99)));
    }

    @Test
    void aCp24TimeThatCannotBeATimeIsReportedAsNoTime() {
        assertNull(Iec60870Connection.convertCp24Time2aToCalendar(cp24(0, 60)));
        assertNull(Iec60870Connection.convertCp24Time2aToCalendar(cp24(60_000, 0)));
    }

    @Test
    void aCp24TimeThatCouldBeATimeIsRead() {
        assertNotNull(Iec60870Connection.convertCp24Time2aToCalendar(cp24(1500, 30)));
        assertNotNull(Iec60870Connection.convertCp24Time2aToCalendar(cp24(59_999, 59)));
    }

    @Test
    void decidingOnATimeNeverThrows() {
        for (int month = 0; month <= 13; month++) {
            int m = month;
            assertDoesNotThrow(() -> Iec60870Connection.convertCp56Time2aToCalendar(cp56(0, 0, 0, 1, m, 24)),
                "a month of " + m + " must be answered, not thrown about");
        }
    }
}
