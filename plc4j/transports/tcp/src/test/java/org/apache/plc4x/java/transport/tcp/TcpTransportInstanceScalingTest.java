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

import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import java.net.InetSocketAddress;
import java.nio.channels.ServerSocketChannel;
import java.nio.channels.SocketChannel;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Concurrency-model probe: opens many idle connections and reports how many carrier (platform)
 * threads the virtual-thread scheduler is using. With the old selector model, each connection's
 * virtual thread blocks in {@code Selector.select()}, which does not unmount the virtual thread,
 * so the scheduler ties up a carrier per connection (toward maxPoolSize=256) — measured at 201
 * carriers for 200 connections on both JDK 21 and JDK 25, i.e. NOT the synchronized-monitor
 * pinning that JEP 491 (JDK 24) removed. With the blocking-read model, the read virtual threads
 * unmount via the NIO poller, so the carrier pool stays small and flat regardless of connection
 * count (measured 2-3 for 200 connections).
 *
 * The test asserts only that all connections stay open (lenient — it is an evidence probe, not a
 * brittle threshold). The before/after signal is the {@code CARRIER_COUNT=...} line printed to
 * stdout and, when run with {@code -Djdk.tracePinnedThreads=full}, the presence/absence of
 * pinned-thread stack traces.
 */
class TcpTransportInstanceScalingTest {

    private static final int CONNECTIONS = 200;

    @Test
    @Disabled("Concurrency-model evidence probe — opens 200 sockets and sleeps; run manually, "
        + "ideally with -Djdk.tracePinnedThreads=full. Not a CI regression test.")
    void manyIdleConnections_carrierThreadUsage() throws Exception {
        ServerSocketChannel server = ServerSocketChannel.open();
        server.bind(new InetSocketAddress("localhost", 0), 256);
        int port = ((InetSocketAddress) server.getLocalAddress()).getPort();

        List<SocketChannel> serverSide = new ArrayList<>();
        volatileFlag.running = true;
        Thread acceptThread = new Thread(() -> {
            try {
                while (volatileFlag.running) {
                    SocketChannel s = server.accept();
                    synchronized (serverSide) {
                        serverSide.add(s);
                    }
                }
            } catch (Exception ignored) {
                // server closed
            }
        }, "scaling-accept");
        acceptThread.setDaemon(true);
        acceptThread.start();

        List<TcpTransportInstance> clients = new ArrayList<>();
        try {
            for (int i = 0; i < CONNECTIONS; i++) {
                TcpTransportConfiguration config = new TcpTransportConfiguration();
                config.receiveBufferSize = 81920;
                config.connectTimeout = 5000;
                clients.add(new TcpTransportInstance(
                    new InetSocketAddress("localhost", port), config, AuditLog.builder().build()));
            }

            // Let every connection's read virtual thread settle into its blocking wait.
            Thread.sleep(3000);

            // One snapshot so carriers and total are derived from the same instant.
            Set<Thread> threads = Thread.getAllStackTraces().keySet();
            // Count only virtual-thread scheduler carriers, not ForkJoinPool.commonPool
            // workers (parallel streams etc.), which would inflate the carrier count.
            long carriers = threads.stream()
                .filter(t -> !t.isVirtual())
                .filter(t -> t.getName().contains("ForkJoinPool"))
                .filter(t -> !t.getName().contains("commonPool"))
                .count();
            long total = threads.size();
            System.out.println("CARRIER_COUNT=" + carriers
                + " TOTAL_THREADS=" + total
                + " CONNECTIONS=" + CONNECTIONS
                + " CPUS=" + Runtime.getRuntime().availableProcessors());

            long openCount = clients.stream().filter(TcpTransportInstance::isOpen).count();
            assertTrue(openCount >= CONNECTIONS - 5,
                "expected ~all connections open, got " + openCount + "/" + CONNECTIONS);
        } finally {
            for (TcpTransportInstance c : clients) {
                try { c.close(); } catch (Exception ignored) { }
            }
            volatileFlag.running = false;
            synchronized (serverSide) {
                for (SocketChannel s : serverSide) {
                    try { s.close(); } catch (Exception ignored) { }
                }
            }
            server.close();
        }
    }

    // tiny holder carrying a volatile stop flag the daemon accept loop polls to shut down
    private static final class Flag { volatile boolean running; }
    private final Flag volatileFlag = new Flag();
}
