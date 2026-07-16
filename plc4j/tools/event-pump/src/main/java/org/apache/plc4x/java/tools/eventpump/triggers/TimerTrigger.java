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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.TimeUnit;

/**
 * A trigger that fires at fixed time intervals.
 * <p>
 * This trigger uses a {@link Timer} to schedule periodic execution.
 * The interval can be configured with any time unit (milliseconds, seconds, minutes, etc.).
 * <p>
 * Example usage:
 * <pre>
 * // Trigger every 5 seconds
 * Trigger trigger = new TimerTrigger(5, TimeUnit.SECONDS);
 * trigger.start(t -> System.out.println("Trigger fired!"));
 * </pre>
 */
public class TimerTrigger implements Trigger {

    private static final Logger LOGGER = LoggerFactory.getLogger(TimerTrigger.class);

    private final long intervalMs;
    private final long initialDelayMs;
    private final Timer timer;
    private final boolean ownTimer;

    private volatile TriggerListener listener;
    private volatile TimerTask currentTask;
    private volatile boolean running = false;
    private volatile boolean closed = false;

    /**
     * Create a timer trigger with the specified interval and no initial delay.
     *
     * @param interval The interval between trigger firings
     * @param unit The time unit of the interval
     */
    public TimerTrigger(long interval, TimeUnit unit) {
        this(interval, 0, unit);
    }

    /**
     * Create a timer trigger with the specified interval and initial delay.
     *
     * @param interval The interval between trigger firings
     * @param initialDelay The initial delay before the first firing
     * @param unit The time unit of the interval and delay
     */
    public TimerTrigger(long interval, long initialDelay, TimeUnit unit) {
        this(interval, initialDelay, unit, null);
    }

    /**
     * Create a timer trigger with a custom timer.
     * Useful for sharing a timer across multiple triggers.
     *
     * @param interval The interval between trigger firings
     * @param initialDelay The initial delay before the first firing
     * @param unit The time unit of the interval and delay
     * @param sharedTimer A shared timer (if null, a new daemon timer will be created)
     */
    public TimerTrigger(long interval, long initialDelay, TimeUnit unit, Timer sharedTimer) {
        if (interval <= 0) {
            throw new IllegalArgumentException("Interval must be positive");
        }
        if (initialDelay < 0) {
            throw new IllegalArgumentException("Initial delay cannot be negative");
        }

        this.intervalMs = unit.toMillis(interval);
        this.initialDelayMs = unit.toMillis(initialDelay);

        if (sharedTimer != null) {
            this.timer = sharedTimer;
            this.ownTimer = false;
        } else {
            this.timer = new Timer("TimerTrigger-" + System.identityHashCode(this), true);
            this.ownTimer = true;
        }

        LOGGER.debug("Created TimerTrigger with interval={}ms, initialDelay={}ms", intervalMs, initialDelayMs);
    }

    @Override
    public void start(TriggerListener listener) {
        if (listener == null) {
            throw new IllegalArgumentException("Listener cannot be null");
        }
        if (closed) {
            throw new IllegalStateException("Trigger is closed");
        }
        if (running) {
            throw new IllegalStateException("Trigger is already running");
        }

        this.listener = listener;

        // Create the timer task
        currentTask = new TimerTask() {
            @Override
            public void run() {
                try {
                    LOGGER.trace("Timer trigger firing");
                    TimerTrigger.this.listener.onTrigger(TimerTrigger.this);
                } catch (Exception e) {
                    if (LOGGER.isTraceEnabled()) {
                        LOGGER.error("Error in trigger listener", e);
                    } else {
                        LOGGER.error("Error in trigger listener: {}", e.getMessage());
                    }
                }
            }
        };

        // Schedule the task
        timer.scheduleAtFixedRate(currentTask, initialDelayMs, intervalMs);
        running = true;

        LOGGER.debug("Started TimerTrigger");
    }

    @Override
    public void stop() {
        if (!running) {
            return;
        }

        if (currentTask != null) {
            currentTask.cancel();
            currentTask = null;
        }

        running = false;
        LOGGER.debug("Stopped TimerTrigger");
    }

    @Override
    public boolean isRunning() {
        return running;
    }

    @Override
    public String getType() {
        return "Timer (interval=" + intervalMs + "ms)";
    }

    @Override
    public void close() {
        if (closed) {
            return;
        }

        stop();

        if (ownTimer) {
            timer.cancel();
        }

        closed = true;
        LOGGER.debug("Closed TimerTrigger");
    }

    /**
     * Get the interval in milliseconds.
     *
     * @return The interval in milliseconds
     */
    public long getIntervalMs() {
        return intervalMs;
    }

    /**
     * Get the initial delay in milliseconds.
     *
     * @return The initial delay in milliseconds
     */
    public long getInitialDelayMs() {
        return initialDelayMs;
    }
}
