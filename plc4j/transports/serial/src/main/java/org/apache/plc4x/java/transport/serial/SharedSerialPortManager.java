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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

/**
 * Manages shared serial ports across multiple transport instances.
 * Allows multiple logical connections to share a single serial port,
 * which is essential for protocols like Modbus RTU that multiplex over one serial line.
 */
public class SharedSerialPortManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(SharedSerialPortManager.class);

    // Map of port name -> shared port wrapper
    private final Map<String, SharedPort> sharedPorts = new ConcurrentHashMap<>();

    public SharedSerialPortManager() {
        // Private constructor for singleton
    }

    /**
     * Acquires a shared serial port for the given port name and configuration.
     * If a port is already open for this name, returns the existing one.
     * Otherwise creates a new one with the provided configuration.
     */
    public SharedPort acquirePort(String portName, SerialPortConfig config) {
        return sharedPorts.compute(portName, (name, existing) -> {
            if (existing != null) {
                existing.incrementRefCount();
                LOGGER.debug("Reusing shared serial port {} (refCount={})",
                    portName, existing.getRefCount());
                return existing;
            } else {
                SerialPort port = SerialPort.getCommPort(portName);

                // Configure port
                port.setBaudRate(config.baudRate);
                port.setNumDataBits(config.dataBits);
                port.setNumStopBits(config.stopBits);
                port.setParity(config.parity);
                port.setFlowControl(config.flowControl);

                // Set timeouts
                port.setComPortTimeouts(
                    SerialPort.TIMEOUT_READ_SEMI_BLOCKING | SerialPort.TIMEOUT_WRITE_BLOCKING,
                    config.readTimeout,
                    config.writeTimeout
                );

                // Open port
                if (!port.openPort()) {
                    throw new RuntimeException("Failed to open serial port: " + portName);
                }

                // Set DTR/RTS
                if (config.dtr) {
                    port.setDTR();
                } else {
                    port.clearDTR();
                }

                if (config.rts) {
                    port.setRTS();
                } else {
                    port.clearRTS();
                }

                LOGGER.info("Opened shared serial port {} at {} baud", portName, config.baudRate);
                return new SharedPort(port, portName, config.interframeDelay);
            }
        });
    }

    /**
     * Releases a shared serial port reference.
     * If this is the last reference, closes the port.
     */
    public void releasePort(SharedPort port) {
        String portName = port.getPortName();

        sharedPorts.compute(portName, (name, existing) -> {
            if (existing == null) {
                LOGGER.warn("Attempted to release non-existent port: {}", portName);
                return null;
            }

            int refCount = existing.decrementRefCount();
            LOGGER.debug("Released shared serial port {} (refCount={})", portName, refCount);

            if (refCount <= 0) {
                existing.getPort().closePort();
                LOGGER.info("Closed shared serial port {}", portName);
                return null; // Remove from map
            }

            return existing;
        });
    }

    /**
     * Configuration for opening a serial port.
     */
    public static class SerialPortConfig {
        public final int baudRate;
        public final int dataBits;
        public final int stopBits;
        public final int parity;
        public final int flowControl;
        public final int readTimeout;
        public final int writeTimeout;
        public final boolean dtr;
        public final boolean rts;
        public final int interframeDelay;

        public SerialPortConfig(int baudRate, int dataBits, int stopBits, int parity,
                                int flowControl, int readTimeout, int writeTimeout,
                                boolean dtr, boolean rts, int interframeDelay) {
            this.baudRate = baudRate;
            this.dataBits = dataBits;
            this.stopBits = stopBits;
            this.parity = parity;
            this.flowControl = flowControl;
            this.readTimeout = readTimeout;
            this.writeTimeout = writeTimeout;
            this.dtr = dtr;
            this.rts = rts;
            this.interframeDelay = interframeDelay;
        }
    }

    /**
     * Wrapper around a SerialPort with reference counting and shared access control.
     */
    public static class SharedPort {
        private final SerialPort port;
        private final String portName;
        private final int interframeDelay;
        private final AtomicInteger refCount = new AtomicInteger(1);
        private final Lock writeLock = new ReentrantLock();
        private volatile long lastWriteTime = 0;

        private SharedPort(SerialPort port, String portName, int interframeDelay) {
            this.port = port;
            this.portName = portName;
            this.interframeDelay = interframeDelay;
        }

        public SerialPort getPort() {
            return port;
        }

        public String getPortName() {
            return portName;
        }

        public int getRefCount() {
            return refCount.get();
        }

        void incrementRefCount() {
            refCount.incrementAndGet();
        }

        int decrementRefCount() {
            return refCount.decrementAndGet();
        }

        /**
         * Acquires the write lock for this port.
         * Ensures only one instance writes at a time and enforces interframe delay.
         */
        public void lockWrite() {
            writeLock.lock();

            // Enforce interframe delay
            if (interframeDelay > 0 && lastWriteTime > 0) {
                long elapsed = System.currentTimeMillis() - lastWriteTime;
                if (elapsed < interframeDelay) {
                    try {
                        Thread.sleep(interframeDelay - elapsed);
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    }
                }
            }
        }

        /**
         * Releases the write lock and updates last write time.
         */
        public void unlockWrite() {
            lastWriteTime = System.currentTimeMillis();
            writeLock.unlock();
        }
    }
}
