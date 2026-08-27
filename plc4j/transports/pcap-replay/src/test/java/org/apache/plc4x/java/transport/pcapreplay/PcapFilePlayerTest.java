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

import org.apache.plc4x.java.utils.testutils.RequirePcap;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.pcap4j.packet.factory.PacketFactories;
import org.pcap4j.packet.namednumber.DataLinkType;
import org.pcap4j.util.MacAddress;

import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for PcapFilePlayer.
 * <p>
 * Every test here drives real pcap playback, so the whole class is skipped when the pcap
 * native library can't be loaded (e.g. on Windows CI without Npcap installed).
 */
@RequirePcap
class PcapFilePlayerTest {

    private PcapFilePlayer player;
    private MacAddress localMac;
    private MacAddress remoteMac;

    @BeforeAll
    static void initPacketFactory() {
        // Force pcap4j to eagerly load the PacketFactoryBinder.
        // Without this, the factory loads lazily on the first player thread,
        // causing a race where early tests finish before parsing is available.
        PacketFactories.getFactory(org.pcap4j.packet.Packet.class, DataLinkType.class);
    }

    @BeforeEach
    void setUp() {
        localMac = MacAddress.getByName("00:11:22:33:44:55");
        remoteMac = MacAddress.getByName("AA:BB:CC:DD:EE:FF");
    }

    @AfterEach
    void tearDown() {
        if (player != null) {
            player.stop();
        }
    }

    @Test
    void testInitialState() {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            1.0,
            false,
            true,
            false,
            100
        );

        assertFalse(player.isPlaying());
        assertEquals(0, player.getPacketsReplayed());
    }

    @Test
    void testStart_changesPlayingState() throws InterruptedException {
        // Use valid test pcap with looping to ensure player stays running
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,  // No protocol filter
            0.0,
            true,  // Enable looping so player doesn't stop immediately
            false,
            false,
            100
        );

        player.start();
        // Give the thread time to start and stabilize
        Thread.sleep(50);
        assertTrue(player.isPlaying());
    }

    @Test
    void testStop_changesPlayingState() throws InterruptedException {
        // Use valid test pcap with looping to ensure player stays running
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,  // No protocol filter
            0.0,
            true,  // Enable looping so player doesn't stop immediately
            false,
            false,
            100
        );

        player.start();
        // Give the thread time to start and stabilize
        Thread.sleep(50);
        assertTrue(player.isPlaying());

        player.stop();
        assertFalse(player.isPlaying());
    }

    @Test
    void testGetNextPacket_timeout() throws InterruptedException {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            0.0,
            false,
            true,
            false,
            100
        );

        // Should timeout since no valid file
        byte[] packet = player.getNextPacket(100, TimeUnit.MILLISECONDS);
        assertNull(packet); // Timeout
    }

    @Test
    void testGetQueueSize_initially() {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            0.0,
            false,
            true,
            false,
            100
        );

        assertEquals(0, player.getQueueSize());
    }

    @Test
    void testPlayback_withValidPcap() throws Exception {
        // test.pcap contains 4 packets (2 incoming, 2 outgoing) with EtherType 0x88B5
        // and the same MAC addresses used in this test class
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            0.0, // As fast as possible
            false,
            true, // Only incoming
            false,
            100
        );

        player.start();
        Thread.sleep(200); // Let it play for a bit

        assertTrue(player.getPacketsReplayed() > 0);

        byte[] packet = player.getNextPacket(1, TimeUnit.SECONDS);
        assertNotNull(packet);
        assertTrue(packet.length > 0);
    }

    @Test
    void testPlayback_withLoop() throws Exception {
        // test.pcap has 2 incoming 0x88B5 packets; with looping, count grows over time.
        // Use a large queue so it doesn't fill up and stall the counter.
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            0.0,
            true, // Loop
            true,
            false,
            100000
        );

        player.start();
        Thread.sleep(200); // Let it loop

        long packets1 = player.getPacketsReplayed();
        Thread.sleep(200);
        long packets2 = player.getPacketsReplayed();

        // Should have replayed more packets (looping)
        assertTrue(packets2 > packets1,
                "Expected more packets after looping: packets1=" + packets1 + ", packets2=" + packets2);
    }

    @Test
    void testPlayback_speedFactor() throws Exception {
        // test.pcap has 2 incoming 0x88B5 packets 20ms apart (at indices 0 and 2).
        // At real-time speed (1.0), replaying should take at least ~20ms.
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            1.0, // Real-time
            false,
            true,
            false,
            100
        );

        long start = System.currentTimeMillis();
        player.start();

        // Wait for the 2 incoming packets to be replayed (with timeout to avoid hang)
        long deadline = start + 5000;
        while (player.getPacketsReplayed() < 2 && System.currentTimeMillis() < deadline) {
            Thread.sleep(10);
        }

        long elapsed = System.currentTimeMillis() - start;

        // At real-time speed, the 20ms inter-packet gap means replay takes nonzero time
        assertTrue(elapsed > 10, "Real-time replay should take measurable time");
        assertTrue(player.getPacketsReplayed() >= 2, "Should have replayed both incoming packets");
    }

    @Test
    void testGetPacketsReplayed_initially() {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            0.0,
            false,
            true,
            false,
            100
        );

        assertEquals(0, player.getPacketsReplayed());
    }

    @Test
    void testConstructor_withAllParameters() {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            2.0,  // Double speed
            true,  // Loop
            false, // Not only incoming
            true,  // Only outgoing
            500    // Queue size
        );

        assertNotNull(player);
        assertFalse(player.isPlaying());
        assertEquals(0, player.getPacketsReplayed());
        assertEquals(0, player.getQueueSize());
    }

    @Test
    void testStart_alreadyPlaying() throws InterruptedException {
        // Use valid test pcap with looping to ensure player stays running
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,  // No protocol filter
            0.0,
            true,  // Enable looping so player doesn't stop immediately
            false,
            false,
            100
        );

        player.start();
        // Give the thread time to start and stabilize
        Thread.sleep(50);
        assertTrue(player.isPlaying());

        // Start again - should be idempotent
        player.start();
        assertTrue(player.isPlaying());
    }

    @Test
    void testStop_notPlaying() {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            0.0,
            false,
            true,
            false,
            100
        );

        // Stop without starting - should not throw
        assertDoesNotThrow(() -> player.stop());
        assertFalse(player.isPlaying());
    }

    @Test
    void testStop_idempotent() throws InterruptedException {
        // Use valid test pcap with looping to ensure player stays running
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,  // No protocol filter to ensure packets match
            0.0,
            true,  // Enable looping so player doesn't stop immediately
            false,
            false,
            100
        );

        player.start();
        // Give the thread time to start and stabilize
        Thread.sleep(50);
        assertTrue(player.isPlaying());

        player.stop();
        assertFalse(player.isPlaying());

        // Stop again - should be idempotent
        player.stop();
        assertFalse(player.isPlaying());
    }

    @Test
    void testConstructor_withSmallQueueSize() {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0,  // No protocol filter
            0.0,
            false,
            false, // Not only incoming
            false, // Not only outgoing
            1    // Minimal queue size
        );

        assertNotNull(player);
        assertEquals(0, player.getQueueSize());
    }

    @Test
    void testGetNextPacket_withoutStart() throws InterruptedException {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            0.0,
            false,
            true,
            false,
            100
        );

        // Should timeout immediately since no packets are being played
        byte[] packet = player.getNextPacket(50, TimeUnit.MILLISECONDS);
        assertNull(packet);
    }

    @Test
    void testConstructor_withBothDirectionFlags() {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            1.5,
            false,
            true,  // Only incoming
            true,  // Only outgoing (contradictory)
            100
        );

        assertNotNull(player);
        // The contradictory flags are handled in the filter logic
    }

    @Test
    void testGetQueueSize_afterStarting() throws InterruptedException {
        player = new PcapFilePlayer(
            "/path/to/test.pcap",
            localMac,
            remoteMac,
            0x88B5,
            0.0,
            false,
            true,
            false,
            100
        );

        player.start();
        Thread.sleep(50); // Give it time to process

        // Queue size should be >= 0
        int size = player.getQueueSize();
        assertTrue(size >= 0);
    }

    @Test
    void testPlayback_withTestPcap() throws Exception {
        // Use the test.pcap from classpath (will be extracted automatically by PcapFilePlayer)
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,  // No protocol filter
            0.0, // As fast as possible
            false,
            false, // Accept all packets
            false,
            1000
        );

        player.start();

        // No isPlaying() check here: test.pcap is 328 bytes and this player replays it once, as
        // fast as it can, so playback can be over before the next statement runs - the flag says
        // more about thread scheduling than about the player. (Every other test in this class
        // that asserts isPlaying() passes loop=true so the player cannot finish on its own; this
        // one is deliberately a single pass.) What the replay did is the thing worth asserting.
        awaitPlaybackEnd(player);

        // Every packet in the file, rather than ">= 0", which a count can never fail.
        assertEquals(4, player.getPacketsReplayed(), "all four packets in test.pcap");

        player.stop();
        assertFalse(player.isPlaying());
    }

    /**
     * Waits for a single-pass replay to finish, rather than sleeping for a fixed time and hoping.
     * A fixed sleep is either too short on a loaded machine - which is when this test used to
     * fail - or slower than it needs to be on an idle one.
     */
    private static void awaitPlaybackEnd(PcapFilePlayer player) throws InterruptedException {
        long deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(5);
        while (player.isPlaying() && (System.nanoTime() < deadline)) {
            Thread.sleep(5);
        }
        assertFalse(player.isPlaying(), "a single pass over a 328 byte file should be long done");
    }

    @Test
    void testPlayback_withProtocolFilter() throws Exception {
        // Test with protocol filter
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0x0800,  // IPv4 protocol
            0.0,
            false,
            false,
            false,
            1000
        );

        player.start();
        Thread.sleep(100);

        long packetsReplayed = player.getPacketsReplayed();
        assertTrue(packetsReplayed >= 0);

        player.stop();
    }

    @Test
    void testPlayback_onlyIncomingPackets() throws Exception {
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,
            0.0,
            false,
            true,  // Only incoming
            false,
            1000
        );

        player.start();
        Thread.sleep(100);

        player.stop();
        assertFalse(player.isPlaying());
    }

    @Test
    void testPlayback_onlyOutgoingPackets() throws Exception {
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,
            0.0,
            false,
            false,
            true,  // Only outgoing
            1000
        );

        player.start();
        Thread.sleep(100);

        player.stop();
        assertFalse(player.isPlaying());
    }

    @Test
    void testPlayback_withBothDirectionFilters() throws Exception {
        // Both flags set - should reject all packets
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,
            0.0,
            false,
            true,  // Only incoming
            true,  // Only outgoing (contradictory)
            1000
        );

        player.start();
        Thread.sleep(100);

        // Should not have replayed any packets due to contradictory filters
        long packetsReplayed = player.getPacketsReplayed();
        assertEquals(0, packetsReplayed);

        player.stop();
    }

    @Test
    void testPlayback_withLooping() throws Exception {
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,
            0.0,
            true,  // Loop enabled
            false,
            false,
            1000
        );

        player.start();
        Thread.sleep(300); // Give time to loop

        long packetsAfterDelay = player.getPacketsReplayed();

        player.stop();
        assertFalse(player.isPlaying());
    }

    @Test
    void testPlayback_withSpeedFactor() throws Exception {
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,
            0.5,  // Half speed
            false,
            false,
            false,
            1000
        );

        long startTime = System.currentTimeMillis();
        player.start();
        Thread.sleep(200);

        player.stop();
    }

    @Test
    void testGetNextPacket_withActualData() throws Exception {
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,
            0.0,
            false,
            false,
            false,
            1000
        );

        player.start();
        Thread.sleep(100); // Wait for some packets

        // Try to get a packet with timeout
        byte[] packet = player.getNextPacket(500, TimeUnit.MILLISECONDS);
        // May be null if no matching packets, or contain data

        player.stop();
    }

    @Test
    void testStop_duringPlayback() throws Exception {
        player = new PcapFilePlayer(
            "/test.pcap",
            localMac,
            remoteMac,
            0,
            0.0,
            true,  // Loop - will run indefinitely
            false,
            false,
            1000
        );

        player.start();
        assertTrue(player.isPlaying());

        Thread.sleep(50);

        player.stop();
        assertFalse(player.isPlaying());

        // Wait for thread to fully stop
        Thread.sleep(50);
        assertFalse(player.isPlaying());
    }

    @Test
    void testPlayback_invalidPcapFile() throws Exception {
        player = new PcapFilePlayer(
            "/nonexistent/invalid.pcap",
            localMac,
            remoteMac,
            0,
            0.0,
            false,
            false,
            false,
            1000
        );

        player.start();
        Thread.sleep(100);

        // Should handle gracefully - player might stop due to error
        // but shouldn't crash
        player.stop();
    }
}
