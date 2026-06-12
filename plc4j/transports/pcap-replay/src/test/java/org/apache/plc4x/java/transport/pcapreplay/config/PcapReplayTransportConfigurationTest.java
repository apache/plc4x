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
package org.apache.plc4x.java.transport.pcapreplay.config;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.pcapreplay.PcapFilePlayer;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.junit.jupiter.api.Assertions.*;

class PcapReplayTransportConfigurationTest {

    @Test
    void testInstantiation() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertNotNull(config);
    }

    @Test
    void testImplementsTransportConfiguration() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertTrue(config instanceof TransportConfiguration);
    }

    @Test
    void testPcapFileField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertNull(config.pcapFile);

        config.pcapFile = "/path/to/test.pcap";
        assertEquals("/path/to/test.pcap", config.pcapFile);
    }

    @Test
    void testLocalAddressField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertNull(config.localAddress);

        config.localAddress = "00:11:22:33:44:55";
        assertEquals("00:11:22:33:44:55", config.localAddress);
    }

    @Test
    void testRemoteAddressField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertNull(config.remoteAddress);

        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        assertEquals("AA:BB:CC:DD:EE:FF", config.remoteAddress);
    }

    @Test
    void testProtocolIdField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertEquals(0, config.protocolId);

        config.protocolId = 0x88B5; // PROFINET
        assertEquals(0x88B5, config.protocolId);
    }

    @Test
    void testSpeedFactorField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertEquals(0.0, config.speedFactor, 0.001);

        config.speedFactor = 1.0; // Real-time
        assertEquals(1.0, config.speedFactor, 0.001);

        config.speedFactor = 2.0; // Double speed
        assertEquals(2.0, config.speedFactor, 0.001);

        config.speedFactor = 0.5; // Half speed
        assertEquals(0.5, config.speedFactor, 0.001);
    }

    @Test
    void testLoopField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertFalse(config.loop);

        config.loop = true;
        assertTrue(config.loop);

        config.loop = false;
        assertFalse(config.loop);
    }

    @Test
    void testOnlyIncomingPacketsField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertFalse(config.onlyIncomingPackets);

        config.onlyIncomingPackets = true;
        assertTrue(config.onlyIncomingPackets);

        config.onlyIncomingPackets = false;
        assertFalse(config.onlyIncomingPackets);
    }

    @Test
    void testOnlyOutgoingPacketsField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertFalse(config.onlyOutgoingPackets);

        config.onlyOutgoingPackets = true;
        assertTrue(config.onlyOutgoingPackets);

        config.onlyOutgoingPackets = false;
        assertFalse(config.onlyOutgoingPackets);
    }

    @Test
    void testMaxFrameSizeField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertEquals(0, config.maxFrameSize);

        config.maxFrameSize = 1500; // Standard Ethernet MTU
        assertEquals(1500, config.maxFrameSize);

        config.maxFrameSize = 9000; // Jumbo frame
        assertEquals(9000, config.maxFrameSize);
    }

    @Test
    void testReadTimeoutField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertEquals(0, config.readTimeout);

        config.readTimeout = 5000; // 5 seconds
        assertEquals(5000, config.readTimeout);

        config.readTimeout = 0; // No timeout
        assertEquals(0, config.readTimeout);
    }

    @Test
    void testPacketQueueSizeField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertEquals(0, config.packetQueueSize);

        config.packetQueueSize = 1000;
        assertEquals(1000, config.packetQueueSize);

        config.packetQueueSize = 100;
        assertEquals(100, config.packetQueueSize);
    }

    @Test
    void testAutoDetectMacAddressesField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertFalse(config.autoDetectMacAddresses);

        config.autoDetectMacAddresses = true;
        assertTrue(config.autoDetectMacAddresses);

        config.autoDetectMacAddresses = false;
        assertFalse(config.autoDetectMacAddresses);
    }

    @Test
    void testAutoStartField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertFalse(config.autoStart);

        config.autoStart = true;
        assertTrue(config.autoStart);

        config.autoStart = false;
        assertFalse(config.autoStart);
    }

    @Test
    void testVlanIdField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertEquals(0, config.vlanId);

        config.vlanId = 100;
        assertEquals(100, config.vlanId);

        config.vlanId = 4095; // Max VLAN ID
        assertEquals(4095, config.vlanId);
    }

    @Test
    void testMockPlayerField() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        assertNull(config.mockPlayer);

        PcapFilePlayer mockPlayer = Mockito.mock(PcapFilePlayer.class);
        config.mockPlayer = mockPlayer;
        assertEquals(mockPlayer, config.mockPlayer);
    }

    @Test
    void testCompleteConfiguration() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();

        // Set all fields
        config.pcapFile = "/path/to/capture.pcap";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.speedFactor = 1.5;
        config.loop = true;
        config.onlyIncomingPackets = false;
        config.onlyOutgoingPackets = true;
        config.maxFrameSize = 1500;
        config.readTimeout = 3000;
        config.packetQueueSize = 500;
        config.autoDetectMacAddresses = false;
        config.autoStart = false;
        config.vlanId = 200;
        PcapFilePlayer mockPlayer = Mockito.mock(PcapFilePlayer.class);
        config.mockPlayer = mockPlayer;

        // Verify all fields
        assertEquals("/path/to/capture.pcap", config.pcapFile);
        assertEquals("00:11:22:33:44:55", config.localAddress);
        assertEquals("AA:BB:CC:DD:EE:FF", config.remoteAddress);
        assertEquals(0x88B5, config.protocolId);
        assertEquals(1.5, config.speedFactor, 0.001);
        assertTrue(config.loop);
        assertFalse(config.onlyIncomingPackets);
        assertTrue(config.onlyOutgoingPackets);
        assertEquals(1500, config.maxFrameSize);
        assertEquals(3000, config.readTimeout);
        assertEquals(500, config.packetQueueSize);
        assertFalse(config.autoDetectMacAddresses);
        assertFalse(config.autoStart);
        assertEquals(200, config.vlanId);
        assertEquals(mockPlayer, config.mockPlayer);
    }

    @Test
    void testMutuallyExclusivePacketDirection() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();

        // Test that both can be set (validation would be done elsewhere)
        config.onlyIncomingPackets = true;
        config.onlyOutgoingPackets = true;

        assertTrue(config.onlyIncomingPackets);
        assertTrue(config.onlyOutgoingPackets);
    }

    @Test
    void testNegativeValues() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();

        // Test that negative values can be set (validation would be done elsewhere)
        config.protocolId = -1;
        config.maxFrameSize = -100;
        config.readTimeout = -1000;
        config.packetQueueSize = -50;
        config.vlanId = -10;

        assertEquals(-1, config.protocolId);
        assertEquals(-100, config.maxFrameSize);
        assertEquals(-1000, config.readTimeout);
        assertEquals(-50, config.packetQueueSize);
        assertEquals(-10, config.vlanId);
    }

    @Test
    void testEdgeCaseSpeedFactors() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();

        // Test edge cases
        config.speedFactor = 0.0; // As fast as possible
        assertEquals(0.0, config.speedFactor, 0.001);

        config.speedFactor = Double.MAX_VALUE;
        assertEquals(Double.MAX_VALUE, config.speedFactor, 0.001);

        config.speedFactor = Double.MIN_VALUE;
        assertEquals(Double.MIN_VALUE, config.speedFactor, 0.001);

        config.speedFactor = -1.0; // Negative (validation would handle this)
        assertEquals(-1.0, config.speedFactor, 0.001);
    }
}
