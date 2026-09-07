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
package org.apache.plc4x.java.s7.readwrite.utils;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.DataItem;
import org.apache.plc4x.java.s7.readwrite.DateAndTime;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.values.PlcDATE_AND_TIME;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import java.time.LocalDateTime;

import static org.junit.jupiter.api.Assertions.*;

/**
 * The S7 DATE_AND_TIME tail is the only place in the repo where BCD is used on fields whose width is
 * not a multiple of 8: a 12 bit millisecond field followed by a 4 bit day-of-week nibble. Those two
 * therefore exercise the partial-read/partial-write alignment of {@code EncodingBCD} and nothing else
 * in the S7 driver does, which is why the day-of-week rotation added in the mspec could compile and
 * pass the existing suite while still putting the wrong nibble on the wire.
 */
class S7DateAndTimeBcdTest {

    /**
     * Siemens numbers the nibble 1 == Sunday .. 7 == Saturday; {@code parseSiemensDayOfWeek} rotates
     * that into the ISO numbering (1 == Monday .. 7 == Sunday) PlcDATE_AND_TIME uses. The dates below
     * are the seven consecutive days of the week starting Sunday 2024-08-18.
     */
    @ParameterizedTest(name = "{0} -> wire nibble {1}, ISO {2}")
    @CsvSource({
        "2024-08-18, 1, 7", // Sunday
        "2024-08-19, 2, 1", // Monday
        "2024-08-20, 3, 2", // Tuesday
        "2024-08-21, 4, 3", // Wednesday
        "2024-08-22, 5, 4", // Thursday
        "2024-08-23, 6, 5", // Friday
        "2024-08-24, 7, 6"  // Saturday
    })
    void dayOfWeekNibbleRoundTrips(String date, int expectedWireNibble, int expectedIsoDayOfWeek) throws Exception {
        LocalDateTime dateTime = LocalDateTime.parse(date + "T00:00:00");
        assertEquals(expectedIsoDayOfWeek, dateTime.getDayOfWeek().getValue(),
            "the test data must agree with java.time about which day this is");

        // Serialize: the nibble is written into the high half of the (byte-aligned) first byte.
        byte[] data = new byte[1];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(data);
        StaticHelper.serializeSiemensDayOfWeek(writeBuffer, new PlcDATE_AND_TIME(dateTime));
        assertEquals((byte) (expectedWireNibble << 4), data[0],
            "day-of-week nibble on the wire for " + date);

        // Parse the very same byte back and land on the ISO numbering again.
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(data);
        assertEquals((short) expectedIsoDayOfWeek, StaticHelper.parseSiemensDayOfWeek(readBuffer),
            "ISO day of week parsed back for " + date);
    }

    @Test
    void parseDayOfWeekRejectsNibblesOutsideOneToSeven() {
        for (int nibble : new int[]{0, 8, 9}) {
            byte[] data = {(byte) (nibble << 4)};
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(data);
            RuntimeException e = assertThrows(RuntimeException.class,
                () -> StaticHelper.parseSiemensDayOfWeek(readBuffer));
            assertTrue(e.getMessage().contains("outside the range [1, 7]"),
                "unexpected message for nibble " + nibble + ": " + e.getMessage());
        }
        // 0xA..0xF are not valid BCD digits at all - EncodingBCD rejects them first.
        byte[] data = {(byte) 0xA0};
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(data);
        assertThrows(RuntimeException.class, () -> StaticHelper.parseSiemensDayOfWeek(readBuffer));
    }

    // ---- The year byte (00-89 == 2000-2089, 90-99 == 1990-1999) ----------------------------------

    /**
     * Both ends of the representable window plus the two years the old {@code year > 2000} test got
     * wrong: 2000 fell into the 1900 branch and asked EncodingBCD for a two digit encoding of 100
     * (an IllegalArgumentException, not a BufferException, so it escaped the method's catch block).
     */
    @ParameterizedTest(name = "{0} -> year byte BCD {1}")
    @CsvSource({
        "1990, 90",
        "1999, 99",
        "2000, 0",
        "2001, 1",
        "2024, 24",
        "2089, 89"
    })
    void yearByteRoundTrips(int year, int bcdDigits) throws Exception {
        byte[] data = new byte[1];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(data);
        StaticHelper.serializeSiemensYear(writeBuffer,
            new PlcDATE_AND_TIME(LocalDateTime.of(year, 8, 19, 14, 35, 7)));
        assertEquals((byte) (((bcdDigits / 10) << 4) | (bcdDigits % 10)), data[0],
            "year byte on the wire for " + year);

        assertEquals((short) year, StaticHelper.parseSiemensYear(new ReadBufferByteBased(data)),
            "year parsed back for " + year);
    }

    /** Outside [1990, 2089] the single byte cannot hold the year, so it must be rejected loudly. */
    @ParameterizedTest(name = "year {0} rejected")
    @CsvSource({"1900", "1989", "2090", "2100"})
    void serializeYearRejectsUnrepresentableYears(int year) {
        byte[] data = new byte[1];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(data);
        RuntimeException e = assertThrows(RuntimeException.class,
            () -> StaticHelper.serializeSiemensYear(writeBuffer,
                new PlcDATE_AND_TIME(LocalDateTime.of(year, 8, 19, 14, 35, 7))));
        assertTrue(e.getMessage().contains("outside the range [1990, 2089]"),
            "unexpected message for year " + year + ": " + e.getMessage());
    }

    /**
     * The full DataItem path for year 2000 - the case that used to throw out of
     * {@code serializeSiemensYear} instead of producing a byte.
     */
    @Test
    void dataItemDateAndTimeRoundTripsTheYearTwoThousand() throws Exception {
        // 2000-08-19 is a Saturday -> Siemens nibble 7.
        LocalDateTime dateTime = LocalDateTime.of(2000, 8, 19, 14, 35, 7, 123_000_000);
        byte[] data = new byte[8];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(data);
        DataItem.staticSerialize(writeBuffer, new PlcDATE_AND_TIME(dateTime),
            "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);

        assertArrayEquals(new byte[]{0x00, 0x08, 0x19, 0x14, 0x35, 0x07, 0x12, 0x37}, data);

        PlcValue reparsed = DataItem.staticParse(new ReadBufferByteBased(data),
            "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);
        assertEquals(dateTime, reparsed.getDateTime());
    }

    // ---- The raw DateAndTime mspec type (12 bit msec + 4 bit dow as plain BCD fields) -------------

    /**
     * 2024-08-19 14:35:07.123 with a day-of-week nibble of 4. Hand-built, not a capture: 4 would mean
     * Wednesday while 2024-08-19 is a Monday (nibble 2). The mismatch is deliberate and harmless here
     * because the raw {@code DateAndTime} mspec type carries {@code dow} as an independent BCD field
     * and never cross-checks it against the date - which is exactly what lets these two tests pin the
     * 12 bit + 4 bit alignment with a nibble value distinct from the millisecond digits around it. The
     * self-consistent vector used for the value-level {@code DataItem} path is {@link #DATA_ITEM_WIRE}.
     */
    private static final byte[] DATE_AND_TIME_WIRE =
        {0x24, 0x08, 0x19, 0x14, 0x35, 0x07, 0x12, 0x34};

    @Test
    void dateAndTimeTypeParsesTheTrailingNibbles() throws Exception {
        DateAndTime parsed = DateAndTime.staticParse(new ReadBufferByteBased(DATE_AND_TIME_WIRE));

        assertEquals((short) 24, parsed.getYear());
        assertEquals((short) 8, parsed.getMonth());
        assertEquals((short) 19, parsed.getDay());
        assertEquals((short) 14, parsed.getHour());
        assertEquals((short) 35, parsed.getMinutes());
        assertEquals((short) 7, parsed.getSeconds());
        // Pre-fix this returned 12 (leading padding nibble mistaken for the first digit).
        assertEquals((short) 123, parsed.getMsec());
        // Pre-fix this returned 0 (the high nibble was read instead of the low one).
        assertEquals((byte) 4, parsed.getDow());
    }

    @Test
    void dateAndTimeTypeReSerializesTheTrailingNibbles() throws Exception {
        DateAndTime value = new DateAndTime((short) 24, (short) 8, (short) 19, (short) 14,
            (short) 35, (short) 7, (short) 123, (byte) 4);

        byte[] data = new byte[8];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(data);
        value.serialize(writeBuffer);

        // Pre-fix the last two bytes came out as 0x23 0x00.
        assertArrayEquals(DATE_AND_TIME_WIRE, data);
    }

    // ---- The DataItem DATE_AND_TIME data item (routes dayOfWeek through the static helpers) -------

    /**
     * 2024-08-19 14:35:07.123 with the day-of-week nibble that actually matches that Monday (2), so
     * the item is self-consistent and the parse -> serialize round trip must be byte-identical.
     * {@code DataItem} recomputes the nibble from the timestamp on serialize, which is exactly the
     * path that goes through {@code serializeSiemensDayOfWeek}.
     */
    private static final byte[] DATA_ITEM_WIRE =
        {0x24, 0x08, 0x19, 0x14, 0x35, 0x07, 0x12, 0x32};

    @Test
    void dataItemDateAndTimeParses() throws Exception {
        PlcValue value = DataItem.staticParse(new ReadBufferByteBased(DATA_ITEM_WIRE),
            "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);

        assertInstanceOf(PlcDATE_AND_TIME.class, value);
        LocalDateTime dateTime = value.getDateTime();
        assertEquals(LocalDateTime.of(2024, 8, 19, 14, 35, 7, 123_000_000), dateTime);
        // The nibble on the wire (2 == Monday) has to agree with the reconstructed timestamp.
        assertEquals(1, dateTime.getDayOfWeek().getValue());
    }

    @Test
    void dataItemDateAndTimeRoundTrips() throws Exception {
        PlcValue value = DataItem.staticParse(new ReadBufferByteBased(DATA_ITEM_WIRE),
            "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);

        byte[] data = new byte[8];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(data);
        DataItem.staticSerialize(writeBuffer, value, "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);

        assertArrayEquals(DATA_ITEM_WIRE, data);
    }

    /** Every day of the week must survive a full DataItem round trip, not just the Monday above. */
    @ParameterizedTest(name = "{0} -> wire nibble {1}")
    @CsvSource({
        "2024-08-18, 1",
        "2024-08-19, 2",
        "2024-08-20, 3",
        "2024-08-21, 4",
        "2024-08-22, 5",
        "2024-08-23, 6",
        "2024-08-24, 7"
    })
    void dataItemDateAndTimeRoundTripsEveryDayOfWeek(String date, int expectedWireNibble) throws Exception {
        LocalDateTime dateTime = LocalDateTime.parse(date + "T14:35:07").withNano(123_000_000);
        PlcValue value = new PlcDATE_AND_TIME(dateTime);

        byte[] data = new byte[8];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(data);
        DataItem.staticSerialize(writeBuffer, value, "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);

        // Trailing byte: high nibble = the last millisecond digit (3), low nibble = the day of week.
        assertEquals((byte) (0x30 | expectedWireNibble), data[7],
            "millisecond digit + day-of-week nibble for " + date);
        assertEquals((byte) 0x12, data[6], "leading millisecond digits for " + date);

        PlcValue reparsed = DataItem.staticParse(new ReadBufferByteBased(data),
            "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);
        assertEquals(dateTime, reparsed.getDateTime());
    }

    /** Guards the millisecond field specifically, across all four digit-count/alignment shapes. */
    @ParameterizedTest(name = "msec {0}")
    @CsvSource({"0", "1", "7", "12", "99", "100", "123", "456", "789", "999"})
    void dataItemDateAndTimeMillisecondsRoundTrip(int millis) throws Exception {
        LocalDateTime dateTime = LocalDateTime.of(2024, 8, 19, 14, 35, 7, millis * 1_000_000);
        PlcValue value = new PlcDATE_AND_TIME(dateTime);

        byte[] data = new byte[8];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(data);
        DataItem.staticSerialize(writeBuffer, value, "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);

        PlcValue reparsed = DataItem.staticParse(new ReadBufferByteBased(data),
            "IEC61131_DATE_AND_TIME", ControllerType.S7_300, 0);
        assertEquals(millis, reparsed.getDateTime().getNano() / 1_000_000);
        assertEquals(dateTime, reparsed.getDateTime());
    }
}
