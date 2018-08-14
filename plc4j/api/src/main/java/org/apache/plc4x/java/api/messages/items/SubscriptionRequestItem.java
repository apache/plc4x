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
package org.apache.plc4x.java.api.messages.items;

import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.api.model.SubscriptionType;

import java.util.Objects;
import java.util.function.Consumer;

public abstract class SubscriptionRequestItem<T> extends RequestItem<T> {

    private SubscriptionType subscriptionType;
    private Consumer<SubscriptionEventItem<T>> consumer;

    public SubscriptionRequestItem(Class<T> datatype, Address address, SubscriptionType subscriptionType, Consumer<SubscriptionEventItem<T>> consumer) {
        super(datatype, address);
        this.subscriptionType = subscriptionType;
        this.consumer = consumer;
    }

    public SubscriptionType getSubscriptionType() {
        return subscriptionType;
    }

    public Consumer<SubscriptionEventItem<T>> getConsumer() {
        return consumer;
    }

    @Override
    public String toString() {
        return "SubscriptionRequestItem{" +
            "subscriptionType=" + subscriptionType +
            ", consumer=" + consumer +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof SubscriptionRequestItem)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        SubscriptionRequestItem<?> that = (SubscriptionRequestItem<?>) o;
        return subscriptionType == that.subscriptionType &&
            Objects.equals(consumer, that.consumer);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), subscriptionType, consumer);
    }
}
