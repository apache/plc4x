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

class FieldWriterSimpleTest {

    @Test
    void writeSimpleFieldDelegatesToDataWriter() throws Exception {
        FieldWriterSimple<String> writer = new FieldWriterSimple<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeSimpleField("test value", dataWriter);

        assertEquals("test value", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeSimpleFieldWithOptions() throws Exception {
        FieldWriterSimple<Integer> writer = new FieldWriterSimple<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("testField");
        writer.writeSimpleField(42, dataWriter, option);

        assertEquals(42, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeSimpleFieldWithNullValue() throws Exception {
        FieldWriterSimple<String> writer = new FieldWriterSimple<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeSimpleField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeSimpleFieldMultipleCalls() throws Exception {
        FieldWriterSimple<Integer> writer = new FieldWriterSimple<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeSimpleField(123, dataWriter);
        writer.writeSimpleField(456, dataWriter);

        assertEquals(456, dataWriter.getLastValue()); // Should be last written value
        assertEquals(2, dataWriter.getWriteCount());
    }

    @Test
    void writeSimpleFieldWithComplexType() throws Exception {
        FieldWriterSimple<byte[]> writer = new FieldWriterSimple<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] testData = {1, 2, 3, 4, 5};
        writer.writeSimpleField(testData, dataWriter);

        assertArrayEquals(testData, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeSimpleFieldWithBooleanType() throws Exception {
        FieldWriterSimple<Boolean> writer = new FieldWriterSimple<>();
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        writer.writeSimpleField(true, dataWriter);

        assertTrue(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }
}