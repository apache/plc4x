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
        byte[] encoded1 = new byte[0];
        assertThrows(ArrayIndexOutOfBoundsException.class, () -> encoding.decodeByte(8, encoded1));

        byte[] encoded2 = new byte[1]; // Need 2 bytes for 16 bits
        assertThrows(ArrayIndexOutOfBoundsException.class, () -> encoding.decodeShort(16, encoded2));

        byte[] encoded3 = new byte[3]; // Need 4 bytes for 32 bits
        assertThrows(ArrayIndexOutOfBoundsException.class, () -> encoding.decodeInt(32, encoded3));
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
}