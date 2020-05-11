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

import java.util.stream.Stream;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.profinet.dcp.readwrite.MacAddress;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;
import org.pcap4j.packet.Packet;

public class ProfinetConfiguration implements Configuration, RawSocketTransportConfiguration {

    @Required
    @ConfigurationParameter
    @ParameterConverter(MacAddressConverter.class)
    protected MacAddress sender;

    public void setSender(MacAddress sender) {
        this.sender = sender;
    }

    public MacAddress getSender() {
        return sender;
    }

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

    public static class MacAddressConverter implements ConfigurationParameterConverter<MacAddress> {
        @Override
        public Class<MacAddress> getType() {
            return MacAddress.class;
        }

        @Override
        public MacAddress convert(String value) {
            String[] split = value.split(":");
            short[] segments = ArrayUtils.toPrimitive(Stream.of(split).map(segment -> Integer.parseInt(segment, 16))
                .map(Integer::shortValue).toArray(Short[]::new));

            if (segments.length != 6) {
                throw new IllegalArgumentException("Value " + value + " is not valid MAC address");
            }

            return new MacAddress(segments[0], segments[1], segments[2], segments[3], segments[4], segments[5]);
        }
    }
}
