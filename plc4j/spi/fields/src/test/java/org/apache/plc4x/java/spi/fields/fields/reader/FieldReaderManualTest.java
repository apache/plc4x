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
import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyReadBuffer;
import org.apache.plc4x.java.spi.fields.data.reader.ParseSupplier;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FieldReaderManualTest {

    @Test
    void readManualFieldExecutesParseFunction() throws Exception {
        FieldReaderManual<String> reader = new FieldReaderManual<>();
        DummyReadBuffer readBuffer = new DummyReadBuffer();
        ParseSupplier<String> parseFunction = () -> "parsed_value";

        String result = reader.readManualField(readBuffer, parseFunction);

        assertEquals("parsed_value", result);
    }

    @Test
    void readManualFieldWithOptions() throws Exception {
        FieldReaderManual<Integer> reader = new FieldReaderManual<>();
        DummyReadBuffer readBuffer = new DummyReadBuffer();
        ParseSupplier<Integer> parseFunction = () -> 42;

        WithOption option = WithOption.WithName("manualField");
        Integer result = reader.readManualField(readBuffer, parseFunction, option);

        assertEquals(42, result);
    }

    @Test
    void readManualFieldManagesContext() throws Exception {
        FieldReaderManual<String> reader = new FieldReaderManual<>();
        DummyReadBuffer readBuffer = new DummyReadBuffer();

        final boolean[] contextPushed = {false};
        final boolean[] contextPopped = {false};

        ParseSupplier<String> parseFunction = () -> {
            // Check that context was pushed before parse function execution
            contextPushed[0] = !readBuffer.context.isEmpty();
            return "test";
        };

        String result = reader.readManualField(readBuffer, parseFunction);

        assertEquals("test", result);
        assertTrue(contextPushed[0], "Context should have been pushed before parse function");
        assertTrue(readBuffer.context.isEmpty(), "Context should be popped after parse function");
    }

    @Test
    void readManualFieldWithException() {
        FieldReaderManual<String> reader = new FieldReaderManual<>();
        DummyReadBuffer readBuffer = new DummyReadBuffer();
        ParseSupplier<String> faultyParseFunction = () -> {
            throw new RuntimeException("Parse error");
        };

        assertThrows(RuntimeException.class, () ->
            reader.readManualField(readBuffer, faultyParseFunction)
        );
    }

    @Test
    void readManualFieldReturnsNullFromParseFunction() throws Exception {
        FieldReaderManual<String> reader = new FieldReaderManual<>();
        DummyReadBuffer readBuffer = new DummyReadBuffer();
        ParseSupplier<String> parseFunction = () -> null;

        String result = reader.readManualField(readBuffer, parseFunction);

        assertNull(result);
    }

    @Test
    void readManualFieldWithComplexParseLogic() throws Exception {
        FieldReaderManual<Integer> reader = new FieldReaderManual<>();
        DummyReadBuffer readBuffer = new DummyReadBuffer();

        ParseSupplier<Integer> complexParseFunction = () -> {
            // Simulate complex parsing logic
            int sum = 0;
            for (int i = 1; i <= 10; i++) {
                sum += i;
            }
            return sum;
        };

        Integer result = reader.readManualField(readBuffer, complexParseFunction);

        assertEquals(55, result); // Sum of 1 to 10 is 55
    }

    @Test
    void readManualFieldMultipleCalls() throws Exception {
        FieldReaderManual<String> reader = new FieldReaderManual<>();
        DummyReadBuffer readBuffer = new DummyReadBuffer();

        int[] counter = {0};
        ParseSupplier<String> parseFunction = () -> "call_" + (++counter[0]);

        String result1 = reader.readManualField(readBuffer, parseFunction);
        String result2 = reader.readManualField(readBuffer, parseFunction);

        assertEquals("call_1", result1);
        assertEquals("call_2", result2);
    }
}