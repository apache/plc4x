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
import org.apache.plc4x.java.transport.tls.config.PskTlsTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.bouncycastle.tls.*;
import org.bouncycastle.tls.crypto.impl.bc.BcTlsCrypto;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.LockSupport;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;

/**
 * TLS Transport instance using Pre-Shared Key (PSK) authentication via Bouncy Castle.
 * Used for connecting to devices (e.g., Beckhoff TwinCAT 3) that support TLS-PSK.
 * <p>
 * Java's built-in SunJSSE provider does not support PSK cipher suites, so this
 * implementation uses Bouncy Castle's TlsClientProtocol + PSKTlsClient directly.
 */
public class PskTlsTransportInstance extends BaseTransportInstance<PskTlsTransportConfiguration> implements AsyncTransportInstance<PskTlsTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(PskTlsTransportInstance.class);
    private static final int DEFAULT_BUFFER_SIZE = 81920;
    private static final byte[] EMPTY_BYTES = new byte[0];

    private final Socket plainSocket;
    private final TlsClientProtocol tlsClientProtocol;
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

    public PskTlsTransportInstance(InetSocketAddress remoteAddress, PskTlsTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        LOGGER.debug("PskTlsTransportInstance: Creating PSK TLS connection to {}", remoteAddress);
        this.ringBuffer = new RingBuffer(configuration.receiveBufferSize);

        auditLog.write(AuditLogEventType.SYSTEM, String.format(
            "TLS-PSK configuration: target=%s:%d, psk-identity=%s, connect-timeout=%d, read-timeout=%d",
            remoteAddress.getHostName(), remoteAddress.getPort(),
            configuration.pskIdentity,
            configuration.connectTimeout, configuration.readTimeout));

        try {
            // Create plain socket and configure TCP options
            this.plainSocket = new Socket();

            if (configuration.localAddress != null && !configuration.localAddress.isEmpty()) {
                InetSocketAddress localAddr = new InetSocketAddress(configuration.localAddress, configuration.localPort);
                plainSocket.bind(localAddr);
                LOGGER.debug("Bound to local address {}:{}", configuration.localAddress, configuration.localPort);
                auditLog.write(AuditLogEventType.SYSTEM, String.format(
                    "Bound to local address: %s:%d", configuration.localAddress, configuration.localPort));
            }

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

            // Create Bouncy Castle TLS-PSK client
            BcTlsCrypto crypto = new BcTlsCrypto(new SecureRandom());
            byte[] identityBytes = configuration.pskIdentity.getBytes(StandardCharsets.UTF_8);
            byte[] keyBytes = configuration.getPskKeyBytes();

            BasicTlsPSKIdentity pskIdentity = new BasicTlsPSKIdentity(identityBytes, keyBytes);

            final int[] negotiatedCipherSuite = new int[]{-1};
            final boolean shouldLogSessionKeys = configuration.logSessionKeys;
            final AuditLog auditLogRef = auditLog;
            PSKTlsClient pskClient = new PSKTlsClient(crypto, pskIdentity) {
                @Override
                protected int[] getSupportedCipherSuites() {
                    // Ordered by preference: GCM first, then CBC-SHA256/384, then CBC-SHA1.
                    // TwinCAT XAE offers the CBC variants (0x00AF, 0x00AE, 0x008D, 0x008C).
                    return new int[]{
                        CipherSuite.TLS_PSK_WITH_AES_256_GCM_SHA384,
                        CipherSuite.TLS_PSK_WITH_AES_128_GCM_SHA256,
                        CipherSuite.TLS_PSK_WITH_AES_256_CBC_SHA384,
                        CipherSuite.TLS_PSK_WITH_AES_128_CBC_SHA256,
                        CipherSuite.TLS_PSK_WITH_AES_256_CBC_SHA,
                        CipherSuite.TLS_PSK_WITH_AES_128_CBC_SHA
                    };
                }

                @Override
                public void notifySelectedCipherSuite(int selectedCipherSuite) {
                    negotiatedCipherSuite[0] = selectedCipherSuite;
                    super.notifySelectedCipherSuite(selectedCipherSuite);
                }

                @Override
                public void notifyHandshakeComplete() throws IOException {
                    super.notifyHandshakeComplete();
                    if (shouldLogSessionKeys) {
                        try {
                            SecurityParameters sp = context.getSecurityParametersConnection();
                            byte[] clientRandom = sp.getClientRandom();
                            byte[] masterSecret = sp.getMasterSecret().extract();
                            String keyLogLine = "CLIENT_RANDOM " + hexEncode(clientRandom) + " " + hexEncode(masterSecret);
                            LOGGER.info("TLS session key (SSLKEYLOGFILE format): {}", keyLogLine);
                            auditLogRef.write(AuditLogEventType.SYSTEM,
                                "SSLKEYLOGFILE: " + keyLogLine);
                        } catch (Exception e) {
                            LOGGER.warn("Failed to extract TLS session keys", e);
                        }
                    }
                }

                private String hexEncode(byte[] bytes) {
                    StringBuilder sb = new StringBuilder(bytes.length * 2);
                    for (byte b : bytes) {
                        sb.append(String.format("%02x", b & 0xff));
                    }
                    return sb.toString();
                }
            };

            // Wrap socket streams with TLS protocol
            auditLog.write(AuditLogEventType.SYSTEM, String.format(
                "TLS-PSK handshake starting: identity=%s, cipher-suites=[PSK_AES_256_GCM_SHA384, PSK_AES_128_GCM_SHA256, PSK_AES_256_CBC_SHA384, PSK_AES_128_CBC_SHA256, PSK_AES_256_CBC_SHA, PSK_AES_128_CBC_SHA]",
                configuration.pskIdentity));

            this.tlsClientProtocol = new TlsClientProtocol(
                plainSocket.getInputStream(), plainSocket.getOutputStream());
            try {
                tlsClientProtocol.connect(pskClient);
            } catch (TlsFatalAlertReceived e) {
                String message = String.format("TLS-PSK connection failed: %s. Please check, if the selected PSK information correctly configured in the target PLC.", e.getMessage());
                auditLog.write(AuditLogEventType.SYSTEM, message);
                closeSocketQuietly();
                throw new TransportException(message, e);
            }

            // Get streams from the TLS protocol
            this.inputStream = tlsClientProtocol.getInputStream();
            this.outputStream = tlsClientProtocol.getOutputStream();

            // Start reader thread
            this.readerThread = new Thread(this::runReaderLoop,
                "TLS-PSK-Reader-" + remoteAddress.getHostName() + ":" + remoteAddress.getPort());
            this.readerThread.start();

            LOGGER.info("TLS-PSK connection established to {}:{}", remoteAddress.getHostName(), remoteAddress.getPort());

            String negotiatedCipher = resolveCipherSuiteName(negotiatedCipherSuite[0]);
            auditLog.write(AuditLogEventType.CONNECT, String.format(
                "TLS-PSK connected to: %s:%d (cipher=%s, identity=%s)",
                remoteAddress.getHostName(), remoteAddress.getPort(),
                negotiatedCipher, configuration.pskIdentity));

        } catch (IOException e) {
            String errorMsg = formatPskError(e, remoteAddress);
            LOGGER.error(errorMsg, e);
            auditLog.write(AuditLogEventType.ERROR, "TLS-PSK connection error: " + errorMsg);
            logPskErrorChain(auditLog, e);
            closeSocketQuietly();
            throw new TransportException(errorMsg, e);
        }
    }

    /**
     * Resolves a BC CipherSuite int constant to a human-readable name.
     */
    static String resolveCipherSuiteName(int cipherSuite) {
        return switch (cipherSuite) {
            case CipherSuite.TLS_PSK_WITH_AES_256_GCM_SHA384 -> "TLS_PSK_WITH_AES_256_GCM_SHA384";
            case CipherSuite.TLS_PSK_WITH_AES_128_GCM_SHA256 -> "TLS_PSK_WITH_AES_128_GCM_SHA256";
            case CipherSuite.TLS_PSK_WITH_AES_256_CBC_SHA384 -> "TLS_PSK_WITH_AES_256_CBC_SHA384";
            case CipherSuite.TLS_PSK_WITH_AES_128_CBC_SHA256 -> "TLS_PSK_WITH_AES_128_CBC_SHA256";
            case CipherSuite.TLS_PSK_WITH_AES_256_CBC_SHA -> "TLS_PSK_WITH_AES_256_CBC_SHA";
            case CipherSuite.TLS_PSK_WITH_AES_128_CBC_SHA -> "TLS_PSK_WITH_AES_128_CBC_SHA";
            default -> "0x" + Integer.toHexString(cipherSuite);
        };
    }

    /**
     * Formats PSK-specific error messages for better user feedback.
     */
    String formatPskError(IOException e, InetSocketAddress remoteAddress) {
        String message = e.getMessage();
        if (message == null) {
            message = e.getClass().getSimpleName();
        }

        if (e instanceof TlsFatalAlert tlsAlert) {
            short desc = tlsAlert.getAlertDescription();
            if (desc == AlertDescription.decrypt_error) {
                return String.format("PSK authentication failed for %s:%d — the PSK key is likely incorrect. Error: %s",
                    remoteAddress.getHostName(), remoteAddress.getPort(), message);
            } else if (desc == AlertDescription.unknown_psk_identity) {
                return String.format("PSK identity not recognized by %s:%d — check the configured psk-identity. Error: %s",
                    remoteAddress.getHostName(), remoteAddress.getPort(), message);
            } else if (desc == AlertDescription.handshake_failure) {
                return String.format("TLS-PSK handshake failed with %s:%d — server may not support PSK cipher suites. Error: %s",
                    remoteAddress.getHostName(), remoteAddress.getPort(), message);
            }
        }

        if (message.contains("Connection refused")) {
            return String.format("Connection refused to %s:%d. Server may not be listening for TLS connections.",
                remoteAddress.getHostName(), remoteAddress.getPort());
        }

        return String.format("Failed to establish TLS-PSK connection to %s:%d — %s",
            remoteAddress.getHostName(), remoteAddress.getPort(), message);
    }

    /**
     * Logs the full exception chain for PSK errors to help diagnose handshake failures.
     */
    static void logPskErrorChain(AuditLog auditLog, Throwable e) {
        int depth = 0;
        Throwable current = e;
        while (current != null && depth < 10) {
            auditLog.write(AuditLogEventType.ERROR, String.format(
                "  TLS-PSK error cause [%d]: %s: %s",
                depth, current.getClass().getName(), current.getMessage()));
            current = current.getCause();
            depth++;
        }
    }

    private void closeSocketQuietly() {
        try {
            if (plainSocket != null && !plainSocket.isClosed()) {
                plainSocket.close();
            }
        } catch (IOException ignored) {}
    }

    public InetSocketAddress getLocalAddress() {
        return (InetSocketAddress) plainSocket.getLocalSocketAddress();
    }

    public InetSocketAddress getRemoteAddress() {
        return (InetSocketAddress) plainSocket.getRemoteSocketAddress();
    }

    @Override
    public boolean isOpen() {
        return open && plainSocket.isConnected() && !plainSocket.isClosed();
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

            LOGGER.trace("Wrote {} bytes over TLS-PSK to {}", bytes.length, getRemoteAddress());

            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Write: " + StaticHelper.ENCODE_HEX(bytes));
            }
        } catch (IOException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in write: " + e.getMessage());
            throw new TransportException("Failed to write data over TLS-PSK", e);
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
                    tlsClientProtocol.close();
                } catch (IOException e) {
                    LOGGER.debug("Error closing TLS protocol", e);
                }

                try {
                    plainSocket.close();
                } catch (IOException e) {
                    LOGGER.debug("Error closing socket", e);
                }

                if (readerThread != null) {
                    try {
                        readerThread.join(1000);
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    }
                }

                LOGGER.debug("TLS-PSK connection closed");
                getAuditLog().write(AuditLogEventType.CLOSE, "TLS-PSK connection closed");

            } finally {
                readLock.unlock();
            }
        } finally {
            writeLock.unlock();
        }
    }

    private void ensureOpen() throws TransportException {
        if (!isOpen()) {
            throw new TransportException("TLS-PSK transport is closed");
        }
    }

    // ========== AsyncTransportInstance Implementation ==========

    @Override
    public void registerDataListener(Runnable listener) {
        this.dataListener = listener;
        LOGGER.debug("Data listener registered for TLS-PSK transport");
    }

    @Override
    public void removeDataListener() {
        this.dataListener = null;
        LOGGER.debug("Data listener removed from TLS-PSK transport");
    }

    @Override
    public void registerDisconnectListener(Consumer<Throwable> listener) {
        this.disconnectListener = listener;
        LOGGER.debug("Disconnect listener registered for TLS-PSK transport");
    }

    @Override
    public void removeDisconnectListener() {
        this.disconnectListener = null;
        LOGGER.debug("Disconnect listener removed from TLS-PSK transport");
    }

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
        LOGGER.debug("TLS-PSK reader loop started");
        byte[] buffer = new byte[DEFAULT_BUFFER_SIZE];

        while (open && !Thread.currentThread().isInterrupted()) {
            try {
                int free = ringBuffer.remainingForWriting();
                if (free == 0) {
                    // The ring is full and whoever consumes it has not caught up. Park briefly and
                    // look again rather than reading bytes we would have to throw away: only the
                    // codec knows where a frame ends, so a byte dropped here takes the middle out
                    // of somebody's message. Never a disconnect - a consumer is allowed to lag.
                    LockSupport.parkNanos(200_000L);
                    continue;
                }

                // Bound the read to the room actually available, so the ring cannot overflow.
                int bytesRead = inputStream.read(buffer, 0, Math.min(buffer.length, free));

                if (bytesRead > 0) {
                    readLock.lock();
                    try {
                        int written = ringBuffer.write(buffer, 0, bytesRead);
                        if (written < bytesRead) {
                            // The read was bounded to the free space, so this cannot happen unless
                            // that reasoning is wrong somewhere - say so loudly rather than
                            // shrugging off lost bytes.
                            LOGGER.error("Wrote only {} of {} bytes into a ring buffer that reported room",
                                written, bytesRead);
                        }

                        safeRun(dataListener);
                    } finally {
                        readLock.unlock();
                    }
                } else if (bytesRead == -1) {
                    LOGGER.info("TLS-PSK connection closed by remote host");
                    getAuditLog().write(AuditLogEventType.SYSTEM,
                        "TLS-PSK connection closed gracefully by remote host");
                    open = false;
                    notifyDisconnect(null);
                    break;
                }
            } catch (IOException e) {
                if (open) {
                    LOGGER.error("Error reading from TLS-PSK socket", e);
                    getAuditLog().write(AuditLogEventType.ERROR, "TLS-PSK read error: " + e.getMessage());
                    open = false;
                    notifyDisconnect(e);
                }
                break;
            }
        }

        LOGGER.debug("TLS-PSK reader loop stopped");
    }

}
