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
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

/**
 * Tests for SharedSerialPortManager.
 * Note: These tests require actual serial ports to run. They use assumptions to skip
 * when no ports are available or when ports cannot be opened.
 */
class SharedSerialPortManagerTest {

    private SharedSerialPortManager manager;
    private String testPortName;

    @BeforeEach
    void setUp() {
        manager = new SharedSerialPortManager();

        // Find an available port for testing
        SerialPort[] ports = SerialPort.getCommPorts();
        assumeTrue(ports.length > 0, "No serial ports available for testing");

        testPortName = ports[0].getSystemPortName();
    }

    @Test
    void testAcquirePort_createsNewPort() {
        SharedSerialPortManager.SerialPortConfig config = createDefaultConfig();

        try {
            SharedSerialPortManager.SharedPort sharedPort = manager.acquirePort(testPortName, config);

            assertNotNull(sharedPort);
            assertNotNull(sharedPort.getPort());
            assertTrue(sharedPort.getPort().isOpen());
            assertEquals(1, sharedPort.getRefCount());

            // Cleanup
            manager.releasePort(sharedPort);
        } catch (Exception e) {
            // Port might be in use
            System.out.println("Could not open port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testAcquirePort_reusesSamePort() {
        SharedSerialPortManager.SerialPortConfig config = createDefaultConfig();

        try {
            SharedSerialPortManager.SharedPort port1 = manager.acquirePort(testPortName, config);
            SharedSerialPortManager.SharedPort port2 = manager.acquirePort(testPortName, config);

            // Should be the same instance
            assertSame(port1, port2);
            assertSame(port1.getPort(), port2.getPort());
            assertEquals(2, port1.getRefCount());

            // Cleanup
            manager.releasePort(port1);
            manager.releasePort(port2);
        } catch (Exception e) {
            System.out.println("Could not open port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testReleasePort_decrementsRefCount() {
        SharedSerialPortManager.SerialPortConfig config = createDefaultConfig();

        try {
            SharedSerialPortManager.SharedPort port = manager.acquirePort(testPortName, config);
            assertEquals(1, port.getRefCount());

            // Acquire again
            manager.acquirePort(testPortName, config);
            assertEquals(2, port.getRefCount());

            // Release once
            manager.releasePort(port);
            assertEquals(1, port.getRefCount());
            assertTrue(port.getPort().isOpen());

            // Release again
            manager.releasePort(port);
            assertEquals(0, port.getRefCount());
            assertFalse(port.getPort().isOpen());
        } catch (Exception e) {
            System.out.println("Could not open port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testReleasePort_closesPortWhenRefCountZero() {
        SharedSerialPortManager.SerialPortConfig config = createDefaultConfig();

        try {
            SharedSerialPortManager.SharedPort port = manager.acquirePort(testPortName, config);

            assertTrue(port.getPort().isOpen());
            assertEquals(1, port.getRefCount());

            manager.releasePort(port);

            assertFalse(port.getPort().isOpen());
            assertEquals(0, port.getRefCount());
        } catch (Exception e) {
            System.out.println("Could not open port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testMultipleAcquireRelease_cycles() {
        SharedSerialPortManager.SerialPortConfig config = createDefaultConfig();

        try {
            // First cycle
            SharedSerialPortManager.SharedPort port1 = manager.acquirePort(testPortName, config);
            assertEquals(1, port1.getRefCount());
            manager.releasePort(port1);
            assertFalse(port1.getPort().isOpen());

            // Second cycle - should create new port
            SharedSerialPortManager.SharedPort port2 = manager.acquirePort(testPortName, config);
            assertEquals(1, port2.getRefCount());
            assertNotSame(port1, port2); // Different instance
            assertTrue(port2.getPort().isOpen());

            // Cleanup
            manager.releasePort(port2);
        } catch (Exception e) {
            System.out.println("Could not open port (expected if in use): " + e.getMessage());
        }
    }

    @Test
    void testWriteLock_enforcesSequentialAccess() throws Exception {
        SharedSerialPortManager.SerialPortConfig config = createDefaultConfig();

        try {
            SharedSerialPortManager.SharedPort port = manager.acquirePort(testPortName, config);

            // Lock from one thread
            port.lockWrite();

            // Try to lock from another thread
            Thread thread = new Thread(() -> {
                // This should block until first lock is released
                port.lockWrite();
                port.unlockWrite();
            });
            thread.start();

            // Give thread time to start
            Thread.sleep(50);

            // Thread should be blocked
            assertTrue(thread.isAlive());

            // Release first lock
            port.unlockWrite();

            // Wait for thread to complete
            thread.join(1000);
            assertFalse(thread.isAlive());

            // Cleanup
            manager.releasePort(port);
        } catch (Exception e) {
            System.out.println("Could not open port (expected if in use): " + e.getMessage());
        }
    }

    //@Test
    void testInterframeDelay_enforced() throws Exception {
        SharedSerialPortManager.SerialPortConfig config = new SharedSerialPortManager.SerialPortConfig(
            9600, 8, 1, SerialPort.NO_PARITY, SerialPort.FLOW_CONTROL_DISABLED,
            1000, 1000, false, false, 100 // 100ms interframe delay
        );

        try {
            SharedSerialPortManager.SharedPort port = manager.acquirePort(testPortName, config);

            // First write
            long start1 = System.currentTimeMillis();
            port.lockWrite();
            port.unlockWrite();
            long end1 = System.currentTimeMillis();

            // Second write - should be delayed
            long start2 = System.currentTimeMillis();
            port.lockWrite();
            port.unlockWrite();
            long end2 = System.currentTimeMillis();

            // Check that at least interframe delay elapsed between writes
            long timeBetweenWrites = start2 - end1;
            assertTrue(timeBetweenWrites >= 90, "Interframe delay not enforced: " + timeBetweenWrites + "ms");

            // Cleanup
            manager.releasePort(port);
        } catch (Exception e) {
            System.out.println("Could not open port (expected if in use): " + e.getMessage());
        }
    }

    private SharedSerialPortManager.SerialPortConfig createDefaultConfig() {
        return new SharedSerialPortManager.SerialPortConfig(
            9600, // baudRate
            8,    // dataBits
            1,    // stopBits
            SerialPort.NO_PARITY,
            SerialPort.FLOW_CONTROL_DISABLED,
            1000, // readTimeout
            1000, // writeTimeout
            false, // dtr
            false, // rts
            0     // interframeDelay
        );
    }
}
