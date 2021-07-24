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
package org.apache.plc4x.java.can.adapter.conversation;

import java.time.Duration;
import java.util.concurrent.TimeoutException;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;
import org.apache.plc4x.java.spi.ConversationContext.ContextHandler;
import org.apache.plc4x.java.spi.ConversationContext.SendRequestContext;
import org.apache.plc4x.java.transport.can.CANTransport.FrameHandler;
import org.apache.plc4x.java.transport.can.FrameData;

public class SendRequestContextWrapper<C, T> implements SendRequestContext<T> {

    private final SendRequestContext<C> delegate;
    private final Class<C> wireType;
    private final Function<C, FrameData> adapter;
    private final FrameHandler<C, T> frameHandler;

    public SendRequestContextWrapper(SendRequestContext<C> delegate, Class<C> wireType, Function<C, FrameData> adapter, FrameHandler<C, T> frameHandler) {
        this.delegate = delegate;
        this.wireType = wireType;
        this.adapter = adapter;
        this.frameHandler = frameHandler;
    }

    @Override
    public SendRequestContext<T> expectResponse(Class<T> clazz, Duration timeout) {
        DeferredErrorHandler<C, ?> errorHandler = new DeferredErrorHandler<>(null);
        DeferredTimeoutHandler<?> timeoutHandler = new DeferredTimeoutHandler<>(null);
        return new ResolvedSendRequestContextWrapper<>(
            delegate.onError(errorHandler)
                .onTimeout(timeoutHandler)
                .expectResponse(wireType, timeout)
                .unwrap(adapter)
                .unwrap(frameHandler::fromCAN)
                .check(f -> clazz.isAssignableFrom(f.getClass())
            ), errorHandler, timeoutHandler
        );
    }

    @Override
    public SendRequestContext<T> check(Predicate<T> checker) {
        throw new IllegalStateException("Not implemented");
    }

    @Override
    public ContextHandler handle(Consumer<T> packetConsumer) {
        throw new IllegalStateException("Not implemented");
    }

    @Override
    public SendRequestContext<T> onTimeout(Consumer<TimeoutException> packetConsumer) {
        throw new IllegalStateException("Not implemented");
    }

    @Override
    public <E extends Throwable> SendRequestContext<T> onError(BiConsumer<T, E> packetConsumer) {
        throw new IllegalStateException("Not implemented");
    }

    @Override
    public <R> SendRequestContext<R> unwrap(Function<T, R> unwrapper) {
        throw new IllegalStateException("Not implemented");
    }

    @Override
    public <R> SendRequestContext<R> only(Class<R> clazz) {
        throw new IllegalStateException("Not implemented");
    }
}
