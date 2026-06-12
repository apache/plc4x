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
package org.apache.plc4x.java.iec608705104.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.Collections;
import java.util.Objects;
import java.util.function.Consumer;

/**
 * Handle for one subscribed IEC 60870-5-104 information object. IEC-104 is
 * push-only on the wire — the controlled station sends ASDUs as state changes
 * occur — so each handle just carries the tag that identifies which incoming
 * objects should be delivered to a registered consumer.
 */
public final class Iec608705104SubscriptionHandle implements PlcSubscriptionHandle {

    private final PlcSubscriber subscriber;
    private final Iec608705104Tag tag;

    public Iec608705104SubscriptionHandle(PlcSubscriber subscriber, Iec608705104Tag tag) {
        this.subscriber = Objects.requireNonNull(subscriber, "subscriber");
        this.tag = Objects.requireNonNull(tag, "tag");
    }

    public Iec608705104Tag getTag() {
        return tag;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return subscriber.registerConsumer(consumer, Collections.singletonList(this));
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Iec608705104SubscriptionHandle other)) {
            return false;
        }
        return Objects.equals(tag, other.tag);
    }

    @Override
    public int hashCode() {
        return Objects.hash(tag);
    }

    @Override
    public String toString() {
        return "Iec608705104SubscriptionHandle[" + tag + "]";
    }

}
