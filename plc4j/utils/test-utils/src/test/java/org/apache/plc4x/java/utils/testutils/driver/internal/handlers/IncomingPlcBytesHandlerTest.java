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
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

class IncomingPlcBytesHandlerTest {

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
    void testExecuteIncomingPlcBytes() throws Exception {
        Element payload = DocumentHelper.createElement("test");
        payload.setText("0102030405");

        IncomingPlcBytesHandler handler = new IncomingPlcBytesHandler(payload);
        assertDoesNotThrow(() -> handler.executeIncomingPlcBytes(testTransport, "BIG_ENDIAN"));

        // Verify bytes were injected
        assertTrue(testTransport.getNumBytesAvailable() > 0);
    }

    @Test
    void testExecuteIncomingPlcBytesWithWhitespace() throws Exception {
        Element payload = DocumentHelper.createElement("test");
        payload.setText("01 02 03 04 05");

        IncomingPlcBytesHandler handler = new IncomingPlcBytesHandler(payload);
        assertDoesNotThrow(() -> handler.executeIncomingPlcBytes(testTransport, "BIG_ENDIAN"));

        assertTrue(testTransport.getNumBytesAvailable() > 0);
    }

    @Test
    void testExecuteIncomingPlcBytesWithInvalidHex() {
        Element payload = DocumentHelper.createElement("test");
        payload.setText("ZZZZ");

        IncomingPlcBytesHandler handler = new IncomingPlcBytesHandler(payload);
        assertThrows(DriverTestsuiteException.class, () ->
            handler.executeIncomingPlcBytes(testTransport, "BIG_ENDIAN"));
    }
}
