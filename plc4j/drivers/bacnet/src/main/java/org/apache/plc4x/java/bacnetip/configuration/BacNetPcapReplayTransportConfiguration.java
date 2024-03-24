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

import org.apache.plc4x.java.bacnetip.readwrite.BacnetConstants;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.transport.pcapreplay.DefaultPcapReplayTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;
import org.pcap4j.packet.Dot1qVlanTagPacket;

public class BacNetPcapReplayTransportConfiguration extends DefaultPcapReplayTransportConfiguration {

    @ConfigurationParameter("support-vlans")
    @BooleanDefaultValue(false)
    @Description("Enables support for VLans")
    private boolean supportVlans;

    @Override
    public boolean getSupportVlans() {
        return supportVlans;
    }

    public void setSupportVlans(boolean supportVlans) {
        this.supportVlans = supportVlans;
    }

    @Override
    public int getDefaultPort() {
        return BacnetConstants.BACNETUDPDEFAULTPORT;
    }

    /**
     * Packet handler to use when running in PCAP mode.
     * In this case all packets are Ethernet frames and we need to first get the
     * IP packet and then the UDP packet and then the raw data from that.
     *
     * @return payload of the packet.
     */
    @Override
    public PacketHandler getPcapPacketHandler() {
        return packet -> {
            // If it's a VLan packet, we need to go one level deeper.
            if (packet.getPayload() instanceof Dot1qVlanTagPacket) {
                return packet.getPayload().getPayload().getPayload().getPayload().getRawData();
            }
            // This is a normal udp packet.
            else {
                if ((packet.getPayload() != null) && (packet.getPayload().getPayload() != null) && (packet.getPayload().getPayload().getPayload() != null)) {
                    return packet.getPayload().getPayload().getPayload().getRawData();
                }
            }
            return null;
        };
    }

}
