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
package org.apache.plc4x.java.profinet;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.gsdml.ProfinetDataItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetDeviceAccessPointItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.profinet.gsdml.ProfinetIoDataInput;
import org.apache.plc4x.java.profinet.gsdml.ProfinetIoDataOutput;
import org.apache.plc4x.java.profinet.gsdml.ProfinetModuleItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetTextIdValue;
import org.apache.plc4x.java.profinet.gsdml.ProfinetVirtualSubmoduleItem;
import org.apache.plc4x.java.profinet.packets.PnDcpPacketFactory;
import org.apache.plc4x.java.profinet.packets.ProfinetConversation;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.profinet.tag.ProfinetTagHandler;
import org.apache.plc4x.java.profinet.utils.ProfinetDataTypeMapper;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.security.SecureRandom;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.Collections;
import java.util.EnumSet;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.Consumer;
import java.util.function.Predicate;

/**
 * PROFINET connection. Ports the legacy {@code ProfinetProtocolLogic} to the
 * new SPI {@link ConnectionBase} model.
 *
 * <p>Note: matches the original driver's behavior exactly, which means the
 * subscribe path still has TODOs — the {@code PlcSubscriptionResponse} future
 * is never completed once the device's {@code ApplicationReady} arrives, and
 * cyclic IO frames are logged rather than dispatched to subscribers. Those
 * gaps were present in the source we ported from.</p>
 */
public class ProfinetConnection extends ConnectionBase<ProfinetConfiguration> implements ProfinetConversation {

    private static final Logger LOGGER = LoggerFactory.getLogger(ProfinetConnection.class);

    private ProfinetMessageCodec messageCodec;
    private final ProfinetDriverContext driverContext = new ProfinetDriverContext();
    private volatile boolean connected = false;

    /** One-shot listeners for matching incoming Ethernet frames. */
    private final CopyOnWriteArrayList<FrameListener> listeners = new CopyOnWriteArrayList<>();
    private final ScheduledExecutorService scheduler =
        Executors.newSingleThreadScheduledExecutor(r -> {
            Thread t = new Thread(r, "profinet-scheduler");
            t.setDaemon(true);
            return t;
        });

    // Cached endpoint info from the raw-socket transport.
    private MacAddress localMac;
    private MacAddress remoteMac;
    private IpAddress localIp;
    private IpAddress remoteIp;

    public ProfinetConnection(ProfinetConfiguration configuration,
                              TransportInstance<?> transportInstance,
                              AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new ProfinetTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    public boolean isConnected() {
        return connected && messageCodec != null && messageCodec.isOpen();
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        if (!(transportInstance instanceof RawSocketTransportInstance rawSocket)) {
            throw new PlcConnectionException(
                "PROFINET requires a raw-socket transport (got " + transportInstance.getClass().getName() + ")");
        }
        localMac = new MacAddress(rawSocket.getLocalMacAddress());
        remoteMac = new MacAddress(rawSocket.getRemoteMacAddress());
        byte[] localIpBytes = rawSocket.getLocalIpAddress();
        if (localIpBytes == null || localIpBytes.length != 4) {
            throw new PlcConnectionException(
                "PROFINET requires an interface with a local IPv4 address");
        }
        localIp = new IpAddress(localIpBytes);
        // Remote IP comes from the connection-string transport-config; for now we
        // pull it from the configuration's ip-address parameter (if present) and
        // otherwise default to all-zeros — the device only needs the MAC for L2
        // addressing during the implicit-communication phase.
        remoteIp = configuration.ipAddress != null
            ? new IpAddress(parseIpv4(configuration.ipAddress))
            : new IpAddress(new byte[]{0, 0, 0, 0});

        messageCodec = new ProfinetMessageCodec(transportInstance, this::handleIncomingFrame);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming PROFINET data", e);
            }
        });

        // Multi-step async handshake. Fail the future and close the channel on
        // any error; only flip {@code connected} once everything succeeded.
        try {
            doHandshake().get(60, TimeUnit.SECONDS);
            connected = true;
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("PROFINET handshake failed", e);
        }
    }

    private CompletableFuture<Void> doHandshake() {
        return PnDcpPacketFactory.sendIdentificationRequest(this, localMac, remoteMac)
            .thenCompose(identifyRes -> {
                extractBlockInfo(identifyRes.getBlocks());
                if (driverContext.getVendorId() == 0 || driverContext.getDeviceId() == 0) {
                    return CompletableFuture.failedFuture(
                        new PlcConnectionException("Unable to determine vendor-id or device-id"));
                }

                ProfinetISO15745Profile deviceProfile =
                    configuration.getGsdProfile(driverContext.getVendorId(), driverContext.getDeviceId());
                if (deviceProfile == null) {
                    return CompletableFuture.failedFuture(new PlcConnectionException(
                        "Unable to find GSD profile for device with vendor-id "
                            + driverContext.getVendorId() + " and device-id " + driverContext.getDeviceId()));
                }

                applyConfiguredDap(deviceProfile);
                if (deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().isEmpty()) {
                    return CompletableFuture.failedFuture(new PlcConnectionException(
                        "GSD descriptor doesn't contain any device access points"));
                }

                return resolveRealIdentification(deviceProfile);
            });
    }

    private void applyConfiguredDap(ProfinetISO15745Profile deviceProfile) {
        if (configuration.dapId != null) {
            for (ProfinetDeviceAccessPointItem item : deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList()) {
                if (item.getId().equalsIgnoreCase(configuration.dapId)) {
                    driverContext.setDap(item);
                    break;
                }
            }
        } else if (deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().size() == 1) {
            driverContext.setDap(deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList().get(0));
        }
    }

    /**
     * Asks the device for its RealIdentificationData (which slots/subslots it
     * actually has). If the device doesn't support that call, fall back to the
     * single-DAP path. If we can't pick a DAP at all, fail.
     */
    private CompletableFuture<Void> resolveRealIdentification(ProfinetISO15745Profile deviceProfile) {
        Map<String, String> textMapping = new HashMap<>();
        for (ProfinetTextIdValue v : deviceProfile.getProfileBody().getApplicationProcess().getExternalTextList().getPrimaryLanguage().getText()) {
            textMapping.put(v.getTextId(), v.getValue());
        }

        return PnDcpPacketFactory.sendRealIdentificationDataRequest(this, driverContext)
            .handle((realIdentificationData, error) -> {
                if (error != null) {
                    // Device doesn't support this call; rely on the
                    // already-configured DAP if there is one.
                    if (driverContext.getDap() == null) {
                        throw new PlcRuntimeException(
                            "Unable to auto-configure connection, please provide the 'dap-id' parameter");
                    }
                    return null;
                }

                indexModulesFromRealIdentification(deviceProfile, realIdentificationData, textMapping);
                if (driverContext.getDap() == null) {
                    throw new PlcRuntimeException(
                        "Unable to auto-detect the device access point, please provide 'dap-id' parameter");
                }
                return null;
            });
    }

    private void indexModulesFromRealIdentification(ProfinetISO15745Profile deviceProfile,
                                                    PnIoCm_Block_RealIdentificationData realIdentificationData,
                                                    Map<String, String> textMapping) {
        Map<Integer, Long> slotModuleIdentificationNumbers = new HashMap<>();
        Map<Integer, Map<Integer, Long>> subslotModuleIdentificationNumbers = new HashMap<>();
        for (PnIoCm_RealIdentificationApi api : realIdentificationData.getApis()) {
            for (PnIoCm_RealIdentificationApi_Slot curSlot : api.getSlots()) {
                slotModuleIdentificationNumbers.put(curSlot.getSlotNumber(), curSlot.getModuleIdentNumber());
                Map<Integer, Long> subslots = subslotModuleIdentificationNumbers
                    .computeIfAbsent(curSlot.getSlotNumber(), k -> new HashMap<>());
                for (PnIoCm_RealIdentificationApi_Subslot curSubslot : curSlot.getSubslots()) {
                    subslots.put(curSubslot.getSubslotNumber(), curSubslot.getSubmoduleIdentNumber());
                }
            }
        }

        long dapModuleIdentificationNumber = slotModuleIdentificationNumbers.getOrDefault(0, 0L);
        if (dapModuleIdentificationNumber == 0) {
            throw new PlcRuntimeException("Unable to detect device access point");
        }

        for (ProfinetDeviceAccessPointItem curDap : deviceProfile.getProfileBody().getApplicationProcess().getDeviceAccessPointList()) {
            if (parseGsdmlHex(curDap.getModuleIdentNumber()) == dapModuleIdentificationNumber) {
                if (driverContext.getDap() != null && !driverContext.getDap().getId().equals(curDap.getId())) {
                    LOGGER.warn("DAP configured in connection string differs from device-configuration.");
                }
                driverContext.setDap(curDap);
                break;
            }
        }

        Map<Integer, ProfinetModuleItem> moduleIndex = new HashMap<>();
        Map<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> submoduleIndex = new HashMap<>();
        for (Map.Entry<Integer, Long> moduleEntry : slotModuleIdentificationNumbers.entrySet()) {
            int curSlot = moduleEntry.getKey();
            if (curSlot == 0) {
                continue;
            }
            long curModuleIdentifier = moduleEntry.getValue();
            for (ProfinetModuleItem curModule : deviceProfile.getProfileBody().getApplicationProcess().getModuleList()) {
                if (parseGsdmlHex(curModule.getModuleIdentNumber()) != curModuleIdentifier) {
                    continue;
                }
                moduleIndex.put(curSlot, curModule);
                Map<Integer, Long> curSubmoduleIds = subslotModuleIdentificationNumbers.get(curSlot);
                for (Map.Entry<Integer, Long> subEntry : curSubmoduleIds.entrySet()) {
                    int curSubslot = subEntry.getKey();
                    long curSubmoduleIdent = subEntry.getValue();
                    for (ProfinetVirtualSubmoduleItem curSubmodule : curModule.getVirtualSubmoduleList()) {
                        if (parseGsdmlHex(curSubmodule.getSubmoduleIdentNumber()) != curSubmoduleIdent) {
                            continue;
                        }
                        submoduleIndex.computeIfAbsent(curSlot, k -> new HashMap<>()).put(curSubslot, curSubmodule);
                        replaceTextIds(curSubmodule, textMapping);
                    }
                }
            }
        }
        driverContext.setModuleIndex(moduleIndex);
        driverContext.setSubmoduleIndex(submoduleIndex);
    }

    private void replaceTextIds(ProfinetVirtualSubmoduleItem submodule, Map<String, String> textMapping) {
        if (submodule.getIoData().getInput() != null) {
            for (ProfinetIoDataInput in : submodule.getIoData().getInput()) {
                for (ProfinetDataItem item : in.getDataItemList()) {
                    if (textMapping.containsKey(item.getTextId())) {
                        item.setTextId(textMapping.get(item.getTextId()));
                    }
                }
            }
        }
        if (submodule.getIoData().getOutput() != null) {
            for (ProfinetIoDataOutput out : submodule.getIoData().getOutput()) {
                for (ProfinetDataItem item : out.getDataItemList()) {
                    if (textMapping.containsKey(item.getTextId())) {
                        item.setTextId(textMapping.get(item.getTextId()));
                    }
                }
            }
        }
    }

    private static long parseGsdmlHex(String s) {
        String stripped = (s.startsWith("0x") || s.startsWith("0X")) ? s.substring(2) : s;
        return Long.parseLong(stripped, 16);
    }

    @Override
    public void close() throws Exception {
        connected = false;
        scheduler.shutdownNow();
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        listeners.clear();
        super.close();
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////
    // ProfinetConversation
    /////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public MacAddress getLocalMacAddress() {
        return localMac;
    }

    @Override
    public MacAddress getRemoteMacAddress() {
        return remoteMac;
    }

    @Override
    public IpAddress getLocalIp() {
        return localIp;
    }

    @Override
    public IpAddress getRemoteIp() {
        return remoteIp;
    }

    @Override
    public void sendToWire(Ethernet_Frame frame) {
        try {
            WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[frame.getLengthInBytes()]);
            frame.serialize(buffer);
            transportInstance.write(buffer.getBytes());
        } catch (Exception e) {
            throw new PlcRuntimeException("Failed to send PROFINET frame", e);
        }
    }

    @Override
    public CompletableFuture<Ethernet_Frame> expect(Predicate<Ethernet_Frame> predicate, Duration timeout) {
        CompletableFuture<Ethernet_Frame> future = new CompletableFuture<>();
        FrameListener listener = new FrameListener(predicate, future);
        listeners.add(listener);
        ScheduledFuture<?> timeoutHandle = scheduler.schedule(() -> {
            if (listeners.remove(listener)) {
                future.completeExceptionally(new TimeoutException("Timed out waiting for matching PROFINET frame"));
            }
        }, timeout.toMillis(), TimeUnit.MILLISECONDS);
        future.whenComplete((f, t) -> {
            timeoutHandle.cancel(false);
            listeners.remove(listener);
        });
        return future;
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming dispatch
    /////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingFrame(Ethernet_Frame frame) {
        for (FrameListener listener : listeners) {
            if (listener.predicate.test(frame)) {
                if (listener.future.complete(frame)) {
                    return;
                }
            }
        }
        // No matching listener — preserve the legacy behavior of logging
        // cyclic IO / ping / unexpected frames rather than dispatching them.
        decodeUnmatched(frame);
    }

    private void decodeUnmatched(Ethernet_Frame frame) {
        // The original protocol logic only handled ping requests here and
        // logged cyclic IO frames to System.out (the subscribe path was
        // incomplete). We mirror that — ping responses must work or the
        // device will tear the AR down.
        if (frame.getPayload() instanceof Ethernet_FramePayload_IPv4 payload) {
            Object dceRpc = payload.getPayload().getPayload();
            if (dceRpc instanceof PnIoCm_Packet_Ping) {
                PnDcpPacketFactory.sendPingResponse(this, driverContext, payload);
                return;
            }
            if (dceRpc instanceof PnIoCm_Packet_ConnectionlessCancel) {
                LOGGER.info("Remote device terminated the PROFINET connection");
                try {
                    close();
                } catch (Exception ignored) {
                    // ignore
                }
                return;
            }
        }
        LOGGER.debug("Unhandled PROFINET frame {}", frame);
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////
    // Browse
    /////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        Map<String, List<PlcBrowseItem>> values = new HashMap<>();
        for (String queryName : browseRequest.getQueryNames()) {
            List<PlcBrowseItem> items = new ArrayList<>();
            for (Map.Entry<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> slotEntry : driverContext.getSubmoduleIndex().entrySet()) {
                int slot = slotEntry.getKey();
                for (Map.Entry<Integer, ProfinetVirtualSubmoduleItem> subslotEntry : slotEntry.getValue().entrySet()) {
                    int subslot = subslotEntry.getKey();
                    ProfinetVirtualSubmoduleItem subslotModule = subslotEntry.getValue();
                    addBrowseItems(items, slot, subslot, subslotModule, ProfinetTag.Direction.INPUT);
                    addBrowseItems(items, slot, subslot, subslotModule, ProfinetTag.Direction.OUTPUT);
                }
            }
            responseCodes.put(queryName, PlcResponseCode.OK);
            values.put(queryName, items);
        }
        return CompletableFuture.completedFuture(new DefaultPlcBrowseResponse(browseRequest, responseCodes, values));
    }

    private void addBrowseItems(List<PlcBrowseItem> items, int slot, int subslot,
                                ProfinetVirtualSubmoduleItem subslotModule, ProfinetTag.Direction direction) {
        if (direction == ProfinetTag.Direction.INPUT) {
            if (subslotModule.getIoData().getInput() == null) return;
            for (ProfinetIoDataInput in : subslotModule.getIoData().getInput()) {
                addDataItems(items, slot, subslot, direction, in.getDataItemList());
            }
        } else {
            if (subslotModule.getIoData().getOutput() == null) return;
            for (ProfinetIoDataOutput out : subslotModule.getIoData().getOutput()) {
                addDataItems(items, slot, subslot, direction, out.getDataItemList());
            }
        }
    }

    private void addDataItems(List<PlcBrowseItem> items, int slot, int subslot,
                              ProfinetTag.Direction direction, List<ProfinetDataItem> dataItems) {
        for (int i = 0; i < dataItems.size(); i++) {
            ProfinetDataItem dataItem = dataItems.get(i);
            ProfinetDataTypeMapper.DataTypeInformation dt = ProfinetDataTypeMapper.getPlcValueType(dataItem);
            items.add(new DefaultPlcBrowseItem(
                new ProfinetTag(slot, subslot, direction, i, dt.getPlcValueType(), dt.getNumElements()),
                dataItem.getTextId(),
                direction == ProfinetTag.Direction.OUTPUT,
                direction == ProfinetTag.Direction.INPUT,
                direction == ProfinetTag.Direction.INPUT
                    ? EnumSet.of(org.apache.plc4x.java.api.types.PlcSubscriptionType.CYCLIC)
                    : Collections.emptySet(),
                false,
                Collections.emptyList(), Collections.emptyMap(), Collections.emptyMap()));
        }
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe (incomplete — same as the original)
    /////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        if (driverContext.getDap() == null) {
            return CompletableFuture.failedFuture(new PlcConnectionException("DAP not set"));
        }

        Map<Integer, Map<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>>> slots = groupTagsBySlot(subscriptionRequest);

        // Build the giant ConnectRequest blocks. Identical layout to the
        // original implementation.
        List<PnIoCm_IoDataObject> inputMessageDataObjects = new ArrayList<>();
        List<PnIoCm_IoCs> inputMessageCs = new ArrayList<>();
        List<PnIoCm_IoDataObject> outputMessageDataObjects = new ArrayList<>();
        List<PnIoCm_IoCs> outputMessageCs = new ArrayList<>();
        int inputFrameOffset = 0;
        int outputFrameOffset = 0;

        // Hard-coded DAP objects (TODO from the original: derive from GSD).
        inputMessageDataObjects.add(new PnIoCm_IoDataObject(0x00, 0x1, inputFrameOffset++));
        outputMessageCs.add(new PnIoCm_IoCs(0x00, 0x1, outputFrameOffset++));
        inputMessageDataObjects.add(new PnIoCm_IoDataObject(0x00, 0x8000, inputFrameOffset++));
        outputMessageCs.add(new PnIoCm_IoCs(0x00, 0x8000, outputFrameOffset++));
        inputMessageDataObjects.add(new PnIoCm_IoDataObject(0x00, 0x8001, inputFrameOffset++));
        outputMessageCs.add(new PnIoCm_IoCs(0x00, 0x8001, outputFrameOffset++));
        inputMessageDataObjects.add(new PnIoCm_IoDataObject(0x00, 0x8002, inputFrameOffset++));
        outputMessageCs.add(new PnIoCm_IoCs(0x00, 0x8002, outputFrameOffset++));

        List<PnIoCm_Block_ExpectedSubmoduleReq> expectedSubmodules = new ArrayList<>();
        expectedSubmodules.add(new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0, Collections.singletonList(
            new PnIoCm_ExpectedSubmoduleBlockReqApi(0x0000, 0x00000010L, 0, Arrays.asList(
                new PnIoCm_Submodule_NoInputNoOutputData(0x0001, 0x00000001L, false, false, false, false),
                new PnIoCm_Submodule_NoInputNoOutputData(0x8000, 0x00000002L, false, false, false, false),
                new PnIoCm_Submodule_NoInputNoOutputData(0x8001, 0x00000003L, false, false, false, false),
                new PnIoCm_Submodule_NoInputNoOutputData(0x8002, 0x00000003L, false, false, false, false)
            ))
        )));

        OffsetCounters off = new OffsetCounters(inputFrameOffset, outputFrameOffset);
        for (Map.Entry<Integer, Map<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>>> slotEntry : slots.entrySet()) {
            int slotNumber = slotEntry.getKey();
            buildSlotSubmodules(slotNumber, slotEntry.getValue(),
                inputMessageDataObjects, inputMessageCs,
                outputMessageDataObjects, outputMessageCs,
                expectedSubmodules, off);
        }

        Ethernet_Frame requestEthernetFrame = buildConnectRequestFrame(
            inputMessageDataObjects, inputMessageCs,
            outputMessageDataObjects, outputMessageCs,
            expectedSubmodules);

        // Same shape as the original: we expect a Res back, then expect an
        // ApplicationReady Req from the device, ack it. The future returned to
        // the caller is intentionally left incomplete — the original code had
        // the same TODO. The connect handshake succeeds and the AR is set up.
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        CompletableFuture<Ethernet_Frame> connectResponse = expect(
            frame -> {
                Ethernet_FramePayload p = unwrapVlan(frame.getPayload());
                if (!(p instanceof Ethernet_FramePayload_IPv4 ipv4)) return false;
                return ipv4.getPayload().getPacketType() == DceRpc_PacketType.RESPONSE;
            }, Duration.ofMillis(1000));
        sendToWire(requestEthernetFrame);

        connectResponse.whenComplete((connectAck, error) -> {
            if (error != null) {
                future.completeExceptionally(error);
                return;
            }
            // Register a listener for the device's ApplicationReady request.
            CompletableFuture<Ethernet_Frame> appReady = expect(
                frame -> {
                    Ethernet_FramePayload p = unwrapVlan(frame.getPayload());
                    if (!(p instanceof Ethernet_FramePayload_IPv4 ipv4)) return false;
                    Object inner = ipv4.getPayload().getPayload();
                    if (!(inner instanceof PnIoCm_Packet_Req req)) return false;
                    return req.getBlocks().size() == 1
                        && req.getBlocks().get(0) instanceof PnIoCm_Control_Request_ApplicationReady;
                }, Duration.ofMillis(500_000));

            appReady.whenComplete((appReadyFrame, appReadyError) -> {
                if (appReadyError != null) {
                    future.completeExceptionally(appReadyError);
                    return;
                }
                Ethernet_FramePayload_IPv4 ipv4 = (Ethernet_FramePayload_IPv4) unwrapVlan(appReadyFrame.getPayload());
                DceRpc_Packet dceRpc = ipv4.getPayload();
                PnIoCm_Control_Request_ApplicationReady block =
                    (PnIoCm_Control_Request_ApplicationReady) ((PnIoCm_Packet_Req) dceRpc.getPayload()).getBlocks().get(0);

                Uuid arUuid = block.getArUuid();
                int sessionKey = block.getSessionKey();
                PnDcpPacketFactory.sendApplicationReadyResponse(this, driverContext,
                    ipv4.getSourcePort(), dceRpc.getActivityUuid(), arUuid, sessionKey);

                // TODO: complete the future with a real PlcSubscriptionResponse
                // (carry-over from the original). For now match its behavior.
                LOGGER.warn("PROFINET subscribe: connect handshake complete but the subscription " +
                    "response is not delivered (port carries over an existing TODO).");
            });

            // Fire the ParameterEnd request in parallel; we don't wait for it
            // before listening for ApplicationReady because the device sends
            // ApplicationReady within ~4ms.
            PnDcpPacketFactory.sendParameterEndRequest(this, driverContext);
        });

        return future;
    }

    private Map<Integer, Map<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>>> groupTagsBySlot(
        PlcSubscriptionRequest request) {
        Map<Integer, Map<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>>> slots = new TreeMap<>();
        for (String tagName : request.getTagNames()) {
            PlcSubscriptionTag subTag = request.getTag(tagName);
            if (!(subTag.getTag() instanceof ProfinetTag profinetTag)) {
                continue;
            }
            slots.computeIfAbsent(profinetTag.getSlot(), k -> new TreeMap<>())
                .computeIfAbsent(profinetTag.getSubSlot(), k -> new TreeMap<>())
                .computeIfAbsent(profinetTag.getDirection(), k -> new TreeMap<>())
                .put(profinetTag.getIndex(), profinetTag);
        }
        return slots;
    }

    private void buildSlotSubmodules(int slotNumber,
                                     Map<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>> subslots,
                                     List<PnIoCm_IoDataObject> inputMessageDataObjects,
                                     List<PnIoCm_IoCs> inputMessageCs,
                                     List<PnIoCm_IoDataObject> outputMessageDataObjects,
                                     List<PnIoCm_IoCs> outputMessageCs,
                                     List<PnIoCm_Block_ExpectedSubmoduleReq> expectedSubmodules,
                                     OffsetCounters off) {
        List<PnIoCm_Submodule> expectedSubmoduleData = new ArrayList<>();
        for (Map.Entry<Integer, Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>>> subslotEntry : subslots.entrySet()) {
            int subslotNumber = subslotEntry.getKey();
            Map<ProfinetTag.Direction, Map<Integer, ProfinetTag>> dirs = subslotEntry.getValue();

            int iocsLength = driverContext.getSubmoduleIndex().get(slotNumber).get(subslotNumber).getIoData().getIocsLength();
            if (iocsLength == 0) iocsLength = 1;
            int iopsLength = driverContext.getSubmoduleIndex().get(slotNumber).get(subslotNumber).getIoData().getIopsLength();
            if (iopsLength == 0) iopsLength = 1;

            PnIoCm_SubmoduleType submoduleType = PnIoCm_SubmoduleType.NO_INPUT_NO_OUTPUT_DATA;
            int inputDataLength = 0;
            int outputDataLength = 0;

            if (dirs.containsKey(ProfinetTag.Direction.INPUT)) {
                submoduleType = PnIoCm_SubmoduleType.INPUT_DATA;
                for (ProfinetTag tag : dirs.get(ProfinetTag.Direction.INPUT).values()) {
                    int dataLength = getDataTypeLengthInBytes(tag.getPlcValueType()) * tag.getNumElements();
                    inputDataLength += dataLength;
                    inputMessageDataObjects.add(new PnIoCm_IoDataObject(slotNumber, subslotNumber, off.input));
                    off.input += dataLength + iocsLength;
                    inputMessageCs.add(new PnIoCm_IoCs(slotNumber, subslotNumber, off.input));
                }
            }
            if (dirs.containsKey(ProfinetTag.Direction.OUTPUT)) {
                submoduleType = submoduleType == PnIoCm_SubmoduleType.NO_INPUT_NO_OUTPUT_DATA
                    ? PnIoCm_SubmoduleType.OUTPUT_DATA
                    : PnIoCm_SubmoduleType.INPUT_AND_OUTPUT_DATA;
                for (ProfinetTag tag : dirs.get(ProfinetTag.Direction.OUTPUT).values()) {
                    outputMessageCs.add(new PnIoCm_IoCs(slotNumber, subslotNumber, off.output));
                    off.output += iocsLength;
                    int dataLength = getDataTypeLengthInBytes(tag.getPlcValueType()) * tag.getNumElements();
                    outputDataLength += dataLength;
                    off.output += dataLength;
                    outputMessageDataObjects.add(new PnIoCm_IoDataObject(slotNumber, subslotNumber, off.output));
                    off.input += iopsLength;
                    off.output += dataLength + iopsLength;
                }
            }

            switch (submoduleType) {
                case NO_INPUT_NO_OUTPUT_DATA ->
                    expectedSubmoduleData.add(new PnIoCm_Submodule_NoInputNoOutputData(subslotNumber, 0x00000010L, false, false, false, false));
                case INPUT_DATA ->
                    expectedSubmoduleData.add(new PnIoCm_Submodule_InputData(subslotNumber, 0x00000010L, false, false, false, false, inputDataLength));
                case OUTPUT_DATA ->
                    expectedSubmoduleData.add(new PnIoCm_Submodule_OutputData(subslotNumber, 0x00000010L, false, false, false, false, outputDataLength));
                case INPUT_AND_OUTPUT_DATA ->
                    expectedSubmoduleData.add(new PnIoCm_Submodule_InputAndOutputData(subslotNumber, 0x00000010L, false, false, false, false, inputDataLength, outputDataLength));
            }
        }
        expectedSubmodules.add(new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
            Collections.singletonList(new PnIoCm_ExpectedSubmoduleBlockReqApi(slotNumber, 0x00000020L, 0x0000, expectedSubmoduleData))));
    }

    private Ethernet_Frame buildConnectRequestFrame(List<PnIoCm_IoDataObject> inputMessageDataObjects,
                                                    List<PnIoCm_IoCs> inputMessageCs,
                                                    List<PnIoCm_IoDataObject> outputMessageDataObjects,
                                                    List<PnIoCm_IoCs> outputMessageCs,
                                                    List<PnIoCm_Block_ExpectedSubmoduleReq> expectedSubmodules) {
        List<PnIoCm_Block> blocks = new ArrayList<>();
        blocks.add(new PnIoCm_Block_ArReq(
            ProfinetDriverContext.BLOCK_VERSION_HIGH, ProfinetDriverContext.BLOCK_VERSION_LOW,
            PnIoCm_ArType.IO_CONTROLLER,
            driverContext.getApplicationRelationUuid(),
            driverContext.getSessionKey(),
            localMac,
            driverContext.getCmInitiatorObjectUuid(),
            false, driverContext.isAdvancedStartupMode(), false, false,
            PnIoCm_CompanionArType.SINGLE_AR,
            false, true, false,
            PnIoCm_State.ACTIVE,
            ProfinetDriverContext.DEFAULT_ACTIVITY_TIMEOUT,
            ProfinetDriverContext.UDP_RT_PORT,
            "plc4x"));
        blocks.add(new PnIoCm_Block_AlarmCrReq(
            ProfinetDriverContext.BLOCK_VERSION_HIGH, ProfinetDriverContext.BLOCK_VERSION_LOW,
            PnIoCm_AlarmCrType.ALARM_CR,
            ProfinetDriverContext.UDP_RT_PORT,
            false, false, 1, 3, 0x0000, 200, 0xC000, 0xA000));
        if (!inputMessageDataObjects.isEmpty() || !inputMessageCs.isEmpty()) {
            blocks.add(buildIoCrBlock(PnIoCm_IoCrType.INPUT_CR, 0x0001, inputMessageDataObjects, inputMessageCs));
        }
        if (!outputMessageDataObjects.isEmpty() || !outputMessageCs.isEmpty()) {
            blocks.add(buildIoCrBlock(PnIoCm_IoCrType.OUTPUT_CR, 0x0002, outputMessageDataObjects, outputMessageCs));
        }
        blocks.addAll(expectedSubmodules);

        PnIoCm_Packet_Req request = new PnIoCm_Packet_Req(16696L, 16696L, 0L, blocks);
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, driverContext.getDeviceId(), driverContext.getVendorId()),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            driverContext.getActivityUuid(),
            0L, 1L,
            DceRpc_Operation.CONNECT,
            (short) 0, request);

        SecureRandom rand = new SecureRandom();
        Ethernet_FramePayload_IPv4 udpFrame = new Ethernet_FramePayload_IPv4(
            rand.nextInt(65536), true, false, (short) 64,
            localIp, remoteIp,
            driverContext.getLocalPort(), driverContext.getRemotePortImplicitCommunication(),
            packet);
        return new Ethernet_Frame(remoteMac, localMac, udpFrame);
    }

    private PnIoCm_Block_IoCrReq buildIoCrBlock(PnIoCm_IoCrType crType, int idBase,
                                                List<PnIoCm_IoDataObject> dataObjects, List<PnIoCm_IoCs> ioCs) {
        return new PnIoCm_Block_IoCrReq(
            ProfinetDriverContext.BLOCK_VERSION_HIGH, ProfinetDriverContext.BLOCK_VERSION_LOW,
            crType, idBase, ProfinetDriverContext.UDP_RT_PORT,
            false, false, false, false,
            PnIoCm_RtClass.RT_CLASS_2,
            ProfinetDriverContext.DEFAULT_IO_DATA_SIZE,
            0x8000 | driverContext.getAndIncrementIdentification(),
            driverContext.getSendClockFactor(),
            driverContext.getReductionRatio(),
            1, 0,
            0xffffffffL,
            driverContext.getWatchdogFactor(),
            driverContext.getDataHoldFactor(),
            0xC000,
            ProfinetDriverContext.DEFAULT_EMPTY_MAC_ADDRESS,
            Collections.singletonList(new PnIoCm_IoCrBlockReqApi(dataObjects, ioCs)));
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest request) {
        return CompletableFuture.completedFuture(new DefaultPlcUnsubscriptionResponse(request));
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer,
                                                        Collection<PlcSubscriptionHandle> handles) {
        // Cyclic IO dispatch is a TODO carried over from the original — no
        // consumers will ever receive events today.
        throw new UnsupportedOperationException("PROFINET cyclic IO dispatch is not yet implemented");
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        // no-op
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    /////////////////////////////////////////////////////////////////////////////////////////////////

    private void extractBlockInfo(List<PnDcp_Block> blocks) {
        Map<String, PnDcp_Block> blockMap = new HashMap<>();
        for (PnDcp_Block block : blocks) {
            String blockName = block.getOption().name() + "-" + block.getSuboption();
            blockMap.put(blockName, block);
        }
        PnDcp_Block block;
        if ((block = blockMap.get("DEVICE_PROPERTIES_OPTION-1")) instanceof PnDcp_Block_DevicePropertiesDeviceVendor v) {
            driverContext.setDeviceType(new String(v.getDeviceVendorValue()));
        }
        if ((block = blockMap.get("DEVICE_PROPERTIES_OPTION-2")) instanceof PnDcp_Block_DevicePropertiesNameOfStation s) {
            driverContext.setDeviceName(new String(s.getNameOfStation()));
        }
        if ((block = blockMap.get("DEVICE_PROPERTIES_OPTION-3")) instanceof PnDcp_Block_DevicePropertiesDeviceId d) {
            driverContext.setVendorId(d.getVendorId());
            driverContext.setDeviceId(d.getDeviceId());
        }
        if ((block = blockMap.get("DEVICE_PROPERTIES_OPTION-4")) instanceof PnDcp_Block_DevicePropertiesDeviceRole r) {
            List<String> roles = new ArrayList<>();
            if (r.getPnioSupervisor()) roles.add("SUPERVISOR");
            if (r.getPnioMultidevive()) roles.add("MULTIDEVICE");
            if (r.getPnioController()) roles.add("CONTROLLER");
            if (r.getPnioDevice()) roles.add("DEVICE");
            driverContext.setRoles(roles);
        }
    }

    private static byte[] parseIpv4(String ip) {
        String[] parts = ip.split("\\.");
        byte[] bytes = new byte[4];
        for (int i = 0; i < 4; i++) {
            bytes[i] = (byte) Integer.parseInt(parts[i]);
        }
        return bytes;
    }

    private static Ethernet_FramePayload unwrapVlan(Ethernet_FramePayload payload) {
        if (payload instanceof Ethernet_FramePayload_VirtualLan vlan) {
            return vlan.getPayload();
        }
        return payload;
    }

    private static int getDataTypeLengthInBytes(PlcValueType dataType) {
        return switch (dataType) {
            case NULL -> 0;
            case BOOL, BYTE, USINT, SINT, CHAR -> 1;
            case WORD, UINT, INT, WCHAR -> 2;
            case DWORD, UDINT, DINT, REAL -> 4;
            case LWORD, ULINT, LINT, LREAL -> 8;
            default -> throw new PlcRuntimeException("Length undefined for type " + dataType.name());
        };
    }

    private static class OffsetCounters {
        int input;
        int output;
        OffsetCounters(int input, int output) { this.input = input; this.output = output; }
    }

    private record FrameListener(Predicate<Ethernet_Frame> predicate, CompletableFuture<Ethernet_Frame> future) {
    }

}
