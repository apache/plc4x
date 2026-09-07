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

package org.apache.plc4x.java.spi.transports.api;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;

import java.util.function.Consumer;

/**
 * Extension of TransportInstance that supports asynchronous, event-driven I/O.
 * This eliminates the need for polling loops.
 * <p>
 * Implementations register with the underlying I/O mechanism (e.g., NIO Selector)
 * and notify registered listeners when data becomes available, enabling zero-polling,
 * event-driven architectures.
 *
 * @param <T> the configuration type
 */
public interface AsyncTransportInstance<T extends TransportConfiguration> extends TransportInstance<T> {

    /**
     * Registers a listener that will be called when data becomes available.
     * This allows for event-driven architectures without polling.
     * <p>
     * The transport implementation is responsible for monitoring the underlying
     * I/O channel and invoking this listener when data arrives.
     *
     * @param listener callback invoked when data is available
     */
    void registerDataListener(Runnable listener);

    /**
     * Removes the data available listener, stopping event notifications.
     */
    void removeDataListener();

    /**
     * Registers a listener that will be called when the transport is disconnected.
     * <p>
     * This is critical for properly handling connection failures - any pending
     * operations waiting for responses should be completed exceptionally when
     * the transport dies, rather than waiting for a timeout.
     * <p>
     * The listener receives the exception that caused the disconnect, or null
     * if the disconnect was graceful (e.g., the remote side closed normally).
     *
     * @param listener callback invoked when the transport disconnects
     */
    default void registerDisconnectListener(Consumer<Throwable> listener) {
        // Default no-op for backward compatibility
    }

    /**
     * Removes the disconnect listener.
     */
    default void removeDisconnectListener() {
        // Default no-op for backward compatibility
    }

}