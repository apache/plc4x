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
package org.apache.plc4x.java.s7.userdata;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.Collections;
import java.util.Objects;
import java.util.function.Consumer;

/**
 * Handle for an active S7Comm alarm subscription. The handle identifies the registered
 * tag name within a single subscribe call; alarm indications are dispatched to all
 * consumers whose registrations include any handle from this connection.
 */
public final class S7AlarmSubscriptionHandle implements PlcSubscriptionHandle {

    private final PlcSubscriber subscriber;
    private final String tagName;

    public S7AlarmSubscriptionHandle(PlcSubscriber subscriber, String tagName) {
        this.subscriber = Objects.requireNonNull(subscriber, "subscriber");
        this.tagName = Objects.requireNonNull(tagName, "tagName");
    }

    public String getTagName() {
        return tagName;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return subscriber.registerConsumer(consumer, Collections.singletonList(this));
    }

    @Override
    public String toString() {
        return "S7AlarmSubscriptionHandle[" + tagName + "]";
    }
}
