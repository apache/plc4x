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
import java.io.File;
import java.io.IOException;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonIgnore;

public class Configuration {

    @JsonIgnore
    private String configFile;

    @JsonProperty
    private String version;

    @JsonProperty(required=true)
    private String dir;

    @JsonProperty
    private String name;

    @JsonProperty
    private String adminUserName;

    @JsonProperty
    private String adminPassword;

    @JsonProperty
    private String securityPassword;

    @JsonProperty
    private List<DeviceConfiguration> devices;

    @JsonProperty
    private Integer tcpPort = 12686;

    @JsonProperty
    private Integer httpPort = 8443;

    public Configuration() {
    }

    public void setConfigFile(String value) {
        configFile = value;
    }

    public String getName() {
        return name;
    }

    public String getAdminUserName() {
        return adminUserName;
    }

    public String getAdminPassword() {
        return adminPassword;
    }

    public String getSecurityPassword() {
        return securityPassword;
    }

    @JsonIgnore
    public void setAdminUserName(String value) throws IOException {
        adminUserName = value;
        ObjectMapper om = new ObjectMapper(new YAMLFactory());
        om.writeValue(new File(configFile), this);
    }

    @JsonIgnore
    public void setAdminPassword(String value) throws IOException {
        //TODO: This need to be encrypted
        adminPassword = value;
        ObjectMapper om = new ObjectMapper(new YAMLFactory());
        om.writeValue(new File(configFile), this);
    }

    @JsonIgnore
    public void setSecurityPassword(String value) throws IOException {
        //TODO: This need to be encrypted
        securityPassword = value;
        ObjectMapper om = new ObjectMapper(new YAMLFactory());
        om.writeValue(new File(configFile), this);
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
        return devices;
    }

}
