package org.apache.plc4x.java.api.messages;
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

import org.apache.plc4x.java.api.messages.items.*;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.concurrent.TimeUnit;
import java.util.function.Consumer;

public class PlcSubscriptionRequest extends PlcRequest<SubscriptionRequestItem<?>> {

    public static PlcSubscriptionRequest.Builder builder() {
        return new PlcSubscriptionRequest.Builder();
    }

    public static class Builder extends PlcRequest.Builder<SubscriptionRequestItem> {

        public final <T> Builder addChangeOfStateItem(Class<T> dataType, PlcField field, Consumer<SubscriptionEventItem<T>> consumer) {
            // As we don't get a list as response rather we have individual consumers we don't need type checking here.
            //checkType(dataType);
            requests.add(new SubscriptionRequestChangeOfStateItem<>(dataType, field, consumer));
            return this;
        }

        public final <T> Builder addCyclicItem(Class<T> dataType, PlcField field, Consumer<SubscriptionEventItem<T>> consumer, TimeUnit timeUnit, int period) {
            // As we don't get a list as response rather we have individual consumers we don't need type checking here.
            //checkType(dataType);
            requests.add(new SubscriptionRequestCyclicItem<>(dataType, field, timeUnit, period, consumer));
            return this;
        }

        public final <T> Builder addEventItem(Class<T> dataType, PlcField field, Consumer<SubscriptionEventItem<T>> consumer) {
            // As we don't get a list as response rather we have individual consumers we don't need type checking here.
            //checkType(dataType);
            requests.add(new SubscriptionRequestEventItem<>(dataType, field, consumer));
            return this;
        }

        public final Builder addItem(SubscriptionRequestItem subscriptionRequestItem) {
            requests.add(subscriptionRequestItem);
            return this;
        }

        public final PlcSubscriptionRequest build() {
            if (requests.isEmpty()) {
                throw new IllegalStateException("No requests added");
            }
            PlcSubscriptionRequest plcSubscriptionRequest = new PlcSubscriptionRequest();
            for (SubscriptionRequestItem request : requests) {
                plcSubscriptionRequest.addItem(request);
            }
            return plcSubscriptionRequest;
        }

    }

    @Override
    public String toString() {
        return "PlcSubscriptionRequest{} " + super.toString();
    }
}
