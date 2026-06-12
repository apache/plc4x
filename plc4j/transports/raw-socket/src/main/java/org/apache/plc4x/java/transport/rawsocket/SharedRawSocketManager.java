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

import org.pcap4j.core.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Manages shared raw socket (pcap) handles across multiple transport instances.
 * Allows multiple logical connections to share a single network interface capture,
 * which is essential for protocols that multiplex over one EtherType.
 */
public class SharedRawSocketManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(SharedRawSocketManager.class);

    // Map of handle key -> shared handle wrapper
    private final Map<HandleKey, SharedHandle> sharedHandles = new ConcurrentHashMap<>();

    public SharedRawSocketManager() {
        // Private constructor for singleton
    }

    /**
     * Acquires a shared pcap handle for the given interface and configuration.
     * If a handle is already open for this interface/protocol, returns the existing one.
     * Otherwise creates a new one.
     */
    public SharedHandle acquireHandle(PcapNetworkInterface nif, PcapConfig config) throws PcapNativeException {
        HandleKey key = new HandleKey(nif.getName(), config.protocolId);

        return sharedHandles.compute(key, (k, existing) -> {
            if (existing != null) {
                existing.incrementRefCount();
                LOGGER.debug("Reusing shared pcap handle on {} for protocol 0x{} (refCount={})",
                    nif.getName(), String.format("%04X", config.protocolId), existing.getRefCount());
                return existing;
            } else {
                try {
                    PcapHandle handle = nif.openLive(
                        config.snapshotLength,
                        config.promiscuousMode ? PcapNetworkInterface.PromiscuousMode.PROMISCUOUS
                            : PcapNetworkInterface.PromiscuousMode.NONPROMISCUOUS,
                        config.captureTimeout
                    );

                    // Set buffer size
                    if (config.bufferSize > 0) {
                        // TODO: Fix this.
                        //handle.setBufferSize(config.bufferSize);
                    }

                    // Apply BPF filter if specified
                    if (config.bpfFilter != null && !config.bpfFilter.isEmpty()) {
                        handle.setFilter(config.bpfFilter, BpfProgram.BpfCompileMode.OPTIMIZE);
                        LOGGER.debug("Applied BPF filter: {}", config.bpfFilter);
                    }

                    LOGGER.info("Opened shared pcap handle on {} for protocol 0x{}",
                        nif.getName(), String.format("%04X", config.protocolId));

                    return new SharedHandle(handle, k);

                } catch (PcapNativeException | NotOpenException e) {
                    throw new RuntimeException("Failed to open pcap handle on " + nif.getName(), e);
                }
            }
        });
    }

    /**
     * Releases a shared pcap handle reference.
     * If this is the last reference, closes the handle.
     */
    public void releaseHandle(SharedHandle handle) {
        HandleKey key = handle.getKey();

        sharedHandles.compute(key, (k, existing) -> {
            if (existing == null) {
                LOGGER.warn("Attempted to release non-existent handle: {}", key);
                return null;
            }

            int refCount = existing.decrementRefCount();
            LOGGER.debug("Released shared pcap handle {} (refCount={})", key, refCount);

            if (refCount <= 0) {
                existing.getHandle().close();
                LOGGER.info("Closed shared pcap handle {}", key);
                return null; // Remove from map
            }

            return existing;
        });
    }

    /**
     * Configuration for opening a pcap handle.
     */
    public static class PcapConfig {
        public final int protocolId;
        public final int snapshotLength;
        public final boolean promiscuousMode;
        public final int captureTimeout;
        public final int bufferSize;
        public final String bpfFilter;

        public PcapConfig(int protocolId, int snapshotLength, boolean promiscuousMode,
                         int captureTimeout, int bufferSize, String bpfFilter) {
            this.protocolId = protocolId;
            this.snapshotLength = snapshotLength;
            this.promiscuousMode = promiscuousMode;
            this.captureTimeout = captureTimeout;
            this.bufferSize = bufferSize;
            this.bpfFilter = bpfFilter;
        }
    }

    /**
     * Wrapper around a PcapHandle with reference counting.
     */
    public static class SharedHandle {
        private final PcapHandle handle;
        private final HandleKey key;
        private final AtomicInteger refCount = new AtomicInteger(1);

        private SharedHandle(PcapHandle handle, HandleKey key) {
            this.handle = handle;
            this.key = key;
        }

        public PcapHandle getHandle() {
            return handle;
        }

        public HandleKey getKey() {
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
     * Key for identifying unique pcap handles.
     */
    public static class HandleKey {
        private final String interfaceName;
        private final int protocolId;

        protected HandleKey(String interfaceName, int protocolId) {
            this.interfaceName = interfaceName;
            this.protocolId = protocolId;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            HandleKey handleKey = (HandleKey) o;
            return protocolId == handleKey.protocolId && interfaceName.equals(handleKey.interfaceName);
        }

        @Override
        public int hashCode() {
            return 31 * interfaceName.hashCode() + protocolId;
        }

        @Override
        public String toString() {
            return interfaceName + ":0x" + String.format("%04X", protocolId);
        }
    }
}
