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

import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.value.PlcValue;

import java.util.Collections;
import java.util.Map;
import java.util.stream.Collectors;

public class DefaultPlcDiscoveryItem implements PlcDiscoveryItem {

    private final String protocolCode;
    private final String transportCode;
    private final String transportUrl;
    private final Map<String, String> options;
    private final String name;
    private final Map<String, PlcValue> attributes;

    public DefaultPlcDiscoveryItem(String protocolCode, String transportCode, String transportUrl,
                                   Map<String, String> options, String name,
                                   Map<String, PlcValue> attributes) {
        this.protocolCode = protocolCode;
        this.transportCode = transportCode;
        this.transportUrl = transportUrl;
        this.options = options != null ? options : Collections.emptyMap();
        this.name = name;
        this.attributes = attributes != null ? attributes : Collections.emptyMap();
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
        StringBuilder sb = new StringBuilder();
        sb.append(protocolCode);
        if (transportCode != null && !transportCode.isEmpty()) {
            sb.append(":").append(transportCode);
        }
        sb.append("://").append(transportUrl);
        if (!options.isEmpty()) {
            sb.append("?");
            sb.append(options.entrySet().stream()
                .map(e -> e.getKey() + "=" + e.getValue())
                .collect(Collectors.joining("&")));
        }
        return sb.toString();
    }

}
