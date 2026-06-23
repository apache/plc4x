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
package org.apache.plc4x.java.plc4x.config;

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;

public class Plc4xConfiguration implements Configuration {

    @ConfigurationParameter("remote-connection-string")
    @Description("URL-Encoded connection string to use on the proxy side to reach the given PLC.")
    private String remoteConnectionString;

    @ConfigurationParameter("request-timeout")
    @IntDefaultValue(5_000)
    @Description("Default timeout for all types of requests.")
    private int requestTimeout;

    @ConfigurationParameter("username")
    @Description("Username for authenticating against the PLC4X proxy server. Authentication is mandatory.")
    private String username;

    @ConfigurationParameter("password")
    @Description("Password for authenticating against the PLC4X proxy server. Authentication is mandatory.")
    private String password;

    public String getRemoteConnectionString() {
        return remoteConnectionString;
    }

    public void setRemoteConnectionString(String remoteConnectionString) {
        this.remoteConnectionString = remoteConnectionString;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public int getRequestTimeout() {
        return requestTimeout;
    }

    public void setRequestTimeout(int requestTimeout) {
        this.requestTimeout = requestTimeout;
    }

}
