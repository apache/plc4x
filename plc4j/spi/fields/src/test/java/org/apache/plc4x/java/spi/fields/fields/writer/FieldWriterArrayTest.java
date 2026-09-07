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

package org.apache.plc4x.java.spi.fields.fields.writer;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyWriteBuffer;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataWriter;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class FieldWriterArrayTest {

    @Test
    void writeByteArrayFieldWritesToDataWriter() throws Exception {
        FieldWriterArray<Byte> writer = new FieldWriterArray<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] testData = {1, 2, 3, 4, 5};
        writer.writeByteArrayField(testData, dataWriter);

        assertArrayEquals(testData, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeByteArrayFieldWithNullArray() throws Exception {
        FieldWriterArray<Byte> writer = new FieldWriterArray<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        writer.writeByteArrayField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount()); // Should not write anything
    }

    @Test
    void writeByteArrayFieldWithEmptyArray() throws Exception {
        FieldWriterArray<Byte> writer = new FieldWriterArray<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] emptyArray = {};
        writer.writeByteArrayField(emptyArray, dataWriter);

        assertArrayEquals(emptyArray, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeSimpleTypeArrayFieldWritesEachElement() throws Exception {
        FieldWriterArray<Integer> writer = new FieldWriterArray<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        List<Integer> values = Arrays.asList(10, 20, 30);
        writer.writeSimpleTypeArrayField(values, dataWriter);

        assertEquals(30, dataWriter.getLastValue()); // Should be last written value
        assertEquals(3, dataWriter.getWriteCount()); // Should write each element
    }

    @Test
    void writeSimpleTypeArrayFieldWithNullList() throws Exception {
        FieldWriterArray<String> writer = new FieldWriterArray<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeSimpleTypeArrayField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount()); // Should not write anything
    }

    @Test
    void writeSimpleTypeArrayFieldWithEmptyList() throws Exception {
        FieldWriterArray<String> writer = new FieldWriterArray<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        List<String> emptyList = Collections.emptyList();
        writer.writeSimpleTypeArrayField(emptyList, dataWriter);

        assertNull(dataWriter.getLastValue()); // No values written
        assertEquals(0, dataWriter.getWriteCount()); // No writes for empty list
    }

    @Test
    void writeSimpleTypeArrayFieldWithSingleElement() throws Exception {
        FieldWriterArray<String> writer = new FieldWriterArray<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        List<String> singletonList = Collections.singletonList("single");
        writer.writeSimpleTypeArrayField(singletonList, dataWriter);

        assertEquals("single", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeComplexTypeArrayFieldWritesEachMessage() throws Exception {
        FieldWriterArray<Message> writer = new FieldWriterArray<>();
        DummyWriteBuffer writeBuffer = new DummyWriteBuffer();

        Message message1 = new TestMessage("msg1");
        Message message2 = new TestMessage("msg2");
        List<Message> messages = Arrays.asList(message1, message2);

        writer.writeComplexTypeArrayField(messages, writeBuffer);

        // The complex array field writes each message by calling serialize()
        // Since our TestMessage writes to writeString, check that
        assertEquals("msg2", writeBuffer.stringWritten); // Last message content written
    }

    @Test
    void writeComplexTypeArrayFieldWithNullList() throws Exception {
        FieldWriterArray<Message> writer = new FieldWriterArray<>();
        DummyWriteBuffer writeBuffer = new DummyWriteBuffer();

        writer.writeComplexTypeArrayField(null, writeBuffer);

        assertNull(writeBuffer.stringWritten); // Should not write anything
    }

    @Test
    void writeComplexTypeArrayFieldWithEmptyList() throws Exception {
        FieldWriterArray<Message> writer = new FieldWriterArray<>();
        DummyWriteBuffer writeBuffer = new DummyWriteBuffer();

        List<Message> emptyList = Collections.emptyList();
        writer.writeComplexTypeArrayField(emptyList, writeBuffer);

        assertNull(writeBuffer.stringWritten); // No messages written
    }

    @Test
    void writeComplexTypeArrayFieldHandlesBufferException() {
        FieldWriterArray<Message> writer = new FieldWriterArray<>();

        WriteBuffer faultyBuffer = new WriteBuffer() {
            @Override
            public void pushContext(WithOption... options) throws BufferException {}

            @Override
            public void popContext(WithOption... options) throws BufferException {}

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
                throw new BufferException("Simulated write error");
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

        List<Message> messages = Collections.singletonList(new TestMessage("test"));

        assertThrows(BufferException.class, () ->
            writer.writeComplexTypeArrayField(messages, faultyBuffer)
        );
    }

    @Test
    void writeArraysWithOptions() throws Exception {
        FieldWriterArray<String> writer = new FieldWriterArray<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("arrayField");
        List<String> values = Arrays.asList("a", "b", "c");
        writer.writeSimpleTypeArrayField(values, dataWriter, option);

        assertEquals("c", dataWriter.getLastValue());
        assertEquals(3, dataWriter.getWriteCount());
    }

    // Test helper class
    private static class TestMessage implements Message {
        private final String content;

        public TestMessage(String content) {
            this.content = content;
        }

        @Override
        public int getLengthInBytes() {
            return content.length();
        }

        @Override
        public int getLengthInBits() {
            return content.length() * 8;
        }

        @Override
        public void serialize(WriteBuffer writeBuffer) throws BufferException {
            writeBuffer.writeString(content.length() * 8, content);
        }

        @Override
        public boolean equals(Object obj) {
            return obj instanceof TestMessage && content.equals(((TestMessage) obj).content);
        }

        @Override
        public int hashCode() {
            return content.hashCode();
        }
    }
}