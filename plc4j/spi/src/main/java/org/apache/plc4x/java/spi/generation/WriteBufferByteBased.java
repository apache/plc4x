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
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.spi.generation.io.MyDefaultBitOutput;

import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.charset.Charset;

import static org.apache.commons.lang3.ArrayUtils.subarray;

public class WriteBufferByteBased implements WriteBuffer {

    private final ByteBuffer bb;
    private final MyDefaultBitOutput bo;
    private ByteOrder byteOrder;

    public WriteBufferByteBased(int size) {
        this(size, ByteOrder.BIG_ENDIAN);
    }

    public WriteBufferByteBased(int size, ByteOrder byteOrder) {
        bb = ByteBuffer.allocate(size);
        BufferByteOutput<?> bbo = new BufferByteOutput<>(bb);
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

    public byte[] getData() {
        return ArrayUtils.subarray(bb.array(), 0, getPos());
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
        // byte buffer need no context handling
    }

    @Override
    public void writeBit(String logicalName, boolean value, WithWriterArgs... writerArgs) throws SerializationException {
        try {
            bo.writeBoolean(value);
        } catch (IOException e) {
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
            bo.writeByte(true, bitLength, value);
        } catch (IOException e) {
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
            bo.writeShort(true, bitLength, value);
        } catch (IOException e) {
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
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                value = Integer.reverseBytes(value) >> 16;
            }
            bo.writeInt(true, bitLength, value);
        } catch (IOException e) {
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
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                value = Long.reverseBytes(value) >> 32;
            }
            bo.writeLong(true, bitLength, value);
        } catch (IOException e) {
            throw new SerializationException("Error writing unsigned long", e);
        }
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        try {
            if (bitLength == 64) {
                if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
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
                throw new SerializationException("Unsigned Big Integer can only contain max 64 bits");
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
        } catch (IOException e) {
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
        } catch (IOException e) {
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
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                value = Integer.reverseBytes(value);
            }
            bo.writeInt(false, bitLength, value);
        } catch (IOException e) {
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
            if (byteOrder == ByteOrder.LITTLE_ENDIAN) {
                value = Long.reverseBytes(value);
            }
            bo.writeLong(false, bitLength, value);
        } catch (IOException e) {
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

    @Override
    public void writeString(String logicalName, int bitLength, String encoding, String value, WithWriterArgs... writerArgs) throws SerializationException {
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
            throw new SerializationException("Error writing string", e);
        }
    }

    @Override
    public void popContext(String logicalName, WithWriterArgs... writerArgs) {
        // byte buffer need no context handling
    }

}
