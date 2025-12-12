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
import java.nio.charset.StandardCharsets;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class EncodingUTF16Test extends BaseEncodingDefaultTest {
    private EncodingUTF16 encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingUTF16();
    }

    @Test
    void getName() {
        assertEquals("UTF16", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingUTF16.optionEncodingUTF16();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("UTF16", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeStringHappyPath() {
        assertEncodeDecode("Hello World!", 416, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeStringWithUnicodeCharacters() {
        assertEncodeDecode("Hello 世界!", 176, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("Café", 160, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("Ñoño", 160, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("Здравствуй", 352, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeStringWithEmojis() {
        assertEncodeDecode("Hello 😀", 288, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("🌍🚀", 160, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeByteHappyPath() {
        assertEncodeDecode((byte) 42, 96, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte));
        assertEncodeDecode((byte) -42, 128, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte));
    }

    @Test
    void encodeAndDecodeShortHappyPath() {
        assertEncodeDecode((short) 12345, 192, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort));
        assertEncodeDecode((short) -12345, 224, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort));
    }

    @Test
    void encodeAndDecodeIntHappyPath() {
        assertEncodeDecode(1234567890, 352, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt));
        assertEncodeDecode(-1234567890, 384, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt));
    }

    @Test
    void encodeAndDecodeLongHappyPath() {
        assertEncodeDecode(1234567890123456789L, 640, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong));
        assertEncodeDecode(-1234567890123456789L, 672, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong));
    }

    @Test
    void encodeAndDecodeBigIntegerHappyPath() {
        assertEncodeDecode(
            new BigInteger("123456789012345678901234567890"),
            992,
            wrapDecoder(encoding::decodeBigInteger),
            wrapEncoder(encoding::encodeBigInteger)
        );
    }

    @Test
    void encodeAndDecodeFloatHappyPath() {
        assertEncodeDecode(123.456f, 256, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f);
        assertEncodeDecode(-123.456f, 288, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f);
    }

    @Test
    void encodeAndDecodeDoubleHappyPath() {
        assertEncodeDecode(123.456789, 352, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001);
        assertEncodeDecode(-123.456789, 384, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001);
    }

    @Test
    void encodeAndDecodeBigDecimalHappyPath() {
        assertEncodeDecode(
            new BigDecimal("123456.789012345678901234567890"),
            512,
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
        assertEncodeDecode("X", 32, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeSingleUnicodeCharacter() {
        assertEncodeDecode("世", 32, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("é", 32, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeZeroValues() {
        assertEncodeDecode((byte) 0, 32, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte));
        assertEncodeDecode((short) 0, 32, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort));
        assertEncodeDecode(0, 32, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt));
        assertEncodeDecode(0L, 32, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong));
        assertEncodeDecode(BigInteger.ZERO, 32, wrapDecoder(encoding::decodeBigInteger), wrapEncoder(encoding::encodeBigInteger));
        assertEncodeDecode(0.0f, 128, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f);
        assertEncodeDecode(0.0, 128, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001);
        assertEncodeDecode(BigDecimal.ZERO, 32, wrapDecoder(encoding::decodeBigDecimal), wrapEncoder(encoding::encodeBigDecimal));
    }

    @Test
    void encodeAndDecodeSpecialCharacters() {
        assertEncodeDecode("!@#$%^&*()", 352, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\\n\\r\\t", 224, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeControlCharacters() {
        assertEncodeDecode("\n", 32, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\r", 32, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\t", 32, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeStringWithSurrogates() {
        // Test surrogate pairs (characters outside BMP)
        assertEncodeDecode("𝄞", 64, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // Musical G clef
        assertEncodeDecode("🌍", 64, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // Earth globe
    }

    @Test
    void encodeAndDecodeStringWithVariableByteLength() {
        String mixed = "A世🌍"; // ASCII + CJK + Emoji
        assertEncodeDecode(mixed, 160, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    // Error Cases
    @Test
    void encodeStringWithInsufficientBits() {
        String input = "Hello";
        int numBits = 32; // Only 2 UTF-16 characters worth of bits for 5 chars
        assertThrows(BufferException.class, () -> encoding.encodeString(numBits, input));
    }

    @Test
    void decodeStringWithInsufficientBytes() {
        byte[] encoded = new byte[]{0, 65, 0}; // Incomplete UTF-16 sequence
        int numBits = 32;
        assertThrows(BufferException.class, () -> encoding.decodeString(numBits, encoded));
    }

    @ParameterizedTest
    @ValueSource(ints = {7, 9, 15, 17, 23, 25})
    void encodeStringWithBitsNotMultipleOf16(int numBits) {
        String input = "A";
        assertThrows(BufferException.class, () -> encoding.encodeString(numBits, input));
    }

    @ParameterizedTest
    @ValueSource(ints = {7, 9, 15, 17, 23, 25})
    void decodeStringWithBitsNotMultipleOf16(int numBits) {
        byte[] encoded = new byte[]{0, 65};
        assertThrows(BufferException.class, () -> encoding.decodeString(numBits, encoded));
    }

    @Test
    void decodeStringWithInvalidNumber() {
        byte[] encoded = "abc".getBytes(StandardCharsets.UTF_16);
        assertThrows(BufferException.class, () -> encoding.decodeInt(64, encoded));
    }

    @Test
    void decodeStringWithOverflowNumber() {
        byte[] encoded = "999999999999999999999".getBytes(StandardCharsets.UTF_16);
        assertThrows(BufferException.class, () -> encoding.decodeLong(352, encoded));
    }

    @Test
    void decodeStringWithInvalidFloat() {
        byte[] encoded = "not.a.float".getBytes(StandardCharsets.UTF_16);
        assertThrows(BufferException.class, () -> encoding.decodeFloat(192, encoded));
    }

    @Test
    void decodeStringWithInvalidDouble() {
        byte[] encoded = "not.a.double".getBytes(StandardCharsets.UTF_16);
        assertThrows(BufferException.class, () -> encoding.decodeDouble(208, encoded));
    }

    @Test
    void decodeStringWithNull() {
        assertThrows(BufferException.class, () -> encoding.decodeString(16, null));
    }

    @Test
    void encodeStringWithNull() {
        assertThrows(BufferException.class, () -> encoding.encodeString(16, null));
    }

    @Test
    void testByteOrderHandling() throws BufferException {
        // UTF-16 (without explicit endianness) uses platform default or BOM
        String input = "A";
        byte[] encoded = encoding.encodeString(64, input);
        assertEquals(8, encoded.length);

        // The exact byte order depends on the platform/implementation
        // but it should be consistent for encode/decode
        String decoded = encoding.decodeString(64, encoded);
        assertEquals(input, decoded);
    }

    @Test
    void testSurrogatePairHandling() throws BufferException {
        // Test that surrogate pairs are handled correctly
        String input = "𝄞🌍"; // Two characters requiring surrogate pairs
        byte[] encoded = encoding.encodeString(160, input);
        assertEquals(20, encoded.length); // 4 bytes per surrogate pair

        String decoded = encoding.decodeString(160, encoded);
        assertEquals(input, decoded);
    }

    @Test
    void testStringPadding() throws BufferException {
        // Test that strings are properly padded when encoded with more bits than needed
        String input = "Hi"; // 2 characters = 4 bytes in UTF-16
        byte[] encoded = encoding.encodeString(64, input); // 8 bytes allocated
        assertEquals(8, encoded.length);

        String decoded = encoding.decodeString(64, encoded);
        // The decoded string should contain the original plus null padding
        assertTrue(decoded.startsWith(input) || decoded.endsWith(input));
    }
}