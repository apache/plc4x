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
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import javax.net.ssl.KeyManagerFactory;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLServerSocket;
import javax.net.ssl.SSLSocket;
import java.io.InputStream;
import java.net.InetSocketAddress;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyStore;
import java.security.SecureRandom;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Chain validation says a certificate was issued by someone we trust. It says nothing about who it
 * was issued <em>to</em> - so on its own, a certificate trusted for any host is accepted for every
 * host, which is exactly what a machine in the middle needs. These connect to 127.0.0.1 with a
 * certificate issued for somewhere else.
 */
class TlsHostnameVerificationTest {

    private static final String PASSWORD = "changeit";

    private static Path wrongNameKeystore;
    private static SSLContext serverContext;

    private SSLServerSocket serverSocket;
    private ExecutorService serverExecutor;
    private int port;

    @BeforeAll
    static void createCertificateForSomewhereElse() throws Exception {
        wrongNameKeystore = Files.createTempFile("tls-wrong-name-", ".p12");
        Files.deleteIfExists(wrongNameKeystore);

        // Issued for a name this test will never connect to.
        ProcessBuilder pb = new ProcessBuilder(
            System.getProperty("java.home") + "/bin/keytool",
            "-genkeypair", "-alias", "elsewhere",
            "-keyalg", "RSA", "-keysize", "2048",
            "-validity", "1",
            "-storepass", PASSWORD, "-keypass", PASSWORD,
            "-storetype", "PKCS12",
            "-keystore", wrongNameKeystore.toString(),
            "-dname", "CN=somewhere.else.invalid,O=Test,L=Test,ST=Test,C=US",
            "-ext", "san=dns:somewhere.else.invalid");
        pb.redirectErrorStream(true);
        Process process = pb.start();
        try (InputStream is = process.getInputStream()) {
            is.readAllBytes();
        }
        if (process.waitFor() != 0) {
            throw new IllegalStateException("keytool failed");
        }

        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        try (InputStream is = Files.newInputStream(wrongNameKeystore)) {
            keyStore.load(is, PASSWORD.toCharArray());
        }
        KeyManagerFactory kmf = KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm());
        kmf.init(keyStore, PASSWORD.toCharArray());
        serverContext = SSLContext.getInstance("TLS");
        serverContext.init(kmf.getKeyManagers(), null, new SecureRandom());
    }

    @AfterAll
    static void removeKeystore() throws Exception {
        if (wrongNameKeystore != null) {
            Files.deleteIfExists(wrongNameKeystore);
        }
    }

    @BeforeEach
    void setUp() throws Exception {
        serverSocket = (SSLServerSocket) serverContext.getServerSocketFactory().createServerSocket();
        serverSocket.bind(new InetSocketAddress("127.0.0.1", 0));
        serverSocket.setSoTimeout(15000);
        port = serverSocket.getLocalPort();
        serverExecutor = Executors.newCachedThreadPool();
        serverExecutor.submit(() -> {
            try {
                SSLSocket accepted = (SSLSocket) serverSocket.accept();
                accepted.startHandshake();
                return accepted;
            } catch (Exception e) {
                return null;
            }
        });
    }

    @AfterEach
    void tearDown() throws Exception {
        if (serverSocket != null && !serverSocket.isClosed()) {
            serverSocket.close();
        }
        if (serverExecutor != null) {
            serverExecutor.shutdownNow();
        }
    }

    /** Trusts the certificate itself, so only the name it was issued for is in question. */
    private TlsTransportConfiguration trustingConfig() {
        TlsTransportConfiguration config = new TlsTransportConfiguration();
        config.connectTimeout = 10000;
        config.readTimeout = 10000;
        config.receiveBufferSize = 8192;
        config.sendBufferSize = 8192;
        config.verifySsl = true;
        config.trustStoreFile = wrongNameKeystore.toString();
        config.trustStorePassword = PASSWORD;
        config.trustStoreType = "PKCS12";
        return config;
    }

    @Test
    void aCertificateIssuedForSomewhereElseIsRefused() {
        TlsTransportConfiguration config = trustingConfig();
        config.ignoreCommonName = false;

        assertThrows(TransportException.class, () -> {
            TlsTransportInstance instance = new TlsTransportInstance(
                new InetSocketAddress("127.0.0.1", port), config, AuditLog.builder().build());
            instance.close();
        }, "a certificate issued for somewhere.else.invalid must not be accepted for 127.0.0.1");
    }

    @Test
    void theSameCertificateIsAcceptedWhenTheCheckIsOptedOutOf() {
        TlsTransportConfiguration config = trustingConfig();
        config.ignoreCommonName = true;

        assertDoesNotThrow(() -> {
            TlsTransportInstance instance = new TlsTransportInstance(
                new InetSocketAddress("127.0.0.1", port), config, AuditLog.builder().build());
            instance.close();
        }, "ignore-common-name is the way out, and it has to work");
    }
}
