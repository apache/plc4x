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
import java.nio.channels.AsynchronousCloseException;
import java.nio.channels.SocketChannel;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.LockSupport;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;

/**
 * Java 21+ TCP transport using one virtual thread per connection doing blocking
 * {@link SocketChannel} reads/writes. On Java 21 a virtual thread blocked in a blocking-mode
 * channel read/write parks and releases its carrier (the JDK registers the fd with the NIO
 * poller), so there is no NIO {@link java.nio.channels.Selector}, no readiness-event loop, and
 * no busy-wait. The public surface and the {@link AsyncTransportInstance} callback contract are
 * unchanged: the read loop fills the {@link RingBuffer} (under {@code readLock}) and invokes the
 * registered data listener exactly as the previous selector loop did.
 */
public class TcpTransportInstance extends BaseTransportInstance<TcpTransportConfiguration> implements AsyncTransportInstance<TcpTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(TcpTransportInstance.class);
    private static final int DEFAULT_BUFFER_SIZE = 81920;
    private static final byte[] EMPTY_BYTES = new byte[0];

    private final SocketChannel socketChannel;
    private final RingBuffer ringBuffer;
    private final ByteBuffer readBuffer;  // Reused per-connection direct buffer for channel reads (confined to the read thread)
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();
    private final AtomicBoolean open = new AtomicBoolean(true);

    // Async support
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;
    private final Thread readThread;

    public TcpTransportInstance(InetSocketAddress remoteAddress, TcpTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        LOGGER.debug("TcpTransportInstance");
        this.ringBuffer = new RingBuffer(configuration.receiveBufferSize);
        this.readBuffer = ByteBuffer.allocateDirect(DEFAULT_BUFFER_SIZE);  // Reused direct buffer for channel reads

        // Opened into a local first, only promoted to the final `socketChannel` field once every
        // setup step (bind, socket options, connect) has succeeded.
        SocketChannel channel = null;
        try {
            // Open socket channel
            channel = SocketChannel.open();

            // Bind to a local address if specified
            if (configuration.localAddress != null && !configuration.localAddress.isEmpty()) {
                SocketAddress localAddr = new InetSocketAddress(configuration.localAddress, configuration.localPort);
                channel.bind(localAddr);
                LOGGER.debug("Bound to local address {}:{}", configuration.localAddress, configuration.localPort);
            }

            // Configure socket options before connecting
            channel.socket().setTcpNoDelay(configuration.tcpNoDelay);
            channel.socket().setKeepAlive(configuration.keepAlive);

            if (configuration.sendBufferSize > 0) {
                channel.socket().setSendBufferSize(configuration.sendBufferSize);
            }
            if (configuration.receiveBufferSize > 0) {
                channel.socket().setReceiveBufferSize(configuration.receiveBufferSize);
            }
            // Note: configuration.readTimeout is intentionally NOT mapped to Socket.setSoTimeout here.
            // SO_TIMEOUT has no effect on blocking SocketChannel reads, so the previous call was a
            // silent no-op. Read timeouts are enforced at the protocol/driver layer (e.g. the S7 driver
            // bounds responses via CompletableFuture timeouts), not by the transport.

            // Connect with timeout
            channel.socket().connect(remoteAddress, configuration.connectTimeout);

            // Blocking mode: on Java 21 a virtual thread blocked in read()/write() parks and
            // releases its carrier, so no selector is needed.
            channel.configureBlocking(true);

            this.socketChannel = channel;

            LOGGER.info("Connected to {}:{} with async support", remoteAddress.getHostName(), remoteAddress.getPort());

            auditLog.write(AuditLogEventType.CONNECT, String.format(
                "Connected to: %s:%d with local address: %s:%d",
                remoteAddress.getHostName(), remoteAddress.getPort(),
                getLocalAddress().getHostName(), getLocalAddress().getPort()));

            // Start the per-connection read loop on a virtual thread (Java 21+) LAST, so an
            // unchecked throw from the logging/audit above cannot leak an already-running thread
            // (the catch only handles IOException and does not stop the read loop).
            this.readThread = Thread.ofVirtual()
                .name("TCP-Read-" + remoteAddress.getHostName() + ":" + remoteAddress.getPort())
                .start(this::runReadLoop);
        } catch (IOException e) {
            String errorMsg = String.format("Failed to connect to %s:%d - %s",
                remoteAddress.getHostName(), remoteAddress.getPort(), e.getMessage());
            LOGGER.error(errorMsg, e);
            // errorMsg already embeds e.getMessage(); a single audit event avoids a duplicate.
            auditLog.write(AuditLogEventType.ERROR, "Error in constructor: " + errorMsg);
            if (channel != null) {
                try {
                    channel.close();
                } catch (IOException closeException) {
                    e.addSuppressed(closeException);
                }
            }
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
        return open.get() && socketChannel.isConnected();
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

            // Blocking write parks the virtual thread until the kernel send buffer accepts the
            // bytes — natural backpressure, no OP_WRITE registration or sleep loop needed.
            // (write() never returns -1; a broken/closed connection surfaces as
            // IOException/AsynchronousCloseException, both handled below.)
            while (writeBuffer.hasRemaining()) {
                socketChannel.write(writeBuffer);
            }

            LOGGER.trace("Wrote {} bytes", bytes.length);

            // Log the bytes to the audit log
            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Write: " + StaticHelper.ENCODE_HEX(bytes));
            }
        } catch (AsynchronousCloseException e) {
            // A concurrent close() closed the channel while we were parked in write(): normal shutdown.
            if (!open.get()) {
                return;
            }
            getAuditLog().write(AuditLogEventType.ERROR, "Error in write: " + e.getMessage());
            throw new TransportException("Failed to write data", e);
        } catch (IOException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in write: " + e.getMessage());
            throw new TransportException("Failed to write data", e);
        } finally {
            writeLock.unlock();
        }
    }

    @Override
    public void close() throws TransportException {
        // CAS so concurrent/repeated close() calls run the shutdown exactly once.
        if (!open.compareAndSet(true, false)) {
            return;
        }

        // Intentionally takes NO locks: closing the channel is what unblocks a parked read()/write().
        // Acquiring writeLock first would deadlock against a writer parked in a blocking write().
        try {
            socketChannel.close();
            // Only log/audit a successful close here; on failure the catch reports ERROR
            // and rethrows, so emitting CLOSE in finally would falsely signal a clean close.
            LOGGER.debug("TCP connection closed");
            getAuditLog().write(AuditLogEventType.CLOSE, "Closed");
        } catch (IOException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in close: " + e.getMessage());
            throw new TransportException("Failed to close connection", e);
        } finally {
            // Always join the read loop, regardless of whether the channel closed cleanly.
            // Skip the self-join when close() runs on the read thread itself (e.g. a
            // disconnect/data listener calls close()) — joining yourself only stalls for the
            // timeout and the loop already exits once open is false.
            if (readThread != null && Thread.currentThread() != readThread) {
                try {
                    readThread.join(1000);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }
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
     * Runs the listener on the read thread, guarding against a misbehaving listener so a thrown
     * exception cannot silently kill the read loop (mirrors {@link #notifyDisconnect}'s posture).
     */
    private void safeRun(Runnable listener) {
        if (listener == null) {
            return;
        }
        try {
            listener.run();
        } catch (Throwable t) {
            LOGGER.error("Data listener failed", t);
        }
    }

    /**
     * Per-connection read loop on a virtual thread: blocking read into the ring buffer, then
     * notify the data listener. No selector, no polling.
     */
    private void runReadLoop() {
        LOGGER.debug("Read loop started");
        try {
            while (open.get()) {
                int free = ringBuffer.remainingForWriting();
                if (free == 0) {
                    // Backpressure: the ring buffer is full and the consumer has not drained yet.
                    // Park briefly and re-check; never disconnect (only the codec knows frame
                    // boundaries, and cross-thread consumers like COTP can legitimately lag).
                    LockSupport.parkNanos(200_000L);
                    continue;
                }

                // Bound the read to free ring-buffer space so the buffer can never overflow.
                readBuffer.clear();
                readBuffer.limit(Math.min(readBuffer.capacity(), free));

                int bytesRead = socketChannel.read(readBuffer);  // parks vthread; releases carrier (JDK21)
                if (bytesRead == -1) {
                    // Connection closed gracefully by remote
                    LOGGER.info("Connection closed by remote");
                    open.set(false);
                    notifyDisconnect(null);
                    break;
                }
                if (bytesRead == 0) {
                    // A blocking read effectively never returns 0; harmless guard.
                    continue;
                }

                readBuffer.flip();
                readLock.lock();
                try {
                    ringBuffer.write(readBuffer);
                } finally {
                    readLock.unlock();
                }

                // Notify OUTSIDE readLock — the data is already in the ring buffer.
                safeRun(dataListener);
            }
        } catch (IOException e) {
            // If open==false, an intentional close() closed the channel (AsynchronousCloseException) —
            // a normal shutdown, not a disconnect. Only a failure while still open is a real disconnect.
            if (open.get()) {
                getAuditLog().write(AuditLogEventType.ERROR, "Error in runReadLoop: " + e.getMessage());
                LOGGER.error("Error in read loop", e);
                open.set(false);
                notifyDisconnect(e);
            }
        }
        LOGGER.debug("Read loop stopped");
    }

}
