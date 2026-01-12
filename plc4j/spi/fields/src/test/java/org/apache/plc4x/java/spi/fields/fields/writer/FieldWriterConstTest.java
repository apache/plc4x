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

class FieldWriterConstTest {

    @Test
    void writeConstFieldDelegatesToDataWriter() throws Exception {
        FieldWriterConst<Integer> writer = new FieldWriterConst<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeConstField(42, dataWriter);

        assertEquals(42, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldWithOptions() throws Exception {
        FieldWriterConst<String> writer = new FieldWriterConst<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("constField");
        writer.writeConstField("constant_value", dataWriter, option);

        assertEquals("constant_value", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldWithNullValue() throws Exception {
        FieldWriterConst<String> writer = new FieldWriterConst<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeConstField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldMultipleCalls() throws Exception {
        FieldWriterConst<Integer> writer = new FieldWriterConst<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeConstField(123, dataWriter);
        writer.writeConstField(456, dataWriter);

        assertEquals(456, dataWriter.getLastValue()); // Should be last written value
        assertEquals(2, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldWithByteArray() throws Exception {
        FieldWriterConst<byte[]> writer = new FieldWriterConst<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] constData = {(byte) 0xFF, (byte) 0xEE, (byte) 0xDD, (byte) 0xCC};
        writer.writeConstField(constData, dataWriter);

        assertArrayEquals(constData, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldWithBooleanType() throws Exception {
        FieldWriterConst<Boolean> writer = new FieldWriterConst<>();
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        writer.writeConstField(true, dataWriter);

        assertTrue(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldWithLongValue() throws Exception {
        FieldWriterConst<Long> writer = new FieldWriterConst<>();
        SimpleDataWriter<Long> dataWriter = new SimpleDataWriter<>();

        Long constValue = 0xDEADBEEFL;
        writer.writeConstField(constValue, dataWriter);

        assertEquals(constValue, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldWithDoubleValue() throws Exception {
        FieldWriterConst<Double> writer = new FieldWriterConst<>();
        SimpleDataWriter<Double> dataWriter = new SimpleDataWriter<>();

        double constValue = Math.PI;
        writer.writeConstField(constValue, dataWriter);

        assertEquals(constValue, dataWriter.getLastValue(), 1e-10);
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeConstFieldWithMultipleOptions() throws Exception {
        FieldWriterConst<Integer> writer = new FieldWriterConst<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        WithOption[] options = {
            WithOption.WithName("constField"),
            WithOption.WithRenderAsList(false)
        };
        writer.writeConstField(999, dataWriter, options);

        assertEquals(999, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }
}