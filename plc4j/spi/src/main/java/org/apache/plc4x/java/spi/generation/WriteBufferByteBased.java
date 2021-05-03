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

import com.github.jinahya.bit.io.BufferByteOutput;
import org.apache.plc4x.java.spi.generation.io.MyDefaultBitOutput;

import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.charset.Charset;

public class WriteBufferByteBased implements WriteBuffer {

    private final ByteBuffer bb;
    private final BufferByteOutput<?> bbo;
    private final MyDefaultBitOutput bo;
    private final boolean littleEndian;

    public WriteBufferByteBased(int size) {
        this(size, false);
    }

    public WriteBufferByteBased(int size, boolean littleEndian) {
        bb = ByteBuffer.allocate(size);
        bbo = new BufferByteOutput<>(bb);
        bo = new MyDefaultBitOutput(bbo);
        this.littleEndian = littleEndian;
    }

    public byte[] getData() {
        return bb.array();
    }

    @Override
    public int getPos() {
        return (int) bo.getPos();
    }

    public byte[] getBytes(int startPos, int endPos) {
        int numBytes = endPos - startPos;
        byte[] data = new byte[numBytes];
        System.arraycopy(bb.array(), startPos, data, 0, numBytes);
        return data;
    }

    @Override
    public void pushContext(String logicalName, WithWriterArgs... writerArgs) {

    }

    @Override
    public void writeBit(String logicalName, boolean value, WithWriterArgs... writerArgs) throws ParseException {
        try {
            bo.writeBoolean(value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeUnsignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("unsigned byte must contain at least 1 bit");
        }
        if (bitLength > 8) {
            throw new ParseException("unsigned byte can only contain max 8 bits");
        }
        try {
            bo.writeByte(true, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeUnsignedShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("unsigned short must contain at least 1 bit");
        }
        if (bitLength > 16) {
            throw new ParseException("unsigned short can only contain max 16 bits");
        }
        try {
            bo.writeShort(true, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeUnsignedInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("unsigned int must contain at least 1 bit");
        }
        if (bitLength > 32) {
            throw new ParseException("unsigned int can only contain max 32 bits");
        }
        try {
            if (littleEndian) {
                value = Integer.reverseBytes(value) >> 16;
            }
            bo.writeInt(true, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeUnsignedLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("unsigned long must contain at least 1 bit");
        }
        if (bitLength > 63) {
            throw new ParseException("unsigned long can only contain max 63 bits");
        }
        try {
            if (littleEndian) {
                value = Long.reverseBytes(value) >> 32;
            }
            bo.writeLong(true, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws ParseException {
        try {
            if (bitLength == 64) {
                if (littleEndian) {
                    if (value.compareTo(BigInteger.valueOf(Long.MAX_VALUE)) >= 0) {
                        writeLong(logicalName, 32, value.longValue());
                        writeLong(logicalName, 32, value.shiftRight(32).longValue());
                    } else {
                        writeLong(logicalName, bitLength, value.longValue());
                    }
                } else {
                    if (value.compareTo(BigInteger.valueOf(Long.MAX_VALUE)) >= 0) {
                        writeLong(logicalName, 32, value.shiftRight(32).longValue());
                        writeLong(logicalName, 32, value.longValue());
                    } else {
                        writeLong(logicalName, bitLength, value.longValue());
                    }
                }
            } else if (bitLength < 64) {
                writeUnsignedLong(logicalName, bitLength, value.longValue());
            } else {
                throw new ParseException("Unsigned Big Integer can only contain max 64 bits");
            }
        } catch (ArithmeticException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("byte must contain at least 1 bit");
        }
        if (bitLength > 8) {
            throw new ParseException("byte can only contain max 8 bits");
        }
        try {
            bo.writeByte(false, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("short must contain at least 1 bit");
        }
        if (bitLength > 16) {
            throw new ParseException("short can only contain max 16 bits");
        }
        try {
            if (littleEndian) {
                value = Short.reverseBytes(value);
            }
            bo.writeShort(false, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("int must contain at least 1 bit");
        }
        if (bitLength > 32) {
            throw new ParseException("int can only contain max 32 bits");
        }
        try {
            if (littleEndian) {
                value = Integer.reverseBytes(value);
            }
            bo.writeInt(false, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws ParseException {
        if (bitLength <= 0) {
            throw new ParseException("long must contain at least 1 bit");
        }
        if (bitLength > 64) {
            throw new ParseException("long can only contain max 64 bits");
        }
        try {
            if (littleEndian) {
                value = Long.reverseBytes(value);
            }
            bo.writeLong(false, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws ParseException {
        try {
            if (bitLength > 64) {
                throw new ParseException("Big Integer can only contain max 64 bits");
            }
            writeLong(logicalName, bitLength, value.longValue());
        } catch (ArithmeticException e) {
            throw new ParseException("Error reading", e);
        }
    }

    @Override
    public void writeFloat(String logicalName, float value, int bitsExponent, int bitsMantissa, WithWriterArgs... writerArgs) throws ParseException {
        if (bitsExponent != 8 || bitsMantissa != 23) {
            throw new UnsupportedOperationException("Exponent and/or Mantissa non standard size");
        }
        writeInt(logicalName, 1 + bitsExponent + bitsMantissa, Float.floatToRawIntBits(value));
    }

    @Override
    public void writeDouble(String logicalName, double value, int bitsExponent, int bitsMantissa, WithWriterArgs... writerArgs) throws ParseException {
        if (bitsExponent != 11 || bitsMantissa != 52) {
            throw new UnsupportedOperationException("Exponent and/or Mantissa non standard size");
        }
        writeLong(logicalName, 1 + bitsExponent + bitsMantissa, Double.doubleToRawLongBits(value));
    }

    @Override
    public void writeBigDecimal(String logicalName, int bitLength, BigDecimal value, WithWriterArgs... writerArgs) throws ParseException {
        throw new UnsupportedOperationException("not implemented yet");
    }

    @Override
    public void writeString(String logicalName, int bitLength, String encoding, String value, WithWriterArgs... writerArgs) throws ParseException {
        final byte[] bytes = value.getBytes(Charset.forName(encoding.replaceAll("[^a-zA-Z0-9]", "")));
        int fixedByteLength = bitLength / 8;

        if (bitLength == 0) {
            fixedByteLength = bytes.length;
        }

        try {
            for (int i = 0; i < fixedByteLength; i++) {
                if (i >= bytes.length) {
                    bo.writeByte(false, 8, (byte) 0x00);
                } else {
                    bo.writeByte(false, 8, bytes[i]);
                }
            }
        } catch (IOException e) {
            throw new ParseException("Error writing string", e);
        }
    }

    @Override
    public void popContext(String logicalName, WithWriterArgs... writerArgs) {

    }

}
