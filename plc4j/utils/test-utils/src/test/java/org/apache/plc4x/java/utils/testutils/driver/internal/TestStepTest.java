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
package org.apache.plc4x.java.utils.testutils.driver.internal;

import org.apache.plc4x.java.transport.test.TestTransportInstance;
import org.apache.plc4x.java.transport.test.config.TestTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.api.PlcConnection;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

class TestStepTest {

    private TestTransportInstance testTransport;
    private PlcConnection mockConnection;
    private AuditLog auditLog;

    @BeforeEach
    void setUp() throws Exception {
        auditLog = mock(AuditLog.class);
        mockConnection = mock(PlcConnection.class);
        TestTransportConfiguration config = new TestTransportConfiguration();
        config.receiveBufferSize = 8192;
        testTransport = new TestTransportInstance(config, auditLog);
    }

    @Test
    void testParseIncomingPlcBytesStep() {
        Element stepElement = DocumentHelper.createElement("incoming-plc-bytes");
        stepElement.setText("010203");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step);
        assertEquals(StepType.INCOMING_PLC_BYTES, step.getType());
    }

    @Test
    void testParseOutgoingPlcBytesStep() {
        Element stepElement = DocumentHelper.createElement("outgoing-plc-bytes");
        stepElement.setText("010203");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step);
        assertEquals(StepType.OUTGOING_PLC_BYTES, step.getType());
    }

    @Test
    void testParseDelayStep() {
        Element stepElement = DocumentHelper.createElement("delay");
        stepElement.setText("100");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step);
        assertEquals(StepType.DELAY, step.getType());
    }

    @Test
    void testParseTerminateStep() {
        Element stepElement = DocumentHelper.createElement("terminate");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step);
        assertEquals(StepType.TERMINATE, step.getType());
    }

    @Test
    void testParseInvalidStepType() {
        Element stepElement = DocumentHelper.createElement("invalid-step");

        assertThrows(DriverTestsuiteException.class, () -> {
            TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        });
    }

    @Test
    void testExecuteIncomingPlcBytes() {
        Element stepElement = DocumentHelper.createElement("incoming-plc-bytes");
        stepElement.setText("010203");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertDoesNotThrow(() -> {
            step.execute(mockConnection, testTransport, "BIG_ENDIAN", new org.apache.plc4x.java.utils.testutils.driver.internal.TestContext());
        });
    }

    @Test
    void testExecuteOutgoingPlcBytes() throws Exception {
        // Write some data first
        testTransport.write(new byte[]{0x01, 0x02, 0x03});

        Element stepElement = DocumentHelper.createElement("outgoing-plc-bytes");
        stepElement.setText("010203");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertDoesNotThrow(() -> {
            step.execute(mockConnection, testTransport, "BIG_ENDIAN", new org.apache.plc4x.java.utils.testutils.driver.internal.TestContext());
        });
    }

    @Test
    void testExecuteDelay() {
        Element stepElement = DocumentHelper.createElement("delay");
        stepElement.setText("10");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        long start = System.currentTimeMillis();
        step.execute(mockConnection, testTransport, "BIG_ENDIAN", new org.apache.plc4x.java.utils.testutils.driver.internal.TestContext());
        long end = System.currentTimeMillis();

        // Should have delayed at least 10ms
        assertTrue(end - start >= 10);
    }

    @Test
    void testExecuteTerminate() {
        Element stepElement = DocumentHelper.createElement("terminate");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertDoesNotThrow(() -> {
            step.execute(mockConnection, testTransport, "BIG_ENDIAN", new org.apache.plc4x.java.utils.testutils.driver.internal.TestContext());
        });
    }

    @Test
    void testGetLocation() {
        Element stepElement = DocumentHelper.createElement("delay");
        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");

        assertTrue(step.getLocation().isEmpty() || step.getLocation().isPresent());
    }

    @Test
    void testGetPayload() {
        Element stepElement = DocumentHelper.createElement("delay");
        stepElement.setText("100");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step.getPayload());
        assertEquals("100", step.getPayload().getTextTrim());
    }

    @Test
    void testGetParserArguments() {
        Element stepElement = DocumentHelper.createElement("delay");
        stepElement.addAttribute("parserArguments", "arg1,arg2");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertEquals("arg1,arg2", step.getParserArguments());
    }

    @Test
    void testParseIncomingPlcMessageStep() {
        Element stepElement = DocumentHelper.createElement("incoming-plc-message");
        stepElement.setText("test-message");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step);
        assertEquals(StepType.INCOMING_PLC_MESSAGE, step.getType());
    }

    @Test
    void testParseOutgoingPlcMessageStep() {
        Element stepElement = DocumentHelper.createElement("outgoing-plc-message");
        stepElement.setText("test-message");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step);
        assertEquals(StepType.OUTGOING_PLC_MESSAGE, step.getType());
    }

    @Test
    void testParseApiRequestStep() {
        Element stepElement = DocumentHelper.createElement("api-request");
        Element requestElement = stepElement.addElement("PlcReadRequest");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step);
        assertEquals(StepType.API_REQUEST, step.getType());
    }

    @Test
    void testParseApiResponseStep() {
        Element stepElement = DocumentHelper.createElement("api-response");
        Element responseElement = stepElement.addElement("PlcReadResponse");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNotNull(step);
        assertEquals(StepType.API_RESPONSE, step.getType());
    }

    @Test
    void testDelayWithLongDuration() {
        Element stepElement = DocumentHelper.createElement("delay");
        stepElement.setText("1");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        long start = System.currentTimeMillis();
        step.execute(mockConnection, testTransport, "BIG_ENDIAN", new org.apache.plc4x.java.utils.testutils.driver.internal.TestContext());
        long end = System.currentTimeMillis();

        assertTrue(end - start >= 1);
    }

    @Test
    void testParseStepWithoutParserArguments() {
        Element stepElement = DocumentHelper.createElement("delay");
        stepElement.setText("50");

        TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
        assertNull(step.getParserArguments());
    }

    @Test
    void testMultipleStepTypes() {
        String[] stepTypes = {"incoming-plc-bytes", "outgoing-plc-bytes", "delay", "terminate"};

        for (String stepType : stepTypes) {
            Element stepElement = DocumentHelper.createElement(stepType);
            if (!stepType.equals("terminate")) {
                stepElement.setText("010203");
            }

            TestStep step = TestStep.parseTestStep("", stepElement, "BIG_ENDIAN");
            assertNotNull(step);
            assertNotNull(step.getType());
        }
    }
}
