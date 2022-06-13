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
package org.apache.plc4x.java.spi.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;

import java.util.Arrays;
import java.util.List;
import java.util.Objects;
import java.util.function.Consumer;

public class DefaultPlcConsumerRegistration implements PlcConsumerRegistration {

    private final PlcSubscriber plcSubscriber;

    private final List<PlcSubscriptionHandle> handles;

    private final int consumerHash;

    public DefaultPlcConsumerRegistration(PlcSubscriber plcSubscriber, Consumer<PlcSubscriptionEvent> consumer, PlcSubscriptionHandle... handles) {
        this.plcSubscriber = plcSubscriber;
        this.handles = Arrays.asList(Objects.requireNonNull(handles));
        this.consumerHash = Objects.requireNonNull(consumer).hashCode();
    }

    @Override
    public Integer getConsumerId() {
        return consumerHash;
    }

    @Override
    public List<PlcSubscriptionHandle> getSubscriptionHandles() {
        return handles;
    }

    @Override
    public void unregister() {
        plcSubscriber.unregister(this);
    }

    @Override
    public final boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof DefaultPlcConsumerRegistration)) {
            return false;
        }
        DefaultPlcConsumerRegistration that = (DefaultPlcConsumerRegistration) o;
        return Objects.equals(plcSubscriber, that.plcSubscriber) &&
            Objects.equals(handles, that.handles) &&
            consumerHash == that.consumerHash;
    }

    @Override
    public final int hashCode() {
        return Objects.hash(plcSubscriber, handles, consumerHash);
    }

    @Override
    public String toString() {
        return "DefaultPlcConsumerRegistration{" +
            "handles=" + handles +
            ", consumerHash=" + consumerHash +
            '}';
    }

}
