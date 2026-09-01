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
package org.apache.plc4x.java.spi.buffers.bytebased;

import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingBCD;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;

import static org.junit.jupiter.api.Assertions.*;

/**
 * End-to-end regression tests for BCD fields whose width is not a multiple of 8.
 * <p>
 * {@code readBits} right-aligns a partial read and {@code writeBits} consumes the value from the same
 * right-aligned bit, but {@code EncodingBCD} used to pack and unpack its digits starting at the high
 * nibble of byte 0. The two conventions agree only when {@code numBits % 8 == 0}, so every field with
 * an odd number of BCD digits silently lost its most significant digit on read and wrote a shifted
 * value on serialize. The concrete victim was the S7 {@code DateAndTime} tail - a 12 bit millisecond
 * field followed by a 4 bit day-of-week field - which is reproduced verbatim below.
 */
class BcdBitAlignmentTest {

    /**
     * The S7 DateAndTime field layout for 2024-08-19 14:35:07.123 with a day-of-week nibble of 4.
     * Hand-built rather than captured - 2024-08-19 is a Monday, whose Siemens nibble is 2 - because
     * what matters here is only that the trailing nibble differs from the millisecond digits next to
     * it, so a misaligned read or write cannot accidentally produce the expected bytes.
     */
    private static final byte[] S7_DATE_AND_TIME =
        {0x24, 0x08, 0x19, 0x14, 0x35, 0x07, 0x12, 0x34};

    @Test
    void readS7DateAndTimeTail() throws Exception {
        ReadBufferByteBased buffer =
            new ReadBufferByteBased(S7_DATE_AND_TIME, EncodingBCD.optionEncodingBCD());

        assertEquals((short) 24, buffer.readUnsignedShort(8));  // year
        assertEquals((short) 8, buffer.readUnsignedShort(8));   // month
        assertEquals((short) 19, buffer.readUnsignedShort(8));  // day
        assertEquals((short) 14, buffer.readUnsignedShort(8));  // hour
        assertEquals((short) 35, buffer.readUnsignedShort(8));  // minutes
        assertEquals((short) 7, buffer.readUnsignedShort(8));   // seconds
        // The 12 bit millisecond field: 0x123 -> 123, NOT 12 (the pre-fix result, which dropped the
        // trailing digit because the leading padding nibble was read as the first digit).
        assertEquals((short) 123, buffer.readUnsignedShort(12));
        // The 4 bit day-of-week nibble: 4, NOT 0 (the pre-fix result read the high nibble).
        assertEquals((short) 4, buffer.readUnsignedShort(4));
        assertEquals(0, buffer.getRemainingBits());
    }

    @Test
    void writeS7DateAndTimeTail() throws Exception {
        byte[] data = new byte[8];
        WriteBufferByteBased buffer =
            new WriteBufferByteBased(data, EncodingBCD.optionEncodingBCD());

        buffer.writeUnsignedShort(8, (short) 24);
        buffer.writeUnsignedShort(8, (short) 8);
        buffer.writeUnsignedShort(8, (short) 19);
        buffer.writeUnsignedShort(8, (short) 14);
        buffer.writeUnsignedShort(8, (short) 35);
        buffer.writeUnsignedShort(8, (short) 7);
        buffer.writeUnsignedShort(12, (short) 123);
        buffer.writeUnsignedShort(4, (short) 4);

        // Pre-fix this produced ... 0x23 0x00: writeBits emitted only the low 12 bits of the
        // left-aligned {0x12, 0x30} the encoding handed it.
        assertArrayEquals(S7_DATE_AND_TIME, data);
    }

    @Test
    void twelveBitFieldRoundTripsEveryValue() throws Exception {
        for (short value = 0; value <= 999; value++) {
            byte[] data = new byte[2];
            WriteBufferByteBased writeBuffer =
                new WriteBufferByteBased(data, EncodingBCD.optionEncodingBCD());
            writeBuffer.writeUnsignedShort(12, value);

            ReadBufferByteBased readBuffer =
                new ReadBufferByteBased(data, EncodingBCD.optionEncodingBCD());
            assertEquals(value, readBuffer.readUnsignedShort(12), "12 bit BCD round trip of " + value);
        }
    }

    /**
     * Note the two distinct alignments at play: within the {@code byte[]} the encoding exchanges with
     * the buffer the field is RIGHT-aligned (a lone digit sits in the low nibble), but the buffer then
     * appends those bits to the STREAM at the current bit position - so a 4 bit field written at bit 0
     * occupies the high nibble of byte 0. Both have to be right for the round trip to work.
     */
    @Test
    void fourBitFieldRoundTripsEveryValue() throws Exception {
        for (short value = 0; value <= 9; value++) {
            byte[] data = new byte[1];
            WriteBufferByteBased writeBuffer =
                new WriteBufferByteBased(data, EncodingBCD.optionEncodingBCD());
            writeBuffer.writeUnsignedShort(4, value);
            assertArrayEquals(new byte[]{(byte) (value << 4)}, data,
                "a 4 bit field written at bit 0 occupies the high nibble of the stream");

            ReadBufferByteBased readBuffer =
                new ReadBufferByteBased(data, EncodingBCD.optionEncodingBCD());
            assertEquals(value, readBuffer.readUnsignedShort(4), "4 bit BCD round trip of " + value);
        }
    }

    /**
     * The same 4 bit field read/written at bit 4 - i.e. in the LOW nibble, which is where the S7
     * day-of-week nibble actually lives. This is the case the pre-fix code got wrong.
     */
    @Test
    void fourBitFieldAtOddNibbleOffsetRoundTripsEveryValue() throws Exception {
        for (short value = 0; value <= 9; value++) {
            byte[] data = new byte[1];
            WriteBufferByteBased writeBuffer =
                new WriteBufferByteBased(data, EncodingBCD.optionEncodingBCD());
            writeBuffer.writeUnsignedShort(4, (short) 0);
            writeBuffer.writeUnsignedShort(4, value);
            assertArrayEquals(new byte[]{(byte) value}, data,
                "a 4 bit field written at bit 4 occupies the low nibble");

            ReadBufferByteBased readBuffer =
                new ReadBufferByteBased(data, EncodingBCD.optionEncodingBCD());
            assertEquals((short) 0, readBuffer.readUnsignedShort(4));
            assertEquals(value, readBuffer.readUnsignedShort(4), "4 bit BCD round trip of " + value);
        }
    }

    /** Two consecutive 4 bit BCD nibbles must land in the high then the low nibble of one byte. */
    @Test
    void consecutiveNibblesPackInOrder() throws Exception {
        byte[] data = new byte[1];
        WriteBufferByteBased writeBuffer =
            new WriteBufferByteBased(data, EncodingBCD.optionEncodingBCD());
        writeBuffer.writeUnsignedShort(4, (short) 7);
        writeBuffer.writeUnsignedShort(4, (short) 2);
        assertArrayEquals(new byte[]{0x72}, data);

        ReadBufferByteBased readBuffer =
            new ReadBufferByteBased(data, EncodingBCD.optionEncodingBCD());
        assertEquals((short) 7, readBuffer.readUnsignedShort(4));
        assertEquals((short) 2, readBuffer.readUnsignedShort(4));
    }

    /** A 20 bit / 5 digit field exercises the offset across more than two bytes. */
    @Test
    void twentyBitFieldRoundTrips() throws Exception {
        byte[] data = new byte[3];
        WriteBufferByteBased writeBuffer =
            new WriteBufferByteBased(data, EncodingBCD.optionEncodingBCD());
        writeBuffer.writeUnsignedInt(20, 12345);
        // Written at bit 0, so the 20 bits occupy bits 0..19 of the stream (the trailing nibble of
        // the 3 byte array stays untouched).
        assertArrayEquals(new byte[]{0x12, 0x34, 0x50}, data);

        ReadBufferByteBased readBuffer =
            new ReadBufferByteBased(data, EncodingBCD.optionEncodingBCD());
        assertEquals(12345, readBuffer.readUnsignedInt(20));
    }

    /** The long and BigInteger code paths share the same packing, so pin them too. */
    @Test
    void oddDigitLongAndBigIntegerFieldsRoundTrip() throws Exception {
        byte[] data = new byte[5];
        WriteBufferByteBased writeBuffer =
            new WriteBufferByteBased(data, EncodingBCD.optionEncodingBCD());
        writeBuffer.writeUnsignedLong(36, 123456789L);
        assertArrayEquals(new byte[]{0x12, 0x34, 0x56, 0x78, (byte) 0x90}, data);

        ReadBufferByteBased readBuffer =
            new ReadBufferByteBased(data, EncodingBCD.optionEncodingBCD());
        assertEquals(123456789L, readBuffer.readUnsignedLong(36));

        byte[] bigData = new byte[5];
        WriteBufferByteBased bigWriteBuffer =
            new WriteBufferByteBased(bigData, EncodingBCD.optionEncodingBCD());
        bigWriteBuffer.writeUnsignedBigInteger(36, BigInteger.valueOf(987654321L));
        assertArrayEquals(new byte[]{(byte) 0x98, 0x76, 0x54, 0x32, 0x10}, bigData);

        ReadBufferByteBased bigReadBuffer =
            new ReadBufferByteBased(bigData, EncodingBCD.optionEncodingBCD());
        assertEquals(BigInteger.valueOf(987654321L), bigReadBuffer.readUnsignedBigInteger(36));
    }
}
