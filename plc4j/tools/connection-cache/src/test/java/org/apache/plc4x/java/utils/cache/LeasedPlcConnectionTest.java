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
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.listener.EventListener;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.util.Collection;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 * Tests for LeasedPlcConnection wrapper.
 */
class LeasedPlcConnectionTest {

    @Mock
    private PlcConnection mockDelegate;

    @Mock
    private PlcConnectionMetadata mockMetadata;

    @Mock
    private PlcTag mockTag;

    @Mock
    private PlcValue mockValue;

    @Mock
    private PlcReadRequest.Builder mockReadBuilder;

    @Mock
    private PlcWriteRequest.Builder mockWriteBuilder;

    @Mock
    private PlcBrowseRequest.Builder mockBrowseBuilder;

    @Mock
    private PlcSubscriptionRequest.Builder mockSubscriptionBuilder;

    @Mock
    private PlcUnsubscriptionRequest.Builder mockUnsubscriptionBuilder;

    @Mock
    private PlcPingResponse mockPingResponse;

    @SuppressWarnings("FieldCanBeLocal")
    private AutoCloseable mocks;

    @BeforeEach
    void setUp() {
        mocks = MockitoAnnotations.openMocks(this);
    }

    @Test
    void testClose_ReturnsToCache() throws Exception {
        // Arrange
        AtomicInteger returnCount = new AtomicInteger(0);
        Runnable returnCallback = returnCount::incrementAndGet;

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (!invalidated) returnCallback.run();
            }
        );

        // Act
        leased.close();

        // Assert
        assertEquals(1, returnCount.get());
        assertTrue(leased.isReturned());
    }

    @Test
    void testClose_MultipleTimesIsSafe() throws Exception {
        // Arrange
        AtomicInteger returnCount = new AtomicInteger(0);
        Runnable returnCallback = returnCount::incrementAndGet;

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (!invalidated) returnCallback.run();
            }
        );

        // Act
        leased.close();
        leased.close();
        leased.close();

        // Assert - Callback only called once
        assertEquals(1, returnCount.get());
        assertTrue(leased.isReturned());
    }

    @Test
    void testConnect_ThrowsException() {
        // Arrange
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act & Assert
        assertThrows(UnsupportedOperationException.class, leased::connect);
    }

    @Test
    void testIsConnected_Delegates() {
        // Arrange
        when(mockDelegate.isConnected()).thenReturn(true);
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        boolean connected = leased.isConnected();

        // Assert
        assertTrue(connected);
        verify(mockDelegate).isConnected();
    }

    @Test
    void testGetMetadata_Delegates() {
        // Arrange
        when(mockDelegate.getMetadata()).thenReturn(mockMetadata);
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcConnectionMetadata metadata = leased.getMetadata();

        // Assert
        assertEquals(mockMetadata, metadata);
        verify(mockDelegate).getMetadata();
    }

    @Test
    void testParseTagAddress_Delegates() {
        // Arrange
        when(mockDelegate.parseTagAddress("tag1")).thenReturn(Optional.of(mockTag));
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        Optional<PlcTag> tag = leased.parseTagAddress("tag1");

        // Assert
        assertTrue(tag.isPresent());
        assertEquals(mockTag, tag.get());
        verify(mockDelegate).parseTagAddress("tag1");
    }

    @Test
    void testParseTagValue_Delegates() {
        // Arrange
        when(mockDelegate.parseTagValue(mockTag, "value")).thenReturn(Optional.of(mockValue));
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        Optional<PlcValue> value = leased.parseTagValue(mockTag, "value");

        // Assert
        assertTrue(value.isPresent());
        assertEquals(mockValue, value.get());
        verify(mockDelegate).parseTagValue(mockTag, "value");
    }

    @Test
    void testBrowseRequestBuilder_Delegates() {
        // Arrange
        when(mockDelegate.browseRequestBuilder()).thenReturn(mockBrowseBuilder);
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcBrowseRequest.Builder builder = leased.browseRequestBuilder();

        // Assert
        // The builder is now wrapped, so we can't check for equality
        // Instead verify the delegate was called and the builder works
        assertNotNull(builder);
        verify(mockDelegate).browseRequestBuilder();
    }

    @Test
    void testReadRequestBuilder_Delegates() {
        // Arrange
        when(mockDelegate.readRequestBuilder()).thenReturn(mockReadBuilder);
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcReadRequest.Builder builder = leased.readRequestBuilder();

        // Assert
        // The builder is now wrapped, so we can't check for equality
        // Instead verify the delegate was called and the builder works
        assertNotNull(builder);
        verify(mockDelegate).readRequestBuilder();
    }

    @Test
    void testWriteRequestBuilder_Delegates() {
        // Arrange
        when(mockDelegate.writeRequestBuilder()).thenReturn(mockWriteBuilder);
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcWriteRequest.Builder builder = leased.writeRequestBuilder();

        // Assert
        // The builder is now wrapped, so we can't check for equality
        // Instead verify the delegate was called and the builder works
        assertNotNull(builder);
        verify(mockDelegate).writeRequestBuilder();
    }

    @Test
    void testSubscriptionRequestBuilder_Delegates() {
        // Arrange
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcSubscriptionRequest.Builder builder = leased.subscriptionRequestBuilder();

        // Assert
        assertNotNull(builder);
        verify(mockDelegate).subscriptionRequestBuilder();
    }

    @Test
    void testUnsubscriptionRequestBuilder_Delegates() {
        // Arrange
        when(mockDelegate.unsubscriptionRequestBuilder()).thenReturn(mockUnsubscriptionBuilder);
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcUnsubscriptionRequest.Builder builder = leased.unsubscriptionRequestBuilder();

        // Assert
        assertNotNull(builder);
        verify(mockDelegate).unsubscriptionRequestBuilder();
    }

    @Test
    void testPing_Delegates() {
        // Arrange
        when(mockDelegate.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        CompletableFuture<? extends PlcPingResponse> result = leased.ping();

        // Assert
        assertNotNull(result);
        verify(mockDelegate).ping();
    }

    @Test
    void testParseTagAddress_ExceptionDoesNotInvalidateConnection() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.parseTagAddress("tag1")).thenThrow(new RuntimeException("Parse error"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert
        RuntimeException thrown = assertThrows(RuntimeException.class, () -> {
            leased.parseTagAddress("tag1");
        });

        assertEquals("Parse error", thrown.getMessage());
        assertEquals(0, invalidationCount.get(), "Connection should NOT be invalidated after parse error (user error)");
        verify(mockDelegate).parseTagAddress("tag1");
    }

    @Test
    void testParseTagValue_ExceptionDoesNotInvalidateConnection() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.parseTagValue(mockTag, "value")).thenThrow(new RuntimeException("Parse value error"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert
        RuntimeException thrown = assertThrows(RuntimeException.class, () -> {
            leased.parseTagValue(mockTag, "value");
        });

        assertEquals("Parse value error", thrown.getMessage());
        assertEquals(0, invalidationCount.get(), "Connection should NOT be invalidated after parse error (user error)");
        verify(mockDelegate).parseTagValue(mockTag, "value");
    }

    @Test
    void testBrowseRequestBuilder_ExceptionInvalidatesConnection() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.browseRequestBuilder()).thenThrow(new RuntimeException("Browse error"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert
        RuntimeException thrown = assertThrows(RuntimeException.class, leased::browseRequestBuilder);

        assertEquals("Browse error", thrown.getMessage());
        // Close triggers invalidation callback
        leased.close();
        
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after exception");
        verify(mockDelegate).browseRequestBuilder();
    }

    @Test
    void testReadRequestBuilder_ExceptionInvalidatesConnection() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.readRequestBuilder()).thenThrow(new RuntimeException("Read builder error"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert
        RuntimeException thrown = assertThrows(RuntimeException.class, leased::readRequestBuilder);

        assertEquals("Read builder error", thrown.getMessage());
        // Close triggers invalidation callback
        leased.close();
        
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after exception");
        verify(mockDelegate).readRequestBuilder();
    }

    @Test
    void testWriteRequestBuilder_ExceptionInvalidatesConnection() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.writeRequestBuilder()).thenThrow(new RuntimeException("Write builder error"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert
        RuntimeException thrown = assertThrows(RuntimeException.class, leased::writeRequestBuilder);

        assertEquals("Write builder error", thrown.getMessage());
        // Close triggers invalidation callback
        leased.close();
        
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after exception");
        verify(mockDelegate).writeRequestBuilder();
    }

    @Test
    void testSubscriptionRequestBuilder_ExceptionInvalidatesConnection() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.subscriptionRequestBuilder()).thenThrow(new RuntimeException("Subscription error"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert
        RuntimeException thrown = assertThrows(RuntimeException.class, leased::subscriptionRequestBuilder);

        assertEquals("Subscription error", thrown.getMessage());
        // Close triggers invalidation callback
        leased.close();
        
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after exception");
        verify(mockDelegate).subscriptionRequestBuilder();
    }

    @Test
    void testUnsubscriptionRequestBuilder_ExceptionInvalidatesConnection() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.unsubscriptionRequestBuilder()).thenThrow(new RuntimeException("Unsubscription error"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert
        RuntimeException thrown = assertThrows(RuntimeException.class, leased::unsubscriptionRequestBuilder);

        assertEquals("Unsubscription error", thrown.getMessage());
        // Close triggers invalidation callback
        leased.close();
        
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after exception");
        verify(mockDelegate).unsubscriptionRequestBuilder();
    }

    @Test
    void testPing_ExceptionInFutureInvalidatesConnection() throws Exception {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.ping()).thenAnswer(inv -> {
            CompletableFuture<PlcPingResponse> failedFuture = new CompletableFuture<>();
            failedFuture.completeExceptionally(new RuntimeException("Ping failed"));
            return failedFuture;
        });

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act
        CompletableFuture<? extends PlcPingResponse> result = leased.ping();

        // Wait for the future to complete
        try {
            result.get();
            fail("Expected ExecutionException");
        } catch (Exception e) {
            // Expected
        }

        // Close triggers invalidation callback
        leased.close();

        // Assert
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after ping failure");
        verify(mockDelegate).ping();
    }

    @Test
    void testPing_SynchronousExceptionInvalidatesConnection() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        when(mockDelegate.ping()).thenThrow(new RuntimeException("Ping threw immediately"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert
        RuntimeException thrown = assertThrows(RuntimeException.class, leased::ping);

        assertEquals("Ping threw immediately", thrown.getMessage());

        // Close triggers invalidation callback
        leased.close();

        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after exception");
        verify(mockDelegate).ping();
    }

    @Test
    void testMultipleExceptions_OnlyInvalidatesOnce() {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        // Use protocol errors (not user errors) to trigger invalidation
        when(mockDelegate.readRequestBuilder()).thenThrow(new RuntimeException("Error 1"));
        when(mockDelegate.writeRequestBuilder()).thenThrow(new RuntimeException("Error 2"));

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act - Trigger multiple protocol exceptions
        assertThrows(RuntimeException.class, leased::readRequestBuilder);
        assertThrows(RuntimeException.class, leased::writeRequestBuilder);

        // Close triggers invalidation callback
        leased.close();

        // Assert - Should only invalidate once even with multiple errors
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated only once");
    }

    @Test
    void testInvalidationCallback_ExceptionHandled() throws Exception {
        // Arrange
        Runnable invalidateCallback = () -> {
            throw new RuntimeException("Invalidation callback failed");
        };

        // Use ping error (protocol error) instead of parseTagAddress (user error)
        when(mockDelegate.ping()).thenAnswer(inv -> {
            CompletableFuture<PlcPingResponse> failedFuture = new CompletableFuture<>();
            failedFuture.completeExceptionally(new RuntimeException("Connection error"));
            return failedFuture;
        });

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act & Assert - Should not propagate exception from callback during close
        try {
            leased.ping().get();
        } catch (Exception e) {
            // Expected - connection error
        }

        // Close should trigger invalidation callback, but not propagate its exception
        leased.close();

        // Verify the delegate method was called
        verify(mockDelegate).ping();
    }

    @Test
    void testClose_AfterInvalidation_DoesNotReturnToCache() throws Exception {
        // Arrange
        AtomicInteger returnCount = new AtomicInteger(0);
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable returnCallback = returnCount::incrementAndGet;
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        // Use ping error (protocol error) instead of parseTagAddress (user error)
        when(mockDelegate.ping()).thenAnswer(inv -> {
            CompletableFuture<PlcPingResponse> failedFuture = new CompletableFuture<>();
            failedFuture.completeExceptionally(new RuntimeException("Connection error"));
            return failedFuture;
        });

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) {
                    invalidateCallback.run();
                } else {
                    returnCallback.run();
                }
            }
        );

        // Act - Trigger protocol error with ping
        try {
            leased.ping().get();
        } catch (Exception e) {
            // Expected
        }
        leased.close();

        // Assert
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated");
        assertEquals(0, returnCount.get(), "Invalidated connection should not be returned to cache");
        assertTrue(leased.isReturned());
    }

    @Test
    void testBrowseRequest_Execute_HandlesExceptions() throws Exception {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        PlcBrowseRequest mockRequest = mock(PlcBrowseRequest.class);
        CompletableFuture<PlcBrowseResponse> failedFuture = new CompletableFuture<>();
        failedFuture.completeExceptionally(new RuntimeException("Browse failed"));

        when(mockDelegate.browseRequestBuilder()).thenReturn(mockBrowseBuilder);
        when(mockBrowseBuilder.addQuery(anyString(), anyString())).thenReturn(mockBrowseBuilder);
        when(mockBrowseBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> failedFuture);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act
        CompletableFuture<?> result = leased.browseRequestBuilder()
            .addQuery("query1", "query")
            .build()
            .execute();

        // Wait for the future to complete
        Thread.sleep(100);

        // Assert
        assertTrue(result.isCompletedExceptionally(), "Future should complete exceptionally");
        // Close triggers invalidation callback
        leased.close();
        
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after browse error");
    }

    @Test
    void testReadRequest_Execute_HandlesExceptions() throws Exception {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        PlcReadRequest mockRequest = mock(PlcReadRequest.class);
        CompletableFuture<PlcReadResponse> failedFuture = new CompletableFuture<>();
        failedFuture.completeExceptionally(new RuntimeException("Read failed"));

        when(mockDelegate.readRequestBuilder()).thenReturn(mockReadBuilder);
        when(mockReadBuilder.addTagAddress(anyString(), anyString())).thenReturn(mockReadBuilder);
        when(mockReadBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> failedFuture);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act
        CompletableFuture<?> result = leased.readRequestBuilder()
            .addTagAddress("tag1", "address")
            .build()
            .execute();

        // Wait for the future to complete
        Thread.sleep(100);

        // Assert
        assertTrue(result.isCompletedExceptionally(), "Future should complete exceptionally");
        // Close triggers invalidation callback
        leased.close();
        
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after read error");
    }

    @Test
    void testWriteRequest_Execute_HandlesExceptions() throws Exception {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        CompletableFuture<PlcWriteResponse> failedFuture = new CompletableFuture<>();
        failedFuture.completeExceptionally(new RuntimeException("Write failed"));

        when(mockDelegate.writeRequestBuilder()).thenReturn(mockWriteBuilder);
        when(mockWriteBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockWriteBuilder);
        when(mockWriteBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> failedFuture);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act
        CompletableFuture<?> result = leased.writeRequestBuilder()
            .addTagAddress("tag1", "address", 42)
            .build()
            .execute();

        // Wait for the future to complete
        Thread.sleep(100);

        // Assert
        assertTrue(result.isCompletedExceptionally(), "Future should complete exceptionally");
        // Close triggers invalidation callback
        leased.close();
        
        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after write error");
    }

    @Test
    void testSubscriptionRequest_Execute_HandlesExceptions() throws Exception {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        PlcSubscriptionRequest mockSubscriptionRequest = mock(PlcSubscriptionRequest.class);
        CompletableFuture<PlcSubscriptionResponse> failedFuture = new CompletableFuture<>();
        failedFuture.completeExceptionally(new RuntimeException("Subscription failed"));
        when(mockSubscriptionRequest.execute()).thenAnswer(inv -> failedFuture);

        when(mockSubscriptionBuilder.build()).thenReturn(mockSubscriptionRequest);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act
        CompletableFuture<? extends PlcSubscriptionResponse> result = leased
            .subscriptionRequestBuilder()
            .build()
            .execute();

        // Wait for the future to complete
        Thread.sleep(100);

        // Assert
        assertTrue(result.isCompletedExceptionally(), "Future should complete exceptionally");
        // Close triggers invalidation callback
        leased.close();

        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after subscription error");
    }

    @Test
    void testSubscription_handleRejectsPostSubscribeRegister() throws Exception {
        // A handle handed out by the cache must reject register(...) after subscribing: such a consumer
        // could not be restored across a cache-managed reconnection. Consumers must be set at subscribe time.
        PlcSubscriptionHandle rawHandle = mock(PlcSubscriptionHandle.class);
        PlcSubscriptionTag rawTag = mock(PlcSubscriptionTag.class);
        PlcSubscriptionResponse rawResponse = mock(PlcSubscriptionResponse.class);

        PlcSubscriptionRequest mockSubscriptionRequest = mock(PlcSubscriptionRequest.class);
        when(mockSubscriptionRequest.getConsumer()).thenReturn(null);
        when(rawResponse.getSubscriptionHandle("tag1")).thenReturn(rawHandle);
        when(rawResponse.getSubscriptionHandles()).thenReturn(List.of(rawHandle));
        when(rawResponse.getRequest()).thenReturn(mockSubscriptionRequest);
        when(rawResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(rawResponse.getTag("tag1")).thenReturn(rawTag);
        when(rawResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);
        doReturn(CompletableFuture.completedFuture(rawResponse)).when(mockSubscriptionRequest).execute();
        when(mockSubscriptionBuilder.build()).thenReturn(mockSubscriptionRequest);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> { });

        PlcSubscriptionResponse response = leased.subscriptionRequestBuilder().build().execute().get();

        // The wrapper delegates all non-handle accessors to the real response...
        assertSame(mockSubscriptionRequest, response.getRequest());
        assertTrue(response.getTagNames().contains("tag1"));
        assertSame(rawTag, response.getTag("tag1"));
        assertEquals(PlcResponseCode.OK, response.getResponseCode("tag1"));

        // ...but every handle it hands out (single or collection) rejects register(...).
        PlcSubscriptionHandle handle = response.getSubscriptionHandle("tag1");
        assertThrows(UnsupportedOperationException.class, () -> handle.register(event -> { }));
        Collection<PlcSubscriptionHandle> handles = response.getSubscriptionHandles();
        assertEquals(1, handles.size());
        assertThrows(UnsupportedOperationException.class, () -> handles.iterator().next().register(event -> { }));
    }

    @Test
    void testUnsubscriptionRequest_Execute_HandlesExceptions() throws Exception {
        // Arrange
        AtomicInteger invalidationCount = new AtomicInteger(0);
        Runnable invalidateCallback = invalidationCount::incrementAndGet;

        PlcUnsubscriptionRequest mockUnsubscriptionRequest = mock(PlcUnsubscriptionRequest.class);
        CompletableFuture<PlcUnsubscriptionResponse> failedFuture = new CompletableFuture<>();
        failedFuture.completeExceptionally(new RuntimeException("Unsubscription failed"));
        when(mockUnsubscriptionRequest.execute()).thenAnswer(inv -> failedFuture);

        when(mockUnsubscriptionBuilder.build()).thenReturn(mockUnsubscriptionRequest);
        when(mockDelegate.unsubscriptionRequestBuilder()).thenReturn(mockUnsubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {
                if (invalidated) invalidateCallback.run();
            }
        );

        // Act
        CompletableFuture<? extends PlcUnsubscriptionResponse> result = leased
            .unsubscriptionRequestBuilder()
            .build()
            .execute();

        // Wait for the future to complete
        Thread.sleep(100);

        // Assert
        assertTrue(result.isCompletedExceptionally(), "Future should complete exceptionally");
        // Close triggers invalidation callback
        leased.close();

        assertEquals(1, invalidationCount.get(), "Connection should be invalidated after unsubscription error");
    }

    @Test
    void testSubscriptionRequestBuilder_Build() {
        // Arrange
        PlcSubscriptionRequest mockSubscriptionRequest = mock(PlcSubscriptionRequest.class);

        when(mockSubscriptionBuilder.build()).thenReturn(mockSubscriptionRequest);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcSubscriptionRequest.Builder builder = leased.subscriptionRequestBuilder();
        PlcSubscriptionRequest request = builder.build();

        // Assert
        assertNotNull(request, "Subscription request should not be null");
        verify(mockSubscriptionBuilder).build();

        leased.close();
    }

    @Test
    void testUnsubscriptionRequestBuilder_Build() {
        // Arrange
        PlcUnsubscriptionRequest mockUnsubscriptionRequest = mock(PlcUnsubscriptionRequest.class);

        when(mockUnsubscriptionBuilder.build()).thenReturn(mockUnsubscriptionRequest);
        when(mockDelegate.unsubscriptionRequestBuilder()).thenReturn(mockUnsubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcUnsubscriptionRequest.Builder builder = leased.unsubscriptionRequestBuilder();
        PlcUnsubscriptionRequest request = builder.build();

        // Assert
        assertNotNull(request, "Unsubscription request should not be null");
        verify(mockUnsubscriptionBuilder).build();

        leased.close();
    }

    @Test
    void testSubscriptionRequestBuilder_AddChangeOfStateTagAddress() {
        // Arrange
        when(mockSubscriptionBuilder.addChangeOfStateTagAddress("tag1", "address1")).thenReturn(mockSubscriptionBuilder);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcSubscriptionRequest.Builder builder = leased.subscriptionRequestBuilder();
        builder.addChangeOfStateTagAddress("tag1", "address1");

        // Assert
        verify(mockSubscriptionBuilder).addChangeOfStateTagAddress("tag1", "address1");

        leased.close();
    }

    @Test
    void testSubscriptionRequestBuilder_AddCyclicTagAddress() {
        // Arrange
        when(mockSubscriptionBuilder.addCyclicTagAddress("tag1", "address1", java.time.Duration.ofSeconds(1))).thenReturn(mockSubscriptionBuilder);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcSubscriptionRequest.Builder builder = leased.subscriptionRequestBuilder();
        builder.addCyclicTagAddress("tag1", "address1", java.time.Duration.ofSeconds(1));

        // Assert
        verify(mockSubscriptionBuilder).addCyclicTagAddress("tag1", "address1", java.time.Duration.ofSeconds(1));

        leased.close();
    }

    @Test
    void testSubscriptionRequestBuilder_AddEventTagAddress() {
        // Arrange
        when(mockSubscriptionBuilder.addEventTagAddress("tag1", "address1")).thenReturn(mockSubscriptionBuilder);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcSubscriptionRequest.Builder builder = leased.subscriptionRequestBuilder();
        builder.addEventTagAddress("tag1", "address1");

        // Assert
        verify(mockSubscriptionBuilder).addEventTagAddress("tag1", "address1");

        leased.close();
    }

    @Test
    void testSubscriptionRequestBuilder_AddChangeOfStateTag() {
        // Arrange
        when(mockSubscriptionBuilder.addChangeOfStateTag("tag1", mockTag)).thenReturn(mockSubscriptionBuilder);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcSubscriptionRequest.Builder builder = leased.subscriptionRequestBuilder();
        builder.addChangeOfStateTag("tag1", mockTag);

        // Assert
        verify(mockSubscriptionBuilder).addChangeOfStateTag("tag1", mockTag);

        leased.close();
    }

    @Test
    void testSubscriptionRequestBuilder_AddCyclicTag() {
        // Arrange
        when(mockSubscriptionBuilder.addCyclicTag("tag1", mockTag, java.time.Duration.ofSeconds(1))).thenReturn(mockSubscriptionBuilder);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcSubscriptionRequest.Builder builder = leased.subscriptionRequestBuilder();
        builder.addCyclicTag("tag1", mockTag, java.time.Duration.ofSeconds(1));

        // Assert
        verify(mockSubscriptionBuilder).addCyclicTag("tag1", mockTag, java.time.Duration.ofSeconds(1));

        leased.close();
    }

    @Test
    void testSubscriptionRequestBuilder_AddEventTag() {
        // Arrange
        when(mockSubscriptionBuilder.addEventTag("tag1", mockTag)).thenReturn(mockSubscriptionBuilder);
        when(mockDelegate.subscriptionRequestBuilder()).thenReturn(mockSubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcSubscriptionRequest.Builder builder = leased.subscriptionRequestBuilder();
        builder.addEventTag("tag1", mockTag);

        // Assert
        verify(mockSubscriptionBuilder).addEventTag("tag1", mockTag);

        leased.close();
    }

    @Test
    void testUnsubscriptionRequest_GetHandles() {
        // Arrange
        PlcUnsubscriptionRequest mockUnsubscriptionRequest = mock(PlcUnsubscriptionRequest.class);
        when(mockUnsubscriptionRequest.getSubscriptionHandles()).thenReturn(java.util.Collections.emptyList());

        when(mockUnsubscriptionBuilder.build()).thenReturn(mockUnsubscriptionRequest);
        when(mockDelegate.unsubscriptionRequestBuilder()).thenReturn(mockUnsubscriptionBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcUnsubscriptionRequest request = leased.unsubscriptionRequestBuilder().build();
        var handles = request.getSubscriptionHandles();

        // Assert
        assertNotNull(handles);

        leased.close();
    }

    @Test
    void testBrowseRequestBuilder_AddQuery() {
        // Arrange
        when(mockBrowseBuilder.addQuery("query1", "address1")).thenReturn(mockBrowseBuilder);
        when(mockDelegate.browseRequestBuilder()).thenReturn(mockBrowseBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcBrowseRequest.Builder builder = leased.browseRequestBuilder();
        builder.addQuery("query1", "address1");

        // Assert
        verify(mockBrowseBuilder).addQuery("query1", "address1");

        leased.close();
    }

    @Test
    void testReadRequestBuilder_AddTagAddress() {
        // Arrange
        when(mockReadBuilder.addTagAddress("tag1", "address1")).thenReturn(mockReadBuilder);
        when(mockDelegate.readRequestBuilder()).thenReturn(mockReadBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcReadRequest.Builder builder = leased.readRequestBuilder();
        builder.addTagAddress("tag1", "address1");

        // Assert
        verify(mockReadBuilder).addTagAddress("tag1", "address1");

        leased.close();
    }

    @Test
    void testReadRequestBuilder_AddTag() {
        // Arrange
        when(mockReadBuilder.addTag("tag1", mockTag)).thenReturn(mockReadBuilder);
        when(mockDelegate.readRequestBuilder()).thenReturn(mockReadBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcReadRequest.Builder builder = leased.readRequestBuilder();
        builder.addTag("tag1", mockTag);

        // Assert
        verify(mockReadBuilder).addTag("tag1", mockTag);

        leased.close();
    }

    @Test
    void testWriteRequestBuilder_AddTagAddress() {
        // Arrange
        when(mockWriteBuilder.addTagAddress("tag1", "address1", 42)).thenReturn(mockWriteBuilder);
        when(mockDelegate.writeRequestBuilder()).thenReturn(mockWriteBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcWriteRequest.Builder builder = leased.writeRequestBuilder();
        builder.addTagAddress("tag1", "address1", 42);

        // Assert
        verify(mockWriteBuilder).addTagAddress("tag1", "address1", 42);

        leased.close();
    }

    @Test
    void testWriteRequestBuilder_AddTag() {
        // Arrange
        when(mockWriteBuilder.addTag("tag1", mockTag, 42)).thenReturn(mockWriteBuilder);
        when(mockDelegate.writeRequestBuilder()).thenReturn(mockWriteBuilder);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act
        PlcWriteRequest.Builder builder = leased.writeRequestBuilder();
        builder.addTag("tag1", mockTag, 42);

        // Assert
        verify(mockWriteBuilder).addTag("tag1", mockTag, 42);

        leased.close();
    }

    /**
     * Test that all operations throw PlcRuntimeException after the connection is closed.
     * This verifies Requirement 7: "If the client keeps the leased connection after calling close(),
     * every operation should fail."
     */
    @Test
    void testOperationsAfterClose_ThrowPlcRuntimeException() {
        // Arrange
        when(mockDelegate.isConnected()).thenReturn(true);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act - close the connection
        leased.close();

        // Assert - all operations should fail with PlcRuntimeException
        PlcRuntimeException ex1 = assertThrows(PlcRuntimeException.class,
            leased::isConnected,
            "isConnected() should throw after close");
        assertTrue(ex1.getMessage().contains("after returning it to the cache"),
            "Exception message should mention cache return");

        PlcRuntimeException ex2 = assertThrows(PlcRuntimeException.class,
            leased::getMetadata,
            "getMetadata() should throw after close");
        assertTrue(ex2.getMessage().contains("after returning it to the cache"));

        PlcRuntimeException ex3 = assertThrows(PlcRuntimeException.class,
            () -> leased.parseTagAddress("test"),
            "parseTagAddress() should throw after close");
        assertTrue(ex3.getMessage().contains("after returning it to the cache"));

        PlcRuntimeException ex4 = assertThrows(PlcRuntimeException.class,
            () -> leased.parseTagValue(mockTag, "test"),
            "parseTagValue() should throw after close");
        assertTrue(ex4.getMessage().contains("after returning it to the cache"));

        PlcRuntimeException ex5 = assertThrows(PlcRuntimeException.class,
            leased::readRequestBuilder,
            "readRequestBuilder() should throw after close");
        assertTrue(ex5.getMessage().contains("after returning it to the cache"));

        PlcRuntimeException ex6 = assertThrows(PlcRuntimeException.class,
            leased::writeRequestBuilder,
            "writeRequestBuilder() should throw after close");
        assertTrue(ex6.getMessage().contains("after returning it to the cache"));

        PlcRuntimeException ex7 = assertThrows(PlcRuntimeException.class,
            leased::subscriptionRequestBuilder,
            "subscriptionRequestBuilder() should throw after close");
        assertTrue(ex7.getMessage().contains("after returning it to the cache"));

        PlcRuntimeException ex8 = assertThrows(PlcRuntimeException.class,
            leased::unsubscriptionRequestBuilder,
            "unsubscriptionRequestBuilder() should throw after close");
        assertTrue(ex8.getMessage().contains("after returning it to the cache"));

        PlcRuntimeException ex9 = assertThrows(PlcRuntimeException.class,
            leased::browseRequestBuilder,
            "browseRequestBuilder() should throw after close");
        assertTrue(ex9.getMessage().contains("after returning it to the cache"));

        PlcRuntimeException ex10 = assertThrows(PlcRuntimeException.class,
            leased::ping,
            "ping() should throw after close");
        assertTrue(ex10.getMessage().contains("after returning it to the cache"));

        // Verify that the delegate was never called after close
        verify(mockDelegate, never()).isConnected();
        verify(mockDelegate, never()).getMetadata();
        verify(mockDelegate, never()).parseTagAddress(anyString());
        verify(mockDelegate, never()).readRequestBuilder();
        verify(mockDelegate, never()).writeRequestBuilder();
        verify(mockDelegate, never()).ping();
    }

    // ========== EventPlcConnection Tests ==========

    @Test
    void testAddEventListener_DelegatesToUnderlyingConnection() {
        // Arrange
        ConnectionStateListener mockListener = mock(ConnectionStateListener.class);

        abstract class MockEventConnection implements PlcConnection, EventPlcConnection {}
        MockEventConnection mockEventDelegate = mock(MockEventConnection.class);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockEventDelegate, (conn, invalidated) -> {}
        );

        // Act
        leased.addEventListener(mockListener);

        // Assert
        verify(mockEventDelegate).addEventListener(mockListener);
    }

    @Test
    void testRemoveEventListener_DelegatesToUnderlyingConnection() {
        // Arrange
        ConnectionStateListener mockListener = mock(ConnectionStateListener.class);

        abstract class MockEventConnection implements PlcConnection, EventPlcConnection {}
        MockEventConnection mockEventDelegate = mock(MockEventConnection.class);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockEventDelegate, (conn, invalidated) -> {}
        );

        // Act
        leased.removeEventListener(mockListener);

        // Assert
        verify(mockEventDelegate).removeEventListener(mockListener);
    }

    @Test
    void testAddEventListener_DoesNotThrowWhenConnectionDoesNotSupportEvents() {
        // Arrange - Use a connection that does NOT implement EventPlcConnection
        ConnectionStateListener mockListener = mock(ConnectionStateListener.class);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act & Assert - Should not throw, just log a warning
        assertDoesNotThrow(() -> leased.addEventListener(mockListener));

        // Verify delegate was NOT called (because it doesn't support events)
        verify(mockDelegate, never()).isConnected(); // Just verify no interaction
    }

    @Test
    void testRemoveEventListener_DoesNotThrowWhenConnectionDoesNotSupportEvents() {
        // Arrange - Use a connection that does NOT implement EventPlcConnection
        ConnectionStateListener mockListener = mock(ConnectionStateListener.class);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, (conn, invalidated) -> {}
        );

        // Act & Assert - Should not throw, just log a warning
        assertDoesNotThrow(() -> leased.removeEventListener(mockListener));
    }

    @Test
    void testAddEventListener_ThrowsAfterClose() {
        // Arrange
        ConnectionStateListener mockListener = mock(ConnectionStateListener.class);

        abstract class MockEventConnection implements PlcConnection, EventPlcConnection {}
        MockEventConnection mockEventDelegate = mock(MockEventConnection.class);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockEventDelegate, (conn, invalidated) -> {}
        );

        // Close the connection
        leased.close();

        // Act & Assert
        PlcRuntimeException ex = assertThrows(PlcRuntimeException.class,
            () -> leased.addEventListener(mockListener),
            "addEventListener() should throw after close");
        assertTrue(ex.getMessage().contains("after returning it to the cache"));
    }

    @Test
    void testRemoveEventListener_ThrowsAfterClose() {
        // Arrange
        ConnectionStateListener mockListener = mock(ConnectionStateListener.class);

        abstract class MockEventConnection implements PlcConnection, EventPlcConnection {}
        MockEventConnection mockEventDelegate = mock(MockEventConnection.class);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockEventDelegate, (conn, invalidated) -> {}
        );

        // Close the connection
        leased.close();

        // Act & Assert
        PlcRuntimeException ex = assertThrows(PlcRuntimeException.class,
            () -> leased.removeEventListener(mockListener),
            "removeEventListener() should throw after close");
        assertTrue(ex.getMessage().contains("after returning it to the cache"));
    }

    @Test
    void testAddEventListener_WithNullListener_DoesNotThrow() {
        // Arrange
        abstract class MockEventConnection implements PlcConnection, EventPlcConnection {}
        MockEventConnection mockEventDelegate = mock(MockEventConnection.class);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockEventDelegate, (conn, invalidated) -> {}
        );

        // Act & Assert - Should handle null gracefully
        assertDoesNotThrow(() -> leased.addEventListener(null));
        verify(mockEventDelegate).addEventListener(null);
    }

    @Test
    void testRemoveEventListener_WithNullListener_DoesNotThrow() {
        // Arrange
        abstract class MockEventConnection implements PlcConnection, EventPlcConnection {}
        MockEventConnection mockEventDelegate = mock(MockEventConnection.class);

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockEventDelegate, (conn, invalidated) -> {}
        );

        // Act & Assert - Should handle null gracefully
        assertDoesNotThrow(() -> leased.removeEventListener(null));
        verify(mockEventDelegate).removeEventListener(null);
    }

    // ========== US2: Lease Identity Tests ==========

    /**
     * T008: Verify that when close() is called, the lease ID is forwarded
     * to the return connection callback.
     */
    @Test
    void testLeaseIdPassedToCallback() {
        // Arrange
        long expectedLeaseId = 42;
        long[] capturedLeaseId = {-1};

        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, null,
            (conn, invalidated, leaseId) -> capturedLeaseId[0] = leaseId,
            expectedLeaseId
        );

        // Act
        leased.close();

        // Assert - the lease ID should be forwarded to the callback
        assertEquals(expectedLeaseId, capturedLeaseId[0]);
    }

    /**
     * Verify that the lease ID is immutable and doesn't change during the connection lifetime.
     */
    @Test
    void testLeaseIdImmutable() {
        // Arrange
        long expectedLeaseId = 99;
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate, null,
            (conn, invalidated, leaseId) -> {},
            expectedLeaseId
        );

        // Act & Assert - lease ID is constant
        assertEquals(expectedLeaseId, leased.getLeaseId());
        // Use the connection (should not change leaseId)
        when(mockDelegate.isConnected()).thenReturn(true);
        leased.isConnected();
        assertEquals(expectedLeaseId, leased.getLeaseId());
    }

    /**
     * Verify backward-compatible constructor still works and uses leaseId=0.
     */
    @Test
    void testBackwardCompatConstructor_LeaseIdZero() {
        // Arrange - use the BiConsumer constructor
        long[] capturedLeaseId = {-1};
        LeasedPlcConnection leased = new LeasedPlcConnection(
            "test:tcp://localhost", mockDelegate,
            (conn, invalidated) -> capturedLeaseId[0] = 0 // BiConsumer doesn't see leaseId
        );

        // Assert - default leaseId is 0
        assertEquals(0, leased.getLeaseId());

        // Close should still work
        leased.close();
    }
}
