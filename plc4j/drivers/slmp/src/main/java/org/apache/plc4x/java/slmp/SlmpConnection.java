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
package org.apache.plc4x.java.slmp;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.slmp.config.SlmpConfiguration;
import org.apache.plc4x.java.slmp.readwrite.SlmpMessage;
import org.apache.plc4x.java.slmp.readwrite.SlmpReadRequest;
import org.apache.plc4x.java.slmp.readwrite.SlmpRequestFrame3E;
import org.apache.plc4x.java.slmp.readwrite.SlmpResponseFrame3E;
import org.apache.plc4x.java.slmp.tag.SlmpTag;
import org.apache.plc4x.java.slmp.tag.SlmpTagHandler;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

public class SlmpConnection extends ConnectionBase<SlmpConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(SlmpConnection.class);

    private SlmpMessageCodec messageCodec;
    private volatile CompletableFuture<SlmpResponseFrame3E> pendingResponse;

    public SlmpConnection(SlmpConfiguration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new SlmpMessageCodec(transportInstance, this::handleIncomingMessage);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming SLMP data", e);
            }
        });
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
        LOGGER.info("SLMP connection established");
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
        failPending(new PlcRuntimeException("Connection closed"));
        super.close();
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");
        failPending(new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause));
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new SlmpTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    protected int getMaxConcurrentRequests() {
        return 1;
    }

    private void handleIncomingMessage(SlmpMessage message) {
        CompletableFuture<SlmpResponseFrame3E> future = pendingResponse;
        pendingResponse = null;
        if (future == null) {
            LOGGER.warn("Received unsolicited SLMP message");
            return;
        }
        if (message instanceof SlmpResponseFrame3E response) {
            future.complete(response);
        } else {
            future.completeExceptionally(new PlcRuntimeException("Unexpected SLMP message type: " + message));
        }
    }

    private void failPending(Throwable cause) {
        CompletableFuture<SlmpResponseFrame3E> future = pendingResponse;
        pendingResponse = null;
        if (future != null) {
            future.completeExceptionally(cause);
        }
    }

    private CompletableFuture<SlmpResponseFrame3E> sendRequest(SlmpRequestFrame3E request) {
        CompletableFuture<SlmpResponseFrame3E> responseFuture = new CompletableFuture<>();
        pendingResponse = responseFuture;
        try {
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingResponse = null;
            responseFuture.completeExceptionally(new PlcRuntimeException("Failed to send SLMP request", e));
            return responseFuture;
        }
        responseFuture.orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
            .whenComplete((r, e) -> {
                if (e instanceof TimeoutException) {
                    pendingResponse = null;
                }
            });
        return responseFuture;
    }

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;

        LinkedHashMap<String, CompletableFuture<PlcResponseItem<PlcValue>>> tagFutures = new LinkedHashMap<>();
        CompletableFuture<Void> chain = CompletableFuture.completedFuture(null);
        for (String tagName : request.getTagNames()) {
            SlmpTag tag = (SlmpTag) request.getTag(tagName);
            CompletableFuture<PlcResponseItem<PlcValue>> tagFuture =
                chain.thenComposeAsync(v -> readSingleTag(tag));
            tagFutures.put(tagName, tagFuture);
            chain = tagFuture.handle((r, e) -> null);
        }

        CompletableFuture<Void> allDone =
            CompletableFuture.allOf(tagFutures.values().toArray(new CompletableFuture[0]));
        return allDone.thenApply(v -> {
            Map<String, PlcResponseItem<PlcValue>> items = new LinkedHashMap<>();
            for (Map.Entry<String, CompletableFuture<PlcResponseItem<PlcValue>>> e : tagFutures.entrySet()) {
                try {
                    items.put(e.getKey(), e.getValue().join());
                } catch (Exception ex) {
                    items.put(e.getKey(), new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
            }
            return (PlcReadResponse) new DefaultPlcReadResponse(request, items);
        });
    }

    private CompletableFuture<PlcResponseItem<PlcValue>> readSingleTag(SlmpTag tag) {
        SlmpReadRequest data = new SlmpReadRequest(
            tag.getDeviceNumber(), tag.getDeviceCode(), tag.getNumberOfPoints());
        SlmpRequestFrame3E frame = new SlmpRequestFrame3E(
            getConfiguration().getMonitoringTimer(), 0x0401, 0x0000, data);
        return executeThrottled(() ->
            sendRequest(frame).thenApply(response ->
                SlmpResponseMapper.mapTag(tag, response.getEndCode(), response.getResponseData())));
    }
}
