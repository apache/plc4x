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
package org.apache.plc4x.java.abeth;

import org.apache.plc4x.java.abeth.configuration.AbEthConfiguration;
import org.apache.plc4x.java.abeth.readwrite.CIPEncapsulationConnectionRequest;
import org.apache.plc4x.java.abeth.readwrite.CIPEncapsulationConnectionResponse;
import org.apache.plc4x.java.abeth.readwrite.CIPEncapsulationPacket;
import org.apache.plc4x.java.abeth.readwrite.CIPEncapsulationReadRequest;
import org.apache.plc4x.java.abeth.readwrite.CIPEncapsulationReadResponse;
import org.apache.plc4x.java.abeth.readwrite.DF1CommandRequestMessage;
import org.apache.plc4x.java.abeth.readwrite.DF1CommandResponseMessageProtectedTypedLogicalRead;
import org.apache.plc4x.java.abeth.readwrite.DF1RequestMessage;
import org.apache.plc4x.java.abeth.readwrite.DF1RequestProtectedTypedLogicalRead;
import org.apache.plc4x.java.abeth.tag.AbEthTag;
import org.apache.plc4x.java.abeth.tag.AbEthTagHandler;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.utils.subscriptionemulation.PollingSubscriptionConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Allen Bradley AB-ETH connection. Direct port of the legacy
 * {@code AbEthProtocolLogic} to the SPI3 {@link ConnectionBase} model.
 *
 * <p>Read-only driver. On connect we send a {@code
 * CIPEncapsulationConnectionRequest} and wait for the
 * {@code CIPEncapsulationConnectionResponse} which carries the session
 * handle; every subsequent read carries that handle and a fresh
 * 16-bit transaction counter used to match responses to in-flight reads.</p>
 */
public class AbEthConnection extends PollingSubscriptionConnectionBase<AbEthConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(AbEthConnection.class);

    /** Sender context the legacy driver used for the connection request. */
    private static final List<Short> CONNECTION_REQUEST_SENDER_CONTEXT = Arrays.asList(
        (short) 0x00, (short) 0x04, (short) 0x00, (short) 0x05,
        (short) 0x00, (short) 0x00, (short) 0x00, (short) 0x00);

    /** Sender context for follow-up read requests. */
    private static final List<Short> EMPTY_SENDER_CONTEXT = Arrays.asList(
        (short) 0x00, (short) 0x00, (short) 0x00, (short) 0x00,
        (short) 0x00, (short) 0x00, (short) 0x00, (short) 0x00);

    private AbEthMessageCodec messageCodec;

    private volatile boolean handshakeComplete = false;
    private volatile long sessionHandle;

    /** Pending connection-response future (handshake). */
    private volatile CompletableFuture<CIPEncapsulationConnectionResponse> pendingConnection;
    /** Pending read futures keyed by transaction counter. */
    private final Map<Integer, CompletableFuture<CIPEncapsulationReadResponse>> pendingReads =
        new ConcurrentHashMap<>();

    private final AtomicInteger transactionCounterGenerator = new AtomicInteger(1);

    public AbEthConnection(AbEthConfiguration configuration,
                           TransportInstance<?> transportInstance,
                           AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new AbEthTagHandler();
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
        messageCodec = new AbEthMessageCodec(transportInstance, this::handleIncomingMessage);

        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming AB-ETH data", e);
            }
        });

        try {
            doHandshake().get(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("Error during AB-ETH connect handshake", e);
        }

        handshakeComplete = true;
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT,
                "AB-ETH connection established (session handle " + sessionHandle + ")");
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    private CompletableFuture<Void> doHandshake() {
        pendingConnection = new CompletableFuture<>();
        CIPEncapsulationConnectionRequest request = new CIPEncapsulationConnectionRequest(
            0L, 0L, CONNECTION_REQUEST_SENDER_CONTEXT, 0L);
        try {
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingConnection.completeExceptionally(new PlcRuntimeException(
                "Failed to send AB-ETH connection request", e));
            return pendingConnection.thenApply(r -> null);
        }
        return pendingConnection
            .orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
            .thenAccept(response -> sessionHandle = response.getSessionHandle());
    }

    @Override
    public void close() throws Exception {
        handshakeComplete = false;
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        PlcRuntimeException closed = new PlcRuntimeException("Connection closed");
        if (pendingConnection != null && !pendingConnection.isDone()) {
            pendingConnection.completeExceptionally(closed);
        }
        pendingReads.values().forEach(f -> f.completeExceptionally(closed));
        pendingReads.clear();
        super.close();
        LOGGER.info("AB-ETH connection closed");
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
        if (pendingConnection != null && !pendingConnection.isDone()) {
            pendingConnection.completeExceptionally(ex);
        }
        pendingReads.values().forEach(f -> f.completeExceptionally(ex));
        pendingReads.clear();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming messages
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingMessage(CIPEncapsulationPacket packet) {
        if (packet instanceof CIPEncapsulationConnectionResponse cr) {
            if (pendingConnection != null && !pendingConnection.isDone()) {
                pendingConnection.complete(cr);
            }
            return;
        }
        if (packet instanceof CIPEncapsulationReadResponse rr) {
            int transactionCounter = rr.getResponse().getTransactionCounter();
            CompletableFuture<CIPEncapsulationReadResponse> future = pendingReads.remove(transactionCounter);
            if (future != null) {
                future.complete(rr);
            } else {
                LOGGER.warn("Received CIPEncapsulationReadResponse for unknown transaction counter {}",
                    transactionCounter);
            }
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        // AB-ETH is single-flight (the legacy driver capped concurrent
        // requests at one), so process tags sequentially.
        return readSequentially(readRequest, new java.util.ArrayList<>(readRequest.getTagNames()),
            new LinkedHashMap<>());
    }

    private CompletableFuture<PlcReadResponse> readSequentially(PlcReadRequest readRequest,
                                                                List<String> remaining,
                                                                Map<String, PlcResponseItem<PlcValue>> results) {
        if (remaining.isEmpty()) {
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse(readRequest, results));
        }
        String tagName = remaining.get(0);
        List<String> next = remaining.subList(1, remaining.size());

        if (!(readRequest.getTag(tagName) instanceof AbEthTag tag)) {
            results.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            return readSequentially(readRequest, next, results);
        }

        return readSingleTag(tag, tagName)
            .handle((item, error) -> {
                if (error != null) {
                    LOGGER.warn("Read of {} failed", tagName, error);
                    PlcResponseCode code = error instanceof TimeoutException
                        ? PlcResponseCode.NOT_FOUND
                        : PlcResponseCode.INTERNAL_ERROR;
                    results.put(tagName, new DefaultPlcResponseItem<>(code, null));
                } else {
                    results.put(tagName, item);
                }
                return null;
            })
            .thenCompose(v -> readSequentially(readRequest, next, results));
    }

    private CompletableFuture<PlcResponseItem<PlcValue>> readSingleTag(AbEthTag tag, String tagName) {
        DF1RequestProtectedTypedLogicalRead logicalRead = new DF1RequestProtectedTypedLogicalRead(
            tag.getByteSize(), tag.getFileNumber(), tag.getFileType().getTypeCode(),
            tag.getElementNumber(), (short) 0);

        int transactionCounter = nextTransactionCounter();
        DF1RequestMessage requestMessage = new DF1CommandRequestMessage(
            (short) getConfiguration().getStation(), (short) 5, (short) 0,
            transactionCounter, logicalRead);
        CIPEncapsulationReadRequest read = new CIPEncapsulationReadRequest(
            sessionHandle, 0L, EMPTY_SENDER_CONTEXT, 0L, requestMessage);

        CompletableFuture<CIPEncapsulationReadResponse> future = new CompletableFuture<>();
        pendingReads.put(transactionCounter, future);
        try {
            messageCodec.send(read);
        } catch (MessageCodecException e) {
            pendingReads.remove(transactionCounter);
            return CompletableFuture.failedFuture(new PlcRuntimeException("Failed to send AB-ETH read request", e));
        }
        return future
            .orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
            .whenComplete((r, e) -> {
                if (e != null) {
                    pendingReads.remove(transactionCounter);
                }
            })
            .thenApply(response -> decodeReadResponse(response, tag, tagName));
    }

    private int nextTransactionCounter() {
        int next = transactionCounterGenerator.getAndIncrement();
        if (next >= 0xFFFF) {
            transactionCounterGenerator.set(1);
            next = transactionCounterGenerator.getAndIncrement();
        }
        return next;
    }

    /**
     * Lifted from the legacy {@code AbEthProtocolLogic#decodeReadResponse}.
     * The wire payload is a list of unsigned bytes; how to assemble them into
     * a PLC value depends entirely on the tag's file type.
     */
    private PlcResponseItem<PlcValue> decodeReadResponse(CIPEncapsulationReadResponse response,
                                                         AbEthTag tag, String tagName) {
        PlcResponseCode responseCode = decodeResponseCode(response.getResponse().getStatus());
        if (responseCode != PlcResponseCode.OK) {
            return new DefaultPlcResponseItem<>(responseCode, null);
        }
        if (!(response.getResponse() instanceof DF1CommandResponseMessageProtectedTypedLogicalRead df1PTLR)) {
            return new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null);
        }
        List<Short> data = df1PTLR.getData();
        PlcValue plcValue = null;
        try {
            switch (tag.getFileType()) {
                case INTEGER:
                    if (data.size() == 1) {
                        plcValue = new PlcINT(data.get(0));
                    } else {
                        plcValue = DefaultPlcValueHandler.of(tag, data);
                    }
                    break;
                case WORD:
                    plcValue = DefaultPlcValueHandler.of(tag, decodeSignedLE(data, 2));
                    break;
                case DWORD:
                    plcValue = DefaultPlcValueHandler.of(tag, decodeSignedLE(data, 4));
                    break;
                case SINGLEBIT:
                    int bit = tag.getBitNumber();
                    int byteIdx = bit < 8 ? 0 : 1;
                    plcValue = DefaultPlcValueHandler.of(tag,
                        (data.get(byteIdx) & (1 << (bit % 8))) != 0);
                    break;
                default:
                    LOGGER.warn("Decoding of file type {} not implemented (tag '{}'): {}",
                        tag.getFileType(), tagName, tag);
            }
        } catch (Exception e) {
            LOGGER.warn("Error decoding tag '{}': {}", tagName, tag, e);
        }
        return new DefaultPlcResponseItem<>(responseCode, plcValue);
    }

    /**
     * Reassembles {@code byteCount} little-endian unsigned bytes into a signed
     * integer. The legacy driver had two ad-hoc copies of this logic (one for
     * WORD, one for DWORD); both branches collapsed to {@code value = bytes
     * interpreted as signed integer of the matching width}.
     */
    private static int decodeSignedLE(List<Short> data, int byteCount) {
        int raw = 0;
        for (int i = byteCount - 1; i >= 0; i--) {
            raw = (raw << 8) | (data.get(i) & 0xFF);
        }
        // Sign-extend from the top bit of byteCount-1 bytes.
        int signBit = 1 << (byteCount * 8 - 1);
        if ((raw & signBit) != 0) {
            raw |= -(signBit << 1);   // fill remaining high bits with 1s
        }
        return raw;
    }

    private static PlcResponseCode decodeResponseCode(short status) {
        return status == 0 ? PlcResponseCode.OK : PlcResponseCode.NOT_FOUND;
    }

}
