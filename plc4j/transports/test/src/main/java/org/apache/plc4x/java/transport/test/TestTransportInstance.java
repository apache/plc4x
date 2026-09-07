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
package org.apache.plc4x.java.transport.test;

import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.transport.test.config.TestTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

/**
 * Test transport instance implementation using in-memory byte streams.
 * This implementation is useful for testing drivers without actual network connections.
 * <p>
 * Key features:
 * - Thread-safe read/write operations
 * - In-memory ring buffer for data storage
 * - Methods to inject test data and retrieve written data
 */
public class TestTransportInstance extends BaseTransportInstance<TestTransportConfiguration>
        implements TransportInstance<TestTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(TestTransportInstance.class);
    private static final byte[] EMPTY_BYTES = new byte[0];

    private final RingBuffer readBuffer;
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();
    private volatile boolean open = true;

    // For testing: store written data separately
    private final RingBuffer writeBuffer;

    public TestTransportInstance(TestTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        this.readBuffer = new RingBuffer(configuration.receiveBufferSize);
        this.writeBuffer = new RingBuffer(configuration.receiveBufferSize);

        LOGGER.debug("Test transport instance created");
        auditLog.write(AuditLogEventType.CONNECT,
            "Test transport instance created");
    }

    /**
     * Injects test data into the read buffer.
     * This simulates receiving data from a remote endpoint.
     *
     * @param data the data to inject
     * @return number of bytes actually written to the buffer
     */
    public int injectTestData(byte[] data) {
        if (data == null || data.length == 0) {
            return 0;
        }

        readLock.lock();
        try {
            int written = readBuffer.write(data);
            LOGGER.debug("Injected {} bytes of test data", written);
            return written;
        } finally {
            readLock.unlock();
        }
    }

    /**
     * Retrieves data that was written to the transport.
     * This is useful for verifying what the driver sent.
     *
     * @param numBytes number of bytes to retrieve
     * @return the retrieved data
     */
    public byte[] getWrittenData(int numBytes) {
        writeLock.lock();
        try {
            if (writeBuffer.availableForReading() < numBytes) {
                return writeBuffer.read(writeBuffer.availableForReading());
            }
            return writeBuffer.read(numBytes);
        } finally {
            writeLock.unlock();
        }
    }

    /**
     * Retrieves all data that was written to the transport.
     *
     * @return all written data
     */
    public byte[] getAllWrittenData() {
        writeLock.lock();
        try {
            return writeBuffer.read(writeBuffer.availableForReading());
        } finally {
            writeLock.unlock();
        }
    }

    /**
     * Waits for the expected number of bytes to be written, then retrieves them.
     * This is useful for test synchronization where we need to wait for the driver to send data.
     *
     * @param expectedBytes the number of bytes to wait for
     * @param timeoutMs     timeout in milliseconds
     * @return the written data
     * @throws TransportException if timeout occurs before expected bytes are available
     */
    public byte[] waitForWrittenData(int expectedBytes, long timeoutMs) throws TransportException {
        long startTime = System.currentTimeMillis();
        long deadline = startTime + timeoutMs;

        while (System.currentTimeMillis() < deadline) {
            writeLock.lock();
            try {
                int available = writeBuffer.availableForReading();
                if (available >= expectedBytes) {
                    byte[] data = writeBuffer.read(expectedBytes);
                    LOGGER.debug("Retrieved {} bytes of written data after {}ms",
                        expectedBytes, System.currentTimeMillis() - startTime);
                    return data;
                }
            } finally {
                writeLock.unlock();
            }

            // Short sleep to avoid busy-waiting
            try {
                Thread.sleep(1);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                throw new TransportException("Interrupted while waiting for written data");
            }
        }

        writeLock.lock();
        try {
            int available = writeBuffer.availableForReading();
            throw new TransportException(String.format(
                "Timeout waiting for written data: expected %d bytes but only %d available after %dms",
                expectedBytes, available, timeoutMs));
        } finally {
            writeLock.unlock();
        }
    }

    /**
     * Returns the number of bytes available in the write buffer.
     *
     * @return number of bytes written
     */
    public int getNumBytesWritten() {
        writeLock.lock();
        try {
            return writeBuffer.availableForReading();
        } finally {
            writeLock.unlock();
        }
    }

    @Override
    public boolean isOpen() {
        return open;
    }

    @Override
    public int getNumBytesAvailable() throws TransportException {
        readLock.lock();
        try {
            if (!isOpen()) {
                return 0;
            }
            return readBuffer.availableForReading();
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

            if (readBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, readBuffer.availableForReading())
                );
            }

            return readBuffer.peek(numBytes);
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

            if (readBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, readBuffer.availableForReading())
                );
            }

            byte[] bytes = readBuffer.read(numBytes);

            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Read: " + StaticHelper.ENCODE_HEX(bytes));
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

            int written = writeBuffer.write(bytes);
            if (written < bytes.length) {
                throw new TransportException(
                    String.format("Could only write %d of %d bytes", written, bytes.length)
                );
            }

            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Write: " + StaticHelper.ENCODE_HEX(bytes));
            }
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
                LOGGER.debug("Test transport closed");
                getAuditLog().write(AuditLogEventType.CLOSE, "Closed");
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

}
