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
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.TestBuffers;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataReader;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FieldReaderPaddingTest {

    @Test
    void readPaddingFieldReadsSpecifiedTimes() throws Exception {
        FieldReaderPadding<Byte> reader = new FieldReaderPadding<>();
        SimpleDataReader<Byte> dataReader = new SimpleDataReader<>(() -> (byte) 0);

        reader.readPaddingField(dataReader, 3);

        assertEquals(3, dataReader.getReadCount());
    }

    @Test
    void readPaddingFieldWithZeroPadding() throws Exception {
        FieldReaderPadding<Byte> reader = new FieldReaderPadding<>();
        SimpleDataReader<Byte> dataReader = new SimpleDataReader<>(() -> (byte) 0);

        reader.readPaddingField(dataReader, 0);

        assertEquals(0, dataReader.getReadCount());
    }

    @Test
    void readPaddingFieldWithNegativePadding() throws Exception {
        FieldReaderPadding<Byte> reader = new FieldReaderPadding<>();
        SimpleDataReader<Byte> dataReader = new SimpleDataReader<>(() -> (byte) 0);

        reader.readPaddingField(dataReader, -1);

        assertEquals(0, dataReader.getReadCount());
    }

    @Test
    void readPaddingFieldIgnoresExceptions() throws Exception {
        FieldReaderPadding<Byte> reader = new FieldReaderPadding<>();

        int[] callCount = {0};
        DataReader<Byte> faultyReader = new DataReader<Byte>() {
            @Override
            public Byte read(WithOption... options) throws BufferException {
                callCount[0]++;
                if (callCount[0] <= 2) {
                    return (byte) 0;
                }
                throw new BufferException("Simulated read error");
            }

            @Override
            public ReadBuffer getReadBuffer() {
                return new TestBuffers.DummyReadBuffer();
            }

            @Override
            public int getPositionInBits() {
                return 0;
            }

            @Override
            public void setPositionInBits(int positionInBits) {}

            @Override
            public void pushContext(WithOption... options) throws BufferException {}

            @Override
            public void popContext(WithOption... options) throws BufferException {}
        };

        // Should not throw exception even though the reader fails on the 3rd call
        assertDoesNotThrow(() -> reader.readPaddingField(faultyReader, 5));

        assertEquals(5, callCount[0]); // All 5 calls should be attempted
    }

    @Test
    void readPaddingFieldWithOptions() throws Exception {
        FieldReaderPadding<Integer> reader = new FieldReaderPadding<>();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 0);

        WithOption option = WithOption.WithName("paddingField");
        reader.readPaddingField(dataReader, 2, option);

        assertEquals(2, dataReader.getReadCount());
    }

    @Test
    void readPaddingFieldManagesContext() throws Exception {
        FieldReaderPadding<Byte> reader = new FieldReaderPadding<>();
        SimpleDataReader<Byte> dataReader = new SimpleDataReader<>(() -> (byte) 0);

        // Context should be empty initially
        assertTrue(dataReader.getReadBuffer().getContext() == null ||
                  dataReader.getReadBuffer().getContext().length == 0);

        reader.readPaddingField(dataReader, 1);

        // Context should be empty after reading (pushed and popped)
        assertTrue(dataReader.getReadBuffer().getContext() == null ||
                  dataReader.getReadBuffer().getContext().length == 0);
    }

    @Test
    void readPaddingFieldIgnoresRuntimeExceptions() throws Exception {
        FieldReaderPadding<String> reader = new FieldReaderPadding<>();

        int[] callCount = {0};
        DataReader<String> faultyReader = new DataReader<String>() {
            @Override
            public String read(WithOption... options) throws BufferException {
                callCount[0]++;
                if (callCount[0] == 2) {
                    throw new RuntimeException("Runtime error");
                }
                return "padding";
            }

            @Override
            public ReadBuffer getReadBuffer() {
                return new TestBuffers.DummyReadBuffer();
            }

            @Override
            public int getPositionInBits() {
                return 0;
            }

            @Override
            public void setPositionInBits(int positionInBits) {}

            @Override
            public void pushContext(WithOption... options) throws BufferException {}

            @Override
            public void popContext(WithOption... options) throws BufferException {}
        };

        // Should not throw exception even though the reader fails with RuntimeException
        assertDoesNotThrow(() -> reader.readPaddingField(faultyReader, 3));

        assertEquals(3, callCount[0]); // All 3 calls should be attempted
    }
}