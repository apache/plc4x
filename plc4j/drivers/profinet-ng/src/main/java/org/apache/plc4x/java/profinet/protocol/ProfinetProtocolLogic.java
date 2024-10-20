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
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.discovery.ProfinetDiscoverer;
import org.apache.plc4x.java.profinet.gsdml.*;
import org.apache.plc4x.java.profinet.packets.PnDcpPacketFactory;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.profinet.tag.ProfinetTagHandler;
import org.apache.plc4x.java.profinet.utils.ProfinetDataTypeMapper;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcBrowseItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetSocketAddress;
import java.security.SecureRandom;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;

public class ProfinetProtocolLogic extends Plc4xProtocolBase<Ethernet_Frame> implements HasConfiguration<ProfinetConfiguration> {

    // This is the minimum cycle time defined in the PN Spec (Page 205): 31,25us
    private static final int NANOS_PER_CYCLE = 31250;

    private ProfinetDriverContext profinetDriverContext;
    private ProfinetConfiguration configuration;

    private Integer cycleOffset = null;
    private boolean connected = false;

    private final Logger logger = LoggerFactory.getLogger(ProfinetProtocolLogic.class);

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);
        if (!(driverContext instanceof ProfinetDriverContext)) {
            throw new PlcRuntimeException(
                "Expecting a driverContext of type ProfinetDriverContext, but got " + driverContext.getClass().getName());
        }
        this.profinetDriverContext = (ProfinetDriverContext) driverContext;
    }

    @Override
    public void setConfiguration(ProfinetConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new ProfinetTagHandler();
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
            if (throwable != null) {
                logger.error("Unable to determine vendor-id or product-id, closing channel...", throwable);
                context.getChannel().close();
                return;
            }

            // Extract required information from the response.
            extractBlockInfo(identifyRes.getBlocks());

            // Check if we actually got the vendor-id and product-id ...
            // without these, we don't know what to do with the device.
            if ((profinetDriverContext.getVendorId() == 0) || (profinetDriverContext.getDeviceId() == 0)) {
                logger.error("Unable to determine vendor-id or product-id, closing channel...");
                context.getChannel().close();
                return;
            }

            // Look up the GSD file for this device ...
            ProfinetISO15745Profile deviceProfile =
                configuration.getGsdProfile(profinetDriverContext.getVendorId(), profinetDriverContext.getDeviceId());
            if (deviceProfile == null) {
                logger.error("Unable to find GSD profile for device with vendor-id {} and device-id {}",
                    profinetDriverContext.getVendorId(), profinetDriverContext.getDeviceId());
                context.getChannel().close();
                return;
            }

            // If the user provided a DAP id in the connection string, use that (after checking that it exists)
            if (configuration.dapId != null) {
                for (ProfinetDeviceAccessPointItem profinetDeviceAccessPointItem : deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList()) {
                    if (profinetDeviceAccessPointItem.getId().equalsIgnoreCase(configuration.dapId)) {
                        profinetDriverContext.setDap(profinetDeviceAccessPointItem);
                        break;
                    }
                }
                if (profinetDriverContext.getDap() == null) {
                    logger.error("Couldn't find requested device access points (DAP): {}", configuration.dapId);
                    context.getChannel().close();
                }
            }
            // If the user didn't define a dap, but the profile only has one, use that.
            else if (deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().size() == 1) {
                profinetDriverContext.setDap(deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().get(0));
            }
            // Otherwise we'll have to use the RealIdentificationDataRequest to fetch all of this information.
            // However, this is not supported by all PN devices.

            // If we've found a device profile with at least one dap, we request the "real identification data" from the
            // device. The response contains information about which slots are present and which module identifiers
            // apply to them as well as which subslots are present and which submodule identifiers these have.
            // In this part of the code, we simply look up the gsd dap, modules and submodules that match these
            // identifiers and save them in an easily accessible format in the deviceContext.
            if (!deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().isEmpty()) {
                // Build an index of the String names.
                Map<String, String> textMapping = new HashMap<>();
                for (ProfinetTextIdValue profinetTextIdValue : deviceProfile.getProfileBody().getApplicationProcess().getExternalTextList().getPrimaryLanguage().getText()) {
                    textMapping.put(profinetTextIdValue.getTextId(), profinetTextIdValue.getValue());
                }

                // Try to read the RealIdentificationData ...
                RawSocketChannel pnChannel = ((RawSocketChannel) context.getChannel());
                CompletableFuture<PnIoCm_Block_RealIdentificationData> future1 =
                    PnDcpPacketFactory.sendRealIdentificationDataRequest(context, pnChannel, profinetDriverContext);
                future1.whenComplete((realIdentificationData, throwable1) -> {
                    // If the device didn't support this, we'll have to handle the
                    // module index and submodule indexes when trying to use them.
                    if (throwable1 != null) {
                        if (profinetDriverContext.getDap() != null) {
                            context.fireConnected();
                        } else {
                            logger.error("Unable to auto-configure connection, please be sure to provide the 'dap-id' connection parameter");
                            context.getChannel().close();
                            return;
                        }
                        return;
                    }

                    // Build an index of all identification numbers for each slot and subslot the device supports.
                    Map<Integer, Long> slotModuleIdentificationNumbers = new HashMap<>();
                    Map<Integer, Map<Integer, Long>> subslotModuleIdentificationNumbers = new HashMap<>();
                    for (PnIoCm_RealIdentificationApi api : realIdentificationData.getApis()) {
                        for (PnIoCm_RealIdentificationApi_Slot curSlot : api.getSlots()) {
                            slotModuleIdentificationNumbers.put(curSlot.getSlotNumber(), curSlot.getModuleIdentNumber());
                            if (!subslotModuleIdentificationNumbers.containsKey(curSlot.getSlotNumber())) {
                                subslotModuleIdentificationNumbers.put(curSlot.getSlotNumber(), new HashMap<>());
                            }
                            for (PnIoCm_RealIdentificationApi_Subslot curSubslot : curSlot.getSubslots()) {
                                subslotModuleIdentificationNumbers.get(curSlot.getSlotNumber()).put(curSubslot.getSubslotNumber(), curSubslot.getSubmoduleIdentNumber());
                            }
                        }
                    }

                    // Get the module identification number of slot 0 (Which is always the DAP)
                    long dapModuleIdentificationNumber = slotModuleIdentificationNumbers.get(0);
                    if (dapModuleIdentificationNumber == 0) {
                        logger.error("Unable to detect device access point, closing channel...");
                        context.getChannel().close();
                        return;
                    }

                    // Iterate through all available DAPs and find the one with the matching module ident number.
                    for (ProfinetDeviceAccessPointItem curDap : deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList()) {
                        String moduleIdentNumberStr = curDap.getModuleIdentNumber();
                        if (moduleIdentNumberStr.startsWith("0x") || moduleIdentNumberStr.startsWith("0X")) {
                            moduleIdentNumberStr = moduleIdentNumberStr.substring(2);
                        }
                        long moduleIdentNumber = Long.parseLong(moduleIdentNumberStr, 16);
                        if (moduleIdentNumber == dapModuleIdentificationNumber) {
                            if ((profinetDriverContext.getDap() != null) &&
                                !profinetDriverContext.getDap().getId().equals(curDap.getId())) {
                                logger.warn("DAP configured in connection string differs from device-configuration.");
                            }
                            profinetDriverContext.setDap(curDap);
                            break;
                        }
                    }
                    // Abort, if we weren't able to detect a DAP.
                    if (profinetDriverContext.getDap() == null) {
                        logger.error("Unable to auto-detect the device access point, please provide \"dap-id\" " +
                            "option in the connection-string. Closing channel...");
                        context.getChannel().close();
                        return;
                    }

                    // Iterate through all available modules and find the ones we're using and build an index of them.
                    Map<Integer, ProfinetModuleItem> moduleIndex = new HashMap<>();
                    Map<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> submoduleIndex = new HashMap<>();
                    for (Map.Entry<Integer, Long> moduleEntry : slotModuleIdentificationNumbers.entrySet()) {
                        int curSlot = moduleEntry.getKey();
                        // Slot 0 is the DAP, which we've already handled, so we'll continue with the next one.
                        if (curSlot == 0) {
                            continue;
                        }
                        long curModuleIdentifier = moduleEntry.getValue();
                        // Find the module that has the given module ident number.
                        for (ProfinetModuleItem curModule : deviceProfile.getProfileBody().getApplicationProcess().getModuleList()) {
                            String moduleIdentNumberStr = curModule.getModuleIdentNumber();
                            if (moduleIdentNumberStr.startsWith("0x") || moduleIdentNumberStr.startsWith("0X")) {
                                moduleIdentNumberStr = moduleIdentNumberStr.substring(2);
                            }
                            long moduleIdentNumber = Long.parseLong(moduleIdentNumberStr, 16);
                            if (curModuleIdentifier == moduleIdentNumber) {
                                moduleIndex.put(curSlot, curModule);

                                // Now get all submodules of this module.
                                Map<Integer, Long> curSubmoduleIndex = subslotModuleIdentificationNumbers.get(curSlot);
                                for (Map.Entry<Integer, Long> submoduleEntry : curSubmoduleIndex.entrySet()) {
                                    int curSubslot = submoduleEntry.getKey();
                                    long curSubmoduleIdentNumber = submoduleEntry.getValue();
                                    for (ProfinetVirtualSubmoduleItem curSubmodule : curModule.getVirtualSubmoduleList()) {
                                        String submoduleIdentNumberStr = curSubmodule.getSubmoduleIdentNumber();
                                        if (submoduleIdentNumberStr.startsWith("0x") || submoduleIdentNumberStr.startsWith("0X")) {
                                            submoduleIdentNumberStr = submoduleIdentNumberStr.substring(2);
                                        }
                                        long submoduleIdentNumber = Long.parseLong(submoduleIdentNumberStr, 16);
                                        if (curSubmoduleIdentNumber == submoduleIdentNumber) {
                                            if (!submoduleIndex.containsKey(curSlot)) {
                                                submoduleIndex.put(curSlot, new HashMap<>());
                                            }
                                            submoduleIndex.get(curSlot).put(curSubslot, curSubmodule);

                                            // Replace the text-ids with readable values
                                            if (curSubmodule.getIoData().getInput() != null) {
                                                for (ProfinetIoDataInput profinetIoDataInput : curSubmodule.getIoData().getInput()) {
                                                    for (ProfinetDataItem profinetDataItem : profinetIoDataInput.getDataItemList()) {
                                                        if (textMapping.containsKey(profinetDataItem.getTextId())) {
                                                            profinetDataItem.setTextId(textMapping.get(profinetDataItem.getTextId()));
                                                        }
                                                    }
                                                }
                                            }
                                            if (curSubmodule.getIoData().getOutput() != null) {
                                                for (ProfinetIoDataOutput profinetIoDataOutput : curSubmodule.getIoData().getOutput()) {
                                                    for (ProfinetDataItem profinetDataItem : profinetIoDataOutput.getDataItemList()) {
                                                        if (textMapping.containsKey(profinetDataItem.getTextId())) {
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
                    }
                    profinetDriverContext.setModuleIndex(moduleIndex);
                    profinetDriverContext.setSubmoduleIndex(submoduleIndex);

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
            } else {
                logger.error("GSD descriptor doesn't contain any device access points");
                context.getChannel().close();
            }
        });
    }

    @Override
    public void close(ConversationContext<Ethernet_Frame> context) {
        context.getChannel().close();
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        Map<String, List<PlcBrowseItem>> values = new HashMap<>();
        for (String queryName : browseRequest.getQueryNames()) {
            List<PlcBrowseItem> items = new ArrayList<>();
            for (Map.Entry<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> slotEntry : profinetDriverContext.getSubmoduleIndex().entrySet()) {
                int slot = slotEntry.getKey();
                for (Map.Entry<Integer, ProfinetVirtualSubmoduleItem> subslotEntry : slotEntry.getValue().entrySet()) {
                    int subslot = subslotEntry.getKey();
                    ProfinetVirtualSubmoduleItem subslotModule = subslotEntry.getValue();

                    // Add all the input tags.
                    if (subslotModule.getIoData().getInput() != null) {
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
                                    name, false, true, true, false,
                                    Collections.emptyList(), Collections.emptyMap(), Collections.emptyMap()));
                            }
                        }
                    }

                    // Add all the output tags.
                    if (subslotModule.getIoData().getOutput() != null) {
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
                                    name, false, true, true, false,
                                    Collections.emptyList(), Collections.emptyMap(), Collections.emptyMap()));
                            }
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
        if (profinetDriverContext.getDap() == null) {
            return CompletableFuture.failedFuture(new PlcConnectionException("DAP not set"));
        }

        // Go through all the tags and build a sorted list of all requested tags.
        Map<Integer, Map<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>>> slots = new TreeMap<>();
        for (String tagName : subscriptionRequest.getTagNames()) {
            PlcSubscriptionTag subscriptionTag = subscriptionRequest.getTag(tagName);
            if (!(subscriptionTag.getTag() instanceof ProfinetTag)) {
                // TODO: Add an error code for this field.
                continue;
            }
            ProfinetTag profinetTag = (ProfinetTag) subscriptionTag.getTag();
            int slot = profinetTag.getSlot();
            int subSlot = profinetTag.getSubSlot();
            ProfinetTag.Direction direction = profinetTag.getDirection();
            int index = profinetTag.getIndex();
            if (!slots.containsKey(slot)) {
                slots.put(slot, new TreeMap<>());
            }
            if (!slots.get(slot).containsKey(subSlot)) {
                slots.get(slot).put(subSlot, new TreeMap<>());
            }
            if (!slots.get(slot).get(subSlot).containsKey(direction)) {
                slots.get(slot).get(subSlot).put(direction, new TreeMap<>());
            }
            slots.get(slot).get(subSlot).get(direction).put(index, profinetTag);
        }

        // Create one PnIoCm_Block_ArReq
        // Create one PnIoCm_IoCrBlockReqApi for input (if there's at least one input)
        // Create one PnIoCm_IoCrBlockReqApi for output (if there's at least one output)
        // Create one PnIoCm_Block_ExpectedSubmoduleReq for every slot being referenced.
        // Create one PnIoCm_Block_AlarmCrReq

        // Go through the sorted slots and subslots and fill the datastructures.
        int inputFrameOffset = 0;
        int outputFrameOffset = 0;
        List<PnIoCm_IoDataObject> inputMessageDataObjects = new ArrayList<>();
        List<PnIoCm_IoCs> inputMessageCs = new ArrayList<>();
        List<PnIoCm_IoDataObject> outputMessageDataObjects = new ArrayList<>();
        List<PnIoCm_IoCs> outputMessageCs = new ArrayList<>();

        // Add the inputMessageDataObjects and outputMessageDataObjects for the DAO object;
        // TODO: Get the structure of this from the GSD file. This is currently just hard-coded.
        inputMessageDataObjects.add(new PnIoCm_IoDataObject(0x00, 0x1, inputFrameOffset));
        inputFrameOffset += 1; // TODO: Add this correctly (Possibly with data from the GSD file.
        outputMessageCs.add(new PnIoCm_IoCs(0x00, 0x1, outputFrameOffset));
        outputFrameOffset += 1; // TODO: Add this correctly (Possibly with data from the GSD file.

        inputMessageDataObjects.add(new PnIoCm_IoDataObject(0x00, 0x8000, inputFrameOffset));
        inputFrameOffset += 1; // TODO: Add this correctly (Possibly with data from the GSD file.
        outputMessageCs.add(new PnIoCm_IoCs(0x00, 0x8000, outputFrameOffset));
        outputFrameOffset += 1; // TODO: Add this correctly (Possibly with data from the GSD file.

        inputMessageDataObjects.add(new PnIoCm_IoDataObject(0x00, 0x8001, inputFrameOffset));
        inputFrameOffset += 1; // TODO: Add this correctly (Possibly with data from the GSD file.
        outputMessageCs.add(new PnIoCm_IoCs(0x00, 0x8001, outputFrameOffset));
        outputFrameOffset += 1; // TODO: Add this correctly (Possibly with data from the GSD file.

        inputMessageDataObjects.add(new PnIoCm_IoDataObject(0x00, 0x8002, inputFrameOffset));
        inputFrameOffset += 1; // TODO: Add this correctly (Possibly with data from the GSD file.
        outputMessageCs.add(new PnIoCm_IoCs(0x00, 0x8002, outputFrameOffset));
        outputFrameOffset += 1; // TODO: Add this correctly (Possibly with data from the GSD file.

        // TODO: this will probably need to be dynamic, based on data in the GSD file.
        List<PnIoCm_Block_ExpectedSubmoduleReq> expectedSubmodules = new ArrayList<>();
        expectedSubmodules.add(new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0, Collections.singletonList(
            new PnIoCm_ExpectedSubmoduleBlockReqApi((short) 0x0000, (short) 0x00000010, 0, Arrays.asList(
                new PnIoCm_Submodule_NoInputNoOutputData((short) 0x0001, (short) 0x00000001, false, false, false, false),
                new PnIoCm_Submodule_NoInputNoOutputData((short) 0x8000, (short) 0x00000002, false, false, false, false),
                new PnIoCm_Submodule_NoInputNoOutputData((short) 0x8001, (short) 0x00000003, false, false, false, false),
                new PnIoCm_Submodule_NoInputNoOutputData((short) 0x8002, (short) 0x00000003, false, false, false, false)
            ))
        )));

        for (Map.Entry<Integer, Map<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>>> slotEntry : slots.entrySet()) {
            int slotNumber = slotEntry.getKey();
            Map<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>> subslot = slotEntry.getValue();
            List<PnIoCm_Submodule> expectedSubmoduleData = new ArrayList<>();
            for (Map.Entry<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>> subslotEntry : subslot.entrySet()) {
                int subslotNumber = subslotEntry.getKey();
                Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>> direction = subslotEntry.getValue();

                int iocsLength = profinetDriverContext.getSubmoduleIndex().get(slotNumber).get(subslotNumber).getIoData().getIocsLength();
                // The default is 1
                if (iocsLength == 0) {
                    iocsLength = 1;
                }
                int iopsLength = profinetDriverContext.getSubmoduleIndex().get(slotNumber).get(subslotNumber).getIoData().getIopsLength();
                // The default is 1
                if (iopsLength == 0) {
                    iopsLength = 1;
                }

                PnIoCm_SubmoduleType submoduleType = PnIoCm_SubmoduleType.NO_INPUT_NO_OUTPUT_DATA;
                int inputDataLength = 0;
                int outputDataLength = 0;

                // Add one PnIoCm_IoDataObject for every input tag (Subscribe).
                // These define the structure of the data in incoming messages from the device.
                if (direction.containsKey(ProfinetTag.Direction.INPUT)) {
                    // Update the type of submodule io.
                    submoduleType = PnIoCm_SubmoduleType.INPUT_DATA;

                    Map<Integer, ProfinetTag> inputTags = direction.get(ProfinetTag.Direction.INPUT);
                    for (Map.Entry<Integer, ProfinetTag> inputTag : inputTags.entrySet()) {
                        ProfinetTag tag = inputTag.getValue();
                        int dataLength = (getDataTypeLengthInBytes(tag.getPlcValueType()) * tag.getNumElements());
                        inputDataLength += dataLength;

                        PnIoCm_IoDataObject input = new PnIoCm_IoDataObject(slotNumber, subslotNumber, inputFrameOffset);
                        inputMessageDataObjects.add(input);

                        // Get the iops-length from the IoData element and the binary length of the input
                        inputFrameOffset += dataLength + iocsLength;

                        PnIoCm_IoCs output = new PnIoCm_IoCs(slotNumber, subslotNumber, inputFrameOffset);
                        inputMessageCs.add(output);
                    }
                }

                // Add one PnIoCm_IoDataObject for every output tag (Publish)
                // These define the structure of the data in outgoing messages sent to the device.
                if (direction.containsKey(ProfinetTag.Direction.OUTPUT)) {
                    // Update the type of submodule io.
                    if (submoduleType == PnIoCm_SubmoduleType.NO_INPUT_NO_OUTPUT_DATA) {
                        submoduleType = PnIoCm_SubmoduleType.OUTPUT_DATA;
                    } else if (submoduleType == PnIoCm_SubmoduleType.INPUT_DATA) {
                        submoduleType = PnIoCm_SubmoduleType.INPUT_AND_OUTPUT_DATA;
                    }

                    Map<Integer, ProfinetTag> outputTags = direction.get(ProfinetTag.Direction.OUTPUT);
                    for (Map.Entry<Integer, ProfinetTag> outputTag : outputTags.entrySet()) {
                        // TODO: Here the offset is wrong (5 instead of 4)
                        PnIoCm_IoCs input = new PnIoCm_IoCs(slotNumber, subslotNumber, outputFrameOffset);
                        outputMessageCs.add(input);
                        outputFrameOffset += iocsLength;

                        ProfinetTag tag = outputTag.getValue();
                        int dataLength = (getDataTypeLengthInBytes(tag.getPlcValueType()) * tag.getNumElements());
                        outputDataLength += dataLength;
                        outputFrameOffset += dataLength;

                        // TODO: Hehe the offset is wrong (5 instead of 9)
                        PnIoCm_IoDataObject output = new PnIoCm_IoDataObject(slotNumber, subslotNumber, outputFrameOffset);
                        outputMessageDataObjects.add(output);

                        // Get the data-length + iocs-length from the IoData element
                        inputFrameOffset += iopsLength;
                        // Get the iops-length and the binary length of the output
                        outputFrameOffset += dataLength + iopsLength;
                    }
                }

                switch (submoduleType) {
                    case NO_INPUT_NO_OUTPUT_DATA: {
                        // TODO: Get the submodule ident number from the GSD file.
                        expectedSubmoduleData.add(new PnIoCm_Submodule_NoInputNoOutputData((short) subslotNumber, (short) 0x00000010, false, false, false, false));
                        break;
                    }
                    case INPUT_DATA: {
                        // TODO: Get the submodule ident number from the GSD file.
                        expectedSubmoduleData.add(new PnIoCm_Submodule_InputData((short) subslotNumber, (short) 0x00000010, false, false, false, false, inputDataLength));
                        break;
                    }
                    case OUTPUT_DATA: {
                        expectedSubmoduleData.add(new PnIoCm_Submodule_OutputData((short) subslotNumber, (short) 0x00000010, false, false, false, false, outputDataLength));
                        break;
                    }
                    case INPUT_AND_OUTPUT_DATA: {
                        expectedSubmoduleData.add(new PnIoCm_Submodule_InputAndOutputData((short) subslotNumber, (short) 0x00000010, false, false, false, false, inputDataLength, outputDataLength));
                        break;
                    }
                }
            }

            // TODO: Get the submodule ident number from the GSD file.
            expectedSubmodules.add(new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0, Collections.singletonList(new PnIoCm_ExpectedSubmoduleBlockReqApi(slotNumber, 0x00000020, 0x0000, expectedSubmoduleData))));
        }

        RawSocketChannel rawSocketChannel = (RawSocketChannel) conversationContext.getChannel();
        MacAddress remoteMacAddress = new MacAddress(rawSocketChannel.getRemoteMacAddress().getAddress());
        InetSocketAddress remoteAddress = (InetSocketAddress) rawSocketChannel.getRemoteAddress();
        MacAddress localMacAddress = new MacAddress(rawSocketChannel.getLocalMacAddress().getAddress());
        InetSocketAddress localAddress = (InetSocketAddress) rawSocketChannel.getLocalAddress();
        List<PnIoCm_Block> blocks = new ArrayList<>();
        blocks.add(new PnIoCm_Block_ArReq(
            ProfinetDriverContext.BLOCK_VERSION_HIGH, ProfinetDriverContext.BLOCK_VERSION_LOW,
            PnIoCm_ArType.IO_CONTROLLER,
            profinetDriverContext.getApplicationRelationUuid(),
            profinetDriverContext.getSessionKey(),
            localMacAddress,
            profinetDriverContext.getCmInitiatorObjectUuid(),
            false,
            profinetDriverContext.isAdvancedStartupMode(),
            false,
            false,
            PnIoCm_CompanionArType.SINGLE_AR,
            false,
            true,
            false,
            PnIoCm_State.ACTIVE,
            ProfinetDriverContext.DEFAULT_ACTIVITY_TIMEOUT,
            ProfinetDriverContext.UDP_RT_PORT,
            "plc4x"));
        blocks.add(new PnIoCm_Block_AlarmCrReq(
            ProfinetDriverContext.BLOCK_VERSION_HIGH, ProfinetDriverContext.BLOCK_VERSION_LOW,
            PnIoCm_AlarmCrType.ALARM_CR,
            ProfinetDriverContext.UDP_RT_PORT,
            false,
            false,
            1,
            3,
            0x0000,
            200,
            0xC000,
            0xA000));
        if (!inputMessageDataObjects.isEmpty() || !inputMessageCs.isEmpty()) {
            blocks.add(new PnIoCm_Block_IoCrReq(
                ProfinetDriverContext.BLOCK_VERSION_HIGH, ProfinetDriverContext.BLOCK_VERSION_LOW,
                PnIoCm_IoCrType.INPUT_CR,
                0x0001,
                ProfinetDriverContext.UDP_RT_PORT,
                false,
                false,
                false,
                false,
                PnIoCm_RtClass.RT_CLASS_2,
                ProfinetDriverContext.DEFAULT_IO_DATA_SIZE,
                // TODO: This differs: Mine is 0x0002 and Ben's is 0x8002 (Probably doesn't matter)
                0x8000 | profinetDriverContext.getAndIncrementIdentification(),
                profinetDriverContext.getSendClockFactor(),
                profinetDriverContext.getReductionRatio(),
                1,
                0,
                0xffffffffL,
                profinetDriverContext.getWatchdogFactor(),
                profinetDriverContext.getDataHoldFactor(),
                0xC000,
                ProfinetDriverContext.DEFAULT_EMPTY_MAC_ADDRESS,
                Collections.singletonList(
                    new PnIoCm_IoCrBlockReqApi(inputMessageDataObjects, inputMessageCs)
                )
            ));
        }
        if (!outputMessageDataObjects.isEmpty() || !outputMessageCs.isEmpty()) {
            blocks.add(new PnIoCm_Block_IoCrReq(
                ProfinetDriverContext.BLOCK_VERSION_HIGH, ProfinetDriverContext.BLOCK_VERSION_LOW,
                PnIoCm_IoCrType.OUTPUT_CR,
                0x0002,
                ProfinetDriverContext.UDP_RT_PORT,
                false,
                false,
                false,
                false,
                PnIoCm_RtClass.RT_CLASS_2,
                ProfinetDriverContext.DEFAULT_IO_DATA_SIZE,
                // TODO: This differs: Mine is 0x0003 and Ben's is 0x8003 (Probably doesn't matter)
                0x8000 | profinetDriverContext.getAndIncrementIdentification(),
                profinetDriverContext.getSendClockFactor(),
                profinetDriverContext.getReductionRatio(),
                1,
                0,
                0xffffffffL,
                profinetDriverContext.getWatchdogFactor(),
                profinetDriverContext.getDataHoldFactor(),
                0xC000,
                ProfinetDriverContext.DEFAULT_EMPTY_MAC_ADDRESS,
                Collections.singletonList(
                    new PnIoCm_IoCrBlockReqApi(outputMessageDataObjects, outputMessageCs)
                )
            ));
        }
        blocks.addAll(expectedSubmodules);
        PnIoCm_Packet_Req request = new PnIoCm_Packet_Req(
            16696L, 16696L, 0L, blocks);
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST,
            true, false, false,
            IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, profinetDriverContext.getDeviceId(), profinetDriverContext.getVendorId()),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            profinetDriverContext.getActivityUuid(),
            0L, 1L,
            DceRpc_Operation.CONNECT,
            (short) 0,
            request
        );
        SecureRandom rand = new SecureRandom();
        // Serialize it to a byte-payload
        Ethernet_FramePayload_IPv4 udpFrame = new Ethernet_FramePayload_IPv4(
            rand.nextInt(65536),
            true,
            false,
            (short) 64,
            new IpAddress(localAddress.getAddress().getAddress()),
            new IpAddress(remoteAddress.getAddress().getAddress()),
            profinetDriverContext.getLocalPort(),
            profinetDriverContext.getRemotePortImplicitCommunication(),
            packet
        );
        Ethernet_Frame requestEthernetFrame = new Ethernet_Frame(
            remoteMacAddress,
            localMacAddress,
            udpFrame);

        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        conversationContext.sendRequest(requestEthernetFrame)
            .name("Expect Subscription response")
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(1000))
            .onTimeout(future::completeExceptionally)
            .onError((responseEthernetFrame, throwable) -> future.completeExceptionally(throwable))
            .unwrap(Ethernet_Frame::getPayload)
            .only(Ethernet_FramePayload_IPv4.class)
            .unwrap(Ethernet_FramePayload_IPv4::getPayload)
            .handle(dceRpcPacket -> {
                if (dceRpcPacket.getPacketType() != DceRpc_PacketType.RESPONSE) {
                    future.completeExceptionally(new PlcException("Expected a response"));
                    return;
                }
                PnIoCm_Packet_Res payload = (PnIoCm_Packet_Res) dceRpcPacket.getPayload();
                // TODO: Maybe do some checks on this.

                // Now we wait for an incoming ApplicationReady request and confirm that.
                conversationContext.expectRequest(Ethernet_Frame.class, Duration.ofMillis(500000))
                    .name("Expect ApplicationReady request")
                    .onTimeout(future::completeExceptionally)
                    .check(ethernetFrame -> ethernetFrame.getPayload() instanceof Ethernet_FramePayload_IPv4)
                    .unwrap(ethernetFrame -> ((Ethernet_FramePayload_IPv4) ethernetFrame.getPayload()))
                    .check(payloadIPv4 -> payloadIPv4.getPayload().getPayload() instanceof PnIoCm_Packet_Req)
                    .check(payloadIPv4 -> ((PnIoCm_Packet_Req) payloadIPv4.getPayload().getPayload()).getBlocks().size() == 1 && ((PnIoCm_Packet_Req) payloadIPv4.getPayload().getPayload()).getBlocks().get(0) instanceof PnIoCm_Control_Request_ApplicationReady)
                    .handle(payloadIPv4 -> {
                        DceRpc_Packet dceRpc_packet = payloadIPv4.getPayload();
                        PnIoCm_Control_Request_ApplicationReady pnIoCmBlock = (PnIoCm_Control_Request_ApplicationReady) ((PnIoCm_Packet_Req) dceRpc_packet.getPayload()).getBlocks().get(0);

                        // TODO: Save these ...
                        Uuid arUuid = pnIoCmBlock.getArUuid();
                        int sessionKey = pnIoCmBlock.getSessionKey();

                        // Send back a response, but this is just a hack ... we need to move this into the subscribe method, or we can't complete the future that we're returning.
                        RawSocketChannel pnChannel = (RawSocketChannel) conversationContext.getChannel();
                        PnDcpPacketFactory.sendApplicationReadyResponse(conversationContext, pnChannel, profinetDriverContext, payloadIPv4.getSourcePort(), dceRpc_packet.getActivityUuid(), arUuid, sessionKey);

                        connected = true;
                        // TODO: Prepare the subscription response.
//                        future.complete(new DefaultPlcSubscriptionResponse(subscriptionRequest, ))
                    });

                // Now send the ParameterEnd request and wait for a response.
                CompletableFuture<PnIoCm_Control_Response_ParameterEnd> parameterEndFuture = PnDcpPacketFactory.sendParameterEndRequest(conversationContext, rawSocketChannel, profinetDriverContext);
                parameterEndFuture.whenComplete((parameterEnd, throwable) -> {
                    // We needed to put the code to expect the ApplicationReady to subscribe before sending,
                    // as the device sends it within 4 ms, and we can't guarantee that we're done setting up
                    // the listener.
                });
            });

        return future;
    }

    @Override
    protected void decode(ConversationContext<Ethernet_Frame> context, Ethernet_Frame msg) throws Exception {
        RawSocketChannel pnChannel = (RawSocketChannel) context.getChannel();
/*        MacAddress srcAddress = new MacAddress(pnChannel.getLocalMacAddress().getAddress());
        MacAddress dstAddress = new MacAddress(pnChannel.getRemoteMacAddress().getAddress());
        byte[] data = new byte[40];
        data[0] = (byte) 0x60;
        data[1] = (byte) 0x60;
        data[2] = (byte) 0x60;
        data[3] = (byte) 0x60;
        data[4] = (byte) 0x60;
        data[6] = (byte) 0x60;*/

        if (msg.getPayload() instanceof Ethernet_FramePayload_PnDcp) {
            Ethernet_FramePayload_PnDcp dcpPacket = (Ethernet_FramePayload_PnDcp) msg.getPayload();
            if (dcpPacket.getPdu() instanceof PnDcp_Pdu_RealTimeCyclic) {
                PnDcp_Pdu_RealTimeCyclic realTimeCyclic = (PnDcp_Pdu_RealTimeCyclic) dcpPacket.getPdu();

/*                // Save the offset of the first packet.
                cycleOffset = realTimeCyclic.getCycleCounter();

                int newCycleCounter = (int) ((((System.nanoTime() - startTime) / (NANOS_PER_CYCLE)) + cycleOffset) % 65536);
                Ethernet_Frame frame = new Ethernet_Frame(
                    dstAddress,
                    srcAddress,
                    new Ethernet_FramePayload_VirtualLan(VirtualLanPriority.INTERNETWORK_CONTROL, false, (short) 0,
                        new Ethernet_FramePayload_PnDcp(
                            new PnDcp_Pdu_RealTimeCyclic(0x8003,
                                new PnIo_CyclicServiceDataUnit(data, (short) 40),
                                newCycleCounter, false, true,
                                true, true, false, true)
                        )
                    )
                );
                System.out.println("Sending");
                context.sendToWire(frame);*/

            } else {
                System.out.println(dcpPacket);
            }
        } else if (msg.getPayload() instanceof Ethernet_FramePayload_IPv4) {
            Ethernet_FramePayload_IPv4 payloadIPv4 = (Ethernet_FramePayload_IPv4) msg.getPayload();
            if (payloadIPv4.getPayload().getPayload() instanceof PnIoCm_Packet_Ping) {
                DceRpc_Packet pingPacket = payloadIPv4.getPayload();
                // Send back a ping response
                PnDcpPacketFactory.sendPingResponse(context, pnChannel, profinetDriverContext, payloadIPv4);
            }
            // The remote device terminated the connection.
            else if (payloadIPv4.getPayload().getPayload() instanceof PnIoCm_Packet_ConnectionlessCancel) {
                context.getChannel().close();
            } else {
                System.out.println(msg);
            }
        }
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
            profinetDriverContext.setDeviceType(new String(block.getDeviceVendorValue()));
        }

        if (blockMap.containsKey(ProfinetDiscoverer.DEVICE_NAME_OF_STATION)) {
            PnDcp_Block_DevicePropertiesNameOfStation block =
                (PnDcp_Block_DevicePropertiesNameOfStation) blockMap.get(ProfinetDiscoverer.DEVICE_NAME_OF_STATION);
            profinetDriverContext.setDeviceName(new String(block.getNameOfStation()));
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
            profinetDriverContext.setRoles(roles);
        }

        if (blockMap.containsKey(ProfinetDiscoverer.DEVICE_ID)) {
            PnDcp_Block_DevicePropertiesDeviceId block =
                (PnDcp_Block_DevicePropertiesDeviceId) blockMap.get(ProfinetDiscoverer.DEVICE_ID);
            profinetDriverContext.setVendorId(block.getVendorId());
            profinetDriverContext.setDeviceId(block.getDeviceId());
        }
    }

    protected int getDataTypeLengthInBytes(PlcValueType dataType) {
        switch (dataType) {
            case NULL:
                return 0;
            case BOOL:
            case BYTE:
            case USINT:
            case SINT:
            case CHAR:
                return 1;
            case WORD:
            case UINT:
            case INT:
            case WCHAR:
                return 2;
            case DWORD:
            case UDINT:
            case DINT:
            case REAL:
                return 4;
            case LWORD:
            case ULINT:
            case LINT:
            case LREAL:
                return 8;
            case STRING:
            case WSTRING:
            case TIME:
            case LTIME:
            case DATE:
            case LDATE:
            case TIME_OF_DAY:
            case LTIME_OF_DAY:
            case DATE_AND_TIME:
            case LDATE_AND_TIME:
            case Struct:
            case List:
            case RAW_BYTE_ARRAY:
                throw new PlcRuntimeException("Length undefined for type " + dataType.name());
        }
        throw new PlcRuntimeException("Length undefined");
    }

}
