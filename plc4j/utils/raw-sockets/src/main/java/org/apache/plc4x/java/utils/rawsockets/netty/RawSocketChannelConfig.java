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
package org.apache.plc4x.java.utils.rawsockets.netty;

import io.netty.channel.*;
import org.pcap4j.packet.Packet;

import java.util.Map;

public class RawSocketChannelConfig extends DefaultChannelConfig implements ChannelConfig {

    private PacketHandler packetHandler;

    public RawSocketChannelConfig(Channel channel) {
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
        return getOptions(super.getOptions(), RawSocketChannelOption.PACKET_HANDLER);
    }

    @Override
    public <T> boolean setOption(ChannelOption<T> option, T value) {
        if(option == RawSocketChannelOption.PACKET_HANDLER) {
            if(value instanceof PacketHandler) {
                packetHandler = (PacketHandler) value;
                return true;
            }
            return false;
        } else {
            return super.setOption(option, value);
        }
    }

    public void setPacketHandler(PacketHandler packetHandler) {
        this.packetHandler = packetHandler;
    }

    public PacketHandler getPacketHandler() {
        return packetHandler;
    }

}
