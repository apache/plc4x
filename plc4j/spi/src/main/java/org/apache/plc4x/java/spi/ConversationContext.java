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

package org.apache.plc4x.java.spi;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.time.Duration;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;

public interface ConversationContext<T> {

    void sendToWire(T msg);

    void fireConnected();

    SendRequestContext<T> sendRequest(T packet);

    interface SendRequestContext<T> {

        SendRequestContext<T> expectResponse(Class<T> clazz, Duration timeout);

        SendRequestContext<T> check(Function<T, Boolean> checker);

        SendRequestContext<T> handle(Consumer<T> packetConsumer);

        <E extends Throwable> SendRequestContext<T> onTimeout(BiConsumer<T, E> packetConsumer);

        <E extends Throwable> SendRequestContext<T> onError(BiConsumer<T, E> packetConsumer);

        SendRequestContext<T> onSuccess(Consumer<T> packetConsumer);

        <R> SendRequestContext<R> unwrap(Function<T, R> unwrapper);

        <R> SendRequestContext<R> unwrap(Class<R> clazz, Function<T, R> unwrapper);

        void finish();
    }

    class PlcCompletionException extends PlcRuntimeException {

        public PlcCompletionException(String message) {
            super(message);
        }

        public PlcCompletionException(String message, Throwable cause) {
            super(message, cause);
        }

        public PlcCompletionException(Throwable cause) {
            super(cause);
        }

        public PlcCompletionException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
            super(message, cause, enableSuppression, writableStackTrace);
        }
    }
}
