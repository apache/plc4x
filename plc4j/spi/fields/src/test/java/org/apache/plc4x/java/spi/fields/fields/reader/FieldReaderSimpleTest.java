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

class FieldReaderSimpleTest {

    @Test
    void readSimpleFieldDelegatesToDataReader() throws Exception {
        FieldReaderSimple<String> reader = new FieldReaderSimple<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "test value");

        String result = reader.readSimpleField(dataReader);

        assertEquals("test value", result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readSimpleFieldWithOptions() throws Exception {
        FieldReaderSimple<Integer> reader = new FieldReaderSimple<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 42);

        WithOption option = WithOption.WithName("testField");
        Integer result = reader.readSimpleField(dataReader, option);

        assertEquals(42, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readSimpleFieldWithNullValue() throws Exception {
        FieldReaderSimple<String> reader = new FieldReaderSimple<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> null);

        String result = reader.readSimpleField(dataReader);

        assertNull(result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readSimpleFieldMultipleCalls() throws Exception {
        FieldReaderSimple<Integer> reader = new FieldReaderSimple<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 123);

        Integer result1 = reader.readSimpleField(dataReader);
        Integer result2 = reader.readSimpleField(dataReader);

        assertEquals(123, result1);
        assertEquals(123, result2);
        assertEquals(2, dataReader.getReadCount());
    }
}