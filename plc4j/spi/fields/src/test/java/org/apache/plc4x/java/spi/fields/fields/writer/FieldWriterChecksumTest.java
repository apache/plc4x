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

class FieldWriterChecksumTest {

    @Test
    void writeChecksumFieldDelegatesToDataWriter() throws Exception {
        FieldWriterChecksum<Integer> writer = new FieldWriterChecksum<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeChecksumField(42, dataWriter);

        assertEquals(42, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeChecksumFieldWithOptions() throws Exception {
        FieldWriterChecksum<String> writer = new FieldWriterChecksum<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("checksumField");
        writer.writeChecksumField("checksum_value", dataWriter, option);

        assertEquals("checksum_value", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeChecksumFieldWithNullValue() throws Exception {
        FieldWriterChecksum<String> writer = new FieldWriterChecksum<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeChecksumField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeChecksumFieldMultipleCalls() throws Exception {
        FieldWriterChecksum<Integer> writer = new FieldWriterChecksum<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeChecksumField(123, dataWriter);
        writer.writeChecksumField(456, dataWriter);

        assertEquals(456, dataWriter.getLastValue()); // Should be last written value
        assertEquals(2, dataWriter.getWriteCount());
    }

    @Test
    void writeChecksumFieldWithByteArray() throws Exception {
        FieldWriterChecksum<byte[]> writer = new FieldWriterChecksum<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] checksum = {(byte) 0xAB, (byte) 0xCD, (byte) 0xEF, 0x12};
        writer.writeChecksumField(checksum, dataWriter);

        assertArrayEquals(checksum, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeChecksumFieldWithLongValue() throws Exception {
        FieldWriterChecksum<Long> writer = new FieldWriterChecksum<>();
        SimpleDataWriter<Long> dataWriter = new SimpleDataWriter<>();

        Long checksumValue = 0xDEADBEEFL;
        writer.writeChecksumField(checksumValue, dataWriter);

        assertEquals(checksumValue, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeChecksumFieldWithBooleanType() throws Exception {
        FieldWriterChecksum<Boolean> writer = new FieldWriterChecksum<>();
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        writer.writeChecksumField(true, dataWriter);

        assertTrue(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeChecksumFieldWithMultipleOptions() throws Exception {
        FieldWriterChecksum<Integer> writer = new FieldWriterChecksum<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        WithOption[] options = {
            WithOption.WithName("checksumField"),
            WithOption.WithRenderAsList(false)
        };
        writer.writeChecksumField(999, dataWriter, options);

        assertEquals(999, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }
}