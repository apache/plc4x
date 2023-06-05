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

package org.apache.plc4x.java.profinet.protocol;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.discovery.ProfinetDiscoverer;
import org.apache.plc4x.java.profinet.gsdml.*;
import org.apache.plc4x.java.profinet.packets.PnDcpPacketFactory;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

public class ProfinetProtocolLogic extends Plc4xProtocolBase<Ethernet_Frame> implements HasConfiguration<ProfinetConfiguration> {

    private ProfinetDriverContext driverContext;
    private ProfinetConfiguration configuration;

    private final Logger logger = LoggerFactory.getLogger(ProfinetProtocolLogic.class);

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);
        if(!(driverContext instanceof ProfinetDriverContext)) {
            throw new PlcRuntimeException(
                "Expecting a driverContext of type ProfinetDriverContext, but got " + driverContext.getClass().getName());
        }
        this.driverContext = (ProfinetDriverContext) driverContext;
    }

    @Override
    public void setConfiguration(ProfinetConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public void onConnect(ConversationContext<Ethernet_Frame> context) {
        super.onConnect(context);

        RawSocketChannel rawSocketChannel = (RawSocketChannel) context.getChannel();

        // The RawSocketChannel is configured to automatically resolve the mac address of the target.
        MacAddress remoteMacAddress = new MacAddress(rawSocketChannel.getRemoteMacAddress().getAddress());
        MacAddress localMacAddress = new MacAddress(rawSocketChannel.getLocalMacAddress().getAddress());

        // Construct and send the search request.
        context.sendRequest(PnDcpPacketFactory.createIdentificationRequest(localMacAddress, remoteMacAddress))
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
            .onTimeout(e -> {
                logger.error("Timeout during Connection establishing, closing channel...");
                context.getChannel().close();
            })
            .onError((ethernetFrame, throwable) -> {
                logger.error("Exception during Connection establishing, closing channel...", throwable);
                context.getChannel().close();
            })
            .unwrap(ethernetFrame -> {
                if(ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                    return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                }
                return ethernetFrame.getPayload();
            })
            .check(ethernetFramePayload -> ethernetFramePayload instanceof Ethernet_FramePayload_PnDcp)
            .unwrap(ethernetFramePayload -> (Ethernet_FramePayload_PnDcp) ethernetFramePayload)
            .unwrap(Ethernet_FramePayload_PnDcp::getPdu)
            .check(pnDcpPdu -> pnDcpPdu instanceof PnDcp_Pdu_IdentifyRes)
            .unwrap(pnDcpPdu -> (PnDcp_Pdu_IdentifyRes) pnDcpPdu)
            .handle(pnDcpPduIdentifyRes -> {
                // Index the blocks of the response
                Map<String, PnDcp_Block> blocks = new HashMap<>();
                for (PnDcp_Block block : pnDcpPduIdentifyRes.getBlocks()) {
                    String blockName = block.getOption().name() + "-" + block.getSuboption().toString();
                    blocks.put(blockName, block);
                }

                if (blocks.containsKey(ProfinetDiscoverer.DEVICE_TYPE_NAME)) {
                    PnDcp_Block_DevicePropertiesDeviceVendor block = (PnDcp_Block_DevicePropertiesDeviceVendor) blocks.get(ProfinetDiscoverer.DEVICE_TYPE_NAME);
                    driverContext.setDeviceType(new String(block.getDeviceVendorValue()));
                }

                if (blocks.containsKey(ProfinetDiscoverer.DEVICE_NAME_OF_STATION)) {
                    PnDcp_Block_DevicePropertiesNameOfStation block = (PnDcp_Block_DevicePropertiesNameOfStation) blocks.get(ProfinetDiscoverer.DEVICE_NAME_OF_STATION);
                    driverContext.setDeviceName(new String(block.getNameOfStation()));
                }

                if (blocks.containsKey(ProfinetDiscoverer.DEVICE_ROLE)) {
                    PnDcp_Block_DevicePropertiesDeviceRole block = (PnDcp_Block_DevicePropertiesDeviceRole) blocks.get(ProfinetDiscoverer.DEVICE_ROLE);
                    List<String> roles = new ArrayList<>();
                    if (block.getPnioSupervisor()) {
                        roles.add("SUPERVISOR");
                    }
                    if (block.getPnioMultidevive()) {
                        roles.add("MULTIDEVICE");
                    }
                    if (block.getPnioController()) {
                        roles.add("CONTROLLER");
                    }
                    if (block.getPnioDevice()) {
                        roles.add("DEVICE");
                    }
                    driverContext.setRoles(roles);
                }

                if (blocks.containsKey(ProfinetDiscoverer.DEVICE_ID)) {
                    PnDcp_Block_DevicePropertiesDeviceId block = (PnDcp_Block_DevicePropertiesDeviceId) blocks.get(ProfinetDiscoverer.DEVICE_ID);
                    driverContext.setVendorId(block.getVendorId());
                    driverContext.setDeviceId(block.getDeviceId());
                }

                // Check if we actually got the vendor-id and product-id ...
                // without these, we don't know what to do with the device.
                if ((driverContext.getVendorId() == 0) || (driverContext.getDeviceId() == 0)) {
                    logger.error("Unable to determine vendor-id or product-id, closing channel...");
                    context.getChannel().close();
                    return;
                }

                // Look up the GSD file for this device ...
                ProfinetISO15745Profile deviceProfile = configuration.getGsdProfile(driverContext.getVendorId(), driverContext.getDeviceId());
                if (deviceProfile == null) {
                    logger.error("Unable to find GSD profile for device with vendor-id {} and device-id {}", driverContext.getVendorId(), driverContext.getDeviceId());
                    context.getChannel().close();
                    return;
                }
                driverContext.setDeviceProfile(deviceProfile);

                // If the user provided a DAP id in the connection string, use that (after checking that it exists)
                if (configuration.dapId != null) {
                    for (ProfinetDeviceAccessPointItem profinetDeviceAccessPointItem : deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList()) {
                        if(profinetDeviceAccessPointItem.getId().equalsIgnoreCase(configuration.dapId)) {
                            driverContext.setDapId(profinetDeviceAccessPointItem.getId());
                            break;
                        }
                    }
                    if(driverContext.getDapId() == null) {
                        logger.error("Couldn't find requested device access points (DAP): {}", configuration.dapId);
                        context.getChannel().close();
                        return;
                    }
                }

                // In theory, we would need to resolve the details for the endpoint for implicit reading in order
                // to correctly read and write asynchronously:
                // https://cache.industry.siemens.com/dl/files/980/109810980/att_1107623/v4/109810980_ImplcitDataRecordHandling_DOC_de_V10.pdf
                // Here and in the recordings of the PRONETA communication, we would
                // - first request a "handle",
                // - then in a second request using the handle request the port and object id
                // - do a third request
                // However it seems that we can simply send to the default udp port and we'll get a response from the
                // correct port, and we can calculate the object id based on the vendor id and device id, which we
                // already have from the discovery.

                // If there's more than one DAP, read the I&M0 block to get the order number,
                // which allows us to find out which DAP to use.
                else if(deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().size() > 1) {
                    RawSocketChannel pnChannel = ((RawSocketChannel) context.getChannel());
                    // Try to read the I&M0 block
                    context.sendRequest(PnDcpPacketFactory.createReadIandM0BlockRequest(pnChannel, driverContext))
                        .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
                        .onTimeout(e -> {
                            logger.error("Timeout during Connection establishing, closing channel...");
                            context.getChannel().close();
                        })
                        .onError((ethernetFrame, throwable) -> {
                            logger.error("Exception during Connection establishing, closing channel...", throwable);
                            context.getChannel().close();
                        })
                        .unwrap(ethernetFrame -> {
                            if(ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                                return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                            }
                            return ethernetFrame.getPayload();
                        })
                        .check(ethernetFramePayload -> ethernetFramePayload instanceof Ethernet_FramePayload_IPv4)
                        .unwrap(ethernetFramePayload -> (Ethernet_FramePayload_IPv4) ethernetFramePayload)
                        .unwrap(Ethernet_FramePayload_IPv4::getPayload)
                        .check((dceRpcPacket -> dceRpcPacket.getPayload() instanceof PnIoCm_Packet_Res))
                        .unwrap(dceRpcPacket -> (PnIoCm_Packet_Res) dceRpcPacket.getPayload())
                        .handle(dceRpcPacketRes -> {
                            if(dceRpcPacketRes.getBlocks().size() != 2) {
                                logger.error("Expected 2 blocks in the response, closing channel...");
                                context.getChannel().close();
                                return;
                            }
                            if(!(dceRpcPacketRes.getBlocks().get(0) instanceof IODReadResponseHeader)) {
                                logger.error("The first block was expected to be of type IODReadResponseHeader, closing channel...");
                                context.getChannel().close();
                                return;
                            }
                            if(!(dceRpcPacketRes.getBlocks().get(1) instanceof PnIoCm_Block_IAndM0)) {
                                logger.error("The second block was expected to be of type PnIoCm_Block_IAndM0, closing channel...");
                                context.getChannel().close();
                                return;
                            }
                            PnIoCm_Block_IAndM0 iAndM0 = (PnIoCm_Block_IAndM0) dceRpcPacketRes.getBlocks().get(1);

                            // Now that we've got the I&M0 block, we can access the order number and use that to match
                            // the right DAP
                            String orderId = iAndM0.getOrderId().trim();

                            // Iterate through all available options and try to find one that matches the
                            // current order id.
                            for (ProfinetDeviceAccessPointItem curDap : deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList()) {
                                if((curDap.getModuleInfo() != null) && (curDap.getModuleInfo().getOrderNumber() != null)) {
                                    String curOrderNumber = curDap.getModuleInfo().getOrderNumber().getValue();
                                    // The order number seems to contain wildcards "*" ... we need to
                                    curOrderNumber = curOrderNumber.replace("*", ".");
                                    if(orderId.matches(curOrderNumber)) {
                                        driverContext.setDapId(curDap.getId());
                                        break;
                                    }
                                }
                            }
                            // Abort, if we weren't able to detect a DAP.
                            if(driverContext.getDapId() == null) {
                                logger.error("Unable to auto-detect the device access point, closing channel...");
                                context.getChannel().close();
                            }
                        });
                }
                // If the current device only has one DAP (like most devices), simply use that.
                else if (deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().size() == 1) {
                    driverContext.setDapId(deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().get(0).getId());
                } else {
                    logger.error("GSD descriptor doesn't contain any device access points");
                    context.getChannel().close();
                    return;
                }

                context.fireConnected();
            });
    }

    @Override
    public void close(ConversationContext<Ethernet_Frame> context) {

    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        if (driverContext.getDeviceProfile() == null) {
            return CompletableFuture.failedFuture(new PlcConnectionException("Unable to find GSD file for given device"));
        }

        ProfinetISO15745Profile deviceProfile = driverContext.getDeviceProfile();
        for (ProfinetModuleItem module : deviceProfile.getProfileBody().getApplicationProcess().getModuleList()) {
            for (ProfinetVirtualSubmoduleItem submodule : module.getVirtualSubmoduleList()) {

            }
        }
        return super.browse(browseRequest);
    }

}
