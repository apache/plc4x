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
package org.apache.plc4x.java.transport.serial;

/**
 * Enforces a minimum silence gap before writes on a serial line, measured
 * from the later of the last write end and the last observed read activity
 * (Modbus RTU t3.5-style semantics). Mirrors plc4go's pacer.
 *
 * A zero (or negative) delay disables pacing. The pacer only tracks
 * timestamps — it does not serialize concurrent writers; callers must hold
 * their own write lock around awaitTurn() + write for the gap guarantee to
 * hold.
 */
final class WritePacer {

    private final long delayMillis;
    private final long delayNanos;
    private final Object lock = new Object();
    private long lastActivityNanos;

    WritePacer(long delayMillis) {
        this.delayMillis = delayMillis;
        this.delayNanos = delayMillis * 1_000_000L;
        // Initialized so a fresh pacer never waits: the first awaitTurn()
        // sees "last activity" as already delayNanos in the past.
        this.lastActivityNanos = System.nanoTime() - delayNanos;
    }

    /**
     * Records line activity (a completed write or received data). Monotone:
     * never moves the timestamp backwards. Uses {@link System#nanoTime()},
     * which is immune to wall-clock adjustments (NTP, DST, manual changes)
     * that would otherwise corrupt the gap computation.
     */
    void noteActivity() {
        if (delayMillis <= 0) {
            return;
        }
        synchronized (lock) {
            long now = System.nanoTime();
            if (now - lastActivityNanos > 0) {
                lastActivityNanos = now;
            }
        }
    }

    /**
     * Sleeps until the configured gap since the last activity has elapsed.
     * Recomputes after every sleep, so activity arriving during the wait
     * extends the gap.
     */
    void awaitTurn() throws InterruptedException {
        if (delayMillis <= 0) {
            return;
        }
        while (true) {
            long ready;
            synchronized (lock) {
                ready = lastActivityNanos + delayNanos;
            }
            long waitNanos = ready - System.nanoTime();
            if (waitNanos <= 0) {
                return;
            }
            Thread.sleep(waitNanos / 1_000_000L, (int) (waitNanos % 1_000_000L));
        }
    }
}
