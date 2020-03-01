/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.java.utils;

import com.github.jinahya.bit.io.BitOutput;
import com.github.jinahya.bit.io.BufferByteOutput;
import com.github.jinahya.bit.io.DefaultBitOutput;

import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.ByteBuffer;

public class WriteBuffer {

    private final ByteBuffer bb;
    private final BufferByteOutput bbo;
    private final BitOutput bo;
    private final boolean littleEndian;

    public WriteBuffer(int size) {
        this(size, true);
    }

    public WriteBuffer(int size, boolean littleEndian) {
        bb = ByteBuffer.allocate(size);
        bbo = new BufferByteOutput<>(bb);
        bo = new DefaultBitOutput<>(bbo);
        this.littleEndian = littleEndian;
    }

    public byte[] getData() {
        return bb.array();
    }

    public void writeBit(boolean value) throws ParseException {
        try {
            bo.writeBoolean(value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeUnsignedByte(int bitLength, byte value) throws ParseException {
        if(bitLength <= 0) {
            throw new ParseException("unsigned byte must contain at least 1 bit");
        }
        if(bitLength > 4) {
            throw new ParseException("unsigned byte can only contain max 4 bits");
        }
        try {
            bo.writeByte(true, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeUnsignedShort(int bitLength, short value) throws ParseException {
        if(bitLength <= 0) {
            throw new ParseException("unsigned short must contain at least 1 bit");
        }
        if(bitLength > 8) {
            throw new ParseException("unsigned short can only contain max 8 bits");
        }
        try {
            bo.writeShort(true, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeUnsignedInt(int bitLength, int value) throws ParseException {
        if(bitLength <= 0) {
            throw new ParseException("unsigned int must contain at least 1 bit");
        }
        if(bitLength > 16) {
            throw new ParseException("unsigned int can only contain max 16 bits");
        }
        try {
            if(!littleEndian) {
                value = Integer.reverseBytes(value) >> 16;
            }
            bo.writeInt(true, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeUnsignedLong(int bitLength, long value) throws ParseException {
        if(bitLength <= 0) {
            throw new ParseException("unsigned long must contain at least 1 bit");
        }
        if(bitLength > 32) {
            throw new ParseException("unsigned long can only contain max 32 bits");
        }
        try {
            if(!littleEndian) {
                value = Long.reverseBytes(value) >> 32;
            }
            bo.writeLong(true, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeUnsignedBigInteger(int bitLength, BigInteger value) throws ParseException {
        throw new UnsupportedOperationException("not implemented yet");
    }

    public void writeByte(int bitLength, byte value) throws ParseException {
        if(bitLength <= 0) {
            throw new ParseException("byte must contain at least 1 bit");
        }
        if(bitLength > 8) {
            throw new ParseException("byte can only contain max 8 bits");
        }
        try {
            bo.writeByte(false, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeShort(int bitLength, short value) throws ParseException {
        if(bitLength <= 0) {
            throw new ParseException("short must contain at least 1 bit");
        }
        if(bitLength > 16) {
            throw new ParseException("short can only contain max 16 bits");
        }
        try {
            if(!littleEndian) {
                value = Short.reverseBytes(value);
            }
            bo.writeShort(false, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeInt(int bitLength, int value) throws ParseException {
        if(bitLength <= 0) {
            throw new ParseException("int must contain at least 1 bit");
        }
        if(bitLength > 32) {
            throw new ParseException("int can only contain max 32 bits");
        }
        try {
            if(!littleEndian) {
                value = Integer.reverseBytes(value);
            }
            bo.writeInt(false, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeLong(int bitLength, long value) throws ParseException {
        if(bitLength <= 0) {
            throw new ParseException("long must contain at least 1 bit");
        }
        if(bitLength > 64) {
            throw new ParseException("long can only contain max 64 bits");
        }
        try {
            if(!littleEndian) {
                value = Long.reverseBytes(value);
            }
            bo.writeLong(false, bitLength, value);
        } catch (IOException e) {
            throw new ParseException("Error reading", e);
        }
    }

    public void writeBigInteger(int bitLength, BigInteger value) throws ParseException {
        throw new UnsupportedOperationException("not implemented yet");
    }

    public void writeFloat(int bitLength, float value) throws ParseException {
        throw new UnsupportedOperationException("not implemented yet");
    }

    public void writeDouble(int bitLength, double value) throws ParseException {
        throw new UnsupportedOperationException("not implemented yet");
    }

    public void writeBigDecimal(int bitLength, BigDecimal value) throws ParseException {
        throw new UnsupportedOperationException("not implemented yet");
    }

}
