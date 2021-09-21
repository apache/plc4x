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
package org.apache.plc4x.java.bacnetip.configuration;

import org.apache.plc4x.java.bacnetip.BacNetIpDriver;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.DoubleDefaultValue;
import org.apache.plc4x.java.transport.pcapreplay.PcapReplayTransportConfiguration;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.transport.udp.UdpTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;
import org.pcap4j.packet.Dot1qVlanTagPacket;

public class BacNetIpConfiguration implements Configuration, UdpTransportConfiguration, RawSocketTransportConfiguration, PcapReplayTransportConfiguration {

    // Path to a single EDE file.
    @ConfigurationParameter("ede-file-path")
    private String edeFilePath;

    // Path to a directory containing many EDE files.
    @ConfigurationParameter("ede-directory-path")
    private String edeDirectoryPath;

    // The speed in which the pcap file is replayed:
    // - 1.0 being the original speed
    // - 0   being as fast as possible (no delays between the packets)
    // - 0.5 being double speed
    // - 2.0 being half speed
    @ConfigurationParameter("pcap-replay-speed")
    @DoubleDefaultValue(1.0F)
    private double pcapReplaySpeed;

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

    public void setPcapReplaySpeed(double pcapReplaySpeed) {
        this.pcapReplaySpeed = pcapReplaySpeed;
    }

    @Override
    public float getReplaySpeedFactor() {
        return (float) pcapReplaySpeed;
    }

    @Override
    public boolean getSupportVlans() {
        return true;
    }

    @Override
    public int getDefaultPort() {
        return BacNetIpDriver.BACNET_IP_PORT;
    }

    @Override
    public Integer getProtocolId() {
        return null;
    }


    /**
     * Packet handler to use when running in PCAP mode.
     * In this case all packets are Ethernet frames and we need to first get the
     * IP packet and then the UDP packet and then the raw data from that.
     * @return payload of the packet.
     */
    @Override
    public PacketHandler getPcapPacketHandler() {
        return packet -> {
            // If it's a VLan packet, we need to go one level deeper.
            if(packet.getPayload() instanceof Dot1qVlanTagPacket) {
                return packet.getPayload().getPayload().getPayload().getPayload().getRawData();
            }
            // This is a normal udp packet.
            else {
                if((packet.getPayload() != null) && (packet.getPayload().getPayload() != null) && (packet.getPayload().getPayload().getPayload() != null)) {
                    return packet.getPayload().getPayload().getPayload().getRawData();
                }
            }
            return null;
        };
    }

}
