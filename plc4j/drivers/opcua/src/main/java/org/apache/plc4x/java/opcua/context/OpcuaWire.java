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
package org.apache.plc4x.java.opcua.context;

import org.apache.plc4x.java.opcua.readwrite.OpcuaAPU;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Predicate;

/**
 * Minimal abstraction over the OPC UA TCP transport that the {@link Conversation}
 * and {@link SecureChannel} use to talk to the server. The new SPI doesn't have
 * the old {@code ConversationContext.sendRequest(...).expectResponse(...).handle(...)}
 * fluent builder, so this interface gives the conversation layer just the three
 * primitives it actually needs:
 *
 * <ul>
 *   <li>{@link #sendToWire(OpcuaAPU)} — fire and forget</li>
 *   <li>{@link #expect(Predicate, Duration)} — one-shot listener for the next
 *       matching incoming frame; completes the returned future when the frame
 *       arrives or with a {@link java.util.concurrent.TimeoutException} otherwise</li>
 *   <li>{@link #subscribe(Predicate, Consumer)} — multi-fire listener that stays
 *       registered until the returned {@link Subscription} is cancelled. Used
 *       for multi-chunk responses where one logical reply arrives as a sequence
 *       of {@code CONTINUE}+{@code FINAL} chunks</li>
 * </ul>
 */
public interface OpcuaWire {

    /** Sends a frame on the bus; does not wait for a response. */
    void sendToWire(OpcuaAPU apu);

    /**
     * Registers a one-shot listener for the next frame matching {@code predicate}.
     * The returned future completes with that frame, or completes exceptionally
     * with {@link java.util.concurrent.TimeoutException} after {@code timeout}.
     */
    CompletableFuture<OpcuaAPU> expect(Predicate<OpcuaAPU> predicate, Duration timeout);

    /**
     * Registers a multi-fire listener. The handler is invoked for every frame
     * matching {@code predicate} until the returned subscription is cancelled.
     * Used by the conversation layer to accumulate chunks of a multi-chunk
     * response.
     */
    Subscription subscribe(Predicate<OpcuaAPU> predicate, Consumer<OpcuaAPU> handler);

    /** Signals the application layer that the underlying transport has closed. */
    void fireDisconnected();

    /** A single listener registration, cancellable. */
    interface Subscription {
        void unsubscribe();
    }

}
