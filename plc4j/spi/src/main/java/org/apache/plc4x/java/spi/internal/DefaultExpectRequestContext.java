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
import org.apache.plc4x.java.spi.ConversationContext;

import java.time.Duration;
import java.time.Instant;
import java.util.Deque;
import java.util.LinkedList;
import java.util.concurrent.TimeoutException;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

public class DefaultExpectRequestContext<T> implements ConversationContext.ExpectRequestContext<T> {

    protected Deque<Either<Function<?, ?>, Predicate<?>>> commands = new LinkedList<>();

    protected final Consumer<HandlerRegistration> finisher;

    private final ConversationContext context;

    protected final Class<?> expectClazz;

    protected Consumer<?> packetConsumer;

    protected Consumer<TimeoutException> onTimeoutConsumer;

    protected BiConsumer<?, ? extends Throwable> errorConsumer;

    protected final Duration timeout;
    private HandlerRegistration registration;

    public DefaultExpectRequestContext(Consumer<HandlerRegistration> finisher, Class<T> expectClazz, Duration timeout, ConversationContext context) {
        this.finisher = finisher;
        this.expectClazz = expectClazz;
        this.timeout = timeout;
        this.context = context;
    }

    protected DefaultExpectRequestContext(Deque<Either<Function<?, ?>, Predicate<?>>> commands, Duration timeout, Consumer<HandlerRegistration> finisher, ConversationContext<?> context, Class<?> expectClazz, Consumer<?> packetConsumer, Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer) {
        this.commands = commands;
        this.timeout = timeout;
        this.finisher = finisher;
        this.context = context;
        this.expectClazz = expectClazz;
        this.packetConsumer = packetConsumer;
        this.onTimeoutConsumer = onTimeoutConsumer;
        this.errorConsumer = errorConsumer;
    }

    @Override
    public ConversationContext.ExpectRequestContext<T> check(Predicate<T> checker) {
        commands.addLast(Either.right(checker));
        return this;
    }

    @Override
    public ConversationContext.ContextHandler handle(Consumer<T> packetConsumer) {
        if (this.packetConsumer != null) {
            throw new ConversationContext.PlcWiringException("can't handle multiple consumers");
        }
        this.packetConsumer = packetConsumer;
        registration = new HandlerRegistration(commands, expectClazz, packetConsumer, onTimeoutConsumer, errorConsumer, timeout);
        finisher.accept(registration);
        return new DefaultContextHandler(registration::hasHandled, registration::cancel);
    }

    @Override
    public ConversationContext.ExpectRequestContext<T> onTimeout(Consumer<TimeoutException> onTimeoutConsumer) {
        if (this.onTimeoutConsumer != null) {
            throw new ConversationContext.PlcWiringException("can't handle multiple timeout consumers");
        }
        this.onTimeoutConsumer = onTimeoutConsumer;
        return this;
    }

    @Override
    public <E extends Throwable> ConversationContext.ExpectRequestContext<T> onError(BiConsumer<T, E> errorConsumer) {
        if (this.errorConsumer != null) {
            throw new ConversationContext.PlcWiringException("can't handle multiple error consumers");
        }
        this.errorConsumer = errorConsumer;
        return this;
    }

    @Override
    public <R> ConversationContext.ExpectRequestContext<R> unwrap(Function<T, R> unwrapper) {
        if (expectClazz == null) {
            throw new ConversationContext.PlcWiringException("expectResponse must be called before first unwrap");
        }
        if (onTimeoutConsumer == null) {
            onTimeoutConsumer = e -> {
                // NOOP
            };
        }
        commands.addLast(Either.left(unwrapper));
        return new DefaultExpectRequestContext<>(commands, timeout, finisher, context, expectClazz, packetConsumer, onTimeoutConsumer, errorConsumer);
    }

}
