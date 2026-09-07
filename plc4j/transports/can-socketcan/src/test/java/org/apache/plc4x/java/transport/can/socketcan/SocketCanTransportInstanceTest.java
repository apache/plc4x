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
package org.apache.plc4x.java.transport.can.socketcan;

import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.can.socketcan.config.SocketCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import tel.schich.javacan.RawCanChannel;

import java.io.IOException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicReference;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for {@link SocketCanTransportInstance}.
 * <p>
 * Uses a testable subclass that bypasses the Linux platform check and uses
 * a mock CAN channel, allowing the transport logic to be tested on any platform.
 */
class SocketCanTransportInstanceTest {

    private AuditLog auditLog;
    private SharedCanManager sharedCanManager;
    private RawCanChannel mockChannel;
    private TestableSocketCanTransportInstance instance;

    @BeforeEach
    void setUp() {
        auditLog = Mockito.mock(AuditLog.class);
        sharedCanManager = new SharedCanManager();
        mockChannel = Mockito.mock(RawCanChannel.class);
    }

    @AfterEach
    void tearDown() throws Exception {
        if (instance != null && instance.isOpen()) {
            instance.close();
        }
    }

    private TestableSocketCanTransportInstance createInstance(SocketCanTransportConfiguration config) throws TransportException {
        instance = new TestableSocketCanTransportInstance(sharedCanManager, config, auditLog, mockChannel);
        return instance;
    }

    private SocketCanTransportConfiguration createConfig(String interfaceName) {
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();
        config.interfaceName = interfaceName;
        return config;
    }

    @Test
    void constructorRejectsNonLinuxPlatform() {
        String osName = System.getProperty("os.name", "").toLowerCase();
        if (osName.contains("linux")) {
            return; // Skip on Linux
        }

        SocketCanTransportConfiguration config = createConfig("can0");

        TransportException ex = assertThrows(TransportException.class, () ->
                new SocketCanTransportInstance(sharedCanManager, config, auditLog));
        assertTrue(ex.getMessage().contains("Linux"));
    }

    @Test
    void constructorRejectsNullInterfaceName() {
        SocketCanTransportConfiguration config = createConfig(null);

        assertThrows(TransportException.class, () -> createInstance(config));
    }

    @Test
    void constructorRejectsEmptyInterfaceName() {
        SocketCanTransportConfiguration config = createConfig("  ");

        assertThrows(TransportException.class, () -> createInstance(config));
    }

    @Test
    void isOpenAfterConstruction() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        assertTrue(instance.isOpen());
    }

    @Test
    void closeMarksNotOpen() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.close();
        assertFalse(instance.isOpen());
    }

    @Test
    void closeNotifiesDisconnectListener() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        AtomicReference<Throwable> receivedCause = new AtomicReference<>(new RuntimeException("not called"));
        instance.registerDisconnectListener(receivedCause::set);

        instance.close();

        // Graceful close passes null
        assertNull(receivedCause.get());
    }

    @Test
    void doubleCloseIsNoOp() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.close();
        assertDoesNotThrow(() -> instance.close());
    }

    @Test
    void readAfterCloseThrows() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.close();
        assertThrows(TransportException.class, () -> instance.read(16));
    }

    @Test
    void writeAfterCloseThrows() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.close();
        assertThrows(TransportException.class, () -> instance.write(new byte[16]));
    }

    @Test
    void peekAfterCloseThrows() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.close();
        assertThrows(TransportException.class, () -> instance.peekReadableBytes(16));
    }

    @Test
    void getNumBytesAvailableAfterCloseThrows() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.close();
        assertThrows(TransportException.class, () -> instance.getNumBytesAvailable());
    }

    @Test
    void registerAndRemoveDataListener() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        AtomicBoolean called = new AtomicBoolean(false);
        instance.registerDataListener(() -> called.set(true));
        instance.removeDataListener();

        // After removal, listener should be null
        assertDoesNotThrow(() -> instance.close());
    }

    @Test
    void registerAndRemoveDisconnectListener() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.registerDisconnectListener(cause -> {});
        instance.removeDisconnectListener();

        // After removal, close should not throw
        assertDoesNotThrow(() -> instance.close());
    }

    @Test
    void writeTooShortFrameThrows() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        assertThrows(TransportException.class, () -> instance.write(new byte[8]));
    }

    @Test
    void getNumBytesAvailableInitiallyZero() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void peekReturnsEmptyWhenNoData() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        byte[] peeked = instance.peekReadableBytes(16);
        assertEquals(0, peeked.length);
    }

    @Test
    void readReturnsEmptyWhenNoData() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        byte[] read = instance.read(16);
        assertEquals(0, read.length);
    }

    @Test
    void configurationFilterIsBuilt() {
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();
        config.interfaceName = "can0";
        config.filterIds = "0x100,0x200";

        var filter = config.buildFilter();
        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x200));
        assertFalse(filter.matches(0x300));
    }

    @Test
    void writeValidFrameCallsChannel() throws Exception {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        // Build a valid 16-byte CAN frame: 4 bytes ID + 1 byte DLC + 3 padding + 8 data
        byte[] frameBytes = new byte[16];
        frameBytes[0] = 0x00;
        frameBytes[1] = 0x00;
        frameBytes[2] = 0x01;
        frameBytes[3] = 0x23; // ID = 0x123
        frameBytes[4] = 0x03; // DLC = 3
        frameBytes[8] = 0x01;
        frameBytes[9] = 0x02;
        frameBytes[10] = 0x03;

        assertDoesNotThrow(() -> instance.write(frameBytes));
        assertTrue(instance.sendFrameCalled);
        assertEquals(0x123, instance.lastSentId);
    }

    @Test
    void writeFrameWithIOExceptionThrowsTransportException() throws Exception {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);
        instance.sendFrameThrows = new IOException("Socket error");

        byte[] frameBytes = new byte[16];
        frameBytes[4] = 0x01; // DLC = 1

        assertThrows(TransportException.class, () -> instance.write(frameBytes));
    }

    @Test
    void writeFrameWithZeroDlc() throws Exception {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        byte[] frameBytes = new byte[16];
        // DLC = 0, no data

        assertDoesNotThrow(() -> instance.write(frameBytes));
        assertTrue(instance.sendFrameCalled);
        assertEquals(0, instance.lastSentData.length);
    }

    @Test
    void writeFrameWithMaxDlc() throws Exception {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        byte[] frameBytes = new byte[16];
        frameBytes[4] = 0x08; // DLC = 8
        for (int i = 0; i < 8; i++) {
            frameBytes[8 + i] = (byte) (i + 1);
        }

        assertDoesNotThrow(() -> instance.write(frameBytes));
        assertTrue(instance.sendFrameCalled);
        assertEquals(8, instance.lastSentData.length);
    }

    @Test
    void closeWithMockChannelDoesNotThrow() throws Exception {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        assertDoesNotThrow(() -> instance.close());
        assertFalse(instance.isOpen());
    }

    @Test
    void closeWithChannelIOExceptionDoesNotPropagate() throws Exception {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        Mockito.doThrow(new IOException("close error")).when(mockChannel).close();

        // Should log warning but not throw
        assertDoesNotThrow(() -> instance.close());
    }

    @Test
    void disconnectListenerExceptionIsSwallowed() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.registerDisconnectListener(cause -> {
            throw new RuntimeException("listener error");
        });

        // close() calls notifyDisconnect which should catch listener exceptions
        assertDoesNotThrow(() -> instance.close());
    }

    @Test
    void processReceivedFrameBuffersMatchingFrame() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.processReceivedFrame(0x123, new byte[]{0x01, 0x02, 0x03});

        assertEquals(16, instance.getNumBytesAvailable());
        byte[] read = instance.read(16);
        assertEquals(0x123, ((read[0] & 0xFF) << 24) | ((read[1] & 0xFF) << 16)
                | ((read[2] & 0xFF) << 8) | (read[3] & 0xFF));
        assertEquals(3, read[4]); // DLC
        assertEquals(0x01, read[8]);
        assertEquals(0x02, read[9]);
        assertEquals(0x03, read[10]);
    }

    @Test
    void processReceivedFrameFiltersNonMatchingIds() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        config.filterIds = "0x100";
        createInstance(config);

        // Send a frame with non-matching ID
        instance.processReceivedFrame(0x200, new byte[]{0x01});

        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void processReceivedFrameNotifiesDataListener() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        AtomicBoolean listenerCalled = new AtomicBoolean(false);
        instance.registerDataListener(() -> listenerCalled.set(true));

        instance.processReceivedFrame(0x100, new byte[]{0x01});

        assertTrue(listenerCalled.get());
    }

    @Test
    void processReceivedFrameWithNoListenerDoesNotThrow() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        assertDoesNotThrow(() -> instance.processReceivedFrame(0x100, new byte[]{0x01}));
    }

    @Test
    void processReceivedFrameEmptyData() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.processReceivedFrame(0x100, new byte[0]);

        assertEquals(16, instance.getNumBytesAvailable());
        byte[] read = instance.read(16);
        assertEquals(0, read[4]); // DLC = 0
    }

    @Test
    void peekDoesNotConsumeData() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.processReceivedFrame(0x100, new byte[]{0x01});

        byte[] peeked = instance.peekReadableBytes(16);
        assertEquals(16, peeked.length);
        // Data should still be available after peek
        assertEquals(16, instance.getNumBytesAvailable());
    }

    @Test
    void multipleFramesBuffered() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        instance.processReceivedFrame(0x100, new byte[]{0x01});
        instance.processReceivedFrame(0x200, new byte[]{0x02});

        assertEquals(32, instance.getNumBytesAvailable());
    }

    @Test
    void getConfigurationReturnsOriginal() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        assertSame(config, instance.getConfiguration());
    }

    @Test
    void writeFrameExtractsCorrectIdAndData() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        createInstance(config);

        byte[] frameBytes = new byte[16];
        // ID = 0x1ABCDEF0 (extended range)
        frameBytes[0] = 0x1A;
        frameBytes[1] = (byte) 0xBC;
        frameBytes[2] = (byte) 0xDE;
        frameBytes[3] = (byte) 0xF0;
        frameBytes[4] = 0x04; // DLC = 4
        frameBytes[8] = 0x11;
        frameBytes[9] = 0x22;
        frameBytes[10] = 0x33;
        frameBytes[11] = 0x44;

        assertDoesNotThrow(() -> instance.write(frameBytes));
        assertEquals(0x1ABCDEF0, instance.lastSentId);
        assertArrayEquals(new byte[]{0x11, 0x22, 0x33, 0x44}, instance.lastSentData);
    }

    @Test
    void processReceivedFrameWithFilterAcceptsMatchingIds() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        config.filterIds = "0x100,0x200";
        createInstance(config);

        instance.processReceivedFrame(0x100, new byte[]{0x01});
        instance.processReceivedFrame(0x200, new byte[]{0x02});
        instance.processReceivedFrame(0x300, new byte[]{0x03}); // filtered out

        assertEquals(32, instance.getNumBytesAvailable()); // Only 2 frames (16 bytes each)
    }

    @Test
    void processReceivedFrameWithRangeFilter() throws TransportException {
        SocketCanTransportConfiguration config = createConfig("can0");
        config.filterRangeStart = 0x100;
        config.filterRangeEnd = 0x1FF;
        createInstance(config);

        instance.processReceivedFrame(0x150, new byte[]{0x01}); // in range
        instance.processReceivedFrame(0x050, new byte[]{0x02}); // out of range
        instance.processReceivedFrame(0x1FF, new byte[]{0x03}); // boundary

        assertEquals(32, instance.getNumBytesAvailable()); // 2 frames accepted
    }

    /**
     * Testable subclass that bypasses platform validation and uses a mock channel.
     */
    static class TestableSocketCanTransportInstance extends SocketCanTransportInstance {

        private final RawCanChannel mockChannel;
        boolean sendFrameCalled;
        int lastSentId;
        byte[] lastSentData;
        IOException sendFrameThrows;

        TestableSocketCanTransportInstance(SharedCanManager sharedCanManager,
                                            SocketCanTransportConfiguration configuration,
                                            AuditLog auditLog,
                                            RawCanChannel mockChannel) throws TransportException {
            super(sharedCanManager, configuration, auditLog);
            this.mockChannel = mockChannel;
        }

        @Override
        protected void validatePlatform() {
            // Skip platform check in tests
        }

        @Override
        protected RawCanChannel openChannel(String interfaceName) {
            return mockChannel;
        }

        @Override
        protected void startReaderThread(String interfaceName) {
            // Don't start background thread in tests
        }

        @Override
        protected void sendFrame(int id, byte[] data) throws IOException {
            if (sendFrameThrows != null) {
                throw sendFrameThrows;
            }
            sendFrameCalled = true;
            lastSentId = id;
            lastSentData = data;
        }
    }
}
