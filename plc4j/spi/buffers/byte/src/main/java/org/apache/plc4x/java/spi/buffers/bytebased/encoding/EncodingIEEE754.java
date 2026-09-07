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
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;

import java.math.BigDecimal;
import java.math.BigInteger;

public class EncodingIEEE754 extends BaseEncodingDefault {

    public static final String NAME = "IEEE754";

    private static final WithOption OPTION = WithOption.WithEncoding(NAME);

    public static WithOption optionEncodingIEEE754() {
        return OPTION;
    }

    @Override
    public String getName() {
        return NAME;
    }

    @Override
    public byte[] encodeFloat(int numBits, float value) {
        if (numBits == 32) {
            // For 32-bit float, all values are supported as IEEE 754 can represent them
            int bits = Float.floatToIntBits(value);
            return new byte[]{
                (byte) ((bits >>> 24) & 0xFF),
                (byte) ((bits >>> 16) & 0xFF),
                (byte) ((bits >>> 8) & 0xFF),
                (byte) (bits & 0xFF)
            };
        } else if (numBits == 16) {
            // For 16-bit float, check if value is within representable range
            float absValue = Math.abs(value);
            if (!Float.isNaN(value) && !Float.isInfinite(value) &&
                absValue != 0.0f &&
                (absValue < 6.104e-5f || absValue > 65504.0f)) {
                throw new IllegalArgumentException(
                    String.format("Value %e cannot be represented as 16-bit float (valid range: ±6.104e-5 to ±65504)",
                        value));
            }
            return encodeHalfPrecision(value);
        } else {
            throw new IllegalArgumentException("Only 16-bit and 32-bit IEEE 754 floats are supported");
        }
    }

    private byte[] encodeHalfPrecision(float value) {
        int intBits = Float.floatToIntBits(value);
        int sign = (intBits >>> 16) & 0x8000; // sign bit

        int exponent = ((intBits >>> 23) & 0xFF) - 127 + 15;
        int mantissa = (intBits >>> 13) & 0x3FF; // 10-bit mantissa

        if (exponent <= 0) {
            // Subnormal or underflow
            if (exponent < -10) {
                return new byte[]{(byte) (sign >>> 8), (byte) sign};
            }
            mantissa = (intBits & 0x7FFFFF) | 0x800000;
            int shift = 1 - exponent;
            mantissa = mantissa >>> (13 + shift);
            return new byte[]{
                (byte) ((sign | mantissa) >>> 8),
                (byte) (sign | mantissa)
            };
        } else if (exponent >= 31) {
            // Infinity or NaN
            int half = sign | 0x7C00 | ((intBits & 0x7FFFFF) != 0 ? 0x200 : 0);
            return new byte[]{(byte) (half >>> 8), (byte) half};
        } else {
            int half = sign | (exponent << 10) | mantissa;
            return new byte[]{(byte) (half >>> 8), (byte) half};
        }
    }

    @Override
    public float decodeFloat(int numBits, byte[] bytes) {
        if (numBits == 32) {
            if (bytes.length < 4) {
                throw new IllegalArgumentException("At least 4 bytes required for 32-bit float");
            }
            int bits = ((bytes[0] & 0xFF) << 24) |
                ((bytes[1] & 0xFF) << 16) |
                ((bytes[2] & 0xFF) << 8) |
                (bytes[3] & 0xFF);
            return Float.intBitsToFloat(bits);
        } else if (numBits == 16) {
            if (bytes.length < 2) {
                throw new IllegalArgumentException("At least 2 bytes required for 16-bit float");
            }
            int bits = ((bytes[0] & 0xFF) << 8) | (bytes[1] & 0xFF);
            return decodeHalfPrecision(bits);
        } else {
            throw new IllegalArgumentException("Only 16-bit and 32-bit IEEE 754 floats are supported");
        }
    }

    private float decodeHalfPrecision(int halfBits) {
        int sign = (halfBits >>> 15) & 0x1;
        int exponent = (halfBits >>> 10) & 0x1F;
        int mantissa = halfBits & 0x3FF;

        int fullBits;
        if (exponent == 0) {
            if (mantissa == 0) {
                // Zero
                fullBits = sign << 31;
            } else {
                // Subnormal
                exponent = 1;
                while ((mantissa & 0x400) == 0) {
                    mantissa <<= 1;
                    exponent--;
                }
                mantissa &= 0x3FF;
                exponent = exponent - 15 + 127;
                fullBits = (sign << 31) | (exponent << 23) | (mantissa << 13);
            }
        } else if (exponent == 0x1F) {
            // Infinity or NaN
            fullBits = (sign << 31) | 0x7F800000 | (mantissa << 13);
        } else {
            // Normalized
            exponent = exponent - 15 + 127;
            fullBits = (sign << 31) | (exponent << 23) | (mantissa << 13);
        }

        return Float.intBitsToFloat(fullBits);
    }

    @Override
    public byte[] encodeDouble(int numBits, double value) {
        if (numBits != 64) {
            throw new IllegalArgumentException("Only 64-bit IEEE 754 doubles are supported");
        }

        long bits = Double.doubleToLongBits(value);
        return new byte[]{
            (byte) ((bits >>> 56) & 0xFF),
            (byte) ((bits >>> 48) & 0xFF),
            (byte) ((bits >>> 40) & 0xFF),
            (byte) ((bits >>> 32) & 0xFF),
            (byte) ((bits >>> 24) & 0xFF),
            (byte) ((bits >>> 16) & 0xFF),
            (byte) ((bits >>> 8) & 0xFF),
            (byte) (bits & 0xFF)
        };
    }

    @Override
    public double decodeDouble(int numBits, byte[] bytes) {
        if (numBits != 64) {
            throw new IllegalArgumentException("Only 64-bit IEEE 754 doubles are supported");
        }
        if (bytes.length < 8) {
            throw new IllegalArgumentException("At least 8 bytes are required to decode a 64-bit double");
        }

        long bits = ((long) (bytes[0] & 0xFF) << 56) |
            ((long) (bytes[1] & 0xFF) << 48) |
            ((long) (bytes[2] & 0xFF) << 40) |
            ((long) (bytes[3] & 0xFF) << 32) |
            ((long) (bytes[4] & 0xFF) << 24) |
            ((long) (bytes[5] & 0xFF) << 16) |
            ((long) (bytes[6] & 0xFF) << 8) |
            ((long) (bytes[7] & 0xFF));

        return Double.longBitsToDouble(bits);
    }

    @Override
    public byte[] encodeBigDecimal(int numBits, BigDecimal value) {
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[(numBits + 7) / 8], EncodingTwosComplement.optionEncodingTwosComplement());

        // Convert BigDecimal to BigInteger by scaling and then getting unscaled value
        BigInteger unscaledValue = value.unscaledValue();
        int scale = value.scale();

        // First 32 bits for scale, rest for unscaled value
        int scaleNumBits = Math.min(32, numBits);
        int valueNumBits = numBits - scaleNumBits;

        try {
            // Write scale
            writeBuffer.writeSignedInt(scaleNumBits, scale);

            // Write unscaled value if there are bits left
            if (valueNumBits > 0) {
                writeBuffer.writeSignedBigInteger(valueNumBits, unscaledValue);
            }
            return writeBuffer.getBytes();
        } catch (BufferException e) {
            throw new RuntimeException("Error encoding BigDecimal", e);
        }
    }

    @Override
    public BigDecimal decodeBigDecimal(int numBits, byte[] bytes) {
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(bytes, EncodingTwosComplement.optionEncodingTwosComplement());

        // First 32 bits for scale. Rest for unscaled value
        int scaleNumBits = Math.min(32, numBits);
        int valueNumBits = numBits - scaleNumBits;

        try {
            // Read scale
            int scale = readBuffer.readSignedInt(scaleNumBits);

            // Read unscaled value if there are bits left
            BigInteger unscaledValue = BigInteger.ZERO;
            if (valueNumBits > 0) {
                unscaledValue = readBuffer.readSignedBigInteger(valueNumBits);
            }

            return new BigDecimal(unscaledValue, scale);
        } catch (BufferException e) {
            throw new RuntimeException("Error decoding BigDecimal", e);
        }
    }

}
