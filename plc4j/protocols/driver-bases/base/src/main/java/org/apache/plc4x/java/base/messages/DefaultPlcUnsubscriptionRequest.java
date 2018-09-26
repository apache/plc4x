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
package org.apache.plc4x.java.base.messages;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.base.model.InternalPlcSubscriptionHandle;

import java.util.*;
import java.util.stream.Collectors;

public class DefaultPlcUnsubscriptionRequest implements InternalPlcUnsubscriptionRequest, InternalPlcFieldRequest {

    private final Collection<? extends InternalPlcSubscriptionHandle> internalPlcSubscriptionHandles;

    public DefaultPlcUnsubscriptionRequest(Collection<? extends InternalPlcSubscriptionHandle> internalPlcSubscriptionHandles) {
        this.internalPlcSubscriptionHandles = internalPlcSubscriptionHandles;
    }

    @Override
    public int getNumberOfFields() {
        throw new IllegalStateException("not available");
    }

    @Override
    public LinkedHashSet<String> getFieldNames() {
        throw new IllegalStateException("not available");
    }

    @Override
    public PlcField getField(String name) {
        throw new IllegalStateException("not available");
    }

    @Override
    public LinkedList<PlcField> getFields() {
        throw new IllegalStateException("not available");
    }

    @Override
    public Collection<? extends InternalPlcSubscriptionHandle> getInternalPlcSubscriptionHandles() {
        return internalPlcSubscriptionHandles;
    }

    @Override
    public LinkedList<Pair<String, PlcField>> getNamedFields() {
        throw new IllegalStateException("not available");
    }

    public static class Builder implements PlcUnsubscriptionRequest.Builder {

        private List<InternalPlcSubscriptionHandle> plcSubscriptionHandles;

        public Builder() {
            plcSubscriptionHandles = new ArrayList<>();
        }

        public PlcUnsubscriptionRequest.Builder addHandle(PlcSubscriptionHandle plcSubscriptionHandle) {
            plcSubscriptionHandles.add((InternalPlcSubscriptionHandle) plcSubscriptionHandle);
            return this;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandle(PlcSubscriptionHandle plcSubscriptionHandle1, PlcSubscriptionHandle... plcSubscriptionHandles) {
            this.plcSubscriptionHandles.add((InternalPlcSubscriptionHandle) plcSubscriptionHandle1);
            this.plcSubscriptionHandles.addAll(Arrays.stream(plcSubscriptionHandles).map(InternalPlcSubscriptionHandle.class::cast).collect(Collectors.toList()));
            return null;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder addHandles(Collection<PlcSubscriptionHandle> plcSubscriptionHandles) {
            this.plcSubscriptionHandles.addAll(plcSubscriptionHandles.stream().map(InternalPlcSubscriptionHandle.class::cast).collect(Collectors.toList()));
            return null;
        }

        @Override
        public PlcUnsubscriptionRequest build() {
            return new DefaultPlcUnsubscriptionRequest(plcSubscriptionHandles);
        }


    }

}
