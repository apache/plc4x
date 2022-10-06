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
package org.apache.plc4x.java.spi;

import io.netty.channel.Channel;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.configuration.Configuration;

import java.time.Duration;
import java.util.concurrent.TimeoutException;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

public interface ConversationContext<T> {

    PlcAuthentication getAuthentication();

    Channel getChannel();

    boolean isPassive();

    void sendToWire(T msg);

    void fireConnected();

    void fireDisconnected();

    void fireDiscovered(Configuration c);

    SendRequestContext<T> sendRequest(T packet);

    interface SendRequestContext<T> {

        SendRequestContext<T> expectResponse(Class<T> clazz, Duration timeout);

        SendRequestContext<T> check(Predicate<T> checker);

        ContextHandler handle(Consumer<T> packetConsumer);

        SendRequestContext<T> onTimeout(Consumer<TimeoutException> packetConsumer);

        <E extends Throwable> SendRequestContext<T> onError(BiConsumer<T, E> packetConsumer);

        <R> SendRequestContext<R> unwrap(Function<T, R> unwrapper);

        <R> SendRequestContext<R> only(Class<R> clazz);
    }

    ExpectRequestContext<T> expectRequest(Class<T> clazz, Duration timeout);

    interface ExpectRequestContext<T> {

        ExpectRequestContext<T> check(Predicate<T> checker);

        ContextHandler handle(Consumer<T> packetConsumer);

        ExpectRequestContext<T> onTimeout(Consumer<TimeoutException> packetConsumer);

        <E extends Throwable> ExpectRequestContext<T> onError(BiConsumer<T, E> packetConsumer);

        <R> ExpectRequestContext<R> unwrap(Function<T, R> unwrapper);
    }

    class PlcWiringException extends PlcRuntimeException {
        public PlcWiringException(String message) {
            super(message);
        }
    }

    interface ContextHandler {

        boolean isDone();

        void cancel();

    }

}
