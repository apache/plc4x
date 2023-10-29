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
package org.apache.plc4x.java.knxnetip.configuration;

import org.apache.plc4x.java.knxnetip.readwrite.KnxLayer;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.spi.configuration.exceptions.ConfigurationException;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationProvider;

public class KnxNetIpConfiguration implements Configuration, TransportConfigurationProvider {

    @ConfigurationParameter("knxproj-file-path")
    public String knxprojFilePath;

    @ConfigurationParameter("knxproj-password")
    public String knxprojPassword;

    @ConfigurationParameter("group-address-num-levels")
    @IntDefaultValue(3)
    public int groupAddressNumLevels = 3;

    @ConfigurationParameter("connection-type")
    @StringDefaultValue("LINK_LAYER")
    public String connectionType = "LINK_LAYER";

    @ComplexConfigurationParameter(prefix = "udp", defaultOverrides = {}, requiredOverrides = {})
    private KnxNetIpUdpTransportConfiguration udpTransportConfiguration;

    @ComplexConfigurationParameter(prefix = "pcap", defaultOverrides = {}, requiredOverrides = {})
    private KnxNetIpPcapReplayTransportConfiguration pcapReplayTransportConfiguration;

    @ComplexConfigurationParameter(prefix = "raw", defaultOverrides = {}, requiredOverrides = {})
    private KnxNetIpRawSocketTransportConfiguration rawSocketTransportConfiguration;

    public String getKnxprojFilePath() {
        return knxprojFilePath;
    }

    public void setKnxprojFilePath(String knxprojFilePath) {
        this.knxprojFilePath = knxprojFilePath;
    }

    public String getKnxprojPassword() {
        return knxprojPassword;
    }

    public void setKnxprojPassword(String knxprojPassword) {
        this.knxprojPassword = knxprojPassword;
    }

    public int getGroupAddressNumLevels() {
        return groupAddressNumLevels;
    }

    public void setGroupAddressNumLevels(int groupAddressNumLevels) {
        this.groupAddressNumLevels = groupAddressNumLevels;
    }

    public String getConnectionType() {
        return connectionType;
    }

    public void setConnectionType(String connectionType) {
        // Try to parse the provided value, if it doesn't match any of the constants,
        // throw an error.
        try {
            KnxLayer.valueOf("TUNNEL_" + connectionType.toUpperCase());
        } catch (IllegalArgumentException e) {
            throw new ConfigurationException("Value provided for connection-type invalid.");
        }
        this.connectionType = connectionType.toUpperCase();
    }

    public KnxNetIpUdpTransportConfiguration getUdpTransportConfiguration() {
        return udpTransportConfiguration;
    }

    public void setUdpTransportConfiguration(KnxNetIpUdpTransportConfiguration udpTransportConfiguration) {
        this.udpTransportConfiguration = udpTransportConfiguration;
    }

    public KnxNetIpPcapReplayTransportConfiguration getPcapReplayTransportConfiguration() {
        return pcapReplayTransportConfiguration;
    }

    public void setPcapReplayTransportConfiguration(KnxNetIpPcapReplayTransportConfiguration pcapReplayTransportConfiguration) {
        this.pcapReplayTransportConfiguration = pcapReplayTransportConfiguration;
    }

    public KnxNetIpRawSocketTransportConfiguration getRawSocketTransportConfiguration() {
        return rawSocketTransportConfiguration;
    }

    public void setRawSocketTransportConfiguration(KnxNetIpRawSocketTransportConfiguration rawSocketTransportConfiguration) {
        this.rawSocketTransportConfiguration = rawSocketTransportConfiguration;
    }

    @Override
    public TransportConfiguration getTransportConfiguration(String transportCode) {
        switch (transportCode) {
            case "udp":
                return udpTransportConfiguration;
            case "pcap":
                return pcapReplayTransportConfiguration;
            case "raw":
                return rawSocketTransportConfiguration;
        }
        return null;
    }

    @Override
    public String toString() {
        return "Configuration{" +
            "knxprojFilePath=" + knxprojFilePath + ", " +
            "groupAddressNumLevels=" + groupAddressNumLevels +
            '}';
    }

}
