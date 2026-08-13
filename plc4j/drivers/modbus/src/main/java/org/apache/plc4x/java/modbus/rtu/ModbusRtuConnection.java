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
package org.apache.plc4x.java.modbus.rtu;

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
import org.apache.plc4x.java.modbus.rtu.config.ModbusRtuConfiguration;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
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
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.apache.commons.lang3.ArrayUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.*;

/**
 * Modbus RTU connection implementation.
 * Handles read, write, and ping operations over Modbus RTU protocol.
 */
public class ModbusRtuConnection extends PollingSubscriptionConnectionBase<ModbusRtuConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(ModbusRtuConnection.class);

    private ModbusRtuMessageCodec messageCodec;
    private final Map<Short, PendingRequest> pendingRequests = new ConcurrentHashMap<>();

    // Serializes request/response transactions. Modbus over serial is a
    // single-outstanding-transaction protocol; without this, concurrent
    // requests to the same unit id silently overwrite each other's entry
    // in pendingRequests and complete with the wrong response.
    private final Object requestChainLock = new Object();
    private CompletableFuture<?> requestTail = CompletableFuture.completedFuture(null);

    public ModbusRtuConnection(ModbusRtuConfiguration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new ModbusRtuMessageCodec(transportInstance, this::handleIncomingMessage);

        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming Modbus data", e);
            }
        });

        LOGGER.info("Modbus RTU connection established");
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT, "Modbus RTU connection established");
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
        pendingRequests.values().forEach(pending ->
            pending.future.completeExceptionally(new PlcRuntimeException("Connection closed")));
        pendingRequests.clear();
        super.close();
        LOGGER.info("Modbus RTU connection closed");
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
            pendingRequests.values().forEach(pending -> pending.future.completeExceptionally(exception));
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

    /**
     * A dispatched request awaiting its response: the future plus the
     * request's function code for response validation.
     */
    private static final class PendingRequest {
        final byte functionFlag;
        final CompletableFuture<ModbusRtuADU> future;
        PendingRequest(byte functionFlag, CompletableFuture<ModbusRtuADU> future) {
            this.functionFlag = functionFlag;
            this.future = future;
        }
    }

    /**
     * Correlates a received frame to the pending request for its unit id.
     * Validation: an exception frame (errorFlag) answers any pending
     * request; a normal frame must carry the SAME function code as the
     * pending request — mismatches (e.g. a late response arriving after
     * its request timed out and a different operation took its place) are
     * discarded without touching the pending entry.
     * <p>
     * Residual limitation: a late response with the SAME function code as
     * the new request is physically indistinguishable — Modbus RTU/ASCII
     * carry no transaction ids. Requests are serialized per connection and
     * shared ports require distinct unit ids per connection, which bounds
     * the exposure to same-fc retry patterns.
     */
    private void handleIncomingMessage(ModbusRtuADU modbusMessage) {
        short address = modbusMessage.getAddress();
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.INCOMING_MESSAGE,
                "Received Modbus RTU response, address=" + address);
        }
        PendingRequest pending = pendingRequests.get(address);
        if (pending == null) {
            LOGGER.warn("Received response for unknown address: {}", address);
            return;
        }
        boolean isException = modbusMessage.getPdu().getErrorFlag();
        if (!isException && modbusMessage.getPdu().getFunctionFlag() != pending.functionFlag) {
            LOGGER.warn("Discarding response for address {} with function code 0x{} while waiting for 0x{} (late response from a timed-out request?)",
                address,
                Integer.toHexString(modbusMessage.getPdu().getFunctionFlag() & 0xFF),
                Integer.toHexString(pending.functionFlag & 0xFF));
            return;
        }
        pendingRequests.remove(address, pending);
        pending.future.complete(modbusMessage);
    }

    // Package-private for tests.
    CompletableFuture<ModbusRtuADU> sendRequest(ModbusRtuADU request, short address) {
        CompletableFuture<ModbusRtuADU> responseFuture = new CompletableFuture<>();

        // Total budget from submission: queueing + dispatch + response.
        long timeoutMs = getConfiguration().getRequestTimeout();
        // Explicit deadline, checked again in dispatchRequest(): the JVM-wide
        // CompletableFuture delay scheduler is a single thread and only gets
        // back around to an already-overdue timeout AFTER it finishes
        // processing the previous request's own completion (cleanup +
        // chain-continuation + submitting the next dispatch) — that hand-off
        // to the async pool routinely wins the race against the scheduler
        // catching up, so responseFuture.isDone() alone can still be false
        // for a request whose enqueue-time budget has, in wall-clock terms,
        // already run out.
        long nowNanos = System.nanoTime();
        long budgetNanos = TimeUnit.MILLISECONDS.toNanos(timeoutMs);
        // A request with only a sliver of its budget left is not worth
        // dispatching either — it has no realistic chance of a round trip
        // completing before its own timeout fires anyway. Backing the
        // dispatch deadline off by this margin also absorbs scheduling
        // jitter between a predecessor's timeout firing (on the JVM-wide
        // single-threaded delay scheduler, see above) and this request's
        // own deadline — queued only a hair later, on the same clock —
        // which would otherwise let a request that is, for all practical
        // purposes, dead slip through a raw "has my deadline passed" check.
        long minViableRemainingNanos = Math.min(TimeUnit.MILLISECONDS.toNanos(50), budgetNanos / 4);
        long dispatchDeadlineNanos = nowNanos + budgetNanos - minViableRemainingNanos;
        responseFuture.orTimeout(timeoutMs, TimeUnit.MILLISECONDS)
            .whenComplete((result, error) -> {
                if (error instanceof TimeoutException) {
                    // Two-arg remove semantics preserved through the holder:
                    // only clear the entry if it is still OUR request.
                    pendingRequests.computeIfPresent(address,
                        (key, pending) -> pending.future == responseFuture ? null : pending);
                }
            });

        CompletableFuture<?> previous;
        synchronized (requestChainLock) {
            previous = requestTail;
            // handle(): the chain continues whether this transaction
            // completes normally, exceptionally, or by timeout.
            requestTail = responseFuture.handle((result, error) -> null);
        }
        if (previous.isDone()) {
            // Hot path: nothing queued — dispatch on the caller thread as before.
            dispatchRequest(request, address, responseFuture, dispatchDeadlineNanos);
        } else {
            // Queued: hop to the async pool so long bursts release the
            // completing thread's stack (a synchronous chain nests one
            // frame per queued request and can silently overflow).
            previous.whenCompleteAsync((ignored, ignoredError) -> dispatchRequest(request, address, responseFuture, dispatchDeadlineNanos));
        }
        return responseFuture;
    }

    private void dispatchRequest(ModbusRtuADU request, short address, CompletableFuture<ModbusRtuADU> responseFuture, long dispatchDeadlineNanos) {
        if (responseFuture.isDone()) {
            // Already completed (timeout or failure) while queued.
            return;
        }
        if (System.nanoTime() >= dispatchDeadlineNanos) {
            // Fail fast: the remaining budget is below the dispatch margin,
            // so the request is semantically dead — don't make the caller
            // wait for the orTimeout backstop.
            responseFuture.completeExceptionally(new TimeoutException(
                "Request timed out while queued (remaining budget below dispatch margin)"));
            return;
        }
        pendingRequests.put(address, new PendingRequest(request.getPdu().getFunctionFlag(), responseFuture));
        try {
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.OUTGOING_MESSAGE,
                    "Sending Modbus RTU request, address=" + address);
            }
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingRequests.computeIfPresent(address, (key, pending) -> pending.future == responseFuture ? null : pending);
            responseFuture.completeExceptionally(new PlcRuntimeException("Failed to send request", e));
        } catch (RuntimeException e) {
            // The responseFuture MUST complete no matter what: the request
            // chain's tail hangs off it, so an unchecked throw here (e.g.
            // from a custom AuditLog) would otherwise silently wedge every
            // subsequent request on this connection.
            pendingRequests.computeIfPresent(address, (key, pending) -> pending.future == responseFuture ? null : pending);
            responseFuture.completeExceptionally(e);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Ping
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        PlcTag pingAddress = new ModbusTagHandler().parseTag(getConfiguration().getPingAddress());
        ModbusPDU readRequestPdu = getReadRequestPdu(pingAddress);
        short unitId = getUnitId(pingAddress);

        ModbusRtuADU modbusRtuADU = new ModbusRtuADU(unitId, readRequestPdu);
        return executeThrottled(() ->
            sendRequest(modbusRtuADU, unitId).thenApply(response ->
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

        // Collect all tags. A tag whose address the builder couldn't parse is kept in the
        // request with an error code and a null tag - it never goes on the wire, its code is
        // reported as-is.
        LinkedHashMap<String, ModbusTag> tagsByName = new LinkedHashMap<>();
        Map<String, PlcResponseItem<PlcValue>> rejectedTags = new LinkedHashMap<>();
        for (String tagName : request.getTagNames()) {
            PlcResponseCode requestCode = request.getTagResponseCode(tagName);
            if (requestCode != PlcResponseCode.OK) {
                rejectedTags.put(tagName, new DefaultPlcResponseItem<>(requestCode, null));
                continue;
            }
            tagsByName.put(tagName, (ModbusTag) request.getTag(tagName));
        }

        // Use the optimizer to merge adjacent tags into block reads
        ModbusReadOptimizer optimizer = new ModbusReadOptimizer(
            getConfiguration().getMaxCoilsPerRequest(),
            getConfiguration().getMaxRegistersPerRequest(),
            getConfiguration().getDefaultPayloadByteOrder());
        List<ModbusReadOptimizer.OptimizedRead> optimizedReads = optimizer.optimizeReads(tagsByName);

        // Execute each optimized block read sequentially (RTU uses unit address for correlation).
        // Chain the block reads so the caller thread returns immediately — executeThrottled()
        // blocks on semaphore acquisition, so iterating synchronously would stall the caller.
        List<CompletableFuture<Map<String, PlcResponseItem<PlcValue>>>> blockFutures = new ArrayList<>();
        CompletableFuture<Void> chain = CompletableFuture.completedFuture(null);
        for (ModbusReadOptimizer.OptimizedRead optimizedRead : optimizedReads) {
            CompletableFuture<Map<String, PlcResponseItem<PlcValue>>> blockFuture =
                chain.thenComposeAsync(v -> executeOptimizedRead(optimizer, optimizedRead));
            blockFutures.add(blockFuture);
            chain = blockFuture.handle((r, e) -> null);
        }

        CompletableFuture<Void> allDone = CompletableFuture.allOf(
            blockFutures.toArray(new CompletableFuture[0]));

        return allDone.thenApply(v -> {
            Map<String, PlcResponseItem<PlcValue>> responseItems = new LinkedHashMap<>(rejectedTags);
            for (CompletableFuture<Map<String, PlcResponseItem<PlcValue>>> blockFuture : blockFutures) {
                try {
                    responseItems.putAll(blockFuture.join());
                } catch (Exception e) {
                    LOGGER.error("Error in optimized block read", e);
                }
            }
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

        ModbusRtuADU modbusRtuADU = new ModbusRtuADU(unitId, requestPdu);
        return executeThrottled(() ->
            sendRequest(modbusRtuADU, unitId).thenApply(responseAdu -> {
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

        if (request.getTagNames().size() == 1) {
            String tagName = request.getTagNames().iterator().next();
            // An unparseable address (or value) is kept in the request with an error code and a
            // null tag - report that code instead of trying to build a PDU from nothing.
            PlcResponseCode requestCode = request.getTagResponseCode(tagName);
            if (requestCode != PlcResponseCode.OK) {
                return CompletableFuture.completedFuture((PlcWriteResponse) new DefaultPlcWriteResponse(
                    request, Collections.singletonMap(tagName, requestCode)));
            }
            PlcTag tag = request.getTag(tagName);
            ModbusPDU requestPdu = getWriteRequestPdu(tag, request.getPlcValue(tagName));
            short unitId = getUnitId(tag);

            ModbusRtuADU modbusRtuADU = new ModbusRtuADU(unitId, requestPdu);
            return executeThrottled(() ->
                sendRequest(modbusRtuADU, unitId).thenApply(responseAdu -> {
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

                    return (PlcWriteResponse) new DefaultPlcWriteResponse(request,
                        Collections.singletonMap(tagName, responseCode));
                })
            );
        } else {
            return CompletableFuture.failedFuture(
                new PlcRuntimeException("Modbus only supports single tag requests"));
        }
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
            return resp.getItems().get(0).getData();
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
            if (((ModbusTag) tag).getDataType() == ModbusDataType.BOOL) {
                // Reverse bits in each byte for coil-style BOOL arrays
                byte[] bytes = new byte[data.length];
                for (int i = 0; i < data.length; i++) {
                    bytes[i] = reverseBitsOfByte(data[i]);
                }
                return bytes;
            }
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
        if (errorResponse.getExceptionCode() == null) {
            LOGGER.warn("Device returned a Modbus exception frame with an unrecognized exception code; mapping to REMOTE_ERROR");
            return PlcResponseCode.REMOTE_ERROR;
        }
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
                new WriteBufferByteBased(new byte[size], WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
            default ->
                new WriteBufferByteBased(new byte[size]);
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
