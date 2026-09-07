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

class EncodingUTF16BETest extends BaseEncodingDefaultTest {

    private EncodingUTF16BE encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingUTF16BE();
    }

    @Test
    void getName() {
        assertEquals("UTF16BE", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingUTF16BE.optionEncodingUTF16BE();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("UTF16BE", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeStringHappyPath() {
        assertEncodeDecode("Hello World!", 192, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeStringWithUnicodeCharacters() {
        assertEncodeDecode("Hello 世界!", 176, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("Café", 64, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("Ñoño", 64, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("Здравствуй", 160, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeStringWithEmojis() {
        assertEncodeDecode("Hello 😀", 256, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("🌍🚀", 64, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeByteHappyPath() {
        assertEncodeDecode((byte) 42, 32, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte));
        assertEncodeDecode((byte) -42, 48, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte));
    }

    @Test
    void encodeAndDecodeShortHappyPath() {
        assertEncodeDecode((short) 12345, 80, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort));
        assertEncodeDecode((short) -12345, 96, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort));
    }

    @Test
    void encodeAndDecodeIntHappyPath() {
        assertEncodeDecode(1234567890, 160, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt));
        assertEncodeDecode(-1234567890, 176, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt));
    }

    @Test
    void encodeAndDecodeLongHappyPath() {
        assertEncodeDecode(1234567890123456789L, 304, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong));
        assertEncodeDecode(-1234567890123456789L, 320, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong));
    }

    @Test
    void encodeAndDecodeBigIntegerHappyPath() {
        assertEncodeDecode(
            new BigInteger("123456789012345678901234567890"),
            480,
            wrapDecoder(encoding::decodeBigInteger),
            wrapEncoder(encoding::encodeBigInteger)
        );
    }

    @Test
    void encodeAndDecodeFloatHappyPath() {
        assertEncodeDecode(123.456f, 112, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f);
        assertEncodeDecode(-123.456f, 128, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f);
    }

    @Test
    void encodeAndDecodeDoubleHappyPath() {
        assertEncodeDecode(123.456789, 160, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001);
        assertEncodeDecode(-123.456789, 176, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001);
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
        assertEncodeDecode("X", 16, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeSingleUnicodeCharacter() {
        assertEncodeDecode("世", 16, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("é", 16, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeZeroValues() {
        assertEncodeDecode((byte) 0, 16, wrapDecoder(encoding::decodeByte), wrapEncoder(encoding::encodeByte));
        assertEncodeDecode((short) 0, 16, wrapDecoder(encoding::decodeShort), wrapEncoder(encoding::encodeShort));
        assertEncodeDecode(0, 16, wrapDecoder(encoding::decodeInt), wrapEncoder(encoding::encodeInt));
        assertEncodeDecode(0L, 16, wrapDecoder(encoding::decodeLong), wrapEncoder(encoding::encodeLong));
        assertEncodeDecode(BigInteger.ZERO, 16, wrapDecoder(encoding::decodeBigInteger), wrapEncoder(encoding::encodeBigInteger));
        assertEncodeDecode(0.0f, 96, wrapDecoder(encoding::decodeFloat), wrapEncoder(encoding::encodeFloat), 0.0001f);
        assertEncodeDecode(0.0, 96, wrapDecoder(encoding::decodeDouble), wrapEncoder(encoding::encodeDouble), 0.0000001);
        assertEncodeDecode(BigDecimal.ZERO, 16, wrapDecoder(encoding::decodeBigDecimal), wrapEncoder(encoding::encodeBigDecimal));
    }

    @Test
    void encodeAndDecodeSpecialCharacters() {
        assertEncodeDecode("!@#$%^&*()", 160, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\\n\\r\\t", 96, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode(" ", 16, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeControlCharacters() {
        assertEncodeDecode("\n", 16, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\r", 16, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
        assertEncodeDecode("\t", 16, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString));
    }

    @Test
    void encodeAndDecodeStringWithSurrogates() {
        // Test surrogate pairs (characters outside BMP)
        assertEncodeDecode("𝄞", 32, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // Musical G clef
        assertEncodeDecode("🌍", 32, wrapDecoder(encoding::decodeString), wrapEncoder(encoding::encodeString)); // Earth globe
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
        byte[] encoded = new byte[]{0, 65, 0}; // Incomplete UTF-16BE sequence
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
        byte[] encoded = "abc".getBytes(StandardCharsets.UTF_16BE);
        assertThrows(BufferException.class, () -> encoding.decodeInt(48, encoded));
    }

    @Test
    void decodeStringWithOverflowNumber() {
        byte[] encoded = "999999999999999999999".getBytes(StandardCharsets.UTF_16BE);
        assertThrows(BufferException.class, () -> encoding.decodeLong(336, encoded));
    }

    @Test
    void decodeStringWithInvalidFloat() {
        byte[] encoded = "not.a.float".getBytes(StandardCharsets.UTF_16BE);
        assertThrows(BufferException.class, () -> encoding.decodeFloat(176, encoded));
    }

    @Test
    void decodeStringWithInvalidDouble() {
        byte[] encoded = "not.a.double".getBytes(StandardCharsets.UTF_16BE);
        assertThrows(BufferException.class, () -> encoding.decodeDouble(192, encoded));
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
    void testByteOrderMark() throws BufferException {
        // UTF-16BE should not include BOM
        String input = "A";
        byte[] encoded = encoding.encodeString(16, input);
        assertEquals(2, encoded.length);
        assertEquals(0, encoded[0]); // High byte
        assertEquals(65, encoded[1]); // Low byte ('A')
    }

    @Test
    void testSurrogatePairEncoding() throws BufferException {
        // Test explicit surrogate pair encoding
        String input = "𝄞"; // Musical G clef (U+1D11E)
        byte[] encoded = encoding.encodeString(32, input);
        assertEquals(4, encoded.length);
        // Should be encoded as surrogate pair D834 DD1E
        assertEquals((byte) 0xD8, encoded[0]);
        assertEquals((byte) 0x34, encoded[1]);
        assertEquals((byte) 0xDD, encoded[2]);
        assertEquals((byte) 0x1E, encoded[3]);
    }
}