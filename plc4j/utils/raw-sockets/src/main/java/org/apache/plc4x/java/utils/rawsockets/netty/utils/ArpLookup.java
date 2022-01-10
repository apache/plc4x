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
package org.apache.plc4x.java.utils.rawsockets.netty.utils;

import org.pcap4j.core.*;
import org.pcap4j.packet.ArpPacket;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.namednumber.ArpHardwareType;
import org.pcap4j.packet.namednumber.ArpOperation;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.ByteArrays;
import org.pcap4j.util.MacAddress;

import java.net.InetSocketAddress;
import java.util.Optional;
import java.util.concurrent.*;

public class ArpLookup {

    public static Optional<MacAddress> resolveMacAddress(PcapNetworkInterface nif, InetSocketAddress remoteAddress, InetSocketAddress localAddress, MacAddress localMacAddress) {
        try {
            // This handle will be used for receiving response packets.
            PcapHandle receivingHandle = nif.openLive(
                65535, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 100);
            // This handle will be used for sending the request packet.
            PcapHandle sendingHandle = nif.openLive(
                65535, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 100);
            // The executor, that handles processing the incoming packets.
            ExecutorService arpExecutor = Executors.newSingleThreadExecutor();
            CompletableFuture<MacAddress> remoteMacAddressFuture = new CompletableFuture<>();
            try {
                // Try to limit the number of processed incoming packets to the minimum.
                // So far we know the source host ip as well as the target ip and mac address.
                receivingHandle.setFilter(
                    String.format("arp and src host %s and dst host %s and ether dst %s",
                        Pcaps.toBpfString(remoteAddress.getAddress()), Pcaps.toBpfString(localAddress.getAddress()),
                        Pcaps.toBpfString(localMacAddress)),
                    BpfProgram.BpfCompileMode.OPTIMIZE);

                // Register the listener, which will be processing all packets that pass
                // the filter (Should actually only be one)
                PacketListener listener =
                    packet -> {
                        if (packet.contains(ArpPacket.class)) {
                            ArpPacket arp = packet.get(ArpPacket.class);
                            if (arp.getHeader().getOperation().equals(ArpOperation.REPLY)) {
                                remoteMacAddressFuture.complete(arp.getHeader().getSrcHardwareAddr());
                            }
                        }
                    };

                // The resolution task actually runs in one of the
                // arpExecutor pools threads and just makes sure the
                // incoming packet is passed to the listener.
                Runnable resolutionTask = () -> {
                    try {
                        receivingHandle.loop(1, listener);
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    } catch (PcapNativeException | NotOpenException e) {
                        remoteMacAddressFuture.completeExceptionally(e);
                    }
                };
                arpExecutor.execute(resolutionTask);

                // Actually assemble the ARP packet.
                ArpPacket.Builder arpBuilder = new ArpPacket.Builder();
                arpBuilder.hardwareType(ArpHardwareType.ETHERNET)
                    .protocolType(EtherType.IPV4)
                    .hardwareAddrLength((byte) MacAddress.SIZE_IN_BYTES)
                    .protocolAddrLength((byte) ByteArrays.INET4_ADDRESS_SIZE_IN_BYTES)
                    .operation(ArpOperation.REQUEST)
                    .srcHardwareAddr(localMacAddress)
                    .srcProtocolAddr(localAddress.getAddress())
                    .dstHardwareAddr(MacAddress.ETHER_BROADCAST_ADDRESS)
                    .dstProtocolAddr(remoteAddress.getAddress());
                EthernetPacket.Builder etherBuilder = new EthernetPacket.Builder();
                etherBuilder
                    .dstAddr(MacAddress.ETHER_BROADCAST_ADDRESS)
                    .srcAddr(localMacAddress)
                    .type(EtherType.ARP)
                    .payloadBuilder(arpBuilder)
                    .paddingAtBuild(true);
                Packet arpRequestPacket = etherBuilder.build();

                // Send the arp lookup packet.
                sendingHandle.sendPacket(arpRequestPacket);

                // Wait for the future to complete (It's completed in the packet listener).
                try {
                    return Optional.of(remoteMacAddressFuture.get(1000, TimeUnit.MILLISECONDS));
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } catch (ExecutionException | TimeoutException e) {
                    return Optional.empty();
                }
            } finally {
                // Gracefully shut down.
                if (receivingHandle.isOpen()) {
                    receivingHandle.close();
                }
                if (sendingHandle.isOpen()) {
                    sendingHandle.close();
                }
                if (!arpExecutor.isShutdown()) {
                    arpExecutor.shutdown();
                }
            }
        } catch (NotOpenException | PcapNativeException e) {
            return Optional.empty();
        }
        return Optional.empty();
    }

}
