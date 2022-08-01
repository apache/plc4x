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

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.profinet.ProfinetDriver;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryResponse;
import org.apache.plc4x.java.spi.messages.PlcDiscoverer;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransport;
import org.pcap4j.core.*;
import org.pcap4j.packet.Dot1qVlanTagPacket;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.IllegalRawDataException;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.LinkLayerAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class ProfinetPlcDiscoverer implements PlcDiscoverer {

    private static final EtherType PN_EtherType = EtherType.getInstance((short) 0x8892);

    // The constants for the different block names and their actual meaning.
    private static final String DEVICE_TYPE_NAME = "DEVICE_PROPERTIES_OPTION-1";
    private static final String DEVICE_NAME_OF_STATION = "DEVICE_PROPERTIES_OPTION-2";
    private static final String DEVICE_ID = "DEVICE_PROPERTIES_OPTION-3";
    private static final String DEVICE_ROLE = "DEVICE_PROPERTIES_OPTION-4";
    private static final String DEVICE_OPTIONS = "DEVICE_PROPERTIES_OPTION-5";
    private static final String DEVICE_INSTANCE = "DEVICE_PROPERTIES_OPTION-7";
    private static final String IP_OPTION_IP = "IP_OPTION-2";

    private final Logger logger = LoggerFactory.getLogger(ProfinetPlcDiscoverer.class);

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return discoverWithHandler(discoveryRequest, null);
    }

    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest, PlcDiscoveryItemHandler handler) {
        CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();
        Set<PcapHandle> openHandles = new HashSet<>();
        List<PlcDiscoveryItem> values = new ArrayList<>();
        try {
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                if (!dev.isLoopBack()) {
                    for (LinkLayerAddress linkLayerAddress : dev.getLinkLayerAddresses()) {
                        org.pcap4j.util.MacAddress macAddress = (org.pcap4j.util.MacAddress) linkLayerAddress;
                        PcapHandle handle = dev.openLive(65536, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
                        openHandles.add(handle);

                        ExecutorService pool = Executors.newSingleThreadExecutor();

                        // Only react on PROFINET DCP packets targeted at our current MAC address.
                        handle.setFilter(
                            "((ether proto 0x8100) or (ether proto 0x8892)) and (ether dst " + Pcaps.toBpfString(macAddress) + ")",
                            BpfProgram.BpfCompileMode.OPTIMIZE);

                        PacketListener listener =
                            packet -> {
                                // EthernetPacket is the highest level of abstraction we can be expecting.
                                // Everything inside this we will have to decode ourselves.
                                if (packet instanceof EthernetPacket) {
                                    EthernetPacket ethernetPacket = (EthernetPacket) packet;
                                    boolean isPnPacket = false;
                                    // I have observed some times the ethernet packets being wrapped inside a VLAN
                                    // Packet, in this case we simply unpack the content.
                                    if (ethernetPacket.getPayload() instanceof Dot1qVlanTagPacket) {
                                        Dot1qVlanTagPacket vlanPacket = (Dot1qVlanTagPacket) ethernetPacket.getPayload();
                                        if (PN_EtherType.equals(vlanPacket.getHeader().getType())) {
                                            isPnPacket = true;
                                        }
                                    } else if (PN_EtherType.equals(ethernetPacket.getHeader().getType())) {
                                        isPnPacket = true;
                                    }

                                    // It's a PROFINET packet.
                                    if (isPnPacket) {
                                        ReadBuffer reader = new ReadBufferByteBased(ethernetPacket.getRawData());
                                        try {
                                            Ethernet_Frame ethernetFrame = Ethernet_Frame.staticParse(reader);
                                            PnDcp_Pdu pdu;
                                            // Access the pdu data (either directly or by
                                            // unpacking the content of the VLAN packet.
                                            if (ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                                                Ethernet_FramePayload_VirtualLan vlefpl = (Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload();
                                                pdu = ((Ethernet_FramePayload_PnDcp) vlefpl.getPayload()).getPdu();
                                            } else {
                                                pdu = ((Ethernet_FramePayload_PnDcp) ethernetFrame.getPayload()).getPdu();
                                            }
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

                                                // Get the Vendor Id and the Device Id
                                                String vendorId = "unknown";
                                                String deviceId = "unknown";
                                                if (blocks.containsKey(DEVICE_ID)) {
                                                    PnDcp_Block_DevicePropertiesDeviceId block = (PnDcp_Block_DevicePropertiesDeviceId) blocks.get(DEVICE_ID);
                                                    vendorId = String.format("%04X", block.getVendorId());
                                                    deviceId = String.format("%04X", block.getDeviceId());
                                                }

                                                Map<String, String> options = new HashMap<>();
                                                options.put("remoteIpAddress", remoteIpAddress);
                                                options.put("remoteSubnetMask", remoteSubnetMask);
                                                options.put("remoteMacAddress", srcAddr.toString());
                                                options.put("localMacAddress", dstAddr.toString());
                                                options.put("deviceTypeName", deviceTypeName);
                                                options.put("deviceName", deviceName);
                                                options.put("vendorId", vendorId);
                                                options.put("deviceId", deviceId);
                                                options.put("role", role);
                                                String name = deviceTypeName + " - " + deviceName;
                                                PlcDiscoveryItem value = new DefaultPlcDiscoveryItem(
                                                    ProfinetDriver.DRIVER_CODE, RawSocketTransport.TRANSPORT_CODE,
                                                    remoteIpAddress, options, name, Collections.emptyMap());
                                                values.add(value);

                                                // If we have a discovery handler, pass it to the handler callback
                                                if (handler != null) {
                                                    handler.handle(value);
                                                }

                                                logger.debug("Found new device: '{}' with connection-url '{}'",
                                                    value.getName(), value.getConnectionUrl());
                                            }
                                        } catch (ParseException e) {
                                            logger.error("Got error decoding packet", e);
                                        }
                                    }
                                }
                            };
                        Task t = new Task(handle, listener);
                        pool.execute(t);

                        // Construct and send the search request.
                        Ethernet_Frame identificationRequest = new Ethernet_Frame(
                            // Pre-Defined PROFINET discovery MAC address
                            new MacAddress(new byte[]{0x01, 0x0E, (byte) 0xCF, 0x00, 0x00, 0x00}),
                            toPlc4xMacAddress(macAddress),
                            new Ethernet_FramePayload_VirtualLan(VirtualLanPriority.BEST_EFFORT, false, 0,
                                new Ethernet_FramePayload_PnDcp(
                                    new PnDcp_Pdu_IdentifyReq(PnDcp_FrameId.DCP_Identify_ReqPDU.getValue(),
                                        1,
                                        256,
                                        Collections.singletonList(
                                            new PnDcp_Block_ALLSelector()
                                        )))));
                        WriteBufferByteBased buffer = new WriteBufferByteBased(34);
                        identificationRequest.serialize(buffer);
                        Packet packet = EthernetPacket.newPacket(buffer.getData(), 0, 34);
                        handle.sendPacket(packet);
                    }
                }
            }
        } catch (IllegalRawDataException | NotOpenException | PcapNativeException | SerializationException e) {
            logger.error("Got an exception while processing raw socket data", e);
            future.completeExceptionally(new PlcException("Got an internal error while performing discovery"));
            for (PcapHandle openHandle : openHandles) {
                openHandle.close();
            }
            return future;
        }

        // Create a timer that completes the future after a given time with all the responses it found till then.
        Timer timer = new Timer("Discovery Timeout");
        timer.schedule(new TimerTask() {
            public void run() {
                PlcDiscoveryResponse response =
                    new DefaultPlcDiscoveryResponse(discoveryRequest, PlcResponseCode.OK, values);
                future.complete(response);
                for (PcapHandle openHandle : openHandles) {
                    openHandle.close();
                }
            }
        }, 5000L);

        return future;
    }

    private static MacAddress toPlc4xMacAddress(org.pcap4j.util.MacAddress pcap4jMacAddress) {
        byte[] address = pcap4jMacAddress.getAddress();
        return new MacAddress(new byte[]{address[0], address[1], address[2], address[3], address[4], address[5]});
    }

    private static class Task implements Runnable {

        private final Logger logger = LoggerFactory.getLogger(Task.class);

        private final PcapHandle handle;
        private final PacketListener listener;

        public Task(PcapHandle handle, PacketListener listener) {
            this.handle = handle;
            this.listener = listener;
        }

        @Override
        public void run() {
            try {
                handle.loop(10, listener);
            } catch (InterruptedException e) {
                logger.error("Got error handling raw socket", e);
                Thread.currentThread().interrupt();
            } catch (PcapNativeException | NotOpenException e) {
                logger.error("Got error handling raw socket", e);
            }
        }
    }

    public static void main(String[] args) throws Exception {
        ProfinetPlcDiscoverer discoverer = new ProfinetPlcDiscoverer();
        discoverer.discover(null);

        Thread.sleep(10000);
    }

}
