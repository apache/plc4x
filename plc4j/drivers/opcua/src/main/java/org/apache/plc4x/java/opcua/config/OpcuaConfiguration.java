/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.opcua.config;


import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.transport.tcp.TcpTransportConfiguration;

public class OpcuaConfiguration implements Configuration, TcpTransportConfiguration {

    private String code;
    private String host;
    private String port;
    private String endpoint;
    private String params;

    @ConfigurationParameter("discovery")
    @BooleanDefaultValue(true)
    private boolean discovery;

    @ConfigurationParameter("username")
    private String username;

    @ConfigurationParameter("password")
    private String password;

    @ConfigurationParameter("certFile")
    private String certFile;

    @ConfigurationParameter("securityPolicy")
    private String securityPolicy;

    @ConfigurationParameter("keyStoreFile")
    private String keyStoreFile;

    public boolean isDiscovery() {
        return discovery;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }

    public String getCertFile() {
        return certFile;
    }

    public String getSecurityPolicy() {
        return securityPolicy;
    }

    public String getKeyStoreFile() {
        return keyStoreFile;
    }

    public void setDiscovery(boolean discovery) {
        this.discovery = discovery;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public void setCertFile(String certFile) {
        this.certFile = certFile;
    }

    public void setSecurityPolicy(String securityPolicy) {
        this.securityPolicy = securityPolicy;
    }

    public void setKeyStoreFile(String keyStoreFile) {
        this.keyStoreFile = keyStoreFile;
    }

    public String getTransportCode() {
        return code;
    }

    public String getHost() {
        return host;
    }

    public String getPort() {
        return port;
    }

    public String getEndpoint() {
        return endpoint;
    }

    public void setTransportCode(String code) {
        this.code = code;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public void setPort(String port) {
        this.port = port;
    }

    public void setEndpoint(String endpoint) {
        this.endpoint = endpoint;
    }



    @Override
    public String toString() {
        return "Configuration{" +
            '}';
    }

}

