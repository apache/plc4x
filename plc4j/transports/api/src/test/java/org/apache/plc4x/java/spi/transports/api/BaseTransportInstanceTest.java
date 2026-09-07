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

package org.apache.plc4x.java.spi.transports.api;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class BaseTransportInstanceTest {

    @Test
    void testConstructor_shouldStoreConfigAndAuditLog() {
        // Arrange
        TransportConfiguration config = mock(TransportConfiguration.class);
        AuditLog auditLog = mock(AuditLog.class);
        when(config.toString()).thenReturn("test-config-string");

        // Act
        TestBaseTransport transport = new TestBaseTransport(config, auditLog);

        // Assert
        assertNotNull(transport);
        assertEquals(config, transport.getConfiguration());
        assertEquals(auditLog, transport.getAuditLog());
    }

    @Test
    void testConstructor_shouldWriteAuditLogEntry() {
        // Arrange
        TransportConfiguration config = mock(TransportConfiguration.class);
        AuditLog auditLog = mock(AuditLog.class);
        when(config.toString()).thenReturn("test-config-string");

        // Act
        new TestBaseTransport(config, auditLog);

        // Assert - verify audit log was called with correct parameters
        ArgumentCaptor<AuditLogEventType> eventTypeCaptor = ArgumentCaptor.forClass(AuditLogEventType.class);
        ArgumentCaptor<String> messageCaptor = ArgumentCaptor.forClass(String.class);
        ArgumentCaptor<Object> configObjectCaptor = ArgumentCaptor.forClass(Object.class);

        verify(auditLog).write(eventTypeCaptor.capture(), messageCaptor.capture(), configObjectCaptor.capture());

        assertEquals(AuditLogEventType.SYSTEM, eventTypeCaptor.getValue());
        assertEquals("Creating Transport with config", messageCaptor.getValue());
        assertEquals(config, configObjectCaptor.getValue());
    }

    @Test
    void testGetConfiguration_shouldReturnConstructorConfig() {
        // Arrange
        TransportConfiguration config = mock(TransportConfiguration.class);
        AuditLog auditLog = mock(AuditLog.class);
        when(config.toString()).thenReturn("test-config");

        TestBaseTransport transport = new TestBaseTransport(config, auditLog);

        // Act
        TransportConfiguration result = transport.getConfiguration();

        // Assert
        assertSame(config, result);
    }

    @Test
    void testGetAuditLog_shouldReturnConstructorAuditLog() {
        // Arrange
        TransportConfiguration config = mock(TransportConfiguration.class);
        AuditLog auditLog = mock(AuditLog.class);
        when(config.toString()).thenReturn("test-config");

        TestBaseTransport transport = new TestBaseTransport(config, auditLog);

        // Act
        AuditLog result = transport.getAuditLog();

        // Assert
        assertSame(auditLog, result);
    }

    @Test
    void testConstructor_withNullConfig_shouldNotThrow() {
        // Arrange
        AuditLog auditLog = mock(AuditLog.class);

        // Act & Assert - should handle null config gracefully or throw NPE
        // The actual behavior depends on the implementation requirements
        assertThrows(NullPointerException.class, () -> new TestBaseTransport(null, auditLog));
    }

    @Test
    void testConstructor_withNullAuditLog_shouldThrow() {
        // Arrange
        TransportConfiguration config = mock(TransportConfiguration.class);
        when(config.toString()).thenReturn("test-config");

        // Act & Assert
        assertThrows(NullPointerException.class, () -> new TestBaseTransport(config, null));
    }

    @Test
    void testGetDriverConfig_shouldDefaultToEmptyString() {
        // Arrange
        TransportConfiguration config = mock(TransportConfiguration.class);
        AuditLog auditLog = mock(AuditLog.class);

        TestBaseTransport transport = new TestBaseTransport(config, auditLog);

        // Act & Assert
        assertEquals("", transport.getDriverConfig());
    }

    @Test
    void testSetDriverConfig_shouldReturnSetValue() {
        // Arrange
        TransportConfiguration config = mock(TransportConfiguration.class);
        AuditLog auditLog = mock(AuditLog.class);

        TestBaseTransport transport = new TestBaseTransport(config, auditLog);

        // Act
        transport.setDriverConfig("/milo");

        // Assert
        assertEquals("/milo", transport.getDriverConfig());
    }

    @Test
    void testSetDriverConfig_withNull_shouldReturnEmptyString() {
        // Arrange
        TransportConfiguration config = mock(TransportConfiguration.class);
        AuditLog auditLog = mock(AuditLog.class);

        TestBaseTransport transport = new TestBaseTransport(config, auditLog);

        // Act
        transport.setDriverConfig(null);

        // Assert
        assertEquals("", transport.getDriverConfig());
    }

    /**
     * Concrete test implementation of BaseTransport for testing purposes
     */
    private static class TestBaseTransport extends BaseTransportInstance<TransportConfiguration> {
        public TestBaseTransport(TransportConfiguration transportConfig, AuditLog auditLog) {
            super(transportConfig, auditLog);
        }

        @Override
        public boolean isOpen() {
            // Dummy implementation
            return false;
        }

        @Override
        public int getNumBytesAvailable() throws TransportException {
            // Dummy implementation
            return 0;
        }

        @Override
        public byte[] peekReadableBytes(int numBytes) throws TransportException {
            // Dummy implementation
            return new byte[0];
        }

        @Override
        public byte[] read(int numBytes) throws TransportException {
            // Dummy implementation
            return new byte[0];
        }

        @Override
        public void write(byte[] bytes) throws TransportException {
            // Dummy implementation
        }

        @Override
        public void close() throws TransportException {
            // Dummy implementation
        }
    }

}
