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

import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcTag;

import java.time.Duration;
import java.util.Collection;
import java.util.function.Consumer;

public class CachedUnsubscriptionRequestBuilder implements PlcUnsubscriptionRequest.Builder {

    private final CachedPlcConnection parent;
    private final PlcUnsubscriptionRequest.Builder builder;

    public CachedUnsubscriptionRequestBuilder(CachedPlcConnection parent, PlcUnsubscriptionRequest.Builder builder) {
        this.parent = parent;
        this.builder = builder;
    }

    @Override
    public PlcUnsubscriptionRequest build() {
        return builder.build();
    }

    @Override
    public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle plcSubscriptionHandle) {
        return builder.addHandles(plcSubscriptionHandle);
    }

    @Override
    public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle plcSubscriptionHandle, PlcSubscriptionHandle... plcSubscriptionHandles) {
        return builder.addHandles(plcSubscriptionHandle, plcSubscriptionHandles);
    }

    @Override
    public PlcUnsubscriptionRequest.Builder addHandles(Collection<PlcSubscriptionHandle> plcSubscriptionHandles) {
        return builder.addHandles(plcSubscriptionHandles);
    }

}
