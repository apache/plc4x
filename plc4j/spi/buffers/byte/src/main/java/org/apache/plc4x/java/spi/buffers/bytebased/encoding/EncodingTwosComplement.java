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

import java.math.BigInteger;

public class EncodingTwosComplement extends BaseEncodingDefault {

    public static final String NAME = "twos-complement";

    private static final WithOption OPTION = WithOption.WithEncoding(NAME);

    public static WithOption optionEncodingTwosComplement() {
        return OPTION;
    }

    @Override
    public String getName() {
        return NAME;
    }

    @Override
    public byte[] encodeByte(int numBits, byte value) {
        if (numBits <= 0 || numBits > 8) {
            throw new IllegalArgumentException("numBits must be between 1 and 8");
        }

        // Add value range validation
        int minValue = -(1 << (numBits - 1));
        int maxValue = (1 << (numBits - 1)) - 1;
        if (value < minValue || value > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (range: %d to %d)",
                    value, numBits, minValue, maxValue));
        }

        // Mask out the bits we want to encode
        int masked = value & ((1 << numBits) - 1);

        // Align into high bits of one byte (big-endian bit packing, left-aligned)
        int shift = 0;//8 - numBits;
        byte encoded = (byte) (masked << shift);

        return new byte[]{encoded};
    }

    @Override
    public byte decodeByte(int numBits, byte[] bytes) {
        if (numBits <= 0 || numBits > 8) {
            throw new IllegalArgumentException("numBits must be between 1 and 8");
        }
        if (bytes.length < 1) {
            throw new IllegalArgumentException("At least 1 byte is required");
        }

        int shift = 0;//8 - numBits;
        int raw = (bytes[0] & 0xFF) >>> shift;

        // Sign-extend to 8 bits if highest bit of encoded value is 1
        if ((raw & (1 << (numBits - 1))) != 0) {
            raw |= -(1 << numBits);  // fill high bits with 1s
        }

        return (byte) raw;
    }

    @Override
    public byte[] encodeShort(int numBits, short value) {
        if (numBits <= 0 || numBits > 16) {
            throw new IllegalArgumentException("numBits must be between 1 and 16");
        }

        // Add value range validation
        int minValue = -(1 << (numBits - 1));
        int maxValue = (1 << (numBits - 1)) - 1;
        if (value < minValue || value > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (range: %d to %d)",
                    value, numBits, minValue, maxValue));
        }

        int numBytes = (numBits + 7) / 8;
        int shift = 0;//numBytes * 8 - numBits;

        // Mask only the lower numBits (handle negative values via two's complement)
        int masked = value & ((1 << numBits) - 1);

        // Align to the most significant bits
        masked <<= shift;

        byte[] result = new byte[numBytes];
        for (int i = 0; i < numBytes; i++) {
            result[i] = (byte) ((masked >> ((numBytes - 1 - i) * 8)) & 0xFF);
        }

        return result;
    }

    @Override
    public short decodeShort(int numBits, byte[] bytes) {
        if (numBits <= 0 || numBits > 16) {
            throw new IllegalArgumentException("numBits must be between 1 and 16");
        }

        int numBytes = (numBits + 7) / 8;
        if (bytes.length < numBytes) {
            throw new IllegalArgumentException("Expected at least " + numBytes + " bytes");
        }

        int shift = 0;//numBytes * 8 - numBits;

        // Assemble bits from bytes
        int raw = 0;
        for (int i = 0; i < numBytes; i++) {
            raw = (raw << 8) | (bytes[i] & 0xFF);
        }

        // Right-align the value
        raw >>>= shift;

        // Sign-extend if necessary
        if ((raw & (1 << (numBits - 1))) != 0) {
            raw |= -(1 << numBits);
        }

        return (short) raw;
    }

    @Override
    public byte[] encodeInt(int numBits, int value) {
        if (numBits <= 0 || numBits > 32) {
            throw new IllegalArgumentException("numBits must be between 1 and 32");
        }

        // Add value range validation
        long minValue = -(1L << (numBits - 1));
        long maxValue = (1L << (numBits - 1)) - 1;
        if (value < minValue || value > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (range: %d to %d)",
                    value, numBits, minValue, maxValue));
        }

        int numBytes = (numBits + 7) / 8;
        int shift = 0;//numBytes * 8 - numBits;

        // Mask value to numBits bits
        byte[] result = new byte[numBytes];
        if (numBits != 32) {
            int masked = value & ((1 << numBits) - 1);
            masked <<= shift;

            for (int i = 0; i < numBytes; i++) {
                result[i] = (byte) ((masked >>> ((numBytes - 1 - i) * 8)) & 0xFF);
            }
        } else {
            result[0] = (byte) ((value >>> 24) & 0xFF);
            result[1] = (byte) ((value >>> 16) & 0xFF);
            result[2] = (byte) ((value >>> 8) & 0xFF);
            result[3] = (byte) (value & 0xFF);
        }

        return result;
    }

    @Override
    public int decodeInt(int numBits, byte[] bytes) {
        if (numBits <= 0 || numBits > 32) {
            throw new IllegalArgumentException("numBits must be between 1 and 32");
        }

        int numBytes = (numBits + 7) / 8;
        if (bytes.length < numBytes) {
            throw new IllegalArgumentException("Expected at least " + numBytes + " bytes");
        }

        int shift = 0;//numBytes * 8 - numBits;

        // Read bits from bytes
        int raw = 0;
        for (int i = 0; i < numBytes; i++) {
            raw = (raw << 8) | (bytes[i] & 0xFF);
        }

        raw >>>= shift;

        // Sign-extend if negative
        if ((raw & (1L << (numBits - 1))) != 0) {
            long longRaw = raw - (1L << numBits);
            raw = (int) longRaw;
        }

        return raw;
    }

    @Override
    public byte[] encodeLong(int numBits, long value) {
        if (numBits <= 0 || numBits > 64) {
            throw new IllegalArgumentException("numBits must be between 1 and 64");
        }

        // Add value range validation
        BigInteger bigValue = BigInteger.valueOf(value);
        BigInteger minValue = BigInteger.ONE.shiftLeft(numBits - 1).negate();
        BigInteger maxValue = BigInteger.ONE.shiftLeft(numBits - 1).subtract(BigInteger.ONE);
        if (bigValue.compareTo(minValue) < 0 || bigValue.compareTo(maxValue) > 0) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (range: %s to %s)",
                    value, numBits, minValue.toString(), maxValue.toString()));
        }

        int numBytes = (numBits + 7) / 8;
        int shift = 0;//numBytes * 8 - numBits;

        byte[] result = new byte[numBytes];

        if (numBits != 64) {
            // Mask to lowest numBits bits
            long masked = value & ((1L << numBits) - 1);
            masked <<= shift;

            for (int i = 0; i < numBytes; i++) {
                result[i] = (byte) ((masked >>> ((numBytes - 1 - i) * 8)) & 0xFF);
            }
        } else {
            result[0] = (byte) ((value >>> 56) & 0xFF);
            result[1] = (byte) ((value >>> 48) & 0xFF);
            result[2] = (byte) ((value >>> 40) & 0xFF);
            result[3] = (byte) ((value >>> 32) & 0xFF);
            result[4] = (byte) ((value >>> 24) & 0xFF);
            result[5] = (byte) ((value >>> 16) & 0xFF);
            result[6] = (byte) ((value >>> 8) & 0xFF);
            result[7] = (byte) (value & 0xFF);
        }

        return result;
    }

    @Override
    public long decodeLong(int numBits, byte[] bytes) {
        if (numBits <= 0 || numBits > 64) {
            throw new IllegalArgumentException("numBits must be between 1 and 64");
        }

        int numBytes = (numBits + 7) / 8;
        if (bytes.length < numBytes) {
            throw new IllegalArgumentException("Expected at least " + numBytes + " bytes");
        }

        int shift = 0;//numBytes * 8 - numBits;

        long raw = 0;
        for (int i = 0; i < numBytes; i++) {
            raw = (raw << 8) | (bytes[i] & 0xFF);
        }

        raw >>>= shift;

        // Sign extend if needed
        BigInteger bigValue = BigInteger.valueOf(raw);
        if (bigValue.and(BigInteger.ONE.shiftLeft(numBits - 1)).compareTo(BigInteger.ZERO) != 0) {
            bigValue = bigValue.or(BigInteger.ONE.negate().shiftLeft(numBits));
            raw = bigValue.longValue();
        }

        return raw;
    }

    @Override
    public byte[] encodeBigInteger(int numBits, BigInteger value) {
        if (numBits <= 0) {
            throw new IllegalArgumentException("numBits must be > 0");
        }

        // Add value range validation
        BigInteger minValue = BigInteger.ONE.shiftLeft(numBits - 1).negate();
        BigInteger maxValue = BigInteger.ONE.shiftLeft(numBits - 1).subtract(BigInteger.ONE);
        if (value.compareTo(minValue) < 0 || value.compareTo(maxValue) > 0) {
            throw new IllegalArgumentException(
                String.format("Value %s cannot be encoded in %d bits (range: %s to %s)",
                    value.toString(), numBits, minValue.toString(), maxValue.toString()));
        }

        int numBytes = (numBits + 7) / 8;
        int shift = 0;//numBytes * 8 - numBits;

        // Modulo 2^numBits to get two's complement representation
        BigInteger mod = BigInteger.ONE.shiftLeft(numBits);
        BigInteger normalized = value.and(mod.subtract(BigInteger.ONE));

        // Shift left to align with MSB
        normalized = normalized.shiftLeft(shift);

        byte[] tmp = normalized.toByteArray();

        // Ensure the result is exactly numBytes long (may require trimming or padding)
        byte[] result = new byte[numBytes];
        int copyFrom = Math.max(0, tmp.length - numBytes);
        int copyTo = Math.max(0, result.length - tmp.length);
        System.arraycopy(tmp, copyFrom, result, copyTo, tmp.length - copyFrom);

        return result;
    }

    @Override
    public BigInteger decodeBigInteger(int numBits, byte[] bytes) {
        if (numBits <= 0) {
            throw new IllegalArgumentException("numBits must be > 0");
        }

        int numBytes = (numBits + 7) / 8;
        if (bytes.length < numBytes) {
            throw new IllegalArgumentException("Expected at least " + numBytes + " bytes");
        }

        int shift = 0;//numBytes * 8 - numBits;

        // Extract bits from the start of the array
        BigInteger raw = new BigInteger(1, bytes).shiftRight(shift);

        // Sign extend if needed
        if (raw.testBit(numBits - 1)) {
            // Negative number: extend with 1s above numBits
            BigInteger signExtension = BigInteger.valueOf(-1).shiftLeft(numBits);
            raw = raw.or(signExtension);
        }

        return raw;
    }

}
