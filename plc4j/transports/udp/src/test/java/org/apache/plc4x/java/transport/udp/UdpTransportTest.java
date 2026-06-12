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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.udp.config.UdpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.channels.DatagramChannel;

import static org.junit.jupiter.api.Assertions.*;

class UdpTransportTest {

    private UdpTransport transport;
    private DatagramChannel testServer;
    private int serverPort;

    @BeforeEach
    void setUp() throws IOException {
        transport = new UdpTransport();

        // Start a test UDP server
        testServer = DatagramChannel.open();
        testServer.bind(new InetSocketAddress("localhost", 0));
        serverPort = ((InetSocketAddress) testServer.getLocalAddress()).getPort();
    }

    @AfterEach
    void tearDown() throws IOException {
        if (testServer != null && testServer.isOpen()) {
            testServer.close();
        }
    }

    @Test
    void testGetTransportCode() {
        assertEquals("udp", transport.getTransportCode());
    }

    @Test
    void testGetTransportName() {
        assertEquals("UDP", transport.getTransportName());
    }

    @Test
    void testGetTransportConfigType() {
        assertEquals(UdpTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void testCreateTransportInstance_dedicated() throws Exception {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.localPort = 0; // Ephemeral port
        config.shareSocket = false;

        TransportInstance<UdpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstance_shared() throws Exception {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.localAddress = "localhost";
        config.localPort = 0; // Will bind to ephemeral port
        config.shareSocket = true;

        TransportInstance<UdpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstance_multipleSharedInstances() throws Exception {
        // Use a specific local port for sharing
        int sharedLocalPort = findAvailablePort();

        UdpTransportConfiguration config1 = new UdpTransportConfiguration();
        config1.localAddress = "localhost";
        config1.localPort = sharedLocalPort;
        config1.shareSocket = true;

        UdpTransportConfiguration config2 = new UdpTransportConfiguration();
        config2.localAddress = "localhost";
        config2.localPort = sharedLocalPort;
        config2.shareSocket = true;

        // Create a first instance
        TransportInstance<UdpTransportConfiguration> instance1 = transport.createTransportInstance(
            "localhost:" + serverPort, config1, AuditLog.builder().build());

        // Create a second instance sharing the same socket
        TransportInstance<UdpTransportConfiguration> instance2 = transport.createTransportInstance(
            "localhost:" + serverPort, config2, AuditLog.builder().build());

        assertTrue(instance1.isOpen());
        assertTrue(instance2.isOpen());

        // Close the first instance - socket should stay open
        instance1.close();
        assertFalse(instance1.isOpen());

        // Close the second instance - socket should close now
        instance2.close();
        assertFalse(instance2.isOpen());
    }

    @Test
    void testCreateTransportInstance_withBroadcast() throws Exception {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.broadcast = true;
        config.shareSocket = false;

        TransportInstance<UdpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstance_withCustomBufferSizes() throws Exception {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.sendBufferSize = 32768;
        config.receiveBufferSize = 32768;
        config.shareSocket = false;

        TransportInstance<UdpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstance_withLocalBinding() throws Exception {
        UdpTransportConfiguration config = new UdpTransportConfiguration();
        config.localAddress = "localhost";
        config.localPort = 0; // Ephemeral
        config.shareSocket = false;

        TransportInstance<UdpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost:" + serverPort, config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstance_wrongConfigurationType() {
        TransportConfiguration wrongConfig = new TransportConfiguration() {};

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("localhost:" + serverPort, wrongConfig, AuditLog.builder().build()));

        assertTrue(exception.getMessage().contains(UdpTransportConfiguration.class.getSimpleName()));
    }

    @Test
    void testCreateTransportInstance_invalidUrl() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();

        assertThrows(PlcRuntimeException.class, () ->
            transport.createTransportInstance("://not-a-valid-url", config, AuditLog.builder().build()));
    }

    @Test
    void testCreateTransportInstance_noPortWithDefaultPort() throws Exception {
        UdpTransportConfiguration config = new UdpTransportConfiguration() {
            @Override
            public int getDefaultPort() {
                return serverPort;
            }
        };
        config.shareSocket = false;

        TransportInstance<UdpTransportConfiguration> instance = transport.createTransportInstance(
            "localhost", config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstance_noPortNoDefaultPort() {
        UdpTransportConfiguration config = new UdpTransportConfiguration();

        PlcRuntimeException exception = assertThrows(PlcRuntimeException.class, () ->
            transport.createTransportInstance("localhost", config, AuditLog.builder().build()));

        assertTrue(exception.getMessage().contains("No port defined"));
    }

    private int findAvailablePort() throws IOException {
        try (DatagramChannel temp = DatagramChannel.open()) {
            temp.bind(new InetSocketAddress("localhost", 0));
            return ((InetSocketAddress) temp.getLocalAddress()).getPort();
        }
    }
}
