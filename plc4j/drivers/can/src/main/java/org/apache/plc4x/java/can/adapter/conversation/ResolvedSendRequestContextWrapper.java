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

public class ResolvedSendRequestContextWrapper<T> implements SendRequestContext<T> {

    private final SendRequestContext<T> delegate;
    private final DeferredErrorHandler<?, ?> errorHandler;
    private final DeferredTimeoutHandler<?> timeoutHandler;

    public ResolvedSendRequestContextWrapper(SendRequestContext<T> delegate, DeferredErrorHandler<?, ?> errorHandler,
        DeferredTimeoutHandler<?> timeoutHandler) {
        this.delegate = delegate;
        this.errorHandler = errorHandler;
        this.timeoutHandler = timeoutHandler;
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
        timeoutHandler.setHandler(packetConsumer);
        return this;
    }

    @Override
    public <E extends Throwable> SendRequestContext<T> onError(BiConsumer<T, E> packetConsumer) {
        errorHandler.setHandler((BiConsumer) packetConsumer);
        return this;
    }

    @Override
    public <R> SendRequestContext<R> unwrap(Function<T, R> unwrapper) {
        return new ResolvedSendRequestContextWrapper<>(delegate.unwrap(unwrapper), errorHandler, timeoutHandler);
    }

    @Override
    public <R> SendRequestContext<R> only(Class<R> clazz) {
        return new ResolvedSendRequestContextWrapper<>(delegate.only(clazz), errorHandler, timeoutHandler);
    }
}
