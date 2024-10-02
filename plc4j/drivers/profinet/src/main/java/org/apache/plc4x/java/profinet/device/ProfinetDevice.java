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
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.profinet.context.ProfinetDeviceContext;
import org.apache.plc4x.java.profinet.gsdml.*;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.DatagramSocket;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.nio.ByteBuffer;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;
import java.util.function.BiFunction;
import java.util.function.Consumer;
import java.util.function.Function;

public class ProfinetDevice implements PlcSubscriber {

    private final Logger logger = LoggerFactory.getLogger(ProfinetDevice.class);
    private static final int DEFAULT_NUMBER_OF_PORTS_TO_SCAN = 100;
    private static final int MIN_CYCLE_NANO_SEC = 31250;
    private final BiFunction<String, String, ProfinetISO15745Profile> gsdHandler;
    private final ProfinetDeviceContext deviceContext = new ProfinetDeviceContext();
    private final MessageWrapper messageWrapper;

    // Each device should create a receiving socket, all the packets are then automatically transferred to the listener for the channel though.
    private DatagramSocket socket = null;
    private String vendorId;
    private String deviceId;
    private Thread eventLoop = null;
    Map<String, List<Consumer<PlcSubscriptionEvent>>> registrations = new HashMap<>();
    private int offset = 0;
    private boolean firstMessage = true;
    private boolean setIpAddress = false;

    public ProfinetDevice(MessageWrapper messageWrapper, String deviceName, String deviceAccess, String subModules, BiFunction<String, String, ProfinetISO15745Profile> gsdHandler) {
        this.messageWrapper = messageWrapper;
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

    private void send(ProfinetCallable<DceRpc_Packet> callable, int sourcePort, int destPort) {
        this.messageWrapper.sendUdpMessage(
            callable,
            deviceContext,
            sourcePort,
            destPort
        );
    }

    private void recordIdAndSend(ProfinetCallable<DceRpc_Packet> callable, int sourcePort, int destPort) {
        deviceContext.addToQueue(callable.getId(), callable);
        this.messageWrapper.sendUdpMessage(
            callable,
            deviceContext,
            sourcePort,
            destPort
        );
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();

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
        // If an explicit address is provided, the driver tries to explicitly configure the device to that address.
        if (this.setIpAddress) {
            deviceContext.setState(ProfinetDeviceState.SET_IP);
        }
        start();
        return true;
    }

    /*
        Starts the device main loop, sending data from controller to device.
     */
    public void start() {
        final long timeout = (long) deviceContext.getConfiguration().getReductionRatio() * deviceContext.getConfiguration().getSendClockFactor() * deviceContext.getConfiguration().getWatchdogFactor() * MIN_CYCLE_NANO_SEC;
        final int cycleTime = (deviceContext.getConfiguration().getSendClockFactor() * deviceContext.getConfiguration().getReductionRatio() * MIN_CYCLE_NANO_SEC) / 1000000;
        Function<Object, Boolean> subscription =
            message -> {
                long startTime = System.nanoTime();
                ScheduledExecutorService ses = Executors.newScheduledThreadPool(1);
                ses.scheduleAtFixedRate(() -> {
                    try {
//                        System.out.println("State: " + deviceContext.getState().name() + " queue length: " + deviceContext.getQueue().size());
                        switch (deviceContext.getState()) {
                            // If an ipAddress is specified in the device config, we use PN DCP to set the IP
                            // address of the PN device identified by the name to that given IP address.
                            case SET_IP:
                                ProfinetMessageDcpIp setIpMessage = new ProfinetMessageDcpIp();
                                messageWrapper.sendPnioMessage(setIpMessage, deviceContext);
                                deviceContext.setState(ProfinetDeviceState.IDLE);
                                break;
                            // Set up a PN-IO connection, subscribing to the stuff passed in with the connection
                            // string and also tell the device about the data we'll be publishing.
                            case IDLE:
                                CreateConnection createConnection = new CreateConnection();
                                // Send the packet and process the response ...
                                recordIdAndSend(createConnection, deviceContext.getSourcePort(), deviceContext.getDestinationPort());

                                // For some reason the first response quite often came in too late,
                                // so we're extending the wait time here.
                                createConnection.getResponseHandled().get(8 * timeout, TimeUnit.NANOSECONDS);
                                break;
                            // TODO: It seems this state is never used?
                            // It seems that in this step we would be setting parameters in the PN device (hereby configuring it)
                            // This should probably be done using the PLC4X Write API anyway.
                            case STARTUP:
                                WriteParameters writeParameters = new WriteParameters();
                                recordIdAndSend(writeParameters, deviceContext.getSourcePort(), deviceContext.getDestinationPort());
                                writeParameters.getResponseHandled().get(timeout, TimeUnit.NANOSECONDS);
                                break;
                            // Send a CONTROL packet telling the device we're done configuring the connection.
                            case PREMED:
                                WriteParametersEnd writeParametersEnd = new WriteParametersEnd();
                                recordIdAndSend(writeParametersEnd, deviceContext.getSourcePort(), deviceContext.getDestinationPort());
                                writeParametersEnd.getResponseHandled().get(timeout, TimeUnit.NANOSECONDS);
                                break;
                            // Here we're waiting for an incoming application-ready request from the device.
                            case WAITAPPLRDY:
                                break;
                            // Here we've received the application-ready request from the device and simply acknowledge
                            // it, which finishes the connection setup.
                            case APPLRDY:
                                ApplicationReadyResponse applicationReadyResponse = new ApplicationReadyResponse(deviceContext.getActivityUuid(), deviceContext.getSequenceNumber());
                                send(applicationReadyResponse, ProfinetDeviceContext.DEFAULT_UDP_PORT, deviceContext.getApplicationResponseDestinationPort());
                                deviceContext.getContext().fireConnected();
                                deviceContext.setState(ProfinetDeviceState.CYCLICDATA);
                                break;
                            // In this state we're receiving data from the remote device and in this part of the
                            // code, we're sending back our data in every cycle.
                            // TODO: Possibly check if, depending on the reduction ratio, we only have to send back data every few cycles.
                            case CYCLICDATA:
                                CyclicData cyclicData = new CyclicData(startTime);
                                messageWrapper.sendPnioMessage(cyclicData, deviceContext);
                                // TODO: Check if we're getting data every cycle ... if not, react.
                                break;
                            case ABORT:
                                // TODO: Handle this
                        }
                    } catch (InterruptedException | ExecutionException | TimeoutException e) {
                        deviceContext.setState(ProfinetDeviceState.ABORT);
                        logger.warn("Got exception", e);
                    }
                }, 0, cycleTime, TimeUnit.MILLISECONDS);
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

        // Look up the human readable text value for the given device identity
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

    public void setIpAddress(String ipAddress) {
        if (ipAddress != null) {
            this.setIpAddress = true;
            this.deviceContext.setIpAddress(ipAddress);
        }
    }

    public void handleResponse(Ethernet_FramePayload_IPv4 packet) {
        logger.debug("Received packet for {}", packet.getPayload().getObjectUuid());
        long objectId = packet.getPayload().getSequenceNumber();
        if (deviceContext.hasSequenceNumberInQueue(objectId)) {
            deviceContext.popFromQueue(objectId).handle(packet.getPayload());
        } else {
            PnIoCm_Packet payloadPacket = packet.getPayload().getPayload();
            deviceContext.setActivityUuid(packet.getPayload().getActivityUuid());
            deviceContext.setSequenceNumber(packet.getPayload().getSequenceNumber());
            if (payloadPacket instanceof PnIoCm_Packet_Req) {
                PnIoCm_Packet_Req req = (PnIoCm_Packet_Req) payloadPacket;
                deviceContext.setMaxArrayCount(req.getArrayMaximumCount());
                deviceContext.setApplicationResponseDestinationPort(packet.getSourcePort());
                for (PnIoCm_Block block : req.getBlocks()) {
                    if (block instanceof PnIoCm_Control_Request_ApplicationReady) {
                        deviceContext.setState(ProfinetDeviceState.APPLRDY);
                    }
                }
            } else if (payloadPacket instanceof PnIoCm_Packet_Fault) {
                DceRpcAck ack = new DceRpcAck(deviceContext.getActivityUuid(), deviceContext.getSequenceNumber());
                recordIdAndSend(ack, deviceContext.getSourcePort(), deviceContext.getDestinationPort());
            } else if (payloadPacket instanceof PnIoCm_Packet_Ping) {
                DceRpcAck ack = new DceRpcAck(deviceContext.getActivityUuid(), deviceContext.getSequenceNumber());
                recordIdAndSend(ack, deviceContext.getSourcePort(), deviceContext.getDestinationPort());
            } else {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Unable to match Response with Requested Profinet packet");
            }
        }
    }

    public void handle(PlcDiscoveryItem item) {
        logger.debug("Received Discovered item at device");
        if (item.getOptions().containsKey("ipAddress") && !this.setIpAddress) {
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
        Map<String, PlcResponseItem<PlcValue>> tags = new HashMap<>();
        ReadBuffer buffer = new ReadBufferByteBased(cyclicPdu.getDataUnit().getData());

        if (firstMessage) {
            offset = cyclicPdu.getCycleCounter();
            firstMessage = false;
        }

        try {
            for (ProfinetModule module : deviceContext.getModules()) {
                module.parseTags(tags, deviceContext.getDeviceName(), buffer);
            }

            Map<Consumer<PlcSubscriptionEvent>, Map<String, PlcResponseItem<PlcValue>>> response = new HashMap<>();
            for (Map.Entry<String, PlcResponseItem<PlcValue>> entry : tags.entrySet()) {
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

            for (Map.Entry<Consumer<PlcSubscriptionEvent>, Map<String, PlcResponseItem<PlcValue>>> entry : response.entrySet()) {
                entry.getKey().accept(new DefaultPlcSubscriptionEvent(Instant.now(), entry.getValue()));
            }
        } catch (ParseException e) {
            deviceContext.setState(ProfinetDeviceState.ABORT);
            logger.error("Error Parsing Cyclic Data from device {}", deviceContext.getDeviceName());
        }
    }

    public void handleAlarmResponse(PnDcp_Pdu_AlarmLow alarmPdu) {
        logger.error("Received Alarm Low packet, attempting to re-connect");
        if (alarmPdu.getVarPart()[3] == 0x18) {
            // Error from the device after not sending anything back ...
            logger.error("- AR RPC-Control Error");
        } else if (alarmPdu.getVarPart()[3] == 0x06) {
            // Switches to the non-working connection here ...
            logger.error("- AR CMI TIMEOUT ...");
        } else {
            logger.error("- Undefined alarm");
        }
        deviceContext.setState(ProfinetDeviceState.IDLE);
    }

    public void handleSetIpAddressResponse(PcDcp_GetSet_Pdu pdu) {
        deviceContext.setState(ProfinetDeviceState.IDLE);
    }

    public void setNetworkInterface(NetworkInterface networkInterface) {
        this.deviceContext.setNetworkInterface(networkInterface);
    }

    public NetworkInterface getNetworkInterface() {
        return this.deviceContext.getNetworkInterface();
    }

    public class CreateConnection implements ProfinetCallable<DceRpc_Packet> {

        final CompletableFuture<Boolean> responseHandled = new CompletableFuture<>();
        private final long id = getObjectId();

        public CompletableFuture<Boolean> getResponseHandled() {
            return responseHandled;
        }

        public long getId() {
            return id;
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
                    new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, Integer.decode("0x" + deviceId), Integer.decode("0x" + vendorId)),
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
                0xffffffffL,
                deviceContext.getConfiguration().getWatchdogFactor(),
                deviceContext.getConfiguration().getDataHoldFactor(),
                0xC000,
                ProfinetDeviceContext.DEFAULT_EMPTY_MAC_ADDRESS,
                inputApis

            ));

            blocks.add(deviceContext.getInputReq());

            List<PnIoCm_IoCrBlockReqApi> outputApis = Collections.singletonList(
                new PnIoCm_IoCrBlockReqApi(
                    deviceContext.getOutputIoPsApiBlocks(),
                    deviceContext.getOutputIoCsApiBlocks()
                )
            );

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
                0xffffffffL,
                deviceContext.getConfiguration().getWatchdogFactor(),
                deviceContext.getConfiguration().getDataHoldFactor(),
                0xC000,
                ProfinetDeviceContext.DEFAULT_EMPTY_MAC_ADDRESS,
                outputApis
            ));

            blocks.add(deviceContext.getOutputReq());

            blocks.addAll(deviceContext.getExpectedSubmoduleReq());

            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST,
                true,
                false,
                false,
                IntegerEncoding.BIG_ENDIAN,
                CharacterEncoding.ASCII,
                FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, Integer.decode("0x" + deviceId), Integer.decode("0x" + vendorId)),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                deviceContext.getUuid(),
                0,
                id,
                DceRpc_Operation.CONNECT,
                (short) 0,
                new PnIoCm_Packet_Req(ProfinetDeviceContext.DEFAULT_ARGS_MAXIMUM, ProfinetDeviceContext.DEFAULT_MAX_ARRAY_COUNT, 0, blocks)
            );
        }

        public void handle(DceRpc_Packet dceRpc_packet) {
            if ((dceRpc_packet.getOperation() == DceRpc_Operation.CONNECT) && (dceRpc_packet.getPacketType() == DceRpc_PacketType.RESPONSE)) {
                if (dceRpc_packet.getPayload().getPacketType() == DceRpc_PacketType.RESPONSE) {
                    final PnIoCm_Packet_Res connectResponse = (PnIoCm_Packet_Res) dceRpc_packet.getPayload();
                    if (connectResponse.getErrorCode() == 0) {
                        // TODO:- Re-enable the Write Parameters step if need be. Need a pcap of a simocode connection.
                        deviceContext.setState(ProfinetDeviceState.PREMED);
                        // Check the types of the block in the response match the expected ones.
                        for (PnIoCm_Block module : connectResponse.getBlocks()) {
                            // TODO: Find out what a MODULE_DIFF_BLOCK is ...
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
                    }
                } else {
                    deviceContext.setState(ProfinetDeviceState.ABORT);
                    logger.error("Received Incorrect Packet Type for Create Connection Response");
                }
            } else if (dceRpc_packet.getPacketType() == DceRpc_PacketType.REJECT) {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Device rejected connection request");
            } else {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Unexpected Response");
            }
            responseHandled.complete(true);
        }
    }

    public class WriteParameters implements ProfinetCallable<DceRpc_Packet> {

        final CompletableFuture<Boolean> responseHandled = new CompletableFuture<>();
        private final long id = getObjectId();

        public CompletableFuture<Boolean> getResponseHandled() {
            return responseHandled;
        }

        public long getId() {
            return id;
        }

        public DceRpc_Packet create() {

            int seqNumber = 0;
            List<PnIoCm_Block> requests = new ArrayList<>();

            // This will be filled in later
            requests.add(null);
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
            long multiWriteRecordLength = 0;
            // The first record isn't included in the overall length
            for (int i = 1; i < requests.size(); i++) {
                multiWriteRecordLength += requests.get(i).getLengthInBytes();
            }

            requests.set(0, new IODWriteRequestHeader(
                (short) 1,
                (short) 0,
                seqNumber,
                ProfinetDeviceContext.ARUUID,
                0x00000000,
                0x0000,
                0x0000,
                0xe040,
                multiWriteRecordLength,
                null
            ));
            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST, true, false, false,
                IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, Integer.decode("0x" + deviceId), Integer.decode("0x" + vendorId)),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                deviceContext.getUuid(),
                0,
                id,
                DceRpc_Operation.WRITE,
                (short) 0,
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

        final CompletableFuture<Boolean> responseHandled = new CompletableFuture<>();
        private final long id = getObjectId();

        public CompletableFuture<Boolean> getResponseHandled() {
            return responseHandled;
        }

        public long getId() {
            return id;
        }

        public DceRpc_Packet create() {
            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST, true, false, false,
                IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, Integer.decode("0x" + deviceId), Integer.decode("0x" + vendorId)),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                deviceContext.getUuid(),
                0,
                id,
                DceRpc_Operation.CONTROL,
                (short) 0,
                new PnIoCm_Packet_Req(16696, 16696, 0,
                    Collections.singletonList(
                        new PnIoCm_Control_Request_ParameterEnd(
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
                    } else {
                        deviceContext.setState(ProfinetDeviceState.ABORT);
                        // TODO:- Introduce the error code lookups
                        logger.error("Error {} - {} in Response from {} during Write Parameters End", connectResponse.getErrorCode1(), connectResponse.getErrorCode2(), deviceContext.getDeviceName());
                    }
                } else {
                    deviceContext.setState(ProfinetDeviceState.ABORT);
                    logger.error("Received Incorrect Packet Type for Write Parameters Ed Response");
                }
            } else if (dceRpc_packet.getPacketType() == DceRpc_PacketType.REJECT) {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Device rejected write parameter end request");
            } else {
                deviceContext.setState(ProfinetDeviceState.ABORT);
                logger.error("Unexpected Response");
            }
            responseHandled.complete(true);
        }
    }

    public class ApplicationReadyResponse implements ProfinetCallable<DceRpc_Packet> {

        private final DceRpc_ActivityUuid activityUuid;
        private final long id;

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


        public DceRpc_Packet create() {
            return new DceRpc_Packet(
                DceRpc_PacketType.RESPONSE,
                false,
                true,
                true,
                IntegerEncoding.BIG_ENDIAN,
                CharacterEncoding.ASCII,
                FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, Integer.decode("0x" + deviceId), Integer.decode("0x" + vendorId)),
                new DceRpc_InterfaceUuid_ControllerInterface(),
                activityUuid,
                0,
                id,
                DceRpc_Operation.CONTROL,
                (short) 0,
                new PnIoCm_Packet_Res(
                    (short) 0,
                    (short) 0,
                    (short) 0,
                    (short) 0,
                    deviceContext.getMaxArrayCount(),
                    0,
                    Collections.singletonList(
                        new PnIoCm_Control_Response_ApplicationReady(
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

    public class DceRpcAck implements ProfinetCallable<DceRpc_Packet> {

        private final DceRpc_ActivityUuid activityUuid;
        private final long id;

        public DceRpcAck(DceRpc_ActivityUuid activityUuid, long seqNumber) {
            this.activityUuid = activityUuid;
            this.id = seqNumber;
        }

        public CompletableFuture<Boolean> getResponseHandled() {
            return null;
        }

        public long getId() {
            return id;
        }

        public DceRpc_Packet create() {
            return new DceRpc_Packet(
                DceRpc_PacketType.NO_CALL,
                false,
                true,
                true,
                IntegerEncoding.BIG_ENDIAN,
                CharacterEncoding.ASCII,
                FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, Integer.decode("0x" + deviceId), Integer.decode("0x" + vendorId)),
                new DceRpc_InterfaceUuid_ControllerInterface(),
                activityUuid,
                0,
                id,
                DceRpc_Operation.CONTROL,
                (short) 0,
                new PnIoCm_Packet_NoCall()
            );
        }

        @Override
        public void handle(DceRpc_Packet packet) {
            logger.debug("Received an unintended packet");
        }
    }

    public class CyclicData implements ProfinetCallable<Ethernet_Frame> {

        private final long startTime;
        private final long id = getObjectId();

        public CyclicData(long startTime) {
            this.startTime = startTime;
        }

        public long getId() {
            return id;
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

                // TODO:- Still having issues with this. For the Simcode after a while we received an Alarm low message, Although it might be related to the ping functionality.
                int elapsedTime = (int) ((((System.nanoTime() - startTime) / (MIN_CYCLE_NANO_SEC)) + offset) % 65536);

                Ethernet_Frame frame = new Ethernet_Frame(
                    deviceContext.getMacAddress(),
                    deviceContext.getLocalMacAddress(),
                    new Ethernet_FramePayload_VirtualLan(
                        VirtualLanPriority.INTERNETWORK_CONTROL,
                        false,
                        (short) 0,
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

                int elapsedTime = (int) ((((System.nanoTime() - startTime) / (MIN_CYCLE_NANO_SEC)) + offset) % 65536);

                Ethernet_Frame frame = new Ethernet_Frame(
                    deviceContext.getMacAddress(),
                    deviceContext.getLocalMacAddress(),
                    new Ethernet_FramePayload_VirtualLan(
                        VirtualLanPriority.INTERNETWORK_CONTROL,
                        false,
                        (short) 0,
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
        public void handle(Ethernet_Frame packet) {
            deviceContext.setState(ProfinetDeviceState.ABORT);
            logger.error("Error Parsing Cyclic Data from device {}", deviceContext.getDeviceName());
        }
    }

    public class ProfinetMessageDcpIp implements ProfinetCallable<Ethernet_Frame> {

        private long id = getObjectId();
        private CompletableFuture<Boolean> responseHandled = new CompletableFuture<>();

        public ProfinetMessageDcpIp() {
        }

        public long getId() {
            return id;
        }

        public void setId(long id) {
            this.id = id;
        }

        public CompletableFuture<Boolean> getResponseHandled() {
            return responseHandled;
        }

        public Ethernet_Frame create() {
            Ethernet_Frame frame = null;
            try {
                frame = new Ethernet_Frame(
                    deviceContext.getMacAddress(),
                    deviceContext.getLocalMacAddress(),
                    new Ethernet_FramePayload_VirtualLan(
                        VirtualLanPriority.INTERNETWORK_CONTROL,
                        false,
                        (short) 0,
                        new Ethernet_FramePayload_PnDcp(
                            new PcDcp_GetSet_Pdu(
                                PnDcp_FrameId.DCP_GetSet_PDU.getValue(),
                                false,
                                false,
                                0x10000001L,
                                Collections.singletonList(
                                    new PnDcp_Block_IpParameter(
                                        false,
                                        false,
                                        true,
                                        deviceContext.getIpAddressAsByteArray(),
                                        deviceContext.getSubnetAsByteArray(),
                                        deviceContext.getGatewayAsByteArray()
                                    )
                                )
                            )
                        )
                    )
                );
            } catch (UnknownHostException e) {
                logger.error("Error parsing IP Address for set ip address phase");
                deviceContext.setState(ProfinetDeviceState.ABORT);
            }

            return frame;
        }

        @Override
        public void handle(Ethernet_Frame ethernetFrame) {
            logger.debug("Received a Set IP Address Response");
        }
    }
}
