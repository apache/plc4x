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
package org.apache.plc4x.java.canopen.conversation;

import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.function.Predicate;

/**
 * Minimal request/response abstraction used by the multi-turn SDO conversations.
 *
 * <p>The new SPI's {@code ConnectionBase} doesn't ship a fluent ConversationContext
 * like the old SPI did. This interface gives the SDO upload/download conversations
 * just two primitives:
 *
 * <ul>
 *   <li>{@link #sendToWire(CANOpenFrame)} — fire-and-forget a frame.</li>
 *   <li>{@link #expect(Predicate, Duration)} — return a future that completes
 *       with the next incoming frame matching {@code predicate}, or completes
 *       exceptionally with {@link java.util.concurrent.TimeoutException} after
 *       {@code timeout}.</li>
 * </ul>
 *
 * <p>Multi-step protocols become straightforward CompletableFuture chains: send,
 * expect, decide, send again, expect again.</p>
 */
public interface CANConversation {

    /** Sends a frame on the bus; does not wait for a response. */
    void sendToWire(CANOpenFrame frame);

    /**
     * Registers a one-shot listener for the next frame matching {@code predicate}.
     * The returned future completes with that frame, or completes exceptionally
     * with {@link java.util.concurrent.TimeoutException} after {@code timeout}.
     */
    CompletableFuture<CANOpenFrame> expect(Predicate<CANOpenFrame> predicate, Duration timeout);

}
