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
package org.apache.plc4x.java.ads;

import org.apache.plc4x.java.ads.configuration.AdsConfiguration;
import org.apache.plc4x.java.ads.discovery.readwrite.*;
import org.apache.plc4x.java.ads.discovery.readwrite.AmsNetId;
import org.apache.plc4x.java.ads.discovery.readwrite.Constants;
import org.apache.plc4x.java.ads.model.AdsSubscriptionHandle;
import org.apache.plc4x.java.ads.readwrite.*;
import org.apache.plc4x.java.ads.resolution.ResolvedAdsTag;
import org.apache.plc4x.java.ads.resolution.TagResolver;
import org.apache.plc4x.java.ads.resolution.ValueDecoder;
import org.apache.plc4x.java.ads.resolution.ValueEncoder;
import org.apache.plc4x.java.ads.tag.AdsTag;
import org.apache.plc4x.java.ads.tag.AdsTagHandler;
import org.apache.plc4x.java.ads.tag.DirectAdsStringTag;
import org.apache.plc4x.java.ads.tag.DirectAdsTag;
import org.apache.plc4x.java.ads.tag.SymbolicAdsTag;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.*;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.*;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.spi.values.PlcUDINT;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

/**
 * ADS over TCP connection.
 *
 * <p>This is a direct port of the original Netty/SPI1-based AdsProtocolLogic to the SPI3
 * connection model — request/response correlation now happens via the invoke-id keyed
 * pendingRequests map, single-flight throttling via {@code executeThrottled} (max-concurrent = 1),
 * and unsolicited AdsDeviceNotificationRequest packets are dispatched directly to registered
 * subscription consumers.
 *
 * <p>Logic semantics (table loading, sum-up reads/writes, recursive subscribe/unsubscribe,
 * symbol/datatype invalidation on online-/symbol-version change, browse, route-setup) are
 * preserved from the original implementation, including its known TODOs around symbolic-tag
 * resolution.
 */
public class AdsTcpConnection extends ConnectionBase<AdsConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsTcpConnection.class);

    private final AtomicLong invokeIdGenerator = new AtomicLong(1);
    private AdsTcpMessageCodec messageCodec;
    private final Map<Long, CompletableFuture<AmsTCPPacket>> pendingRequests = new ConcurrentHashMap<>();

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();
    private final ConcurrentHashMap<SymbolicAdsTag, CompletableFuture<Void>> pendingResolutionRequests = new ConcurrentHashMap<>();

    private String adsVersion;
    private String deviceName;
    private int symbolVersion;
    private long onlineVersion;
    private final Map<String, AdsSymbolTableEntry> symbolTable = new HashMap<>();
    private final Map<String, AdsDataTypeTableEntry> dataTypeTable = new HashMap<>();
    private final ReentrantLock invalidationLock = new ReentrantLock();

    public AdsTcpConnection(AdsConfiguration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new AdsTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    protected int getMaxConcurrentRequests() {
        return 1;
    }

    @Override
    public boolean isConnected() {
        return messageCodec != null && messageCodec.isOpen();
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new AdsTcpMessageCodec(transportInstance, this::handleIncomingMessage);

        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming ADS data", e);
            }
        });

        // Step 1: optional AMS route setup if username/password authentication is provided.
        // (The original implementation triggered this from onConnect; we don't have access to
        //  authentication here as ConnectionBase doesn't propagate it. Skipping unless needed.)
        // If auth needed: see setupAmsRoute() — currently disabled because authentication is
        // not surfaced through the SPI3 ConnectionBase.

        // Step 2: when the configuration asks us to load the symbol/data-type tables, do so;
        //         otherwise mark the connection as connected immediately.
        if (!getConfiguration().isLoadSymbolAndDataTypeTables()) {
            fireConnectedAuditAndState();
            return;
        }

        try {
            doConnectHandshake().get(getConfiguration().getTimeoutRequest() * 5L, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("Error during ADS connect handshake", e);
        }

        fireConnectedAuditAndState();
    }

    private void fireConnectedAuditAndState() {
        LOGGER.info("ADS TCP connection established");
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT, "ADS TCP connection established");
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    private CompletableFuture<Void> doConnectHandshake() {
        // ReadDeviceInfo
        AmsPacket readDeviceInfoRequest = new AdsReadDeviceInfoRequest(
            getConfiguration().getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
            getConfiguration().getSourceAmsNetId(), 800, ReturnCode.OK, getInvokeId());

        return sendAmsRequest(readDeviceInfoRequest, AdsReadDeviceInfoResponse.class).thenCompose(readDeviceInfoResponse -> {
            if (readDeviceInfoResponse.getResult() != ReturnCode.OK) {
                return CompletableFuture.failedFuture(new PlcConnectionException(
                    "Error reading device info. Got: " + readDeviceInfoResponse.getResult()));
            }
            adsVersion = String.format("%d.%d.%d", readDeviceInfoResponse.getMajorVersion(),
                readDeviceInfoResponse.getMinorVersion(), readDeviceInfoResponse.getVersion());
            deviceName = new String(readDeviceInfoResponse.getDevice()).trim();

            // Read online version (sym-by-name "TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt")
            AmsPacket readOnlineVersionRequest = new AdsReadWriteRequest(
                getConfiguration().getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
                getConfiguration().getSourceAmsNetId(), 800, ReturnCode.OK, getInvokeId(),
                ReservedIndexGroups.ADSIGRP_SYM_VALBYNAME.getValue(), 0L, 4L, null,
                "TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt".getBytes(StandardCharsets.UTF_8));
            return sendAmsRequest(readOnlineVersionRequest, AdsReadWriteResponse.class);
        }).thenCompose(readOnlineVersionResponse -> {
            if (readOnlineVersionResponse.getResult() != ReturnCode.OK) {
                return CompletableFuture.failedFuture(new PlcConnectionException(
                    "Error reading online version number. Got: " + readOnlineVersionResponse.getResult()));
            }
            try {
                ReadBuffer rb = new ReadBufferByteBased(readOnlineVersionResponse.getData());
                onlineVersion = rb.readUnsignedLong(32, WithOption.WithUnsignedIntegerEncoding("unsigned-binary"));
            } catch (BufferException e) {
                return CompletableFuture.failedFuture(new PlcConnectionException(
                    "Error parsing online version number response.", e));
            }
            // Read symbol version (Group 0xF008, Offset 0, length 1)
            AmsPacket readSymbolVersionRequest = new AdsReadRequest(
                getConfiguration().getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
                getConfiguration().getSourceAmsNetId(), 800, ReturnCode.OK, getInvokeId(),
                ReservedIndexGroups.ADSIGRP_SYM_VERSION.getValue(), 0L, 1L);
            return sendAmsRequest(readSymbolVersionRequest, AdsReadResponse.class);
        }).thenCompose(readSymbolVersionResponse -> {
            if (readSymbolVersionResponse.getResult() != ReturnCode.OK) {
                return CompletableFuture.failedFuture(new PlcConnectionException(
                    "Error reading symbol version number. Got: " + readSymbolVersionResponse.getResult()));
            }
            try {
                ReadBuffer rb = new ReadBufferByteBased(readSymbolVersionResponse.getData());
                symbolVersion = rb.readUnsignedInt(8, WithOption.WithUnsignedIntegerEncoding("unsigned-binary"));
            } catch (BufferException e) {
                return CompletableFuture.failedFuture(new PlcConnectionException(
                    "Error parsing symbol version number response.", e));
            }
            LOGGER.debug("Fetching sizes of symbol and datatype table sizes.");
            return readSymbolTableAndDatatypeTable();
        });
    }

    private CompletableFuture<Void> readSymbolTableAndDatatypeTable() {
        // Read sizes
        AmsPacket sizesRequest = new AdsReadRequest(
            getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
            getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(), ReturnCode.OK, getInvokeId(),
            ReservedIndexGroups.ADSIGRP_SYMBOL_AND_DATA_TYPE_SIZES.getValue(), 0L, 24L);

        return sendAmsRequest(sizesRequest, AdsReadResponse.class).thenCompose(sizesResponse -> {
            if (sizesResponse.getResult() != ReturnCode.OK) {
                return CompletableFuture.failedFuture(new PlcException(
                    "Reading data type and symbol table sizes failed: " + sizesResponse.getResult()));
            }
            final AdsTableSizes sizes;
            try {
                ReadBuffer readBuffer = new ReadBufferByteBased(sizesResponse.getData());
                sizes = AdsTableSizes.staticParse(readBuffer);
            } catch (BufferException e) {
                return CompletableFuture.failedFuture(new PlcException("Error loading the table sizes", e));
            }
            LOGGER.debug("PLC contains {} symbols and {} data-types", sizes.getSymbolCount(), sizes.getDataTypeCount());

            // Read data type table
            AmsPacket dataTypeRequest = new AdsReadRequest(
                getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
                getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(), ReturnCode.OK, getInvokeId(),
                ReservedIndexGroups.ADSIGRP_DATA_TYPE_TABLE_UPLOAD.getValue(), 0L, sizes.getDataTypeLength());
            return sendAmsRequest(dataTypeRequest, AdsReadResponse.class).thenCompose(dataTypeResponse -> {
                if (dataTypeResponse.getResult() != ReturnCode.OK) {
                    return CompletableFuture.failedFuture(new PlcException(
                        "Reading data type table failed: " + dataTypeResponse.getResult()));
                }
                ReadBuffer rb = new ReadBufferByteBased(dataTypeResponse.getData());
                for (int i = 0; i < sizes.getDataTypeCount(); i++) {
                    try {
                        AdsDataTypeTableEntry entry = AdsDataTypeTableEntry.staticParse(rb);
                        dataTypeTable.put(entry.getMainName(), entry);
                    } catch (BufferException e) {
                        return CompletableFuture.failedFuture(new RuntimeException(e));
                    }
                }

                // Read symbol table
                AmsPacket symbolRequest = new AdsReadRequest(
                    getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
                    getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(), ReturnCode.OK, getInvokeId(),
                    ReservedIndexGroups.ADSIGRP_SYM_UPLOAD.getValue(), 0L, sizes.getSymbolLength());
                return sendAmsRequest(symbolRequest, AdsReadResponse.class).thenCompose(symbolResponse -> {
                    if (symbolResponse.getResult() != ReturnCode.OK) {
                        return CompletableFuture.failedFuture(new PlcException(
                            "Reading symbol table failed: " + symbolResponse.getResult()));
                    }
                    ReadBuffer rb2 = new ReadBufferByteBased(symbolResponse.getData());
                    for (int i = 0; i < sizes.getSymbolCount(); i++) {
                        try {
                            AdsSymbolTableEntry entry = AdsSymbolTableEntry.staticParse(rb2);
                            symbolTable.put(entry.getName(), entry);
                        } catch (BufferException e) {
                            return CompletableFuture.failedFuture(new RuntimeException(e));
                        }
                    }

                    // Subscribe to online + symbol version invalidation events.
                    LinkedHashMap<String, PlcSubscriptionTag> subTags = new LinkedHashMap<>();
                    subTags.put("onlineVersion", new DefaultPlcSubscriptionTag(
                        PlcSubscriptionType.CHANGE_OF_STATE,
                        new SymbolicAdsTag("TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt",
                            PlcValueType.UDINT, Collections.emptyList()),
                        java.time.Duration.ofMillis(1000)));
                    subTags.put("symbolVersion", new DefaultPlcSubscriptionTag(
                        PlcSubscriptionType.CHANGE_OF_STATE,
                        new DirectAdsTag(0xF008, 0x0000, "USINT", 1),
                        java.time.Duration.ofMillis(1000)));

                    Consumer<PlcSubscriptionEvent> consumer = ev -> {
                        for (String tagName : ev.getTagNames()) {
                            switch (tagName) {
                                case "onlineVersion": {
                                    long newVersion = ev.getPlcValue("onlineVersion").getLong();
                                    if (onlineVersion != newVersion) {
                                        if (invalidationLock.tryLock()) {
                                            LOGGER.info("Detected change of the 'online-version', invalidating data type and symbol information.");
                                            readSymbolTableAndDatatypeTable().whenComplete((u, t) -> {
                                                if (t != null) LOGGER.error("Error reloading data type and symbol data", t);
                                                invalidationLock.unlock();
                                            });
                                        }
                                    }
                                    break;
                                }
                                case "symbolVersion": {
                                    int newVersion = ev.getPlcValue("symbolVersion").getInteger();
                                    if (symbolVersion != newVersion) {
                                        if (invalidationLock.tryLock()) {
                                            LOGGER.info("Detected change of the 'symbol-version', invalidating data type and symbol information.");
                                            readSymbolTableAndDatatypeTable().whenComplete((u, t) -> {
                                                if (t != null) LOGGER.error("Error reloading data type and symbol data", t);
                                                invalidationLock.unlock();
                                            });
                                        }
                                    }
                                    break;
                                }
                            }
                        }
                    };

                    DefaultPlcSubscriptionRequest bootstrapReq = new DefaultPlcSubscriptionRequest(
                        this, subTags, consumer, Collections.emptyMap());
                    return subscribe(bootstrapReq).thenAccept(resp -> {
                        // Wire up the global consumer to the returned handles.
                        registerConsumer(consumer, resp.getSubscriptionHandles());
                    });
                });
            });
        });
    }

    @Override
    public void close() throws Exception {
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        pendingRequests.values().forEach(future ->
            future.completeExceptionally(new PlcRuntimeException("Connection closed")));
        pendingRequests.clear();
        super.close();
        LOGGER.info("ADS TCP connection closed");
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");

        PlcRuntimeException exception = new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause);
        if (!pendingRequests.isEmpty()) {
            LOGGER.warn("Failing {} pending requests due to transport disconnect", pendingRequests.size());
            pendingRequests.values().forEach(f -> f.completeExceptionally(exception));
            pendingRequests.clear();
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Request/response correlation
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private long getInvokeId() {
        long invokeId = invokeIdGenerator.getAndIncrement();
        if (invokeIdGenerator.get() == 0xFFFFFFFFL) {
            invokeIdGenerator.set(1);
        }
        return invokeId;
    }

    private void handleIncomingMessage(AmsTCPPacket packet) {
        AmsPacket userdata = packet.getUserdata();

        // Unsolicited device-notification: dispatch to consumers.
        if (userdata instanceof AdsDeviceNotificationRequest notification) {
            try {
                dispatchDeviceNotification(notification);
            } catch (Exception e) {
                LOGGER.error("Error dispatching ADS device notification", e);
            }
            return;
        }

        long invokeId = userdata.getInvokeId();
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.INCOMING_MESSAGE, "Received ADS response, invokeId=" + invokeId);
        }
        CompletableFuture<AmsTCPPacket> future = pendingRequests.remove(invokeId);
        if (future != null) {
            future.complete(packet);
        } else {
            LOGGER.warn("Received response for unknown invoke ID: {}", invokeId);
        }
    }

    private CompletableFuture<AmsTCPPacket> sendAmsTCPPacket(AmsTCPPacket request, long invokeId) {
        CompletableFuture<AmsTCPPacket> response = new CompletableFuture<>();
        pendingRequests.put(invokeId, response);
        try {
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.OUTGOING_MESSAGE, "Sending ADS request, invokeId=" + invokeId);
            }
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingRequests.remove(invokeId);
            response.completeExceptionally(new PlcRuntimeException("Failed to send request", e));
            return response;
        }
        long timeoutMs = getConfiguration().getTimeoutRequest();
        response.orTimeout(timeoutMs, TimeUnit.MILLISECONDS).whenComplete((r, e) -> {
            if (e instanceof TimeoutException) {
                pendingRequests.remove(invokeId);
            }
        });
        return response;
    }

    @SuppressWarnings("unchecked")
    private <T extends AmsPacket> CompletableFuture<T> sendAmsRequest(AmsPacket request, Class<T> expectedResponseType) {
        long invokeId = request.getInvokeId();
        AmsTCPPacket tcpPacket = new AmsTCPPacket(request);
        return executeThrottled(() -> sendAmsTCPPacket(tcpPacket, invokeId).thenApply(resp -> {
            AmsPacket userdata = resp.getUserdata();
            if (!expectedResponseType.isInstance(userdata)) {
                throw new PlcRuntimeException("Unexpected response type: " + userdata.getClass().getSimpleName()
                    + " (expected " + expectedResponseType.getSimpleName() + ")");
            }
            return (T) userdata;
        }));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Ping
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        AmsPacket pingPacket = new AdsReadDeviceInfoRequest(
            getConfiguration().getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
            getConfiguration().getSourceAmsNetId(), 800, ReturnCode.OK, getInvokeId());
        return sendAmsRequest(pingPacket, AdsReadDeviceInfoResponse.class)
            .thenApply(r -> (PlcPingResponse) new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        // Resolve every tag against the currently-loaded symbol/data-type tables. We don't go
        // through the legacy ADSIGRP_SYM_HNDBYNAME round-trip — direct (group, offset, size)
        // resolution from the locally-cached tables is enough and lets us support the full
        // path syntax (".field", "[i]", multi-dim, mixed).
        TagResolver resolver = new TagResolver(symbolTable, dataTypeTable);
        Map<String, ResolvedAdsTag> resolved = new LinkedHashMap<>();
        Map<String, PlcResponseCode> initialFailures = new HashMap<>();
        for (String tagName : readRequest.getTagNames()) {
            PlcTag tag = readRequest.getTag(tagName);
            try {
                resolved.put(tagName, resolveForReadOrWrite(resolver, tag));
            } catch (PlcInvalidTagException e) {
                LOGGER.debug("Cannot resolve tag {}: {}", tagName, e.getMessage());
                initialFailures.put(tagName, PlcResponseCode.INVALID_ADDRESS);
            } catch (Exception e) {
                LOGGER.warn("Error resolving tag {}", tagName, e);
                initialFailures.put(tagName, PlcResponseCode.INTERNAL_ERROR);
            }
        }
        if (resolved.isEmpty()) {
            return CompletableFuture.completedFuture(buildFailureReadResponse(readRequest, initialFailures));
        }
        return executeRead(readRequest, resolved, initialFailures);
    }

    private ResolvedAdsTag resolveForReadOrWrite(TagResolver resolver, PlcTag tag) {
        if (tag == null) {
            throw new PlcInvalidTagException("Tag could not be parsed");
        }
        if (tag instanceof SymbolicAdsTag s) {
            return resolver.resolve(s);
        }
        if (tag instanceof DirectAdsStringTag s) {
            return new ResolvedAdsTag(s.getIndexGroup(), s.getIndexOffset(),
                computeDirectSize(s.getPlcDataType(), s.getStringLength(), s.getNumberOfElements()),
                s.getPlcDataType(), TagResolver.plcValueTypeForName(s.getPlcDataType(), null),
                s.getStringLength(), Collections.emptyList());
        }
        if (tag instanceof DirectAdsTag d) {
            return new ResolvedAdsTag(d.getIndexGroup(), d.getIndexOffset(),
                computeDirectSize(d.getPlcDataType(), 0, d.getNumberOfElements()),
                d.getPlcDataType(), TagResolver.plcValueTypeForName(d.getPlcDataType(), null),
                0, Collections.emptyList());
        }
        throw new PlcInvalidTagException("Unsupported tag type: " + tag.getClass().getName());
    }

    private long computeDirectSize(String typeName, int stringLength, int numberOfElements) {
        Optional<AdsDataTypeTableEntry> dt = getDataTypeTableEntry(typeName);
        long size = dt.map(AdsDataTypeTableEntry::getSize).orElseGet(() -> {
            if ("STRING".equals(typeName)) return (long) (stringLength + 1);
            if ("WSTRING".equals(typeName)) return (long) (stringLength + 1) * 2;
            try {
                return (long) AdsDataType.valueOf(typeName).getNumBytes();
            } catch (IllegalArgumentException ignored) {
                return 0L;
            }
        });
        return size * Math.max(1, numberOfElements);
    }

    private PlcReadResponse buildFailureReadResponse(PlcReadRequest req, Map<String, PlcResponseCode> codes) {
        Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
        for (String name : req.getTagNames()) {
            PlcResponseCode code = codes.getOrDefault(name, PlcResponseCode.INTERNAL_ERROR);
            values.put(name, new DefaultPlcResponseItem<>(code, null));
        }
        return new DefaultPlcReadResponse(req, values);
    }

    private CompletableFuture<PlcReadResponse> executeRead(PlcReadRequest readRequest,
                                                           Map<String, ResolvedAdsTag> resolved,
                                                           Map<String, PlcResponseCode> initialFailures) {
        if (resolved.size() == 1) {
            Map.Entry<String, ResolvedAdsTag> only = resolved.entrySet().iterator().next();
            return singleRead(readRequest, only.getKey(), only.getValue(), initialFailures);
        }
        return multiRead(readRequest, resolved, initialFailures);
    }

    private CompletableFuture<PlcReadResponse> singleRead(PlcReadRequest readRequest,
                                                          String tagName, ResolvedAdsTag tag,
                                                          Map<String, PlcResponseCode> initialFailures) {
        AmsPacket request = new AdsReadRequest(
            getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
            getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(),
            ReturnCode.OK, getInvokeId(),
            tag.indexGroup(), tag.indexOffset(), tag.sizeInBytes());

        return sendAmsRequest(request, AdsReadResponse.class).thenApply(response -> {
            Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
            for (String name : readRequest.getTagNames()) {
                if (initialFailures.containsKey(name)) {
                    values.put(name, new DefaultPlcResponseItem<>(initialFailures.get(name), null));
                    continue;
                }
                if (!name.equals(tagName)) {
                    // Shouldn't happen — single-read path means exactly one resolved tag.
                    values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                    continue;
                }
                if (response.getResult() != ReturnCode.OK) {
                    values.put(name, new DefaultPlcResponseItem<>(parsePlcResponseCode(response.getResult()), null));
                    continue;
                }
                values.put(name, decodePayload(tag, response.getData()));
            }
            return new DefaultPlcReadResponse(readRequest, values);
        });
    }

    private CompletableFuture<PlcReadResponse> multiRead(PlcReadRequest readRequest,
                                                         Map<String, ResolvedAdsTag> resolved,
                                                         Map<String, PlcResponseCode> initialFailures) {
        // Tags participating in the sum-up read, in request order.
        List<String> orderedNames = readRequest.getTagNames().stream()
            .filter(resolved::containsKey).collect(Collectors.toList());

        long expectedDataSize = orderedNames.stream()
            .mapToLong(n -> 4L + resolved.get(n).sizeInBytes())  // 4-byte per-item return code + data
            .sum();

        AmsPacket request = new AdsReadWriteRequest(
            getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
            getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(),
            ReturnCode.OK, getInvokeId(),
            ReservedIndexGroups.ADSIGRP_MULTIPLE_READ.getValue(), (long) orderedNames.size(),
            expectedDataSize,
            orderedNames.stream().map(n -> {
                ResolvedAdsTag t = resolved.get(n);
                return new AdsMultiRequestItemRead(t.indexGroup(), t.indexOffset(), t.sizeInBytes());
            }).collect(Collectors.toList()),
            null);

        return sendAmsRequest(request, AdsReadWriteResponse.class).thenApply(response -> {
            Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
            if (response.getResult() != ReturnCode.OK) {
                for (String name : readRequest.getTagNames()) {
                    PlcResponseCode code = initialFailures.getOrDefault(name, parsePlcResponseCode(response.getResult()));
                    values.put(name, new DefaultPlcResponseItem<>(code, null));
                }
                return new DefaultPlcReadResponse(readRequest, values);
            }

            // Sum-up response layout: N x ReturnCode (4 bytes each), then concatenated data.
            byte[] data = response.getData();
            ReadBufferByteBased rb = newLittleEndianBuffer(data);
            // Read per-item return codes.
            Map<String, ReturnCode> resultPerTag = new LinkedHashMap<>();
            try {
                for (String n : orderedNames) {
                    resultPerTag.put(n, ReturnCode.enumForValue(rb.readUnsignedLong(32)));
                }
            } catch (BufferException e) {
                return buildFailureReadResponse(readRequest, Collections.emptyMap());
            }
            // Now decode each item's payload — for items whose result wasn't OK, we still need
            // to advance the cursor by sizeInBytes to keep alignment for the rest.
            for (String name : readRequest.getTagNames()) {
                if (initialFailures.containsKey(name)) {
                    values.put(name, new DefaultPlcResponseItem<>(initialFailures.get(name), null));
                    continue;
                }
                ResolvedAdsTag tag = resolved.get(name);
                ReturnCode code = resultPerTag.get(name);
                if (code != ReturnCode.OK) {
                    skipBytes(rb, tag.sizeInBytes());
                    values.put(name, new DefaultPlcResponseItem<>(parsePlcResponseCode(code), null));
                    continue;
                }
                try {
                    int posBefore = rb.getPositionInBits() / 8;
                    PlcValue v = new ValueDecoder(dataTypeTable).decode(rb, tag);
                    int consumed = (rb.getPositionInBits() / 8) - posBefore;
                    long expected = tag.sizeInBytes();
                    if (consumed < expected) {
                        // Top up to match the per-item size to stay aligned.
                        skipBytes(rb, expected - consumed);
                    }
                    values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.OK, v));
                } catch (Exception e) {
                    LOGGER.warn("Error decoding tag {}", name, e);
                    skipBytes(rb, tag.sizeInBytes());
                    values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
            }
            return new DefaultPlcReadResponse(readRequest, values);
        });
    }

    private PlcResponseItem<PlcValue> decodePayload(ResolvedAdsTag tag, byte[] data) {
        try {
            ReadBufferByteBased rb = newLittleEndianBuffer(data);
            PlcValue v = new ValueDecoder(dataTypeTable).decode(rb, tag);
            return new DefaultPlcResponseItem<>(PlcResponseCode.OK, v);
        } catch (Exception e) {
            LOGGER.warn("Error decoding response", e);
            return new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null);
        }
    }

    private static ReadBufferByteBased newLittleEndianBuffer(byte[] data) {
        return new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
    }

    private static void skipBytes(ReadBuffer rb, long bytes) {
        try {
            if (bytes > 0) rb.readBits((int) bytes * 8);
        } catch (BufferException ignored) {
            // best-effort; cursor remains where it was
        }
    }

    private PlcResponseCode parsePlcResponseCode(ReturnCode adsResult) {
        if (adsResult == ReturnCode.OK) return PlcResponseCode.OK;
        if (adsResult == ReturnCode.ADSERR_DEVICE_SYMBOLNOTFOUND) return PlcResponseCode.INVALID_ADDRESS;
        return PlcResponseCode.INTERNAL_ERROR;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        TagResolver resolver = new TagResolver(symbolTable, dataTypeTable);
        Map<String, ResolvedAdsTag> resolved = new LinkedHashMap<>();
        Map<String, PlcResponseCode> initialFailures = new HashMap<>();
        for (String tagName : writeRequest.getTagNames()) {
            PlcTag tag = writeRequest.getTag(tagName);
            try {
                resolved.put(tagName, resolveForReadOrWrite(resolver, tag));
            } catch (PlcInvalidTagException e) {
                LOGGER.debug("Cannot resolve tag {}: {}", tagName, e.getMessage());
                initialFailures.put(tagName, PlcResponseCode.INVALID_ADDRESS);
            } catch (Exception e) {
                LOGGER.warn("Error resolving tag {}", tagName, e);
                initialFailures.put(tagName, PlcResponseCode.INTERNAL_ERROR);
            }
        }
        if (resolved.isEmpty()) {
            Map<String, PlcResponseCode> codes = new LinkedHashMap<>();
            for (String name : writeRequest.getTagNames()) {
                codes.put(name, initialFailures.getOrDefault(name, PlcResponseCode.INTERNAL_ERROR));
            }
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, codes));
        }
        return executeWrite(writeRequest, resolved, initialFailures);
    }

    private CompletableFuture<PlcWriteResponse> executeWrite(PlcWriteRequest writeRequest,
                                                             Map<String, ResolvedAdsTag> resolved,
                                                             Map<String, PlcResponseCode> initialFailures) {
        // Serialize each tag's payload up-front so we can size the request accurately and
        // surface serialization failures cleanly before going on the wire.
        Map<String, byte[]> serialized = new LinkedHashMap<>();
        for (Map.Entry<String, ResolvedAdsTag> e : resolved.entrySet()) {
            String name = e.getKey();
            ResolvedAdsTag tag = e.getValue();
            try {
                serialized.put(name, serializeValue(tag, writeRequest.getPlcValue(name)));
            } catch (Exception ex) {
                LOGGER.warn("Error serializing tag {}", name, ex);
                initialFailures.put(name, PlcResponseCode.INVALID_DATA);
            }
        }
        // Drop any that failed to serialize.
        serialized.keySet().forEach(name -> {
            if (initialFailures.containsKey(name)) {
                serialized.remove(name);
            }
        });
        if (serialized.isEmpty()) {
            Map<String, PlcResponseCode> codes = new LinkedHashMap<>();
            for (String name : writeRequest.getTagNames()) {
                codes.put(name, initialFailures.getOrDefault(name, PlcResponseCode.INTERNAL_ERROR));
            }
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, codes));
        }
        if (serialized.size() == 1) {
            Map.Entry<String, byte[]> only = serialized.entrySet().iterator().next();
            return singleWrite(writeRequest, only.getKey(), resolved.get(only.getKey()), only.getValue(), initialFailures);
        }
        return multiWrite(writeRequest, resolved, serialized, initialFailures);
    }

    private byte[] serializeValue(ResolvedAdsTag tag, PlcValue value) throws BufferException {
        if (value == null) {
            throw new BufferException("Null PlcValue for tag");
        }
        WriteBufferByteBased wb = new WriteBufferByteBased(new byte[(int) tag.sizeInBytes()],
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
        new ValueEncoder(dataTypeTable).encode(wb, tag, value);
        return wb.getBytes();
    }

    private CompletableFuture<PlcWriteResponse> singleWrite(PlcWriteRequest writeRequest,
                                                            String tagName, ResolvedAdsTag tag, byte[] data,
                                                            Map<String, PlcResponseCode> initialFailures) {
        AmsPacket request = new AdsWriteRequest(
            getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
            getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(),
            ReturnCode.OK, getInvokeId(),
            tag.indexGroup(), tag.indexOffset(), data);
        return sendAmsRequest(request, AdsWriteResponse.class).thenApply(response -> {
            Map<String, PlcResponseCode> codes = new LinkedHashMap<>();
            for (String name : writeRequest.getTagNames()) {
                if (initialFailures.containsKey(name)) {
                    codes.put(name, initialFailures.get(name));
                } else if (name.equals(tagName)) {
                    codes.put(name, parsePlcResponseCode(response.getResult()));
                } else {
                    codes.put(name, PlcResponseCode.INTERNAL_ERROR);
                }
            }
            return new DefaultPlcWriteResponse(writeRequest, codes);
        });
    }

    private CompletableFuture<PlcWriteResponse> multiWrite(PlcWriteRequest writeRequest,
                                                           Map<String, ResolvedAdsTag> resolved,
                                                           Map<String, byte[]> serialized,
                                                           Map<String, PlcResponseCode> initialFailures) {
        // Concatenate per-tag payloads in request order.
        int totalDataBytes = serialized.values().stream().mapToInt(b -> b.length).sum();
        byte[] payload = new byte[totalDataBytes];
        int p = 0;
        List<String> names = new ArrayList<>(serialized.keySet());
        for (String name : names) {
            byte[] b = serialized.get(name);
            System.arraycopy(b, 0, payload, p, b.length);
            p += b.length;
        }

        AmsPacket request = new AdsReadWriteRequest(
            getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
            getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(),
            ReturnCode.OK, getInvokeId(),
            ReservedIndexGroups.ADSIGRP_MULTIPLE_WRITE.getValue(), (long) names.size(),
            (long) names.size() * 4L,
            names.stream().map(n -> {
                ResolvedAdsTag t = resolved.get(n);
                return new AdsMultiRequestItemWrite(t.indexGroup(), t.indexOffset(), t.sizeInBytes());
            }).collect(Collectors.toList()),
            payload);

        return sendAmsRequest(request, AdsReadWriteResponse.class).thenApply(response -> {
            Map<String, PlcResponseCode> codes = new LinkedHashMap<>();
            if (response.getResult() != ReturnCode.OK) {
                for (String n : writeRequest.getTagNames()) {
                    codes.put(n, initialFailures.getOrDefault(n, parsePlcResponseCode(response.getResult())));
                }
                return new DefaultPlcWriteResponse(writeRequest, codes);
            }
            ReadBufferByteBased rb = newLittleEndianBuffer(response.getData());
            Map<String, ReturnCode> perTag = new LinkedHashMap<>();
            try {
                for (String n : names) {
                    perTag.put(n, ReturnCode.enumForValue(rb.readUnsignedLong(32)));
                }
            } catch (BufferException e) {
                for (String n : writeRequest.getTagNames()) {
                    codes.put(n, initialFailures.getOrDefault(n, PlcResponseCode.INTERNAL_ERROR));
                }
                return new DefaultPlcWriteResponse(writeRequest, codes);
            }
            for (String n : writeRequest.getTagNames()) {
                if (initialFailures.containsKey(n)) {
                    codes.put(n, initialFailures.get(n));
                } else if (perTag.containsKey(n)) {
                    codes.put(n, parsePlcResponseCode(perTag.get(n)));
                } else {
                    codes.put(n, PlcResponseCode.INTERNAL_ERROR);
                }
            }
            return new DefaultPlcWriteResponse(writeRequest, codes);
        });
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe / Unsubscribe
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        List<PlcTag> innerTags = subscriptionRequest.getTags().stream()
            .map(t -> ((PlcSubscriptionTag) t).getTag())
            .collect(Collectors.toList());

        return getDirectAddresses(innerTags).thenCompose(resolvedTags -> {
            if (resolvedTags == null) {
                return CompletableFuture.failedFuture(new PlcException("Tags are null"));
            }
            return executeSubscribe(subscriptionRequest, resolvedTags);
        });
    }

    private CompletableFuture<PlcSubscriptionResponse> executeSubscribe(PlcSubscriptionRequest subscriptionRequest, Map<AdsTag, DirectAdsTag> resolvedTags) {
        List<AmsTCPPacket> amsTCPPackets = new ArrayList<>();
        for (PlcSubscriptionTag t : subscriptionRequest.getTags()) {
            String dataTypeName = resolvedTags.get((AdsTag) t.getTag()).getPlcDataType();
            Optional<AdsDataTypeTableEntry> dataTypeOpt = getDataTypeTableEntry(dataTypeName);
            if (dataTypeOpt.isEmpty()) {
                amsTCPPackets.add(null);
                continue;
            }
            AdsDataTypeTableEntry dataType = dataTypeOpt.get();
            DirectAdsTag directAdsTag = getDirectAdsTagForSymbolicName(t.getTag());
            // TODO: We should implement multi-dimensional arrays here ...
            int numberOfElements = (t.getArrayInfo().isEmpty()) ? 1 : t.getArrayInfo().get(0).getSize();
            amsTCPPackets.add(new AmsTCPPacket(new AdsAddDeviceNotificationRequest(
                getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
                getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(), ReturnCode.OK, getInvokeId(),
                directAdsTag.getIndexGroup(), directAdsTag.getIndexOffset(),
                dataType.getSize() * numberOfElements,
                t.getPlcSubscriptionType() == PlcSubscriptionType.CYCLIC ? AdsTransMode.CYCLIC : AdsTransMode.ON_CHANGE,
                0L,
                t.getDuration().orElse(java.time.Duration.ZERO).toMillis())));
        }

        Map<String, PlcResponseItem<PlcSubscriptionHandle>> responses = new LinkedHashMap<>();
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        Iterator<String> nameIt = subscriptionRequest.getTagNames().iterator();
        Iterator<AmsTCPPacket> packetIt = amsTCPPackets.iterator();
        subscribeRecursively(subscriptionRequest, nameIt, resolvedTags, responses, future, packetIt);
        return future;
    }

    private void subscribeRecursively(PlcSubscriptionRequest subscriptionRequest,
                                      Iterator<String> tagNames,
                                      Map<AdsTag, DirectAdsTag> resolvedTags,
                                      Map<String, PlcResponseItem<PlcSubscriptionHandle>> responses,
                                      CompletableFuture<PlcSubscriptionResponse> future,
                                      Iterator<AmsTCPPacket> amsTCPPackets) {
        if (!amsTCPPackets.hasNext()) {
            future.complete(new DefaultPlcSubscriptionResponse(subscriptionRequest, responses));
            return;
        }
        AmsTCPPacket packet = amsTCPPackets.next();
        boolean hasMorePackets = amsTCPPackets.hasNext();
        String tagName = tagNames.next();
        if (packet == null) {
            responses.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
            subscribeRecursively(subscriptionRequest, tagNames, resolvedTags, responses, future, amsTCPPackets);
            return;
        }
        long invokeId = packet.getUserdata().getInvokeId();
        executeThrottled(() -> sendAmsTCPPacket(packet, invokeId)).thenAccept(resp -> {
            AmsPacket userdata = resp.getUserdata();
            if (!(userdata instanceof AdsAddDeviceNotificationResponse response)) {
                future.completeExceptionally(new PlcException("Unexpected response type: " + userdata.getClass().getSimpleName()));
                return;
            }
            if (response.getResult() == ReturnCode.OK) {
                PlcSubscriptionTag subscriptionTag = subscriptionRequest.getTag(tagName);
                String dataTypeName = resolvedTags.get((AdsTag) subscriptionTag.getTag()).getPlcDataType();
                Optional<AdsDataTypeTableEntry> dt = getDataTypeTableEntry(dataTypeName);
                if (dt.isEmpty()) {
                    future.completeExceptionally(new PlcRuntimeException("Could not find data type: " + dataTypeName));
                    return;
                }
                responses.put(tagName, new DefaultPlcResponseItem<>(
                    parsePlcResponseCode(response.getResult()),
                    new AdsSubscriptionHandle(this, tagName, dt.get(), response.getNotificationHandle())));
                if (!hasMorePackets) {
                    future.complete(new DefaultPlcSubscriptionResponse(subscriptionRequest, responses));
                    return;
                }
                subscribeRecursively(subscriptionRequest, tagNames, resolvedTags, responses, future, amsTCPPackets);
            } else {
                if (response.getResult() == ReturnCode.ADSERR_DEVICE_INVALIDSIZE) {
                    future.completeExceptionally(new PlcException("The parameter size was not correct (Internal error)"));
                } else {
                    future.completeExceptionally(new PlcException("Unexpected result " + response.getResult()));
                }
            }
        }).exceptionally(t -> {
            future.completeExceptionally(t);
            return null;
        });
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        List<Long> notificationHandles = new ArrayList<>();
        unsubscriptionRequest.getSubscriptionHandles().stream()
            .filter(handle -> handle instanceof AdsSubscriptionHandle)
            .map(handle -> (AdsSubscriptionHandle) handle)
            .forEach(adsHandle -> {
                notificationHandles.add(adsHandle.getNotificationHandle());
                consumers.keySet().stream()
                    .filter(reg -> reg.getSubscriptionHandles().contains(adsHandle))
                    .forEach(DefaultPlcConsumerRegistration::unregister);
            });

        List<AmsTCPPacket> amsTCPPackets = notificationHandles.stream().map(handle -> new AmsTCPPacket(
            new AdsDeleteDeviceNotificationRequest(
                getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
                getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(), ReturnCode.OK, getInvokeId(), handle)))
            .collect(Collectors.toList());

        CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<>();
        Iterator<AmsTCPPacket> it = amsTCPPackets.iterator();
        unsubscribeRecursively(unsubscriptionRequest, future, it);
        return future;
    }

    private void unsubscribeRecursively(PlcUnsubscriptionRequest unsubscriptionRequest,
                                        CompletableFuture<PlcUnsubscriptionResponse> future,
                                        Iterator<AmsTCPPacket> amsTCPPackets) {
        if (!amsTCPPackets.hasNext()) {
            future.complete(new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest));
            return;
        }
        AmsTCPPacket packet = amsTCPPackets.next();
        boolean hasMorePackets = amsTCPPackets.hasNext();
        long invokeId = packet.getUserdata().getInvokeId();
        executeThrottled(() -> sendAmsTCPPacket(packet, invokeId)).thenAccept(resp -> {
            AmsPacket userdata = resp.getUserdata();
            if (!(userdata instanceof AdsDeleteDeviceNotificationResponse response)) {
                future.completeExceptionally(new PlcException("Unexpected response type: " + userdata.getClass().getSimpleName()));
                return;
            }
            if (response.getResult() == ReturnCode.OK) {
                if (!hasMorePackets) {
                    future.complete(new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest));
                    return;
                }
                unsubscribeRecursively(unsubscriptionRequest, future, amsTCPPackets);
            } else if (response.getResult() == ReturnCode.ADSERR_DEVICE_NOTIFYHNDINVALID) {
                future.completeExceptionally(new PlcException("The notification handle is invalid (Internal error)"));
            } else {
                future.completeExceptionally(new PlcException("Unexpected result " + response.getResult()));
            }
        }).exceptionally(t -> {
            future.completeExceptionally(t);
            return null;
        });
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        DefaultPlcConsumerRegistration registration = new DefaultPlcConsumerRegistration(
            this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(registration, consumer);
        return registration;
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        if (registration instanceof DefaultPlcConsumerRegistration r) {
            consumers.remove(r);
        }
    }

    private void dispatchDeviceNotification(AdsDeviceNotificationRequest notification) throws BufferException {
        long receiveTs = System.currentTimeMillis();
        for (AdsStampHeader stamp : notification.getAdsStampHeaders()) {
            // Convert Windows FILETIME to unix epoch ms.
            long unixEpochTimestamp = stamp.getTimestamp().divide(java.math.BigInteger.valueOf(10000L)).longValue() - 11644473600000L;
            for (AdsNotificationSample sample : stamp.getAdsNotificationSamples()) {
                long handle = sample.getNotificationHandle();
                for (DefaultPlcConsumerRegistration registration : consumers.keySet()) {
                    for (PlcSubscriptionHandle subscriptionHandle : registration.getSubscriptionHandles()) {
                        if (subscriptionHandle instanceof AdsSubscriptionHandle adsHandle
                            && adsHandle.getNotificationHandle() == handle) {
                            Map<String, PlcResponseItem<PlcValue>> values = convertSampleToPlc4XResult(adsHandle, sample.getData());
                            DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.ofEpochMilli(unixEpochTimestamp), values);
                            consumers.get(registration).accept(event);
                        }
                    }
                }
            }
        }
    }

    private Map<String, PlcResponseItem<PlcValue>> convertSampleToPlc4XResult(AdsSubscriptionHandle subscriptionHandle, byte[] data) throws BufferException {
        Map<String, PlcResponseItem<PlcValue>> values = new HashMap<>();
        ReadBufferByteBased rb = new ReadBufferByteBased(data, WithOption.WithUnsignedIntegerEncoding("unsigned-binary"), WithOption.WithSignedIntegerEncoding("twos-complement"), WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
        values.put(subscriptionHandle.getTagName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK,
            DataItem.staticParse(rb, getPlcValueTypeForAdsDataType(subscriptionHandle.getAdsDataType()), data.length)));
        return values;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Browse
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        return onBrowseWithInterceptor(browseRequest, (queryName, query, item) -> true);
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        Map<String, List<PlcBrowseItem>> values = new HashMap<>();
        for (String queryName : browseRequest.getQueryNames()) {
            PlcQuery query = browseRequest.getQuery(queryName);
            List<PlcBrowseItem> resultsForQuery = new ArrayList<>();
            for (AdsSymbolTableEntry symbol : symbolTable.values()) {
                Optional<AdsDataTypeTableEntry> dataTypeOpt = getDataTypeTableEntry(symbol.getDataTypeName());
                if (dataTypeOpt.isEmpty()) {
                    LOGGER.warn("couldn't find datatype: {}", symbol.getDataTypeName());
                    continue;
                }
                AdsDataTypeTableEntry dataType = dataTypeOpt.get();
                PlcValueType plcValueType = getPlcValueTypeForAdsDataTypeForBrowse(dataType);

                List<PlcBrowseItem> children = getBrowseItems(symbol.getName(), symbol.getGroup(), symbol.getOffset(), !symbol.getFlagReadOnly(), dataType);
                Map<String, PlcBrowseItem> childMap = new HashMap<>();
                for (PlcBrowseItem child : children) {
                    childMap.put(child.getName(), child);
                }

                Map<String, PlcValue> options = new HashMap<>();
                options.put("comment", new PlcSTRING(symbol.getComment()));
                options.put("group-id", new PlcUDINT(symbol.getGroup()));
                options.put("offset", new PlcUDINT(symbol.getOffset()));
                options.put("size-in-bytes", new PlcUDINT(symbol.getSize()));

                List<ArrayInfo> arrayInfo = new ArrayList<>(dataType.getArrayInfo().size());
                List<ArrayInfo> itemArrayInfo = new ArrayList<>(dataType.getArrayInfo().size());
                for (AdsDataTypeArrayInfo a : dataType.getArrayInfo()) {
                    arrayInfo.add(new DefaultArrayInfo((int) a.getLowerBound(), (int) a.getUpperBound()));
                    itemArrayInfo.add(new DefaultArrayInfo((int) a.getLowerBound(), (int) a.getUpperBound()));
                }
                DefaultPlcBrowseItem item = new DefaultPlcBrowseItem(
                    new SymbolicAdsTag(symbol.getName(), plcValueType, arrayInfo), symbol.getName(),
                    true, !symbol.getFlagReadOnly(),
                    EnumSet.of(PlcSubscriptionType.CYCLIC, PlcSubscriptionType.CHANGE_OF_STATE),
                    false, itemArrayInfo, childMap, options);

                if (interceptor.intercept(queryName, query, item)) {
                    resultsForQuery.add(item);
                }
            }
            responseCodes.put(queryName, PlcResponseCode.OK);
            values.put(queryName, resultsForQuery);
        }
        return CompletableFuture.completedFuture(new DefaultPlcBrowseResponse(browseRequest, responseCodes, values));
    }

    private List<PlcBrowseItem> getBrowseItems(String basePath, long baseGroupId, long baseOffset, boolean parentWritable, AdsDataTypeTableEntry dataType) {
        if (dataType.getArrayDimensions() > 0) {
            Optional<AdsDataTypeTableEntry> opt = getDataTypeTableEntry(dataType.getMainName());
            if (opt.isEmpty()) {
                LOGGER.warn("couldn't find datatype: {}", dataType.getMainName());
                return Collections.emptyList();
            }
            dataType = opt.get();
        }
        if (dataType.getNumChildren() == 0) {
            return Collections.emptyList();
        }
        List<PlcBrowseItem> values = new ArrayList<>(dataType.getNumChildren());
        for (AdsDataTypeTableEntry child : dataType.getChildren()) {
            Optional<AdsDataTypeTableEntry> opt = getDataTypeTableEntry(child.getSecondaryName());
            if (opt.isEmpty()) {
                LOGGER.warn("couldn't find datatype: {} for child {}", dataType.getSecondaryName(), child.getMainName());
                continue;
            }
            AdsDataTypeTableEntry childDataType = opt.get();
            String itemAddress = basePath + "." + child.getMainName();
            PlcValueType plc4xPlcValueType = PlcValueType.valueOf(getPlcValueTypeForAdsDataTypeForBrowse(childDataType).toString());

            List<PlcBrowseItem> children = getBrowseItems(itemAddress, baseGroupId, baseOffset + child.getOffset(), parentWritable, childDataType);
            Map<String, PlcBrowseItem> childMap = new HashMap<>();
            for (PlcBrowseItem ch : children) {
                childMap.put(ch.getName(), ch);
            }

            Map<String, PlcValue> options = new HashMap<>();
            options.put("comment", new PlcSTRING(child.getComment()));
            options.put("group-id", new PlcUDINT(baseGroupId));
            options.put("offset", new PlcUDINT(baseOffset + child.getOffset()));
            options.put("size-in-bytes", new PlcUDINT(childDataType.getSize()));

            List<ArrayInfo> arrayInfo = new ArrayList<>(childDataType.getArrayInfo().size());
            List<ArrayInfo> itemArrayInfo = new ArrayList<>(childDataType.getArrayInfo().size());
            for (AdsDataTypeArrayInfo a : childDataType.getArrayInfo()) {
                arrayInfo.add(new DefaultArrayInfo((int) a.getLowerBound(), (int) a.getUpperBound()));
                itemArrayInfo.add(new DefaultArrayInfo((int) a.getLowerBound(), (int) a.getUpperBound()));
            }
            values.add(new DefaultPlcBrowseItem(
                new SymbolicAdsTag(basePath + "." + child.getMainName(), plc4xPlcValueType, arrayInfo),
                child.getMainName(),
                true, parentWritable,
                EnumSet.of(PlcSubscriptionType.CYCLIC, PlcSubscriptionType.CHANGE_OF_STATE),
                false, itemArrayInfo, childMap, options));
        }
        return values;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Symbolic resolution
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private CompletableFuture<Map<AdsTag, DirectAdsTag>> getDirectAddresses(List<PlcTag> tags) {
        CompletableFuture<Map<AdsTag, DirectAdsTag>> future = new CompletableFuture<>();

        List<SymbolicAdsTag> referencedSymbolicTags = tags.stream()
            .filter(SymbolicAdsTag.class::isInstance)
            .map(SymbolicAdsTag.class::cast)
            .collect(Collectors.toList());

        List<SymbolicAdsTag> symbolicTagsNeedingResolution = referencedSymbolicTags.stream()
            .filter(t -> getDirectAdsTagForSymbolicName(t) == null)
            .collect(Collectors.toList());

        if (!symbolicTagsNeedingResolution.isEmpty()) {
            List<SymbolicAdsTag> requiredResolutionTags = symbolicTagsNeedingResolution.stream()
                .filter(t -> !pendingResolutionRequests.containsKey(t))
                .collect(Collectors.toList());
            if (!requiredResolutionTags.isEmpty()) {
                CompletableFuture<Void> resolutionFuture;
                if (requiredResolutionTags.size() == 1) {
                    SymbolicAdsTag t = requiredResolutionTags.get(0);
                    resolutionFuture = resolveSingleSymbolicAddress(t);
                    pendingResolutionRequests.put(t, resolutionFuture);
                } else {
                    resolutionFuture = resolveMultipleSymbolicAddresses(requiredResolutionTags);
                    for (SymbolicAdsTag t : requiredResolutionTags) {
                        pendingResolutionRequests.put(t, resolutionFuture);
                    }
                }
            }

            CompletableFuture<Void> resolutionComplete = CompletableFuture.allOf(symbolicTagsNeedingResolution.stream()
                .map(pendingResolutionRequests::get)
                .toArray(CompletableFuture[]::new));

            resolutionComplete.handleAsync((unused, throwable) -> {
                Map<AdsTag, DirectAdsTag> mapping = new HashMap<>(tags.size());
                for (PlcTag tag : tags) {
                    if (tag instanceof SymbolicAdsTag) {
                        if (pendingResolutionRequests.containsKey(tag) && pendingResolutionRequests.get(tag).isCompletedExceptionally()) {
                            mapping.put((AdsTag) tag, null);
                        } else {
                            mapping.put((AdsTag) tag, getDirectAdsTagForSymbolicName(tag));
                        }
                    } else {
                        mapping.put((AdsTag) tag, (DirectAdsTag) tag);
                    }
                }
                return future.complete(mapping);
            });
        } else {
            Map<AdsTag, DirectAdsTag> mapping = new HashMap<>(tags.size());
            for (PlcTag tag : tags) {
                if (tag instanceof SymbolicAdsTag) {
                    mapping.put((AdsTag) tag, getDirectAdsTagForSymbolicName(tag));
                } else {
                    mapping.put((AdsTag) tag, (DirectAdsTag) tag);
                }
            }
            future.complete(mapping);
        }
        return future;
    }

    private CompletableFuture<Void> resolveSingleSymbolicAddress(SymbolicAdsTag symbolicAdsTag) {
        AmsPacket request = new AdsReadWriteRequest(
            getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
            getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(), ReturnCode.OK, getInvokeId(),
            ReservedIndexGroups.ADSIGRP_SYM_HNDBYNAME.getValue(), 0L, 4L, null,
            getNullByteTerminatedArray(symbolicAdsTag.getSymbolicAddress()));

        return sendAmsRequest(request, AdsReadWriteResponse.class).thenCompose(response -> {
            if (response.getResult() != ReturnCode.OK) {
                return CompletableFuture.failedFuture(new PlcException(
                    "Couldn't retrieve handle for symbolic tag " + symbolicAdsTag.getSymbolicAddress()
                        + " got return code " + response.getResult().name()));
            }
            try {
                ReadBuffer rb = new ReadBufferByteBased(response.getData(), WithOption.WithUnsignedIntegerEncoding("unsigned-binary"), WithOption.WithSignedIntegerEncoding("twos-complement"), WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
                long handle = rb.readUnsignedLong(32);
                // TODO: Find out how to read the datatype for the given symbolic tag
                //       (the original implementation doesn't actually populate the symbolic→direct
                //        mapping after retrieving the handle, so subscribe-on-symbolic isn't fully
                //        functional — preserved as-is.)
                return CompletableFuture.completedFuture(null);
            } catch (BufferException e) {
                return CompletableFuture.failedFuture(e);
            }
        });
    }

    private CompletableFuture<Void> resolveMultipleSymbolicAddresses(List<SymbolicAdsTag> symbolicAdsTags) {
        long expectedResponseDataSize = (long) symbolicAdsTags.size() * 12;
        byte[] addressData = symbolicAdsTags.stream()
            .map(SymbolicAdsTag::getSymbolicAddress).collect(Collectors.joining("")).getBytes();

        AmsPacket request = new AdsReadWriteRequest(
            getConfiguration().getTargetAmsNetId(), getConfiguration().getTargetAmsPort(),
            getConfiguration().getSourceAmsNetId(), getConfiguration().getSourceAmsPort(), ReturnCode.OK, getInvokeId(),
            ReservedIndexGroups.ADSIGRP_MULTIPLE_READ_WRITE.getValue(),
            (long) symbolicAdsTags.size(), expectedResponseDataSize,
            symbolicAdsTags.stream().map(t -> new AdsMultiRequestItemReadWrite(
                ReservedIndexGroups.ADSIGRP_SYM_HNDBYNAME.getValue(), 0L, 4L, (long) t.getSymbolicAddress().length()))
                .collect(Collectors.toList()),
            addressData);

        return sendAmsRequest(request, AdsReadWriteResponse.class).thenAccept(response -> {
            ReadBuffer rb = new ReadBufferByteBased(response.getData(), WithOption.WithUnsignedIntegerEncoding("unsigned-binary"), WithOption.WithSignedIntegerEncoding("twos-complement"), WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
            Map<SymbolicAdsTag, Long> returnCodes = new HashMap<>();
            symbolicAdsTags.forEach(t -> {
                try {
                    long returnCode = rb.readUnsignedLong(32);
                    long itemLength = rb.readUnsignedLong(32);
                    assert itemLength == 4;
                    returnCodes.put(t, returnCode);
                } catch (BufferException e) {
                    throw new PlcRuntimeException(e);
                }
            });
            symbolicAdsTags.forEach(t -> {
                try {
                    ReturnCode returnCode = ReturnCode.enumForValue(returnCodes.get(t));
                    if (returnCode == ReturnCode.OK) {
                        long handle = rb.readUnsignedLong(32);
                        // TODO: Finish parsing of the response and possibly the reading of
                        //       datatype information for the current tag (preserved as-is).
                    } else if (returnCode == ReturnCode.ADSERR_DEVICE_SYMBOLNOTFOUND) {
                        pendingResolutionRequests.put(t, CompletableFuture.failedFuture(
                            new PlcInvalidTagException("Could not resolve tag " + t.getSymbolicAddress())));
                    }
                } catch (BufferException e) {
                    throw new PlcRuntimeException(e);
                }
            });
        });
    }

    private DirectAdsTag getDirectAdsTagForSymbolicName(PlcTag tag) {
        if (tag instanceof DirectAdsTag d) {
            return d;
        }
        SymbolicAdsTag symbolicAdsTag = (SymbolicAdsTag) tag;
        String symbolicAddress = symbolicAdsTag.getSymbolicAddress();
        String[] addressParts = symbolicAddress.split("\\.");

        if (addressParts.length < 2) {
            if (!symbolTable.containsKey(symbolicAddress)) return null;
            AdsSymbolTableEntry entry = symbolTable.get(symbolicAddress);
            Optional<AdsDataTypeTableEntry> opt = getDataTypeTableEntry(entry.getDataTypeName());
            if (opt.isEmpty()) return null;
            AdsDataTypeTableEntry dataType = opt.get();
            return new DirectAdsTag(entry.getGroup(), entry.getOffset(),
                dataType.getMainName(), dataType.getArrayDimensions());
        } else {
            String symbolName = addressParts[0] + "." + addressParts[1];
            if (!symbolTable.containsKey(symbolName)) return null;
            AdsSymbolTableEntry entry = symbolTable.get(symbolName);
            Optional<AdsDataTypeTableEntry> opt = getDataTypeTableEntry(entry.getDataTypeName());
            if (opt.isEmpty()) return null;
            AdsDataTypeTableEntry dataType = opt.get();
            return resolveDirectAdsTagForSymbolicNameFromDataType(
                Arrays.asList(addressParts).subList(2, addressParts.length),
                entry.getGroup(), entry.getOffset(), dataType);
        }
    }

    private DirectAdsTag resolveDirectAdsTagForSymbolicNameFromDataType(List<String> remainingAddressParts, long currentGroup, long currentOffset, AdsDataTypeTableEntry adsDataTypeTableEntry) {
        if (remainingAddressParts.isEmpty()) {
            // TODO: Implement Array support
            if (adsDataTypeTableEntry.getDataType().getValue() == AdsDataType.CHAR.getValue()) {
                int stringLength = (int) adsDataTypeTableEntry.getSize() - 1;
                return new DirectAdsStringTag(currentGroup, currentOffset, adsDataTypeTableEntry.getMainName(), stringLength, 1);
            } else if (adsDataTypeTableEntry.getDataType().getValue() == AdsDataType.WCHAR.getValue()) {
                int stringLength = (int) (adsDataTypeTableEntry.getSize() - 2) / 2;
                return new DirectAdsStringTag(currentGroup, currentOffset, adsDataTypeTableEntry.getMainName(), stringLength, 1);
            } else {
                return new DirectAdsTag(currentGroup, currentOffset, adsDataTypeTableEntry.getMainName(), 1);
            }
        }
        for (AdsDataTypeTableEntry child : adsDataTypeTableEntry.getChildren()) {
            if (child.getMainName().equals(remainingAddressParts.get(0))) {
                Optional<AdsDataTypeTableEntry> opt = getDataTypeTableEntry(child.getSecondaryName());
                if (opt.isEmpty()) {
                    throw new PlcRuntimeException("Could not resolve data type " + child.getSecondaryName());
                }
                return resolveDirectAdsTagForSymbolicNameFromDataType(
                    remainingAddressParts.subList(1, remainingAddressParts.size()),
                    currentGroup, currentOffset + child.getOffset(), opt.get());
            }
        }
        throw new PlcRuntimeException(String.format("Couldn't find child with name '%s' for type '%s'",
            remainingAddressParts.get(0), adsDataTypeTableEntry.getMainName()));
    }

    private PlcValueType getPlcValueTypeForAdsDataTypeForBrowse(AdsDataTypeTableEntry dataTypeTableEntry) {
        String dataTypeName = (!dataTypeTableEntry.getSecondaryName().isEmpty() && !dataTypeTableEntry.getMainName().equals("BOOL")) ?
            dataTypeTableEntry.getSecondaryName() : dataTypeTableEntry.getMainName();
        if (dataTypeName.startsWith("STRING(")) dataTypeName = "STRING";
        else if (dataTypeName.startsWith("WSTRING(")) dataTypeName = "WSTRING";
        try {
            return PlcValueType.valueOf(dataTypeName);
        } catch (IllegalArgumentException e) {
            return PlcValueType.Struct;
        }
    }

    private PlcValueType getPlcValueTypeForAdsDataType(AdsDataTypeTableEntry dataTypeTableEntry) {
        String dataTypeName = dataTypeTableEntry.getMainName();
        if (dataTypeName.startsWith("STRING(")) dataTypeName = "STRING";
        else if (dataTypeName.startsWith("WSTRING(")) dataTypeName = "WSTRING";
        try {
            return PlcValueType.valueOf(dataTypeName);
        } catch (IllegalArgumentException e) {
            if (dataTypeTableEntry.getArrayDimensions() > 0) return PlcValueType.List;
            if (dataTypeTableEntry.getChildren().isEmpty()) {
                try {
                    dataTypeName = dataTypeTableEntry.getSecondaryName();
                    if (dataTypeName.startsWith("STRING(")) dataTypeName = "STRING";
                    else if (dataTypeName.startsWith("WSTRING(")) dataTypeName = "WSTRING";
                    return PlcValueType.valueOf(dataTypeName);
                } catch (IllegalArgumentException e2) {
                    return PlcValueType.NULL;
                }
            }
            return PlcValueType.Struct;
        }
    }

    private Optional<AdsDataTypeTableEntry> getDataTypeTableEntry(String name) {
        if (dataTypeTable.containsKey(name)) {
            return Optional.of(dataTypeTable.get(name));
        }
        try {
            AdsDataType adsDataType;
            int numBytes;
            if (name.startsWith("STRING(")) {
                adsDataType = AdsDataType.valueOf("CHAR");
                numBytes = Integer.parseInt(name.substring(7, name.length() - 1)) + 1;
            } else if (name.startsWith("WSTRING(")) {
                adsDataType = AdsDataType.valueOf("WCHAR");
                numBytes = Integer.parseInt(name.substring(8, name.length() - 1)) * 2 + 2;
            } else {
                adsDataType = AdsDataType.valueOf(name);
                numBytes = adsDataType.getNumBytes();
            }
            return Optional.of(new AdsDataTypeTableEntry(
                128L, 1L, 0L, 0L, (long) numBytes, 0L,
                AdsDatatypeId.enumForValue(adsDataType.getValue()),
                false, false, false, false, false, false, false, false, false, false, false, false,
                false, false, false, false, false, false, false, false, false, false, false,
                0, 0,
                name, "", "",
                Collections.emptyList(), Collections.emptyList(), new byte[0],
                null, null, null, new byte[0]));
        } catch (IllegalArgumentException e) {
            return Optional.empty();
        }
    }

    private byte[] getNullByteTerminatedArray(String value) {
        byte[] valueBytes = value.getBytes();
        byte[] nullTerminatedBytes = new byte[valueBytes.length + 1];
        System.arraycopy(valueBytes, 0, nullTerminatedBytes, 0, valueBytes.length);
        return nullTerminatedBytes;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // AMS route setup (preserved from original; not currently invoked because authentication
    // is not surfaced through ConnectionBase). Kept here for reference / future re-enablement.
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @SuppressWarnings("unused")
    private CompletableFuture<Void> setupAmsRoute(PlcUsernamePasswordAuthentication authentication, InetAddress remoteAddress) {
        CompletableFuture<Void> future = new CompletableFuture<>();
        new Thread(() -> {
            try {
                // The local address determination is no longer trivially available here; if the
                // transport layer can supply a SocketAddress in future versions, plug it in here.
                AmsNetId sourceAmsNetId = new AmsNetId(
                    getConfiguration().getSourceAmsNetId().getOctet1(), getConfiguration().getSourceAmsNetId().getOctet2(),
                    getConfiguration().getSourceAmsNetId().getOctet3(), getConfiguration().getSourceAmsNetId().getOctet4(),
                    getConfiguration().getSourceAmsNetId().getOctet5(), getConfiguration().getSourceAmsNetId().getOctet6());
                String routeName = String.format("PLC4X-%d.%d.%d.%d.%d.%d",
                    sourceAmsNetId.getOctet1(), sourceAmsNetId.getOctet2(), sourceAmsNetId.getOctet3(),
                    sourceAmsNetId.getOctet4(), sourceAmsNetId.getOctet5(), sourceAmsNetId.getOctet6());

                AdsDiscovery req = new AdsDiscovery(getInvokeId(), Operation.ADD_OR_UPDATE_ROUTE_REQUEST,
                    sourceAmsNetId, AdsPortNumbers.SYSTEM_SERVICE,
                    Arrays.asList(new AdsDiscoveryBlockRouteName(new AmsString(routeName)),
                        new AdsDiscoveryBlockAmsNetId(sourceAmsNetId),
                        new AdsDiscoveryBlockUserName(new AmsString(authentication.getUsername())),
                        new AdsDiscoveryBlockPassword(new AmsString(authentication.getPassword())),
                        new AdsDiscoveryBlockHostName(new AmsString(remoteAddress.getHostAddress()))));

                try (DatagramSocket socket = new DatagramSocket(Constants.ADSDISCOVERYUDPDEFAULTPORT)) {
                    WriteBufferByteBased wb = new WriteBufferByteBased(new byte[req.getLengthInBytes()],
                        WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
                    req.serialize(wb);

                    DatagramPacket pkt = new DatagramPacket(wb.getBytes(), wb.getBytes().length,
                        remoteAddress, Constants.ADSDISCOVERYUDPDEFAULTPORT);
                    socket.send(pkt);

                    byte[] buf = new byte[100];
                    DatagramPacket responsePacket = new DatagramPacket(buf, buf.length);
                    socket.setSoTimeout(getConfiguration().getTimeoutRequest());
                    socket.receive(responsePacket);

                    ReadBuffer rb = new ReadBufferByteBased(responsePacket.getData(),
                        WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
                    AdsDiscovery resp = AdsDiscovery.staticParse(rb);
                    if (resp.getRequestId() == 1) {
                        for (AdsDiscoveryBlock block : resp.getBlocks()) {
                            if (block.getBlockType() == AdsDiscoveryBlockType.STATUS) {
                                AdsDiscoveryBlockStatus status = (AdsDiscoveryBlockStatus) block;
                                if (status.getStatus() != Status.SUCCESS) {
                                    future.completeExceptionally(new PlcException("Error adding AMS route"));
                                    return;
                                }
                            }
                        }
                    }
                    future.complete(null);
                }
            } catch (Exception e) {
                future.completeExceptionally(new PlcException("Error adding AMS route", e));
            }
        }).start();
        return future;
    }

}
