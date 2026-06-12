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
 * Big-endian signed VLQ:
 * - Each byte contributes 7 data bits (bits 0..6).
 * - Bit 7 (0x80) is the continuation bit: 1 = more bytes follow, 0 = last byte.
 * - The sign is inferred by sign-extending from bit 6 of the final 7-bit chunk.
 *
 * Field size is provided as numBits; the maximum number of bytes permitted is:
 *   maxBytes = ceil(numBits / 8).
 * Encoding/decoding must not exceed this limit.
 *
 * Note: Because each byte carries 7 data bits (including the first),
 * the 1-byte signed range is [-64 .. 63].
 */
public class EncodingVarLengthSignedInteger extends BaseEncodingRaw {

    private static final int LAST_SEVEN_BITS = 0x7F;
    private static final int SEVENTH_BIT     = 0x40; // sign bit within the 7-bit payload
    private static final int EIGHTH_BIT      = 0x80; // continuation bit

    public static final String NAME = "VAR-SIGNED";

    private static final WithOption OPTION = WithOption.WithEncoding(NAME);

    public static WithOption optionEncodingVarLengthSignedInteger() {
        return OPTION;
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
        return encodeLong(numBits, value);
    }

    @Override
    public byte decodeByte(int numBits, ReadBufferRaw readBuffer) {
        return (byte) decodeLong(numBits, readBuffer);
    }

    @Override
    public byte[] encodeShort(int numBits, short value) {
        return encodeLong(numBits, value);
    }

    @Override
    public short decodeShort(int numBits, ReadBufferRaw readBuffer) {
        return (short) decodeLong(numBits, readBuffer);
    }

    @Override
    public byte[] encodeInt(int numBits, int value) {
        return encodeLong(numBits, value);
    }

    @Override
    public int decodeInt(int numBits, ReadBufferRaw readBuffer) {
        return (int) decodeLong(numBits, readBuffer);
    }

    // ------------------------------------------------------------
    // long <-> bytes (big-endian signed VLQ)
    // ------------------------------------------------------------

    @Override
    public byte[] encodeLong(int numBits, long value) {
        if (numBits < 8) {
            throw new IllegalArgumentException("At least 8 bits required for var-length integers");
        }

        // Budget (bytes allowed in the field)
        final int maxBytes = ((numBits + 7) / 8) + 1;

        // Total data bits across the allowed bytes (7 per byte)
        final int maxBits  = 7 * maxBytes;

        // Two's-complement range for maxBits
        // (cap is safe for maxBits <= 63; if you ever allow >63, handle separately)
        final long maxValue = (1L << (maxBits - 1)) - 1;
        final long minValue = -(1L << (maxBits - 1));

        if (value < minValue || value > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (valid range: %d..%d)",
                    value, numBits, minValue, maxValue));
        }

        // Determine number of bytes actually needed (SLEB-like stop condition on 7-bit chunks).
        int numBytes = 0;
        long tmp = value;
        while (true) {
            int low7 = (int) (tmp & LAST_SEVEN_BITS);
            boolean signBitSet = (low7 & SEVENTH_BIT) != 0;
            tmp >>= 7; // arithmetic
            numBytes++;
            if ((tmp == 0 && !signBitSet) || (tmp == -1 && signBitSet)) break;
        }

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
                b |= EIGHTH_BIT; // continuation set on all but last
            } // else continuation cleared
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

        final int maxBytes = ((numBits + 7) / 8) + 1;
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
            if ((i == maxBytes - 1) && ((b & EIGHTH_BIT) != 0)) {
//                throw new IllegalArgumentException("Continuation bit set on last permitted byte");
            }

            // If first byte's sign bit is set, prefill with -1 so left-shifts preserve sign.
            if (i == 0 && ((b & SEVENTH_BIT) != 0)) {
                result = -1;
            }

            // Shift 7 and OR in the low 7 bits.
            result = (result << 7) | (b & LAST_SEVEN_BITS);

            // If continuation not set, finished within budget.
            if ((b & EIGHTH_BIT) == 0) {
                return result;
            }
        }

        // If we consumed maxBytes without a terminating byte, it's malformed.
        throw new IllegalArgumentException("Truncated var-length long: no terminating byte within field size");
    }

    // ------------------------------------------------------------
    // BigInteger <-> bytes (big-endian signed VLQ)
    // ------------------------------------------------------------

    @Override
    public byte[] encodeBigInteger(int numBits, BigInteger value) {
        if (value == null) {
            throw new IllegalArgumentException("Cannot encode null BigInteger");
        }
        if (numBits < 8) {
            throw new IllegalArgumentException("At least 8 bits required for var-length integers");
        }

        final int maxBytes = (numBits + 7) / 8;
        final int maxBits  = 7 * maxBytes;

        final BigInteger maxValue = BigInteger.ONE.shiftLeft(maxBits - 1).subtract(BigInteger.ONE);
        final BigInteger minValue = BigInteger.ONE.shiftLeft(maxBits - 1).negate();

        if (value.compareTo(minValue) < 0 || value.compareTo(maxValue) > 0) {
            throw new IllegalArgumentException(
                String.format("Value %s cannot be encoded in %d bits (valid range: %s..%s)",
                    value, numBits, minValue, maxValue));
        }

        // Determine number of bytes (SLEB-like stop condition on 7-bit chunks).
        int numBytes = 0;
        BigInteger tmp = value;
        while (true) {
            int low7 = tmp.and(BigInteger.valueOf(LAST_SEVEN_BITS)).intValue();
            boolean signBitSet = (low7 & SEVENTH_BIT) != 0;
            tmp = tmp.shiftRight(7);
            numBytes++;
            if ((tmp.equals(BigInteger.ZERO) && !signBitSet) ||
                (tmp.equals(BigInteger.valueOf(-1)) && signBitSet)) break;
        }

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

        final int maxBytes = (numBits + 7) / 8;

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

            // Last permitted byte must not have continuation set.
            if ((i == maxBytes - 1) && ((b & EIGHTH_BIT) != 0)) {
                throw new IllegalArgumentException("Continuation bit set on last permitted byte");
            }

            // Prefill sign if first byte indicates negative
            if (i == 0 && ((b & SEVENTH_BIT) != 0)) {
                result = BigInteger.ONE.negate(); // -1
            }

            // Shift and OR low 7
            result = result.shiftLeft(7)
                .or(BigInteger.valueOf(b & LAST_SEVEN_BITS));

            if ((b & EIGHTH_BIT) == 0) {
                return result;
            }
        }

        throw new IllegalArgumentException("Truncated var-length BigInteger: no terminating byte within field size");
    }
}
