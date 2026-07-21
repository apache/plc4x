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
package org.apache.plc4x.java.slmp;

import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.slmp.config.SlmpConfiguration;
import org.apache.plc4x.java.slmp.readwrite.SlmpResponseFrame3E;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.BooleanSupplier;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.fail;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Failure-path coverage for the single-slot request/response correlation in
 * {@link SlmpConnection}: 3E frames carry no correlation id, so a timeout, a
 * lost transport or an unsolicited frame must degrade exactly one tag (or
 * nothing) instead of poisoning the connection. Mirrors the style of
 * {@code ModbusRtuConnectionRequestChainTest}.
 */
class SlmpConnectionFailurePathTest {

    @Test
    void timedOutTagIsIsolatedToRemoteErrorWhileOthersSucceed() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 200);

        CompletableFuture<? extends PlcReadResponse> responseFuture = connection.readRequestBuilder()
            .addTagAddress("answered", "D350")
            .addTagAddress("silent", "D351")
            .build().execute();

        // Sequential per-tag chain: only the first read may be on the wire.
        awaitTrue(() -> transport.writeCount() == 1, 2, TimeUnit.SECONDS);
        transport.deliver(responseFrame(0x0000, new byte[]{(byte) 0xAB, 0x56}));
        transport.runDataListener();

        // The second read dispatches once the first completes — and is never answered.
        awaitTrue(() -> transport.writeCount() == 2, 2, TimeUnit.SECONDS);

        PlcReadResponse response = responseFuture.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.OK, response.getResponseCode("answered"));
        assertEquals(0x56AB, response.getPlcValue("answered").getInteger());
        assertEquals(PlcResponseCode.REMOTE_ERROR, response.getResponseCode("silent"),
            "a device that does not answer in time must surface as REMOTE_ERROR for that tag only");
    }

    @Test
    void transportDisconnectFailsThePendingTagAsInternalError() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 5000);

        CompletableFuture<? extends PlcReadResponse> responseFuture = connection.readRequestBuilder()
            .addTagAddress("pending", "D350")
            .build().execute();
        awaitTrue(() -> transport.writeCount() == 1, 2, TimeUnit.SECONDS);

        connection.onTransportDisconnected(new PlcRuntimeException("link down"));

        PlcReadResponse response = responseFuture.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.INTERNAL_ERROR, response.getResponseCode("pending"),
            "losing the transport must fail the in-flight tag instead of stranding its future");
    }

    @Test
    void failedSendFailsOnlyThatTag() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 5000);

        transport.failNextWrite();
        CompletableFuture<? extends PlcReadResponse> responseFuture = connection.readRequestBuilder()
            .addTagAddress("unsendable", "D350")
            .build().execute();

        PlcReadResponse response = responseFuture.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.INTERNAL_ERROR, response.getResponseCode("unsendable"));
    }

    @Test
    void unsolicitedFrameIsIgnoredAndTheConnectionStaysUsable() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 5000);

        // A frame nobody asked for must be dropped (logged), not blow up the receiver.
        transport.deliver(responseFrame(0x0000, new byte[]{0x01, 0x02}));
        transport.runDataListener();

        // A regular read afterwards still works — the slot was not poisoned.
        CompletableFuture<? extends PlcReadResponse> responseFuture = connection.readRequestBuilder()
            .addTagAddress("after", "D350")
            .build().execute();
        awaitTrue(() -> transport.writeCount() == 1, 2, TimeUnit.SECONDS);
        transport.deliver(responseFrame(0x0000, new byte[]{(byte) 0xAB, 0x56}));
        transport.runDataListener();

        PlcReadResponse response = responseFuture.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.OK, response.getResponseCode("after"));
        assertEquals(0x56AB, response.getPlcValue("after").getInteger());
    }

    @Test
    void connectRejectsMonitoringTimerOutsideUnsigned16Range() {
        SlmpConfiguration config = new SlmpConfiguration();
        config.setRequestTimeout(5000);
        config.setMonitoringTimer(0x1_0000); // one past the uint16 ceiling serialized into the 3E frame

        SlmpConnection connection = newDisconnectedConnection(config);
        assertThrows(PlcConnectionException.class, connection::connect,
            "an out-of-range monitoring-timer must be rejected at connect, before it can be truncated on the wire");
    }

    @Test
    void connectRejectsNonPositiveRequestTimeout() {
        SlmpConfiguration config = new SlmpConfiguration();
        config.setRequestTimeout(0); // orTimeout(0) would time out every request immediately
        config.setMonitoringTimer(0x0000);

        SlmpConnection connection = newDisconnectedConnection(config);
        assertThrows(PlcConnectionException.class, connection::connect,
            "a non-positive request-timeout must be rejected at connect rather than failing every read");
    }

    @Test
    void builderErrorItemIsEchoedForReads() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 5000);
        CompletableFuture<? extends PlcReadResponse> future = connection.readRequestBuilder()
            .addTagAddress("bad", "Z99")     // unparseable device -> builder error item
            .addTagAddress("good", "D350")
            .build().execute();
        awaitTrue(() -> transport.writeCount() == 1, 2, TimeUnit.SECONDS);   // only the good tag reaches the wire
        transport.deliver(responseFrame(0x0000, new byte[]{(byte) 0xAB, 0x56}));
        transport.runDataListener();
        PlcReadResponse response = future.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.INVALID_ADDRESS, response.getResponseCode("bad"));
        assertEquals(PlcResponseCode.OK, response.getResponseCode("good"));
    }

    private static SlmpConnection newDisconnectedConnection(SlmpConfiguration config) {
        AuditLog auditLog = mock(AuditLog.class);
        when(auditLog.isEnabled()).thenReturn(false);
        return new SlmpConnection(config, new ScriptedAsyncTransport(), auditLog);
    }

    private static SlmpConnection newConnectedConnection(ScriptedAsyncTransport transport, int requestTimeoutMs)
            throws Exception {
        SlmpConfiguration config = new SlmpConfiguration();
        config.setRequestTimeout(requestTimeoutMs);
        config.setMonitoringTimer(0x0000);

        AuditLog auditLog = mock(AuditLog.class);
        when(auditLog.isEnabled()).thenReturn(false);

        SlmpConnection connection = new SlmpConnection(config, transport, auditLog);
        connection.connect();
        return connection;
    }

    private static byte[] responseFrame(int endCode, byte[] data) throws Exception {
        SlmpResponseFrame3E frame = new SlmpResponseFrame3E(endCode, data);
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[frame.getLengthInBytes()]);
        frame.serialize(writeBuffer);
        return writeBuffer.getBytes();
    }

    private static void awaitTrue(BooleanSupplier condition, long timeout, TimeUnit unit) throws InterruptedException {
        long deadline = System.currentTimeMillis() + unit.toMillis(timeout);
        while (!condition.getAsBoolean()) {
            if (System.currentTimeMillis() > deadline) {
                fail("condition not met within timeout");
            }
            Thread.sleep(10);
        }
    }

    /**
     * In-memory transport double implementing the async listener-registration
     * contract {@code ConnectionBase.startReceiving} relies on. Bytes queued
     * via {@link #deliver(byte[])} are consumed through the real
     * TransportInstance read contract by the codec.
     */
    static final class ScriptedAsyncTransport implements AsyncTransportInstance<TransportConfiguration> {
        private final ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        private int readPosition;
        private boolean open = true;
        private final AtomicInteger writeCount = new AtomicInteger();
        private final AtomicBoolean failNextWrite = new AtomicBoolean(false);
        private final AtomicReference<Runnable> dataListener = new AtomicReference<>();

        void deliver(byte[] bytes) {
            buffer.writeBytes(bytes);
        }

        int writeCount() {
            return writeCount.get();
        }

        void failNextWrite() {
            failNextWrite.set(true);
        }

        void runDataListener() {
            Runnable listener = dataListener.get();
            if (listener == null) {
                throw new IllegalStateException("No data listener registered");
            }
            listener.run();
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
            return buffer.size() - readPosition;
        }

        @Override
        public byte[] peekReadableBytes(int numBytes) throws TransportException {
            if (numBytes > getNumBytesAvailable()) {
                throw new TransportException("peek beyond available: " + numBytes);
            }
            byte[] all = buffer.toByteArray();
            byte[] result = new byte[numBytes];
            System.arraycopy(all, readPosition, result, 0, numBytes);
            return result;
        }

        @Override
        public byte[] read(int numBytes) throws TransportException {
            byte[] result = peekReadableBytes(numBytes);
            readPosition += numBytes;
            return result;
        }

        @Override
        public void write(byte[] bytes) throws TransportException {
            writeCount.incrementAndGet();
            if (failNextWrite.compareAndSet(true, false)) {
                throw new TransportException("scripted write failure");
            }
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
