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

class FieldReaderEnumTest {

    enum TestEnum {
        VALUE_ONE,
        VALUE_TWO
    }

    @Test
    void readEnumFieldDelegatesToDataReader() throws Exception {
        FieldReaderEnum<TestEnum> reader = new FieldReaderEnum<>();
        SimpleDataReader<TestEnum> dataReader = new SimpleDataReader<>(() -> TestEnum.VALUE_ONE);

        TestEnum result = reader.readEnumField(dataReader);

        assertEquals(TestEnum.VALUE_ONE, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readEnumFieldWithOptions() throws Exception {
        FieldReaderEnum<TestEnum> reader = new FieldReaderEnum<>();
        SimpleDataReader<TestEnum> dataReader = new SimpleDataReader<>(() -> TestEnum.VALUE_TWO);

        WithOption option = WithOption.WithName("enumField");
        TestEnum result = reader.readEnumField(dataReader, option);

        assertEquals(TestEnum.VALUE_TWO, result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readEnumFieldManagesContext() throws Exception {
        FieldReaderEnum<String> reader = new FieldReaderEnum<>();
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "ENUM_VALUE");

        // Context should be empty initially
        assertTrue(dataReader.getReadBuffer().getContext() == null ||
                  dataReader.getReadBuffer().getContext().length == 0);

        String result = reader.readEnumField(dataReader);

        assertEquals("ENUM_VALUE", result);

        // Context should be empty after reading (pushed and popped)
        assertTrue(dataReader.getReadBuffer().getContext() == null ||
                  dataReader.getReadBuffer().getContext().length == 0);
    }

    @Test
    void readEnumFieldWithNullValue() throws Exception {
        FieldReaderEnum<TestEnum> reader = new FieldReaderEnum<>();
        SimpleDataReader<TestEnum> dataReader = new SimpleDataReader<>(() -> null);

        TestEnum result = reader.readEnumField(dataReader);

        assertNull(result);
        assertEquals(1, dataReader.getReadCount());
    }

    @Test
    void readEnumFieldMultipleCalls() throws Exception {
        FieldReaderEnum<TestEnum> reader = new FieldReaderEnum<>();
        SimpleDataReader<TestEnum> dataReader = new SimpleDataReader<>(() -> TestEnum.VALUE_ONE);

        TestEnum result1 = reader.readEnumField(dataReader);
        TestEnum result2 = reader.readEnumField(dataReader);

        assertEquals(TestEnum.VALUE_ONE, result1);
        assertEquals(TestEnum.VALUE_ONE, result2);
        assertEquals(2, dataReader.getReadCount());
    }
}