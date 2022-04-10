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

import org.apache.commons.net.util.SubnetUtils;
import org.pcap4j.core.*;
import org.pcap4j.packet.ArpPacket;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.namednumber.ArpHardwareType;
import org.pcap4j.packet.namednumber.ArpOperation;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.ByteArrays;
import org.pcap4j.util.MacAddress;

import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.UnknownHostException;
import java.util.*;
import java.util.concurrent.*;
import java.util.stream.Collectors;

public class ArpUtils {

    public static Set<InetAddress> scanNetworkDevice(PcapNetworkInterface nif) {
        Set<InetAddress> foundAddresses = new HashSet<>();
        try{
            // Calculate all ip addresses, this device can reach.
            Map<String, List<String>> addresses = new HashMap<>();
            for (PcapAddress address : nif.getAddresses()) {
                if(address instanceof PcapIpV4Address) {
                    final PcapIpV4Address ipV4Address = (PcapIpV4Address) address;
                    SubnetUtils su = new SubnetUtils(ipV4Address.getAddress().getHostAddress(), ipV4Address.getNetmask().getHostAddress());
                    final String currentAddress = ipV4Address.getAddress().getHostAddress();
                    final List<String> reachableAddresses = new ArrayList<>(Arrays.asList(su.getInfo().getAllAddresses()));
                    // Remove the current address from the list.
                    reachableAddresses.remove(currentAddress);
                    addresses.put(currentAddress, reachableAddresses);
                }
            }
            // If this device doesn't have any addresses, abort.
            if(addresses.isEmpty()) {
                return Collections.emptySet();
            }

            final Optional<MacAddress> first = nif.getLinkLayerAddresses().stream()
                .filter(linkLayerAddress -> linkLayerAddress instanceof MacAddress)
                .map(linkLayerAddress -> (MacAddress) linkLayerAddress).findFirst();
            // If we couldn't find a local mac address, abort.
            //noinspection SimplifyOptionalCallChains (Not compatible with Java 8)
            if(!first.isPresent()) {
                return Collections.emptySet();
            }
            final MacAddress localMacAddress = first.get();

            // This handle will be used for receiving response packets.
            PcapHandle receivingHandle = nif.openLive(
                65535, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 100);
            // This handle will be used for sending the request packet.
            PcapHandle sendingHandle = nif.openLive(
                65535, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 100);
            // The executor, that handles processing the incoming packets.
            ExecutorService arpExecutor = Executors.newSingleThreadExecutor();
            try {
                StringBuilder sb = new StringBuilder("arp");
                sb.append(" and ether dst ").append(Pcaps.toBpfString(localMacAddress)).append(" and (");
                boolean firstAddress = true;
                for (String localAddress : addresses.keySet()) {
                    if(!firstAddress) {
                        sb.append(" or ");
                    }
                    sb.append("(dst host ").append(localAddress).append(")");
                    firstAddress = false;
                }
                sb.append(")");
                // In this case we simply accept any arp traffic.
                receivingHandle.setFilter(sb.toString(), BpfProgram.BpfCompileMode.OPTIMIZE);

                // Register the listener for all incoming arp packets.
                PacketListener listener =
                    packet -> {
                        if (packet.contains(ArpPacket.class)) {
                            ArpPacket arp = packet.get(ArpPacket.class);
                            if (arp.getHeader().getOperation().equals(ArpOperation.REPLY)) {
                                foundAddresses.add(arp.getHeader().getSrcProtocolAddr());
                            }
                        }
                    };

                // The resolution task actually runs in one of the
                // arpExecutor pools threads and just makes sure the
                // incoming packet is passed to the listener.
                Runnable resolutionTask = () -> {
                    try {
                        while (receivingHandle.isOpen()) {
                            final Packet nextPacket = receivingHandle.getNextPacket();
                            if(nextPacket != null) {
                                listener.gotPacket(nextPacket);
                            }
                        }
                    } catch (NotOpenException e) {
                        // Ignore.
                    }
                };
                arpExecutor.execute(resolutionTask);

                for (Map.Entry<String, List<String>> stringListEntry : addresses.entrySet()) {
                    InetAddress localAddress = InetAddress.getByName(stringListEntry.getKey());
                    List<String> remoteAddresses = stringListEntry.getValue();
                    for (String remoteAddressString : remoteAddresses) {
                        InetAddress remoteAddress = InetAddress.getByName(remoteAddressString);
                        ArpPacket.Builder arpBuilder = new ArpPacket.Builder();
                        arpBuilder.hardwareType(ArpHardwareType.ETHERNET)
                            .protocolType(EtherType.IPV4)
                            .hardwareAddrLength((byte) MacAddress.SIZE_IN_BYTES)
                            .protocolAddrLength((byte) ByteArrays.INET4_ADDRESS_SIZE_IN_BYTES)
                            .operation(ArpOperation.REQUEST)
                            .srcHardwareAddr(localMacAddress)
                            .srcProtocolAddr(localAddress)
                            .dstHardwareAddr(MacAddress.ETHER_BROADCAST_ADDRESS)
                            .dstProtocolAddr(remoteAddress);
                        EthernetPacket.Builder etherBuilder = new EthernetPacket.Builder();
                        etherBuilder
                            .dstAddr(MacAddress.ETHER_BROADCAST_ADDRESS)
                            .srcAddr(localMacAddress)
                            .type(EtherType.ARP)
                            .payloadBuilder(arpBuilder)
                            .paddingAtBuild(true);
                        Packet arpRequestPacket = etherBuilder.build();

                        //System.out.println("Checking " + remoteAddressString + " from ip " + stringListEntry.getKey());
                        // Send the arp lookup packet.
                        sendingHandle.sendPacket(arpRequestPacket);
                    }
                }

                try {
                    Thread.sleep(200);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            } catch (UnknownHostException e) {
                e.printStackTrace();
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
            return Collections.emptySet();
        }
        return foundAddresses;
    }

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

    public static void main(String[] args) throws Exception {
        for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
            final Set<InetAddress> inetAddresses = scanNetworkDevice(dev);
            final List<Integer> inetAddresses1 = inetAddresses.stream().map(address -> (short) address.getAddress()[3] & 0xFF).sorted().collect(Collectors.toList());
            System.out.printf("Found %d ip addresses with device %s:\n  %s%n", inetAddresses.size(), dev, inetAddresses1);
        }
    }

}
