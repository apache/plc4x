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
package org.apache.plc4x.java.spi.messages;

import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.Collection;
import java.util.Map;
import java.util.stream.Collectors;

public class DefaultPlcSubscriptionResponse implements PlcSubscriptionResponse, PlcResponse, Serializable {

    private final PlcSubscriptionRequest request;

    private final Map<String, ResponseItem<PlcSubscriptionHandle>> values;

    public DefaultPlcSubscriptionResponse(PlcSubscriptionRequest request,
                                          Map<String, ResponseItem<PlcSubscriptionHandle>> values) {
        this.request = request;
        this.values = values;
        request.getPreRegisteredConsumers().forEach((subscriptionTagName, consumers) -> {
            PlcSubscriptionHandle subscriptionHandle = getSubscriptionHandle(subscriptionTagName);
            if (subscriptionHandle == null) {
                throw new PlcRuntimeException("PlcSubscriptionHandle for " + subscriptionTagName + " not found");
            }
            consumers.forEach(subscriptionHandle::register);
        });
    }

    @Override
    public PlcSubscriptionHandle getSubscriptionHandle(String name) {
        ResponseItem<PlcSubscriptionHandle> response = values.get(name);
        if (response == null) {
            return null;
        }
        if (response.getCode() != PlcResponseCode.OK) {
            throw new PlcRuntimeException("Item " + name + " failed to subscribe: " + response.getCode());
        }
        return response.getValue();
    }

    @Override
    public Collection<String> getTagNames() {
        return values.keySet();
    }

    @Override
    public PlcSubscriptionTag getTag(String name) {
        throw new PlcNotImplementedException("tag access not possible as these come async");
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        ResponseItem<PlcSubscriptionHandle> response = values.get(name);
        if (response == null) {
            return null;
        }
        return response.getCode();
    }

    @Override
    public PlcSubscriptionRequest getRequest() {
        return request;
    }

    @Override
    public Collection<PlcSubscriptionHandle> getSubscriptionHandles() {
        return values.values().stream().map(ResponseItem::getValue).collect(Collectors.toList());
    }

    public Map<String, ResponseItem<PlcSubscriptionHandle>> getValues() {
        return values;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcSubscriptionResponse");

        if (request instanceof Serializable) {
            ((Serializable) request).serialize(writeBuffer);
        }
        writeBuffer.pushContext("values");
        for (Map.Entry<String, ResponseItem<PlcSubscriptionHandle>> valueEntry : values.entrySet()) {
            String tagName = valueEntry.getKey();
            writeBuffer.pushContext(tagName);
            ResponseItem<PlcSubscriptionHandle> valueResponse = valueEntry.getValue();
            valueResponse.serialize(writeBuffer);
            writeBuffer.pushContext(tagName);
        }
        writeBuffer.popContext("values");

        writeBuffer.popContext("PlcSubscriptionResponse");
    }

}
