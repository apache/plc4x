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
import org.apache.plc4x.java.base.messages.PlcSubscriber;

import java.util.Arrays;
import java.util.Collection;
import java.util.Objects;
import java.util.function.Consumer;

public class DefaultPlcConsumerRegistration implements InternalPlcConsumerRegistration {

    private final PlcSubscriber plcSubscriber;

    private final Collection<? extends InternalPlcSubscriptionHandle> handles;

    private final int consumerHash;

    public DefaultPlcConsumerRegistration(PlcSubscriber plcSubscriber, Consumer<PlcSubscriptionEvent> consumer, InternalPlcSubscriptionHandle... handles) {
        this.plcSubscriber = plcSubscriber;
        this.handles = Arrays.asList(Objects.requireNonNull(handles));
        this.consumerHash = Objects.requireNonNull(consumer).hashCode();
    }

    @Override
    public int getConsumerHash() {
        return consumerHash;
    }

    @Override
    public Collection<? extends InternalPlcSubscriptionHandle> getAssociatedHandles() {
        return handles;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof DefaultPlcConsumerRegistration)) {
            return false;
        }
        DefaultPlcConsumerRegistration that = (DefaultPlcConsumerRegistration) o;
        return consumerHash == that.consumerHash &&
            Objects.equals(handles, that.handles);
    }

    @Override
    public int hashCode() {
        return Objects.hash(handles, consumerHash);
    }

    @Override
    public String toString() {
        return "DefaultPlcConsumerRegistration{" +
            "handles=" + handles +
            ", consumerHash=" + consumerHash +
            '}';
    }

    @Override
    public void unregister() {
        plcSubscriber.unregister(this);
    }
}
