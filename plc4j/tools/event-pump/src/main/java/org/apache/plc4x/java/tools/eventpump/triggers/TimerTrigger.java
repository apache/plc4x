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
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.ThreadPoolExecutor;
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

    // Listener invocations are handed to this executor rather than run on the timer
    // thread. A listener may block (leasing a connection from the connection manager
    // does), and blocking the timer thread delays this trigger's own next tick — and,
    // when a Timer is shared across triggers, every other batch on it as well.
    // The queue holds a single pending firing and the rest are discarded: a listener
    // that outruns the interval must not build an unbounded backlog.
    private final ThreadPoolExecutor dispatcher;

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

        this.dispatcher = new ThreadPoolExecutor(1, 1, 0L, TimeUnit.MILLISECONDS,
            new ArrayBlockingQueue<>(1),
            runnable -> {
                Thread thread = new Thread(runnable, "TimerTrigger-dispatch-" + System.identityHashCode(this));
                thread.setDaemon(true);
                return thread;
            },
            new ThreadPoolExecutor.DiscardPolicy());

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

        // Create the timer task. It only hands the firing off to the dispatcher, so the
        // timer thread itself never runs listener code and never blocks.
        currentTask = new TimerTask() {
            @Override
            public void run() {
                if (!running) {
                    return;
                }
                dispatcher.execute(() -> {
                    if (!running) {
                        return;
                    }
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
                });
            }
        };

        // Fixed delay rather than fixed rate: if a fetch cycle overruns the interval,
        // fixed-rate scheduling fires the missed ticks back-to-back in a burst the moment
        // it catches up, which just piles more pressure on a PLC that is already slow.
        timer.schedule(currentTask, initialDelayMs, intervalMs);
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
        dispatcher.shutdown();

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
