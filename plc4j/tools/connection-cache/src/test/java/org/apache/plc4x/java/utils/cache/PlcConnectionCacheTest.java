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

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionFactory;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.messages.PlcPingResponse;
import org.apache.plc4x.java.utils.cache.exceptions.PlcConnectionCacheClosedException;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

/**
 * Tests for PlcConnectionCache.
 * Verifies connection caching, timeouts, thread safety, and proper cleanup.
 */
class PlcConnectionCacheTest {

    @Mock
    private PlcConnectionFactory mockConnectionFactory;

    @Mock
    private PlcConnection mockConnection;

    @Mock
    private PlcPingResponse mockPingResponse;

    private PlcConnectionCache cache;
    private ScheduledExecutorService scheduler;
    private AutoCloseable mocks;

    @BeforeEach
    void setUp() {
        mocks = MockitoAnnotations.openMocks(this);
        scheduler = Executors.newSingleThreadScheduledExecutor();
    }

    @AfterEach
    void tearDown() throws Exception {
        if (cache != null) {
            cache.close();
        }
        scheduler.shutdownNow();
        if (mocks != null) {
            mocks.close();
        }
    }

    /**
     * Test basic connection acquisition and return.
     */
    @Test
    void testGetConnection_BasicAcquisition() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        // Act
        PlcConnection connection = cache.getConnection("test:tcp://localhost");

        // Assert
        assertNotNull(connection);
        assertInstanceOf(LeasedPlcConnection.class, connection);
        verify(mockConnectionFactory, times(1)).getConnection("test:tcp://localhost");
        assertEquals(1, cache.getCachedConnectionCount());
        assertEquals(1, cache.getActiveLeaseCount());

        // Return connection
        connection.close();
        assertEquals(0, cache.getActiveLeaseCount());
    }

    /**
     * Test that connections are reused from cache.
     */
    @Test
    void testGetConnection_Reuse() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Act - Second acquisition (should reuse)
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");
        connection2.close();

        // Assert
        verify(mockConnectionFactory, times(1)).getConnection("test:tcp://localhost"); // Only created once
        assertEquals(1, cache.getCachedConnectionCount());
    }

    /**
     * Test idle timeout closes unused connections.
     */
    @Test
    void testIdleTimeout() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withMaxIdleTime(500, TimeUnit.MILLISECONDS) // Short timeout for testing
            .build();

        // Act
        PlcConnection connection = cache.getConnection("test:tcp://localhost");
        connection.close(); // Return to cache

        assertEquals(1, cache.getCachedConnectionCount());

        // Wait for idle timeout
        Thread.sleep(1000);

        // Assert - Connection should be closed and removed from cache
        verify(mockConnection, times(1)).close();
        // Note: The cache still contains the entry, but it's marked as closed
        // The next getConnection will remove it
    }

    /**
     * Test max lease timeout force-returns leased connections to the pool.
     */
    @Test
    void testMaxLeaseTimeout() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withMaxLeaseTime(500, TimeUnit.MILLISECONDS) // Short timeout for testing
            .build();

        // Act - Lease but don't return
        PlcConnection connection = cache.getConnection("test:tcp://localhost");

        assertEquals(1, cache.getActiveLeaseCount());

        // Wait for max lease timeout
        Thread.sleep(1000);

        // Assert - Connection should be force-returned to pool (not closed, but available)
        // The lease count should be back to 0
        assertEquals(0, cache.getActiveLeaseCount());

        // Connection should still be alive in the pool, not closed
        verify(mockConnection, never()).close();
    }

    /**
     * Test connection validation via ping.
     */
    @Test
    void testConnectionValidation_PingFails() throws Exception {
        // Arrange
        PlcConnection goodConnection = mock(PlcConnection.class);
        PlcConnection badConnection = mock(PlcConnection.class);

        when(mockConnectionFactory.getConnection(anyString()))
            .thenReturn(badConnection)
            .thenReturn(goodConnection);

        // First connection: ping fails
        when(badConnection.isConnected()).thenReturn(true);
        when(badConnection.ping()).thenAnswer(inv -> CompletableFuture.failedFuture(new RuntimeException("Ping failed")));

        // Second connection: ping succeeds
        when(goodConnection.isConnected()).thenReturn(true);
        when(goodConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withIdlePingThreshold(0, TimeUnit.SECONDS) // Always validate
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Simulate the cached connection going bad
        when(badConnection.ping()).thenAnswer(inv -> CompletableFuture.failedFuture(new RuntimeException("Ping failed")));

        // Act - Second acquisition should detect bad connection and create new one
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");

        // Assert - Should have created TWO connections (first one failed validation)
        verify(mockConnectionFactory, times(2)).getConnection("test:tcp://localhost");
        verify(badConnection, times(1)).close(); // Bad connection was closed

        connection2.close();
    }

    /**
     * Test connection validation when connection reports not connected.
     */
    @Test
    void testConnectionValidation_NotConnected() throws Exception {
        // Arrange
        PlcConnection badConnection = mock(PlcConnection.class);
        PlcConnection goodConnection = mock(PlcConnection.class);

        when(mockConnectionFactory.getConnection(anyString()))
            .thenReturn(badConnection)
            .thenReturn(goodConnection);

        // First connection initially connected
        when(badConnection.isConnected()).thenReturn(true);
        when(badConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        // Second connection good
        when(goodConnection.isConnected()).thenReturn(true);
        when(goodConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withIdlePingThreshold(0, TimeUnit.SECONDS) // Always validate
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Simulate connection going offline
        when(badConnection.isConnected()).thenReturn(false);

        // Act - Second acquisition should detect offline connection
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");

        // Assert
        verify(mockConnectionFactory, times(2)).getConnection("test:tcp://localhost");
        verify(badConnection, times(1)).close();

        connection2.close();
    }

    /**
     * Test thread safety with concurrent connection requests.
     * With the new blocking behavior, only ONE client can use a connection at a time.
     * Other clients must wait their turn.
     */
    @Test
    void testConcurrentAccess() throws Exception {
        // Arrange
        AtomicInteger connectionCount = new AtomicInteger(0);

        when(mockConnectionFactory.getConnection(anyString())).thenAnswer(invocation -> {
            connectionCount.incrementAndGet();
            PlcConnection conn = mock(PlcConnection.class);
            when(conn.isConnected()).thenReturn(true);
            when(conn.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));
            return conn;
        });

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        // Act - Multiple threads requesting same connection
        // Each thread gets the connection, uses it briefly, then returns it
        int threadCount = 10;
        ExecutorService executor = Executors.newFixedThreadPool(threadCount);
        List<Future<Boolean>> futures = new ArrayList<>();

        for (int i = 0; i < threadCount; i++) {
            futures.add(executor.submit(() -> {
                try (PlcConnection conn = cache.getConnection("test:tcp://localhost")) {
                    // Use the connection briefly
                    Thread.sleep(10);
                    return true;
                } catch (Exception e) {
                    return false;
                }
            }));
        }

        // Wait for all threads
        for (Future<Boolean> future : futures) {
            assertTrue(future.get(10, TimeUnit.SECONDS), "Thread should successfully get and release connection");
        }

        // Assert - Only ONE physical connection should have been created
        // because all threads waited their turn to use it
        assertEquals(1, connectionCount.get(), "Should create exactly 1 connection (all threads reused it)");
        assertEquals(1, cache.getCachedConnectionCount(), "Should have 1 cached connection");

        // Cleanup
        executor.shutdown();
        assertTrue(executor.awaitTermination(5, TimeUnit.SECONDS));
    }

    /**
     * Test that the cache can be closed and prevents further use.
     */
    @Test
    void testClose() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        PlcConnection connection = cache.getConnection("test:tcp://localhost");
        connection.close();

        // Act
        cache.close();

        // Assert
        verify(mockConnection, times(1)).close();
        assertEquals(0, cache.getCachedConnectionCount());

        // Verify cannot get connections after close
        assertThrows(PlcConnectionCacheClosedException.class, () -> cache.getConnection("test:tcp://localhost"));
    }

    /**
     * The caller that needs the teardown (GH-1399) usually holds the interface type, not this
     * class: closing through PlcConnectionManager in a try-with-resources block has to release the
     * cached connections just the same, without any downcast.
     */
    @Test
    void testCloseThroughTheConnectionManagerInterface() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        // Act - the reference is the interface, and it is the try-with-resources resource
        try (PlcConnectionManager connectionManager = cache) {
            PlcConnection connection = connectionManager.getConnection("test:tcp://localhost");
            connection.close();
            assertEquals(1, cache.getCachedConnectionCount());
        }

        // Assert
        verify(mockConnection, times(1)).close();
        assertEquals(0, cache.getCachedConnectionCount());
        assertThrows(PlcConnectionCacheClosedException.class,
            () -> cache.getConnection("test:tcp://localhost"));
    }

    /**
     * Test that connect() cannot be called on leased connections.
     */
    @Test
    void testLeasedConnection_ConnectNotAllowed() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        // Act
        PlcConnection connection = cache.getConnection("test:tcp://localhost");

        // Assert
        assertThrows(UnsupportedOperationException.class, connection::connect);

        connection.close();
    }

    /**
     * Test that closing a leased connection multiple times is safe.
     */
    @Test
    void testLeasedConnection_MultipleClose() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        // Act
        PlcConnection connection = cache.getConnection("test:tcp://localhost");
        assertEquals(1, cache.getActiveLeaseCount());

        connection.close();
        assertEquals(0, cache.getActiveLeaseCount());

        // Second close should be safe (no-op)
        connection.close();
        assertEquals(0, cache.getActiveLeaseCount());
    }

    /**
     * Test builder validation.
     */
    @Test
    void testBuilder_RequiresDriver() {
        // Act & Assert
        assertThrows(IllegalStateException.class, () ->
            PlcConnectionCache.getBuilder().build()
        );
    }

    /**
     * Test different connection strings get different connections.
     */
    @Test
    void testMultipleConnectionStrings() throws Exception {
        // Arrange
        PlcConnection connection1 = mock(PlcConnection.class);
        PlcConnection connection2 = mock(PlcConnection.class);

        when(mockConnectionFactory.getConnection("test:tcp://host1")).thenReturn(connection1);
        when(mockConnectionFactory.getConnection("test:tcp://host2")).thenReturn(connection2);

        when(connection1.isConnected()).thenReturn(true);
        when(connection1.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));
        when(connection2.isConnected()).thenReturn(true);
        when(connection2.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        // Act
        PlcConnection conn1 = cache.getConnection("test:tcp://host1");
        PlcConnection conn2 = cache.getConnection("test:tcp://host2");

        // Assert
        assertEquals(2, cache.getCachedConnectionCount());
        assertEquals(2, cache.getActiveLeaseCount());

        conn1.close();
        conn2.close();
    }

    /**
     * Test connection with authentication.
     */
    @Test
    void testGetConnection_WithAuthentication() throws Exception {
        // Arrange
        PlcAuthentication mockAuth = mock(PlcAuthentication.class);
        PlcConnection authConnection = mock(PlcConnection.class);

        when(mockConnectionFactory.getConnection(anyString(), any(PlcAuthentication.class))).thenReturn(authConnection);
        when(authConnection.isConnected()).thenReturn(true);
        when(authConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        // Act
        PlcConnection connection = cache.getConnection("test:tcp://localhost", mockAuth);

        // Assert
        assertNotNull(connection);
        verify(mockConnectionFactory, times(1)).getConnection("test:tcp://localhost", mockAuth);
        assertEquals(1, cache.getCachedConnectionCount());
        assertEquals(1, cache.getActiveLeaseCount());

        connection.close();
    }

    /**
     * Test that authentication is used when reconnecting after validation failure.
     */
    @Test
    void testGetConnection_AuthenticationUsedOnReconnect() throws Exception {
        // Arrange
        PlcAuthentication mockAuth = mock(PlcAuthentication.class);
        PlcConnection badConnection = mock(PlcConnection.class);
        PlcConnection goodConnection = mock(PlcConnection.class);

        when(mockConnectionFactory.getConnection(anyString(), any(PlcAuthentication.class)))
            .thenReturn(badConnection)
            .thenReturn(goodConnection);

        // First connection fails validation
        when(badConnection.isConnected()).thenReturn(true);
        when(badConnection.ping()).thenAnswer(inv -> CompletableFuture.failedFuture(new RuntimeException("Ping failed")));

        // Second connection succeeds
        when(goodConnection.isConnected()).thenReturn(true);
        when(goodConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withIdlePingThreshold(0, TimeUnit.SECONDS) // Always ping
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost", mockAuth);
        connection1.close();

        // Simulate validation failure
        when(badConnection.ping()).thenAnswer(inv -> CompletableFuture.failedFuture(new RuntimeException("Ping failed")));

        // Act - Second acquisition should use authentication when creating new connection
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost", mockAuth);

        // Assert - Should have called getConnection with authentication twice
        verify(mockConnectionFactory, times(2)).getConnection("test:tcp://localhost", mockAuth);
        verify(badConnection, times(1)).close();

        connection2.close();
    }

    /**
     * Test that connections are NOT pinged when recently used (below idle threshold).
     */
    @Test
    void testIdlePingThreshold_RecentlyUsedSkipsPing() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withIdlePingThreshold(5, TimeUnit.SECONDS) // 5 second threshold
            .build();

        // Act - First acquisition (will ping because new connection)
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Reset ping mock to track subsequent calls
        reset(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        // Act - Second acquisition immediately (should NOT ping)
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");

        // Assert - ping() should NOT be called because connection was just used
        verify(mockConnection, never()).ping();

        connection2.close();
    }

    /**
     * Test that connections ARE pinged when idle beyond threshold.
     */
    @Test
    void testIdlePingThreshold_IdleConnectionIsPinged() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withIdlePingThreshold(500, TimeUnit.MILLISECONDS) // 500ms threshold
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Wait for the connection to become idle beyond the 500ms threshold. The ping decision is
        // made deterministically from the idle duration at lease() time (no scheduled task
        // involved), so any wall-clock margin over the threshold suffices; 800ms is comfortably
        // clear of scheduling jitter on a loaded CI runner.
        Thread.sleep(800);

        // Reset ping mock to track subsequent calls
        reset(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        // Act - Second acquisition after idle period (should ping)
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");

        // Assert - ping() SHOULD be called because connection was idle too long
        verify(mockConnection, times(1)).ping();

        connection2.close();
    }

    /**
     * Test that idle connections that fail ping validation are replaced.
     */
    @Test
    void testIdlePingThreshold_FailedPingReplacesConnection() throws Exception {
        // Arrange
        PlcConnection badConnection = mock(PlcConnection.class);
        PlcConnection goodConnection = mock(PlcConnection.class);

        when(mockConnectionFactory.getConnection(anyString()))
            .thenReturn(badConnection)
            .thenReturn(goodConnection);

        // First connection initially works
        when(badConnection.isConnected()).thenReturn(true);
        when(badConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        // Second connection works
        when(goodConnection.isConnected()).thenReturn(true);
        when(goodConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withIdlePingThreshold(500, TimeUnit.MILLISECONDS) // 500ms threshold
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Wait for connection to become idle
        Thread.sleep(600);

        // Simulate connection going bad
        when(badConnection.ping()).thenAnswer(inv -> CompletableFuture.failedFuture(new RuntimeException("Ping failed")));

        // Act - Second acquisition should detect bad connection and create new one
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");

        // Assert - Should have created TWO connections (first one failed validation)
        verify(mockConnectionFactory, times(2)).getConnection("test:tcp://localhost");
        verify(badConnection, times(1)).close(); // Bad connection was closed

        connection2.close();
    }

    /**
     * Test custom idle ping threshold configuration.
     */
    @Test
    void testIdlePingThreshold_CustomConfiguration() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        // Build with custom 1 second threshold
        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withIdlePingThreshold(1, TimeUnit.SECONDS)
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Reset mock
        reset(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        // Act - Immediate reacquisition (should NOT ping)
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");
        verify(mockConnection, never()).ping();
        connection2.close();

        // Wait for custom threshold to pass
        Thread.sleep(1100);

        // Reset mock again
        reset(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        // Act - Third acquisition after idle period (should ping)
        PlcConnection connection3 = cache.getConnection("test:tcp://localhost");
        verify(mockConnection, times(1)).ping();

        connection3.close();
    }

    /**
     * Test that idle threshold of 0 always validates connections.
     */
    @Test
    void testIdlePingThreshold_ZeroAlwaysPings() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .withIdlePingThreshold(0, TimeUnit.SECONDS) // Always validate
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Reset ping mock to track subsequent calls
        reset(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenAnswer(inv -> CompletableFuture.completedFuture(mockPingResponse));

        // Act - Second acquisition immediately (should STILL ping because threshold is 0)
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");

        // Assert - ping() SHOULD be called even though connection was just used
        verify(mockConnection, times(1)).ping();

        connection2.close();
    }

    @Test
    void testBuilder_withMaxWaitTime() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withMaxWaitTime(100, TimeUnit.MILLISECONDS) // Very short wait time
            .build();

        // Act - Acquire connection and hold it
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");

        // Try to acquire in another thread - should timeout quickly
        ExecutorService executor = Executors.newSingleThreadExecutor();
        Future<Exception> future = executor.submit(() -> {
            try {
                PlcConnection connection2 = cache.getConnection("test:tcp://localhost");
                connection2.close();
                return null;
            } catch (Exception e) {
                return e;
            }
        });

        // Assert - Should timeout after ~100ms
        Exception exception = future.get(2, TimeUnit.SECONDS);
        assertNotNull(exception, "Expected timeout exception but got successful connection");
        assertNotNull(exception.getMessage(), "Exception message is null");
        assertTrue(exception.getMessage().contains("Error acquiring lease"),
            "Expected 'Error acquiring lease' in message but got: " + exception.getMessage());
        // Check that the cause chain contains a TimeoutException
        // The exception structure is: PlcConnectionException -> ExecutionException -> TimeoutException
        Throwable cause = exception.getCause();
        assertNotNull(cause, "Exception should have a cause");
        // The cause could be ExecutionException wrapping TimeoutException, or TimeoutException directly
        if (cause instanceof ExecutionException) {
            assertInstanceOf(TimeoutException.class, cause.getCause(), "Expected TimeoutException as nested cause but got: " + cause.getCause());
        } else {
            assertInstanceOf(TimeoutException.class, cause, "Expected TimeoutException as cause but got: " + cause);
        }

        // Cleanup
        connection1.close();
        executor.shutdown();
    }

    @Test
    void testBuilder_withPingTimeout() throws Exception {
        // Arrange
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);
        // Mock ping to timeout - never complete the future
        when(mockConnection.ping()).thenAnswer(inv -> new CompletableFuture<PlcPingResponse>());

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withIdlePingThreshold(0, TimeUnit.MILLISECONDS) // Always ping
            .withPingTimeout(100, TimeUnit.MILLISECONDS) // Very short ping timeout
            .build();

        // Act - First acquisition
        PlcConnection connection1 = cache.getConnection("test:tcp://localhost");
        connection1.close();

        // Wait a bit to ensure connection is idle
        Thread.sleep(50);

        // Act - Second acquisition should ping and timeout
        PlcConnection connection2 = cache.getConnection("test:tcp://localhost");

        // Assert - New connection should have been created due to ping timeout
        verify(mockConnectionFactory, times(2)).getConnection(anyString());
        verify(mockConnection, times(1)).ping();

        connection2.close();
    }

    /**
     * getCachedConnections() returns a snapshot of the cached connection strings.
     */
    @Test
    void testGetCachedConnections_ReturnsSnapshot() throws Exception {
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        assertTrue(cache.getCachedConnections().isEmpty());

        cache.getConnection("test:tcp://a").close();
        cache.getConnection("test:tcp://b").close();

        java.util.Set<String> cached = cache.getCachedConnections();
        assertEquals(2, cached.size());
        assertTrue(cached.contains("test:tcp://a"));
        assertTrue(cached.contains("test:tcp://b"));

        // The returned set is a copy — mutating it must not affect the cache.
        cached.clear();
        assertEquals(2, cache.getCachedConnectionCount());
    }

    /**
     * removeCachedConnection() evicts and closes a cached (idle) connection; a subsequent
     * getConnection() transparently establishes a fresh one.
     */
    @Test
    void testRemoveCachedConnection_EvictsAndCloses() throws Exception {
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        cache.getConnection("test:tcp://localhost").close(); // return to cache (idle)
        assertEquals(1, cache.getCachedConnectionCount());

        boolean removed = cache.removeCachedConnection("test:tcp://localhost");

        assertTrue(removed);
        assertEquals(0, cache.getCachedConnectionCount());
        assertFalse(cache.getCachedConnections().contains("test:tcp://localhost"));
        verify(mockConnection, atLeastOnce()).close();

        // Next acquisition establishes a brand-new connection.
        cache.getConnection("test:tcp://localhost").close();
        verify(mockConnectionFactory, times(2)).getConnection("test:tcp://localhost");
    }

    /**
     * removeCachedConnection() force-closes even a currently-leased connection.
     */
    @Test
    void testRemoveCachedConnection_ForceClosesLeased() throws Exception {
        when(mockConnectionFactory.getConnection(anyString())).thenReturn(mockConnection);
        when(mockConnection.isConnected()).thenReturn(true);

        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        PlcConnection leased = cache.getConnection("test:tcp://localhost"); // held, not returned
        assertEquals(1, cache.getActiveLeaseCount());

        boolean removed = cache.removeCachedConnection("test:tcp://localhost");

        assertTrue(removed);
        assertEquals(0, cache.getCachedConnectionCount());
        verify(mockConnection, atLeastOnce()).close();
        assertNotNull(leased); // the now-stale lease reference still exists; its underlying connection is closed
    }

    /**
     * removeCachedConnection() returns false when nothing is cached for the string.
     */
    @Test
    void testRemoveCachedConnection_AbsentReturnsFalse() throws Exception {
        cache = PlcConnectionCache.getBuilder()
            .withConnectionFactory(mockConnectionFactory)
            .withScheduler(scheduler)
            .build();

        assertFalse(cache.removeCachedConnection("test:tcp://never-cached"));
    }
}
