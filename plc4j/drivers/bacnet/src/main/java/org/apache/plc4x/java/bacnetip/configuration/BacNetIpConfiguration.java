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
package org.apache.plc4x.java.bacnetip.configuration;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameterDefaultOverride;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationProvider;

public class BacNetIpConfiguration implements Configuration, TransportConfigurationProvider {

    // Path to a single EDE file.
    @ConfigurationParameter("ede-file-path")
    private String edeFilePath;

    // Path to a directory containing many EDE files.
    @ConfigurationParameter("ede-directory-path")
    private String edeDirectoryPath;

    @ComplexConfigurationParameter(prefix = "udp", defaultOverrides = {}, requiredOverrides = {})
    private BacNetUdpTransportConfiguration udpTransportConfiguration;

    @ComplexConfigurationParameter(prefix = "pcap", defaultOverrides = {
        @ComplexConfigurationParameterDefaultOverride(name = "support-vlans", value = "true")
    }, requiredOverrides = {})
    private BacNetPcapReplayTransportConfiguration pcapReplayTransportConfiguration;

    @ComplexConfigurationParameter(prefix = "raw", defaultOverrides = {
        @ComplexConfigurationParameterDefaultOverride(name = "resolve-mac-address", value = "true")
    }, requiredOverrides = {})
    private BacNetRawSocketTransportConfiguration rawSocketTransportConfiguration;

    public String getEdeFilePath() {
        return edeFilePath;
    }

    public void setEdeFilePath(String edeFilePath) {
        this.edeFilePath = edeFilePath;
    }

    public String getEdeDirectoryPath() {
        return edeDirectoryPath;
    }

    public void setEdeDirectoryPath(String edeDirectoryPath) {
        this.edeDirectoryPath = edeDirectoryPath;
    }

    public BacNetUdpTransportConfiguration getUdpTransportConfiguration() {
        return udpTransportConfiguration;
    }

    public void setUdpTransportConfiguration(BacNetUdpTransportConfiguration udpTransportConfiguration) {
        this.udpTransportConfiguration = udpTransportConfiguration;
    }

    public BacNetPcapReplayTransportConfiguration getPcapReplayTransportConfiguration() {
        return pcapReplayTransportConfiguration;
    }

    public void setPcapReplayTransportConfiguration(BacNetPcapReplayTransportConfiguration pcapReplayTransportConfiguration) {
        this.pcapReplayTransportConfiguration = pcapReplayTransportConfiguration;
    }

    public BacNetRawSocketTransportConfiguration getRawSocketTransportConfiguration() {
        return rawSocketTransportConfiguration;
    }

    public void setRawSocketTransportConfiguration(BacNetRawSocketTransportConfiguration rawSocketTransportConfiguration) {
        this.rawSocketTransportConfiguration = rawSocketTransportConfiguration;
    }

    @Override
    public TransportConfiguration getTransportConfiguration(String transportCode) {
        switch (transportCode) {
            case "udp":
                return udpTransportConfiguration;
            case "raw":
                return rawSocketTransportConfiguration;
            case "pcap":
                return pcapReplayTransportConfiguration;
        }
        return null;
    }

}
