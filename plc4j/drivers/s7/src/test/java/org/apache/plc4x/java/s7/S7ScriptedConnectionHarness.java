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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageResponseData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterSetupCommunication;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.io.ByteArrayOutputStream;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;

/**
 * Brings up a real {@link S7CotpConnection} against an in-memory transport, so behaviour that
 * only shows up on a connected driver (per-tag response codes, what actually reaches the wire)
 * can be tested without a PLC.
 * <p>
 * Connecting an S7 driver is a three step conversation - COTP connection request, S7 setup
 * communication, and an SZL identification probe - and the connection blocks until it is done.
 * A responder thread answers each outgoing frame with a canned reply, which is enough to get
 * the connection into its normal operating state.
 */
public final class S7ScriptedConnectionHarness {

    private S7ScriptedConnectionHarness() {
    }

    /**
     * Connects an S7 connection wired to the given transport, scripting the handshake.
     */
    public static S7CotpConnection newConnectedConnection(ScriptedS7Transport transport) throws Exception {
        S7Configuration configuration = new S7Configuration();
        // Pinning the controller type makes the driver skip the SZL capability probe, so the
        // handshake is just COTP connect + S7 setup communication.
        configuration.setControllerType(ControllerType.S7_300);
        S7CotpConnection connection = new S7CotpConnection(configuration, transport,
            AuditLog.builder().build());

        // The handshake blocks, so it has to run while the responder feeds it answers.
        AtomicReference<Exception> failure = new AtomicReference<>();
        Thread connectThread = new Thread(() -> {
            try {
                connection.connect();
            } catch (Exception e) {
                failure.set(e);
            }
        }, "s7-harness-connect");
        connectThread.setDaemon(true);

        AtomicBoolean handshakeDone = new AtomicBoolean(false);
        Thread responder = new Thread(() -> {
            int answered = 0;
            while (!handshakeDone.get()) {
                if (transport.writeCount() > answered) {
                    byte[] reply = replyFor(answered, transport.writtenFrames().get(answered));
                    answered++;
                    if (reply != null) {
                        if (Boolean.getBoolean("s7harness.debug")) {
                            System.out.println("HARNESS: replying to frame " + (answered - 1)
                                + " with " + java.util.HexFormat.of().formatHex(reply)
                                + " (listener registered: " + transport.hasDataListener() + ")");
                        }
                        transport.deliver(reply);
                        transport.runDataListener();
                    }
                }
                try {
                    Thread.sleep(5);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    return;
                }
            }
        }, "s7-harness-responder");
        responder.setDaemon(true);

        responder.start();
        connectThread.start();
        connectThread.join(15_000);
        handshakeDone.set(true);

        if (failure.get() != null) {
            throw failure.get();
        }
        if (connectThread.isAlive()) {
            throw new IllegalStateException("S7 handshake did not complete");
        }
        transport.resetCounters();
        return connection;
    }

    /**
     * The canned reply for the n-th frame the driver sends during the handshake. The third
     * frame is the SZL capability probe, which the driver is designed to survive failing - it
     * just marks UserData services unavailable - so it is left unanswered.
     */
    private static byte[] replyFor(int frameIndex, byte[] request) {
        try {
            if (frameIndex == 0) {
                // Responses are correlated by tpdu reference, so echo the one the driver used.
                return s7SetupCommunicationResponse(tpduReferenceOf(request));
            }
            return null;
        } catch (Exception e) {
            throw new IllegalStateException("Could not build handshake reply " + frameIndex, e);
        }
    }

    /** Digs the S7 tpdu reference out of a frame the driver just sent. */
    private static int tpduReferenceOf(byte[] frame) throws Exception {
        return S7Message.staticParse(readBuffer(frame)).getTpduReference();
    }

    /** Same, for tests that need to answer a specific request they watched go out. */
    public static int tpduReferenceOfFrame(byte[] frame) throws Exception {
        return tpduReferenceOf(frame);
    }

    /**
     * A response carrying an error class/code in the S7 header and nothing else, which is what a
     * controller sends when it refuses the whole request rather than individual items.
     */
    public static byte[] headerErrorResponse(int tpduReference, int errorClass, int errorCode) throws Exception {
        return wireBytes(new S7MessageResponseData(tpduReference, null, null,
            (short) errorClass, (short) errorCode));
    }

    private static byte[] s7SetupCommunicationResponse(int tpduReference) throws Exception {
        S7Message message = new S7MessageResponseData(tpduReference,
            new S7ParameterSetupCommunication(8, 8, 240), null, (short) 0, (short) 0);
        return wireBytes(message);
    }

    /**
     * The COTP/TPKT framing is done by the transport this driver sits on, so what the driver
     * itself reads and writes are plain S7 messages.
     */
    private static byte[] wireBytes(S7Message message) throws Exception {
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[message.getLengthInBytes()],
            ENCODING_OPTIONS);
        message.serialize(writeBuffer);
        return writeBuffer.getBytes();
    }

    private static ReadBufferByteBased readBuffer(byte[] data) {
        return new ReadBufferByteBased(data, ENCODING_OPTIONS);
    }

    /**
     * The generated S7 parsers/serializers require these to be set on the buffer - the driver's
     * message codec configures the very same ones.
     */
    private static final WithOption[] ENCODING_OPTIONS = {
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithByteBasedOption.WithByteOrder("BIG_ENDIAN")
    };

    /**
     * In-memory transport double: records everything the driver writes and lets a test feed
     * bytes back in. Mirrors the scripted transport used by the Modbus tests.
     */
    public static final class ScriptedS7Transport implements AsyncTransportInstance<TransportConfiguration> {
        private final ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        private int readPosition;
        private boolean open = true;
        private final AtomicInteger writeCount = new AtomicInteger();
        private final List<byte[]> writtenFrames = new CopyOnWriteArrayList<>();
        private final AtomicReference<Runnable> dataListener = new AtomicReference<>();

        public void deliver(byte[] bytes) {
            synchronized (buffer) {
                buffer.writeBytes(bytes);
            }
        }

        public int writeCount() {
            return writeCount.get();
        }

        public List<byte[]> writtenFrames() {
            return writtenFrames;
        }

        /** Forgets the handshake traffic so a test only sees what its own request produced. */
        public void resetCounters() {
            writeCount.set(0);
            writtenFrames.clear();
        }

        public boolean hasDataListener() {
            return dataListener.get() != null;
        }

        public void runDataListener() {
            Runnable listener = dataListener.get();
            if (listener != null) {
                listener.run();
            }
        }

        @Override
        public TransportConfiguration getConfiguration() {
            return null;
        }

        @Override
        public boolean isOpen() {
            return open;
        }

        @Override
        public int getNumBytesAvailable() {
            synchronized (buffer) {
                return buffer.size() - readPosition;
            }
        }

        @Override
        public byte[] peekReadableBytes(int numBytes) throws TransportException {
            synchronized (buffer) {
                if (numBytes > getNumBytesAvailable()) {
                    throw new TransportException("peek beyond available: " + numBytes);
                }
                byte[] all = buffer.toByteArray();
                byte[] result = new byte[numBytes];
                System.arraycopy(all, readPosition, result, 0, numBytes);
                return result;
            }
        }

        @Override
        public byte[] read(int numBytes) throws TransportException {
            synchronized (buffer) {
                byte[] result = peekReadableBytes(numBytes);
                readPosition += numBytes;
                return result;
            }
        }

        @Override
        public void write(byte[] bytes) {
            writtenFrames.add(bytes.clone());
            writeCount.incrementAndGet();
        }

        @Override
        public void close() {
            open = false;
        }

        @Override
        public void registerDataListener(Runnable listener) {
            dataListener.set(listener);
        }

        @Override
        public void removeDataListener() {
            dataListener.set(null);
        }
    }
}
