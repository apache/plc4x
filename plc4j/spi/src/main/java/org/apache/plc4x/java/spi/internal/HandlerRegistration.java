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
package org.apache.plc4x.java.spi.internal;

import io.vavr.control.Either;

import java.time.Duration;
import java.time.Instant;
import java.util.Deque;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

public class HandlerRegistration {

    private static int counter = 0;

    private final int id = counter++;

    private final Deque<Either<Function<?, ?>, Predicate<?>>> commands;

    private final Class<?> expectClazz;

    private final Consumer<?> packetConsumer;

    private final Consumer<TimeoutException> onTimeoutConsumer;

    private final BiConsumer<?, ? extends Throwable> errorConsumer;
    private final Duration timeout;
    private final Instant timeoutAt;

    private volatile boolean cancelled = false;
    private volatile boolean handled = false;

    public HandlerRegistration(Deque<Either<Function<?, ?>, Predicate<?>>> commands, Class<?> expectClazz, Consumer<?> packetConsumer, Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer, Duration timeout) {
        this.commands = commands;
        this.expectClazz = expectClazz;
        this.packetConsumer = packetConsumer;
        this.onTimeoutConsumer = onTimeoutConsumer;
        this.errorConsumer = errorConsumer;
        this.timeout = timeout;
        this.timeoutAt = Instant.now().plus(timeout);
    }

    public Deque<Either<Function<?, ?>, Predicate<?>>> getCommands() {
        return commands;
    }

    public Class<?> getExpectClazz() {
        return expectClazz;
    }

    public Consumer<?> getPacketConsumer() {
        return packetConsumer;
    }

    public Consumer<TimeoutException> getOnTimeoutConsumer() {
        return onTimeoutConsumer;
    }

    public BiConsumer<?, ? extends Throwable> getErrorConsumer() {
        return errorConsumer;
    }

    public Duration getTimeout() {
        return timeout;
    }

    public Instant getTimeoutAt() {
        return timeoutAt;
    }

    public void cancel() {
        this.cancelled = true;
    }

    public boolean isCancelled() {
        return this.cancelled;
    }

    public void confirmHandled() {
        this.handled = true;
    }

    public boolean hasHandled() {
        return this.handled;
    }

    @Override
    public String toString() {
        return "HandlerRegistration#" + id;
    }
}
