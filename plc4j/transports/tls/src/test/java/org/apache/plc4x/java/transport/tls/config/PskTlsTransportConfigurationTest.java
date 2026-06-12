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

class PskTlsTransportConfigurationTest {

    @Test
    void testInstantiation() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        assertNotNull(config);
    }

    @Test
    void testImplementsTransportConfiguration() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        assertInstanceOf(TransportConfiguration.class, config);
    }

    @Test
    void testExtendsTcpTransportConfiguration() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        assertInstanceOf(TcpTransportConfiguration.class, config);
    }

    @Test
    void testPskFieldsDefaultToNull() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        assertNull(config.pskIdentity);
        assertNull(config.pskKey);
    }

    @Test
    void testLogSessionKeysDefaultToFalse() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        assertFalse(config.logSessionKeys);
    }

    @Test
    void testInheritedGetDefaultPort() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        assertEquals(TcpTransportConfiguration.NO_DEFAULT_PORT, config.getDefaultPort());
    }

    @Test
    void testInheritedTcpConfiguration() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();

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

    // ========== PSK Validation Tests ==========

    @Test
    void testValidatePskConfiguration_passesWhenBothSetWithValidHex() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskIdentity = "MyPskUser";
        config.pskKey = "0123456789abcdef";
        assertDoesNotThrow(config::validatePskConfiguration);
    }

    @Test
    void testValidatePskConfiguration_throwsForNullIdentity() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskKey = "0123456789abcdef";
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, config::validatePskConfiguration);
        assertTrue(ex.getMessage().contains("psk-identity must not be empty"));
    }

    @Test
    void testValidatePskConfiguration_throwsForEmptyIdentity() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskIdentity = "";
        config.pskKey = "0123456789abcdef";
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, config::validatePskConfiguration);
        assertTrue(ex.getMessage().contains("psk-identity must not be empty"));
    }

    @Test
    void testValidatePskConfiguration_throwsForNullKey() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskIdentity = "MyPskUser";
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, config::validatePskConfiguration);
        assertTrue(ex.getMessage().contains("psk-key must not be empty"));
    }

    @Test
    void testValidatePskConfiguration_throwsForEmptyKey() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskIdentity = "MyPskUser";
        config.pskKey = "";
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, config::validatePskConfiguration);
        assertTrue(ex.getMessage().contains("psk-key must not be empty"));
    }

    @Test
    void testValidatePskConfiguration_throwsForOddLengthHex() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskIdentity = "MyPskUser";
        config.pskKey = "012";
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, config::validatePskConfiguration);
        assertTrue(ex.getMessage().contains("even length"));
    }

    @Test
    void testValidatePskConfiguration_throwsForInvalidHexCharacters() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskIdentity = "MyPskUser";
        config.pskKey = "ZZZZ";
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, config::validatePskConfiguration);
        assertTrue(ex.getMessage().contains("hexadecimal"));
    }

    // ========== Key Bytes Conversion Tests ==========

    @Test
    void testGetPskKeyBytes() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskKey = "0123456789abcdef";
        byte[] expected = new byte[]{0x01, 0x23, 0x45, 0x67, (byte) 0x89, (byte) 0xab, (byte) 0xcd, (byte) 0xef};
        assertArrayEquals(expected, config.getPskKeyBytes());
    }

    @Test
    void testGetPskKeyBytes_upperCaseHex() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.pskKey = "AABBCCDD";
        byte[] expected = new byte[]{(byte) 0xAA, (byte) 0xBB, (byte) 0xCC, (byte) 0xDD};
        assertArrayEquals(expected, config.getPskKeyBytes());
    }

    @Test
    void testGetPskKeyBytes_throwsWhenKeyNull() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        assertThrows(IllegalStateException.class, config::getPskKeyBytes);
    }
}
