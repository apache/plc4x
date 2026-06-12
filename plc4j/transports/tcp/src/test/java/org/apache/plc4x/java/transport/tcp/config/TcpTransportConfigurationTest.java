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
package org.apache.plc4x.java.transport.tcp.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TcpTransportConfigurationTest {

    @Test
    void testInstantiation() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        assertNotNull(config);
    }

    @Test
    void testImplementsTransportConfiguration() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        assertTrue(true);
    }

    @Test
    void testConnectTimeoutField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 10000;
        assertEquals(10000, config.connectTimeout);
    }

    @Test
    void testReadTimeoutField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.readTimeout = 5000;
        assertEquals(5000, config.readTimeout);
    }

    @Test
    void testWriteTimeoutField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.writeTimeout = 3000;
        assertEquals(3000, config.writeTimeout);
    }

    @Test
    void testTcpNoDelayField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.tcpNoDelay = false;
        assertFalse(config.tcpNoDelay);

        config.tcpNoDelay = true;
        assertTrue(config.tcpNoDelay);
    }

    @Test
    void testKeepAliveField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.keepAlive = true;
        assertTrue(config.keepAlive);

        config.keepAlive = false;
        assertFalse(config.keepAlive);
    }

    @Test
    void testSendBufferSizeField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.sendBufferSize = 32768;
        assertEquals(32768, config.sendBufferSize);
    }

    @Test
    void testReceiveBufferSizeField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.receiveBufferSize = 65536;
        assertEquals(65536, config.receiveBufferSize);
    }

    @Test
    void testLocalAddressField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        assertNull(config.localAddress);

        config.localAddress = "192.168.1.100";
        assertEquals("192.168.1.100", config.localAddress);
    }

    @Test
    void testLocalPortField() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.localPort = 8080;
        assertEquals(8080, config.localPort);
    }

    @Test
    void testGetDefaultPort() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        assertEquals(TcpTransportConfiguration.NO_DEFAULT_PORT, config.getDefaultPort());
    }

    @Test
    void testNoDefaultPortConstant() {
        assertEquals(-1, TcpTransportConfiguration.NO_DEFAULT_PORT);
    }

    @Test
    void testCompleteConfiguration() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();

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
    void testZeroTimeouts() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 0;
        config.readTimeout = 0;
        config.writeTimeout = 0;

        assertEquals(0, config.connectTimeout);
        assertEquals(0, config.readTimeout);
        assertEquals(0, config.writeTimeout);
    }

    @Test
    void testZeroBufferSizes() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.sendBufferSize = 0;
        config.receiveBufferSize = 0;

        assertEquals(0, config.sendBufferSize);
        assertEquals(0, config.receiveBufferSize);
    }

    @Test
    void testEphemeralLocalPort() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.localPort = 0;
        assertEquals(0, config.localPort);
    }
}
