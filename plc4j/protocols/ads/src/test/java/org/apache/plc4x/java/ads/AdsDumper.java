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
package org.apache.plc4x.java.ads;

import org.apache.commons.lang3.RandomStringUtils;
import org.apache.plc4x.java.ads.api.commands.AdsWriteRequest;
import org.apache.plc4x.java.ads.api.commands.types.Data;
import org.apache.plc4x.java.ads.api.commands.types.IndexGroup;
import org.apache.plc4x.java.ads.api.commands.types.IndexOffset;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.tcp.AmsTCPPacket;
import org.apache.plc4x.java.ads.api.tcp.types.UserData;
import org.apache.plc4x.java.mock.connection.tcp.TcpHexDumper;
import org.pcap4j.core.PcapDumper;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.Pcaps;
import org.pcap4j.packet.*;
import org.pcap4j.packet.namednumber.*;
import org.pcap4j.util.MacAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.Socket;
import java.nio.file.Files;
import java.nio.file.Path;

public class AdsDumper {

    private static final Logger logger = LoggerFactory.getLogger(AdsDumper.class);

    public static void main(String... args) throws Exception {
        Path dumpFile = Files.createTempFile("pcapdump", ".pcap");

        try (PcapHandle handle = Pcaps.openDead(DataLinkType.EN10MB, 65536);
             PcapDumper dumper = handle.dumpOpen(dumpFile.toAbsolutePath().toString())) {

            String randomString = RandomStringUtils.randomAscii(1024);

            AdsWriteRequest adsWriteRequest = AdsWriteRequest.of(
                AmsNetId.of("192.168.99.101.1.1"),
                AmsPort.of(851),
                AmsNetId.of("192.168.99.1.1.1"),
                AmsPort.of(14),
                Invoke.of(0),
                IndexGroup.of(1),
                IndexOffset.of(3),
                Data.of(randomString.getBytes())
            );
            AmsTCPPacket amsTCPPacket = AmsTCPPacket.of(UserData.of(adsWriteRequest.getBytes()));

            //try (TcpHexDumper tcpHexDumper = TcpHexDumper.runOn(55862); Socket localhost = new Socket("localhost", tcpHexDumper.getPort())) {
            try (TcpHexDumper tcpHexDumper = TcpHexDumper.runOn(55862); Socket localhost = new Socket("192.168.99.100", 48898)) {
                localhost.getOutputStream().write(amsTCPPacket.getBytes());
            }

            UnknownPacket.Builder amsPacket = new UnknownPacket.Builder();
            amsPacket.rawData(amsTCPPacket.getBytes());

            TcpPacket.Builder tcpPacketBuilder = new TcpPacket.Builder();
            tcpPacketBuilder
                .srcAddr(InetAddress.getLocalHost())
                .srcPort(TcpPort.getInstance((short) 13))
                .dstAddr(InetAddress.getLocalHost())
                .dstPort(TcpPort.getInstance((short) 48898))
                .payloadBuilder(amsPacket)
                .correctChecksumAtBuild(true)
                .correctLengthAtBuild(true);

            IpV4Packet.Builder ipv4PacketBuilder = new IpV4Packet.Builder();
            ipv4PacketBuilder
                .version(IpVersion.IPV4)
                .tos(IpV4Rfc1349Tos.newInstance((byte) 0x75))
                .protocol(IpNumber.TCP)
                .srcAddr((Inet4Address) InetAddress.getLocalHost())
                .dstAddr((Inet4Address) InetAddress.getLocalHost())
                .payloadBuilder(tcpPacketBuilder)
                .correctChecksumAtBuild(true)
                .correctLengthAtBuild(true);

            EthernetPacket.Builder etherPacketBuilder = new EthernetPacket.Builder();
            etherPacketBuilder
                .srcAddr(MacAddress.getByName("fe:00:00:00:00:01"))
                .dstAddr(MacAddress.getByName("fe:00:00:00:00:02"))
                .type(EtherType.IPV4)
                .payloadBuilder(ipv4PacketBuilder)
                .paddingAtBuild(true);

            dumper.dump(etherPacketBuilder.build());
            dumper.flush();

            logger.info("Wrote {}", dumpFile);
        }
    }
}
