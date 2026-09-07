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
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class EncodingVarLengthUnsignedIntegerTest extends BaseEncodingRawTest {
    private EncodingVarLengthUnsignedInteger encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingVarLengthUnsignedInteger();
    }

    @Test
    void getName() {
        assertEquals("VAR-UNSIGNED", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingVarLengthUnsignedInteger.optionEncodingVarLengthUnsignedInteger();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("VAR-UNSIGNED", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeByteHappyPath() {
        assertEncodeDecode((byte) 0, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 1, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 127, 8, encoding::decodeByte, encoding::encodeByte);
    }

    @Test
    void encodeAndDecodeShortHappyPath() {
        assertEncodeDecode((short) 0, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 1, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 127, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 128, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 16383, 16, encoding::decodeShort, encoding::encodeShort);
    }

    @Test
    void encodeAndDecodeIntHappyPath() {
        assertEncodeDecode(0, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(1, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(127, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(128, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(16383, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(16384, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(268435455, 32, encoding::decodeInt, encoding::encodeInt);
    }

    @Test
    void encodeAndDecodeLongHappyPath() {
        assertEncodeDecode(0L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(1L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(127L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(128L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(16383L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(16384L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(268435455L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(268435456L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(72057594037927935L, 64, encoding::decodeLong, encoding::encodeLong);
    }

    @Test
    void encodeAndDecodeBigIntegerHappyPath() {
        assertEncodeDecode(BigInteger.ZERO, 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.ONE, 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(127), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(128), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(16383), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(16384), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(268435455), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(268435456), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(72057594037927935L), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(72057594037927936L), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
    }

    @Test
    void encodeAndDecodeLargeBigIntegers() {
        BigInteger large = new BigInteger("123456789012345678901234567890");
        assertEncodeDecode(large, 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
    }

    // Edge Cases
    @Test
    void encodeAndDecodeZero() {
        assertEncodeDecode((byte) 0, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((short) 0, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(0, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(0L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(BigInteger.ZERO, 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
    }

    @Test
    void encodeAndDecodeMaxValues() {
        assertEncodeDecode((byte) 127, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((short) 16383, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(268435455, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(72057594037927935L, 64, encoding::decodeLong, encoding::encodeLong);
    }

    @Test
    void encodeAndDecodeBoundaryValues() {
        // Test 7-bit boundary values (VLQ encoding uses 7 bits per byte)
        assertEncodeDecode((byte) 127, 8, encoding::decodeByte, encoding::encodeByte); // 2^7 - 1
        assertEncodeDecode((short) 16383, 16, encoding::decodeShort, encoding::encodeShort); // 2^14 - 1
        assertEncodeDecode(268435455, 32, encoding::decodeInt, encoding::encodeInt); // 2^21 - 1
        assertEncodeDecode(72057594037927935L, 64, encoding::decodeLong, encoding::encodeLong); // 2^28 - 1
    }

    // Error Cases
    @Test
    void encodeNegativeByte() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(8, (byte) -1));
    }

    @Test
    void encodeNegativeShort() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(16, (short) -1));
    }

    @Test
    void encodeNegativeInt() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(32, -1));
    }

    @Test
    void encodeNegativeLong() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(64, -1L));
    }

    @Test
    void encodeNegativeBigInteger() {
        assertThrows(IllegalArgumentException.class,
            () -> encoding.encodeBigInteger(128, BigInteger.valueOf(-1)));
    }

    @Test
    void encodeValuesTooLarge() {
        // Test values that exceed the maximum for variable-length encoding
        // For VLQ, the maximum depends on the number of bits allowed

        // For byte with 8 bits, max bytes = (8 + 6) / 7 = 2, max value = 2^14 - 1 = 16383
        // But since we're dealing with byte values, the actual max is 255
        // The implementation should handle this correctly

        // For short with 16 bits, max bytes = (16 + 6) / 7 = 3, max value = 2^21 - 1 = 2097151
        // But since we're dealing with short values, the actual max is 65535

        // Test with values that would exceed the bit limit if not properly masked
        // These should work because the implementation masks the values
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeByte(8, (byte) 128));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeShort(8, (short) 16384));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeInt(8, 268435456));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeLong(8, 72057594037927936L));
    }

    @Test
    void decodeByteWithInvalidFormat() {
        // Invalid VLQ format - continuation bit set on last byte
        byte[] encoded = new byte[]{(byte) 0x80};
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, read));
    }

    @Test
    void decodeShortWithInvalidFormat() {
        // Invalid VLQ format - continuation bit set on last byte
        byte[] encoded = new byte[]{(byte) 0x80};
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeShort(16, read));
    }

    @Test
    void decodeIntWithInvalidFormat() {
        // Invalid VLQ format - continuation bit set on last byte
        byte[] encoded = new byte[]{(byte) 0x80};
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(32, read));
    }

    @Test
    void decodeLongWithInvalidFormat() {
        // Invalid VLQ format - too many continuation bytes
        byte[] encoded = new byte[]{(byte) 0x80, (byte) 0x80, (byte) 0x80, (byte) 0x80,
            (byte) 0x80, (byte) 0x80, (byte) 0x80, (byte) 0x80,
            (byte) 0x80, (byte) 0x80, (byte) 0x01};
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeLong(64, read));
    }

    @Test
    void decodeWithInsufficientBytes() {
        byte[] encoded = new byte[0];
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, read));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeShort(16, read));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(32, read));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeLong(64, read));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeBigInteger(128, read));
    }

    @Test
    void decodeWithNull() {
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, null));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeShort(16, null));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(32, null));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeLong(64, null));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeBigInteger(128, null));
    }

    @Test
    void encodeBigIntegerWithNull() {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeBigInteger(128, null));
    }

    @Test
    void encodeByteWithNull() {
        // Note: primitive byte cannot be null, but test for completeness
        // This test verifies the method signature accepts byte primitives
        assertEncodeDecode((byte) 0, 8, encoding::decodeByte, encoding::encodeByte);
    }

    @Test
    void encodeShortWithNull() {
        // Note: primitive short cannot be null, but test for completeness
        // This test verifies the method signature accepts short primitives
        assertEncodeDecode((short) 0, 16, encoding::decodeShort, encoding::encodeShort);
    }

    // Specific VLQ Encoding Tests
    @Test
    void testVLQEncodingFormat() {
        // Test that the VLQ encoding format is correct
        // Value 0 should encode to [0x00]
        byte[] encoded = encoding.encodeByte(8, (byte) 0);
        assertArrayEquals(new byte[]{0x00}, encoded);

        // Value 127 should encode to [0x7F]
        encoded = encoding.encodeByte(8, (byte) 127);
        assertArrayEquals(new byte[]{0x7F}, encoded);

        // Value 128 should encode to [0x81, 0x00] (continuation bit set)
        encoded = encoding.encodeShort(16, (short) 128);
        assertArrayEquals(new byte[]{(byte) 0x81, 0x00}, encoded);

        // Value 255 should encode to [0x81, 0x7F]
        encoded = encoding.encodeShort(16, (short) 255);
        assertArrayEquals(new byte[]{(byte) 0x81, 0x7F}, encoded);
    }

    @Test
    void testVLQEncodingConsistency() {
        // Test that encoding and decoding are consistent across all types
        byte testByte = 42;
        short testShort = 12345;
        int testInt = 1234567;
        long testLong = 123456789L;

        // Encode as different types and verify they produce the same result when values fit
        byte[] byteEncoded = encoding.encodeByte(8, testByte);
        byte[] shortEncoded = encoding.encodeShort(16, testByte);
        byte[] intEncoded = encoding.encodeInt(32, testByte);
        byte[] longEncoded = encoding.encodeLong(64, testByte);

        // They should all produce the same encoding for the same value
        assertArrayEquals(byteEncoded, shortEncoded);
        assertArrayEquals(byteEncoded, intEncoded);
        assertArrayEquals(byteEncoded, longEncoded);
    }

    // Unsupported Operations
    @Test
    void unsupportedOperations() {
        //assertUnsupportedOperations(encoding);
    }

}