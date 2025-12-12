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

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferValueException;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.Charset;

public abstract class BaseStringEncoding implements EncodingDefault {

    protected abstract int getBitsPerCharacter();

    protected abstract Charset getCharset();

    @Override
    public byte[] encodeByte(int numBits, byte value) throws BufferException {
        String stringValue = Byte.toString(value);
        return encodeString(numBits, stringValue, true);
    }

    @Override
    public byte decodeByte(int numBits, byte[] bytes) throws BufferException {
        String stringValue = decodeString(numBits, bytes);
        try {
            return Byte.parseByte(stringValue.trim());
        } catch (NumberFormatException e) {
            throw new BufferValueException("Byte value cannot be parsed from string: " + stringValue, bytes);
        }
    }

    @Override
    public byte[] encodeShort(int numBits, short value) throws BufferException {
        String stringValue = Short.toString(value);
        return encodeString(numBits, stringValue, true);
    }

    @Override
    public short decodeShort(int numBits, byte[] bytes) throws BufferException {
        String stringValue = decodeString(numBits, bytes);
        try {
            return Short.parseShort(stringValue.trim());
        } catch (NumberFormatException e) {
            throw new BufferValueException("Short value cannot be parsed from string: " + stringValue, bytes);
        }
    }

    @Override
    public byte[] encodeInt(int numBits, int value) throws BufferException {
        String stringValue = Integer.toString(value);
        return encodeString(numBits, stringValue, true);
    }

    @Override
    public int decodeInt(int numBits, byte[] bytes) throws BufferException {
        String stringValue = decodeString(numBits, bytes);
        try {
            return Integer.parseInt(stringValue.trim());
        } catch (NumberFormatException e) {
            throw new BufferValueException("Integer value cannot be parsed from string: " + stringValue, bytes);
        }
    }

    @Override
    public byte[] encodeLong(int numBits, long value) throws BufferException {
        String stringValue = Long.toString(value);
        return encodeString(numBits, stringValue, true);
    }

    @Override
    public long decodeLong(int numBits, byte[] bytes) throws BufferException {
        String stringValue = decodeString(numBits, bytes);
        try {
            return Long.parseLong(stringValue.trim());
        } catch (NumberFormatException e) {
            throw new BufferValueException("Long value cannot be parsed from string: " + stringValue, bytes);
        }
    }

    @Override
    public byte[] encodeBigInteger(int numBits, BigInteger value) throws BufferException {
        String stringValue = value.toString();
        return encodeString(numBits, stringValue, true);
    }

    @Override
    public BigInteger decodeBigInteger(int numBits, byte[] bytes) throws BufferException {
        String stringValue = decodeString(numBits, bytes);
        try {
            return new BigInteger(stringValue.trim());
        } catch (NumberFormatException e) {
            throw new BufferValueException("BigInteger value cannot be parsed from string: " + stringValue, bytes);
        }
    }

    @Override
    public byte[] encodeFloat(int numBits, float value) throws BufferException {
        String stringValue = Float.toString(value);
        return encodeString(numBits, stringValue, true);
    }

    @Override
    public float decodeFloat(int numBits, byte[] bytes) throws BufferException {
        String stringValue = decodeString(numBits, bytes);
        try {
            return Float.parseFloat(stringValue.trim());
        } catch (NumberFormatException e) {
            throw new BufferValueException("Float value cannot be parsed from string: " + stringValue, bytes);
        }
    }

    @Override
    public byte[] encodeDouble(int numBits, double value) throws BufferException {
        String stringValue = Double.toString(value);
        return encodeString(numBits, stringValue, true);
    }

    @Override
    public double decodeDouble(int numBits, byte[] bytes) throws BufferException {
        String stringValue = decodeString(numBits, bytes);
        try {
            return Double.parseDouble(stringValue.trim());
        } catch (NumberFormatException e) {
            throw new BufferValueException("Double value cannot be parsed from string: " + stringValue, bytes);
        }
    }

    @Override
    public byte[] encodeBigDecimal(int numBits, BigDecimal value) throws BufferException {
        String stringValue = value.toString();
        return encodeString(numBits, stringValue, true);
    }

    @Override
    public BigDecimal decodeBigDecimal(int numBits, byte[] bytes) throws BufferException {
        String stringValue = decodeString(numBits, bytes);
        try {
            return new BigDecimal(stringValue.trim());
        } catch (NumberFormatException e) {
            throw new BufferValueException("BigDecimal value cannot be parsed from string: " + stringValue, bytes);
        }
    }

    @Override
    public byte[] encodeString(int numBits, String value) throws BufferException {
        return encodeString(numBits, value, false);
    }

    public byte[] encodeString(int numBits, String value, boolean numericValue) throws BufferException {
        if (value == null) {
            throw new BufferException("value must not be null");
        }
        if (numBits % getBitsPerCharacter() != 0) {
            throw new BufferException("numBits must be a multiple of " + getBitsPerCharacter());
        }
        int numBytes = numBits / 8;
        byte[] bytes = value.getBytes(getCharset());

        // Check if the string can be properly encoded in the given charset
        String roundTrip = new String(bytes, getCharset());
        if (!roundTrip.equals(value)) {
            throw new BufferException(
                String.format("String contains characters that cannot be encoded in %s",
                    getCharset().displayName()));
        }

        // If it's a negative value, strip the leading '-' character
        boolean negative = numericValue && value.charAt(0) == '-';
        if (negative) {
            bytes = value.substring(1).getBytes(getCharset());
        }

        if (bytes.length > numBytes) {
            throw new BufferException(
                String.format("String requires %d bits but only %d bits available",
                    bytes.length * getBitsPerCharacter(), numBits));
        }
        byte[] newBytes = new byte[numBytes];
        // Numeric values are left padded with 0s
        if (numericValue) {
            // Fill with the charset-specific encoding of a single space character
            byte[] spaceBytes = "0".getBytes(getCharset());
            for (int i = 0; i < numBytes; i += spaceBytes.length) {
                System.arraycopy(spaceBytes, 0, newBytes, i, Math.min(spaceBytes.length, numBytes - i));
            }
            // Negative values must have their negativity sign in the first char.
            if (negative) {
                byte[] minusBytes = "-".getBytes(getCharset());
                System.arraycopy(minusBytes, 0, newBytes, 0, minusBytes.length);
            }
        }
        System.arraycopy(bytes, 0, newBytes, numericValue ? numBytes - bytes.length : 0, bytes.length);
        return newBytes;
    }

    @Override
    public String decodeString(int numBits, byte[] bytes) throws BufferException {
        if (bytes == null) {
            throw new BufferException("Cannot decode null byte array");
        }
        if (numBits % getBitsPerCharacter() != 0) {
            throw new BufferException("numBits must be a multiple of " + getBitsPerCharacter());
        }
        int numBytes = numBits / 8;
        if (bytes.length > numBytes) {
            throw new BufferException("byte array is too long to be decoded using " + numBits + " bits. Byte array length is " + (bytes.length * 8) + " bits.");
        } else if (bytes.length < numBytes) {
            throw new BufferException("byte array is too short to be decoded using " + numBits + " bits. Byte array length is " + (bytes.length * 8) + " bits.");
        }
        // If the entire byte array is composed only of encoded space characters, keep it as-is (do not trim)
        byte[] spaceBytes = " ".getBytes(getCharset());
        boolean onlySpaces = bytes.length > 0 && (bytes.length % spaceBytes.length == 0);
        if (onlySpaces) {
            for (int i = 0; i < bytes.length; i += spaceBytes.length) {
                for (int j = 0; j < spaceBytes.length; j++) {
                    if (i + j >= bytes.length || bytes[i + j] != spaceBytes[j]) {
                        onlySpaces = false;
                        break;
                    }
                }
                if (!onlySpaces) break;
            }
        }

        String decoded = new String(bytes, getCharset());//.replace("\u0000", "");
        if (onlySpaces) {
            return decoded;
        }
        if(decoded.contains("\u0000")) {
            decoded = decoded.substring(0, decoded.indexOf("\u0000"));
        }
        // Remove trailing padding spaces that may have been added during encoding
        int end = decoded.length();
        while (end > 0 && decoded.charAt(end - 1) == ' ') {
            end--;
        }
        return ((end == decoded.length()) ? decoded : decoded.substring(0, end)).replace("\uFEFF", "");
    }

}
