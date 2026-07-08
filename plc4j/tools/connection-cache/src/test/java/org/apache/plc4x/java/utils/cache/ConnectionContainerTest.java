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
import org.apache.plc4x.java.api.messages.PlcPingResponse;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static org.awaitility.Awaitility.await;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 * Tests for ConnectionContainer.
 */
class ConnectionContainerTest {

    @Mock
    private PlcConnection mockConnection;

    private ScheduledExecutorService scheduler;
    private AutoCloseable mocks;

    @BeforeEach
    void setUp() {
        mocks = MockitoAnnotations.openMocks(this);
        scheduler = Executors.newSingleThreadScheduledExecutor();

        // Setup default ping behavior
        when(mockConnection.isConnected()).thenReturn(true);
        when(mockConnection.ping()).thenReturn(CompletableFuture.completedFuture(null));
    }

    @AfterEach
    void tearDown() throws Exception {
        scheduler.shutdownNow();
        if (mocks != null) {
            mocks.close();
        }
    }

    @Test
    void testLease_BasicFunctionality() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000, // maxIdleTimeMs
            1000, // maxLeaseTimeMs
            30000, // maxWaitTimeMs
            5000, // pingTimeoutMs
            30000 // idlePingThresholdMs
        );

        // Act
        Future<PlcConnection> leaseFuture = container.lease();
        PlcConnection leased = leaseFuture.get();

        // Assert
        assertNotNull(leased);
        assertTrue(container.isLeased());
        assertFalse(container.isClosed());

        // Cleanup
        leased.close();
    }

    @Test
    void testLease_SecondLeaseWaitsInQueue() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,
            1000,
            30000,
            5000,
            30000
        );

        // Act - Lease connection in first thread
        Future<PlcConnection> future1 = container.lease();
        PlcConnection leased1 = future1.get();
        assertTrue(container.isLeased());

        // Try to lease in second thread - should be queued
        Future<PlcConnection> future2 = container.lease();
        assertFalse(future2.isDone()); // Should be waiting

        // Return first lease - this should complete future2
        leased1.close();

        PlcConnection leased2 = future2.get(1, TimeUnit.SECONDS);
        assertNotNull(leased2);
        assertTrue(container.isLeased());

        // Cleanup
        leased2.close();
    }

    @Test
    void testClose_ClosesUnderlyingConnection() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,
            1000,
            30000,
            5000,
            30000
        );

        // Create connection by leasing and returning
        PlcConnection leased = container.lease().get();
        leased.close();

        // Act
        container.close();

        // Assert
        verify(mockConnection, times(1)).close();
        assertTrue(container.isClosed());
    }

    @Test
    void testIdleTimeout_ClosesConnection() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            500, // Short idle timeout for testing
            10000,
            30000,
            5000,
            30000
        );

        // Act - Lease and return
        PlcConnection leased = container.lease().get();
        leased.close();

        assertFalse(container.isClosed());

        // Wait for idle timeout
        Thread.sleep(1000);

        // Assert
        verify(mockConnection, times(1)).close();
        assertTrue(container.isClosed());
    }

    @Test
    void testMaxLeaseTimeout_ForcesReturn() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            10000,
            500, // Short max lease timeout for testing
            30000,
            5000,
            30000
        );

        // Act - Lease but don't return
        PlcConnection leased = container.lease().get();
        assertTrue(container.isLeased());

        // Wait for max lease timeout
        Thread.sleep(1000);

        // Assert - Connection should be returned automatically
        assertFalse(container.isLeased());
    }

    @Test
    void testConnectionValidation_PingsIdleConnection() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,
            1000,
            30000,
            5000,
            100 // Very short idle ping threshold
        );

        // Act - Lease, return, then re-lease once the connection is due for validation.
        PlcConnection leased1 = container.lease().get();
        leased1.close();

        // The container arms validation via a scheduled task ~idlePingThreshold after the return.
        // That scheduler tick can lag under load, so poll the lease/validate path until the idle
        // connection is actually pinged instead of assuming a fixed wall-clock delay.
        await().atMost(Duration.ofSeconds(5)).untilAsserted(() -> {
            container.lease().get().close();
            verify(mockConnection, atLeastOnce()).ping();
        });
    }

    @Test
    void testConnectionValidation_FailedPingCreatesNewConnection() throws Exception {
        // Arrange
        PlcConnection mockConnection2 = mock(PlcConnection.class);
        when(mockConnection2.isConnected()).thenReturn(true);
        when(mockConnection2.ping()).thenReturn(CompletableFuture.completedFuture(null));

        int[] callCount = {0};
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> {
                callCount[0]++;
                return callCount[0] == 1 ? mockConnection : mockConnection2;
            },
            scheduler,
            5000,
            1000,
            30000,
            5000,
            100 // Very short idle ping threshold
        );

        // Make first connection's ping fail
        when(mockConnection.ping()).thenReturn(CompletableFuture.failedFuture(new RuntimeException("Ping failed")));

        // Act - Lease, return, then re-lease once the connection is due for validation.
        PlcConnection leased1 = container.lease().get();
        leased1.close();

        // Validation is armed by a scheduled task ~idlePingThreshold after the return, which can
        // lag under load. Poll the lease/validate path until the (failing) ping has run; once the
        // first connection's ping fails the container discards it and builds a replacement.
        await().atMost(Duration.ofSeconds(5)).untilAsserted(() -> {
            container.lease().get().close();
            verify(mockConnection, atLeastOnce()).ping();
        });

        // Assert - the failed ping replaced the first connection exactly once.
        assertEquals(2, callCount[0]);
        verify(mockConnection, times(1)).close(); // Old connection closed
    }

    @Test
    void testWaitTimeout_ThrowsTimeoutException() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            10000,
            10000,
            500, // Very short maxWaitTime for testing
            5000,
            30000
        );

        // Act - Lease connection and hold it
        PlcConnection leased1 = container.lease().get();
        assertTrue(container.isLeased());

        // Try to lease in second thread - should timeout
        Future<PlcConnection> future2 = container.lease();

        // Assert - Should timeout after 500ms
        Thread.sleep(1000);
        assertTrue(future2.isDone());

        // Verify it completed exceptionally by trying to get the result
        try {
            future2.get();
            fail("Expected ExecutionException due to timeout");
        } catch (java.util.concurrent.ExecutionException e) {
            // Expected - should wrap a TimeoutException
            assertInstanceOf(TimeoutException.class, e.getCause());
        }

        // Cleanup
        leased1.close();
    }

    @Test
    void testConnectionString_ReturnsCorrectValue() {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,
            1000,
            30000,
            5000,
            30000
        );

        // Act
        String connectionString = container.getConnectionString();

        // Assert
        assertEquals("test:tcp://localhost", connectionString);
    }

    /**
     * T001: Verify that after max lease timeout fires, the returned connection
     * requires validation before reuse. With non-zero idlePingThresholdMs,
     * a normal return would NOT force validation — but a timeout return MUST.
     */
    @Test
    void testReturnLeaseFromTimeoutForcesValidation() throws Exception {
        // Arrange - use non-zero idlePingThresholdMs so normal return wouldn't force validation
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,  // maxIdleTimeMs
            500,   // maxLeaseTimeMs - short so timeout fires quickly
            30000, // maxWaitTimeMs
            5000,  // pingTimeoutMs
            30000  // idlePingThresholdMs - long enough that normal return wouldn't trigger validation
        );

        // Act - Lease but don't return; let the max lease timeout fire
        PlcConnection leased1 = container.lease().get();
        Thread.sleep(1000); // Wait for max lease timeout (500ms)

        // Connection should have been returned by timeout
        assertFalse(container.isLeased());

        // Lease again immediately — validation should be forced by the timeout return
        PlcConnection leased2 = container.lease().get();

        // Assert - ping should have been called because timeout forced validation
        verify(mockConnection, atLeastOnce()).ping();

        // Cleanup
        leased2.close();
    }

    /**
     * T002: Verify that normal returnLease (via close() on LeasedPlcConnection)
     * follows the existing validation policy based on idlePingThresholdMs.
     * When returned normally and re-leased immediately, no validation should occur
     * if the idle ping threshold hasn't elapsed.
     */
    @Test
    void testReturnLeaseNormalDoesNotForceValidation() throws Exception {
        // Arrange - use non-zero idlePingThresholdMs
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,  // maxIdleTimeMs
            10000, // maxLeaseTimeMs
            30000, // maxWaitTimeMs
            5000,  // pingTimeoutMs
            30000  // idlePingThresholdMs - long enough that it won't fire during test
        );

        // Act - Lease and return normally
        PlcConnection leased1 = container.lease().get();
        leased1.close(); // Normal return

        // Lease again immediately (before idlePingThresholdMs fires)
        PlcConnection leased2 = container.lease().get();

        // Assert - ping should NOT have been called since we're within idlePingThresholdMs
        verify(mockConnection, never()).ping();

        // Cleanup
        leased2.close();
    }

    /**
     * T003: Full flow test: timeout fires → connection returned → immediate lease →
     * validation ping fails → new connection created.
     */
    @Test
    void testTimeoutReturnThenImmediateLease() throws Exception {
        // Arrange - two connections: first will fail validation, second will succeed
        PlcConnection mockConnection2 = mock(PlcConnection.class);
        when(mockConnection2.isConnected()).thenReturn(true);
        when(mockConnection2.ping()).thenReturn(CompletableFuture.completedFuture(null));

        int[] callCount = {0};
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> {
                callCount[0]++;
                return callCount[0] == 1 ? mockConnection : mockConnection2;
            },
            scheduler,
            5000,  // maxIdleTimeMs
            500,   // maxLeaseTimeMs - short so timeout fires quickly
            30000, // maxWaitTimeMs
            5000,  // pingTimeoutMs
            30000  // idlePingThresholdMs
        );

        // Make first connection report as disconnected so validation fails
        when(mockConnection.isConnected()).thenReturn(true); // Initially connected
        PlcConnection leased1 = container.lease().get();

        // Now make it report as not connected — simulates a network outage
        when(mockConnection.isConnected()).thenReturn(false);

        // Wait for max lease timeout to fire
        Thread.sleep(1000);
        assertFalse(container.isLeased());

        // Lease again — should validate, fail, close old connection, create new
        PlcConnection leased2 = container.lease().get();

        // Assert - a new connection should have been created
        assertEquals(2, callCount[0]);
        verify(mockConnection, times(1)).close(); // Old connection closed

        // Cleanup
        leased2.close();
    }

    @Test
    void testInvalidation_ClosesConnectionAndCreatesNew() throws Exception {
        // Arrange
        PlcConnection mockConnection2 = mock(PlcConnection.class);
        when(mockConnection2.isConnected()).thenReturn(true);
        when(mockConnection2.ping()).thenReturn(CompletableFuture.completedFuture(null));

        int[] callCount = {0};
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> {
                callCount[0]++;
                return callCount[0] == 1 ? mockConnection : mockConnection2;
            },
            scheduler,
            5000,
            1000,
            30000,
            5000,
            30000
        );

        // Act - Lease connection and trigger invalidation through a failed operation
        PlcConnection leased1 = container.lease().get();

        // Simulate an error that invalidates the connection by making ping fail
        when(mockConnection.ping()).thenAnswer(invocation -> {
            CompletableFuture<PlcPingResponse> future = new CompletableFuture<>();
            future.completeExceptionally(new RuntimeException("Connection lost"));
            return future;
        });

        // Trigger the ping to cause invalidation
        try {
            leased1.ping().get();
        } catch (Exception e) {
            // Expected - this will mark the connection for invalidation
        }

        // Close the leased connection (will be invalidated due to the error)
        leased1.close();

        // Next lease should create new connection
        PlcConnection leased2 = container.lease().get();

        // Assert
        assertEquals(2, callCount[0]);
        verify(mockConnection, times(1)).close();

        // Cleanup
        leased2.close();
    }

    // ========== US2: Lease Identity Protection Tests ==========

    /**
     * T009: Verify that invalidation proceeds when the lease ID matches the current active lease.
     * This is the normal invalidation flow — the active lease reports an error and triggers invalidation.
     */
    @Test
    void testInvalidateLeaseMatchingId() throws Exception {
        // Arrange
        PlcConnection mockConnection2 = mock(PlcConnection.class);
        when(mockConnection2.isConnected()).thenReturn(true);
        when(mockConnection2.ping()).thenReturn(CompletableFuture.completedFuture(null));

        int[] callCount = {0};
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> {
                callCount[0]++;
                return callCount[0] == 1 ? mockConnection : mockConnection2;
            },
            scheduler,
            5000,
            10000, // Long max lease timeout — we don't want it to fire
            30000,
            5000,
            30000
        );

        // Act - Lease connection and trigger invalidation (matching lease ID)
        PlcConnection leased = container.lease().get();

        // Make ping fail to trigger invalidation
        when(mockConnection.ping()).thenReturn(CompletableFuture.failedFuture(new RuntimeException("Error")));
        try { leased.ping().get(); } catch (Exception e) { /* expected */ }

        // Close triggers invalidateLease with matching ID
        leased.close();

        // Assert - connection should have been invalidated (closed)
        assertFalse(container.isLeased());
        assertTrue(container.isClosed()); // Connection was closed by invalidation
        verify(mockConnection, times(1)).close();
    }

    /**
     * T010: Verify that stale invalidation (lease ID mismatch) is ignored
     * and the current connection remains open.
     * <p>
     * Scenario: Lease A → timeout → Lease B → Lease A's stale invalidation fires → ignored.
     */
    @Test
    void testInvalidateLeaseStaleId() throws Exception {
        // Arrange — two connections: first will fail validation after timeout
        PlcConnection mockConnection2 = mock(PlcConnection.class);
        when(mockConnection2.isConnected()).thenReturn(true);
        when(mockConnection2.ping()).thenReturn(CompletableFuture.completedFuture(null));

        int[] callCount = {0};
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> {
                callCount[0]++;
                return callCount[0] == 1 ? mockConnection : mockConnection2;
            },
            scheduler,
            5000,
            500,   // Short max lease timeout
            30000,
            5000,
            30000
        );

        // Lease A (leaseId=1)
        PlcConnection leased1 = container.lease().get();

        // Simulate network outage — make validation fail after timeout
        when(mockConnection.isConnected()).thenReturn(false);

        // Wait for max lease timeout to fire
        Thread.sleep(1000);
        assertFalse(container.isLeased());

        // Lease B (leaseId=2) — validation fails, new connection created
        PlcConnection leased2 = container.lease().get();
        assertTrue(container.isLeased());
        assertEquals(2, callCount[0]);

        // Now trigger stale invalidation from Lease A
        when(mockConnection.ping()).thenReturn(CompletableFuture.failedFuture(new RuntimeException("Old error")));
        try { leased1.ping().get(); } catch (Exception e) { /* expected */ }

        // Close Lease A — calls invalidateLease(leaseId=1), but currentLeaseId=2
        leased1.close();

        // Assert — Lease B's connection is NOT affected
        assertTrue(container.isLeased());
        assertFalse(container.isClosed());
        verify(mockConnection2, never()).close(); // Lease B's connection not closed

        // Cleanup
        leased2.close();
    }

    /**
     * T011: Verify that stale return (lease ID mismatch) is ignored.
     * <p>
     * Scenario: Lease A → timeout return → Lease B → Lease A's stale close fires → ignored.
     */
    @Test
    void testReturnLeaseStaleId() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,
            500,   // Short max lease timeout
            30000,
            5000,
            30000
        );

        // Lease A (leaseId=1)
        PlcConnection leased1 = container.lease().get();

        // Wait for max lease timeout
        Thread.sleep(1000);
        assertFalse(container.isLeased());

        // Lease B (leaseId=2) — same connection since ping succeeds
        PlcConnection leased2 = container.lease().get();
        assertTrue(container.isLeased());

        // Close Lease A (normal return, no invalidation)
        // Calls returnLease(leaseId=1, false) — stale, should be ignored
        leased1.close();

        // Assert — Lease B is still active
        assertTrue(container.isLeased());

        // Cleanup
        leased2.close();
    }

    /**
     * T012: Verify no state corruption when lease and invalidation operations
     * execute concurrently on the same container.
     */
    @Test
    void testConcurrentLeaseAndInvalidate() throws Exception {
        // Arrange
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,
            500,   // Short max lease timeout
            30000,
            5000,
            30000
        );

        int iterations = 50;
        java.util.concurrent.ExecutorService executor = java.util.concurrent.Executors.newFixedThreadPool(4);
        java.util.concurrent.CountDownLatch latch = new java.util.concurrent.CountDownLatch(iterations);

        // Act — run many concurrent lease/return cycles
        for (int i = 0; i < iterations; i++) {
            executor.submit(() -> {
                try {
                    Future<PlcConnection> future = container.lease();
                    PlcConnection conn = future.get(2, TimeUnit.SECONDS);
                    Thread.sleep(5); // Brief use
                    conn.close();
                } catch (Exception e) {
                    // Timeouts and other errors are acceptable under contention
                } finally {
                    latch.countDown();
                }
            });
        }

        // Assert — all tasks complete without deadlock
        assertTrue(latch.await(30, TimeUnit.SECONDS), "All tasks should complete without deadlock");
        executor.shutdown();
    }

    /**
     * A hanging connection.close() must not block the container indefinitely: with a close timeout
     * configured, close() is abandoned to the background after the timeout and the container recovers.
     */
    @Test
    void testCloseTimeout_AbandonsHangingClose() throws Exception {
        CountDownLatch closeStarted = new CountDownLatch(1);
        // Simulate a wedged socket: close() blocks far longer than the configured close timeout.
        doAnswer(inv -> {
            closeStarted.countDown();
            Thread.sleep(60_000);
            return null;
        }).when(mockConnection).close();

        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> mockConnection,
            scheduler,
            5000,  // maxIdleTimeMs
            1000,  // maxLeaseTimeMs
            30000, // maxWaitTimeMs
            5000,  // pingTimeoutMs
            30000, // idlePingThresholdMs
            200    // closeTimeoutMs — bounded
        );

        // Establish the connection and return it to the cache (idle).
        container.lease().get(1, TimeUnit.SECONDS).close();

        long start = System.currentTimeMillis();
        container.close(); // closeInternal -> bounded close; must not wait on the 60s hang
        long elapsed = System.currentTimeMillis() - start;

        assertTrue(closeStarted.await(1, TimeUnit.SECONDS), "close() should have been invoked");
        assertTrue(elapsed < 2000,
            "container.close() must be bounded by the close timeout, but took " + elapsed + "ms");
    }

    /**
     * Regression test for issue #2631: a hanging driver {@code connect()} must not hold the container
     * lock indefinitely and deadlock other callers. The connect runs under the lock, so it is bounded
     * by {@code maxWaitTimeMs}; after that the attempt fails and the lock is released, letting the next
     * caller proceed instead of blocking forever on {@code lock.lock()} inside {@code lease()}.
     */
    @Test
    void testConnectTimeout_AbandonsHangingConnect() throws Exception {
        CountDownLatch connectStarted = new CountDownLatch(1);
        // Simulate a wedged connect that never completes (far longer than the configured max wait).
        ConnectionContainer container = new ConnectionContainer(
            "test:tcp://localhost",
            () -> {
                connectStarted.countDown();
                try {
                    Thread.sleep(60_000);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
                return mockConnection;
            },
            scheduler,
            5000,  // maxIdleTimeMs
            1000,  // maxLeaseTimeMs
            300,   // maxWaitTimeMs — bounds the connect
            5000,  // pingTimeoutMs
            30000  // idlePingThresholdMs
        );

        // A first caller triggers the (hanging) connect; lease() must not block past the max wait.
        long start = System.currentTimeMillis();
        Future<PlcConnection> future1 = container.lease();
        long elapsed = System.currentTimeMillis() - start;

        assertTrue(connectStarted.await(1, TimeUnit.SECONDS), "connect should have been invoked");
        assertTrue(elapsed < 2000,
            "lease() must be bounded by the max wait time, but took " + elapsed + "ms");
        // The bounded connect failed, so the lease future completes exceptionally rather than hanging.
        assertTrue(future1.isDone());
        assertThrows(java.util.concurrent.ExecutionException.class, future1::get);

        // The container recovered: it is not stuck leased and a subsequent caller can proceed.
        assertFalse(container.isLeased());
    }
}
