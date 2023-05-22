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
package org.apache.plc4x.java.can.adapter.conversation;

import java.time.Duration;
import java.util.concurrent.TimeoutException;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;
import org.apache.plc4x.java.spi.ConversationContext.ContextHandler;
import org.apache.plc4x.java.spi.ConversationContext.SendRequestContext;

public class DeferredRequestContextWrapper<T> implements SendRequestContext<T> {

    private final SendRequestContext<T> delegate;
    private final Function<SendRequestContext<?>, SendRequestContext<?>> completer;
    protected final Consumer<TimeoutException> onTimeoutConsumer;
    protected final BiConsumer<?, ? extends Throwable> errorConsumer;

    public DeferredRequestContextWrapper(SendRequestContext<T> delegate, Function<SendRequestContext<?>, SendRequestContext<?>> completer,
        Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer) {
        this.delegate = delegate;
        this.completer = completer;
        this.onTimeoutConsumer = onTimeoutConsumer;
        this.errorConsumer = errorConsumer;
    }

    @Override
    public SendRequestContext<T> expectResponse(Class<T> clazz, Duration timeout) {
        throw new IllegalStateException("Response type is already specified");
    }

    @Override
    public SendRequestContext<T> check(Predicate<T> checker) {
        delegate.check(checker);
        return this;
    }

    @Override
    public ContextHandler handle(Consumer<T> packetConsumer) {
        return delegate.handle(packetConsumer);
    }

    @Override
    public SendRequestContext<T> onTimeout(Consumer<TimeoutException> packetConsumer) {
        return new DeferredRequestContextWrapper<>(delegate, completer, packetConsumer, errorConsumer);
    }

    @Override
    public <E extends Throwable> SendRequestContext<T> onError(BiConsumer<T, E> packetConsumer) {
        return new DeferredRequestContextWrapper<>(delegate, completer, onTimeoutConsumer, errorConsumer);
    }

    @Override
    public <R> SendRequestContext<R> unwrap(Function<T, R> unwrapper) {
        return resolve(new DeferredRequestContextWrapper<>(delegate.unwrap(unwrapper), completer, onTimeoutConsumer, errorConsumer));
    }

    private <R> SendRequestContext<R> resolve(DeferredRequestContextWrapper<R> contextWrapper) {
        if (contextWrapper.errorConsumer != null && contextWrapper.onTimeoutConsumer != null) {
            return (SendRequestContext<R>) contextWrapper.completer.apply(contextWrapper);
        }
        return contextWrapper;
    }

    @Override
    public <R> SendRequestContext<R> only(Class<R> clazz) {
        return resolve(new DeferredRequestContextWrapper<>(delegate.only(clazz), completer, onTimeoutConsumer, errorConsumer));
    }
}
