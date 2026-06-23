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
package org.apache.plc4x.java.plc4x;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.plc4x.config.Plc4xConfiguration;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xAuthRequest;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xAuthResponse;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xConnectRequest;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xConnectResponse;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xMessage;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xReadRequest;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xReadResponse;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xResponseCode;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xTagRequest;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xTagResponse;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xTagValueRequest;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xTagValueResponse;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xValueType;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xWriteRequest;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xWriteResponse;
import org.apache.plc4x.java.plc4x.tag.Plc4XTagHandler;
import org.apache.plc4x.java.plc4x.tag.Plc4xTag;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * PLC4X proxy connection.
 *
 * <p>On connect we send a {@code Plc4xConnectRequest} containing the remote
 * connection string we want the proxy to open on our behalf. The response
 * carries a {@code connectionId} that has to be quoted on every subsequent
 * read/write so the proxy knows which physical PLC the request is destined for.</p>
 *
 * <p>Reads/writes are correlated by {@code requestId}: each outgoing request
 * registers a {@link CompletableFuture} keyed by its tx id; the receive loop
 * looks the id up on every incoming message and completes the matching future.</p>
 */
public class Plc4xConnection extends ConnectionBase<Plc4xConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4xConnection.class);

    private Plc4xMessageCodec messageCodec;
    private volatile boolean handshakeComplete = false;
    private int connectionId;

    private final AtomicInteger txIdGenerator = new AtomicInteger(1);
    private final Map<Integer, CompletableFuture<Plc4xMessage>> pendingResponses = new ConcurrentHashMap<>();

    public Plc4xConnection(Plc4xConfiguration configuration,
                           TransportInstance<?> transportInstance,
                           AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new Plc4XTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    public boolean isConnected() {
        return handshakeComplete && messageCodec != null && messageCodec.isOpen();
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new Plc4xMessageCodec(transportInstance, this::handleIncomingMessage);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming PLC4X proxy data", e);
            }
        });

        // Authenticate first. The proxy mandates username/password auth; no operation is
        // permitted until this exchange succeeds. We never log the credentials.
        authenticate();

        // Open the underlying proxied connection.
        int requestId = txIdGenerator.getAndIncrement();
        CompletableFuture<Plc4xMessage> future = registerPending(requestId);
        try {
            messageCodec.send(new Plc4xConnectRequest(requestId, configuration.getRemoteConnectionString()));
        } catch (MessageCodecException e) {
            pendingResponses.remove(requestId);
            throw new PlcConnectionException("Failed to send proxy connect request", e);
        }

        try {
            Plc4xMessage response = future
                .orTimeout(configuration.getRequestTimeout(), TimeUnit.MILLISECONDS)
                .get();
            if (!(response instanceof Plc4xConnectResponse connectResponse)) {
                throw new PlcConnectionException("Unexpected response to proxy connect: " + response);
            }
            connectionId = connectResponse.getConnectionId();
            handshakeComplete = true;
        } catch (PlcConnectionException e) {
            throw e;
        } catch (Exception e) {
            throw new PlcConnectionException("Error establishing proxy connection", e);
        }
    }

    /**
     * Performs the mandatory username/password handshake with the proxy. Throws if the
     * server rejects the credentials or the exchange does not complete in time. Credentials
     * are taken from the connection configuration and are never logged.
     */
    private void authenticate() throws PlcConnectionException {
        int requestId = txIdGenerator.getAndIncrement();
        CompletableFuture<Plc4xMessage> future = registerPending(requestId);
        String username = configuration.getUsername();
        String password = configuration.getPassword();
        if (username == null || password == null) {
            pendingResponses.remove(requestId);
            throw new PlcConnectionException(
                "Username and password are required to connect to a PLC4X proxy server");
        }
        try {
            messageCodec.send(new Plc4xAuthRequest(requestId, username, password));
        } catch (MessageCodecException e) {
            pendingResponses.remove(requestId);
            throw new PlcConnectionException("Failed to send proxy authentication request", e);
        }
        try {
            Plc4xMessage response = future
                .orTimeout(configuration.getRequestTimeout(), TimeUnit.MILLISECONDS)
                .get();
            if (!(response instanceof Plc4xAuthResponse authResponse)) {
                throw new PlcConnectionException("Unexpected response to proxy authentication: " + response);
            }
            if (authResponse.getResponseCode() != Plc4xResponseCode.OK) {
                throw new PlcConnectionException(
                    "Authentication against PLC4X proxy server failed: " + authResponse.getResponseCode());
            }
        } catch (PlcConnectionException e) {
            throw e;
        } catch (Exception e) {
            throw new PlcConnectionException("Error during proxy authentication", e);
        }
    }

    @Override
    public void close() throws Exception {
        handshakeComplete = false;
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        pendingResponses.values().forEach(f ->
            f.completeExceptionally(new PlcRuntimeException("Connection closed")));
        pendingResponses.clear();
        super.close();
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        handshakeComplete = false;
        PlcRuntimeException ex = new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause);
        pendingResponses.values().forEach(f -> f.completeExceptionally(ex));
        pendingResponses.clear();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming dispatch
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingMessage(Plc4xMessage message) {
        CompletableFuture<Plc4xMessage> future = pendingResponses.remove(message.getRequestId());
        if (future != null) {
            future.complete(message);
        } else {
            LOGGER.debug("Received PLC4X proxy response for unknown request id {}", message.getRequestId());
        }
    }

    private CompletableFuture<Plc4xMessage> registerPending(int requestId) {
        CompletableFuture<Plc4xMessage> future = new CompletableFuture<>();
        pendingResponses.put(requestId, future);
        return future;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        List<Plc4xTagRequest> wireTags = new ArrayList<>(readRequest.getNumberOfTags());
        for (String tagName : readRequest.getTagNames()) {
            Plc4xTag plc4xTag = (Plc4xTag) readRequest.getTag(tagName);
            wireTags.add(new Plc4xTagRequest(
                new org.apache.plc4x.java.plc4x.readwrite.Plc4xTag(
                    tagName, plc4xTag.getAddressString() + ":" + plc4xTag.getPlcValueType().name())));
        }
        int requestId = txIdGenerator.getAndIncrement();
        Plc4xReadRequest wireRequest = new Plc4xReadRequest(requestId, connectionId, wireTags);

        CompletableFuture<Plc4xMessage> responseFuture = registerPending(requestId);
        try {
            messageCodec.send(wireRequest);
        } catch (MessageCodecException e) {
            pendingResponses.remove(requestId);
            return CompletableFuture.failedFuture(e);
        }

        return responseFuture
            .orTimeout(configuration.getRequestTimeout(), TimeUnit.MILLISECONDS)
            .whenComplete((r, t) -> {
                if (t != null) {
                    pendingResponses.remove(requestId);
                }
            })
            .thenApply(message -> {
                if (!(message instanceof Plc4xReadResponse readResponse)) {
                    throw new PlcRuntimeException("Unexpected response to read: " + message);
                }
                Map<String, PlcResponseItem<PlcValue>> apiResponses = new LinkedHashMap<>();
                for (Plc4xTagValueResponse wireTag : readResponse.getTags()) {
                    PlcResponseCode code = PlcResponseCode.valueOf(wireTag.getResponseCode().name());
                    apiResponses.put(wireTag.getTag().getName(),
                        new DefaultPlcResponseItem<>(code, wireTag.getValue()));
                }
                return new DefaultPlcReadResponse(readRequest, apiResponses);
            });
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        List<Plc4xTagValueRequest> wireTags = new ArrayList<>(writeRequest.getNumberOfTags());
        for (String tagName : writeRequest.getTagNames()) {
            Plc4xTag plc4xTag = (Plc4xTag) writeRequest.getTag(tagName);
            Plc4xValueType wireType = Plc4xValueType.valueOf(plc4xTag.getPlcValueType().name());
            PlcValue plcValue = writeRequest.getPlcValue(tagName);
            wireTags.add(new Plc4xTagValueRequest(
                new org.apache.plc4x.java.plc4x.readwrite.Plc4xTag(
                    tagName, plc4xTag.getAddressString() + ":" + plc4xTag.getPlcValueType().name()),
                wireType, plcValue));
        }
        int requestId = txIdGenerator.getAndIncrement();
        Plc4xWriteRequest wireRequest = new Plc4xWriteRequest(requestId, connectionId, wireTags);

        CompletableFuture<Plc4xMessage> responseFuture = registerPending(requestId);
        try {
            messageCodec.send(wireRequest);
        } catch (MessageCodecException e) {
            pendingResponses.remove(requestId);
            return CompletableFuture.failedFuture(e);
        }

        return responseFuture
            .orTimeout(configuration.getRequestTimeout(), TimeUnit.MILLISECONDS)
            .whenComplete((r, t) -> {
                if (t != null) {
                    pendingResponses.remove(requestId);
                }
            })
            .thenApply(message -> {
                if (!(message instanceof Plc4xWriteResponse writeResponse)) {
                    throw new PlcRuntimeException("Unexpected response to write: " + message);
                }
                Map<String, PlcResponseCode> apiResponses = new HashMap<>();
                for (Plc4xTagResponse wireTag : writeResponse.getTags()) {
                    Plc4xResponseCode code = wireTag.getResponseCode();
                    apiResponses.put(wireTag.getTag().getName(), PlcResponseCode.valueOf(code.name()));
                }
                return new DefaultPlcWriteResponse(writeRequest, apiResponses);
            });
    }

}
