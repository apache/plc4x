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
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class EncodingUTF8Test extends BaseEncodingDefaultTest {

    private EncodingUTF8 encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingUTF8();
    }

    @Test
    void getName() {
        assertEquals("UTF8", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingUTF8.optionEncodingUTF8();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("UTF8", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeStringHappyPath() {
        assertEncodeDecode("Hello World!", 96, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeStringWithUnicodeCharacters() {
        // Test various Unicode characters
        assertEncodeDecode("Hello 世界!", 104, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // Chinese characters
        assertEncodeDecode("Café", 40, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // Accented characters
        assertEncodeDecode("Ñoño", 48, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // Spanish characters
        assertEncodeDecode("Здравствуй", 160, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // Cyrillic
    }

    @Test
    void encodeAndDecodeStringWithEmojis() {
        // Test emoji characters (4-byte UTF-8 sequences)
        assertEncodeDecode("Hello 😀", 80, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("🌍🚀", 64, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeByteHappyPath() {
        assertEncodeDecode((byte) 42, 24, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte)); // "42" = 3 chars * 8 bits
        assertEncodeDecode((byte) -42, 32, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte)); // "-42" = 4 chars * 8 bits
    }

    @Test
    void encodeAndDecodeShortHappyPath() {
        assertEncodeDecode((short) 12345, 40, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort)); // "12345" = 5 chars * 8 bits
        assertEncodeDecode((short) -12345, 48, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort)); // "-12345" = 6 chars * 8 bits
    }

    @Test
    void encodeAndDecodeIntHappyPath() {
        assertEncodeDecode(1234567890, 80, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt)); // "1234567890" = 10 chars * 8 bits
        assertEncodeDecode(-1234567890, 88, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt)); // "-1234567890" = 11 chars * 8 bits
    }

    @Test
    void encodeAndDecodeLongHappyPath() {
        assertEncodeDecode(1234567890123456789L, 152, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong)); // "1234567890123456789" = 19 chars * 8 bits
        assertEncodeDecode(-1234567890123456789L, 160, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong)); // "-1234567890123456789" = 20 chars * 8 bits
    }

    @Test
    void encodeAndDecodeBigIntegerHappyPath() {
        assertEncodeDecode(
            new BigInteger("123456789012345678901234567890"),
            240, // "123456789012345678901234567890" = 30 chars * 8 bits
            wrapDecoder(encoding::decodeBigInteger),
            wrapEncoder(encoding::encodeBigInteger)
        );
    }

    @Test
    void encodeAndDecodeFloatHappyPath() {
        assertEncodeDecode(123.456f, 56, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f); // "123.456" = 7 chars * 8 bits
        assertEncodeDecode(-123.456f, 64, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f); // "-123.456" = 8 chars * 8 bits
    }

    @Test
    void encodeAndDecodeDoubleHappyPath() {
        assertEncodeDecode(123.456789, 80, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001); // "123.456789" = 10 chars * 8 bits
        assertEncodeDecode(-123.456789, 88, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001); // "-123.456789" = 11 chars * 8 bits
    }

    @Test
    void encodeAndDecodeBigDecimalHappyPath() {
        assertEncodeDecode(
            new BigDecimal("123456.789012345678901234567890"),
            256, // "123456.789012345678901234567890" = 32 chars * 8 bits
            wrapDecoder(encoding::decodeBigDecimal),
            wrapEncoder(encoding::encodeBigDecimal)
        );
    }

    // Edge Cases
    @Test
    void encodeAndDecodeEmptyString() {
        assertEncodeDecode("", 0, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeSingleCharacter() {
        assertEncodeDecode("X", 8, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeSingleUnicodeCharacter() {
        assertEncodeDecode("世", 24, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // 3-byte UTF-8 character
        assertEncodeDecode("é", 16, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // 2-byte UTF-8 character
    }

    @Test
    void encodeAndDecodeZeroValues() {
        assertEncodeDecode((byte) 0, 8, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte));
        assertEncodeDecode((short) 0, 8, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort));
        assertEncodeDecode(0, 8, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt));
        assertEncodeDecode(0L, 8, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong));
        assertEncodeDecode(BigInteger.ZERO, 8, wrapDecoder(encoding::decodeBigInteger), wrapEncoder(encoding::encodeBigInteger));
        assertEncodeDecode(0.0f, 24, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f);
        assertEncodeDecode(0.0, 24, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001);
        assertEncodeDecode(BigDecimal.ZERO, 8, wrapDecoder(encoding::decodeBigDecimal), wrapEncoder(encoding::encodeBigDecimal));
    }

    @Test
    void encodeAndDecodeSpecialCharacters() {
        assertEncodeDecode("!@#$%^&*()", 80, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\\n\\r\\t", 48, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode(" ", 8, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeControlCharacters() {
        assertEncodeDecode("\n", 8, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\r", 8, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\t", 8, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeStringWithVariableByteLength() {
        // UTF-8 characters can be 1-4 bytes long
        String mixed = "A世🌍"; // 1 byte + 3 bytes + 4 bytes = 8 bytes total
        assertEncodeDecode(mixed, 64, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    // Error Cases
    @Test
    void encodeStringWithInsufficientBits() {
        String input = "Hello";
        int numBits = 32; // Only 4 bytes for 5 characters
        assertThrows(BufferException.class, () -> encoding.encodeString(numBits, input));
    }

    @Test
    void decodeStringWithInsufficientBytes() {
        byte[] encoded = "Hello".getBytes();
        int numBits = 32; // Expecting 5 bytes, got 4
        assertThrows(BufferException.class, () -> encoding.decodeString(numBits, encoded));
    }

    @ParameterizedTest
    @ValueSource(ints = {7, 9, 15, 17})
    void encodeStringWithBitsNotMultipleOf8(int numBits) {
        String input = "A";
        assertThrows(BufferException.class, () -> encoding.encodeString(numBits, input));
    }

    @ParameterizedTest
    @ValueSource(ints = {7, 9, 15, 17})
    void decodeStringWithBitsNotMultipleOf8(int numBits) {
        byte[] encoded = new byte[]{0x41}; // "A"
        assertThrows(BufferException.class, () -> encoding.decodeString(numBits, encoded));
    }

    @Test
    void decodeStringWithInvalidNumber() {
        byte[] encoded = "abc".getBytes(); // Not a valid number
        assertThrows(BufferException.class, () -> encoding.decodeInt(24, encoded));
    }

    @Test
    void decodeStringWithOverflowNumber() {
        byte[] encoded = "999999999999999999999".getBytes(); // Too large for Long
        assertThrows(BufferException.class, () -> encoding.decodeLong(168, encoded));
    }

    @Test
    void decodeStringWithInvalidFloat() {
        byte[] encoded = "not.a.float".getBytes();
        assertThrows(BufferException.class, () -> encoding.decodeFloat(88, encoded));
    }

    @Test
    void decodeStringWithInvalidDouble() {
        byte[] encoded = "not.a.double".getBytes();
        assertThrows(BufferException.class, () -> encoding.decodeDouble(96, encoded));
    }

    @Test
    void decodeStringWithNull() {
        assertThrows(BufferException.class, () -> encoding.decodeString(8, null));
    }

    @Test
    void encodeStringWithNull() {
        assertThrows(BufferException.class, () -> encoding.encodeString(8, null));
    }

    @Test
    void testUTF8ByteLengthCalculation() throws BufferException {
        // Test that UTF-8 encoding correctly calculates byte lengths
        String ascii = "Hello"; // 5 bytes
        byte[] encoded = encoding.encodeString(40, ascii);
        assertEquals(5, encoded.length);

        String unicode = "世界"; // 6 bytes (3 bytes each)
        encoded = encoding.encodeString(48, unicode);
        assertEquals(6, encoded.length);

        String emoji = "😀"; // 4 bytes
        encoded = encoding.encodeString(32, emoji);
        assertEquals(4, encoded.length);
    }

}