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
package org.apache.plc4x.java.spi.internal;

import io.vavr.control.Either;

import java.time.Duration;
import java.util.Deque;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

public class HandlerRegistration implements Future<Void> {

    private static int counter = 0;

    private final int id = counter++;

    private final String name;

    private final Deque<Either<Function<?, ?>, Predicate<?>>> commands;

    private final Class<?> expectClazz;

    private final Consumer<?> packetConsumer;

    private final Consumer<TimeoutException> onTimeoutConsumer;

    private final BiConsumer<?, ? extends Throwable> errorConsumer;
    private final Runnable onHandled;
    private final Runnable onError;
    private final Runnable onCancelled;
    private final Duration timeout;

    private final CompletableFuture<Void> handled = new CompletableFuture<>();

    public HandlerRegistration(String name, Deque<Either<Function<?, ?>, Predicate<?>>> commands, Class<?> expectClazz, Consumer<?> packetConsumer, Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer, Duration timeout) {
        this(
            name,
            commands,
            expectClazz,
            packetConsumer,
            onTimeoutConsumer,
            errorConsumer,
            () -> {},
            () -> {},
            () -> {},
            timeout
        );
    }

    public HandlerRegistration(String name, Deque<Either<Function<?, ?>, Predicate<?>>> commands, Class<?> expectClazz, Consumer<?> packetConsumer, Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer, Runnable onHandled, Runnable onError, Runnable onCancelled, Duration timeout) {
        this.name = name;
        this.commands = commands;
        this.expectClazz = expectClazz;
        this.packetConsumer = packetConsumer;
        this.onTimeoutConsumer = onTimeoutConsumer;
        this.errorConsumer = errorConsumer;
        this.onHandled = onHandled;
        this.onError = onError;
        this.onCancelled = onCancelled;
        this.timeout = timeout;
    }

    public String getName() {
        return name;
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

    public void cancel() {
        handled.cancel(true);
        onCancelled.run();
    }

    @Override
    public boolean cancel(boolean ignored) {
        if (isCancelled()) {
            return false;
        } else {
            cancel();
            return true;
        }
    }

    public boolean isCancelled() {
        return handled.isCancelled();
    }

    @Override
    public boolean isDone() {
        return hasHandled();
    }

    @Override
    public Void get() throws InterruptedException, ExecutionException {
        return handled.get();
    }

    @Override
    public Void get(long amount, TimeUnit timeUnit) throws InterruptedException, ExecutionException, TimeoutException {
        return handled.get(amount, timeUnit);
    }

    public void confirmHandled() {
        confirmCompleted();
        this.onHandled.run();
    }

    public void confirmError() {
        confirmCompleted();
        this.onError.run();
    }

    public void confirmCompleted() {
        this.handled.complete(null);
    }

    public boolean hasHandled() {
        return this.handled.isDone();
    }

    @Override
    public String toString() {
        return "HandlerRegistration#" + id;
    }

}
