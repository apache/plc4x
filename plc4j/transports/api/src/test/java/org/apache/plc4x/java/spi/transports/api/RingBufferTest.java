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

package org.apache.plc4x.java.spi.transports.api;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeEach;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.concurrent.*;

import static org.junit.jupiter.api.Assertions.*;

class RingBufferTest {

    private RingBuffer buffer;

    @BeforeEach
    void setUp() {
        buffer = new RingBuffer(10);
    }

    @Test
    void testConstructorWithValidCapacity() {
        RingBuffer rb = new RingBuffer(100);
        assertEquals(100, rb.capacity());
        assertEquals(0, rb.availableForReading());
        assertEquals(100, rb.remainingForWriting());
    }

    @Test
    void testConstructorWithInvalidCapacity() {
        assertThrows(IllegalArgumentException.class, () -> new RingBuffer(0));
        assertThrows(IllegalArgumentException.class, () -> new RingBuffer(-1));
    }

    @Test
    void testWriteAndRead() {
        byte[] data = {1, 2, 3, 4, 5};
        int written = buffer.write(data);

        assertEquals(5, written);
        assertEquals(5, buffer.availableForReading());
        assertEquals(5, buffer.remainingForWriting());

        byte[] read = buffer.read(5);
        assertArrayEquals(data, read);
        assertEquals(0, buffer.availableForReading());
        assertEquals(10, buffer.remainingForWriting());
    }

    @Test
    void testWriteWithOffset() {
        byte[] data = {1, 2, 3, 4, 5, 6, 7, 8};
        int written = buffer.write(data, 2, 4);

        assertEquals(4, written);
        assertEquals(4, buffer.availableForReading());

        byte[] read = buffer.read(4);
        assertArrayEquals(new byte[]{3, 4, 5, 6}, read);
    }

    @Test
    void testWriteMoreThanCapacity() {
        byte[] data = new byte[15];
        for (int i = 0; i < 15; i++) {
            data[i] = (byte) i;
        }

        // TODO: This should fail and throw an exception
        int written = buffer.write(data);
        assertEquals(10, written); // Only capacity amount written
        assertEquals(10, buffer.availableForReading());
        assertEquals(0, buffer.remainingForWriting());
    }

    @Test
    void testWriteNull() {
        assertThrows(IllegalArgumentException.class, () -> buffer.write((byte[]) null));
        assertThrows(IllegalArgumentException.class, () -> buffer.write(null, 0, 0));
    }

    @Test
    void testWriteWithInvalidOffsetLength() {
        byte[] data = {1, 2, 3, 4, 5};
        assertThrows(IllegalArgumentException.class, () -> buffer.write(data, -1, 3));
        assertThrows(IllegalArgumentException.class, () -> buffer.write(data, 0, -1));
        assertThrows(IllegalArgumentException.class, () -> buffer.write(data, 0, 10));
        assertThrows(IllegalArgumentException.class, () -> buffer.write(data, 3, 5));
    }

    @Test
    void testPeek() {
        byte[] data = {1, 2, 3, 4, 5};
        buffer.write(data);

        byte[] peeked1 = buffer.peek(3);
        assertArrayEquals(new byte[]{1, 2, 3}, peeked1);
        assertEquals(5, buffer.availableForReading()); // Available unchanged

        byte[] peeked2 = buffer.peek(3);
        assertArrayEquals(new byte[]{1, 2, 3}, peeked2);
        assertEquals(5, buffer.availableForReading()); // Still unchanged

        byte[] read = buffer.read(5);
        assertArrayEquals(data, read);
    }

    @Test
    void testPeekMoreThanAvailableForReading() {
        byte[] data = {1, 2, 3};
        buffer.write(data);

        byte[] peeked = buffer.peek(10);
        assertEquals(3, peeked.length);
        assertArrayEquals(data, peeked);
    }

    @Test
    void testPeekEmpty() {
        byte[] peeked = buffer.peek(5);
        assertEquals(0, peeked.length);
    }

    @Test
    void testPeekNegative() {
        assertThrows(IllegalArgumentException.class, () -> buffer.peek(-1));
    }

    @Test
    void testReadMoreThanAvailableForReading() {
        byte[] data = {1, 2, 3};
        buffer.write(data);

        byte[] read = buffer.read(10);
        assertEquals(3, read.length);
        assertArrayEquals(data, read);
        assertEquals(0, buffer.availableForReading());
    }

    @Test
    void testReadEmpty() {
        byte[] read = buffer.read(5);
        assertEquals(0, read.length);
    }

    @Test
    void testReadNegative() {
        assertThrows(IllegalArgumentException.class, () -> buffer.read(-1));
    }

    @Test
    void testReadZero() {
        byte[] data = {1, 2, 3};
        buffer.write(data);

        byte[] read = buffer.read(0);
        assertEquals(0, read.length);
        assertEquals(3, buffer.availableForReading());
    }

    @Test
    void testClear() {
        byte[] data = {1, 2, 3, 4, 5};
        buffer.write(data);

        assertEquals(5, buffer.availableForReading());
        buffer.clear();
        assertEquals(0, buffer.availableForReading());
        assertEquals(10, buffer.remainingForWriting());
    }

    @Test
    void testSkip() {
        byte[] data = {1, 2, 3, 4, 5, 6, 7, 8};
        buffer.write(data);

        int skipped = buffer.skip(3);
        assertEquals(3, skipped);
        assertEquals(5, buffer.availableForReading());

        byte[] read = buffer.read(5);
        assertArrayEquals(new byte[]{4, 5, 6, 7, 8}, read);
    }

    @Test
    void testSkipMoreThanAvailableForReading() {
        byte[] data = {1, 2, 3};
        buffer.write(data);

        int skipped = buffer.skip(10);
        assertEquals(3, skipped);
        assertEquals(0, buffer.availableForReading());
    }

    @Test
    void testSkipNegative() {
        assertThrows(IllegalArgumentException.class, () -> buffer.skip(-1));
    }

    @Test
    void testWrapAround() {
        // Fill buffer
        byte[] data1 = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
        buffer.write(data1);

        // Read some to create space at the beginning
        buffer.read(6);
        assertEquals(4, buffer.availableForReading());
        assertEquals(6, buffer.remainingForWriting());

        // Write data that will wrap around
        byte[] data2 = {11, 12, 13, 14, 15, 16};
        int written = buffer.write(data2);
        assertEquals(6, written);
        assertEquals(10, buffer.availableForReading());

        // Read all and verify order
        byte[] read = buffer.read(10);
        assertArrayEquals(new byte[]{7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, read);
    }

    @Test
    void testWrapAroundPeek() {
        // Fill buffer
        buffer.write(new byte[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10});

        // Read to create wrap-around condition
        buffer.read(7);

        // Write data that wraps
        buffer.write(new byte[]{11, 12, 13, 14, 15, 16, 17});

        // Peek across the wrap boundary
        byte[] peeked = buffer.peek(10);
        assertArrayEquals(new byte[]{8, 9, 10, 11, 12, 13, 14, 15, 16, 17}, peeked);
    }

    @Test
    void testCompactNoWrapAround() {
        // Write and read to move positions forward
        buffer.write(new byte[]{1, 2, 3, 4, 5});
        buffer.read(2);

        // Compact should move data to beginning
        buffer.compact();

        assertEquals(3, buffer.availableForReading());
        byte[] read = buffer.read(3);
        assertArrayEquals(new byte[]{3, 4, 5}, read);
    }

    @Test
    void testCompactWithWrapAround() {
        // Create wrap-around condition
        buffer.write(new byte[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10});
        buffer.read(7);
        buffer.write(new byte[]{11, 12, 13, 14, 15, 16, 17});

        // Compact should handle wrap-around
        buffer.compact();

        assertEquals(10, buffer.availableForReading());
        byte[] read = buffer.read(10);
        assertArrayEquals(new byte[]{8, 9, 10, 11, 12, 13, 14, 15, 16, 17}, read);
    }

    @Test
    void testCompactEmpty() {
        buffer.compact();
        assertEquals(0, buffer.availableForReading());
        assertEquals(10, buffer.remainingForWriting());
    }

    @Test
    void testMultipleWriteReadCycles() {
        for (int i = 0; i < 100; i++) {
            byte[] data = {(byte) i, (byte) (i + 1), (byte) (i + 2)};
            buffer.write(data);

            byte[] read = buffer.read(3);
            assertArrayEquals(data, read);
            assertEquals(0, buffer.availableForReading());
        }
    }

    @Test
    void testPartialReadsAndWrites() {
        buffer.write(new byte[]{1, 2, 3, 4, 5});
        buffer.read(2);
        buffer.write(new byte[]{6, 7, 8});
        buffer.read(3);
        buffer.write(new byte[]{9, 10, 11, 12});

        assertEquals(7, buffer.availableForReading());
        byte[] read = buffer.read(7);
        assertArrayEquals(new byte[]{6, 7, 8, 9, 10, 11, 12}, read);
    }

    @Test
    void testConcurrentWriteAndRead() throws InterruptedException {
        RingBuffer concurrentBuffer = new RingBuffer(1000);
        int numOperations = 1000;
        CountDownLatch startLatch = new CountDownLatch(1);
        CountDownLatch doneLatch = new CountDownLatch(2);
        List<Throwable> errors = new CopyOnWriteArrayList<>();

        // Writer thread
        Thread writer = new Thread(() -> {
            try {
                startLatch.await();
                Random random = new Random(42);
                for (int i = 0; i < numOperations; i++) {
                    byte[] data = new byte[random.nextInt(10) + 1];
                    random.nextBytes(data);
                    while (concurrentBuffer.write(data) < data.length) {
                        Thread.sleep(1);
                    }
                }
            } catch (Throwable t) {
                errors.add(t);
            } finally {
                doneLatch.countDown();
            }
        });

        // Reader thread
        Thread reader = new Thread(() -> {
            try {
                startLatch.await();
                Random random = new Random(43);
                for (int i = 0; i < numOperations; i++) {
                    int toRead = random.nextInt(10) + 1;
                    while (concurrentBuffer.availableForReading() < toRead) {
                        Thread.sleep(1);
                    }
                    byte[] data = concurrentBuffer.read(toRead);
                    assertNotNull(data);
                }
            } catch (Throwable t) {
                errors.add(t);
            } finally {
                doneLatch.countDown();
            }
        });

        writer.start();
        reader.start();
        startLatch.countDown();

        assertTrue(doneLatch.await(30, TimeUnit.SECONDS), "Concurrent test timed out");

        if (!errors.isEmpty()) {
            throw new AssertionError("Concurrent test failed: " + errors.get(0).getMessage(), errors.get(0));
        }
    }

    @Test
    void testMultipleReaders() throws InterruptedException, ExecutionException {
        RingBuffer sharedBuffer = new RingBuffer(1000);
        sharedBuffer.write(new byte[1000]); // Fill the buffer

        int numReaders = 10;
        ExecutorService executor = Executors.newFixedThreadPool(numReaders);
        List<Future<Integer>> futures = new ArrayList<>();

        for (int i = 0; i < numReaders; i++) {
            futures.add(executor.submit(() -> {
                int totalRead = 0;
                while (sharedBuffer.availableForReading() > 0) {
                    byte[] data = sharedBuffer.read(10);
                    totalRead += data.length;
                }
                return totalRead;
            }));
        }

        int totalReadByAll = 0;
        for (Future<Integer> future : futures) {
            totalReadByAll += future.get();
        }

        executor.shutdown();
        assertTrue(executor.awaitTermination(5, TimeUnit.SECONDS));

        assertEquals(1000, totalReadByAll);
        assertEquals(0, sharedBuffer.availableForReading());
    }

    @Test
    void testStressTest() {
        RingBuffer stressBuffer = new RingBuffer(100);
        Random random = new Random(12345);

        for (int i = 0; i < 10000; i++) {
            // Random operation
            int op = random.nextInt(5);
            switch (op) {
                case 0: // Write
                    byte[] writeData = new byte[random.nextInt(20)];
                    random.nextBytes(writeData);
                    stressBuffer.write(writeData);
                    break;
                case 1: // Read
                    stressBuffer.read(random.nextInt(20));
                    break;
                case 2: // Peek
                    stressBuffer.peek(random.nextInt(20));
                    break;
                case 3: // Skip
                    stressBuffer.skip(random.nextInt(20));
                    break;
                case 4: // Clear
                    if (random.nextInt(10) == 0) {
                        stressBuffer.clear();
                    }
                    break;
            }

            // Verify invariants
            assertTrue(stressBuffer.availableForReading() >= 0);
            assertTrue(stressBuffer.availableForReading() <= stressBuffer.capacity());
            assertTrue(stressBuffer.remainingForWriting() >= 0);
            assertTrue(stressBuffer.remainingForWriting() <= stressBuffer.capacity());
            assertEquals(stressBuffer.capacity(), stressBuffer.availableForReading() + stressBuffer.remainingForWriting());
        }
    }

    // Tests for ByteBuffer write method

    @Test
    void testWriteFromByteBuffer() {
        java.nio.ByteBuffer byteBuffer = java.nio.ByteBuffer.allocate(5);
        byteBuffer.put(new byte[]{1, 2, 3, 4, 5});
        byteBuffer.flip();

        int written = buffer.write(byteBuffer);

        assertEquals(5, written);
        assertEquals(5, buffer.availableForReading());
        assertEquals(5, buffer.remainingForWriting());
        assertEquals(0, byteBuffer.remaining()); // ByteBuffer should be consumed

        byte[] read = buffer.read(5);
        assertArrayEquals(new byte[]{1, 2, 3, 4, 5}, read);
    }

    @Test
    void testWriteFromDirectByteBuffer() {
        java.nio.ByteBuffer byteBuffer = java.nio.ByteBuffer.allocateDirect(5);
        byteBuffer.put(new byte[]{10, 20, 30, 40, 50});
        byteBuffer.flip();

        int written = buffer.write(byteBuffer);

        assertEquals(5, written);
        assertEquals(5, buffer.availableForReading());
        assertEquals(0, byteBuffer.remaining());

        byte[] read = buffer.read(5);
        assertArrayEquals(new byte[]{10, 20, 30, 40, 50}, read);
    }

    @Test
    void testWriteFromByteBufferWithWrapAround() {
        // Fill the buffer partially
        buffer.write(new byte[]{1, 2, 3, 4, 5, 6, 7, 8});
        buffer.read(5); // Read 5 bytes to create space at the beginning

        // Now write 7 bytes which will wrap around
        java.nio.ByteBuffer byteBuffer = java.nio.ByteBuffer.allocate(7);
        byteBuffer.put(new byte[]{10, 20, 30, 40, 50, 60, 70});
        byteBuffer.flip();

        int written = buffer.write(byteBuffer);

        assertEquals(7, written);
        assertEquals(10, buffer.availableForReading()); // 3 old + 7 new
        assertEquals(0, byteBuffer.remaining());

        // Read all and verify
        byte[] read = buffer.read(10);
        assertArrayEquals(new byte[]{6, 7, 8, 10, 20, 30, 40, 50, 60, 70}, read);
    }

    @Test
    void testWriteFromByteBufferWhenBufferFull() {
        // Fill the buffer
        buffer.write(new byte[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10});

        java.nio.ByteBuffer byteBuffer = java.nio.ByteBuffer.allocate(5);
        byteBuffer.put(new byte[]{11, 12, 13, 14, 15});
        byteBuffer.flip();

        int written = buffer.write(byteBuffer);

        assertEquals(0, written); // Nothing written
        assertEquals(10, buffer.availableForReading());
        assertEquals(5, byteBuffer.remaining()); // ByteBuffer unchanged
    }

    @Test
    void testWriteFromByteBufferPartial() {
        // Fill buffer with 7 bytes, leaving only 3 bytes available
        buffer.write(new byte[]{1, 2, 3, 4, 5, 6, 7});

        java.nio.ByteBuffer byteBuffer = java.nio.ByteBuffer.allocate(5);
        byteBuffer.put(new byte[]{10, 20, 30, 40, 50});
        byteBuffer.flip();

        int written = buffer.write(byteBuffer);

        assertEquals(3, written); // Only 3 bytes written
        assertEquals(10, buffer.availableForReading());
        assertEquals(2, byteBuffer.remaining()); // 2 bytes left in ByteBuffer
    }

    @Test
    void testWriteFromEmptyByteBuffer() {
        java.nio.ByteBuffer byteBuffer = java.nio.ByteBuffer.allocate(0);

        int written = buffer.write(byteBuffer);

        assertEquals(0, written);
        assertEquals(0, buffer.availableForReading());
    }

    @Test
    void testWriteFromNullByteBuffer() {
        assertThrows(IllegalArgumentException.class, () -> buffer.write((java.nio.ByteBuffer) null));
    }

    @Test
    void testWriteFromByteBufferWithLimit() {
        java.nio.ByteBuffer byteBuffer = java.nio.ByteBuffer.allocate(10);
        byteBuffer.put(new byte[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10});
        byteBuffer.flip();
        byteBuffer.limit(5); // Limit to first 5 bytes

        int written = buffer.write(byteBuffer);

        assertEquals(5, written);
        assertEquals(5, buffer.availableForReading());
        assertEquals(0, byteBuffer.remaining());

        byte[] read = buffer.read(5);
        assertArrayEquals(new byte[]{1, 2, 3, 4, 5}, read);
    }
}