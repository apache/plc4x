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

class EncodingUnsignedBinaryTest extends BaseEncodingDefaultTest {

    private EncodingUnsignedBinary encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingUnsignedBinary();
    }

    @Test
    void getName() {
        assertEquals("unsigned-binary", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingUnsignedBinary.optionEncodingUnsignedBinary();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("unsigned-binary", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeByteHappyPath() {
        assertEncodeDecode((byte) 42, 7, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 127, 7, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 0, 7, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 1, 1, encoding::decodeByte, encoding::encodeByte);
    }

    @Test
    void encodeAndDecodeShortHappyPath() {
        assertEncodeDecode((short) 12345, 15, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 32767, 15, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 0, 15, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 1, 1, encoding::decodeShort, encoding::encodeShort);
    }

    @Test
    void encodeAndDecodeIntHappyPath() {
        assertEncodeDecode(1234567890, 31, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(Integer.MAX_VALUE, 31, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(0, 31, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(1, 1, encoding::decodeInt, encoding::encodeInt);
    }

    @Test
    void encodeAndDecodeLongHappyPath() {
        assertEncodeDecode(1234567890123456789L, 63, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(Long.MAX_VALUE, 63, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(0L, 63, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(1L, 1, encoding::decodeLong, encoding::encodeLong);
    }

    @Test
    void encodeAndDecodeBigIntegerHappyPath() {
        BigInteger[] testValues = {
            BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.valueOf(2)),
            new BigInteger("123456789012345678901234567890"),
            BigInteger.ZERO,
            BigInteger.ONE
        };

        for (BigInteger value : testValues) {
            int numBits = Math.max(1, value.bitLength());
            assertEncodeDecode(value, numBits, encoding::decodeBigInteger, encoding::encodeBigInteger);
        }
    }

    // Edge Cases
    @Test
    void encodeAndDecodeZeroValues() {
        assertEncodeDecode((byte) 0, 1, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((short) 0, 1, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(0, 1, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(0L, 1, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(BigInteger.ZERO, 1, encoding::decodeBigInteger, encoding::encodeBigInteger);
    }

    @Test
    void encodeAndDecodeSingleBit() {
        // Test encoding/decoding with just 1 bit (can only represent 0 and 1)
        byte[] encoded = encoding.encodeByte(1, (byte) 0);
        assertEquals(0, encoding.decodeByte(1, encoded));

        encoded = encoding.encodeByte(1, (byte) 1);
        assertEquals(1, encoding.decodeByte(1, encoded));
    }

    @Test
    void encodeAndDecodeMaxValues() {
        // Test maximum values for each bit width
        assertEncodeDecode((byte) 127, 7, encoding::decodeByte, encoding::encodeByte); // 7 bits max
        assertEncodeDecode((short) 32767, 15, encoding::decodeShort, encoding::encodeShort); // 15 bits max
        assertEncodeDecode(Integer.MAX_VALUE, 31, encoding::decodeInt, encoding::encodeInt); // 31 bits max
        assertEncodeDecode(Long.MAX_VALUE, 63, encoding::decodeLong, encoding::encodeLong); // 63 bits max
    }

    @Test
    void testBitAlignment() {
        // Test that values are properly aligned when not using full bytes
        int value = 0x0F; // 0000 1111
        byte[] encoded = encoding.encodeInt(4, value); // Should use only 4 bits
        assertEquals(1, encoded.length); // Should still use a full byte
        assertEquals(value, encoding.decodeInt(4, encoded));

        // Test alignment with multiple bytes
        value = 0x0FFF; // 0000 1111 1111 1111
        encoded = encoding.encodeInt(12, value); // Should use 12 bits
        assertEquals(2, encoded.length); // Should use 2 bytes
        assertEquals(value, encoding.decodeInt(12, encoded));
    }

    // Error Cases
    @ParameterizedTest
    @ValueSource(ints = {0, -1, 8, 9, 16})
    void encodeByteInvalidBitCount(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(numBits, (byte) 1));
    }

    @ParameterizedTest
    @ValueSource(ints = {0, -1, 16, 17, 32})
    void encodeShortInvalidBitCount(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(numBits, (short) 1));
    }

    @ParameterizedTest
    @ValueSource(ints = {0, -1, 32, 33, 64})
    void encodeIntInvalidBitCount(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(numBits, 1));
    }

    @ParameterizedTest
    @ValueSource(ints = {0, -1, 64, 65, 128})
    void encodeLongInvalidBitCount(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(numBits, 1L));
    }

    @Test
    void encodeBigIntegerInvalidBitCount() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(0, BigInteger.ONE));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(-1, BigInteger.ONE));
    }

    @Test
    void encodeNegativeValues() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(7, (byte) -1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(15, (short) -1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(31, -1));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(63, -1L));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(128, BigInteger.valueOf(-1)));
    }

    @Test
    void encodeValuesTooLarge() {
        // Test values that exceed the specified number of bits
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(3, (byte) 8)); // 3 bits can only store 0-7
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(8, (short) 256)); // 8 bits can only store 0-255
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(16, 65536)); // 16 bits can only store 0-65535
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(32, 1L << 32)); // 32 bits can only store 0-4294967295
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(64, BigInteger.ONE.shiftLeft(64)));
    }

    @Test
    void decodeWithInsufficientBytes() {
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(7, new byte[0]));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeShort(15, new byte[1]));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(31, new byte[3]));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeLong(63, new byte[7]));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeBigInteger(128, new byte[15]));
    }

    @Test
    void decodeWithNull() {
        assertThrows(NullPointerException.class, () -> encoding.decodeByte(7, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeShort(15, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeInt(31, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeLong(63, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeBigInteger(128, null));
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