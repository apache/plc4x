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
package org.apache.plc4x.java.spi.messages;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcUnsubscriptionRequest implements PlcUnsubscriptionRequest, PlcRequest, Serializable {

    private final PlcSubscriber subscriber;

    private final List<PlcSubscriptionHandle> plcSubscriptionHandles;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcUnsubscriptionRequest(@JsonProperty("subscriber") PlcSubscriber subscriber,
                                           @JsonProperty("internalPlcSubscriptionHandles") List<PlcSubscriptionHandle> plcSubscriptionHandles) {
        this.subscriber = subscriber;
        this.plcSubscriptionHandles = plcSubscriptionHandles;
    }

    @Override
    @JsonIgnore
    public CompletableFuture<PlcUnsubscriptionResponse> execute() {
        return subscriber.unsubscribe(this);
    }

    @Override
    public List<PlcSubscriptionHandle> getSubscriptionHandles() {
        return plcSubscriptionHandles;
    }

    public PlcSubscriber getSubscriber() {
        return subscriber;
    }

    public Collection<PlcSubscriptionHandle> getPlcSubscriptionHandles() {
        return plcSubscriptionHandles;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        // TODO: Implement
    }

    public static class Builder implements PlcUnsubscriptionRequest.Builder {

        private final PlcSubscriber subscriber;
        private List<PlcSubscriptionHandle> plcSubscriptionHandles;

        public Builder(PlcSubscriber subscriber) {
            this.subscriber = subscriber;
            plcSubscriptionHandles = new ArrayList<>();
        }

        public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle plcSubscriptionHandle) {
            plcSubscriptionHandles.add(plcSubscriptionHandle);
            return this;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle plcSubscriptionHandle1, PlcSubscriptionHandle... plcSubscriptionHandles) {
            this.plcSubscriptionHandles.add(plcSubscriptionHandle1);
            this.plcSubscriptionHandles.addAll(Arrays.stream(plcSubscriptionHandles).map(PlcSubscriptionHandle.class::cast).collect(Collectors.toList()));
            return this;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(Collection<PlcSubscriptionHandle> plcSubscriptionHandles) {
            this.plcSubscriptionHandles.addAll(plcSubscriptionHandles.stream().map(PlcSubscriptionHandle.class::cast).collect(Collectors.toList()));
            return this;
        }

        @Override
        public PlcUnsubscriptionRequest build() {
            return new DefaultPlcUnsubscriptionRequest(subscriber, plcSubscriptionHandles);
        }

    }

}
