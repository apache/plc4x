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

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.util.Collection;
import java.util.LinkedHashSet;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 * Tests for SubscriptionRecord.
 */
class SubscriptionRecordTest {

    @Mock
    private PlcSubscriptionRequest mockRequest;

    @Mock
    private PlcSubscriptionResponse mockResponse;

    @Mock
    private PlcSubscriptionHandle mockHandle1;

    @Mock
    private PlcSubscriptionHandle mockHandle2;

    @Mock
    private Consumer<PlcSubscriptionEvent> mockConsumer;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
    }

    @Test
    void testConstructor_WithResponse() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        // Act
        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Assert
        assertNotNull(record);
        assertEquals(mockRequest, record.getOriginalRequest());
        assertEquals(mockConsumer, record.getConsumer());
        assertEquals(mockResponse, record.getCurrentResponse());
    }

    @Test
    void testConstructor_NullResponse() {
        // Act
        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, null);

        // Assert
        assertNotNull(record);
        assertEquals(mockRequest, record.getOriginalRequest());
        assertEquals(mockConsumer, record.getConsumer());
        assertNull(record.getCurrentResponse());
    }

    @Test
    void testGetCurrentHandles() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Act
        Collection<PlcSubscriptionHandle> handles = record.getCurrentHandles();

        // Assert
        assertNotNull(handles);
        assertEquals(1, handles.size());
        assertTrue(handles.contains(mockHandle1));
    }

    @Test
    void testContainsHandle_True() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Act & Assert
        assertTrue(record.containsHandle(mockHandle1));
    }

    @Test
    void testContainsHandle_False() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Act & Assert
        assertFalse(record.containsHandle(mockHandle2));
    }

    @Test
    void testContainsHandle_Null() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Act & Assert
        assertFalse(record.containsHandle(null));
    }

    @Test
    void testGetCurrentHandle_ReturnsOriginal() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Act
        PlcSubscriptionHandle result = record.getCurrentHandle(mockHandle1);

        // Assert
        assertEquals(mockHandle1, result);
    }

    @Test
    void testGetCurrentHandle_Null() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Act
        PlcSubscriptionHandle result = record.getCurrentHandle(null);

        // Assert
        assertNull(result);
    }

    @Test
    void testGetCurrentHandle_UnknownHandle() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Act - handle2 was never recorded
        PlcSubscriptionHandle result = record.getCurrentHandle(mockHandle2);

        // Assert - returns original since no mapping exists
        assertEquals(mockHandle2, result);
    }

    @Test
    void testUpdateAfterReconnection_NullResponse() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Act - should not throw
        record.updateAfterReconnection(null);

        // Assert - response should remain unchanged
        assertEquals(mockResponse, record.getCurrentResponse());
    }

    @Test
    void testUpdateAfterReconnection() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);

        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Create new response with new handle
        PlcSubscriptionResponse newResponse = mock(PlcSubscriptionResponse.class);
        LinkedHashSet<String> newTagNames = new LinkedHashSet<>();
        newTagNames.add("tag1");
        when(newResponse.getTagNames()).thenReturn(newTagNames);
        when(newResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle2);

        // Act
        record.updateAfterReconnection(newResponse);

        // Assert
        assertEquals(newResponse, record.getCurrentResponse());
    }

    @Test
    void testMultipleTags() {
        // Arrange
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        tagNames.add("tag2");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);
        when(mockResponse.getSubscriptionHandle("tag2")).thenReturn(mockHandle2);

        // Act
        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Assert
        assertTrue(record.containsHandle(mockHandle1));
        assertTrue(record.containsHandle(mockHandle2));
        assertEquals(2, record.getCurrentHandles().size());
    }

    @Test
    void testNullHandleInResponse() {
        // Arrange - response returns null for one handle
        LinkedHashSet<String> tagNames = new LinkedHashSet<>();
        tagNames.add("tag1");
        tagNames.add("tag2");
        when(mockResponse.getTagNames()).thenReturn(tagNames);
        when(mockResponse.getSubscriptionHandle("tag1")).thenReturn(mockHandle1);
        when(mockResponse.getSubscriptionHandle("tag2")).thenReturn(null);

        // Act
        SubscriptionRecord record = new SubscriptionRecord(mockRequest, mockConsumer, mockResponse);

        // Assert
        assertTrue(record.containsHandle(mockHandle1));
        assertEquals(1, record.getCurrentHandles().size());
    }
}
