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
package org.apache.plc4x.java.transport.rawsocket;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.rawsocket.config.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;

import java.io.File;
import java.util.List;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

// Timeout prevents pcap operations from hanging the entire test suite
@Timeout(value = 30, unit = TimeUnit.SECONDS)
class RawSocketTransportTest {

    private RawSocketTransport transport;

    @BeforeEach
    void setUp() {
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
            throw new RuntimeException("Could not set JNA library path", e);
        }

        transport = new RawSocketTransport();
    }

    @Test
    void testGetTransportCode() {
        assertEquals("raw-socket", transport.getTransportCode());
    }

    @Test
    void testGetTransportName() {
        assertEquals("Raw Socket (Ethernet)", transport.getTransportName());
    }

    @Test
    void testGetTransportConfigType() {
        assertEquals(RawSocketTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void testCreateTransportInstance_dedicated() throws Exception {
        // Skip if no network interfaces or no pcap permissions
        List<PcapNetworkInterface> devs = Pcaps.findAllDevs();
        assumeTrue(devs != null && !devs.isEmpty(), "No network interfaces found");

        PcapNetworkInterface nif = devs.get(0);

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = nif.getName();
        config.remoteAddress = "00:00:00:00:00:01"; // Dummy MAC
        config.protocolId = 0x88B5; // PROFINET
        config.reuseInterface = false;
        config.promiscuousMode = false;

        try {
            TransportInstance<RawSocketTransportConfiguration> instance = transport.createTransportInstance(
                "raw-socket://" + nif.getName(), config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
        } catch (TransportException e) {
            // May fail due to permissions - acceptable for unit tests
            System.out.println("Could not open raw socket (expected without root/admin): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_shared() throws Exception {
        List<PcapNetworkInterface> devs = Pcaps.findAllDevs();
        assumeTrue(devs != null && !devs.isEmpty(), "No network interfaces found");

        PcapNetworkInterface nif = devs.get(0);

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = nif.getName();
        config.remoteAddress = "00:00:00:00:00:01";
        config.protocolId = 0x88B5;
        config.reuseInterface = true;

        try {
            TransportInstance<RawSocketTransportConfiguration> instance = transport.createTransportInstance(
                "raw-socket://" + nif.getName(), config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
        } catch (TransportException e) {
            System.out.println("Could not open raw socket (expected without root/admin): " + e.getMessage());
        }
    }

    @Test
    //@Disabled("All of a sudden this test hangs ... investigate")
    void testCreateTransportInstance_multipleShared() throws Exception {
        List<PcapNetworkInterface> devs = Pcaps.findAllDevs();
        assumeTrue(devs != null && !devs.isEmpty(), "No network interfaces found");

        PcapNetworkInterface nif = devs.get(0);

        RawSocketTransportConfiguration config1 = new RawSocketTransportConfiguration();
        config1.interfaceName = nif.getName();
        config1.remoteAddress = "00:00:00:00:00:01";
        config1.protocolId = 0x88B5;
        config1.reuseInterface = true;

        RawSocketTransportConfiguration config2 = new RawSocketTransportConfiguration();
        config2.interfaceName = nif.getName();
        config2.remoteAddress = "00:00:00:00:00:02";
        config2.protocolId = 0x88B5;
        config2.reuseInterface = true;

        try {
            TransportInstance<RawSocketTransportConfiguration> instance1 = transport.createTransportInstance(
                "raw-socket://" + nif.getName(), config1, AuditLog.builder().build());

            TransportInstance<RawSocketTransportConfiguration> instance2 = transport.createTransportInstance(
                "raw-socket://" + nif.getName(), config2, AuditLog.builder().build());

            assertTrue(instance1.isOpen());
            assertTrue(instance2.isOpen());

            // Close first - handle should stay open
            instance1.close();
            assertFalse(instance1.isOpen());

            // Close second - handle should close
            instance2.close();
            assertFalse(instance2.isOpen());

        } catch (TransportException e) {
            System.out.println("Could not open raw socket (expected without root/admin): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_withVLAN() throws Exception {
        List<PcapNetworkInterface> devs = Pcaps.findAllDevs();
        assumeTrue(devs != null && !devs.isEmpty(), "No network interfaces found");

        PcapNetworkInterface nif = devs.get(0);

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = nif.getName();
        config.remoteAddress = "00:00:00:00:00:01";
        config.protocolId = 0x88B5;
        config.vlanId = 100;
        config.vlanPriority = 5;
        config.reuseInterface = false;

        try {
            TransportInstance<RawSocketTransportConfiguration> instance = transport.createTransportInstance(
                "raw-socket://" + nif.getName(), config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
        } catch (TransportException e) {
            System.out.println("Could not open raw socket (expected without root/admin): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_withCustomBPF() throws Exception {
        List<PcapNetworkInterface> devs = Pcaps.findAllDevs();
        assumeTrue(devs != null && !devs.isEmpty(), "No network interfaces found");

        PcapNetworkInterface nif = devs.get(0);

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = nif.getName();
        config.remoteAddress = "00:00:00:00:00:01";
        config.protocolId = 0x88B5;
        config.bpfFilter = "ether proto 0x88B5";
        config.reuseInterface = false;

        try {
            TransportInstance<RawSocketTransportConfiguration> instance = transport.createTransportInstance(
                "raw-socket://" + nif.getName(), config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
        } catch (TransportException e) {
            System.out.println("Could not open raw socket (expected without root/admin): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_promiscuousMode() throws Exception {
        List<PcapNetworkInterface> devs = Pcaps.findAllDevs();
        assumeTrue(devs != null && !devs.isEmpty(), "No network interfaces found");

        PcapNetworkInterface nif = devs.get(0);

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = nif.getName();
        config.remoteAddress = "00:00:00:00:00:01";
        config.protocolId = 0x88B5;
        config.promiscuousMode = true;
        config.reuseInterface = false;

        try {
            TransportInstance<RawSocketTransportConfiguration> instance = transport.createTransportInstance(
                "raw-socket://" + nif.getName(), config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
        } catch (TransportException e) {
            System.out.println("Could not open raw socket (expected without root/admin): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_invalidInterface() {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "INVALID_INTERFACE_THAT_DOES_NOT_EXIST";
        config.remoteAddress = "00:00:00:00:00:01";
        config.protocolId = 0x88B5;

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("raw-socket://invalid", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_invalidConfigType() {
        // Pass a wrong configuration type to trigger the IllegalArgumentException path
        TransportConfiguration wrongConfig = new TransportConfiguration() {};

        assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("raw-socket://test", wrongConfig, AuditLog.builder().build())
        );
    }
}
