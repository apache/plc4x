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
import com.fazecast.jSerialComm.SerialPortDataListener;
import com.fazecast.jSerialComm.SerialPortEvent;
import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.transport.serial.config.SerialTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.OutputStream;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;

/**
 * Serial transport implementation using jSerialComm with async support.
 * Supports both dedicated and shared port modes for optimal resource usage.
 * Implements AsyncTransportInstance for event-driven I/O without polling.
 */
public class SerialTransportInstance extends BaseTransportInstance<SerialTransportConfiguration>
        implements AsyncTransportInstance<SerialTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(SerialTransportInstance.class);
    private static final int DEFAULT_BUFFER_SIZE = 8192;
    private static final byte[] EMPTY_BYTES = new byte[0];

    private final SharedSerialPortManager sharedSerialPortManager;
    private final SerialPort port;
    private final SharedSerialPortManager.SharedPort sharedPort; // null if not shared
    private final OutputStream outputStream;
    private final RingBuffer ringBuffer;
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock(); // Only used for non-shared ports
    private volatile boolean open = true;

    // Async support
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;
    private final SerialPortDataListener serialPortDataListener;
    private volatile Thread readerThread;

    public SerialTransportInstance(SharedSerialPortManager sharedSerialPortManager, String port, SerialTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        this.sharedSerialPortManager = sharedSerialPortManager;
        this.ringBuffer = new RingBuffer(DEFAULT_BUFFER_SIZE);

        try {
            SerialPort tempPort;
            SharedSerialPortManager.SharedPort tempSharedPort = null;

            if (configuration.reusePort) {
                // Use shared port manager
                SharedSerialPortManager.SerialPortConfig portConfig = new SharedSerialPortManager.SerialPortConfig(
                    configuration.baudRate,
                    configuration.dataBits,
                    configuration.stopBits,
                    parseParity(configuration.parity),
                    parseFlowControl(configuration.flowControl),
                    configuration.readTimeout,
                    configuration.writeTimeout,
                    configuration.dtr,
                    configuration.rts,
                    configuration.interframeDelay
                );

                tempSharedPort = sharedSerialPortManager.acquirePort(port, portConfig);
                tempPort = tempSharedPort.getPort();

                LOGGER.debug("Using shared serial port {}",port);

            } else {
                // Create a dedicated port
                tempPort = SerialPort.getCommPort(port);

                // Configure port
                tempPort.setBaudRate(configuration.baudRate);
                tempPort.setNumDataBits(configuration.dataBits);
                tempPort.setNumStopBits(configuration.stopBits);
                tempPort.setParity(parseParity(configuration.parity));
                tempPort.setFlowControl(parseFlowControl(configuration.flowControl));

                // Set timeouts
                tempPort.setComPortTimeouts(
                    SerialPort.TIMEOUT_READ_SEMI_BLOCKING | SerialPort.TIMEOUT_WRITE_BLOCKING,
                    configuration.readTimeout,
                    configuration.writeTimeout
                );

                // Open port
                if (!tempPort.openPort()) {
                    throw new TransportException("Failed to open serial port: " + port);
                }

                // Set DTR/RTS
                if (configuration.dtr) {
                    tempPort.setDTR();
                } else {
                    tempPort.clearDTR();
                }

                if (configuration.rts) {
                    tempPort.setRTS();
                } else {
                    tempPort.clearRTS();
                }

                LOGGER.debug("Created dedicated serial port bound to {}", port);
            }

            LOGGER.info("Serial transport opened on {} at {} baud, {}{}{}, flow control: {}",
                port,
                configuration.baudRate,
                configuration.dataBits,
                configuration.parity.charAt(0),
                configuration.stopBits,
                configuration.flowControl);

            this.port = tempPort;
            this.sharedPort = tempSharedPort;
            this.outputStream = tempPort.getOutputStream();

            // Create the serial port data listener for async I/O
            // Note: Some platforms/devices don't properly support SerialPortDataListener events,
            // so we use a background reader thread as a fallback
            this.serialPortDataListener = new SerialPortDataListener() {
                @Override
                public int getListeningEvents() {
                    return SerialPort.LISTENING_EVENT_DATA_AVAILABLE;
                }

                @Override
                public void serialEvent(SerialPortEvent event) {
                    if (event.getEventType() != SerialPort.LISTENING_EVENT_DATA_AVAILABLE) {
                        return;
                    }
                    readFromPort();
                }
            };

            // Try to register the event listener first
            boolean eventListenerRegistered = tempPort.addDataListener(serialPortDataListener);
            LOGGER.debug("Serial port event listener registration result: {}", eventListenerRegistered);

            // Always start a background reader thread as fallback
            // This ensures data is read even on platforms where events don't fire properly
            this.readerThread = new Thread(this::readerLoop, "Serial-Reader-" + port);
            this.readerThread.setDaemon(true);
            this.readerThread.start();
            LOGGER.debug("Serial port reader thread started");

            getAuditLog().write(AuditLogEventType.CONNECT, String.format(
                "Serial port opened on %s at %d baud, %d%s%d, flow control: %s",
                port, configuration.baudRate, configuration.dataBits,
                configuration.parity.charAt(0), configuration.stopBits, configuration.flowControl));
        } catch (Exception e) {
            String errorMsg = String.format("Failed to create serial transport for %s - %s",
                port, e.getMessage());
            LOGGER.error(errorMsg, e);
            getAuditLog().write(AuditLogEventType.ERROR, "Error in constructor: " + errorMsg);
            throw new TransportException(errorMsg, e);
        }
    }

    /**
     * Converts parity string to jSerialComm constant.
     */
    private int parseParity(String parity) {
        return switch (parity.toUpperCase()) {
            case "NONE" -> SerialPort.NO_PARITY;
            case "ODD" -> SerialPort.ODD_PARITY;
            case "EVEN" -> SerialPort.EVEN_PARITY;
            case "MARK" -> SerialPort.MARK_PARITY;
            case "SPACE" -> SerialPort.SPACE_PARITY;
            default -> {
                LOGGER.warn("Unknown parity '{}', using NONE", parity);
                yield SerialPort.NO_PARITY;
            }
        };
    }

    /**
     * Converts flow control string to jSerialComm constant.
     */
    private int parseFlowControl(String flowControl) {
        return switch (flowControl.toUpperCase()) {
            case "NONE" -> SerialPort.FLOW_CONTROL_DISABLED;
            case "RTS_CTS", "RTSCTS" -> SerialPort.FLOW_CONTROL_RTS_ENABLED | SerialPort.FLOW_CONTROL_CTS_ENABLED;
            case "XON_XOFF", "XONXOFF" -> SerialPort.FLOW_CONTROL_XONXOFF_IN_ENABLED | SerialPort.FLOW_CONTROL_XONXOFF_OUT_ENABLED;
            case "RTS_CTS_XON_XOFF" -> SerialPort.FLOW_CONTROL_RTS_ENABLED | SerialPort.FLOW_CONTROL_CTS_ENABLED |
                    SerialPort.FLOW_CONTROL_XONXOFF_IN_ENABLED | SerialPort.FLOW_CONTROL_XONXOFF_OUT_ENABLED;
            default -> {
                LOGGER.warn("Unknown flow control '{}', using NONE", flowControl);
                yield SerialPort.FLOW_CONTROL_DISABLED;
            }
        };
    }

    @Override
    public boolean isOpen() {
        return open && port.isOpen();
    }

    @Override
    public int getNumBytesAvailable() throws TransportException {
        readLock.lock();
        try {
            if (!isOpen()) {
                return 0;
            }

            return ringBuffer.availableForReading();
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public byte[] peekReadableBytes(int numBytes) throws TransportException {
        if (numBytes <= 0) {
            return EMPTY_BYTES;
        }

        readLock.lock();
        try {
            ensureOpen();

            // Fill the ring buffer if necessary
            getNumBytesAvailable();

            if (ringBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, ringBuffer.availableForReading())
                );
            }

            // Peek without consuming
            return ringBuffer.peek(numBytes);
        } catch (TransportException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in peekReadableBytes: " + e.getMessage());
            throw e;
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public byte[] read(int numBytes) throws TransportException {
        if (numBytes <= 0) {
            return EMPTY_BYTES;
        }

        readLock.lock();
        try {
            ensureOpen();

            // Fill the ring buffer if necessary
            getNumBytesAvailable();

            if (ringBuffer.availableForReading() < numBytes) {
                throw new TransportException(
                    String.format("Requested %d bytes but only %d available", numBytes, ringBuffer.availableForReading())
                );
            }

            // Read and consume bytes
            byte[] bytes = ringBuffer.read(numBytes);

            // Log the bytes to the audit log
            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.INCOMING_BYTES, StaticHelper.ENCODE_HEX(bytes));
            }

            return bytes;
        } catch (TransportException e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in read: " + e.getMessage());
            throw e;
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public void write(byte[] bytes) throws TransportException {
        if (bytes == null || bytes.length == 0) {
            return;
        }

        // Use shared port's write lock if available, otherwise use local lock
        if (sharedPort != null) {
            sharedPort.lockWrite();
        } else {
            writeLock.lock();
        }

        try {
            ensureOpen();

            outputStream.write(bytes);
            outputStream.flush();

            LOGGER.trace("Wrote {} bytes to serial port", bytes.length);

            // Log the bytes to the audit log
            if (getAuditLog().isEnabled()) {
                getAuditLog().write(AuditLogEventType.OUTGOING_BYTES, "Write: " + StaticHelper.ENCODE_HEX(bytes));
            }
        } catch (Exception e) {
            getAuditLog().write(AuditLogEventType.ERROR, "Error in write: " + e.getMessage());
            throw new TransportException("Failed to write to serial port", e);
        } finally {
            if (sharedPort != null) {
                sharedPort.unlockWrite();
            } else {
                writeLock.unlock();
            }
        }
    }

    @Override
    public void close() throws TransportException {
        if (!open) {
            return;
        }

        // Set open to false first to stop the reader thread
        open = false;

        // Stop the reader thread
        if (readerThread != null) {
            readerThread.interrupt();
            try {
                readerThread.join(1000);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }

        // Remove the data listener
        port.removeDataListener();
        LOGGER.debug("Serial port data listener removed");

        if (sharedPort != null) {
            // Don't lock for shared ports - manager handles it
            readLock.lock();
            try {
                open = false;
                sharedSerialPortManager.releasePort(sharedPort);
                LOGGER.debug("Released shared serial port");
                getAuditLog().write(AuditLogEventType.CLOSE, "Released shared serial port");
            } finally {
                readLock.unlock();
            }
        } else {
            // Lock for dedicated ports
            writeLock.lock();
            try {
                readLock.lock();
                try {
                    open = false;
                    port.closePort();
                    LOGGER.debug("Closed dedicated serial port");
                    getAuditLog().write(AuditLogEventType.CLOSE, "Closed dedicated serial port");
                } finally {
                    readLock.unlock();
                }
            } finally {
                writeLock.unlock();
            }
        }
    }

    // ========== AsyncTransportInstance Implementation ==========

    @Override
    public void registerDataListener(Runnable listener) {
        this.dataListener = listener;
        LOGGER.debug("Data listener registered");
    }

    @Override
    public void removeDataListener() {
        this.dataListener = null;
        LOGGER.debug("Data listener removed");
    }

    @Override
    public void registerDisconnectListener(Consumer<Throwable> listener) {
        this.disconnectListener = listener;
        LOGGER.debug("Disconnect listener registered");
    }

    @Override
    public void removeDisconnectListener() {
        this.disconnectListener = null;
        LOGGER.debug("Disconnect listener removed");
    }

    /**
     * Notifies the disconnect listener if one is registered.
     *
     * @param cause the exception that caused the disconnect, or null for graceful close
     */
    private void notifyDisconnect(Throwable cause) {
        Consumer<Throwable> listener = disconnectListener;
        if (listener != null) {
            try {
                listener.accept(cause);
            } catch (Exception e) {
                LOGGER.error("Error in disconnect listener", e);
            }
        }
    }

    /**
     * Ensures the transport is still open, throws exception otherwise.
     */
    private void ensureOpen() throws TransportException {
        if (!isOpen()) {
            throw new TransportException("Transport is closed");
        }
    }

    /**
     * Reads available data from the serial port into the ring buffer.
     * This method is called both from the event listener and the reader thread.
     */
    private void readFromPort() {
        readLock.lock();
        try {
            if (!open) {
                return;
            }

            int available = port.bytesAvailable();
            if (available <= 0) {
                return;
            }

            // Check available space in ring buffer
            int availableSpace = ringBuffer.remainingForWriting();
            if (availableSpace == 0) {
                LOGGER.warn("Ring buffer is full, discarding {} bytes from serial port", available);
                return;
            }

            // Read only as much as we can store
            int bytesToRead = Math.min(available, availableSpace);
            byte[] buffer = new byte[bytesToRead];
            int bytesRead = port.readBytes(buffer, bytesToRead);

            if (bytesRead > 0) {
                ringBuffer.write(buffer, 0, bytesRead);
                LOGGER.debug("Read {} bytes from serial port into ring buffer, buffer now has {} bytes",
                    bytesRead, ringBuffer.availableForReading());

                // Notify the data listener if registered
                Runnable listener = dataListener;
                if (listener != null) {
                    listener.run();
                }
            }
        } finally {
            readLock.unlock();
        }
    }

    /**
     * Background reader loop that polls for data from the serial port.
     * This is used as a fallback for platforms where SerialPortDataListener events
     * don't fire properly (e.g., PTY devices on macOS).
     */
    private void readerLoop() {
        LOGGER.debug("Serial reader loop started");

        while (open && !Thread.currentThread().isInterrupted()) {
            try {
                readFromPort();

                // Small sleep to prevent busy-waiting
                // The timeout on the serial port helps, but we add a small delay
                // to reduce CPU usage when no data is available
                Thread.sleep(10);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break;
            } catch (Exception e) {
                if (open) {
                    LOGGER.error("Error in serial reader loop", e);
                }
            }
        }

        LOGGER.debug("Serial reader loop stopped");
    }
}
