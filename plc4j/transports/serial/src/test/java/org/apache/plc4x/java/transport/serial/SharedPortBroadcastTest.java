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
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class SharedPortBroadcastTest {

    private static SharedSerialPortManager.SerialPortConfig config(int baudRate) {
        return new SharedSerialPortManager.SerialPortConfig(
            baudRate, 8, 1, SerialPort.NO_PARITY, SerialPort.FLOW_CONTROL_DISABLED,
            1000, 1000, false, false, 0);
    }

    static final class MockPortFactory implements java.util.function.Function<String, SerialPort> {
        final Map<String, SerialPort> created = new HashMap<>();

        @Override
        public SerialPort apply(String portName) {
            SerialPort port = mock(SerialPort.class);
            when(port.openPort()).thenReturn(true);
            when(port.getOutputStream()).thenReturn(new ByteArrayOutputStream());
            when(port.bytesAvailable()).thenReturn(0);
            created.put(portName, port);
            return port;
        }
    }

    @Test
    void identicalConfigSharesOnePhysicalPort() throws Exception {
        MockPortFactory factory = new MockPortFactory();
        SharedSerialPortManager manager = new SharedSerialPortManager(factory);

        SharedSerialPortManager.SharedPort first = manager.acquirePort("COMX", config(9600));
        SharedSerialPortManager.SharedPort second = manager.acquirePort("COMX", config(9600));

        assertSame(first, second, "identical config must share the entry");
        assertEquals(1, factory.created.size(), "only one physical open");
    }

    @Test
    void configMismatchThrowsNamingThePort() throws Exception {
        SharedSerialPortManager manager = new SharedSerialPortManager(new MockPortFactory());
        manager.acquirePort("COMY", config(9600));

        TransportException e = assertThrows(TransportException.class,
            () -> manager.acquirePort("COMY", config(19200)));
        assertTrue(e.getMessage().contains("COMY"), "error must name the port: " + e.getMessage());
    }

    @Test
    void lastReleaseClosesThePort() throws Exception {
        MockPortFactory factory = new MockPortFactory();
        SharedSerialPortManager manager = new SharedSerialPortManager(factory);
        SharedSerialPortManager.SharedPort shared = manager.acquirePort("COMZ", config(9600));
        manager.acquirePort("COMZ", config(9600));

        manager.releasePort(shared);
        verify(factory.created.get("COMZ"), never()).closePort();

        manager.releasePort(shared);
        verify(factory.created.get("COMZ")).closePort();

        // A fresh acquire after full release opens a new physical port.
        manager.acquirePort("COMZ", config(9600));
        assertEquals(1, factory.created.size(), "map replaced same-name entry");
    }

    static final class CollectingSubscriber implements SharedPortSubscriber {
        final ByteArrayOutputStream received = new ByteArrayOutputStream();
        final java.util.concurrent.CountDownLatch dataLatch = new java.util.concurrent.CountDownLatch(1);
        final java.util.concurrent.CountDownLatch failureLatch = new java.util.concurrent.CountDownLatch(1);
        volatile Throwable failure;

        @Override
        public synchronized void onData(byte[] data, int offset, int length) {
            received.write(data, offset, length);
            dataLatch.countDown();
        }

        @Override
        public void onFailure(Throwable cause) {
            failure = cause;
            failureLatch.countDown();
        }

        synchronized byte[] bytes() {
            return received.toByteArray();
        }
    }

    @Test
    void readerBroadcastsToAllSubscribers() throws Exception {
        MockPortFactory factory = new MockPortFactory();
        SharedSerialPortManager manager = new SharedSerialPortManager(factory);
        SharedSerialPortManager.SharedPort shared = manager.acquirePort("COMB", config(9600));
        SerialPort port = factory.created.get("COMB");

        CollectingSubscriber first = new CollectingSubscriber();
        CollectingSubscriber second = new CollectingSubscriber();
        shared.addSubscriber(first);
        shared.addSubscriber(second);

        byte[] payload = {0x01, 0x02, 0x03};
        // The shared reader thread is already running and polling the mock
        // concurrently. Mockito's stubbing bookkeeping (which invocation a
        // "when(...)" call attaches its answer to) is shared per-mock, not
        // thread-local, so restubbing here while the reader thread is free
        // to call the same mock races: a reader-thread invocation of
        // bytesAvailable() can interleave with the test thread evaluating
        // "when(port.readBytes(...))" and end up attached to the wrong
        // invocation, corrupting the mock's stub table. Take SharedPort's
        // own monitor (the same lock readFromPort() synchronizes on) to
        // block the reader for the duration of the stubbing calls, so the
        // reorder below is deterministic rather than merely probable.
        synchronized (shared) {
            when(port.readBytes(any(byte[].class), anyInt())).thenAnswer(invocation -> {
                byte[] buffer = invocation.getArgument(0);
                System.arraycopy(payload, 0, buffer, 0, payload.length);
                return payload.length;
            });
            when(port.bytesAvailable()).thenReturn(payload.length, 0);
        }

        assertTrue(first.dataLatch.await(5, java.util.concurrent.TimeUnit.SECONDS), "first subscriber got data");
        assertTrue(second.dataLatch.await(5, java.util.concurrent.TimeUnit.SECONDS), "second subscriber got data");
        assertArrayEquals(payload, first.bytes());
        assertArrayEquals(payload, second.bytes());

        manager.releasePort(shared);
    }

    @Test
    void fatalReadErrorNotifiesSubscribersAndEvicts() throws Exception {
        MockPortFactory factory = new MockPortFactory();
        SharedSerialPortManager manager = new SharedSerialPortManager(factory);
        SharedSerialPortManager.SharedPort shared = manager.acquirePort("COMF", config(9600));
        SerialPort port = factory.created.get("COMF");

        CollectingSubscriber subscriber = new CollectingSubscriber();
        shared.addSubscriber(subscriber);

        // Stub under the SharedPort monitor: readFromPort() synchronizes on
        // it, so this blocks the live reader while Mockito's (non-thread-safe)
        // stub bookkeeping is mutated. readBytes is stubbed first so the
        // reader can never see a non-zero bytesAvailable with an unstubbed read.
        synchronized (shared) {
            when(port.readBytes(any(byte[].class), anyInt())).thenReturn(-1); // fatal
            when(port.bytesAvailable()).thenReturn(1);
        }

        assertTrue(subscriber.failureLatch.await(5, java.util.concurrent.TimeUnit.SECONDS), "onFailure fired");
        assertNotNull(subscriber.failure);

        // Evicted: a fresh acquire opens a new physical port.
        manager.acquirePort("COMF", config(9600));
        assertEquals(1, factory.created.size(), "entry was evicted and re-created under the same name");
    }

    @Test
    void lastReleaseStopsReaderBeforeClosingPort() throws Exception {
        MockPortFactory factory = new MockPortFactory();
        SharedSerialPortManager manager = new SharedSerialPortManager(factory);
        SharedSerialPortManager.SharedPort shared = manager.acquirePort("COMS", config(9600));
        SerialPort port = factory.created.get("COMS");

        manager.releasePort(shared);

        org.mockito.InOrder inOrder = inOrder(port);
        inOrder.verify(port).removeDataListener();
        inOrder.verify(port).closePort();
    }

    @Test
    void sharedWritesArePaced() throws Exception {
        MockPortFactory factory = new MockPortFactory();
        SharedSerialPortManager manager = new SharedSerialPortManager(factory);
        SharedSerialPortManager.SerialPortConfig paced = new SharedSerialPortManager.SerialPortConfig(
            9600, 8, 1, SerialPort.NO_PARITY, SerialPort.FLOW_CONTROL_DISABLED,
            1000, 1000, false, false, 60);
        SharedSerialPortManager.SharedPort shared = manager.acquirePort("COMP", paced);

        shared.lockWrite();
        shared.unlockWrite();
        long start = System.currentTimeMillis();
        shared.lockWrite();
        shared.unlockWrite();
        assertTrue(System.currentTimeMillis() - start >= 50,
            "second write must wait out the inter-frame gap");

        manager.releasePort(shared);
    }
}
