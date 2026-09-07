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

class IncomingPlcMessageHandlerTest {

    private TestTransportInstance testTransport;
    private AuditLog auditLog;

    @BeforeEach
    void setUp() throws Exception {
        auditLog = mock(AuditLog.class);
        TestTransportConfiguration config = new TestTransportConfiguration();
        config.receiveBufferSize = 8192;
        testTransport = new TestTransportInstance(config, auditLog);
    }

    // Tests for parseDynamic method
    @Test
    void testParseDynamic_Integer() {
        Object result = IncomingPlcMessageHandler.parseDynamic(Integer.class, "42");
        assertEquals(42, result);
    }

    @Test
    void testParseDynamic_int() {
        Object result = IncomingPlcMessageHandler.parseDynamic(int.class, "42");
        assertEquals(42, result);
    }

    @Test
    void testParseDynamic_Long() {
        Object result = IncomingPlcMessageHandler.parseDynamic(Long.class, "1234567890");
        assertEquals(1234567890L, result);
    }

    @Test
    void testParseDynamic_long() {
        Object result = IncomingPlcMessageHandler.parseDynamic(long.class, "9876543210");
        assertEquals(9876543210L, result);
    }

    @Test
    void testParseDynamic_Short() {
        Object result = IncomingPlcMessageHandler.parseDynamic(Short.class, "100");
        assertEquals((short) 100, result);
    }

    @Test
    void testParseDynamic_short() {
        Object result = IncomingPlcMessageHandler.parseDynamic(short.class, "200");
        assertEquals((short) 200, result);
    }

    @Test
    void testParseDynamic_Byte() {
        Object result = IncomingPlcMessageHandler.parseDynamic(Byte.class, "127");
        assertEquals((byte) 127, result);
    }

    @Test
    void testParseDynamic_byte() {
        Object result = IncomingPlcMessageHandler.parseDynamic(byte.class, "64");
        assertEquals((byte) 64, result);
    }

    @Test
    void testParseDynamic_Double() {
        Object result = IncomingPlcMessageHandler.parseDynamic(Double.class, "3.14159");
        assertEquals(3.14159, result);
    }

    @Test
    void testParseDynamic_double() {
        Object result = IncomingPlcMessageHandler.parseDynamic(double.class, "2.71828");
        assertEquals(2.71828, result);
    }

    @Test
    void testParseDynamic_Float() {
        Object result = IncomingPlcMessageHandler.parseDynamic(Float.class, "1.5");
        assertEquals(1.5f, result);
    }

    @Test
    void testParseDynamic_float() {
        Object result = IncomingPlcMessageHandler.parseDynamic(float.class, "2.5");
        assertEquals(2.5f, result);
    }

    @Test
    void testParseDynamic_Boolean_true() {
        Object result = IncomingPlcMessageHandler.parseDynamic(Boolean.class, "true");
        assertEquals(true, result);
    }

    @Test
    void testParseDynamic_Boolean_false() {
        Object result = IncomingPlcMessageHandler.parseDynamic(Boolean.class, "false");
        assertEquals(false, result);
    }

    @Test
    void testParseDynamic_boolean() {
        Object result = IncomingPlcMessageHandler.parseDynamic(boolean.class, "true");
        assertEquals(true, result);
    }

    @Test
    void testParseDynamic_String() {
        Object result = IncomingPlcMessageHandler.parseDynamic(String.class, "hello");
        assertEquals("hello", result);
    }

    @Test
    void testExecuteIncomingPlcMessage_noMessageElement() {
        // Create XML without a message element (only the root)
        Element messageXml = DocumentHelper.createElement("incoming-plc-message");

        IncomingPlcMessageHandler handler = new IncomingPlcMessageHandler("org.apache.plc4x.java.modbus.readwrite", messageXml);
        // Should throw because there's no message element
        RuntimeException exception = assertThrows(RuntimeException.class,
            () -> handler.executeIncomingPlcMessage(testTransport, "BIG_ENDIAN"));
        assertEquals("No message element found in reference XML", exception.getMessage());
    }

    @Test
    void testExecuteIncomingPlcMessage_withParserArgumentsOnly() {
        // Create XML with only parser-arguments (no actual message)
        Element messageXml = DocumentHelper.createElement("incoming-plc-message");
        Element parserArgs = messageXml.addElement("parser-arguments");
        parserArgs.addElement("driverType").setText("MODBUS_TCP");

        IncomingPlcMessageHandler handler = new IncomingPlcMessageHandler("org.apache.plc4x.java.modbus.readwrite", messageXml);
        // Should throw because parser-arguments is filtered out
        RuntimeException exception = assertThrows(RuntimeException.class,
            () -> handler.executeIncomingPlcMessage(testTransport, "BIG_ENDIAN"));
        assertEquals("No message element found in reference XML", exception.getMessage());
    }

    @Test
    void testExecuteIncomingPlcMessage_withUnknownMessageClass() {
        // Create XML with a message element that doesn't have a corresponding class
        Element messageXml = DocumentHelper.createElement("incoming-plc-message");
        Element testMessage = messageXml.addElement("NonExistentMessage");
        testMessage.addText("test");

        IncomingPlcMessageHandler handler = new IncomingPlcMessageHandler("org.apache.plc4x.java.modbus.readwrite", messageXml);
        // Should throw because the class doesn't exist
        RuntimeException exception = assertThrows(RuntimeException.class,
            () -> handler.executeIncomingPlcMessage(testTransport, "BIG_ENDIAN"));
        assertTrue(exception.getMessage().contains("Failed to load class"));
    }
}
