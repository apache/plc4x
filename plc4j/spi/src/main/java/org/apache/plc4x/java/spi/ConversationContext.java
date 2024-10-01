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
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
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

    void fireDiscovered(PlcConnectionConfiguration c);

    SendRequestContext<T> sendRequest(T packet);

    interface SendRequestContext<T> {

        /**
         * names this conversation
         *
         * @param name the name of this conversation
         * @return this
         */
        SendRequestContext<T> name(String name);

        /**
         * defines a response type to be expected with an additional timeout
         *
         * @param clazz   the type of the response
         * @param timeout the timeout
         * @return this
         */
        SendRequestContext<T> expectResponse(Class<T> clazz, Duration timeout);

        /**
         * checks the message using the supplied {@code checker}
         *
         * @param checker function to check the message. Should return true if ok.
         * @return this
         */
        SendRequestContext<T> check(Predicate<T> checker);

        /**
         * final message handle
         *
         * @param packetConsumer consumer used to handle the message
         * @return this
         */
        ContextHandler handle(Consumer<T> packetConsumer);

        CompletableFuture<T> toFuture();

        /**
         * allows to define a timeout handler which then calls {@code packetConsumer}
         *
         * @param packetConsumer the timeout handler
         * @return this
         */
        SendRequestContext<T> onTimeout(Consumer<TimeoutException> packetConsumer);

        /**
         * allows to define an error handler which then calls {@code packetConsumer}
         *
         * @param packetConsumer the error handler
         * @param <E>            the error
         * @return this
         */
        <E extends Throwable> SendRequestContext<T> onError(BiConsumer<T, E> packetConsumer);

        /**
         * unwraps {@code T} and returns {@code R} transformed by {@code unwrapper}
         *
         * @param unwrapper the function used for the transformation
         * @param <R>       the unwrapped type
         * @return this
         */
        <R> SendRequestContext<R> unwrap(Function<T, R> unwrapper);

        /**
         * combines {@link #check } with {@link #unwrap}
         *
         * @param clazz the {@link Class} to be checked
         * @param <R>   type for the {@link Class}
         * @return this
         */
        <R> SendRequestContext<R> only(Class<R> clazz);
    }

    ExpectRequestContext<T> expectRequest(Class<T> clazz, Duration timeout);

    interface ExpectRequestContext<T> {

        ExpectRequestContext<T> name(String name);

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

        void await() throws InterruptedException, ExecutionException;

    }

}
