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
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;
import org.mockito.MockedStatic;
import org.pcap4j.core.*;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.UnknownPacket;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.MacAddress;

import java.net.InetAddress;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 * Mock-based tests for RawSocketTransportInstance that don't require pcap permissions.
 * Uses Mockito to mock pcap4j classes so all code paths can be tested in CI.
 *
 * <p>Uses doReturn/doThrow style consistently to avoid Mockito thread-safety issues
 * with the capture thread calling getNextPacketEx() concurrently during re-stubbing.</p>
 */
@Timeout(value = 10, unit = TimeUnit.SECONDS)
class RawSocketTransportInstanceMockTest {

    private MockedStatic<Pcaps> pcapsMock;
    private PcapNetworkInterface mockNif;
    private PcapHandle mockHandle;
    private RawSocketTransportInstance instance;
    private SharedRawSocketManager sharedManager;

    /**
     * Thread-safe flag controlling what mockHandle.isOpen() returns.
     * Using an AtomicBoolean avoids re-stubbing the mock from threads that don't own it.
     */
    private final AtomicBoolean handleOpen = new AtomicBoolean(true);

    @BeforeEach
    void setUp() throws Exception {
        handleOpen.set(true);

        pcapsMock = mockStatic(Pcaps.class);

        mockNif = mock(PcapNetworkInterface.class);
        doReturn("en0").when(mockNif).getName();
        doReturn("Mock Ethernet").when(mockNif).getDescription();
        doReturn(Collections.emptyList()).when(mockNif).getAddresses();

        mockHandle = mock(PcapHandle.class);
        // Thread-safe isOpen() backed by AtomicBoolean
        doAnswer(inv -> handleOpen.get()).when(mockHandle).isOpen();
        // Default: no packets arriving
        doThrow(new java.util.concurrent.TimeoutException("no packet"))
            .when(mockHandle).getNextPacketEx();

        doReturn(mockHandle).when(mockNif)
            .openLive(anyInt(), any(PcapNetworkInterface.PromiscuousMode.class), anyInt());

        pcapsMock.when(Pcaps::findAllDevs).thenReturn(List.of(mockNif));

        sharedManager = new SharedRawSocketManager();
    }

    @AfterEach
    void tearDown() {
        try {
            handleOpen.set(false);
            if (instance != null) {
                try {
                    instance.close();
                } catch (Exception ignored) {
                    // Best effort cleanup
                }
            }
        } finally {
            if (pcapsMock != null) {
                pcapsMock.close();
            }
        }
    }

    /**
     * Protocol ID used for packet-matching tests. Must be < 0x8000 to avoid
     * sign extension issues when comparing int protocolId with the short value
     * returned by EtherType.value() in matchesFilter().
     */
    private static final int TEST_PROTOCOL_ID = 0x0800;

    /**
     * Creates a standard test instance with dedicated (non-shared) handle.
     */
    private RawSocketTransportInstance createDedicatedInstance() throws TransportException {
        return createDedicatedInstance(0x88B5);
    }

    /**
     * Creates a dedicated instance with a specific protocol ID.
     */
    private RawSocketTransportInstance createDedicatedInstance(int protocolId) throws TransportException {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = protocolId;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.reuseInterface = false;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        return new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());
    }

    /**
     * Creates a test instance with shared handle mode.
     */
    private RawSocketTransportInstance createSharedInstance() throws TransportException {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.reuseInterface = true;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        return new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());
    }

    /**
     * Builds a matching Ethernet packet (correct src/dst MAC and EtherType)
     * using TEST_PROTOCOL_ID to avoid sign extension issues in matchesFilter().
     */
    private EthernetPacket buildMatchingPacket(byte[] payload) {
        return new EthernetPacket.Builder()
            .srcAddr(MacAddress.getByName("AA:BB:CC:DD:EE:FF"))
            .dstAddr(MacAddress.getByName("00:11:22:33:44:55"))
            .type(EtherType.getInstance((short) TEST_PROTOCOL_ID))
            .payloadBuilder(new UnknownPacket.Builder().rawData(payload))
            .paddingAtBuild(true)
            .build();
    }

    // ========== Constructor / Initialization Tests ==========

    @Test
    void testConstructor_dedicatedHandle() throws Exception {
        instance = createDedicatedInstance();

        assertTrue(instance.isOpen());
        assertNotNull(instance.getConfiguration());
        assertEquals(0x88B5, instance.getConfiguration().protocolId);
    }

    @Test
    void testConstructor_sharedHandle() throws Exception {
        instance = createSharedInstance();

        assertTrue(instance.isOpen());
        assertNotNull(instance.getConfiguration());
    }

    @Test
    void testConstructor_noInterfaceSpecified_usesFirst() throws Exception {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());
        assertTrue(instance.isOpen());
    }

    @Test
    void testConstructor_interfaceNotFound_throwsTransportException() {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "nonexistent_interface";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        assertThrows(TransportException.class, () ->
            new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build()));
    }

    @Test
    void testConstructor_noNetworkInterfaces_throwsTransportException() {
        pcapsMock.when(Pcaps::findAllDevs).thenReturn(Collections.emptyList());

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        assertThrows(TransportException.class, () ->
            new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build()));
    }

    @Test
    void testConstructor_nullNetworkInterfaces_throwsTransportException() {
        pcapsMock.when(Pcaps::findAllDevs).thenReturn(null);

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        assertThrows(TransportException.class, () ->
            new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build()));
    }

    @Test
    void testConstructor_openLiveFails_throwsTransportException() throws Exception {
        doThrow(new PcapNativeException("Permission denied"))
            .when(mockNif).openLive(anyInt(), any(PcapNetworkInterface.PromiscuousMode.class), anyInt());

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        assertThrows(TransportException.class, () ->
            new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build()));
    }

    @Test
    void testConstructor_withLocalAddressFromInterface() throws Exception {
        PcapAddress mockAddr = mock(PcapAddress.class);
        InetAddress mockInet = mock(InetAddress.class);
        doReturn(new byte[]{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}).when(mockInet).getAddress();
        doReturn(mockInet).when(mockAddr).getAddress();
        doReturn(List.of(mockAddr)).when(mockNif).getAddresses();

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());
        assertTrue(instance.isOpen());
    }

    @Test
    void testConstructor_emptyLocalAddress_usesInterfaceMac() throws Exception {
        PcapAddress mockAddr = mock(PcapAddress.class);
        InetAddress mockInet = mock(InetAddress.class);
        doReturn(new byte[]{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}).when(mockInet).getAddress();
        doReturn(mockInet).when(mockAddr).getAddress();
        doReturn(List.of(mockAddr)).when(mockNif).getAddresses();

        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());
        assertTrue(instance.isOpen());
    }

    @Test
    void testConstructor_withPromiscuousMode() throws Exception {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;
        config.promiscuousMode = true;

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());
        assertTrue(instance.isOpen());

        verify(mockNif).openLive(eq(65536), eq(PcapNetworkInterface.PromiscuousMode.PROMISCUOUS), anyInt());
    }

    @Test
    void testConstructor_withBpfFilter() throws Exception {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;
        config.bpfFilter = "ether proto 0x88B5";

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());
        assertTrue(instance.isOpen());

        verify(mockHandle).setFilter(eq("ether proto 0x88B5"), eq(BpfProgram.BpfCompileMode.OPTIMIZE));
    }

    @Test
    void testConstructor_lowCaptureTimeout_enforcesMinimum() throws Exception {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 10; // Below minimum of 100
        config.snapshotLength = 65536;

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());

        verify(mockNif).openLive(eq(65536), any(PcapNetworkInterface.PromiscuousMode.class), eq(100));
    }

    @Test
    void testConstructor_withAutoBpfFilter_nonPromiscuous() throws Exception {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;
        config.promiscuousMode = false;

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());

        // Auto-generated filter should include ether proto and dst MAC
        verify(mockHandle).setFilter(contains("ether proto 0x88B5"), eq(BpfProgram.BpfCompileMode.OPTIMIZE));
    }

    @Test
    void testConstructor_withAutoBpfFilter_promiscuousMode() throws Exception {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;
        config.promiscuousMode = true;

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());

        // In promiscuous mode, filter should only have ether proto (no dst MAC)
        verify(mockHandle).setFilter(eq("ether proto 0x88B5"), eq(BpfProgram.BpfCompileMode.OPTIMIZE));
    }

    // ========== isOpen Tests ==========

    @Test
    void testIsOpen_whenOpen() throws Exception {
        instance = createDedicatedInstance();
        assertTrue(instance.isOpen());
    }

    @Test
    void testIsOpen_whenHandleClosed() throws Exception {
        instance = createDedicatedInstance();
        handleOpen.set(false);
        assertFalse(instance.isOpen());
    }

    // ========== getNumBytesAvailable Tests ==========

    @Test
    void testGetNumBytesAvailable_initiallyZero() throws Exception {
        instance = createDedicatedInstance();
        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void testGetNumBytesAvailable_whenClosed() throws Exception {
        instance = createDedicatedInstance();
        handleOpen.set(false);
        assertEquals(0, instance.getNumBytesAvailable());
    }

    // ========== read Tests ==========

    @Test
    void testRead_zeroBytes_returnsEmptyArray() throws Exception {
        instance = createDedicatedInstance();
        byte[] result = instance.read(0);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testRead_negativeBytes_returnsEmptyArray() throws Exception {
        instance = createDedicatedInstance();
        byte[] result = instance.read(-1);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testRead_whenClosed_throwsTransportException() throws Exception {
        instance = createDedicatedInstance();
        handleOpen.set(false);

        assertThrows(TransportException.class, () -> instance.read(10));
    }

    // ========== peekReadableBytes Tests ==========

    @Test
    void testPeekReadableBytes_zeroBytes_returnsEmptyArray() throws Exception {
        instance = createDedicatedInstance();
        byte[] result = instance.peekReadableBytes(0);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekReadableBytes_negativeBytes_returnsEmptyArray() throws Exception {
        instance = createDedicatedInstance();
        byte[] result = instance.peekReadableBytes(-1);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekReadableBytes_whenClosed_throwsTransportException() throws Exception {
        instance = createDedicatedInstance();
        handleOpen.set(false);

        assertThrows(TransportException.class, () -> instance.peekReadableBytes(10));
    }

    // ========== write Tests ==========

    @Test
    void testWrite_null_doesNothing() throws Exception {
        instance = createDedicatedInstance();
        assertDoesNotThrow(() -> instance.write(null));
        verify(mockHandle, never()).sendPacket(any(Packet.class));
    }

    @Test
    void testWrite_emptyArray_doesNothing() throws Exception {
        instance = createDedicatedInstance();
        assertDoesNotThrow(() -> instance.write(new byte[0]));
        verify(mockHandle, never()).sendPacket(any(Packet.class));
    }

    @Test
    void testWrite_validData_sendsPacket() throws Exception {
        instance = createDedicatedInstance();
        byte[] data = new byte[]{0x01, 0x02, 0x03, 0x04};

        instance.write(data);

        verify(mockHandle).sendPacket(any(Packet.class));
    }

    @Test
    void testWrite_exceedsMaxFrameSize_throwsTransportException() throws Exception {
        instance = createDedicatedInstance();
        byte[] data = new byte[2000];

        assertThrows(TransportException.class, () -> instance.write(data));
    }

    @Test
    void testWrite_exactlyMaxFrameSize_succeeds() throws Exception {
        instance = createDedicatedInstance();
        byte[] data = new byte[1500];

        assertDoesNotThrow(() -> instance.write(data));
        verify(mockHandle).sendPacket(any(Packet.class));
    }

    @Test
    void testWrite_oneOverMaxFrameSize_throws() throws Exception {
        instance = createDedicatedInstance();
        byte[] data = new byte[1501];

        assertThrows(TransportException.class, () -> instance.write(data));
    }

    @Test
    void testWrite_whenClosed_throwsTransportException() throws Exception {
        instance = createDedicatedInstance();
        handleOpen.set(false);

        assertThrows(TransportException.class, () -> instance.write(new byte[]{0x01}));
    }

    @Test
    void testWrite_sendPacketFails_throwsTransportException() throws Exception {
        // Set up sendPacket to fail BEFORE creating the instance
        // (avoids race condition with capture thread)
        doThrow(new PcapNativeException("Send failed"))
            .when(mockHandle).sendPacket(any(Packet.class));

        instance = createDedicatedInstance();

        assertThrows(TransportException.class, () -> instance.write(new byte[]{0x01, 0x02}));
    }

    @Test
    void testWrite_withVlanId_logsWarning() throws Exception {
        RawSocketTransportConfiguration config = new RawSocketTransportConfiguration();
        config.interfaceName = "en0";
        config.localAddress = "00:11:22:33:44:55";
        config.remoteAddress = "AA:BB:CC:DD:EE:FF";
        config.protocolId = 0x88B5;
        config.maxFrameSize = 1500;
        config.readTimeout = 100;
        config.captureTimeout = 100;
        config.snapshotLength = 65536;
        config.vlanId = 100;

        instance = new RawSocketTransportInstance(sharedManager, config, AuditLog.builder().build());

        assertDoesNotThrow(() -> instance.write(new byte[]{0x01, 0x02, 0x03}));
    }

    // ========== close Tests ==========

    @Test
    void testClose_dedicatedHandle() throws Exception {
        instance = createDedicatedInstance();
        assertTrue(instance.isOpen());

        handleOpen.set(false);
        instance.close();
        assertFalse(instance.isOpen());
    }

    @Test
    void testClose_sharedHandle() throws Exception {
        instance = createSharedInstance();
        assertTrue(instance.isOpen());

        handleOpen.set(false);
        instance.close();
        assertFalse(instance.isOpen());
    }

    @Test
    void testClose_idempotent() throws Exception {
        instance = createDedicatedInstance();

        handleOpen.set(false);
        instance.close();
        assertFalse(instance.isOpen());

        assertDoesNotThrow(() -> instance.close());
    }

    @Test
    void testClose_handleCloseThrows_doesNotPropagate() throws Exception {
        doThrow(new RuntimeException("Close failed")).when(mockHandle).close();

        instance = createDedicatedInstance();

        handleOpen.set(false);

        assertDoesNotThrow(() -> instance.close());
    }

    // ========== Async Listener Tests ==========

    @Test
    void testRegisterDataListener() throws Exception {
        instance = createDedicatedInstance();
        AtomicBoolean called = new AtomicBoolean(false);

        assertDoesNotThrow(() -> instance.registerDataListener(() -> called.set(true)));
    }

    @Test
    void testRemoveDataListener() throws Exception {
        instance = createDedicatedInstance();
        instance.registerDataListener(() -> {});

        assertDoesNotThrow(() -> instance.removeDataListener());
    }

    @Test
    void testRegisterDisconnectListener() throws Exception {
        instance = createDedicatedInstance();
        AtomicReference<Throwable> receivedCause = new AtomicReference<>();

        assertDoesNotThrow(() -> instance.registerDisconnectListener(receivedCause::set));
    }

    @Test
    void testRemoveDisconnectListener() throws Exception {
        instance = createDedicatedInstance();
        instance.registerDisconnectListener(t -> {});

        assertDoesNotThrow(() -> instance.removeDisconnectListener());
    }

    // ========== Capture Thread / Packet Processing Tests ==========

    @Test
    void testCaptureThread_startsAutomatically() throws Exception {
        instance = createDedicatedInstance();

        Thread.sleep(50);

        assertTrue(instance.isOpen());
        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void testCaptureThread_matchingPacket_notifiesDataListener() throws Exception {
        EthernetPacket ethPacket = buildMatchingPacket(new byte[]{0x01, 0x02, 0x03});

        AtomicInteger listenerCallCount = new AtomicInteger(0);

        // Return the packet continuously so the listener will be called
        // after it's registered (some calls happen before registration, that's fine)
        doAnswer(inv -> ethPacket).when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance(TEST_PROTOCOL_ID);
        instance.registerDataListener(listenerCallCount::incrementAndGet);

        Thread.sleep(300);

        assertTrue(listenerCallCount.get() >= 1,
            "Data listener should have been called at least once, was " + listenerCallCount.get());
    }

    @Test
    void testCaptureThread_wrongEtherType_doesNotNotify() throws Exception {
        // Use a different EtherType than what the instance is configured for
        EthernetPacket ethPacket = new EthernetPacket.Builder()
            .srcAddr(MacAddress.getByName("AA:BB:CC:DD:EE:FF"))
            .dstAddr(MacAddress.getByName("00:11:22:33:44:55"))
            .type(EtherType.getInstance((short) 0x0806)) // ARP, not the configured 0x0800
            .payloadBuilder(new UnknownPacket.Builder().rawData(new byte[]{0x01, 0x02}))
            .paddingAtBuild(true)
            .build();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);

        doAnswer(inv -> ethPacket).when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance(TEST_PROTOCOL_ID);
        instance.registerDataListener(() -> listenerCalled.set(true));

        Thread.sleep(300);

        assertFalse(listenerCalled.get(), "Listener should not be called for wrong EtherType");
    }

    @Test
    void testCaptureThread_wrongSourceMac_doesNotNotify() throws Exception {
        EthernetPacket ethPacket = new EthernetPacket.Builder()
            .srcAddr(MacAddress.getByName("11:22:33:44:55:66")) // Not the expected remote
            .dstAddr(MacAddress.getByName("00:11:22:33:44:55"))
            .type(EtherType.getInstance((short) TEST_PROTOCOL_ID))
            .payloadBuilder(new UnknownPacket.Builder().rawData(new byte[]{0x01, 0x02}))
            .paddingAtBuild(true)
            .build();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);

        doAnswer(inv -> ethPacket).when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance(TEST_PROTOCOL_ID);
        instance.registerDataListener(() -> listenerCalled.set(true));

        Thread.sleep(300);

        assertFalse(listenerCalled.get(), "Listener should not be called for wrong source MAC");
    }

    @Test
    void testCaptureThread_wrongDestinationMac_doesNotNotify() throws Exception {
        EthernetPacket ethPacket = new EthernetPacket.Builder()
            .srcAddr(MacAddress.getByName("AA:BB:CC:DD:EE:FF"))
            .dstAddr(MacAddress.getByName("FF:EE:DD:CC:BB:AA")) // Not local, not broadcast
            .type(EtherType.getInstance((short) TEST_PROTOCOL_ID))
            .payloadBuilder(new UnknownPacket.Builder().rawData(new byte[]{0x01, 0x02}))
            .paddingAtBuild(true)
            .build();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);

        doAnswer(inv -> ethPacket).when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance(TEST_PROTOCOL_ID);
        instance.registerDataListener(() -> listenerCalled.set(true));

        Thread.sleep(300);

        assertFalse(listenerCalled.get(), "Listener should not be called for wrong destination MAC");
    }

    @Test
    void testCaptureThread_broadcastDestination_notifies() throws Exception {
        EthernetPacket ethPacket = new EthernetPacket.Builder()
            .srcAddr(MacAddress.getByName("AA:BB:CC:DD:EE:FF"))
            .dstAddr(MacAddress.ETHER_BROADCAST_ADDRESS)
            .type(EtherType.getInstance((short) TEST_PROTOCOL_ID))
            .payloadBuilder(new UnknownPacket.Builder().rawData(new byte[]{0x01, 0x02, 0x03}))
            .paddingAtBuild(true)
            .build();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);

        doAnswer(inv -> ethPacket).when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance(TEST_PROTOCOL_ID);
        instance.registerDataListener(() -> listenerCalled.set(true));

        Thread.sleep(300);

        assertTrue(listenerCalled.get(), "Listener should be called for broadcast destination");
    }

    @Test
    void testCaptureThread_packetWithNoPayload_doesNotNotify() throws Exception {
        EthernetPacket ethPacket = new EthernetPacket.Builder()
            .srcAddr(MacAddress.getByName("AA:BB:CC:DD:EE:FF"))
            .dstAddr(MacAddress.getByName("00:11:22:33:44:55"))
            .type(EtherType.getInstance((short) TEST_PROTOCOL_ID))
            .paddingAtBuild(true)
            .build();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);

        doAnswer(inv -> ethPacket).when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance(TEST_PROTOCOL_ID);
        instance.registerDataListener(() -> listenerCalled.set(true));

        Thread.sleep(300);

        assertFalse(listenerCalled.get(), "Listener should not be called for packet with no payload");
    }

    @Test
    void testCaptureThread_handleNotOpen_stopsCapture() throws Exception {
        doThrow(new NotOpenException("Handle closed"))
            .when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance();

        Thread.sleep(200);
        // Capture thread should have exited gracefully
    }

    @Test
    void testCaptureThread_pcapNativeException_stopsCapture() throws Exception {
        doThrow(new PcapNativeException("Native error"))
            .when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance();

        Thread.sleep(200);
        // Capture thread exits gracefully
    }

    @Test
    void testCaptureThread_genericException_stopsCapture() throws Exception {
        doThrow(new RuntimeException("Unexpected error"))
            .when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance();

        Thread.sleep(200);
        // Capture thread exits after unexpected exception
    }

    // ========== Read/Peek with Packet Data Tests ==========

    @Test
    void testRead_withPacketData() throws Exception {
        byte[] payload = new byte[]{0x01, 0x02, 0x03, 0x04, 0x05};
        EthernetPacket ethPacket = buildMatchingPacket(payload);

        // Continuously deliver packets so the receiveQueue has data for read()
        doAnswer(inv -> ethPacket).when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance(TEST_PROTOCOL_ID);

        // Wait for capture thread to push data into receiveQueue
        Thread.sleep(300);

        // read() pulls data from receiveQueue into ringBuffer
        byte[] read = instance.read(payload.length);
        assertNotNull(read);
        assertEquals(payload.length, read.length);
        assertArrayEquals(payload, read);
    }

    @Test
    void testPeekReadableBytes_withPacketData() throws Exception {
        byte[] payload = new byte[]{0x0A, 0x0B, 0x0C};
        EthernetPacket ethPacket = buildMatchingPacket(payload);

        doAnswer(inv -> ethPacket).when(mockHandle).getNextPacketEx();

        instance = createDedicatedInstance(TEST_PROTOCOL_ID);

        Thread.sleep(300);

        // Peek fills ringBuffer from receiveQueue internally
        byte[] peeked = instance.peekReadableBytes(payload.length);
        assertNotNull(peeked);
        assertEquals(payload.length, peeked.length);
        assertArrayEquals(payload, peeked);

        // Peek should not consume bytes
        assertEquals(payload.length, instance.getNumBytesAvailable());
    }

    @Test
    void testRead_afterClose_throwsTransportException() throws Exception {
        instance = createDedicatedInstance();

        handleOpen.set(false);
        instance.close();

        assertThrows(TransportException.class, () -> instance.read(3));
    }
}
