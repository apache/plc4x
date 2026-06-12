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

import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.udp.config.UdpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.DatagramChannel;
import java.util.concurrent.atomic.AtomicReference;

import static org.junit.jupiter.api.Assertions.*;

class UdpTransportInstanceTest {

    private DatagramChannel serverChannel;

    private UdpTransportInstance transportInstance;
    private InetSocketAddress clientAddress = new InetSocketAddress("localhost", 0);

    @BeforeEach
    void setUp() throws Exception {
        // Create server channel
        serverChannel = DatagramChannel.open();
        serverChannel.bind(new InetSocketAddress("localhost", 0));
        int serverPort = ((InetSocketAddress) serverChannel.getLocalAddress()).getPort();

        // Create transport-instance
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        InetSocketAddress remoteAddress = new InetSocketAddress("localhost", serverPort);

        transportInstance = new UdpTransportInstance(remoteAddress, config, new SharedUdpSocketManager(), AuditLog.builder().build());

        // Get the local address of the client.
        clientAddress = transportInstance.getLocalAddress().orElseThrow();
    }

    @AfterEach
    void tearDown() throws Exception {
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }
        if (serverChannel != null && serverChannel.isOpen()) {
            serverChannel.close();
        }
    }

    @Test
    void testGetConfiguration() {
        UdpTransportConfiguration config = transportInstance.getConfiguration();
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
        byte[] data = "Hello UDP!".getBytes();

        transportInstance.write(data);

        // Receive on serverside
        ByteBuffer buffer = ByteBuffer.allocate(1024);
        InetSocketAddress source = (InetSocketAddress) serverChannel.receive(buffer);

        assertNotNull(source);
        assertEquals(data.length, buffer.position());

        buffer.flip();
        byte[] received = new byte[buffer.remaining()];
        buffer.get(received);
        assertArrayEquals(data, received);
    }

    @Test
    void testWrite_maxPacketSize() throws Exception {
        byte[] data = new byte[1024];
        for (int i = 0; i < data.length; i++) {
            data[i] = (byte) (i % 256);
        }

        transportInstance.write(data);

        ByteBuffer buffer = ByteBuffer.allocate(2048);
        serverChannel.receive(buffer);

        buffer.flip();
        assertEquals(data.length, buffer.remaining());
    }

    @Test
    void testWrite_exceedsMaxPacketSize() {
        byte[] data = new byte[70000]; // Exceeds default max

        assertThrows(TransportException.class, () ->
            transportInstance.write(data)
        );
    }

    @Test
    void testRead_successful() throws Exception {
        byte[] data = "Test Data".getBytes();

        // Send from server
        serverChannel.send(ByteBuffer.wrap(data), clientAddress);

        // Wait a bit for packet to arrive
        Thread.sleep(50);

        // Read on client
        byte[] result = transportInstance.read(data.length);

        assertArrayEquals(data, result);
    }

    @Test
    void testRead_partialPacket() throws Exception {
        byte[] data = "0123456789".getBytes();

        // Send from server
        serverChannel.send(ByteBuffer.wrap(data), clientAddress);

        // Wait for packet
        Thread.sleep(50);

        // Read only part of the packet
        byte[] result = transportInstance.read(5);

        assertEquals(5, result.length);
        assertArrayEquals("01234".getBytes(), result);
    }

    @Test
    void testPeekReadableBytes_doesNotConsume() throws Exception {
        byte[] data = "Peek Test".getBytes();

        // Send from server
        serverChannel.send(ByteBuffer.wrap(data), clientAddress);

        // Wait for packet
        Thread.sleep(50);

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
        assertEquals(0, initialAvailable);

        // Send from server
        serverChannel.send(ByteBuffer.wrap(data), clientAddress);

        // Wait for packet
        Thread.sleep(50);

        // Check availability
        int available = transportInstance.getNumBytesAvailable();
        assertEquals(data.length, available);
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
        assertDoesNotThrow(() -> transportInstance.write(new byte[0]));
    }

    @Test
    void testWrite_nullArray() throws TransportException {
        assertDoesNotThrow(() -> transportInstance.write(null));
    }

    @Test
    void testRead_zeroBytes() throws TransportException {
        byte[] result = transportInstance.read(0);
        assertEquals(0, result.length);
    }

    @Test
    void testMultiplePackets() throws Exception {
        // Send multiple packets
        for (int i = 0; i < 5; i++) {
            byte[] data = ("Packet " + i).getBytes();
            serverChannel.send(ByteBuffer.wrap(data), clientAddress);
            Thread.sleep(10);
        }

        // Read packets
        for (int i = 0; i < 5; i++) {
            byte[] expected = ("Packet " + i).getBytes();
            byte[] received = transportInstance.read(expected.length);
            assertArrayEquals(expected, received);
        }
    }

    @Test
    void testPacketFiltering_ignoresWrongSource() throws Exception {
        // Create another sender
        DatagramChannel otherSender = DatagramChannel.open();

        try (otherSender) {
            otherSender.bind(new InetSocketAddress("localhost", 0));
            // Send a message from the wrong source
            otherSender.send(ByteBuffer.wrap("Wrong".getBytes()), clientAddress);

            // Send a message from the correct source
            serverChannel.send(ByteBuffer.wrap("Right".getBytes()), clientAddress);

            Thread.sleep(50);

            // Should only receive from the correct source
            byte[] result = transportInstance.read(5);
            assertArrayEquals("Right".getBytes(), result);
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

        // Send packet from server
        serverChannel.send(ByteBuffer.wrap("Test".getBytes()), clientAddress);

        // Wait for selector to process
        Thread.sleep(100);

        assertTrue(listenerCalled[0], "Listener should be called when packet arrives");
    }

    @Test
    void testRegisterDataListener_multipleNotifications() throws Exception {
        final int[] callCount = {0};

        // Register listener
        transportInstance.registerDataListener(() -> callCount[0]++);

        // Send multiple packets
        for (int i = 0; i < 3; i++) {
            serverChannel.send(ByteBuffer.wrap(("Packet" + i).getBytes()), clientAddress);
            Thread.sleep(50);
        }

        assertTrue(callCount[0] >= 3, "Listener should be called at least 3 times");
    }

    @Test
    void testRemoveDataListener_stopsNotifications() throws Exception {
        final int[] callCount = {0};

        // Register listener
        transportInstance.registerDataListener(() -> callCount[0]++);

        // Send packet
        serverChannel.send(ByteBuffer.wrap("Test1".getBytes()), clientAddress);
        Thread.sleep(50);

        int firstCount = callCount[0];
        assertTrue(firstCount > 0, "Listener should be called initially");

        // Remove listener
        transportInstance.removeDataListener();

        // Send more packets
        serverChannel.send(ByteBuffer.wrap("Test2".getBytes()), clientAddress);
        Thread.sleep(50);

        // Count should not increase
        assertEquals(firstCount, callCount[0], "Listener should not be called after removal");
    }

    @Test
    void testSelectorThread_startsAutomatically() throws Exception {
        // Selector thread should be running
        assertTrue(transportInstance.isOpen());

        // Should be able to receive packet immediately
        serverChannel.send(ByteBuffer.wrap("SelectorTest".getBytes()), clientAddress);
        Thread.sleep(100);

        int available = transportInstance.getNumBytesAvailable();
        assertTrue(available >= 12, "Selector should have read packet into buffer");
    }

    @Test
    void testSelectorThread_stopsOnClose() throws Exception {
        transportInstance.close();

        // Wait for selector thread to stop
        Thread.sleep(200);

        assertFalse(transportInstance.isOpen());
    }

    @Test
    void testPacketFiltering_asyncNotification() throws Exception {
        final boolean[] listenerCalled = {false};

        // Register listener
        transportInstance.registerDataListener(() -> listenerCalled[0] = true);

        // Create another sender (wrong source)
        DatagramChannel otherSender = DatagramChannel.open();

        try (otherSender) {
            otherSender.bind(new InetSocketAddress("localhost", 0));
            // Send a message from the wrong source (should be filtered)
            otherSender.send(ByteBuffer.wrap("Wrong".getBytes()), clientAddress);
            Thread.sleep(50);

            assertFalse(listenerCalled[0], "Listener should not be called for filtered packets");

            // Send a message from the correct source
            serverChannel.send(ByteBuffer.wrap("Right".getBytes()), clientAddress);
            Thread.sleep(50);

            assertTrue(listenerCalled[0], "Listener should be called for valid packets");
        }
    }

    // ========== Additional Coverage Tests ==========

    @Test
    void testGetNumBytesAvailable_whenClosed() throws TransportException {
        transportInstance.close();
        assertEquals(0, transportInstance.getNumBytesAvailable());
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
    void testPeekReadableBytes_whenClosed() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.peekReadableBytes(10)
        );
    }

    @Test
    void testPeekReadableBytes_insufficientData() throws Exception {
        // Send a small packet
        serverChannel.send(ByteBuffer.wrap("Hi".getBytes()), clientAddress);
        Thread.sleep(50);

        // Request more bytes than available
        assertThrows(TransportException.class, () ->
            transportInstance.peekReadableBytes(100)
        );
    }

    @Test
    void testRead_insufficientData() throws Exception {
        // Send a small packet
        serverChannel.send(ByteBuffer.wrap("Hi".getBytes()), clientAddress);
        Thread.sleep(50);

        // Request more bytes than available
        assertThrows(TransportException.class, () ->
            transportInstance.read(100)
        );
    }

    @Test
    void testRead_negativeBytes() throws TransportException {
        byte[] result = transportInstance.read(-1);
        assertEquals(0, result.length);
    }

    @Test
    void testRegisterDisconnectListener() {
        AtomicReference<Throwable> received = new AtomicReference<>();

        assertDoesNotThrow(() ->
            transportInstance.registerDisconnectListener(received::set)
        );
    }

    @Test
    void testRemoveDisconnectListener() {
        transportInstance.registerDisconnectListener(t -> {});

        assertDoesNotThrow(() ->
            transportInstance.removeDisconnectListener()
        );
    }

    @Test
    void testConstructor_withLocalAddress() throws Exception {
        DatagramChannel server = DatagramChannel.open();
        server.bind(new InetSocketAddress("localhost", 0));
        int port = ((InetSocketAddress) server.getLocalAddress()).getPort();

        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.localAddress = "localhost";
        config.localPort = 0;

        UdpTransportInstance instance = new UdpTransportInstance(
            new InetSocketAddress("localhost", port), config,
            new SharedUdpSocketManager(), AuditLog.builder().build());

        try {
            assertTrue(instance.isOpen());
            assertTrue(instance.getLocalAddress().isPresent());
        } finally {
            instance.close();
            server.close();
        }
    }

    @Test
    void testConstructor_withSocketOptions() throws Exception {
        DatagramChannel server = DatagramChannel.open();
        server.bind(new InetSocketAddress("localhost", 0));
        int port = ((InetSocketAddress) server.getLocalAddress()).getPort();

        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.sendBufferSize = 32768;
        config.receiveBufferSize = 32768;
        config.readTimeout = 1000;

        UdpTransportInstance instance = new UdpTransportInstance(
            new InetSocketAddress("localhost", port), config,
            new SharedUdpSocketManager(), AuditLog.builder().build());

        try {
            assertTrue(instance.isOpen());
        } finally {
            instance.close();
            server.close();
        }
    }

    @Test
    void testConstructor_invalidLocalAddress() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        // Bind to a non-local address to trigger IOException at bind()
        config.localAddress = "192.0.2.1";
        config.localPort = 12345;

        assertThrows(TransportException.class, () ->
            new UdpTransportInstance(
                new InetSocketAddress("localhost", 9999),
                config, new SharedUdpSocketManager(), AuditLog.builder().build())
        );
    }

    @Test
    void testWrite_afterChannelClosed() throws Exception {
        // Close the underlying channel to trigger IOException on write
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.write("test".getBytes())
        );
    }

    @Test
    void testGetLocalAddress() {
        assertTrue(transportInstance.getLocalAddress().isPresent());
        InetSocketAddress local = transportInstance.getLocalAddress().get();
        assertNotNull(local);
        assertTrue(local.getPort() > 0);
    }

}
