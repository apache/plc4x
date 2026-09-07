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

class FieldWriterOptionalTest {

    @Test
    void writeOptionalFieldWithTrueConditionAndNonNullValue() throws Exception {
        FieldWriterOptional<String> writer = new FieldWriterOptional<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField("optional_value", dataWriter, true);

        assertEquals("optional_value", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalFieldWithFalseCondition() throws Exception {
        FieldWriterOptional<String> writer = new FieldWriterOptional<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField("should_not_write", dataWriter, false);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount()); // Should not write
    }

    @Test
    void writeOptionalFieldWithNullValue() throws Exception {
        FieldWriterOptional<String> writer = new FieldWriterOptional<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField(null, dataWriter, true);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount()); // Should not write
    }

    @Test
    void writeOptionalFieldWithTrueConditionAndNullValue() throws Exception {
        FieldWriterOptional<Integer> writer = new FieldWriterOptional<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField(null, dataWriter, true);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount()); // Should not write
    }

    @Test
    void writeOptionalFieldWithFalseConditionAndNonNullValue() throws Exception {
        FieldWriterOptional<Integer> writer = new FieldWriterOptional<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField(42, dataWriter, false);

        assertNull(dataWriter.getLastValue());
        assertEquals(0, dataWriter.getWriteCount()); // Should not write
    }

    @Test
    void writeOptionalFieldWithOptions() throws Exception {
        FieldWriterOptional<String> writer = new FieldWriterOptional<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("optionalField");
        writer.writeOptionalField("test", dataWriter, true, option);

        assertEquals("test", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalFieldMultipleCallsWithDifferentConditions() throws Exception {
        FieldWriterOptional<Integer> writer = new FieldWriterOptional<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField(1, dataWriter, true);   // Should write
        writer.writeOptionalField(2, dataWriter, false);  // Should not write
        writer.writeOptionalField(3, dataWriter, true);   // Should write

        assertEquals(3, dataWriter.getLastValue()); // Last written value
        assertEquals(2, dataWriter.getWriteCount()); // Only 2 writes
    }

    @Test
    void writeOptionalFieldWithBooleanValue() throws Exception {
        FieldWriterOptional<Boolean> writer = new FieldWriterOptional<>();
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField(true, dataWriter, true);

        assertTrue(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalFieldWithFalseBooleanValue() throws Exception {
        FieldWriterOptional<Boolean> writer = new FieldWriterOptional<>();
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField(false, dataWriter, true);

        assertFalse(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalFieldWithByteArray() throws Exception {
        FieldWriterOptional<byte[]> writer = new FieldWriterOptional<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] testData = {1, 2, 3, 4, 5};
        writer.writeOptionalField(testData, dataWriter, true);

        assertArrayEquals(testData, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeOptionalFieldWithEmptyString() throws Exception {
        FieldWriterOptional<String> writer = new FieldWriterOptional<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField("", dataWriter, true);

        assertEquals("", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount()); // Empty string is not null
    }

    @Test
    void writeOptionalFieldWithZeroNumber() throws Exception {
        FieldWriterOptional<Integer> writer = new FieldWriterOptional<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeOptionalField(0, dataWriter, true);

        assertEquals(0, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount()); // Zero is not null
    }

    @Test
    void writeOptionalFieldContextualBehavior() throws Exception {
        FieldWriterOptional<String> writer = new FieldWriterOptional<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        // Test that behavior is consistent across multiple optional field writes
        writer.writeOptionalField("first", dataWriter, true);
        writer.writeOptionalField("second", dataWriter, false);
        writer.writeOptionalField("third", dataWriter, true);
        writer.writeOptionalField("fourth", dataWriter, false);
        writer.writeOptionalField("fifth", dataWriter, true);

        assertEquals("fifth", dataWriter.getLastValue());
        assertEquals(3, dataWriter.getWriteCount()); // Only the true conditions
    }
}