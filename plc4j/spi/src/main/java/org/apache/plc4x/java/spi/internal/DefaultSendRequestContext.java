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

public class DefaultSendRequestContext<T> implements ConversationContext.SendRequestContext<T> {

    protected Deque<Either<Function<?, ?>, Predicate<?>>> commands = new LinkedList<>();

    protected final Consumer<HandlerRegistration> finisher;

    private final Object request;

    private final ConversationContext context;

    protected Class<?> expectClazz;

    protected Consumer<?> packetConsumer;

    protected Consumer<TimeoutException> onTimeoutConsumer;

    protected BiConsumer<?, ? extends Throwable> errorConsumer;

    protected Duration timeout = Duration.ofMillis(1000);

    public DefaultSendRequestContext(Consumer<HandlerRegistration> finisher, T request, ConversationContext<T> context) {
        this.finisher = finisher;
        this.request = request;
        this.context = context;
    }

    protected DefaultSendRequestContext(Deque<Either<Function<?, ?>, Predicate<?>>> commands, Duration timeout, Consumer<HandlerRegistration> finisher, Object request, ConversationContext<?> context, Class<?> expectClazz, Consumer<?> packetConsumer, Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer) {
        this.commands = commands;
        this.timeout = timeout;
        this.finisher = finisher;
        this.request = request;
        this.context = context;
        this.expectClazz = expectClazz;
        this.packetConsumer = packetConsumer;
        this.onTimeoutConsumer = onTimeoutConsumer;
        this.errorConsumer = errorConsumer;
    }

    @Override
    public ConversationContext.SendRequestContext<T> expectResponse(Class<T> clazz, Duration timeout) {
        this.timeout = timeout;
        if (expectClazz != null) {
            throw new ConversationContext.PlcWiringException("can't expect class of type " + clazz + " as we already expecting clazz of type " + expectClazz);
        }
        expectClazz = clazz;
        commands.addLast(Either.right(new TypePredicate<>(clazz)));
        return this;
    }

    @Override
    public ConversationContext.SendRequestContext<T> check(Predicate<T> checker) {
        commands.addLast(Either.right(checker));
        return this;
    }

    @Override
    public DefaultContextHandler handle(Consumer<T> packetConsumer) {
        if (this.packetConsumer != null) {
            throw new ConversationContext.PlcWiringException("can't handle multiple consumers");
        }
        this.packetConsumer = packetConsumer;
        final HandlerRegistration registration = new HandlerRegistration(commands, expectClazz, packetConsumer,
            onTimeoutConsumer, errorConsumer, timeout);
        finisher.accept(registration);
        context.sendToWire(request);
        return new DefaultContextHandler(registration::hasHandled, registration::cancel);
    }

    @Override
    public ConversationContext.SendRequestContext<T> onTimeout(Consumer<TimeoutException> onTimeoutConsumer) {
        if (this.onTimeoutConsumer != null) {
            throw new ConversationContext.PlcWiringException("can't handle multiple timeout consumers");
        }
        this.onTimeoutConsumer = onTimeoutConsumer;
        return this;
    }

    @Override
    public <E extends Throwable> ConversationContext.SendRequestContext<T> onError(BiConsumer<T, E> errorConsumer) {
        if (this.errorConsumer != null) {
            throw new ConversationContext.PlcWiringException("can't handle multiple error consumers");
        }
        this.errorConsumer = errorConsumer;
        return this;
    }

    @Override
    public <R> ConversationContext.SendRequestContext<R> unwrap(Function<T, R> unwrapper) {
        if (expectClazz == null) {
            throw new ConversationContext.PlcWiringException("expectResponse must be called before first unwrap");
        }
        if (onTimeoutConsumer == null) {
            onTimeoutConsumer = new NoopTimeoutConsumer();
        }
        commands.addLast(Either.left(unwrapper));
        return new DefaultSendRequestContext<>(commands, timeout, finisher, request, context, expectClazz, packetConsumer, onTimeoutConsumer, errorConsumer);
    }

    @Override
    public <R> ConversationContext.SendRequestContext<R> only(Class<R> clazz) {
        this.check(new TypePredicate<>(clazz));
        return this.unwrap(new CastFunction<>(clazz));
    }

    static class TypePredicate<T, R> implements Predicate<R> {

        private final Class<T> type;

        TypePredicate(Class<T> type) {
            this.type = type;
        }

        @Override
        public boolean test(R value) {
            return type.isInstance(value);
        }
    }

    static class CastFunction<T, R> implements Function<R, T> {

        private final Class<T> type;

        CastFunction(Class<T> type) {
            this.type = type;
        }

        @Override
        public T apply(R value) {
            return type.cast(value);
        }

    }

    static class NoopTimeoutConsumer implements Consumer<TimeoutException> {

        @Override
        public void accept(TimeoutException e) {

        }
    }
}
