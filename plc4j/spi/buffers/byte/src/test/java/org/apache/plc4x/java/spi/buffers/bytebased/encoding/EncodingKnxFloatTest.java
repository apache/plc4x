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

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class EncodingKnxFloatTest extends BaseEncodingDefaultTest {

    private EncodingKnxFloat encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingKnxFloat();
    }

    @Test
    void getName() {
        assertEquals("KNXFloat", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingKnxFloat.optionEncodingKNXFloat();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("KNXFloat", stringEncoding.get());
    }

    // Happy Path Tests
    @Test
    void encodeAndDecodeFloatHappyPath() {
        // Test typical KNX float values
        //assertEncodeDecode(0.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.01f);
        //assertEncodeDecode(1.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.01f);
        assertEncodeDecode(-1.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.01f);
        assertEncodeDecode(10.5f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.01f);
        assertEncodeDecode(-10.5f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.01f);
        assertEncodeDecode(100.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.1f);
        assertEncodeDecode(-100.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.1f);
    }

    @Test
    void encodeAndDecodeSmallValues() {
        // Test small values that should work with KNX float precision
        assertEncodeDecode(0.01f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
        assertEncodeDecode(0.1f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
        assertEncodeDecode(0.5f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
        assertEncodeDecode(-0.01f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
        assertEncodeDecode(-0.1f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
        assertEncodeDecode(-0.5f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
    }

    @Test
    void encodeAndDecodeLargeValues() {
        // Test larger values within KNX float range
        assertEncodeDecode(1000.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 1.0f);
        assertEncodeDecode(-1000.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 1.0f);
        assertEncodeDecode(5000.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 10.0f);
        assertEncodeDecode(-5000.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 10.0f);
    }

    // Edge Cases
    @Test
    void encodeAndDecodeZero() {
        byte[] encoded = encoding.encodeFloat(16, 0.0f);
        assertNotNull(encoded);
        assertEquals(2, encoded.length);
        assertEquals(0.0f, encoding.decodeFloat(16, encoded), 0.0f);
    }

    @Test
    void encodeAndDecodePositiveZero() {
        float positiveZero = +0.0f;
        byte[] encoded = encoding.encodeFloat(16, positiveZero);
        float decoded = encoding.decodeFloat(16, encoded);
        assertEquals(0.0f, decoded, 0.0f);
    }

    @Test
    void encodeAndDecodeNegativeZero() {
        float negativeZero = -0.0f;
        byte[] encoded = encoding.encodeFloat(16, negativeZero);
        float decoded = encoding.decodeFloat(16, encoded);
        assertEquals(0.0f, decoded, 0.0f);
    }

    @Test
    void encodeAndDecodeMaxMantissaValue() {
        // Test value that uses maximum mantissa (2047 * 0.01 = 20.47)
        float maxMantissa = 20.47f;
        byte[] encoded = encoding.encodeFloat(16, maxMantissa);
        float decoded = encoding.decodeFloat(16, encoded);
        assertEquals(maxMantissa, decoded, 0.01f);
    }

    @Test
    void encodeAndDecodeValuesRequiringExponent() {
        // Test values that require exponent scaling
        float largeValue = 4000.0f; // Should require exponent > 0
        byte[] encoded = encoding.encodeFloat(16, largeValue);
        float decoded = encoding.decodeFloat(16, encoded);
        assertEquals(largeValue, decoded, 50.0f); // Allow larger delta due to scaling
    }

    // Error Cases
    @ParameterizedTest
    @ValueSource(ints = {8, 24, 32, 64})
    void encodeFloatWithUnsupportedBitLength(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeFloat(numBits, 1.0f));
    }

    @ParameterizedTest
    @ValueSource(ints = {8, 24, 32, 64})
    void decodeFloatWithUnsupportedBitLength(int numBits) {
        byte[] bytes = new byte[numBits / 8];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeFloat(numBits, bytes));
    }

    @Test
    void encodeFloatValueOutOfRange() {
        // Test values that are too large for KNX float format
        float tooLarge = 1000000.0f; // Should exceed maximum exponent
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeFloat(16, tooLarge));
    }

    @Test
    void encodeFloatExponentExceedsMax() {
        // Any value that would require an exponent > 0x0F (e.g., beyond the max representable range)
        // should cause an IllegalArgumentException during encoding.
        float requiresExponentBeyond15 = 700_000.0f; // This would require exponent 16+
        assertThrows(IllegalArgumentException.class,
            () -> encoding.encodeFloat(16, requiresExponentBeyond15));
    }

    @Test
    void decodeFloatWithInvalidByteCount() {
        // KNX float requires exactly 2 bytes
        byte[] oneByteOnly = new byte[1];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeFloat(16, oneByteOnly));

        byte[] threeBytes = new byte[3];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeFloat(16, threeBytes));

        byte[] emptyBytes = new byte[0];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeFloat(16, emptyBytes));
    }

    @Test
    void decodeFloatWithNull() {
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeFloat(16, null));
    }

    @Test
    void encodeFloatSpecialValues() {
        // Test special float values
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeFloat(16, Float.POSITIVE_INFINITY));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeFloat(16, Float.NEGATIVE_INFINITY));
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeFloat(16, Float.NaN));
    }

    @Test
    void testPrecisionLimits() {
        // Test that KNX float has limited precision (scaled by 100)
        float preciseValue = 1.234567f;
        byte[] encoded = encoding.encodeFloat(16, preciseValue);
        float decoded = encoding.decodeFloat(16, encoded);

        // KNX float should round to nearest 0.01
        assertEquals(1.23f, decoded, 0.01f);
    }

    @Test
    void testNegativeValueEncoding() {
        // Test that negative values are properly encoded using two's complement
        float negativeValue = -123.45f;
        byte[] encoded = encoding.encodeFloat(16, negativeValue);
        assertNotNull(encoded);
        assertEquals(2, encoded.length);

        // First bit should be set for negative values
        assertTrue((encoded[0] & 0x80) != 0);

        float decoded = encoding.decodeFloat(16, encoded);
        assertEquals(negativeValue, decoded, 0.1f);
    }

    @Test
    void testExponentScaling() {
        // Test values that require different exponent values
        float[] testValues = {20.0f, 40.0f, 80.0f, 160.0f, 320.0f};

        for (float value : testValues) {
            byte[] encoded = encoding.encodeFloat(16, value);
            float decoded = encoding.decodeFloat(16, encoded);
            assertEquals(value, decoded, value * 0.05f); // Allow 5% error for scaled values
        }
    }

    // Unsupported Operations
    @Test
    void unsupportedOperations() {
        assertUnsupportedOperation(() -> encoding.encodeByte(8, (byte) 1));
        assertUnsupportedOperation(() -> encoding.decodeByte(8, new byte[1]));

        assertUnsupportedOperation(() -> encoding.encodeShort(16, (short) 1));
        assertUnsupportedOperation(() -> encoding.decodeShort(16, new byte[2]));

        assertUnsupportedOperation(() -> encoding.encodeInt(32, 1));
        assertUnsupportedOperation(() -> encoding.decodeInt(32, new byte[4]));

        assertUnsupportedOperation(() -> encoding.encodeLong(64, 1L));
        assertUnsupportedOperation(() -> encoding.decodeLong(64, new byte[8]));

        assertUnsupportedOperation(() -> encoding.encodeBigInteger(128, java.math.BigInteger.ONE));
        assertUnsupportedOperation(() -> encoding.decodeBigInteger(128, new byte[16]));

        assertUnsupportedOperation(() -> encoding.encodeDouble(64, 1.0));
        assertUnsupportedOperation(() -> encoding.decodeDouble(64, new byte[8]));

        assertUnsupportedOperation(() -> encoding.encodeBigDecimal(128, java.math.BigDecimal.ONE));
        assertUnsupportedOperation(() -> encoding.decodeBigDecimal(128, new byte[16]));

        assertUnsupportedOperation(() -> encoding.encodeString(64, "test"));
        assertUnsupportedOperation(() -> encoding.decodeString(64, new byte[8]));
    }
}