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

public class ReadBufferByteBased implements ReadBuffer {

    private final MyDefaultBitInput bi;
    private ByteOrder byteOrder;
    private final int totalBytes;

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
        try {
            return bi.readBoolean();
        } catch (IOException e) {
            throw new ParseException("Error reading bit", e);
        }
    }

    @Override
    public byte readByte(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        return readSignedByte(logicalName, 8, readerArgs);
    }

    @Override
    public byte[] readByteArray(String logicalName, int numberOfBytes, WithReaderArgs... readerArgs) throws ParseException {
        byte[] bytes = new byte[numberOfBytes];
        for (int i = 0; i < numberOfBytes; i++) {
            bytes[i] = readByte();
        }
        return bytes;
    }

    @Override
    public byte readUnsignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("unsigned byte must contain at least 1 bit");
        }
        if (bitLength > 7) {
            throw new ParseException("unsigned byte can only contain max 4 bits");
        }
        try {
            return bi.readByte(true, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned byte", e);
        }
    }

    @Override
    public short readUnsignedShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("unsigned short must contain at least 1 bit");
        }
        if (bitLength > 15) {
            throw new ParseException("unsigned short can only contain max 8 bits");
        }
        try {
            // No need to flip here as we're only reading one byte.
            return bi.readShort(true, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned short", e);
        }
    }

    @Override
    public int readUnsignedInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("unsigned int must contain at least 1 bit");
        }
        if (bitLength > 31) {
            throw new ParseException("unsigned int can only contain max 16 bits");
        }
        try {
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                int intValue = bi.readInt(true, bitLength);
                return Integer.reverseBytes(intValue) >>> 16;
            }
            return bi.readInt(true, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned int", e);
        }
    }

    @Override
    public long readUnsignedLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("unsigned long must contain at least 1 bit");
        }
        if (bitLength > 63) {
            throw new ParseException("unsigned long can only contain max 32 bits");
        }
        try {
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                final long longValue = bi.readLong(true, bitLength);
                return Long.reverseBytes(longValue) >>> 32;
            }
            return bi.readLong(true, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned long", e);
        }
    }

    @Override
    public BigInteger readUnsignedBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        //Support specific case where value less than 64 bits and big endian.
        if (bitLength <= 0) {
            throw new ParseException("unsigned long must contain at least 1 bit");
        }
        if (bitLength > 64) {
            throw new ParseException("unsigned long can only contain max 64 bits");
        }
        try {
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
        } catch (IOException e) {
            throw new ParseException("Error reading unsigned big integer", e);
        }
    }

    @Override
    public byte readSignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
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
        if (bitLength <= 0) {
            throw new ParseException("int must contain at least 1 bit");
        }
        if (bitLength > 32) {
            throw new ParseException("int can only contain max 32 bits");
        }
        try {
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                return Integer.reverseBytes(bi.readInt(false, bitLength));
            }
            return bi.readInt(false, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading signed int", e);
        }
    }

    @Override
    public long readLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("long must contain at least 1 bit");
        }
        if (bitLength > 64) {
            throw new ParseException("long can only contain max 64 bits");
        }
        try {
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                return Long.reverseBytes(bi.readLong(false, bitLength));
            }
            return bi.readLong(false, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading signed long", e);
        }
    }

    @Override
    public BigInteger readBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        throw new UnsupportedOperationException("not implemented yet");
    }

    @Override
    public float readFloat(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        try {
            if (bitLength == 16) {
                return readFloat16();
            } else if (bitLength == 32) {
                return readFloat32(logicalName);
            } else {
                throw new UnsupportedOperationException("unsupported bit length (only 16 and 32 supported)");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading float", e);
        }
    }

    private float readFloat16() throws IOException {
        // NOTE: KNX uses 4 bits as exponent and 11 as fraction
        final boolean sign = bi.readBoolean();
        final byte exponent = bi.readByte(true, 4);
        short fraction = bi.readShort(true, 11);
        // This is a 12-bit 2's complement notation ... the first bit belongs to the last 11 bits.
        // If the first bit is set, then we need to also set the upper 5 bits of the fraction part.
        if(sign) {
            fraction = (short) (fraction | 0xF800);
        }
        if ((exponent >= 1) && (exponent < 15)) {
            return (float) (0.01 * fraction * Math.pow(2, exponent));
        }
        if (exponent == 0) {
            if (fraction == 0) {
                return 0.0f;
            } else {
                return (2 ^ (-14)) * (fraction / 10f);
            }
        }
        if (exponent == 15) {
            if (fraction == 0) {
                return sign ? Float.POSITIVE_INFINITY : Float.NEGATIVE_INFINITY;
            } else {
                return Float.NaN;
            }
        }
        throw new NumberFormatException();
    }

    /*private float readFloat16() throws IOException {
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
    }*/

    private float readFloat32(String logicalName) throws ParseException {
        int intValue = readInt(logicalName, 32);
        return Float.intBitsToFloat(intValue);
    }

    @Override
    public double readDouble(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        if (bitLength == 64) {
            long longValue = readLong(logicalName, 64);
            return Double.longBitsToDouble(longValue);
        } else {
            throw new UnsupportedOperationException("Error reading double: unsupported bit length (only 64 supported)");
        }
    }

    @Override
    public BigDecimal readBigDecimal(String logicalName, int bitLength, WithReaderArgs... readerArgs) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /*
     * When encoding strings we currently implement a sort of 0-terminated string. If the string is shorter than the
     * max bit-length, we fill it up with 0x00, which makes it 0-terminated. If it exactly fits, then there is no
     * 0-termination.
     */
    @Override
    public String readString(String logicalName, int bitLength, String encoding, WithReaderArgs... readerArgs) throws ParseException {
        encoding = encoding.replaceAll("[^a-zA-Z0-9]", "");
        switch (encoding.toUpperCase()) {
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
                return new String(strBytes, StandardCharsets.UTF_8).substring(0, realLength);
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
                        } else if (!finishedReading){
                            strBytes[(i * 2)] = b1;
                            strBytes[(i * 2) + 1] = b2;
                            realLength++;
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
                return new String(strBytes, charset).substring(0, realLength);
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
