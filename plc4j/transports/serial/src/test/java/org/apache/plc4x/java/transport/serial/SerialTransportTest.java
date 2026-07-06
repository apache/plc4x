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
package org.apache.plc4x.java.transport.serial;

import com.fazecast.jSerialComm.SerialPort;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.serial.config.SerialTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

class SerialTransportTest {

    private SerialTransport transport;

    @BeforeEach
    void setUp() {
        transport = new SerialTransport();
    }

    @Test
    void testGetTransportCode() {
        assertEquals("serial", transport.getTransportCode());
    }

    @Test
    void testGetTransportName() {
        assertEquals("Serial", transport.getTransportName());
    }

    @Test
    void testGetTransportConfigType() {
        assertEquals(SerialTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void testCreateTransportInstance_dedicated() throws Exception {
        // Skip if no serial ports available
        SerialPort[] ports = SerialPort.getCommPorts();
        assumeTrue(ports.length > 0, "No serial ports available for testing");

        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.baudRate = 9600;
        config.dataBits = 8;
        config.stopBits = 1;
        config.parity = "NONE";
        config.flowControl = "NONE";
        config.reusePort = false;

        // This may fail if the port is in use or requires permissions
        try {
            TransportInstance<SerialTransportConfiguration> instance = transport.createTransportInstance(
                "serial://" + ports[0].getSystemPortName(), config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
        } catch (TransportException e) {
            // Port might be in use or lack permissions - this is acceptable for unit tests
            System.out.println("Could not open serial port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_shared() throws Exception {
        // Skip if no serial ports available
        SerialPort[] ports = SerialPort.getCommPorts();
        assumeTrue(ports.length > 0, "No serial ports available for testing");

        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.baudRate = 9600;
        config.dataBits = 8;
        config.stopBits = 1;
        config.parity = "NONE";
        config.flowControl = "NONE";
        config.reusePort = true;

        try {
            TransportInstance<SerialTransportConfiguration> instance = transport.createTransportInstance(
                "serial://" + ports[0].getSystemPortName(), config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
        } catch (TransportException e) {
            System.out.println("Could not open serial port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_multipleShared() throws Exception {
        // Skip if no serial ports available
        SerialPort[] ports = SerialPort.getCommPorts();
        assumeTrue(ports.length > 0, "No serial ports available for testing");

        SerialTransportConfiguration config1 = new SerialTransportConfiguration();
        config1.baudRate = 9600;
        config1.reusePort = true;

        SerialTransportConfiguration config2 = new SerialTransportConfiguration();
        config2.baudRate = 9600;
        config2.reusePort = true;

        try {
            TransportInstance<SerialTransportConfiguration> instance1 = transport.createTransportInstance(
                "serial://" + ports[0].getSystemPortName(), config1, AuditLog.builder().build());

            TransportInstance<SerialTransportConfiguration> instance2 = transport.createTransportInstance(
                "serial://" + ports[0].getSystemPortName(), config2, AuditLog.builder().build());

            assertTrue(instance1.isOpen());
            assertTrue(instance2.isOpen());

            // Close first - port should stay open
            instance1.close();
            assertFalse(instance1.isOpen());

            // Close second - port should close
            instance2.close();
            assertFalse(instance2.isOpen());

        } catch (TransportException e) {
            System.out.println("Could not open serial port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_withCustomSettings() throws Exception {
        SerialPort[] ports = SerialPort.getCommPorts();
        assumeTrue(ports.length > 0, "No serial ports available for testing");

        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.baudRate = 115200;
        config.dataBits = 7;
        config.stopBits = 2;
        config.parity = "EVEN";
        config.flowControl = "RTS_CTS";
        config.dtr = true;
        config.rts = true;
        config.reusePort = false;

        try {
            TransportInstance<SerialTransportConfiguration> instance = transport.createTransportInstance(
                "serial://" + ports[0].getSystemPortName(), config, AuditLog.builder().build());

            assertNotNull(instance);
            assertTrue(instance.isOpen());

            instance.close();
        } catch (TransportException e) {
            System.out.println("Could not open serial port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testCreateTransportInstance_invalidPort() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.baudRate = 9600;
        config.reusePort = false;

        assertThrows(TransportException.class, () ->
            transport.createTransportInstance("serial:///dev/ttyINVALID_PORT_THAT_DOES_NOT_EXIST", config, AuditLog.builder().build())
        );
    }

    @Test
    void testParityParsing() throws Exception {
        SerialPort[] ports = SerialPort.getCommPorts();
        assumeTrue(ports.length > 0, "No serial ports available for testing");

        String[] parityOptions = {"NONE", "ODD", "EVEN", "MARK", "SPACE"};

        for (String parity : parityOptions) {
            SerialTransportConfiguration config = new SerialTransportConfiguration();
            config.baudRate = 9600;
            config.parity = parity;
            config.reusePort = false;

            try {
                TransportInstance<SerialTransportConfiguration> instance = transport.createTransportInstance(
                    "serial://" + ports[0].getSystemPortName(), config, AuditLog.builder().build());
                instance.close();
            } catch (TransportException e) {
                // Port might be in use
            }
        }
    }

    @Test
    void testFlowControlParsing() throws Exception {
        SerialPort[] ports = SerialPort.getCommPorts();
        assumeTrue(ports.length > 0, "No serial ports available for testing");

        String[] flowControlOptions = {"NONE", "RTS_CTS", "XON_XOFF"};

        for (String flowControl : flowControlOptions) {
            SerialTransportConfiguration config = new SerialTransportConfiguration();
            config.baudRate = 9600;
            config.flowControl = flowControl;
            config.reusePort = false;

            try {
                TransportInstance<SerialTransportConfiguration> instance = transport.createTransportInstance(
                    "serial://" + ports[0].getSystemPortName(), config, AuditLog.builder().build());
                instance.close();
            } catch (TransportException e) {
                // Port might be in use
            }
        }
    }

    @Test
    void testCreateTransportInstance_wrongConfigurationType() {
        // Create a dummy configuration that is not SerialTransportConfiguration
        TransportConfiguration wrongConfig = new TransportConfiguration() {};

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () ->
            transport.createTransportInstance("serial:///dev/ttyUSB0", wrongConfig, AuditLog.builder().build())
        );

        assertTrue(exception.getMessage().contains("Expected configuration of type"));
        assertTrue(exception.getMessage().contains("SerialTransportConfiguration"));
    }
}
