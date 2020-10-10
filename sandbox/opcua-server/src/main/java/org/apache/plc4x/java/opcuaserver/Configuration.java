/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.opcuaserver;

import java.util.List;

import com.fasterxml.jackson.annotation.JsonProperty;

public class Configuration {

    @JsonProperty
    private String version;

    @JsonProperty(required=true)
    private String dir;

    @JsonProperty
    private String name;

    @JsonProperty
    private List<DeviceConfiguration> modbusDevices;

    @JsonProperty
    private Integer tcpPort = 12686;

    @JsonProperty
    private Integer httpPort = 8443;

    public Configuration() {
    }

    public String getName() {
        return name;
    }

    public String getDir() {
        return dir;
    }

    public Integer getTcpPort() {
        return tcpPort;
    }

    public Integer getHttpPort() {
        return httpPort;
    }

    public List<DeviceConfiguration> getDevices() {
        return modbusDevices;
    }

}
