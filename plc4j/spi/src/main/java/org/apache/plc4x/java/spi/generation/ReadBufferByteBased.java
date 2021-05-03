/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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

public class ReadBufferByteBased implements ReadBuffer {

    private final MyDefaultBitInput bi;
    private final boolean littleEndian;
    private final long totalBytes;

    public ReadBufferByteBased(byte[] input) {
        this(input, false);
    }

    public ReadBufferByteBased(byte[] input, boolean littleEndian) {
        ArrayByteInput abi = new ArrayByteInput(input);
        this.bi = new MyDefaultBitInput(abi);
        this.littleEndian = littleEndian;
        this.totalBytes = input.length;
    }

    @Override
    public int getPos() {
        return (int) bi.getPos();
    }

    public byte[] getBytes(int startPos, int endPos) {
        int numBytes = endPos - startPos;
        byte[] data = new byte[numBytes];
        System.arraycopy(bi.getDelegate().getSource(), startPos, data, 0, numBytes);
        return data;
    }

    public long getTotalBytes() {
        return totalBytes;
    }

    @Override
    public boolean hasMore(int numBits) {
        return (numBits / 8) <= (totalBytes - getPos());
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
            throw new ParseException("Error reading", e);
        } finally {
            // Reset the delegate to the old index.
            bi.getDelegate().index(oldIndex);
        }
    }

    @Override
    public void pullContext(String logicalName, WithReaderArgs... readerArgs) {
    }

    @Override
    public boolean readBit(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        try {
            return bi.readBoolean();
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
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
            throw new ParseException("Error reading", e);
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
            throw new ParseException("Error reading", e);
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
            if (littleEndian) {
                int intValue = bi.readInt(true, bitLength);
                return Integer.reverseBytes(intValue) >>> 16;
            }
            return bi.readInt(true, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
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
            if (littleEndian) {
                final long longValue = bi.readLong(true, bitLength);
                return Long.reverseBytes(longValue) >>> 32;
            }
            return bi.readLong(true, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
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
            Long val = bi.readLong(false, bitLength);
            if (littleEndian) {
                val = Long.reverseBytes(val);
            }
            if (val >= 0) {
                return BigInteger.valueOf(val);
            } else {
                BigInteger constant = BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.valueOf(2)).add(BigInteger.valueOf(2));
                return BigInteger.valueOf(val).add(constant);
            }
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public byte readByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("byte must contain at least 1 bit");
        }
        if (bitLength > 8) {
            throw new ParseException("byte can only contain max 8 bits");
        }
        try {
            return bi.readByte(false, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
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
            if (littleEndian) {
                return Short.reverseBytes(bi.readShort(false, bitLength));
            }
            return bi.readShort(false, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
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
            if (littleEndian) {
                return Integer.reverseBytes(bi.readInt(false, bitLength));
            }
            return bi.readInt(false, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
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
            if (littleEndian) {
                return Long.reverseBytes(bi.readLong(false, bitLength));
            }
            return bi.readLong(false, bitLength);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
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
                // https://en.wikipedia.org/wiki/Half-precision_floating-point_format
                final boolean sign = bi.readBoolean();
                final byte exponent = bi.readByte(true, 5);
                final short fraction = bi.readShort(true, 10);
                if ((exponent >= 1) && (exponent <= 30)) {
                    return (sign ? 1 : -1) * (2 ^ (exponent - 15)) * (1 + (fraction / 10f));
                } else if (exponent == 0) {
                    if (fraction == 0) {
                        return 0.0f;
                    } else {
                        return (sign ? 1 : -1) * (2 ^ (-14)) * (fraction / 10f);
                    }
                } else if (exponent == 31) {
                    if (fraction == 0) {
                        return sign ? Float.POSITIVE_INFINITY : Float.NEGATIVE_INFINITY;
                    } else {
                        return Float.NaN;
                    }
                } else {
                    throw new NumberFormatException();
                }
            } else if (bitLength == 32) {
                int intValue = readInt(logicalName,32);
                return Float.intBitsToFloat(intValue);
            } else {
                throw new UnsupportedOperationException("unsupported bit length (only 16 and 32 supported)");
            }
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public double readDouble(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        if(bitLength == 64) {
            long longValue = readLong(logicalName,64);
            return Double.longBitsToDouble(longValue);
        } else {
            throw new UnsupportedOperationException("unsupported bit length (only 64 supported)");
        }
    }

    @Override
    public BigDecimal readBigDecimal(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException {
        throw new UnsupportedOperationException("not implemented yet");
    }

    @Override
    public String readString(String logicalName, int bitLength, String encoding, WithReaderArgs... readerArgs) {
        byte[] strBytes = new byte[bitLength / 8];
        for (int i = 0; (i < (bitLength / 8)) && hasMore(8); i++) {
            try {
                strBytes[i] = readByte(logicalName,8);
            } catch (Exception e) {
                throw new PlcRuntimeException(e);
            }
        }
        //replaceAll function removes and leading ' char or hypens.
        return new String(strBytes, Charset.forName(encoding.replaceAll("[^a-zA-Z0-9]","")));
    }


    @Override
    public void closeContext(String logicalName, WithReaderArgs... readerArgs) {
    }

}
