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
package org.apache.plc4x.java.bacnetip.configuration;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.transport.pcap.PcapTransportConfiguration;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.transport.udp.UdpTransportConfiguration;
import org.apache.plc4x.java.utils.pcapsockets.netty.handlers.PacketHandler;
import org.pcap4j.packet.Packet;

public class PassiveBacNetIpConfiguration implements Configuration, UdpTransportConfiguration, RawSocketTransportConfiguration, PcapTransportConfiguration {

    @ConfigurationParameter("ede-file-path")
    public String edeFilePath;

    public String getEdeFilePath() {
        return edeFilePath;
    }

    public void setEdeFilePath(String edeFilePath) {
        this.edeFilePath = edeFilePath;
    }

    @Override
    public int getDefaultPort() {
        return 47808;
    }

    @Override
    public Integer getProtocolId() {
        return null;
    }

    @Override
    public float getReplaySpeedFactor() {
        return 0;
    }

    /**
     * Packet handler to use when running in PCAP mode.
     * In this case all packets are Ethernet frames and we need to first get the
     * IP packet and then the UDP packet and then the raw data from that.
     * @return payload of the packet.
     */
    @Override
    public PacketHandler getPcapPacketHandler() {
        return new PacketHandler() {
            @Override
            public byte[] getData(Packet packet) {
                return packet.getPayload().getPayload().getPayload().getRawData();
            }
        };
    }

}
