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

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyWriteBuffer;
import org.apache.plc4x.java.spi.fields.data.reader.DataIoSerializerFunction;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.values.PlcUINT;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DataWriterDataIoDefaultTest {

    @Test
    void writeValueUsesSerializer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        final boolean[] serializerCalled = {false};

        DataIoSerializerFunction<WriteBuffer> serializer = (writeBuffer, value) -> {
            serializerCalled[0] = true;
            writeBuffer.writeUnsignedInt(32, 42);
        };

        DataWriterDataIoDefault writer = new DataWriterDataIoDefault(wb, serializer);
        PlcValue testValue = new PlcUINT(123);

        writer.write(testValue);

        assertTrue(serializerCalled[0]);
        assertEquals(42, (int) wb.unsignedIntWritten);
    }

    @Test
    void writeWithOptionsPassesOptionsToBuffer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        final boolean[] serializerCalled = {false};

        DataIoSerializerFunction<WriteBuffer> serializer = (writeBuffer, value) -> {
            serializerCalled[0] = true;
            writeBuffer.writeSignedByte(8, (byte) 0x42);
        };

        DataWriterDataIoDefault writer = new DataWriterDataIoDefault(wb, serializer);
        PlcValue testValue = new PlcSTRING("test");

        WithOption option = WithOption.WithName("testField");
        writer.write(testValue, option);

        assertTrue(serializerCalled[0]);
        assertEquals((byte) 0x42, wb.signedByteWritten);
    }

    @Test
    void contextMethodsDelegateToWriteBuffer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataIoSerializerFunction<WriteBuffer> serializer = (writeBuffer, value) -> {};

        DataWriterDataIoDefault writer = new DataWriterDataIoDefault(wb, serializer);

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
        DataIoSerializerFunction<WriteBuffer> serializer = (writeBuffer, value) -> {};

        DataWriterDataIoDefault writer = new DataWriterDataIoDefault(wb, serializer);

        assertSame(wb, writer.getWriteBuffer());
    }

    @Test
    void contextExceptionIsWrappedInRuntimeException() {
        WriteBuffer faultyBuffer = new WriteBuffer() {
            @Override
            public void pushContext(WithOption... options) throws BufferException {
                throw new BufferException("Push context exception");
            }

            @Override
            public void popContext(WithOption... options) throws BufferException {
                throw new BufferException("Pop context exception");
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

        DataIoSerializerFunction<WriteBuffer> serializer = (writeBuffer, value) -> {};
        DataWriterDataIoDefault writer = new DataWriterDataIoDefault(faultyBuffer, serializer);

        assertThrows(RuntimeException.class, () -> writer.pushContext());
        assertThrows(RuntimeException.class, () -> writer.popContext());
    }

    @Test
    void serializerExceptionPropagates() {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataIoSerializerFunction<WriteBuffer> faultySerializer = (writeBuffer, value) -> {
            throw new BufferException("Serializer exception");
        };

        DataWriterDataIoDefault writer = new DataWriterDataIoDefault(wb, faultySerializer);
        PlcValue testValue = new PlcUINT(123);

        assertThrows(BufferException.class, () -> writer.write(testValue));
    }
}