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
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.Collection;
import java.util.Map;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcWriteResponse implements PlcWriteResponse, Serializable {

    private final PlcWriteRequest request;
    private final Map<String, PlcResponseCode> responseCodes;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcWriteResponse(@JsonProperty("request") PlcWriteRequest request,
                                   @JsonProperty("responseCodes") Map<String, PlcResponseCode> responseCodes) {
        this.request = request;
        this.responseCodes = responseCodes;
    }

    @Override
    public PlcWriteRequest getRequest() {
        return request;
    }

    @Override
    @JsonIgnore
    public Collection<String> getTagNames() {
        return request.getTagNames();
    }

    @Override
    @JsonIgnore
    public PlcTag getTag(String name) {
        return request.getTag(name);
    }

    @Override
    @JsonIgnore
    public PlcResponseCode getResponseCode(String name) {
        return responseCodes.get(name);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcWriteResponse");

        if (request instanceof Serializable) {
            ((Serializable) request).serialize(writeBuffer);
        }
        writeBuffer.pushContext("tags");
        for (Map.Entry<String, PlcResponseCode> tagEntry : responseCodes.entrySet()) {
            String tagName = tagEntry.getKey();
            final PlcResponseCode tagResponseCode = tagEntry.getValue();
            String result = tagResponseCode.name();
            writeBuffer.writeString(tagName,
                result.getBytes(StandardCharsets.UTF_8).length * 8,
                result, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        }
        writeBuffer.popContext("tags");

        writeBuffer.popContext("PlcWriteResponse");
    }

}
