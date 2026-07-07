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
package org.apache.plc4x.java.modbus.rtu;

import org.apache.plc4x.java.modbus.readwrite.ModbusErrorCode;
import org.apache.plc4x.java.modbus.readwrite.ModbusPDU;
import org.apache.plc4x.java.modbus.readwrite.ModbusPDUError;
import org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadHoldingRegistersRequest;
import org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadHoldingRegistersResponse;
import org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleRegisterResponse;
import org.apache.plc4x.java.modbus.readwrite.ModbusRtuADU;
import org.apache.plc4x.java.modbus.rtu.config.ModbusRtuConfiguration;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.BooleanSupplier;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNotSame;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Regression coverage for per-connection request serialization: Modbus over
 * serial is a single-outstanding-transaction protocol, so concurrent
 * {@code sendRequest} calls to the same unit id must not be allowed to race
 * on {@code pendingRequests} and receive each other's responses.
 */
class ModbusRtuConnectionRequestChainTest {

    @Test
    void concurrentSameUnitRequestsGetTheirOwnResponses() throws Exception {
        // Two requests to unit 1, fired without awaiting. Pre-fix,
        // pendingRequests.put overwrites request 1's future: response 1
        // completes request 2 (wrong data), request 1 only ever times out.
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        ModbusRtuConnection connection = newConnectedConnection(transport); // helper: construct + connect

        ModbusRtuADU request1 = readHoldingRegistersRequestAdu(1, 0x0000, 1); // helper building a real request ADU
        ModbusRtuADU request2 = readHoldingRegistersRequestAdu(1, 0x0010, 1);
        CompletableFuture<ModbusRtuADU> future1 = connection.sendRequest(request1, (short) 1);
        CompletableFuture<ModbusRtuADU> future2 = connection.sendRequest(request2, (short) 1);

        // Only request 1 may be on the wire so far (single outstanding
        // transaction); respond to it with a DISTINGUISHABLE payload.
        assertEquals(1, transport.writeCount(), "second request must wait for the first transaction");
        transport.deliver(readResponseFrame(1, new byte[]{0x11, 0x11}));
        transport.runDataListener(); // pump the codec

        ModbusRtuADU response1 = future1.get(2, TimeUnit.SECONDS);
        assertArrayEquals(new byte[]{0x11, 0x11}, extractRegisterBytes(response1)); // helper on the response PDU

        // Completing transaction 1 releases the chain: request 2 goes out.
        // (The chained send may run on the completing thread; wait briefly.)
        awaitTrue(() -> transport.writeCount() == 2, 2, TimeUnit.SECONDS);
        transport.deliver(readResponseFrame(1, new byte[]{0x22, 0x22}));
        transport.runDataListener();

        ModbusRtuADU response2 = future2.get(2, TimeUnit.SECONDS);
        assertArrayEquals(new byte[]{0x22, 0x22}, extractRegisterBytes(response2));
        assertNotSame(response1, response2);
    }

    @Test
    void failedTransactionDoesNotPoisonTheChain() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        ModbusRtuConnection connection = newConnectedConnection(transport);

        transport.failNextWrite(); // fake: next write() throws TransportException
        CompletableFuture<ModbusRtuADU> failed = connection.sendRequest(
            readHoldingRegistersRequestAdu(1, 0x0000, 1), (short) 1);
        assertThrows(Exception.class, () -> failed.get(2, TimeUnit.SECONDS));

        CompletableFuture<ModbusRtuADU> next = connection.sendRequest(
            readHoldingRegistersRequestAdu(1, 0x0000, 1), (short) 1);
        awaitTrue(() -> transport.writeCount() >= 1, 2, TimeUnit.SECONDS);
        transport.deliver(readResponseFrame(1, new byte[]{0x33, 0x33}));
        transport.runDataListener();
        assertNotNull(next.get(2, TimeUnit.SECONDS), "chain must continue past a failed transaction");
    }

    @Test
    void longQueuedBurstSurvivesWithoutStackOverflowOrStranding() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        ModbusRtuConnection connection = newConnectedConnection(transport);

        // Head request dispatches successfully and stays pending (no response
        // delivered yet) — this is what actually holds the chain open so a
        // backlog can build up behind it. (Note: setting failAllWrites() before
        // this call would make the head itself fail synchronously within its
        // own sendRequest() call, resolving the chain before the loop below
        // even starts — no backlog would ever form and the bug wouldn't be
        // exercised.)
        CompletableFuture<ModbusRtuADU> head = connection.sendRequest(
            readHoldingRegistersRequestAdu(1, 0x0000, 1), (short) 1);
        assertEquals(1, transport.writeCount());

        // Now make every subsequent write fail immediately at dispatch time,
        // then queue a long burst behind the still-pending head. Each of these
        // sendRequest calls only registers a whenComplete on the not-yet-done
        // chain tail and returns — none of them dispatch yet.
        transport.failAllWrites();
        List<CompletableFuture<ModbusRtuADU>> futures = new ArrayList<>();
        for (int i = 0; i < 10_000; i++) {
            futures.add(connection.sendRequest(readHoldingRegistersRequestAdu(1, i & 0xFFFF, 1), (short) 1));
        }

        // Release the chain: completing the head cascades synchronously through
        // every queued dispatch (each of which fails immediately), which is
        // exactly the recursive-stack scenario under test.
        head.completeExceptionally(new RuntimeException("induced head failure"));

        for (CompletableFuture<ModbusRtuADU> future : futures) {
            ExecutionException e = assertThrows(ExecutionException.class, () -> future.get(10, TimeUnit.SECONDS));
            assertNotNull(e.getCause(), "every queued request must complete exceptionally — none stranded");
        }
    }

    @Test
    void mismatchedFunctionCodeIsDiscardedAndPendingSurvives() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        ModbusRtuConnection connection = newConnectedConnection(transport);

        CompletableFuture<ModbusRtuADU> pendingRead = connection.sendRequest(
            readHoldingRegistersRequestAdu(1, 0x0000, 1), (short) 1); // fc 0x03

        // A write-echo frame (fc 0x06) for the same unit must NOT complete it.
        transport.deliver(wireBytes(new ModbusRtuADU((short) 1,
            new ModbusPDUWriteSingleRegisterResponse(0x0010, 0x1234))));
        transport.runDataListener();
        assertFalse(pendingRead.isDone(), "mismatched fc must be discarded, pending must survive");

        // The genuine 0x03 response still completes it.
        transport.deliver(readResponseFrame(1, new byte[]{0x11, 0x22}));
        transport.runDataListener();
        assertArrayEquals(new byte[]{0x11, 0x22},
            extractRegisterBytes(pendingRead.get(2, TimeUnit.SECONDS)));
    }

    @Test
    void exceptionResponseCompletesAnyPendingRequest() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        ModbusRtuConnection connection = newConnectedConnection(transport);

        CompletableFuture<ModbusRtuADU> pendingRead = connection.sendRequest(
            readHoldingRegistersRequestAdu(1, 0x0000, 1), (short) 1);

        transport.deliver(wireBytes(new ModbusRtuADU((short) 1,
            new ModbusPDUError(ModbusErrorCode.ILLEGAL_DATA_ADDRESS))));
        transport.runDataListener();

        ModbusRtuADU response = pendingRead.get(2, TimeUnit.SECONDS);
        assertTrue(response.getPdu().getErrorFlag(), "exception frames answer any pending request");
    }

    @Test
    void requestTimedOutWhileQueuedNeverReachesTheWire() throws Exception {
        ScriptedAsyncTransport transport = new ScriptedAsyncTransport();
        // Short request timeout so the queued request's total budget expires
        // while the head transaction is still outstanding.
        ModbusRtuConnection connection = newConnectedConnection(transport, 200);

        // Head request dispatches and stays pending (never answered).
        CompletableFuture<ModbusRtuADU> head = connection.sendRequest(
            readHoldingRegistersRequestAdu(1, 0x0000, 1), (short) 1);
        awaitTrue(() -> transport.writeCount() == 1, 2, TimeUnit.SECONDS);

        // Second request queues behind it and must time out WHILE QUEUED.
        CompletableFuture<ModbusRtuADU> queued = connection.sendRequest(
            readHoldingRegistersRequestAdu(1, 0x0010, 1), (short) 1);
        assertThrows(ExecutionException.class, () -> queued.get(2, TimeUnit.SECONDS),
            "queued request must time out from its enqueue-time clock");

        // Head times out too (same 200ms clock); the chain then dispatches the
        // dead queued request — which must be SKIPPED: no second wire write.
        assertThrows(ExecutionException.class, () -> head.get(2, TimeUnit.SECONDS));
        Thread.sleep(200); // give the chain's async dispatch a beat
        assertEquals(1, transport.writeCount(),
            "a request that died in the queue must never reach the wire");
    }

    /**
     * Builds a connection wired to the given fake transport, and drives it
     * through the same construction/connect path as {@code ModbusRtuConnectionTest}.
     */
    private static ModbusRtuConnection newConnectedConnection(ScriptedAsyncTransport transport) throws Exception {
        return newConnectedConnection(transport, 5000);
    }

    private static ModbusRtuConnection newConnectedConnection(ScriptedAsyncTransport transport, int requestTimeoutMs) throws Exception {
        ModbusRtuConfiguration config = new ModbusRtuConfiguration();
        config.setRequestTimeout(requestTimeoutMs);
        config.setDefaultUnitIdentifier(1);
        config.setPingAddress("4x00001:BOOL");
        config.setDefaultPayloadByteOrder(ModbusByteOrder.BIG_ENDIAN);
        config.setMaxCoilsPerRequest(2000);
        config.setMaxRegistersPerRequest(125);

        AuditLog auditLog = mock(AuditLog.class);
        when(auditLog.isEnabled()).thenReturn(false);

        ModbusRtuConnection connection = new ModbusRtuConnection(config, transport, auditLog);
        connection.connect();
        return connection;
    }

    private static ModbusRtuADU readHoldingRegistersRequestAdu(int unit, int address, int count) {
        return new ModbusRtuADU((short) unit, new ModbusPDUReadHoldingRegistersRequest(address, count));
    }

    private static byte[] extractRegisterBytes(ModbusRtuADU responseAdu) {
        ModbusPDU pdu = responseAdu.getPdu();
        return ((ModbusPDUReadHoldingRegistersResponse) pdu).getValue();
    }

    private static byte[] wireBytes(ModbusRtuADU adu) throws Exception {
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[adu.getLengthInBytes()]);
        adu.serialize(writeBuffer);
        return writeBuffer.getBytes();
    }

    private static byte[] readResponseFrame(int unitId, byte[] registers) throws Exception {
        return wireBytes(new ModbusRtuADU((short) unitId, new ModbusPDUReadHoldingRegistersResponse(registers)));
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
     * In-memory transport double supporting the async listener-registration
     * contract the connection relies on ({@code startReceiving} registers a
     * data listener rather than polling). Bytes are queued via
     * {@code deliver(...)} and consumed through the real TransportInstance
     * contract, mirroring {@code ModbusRtuMessageCodecTest.ScriptedTransport}.
     */
    static final class ScriptedAsyncTransport implements AsyncTransportInstance<TransportConfiguration> {
        private final ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        private int readPosition;
        private boolean open = true;
        private final AtomicInteger writeCount = new AtomicInteger();
        private final AtomicBoolean failNextWrite = new AtomicBoolean(false);
        private final AtomicBoolean failAllWrites = new AtomicBoolean(false);
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

        void failAllWrites() {
            failAllWrites.set(true);
        }

        void runDataListener() {
            Runnable listener = dataListener.get();
            if (listener == null) {
                throw new IllegalStateException("No data listener registered");
            }
            listener.run();
        }

        private byte[] bufferedBytes() {
            return buffer.toByteArray();
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
            byte[] all = bufferedBytes();
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
            if (failAllWrites.get() || failNextWrite.compareAndSet(true, false)) {
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
