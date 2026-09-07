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

class FieldReaderUnknownTest {

    @Test
    void readUnknownFieldDelegatesToDataReader() throws Exception {
        FieldReaderUnknown<String> reader = new FieldReaderUnknown<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "unknown_value");

        String result = reader.readUnknownField(dataReader);

        assertEquals("unknown_value", result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readUnknownFieldWithOptions() throws Exception {
        FieldReaderUnknown<Integer> reader = new FieldReaderUnknown<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 999);

        WithOption option = WithOption.WithName("unknownField");
        Integer result = reader.readUnknownField(dataReader, option);

        assertEquals(999, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readUnknownFieldWithNullValue() throws Exception {
        FieldReaderUnknown<Object> reader = new FieldReaderUnknown<>();
        SimpleDataReader<Object> dataReader = new SimpleDataReader<>(() -> null);

        Object result = reader.readUnknownField(dataReader);

        assertNull(result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readUnknownFieldMultipleCalls() throws Exception {
        FieldReaderUnknown<Boolean> reader = new FieldReaderUnknown<>();
        SimpleDataReader<Boolean> dataReader = new SimpleDataReader<>(() -> false);

        Boolean result1 = reader.readUnknownField(dataReader);
        Boolean result2 = reader.readUnknownField(dataReader);

        assertFalse(result1);
        assertFalse(result2);
        assertEquals(2, dataReader.getReadCount());
    }

    @Test
    void readUnknownFieldWithComplexType() throws Exception {
        FieldReaderUnknown<byte[]> reader = new FieldReaderUnknown<>();
        byte[] unknownData = {0x41, 0x42, 0x43};
        SimpleDataReader<byte[]> dataReader = new SimpleDataReader<>(() -> unknownData);

        byte[] result = reader.readUnknownField(dataReader);

        assertArrayEquals(unknownData, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readUnknownFieldWithVaryingValues() throws Exception {
        FieldReaderUnknown<String> reader = new FieldReaderUnknown<>();
        int[] callCount = {0};
        String[] values = {"first", "second", "third"};

        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> {
            int index = callCount[0] % values.length;
            callCount[0]++;
            return values[index];
        });

        String result1 = reader.readUnknownField(dataReader);
        String result2 = reader.readUnknownField(dataReader);
        String result3 = reader.readUnknownField(dataReader);

        assertEquals("first", result1);
        assertEquals("second", result2);
        assertEquals("third", result3);
        assertEquals(3, dataReader.getReadCount());
    }

    @Test
    void readUnknownFieldWithNumericTypes() throws Exception {
        FieldReaderUnknown<Double> reader = new FieldReaderUnknown<>();
        SimpleDataReader<Double> dataReader = new SimpleDataReader<>(() -> 3.141592653589793);

        Double result = reader.readUnknownField(dataReader);

        assertEquals(3.141592653589793, result, 1e-15);
        assertEquals(1, dataReader.getReadCount());
    }
}