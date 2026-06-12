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

import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.rawsocket.config.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.testutils.RequirePcap;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;
import org.pcap4j.util.MacAddress;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

/**
 * Tests for RawSocketTransportInstance.
 * Note: These tests require the pcap native library (skipped via {@link RequirePcap} when
 * absent, e.g. on Windows CI without Npcap) plus pcap permissions (root/admin); assumptions
 * skip individual tests when permissions are not enough.
 */
@RequirePcap
class RawSocketTransportInstanceTest {

    private PcapHandle handle;
    private PcapNetworkInterface nif;
    private RawSocketTransportInstance transportInstance;
    private RawSocketTransportConfiguration config;
    private MacAddress localMac;
    private MacAddress remoteMac;

    @BeforeEach
    void setUp() {
        try {
            // Find available network interface
            List<PcapNetworkInterface> devs = Pcaps.findAllDevs();
            assumeTrue(devs != null && !devs.isEmpty(), "No network interfaces found");

            // Find an Ethernet-capable interface (not loopback, not tunnel)
            PcapNetworkInterface selectedNif = null;
            for (PcapNetworkInterface dev : devs) {
                String name = dev.getName().toLowerCase();
                // Skip loopback (lo0, lo), tunnel interfaces (utun, gif, stf), and bridge interfaces (bridge)
                if (name.contains("lo") || name.contains("utun") || name.contains("gif") ||
                    name.contains("stf") || name.contains("bridge") || name.contains("awdl")) {
                    continue;
                }

                // Prefer en0 (primary Ethernet/WiFi), otherwise take first valid interface
                if (name.equals("en0")) {
                    selectedNif = dev;
                    break;
                } else if (selectedNif == null && (name.startsWith("en") || name.startsWith("eth"))) {
                    selectedNif = dev;
                }
            }

            assumeTrue(selectedNif != null, "No Ethernet-capable network interface found (tried en*, eth*)");
            nif = selectedNif;

            // Try to open a handle
            handle = nif.openLive(65536, PcapNetworkInterface.PromiscuousMode.NONPROMISCUOUS, 10);

            // Set up MAC addresses
            localMac = MacAddress.getByName("00:11:22:33:44:55");
            remoteMac = MacAddress.getByName("AA:BB:CC:DD:EE:FF");

            // Create configuration
            config = new RawSocketTransportConfiguration();
            config.interfaceName = nif.getName();
            config.localAddress = localMac.toString();
            config.remoteAddress = remoteMac.toString();
            config.protocolId = 0x88B5; // PROFINET
            config.maxFrameSize = 1500;
            config.readTimeout = 100;

            // Create a transport instance
            transportInstance = new RawSocketTransportInstance(new SharedRawSocketManager(), config, AuditLog.builder().build());

            Thread.sleep(100);
        } catch (Exception e) {
            assumeTrue(false, "Could not initialize test (requires pcap permissions): " + e.getMessage());
        }
    }

    @AfterEach
    void tearDown() throws Exception {
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }

        // Close the pcap handle to stop the capture loop
        if (handle != null && handle.isOpen()) {
            handle.breakLoop();  // Stop the packet capture loop
            handle.close();
        }
    }

    @Test
    void testGetConfiguration() {
        RawSocketTransportConfiguration config = transportInstance.getConfiguration();
        assertNotNull(config);
    }

    @Test
    void testIsOpen_whenConnected() {
        assertTrue(transportInstance.isOpen());
    }

    //@Test
    void testIsOpen_whenClosed() throws TransportException {
        transportInstance.close();
        assertFalse(transportInstance.isOpen());
    }

    @Test
    void testGetNumBytesAvailable() throws TransportException {
        // Initially should return 0
        int available = transportInstance.getNumBytesAvailable();
        assertEquals(0, available);
    }

    @Test
    void testWrite_successful() throws TransportException {
        byte[] data = new byte[]{0x01, 0x02, 0x03, 0x04};

        // Should not throw
        assertDoesNotThrow(() -> transportInstance.write(data));
    }

    //@Test
    void testWrite_exceedsMaxFrameSize() {
        byte[] data = new byte[2000]; // Exceeds default 1500

        assertThrows(TransportException.class, () ->
            transportInstance.write(data)
        );
    }

    @Test
    void testWrite_emptyArray() throws TransportException {
        assertDoesNotThrow(() -> transportInstance.write(new byte[0]));
    }

    @Test
    void testWrite_nullArray() throws TransportException {
        assertDoesNotThrow(() -> transportInstance.write(null));
    }

    @Test
    void testRead_zeroBytes() throws TransportException {
        byte[] result = transportInstance.read(0);
        assertEquals(0, result.length);
    }

    @Test
    void testRead_whenClosed_throwsException() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.read(10)
        );
    }

    @Test
    void testWrite_whenClosed_throwsException() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.write(new byte[]{0x01})
        );
    }

//    @Test
    void testClose_idempotent() throws TransportException {
        transportInstance.close();
        assertFalse(transportInstance.isOpen());

        // Should not throw
        assertDoesNotThrow(() -> transportInstance.close());
        assertFalse(transportInstance.isOpen());
    }

    //@Test
    void testPeekReadableBytes_noData() throws TransportException {
        // When no data available, should timeout or throw
        assertThrows(TransportException.class, () ->
            transportInstance.peekReadableBytes(10)
        );
    }

    /**
     * Test with loopback if available.
     * This test is disabled by default as it requires specific network setup.
     */
    @Test
    @Disabled("Requires network loopback setup - enable manually for hardware testing")
    void testLoopback_writeAndRead() throws Exception {
        byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04, 0x05};

        // Write data
        transportInstance.write(testData);

        // Wait for loopback
        Thread.sleep(100);

        // Read data back
        byte[] received = transportInstance.read(testData.length);

        assertArrayEquals(testData, received);
    }

    /**
     * Test peek functionality with loopback.
     * This test is disabled by default as it requires specific network setup.
     */
    @Test
    @Disabled("Requires network loopback setup - enable manually for hardware testing")
    void testLoopback_peekAndRead() throws Exception {
        byte[] testData = new byte[]{(byte) 0xAA, (byte) 0xBB, (byte) 0xCC, (byte) 0xDD};

        // Write data
        transportInstance.write(testData);

        // Wait for loopback
        Thread.sleep(100);

        // Peek data
        byte[] peeked = transportInstance.peekReadableBytes(testData.length);
        assertArrayEquals(testData, peeked);

        // Read data - should still be available
        byte[] read = transportInstance.read(testData.length);
        assertArrayEquals(testData, read);
    }

    // ========== Async Transport Tests ==========

    @Test
    void testRegisterDataListener_listenerRegistration() {
        final boolean[] listenerCalled = {false};

        // Register listener - should not throw
        assertDoesNotThrow(() ->
            transportInstance.registerDataListener(() -> listenerCalled[0] = true)
        );
    }

    @Test
    @Disabled("Requires network loopback setup - enable manually for hardware testing")
    void testOnDataAvailable_listenerTriggeredRegisterPacket() throws Exception {
        final boolean[] listenerCalled = {false};

        // Register listener
        transportInstance.registerDataListener(() -> listenerCalled[0] = true);

        // Send packet (with loopback)
        byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04};
        transportInstance.write(testData);

        // Wait for pcap callback to process
        Thread.sleep(100);

        assertTrue(listenerCalled[0], "Listener should be called when packet arrives");
    }

    @Test
    @Disabled("Requires network loopback setup - enable manually for hardware testing")
    void testRegisterDataListener_immediateNotification() throws Exception {
        final long[] notificationTime = {0};
        final long[] packetSentTime = {0};

        // Register listener that records timestamp
        transportInstance.registerDataListener(() -> notificationTime[0] = System.nanoTime());

        // Send packet and record time
        byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04};
        packetSentTime[0] = System.nanoTime();
        transportInstance.write(testData);

        // Wait for notification
        Thread.sleep(100);

        assertTrue(notificationTime[0] > 0, "Listener should have been called");

        // Calculate latency in milliseconds
        long latencyMs = (notificationTime[0] - packetSentTime[0]) / 1_000_000;

        // Should be notified within reasonable time (< 50ms for sub-millisecond timing critical protocols)
        assertTrue(latencyMs < 50,
            "Notification should be immediate for timing-critical protocols (was " + latencyMs + "ms)");
    }

    @Test
    void testRemoveDataListener() {
        // Register listener
        transportInstance.registerDataListener(() -> {});

        // Remove listener - should not throw
        assertDoesNotThrow(() -> transportInstance.removeDataListener());
    }

    @Test
    @Disabled("Requires network loopback setup - enable manually for hardware testing")
    void testRemoveDataListener_stopsNotifications() throws Exception {
        final int[] callCount = {0};

        // Register listener
        transportInstance.registerDataListener(() -> callCount[0]++);

        // Send packet
        transportInstance.write(new byte[]{0x01, 0x02, 0x03});
        Thread.sleep(50);

        int firstCount = callCount[0];
        assertTrue(firstCount > 0, "Listener should be called initially");

        // Remove listener
        transportInstance.removeDataListener();

        // Send another packet
        transportInstance.write(new byte[]{0x04, 0x05, 0x06});
        Thread.sleep(50);

        // Count should not increase
        assertEquals(firstCount, callCount[0], "Listener should not be called after removal");
    }

    @Test
    void testPcapThread_startsAutomatically() {
        // pcap thread should be running after construction
        assertTrue(transportInstance.isOpen());

        // Transport should be ready to receive packets
        assertEquals(0, transportInstance.getNumBytesAvailable());
    }

    @Test
    void testPcapThread_continuousCapture() throws Exception {
        // pcap thread should continuously capture packets
        assertTrue(transportInstance.isOpen());

        // Wait a bit to ensure thread is running
        Thread.sleep(100);

        // Should still be open and capturing
        assertTrue(transportInstance.isOpen());
    }

    /**
     * Test that async notification works with PROFINET-like timing requirements.
     * This verifies sub-millisecond notification latency critical for industrial protocols.
     */
    @Test
    @Disabled("Requires network loopback setup - enable manually for hardware testing")
    void testTimingCritical_subMillisecondNotification() throws Exception {
        final long[] notificationTimes = new long[10];
        final int[] notificationCount = {0};

        // Register listener that tracks timing
        transportInstance.registerDataListener(() -> {
            int idx = notificationCount[0]++;
            if (idx < notificationTimes.length) {
                notificationTimes[idx] = System.nanoTime();
            }
        });

        // Send multiple packets rapidly (simulating PROFINET cycle)
        long startTime = System.nanoTime();
        for (int i = 0; i < 10; i++) {
            transportInstance.write(new byte[]{(byte) i, 0x01, 0x02, 0x03});
            Thread.sleep(1); // 1ms cycle time (typical for PROFINET)
        }

        // Wait for all notifications
        Thread.sleep(100);

        // Verify all packets triggered notifications
        assertEquals(10, notificationCount[0], "All packets should trigger notifications");

        // Verify timing: each notification should be close to when packet was sent
        for (int i = 0; i < 10; i++) {
            long latency = (notificationTimes[i] - startTime) / 1_000_000;
            assertTrue(latency < 100,
                "Notification " + i + " should arrive quickly (was " + latency + "ms)");
        }
    }
}
