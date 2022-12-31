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
package org.apache.plc4x.java.utils.connectionpool2;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcTag;

import java.time.Duration;
import java.util.function.Consumer;

public class CachedSubscriptionRequestBuilder implements PlcSubscriptionRequest.Builder {

    private final CachedPlcConnection parent;
    private final PlcSubscriptionRequest.Builder builder;

    public CachedSubscriptionRequestBuilder(CachedPlcConnection parent, PlcSubscriptionRequest.Builder builder) {
        this.parent = parent;
        this.builder = builder;
    }

    @Override
    public PlcSubscriptionRequest build() {
        return builder.build();
    }

    @Override
    public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval) {
        return builder.addCyclicTagAddress(name, tagAddress, pollingInterval);
    }

    @Override
    public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval) {
        return builder.addCyclicTag(name, tag, pollingInterval);
    }

    @Override
    public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress) {
        return builder.addChangeOfStateTagAddress(name, tagAddress);
    }

    @Override
    public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag) {
        return builder.addChangeOfStateTag(name, tag);
    }

    @Override
    public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress) {
        return builder.addEventTagAddress(name, tagAddress);
    }

    @Override
    public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag) {
        return builder.addEventTag(name, tag);
    }

    @Override
    public PlcSubscriptionRequest.Builder addPreRegisteredConsumer(String name, Consumer<PlcSubscriptionEvent> preRegisteredConsumer) {
        return builder.addPreRegisteredConsumer(name, preRegisteredConsumer);
    }

}
