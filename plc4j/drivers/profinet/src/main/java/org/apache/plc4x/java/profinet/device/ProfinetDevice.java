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

package org.apache.plc4x.java.profinet.device;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.profinet.context.ProfinetDeviceContext;
import org.apache.plc4x.java.profinet.gsdml.*;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.*;
import java.nio.ByteBuffer;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.BiFunction;
import java.util.function.Consumer;
import java.util.function.Function;

public class ProfinetDevice implements PlcSubscriber{

    private final Logger logger = LoggerFactory.getLogger(ProfinetDevice.class);
    private static final int DEFAULT_NUMBER_OF_PORTS_TO_SCAN = 100;
    private static final int MIN_CYCLE_NANO_SEC = 31250;
    private final BiFunction<String, String, ProfinetISO15745Profile> gsdHandler;
    private final ProfinetDeviceContext deviceContext = new ProfinetDeviceContext();

    // Each device should create a receiving socket, all the packets are then automatically transferred to the listener for the channel though.
    private DatagramSocket socket = null;
    private String vendorId;
    private String deviceId;
    private Thread eventLoop = null;
    Map<String, List<Consumer<PlcSubscriptionEvent>>> registrations = new HashMap<>();

    public ProfinetDevice(String deviceName, String deviceAccess, String subModules, BiFunction<String, String, ProfinetISO15745Profile> gsdHandler)  {
        this.gsdHandler = gsdHandler;
        deviceContext.setDeviceAccess(deviceAccess);
        deviceContext.setSubModules(subModules);
        deviceContext.setDeviceName(deviceName);
        openDeviceUdpPort();
    }

    private void openDeviceUdpPort() {
        // Open the receiving UDP port.
        int count = 0;
        int port = ProfinetDeviceContext.DEFAULT_SEND_UDP_PORT;
        boolean portFound = false;
        while (!portFound && count < DEFAULT_NUMBER_OF_PORTS_TO_SCAN) {
            try {
                socket = new DatagramSocket(port + count);
                portFound = true;
            } catch (SocketException e) {
                count += 1;
                port += 1;
            }
        }
        if (!portFound) {
            throw new RuntimeException("Unable to find free open port");
        }
    }

    private long getObjectId() {
        return deviceContext.getAndIncrementIdentification();
    }

    public String getVendorId() {
        return vendorId;
    }

    public String getDeviceId() {
        return deviceId;
    }

    public void setVendorDeviceId(String vendorId, String deviceId) {
        try {
            this.vendorId = vendorId;
            this.deviceId = deviceId;
            if (deviceContext.getGsdFile() == null) {
                deviceContext.setGsdFile(gsdHandler.apply(vendorId, deviceId));
            }
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
    }

    private void recordIdAndSend(ProfinetCallable<DceRpc_Packet> callable) {
        deviceContext.getQueue().put(callable.getId(), callable);
        ProfinetMessageWrapper.sendUdpMessage(
            callable,
            deviceContext
        );
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Map<String, ResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();

            return new DefaultPlcSubscriptionResponse(subscriptionRequest, values);
        });
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return null;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        // Register the current consumer for each of the given subscription handles
        for (PlcSubscriptionHandle subscriptionHandle : handles) {
            logger.debug("Registering Consumer");
            ProfinetSubscriptionHandle profinetHandle = (ProfinetSubscriptionHandle) subscriptionHandle;
            if (registrations.containsKey(profinetHandle.getAddressString())) {
                registrations.get(profinetHandle.getAddressString()).add(consumer);
            } else {
                List<Consumer<PlcSubscriptionEvent>> consumers = new ArrayList<>();
                consumers.add(consumer);
                registrations.put(profinetHandle.getAddressString(), consumers);
            }
        }
        return new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {

    }

    public boolean onConnect() throws ExecutionException, InterruptedException, TimeoutException {
        start();
        return true;
    }

    /*
        Starts the device main loop, sending data from controller to device.
     */
    public void start() {
        final long timeout = (long) deviceContext.getConfiguration().getReductionRatio() * deviceContext.getConfiguration().getSendClockFactor() * deviceContext.getConfiguration().getWatchdogFactor() * MIN_CYCLE_NANO_SEC;
        final int cycleTime = (int) (deviceContext.getConfiguration().getSendClockFactor() * deviceContext.getConfiguration().getReductionRatio() * (MIN_CYCLE_NANO_SEC/1000000.0));
        Function<Object, Boolean> subscription =
            message -> {
                long startTime = System.nanoTime();
                while (deviceContext.getState() != ProfinetDeviceState.ABORT) {
                    try {
                        switch(deviceContext.getState()) {
                            case IDLE:
                                CreateConnection createConnection = new CreateConnection();
                                recordIdAndSend(createConnection);
                                createConnection.getResponseHandled().get(timeout, TimeUnit.NANOSECONDS);
                                break;
                            case STARTUP:
                                WriteParameters writeParameters = new WriteParameters();
                                recordIdAndSend(writeParameters);
                                writeParameters.getResponseHandled().get(timeout, TimeUnit.NANOSECONDS);
                                break;
                            case PREMED:
                                WriteParametersEnd writeParametersEnd = new WriteParametersEnd();
                                recordIdAndSend(writeParametersEnd);
                                writeParametersEnd.getResponseHandled().get(timeout, TimeUnit.NANOSECONDS);
                                break;
                            case WAITAPPLRDY:
                                Thread.sleep(cycleTime);
                            case APPLRDY:
                                ApplicationReadyResponse applicationReadyResponse = new ApplicationReadyResponse(deviceContext.getActivityUuid(), deviceContext.getSequenceNumber());
                                recordIdAndSend(applicationReadyResponse);
                                deviceContext.getContext().fireConnected();
                                deviceContext.setState(ProfinetDeviceState.CYCLICDATA);
                                break;
                            case CYCLICDATA:
                                CyclicData cyclicData = new CyclicData(startTime);
                                ProfinetMessageWrapper.sendPnioMessage(cyclicData, deviceContext);
                                Thread.sleep(cycleTime);
                                break;
                        }
                    } catch (InterruptedException | ExecutionException | TimeoutException e) {
                        deviceContext.setState(ProfinetDeviceState.ABORT);
                    }
                }
                return null;
            };

        eventLoop = new Thread(new ProfinetRunnable(null, subscription));
        eventLoop.start();
    }

    /*
        Return metadata about the device. This is sourced from the connection string as well as GSD file.
     */
    public Map<String, PlcValue> getDeviceInfo() {
        Map<String, PlcValue> options = new HashMap<>();
        ProfinetDeviceIdentity deviceIdentity = this.deviceContext.getGsdFile().getProfileBody().getDeviceIdentity();
        options.put("device_id", new PlcSTRING(deviceIdentity.getDeviceID()));
        options.put("vendor_id", new PlcSTRING(deviceIdentity.getVendorId()));
        options.put("vendor_name", new PlcSTRING(deviceIdentity.getVendorName().getValue()));
        if (deviceIdentity.getInfoText() != null && deviceIdentity.getInfoText().getTextId() != null) {
            String key = deviceIdentity.getInfoText().getTextId();
            ProfinetExternalTextList externaltextList = this.deviceContext.getGsdFile().getProfileBody().getApplicationProcess().getExternalTextList();
            for (ProfinetTextIdValue s : externaltextList.getPrimaryLanguage().getText()) {
                if (key.equals(s.getTextId())) {
                    options.put("info_text", new PlcSTRING(s.getValue()));
                    break;
                }
            }
        }
        return options;
    }

    /*
        Create a structure including all the devices tags and child tags.
        The options include metadata about the device.
        The children are a list of configured submodules, with the same format as the parent.
        Each address of the children is formatted with the format i.e. parent.submodule.chiildtag
     */
    public List<PlcBrowseItem> browseTags(List<PlcBrowseItem> browseItems) {
        Map<String, PlcValue> options = getDeviceInfo();
        for (ProfinetModule module : deviceContext.getModules()) {
            browseItems = module.browseTags(browseItems, deviceContext.getDeviceName(), options);
        }
        return browseItems;
    }

    public boolean hasLldpPdu() {
        return deviceContext.isLldpReceived();
    }

    public boolean hasDcpPdu() {
        return deviceContext.isDcpReceived();
    }

    public void handleResponse(Ethernet_FramePayload_IPv4 packet) {
        logger.debug("Received packet for {}", packet.getPayload().getObjectUuid());
        long objectId = packet.getPayload().getSequenceNumber();
        if (deviceContext.getQueue().containsKey(objectId)) {
            deviceContext.getQueue().get(objectId).handle(packet.getPayload());
        } else {
            PnIoCm_Packet payloadPacket = packet.getPayload().getPayload();
            deviceContext.setActivityUuid(packet.getPayload().getActivityUuid());
            deviceContext.setSequenceNumber(packet.getPayload().getSequenceNumber());
            if (payloadPacket instanceof PnIoCm_Packet_Req) {
                PnIoCm_Packet_Req req = (PnIoCm_Packet_Req) payloadPacket;
                for (PnIoCm_Block block : req.getBlocks()) {
                    if (block instanceof PnIoCM_Block_Request) {
                        deviceContext.setState(ProfinetDeviceState.APPLRDY);
                    }
                }
            } else {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Unable to match Response with Requested Profinet packet");
            }
        }
    }

    public void handle(PlcDiscoveryItem item) {
        logger.debug("Received Discovered item at device");
        if (item.getOptions().containsKey("ipAddress")) {
            deviceContext.setIpAddress(item.getOptions().get("ipAddress"));
        }
        if (item.getOptions().containsKey("portId")) {
            deviceContext.setPortId(item.getOptions().get("portId"));
        }
        if (item.getOptions().containsKey("deviceTypeName")) {
            deviceContext.setDeviceTypeName(item.getOptions().get("deviceTypeName"));
        }
        if (item.getOptions().containsKey("vendorId") && item.getOptions().containsKey("deviceId")) {
            setVendorDeviceId(item.getOptions().get("vendorId"), item.getOptions().get("deviceId"));
        }
        if (item.getOptions().containsKey("deviceName")) {
            deviceContext.setDeviceName(item.getOptions().get("deviceName"));
        }
        if (item.getOptions().containsKey("localMacAddress")) {
            String macString = item.getOptions().get("localMacAddress").replace(":", "");
            try {
                deviceContext.setLocalMacAddress(new MacAddress(Hex.decodeHex(macString)));
            } catch (DecoderException e) {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Error Decode Local Mac Address from PN-DCP {}", item.getOptions().get("localMacAddress"));
            }
        }
        if (item.getOptions().containsKey("macAddress")) {
            String macString = item.getOptions().get("macAddress").replace(":", "");
            try {
                deviceContext.setMacAddress(new MacAddress(Hex.decodeHex(macString)));
            } catch (DecoderException e) {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Error Decode Mac Address from PN-DCP {}", item.getOptions().get("macAddress"));
            }
        }
        if (item.getOptions().containsKey("packetType")) {
            if (item.getOptions().get("packetType").equals("lldp")) {
                deviceContext.setLldpReceived(true);
            }
            if (item.getOptions().get("packetType").equals("dcp")) {
                deviceContext.setDcpReceived(true);
            }
        }
    }

    public void setContext(ConversationContext<Ethernet_Frame> context, ProfinetChannel channel) {
        deviceContext.setContext(context);
        deviceContext.setChannel(channel);
    }

    public ProfinetDeviceContext getDeviceContext() {
        return deviceContext;
    }

    public void handleRealTimeResponse(PnDcp_Pdu_RealTimeCyclic cyclicPdu) {
        Map<String, ResponseItem<PlcValue>> tags = new HashMap<>();
        ReadBuffer buffer = new ReadBufferByteBased(cyclicPdu.getDataUnit().getData());

        try {
            for (ProfinetModule module : deviceContext.getModules()) {
                module.parseTags(tags, deviceContext.getDeviceName(), buffer);
            }

            Map<Consumer<PlcSubscriptionEvent>, Map<String, ResponseItem<PlcValue>>> response = new HashMap<>();
            for (Map.Entry<String, ResponseItem<PlcValue>> entry : tags.entrySet()) {
                boolean processTag = false;
                ProfinetSubscriptionHandle handle = deviceContext.getSubscriptionHandle(entry.getKey());
                if (handle != null) {
                    if (handle.getLastValue() == null) {
                        processTag = true;
                        handle.setLastValue(entry.getValue().getValue());
                    } else if (handle.getSubscriptionType() == PlcSubscriptionType.CHANGE_OF_STATE && !entry.getValue().getValue().toString().equals(handle.getLastValue().toString())) {
                        processTag = true;
                        handle.setLastValue(entry.getValue().getValue());
                    }
                    if (handle.getSubscriptionType() == PlcSubscriptionType.CYCLIC) {
                        processTag = true;
                    }
                    if (handle.getSubscriptionType() == PlcSubscriptionType.EVENT) {
                        processTag = true;
                    }
                }
                if (registrations.containsKey(entry.getKey()) && processTag) {
                    List<Consumer<PlcSubscriptionEvent>> selectedRegistrations = registrations.get(entry.getKey());
                    for (Consumer<PlcSubscriptionEvent> reg : selectedRegistrations) {
                        if (response.containsKey(reg)) {
                            response.get(reg).put(deviceContext.getSubscriptionHandle(entry.getKey()).getTag(), entry.getValue());
                        } else {
                            response.put(reg, new HashMap<>());
                            response.get(reg).put(deviceContext.getSubscriptionHandle(entry.getKey()).getTag(), entry.getValue());
                        }
                    }
                }
            }

            for (Map.Entry<Consumer<PlcSubscriptionEvent>, Map<String, ResponseItem<PlcValue>>> entry : response.entrySet()) {
                entry.getKey().accept(new DefaultPlcSubscriptionEvent(Instant.now(), entry.getValue()));
            }
        } catch (ParseException e) {
            deviceContext.setState(ProfinetDeviceState.ABORT);
            logger.error("Error Parsing Cyclic Data from device {}", deviceContext.getDeviceName());
        }
    }

    public void handleAlarmResponse(PnDcp_Pdu_AlarmLow alarmPdu) {
        logger.error("Received Alarm Low packet, attempting to re-connect");
        deviceContext.setState(ProfinetDeviceState.IDLE);
    }

    public class CreateConnection implements ProfinetCallable<DceRpc_Packet> {

        CompletableFuture<Boolean> responseHandled = new CompletableFuture<>();
        private long id = getObjectId();
        public CompletableFuture<Boolean> getResponseHandled() {
            return responseHandled;
        }
        public long getId() {
            return id;
        }
        public void setId(long id) {
            this.id = id;
        }

        public DceRpc_Packet create() {
            deviceContext.setSessionKey(deviceContext.getAndIncrementSessionKey());

            List<PnIoCm_Block> blocks = new ArrayList<>();
            blocks.add(new PnIoCm_Block_ArReq(
                    ProfinetDeviceContext.BLOCK_VERSION_HIGH,
                    ProfinetDeviceContext.BLOCK_VERSION_LOW,
                    PnIoCm_ArType.IO_CONTROLLER,
                    ProfinetDeviceContext.ARUUID,
                    deviceContext.getSessionKey(),
                    deviceContext.getLocalMacAddress(),
                    new DceRpc_ObjectUuid((byte) 0x00, 0x0001, Integer.valueOf(deviceId), Integer.valueOf(vendorId)),
                    false,
                    deviceContext.isNonLegacyStartupMode(),
                    false,
                    false,
                    PnIoCm_CompanionArType.SINGLE_AR,
                    false,
                    true,
                    false,
                    PnIoCm_State.ACTIVE,
                    ProfinetDeviceContext.DEFAULT_ACTIVITY_TIMEOUT,
                    ProfinetDeviceContext.UDP_RT_PORT,
                    ProfinetDeviceContext.DEFAULT_PLC4X_STATION_NAME
                )
            );

            blocks.add(
                new PnIoCm_Block_AlarmCrReq(
                    (short) 1,
                    (short) 0,
                    PnIoCm_AlarmCrType.ALARM_CR,
                    0x8892,
                    false,
                    false,
                    1,
                    3,
                    0x0000,
                    200,
                    0xC000,
                    0xA000
                )
            );

            List<PnIoCm_IoCrBlockReqApi> inputApis = Collections.singletonList(
                new PnIoCm_IoCrBlockReqApi(
                    deviceContext.getInputIoPsApiBlocks(),
                    deviceContext.getInputIoCsApiBlocks())
            );

            List<PnIoCm_IoCrBlockReqApi> outputApis = Collections.singletonList(
                new PnIoCm_IoCrBlockReqApi(
                    deviceContext.getOutputIoPsApiBlocks(),
                    deviceContext.getOutputIoCsApiBlocks()
                )
            );

            deviceContext.setInputReq(new PnIoCm_Block_IoCrReq(
                (short) 1,
                (short) 0,
                PnIoCm_IoCrType.INPUT_CR,
                0x0001,
                ProfinetDeviceContext.UDP_RT_PORT,
                false,
                false,
                false,
                false,
                PnIoCm_RtClass.RT_CLASS_2,
                ProfinetDeviceContext.DEFAULT_IO_DATA_SIZE,
                deviceContext.getIncrementAndGetFrameId(),
                deviceContext.getConfiguration().getSendClockFactor(),
                deviceContext.getConfiguration().getReductionRatio(),
                1,
                0,
                0xffffffff,
                deviceContext.getConfiguration().getWatchdogFactor(),
                deviceContext.getConfiguration().getDataHoldFactor(),
                0xC000,
                ProfinetDeviceContext.DEFAULT_EMPTY_MAC_ADDRESS,
                inputApis

            ));

            blocks.add(deviceContext.getInputReq());

            deviceContext.setOutputReq(new PnIoCm_Block_IoCrReq(
                (short) 1,
                (short) 0,
                PnIoCm_IoCrType.OUTPUT_CR,
                0x0002,
                ProfinetDeviceContext.UDP_RT_PORT,
                false,
                false,
                false,
                false,
                PnIoCm_RtClass.RT_CLASS_2,
                ProfinetDeviceContext.DEFAULT_IO_DATA_SIZE,
                deviceContext.getIncrementAndGetFrameId(),
                deviceContext.getConfiguration().getSendClockFactor(),
                deviceContext.getConfiguration().getReductionRatio(),
                1,
                0,
                0xffffffff,
                deviceContext.getConfiguration().getWatchdogFactor(),
                deviceContext.getConfiguration().getDataHoldFactor(),
                0xC000,
                ProfinetDeviceContext.DEFAULT_EMPTY_MAC_ADDRESS,
                outputApis
            ));

            blocks.add(deviceContext.getOutputReq());

            for (PnIoCm_Block_ExpectedSubmoduleReq expectedSubModuleApiBlocksReq : deviceContext.getExpectedSubmoduleReq()) {
                blocks.add(expectedSubModuleApiBlocksReq);
            }

            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST,
                true,
                false,
                false,
                IntegerEncoding.BIG_ENDIAN,
                CharacterEncoding.ASCII,
                FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, 0x0001, Integer.valueOf(deviceId), Integer.valueOf(vendorId)),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                deviceContext.getUuid(),
                0,
                id,
                DceRpc_Operation.CONNECT,
                new PnIoCm_Packet_Req(ProfinetDeviceContext.DEFAULT_ARGS_MAXIMUM, ProfinetDeviceContext.DEFAULT_MAX_ARRAY_COUNT, 0, blocks)
            );
        }

        public void handle(DceRpc_Packet dceRpc_packet) {
            if ((dceRpc_packet.getOperation() == DceRpc_Operation.CONNECT) && (dceRpc_packet.getPacketType() == DceRpc_PacketType.RESPONSE)) {
                if (dceRpc_packet.getPayload().getPacketType() == DceRpc_PacketType.RESPONSE) {
                    final PnIoCm_Packet_Res connectResponse = (PnIoCm_Packet_Res) dceRpc_packet.getPayload();
                    if (connectResponse.getErrorCode() == 0) {
                        deviceContext.setState(ProfinetDeviceState.STARTUP);
                        responseHandled.complete(true);
                        for (PnIoCm_Block module : connectResponse.getBlocks()) {
                            if (module.getBlockType() == PnIoCm_BlockType.MODULE_DIFF_BLOCK) {
                                PnIoCm_Block_ModuleDiff diffModule = (PnIoCm_Block_ModuleDiff) module;
                                logger.error("Module is different to what is expected in slot {}", diffModule.getApis().get(0).getModules().get(0).getSlotNumber());
                                deviceContext.setState(ProfinetDeviceState.ABORT);
                            }
                        }
                    } else {
                        deviceContext.setState(ProfinetDeviceState.ABORT);
                        // TODO:- Introduce the error code lookups
                        logger.error("Error {} - {} in Response from {} ", connectResponse.getErrorCode1(), connectResponse.getErrorCode2(), deviceContext.getDeviceName());
                        responseHandled.complete(true);
                    }
                } else {
                    deviceContext.setState(ProfinetDeviceState.ABORT);
                    logger.error("Received Incorrect Packet Type for Create Connection Response");
                    responseHandled.complete(true);
                }
            } else if (dceRpc_packet.getPacketType() == DceRpc_PacketType.REJECT) {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Device rejected connection request");
                responseHandled.complete(true);
            } else {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Unexpected Response");
                responseHandled.complete(true);
            }
        }
    }

    public class WriteParameters implements ProfinetCallable<DceRpc_Packet> {

        CompletableFuture<Boolean> responseHandled = new CompletableFuture<>();
        private long id = getObjectId();

        public CompletableFuture<Boolean> getResponseHandled() {
            return responseHandled;
        }

        public long getId() {
            return id;
        }

        public void setId(long id) {
            this.id = id;
        }

        public DceRpc_Packet create() {

            int seqNumber = 0;
            List<PnIoCm_Block> requests = new ArrayList<>();
            requests.add(
                new IODWriteRequestHeader(
                    (short) 1,
                    (short) 0,
                    seqNumber,
                    ProfinetDeviceContext.ARUUID,
                    0x00000000,
                    0x0000,
                    0x0000,
                    0xe040,
                    180,
                    null

                ));
            seqNumber += 1;
            for (ProfinetInterfaceSubmoduleItem interfaceModule : deviceContext.getInterfaceSubModules()) {
                requests.add(
                    new IODWriteRequestHeader(
                        (short) 1,
                        (short) 0,
                        seqNumber,
                        ProfinetDeviceContext.ARUUID,
                        0x00000000,
                        0x0000,
                        interfaceModule.getSubslotNumber(),
                        0x8071,
                        12,
                        null
                    ));
                requests.add(
                    new PDInterfaceAdjust(
                        (short) 1,
                        (short) 0,
                        MultipleInterfaceModeNameOfDevice.NAME_PROVIDED_BY_LLDP
                    )
                );
                seqNumber += 1;
            }

            int index = 1;
            for (String submodule : deviceContext.getSubModules()) {
                ProfinetModuleItem foundModule = null;
                for (ProfinetModuleItem module : deviceContext.getGsdFile().getProfileBody().getApplicationProcess().getModuleList()) {
                    if (module.getId().equals(submodule)) {
                        foundModule = module;
                        break;
                    }
                }

                if (foundModule != null && foundModule.getVirtualSubmoduleList().get(0).getRecordDataList() != null) {
                    for (ProfinetParameterRecordDataItem record : foundModule.getVirtualSubmoduleList().get(0).getRecordDataList()) {
                        requests.add(
                            new IODWriteRequestHeader(
                                (short) 1,
                                (short) 0,
                                seqNumber,
                                ProfinetDeviceContext.ARUUID,
                                0x00000000,
                                index,
                                0x0001,
                                record.getIndex(),
                                record.getLength(),
                                new UserData(ByteBuffer.allocate(4).putInt(Integer.valueOf(record.getRef().getDefaultValue())).array(), (long) record.getLength())
                            ));
                        seqNumber += 1;
                    }
                }
                index += 1;
            }

            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST, true, false, false,
                IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, 0x0001, Integer.valueOf(deviceId), Integer.valueOf(vendorId)),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                deviceContext.getUuid(),
                0,
                id,
                DceRpc_Operation.WRITE,
                new PnIoCm_Packet_Req(16696, 16696, 0,
                    requests)
            );
        }

        @Override
        public void handle(DceRpc_Packet dceRpc_packet) {
            logger.debug("Received a Write Parameter Response");
            if ((dceRpc_packet.getOperation() == DceRpc_Operation.WRITE) && (dceRpc_packet.getPacketType() == DceRpc_PacketType.RESPONSE)) {
                if (dceRpc_packet.getPayload().getPacketType() == DceRpc_PacketType.RESPONSE) {
                    final PnIoCm_Packet_Res connectResponse = (PnIoCm_Packet_Res) dceRpc_packet.getPayload();
                    if (connectResponse.getErrorCode() == 0) {
                        deviceContext.setState(ProfinetDeviceState.PREMED);
                        responseHandled.complete(true);
                    } else {
                        deviceContext.setState(ProfinetDeviceState.ABORT);
                        // TODO:- Introduce the error code lookups
                        logger.error("Error {} - {} in Response from {} during Write Parameters ", connectResponse.getErrorCode1(), connectResponse.getErrorCode2(), deviceContext.getDeviceName());
                        responseHandled.complete(true);
                    }
                } else {
                    deviceContext.setState(ProfinetDeviceState.ABORT);
                    logger.error("Received Incorrect Packet Type for Write Parameters Response");
                    responseHandled.complete(true);
                }
            } else if (dceRpc_packet.getPacketType() == DceRpc_PacketType.REJECT) {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Device rejected write parameter request");
                responseHandled.complete(true);
            } else {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Unexpected Response");
                responseHandled.complete(true);
            }
        }
    }

    public class WriteParametersEnd implements ProfinetCallable<DceRpc_Packet> {

        CompletableFuture<Boolean> responseHandled = new CompletableFuture<>();
        private long id = getObjectId();

        public CompletableFuture<Boolean> getResponseHandled() {
            return responseHandled;
        }

        public long getId() {
            return id;
        }

        public void setId(long id) {
            this.id = id;
        }

        public DceRpc_Packet create() {
            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST, true, false, false,
                IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, 0x0001, Integer.valueOf(deviceId), Integer.valueOf(vendorId)),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                deviceContext.getUuid(),
                0,
                id,
                DceRpc_Operation.CONTROL,
                new PnIoCm_Packet_Req(16696, 16696, 0,
                    List.of(
                        new PnIoCm_Control_Request(
                            (short) 1,
                            (short) 0,
                            ProfinetDeviceContext.ARUUID,
                            deviceContext.getSessionKey(),
                            0x0001
                        )
                    ))
            );
        }

        @Override
        public void handle(DceRpc_Packet dceRpc_packet) {
            logger.debug("Received a Write Parameter End Response");
            if ((dceRpc_packet.getOperation() == DceRpc_Operation.CONTROL) && (dceRpc_packet.getPacketType() == DceRpc_PacketType.RESPONSE)) {
                if (dceRpc_packet.getPayload().getPacketType() == DceRpc_PacketType.RESPONSE) {
                    final PnIoCm_Packet_Res connectResponse = (PnIoCm_Packet_Res) dceRpc_packet.getPayload();
                    if (connectResponse.getErrorCode() == 0) {
                        deviceContext.setState(ProfinetDeviceState.WAITAPPLRDY);
                        responseHandled.complete(true);
                    } else {
                        deviceContext.setState(ProfinetDeviceState.ABORT);
                        // TODO:- Introduce the error code lookups
                        logger.error("Error {} - {} in Response from {} during Write Parameters End", connectResponse.getErrorCode1(), connectResponse.getErrorCode2(), deviceContext.getDeviceName());
                        responseHandled.complete(true);
                    }
                } else {
                    deviceContext.setState(ProfinetDeviceState.ABORT);
                    logger.error("Received Incorrect Packet Type for Write Parameters Ed Response");
                    responseHandled.complete(true);
                }
            } else if (dceRpc_packet.getPacketType() == DceRpc_PacketType.REJECT) {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Device rejected write parameter end request");
                responseHandled.complete(true);
            } else {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Unexpected Response");
                responseHandled.complete(true);
            }
        }
    }

    public class ApplicationReadyResponse implements ProfinetCallable<DceRpc_Packet> {

        private final DceRpc_ActivityUuid activityUuid;
        private long id;

        public ApplicationReadyResponse(DceRpc_ActivityUuid activityUuid, long seqNumber) {
            this.activityUuid = activityUuid;
            this.id = seqNumber;
        }

        public CompletableFuture<Boolean> getResponseHandled() {
            return null;
        }

        public long getId() {
            return id;
        }

        public void setId(long id) {
            this.id = id;
        }

        public DceRpc_Packet create() {
            return new DceRpc_Packet(
                DceRpc_PacketType.RESPONSE,
                false,
                true,
                true,
                IntegerEncoding.BIG_ENDIAN,
                CharacterEncoding.ASCII,
                FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, 0x0001, Integer.valueOf(deviceId), Integer.valueOf(vendorId)),
                new DceRpc_InterfaceUuid_ControllerInterface(),
                activityUuid,
                0,
                id,
                DceRpc_Operation.CONTROL,
                new PnIoCm_Packet_Res(
                    (short) 0,
                    (short) 0,
                    (short) 0,
                    (short) 0,
                    ProfinetDeviceContext.DEFAULT_MAX_ARRAY_COUNT,
                    0,
                    List.of(
                        new PnIoCM_Block_Response(
                            (short) 1,
                            (short) 0,
                            ProfinetDeviceContext.ARUUID,
                            deviceContext.getSessionKey(),
                            0x0008,
                            0x0000
                        )
                    ))
            );
        }

        @Override
        public void handle(DceRpc_Packet packet) {
            logger.debug("Received an unintented packet - We were expecting a response for an Application Ready Response");
        }
    }

    public class CyclicData implements ProfinetCallable<Ethernet_Frame> {

        private final long startTime;
        private long id = getObjectId();

        public CyclicData(long startTime) {
            this.startTime = startTime;
        }

        public long getId() {
            return id;
        }

        public void setId(long id) {
            this.id = id;
        }

        public Ethernet_Frame create() {

            WriteBufferByteBased buffer = new WriteBufferByteBased(deviceContext.getOutputReq().getDataLength());
            PnIoCm_IoCrBlockReqApi api = deviceContext.getOutputReq().getApis().get(0);
            try {
                for (PnIoCm_IoCs iocs : api.getIoCss()) {
                    PnIoCm_DataUnitIoCs ioc = new PnIoCm_DataUnitIoCs(false, (byte) 0x03, false);
                    ioc.serialize(buffer);
                }

                for (PnIoCm_IoDataObject dataObject : api.getIoDataObjects()) {
                    // TODO: Need to specify the datatype length based on the gsd file
                    PnIoCm_DataUnitDataObject ioc = new PnIoCm_DataUnitDataObject(
                        new byte[1],
                        new PnIoCm_DataUnitIoCs(false, (byte) 0x03, false),
                        1
                    );
                    ioc.serialize(buffer);
                }

                while (buffer.getPos() < deviceContext.getOutputReq().getDataLength()) {
                    buffer.writeByte((byte) 0x00);
                }

                int elapsedTime = (int) ((((System.nanoTime() - startTime)/(MIN_CYCLE_NANO_SEC))) % 65536);

                Ethernet_Frame frame = new Ethernet_Frame(
                    deviceContext.getMacAddress(),
                    deviceContext.getLocalMacAddress(),
                    new Ethernet_FramePayload_VirtualLan(
                        VirtualLanPriority.INTERNETWORK_CONTROL,
                        false,
                        0,
                        new Ethernet_FramePayload_PnDcp(
                            new PnDcp_Pdu_RealTimeCyclic(
                                deviceContext.getOutputReq().getFrameId(),
                                new PnIo_CyclicServiceDataUnit(buffer.getBytes(), (short) deviceContext.getOutputReq().getDataLength()),
                                elapsedTime,
                                false,
                                true,
                                true,
                                true,
                                false,
                                true))
                    ));
                return frame;
            } catch (SerializationException e) {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Error serializing cyclic data for device {}", deviceContext.getDeviceName());

                int elapsedTime = (int) ((((System.nanoTime() - startTime)/(MIN_CYCLE_NANO_SEC))) % 65536);

                Ethernet_Frame frame = new Ethernet_Frame(
                    deviceContext.getMacAddress(),
                    deviceContext.getLocalMacAddress(),
                    new Ethernet_FramePayload_VirtualLan(
                        VirtualLanPriority.INTERNETWORK_CONTROL,
                        false,
                        0,
                        new Ethernet_FramePayload_PnDcp(
                            new PnDcp_Pdu_RealTimeCyclic(
                                deviceContext.getOutputReq().getFrameId(),
                                new PnIo_CyclicServiceDataUnit(new byte[]{}, (short) 0),
                                elapsedTime,
                                false,
                                true,
                                true,
                                true,
                                false,
                                true))
                    ));
                return frame;
            }
        }

        @Override
        public void handle(Ethernet_Frame packet)  {
            deviceContext.setState(ProfinetDeviceState.ABORT);
            logger.error("Error Parsing Cyclic Data from device {}", deviceContext.getDeviceName());
        }
    }
}
