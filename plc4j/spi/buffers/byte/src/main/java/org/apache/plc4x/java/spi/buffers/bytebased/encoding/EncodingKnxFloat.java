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

/**
 * TODO: This should be moved into the KNX driver
 */
public class EncodingKnxFloat extends BaseEncodingDefault {

    public static final String NAME = "KNXFloat";

    private static final WithOption OPTION = WithOption.WithEncoding(NAME);

    public static WithOption optionEncodingKNXFloat() {
        return OPTION;
    }

    @Override
    public String getName() {
        return NAME;
    }

    @Override
    public byte[] encodeFloat(int numBits, float value) {
        if (numBits != 16) {
            throw new IllegalArgumentException("Only 16 bit encoding supported");
        }
        if (Float.isNaN(value)) {
            throw new IllegalArgumentException("KNX float does not support NaN");
        }
        if (Float.isInfinite(value)) {
            throw new IllegalArgumentException("KNX float does not support Infinity");
        }

        // Calculate the maximum representable value
        // Max mantissa is 0x07FF (2047) with max exponent 0x0F (15)
        float maxValue = 2047.0f * (float) Math.pow(2, 15) / 100.0f; // Scale back from fixed-point
        if (Math.abs(value) > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %e is outside the valid range for KNX float (±%e)",
                    value, maxValue));
        }

        // Scale the value to get the mantissa (fixed-point with 2 decimal places)
        int raw = Math.round(value * 100.0f);

        boolean sign = raw < 0;
        int mantissa = Math.abs(raw);
        int exponent = 0;

        // Normalize mantissa to fit in 11 bits (0x7FF = 2047)
        while (mantissa > 0x07FF) {
            mantissa >>= 1;
            exponent++;
        }

        if (sign) {
            mantissa = (~mantissa + 1) & 0x07FF; // Two's complement in 11 bits
        }

        int result = 0;
        result |= (sign ? 0x8000 : 0x0000);           // Bit 15: sign
        result |= (exponent & 0x0F) << 11;        // Bits 14-11: exponent
        result |= (mantissa & 0x07FF);            // Bits 10-0: mantissa/fraction

        return new byte[]{
            (byte) ((result >> 8) & 0xFF),
            (byte) (result & 0xFF)
        };
    }

    @Override
    public float decodeFloat(int numBits, byte[] bytes) {
        if (numBits != 16) {
            throw new IllegalArgumentException("Only 16 bit encoding supported");
        }
        if (bytes == null) {
            throw new IllegalArgumentException("bytes cannot be null");
        }
        if (bytes.length != 2) {
            throw new IllegalArgumentException("Invalid number of bytes");
        }

        final boolean sign = (bytes[0] & (byte) 0x80) != 0;
        final byte exponent = (byte) ((bytes[0] & (byte) 0x78) >> 3);
        short fraction = (short) ((short) ((bytes[0] & (byte) 0x07) << 8) | (short) (bytes[1] & 0xFF));
        if (sign) {
            fraction = (short) (fraction | 0xF800);
        }
        return (float) (0.01 * fraction * Math.pow(2, exponent));
    }

}
