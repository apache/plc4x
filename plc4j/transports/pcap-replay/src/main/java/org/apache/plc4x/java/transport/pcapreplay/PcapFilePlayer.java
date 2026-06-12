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

import org.pcap4j.core.*;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.Packet;
import org.pcap4j.util.MacAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.EOFException;
import java.io.File;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicLong;

/**
 * Plays back packets from a PCAP file, simulating network traffic.
 * Supports speed control, looping, and MAC-address/protocol filtering.
 * Can read from both filesystem paths and classpath resources.
 */
public class PcapFilePlayer {

    private static final Logger LOGGER = LoggerFactory.getLogger(PcapFilePlayer.class);

    private final String pcapFilePath;
    private final MacAddress localMac;
    private final MacAddress remoteMac;
    private final int protocolId;
    private final double speedFactor;
    private final boolean loop;
    private final boolean onlyIncoming;
    private final boolean onlyOutgoing;
    private final int queueSize;
    private Path tempPcapFile; // Temporary file for classpath resources

    private final BlockingQueue<byte[]> packetQueue;
    private final AtomicBoolean playing = new AtomicBoolean(false);
    private final AtomicBoolean stopped = new AtomicBoolean(false);
    private final AtomicLong packetsReplayed = new AtomicLong(0);
    private Thread playerThread;

    public PcapFilePlayer(String pcapFilePath, MacAddress localMac, MacAddress remoteMac,
                          int protocolId, double speedFactor, boolean loop,
                          boolean onlyIncoming, boolean onlyOutgoing, int queueSize) {
        LOGGER.debug("PcapFilePlayer: Pre Java 21 version");
        this.pcapFilePath = resolvePcapFile(pcapFilePath);
        this.localMac = localMac;
        this.remoteMac = remoteMac;
        this.protocolId = protocolId;
        this.speedFactor = speedFactor;
        this.loop = loop;
        this.onlyIncoming = onlyIncoming;
        this.onlyOutgoing = onlyOutgoing;
        this.queueSize = queueSize;
        this.packetQueue = new LinkedBlockingQueue<>(queueSize);
    }

    /**
     * Starts replaying the PCAP file.
     */
    public void start() {
        if (playing.compareAndSet(false, true)) {
            stopped.set(false);
            playerThread = new Thread(this::playbackLoop);
            playerThread.start();
            LOGGER.info("Started PCAP replay from: {}", pcapFilePath);
        }
    }

    /**
     * Stops replaying the PCAP file.
     */
    public void stop() {
        if (stopped.compareAndSet(false, true)) {
            playing.set(false);
            if (playerThread != null) {
                playerThread.interrupt();
            }
            cleanupTempFile();
            LOGGER.info("Stopped PCAP replay (replayed {} packets)", packetsReplayed.get());
        }
    }

    /**
     * Checks if the player is currently playing.
     */
    public boolean isPlaying() {
        return playing.get();
    }

    /**
     * Gets the next packet from the queue (blocking).
     */
    public byte[] getNextPacket(long timeout, TimeUnit unit) throws InterruptedException {
        return packetQueue.poll(timeout, unit);
    }

    /**
     * Gets the number of packets currently in the queue.
     */
    public int getQueueSize() {
        return packetQueue.size();
    }

    /**
     * Main playback loop.
     */
    private void playbackLoop() {
        try {
            do {
                replayPcapFile();
            } while (loop && !stopped.get());

        } catch (Exception e) {
            if (!stopped.get()) {
                LOGGER.error("Error during PCAP replay", e);
            }
        } finally {
            playing.set(false);
            LOGGER.debug("PCAP playback thread stopped");
        }
    }

    /**
     * Replays the PCAP file once.
     */
    private void replayPcapFile() throws Exception {
        PcapHandle handle = null;

        try {
            handle = Pcaps.openOffline(pcapFilePath, PcapHandle.TimestampPrecision.MICRO);

            long previousTimestamp = -1;
            Packet packet;

            while (!stopped.get() && (packet = handle.getNextPacketEx()) != null) {
                PcapHandle.TimestampPrecision precision = handle.getTimestampPrecision();
                long timestamp = handle.getTimestamp().getTime() * 1000; // Convert to microseconds

                // Handle timing for realistic replay
                if (speedFactor > 0 && previousTimestamp >= 0) {
                    long delay = (long) ((timestamp - previousTimestamp) / speedFactor);
                    if (delay > 0) {
                        TimeUnit.MICROSECONDS.sleep(delay);
                    }
                }
                previousTimestamp = timestamp;

                // Process Ethernet packet
                if (packet.contains(EthernetPacket.class)) {
                    EthernetPacket ethPacket = packet.get(EthernetPacket.class);

                    if (matchesFilter(ethPacket)) {
                        byte[] payload = extractPayload(ethPacket);
                        if (payload != null && payload.length > 0) {
                            // Try to add to queue, drop if full
                            if (!packetQueue.offer(payload)) {
                                LOGGER.warn("Packet queue full, dropping packet");
                            } else {
                                packetsReplayed.incrementAndGet();
                                LOGGER.trace("Replayed packet: {} bytes", payload.length);
                            }
                        }
                    }
                }
            }

        } catch (EOFException e) {
            LOGGER.debug("Reached end of PCAP file");
        } finally {
            if (handle != null && handle.isOpen()) {
                handle.close();
            }
        }
    }

    /**
     * Checks if a packet matches filter criteria.
     */
    private boolean matchesFilter(EthernetPacket ethPacket) {
        EthernetPacket.EthernetHeader header = ethPacket.getHeader();

        // Check protocol if specified
        if (protocolId > 0) {
            // Mask to 16 bits to avoid sign-extension from EtherType's signed short
            int packetProto = header.getType().value() & 0xFFFF;
            if (packetProto != protocolId) {
                return false;
            }
        }

        MacAddress src = header.getSrcAddr();
        MacAddress dst = header.getDstAddr();

        // Filter by direction
        if (onlyIncoming && onlyOutgoing) {
            // Both flags set - this is contradictory, reject
            return false;
        }

        if (onlyIncoming) {
            // Only packets from remote to local
            return src.equals(remoteMac) && dst.equals(localMac);
        } else if (onlyOutgoing) {
            // Only packets from local to remote
            return src.equals(localMac) && dst.equals(remoteMac);
        } else {
            // Any packet involving local or remote MAC addresses
            return src.equals(localMac) || src.equals(remoteMac) ||
                dst.equals(localMac) || dst.equals(remoteMac);
        }
    }

    /**
     * Extracts payload from an Ethernet packet.
     */
    private byte[] extractPayload(EthernetPacket ethPacket) {
        Packet payload = ethPacket.getPayload();
        if (payload != null) {
            return payload.getRawData();
        }
        return null;
    }

    /**
     * Gets statistics about replay.
     */
    public long getPacketsReplayed() {
        return packetsReplayed.get();
    }

    /**
     * Resolves a PCAP file path, handling both filesystem paths and classpath resources.
     * If the path starts with "/", it's treated as a classpath resource and extracted to a temp file.
     *
     * @param pcapFile The PCAP file path (filesystem or classpath)
     * @return The resolved filesystem path
     */
    private String resolvePcapFile(String pcapFile) {
        // Check if it's a classpath resource (starts with '/')
        if (pcapFile.startsWith("/")) {
            return extractClasspathResource(pcapFile);
        }

        // Treat as a filesystem path
        File file = new File(pcapFile);
        if (!file.exists()) {
            LOGGER.warn("PCAP file not found: {}", pcapFile);
            return pcapFile; // Return as-is, error will occur when trying to open
        }
        return file.getAbsolutePath();
    }

    /**
     * Extracts a classpath resource to a temporary file.
     *
     * @param resourcePath The classpath resource path
     * @return The absolute path to the temporary file
     */
    private String extractClasspathResource(String resourcePath) {
        try {
            // Load resource from the classpath
            InputStream resourceStream = getClass().getResourceAsStream(resourcePath);
            if (resourceStream == null) {
                LOGGER.error("Classpath resource not found: {}", resourcePath);
                return resourcePath; // Return as-is, error will occur when trying to open
            }

            // Create temporary file
            String fileName = resourcePath.substring(resourcePath.lastIndexOf('/') + 1);
            String baseName = fileName.contains(".") ? fileName.substring(0, fileName.lastIndexOf('.')) : fileName;
            String extension = fileName.contains(".") ? fileName.substring(fileName.lastIndexOf('.')) : ".pcap";

            tempPcapFile = Files.createTempFile(baseName + "_", extension);
            Files.copy(resourceStream, tempPcapFile, StandardCopyOption.REPLACE_EXISTING);
            resourceStream.close();

            String tempPath = tempPcapFile.toAbsolutePath().toString();
            LOGGER.debug("Extracted classpath resource {} to temporary file: {}", resourcePath, tempPath);

            return tempPath;

        } catch (Exception e) {
            LOGGER.error("Failed to extract classpath resource: " + resourcePath, e);
            return resourcePath; // Return as-is, error will occur when trying to open
        }
    }

    /**
     * Cleans up temporary files created from classpath resources.
     */
    private void cleanupTempFile() {
        if (tempPcapFile != null) {
            try {
                if (Files.exists(tempPcapFile)) {
                    Files.delete(tempPcapFile);
                    LOGGER.debug("Deleted temporary PCAP file: {}", tempPcapFile);
                }
            } catch (Exception e) {
                LOGGER.warn("Failed to delete temporary PCAP file: {}", tempPcapFile, e);
            }
            tempPcapFile = null;
        }
    }
}
