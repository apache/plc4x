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
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferRaw;

import java.math.BigInteger;

/**
 * Big-endian variable-length unsigned integer (VLQ):
 * - Each byte contributes 7 data bits (bits 0..6).
 * - Bit 7 (0x80) is the continuation bit: 1 = more bytes follow, 0 = last byte.
 *
 * Field size is provided as numBits; the maximum number of bytes permitted in the field is:
 *   maxBytes = ceil(numBits / 8).
 * Encoding/decoding must not exceed this limit.
 */
public class EncodingVarLengthUnsignedInteger extends BaseEncodingRaw {

    private static final int LAST_SEVEN_BITS = 0x7F;
    private static final int EIGHTH_BIT      = 0x80; // continuation bit

    public static final String NAME = "VAR-UNSIGNED";

    public static WithOption optionEncodingVarLengthUnsignedInteger() {
        return WithOption.WithEncoding(NAME);
    }

    @Override
    public String getName() {
        return NAME;
    }

    // ------------------------------------------------------------
    // byte / short / int variants delegate to long
    // ------------------------------------------------------------

    @Override
    public byte[] encodeByte(int numBits, byte value) {
        return encodeLong(numBits, value & 0xFFL);
    }

    @Override
    public byte decodeByte(int numBits, ReadBufferRaw readBuffer) {
        return (byte) decodeLong(numBits, readBuffer);
    }

    @Override
    public byte[] encodeShort(int numBits, short value) {
        return encodeLong(numBits, value & 0xFFFFL);
    }

    @Override
    public short decodeShort(int numBits, ReadBufferRaw readBuffer) {
        return (short) decodeLong(numBits, readBuffer);
    }

    @Override
    public byte[] encodeInt(int numBits, int value) {
        return encodeLong(numBits, value & 0xFFFFFFFFL);
    }

    @Override
    public int decodeInt(int numBits, ReadBufferRaw readBuffer) {
        return (int) decodeLong(numBits, readBuffer);
    }

    // ------------------------------------------------------------
    // long <-> bytes (big-endian unsigned VLQ)
    // ------------------------------------------------------------

    @Override
    public byte[] encodeLong(int numBits, long value) {
        if (value < 0) {
            throw new IllegalArgumentException("Cannot encode negative value in unsigned format: " + value);
        }
        if (numBits < 8) {
            throw new IllegalArgumentException("At least 8 bits required for var-length integers");
        }

        // Budget (bytes allowed in the field)
        final int maxBytes = (numBits + 7) / 8;

        // Maximum encodable value with maxBytes groups of 7 bits each.
        // For long, cap at 63 data bits to avoid (1L << 64) overflow.
        final int dataBits = Math.min(63, maxBytes * 7);
        final long maxValue = (dataBits == 63) ? Long.MAX_VALUE : ((1L << dataBits) - 1);

        if (value > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (max value: %d)", value, numBits, maxValue));
        }

        // Determine number of 7-bit groups required.
        int numBytes = 0;
        long tmp = value;
        do {
            numBytes++;
            tmp >>>= 7;
        } while (tmp != 0);

        if (numBytes > maxBytes) {
            throw new IllegalArgumentException(
                String.format("Value %d requires %d bytes but field allows only %d bytes (numBits=%d)",
                    value, numBytes, maxBytes, numBits));
        }

        // Serialize big-endian: most-significant 7-bit group first.
        byte[] out = new byte[numBytes];
        for (int i = numBytes - 1; i >= 0; i--) {
            int shift = i * 7;
            int b = (int) ((value >> shift) & LAST_SEVEN_BITS);
            if (i > 0) {
                b |= EIGHTH_BIT; // continuation on all but last
            }
            out[(numBytes - 1) - i] = (byte) b;
        }
        return out;
    }

    @Override
    public long decodeLong(int numBits, ReadBufferRaw readBuffer) {
        if (readBuffer == null) {
            throw new IllegalArgumentException("Cannot decode with null ReadBufferRaw");
        }
        if (numBits < 8) {
            throw new IllegalArgumentException("At least 8 bits required for reading var-length integers");
        }

        final int maxBytes = (numBits + 7) / 8;

        long result = 0;
        for (int i = 0; i < maxBytes; i++) {
            if (readBuffer.getRemainingBits() < 8) {
                throw new IllegalArgumentException("Not enough bits remaining to read next byte");
            }

            final byte b;
            try {
                b = readBuffer.readBits(8)[0];
            } catch (BufferException e) {
                throw new IllegalArgumentException("Cannot decode byte array", e);
            }

            // Last permitted byte must not have continuation set.
            if ((i == maxBytes - 1) && ((b & 0x80) != 0)) {
                throw new IllegalArgumentException("Continuation bit set on last permitted byte");
            }

            // Big-endian VLQ: shift by 7 each step, OR the low 7 bits.
            result = (result << 7) | (b & 0x7F);

            // Stop on last (continuation=0)
            if ((b & 0x80) == 0) {
                return result;
            }
        }

        throw new IllegalArgumentException("Truncated var-length unsigned long: no terminating byte within field size");
    }

    // ------------------------------------------------------------
    // BigInteger <-> bytes (big-endian unsigned VLQ)
    // ------------------------------------------------------------

    @Override
    public byte[] encodeBigInteger(int numBits, BigInteger value) {
        if (value == null) {
            throw new IllegalArgumentException("Cannot encode null BigInteger");
        }
        if (value.signum() < 0) {
            throw new IllegalArgumentException("Cannot encode negative value in unsigned format: " + value);
        }
        if (numBits < 8) {
            throw new IllegalArgumentException("At least 8 bits required for var-length integers");
        }

        final int maxBytes = ((numBits + 7) / 8) + 1;
        final BigInteger maxValue = BigInteger.ONE.shiftLeft(maxBytes * 7).subtract(BigInteger.ONE);

        if (value.compareTo(maxValue) > 0) {
            throw new IllegalArgumentException(
                String.format("Value %s cannot be encoded in %d bits (max value: %s)", value, numBits, maxValue));
        }

        // Determine number of 7-bit groups required.
        int numBytes = 0;
        BigInteger tmp = value;
        do {
            numBytes++;
            tmp = tmp.shiftRight(7);
        } while (!tmp.equals(BigInteger.ZERO));

        if (numBytes > maxBytes) {
            throw new IllegalArgumentException(
                String.format("Value %s requires %d bytes but field allows only %d bytes (numBits=%d)",
                    value, numBytes, maxBytes, numBits));
        }

        // Serialize big-endian groups.
        byte[] out = new byte[numBytes];
        for (int i = numBytes - 1; i >= 0; i--) {
            int shift = i * 7;
            int b = value.shiftRight(shift).and(BigInteger.valueOf(LAST_SEVEN_BITS)).intValue();
            if (i > 0) {
                b |= EIGHTH_BIT;
            }
            out[(numBytes - 1) - i] = (byte) b;
        }
        return out;
    }

    @Override
    public BigInteger decodeBigInteger(int numBits, ReadBufferRaw readBuffer) {
        if (readBuffer == null) {
            throw new IllegalArgumentException("Cannot decode with null ReadBufferRaw");
        }
        if (numBits < 8) {
            throw new IllegalArgumentException("At least 8 bits required for reading var-length integers");
        }

        final int maxBytes = ((numBits + 7) / 8) + 1;

        BigInteger result = BigInteger.ZERO;
        for (int i = 0; i < maxBytes; i++) {
            if (readBuffer.getRemainingBits() < 8) {
                throw new IllegalArgumentException("Not enough bits remaining to read next byte");
            }

            final byte b;
            try {
                b = readBuffer.readBits(8)[0];
            } catch (BufferException e) {
                throw new IllegalArgumentException("Cannot decode byte array", e);
            }

            if ((i == maxBytes - 1) && ((b & 0x80) != 0)) {
//                throw new IllegalArgumentException("Continuation bit set on last permitted byte");
            }

            result = result.shiftLeft(7)
                .or(BigInteger.valueOf(b & 0x7F));

            if ((b & 0x80) == 0) {
                return result;
            }
        }

        throw new IllegalArgumentException("Truncated var-length unsigned BigInteger: no terminating byte within field size");
    }

}
