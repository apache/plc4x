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
package org.apache.plc4x.java.transport.test;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.test.config.TestTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

/**
 * Tests for TestTransport.
 */
class TestTransportTest {

    private TestTransport transport;
    private AuditLog mockAuditLog;

    @BeforeEach
    void setUp() {
        transport = new TestTransport();
        mockAuditLog = mock(AuditLog.class);
    }

    @Test
    void testGetTransportCode() {
        assertEquals("test", transport.getTransportCode());
    }

    @Test
    void testGetTransportName() {
        assertEquals("Test Transport", transport.getTransportName());
    }

    @Test
    void testGetTransportConfigType() {
        assertEquals(TestTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void testCreateTransportInstanceWithValidConfiguration() throws TransportException {
        TestTransportConfiguration config = new TestTransportConfiguration();
        config.receiveBufferSize = 81920;

        TransportInstance<TestTransportConfiguration> instance =
            transport.createTransportInstance("test://localhost", config, mockAuditLog);

        assertNotNull(instance);
        assertInstanceOf(TestTransportInstance.class, instance);
        assertTrue(instance.isOpen());
    }

    @Test
    void testCreateTransportInstanceWithWrongConfigurationType() {
        // Create a mock configuration that is not TestTransportConfiguration
        TransportConfiguration wrongConfig = mock(TransportConfiguration.class);

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("test://localhost", wrongConfig, mockAuditLog)
        );

        assertTrue(exception.getMessage().contains("Expected configuration of type"));
        assertTrue(exception.getMessage().contains("TestTransportConfiguration"));
    }

    @Test
    void testCreateMultipleInstances() throws TransportException {
        TestTransportConfiguration config = new TestTransportConfiguration();
        config.receiveBufferSize = 81920;

        TransportInstance<TestTransportConfiguration> instance1 =
            transport.createTransportInstance("test://endpoint1", config, mockAuditLog);
        TransportInstance<TestTransportConfiguration> instance2 =
            transport.createTransportInstance("test://endpoint2", config, mockAuditLog);

        assertNotNull(instance1);
        assertNotNull(instance2);
        assertNotSame(instance1, instance2);
        assertTrue(instance1.isOpen());
        assertTrue(instance2.isOpen());
    }

}
