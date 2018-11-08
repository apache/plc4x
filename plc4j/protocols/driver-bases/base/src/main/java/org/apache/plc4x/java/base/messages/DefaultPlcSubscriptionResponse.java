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
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.Collection;
import java.util.Map;
import java.util.stream.Collectors;

public class DefaultPlcSubscriptionResponse implements InternalPlcSubscriptionResponse {

    private final InternalPlcSubscriptionRequest request;

    private final Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> values;

    public DefaultPlcSubscriptionResponse(InternalPlcSubscriptionRequest request, Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> values) {
        this.request = request;
        this.values = values;
    }

    @Override
    public PlcSubscriptionHandle getSubscriptionHandle(String name) {
        Pair<PlcResponseCode, PlcSubscriptionHandle> response = values.get(name);
        if (response == null) {
            return null;
        }
        if (response.getKey() != PlcResponseCode.OK) {
            throw new PlcRuntimeException("Item " + name + " failed to subscribe: " + response.getKey());
        }
        return response.getValue();
    }

    @Override
    public Collection<String> getFieldNames() {
        return values.keySet();
    }

    @Override
    public PlcField getField(String name) {
        throw new PlcNotImplementedException("field access not possible as these come async");
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        Pair<PlcResponseCode, PlcSubscriptionHandle> response = values.get(name);
        if (response == null) {
            return null;
        }
        return response.getKey();
    }

    @Override
    public PlcSubscriptionRequest getRequest() {
        return request;
    }

    @Override
    public Collection<PlcSubscriptionHandle> getSubscriptionHandles() {
        return values.values().stream().map(Pair::getValue).collect(Collectors.toList());
    }

    @Override
    public Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> getValues() {
        return values;
    }

}
