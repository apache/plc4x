/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.api.connection;

import org.apache.plc4x.java.api.messages.PlcNotification;
import org.apache.plc4x.java.api.model.Address;

import java.util.function.Consumer;

/**
 * Interface implemented by all PlcConnections that are able to receive notifications from remote resources.
 */
public interface PlcSubscriber {

    /**
     * Subscribes a {@code consumer} to a {@code address} parsing values as {@code dataType}.
     * {@code consumer} and {@code address} are used as unique identification.
     *
     * @param consumer to be subscribed.
     * @param address  to be read.
     * @param dataType to be decoded.
     */
    void subscribe(Consumer<PlcNotification<?>> consumer, Address address, Class<?> dataType);


    /**
     * Unsubscribes a {@code consumer}.
     * {@code consumer} and {@code address} are used as unique identification.
     *
     * @param consumer to be unsubscribed.
     */
    void unsubscribe(Consumer<PlcNotification<?>> consumer, Address address);
}
