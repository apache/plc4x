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

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.fields.data.TestBuffers;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataReader;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

class FieldReaderArrayTest {

    @Test
    void readArrayByCountReadsThatManyItems() throws Exception {
        AtomicInteger i = new AtomicInteger();
        SimpleDataReader<Integer> dr = new SimpleDataReader<>(i::getAndIncrement);
        FieldReaderArray<Integer> reader = new FieldReaderArray<>();
        List<Integer> list = reader.readArrayFieldCount(dr, 3);
        assertEquals(List.of(0,1,2), list);
        assertEquals(3, dr.getReadCount());
    }

    @Test
    void readArrayByCountNegativeReturnsNull() throws Exception {
        SimpleDataReader<Integer> dr = new SimpleDataReader<>(() -> 1);
        FieldReaderArray<Integer> reader = new FieldReaderArray<>();
        assertNull(reader.readArrayFieldCount(dr, -1));
    }

    @Test
    void readArrayCountExceedsMaxThrows() {
        SimpleDataReader<Integer> dr = new SimpleDataReader<>(() -> 1);
        FieldReaderArray<Integer> reader = new FieldReaderArray<>();
        assertThrows(BufferException.class, () -> reader.readArrayFieldCount(dr, (long)Integer.MAX_VALUE + 1));
    }

    @Test
    void readArrayByLengthConsumesUntilLength() throws Exception {
        AtomicInteger i = new AtomicInteger();
        // Custom DataReader that advances position by 8 bits per read
        DataReader<Integer> dr = new DataReader<>() {
            final TestBuffers.DummyReadBuffer rb = new TestBuffers.DummyReadBuffer();
            @Override public Integer read(WithOption... options) { rb.setPositionInBits(rb.getPositionInBits() + 8); return i.getAndIncrement(); }
            @Override public ReadBuffer getReadBuffer() { return rb; }
            @Override public int getPositionInBits() { return rb.getPositionInBits(); }
            @Override public void setPositionInBits(int positionInBits) { rb.setPositionInBits(positionInBits); }
            @Override public void pushContext(WithOption... options) {}
            @Override public void popContext(WithOption... options) {}
        };
        FieldReaderArray<Integer> reader = new FieldReaderArray<>();
        // length in bytes = 3, with 1 byte per element -> 3 elements
        List<Integer> list = reader.readArrayFieldLength(dr, 3);
        assertEquals(List.of(0,1,2), list);
    }

    @Test
    void readArrayByTerminationStopsOnCondition() throws Exception {
        AtomicInteger i = new AtomicInteger();
        SimpleDataReader<Integer> dr = new SimpleDataReader<>(i::getAndIncrement);
        FieldReaderArray<Integer> reader = new FieldReaderArray<>();
        List<Integer> list = reader.readArrayFieldTerminated(dr, () -> i.get() >= 3);
        assertEquals(List.of(0,1,2), list);
    }
}
