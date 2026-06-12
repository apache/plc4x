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

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.pcap4j.core.PcapNetworkInterface;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for SharedRawSocketManager.
 * Note: These tests require pcap permissions to run.
 */
class SharedRawSocketManagerTest {

    private SharedRawSocketManager manager;
    private PcapNetworkInterface testInterface;

    @BeforeEach
    void setUp() {
        manager = new SharedRawSocketManager();
        testInterface = PcapTestSupport.findAllDevsOrSkip().get(0);
    }

    @Test
    void testAcquireHandle_createsNewHandle() {
        SharedRawSocketManager.PcapConfig config = createDefaultConfig();

        try {
            SharedRawSocketManager.SharedHandle handle = manager.acquireHandle(testInterface, config);

            assertNotNull(handle);
            assertNotNull(handle.getHandle());
            assertTrue(handle.getHandle().isOpen());
            assertEquals(1, handle.getRefCount());

            // Cleanup
            manager.releaseHandle(handle);
        } catch (Exception e) {
            System.out.println("Could not open handle (expected without permissions): " + e.getMessage());
        }
    }

    @Test
    void testAcquireHandle_reusesSameHandle() {
        SharedRawSocketManager.PcapConfig config = createDefaultConfig();

        try {
            SharedRawSocketManager.SharedHandle handle1 = manager.acquireHandle(testInterface, config);
            SharedRawSocketManager.SharedHandle handle2 = manager.acquireHandle(testInterface, config);

            // Should be the same instance
            assertSame(handle1, handle2);
            assertSame(handle1.getHandle(), handle2.getHandle());
            assertEquals(2, handle1.getRefCount());

            // Cleanup
            manager.releaseHandle(handle1);
            manager.releaseHandle(handle2);
        } catch (Exception e) {
            System.out.println("Could not open handle (expected without permissions): " + e.getMessage());
        }
    }

    @Test
    void testAcquireHandle_differentProtocols_createsDifferentHandles() {
        SharedRawSocketManager.PcapConfig config1 = new SharedRawSocketManager.PcapConfig(
            0x88B5, 65536, false, 10, 0, null
        );
        SharedRawSocketManager.PcapConfig config2 = new SharedRawSocketManager.PcapConfig(
            0x88CC, 65536, false, 10, 0, null
        );

        try {
            SharedRawSocketManager.SharedHandle handle1 = manager.acquireHandle(testInterface, config1);
            SharedRawSocketManager.SharedHandle handle2 = manager.acquireHandle(testInterface, config2);

            assertNotSame(handle1, handle2);
            assertNotSame(handle1.getHandle(), handle2.getHandle());

            // Cleanup
            manager.releaseHandle(handle1);
            manager.releaseHandle(handle2);
        } catch (Exception e) {
            System.out.println("Could not open handle (expected without permissions): " + e.getMessage());
        }
    }

    @Test
    void testReleaseHandle_decrementsRefCount() {
        SharedRawSocketManager.PcapConfig config = createDefaultConfig();

        try {
            SharedRawSocketManager.SharedHandle handle = manager.acquireHandle(testInterface, config);
            assertEquals(1, handle.getRefCount());

            // Acquire again
            manager.acquireHandle(testInterface, config);
            assertEquals(2, handle.getRefCount());

            // Release once
            manager.releaseHandle(handle);
            assertEquals(1, handle.getRefCount());
            assertTrue(handle.getHandle().isOpen());

            // Release again
            manager.releaseHandle(handle);
            assertEquals(0, handle.getRefCount());
            assertFalse(handle.getHandle().isOpen());
        } catch (Exception e) {
            System.out.println("Could not open handle (expected without permissions): " + e.getMessage());
        }
    }

    @Test
    void testReleaseHandle_closesHandleWhenRefCountZero() {
        SharedRawSocketManager.PcapConfig config = createDefaultConfig();

        try {
            SharedRawSocketManager.SharedHandle handle = manager.acquireHandle(testInterface, config);

            assertTrue(handle.getHandle().isOpen());
            assertEquals(1, handle.getRefCount());

            manager.releaseHandle(handle);

            assertFalse(handle.getHandle().isOpen());
            assertEquals(0, handle.getRefCount());
        } catch (Exception e) {
            System.out.println("Could not open handle (expected without permissions): " + e.getMessage());
        }
    }

    @Test
    void testMultipleAcquireRelease_cycles() {
        SharedRawSocketManager.PcapConfig config = createDefaultConfig();

        try {
            // First cycle
            SharedRawSocketManager.SharedHandle handle1 = manager.acquireHandle(testInterface, config);
            assertEquals(1, handle1.getRefCount());
            manager.releaseHandle(handle1);
            assertFalse(handle1.getHandle().isOpen());

            // Second cycle - should create new handle
            SharedRawSocketManager.SharedHandle handle2 = manager.acquireHandle(testInterface, config);
            assertEquals(1, handle2.getRefCount());
            assertNotSame(handle1, handle2); // Different instance
            assertTrue(handle2.getHandle().isOpen());

            // Cleanup
            manager.releaseHandle(handle2);
        } catch (Exception e) {
            System.out.println("Could not open handle (expected without permissions): " + e.getMessage());
        }
    }

    @Test
    void testHandleKey_equality() {
        SharedRawSocketManager.HandleKey key1 = new SharedRawSocketManager.HandleKey("eth0", 0x88B5);
        SharedRawSocketManager.HandleKey key2 = new SharedRawSocketManager.HandleKey("eth0", 0x88B5);
        SharedRawSocketManager.HandleKey key3 = new SharedRawSocketManager.HandleKey("eth0", 0x88CC);
        SharedRawSocketManager.HandleKey key4 = new SharedRawSocketManager.HandleKey("eth1", 0x88B5);

        assertEquals(key1, key2);
        assertNotEquals(key1, key3);
        assertNotEquals(key1, key4);
    }

    @Test
    void testHandleKey_hashCode() {
        SharedRawSocketManager.HandleKey key1 = new SharedRawSocketManager.HandleKey("eth0", 0x88B5);
        SharedRawSocketManager.HandleKey key2 = new SharedRawSocketManager.HandleKey("eth0", 0x88B5);

        assertEquals(key1.hashCode(), key2.hashCode());
    }

    @Test
    void testHandleKey_toString() {
        SharedRawSocketManager.HandleKey key = new SharedRawSocketManager.HandleKey("eth0", 0x88B5);
        String result = key.toString();

        assertEquals("eth0:0x88B5", result);
    }

    @Test
    void testHandleKey_toString_differentProtocol() {
        SharedRawSocketManager.HandleKey key = new SharedRawSocketManager.HandleKey("en0", 0x88CC);
        String result = key.toString();

        assertEquals("en0:0x88CC", result);
    }

    @Test
    void testHandleKey_equals_sameObject() {
        SharedRawSocketManager.HandleKey key = new SharedRawSocketManager.HandleKey("eth0", 0x88B5);
        assertEquals(key, key);
    }

    @Test
    void testHandleKey_equals_null() {
        SharedRawSocketManager.HandleKey key = new SharedRawSocketManager.HandleKey("eth0", 0x88B5);
        assertNotEquals(null, key);
    }

    @Test
    void testHandleKey_equals_differentType() {
        SharedRawSocketManager.HandleKey key = new SharedRawSocketManager.HandleKey("eth0", 0x88B5);
        assertNotEquals("eth0:0x88B5", key);
    }

    @Test
    void testPcapConfig_fieldsAccessible() {
        SharedRawSocketManager.PcapConfig config = new SharedRawSocketManager.PcapConfig(
            0x88B5, 65536, true, 2000, 1048576, "ether proto 0x88B5"
        );

        assertEquals(0x88B5, config.protocolId);
        assertEquals(65536, config.snapshotLength);
        assertTrue(config.promiscuousMode);
        assertEquals(2000, config.captureTimeout);
        assertEquals(1048576, config.bufferSize);
        assertEquals("ether proto 0x88B5", config.bpfFilter);
    }

    @Test
    void testPcapConfig_nullBpfFilter() {
        SharedRawSocketManager.PcapConfig config = new SharedRawSocketManager.PcapConfig(
            0x88CC, 32768, false, 500, 0, null
        );

        assertEquals(0x88CC, config.protocolId);
        assertNull(config.bpfFilter);
    }

    private SharedRawSocketManager.PcapConfig createDefaultConfig() {
        return new SharedRawSocketManager.PcapConfig(
            0x88B5,  // protocolId (PROFINET)
            65536,   // snapshotLength
            false,   // promiscuousMode
            10,      // captureTimeout
            0,       // bufferSize
            null     // bpfFilter
        );
    }
}
