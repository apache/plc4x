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

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.profinet.ProfinetDriver;
import org.apache.plc4x.java.profinet.device.ProfinetChannel;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryResponse;
import org.apache.plc4x.java.spi.messages.PlcDiscoverer;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransport;
import org.pcap4j.core.*;
import org.pcap4j.packet.*;
import org.pcap4j.packet.namednumber.EtherType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.function.Function;

public class ProfinetPlcDiscoverer implements PlcDiscoverer {

    // The constants for the different block names and their actual meaning.
    private static final String DEVICE_TYPE_NAME = "DEVICE_PROPERTIES_OPTION-1";
    private static final String DEVICE_NAME_OF_STATION = "DEVICE_PROPERTIES_OPTION-2";
    private static final String PLC4X_LLDP_IDENTIFIER = "PLC4X PROFINET Controller Client";
    private static final String PLC4X_LLDP_PORT = "port001.plc4x";
    private static final String DEVICE_ID = "DEVICE_PROPERTIES_OPTION-3";
    private static final String DEVICE_ROLE = "DEVICE_PROPERTIES_OPTION-4";
    private static final String DEVICE_OPTIONS = "DEVICE_PROPERTIES_OPTION-5";
    private static final String DEVICE_INSTANCE = "DEVICE_PROPERTIES_OPTION-7";
    private static final String IP_OPTION_IP = "IP_OPTION-2";
    // Pre-Defined PROFINET discovery MAC address
    private static final MacAddress PROFINET_BROADCAST_MAC_ADDRESS = new MacAddress(new byte[]{0x01, 0x0E, (byte) 0xCF, 0x00, 0x00, 0x00});
    // Pre-Defined LLDP discovery MAC address
    private static final MacAddress LLDP_BROADCAST_MAC_ADDRESS = new MacAddress(new byte[]{0x01, (byte) 0x80, (byte) 0xc2, 0x00, 0x00, 0x0e});

    Map<MacAddress, PcapHandle> openHandles;
    List<PlcDiscoveryItem> values = new ArrayList<>();
    Set<Timer> periodicTimers = new HashSet<>();
    private final Logger logger = LoggerFactory.getLogger(ProfinetPlcDiscoverer.class);
    private PlcDiscoveryItemHandler handler;

    public ProfinetPlcDiscoverer(ProfinetChannel channel) {
        this.openHandles = channel.getOpenHandles();
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
                for (Map.Entry<MacAddress, PcapHandle> entry : openHandles.entrySet()) {
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
        startLldpPoll(5000L);
        startPnDcpPoll(30000L);
        CompletableFuture<PlcDiscoveryResponse> future = setDiscoveryEndTimer(discoveryRequest, 10000L);
        return future;
    }

    public void ongoingDiscoverWithHandler(PlcDiscoveryRequest discoveryRequest, PlcDiscoveryItemHandler handler, long lldpPeriod, long dcpPeriod) {
        this.handler = handler;
        startLldpPoll(lldpPeriod);
        startPnDcpPoll(dcpPeriod);
    }

    public void processPnDcp(PnDcp_Pdu pdu, EthernetPacket ethernetPacket) {
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
                deviceTypeName = new String(block.getDeviceVendorValue()).replace(" ", "%20");
            }

            String deviceName = "unknown";
            if (blocks.containsKey(DEVICE_NAME_OF_STATION)) {
                PnDcp_Block_DevicePropertiesNameOfStation block = (PnDcp_Block_DevicePropertiesNameOfStation) blocks.get(DEVICE_NAME_OF_STATION);
                deviceName = new String(block.getNameOfStation()).replace(" ", "%20");
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
            options.put("ipAddress", remoteIpAddress);
            options.put("subnetMask", remoteSubnetMask);
            options.put("macAddress", srcAddr.toString());
            options.put("localMacAddress", dstAddr.toString());
            options.put("deviceTypeName", deviceTypeName);
            options.put("deviceName", deviceName);
            options.put("vendorId", vendorId);
            options.put("deviceId", deviceId);
            options.put("role", role);
            options.put("packetType", "dcp");
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
    }

    public void processLldp(Lldp_Pdu pdu) {

        Map<String, String> options = new HashMap<>();
        boolean profibusDevice = false;

        List<LldpUnit> units = pdu.getLldpParameters();
        for (LldpUnit unit : units) {
            if (unit instanceof TlvPortId) {
                TlvPortId portIdPacket = (TlvPortId) unit;
                options.put("portId", portIdPacket.getPortId());
                if (portIdPacket.getPortId().contains(".")) {
                    options.put("deviceName", portIdPacket.getPortId().split("\\.")[1]);
                } else {
                    options.put("deviceName", portIdPacket.getPortId());
                }
            } else if (unit instanceof TlvChassisId) {
                TlvChassisId chassisIdPacket = (TlvChassisId) unit;
                options.put("chassisId", chassisIdPacket.getChassisId());
            } else if (unit instanceof TlvManagementAddress) {
                TlvManagementAddress managementAddressPacket = (TlvManagementAddress) unit;
                try {
                    String ipAddress = InetAddress.getByAddress(managementAddressPacket.getIpAddress().getData()).getHostAddress();
                    options.put("ipAddress", ipAddress);
                } catch (UnknownHostException e) {
                    throw new RuntimeException(e);
                }

            } else if (unit instanceof TlvOrganizationSpecific) {
                TlvOrganizationSpecific orgSpecific = (TlvOrganizationSpecific) unit;
                if (orgSpecific.getOrganizationSpecificUnit().getUniqueCode() == 0x000ECF) {
                    TlvOrgSpecificProfibus specificProfibus = (TlvOrgSpecificProfibus) orgSpecific.getOrganizationSpecificUnit();
                    TlvOrgSpecificProfibusUnit specificProfibusUnit = specificProfibus.getSpecificUnit();
                    switch (specificProfibusUnit.getSubType()) {
                        case CHASSIS_MAC:
                            TlvProfibusSubTypeChassisMac chassisMac = (TlvProfibusSubTypeChassisMac) specificProfibusUnit;
                            options.put("macAddress", Hex.encodeHexString(chassisMac.getMacAddress().getAddress()));
                            profibusDevice = true;
                            break;
                        case PORT_STATUS:
                            break;
                    }
                }
            }
        }

        options.put("packetType", "lldp");

        if (profibusDevice) {
            PlcDiscoveryItem value = new DefaultPlcDiscoveryItem(
                ProfinetDriver.DRIVER_CODE, RawSocketTransport.TRANSPORT_CODE,
                "lldp_response_packet", options, options.get("portId"), Collections.emptyMap());
            values.add(value);

            // If we have a discovery handler, pass it to the handler callback
            if (handler != null) {
                handler.handle(value);
            }

            logger.debug("Found new device: '{}' via an LLDP broardcast",
                options.get("portId"));
        }
    }

    public void startPnDcpPoll(long period) {
        for (Map.Entry<MacAddress, PcapHandle> entry : openHandles.entrySet()) {
            PcapHandle handle = entry.getValue();
            MacAddress macAddress = entry.getKey();
            // Construct and send the search request.

            Function<Object, Boolean> pnDcpTimer =
                message -> {
                    Ethernet_Frame identificationRequest = new Ethernet_Frame(
                        PROFINET_BROADCAST_MAC_ADDRESS,
                        macAddress,
                        new Ethernet_FramePayload_VirtualLan(VirtualLanPriority.BEST_EFFORT, false, 0,
                            new Ethernet_FramePayload_PnDcp(
                                new PnDcp_Pdu_IdentifyReq(PnDcp_FrameId.DCP_Identify_ReqPDU.getValue(),
                                    1,
                                    256,
                                    Collections.singletonList(
                                        new PnDcp_Block_ALLSelector()
                                    )))));
                    WriteBufferByteBased buffer = new WriteBufferByteBased(identificationRequest.getLengthInBytes());
                    try {
                        identificationRequest.serialize(buffer);
                    } catch (SerializationException e) {
                        throw new RuntimeException(e);
                    }
                    Packet packet = null;
                    try {
                        packet = EthernetPacket.newPacket(buffer.getBytes(), 0, identificationRequest.getLengthInBytes());
                    } catch (IllegalRawDataException e) {
                        throw new RuntimeException(e);
                    }
                    try {
                        handle.sendPacket(packet);
                    } catch (PcapNativeException | NotOpenException e) {
                        throw new RuntimeException(e);
                    }
                    return null;
                };

            Timer timer = new Timer();
            periodicTimers.add(timer);

            // Schedule to run after every 3 second(3000 millisecond)
            timer.scheduleAtFixedRate(
                new PeriodicTask(handle, pnDcpTimer),
                0,
                period);
        }
    }

    public void startLldpPoll(long period) {
        for (Map.Entry<MacAddress, PcapHandle> entry : openHandles.entrySet()) {
            PcapHandle handle = entry.getValue();
            MacAddress macAddress = entry.getKey();

            Function<Object, Boolean> lldpTimer =
                message -> {
                    // Construct and send the LLDP Probe
                    TlvOrgSpecificProfibus portStatus = new TlvOrgSpecificProfibus(
                        new TlvProfibusSubTypePortStatus(0x00, false, false, (byte) 0x00)
                    );

                    TlvOrgSpecificProfibus chassisMac = new TlvOrgSpecificProfibus(
                        new TlvProfibusSubTypeChassisMac(macAddress)
                    );

                    TlvIeee8023MacPhyConfigStatus ieee = new TlvIeee8023MacPhyConfigStatus(
                        (short) 0x03,
                        0x0020,
                        0x0010
                    );

                    Ethernet_Frame identificationRequest = null;
                    try {
                        identificationRequest = new Ethernet_Frame(
                            LLDP_BROADCAST_MAC_ADDRESS,
                            macAddress,
                            new Ethernet_FramePayload_LLDP(
                                new Lldp_Pdu(
                                    Arrays.asList(
                                        new TlvChassisId(
                                            PLC4X_LLDP_IDENTIFIER.length() + 1,
                                            (short) 7,
                                            PLC4X_LLDP_IDENTIFIER
                                        ),
                                        new TlvPortId(
                                            PLC4X_LLDP_PORT.length() + 1,
                                            (short) 7,
                                            PLC4X_LLDP_PORT
                                        ),
                                        new TlvTimeToLive(2, 20),
                                        new TlvOrganizationSpecific(
                                            portStatus.getLengthInBytes(),
                                            portStatus
                                        ),
                                        new TlvOrganizationSpecific(
                                            chassisMac.getLengthInBytes(),
                                            chassisMac
                                        ),
                                        new TlvOrganizationSpecific(
                                            ieee.getLengthInBytes(),
                                            new TlvOrgSpecificIeee8023(ieee)
                                        ),
                                        new TlvManagementAddress(
                                            12,
                                            ManagementAddressSubType.IPV4,
                                            new IpAddress(Hex.decodeHex("c0a85a6e")),
                                            (short) 0x03,
                                            0x01L,
                                            (short) 0x00
                                        ),
                                        new EndOfLldp(0)
                                    )
                                )));
                    } catch (DecoderException e) {
                        throw new RuntimeException(e);
                    }
                    WriteBufferByteBased buffer = new WriteBufferByteBased(identificationRequest.getLengthInBytes());
                    try {
                        identificationRequest.serialize(buffer);
                        Packet packet = EthernetPacket.newPacket(buffer.getBytes(), 0, identificationRequest.getLengthInBytes());
                        handle.sendPacket(packet);
                    } catch (PcapNativeException | NotOpenException | SerializationException | IllegalRawDataException e) {
                        throw new RuntimeException(e);
                    }
                    return null;
                };
            Timer timer = new Timer();
            periodicTimers.add(timer);

            // Schedule to run after every 3 second(3000 millisecond)
            timer.scheduleAtFixedRate(
                new PeriodicTask(handle, lldpTimer),
                0,
                period);
        }
    }

    private static class PeriodicTask extends TimerTask {

        private final PcapHandle handle;
        private final Function<Object, Boolean> operator;

        public PeriodicTask(PcapHandle handle, Function<Object, Boolean> operator) {
            this.handle = handle;
            this.operator = operator;
        }

        @Override
        public void run() {
            operator.apply(null);
        }

    }

}