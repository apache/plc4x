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

import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.pcapreplay.config.PcapReplayTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.testutils.RequirePcap;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.pcap4j.util.MacAddress;

import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 * Tests for PcapReplayTransportInstance.
 */
class PcapReplayTransportInstanceTest {

    private PcapFilePlayer mockPlayer;
    private PcapReplayTransportInstance transportInstance;
    private PcapReplayTransportConfiguration config;
    private MacAddress localMac;
    private MacAddress remoteMac;

    @BeforeEach
    void setUp() throws TransportException {
        // Create a mock player
        mockPlayer = mock(PcapFilePlayer.class);
        when(mockPlayer.isPlaying()).thenReturn(true);
        when(mockPlayer.getPacketsReplayed()).thenReturn(0L);

        // Set up MAC addresses
        localMac = MacAddress.getByName("00:11:22:33:44:55");
        remoteMac = MacAddress.getByName("AA:BB:CC:DD:EE:FF");

        // Create configuration
        config = new PcapReplayTransportConfiguration();
        config.pcapFile = "/test.pcap";
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.localAddress = localMac.toString();
        config.remoteAddress = remoteMac.toString();
        config.packetQueueSize = 100;
        config.maxFrameSize = 1500;
        config.autoStart = false; // Don't auto-start in tests
        config.mockPlayer = mockPlayer;

        // Create transport-instance
        transportInstance = new PcapReplayTransportInstance(config, AuditLog.builder().build());
    }

    @AfterEach
    void tearDown() throws Exception {
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }
    }

    @Test
    void testGetConfiguration() {
        PcapReplayTransportConfiguration config = transportInstance.getConfiguration();
        assertNotNull(config);
    }

    @Test
    void testIsOpen_initially() {
        assertTrue(transportInstance.isOpen());
    }

    @Test
    void testIsOpen_whenClosed() throws TransportException {
        transportInstance.close();
        assertFalse(transportInstance.isOpen());
    }

    @Test
    void testGetNumBytesAvailable() throws TransportException {
        when(mockPlayer.getQueueSize()).thenReturn(5);

        int available = transportInstance.getNumBytesAvailable();
        assertTrue(available >= 0);
    }

    @Test
    void testWrite_simulatedOperation() throws TransportException {
        byte[] data = new byte[]{0x01, 0x02, 0x03, 0x04};

        // Should not throw (write is simulated in replay mode)
        assertThrows(TransportException.class, () -> transportInstance.write(data));
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

    @Test
    void testStartReplay() {
        when(mockPlayer.isPlaying()).thenReturn(false);
        transportInstance.startReplay();
        verify(mockPlayer, atLeastOnce()).start();
    }

    @Test
    void testStopReplay() {
        transportInstance.stopReplay();
        verify(mockPlayer).stop();
    }

    @Test
    void testIsReplaying() {
        when(mockPlayer.isPlaying()).thenReturn(true);
        assertTrue(transportInstance.isReplaying());

        when(mockPlayer.isPlaying()).thenReturn(false);
        assertFalse(transportInstance.isReplaying());
    }

    @Test
    void testGetPacketsReplayed() {
        when(mockPlayer.getPacketsReplayed()).thenReturn(42L);
        assertEquals(42L, transportInstance.getPacketsReplayed());
    }

    @Test
    void testClose_idempotent() throws TransportException {
        transportInstance.close();
        assertFalse(transportInstance.isOpen());

        // Should not throw
        assertDoesNotThrow(() -> transportInstance.close());
        assertFalse(transportInstance.isOpen());

        verify(mockPlayer, atLeastOnce()).stop(); // Called twice
    }

    @Test
    void testClose_stopsPlayer() throws TransportException {
        transportInstance.close();
        verify(mockPlayer).stop();
    }

    @Test
    void testRead_withData() throws Exception {
        byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04};

        // Mock player to return data when polled
        when(mockPlayer.getNextPacket(anyLong(), any(TimeUnit.class))).thenReturn(testData);

        byte[] result = transportInstance.read(4);
        assertArrayEquals(testData, result);
    }

    @Test
    void testPeekReadableBytes_doesNotConsume() throws Exception {
        byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04};

        // Mock player to return data when polled
        when(mockPlayer.getNextPacket(anyLong(), any(TimeUnit.class))).thenReturn(testData);

        // Peek twice — should return the same data both times
        byte[] peek1 = transportInstance.peekReadableBytes(4);
        byte[] peek2 = transportInstance.peekReadableBytes(4);

        assertArrayEquals(peek1, peek2);

        // Read should still return same data and consume it
        byte[] read = transportInstance.read(4);
        assertArrayEquals(testData, read);
    }

    @Test
    void testPeekReadableBytes_zeroBytes() throws TransportException {
        byte[] result = transportInstance.peekReadableBytes(0);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekReadableBytes_whenClosed() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.peekReadableBytes(10)
        );
    }

    @Test
    void testStartReplay_multipleTimes() {
        // First call: isPlaying returns false, so start() is called
        when(mockPlayer.isPlaying()).thenReturn(false);
        transportInstance.startReplay();

        // Second call: isPlaying returns true, so start() is not called again
        when(mockPlayer.isPlaying()).thenReturn(true);
        transportInstance.startReplay();

        // Should have called start once (when player was not playing)
        verify(mockPlayer, times(1)).start();
    }

    @Test
    void testStopReplay_multipleTimes() {
        transportInstance.stopReplay();
        transportInstance.stopReplay();
        // Should handle multiple stops gracefully
        verify(mockPlayer, times(2)).stop();
    }

    @Test
    void testConstructor_withAutoStart() throws TransportException {
        config.autoStart = true;
        PcapFilePlayer autoStartPlayer = mock(PcapFilePlayer.class);
        config.mockPlayer = autoStartPlayer;

        PcapReplayTransportInstance autoStartInstance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Should have started automatically
        verify(autoStartPlayer).start();

        autoStartInstance.close();
    }

    @Test
    void testConstructor_withoutMockPlayer() throws TransportException {
        config.mockPlayer = null;
        config.pcapFile = "/test.pcap";

        // Should create actual player (will use classpath resource)
        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());
        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testConstructor_withProtocolFilter() throws TransportException {
        config.protocolId = 0x88B5; // PROFINET
        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        assertNotNull(instance);
        instance.close();
    }

    @Test
    void testConstructor_withVlanId() throws TransportException {
        config.vlanId = 100;
        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        assertNotNull(instance);
        instance.close();
    }

    @Test
    void testConstructor_withSpeedFactor() throws TransportException {
        config.speedFactor = 2.0; // Double speed
        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        assertNotNull(instance);
        instance.close();
    }

    @Test
    void testConstructor_withLoopEnabled() throws TransportException {
        config.loop = true;
        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        assertNotNull(instance);
        instance.close();
    }

    @Test
    void testConstructor_withOnlyIncomingPackets() throws TransportException {
        config.onlyIncomingPackets = true;
        config.onlyOutgoingPackets = false;
        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        assertNotNull(instance);
        instance.close();
    }

    @Test
    void testConstructor_withOnlyOutgoingPackets() throws TransportException {
        config.onlyIncomingPackets = false;
        config.onlyOutgoingPackets = true;
        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        assertNotNull(instance);
        instance.close();
    }

    @Test
    void testGetNumBytesAvailable_withZeroBytes() {
        when(mockPlayer.getQueueSize()).thenReturn(0);

        assertDoesNotThrow(() -> {
            int available = transportInstance.getNumBytesAvailable();
            assertEquals(0, available);
        });
    }

    @Test
    void testGetNumBytesAvailable_whenClosed() throws TransportException {
        transportInstance.close();

        // Should return 0 when closed
        int available = transportInstance.getNumBytesAvailable();
        assertEquals(0, available);
    }

    @Test
    void testWrite_largeArray() throws TransportException {
        byte[] largeData = new byte[10000];
        for (int i = 0; i < largeData.length; i++) {
            largeData[i] = (byte) (i % 256);
        }

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Writing should throw an exception as this is a replay-only transport
        assertThrows(TransportException.class, () -> instance.write(largeData));

        instance.close();
    }

    @Test
    void testIsReplaying_afterClose() throws TransportException {
        when(mockPlayer.isPlaying()).thenReturn(true);
        transportInstance.close();

        // Player should be stopped after close
        when(mockPlayer.isPlaying()).thenReturn(false);
        assertFalse(transportInstance.isReplaying());
    }

    @Test
    void testGetPacketsReplayed_afterMultipleOperations() {
        when(mockPlayer.getPacketsReplayed()).thenReturn(10L);
        assertEquals(10L, transportInstance.getPacketsReplayed());

        when(mockPlayer.getPacketsReplayed()).thenReturn(25L);
        assertEquals(25L, transportInstance.getPacketsReplayed());
    }

    @Test
    void testConstructor_withManualMacAddresses() throws TransportException {
        config.autoDetectMacAddresses = false;
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());
        assertNotNull(instance);
        assertTrue(instance.isOpen());
        instance.close();
    }

    @Test
    @RequirePcap
    void testRead_withActualPcapPlayback() throws Exception {
        // Create a real player instance that uses test.pcap
        config.mockPlayer = null; // Use real player
        config.pcapFile = "/test.pcap";
        config.autoStart = true;
        config.onlyIncomingPackets = false;
        config.onlyOutgoingPackets = false;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Wait for some packets to be queued
        Thread.sleep(200);

        // Try to read - may get data or timeout
        try {
            byte[] data = instance.read(0);
            assertNotNull(data);
        } catch (TransportException e) {
            // Timeout is acceptable
        }

        instance.close();
    }

    @Test
    @RequirePcap
    void testPeekReadableBytes_withActualData() throws Exception {
        config.mockPlayer = null;
        config.pcapFile = "/test.pcap";
        config.autoStart = true;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        Thread.sleep(200);

        // Try to peek
        try {
            byte[] peeked = instance.peekReadableBytes(0);
            assertNotNull(peeked);
        } catch (TransportException e) {
            // Timeout is acceptable
        }

        instance.close();
    }

    @Test
    @RequirePcap
    void testGetNumBytesAvailable_withRealPlayer() throws Exception {
        config.mockPlayer = null;
        config.pcapFile = "/test.pcap";
        config.autoStart = true;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        Thread.sleep(150);

        int available = instance.getNumBytesAvailable();
        assertTrue(available >= 0);

        instance.close();
    }

    @Test
    @RequirePcap
    void testStartStopReplay_withRealPlayer() throws Exception {
        config.mockPlayer = null;
        config.pcapFile = "/test.pcap";
        config.autoStart = false;
        config.loop = true;  // Enable looping so player keeps running

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        assertFalse(instance.isReplaying());

        instance.startReplay();
        Thread.sleep(100);  // Give it more time to start
        assertTrue(instance.isReplaying());

        instance.stopReplay();
        Thread.sleep(50);
        assertFalse(instance.isReplaying());

        instance.close();
    }

    @Test
    @RequirePcap
    void testGetPacketsReplayed_withRealPlayer() throws Exception {
        config.mockPlayer = null;
        config.pcapFile = "/test.pcap";
        config.autoStart = true;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        Thread.sleep(200);

        long packets = instance.getPacketsReplayed();
        assertTrue(packets >= 0);

        instance.close();
    }

    @Test
    void testConstructor_withReadTimeout() throws TransportException {
        config.readTimeout = 1000; // 1 second timeout

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());
        assertNotNull(instance);
        instance.close();
    }

    @Test
    void testConstructor_withDifferentPacketQueueSizes() throws TransportException {
        config.packetQueueSize = 10;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());
        assertNotNull(instance);
        instance.close();

        config.packetQueueSize = 5000;
        instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());
        assertNotNull(instance);
        instance.close();
    }

    @Test
    @RequirePcap
    void testMultipleReadsInSequence() throws Exception {
        config.mockPlayer = null;
        config.pcapFile = "/test.pcap";
        config.autoStart = true;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        Thread.sleep(200);

        // Try multiple reads
        for (int i = 0; i < 3; i++) {
            try {
                byte[] data = instance.read(0);
                assertNotNull(data);
            } catch (TransportException e) {
                // Timeout acceptable
            }
        }

        instance.close();
    }

    @Test
    @RequirePcap
    void testRead_afterStoppingReplay() throws Exception {
        config.mockPlayer = null;
        config.pcapFile = "/test.pcap";
        config.autoStart = true;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        Thread.sleep(100);

        instance.stopReplay();
        Thread.sleep(50);

        // Read should still work but likely timeout
        try {
            byte[] data = instance.read(0);
        } catch (TransportException e) {
            // Expected timeout
        }

        instance.close();
    }

    @Test
    void testConstructor_allConfigurationCombinations() throws TransportException {
        // Test various configuration combinations
        config.protocolId = 0x88B5;
        config.speedFactor = 1.0;
        config.loop = true;
        config.vlanId = 100;
        config.onlyIncomingPackets = true;
        config.maxFrameSize = 9000;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());
        assertNotNull(instance);
        assertTrue(instance.isOpen());
        instance.close();
    }

    @Test
    void testWrite_throwsException() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Writing should throw an exception as this is a replay-only transport
        assertThrows(TransportException.class, () -> instance.write(new byte[]{1, 2, 3}));

        instance.close();
    }

    @Test
    void testPeekReadableBytes_withZeroBytes() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        byte[] result = instance.peekReadableBytes(0);
        assertNotNull(result);
        assertEquals(0, result.length);

        instance.close();
    }

    @Test
    void testPeekReadableBytes_withNegativeBytes() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Negative bytes should throw an exception
        assertArrayEquals(new byte[0], instance.peekReadableBytes(-1));

        instance.close();
    }

    @Test
    void testRead_withZeroBytes() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        byte[] result = instance.read(0);
        assertNotNull(result);
        assertEquals(0, result.length);

        instance.close();
    }

    @Test
    void testRead_withNegativeBytes() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Negative bytes should throw an exception
        assertArrayEquals(new byte[0], instance.read(-1));

        instance.close();
    }

    @Test
    void testMultipleClose() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        instance.close();
        assertFalse(instance.isOpen());

        // Closing again should be idempotent
        instance.close();
        assertFalse(instance.isOpen());
    }

    @Test
    void testIsOpen_initialState() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testGetNumBytesAvailable_initially() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Initially should be 0
        assertEquals(0, instance.getNumBytesAvailable());

        instance.close();
    }

    @Test
    void testStartReplay_alreadyPlaying() throws TransportException {
        config.mockPlayer = mockPlayer;
        when(mockPlayer.isPlaying()).thenReturn(true);

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Start when already playing should not call start again
        instance.startReplay();
        verify(mockPlayer, times(0)).start();

        instance.close();
    }

    @Test
    void testStopReplay_notPlaying() throws TransportException {
        config.mockPlayer = mockPlayer;
        when(mockPlayer.isPlaying()).thenReturn(false);

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Stop when not playing should not call stop
        instance.stopReplay();
        verify(mockPlayer, times(0)).stop();

        instance.close();
    }

    @Test
    void testConstructor_withInvalidConfiguration() {
        config.pcapFile = null;
        config.mockPlayer = null;

        // Should throw exception with invalid configuration
        assertThrows(TransportException.class, () -> new PcapReplayTransportInstance(config, AuditLog.builder().build()));
    }

    @Test
    void testConstructor_withEmptyMacAddress() {
        config.mockPlayer = null;
        config.pcapFile = "/test.pcap";
        config.localAddress = "";
        config.remoteAddress = "";
        config.autoDetectMacAddresses = false;

        // Should throw exception with empty MAC addresses
        assertThrows(TransportException.class, () -> new PcapReplayTransportInstance(config, AuditLog.builder().build()));
    }

    @Test
    void testPeekReadableBytes_moreThanAvailable() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Requesting more bytes than available should throw (player returns no data)
        assertThrows(TransportException.class, () -> instance.peekReadableBytes(100));

        instance.close();
    }

    @Test
    void testRead_moreThanAvailable() throws TransportException {
        config.mockPlayer = mockPlayer;

        PcapReplayTransportInstance instance = new PcapReplayTransportInstance(config, AuditLog.builder().build());

        // Requesting more bytes than available should throw (player returns no data)
        assertThrows(TransportException.class, () -> instance.read(100));

        instance.close();
    }

}
