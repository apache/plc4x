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
package org.apache.plc4x.java.transport.udp;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.net.InetSocketAddress;

import static org.junit.jupiter.api.Assertions.*;

class SharedUdpSocketManagerTest {

    private SharedUdpSocketManager manager;

    @BeforeEach
    void setUp() {
        manager = new SharedUdpSocketManager();
    }

    @Test
    void testAcquireSocket_createsNewSocket() throws IOException {
        SharedUdpSocketManager.SharedSocket socket = manager.acquireSocket("localhost", 0);

        assertNotNull(socket);
        assertNotNull(socket.getChannel());
        assertTrue(socket.getChannel().isOpen());
        assertEquals(1, socket.getRefCount());

        // Cleanup
        manager.releaseSocket(socket);
    }

    @Test
    void testAcquireSocket_reusesSameSocket() throws IOException {
        // Acquire first socket on specific port
        int port = findAvailablePort();
        SharedUdpSocketManager.SharedSocket socket1 = manager.acquireSocket("localhost", port);

        // Acquire second socket on same address/port
        SharedUdpSocketManager.SharedSocket socket2 = manager.acquireSocket("localhost", port);

        // Should be the same socket
        assertSame(socket1, socket2);
        assertSame(socket1.getChannel(), socket2.getChannel());
        assertEquals(2, socket1.getRefCount());

        // Cleanup
        manager.releaseSocket(socket1);
        manager.releaseSocket(socket2);
    }

    @Test
    void testAcquireSocket_differentPorts_createsDifferentSockets() throws IOException {
        int port1 = findAvailablePort();
        int port2 = findAvailablePort();

        SharedUdpSocketManager.SharedSocket socket1 = manager.acquireSocket("localhost", port1);
        SharedUdpSocketManager.SharedSocket socket2 = manager.acquireSocket("localhost", port2);

        assertNotSame(socket1, socket2);
        assertNotSame(socket1.getChannel(), socket2.getChannel());

        // Cleanup
        manager.releaseSocket(socket1);
        manager.releaseSocket(socket2);
    }

    @Test
    void testReleaseSocket_decrementsRefCount() throws IOException {
        int port = findAvailablePort();
        SharedUdpSocketManager.SharedSocket socket = manager.acquireSocket("localhost", port);

        assertEquals(1, socket.getRefCount());

        // Acquire again
        manager.acquireSocket("localhost", port);
        assertEquals(2, socket.getRefCount());

        // Release once
        manager.releaseSocket(socket);
        assertEquals(1, socket.getRefCount());
        assertTrue(socket.getChannel().isOpen()); // Should still be open

        // Release again
        manager.releaseSocket(socket);
        assertEquals(0, socket.getRefCount());
        assertFalse(socket.getChannel().isOpen()); // Should be closed now
    }

    @Test
    void testReleaseSocket_closesSocketWhenRefCountZero() throws IOException {
        int port = findAvailablePort();
        SharedUdpSocketManager.SharedSocket socket = manager.acquireSocket("localhost", port);

        assertTrue(socket.getChannel().isOpen());
        assertEquals(1, socket.getRefCount());

        manager.releaseSocket(socket);

        assertFalse(socket.getChannel().isOpen());
        assertEquals(0, socket.getRefCount());
    }

    @Test
    void testReleaseSocket_nonExistent_doesNotThrow() {
        int port = 12345;
        SharedUdpSocketManager.SharedSocket socket = new SharedUdpSocketManager.SharedSocket(
            null, new SharedUdpSocketManager.SocketKey("localhost", port)
        );

        // Should not throw
        assertDoesNotThrow(() -> manager.releaseSocket(socket));
    }

    @Test
    void testMultipleAcquireRelease_cycles() throws IOException {
        int port = findAvailablePort();

        // First cycle
        SharedUdpSocketManager.SharedSocket socket1 = manager.acquireSocket("localhost", port);
        assertEquals(1, socket1.getRefCount());
        manager.releaseSocket(socket1);
        assertFalse(socket1.getChannel().isOpen());

        // Second cycle - should create new socket
        SharedUdpSocketManager.SharedSocket socket2 = manager.acquireSocket("localhost", port);
        assertEquals(1, socket2.getRefCount());
        assertNotSame(socket1, socket2); // Different instance
        assertTrue(socket2.getChannel().isOpen());

        // Cleanup
        manager.releaseSocket(socket2);
    }

    @Test
    void testConcurrentAcquireRelease() throws Exception {
        int port = findAvailablePort();

        Thread[] threads = new Thread[10];
        for (int i = 0; i < threads.length; i++) {
            threads[i] = new Thread(() -> {
                try {
                    for (int j = 0; j < 100; j++) {
                        SharedUdpSocketManager.SharedSocket socket = manager.acquireSocket("localhost", port);
                        // Simulate some work
                        Thread.sleep(1);
                        manager.releaseSocket(socket);
                    }
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });
            threads[i].start();
        }

        // Wait for all threads
        for (Thread thread : threads) {
            thread.join();
        }

        // All sockets should be released
        // (Note: This is difficult to assert directly without exposing internal state)
    }

    @Test
    void testSocketKey_equality() {
        SharedUdpSocketManager.SocketKey key1 = new SharedUdpSocketManager.SocketKey("localhost", 8080);
        SharedUdpSocketManager.SocketKey key2 = new SharedUdpSocketManager.SocketKey("localhost", 8080);
        SharedUdpSocketManager.SocketKey key3 = new SharedUdpSocketManager.SocketKey("localhost", 8081);
        SharedUdpSocketManager.SocketKey key4 = new SharedUdpSocketManager.SocketKey("127.0.0.1", 8080);

        assertEquals(key1, key2);
        assertNotEquals(key1, key3);
        assertNotEquals(key1, key4);
    }

    @Test
    void testSocketKey_hashCode() {
        SharedUdpSocketManager.SocketKey key1 = new SharedUdpSocketManager.SocketKey("localhost", 8080);
        SharedUdpSocketManager.SocketKey key2 = new SharedUdpSocketManager.SocketKey("localhost", 8080);

        assertEquals(key1.hashCode(), key2.hashCode());
    }

    private int findAvailablePort() throws IOException {
        try (var channel = java.nio.channels.DatagramChannel.open()) {
            channel.bind(new InetSocketAddress("localhost", 0));
            return ((InetSocketAddress) channel.getLocalAddress()).getPort();
        }
    }
}
