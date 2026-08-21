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

import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.udp.config.UdpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.DatagramChannel;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.util.Optional;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;

/**
 * Java 21+ optimized UDP transport (using virtual threads) implementation using NIO DatagramChannel with async support.
 * Implements AsyncTransportInstance for event-driven I/O without polling.
 * Supports both dedicated and shared socket modes for optimal resource usage.
 */
public class UdpTransportInstance implements AsyncTransportInstance<UdpTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(UdpTransportInstance.class);
    private static final int DEFAULT_BUFFER_SIZE = 8192;
    private static final byte[] EMPTY_BYTES = new byte[0];

    private final UdpTransportConfiguration configuration;
    private final SharedUdpSocketManager sharedUdpSocketManager;
    private final DatagramChannel channel;
    private final InetSocketAddress remoteAddress;
    private final SharedUdpSocketManager.SharedSocket sharedSocket; // null if not shared
    private final RingBuffer ringBuffer;
    private final ByteBuffer readBuffer;  // Pre-allocated direct buffer for zero-copy I/O
    private final int maxPacketSize;
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();
    private volatile boolean open = true;

    // Async support
    private final Selector selector;
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;
    private final Thread selectorThread;

    public UdpTransportInstance(InetSocketAddress remoteAddress, UdpTransportConfiguration configuration, SharedUdpSocketManager sharedUdpSocketManager, AuditLog auditLog) throws TransportException {
        LOGGER.debug("UdpTransportInstance: Java 21 version");
        this.configuration = configuration;
        this.sharedUdpSocketManager = sharedUdpSocketManager;

        // Use the default max packet size if not configured (handles test cases where defaults aren't applied)
        this.maxPacketSize = configuration.maxPacketSize > 0 ? configuration.maxPacketSize : 65507;
        this.ringBuffer = new RingBuffer(DEFAULT_BUFFER_SIZE);
        this.readBuffer = ByteBuffer.allocateDirect(DEFAULT_BUFFER_SIZE);

        try {
            DatagramChannel tempChannel;
            SharedUdpSocketManager.SharedSocket tempSharedSocket = null;

            if (configuration.shareSocket) {
                // Use shared socket-manager
                tempSharedSocket = sharedUdpSocketManager
                    .acquireSocket(configuration.localAddress, configuration.localPort);
                tempChannel = tempSharedSocket.getChannel();

                LOGGER.debug("Using shared UDP socket on {}:{}",
                    configuration.localAddress, configuration.localPort);
            } else {
                // Create a dedicated socket
                tempChannel = DatagramChannel.open();

                // Bind to a local address and port
                SocketAddress bindAddress;
                if (configuration.localAddress != null && !configuration.localAddress.isEmpty()) {
                    bindAddress = new InetSocketAddress(configuration.localAddress, configuration.localPort);
                } else {
                    bindAddress = new InetSocketAddress(configuration.localPort);
                }
                tempChannel.bind(bindAddress);

                LOGGER.debug("Created dedicated UDP socket bound to {}", bindAddress);
            }

            // Configure socket options
            tempChannel.socket().setBroadcast(configuration.broadcast);
            tempChannel.socket().setReuseAddress(configuration.reuseAddress);

            if (configuration.sendBufferSize > 0) {
                tempChannel.socket().setSendBufferSize(configuration.sendBufferSize);
            }
            if (configuration.receiveBufferSize > 0) {
                tempChannel.socket().setReceiveBufferSize(configuration.receiveBufferSize);
            }
            if (configuration.readTimeout > 0) {
                tempChannel.socket().setSoTimeout(configuration.readTimeout);
            }

            if(!tempChannel.isConnected()) {
                tempChannel.connect(remoteAddress);
            }

            // Configure non-blocking mode for NIO selector
            tempChannel.configureBlocking(false);

            // Create a selector for async I/O
            this.selector = Selector.open();
            tempChannel.register(selector, SelectionKey.OP_READ);

            // Start selector thread using virtual thread (Java 21+)
            this.selectorThread = Thread.ofVirtual()
                .name("UDP-Selector-" + remoteAddress.getHostName() + ":" + remoteAddress.getPort())
            .start(this::runSelectorLoop);

            this.channel = tempChannel;
            this.remoteAddress = remoteAddress;
            this.sharedSocket = tempSharedSocket;

            LOGGER.info("Connected to {}:{} with async support", remoteAddress.getHostName(), remoteAddress.getPort());
            auditLog.write(AuditLogEventType.CONNECT, String.format(
                "Connected to: %s:%d with local address: %s:%d",
                remoteAddress.getHostName(), remoteAddress.getPort(),
                getLocalAddress().orElseThrow().getHostName(), getLocalAddress().orElseThrow().getPort()));
        } catch (IOException e) {
            String errorMsg = String.format("Failed to create UDP transport for %s:%d - %s",
                remoteAddress.getHostName(), remoteAddress.getPort(), e.getMessage());
            LOGGER.error(errorMsg, e);
            auditLog.write(AuditLogEventType.ERROR, "Error in constructor: " + errorMsg);
            auditLog.write(AuditLogEventType.ERROR, "Error in constructor: " + e.getMessage());
            throw new TransportException(errorMsg, e);
        }
    }

    @Override
    public UdpTransportConfiguration getConfiguration() {
        return configuration;
    }

    @Override
    public boolean isOpen() {
        return open && channel.isOpen();
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
            return EMPTY_BYTES;
        }

        readLock.lock();
        try {
            ensureOpen();

            // Fill the ring buffer if necessary
            getNumBytesAvailable();

            if (ringBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, ringBuffer.availableForReading())
                );
            }

            // Peek without consuming
            return ringBuffer.peek(numBytes);
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public byte[] read(int numBytes) throws TransportException {
        if (numBytes <= 0) {
            return EMPTY_BYTES;
        }

        readLock.lock();
        try {
            ensureOpen();

            // Fill the ring buffer if necessary
            getNumBytesAvailable();

            if (ringBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, ringBuffer.availableForReading())
                );
            }

            // Read and consume bytes
            return ringBuffer.read(numBytes);

        } finally {
            readLock.unlock();
        }
    }

    @Override
    public void write(byte[] bytes) throws TransportException {
        if (bytes == null || bytes.length == 0) {
            return;
        }

        if (bytes.length > maxPacketSize) {
            throw new TransportException(
                String.format("Packet size %d exceeds maximum %d", bytes.length, maxPacketSize)
            );
        }

        writeLock.lock();
        try {
            ensureOpen();

            ByteBuffer writeBuffer = ByteBuffer.wrap(bytes);
            int sent = channel.send(writeBuffer, remoteAddress);

            if (sent != bytes.length) {
                throw new TransportException(
                    String.format("Failed to send complete packet: sent %d of %d bytes", sent, bytes.length)
                );
            }

            LOGGER.trace("Sent {} bytes to {}", sent, remoteAddress);

        } catch (IOException e) {
            throw new TransportException("Failed to send UDP packet", e);
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

                // Wake up selector
                selector.wakeup();

                if (sharedSocket != null) {
                    // Release shared socket (may not close if other instances use it)
                    sharedUdpSocketManager.releaseSocket(sharedSocket);
                    LOGGER.debug("Released shared UDP socket");
                } else {
                    // Close channel and selector
                    channel.close();
                    selector.close();

                    // Wait for selector thread to finish
                    if (selectorThread != null) {
                        try {
                            selectorThread.join(1000);
                        } catch (InterruptedException e) {
                            Thread.currentThread().interrupt();
                        }
                    }

                    LOGGER.debug("Closed dedicated UDP socket");
                }

            } catch (IOException e) {
                throw new TransportException("Failed to close UDP transport", e);
            } finally {
                readLock.unlock();
            }
        } finally {
            writeLock.unlock();
        }
    }

    /**
     * Receives a packet from the channel and stores it in the ring buffer.
     * Filters packets to only accept from the connected remote address.
     */
    private void receivePacket() throws TransportException {
        try {
            readBuffer.clear();
            SocketAddress source = channel.receive(readBuffer);

            if (source == null) {
                // No packet received (non-blocking mode)
                return;
            }

            readBuffer.flip();

            // Filter packets: only accept from our connected remote address
            if (shouldAcceptPacket(source)) {
                // Copy to ring buffer
                byte[] data = new byte[readBuffer.remaining()];
                readBuffer.get(data);
                ringBuffer.write(data);

                LOGGER.trace("Received packet of {} bytes from {}", data.length, source);
            } else {
                LOGGER.trace("Ignored packet from {} (expected {})", source, remoteAddress);
            }

        } catch (IOException e) {
            throw new TransportException("Failed to receive UDP packet", e);
        }
    }

    /**
     * Checks if a packet from the given source should be accepted.
     * When using shared sockets, multiple instances may receive packets for different remote addresses.
     */
    private boolean shouldAcceptPacket(SocketAddress source) {
        if (!(source instanceof InetSocketAddress inetSource)) {
            return false;
        }

        // Accept if from our remote address
        return inetSource.getAddress().equals(remoteAddress.getAddress()) &&
               inetSource.getPort() == remoteAddress.getPort();
    }

    /**
     * Ensures the transport is still open, throws exception otherwise.
     */
    private void ensureOpen() throws TransportException {
        if (!isOpen()) {
            throw new TransportException("Transport is closed");
        }
    }

    public Optional<InetSocketAddress> getLocalAddress() {
        if(channel != null) {
            return Optional.of((InetSocketAddress) channel.socket().getLocalSocketAddress());
        }
        return Optional.empty();
    }

    // ========== AsyncTransportInstance Implementation ==========

    @Override
    public void registerDataListener(Runnable listener) {
        this.dataListener = listener;
        LOGGER.debug("Data listener registered");
    }

    @Override
    public void removeDataListener() {
        this.dataListener = null;
        LOGGER.debug("Data listener removed");
    }

    @Override
    public void registerDisconnectListener(Consumer<Throwable> listener) {
        this.disconnectListener = listener;
        LOGGER.debug("Disconnect listener registered");
    }

    @Override
    public void removeDisconnectListener() {
        this.disconnectListener = null;
        LOGGER.debug("Disconnect listener removed");
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

    /**
     * Selector loop that runs in a virtual thread and notifies listeners when data arrives.
     * This is the core of the async implementation - no polling needed in the driver!
     */
    private void runSelectorLoop() {
        LOGGER.debug("UDP selector loop started");

        try {
            runSelectorLoopInternal();
        } catch (Throwable t) {
            // This loop is the only thing reading the channel. If it ever exits abnormally
            // there is nobody left to notice, so the transport must stop claiming to be open
            // instead of leaving callers waiting on a connection that can no longer receive.
            LOGGER.error("UDP selector loop failed", t);
            if (open) {
                open = false;
                notifyDisconnect(t);
            }
        }

        LOGGER.debug("UDP selector loop stopped");
    }

    private void runSelectorLoopInternal() {
        while (open && !Thread.currentThread().isInterrupted()) {
            try {
                // Block until events are available (no CPU waste!)
                int readyChannels = selector.select();

                if (readyChannels == 0) {
                    continue;
                }

                var selectedKeys = selector.selectedKeys();
                var iterator = selectedKeys.iterator();

                while (iterator.hasNext()) {
                    SelectionKey key = iterator.next();
                    iterator.remove();

                    if (!key.isValid()) {
                        continue;
                    }

                    if (key.isReadable()) {
                        // Data available - read it into the ring buffer
                        readLock.lock();
                        try {
                            receivePacket(); // Existing method

                            // Notify the listener that data is available
                            Runnable listener = dataListener;
                            if (listener != null) {
                                listener.run();
                            }
                        } catch (TransportException e) {
                            if (open) {
                                LOGGER.error("Error receiving packet", e);
                            }
                        } catch (RuntimeException e) {
                            // Any host can send a datagram to a UDP port, so one we cannot
                            // handle costs us this receive cycle and nothing more.
                            if (open) {
                                LOGGER.error("Error handling received packet", e);
                            }
                        } finally {
                            readLock.unlock();
                        }
                    }
                }

            } catch (IOException e) {
                if (open) {
                    LOGGER.error("Error in selector loop", e);
                    open = false;
                    notifyDisconnect(e);
                }
                break;
            }
        }
    }

}
