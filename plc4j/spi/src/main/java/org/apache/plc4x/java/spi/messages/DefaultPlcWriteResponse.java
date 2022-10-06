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
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.Collection;
import java.util.Map;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcWriteResponse implements PlcWriteResponse, Serializable {

    private final PlcWriteRequest request;
    private final Map<String, PlcResponseCode> responses;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcWriteResponse(@JsonProperty("request") PlcWriteRequest request,
                                   @JsonProperty("values") Map<String, PlcResponseCode> responses) {
        this.request = request;
        this.responses = responses;
    }

    @Override
    public PlcWriteRequest getRequest() {
        return request;
    }

    @Override
    @JsonIgnore
    public Collection<String> getFieldNames() {
        return request.getFieldNames();
    }

    @Override
    @JsonIgnore
    public PlcField getField(String name) {
        return request.getField(name);
    }

    @Override
    @JsonIgnore
    public PlcResponseCode getResponseCode(String name) {
        return responses.get(name);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcWriteResponse");

        if (request instanceof Serializable) {
            ((Serializable) request).serialize(writeBuffer);
        }
        writeBuffer.pushContext("fields");
        for (Map.Entry<String, PlcResponseCode> fieldEntry : responses.entrySet()) {
            String fieldName = fieldEntry.getKey();
            final PlcResponseCode fieldResponseCode = fieldEntry.getValue();
            String result = fieldResponseCode.name();
            writeBuffer.writeString(fieldName, result.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), result);
        }
        writeBuffer.popContext("fields");

        writeBuffer.popContext("PlcWriteResponse");
    }

}
