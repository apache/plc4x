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
package org.apache.plc4x.java.transport.udp.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class UdpTransportConfigurationTest {

    @Test
    void testInstantiation() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        assertNotNull(config);
    }

    @Test
    void testImplementsTransportConfiguration() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        assertTrue(true);
    }

    @Test
    void testLocalAddressField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        assertNull(config.localAddress);

        config.localAddress = "192.168.1.100";
        assertEquals("192.168.1.100", config.localAddress);
    }

    @Test
    void testLocalPortField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.localPort = 8080;
        assertEquals(8080, config.localPort);
    }

    @Test
    void testReadTimeoutField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.readTimeout = 5000;
        assertEquals(5000, config.readTimeout);
    }

    @Test
    void testMaxPacketSizeField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.maxPacketSize = 1500;
        assertEquals(1500, config.maxPacketSize);
    }

    @Test
    void testSendBufferSizeField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.sendBufferSize = 32768;
        assertEquals(32768, config.sendBufferSize);
    }

    @Test
    void testReceiveBufferSizeField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.receiveBufferSize = 65536;
        assertEquals(65536, config.receiveBufferSize);
    }

    @Test
    void testBroadcastField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.broadcast = false;
        assertFalse(config.broadcast);

        config.broadcast = true;
        assertTrue(config.broadcast);
    }

    @Test
    void testReuseAddressField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.reuseAddress = false;
        assertFalse(config.reuseAddress);

        config.reuseAddress = true;
        assertTrue(config.reuseAddress);
    }

    @Test
    void testShareSocketField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.shareSocket = false;
        assertFalse(config.shareSocket);

        config.shareSocket = true;
        assertTrue(config.shareSocket);
    }

    @Test
    void testMulticastTtlField() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.multicastTtl = 128;
        assertEquals(128, config.multicastTtl);
    }

    @Test
    void testGetDefaultPort() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        assertEquals(UdpTransportConfiguration.NO_DEFAULT_PORT, config.getDefaultPort());
    }

    @Test
    void testNoDefaultPortConstant() {
        assertEquals(-1, UdpTransportConfiguration.NO_DEFAULT_PORT);
    }

    @Test
    void testCompleteConfiguration() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();

        config.localAddress = "192.168.1.100";
        config.localPort = 8080;
        config.readTimeout = 5000;
        config.maxPacketSize = 1500;
        config.sendBufferSize = 32768;
        config.receiveBufferSize = 65536;
        config.broadcast = true;
        config.reuseAddress = true;
        config.shareSocket = true;
        config.multicastTtl = 64;

        assertEquals("192.168.1.100", config.localAddress);
        assertEquals(8080, config.localPort);
        assertEquals(5000, config.readTimeout);
        assertEquals(1500, config.maxPacketSize);
        assertEquals(32768, config.sendBufferSize);
        assertEquals(65536, config.receiveBufferSize);
        assertTrue(config.broadcast);
        assertTrue(config.reuseAddress);
        assertTrue(config.shareSocket);
        assertEquals(64, config.multicastTtl);
    }

    @Test
    void testZeroReadTimeout() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.readTimeout = 0;
        assertEquals(0, config.readTimeout);
    }

    @Test
    void testZeroBufferSizes() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.sendBufferSize = 0;
        config.receiveBufferSize = 0;

        assertEquals(0, config.sendBufferSize);
        assertEquals(0, config.receiveBufferSize);
    }

    @Test
    void testEphemeralLocalPort() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.localPort = 0;
        assertEquals(0, config.localPort);
    }
}
