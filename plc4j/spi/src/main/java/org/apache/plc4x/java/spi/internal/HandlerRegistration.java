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

package org.apache.plc4x.java.spi.internal;

import io.vavr.control.Either;

import java.util.Deque;
import java.util.concurrent.TimeoutException;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

public class HandlerRegistration {
    private final Deque<Either<Function<?, ?>, Predicate<?>>> commands;

    private final Class<?> expectClazz;

    private final Consumer<?> packetConsumer;

    private final Consumer<TimeoutException> onTimeoutConsumer;

    private final BiConsumer<?, ? extends Throwable> errorConsumer;

    public HandlerRegistration(Deque<Either<Function<?, ?>, Predicate<?>>> commands, Class<?> expectClazz, Consumer<?> packetConsumer, Consumer<TimeoutException> onTimeoutConsumer, BiConsumer<?, ? extends Throwable> errorConsumer) {
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

    @Override public String toString() {
        return "HandlerRegistration{" +
            "commands=" + commands +
            ", expectClazz=" + expectClazz +
            ", packetConsumer=" + packetConsumer +
            ", onTimeoutConsumer=" + onTimeoutConsumer +
            ", errorConsumer=" + errorConsumer +
            '}';
    }
}
