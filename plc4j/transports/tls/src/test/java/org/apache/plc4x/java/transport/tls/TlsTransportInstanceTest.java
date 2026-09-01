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
import org.apache.plc4x.java.transport.tls.config.TlsTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.*;

import javax.net.ssl.*;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.lang.reflect.Method;
import java.net.InetSocketAddress;
import java.net.ServerSocket;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyStore;
import java.security.SecureRandom;
import java.security.cert.X509Certificate;
import java.util.concurrent.*;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for TlsTransportInstance covering both success and failure paths.
 * Integration tests use a local TLS server with a self-signed certificate.
 */
class TlsTransportInstanceTest {

    private static Path testKeystorePath;
    private static final String KEYSTORE_PASSWORD = "changeit";

    private static SSLContext serverSslContext;
    private ServerSocket plainServerSocket;
    private int serverPort;
    private ExecutorService serverExecutor;

    @BeforeAll
    static void createTestKeystore() throws Exception {
        testKeystorePath = Files.createTempFile("test-tls-", ".p12");
        // keytool won't overwrite existing file
        Files.deleteIfExists(testKeystorePath);

        ProcessBuilder pb = new ProcessBuilder(
            System.getProperty("java.home") + "/bin/keytool",
            "-genkeypair", "-alias", "testserver",
            "-keyalg", "RSA", "-keysize", "2048",
            "-validity", "1",
            "-storepass", KEYSTORE_PASSWORD,
            "-keypass", KEYSTORE_PASSWORD,
            "-storetype", "PKCS12",
            "-keystore", testKeystorePath.toString(),
            "-dname", "CN=localhost,O=Test,L=Test,ST=Test,C=US",
            "-ext", "san=dns:localhost,ip:127.0.0.1"
        );
        pb.redirectErrorStream(true);
        Process p = pb.start();
        try (InputStream is = p.getInputStream()) {
            is.readAllBytes();
        }
        int exitCode = p.waitFor();
        if (exitCode != 0) {
            throw new RuntimeException("keytool failed with exit code " + exitCode);
        }

        // Create the server SSL context once (shared across all tests)
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        try (InputStream fis = Files.newInputStream(testKeystorePath)) {
            keyStore.load(fis, KEYSTORE_PASSWORD.toCharArray());
        }
        KeyManagerFactory kmf = KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm());
        kmf.init(keyStore, KEYSTORE_PASSWORD.toCharArray());
        serverSslContext = SSLContext.getInstance("TLS");
        TrustManager[] trustAll = new TrustManager[]{
            new X509ExtendedTrustManager() {
                @Override public void checkClientTrusted(X509Certificate[] chain, String authType) {}
                @Override public void checkServerTrusted(X509Certificate[] chain, String authType) {}
                @Override public X509Certificate[] getAcceptedIssuers() { return new X509Certificate[0]; }
                @Override public void checkClientTrusted(X509Certificate[] chain, String authType, java.net.Socket socket) {}
                @Override public void checkServerTrusted(X509Certificate[] chain, String authType, java.net.Socket socket) {}
                @Override public void checkClientTrusted(X509Certificate[] chain, String authType, SSLEngine engine) {}
                @Override public void checkServerTrusted(X509Certificate[] chain, String authType, SSLEngine engine) {}
            }
        };
        serverSslContext.init(kmf.getKeyManagers(), trustAll, new SecureRandom());
    }

    @AfterAll
    static void cleanupKeystore() throws Exception {
        if (testKeystorePath != null) {
            Files.deleteIfExists(testKeystorePath);
        }
    }

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

    private TlsTransportConfiguration createConfig() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        config.connectTimeout = 10000;
        config.readTimeout = 10000;
        config.receiveBufferSize = 8192;
        config.sendBufferSize = 8192;
        config.tcpNoDelay = true;
        config.keepAlive = true;
        config.verifySsl = false;
        return config;
    }

    /**
     * Starts accepting a connection on the plain server socket and wraps it
     * in a server-mode SSLSocket with the test certificate. This avoids
     * SSLServerSocket issues under surefire's Mockito agent.
     */
    private Future<SSLSocket> acceptConnection() {
        CountDownLatch ready = new CountDownLatch(1);
        Future<SSLSocket> future = serverExecutor.submit(() -> {
            ready.countDown();
            java.net.Socket rawSocket = plainServerSocket.accept();
            SSLSocket serverSsl = (SSLSocket) serverSslContext.getSocketFactory().createSocket(
                rawSocket, null, rawSocket.getPort(), true);
            serverSsl.setUseClientMode(false);
            serverSsl.startHandshake();
            return serverSsl;
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
    void testSuccessfulConnection() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            assertTrue(instance.isOpen());
            assertNotNull(instance.getRemoteAddress());
            assertNotNull(instance.getLocalAddress());
            assertNotNull(instance.getTlsProtocol());
            assertNotNull(instance.getCipherSuite());
            assertEquals(0, instance.getNumBytesAvailable());

            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);
            assertNotNull(serverSocket);
            serverSocket.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testConnectionWithExplicitTlsVersion() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportConfiguration config = createConfig();
        config.tlsVersion = "TLSv1.2";

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            config, AuditLog.builder().build());

        try {
            assertTrue(instance.isOpen());
            assertEquals("TLSv1.2", instance.getTlsProtocol());
            serverFuture.get(5, TimeUnit.SECONDS).close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testWriteData() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04};
            instance.write(testData);

            InputStream serverIn = serverSocket.getInputStream();
            byte[] received = new byte[4];
            int totalRead = 0;
            while (totalRead < 4) {
                int bytesRead = serverIn.read(received, totalRead, 4 - totalRead);
                if (bytesRead == -1) break;
                totalRead += bytesRead;
            }
            assertEquals(4, totalRead);
            assertArrayEquals(testData, received);

            serverSocket.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testReadData() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            byte[] testData = new byte[]{0x0A, 0x0B, 0x0C, 0x0D};
            OutputStream serverOut = serverSocket.getOutputStream();
            serverOut.write(testData);
            serverOut.flush();

            // Wait for data to arrive in ring buffer via reader loop
            long deadline = System.currentTimeMillis() + 3000;
            while (instance.getNumBytesAvailable() < 4 && System.currentTimeMillis() < deadline) {
                Thread.sleep(50);
            }

            assertTrue(instance.getNumBytesAvailable() >= 4,
                "Expected at least 4 bytes available but got " + instance.getNumBytesAvailable());
            byte[] read = instance.read(4);
            assertArrayEquals(testData, read);

            serverSocket.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testPeekReadableBytes() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            byte[] testData = new byte[]{0x01, 0x02, 0x03};
            OutputStream serverOut = serverSocket.getOutputStream();
            serverOut.write(testData);
            serverOut.flush();

            long deadline = System.currentTimeMillis() + 3000;
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

            serverSocket.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testCloseConnection() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

        assertTrue(instance.isOpen());
        instance.close();
        assertFalse(instance.isOpen());

        serverSocket.close();
    }

    @Test
    void testCloseAlreadyClosed() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        serverFuture.get(5, TimeUnit.SECONDS).close();

        instance.close();
        // Second close should be a no-op
        instance.close();
        assertFalse(instance.isOpen());
    }

    @Test
    void testDataListenerNotification() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            CountDownLatch latch = new CountDownLatch(1);
            instance.registerDataListener(latch::countDown);

            OutputStream serverOut = serverSocket.getOutputStream();
            serverOut.write(new byte[]{0x01});
            serverOut.flush();

            assertTrue(latch.await(5, TimeUnit.SECONDS), "Data listener should have been called");

            instance.removeDataListener();
            serverSocket.close();
        } finally {
            instance.close();
        }
    }

    /**
     * The reader loop is the only thing reading the socket, so a data listener that fails on one
     * frame must cost that frame and nothing else - the peer chooses what arrives, and a frame the
     * codec cannot parse is a frame the peer can send again.
     */
    @Test
    void testDataListenerThrowingKeepsTransportReceiving() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            CountDownLatch secondNotification = new CountDownLatch(2);
            java.util.concurrent.atomic.AtomicInteger notifications =
                new java.util.concurrent.atomic.AtomicInteger();
            instance.registerDataListener(() -> {
                secondNotification.countDown();
                if (notifications.incrementAndGet() == 1) {
                    throw new IllegalStateException("simulated parse failure");
                }
            });

            OutputStream serverOut = serverSocket.getOutputStream();
            serverOut.write(new byte[]{0x01});
            serverOut.flush();
            Thread.sleep(200);
            serverOut.write(new byte[]{0x02});
            serverOut.flush();

            assertTrue(secondNotification.await(5, TimeUnit.SECONDS),
                "the reader loop must survive a listener that threw and deliver the next frame");
            assertTrue(instance.isOpen(), "one unparseable frame must not close the transport");

            instance.removeDataListener();
            serverSocket.close();
        } finally {
            instance.close();
        }
    }

    /**
     * If the loop ever does stop while the transport still believes it is open, whoever is waiting
     * on a response has to be told - otherwise they wait on a connection that can no longer
     * receive.
     */
    @Test
    void testReaderLoopExitingWhileOpenReportsDisconnect() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            CountDownLatch disconnected = new CountDownLatch(1);
            instance.registerDisconnectListener(cause -> disconnected.countDown());

            // The server going away ends the loop; the transport must not stay "open".
            serverSocket.close();

            assertTrue(disconnected.await(5, TimeUnit.SECONDS),
                "a reader loop that stops must report the disconnection");
            assertFalse(instance.isOpen());
        } finally {
            instance.close();
        }
    }

    /**
     * More bytes can arrive than the ring buffer has room for, and the codec is the only thing that
     * knows where a frame ends - so a byte read off the socket and then dropped takes the middle out
     * of somebody's message. Every byte the peer sent has to turn up, in the order it was sent.
     */
    @Test
    void testReadingMoreThanTheRingHoldsLosesNothing() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportConfiguration config = createConfig();
        // Small enough that a single send cannot fit, so the loop has to wait for room.
        config.receiveBufferSize = 2048;

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort), config, AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            int total = 16384;
            byte[] sent = new byte[total];
            for (int i = 0; i < total; i++) {
                sent[i] = (byte) (i % 251);
            }

            java.io.ByteArrayOutputStream drained = new java.io.ByteArrayOutputStream();
            instance.registerDataListener(() -> { });

            OutputStream serverOut = serverSocket.getOutputStream();
            Thread sender = new Thread(() -> {
                try {
                    serverOut.write(sent);
                    serverOut.flush();
                } catch (IOException e) {
                    // the assertions below report what actually arrived
                }
            });
            sender.start();

            long deadline = System.currentTimeMillis() + 15000;
            while (drained.size() < total && System.currentTimeMillis() < deadline) {
                int available = instance.getNumBytesAvailable();
                if (available > 0) {
                    drained.write(instance.read(available));
                } else {
                    Thread.sleep(5);
                }
            }
            sender.join(5000);

            assertEquals(total, drained.size(),
                "every byte the peer sent must be delivered, not dropped when the ring filled");
            assertArrayEquals(sent, drained.toByteArray(),
                "the bytes must arrive in the order they were sent");

            instance.removeDataListener();
            serverSocket.close();
        } finally {
            instance.close();
        }
    }

    @Test
    void testDisconnectListenerOnServerClose() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            CountDownLatch latch = new CountDownLatch(1);
            instance.registerDisconnectListener(cause -> latch.countDown());

            // Close server side to trigger disconnect in reader loop
            serverSocket.close();

            assertTrue(latch.await(5, TimeUnit.SECONDS), "Disconnect listener should have been called");
        } finally {
            instance.close();
        }
    }

    @Test
    void testDisconnectListenerThatThrowsException() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            SSLSocket serverSocket = serverFuture.get(5, TimeUnit.SECONDS);

            CountDownLatch latch = new CountDownLatch(1);
            instance.registerDisconnectListener(cause -> {
                latch.countDown();
                throw new RuntimeException("listener error");
            });

            serverSocket.close();

            // Listener exception should be caught, not propagated
            assertTrue(latch.await(5, TimeUnit.SECONDS), "Disconnect listener should have been called despite throwing");
        } finally {
            instance.close();
        }
    }

    @Test
    void testRegisterAndRemoveListeners() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS).close();

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
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS);
            // Should not throw — null and empty are no-ops
            instance.write(null);
            instance.write(new byte[0]);
        } finally {
            instance.close();
        }
    }

    @Test
    void testReadZeroBytes() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS);
            byte[] result = instance.read(0);
            assertNotNull(result);
            assertEquals(0, result.length);
        } finally {
            instance.close();
        }
    }

    @Test
    void testPeekZeroBytes() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS);
            byte[] result = instance.peekReadableBytes(0);
            assertNotNull(result);
            assertEquals(0, result.length);
        } finally {
            instance.close();
        }
    }

    @Test
    void testReadNegativeBytes() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS);
            byte[] result = instance.read(-1);
            assertNotNull(result);
            assertEquals(0, result.length);
        } finally {
            instance.close();
        }
    }

    @Test
    void testPeekNegativeBytes() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS);
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
        TlsTransportConfiguration config = createConfig();
        config.connectTimeout = 1000;

        assertThrows(TransportException.class, () ->
            new TlsTransportInstance(new InetSocketAddress("127.0.0.1", 1), config, AuditLog.builder().build())
        );
    }

    @Test
    void testConnectionToInvalidHost() {
        TlsTransportConfiguration config = createConfig();
        config.connectTimeout = 1000;

        assertThrows(TransportException.class, () ->
            new TlsTransportInstance(
                new InetSocketAddress("invalid.host.does.not.exist.local", 8016),
                config, AuditLog.builder().build())
        );
    }

    @Test
    void testConnectionErrorMessageContainsHost() {
        TlsTransportConfiguration config = createConfig();
        config.connectTimeout = 1000;

        TransportException exception = assertThrows(TransportException.class, () ->
            new TlsTransportInstance(new InetSocketAddress("127.0.0.1", 1), config, AuditLog.builder().build())
        );

        assertNotNull(exception.getMessage());
        assertTrue(exception.getMessage().contains("127.0.0.1") ||
                   exception.getMessage().contains("Connection refused") ||
                   exception.getMessage().contains("connect"),
            "Error message should contain connection info: " + exception.getMessage());
    }

    @Test
    void testCertificateVerificationFailure() {
        // Connect with verify=true to our self-signed server — should fail with PKIX error
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportConfiguration config = createConfig();
        config.verifySsl = true;
        config.connectTimeout = 5000;

        TransportException exception = assertThrows(TransportException.class, () ->
            new TlsTransportInstance(
                new InetSocketAddress("127.0.0.1", serverPort),
                config, AuditLog.builder().build())
        );

        // formatTlsError should produce a message mentioning the host
        assertNotNull(exception.getMessage());
        assertTrue(exception.getMessage().contains("127.0.0.1") ||
                   exception.getMessage().contains("localhost"),
            "Error message should contain host: " + exception.getMessage());
    }

    @Test
    void testHandshakeErrorWithNonTlsServer() throws Exception {
        // Connect to a plain TCP server to trigger TLS handshake failure
        try (ServerSocket plainServer = new ServerSocket()) {
            plainServer.bind(new InetSocketAddress("127.0.0.1", 0));
            int plainPort = plainServer.getLocalPort();

            serverExecutor.submit(() -> {
                try {
                    java.net.Socket socket = plainServer.accept();
                    // Send non-TLS data to cause handshake failure
                    socket.getOutputStream().write("NOT-TLS\r\n".getBytes());
                    socket.getOutputStream().flush();
                    Thread.sleep(2000);
                    socket.close();
                } catch (Exception ignored) {}
                return null;
            });

            TlsTransportConfiguration config = createConfig();
            config.connectTimeout = 5000;
            config.readTimeout = 3000;

            TransportException exception = assertThrows(TransportException.class, () ->
                new TlsTransportInstance(
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
    void testLocalAddressBinding() {
        TlsTransportConfiguration config = createConfig();
        config.connectTimeout = 1000;
        config.localAddress = "127.0.0.1";
        config.localPort = 0;

        assertThrows(TransportException.class, () ->
            new TlsTransportInstance(new InetSocketAddress("127.0.0.1", 1), config, AuditLog.builder().build())
        );
    }

    @Test
    void testReadMoreThanAvailable() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS);
            assertThrows(TransportException.class, () -> instance.read(100));
        } finally {
            instance.close();
        }
    }

    @Test
    void testPeekMoreThanAvailable() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS);
            assertThrows(TransportException.class, () -> instance.peekReadableBytes(100));
        } finally {
            instance.close();
        }
    }

    @Test
    void testWriteAfterClose() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        serverFuture.get(5, TimeUnit.SECONDS).close();
        instance.close();

        assertThrows(TransportException.class, () -> instance.write(new byte[]{0x01}));
    }

    @Test
    void testReadAfterClose() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        serverFuture.get(5, TimeUnit.SECONDS).close();
        instance.close();

        assertThrows(TransportException.class, () -> instance.read(1));
    }

    @Test
    void testPeekAfterClose() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        serverFuture.get(5, TimeUnit.SECONDS).close();
        instance.close();

        assertThrows(TransportException.class, () -> instance.peekReadableBytes(1));
    }

    @Test
    void testGetNumBytesAvailableAfterClose() throws Exception {
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        serverFuture.get(5, TimeUnit.SECONDS).close();
        instance.close();

        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void testSocketOptions() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 16384;
        config.sendBufferSize = 16384;
        config.tcpNoDelay = true;
        config.keepAlive = true;
        config.readTimeout = 1000;
        config.verifySsl = false;

        // Will fail to connect but exercises socket option setting code
        assertThrows(TransportException.class, () ->
            new TlsTransportInstance(new InetSocketAddress("127.0.0.1", 1), config, AuditLog.builder().build())
        );
    }

    // ========== Reflection Tests for Private Utility Methods ==========

    @Test
    void testBytesToHex() throws Exception {
        Method method = TlsTransportInstance.class.getDeclaredMethod("bytesToHex", byte[].class);
        method.setAccessible(true);

        assertEquals("(empty)", method.invoke(null, (Object) null));
        assertEquals("(empty)", method.invoke(null, (Object) new byte[0]));
        assertEquals("0102ff", method.invoke(null, (Object) new byte[]{0x01, 0x02, (byte) 0xFF}));
        assertEquals("00", method.invoke(null, (Object) new byte[]{0x00}));
        assertEquals("7f80", method.invoke(null, (Object) new byte[]{0x7F, (byte) 0x80}));
    }

    @Test
    void testLogTlsErrorChain() throws Exception {
        Method method = TlsTransportInstance.class.getDeclaredMethod(
            "logTlsErrorChain", AuditLog.class, Throwable.class);
        method.setAccessible(true);

        AuditLog auditLog = AuditLog.builder().build();

        // Chain of 3 causes
        Exception inner = new RuntimeException("root cause");
        Exception middle = new IOException("middle", inner);
        Exception outer = new IOException("outer error", middle);

        // Should not throw
        method.invoke(null, auditLog, outer);

        // Single exception (no chain)
        method.invoke(null, auditLog, new IOException("single error"));
    }

    @Test
    void testFormatTlsError_allBranches() throws Exception {
        // Create a successful instance to call the private method on
        Future<SSLSocket> serverFuture = acceptConnection();

        TlsTransportInstance instance = new TlsTransportInstance(
            new InetSocketAddress("127.0.0.1", serverPort),
            createConfig(), AuditLog.builder().build());

        try {
            serverFuture.get(5, TimeUnit.SECONDS);

            Method method = TlsTransportInstance.class.getDeclaredMethod(
                "formatTlsError", IOException.class, InetSocketAddress.class);
            method.setAccessible(true);

            InetSocketAddress addr = new InetSocketAddress("testhost", 1234);

            // PKIX error branch
            String pkixResult = (String) method.invoke(instance,
                new IOException("PKIX path building failed"), addr);
            assertTrue(pkixResult.contains("certificate not trusted"));
            assertTrue(pkixResult.contains("testhost"));
            assertTrue(pkixResult.contains("1234"));

            // Certificate error branch
            String certResult = (String) method.invoke(instance,
                new IOException("No subject alternative certificate names"), addr);
            assertTrue(certResult.contains("certificate not trusted"));

            // Handshake error branch
            String handshakeResult = (String) method.invoke(instance,
                new IOException("handshake_failure"), addr);
            assertTrue(handshakeResult.contains("handshake failed"));
            assertTrue(handshakeResult.contains("testhost"));

            // Handshake with capital H
            String handshakeResult2 = (String) method.invoke(instance,
                new IOException("Handshake alert: something"), addr);
            assertTrue(handshakeResult2.contains("handshake failed"));

            // Connection refused branch
            String refusedResult = (String) method.invoke(instance,
                new IOException("Connection refused"), addr);
            assertTrue(refusedResult.contains("Connection refused"));
            assertTrue(refusedResult.contains("testhost"));

            // Generic error branch
            String genericResult = (String) method.invoke(instance,
                new IOException("some unexpected error"), addr);
            assertTrue(genericResult.contains("testhost"));
            assertTrue(genericResult.contains("1234"));
            assertTrue(genericResult.contains("some unexpected error"));

            // Null message branch
            String nullResult = (String) method.invoke(instance,
                new IOException((String) null), addr);
            assertTrue(nullResult.contains("testhost"));
            assertTrue(nullResult.contains("IOException")); // Falls back to class name
        } finally {
            instance.close();
        }
    }

    @Test
    void testFormatSubjectAltNames_withoutSans() throws Exception {
        Method method = TlsTransportInstance.class.getDeclaredMethod(
            "formatSubjectAltNames", X509Certificate.class);
        method.setAccessible(true);

        // Create a mock-like cert that returns null SANs using a dynamic proxy
        // We can test with a cert that has no SANs by loading one from a keystore
        // But our test keystore has SANs. Generate a second keystore without SANs.
        Path noSansKeystorePath = Files.createTempFile("test-tls-nosans-", ".p12");
        Files.deleteIfExists(noSansKeystorePath);
        try {
            ProcessBuilder pb = new ProcessBuilder(
                System.getProperty("java.home") + "/bin/keytool",
                "-genkeypair", "-alias", "nosans",
                "-keyalg", "RSA", "-keysize", "2048",
                "-validity", "1",
                "-storepass", KEYSTORE_PASSWORD,
                "-keypass", KEYSTORE_PASSWORD,
                "-storetype", "PKCS12",
                "-keystore", noSansKeystorePath.toString(),
                "-dname", "CN=localhost"
                // No -ext san=... so no SANs
            );
            pb.redirectErrorStream(true);
            Process p = pb.start();
            try (InputStream is = p.getInputStream()) { is.readAllBytes(); }
            p.waitFor();

            KeyStore ks = KeyStore.getInstance("PKCS12");
            try (InputStream fis = Files.newInputStream(noSansKeystorePath)) {
                ks.load(fis, KEYSTORE_PASSWORD.toCharArray());
            }
            X509Certificate cert = (X509Certificate) ks.getCertificate("nosans");

            String result = (String) method.invoke(null, cert);
            assertEquals("", result, "Should return empty string for cert without SANs");
        } finally {
            Files.deleteIfExists(noSansKeystorePath);
        }
    }

    @Test
    void testFormatSubjectAltNames_withSans() throws Exception {
        Method method = TlsTransportInstance.class.getDeclaredMethod(
            "formatSubjectAltNames", X509Certificate.class);
        method.setAccessible(true);

        // Load our test keystore which has SANs (dns:localhost, ip:127.0.0.1)
        KeyStore ks = KeyStore.getInstance("PKCS12");
        try (InputStream fis = Files.newInputStream(testKeystorePath)) {
            ks.load(fis, KEYSTORE_PASSWORD.toCharArray());
        }
        X509Certificate cert = (X509Certificate) ks.getCertificate("testserver");

        String result = (String) method.invoke(null, cert);
        assertNotNull(result);
        assertTrue(result.contains("SANs="), "Should contain SANs section: " + result);
        assertTrue(result.contains("dNSName=localhost"), "Should contain DNS SAN: " + result);
        assertTrue(result.contains("IP="), "Should contain IP SAN: " + result);
    }
}
