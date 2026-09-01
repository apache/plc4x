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

import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;

import java.util.*;

public class DefaultPlcSubscriptionResponse implements PlcSubscriptionResponse {

    private final PlcSubscriptionRequest request;
    private final Map<String, PlcResponseItem<PlcSubscriptionHandle>> values;

    public DefaultPlcSubscriptionResponse(PlcSubscriptionRequest request,
                                          Map<String, PlcResponseItem<PlcSubscriptionHandle>> values) {
        this.request = request;
        this.values = values;
    }

    @Override
    public PlcSubscriptionRequest getRequest() {
        return request;
    }

    @Override
    public Collection<String> getTagNames() {
        return request.getTagNames();
    }

    @Override
    public PlcSubscriptionTag getTag(String name) {
        return request.getTag(name);
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        if (values.get(name) == null) {
            return PlcResponseCode.NOT_FOUND;
        }
        return values.get(name).getResponseCode();
    }

    @Override
    public PlcSubscriptionHandle getSubscriptionHandle(String name) {
        if (values.get(name) == null) {
            return null;
        }
        return values.get(name).getValue();
    }

    @Override
    public Collection<PlcSubscriptionHandle> getSubscriptionHandles() {
        // A single handle can cover several tags - drivers commonly create one handle for the whole
        // request and map it to every tag name. Returning it once per tag would make the usual
        // "register a consumer on every handle" loop deliver each event N times, see GH-1896.
        // Duplicates are detected by identity: two distinct handles that happen to compare equal
        // are still two subscriptions.
        Map<PlcSubscriptionHandle, Boolean> seen = new IdentityHashMap<>();
        List<PlcSubscriptionHandle> handles = new ArrayList<>(values.size());
        for (PlcResponseItem<PlcSubscriptionHandle> item : values.values()) {
            if (item != null && item.getValue() != null
                && seen.put(item.getValue(), Boolean.TRUE) == null) {
                handles.add(item.getValue());
            }
        }
        return handles;
    }

}
