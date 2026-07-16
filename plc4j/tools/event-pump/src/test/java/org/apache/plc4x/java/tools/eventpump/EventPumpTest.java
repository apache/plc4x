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

package org.apache.plc4x.java.tools.eventpump;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.tools.eventpump.triggers.TimerTrigger;
import org.apache.plc4x.java.tools.eventpump.triggers.Trigger;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Tests for EventPump.
 */
class EventPumpTest {

    private PlcConnectionManager connectionManager;
    private String connectionString;
    private PlcReadResponse mockReadResponse;

    @BeforeEach
    void setUp() throws Exception {
        // Create mocks for the request builders
        PlcReadRequest.Builder mockReadBuilder = mock(PlcReadRequest.Builder.class);
        PlcReadRequest mockReadRequest = mock(PlcReadRequest.class);
        mockReadResponse = mock(PlcReadResponse.class);

        // Setup mock behavior
        when(mockReadBuilder.addTagAddress(any(), any())).thenReturn(mockReadBuilder);
        when(mockReadBuilder.build()).thenReturn(mockReadRequest);
        when(mockReadRequest.execute()).thenAnswer(inv ->
            CompletableFuture.completedFuture(mockReadResponse));

        // Use stub instead of mocking PlcConnection
        PlcConnection connection = new StubPlcConnection(mockReadBuilder);

        // Create a mock connection manager that returns our stub connection
        connectionManager = mock(PlcConnectionManager.class);
        when(connectionManager.getConnection(any())).thenReturn(connection);

        connectionString = "test://localhost";
    }

    @Test
    void testAddAndRemoveBatch() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        pump.addBatch(batch);

        // Assert
        assertEquals(1, pump.getBatchCount());
        assertNotNull(pump.getBatch("batch1"));

        // Remove
        TagBatch removed = pump.removeBatch("batch1");
        assertEquals(batch, removed);
        assertEquals(0, pump.getBatchCount());
        assertNull(pump.getBatch("batch1"));

        // Cleanup
        pump.close();
    }

    @Test
    void testAddDuplicateBatch() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger1 = new TimerTrigger(1, TimeUnit.SECONDS);
        Trigger trigger2 = new TimerTrigger(1, TimeUnit.SECONDS);

        TagBatch batch1 = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger1)
            .withListener((b, r) -> {})
            .build();

        TagBatch batch2 = TagBatch.builder()
            .withBatchId("batch1") // Same ID
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger2)
            .withListener((b, r) -> {})
            .build();

        pump.addBatch(batch1);

        // Act & Assert - Adding duplicate should fail
        assertThrows(IllegalArgumentException.class, () -> {
            pump.addBatch(batch2);
        });

        // Cleanup
        trigger2.close();
        pump.close();
    }

    @Test
    void testStartAndStopBatch() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        pump.addBatch(batch);

        // Act
        pump.startBatch("batch1");

        // Assert
        assertTrue(batch.isStarted());
        assertEquals(1, pump.getStartedBatchCount());

        // Stop
        pump.stopBatch("batch1");
        assertFalse(batch.isStarted());
        assertEquals(0, pump.getStartedBatchCount());

        // Cleanup
        pump.close();
    }

    @Test
    void testStartAll() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger1 = new TimerTrigger(1, TimeUnit.SECONDS);
        Trigger trigger2 = new TimerTrigger(1, TimeUnit.SECONDS);

        TagBatch batch1 = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger1)
            .withListener((b, r) -> {})
            .build();

        TagBatch batch2 = TagBatch.builder()
            .withBatchId("batch2")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger2)
            .withListener((b, r) -> {})
            .build();

        pump.addBatch(batch1);
        pump.addBatch(batch2);

        // Act
        pump.startAll();

        // Assert
        assertTrue(batch1.isStarted());
        assertTrue(batch2.isStarted());
        assertEquals(2, pump.getStartedBatchCount());

        // Cleanup
        pump.close();
    }

    @Test
    void testStopAll() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger1 = new TimerTrigger(1, TimeUnit.SECONDS);
        Trigger trigger2 = new TimerTrigger(1, TimeUnit.SECONDS);

        TagBatch batch1 = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger1)
            .withListener((b, r) -> {})
            .build();

        TagBatch batch2 = TagBatch.builder()
            .withBatchId("batch2")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger2)
            .withListener((b, r) -> {})
            .build();

        pump.addBatch(batch1);
        pump.addBatch(batch2);
        pump.startAll();

        // Act
        pump.stopAll();

        // Assert
        assertFalse(batch1.isStarted());
        assertFalse(batch2.isStarted());
        assertEquals(0, pump.getStartedBatchCount());

        // Cleanup
        pump.close();
    }

    @Test
    void testGetAllBatches() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger1 = new TimerTrigger(1, TimeUnit.SECONDS);
        Trigger trigger2 = new TimerTrigger(1, TimeUnit.SECONDS);

        TagBatch batch1 = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger1)
            .withListener((b, r) -> {})
            .build();

        TagBatch batch2 = TagBatch.builder()
            .withBatchId("batch2")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger2)
            .withListener((b, r) -> {})
            .build();

        pump.addBatch(batch1);
        pump.addBatch(batch2);

        // Act
        Map<String, TagBatch> allBatches = pump.getAllBatches();

        // Assert
        assertEquals(2, allBatches.size());
        assertTrue(allBatches.containsKey("batch1"));
        assertTrue(allBatches.containsKey("batch2"));

        // Cleanup
        pump.close();
    }

    @Test
    void testGetBatchIds() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger1 = new TimerTrigger(1, TimeUnit.SECONDS);
        Trigger trigger2 = new TimerTrigger(1, TimeUnit.SECONDS);
        Trigger trigger3 = new TimerTrigger(1, TimeUnit.SECONDS);

        TagBatch batch1 = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger1)
            .withListener((b, r) -> {})
            .build();

        TagBatch batch2 = TagBatch.builder()
            .withBatchId("batch2")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger2)
            .withListener((b, r) -> {})
            .build();

        TagBatch batch3 = TagBatch.builder()
            .withBatchId("batch3")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger3)
            .withListener((b, r) -> {})
            .build();

        pump.addBatch(batch1);
        pump.addBatch(batch2);
        pump.addBatch(batch3);

        // Act
        java.util.Set<String> batchIds = pump.getBatchIds();

        // Assert
        assertEquals(3, batchIds.size());
        assertTrue(batchIds.contains("batch1"));
        assertTrue(batchIds.contains("batch2"));
        assertTrue(batchIds.contains("batch3"));

        // Cleanup
        pump.close();
    }

    @Test
    void testClose() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        pump.addBatch(batch);
        pump.startAll();

        // Act
        pump.close();

        // Assert
        assertTrue(pump.isClosed());
        assertEquals(0, pump.getBatchCount());
        assertFalse(batch.isStarted());
    }

    @Test
    void testOperationsAfterClose() {
        // Arrange
        EventPump pump = new EventPump();
        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "MAIN.tag1");

        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddresses(tags)
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        pump.close();

        // Act & Assert - Operations after close should fail
        assertThrows(IllegalStateException.class, () -> pump.addBatch(batch));
        assertThrows(IllegalStateException.class, pump::startAll);

        // Cleanup
        trigger.close();
    }

    @Test
    void testMultipleClose() {
        // Arrange
        EventPump pump = new EventPump();

        // Act - Multiple closes should be safe
        pump.close();
        pump.close();
        pump.close();

        // Assert - Should not throw exception
        assertTrue(pump.isClosed());
    }

    @Test
    void testStartNonExistentBatch() {
        // Arrange
        EventPump pump = new EventPump();

        // Act & Assert - Starting non-existent batch should fail
        assertThrows(IllegalArgumentException.class, () -> {
            pump.startBatch("nonexistent");
        });

        // Cleanup
        pump.close();
    }

    @Test
    void testStopNonExistentBatch() {
        // Arrange
        EventPump pump = new EventPump();

        // Act & Assert - Stopping non-existent batch should fail
        assertThrows(IllegalArgumentException.class, () -> {
            pump.stopBatch("nonexistent");
        });

        // Cleanup
        pump.close();
    }
}
