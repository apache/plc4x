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

package org.apache.plc4x.java.spi.fields.data;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Stack;

/**
 * Minimal dummy ReadBuffer/WriteBuffer implementations for testing DataReaders/DataWriters
 * without depending on concrete byte-based buffer modules.
 */
public class TestBuffers {

    public static class DummyReadBuffer implements ReadBuffer {
        public boolean bitValue;
        public byte[] bitsValue;
        public byte unsignedByteValue;
        public short unsignedShortValue;
        public int unsignedIntValue;
        public long unsignedLongValue;
        public BigInteger unsignedBigIntegerValue;
        public byte signedByteValue;
        public short signedShortValue;
        public int signedIntValue;
        public long signedLongValue;
        public BigInteger signedBigIntegerValue;
        public float floatValue;
        public double doubleValue;
        public BigDecimal bigDecimalValue;
        public String stringValue;

        public Stack<WithOption[]> context = new Stack<>();

        private int positionInBits = 0;

        @Override
        public boolean readBit(WithOption... options) { return bitValue; }

        private int bitsCursor = 0;
        public byte[] readBits(int numBits, WithOption... options) {
            if (bitsValue != null && numBits == 8) {
                // Return next byte sequentially
                int idx = Math.min(bitsCursor, bitsValue.length - 1);
                byte b = bitsValue[idx];
                bitsCursor = Math.min(bitsCursor + 1, bitsValue.length);
                return new byte[] { b };
            }
            return bitsValue != null ? bitsValue : new byte[(numBits+7)/8];
        }

        @Override
        public byte readUnsignedByte(int numBits, WithOption... options) { return unsignedByteValue; }

        @Override
        public short readUnsignedShort(int numBits, WithOption... options) { return unsignedShortValue; }

        @Override
        public int readUnsignedInt(int numBits, WithOption... options) { return unsignedIntValue; }

        @Override
        public long readUnsignedLong(int numBits, WithOption... options) { return unsignedLongValue; }

        @Override
        public BigInteger readUnsignedBigInteger(int bitLength, WithOption... options) { return unsignedBigIntegerValue; }

        @Override
        public byte readSignedByte(int numBits, WithOption... options) { return signedByteValue; }

        @Override
        public short readSignedShort(int numBits, WithOption... options) { return signedShortValue; }

        @Override
        public int readSignedInt(int numBits, WithOption... options) { return signedIntValue; }

        @Override
        public long readSignedLong(int numBits, WithOption... options) { return signedLongValue; }

        @Override
        public BigInteger readSignedBigInteger(int bitLength, WithOption... options) { return signedBigIntegerValue; }

        @Override
        public float readFloat(int numBits, WithOption... options) { return floatValue; }

        @Override
        public double readDouble(int numBits, WithOption... options) { return doubleValue; }

        @Override
        public BigDecimal readBigDecimal(int numBits, WithOption... options) { return bigDecimalValue; }

        @Override
        public String readString(int numBits, WithOption... options) { return stringValue; }

        @Override
        public ReadBuffer createSubBuffer(int numBits, WithOption... options) { return this; }

        @Override
        public int getPositionInBits() { return positionInBits; }

        @Override
        public int getRemainingBits() { return Integer.MAX_VALUE; }

        @Override
        public void setPositionInBits(int positionInBits) { this.positionInBits = positionInBits; }

        @Override
        public void pushContext(WithOption... options) throws BufferException {
            context.push(options);
        }

        @Override
        public void popContext(WithOption... options) throws BufferException {
            context.pop();
        }

        @Override
        public WithOption[] getContext() {
            if (context.isEmpty()) {
                return new WithOption[0];
            }
            return context.peek();
        }
    }

    public static class DummyWriteBuffer implements WriteBuffer {
        public Boolean bitWritten;
        public byte[] bitsWritten;
        public Message messageWritten;
        public Byte unsignedByteWritten;
        public Short unsignedShortWritten;
        public Integer unsignedIntWritten;
        public Long unsignedLongWritten;
        public BigInteger unsignedBigIntegerWritten;
        public Byte signedByteWritten;
        public Short signedShortWritten;
        public Integer signedIntWritten;
        public Long signedLongWritten;
        public BigInteger signedBigIntegerWritten;
        public Float floatWritten;
        public Double doubleWritten;
        public BigDecimal bigDecimalWritten;
        public String stringWritten;

        public Stack<WithOption[]> context = new Stack<>();

        private int positionInBits = 0;
        private byte[] bytes = new byte[0];

        @Override
        public void writeBit(boolean value, WithOption... options) { bitWritten = value; }

        @Override
        public void writeBits(int numBits, byte[] value, WithOption... options) { bitsWritten = value; }

        @Override
        public void writeUnsignedByte(int numBits, byte value, WithOption... options) { unsignedByteWritten = value; }

        @Override
        public void writeUnsignedShort(int numBits, short value, WithOption... options) { unsignedShortWritten = value; }

        @Override
        public void writeUnsignedInt(int numBits, int value, WithOption... options) { unsignedIntWritten = value; }

        @Override
        public void writeUnsignedLong(int numBits, long value, WithOption... options) { unsignedLongWritten = value; }

        @Override
        public void writeUnsignedBigInteger(int numBits, BigInteger value, WithOption... options) { unsignedBigIntegerWritten = value; }

        @Override
        public void writeSignedByte(int numBits, byte value, WithOption... options) { signedByteWritten = value; }

        @Override
        public void writeSignedShort(int numBits, short value, WithOption... options) { signedShortWritten = value; }

        @Override
        public void writeSignedInt(int numBits, int value, WithOption... options) { signedIntWritten = value; }

        @Override
        public void writeSignedLong(int numBits, long value, WithOption... options) { signedLongWritten = value; }

        @Override
        public void writeSignedBigInteger(int numBits, BigInteger value, WithOption... options) { signedBigIntegerWritten = value; }

        @Override
        public void writeFloat(int numBits, float value, WithOption... options) { floatWritten = value; }

        @Override
        public void writeDouble(int numBits, double value, WithOption... options) { doubleWritten = value; }

        @Override
        public void writeBigDecimal(int numBits, BigDecimal value, WithOption... options) { bigDecimalWritten = value; }

        @Override
        public void writeString(int numBits, String value, WithOption... options) { stringWritten = value; }

        @Override
        public void writeMessage(Message message) { messageWritten = message; }

        @Override
        public WriteBuffer createSubBuffer(int numBits, WithOption... options) { return this; }

        @Override
        public int getPositionInBits() { return positionInBits; }

        @Override
        public int getRemainingBits() { return Integer.MAX_VALUE; }

        @Override
        public byte[] getBytes() { return bytes; }

        @Override
        public boolean isByteBased() { return false; }

        public void setPositionInBits(int positionInBits) { this.positionInBits = positionInBits; }

        @Override
        public void pushContext(WithOption... options) throws BufferException {
            context.push(options);
        }

        @Override
        public void popContext(WithOption... options) throws BufferException {
            context.pop();
        }

        @Override
        public WithOption[] getContext() { return context.peek(); }
    }
}
