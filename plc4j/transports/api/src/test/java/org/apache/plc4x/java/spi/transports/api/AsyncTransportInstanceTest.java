/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.spi.transports.api;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.junit.jupiter.api.Test;

import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.*;

class AsyncTransportInstanceTest {

    @Test
    void testDefaultDisconnectListenerMethods_noOp() {
        // Create a minimal implementation that uses the default disconnect listener methods
        AsyncTransportInstance<TransportConfiguration> instance = new MinimalAsyncTransportInstance();

        // Default methods should not throw - they are no-ops for backward compatibility
        assertDoesNotThrow(() -> instance.registerDisconnectListener(cause -> {}));
        assertDoesNotThrow(() -> instance.removeDisconnectListener());
    }

    @Test
    void testDefaultDisconnectListenerMethods_withNullListener() {
        AsyncTransportInstance<TransportConfiguration> instance = new MinimalAsyncTransportInstance();

        // Should handle null listener gracefully (no-op implementation)
        assertDoesNotThrow(() -> instance.registerDisconnectListener(null));
    }

    @Test
    void testCustomDisconnectListenerImplementation() {
        // Create an implementation that actually uses the disconnect listener
        CustomAsyncTransportInstance instance = new CustomAsyncTransportInstance();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);
        AtomicReference<Throwable> receivedCause = new AtomicReference<>();

        // Register a listener
        instance.registerDisconnectListener(cause -> {
            listenerCalled.set(true);
            receivedCause.set(cause);
        });

        // Simulate a disconnect
        RuntimeException testException = new RuntimeException("Test disconnect");
        instance.simulateDisconnect(testException);

        // Verify the listener was called with the correct cause
        assertTrue(listenerCalled.get());
        assertEquals(testException, receivedCause.get());
    }

    @Test
    void testCustomDisconnectListenerImplementation_gracefulDisconnect() {
        CustomAsyncTransportInstance instance = new CustomAsyncTransportInstance();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);
        // Use a sentinel exception to detect if the listener was called with null
        RuntimeException sentinel = new RuntimeException("SENTINEL");
        AtomicReference<Throwable> receivedCause = new AtomicReference<>(sentinel);

        instance.registerDisconnectListener(cause -> {
            listenerCalled.set(true);
            receivedCause.set(cause);
        });

        // Simulate a graceful disconnect (null cause)
        instance.simulateDisconnect(null);

        assertTrue(listenerCalled.get());
        assertNull(receivedCause.get());
    }

    @Test
    void testRemoveDisconnectListener() {
        CustomAsyncTransportInstance instance = new CustomAsyncTransportInstance();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);

        instance.registerDisconnectListener(cause -> listenerCalled.set(true));
        instance.removeDisconnectListener();

        // After removing, disconnect should not call the listener
        instance.simulateDisconnect(new RuntimeException("Test"));

        assertFalse(listenerCalled.get());
    }

    @Test
    void testDataListenerMethods() {
        CustomAsyncTransportInstance instance = new CustomAsyncTransportInstance();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);

        // Register a data listener
        instance.registerDataListener(() -> listenerCalled.set(true));

        // Simulate data available
        instance.simulateDataAvailable();

        assertTrue(listenerCalled.get());
    }

    @Test
    void testRemoveDataListener() {
        CustomAsyncTransportInstance instance = new CustomAsyncTransportInstance();

        AtomicBoolean listenerCalled = new AtomicBoolean(false);

        instance.registerDataListener(() -> listenerCalled.set(true));
        instance.removeDataListener();

        // After removing, data events should not call the listener
        instance.simulateDataAvailable();

        assertFalse(listenerCalled.get());
    }

    /**
     * Minimal implementation that uses the default disconnect listener methods (no-op)
     */
    private static class MinimalAsyncTransportInstance implements AsyncTransportInstance<TransportConfiguration> {
        private Runnable dataListener;

        @Override
        public TransportConfiguration getConfiguration() {
            return null;
        }

        @Override
        public void registerDataListener(Runnable listener) {
            this.dataListener = listener;
        }

        @Override
        public void removeDataListener() {
            this.dataListener = null;
        }

        // Note: registerDisconnectListener and removeDisconnectListener use default implementations

        @Override
        public int getNumBytesAvailable() {
            return 0;
        }

        @Override
        public boolean isOpen() {
            return true;
        }

        @Override
        public byte[] peekReadableBytes(int numBytes) throws TransportException {
            return new byte[0];
        }

        @Override
        public byte[] read(int numBytes) throws TransportException {
            return new byte[0];
        }

        @Override
        public void write(byte[] bytes) throws TransportException {
        }

        @Override
        public void close() throws TransportException {
        }
    }

    /**
     * Implementation that provides actual disconnect listener functionality for testing
     */
    private static class CustomAsyncTransportInstance implements AsyncTransportInstance<TransportConfiguration> {
        private Runnable dataListener;
        private Consumer<Throwable> disconnectListener;

        @Override
        public TransportConfiguration getConfiguration() {
            return null;
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

        public void simulateDataAvailable() {
            if (dataListener != null) {
                dataListener.run();
            }
        }

        public void simulateDisconnect(Throwable cause) {
            if (disconnectListener != null) {
                disconnectListener.accept(cause);
            }
        }

        @Override
        public int getNumBytesAvailable() {
            return 0;
        }

        @Override
        public boolean isOpen() {
            return true;
        }

        @Override
        public byte[] peekReadableBytes(int numBytes) throws TransportException {
            return new byte[0];
        }

        @Override
        public byte[] read(int numBytes) throws TransportException {
            return new byte[0];
        }

        @Override
        public void write(byte[] bytes) throws TransportException {
        }

        @Override
        public void close() throws TransportException {
        }
    }
}
