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
import org.apache.plc4x.java.spi.buffers.api.MessageInput;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareDocumentFactory;
import org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareElement;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.dom4j.QName;
import org.junit.jupiter.api.Test;

import java.net.URI;
import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

class MessageValidatorAndMigratorTest {

    public static class SimpleTestMessage implements Message {
        private final String value;

        public SimpleTestMessage(String value) {
            this.value = value;
        }

        @Override
        public void serialize(WriteBuffer buffer) throws BufferException {
            // Write the value as XML-style content
        }

        @Override
        public int getLengthInBytes() {
            return value.length();
        }

        @Override
        public int getLengthInBits() {
            return value.length() * 8;
        }

        public String getValue() {
            return value;
        }
    }

    @Test
    void testValidateOutboundMessageAndMigrate() {
        Message mockMessage = mock(Message.class);
        Element referenceXml = DocumentHelper.createElement("test");
        byte[] data = new byte[]{0x01, 0x02};
        URI uri = URI.create("test://test.xml");

        // Should throw because mock can't be serialized properly
        assertThrows(RuntimeException.class, () -> {
            MessageValidatorAndMigrator.validateOutboundMessageAndMigrate(
                "TestCase",
                mockMessage,
                referenceXml,
                data,
                "BIG_ENDIAN",
                false,
                uri
            );
        });
    }

    @Test
    void testValidateWithNullMessage() {
        Element referenceXml = DocumentHelper.createElement("test");
        byte[] data = new byte[]{};
        URI uri = URI.create("test://test.xml");

        assertThrows(Exception.class, () -> {
            MessageValidatorAndMigrator.validateOutboundMessageAndMigrate(
                "TestCase",
                null,
                referenceXml,
                data,
                "BIG_ENDIAN",
                false,
                uri
            );
        });
    }

    @Test
    void testValidateInboundMessageAndGetWithMessageInput() {
        // Create a simple message input that returns a fixed message
        MessageInput<SimpleTestMessage> messageInput = new MessageInput<SimpleTestMessage>() {
            @Override
            public SimpleTestMessage parse(ReadBuffer buffer) throws BufferException {
                return new SimpleTestMessage("test-value");
            }
        };

        LocationAwareDocumentFactory factory = new LocationAwareDocumentFactory();
        Element referenceXml = factory.createElement(new QName("SimpleTestMessage"));
        referenceXml.addElement("value").setText("test-value");

        Message result = MessageValidatorAndMigrator.validateInboundMessageAndGet(messageInput, referenceXml);

        assertNotNull(result);
        assertInstanceOf(SimpleTestMessage.class, result);
        assertEquals("test-value", ((SimpleTestMessage) result).getValue());
    }

    @Test
    void testValidateInboundMessageAndGetWithOptions() {
        Map<String, String> options = new HashMap<>();
        options.put("package", "invalid.package");

        LocationAwareDocumentFactory factory = new LocationAwareDocumentFactory();
        Element referenceXml = factory.createElement(new QName("TestMessage"));

        // Should throw because the package doesn't exist
        assertThrows(Exception.class, () -> {
            MessageValidatorAndMigrator.validateInboundMessageAndGet(options, referenceXml);
        });
    }

    @Test
    void testValidateInboundMessageWithInvalidXml() {
        MessageInput<SimpleTestMessage> messageInput = new MessageInput<SimpleTestMessage>() {
            @Override
            public SimpleTestMessage parse(ReadBuffer buffer) throws BufferException {
                throw new BufferException("Parse failed");
            }
        };

        Element referenceXml = DocumentHelper.createElement("SimpleTestMessage");

        assertThrows(Exception.class, () -> {
            MessageValidatorAndMigrator.validateInboundMessageAndGet(messageInput, referenceXml);
        });
    }

    @Test
    void testValidateOutboundWithLocationAwareElement() {
        Message mockMessage = mock(Message.class);

        // Create location-aware element directly with location
        LocationAwareElement referenceXml = new LocationAwareElement(
            new QName("test"),
            new org.apache.plc4x.java.utils.testutils.utils.model.Location(10, 5)
        );

        byte[] data = new byte[]{0x01, 0x02};
        URI uri = URI.create("test://test.xml");

        assertThrows(RuntimeException.class, () -> {
            MessageValidatorAndMigrator.validateOutboundMessageAndMigrate(
                "TestCase",
                mockMessage,
                referenceXml,
                data,
                "BIG_ENDIAN",
                false,
                uri
            );
        });
    }
}
