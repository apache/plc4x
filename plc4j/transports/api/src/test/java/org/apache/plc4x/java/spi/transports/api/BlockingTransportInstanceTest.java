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

import java.time.Duration;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

class BlockingTransportInstanceTest {

    @Test
    void testWaitForBytes_success() throws Exception {
        // Create a test implementation that simulates bytes becoming available
        BlockingTransportInstance<TransportConfiguration> instance = new TestBlockingTransportInstance(10);

        // Should complete immediately since 10 bytes are available
        assertDoesNotThrow(() -> instance.waitForBytes(5, Duration.ofSeconds(1)));
    }

    @Test
    void testWaitForBytes_timeout() {
        // Create a test implementation with only 5 bytes available
        BlockingTransportInstance<TransportConfiguration> instance = new TestBlockingTransportInstance(5);

        // Should timeout waiting for 10 bytes
        TimeoutException exception = assertThrows(TimeoutException.class, () -> {
            instance.waitForBytes(10, Duration.ofMillis(100));
        });

        assertTrue(exception.getMessage().contains("Timeout waiting for 10 bytes"));
    }

    @Test
    void testWaitForBytes_exactMatch() throws Exception {
        // Create a test implementation with exactly the number of bytes needed
        BlockingTransportInstance<TransportConfiguration> instance = new TestBlockingTransportInstance(10);

        // Should complete immediately
        assertDoesNotThrow(() -> instance.waitForBytes(10, Duration.ofSeconds(1)));
    }

    @Test
    void testWaitForBytes_zeroBytes() throws Exception {
        // Create a test implementation with no bytes
        BlockingTransportInstance<TransportConfiguration> instance = new TestBlockingTransportInstance(0);

        // Waiting for 0 bytes should succeed immediately
        assertDoesNotThrow(() -> instance.waitForBytes(0, Duration.ofMillis(100)));
    }

    @Test
    void testWaitForBytes_bytesIncreasingGradually() throws Exception {
        // Test implementation that simulates bytes becoming available gradually
        GradualBytesTransportInstance instance = new GradualBytesTransportInstance();

        // Should eventually get 10 bytes
        assertDoesNotThrow(() -> instance.waitForBytes(10, Duration.ofSeconds(2)));
    }

    /**
     * Test implementation that always reports a fixed number of available bytes
     */
    private static class TestBlockingTransportInstance implements BlockingTransportInstance<TransportConfiguration> {
        private final int availableBytes;

        TestBlockingTransportInstance(int availableBytes) {
            this.availableBytes = availableBytes;
        }

        @Override
        public TransportConfiguration getConfiguration() {
            return null;
        }

        @Override
        public void waitForData(Duration timeout) throws TransportException, TimeoutException, InterruptedException {
            // Simulate a short wait
            Thread.sleep(10);
        }

        @Override
        public int getNumBytesAvailable() {
            return availableBytes;
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
     * Test implementation that simulates bytes becoming available gradually
     */
    private static class GradualBytesTransportInstance implements BlockingTransportInstance<TransportConfiguration> {
        private final AtomicInteger availableBytes = new AtomicInteger(0);

        @Override
        public TransportConfiguration getConfiguration() {
            return null;
        }

        @Override
        public void waitForData(Duration timeout) throws TransportException, TimeoutException, InterruptedException {
            // Simulate data arriving gradually
            Thread.sleep(50);
            availableBytes.addAndGet(3);
        }

        @Override
        public int getNumBytesAvailable() {
            return availableBytes.get();
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
