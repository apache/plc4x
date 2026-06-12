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

package org.apache.plc4x.java.spi.fields.data.writer;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyWriteBuffer;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;

import static org.junit.jupiter.api.Assertions.*;

class DataWriterSimpleBigDecimalTest {

    @Test
    void writeBigDecimalDelegatesToWriteBuffer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterSimpleBigDecimal writer = new DataWriterSimpleBigDecimal(wb, 128);

        BigDecimal testValue = new BigDecimal("123.456789");
        writer.write(testValue);

        assertEquals(testValue, wb.bigDecimalWritten);
    }

    @Test
    void writeBigDecimalWithOptions() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterSimpleBigDecimal writer = new DataWriterSimpleBigDecimal(wb, 128);

        BigDecimal testValue = new BigDecimal("999.999");
        WithOption option = WithOption.WithName("testDecimal");
        writer.write(testValue, option);

        assertEquals(testValue, wb.bigDecimalWritten);
    }

    @Test
    void contextMethodsDelegateToWriteBuffer() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterSimpleBigDecimal writer = new DataWriterSimpleBigDecimal(wb, 128);

        WithOption option = WithOption.WithName("test");

        assertTrue(wb.context.isEmpty());
        writer.pushContext(option);
        assertFalse(wb.context.isEmpty());
        writer.popContext(option);
        assertTrue(wb.context.isEmpty());
    }

    @Test
    void getWriteBufferReturnsOriginalBuffer() {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterSimpleBigDecimal writer = new DataWriterSimpleBigDecimal(wb, 128);

        assertSame(wb, writer.getWriteBuffer());
    }

    @Test
    void bitLengthIsStored() {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterSimpleBigDecimal writer = new DataWriterSimpleBigDecimal(wb, 256);

        assertEquals(256, writer.bitLength);
    }

    @Test
    void nullValueHandling() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterSimpleBigDecimal writer = new DataWriterSimpleBigDecimal(wb, 128);

        writer.write(null);

        assertNull(wb.bigDecimalWritten);
    }
}