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
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DataWriterEnumDefaultTest {

    enum TestEnum {
        FIRST(1, "First"),
        SECOND(2, "Second"),
        THIRD(3, "Third");

        private final int value;
        private final String description;

        TestEnum(int value, String description) {
            this.value = value;
            this.description = description;
        }

        public int getValue() {
            return value;
        }

        public String getDescription() {
            return description;
        }
    }

    @Test
    void writeEnumValueUsesSerializer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriter<Integer> intWriter = DataWriterFactory.writeUnsignedInt(wb, 32);

        DataWriterEnumDefault<TestEnum, Integer> enumWriter = new DataWriterEnumDefault<>(
            TestEnum::getValue,
            TestEnum::getDescription,
            intWriter
        );

        enumWriter.write(TestEnum.SECOND);

        assertEquals(2, (int) wb.unsignedIntWritten);
    }

    @Test
    void writeNullValueLogsWarningAndReturns() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriter<Integer> intWriter = DataWriterFactory.writeUnsignedInt(wb, 32);

        DataWriterEnumDefault<TestEnum, Integer> enumWriter = new DataWriterEnumDefault<>(
            TestEnum::getValue,
            TestEnum::getDescription,
            intWriter
        );

        enumWriter.write(null);

        assertNull(wb.unsignedIntWritten);
    }

    @Test
    void writeWithOptionsAddsStringRepresentation() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriter<Integer> intWriter = DataWriterFactory.writeUnsignedInt(wb, 32);

        DataWriterEnumDefault<TestEnum, Integer> enumWriter = new DataWriterEnumDefault<>(
            TestEnum::getValue,
            TestEnum::getDescription,
            intWriter
        );

        WithOption option = WithOption.WithName("testField");
        enumWriter.write(TestEnum.THIRD, option);

        assertEquals(3, (int) wb.unsignedIntWritten);
    }

    @Test
    void contextMethodsDelegateToDataWriter() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriter<Integer> intWriter = DataWriterFactory.writeUnsignedInt(wb, 32);

        DataWriterEnumDefault<TestEnum, Integer> enumWriter = new DataWriterEnumDefault<>(
            TestEnum::getValue,
            TestEnum::getDescription,
            intWriter
        );

        WithOption option = WithOption.WithName("test");

        assertTrue(wb.context.isEmpty());
        enumWriter.pushContext(option);
        assertFalse(wb.context.isEmpty());
        enumWriter.popContext(option);
        assertTrue(wb.context.isEmpty());
    }

    @Test
    void getWriteBufferDelegatesToDataWriter() {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriter<Integer> intWriter = DataWriterFactory.writeUnsignedInt(wb, 32);

        DataWriterEnumDefault<TestEnum, Integer> enumWriter = new DataWriterEnumDefault<>(
            TestEnum::getValue,
            TestEnum::getDescription,
            intWriter
        );

        assertSame(wb, enumWriter.getWriteBuffer());
    }

    @Test
    void writeWithCustomSerializerAndNamer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriter<Integer> intWriter = DataWriterFactory.writeUnsignedInt(wb, 32);

        DataWriterEnumDefault<TestEnum, Integer> enumWriter = new DataWriterEnumDefault<>(
            TestEnum::getValue,
            TestEnum::getDescription,
            intWriter
        );

        // Test the public write method with custom serializer and namer
        enumWriter.write(
            TestEnum.FIRST,
            e -> e.getValue() * 10,  // Custom serializer
            e -> "Custom_" + e.name(), // Custom namer
            intWriter
        );

        assertEquals(10, (int) wb.unsignedIntWritten);
    }

    @Test
    void exceptionPropagationFromDataWriter() throws Exception {
        DataWriter<Integer> faultyWriter = new DataWriter<Integer>() {
            @Override
            public void write(Integer value, WithOption... options) throws BufferException {
                throw new BufferException("Test exception");
            }

            @Override
            public WriteBuffer getWriteBuffer() {
                return null;
            }

            @Override
            public void pushContext(WithOption... options) throws BufferException {
                throw new BufferException("Push context exception");
            }

            @Override
            public void popContext(WithOption... options) throws BufferException {
                throw new BufferException("Pop context exception");
            }
        };

        DataWriterEnumDefault<TestEnum, Integer> enumWriter = new DataWriterEnumDefault<>(
            TestEnum::getValue,
            TestEnum::getDescription,
            faultyWriter
        );

        assertThrows(BufferException.class, () -> enumWriter.write(TestEnum.FIRST));
        assertThrows(BufferException.class, () -> enumWriter.pushContext());
        assertThrows(BufferException.class, () -> enumWriter.popContext());
    }
}