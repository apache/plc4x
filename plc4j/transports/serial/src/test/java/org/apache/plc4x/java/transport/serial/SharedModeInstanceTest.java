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
import org.apache.plc4x.java.transport.serial.config.SerialTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.Function;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 * Instance-level integration test for shared-port ("reuse-port") mode: two
 * real {@link SerialTransportInstance}s backed by a single mocked physical
 * port, verifying that both receive a broadcast, that closing one instance
 * does not disturb the other, and that the physical port is only closed once
 * the last instance releases it.
 */
class SharedModeInstanceTest {

    private static SerialTransportConfiguration config() {
        SerialTransportConfiguration config = new SerialTransportConfiguration();
        config.baudRate = 9600;
        config.dataBits = 8;
        config.stopBits = 1;
        config.parity = "NONE";
        config.flowControl = "NONE";
        config.readTimeout = 1000;
        config.writeTimeout = 1000;
        config.reusePort = true;
        return config;
    }

    /**
     * Mirrors SharedPortBroadcastTest's MockPortFactory, extended with a
     * one-shot "broadcast" mechanism. All stubs are attached once, at mock
     * creation time (i.e. before the SharedPort's reader thread is started
     * by SharedSerialPortManager.acquirePort()), so there is no live
     * re-stubbing of a mock the reader thread might be concurrently calling.
     */
    static final class MockPortFactory implements Function<String, SerialPort> {
        final Map<String, SerialPort> created = new HashMap<>();
        private final AtomicInteger bytesAvailable = new AtomicInteger(0);
        private final AtomicReference<byte[]> pendingPayload = new AtomicReference<>();

        @Override
        public SerialPort apply(String portName) {
            SerialPort port = mock(SerialPort.class);
            when(port.openPort()).thenReturn(true);
            // SerialTransportInstance.isOpen() checks `open && port.isOpen()`;
            // without this stub the mock's default `false` would make
            // getNumBytesAvailable()/read() short-circuit to empty forever.
            when(port.isOpen()).thenReturn(true);
            when(port.getOutputStream()).thenReturn(new ByteArrayOutputStream());
            // Stub readBytes first, then bytesAvailable, both at mock-creation
            // time - before the reader thread exists - so there is no window
            // where the reader could observe a non-zero bytesAvailable with
            // an unstubbed (default) readBytes.
            when(port.readBytes(any(byte[].class), anyInt())).thenAnswer(invocation -> {
                byte[] buffer = invocation.getArgument(0);
                byte[] payload = pendingPayload.get();
                if (payload == null) {
                    return 0;
                }
                System.arraycopy(payload, 0, buffer, 0, payload.length);
                return payload.length;
            });
            when(port.bytesAvailable()).thenAnswer(inv -> bytesAvailable.getAndSet(0));
            created.put(portName, port);
            return port;
        }

        /**
         * Arms a one-shot broadcast: the next poll of the shared reader
         * thread will read exactly this payload and hand it to all
         * subscribers, then go quiet again.
         */
        void broadcast(byte[] payload) {
            pendingPayload.set(payload);
            bytesAvailable.set(payload.length);
        }
    }

    private static void awaitBytesAvailable(SerialTransportInstance instance, int expected) throws Exception {
        long deadline = System.currentTimeMillis() + 5000;
        while (System.currentTimeMillis() < deadline) {
            if (instance.getNumBytesAvailable() >= expected) {
                return;
            }
            Thread.sleep(10);
        }
        fail("Timed out waiting for " + expected + " bytes to become available");
    }

    /**
     * Functional interface mirroring {@link java.util.function.BooleanSupplier}
     * but allowing checked exceptions, so conditions can call methods like
     * {@link SerialTransportInstance#getNumBytesAvailable()} directly.
     */
    @FunctionalInterface
    private interface ThrowingBooleanSupplier {
        boolean getAsBoolean() throws Exception;
    }

    /**
     * Generic poll-sleep helper: repeatedly evaluates {@code condition} until
     * it returns true or {@code timeout} elapses, failing the test otherwise.
     */
    private static void awaitTrue(ThrowingBooleanSupplier condition, long timeout, TimeUnit unit) throws Exception {
        long deadlineNanos = System.nanoTime() + unit.toNanos(timeout);
        while (System.nanoTime() < deadlineNanos) {
            if (condition.getAsBoolean()) {
                return;
            }
            Thread.sleep(10);
        }
        fail("Timed out waiting for condition to become true");
    }

    @Test
    void twoSharedInstancesReceiveBroadcastsAndCloseIndependently() throws Exception {
        MockPortFactory factory = new MockPortFactory();
        SharedSerialPortManager manager = new SharedSerialPortManager(factory);

        SerialTransportInstance first = new SerialTransportInstance(
            manager, "COMSHARED", config(), AuditLog.builder().build());
        SerialTransportInstance second = new SerialTransportInstance(
            manager, "COMSHARED", config(), AuditLog.builder().build());

        assertEquals(1, factory.created.size(), "both instances must share one physical port");
        SerialPort mockPort = factory.created.get("COMSHARED");

        byte[] payload = {0x01, 0x02, 0x03};
        factory.broadcast(payload);

        awaitBytesAvailable(first, payload.length);
        awaitBytesAvailable(second, payload.length);

        assertArrayEquals(payload, first.read(payload.length), "first instance must receive the broadcast");
        assertArrayEquals(payload, second.read(payload.length), "second instance must receive the broadcast");

        // Closing the first instance must not affect the second: no join
        // stall (Finding 1) and the physical port must stay open since the
        // second instance still holds a reference.
        first.close();
        assertFalse(first.isOpen());
        verify(mockPort, never()).closePort();

        byte[] secondPayload = {0x04, 0x05};
        factory.broadcast(secondPayload);

        awaitBytesAvailable(second, secondPayload.length);
        assertArrayEquals(secondPayload, second.read(secondPayload.length),
            "second instance must keep receiving broadcasts after the first instance closed");

        // Closing the last instance must release the physical port.
        second.close();
        assertFalse(second.isOpen());
        verify(mockPort).closePort();
    }

    /**
     * Regression test for shared-mode dispatch isolation: the shared reader
     * thread iterates all subscribers' onData() inline (SharedPort.readFromPort()).
     * If one instance's dataListener blocks, an inline dispatch would stall
     * the reader thread before it ever reaches the next subscriber, starving
     * every other connection sharing the physical port.
     */
    @Test
    void blockedListenerOnOneConnectionDoesNotStallTheOther() throws Exception {
        MockPortFactory factory = new MockPortFactory();
        SharedSerialPortManager manager = new SharedSerialPortManager(factory);

        SerialTransportInstance first = new SerialTransportInstance(
            manager, "COMSHARED", config(), AuditLog.builder().build());
        SerialTransportInstance second = new SerialTransportInstance(
            manager, "COMSHARED", config(), AuditLog.builder().build());

        try {
            CountDownLatch blockA = new CountDownLatch(1);
            CountDownLatch aEntered = new CountDownLatch(1);
            AtomicInteger bNotifications = new AtomicInteger();

            first.registerDataListener(() -> {
                aEntered.countDown();
                try {
                    blockA.await(10, TimeUnit.SECONDS); // A's callback wedges
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            });
            second.registerDataListener(bNotifications::incrementAndGet);

            factory.broadcast(new byte[]{0x01, 0x02, 0x03, 0x04});
            assertTrue(aEntered.await(5, TimeUnit.SECONDS), "A's listener must have been invoked");

            // While A is wedged, B must still be notified and readable.
            awaitTrue(() -> bNotifications.get() >= 1, 5, TimeUnit.SECONDS);
            awaitTrue(() -> second.getNumBytesAvailable() >= 4, 5, TimeUnit.SECONDS);

            // A second broadcast must still reach B while A stays wedged —
            // both the ring content AND a second dispatch notification
            // (pins the coalescing flag being re-armed after each run).
            factory.broadcast(new byte[]{0x05, 0x06, 0x07, 0x08});
            awaitTrue(() -> second.getNumBytesAvailable() >= 8, 5, TimeUnit.SECONDS);
            awaitTrue(() -> bNotifications.get() >= 2, 5, TimeUnit.SECONDS);

            blockA.countDown(); // release A; its coalesced dispatch drains
            awaitTrue(() -> first.getNumBytesAvailable() >= 8, 5, TimeUnit.SECONDS);
        } finally {
            first.close();
            second.close();
        }
    }
}
