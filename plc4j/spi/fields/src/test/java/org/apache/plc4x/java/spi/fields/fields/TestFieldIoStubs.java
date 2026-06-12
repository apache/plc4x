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

package org.apache.plc4x.java.spi.fields.fields;

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.TestBuffers;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.apache.plc4x.java.spi.fields.data.writer.DataWriter;

import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Supplier;

/**
 * Minimal stub implementations for DataReader/DataWriter used in field reader/writer tests.
 */
public class TestFieldIoStubs {

    public static class SimpleDataReader<T> implements DataReader<T> {
        private final TestBuffers.DummyReadBuffer rb = new TestBuffers.DummyReadBuffer();
        private final Supplier<T> supplier;
        private final AtomicInteger readCount = new AtomicInteger();

        public SimpleDataReader(Supplier<T> supplier) {
            this.supplier = supplier;
        }

        @Override
        public T read(WithOption... options) throws BufferException {
            readCount.incrementAndGet();
            return supplier.get();
        }

        @Override
        public ReadBuffer getReadBuffer() {
            return rb;
        }

        @Override
        public int getPositionInBits() {
            return rb.getPositionInBits();
        }

        @Override
        public void setPositionInBits(int positionInBits) {
            rb.setPositionInBits(positionInBits);
        }

        @Override
        public void pushContext(WithOption... options) throws BufferException {
            rb.pushContext(options);
        }

        @Override
        public void popContext(WithOption... options) throws BufferException {
            rb.popContext(options);
        }

        public int getReadCount() { return readCount.get(); }
    }

    public static class SimpleDataWriter<T> implements DataWriter<T> {
        private final TestBuffers.DummyWriteBuffer wb = new TestBuffers.DummyWriteBuffer();
        private T lastValue;
        private final AtomicInteger writeCount = new AtomicInteger();

        @Override
        public void write(T value, WithOption... options) throws BufferException {
            this.lastValue = value;
            writeCount.incrementAndGet();
        }

        @Override
        public WriteBuffer getWriteBuffer() {
            return wb;
        }

        public T getLastValue() { return lastValue; }
        public int getWriteCount() { return writeCount.get(); }


        @Override
        public void pushContext(WithOption... options) throws BufferException { wb.pushContext(options); }

        @Override
        public void popContext(WithOption... options) throws BufferException { wb.popContext(options); }
    }
}
