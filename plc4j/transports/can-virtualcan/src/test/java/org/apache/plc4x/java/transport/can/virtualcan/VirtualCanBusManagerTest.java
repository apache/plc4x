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

import org.apache.plc4x.java.transport.can.virtualcan.config.VirtualCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for {@link VirtualCanBusManager}: bus creation, instance management,
 * frame broadcasting, and cleanup.
 */
class VirtualCanBusManagerTest {

    private AuditLog auditLog;

    @BeforeEach
    void setUp() {
        VirtualCanBusManager.reset();
        auditLog = AuditLog.builder().build();
    }

    @AfterEach
    void tearDown() {
        VirtualCanBusManager.reset();
    }

    /**
     * Helper to create a VirtualCanTransportInstance with a given bus name.
     */
    private VirtualCanTransportInstance createInstance(String busName) {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.busName = busName;
        return new VirtualCanTransportInstance(config, auditLog);
    }

    @Test
    void testConnectCreatesBus() {
        VirtualCanTransportInstance instance = createInstance("bus1");
        assertTrue(instance.isOpen());
        instance.close();
    }

    @Test
    void testMultipleInstancesOnSameBus() {
        VirtualCanTransportInstance instance1 = createInstance("bus1");
        VirtualCanTransportInstance instance2 = createInstance("bus1");

        assertTrue(instance1.isOpen());
        assertTrue(instance2.isOpen());

        instance1.close();
        instance2.close();
    }

    @Test
    void testDisconnectRemovesInstance() {
        VirtualCanTransportInstance instance1 = createInstance("bus1");
        VirtualCanTransportInstance instance2 = createInstance("bus1");

        instance1.close();
        assertFalse(instance1.isOpen());
        // Second instance should still be open
        assertTrue(instance2.isOpen());

        instance2.close();
    }

    @Test
    void testDisconnectLastRemovesBus() {
        VirtualCanTransportInstance instance = createInstance("bus1");
        instance.close();

        // After the last instance disconnects, the bus should be gone.
        // We verify this indirectly: broadcast to a removed bus should be a no-op.
        // If the bus still existed, this would fail or cause unexpected behavior.
        VirtualCanBusManager.broadcast("bus1", instance, new byte[]{1, 2, 3});
        // No exception means it handled the missing bus gracefully
    }

    @Test
    void testBroadcastDeliversToAllExceptSender() throws Exception {
        VirtualCanTransportInstance sender = createInstance("bus1");
        VirtualCanTransportInstance receiver1 = createInstance("bus1");
        VirtualCanTransportInstance receiver2 = createInstance("bus1");

        byte[] frame = {0x01, 0x02, 0x03, 0x04};
        VirtualCanBusManager.broadcast("bus1", sender, frame);

        // Sender should not have received the frame
        assertEquals(0, sender.getNumBytesAvailable());

        // Both receivers should have the frame
        assertEquals(4, receiver1.getNumBytesAvailable());
        assertArrayEquals(frame, receiver1.read(4));

        assertEquals(4, receiver2.getNumBytesAvailable());
        assertArrayEquals(frame, receiver2.read(4));

        sender.close();
        receiver1.close();
        receiver2.close();
    }

    @Test
    void testBroadcastToEmptyBusNoError() {
        // Broadcasting to a bus that does not exist should not throw
        assertDoesNotThrow(() ->
                VirtualCanBusManager.broadcast("nonexistent", null, new byte[]{1, 2, 3}));
    }

    @Test
    void testBroadcastToSingleInstanceBus() throws Exception {
        // Only one instance on the bus — broadcast should be a no-op
        VirtualCanTransportInstance alone = createInstance("lonely-bus");

        VirtualCanBusManager.broadcast("lonely-bus", alone, new byte[]{1, 2, 3});
        assertEquals(0, alone.getNumBytesAvailable());

        alone.close();
    }

    @Test
    void testResetClearsAllBuses() throws Exception {
        VirtualCanTransportInstance instance1 = createInstance("bus1");
        VirtualCanTransportInstance instance2 = createInstance("bus2");

        VirtualCanBusManager.reset();

        // After reset, broadcasting should be a no-op (buses are gone)
        VirtualCanBusManager.broadcast("bus1", instance1, new byte[]{1});
        VirtualCanBusManager.broadcast("bus2", instance2, new byte[]{2});

        assertEquals(0, instance1.getNumBytesAvailable());
        assertEquals(0, instance2.getNumBytesAvailable());

        // Clean up the instances (they're still "open" from their perspective)
        instance1.close();
        instance2.close();
    }
}
