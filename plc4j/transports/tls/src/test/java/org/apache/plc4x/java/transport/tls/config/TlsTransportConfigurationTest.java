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
package org.apache.plc4x.java.transport.tls.config;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TlsTransportConfigurationTest {

    @Test
    void testInstantiation() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        assertNotNull(config);
    }

    @Test
    void testImplementsTransportConfiguration() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        assertInstanceOf(TransportConfiguration.class, config);
    }

    @Test
    void testExtendsTcpTransportConfiguration() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        assertInstanceOf(TcpTransportConfiguration.class, config);
    }

    @Test
    void testVerifySslFieldDefaultValue() {
        // The default value from @BooleanDefaultValue(true) annotation
        // Note: The annotation provides metadata, actual default is false until populated
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        // Default for uninitialized boolean is false
        assertFalse(config.verifySsl);
    }

    @Test
    void testVerifySslFieldCanBeSetTrue() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        config.verifySsl = true;
        assertTrue(config.verifySsl);
        assertTrue(config.isVerifySsl());
    }

    @Test
    void testVerifySslFieldCanBeSetFalse() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        config.verifySsl = false;
        assertFalse(config.verifySsl);
        assertFalse(config.isVerifySsl());
    }

    @Test
    void testIsVerifySslMethod() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();

        config.verifySsl = true;
        assertTrue(config.isVerifySsl());

        config.verifySsl = false;
        assertFalse(config.isVerifySsl());
    }

    @Test
    void testInheritedTcpConfiguration() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();

        // Test inherited TCP configuration fields
        config.connectTimeout = 10000;
        config.readTimeout = 5000;
        config.writeTimeout = 3000;
        config.tcpNoDelay = true;
        config.keepAlive = true;
        config.sendBufferSize = 32768;
        config.receiveBufferSize = 65536;
        config.localAddress = "192.168.1.100";
        config.localPort = 8080;

        assertEquals(10000, config.connectTimeout);
        assertEquals(5000, config.readTimeout);
        assertEquals(3000, config.writeTimeout);
        assertTrue(config.tcpNoDelay);
        assertTrue(config.keepAlive);
        assertEquals(32768, config.sendBufferSize);
        assertEquals(65536, config.receiveBufferSize);
        assertEquals("192.168.1.100", config.localAddress);
        assertEquals(8080, config.localPort);
    }

    @Test
    void testInheritedGetDefaultPort() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        // TLS configuration inherits TCP's default port (-1 meaning no default)
        assertEquals(TcpTransportConfiguration.NO_DEFAULT_PORT, config.getDefaultPort());
    }

    @Test
    void testCompleteConfiguration() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();

        // Set TLS-specific configuration
        config.verifySsl = false;

        // Set inherited TCP configuration
        config.connectTimeout = 10000;
        config.readTimeout = 5000;
        config.writeTimeout = 3000;
        config.tcpNoDelay = true;
        config.keepAlive = true;
        config.sendBufferSize = 32768;
        config.receiveBufferSize = 65536;

        // Verify all settings
        assertFalse(config.verifySsl);
        assertEquals(10000, config.connectTimeout);
        assertEquals(5000, config.readTimeout);
        assertEquals(3000, config.writeTimeout);
        assertTrue(config.tcpNoDelay);
        assertTrue(config.keepAlive);
        assertEquals(32768, config.sendBufferSize);
        assertEquals(65536, config.receiveBufferSize);
    }

    @Test
    void testVerifySslWithDevelopmentMode() {
        // Simulates development mode where self-signed certs are used
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        config.verifySsl = false;

        assertFalse(config.isVerifySsl());
    }

    @Test
    void testVerifySslWithProductionMode() {
        // Simulates production mode where proper certificates are used
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        config.verifySsl = true;

        assertTrue(config.isVerifySsl());
    }

    @Test
    void testTlsVersionFieldDefaultValue() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        assertNull(config.tlsVersion);
        assertNull(config.getTlsVersion());
    }

    @Test
    void testTlsVersionFieldCanBeSet() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();

        config.tlsVersion = "TLSv1.2";
        assertEquals("TLSv1.2", config.tlsVersion);
        assertEquals("TLSv1.2", config.getTlsVersion());

        config.tlsVersion = "TLSv1.3";
        assertEquals("TLSv1.3", config.getTlsVersion());
    }

}
