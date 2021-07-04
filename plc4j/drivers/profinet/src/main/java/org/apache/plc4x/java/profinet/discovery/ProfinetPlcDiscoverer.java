/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.profinet.discovery;

import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.readwrite.io.EthernetFrameIO;
import org.apache.plc4x.java.profinet.readwrite.types.VirtualLanPriority;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.apache.plc4x.java.spi.messages.PlcDiscoverer;
import org.pcap4j.core.*;
import org.pcap4j.packet.Dot1qVlanTagPacket;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.IllegalRawDataException;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.namednumber.EtherType;
import org.pcap4j.util.LinkLayerAddress;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class ProfinetPlcDiscoverer implements PlcDiscoverer {

    private static final EtherType PN_EtherType = EtherType.getInstance((short) 0x8892);

    private static final String DEVICE_TYPE_NAME = "DEVICE_PROPERTIES_OPTION-1";
    private static final String DEVICE_NAME_OF_STATION = "DEVICE_PROPERTIES_OPTION-2";
    private static final String DEVICE_ID = "DEVICE_PROPERTIES_OPTION-3";
    private static final String DEVICE_ROLE = "DEVICE_PROPERTIES_OPTION-4";
    private static final String DEVICE_OPTIONS = "DEVICE_PROPERTIES_OPTION-5";
    private static final String DEVICE_INSTANCE = "DEVICE_PROPERTIES_OPTION-7";
    private static final String IP_OPTION_IP = "IP_OPTION-2";

    public ProfinetPlcDiscoverer() {
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        Map<String, DCP_Identify_ResPDU> pnDevices = new HashMap<>();
        try {
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                if(!dev.isLoopBack() && dev.isRunning()) {
                    for (LinkLayerAddress linkLayerAddress : dev.getLinkLayerAddresses()) {
                        org.pcap4j.util.MacAddress macAddress = (org.pcap4j.util.MacAddress) linkLayerAddress;
                        PcapHandle handle = dev.openLive(65536, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
                        PcapHandle sendHandle = dev.openLive(65536, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
                        ExecutorService pool = Executors.newSingleThreadExecutor();

                        // Only react on PROFINET DCP packets targeted at our current MAC address.
                        // TODO: Find out how to filter based on the ether frame type ...
                        handle.setFilter(
                            "ether dst " + Pcaps.toBpfString(macAddress),
                            BpfProgram.BpfCompileMode.OPTIMIZE);

                        PacketListener listener =
                            packet -> {
                                if(packet instanceof EthernetPacket) {
                                    EthernetPacket ethernetPacket = (EthernetPacket) packet;
                                    boolean isPnPacket = false;
                                    if(ethernetPacket.getPayload() instanceof Dot1qVlanTagPacket) {
                                        Dot1qVlanTagPacket vlanPacket = (Dot1qVlanTagPacket) ethernetPacket.getPayload();
                                        if(PN_EtherType.equals(vlanPacket.getHeader().getType())) {
                                            isPnPacket = true;
                                        }
                                    } else if(PN_EtherType.equals(ethernetPacket.getHeader().getType())) {
                                        isPnPacket = true;
                                    }

                                    // It's a PROFINET packet.
                                    if(isPnPacket) {
                                        ReadBuffer reader = new ReadBufferByteBased(ethernetPacket.getRawData());
                                        try {
                                            EthernetFrame ethernetFrame = EthernetFrameIO.staticParse(reader);
                                            String sourceMacAddress = toMacAddressString(ethernetFrame.getSource());
                                            DCP_PDU pdu;
                                            if(ethernetFrame.getPayload() instanceof VirtualLanEthernetFramePayload) {
                                                VirtualLanEthernetFramePayload vlefpl = (VirtualLanEthernetFramePayload) ethernetFrame.getPayload();
                                                pdu = ((ProfinetEthernetFramePayload) vlefpl.getPayload()).getPdu();
                                            } else {
                                                pdu = ((ProfinetEthernetFramePayload) ethernetFrame.getPayload()).getPdu();
                                            }
                                            if(pdu instanceof DCP_Identify_ResPDU) {
                                                DCP_Identify_ResPDU identify_resPDU = (DCP_Identify_ResPDU) pdu;
                                                if(!pnDevices.containsKey(sourceMacAddress)) {
                                                    pnDevices.put(sourceMacAddress, identify_resPDU);
                                                }
                                            }
                                        } catch (ParseException e) {
                                            e.printStackTrace();
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
                        sendHandle.sendPacket(packet);
                    }
                }
            }
        } catch (PcapNativeException | ParseException e) {
            e.printStackTrace();
        } catch (NotOpenException e) {
            e.printStackTrace();
        } catch (IllegalRawDataException e) {
            e.printStackTrace();
        }
        try {
            Thread.sleep(5000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        System.out.println(String.format("Found %d PROFINET devices:", pnDevices.size()));
        for (DCP_Identify_ResPDU pnDevice : pnDevices.values()) {
            outputPnDevice(pnDevice);
        }

        return null;
    }

    private static MacAddress toPlc4xMacAddress(org.pcap4j.util.MacAddress pcap4jMacAddress) {
        byte[] address = pcap4jMacAddress.getAddress();
        return new MacAddress(new short[]{ (short) address[0], (short) address[1], (short) address[2], (short) address[3], (short) address[4], (short) address[5]});
    }

    private static String toMacAddressString(MacAddress macAddress) {
        return String.format("%x2:%x2:%x2:%x2:%x2:%x2", macAddress.getAddress()[0], macAddress.getAddress()[1],
            macAddress.getAddress()[2], macAddress.getAddress()[3], macAddress.getAddress()[4], macAddress.getAddress()[5]);
    }

    private static void outputPnDevice(DCP_Identify_ResPDU pnDevice) {
        Map<String, DCP_Block> blocks = new HashMap<>();
        for (DCP_Block block : pnDevice.getBlocks()) {
            String blockName = block.getOption().name() + "-" + block.getSuboption().toString();
            blocks.put(blockName, block);
        }

        String deviceTypeName = "unknown";
        if(blocks.containsKey(DEVICE_TYPE_NAME)) {
            DCP_BlockDevicePropertiesDeviceVendor block = (DCP_BlockDevicePropertiesDeviceVendor) blocks.get(DEVICE_TYPE_NAME);
            deviceTypeName = new String(block.getDeviceVendorValue());
        }
        String deviceName = "unknown";
        if(blocks.containsKey(DEVICE_NAME_OF_STATION)) {
            DCP_BlockDevicePropertiesNameOfStation block = (DCP_BlockDevicePropertiesNameOfStation) blocks.get(DEVICE_NAME_OF_STATION);
            deviceName = new String(block.getNameOfStation());
        }
        String ipAddress = "unknown";
        if(blocks.containsKey(IP_OPTION_IP)) {
            DCP_BlockIpIpParameter block = (DCP_BlockIpIpParameter) blocks.get(IP_OPTION_IP);
            ipAddress = String.format("%d.%d.%d.%d", block.getIpAddress()[0], block.getIpAddress()[1], block.getIpAddress()[2], block.getIpAddress()[3]);
        }
        System.out.println(String.format("Found '%s' with name '%s' on IP: %s%n\t%s", deviceTypeName, deviceName, ipAddress, pnDevice));
    }

    private static class Task implements Runnable {

        private PcapHandle handle;
        private PacketListener listener;

        public Task(PcapHandle handle, PacketListener listener) {
            this.handle = handle;
            this.listener = listener;
        }

        @Override
        public void run() {
            try {
                handle.loop(10, listener);
            } catch (PcapNativeException e) {
                e.printStackTrace();
            } catch (InterruptedException e) {
                e.printStackTrace();
            } catch (NotOpenException e) {
                e.printStackTrace();
            }
        }
    }

    public static void main(String[] args) throws Exception {
        ProfinetPlcDiscoverer discoverer = new ProfinetPlcDiscoverer();
        discoverer.discover(null);

        Thread.sleep(10000);
    }

}
