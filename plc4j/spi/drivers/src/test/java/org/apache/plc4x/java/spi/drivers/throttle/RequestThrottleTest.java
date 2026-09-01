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
package org.apache.plc4x.java.spi.drivers.throttle;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

class RequestThrottleTest {

    @Test
    @DisplayName("Should initialize with correct max concurrent requests")
    void testInitialization() {
        RequestThrottle throttle = new RequestThrottle(3);

        assertEquals(3, throttle.getMaxConcurrentRequests());
        assertEquals(3, throttle.getAvailablePermits());
        assertEquals(0, throttle.getInFlightRequests());
    }

    @Test
    @DisplayName("Should reject invalid max concurrent requests")
    void testInvalidInitialization() {
        assertThrows(IllegalArgumentException.class, () -> new RequestThrottle(0));
        assertThrows(IllegalArgumentException.class, () -> new RequestThrottle(-1));
    }

    @Test
    @DisplayName("Should throttle concurrent requests to max limit")
    void testThrottling() throws Exception {
        RequestThrottle throttle = new RequestThrottle(2);
        AtomicInteger concurrentCount = new AtomicInteger(0);
        AtomicInteger maxConcurrent = new AtomicInteger(0);

        // Create 5 slow requests
        CompletableFuture<?>[] futures = new CompletableFuture[5];
        for (int i = 0; i < 5; i++) {
            futures[i] = throttle.execute(() -> {
                int current = concurrentCount.incrementAndGet();
                maxConcurrent.updateAndGet(max -> Math.max(max, current));

                return CompletableFuture.runAsync(() -> {
                    try {
                        Thread.sleep(50); // Simulate work
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    }
                }).thenRun(concurrentCount::decrementAndGet);
            });
        }

        // Wait for all to complete
        CompletableFuture.allOf(futures).get();

        // Verify max concurrent never exceeded 2
        assertEquals(2, maxConcurrent.get(), "Max concurrent requests should not exceed throttle limit");
        assertEquals(0, concurrentCount.get(), "All requests should have completed");
    }

    @Test
    @DisplayName("Should release permit on request completion")
    void testPermitRelease() throws Exception {
        RequestThrottle throttle = new RequestThrottle(1);

        CompletableFuture<String> future = throttle.execute(() ->
            CompletableFuture.completedFuture("test"));

        future.get();

        assertEquals(1, throttle.getAvailablePermits());
        assertEquals(0, throttle.getInFlightRequests());
    }

    @Test
    @DisplayName("Should release permit on request failure")
    void testPermitReleaseOnFailure() throws Exception {
        RequestThrottle throttle = new RequestThrottle(1);

        CompletableFuture<String> future = throttle.execute(() ->
            CompletableFuture.failedFuture(new RuntimeException("test error")));

        try {
            future.get();
            fail("Should have thrown exception");
        } catch (Exception e) {
            // Expected
        }

        assertEquals(1, throttle.getAvailablePermits());
        assertEquals(0, throttle.getInFlightRequests());
    }

    @Test
    @DisplayName("Should increase permits when adjusting max upwards")
    void testAdjustMaxUpwards() {
        RequestThrottle throttle = new RequestThrottle(2);

        throttle.adjustMaxConcurrentRequests(5);

        assertEquals(5, throttle.getMaxConcurrentRequests());
        assertEquals(5, throttle.getAvailablePermits());
    }

    @Test
    @DisplayName("Should decrease permits when adjusting max downwards")
    void testAdjustMaxDownwards() {
        RequestThrottle throttle = new RequestThrottle(5);

        throttle.adjustMaxConcurrentRequests(2);

        assertEquals(2, throttle.getMaxConcurrentRequests());
        assertEquals(2, throttle.getAvailablePermits());
    }

    @Test
    @DisplayName("Should reject invalid adjustment")
    void testInvalidAdjustment() {
        RequestThrottle throttle = new RequestThrottle(2);
        assertThrows(IllegalArgumentException.class, () -> throttle.adjustMaxConcurrentRequests(0));
    }

    @Test
    @DisplayName("Should handle adjustment with in-flight requests")
    void testAdjustWithInFlightRequests() throws Exception {
        RequestThrottle throttle = new RequestThrottle(3);

        // Start 2 long-running requests
        CompletableFuture<Void> blocker1 = new CompletableFuture<>();
        CompletableFuture<Void> blocker2 = new CompletableFuture<>();

        throttle.execute(() -> blocker1);
        throttle.execute(() -> blocker2);

        // Should have 1 permit available, 2 in flight
        assertEquals(1, throttle.getAvailablePermits());
        assertEquals(2, throttle.getInFlightRequests());

        // Adjust down to 2 (should drain the 1 available permit)
        throttle.adjustMaxConcurrentRequests(2);

        assertEquals(2, throttle.getMaxConcurrentRequests());
        assertEquals(0, throttle.getAvailablePermits());
        assertEquals(2, throttle.getInFlightRequests());

        // Complete the requests
        blocker1.complete(null);
        blocker2.complete(null);

        // Give a moment for cleanup
        Thread.sleep(10);

        assertEquals(0, throttle.getInFlightRequests());
        // After adjustment down to 2 and completion of 2 requests, we should have 2 permits
        assertEquals(2, throttle.getAvailablePermits());
    }

    @Test
    @DisplayName("Should track in-flight requests correctly")
    void testInFlightTracking() throws Exception {
        RequestThrottle throttle = new RequestThrottle(2);

        CompletableFuture<Void> blocker = new CompletableFuture<>();

        throttle.execute(() -> blocker);

        assertTrue(throttle.getInFlightRequests() > 0);
        assertEquals(1, throttle.getInFlightRequests());
        assertEquals(1, throttle.getAvailablePermits());

        blocker.complete(null);

        // Give a moment for cleanup
        Thread.sleep(10);

        assertEquals(0, throttle.getInFlightRequests());
        assertEquals(2, throttle.getAvailablePermits());
    }

    @Test
    @DisplayName("Should release permit even when the supplier throws")
    void testSupplierException() throws Exception {
        RequestThrottle throttle = new RequestThrottle(1);

        CompletableFuture<String> future = throttle.execute(() -> {
            throw new RuntimeException("Supplier failed");
        });

        try {
            future.get();
            fail("Should have thrown exception");
        } catch (Exception e) {
            assertInstanceOf(RuntimeException.class, e.getCause());
            assertEquals("Supplier failed", e.getCause().getMessage());
        }

        // Permit should be released even though supplier threw
        assertEquals(1, throttle.getAvailablePermits());
        assertEquals(0, throttle.getInFlightRequests());
    }

    /**
     * Submitting a request must never park the calling thread. A driver that chains requests with
     * thenCompose runs the follow-up on the thread that completed the previous stage - typically
     * the connection's receive thread. If submitting blocked there while all permits were held,
     * that thread could no longer deliver the response that would free a permit, and the
     * connection stalled until an unrelated request timed out.
     */
    @Test
    void submittingDoesNotBlockWhenNoPermitIsAvailable() throws Exception {
        RequestThrottle throttle = new RequestThrottle(1);
        CompletableFuture<String> blocker = new CompletableFuture<>();
        throttle.execute(() -> blocker);

        // No permit left. This must return immediately rather than parking this thread.
        AtomicBoolean started = new AtomicBoolean();
        CompletableFuture<String> queued = throttle.execute(() -> {
            started.set(true);
            return CompletableFuture.completedFuture("second");
        });

        assertFalse(started.get(), "the queued request must not start before a permit is free");
        assertFalse(queued.isDone());

        blocker.complete("first");

        assertEquals("second", queued.get(5, TimeUnit.SECONDS));
        assertEquals(1, throttle.getAvailablePermits());
    }

    /**
     * The deadlock itself: the queued request is released by the completion of the first one, from
     * the very thread that is inside the first request's completion callback.
     */
    @Test
    void requestCompletingFromWithinAnotherRequestDoesNotDeadlock() throws Exception {
        RequestThrottle throttle = new RequestThrottle(1);
        CompletableFuture<String> first = new CompletableFuture<>();
        CompletableFuture<String> firstResult = throttle.execute(() -> first);

        // Chain a second request onto the first, the way a driver splits a large request.
        CompletableFuture<String> chained = firstResult.thenCompose(
            ignored -> throttle.execute(() -> CompletableFuture.completedFuture("chained")));

        // Completing from this thread runs the continuation - and therefore the second submit -
        // on this thread, while the permit is still held by the first request.
        first.complete("first");

        assertEquals("chained", chained.get(5, TimeUnit.SECONDS));
        assertEquals(1, throttle.getAvailablePermits());
        assertEquals(0, throttle.getInFlightRequests());
    }

    /**
     * Queued requests start in submission order.
     */
    @Test
    void queuedRequestsRunInOrder() throws Exception {
        RequestThrottle throttle = new RequestThrottle(1);
        CompletableFuture<String> blocker = new CompletableFuture<>();
        throttle.execute(() -> blocker);

        List<Integer> order = Collections.synchronizedList(new ArrayList<>());
        List<CompletableFuture<String>> queued = new ArrayList<>();
        for (int i = 0; i < 5; i++) {
            int index = i;
            queued.add(throttle.execute(() -> {
                order.add(index);
                return CompletableFuture.completedFuture("done" + index);
            }));
        }

        blocker.complete("first");
        for (CompletableFuture<String> future : queued) {
            future.get(5, TimeUnit.SECONDS);
        }

        assertEquals(List.of(0, 1, 2, 3, 4), order);
        assertEquals(1, throttle.getAvailablePermits());
    }

    /**
     * Raising the limit has to let already-queued requests start.
     */
    @Test
    void raisingTheLimitStartsQueuedRequests() throws Exception {
        RequestThrottle throttle = new RequestThrottle(1);
        throttle.execute(() -> new CompletableFuture<>());

        CompletableFuture<String> queued = throttle.execute(() -> CompletableFuture.completedFuture("ok"));
        assertFalse(queued.isDone());

        throttle.adjustMaxConcurrentRequests(2);

        assertEquals("ok", queued.get(5, TimeUnit.SECONDS));
    }
}