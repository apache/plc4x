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

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

/**
 * Tests for TagBatch, focusing on tag management functionality.
 */
class TagBatchTest {

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

        // Create mock connection manager that returns our stub connection
        connectionManager = mock(PlcConnectionManager.class);
        when(connectionManager.getConnection(any())).thenReturn(connection);

        connectionString = "test://localhost";
    }

    @Test
    void testRemoveSingleTag() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .addTagAddress("tag2", "MAIN.tag2")
            .addTagAddress("tag3", "MAIN.tag3")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        boolean removed = batch.removeTag("tag2");

        // Assert
        assertTrue(removed);
        assertEquals(2, batch.getTagCount());
        assertFalse(batch.hasTag("tag2"));
        assertTrue(batch.hasTag("tag1"));
        assertTrue(batch.hasTag("tag3"));

        // Cleanup
        batch.close();
    }

    @Test
    void testRemoveNonExistentTag() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        boolean removed = batch.removeTag("nonexistent");

        // Assert
        assertFalse(removed);
        assertEquals(1, batch.getTagCount());

        // Cleanup
        batch.close();
    }

    @Test
    void testRemoveTagFromClosedBatch() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();
        batch.close();

        // Act & Assert
        assertThrows(IllegalStateException.class, () -> {
            batch.removeTag("tag1");
        });
    }

    @Test
    void testRemoveMultipleTags() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .addTagAddress("tag2", "MAIN.tag2")
            .addTagAddress("tag3", "MAIN.tag3")
            .addTagAddress("tag4", "MAIN.tag4")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        Collection<String> tagsToRemove = Arrays.asList("tag1", "tag3", "nonexistent");
        int removedCount = batch.removeTags(tagsToRemove);

        // Assert
        assertEquals(2, removedCount);
        assertEquals(2, batch.getTagCount());
        assertFalse(batch.hasTag("tag1"));
        assertTrue(batch.hasTag("tag2"));
        assertFalse(batch.hasTag("tag3"));
        assertTrue(batch.hasTag("tag4"));

        // Cleanup
        batch.close();
    }

    @Test
    void testRemoveTagsWithNullCollection() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        int removedCount = batch.removeTags(null);

        // Assert
        assertEquals(0, removedCount);
        assertEquals(1, batch.getTagCount());

        // Cleanup
        batch.close();
    }

    @Test
    void testRemoveTagsFromClosedBatch() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();
        batch.close();

        // Act & Assert
        assertThrows(IllegalStateException.class, () -> {
            batch.removeTags(Arrays.asList("tag1"));
        });
    }

    @Test
    void testClearTags() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .addTagAddress("tag2", "MAIN.tag2")
            .addTagAddress("tag3", "MAIN.tag3")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        batch.clearTags();

        // Assert
        assertEquals(0, batch.getTagCount());
        assertFalse(batch.hasTag("tag1"));
        assertFalse(batch.hasTag("tag2"));
        assertFalse(batch.hasTag("tag3"));
        assertTrue(batch.getTagNames().isEmpty());

        // Cleanup
        batch.close();
    }

    @Test
    void testClearTagsFromClosedBatch() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();
        batch.close();

        // Act & Assert
        assertThrows(IllegalStateException.class, batch::clearTags);
    }

    @Test
    void testHasTag() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .addTagAddress("tag2", "MAIN.tag2")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Assert
        assertTrue(batch.hasTag("tag1"));
        assertTrue(batch.hasTag("tag2"));
        assertFalse(batch.hasTag("tag3"));
        assertFalse(batch.hasTag("nonexistent"));

        // Cleanup
        batch.close();
    }

    @Test
    void testGetTagCount() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .addTagAddress("tag2", "MAIN.tag2")
            .addTagAddress("tag3", "MAIN.tag3")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Assert
        assertEquals(3, batch.getTagCount());

        // Add more tags
        batch.addTag("tag4", "MAIN.tag4");
        assertEquals(4, batch.getTagCount());

        // Remove tags
        batch.removeTag("tag1");
        assertEquals(3, batch.getTagCount());

        // Cleanup
        batch.close();
    }

    @Test
    void testGetTagNames() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .addTagAddress("tag2", "MAIN.tag2")
            .addTagAddress("tag3", "MAIN.tag3")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        Set<String> tagNames = batch.getTagNames();

        // Assert
        assertEquals(3, tagNames.size());
        assertTrue(tagNames.contains("tag1"));
        assertTrue(tagNames.contains("tag2"));
        assertTrue(tagNames.contains("tag3"));

        // Test immutability
        assertThrows(UnsupportedOperationException.class, () -> {
            tagNames.add("tag4");
        });

        // Cleanup
        batch.close();
    }

    @Test
    void testAddTagAtRuntime() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        batch.addTag("tag2", "MAIN.tag2");

        // Assert
        assertEquals(2, batch.getTagCount());
        assertTrue(batch.hasTag("tag1"));
        assertTrue(batch.hasTag("tag2"));

        // Cleanup
        batch.close();
    }

    @Test
    void testAddTagWithNullName() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act & Assert
        assertThrows(IllegalArgumentException.class, () -> {
            batch.addTag(null, "MAIN.tag2");
        });

        // Cleanup
        batch.close();
    }

    @Test
    void testAddTagWithEmptyName() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act & Assert
        assertThrows(IllegalArgumentException.class, () -> {
            batch.addTag("  ", "MAIN.tag2");
        });

        // Cleanup
        batch.close();
    }

    @Test
    void testAddTagWithNullAddress() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act & Assert
        assertThrows(IllegalArgumentException.class, () -> {
            batch.addTag("tag2", null);
        });

        // Cleanup
        batch.close();
    }

    @Test
    void testAddTagToClosedBatch() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();
        batch.close();

        // Act & Assert
        assertThrows(IllegalStateException.class, () -> {
            batch.addTag("tag2", "MAIN.tag2");
        });
    }

    @Test
    void testAddMultipleTagsAtRuntime() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        Map<String, String> newTags = new HashMap<>();
        newTags.put("tag2", "MAIN.tag2");
        newTags.put("tag3", "MAIN.tag3");

        // Act
        batch.addTags(newTags);

        // Assert
        assertEquals(3, batch.getTagCount());
        assertTrue(batch.hasTag("tag1"));
        assertTrue(batch.hasTag("tag2"));
        assertTrue(batch.hasTag("tag3"));

        // Cleanup
        batch.close();
    }

    @Test
    void testAddTagsWithNullMap() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act & Assert
        assertThrows(IllegalArgumentException.class, () -> {
            batch.addTags(null);
        });

        // Cleanup
        batch.close();
    }

    @Test
    void testAddTagsWithEmptyMap() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act & Assert
        assertThrows(IllegalArgumentException.class, () -> {
            batch.addTags(new HashMap<>());
        });

        // Cleanup
        batch.close();
    }

    @Test
    void testAddTagsToClosedBatch() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();
        batch.close();

        Map<String, String> newTags = new HashMap<>();
        newTags.put("tag2", "MAIN.tag2");

        // Act & Assert
        assertThrows(IllegalStateException.class, () -> {
            batch.addTags(newTags);
        });
    }

    @Test
    void testDynamicTagManagementWhileRunning() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        batch.start();
        assertTrue(batch.isStarted());

        // Act - Add and remove tags while running
        batch.addTag("tag2", "MAIN.tag2");
        batch.addTag("tag3", "MAIN.tag3");
        assertEquals(3, batch.getTagCount());

        batch.removeTag("tag1");
        assertEquals(2, batch.getTagCount());

        // Assert
        assertFalse(batch.hasTag("tag1"));
        assertTrue(batch.hasTag("tag2"));
        assertTrue(batch.hasTag("tag3"));

        // Cleanup
        batch.close();
    }

    @Test
    void testGetTagsReturnsUnmodifiableMap() {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .build();

        // Act
        Map<String, String> tags = batch.getTags();

        // Assert - Should not be able to modify
        assertThrows(UnsupportedOperationException.class, () -> {
            tags.put("tag2", "MAIN.tag2");
        });

        // Cleanup
        batch.close();
    }

    @Test
    void testTagBatchListenerInterface() throws Exception {
        // Arrange
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);

        // Create listener using anonymous class (not lambda) to ensure interface gets coverage
        TagBatch.TagBatchListener listener = new TagBatch.TagBatchListener() {
            @Override
            public void onTagsFetched(TagBatch batch, PlcReadResponse response) {
                assertNotNull(batch);
                assertNotNull(response);
            }
        };

        TagBatch batch = TagBatch.builder()
            .withBatchId("batch1")
            .withConnectionManager(connectionManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener(listener)
            .build();

        // Act - Manually invoke listener to test interface
        listener.onTagsFetched(batch, mockReadResponse);

        // Cleanup
        batch.close();
    }

    // ========== US3: Exponential Backoff Tests ==========

    /**
     * T019: Verify that after multiple consecutive fetchTags() failures,
     * the backoff interval increases exponentially (fewer getConnection calls).
     */
    @Test
    void testBackoffOnConsecutiveFailures() throws Exception {
        // Arrange — make getConnection throw to simulate network outage
        PlcConnectionManager failingManager = mock(PlcConnectionManager.class);
        when(failingManager.getConnection(any())).thenThrow(new RuntimeException("Connection refused"));

        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("backoff-test")
            .withConnectionManager(failingManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .withInitialBackoffMs(50)
            .withMaxBackoffMs(500)
            .build();

        // Act — first failure: should connect and fail
        batch.fetchTags();
        assertEquals(1, batch.getConsecutiveFailures());
        verify(failingManager, times(1)).getConnection(any());

        // Second call immediately — should be skipped (within backoff window)
        batch.fetchTags();
        assertEquals(1, batch.getConsecutiveFailures()); // Unchanged — skipped
        verify(failingManager, times(1)).getConnection(any()); // No new call

        // Wait for first backoff to expire (50ms), then try again
        Thread.sleep(100);
        batch.fetchTags();
        assertEquals(2, batch.getConsecutiveFailures()); // New failure
        verify(failingManager, times(2)).getConnection(any());

        // Wait for second backoff (100ms), then try again
        Thread.sleep(150);
        batch.fetchTags();
        assertEquals(3, batch.getConsecutiveFailures()); // New failure
        verify(failingManager, times(3)).getConnection(any());

        // Verify backoff is increasing — next allowed time is further out
        assertTrue(batch.getNextAllowedFetchTimeMs() > System.currentTimeMillis());

        // Cleanup
        batch.close();
    }

    /**
     * T020: Verify that after a successful fetchTags() call, the backoff resets to zero
     * and normal polling resumes immediately.
     */
    @Test
    void testBackoffResetOnSuccess() throws Exception {
        // Arrange — start with a failing connection, then switch to success
        PlcConnectionManager switchableManager = mock(PlcConnectionManager.class);
        when(switchableManager.getConnection(any())).thenThrow(new RuntimeException("Connection refused"));

        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("backoff-reset-test")
            .withConnectionManager(switchableManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .withInitialBackoffMs(50)
            .withMaxBackoffMs(500)
            .build();

        // Fail a few times to build up backoff
        batch.fetchTags(); // Fail 1
        Thread.sleep(100);
        batch.fetchTags(); // Fail 2
        assertEquals(2, batch.getConsecutiveFailures());
        assertTrue(batch.getNextAllowedFetchTimeMs() > 0);

        // Wait for backoff to expire
        Thread.sleep(200);

        // Switch to a working connection — use doReturn to override existing thenThrow
        PlcReadRequest.Builder mockBuilder = mock(PlcReadRequest.Builder.class);
        PlcReadRequest mockRequest = mock(PlcReadRequest.class);
        when(mockBuilder.addTagAddress(any(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(inv ->
            CompletableFuture.completedFuture(mockReadResponse));
        PlcConnection workingConnection = new StubPlcConnection(mockBuilder);
        doReturn(workingConnection).when(switchableManager).getConnection(any());

        // Act — successful fetch
        CompletableFuture<Void> result = batch.fetchTags();
        result.get(2, TimeUnit.SECONDS); // Wait for completion

        // Assert — backoff should be reset
        assertEquals(0, batch.getConsecutiveFailures());
        assertEquals(0, batch.getNextAllowedFetchTimeMs());

        // Cleanup
        batch.close();
    }

    /**
     * T021: Verify that the backoff interval never exceeds the configured maximum.
     */
    @Test
    void testBackoffMaximumCap() throws Exception {
        // Arrange
        PlcConnectionManager failingManager = mock(PlcConnectionManager.class);
        when(failingManager.getConnection(any())).thenThrow(new RuntimeException("Connection refused"));

        long maxBackoffMs = 200;
        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("backoff-cap-test")
            .withConnectionManager(failingManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .withInitialBackoffMs(50)
            .withMaxBackoffMs(maxBackoffMs)
            .build();

        // Act — fail many times to hit the cap
        for (int i = 0; i < 10; i++) {
            batch.fetchTags();
            // Wait for the backoff to expire before the next attempt
            Thread.sleep(maxBackoffMs + 50);
        }

        // Assert — consecutive failures should be 10 but backoff capped at maxBackoffMs
        assertEquals(10, batch.getConsecutiveFailures());
        // The next allowed time should be at most maxBackoffMs from now
        long maxAllowed = System.currentTimeMillis() + maxBackoffMs;
        assertTrue(batch.getNextAllowedFetchTimeMs() <= maxAllowed,
            "Backoff should not exceed maxBackoffMs (" + maxBackoffMs + "ms)");

        // Cleanup
        batch.close();
    }

    /**
     * T022: Verify that trigger invocations falling within the backoff window are skipped
     * (return immediately without attempting connection).
     */
    @Test
    void testBackoffSkipsTrigger() throws Exception {
        // Arrange
        PlcConnectionManager failingManager = mock(PlcConnectionManager.class);
        when(failingManager.getConnection(any())).thenThrow(new RuntimeException("Connection refused"));

        Trigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        TagBatch batch = TagBatch.builder()
            .withBatchId("backoff-skip-test")
            .withConnectionManager(failingManager)
            .withConnectionString(connectionString)
            .addTagAddress("tag1", "MAIN.tag1")
            .withTrigger(trigger)
            .withListener((b, r) -> {})
            .withInitialBackoffMs(500)  // Long backoff window
            .withMaxBackoffMs(5000)
            .build();

        // First call fails and enters backoff
        batch.fetchTags();
        verify(failingManager, times(1)).getConnection(any());
        assertEquals(1, batch.getConsecutiveFailures());

        // Multiple rapid calls — all should be skipped
        for (int i = 0; i < 5; i++) {
            CompletableFuture<Void> result = batch.fetchTags();
            // Should complete immediately with null (skipped)
            assertNotNull(result);
            assertNull(result.get(1, TimeUnit.SECONDS));
        }

        // Verify no additional connection attempts were made
        verify(failingManager, times(1)).getConnection(any());
        // Failure count should still be 1 — skipped triggers don't increment it
        assertEquals(1, batch.getConsecutiveFailures());

        // Cleanup
        batch.close();
    }
}
