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
package org.apache.plc4x.java.spi.generation;

import com.github.jinahya.bit.io.ArrayByteInput;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.generation.io.MyDefaultBitInput;

import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.util.Objects;

public class ReadBufferByteBased implements ReadBuffer, BufferCommons {

    public static final long LAST_SEVEN_BITS = (byte) 0x7F;
    public static final long SEVENTH_BIT = (byte) 0x40;
    public static final long EIGHTH_BIT = (byte) 0x80;

    private final MyDefaultBitInput bi;
    private ByteOrder byteOrder;
    private final int totalBytes;
    private int bitPos; // This is useful when debugging so you can see the actual bit position

    public ReadBufferByteBased(byte[] input) {
        this(input, ByteOrder.BIG_ENDIAN);
    }

    public ReadBufferByteBased(byte[] input, ByteOrder byteOrder) {
        Objects.requireNonNull(input);
        Objects.requireNonNull(byteOrder);
        ArrayByteInput abi = new ArrayByteInput(input);
        this.bi = new MyDefaultBitInput(abi);
        this.byteOrder = byteOrder;
        this.totalBytes = input.length;
    }

    @Override
    public int getPos() {
        return (int) bi.getPos();
    }

    @Override
    public void reset(int pos) {
        bi.reset(pos);
        bitPos = pos*8;
    }

    public byte[] getBytes(int startPos, int endPos) {
        int numBytes = endPos - startPos;
        byte[] data = new byte[numBytes];
        System.arraycopy(bi.getDelegate().getSource(), startPos, data, 0, numBytes);
        return data;
    }

    public int getTotalBytes() {
        return totalBytes;
    }

    @Override
    public boolean hasMore(int numBits) {
        return (numBits / 8) <= (totalBytes - getPos());
    }

    @Override
    public ByteOrder getByteOrder() {
        return byteOrder;
    }

    @Override
    public void setByteOrder(ByteOrder byteOrder) {
        this.byteOrder = byteOrder;
    }

    public byte peekByte(int offset) throws ParseException {
        // Remember the old index.
        int oldIndex = bi.getDelegate().getIndex();
        try {
            // Set the delegate to the desired position.
            bi.getDelegate().index(oldIndex + offset);
            // Read the byte.
            return bi.readByte(false, 8);
        } catch (IOException e) {
            throw new ParseException("Error peeking byte", e);
        } finally {
            // Reset the delegate to the old index.
            bi.getDelegate().index(oldIndex);
        }
    }

    @Override
    public void pullContext(String logicalName, WithReaderArgs... readerArgs) {
        // byte buffer need no context handling
    }

    @Override
    public boolean readBit(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        bitPos++;
        try {
            return bi.readBoolean();
        } catch (IOException e) {
            throw new ParseException("Error reading bit", e);
        }
    }

    @Override
    public byte readByte(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=8;
        return readSignedByte(logicalName, 8, readerArgs);
    }

    @Override
    public byte[] readByteArray(String logicalName, int numberOfBytes, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=numberOfBytes*8;
        byte[] bytes = new byte[numberOfBytes];
        for (int i = 0; i < numberOfBytes; i++) {
            bytes[i] = readByte();
        }
        return bytes;
    }

    @Override
    public byte readUnsignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength <= 0) {
            throw new ParseException("unsigned byte must contain at least 1 bit");
        }
        if (bitLength > 7) {
            throw new ParseException("unsigned byte can only contain max 4 bits");
        }
        try {
            String encoding = extractEncoding(readerArgs).orElse("default");
            switch (encoding) {
                case "default":
                    return bi.readByte(true, bitLength);
                // BCD = Binary Encoded Decimal (A decimal number is represented by a sequence of 4 bit hexadecimal values from 0-9.
                // https://www.elektronik-kompendium.de/sites/dig/1010311.htm
                case "BCD":
                    if (bitLength % 4 != 0) {
                        throw new ParseException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    byte digit = bi.readByte(true, 4);
                    if ((digit < 0) || (digit > 9)) {
                        throw new ParseException("'BCD' encoded value is not a correctly encoded BCD value");
                    }
                    return digit;
                default:
                    throw new ParseException("unsupported encoding '" + encoding + "'");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned byte", e);
        }
    }

    @Override
    public short readUnsignedShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength <= 0) {
            throw new ParseException("unsigned short must contain at least 1 bit");
        }
        if (bitLength > 15) {
            throw new ParseException("unsigned short can only contain max 8 bits");
        }
        try {
            String encoding = extractEncoding(readerArgs).orElse("default");
            switch (encoding) {
                case "ASCII":
                    // AsciiUint can only decode values that have a multiple of 8 length.
                    if (bitLength % 8 != 0) {
                        throw new ParseException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }
                    int charLen = bitLength / 8;
                    byte[] stringBytes = new byte[charLen];
                    for (int i = 0; i < charLen; i++) {
                        stringBytes[i] = bi.readByte(false, 8);
                    }
                    String stringValue = new String(stringBytes, StandardCharsets.US_ASCII);
                    stringValue = stringValue.trim();
                    return Short.parseShort(stringValue);
                case "BCD":
                    if (bitLength % 4 != 0) {
                        throw new ParseException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    int numDigits = bitLength / 4;
                    short value = 0;
                    for (int i = numDigits - 1; i >= 0; i--) {
                        byte digit = bi.readByte(true, 4);
                        if ((digit < 0) || (digit > 9)) {
                            throw new ParseException("'BCD' encoded value is not a correctly encoded BCD value");
                        }
                        // Shift the current digit to the required position and add it to the rest.
                        value += (short) (digit * Math.pow(10, i));
                    }
                    return value;
                case "default":
                    // No need to flip here as we're only reading one byte.
                    return bi.readShort(true, bitLength);
                default:
                    throw new ParseException("unsupported encoding '" + encoding + "'");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned short", e);
        }
    }

    @Override
    public int readUnsignedInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength <= 0) {
            throw new ParseException("unsigned int must contain at least 1 bit");
        }
        if (bitLength > 31) {
            throw new ParseException("unsigned int can only contain max 16 bits");
        }
        try {
            String encoding = extractEncoding(readerArgs).orElse("default");
            switch (encoding) {
                case "ASCII":
                    // AsciiUint can only decode values that have a multiple of 8 length.
                    if (bitLength % 8 != 0) {
                        throw new ParseException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }
                    int charLen = bitLength / 8;
                    byte[] stringBytes = new byte[charLen];
                    for (int i = 0; i < charLen; i++) {
                        stringBytes[i] = bi.readByte(false, 8);
                    }
                    String stringValue = new String(stringBytes, StandardCharsets.US_ASCII);
                    stringValue = stringValue.trim();
                    return Integer.parseInt(stringValue);
                case "BCD":
                    if (bitLength % 4 != 0) {
                        throw new ParseException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    int numDigits = bitLength / 4;
                    int value = 0;
                    for (int i = numDigits - 1; i >= 0; i--) {
                        byte digit = bi.readByte(true, 4);
                        if ((digit < 0) || (digit > 9)) {
                            throw new ParseException("'BCD' encoded value is not a correctly encoded BCD value");
                        }
                        // Shift the current digit to the required position and add it to the rest.
                        value += (int) (digit * Math.pow(10, i));
                    }
                    return value;
                case "default":
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        final int longValue = bi.readInt(true, bitLength);
                        return Integer.reverseBytes(longValue) >>> (32 - bitLength);
                    }
                    return bi.readInt(true, bitLength);
                default:
                    throw new ParseException("unsupported encoding '" + encoding + "'");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned int", e);
        }
    }

    @Override
    public long readUnsignedLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength <= 0) {
            throw new ParseException("unsigned long must contain at least 1 bit");
        }
        if (bitLength > 63) {
            throw new ParseException("unsigned long can only contain max 32 bits");
        }
        try {
            String encoding = extractEncoding(readerArgs).orElse("default");
            switch (encoding) {
                case "ASCII":
                    // AsciiUint can only decode values that have a multiple of 8 length.
                    if (bitLength % 8 != 0) {
                        throw new ParseException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }
                    int charLen = bitLength / 8;
                    byte[] stringBytes = new byte[charLen];
                    for (int i = 0; i < charLen; i++) {
                        stringBytes[i] = bi.readByte(false, 8);
                    }
                    String stringValue = new String(stringBytes, StandardCharsets.US_ASCII);
                    stringValue = stringValue.trim();
                    return Long.parseLong(stringValue);
                case "BCD":
                    if (bitLength % 4 != 0) {
                        throw new ParseException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    int numDigits = bitLength / 4;
                    long value = 0;
                    for (int i = numDigits - 1; i >= 0; i--) {
                        byte digit = bi.readByte(true, 4);
                        if ((digit < 0) || (digit > 9)) {
                            throw new ParseException("'BCD' encoded value is not a correctly encoded BCD value");
                        }
                        // Shift the current digit to the required position and add it to the rest.
                        value += (long) (digit * Math.pow(10, i));
                    }
                    return value;
                case "VARUDINT": {
                    long result = 0;
                    int shift = 0;
                    for (int i = 0; i < 4; i++) {
                        short b = bi.readShort(true, 8);
                        // Add the lower 7 bits of b, shifted appropriately.
                        result = result << shift;
                        result |= ((long) b & 0x0000007F);
                        // If the most significant bit is 0, this is the last byte.
                        if ((b & 0x80) == 0) {
                            break;
                        }
                        shift = 7;
                        // Ensure we do not exceed the maximum allowed bit length.
                        if (shift >= bitLength) {
                            throw new ParseException("var-length-uint exceeds allowed bit length " + bitLength);
                        }
                    }
                    return result;
                }
                case "default":
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        final long longValue = bi.readLong(true, bitLength);
                        return Long.reverseBytes(longValue) >>> 32;
                    }
                    return bi.readLong(true, bitLength);
                default:
                    throw new ParseException("unsupported encoding '" + encoding + "'");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned long", e);
        }
    }

    @Override
    public BigInteger readUnsignedBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        //Support specific case where value less than 64 bits and big endian.
        if (bitLength <= 0) {
            throw new ParseException("unsigned long must contain at least 1 bit");
        }
        try {
            String encoding = extractEncoding(readerArgs).orElse("default");
            switch (encoding) {
                case "ASCII":
                    // AsciiUint can only decode values that have a multiple of 8 length.
                    if (bitLength % 8 != 0) {
                        throw new ParseException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }
                    int charLen = bitLength / 8;
                    byte[] stringBytes = new byte[charLen];
                    for (int i = 0; i < charLen; i++) {
                        stringBytes[i] = bi.readByte(false, 8);
                    }
                    String stringValue = new String(stringBytes, StandardCharsets.US_ASCII);
                    stringValue = stringValue.trim();
                    return new BigInteger(stringValue);
                case "BCD":
                    if (bitLength % 4 != 0) {
                        throw new ParseException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    int numDigits = bitLength / 4;
                    BigInteger value = BigInteger.ZERO;
                    for (int i = numDigits - 1; i >= 0; i--) {
                        byte digit = bi.readByte(true, 4);
                        if ((digit < 0) || (digit > 9)) {
                            throw new ParseException("'BCD' encoded value is not a correctly encoded BCD value");
                        }
                        // Shift the current digit to the required position and add it to the rest.
                        value = value.add(BigInteger.valueOf(digit).multiply(BigInteger.valueOf(10).pow(i)));
                    }
                    return value;
                case "VARUDINT": {
                    long result = 0;
                    for (int i = 0; i < 16; i++) {
                        short b = bi.readShort(true, 8);

                        // if this is the first byte, and it's negative (the 7th bit is true)
                        // initialize the result with a value where all bits are 1
                        if((i == 0) && ((b & SEVENTH_BIT) != 0)) {
                            result = -1;
                        }

                        // Add the lower 7 bits of b, shifted appropriately.
                        result = result << 7;
                        result |= (int) (b & LAST_SEVEN_BITS);
                        // If the most significant bit is 0, this is the last byte.
                        if ((b & EIGHTH_BIT) == 0) {
                            break;
                        }
                    }
                    return BigInteger.valueOf(result);
                }
                case "default":
                    // Read as signed value
                    long val = bi.readLong(false, bitLength);
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        val = Long.reverseBytes(val);
                    }
                    if (val >= 0) {
                        return BigInteger.valueOf(val);
                    } else {
                        BigInteger constant = BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.valueOf(2)).add(BigInteger.valueOf(2));
                        return BigInteger.valueOf(val).add(constant);
                    }
                default:
                    throw new ParseException("unsupported encoding '" + encoding + "'");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned big integer", e);
        }
    }

    @Override
    public byte readSignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength <= 0) {
            throw new ParseException("byte must contain at least 1 bit");
        }
        if (bitLength > 8) {
            throw new ParseException("byte can only contain max 8 bits");
        }
        try {
            return bi.readByte(false, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading signed byte", e);
        }
    }

    @Override
    public short readShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength <= 0) {
            throw new ParseException("short must contain at least 1 bit");
        }
        if (bitLength > 16) {
            throw new ParseException("short can only contain max 16 bits");
        }
        try {
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                return Short.reverseBytes(bi.readShort(false, bitLength));
            }
            return bi.readShort(false, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading signed short", e);
        }
    }

    @Override
    public int readInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength <= 0) {
            throw new ParseException("int must contain at least 1 bit");
        }
        if (bitLength > 32) {
            throw new ParseException("int can only contain max 32 bits");
        }
        try {
            String encoding = extractEncoding(readerArgs).orElse("default");
            switch (encoding) {
                case "VARDINT": {
                    int result = 0;
                    for (int i = 0; i < 5; i++) {
                        short b = bi.readShort(true, 8);

                        // if this is the first byte, and it's negative (the 7th bit is true)
                        // initialize the result with a value where all bits are 1
                        if((i == 0) && ((b & SEVENTH_BIT) != 0)) {
                            result = -1;
                        }

                        // Add the lower 7 bits of b, shifted appropriately.
                        result = result << 7;
                        result |= (int) (b & LAST_SEVEN_BITS);
                        // If the most significant bit is 0, this is the last byte.
                        if ((b & EIGHTH_BIT) == 0) {
                            break;
                        }
                    }
                    return result;
                }
                case "default":
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        return Integer.reverseBytes(bi.readInt(false, bitLength));
                    }
                    return bi.readInt(false, bitLength);

                default:
                    throw new ParseException("unsupported encoding '" + encoding + "'");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading signed int", e);
        }
    }

    @Override
    public long readLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength <= 0) {
            throw new ParseException("long must contain at least 1 bit");
        }
        if (bitLength > 64) {
            throw new ParseException("long can only contain max 64 bits");
        }
        try {
            String encoding = extractEncoding(readerArgs).orElse("default");
            switch (encoding) {
                case "default":
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        return Long.reverseBytes(bi.readLong(false, bitLength));
                    }
                    return bi.readLong(false, bitLength);

                default:
                    throw new ParseException("unsupported encoding '" + encoding + "'");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading signed long", e);
        }
    }

    @Override
    public BigInteger readBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        throw new UnsupportedOperationException("not implemented yet");
    }

    @Override
    public float readFloat(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        String encoding = extractEncoding(readerArgs).orElse("UTF-8");
        try {
            if (bitLength == 16) {
                if ("KNXFloat".equals(encoding)) {
                    return readKnxFloat16();
                } else {
                    return readFloat16();
                }
            } else if (bitLength == 32) {
                return readFloat32(logicalName);
            } else {
                throw new UnsupportedOperationException("unsupported bit length (only 16 and 32 supported)");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading float", e);
        }
    }

    private float readKnxFloat16() throws IOException {
        // NOTE: KNX uses 4 bits as exponent and 11 as fraction
        final boolean sign = bi.readBoolean();
        final byte exponent = bi.readByte(true, 4);
        short fraction = bi.readShort(true, 11);
        if (sign) {
            fraction = (short) (fraction | 0xF800);
        }
        return (float) (0.01 * fraction * Math.pow(2, exponent));
    }

    private float readFloat16() throws IOException {
        // https://en.wikipedia.org/wiki/Half-precision_floating-point_format
        final boolean sign = bi.readBoolean();
        final byte exponent = bi.readByte(true, 5);
        final short fraction = bi.readShort(true, 10);
        final int signMultiplication = sign ? 1 : -1;
        if ((exponent >= 1) && (exponent <= 30)) {
            return signMultiplication * (2 ^ (exponent - 15)) * (1 + (fraction / 10f));
        }
        if (exponent == 0) {
            if (fraction == 0) {
                return 0.0f;
            } else {
                return signMultiplication * (2 ^ (-14)) * (fraction / 10f);
            }
        }
        if (exponent == 31) {
            if (fraction == 0) {
                return sign ? Float.POSITIVE_INFINITY : Float.NEGATIVE_INFINITY;
            } else {
                return Float.NaN;
            }
        }
        throw new NumberFormatException();
    }

    private float readFloat32(String logicalName) throws ParseException {
        bitPos+=32;
        int intValue = readInt(logicalName, 32);
        return Float.intBitsToFloat(intValue);
    }

    @Override
    public double readDouble(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        if (bitLength == 64) {
            long longValue = readLong(logicalName, 64);
            return Double.longBitsToDouble(longValue);
        } else {
            throw new UnsupportedOperationException("Error reading double: unsupported bit length (only 64 supported)");
        }
    }

    @Override
    public BigDecimal readBigDecimal(String logicalName, int bitLength, WithReaderArgs... readerArgs) {
        bitPos+=bitLength;
        throw new UnsupportedOperationException("not implemented yet");
    }

    /*
     * When encoding strings we currently implement a sort of 0-terminated string. If the string is shorter than the
     * max bit-length, we fill it up with 0x00, which makes it 0-terminated. If it exactly fits, then there is no
     * 0-termination.
     */
    @Override
    public String readString(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        bitPos+=bitLength;
        String encoding = extractEncoding(readerArgs).orElse("UTF-8");
        encoding = encoding.replaceAll("[^a-zA-Z0-9]", "");
        encoding = encoding.toUpperCase();
        switch (encoding) {
            case "ASCII":
            case "WINDOWS1252":
            case "UTF8": {
                byte[] strBytes = new byte[bitLength / 8];
                int realLength = 0;
                boolean finishedReading = false;
                for (int i = 0; (i < (bitLength / 8)) && hasMore(8); i++) {
                    try {
                        byte b = readByte(logicalName);
                        if (!disable0Termination() && (b == 0x00)) {
                            finishedReading = true;
                        } else if (!finishedReading) {
                            strBytes[i] = b;
                            realLength++;
                        }
                    } catch (Exception e) {
                        throw new PlcRuntimeException(e);
                    }
                }
                Charset charset;
                switch (encoding) {
                    case "UTF8":
                        charset = StandardCharsets.UTF_8;
                        break;
                    case "WINDOWS1252":
                        charset = Charset.forName("windows-1252");
                        break;
                    default:
                        charset = StandardCharsets.US_ASCII;
                }
                return new String(strBytes, 0, realLength, charset);
            }
            case "UTF16":
            case "UTF16LE":
            case "UTF16BE": {
                byte[] strBytes = new byte[bitLength / 8];
                int realLength = 0;
                boolean finishedReading = false;
                for (int i = 0; (i < (bitLength / 16)) && hasMore(16); i++) {
                    try {
                        byte b1 = readByte(logicalName);
                        byte b2 = readByte(logicalName);
                        if (!disable0Termination() && (b1 == 0x00) && (b2 == 0x00)) {
                            finishedReading = true;
                        } else if (!finishedReading) {
                            strBytes[(i * 2)] = b1;
                            strBytes[(i * 2) + 1] = b2;
                            realLength += 2;
                        }
                    } catch (Exception e) {
                        throw new PlcRuntimeException(e);
                    }
                }
                Charset charset;
                switch (encoding) {
                    case "UTF16LE":
                        charset = StandardCharsets.UTF_16LE;
                        break;
                    case "UTF16BE":
                        charset = StandardCharsets.UTF_16BE;
                        break;
                    default:
                        charset = StandardCharsets.UTF_16;
                }
                return new String(strBytes, 0, realLength, charset);
            }
            default:
                throw new ParseException("Unsupported encoding: " + encoding);
        }
    }


    @Override
    public void closeContext(String logicalName, WithReaderArgs... readerArgs) {
        // byte buffer need no context handling
    }

    private boolean disable0Termination() {
        return Boolean.parseBoolean(System.getProperty("disable-string-0-termination", "false"));
    }

}
