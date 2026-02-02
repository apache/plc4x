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

package org.apache.plc4x.java.transports.api;

import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

/**
 * Thread-safe ring buffer implementation for efficient byte storage and retrieval.
 * This implementation uses a circular buffer with separate read and write positions,
 * allowing for efficient peek and read operations without data copying.
 * <p>
 * Key features:
 * - Thread-safe operations with fine-grained locking
 * - Zero-copy peek operations
 * - Automatic wrap-around handling
 * - Dynamic capacity management
 */
public class RingBuffer {

    private final byte[] buffer;
    private final int capacity;
    private final Lock lock = new ReentrantLock();

    private int readPosition = 0;
    private int writePosition = 0;
    private int available = 0;

    /**
     * Creates a new ring buffer with the specified capacity.
     *
     * @param capacity the maximum number of bytes the buffer can hold
     * @throws IllegalArgumentException if capacity is less than or equal to zero
     */
    public RingBuffer(int capacity) {
        if (capacity <= 0) {
            throw new IllegalArgumentException("Capacity must be greater than zero");
        }
        this.capacity = capacity;
        this.buffer = new byte[capacity];
    }

    /**
     * Returns the number of bytes available for reading.
     *
     * @return number of bytes available
     */
    public int availableForReading() {
        lock.lock();
        try {
            return available;
        } finally {
            lock.unlock();
        }
    }

    /**
     * Returns the number of bytes that can be written to the buffer.
     *
     * @return number of bytes of free space
     */
    public int remainingForWriting() {
        lock.lock();
        try {
            return capacity - available;
        } finally {
            lock.unlock();
        }
    }

    /**
     * Returns the total capacity of the buffer.
     *
     * @return buffer capacity in bytes
     */
    public int capacity() {
        return capacity;
    }

    /**
     * Clears all data from the buffer, resetting read and write positions.
     */
    public void clear() {
        lock.lock();
        try {
            readPosition = 0;
            writePosition = 0;
            available = 0;
        } finally {
            lock.unlock();
        }
    }

    /**
     * Writes bytes to the buffer.
     *
     * @param data the byte array to write
     * @return number of bytes actually written
     * @throws IllegalArgumentException if data is null
     */
    public int write(byte[] data) {
        if (data == null) {
            throw new IllegalArgumentException("Data cannot be null");
        }
        return write(data, 0, data.length);
    }

    /**
     * Writes bytes to the buffer from the specified offset and length.
     *
     * @param data the byte array to write from
     * @param offset the offset in the data array to start from
     * @param length the number of bytes to write
     * @return number of bytes actually written
     * @throws IllegalArgumentException if data is null or offset/length are invalid
     */
    public int write(byte[] data, int offset, int length) {
        if (data == null) {
            throw new IllegalArgumentException("Data cannot be null");
        }
        if (offset < 0 || length < 0 || offset + length > data.length) {
            throw new IllegalArgumentException("Invalid offset or length");
        }
        if (length == 0) {
            return 0;
        }

        lock.lock();
        try {
            int bytesToWrite = Math.min(length, capacity - available);
            if (bytesToWrite == 0) {
                return 0;
            }

            // Handle wrap-around
            int firstChunk = Math.min(bytesToWrite, capacity - writePosition);
            System.arraycopy(data, offset, buffer, writePosition, firstChunk);

            if (firstChunk < bytesToWrite) {
                // Wrapped around to the beginning
                int secondChunk = bytesToWrite - firstChunk;
                System.arraycopy(data, offset + firstChunk, buffer, 0, secondChunk);
                writePosition = secondChunk;
            } else {
                writePosition = (writePosition + firstChunk) % capacity;
            }

            available += bytesToWrite;
            return bytesToWrite;

        } finally {
            lock.unlock();
        }
    }

    /**
     * Writes bytes from a ByteBuffer to the ring buffer.
     * This method avoids creating an intermediate byte array, making it more efficient
     * for NIO operations. The ByteBuffer's position will be advanced by the number of
     * bytes written.
     *
     * @param byteBuffer the ByteBuffer to read from (must be in read mode with position at start of data)
     * @return number of bytes actually written
     * @throws IllegalArgumentException if byteBuffer is null
     */
    public int write(java.nio.ByteBuffer byteBuffer) {
        if (byteBuffer == null) {
            throw new IllegalArgumentException("ByteBuffer cannot be null");
        }

        int length = byteBuffer.remaining();
        if (length == 0) {
            return 0;
        }

        lock.lock();
        try {
            int bytesToWrite = Math.min(length, capacity - available);
            if (bytesToWrite == 0) {
                return 0;
            }

            // Handle wrap-around - write directly from ByteBuffer to internal buffer
            int firstChunk = Math.min(bytesToWrite, capacity - writePosition);
            byteBuffer.get(buffer, writePosition, firstChunk);

            if (firstChunk < bytesToWrite) {
                // Wrapped around to the beginning
                int secondChunk = bytesToWrite - firstChunk;
                byteBuffer.get(buffer, 0, secondChunk);
                writePosition = secondChunk;
            } else {
                writePosition = (writePosition + firstChunk) % capacity;
            }

            available += bytesToWrite;
            return bytesToWrite;

        } finally {
            lock.unlock();
        }
    }

    /**
     * Peeks at bytes in the buffer without consuming them.
     * Returns at most numBytes, but may return fewer if not enough data is available.
     *
     * @param numBytes the number of bytes to peek
     * @return byte array containing the peeked data (may be smaller than requested)
     * @throws IllegalArgumentException if numBytes is negative
     */
    public byte[] peek(int numBytes) {
        if (numBytes < 0) {
            throw new IllegalArgumentException("Number of bytes cannot be negative");
        }
        if (numBytes == 0) {
            return new byte[0];
        }

        lock.lock();
        try {
            int bytesToPeek = Math.min(numBytes, available);
            if (bytesToPeek == 0) {
                return new byte[0];
            }

            byte[] result = new byte[bytesToPeek];

            // Handle wrap-around
            int firstChunk = Math.min(bytesToPeek, capacity - readPosition);
            System.arraycopy(buffer, readPosition, result, 0, firstChunk);

            if (firstChunk < bytesToPeek) {
                // Wrapped around to the beginning
                int secondChunk = bytesToPeek - firstChunk;
                System.arraycopy(buffer, 0, result, firstChunk, secondChunk);
            }

            return result;

        } finally {
            lock.unlock();
        }
    }

    /**
     * Reads and consumes bytes from the buffer.
     * Returns at most numBytes, but may return fewer if not enough data is available.
     *
     * @param numBytes the number of bytes to read
     * @return byte array containing the read data (may be smaller than requested)
     * @throws IllegalArgumentException if numBytes is negative
     */
    public byte[] read(int numBytes) {
        if (numBytes < 0) {
            throw new IllegalArgumentException("Number of bytes cannot be negative");
        }
        if (numBytes == 0) {
            return new byte[0];
        }

        lock.lock();
        try {
            int bytesToRead = Math.min(numBytes, available);
            if (bytesToRead == 0) {
                return new byte[0];
            }

            byte[] result = new byte[bytesToRead];

            // Handle wrap-around
            int firstChunk = Math.min(bytesToRead, capacity - readPosition);
            System.arraycopy(buffer, readPosition, result, 0, firstChunk);

            if (firstChunk < bytesToRead) {
                // Wrapped around to the beginning
                int secondChunk = bytesToRead - firstChunk;
                System.arraycopy(buffer, 0, result, firstChunk, secondChunk);
                readPosition = secondChunk;
            } else {
                readPosition = (readPosition + firstChunk) % capacity;
            }

            available -= bytesToRead;
            return result;

        } finally {
            lock.unlock();
        }
    }

    /**
     * Skips (discards) up to numBytes from the buffer.
     *
     * @param numBytes the number of bytes to skip
     * @return the actual number of bytes skipped
     * @throws IllegalArgumentException if numBytes is negative
     */
    public int skip(int numBytes) {
        if (numBytes < 0) {
            throw new IllegalArgumentException("Number of bytes cannot be negative");
        }
        if (numBytes == 0) {
            return 0;
        }

        lock.lock();
        try {
            int bytesToSkip = Math.min(numBytes, available);
            if (bytesToSkip == 0) {
                return 0;
            }

            readPosition = (readPosition + bytesToSkip) % capacity;
            available -= bytesToSkip;
            return bytesToSkip;

        } finally {
            lock.unlock();
        }
    }

    /**
     * Compacts the buffer by moving all available data to the beginning.
     * This is useful for ensuring contiguous space for writing.
     * Note: This operation is O(n) and should be used sparingly.
     */
    public void compact() {
        lock.lock();
        try {
            if (available == 0 || readPosition == 0) {
                // Nothing to compact or already at the beginning
                writePosition = available;
                readPosition = 0;
                return;
            }

            if (readPosition < writePosition) {
                // No wrap-around, simple case
                if (readPosition > 0) {
                    System.arraycopy(buffer, readPosition, buffer, 0, available);
                }
            } else {
                // Wrap-around case
                byte[] temp = new byte[available];
                int firstChunk = capacity - readPosition;
                System.arraycopy(buffer, readPosition, temp, 0, firstChunk);
                if (firstChunk < available) {
                    System.arraycopy(buffer, 0, temp, firstChunk, available - firstChunk);
                }
                System.arraycopy(temp, 0, buffer, 0, available);
            }

            readPosition = 0;
            writePosition = available;

        } finally {
            lock.unlock();
        }
    }
}