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
package org.apache.plc4x.java.transport.can.virtualcan;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.can.virtualcan.config.VirtualCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for {@link VirtualCanTransport}: transport metadata, config type,
 * and instance creation (including error paths).
 */
class VirtualCanTransportTest {

    private VirtualCanTransport transport;

    @BeforeEach
    void setUp() {
        VirtualCanBusManager.reset();
        transport = new VirtualCanTransport();
    }

    @AfterEach
    void tearDown() {
        VirtualCanBusManager.reset();
    }

    @Test
    void testGetTransportCode() {
        assertEquals("can-virtualcan", transport.getTransportCode());
    }

    @Test
    void testGetTransportName() {
        assertEquals("Virtual CAN", transport.getTransportName());
    }

    @Test
    void testGetTransportConfigType() {
        assertEquals(VirtualCanTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void testCreateTransportInstance() throws Exception {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.busName = "test-bus";

        TransportInstance<VirtualCanTransportConfiguration> instance =
                transport.createTransportInstance("", config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());
        assertInstanceOf(VirtualCanTransportInstance.class, instance);

        instance.close();
    }

    @Test
    void testCreateTransportInstanceWithDefaultBus() throws Exception {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();

        TransportInstance<VirtualCanTransportConfiguration> instance =
                transport.createTransportInstance("can-virtualcan://default", config, AuditLog.builder().build());

        assertNotNull(instance);
        assertTrue(instance.isOpen());

        instance.close();
    }

    @Test
    void testCreateTransportInstanceWrongConfigType() {
        TransportConfiguration wrongConfig = new TransportConfiguration() {};

        assertThrows(IllegalArgumentException.class, () ->
                transport.createTransportInstance("can-virtualcan://test", wrongConfig, AuditLog.builder().build()));
    }

    @Test
    void testAddressSegmentNamesTheBus() throws Exception {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();

        TransportInstance<VirtualCanTransportConfiguration> instance =
                transport.createTransportInstance("addressed-bus", config, AuditLog.builder().build());

        assertEquals("addressed-bus", config.busName);

        instance.close();
    }

    @Test
    void testBusNameOptionUsedWhenAddressSegmentIsEmpty() throws Exception {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.busName = "configured-bus";

        TransportInstance<VirtualCanTransportConfiguration> instance =
                transport.createTransportInstance("", config, AuditLog.builder().build());

        assertEquals("configured-bus", config.busName);

        instance.close();
    }

    @Test
    void testAddressSegmentWinsOverBusNameOption() throws Exception {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.busName = "configured-bus";

        TransportInstance<VirtualCanTransportConfiguration> instance =
                transport.createTransportInstance("addressed-bus", config, AuditLog.builder().build());

        assertEquals("addressed-bus", config.busName);

        instance.close();
    }

    @Test
    void testInstancesAddressedToDifferentBusesAreIsolated() throws Exception {
        VirtualCanTransportInstance instanceA = (VirtualCanTransportInstance) transport.createTransportInstance(
                "bus-alpha", new VirtualCanTransportConfiguration(), AuditLog.builder().build());
        VirtualCanTransportInstance instanceB = (VirtualCanTransportInstance) transport.createTransportInstance(
                "bus-beta", new VirtualCanTransportConfiguration(), AuditLog.builder().build());

        try {
            instanceA.write(new byte[]{0x0A, 0x0B});

            // Before the address segment was honoured both of these landed on the bus named
            // "default" and this frame arrived at instanceB.
            assertEquals(0, instanceB.getNumBytesAvailable());
        } finally {
            instanceA.close();
            instanceB.close();
        }
    }

    @Test
    void testInstancesAddressedToTheSameBusExchangeFrames() throws Exception {
        VirtualCanTransportInstance sender = (VirtualCanTransportInstance) transport.createTransportInstance(
                "shared-bus", new VirtualCanTransportConfiguration(), AuditLog.builder().build());
        VirtualCanTransportInstance receiver = (VirtualCanTransportInstance) transport.createTransportInstance(
                "shared-bus", new VirtualCanTransportConfiguration(), AuditLog.builder().build());

        try {
            sender.write(new byte[]{0x0A, 0x0B});

            assertEquals(2, receiver.getNumBytesAvailable());
        } finally {
            sender.close();
            receiver.close();
        }
    }
}
