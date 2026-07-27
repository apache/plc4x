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
package org.apache.plc4x.java.modbus.tcp;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.optimizer.ModbusReadOptimizer;
import org.apache.plc4x.java.modbus.base.tag.*;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.modbus.tcp.config.ModbusTcpConfiguration;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.utils.subscriptionemulation.PollingSubscriptionConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.*;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.apache.commons.lang3.ArrayUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Modbus TCP connection implementation.
 * Handles read, write, and ping operations over Modbus TCP protocol.
 */
public class ModbusTcpConnection extends PollingSubscriptionConnectionBase<ModbusTcpConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(ModbusTcpConnection.class);

    private final AtomicInteger transactionIdentifierGenerator = new AtomicInteger(1);
    private ModbusTcpMessageCodec messageCodec;
    private final Map<Integer, CompletableFuture<ModbusTcpADU>> pendingRequests = new ConcurrentHashMap<>();

    public ModbusTcpConnection(ModbusTcpConfiguration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new ModbusTcpMessageCodec(transportInstance, this::handleIncomingMessage);

        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming Modbus data", e);
            }
        });

        LOGGER.info("Modbus TCP connection established");
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT, "Modbus TCP connection established");
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    @Override
    public boolean isConnected() {
        return messageCodec != null && messageCodec.isOpen();
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
        LOGGER.info("Modbus TCP connection closed");
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");

        PlcRuntimeException exception = new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause);
        int pendingCount = pendingRequests.size();
        if (pendingCount > 0) {
            LOGGER.warn("Failing {} pending requests due to transport disconnect", pendingCount);
            pendingRequests.values().forEach(future -> future.completeExceptionally(exception));
            pendingRequests.clear();
        }
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new ModbusTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    protected int getMaxConcurrentRequests() {
        return 1;
    }

    private void handleIncomingMessage(ModbusTcpADU modbusMessage) {
        int transactionId = modbusMessage.getTransactionIdentifier();
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.INCOMING_MESSAGE,
                "Received Modbus TCP response, txId=" + transactionId);
        }
        CompletableFuture<ModbusTcpADU> future = pendingRequests.remove(transactionId);
        if (future != null) {
            future.complete(modbusMessage);
        } else {
            LOGGER.warn("Received response for unknown transaction ID: {}", transactionId);
        }
    }

    private int nextTransactionId() {
        int id = transactionIdentifierGenerator.getAndIncrement();
        if (transactionIdentifierGenerator.get() == 0xFFFF) {
            transactionIdentifierGenerator.set(1);
        }
        return id;
    }

    private CompletableFuture<ModbusTcpADU> sendRequest(ModbusTcpADU request, int transactionId) {
        CompletableFuture<ModbusTcpADU> responseFuture = new CompletableFuture<>();
        pendingRequests.put(transactionId, responseFuture);

        try {
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.OUTGOING_MESSAGE,
                    "Sending Modbus TCP request, txId=" + transactionId);
            }
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingRequests.remove(transactionId);
            responseFuture.completeExceptionally(new PlcRuntimeException("Failed to send request", e));
            return responseFuture;
        }

        long timeoutMs = getConfiguration().getRequestTimeout();
        responseFuture.orTimeout(timeoutMs, TimeUnit.MILLISECONDS)
            .whenComplete((result, error) -> {
                if (error instanceof TimeoutException) {
                    pendingRequests.remove(transactionId);
                }
            });

        return responseFuture;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Ping
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        PlcTag pingAddress = new ModbusTagHandler().parseTag(getConfiguration().getPingAddress());
        ModbusPDU readRequestPdu = getReadRequestPdu(pingAddress);
        short unitId = getUnitId(pingAddress);
        int transactionId = nextTransactionId();

        ModbusTcpADU modbusTcpADU = new ModbusTcpADU(transactionId, unitId, readRequestPdu);
        return executeThrottled(() ->
            sendRequest(modbusTcpADU, transactionId).thenApply(response ->
                new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK)
            )
        );
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;

        // Collect all tags
        LinkedHashMap<String, ModbusTag> tagsByName = new LinkedHashMap<>();
        for (String tagName : request.getTagNames()) {
            tagsByName.put(tagName, (ModbusTag) request.getTag(tagName));
        }

        // Use the optimizer to merge adjacent tags into block reads
        ModbusReadOptimizer optimizer = new ModbusReadOptimizer(
            getConfiguration().getMaxCoilsPerRequest(),
            getConfiguration().getMaxRegistersPerRequest(),
            getConfiguration().getDefaultPayloadByteOrder());
        List<ModbusReadOptimizer.OptimizedRead> optimizedReads = optimizer.optimizeReads(tagsByName);

        // Execute each optimized block read, chaining them so the caller thread returns
        // immediately — executeThrottled() blocks on semaphore acquisition, so iterating
        // synchronously would stall the caller until all block reads complete.
        List<CompletableFuture<Map<String, PlcResponseItem<PlcValue>>>> blockFutures = new ArrayList<>();
        CompletableFuture<Void> chain = CompletableFuture.completedFuture(null);
        for (ModbusReadOptimizer.OptimizedRead optimizedRead : optimizedReads) {
            CompletableFuture<Map<String, PlcResponseItem<PlcValue>>> blockFuture =
                chain.thenComposeAsync(v -> executeOptimizedRead(optimizer, optimizedRead));
            blockFutures.add(blockFuture);
            chain = blockFuture.handle((r, e) -> null);
        }

        // Merge all results
        CompletableFuture<Void> allDone = CompletableFuture.allOf(
            blockFutures.toArray(new CompletableFuture[0]));

        return allDone.thenApply(v -> {
            Map<String, PlcResponseItem<PlcValue>> responseItems = new LinkedHashMap<>();
            for (CompletableFuture<Map<String, PlcResponseItem<PlcValue>>> blockFuture : blockFutures) {
                try {
                    responseItems.putAll(blockFuture.join());
                } catch (Exception e) {
                    LOGGER.error("Error in optimized block read", e);
                }
            }
            // Ensure all original tags have a response (in case of optimizer issues)
            for (String tagName : request.getTagNames()) {
                if (!responseItems.containsKey(tagName)) {
                    responseItems.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
            }
            return (PlcReadResponse) new DefaultPlcReadResponse(request, responseItems);
        });
    }

    private CompletableFuture<Map<String, PlcResponseItem<PlcValue>>> executeOptimizedRead(
            ModbusReadOptimizer optimizer, ModbusReadOptimizer.OptimizedRead optimizedRead) {
        ModbusTag mergedTag = optimizedRead.mergedTag;
        ModbusPDU requestPdu = getReadRequestPdu(mergedTag);
        short unitId = getUnitId(mergedTag);
        int transactionId = nextTransactionId();

        ModbusTcpADU modbusTcpADU = new ModbusTcpADU(transactionId, unitId, requestPdu);
        return executeThrottled(() ->
            sendRequest(modbusTcpADU, transactionId).thenApply(responseAdu -> {
                ModbusPDU responsePdu = responseAdu.getPdu();

                if (responsePdu instanceof ModbusPDUError errorResponse) {
                    PlcResponseCode errorCode = getErrorCode(errorResponse);
                    return optimizer.splitResponse(optimizedRead, errorCode, null);
                }

                byte[] blockData = extractResponseData(requestPdu, responsePdu);
                if (blockData == null) {
                    return optimizer.splitResponse(optimizedRead, PlcResponseCode.INTERNAL_ERROR, null);
                }

                if (auditLog.isEnabled()) {
                    auditLog.write(AuditLogEventType.API_RESPONSE,
                        "Block read response: " + blockData.length + " bytes for " +
                            optimizedRead.originalTagNames.size() + " tags");
                }

                return optimizer.splitResponse(optimizedRead, PlcResponseCode.OK, blockData);
            })
        );
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;

        // Modbus only supports one operation per PDU, so we send each tag as a separate request
        // and merge the results into a single response. Chain the per-tag writes so the caller
        // thread returns immediately — executeThrottled() blocks on semaphore acquisition, so
        // iterating synchronously would stall the caller until all tags complete.
        LinkedHashMap<String, CompletableFuture<PlcResponseCode>> tagFutures = new LinkedHashMap<>();
        CompletableFuture<Void> chain = CompletableFuture.completedFuture(null);
        for (String tagName : request.getTagNames()) {
            PlcTag tag = request.getTag(tagName);
            PlcValue value = request.getPlcValue(tagName);
            CompletableFuture<PlcResponseCode> tagFuture =
                chain.thenComposeAsync(v -> writeSingleTag(tagName, tag, value));
            tagFutures.put(tagName, tagFuture);
            chain = tagFuture.handle((r, e) -> null);
        }

        CompletableFuture<Void> allDone = CompletableFuture.allOf(
            tagFutures.values().toArray(new CompletableFuture[0]));

        return allDone.thenApply(v -> {
            Map<String, PlcResponseCode> responseCodes = new LinkedHashMap<>();
            for (Map.Entry<String, CompletableFuture<PlcResponseCode>> entry : tagFutures.entrySet()) {
                try {
                    responseCodes.put(entry.getKey(), entry.getValue().join());
                } catch (Exception e) {
                    responseCodes.put(entry.getKey(), PlcResponseCode.INTERNAL_ERROR);
                }
            }
            return (PlcWriteResponse) new DefaultPlcWriteResponse(request, responseCodes);
        });
    }

    private CompletableFuture<PlcResponseCode> writeSingleTag(String tagName, PlcTag tag, PlcValue plcValue) {
        ModbusPDU requestPdu = getWriteRequestPdu(tag, plcValue);
        short unitId = getUnitId(tag);
        int transactionId = nextTransactionId();

        ModbusTcpADU modbusTcpADU = new ModbusTcpADU(transactionId, unitId, requestPdu);
        return executeThrottled(() ->
            sendRequest(modbusTcpADU, transactionId).thenApply(responseAdu -> {
                ModbusPDU responsePdu = responseAdu.getPdu();
                PlcResponseCode responseCode;

                if (responsePdu instanceof ModbusPDUError errorResponse) {
                    responseCode = getErrorCode(errorResponse);
                } else {
                    responseCode = PlcResponseCode.OK;
                    if (responsePdu instanceof ModbusPDUWriteSingleCoilResponse response) {
                        ModbusPDUWriteSingleCoilRequest requestSingleCoil = (ModbusPDUWriteSingleCoilRequest) requestPdu;
                        if (!((response.getValue() == requestSingleCoil.getValue()) &&
                            (response.getAddress() == requestSingleCoil.getAddress()))) {
                            responseCode = PlcResponseCode.REMOTE_ERROR;
                        }
                    }
                }

                if (auditLog.isEnabled()) {
                    auditLog.write(AuditLogEventType.API_RESPONSE,
                        "Write response for '" + tagName + "': " + responseCode);
                }

                return responseCode;
            })
        );
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus Protocol Helpers (inlined from the old ModbusProtocolLogic)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private short getUnitId(PlcTag tag) {
        if (tag instanceof ModbusTag modbusTag) {
            Short unitId = modbusTag.getUnitId();
            return unitId != null ? unitId : (short) getConfiguration().getDefaultUnitIdentifier();
        }
        return (short) getConfiguration().getDefaultUnitIdentifier();
    }

    private ModbusPDU getReadRequestPdu(PlcTag tag) {
        if (tag instanceof ModbusTagDiscreteInput discreteInput) {
            return new ModbusPDUReadDiscreteInputsRequest(discreteInput.getAddress(), discreteInput.getNumberOfElements());
        } else if (tag instanceof ModbusTagCoil coil) {
            return new ModbusPDUReadCoilsRequest(coil.getAddress(), coil.getNumberOfElements());
        } else if (tag instanceof ModbusTagInputRegister inputRegister) {
            return new ModbusPDUReadInputRegistersRequest(inputRegister.getAddress(), Math.max(inputRegister.getLengthWords(), 1));
        } else if (tag instanceof ModbusTagHoldingRegister holdingRegister) {
            return new ModbusPDUReadHoldingRegistersRequest(holdingRegister.getAddress(), Math.max(holdingRegister.getLengthWords(), 1));
        } else if (tag instanceof ModbusTagExtendedRegister extendedRegister) {
            int group1Address = extendedRegister.getAddress() % 10000;
            int group1Quantity;
            int group1FileNumber = (int) (Math.floor((float) extendedRegister.getAddress() / 10000) + 1);
            List<ModbusPDUReadFileRecordRequestItem> itemArray;

            if ((group1Address + extendedRegister.getLengthWords()) <= 10000) {
                group1Quantity = extendedRegister.getLengthWords();
                ModbusPDUReadFileRecordRequestItem group1 =
                    new ModbusPDUReadFileRecordRequestItem((short) 6, group1FileNumber, group1Address, group1Quantity);
                itemArray = Collections.singletonList(group1);
            } else {
                group1Quantity = 10000 - group1Address;
                int group2Quantity = extendedRegister.getLengthWords() - group1Quantity;
                int group2FileNumber = group1FileNumber + 1;
                ModbusPDUReadFileRecordRequestItem group1 =
                    new ModbusPDUReadFileRecordRequestItem((short) 6, group1FileNumber, group1Address, group1Quantity);
                ModbusPDUReadFileRecordRequestItem group2 =
                    new ModbusPDUReadFileRecordRequestItem((short) 6, group2FileNumber, 0, group2Quantity);
                itemArray = Arrays.asList(group1, group2);
            }
            return new ModbusPDUReadFileRecordRequest(itemArray);
        }
        throw new PlcRuntimeException("Unsupported read tag type " + tag.getClass().getName());
    }

    private ModbusPDU getWriteRequestPdu(PlcTag tag, PlcValue plcValue) {
        if (tag instanceof ModbusTagCoil coil) {
            ModbusByteOrder byteOrder = getEffectiveByteOrder(coil);
            ModbusPDUWriteMultipleCoilsRequest request =
                new ModbusPDUWriteMultipleCoilsRequest(coil.getAddress(), coil.getNumberOfElements(),
                    fromPlcValue(tag, plcValue, byteOrder));
            if (request.getQuantity() == coil.getNumberOfElements()) {
                return request;
            } else {
                throw new PlcRuntimeException("Number of requested bytes (" + request.getQuantity() +
                    ") doesn't match number of requested addresses (" + coil.getNumberOfElements() + ")");
            }
        } else if (tag instanceof ModbusTagHoldingRegister holdingRegister) {
            ModbusByteOrder byteOrder = getEffectiveByteOrder(holdingRegister);
            byte[] bytes = fromPlcValue(tag, plcValue, byteOrder);
            if (bytes.length == 2) {
                int value = ((bytes[0] & 0xFF) << 8) | (bytes[1] & 0xFF);
                return new ModbusPDUWriteSingleRegisterRequest(holdingRegister.getAddress(), value);
            } else {
                ModbusPDUWriteMultipleHoldingRegistersRequest request =
                    new ModbusPDUWriteMultipleHoldingRegistersRequest(holdingRegister.getAddress(),
                        holdingRegister.getLengthWords(), bytes);
                if (request.getValue().length == holdingRegister.getLengthWords() * 2) {
                    return request;
                } else {
                    throw new PlcRuntimeException("Number of requested values (" + request.getValue().length / 2 +
                        ") doesn't match number of requested addresses (" + holdingRegister.getLengthWords() + ")");
                }
            }
        } else if (tag instanceof ModbusTagExtendedRegister extendedRegister) {
            ModbusByteOrder byteOrder = getEffectiveByteOrder(extendedRegister);
            int group1Address = extendedRegister.getAddress() % FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH;
            int group1FileNumber = (int) (Math.floor((float) extendedRegister.getAddress() / FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH) + 1);
            List<ModbusPDUWriteFileRecordRequestItem> itemArray;
            if ((group1Address + extendedRegister.getLengthWords()) <= FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH) {
                ModbusPDUWriteFileRecordRequestItem group1 = new ModbusPDUWriteFileRecordRequestItem(
                    (short) 6, group1FileNumber, group1Address, fromPlcValue(tag, plcValue, byteOrder));
                itemArray = Collections.singletonList(group1);
            } else {
                int group1Quantity = FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH - group1Address;
                int group2FileNumber = group1FileNumber + 1;
                byte[] allBytes = fromPlcValue(tag, plcValue, byteOrder);
                byte[] plcValue1 = ArrayUtils.subarray(allBytes, 0, group1Quantity * 2);
                byte[] plcValue2 = ArrayUtils.subarray(allBytes, group1Quantity * 2, allBytes.length);
                ModbusPDUWriteFileRecordRequestItem group1 = new ModbusPDUWriteFileRecordRequestItem(
                    (short) 6, group1FileNumber, group1Address, plcValue1);
                ModbusPDUWriteFileRecordRequestItem group2 = new ModbusPDUWriteFileRecordRequestItem(
                    (short) 6, group2FileNumber, 0, plcValue2);
                itemArray = Arrays.asList(group1, group2);
            }
            return new ModbusPDUWriteFileRecordRequest(itemArray);
        }
        throw new PlcRuntimeException("Unsupported write tag type " + tag.getClass().getName());
    }

    private PlcValue toPlcValue(ModbusPDU request, ModbusPDU response, ModbusDataType dataType, int numberOfElements, ModbusByteOrder byteOrder) throws BufferException {
        byte[] responseData = extractResponseData(request, response);
        if (responseData == null) {
            return null;
        }

        // Apply byte-swap if needed (before parsing)
        if (byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP || byteOrder == ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP) {
            responseData = byteSwap(responseData);
        }

        boolean bigEndian = (byteOrder == ModbusByteOrder.BIG_ENDIAN || byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP);
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(responseData);
        return DataItem.staticParse(readBuffer, dataType, numberOfElements, bigEndian);
    }

    private byte[] extractResponseData(ModbusPDU request, ModbusPDU response) {
        if (request instanceof ModbusPDUReadDiscreteInputsRequest && response instanceof ModbusPDUReadDiscreteInputsResponse resp) {
            return resp.getValue();
        } else if (request instanceof ModbusPDUReadCoilsRequest && response instanceof ModbusPDUReadCoilsResponse resp) {
            return resp.getValue();
        } else if (request instanceof ModbusPDUReadInputRegistersRequest && response instanceof ModbusPDUReadInputRegistersResponse resp) {
            return resp.getValue();
        } else if (request instanceof ModbusPDUReadHoldingRegistersRequest && response instanceof ModbusPDUReadHoldingRegistersResponse resp) {
            return resp.getValue();
        } else if (request instanceof ModbusPDUReadFileRecordRequest && response instanceof ModbusPDUReadFileRecordResponse resp) {
            return resp.getItems().getFirst().getData();
        }
        LOGGER.warn("Unexpected response type {} for request type {}", response.getClass().getSimpleName(), request.getClass().getSimpleName());
        return null;
    }

    private byte[] fromPlcValue(PlcTag tag, PlcValue plcValue, ModbusByteOrder byteOrder) {
        ModbusDataType tagDataType = ((ModbusTag) tag).getDataType();
        try {
            if (tag instanceof ModbusTagCoil) {
                return fromPlcValueCoil(plcValue, byteOrder);
            }
            boolean bigEndian = (byteOrder == ModbusByteOrder.BIG_ENDIAN || byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP);
            int size = DataItem.getLengthInBytes(plcValue, tagDataType, plcValue.getLength(), bigEndian);
            WriteBufferByteBased writeBuffer = createWriteBuffer(size, byteOrder);
            DataItem.staticSerialize(writeBuffer, plcValue, tagDataType, plcValue.getLength(), bigEndian);
            byte[] data = writeBuffer.getBytes();
            if (byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP || byteOrder == ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP) {
                data = byteSwap(data);
            }
            // Note: bit reversal for coil BOOL arrays is handled in fromPlcValueCoil, not here.
            // Holding register BOOLs use DataItem serialization which handles them correctly.
            return data;
        } catch (BufferException e) {
            throw new PlcRuntimeException("Unable to serialize PlcValue: " + e.getMessage(), e);
        }
    }

    private byte[] fromPlcValueCoil(PlcValue plcValue, ModbusByteOrder byteOrder) throws BufferException {
        if (plcValue instanceof PlcBOOL) {
            return new byte[]{(byte) (plcValue.getBoolean() ? 1 : 0)};
        } else if (plcValue instanceof PlcList valueList) {
            WriteBufferByteBased wb = createWriteBuffer(((plcValue.getLength() - 1) / 8) + 1, byteOrder);
            int paddingBits = 8 - (plcValue.getLength() % 8);
            if (paddingBits < 8) {
                for (int i = 0; i < paddingBits; i++) {
                    wb.writeBit(false);
                }
            }
            for (int i = 0; i < plcValue.getLength(); i++) {
                PlcValue value = valueList.getIndex((plcValue.getLength() - 1) - i);
                if (!(value instanceof PlcBOOL)) {
                    throw new PlcRuntimeException("Expecting only BOOL values when writing coils.");
                }
                wb.writeBit(((PlcBOOL) value).getBoolean());
            }
            byte[] bytes = wb.getBytes();
            if (byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP || byteOrder == ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP) {
                bytes = byteSwap(bytes);
            }
            ArrayUtils.reverse(bytes);
            return bytes;
        }
        throw new PlcRuntimeException("Expecting only BOOL or List values when writing coils.");
    }

    private PlcResponseCode getErrorCode(ModbusPDUError errorResponse) {
        return switch (errorResponse.getExceptionCode()) {
            case ILLEGAL_FUNCTION -> PlcResponseCode.UNSUPPORTED;
            case ILLEGAL_DATA_ADDRESS -> PlcResponseCode.INVALID_ADDRESS;
            case ILLEGAL_DATA_VALUE -> PlcResponseCode.INVALID_DATA;
            case SLAVE_DEVICE_FAILURE -> PlcResponseCode.REMOTE_ERROR;
            case ACKNOWLEDGE -> PlcResponseCode.OK;
            case SLAVE_DEVICE_BUSY -> PlcResponseCode.REMOTE_BUSY;
            case NEGATIVE_ACKNOWLEDGE -> PlcResponseCode.REMOTE_ERROR;
            case MEMORY_PARITY_ERROR -> PlcResponseCode.INTERNAL_ERROR;
            case GATEWAY_PATH_UNAVAILABLE -> PlcResponseCode.INTERNAL_ERROR;
            case GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND -> PlcResponseCode.REMOTE_ERROR;
        };
    }

    private ModbusByteOrder getEffectiveByteOrder(ModbusTag tag) {
        ModbusByteOrder byteOrder = getConfiguration().getDefaultPayloadByteOrder();
        if (tag.getByteOrder() != null) {
            byteOrder = tag.getByteOrder();
        }
        return byteOrder;
    }

    private WriteBufferByteBased createWriteBuffer(int size, ModbusByteOrder byteOrder) {
        return switch (byteOrder) {
            case LITTLE_ENDIAN, LITTLE_ENDIAN_BYTE_SWAP ->
                new WriteBufferByteBased(new byte[size],
                    WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
                    WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                    WithOption.WithSignedIntegerEncoding("twos-complement"),
                    WithOption.WithFloatEncoding("IEEE754"),
                    WithOption.WithStringEncoding("UTF8"));
            default ->
                new WriteBufferByteBased(new byte[size],
                    WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                    WithOption.WithSignedIntegerEncoding("twos-complement"),
                    WithOption.WithFloatEncoding("IEEE754"),
                    WithOption.WithStringEncoding("UTF8"));
        };
    }

    private static byte[] byteSwap(byte[] in) {
        byte[] out = new byte[in.length];
        for (int i = 0; i < out.length - 1; i += 2) {
            out[i] = in[i + 1];
            out[i + 1] = in[i];
        }
        // Handle odd-length arrays
        if (in.length % 2 != 0) {
            out[in.length - 1] = in[in.length - 1];
        }
        return out;
    }

    private static byte reverseBitsOfByte(byte b) {
        java.util.BitSet bits = java.util.BitSet.valueOf(new byte[]{b});
        java.util.BitSet reverse = java.util.BitSet.valueOf(new byte[]{(byte) 0xFF});
        for (int j = 0; j < 8; j++) {
            reverse.set(j, bits.get(7 - j));
        }
        return java.util.Arrays.copyOf(reverse.toByteArray(), 1)[0];
    }

    private static final int FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH = 10000;

}
