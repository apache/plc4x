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
package org.apache.plc4x.java.utils.testutils.driver.internal.handlers;

import org.apache.plc4x.java.transport.test.TestTransportInstance;
import org.apache.plc4x.java.transport.test.config.TestTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

class OutgoingPlcMessageHandlerTest {

    private TestTransportInstance testTransport;
    private AuditLog auditLog;

    @BeforeEach
    void setUp() throws Exception {
        auditLog = mock(AuditLog.class);
        TestTransportConfiguration config = new TestTransportConfiguration();
        config.receiveBufferSize = 8192;
        testTransport = new TestTransportInstance(config, auditLog);
    }

    @Test
    void testExecuteOutgoingPlcMessage_noMessageElement() {
        // Create XML without a message element (only the root)
        Element referenceXml = DocumentHelper.createElement("outgoing-plc-message");

        OutgoingPlcMessageHandler handler = new OutgoingPlcMessageHandler("", referenceXml);
        // Should throw because there's no message element
        RuntimeException exception = assertThrows(RuntimeException.class,
            () -> handler.executeOutgoingPlcMessage(testTransport, "BIG_ENDIAN"));
        assertEquals("No message element found in reference XML", exception.getMessage());
    }

    @Test
    void testExecuteOutgoingPlcMessage_withParserArgumentsOnly() {
        // Create XML with only parser-arguments (no actual message)
        Element referenceXml = DocumentHelper.createElement("outgoing-plc-message");
        Element parserArgs = referenceXml.addElement("parser-arguments");
        parserArgs.addElement("driverType").setText("MODBUS_TCP");

        OutgoingPlcMessageHandler handler = new OutgoingPlcMessageHandler("", referenceXml);
        // Should throw because parser-arguments is filtered out
        RuntimeException exception = assertThrows(RuntimeException.class,
            () -> handler.executeOutgoingPlcMessage(testTransport, "BIG_ENDIAN"));
        assertEquals("No message element found in reference XML", exception.getMessage());
    }

    @Test
    void testExecuteOutgoingPlcMessage_withUnknownMessageClass() {
        // Create XML with a message element that doesn't have a corresponding class
        Element referenceXml = DocumentHelper.createElement("outgoing-plc-message");
        Element testMessage = referenceXml.addElement("NonExistentMessage");
        testMessage.addText("test");

        OutgoingPlcMessageHandler handler = new OutgoingPlcMessageHandler("", referenceXml);
        // Should throw because the class doesn't exist
        RuntimeException exception = assertThrows(RuntimeException.class,
            () -> handler.executeOutgoingPlcMessage(testTransport, "BIG_ENDIAN"));
        assertTrue(exception.getMessage().contains("Failed to load class"));
    }

    @Test
    void testParseDynamic_withString() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(String.class, "test");
        assertEquals("test", result);
    }

    @Test
    void testParseDynamic_withInteger() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(Integer.class, "42");
        assertEquals(42, result);
    }

    @Test
    void testParseDynamic_withIntPrimitive() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(int.class, "42");
        assertEquals(42, result);
    }

    @Test
    void testParseDynamic_withLong() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(Long.class, "12345678901");
        assertEquals(12345678901L, result);
    }

    @Test
    void testParseDynamic_withBoolean() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(Boolean.class, "true");
        assertEquals(true, result);
    }

    @Test
    void testParseDynamic_withBooleanPrimitive() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(boolean.class, "false");
        assertEquals(false, result);
    }

    @Test
    void testParseDynamic_withDouble() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(Double.class, "3.14159");
        assertEquals(3.14159, result);
    }

    @Test
    void testParseDynamic_withFloat() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(Float.class, "3.14");
        assertEquals(3.14f, result);
    }

    @Test
    void testParseDynamic_withShort() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(Short.class, "123");
        assertEquals((short) 123, result);
    }

    @Test
    void testParseDynamic_withByte() {
        Object result = OutgoingPlcMessageHandler.parseDynamic(Byte.class, "42");
        assertEquals((byte) 42, result);
    }

    @Test
    void testParseDynamic_withUnsupportedType() {
        // Should throw for types that can't be parsed
        assertThrows(RuntimeException.class,
            () -> OutgoingPlcMessageHandler.parseDynamic(OutgoingPlcMessageHandler.class, "test"));
    }
}
