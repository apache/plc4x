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
package org.apache.plc4x.java.s7.discovery;

import org.apache.plc4x.java.s7discovery.readwrite.Ethernet_Frame;
import org.apache.plc4x.java.s7discovery.readwrite.Ethernet_FramePayload;
import org.apache.plc4x.java.s7discovery.readwrite.Ethernet_FramePayload_VirtualLan;
import org.apache.plc4x.java.s7discovery.readwrite.MacAddress;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.pcap4j.core.BpfProgram;
import org.pcap4j.core.NotOpenException;
import org.pcap4j.core.PacketListener;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.PcapNativeException;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;
import org.pcap4j.packet.Dot1qVlanTagPacket;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.IllegalRawDataException;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.UdpPacket;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.LinkLayerAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.function.BiConsumer;

/**
 * Raw L2 channel used by the S7 driver's PROFINET DCP-based discovery — opens a pcap4j
 * handle on every non-loopback network interface, filters for PROFINET (EtherType 0x8892
 * including VLAN-tagged variants) plus LLDP (0x88cc), and dispatches decoded
 * {@link Ethernet_FramePayload}s to registered listeners.
 *
 * <p>Stripped down from a fuller PROFINET implementation; only the discovery-relevant code
 * paths remain.
 */
public class ProfinetChannel {

    private final Logger logger = LoggerFactory.getLogger(ProfinetChannel.class);
    private static final EtherType PN_EtherType = EtherType.getInstance((short) 0x8892);
    private static final EtherType LLDP_EtherType = EtherType.getInstance((short) 0x88cc);

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
            try {
                WriteBufferByteBased buffer = newWriteBuffer(ethFrame.getLengthInBytes());
                ethFrame.serialize(buffer);
                Packet packet = EthernetPacket.newPacket(buffer.getBytes(), 0, ethFrame.getLengthInBytes());
                handle.sendPacket(packet);
            } catch (PcapNativeException | NotOpenException | BufferException | IllegalRawDataException e) {
                throw new RuntimeException(e);
            }
        }
    }

    /**
     * Build a {@link WriteBufferByteBased} configured for PROFINET wire encoding —
     * big-endian byte order, plain unsigned-binary / twos-complement / IEEE754 numerics.
     */
    public static WriteBufferByteBased newWriteBuffer(int size) {
        return new WriteBufferByteBased(new byte[size],
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
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
                    logger.debug("Raw socket loop ended", e);
                }
            });
            thread.setDaemon(true);
            thread.setName("ProfinetChannel-RawSocket-" + handle.getFilteringExpression());
            thread.start();
        }
    }

    public PacketListener createListener() {
        return packet -> {
            // EthernetPacket is the highest level of abstraction we can get from pcap4j.
            // Any PROFINET payload inside has to be decoded against the s7discovery mspec.
            if (!(packet instanceof EthernetPacket ethernetPacket)) {
                return;
            }
            boolean isPnPacket = false;
            // Sometimes Ethernet packets are wrapped in a VLAN tag — unwrap and check.
            if (ethernetPacket.getPayload() instanceof Dot1qVlanTagPacket vlanPacket) {
                if (PN_EtherType.equals(vlanPacket.getHeader().getType())
                    || LLDP_EtherType.equals(vlanPacket.getHeader().getType())) {
                    isPnPacket = true;
                }
            } else if (PN_EtherType.equals(ethernetPacket.getHeader().getType())
                || LLDP_EtherType.equals(ethernetPacket.getHeader().getType())) {
                isPnPacket = true;
            } else if (ethernetPacket.getHeader().getType() == EtherType.IPV4
                && ethernetPacket.getPayload().getPayload() instanceof UdpPacket udp) {
                int dstPort = udp.getHeader().getDstPort().value();
                if (dstPort == -30572 || dstPort == -15536 || dstPort == -15535) {
                    isPnPacket = true;
                }
            }

            if (!isPnPacket) {
                return;
            }
            try {
                ReadBuffer reader = new ReadBufferByteBased(ethernetPacket.getRawData(),
                    WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                    WithOption.WithSignedIntegerEncoding("twos-complement"),
                    WithOption.WithFloatEncoding("IEEE754"),
                    WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
                Ethernet_Frame ethernetFrame = Ethernet_Frame.staticParse(reader);
                Ethernet_FramePayload payload = ethernetFrame.getPayload();
                // Unwrap a single VLAN layer if present.
                if (payload instanceof Ethernet_FramePayload_VirtualLan vlan) {
                    payload = vlan.getPayload();
                }
                for (BiConsumer<Ethernet_FramePayload, EthernetPacket> listener : packetListeners) {
                    listener.accept(payload, ethernetPacket);
                }
            } catch (BufferException e) {
                logger.debug("Got error decoding PROFINET packet", e);
            }
        };
    }

    public Map<MacAddress, PcapHandle> getInterfaceHandles(List<PcapNetworkInterface> devs) {
        Map<MacAddress, PcapHandle> handles = new HashMap<>();
        try {
            for (PcapNetworkInterface dev : devs) {
                // Some macOS interfaces without IP addresses fail at filter compile time —
                // skip them up-front rather than catching afterwards.
                if (dev.getAddresses().isEmpty() || dev.isLoopBack()) {
                    continue;
                }
                for (LinkLayerAddress linkLayerAddress : dev.getLinkLayerAddresses()) {
                    org.pcap4j.util.MacAddress macAddress = (org.pcap4j.util.MacAddress) linkLayerAddress;
                    PcapHandle handle = dev.openLive(65536, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
                    handles.put(toPlc4xMacAddress(macAddress), handle);

                    // Filter on PROFINET (0x8892) and the VLAN ethertype (0x8100), restricted
                    // to traffic destined to our local MAC. PROFINET broadcasts arrive with a
                    // multicast destination address that the kernel also delivers here.
                    handle.setFilter(
                        "((ether proto 0x8100) or (ether proto 0x8892)) and (ether dst "
                            + Pcaps.toBpfString(macAddress) + ")",
                        BpfProgram.BpfCompileMode.OPTIMIZE);
                }
            }
        } catch (NotOpenException | PcapNativeException e) {
            logger.error("Got an exception while opening raw sockets", e);
            for (Map.Entry<MacAddress, PcapHandle> entry : handles.entrySet()) {
                try {
                    entry.getValue().breakLoop();
                    entry.getValue().close();
                } catch (NotOpenException ignored) {
                    // Already closed.
                }
            }
        }
        return handles;
    }

    public Map<MacAddress, PcapHandle> getOpenHandles() {
        return openHandles;
    }

    private static MacAddress toPlc4xMacAddress(org.pcap4j.util.MacAddress pcap4jMacAddress) {
        byte[] address = pcap4jMacAddress.getAddress();
        return new MacAddress(new byte[] {
            address[0], address[1], address[2], address[3], address[4], address[5]
        });
    }
}
