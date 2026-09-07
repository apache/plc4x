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

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.fields.data.TestBuffers;
import org.apache.plc4x.java.spi.fields.exceptions.ParseAssertException;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataReader;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FieldReaderPeekTest {

    @Test
    void peekSuccessReturnsValueAndRestoresPosition() throws Exception {
        FieldReaderPeek<Integer> reader = new FieldReaderPeek<>();
        SimpleDataReader<Integer> dr = new SimpleDataReader<>(() -> 5);
        dr.setPositionInBits(10);
        Integer val = reader.readPeekField(dr, 0);
        assertEquals(5, val);
        assertEquals(10, dr.getPositionInBits());
    }

    @Test
    void peekFailureReturnsNullAndRestoresPosition() throws Exception {
        FieldReaderPeek<Integer> reader = new FieldReaderPeek<>();
        DataReader<Integer> dr = new DataReader<>() {
            final TestBuffers.DummyReadBuffer rb = new TestBuffers.DummyReadBuffer();
            @Override public Integer read(WithOption... options) throws org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException { rb.setPositionInBits(123); throw new ParseAssertException("fail"); }
            @Override public ReadBuffer getReadBuffer() { return rb; }
            @Override public int getPositionInBits() { return rb.getPositionInBits(); }
            @Override public void setPositionInBits(int positionInBits) { rb.setPositionInBits(positionInBits); }
            @Override public void pushContext(WithOption... options) {}
            @Override public void popContext(WithOption... options) {}
        };
        int pos = dr.getPositionInBits();
        Integer val = reader.readPeekField(dr, 0);
        assertNull(val);
        assertEquals(pos, dr.getPositionInBits());
    }
}
