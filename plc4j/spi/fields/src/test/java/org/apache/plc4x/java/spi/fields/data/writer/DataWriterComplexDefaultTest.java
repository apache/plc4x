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

package org.apache.plc4x.java.spi.fields.data.writer;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyWriteBuffer;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DataWriterComplexDefaultTest {

    @Test
    void writeNullValueLogsWarningButContinues() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterComplexDefault<Message> writer = new DataWriterComplexDefault<>(wb);

        writer.write(null);

        assertNull(wb.messageWritten);
    }

    @Test
    void writeValidMessageDelegatesToWriteBuffer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterComplexDefault<Message> writer = new DataWriterComplexDefault<>(wb);

        Message testMessage = new Message() {
            @Override
            public int getLengthInBytes() {
                return 4;
            }

            @Override
            public int getLengthInBits() {
                return 32;
            }

            @Override
            public void serialize(WriteBuffer writeBuffer) throws BufferException {
                writeBuffer.writeUnsignedInt(32, 0xDEADBEEF);
            }
        };

        writer.write(testMessage);

        assertEquals(testMessage, wb.messageWritten);
    }

    @Test
    void writeWithOptionsPassesOptionsToBuffer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterComplexDefault<Message> writer = new DataWriterComplexDefault<>(wb);

        Message testMessage = new Message() {
            @Override
            public int getLengthInBytes() {
                return 1;
            }

            @Override
            public int getLengthInBits() {
                return 8;
            }

            @Override
            public void serialize(WriteBuffer writeBuffer) throws BufferException {
                writeBuffer.writeSignedByte(8, (byte) 0x42);
            }
        };

        WithOption option = WithOption.WithName("test");
        writer.write(testMessage, option);

        assertEquals(testMessage, wb.messageWritten);
    }

    @Test
    void contextMethodsDelegateToWriteBuffer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterComplexDefault<Message> writer = new DataWriterComplexDefault<>(wb);

        WithOption option = WithOption.WithName("test");

        assertTrue(wb.context.isEmpty());
        writer.pushContext(option);
        assertFalse(wb.context.isEmpty());
        writer.popContext(option);
        assertTrue(wb.context.isEmpty());
    }

    @Test
    void getWriteBufferReturnsOriginalBuffer() {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterComplexDefault<Message> writer = new DataWriterComplexDefault<>(wb);

        assertSame(wb, writer.getWriteBuffer());
    }

    @Test
    void contextExceptionIsWrappedInRuntimeException() {
        WriteBuffer faultyBuffer = new WriteBuffer() {
            @Override
            public void pushContext(WithOption... options) throws BufferException {
                throw new BufferException("Test exception");
            }

            @Override
            public void popContext(WithOption... options) throws BufferException {
                throw new BufferException("Test exception");
            }

            @Override
            public WithOption[] getContext() {
                return new WithOption[0];
            }

            // Other required methods would throw UnsupportedOperationException
            @Override
            public void writeBit(boolean value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeBits(int bitLength, byte[] bytes, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeUnsignedByte(int bitLength, byte value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeUnsignedShort(int bitLength, short value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeUnsignedInt(int bitLength, int value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeUnsignedLong(int bitLength, long value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeUnsignedBigInteger(int bitLength, java.math.BigInteger value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeSignedByte(int bitLength, byte value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeSignedShort(int bitLength, short value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeSignedInt(int bitLength, int value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeSignedLong(int bitLength, long value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeSignedBigInteger(int bitLength, java.math.BigInteger value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeFloat(int bitLength, float value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeDouble(int bitLength, double value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeBigDecimal(int bitLength, java.math.BigDecimal value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public void writeString(int bitLength, String value, WithOption... options) throws BufferException {
                throw new UnsupportedOperationException();
            }

            @Override
            public WriteBuffer createSubBuffer(int numBits, WithOption... options) throws BufferException {
                return null;
            }

            @Override
            public int getPositionInBits() {
                return 0;
            }

            @Override
            public int getRemainingBits() {
                return 0;
            }

            @Override
            public byte[] getBytes() {
                return new byte[0];
            }

        };

        DataWriterComplexDefault<Message> writer = new DataWriterComplexDefault<>(faultyBuffer);

        assertThrows(RuntimeException.class, () -> writer.pushContext());
        assertThrows(RuntimeException.class, () -> writer.popContext());
    }
}