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

class EncodingIEEE754Test extends BaseEncodingDefaultTest {

    private EncodingIEEE754 encoding;

    @BeforeEach
    void setUp() {
        encoding = new EncodingIEEE754();
    }

    @Test
    void getName() {
        assertEquals("IEEE754", encoding.getName());
    }

    @Test
    void optionFactory() {
        WithOption withOption = EncodingIEEE754.optionEncodingIEEE754();
        Optional<String> stringEncoding = WithOption.extractEncoding(new WithOption[]{withOption});
        assertTrue(stringEncoding.isPresent());
        assertEquals("IEEE754", stringEncoding.get());
    }

    // Happy Path Tests - 32-bit Float
    @Test
    void encodeAndDecodeFloat32HappyPath() {
        assertEncodeDecode(123.456f, 32, encoding::decodeFloat, encoding::encodeFloat, 0.0001f);
        assertEncodeDecode(-123.456f, 32, encoding::decodeFloat, encoding::encodeFloat, 0.0001f);
        assertEncodeDecode(0.0f, 32, encoding::decodeFloat, encoding::encodeFloat, 0.0f);
        assertEncodeDecode(1.0f, 32, encoding::decodeFloat, encoding::encodeFloat, 0.0f);
        assertEncodeDecode(-1.0f, 32, encoding::decodeFloat, encoding::encodeFloat, 0.0f);
    }

    // Happy Path Tests - 16-bit Float (Half Precision)
    @Test
    void encodeAndDecodeFloat16HappyPath() {
        assertEncodeDecode(1.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
        assertEncodeDecode(-1.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
        assertEncodeDecode(0.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.0f);
        assertEncodeDecode(2.0f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
        assertEncodeDecode(0.5f, 16, encoding::decodeFloat, encoding::encodeFloat, 0.001f);
    }

    // Happy Path Tests - 64-bit Double
    @Test
    void encodeAndDecodeDouble64HappyPath() {
        assertEncodeDecode(123.456789012345, 64, encoding::decodeDouble, encoding::encodeDouble, 0.0000000000001);
        assertEncodeDecode(-123.456789012345, 64, encoding::decodeDouble, encoding::encodeDouble, 0.0000000000001);
        assertEncodeDecode(0.0, 64, encoding::decodeDouble, encoding::encodeDouble, 0.0);
        assertEncodeDecode(1.0, 64, encoding::decodeDouble, encoding::encodeDouble, 0.0);
        assertEncodeDecode(-1.0, 64, encoding::decodeDouble, encoding::encodeDouble, 0.0);
    }

    // Edge Cases - Special Values
    @Test
    void encodeAndDecodeSpecialFloatValues() {
        // Positive and negative zero
        float positiveZero = 0.0f;
        float negativeZero = -0.0f;

        byte[] encodedPosZero = encoding.encodeFloat(32, positiveZero);
        byte[] encodedNegZero = encoding.encodeFloat(32, negativeZero);

        assertEquals(positiveZero, encoding.decodeFloat(32, encodedPosZero));
        assertEquals(negativeZero, encoding.decodeFloat(32, encodedNegZero));
        assertTrue(Float.compare(positiveZero, encoding.decodeFloat(32, encodedPosZero)) == 0);
        assertTrue(Float.compare(negativeZero, encoding.decodeFloat(32, encodedNegZero)) == 0);

        // Infinity
        float posInf = Float.POSITIVE_INFINITY;
        float negInf = Float.NEGATIVE_INFINITY;

        byte[] encodedPosInf = encoding.encodeFloat(32, posInf);
        byte[] encodedNegInf = encoding.encodeFloat(32, negInf);

        assertEquals(posInf, encoding.decodeFloat(32, encodedPosInf));
        assertEquals(negInf, encoding.decodeFloat(32, encodedNegInf));

        // NaN
        float nan = Float.NaN;
        byte[] encodedNaN = encoding.encodeFloat(32, nan);
        assertTrue(Float.isNaN(encoding.decodeFloat(32, encodedNaN)));
    }

    @Test
    void encodeAndDecodeSpecialDoubleValues() {
        // Positive and negative zero
        double positiveZero = 0.0;
        double negativeZero = -0.0;

        byte[] encodedPosZero = encoding.encodeDouble(64, positiveZero);
        byte[] encodedNegZero = encoding.encodeDouble(64, negativeZero);

        assertEquals(positiveZero, encoding.decodeDouble(64, encodedPosZero));
        assertEquals(negativeZero, encoding.decodeDouble(64, encodedNegZero));
        assertTrue(Double.compare(positiveZero, encoding.decodeDouble(64, encodedPosZero)) == 0);
        assertTrue(Double.compare(negativeZero, encoding.decodeDouble(64, encodedNegZero)) == 0);

        // Infinity
        double posInf = Double.POSITIVE_INFINITY;
        double negInf = Double.NEGATIVE_INFINITY;

        byte[] encodedPosInf = encoding.encodeDouble(64, posInf);
        byte[] encodedNegInf = encoding.encodeDouble(64, negInf);

        assertEquals(posInf, encoding.decodeDouble(64, encodedPosInf));
        assertEquals(negInf, encoding.decodeDouble(64, encodedNegInf));

        // NaN
        double nan = Double.NaN;
        byte[] encodedNaN = encoding.encodeDouble(64, nan);
        assertTrue(Double.isNaN(encoding.decodeDouble(64, encodedNaN)));
    }

    @Test
    void encodeAndDecodeMinMaxFloatValues() {
        // Min and max normal values
        float minNormal = Float.MIN_NORMAL;
        float maxValue = Float.MAX_VALUE;

        byte[] encodedMin = encoding.encodeFloat(32, minNormal);
        byte[] encodedMax = encoding.encodeFloat(32, maxValue);

        assertEquals(minNormal, encoding.decodeFloat(32, encodedMin));
        assertEquals(maxValue, encoding.decodeFloat(32, encodedMax));

        // Smallest positive value (subnormal)
        float minValue = Float.MIN_VALUE;
        byte[] encodedMinValue = encoding.encodeFloat(32, minValue);
        assertEquals(minValue, encoding.decodeFloat(32, encodedMinValue));
    }

    @Test
    void encodeAndDecodeMinMaxDoubleValues() {
        // Min and max normal values
        double minNormal = Double.MIN_NORMAL;
        double maxValue = Double.MAX_VALUE;

        byte[] encodedMin = encoding.encodeDouble(64, minNormal);
        byte[] encodedMax = encoding.encodeDouble(64, maxValue);

        assertEquals(minNormal, encoding.decodeDouble(64, encodedMin));
        assertEquals(maxValue, encoding.decodeDouble(64, encodedMax));

        // Smallest positive value (subnormal)
        double minValue = Double.MIN_VALUE;
        byte[] encodedMinValue = encoding.encodeDouble(64, minValue);
        assertEquals(minValue, encoding.decodeDouble(64, encodedMinValue));
    }

    // Error Cases
    @ParameterizedTest
    @ValueSource(ints = {8, 24, 48, 128})
    void encodeFloatWithUnsupportedBitLength(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeFloat(numBits, 1.0f));
    }

    @ParameterizedTest
    @ValueSource(ints = {8, 16, 32, 48, 128})
    void encodeDoubleWithUnsupportedBitLength(int numBits) {
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeDouble(numBits, 1.0));
    }

    @ParameterizedTest
    @ValueSource(ints = {8, 24, 48, 128})
    void decodeFloatWithUnsupportedBitLength(int numBits) {
        byte[] bytes = new byte[numBits / 8];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeFloat(numBits, bytes));
    }

    @ParameterizedTest
    @ValueSource(ints = {8, 16, 32, 48, 128})
    void decodeDoubleWithUnsupportedBitLength(int numBits) {
        byte[] bytes = new byte[numBits / 8];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeDouble(numBits, bytes));
    }

    @Test
    void decodeFloatWithInsufficientBytes() {
        // 32-bit float needs 4 bytes
        byte[] insufficientBytes = new byte[3];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeFloat(32, insufficientBytes));

        // 16-bit float needs 2 bytes
        byte[] oneByteOnly = new byte[1];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeFloat(16, oneByteOnly));
    }

    @Test
    void decodeDoubleWithInsufficientBytes() {
        // 64-bit double needs 8 bytes
        byte[] insufficientBytes = new byte[7];
        assertThrows(IllegalArgumentException.class, () -> encoding.decodeDouble(64, insufficientBytes));
    }

    @Test
    void decodeWithNullBytes() {
        assertThrows(NullPointerException.class, () -> encoding.decodeFloat(32, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeFloat(16, null));
        assertThrows(NullPointerException.class, () -> encoding.decodeDouble(64, null));
    }

    @Test
    void encodeDecodeFloat16Precision() {
        // Test precision limits of 16-bit floats
        float value = 1.5f;
        byte[] encoded = encoding.encodeFloat(16, value);
        float decoded = encoding.decodeFloat(16, encoded);
        assertEquals(value, decoded, 0.001f);

        // Test value that loses precision in 16-bit
        float preciseValue = 1.234567f;
        encoded = encoding.encodeFloat(16, preciseValue);
        decoded = encoding.decodeFloat(16, encoded);
        // 16-bit float has less precision, so we allow larger delta
        assertEquals(preciseValue, decoded, 0.01f);
    }

    @Test
    void encodeDecodeFloat16Overflow() {
        // Values too large for 16-bit should become infinity
        float largeValue = 100000.0f;
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeFloat(16, largeValue));
    }

    @Test
    void encodeDecodeFloat16Underflow() {
        // Very small values should become zero or subnormal
        float smallValue = 1e-10f;
        assertThrows(IllegalArgumentException.class, () -> encoding.encodeFloat(16, smallValue));
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

        //assertUnsupportedOperation(() -> encoding.encodeBigDecimal(128, java.math.BigDecimal.ONE));
        //assertUnsupportedOperation(() -> encoding.decodeBigDecimal(128, new byte[16]));

        assertUnsupportedOperation(() -> encoding.encodeString(64, "test"));
        assertUnsupportedOperation(() -> encoding.decodeString(64, new byte[8]));
    }
}