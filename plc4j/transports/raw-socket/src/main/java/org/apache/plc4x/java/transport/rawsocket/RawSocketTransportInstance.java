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

import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.transport.rawsocket.config.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.pcap4j.core.*;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.UnknownPacket;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.MacAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;

import static java.lang.Math.max;

/**
 * Java 21 efficient Raw Socket transport implementation using pcap4j with async support.
 * Implements AsyncTransportInstance for event-driven I/O, critical for timing-sensitive
 * protocols like PROFINET, EtherNet/IP, and EtherCAT.
 * Supports both dedicated and shared interface modes for optimal resource usage.
 */
public class RawSocketTransportInstance extends BaseTransportInstance<RawSocketTransportConfiguration> implements AsyncTransportInstance<RawSocketTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(RawSocketTransportInstance.class);
    private static final int DEFAULT_BUFFER_SIZE = 8192;
    /** Minimum capture timeout to ensure the capture thread can check the open flag regularly */
    private static final int MIN_CAPTURE_TIMEOUT_MS = 100;

    private final SharedRawSocketManager sharedRawSocketManager;
    private final PcapHandle handle;
    private final PcapNetworkInterface networkInterface;
    private final MacAddress localMac;
    private final MacAddress remoteMac;
    private final SharedRawSocketManager.SharedHandle sharedHandle; // null if not shared
    private final BlockingQueue<byte[]> receiveQueue;
    private final RingBuffer ringBuffer;
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();
    private volatile boolean open = true;
    private volatile Thread captureThread;

    // Async support
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;

    public RawSocketTransportInstance(SharedRawSocketManager sharedRawSocketManager, RawSocketTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        LOGGER.debug("RawSocketTransportInstance: Pre Java 21 version");
        this.sharedRawSocketManager = sharedRawSocketManager;
        this.receiveQueue = new LinkedBlockingQueue<>();
        this.ringBuffer = new RingBuffer(DEFAULT_BUFFER_SIZE);

        try {
            // Find network interface
            PcapNetworkInterface nif = findNetworkInterface(configuration.interfaceName);
            LOGGER.debug("Using network interface: {} ({})", nif.getName(), nif.getDescription());

            // Parse MAC addresses
            MacAddress tempLocalMac = configuration.localAddress != null && !configuration.localAddress.isEmpty()
                ? MacAddress.getByName(configuration.localAddress)
                : getInterfaceMacAddress(nif);

            MacAddress tempRemoteMac = MacAddress.getByName(configuration.remoteAddress);

            LOGGER.debug("Local MAC: {}, Remote MAC: {}", tempLocalMac, tempRemoteMac);

            // Create or acquire pcap handle
            PcapHandle tempHandle;
            SharedRawSocketManager.SharedHandle tempSharedHandle = null;

            // Ensure a minimum capture timeout so the capture thread can check the open flag
            int effectiveCaptureTimeout = max(configuration.captureTimeout, MIN_CAPTURE_TIMEOUT_MS);

            if (configuration.reuseInterface) {
                // Use a shared handle manager
                SharedRawSocketManager.PcapConfig pcapConfig = new SharedRawSocketManager.PcapConfig(
                    configuration.protocolId,
                    configuration.snapshotLength,
                    configuration.promiscuousMode,
                    effectiveCaptureTimeout,
                    configuration.bufferSize,
                    buildBpfFilter(configuration, tempRemoteMac, tempLocalMac)
                );

                tempSharedHandle = sharedRawSocketManager.acquireHandle(nif, pcapConfig);
                tempHandle = tempSharedHandle.getHandle();

                LOGGER.debug("Using shared pcap handle on {}", nif.getName());

            } else {
                // Create a dedicated handle
                tempHandle = nif.openLive(
                    configuration.snapshotLength,
                    configuration.promiscuousMode ? PcapNetworkInterface.PromiscuousMode.PROMISCUOUS
                        : PcapNetworkInterface.PromiscuousMode.NONPROMISCUOUS,
                    effectiveCaptureTimeout
                );

                // Apply BPF filter
                String bpfFilter = buildBpfFilter(configuration, tempRemoteMac, tempLocalMac);
                if (bpfFilter != null && !bpfFilter.isEmpty()) {
                    tempHandle.setFilter(bpfFilter, BpfProgram.BpfCompileMode.OPTIMIZE);
                    LOGGER.debug("Applied BPF filter: {}", bpfFilter);
                }

                LOGGER.debug("Created dedicated pcap handle on {}", nif.getName());
            }

            LOGGER.info("Raw socket transport opened on {} for protocol 0x{}, local={}, remote={}",
                nif.getName(), String.format("%04X", configuration.protocolId), tempLocalMac, tempRemoteMac);

            this.handle = tempHandle;
            this.networkInterface = nif;
            this.localMac = tempLocalMac;
            this.remoteMac = tempRemoteMac;
            this.sharedHandle = tempSharedHandle;

            // Start a packet capture thread
            startCaptureThread();

            getAuditLog().write(AuditLogEventType.CONNECT, String.format(
                "Raw socket opened on %s for protocol 0x%04X, local=%s, remote=%s",
                nif.getName(), configuration.protocolId, tempLocalMac, tempRemoteMac));
        } catch (Exception e) {
            String errorMsg = String.format("Failed to create raw socket transport for protocol 0x%04X - %s",
                configuration.protocolId, e.getMessage());
            LOGGER.error(errorMsg, e);
            getAuditLog().write(AuditLogEventType.ERROR, "Error in constructor: " + errorMsg);
            throw new TransportException(errorMsg, e);
        }
    }

    /**
     * Finds the network interface by name or returns the first available if not specified.
     */
    private PcapNetworkInterface findNetworkInterface(String interfaceName) throws PcapNativeException, TransportException {
        List<PcapNetworkInterface> allDevs = Pcaps.findAllDevs();

        if (allDevs == null || allDevs.isEmpty()) {
            throw new TransportException("No network interfaces found");
        }

        if (interfaceName != null && !interfaceName.isEmpty()) {
            // Find a specific interface
            for (PcapNetworkInterface nif : allDevs) {
                if (nif.getName().equals(interfaceName)) {
                    return nif;
                }
            }
            throw new TransportException("Network interface not found: " + interfaceName);
        } else {
            // Return the first available interface
            PcapNetworkInterface nif = allDevs.get(0);
            LOGGER.debug("No interface specified, using first available: {}", nif.getName());
            return nif;
        }
    }

    /**
     * Gets the MAC address of the network interface.
     */
    private MacAddress getInterfaceMacAddress(PcapNetworkInterface nif) throws TransportException {
        try {
            // Try to get the MAC address from link layer addresses
            for (PcapAddress addr : nif.getAddresses()) {
                if (addr.getAddress() != null && addr.getAddress().getAddress().length == 6) {
                    return MacAddress.getByAddress(addr.getAddress().getAddress());
                }
            }

            // Fallback: try to get from system network interface
            java.net.NetworkInterface javaIf = java.net.NetworkInterface.getByName(nif.getName());
            if (javaIf != null && javaIf.getHardwareAddress() != null) {
                return MacAddress.getByAddress(javaIf.getHardwareAddress());
            }

            throw new TransportException("Could not determine MAC address for interface: " + nif.getName());

        } catch (Exception e) {
            throw new TransportException("Failed to get MAC address: " + e.getMessage(), e);
        }
    }

    /**
     * Builds a BPF filter expression for packet capture.
     */
    private String buildBpfFilter(RawSocketTransportConfiguration config, MacAddress remoteMac, MacAddress localMac) {
        if (config.bpfFilter != null && !config.bpfFilter.isEmpty()) {
            return config.bpfFilter;
        }

        // Build automatic filter based on EtherType and MAC addresses
        StringBuilder filter = new StringBuilder();

        // Filter by EtherType
        filter.append(String.format("ether proto 0x%04X", config.protocolId));

        // Filter by destination MAC address (our local address)
        if (!config.promiscuousMode) {
            filter.append(String.format(" and ether dst %s", localMac.toString().replace('-', ':').toLowerCase()));
        }

        // Optionally filter by source MAC address (remote address)
        // Note: commented out as it might be too restrictive for some protocols
        // filter.append(String.format(" and ether src %s", remoteMac.toString().replace('-', ':').toLowerCase()));

        return filter.toString();
    }

    /**
     * Starts the background thread that captures packets.
     */
    /**
     * Starts the background thread that captures packets.
     * Uses getNextPacketEx() instead of loop() to avoid issues with shared handles
     * where breakLoop() can only break one of multiple threads calling loop().
     */
    private void startCaptureThread() {
        captureThread = new Thread(() -> {
            LOGGER.debug("Packet capture thread started");

            try {
                // Poll for packets one at a time, checking the open flag between calls.
                // getNextPacketEx() returns after the capture timeout, so the thread
                // can exit promptly when close() sets open=false.
                while (open && !Thread.currentThread().isInterrupted()) {
                    try {
                        Packet packet = handle.getNextPacketEx();
                        if (packet != null && packet.contains(EthernetPacket.class)) {
                            processPacket(packet.get(EthernetPacket.class));
                        }
                    } catch (TimeoutException e) {
                        // Expected when no packets arrive within the capture timeout
                    } catch (NotOpenException e) {
                        LOGGER.debug("Pcap handle closed");
                        break;
                    } catch (PcapNativeException e) {
                        if (open) {
                            LOGGER.error("Error in packet capture", e);
                        }
                        break;
                    } catch (Exception e) {
                        if (open) {
                            LOGGER.error("Error in packet capture", e);
                        }
                        break;
                    }
                }
            } finally {
                LOGGER.debug("Packet capture thread stopped");
            }
        });
        captureThread.start();
    }

    /**
     * Processes a captured Ethernet packet: filters, extracts payload, and notifies listeners.
     */
    private void processPacket(EthernetPacket ethPacket) {
        try {
            if (matchesFilter(ethPacket)) {
                byte[] payload = extractPayload(ethPacket);
                if (payload != null && payload.length > 0) {
                    receiveQueue.offer(payload);
                    LOGGER.trace("Captured packet: {} bytes from {}", payload.length,
                        ethPacket.getHeader().getSrcAddr());

                    // Notify async listener immediately (critical for timing-sensitive protocols!)
                    Runnable asyncListener = dataListener;
                    if (asyncListener != null) {
                        asyncListener.run();
                    }
                }
            }
        } catch (Exception e) {
            LOGGER.warn("Error processing captured packet", e);
        }
    }

    /**
     * Checks if the Ethernet packet matches our filter criteria.
     */
    private boolean matchesFilter(EthernetPacket ethPacket) {
        EthernetPacket.EthernetHeader header = ethPacket.getHeader();

        // Check destination MAC address (should be us or broadcast)
        MacAddress dstMac = header.getDstAddr();
        if (!dstMac.equals(localMac) && !dstMac.equals(MacAddress.ETHER_BROADCAST_ADDRESS)) {
            return false;
        }

        // Check source MAC address (should be remote)
        MacAddress srcMac = header.getSrcAddr();
        if (!srcMac.equals(remoteMac)) {
            return false;
        }

        // Check EtherType
        int etherType = header.getType().value();
        return etherType == getConfiguration().protocolId;
    }

    /**
     * Extracts the bytes we hand to the driver. In default mode this is the
     * Ethernet payload (the driver speaks an upper-layer protocol). With
     * {@code includeEthernetHeader} true (L2 protocols like PROFINET) we
     * return the full frame so the driver can inspect MAC addresses and the
     * EtherType itself.
     */
    private byte[] extractPayload(EthernetPacket ethPacket) {
        if (getConfiguration().includeEthernetHeader) {
            return ethPacket.getRawData();
        }
        Packet payload = ethPacket.getPayload();
        if (payload != null) {
            return payload.getRawData();
        }
        return null;
    }

    @Override
    public boolean isOpen() {
        return open && handle.isOpen();
    }

    /**
     * Local MAC address of the interface this transport is bound to.
     * Needed by L2 protocols like PROFINET that build Ethernet frames manually.
     */
    public byte[] getLocalMacAddress() {
        return localMac.getAddress();
    }

    /**
     * Remote MAC address this transport is configured to communicate with.
     * Needed by L2 protocols like PROFINET that build Ethernet frames manually.
     */
    public byte[] getRemoteMacAddress() {
        return remoteMac.getAddress();
    }

    /**
     * Local IPv4 address of the interface this transport is bound to, if available.
     * Returns {@code null} when the interface has no IPv4 address (for example a
     * pure L2 raw-socket interface).
     */
    public byte[] getLocalIpAddress() {
        for (org.pcap4j.core.PcapAddress address : networkInterface.getAddresses()) {
            if (address instanceof org.pcap4j.core.PcapIpV4Address && address.getAddress() != null) {
                return address.getAddress().getAddress();
            }
        }
        return null;
    }

    @Override
    public int getNumBytesAvailable() {
        readLock.lock();
        try {
            if (!isOpen()) {
                return 0;
            }

            return ringBuffer.availableForReading();
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public byte[] peekReadableBytes(int numBytes) throws TransportException {
        if (numBytes <= 0) {
            return new byte[0];
        }

        readLock.lock();
        try {
            ensureOpen();

            // Ensure we have enough data in the buffer
            while (ringBuffer.availableForReading() < numBytes && isOpen()) {
                fillRingBuffer(numBytes - ringBuffer.availableForReading());
            }

            if (ringBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, ringBuffer.availableForReading())
                );
            }

            // Peek without consuming
            return ringBuffer.peek(numBytes);
        } catch (TransportException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in peekReadableBytes: " + e.getMessage());
            throw e;
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public byte[] read(int numBytes) throws TransportException {
        if (numBytes <= 0) {
            return new byte[0];
        }

        readLock.lock();
        try {
            ensureOpen();

            // Ensure we have enough data in the buffer
            while (ringBuffer.availableForReading() < numBytes && isOpen()) {
                fillRingBuffer(numBytes - ringBuffer.availableForReading());
            }

            if (ringBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, ringBuffer.availableForReading())
                );
            }

            // Read and consume bytes
            byte[] bytes = ringBuffer.read(numBytes);

            // Log the bytes to the audit log
            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.INCOMING_BYTES, StaticHelper.ENCODE_HEX(bytes));
            }

            return bytes;
        } catch (TransportException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in read: " + e.getMessage());
            throw e;
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public void write(byte[] bytes) throws TransportException {
        if (bytes == null || bytes.length == 0) {
            return;
        }

        if (bytes.length > getConfiguration().maxFrameSize) {
            throw new TransportException(
                String.format("Frame size %d exceeds maximum %d", bytes.length, getConfiguration().maxFrameSize)
            );
        }

        writeLock.lock();
        try {
            ensureOpen();

            // L2 mode: the driver hands us a fully-formed Ethernet frame; send
            // it as-is without re-wrapping.
            if (getConfiguration().includeEthernetHeader) {
                try {
                    Packet rawPacket = EthernetPacket.newPacket(bytes, 0, bytes.length);
                    handle.sendPacket(rawPacket);
                } catch (org.pcap4j.packet.IllegalRawDataException e) {
                    throw new TransportException("Failed to parse outgoing Ethernet frame", e);
                }
                if (getAuditLog().isEnabled()) {
                    getAuditLog().write(AuditLogEventType.OUTGOING_BYTES,
                        "Write (raw L2): " + StaticHelper.ENCODE_HEX(bytes));
                }
                return;
            }

            // Build Ethernet frame
            EthernetPacket.Builder ethBuilder = new EthernetPacket.Builder()
                .srcAddr(localMac)
                .dstAddr(remoteMac)
                .type(EtherType.getInstance((short) getConfiguration().protocolId))
                .payloadBuilder(new UnknownPacket.Builder().rawData(bytes))
                .paddingAtBuild(true);

            // Add a VLAN tag if configured
            if (getConfiguration().vlanId > 0) {
                // Note: VLAN support would require Dot1qVlanTagPacket wrapper
                LOGGER.warn("VLAN tagging not yet fully implemented");
            }

            EthernetPacket packet = ethBuilder.build();

            // Send a packet
            handle.sendPacket(packet);

            LOGGER.trace("Sent {} bytes via raw socket to {}", bytes.length, remoteMac);

            // Log the bytes to the audit log
            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Write: " + StaticHelper.ENCODE_HEX(bytes));
            }
        } catch (Exception e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in write: " + e.getMessage());
            throw new TransportException("Failed to send packet", e);
        } finally {
            writeLock.unlock();
        }
    }

    @Override
    public void close() throws TransportException {
        if (!open) {
            return;
        }

        Thread captureThreadRef;
        boolean isShared;

        writeLock.lock();
        try {
            readLock.lock();
            try {
                open = false;

                captureThreadRef = captureThread;
                isShared = (sharedHandle != null);

            } finally {
                readLock.unlock();
            }
        } finally {
            writeLock.unlock();
        }

        // Stop capture thread outside of locks
        if (captureThreadRef != null) {
            captureThreadRef.interrupt();
            try {
                // Wait for the capture-thread to finish before closing the handle
                captureThreadRef.join(1000); // Wait max 1 second
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                LOGGER.warn("Interrupted while waiting for capture thread to finish");
            }
        }

        // Close the handle outside locks to avoid deadlock
        if (isShared) {
            // Release shared handle
            sharedRawSocketManager.releaseHandle(sharedHandle);
            LOGGER.debug("Released shared pcap handle");
        } else {
            // Close dedicated handle in a separate thread with timeout
            // This works around a known pcap issue on some platforms where close() can hang
            Thread closeThread = new Thread(() -> {
                try {
                    handle.close();
                    LOGGER.debug("Closed dedicated pcap handle");
                } catch (Exception e) {
                    LOGGER.warn("Error closing pcap handle: {}", e.getMessage());
                }
            });
            closeThread.start();

            try {
                closeThread.join(500); // Wait max 500ms for close
                if (closeThread.isAlive()) {
                    LOGGER.warn("Pcap handle close timed out, handle will be garbage collected");
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                LOGGER.warn("Interrupted while closing pcap handle");
            }
        }

        getAuditLog().write(AuditLogEventType.CLOSE, "Closed");
    }

    /**
     * Fills the ring buffer with at least minBytes from the receive-queue.
     */
    private void fillRingBuffer(int minBytes) throws TransportException {
        try {
            int totalRead = 0;

            while (totalRead < minBytes && isOpen()) {
                // Try to get a packet from the queue
                byte[] packet = getConfiguration().readTimeout > 0
                    ? receiveQueue.poll(getConfiguration().readTimeout, TimeUnit.MILLISECONDS)
                    : receiveQueue.take();

                if (packet == null) {
                    // Timeout or interrupted
                    break;
                }

                // Write packet data to the ring-buffer
                int written = ringBuffer.write(packet);
                totalRead += written;

                // If the packet didn't fit completely, this is a problem
                if (written < packet.length) {
                    LOGGER.warn("Packet truncated: {} bytes lost", packet.length - written);
                }
            }

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new TransportException("Read interrupted", e);
        }
    }

    /**
     * Ensures the transport is still open, throws exception otherwise.
     */
    private void ensureOpen() throws TransportException {
        if (!isOpen()) {
            throw new TransportException("Transport is closed");
        }
    }

    // ========== AsyncTransportInstance Implementation ==========

    @Override
    public void registerDataListener(Runnable listener) {
        this.dataListener = listener;
        LOGGER.debug("Data listener registered for raw socket transport");
    }

    @Override
    public void removeDataListener() {
        this.dataListener = null;
        LOGGER.debug("Data listener removed from raw socket transport");
    }

    @Override
    public void registerDisconnectListener(Consumer<Throwable> listener) {
        this.disconnectListener = listener;
        LOGGER.debug("Disconnect listener registered for raw socket transport");
    }

    @Override
    public void removeDisconnectListener() {
        this.disconnectListener = null;
        LOGGER.debug("Disconnect listener removed from raw socket transport");
    }

    /**
     * Notifies the disconnect listener if one is registered.
     *
     * @param cause the exception that caused the disconnect, or null for graceful close
     */
    private void notifyDisconnect(Throwable cause) {
        Consumer<Throwable> listener = disconnectListener;
        if (listener != null) {
            try {
                listener.accept(cause);
            } catch (Exception e) {
                LOGGER.error("Error in disconnect listener", e);
            }
        }
    }

}
