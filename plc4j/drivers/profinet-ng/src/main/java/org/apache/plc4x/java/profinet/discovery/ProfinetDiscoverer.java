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
package org.apache.plc4x.java.profinet.discovery;

import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.profinet.ProfinetDriver;
import org.apache.plc4x.java.profinet.channel.ProfinetChannel;
import org.apache.plc4x.java.profinet.packets.PnDcpPacketFactory;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryResponse;
import org.apache.plc4x.java.spi.messages.PlcDiscoverer;
import org.apache.plc4x.java.spi.values.PlcValues;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransport;
import org.pcap4j.core.NotOpenException;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.PcapNativeException;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.IllegalRawDataException;
import org.pcap4j.packet.Packet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.*;
import java.util.concurrent.CompletableFuture;

public class ProfinetDiscoverer implements PlcDiscoverer {

    // The constants for the different block names and their actual meaning.
    public static final String DEVICE_TYPE_NAME = "DEVICE_PROPERTIES_OPTION-1";
    public static final String DEVICE_NAME_OF_STATION = "DEVICE_PROPERTIES_OPTION-2";
    public static final String DEVICE_ID = "DEVICE_PROPERTIES_OPTION-3";
    public static final String DEVICE_ROLE = "DEVICE_PROPERTIES_OPTION-4";
    public static final String DEVICE_OPTIONS = "DEVICE_PROPERTIES_OPTION-5";
    public static final String DEVICE_INSTANCE = "DEVICE_PROPERTIES_OPTION-7";
    public static final String IP_OPTION_IP = "IP_OPTION-2";
    // Pre-Defined PROFINET discovery MAC address
    private static final MacAddress PROFINET_BROADCAST_MAC_ADDRESS = new MacAddress(new byte[]{0x01, 0x0E, (byte) 0xCF, 0x00, 0x00, 0x00});
    final private ProfinetChannel channel;
    final List<PlcDiscoveryItem> values = new ArrayList<>();
    final Set<Timer> periodicTimers = new HashSet<>();
    private final Logger logger = LoggerFactory.getLogger(ProfinetDiscoverer.class);
    private PlcDiscoveryItemHandler handler;

    public ProfinetDiscoverer(ProfinetChannel channel) {
        this.channel = channel;
        channel.addPacketListener(this::handleIncomingPacket);
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return discoverWithHandler(discoveryRequest, null);
    }

    public CompletableFuture<PlcDiscoveryResponse> setDiscoveryEndTimer(PlcDiscoveryRequest discoveryRequest, long delay) {
        CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();

        // Create a timer that completes the future after a given time with all the responses it found till then.
        Timer timer = new Timer("Discovery Timeout");
        timer.schedule(new TimerTask() {
            public void run() {
                PlcDiscoveryResponse response =
                    new DefaultPlcDiscoveryResponse(discoveryRequest, PlcResponseCode.OK, values);
                for (Map.Entry<MacAddress, PcapHandle> entry : channel.getOpenHandles().entrySet()) {
                    PcapHandle openHandle = entry.getValue();
                    try {
                        openHandle.breakLoop();
                        openHandle.close();
                    } catch (Exception e) {
                        logger.error("Error occurred while closing handle");
                    }
                }
                for (Timer timer : periodicTimers) {
                    timer.cancel();
                    timer.purge();
                }
                future.complete(response);
            }
        }, delay);

        return future;
    }

    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest, PlcDiscoveryItemHandler handler) {
        this.handler = handler;
        sendPnDcpDiscoveryRequest();
        return setDiscoveryEndTimer(discoveryRequest, 10000L);
    }

    public void sendPnDcpDiscoveryRequest() {
        for (Map.Entry<MacAddress, PcapHandle> entry : channel.getOpenHandles().entrySet()) {
            MacAddress localMacAddress = entry.getKey();
            PcapHandle handle = entry.getValue();

            // Construct and send the search request.
            Ethernet_Frame identificationRequest = PnDcpPacketFactory.createIdentificationRequest(localMacAddress, PROFINET_BROADCAST_MAC_ADDRESS);
            WriteBufferByteBased buffer = new WriteBufferByteBased(identificationRequest.getLengthInBytes());
            try {
                identificationRequest.serialize(buffer);
            } catch (SerializationException e) {
                throw new RuntimeException(e);
            }
            try {
                Packet packet = EthernetPacket.newPacket(buffer.getBytes(), 0, identificationRequest.getLengthInBytes());
                handle.sendPacket(packet);
            } catch (PcapNativeException | NotOpenException | IllegalRawDataException e) {
                throw new RuntimeException(e);
            }
        }
    }

    protected void handleIncomingPacket(Ethernet_FramePayload frame, EthernetPacket ethernetPacket) {
        if (frame instanceof Ethernet_FramePayload_PnDcp) {
            PnDcp_Pdu pdu = ((Ethernet_FramePayload_PnDcp) frame).getPdu();
            if (pdu.getFrameId() == PnDcp_FrameId.DCP_Identify_ResPDU) {
                handlePnDcpPacket(pdu, ethernetPacket);
            }
        }
    }

    public void handlePnDcpPacket(PnDcp_Pdu pdu, EthernetPacket ethernetPacket) {
        // Inspect the PDU itself
        // (in this case we only process identify response packets)
        if (pdu instanceof PnDcp_Pdu_IdentifyRes) {
            PnDcp_Pdu_IdentifyRes identifyResPDU = (PnDcp_Pdu_IdentifyRes) pdu;

            Map<String, PnDcp_Block> blocks = new HashMap<>();
            for (PnDcp_Block block : identifyResPDU.getBlocks()) {
                String blockName = block.getOption().name() + "-" + block.getSuboption().toString();
                blocks.put(blockName, block);
            }

            // The mac address of the device we found
            org.pcap4j.util.MacAddress srcAddr = ethernetPacket.getHeader().getSrcAddr();
            // The mac address of the local network device
            org.pcap4j.util.MacAddress dstAddr = ethernetPacket.getHeader().getDstAddr();

            String deviceTypeName = "unknown";
            if (blocks.containsKey(DEVICE_TYPE_NAME)) {
                PnDcp_Block_DevicePropertiesDeviceVendor block = (PnDcp_Block_DevicePropertiesDeviceVendor) blocks.get(DEVICE_TYPE_NAME);
                deviceTypeName = new String(block.getDeviceVendorValue());
            }

            String deviceName = "unknown";
            if (blocks.containsKey(DEVICE_NAME_OF_STATION)) {
                PnDcp_Block_DevicePropertiesNameOfStation block = (PnDcp_Block_DevicePropertiesNameOfStation) blocks.get(DEVICE_NAME_OF_STATION);
                deviceName = new String(block.getNameOfStation());
            }

            String role = "unknown";
            if (blocks.containsKey(DEVICE_ROLE)) {
                role = "";
                PnDcp_Block_DevicePropertiesDeviceRole block = (PnDcp_Block_DevicePropertiesDeviceRole) blocks.get(DEVICE_ROLE);
                if (block.getPnioSupervisor()) {
                    role += ",SUPERVISOR";
                }
                if (block.getPnioMultidevive()) {
                    role += ",MULTIDEVICE";
                }
                if (block.getPnioController()) {
                    role += ",CONTROLLER";
                }
                if (block.getPnioDevice()) {
                    role += ",DEVICE";
                }
                // Cut off the first comma
                if (role.length() > 0) {
                    role = role.substring(1);
                } else {
                    role = "unknown";
                }
            }

            String remoteIpAddress = "unknown";
            String remoteSubnetMask = "unknown";
            if (blocks.containsKey(IP_OPTION_IP)) {
                PnDcp_Block_IpParameter block = (PnDcp_Block_IpParameter) blocks.get(IP_OPTION_IP);
                try {
                    InetAddress addr = InetAddress.getByAddress(block.getIpAddress());
                    remoteIpAddress = addr.getHostAddress();
                    InetAddress netMask = InetAddress.getByAddress(block.getSubnetMask());
                    remoteSubnetMask = netMask.getHostAddress();
                } catch (UnknownHostException e) {
                    remoteIpAddress = "invalid";
                }
            }

            // Get the Vendor-Id and the Device-Id
            String vendorId = "unknown";
            String deviceId = "unknown";
            if (blocks.containsKey(DEVICE_ID)) {
                PnDcp_Block_DevicePropertiesDeviceId block = (PnDcp_Block_DevicePropertiesDeviceId) blocks.get(DEVICE_ID);
                vendorId = String.format("%04X", block.getVendorId());
                deviceId = String.format("%04X", block.getDeviceId());
            }

            Map<String, PlcValue> attributes = new HashMap<>();
            attributes.put("ipAddress", PlcValues.of(remoteIpAddress));
            attributes.put("subnetMask", PlcValues.of(remoteSubnetMask));
            attributes.put("macAddress", PlcValues.of(srcAddr.toString()));
            attributes.put("localMacAddress", PlcValues.of(dstAddr.toString()));
            attributes.put("deviceTypeName", PlcValues.of(deviceTypeName));
            attributes.put("deviceName", PlcValues.of(deviceName));
            attributes.put("vendorId", PlcValues.of(vendorId));
            attributes.put("deviceId", PlcValues.of(deviceId));
            attributes.put("role", PlcValues.of(role));
            attributes.put("packetType", PlcValues.of("dcp"));

            String name = deviceTypeName + " - " + deviceName;

            PlcDiscoveryItem value = new DefaultPlcDiscoveryItem(
                ProfinetDriver.DRIVER_CODE, RawSocketTransport.TRANSPORT_CODE,
                remoteIpAddress, Collections.emptyMap(), name, attributes);
            values.add(value);

            // If we have a discovery handler, pass it to the handler callback
            if (handler != null) {
                handler.handle(value);
            }

            logger.debug("Found new device: '{}' with connection-url '{}'",
                value.getName(), value.getConnectionUrl());
        }
    }

}