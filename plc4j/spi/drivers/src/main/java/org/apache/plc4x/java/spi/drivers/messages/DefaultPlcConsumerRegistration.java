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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.Arrays;
import java.util.List;
import java.util.Objects;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class DefaultPlcConsumerRegistration implements PlcConsumerRegistration {

    private static final AtomicInteger ID_GENERATOR = new AtomicInteger(0);

    private final int consumerId;
    private final PlcSubscriber subscriber;
    private final Consumer<PlcSubscriptionEvent> consumer;
    private final List<PlcSubscriptionHandle> handles;

    public DefaultPlcConsumerRegistration(PlcSubscriber subscriber, Consumer<PlcSubscriptionEvent> consumer, PlcSubscriptionHandle... handles) {
        this.consumerId = ID_GENERATOR.incrementAndGet();
        this.subscriber = subscriber;
        this.consumer = consumer;
        this.handles = Arrays.asList(handles);
    }

    @Override
    public Integer getConsumerId() {
        return consumerId;
    }

    public Consumer<PlcSubscriptionEvent> getConsumer() {
        return consumer;
    }

    @Override
    public List<PlcSubscriptionHandle> getSubscriptionHandles() {
        return handles;
    }

    @Override
    public void unregister() {
        subscriber.unregisterConsumer(this);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof DefaultPlcConsumerRegistration that)) return false;
        return consumerId == that.consumerId;
    }

    @Override
    public int hashCode() {
        return Objects.hash(consumerId);
    }

}
