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

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataWriter;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FieldWriterEnumTest {

    enum TestEnum {
        VALUE_ONE,
        VALUE_TWO
    }

    @Test
    void writeEnumFieldDelegatesToDataWriter() throws Exception {
        FieldWriterEnum<TestEnum> writer = new FieldWriterEnum<>();
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        writer.writeEnumField(TestEnum.VALUE_ONE, dataWriter);

        assertEquals(TestEnum.VALUE_ONE, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeEnumFieldWithOptions() throws Exception {
        FieldWriterEnum<TestEnum> writer = new FieldWriterEnum<>();
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("enumField");
        writer.writeEnumField(TestEnum.VALUE_TWO, dataWriter, option);

        assertEquals(TestEnum.VALUE_TWO, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeEnumFieldManagesContext() throws Exception {
        FieldWriterEnum<String> writer = new FieldWriterEnum<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeEnumField("ENUM_VALUE", dataWriter);

        assertEquals("ENUM_VALUE", dataWriter.getLastValue());
    }

    @Test
    void writeEnumFieldWithNullValue() throws Exception {
        FieldWriterEnum<TestEnum> writer = new FieldWriterEnum<>();
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        writer.writeEnumField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeEnumFieldMultipleCalls() throws Exception {
        FieldWriterEnum<TestEnum> writer = new FieldWriterEnum<>();
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        writer.writeEnumField(TestEnum.VALUE_ONE, dataWriter);
        writer.writeEnumField(TestEnum.VALUE_TWO, dataWriter);

        assertEquals(TestEnum.VALUE_TWO, dataWriter.getLastValue());
        assertEquals(2, dataWriter.getWriteCount());
    }

    @Test
    void writeEnumFieldWithStringEnum() throws Exception {
        FieldWriterEnum<String> writer = new FieldWriterEnum<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeEnumField("CUSTOM_ENUM_VALUE", dataWriter);

        assertEquals("CUSTOM_ENUM_VALUE", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeEnumFieldWithIntegerEnum() throws Exception {
        FieldWriterEnum<Integer> writer = new FieldWriterEnum<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeEnumField(42, dataWriter);

        assertEquals(42, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeEnumFieldWithMultipleOptions() throws Exception {
        FieldWriterEnum<TestEnum> writer = new FieldWriterEnum<>();
        SimpleDataWriter<TestEnum> dataWriter = new SimpleDataWriter<>();

        WithOption[] options = {
            WithOption.WithName("enumField"),
            WithOption.WithRenderAsList(false)
        };
        writer.writeEnumField(TestEnum.VALUE_ONE, dataWriter, options);

        assertEquals(TestEnum.VALUE_ONE, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeEnumFieldContextStackBehavior() throws Exception {
        FieldWriterEnum<String> writer = new FieldWriterEnum<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("testEnum");

        // Verify that context is properly managed for nested calls
        writer.writeEnumField("first", dataWriter, option);
        assertEquals("first", dataWriter.getLastValue());

        writer.writeEnumField("second", dataWriter, option);
        assertEquals("second", dataWriter.getLastValue());

        assertEquals(2, dataWriter.getWriteCount());
    }
}