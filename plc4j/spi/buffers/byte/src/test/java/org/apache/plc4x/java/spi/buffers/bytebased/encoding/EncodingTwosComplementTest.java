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

class EncodingTwosComplementTest extends BaseEncodingDefaultTest {

    private EncodingTwosComplement encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingTwosComplement();
    }

    @Test
    void getName() {
        assertEquals("twos-complement", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingTwosComplement.optionEncodingTwosComplement();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("twos-complement", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeByteHappyPath() {
        assertEncodeDecode((byte) 42, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) -42, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 127, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) -128, 8, encoding::decodeByte, encoding::encodeByte);
    }

    @Test
    void encodeAndDecodeShortHappyPath() {
        assertEncodeDecode((short) 12345, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) -12345, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(Short.MAX_VALUE, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(Short.MIN_VALUE, 16, encoding::decodeShort, encoding::encodeShort);
    }

    @Test
    void encodeAndDecodeIntHappyPath() {
        assertEncodeDecode(1234567890, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(-1234567890, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(Integer.MAX_VALUE, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(Integer.MIN_VALUE, 32, encoding::decodeInt, encoding::encodeInt);
    }

    @Test
    void encodeAndDecodeLongHappyPath() {
        assertEncodeDecode(1234567890123456789L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-1234567890123456789L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(Long.MAX_VALUE, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(Long.MIN_VALUE, 64, encoding::decodeLong, encoding::encodeLong);
    }

    @Test
    void encodeAndDecodeBigIntegerHappyPath() {
        BigInteger[] testValues = {
            BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.valueOf(2)),
            BigInteger.valueOf(Long.MIN_VALUE).multiply(BigInteger.valueOf(2)),
            new BigInteger("123456789012345678901234567890"),
            new BigInteger("-123456789012345678901234567890")
        };

        for (BigInteger value : testValues) {
            int numBits = value.bitLength() + 1; // +1 for sign bit
            assertEncodeDecode(value, numBits, encoding::decodeBigInteger, encoding::encodeBigInteger);
        }
    }

    // Edge Cases
    @Test
    void encodeAndDecodeZeroValues() {
        assertEncodeDecode((byte) 0, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((short) 0, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(0, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(0L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(BigInteger.ZERO, 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
    }

    @Test
    void encodeAndDecodeSingleBit() {
        // Test encoding/decoding with just 1 bit (can only represent 0 and -1)
        byte[] encoded = encoding.encodeByte(1, (byte) 0);
        assertEquals(0, encoding.decodeByte(1, encoded));

        encoded = encoding.encodeByte(1, (byte) -1);
        assertEquals(-1, encoding.decodeByte(1, encoded));
    }

    @Test
    void encodeAndDecodePartialBytes() {
        // Test encoding with non-byte-aligned bit counts
        byte[] encoded = encoding.encodeInt(12, 2047); // 12 bits can represent -2048 to 2047
        assertEquals(2047, encoding.decodeInt(12, encoded));

        encoded = encoding.encodeInt(12, -2048);
        assertEquals(-2048, encoding.decodeInt(12, encoded));
    }

    @Test
    void testSignExtension() {
        // Test that negative values are properly sign-extended
        byte[] encoded = encoding.encodeByte(4, (byte) -1); // 4-bit -1 should become 8-bit -1
        assertEquals(-1, encoding.decodeByte(4, encoded));

        encoded = encoding.encodeShort(12, (short) -1);
        assertEquals(-1, encoding.decodeShort(12, encoded));

        encoded = encoding.encodeInt(24, -1);
        assertEquals(-1, encoding.decodeInt(24, encoded));
    }

    // Error Cases
    @ParameterizedTest
    @ValueSource(ints = {0, -1, 9, 10, 16})
    void encodeByteInvalidBitCount(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(numBits, (byte) 1));
    }

    @ParameterizedTest
    @ValueSource(ints = {0, -1, 17, 20, 32})
    void encodeShortInvalidBitCount(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(numBits, (short) 1));
    }

    @ParameterizedTest
    @ValueSource(ints = {0, -1, 33, 40, 64})
    void encodeIntInvalidBitCount(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(numBits, 1));
    }

    @ParameterizedTest
    @ValueSource(ints = {0, -1, 65, 128})
    void encodeLongInvalidBitCount(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(numBits, 1L));
    }

    @Test
    void encodeBigIntegerInvalidBitCount() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(0, BigInteger.ONE));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(-1, BigInteger.ONE));
    }

    @Test
    void decodeWithInsufficientBytes() {
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, new byte[0]));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeShort(16, new byte[1]));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(32, new byte[3]));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeLong(64, new byte[7]));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeBigInteger(128, new byte[15]));
    }

    @Test
    void decodeWithNull() {
        assertThrows(NullPointerException.class, () -> encoding.decodeByte(8, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeShort(16, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeInt(32, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeLong(64, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeBigInteger(128, null));
    }

    @Test
    void testValueRangeEnforcement() {
        // Test that values too large for the specified number of bits are properly masked
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(4, (byte) 31));

        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(8, (short) 1000));
    }

    @Test
    void testBitAlignment() {
        // Test that values are properly aligned when not using full bytes
        int value = -1;
        byte[] encoded = encoding.encodeInt(4, value); // Should use only 4 bits
        assertEquals(1, encoded.length); // Should still use a full byte
        // We wrote 4 bits of 0x0F (15) as signed int, that makes this to -1, when parsing the value back, it
        // is correctly parsed as -1
        assertEquals(-1, encoding.decodeInt(4, encoded));

        // Test alignment with multiple bytes
        value = -1;
        encoded = encoding.encodeInt(12, value); // Should use 12 bits
        assertEquals(2, encoded.length); // Should use 2 bytes
        // We wrote 12 bits of 0x0FFF (4095) as signed int, that makes this to -1, when parsing the value back, it
        // is correctly parsed as -1
        assertEquals(-1, encoding.decodeInt(12, encoded));
    }

    // Unsupported Operations
    @Test
    void unsupportedOperations() {
        assertUnsupportedOperation(() -> encoding.encodeFloat(32, 1.0f));
        assertUnsupportedOperation(() -> encoding.decodeFloat(32, new byte[4]));

        assertUnsupportedOperation(() -> encoding.encodeDouble(64, 1.0));
        assertUnsupportedOperation(() -> encoding.decodeDouble(64, new byte[8]));

        assertUnsupportedOperation(() -> encoding.encodeBigDecimal(128, java.math.BigDecimal.ONE));
        assertUnsupportedOperation(() -> encoding.decodeBigDecimal(128, new byte[16]));

        assertUnsupportedOperation(() -> encoding.encodeString(64, "test"));
        assertUnsupportedOperation(() -> encoding.decodeString(64, new byte[8]));
    }
}