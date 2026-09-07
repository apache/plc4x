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
package org.apache.plc4x.java.utils.testutils.utils.migration;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Method;
import java.util.*;

import static org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption.WithByteOrder;
import static org.junit.jupiter.api.Assertions.*;

class MessageResolverTest {

    @Test
    void testGetMessageInputWithInvalidPackage() {
        Map<String, String> options = new HashMap<>();
        options.put("package", "invalid.package.name");

        List<String> parserArgs = new ArrayList<>();

        assertThrows(DriverTestsuiteException.class, () -> {
            MessageResolver.getMessageInput(options, "TestMessage", parserArgs);
        });
    }

    @Test
    void testGetMessageInputWithMissingProtocol() {
        Map<String, String> options = new HashMap<>();
        // No package, protocolName, or outputFlavor

        List<String> parserArgs = new ArrayList<>();

        assertThrows(Exception.class, () -> {
            MessageResolver.getMessageInput(options, "TestMessage", parserArgs);
        });
    }

    @Test
    void testGetMessageIOStaticLinkedWithInvalidOptions() {
        Map<String, String> options = new HashMap<>();
        options.put("protocolName", "nonexistent");
        options.put("outputFlavor", "readwrite");

        List<String> parserArgs = new ArrayList<>();

        assertThrows(Exception.class, () -> {
            MessageResolver.getMessageIOStaticLinked(options, "TestMessage", parserArgs);
        });
    }

    // Tests for DeferredMessageInput inner class

    public static class TestMessage implements Message {
        public final boolean boolVal;
        public final byte byteVal;
        public final short shortVal;
        public final int intVal;
        public final long longVal;
        public final float floatVal;
        public final double doubleVal;
        public final String stringVal;

        public TestMessage(boolean boolVal, byte byteVal, short shortVal, int intVal,
                          long longVal, float floatVal, double doubleVal, String stringVal) {
            this.boolVal = boolVal;
            this.byteVal = byteVal;
            this.shortVal = shortVal;
            this.intVal = intVal;
            this.longVal = longVal;
            this.floatVal = floatVal;
            this.doubleVal = doubleVal;
            this.stringVal = stringVal;
        }

        public static TestMessage staticParse(ReadBuffer io, boolean b, byte by, short s,
                                             int i, long l, float f, double d, String str) {
            return new TestMessage(b, by, s, i, l, f, d, str);
        }

        @Override
        public void serialize(org.apache.plc4x.java.spi.buffers.api.WriteBuffer buffer) throws BufferException {}

        @Override
        public int getLengthInBytes() { return 0; }

        @Override
        public int getLengthInBits() { return 0; }
    }

    public enum TestEnum {
        VALUE1, VALUE2
    }

    public static class EnumTestMessage implements Message {
        public final TestEnum enumVal;

        public EnumTestMessage(TestEnum enumVal) {
            this.enumVal = enumVal;
        }

        public static EnumTestMessage staticParse(ReadBuffer io, TestEnum e) {
            return new EnumTestMessage(e);
        }

        @Override
        public void serialize(org.apache.plc4x.java.spi.buffers.api.WriteBuffer buffer) throws BufferException {}

        @Override
        public int getLengthInBytes() { return 0; }

        @Override
        public int getLengthInBits() { return 0; }
    }

    @Test
    void testDeferredMessageInputWithAllParameterTypes() throws Exception {
        Method method = TestMessage.class.getMethod("staticParse", ReadBuffer.class,
            boolean.class, byte.class, short.class, int.class, long.class,
            float.class, double.class, String.class);

        List<String> args = Arrays.asList("true", "42", "1000", "50000", "1000000", "3.14", "2.718", "test");

        MessageResolver.DeferredMessageInput input = new MessageResolver.DeferredMessageInput(method, args);

        ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[]{}, WithByteOrder("BIG_ENDIAN"));
        Message result = input.parse(buffer);

        assertNotNull(result);
        assertInstanceOf(TestMessage.class, result);
        TestMessage msg = (TestMessage) result;
        assertTrue(msg.boolVal);
        assertEquals(42, msg.byteVal);
        assertEquals(1000, msg.shortVal);
        assertEquals(50000, msg.intVal);
        assertEquals(1000000L, msg.longVal);
        assertEquals(3.14f, msg.floatVal, 0.01);
        assertEquals(2.718, msg.doubleVal, 0.001);
        assertEquals("test", msg.stringVal);
    }

    @Test
    void testDeferredMessageInputWithEnumParameter() throws Exception {
        Method method = EnumTestMessage.class.getMethod("staticParse", ReadBuffer.class, TestEnum.class);

        List<String> args = Collections.singletonList("VALUE1");

        MessageResolver.DeferredMessageInput input = new MessageResolver.DeferredMessageInput(method, args);

        ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[]{}, WithByteOrder("BIG_ENDIAN"));
        Message result = input.parse(buffer);

        assertNotNull(result);
        assertInstanceOf(EnumTestMessage.class, result);
        assertEquals(TestEnum.VALUE1, ((EnumTestMessage) result).enumVal);
    }

    @Test
    void testDeferredMessageInputWithWrongParameterCount() throws Exception {
        Method method = TestMessage.class.getMethod("staticParse", ReadBuffer.class,
            boolean.class, byte.class, short.class, int.class, long.class,
            float.class, double.class, String.class);

        // Only 3 args instead of 8
        List<String> args = Arrays.asList("true", "42", "1000");

        MessageResolver.DeferredMessageInput input = new MessageResolver.DeferredMessageInput(method, args);

        ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[]{}, WithByteOrder("BIG_ENDIAN"));

        BufferException exception = assertThrows(BufferException.class, () -> {
            input.parse(buffer);
        });

        assertTrue(exception.getMessage().contains("Invalid parameters detected"));
        assertTrue(exception.getMessage().contains("expected 8, not 3"));
    }

    @Test
    void testDeferredMessageInputWithBoxedTypes() throws Exception {
        Method method = TestMessage.class.getMethod("staticParse", ReadBuffer.class,
            boolean.class, byte.class, short.class, int.class, long.class,
            float.class, double.class, String.class);

        List<String> args = Arrays.asList("false", "0", "0", "0", "0", "0.0", "0.0", "");

        MessageResolver.DeferredMessageInput input = new MessageResolver.DeferredMessageInput(method, args);

        ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[]{}, WithByteOrder("BIG_ENDIAN"));
        Message result = input.parse(buffer);

        assertNotNull(result);
    }

    public static class UnsupportedTypeMessage implements Message {
        public static UnsupportedTypeMessage staticParse(ReadBuffer io, Object unsupported) {
            return new UnsupportedTypeMessage();
        }

        @Override
        public void serialize(org.apache.plc4x.java.spi.buffers.api.WriteBuffer buffer) throws BufferException {}

        @Override
        public int getLengthInBytes() { return 0; }

        @Override
        public int getLengthInBits() { return 0; }
    }

    @Test
    void testDeferredMessageInputWithUnsupportedParameterType() throws Exception {
        Method method = UnsupportedTypeMessage.class.getMethod("staticParse", ReadBuffer.class, Object.class);

        List<String> args = Collections.singletonList("someValue");

        MessageResolver.DeferredMessageInput input = new MessageResolver.DeferredMessageInput(method, args);

        ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[]{}, WithByteOrder("BIG_ENDIAN"));

        BufferException exception = assertThrows(BufferException.class, () -> {
            input.parse(buffer);
        });

        assertTrue(exception.getMessage().contains("unsupported parameter type"));
    }

    public static class FailingMessage implements Message {
        public static FailingMessage staticParse(ReadBuffer io) {
            throw new RuntimeException("Parse intentionally failed");
        }

        @Override
        public void serialize(org.apache.plc4x.java.spi.buffers.api.WriteBuffer buffer) throws BufferException {}

        @Override
        public int getLengthInBytes() { return 0; }

        @Override
        public int getLengthInBits() { return 0; }
    }

    @Test
    void testDeferredMessageInputWithParseException() throws Exception {
        Method method = FailingMessage.class.getMethod("staticParse", ReadBuffer.class);

        List<String> args = Collections.emptyList();

        MessageResolver.DeferredMessageInput input = new MessageResolver.DeferredMessageInput(method, args);

        ReadBufferByteBased buffer = new ReadBufferByteBased(new byte[]{}, WithByteOrder("BIG_ENDIAN"));

        BufferException exception = assertThrows(BufferException.class, () -> {
            input.parse(buffer);
        });

        assertTrue(exception.getMessage().contains("Could not parse payload"));
    }
}
