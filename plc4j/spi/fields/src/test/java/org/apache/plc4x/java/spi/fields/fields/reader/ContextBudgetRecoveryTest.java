/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferUnderflowException;
import org.apache.plc4x.java.spi.fields.data.TestBuffers;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.apache.plc4x.java.spi.fields.exceptions.ParseAssertException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * A complex field opens a context on the way into it and closes it on the way out, so a read that
 * fails part-way leaves its context open. That costs nothing when the failure ends the parse, but
 * an optional or peeked field recovers and carries on - so it has to give back the nesting budget
 * the failed read took, the same way it gives back the read position.
 */
class ContextBudgetRecoveryTest {

    /** A buffer that counts its open contexts, which is all these readers need to be judged on. */
    private static class ContextCountingReadBuffer extends TestBuffers.DummyReadBuffer {

        private int depth;

        @Override
        public void pushContext(WithOption... options) {
            depth++;
        }

        @Override
        public void popContext(WithOption... options) {
            if (depth > 0) {
                depth--;
            }
        }

        @Override
        public int getContextDepth() {
            return depth;
        }

        @Override
        public void resetContextDepth(int contextDepth) {
            depth = Math.min(depth, contextDepth);
        }
    }

    /**
     * Stands in for a complex field's reader: opens a context the way a generated parser does,
     * then fails before it can close it again.
     */
    private static class FailingComplexReader implements DataReader<Object> {

        private final ReadBuffer readBuffer;
        private final BufferException failure;

        FailingComplexReader(ReadBuffer readBuffer, BufferException failure) {
            this.readBuffer = readBuffer;
            this.failure = failure;
        }

        @Override
        public Object read(WithOption... options) throws BufferException {
            readBuffer.pushContext(WithOption.WithName("someComplexType"));
            throw failure;
        }

        @Override
        public ReadBuffer getReadBuffer() {
            return readBuffer;
        }

        @Override
        public int getPositionInBits() {
            return readBuffer.getPositionInBits();
        }

        @Override
        public void setPositionInBits(int positionInBits) {
            readBuffer.setPositionInBits(positionInBits);
        }

        @Override
        public void pushContext(WithOption... options) throws BufferException {
            readBuffer.pushContext(options);
        }

        @Override
        public void popContext(WithOption... options) throws BufferException {
            readBuffer.popContext(options);
        }
    }

    @Test
    void anOptionalThatRanOutOfDataGivesBackTheNestingBudget() throws BufferException {
        ContextCountingReadBuffer readBuffer = new ContextCountingReadBuffer();
        for (int i = 0; i < 100; i++) {
            new FieldReaderOptional<>().readOptionalField(
                new FailingComplexReader(readBuffer, new BufferUnderflowException("ran out")), true);
        }
        assertEquals(0, readBuffer.getContextDepth(),
            "a hundred absent optional fields must not spend a hundred levels of nesting budget");
    }

    @Test
    void anOptionalWhoseAssertionFailedGivesBackTheNestingBudget() throws BufferException {
        ContextCountingReadBuffer readBuffer = new ContextCountingReadBuffer();
        for (int i = 0; i < 100; i++) {
            new FieldReaderOptional<>().readOptionalField(
                new FailingComplexReader(readBuffer, new ParseAssertException("no match")), true);
        }
        assertEquals(0, readBuffer.getContextDepth());
    }

    @Test
    void aFailedPeekGivesBackTheNestingBudget() throws BufferException {
        ContextCountingReadBuffer readBuffer = new ContextCountingReadBuffer();
        for (int i = 0; i < 100; i++) {
            new FieldReaderPeek<>().readPeekField(
                new FailingComplexReader(readBuffer, new ParseAssertException("no match")), 0);
        }
        assertEquals(0, readBuffer.getContextDepth(),
            "a hundred failed peeks must not spend a hundred levels of nesting budget");
    }
}
