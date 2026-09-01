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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Queue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.Semaphore;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.function.Supplier;

/**
 * Manages concurrent request throttling to prevent overwhelming PLCs/devices.
 * Uses a semaphore to limit the number of requests that can be in-flight simultaneously.
 */
public class RequestThrottle {

    private static final Logger LOGGER = LoggerFactory.getLogger(RequestThrottle.class);

    private final Semaphore semaphore;
    private volatile int maxConcurrentRequests;
    /** Requests waiting for a permit, started in the order they were submitted. */
    private final Queue<Runnable> pending = new ConcurrentLinkedQueue<>();
    private final AtomicBoolean draining = new AtomicBoolean();

    public RequestThrottle(int maxConcurrentRequests) {
        if (maxConcurrentRequests < 1) {
            throw new IllegalArgumentException("Max concurrent requests must be at least 1, got: " + maxConcurrentRequests);
        }
        this.maxConcurrentRequests = maxConcurrentRequests;
        this.semaphore = new Semaphore(maxConcurrentRequests, true);
        LOGGER.debug("Request throttle initialized with {} max concurrent requests", maxConcurrentRequests);
    }

    /**
     * Acquires a permit, blocking until one is available.
     * Must be paired with {@link #release()} — typically in a response callback.
     */
    public void acquire() throws InterruptedException {
        semaphore.acquire();
    }

    /**
     * Releases a permit previously acquired via {@link #acquire()}.
     */
    public void release() {
        semaphore.release();
    }

    public <T> CompletableFuture<T> execute(Supplier<CompletableFuture<T>> requestSupplier) {
        CompletableFuture<T> result = new CompletableFuture<>();
        pending.add(() -> start(requestSupplier, result));
        drain();
        return result;
    }

    /**
     * Starts as many queued requests as there are permits. Only one thread drains at a time; any
     * other thread hands its work to whoever holds the drain and returns immediately, so a request
     * that completes inline cannot recurse back into the queue.
     */
    private void drain() {
        while (true) {
            if (!draining.compareAndSet(false, true)) {
                return;
            }
            try {
                while (!pending.isEmpty() && semaphore.tryAcquire()) {
                    Runnable request = pending.poll();
                    if (request == null) {
                        // Another thread took the entry between the check and the poll.
                        semaphore.release();
                        break;
                    }
                    request.run();
                }
            } finally {
                draining.set(false);
            }
            // Something may have been queued by another thread while we were finishing the loop
            // above; without this re-check it would sit there until the next permit is released.
            if (pending.isEmpty() || semaphore.availablePermits() == 0) {
                return;
            }
        }
    }

    private <T> void start(Supplier<CompletableFuture<T>> requestSupplier, CompletableFuture<T> result) {
        CompletableFuture<T> requestFuture;
        try {
            requestFuture = requestSupplier.get();
        } catch (Exception e) {
            releaseAndDrain();
            result.completeExceptionally(e);
            return;
        }
        requestFuture.whenComplete((value, error) -> {
            releaseAndDrain();
            if (error != null) {
                result.completeExceptionally(error);
            } else {
                result.complete(value);
            }
        });
    }

    private void releaseAndDrain() {
        semaphore.release();
        drain();
    }

    public synchronized void adjustMaxConcurrentRequests(int newMax) {
        if (newMax < 1) {
            throw new IllegalArgumentException("Max concurrent requests must be at least 1, got: " + newMax);
        }
        if (newMax == this.maxConcurrentRequests) {
            return;
        }

        int difference = newMax - this.maxConcurrentRequests;
        if (difference > 0) {
            semaphore.release(difference);
            drain();
        } else {
            int toDrain = Math.min(-difference, semaphore.availablePermits());
            if (toDrain > 0) {
                semaphore.acquireUninterruptibly(toDrain);
            }
        }
        this.maxConcurrentRequests = newMax;
        LOGGER.info("Adjusted max concurrent requests to {}", newMax);
    }

    public int getAvailablePermits() {
        return semaphore.availablePermits();
    }

    public int getMaxConcurrentRequests() {
        return maxConcurrentRequests;
    }

    public int getInFlightRequests() {
        return maxConcurrentRequests - semaphore.availablePermits();
    }

}
