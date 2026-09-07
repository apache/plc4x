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
package org.apache.plc4x.java.canopen;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.canopen.tag.CANOpenSubscriptionTag;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.Collections;
import java.util.function.Consumer;

public class CANOpenSubscriptionHandle implements PlcSubscriptionHandle {

    private final PlcSubscriber subscriber;
    private final String name;
    private final CANOpenSubscriptionTag tag;

    public CANOpenSubscriptionHandle(PlcSubscriber subscriber, String name, CANOpenSubscriptionTag tag) {
        this.subscriber = subscriber;
        this.name = name;
        this.tag = tag;
    }

    public boolean matches(CANOpenService service, int identifier) {
        if (tag.getService() != service) {
            return false;
        }
        return tag.isWildcard() || tag.getNodeId() == identifier;
    }

    public String getName() {
        return name;
    }

    public CANOpenSubscriptionTag getTag() {
        return tag;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return subscriber.registerConsumer(consumer, Collections.singletonList(this));
    }

    @Override
    public String toString() {
        return "CANopenSubscriptionHandle [service=" + tag.getService()
            + ", node=" + intAndHex(tag.getNodeId())
            + ", cob=" + intAndHex(tag.getService().getMin() + tag.getNodeId()) + "]";
    }

    private static String intAndHex(int val) {
        return val + "(0x" + Integer.toHexString(val) + ")";
    }

}
