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

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.mock;

class OutgoingPlcBytesHandlerTest {

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
    void testExecuteOutgoingPlcBytesWithMatchingData() throws Exception {
        byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04, 0x05};
        testTransport.write(testData);

        Element payload = DocumentHelper.createElement("test");
        payload.setText("0102030405");

        OutgoingPlcBytesHandler handler = new OutgoingPlcBytesHandler("com.example.test", payload);
        assertDoesNotThrow(() -> handler.executeOutgoingPlcBytes(testTransport, "BIG_ENDIAN"));
    }

    @Test
    void testExecuteOutgoingPlcBytesWithMismatchedData() throws Exception {
        byte[] testData = new byte[]{0x01, 0x02, 0x03};
        testTransport.write(testData);

        Element payload = DocumentHelper.createElement("test");
        payload.setText("0102FF");

        OutgoingPlcBytesHandler handler = new OutgoingPlcBytesHandler("com.example.test", payload);
        assertThrows(DriverTestsuiteException.class, () ->
            handler.executeOutgoingPlcBytes(testTransport, "BIG_ENDIAN"));
    }

    @Test
    void testExecuteOutgoingPlcBytesWithWhitespace() throws Exception {
        byte[] testData = new byte[]{0x01, 0x02, 0x03};
        testTransport.write(testData);

        Element payload = DocumentHelper.createElement("test");
        payload.setText("01 02 03");

        OutgoingPlcBytesHandler handler = new OutgoingPlcBytesHandler("com.example.test", payload);
        assertDoesNotThrow(() -> handler.executeOutgoingPlcBytes(testTransport, "BIG_ENDIAN"));
    }
}
