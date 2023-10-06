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

package org.apache.plc4x.java.profinet.channel;

import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.generation.*;
import org.pcap4j.core.*;
import org.pcap4j.packet.*;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.LinkLayerAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

public class ProfinetChannel {

    private final Logger logger = LoggerFactory.getLogger(ProfinetChannel.class);
    private static final EtherType PN_EtherType = EtherType.getInstance((short) 0x8892);
    private static final EtherType LLDP_EtherType = EtherType.getInstance((short) 0x88cc);

    private static final EtherType IPv4_EtherType = EtherType.getInstance((short) 0x0800);

    private final Map<MacAddress, PcapHandle> openHandles;
    private final Set<BiConsumer<Ethernet_FramePayload, EthernetPacket>> packetListeners;

    public ProfinetChannel(List<PcapNetworkInterface> devs) {
        this.openHandles = getInterfaceHandles(devs);
        this.packetListeners = new HashSet<>();
        startListener();
    }

    public void addPacketListener(BiConsumer<Ethernet_FramePayload, EthernetPacket> packetListener) {
        this.packetListeners.add(packetListener);
    }

    public void send(Ethernet_Frame ethFrame) {
        for (Map.Entry<MacAddress, PcapHandle> entry : openHandles.entrySet()) {
            PcapHandle handle = entry.getValue();
            WriteBufferByteBased buffer = new WriteBufferByteBased(ethFrame.getLengthInBytes());
            try {
                ethFrame.serialize(buffer);
                Packet packet = EthernetPacket.newPacket(buffer.getBytes(), 0, ethFrame.getLengthInBytes());
                handle.sendPacket(packet);
            } catch (PcapNativeException | NotOpenException | SerializationException | IllegalRawDataException e) {
                throw new RuntimeException(e);
            }
        }
    }

    protected void startListener() {
        for (PcapHandle handle : openHandles.values()) {
            PacketListener listener = createListener();
            Thread thread = new Thread(() -> {
                try {
                    handle.loop(-1, listener);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } catch (PcapNativeException | NotOpenException e) {
                    logger.error("Got error handling raw socket", e);
                }
            });
            thread.setDaemon(true);
            thread.setName("RawSocket handler " + handle.getFilteringExpression());
            thread.start();
            // TODO: This might be leaking the thread after finishing.
        }
    }

    public PacketListener createListener() {
        return packet -> {
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
                } else if (ethernetPacket.getHeader().getType() == EtherType.IPV4 && ethernetPacket.getPayload().getPayload() instanceof UdpPacket) {
                    UdpPacket payload = (UdpPacket) ethernetPacket.getPayload().getPayload();
                    // Check if it's a PROFINET packet
                    if (payload.getHeader().getDstPort().value() == -30572 || payload.getHeader().getDstPort().value() == -15536 || payload.getHeader().getDstPort().value() == -15535) {
                        isPnPacket = true;
                    }
                }

                if (isPnPacket) {
                    ReadBuffer reader = new ReadBufferByteBased(ethernetPacket.getRawData());
                    try {
                        Ethernet_Frame ethernetFrame = Ethernet_Frame.staticParse(reader);
                        Ethernet_FramePayload payload = ethernetFrame.getPayload();

                        // Virtual lan packets simply contain normal payloads, so we simply unpack them
                        if (payload instanceof Ethernet_FramePayload_VirtualLan) {
                            payload = ((Ethernet_FramePayload_VirtualLan) payload).getPayload();
                            // TODO: Possibly we could have multiple levels of VLans ...
                        }

                        // Pass the PN frame to the listeners.
                        for (BiConsumer<Ethernet_FramePayload, EthernetPacket> packetListener : packetListeners) {
                            packetListener.accept(payload, ethernetPacket);
                        }

                        /*if (payload instanceof Ethernet_FramePayload_PnDcp) {
                            PnDcp_Pdu pdu = ((Ethernet_FramePayload_PnDcp) payload).getPdu();
                            if (pdu.getFrameId() == PnDcp_FrameId.DCP_Identify_ResPDU) {
                                if (discoverer != null) {
                                    discoverer.processPnDcp(pdu, ethernetPacket);
                                }
                            } else if (pdu.getFrameId() == PnDcp_FrameId.DCP_GetSet_PDU) {
                                for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
                                    if (Arrays.equals(device.getValue().getDeviceContext().getMacAddress().getAddress(), ethernetFrame.getSource().getAddress())) {
                                        PcDcp_GetSet_Pdu getSetPdu = (PcDcp_GetSet_Pdu) pdu;
                                        device.getValue().handleSetIpAddressResponse(getSetPdu);
                                    }
                                }
                            } else if (pdu.getFrameId() == PnDcp_FrameId.Alarm_Low) {
                                    for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
                                        if (Arrays.equals(device.getValue().getDeviceContext().getMacAddress().getAddress(), ethernetFrame.getSource().getAddress())) {
                                            PnDcp_Pdu_AlarmLow alarmPdu = (PnDcp_Pdu_AlarmLow) pdu;
                                            device.getValue().handleAlarmResponse(alarmPdu);
                                        }
                                    }
                                }
                                else if (pdu.getFrameId() == PnDcp_FrameId.RT_CLASS_1) {
                                    for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
                                        if (device.getValue().getDeviceContext().getMacAddress() == null) {
                                            logger.info("Hurz");
                                        } else if (Arrays.equals(device.getValue().getDeviceContext().getMacAddress().getAddress(), ethernetFrame.getSource().getAddress())) {
                                            PnDcp_Pdu_RealTimeCyclic cyclicPdu = (PnDcp_Pdu_RealTimeCyclic) pdu;
                                            device.getValue().handleRealTimeResponse(cyclicPdu);
                                        }
                                    }
                                }
                        } else if (payload instanceof Ethernet_FramePayload_LLDP) {
                            Lldp_Pdu pdu = ((Ethernet_FramePayload_LLDP) payload).getPdu();
                            if (discoverer != null) {
                                discoverer.processLldp(pdu);
                            }
                        } else if (payload instanceof Ethernet_FramePayload_IPv4) {
                            Ethernet_FramePayload_IPv4 payloadIPv4 = (Ethernet_FramePayload_IPv4) payload;
                            if (payloadIPv4.getPayload().getPacketType() == DceRpc_PacketType.PING) {
                                DceRpc_Packet pingRequest = payloadIPv4.getPayload();
                                // Intercept ping packets that originated from the PN device itself.
                                // TODO: Find out how to react to PING messages.
                                // According to https://pubs.opengroup.org/onlinepubs/9629399/chap12.htm the correct response for us to such a ping message would be a "working" response
                                Ethernet_Frame pingResponse = new Ethernet_Frame(ethernetFrame.getSource(), ethernetFrame.getDestination(),
                                    new Ethernet_FramePayload_IPv4(payloadIPv4.getIdentification(), false, false,
                                        payloadIPv4.getTimeToLive(), payloadIPv4.getDestinationAddress(),
                                        payloadIPv4.getSourceAddress(), payloadIPv4.getDestinationPort(),
                                        payloadIPv4.getSourcePort(), new DceRpc_Packet(DceRpc_PacketType.WORKING,
                                        false, false, false,
                                        IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                                        pingRequest.getObjectUuid(), pingRequest.getInterfaceUuid(), pingRequest.getActivityUuid(),
                                        0L, 0L, DceRpc_Operation.CONNECT, (short) 0, new PnIoCm_Packet_Working())
                                        ));
                                this.send(pingResponse);

                                logger.info("Received PING packet: {}", packet);
                            } else {
                                for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
                                    if (Arrays.equals(device.getValue().getDeviceContext().getMacAddress().getAddress(), ethernetFrame.getSource().getAddress())) {
                                        device.getValue().handleResponse(payloadIPv4);
                                    }
                                }
                            }
                        }*/
                    } catch (ParseException e) {
                        logger.error("Got error decoding packet", e);
                    }
                }
            }
        };
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
                            "(ether proto 0x0800) or (((ether proto 0x8100) or (ether proto 0x8892)) and (ether dst " + Pcaps.toBpfString(macAddress) + ")) or (ether proto 0x88cc)",
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

    public Map<MacAddress, PcapHandle> getOpenHandles() {
        return openHandles;
    }

    private static MacAddress toPlc4xMacAddress(org.pcap4j.util.MacAddress pcap4jMacAddress) {
        byte[] address = pcap4jMacAddress.getAddress();
        return new MacAddress(new byte[]{address[0], address[1], address[2], address[3], address[4], address[5]});
    }

}
