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
 */
public class EncodingBCD extends BaseEncodingDefault {

    public static final String NAME = "BCD";

    public static WithOption optionEncodingBCD() {
        return WithOption.WithEncoding(NAME);
    }

    @Override
    public String getName() {
        return NAME;
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
        if (numBits % 4 != 0) {
            throw new IllegalArgumentException("numBits must be a multiple of 4");
        }
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

        // Pack digits into bytes (two digits per byte)
        byte[] result = new byte[(numDigits + 1) / 2];
        for (int i = 0; i < numDigits; i += 2) {
            int high = digits[i];
            int low = (i + 1 < numDigits) ? digits[i + 1] : 0;
            result[i / 2] = (byte) ((high << 4) | low);
        }

        return result;
    }

    @Override
    public int decodeInt(int numBits, byte[] bytes) {
        if (numBits % 4 != 0) {
            throw new IllegalArgumentException("numBits must be a multiple of 4");
        }

        int numDigits = numBits / 4;
        int value = 0;

        for (int i = 0; i < numDigits; i++) {
            int byteIndex = i / 2;
            boolean evenDigitNumber = (i % 2 == 0);

            int digit;
            if (evenDigitNumber) {
                digit = (bytes[byteIndex] >> 4) & 0x0F;
            } else {
                digit = bytes[byteIndex] & 0x0F;
            }

            if (digit > 9) {
                throw new IllegalArgumentException("Invalid BCD digit: " + digit);
            }

            value = value * 10 + digit;
        }
        return value;
    }

    @Override
    public byte[] encodeLong(int numBits, long value) {
        if (numBits % 4 != 0) {
            throw new IllegalArgumentException("numBits must be a multiple of 4");
        }
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

        // Pack digits into bytes (two digits per byte)
        byte[] result = new byte[(numDigits + 1) / 2];
        for (int i = 0; i < numDigits; i += 2) {
            int high = digits[i];
            int low = (i + 1 < numDigits) ? digits[i + 1] : 0;
            result[i / 2] = (byte) ((high << 4) | low);
        }

        return result;
    }

    @Override
    public long decodeLong(int numBits, byte[] bytes) {
        if (numBits % 4 != 0) {
            throw new IllegalArgumentException("numBits must be a multiple of 4");
        }

        int numDigits = numBits / 4;
        long value = 0;

        for (int i = 0; i < numDigits; i++) {
            int byteIndex = i / 2;
            boolean evenDigitNumber = (i % 2 == 0);

            int digit;
            if (evenDigitNumber) {
                digit = (bytes[byteIndex] >> 4) & 0x0F;
            } else {
                digit = bytes[byteIndex] & 0x0F;
            }

            if (digit > 9) {
                throw new IllegalArgumentException("Invalid BCD digit: " + digit);
            }

            value = value * 10 + digit;
        }
        return value;
    }


    @Override
    public byte[] encodeBigInteger(int numBits, BigInteger value) {
        if (numBits % 4 != 0) {
            throw new IllegalArgumentException("numBits must be a multiple of 4");
        }
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

        // Pack digits into bytes (two digits per byte)
        byte[] result = new byte[(numDigits + 1) / 2];
        for (int i = 0; i < numDigits; i += 2) {
            int high = digits[i];
            int low = (i + 1 < numDigits) ? digits[i + 1] : 0;
            result[i / 2] = (byte) ((high << 4) | low);
        }

        return result;
    }

    @Override
    public BigInteger decodeBigInteger(int numBits, byte[] bytes) {
        if (numBits % 4 != 0) {
            throw new IllegalArgumentException("numBits must be a multiple of 4");
        }

        int numDigits = numBits / 4;
        BigInteger value = BigInteger.ZERO;

        for (int i = 0; i < numDigits; i++) {
            int byteIndex = i / 2;
            boolean evenDigitNumber = (i % 2 == 0);

            int digit;
            if (evenDigitNumber) {
                digit = (bytes[byteIndex] >> 4) & 0x0F;
            } else {
                digit = bytes[byteIndex] & 0x0F;
            }

            if (digit > 9) {
                throw new IllegalArgumentException("Invalid BCD digit: " + digit);
            }

            value = value.multiply(BigInteger.TEN).add(BigInteger.valueOf(digit));
        }
        return value;
    }

}
