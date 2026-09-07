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

/**
 * BCD = Binary Encoded Decimal (A decimal number is represented by a sequence of 4-bit hexadecimal values from 0-9)
 * <a href="https://www.elektronik-kompendium.de/sites/dig/1010311.htm">...</a>
 * <p>
 * IMPORTANT - bit alignment: the byte-based buffers exchange partial (non byte multiple) fields with the
 * encodings RIGHT-aligned in the returned/accepted {@code byte[]}. {@code ReadBufferByteBased.readBits}
 * starts filling the result at bit {@code (8 - (numBits % 8)) % 8} and {@code WriteBufferByteBased.writeBits}
 * consumes the value starting at that same bit, so a 12 bit field carrying the digits {@code 123} is the
 * array {@code {0x01, 0x23}} - not {@code {0x12, 0x30}}. {@code EncodingUnsignedBinary} follows the same
 * convention (its {@code shift} is deliberately 0). This class therefore skips the leading padding nibbles
 * instead of packing/unpacking digits from the high nibble of byte 0: for a field whose digit count is odd
 * the first digit lives in the LOW nibble of the first byte. Getting this wrong is invisible whenever
 * {@code numBits % 8 == 0} (offset 0) and corrupts every odd-digit-count field otherwise - which is exactly
 * how the S7 {@code DateAndTime} 12 bit millisecond and 4 bit day-of-week fields used to break.
 * <p>
 * CAVEAT - byte order: the buffers put {@code byteOrder.process(bytes)} between this encoding and the
 * bit stream. Reversing a right-aligned {@code byte[]} moves the leading padding nibble into the middle
 * of the field, so for a field whose digit count is odd the contract above only holds under a
 * big-endian byte order (the default, and the only one any mspec currently pairs with BCD): under
 * {@code LITTLE_ENDIAN} a 12 bit write of {@code 123} reads back as {@code 103}. Reversing whole bytes
 * is simply not a meaningful operation on a nibble stream that does not start on a byte boundary; if a
 * protocol ever needs byte-swapped BCD it has to swap complete BCD bytes, not the encoded field.
 * <p>
 * Rejecting that combination outright would be preferable to corrupting the value silently, but the
 * {@code decode*}/{@code encode*} contract passes no byte order down here, so the guard would have to
 * live in the buffers next to their {@code byteOrder.process} call - a behaviour change (a hard failure
 * where callers see none today) and hence deliberately not made as part of the alignment fix.
 */
public class EncodingBCD extends BaseEncodingDefault {

    public static final String NAME = "BCD";

    private static final WithOption OPTION = WithOption.WithEncoding(NAME);

    public static WithOption optionEncodingBCD() {
        return OPTION;
    }

    @Override
    public String getName() {
        return NAME;
    }

    /**
     * Number of bytes a {@code numBits} wide field occupies in the buffer representation.
     */
    private static int numBytesFor(int numBits) {
        return (numBits + 7) / 8;
    }

    /**
     * Index of the nibble holding the most significant BCD digit. Non-zero (namely 1) exactly when the
     * field has an odd number of digits, in which case the leading nibble of the first byte is padding.
     */
    private static int firstNibbleIndex(int numBits) {
        return (numBytesFor(numBits) * 2) - (numBits / 4);
    }

    private static void requireMultipleOfFour(int numBits) {
        if (numBits % 4 != 0) {
            throw new IllegalArgumentException("numBits must be a multiple of 4");
        }
    }

    /**
     * Same guard (and same message) {@code EncodingUnsignedBinary} applies, so a short array reports
     * what is missing instead of throwing a raw {@link ArrayIndexOutOfBoundsException} out of the
     * nibble unpacking. A {@code null} array still fails with a {@link NullPointerException}.
     */
    private static void requireEnoughBytes(int numBits, byte[] bytes) {
        int numBytes = numBytesFor(numBits);
        if (bytes.length < numBytes) {
            throw new IllegalArgumentException("Expected at least " + numBytes + " bytes");
        }
    }

    /**
     * Reads the BCD digit stored in the given nibble, validating it is in the range 0-9.
     */
    private static int readDigit(byte[] bytes, int nibbleIndex) {
        int currentByte = bytes[nibbleIndex / 2];
        int digit = ((nibbleIndex % 2) == 0) ? ((currentByte >> 4) & 0x0F) : (currentByte & 0x0F);
        if (digit > 9) {
            throw new IllegalArgumentException("Invalid BCD digit: " + digit);
        }
        return digit;
    }

    /**
     * Packs the given digits (most significant first) right-aligned into a {@code numBits} wide field.
     */
    private static byte[] packDigits(int numBits, byte[] digits) {
        byte[] result = new byte[numBytesFor(numBits)];
        int nibbleIndex = firstNibbleIndex(numBits);
        for (byte digit : digits) {
            if ((nibbleIndex % 2) == 0) {
                result[nibbleIndex / 2] |= (byte) (digit << 4);
            } else {
                result[nibbleIndex / 2] |= digit;
            }
            nibbleIndex++;
        }
        return result;
    }

    @Override
    public byte[] encodeByte(int numBits, byte value) {
        return encodeInt(numBits, value & 0xFF);
    }

    @Override
    public byte decodeByte(int numBits, byte[] bytes) {
        int intValue = decodeInt(numBits, bytes);

        if (intValue > 255) {
            throw new IllegalArgumentException("Decoded value too large for byte: " + intValue);
        }

        return (byte) intValue;
    }

    @Override
    public byte[] encodeShort(int numBits, short value) {
        return encodeInt(numBits, value & 0xFFFF);
    }

    @Override
    public short decodeShort(int numBits, byte[] bytes) {
        int intValue = decodeInt(numBits, bytes);

        if (intValue > 65535) {
            throw new IllegalArgumentException("Decoded value too large for short: " + intValue);
        }

        return (short) intValue;
    }

    @Override
    public byte[] encodeInt(int numBits, int value) {
        requireMultipleOfFour(numBits);
        if (value < 0) {
            throw new IllegalArgumentException("BCD encoding only supports non-negative integers");
        }

        // Add value range validation
        int numDigits = numBits / 4;
        int maxValue = (int) Math.pow(10, numDigits);
        if (value >= maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d BCD digits (max value: %d)",
                    value, numDigits, maxValue - 1));
        }

        // Extract decimal digits from least significant to most
        byte[] digits = new byte[numDigits];
        for (int i = numDigits - 1; i >= 0; i--) {
            digits[i] = (byte) (value % 10);
            value /= 10;
        }

        return packDigits(numBits, digits);
    }

    @Override
    public int decodeInt(int numBits, byte[] bytes) {
        requireMultipleOfFour(numBits);
        requireEnoughBytes(numBits, bytes);

        int numDigits = numBits / 4;
        int firstNibble = firstNibbleIndex(numBits);
        int value = 0;

        for (int i = 0; i < numDigits; i++) {
            value = value * 10 + readDigit(bytes, firstNibble + i);
        }
        return value;
    }

    @Override
    public byte[] encodeLong(int numBits, long value) {
        requireMultipleOfFour(numBits);
        if (value < 0) {
            throw new IllegalArgumentException("BCD encoding only supports non-negative integers");
        }

        // Add value range validation
        int numDigits = numBits / 4;
        long maxValue = (long) Math.pow(10, numDigits);
        if (value >= maxValue) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d BCD digits (max value: %d)",
                    value, numDigits, maxValue - 1));
        }

        // Extract decimal digits from least significant to most
        byte[] digits = new byte[numDigits];
        for (int i = numDigits - 1; i >= 0; i--) {
            digits[i] = (byte) (value % 10);
            value /= 10;
        }

        return packDigits(numBits, digits);
    }

    @Override
    public long decodeLong(int numBits, byte[] bytes) {
        requireMultipleOfFour(numBits);
        requireEnoughBytes(numBits, bytes);

        int numDigits = numBits / 4;
        int firstNibble = firstNibbleIndex(numBits);
        long value = 0;

        for (int i = 0; i < numDigits; i++) {
            value = value * 10 + readDigit(bytes, firstNibble + i);
        }
        return value;
    }


    @Override
    public byte[] encodeBigInteger(int numBits, BigInteger value) {
        requireMultipleOfFour(numBits);
        if (value.compareTo(BigInteger.ZERO) < 0) {
            throw new IllegalArgumentException("BCD encoding only supports non-negative integers");
        }

        // Add value range validation
        int numDigits = numBits / 4;
        BigInteger maxValue = BigInteger.TEN.pow(numDigits);
        if (value.compareTo(maxValue) >= 0) {
            throw new IllegalArgumentException(
                String.format("Value %d cannot be encoded in %d BCD digits (max value: %d)",
                    value, numDigits, maxValue.subtract(BigInteger.ONE)));
        }

        // Extract decimal digits from least significant to most
        byte[] digits = new byte[numDigits];
        for (int i = numDigits - 1; i >= 0; i--) {
            digits[i] = value.mod(BigInteger.TEN).byteValue();
            value = value.divide(BigInteger.TEN);
        }

        return packDigits(numBits, digits);
    }

    @Override
    public BigInteger decodeBigInteger(int numBits, byte[] bytes) {
        requireMultipleOfFour(numBits);
        requireEnoughBytes(numBits, bytes);

        int numDigits = numBits / 4;
        int firstNibble = firstNibbleIndex(numBits);
        BigInteger value = BigInteger.ZERO;

        for (int i = 0; i < numDigits; i++) {
            value = value.multiply(BigInteger.TEN).add(BigInteger.valueOf(readDigit(bytes, firstNibble + i)));
        }
        return value;
    }

}
