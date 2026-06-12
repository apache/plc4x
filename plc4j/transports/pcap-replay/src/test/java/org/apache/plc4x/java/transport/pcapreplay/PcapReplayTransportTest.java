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
package org.apache.plc4x.java.transport.pcapreplay;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.pcapreplay.config.PcapReplayTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.pcap4j.packet.factory.PacketFactories;
import org.pcap4j.packet.namednumber.DataLinkType;

import java.io.File;
import java.io.IOException;

import static org.junit.jupiter.api.Assertions.*;

class PcapReplayTransportTest {

    private PcapReplayTransport transport;
    private File tempPcapFile;

    @BeforeAll
    static void initPacketFactory() {
        // Force pcap4j to eagerly load the PacketFactoryBinder.
        // Without this, the factory loads lazily on the first player thread,
        // causing a race where early tests finish before parsing is available.
        PacketFactories.getFactory(org.pcap4j.packet.Packet.class, DataLinkType.class);
    }

    @BeforeEach
    void setUp() throws IOException {
        try {
            // For some reason it doesn't work if we pass this in from the outside.
            //if (os == "mac") {
            // On my Intel Mac I found the libs in: "/usr/local/Cellar/libpcap/1.10.1/lib"
            // On my M1 Mac I found the libs in: "/opt/homebrew/Cellar/libpcap/1.10.1/lib"
            if (new File("/usr/local/Cellar/libpcap/1.10.1/lib").exists()) {
                System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.1/lib");
            } else if (new File("/usr/local/Cellar/libpcap/1.10.5/lib").exists()) {
                System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.5/lib");
            } else if (new File("/opt/homebrew/opt/libpcap/lib").exists()) {
                System.getProperties().setProperty("jna.library.path", "/opt/homebrew/opt/libpcap/lib");
            }
            //}
        } catch (Error e) {
            e.printStackTrace();
        }

        transport = new PcapReplayTransport();

        // Create a minimal test PCAP file
        tempPcapFile = File.createTempFile("test", ".pcap");
        // Note: For real tests, you'd need a valid PCAP file
        // This is just for basic structure testing
    }

    @AfterEach
    void tearDown() {
        if (tempPcapFile != null && tempPcapFile.exists()) {
            tempPcapFile.delete();
        }
    }

    @Test
    void testGetTransportCode() {
        assertEquals("pcap-replay", transport.getTransportCode());
    }

    @Test
    void testGetTransportName() {
        assertEquals("PCAP Replay", transport.getTransportName());
    }

    @Test
    void testGetTransportConfigType() {
        assertEquals(PcapReplayTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void testCreateTransportInstance_fileNotFound() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        config.pcapFile = "/nonexistent/path/to/file.pcap";
        config.autoDetectMacAddresses = false;
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("pcap-replay://test", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_missingMacAddresses() {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        config.pcapFile = tempPcapFile.getAbsolutePath();
        config.autoDetectMacAddresses = false;
        config.maxFrameSize = 1500;
        // Missing MAC addresses

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("pcap-replay://test", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withManualMacAddresses() {
        // This test would require a valid PCAP file
        // Skipping actual creation due to need for valid PCAP content

        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        config.pcapFile = tempPcapFile.getAbsolutePath();
        config.autoDetectMacAddresses = false;
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.autoStart = false;
        config.maxFrameSize = 1500;

        // Would fail due to invalid PCAP format, but tests configuration
        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("pcap-replay://test", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withValidPcap() throws Exception {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        config.pcapFile = "/test.pcap";
        config.autoDetectMacAddresses = true;
        config.protocolId = 0x88B5;
        config.speedFactor = 0; // As fast as possible
        config.loop = false;
        config.autoStart = true;
        config.maxFrameSize = 1500;
        config.packetQueueSize = 1000;

        TransportInstance<PcapReplayTransportConfiguration> instance =
            transport.createTransportInstance("pcap-replay://test", config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstance_withLoop() throws Exception {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        config.pcapFile = "/test.pcap";
        config.autoDetectMacAddresses = true;
        config.protocolId = 0x88B5;
        config.speedFactor = 2.0; // Double speed
        config.loop = true;
        config.autoStart = true;
        config.maxFrameSize = 1500;
        config.packetQueueSize = 1000;

        TransportInstance<PcapReplayTransportConfiguration> instance =
            transport.createTransportInstance("pcap-replay://test", config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstance_onlyIncoming() throws Exception {
        PcapReplayTransportConfiguration config = new PcapReplayTransportConfiguration();
        config.pcapFile = "/test.pcap";
        config.autoDetectMacAddresses = true;
        config.protocolId = 0x88B5;
        config.onlyIncomingPackets = true;
        config.onlyOutgoingPackets = false;
        config.maxFrameSize = 1500;
        config.packetQueueSize = 1000;

        TransportInstance<PcapReplayTransportConfiguration> instance =
            transport.createTransportInstance("pcap-replay://test", config, AuditLog.builder().build());

        assertNotNull(instance);
        instance.close();
    }
}
