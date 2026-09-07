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

import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.slmp.config.SlmpConfiguration;
import org.apache.plc4x.java.slmp.readwrite.SlmpResponseFrame3E;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.BooleanSupplier;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.fail;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Happy-path and error-mapping coverage for the SLMP Batch Write (0x1401) path
 * in {@link SlmpConnection}: a successful device acknowledgement maps to OK, a
 * non-zero endCode maps to REMOTE_ERROR, and a device that never answers times
 * out to REMOTE_ERROR for that tag. Mirrors {@code SlmpConnectionFailurePathTest}.
 */
class SlmpConnectionWritePathTest {

    @Test
    void successfulWriteMapsToOk() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 5000);
        CompletableFuture<? extends PlcWriteResponse> future = connection.writeRequestBuilder()
            .addTagAddress("v", "D350", 0x1234)
            .build().execute();
        awaitTrue(() -> transport.writeCount() == 1, 2, TimeUnit.SECONDS);
        transport.deliver(responseFrame(0x0000, new byte[0])); // Batch Write success: empty payload
        transport.runDataListener();
        PlcWriteResponse response = future.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.OK, response.getResponseCode("v"));
    }

    @Test
    void deviceErrorMapsToRemoteError() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 5000);
        CompletableFuture<? extends PlcWriteResponse> future = connection.writeRequestBuilder()
            .addTagAddress("v", "D350", 0x1234)
            .build().execute();
        awaitTrue(() -> transport.writeCount() == 1, 2, TimeUnit.SECONDS);
        transport.deliver(responseFrame(0xC059, new byte[0]));
        transport.runDataListener();
        PlcWriteResponse response = future.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.REMOTE_ERROR, response.getResponseCode("v"));
    }

    @Test
    void timedOutWriteIsRemoteError() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 200);
        CompletableFuture<? extends PlcWriteResponse> future = connection.writeRequestBuilder()
            .addTagAddress("v", "D350", 0x1234)
            .build().execute();
        PlcWriteResponse response = future.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.REMOTE_ERROR, response.getResponseCode("v"));
    }

    @Test
    void builderErrorItemsAreEchoedNotMaskedAsInternalError() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        SlmpConnection connection = newConnectedConnection(transport, 5000);

        CompletableFuture<? extends PlcWriteResponse> future = connection.writeRequestBuilder()
            .addTagAddress("badValue", "D350", 70000)       // out of WORD range -> builder error item
            .addTagAddress("badAddress", "Z99", 1)          // unparseable device -> builder error item
            .addTagAddress("good", "D351", 0x1234)
            .build().execute();

        awaitTrue(() -> transport.writeCount() == 1, 2, TimeUnit.SECONDS);   // only the good tag reaches the wire
        transport.deliver(responseFrame(0x0000, new byte[0]));
        transport.runDataListener();

        PlcWriteResponse response = future.get(5, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.INVALID_DATA, response.getResponseCode("badValue"));
        assertEquals(PlcResponseCode.INVALID_ADDRESS, response.getResponseCode("badAddress"));
        assertEquals(PlcResponseCode.OK, response.getResponseCode("good"));
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
        private final AtomicReference<Runnable> dataListener = new AtomicReference<>();

        void deliver(byte[] bytes) {
            buffer.writeBytes(bytes);
        }

        int writeCount() {
            return writeCount.get();
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
