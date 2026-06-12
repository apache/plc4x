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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.CompletableFuture;

public class DefaultPlcUnsubscriptionRequest implements PlcUnsubscriptionRequest {

    private final PlcSubscriber subscriber;
    private final List<PlcSubscriptionHandle> handles;

    public DefaultPlcUnsubscriptionRequest(PlcSubscriber subscriber, List<PlcSubscriptionHandle> handles) {
        this.subscriber = subscriber;
        this.handles = handles;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> execute() {
        return subscriber.unsubscribe(this);
    }

    @Override
    public List<PlcSubscriptionHandle> getSubscriptionHandles() {
        return handles;
    }

    public static class Builder implements PlcUnsubscriptionRequest.Builder {

        private final PlcSubscriber subscriber;
        private final List<PlcSubscriptionHandle> handles = new ArrayList<>();

        public Builder(PlcSubscriber subscriber) {
            this.subscriber = subscriber;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle handle) {
            handles.add(handle);
            return this;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle handle1, PlcSubscriptionHandle... moreHandles) {
            handles.add(handle1);
            Collections.addAll(handles, moreHandles);
            return this;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(Collection<PlcSubscriptionHandle> handles) {
            this.handles.addAll(handles);
            return this;
        }

        @Override
        public PlcUnsubscriptionRequest build() {
            return new DefaultPlcUnsubscriptionRequest(subscriber, handles);
        }
    }

}
