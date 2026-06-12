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
package org.apache.plc4x.java.transport.serial.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class SerialTransportConfigurationTest {

    @Test
    void testInstantiation() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        assertNotNull(config);
    }

    @Test
    void testImplementsTransportConfiguration() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        assertTrue(true);
    }

    @Test
    void testBaudRateField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.baudRate = 115200;
        assertEquals(115200, config.baudRate);
    }

    @Test
    void testDataBitsField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.dataBits = 7;
        assertEquals(7, config.dataBits);
    }

    @Test
    void testStopBitsField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.stopBits = 2;
        assertEquals(2, config.stopBits);
    }

    @Test
    void testParityField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        assertNull(config.parity);

        config.parity = "ODD";
        assertEquals("ODD", config.parity);
    }

    @Test
    void testFlowControlField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        assertNull(config.flowControl);

        config.flowControl = "RTS_CTS";
        assertEquals("RTS_CTS", config.flowControl);
    }

    @Test
    void testReadTimeoutField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.readTimeout = 5000;
        assertEquals(5000, config.readTimeout);
    }

    @Test
    void testWriteTimeoutField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.writeTimeout = 3000;
        assertEquals(3000, config.writeTimeout);
    }

    @Test
    void testDtrField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.dtr = false;
        assertFalse(config.dtr);

        config.dtr = true;
        assertTrue(config.dtr);
    }

    @Test
    void testRtsField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.rts = false;
        assertFalse(config.rts);

        config.rts = true;
        assertTrue(config.rts);
    }

    @Test
    void testReusePortField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.reusePort = false;
        assertFalse(config.reusePort);

        config.reusePort = true;
        assertTrue(config.reusePort);
    }

    @Test
    void testReceiveBufferSizeField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.receiveBufferSize = 65536;
        assertEquals(65536, config.receiveBufferSize);
    }

    @Test
    void testSendBufferSizeField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.sendBufferSize = 32768;
        assertEquals(32768, config.sendBufferSize);
    }

    @Test
    void testBreakEnabledField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.breakEnabled = false;
        assertFalse(config.breakEnabled);

        config.breakEnabled = true;
        assertTrue(config.breakEnabled);
    }

    @Test
    void testInterframeDelayField() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.interframeDelay = 50;
        assertEquals(50, config.interframeDelay);
    }

    @Test
    void testCompleteConfiguration() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();

        config.baudRate = 115200;
        config.dataBits = 8;
        config.stopBits = 1;
        config.parity = "EVEN";
        config.flowControl = "XON_XOFF";
        config.readTimeout = 2000;
        config.writeTimeout = 1500;
        config.dtr = true;
        config.rts = true;
        config.reusePort = true;
        config.receiveBufferSize = 8192;
        config.sendBufferSize = 8192;
        config.breakEnabled = true;
        config.interframeDelay = 10;

        assertEquals(115200, config.baudRate);
        assertEquals(8, config.dataBits);
        assertEquals(1, config.stopBits);
        assertEquals("EVEN", config.parity);
        assertEquals("XON_XOFF", config.flowControl);
        assertEquals(2000, config.readTimeout);
        assertEquals(1500, config.writeTimeout);
        assertTrue(config.dtr);
        assertTrue(config.rts);
        assertTrue(config.reusePort);
        assertEquals(8192, config.receiveBufferSize);
        assertEquals(8192, config.sendBufferSize);
        assertTrue(config.breakEnabled);
        assertEquals(10, config.interframeDelay);
    }

    @Test
    void testZeroTimeouts() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.readTimeout = 0;
        config.writeTimeout = 0;

        assertEquals(0, config.readTimeout);
        assertEquals(0, config.writeTimeout);
    }

    @Test
    void testZeroBufferSizes() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.sendBufferSize = 0;
        config.receiveBufferSize = 0;

        assertEquals(0, config.sendBufferSize);
        assertEquals(0, config.receiveBufferSize);
    }

    @Test
    void testZeroInterframeDelay() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.interframeDelay = 0;
        assertEquals(0, config.interframeDelay);
    }
}
