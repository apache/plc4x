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

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.channels.ServerSocketChannel;

import static org.junit.jupiter.api.Assertions.*;

class TcpTransportTest {

    private TcpTransport transport;
    private ServerSocketChannel serverChannel;
    private int serverPort;

    @BeforeEach
    void setUp() throws IOException {
        transport = new TcpTransport();

        // Start a simple TCP server for testing
        serverChannel = ServerSocketChannel.open();
        serverChannel.bind(new InetSocketAddress("localhost", 0));
        serverPort = ((InetSocketAddress) serverChannel.getLocalAddress()).getPort();
    }

    @AfterEach
    void tearDown() throws IOException {
        if (serverChannel != null && serverChannel.isOpen()) {
            serverChannel.close();
        }
    }

    @Test
    void testGetTransportCode() {
        assertEquals("tcp", transport.getTransportCode());
    }

    @Test
    void testGetTransportName() {
        assertEquals("TCP", transport.getTransportName());
    }

    @Test
    void testGetTransportConfigType() {
        assertEquals(TcpTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void testCreateTransportInstance_successful() throws Exception {
        // Accept connection in the background
        Thread serverThread = new Thread(() -> {
            try {
                serverChannel.accept();
            } catch (IOException e) {
                // Expected when test completes
            }
        });
        serverThread.start();

        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 5000;
        config.receiveBufferSize = 81920;

        TransportInstance<TcpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
        serverThread.join(1000);
    }

    @Test
    void testCreateTransportInstance_connectionRefused() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("localhost:1", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withInvalidHost() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 1000;
        config.receiveBufferSize = 81920;

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance(
                "invalid.host.that.does.not.exist.local:80", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withLocalBinding() throws Exception {
        // Accept connection in the background
        Thread serverThread = new Thread(() -> {
            try {
                serverChannel.accept();
            } catch (IOException e) {
                // Expected when test completes
            }
        });
        serverThread.start();

        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 5000;
        config.receiveBufferSize = 81920;
        config.localAddress = "localhost";
        config.localPort = 0; // Use ephemeral port

        TransportInstance<TcpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
        serverThread.join(1000);
    }

    @Test
    void testCreateTransportInstance_withCustomSocketOptions() throws Exception {
        // Accept connection in the background
        Thread serverThread = new Thread(() -> {
            try {
                serverChannel.accept();
            } catch (IOException e) {
                // Expected when test completes
            }
        });
        serverThread.start();

        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 5000;
        config.receiveBufferSize = 81920;
        config.tcpNoDelay = true;
        config.keepAlive = true;
        config.sendBufferSize = 16384;
        config.receiveBufferSize = 16384;

        TransportInstance<TcpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
        serverThread.join(1000);
    }

    @Test
    void testCreateTransportInstance_withInvalidUrl() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();

        // Invalid URL format (no host/IP)
        assertThrows(Exception.class, () ->
            transport.createTransportInstance(":", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withoutPort() {
        TcpTransportConfiguration config = new TcpTransportConfiguration();

        // No port in URL and no default port configured
        assertThrows(Exception.class, () ->
            transport.createTransportInstance("localhost", config, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withWrongConfigurationType() {
        TransportConfiguration wrongConfig = new TransportConfiguration() {};

        // Wrong configuration type should throw IllegalArgumentException
        assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("localhost:8080", wrongConfig, AuditLog.builder().build())
        );
    }

    @Test
    void testCreateTransportInstance_withIpAddress() throws Exception {
        // Accept connection in the background
        Thread serverThread = new Thread(() -> {
            try {
                serverChannel.accept();
            } catch (IOException e) {
                // Expected when test completes
            }
        });
        serverThread.start();

        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 5000;
        config.receiveBufferSize = 81920;

        // Test with IP address format
        TransportInstance<TcpTransportConfiguration> instance = transport.createTransportInstance(
            "127.0.0.1:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
        serverThread.join(1000);
    }

    @Test
    void testCreateTransportInstance_withHostnameAndNoPort() {
        TcpTransportConfiguration config = new TcpTransportConfiguration() {
            @Override
            public int getDefaultPort() {
                return serverPort; // Provide default port
            }
        };
        config.receiveBufferSize = 81920;

        // Accept connection in the background
        Thread serverThread = new Thread(() -> {
            try {
                serverChannel.accept();
            } catch (IOException e) {
                // Expected when test completes
            }
        });
        serverThread.start();

        try {
            // Should use the default port from configuration
            TransportInstance<TcpTransportConfiguration> instance = transport.createTransportInstance(
                "localhost", config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
            serverThread.join(1000);
        } catch (Exception e) {
            fail("Should use default port from configuration: " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_urlWithExtraParameters() throws Exception {
        // Accept connection in the background
        Thread serverThread = new Thread(() -> {
            try {
                serverChannel.accept();
            } catch (IOException e) {
                // Expected when test completes
            }
        });
        serverThread.start();

        TcpTransportConfiguration config = new TcpTransportConfiguration();
        config.connectTimeout = 5000;
        config.receiveBufferSize = 81920;

        // URL with extra query parameters should still work
        TransportInstance<TcpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort + "?param=value", config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
        serverThread.join(1000);
    }
}

