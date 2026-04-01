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
package org.apache.plc4x.java.umas.readwrite.protocol;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.apache.plc4x.java.umas.readwrite.*;
import org.apache.plc4x.java.umas.readwrite.UmasFunctionKeyTracker;
import org.apache.plc4x.java.umas.readwrite.configuration.UmasConfiguration;
import org.apache.plc4x.java.umas.readwrite.context.UmasDriverContext;
import org.apache.plc4x.java.umas.readwrite.tag.SymbolicUmasTag;
import org.apache.plc4x.java.umas.readwrite.tag.UmasTag;
import org.apache.plc4x.java.umas.readwrite.tag.UmasTagHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

/**
 * Protocol logic for the UMAS driver.
 * Handles connection handshake, read, write, browse, and ping operations.
 */
public class UmasProtocolLogic extends Plc4xProtocolBase<ModbusTcpADU> implements HasConfiguration<UmasConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(UmasProtocolLogic.class);

    private static final int CUSTOM_TYPE_THRESHOLD = 0x1A;
    private static final int RECORD_TYPE_DD02 = 0xDD02;
    private static final int RECORD_TYPE_DD03 = 0xDD03;
    private static final int SYMBOL_TABLE_BLOCK = 0xFFFF;

    private UmasConfiguration configuration;
    private UmasDriverContext umasDriverContext;
    private RequestTransactionManager tm;
    private Duration requestTimeout;
    private short unitIdentifier;

    @Override
    public void setConfiguration(UmasConfiguration configuration) {
        this.configuration = configuration;
        this.requestTimeout = Duration.ofMillis(configuration.getRequestTimeout());
        this.unitIdentifier = (short) configuration.getUnitIdentifier();
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void setDriverContext(org.apache.plc4x.java.spi.context.DriverContext driverContext) {
        super.setDriverContext(driverContext);
        this.umasDriverContext = (UmasDriverContext) driverContext;
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new UmasTagHandler();
    }

    @Override
    public void close(ConversationContext<ModbusTcpADU> context) {
        if (tm != null) {
            tm.shutdown();
        }
    }

    // ========================================================================
    // Connection handshake (fully async — must not block the Netty event loop)
    // ========================================================================

    @Override
    public void onConnect(ConversationContext<ModbusTcpADU> context) {
        // Chain all handshake steps asynchronously to avoid blocking the Netty event loop.
        performPlcIdentAsync(context)
            .thenCompose(v -> performInitCommsAsync(context))
            .thenCompose(v -> performRepeatAsync(context))
            .thenCompose(v -> performReadMemoryBlockAsync(context, 0x30, 0, 33))
            .thenCompose(v -> performProjectInfoAsync(context, (short) 1))
            .thenCompose(v -> performReadMemoryBlockAsync(context, 0x13, 0, 33))
            .thenCompose(v -> performProjectInfoAsync(context, (short) 0))
            .thenCompose(v -> performProjectInfoAsync(context, (short) 4))
            .thenCompose(v -> performProjectInfoAsync(context, (short) 1))
            .thenCompose(v -> performProjectInfoAsync(context, (short) 3))
            .thenAccept(v -> {
                LOGGER.info("UMAS connection established to PLC: hostname={}, model={}, firmware={}",
                    umasDriverContext.getPlcHostname(), umasDriverContext.getPlcModel(),
                    umasDriverContext.getPlcFirmwareVersion());
                context.fireConnected();
            })
            .exceptionally(e -> {
                LOGGER.error("UMAS handshake failed", e);
                context.getChannel().close();
                return null;
            });
    }

    @Override
    public void onDisconnect(ConversationContext<ModbusTcpADU> context) {
        context.fireDisconnected();
    }

    private CompletableFuture<Void> performPlcIdentAsync(ConversationContext<ModbusTcpADU> context) {
        return sendAsyncRequest(context, new UmasPDUPlcIdentRequest((short) 0), "PlcIdent")
            .thenAccept(response -> {
                if (response instanceof UmasPDUPlcIdentResponse identResponse) {
                    umasDriverContext.setPlcHostname(identResponse.getHostname());
                    umasDriverContext.setPlcModel(identResponse.getModel());
                    umasDriverContext.setPlcFirmwareVersion(identResponse.getComVersion());
                    LOGGER.info("PlcIdent: hostname={}, model={}, comVersion={}",
                        identResponse.getHostname(), identResponse.getModel(), identResponse.getComVersion());
                } else {
                    throw new PlcRuntimeException("PlcIdent: unexpected response type: " + response.getClass().getSimpleName());
                }
            });
    }

    private CompletableFuture<Void> performInitCommsAsync(ConversationContext<ModbusTcpADU> context) {
        return sendAsyncRequest(context, new UmasInitCommsRequest((short) 0, (short) 0x00), "InitComms")
            .thenAccept(response -> {
                if (response instanceof UmasInitCommsResponse initResponse) {
                    umasDriverContext.setMaxFrameSize(initResponse.getMaxFrameSize());
                    umasDriverContext.setPlcFirmwareVersion(initResponse.getFirmwareVersion());
                    LOGGER.info("InitComms: maxFrameSize={}, firmwareVersion={}",
                        initResponse.getMaxFrameSize(), initResponse.getFirmwareVersion());
                } else {
                    throw new PlcRuntimeException("InitComms: unexpected response type: " + response.getClass().getSimpleName());
                }
            });
    }

    private CompletableFuture<Void> performRepeatAsync(ConversationContext<ModbusTcpADU> context) {
        int echoSize = umasDriverContext.getMaxFrameSize() - 3;
        byte[] echoData = new byte[echoSize];
        java.util.Arrays.fill(echoData, 1, echoSize, (byte) 0x54);
        return sendAsyncRequest(context, new UmasPDURepeatRequest((short) 0, echoData), "Repeat")
            .thenAccept(response -> {
                if (response instanceof UmasPDURepeatResponse repeatResponse) {
                    LOGGER.info("Repeat: echo OK, {} bytes returned",
                        repeatResponse.getBlock() != null ? repeatResponse.getBlock().length : 0);
                } else {
                    throw new PlcRuntimeException("Repeat: unexpected response type: " + response.getClass().getSimpleName());
                }
            });
    }

    private CompletableFuture<Void> performReadMemoryBlockAsync(ConversationContext<ModbusTcpADU> context, int blockNumber, int offset, int numberOfBytes) {
        String stepName = "ReadMemoryBlock(0x" + String.format("%02X", blockNumber) + ")";
        UmasPDUItem request = new UmasPDUReadMemoryBlockRequest(
            (short) 0, (short) 0x01, blockNumber, offset, 0, numberOfBytes);
        return sendAsyncRequest(context, request, stepName)
            .thenAccept(response -> {
                if (response instanceof UmasPDUReadMemoryBlockResponse readResponse) {
                    byte[] block = readResponse.getBlock();
                    LOGGER.info("{}: {} bytes received", stepName, block != null ? block.length : 0);
                    if (blockNumber == 0x30 && block != null && block.length >= 9) {
                        long hardwareId = (block[5] & 0xFFL) | ((block[6] & 0xFFL) << 8)
                            | ((block[7] & 0xFFL) << 16) | ((block[8] & 0xFFL) << 24);
                        umasDriverContext.setHardwareId(hardwareId);
                        LOGGER.info("{}: extracted hardwareId=0x{}", stepName, String.format("%08X", hardwareId));
                    }
                } else {
                    LOGGER.warn("{}: unexpected response type: {}", stepName, response.getClass().getSimpleName());
                }
            });
    }

    private CompletableFuture<Void> performProjectInfoAsync(ConversationContext<ModbusTcpADU> context, short subcode) {
        return sendAsyncRequest(context, new UmasPDUProjectInfoRequest((short) 0, subcode), "ProjectInfo(subcode=" + subcode + ")")
            .thenAccept(response -> {
                if (response instanceof UmasPDUProjectInfoResponse projectInfoResponse) {
                    byte[] block = projectInfoResponse.getBlock();
                    LOGGER.info("ProjectInfo(subcode={}): {} bytes", subcode, block != null ? block.length : 0);
                } else {
                    LOGGER.warn("ProjectInfo(subcode={}): unexpected response type: {}", subcode, response.getClass().getSimpleName());
                }
            });
    }

    // ========================================================================
    // Async request/response helper (non-blocking, safe for Netty event loop)
    // ========================================================================

    private CompletableFuture<UmasPDUItem> sendAsyncRequest(ConversationContext<ModbusTcpADU> context, UmasPDUItem item, String stepName) {
        int transactionId = umasDriverContext.getNextTransactionId();
        ModbusTcpADU request = buildModbusTcpADU(transactionId, item);

        CompletableFuture<UmasPDUItem> future = new CompletableFuture<>();
        context.sendRequest(request)
            .expectResponse(ModbusTcpADU.class, requestTimeout)
            .onTimeout(e -> future.completeExceptionally(new PlcConnectionException(stepName + " timed out")))
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p.getTransactionIdentifier() == transactionId)
            .handle(p -> {
                try {
                    future.complete(extractUmasResponse(p, stepName));
                } catch (PlcConnectionException e) {
                    future.completeExceptionally(e);
                }
            });

        return future;
    }

    // ========================================================================
    // Ping
    // ========================================================================

    @Override
    public CompletableFuture<PlcPingResponse> ping(PlcPingRequest pingRequest) {
        CompletableFuture<PlcPingResponse> future = new CompletableFuture<>();
        int transactionId = umasDriverContext.getNextTransactionId();
        UmasPDUItem statusRequest = new UmasPDUPlcStatusRequest(umasDriverContext.getPairingKey());
        ModbusTcpADU request = buildModbusTcpADU(transactionId, statusRequest);


        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(request)
            .expectResponse(ModbusTcpADU.class, requestTimeout)
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p.getTransactionIdentifier() == transactionId)
            .handle(p -> {
                transaction.endRequest();
                future.complete(new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK));
            }));
        return future;
    }

    // ========================================================================
    // Read
    // ========================================================================

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;

        // Process tags sequentially via the transaction manager
        CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseItem<PlcValue>> responseItems = new LinkedHashMap<>();
            for (String tagName : request.getTagNames()) {
                PlcTag tag = request.getTag(tagName);
                responseItems.put(tagName, readSingleTag(tagName, tag));
            }
            return new DefaultPlcReadResponse(request, responseItems);
        }).whenComplete((response, throwable) -> {
            if (throwable != null) {
                future.completeExceptionally(throwable);
            } else {
                future.complete(response);
            }
        });

        return future;
    }

    private PlcResponseItem<PlcValue> readSingleTag(String tagName, PlcTag tag) {
        if (!(tag instanceof SymbolicUmasTag symbolicTag)) {
            LOGGER.warn("Read tag '{}' is not a SymbolicUmasTag: {}", tagName, tag.getClass().getSimpleName());
            return new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null);
        }

        String symbolicAddress = symbolicTag.getSymbolicAddress().toLowerCase();
        Optional<UmasUnlocatedVariableReference> symbolOpt = umasDriverContext.getSymbol(symbolicAddress);
        if (symbolOpt.isEmpty()) {
            LOGGER.warn("Read tag '{}': symbol '{}' not found in symbol table", tagName, symbolicAddress);
            return new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null);
        }

        UmasUnlocatedVariableReference symbol = symbolOpt.get();

        try {
            VariableReadRequestReference readRef = buildReadReference(symbol);

            int transactionId = umasDriverContext.getNextTransactionId();
            UmasPDUItem readReq = new UmasPDUReadVariableRequest(
                umasDriverContext.getPairingKey(),
                umasDriverContext.getProjectCrc(),
                (short) 1,
                List.of(readRef));

            ModbusTcpADU modbusTcpADU = buildModbusTcpADU(transactionId, readReq);

            CompletableFuture<ModbusTcpADU> responseFuture = new CompletableFuture<>();
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> conversationContext.sendRequest(modbusTcpADU)
                .expectResponse(ModbusTcpADU.class, requestTimeout)
                .onTimeout(responseFuture::completeExceptionally)
                .onError((p, e) -> responseFuture.completeExceptionally(e))
                .check(p -> p.getTransactionIdentifier() == transactionId)
                .handle(p -> {
                    transaction.endRequest();
                    responseFuture.complete(p);
                }));

            ModbusTcpADU response = responseFuture.get(configuration.getRequestTimeout() + 1000, TimeUnit.MILLISECONDS);
            UmasPDUItem responseItem = extractUmasResponse(response, "ReadVariable(" + tagName + ")");

            if (responseItem instanceof UmasPDUReadVariableResponse readResponse) {
                PlcValue value = parseReadResponse(symbol, readResponse.getBlock());
                return new DefaultPlcResponseItem<>(PlcResponseCode.OK, value);
            } else if (responseItem instanceof UmasPDUErrorResponse) {
                return new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null);
            } else {
                return new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null);
            }
        } catch (Exception e) {
            LOGGER.error("Read tag '{}' failed: {}", tagName, e.getMessage());
            return new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null);
        }
    }

    private VariableReadRequestReference buildReadReference(UmasUnlocatedVariableReference symbol) {
        int dataTypeId = symbol.getDataType();
        byte dataSizeIndex;
        if (UmasDataType.isDefined((short) dataTypeId)) {
            UmasDataType umasType = UmasDataType.enumForValue((short) dataTypeId);
            dataSizeIndex = (byte) umasType.getRequestSize();
        } else {
            dataSizeIndex = (byte) 3;
        }
        return new VariableReadRequestReference(
            (byte) 0, dataSizeIndex, symbol.getBlock(),
            (int) symbol.getOffset(), (short) 0, null);
    }

    private PlcValue parseReadResponse(UmasUnlocatedVariableReference symbol, byte[] block) throws Exception {
        if (block == null || block.length == 0) {
            throw new PlcConnectionException("Read response has empty data block");
        }
        int dataTypeId = symbol.getDataType();
        if (UmasDataType.isDefined((short) dataTypeId)) {
            UmasDataType umasType = UmasDataType.enumForValue((short) dataTypeId);
            ReadBuffer readBuffer = new ReadBufferByteBased(block, ByteOrder.LITTLE_ENDIAN);
            return DataItem.staticParse(readBuffer, umasType, 1);
        }
        return new PlcRawByteArray(block);
    }

    // ========================================================================
    // Write
    // ========================================================================

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;

        CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseCode> responseCodes = new LinkedHashMap<>();
            for (String tagName : request.getTagNames()) {
                PlcTag tag = request.getTag(tagName);
                PlcValue value = request.getPlcValue(tagName);
                responseCodes.put(tagName, writeSingleTag(tagName, tag, value));
            }
            return new DefaultPlcWriteResponse(request, responseCodes);
        }).whenComplete((response, throwable) -> {
            if (throwable != null) {
                future.completeExceptionally(throwable);
            } else {
                future.complete(response);
            }
        });

        return future;
    }

    private PlcResponseCode writeSingleTag(String tagName, PlcTag tag, PlcValue value) {
        if (!(tag instanceof SymbolicUmasTag symbolicTag)) {
            return PlcResponseCode.INVALID_ADDRESS;
        }

        String symbolicAddress = symbolicTag.getSymbolicAddress().toLowerCase();
        Optional<UmasUnlocatedVariableReference> symbolOpt = umasDriverContext.getSymbol(symbolicAddress);
        if (symbolOpt.isEmpty()) {
            return PlcResponseCode.NOT_FOUND;
        }

        UmasUnlocatedVariableReference symbol = symbolOpt.get();

        try {
            byte[] serializedData = serializeValue(symbol, value);
            VariableWriteRequestReference writeRef = buildWriteReference(symbol, serializedData);

            int transactionId = umasDriverContext.getNextTransactionId();
            UmasPDUItem writeReq = new UmasPDUWriteVariableRequest(
                umasDriverContext.getPairingKey(),
                umasDriverContext.getProjectCrc(),
                (short) 1,
                List.of(writeRef));

            ModbusTcpADU modbusTcpADU = buildModbusTcpADU(transactionId, writeReq);

            CompletableFuture<ModbusTcpADU> responseFuture = new CompletableFuture<>();
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> conversationContext.sendRequest(modbusTcpADU)
                .expectResponse(ModbusTcpADU.class, requestTimeout)
                .onTimeout(responseFuture::completeExceptionally)
                .onError((p, e) -> responseFuture.completeExceptionally(e))
                .check(p -> p.getTransactionIdentifier() == transactionId)
                .handle(p -> {
                    transaction.endRequest();
                    responseFuture.complete(p);
                }));

            ModbusTcpADU response = responseFuture.get(configuration.getRequestTimeout() + 1000, TimeUnit.MILLISECONDS);
            UmasPDUItem responseItem = extractUmasResponse(response, "WriteVariable(" + tagName + ")");

            if (responseItem instanceof UmasPDUWriteVariableResponse) {
                return PlcResponseCode.OK;
            } else if (responseItem instanceof UmasPDUErrorResponse) {
                return PlcResponseCode.REMOTE_ERROR;
            } else {
                return PlcResponseCode.INTERNAL_ERROR;
            }
        } catch (Exception e) {
            LOGGER.error("Write tag '{}' failed: {}", tagName, e.getMessage());
            return PlcResponseCode.REMOTE_ERROR;
        }
    }

    private VariableWriteRequestReference buildWriteReference(UmasUnlocatedVariableReference symbol, byte[] data) {
        int dataTypeId = symbol.getDataType();
        byte dataSizeIndex;
        if (UmasDataType.isDefined((short) dataTypeId)) {
            UmasDataType umasType = UmasDataType.enumForValue((short) dataTypeId);
            dataSizeIndex = (byte) umasType.getDataTypeSize();
        } else {
            dataSizeIndex = (byte) data.length;
        }
        return new VariableWriteRequestReference(
            (byte) 0, dataSizeIndex, symbol.getBlock(),
            (int) symbol.getOffset(), 0, null, data);
    }

    private byte[] serializeValue(UmasUnlocatedVariableReference symbol, PlcValue value) throws PlcConnectionException {
        int dataTypeId = symbol.getDataType();
        if (!UmasDataType.isDefined((short) dataTypeId)) {
            if (value.getRaw() != null) {
                return value.getRaw();
            }
            throw new PlcConnectionException("Cannot serialize value for unknown data type: " + dataTypeId);
        }

        UmasDataType umasType = UmasDataType.enumForValue((short) dataTypeId);
        return serializeForType(umasType, value);
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
            case UINT -> {
                ByteBuffer buf = ByteBuffer.allocate(2).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putShort((short) (value.getInteger() & 0xFFFF));
                yield buf.array();
            }
            case DINT -> {
                ByteBuffer buf = ByteBuffer.allocate(4).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putInt(value.getInteger());
                yield buf.array();
            }
            case UDINT -> {
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
            case TIME, DATE, TOD -> {
                ByteBuffer buf = ByteBuffer.allocate(4).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putInt((int) (value.getLong() & 0xFFFFFFFFL));
                yield buf.array();
            }
            case DATE_AND_TIME -> {
                ByteBuffer buf = ByteBuffer.allocate(8).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putLong(value.getLong());
                yield buf.array();
            }
            case WORD -> {
                ByteBuffer buf = ByteBuffer.allocate(2).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putShort((short) (value.getInteger() & 0xFFFF));
                yield buf.array();
            }
            case DWORD -> {
                ByteBuffer buf = ByteBuffer.allocate(4).order(java.nio.ByteOrder.LITTLE_ENDIAN);
                buf.putInt((int) (value.getLong() & 0xFFFFFFFFL));
                yield buf.array();
            }
        };
    }

    // ========================================================================
    // Browse
    // ========================================================================

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        CompletableFuture<PlcBrowseResponse> future = new CompletableFuture<>();

        CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseCode> responseCodes = new LinkedHashMap<>();
            Map<String, List<PlcBrowseItem>> values = new LinkedHashMap<>();

            for (String queryName : browseRequest.getQueryNames()) {
                try {
                    List<PlcBrowseItem> items = executeBrowse();
                    responseCodes.put(queryName, PlcResponseCode.OK);
                    values.put(queryName, items);
                } catch (Exception e) {
                    LOGGER.error("Browse query '{}' failed: {}", queryName, e.getMessage());
                    responseCodes.put(queryName, PlcResponseCode.REMOTE_ERROR);
                    values.put(queryName, Collections.emptyList());
                }
            }
            return new DefaultPlcBrowseResponse(browseRequest, responseCodes, values);
        }).whenComplete((response, throwable) -> {
            if (throwable != null) {
                future.completeExceptionally(throwable);
            } else {
                future.complete(response);
            }
        });

        return future;
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        CompletableFuture<PlcBrowseResponse> future = new CompletableFuture<>();

        CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseCode> responseCodes = new LinkedHashMap<>();
            Map<String, List<PlcBrowseItem>> values = new LinkedHashMap<>();

            for (String queryName : browseRequest.getQueryNames()) {
                try {
                    List<PlcBrowseItem> items = executeBrowse();
                    // Deliver each item through the interceptor
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
        }).whenComplete((response, throwable) -> {
            if (throwable != null) {
                future.completeExceptionally(throwable);
            } else {
                future.complete(response);
            }
        });

        return future;
    }

    private List<PlcBrowseItem> executeBrowse() throws Exception {
        // Phase 1: Download datatype names
        List<UmasDatatypeReference> datatypeRefs = downloadDatatypeNames();
        LOGGER.info("Browse: downloaded {} datatype references", datatypeRefs.size());

        // Phase 2: Resolve custom types
        resolveCustomTypes(datatypeRefs);

        // Phase 3: Download symbol table
        List<UmasUnlocatedVariableReference> symbols = downloadSymbolTable();
        LOGGER.info("Browse: downloaded {} symbols", symbols.size());

        for (UmasUnlocatedVariableReference symbol : symbols) {
            umasDriverContext.addSymbol(symbol.getValue(), symbol);
        }

        return convertToBrowseItems(symbols);
    }

    private List<UmasDatatypeReference> downloadDatatypeNames() throws Exception {
        int transactionId = umasDriverContext.getNextTransactionId();
        UmasPDUItem request = new UmasPDUReadUnlocatedVariableNamesRequest(
            umasDriverContext.getPairingKey(), RECORD_TYPE_DD03, (short) 0x03,
            umasDriverContext.getHardwareId(), 0, 0, null);

        ModbusTcpADU modbusTcpADU = buildModbusTcpADU(transactionId, request);

        CompletableFuture<ModbusTcpADU> responseFuture = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(modbusTcpADU)
            .expectResponse(ModbusTcpADU.class, requestTimeout)
            .onTimeout(responseFuture::completeExceptionally)
            .onError((p, e) -> responseFuture.completeExceptionally(e))
            .check(p -> p.getTransactionIdentifier() == transactionId)
            .handle(p -> {
                transaction.endRequest();
                responseFuture.complete(p);
            }));

        ModbusTcpADU response = responseFuture.get(configuration.getRequestTimeout() + 1000, TimeUnit.MILLISECONDS);
        UmasPDUItem responseItem = extractUmasResponse(response, "BrowseDatatypeNames");

        if (!(responseItem instanceof UmasPDUReadUnlocatedVariableResponse unlocatedResponse)) {
            throw new PlcConnectionException("BrowseDatatypeNames: unexpected response: " + responseItem.getClass().getSimpleName());
        }

        byte[] block = unlocatedResponse.getBlock();
        if (block == null || block.length == 0) {
            return Collections.emptyList();
        }

        ReadBuffer readBuffer = new ReadBufferByteBased(block, ByteOrder.LITTLE_ENDIAN);
        UmasPDUReadDatatypeNamesResponse parsed = UmasPDUReadDatatypeNamesResponse.staticParse(readBuffer);
        return parsed.getRecords();
    }

    private void resolveCustomTypes(List<UmasDatatypeReference> datatypeRefs) throws Exception {
        for (int i = 0; i < datatypeRefs.size(); i++) {
            UmasDatatypeReference ref = datatypeRefs.get(i);
            int typeId = CUSTOM_TYPE_THRESHOLD + i;
            short primitiveId = ref.getDataType();
            if (UmasDataType.isDefined(primitiveId)) {
                umasDriverContext.addDataType(typeId, UmasDataType.enumForValue(primitiveId));
            }
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
        int transactionId = umasDriverContext.getNextTransactionId();
        UmasPDUItem request = new UmasPDUReadUnlocatedVariableNamesRequest(
            umasDriverContext.getPairingKey(), RECORD_TYPE_DD02, (short) 0x03,
            umasDriverContext.getHardwareId(), typeIndex, 0, 0);

        ModbusTcpADU modbusTcpADU = buildModbusTcpADU(transactionId, request);

        CompletableFuture<ModbusTcpADU> responseFuture = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(modbusTcpADU)
            .expectResponse(ModbusTcpADU.class, requestTimeout)
            .onTimeout(responseFuture::completeExceptionally)
            .onError((p, e) -> responseFuture.completeExceptionally(e))
            .check(p -> p.getTransactionIdentifier() == transactionId)
            .handle(p -> {
                transaction.endRequest();
                responseFuture.complete(p);
            }));

        ModbusTcpADU response = responseFuture.get(configuration.getRequestTimeout() + 1000, TimeUnit.MILLISECONDS);
        UmasPDUItem responseItem = extractUmasResponse(response, "ResolveType(" + ref.getValue() + ")");

        if (!(responseItem instanceof UmasPDUReadUnlocatedVariableResponse unlocatedResponse)) {
            LOGGER.warn("ResolveType({}): unexpected response: {}", ref.getValue(), responseItem.getClass().getSimpleName());
            return;
        }

        byte[] block = unlocatedResponse.getBlock();
        if (block == null || block.length < 2) {
            return;
        }

        parseCustomTypeBlock(typeIndex, ref, block);
    }

    private void parseCustomTypeBlock(int typeIndex, UmasDatatypeReference ref, byte[] block) throws Exception {
        int classId = block[0] & 0xFF;

        if (classId == 0x04) {
            ReadBuffer readBuffer = new ReadBufferByteBased(block, ByteOrder.LITTLE_ENDIAN);
            UmasArrayTypeDefinition arrayDef = UmasArrayTypeDefinition.staticParse(readBuffer);
            umasDriverContext.addArrayType(typeIndex, ref.getValue(),
                arrayDef.getElementTypeId(), arrayDef.getDimensions());
        } else {
            ReadBuffer readBuffer = new ReadBufferByteBased(block, ByteOrder.LITTLE_ENDIAN);
            UmasPDUReadUmasUDTDefinitionResponse udtResponse =
                UmasPDUReadUmasUDTDefinitionResponse.staticParse(readBuffer);
            umasDriverContext.addCustomType(typeIndex, ref.getValue(), udtResponse.getRecords());
        }
    }

    private List<UmasUnlocatedVariableReference> downloadSymbolTable() throws Exception {
        int transactionId = umasDriverContext.getNextTransactionId();
        UmasPDUItem request = new UmasPDUReadUnlocatedVariableNamesRequest(
            umasDriverContext.getPairingKey(), RECORD_TYPE_DD02, (short) 0x03,
            umasDriverContext.getHardwareId(), SYMBOL_TABLE_BLOCK, 0, 0);

        ModbusTcpADU modbusTcpADU = buildModbusTcpADU(transactionId, request);

        CompletableFuture<ModbusTcpADU> responseFuture = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(modbusTcpADU)
            .expectResponse(ModbusTcpADU.class, requestTimeout)
            .onTimeout(responseFuture::completeExceptionally)
            .onError((p, e) -> responseFuture.completeExceptionally(e))
            .check(p -> p.getTransactionIdentifier() == transactionId)
            .handle(p -> {
                transaction.endRequest();
                responseFuture.complete(p);
            }));

        ModbusTcpADU response = responseFuture.get(configuration.getRequestTimeout() + 1000, TimeUnit.MILLISECONDS);
        UmasPDUItem responseItem = extractUmasResponse(response, "BrowseSymbolTable");

        if (!(responseItem instanceof UmasPDUReadUnlocatedVariableResponse unlocatedResponse)) {
            throw new PlcConnectionException("BrowseSymbolTable: unexpected response: " + responseItem.getClass().getSimpleName());
        }

        byte[] block = unlocatedResponse.getBlock();
        if (block == null || block.length == 0) {
            return Collections.emptyList();
        }

        ReadBuffer readBuffer = new ReadBufferByteBased(block, ByteOrder.LITTLE_ENDIAN);
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

        Optional<Integer> elementTypeId = umasDriverContext.getArrayElementTypeId(dataTypeId);
        if (elementTypeId.isPresent()) {
            plcValueType = resolveValueType(elementTypeId.get());
            arrayInfo = buildArrayInfo(dataTypeId);
            children = buildStructChildren(elementTypeId.get());
        } else if (umasDriverContext.getCustomTypeFields(dataTypeId).isPresent()) {
            plcValueType = PlcValueType.Struct;
            children = buildStructChildren(dataTypeId);
        } else {
            plcValueType = mapToPlcValueType(dataTypeId);
        }

        SymbolicUmasTag tag = new SymbolicUmasTag(name, plcValueType, Collections.emptyList());
        return new DefaultPlcBrowseItem(tag, name, true, true, true, false,
            arrayInfo, children, Collections.emptyMap());
    }

    private PlcValueType resolveValueType(int typeId) {
        if (umasDriverContext.getCustomTypeFields(typeId).isPresent()) {
            return PlcValueType.Struct;
        }
        if (umasDriverContext.getArrayElementTypeId(typeId).isPresent()) {
            return resolveValueType(umasDriverContext.getArrayElementTypeId(typeId).get());
        }
        return mapToPlcValueType(typeId);
    }

    private List<ArrayInfo> buildArrayInfo(int typeId) {
        Optional<List<UmasArrayDimension>> dims = umasDriverContext.getArrayDimensions(typeId);
        if (dims.isEmpty() || dims.get().isEmpty()) {
            return Collections.emptyList();
        }
        List<ArrayInfo> result = new ArrayList<>();
        for (UmasArrayDimension dim : dims.get()) {
            result.add(new DefaultArrayInfo((int) dim.getStartIndex(), (int) dim.getUpperBound()));
        }
        return result;
    }

    private Map<String, PlcBrowseItem> buildStructChildren(int typeId) {
        Optional<List<UmasUDTDefinition>> fields = umasDriverContext.getCustomTypeFields(typeId);
        if (fields.isEmpty() || fields.get().isEmpty()) {
            return Collections.emptyMap();
        }
        Map<String, PlcBrowseItem> children = new LinkedHashMap<>();
        for (UmasUDTDefinition field : fields.get()) {
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

    // ========================================================================
    // Helpers
    // ========================================================================

    private ModbusTcpADU buildModbusTcpADU(int transactionId, UmasPDUItem item) {
        // Track the function key for response discrimination in the parser
        UmasFunctionKeyTracker.trackRequest(transactionId, item.getUmasFunctionKey());
        UmasPDU umasPdu = new UmasPDU(item);
        return new ModbusTcpADU(transactionId, unitIdentifier, umasPdu);
    }

    private UmasPDUItem extractUmasResponse(ModbusTcpADU response, String stepName) throws PlcConnectionException {
        ModbusPDU pdu = response.getPdu();
        if (pdu instanceof ModbusPDUError errorPdu) {
            throw new PlcConnectionException(stepName + " received Modbus error: " + errorPdu.getExceptionCode());
        }
        if (!(pdu instanceof UmasPDU umasPdu)) {
            throw new PlcConnectionException(stepName + " received unexpected PDU type: " + pdu.getClass().getSimpleName());
        }
        return umasPdu.getItem();
    }

}
