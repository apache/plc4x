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
package org.apache.plc4x.java.spi.buffers.bytebased.encoding;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import java.math.BigInteger;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class EncodingBCDTest extends BaseEncodingDefaultTest {

    private EncodingBCD encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingBCD();
    }

    @Test
    void getName() {
        assertEquals("BCD", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingBCD.optionEncodingBCD();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("BCD", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeByteHappyPath() {
        assertEncodeDecode((byte) 42, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 99, 8, encoding::decodeByte, encoding::encodeByte);
    }

    @Test
    void encodeAndDecodeShortHappyPath() {
        assertEncodeDecode((short) 1234, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 9999, 16, encoding::decodeShort, encoding::encodeShort);
    }

    @Test
    void encodeAndDecodeIntHappyPath() {
        assertEncodeDecode(12345678, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(99999999, 32, encoding::decodeInt, encoding::encodeInt);
    }

    @Test
    void encodeAndDecodeLongHappyPath() {
        assertEncodeDecode(1234567890123456L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(9999999999999999L, 64, encoding::decodeLong, encoding::encodeLong);
    }

    @Test
    void encodeAndDecodeBigIntegerHappyPath() {
        assertEncodeDecode(
            new BigInteger("12345678901234567890"),
            80, // 20 digits = 80 bits (4 bits per digit)
            encoding::decodeBigInteger,
            encoding::encodeBigInteger
        );
    }

    // Edge Cases
    @Test
    void encodeAndDecodeZeroValues() {
        assertEncodeDecode((byte) 0, 4, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((short) 0, 4, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(0, 4, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(0L, 4, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(BigInteger.ZERO, 4, encoding::decodeBigInteger, encoding::encodeBigInteger);
    }

    @Test
    void encodeAndDecodeSingleDigit() {
        assertEncodeDecode((byte) 5, 4, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((short) 7, 4, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(9, 4, encoding::decodeInt, encoding::encodeInt);
    }

    @Test
    void encodeAndDecodeMaxValues() {
        // Max 2-digit BCD (8 bits)
        assertEncodeDecode((byte) 99, 8, encoding::decodeByte, encoding::encodeByte);

        // Max 4-digit BCD (16 bits)
        assertEncodeDecode((short) 9999, 16, encoding::decodeShort, encoding::encodeShort);

        // Max 8-digit BCD (32 bits)
        assertEncodeDecode(99999999, 32, encoding::decodeInt, encoding::encodeInt);

        // Max 16-digit BCD (64 bits)
        assertEncodeDecode(9999999999999999L, 64, encoding::decodeLong, encoding::encodeLong);
    }

    // Error Cases
    @ParameterizedTest
    @ValueSource(ints = {1, 2, 3, 5, 6, 7})
    void encodeWithBitsNotMultipleOf4(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(numBits, (byte) 1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(numBits, (short) 1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(numBits, 1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(numBits, 1L));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(numBits, BigInteger.ONE));
    }

    @Test
    void encodeNegativeValues() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(8, (byte) -1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(16, (short) -1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(32, -1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(64, -1L));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(128, BigInteger.valueOf(-1)));
    }

    @Test
    void encodeValuesTooLargeForBits() {
        // 1 digit (4 bits) can only store 0-9
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(4, (byte) 10));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(4, (short) 10));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(4, 10));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(4, 10L));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(4, BigInteger.TEN));

        // 2 digits (8 bits) can only store 0-99
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(8, (byte) 100));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(8, (short) 100));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(8, 100));
    }

    @Test
    void decodeInvalidBCDDigits() {
        // 0xA through 0xF are not valid BCD digits
        byte[] encoded1 = new byte[]{(byte) 0xA1}; // First nibble is invalid
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, encoded1));

        byte[] encoded2 = new byte[]{(byte) 0x1A}; // Second nibble is invalid
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, encoded2));
    }

    @Test
    void decodeWithInsufficientBytes() {
        // Same guard and same message as EncodingUnsignedBinary, instead of the raw
        // ArrayIndexOutOfBoundsException the nibble unpacking would otherwise throw.
        byte[] encoded1 = new byte[0];
        IllegalArgumentException e1 =
            assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, encoded1));
        assertEquals("Expected at least 1 bytes", e1.getMessage());

        byte[] encoded2 = new byte[1]; // Need 2 bytes for 16 bits
        IllegalArgumentException e2 =
            assertThrows(IllegalArgumentException.class, () -> encoding.decodeShort(16, encoded2));
        assertEquals("Expected at least 2 bytes", e2.getMessage());

        byte[] encoded3 = new byte[3]; // Need 4 bytes for 32 bits
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(32, encoded3));

        byte[] encoded4 = new byte[1]; // Need 2 bytes for a 12 bit (odd digit count) field
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(12, encoded4));

        byte[] encoded5 = new byte[7];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeLong(64, encoded5));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeBigInteger(64, encoded5));
    }

    @Test
    void decodeWithNull() {
        assertThrows(NullPointerException.class, () -> encoding.decodeByte(8, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeShort(16, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeInt(32, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeLong(64, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeBigInteger(128, null));
    }

    // Unsupported Operations
    @Test
    void unsupportedOperations() {
        // TODO: Update this list of unsupported operations.
        //assertUnsupportedOperations(encoding);
    }

    // ---- Bit alignment (see the class javadoc of EncodingBCD) ----------------------------------
    // The byte-based buffers hand partial fields to the encodings RIGHT-aligned: readBits starts
    // filling the result at bit (8 - (numBits % 8)) % 8 and writeBits consumes it from that same
    // bit. BCD used to pack/unpack from the high nibble of byte 0 regardless, which silently
    // corrupted every field with an odd number of digits (the S7 DateAndTime 12 bit millisecond
    // and 4 bit day-of-week fields).

    @Test
    void encodeOddDigitCountIsRightAligned() {
        // 3 digits in a 12 bit field -> the leading nibble of byte 0 is padding.
        assertArrayEquals(new byte[]{0x01, 0x23}, encoding.encodeInt(12, 123));
        assertArrayEquals(new byte[]{0x00, 0x07}, encoding.encodeInt(12, 7));
        assertArrayEquals(new byte[]{0x09, (byte) 0x99}, encoding.encodeInt(12, 999));
        // 1 digit in a 4 bit field -> the digit is the LOW nibble.
        assertArrayEquals(new byte[]{0x04}, encoding.encodeInt(4, 4));
        assertArrayEquals(new byte[]{0x00}, encoding.encodeInt(4, 0));
        assertArrayEquals(new byte[]{0x09}, encoding.encodeInt(4, 9));
        // 5 digits in a 20 bit field.
        assertArrayEquals(new byte[]{0x01, 0x23, 0x45}, encoding.encodeInt(20, 12345));
        // Same for the other widths sharing the packing code.
        assertArrayEquals(new byte[]{0x01, 0x23}, encoding.encodeShort(12, (short) 123));
        assertArrayEquals(new byte[]{0x04}, encoding.encodeByte(4, (byte) 4));
        assertArrayEquals(new byte[]{0x01, 0x23}, encoding.encodeLong(12, 123L));
        assertArrayEquals(new byte[]{0x01, 0x23}, encoding.encodeBigInteger(12, BigInteger.valueOf(123)));
    }

    @Test
    void decodeOddDigitCountIsRightAligned() {
        assertEquals(123, encoding.decodeInt(12, new byte[]{0x01, 0x23}));
        assertEquals(7, encoding.decodeInt(12, new byte[]{0x00, 0x07}));
        assertEquals(999, encoding.decodeInt(12, new byte[]{0x09, (byte) 0x99}));
        assertEquals(4, encoding.decodeInt(4, new byte[]{0x04}));
        assertEquals(9, encoding.decodeInt(4, new byte[]{0x09}));
        assertEquals(12345, encoding.decodeInt(20, new byte[]{0x01, 0x23, 0x45}));
        assertEquals((short) 123, encoding.decodeShort(12, new byte[]{0x01, 0x23}));
        assertEquals((byte) 4, encoding.decodeByte(4, new byte[]{0x04}));
        assertEquals(123L, encoding.decodeLong(12, new byte[]{0x01, 0x23}));
        assertEquals(BigInteger.valueOf(123), encoding.decodeBigInteger(12, new byte[]{0x01, 0x23}));
    }

    @Test
    void evenDigitCountIsUnchanged() {
        // numBits % 8 == 0 means offset 0, i.e. exactly the previous behaviour - these must not move.
        assertArrayEquals(new byte[]{0x12}, encoding.encodeInt(8, 12));
        assertArrayEquals(new byte[]{0x12, 0x34}, encoding.encodeInt(16, 1234));
        assertArrayEquals(new byte[]{0x12, 0x34, 0x56, 0x78}, encoding.encodeInt(32, 12345678));
        assertEquals(12, encoding.decodeInt(8, new byte[]{0x12}));
        assertEquals(1234, encoding.decodeInt(16, new byte[]{0x12, 0x34}));
        assertEquals(12345678, encoding.decodeInt(32, new byte[]{0x12, 0x34, 0x56, 0x78}));
    }

    @Test
    void oddDigitCountRoundTripsEveryValue() {
        for (int value = 0; value <= 9; value++) {
            assertArrayEquals(new byte[]{(byte) value}, encoding.encodeInt(4, value));
            assertEquals(value, encoding.decodeInt(4, encoding.encodeInt(4, value)));
        }
        for (int value = 0; value <= 999; value++) {
            assertEquals(value, encoding.decodeInt(12, encoding.encodeInt(12, value)));
        }
    }

    @Test
    void decodeRejectsInvalidNibblesInOddDigitFields() {
        // The padding nibble is ignored, but every DIGIT nibble is still validated.
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(4, new byte[]{0x0A}));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(12, new byte[]{0x0A, 0x23}));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(12, new byte[]{0x01, (byte) 0xB3}));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(12, new byte[]{0x01, 0x2C}));
        // ... and the padding nibble genuinely is ignored (readBits leaves it zeroed, but a caller
        // handing us a dirty buffer must not be able to change the decoded value).
        assertEquals(4, encoding.decodeInt(4, new byte[]{(byte) 0xF4}));
        assertEquals(123, encoding.decodeInt(12, new byte[]{(byte) 0xF1, 0x23}));
    }

}
