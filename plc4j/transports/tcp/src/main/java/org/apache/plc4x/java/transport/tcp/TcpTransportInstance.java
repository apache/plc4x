/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.transport.tcp;

import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.SocketChannel;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;

/**
 * Java 21+ optimized version using virtual threads - TCP transport implementation using NIO SocketChannel with async support.
 * Implements AsyncTransportInstance for event-driven I/O without polling.
 */
public class TcpTransportInstance extends BaseTransportInstance<TcpTransportConfiguration> implements AsyncTransportInstance<TcpTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(TcpTransportInstance.class);
    private static final int DEFAULT_BUFFER_SIZE = 81920;
    private static final byte[] EMPTY_BYTES = new byte[0];

    private final SocketChannel socketChannel;
    private final RingBuffer ringBuffer;
    private final ByteBuffer readBuffer;  // Pre-allocated direct buffer for zero-copy I/O
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();
    private volatile boolean open = true;

    // Async support
    private final Selector selector;
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;
    private final Thread selectorThread;

    public TcpTransportInstance(InetSocketAddress remoteAddress, TcpTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        LOGGER.debug("TcpTransportInstance");
        this.ringBuffer = new RingBuffer(configuration.receiveBufferSize);
        this.readBuffer = ByteBuffer.allocateDirect(DEFAULT_BUFFER_SIZE);  // Direct buffer for zero-copy

        try {
            // Open socket channel
            this.socketChannel = SocketChannel.open();

            // Bind to a local address if specified
            if (configuration.localAddress != null && !configuration.localAddress.isEmpty()) {
                SocketAddress localAddr = new InetSocketAddress(configuration.localAddress, configuration.localPort);
                socketChannel.bind(localAddr);
                LOGGER.debug("Bound to local address {}:{}", configuration.localAddress, configuration.localPort);
            }

            // Configure socket options before connecting
            socketChannel.socket().setTcpNoDelay(configuration.tcpNoDelay);
            socketChannel.socket().setKeepAlive(configuration.keepAlive);

            if (configuration.sendBufferSize > 0) {
                socketChannel.socket().setSendBufferSize(configuration.sendBufferSize);
            }
            if (configuration.receiveBufferSize > 0) {
                socketChannel.socket().setReceiveBufferSize(configuration.receiveBufferSize);
            }
            if (configuration.readTimeout > 0) {
                socketChannel.socket().setSoTimeout(configuration.readTimeout);
            }

            // Connect with timeout
            socketChannel.socket().connect(remoteAddress, configuration.connectTimeout);

            // Configure non-blocking mode for NIO selector
            socketChannel.configureBlocking(false);

            // Create a selector for async I/O
            this.selector = Selector.open();
            socketChannel.register(selector, SelectionKey.OP_READ);

            // Start selector thread using virtual thread (Java 21+)
            this.selectorThread = Thread.ofVirtual()
                .name("TCP-Selector-" + remoteAddress.getHostName() + ":" + remoteAddress.getPort())
                .start(this::runSelectorLoop);

            LOGGER.info("Connected to {}:{} with async support", remoteAddress.getHostName(), remoteAddress.getPort());

            auditLog.write(AuditLogEventType.CONNECT, String.format(
                "Connected to: %s:%d with local address: %s:%d",
                remoteAddress.getHostName(), remoteAddress.getPort(),
                getLocalAddress().getHostName(), getLocalAddress().getPort()));
        } catch (IOException e) {
            String errorMsg = String.format("Failed to connect to %s:%d - %s",
                remoteAddress.getHostName(), remoteAddress.getPort(), e.getMessage());
            LOGGER.error(errorMsg, e);
            auditLog.write(AuditLogEventType.ERROR, "Error in constructor: " + errorMsg);
            auditLog.write(AuditLogEventType.ERROR, "Error in constructor: " + e.getMessage());
            throw new TransportException(errorMsg, e);
        }
    }

    public InetSocketAddress getRemoteAddress() {
        return (InetSocketAddress) socketChannel.socket().getRemoteSocketAddress();
    }

    public InetSocketAddress getLocalAddress() {
        return (InetSocketAddress) socketChannel.socket().getLocalSocketAddress();
    }

    @Override
    public boolean isOpen() {
        return open && socketChannel.isConnected();
    }

    @Override
    public int getNumBytesAvailable() throws TransportException {
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
            return EMPTY_BYTES;
        }

        readLock.lock();
        try {
            ensureOpen();

            if (ringBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, ringBuffer.availableForReading())
                );
            }

            // Read and consume bytes
            byte[] bytes = ringBuffer.read(numBytes);

            // Re-enable read operations if they were disabled due to full buffer
            // Now that we've freed up space, the selector can read more data
            reEnableReadIfNeeded();

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

            ByteBuffer writeBuffer = ByteBuffer.wrap(bytes);

            while (writeBuffer.hasRemaining()) {
                int written = socketChannel.write(writeBuffer);
                if (written == -1) {
                    open = false;
                    throw new TransportException("Connection closed while writing");
                }

                // If no bytes were written and buffer is full, register for write operations
                // and wait until the channel becomes writable (prevents CPU spinning)
                if (written == 0 && writeBuffer.hasRemaining()) {
                    try {
                        // Temporarily register interest in write operations
                        SelectionKey key = socketChannel.keyFor(selector);
                        if (key != null && key.isValid()) {
                            key.interestOps(key.interestOps() | SelectionKey.OP_WRITE);
                            selector.wakeup();

                            // Wait a short time for the socket to become writable
                            // This prevents tight CPU spinning when the socket buffer is full
                            Thread.sleep(1);

                            // Remove write interest to avoid unnecessary wake-ups
                            key.interestOps(SelectionKey.OP_READ);
                        }
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        throw new TransportException("Write interrupted", e);
                    }
                }
            }

            LOGGER.trace("Wrote {} bytes to {}", bytes.length, socketChannel.getRemoteAddress());

            // Log the bytes to the audit log
            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Write: " + StaticHelper.ENCODE_HEX(bytes));
            }
        } catch (TransportException | IOException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in write: " + e.getMessage());
            throw new TransportException("Failed to write data", e);
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

                // Close socket channel
                socketChannel.close();

                // Close selector
                selector.close();

                // Wait for the selector thread to finish
                if (selectorThread != null) {
                    try {
                        selectorThread.join(1000);
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    }
                }

                LOGGER.debug("TCP connection closed");
                getAuditLog().write(AuditLogEventType.CLOSE, "Closed");
            } catch (IOException e) {
                getAuditLog().write(AuditLogEventType.ERROR, "Error in close: " + e.getMessage());
                throw new TransportException("Failed to close connection", e);
            } finally {
                readLock.unlock();
            }
        } finally {
            writeLock.unlock();
        }
    }

    /**
     * Ensures the connection is still open, throws exception otherwise.
     */
    private void ensureOpen() throws TransportException {
        if (!isOpen()) {
            throw new TransportException("Transport is closed");
        }
    }

    /**
     * Re-enables read operations on the selector if they were previously disabled
     * due to a full ring buffer. This should be called after reading from the ring
     * buffer to allow new data to be received.
     */
    private void reEnableReadIfNeeded() {
        try {
            SelectionKey key = socketChannel.keyFor(selector);
            if (key != null && key.isValid()) {
                // Check if read interest is currently disabled
                if ((key.interestOps() & SelectionKey.OP_READ) == 0) {
                    // Re-enable read operations
                    key.interestOps(key.interestOps() | SelectionKey.OP_READ);
                    // Wake up the selector to process the new interest ops
                    selector.wakeup();
                    LOGGER.debug("Re-enabled read operations after buffer space freed");
                }
            }
        } catch (Exception e) {
            // Log but don't throw - this is a best-effort operation
            LOGGER.warn("Failed to re-enable read operations", e);
        }
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
        LOGGER.debug("Selector loop started");

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
                        boolean notifyListener = false;
                        boolean connectionClosed = false;
                        readLock.lock();
                        try {
                            // Check available space in ring buffer before reading
                            int availableSpace = ringBuffer.remainingForWriting();
                            if (availableSpace == 0) {
                                LOGGER.warn("Ring buffer is full, temporarily disabling read operations");
                                key.interestOps(key.interestOps() & ~SelectionKey.OP_READ);
                                continue;
                            }

                            // Limit read buffer to available space in ring-buffer to prevent data loss
                            readBuffer.clear();
                            readBuffer.limit(Math.min(readBuffer.capacity(), availableSpace));

                            int bytesRead = socketChannel.read(readBuffer);

                            if (bytesRead > 0) {
                                readBuffer.flip();

                                // Write directly from ByteBuffer to ring buffer (avoiding intermediate byte array allocation)
                                int bytesWritten = ringBuffer.write(readBuffer);
                                if (bytesWritten < bytesRead) {
                                    String message = String.format("Ring buffer write incomplete. Expected to write " +
                                            "%d bytes but only wrote %d bytes. This should not happen.",
                                        bytesRead, bytesWritten);
                                    LOGGER.error(message);
                                    getAuditLog().write(AuditLogEventType.ERROR, message);
                                }

                                notifyListener = true;
                            } else if (bytesRead == -1) {
                                // Connection closed gracefully by remote
                                LOGGER.info("Connection closed by remote");
                                open = false;
                                connectionClosed = true;
                            }
                        } finally {
                            readLock.unlock();
                        }

                        // Notify listener OUTSIDE the readLock — the data is already
                        // in the ring buffer, so holding the lock during response
                        // processing would unnecessarily block the I/O path and cause
                        // reentrant lock overhead in processIncomingData().
                        if (notifyListener) {
                            Runnable listener = dataListener;
                            if (listener != null) {
                                listener.run();
                            }
                        } else if (connectionClosed) {
                            notifyDisconnect(null);
                            break;
                        }
                    }
                }

            } catch (IOException e) {
                getAuditLog().write(AuditLogEventType.ERROR, "Error in runSelectorLoop: " + e.getMessage());
                if (open) {
                    LOGGER.error("Error in selector loop", e);
                    open = false;
                    notifyDisconnect(e);
                }
                break;
            }
        }

        LOGGER.debug("Selector loop stopped");
    }

}
