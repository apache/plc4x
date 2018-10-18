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

package org.apache.plc4x.java.base.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.base.messages.PlcSubscriber;

import java.util.Collections;
import java.util.Objects;
import java.util.function.Consumer;

public class DefaultPlcSubscriptionHandle implements InternalPlcSubscriptionHandle {

    private final PlcSubscriber plcSubscriber;

    public DefaultPlcSubscriptionHandle(PlcSubscriber plcSubscriber) {
        this.plcSubscriber = plcSubscriber;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return plcSubscriber.register(consumer, Collections.singletonList(this));
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof DefaultPlcSubscriptionHandle)) {
            return false;
        }
        DefaultPlcSubscriptionHandle that = (DefaultPlcSubscriptionHandle) o;
        return Objects.equals(plcSubscriber, that.plcSubscriber);
    }

    @Override
    public int hashCode() {
        return Objects.hash(plcSubscriber);
    }

    @Override
    public String toString() {
        return "DefaultPlcSubscriptionHandle{" +
            "plcSubscriber=" + plcSubscriber +
            '}';
    }
}
