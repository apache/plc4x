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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.channels.DatagramChannel;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Manages shared UDP sockets across multiple transport instances.
 * Allows multiple logical connections to share a single UDP socket,
 * which is useful for protocols that multiplex connections over one port.
 */
public class SharedUdpSocketManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(SharedUdpSocketManager.class);

    // Map of socket key -> shared socket wrapper
    private final Map<SocketKey, SharedSocket> sharedSockets = new ConcurrentHashMap<>();

    public SharedUdpSocketManager() {
        // Private constructor for singleton
    }

    /**
     * Acquires a shared socket for the given local address and port.
     * If a socket already exists for this address/port, returns the existing one.
     * Otherwise, creates a new one.
     */
    public SharedSocket acquireSocket(String localAddress, int localPort) throws IOException {
        SocketKey key = new SocketKey(localAddress, localPort);

        return sharedSockets.compute(key, (k, existing) -> {
            if (existing != null) {
                existing.incrementRefCount();
                LOGGER.debug("Reusing shared UDP socket on {}:{} (refCount={})",
                    localAddress, localPort, existing.getRefCount());
                return existing;
            } else {
                try {
                    DatagramChannel channel = DatagramChannel.open();

                    // Bind to specified address/port
                    SocketAddress bindAddress;
                    if (localAddress != null && !localAddress.isEmpty()) {
                        bindAddress = new InetSocketAddress(localAddress, localPort);
                    } else {
                        bindAddress = new InetSocketAddress(localPort);
                    }
                    channel.bind(bindAddress);

                    LOGGER.info("Created new shared UDP socket on {}:{}", localAddress, localPort);
                    return new SharedSocket(channel, k);

                } catch (IOException e) {
                    throw new RuntimeException("Failed to create UDP socket", e);
                }
            }
        });
    }

    /**
     * Releases a shared socket reference.
     * If this is the last reference, closes the socket.
     */
    public void releaseSocket(SharedSocket socket) {
        SocketKey key = socket.getKey();

        sharedSockets.compute(key, (k, existing) -> {
            if (existing == null) {
                LOGGER.warn("Attempted to release non-existent socket for {}", key);
                return null;
            }

            int refCount = existing.decrementRefCount();
            LOGGER.debug("Released shared UDP socket on {} (refCount={})", key, refCount);

            if (refCount <= 0) {
                try {
                    existing.getChannel().close();
                    LOGGER.info("Closed shared UDP socket on {}", key);
                } catch (IOException e) {
                    LOGGER.error("Failed to close UDP socket", e);
                }
                return null; // Remove from map
            }

            return existing;
        });
    }

    /**
     * Wrapper around a DatagramChannel with reference counting.
     */
    public static class SharedSocket {
        private final DatagramChannel channel;
        private final SocketKey key;
        private final AtomicInteger refCount = new AtomicInteger(1);

        public SharedSocket(DatagramChannel channel, SocketKey key) {
            this.channel = channel;
            this.key = key;
        }

        public DatagramChannel getChannel() {
            return channel;
        }

        public SocketKey getKey() {
            return key;
        }

        public int getRefCount() {
            return refCount.get();
        }

        void incrementRefCount() {
            refCount.incrementAndGet();
        }

        int decrementRefCount() {
            return refCount.decrementAndGet();
        }
    }

    /**
     * Key for identifying unique socket bindings.
     */
    public static class SocketKey {
        private final String address;
        private final int port;

        public SocketKey(String address, int port) {
            this.address = address == null ? "" : address;
            this.port = port;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            SocketKey socketKey = (SocketKey) o;
            return port == socketKey.port && address.equals(socketKey.address);
        }

        @Override
        public int hashCode() {
            return 31 * address.hashCode() + port;
        }

        @Override
        public String toString() {
            return address.isEmpty() ? "0.0.0.0:" + port : address + ":" + port;
        }
    }
}
