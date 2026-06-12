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
package org.apache.plc4x.java.transport.can.virtualcan;

import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.transport.can.virtualcan.config.VirtualCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;

/**
 * In-memory virtual CAN transport instance for testing.
 * <p>
 * Frames written by one instance are synchronously broadcast to every other
 * instance on the same virtual bus via {@link VirtualCanBusManager}. No
 * background thread is needed because delivery happens inside the sender's
 * {@link #write(byte[])} call.
 * <p>
 * Implements {@link AsyncTransportInstance} so that protocol layers can register
 * a data-available listener, which is notified whenever a frame is received from
 * another instance on the bus.
 */
public class VirtualCanTransportInstance
        extends BaseTransportInstance<VirtualCanTransportConfiguration>
        implements AsyncTransportInstance<VirtualCanTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(VirtualCanTransportInstance.class);

    /** Default ring-buffer capacity — large enough for many queued CAN frames. */
    private static final int DEFAULT_BUFFER_CAPACITY = 65536;

    private volatile boolean open;
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;

    private final RingBuffer ringBuffer;
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();
    private final String busName;

    /**
     * Creates a new virtual CAN transport instance and connects it to the
     * configured virtual bus.
     *
     * @param config   the transport configuration (supplies the bus name)
     * @param auditLog the audit log for recording transport events
     */
    public VirtualCanTransportInstance(VirtualCanTransportConfiguration config, AuditLog auditLog) {
        super(config, auditLog);
        LOGGER.debug("VirtualCanTransportInstance: Pre Java 21 version");

        this.ringBuffer = new RingBuffer(DEFAULT_BUFFER_CAPACITY);
        this.busName = config.busName;
        this.open = true;

        VirtualCanBusManager.connect(busName, this);

        getAuditLog().write(AuditLogEventType.CONNECT,
                String.format("Virtual CAN transport opened on bus '%s'", busName));
        LOGGER.info("Virtual CAN transport opened on bus '{}'", busName);
    }

    @Override
    public boolean isOpen() {
        return open;
    }

    /**
     * {@inheritDoc}
     * <p>
     * Returns the number of bytes currently buffered and available for reading.
     */
    @Override
    public int getNumBytesAvailable() throws TransportException {
        readLock.lock();
        try {
            if (!open) {
                return 0;
            }
            return ringBuffer.availableForReading();
        } finally {
            readLock.unlock();
        }
    }

    /**
     * {@inheritDoc}
     * <p>
     * Peeks at buffered bytes without consuming them.
     */
    @Override
    public byte[] peekReadableBytes(int numBytes) throws TransportException {
        if (numBytes <= 0) {
            return new byte[0];
        }
        readLock.lock();
        try {
            ensureOpen();
            return ringBuffer.peek(numBytes);
        } finally {
            readLock.unlock();
        }
    }

    /**
     * {@inheritDoc}
     * <p>
     * Reads and consumes bytes from the internal ring buffer.
     */
    @Override
    public byte[] read(int numBytes) throws TransportException {
        if (numBytes <= 0) {
            return new byte[0];
        }
        readLock.lock();
        try {
            ensureOpen();
            byte[] data = ringBuffer.read(numBytes);

            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.INCOMING_BYTES, StaticHelper.ENCODE_HEX(data));
            }
            return data;
        } finally {
            readLock.unlock();
        }
    }

    /**
     * {@inheritDoc}
     * <p>
     * Writes frame bytes to the virtual bus. The bytes are broadcast synchronously
     * to every other instance connected to the same bus.
     */
    @Override
    public void write(byte[] bytes) throws TransportException {
        if (bytes == null || bytes.length == 0) {
            return;
        }
        writeLock.lock();
        try {
            ensureOpen();

            VirtualCanBusManager.broadcast(busName, this, bytes);

            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES,
                        "Write: " + StaticHelper.ENCODE_HEX(bytes));
            }
            LOGGER.trace("Sent {} bytes on virtual CAN bus '{}'", bytes.length, busName);
        } finally {
            writeLock.unlock();
        }
    }

    /**
     * {@inheritDoc}
     * <p>
     * Closes this transport instance, disconnects from the virtual bus, and
     * notifies the disconnect listener (if registered) with {@code null} to
     * indicate a graceful close.
     */
    @Override
    public void close() {
        if (!open) {
            return;
        }
        open = false;

        VirtualCanBusManager.disconnect(busName, this);

        // Notify the disconnect listener about the graceful close
        Consumer<Throwable> listener = disconnectListener;
        if (listener != null) {
            try {
                listener.accept(null);
            } catch (Exception e) {
                LOGGER.error("Error in disconnect listener", e);
            }
        }

        getAuditLog().write(AuditLogEventType.CLOSE,
                String.format("Virtual CAN transport closed on bus '%s'", busName));
        LOGGER.info("Virtual CAN transport closed on bus '{}'", busName);
    }

    // ========== AsyncTransportInstance Implementation ==========

    @Override
    public void registerDataListener(Runnable listener) {
        this.dataListener = listener;
        LOGGER.debug("Data listener registered for virtual CAN transport on bus '{}'", busName);
    }

    @Override
    public void removeDataListener() {
        this.dataListener = null;
        LOGGER.debug("Data listener removed from virtual CAN transport on bus '{}'", busName);
    }

    @Override
    public void registerDisconnectListener(Consumer<Throwable> listener) {
        this.disconnectListener = listener;
        LOGGER.debug("Disconnect listener registered for virtual CAN transport on bus '{}'", busName);
    }

    @Override
    public void removeDisconnectListener() {
        this.disconnectListener = null;
        LOGGER.debug("Disconnect listener removed from virtual CAN transport on bus '{}'", busName);
    }

    // ========== Package-Private Callback ==========

    /**
     * Called by {@link VirtualCanBusManager} when another instance on the same bus
     * writes a frame. Buffers the frame bytes and notifies the data listener.
     *
     * @param frameBytes the raw frame bytes received from another instance
     */
    void onFrameReceived(byte[] frameBytes) {
        if (!open) {
            return;
        }

        readLock.lock();
        try {
            int written = ringBuffer.write(frameBytes);
            if (written < frameBytes.length) {
                LOGGER.warn("Ring buffer overflow on bus '{}': {} of {} bytes written",
                        busName, written, frameBytes.length);
            }
        } finally {
            readLock.unlock();
        }

        // Notify the data listener outside the lock to avoid potential deadlock
        Runnable listener = dataListener;
        if (listener != null) {
            try {
                listener.run();
            } catch (Exception e) {
                LOGGER.warn("Error in data listener on bus '{}'", busName, e);
            }
        }
    }

    /**
     * Checks that the transport is still open and throws if it is not.
     *
     * @throws TransportException if the transport has been closed
     */
    private void ensureOpen() throws TransportException {
        if (!open) {
            throw new TransportException("Transport is closed");
        }
    }
}
