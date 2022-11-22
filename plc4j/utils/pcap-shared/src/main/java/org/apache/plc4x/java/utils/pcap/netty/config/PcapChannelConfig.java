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
package org.apache.plc4x.java.utils.pcap.netty.config;

import io.netty.channel.Channel;
import io.netty.channel.ChannelConfig;
import io.netty.channel.ChannelOption;
import io.netty.channel.DefaultChannelConfig;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;
import org.pcap4j.core.Pcaps;
import org.pcap4j.packet.Packet;
import org.pcap4j.util.MacAddress;

import java.net.SocketAddress;
import java.util.Map;

public class PcapChannelConfig extends DefaultChannelConfig implements ChannelConfig {

    public static final int ALL_PROTOCOLS = -1;
    public static final int ALL_PORTS = -1;

    private boolean supportVlans = false;
    private int protocolId = ALL_PROTOCOLS;
    private int port = ALL_PORTS;
    private PacketHandler packetHandler = Packet::getRawData;
    private boolean resolveMacAddress = false;

    public PcapChannelConfig(Channel channel) {
        super(channel);
    }

    public PcapChannelConfig clone() {
        PcapChannelConfig clone = new PcapChannelConfig(channel);
        clone.supportVlans = this.supportVlans;
        clone.protocolId = this.protocolId;
        clone.port = this.port;
        clone.packetHandler = this.packetHandler;
        clone.resolveMacAddress = this.resolveMacAddress;
        return clone;
    }

    @Override
    public Map<ChannelOption<?>, Object> getOptions() {
        return getOptions(super.getOptions(),
            PcapChannelOption.SUPPORT_VLANS, PcapChannelOption.PORT, PcapChannelOption.PROTOCOL_ID,
            PcapChannelOption.PACKET_HANDLER, PcapChannelOption.RESOLVE_MAC_ADDRESS);
    }

    @Override
    public <T> boolean setOption(ChannelOption<T> option, T value) {
        if (option == PcapChannelOption.SUPPORT_VLANS) {
            if (value instanceof Boolean) {
                supportVlans = (Boolean) value;
                return true;
            }
            return false;
        } else if (option == PcapChannelOption.PORT) {
            if (value instanceof Integer) {
                port = (Integer) value;
                return true;
            }
            return false;
        } else if (option == PcapChannelOption.PROTOCOL_ID) {
            if (value instanceof Integer) {
                protocolId = (Integer) value;
                return true;
            }
            return false;
        } else if (option == PcapChannelOption.PACKET_HANDLER) {
            if (value instanceof PacketHandler) {
                packetHandler = (PacketHandler) value;
                return true;
            }
            return false;
        } else if (option == PcapChannelOption.RESOLVE_MAC_ADDRESS) {
            if (value instanceof Boolean) {
                resolveMacAddress = (Boolean) value;
                return true;
            }
            return false;
        } else {
            return super.setOption(option, value);
        }
    }

    public boolean isSupportVlans() {
        return supportVlans;
    }

    public int getPort() {
        return port;
    }

    public void setPort(int port) {
        this.port = port;
    }

    public int getProtocolId() {
        return protocolId;
    }

    public void setProtocolId(int protocolId) {
        this.protocolId = protocolId;
    }

    public void setPacketHandler(PacketHandler packetHandler) {
        this.packetHandler = packetHandler;
    }

    public PacketHandler getPacketHandler() {
        return packetHandler;
    }

    public boolean isResolveMacAddress() {
        return resolveMacAddress;
    }

    public String getMacBasedFilterString(MacAddress localMacAddress, MacAddress remoteMacAddress) {
        StringBuilder sb = new StringBuilder();
        if (getProtocolId() != ALL_PROTOCOLS) {
            sb.append(" and (ether proto ").append(getProtocolId()).append(")");
        }
        // Add a filter for TCP or UDP port.
        if (getPort() != ALL_PORTS) {
            sb.append(" and (port ").append(getPort()).append(")");
        }
        // Add a filter for source or target address.
        if(localMacAddress != null) {
            sb.append(" and (ether dst ").append(Pcaps.toBpfString(localMacAddress)).append(")");
        }
        // Add a filter for source or target address.
        if(remoteMacAddress != null) {
            sb.append(" and (ether src ").append(Pcaps.toBpfString(remoteMacAddress)).append(")");
        }
        return (sb.length() > 0) ? sb.substring(" and ".length()) : "";
    }

    public String getMacBasedFilterString(SocketAddress localAddress, SocketAddress remoteAddress) {
        StringBuilder sb = new StringBuilder();
        if (isSupportVlans()) {
            final PcapChannelConfig clone = this.clone();
            clone.supportVlans = false;
            String subFilterString = clone.getMacBasedFilterString(localAddress, remoteAddress);
            if (subFilterString.isEmpty()) {
                sb.append(" and (vlan)");
            } else {
                sb.append(" and ((vlan and ").append(subFilterString).append(") or (").append(subFilterString).append("))");
            }
        } else {
            if (getProtocolId() != ALL_PROTOCOLS) {
                sb.append(" and (ether proto ").append(getProtocolId()).append(")");
            }
            // Add a filter for TCP or UDP port.
            if (getPort() != ALL_PORTS) {
                sb.append(" and (port ").append(getPort()).append(")");
            }
            // Add a filter for source or target address.
            /*if(localAddress != null) {
                sb.append(" and (host ").append(localAddress.getHostAddress()).append(")");
            }
            // Add a filter for source or target address.
            if(remoteAddress != null) {
                sb.append(" and (host ").append(localAddress.getHostAddress()).append(")");
            }*/
        }
        return (sb.length() > 0) ? sb.substring(" and ".length()) : "";
    }

}
