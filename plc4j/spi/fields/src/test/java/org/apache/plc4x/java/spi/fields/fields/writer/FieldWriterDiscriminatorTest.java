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

class FieldWriterDiscriminatorTest {

    @Test
    void writeDiscriminatorFieldDelegatesToDataWriter() throws Exception {
        FieldWriterDiscriminator<Integer> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        writer.writeDiscriminatorField(123, dataWriter);

        assertEquals(123, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldWithOptions() throws Exception {
        FieldWriterDiscriminator<String> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        WithOption option = WithOption.WithName("discriminatorField");
        writer.writeDiscriminatorField("discriminator_value", dataWriter, option);

        assertEquals("discriminator_value", dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldWithNullValue() throws Exception {
        FieldWriterDiscriminator<String> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<String> dataWriter = new SimpleDataWriter<>();

        writer.writeDiscriminatorField(null, dataWriter);

        assertNull(dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldMultipleCalls() throws Exception {
        FieldWriterDiscriminator<Boolean> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<Boolean> dataWriter = new SimpleDataWriter<>();

        writer.writeDiscriminatorField(true, dataWriter);
        writer.writeDiscriminatorField(false, dataWriter);

        assertFalse(dataWriter.getLastValue());
        assertEquals(2, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldWithComplexType() throws Exception {
        FieldWriterDiscriminator<byte[]> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<byte[]> dataWriter = new SimpleDataWriter<>();

        byte[] discriminatorData = {1, 2, 3, 4};
        writer.writeDiscriminatorField(discriminatorData, dataWriter);

        assertArrayEquals(discriminatorData, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldWithEnumValue() throws Exception {
        FieldWriterDiscriminator<TestDiscriminator> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<TestDiscriminator> dataWriter = new SimpleDataWriter<>();

        writer.writeDiscriminatorField(TestDiscriminator.TYPE_A, dataWriter);

        assertEquals(TestDiscriminator.TYPE_A, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldWithNumericValues() throws Exception {
        FieldWriterDiscriminator<Long> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<Long> dataWriter = new SimpleDataWriter<>();

        long discriminatorValue = 0x12345678L;
        writer.writeDiscriminatorField(discriminatorValue, dataWriter);

        assertEquals(discriminatorValue, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldWithFloatingPoint() throws Exception {
        FieldWriterDiscriminator<Double> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<Double> dataWriter = new SimpleDataWriter<>();

        double discriminatorValue = 3.141592653589793;
        writer.writeDiscriminatorField(discriminatorValue, dataWriter);

        assertEquals(discriminatorValue, dataWriter.getLastValue(), 1e-15);
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldWithMultipleOptions() throws Exception {
        FieldWriterDiscriminator<Integer> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        WithOption[] options = {
            WithOption.WithName("discriminatorField"),
            WithOption.WithRenderAsList(true)
        };
        writer.writeDiscriminatorField(999, dataWriter, options);

        assertEquals(999, dataWriter.getLastValue());
        assertEquals(1, dataWriter.getWriteCount());
    }

    @Test
    void writeDiscriminatorFieldSequentialValues() throws Exception {
        FieldWriterDiscriminator<Integer> writer = new FieldWriterDiscriminator<>();
        SimpleDataWriter<Integer> dataWriter = new SimpleDataWriter<>();

        // Test writing different discriminator values sequentially
        writer.writeDiscriminatorField(1, dataWriter);
        writer.writeDiscriminatorField(2, dataWriter);
        writer.writeDiscriminatorField(3, dataWriter);

        assertEquals(3, dataWriter.getLastValue()); // Last value written
        assertEquals(3, dataWriter.getWriteCount()); // Total writes
    }

    // Test enum for discriminator values
    enum TestDiscriminator {
        TYPE_A,
        TYPE_B,
        TYPE_C
    }
}