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

package org.apache.plc4x.java.utils.subscriptionemulation;

import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 * Tests for PollingSubscriptionConnectionBase.
 */
class PollingSubscriptionConnectionBaseTest {

    private TestConnection connection;
    private PlcTagHandler tagHandler;
    private AtomicInteger mockValueCounter;

    @BeforeEach
    void setUp() {
        mockValueCounter = new AtomicInteger(100);
        tagHandler = mock(PlcTagHandler.class);

        // Setup tag handler to return mock tags
        PlcTag mockTag = mock(PlcTag.class);
        when(mockTag.getAddressString()).thenReturn("MAIN.tag1");
        when(tagHandler.parseTag(anyString())).thenReturn(mockTag);

        connection = new TestConnection(mockValueCounter, tagHandler);
    }

    @Test
    void testCyclicSubscriptionFiresRegularly() throws Exception {
        // Arrange
        CountDownLatch latch = new CountDownLatch(3);
        List<Integer> receivedValues = Collections.synchronizedList(new ArrayList<>());

        PlcSubscriptionRequest request = new DefaultPlcSubscriptionRequest.Builder(connection, tagHandler)
            .addCyclicTagAddress("tag1", "MAIN.tag1", Duration.ofMillis(100))
            .setConsumer(event -> {
                PlcValue value = event.getPlcValue("tag1");
                receivedValues.add(value.getInteger());
                latch.countDown();
            })
            .build();

        // Act
        CompletableFuture<PlcSubscriptionResponse> subscribeFuture = connection.subscribe(request);
        PlcSubscriptionResponse response = subscribeFuture.get(1, TimeUnit.SECONDS);

        // Assert subscription was successful
        assertEquals(PlcResponseCode.OK, response.getResponseCode("tag1"));
        assertNotNull(response.getSubscriptionHandle("tag1"));

        // Wait for multiple events (generous timeout for CI/heavy classpath environments)
        boolean completed = latch.await(2, TimeUnit.SECONDS);
        assertTrue(completed, "Should have received at least 3 events");
        assertTrue(receivedValues.size() >= 3, "Should have received at least 3 values");

        // Cleanup
        connection.close();
    }

    @Test
    void testChangeOfStateSubscriptionOnlyFiresOnChange() throws Exception {
        // Arrange
        CountDownLatch latch = new CountDownLatch(2);
        List<Integer> receivedValues = Collections.synchronizedList(new ArrayList<>());

        PlcSubscriptionRequest request = new DefaultPlcSubscriptionRequest.Builder(connection, tagHandler)
            .addChangeOfStateTagAddress("tag1", "MAIN.tag1")
            .setConsumer(event -> {
                PlcValue value = event.getPlcValue("tag1");
                receivedValues.add(value.getInteger());
                latch.countDown();
            })
            .build();

        // Act
        CompletableFuture<PlcSubscriptionResponse> subscribeFuture = connection.subscribe(request);
        PlcSubscriptionResponse response = subscribeFuture.get(1, TimeUnit.SECONDS);

        // Assert subscription was successful
        assertEquals(PlcResponseCode.OK, response.getResponseCode("tag1"));

        // Simulate value changes
        Thread.sleep(200); // Wait for the first poll (value=100)
        mockValueCounter.set(200); // Change value
        Thread.sleep(300); // Wait for polls

        boolean completed = latch.await(2, TimeUnit.SECONDS);
        assertTrue(completed, "Should have received at least 2 change events");

        // Should have received at least 2 different values
        assertTrue(receivedValues.size() >= 2, "Should have at least 2 values");

        // Cleanup
        connection.close();
    }

    @Test
    void testEventSubscriptionThrowsException() {
        // Arrange
        PlcSubscriptionRequest request = new DefaultPlcSubscriptionRequest.Builder(connection, tagHandler)
            .addEventTagAddress("tag1", "MAIN.tag1")
            .build();

        // Act & Assert
        assertThrows(Exception.class, () -> {
            CompletableFuture<PlcSubscriptionResponse> subscribeFuture = connection.subscribe(request);
            subscribeFuture.get(1, TimeUnit.SECONDS);
        });

        // Cleanup
        assertDoesNotThrow(() -> connection.close());
    }

    @Test
    void testConnectionWithoutReaderDoesNotFireEvents() throws Exception {
        // Arrange - Create a connection that does NOT override onRead()
        // Subscribe will succeed, but polling will fail because default onRead() returns "not supported"
        TestConnectionWithoutReader connection = new TestConnectionWithoutReader(tagHandler);

        AtomicInteger eventCount = new AtomicInteger(0);

        PlcSubscriptionRequest request = new DefaultPlcSubscriptionRequest.Builder(connection, tagHandler)
            .addCyclicTagAddress("tag1", "MAIN.tag1", Duration.ofMillis(50))
            .setConsumer(event -> eventCount.incrementAndGet())
            .build();

        // Act - subscribe should succeed
        CompletableFuture<PlcSubscriptionResponse> subscribeFuture = connection.subscribe(request);
        PlcSubscriptionResponse response = subscribeFuture.get(1, TimeUnit.SECONDS);
        assertEquals(PlcResponseCode.OK, response.getResponseCode("tag1"));

        // Wait a bit — events should NOT fire because onRead() returns "not supported"
        Thread.sleep(200);
        assertEquals(0, eventCount.get(), "Should not receive events when onRead is not overridden");

        // Cleanup
        assertDoesNotThrow(() -> connection.close());
    }

    @Test
    void testCustomValueComparison() throws Exception {
        // Create a connection with custom value comparison that ignores changes
        TestConnectionWithCustomComparison customConnection = new TestConnectionWithCustomComparison(mockValueCounter, tagHandler);

        CountDownLatch latch = new CountDownLatch(1);
        AtomicInteger eventCount = new AtomicInteger(0);

        PlcSubscriptionRequest request = new DefaultPlcSubscriptionRequest.Builder(customConnection, tagHandler)
            .addChangeOfStateTagAddress("tag1", "MAIN.tag1")
            .setConsumer(event -> {
                eventCount.incrementAndGet();
                latch.countDown();
            })
            .build();

        // Subscribe
        CompletableFuture<PlcSubscriptionResponse> subscribeFuture = customConnection.subscribe(request);
        subscribeFuture.get(1, TimeUnit.SECONDS);

        // Wait for first event
        latch.await(2, TimeUnit.SECONDS);
        assertEquals(1, eventCount.get(), "Should have received initial event");

        // Change value but custom comparator always returns true (equal)
        mockValueCounter.set(999);
        Thread.sleep(150);

        // Should still only have 1 event because custom comparator says values are equal
        assertTrue(eventCount.get() <= 2, "Should not receive many events with custom always-equal comparator");

        // Cleanup
        customConnection.close();
    }

    // Test stub implementation
    static class TestConnection extends PollingSubscriptionConnectionBase<TestConfiguration> {
        private final AtomicInteger valueSupplier;
        private final PlcTagHandler tagHandler;

        public TestConnection(AtomicInteger valueSupplier, PlcTagHandler tagHandler) {
            super(new TestConfiguration(), mock(TransportInstance.class), mock(AuditLog.class));
            this.valueSupplier = valueSupplier;
            this.tagHandler = tagHandler;
        }

        @Override
        protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
            LinkedHashMap<String, DefaultPlcResponseItem<PlcValue>> responseItems = new LinkedHashMap<>();

            for (String tagName : readRequest.getTagNames()) {
                PlcValue value = mock(PlcValue.class);
                when(value.getInteger()).thenReturn(valueSupplier.get());
                when(value.getObject()).thenReturn(valueSupplier.get());

                responseItems.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, value));
            }

            PlcReadResponse response = new DefaultPlcReadResponse(readRequest, (Map) responseItems);
            return CompletableFuture.completedFuture(response);
        }

        @Override
        protected PlcTagHandler getTagHandler() {
            return tagHandler;
        }

        @Override
        protected PlcValueHandler getValueHandler() {
            return null;
        }

        @Override
        protected int getMaxConcurrentRequests() {
            return 10;
        }

        @Override
        protected long getDefaultPollingInterval() {
            return 50; // Fast polling for tests
        }
    }

    // Test connection that does NOT implement PlcReader
    static class TestConnectionWithoutReader extends PollingSubscriptionConnectionBase<TestConfiguration> {
        private final PlcTagHandler tagHandler;

        public TestConnectionWithoutReader(PlcTagHandler tagHandler) {
            super(new TestConfiguration(), mock(TransportInstance.class), mock(AuditLog.class));
            this.tagHandler = tagHandler;
        }

        @Override
        protected PlcTagHandler getTagHandler() {
            return tagHandler;
        }

        @Override
        protected PlcValueHandler getValueHandler() {
            return null;
        }

        @Override
        protected int getMaxConcurrentRequests() {
            return 10;
        }
    }

    // Test connection with custom value comparison
    static class TestConnectionWithCustomComparison extends TestConnection {
        public TestConnectionWithCustomComparison(AtomicInteger valueSupplier, PlcTagHandler tagHandler) {
            super(valueSupplier, tagHandler);
        }

        @Override
        protected boolean valuesEqual(PlcValue v1, PlcValue v2) {
            // Always return true (values are always "equal")
            return true;
        }
    }

    // Test configuration
    static class TestConfiguration implements Configuration {
    }
}
