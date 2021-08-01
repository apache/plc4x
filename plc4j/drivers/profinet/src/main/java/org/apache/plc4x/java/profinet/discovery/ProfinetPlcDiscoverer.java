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
package org.apache.plc4x.java.profinet.discovery;

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.profinet.ProfinetDriver;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.readwrite.io.EthernetFrameIO;
import org.apache.plc4x.java.profinet.readwrite.types.VirtualLanPriority;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
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
                if (!dev.isLoopBack() && dev.isRunning()) {
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
                                            EthernetFrame ethernetFrame = EthernetFrameIO.staticParse(reader);
                                            DCP_PDU pdu;
                                            // Access the pdu data (either directly or by
                                            // unpacking the content of the VLAN packet.
                                            if (ethernetFrame.getPayload() instanceof VirtualLanEthernetFramePayload) {
                                                VirtualLanEthernetFramePayload vlefpl = (VirtualLanEthernetFramePayload) ethernetFrame.getPayload();
                                                pdu = ((ProfinetEthernetFramePayload) vlefpl.getPayload()).getPdu();
                                            } else {
                                                pdu = ((ProfinetEthernetFramePayload) ethernetFrame.getPayload()).getPdu();
                                            }
                                            // Inspect the PDU itself
                                            // (in this case we only process identify response packets)
                                            if (pdu instanceof DCP_Identify_ResPDU) {
                                                DCP_Identify_ResPDU identifyResPDU = (DCP_Identify_ResPDU) pdu;

                                                Map<String, DCP_Block> blocks = new HashMap<>();
                                                for (DCP_Block block : identifyResPDU.getBlocks()) {
                                                    String blockName = block.getOption().name() + "-" + block.getSuboption().toString();
                                                    blocks.put(blockName, block);
                                                }

                                                // The mac address of the device we found
                                                org.pcap4j.util.MacAddress srcAddr = ethernetPacket.getHeader().getSrcAddr();
                                                // The mac address of the local network device
                                                org.pcap4j.util.MacAddress dstAddr = ethernetPacket.getHeader().getDstAddr();

                                                String deviceTypeName = "unknown";
                                                if (blocks.containsKey(DEVICE_TYPE_NAME)) {
                                                    DCP_BlockDevicePropertiesDeviceVendor block = (DCP_BlockDevicePropertiesDeviceVendor) blocks.get(DEVICE_TYPE_NAME);
                                                    deviceTypeName = new String(block.getDeviceVendorValue());
                                                }
                                                String deviceName = "unknown";
                                                if (blocks.containsKey(DEVICE_NAME_OF_STATION)) {
                                                    DCP_BlockDevicePropertiesNameOfStation block = (DCP_BlockDevicePropertiesNameOfStation) blocks.get(DEVICE_NAME_OF_STATION);
                                                    deviceName = new String(block.getNameOfStation());
                                                }

                                                String transportUrl = srcAddr.toString();
                                                Map<String, String> options =
                                                    Collections.singletonMap("localMacAddress", dstAddr.toString());
                                                String name = deviceTypeName + " - " + deviceName;
                                                PlcDiscoveryItem value = new DefaultPlcDiscoveryItem(
                                                    ProfinetDriver.DRIVER_CODE, RawSocketTransport.TRANSPORT_CODE,
                                                    transportUrl, options, name);
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
                        EthernetFrame identificationRequest = new EthernetFrame(
                            // Pre-Defined PROFINET discovery MAC address
                            new MacAddress(new short[]{0x01, 0x0E, 0xCF, 0x00, 0x00, 0x00}),
                            toPlc4xMacAddress(macAddress),
                            new VirtualLanEthernetFramePayload(VirtualLanPriority.BEST_EFFORT, false, 0,
                                new ProfinetEthernetFramePayload(
                                    new DCP_Identify_ReqPDU(
                                        new ServiceType(false, false),
                                        1,
                                        256,
                                        new DCP_Block[]{
                                            new DCP_BlockALLSelector()
                                        }))));
                        WriteBufferByteBased buffer = new WriteBufferByteBased(34);
                        EthernetFrameIO.staticSerialize(buffer, identificationRequest);
                        Packet packet = EthernetPacket.newPacket(buffer.getData(), 0, 34);
                        handle.sendPacket(packet);
                    }
                }
            }
        } catch (IllegalRawDataException | NotOpenException | PcapNativeException | ParseException e) {
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
        return new MacAddress(new short[]{address[0], address[1], address[2], address[3], address[4], address[5]});
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
