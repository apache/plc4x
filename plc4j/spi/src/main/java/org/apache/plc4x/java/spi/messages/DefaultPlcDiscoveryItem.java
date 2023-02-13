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
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.Map;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcDiscoveryItem implements PlcDiscoveryItem, Serializable {

    private final String protocolCode;
    private final String transportCode;
    private final String transportUrl;
    private final Map<String, String> options;
    private final String name;

    private final Map<String, PlcValue> attributes;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcDiscoveryItem(@JsonProperty("protocolCode") String protocolCode,
                                   @JsonProperty("transportCode") String transportCode,
                                   @JsonProperty("transportUrl") String transportUrl,
                                   @JsonProperty("options") Map<String, String> options,
                                   @JsonProperty("name") String name,
                                   @JsonProperty("options") Map<String, PlcValue> attributes) {
        this.protocolCode = protocolCode;
        this.transportCode = transportCode;
        this.transportUrl = transportUrl;
        this.options = options;
        this.name = name;
        this.attributes = attributes;
    }

    @Override
    public String getProtocolCode() {
        return protocolCode;
    }

    @Override
    public String getTransportCode() {
        return transportCode;
    }

    @Override
    public String getTransportUrl() {
        return transportUrl;
    }

    @Override
    public Map<String, String> getOptions() {
        return options;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public Map<String, PlcValue> getAttributes() {
        return attributes;
    }

    @Override
    public String getConnectionUrl() {
        StringBuilder sb = new StringBuilder(String.format("%s:%s://%s",
            protocolCode, transportCode, transportUrl.toString()));
        if(options != null && !options.isEmpty()) {
            boolean first = true;
            for (Map.Entry<String, String> optionEntry : options.entrySet()) {
                if(first) {
                    sb.append("?");
                    first = false;
                } else {
                    sb.append("&");
                }
                sb.append(optionEntry.getKey()).append("=").append(optionEntry.getValue());
            }
        }
        return sb.toString();
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeString("protocolCode",
            protocolCode.getBytes(StandardCharsets.UTF_8).length * 8,
            protocolCode, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        writeBuffer.writeString("transportCode",
            transportCode.getBytes(StandardCharsets.UTF_8).length * 8,
            transportCode, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        writeBuffer.writeString("transportUrl",
            transportUrl.toString().getBytes(StandardCharsets.UTF_8).length * 8,
            transportUrl.toString(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        if(options != null && !options.isEmpty()) {
            writeBuffer.pushContext("options");
            for (Map.Entry<String, String> optionEntry : options.entrySet()) {
                writeBuffer.pushContext("option");
                writeBuffer.writeString("name",
                    optionEntry.getKey().getBytes(StandardCharsets.UTF_8).length * 8,
                    optionEntry.getKey(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
                writeBuffer.writeString("value",
                    optionEntry.getValue().getBytes(StandardCharsets.UTF_8).length * 8,
                    optionEntry.getValue(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
                writeBuffer.popContext("option");
            }
            writeBuffer.popContext("options");
        }
        if(name != null && !name.isEmpty()) {
            writeBuffer.writeString("name",
                name.getBytes(StandardCharsets.UTF_8).length * 8,
                name, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        }

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
