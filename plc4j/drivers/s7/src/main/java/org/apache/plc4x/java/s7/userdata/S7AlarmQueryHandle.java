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
 * Handle for a one-shot {@code QUERY:ALARM_S}/{@code QUERY:ALARM_8} fetch. Holds the raw
 * byte payload returned by the PLC at subscribe time; when a consumer registers, the
 * connection delivers a single {@link PlcSubscriptionEvent} with that payload and the
 * handle goes inert.
 *
 * <p>Conceptually a Read disguised as a Subscribe — that's how the SPI1 driver exposed it,
 * and we preserve the same user-facing API.
 */
public final class S7AlarmQueryHandle implements PlcSubscriptionHandle {

    private final PlcSubscriber subscriber;
    private final String tagName;
    private final byte[] payload;

    public S7AlarmQueryHandle(PlcSubscriber subscriber, String tagName, byte[] payload) {
        this.subscriber = Objects.requireNonNull(subscriber, "subscriber");
        this.tagName = Objects.requireNonNull(tagName, "tagName");
        this.payload = payload == null ? new byte[0] : payload;
    }

    public String getTagName() {
        return tagName;
    }

    public byte[] getPayload() {
        return payload;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return subscriber.registerConsumer(consumer, Collections.singletonList(this));
    }

    @Override
    public String toString() {
        return "S7AlarmQueryHandle[" + tagName + ", " + payload.length + " bytes]";
    }
}
