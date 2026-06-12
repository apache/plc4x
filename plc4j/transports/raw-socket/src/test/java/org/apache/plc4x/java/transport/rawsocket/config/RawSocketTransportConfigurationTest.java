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
package org.apache.plc4x.java.transport.rawsocket.config;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for RawSocketTransportConfiguration field defaults and type hierarchy.
 */
class RawSocketTransportConfigurationTest {

    @Test
    void testImplementsTransportConfiguration() {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        assertInstanceOf(TransportConfiguration.class, config);
    }

    @Test
    void testDefaultValues() {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();

        assertNull(config.interfaceName);
        assertNull(config.localAddress);
        assertNull(config.remoteAddress);
        assertEquals(0, config.protocolId);
        assertFalse(config.promiscuousMode);
        assertEquals(0, config.captureTimeout);
        assertEquals(0, config.snapshotLength);
        assertEquals(0, config.bufferSize);
        assertEquals(0, config.vlanId);
        assertEquals(0, config.vlanPriority);
        assertFalse(config.reuseInterface);
        assertNull(config.bpfFilter);
        assertEquals(0, config.maxFrameSize);
        assertEquals(0, config.readTimeout);
    }

    @Test
    void testFieldsAreSettable() {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();

        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.promiscuousMode = true;
        config.captureTimeout = 2000;
        config.snapshotLength = 32768;
        config.bufferSize = 2097152;
        config.vlanId = 100;
        config.vlanPriority = 5;
        config.reuseInterface = true;
        config.bpfFilter = "ether proto 0x88B5";
        config.maxFrameSize = 9000;
        config.readTimeout = 5000;

        assertEquals("en0", config.interfaceName);
        assertEquals("00:11:22:33:44:55", config.localAddress);
        assertEquals("AA:BB:CC:DD:EE:FF", config.remoteAddress);
        assertEquals(0x88B5, config.protocolId);
        assertTrue(config.promiscuousMode);
        assertEquals(2000, config.captureTimeout);
        assertEquals(32768, config.snapshotLength);
        assertEquals(2097152, config.bufferSize);
        assertEquals(100, config.vlanId);
        assertEquals(5, config.vlanPriority);
        assertTrue(config.reuseInterface);
        assertEquals("ether proto 0x88B5", config.bpfFilter);
        assertEquals(9000, config.maxFrameSize);
        assertEquals(5000, config.readTimeout);
    }
}
