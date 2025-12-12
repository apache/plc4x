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
import java.util.Arrays;

public class EncodingUnsignedBinary extends BaseEncodingDefault {

    public static final String NAME = "unsigned-binary";

    public static WithOption optionEncodingUnsignedBinary() {
        return WithOption.WithEncoding(NAME);
    }

    @Override
    public String getName() {
        return NAME;
    }

    @Override
    public byte[] encodeByte(int numBits, byte value) {
        if (numBits <= 0 || numBits > 7) {
            throw new IllegalArgumentException("numBits must be between 1 and 7 to safely fit in an unsigned byte");
        }

        // Add value range validation
        int maxValue = (1 << numBits) - 1;
        int unsignedValue = value & 0xFF;
        if (unsignedValue > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (max value: %d)",
                    unsignedValue, numBits, maxValue));
        }

        int shift = 0;//8 - numBits;
        byte encoded = (byte) (unsignedValue << shift); // align MSB

        return new byte[]{encoded};
    }

    @Override
    public byte decodeByte(int numBits, byte[] bytes) {
        if (numBits <= 0 || numBits > 7) {
            throw new IllegalArgumentException("numBits must be between 1 and 7 to safely return an unsigned byte");
        }
        if (((numBits + 7) / 8) > bytes.length) {
            throw new IllegalArgumentException("bytes should have at least " + ((numBits + 7) / 8) + " bytes to decode a byte of " + numBits + " bits");
        }

        int shift = 0;//8 - numBits;
        int raw = (bytes[0] & 0xFF) >>> shift;
        return (byte) raw;
    }

    @Override
    public byte[] encodeShort(int numBits, short value) {
        if (numBits <= 0 || numBits > 15) {
            throw new IllegalArgumentException("numBits must be between 1 and 15 to safely fit in an unsigned short");
        }

        // Add value range validation
        int maxValue = (1 << numBits) - 1;
        int unsignedValue = value & 0xFFFF;
        if (unsignedValue > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (max value: %d)",
                    unsignedValue, numBits, maxValue));
        }

        int unsigned = value & 0xFFFF; // ensure it's treated as unsigned
        if (unsigned >= (1 << numBits)) {
            throw new IllegalArgumentException("Value " + unsigned + " doesn't fit in " + numBits + " bits");
        }

        int numBytes = (numBits + 7) / 8;
        int shift = 0;//numBytes * 8 - numBits;

        int aligned = unsigned << shift;

        byte[] result = new byte[numBytes];
        for (int i = 0; i < numBytes; i++) {
            result[i] = (byte) ((aligned >>> ((numBytes - 1 - i) * 8)) & 0xFF);
        }

        return result;
    }

    @Override
    public short decodeShort(int numBits, byte[] bytes) {
        if (numBits <= 0 || numBits > 15) {
            throw new IllegalArgumentException("numBits must be between 1 and 15 to safely return an unsigned short");
        }

        int numBytes = (numBits + 7) / 8;
        if (bytes.length < numBytes) {
            throw new IllegalArgumentException("Expected at least " + numBytes + " bytes");
        }

        int shift = 0;//numBytes * 8 - numBits;

        int raw = 0;
        for (int i = 0; i < numBytes; i++) {
            raw = (raw << 8) | (bytes[i] & 0xFF);
        }

        raw >>>= shift;

        return (short) raw; // valid since max value is < 32768
    }

    @Override
    public byte[] encodeInt(int numBits, int value) {
        if (numBits <= 0 || numBits > 31) {
            throw new IllegalArgumentException("numBits must be between 1 and 31 (max value = 2^31 - 1)");
        }

        // Add value range validation
        long maxValue = (1L << numBits) - 1;
        long unsignedValue = value & 0xFFFFFFFFL;
        if (unsignedValue > maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (max value: %d)",
                    unsignedValue, numBits, maxValue));
        }

        int numBytes = (numBits + 7) / 8;
        int shift = 0;//numBytes * 8 - numBits;

        int aligned = value << shift;

        byte[] result = new byte[numBytes];
        for (int i = 0; i < numBytes; i++) {
            result[i] = (byte) ((aligned >>> ((numBytes - 1 - i) * 8)) & 0xFF);
        }

        return result;
    }

    @Override
    public int decodeInt(int numBits, byte[] bytes) {
        if (numBits <= 0 || numBits > 31) {
            throw new IllegalArgumentException("numBits must be between 1 and 31 to fit in signed int");
        }

        int numBytes = (numBits + 7) / 8;
        if (bytes.length < numBytes) {
            throw new IllegalArgumentException("Expected at least " + numBytes + " bytes");
        }

        int shift = 0;//numBytes * 8 - numBits;

        int raw = 0;
        for (int i = 0; i < numBytes; i++) {
            raw = (raw << 8) | (bytes[i] & 0xFF);
        }

        return raw >>> shift;
    }

    @Override
    public byte[] encodeLong(int numBits, long value) {
        if (numBits <= 0 || numBits > 63) {
            throw new IllegalArgumentException("numBits must be between 1 and 63 to fit in a signed long");
        }

        // Add value range validation
        if (value < 0) {
            throw new IllegalArgumentException("Value " + value + " is negative and cannot be encoded in unsigned long");
        }
        BigInteger maxValue = BigInteger.ONE.shiftLeft(numBits).subtract(BigInteger.ONE);
        BigInteger unsignedValue = BigInteger.valueOf(value & Long.MAX_VALUE);
        if (unsignedValue.compareTo(maxValue) > 0) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d bits (max value: %s)",
                    unsignedValue, numBits, maxValue.toString()));
        }

        int numBytes = (numBits + 7) / 8;
        int shift = 0;//numBytes * 8 - numBits;

        long aligned = value << shift;

        byte[] result = new byte[numBytes];
        for (int i = 0; i < numBytes; i++) {
            result[i] = (byte) ((aligned >>> ((numBytes - 1 - i) * 8)) & 0xFF);
        }

        return result;
    }

    @Override
    public long decodeLong(int numBits, byte[] bytes) {
        if (numBits <= 0 || numBits > 63) {
            throw new IllegalArgumentException("numBits must be between 1 and 63 to fit in a signed long");
        }

        int numBytes = (numBits + 7) / 8;
        if (bytes.length < numBytes) {
            throw new IllegalArgumentException("Expected at least " + numBytes + " bytes");
        }

        int shift = 0;//numBytes * 8 - numBits;

        long raw = 0;
        for (int i = 0; i < numBytes; i++) {
            raw = (raw << 8) | (bytes[i] & 0xFFL);
        }

        return raw >>> shift;
    }

    @Override
    public byte[] encodeBigInteger(int numBits, BigInteger value) {
        if (numBits <= 0) {
            throw new IllegalArgumentException("numBits must be > 0");
        }
        if (value == null) {
            throw new NullPointerException("value must not be null");
        }

        // Add value range validation
        if (value.signum() < 0) {
            throw new IllegalArgumentException("Value must be non-negative for unsigned encoding");
        }
        if (value.bitLength() > numBits) {
            throw new IllegalArgumentException("Value does not fit in " + numBits + " bits");
        }

        int numBytes = (numBits + 7) / 8;
        int shift = 0;//numBytes * 8 - numBits;

        // Left-align bits in the byte array
        BigInteger aligned = value.shiftLeft(shift);
        byte[] raw = aligned.toByteArray();

        byte[] result = new byte[numBytes];

        // Copy from raw to result (right-aligned)
        int copyFrom = Math.max(0, raw.length - numBytes);
        int copyTo = Math.max(0, result.length - raw.length);
        System.arraycopy(raw, copyFrom, result, copyTo, raw.length - copyFrom);

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

        byte[] relevant = Arrays.copyOfRange(bytes, 0, numBytes);
        BigInteger raw = new BigInteger(1, relevant); // unsigned
        return raw.shiftRight(shift);
    }

}
