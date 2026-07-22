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

import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.ServerSocketChannel;
import java.nio.channels.SocketChannel;

import static org.junit.jupiter.api.Assertions.*;

class TcpTransportInstanceTest {

    private ServerSocketChannel serverChannel;
    private SocketChannel serverSideChannel;

    private TcpTransportInstance transportInstance;

    @BeforeEach
    void setUp() throws Exception {
        // Start test-server
        serverChannel = ServerSocketChannel.open();
        serverChannel.bind(new InetSocketAddress("localhost", 0));
        int serverPort = ((InetSocketAddress) serverChannel.getLocalAddress()).getPort();

        Thread acceptThread = new Thread(() -> {
            try {
                serverSideChannel = serverChannel.accept();
            } catch (Exception e) {
                e.printStackTrace();
            }
        });
        acceptThread.start();

        // Create transport-instance
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.receiveBufferSize = 81920;
        InetSocketAddress remoteAddress = new InetSocketAddress("localhost", serverPort);

        transportInstance = new TcpTransportInstance(remoteAddress, config, AuditLog.builder().build());

        // Wait for the server to accept connection
        int retries = 0;
        while (serverSideChannel == null && retries < 100) {
            Thread.sleep(10);
            retries++;
        }
        if (serverSideChannel == null) {
            throw new RuntimeException("Server failed to accept connection");
        }
    }

    @AfterEach
    void tearDown() throws Exception {
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }
        if (serverSideChannel != null && serverSideChannel.isOpen()) {
            serverSideChannel.close();
        }
        if (serverChannel != null && serverChannel.isOpen()) {
            serverChannel.close();
        }
    }

    @Test
    void testGetConfiguration() {
        TcpTransportConfiguration config = transportInstance.getConfiguration();
        assertNotNull(config);
    }

    @Test
    void testIsOpen_whenConnected() {
        assertTrue(transportInstance.isOpen());
    }

    @Test
    void testIsOpen_whenClosed() throws TransportException {
        transportInstance.close();
        assertFalse(transportInstance.isOpen());
    }

    @Test
    void testWrite_successful() throws Exception {
        byte[] data = "Hello, TCP!".getBytes();

        transportInstance.write(data);

        // Read on serverside
        ByteBuffer buffer = ByteBuffer.allocate(1024);
        int read = serverSideChannel.read(buffer);

        assertEquals(data.length, read);
        assertArrayEquals(data, java.util.Arrays.copyOf(buffer.array(), read));
    }

    @Test
    void testWrite_largeData() throws Exception {
        byte[] data = new byte[65536]; // 64KB
        for (int i = 0; i < data.length; i++) {
            data[i] = (byte) (i % 256);
        }

        transportInstance.write(data);

        // Read on serverside
        ByteBuffer buffer = ByteBuffer.allocate(data.length);
        int totalRead = 0;
        while (totalRead < data.length) {
            int read = serverSideChannel.read(buffer);
            if (read == -1) break;
            totalRead += read;
        }

        assertEquals(data.length, totalRead);
        assertArrayEquals(data, buffer.array());
    }

    @Test
    void testRead_successful() throws Exception {
        byte[] data = "Test Data".getBytes();

        // Write from serverside
        ByteBuffer buffer = ByteBuffer.wrap(data);
        serverSideChannel.write(buffer);

        // Wait a bit for data to arrive
        Thread.sleep(100);

        // Read on clientside
        byte[] result = transportInstance.read(data.length);

        assertArrayEquals(data, result);
    }

    @Test
    void testPeekReadableBytes_doesNotConsume() throws Exception {
        byte[] data = "Peek Test".getBytes();

        // Write from serverside
        ByteBuffer buffer = ByteBuffer.wrap(data);
        serverSideChannel.write(buffer);

        // Wait a bit for data to arrive
        Thread.sleep(100);

        // Peek multiple times
        byte[] peek1 = transportInstance.peekReadableBytes(4);
        byte[] peek2 = transportInstance.peekReadableBytes(4);

        assertArrayEquals(peek1, peek2);
        assertArrayEquals("Peek".getBytes(), peek1);

        // Now actually read
        byte[] read = transportInstance.read(data.length);
        assertArrayEquals(data, read);
    }

    @Test
    void testGetNumBytesAvailable() throws Exception {
        byte[] data = "Available".getBytes();

        // Initially no bytes available
        int initialAvailable = transportInstance.getNumBytesAvailable();
        assertTrue(initialAvailable >= 0);

        // Write from serverside
        ByteBuffer buffer = ByteBuffer.wrap(data);
        serverSideChannel.write(buffer);

        // Wait for data to arrive
        Thread.sleep(100);

        // Check availability
        int available = transportInstance.getNumBytesAvailable();
        assertTrue(available >= data.length);
    }

    @Test
    void testRead_whenClosed_throwsException() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.read(10)
        );
    }

    @Test
    void testWrite_whenClosed_throwsException() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.write("test".getBytes())
        );
    }

    @Test
    void testWrite_emptyArray() throws TransportException {
        // Should not throw
        assertDoesNotThrow(() -> transportInstance.write(new byte[0]));
    }

    @Test
    void testWrite_nullArray() throws TransportException {
        // Should not throw
        assertDoesNotThrow(() -> transportInstance.write(null));
    }

    @Test
    void testRead_zeroBytes() throws TransportException {
        byte[] result = transportInstance.read(0);
        assertEquals(0, result.length);
    }

    /*@Test
    void testConcurrentReadWrite() throws Exception {
        CountDownLatch startLatch = new CountDownLatch(1);
        CountDownLatch doneLatch = new CountDownLatch(2);

        // Writer thread
        Thread writer = new Thread(() -> {
            try {
                startLatch.await();
                for (int i = 0; i < 100; i++) {
                    String message = "Message " + i + "\n";
                    transportInstance.write((message).getBytes());
                    System.out.println("Written: " + message);
                }
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
                doneLatch.countDown();
            }
        });
        writer.start();

        // Server echo thread
        Thread echo = new Thread(() -> {
            try {
                startLatch.await();
                ByteBuffer buffer = ByteBuffer.allocate(1024);
                int i = 0;
                while (i < 99) {
                    buffer.clear();
                    int read = serverSideChannel.read(buffer);
                    if (read > 0) {
                        String message = new String(buffer.array());
                        System.out.println("Reading: " + message);
                        String numberString = message.substring(message.lastIndexOf(" ") + 1);
                        numberString = numberString.substring(0, numberString.indexOf("\n"));
                        i = Integer.parseInt(numberString);
                        System.out.println(i);

                        //buffer.flip();
                        //serverSideChannel.write(buffer);
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
                doneLatch.countDown();
            }
        });
        echo.start();

        startLatch.countDown();
        assertTrue(doneLatch.await(10, TimeUnit.SECONDS));
    }*/

    @Test
    void testConstructor_failedBind_doesNotLeakFileDescriptors() throws Exception {
        // getOpenFileDescriptorCount() is only exposed on Unix-flavoured JVMs
        java.lang.management.OperatingSystemMXBean rawOsBean =
            java.lang.management.ManagementFactory.getOperatingSystemMXBean();
        org.junit.jupiter.api.Assumptions.assumeTrue(
            rawOsBean instanceof com.sun.management.UnixOperatingSystemMXBean,
            "Open file descriptor count not available on this platform");
        com.sun.management.UnixOperatingSystemMXBean osBean =
            (com.sun.management.UnixOperatingSystemMXBean) rawOsBean;

        // Occupy a specific local port so binding our own channel to it deterministically fails
        // with BindException
        try (ServerSocketChannel occupier = ServerSocketChannel.open()) {
            occupier.bind(new InetSocketAddress("127.0.0.1", 0));
            int occupiedPort = ((InetSocketAddress) occupier.getLocalAddress()).getPort();

            TcpTransportConfiguration config = new TcpTransportConfiguration();
            config.receiveBufferSize = 81920;
            config.localAddress = "127.0.0.1";
            config.localPort = occupiedPort;
            // Never actually reached - bind() fails first - but must be a well-formed address.
            InetSocketAddress remoteAddress = new InetSocketAddress("127.0.0.1", occupiedPort);

            // Sample the fd count after every attempt rather than just before/after
            int attempts = 20;
            long previous = osBean.getOpenFileDescriptorCount();
            int increases = 0;
            for (int i = 0; i < attempts; i++) {
                assertThrows(TransportException.class, () ->
                    new TcpTransportInstance(remoteAddress, config, AuditLog.builder().build())
                );
                long current = osBean.getOpenFileDescriptorCount();
                if (current > previous) {
                    increases++;
                }
                previous = current;
            }

            // Without the fix each failed constructor call leaks the SocketChannel it opened before
            // bind() failed, so the fd count grows on essentially every attempt. With the fix it doesn't.
            double leakThreshold = 0.75;
            assertTrue(increases < attempts * leakThreshold,
                "Open file descriptor count increased on " + increases + "/" + attempts
                    + " failed bind attempts; expected only incidental growth, not a consistent"
                    + " per-attempt increase (possible SocketChannel leak)");
        }
    }

    @Test
    void testClose_idempotent() throws TransportException {
        transportInstance.close();
        assertFalse(transportInstance.isOpen());

        // Should not throw
        assertDoesNotThrow(() -> transportInstance.close());
        assertFalse(transportInstance.isOpen());
    }

    // ========== Async Transport Tests ==========

    @Test
    void testOnDataAvailable_listenerTriggeredRegisterData() throws Exception {
        final boolean[] listenerCalled = {false};

        // Register listener
        transportInstance.registerDataListener(() -> listenerCalled[0] = true);

        // Send data from server
        ByteBuffer buffer = ByteBuffer.wrap("Test".getBytes());
        serverSideChannel.write(buffer);

        // Wait for selector to process
        Thread.sleep(100);

        assertTrue(listenerCalled[0], "Listener should be called when data arrives");
    }

    @Test
    void testRegisterDataListener_multipleNotifications() throws Exception {
        final int[] callCount = {0};

        // Register listener
        transportInstance.registerDataListener(() -> callCount[0]++);

        // Send multiple messages
        for (int i = 0; i < 3; i++) {
            serverSideChannel.write(ByteBuffer.wrap(("Msg" + i).getBytes()));
            Thread.sleep(50);
        }

        assertTrue(callCount[0] >= 3, "Listener should be called at least 3 times");
    }

    @Test
    void testRemoveDataListener_stopsNotifications() throws Exception {
        final int[] callCount = {0};

        // Register listener
        transportInstance.registerDataListener(() -> callCount[0]++);

        // Send data
        serverSideChannel.write(ByteBuffer.wrap("Test1".getBytes()));
        Thread.sleep(50);

        int firstCount = callCount[0];
        assertTrue(firstCount > 0, "Listener should be called initially");

        // Remove listener
        transportInstance.removeDataListener();

        // Send more data
        serverSideChannel.write(ByteBuffer.wrap("Test2".getBytes()));
        Thread.sleep(50);

        // Count should not increase
        assertEquals(firstCount, callCount[0], "Listener should not be called after removal");
    }

    @Test
    void testSelectorThread_startsAutomatically() throws Exception {
        // Selector thread should be running
        assertTrue(transportInstance.isOpen());

        // Should be able to receive data immediately
        serverSideChannel.write(ByteBuffer.wrap("SelectorTest".getBytes()));
        Thread.sleep(100);

        int available = transportInstance.getNumBytesAvailable();
        assertTrue(available >= 12, "Selector should have read data into buffer");
    }

    @Test
    void testSelectorThread_stopsOnClose() throws Exception {
        transportInstance.close();

        // Wait for selector thread to stop
        Thread.sleep(200);

        assertFalse(transportInstance.isOpen());
    }

    @Test
    void testRead_insufficientBytes_throwsException() throws Exception {
        byte[] data = "Short".getBytes(); // 5 bytes

        // Write from serverside
        ByteBuffer buffer = ByteBuffer.wrap(data);
        serverSideChannel.write(buffer);

        // Wait for data to arrive
        Thread.sleep(100);

        // Try to read more bytes than available
        assertThrows(TransportException.class, () ->
            transportInstance.read(100)
        );
    }

    @Test
    void testPeekReadableBytes_insufficientBytes_throwsException() throws Exception {
        byte[] data = "Small".getBytes(); // 5 bytes

        // Write from serverside
        ByteBuffer buffer = ByteBuffer.wrap(data);
        serverSideChannel.write(buffer);

        // Wait for data to arrive
        Thread.sleep(100);

        // Try to peek more bytes than available
        assertThrows(TransportException.class, () ->
            transportInstance.peekReadableBytes(100)
        );
    }

    @Test
    void testPeekReadableBytes_zeroBytes() throws TransportException {
        byte[] result = transportInstance.peekReadableBytes(0);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekReadableBytes_negativeBytes() throws TransportException {
        byte[] result = transportInstance.peekReadableBytes(-1);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekReadableBytes_whenClosed_throwsException() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.peekReadableBytes(10)
        );
    }

    @Test
    void testGetNumBytesAvailable_whenClosed_returnsZero() throws TransportException {
        transportInstance.close();

        int available = transportInstance.getNumBytesAvailable();
        assertEquals(0, available);
    }

    @Test
    void testRead_negativeBytes() throws TransportException {
        byte[] result = transportInstance.read(-5);
        assertEquals(0, result.length);
    }

    // ========== Ring Buffer Flow Control Tests ==========

    @Test
    void testRingBufferFlowControl_fillsBufferCompletely() throws Exception {
        // Get the ring buffer size from config
        int bufferSize = transportInstance.getConfiguration().receiveBufferSize;

        // Fill the ring buffer completely by sending more data than buffer can hold
        byte[] largeData = new byte[bufferSize + 10000];
        for (int i = 0; i < largeData.length; i++) {
            largeData[i] = (byte) (i % 256);
        }

        // Write from server side
        ByteBuffer buffer = ByteBuffer.wrap(largeData);
        serverSideChannel.write(buffer);

        // Wait for data to be read into ring buffer
        Thread.sleep(200);

        // Ring buffer should be full or nearly full
        int available = transportInstance.getNumBytesAvailable();
        assertTrue(available > 0, "Some data should be available");
        assertTrue(available <= bufferSize, "Available bytes should not exceed buffer size");
    }

    @Test
    void testRingBufferFlowControl_reEnablesReadsAfterConsumingData() throws Exception {
        // This test verifies flow control works by:
        // 1. Filling the buffer exactly (no overflow)
        // 2. Draining the buffer (triggers read re-enable via reEnableReadIfNeeded())
        // 3. Sending more data and verifying it arrives (proves re-enable worked)
        //
        // If the bug exists (reads never re-enabled), step 3 would fail/timeout.

        int bufferSize = transportInstance.getConfiguration().receiveBufferSize;

        // Fill buffer exactly to capacity (don't overflow to avoid pending data in socket)
        byte[] fillData = new byte[bufferSize];
        for (int i = 0; i < fillData.length; i++) {
            fillData[i] = (byte) (i % 256);
        }

        serverSideChannel.write(ByteBuffer.wrap(fillData));
        Thread.sleep(200); // Let buffer fill

        int beforeDrain = transportInstance.getNumBytesAvailable();
        assertTrue(beforeDrain > 0, "Buffer should have data");

        // Completely drain the buffer (this calls reEnableReadIfNeeded())
        transportInstance.read(beforeDrain);

        // Give the selector thread time to process the re-enable
        Thread.sleep(50);

        // Buffer should be empty (or close to it) since we sent exactly buffer size
        int afterDrain = transportInstance.getNumBytesAvailable();
        assertTrue(afterDrain < bufferSize / 2, "Buffer should be mostly empty after drain, got " + afterDrain);

        // Send fresh data AFTER draining
        byte[] freshData = new byte[1000];
        for (int i = 0; i < freshData.length; i++) {
            freshData[i] = (byte) (i % 256);
        }

        serverSideChannel.write(ByteBuffer.wrap(freshData));

        // Wait for data to arrive - if reads weren't re-enabled, this would timeout
        int retries = 0;
        int received = afterDrain;
        int target = afterDrain + freshData.length;
        while (received < target && retries < 100) {
            Thread.sleep(20);
            received = transportInstance.getNumBytesAvailable();
            retries++;
        }

        // The key assertion: we received more data after draining
        // This proves reads were re-enabled
        assertTrue(received >= freshData.length,
            "Should receive fresh data after draining buffer. Got " + received + " bytes after " +
            retries + " retries. If this fails, reEnableReadIfNeeded() didn't work.");
    }

    @Test
    void testRingBufferFlowControl_continuousReadWriteCycle() throws Exception {
        int iterations = 10;
        int chunkSize = 5000;

        for (int i = 0; i < iterations; i++) {
            // Send a chunk of data
            byte[] sendData = new byte[chunkSize];
            for (int j = 0; j < chunkSize; j++) {
                sendData[j] = (byte) ((i * chunkSize + j) % 256);
            }

            ByteBuffer sendBuffer = ByteBuffer.wrap(sendData);
            serverSideChannel.write(sendBuffer);

            // Wait for data to arrive
            Thread.sleep(50);

            // Read the data
            int available = transportInstance.getNumBytesAvailable();
            if (available > 0) {
                byte[] receivedData = transportInstance.read(Math.min(available, chunkSize));
                assertNotNull(receivedData, "Should receive data");
                assertTrue(receivedData.length > 0, "Should receive some bytes");
            }
        }

        // After all iterations, transport should still be functional
        assertTrue(transportInstance.isOpen(), "Transport should still be open");
    }

    @Test
    void testRingBufferFlowControl_handlesPressureGracefully() throws Exception {
        // Get buffer size
        int bufferSize = transportInstance.getConfiguration().receiveBufferSize;

        // Send data in bursts to create pressure
        byte[] burstData = new byte[bufferSize / 2];
        for (int i = 0; i < burstData.length; i++) {
            burstData[i] = (byte) (i % 256);
        }

        // Send multiple bursts without reading
        for (int burst = 0; burst < 3; burst++) {
            ByteBuffer buffer = ByteBuffer.wrap(burstData);
            serverSideChannel.write(buffer);
            Thread.sleep(50);
        }

        // Buffer should be full or nearly full, but not crashed
        assertTrue(transportInstance.isOpen(), "Transport should handle buffer pressure");

        // Now drain the buffer
        int totalRead = 0;
        while (transportInstance.getNumBytesAvailable() > 0) {
            int available = transportInstance.getNumBytesAvailable();
            byte[] data = transportInstance.read(Math.min(available, 1000));
            totalRead += data.length;
        }

        assertTrue(totalRead > 0, "Should have read data from buffer");

        // After draining, should be able to receive more data
        byte[] newData = new byte[100];
        for (int i = 0; i < newData.length; i++) {
            newData[i] = (byte) (200 + i);
        }

        ByteBuffer newBuffer = ByteBuffer.wrap(newData);
        serverSideChannel.write(newBuffer);
        Thread.sleep(100);

        assertTrue(transportInstance.getNumBytesAvailable() > 0, "Should receive new data after draining");
    }

    @Test
    void testRingBufferFlowControl_smallReadsFromFullBuffer() throws Exception {
        // Fill the buffer
        int bufferSize = transportInstance.getConfiguration().receiveBufferSize;
        byte[] fillData = new byte[bufferSize];
        for (int i = 0; i < fillData.length; i++) {
            fillData[i] = (byte) (i % 256);
        }

        ByteBuffer buffer = ByteBuffer.wrap(fillData);
        serverSideChannel.write(buffer);
        Thread.sleep(200);

        // Read in small chunks
        int totalRead = 0;
        int chunkSize = 100;

        while (transportInstance.getNumBytesAvailable() > 0 && totalRead < bufferSize) {
            int available = transportInstance.getNumBytesAvailable();
            int toRead = Math.min(available, chunkSize);

            byte[] chunk = transportInstance.read(toRead);
            assertEquals(toRead, chunk.length, "Should read requested chunk size");

            // Verify data integrity
            for (int i = 0; i < chunk.length; i++) {
                assertEquals((byte) ((totalRead + i) % 256), chunk[i],
                    "Data integrity check failed at byte " + (totalRead + i));
            }

            totalRead += chunk.length;

            // Small delay to let flow control work
            Thread.sleep(10);
        }

        assertTrue(totalRead > 0, "Should have read data from buffer");
    }
}
