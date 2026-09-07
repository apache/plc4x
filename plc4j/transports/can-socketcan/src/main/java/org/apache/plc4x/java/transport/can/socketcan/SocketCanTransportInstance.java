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
package org.apache.plc4x.java.transport.can.socketcan;

import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.can.CanIdFilter;
import org.apache.plc4x.java.transport.can.socketcan.config.SocketCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tel.schich.javacan.CanChannels;
import tel.schich.javacan.CanFrame;
import tel.schich.javacan.RawCanChannel;

import java.io.IOException;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;

/**
 * Transport instance for Linux SocketCAN communication.
 * <p>
 * Opens a CAN socket (either dedicated or shared via {@link SharedCanManager}),
 * starts a background reader thread to receive CAN frames, applies CAN ID filtering,
 * and buffers matching frames in a {@link org.apache.plc4x.java.spi.transports.api.RingBuffer} for reading.
 * <p>
 * Thread safety:
 * <ul>
 *   <li>Listeners stored in volatile fields</li>
 *   <li>Read/write operations protected by ReentrantLocks</li>
 *   <li>Open state tracked by volatile boolean</li>
 * </ul>
 */
public class SocketCanTransportInstance extends BaseTransportInstance<SocketCanTransportConfiguration>
        implements AsyncTransportInstance<SocketCanTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(SocketCanTransportInstance.class);

    /** Size of a CAN 2.0 frame in the SocketCAN wire format (16 bytes). */
    private static final int CAN_FRAME_SIZE = 16;

    private final SharedCanManager sharedCanManager;
    private final CanIdFilter filter;
    private final RingBuffer ringBuffer = new RingBuffer(65536);
    private final Lock readLock = new ReentrantLock();
    private final Lock writeLock = new ReentrantLock();

    private volatile boolean open;
    private volatile Thread readerThread;
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;

    // Either a shared handle (if reuseInterface) or a dedicated channel
    private SharedCanManager.SharedCanHandle sharedHandle;
    private RawCanChannel channel;

    /**
     * Creates a new SocketCAN transport instance.
     *
     * @param sharedCanManager the shared handle manager
     * @param configuration    transport configuration
     * @param auditLog         audit log for connection events
     * @throws TransportException if the platform is not Linux or the CAN interface cannot be opened
     */
    public SocketCanTransportInstance(SharedCanManager sharedCanManager,
                                      SocketCanTransportConfiguration configuration,
                                      AuditLog auditLog) throws TransportException {
        super(configuration, auditLog);
        this.sharedCanManager = sharedCanManager;
        this.filter = configuration.buildFilter();

        LOGGER.debug("SocketCanTransportInstance: Java 17 version");

        // Validate Linux platform
        validatePlatform();

        String interfaceName = configuration.interfaceName;
        if (interfaceName == null || interfaceName.trim().isEmpty()) {
            throw new TransportException("CAN interface name must be specified");
        }

        try {
            if (configuration.reuseInterface) {
                // Shared mode: acquire handle from manager
                this.sharedHandle = sharedCanManager.acquireHandle(interfaceName);
                this.channel = sharedHandle.getChannel();
                LOGGER.debug("Using shared CAN socket on {} (refCount={})",
                        interfaceName, sharedHandle.getRefCount());
            } else {
                // Dedicated mode: open own channel
                this.channel = openChannel(interfaceName);
                LOGGER.debug("Opened dedicated CAN socket on {}", interfaceName);
            }
        } catch (IOException e) {
            throw new TransportException("Failed to open CAN socket on " + interfaceName, e);
        }

        this.open = true;
        startReaderThread(interfaceName);
    }

    /**
     * Validates that the current platform supports SocketCAN.
     * Overridable for testing on non-Linux platforms.
     *
     * @throws TransportException if the platform is not Linux
     */
    protected void validatePlatform() throws TransportException {
        String osName = System.getProperty("os.name", "").toLowerCase();
        if (!osName.contains("linux")) {
            throw new TransportException(
                    "SocketCAN transport is only available on Linux. Current OS: " + System.getProperty("os.name"));
        }
    }

    /**
     * Opens a dedicated CAN channel for the given interface.
     * Overridable for testing without actual CAN hardware.
     *
     * @param interfaceName the CAN interface name
     * @return the opened channel
     * @throws IOException if the channel cannot be opened
     */
    protected RawCanChannel openChannel(String interfaceName) throws IOException {
        return CanChannels.newRawChannel(interfaceName);
    }

    /**
     * Starts the background reader thread that receives CAN frames from the socket.
     *
     * @param interfaceName the CAN interface name for thread naming
     */
    protected void startReaderThread(String interfaceName) {
        readerThread = new Thread(this::runReaderLoop, "SocketCAN-Reader-" + interfaceName);
        readerThread.setDaemon(true);
        readerThread.start();
    }

    /**
     * Main reader loop: continuously reads CAN frames from the socket,
     * applies the CAN ID filter, and buffers matching frames.
     */
    private void runReaderLoop() {
        while (open && !Thread.currentThread().isInterrupted()) {
            try {
                CanFrame frame = channel.read();
                processReceivedFrame(frame.getId(), frame.getDataLength(), frame);
            } catch (IOException e) {
                if (open) {
                    LOGGER.error("Error reading from CAN socket: {}", e.getMessage());
                    open = false;
                    notifyDisconnect(e);
                }
                break;
            } catch (Exception e) {
                if (open) {
                    LOGGER.warn("Unexpected error in CAN reader: {}", e.getMessage());
                }
            }
        }
    }

    /**
     * Processes a received CAN frame: applies the ID filter, serializes to the ring buffer,
     * and notifies the data listener. Package-private for testability.
     *
     * @param canId      the CAN identifier from the frame
     * @param dataLength the data payload length
     * @param frame      the native CAN frame (may be null in tests using the byte-array overload)
     */
    void processReceivedFrame(int canId, int dataLength, CanFrame frame) {
        if (!filter.matches(canId)) {
            return; // Skip frames that don't match the filter
        }

        byte[] frameBytes = new byte[CAN_FRAME_SIZE];
        frameBytes[0] = (byte) ((canId >> 24) & 0xFF);
        frameBytes[1] = (byte) ((canId >> 16) & 0xFF);
        frameBytes[2] = (byte) ((canId >> 8) & 0xFF);
        frameBytes[3] = (byte) (canId & 0xFF);
        frameBytes[4] = (byte) dataLength;
        if (frame != null && dataLength > 0) {
            byte[] data = new byte[dataLength];
            frame.getData(data, 0, dataLength);
            System.arraycopy(data, 0, frameBytes, 8, dataLength);
        }

        ringBuffer.write(frameBytes);

        Runnable listener = dataListener;
        if (listener != null) {
            listener.run();
        }
    }

    /**
     * Processes a received CAN frame from raw bytes (for testing without native frames).
     *
     * @param canId      the CAN identifier
     * @param data       the data payload
     */
    void processReceivedFrame(int canId, byte[] data) {
        if (!filter.matches(canId)) {
            return;
        }

        byte[] frameBytes = new byte[CAN_FRAME_SIZE];
        frameBytes[0] = (byte) ((canId >> 24) & 0xFF);
        frameBytes[1] = (byte) ((canId >> 16) & 0xFF);
        frameBytes[2] = (byte) ((canId >> 8) & 0xFF);
        frameBytes[3] = (byte) (canId & 0xFF);
        frameBytes[4] = (byte) data.length;
        System.arraycopy(data, 0, frameBytes, 8, Math.min(data.length, 8));

        ringBuffer.write(frameBytes);

        Runnable listener = dataListener;
        if (listener != null) {
            listener.run();
        }
    }

    @Override
    public boolean isOpen() {
        return open;
    }

    @Override
    public int getNumBytesAvailable() throws TransportException {
        if (!open) {
            throw new TransportException("Transport is closed");
        }
        readLock.lock();
        try {
            return ringBuffer.availableForReading();
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public byte[] peekReadableBytes(int numBytes) throws TransportException {
        if (!open) {
            throw new TransportException("Transport is closed");
        }
        readLock.lock();
        try {
            return ringBuffer.peek(numBytes);
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public byte[] read(int numBytes) throws TransportException {
        if (!open) {
            throw new TransportException("Transport is closed");
        }
        readLock.lock();
        try {
            return ringBuffer.read(numBytes);
        } finally {
            readLock.unlock();
        }
    }

    @Override
    public void write(byte[] bytes) throws TransportException {
        if (!open) {
            throw new TransportException("Transport is closed");
        }
        writeLock.lock();
        try {
            // Deserialize bytes to a CAN frame and send via channel
            if (bytes.length < CAN_FRAME_SIZE) {
                throw new TransportException("CAN frame must be at least " + CAN_FRAME_SIZE + " bytes");
            }
            int id = ((bytes[0] & 0xFF) << 24) | ((bytes[1] & 0xFF) << 16)
                    | ((bytes[2] & 0xFF) << 8) | (bytes[3] & 0xFF);
            int dlc = bytes[4] & 0xFF;
            byte[] data = new byte[Math.min(dlc, 8)];
            System.arraycopy(bytes, 8, data, 0, data.length);

            // Build and send the frame
            sendFrame(id, data);
        } catch (IOException e) {
            throw new TransportException("Failed to write CAN frame", e);
        } finally {
            writeLock.unlock();
        }
    }

    /**
     * Sends a CAN frame via the channel. Overridable for testing without native CAN support.
     *
     * @param id   the CAN identifier
     * @param data the frame data payload
     * @throws IOException if the write fails
     */
    protected void sendFrame(int id, byte[] data) throws IOException {
        CanFrame frame = CanFrame.create(id, (byte) 0, data);
        channel.write(frame);
    }

    @Override
    public void close() throws TransportException {
        if (!open) {
            return;
        }
        open = false;

        // Stop reader thread
        Thread reader = readerThread;
        if (reader != null) {
            reader.interrupt();
            try {
                reader.join(2000);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }

        // Release or close the channel
        if (sharedHandle != null) {
            sharedCanManager.releaseHandle(sharedHandle);
        } else if (channel != null) {
            try {
                channel.close();
            } catch (IOException e) {
                LOGGER.warn("Error closing CAN socket: {}", e.getMessage());
            }
        }

        // Notify disconnect listener
        notifyDisconnect(null);
    }

    @Override
    public void registerDataListener(Runnable listener) {
        this.dataListener = listener;
    }

    @Override
    public void removeDataListener() {
        this.dataListener = null;
    }

    @Override
    public void registerDisconnectListener(Consumer<Throwable> listener) {
        this.disconnectListener = listener;
    }

    @Override
    public void removeDisconnectListener() {
        this.disconnectListener = null;
    }

    /**
     * Notifies the registered disconnect listener, if any.
     *
     * @param cause the exception that caused the disconnect, or null for graceful close
     */
    private void notifyDisconnect(Throwable cause) {
        Consumer<Throwable> listener = disconnectListener;
        if (listener != null) {
            try {
                listener.accept(cause);
            } catch (Exception e) {
                LOGGER.warn("Error in disconnect listener: {}", e.getMessage());
            }
        }
    }
}
