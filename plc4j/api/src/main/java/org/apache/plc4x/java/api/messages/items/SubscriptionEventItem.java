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

import java.util.Calendar;
import java.util.List;
import java.util.Objects;

public class SubscriptionEventItem<T> {

    private SubscriptionRequestItem<T> subscriptionRequestItem;
    private Calendar timestamp;
    private List<T> values;

    public SubscriptionEventItem(SubscriptionRequestItem<T> subscriptionRequestItem, Calendar timestamp, List<T> values) {
        this.subscriptionRequestItem = subscriptionRequestItem;
        this.timestamp = timestamp;
        this.values = values;
    }

    public SubscriptionRequestItem<T> getSubscriptionRequestItem() {
        return subscriptionRequestItem;
    }

    public Calendar getTimestamp() {
        return timestamp;
    }

    public List<T> getValues() {
        return values;
    }

    @Override
    public String toString() {
        return "SubscriptionEventItem{" +
            "subscriptionRequestItem=" + subscriptionRequestItem +
            ", timestamp=" + timestamp +
            ", values=" + values +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof SubscriptionEventItem)) {
            return false;
        }
        SubscriptionEventItem<?> that = (SubscriptionEventItem<?>) o;
        return Objects.equals(subscriptionRequestItem, that.subscriptionRequestItem) &&
            Objects.equals(timestamp, that.timestamp) &&
            Objects.equals(values, that.values);
    }

    @Override
    public int hashCode() {
        return Objects.hash(subscriptionRequestItem, timestamp, values);
    }
}
