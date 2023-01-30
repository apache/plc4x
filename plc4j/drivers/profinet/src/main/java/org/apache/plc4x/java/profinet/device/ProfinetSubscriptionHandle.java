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

package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;
import java.util.*;
import java.util.function.Consumer;

public class ProfinetSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private final PlcSubscriptionTag tag;
    private final String address;
    private PlcValue lastValue;

    public ProfinetSubscriptionHandle(PlcSubscriber plcSubscriber, String address, PlcSubscriptionTag tag) {
        super(plcSubscriber);
        this.address = address;
        this.tag = tag;
    }

    public String getTag() {
        return address;
    }

    public String getAddressString() {
        return tag.getAddressString();
    }

    public PlcSubscriptionType getSubscriptionType() {
        return tag.getPlcSubscriptionType();
    }

    public PlcValue getLastValue() {
        return lastValue;
    }

    public void setLastValue(PlcValue lastValue) {
        this.lastValue = lastValue;
    }
}
