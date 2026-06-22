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

import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7discovery.readwrite.Ethernet_Frame;
import org.apache.plc4x.java.s7discovery.readwrite.Ethernet_FramePayload;
import org.apache.plc4x.java.s7discovery.readwrite.Ethernet_FramePayload_PnDcp;
import org.apache.plc4x.java.s7discovery.readwrite.Ethernet_FramePayload_VirtualLan;
import org.apache.plc4x.java.s7discovery.readwrite.MacAddress;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Block;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Block_ALLSelector;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Block_DevicePropertiesDeviceId;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Block_DevicePropertiesDeviceRole;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Block_DevicePropertiesDeviceVendor;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Block_DevicePropertiesNameOfStation;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Block_IpParameter;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Pdu;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Pdu_IdentifyReq;
import org.apache.plc4x.java.s7discovery.readwrite.PnDcp_Pdu_IdentifyRes;
import org.apache.plc4x.java.s7discovery.readwrite.VirtualLanPriority;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.functions.PlcDiscoverer;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryResponse;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.pcap4j.core.NotOpenException;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.PcapNativeException;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.IllegalRawDataException;
import org.pcap4j.packet.Packet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.CompletableFuture;

/**
 * S7 device discoverer using PROFINET DCP (Discovery and Configuration Protocol). Sends an
 * {@code IdentifyAll} broadcast on every available network interface and collects responses
 * for a fixed window (10 s by default), filtering for devices that we know support the
 * S7Comm protocol (currently Siemens vendor ID {@code 0x002A}, S7-1200 device IDs
 * {@code 0x010D} / {@code 0x010E}; extend as more device types are validated).
 *
 * <p>Requires libpcap on the host (and the {@code pcap4j} runtime dependency on the
 * classpath). Loopback and IP-less interfaces are skipped.
 */
public class S7PlcDiscoverer implements PlcDiscoverer {

    private static final Logger LOGGER = LoggerFactory.getLogger(S7PlcDiscoverer.class);

    // The PROFINET DCP block names follow the {option}-{suboption} convention.
    public static final String DEVICE_TYPE_NAME       = "DEVICE_PROPERTIES_OPTION-1";
    public static final String DEVICE_NAME_OF_STATION = "DEVICE_PROPERTIES_OPTION-2";
    public static final String DEVICE_ID              = "DEVICE_PROPERTIES_OPTION-3";
    public static final String DEVICE_ROLE            = "DEVICE_PROPERTIES_OPTION-4";
    public static final String IP_OPTION_IP           = "IP_OPTION-2";

    // Pre-defined PROFINET multicast MAC for discovery broadcasts.
    private static final MacAddress PROFINET_BROADCAST_MAC_ADDRESS =
        new MacAddress(new byte[] { 0x01, 0x0E, (byte) 0xCF, 0x00, 0x00, 0x00 });

    private static final long DEFAULT_DISCOVERY_TIMEOUT_MS = 10_000L;

    private final ProfinetChannel channel;
    private final List<PlcDiscoveryItem> values = new ArrayList<>();
    private PlcDiscoveryItemHandler handler;

    /** Open a discoverer using all non-loopback interfaces with link-layer addresses. */
    public S7PlcDiscoverer() throws PcapNativeException {
        this(new ProfinetChannel(Pcaps.findAllDevs()));
    }

    /** Test-friendly constructor; pre-built channel allows wiring a fake for unit tests. */
    public S7PlcDiscoverer(ProfinetChannel channel) {
        this.channel = channel;
        channel.addPacketListener(this::handleIncomingPacket);
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return discoverWithHandler(discoveryRequest, null);
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest,
                                                                       PlcDiscoveryItemHandler handler) {
        this.handler = handler;
        sendPnDcpDiscoveryRequest();
        return setDiscoveryEndTimer(discoveryRequest, DEFAULT_DISCOVERY_TIMEOUT_MS);
    }

    public void sendPnDcpDiscoveryRequest() {
        for (Map.Entry<MacAddress, PcapHandle> entry : channel.getOpenHandles().entrySet()) {
            MacAddress localMacAddress = entry.getKey();
            PcapHandle handle = entry.getValue();

            Ethernet_Frame identificationRequest = new Ethernet_Frame(
                PROFINET_BROADCAST_MAC_ADDRESS,
                localMacAddress,
                new Ethernet_FramePayload_VirtualLan(VirtualLanPriority.BEST_EFFORT, false, (short) 0,
                    new Ethernet_FramePayload_PnDcp(
                        new PnDcp_Pdu_IdentifyReq(0xFEFE,
                            1L,
                            256,
                            Collections.singletonList(new PnDcp_Block_ALLSelector())))));

            try {
                WriteBufferByteBased buffer = ProfinetChannel.newWriteBuffer(identificationRequest.getLengthInBytes());
                identificationRequest.serialize(buffer);
                Packet packet = EthernetPacket.newPacket(buffer.getBytes(), 0, identificationRequest.getLengthInBytes());
                handle.sendPacket(packet);
            } catch (BufferException | IllegalRawDataException | NotOpenException e) {
                throw new RuntimeException("Failed to send PROFINET DCP discovery request", e);
            } catch (PcapNativeException e) {
                // "Network is down" happens when e.g. the Wi-Fi adapter has no carrier — skip.
                if (e.getMessage() == null || !e.getMessage().contains("Network is down")) {
                    throw new RuntimeException(e);
                }
            }
        }
    }

    public CompletableFuture<PlcDiscoveryResponse> setDiscoveryEndTimer(PlcDiscoveryRequest discoveryRequest, long delay) {
        CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();
        Timer timer = new Timer("S7-Discovery-Timeout", true);
        timer.schedule(new TimerTask() {
            public void run() {
                PlcDiscoveryResponse response = new DefaultPlcDiscoveryResponse(
                    discoveryRequest, PlcResponseCode.OK, values);
                for (Map.Entry<MacAddress, PcapHandle> entry : channel.getOpenHandles().entrySet()) {
                    PcapHandle openHandle = entry.getValue();
                    try {
                        openHandle.breakLoop();
                        openHandle.close();
                    } catch (Exception e) {
                        LOGGER.debug("Error closing pcap handle on discovery end", e);
                    }
                }
                timer.cancel();
                timer.purge();
                future.complete(response);
            }
        }, delay);
        return future;
    }

    protected void handleIncomingPacket(Ethernet_FramePayload frame, EthernetPacket ethernetPacket) {
        if (frame instanceof Ethernet_FramePayload_PnDcp pnDcp) {
            PnDcp_Pdu pdu = pnDcp.getPdu();
            // FrameId 0xFEFF identifies a PROFINET DCP IdentifyRes; ignore everything else.
            if (pdu.getFrameIdValue() == 0xFEFF) {
                handlePnDcpPacket(pdu, ethernetPacket);
            }
        }
    }

    public void handlePnDcpPacket(PnDcp_Pdu pdu, EthernetPacket ethernetPacket) {
        if (!(pdu instanceof PnDcp_Pdu_IdentifyRes identifyResPDU)) {
            return;
        }

        Map<String, PnDcp_Block> blocks = new HashMap<>();
        for (PnDcp_Block block : identifyResPDU.getBlocks()) {
            String blockName = block.getOption().name() + "-" + (block.getSuboption() & 0xFF);
            blocks.put(blockName, block);
        }

        org.pcap4j.util.MacAddress srcAddr = ethernetPacket.getHeader().getSrcAddr();    // device's MAC
        org.pcap4j.util.MacAddress dstAddr = ethernetPacket.getHeader().getDstAddr();    // our local MAC

        String deviceTypeName = "unknown";
        if (blocks.get(DEVICE_TYPE_NAME) instanceof PnDcp_Block_DevicePropertiesDeviceVendor block) {
            deviceTypeName = new String(block.getDeviceVendorValue());
        }

        String deviceName = "unknown";
        if (blocks.get(DEVICE_NAME_OF_STATION) instanceof PnDcp_Block_DevicePropertiesNameOfStation block) {
            deviceName = new String(block.getNameOfStation());
        }

        String role = decodeDeviceRole(blocks.get(DEVICE_ROLE));

        IpInfo ip = decodeIp(blocks.get(IP_OPTION_IP));
        Map<String, String> options = Collections.emptyMap();
        if ("0.0.0.0".equals(ip.address)) {
            // Some devices ship without an IP and expect to be assigned one via PN-IO. Use
            // the MAC as the connection key so callers can at least identify the device.
            ip = new IpInfo(srcAddr.toString(), ip.subnetMask);
            options = Collections.singletonMap("ip-address", "{some-ip-address}");
        }

        // Vendor + Device IDs determine whether we recognise this as an S7Comm-capable PLC.
        String vendorId = "unknown";
        String deviceId = "unknown";
        if (blocks.get(DEVICE_ID) instanceof PnDcp_Block_DevicePropertiesDeviceId block) {
            vendorId = String.format("%04X", block.getVendorId());
            deviceId = String.format("%04X", block.getDeviceId());
        }
        if (!supportsS7Comm(vendorId, deviceTypeName)) {
            return;
        }

        Map<String, PlcValue> attributes = new HashMap<>();
        attributes.put("ipAddress",       new PlcSTRING(ip.address));
        attributes.put("subnetMask",      new PlcSTRING(ip.subnetMask));
        attributes.put("macAddress",      new PlcSTRING(srcAddr.toString()));
        attributes.put("localMacAddress", new PlcSTRING(dstAddr.toString()));
        attributes.put("deviceTypeName",  new PlcSTRING(deviceTypeName));
        attributes.put("deviceName",      new PlcSTRING(deviceName));
        attributes.put("vendorId",        new PlcSTRING(vendorId));
        attributes.put("deviceId",        new PlcSTRING(deviceId));
        attributes.put("role",            new PlcSTRING(role));
        attributes.put("packetType",      new PlcSTRING("dcp"));

        String displayName = deviceTypeName + " - " + deviceName;
        PlcDiscoveryItem item = new DefaultPlcDiscoveryItem(
            "s7", "cotp", ip.address, options, displayName, attributes);
        values.add(item);
        if (handler != null) {
            handler.handle(item);
        }
        LOGGER.debug("Discovered S7 device: '{}' at '{}'", item.getName(), item.getConnectionUrl());
    }

    private static String decodeDeviceRole(PnDcp_Block block) {
        if (!(block instanceof PnDcp_Block_DevicePropertiesDeviceRole role)) {
            return "unknown";
        }
        StringBuilder sb = new StringBuilder();
        if (role.getPnioSupervisor()) sb.append(",SUPERVISOR");
        if (role.getPnioMultidevive()) sb.append(",MULTIDEVICE");
        if (role.getPnioController()) sb.append(",CONTROLLER");
        if (role.getPnioDevice()) sb.append(",DEVICE");
        return sb.length() == 0 ? "unknown" : sb.substring(1);
    }

    private static IpInfo decodeIp(PnDcp_Block block) {
        if (!(block instanceof PnDcp_Block_IpParameter ip)) {
            return new IpInfo("unknown", "unknown");
        }
        try {
            return new IpInfo(
                InetAddress.getByAddress(ip.getIpAddress()).getHostAddress(),
                InetAddress.getByAddress(ip.getSubnetMask()).getHostAddress());
        } catch (UnknownHostException e) {
            return new IpInfo("invalid", "invalid");
        }
    }

    /**
     * Filter to Siemens products that speak S7Comm on ISO-on-TCP port 102. The PROFINET
     * DCP {@code DEVICE_PROPERTIES_OPTION-1} block carries the human-readable product
     * family ("S7-1200", "S7-300", "LOGO!", etc.) — matching on a small prefix list lets
     * us include all S7 PLC families without device-ID lookups while excluding Siemens
     * non-PLC products (Simocode motor managers, ET200 IO modules, switches, etc.) that
     * share the vendor ID but don't accept S7Comm connections.
     *
     * <p>Recognised product families (case-insensitive prefix match):
     * <ul>
     *   <li>S7-300, S7-400, S7-1200, S7-1500 — all match {@code "S7-"}.</li>
     *   <li>LOGO! 8.x with PROFINET — matches {@code "LOGO"}.</li>
     * </ul>
     */
    private static boolean supportsS7Comm(String vendorId, String deviceTypeName) {
        if (!"002A".equals(vendorId) || deviceTypeName == null) {
            return false;
        }
        String upper = deviceTypeName.toUpperCase(java.util.Locale.ROOT);
        return upper.startsWith("S7-") || upper.startsWith("LOGO");
    }

    private record IpInfo(String address, String subnetMask) {
    }
}
