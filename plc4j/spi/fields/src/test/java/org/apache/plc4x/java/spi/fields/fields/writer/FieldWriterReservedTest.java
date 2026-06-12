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

class FieldWriterReservedTest {

    @Test
    void writeReservedFieldDelegatesToDataWriter() throws Exception {
        FieldWriterReserved<Integer> writer = new FieldWriterReserved<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeReservedField(0, dataWriter);

        assertEquals(0, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldWithOptions() throws Exception {
        FieldWriterReserved<String> writer = new FieldWriterReserved<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("reservedField");
        writer.writeReservedField("reserved", dataWriter, option);

        assertEquals("reserved", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldUnnamedGetsReservedName() throws Exception {
        FieldWriterReserved<Integer> writer = new FieldWriterReserved<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        WithOption unnamedOption = WithOption.WithName("unnamed");
        writer.writeReservedField(123, dataWriter, unnamedOption);

        assertEquals(123, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldUnnamedCaseInsensitive() throws Exception {
        FieldWriterReserved<Integer> writer = new FieldWriterReserved<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        WithOption unnamedOption = WithOption.WithName("UNNAMED");
        writer.writeReservedField(456, dataWriter, unnamedOption);

        assertEquals(456, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldWithNullValue() throws Exception {
        FieldWriterReserved<String> writer = new FieldWriterReserved<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeReservedField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldMultipleCalls() throws Exception {
        FieldWriterReserved<Integer> writer = new FieldWriterReserved<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeReservedField(100, dataWriter);
        writer.writeReservedField(200, dataWriter);

        assertEquals(200, dataWriter.getLastValue()); // Should be last written value
        assertEquals(2, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldWithByteArray() throws Exception {
        FieldWriterReserved<byte[]> writer = new FieldWriterReserved<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] reservedData = {0, 0, 0, 0};
        writer.writeReservedField(reservedData, dataWriter);

        assertArrayEquals(reservedData, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldUnnamedVariationCase() throws Exception {
        FieldWriterReserved<String> writer = new FieldWriterReserved<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        // Test different case variations
        WithOption[] testCases = {
            WithOption.WithName("Unnamed"),
            WithOption.WithName("UnNaMeD"),
            WithOption.WithName("uNnAmEd")
        };

        for (WithOption option : testCases) {
            SimpleDataWriter<String> testDataWriter = new SimpleDataWriter<>();
            writer.writeReservedField("test", testDataWriter, option);

            assertEquals("test", testDataWriter.getLastValue());
            assertEquals(1, testDataWriter.getWriteCount());
        }
    }

    @Test
    void writeReservedFieldWithBooleanType() throws Exception {
        FieldWriterReserved<Boolean> writer = new FieldWriterReserved<>();
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        writer.writeReservedField(false, dataWriter);

        assertFalse(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeReservedFieldWithMultipleOptions() throws Exception {
        FieldWriterReserved<Integer> writer = new FieldWriterReserved<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        WithOption[] options = {
            WithOption.WithName("reservedField"),
            WithOption.WithRenderAsList(true)
        };
        writer.writeReservedField(777, dataWriter, options);

        assertEquals(777, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }
}