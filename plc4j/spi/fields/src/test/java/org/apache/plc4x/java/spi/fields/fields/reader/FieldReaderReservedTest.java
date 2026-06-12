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

class FieldReaderReservedTest {

    @Test
    void readReservedFieldMatchingValueReturnsValue() throws Exception {
        FieldReaderReserved<Integer> reader = new FieldReaderReserved<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 42);

        Integer result = reader.readReservedField(dataReader, 42);

        assertEquals(42, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readReservedFieldMismatchLogsInfoButReturnsReadValue() throws Exception {
        FieldReaderReserved<Integer> reader = new FieldReaderReserved<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 99);

        Integer result = reader.readReservedField(dataReader, 42);

        // Should return the read value, not the reference value
        assertEquals(99, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readReservedFieldWithOptions() throws Exception {
        FieldReaderReserved<String> reader = new FieldReaderReserved<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "reserved");

        WithOption option = WithOption.WithName("reservedField");
        String result = reader.readReservedField(dataReader, "reserved", option);

        assertEquals("reserved", result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readReservedFieldUnnamedGetsReservedName() throws Exception {
        FieldReaderReserved<Integer> reader = new FieldReaderReserved<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 0);

        WithOption unnamedOption = WithOption.WithName("unnamed");
        Integer result = reader.readReservedField(dataReader, 0, unnamedOption);

        assertEquals(0, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readReservedFieldUnnamedCaseInsensitive() throws Exception {
        FieldReaderReserved<Integer> reader = new FieldReaderReserved<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 123);

        WithOption unnamedOption = WithOption.WithName("UNNAMED");
        Integer result = reader.readReservedField(dataReader, 123, unnamedOption);

        assertEquals(123, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readReservedFieldWithNullValues() throws Exception {
        FieldReaderReserved<String> reader = new FieldReaderReserved<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> null);

        String result = reader.readReservedField(dataReader, null);

        assertNull(result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readReservedFieldNullReadValueMismatch() throws Exception {
        FieldReaderReserved<String> reader = new FieldReaderReserved<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> null);

        String result = reader.readReservedField(dataReader, "expected");

        // Should return the read value (null), not the reference value
        assertNull(result);
        assertEquals(1, dataReader.getReadCount());
    }
}