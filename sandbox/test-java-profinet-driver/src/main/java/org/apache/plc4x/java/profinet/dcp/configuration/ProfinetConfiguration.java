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
package org.apache.plc4x.java.profinet.dcp.configuration;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;
import org.pcap4j.packet.Packet;

public class ProfinetConfiguration implements Configuration, RawSocketTransportConfiguration {

    @Override
    public boolean getSupportVlans() {
        return false;
    }

    @Override
    public int getDefaultPort() {
        return -1;
    }

    @Override
    public Integer getProtocolId() {
        return 0x8892;
    }

    @Override
    public PacketHandler getPcapPacketHandler() {
        return new PacketHandler() {
            @Override
            public byte[] getData(Packet packet) {
                // We rely directly on the ethernet frame, so we just need everything.
                return packet.getRawData();
            }
        };
    }
}
