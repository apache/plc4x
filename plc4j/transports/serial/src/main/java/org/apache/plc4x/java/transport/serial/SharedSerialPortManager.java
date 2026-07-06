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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Function;

/**
 * Manages shared serial ports across multiple transport instances.
 * Allows multiple logical connections to share a single serial port,
 * which is essential for protocols like Modbus RTU that multiplex over one serial line.
 */
public class SharedSerialPortManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(SharedSerialPortManager.class);

    // Map of port name -> shared port wrapper
    private final Map<String, SharedPort> sharedPorts = new HashMap<>();

    private final Function<String, SerialPort> portFactory;

    public SharedSerialPortManager() {
        this(SerialPort::getCommPort);
    }

    // Visible for tests: lets unit tests supply mocked physical ports.
    SharedSerialPortManager(Function<String, SerialPort> portFactory) {
        this.portFactory = portFactory;
    }

    /**
     * Acquires a shared serial port. Joining an existing entry requires an
     * IDENTICAL configuration; a mismatch is a hard error (no silent
     * first-config-wins).
     */
    public synchronized SharedPort acquirePort(String portName, SerialPortConfig config) throws TransportException {
        SharedPort existing = sharedPorts.get(portName);
        if (existing != null && existing.isClosed()) {
            // An acquire can race the fail() window where the entry is
            // dying (closed CAS already flipped) but not yet evicted from
            // the map (evict() and this method both synchronize on the
            // manager, but fail() may not have reached evict() yet, or a
            // stale entry may still linger under this name). Treat it as
            // absent and fall through to opening a fresh port.
            sharedPorts.remove(portName);
            existing = null;
        }
        if (existing != null) {
            if (!existing.getConfig().matches(config)) {
                throw new TransportException(String.format(
                    "Serial port %s is already shared with a different configuration (existing: %s, requested: %s)",
                    portName, existing.getConfig(), config));
            }
            existing.incrementRefCount();
            LOGGER.debug("Reusing shared serial port {} (refCount={})", portName, existing.getRefCount());
            return existing;
        }

        SerialPort port = portFactory.apply(portName);

        // Configure port
        port.setBaudRate(config.baudRate);
        port.setNumDataBits(config.dataBits);
        port.setNumStopBits(config.stopBits);
        port.setParity(config.parity);
        port.setFlowControl(config.flowControl);
        port.setComPortTimeouts(
            SerialPort.TIMEOUT_READ_SEMI_BLOCKING | SerialPort.TIMEOUT_WRITE_BLOCKING,
            config.readTimeout,
            config.writeTimeout
        );

        if (!port.openPort()) {
            throw new TransportException("Failed to open serial port: " + portName);
        }

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
        SharedPort sharedPort = new SharedPort(port, portName, config, this);
        sharedPorts.put(portName, sharedPort);
        sharedPort.startReader();
        return sharedPort;
    }

    /**
     * Releases a shared serial port reference.
     * If this is the last reference, closes the port.
     */
    public void releasePort(SharedPort port) {
        String portName = port.getPortName();
        boolean shouldShutdown;
        synchronized (this) {
            SharedPort existing = sharedPorts.get(portName);
            if (existing != port) {
                // The entry was already evicted (fatal failure) and possibly
                // replaced by a fresh acquire under the same name. A stale
                // holder releasing here must not touch the new entry.
                LOGGER.warn("releasing an already-evicted shared serial port {}", portName);
                return;
            }
            int refCount = existing.decrementRefCount();
            LOGGER.debug("Released shared serial port {} (refCount={})", portName, refCount);
            shouldShutdown = refCount <= 0;
            if (shouldShutdown) {
                sharedPorts.remove(portName);
            }
        }
        // shutdown() joins the reader thread; it must run OUTSIDE the
        // manager lock. Otherwise, a reader concurrently inside fail()
        // (which calls the synchronized manager.evict(this)) would stall the
        // join while holding up this lock, letting closePort() race with the
        // still-in-progress fail(). The closed-guard in SharedPort makes the
        // ordering safe regardless, but staying off the lock avoids the stall.
        if (shouldShutdown) {
            port.shutdown();
            LOGGER.info("Closed shared serial port {}", portName);
        }
    }

    /**
     * Evicts a shared port entry after a fatal failure so that a subsequent
     * acquire re-opens the physical device instead of reusing a dead entry.
     */
    synchronized void evict(SharedPort port) {
        if (sharedPorts.get(port.getPortName()) == port) {
            sharedPorts.remove(port.getPortName());
        }
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

        public boolean matches(SerialPortConfig other) {
            return other != null
                && baudRate == other.baudRate
                && dataBits == other.dataBits
                && stopBits == other.stopBits
                && parity == other.parity
                && flowControl == other.flowControl
                && readTimeout == other.readTimeout
                && writeTimeout == other.writeTimeout
                && dtr == other.dtr
                && rts == other.rts
                && interframeDelay == other.interframeDelay;
        }

        @Override
        public String toString() {
            return String.format(
                "baud=%d dataBits=%d stopBits=%d parity=%d flowControl=%d readTimeout=%d writeTimeout=%d dtr=%b rts=%b interframeDelay=%d",
                baudRate, dataBits, stopBits, parity, flowControl, readTimeout, writeTimeout, dtr, rts, interframeDelay);
        }
    }

    /**
     * Wrapper around a SerialPort with reference counting and shared access control.
     */
    public static class SharedPort {
        private static final int READER_CHUNK = 4096;

        private final SerialPort port;
        private final String portName;
        private final SerialPortConfig config;
        private final AtomicInteger refCount = new AtomicInteger(1);
        private final Lock writeLock = new ReentrantLock();

        private final java.util.List<SharedPortSubscriber> subscribers = new java.util.concurrent.CopyOnWriteArrayList<>();
        private final WritePacer pacer;
        private volatile boolean open = true;
        private volatile Thread readerThread;
        private final SharedSerialPortManager manager;
        private final java.util.concurrent.atomic.AtomicBoolean closed = new java.util.concurrent.atomic.AtomicBoolean(false);

        private SharedPort(SerialPort port, String portName, SerialPortConfig config, SharedSerialPortManager manager) {
            this.port = port;
            this.portName = portName;
            this.config = config;
            this.manager = manager;
            this.pacer = new WritePacer(config.interframeDelay);
        }

        public SerialPort getPort() {
            return port;
        }

        public String getPortName() {
            return portName;
        }

        public SerialPortConfig getConfig() {
            return config;
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
         * True once this entry has started (or finished) tearing down, via
         * either {@link #fail(Throwable)} or {@link #shutdown()}. Lets
         * {@link SharedSerialPortManager#acquirePort} detect a dying entry
         * that hasn't been evicted from the map yet.
         */
        boolean isClosed() {
            return closed.get();
        }

        void addSubscriber(SharedPortSubscriber subscriber) {
            subscribers.add(subscriber);
        }

        void removeSubscriber(SharedPortSubscriber subscriber) {
            subscribers.remove(subscriber);
        }

        /**
         * Starts the single reader for this physical port: an event listener
         * plus a polling fallback thread (some platforms/devices never fire
         * events). Both funnel through readFromPort(), which is what
         * broadcasts to subscribers — the per-port reader replaces the old
         * per-instance competing readers.
         */
        void startReader() {
            port.addDataListener(new com.fazecast.jSerialComm.SerialPortDataListener() {
                @Override
                public int getListeningEvents() {
                    return SerialPort.LISTENING_EVENT_DATA_AVAILABLE;
                }

                @Override
                public void serialEvent(com.fazecast.jSerialComm.SerialPortEvent event) {
                    if (event.getEventType() == SerialPort.LISTENING_EVENT_DATA_AVAILABLE) {
                        readFromPort();
                    }
                }
            });
            readerThread = new Thread(this::readerLoop, "Shared-Serial-Reader-" + portName);
            readerThread.setDaemon(true);
            readerThread.start();
        }

        private void readerLoop() {
            while (open && !Thread.currentThread().isInterrupted()) {
                try {
                    readFromPort();
                    Thread.sleep(10);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    return;
                } catch (Exception e) {
                    // An unexpected RuntimeException out of readFromPort()
                    // must escalate to a fatal failure instead of silently
                    // killing the polling thread and leaving subscribers
                    // waiting on a port that looks alive but has no reader.
                    fail(e);
                    return;
                }
            }
        }

        private synchronized void readFromPort() {
            if (!open) {
                return;
            }
            int available = port.bytesAvailable();
            if (available < 0) {
                fail(new java.io.IOException("Serial port " + portName + " reported an error (bytesAvailable=" + available + ")"));
                return;
            }
            if (available == 0) {
                return;
            }
            byte[] buffer = new byte[Math.min(available, READER_CHUNK)];
            int bytesRead = port.readBytes(buffer, buffer.length);
            if (bytesRead < 0) {
                fail(new java.io.IOException("Serial port " + portName + " read failed (" + bytesRead + ")"));
                return;
            }
            if (bytesRead == 0) {
                return;
            }
            pacer.noteActivity();
            for (SharedPortSubscriber subscriber : subscribers) {
                try {
                    subscriber.onData(buffer, 0, bytesRead);
                } catch (Exception e) {
                    LOGGER.warn("Shared serial subscriber threw in onData on {}", portName, e);
                }
            }
        }

        /**
         * Fatal path: evict the entry and release the physical port FIRST,
         * THEN notify subscribers.
         *
         * The fd must be released before listeners run: a disconnect
         * listener can synchronously attempt to reconnect, and jSerialComm
         * holds an exclusive lock on the device, so a reopen would fail if
         * this port hadn't already released it. Ordering: evict (so the
         * reconnect's acquirePort() re-opens the device instead of reusing
         * the dying entry) -> removeDataListener/closePort (release the fd)
         * -> onFailure to all subscribers.
         *
         * Guarded by {@link #closed} so that a concurrent {@link #shutdown()}
         * (triggered by a release racing this failure) can't result in both
         * paths running: whichever loses the CAS no-ops entirely, avoiding a
         * stale onFailure delivery or a double close.
         */
        private void fail(Throwable cause) {
            if (!closed.compareAndSet(false, true)) {
                return;
            }
            open = false;
            LOGGER.warn("Shared serial port {} failed", portName, cause);
            manager.evict(this);
            port.removeDataListener();
            port.closePort();
            for (SharedPortSubscriber subscriber : subscribers) {
                try {
                    subscriber.onFailure(cause);
                } catch (Exception e) {
                    LOGGER.warn("Shared serial subscriber threw in onFailure on {}", portName, e);
                }
            }
        }

        /**
         * Stops the reader (listener + thread) and closes the port.
         *
         * Ordering rationale: the port is closed BEFORE interrupting/joining
         * the reader thread. A reader thread can be parked inside a native
         * jSerialComm blocking read; a plain Java {@code Thread.interrupt()}
         * does not unblock that native call, but closing the underlying port
         * does. Once closed, the reader loop's {@code open} check (both at
         * the top of {@link #readFromPort()} and in the {@code readerLoop}
         * condition) causes it to exit quietly rather than treating the
         * close as a fatal failure. The {@link #closed} CAS guard also makes
         * this safe against a concurrent {@link #fail()}: whichever call
         * wins the race performs the teardown, the other is a no-op, so
         * shutdown() never has to wait on the manager lock that a
         * concurrent fail()->manager.evict() might be holding.
         */
        void shutdown() {
            if (!closed.compareAndSet(false, true)) {
                return;
            }
            open = false;
            port.removeDataListener();
            port.closePort();
            Thread thread = readerThread;
            if (thread != null) {
                thread.interrupt();
                try {
                    thread.join(1000);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }
        }

        /**
         * Acquires the write lock for this port.
         * Ensures only one instance writes at a time and enforces interframe delay.
         */
        public void lockWrite() {
            writeLock.lock();
            try {
                pacer.awaitTurn();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }

        /**
         * Releases the write lock and updates last write time.
         */
        public void unlockWrite() {
            pacer.noteActivity();
            writeLock.unlock();
        }
    }
}
