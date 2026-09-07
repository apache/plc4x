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
package org.apache.plc4x.java.transport.tls;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.tls.config.PskTlsTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class PskTlsTransportTest {

    private PskTlsTransport transport;

    @BeforeEach
    void setUp() {
        transport = new PskTlsTransport();
    }

    @Test
    void testGetTransportCode() {
        assertEquals("tls-psk", transport.getTransportCode());
    }

    @Test
    void testGetTransportName() {
        assertEquals("TLS-PSK", transport.getTransportName());
    }

    @Test
    void testGetTransportConfigType() {
        assertEquals(PskTlsTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void testCreateTransportInstance_withWrongConfigurationType() {
        TransportConfiguration wrongConfig = new TransportConfiguration() {};

        assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("127.0.0.1:8016", wrongConfig, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withInvalidUrl() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskIdentity = "TestUser";
        config.pskKey = "0123456789abcdef";

        assertThrows(Exception.class, () ->
            transport.createTransportInstance(":", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withoutPort() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskIdentity = "TestUser";
        config.pskKey = "0123456789abcdef";

        assertThrows(Exception.class, () ->
            transport.createTransportInstance("127.0.0.1", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_connectionRefused() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;
        config.pskIdentity = "TestUser";
        config.pskKey = "0123456789abcdef";

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("127.0.0.1:1", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_validationFailsEmptyIdentity() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;
        config.pskIdentity = "";
        config.pskKey = "0123456789abcdef";

        assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("127.0.0.1:9999", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_validationFailsEmptyKey() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;
        config.pskIdentity = "TestUser";
        config.pskKey = "";

        assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("127.0.0.1:9999", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_validationFailsInvalidHex() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;
        config.pskIdentity = "TestUser";
        config.pskKey = "ZZZ";

        assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("127.0.0.1:9999", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withDefaultPort() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration() {
            @Override
            public int getDefaultPort() {
                return 8016;
            }
        };
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;
        config.pskIdentity = "TestUser";
        config.pskKey = "0123456789abcdef";

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("127.0.0.1", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withIpAddress() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;
        config.pskIdentity = "TestUser";
        config.pskKey = "0123456789abcdef";

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("127.0.0.1:9999", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withExplicitPort() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;
        config.pskIdentity = "TestUser";
        config.pskKey = "0123456789abcdef";

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("127.0.0.1:9999", config, AuditLog.builder().build())
        );
    }
}
