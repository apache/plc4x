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
package org.apache.plc4x.java.transport.can.virtualcan;

import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.can.virtualcan.config.VirtualCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicReference;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for {@link VirtualCanTransportInstance}: frame round-trips, self-exclusion,
 * lifecycle management, listener notification, and bus isolation.
 */
class VirtualCanTransportInstanceTest {

    private AuditLog auditLog;
    private final List<VirtualCanTransportInstance> openInstances = new ArrayList<>();

    @BeforeEach
    void setUp() {
        VirtualCanBusManager.reset();
        auditLog = AuditLog.builder().build();
    }

    @AfterEach
    void tearDown() {
        // Close all instances that tests may have left open
        for (VirtualCanTransportInstance instance : openInstances) {
            if (instance.isOpen()) {
                instance.close();
            }
        }
        openInstances.clear();
        VirtualCanBusManager.reset();
    }

    /**
     * Helper to create and track an instance for automatic cleanup.
     */
    private VirtualCanTransportInstance createInstance(String busName) {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.busName = busName;
        VirtualCanTransportInstance instance = new VirtualCanTransportInstance(config, auditLog);
        openInstances.add(instance);
        return instance;
    }

    @Test
    void testFrameRoundTripBetweenTwoInstances() throws Exception {
        VirtualCanTransportInstance sender = createInstance("bus1");
        VirtualCanTransportInstance receiver = createInstance("bus1");

        byte[] frame = {0x10, 0x20, 0x30, 0x40, 0x50};
        sender.write(frame);

        assertEquals(5, receiver.getNumBytesAvailable());
        byte[] received = receiver.read(5);
        assertArrayEquals(frame, received);
    }

    @Test
    void testSelfExclusion() throws Exception {
        VirtualCanTransportInstance instance = createInstance("bus1");

        byte[] frame = {0x01, 0x02, 0x03};
        instance.write(frame);

        // The sender should not receive its own frame
        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void testCloseSetIsOpenFalse() {
        VirtualCanTransportInstance instance = createInstance("bus1");
        assertTrue(instance.isOpen());

        instance.close();
        assertFalse(instance.isOpen());
    }

    @Test
    void testDoubleCloseIsHarmless() {
        VirtualCanTransportInstance instance = createInstance("bus1");
        instance.close();
        // Second close should not throw
        assertDoesNotThrow(instance::close);
        assertFalse(instance.isOpen());
    }

    @Test
    void testWriteAfterCloseThrows() {
        VirtualCanTransportInstance instance = createInstance("bus1");
        instance.close();

        assertThrows(TransportException.class, () -> instance.write(new byte[]{1, 2, 3}));
    }

    @Test
    void testReadAfterCloseThrows() {
        VirtualCanTransportInstance instance = createInstance("bus1");
        instance.close();

        assertThrows(TransportException.class, () -> instance.read(1));
    }

    @Test
    void testPeekAfterCloseThrows() {
        VirtualCanTransportInstance instance = createInstance("bus1");
        instance.close();

        assertThrows(TransportException.class, () -> instance.peekReadableBytes(1));
    }

    @Test
    void testGetNumBytesAvailableAfterCloseReturnsZero() throws Exception {
        VirtualCanTransportInstance instance = createInstance("bus1");
        instance.close();

        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void testDisconnectListenerCalledOnClose() {
        VirtualCanTransportInstance instance = createInstance("bus1");

        AtomicBoolean listenerCalled = new AtomicBoolean(false);
        AtomicReference<Throwable> receivedCause = new AtomicReference<>(new RuntimeException("sentinel"));

        instance.registerDisconnectListener(cause -> {
            listenerCalled.set(true);
            receivedCause.set(cause);
        });

        instance.close();

        assertTrue(listenerCalled.get());
        // Graceful close passes null
        assertNull(receivedCause.get());
    }

    @Test
    void testDisconnectListenerNotCalledAfterRemoval() {
        VirtualCanTransportInstance instance = createInstance("bus1");

        AtomicBoolean listenerCalled = new AtomicBoolean(false);
        instance.registerDisconnectListener(cause -> listenerCalled.set(true));
        instance.removeDisconnectListener();

        instance.close();

        assertFalse(listenerCalled.get());
    }

    @Test
    void testDifferentBusNamesIsolateTraffic() throws Exception {
        VirtualCanTransportInstance instanceA = createInstance("bus-a");
        VirtualCanTransportInstance instanceB = createInstance("bus-b");

        byte[] frame = {0x0A, 0x0B};
        instanceA.write(frame);

        // Instance on a different bus should not receive the frame
        assertEquals(0, instanceB.getNumBytesAvailable());
    }

    @Test
    void testDataListenerNotification() throws Exception {
        VirtualCanTransportInstance sender = createInstance("bus1");
        VirtualCanTransportInstance receiver = createInstance("bus1");

        AtomicBoolean listenerNotified = new AtomicBoolean(false);
        receiver.registerDataListener(() -> listenerNotified.set(true));

        sender.write(new byte[]{0x01, 0x02});

        assertTrue(listenerNotified.get());
    }

    @Test
    void testDataListenerNotCalledAfterRemoval() throws Exception {
        VirtualCanTransportInstance sender = createInstance("bus1");
        VirtualCanTransportInstance receiver = createInstance("bus1");

        AtomicBoolean listenerNotified = new AtomicBoolean(false);
        receiver.registerDataListener(() -> listenerNotified.set(true));
        receiver.removeDataListener();

        sender.write(new byte[]{0x01, 0x02});

        assertFalse(listenerNotified.get());
    }

    @Test
    void testNoErrorSendingToBusWithNoOtherParticipants() throws Exception {
        VirtualCanTransportInstance alone = createInstance("lonely-bus");

        // Writing when no other instances are on the bus should not throw
        assertDoesNotThrow(() -> alone.write(new byte[]{0x01, 0x02, 0x03}));

        // The sender should not receive its own data
        assertEquals(0, alone.getNumBytesAvailable());
    }

    @Test
    void testWriteNullIsNoOp() throws Exception {
        VirtualCanTransportInstance instance = createInstance("bus1");
        // Writing null bytes should silently return
        assertDoesNotThrow(() -> instance.write(null));
    }

    @Test
    void testWriteEmptyArrayIsNoOp() throws Exception {
        VirtualCanTransportInstance instance = createInstance("bus1");
        // Writing empty bytes should silently return
        assertDoesNotThrow(() -> instance.write(new byte[0]));
    }

    @Test
    void testReadZeroBytesReturnsEmptyArray() throws Exception {
        VirtualCanTransportInstance instance = createInstance("bus1");
        byte[] result = instance.read(0);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekZeroBytesReturnsEmptyArray() throws Exception {
        VirtualCanTransportInstance instance = createInstance("bus1");
        byte[] result = instance.peekReadableBytes(0);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekDoesNotConsumeData() throws Exception {
        VirtualCanTransportInstance sender = createInstance("bus1");
        VirtualCanTransportInstance receiver = createInstance("bus1");

        byte[] frame = {0x0A, 0x0B, 0x0C};
        sender.write(frame);

        // Peek should not consume the data
        byte[] peeked = receiver.peekReadableBytes(3);
        assertArrayEquals(frame, peeked);
        assertEquals(3, receiver.getNumBytesAvailable());

        // Read should consume the data
        byte[] read = receiver.read(3);
        assertArrayEquals(frame, read);
        assertEquals(0, receiver.getNumBytesAvailable());
    }

    @Test
    void testMultipleFramesAccumulate() throws Exception {
        VirtualCanTransportInstance sender = createInstance("bus1");
        VirtualCanTransportInstance receiver = createInstance("bus1");

        sender.write(new byte[]{0x01, 0x02});
        sender.write(new byte[]{0x03, 0x04});

        assertEquals(4, receiver.getNumBytesAvailable());
        byte[] allData = receiver.read(4);
        assertArrayEquals(new byte[]{0x01, 0x02, 0x03, 0x04}, allData);
    }

    @Test
    void testGetConfigurationReturnsProvidedConfig() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.busName = "my-bus";
        VirtualCanTransportInstance instance = new VirtualCanTransportInstance(config, auditLog);
        openInstances.add(instance);

        assertSame(config, instance.getConfiguration());
        assertEquals("my-bus", instance.getConfiguration().busName);
    }

    @Test
    void testClosedInstanceDoesNotReceiveFrames() throws Exception {
        VirtualCanTransportInstance sender = createInstance("bus1");
        VirtualCanTransportInstance receiver = createInstance("bus1");

        receiver.close();

        // Writing after receiver is closed should not throw
        assertDoesNotThrow(() -> sender.write(new byte[]{0x01}));
    }

    @Test
    void testThreeInstanceBroadcast() throws Exception {
        VirtualCanTransportInstance a = createInstance("bus1");
        VirtualCanTransportInstance b = createInstance("bus1");
        VirtualCanTransportInstance c = createInstance("bus1");

        byte[] frame = {(byte) 0xCA, (byte) 0xFE};
        a.write(frame);

        // a should not receive, b and c should
        assertEquals(0, a.getNumBytesAvailable());
        assertEquals(2, b.getNumBytesAvailable());
        assertEquals(2, c.getNumBytesAvailable());

        assertArrayEquals(frame, b.read(2));
        assertArrayEquals(frame, c.read(2));
    }
}
