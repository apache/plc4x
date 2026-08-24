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
package org.apache.plc4x.java.umas;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseRequestInterceptor;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.messages.PlcPingRequest;
import org.apache.plc4x.java.api.messages.PlcPingResponse;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
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
import org.apache.plc4x.java.utils.subscriptionemulation.PollingSubscriptionConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcPingResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcDATE;
import org.apache.plc4x.java.spi.values.PlcDATE_AND_TIME;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.values.PlcTIME;
import org.apache.plc4x.java.spi.values.PlcTIME_OF_DAY;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.umas.configuration.UmasConfiguration;
import org.apache.plc4x.java.umas.readwrite.DataItem;
import org.apache.plc4x.java.umas.readwrite.ModbusPDU;
import org.apache.plc4x.java.umas.readwrite.ModbusPDUError;
import org.apache.plc4x.java.umas.readwrite.ModbusTcpADU;
import org.apache.plc4x.java.umas.readwrite.UmasArrayDimension;
import org.apache.plc4x.java.umas.readwrite.UmasArrayTypeDefinition;
import org.apache.plc4x.java.umas.readwrite.UmasDataType;
import org.apache.plc4x.java.umas.readwrite.UmasDatatypeReference;
import org.apache.plc4x.java.umas.readwrite.UmasInitCommsRequest;
import org.apache.plc4x.java.umas.readwrite.UmasInitCommsResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDU;
import org.apache.plc4x.java.umas.readwrite.UmasPDUErrorResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUItem;
import org.apache.plc4x.java.umas.readwrite.UmasPDUPlcIdentRequest;
import org.apache.plc4x.java.umas.readwrite.UmasPDUPlcIdentResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUPlcStatusRequest;
import org.apache.plc4x.java.umas.readwrite.UmasPDUProjectInfoRequest;
import org.apache.plc4x.java.umas.readwrite.UmasPDUProjectInfoResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadDatatypeNamesResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadMemoryBlockRequest;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadMemoryBlockResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadUmasUDTDefinitionResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadUnlocatedVariableNamesRequest;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadUnlocatedVariableNamesResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadUnlocatedVariableResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadVariableRequest;
import org.apache.plc4x.java.umas.readwrite.UmasPDUReadVariableResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDURepeatRequest;
import org.apache.plc4x.java.umas.readwrite.UmasPDURepeatResponse;
import org.apache.plc4x.java.umas.readwrite.UmasPDUWriteVariableRequest;
import org.apache.plc4x.java.umas.readwrite.UmasPDUWriteVariableResponse;
import org.apache.plc4x.java.umas.readwrite.UmasUDTDefinition;
import org.apache.plc4x.java.umas.readwrite.UmasUnlocatedVariableReference;
import org.apache.plc4x.java.umas.readwrite.VariableReadRequestReference;
import org.apache.plc4x.java.umas.readwrite.VariableWriteRequestReference;
import org.apache.plc4x.java.umas.tag.SymbolicUmasTag;
import org.apache.plc4x.java.umas.tag.UmasTagHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Comparator;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * UMAS connection. Direct port of the legacy {@code UmasProtocolLogic} +
 * {@code UmasDriverContext} to the SPI3 {@link ConnectionBase} model.
 *
 * <p>UMAS is tunneled inside Modbus/TCP (function code 0x5A). On connect this
 * class runs a multi-step async handshake against the PLC (PlcIdent →
 * InitComms → Repeat → memory-block reads → project-info queries) and then
 * eagerly downloads the data dictionary (symbol + type tables) so that
 * read/write/browse have everything they need immediately.</p>
 *
 * <p>Requests are matched to responses by Modbus transaction ID. Each outgoing
 * UMAS PDU is recorded in {@link UmasFunctionKeyTracker} so the codec can
 * resolve the response subtype (UMAS responses share function key 0xFE).</p>
 */
public class UmasConnection extends PollingSubscriptionConnectionBase<UmasConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(UmasConnection.class);

    private static final int CUSTOM_TYPE_THRESHOLD = 0x1A;
    private static final int RECORD_TYPE_DD02 = 0xDD02;
    private static final int RECORD_TYPE_DD03 = 0xDD03;
    private static final int SYMBOL_TABLE_BLOCK = 0xFFFF;
    private static final int DEFAULT_STRING_BUFFER_SIZE = 254;

    /** Options for raw-bytes parsing of UMAS payload chunks (little-endian, all numerics binary). */
    private static final WithOption[] LITTLE_ENDIAN_OPTIONS = {
        WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    private UmasMessageCodec messageCodec;

    /**
     * This connection's request function keys. One per connection, because a transaction id is a
     * counter belonging to a connection and means nothing outside it.
     */
    private final UmasFunctionKeyTracker functionKeyTracker = new UmasFunctionKeyTracker();
    private final short unitIdentifier;
    private final long requestTimeoutMs;

    // -------- Session state (merged from old UmasDriverContext) --------

    private final AtomicInteger transactionIdGenerator = new AtomicInteger(1);
    private volatile boolean handshakeComplete = false;

    private volatile String plcHostname;
    private volatile int plcModel;
    private volatile int plcFirmwareVersion;
    private volatile int maxFrameSize;
    private volatile short pairingKey;
    private volatile long hardwareId;
    private volatile long projectCrc;

    private final Map<Integer, List<UmasUDTDefinition>> customTypeFields = new ConcurrentHashMap<>();
    private final Map<Integer, String> customTypeNames = new ConcurrentHashMap<>();
    private final Map<Integer, Integer> customTypeElementTypeIds = new ConcurrentHashMap<>();
    private final Map<Integer, List<UmasArrayDimension>> customTypeDimensions = new ConcurrentHashMap<>();
    private final Map<String, UmasUnlocatedVariableReference> symbolTable = new ConcurrentHashMap<>();
    private final Map<Integer, UmasDataType> dataTypeTable = new ConcurrentHashMap<>();
    private final Map<Integer, Integer> dataTypeSizes = new ConcurrentHashMap<>();
    private final Map<String, Integer> symbolSizes = new ConcurrentHashMap<>();

    /** Pending responses keyed by Modbus transaction ID. */
    private final Map<Integer, CompletableFuture<ModbusTcpADU>> pendingResponses = new ConcurrentHashMap<>();

    public UmasConnection(UmasConfiguration configuration,
                          TransportInstance<?> transportInstance,
                          AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
        this.unitIdentifier = (short) configuration.getUnitIdentifier();
        this.requestTimeoutMs = configuration.getRequestTimeout();
        this.maxFrameSize = configuration.getMaxFrameSize();
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new UmasTagHandler();
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
        return handshakeComplete && messageCodec != null && messageCodec.isOpen();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Connect / close
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new UmasMessageCodec(transportInstance, this::handleIncomingMessage, functionKeyTracker);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming UMAS data", e);
            }
        });

        try {
            doHandshake().get(requestTimeoutMs * 10L, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("Error during UMAS connect handshake", e);
        }

        handshakeComplete = true;
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT,
                String.format("UMAS connection established to '%s' (model=%d, firmware=%d)",
                    plcHostname, plcModel, plcFirmwareVersion));
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    /**
     * Multi-step handshake. Each call returns a {@link CompletableFuture} that
     * completes when the matching response indication arrives. We chain them
     * together with {@code thenCompose} so a failure aborts the rest.
     */
    private CompletableFuture<Void> doHandshake() {
        return performPlcIdent()
            .thenCompose(v -> performInitComms())
            .thenCompose(v -> performRepeat())
            .thenCompose(v -> performReadMemoryBlock(0x30, 0, 33))
            .thenCompose(v -> performProjectInfo((short) 1))
            .thenCompose(v -> performReadMemoryBlock(0x13, 0, 33))
            .thenCompose(v -> performProjectInfo((short) 0))
            .thenCompose(v -> performProjectInfo((short) 4))
            .thenCompose(v -> performProjectInfo((short) 1))
            .thenCompose(v -> performProjectInfo((short) 3))
            // Eagerly download data dictionary so reads/writes/browse work
            // immediately. The dictionary downloads use blocking get() calls
            // so we offload them to a separate thread to avoid deadlocking
            // the receive loop.
            .thenCompose(v -> CompletableFuture.runAsync(() -> {
                try {
                    loadDataDictionary();
                } catch (Exception e) {
                    LOGGER.warn("Failed to eagerly load data dictionary during connect: {}", e.getMessage());
                }
            }))
            .thenAccept(v -> LOGGER.info(
                "UMAS connection established: hostname={}, model={}, firmware={}",
                plcHostname, plcModel, plcFirmwareVersion));
    }

    private CompletableFuture<Void> performPlcIdent() {
        return sendRequest(new UmasPDUPlcIdentRequest((short) 0), "PlcIdent")
            .thenAccept(response -> {
                if (response instanceof UmasPDUPlcIdentResponse ident) {
                    plcHostname = ident.getHostname();
                    plcModel = ident.getModel();
                    plcFirmwareVersion = ident.getComVersion();
                    LOGGER.info("PlcIdent: hostname={}, model={}, comVersion={}",
                        ident.getHostname(), ident.getModel(), ident.getComVersion());
                } else {
                    throw new PlcRuntimeException("PlcIdent: unexpected response type " + response.getClass().getSimpleName());
                }
            });
    }

    private CompletableFuture<Void> performInitComms() {
        return sendRequest(new UmasInitCommsRequest((short) 0, (short) 0x00), "InitComms")
            .thenAccept(response -> {
                if (response instanceof UmasInitCommsResponse init) {
                    maxFrameSize = init.getMaxFrameSize();
                    plcFirmwareVersion = init.getFirmwareVersion();
                    LOGGER.info("InitComms: maxFrameSize={}, firmwareVersion={}",
                        init.getMaxFrameSize(), init.getFirmwareVersion());
                } else {
                    throw new PlcRuntimeException("InitComms: unexpected response type " + response.getClass().getSimpleName());
                }
            });
    }

    private CompletableFuture<Void> performRepeat() {
        int echoSize = maxFrameSize - 3;
        byte[] echoData = new byte[echoSize];
        Arrays.fill(echoData, 1, echoSize, (byte) 0x54);
        return sendRequest(new UmasPDURepeatRequest((short) 0, echoData), "Repeat")
            .thenAccept(response -> {
                if (response instanceof UmasPDURepeatResponse repeat) {
                    LOGGER.info("Repeat: echo OK, {} bytes returned",
                        repeat.getBlock() != null ? repeat.getBlock().length : 0);
                } else {
                    throw new PlcRuntimeException("Repeat: unexpected response type " + response.getClass().getSimpleName());
                }
            });
    }

    private CompletableFuture<Void> performReadMemoryBlock(int blockNumber, int offset, int numberOfBytes) {
        String stepName = "ReadMemoryBlock(0x" + String.format("%02X", blockNumber) + ")";
        UmasPDUItem request = new UmasPDUReadMemoryBlockRequest(
            (short) 0, (short) 0x01, blockNumber, offset, 0, numberOfBytes);
        return sendRequest(request, stepName)
            .thenAccept(response -> {
                if (!(response instanceof UmasPDUReadMemoryBlockResponse memResponse)) {
                    LOGGER.warn("{}: unexpected response type {}", stepName, response.getClass().getSimpleName());
                    return;
                }
                byte[] block = memResponse.getBlock();
                LOGGER.info("{}: {} bytes received", stepName, block != null ? block.length : 0);
                if (blockNumber == 0x30 && block != null && block.length >= 17) {
                    // UmasMemoryBlockBasicInfo: range(2) + notSure(2) + index(1) + hardwareId(4) = 9 bytes
                    long hwId = (block[5] & 0xFFL) | ((block[6] & 0xFFL) << 8)
                        | ((block[7] & 0xFFL) << 16) | ((block[8] & 0xFFL) << 24);
                    hardwareId = hwId;
                    // After the 9-byte header come two 32-bit LE hashes. The
                    // project CRC used in FC 0x22/0x23 read/write requests is
                    // hash1 + hash2 (discovered by comparing working Schneider
                    // OPC UA Server traffic against the raw block 0x30).
                    long hash1 = (block[9] & 0xFFL) | ((block[10] & 0xFFL) << 8)
                        | ((block[11] & 0xFFL) << 16) | ((block[12] & 0xFFL) << 24);
                    long hash2 = (block[13] & 0xFFL) | ((block[14] & 0xFFL) << 8)
                        | ((block[15] & 0xFFL) << 16) | ((block[16] & 0xFFL) << 24);
                    projectCrc = (hash1 + hash2) & 0xFFFFFFFFL;
                    LOGGER.info("{}: hardwareId=0x{}, hash1=0x{}, hash2=0x{}, projectCRC=0x{}",
                        stepName, String.format("%08X", hwId),
                        String.format("%08X", hash1), String.format("%08X", hash2),
                        String.format("%08X", projectCrc));
                }
            });
    }

    private CompletableFuture<Void> performProjectInfo(short subcode) {
        String stepName = "ProjectInfo(subcode=" + subcode + ")";
        return sendRequest(new UmasPDUProjectInfoRequest((short) 0, subcode), stepName)
            .thenAccept(response -> {
                if (response instanceof UmasPDUProjectInfoResponse projInfo) {
                    byte[] block = projInfo.getBlock();
                    LOGGER.info("{}: {} bytes", stepName, block != null ? block.length : 0);
                } else {
                    LOGGER.warn("{}: unexpected response type {}", stepName, response.getClass().getSimpleName());
                }
            });
    }

    @Override
    public void close() throws Exception {
        handshakeComplete = false;
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        PlcRuntimeException closed = new PlcRuntimeException("Connection closed");
        pendingResponses.values().forEach(f -> f.completeExceptionally(closed));
        pendingResponses.clear();
        functionKeyTracker.clear();
        super.close();
        LOGGER.info("UMAS connection closed");
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        handshakeComplete = false;
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");
        PlcRuntimeException ex = new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause);
        pendingResponses.values().forEach(f -> f.completeExceptionally(ex));
        pendingResponses.clear();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming dispatch
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingMessage(ModbusTcpADU response) {
        int transactionId = response.getTransactionIdentifier();
        CompletableFuture<ModbusTcpADU> future = pendingResponses.remove(transactionId);
        if (future != null) {
            future.complete(response);
        } else {
            LOGGER.debug("Received UMAS response for unknown transaction id {}", transactionId);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Core request/response (replaces the old sendAsyncRequest helper)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * Sends a UMAS PDU and returns a future that completes with the parsed
     * response item. Wraps the PDU in a ModbusTcpADU, records the function
     * key in the tracker so the codec can resolve the response subtype, and
     * stores a pending future keyed by transaction ID.
     */
    private CompletableFuture<UmasPDUItem> sendRequest(UmasPDUItem item, String stepName) {
        int transactionId = nextTransactionId();
        functionKeyTracker.trackRequest(transactionId, item.getUmasFunctionKey());
        UmasPDU pdu = new UmasPDU(item);
        ModbusTcpADU adu = new ModbusTcpADU(transactionId, unitIdentifier, pdu);

        CompletableFuture<ModbusTcpADU> future = new CompletableFuture<>();
        pendingResponses.put(transactionId, future);
        try {
            messageCodec.send(adu);
        } catch (MessageCodecException e) {
            pendingResponses.remove(transactionId);
            // Nothing was sent, so nothing will answer: the function key would sit here until a
            // later request reached the same id and found it.
            functionKeyTracker.forget(transactionId);
            return CompletableFuture.failedFuture(new PlcRuntimeException("Failed to send " + stepName, e));
        }

        return future
            .orTimeout(requestTimeoutMs, TimeUnit.MILLISECONDS)
            .whenComplete((r, e) -> {
                if (e != null) {
                    pendingResponses.remove(transactionId);
                    // A request that gave up leaves its function key behind otherwise, and the
                    // transaction id counter comes round again.
                    functionKeyTracker.forget(transactionId);
                }
            })
            .thenApply(response -> extractUmasResponse(response, stepName));
    }

    private int nextTransactionId() {
        int id = transactionIdGenerator.getAndIncrement();
        if (id > 0xFFFF) {
            transactionIdGenerator.compareAndSet(id + 1, 1);
            return id & 0xFFFF;
        }
        return id;
    }

    private static UmasPDUItem extractUmasResponse(ModbusTcpADU response, String stepName) {
        ModbusPDU pdu = response.getPdu();
        if (pdu instanceof ModbusPDUError errorPdu) {
            throw new PlcRuntimeException(stepName + " received Modbus error: " + errorPdu.getExceptionCode());
        }
        if (!(pdu instanceof UmasPDU umasPdu)) {
            throw new PlcRuntimeException(stepName + " received unexpected PDU type: " + pdu.getClass().getSimpleName());
        }
        return umasPdu.getItem();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Ping
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        return sendRequest(new UmasPDUPlcStatusRequest(pairingKey), "Ping")
            .handle((response, error) -> {
                if (error != null) {
                    LOGGER.debug("UMAS ping failed", error);
                    return new DefaultPlcPingResponse(pingRequest, PlcResponseCode.REMOTE_ERROR);
                }
                return new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK);
            });
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseItem<PlcValue>> responseItems = new LinkedHashMap<>();
            for (String tagName : readRequest.getTagNames()) {
                PlcTag tag = readRequest.getTag(tagName);
                responseItems.put(tagName, readSingleTag(tagName, tag));
            }
            return new DefaultPlcReadResponse(readRequest, responseItems);
        });
    }

    private PlcResponseItem<PlcValue> readSingleTag(String tagName, PlcTag tag) {
        if (!(tag instanceof SymbolicUmasTag symbolicTag)) {
            LOGGER.warn("Read tag '{}' is not a SymbolicUmasTag: {}", tagName, tag.getClass().getSimpleName());
            return new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null);
        }

        String symbolicAddress = symbolicTag.getSymbolicAddress().toLowerCase();
        Optional<UmasUnlocatedVariableReference> symbolOpt = Optional.ofNullable(symbolTable.get(symbolicAddress));
        if (symbolOpt.isEmpty()) {
            LOGGER.warn("Read tag '{}': symbol '{}' not found in symbol table", tagName, symbolicAddress);
            return new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null);
        }

        UmasUnlocatedVariableReference symbol = symbolOpt.get();
        try {
            VariableReadRequestReference readRef = buildReadReference(symbol);
            UmasPDUItem request = new UmasPDUReadVariableRequest(
                pairingKey, projectCrc, (short) 1, List.of(readRef));

            UmasPDUItem responseItem = sendRequest(request, "ReadVariable(" + tagName + ")")
                .get(requestTimeoutMs + 1000, TimeUnit.MILLISECONDS);

            if (responseItem instanceof UmasPDUReadVariableResponse readResponse) {
                PlcValue value = parseReadResponse(symbol, readResponse.getBlock());
                return new DefaultPlcResponseItem<>(PlcResponseCode.OK, value);
            }
            if (responseItem instanceof UmasPDUErrorResponse) {
                return new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null);
            }
            return new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null);
        } catch (Exception e) {
            LOGGER.error("Read tag '{}' failed: {}", tagName, e.getMessage());
            return new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null);
        }
    }

    private VariableReadRequestReference buildReadReference(UmasUnlocatedVariableReference symbol) {
        int dataTypeId = symbol.getDataType();
        // The symbol's 32-bit offset encodes two fields:
        //   - lower 8 bits  → offset (uint 8 in VariableReadRequestReference)
        //   - upper bits     → baseOffset (uint 16 in VariableReadRequestReference)
        long symbolOffset = symbol.getOffset();
        int baseOffset = (int) (symbolOffset >> 8);
        short offset = (short) (symbolOffset & 0xFF);

        // STRING: requestSize=17 doesn't fit in the 4-bit dataSizeIndex field.
        // Read as a byte array instead: isArray=1, dataSizeIndex=1, arrayLength=bufferSize.
        if (UmasDataType.isDefined((short) dataTypeId)
                && UmasDataType.enumForValue((short) dataTypeId) == UmasDataType.STRING) {
            int stringSize = Optional.ofNullable(dataTypeSizes.get(dataTypeId))
                .or(() -> Optional.ofNullable(symbolSizes.get(symbol.getValue().toLowerCase())))
                .orElse(DEFAULT_STRING_BUFFER_SIZE);
            return new VariableReadRequestReference(
                (byte) 1, (byte) 1, symbol.getBlock(),
                baseOffset, offset, stringSize);
        }

        byte dataSizeIndex;
        if (UmasDataType.isDefined((short) dataTypeId)) {
            UmasDataType umasType = UmasDataType.enumForValue((short) dataTypeId);
            dataSizeIndex = (byte) umasType.getRequestSize();
        } else {
            dataSizeIndex = (byte) 3;
        }
        return new VariableReadRequestReference(
            (byte) 0, dataSizeIndex, symbol.getBlock(), baseOffset, offset, null);
    }

    private PlcValue parseReadResponse(UmasUnlocatedVariableReference symbol, byte[] block) throws BufferException {
        if (block == null || block.length == 0) {
            throw new PlcRuntimeException("Read response has empty data block");
        }
        int dataTypeId = symbol.getDataType();
        if (!UmasDataType.isDefined((short) dataTypeId)) {
            return new PlcRawByteArray(block);
        }
        UmasDataType umasType = UmasDataType.enumForValue((short) dataTypeId);

        // Types we parse by hand because the generated DataItem code either
        // returns no PlcValue (temporal types) or uses the wrong read strategy
        // for STRING.
        switch (umasType) {
            case STRING:
                for (int i = 0; i < block.length; i++) {
                    if (block[i] == 0x00) {
                        return new PlcSTRING(new String(block, 0, i, StandardCharsets.UTF_8));
                    }
                }
                return new PlcSTRING(new String(block, StandardCharsets.UTF_8));
            case TIME:
                // uint32 milliseconds, little-endian.
                return new PlcTIME(readUint32LE(block));
            case DATE: {
                // BCD: day(1) + month(1) + year(2 LE).
                int day = decodeBcdByte(block[0]);
                int month = decodeBcdByte(block[1]);
                int year = decodeBcd16(block[2], block[3]);
                return new PlcDATE(java.time.LocalDate.of(year, month, day));
            }
            case TOD: {
                // BCD: centiseconds(1) + seconds(1) + minutes(1) + hours(1).
                int secs = decodeBcdByte(block[1]);
                int mins = decodeBcdByte(block[2]);
                int hours = decodeBcdByte(block[3]);
                long totalSeconds = hours * 3600L + mins * 60L + secs;
                return new PlcTIME_OF_DAY(totalSeconds);
            }
            case DATE_AND_TIME: {
                // 8 bytes: reserved(1) + seconds(1 BCD) + minutes(1 BCD) +
                // hour(1 BCD) + day(1 BCD) + month(1 BCD) + year(2 BCD LE).
                int seconds = decodeBcdByte(block[1]);
                int minutes = decodeBcdByte(block[2]);
                int hour = decodeBcdByte(block[3]);
                int dtDay = decodeBcdByte(block[4]);
                int dtMonth = decodeBcdByte(block[5]);
                int dtYear = decodeBcd16(block[6], block[7]);
                return new PlcDATE_AND_TIME(java.time.LocalDateTime.of(
                    dtYear, dtMonth, dtDay, hour, minutes, seconds));
            }
            default: {
                ReadBuffer readBuffer = new ReadBufferByteBased(block, LITTLE_ENDIAN_OPTIONS);
                return DataItem.staticParse(readBuffer, umasType, 1);
            }
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseCode> responseCodes = new LinkedHashMap<>();
            for (String tagName : writeRequest.getTagNames()) {
                PlcTag tag = writeRequest.getTag(tagName);
                PlcValue value = writeRequest.getPlcValue(tagName);
                responseCodes.put(tagName, writeSingleTag(tagName, tag, value));
            }
            return new DefaultPlcWriteResponse(writeRequest, responseCodes);
        });
    }

    private PlcResponseCode writeSingleTag(String tagName, PlcTag tag, PlcValue value) {
        if (!(tag instanceof SymbolicUmasTag symbolicTag)) {
            return PlcResponseCode.INVALID_ADDRESS;
        }
        String symbolicAddress = symbolicTag.getSymbolicAddress().toLowerCase();
        Optional<UmasUnlocatedVariableReference> symbolOpt = Optional.ofNullable(symbolTable.get(symbolicAddress));
        if (symbolOpt.isEmpty()) {
            return PlcResponseCode.NOT_FOUND;
        }

        UmasUnlocatedVariableReference symbol = symbolOpt.get();
        try {
            byte[] serializedData = serializeValue(symbol, value);
            VariableWriteRequestReference writeRef = buildWriteReference(symbol, serializedData);

            UmasPDUItem request = new UmasPDUWriteVariableRequest(
                pairingKey, projectCrc, (short) 1, List.of(writeRef));

            UmasPDUItem responseItem = sendRequest(request, "WriteVariable(" + tagName + ")")
                .get(requestTimeoutMs + 1000, TimeUnit.MILLISECONDS);

            if (responseItem instanceof UmasPDUWriteVariableResponse) {
                return PlcResponseCode.OK;
            }
            if (responseItem instanceof UmasPDUErrorResponse) {
                return PlcResponseCode.REMOTE_ERROR;
            }
            return PlcResponseCode.INTERNAL_ERROR;
        } catch (Exception e) {
            LOGGER.error("Write tag '{}' failed: {}", tagName, e.getMessage());
            return PlcResponseCode.REMOTE_ERROR;
        }
    }

    private VariableWriteRequestReference buildWriteReference(UmasUnlocatedVariableReference symbol, byte[] data) {
        int dataTypeId = symbol.getDataType();
        // Write references use a different offset encoding than read references:
        //   Read:  baseOffset = offset >> 8  (high 16 bits), offset = offset & 0xFF  (low 8 bits)
        //   Write: baseOffset = offset & 0xFFFF (low 16 bits), offset = (offset >> 16) & 0xFFFF (high 16 bits)
        long symbolOffset = symbol.getOffset();
        int baseOffset = (int) (symbolOffset & 0xFFFF);
        int offset = (int) ((symbolOffset >> 16) & 0xFFFF);

        if (UmasDataType.isDefined((short) dataTypeId)
                && UmasDataType.enumForValue((short) dataTypeId) == UmasDataType.STRING) {
            // STRING: write as byte array since requestSize=17 doesn't fit in 4-bit dataSizeIndex.
            return new VariableWriteRequestReference(
                (byte) 1, (byte) 1, symbol.getBlock(),
                baseOffset, offset, data.length, data);
        }
        byte dataSizeIndex;
        if (UmasDataType.isDefined((short) dataTypeId)) {
            UmasDataType umasType = UmasDataType.enumForValue((short) dataTypeId);
            dataSizeIndex = (byte) umasType.getRequestSize();
        } else {
            dataSizeIndex = (byte) 3;
        }
        return new VariableWriteRequestReference(
            (byte) 0, dataSizeIndex, symbol.getBlock(),
            baseOffset, offset, null, data);
    }

    private byte[] serializeValue(UmasUnlocatedVariableReference symbol, PlcValue value) {
        int dataTypeId = symbol.getDataType();
        if (!UmasDataType.isDefined((short) dataTypeId)) {
            if (value.getRaw() != null) {
                return value.getRaw();
            }
            throw new PlcRuntimeException("Cannot serialize value for unknown data type: " + dataTypeId);
        }
        return serializeForType(UmasDataType.enumForValue((short) dataTypeId), value);
    }

    private static byte[] serializeForType(UmasDataType umasType, PlcValue value) {
        return switch (umasType) {
            case BOOL, EBOOL, UNKNOWN2, UNKNOWN3 ->
                new byte[]{(byte) (value.getBoolean() ? 1 : 0)};
            case BYTE, UNKNOWN11, UNKNOWN12, UNKNOWN13, UNKNOWN17, UNKNOWN18, UNKNOWN19, UNKNOWN20, UNKNOWN24 ->
                new byte[]{value.getByte()};
            case INT -> {
                ByteBuffer buf = ByteBuffer.allocate(2).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putShort(value.getShort());
                yield buf.array();
            }
            case UINT, WORD -> {
                ByteBuffer buf = ByteBuffer.allocate(2).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putShort((short) (value.getInteger() & 0xFFFF));
                yield buf.array();
            }
            case DINT -> {
                ByteBuffer buf = ByteBuffer.allocate(4).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putInt(value.getInteger());
                yield buf.array();
            }
            case UDINT, DWORD -> {
                ByteBuffer buf = ByteBuffer.allocate(4).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putInt((int) (value.getLong() & 0xFFFFFFFFL));
                yield buf.array();
            }
            case REAL -> {
                ByteBuffer buf = ByteBuffer.allocate(4).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putFloat(value.getFloat());
                yield buf.array();
            }
            case STRING -> {
                byte[] strBytes = value.getString().getBytes(StandardCharsets.US_ASCII);
                byte[] result = new byte[strBytes.length + 1];
                System.arraycopy(strBytes, 0, result, 0, strBytes.length);
                yield result;
            }
            case TIME -> {
                // uint32 milliseconds, not BCD.
                ByteBuffer buf = ByteBuffer.allocate(4).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putInt((int) (value.getLong() & 0xFFFFFFFFL));
                yield buf.array();
            }
            case DATE -> {
                java.time.LocalDate date = value.getDate();
                byte[] result = new byte[4];
                result[0] = encodeBcd(date.getDayOfMonth());
                result[1] = encodeBcd(date.getMonthValue());
                int year = date.getYear();
                result[2] = encodeBcd(year % 100);
                result[3] = encodeBcd(year / 100);
                yield result;
            }
            case TOD -> {
                java.time.LocalTime time = value.getTime();
                byte[] result = new byte[4];
                result[0] = encodeBcd((int) ((time.toNanoOfDay() / 10_000_000) % 100));
                result[1] = encodeBcd(time.getSecond());
                result[2] = encodeBcd(time.getMinute());
                result[3] = encodeBcd(time.getHour());
                yield result;
            }
            case DATE_AND_TIME -> {
                java.time.LocalDateTime dt = value.getDateTime();
                byte[] result = new byte[8];
                result[0] = 0x00;
                result[1] = encodeBcd(dt.getSecond());
                result[2] = encodeBcd(dt.getMinute());
                result[3] = encodeBcd(dt.getHour());
                result[4] = encodeBcd(dt.getDayOfMonth());
                result[5] = encodeBcd(dt.getMonthValue());
                int dtYear = dt.getYear();
                result[6] = encodeBcd(dtYear % 100);
                result[7] = encodeBcd(dtYear / 100);
                yield result;
            }
        };
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Browse
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        return executeBrowse(browseRequest, null);
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowseWithInterceptor(PlcBrowseRequest browseRequest,
                                                                          PlcBrowseRequestInterceptor interceptor) {
        return executeBrowse(browseRequest, interceptor);
    }

    private CompletableFuture<PlcBrowseResponse> executeBrowse(PlcBrowseRequest browseRequest,
                                                               PlcBrowseRequestInterceptor interceptor) {
        return CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseCode> responseCodes = new LinkedHashMap<>();
            Map<String, List<PlcBrowseItem>> values = new LinkedHashMap<>();
            for (String queryName : browseRequest.getQueryNames()) {
                try {
                    if (symbolTable.isEmpty()) {
                        loadDataDictionary();
                    }
                    List<PlcBrowseItem> items = convertToBrowseItems(new ArrayList<>(symbolTable.values()));
                    if (interceptor != null) {
                        for (PlcBrowseItem item : items) {
                            interceptor.intercept(queryName, browseRequest.getQuery(queryName), item);
                        }
                    }
                    responseCodes.put(queryName, PlcResponseCode.OK);
                    values.put(queryName, items);
                } catch (Exception e) {
                    LOGGER.error("Browse query '{}' failed: {}", queryName, e.getMessage());
                    responseCodes.put(queryName, PlcResponseCode.REMOTE_ERROR);
                    values.put(queryName, Collections.emptyList());
                }
            }
            return new DefaultPlcBrowseResponse(browseRequest, responseCodes, values);
        });
    }

    /**
     * Downloads the data dictionary (types + symbols) into the session state.
     * Called during handshake so the symbol table is available for
     * read/write/browse immediately.
     */
    private void loadDataDictionary() throws Exception {
        LOGGER.info("Loading data dictionary (types + symbols)...");
        List<UmasDatatypeReference> datatypeRefs = downloadDatatypeNames();
        LOGGER.info("Data dictionary: downloaded {} datatype references", datatypeRefs.size());

        resolveCustomTypes(datatypeRefs);

        List<UmasUnlocatedVariableReference> symbols = downloadSymbolTable();
        LOGGER.info("Data dictionary: downloaded {} symbols", symbols.size());

        for (UmasUnlocatedVariableReference symbol : symbols) {
            symbolTable.put(symbol.getValue().toLowerCase(), symbol);
        }
        computeSymbolSizes();

        LOGGER.info("Data dictionary loaded: {} symbols", symbolTable.size());
    }

    private List<UmasDatatypeReference> downloadDatatypeNames() throws Exception {
        UmasPDUItem request = new UmasPDUReadUnlocatedVariableNamesRequest(
            pairingKey, RECORD_TYPE_DD03, (short) 0x03, hardwareId, 0, 0, null);
        UmasPDUItem responseItem = sendRequest(request, "BrowseDatatypeNames")
            .get(requestTimeoutMs + 1000, TimeUnit.MILLISECONDS);
        if (!(responseItem instanceof UmasPDUReadUnlocatedVariableResponse unlocated)) {
            throw new PlcConnectionException("BrowseDatatypeNames: unexpected response " + responseItem.getClass().getSimpleName());
        }
        byte[] block = unlocated.getBlock();
        if (block == null || block.length == 0) {
            return Collections.emptyList();
        }
        ReadBuffer readBuffer = new ReadBufferByteBased(block, LITTLE_ENDIAN_OPTIONS);
        UmasPDUReadDatatypeNamesResponse parsed = UmasPDUReadDatatypeNamesResponse.staticParse(readBuffer);
        return parsed.getRecords();
    }

    private void resolveCustomTypes(List<UmasDatatypeReference> datatypeRefs) throws Exception {
        for (int i = 0; i < datatypeRefs.size(); i++) {
            UmasDatatypeReference ref = datatypeRefs.get(i);
            int typeId = CUSTOM_TYPE_THRESHOLD + i;
            short primitiveId = ref.getDataType();
            if (UmasDataType.isDefined(primitiveId)) {
                dataTypeTable.put(typeId, UmasDataType.enumForValue(primitiveId));
            }
            // Store the allocated byte size — important for STRING and custom
            // struct types where the size isn't derivable from UmasDataType alone.
            dataTypeSizes.put(typeId, ref.getDataSize());
        }
        for (int i = 0; i < datatypeRefs.size(); i++) {
            UmasDatatypeReference ref = datatypeRefs.get(i);
            int typeId = CUSTOM_TYPE_THRESHOLD + i;
            if (ref.getClassIdentifier() != 0) {
                resolveCustomType(typeId, ref);
            }
        }
    }

    private void resolveCustomType(int typeIndex, UmasDatatypeReference ref) throws Exception {
        UmasPDUItem request = new UmasPDUReadUnlocatedVariableNamesRequest(
            pairingKey, RECORD_TYPE_DD02, (short) 0x03, hardwareId, typeIndex, 0, 0);
        UmasPDUItem responseItem = sendRequest(request, "ResolveType(" + ref.getValue() + ")")
            .get(requestTimeoutMs + 1000, TimeUnit.MILLISECONDS);
        if (!(responseItem instanceof UmasPDUReadUnlocatedVariableResponse unlocated)) {
            LOGGER.warn("ResolveType({}): unexpected response {}", ref.getValue(), responseItem.getClass().getSimpleName());
            return;
        }
        byte[] block = unlocated.getBlock();
        if (block == null || block.length < 2) {
            return;
        }
        int classId = block[0] & 0xFF;
        if (classId == 0x04) {
            ReadBuffer readBuffer = new ReadBufferByteBased(block, LITTLE_ENDIAN_OPTIONS);
            UmasArrayTypeDefinition arrayDef = UmasArrayTypeDefinition.staticParse(readBuffer);
            customTypeNames.put(typeIndex, ref.getValue());
            customTypeFields.put(typeIndex, Collections.emptyList());
            customTypeElementTypeIds.put(typeIndex, arrayDef.getElementTypeId());
            customTypeDimensions.put(typeIndex, arrayDef.getDimensions());
        } else {
            ReadBuffer readBuffer = new ReadBufferByteBased(block, LITTLE_ENDIAN_OPTIONS);
            UmasPDUReadUmasUDTDefinitionResponse udtResponse = UmasPDUReadUmasUDTDefinitionResponse.staticParse(readBuffer);
            customTypeNames.put(typeIndex, ref.getValue());
            customTypeFields.put(typeIndex, udtResponse.getRecords());
        }
    }

    private List<UmasUnlocatedVariableReference> downloadSymbolTable() throws Exception {
        UmasPDUItem request = new UmasPDUReadUnlocatedVariableNamesRequest(
            pairingKey, RECORD_TYPE_DD02, (short) 0x03, hardwareId, SYMBOL_TABLE_BLOCK, 0, 0);
        UmasPDUItem responseItem = sendRequest(request, "BrowseSymbolTable")
            .get(requestTimeoutMs + 1000, TimeUnit.MILLISECONDS);
        if (!(responseItem instanceof UmasPDUReadUnlocatedVariableResponse unlocated)) {
            throw new PlcConnectionException("BrowseSymbolTable: unexpected response " + responseItem.getClass().getSimpleName());
        }
        byte[] block = unlocated.getBlock();
        if (block == null || block.length == 0) {
            return Collections.emptyList();
        }
        ReadBuffer readBuffer = new ReadBufferByteBased(block, LITTLE_ENDIAN_OPTIONS);
        UmasPDUReadUnlocatedVariableNamesResponse parsed =
            UmasPDUReadUnlocatedVariableNamesResponse.staticParse(readBuffer);
        return parsed.getRecords();
    }

    private List<PlcBrowseItem> convertToBrowseItems(List<UmasUnlocatedVariableReference> symbols) {
        List<PlcBrowseItem> items = new ArrayList<>(symbols.size());
        for (UmasUnlocatedVariableReference symbol : symbols) {
            items.add(buildBrowseItem(symbol.getValue(), symbol.getDataType()));
        }
        return items;
    }

    private PlcBrowseItem buildBrowseItem(String name, int dataTypeId) {
        PlcValueType plcValueType;
        List<ArrayInfo> arrayInfo = Collections.emptyList();
        Map<String, PlcBrowseItem> children = Collections.emptyMap();

        Optional<Integer> elementTypeId = Optional.ofNullable(customTypeElementTypeIds.get(dataTypeId));
        if (elementTypeId.isPresent()) {
            plcValueType = resolveValueType(elementTypeId.get());
            arrayInfo = buildArrayInfo(dataTypeId);
            children = buildStructChildren(elementTypeId.get());
        } else if (customTypeFields.containsKey(dataTypeId)) {
            plcValueType = PlcValueType.Struct;
            children = buildStructChildren(dataTypeId);
        } else {
            plcValueType = mapToPlcValueType(dataTypeId);
        }
        SymbolicUmasTag tag = new SymbolicUmasTag(name, plcValueType, Collections.emptyList());
        return new DefaultPlcBrowseItem(
            tag, name,
            true,   // readable
            true,   // writable
            Collections.emptySet(),   // no subscriptions
            false,  // not publishable
            arrayInfo, children, Collections.emptyMap());
    }

    private PlcValueType resolveValueType(int typeId) {
        if (customTypeFields.containsKey(typeId)) {
            return PlcValueType.Struct;
        }
        Integer elementTypeId = customTypeElementTypeIds.get(typeId);
        if (elementTypeId != null) {
            return resolveValueType(elementTypeId);
        }
        return mapToPlcValueType(typeId);
    }

    private List<ArrayInfo> buildArrayInfo(int typeId) {
        List<UmasArrayDimension> dims = customTypeDimensions.get(typeId);
        if (dims == null || dims.isEmpty()) {
            return Collections.emptyList();
        }
        List<ArrayInfo> result = new ArrayList<>();
        for (UmasArrayDimension dim : dims) {
            result.add(new DefaultArrayInfo((int) dim.getStartIndex(), (int) dim.getUpperBound()));
        }
        return result;
    }

    private Map<String, PlcBrowseItem> buildStructChildren(int typeId) {
        List<UmasUDTDefinition> fields = customTypeFields.get(typeId);
        if (fields == null || fields.isEmpty()) {
            return Collections.emptyMap();
        }
        Map<String, PlcBrowseItem> children = new LinkedHashMap<>();
        for (UmasUDTDefinition field : fields) {
            children.put(field.getValue(), buildBrowseItem(field.getValue(), field.getDataType()));
        }
        return children;
    }

    private static PlcValueType mapToPlcValueType(int dataTypeId) {
        if (!UmasDataType.isDefined((short) dataTypeId)) {
            return PlcValueType.RAW_BYTE_ARRAY;
        }
        UmasDataType umasType = UmasDataType.enumForValue((short) dataTypeId);
        return switch (umasType) {
            case BOOL, EBOOL, UNKNOWN2, UNKNOWN3 -> PlcValueType.BOOL;
            case BYTE, UNKNOWN11, UNKNOWN12, UNKNOWN13, UNKNOWN17, UNKNOWN18,
                 UNKNOWN19, UNKNOWN20, UNKNOWN24 -> PlcValueType.BYTE;
            case INT -> PlcValueType.INT;
            case UINT -> PlcValueType.UINT;
            case DINT -> PlcValueType.DINT;
            case UDINT -> PlcValueType.UDINT;
            case REAL -> PlcValueType.REAL;
            case STRING -> PlcValueType.STRING;
            case TIME -> PlcValueType.TIME;
            case DATE -> PlcValueType.DATE;
            case TOD -> PlcValueType.TIME_OF_DAY;
            case DATE_AND_TIME -> PlcValueType.DATE_AND_TIME;
            case WORD -> PlcValueType.WORD;
            case DWORD -> PlcValueType.DWORD;
        };
    }

    /**
     * Computes per-symbol allocated sizes from the symbol-table layout. For
     * each memory block, sort symbols by offset and use the gap to the next
     * symbol as the size. The last symbol in each block doesn't get a size.
     */
    private void computeSymbolSizes() {
        symbolSizes.clear();
        Map<Integer, List<Map.Entry<String, UmasUnlocatedVariableReference>>> byBlock = new HashMap<>();
        for (Map.Entry<String, UmasUnlocatedVariableReference> entry : symbolTable.entrySet()) {
            byBlock.computeIfAbsent(entry.getValue().getBlock(), k -> new ArrayList<>()).add(entry);
        }
        for (var blockEntry : byBlock.entrySet()) {
            var symbols = blockEntry.getValue();
            symbols.sort(Comparator.comparingLong(e -> e.getValue().getOffset()));
            for (int i = 0; i < symbols.size() - 1; i++) {
                var current = symbols.get(i);
                var next = symbols.get(i + 1);
                int size = (int) (next.getValue().getOffset() - current.getValue().getOffset());
                if (size > 0) {
                    symbolSizes.put(current.getKey(), size);
                }
            }
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    /** Reads a little-endian uint32 from the first 4 bytes of a block. */
    private static long readUint32LE(byte[] block) {
        return (block[0] & 0xFFL) | ((block[1] & 0xFFL) << 8)
            | ((block[2] & 0xFFL) << 16) | ((block[3] & 0xFFL) << 24);
    }

    /** Decodes a single BCD-encoded byte (e.g. 0x25 → 25). */
    private static int decodeBcdByte(byte b) {
        return ((b >> 4) & 0x0F) * 10 + (b & 0x0F);
    }

    /** Decodes a BCD-encoded uint16 LE from 2 bytes (e.g. 0x20, 0x25 → 2025). */
    private static int decodeBcd16(byte lo, byte hi) {
        return decodeBcdByte(hi) * 100 + decodeBcdByte(lo);
    }

    /** Encodes a decimal value (0-99) into a BCD byte. e.g. 25 → 0x25. */
    private static byte encodeBcd(int value) {
        return (byte) (((value / 10) << 4) | (value % 10));
    }

}
