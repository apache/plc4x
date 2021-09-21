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
package org.apache.plc4x.java.opcuaserver.configuration;

import org.eclipse.jetty.util.security.Password;

import java.util.HashMap;
import java.util.Map;
import java.io.File;
import java.io.IOException;

import java.nio.file.Path;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonIgnore;

public class PasswordConfiguration {

    @JsonIgnore
    private String passwordConfigFile;

    @JsonProperty
    private String version;

    @JsonProperty
    private String securityPassword;

    @JsonProperty
    private Map<String, User> users = new HashMap<>();

    public PasswordConfiguration() {
    }

    public boolean checkPassword(String username, String password) {
        if (users.containsKey(username)) {
            return users.get(username).checkPassword(password);
        }
        return false;
    }

    public void createUser(String username, String password, String security) throws IOException {
        User user = new User(username, password, security);
        users.put(username, user);
        ObjectMapper om = new ObjectMapper(new YAMLFactory());
        om.writeValue(new File(passwordConfigFile), this);
    }

    @JsonIgnore
    public String getSecurityPassword() {
        return Password.deobfuscate("OBF:" + securityPassword);
    }

    @JsonIgnore
    public void setSecurityPassword(String value) throws IOException {
        securityPassword = Password.obfuscate(value).substring(4);
        ObjectMapper om = new ObjectMapper(new YAMLFactory());
        om.writeValue(new File(passwordConfigFile), this);
    }

    public void setPasswordConfigFile(Path value) {
        passwordConfigFile = value.toString();
    }

    public void setVersion(String value) {
        version = value;
    }

}
