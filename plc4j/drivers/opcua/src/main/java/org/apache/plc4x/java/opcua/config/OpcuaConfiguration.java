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
package org.apache.plc4x.java.opcua.config;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationProvider;
import org.apache.plc4x.java.transport.tcp.DefaultTcpTransportConfiguration;

public class OpcuaConfiguration implements Configuration, TransportConfigurationProvider {

    @ConfigurationParameter("protocolCode")
    private String protocolCode;

    @ConfigurationParameter("transportCode")
    private String transportCode;

    @ConfigurationParameter("transportConfig")
    private String transportConfig;

    @ConfigurationParameter("discovery")
    @BooleanDefaultValue(true)
    private boolean discovery;

    @ConfigurationParameter("username")
    private String username;

    @ConfigurationParameter("password")
    private String password;

    @ConfigurationParameter("securityPolicy")
    @StringDefaultValue("None")
    private String securityPolicy;

    @ConfigurationParameter("keyStoreFile")
    private String keyStoreFile;

    @ConfigurationParameter("certDirectory")
    private String certDirectory;

    @ConfigurationParameter("keyStorePassword")
    private String keyStorePassword;

    @ComplexConfigurationParameter(prefix = "tcp", defaultOverrides = {}, requiredOverrides = {})
    private DefaultTcpTransportConfiguration tcpTransportConfiguration;

    public String getProtocolCode() {
        return protocolCode;
    }

    public String getTransportCode() {
        return transportCode;
    }

    public String getTransportConfig() {
        return transportConfig;
    }

    public boolean isDiscovery() {
        return discovery;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }

    public String getCertDirectory() {
        return certDirectory;
    }

    public String getSecurityPolicy() {
        return securityPolicy;
    }

    public String getKeyStoreFile() {
        return keyStoreFile;
    }

    public String getKeyStorePassword() {
        return keyStorePassword;
    }

    public DefaultTcpTransportConfiguration getTcpTransportConfiguration() {
        return tcpTransportConfiguration;
    }

    public void setTcpTransportConfiguration(DefaultTcpTransportConfiguration tcpTransportConfiguration) {
        this.tcpTransportConfiguration = tcpTransportConfiguration;
    }

    @Override
    public TransportConfiguration getTransportConfiguration(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return tcpTransportConfiguration;
        }
        return null;
    }

    @Override
    public String toString() {
        return "OpcuaConfiguration{" +
            "discovery=" + discovery +
            ", username='" + username + '\'' +
            ", password='" + (password != null ? "******" : null) + '\'' +
            ", securityPolicy='" + securityPolicy + '\'' +
            ", keyStoreFile='" + keyStoreFile + '\'' +
            ", certDirectory='" + certDirectory + '\'' +
            ", keyStorePassword='" + (keyStorePassword != null ? "******" : null) + '\'' +
            '}';
    }
}

