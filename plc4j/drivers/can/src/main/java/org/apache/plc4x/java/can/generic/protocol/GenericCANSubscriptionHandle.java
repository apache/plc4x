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
package org.apache.plc4x.java.can.generic.protocol;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.can.generic.tag.GenericCANTag;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.function.Consumer;

public class GenericCANSubscriptionHandle implements PlcSubscriptionHandle {

    private final PlcSubscriber subscriber;
    private final Map<String, GenericCANTag> tags = new LinkedHashMap<>();
    private final Integer nodeId;

    public GenericCANSubscriptionHandle(PlcSubscriber subscriber, Integer nodeId) {
        this.subscriber = subscriber;
        this.nodeId = nodeId;
    }

    public boolean matches(int identifier) {
        return nodeId == identifier;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return subscriber.registerConsumer(consumer, Collections.singletonList(this));
    }

    public void add(String name, GenericCANTag tag) {
        tags.put(name, tag);
    }

    public Map<String, GenericCANTag> getTags() {
        return Collections.unmodifiableMap(tags);
    }

    @Override
    public String toString() {
        return "GenericCANSubscriptionHandle[node=" + nodeId + " (0x" + Integer.toHexString(nodeId) + ")]";
    }
}
