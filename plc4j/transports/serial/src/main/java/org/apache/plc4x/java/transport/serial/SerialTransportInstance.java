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
import java.util.Locale;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.RejectedExecutionException;
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
    // Per-connection receive ring; matches the plc4go shared-port
    // subscriber rings (64 KiB) so burst tolerance is consistent.
    private static final int DEFAULT_BUFFER_SIZE = 65536;
    private static final byte[] EMPTY_BYTES = new byte[0];

    private final SharedSerialPortManager sharedSerialPortManager;
    private final SerialPort port;
    private final SharedSerialPortManager.SharedPort sharedPort; // null if not shared
    private final OutputStream outputStream;
    private final RingBuffer ringBuffer;
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock(); // Only used for non-shared ports
    private volatile boolean open = true;
    // Guards close() so two concurrent callers can't both pass the
    // check-then-act on `open` and both proceed to release resources (e.g.
    // double-decrementing the shared port's refcount).
    private final java.util.concurrent.atomic.AtomicBoolean closeGuard = new java.util.concurrent.atomic.AtomicBoolean(false);

    // Async support
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;
    private SerialPortDataListener serialPortDataListener;
    private volatile Thread readerThread;
    private final WritePacer writePacer; // dedicated-path pacing; no-op pacer in shared mode
    private SharedPortSubscriber sharedSubscriber;
    private long droppedBytes;
    private long lastWarnedDroppedBytes;
    // Shared mode only: dispatches dataListener callbacks off the shared
    // reader thread so one blocking callback stalls only this connection
    // (dedicated mode already has this isolation via its own reader
    // thread). Single-threaded => notifications stay ordered.
    private volatile ExecutorService sharedDispatchExecutor;
    private final java.util.concurrent.atomic.AtomicBoolean dispatchPending = new java.util.concurrent.atomic.AtomicBoolean(false);

    public SerialTransportInstance(SharedSerialPortManager sharedSerialPortManager, String port, SerialTransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        this.sharedSerialPortManager = sharedSerialPortManager;
        this.ringBuffer = new RingBuffer(DEFAULT_BUFFER_SIZE);

        // Hoisted so the catch block below can see whatever got created
        // before the failure, for best-effort cleanup.
        SerialPort tempPort = null;
        SharedSerialPortManager.SharedPort tempSharedPort = null;

        try {
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
            this.writePacer = new WritePacer(tempSharedPort != null ? 0 : configuration.interframeDelay);

            if (tempSharedPort != null) {
                // Shared mode: the SharedPort owns the single reader; this
                // instance only subscribes to the broadcast.
                this.serialPortDataListener = null;
                this.readerThread = null;
                this.sharedSubscriber = new SharedPortSubscriber() {
                    @Override
                    public void onData(byte[] data, int offset, int length) {
                        deliverSharedData(data, offset, length);
                    }

                    @Override
                    public void onFailure(Throwable cause) {
                        Consumer<Throwable> listener = disconnectListener;
                        if (listener != null) {
                            listener.accept(cause);
                        }
                    }
                };
                tempSharedPort.addSubscriber(this.sharedSubscriber);
                this.sharedDispatchExecutor = Executors.newSingleThreadExecutor(runnable -> {
                    Thread thread = new Thread(runnable, "Serial-Shared-Dispatch-" + port);
                    thread.setDaemon(true);
                    return thread;
                });
            } else {
                // Dedicated mode: keep the per-instance reader exactly as before.
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
                boolean eventListenerRegistered = tempPort.addDataListener(serialPortDataListener);
                LOGGER.debug("Serial port event listener registration result: {}", eventListenerRegistered);
                this.readerThread = new Thread(this::readerLoop, "Serial-Reader-" + port);
                this.readerThread.setDaemon(true);
                this.readerThread.start();
                LOGGER.debug("Serial port reader thread started");
            }

            getAuditLog().write(AuditLogEventType.CONNECT, String.format(
                "Serial port opened on %s at %d baud, %d%s%d, flow control: %s",
                port, configuration.baudRate, configuration.dataBits,
                configuration.parity.charAt(0), configuration.stopBits, configuration.flowControl));
        } catch (Exception e) {
            String errorMsg = String.format("Failed to create serial transport for %s - %s",
                port, e.getMessage());
            LOGGER.error(errorMsg, e);

            // Best-effort cleanup FIRST: never leak a subscriber/refcount
            // (shared) or an open port (dedicated) out of a failed
            // constructor — and never let a throwing audit log skip it.
            try {
                if (tempSharedPort != null) {
                    if (this.sharedSubscriber != null) {
                        tempSharedPort.removeSubscriber(this.sharedSubscriber);
                    }
                    sharedSerialPortManager.releasePort(tempSharedPort);
                } else if (tempPort != null) {
                    tempPort.closePort();
                }
                if (sharedDispatchExecutor != null) {
                    sharedDispatchExecutor.shutdown();
                }
            } catch (Exception cleanupError) {
                LOGGER.warn("Cleanup after failed serial transport construction also failed", cleanupError);
            }

            try {
                getAuditLog().write(AuditLogEventType.ERROR, "Error in constructor: " + errorMsg);
            } catch (Exception auditError) {
                LOGGER.warn("Audit log write after failed serial transport construction also failed", auditError);
            }

            throw new TransportException(errorMsg, e);
        }
    }

    /**
     * Converts a parity option value to the matching jSerialComm constant.
     * Values are case-insensitive and accept "-" or "_" as separator.
     */
    static int parseParity(String parity) throws TransportException {
        return switch (normalizeOptionValue(parity)) {
            case "none" -> SerialPort.NO_PARITY;
            case "odd" -> SerialPort.ODD_PARITY;
            case "even" -> SerialPort.EVEN_PARITY;
            case "mark" -> SerialPort.MARK_PARITY;
            case "space" -> SerialPort.SPACE_PARITY;
            default -> throw new TransportException(
                "Invalid value '" + parity + "' for option 'parity' (must be one of: none, odd, even, mark, space)");
        };
    }

    /**
     * Converts a flow-control option value to the matching jSerialComm
     * constant. Values are case-insensitive and accept "-" or "_" as
     * separator. Combining hardware and software flow control is not
     * supported (matching the plc4go serial transport).
     */
    static int parseFlowControl(String flowControl) throws TransportException {
        return switch (normalizeOptionValue(flowControl)) {
            case "none" -> SerialPort.FLOW_CONTROL_DISABLED;
            case "rts-cts", "rtscts" -> SerialPort.FLOW_CONTROL_RTS_ENABLED | SerialPort.FLOW_CONTROL_CTS_ENABLED;
            case "xon-xoff", "xonxoff" -> SerialPort.FLOW_CONTROL_XONXOFF_IN_ENABLED | SerialPort.FLOW_CONTROL_XONXOFF_OUT_ENABLED;
            default -> throw new TransportException(
                "Invalid value '" + flowControl + "' for option 'flow-control' (must be one of: none, rts-cts, xon-xoff)");
        };
    }

    private static String normalizeOptionValue(String value) {
        return value.toLowerCase(Locale.ROOT).replace('_', '-');
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
            try {
                writePacer.awaitTurn();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
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
                writePacer.noteActivity();
                writeLock.unlock();
            }
        }
    }

    @Override
    public void close() throws TransportException {
        if (!closeGuard.compareAndSet(false, true)) {
            // Already closed (or a concurrent close is in flight); avoid
            // double teardown, e.g. a double decrement of the shared port's
            // refcount.
            return;
        }

        // Set open to false first to stop the reader thread
        open = false;

        if (sharedPort != null) {
            // Don't lock for shared ports - manager handles it. No reader
            // thread/listener to stop here: the SharedPort owns those and
            // must keep serving the other subscribers. Only remove our own
            // subscription.
            //
            // Only flip `open` under readLock (so an in-flight
            // deliverSharedData() call observes the close and bails out
            // cleanly). removeSubscriber() and releasePort() must run
            // OUTSIDE the lock: on the last release, releasePort() ->
            // SharedPort.shutdown() joins the shared reader thread, and that
            // reader may be blocked inside deliverSharedData() waiting for
            // this very readLock. Calling releasePort() while holding the
            // lock would make every close racing live data pay the ~1s join
            // timeout (mirrors the rationale already documented on
            // SharedSerialPortManager.releasePort()).
            readLock.lock();
            try {
                open = false;
            } finally {
                readLock.unlock();
            }
            if (sharedSubscriber != null) {
                sharedPort.removeSubscriber(sharedSubscriber);
            }
            sharedSerialPortManager.releasePort(sharedPort);
            ExecutorService executor = sharedDispatchExecutor;
            if (executor != null) {
                executor.shutdown(); // in-flight dispatch may finish; no await needed
            }
            LOGGER.debug("Released shared serial port");
            getAuditLog().write(AuditLogEventType.CLOSE, "Released shared serial port");
        } else {
            // Stop the reader thread (dedicated only)
            if (readerThread != null) {
                readerThread.interrupt();
                try {
                    readerThread.join(1000);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }

            // Remove the data listener (dedicated only)
            port.removeDataListener();
            LOGGER.debug("Serial port data listener removed");

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
                writePacer.noteActivity();

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
     * Shared-mode data path: writes a broadcast chunk into this instance's
     * ring buffer with drop-OLDEST overflow semantics (mirroring plc4go) and
     * wakes the codec-facing data listener. Called from the shared reader
     * thread.
     */
    private void deliverSharedData(byte[] data, int offset, int length) {
        readLock.lock();
        try {
            if (!open) {
                return;
            }
            int dropped = writeDroppingOldest(ringBuffer, data, offset, length);
            if (dropped > 0) {
                droppedBytes += dropped;
                if (lastWarnedDroppedBytes == 0 || droppedBytes >= 2 * lastWarnedDroppedBytes) {
                    lastWarnedDroppedBytes = droppedBytes;
                    LOGGER.warn("Ring buffer overflow on shared serial port, dropped {} oldest bytes (total dropped: {})",
                        dropped, droppedBytes);
                }
            }
            if (dataListener != null && dispatchPending.compareAndSet(false, true)) {
                try {
                    sharedDispatchExecutor.execute(this::runSharedDispatch);
                } catch (RejectedExecutionException e) {
                    dispatchPending.set(false); // closing; nothing to dispatch to
                }
            }
        } finally {
            readLock.unlock();
        }
    }

    /**
     * Runs the dataListener callback off the shared reader thread (shared
     * mode only). Submitted to {@link #sharedDispatchExecutor}, a
     * single-threaded executor, so callback invocations stay ordered and
     * never overlap.
     */
    private void runSharedDispatch() {
        // Clear BEFORE running: data arriving during the run triggers
        // exactly one follow-up task; the listener drains everything
        // available per invocation, so coalescing loses nothing.
        dispatchPending.set(false);
        Runnable listener = dataListener;
        if (listener != null && open) {
            try {
                listener.run();
            } catch (Exception e) {
                // Pre-executor, SharedPort's broadcast wrapper caught and
                // logged listener exceptions; keep that observability here.
                LOGGER.warn("Serial data listener threw during shared-mode dispatch", e);
            }
        }
    }

    /**
     * Writes a chunk into the ring buffer, dropping the OLDEST bytes when
     * capacity would be exceeded (mirroring the plc4go shared-port ring
     * semantics — RingBuffer.write alone would drop the newest). Returns
     * how many bytes were dropped.
     */
    static int writeDroppingOldest(RingBuffer ringBuffer, byte[] data, int offset, int length) {
        int capacity = ringBuffer.capacity();
        if (length >= capacity) {
            int dropped = ringBuffer.availableForReading() + (length - capacity);
            ringBuffer.clear();
            ringBuffer.write(data, offset + (length - capacity), capacity);
            return dropped;
        }
        int remaining = ringBuffer.remainingForWriting();
        int dropped = 0;
        if (remaining < length) {
            dropped = length - remaining;
            ringBuffer.skip(dropped);
        }
        ringBuffer.write(data, offset, length);
        return dropped;
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
