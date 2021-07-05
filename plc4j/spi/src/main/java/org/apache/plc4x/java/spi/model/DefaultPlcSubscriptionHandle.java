/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.spi.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;

import java.util.Collections;
import java.util.function.Consumer;

// Warning: do not override equals and hashCode as these should not include the plcSubscriber.
public class DefaultPlcSubscriptionHandle implements PlcSubscriptionHandle {

    private final transient PlcSubscriber plcSubscriber;

    public DefaultPlcSubscriptionHandle(PlcSubscriber plcSubscriber) {
        this.plcSubscriber = plcSubscriber;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return plcSubscriber.register(consumer, Collections.singletonList(this));
    }

    @Override
    public boolean equals(Object obj) {
        if (!(obj instanceof DefaultPlcSubscriptionHandle)) {
            return false;
        }
        // A handle is unique therefore we use the default implementation from Object
        return (this == obj);
    }

    @Override
    public int hashCode() {
        // A handle is unique therefore we use the default implementation from Object
        return System.identityHashCode(this);
    }

    @Override
    public String toString() {
        return "DefaultPlcSubscriptionHandle{" +
            "plcSubscriber=" + plcSubscriber +
            '}';
    }

}
