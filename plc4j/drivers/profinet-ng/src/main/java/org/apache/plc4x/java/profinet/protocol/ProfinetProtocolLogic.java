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
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.discovery.ProfinetDiscoverer;
import org.apache.plc4x.java.profinet.gsdml.*;
import org.apache.plc4x.java.profinet.packets.PnDcpPacketFactory;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.profinet.utils.ProfinetDataTypeMapper;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcBrowseItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
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
        CompletableFuture<PnDcp_Pdu_IdentifyRes> future =
            PnDcpPacketFactory.sendIdentificationRequest(context, localMacAddress, remoteMacAddress);
        future.whenComplete((identifyRes, throwable) -> {
            if(throwable != null) {
                logger.error("Unable to determine vendor-id or product-id, closing channel...", throwable);
                context.getChannel().close();
                return;
            }

            // Extract required information from the response.
            extractBlockInfo(identifyRes.getBlocks());

            // Check if we actually got the vendor-id and product-id ...
            // without these, we don't know what to do with the device.
            if ((driverContext.getVendorId() == 0) || (driverContext.getDeviceId() == 0)) {
                logger.error("Unable to determine vendor-id or product-id, closing channel...");
                context.getChannel().close();
                return;
            }

            // Look up the GSD file for this device ...
            ProfinetISO15745Profile deviceProfile =
                configuration.getGsdProfile(driverContext.getVendorId(), driverContext.getDeviceId());
            if (deviceProfile == null) {
                logger.error("Unable to find GSD profile for device with vendor-id {} and device-id {}",
                    driverContext.getVendorId(), driverContext.getDeviceId());
                context.getChannel().close();
                return;
            }

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

            // If we've found a device profile with at least one dap, we request the "real identification data" from the
            // device. The response contains information about which slots are present and which module identifiers
            // apply to them as well as which subslots are present and which submodule identifiers these have.
            // In this part of the code, we simply look up the gsd dap, modules and submodules that match these
            // identifiers and save them in an easily accessible format in the deviceContext.
            else if(deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().size() > 0) {
                // Build an index of the String names.
                Map<String, String> textMapping = new HashMap<>();
                for (ProfinetTextIdValue profinetTextIdValue : deviceProfile.getProfileBody().getApplicationProcess().getExternalTextList().getPrimaryLanguage().getText()) {
                    textMapping.put(profinetTextIdValue.getTextId(), profinetTextIdValue.getValue());
                }

                // Try to read the RealIdentificationData ...
                RawSocketChannel pnChannel = ((RawSocketChannel) context.getChannel());
                CompletableFuture<PnIoCm_Block_RealIdentificationData> future1 =
                    PnDcpPacketFactory.sendRealIdentificationDataRequest(context, pnChannel, driverContext);
                future1.whenComplete((realIdentificationData, throwable1) -> {
                    if(throwable1 != null) {
                        logger.error("Unable to detect device access point, closing channel...", throwable1);
                        context.getChannel().close();
                        return;
                    }

                    // Build an index of all identification numbers for each slot and subslot the device supports.
                    Map<Integer, Long> slotModuleIdentificationNumbers = new HashMap<>();
                    Map<Integer, Map<Integer, Long>> subslotModuleIdentificationNumbers = new HashMap<>();
                    for (PnIoCm_RealIdentificationApi api : realIdentificationData.getApis()) {
                        for (PnIoCm_RealIdentificationApi_Slot curSlot : api.getSlots()) {
                            slotModuleIdentificationNumbers.put(curSlot.getSlotNumber(), curSlot.getModuleIdentNumber());
                            if(!subslotModuleIdentificationNumbers.containsKey(curSlot.getSlotNumber())) {
                                subslotModuleIdentificationNumbers.put(curSlot.getSlotNumber(), new HashMap<>());
                            }
                            for (PnIoCm_RealIdentificationApi_Subslot curSubslot : curSlot.getSubslots()) {
                                subslotModuleIdentificationNumbers.get(curSlot.getSlotNumber()).put(curSubslot.getSubslotNumber(), curSubslot.getSubmoduleIdentNumber());
                            }
                        }
                    }

                    // Get the module identification number of slot 0 (Which is always the DAP)
                    long dapModuleIdentificationNumber = slotModuleIdentificationNumbers.get(0);
                    if(dapModuleIdentificationNumber == 0){
                        logger.error("Unable to detect device access point, closing channel...");
                        context.getChannel().close();
                        return;
                    }

                    // Iterate through all available DAPs and find the one with the matching module ident number.
                    for (ProfinetDeviceAccessPointItem curDap : deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList()) {
                        String moduleIdentNumberStr = curDap.getModuleIdentNumber();
                        if(moduleIdentNumberStr.startsWith("0x") || moduleIdentNumberStr.startsWith("0X")) {
                            moduleIdentNumberStr = moduleIdentNumberStr.substring(2);
                        }
                        long moduleIdentNumber = Long.parseLong(moduleIdentNumberStr, 16);
                        if(moduleIdentNumber == dapModuleIdentificationNumber) {
                            driverContext.setDap(curDap);
                            break;
                        }
                    }
                    // Abort, if we weren't able to detect a DAP.
                    if(driverContext.getDap() == null) {
                        logger.error("Unable to auto-detect the device access point, closing channel...");
                        context.getChannel().close();
                        return;
                    }

                    // Iterate through all available modules and find the ones we're using and build an index of them.
                    Map<Integer, ProfinetModuleItem> moduleIndex = new HashMap<>();
                    Map<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> submoduleIndex = new HashMap<>();
                    for (Map.Entry<Integer, Long> moduleEntry : slotModuleIdentificationNumbers.entrySet()) {
                        int curSlot = moduleEntry.getKey();
                        // Slot 0 is the DAP, so we'll continue with the next one.
                        if(curSlot == 0) {
                            continue;
                        }
                        long curModuleIdentifier = moduleEntry.getValue();
                        // Find the module that has the given module ident number.
                        for (ProfinetModuleItem curModule : deviceProfile.getProfileBody().getApplicationProcess().getModuleList()) {
                            String moduleIdentNumberStr = curModule.getModuleIdentNumber();
                            if(moduleIdentNumberStr.startsWith("0x") || moduleIdentNumberStr.startsWith("0X")) {
                                moduleIdentNumberStr = moduleIdentNumberStr.substring(2);
                            }
                            long moduleIdentNumber = Long.parseLong(moduleIdentNumberStr, 16);
                            if(curModuleIdentifier == moduleIdentNumber) {
                                moduleIndex.put(curSlot, curModule);

                                // Now get all submodules of this module.
                                Map<Integer, Long> curSubmoduleIndex = subslotModuleIdentificationNumbers.get(curSlot);
                                for (Map.Entry<Integer, Long> submoduleEntry : curSubmoduleIndex.entrySet()) {
                                    int curSubslot = submoduleEntry.getKey();
                                    long curSubmoduleIdentNumber = submoduleEntry.getValue();
                                    for (ProfinetVirtualSubmoduleItem curSubmodule : curModule.getVirtualSubmoduleList()) {
                                        String submoduleIdentNumberStr = curSubmodule.getSubmoduleIdentNumber();
                                        if(submoduleIdentNumberStr.startsWith("0x") || submoduleIdentNumberStr.startsWith("0X")) {
                                            submoduleIdentNumberStr = submoduleIdentNumberStr.substring(2);
                                        }
                                        long submoduleIdentNumber = Long.parseLong(submoduleIdentNumberStr, 16);
                                        if(curSubmoduleIdentNumber == submoduleIdentNumber) {
                                            if(!submoduleIndex.containsKey(curSlot)) {
                                                submoduleIndex.put(curSlot, new HashMap<>());
                                            }
                                            submoduleIndex.get(curSlot).put(curSubslot, curSubmodule);

                                            // Replace the text-ids with readable values
                                            for (ProfinetIoDataInput profinetIoDataInput : curSubmodule.getIoData().getInput()) {
                                                for (ProfinetDataItem profinetDataItem : profinetIoDataInput.getDataItemList()) {
                                                    if(textMapping.containsKey(profinetDataItem.getTextId())) {
                                                        profinetDataItem.setTextId(textMapping.get(profinetDataItem.getTextId()));
                                                    }
                                                }
                                            }
                                            for (ProfinetIoDataOutput profinetIoDataOutput : curSubmodule.getIoData().getOutput()) {
                                                for (ProfinetDataItem profinetDataItem : profinetIoDataOutput.getDataItemList()) {
                                                    if(textMapping.containsKey(profinetDataItem.getTextId())) {
                                                        profinetDataItem.setTextId(textMapping.get(profinetDataItem.getTextId()));
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                    driverContext.setModuleIndex(moduleIndex);
                    driverContext.setSubmoduleIndex(submoduleIndex);

                    context.fireConnected();
                });
                // Try to read the I&M0 block
                // * Commented out this code, as here it performed the same task as the ReadRealIdentificationData approach and If the other has side-effects, I can roll-back to this.
                /*CompletableFuture<PnIoCm_Block_IAndM0> pnIoCmBlockIAndM0CompletableFuture = PnDcpPacketFactory.sendReadIAndM0BlockRequest(context, pnChannel, driverContext);
                pnIoCmBlockIAndM0CompletableFuture.whenComplete((pnIoCmBlockIAndM0, throwable1) -> {
                    if(throwable1 != null) {
                        logger.error("Unable to detect device access point, closing channel...", throwable1);
                        context.getChannel().close();
                        return;
                    }

                    // Now that we've got the I&M0 block, we can access the order number and use that to match
                    // the right DAP
                    String orderId = pnIoCmBlockIAndM0.getOrderId().trim();

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
                        return;
                    }

                    context.fireConnected();
                });*/
            }

            else {
                logger.error("GSD descriptor doesn't contain any device access points");
                context.getChannel().close();
            }
        });
    }

    @Override
    public void close(ConversationContext<Ethernet_Frame> context) {

    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        Map<String, List<PlcBrowseItem>> values = new HashMap<>();
        for (String queryName : browseRequest.getQueryNames()) {
            List<PlcBrowseItem> items = new ArrayList<>();
            for(Map.Entry<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> slotEntry : driverContext.getSubmoduleIndex().entrySet()) {
                int slot = slotEntry.getKey();
                for(Map.Entry<Integer, ProfinetVirtualSubmoduleItem> subslotEntry: slotEntry.getValue().entrySet()) {
                    int subslot = subslotEntry.getKey();
                    ProfinetVirtualSubmoduleItem subslotModule = subslotEntry.getValue();

                    // Add all the input tags.
                    for (ProfinetIoDataInput profinetIoDataInput : subslotModule.getIoData().getInput()) {
                        for (int i = 0; i < profinetIoDataInput.getDataItemList().size(); i++) {
                            ProfinetDataItem profinetDataItem = profinetIoDataInput.getDataItemList().get(i);
                            ProfinetDataTypeMapper.DataTypeInformation dataTypeInformation =
                                ProfinetDataTypeMapper.getPlcValueType(profinetDataItem);
                            // The ids have been replaced by real textual values in the connection phase.
                            String name = profinetDataItem.getTextId();
                            items.add(new DefaultPlcBrowseItem(new ProfinetTag(
                                slot, subslot, ProfinetTag.Direction.INPUT,
                                i, dataTypeInformation.getPlcValueType(), dataTypeInformation.getNumElements()),
                                name, false, true, true,
                                Collections.emptyMap(), Collections.emptyMap()));
                        }
                    }

                    // Add all the output tags.
                    for (ProfinetIoDataOutput profinetIoDataOutput : subslotModule.getIoData().getOutput()) {
                        for (int i = 0; i < profinetIoDataOutput.getDataItemList().size(); i++) {
                            ProfinetDataItem profinetDataItem = profinetIoDataOutput.getDataItemList().get(i);
                            ProfinetDataTypeMapper.DataTypeInformation dataTypeInformation =
                                ProfinetDataTypeMapper.getPlcValueType(profinetDataItem);
                            // The ids have been replaced by real textual values in the connection phase.
                            String name = profinetDataItem.getTextId();
                            items.add(new DefaultPlcBrowseItem(new ProfinetTag(
                                slot, subslot, ProfinetTag.Direction.OUTPUT,
                                i, dataTypeInformation.getPlcValueType(), dataTypeInformation.getNumElements()),
                                name, false, true, true,
                                Collections.emptyMap(), Collections.emptyMap()));
                        }
                    }
                }
            }
            responseCodes.put(queryName, PlcResponseCode.OK);
            values.put(queryName, items);
        }
        PlcBrowseResponse response = new DefaultPlcBrowseResponse(browseRequest, responseCodes, values);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        // When subscribing, we actually set up the PN IO Application Relation and make the remote device start sending data.
        if (driverContext.getDap() == null) {
            return CompletableFuture.failedFuture(new PlcConnectionException("DAP not set"));
        }

        // Find the matching data in the device descriptor.
        for (String tagName : subscriptionRequest.getTagNames()) {
            PlcSubscriptionTag tag = subscriptionRequest.getTag(tagName);
            if(!(tag instanceof ProfinetTag)) {
                // TODO: Add an error code for this field.
                continue;
            }
            ProfinetTag profinetTag = (ProfinetTag) tag;
            int slot = profinetTag.getSlot();
            int subSlot = profinetTag.getSubSlot();

        }
        return null;
    }

    protected void extractBlockInfo(List<PnDcp_Block> blocks) {
        // Index the blocks of the response
        Map<String, PnDcp_Block> blockMap = new HashMap<>();
        for (PnDcp_Block block : blocks) {
            String blockName = block.getOption().name() + "-" + block.getSuboption().toString();
            blockMap.put(blockName, block);
        }

        if (blockMap.containsKey(ProfinetDiscoverer.DEVICE_TYPE_NAME)) {
            PnDcp_Block_DevicePropertiesDeviceVendor block =
                (PnDcp_Block_DevicePropertiesDeviceVendor) blockMap.get(ProfinetDiscoverer.DEVICE_TYPE_NAME);
            driverContext.setDeviceType(new String(block.getDeviceVendorValue()));
        }

        if (blockMap.containsKey(ProfinetDiscoverer.DEVICE_NAME_OF_STATION)) {
            PnDcp_Block_DevicePropertiesNameOfStation block =
                (PnDcp_Block_DevicePropertiesNameOfStation) blockMap.get(ProfinetDiscoverer.DEVICE_NAME_OF_STATION);
            driverContext.setDeviceName(new String(block.getNameOfStation()));
        }

        if (blockMap.containsKey(ProfinetDiscoverer.DEVICE_ROLE)) {
            PnDcp_Block_DevicePropertiesDeviceRole block =
                (PnDcp_Block_DevicePropertiesDeviceRole) blockMap.get(ProfinetDiscoverer.DEVICE_ROLE);
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

        if (blockMap.containsKey(ProfinetDiscoverer.DEVICE_ID)) {
            PnDcp_Block_DevicePropertiesDeviceId block =
                (PnDcp_Block_DevicePropertiesDeviceId) blockMap.get(ProfinetDiscoverer.DEVICE_ID);
            driverContext.setVendorId(block.getVendorId());
            driverContext.setDeviceId(block.getDeviceId());
        }
    }

}
