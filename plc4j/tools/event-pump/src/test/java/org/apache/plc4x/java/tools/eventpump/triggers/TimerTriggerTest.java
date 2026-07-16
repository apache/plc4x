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

package org.apache.plc4x.java.tools.eventpump.triggers;

import org.junit.jupiter.api.Test;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for TimerTrigger.
 */
class TimerTriggerTest {

    @Test
    void testBasicFunctionality() throws Exception {
        // Arrange
        CountDownLatch latch = new CountDownLatch(3); // Wait for 3 firings
        AtomicInteger fireCount = new AtomicInteger(0);

        TimerTrigger trigger = new TimerTrigger(100, TimeUnit.MILLISECONDS);

        // Act
        trigger.start(t -> {
            fireCount.incrementAndGet();
            latch.countDown();
        });

        // Wait for at least 3 firings
        assertTrue(latch.await(2, TimeUnit.SECONDS), "Trigger should fire at least 3 times");

        // Assert
        assertTrue(trigger.isRunning());
        assertTrue(fireCount.get() >= 3, "Fire count should be at least 3, was: " + fireCount.get());

        // Cleanup
        trigger.close();
        assertFalse(trigger.isRunning());
    }

    @Test
    void testStopAndRestart() throws Exception {
        // Arrange
        CountDownLatch latch1 = new CountDownLatch(2);
        CountDownLatch latch2 = new CountDownLatch(2);
        AtomicInteger fireCount = new AtomicInteger(0);

        TimerTrigger trigger = new TimerTrigger(100, TimeUnit.MILLISECONDS);

        // Act - First start
        trigger.start(t -> {
            fireCount.incrementAndGet();
            latch1.countDown();
        });

        assertTrue(latch1.await(1, TimeUnit.SECONDS));
        int countAfterFirstRun = fireCount.get();

        // Stop
        trigger.stop();
        assertFalse(trigger.isRunning());

        // Wait a bit - fire count should not increase
        Thread.sleep(300);
        assertEquals(countAfterFirstRun, fireCount.get(), "Fire count should not increase while stopped");

        // Restart
        trigger.start(t -> {
            fireCount.incrementAndGet();
            latch2.countDown();
        });

        assertTrue(latch2.await(1, TimeUnit.SECONDS));
        assertTrue(fireCount.get() > countAfterFirstRun, "Fire count should increase after restart");

        // Cleanup
        trigger.close();
    }

    @Test
    void testInitialDelay() throws Exception {
        // Arrange
        CountDownLatch latch = new CountDownLatch(1);
        AtomicInteger fireCount = new AtomicInteger(0);
        long startTime = System.currentTimeMillis();

        TimerTrigger trigger = new TimerTrigger(500, 300, TimeUnit.MILLISECONDS);

        // Act
        trigger.start(t -> {
            fireCount.incrementAndGet();
            latch.countDown();
        });

        // Wait for first firing
        assertTrue(latch.await(1, TimeUnit.SECONDS));
        long firstFiringTime = System.currentTimeMillis() - startTime;

        // Assert - first firing should occur after initial delay
        assertTrue(firstFiringTime >= 250, "First firing should occur after initial delay (~300ms), was: " + firstFiringTime + "ms");

        // Cleanup
        trigger.close();
    }

    @Test
    void testMultipleClose() {
        // Arrange
        TimerTrigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);

        // Act - Close multiple times should be safe
        trigger.close();
        trigger.close();
        trigger.close();

        // Assert - Should not throw exception
        assertFalse(trigger.isRunning());
    }

    @Test
    void testStartAfterClose() {
        // Arrange
        TimerTrigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        trigger.close();

        // Act & Assert - Starting after close should fail
        assertThrows(IllegalStateException.class, () -> {
            trigger.start(t -> {});
        });
    }

    @Test
    void testStartTwice() {
        // Arrange
        TimerTrigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);
        trigger.start(t -> {});

        // Act & Assert - Starting twice should fail
        assertThrows(IllegalStateException.class, () -> {
            trigger.start(t -> {});
        });

        // Cleanup
        trigger.close();
    }

    @Test
    void testNullListener() {
        // Arrange
        TimerTrigger trigger = new TimerTrigger(1, TimeUnit.SECONDS);

        // Act & Assert - Null listener should fail
        assertThrows(IllegalArgumentException.class, () -> {
            trigger.start(null);
        });

        // Cleanup
        trigger.close();
    }

    @Test
    void testInvalidInterval() {
        // Act & Assert - Zero interval should fail
        assertThrows(IllegalArgumentException.class, () -> {
            //noinspection resource
            new TimerTrigger(0, TimeUnit.SECONDS);
        });

        // Negative interval should fail
        assertThrows(IllegalArgumentException.class, () -> {
            //noinspection resource
            new TimerTrigger(-1, TimeUnit.SECONDS);
        });
    }

    @Test
    void testInvalidInitialDelay() {
        // Act & Assert - Negative initial delay should fail
        assertThrows(IllegalArgumentException.class, () -> {
            //noinspection resource
            new TimerTrigger(1, -1, TimeUnit.SECONDS);
        });
    }

    @Test
    void testGetType() {
        // Arrange
        TimerTrigger trigger = new TimerTrigger(5, TimeUnit.SECONDS);

        // Act
        String type = trigger.getType();

        // Assert
        assertNotNull(type);
        assertTrue(type.contains("Timer"));
        assertTrue(type.contains("5000")); // 5 seconds = 5000ms

        // Cleanup
        trigger.close();
    }

    @Test
    void testGetters() {
        // Arrange
        TimerTrigger trigger = new TimerTrigger(5, 2, TimeUnit.SECONDS);

        // Act & Assert
        assertEquals(5000, trigger.getIntervalMs());
        assertEquals(2000, trigger.getInitialDelayMs());

        // Cleanup
        trigger.close();
    }
}
