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

import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.transport.tls.config.TlsTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.net.ssl.*;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.security.KeyManagementException;
import java.security.KeyStore;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.security.cert.CertificateParsingException;
import java.security.cert.X509Certificate;
import java.util.Arrays;
import java.util.Collection;
import java.util.List;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;
import java.util.stream.Collectors;

/**
 * TLS Transport instance that provides encrypted communication using SSL/TLS.
 * Supports both standard certificate validation and self-signed certificate bypass.
 * <p>
 * This implementation uses blocking SSLSocket for simplicity and compatibility.
 * For high-performance async I/O, see the Java 21 variant with virtual threads.
 */
public class TlsTransportInstance extends BaseTransportInstance<TlsTransportConfiguration> implements AsyncTransportInstance<TlsTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(TlsTransportInstance.class);
    private static final int DEFAULT_BUFFER_SIZE = 81920;
    private static final byte[] EMPTY_BYTES = new byte[0];

    private final SSLSocket sslSocket;
    private final InputStream inputStream;
    private final OutputStream outputStream;
    private final RingBuffer ringBuffer;
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();
    private volatile boolean open = true;

    // Async support
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;
    private Thread readerThread;

    public TlsTransportInstance(InetSocketAddress remoteAddress, TlsTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        LOGGER.debug("TlsTransportInstance: Creating TLS connection (verify-ssl={})", configuration.isVerifySsl());
        this.ringBuffer = new RingBuffer(configuration.receiveBufferSize);

        auditLog.write(AuditLogEventType.SYSTEM, String.format(
            "TLS configuration: target=%s:%d, verify-ssl=%s, tls-version=%s, connect-timeout=%d, read-timeout=%d",
            remoteAddress.getHostName(), remoteAddress.getPort(),
            configuration.isVerifySsl(),
            configuration.getTlsVersion() != null ? configuration.getTlsVersion() : "auto (TLSv1.3/TLSv1.2)",
            configuration.connectTimeout, configuration.readTimeout));

        try {
            // Create SSLContext with appropriate trust and key configuration
            SSLContext sslContext = createSslContext(configuration);
            SSLSocketFactory socketFactory = sslContext.getSocketFactory();

            auditLog.write(AuditLogEventType.SYSTEM, String.format(
                "SSLContext initialized: provider=%s, protocol=%s, verify-ssl=%s",
                sslContext.getProvider().getName(), sslContext.getProtocol(), configuration.isVerifySsl()));

            // Create plain socket first for timeout control
            Socket plainSocket = new Socket();

            // Bind to a local address if specified
            if (configuration.localAddress != null && !configuration.localAddress.isEmpty()) {
                InetSocketAddress localAddr = new InetSocketAddress(configuration.localAddress, configuration.localPort);
                plainSocket.bind(localAddr);
                LOGGER.debug("Bound to local address {}:{}", configuration.localAddress, configuration.localPort);
                auditLog.write(AuditLogEventType.SYSTEM, String.format(
                    "Bound to local address: %s:%d", configuration.localAddress, configuration.localPort));
            }

            // Configure socket options before connecting
            plainSocket.setTcpNoDelay(configuration.tcpNoDelay);
            plainSocket.setKeepAlive(configuration.keepAlive);

            if (configuration.sendBufferSize > 0) {
                plainSocket.setSendBufferSize(configuration.sendBufferSize);
            }
            if (configuration.receiveBufferSize > 0) {
                plainSocket.setReceiveBufferSize(configuration.receiveBufferSize);
            }
            if (configuration.readTimeout > 0) {
                plainSocket.setSoTimeout(configuration.readTimeout);
            }

            // Connect with timeout
            auditLog.write(AuditLogEventType.SYSTEM, String.format(
                "TCP connecting to %s:%d (timeout=%dms)",
                remoteAddress.getHostName(), remoteAddress.getPort(), configuration.connectTimeout));
            plainSocket.connect(remoteAddress, configuration.connectTimeout);
            auditLog.write(AuditLogEventType.SYSTEM, String.format(
                "TCP connected, local=%s:%d -> remote=%s:%d",
                plainSocket.getLocalAddress().getHostAddress(), plainSocket.getLocalPort(),
                remoteAddress.getHostName(), remoteAddress.getPort()));

            // Wrap in SSLSocket
            this.sslSocket = (SSLSocket) socketFactory.createSocket(
                plainSocket,
                remoteAddress.getHostName(),
                remoteAddress.getPort(),
                true // auto-close underlying socket
            );

            // Configure TLS protocol versions
            String configuredTlsVersion = configuration.getTlsVersion();
            String[] enabledProtocols;
            if (configuredTlsVersion != null && !configuredTlsVersion.isEmpty()) {
                // Use specifically configured TLS version
                LOGGER.debug("Using configured TLS version: {}", configuredTlsVersion);
                enabledProtocols = new String[]{configuredTlsVersion};
            } else {
                // Default: prefer TLS 1.3, fall back to 1.2
                enabledProtocols = new String[]{"TLSv1.3", "TLSv1.2"};
            }
            sslSocket.setEnabledProtocols(enabledProtocols);

            auditLog.write(AuditLogEventType.SYSTEM, String.format(
                "TLS handshake starting: enabled-protocols=%s, enabled-ciphers=%s",
                Arrays.toString(enabledProtocols),
                Arrays.toString(sslSocket.getEnabledCipherSuites())));

            // Perform TLS handshake
            LOGGER.debug("Starting TLS handshake with {}", remoteAddress);
            sslSocket.startHandshake();

            // Log detailed TLS session info
            logTlsSessionInfo();
            logCertificateChain(auditLog);

            // Get streams
            this.inputStream = sslSocket.getInputStream();
            this.outputStream = sslSocket.getOutputStream();

            // Start reader thread for async support
            this.readerThread = new Thread(this::runReaderLoop,
                "TLS-Reader-" + remoteAddress.getHostName() + ":" + remoteAddress.getPort());
            this.readerThread.start();

            LOGGER.info("TLS connection established to {}:{}", remoteAddress.getHostName(), remoteAddress.getPort());

            SSLSession session = sslSocket.getSession();
            auditLog.write(AuditLogEventType.CONNECT, String.format(
                "TLS connected to: %s:%d (protocol=%s, cipher=%s, session-id=%s)",
                remoteAddress.getHostName(), remoteAddress.getPort(),
                session.getProtocol(), session.getCipherSuite(),
                bytesToHex(session.getId())));

        } catch (IOException e) {
            String errorMsg = formatTlsError(e, remoteAddress);
            LOGGER.error(errorMsg, e);
            auditLog.write(AuditLogEventType.ERROR, "TLS connection error: " + errorMsg);
            logTlsErrorChain(auditLog, e);
            throw new TransportException(errorMsg, e);
        } catch (NoSuchAlgorithmException | KeyManagementException e) {
            String errorMsg = "Failed to initialize TLS context: " + e.getMessage();
            LOGGER.error(errorMsg, e);
            auditLog.write(AuditLogEventType.ERROR, errorMsg);
            throw new TransportException(errorMsg, e);
        }
    }

    /**
     * Creates an SSLContext with the appropriate KeyManager and TrustManager configuration.
     *
     * @param configuration the TLS transport configuration
     * @return configured SSLContext
     */
    private SSLContext createSslContext(TlsTransportConfiguration configuration)
            throws NoSuchAlgorithmException, KeyManagementException, TransportException {
        SSLContext sslContext = SSLContext.getInstance("TLS");

        // Load client keystore if configured (for mutual TLS / client certificate authentication)
        KeyManager[] keyManagers = null;
        if (configuration.keystore != null && !configuration.keystore.isEmpty()) {
            keyManagers = loadKeyManagers(configuration);
        }

        TrustManager[] trustManagers = null;
        if (!configuration.isVerifySsl()) {
            // Use trust-all TrustManager for self-signed certificates
            LOGGER.warn("SSL certificate verification is DISABLED. Connection is encrypted but NOT protected against MITM attacks.");
            trustManagers = new TrustManager[]{
                new X509ExtendedTrustManager() {
                    @Override
                    public void checkClientTrusted(X509Certificate[] chain, String authType) {}
                    @Override
                    public void checkServerTrusted(X509Certificate[] chain, String authType) {}
                    @Override
                    public X509Certificate[] getAcceptedIssuers() { return new X509Certificate[0]; }
                    @Override
                    public void checkClientTrusted(X509Certificate[] chain, String authType, Socket socket) {}
                    @Override
                    public void checkServerTrusted(X509Certificate[] chain, String authType, Socket socket) {}
                    @Override
                    public void checkClientTrusted(X509Certificate[] chain, String authType, SSLEngine engine) {}
                    @Override
                    public void checkServerTrusted(X509Certificate[] chain, String authType, SSLEngine engine) {}
                }
            };
        }

        sslContext.init(keyManagers, trustManagers, new SecureRandom());
        return sslContext;
    }

    /**
     * Loads KeyManagers from the configured keystore for client certificate authentication.
     */
    private KeyManager[] loadKeyManagers(TlsTransportConfiguration configuration) throws TransportException {
        String keystoreType = configuration.keystoreType != null ? configuration.keystoreType : "PKCS12";
        char[] password = configuration.keystorePassword != null
            ? configuration.keystorePassword.toCharArray()
            : new char[0];

        try {
            KeyStore keyStore = KeyStore.getInstance(keystoreType);
            try (FileInputStream fis = new FileInputStream(configuration.keystore)) {
                keyStore.load(fis, password);
            }

            KeyManagerFactory kmf = KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm());
            kmf.init(keyStore, password);

            LOGGER.info("Loaded client keystore: type={}, path={}", keystoreType, configuration.keystore);
            getAuditLog().write(AuditLogEventType.SYSTEM, String.format(
                "Client keystore loaded: type=%s, path=%s, aliases=%d",
                keystoreType, configuration.keystore, keyStore.size()));

            return kmf.getKeyManagers();
        } catch (Exception e) {
            throw new TransportException(String.format(
                "Failed to load client keystore '%s' (type=%s): %s",
                configuration.keystore, keystoreType, e.getMessage()), e);
        }
    }

    /**
     * Formats TLS-specific error messages for better user feedback.
     */
    private String formatTlsError(IOException e, InetSocketAddress remoteAddress) {
        String message = e.getMessage();
        if (message == null) {
            message = e.getClass().getSimpleName();
        }

        if (message.contains("PKIX") || message.contains("certificate")) {
            return String.format("Server certificate not trusted for %s:%d. Use verify-ssl=false for self-signed certificates. Error: %s",
                remoteAddress.getHostName(), remoteAddress.getPort(), message);
        } else if (message.contains("handshake") || message.contains("Handshake")) {
            return String.format("TLS handshake failed with %s:%d. Server may not support TLS on this port. Error: %s",
                remoteAddress.getHostName(), remoteAddress.getPort(), message);
        } else if (message.contains("Connection refused")) {
            return String.format("Connection refused to %s:%d. Server may not be listening for TLS connections.",
                remoteAddress.getHostName(), remoteAddress.getPort());
        } else {
            return String.format("Failed to establish TLS connection to %s:%d - %s",
                remoteAddress.getHostName(), remoteAddress.getPort(), message);
        }
    }

    /**
     * Logs TLS session information for debugging and auditing.
     */
    private void logTlsSessionInfo() {
        SSLSession session = sslSocket.getSession();
        LOGGER.debug("TLS Session Info - Protocol: {}, Cipher Suite: {}, Peer Host: {}",
            session.getProtocol(), session.getCipherSuite(), session.getPeerHost());
    }

    /**
     * Logs the server certificate chain details to the audit log after a successful handshake.
     */
    private void logCertificateChain(AuditLog auditLog) {
        try {
            SSLSession session = sslSocket.getSession();
            java.security.cert.Certificate[] peerCerts = session.getPeerCertificates();

            auditLog.write(AuditLogEventType.SYSTEM, String.format(
                "Server certificate chain: %d certificate(s)", peerCerts.length));

            for (int i = 0; i < peerCerts.length; i++) {
                if (peerCerts[i] instanceof X509Certificate x509) {
                    String label = (i == 0) ? "Server" : (i == peerCerts.length - 1) ? "Root CA" : "Intermediate CA";
                    auditLog.write(AuditLogEventType.SYSTEM, String.format(
                        "  [%d] %s: subject=%s, issuer=%s, serial=%s, valid=%s to %s, sig-alg=%s%s",
                        i, label,
                        x509.getSubjectX500Principal().getName(),
                        x509.getIssuerX500Principal().getName(),
                        x509.getSerialNumber().toString(16),
                        x509.getNotBefore(), x509.getNotAfter(),
                        x509.getSigAlgName(),
                        formatSubjectAltNames(x509)));
                }
            }
        } catch (SSLPeerUnverifiedException e) {
            auditLog.write(AuditLogEventType.SYSTEM,
                "Server certificate chain: peer not verified (certificate verification may be disabled)");
        }
    }

    /**
     * Formats Subject Alternative Names (SANs) from a certificate for logging.
     */
    private static String formatSubjectAltNames(X509Certificate cert) {
        try {
            Collection<List<?>> sans = cert.getSubjectAlternativeNames();
            if (sans == null || sans.isEmpty()) {
                return "";
            }
            String sanStr = sans.stream()
                .map(san -> {
                    int type = (Integer) san.get(0);
                    String typeName = switch (type) {
                        case 0 -> "otherName";
                        case 1 -> "rfc822Name";
                        case 2 -> "dNSName";
                        case 4 -> "directoryName";
                        case 6 -> "URI";
                        case 7 -> "IP";
                        default -> "type" + type;
                    };
                    return typeName + "=" + san.get(1);
                })
                .collect(Collectors.joining(", "));
            return ", SANs=[" + sanStr + "]";
        } catch (CertificateParsingException e) {
            return ", SANs=[parse error: " + e.getMessage() + "]";
        }
    }

    /**
     * Logs the full exception chain for TLS errors to help diagnose handshake failures.
     */
    private static void logTlsErrorChain(AuditLog auditLog, Throwable e) {
        int depth = 0;
        Throwable current = e;
        while (current != null && depth < 10) {
            auditLog.write(AuditLogEventType.ERROR, String.format(
                "  TLS error cause [%d]: %s: %s",
                depth, current.getClass().getName(), current.getMessage()));
            current = current.getCause();
            depth++;
        }
    }

    /**
     * Converts a byte array to a hex string for session ID logging.
     */
    private static String bytesToHex(byte[] bytes) {
        if (bytes == null || bytes.length == 0) {
            return "(empty)";
        }
        StringBuilder sb = new StringBuilder(bytes.length * 2);
        for (byte b : bytes) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    public InetSocketAddress getRemoteAddress() {
        return (InetSocketAddress) sslSocket.getRemoteSocketAddress();
    }

    public InetSocketAddress getLocalAddress() {
        return (InetSocketAddress) sslSocket.getLocalSocketAddress();
    }

    /**
     * Returns the TLS protocol version being used.
     */
    public String getTlsProtocol() {
        return sslSocket.getSession().getProtocol();
    }

    /**
     * Returns the cipher suite being used.
     */
    public String getCipherSuite() {
        return sslSocket.getSession().getCipherSuite();
    }

    @Override
    public boolean isOpen() {
        return open && sslSocket.isConnected() && !sslSocket.isClosed();
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

            byte[] bytes = ringBuffer.read(numBytes);

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

            outputStream.write(bytes);
            outputStream.flush();

            LOGGER.trace("Wrote {} bytes over TLS to {}", bytes.length, getRemoteAddress());

            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Write: " + StaticHelper.ENCODE_HEX(bytes));
            }
        } catch (IOException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in write: " + e.getMessage());
            throw new TransportException("Failed to write data over TLS", e);
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

                // Close streams and socket
                try {
                    inputStream.close();
                } catch (IOException e) {
                    LOGGER.debug("Error closing input stream", e);
                }

                try {
                    outputStream.close();
                } catch (IOException e) {
                    LOGGER.debug("Error closing output stream", e);
                }

                try {
                    sslSocket.close();
                } catch (IOException e) {
                    LOGGER.debug("Error closing SSL socket", e);
                }

                // Wait for reader thread to finish
                if (readerThread != null) {
                    try {
                        readerThread.join(1000);
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    }
                }

                LOGGER.debug("TLS connection closed");
                getAuditLog().write(AuditLogEventType.CLOSE, String.format(
                    "TLS connection closed (protocol=%s, cipher=%s)",
                    sslSocket.getSession().getProtocol(), sslSocket.getSession().getCipherSuite()));

            } finally {
                readLock.unlock();
            }
        } finally {
            writeLock.unlock();
        }
    }

    private void ensureOpen() throws TransportException {
        if (!isOpen()) {
            throw new TransportException("TLS transport is closed");
        }
    }

    // ========== AsyncTransportInstance Implementation ==========

    @Override
    public void registerDataListener(Runnable listener) {
        this.dataListener = listener;
        LOGGER.debug("Data listener registered for TLS transport");
    }

    @Override
    public void removeDataListener() {
        this.dataListener = null;
        LOGGER.debug("Data listener removed from TLS transport");
    }

    @Override
    public void registerDisconnectListener(Consumer<Throwable> listener) {
        this.disconnectListener = listener;
        LOGGER.debug("Disconnect listener registered for TLS transport");
    }

    @Override
    public void removeDisconnectListener() {
        this.disconnectListener = null;
        LOGGER.debug("Disconnect listener removed from TLS transport");
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
     * Runs a listener so that a failure inside it costs the frame it was handling and nothing more.
     * The reader loop is the only thing reading the socket, and the peer decides what arrives, so
     * one frame the codec cannot make sense of must not end the loop that would have read the next.
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
     * Reader loop that continuously reads from the TLS socket and populates the ring buffer.
     * Notifies listeners when data arrives.
     */
    private void runReaderLoop() {
        try {
            runReaderLoopInternal();
        } catch (Throwable t) {
            // Nothing else reads this socket. If the loop ever leaves by an unexpected route the
            // transport must stop claiming to be open, rather than leave callers waiting on a
            // connection that can no longer receive.
            LOGGER.error("Reader loop failed", t);
        }
        if (open) {
            open = false;
            notifyDisconnect(null);
        }
    }

    private void runReaderLoopInternal() {
        LOGGER.debug("TLS reader loop started");
        byte[] buffer = new byte[DEFAULT_BUFFER_SIZE];

        while (open && !Thread.currentThread().isInterrupted()) {
            try {
                int bytesRead = inputStream.read(buffer);

                if (bytesRead > 0) {
                    readLock.lock();
                    try {
                        // Write to ring buffer
                        int written = ringBuffer.write(buffer, 0, bytesRead);
                        if (written < bytesRead) {
                            LOGGER.warn("Ring buffer full, lost {} bytes", bytesRead - written);
                        }

                        // Notify listener
                        safeRun(dataListener);
                    } finally {
                        readLock.unlock();
                    }
                } else if (bytesRead == -1) {
                    // Connection closed gracefully by server
                    LOGGER.info("TLS connection closed by remote host");
                    getAuditLog().write(AuditLogEventType.SYSTEM,
                        "TLS connection closed gracefully by remote host");
                    open = false;
                    notifyDisconnect(null);  // null indicates graceful close
                    break;
                }
            } catch (IOException e) {
                if (open) {
                    LOGGER.error("Error reading from TLS socket", e);
                    getAuditLog().write(AuditLogEventType.ERROR, "TLS read error: " + e.getMessage());
                    open = false;
                    notifyDisconnect(e);
                }
                break;
            }
        }

        LOGGER.debug("TLS reader loop stopped");
    }

}
