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

class FieldWriterPaddingTest {

    @Test
    void writePaddingFieldWritesSpecifiedNumberOfTimes() throws Exception {
        FieldWriterPadding<Byte> writer = new FieldWriterPadding<>();
        SimpleDataWriter<Byte> dataWriter = new SimpleDataWriter<>();

        byte paddingValue = 0;
        writer.writePaddingField(3, paddingValue, dataWriter);

        assertEquals(paddingValue, dataWriter.getLastValue());
        assertEquals(3, dataWriter.getWriteCount()); // Should write 3 times
    }

    @Test
    void writePaddingFieldWithZeroPadding() throws Exception {
        FieldWriterPadding<Byte> writer = new FieldWriterPadding<>();
        SimpleDataWriter<Byte> dataWriter = new SimpleDataWriter<>();

        byte paddingValue = 0;
        writer.writePaddingField(0, paddingValue, dataWriter);

        assertNull(dataWriter.getLastValue()); // No writes
        assertEquals(0, dataWriter.getWriteCount()); // Should not write anything
    }

    @Test
    void writePaddingFieldWithNegativePadding() throws Exception {
        FieldWriterPadding<Byte> writer = new FieldWriterPadding<>();
        SimpleDataWriter<Byte> dataWriter = new SimpleDataWriter<>();

        byte paddingValue = 0;
        writer.writePaddingField(-1, paddingValue, dataWriter);

        assertNull(dataWriter.getLastValue()); // No writes
        assertEquals(0, dataWriter.getWriteCount()); // Should not write anything
    }

    @Test
    void writePaddingFieldWithSinglePadding() throws Exception {
        FieldWriterPadding<Integer> writer = new FieldWriterPadding<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        int paddingValue = 999;
        writer.writePaddingField(1, paddingValue, dataWriter);

        assertEquals(paddingValue, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writePaddingFieldWithLargePadding() throws Exception {
        FieldWriterPadding<String> writer = new FieldWriterPadding<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        String paddingValue = "pad";
        writer.writePaddingField(10, paddingValue, dataWriter);

        assertEquals(paddingValue, dataWriter.getLastValue());
        assertEquals(10, dataWriter.getWriteCount()); // Should write 10 times
    }

    @Test
    void writePaddingFieldWithOptions() throws Exception {
        FieldWriterPadding<Character> writer = new FieldWriterPadding<>();
        SimpleDataWriter<Character> dataWriter = new SimpleDataWriter<>();

        char paddingValue = ' ';
        WithOption option = WithOption.WithName("paddingField");
        writer.writePaddingField(5, paddingValue, dataWriter, option);

        assertEquals(paddingValue, dataWriter.getLastValue());
        assertEquals(5, dataWriter.getWriteCount());
    }

    @Test
    void writePaddingFieldManagesContext() throws Exception {
        FieldWriterPadding<Byte> writer = new FieldWriterPadding<>();
        SimpleDataWriter<Byte> dataWriter = new SimpleDataWriter<>();

        byte paddingValue = 42;
        writer.writePaddingField(2, paddingValue, dataWriter);

        assertEquals(paddingValue, dataWriter.getLastValue());
        assertEquals(2, dataWriter.getWriteCount());
    }

    @Test
    void writePaddingFieldWithNullPaddingValue() throws Exception {
        FieldWriterPadding<String> writer = new FieldWriterPadding<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writePaddingField(3, null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(3, dataWriter.getWriteCount()); // Should still write 3 times with null
    }

    @Test
    void writePaddingFieldWithBooleanPadding() throws Exception {
        FieldWriterPadding<Boolean> writer = new FieldWriterPadding<>();
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        writer.writePaddingField(4, true, dataWriter);

        assertTrue(dataWriter.getLastValue());
        assertEquals(4, dataWriter.getWriteCount());
    }

    @Test
    void writePaddingFieldWithArrayPadding() throws Exception {
        FieldWriterPadding<byte[]> writer = new FieldWriterPadding<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] paddingValue = {0, 0, 0, 0};
        writer.writePaddingField(2, paddingValue, dataWriter);

        assertArrayEquals(paddingValue, dataWriter.getLastValue());
        assertEquals(2, dataWriter.getWriteCount());
    }

    @Test
    void writePaddingFieldMultipleOptions() throws Exception {
        FieldWriterPadding<Integer> writer = new FieldWriterPadding<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        WithOption[] options = {
            WithOption.WithName("paddingField"),
            WithOption.WithRenderAsList(false)
        };
        writer.writePaddingField(3, 123, dataWriter, options);

        assertEquals(123, dataWriter.getLastValue());
        assertEquals(3, dataWriter.getWriteCount());
    }

    @Test
    void writePaddingFieldContextualBehavior() throws Exception {
        FieldWriterPadding<String> writer = new FieldWriterPadding<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        // Test multiple padding operations
        writer.writePaddingField(2, "first", dataWriter);
        writer.writePaddingField(3, "second", dataWriter);
        writer.writePaddingField(1, "third", dataWriter);

        assertEquals("third", dataWriter.getLastValue()); // Last value written
        assertEquals(6, dataWriter.getWriteCount()); // 2 + 3 + 1 = 6 total writes
    }
}