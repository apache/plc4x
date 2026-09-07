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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tel.schich.javacan.CanChannels;
import tel.schich.javacan.RawCanChannel;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Manages shared CAN socket handles across multiple transport instances.
 * <p>
 * When multiple transport instances need to communicate on the same CAN interface,
 * sharing a single socket is more efficient than opening separate sockets. This manager
 * uses reference counting to ensure the shared socket stays open as long as at least one
 * instance needs it, and is closed when the last instance releases it.
 * <p>
 * Thread-safe: uses {@link ConcurrentHashMap} for the handle registry and
 * {@link AtomicInteger} for reference counting.
 *
 * @see SharedCanHandle
 */
public class SharedCanManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(SharedCanManager.class);

    private final Map<String, SharedCanHandle> sharedHandles = new ConcurrentHashMap<>();

    /**
     * Acquires a shared CAN socket for the given interface name.
     * <p>
     * If a socket is already open for this interface, returns the existing handle
     * with an incremented reference count. Otherwise, opens a new CAN socket.
     *
     * @param interfaceName the CAN interface (e.g., "can0", "vcan0")
     * @return the shared handle wrapping the CAN socket
     * @throws IOException if the CAN socket cannot be opened
     */
    public SharedCanHandle acquireHandle(String interfaceName) throws IOException {
        return sharedHandles.compute(interfaceName, (key, existing) -> {
            if (existing != null) {
                existing.incrementRefCount();
                LOGGER.debug("Reusing shared CAN handle on {} (refCount={})",
                        interfaceName, existing.getRefCount());
                return existing;
            } else {
                try {
                    RawCanChannel channel = CanChannels.newRawChannel(interfaceName);
                    LOGGER.info("Opened shared CAN socket on {}", interfaceName);
                    return new SharedCanHandle(channel, interfaceName);
                } catch (IOException e) {
                    throw new RuntimeException("Failed to open CAN socket on " + interfaceName, e);
                }
            }
        });
    }

    /**
     * Releases a shared CAN socket reference.
     * <p>
     * Decrements the reference count. If this was the last reference,
     * closes the underlying CAN socket and removes it from the registry.
     *
     * @param handle the handle to release
     */
    public void releaseHandle(SharedCanHandle handle) {
        String interfaceName = handle.getInterfaceName();

        sharedHandles.compute(interfaceName, (key, existing) -> {
            if (existing == null) {
                LOGGER.warn("Attempted to release non-existent handle for: {}", interfaceName);
                return null;
            }

            int refCount = existing.decrementRefCount();
            LOGGER.debug("Released shared CAN handle {} (refCount={})", interfaceName, refCount);

            if (refCount <= 0) {
                try {
                    existing.getChannel().close();
                    LOGGER.info("Closed shared CAN socket on {}", interfaceName);
                } catch (IOException e) {
                    LOGGER.warn("Error closing shared CAN socket on {}: {}", interfaceName, e.getMessage());
                }
                return null; // Remove from map
            }

            return existing;
        });
    }

    /**
     * Reference-counted wrapper around a shared CAN socket channel.
     */
    public static class SharedCanHandle {
        private final RawCanChannel channel;
        private final String interfaceName;
        private final AtomicInteger refCount = new AtomicInteger(1);

        SharedCanHandle(RawCanChannel channel, String interfaceName) {
            this.channel = channel;
            this.interfaceName = interfaceName;
        }

        /**
         * Returns the underlying CAN channel.
         *
         * @return the raw CAN channel
         */
        public RawCanChannel getChannel() {
            return channel;
        }

        /**
         * Returns the interface name this handle is bound to.
         *
         * @return the CAN interface name
         */
        public String getInterfaceName() {
            return interfaceName;
        }

        /**
         * Returns the current reference count.
         *
         * @return the number of active references
         */
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
}
