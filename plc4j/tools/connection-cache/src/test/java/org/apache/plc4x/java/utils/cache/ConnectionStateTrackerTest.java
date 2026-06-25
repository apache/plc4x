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
package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.api.EventPlcConnection;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.time.Duration;
import java.util.Collections;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 * Tests for ConnectionStateTracker.
 */
class ConnectionStateTrackerTest {

    private ConnectionStateTracker tracker;

    @Mock
    private PlcSubscriptionRequest mockRequest;

    @Mock
    private PlcSubscriptionResponse mockResponse;

    @Mock
    private PlcSubscriptionHandle mockHandle;

    @Mock
    private Consumer<PlcSubscriptionEvent> mockConsumer;

    @Mock
    private ConnectionStateListener mockListener;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
        tracker = new ConnectionStateTracker("test:tcp://localhost");
    }

    @Test
    void testInitialState() {
        assertEquals(0, tracker.getSubscriptionCount());
        assertEquals(0, tracker.getEventListenerCount());
        assertFalse(tracker.hasStateToRestore());
    }

    @Test
    void testRecordSubscription() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle);

        // Act
        String subscriptionId = tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);

        // Assert
        assertNotNull(subscriptionId);
        assertTrue(subscriptionId.startsWith("sub-"));
        assertEquals(1, tracker.getSubscriptionCount());
        assertTrue(tracker.hasStateToRestore());
    }

    @Test
    void testRemoveSubscription() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle);

        tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);
        assertEquals(1, tracker.getSubscriptionCount());

        // Act
        tracker.removeSubscription(List.of(mockHandle));

        // Assert
        assertEquals(0, tracker.getSubscriptionCount());
    }

    @Test
    void testRemoveSubscription_EmptyList() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle);

        tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);

        // Act - empty list should not remove anything
        tracker.removeSubscription(Collections.emptyList());

        // Assert
        assertEquals(1, tracker.getSubscriptionCount());
    }

    @Test
    void testRemoveSubscription_NullList() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle);

        tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);

        // Act - null list should not remove anything
        tracker.removeSubscription(null);

        // Assert
        assertEquals(1, tracker.getSubscriptionCount());
    }

    @Test
    void testGetCurrentHandle_NoMapping() {
        // Act - handle not tracked
        PlcSubscriptionHandle result = tracker.getCurrentHandle(mockHandle);

        // Assert - returns original if not found
        assertEquals(mockHandle, result);
    }

    @Test
    void testGetCurrentHandle_Null() {
        // Act
        PlcSubscriptionHandle result = tracker.getCurrentHandle(null);

        // Assert
        assertNull(result);
    }

    @Test
    void testAddEventListener() {
        // Act
        tracker.addEventListener(mockListener);

        // Assert
        assertEquals(1, tracker.getEventListenerCount());
        assertTrue(tracker.hasStateToRestore());
    }

    @Test
    void testAddEventListener_Duplicate() {
        // Act - add same listener twice
        tracker.addEventListener(mockListener);
        tracker.addEventListener(mockListener);

        // Assert - should only be added once
        assertEquals(1, tracker.getEventListenerCount());
    }

    @Test
    void testAddEventListener_Null() {
        // Act
        tracker.addEventListener(null);

        // Assert
        assertEquals(0, tracker.getEventListenerCount());
    }

    @Test
    void testRemoveEventListener() {
        // Arrange
        tracker.addEventListener(mockListener);
        assertEquals(1, tracker.getEventListenerCount());

        // Act
        tracker.removeEventListener(mockListener);

        // Assert
        assertEquals(0, tracker.getEventListenerCount());
    }

    @Test
    void testRemoveEventListener_Null() {
        // Arrange
        tracker.addEventListener(mockListener);

        // Act
        tracker.removeEventListener(null);

        // Assert - should not throw, should not remove anything
        assertEquals(1, tracker.getEventListenerCount());
    }

    @Test
    void testRestoreState_NullConnection() {
        // Arrange
        tracker.addEventListener(mockListener);

        // Act & Assert - should not throw
        assertDoesNotThrow(() -> tracker.restoreState(null));
    }

    @Test
    void testRestoreState_EventListeners() {
        // Arrange
        abstract class MockEventConnection implements PlcConnection, EventPlcConnection {}
        MockEventConnection mockConnection = mock(MockEventConnection.class);

        tracker.addEventListener(mockListener);

        // Act
        tracker.restoreState(mockConnection);

        // Assert
        verify(mockConnection).addEventListener(mockListener);
    }

    @Test
    void testRestoreState_EventListeners_ConnectionDoesNotSupportEvents() {
        // Arrange
        PlcConnection mockConnection = mock(PlcConnection.class);
        tracker.addEventListener(mockListener);

        // Act & Assert - should not throw, just log warning
        assertDoesNotThrow(() -> tracker.restoreState(mockConnection));
    }

    @Test
    void testRestoreState_EventListeners_ExceptionDuringRestore() {
        // Arrange
        abstract class MockEventConnection implements PlcConnection, EventPlcConnection {}
        MockEventConnection mockConnection = mock(MockEventConnection.class);
        doThrow(new RuntimeException("Test exception")).when(mockConnection).addEventListener(any());

        tracker.addEventListener(mockListener);

        // Act & Assert - should not throw
        assertDoesNotThrow(() -> tracker.restoreState(mockConnection));
    }

    @Test
    void testRestoreState_failedResubscribeKeepsRecordForRetry() {
        // A recorded subscription restored against a connection that cannot re-subscribe
        // (bare mock — subscriptionRequestBuilder() returns null) must be KEPT, not discarded,
        // so a later reconnection can retry it. (The previous implementation wrongly cleared it.)
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle);

        tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);
        assertEquals(1, tracker.getSubscriptionCount());

        tracker.restoreState(mock(PlcConnection.class));

        assertEquals(1, tracker.getSubscriptionCount());
    }

    @Test
    void testRestoreState_rebuildsSubscriptionAndRemapsHandle() {
        // Original subscription: one CHANGE_OF_STATE tag; the client's original handle is mockHandle.
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle);
        tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);

        // Describe the original tag so it can be replayed against the new connection.
        PlcSubscriptionTag subTag = mock(PlcSubscriptionTag.class);
        PlcTag innerTag = mock(PlcTag.class);
        when(mockRequest.getTagNames()).thenReturn(tagNames);
        when(mockRequest.getTag("tag1")).thenReturn(subTag);
        when(mockRequest.getTagConsumer("tag1")).thenReturn(null);
        when(subTag.getPlcSubscriptionType()).thenReturn(PlcSubscriptionType.CHANGE_OF_STATE);
        when(subTag.getTag()).thenReturn(innerTag);
        when(subTag.getDuration()).thenReturn(Optional.empty());

        // The new connection's builder chain yields a new response carrying a NEW handle.
        PlcSubscriptionHandle newHandle = mock(PlcSubscriptionHandle.class);
        PlcSubscriptionResponse newResponse = mock(PlcSubscriptionResponse.class);
        when(newResponse.getTagNames()).thenReturn(tagNames);
        when(newResponse.getSubscriptionHandle("tag1")).thenReturn(newHandle);

        PlcSubscriptionRequest rebuilt = mock(PlcSubscriptionRequest.class);
        doReturn(CompletableFuture.completedFuture(newResponse)).when(rebuilt).execute();

        PlcSubscriptionRequest.Builder builder = mock(PlcSubscriptionRequest.Builder.class);
        when(builder.setConsumer(any())).thenReturn(builder);
        when(builder.addChangeOfStateTag(eq("tag1"), eq(innerTag))).thenReturn(builder);
        when(builder.build()).thenReturn(rebuilt);

        PlcConnection newConnection = mock(PlcConnection.class);
        when(newConnection.subscriptionRequestBuilder()).thenReturn(builder);

        // Act
        tracker.restoreState(newConnection);

        // Assert: the subscription was actually rebuilt + re-executed on the new connection...
        verify(builder).addChangeOfStateTag("tag1", innerTag);
        verify(rebuilt).execute();
        // ...the record is retained...
        assertEquals(1, tracker.getSubscriptionCount());
        // ...and a client still holding the ORIGINAL handle now resolves to the NEW handle.
        assertSame(newHandle, tracker.getCurrentHandle(mockHandle));
    }

    @Test
    void testRestoreState_replaysAllTagTypes() {
        // One subscription with a CYCLIC (consumer + interval), an EVENT (consumer), and a
        // CHANGE_OF_STATE (interval, no consumer) tag — exercises every re-subscribe arm.
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("cyc");
        tagNames.add("evt");
        tagNames.add("cos");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle(anyString())).thenReturn(mockHandle);
        tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);

        PlcTag innerCyc = mock(PlcTag.class);
        PlcTag innerEvt = mock(PlcTag.class);
        PlcTag innerCos = mock(PlcTag.class);
        PlcSubscriptionTag cycTag = mock(PlcSubscriptionTag.class);
        PlcSubscriptionTag evtTag = mock(PlcSubscriptionTag.class);
        PlcSubscriptionTag cosTag = mock(PlcSubscriptionTag.class);
        when(cycTag.getPlcSubscriptionType()).thenReturn(PlcSubscriptionType.CYCLIC);
        when(cycTag.getTag()).thenReturn(innerCyc);
        when(cycTag.getDuration()).thenReturn(Optional.of(Duration.ofSeconds(2)));
        when(evtTag.getPlcSubscriptionType()).thenReturn(PlcSubscriptionType.EVENT);
        when(evtTag.getTag()).thenReturn(innerEvt);
        when(evtTag.getDuration()).thenReturn(Optional.empty());
        when(cosTag.getPlcSubscriptionType()).thenReturn(PlcSubscriptionType.CHANGE_OF_STATE);
        when(cosTag.getTag()).thenReturn(innerCos);
        when(cosTag.getDuration()).thenReturn(Optional.of(Duration.ofSeconds(1)));
        when(mockRequest.getTagNames()).thenReturn(tagNames);
        when(mockRequest.getTag("cyc")).thenReturn(cycTag);
        when(mockRequest.getTag("evt")).thenReturn(evtTag);
        when(mockRequest.getTag("cos")).thenReturn(cosTag);
        when(mockRequest.getTagConsumer("cyc")).thenReturn(mockConsumer);
        when(mockRequest.getTagConsumer("evt")).thenReturn(mockConsumer);
        when(mockRequest.getTagConsumer("cos")).thenReturn(null);

        PlcSubscriptionResponse newResponse = mock(PlcSubscriptionResponse.class);
        when(newResponse.getTagNames()).thenReturn(tagNames);
        when(newResponse.getSubscriptionHandle(anyString())).thenReturn(mock(PlcSubscriptionHandle.class));
        PlcSubscriptionRequest rebuilt = mock(PlcSubscriptionRequest.class);
        doReturn(CompletableFuture.completedFuture(newResponse)).when(rebuilt).execute();

        PlcSubscriptionRequest.Builder builder = mock(PlcSubscriptionRequest.Builder.class);
        when(builder.setConsumer(any())).thenReturn(builder);
        when(builder.addCyclicTag(eq("cyc"), eq(innerCyc), eq(Duration.ofSeconds(2)), any())).thenReturn(builder);
        when(builder.addEventTag(eq("evt"), eq(innerEvt), any())).thenReturn(builder);
        when(builder.addChangeOfStateTag(eq("cos"), eq(innerCos), eq(Duration.ofSeconds(1)))).thenReturn(builder);
        when(builder.build()).thenReturn(rebuilt);
        PlcConnection newConnection = mock(PlcConnection.class);
        when(newConnection.subscriptionRequestBuilder()).thenReturn(builder);

        tracker.restoreState(newConnection);

        verify(builder).addCyclicTag(eq("cyc"), eq(innerCyc), eq(Duration.ofSeconds(2)), any());
        verify(builder).addEventTag(eq("evt"), eq(innerEvt), any());
        verify(builder).addChangeOfStateTag("cos", innerCos, Duration.ofSeconds(1));
        verify(rebuilt).execute();
    }

    @Test
    void testClear() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle);

        tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);
        tracker.addEventListener(mockListener);

        // Act
        tracker.clear();

        // Assert
        assertEquals(0, tracker.getSubscriptionCount());
        assertEquals(0, tracker.getEventListenerCount());
        assertFalse(tracker.hasStateToRestore());
    }

    @Test
    void testHasStateToRestore_OnlySubscriptions() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle);

        tracker.recordSubscription(mockRequest, mockConsumer, mockResponse);

        // Assert
        assertTrue(tracker.hasStateToRestore());
    }

    @Test
    void testHasStateToRestore_OnlyListeners() {
        // Arrange
        tracker.addEventListener(mockListener);

        // Assert
        assertTrue(tracker.hasStateToRestore());
    }
}
