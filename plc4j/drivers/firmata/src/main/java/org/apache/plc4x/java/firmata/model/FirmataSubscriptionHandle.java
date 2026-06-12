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
package org.apache.plc4x.java.firmata.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.firmata.tag.FirmataTag;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.Collections;
import java.util.Objects;
import java.util.function.Consumer;

/**
 * Handle for one subscribed Firmata tag (digital or analog). Firmata is push-only:
 * the board emits unsolicited update messages whenever a configured pin changes,
 * so each handle only carries the metadata the connection needs to filter incoming
 * events against the active subscriptions.
 */
public final class FirmataSubscriptionHandle implements PlcSubscriptionHandle {

    private final PlcSubscriber subscriber;
    private final String name;
    private final FirmataTag tag;

    public FirmataSubscriptionHandle(PlcSubscriber subscriber, String name, FirmataTag tag) {
        this.subscriber = Objects.requireNonNull(subscriber, "subscriber");
        this.name = Objects.requireNonNull(name, "name");
        this.tag = Objects.requireNonNull(tag, "tag");
    }

    public String getName() {
        return name;
    }

    public FirmataTag getTag() {
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
        if (!(o instanceof FirmataSubscriptionHandle other)) {
            return false;
        }
        return Objects.equals(tag, other.tag) && Objects.equals(name, other.name);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name, tag);
    }

    @Override
    public String toString() {
        return "FirmataSubscriptionHandle[" + name + "=" + tag + "]";
    }

}
