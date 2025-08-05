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

import com.github.jinahya.bit.io.BufferByteOutput;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.spi.generation.io.MyDefaultBitOutput;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

public class WriteBufferByteBased implements WriteBuffer, BufferCommons {

    private final ByteBuffer bb;
    private final MyDefaultBitOutput bo;
    private ByteOrder byteOrder;

    public WriteBufferByteBased(int size) {
        this(size, ByteOrder.BIG_ENDIAN);
    }

    public WriteBufferByteBased(int size, ByteOrder byteOrder) {
        bb = ByteBuffer.allocate(size);
        BufferByteOutput<ByteBuffer> bbo = new BufferByteOutput<>(bb);
        bo = new MyDefaultBitOutput(bbo);
        this.byteOrder = byteOrder;
    }

    public ByteOrder getByteOrder() {
        return byteOrder;
    }

    public void setByteOrder(ByteOrder byteOrder) {
        this.byteOrder = byteOrder;
    }

    public void setPos(int position) {
        bb.position(position);
    }

    public byte[] getBytes() {
        return ArrayUtils.subarray(bb.array(), 0, getPos());
    }

    public byte[] getBytes(int start, int end) {
        return ArrayUtils.subarray(bb.array(), start, end);
    }

    @Override
    public int getPos() {
        return (int) bo.getPos();
    }

    @Override
    public void pushContext(String logicalName, WithWriterArgs... writerArgs) {
        // byte buffer need no context handling
    }

    @Override
    public void writeBit(String logicalName, boolean value, WithWriterArgs... writerArgs) throws SerializationException {
        try {
            bo.writeBoolean(value);
        } catch (Exception e) {
            throw new SerializationException("Error writing bit", e);
        }
    }

    @Override
    public void writeByte(String logicalName, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        writeSignedByte(logicalName, 8, value, writerArgs);
    }

    @Override
    public void writeByteArray(String logicalName, byte[] bytes, WithWriterArgs... writerArgs) throws SerializationException {
        for (byte aByte : bytes) {
            writeSignedByte(logicalName, 8, aByte, writerArgs);
        }
    }

    @Override
    public void writeUnsignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength <= 0) {
            throw new SerializationException("unsigned byte must contain at least 1 bit");
        }
        if (bitLength > 8) {
            throw new SerializationException("unsigned byte can only contain max 8 bits");
        }
        try {
            String encoding = extractEncoding(writerArgs).orElse("default");
            switch (encoding) {
                case "BCD":
                    if(bitLength % 4 != 0) {
                        throw new SerializationException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    if((value < 0) || (value > 9)) {
                        throw new SerializationException("'BCD' encoded value must be only one hexadecimal digit long");
                    }
                    bo.writeByte(true, bitLength, value);
                    break;
                case "default":
                    bo.writeByte(true, bitLength, value);
                    break;
                default:
                    throw new ParseException("unsupported encoding '" + encoding + "'");
            }
        } catch (Exception e) {
            throw new SerializationException("Error writing unsigned byte", e);
        }
    }

    @Override
    public void writeUnsignedShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength <= 0) {
            throw new SerializationException("unsigned short must contain at least 1 bit");
        }
        if (bitLength > 16) {
            throw new SerializationException("unsigned short can only contain max 16 bits");
        }
        try {
            String encoding = extractEncoding(writerArgs).orElse("default");
            switch (encoding) {
                case "ASCII": {
                    // AsciiUint can only decode values that have a multiple of 8 length.
                    if (bitLength % 8 != 0) {
                        throw new SerializationException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }
                    int charLen = bitLength / 8;
                    short maxValue = (short) (Math.pow(10, charLen) - 1);
                    if (value > maxValue) {
                        throw new SerializationException("Provided value of " + value + " exceeds the max value of " + maxValue);
                    }
                    String stringValue = String.format("%0" + charLen + "d", value);
                    for (byte curByte : stringValue.getBytes(StandardCharsets.US_ASCII)) {
                        bo.writeByte(false, 8, curByte);
                    }
                    break;
                }
                case "BCD": {
                    if (bitLength % 4 != 0) {
                        throw new SerializationException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    int numDigits = bitLength / 4;
                    short maxValue = (short) (Math.pow(10, numDigits) - 1);
                    if (value > maxValue) {
                        throw new SerializationException("Provided value of " + value + " exceeds the max value of " + maxValue);
                    }
                    // Write all but the last digit, by dividing the number
                    // by powers of 10 and writing the last number.
                    for(int i = numDigits - 1; i >= 0; i--) {
                        short divisor = (short) Math.pow(10, i);
                        bo.writeByte(false, 4, (byte) ((value / divisor) % 10));
                    }
                    break;
                }
                case "default":
                    bo.writeShort(true, bitLength, value);
                    break;
                default:
                    throw new SerializationException("unsupported encoding '" + encoding + "'");
            }
        } catch (Exception e) {
            throw new SerializationException("Error writing unsigned short", e);
        }
    }

    @Override
    public void writeUnsignedInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength <= 0) {
            throw new SerializationException("unsigned int must contain at least 1 bit");
        }
        if (bitLength > 32) {
            throw new SerializationException("unsigned int can only contain max 32 bits");
        }
        try {
            String encoding = extractEncoding(writerArgs).orElse("default");
            switch (encoding) {
                case "ASCII": {
                    // AsciiUint can only decode values that have a multiple of 8 length.
                    if (bitLength % 8 != 0) {
                        throw new SerializationException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }
                    int charLen = bitLength / 8;
                    int maxValue = (int) (Math.pow(10, charLen) - 1);
                    if (value > maxValue) {
                        throw new SerializationException("Provided value of " + value + " exceeds the max value of " + maxValue);
                    }
                    String stringValue = String.format("%0" + charLen + "d", value);
                    for (byte curByte : stringValue.getBytes(StandardCharsets.US_ASCII)) {
                        bo.writeByte(false, 8, curByte);
                    }
                    break;
                }
                case "BCD": {
                    if (bitLength % 4 != 0) {
                        throw new SerializationException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    int numDigits = bitLength / 4;
                    int maxValue = (int) (Math.pow(10, numDigits) - 1);
                    if (value > maxValue) {
                        throw new SerializationException("Provided value of " + value + " exceeds the max value of " + maxValue);
                    }
                    // Write all but the last digit, by dividing the number
                    // by powers of 10 and writing the last number.
                    for(int i = numDigits - 1; i >= 0; i--) {
                        int divisor = (int) Math.pow(10, i);
                        bo.writeByte(false, 4, (byte) ((value / divisor) % 10));
                    }
                    break;
                }
                case "default":
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        value = Integer.reverseBytes(value) >> (32 - bitLength);
                    }
                    bo.writeInt(true, bitLength, value);
                    break;
                default:
                    throw new SerializationException("unsupported encoding '" + encoding + "'");
            }
        } catch (Exception e) {
            throw new SerializationException("Error writing unsigned int", e);
        }
    }

    @Override
    public void writeUnsignedLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength <= 0) {
            throw new SerializationException("unsigned long must contain at least 1 bit");
        }
        if (bitLength > 63) {
            throw new SerializationException("unsigned long can only contain max 63 bits");
        }
        try {
            String encoding = extractEncoding(writerArgs).orElse("default");
            switch (encoding) {
                case "ASCII": {
                    // AsciiUint can only decode values that have a multiple of 8 length.
                    if (bitLength % 8 != 0) {
                        throw new SerializationException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }
                    int charLen = bitLength / 8;
                    long maxValue = (long) (Math.pow(10, charLen) - 1);
                    if (value > maxValue) {
                        throw new SerializationException("Provided value of " + value + " exceeds the max value of " + maxValue);
                    }
                    String stringValue = String.format("%0" + charLen + "d", value);
                    for (byte curByte : stringValue.getBytes(StandardCharsets.US_ASCII)) {
                        bo.writeByte(false, 8, curByte);
                    }
                    break;
                }
                case "BCD": {
                    if (bitLength % 4 != 0) {
                        throw new ParseException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    int numDigits = bitLength / 4;
                    long maxValue = (long) (Math.pow(10, numDigits) - 1);
                    if (value > maxValue) {
                        throw new SerializationException("Provided value of " + value + " exceeds the max value of " + maxValue);
                    }
                    // Write all but the last digit, by dividing the number
                    // by powers of 10 and writing the last number.
                    for(int i = numDigits - 1; i >= 0; i--) {
                        long divisor = (long) Math.pow(10, i);
                        bo.writeByte(false, 4, (byte) ((value / divisor) % 10));
                    }
                    break;
                }
                // It seems that normally var-length unsigned integers would be encoded little-endian.
                // However, for S7CommPlus we have a big-endian variant. If we ever encounter a LE
                // Protocol, we will need to update this into BE and LE variants.
                // https://en.wikipedia.org/wiki/Variable-length_quantity
                case "VARUDINT": {
                    // Check that the provided value fits in the allowed bit length.
                    if (value < 0) {
                        throw new SerializationException("Provided value of " + value + " exceeds the min value of 0");
                    }
                    if (value > 0xFFFFFF7FL) {
                        throw new SerializationException("Provided value of " + value + " exceeds the max value of " + 0xFFFFFF7FL);
                    }
                    // Determine the number of 7-bit groups (bytes) required.
                    int numBytes = 0;
                    long temp = value;
                    do {
                        numBytes++;
                        temp >>>= 7;
                    } while (temp != 0);

                    // Write each 7-bit group starting from the most significant.
                    for (int i = numBytes - 1; i >= 0; i--) {
                        int shift = i * 7;
                        int b = (int) ((value >> shift) & 0x7F);
                        // Set the continuation bit for all but the last (least significant) group.
                        if (i > 0) {
                            b |= 0x80;
                        }
                        bo.writeByte(false, 8, (byte) b);
                    }
                    break;
                }
                case "default":
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        value = Long.reverseBytes(value) >> 32;
                    }
                    bo.writeLong(true, bitLength, value);
                    break;
                default:
                    throw new SerializationException("unsupported encoding '" + encoding + "'");
            }
        } catch (Exception e) {
            throw new SerializationException("Error writing unsigned long", e);
        }
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        try {
            String encoding = extractEncoding(writerArgs).orElse("default");
            switch (encoding) {
                case "ASCII":
                    // AsciiUint can only decode values that have a multiple of 8 length.
                    if (bitLength % 8 != 0) {
                        throw new SerializationException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }

                    String stringValue = value.toString();
                    if(stringValue.length() > (bitLength / 8)) {
                        throw new SerializationException("'ASCII' encoded fields must have a length that is a multiple of 8 bits long");
                    }
                    writeString(stringValue.length() * 8, stringValue);
                    return;
                case "BCD":
                    if (bitLength % 4 != 0) {
                        throw new SerializationException("'BCD' encoded fields must have a length that is a multiple of 4 bits long");
                    }
                    int numDigits = bitLength / 4;
                    for (int i = numDigits - 2; i >= 0; i = i - 2) {
                        int twoDigits = value.divide(BigInteger.valueOf(10).pow(i)).mod(BigInteger.valueOf(100)).intValue();
                        byte bcdDigits = (byte) ((twoDigits / 10) << 4 | twoDigits % 10);
                        writeByte(bcdDigits);
                    }
                    return;
                case "VARUDINT":
                    // Check that the provided value fits in the allowed bit length.
                    if (value.compareTo(BigInteger.ZERO) < 0) {
                        throw new SerializationException("Provided value of " + value + " exceeds the min value of 0");
                    }
                    /*if (value.compareTo(BigInteger.valueOf(0xFFFFFF7FL)) > 0) {
                        throw new SerializationException("Provided value of " + value + " exceeds the max value of " + 0xFFFFFF7FL);
                    }*/
                    // Determine the number of 7-bit groups (bytes) required.
                    int numBytes = 0;
                    long temp = value.longValue();
                    do {
                        numBytes++;
                        temp >>>= 7;
                    } while (temp != 0);

                    // Write each 7-bit group starting from the most significant.
                    long longValue = value.longValue();
                    for (int i = numBytes - 1; i >= 0; i--) {
                        int shift = i * 7;
                        int b = (int) ((longValue >> shift) & 0x7F);
                        // Set the continuation bit for all but the last (least significant) group.
                        if (i > 0) {
                            b |= 0x80;
                        }
                        writeByte((byte) b);
                    }
                    return;
                case "default":
                    if (bitLength == 64) {
                        if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                            if (value.compareTo(BigInteger.valueOf(Long.MAX_VALUE)) >= 0) {
                                writeLong(logicalName, 32, value.longValue(), writerArgs);
                                writeLong(logicalName, 32, value.shiftRight(32).longValue(), writerArgs);
                            } else {
                                writeLong(logicalName, bitLength, value.longValue(), writerArgs);
                            }
                        } else {
                            if (value.compareTo(BigInteger.valueOf(Long.MAX_VALUE)) >= 0) {
                                writeLong(logicalName, 32, value.shiftRight(32).longValue(), writerArgs);
                                writeLong(logicalName, 32, value.longValue(), writerArgs);
                            } else {
                                writeLong(logicalName, bitLength, value.longValue(), writerArgs);
                            }
                        }
                    } else if (bitLength < 64) {
                        writeUnsignedLong(logicalName, bitLength, value.longValue(), writerArgs);
                    } else {
                        throw new SerializationException("Unsigned Big Integer can only contain max 64 bits");
                    }
            }
        } catch (ArithmeticException e) {
            throw new SerializationException("Error writing unsigned big integer", e);
        }
    }

    @Override
    public void writeSignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength <= 0) {
            throw new SerializationException("byte must contain at least 1 bit");
        }
        if (bitLength > 8) {
            throw new SerializationException("byte can only contain max 8 bits");
        }
        try {
            bo.writeByte(false, bitLength, value);
        } catch (Exception e) {
            throw new SerializationException("Error writing signed byte", e);
        }
    }

    @Override
    public void writeShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength <= 0) {
            throw new SerializationException("short must contain at least 1 bit");
        }
        if (bitLength > 16) {
            throw new SerializationException("short can only contain max 16 bits");
        }
        try {
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                value = Short.reverseBytes(value);
            }
            bo.writeShort(false, bitLength, value);
        } catch (Exception e) {
            throw new SerializationException("Error writing signed short", e);
        }
    }

    @Override
    public void writeInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength <= 0) {
            throw new SerializationException("int must contain at least 1 bit");
        }
        if (bitLength > 32) {
            throw new SerializationException("int can only contain max 32 bits");
        }
        try {
            String encoding = extractEncoding(writerArgs).orElse("default");
            switch (encoding) {
                // https://en.wikipedia.org/wiki/Variable-length_quantity
                // The first byte of a var-length signed integer contains only 6 bits (the last 6)
                // the seventh bit more or less defines the sign (1 = negative, 0 = positive)
                // If the number fits in 6 bits, the eighth bit is not set and we're done.
                // If not the first 6 bits are output, the eighth bit
                case "VARDINT": {
                    // Find out how any bytes are needed to serialize the current value
                    boolean positive = value >= 0;
                    int numBytes = 1;
                    long tmpValue = value;
                    while (tmpValue >> 6 != (positive ? 0 : -1)) {
                        numBytes++;
                        tmpValue >>= 7;
                    }

                    // Serialise the bytes
                    for (int i = numBytes - 1; i >= 0; i--) {
                        tmpValue = value >> (7 * i) & 0x7F;
                        if (i > 0) {
                            tmpValue |= 0x80;
                        } else {
                            tmpValue &= 0x7F;
                        }
                        bo.writeShort(false, 8, (short) tmpValue);
                    }
                    break;
                }
                case "default":
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        value = Integer.reverseBytes(value);
                    }
                    bo.writeInt(false, bitLength, value);
                    break;
                default:
                    throw new SerializationException("unsupported encoding '" + encoding + "'");
            }

        } catch (Exception e) {
            throw new SerializationException("Error writing signed int", e);
        }
    }

    @Override
    public void writeLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength <= 0) {
            throw new SerializationException("long must contain at least 1 bit");
        }
        if (bitLength > 64) {
            throw new SerializationException("long can only contain max 64 bits");
        }
        try {
            String encoding = extractEncoding(writerArgs).orElse("default");
            switch (encoding) {
                case "default":
                    if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                        value = Long.reverseBytes(value);
                    }
                    bo.writeLong(false, bitLength, value);
                    break;
                default:
                    throw new SerializationException("unsupported encoding '" + encoding + "'");
            }
        } catch (Exception e) {
            throw new SerializationException("Error writing signed long", e);
        }
    }

    @Override
    public void writeBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        try {
            if (bitLength > 64) {
                throw new SerializationException("Big Integer can only contain max 64 bits");
            }
            writeLong(logicalName, bitLength, value.longValue());
        } catch (ArithmeticException e) {
            throw new SerializationException("Error writing big integer", e);
        }
    }

    @Override
    public void writeFloat(String logicalName, int bitLength, float value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength != 32) {
            throw new UnsupportedOperationException("Error writing float: Exponent and/or Mantissa non standard size");
        }
        // TODO: This assumes the default encoding of IEEE 32 bit floating point
        writeInt(logicalName, bitLength, Float.floatToRawIntBits(value));
    }

    @Override
    public void writeDouble(String logicalName, int bitLength, double value, WithWriterArgs... writerArgs) throws SerializationException {
        if (bitLength != 64) {
            throw new UnsupportedOperationException("Error writing double: Exponent and/or Mantissa non standard size");
        }
        // TODO: This assumes the default encoding of IEEE 64 bit floating point
        writeLong(logicalName, bitLength, Double.doubleToRawLongBits(value));
    }

    @Override
    public void writeBigDecimal(String logicalName, int bitLength, BigDecimal value, WithWriterArgs... writerArgs) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /*
     * When encoding strings we currently implement a sort of 0-terminated string. If the string is shorter than the
     * max bit-length, we fill it up with 0x00, which makes it 0-terminated. If it exactly fits, then there is no
     * 0-termination.
     */
    @Override
    public void writeString(String logicalName, int bitLength, String value, WithWriterArgs... writerArgs) throws SerializationException {
        byte[] bytes;
        String encoding = extractEncoding(writerArgs).orElse("UTF-8");
        encoding = encoding.replaceAll("[^a-zA-Z0-9]", "");
        switch (encoding.toUpperCase()) {
            case "ASCII": {
                bytes = value.getBytes(StandardCharsets.US_ASCII);
                break;
            }
            case "WINDOWS1252": {
                bytes = value.getBytes(Charset.forName("windows-1252"));
                break;
            }
            case "UTF8": {
                bytes = value.getBytes(StandardCharsets.UTF_8);
                break;
            }
            case "UTF16":
            case "UTF16BE": {
                bytes = value.getBytes(StandardCharsets.UTF_16);
                if(bytes.length > 2) {
                    bytes = new byte[] {
                        bytes[2], bytes[3]
                    };
                }
                break;
            }
            case "UTF16LE":
                bytes = value.getBytes(StandardCharsets.UTF_16LE);
                break;
            default:
                throw new SerializationException("Unsupported encoding: " + encoding);
        }

        int fixedByteLength = (int) Math.ceil((float) bitLength / 8.0);
        if (bitLength == 0) {
            fixedByteLength = bytes.length;
        }

        try {
            int numStringBytes = Math.min(bytes.length, fixedByteLength);
            int numZeroBytes = fixedByteLength - numStringBytes;
            // Output the string data
            for (int i = 0; i < numStringBytes; i++) {
                bo.writeByte(false, 8, bytes[i]);
            }
            // Fill up with empty bytes
            for (int i = 0; i < numZeroBytes; i++) {
                bo.writeByte(false, 8, (byte) 0x00);
            }
        } catch (Exception e) {
            throw new SerializationException("Error writing string", e);
        }
    }

    @Override
    public void popContext(String logicalName, WithWriterArgs... writerArgs) {
        // byte buffer need no context handling
    }

}
