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
package org.apache.plc4x.java.transport.tls;

import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.tls.config.PskTlsTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.bouncycastle.tls.*;
import org.bouncycastle.tls.crypto.impl.bc.BcTlsCrypto;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.ServerSocket;
import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.util.concurrent.*;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for PskTlsTransportInstance covering integration and error paths.
 * Uses a local Bouncy Castle PSK TLS server for integration tests.
 */
class PskTlsTransportInstanceTest {

    private static final String TEST_PSK_IDENTITY = "TestPskUser";
    private static final String TEST_PSK_KEY_HEX = "0123456789abcdef0123456789abcdef";
    private static final byte[] TEST_PSK_KEY_BYTES = hexToBytes(TEST_PSK_KEY_HEX);

    private ServerSocket plainServerSocket;
    private int serverPort;
    private ExecutorService serverExecutor;

    @BeforeEach
    void setUp() throws Exception {
        serverExecutor = Executors.newCachedThreadPool();
        plainServerSocket = new ServerSocket();
        plainServerSocket.bind(new InetSocketAddress("127.0.0.1", 0));
        plainServerSocket.setSoTimeout(15000);
        serverPort = plainServerSocket.getLocalPort();
    }

    @AfterEach
    void tearDown() throws Exception {
        if (plainServerSocket != null && !plainServerSocket.isClosed()) {
            plainServerSocket.close();
        }
        if (serverExecutor != null) {
            serverExecutor.shutdownNow();
            serverExecutor.awaitTermination(5, TimeUnit.SECONDS);
        }
    }

    private PskTlsTransportConfiguration createPskConfig() {
        PskTlsTransportConfiguration config = new PskTlsTransportConfiguration();
        config.connectTimeout = 10000;
        config.readTimeout = 10000;
        config.receiveBufferSize = 8192;
        config.sendBufferSize = 8192;
        config.tcpNoDelay = true;
        config.keepAlive = true;
        config.pskIdentity = TEST_PSK_IDENTITY;
        config.pskKey = TEST_PSK_KEY_HEX;
        return config;
    }

    /**
     * Starts a Bouncy Castle PSK TLS server on the plain server socket.
     * Returns the server-side TlsServerProtocol for reading/writing.
     */
    private Future<TlsServerProtocol> acceptPskConnection() {
        CountDownLatch ready = new CountDownLatch(1);
        Future<TlsServerProtocol> future = serverExecutor.submit(() -> {
            ready.countDown();
            java.net.Socket rawSocket = plainServerSocket.accept();
            rawSocket.setSoTimeout(10000);

            BcTlsCrypto crypto = new BcTlsCrypto(new SecureRandom());
            TlsServerProtocol serverProtocol = new TlsServerProtocol(
                rawSocket.getInputStream(), rawSocket.getOutputStream());

            PSKTlsServer pskServer = new PSKTlsServer(crypto, new TlsPSKIdentityManager() {
                @Override
                public byte[] getHint() {
                    return null;
                }

                @Override
                public byte[] getPSK(byte[] identity) {
                    String receivedIdentity = new String(identity, StandardCharsets.UTF_8);
                    if (TEST_PSK_IDENTITY.equals(receivedIdentity)) {
                        return TEST_PSK_KEY_BYTES.clone();
                    }
                    return null;
                }
            }) {
                @Override
                protected int[] getSupportedCipherSuites() {
                    return new int[]{
                        CipherSuite.TLS_PSK_WITH_AES_256_GCM_SHA384,
                        CipherSuite.TLS_PSK_WITH_AES_128_GCM_SHA256,
                        CipherSuite.TLS_PSK_WITH_AES_256_CBC_SHA384,
                        CipherSuite.TLS_PSK_WITH_AES_128_CBC_SHA256
                    };
                }
            };

            serverProtocol.accept(pskServer);
            return serverProtocol;
        });

        try {
            ready.await(5, TimeUnit.SECONDS);
            Thread.sleep(50);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        return future;
    }

    // ========== Integration Tests (Success Paths) ==========

    @Test
    void testSuccessfulPskConnection() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            assertTrue(instance.isOpen());
            assertNotNull(instance.getRemoteAddress());
            assertEquals(0, instance.getNumBytesAvailable());

            TlsServerProtocol server = serverFuture.get(10, TimeUnit.SECONDS);
            assertNotNull(server);
            server.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testWriteData() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            TlsServerProtocol server = serverFuture.get(10, TimeUnit.SECONDS);

            byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04};
            instance.write(testData);

            InputStream serverIn = server.getInputStream();
            byte[] received = new byte[4];
            int totalRead = 0;
            while (totalRead < 4) {
                int bytesRead = serverIn.read(received, totalRead, 4 - totalRead);
                if (bytesRead == -1) break;
                totalRead += bytesRead;
            }
            assertEquals(4, totalRead);
            assertArrayEquals(testData, received);

            server.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testReadData() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            TlsServerProtocol server = serverFuture.get(10, TimeUnit.SECONDS);

            byte[] testData = new byte[]{0x0A, 0x0B, 0x0C, 0x0D};
            OutputStream serverOut = server.getOutputStream();
            serverOut.write(testData);
            serverOut.flush();

            // Wait for data to arrive
            long deadline = System.currentTimeMillis() + 5000;
            while (instance.getNumBytesAvailable() < 4 && System.currentTimeMillis() < deadline) {
                Thread.sleep(50);
            }

            assertTrue(instance.getNumBytesAvailable() >= 4,
                "Expected at least 4 bytes available but got " + instance.getNumBytesAvailable());
            byte[] read = instance.read(4);
            assertArrayEquals(testData, read);

            server.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testPeekReadableBytes() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            TlsServerProtocol server = serverFuture.get(10, TimeUnit.SECONDS);

            byte[] testData = new byte[]{0x01, 0x02, 0x03};
            OutputStream serverOut = server.getOutputStream();
            serverOut.write(testData);
            serverOut.flush();

            long deadline = System.currentTimeMillis() + 5000;
            while (instance.getNumBytesAvailable() < 3 && System.currentTimeMillis() < deadline) {
                Thread.sleep(50);
            }

            // Peek should not consume data
            byte[] peeked = instance.peekReadableBytes(3);
            assertArrayEquals(testData, peeked);
            assertEquals(3, instance.getNumBytesAvailable());

            // Read should consume data
            byte[] read = instance.read(3);
            assertArrayEquals(testData, read);
            assertEquals(0, instance.getNumBytesAvailable());

            server.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testCloseConnection() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        TlsServerProtocol server = serverFuture.get(10, TimeUnit.SECONDS);

        assertTrue(instance.isOpen());
        instance.close();
        assertFalse(instance.isOpen());

        server.close();
    }

    @Test
    void testCloseAlreadyClosed() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        serverFuture.get(10, TimeUnit.SECONDS).close();

        instance.close();
        // Second close should be a no-op
        instance.close();
        assertFalse(instance.isOpen());
    }

    @Test
    void testDataListenerNotification() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            TlsServerProtocol server = serverFuture.get(10, TimeUnit.SECONDS);

            CountDownLatch latch = new CountDownLatch(1);
            instance.registerDataListener(latch::countDown);

            OutputStream serverOut = server.getOutputStream();
            serverOut.write(new byte[]{0x01});
            serverOut.flush();

            assertTrue(latch.await(5, TimeUnit.SECONDS), "Data listener should have been called");

            instance.removeDataListener();
            server.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testDisconnectListenerOnServerClose() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            TlsServerProtocol server = serverFuture.get(10, TimeUnit.SECONDS);

            CountDownLatch latch = new CountDownLatch(1);
            instance.registerDisconnectListener(cause -> latch.countDown());

            // Close server side to trigger disconnect
            server.close();

            assertTrue(latch.await(5, TimeUnit.SECONDS), "Disconnect listener should have been called");
        } finally {
            instance.close();
        }
    }

    @Test
    void testDisconnectListenerThatThrowsException() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            TlsServerProtocol server = serverFuture.get(10, TimeUnit.SECONDS);

            CountDownLatch latch = new CountDownLatch(1);
            instance.registerDisconnectListener(cause -> {
                latch.countDown();
                throw new RuntimeException("listener error");
            });

            server.close();

            assertTrue(latch.await(5, TimeUnit.SECONDS), "Disconnect listener should have been called despite throwing");
        } finally {
            instance.close();
        }
    }

    @Test
    void testRegisterAndRemoveListeners() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS).close();

            instance.registerDataListener(() -> {});
            instance.removeDataListener();
            instance.registerDisconnectListener(t -> {});
            instance.removeDisconnectListener();
        } finally {
            instance.close();
        }
    }

    // ========== Boundary / Guard Tests ==========

    @Test
    void testWriteNullBytes() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);
            instance.write(null);
            instance.write(new byte[0]);
        } finally {
            instance.close();
        }
    }

    @Test
    void testReadZeroBytes() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);
            byte[] result = instance.read(0);
            assertNotNull(result);
            assertEquals(0, result.length);
        } finally {
            instance.close();
        }
    }

    @Test
    void testPeekZeroBytes() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);
            byte[] result = instance.peekReadableBytes(0);
            assertNotNull(result);
            assertEquals(0, result.length);
        } finally {
            instance.close();
        }
    }

    @Test
    void testReadNegativeBytes() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);
            byte[] result = instance.read(-1);
            assertNotNull(result);
            assertEquals(0, result.length);
        } finally {
            instance.close();
        }
    }

    @Test
    void testPeekNegativeBytes() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);
            byte[] result = instance.peekReadableBytes(-1);
            assertNotNull(result);
            assertEquals(0, result.length);
        } finally {
            instance.close();
        }
    }

    // ========== Error Path Tests ==========

    @Test
    void testConnectionRefused() {
        PskTlsTransportConfiguration config = createPskConfig();
        config.connectTimeout = 1000;

        assertThrows(TransportException.class, () ->
            new PskTlsTransportInstance(new InetSocketAddress("127.0.0.1", 1), config, AuditLog.builder().build())
        );
    }

    @Test
    void testConnectionToInvalidHost() {
        PskTlsTransportConfiguration config = createPskConfig();
        config.connectTimeout = 1000;

        assertThrows(TransportException.class, () ->
            new PskTlsTransportInstance(
                new InetSocketAddress("invalid.host.does.not.exist.local", 8016),
                config, AuditLog.builder().build())
        );
    }

    @Test
    void testConnectionRefusedErrorMessage() {
        PskTlsTransportConfiguration config = createPskConfig();
        config.connectTimeout = 1000;

        TransportException exception = assertThrows(TransportException.class, () ->
            new PskTlsTransportInstance(new InetSocketAddress("127.0.0.1", 1), config, AuditLog.builder().build())
        );

        assertNotNull(exception.getMessage());
        assertTrue(exception.getMessage().contains("127.0.0.1") ||
                   exception.getMessage().contains("localhost") ||
                   exception.getMessage().contains("Connection refused"),
            "Error message should contain connection info: " + exception.getMessage());
    }

    @Test
    void testHandshakeWithNonPskServer() throws Exception {
        // Connect to a plain TCP server to trigger PSK handshake failure
        try (ServerSocket plainServer = new ServerSocket()) {
            plainServer.bind(new InetSocketAddress("127.0.0.1", 0));
            int plainPort = plainServer.getLocalPort();

            serverExecutor.submit(() -> {
                try {
                    java.net.Socket socket = plainServer.accept();
                    socket.getOutputStream().write("NOT-TLS\r\n".getBytes());
                    socket.getOutputStream().flush();
                    Thread.sleep(2000);
                    socket.close();
                } catch (Exception ignored) {}
                return null;
            });

            PskTlsTransportConfiguration config = createPskConfig();
            config.connectTimeout = 5000;
            config.readTimeout = 3000;

            TransportException exception = assertThrows(TransportException.class, () ->
                new PskTlsTransportInstance(
                    new InetSocketAddress("127.0.0.1", plainPort),
                    config, AuditLog.builder().build())
            );

            assertNotNull(exception.getMessage());
            assertTrue(exception.getMessage().contains("127.0.0.1") ||
                       exception.getMessage().contains("localhost"),
                "Error message should contain host: " + exception.getMessage());
        }
    }

    @Test
    void testWrongPskKey() throws Exception {
        // Server expects TEST_PSK_KEY but client provides wrong key
        CountDownLatch ready = new CountDownLatch(1);
        serverExecutor.submit(() -> {
            ready.countDown();
            try {
                java.net.Socket rawSocket = plainServerSocket.accept();
                rawSocket.setSoTimeout(10000);
                BcTlsCrypto crypto = new BcTlsCrypto(new SecureRandom());
                TlsServerProtocol serverProtocol = new TlsServerProtocol(
                    rawSocket.getInputStream(), rawSocket.getOutputStream());

                PSKTlsServer pskServer = new PSKTlsServer(crypto, new TlsPSKIdentityManager() {
                    @Override
                    public byte[] getHint() { return null; }

                    @Override
                    public byte[] getPSK(byte[] identity) {
                        return TEST_PSK_KEY_BYTES.clone();
                    }
                }) {
                    @Override
                    protected int[] getSupportedCipherSuites() {
                        return new int[]{CipherSuite.TLS_PSK_WITH_AES_128_GCM_SHA256};
                    }
                };

                serverProtocol.accept(pskServer);
            } catch (IOException ignored) {
                // Expected - handshake should fail
            }
            return null;
        });
        ready.await(5, TimeUnit.SECONDS);
        Thread.sleep(50);

        PskTlsTransportConfiguration config = createPskConfig();
        config.pskKey = "ffffffffffffffffffffffffffffffff"; // Wrong key

        assertThrows(TransportException.class, () ->
            new PskTlsTransportInstance(
                new InetSocketAddress("127.0.0.1", serverPort),
                config, AuditLog.builder().build())
        );
    }

    @Test
    void testLocalAddressBinding() {
        PskTlsTransportConfiguration config = createPskConfig();
        config.connectTimeout = 1000;
        config.localAddress = "127.0.0.1";
        config.localPort = 0;

        assertThrows(TransportException.class, () ->
            new PskTlsTransportInstance(new InetSocketAddress("127.0.0.1", 1), config, AuditLog.builder().build())
        );
    }

    @Test
    void testReadMoreThanAvailable() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);
            assertThrows(TransportException.class, () -> instance.read(100));
        } finally {
            instance.close();
        }
    }

    @Test
    void testPeekMoreThanAvailable() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);
            assertThrows(TransportException.class, () -> instance.peekReadableBytes(100));
        } finally {
            instance.close();
        }
    }

    @Test
    void testWriteAfterClose() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        serverFuture.get(10, TimeUnit.SECONDS).close();
        instance.close();

        assertThrows(TransportException.class, () -> instance.write(new byte[]{0x01}));
    }

    @Test
    void testReadAfterClose() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        serverFuture.get(10, TimeUnit.SECONDS).close();
        instance.close();

        assertThrows(TransportException.class, () -> instance.read(1));
    }

    @Test
    void testPeekAfterClose() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        serverFuture.get(10, TimeUnit.SECONDS).close();
        instance.close();

        assertThrows(TransportException.class, () -> instance.peekReadableBytes(1));
    }

    @Test
    void testGetNumBytesAvailableAfterClose() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        serverFuture.get(10, TimeUnit.SECONDS).close();
        instance.close();

        assertEquals(0, instance.getNumBytesAvailable());
    }

    // ========== Static / Package-Private Utility Method Tests ==========

    @Test
    void testResolveCipherSuiteName_knownSuites() {
        assertEquals("TLS_PSK_WITH_AES_256_GCM_SHA384",
            PskTlsTransportInstance.resolveCipherSuiteName(CipherSuite.TLS_PSK_WITH_AES_256_GCM_SHA384));
        assertEquals("TLS_PSK_WITH_AES_128_GCM_SHA256",
            PskTlsTransportInstance.resolveCipherSuiteName(CipherSuite.TLS_PSK_WITH_AES_128_GCM_SHA256));
        assertEquals("TLS_PSK_WITH_AES_256_CBC_SHA384",
            PskTlsTransportInstance.resolveCipherSuiteName(CipherSuite.TLS_PSK_WITH_AES_256_CBC_SHA384));
        assertEquals("TLS_PSK_WITH_AES_128_CBC_SHA256",
            PskTlsTransportInstance.resolveCipherSuiteName(CipherSuite.TLS_PSK_WITH_AES_128_CBC_SHA256));
    }

    @Test
    void testResolveCipherSuiteName_unknownSuite() {
        String result = PskTlsTransportInstance.resolveCipherSuiteName(0x9999);
        assertTrue(result.startsWith("0x"), "Unknown suite should be hex: " + result);
    }

    @Test
    void testLogPskErrorChain_singleException() {
        AuditLog auditLog = AuditLog.builder().build();
        // Should not throw
        PskTlsTransportInstance.logPskErrorChain(auditLog, new IOException("single error"));
    }

    @Test
    void testLogPskErrorChain_chainedExceptions() {
        AuditLog auditLog = AuditLog.builder().build();
        Exception inner = new RuntimeException("root cause");
        Exception middle = new IOException("middle", inner);
        Exception outer = new IOException("outer error", middle);

        // Should not throw
        PskTlsTransportInstance.logPskErrorChain(auditLog, outer);
    }

    @Test
    void testFormatPskError_decryptError() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);

            InetSocketAddress addr = new InetSocketAddress("testhost", 1234);

            String result = instance.formatPskError(
                new TlsFatalAlert(AlertDescription.decrypt_error), addr);
            assertTrue(result.contains("PSK key is likely incorrect"), "Should mention key: " + result);
            assertTrue(result.contains("testhost"), "Should contain host: " + result);
        } finally {
            instance.close();
        }
    }

    @Test
    void testFormatPskError_unknownPskIdentity() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);

            InetSocketAddress addr = new InetSocketAddress("testhost", 1234);

            String result = instance.formatPskError(
                new TlsFatalAlert(AlertDescription.unknown_psk_identity), addr);
            assertTrue(result.contains("identity not recognized"), "Should mention identity: " + result);
            assertTrue(result.contains("testhost"), "Should contain host: " + result);
        } finally {
            instance.close();
        }
    }

    @Test
    void testFormatPskError_handshakeFailure() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);

            InetSocketAddress addr = new InetSocketAddress("testhost", 1234);

            String result = instance.formatPskError(
                new TlsFatalAlert(AlertDescription.handshake_failure), addr);
            assertTrue(result.contains("handshake failed"), "Should mention handshake: " + result);
            assertTrue(result.contains("testhost"), "Should contain host: " + result);
        } finally {
            instance.close();
        }
    }

    @Test
    void testFormatPskError_connectionRefused() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);

            InetSocketAddress addr = new InetSocketAddress("testhost", 1234);

            String result = instance.formatPskError(
                new IOException("Connection refused"), addr);
            assertTrue(result.contains("Connection refused"), "Should mention refused: " + result);
            assertTrue(result.contains("testhost"), "Should contain host: " + result);
        } finally {
            instance.close();
        }
    }

    @Test
    void testFormatPskError_genericError() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);

            InetSocketAddress addr = new InetSocketAddress("testhost", 1234);

            String result = instance.formatPskError(
                new IOException("some unexpected error"), addr);
            assertTrue(result.contains("testhost"), "Should contain host: " + result);
            assertTrue(result.contains("1234"), "Should contain port: " + result);
            assertTrue(result.contains("some unexpected error"), "Should contain error: " + result);
        } finally {
            instance.close();
        }
    }

    @Test
    void testFormatPskError_nullMessage() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);

            InetSocketAddress addr = new InetSocketAddress("testhost", 1234);

            String result = instance.formatPskError(
                new IOException((String) null), addr);
            assertTrue(result.contains("testhost"), "Should contain host: " + result);
            assertTrue(result.contains("IOException"), "Should fall back to class name: " + result);
        } finally {
            instance.close();
        }
    }

    @Test
    void testFormatPskError_tlsFatalAlertOtherDescription() throws Exception {
        Future<TlsServerProtocol> serverFuture = acceptPskConnection();

        PskTlsTransportInstance instance = new PskTlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createPskConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(10, TimeUnit.SECONDS);

            InetSocketAddress addr = new InetSocketAddress("testhost", 1234);

            // TlsFatalAlert with a different description (not decrypt/unknown_psk/handshake)
            String result = instance.formatPskError(
                new TlsFatalAlert(AlertDescription.internal_error), addr);
            assertTrue(result.contains("testhost"), "Should contain host: " + result);
        } finally {
            instance.close();
        }
    }

    // ========== Helper ==========

    private static byte[] hexToBytes(String hex) {
        byte[] bytes = new byte[hex.length() / 2];
        for (int i = 0; i < bytes.length; i++) {
            bytes[i] = (byte) Integer.parseInt(hex.substring(i * 2, i * 2 + 2), 16);
        }
        return bytes;
    }
}
