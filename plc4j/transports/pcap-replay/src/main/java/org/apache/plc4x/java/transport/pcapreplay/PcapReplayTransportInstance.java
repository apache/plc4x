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
package org.apache.plc4x.java.transport.pcapreplay;

import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.transport.pcapreplay.config.PcapReplayTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.Pcaps;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.Packet;
import org.pcap4j.util.MacAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.Arrays;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

/**
 * Java 21 PCAP Replay transport instance that simulates a real Ethernet device
 * by replaying packets from a Wireshark capture file.
 */
public class PcapReplayTransportInstance extends BaseTransportInstance<PcapReplayTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(PcapReplayTransportInstance.class);

    private final PcapFilePlayer player;
    private final MacAddress localMac;
    private final MacAddress remoteMac;
    private final RingBuffer ringBuffer;
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();
    private volatile boolean open = true;
    private Path tempPcapFile; // Temporary file for classpath resources

    public PcapReplayTransportInstance(PcapReplayTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        this.ringBuffer = new RingBuffer(configuration.maxFrameSize);

        try {
            // Resolve PCAP file from the filesystem or classpath
            String pcapFilePath = resolvePcapFile(configuration.pcapFile);

            // Auto-detect MAC addresses if needed
            MacAddress tempLocalMac;
            MacAddress tempRemoteMac;

            if (configuration.autoDetectMacAddresses) {
                MacAddress[] detectedMacs = detectMacAddresses(pcapFilePath, configuration.protocolId);
                tempLocalMac = configuration.localAddress != null && !configuration.localAddress.isEmpty()
                    ? MacAddress.getByName(configuration.localAddress)
                    : detectedMacs[0];
                tempRemoteMac = configuration.remoteAddress != null && !configuration.remoteAddress.isEmpty()
                    ? MacAddress.getByName(configuration.remoteAddress)
                    : detectedMacs[1];
            } else {
                if (configuration.localAddress == null || configuration.localAddress.isEmpty()) {
                    throw new TransportException("Local MAC address must be specified when auto-detect is disabled");
                }
                if (configuration.remoteAddress == null || configuration.remoteAddress.isEmpty()) {
                    throw new TransportException("Remote MAC address must be specified when auto-detect is disabled");
                }
                tempLocalMac = MacAddress.getByName(configuration.localAddress);
                tempRemoteMac = MacAddress.getByName(configuration.remoteAddress);
            }

            LOGGER.debug("PCAP replay - Local MAC: {}, Remote MAC: {}", tempLocalMac, tempRemoteMac);

            if (configuration.mockPlayer == null) {
                // Create player
                PcapFilePlayer tempPlayer = new PcapFilePlayer(
                    pcapFilePath,
                    tempLocalMac,
                    tempRemoteMac,
                    configuration.protocolId,
                    configuration.speedFactor,
                    configuration.loop,
                    configuration.onlyIncomingPackets,
                    configuration.onlyOutgoingPackets,
                    configuration.packetQueueSize
                );

                LOGGER.info("PCAP replay transport created for {} (speed: {}x, loop: {})",
                    configuration.pcapFile,
                    configuration.speedFactor > 0 ? configuration.speedFactor : "max",
                    configuration.loop);

                this.player = tempPlayer;
            } else {
                this.player = configuration.mockPlayer;
            }
            this.localMac = tempLocalMac;
            this.remoteMac = tempRemoteMac;

            // Auto-start if configured
            if (configuration.autoStart) {
                player.start();
            }

            getAuditLog().write(AuditLogEventType.CONNECT, String.format(
                "PCAP replay transport created for %s (speed: %sx, loop: %s)",
                configuration.pcapFile,
                configuration.speedFactor > 0 ? String.valueOf(configuration.speedFactor) : "max",
                configuration.loop));
        } catch (Exception e) {
            String errorMsg = String.format("Failed to create PCAP replay transport for %s - %s",
                configuration.pcapFile, e.getMessage());
            LOGGER.error(errorMsg, e);
            getAuditLog().write(AuditLogEventType.ERROR, "Error in constructor: " + errorMsg);
            throw new TransportException(errorMsg, e);
        }
    }

    /**
     * Detects local and remote MAC addresses from the first suitable packet in the PCAP.
     */
    private MacAddress[] detectMacAddresses(String pcapFile, int protocolId) throws Exception {
        PcapHandle handle = null;
        try {
            handle = Pcaps.openOffline(pcapFile, PcapHandle.TimestampPrecision.MICRO);

            Packet packet;
            while ((packet = handle.getNextPacketEx()) != null) {
                if (packet.contains(EthernetPacket.class)) {
                    EthernetPacket ethPacket = packet.get(EthernetPacket.class);
                    EthernetPacket.EthernetHeader header = ethPacket.getHeader();

                    // Check if matches protocol filter
                    if (protocolId > 0) {
                        // Mask to 16 bits to avoid sign-extension from EtherType's signed short
                        int packetProto = header.getType().value() & 0xFFFF;
                        if (packetProto != protocolId) {
                            continue;
                        }
                    }

                    // Use first packet's addresses
                    MacAddress src = header.getSrcAddr();
                    MacAddress dst = header.getDstAddr();

                    // Avoid broadcast addresses
                    if (dst.equals(MacAddress.ETHER_BROADCAST_ADDRESS)) {
                        continue;
                    }

                    LOGGER.debug("Auto-detected MAC addresses from PCAP: local={}, remote={}", dst, src);
                    return new MacAddress[]{dst, src}; // Assume dst is local, src is remote
                }
            }

            throw new TransportException("No suitable packets found in PCAP file for MAC detection");

        } finally {
            if (handle != null && handle.isOpen()) {
                handle.close();
            }
        }
    }

    @Override
    public boolean isOpen() {
        return open;
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
                int beforeFill = ringBuffer.availableForReading();
                fillRingBuffer(numBytes - beforeFill);
                if (ringBuffer.availableForReading() == beforeFill) {
                    // No progress — player has no more data to provide
                    break;
                }
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
                int beforeFill = ringBuffer.availableForReading();
                fillRingBuffer(numBytes - beforeFill);
                if (ringBuffer.availableForReading() == beforeFill) {
                    // No progress — player has no more data to provide
                    break;
                }
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

        writeLock.lock();
        try {
            ensureOpen();

            // PCAP replay is read-only - writing is a no-op or could be logged
            LOGGER.debug("Write operation on PCAP replay transport (simulated, {} bytes)", bytes.length);

            // Log the write attempt to audit log
            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Write (simulated): " + StaticHelper.ENCODE_HEX(bytes));
            }

            // In a real implementation, you might want to:
            // 1. Store written data for validation
            // 2. Trigger specific responses from the PCAP based on written data
            // 3. Log the write operation for debugging purposes
            throw new TransportException("PCAP replay transport is read-only");
        } catch (TransportException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in write: " + e.getMessage());
            throw e;
        } finally {
            writeLock.unlock();
        }
    }

    @Override
    public void close() throws TransportException {
        if (!open) {
            return;
        }

        writeLock.lock();
        try {
            readLock.lock();
            try {
                open = false;
                player.stop();
                LOGGER.debug("Closed PCAP replay transport (replayed {} packets)", player.getPacketsReplayed());

                // Clean up temporary file if created
                if (tempPcapFile != null) {
                    try {
                        Files.deleteIfExists(tempPcapFile);
                        LOGGER.debug("Deleted temporary PCAP file: {}", tempPcapFile);
                    } catch (Exception e) {
                        LOGGER.warn("Failed to delete temporary PCAP file: {}", e.getMessage());
                    }
                }

                getAuditLog().write(AuditLogEventType.CLOSE, String.format(
                    "Closed PCAP replay transport (replayed %d packets)", player.getPacketsReplayed()));
            } finally {
                readLock.unlock();
            }
        } finally {
            writeLock.unlock();
        }
    }

    /**
     * Starts the PCAP replay if not already started.
     */
    public void startReplay() {
        if (!player.isPlaying()) {
            player.start();
        }
    }

    /**
     * Stops the PCAP replay.
     */
    public void stopReplay() {
        if (player.isPlaying()) {
            player.stop();
        }
    }

    /**
     * Checks if replay is currently active.
     */
    public boolean isReplaying() {
        return player.isPlaying();
    }

    /**
     * Gets the number of packets replayed so far.
     */
    public long getPacketsReplayed() {
        return player.getPacketsReplayed();
    }

    /**
     * Fills the ring buffer with at least minBytes from the player.
     */
    /**
     * What is left of a packet too big for the buffer, waiting for room. Held so the replayed byte
     * stream stays the stream that was captured; a frame larger than the buffer is otherwise
     * delivered with its tail missing and everything after it misread.
     */
    private byte[] pendingFrameRemainder;

    private void fillRingBuffer(int minBytes) throws TransportException {
        try {
            int totalRead = 0;

            while (totalRead < minBytes && isOpen()) {
                // Anything left over from a packet that did not fit last time goes first. A replay
                // is worth having only if it replays what was captured, and the codec reading this
                // decides where one message ends and the next begins - so a gap in the middle of a
                // packet is not a lost packet, it is every value after it read from the wrong
                // offset.
                if (pendingFrameRemainder != null) {
                    int carried = ringBuffer.write(pendingFrameRemainder);
                    totalRead += carried;
                    pendingFrameRemainder = carried < pendingFrameRemainder.length
                        ? Arrays.copyOfRange(pendingFrameRemainder, carried, pendingFrameRemainder.length)
                        : null;
                    if (pendingFrameRemainder != null) {
                        // Still no room. The consumer has to drain before there will be.
                        break;
                    }
                    continue;
                }

                // Get packet from player
                byte[] packet = getConfiguration().readTimeout > 0
                    ? player.getNextPacket(getConfiguration().readTimeout, TimeUnit.MILLISECONDS)
                    : player.getNextPacket(Long.MAX_VALUE, TimeUnit.DAYS);

                if (packet == null) {
                    // Timeout or no more packets
                    break;
                }

                // Write packet data to the ring buffer
                int written = ringBuffer.write(packet);
                totalRead += written;

                if (written < packet.length) {
                    // Keep what did not fit and deliver it before anything newer, rather than
                    // logging the loss and carrying on with a stream that no longer matches the
                    // capture.
                    pendingFrameRemainder = Arrays.copyOfRange(packet, written, packet.length);
                    LOGGER.debug("Parking {} bytes of a {} byte packet until the buffer drains",
                        pendingFrameRemainder.length, packet.length);
                    break;
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

    /**
     * Resolves a PCAP file path, supporting both filesystem and classpath resources.
     * If the path starts with '/', it's treated as a classpath resource and extracted
     * to a temporary file. Otherwise, it's treated as a filesystem path.
     *
     * @param pcapFile The PCAP file path (filesystem or classpath resource starting with '/')
     * @return The resolved absolute file path
     * @throws TransportException If the file cannot be found or accessed
     */
    private String resolvePcapFile(String pcapFile) throws TransportException {
        // Check if it's a classpath resource (starts with '/')
        if (pcapFile.startsWith("/")) {
            return extractClasspathResource(pcapFile);
        }

        // Treat as a filesystem path
        File file = new File(pcapFile);
        if (!file.exists()) {
            throw new TransportException("PCAP file not found: " + pcapFile);
        }
        if (!file.canRead()) {
            throw new TransportException("Cannot read PCAP file: " + pcapFile);
        }

        return file.getAbsolutePath();
    }

    /**
     * Extracts a classpath resource to a temporary file.
     *
     * @param resourcePath The classpath resource path
     * @return The absolute path to the temporary file
     * @throws TransportException If the resource cannot be found or extracted
     */
    private String extractClasspathResource(String resourcePath) throws TransportException {
        try {
            // Load resource from the classpath
            InputStream resourceStream = getClass().getResourceAsStream(resourcePath);
            if (resourceStream == null) {
                throw new TransportException("Classpath resource not found: " + resourcePath);
            }

            // Create a temporary file
            String fileName = resourcePath.substring(resourcePath.lastIndexOf('/') + 1);
            String prefix = fileName.contains(".") ? fileName.substring(0, fileName.lastIndexOf('.')) : fileName;
            String suffix = fileName.contains(".") ? fileName.substring(fileName.lastIndexOf('.')) : ".tmp";

            tempPcapFile = Files.createTempFile(prefix + "_", suffix);

            // Copy resource to the temporary file
            Files.copy(resourceStream, tempPcapFile, StandardCopyOption.REPLACE_EXISTING);
            resourceStream.close();

            // Mark for deletion on JVM exit (backup cleanup)
            tempPcapFile.toFile().deleteOnExit();

            LOGGER.debug("Extracted classpath resource {} to temporary file: {}", resourcePath, tempPcapFile);

            return tempPcapFile.toAbsolutePath().toString();

        } catch (Exception e) {
            throw new TransportException("Failed to extract classpath resource " + resourcePath + ": " + e.getMessage(), e);
        }
    }
}
