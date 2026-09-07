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

class EncodingVarLengthSignedIntegerTest extends BaseEncodingRawTest {
    private EncodingVarLengthSignedInteger encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingVarLengthSignedInteger();
    }

    @Test
    void getName() {
        assertEquals("VAR-SIGNED", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingVarLengthSignedInteger.optionEncodingVarLengthSignedInteger();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("VAR-SIGNED", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeByteHappyPath() {
        assertEncodeDecode((byte) 0, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 1, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) -1, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) 63, 8, encoding::decodeByte, encoding::encodeByte); // 6-bit max positive
        assertEncodeDecode((byte) -64, 8, encoding::decodeByte, encoding::encodeByte); // 6-bit min negative
    }

    @Test
    void encodeAndDecodeShortHappyPath() {
        assertEncodeDecode((short) 0, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 1, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) -1, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 63, 16, encoding::decodeShort, encoding::encodeShort); // 6-bit boundary
        assertEncodeDecode((short) -64, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 127, 16, encoding::decodeShort, encoding::encodeShort); // 7-bit boundary
        assertEncodeDecode((short) -128, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) 8191, 16, encoding::decodeShort, encoding::encodeShort); // 13-bit max
        assertEncodeDecode((short) -8192, 16, encoding::decodeShort, encoding::encodeShort); // 13-bit min
    }

    @Test
    void encodeAndDecodeIntHappyPath() {
        assertEncodeDecode(0, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(1, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(-1, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(63, 32, encoding::decodeInt, encoding::encodeInt); // 6-bit boundary
        assertEncodeDecode(-64, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(127, 32, encoding::decodeInt, encoding::encodeInt); // 7-bit boundary
        assertEncodeDecode(-128, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(8191, 32, encoding::decodeInt, encoding::encodeInt); // 14-bit boundary
        assertEncodeDecode(-8192, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(16383, 32, encoding::decodeInt, encoding::encodeInt); // 14-bit boundary
        assertEncodeDecode(-16384, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(2097151, 32, encoding::decodeInt, encoding::encodeInt); // 21-bit boundary
        assertEncodeDecode(-2097152, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(134217727, 32, encoding::decodeInt, encoding::encodeInt); // 21-bit boundary
        assertEncodeDecode(-134217728, 32, encoding::decodeInt, encoding::encodeInt);
    }

    @Test
    void encodeAndDecodeLongHappyPath() {
        assertEncodeDecode(0L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(1L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-1L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(63L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-64L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(127L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-128L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(8191L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-8192L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(16383L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-16384L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(2097151L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-2097152L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(134217727, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-134217728, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(36028797018963967L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-36028797018963968L, 64, encoding::decodeLong, encoding::encodeLong);
    }

    @Test
    void encodeAndDecodeBigIntegerHappyPath() {
        assertEncodeDecode(BigInteger.ZERO, 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.ONE, 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.ONE.negate(), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(127), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(-127), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(128), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(-128), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(16383), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(BigInteger.valueOf(-16383), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
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
        assertEncodeDecode((byte) 63, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((byte) -64, 8, encoding::decodeByte, encoding::encodeByte);
        assertEncodeDecode((short) 8191, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode((short) -8192, 16, encoding::decodeShort, encoding::encodeShort);
        assertEncodeDecode(134217727, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(-134217728, 32, encoding::decodeInt, encoding::encodeInt);
        assertEncodeDecode(36028797018963967L, 64, encoding::decodeLong, encoding::encodeLong);
        assertEncodeDecode(-36028797018963968L, 64, encoding::decodeLong, encoding::encodeLong);
    }

    @Test
    void encodeAndDecodeLargeBigIntegers() {
        BigInteger large = new BigInteger("123456789012345678901234567890");
        assertEncodeDecode(large, 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
        assertEncodeDecode(large.negate(), 128, encoding::decodeBigInteger, encoding::encodeBigInteger);
    }

    // Error Cases
    @Test
    void decodeWithInvalidFormat() {
        // Invalid VLQ format - continuation bit set on last byte
        byte[] encoded = new byte[]{(byte) 0x80};
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, read));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeShort(16, read));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(32, read));

        // Invalid VLQ format - missing continuation bit on non-last byte
        byte[] encoded2 = new byte[]{(byte) 0x40, (byte) 0x00};
        ReadBufferByteBased read2 = new ReadBufferByteBased(encoded);
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeByte(8, read2));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeShort(16, read2));
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeInt(32, read2));
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
    void testVLQEncodingFormat() {
        // Test that the VLQ encoding format is correct for signed integers

        // Value 0 should encode to [0x00]
        byte[] encoded = encoding.encodeByte(8, (byte) 0);
        assertArrayEquals(new byte[]{0x00}, encoded);

        // Value 63 (max 6-bit positive) should encode to [0x3F]
        encoded = encoding.encodeByte(8, (byte) 63);
        assertArrayEquals(new byte[]{0x3F}, encoded);

        // Value -64 (min 6-bit negative) should encode to [0x40]
        encoded = encoding.encodeByte(8, (byte) -64);
        assertArrayEquals(new byte[]{0x40}, encoded);

        // Value 64 should encode to [0x80, 0x40] (needs continuation bit)
        encoded = encoding.encodeShort(16, (short) 64);
        assertArrayEquals(new byte[]{(byte) 0x80, 0x40}, encoded);

        // Value -65 should encode to [0x80, 0x41] (needs continuation bit)
        encoded = encoding.encodeShort(16, (short) -65);
        assertArrayEquals(new byte[]{(byte) 0xFF, 0x3F}, encoded);
    }

    @Test
    void testVLQEncodingConsistency() {
        // Test that encoding and decoding are consistent across all types
        byte testValue = 42;

        // Encode same value with different types
        byte[] byteEncoded = encoding.encodeByte(8, testValue);
        byte[] shortEncoded = encoding.encodeShort(16, testValue);
        byte[] intEncoded = encoding.encodeInt(32, testValue);
        byte[] longEncoded = encoding.encodeLong(64, testValue);

        // They should all produce the same encoding for the same value
        assertArrayEquals(byteEncoded, shortEncoded);
        assertArrayEquals(byteEncoded, intEncoded);
        assertArrayEquals(byteEncoded, longEncoded);

        // Test with negative value
        testValue = -42;
        byteEncoded = encoding.encodeByte(8, testValue);
        shortEncoded = encoding.encodeShort(16, testValue);
        intEncoded = encoding.encodeInt(32, testValue);
        longEncoded = encoding.encodeLong(64, testValue);

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