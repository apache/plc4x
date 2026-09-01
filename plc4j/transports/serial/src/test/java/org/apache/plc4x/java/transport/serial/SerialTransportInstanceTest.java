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
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.serial.config.SerialTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Instant;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

/**
 * Tests for SerialTransportInstance.
 * Note: These tests require actual serial ports to run. They use assumptions to skip
 * when no ports are available or when ports cannot be opened.
 */
class SerialTransportInstanceTest {

    private Process process;
    private SerialPort port;
    private SerialTransportInstance transportInstance;
    private SerialTransportConfiguration config;

    private Path masterLink;
    private Path slaveLink;
    /** Path of the pty this test talks to; unique per test, see {@link #setUp()}. */
    private String portPath;

    @BeforeEach
    void setUp() throws Exception {
        boolean socatAvailable = false;
        if (!System.getProperty("os.name").toLowerCase().contains("win")) {
            try {
                Process check = new ProcessBuilder("which", "socat").start();
                socatAvailable = (check.waitFor() == 0);
            } catch (IOException e) {
                socatAvailable = false;
            }
        }
        assumeTrue(socatAvailable, "socat is not installed - skipping serial transport integration tests (install socat to enable)");

        // Every test gets its own pty pair. These used to be two fixed paths shared by all of
        // them, and since tearDown did not wait for socat to exit, one still shutting down could
        // remove the link the next test's socat had just created - which showed up as rare
        // "Failed to open serial port" errors.
        String unique = "plc4x-serial-test-" + UUID.randomUUID();
        Path tempDir = Path.of(System.getProperty("java.io.tmpdir"));
        masterLink = tempDir.resolve(unique + "-a");
        slaveLink = tempDir.resolve(unique + "-b");
        portPath = masterLink.toString();

        process = new ProcessBuilder("socat", "-d", "-d",
            "pty,raw,echo=0,link=" + masterLink,
            "pty,raw,echo=0,link=" + slaveLink).start();
        awaitPtyLinks();

        config = new SerialTransportConfiguration();
        config.baudRate = 9600;
        config.dataBits = 8;
        config.stopBits = 1;
        config.parity = "NONE";
        config.flowControl = "NONE";

        transportInstance = new SerialTransportInstance(new SharedSerialPortManager(), portPath, config, AuditLog.builder().build());
    }

    @AfterEach
    void tearDown() throws Exception {
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }

        if (process != null) {
            // Wait for socat to really be gone: it removes its links while shutting down, so
            // returning before that finishes is what let one test disturb the next.
            process.destroy();
            if (!process.waitFor(5, TimeUnit.SECONDS)) {
                process.destroyForcibly().waitFor(5, TimeUnit.SECONDS);
            }
        }
        // socat normally removes these itself; clean up in case it was killed before it could.
        if (masterLink != null) {
            Files.deleteIfExists(masterLink);
        }
        if (slaveLink != null) {
            Files.deleteIfExists(slaveLink);
        }
    }

    /**
     * Waits until socat has published both pty links, rather than assuming a fixed pause is long
     * enough. Fails with the reason rather than letting the port open fail further down.
     */
    private void awaitPtyLinks() throws Exception {
        Instant deadline = Instant.now().plusSeconds(10);
        while (Instant.now().isBefore(deadline)) {
            if (Files.exists(masterLink) && Files.exists(slaveLink)) {
                return;
            }
            if (!process.isAlive()) {
                fail("socat exited with code " + process.exitValue() + " before creating the pty links");
            }
            Thread.sleep(20);
        }
        fail("socat did not create the pty links " + masterLink + " and " + slaveLink + " within 10s");
    }

    @Test
    void testGetConfiguration() {
        SerialTransportConfiguration config = transportInstance.getConfiguration();
        assertNotNull(config);
    }

    @Test
    void testIsOpen_whenConnected() {
        assertTrue(transportInstance.isOpen());
    }

    @Test
    void testIsOpen_whenClosed() throws TransportException {
        transportInstance.close();
        assertFalse(transportInstance.isOpen());
    }

    @Test
    void testGetNumBytesAvailable() throws TransportException {
        // Initially should return 0 or small number
        int available = transportInstance.getNumBytesAvailable();
        assertTrue(available >= 0);
    }

    @Test
    void testWrite_successful() throws TransportException {
        byte[] data = "TEST".getBytes();

        // Should not throw
        assertDoesNotThrow(() -> transportInstance.write(data));
    }

    @Test
    void testWrite_emptyArray() throws TransportException {
        assertDoesNotThrow(() -> transportInstance.write(new byte[0]));
    }

    @Test
    void testWrite_nullArray() throws TransportException {
        assertDoesNotThrow(() -> transportInstance.write(null));
    }

    @Test
    void testRead_zeroBytes() throws TransportException {
        byte[] result = transportInstance.read(0);
        assertEquals(0, result.length);
    }

    @Test
    void testRead_whenClosed_throwsException() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.read(10)
        );
    }

    @Test
    void testWrite_whenClosed_throwsException() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.write("test".getBytes())
        );
    }

    @Test
    void testClose_idempotent() throws TransportException {
        transportInstance.close();
        assertFalse(transportInstance.isOpen());

        // Should not throw
        assertDoesNotThrow(() -> transportInstance.close());
        assertFalse(transportInstance.isOpen());
    }

    @Test
    void testPeekReadableBytes_emptyBuffer() throws TransportException {
        // When no data available, should throw
        assertThrows(TransportException.class, () ->
            transportInstance.peekReadableBytes(10)
        );
    }

    @Test
    void testWrite_largeData() throws TransportException {
        byte[] data = new byte[1024];
        for (int i = 0; i < data.length; i++) {
            data[i] = (byte) (i % 256);
        }

        // Should not throw
        assertDoesNotThrow(() -> transportInstance.write(data));
    }

    /**
     * Test with loopback cable if available.
     * This test is disabled by default as it requires hardware setup.
     */
    @Test
    @Disabled("Requires loopback cable - enable manually for hardware testing")
    void testLoopback_writeAndRead() throws Exception {
        byte[] testData = "LOOPBACK TEST".getBytes();

        // Write data
        transportInstance.write(testData);

        // Wait for loopback
        Thread.sleep(100);

        // Read data back
        byte[] received = transportInstance.read(testData.length);

        assertArrayEquals(testData, received);
    }

    /**
     * Test peek functionality with loopback cable.
     * This test is disabled by default as it requires hardware setup.
     */
    @Test
    @Disabled("Requires loopback cable - enable manually for hardware testing")
    void testLoopback_peekAndRead() throws Exception {
        byte[] testData = "PEEK TEST".getBytes();

        // Write data
        transportInstance.write(testData);

        // Wait for loopback
        Thread.sleep(100);

        // Peek data
        byte[] peeked = transportInstance.peekReadableBytes(testData.length);
        assertArrayEquals(testData, peeked);

        // Read data - should still be available
        byte[] read = transportInstance.read(testData.length);
        assertArrayEquals(testData, read);
    }

    @Test
    void testRegisterDataListener() {
        // Test that registering a data listener doesn't throw
        Runnable listener = () -> {};
        assertDoesNotThrow(() -> transportInstance.registerDataListener(listener));
    }

    @Test
    void testRemoveDataListener() {
        // Register first then remove
        Runnable listener = () -> {};
        transportInstance.registerDataListener(listener);

        // Test that removing a data listener doesn't throw
        assertDoesNotThrow(() -> transportInstance.removeDataListener());
    }

    @Test
    void testPeekReadableBytes_zeroBytes() throws TransportException {
        // Requesting 0 bytes should return empty array
        byte[] result = transportInstance.peekReadableBytes(0);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekReadableBytes_negativeBytes() throws TransportException {
        // Requesting negative bytes should return empty array
        byte[] result = transportInstance.peekReadableBytes(-1);
        assertEquals(0, result.length);
    }

    @Test
    void testRead_negativeBytes() throws TransportException {
        // Requesting negative bytes should return empty array
        byte[] result = transportInstance.read(-5);
        assertEquals(0, result.length);
    }

    @Test
    void testGetNumBytesAvailable_whenClosed() throws TransportException {
        transportInstance.close();

        // Should return 0 when closed, not throw
        int available = transportInstance.getNumBytesAvailable();
        assertEquals(0, available);
    }

    @Test
    void testPeekReadableBytes_whenClosed_throwsException() throws TransportException {
        transportInstance.close();

        assertThrows(TransportException.class, () ->
            transportInstance.peekReadableBytes(10)
        );
    }

    @Test
    void testWithDtrAndRtsEnabled() throws Exception {
        // Clean up current instance first
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }

        // Create config with DTR and RTS enabled
        SerialTransportConfiguration configWithSignals = new SerialTransportConfiguration();
        configWithSignals.baudRate = 9600;
        configWithSignals.dataBits = 8;
        configWithSignals.stopBits = 1;
        configWithSignals.parity = "NONE";
        configWithSignals.flowControl = "NONE";
        configWithSignals.dtr = true;  // Enable DTR
        configWithSignals.rts = true;  // Enable RTS
        configWithSignals.reusePort = false;

        SerialTransportInstance instanceWithSignals = new SerialTransportInstance(
            new SharedSerialPortManager(), portPath, configWithSignals, AuditLog.builder().build());

        assertTrue(instanceWithSignals.isOpen());
        instanceWithSignals.close();
    }

    @Test
    void testWithUnknownParity() throws Exception {
        // Clean up current instance first
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }

        // Create config with unknown parity - should now fail fast
        SerialTransportConfiguration configWithUnknownParity = new SerialTransportConfiguration();
        configWithUnknownParity.baudRate = 9600;
        configWithUnknownParity.dataBits = 8;
        configWithUnknownParity.stopBits = 1;
        configWithUnknownParity.parity = "UNKNOWN_PARITY";  // Unknown parity
        configWithUnknownParity.flowControl = "NONE";
        configWithUnknownParity.reusePort = false;

        assertThrows(TransportException.class, () -> new SerialTransportInstance(
            new SharedSerialPortManager(), portPath, configWithUnknownParity, AuditLog.builder().build()));
    }

    @Test
    void testWithUnknownFlowControl() throws Exception {
        // Clean up current instance first
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }

        // Create config with unknown flow control - should now fail fast
        SerialTransportConfiguration configWithUnknownFlowControl = new SerialTransportConfiguration();
        configWithUnknownFlowControl.baudRate = 9600;
        configWithUnknownFlowControl.dataBits = 8;
        configWithUnknownFlowControl.stopBits = 1;
        configWithUnknownFlowControl.parity = "NONE";
        configWithUnknownFlowControl.flowControl = "UNKNOWN_FLOW_CONTROL";  // Unknown flow control
        configWithUnknownFlowControl.reusePort = false;

        assertThrows(TransportException.class, () -> new SerialTransportInstance(
            new SharedSerialPortManager(), portPath, configWithUnknownFlowControl, AuditLog.builder().build()));
    }

    @Test
    void testWithXonXoffFlowControl() throws Exception {
        // Clean up current instance first
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }

        // Create config with XON/XOFF flow control
        SerialTransportConfiguration configWithXonXoff = new SerialTransportConfiguration();
        configWithXonXoff.baudRate = 9600;
        configWithXonXoff.dataBits = 8;
        configWithXonXoff.stopBits = 1;
        configWithXonXoff.parity = "NONE";
        configWithXonXoff.flowControl = "XON_XOFF";
        configWithXonXoff.reusePort = false;

        SerialTransportInstance instanceWithXonXoff = new SerialTransportInstance(
            new SharedSerialPortManager(), portPath, configWithXonXoff, AuditLog.builder().build());

        assertTrue(instanceWithXonXoff.isOpen());
        instanceWithXonXoff.close();
    }

    @Test
    void testWithRtsCtsXonXoffFlowControl() throws Exception {
        // Clean up current instance first
        if (transportInstance != null && transportInstance.isOpen()) {
            transportInstance.close();
        }

        // Create config with RTS_CTS_XON_XOFF flow control - combined mode is no longer supported
        SerialTransportConfiguration configWithCombinedFlowControl = new SerialTransportConfiguration();
        configWithCombinedFlowControl.baudRate = 9600;
        configWithCombinedFlowControl.dataBits = 8;
        configWithCombinedFlowControl.stopBits = 1;
        configWithCombinedFlowControl.parity = "NONE";
        configWithCombinedFlowControl.flowControl = "RTS_CTS_XON_XOFF";
        configWithCombinedFlowControl.reusePort = false;

        assertThrows(TransportException.class, () -> new SerialTransportInstance(
            new SharedSerialPortManager(), portPath, configWithCombinedFlowControl, AuditLog.builder().build()));
    }

    @Test
    void testWithAllParityOptions() throws Exception {
        String[] parityOptions = {"ODD", "EVEN", "MARK", "SPACE"};

        for (String parity : parityOptions) {
            // Clean up current instance first
            if (transportInstance != null && transportInstance.isOpen()) {
                transportInstance.close();
            }

            SerialTransportConfiguration configWithParity = new SerialTransportConfiguration();
            configWithParity.baudRate = 9600;
            configWithParity.dataBits = 8;
            configWithParity.stopBits = 1;
            configWithParity.parity = parity;
            configWithParity.flowControl = "NONE";
            configWithParity.reusePort = false;

            SerialTransportInstance instanceWithParity = new SerialTransportInstance(
                new SharedSerialPortManager(), portPath, configWithParity, AuditLog.builder().build());

            assertTrue(instanceWithParity.isOpen(), "Instance with parity " + parity + " should be open");
            instanceWithParity.close();
        }
    }

    @Test
    void parseParityAcceptsCaseInsensitiveForms() throws Exception {
        assertEquals(SerialPort.NO_PARITY, SerialTransportInstance.parseParity("none"));
        assertEquals(SerialPort.EVEN_PARITY, SerialTransportInstance.parseParity("even"));
        assertEquals(SerialPort.EVEN_PARITY, SerialTransportInstance.parseParity("EVEN"));
        assertEquals(SerialPort.ODD_PARITY, SerialTransportInstance.parseParity("Odd"));
        assertEquals(SerialPort.MARK_PARITY, SerialTransportInstance.parseParity("MARK"));
        assertEquals(SerialPort.SPACE_PARITY, SerialTransportInstance.parseParity("space"));
    }

    @Test
    void parseParityRejectsUnknownValues() {
        TransportException e = assertThrows(TransportException.class,
            () -> SerialTransportInstance.parseParity("strong"));
        assertTrue(e.getMessage().contains("parity"));
        assertTrue(e.getMessage().contains("strong"));
    }

    @Test
    void parseFlowControlAcceptsCanonicalAndLegacyForms() throws Exception {
        int rtsCts = SerialPort.FLOW_CONTROL_RTS_ENABLED | SerialPort.FLOW_CONTROL_CTS_ENABLED;
        assertEquals(rtsCts, SerialTransportInstance.parseFlowControl("rts-cts"));
        assertEquals(rtsCts, SerialTransportInstance.parseFlowControl("RTS_CTS"));
        assertEquals(rtsCts, SerialTransportInstance.parseFlowControl("RTSCTS"));
        int xonXoff = SerialPort.FLOW_CONTROL_XONXOFF_IN_ENABLED | SerialPort.FLOW_CONTROL_XONXOFF_OUT_ENABLED;
        assertEquals(xonXoff, SerialTransportInstance.parseFlowControl("xon-xoff"));
        assertEquals(xonXoff, SerialTransportInstance.parseFlowControl("XON_XOFF"));
        assertEquals(SerialPort.FLOW_CONTROL_DISABLED, SerialTransportInstance.parseFlowControl("NONE"));
    }

    @Test
    void parseFlowControlRejectsCombinedMode() {
        TransportException e = assertThrows(TransportException.class,
            () -> SerialTransportInstance.parseFlowControl("RTS_CTS_XON_XOFF"));
        assertTrue(e.getMessage().contains("flow-control"));
    }

    @Test
    void parseFlowControlRejectsUnknownValues() {
        assertThrows(TransportException.class,
            () -> SerialTransportInstance.parseFlowControl("magic"));
    }
}
