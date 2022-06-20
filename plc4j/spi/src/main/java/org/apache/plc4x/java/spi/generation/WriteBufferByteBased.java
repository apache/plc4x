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

    /**
     * @deprecated use {@link WriteBufferByteBased#getBytes()}
     */
    @Deprecated
    public byte[] getData() {
        return getBytes();
    }

    public byte[] getBytes() {
        return ArrayUtils.subarray(bb.array(), 0, getPos());
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
        if (bitLength <= 0) {
            throw new SerializationException("long must contain at least 1 bit");
        }
        if (value.compareTo(BigInteger.ZERO) < 0) {
            throw new SerializationException("value " + value + " is below 0");
        }
        int actualBitLength = value.bitLength();
        if (bitLength < actualBitLength) {
            throw new SerializationException("bit length" + actualBitLength + " exceeds supplied bit length " + bitLength);
        }
        byte[] bytes = value.toByteArray();
        int remainingBitLength = actualBitLength;
        try {
            if (byteOrder != ByteOrder.LITTLE_ENDIAN) {
                // MSB in 0
                for (int i = 0; i < bytes.length; i++) {
                    int bitsToWrite = Math.min(remainingBitLength, 8);
                    bo.writeByte(false, bitsToWrite, bytes[i]);
                    remainingBitLength -= bitsToWrite;
                }
            } else {
                // MSB in bytes.length
                for (int i = bytes.length - 1; i >= 0; i--) {
                    int bitsToWrite = Math.min(remainingBitLength, 8);
                    bo.writeByte(false, bitsToWrite, bytes[i]);
                    remainingBitLength -= bitsToWrite;
                }
            }
        } catch (IOException e) {
            throw new SerializationException("Error reading", e);
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
        if (bitLength <= 0) {
            throw new SerializationException("long must contain at least 1 bit");
        }
        int actualBitLength = value.bitLength();
        boolean negative = value.compareTo(BigInteger.ZERO) < 0;
        int bitLengthIncludingPossibleSign = actualBitLength + (negative ? 1 : 0);
        if (bitLength < bitLengthIncludingPossibleSign) {
            throw new SerializationException("bit length including possible sign " + bitLengthIncludingPossibleSign + " exceeds supplied bit length " + bitLength);
        }
        byte[] bytes = value.toByteArray();
        int remainingBitLength = bitLengthIncludingPossibleSign;
        try {
            if (byteOrder != ByteOrder.LITTLE_ENDIAN) {
                // MSB in 0
                for (int i = 0; i < bytes.length; i++) {
                    int bitsToWrite = Math.min(remainingBitLength, 8);
                    bo.writeByte(false, bitsToWrite, bytes[i]);
                    remainingBitLength -= bitsToWrite;
                }
            } else {
                // MSB in bytes.length
                for (int i = bytes.length - 1; i >= 0; i--) {
                    int bitsToWrite = Math.min(remainingBitLength, 8);
                    bo.writeByte(false, bitsToWrite, bytes[i]);
                    remainingBitLength -= bitsToWrite;
                }
            }
        } catch (IOException e) {
            throw new SerializationException("Error reading", e);
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
        int fixedByteLength = (int) Math.ceil((float) bitLength / 8.0);

        if (bitLength == 0) {
            fixedByteLength = bytes.length;
        }

        try {
            int offset = bytes.length - fixedByteLength;
            while (offset < 0) {
                bo.writeByte(false, 8, (byte) 0x00);
                offset++;
            }
            for (int i = offset; i < bytes.length; i++) {
                bo.writeByte(false, 8, bytes[i]);
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
