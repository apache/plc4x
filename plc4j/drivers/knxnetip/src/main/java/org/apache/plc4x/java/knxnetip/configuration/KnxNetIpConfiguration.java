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
package org.apache.plc4x.java.knxnetip.configuration;

import org.apache.plc4x.java.knxnetip.KnxNetIpDriver;
import org.apache.plc4x.java.knxnetip.readwrite.KnxLayer;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.FloatDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.spi.configuration.exceptions.ConfigurationException;
import org.apache.plc4x.java.transport.pcapreplay.PcapReplayTransportConfiguration;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.transport.udp.UdpTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.config.PcapChannelConfig;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;

public class KnxNetIpConfiguration implements Configuration, UdpTransportConfiguration, PcapReplayTransportConfiguration, RawSocketTransportConfiguration {

    @ConfigurationParameter("knxproj-file-path")
    public String knxprojFilePath;

    @ConfigurationParameter("group-address-num-levels")
    @IntDefaultValue(3)
    public int groupAddressNumLevels = 3;

    @ConfigurationParameter("connection-type")
    @StringDefaultValue("LINK_LAYER")
    public String connectionType = "LINK_LAYER";

    @ConfigurationParameter("replay-speed-factor")
    @FloatDefaultValue(1.0f)
    public float replaySpeedFactor = 1.0f;

    @ConfigurationParameter("loop")
    @BooleanDefaultValue(false)
    public boolean loop = false;

    public String getKnxprojFilePath() {
        return knxprojFilePath;
    }

    public void setKnxprojFilePath(String knxprojFilePath) {
        this.knxprojFilePath = knxprojFilePath;
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

    @Override
    public float getReplaySpeedFactor() {
        return replaySpeedFactor;
    }

    public void setReplaySpeedFactor(float replaySpeedFactor) {
        this.replaySpeedFactor = replaySpeedFactor;
    }

    @Override
    public boolean isLoop() {
        return loop;
    }

    public void setLoop(boolean loop) {
        this.loop = loop;
    }

    @Override
    public int getDefaultPort() {
        return KnxNetIpDriver.KNXNET_IP_PORT;
    }

    @Override
    public Integer getProtocolId() {
        return PcapChannelConfig.ALL_PROTOCOLS;
    }

    @Override
    public PacketHandler getPcapPacketHandler() {
        return packet -> packet.getPayload().getPayload().getPayload().getRawData();
    }

    @Override
    public String toString() {
        return "Configuration{" +
            "knxprojFilePath=" + knxprojFilePath + ", " +
            "groupAddressNumLevels=" + groupAddressNumLevels +
            '}';
    }

}
