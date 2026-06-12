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

package org.apache.plc4x.java.spi.fields.fields.reader;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataReader;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FieldReaderDiscriminatorTest {

    @Test
    void readDiscriminatorFieldDelegatesToDataReader() throws Exception {
        FieldReaderDiscriminator<Integer> reader = new FieldReaderDiscriminator<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 123);

        Integer result = reader.readDiscriminatorField(dataReader);

        assertEquals(123, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readDiscriminatorFieldWithOptions() throws Exception {
        FieldReaderDiscriminator<String> reader = new FieldReaderDiscriminator<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "discriminator_value");

        WithOption option = WithOption.WithName("discriminatorField");
        String result = reader.readDiscriminatorField(dataReader, option);

        assertEquals("discriminator_value", result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readDiscriminatorFieldWithNullValue() throws Exception {
        FieldReaderDiscriminator<String> reader = new FieldReaderDiscriminator<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> null);

        String result = reader.readDiscriminatorField(dataReader);

        assertNull(result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readDiscriminatorFieldMultipleCalls() throws Exception {
        FieldReaderDiscriminator<Boolean> reader = new FieldReaderDiscriminator<>();
        SimpleDataReader<Boolean> dataReader = new SimpleDataReader<>(() -> true);

        Boolean result1 = reader.readDiscriminatorField(dataReader);
        Boolean result2 = reader.readDiscriminatorField(dataReader);

        assertTrue(result1);
        assertTrue(result2);
        assertEquals(2, dataReader.getReadCount());
    }

    @Test
    void readDiscriminatorFieldWithComplexType() throws Exception {
        FieldReaderDiscriminator<byte[]> reader = new FieldReaderDiscriminator<>();
        byte[] expectedData = {1, 2, 3, 4};
        SimpleDataReader<byte[]> dataReader = new SimpleDataReader<>(() -> expectedData);

        byte[] result = reader.readDiscriminatorField(dataReader);

        assertArrayEquals(expectedData, result);
        assertEquals(1, dataReader.getReadCount());
    }
}