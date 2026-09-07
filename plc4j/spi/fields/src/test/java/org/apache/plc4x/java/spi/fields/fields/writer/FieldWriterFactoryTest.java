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
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataWriter;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class FieldWriterFactoryTest {

    @Test
    void writeSimpleTypeArrayFieldDelegates() throws Exception {
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();
        List<String> values = Arrays.asList("a", "b", "c");

        FieldWriterFactory.writeSimpleTypeArrayField(values, dataWriter);

        assertEquals("c", dataWriter.getLastValue()); // Last written value
        assertEquals(3, dataWriter.getWriteCount());
    }

    @Test
    void writeComplexTypeArrayFieldDelegates() throws Exception {
        DummyWriteBuffer writeBuffer = new DummyWriteBuffer();
        List<Message> messages = Arrays.asList(
            new TestMessage("msg1"),
            new TestMessage("msg2")
        );

        FieldWriterFactory.writeComplexTypeArrayField(messages, writeBuffer);

        assertEquals("msg2", writeBuffer.stringWritten); // Last message content
    }

    @Test
    void writeByteArrayFieldDelegates() throws Exception {
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();
        byte[] values = {1, 2, 3, 4, 5};

        FieldWriterFactory.writeByteArrayField(values, dataWriter);

        assertArrayEquals(values, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeChecksumFieldDelegates() throws Exception {
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeChecksumField(42, dataWriter);

        assertEquals(42, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldDelegates() throws Exception {
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeConstField("constant", dataWriter);

        assertEquals("constant", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeEnumFieldDelegates() throws Exception {
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeEnumField(TestEnum.VALUE_ONE, dataWriter);

        assertEquals(TestEnum.VALUE_ONE, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldDelegates() throws Exception {
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeDiscriminatorField(999, dataWriter);

        assertEquals(999, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorEnumFieldDelegates() throws Exception {
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeDiscriminatorEnumField(TestEnum.VALUE_TWO, dataWriter);

        assertEquals(TestEnum.VALUE_TWO, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeImplicitFieldDelegates() throws Exception {
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeImplicitField(true, dataWriter);

        assertTrue(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeManualFieldDelegates() throws Exception {
        DummyWriteBuffer writeBuffer = new DummyWriteBuffer();
        final boolean[] runnableCalled = {false};

        FieldCommons.RunSerializeWrapped runnable = () -> {
            runnableCalled[0] = true;
            writeBuffer.writeUnsignedInt(32, 42);
        };

        FieldWriterFactory.writeManualField(runnable, writeBuffer);

        assertTrue(runnableCalled[0]);
        assertEquals(42, writeBuffer.unsignedIntWritten);
    }

    @Test
    void writeManualArrayFieldWithByteArray() throws Exception {
        DummyWriteBuffer writeBuffer = new DummyWriteBuffer();
        final int[] callCount = {0};

        FieldCommons.ConsumeSerializeWrapped<Byte> runnable = (value) -> {
            callCount[0]++;
            writeBuffer.writeUnsignedByte(8, value);
        };

        byte[] bytes = {1, 2, 3};
        FieldWriterFactory.writeManualArrayField(bytes, runnable, writeBuffer);

        assertEquals(3, callCount[0]); // Should be called for each byte
        assertEquals((byte) 3, writeBuffer.unsignedByteWritten); // Last byte written
    }

    @Test
    void writeManualArrayFieldWithList() throws Exception {
        DummyWriteBuffer writeBuffer = new DummyWriteBuffer();
        final int[] callCount = {0};

        FieldCommons.ConsumeSerializeWrapped<String> runnable = (value) -> {
            callCount[0]++;
            writeBuffer.writeString(value.length() * 8, value);
        };

        List<String> values = Arrays.asList("a", "b", "c");
        FieldWriterFactory.writeManualArrayField(values, runnable, writeBuffer);

        assertEquals(3, callCount[0]); // Should be called for each string
        assertEquals("c", writeBuffer.stringWritten); // Last string written
    }

    @Test
    void writeOptionalFieldWithDefaultCondition() throws Exception {
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeOptionalField("optional", dataWriter);

        assertEquals("optional", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalFieldWithCondition() throws Exception {
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeOptionalField("conditional", dataWriter, true);

        assertEquals("conditional", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalFieldWithFalseCondition() throws Exception {
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeOptionalField("should_not_write", dataWriter, false);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalEnumFieldWithDefaultCondition() throws Exception {
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeOptionalEnumField(TestEnum.VALUE_ONE, dataWriter);

        assertEquals(TestEnum.VALUE_ONE, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalEnumFieldWithCondition() throws Exception {
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeOptionalEnumField(TestEnum.VALUE_TWO, dataWriter, true);

        assertEquals(TestEnum.VALUE_TWO, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalEnumFieldWithFalseCondition() throws Exception {
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeOptionalEnumField(TestEnum.VALUE_ONE, dataWriter, false);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount());
    }

    @Test
    void writePaddingFieldDelegates() throws Exception {
        SimpleDataWriter<Byte> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writePaddingField(3, (byte) 0, dataWriter);

        assertEquals((byte) 0, dataWriter.getLastValue());
        assertEquals(3, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldDelegates() throws Exception {
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeReservedField(0, dataWriter);

        assertEquals(0, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeSimpleFieldDelegates() throws Exception {
        SimpleDataWriter<Double> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeSimpleField(3.14, dataWriter);

        assertEquals(3.14, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeSimpleEnumFieldDelegates() throws Exception {
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeSimpleEnumField(TestEnum.VALUE_ONE, dataWriter);

        assertEquals(TestEnum.VALUE_ONE, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void factoryMethodsWithOptions() throws Exception {
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();
        WithOption option = WithOption.WithName("testField");

        FieldWriterFactory.writeSimpleField("test", dataWriter, option);

        assertEquals("test", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void factoryMethodsWithNullValues() throws Exception {
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        FieldWriterFactory.writeSimpleField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void factoryMethodsWithEmptyCollections() throws Exception {
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();
        List<String> emptyList = Collections.emptyList();

        FieldWriterFactory.writeSimpleTypeArrayField(emptyList, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount());
    }

    // Test helper classes
    enum TestEnum {
        VALUE_ONE,
        VALUE_TWO
    }

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