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
package org.apache.plc4x.java.utils.pcapsockets.netty;

import io.netty.channel.Channel;
import io.netty.channel.ChannelConfig;
import io.netty.channel.ChannelOption;
import io.netty.channel.DefaultChannelConfig;
import org.apache.plc4x.java.utils.pcapsockets.netty.handlers.PacketHandler;
import org.pcap4j.packet.Packet;

import java.util.Map;

public class PcapSocketChannelConfig extends DefaultChannelConfig implements ChannelConfig {

    public static float SPEED_SLOW_HALF = 2f;
    public static float SPEED_REALTIME = -1f;
    public static float SPEED_FAST_DOUBLE = 0.5f;
    public static float SPEED_FAST_FULL = 0f;

    private int port = PcapSocketAddress.ALL_PORTS;
    private int protocolId = PcapSocketAddress.ALL_PROTOCOLS;
    private float speedFactor = SPEED_REALTIME;
    private PacketHandler packetHandler;

    public PcapSocketChannelConfig(Channel channel) {
        super(channel);
        packetHandler = new PacketHandler() {
            @Override
            public byte[] getData(Packet packet) {
                return packet.getRawData();
            }
        };
    }

    @Override
    public Map<ChannelOption<?>, Object> getOptions() {
        return getOptions(super.getOptions(),
            PcapSocketChannelOption.PORT, PcapSocketChannelOption.PROTOCOL_ID,
            PcapSocketChannelOption.SPEED_FACTOR, PcapSocketChannelOption.PACKET_HANDLER);
    }

    @Override
    public <T> boolean setOption(ChannelOption<T> option, T value) {
        if(option == PcapSocketChannelOption.PORT) {
            if(value instanceof Integer) {
                port = (Integer) value;
                return true;
            }
            return false;
        } else if(option == PcapSocketChannelOption.PROTOCOL_ID) {
            if(value instanceof Integer) {
                protocolId = (Integer) value;
                return true;
            }
            return false;
        } else if(option == PcapSocketChannelOption.SPEED_FACTOR) {
            if(value instanceof Float) {
                speedFactor = (Float) value;
                if(speedFactor > 0) {
                    return true;
                }
            }
            return false;
        } else if(option == PcapSocketChannelOption.PACKET_HANDLER) {
            if(value instanceof PacketHandler) {
                packetHandler = (PacketHandler) value;
                return true;
            }
            return false;
        } else {
            return super.setOption(option, value);
        }
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

    public float getSpeedFactor() {
        return speedFactor;
    }

    public void setSpeedFactor(float speedFactor) {
        this.speedFactor = speedFactor;
    }

    public void setPacketHandler(PacketHandler packetHandler) {
        this.packetHandler = packetHandler;
    }

    public PacketHandler getPacketHandler() {
        return packetHandler;
    }

}
