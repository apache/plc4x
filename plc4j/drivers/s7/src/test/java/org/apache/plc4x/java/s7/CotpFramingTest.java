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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.transport.cotp.CotpTransportInstance;

import org.apache.plc4x.java.transport.cotp.config.CotpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The peer writes the TPKT length field, so it decides how the transport frames what follows. A
 * length that cannot describe a packet has to be reported, not waited on, and a packet taken off
 * the TCP stream has to be accounted for whether or not there is room for its payload.
 */
class CotpFramingTest {

    /** A real Connection Confirm from the S7 parser testsuite, so the handshake completes. */
    private static final byte[] CONNECTION_CONFIRM = decodeHex("0300001611d0000f000b00c0010ac1020311c2020100");

    private ServerSocket serverSocket;
    private ExecutorService serverExecutor;
    private int port;

    private static byte[] decodeHex(String hex) {
        byte[] out = new byte[hex.length() / 2];
        for (int i = 0; i < out.length; i++) {
            out[i] = (byte) Integer.parseInt(hex.substring(i * 2, i * 2 + 2), 16);
        }
        return out;
    }

    @BeforeEach
    void setUp() throws Exception {
        serverSocket = new ServerSocket();
        serverSocket.bind(new InetSocketAddress("127.0.0.1", 0));
        port = serverSocket.getLocalPort();
        serverExecutor = Executors.newSingleThreadExecutor();
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

    private CotpTransportConfiguration config() {
        CotpTransportConfiguration configuration = new CotpTransportConfiguration();
        configuration.cotpConnectionTimeout = 3000;
        // The driver normally fills these in from the defaults.
        configuration.connectTimeout = 5000;
        configuration.receiveBufferSize = 8192;
        configuration.sendBufferSize = 8192;
        return configuration;
    }

    /**
     * Answers the connection request, then writes whatever the test wants to send afterwards.
     */
    private CompletableFuture<Socket> peerThatConfirmsThenSends(byte[] afterHandshake) {
        CompletableFuture<Socket> accepted = new CompletableFuture<>();
        serverExecutor.submit(() -> {
            try {
                Socket socket = serverSocket.accept();
                InputStream in = socket.getInputStream();
                // Read the connection request (22 bytes) before answering.
                byte[] request = new byte[22];
                int read = 0;
                while (read < request.length) {
                    int n = in.read(request, read, request.length - read);
                    if (n < 0) {
                        break;
                    }
                    read += n;
                }
                OutputStream out = socket.getOutputStream();
                out.write(CONNECTION_CONFIRM);
                out.flush();
                if (afterHandshake != null && afterHandshake.length > 0) {
                    Thread.sleep(150);
                    out.write(afterHandshake);
                    out.flush();
                }
                accepted.complete(socket);
            } catch (Exception e) {
                accepted.completeExceptionally(e);
            }
            return null;
        });
        return accepted;
    }

    /** A real TPKT+COTP data packet from the testsuite, carrying an 18-byte S7 payload. */
    private static final byte[] GOOD_DATA_PACKET =
        decodeHex("0300001902f08132010000000000080000f0000008000803f0");

    /**
     * A TPKT header declaring a length of zero describes no packet at all, and reading zero bytes
     * consumes nothing - so the same header was framed again on the next pass, and the one after.
     * The packet queued behind it never got a turn.
     */
    @Test
    void aTpktLengthThatDescribesNoPacketDoesNotStall() throws Exception {
        byte[] crafted = new byte[4 + GOOD_DATA_PACKET.length];
        // Length 0x0000, then a packet that is perfectly readable.
        crafted[0] = 0x03;
        crafted[1] = 0x00;
        crafted[2] = 0x00;
        crafted[3] = 0x00;
        System.arraycopy(GOOD_DATA_PACKET, 0, crafted, 4, GOOD_DATA_PACKET.length);
        peerThatConfirmsThenSends(crafted);

        CotpTransportInstance instance = new CotpTransportInstance("127.0.0.1", port, config(),
            AuditLog.builder().build());
        try {
            long deadline = System.currentTimeMillis() + 4000;
            int available = 0;
            while (available == 0 && System.currentTimeMillis() < deadline) {
                try {
                    available = instance.getNumBytesAvailable();
                } catch (RuntimeException e) {
                    // Reporting the bad header is fine; never getting past it is not.
                }
                Thread.sleep(20);
            }
            assertTrue(available > 0,
                "the packet behind an unusable TPKT header must still be delivered");
        } finally {
            instance.close();
        }
    }

    /**
     * The connect phase waits for the Connection Confirm. It used to do so by asking the transport
     * how many bytes were available, over and over, with nothing between the asks.
     */
    @Test
    void waitingForTheConnectionConfirmDoesNotBurnACore() throws Exception {
        CompletableFuture<Socket> accepted = new CompletableFuture<>();
        serverExecutor.submit(() -> {
            try {
                Socket socket = serverSocket.accept();
                InputStream in = socket.getInputStream();
                byte[] request = new byte[22];
                int read = 0;
                while (read < request.length) {
                    int n = in.read(request, read, request.length - read);
                    if (n < 0) {
                        break;
                    }
                    read += n;
                }
                // Make the connect phase wait a while before answering.
                Thread.sleep(1200);
                OutputStream out = socket.getOutputStream();
                out.write(CONNECTION_CONFIRM);
                out.flush();
                accepted.complete(socket);
            } catch (Exception e) {
                accepted.completeExceptionally(e);
            }
            return null;
        });

        long cpuBefore = threadCpuTimeOfThisProcess();
        CotpTransportInstance instance = new CotpTransportInstance("127.0.0.1", port, config(),
            AuditLog.builder().build());
        long cpuSpent = threadCpuTimeOfThisProcess() - cpuBefore;
        try {
            accepted.get(5, TimeUnit.SECONDS);
            // 1.2 seconds of waiting should cost almost nothing. A spin costs about a second of it.
            assertTrue(cpuSpent < 500_000_000L,
                "waiting for the confirm spent " + (cpuSpent / 1_000_000) + " ms of CPU on a 1200 ms wait");
        } finally {
            instance.close();
        }
    }

    private static long threadCpuTimeOfThisProcess() {
        java.lang.management.ThreadMXBean bean = java.lang.management.ManagementFactory.getThreadMXBean();
        long total = 0;
        for (long id : bean.getAllThreadIds()) {
            long cpu = bean.getThreadCpuTime(id);
            if (cpu > 0) {
                total += cpu;
            }
        }
        return total;
    }

    @Test
    void aWellFormedPacketStillConnects() {
        peerThatConfirmsThenSends(null);
        assertDoesNotThrow(() -> {
            CotpTransportInstance instance = new CotpTransportInstance("127.0.0.1", port, config(),
                AuditLog.builder().build());
            instance.close();
        });
    }
}
