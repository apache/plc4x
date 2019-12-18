/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.java.spi.internal;

import io.vavr.control.Either;
import org.apache.plc4x.java.spi.ConversationContext;

import java.time.Duration;
import java.util.Deque;
import java.util.LinkedList;
import java.util.concurrent.TimeoutException;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

public class DefaultSendRequestContext<T> implements ConversationContext.SendRequestContext<T> {

    protected Deque<Either<Function<?, ?>, Predicate<?>>> commands = new LinkedList<>();

    protected final Consumer<FinalContext> finisher;

    protected Class<?> expectClazz;

    protected Consumer<?> packetConsumer;

    protected Consumer<TimeoutException> onTimeoutConsumer;

    protected BiConsumer<?, ? extends Throwable> errorConsumer;

    public DefaultSendRequestContext(Consumer<FinalContext> finisher) {
        this.finisher = finisher;
    }

    protected DefaultSendRequestContext(Deque<Either<Function<?, ?>, Predicate<?>>> commands, Consumer<FinalContext> finisher, Class<?> expectClazz, Consumer<?> packetConsumer, Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer) {
        this.commands = commands;
        this.finisher = finisher;
        this.expectClazz = expectClazz;
        this.packetConsumer = packetConsumer;
        this.onTimeoutConsumer = onTimeoutConsumer;
        this.errorConsumer = errorConsumer;
    }

    @Override
    public ConversationContext.SendRequestContext<T> expectResponse(Class<T> clazz, Duration timeout) {
        if (expectClazz != null) {
            throw new ConversationContext.PlcWiringException("can't expect class of type " + clazz + " as we already expecting clazz of type " + expectClazz);
        }
        expectClazz = clazz;
        commands.addLast(Either.right(clazz::isInstance));
        return this;
    }

    @Override
    public ConversationContext.SendRequestContext<T> check(Predicate<T> checker) {
        commands.addLast(Either.right(checker));
        return this;
    }

    @Override
    public void handle(Consumer<T> packetConsumer) {
        if (this.packetConsumer != null) {
            throw new ConversationContext.PlcWiringException("can't handle multiple consumers");
        }
        this.packetConsumer = packetConsumer;
        finisher.accept(new FinalContext(commands, expectClazz, packetConsumer, onTimeoutConsumer, errorConsumer));
    }

    @Override
    public <E extends Throwable> ConversationContext.SendRequestContext<T> onTimeout(Consumer<TimeoutException> onTimeoutConsumer) {
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
            throw new ConversationContext.PlcWiringException("onTimeout must be called before first unwrap");
        }
        return new DefaultSendRequestContext<>(commands, finisher, expectClazz, packetConsumer, onTimeoutConsumer, errorConsumer);
    }

    public static class FinalContext {
        private final Deque<Either<Function<?, ?>, Predicate<?>>> commands;

        private final Class<?> expectClazz;

        private final Consumer<?> packetConsumer;

        private final Consumer<TimeoutException> onTimeoutConsumer;

        private final BiConsumer<?, ? extends Throwable> errorConsumer;

        public FinalContext(Deque<Either<Function<?, ?>, Predicate<?>>> commands, Class<?> expectClazz, Consumer<?> packetConsumer, Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer) {
            this.commands = commands;
            this.expectClazz = expectClazz;
            this.packetConsumer = packetConsumer;
            this.onTimeoutConsumer = onTimeoutConsumer;
            this.errorConsumer = errorConsumer;
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
    }

}
