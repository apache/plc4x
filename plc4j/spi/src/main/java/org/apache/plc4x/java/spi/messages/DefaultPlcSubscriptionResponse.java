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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.w3c.dom.Element;

import java.util.Collection;
import java.util.Map;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcSubscriptionResponse implements PlcSubscriptionResponse, PlcResponse, Serializable {

    private final PlcSubscriptionRequest request;

    private final Map<String, ResponseItem<PlcSubscriptionHandle>> values;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcSubscriptionResponse(@JsonProperty("request") PlcSubscriptionRequest request,
                                          @JsonProperty("values") Map<String, ResponseItem<PlcSubscriptionHandle>> values) {
        this.request = request;
        this.values = values;
        request.getPreRegisteredConsumers().forEach((subscriptionFieldName, consumers) -> {
            PlcSubscriptionHandle subscriptionHandle = getSubscriptionHandle(subscriptionFieldName);
            if (subscriptionHandle == null) {
                throw new PlcRuntimeException("PlcSubscriptionHandle for " + subscriptionFieldName + " not found");
            }
            consumers.forEach(subscriptionHandle::register);
        });
    }

    @Override
    @JsonIgnore
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
    @JsonIgnore
    public Collection<String> getFieldNames() {
        return values.keySet();
    }

    @Override
    @JsonIgnore
    public PlcSubscriptionField getField(String name) {
        throw new PlcNotImplementedException("field access not possible as these come async");
    }

    @Override
    @JsonIgnore
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
    @JsonIgnore
    public Collection<PlcSubscriptionHandle> getSubscriptionHandles() {
        return values.values().stream().map(ResponseItem<PlcSubscriptionHandle>::getValue).collect(Collectors.toList());
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
            String fieldName = valueEntry.getKey();
            writeBuffer.pushContext(fieldName);
            ResponseItem<PlcSubscriptionHandle> valueResponse = valueEntry.getValue();
            valueResponse.serialize(writeBuffer);
            writeBuffer.pushContext(fieldName);
        }
        writeBuffer.popContext("values");

        writeBuffer.popContext("PlcSubscriptionResponse");
    }

}
