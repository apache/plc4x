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

package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.discovery.ProfinetPlcDiscoverer;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.pcap4j.core.*;
import org.pcap4j.packet.Dot1qVlanTagPacket;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.LinkLayerAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.function.Function;


public class ProfinetChannel {

    private final Logger logger = LoggerFactory.getLogger(ProfinetChannel.class);
    private static final EtherType PN_EtherType = EtherType.getInstance((short) 0x8892);
    private static final EtherType LLDP_EtherType = EtherType.getInstance((short) 0x88cc);
    private ProfinetPlcDiscoverer discoverer = null;
    private ProfinetConfiguration configuration = null;
    private Map<MacAddress, PcapHandle> openHandles;

    public ProfinetChannel(List<PcapNetworkInterface> devs) {
        this.openHandles = getInterfaceHandles(devs);
        startListener();
    }

    public void startListener() {
        for (Map.Entry<MacAddress, PcapHandle> entry : openHandles.entrySet()) {
            PcapHandle handle = entry.getValue();

            Function<Object, Boolean> runtimeHandler =
                message -> {
                    PacketListener listener = createListener();
                    try {
                        handle.loop(-1, listener);
                    } catch (InterruptedException e) {
                        logger.error("Got error handling raw socket", e);
                        Thread.currentThread().interrupt();
                    } catch (PcapNativeException | NotOpenException e) {
                        logger.error("Got error handling raw socket", e);
                    }
                    return null;
                };

            Thread thread = new Thread(new ProfinetRunnable(handle, runtimeHandler));
            thread.start();
        }
    }

    private static class ProfinetRunnable implements Runnable {

        private final Logger logger = LoggerFactory.getLogger(ProfinetRunnable.class);
        private final PcapHandle handle;
        private final Function<Object, Boolean> operator;

        public ProfinetRunnable(PcapHandle handle, Function<Object, Boolean> operator) {
            this.handle = handle;
            this.operator = operator;
        }

        @Override
        public void run() {
            operator.apply(null);
        }

    }

    public PacketListener createListener() {
        PacketListener listener =
            packet -> {
                // EthernetPacket is the highest level of abstraction we can be expecting.
                // Everything inside this we will have to decode ourselves.
                if (packet instanceof EthernetPacket) {
                    EthernetPacket ethernetPacket = (EthernetPacket) packet;
                    boolean isPnPacket = false;
                    // I have observed sometimes the ethernet packets being wrapped inside a VLAN
                    // Packet, in this case we simply unpack the content.
                    if (ethernetPacket.getPayload() instanceof Dot1qVlanTagPacket) {
                        Dot1qVlanTagPacket vlanPacket = (Dot1qVlanTagPacket) ethernetPacket.getPayload();
                        if (PN_EtherType.equals(vlanPacket.getHeader().getType()) || LLDP_EtherType.equals(vlanPacket.getHeader().getType())) {
                            isPnPacket = true;
                        }
                    } else if (PN_EtherType.equals(ethernetPacket.getHeader().getType()) || LLDP_EtherType.equals(ethernetPacket.getHeader().getType())) {
                        isPnPacket = true;
                    }

                    if (isPnPacket) {
                        ReadBuffer reader = new ReadBufferByteBased(ethernetPacket.getRawData());
                        try {
                            Ethernet_Frame ethernetFrame = Ethernet_Frame.staticParse(reader);
                            Ethernet_FramePayload payload = ethernetFrame.getPayload();
                            if (payload instanceof Ethernet_FramePayload_VirtualLan) {
                                payload = ((Ethernet_FramePayload_VirtualLan) payload).getPayload();
                            }

                            if (ethernetFrame.getPayload() instanceof Ethernet_FramePayload_PnDcp) {
                                PnDcp_Pdu pdu = ((Ethernet_FramePayload_PnDcp) payload).getPdu();
                                if (discoverer != null) {
                                    discoverer.processPnDcp(pdu, ethernetPacket);
                                }
                            } else if (ethernetFrame.getPayload() instanceof Ethernet_FramePayload_LLDP) {
                                Lldp_Pdu pdu = ((Ethernet_FramePayload_LLDP) payload).getPdu();
                                if (discoverer != null) {
                                    discoverer.processLldp(pdu);
                                }
                            } else if (payload instanceof Ethernet_FramePayload_IPv4) {
                                logger.debug("Udp Packet Found");
                            }

                        } catch (ParseException e) {
                            logger.error("Got error decoding packet", e);
                        }
                    }
                }
            };
        return listener;
    }

    public Map<MacAddress, PcapHandle> getInterfaceHandles(List<PcapNetworkInterface> devs) {
        Map<MacAddress, PcapHandle> openHandles = new HashMap<>();
        try {
            for (PcapNetworkInterface dev : devs) {
                // It turned out on some MAC network devices without any ip addresses
                // the compiling of the filter expression was causing errors. As
                // currently there was no other way to detect this, this check seems
                // to be sufficient.
                if (dev.getAddresses().size() == 0) {
                    continue;
                }
                if (!dev.isLoopBack()) {
                    for (LinkLayerAddress linkLayerAddress : dev.getLinkLayerAddresses()) {
                        org.pcap4j.util.MacAddress macAddress = (org.pcap4j.util.MacAddress) linkLayerAddress;
                        PcapHandle handle = dev.openLive(65536, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
                        openHandles.put(toPlc4xMacAddress(macAddress), handle);

                        // Only react on PROFINET, UDP or LLDP packets targeted at our current MAC address.
                        handle.setFilter(
                            "(((ether proto 0x8100) or (ether proto 0x8892)) and (ether dst " + Pcaps.toBpfString(macAddress) + ")) or (ether proto 0x88cc) or (ether proto 0x0800)",
                            BpfProgram.BpfCompileMode.OPTIMIZE);
                    }
                }
            }
        } catch (NotOpenException | PcapNativeException e) {
            logger.error("Got an exception while processing raw socket data", e);
            for (Map.Entry<MacAddress, PcapHandle> entry : openHandles.entrySet()) {
                PcapHandle openHandle = entry.getValue();
                try {
                    openHandle.breakLoop();
                    openHandle.close();
                } catch (NotOpenException error) {
                    logger.info("Handle already closed.");
                }
            }
        }
        return openHandles;
    }

    public ProfinetPlcDiscoverer getDiscoverer() {
        return discoverer;
    }

    public void setDiscoverer(ProfinetPlcDiscoverer discoverer) {
        this.discoverer = discoverer;
    }

    public void setConfiguration(ProfinetConfiguration configuration) {
        this.configuration = configuration;
    }

    public ProfinetConfiguration getConfiguration() {
        return configuration;
    }

    public Map<MacAddress, PcapHandle> getOpenHandles() {
        return openHandles;
    }

    private static MacAddress toPlc4xMacAddress(org.pcap4j.util.MacAddress pcap4jMacAddress) {
        byte[] address = pcap4jMacAddress.getAddress();
        return new MacAddress(new byte[]{address[0], address[1], address[2], address[3], address[4], address[5]});
    }

}
