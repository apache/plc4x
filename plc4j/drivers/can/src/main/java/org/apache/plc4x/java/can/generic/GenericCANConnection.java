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
package org.apache.plc4x.java.can.generic;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.can.generic.configuration.GenericCANConfiguration;
import org.apache.plc4x.java.can.generic.protocol.GenericCANSubscriptionHandle;
import org.apache.plc4x.java.can.generic.tag.GenericCANTag;
import org.apache.plc4x.java.can.generic.tag.GenericCANTagHandler;
import org.apache.plc4x.java.can.generic.transport.GenericFrame;
import org.apache.plc4x.java.genericcan.readwrite.DataItem;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcUnsubscriptionResponse;
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

import java.time.Instant;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Map.Entry;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

/**
 * Generic CAN connection. Direct port of the legacy {@code GenericCANProtocolLogic} +
 * {@code CANDriverAdapter} stack to the SPI3 {@link ConnectionBase} model.
 *
 * <p>The adapter pattern from the old SPI (Plc4xCANProtocolBase + ConversationContext
 * wrappers) is collapsed into this class: incoming frames come straight off the codec,
 * outgoing writes go straight to it. Each subscription holds the list of tags it cares
 * about for a particular CAN node id; when a matching frame arrives we decode each tag
 * out of the payload and deliver one {@code DefaultPlcSubscriptionEvent} to every
 * registered consumer.</p>
 */
public class GenericCANConnection extends ConnectionBase<GenericCANConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(GenericCANConnection.class);

    static final WithOption[] LITTLE_ENDIAN_OPTIONS = {
        WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    private GenericCANMessageCodec messageCodec;
    private volatile boolean connected = false;

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    public GenericCANConnection(GenericCANConfiguration configuration,
                                TransportInstance<?> transportInstance,
                                AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new GenericCANTagHandler();
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
    protected void onConnect() {
        messageCodec = new GenericCANMessageCodec(transportInstance, this::handleIncomingFrame);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming CAN data", e);
            }
        });
        connected = true;
    }

    @Override
    public void close() throws Exception {
        connected = false;
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        consumers.clear();
        super.close();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming dispatch
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingFrame(GenericFrame frame) {
        for (Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            DefaultPlcConsumerRegistration registration = entry.getKey();
            Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
                if (!(handle instanceof GenericCANSubscriptionHandle subscription)) {
                    continue;
                }
                if (!subscription.matches(frame.getNodeId())) {
                    continue;
                }
                Map<String, PlcResponseItem<PlcValue>> tags = new LinkedHashMap<>();
                ReadBufferByteBased readBuffer = new ReadBufferByteBased(frame.getData(), LITTLE_ENDIAN_OPTIONS);
                for (Entry<String, GenericCANTag> tag : subscription.getTags().entrySet()) {
                    try {
                        PlcValue value = DataItem.staticParse(readBuffer, tag.getValue().getDataType(), frame.getData().length);
                        tags.put(tag.getKey(), new DefaultPlcResponseItem<>(
                            value != null ? PlcResponseCode.OK : PlcResponseCode.INTERNAL_ERROR, value));
                    } catch (BufferException e) {
                        tags.put(tag.getKey(), new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_DATA, null));
                    }
                }
                consumer.accept(new DefaultPlcSubscriptionEvent(Instant.now(), tags));
            }
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        return CompletableFuture.supplyAsync(() -> {
            // Pack all writes for a given CAN node id into one outgoing frame
            // — that's what the legacy driver did, and what makes sense given
            // each CAN frame can carry up to 8 bytes of payload.
            Map<Integer, WriteBufferByteBased> messages = new LinkedHashMap<>();
            Map<Integer, Map<String, PlcResponseCode>> responses = new HashMap<>();

            for (String tagName : writeRequest.getTagNames()) {
                PlcTag plcTag = writeRequest.getTag(tagName);
                if (!(plcTag instanceof GenericCANTag canTag)) {
                    responses.computeIfAbsent(-1, n -> new HashMap<>()).put(tagName, PlcResponseCode.UNSUPPORTED);
                    continue;
                }
                WriteBufferByteBased buffer = messages.computeIfAbsent(canTag.getNodeId(),
                    n -> new WriteBufferByteBased(new byte[8], LITTLE_ENDIAN_OPTIONS));
                Map<String, PlcResponseCode> statusMap = responses.computeIfAbsent(canTag.getNodeId(),
                    n -> new HashMap<>());
                PlcValue value = writeRequest.getPlcValue(tagName);
                try {
                    DataItem.staticSerialize(buffer, value, canTag.getDataType(), value.getLength());
                    statusMap.put(tagName, PlcResponseCode.OK);
                } catch (BufferException e) {
                    statusMap.put(tagName, PlcResponseCode.INVALID_DATA);
                }
            }

            Map<String, PlcResponseCode> codes = new LinkedHashMap<>();
            for (Entry<Integer, WriteBufferByteBased> message : messages.entrySet()) {
                int nodeId = message.getKey();
                Map<String, PlcResponseCode> nodeResponses = responses.get(nodeId);
                boolean discarded = false;
                for (Entry<String, PlcResponseCode> entry : nodeResponses.entrySet()) {
                    codes.put(entry.getKey(), entry.getValue());
                    if (!discarded && entry.getValue() != PlcResponseCode.OK) {
                        LOGGER.info("Discarding frame for node {}: a tag failed to serialize ({})",
                            nodeId, entry.getKey());
                        discarded = true;
                    }
                }
                if (!discarded) {
                    byte[] data = message.getValue().getBytes();
                    try {
                        messageCodec.sendFrame(new GenericFrame(nodeId, data));
                    } catch (MessageCodecException e) {
                        LOGGER.error("Failed to send CAN frame for node {}", nodeId, e);
                        for (String tagName : nodeResponses.keySet()) {
                            codes.put(tagName, PlcResponseCode.REMOTE_ERROR);
                        }
                    }
                }
            }
            // Tags with no node (UNSUPPORTED) need to land in the response too.
            Map<String, PlcResponseCode> unsupported = responses.get(-1);
            if (unsupported != null) {
                codes.putAll(unsupported);
            }

            return new DefaultPlcWriteResponse(writeRequest, codes);
        });
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe / consumer registry
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest request) {
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> answers = new LinkedHashMap<>();
        Map<Integer, GenericCANSubscriptionHandle> handles = new LinkedHashMap<>();
        for (String key : request.getTagNames()) {
            DefaultPlcSubscriptionTag subscription = (DefaultPlcSubscriptionTag) request.getTag(key);
            if (subscription.getPlcSubscriptionType() != PlcSubscriptionType.EVENT) {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.UNSUPPORTED, null));
            } else if (subscription.getTag() instanceof GenericCANTag canTag) {
                GenericCANSubscriptionHandle handle = handles.computeIfAbsent(canTag.getNodeId(),
                    node -> new GenericCANSubscriptionHandle(this, node));
                handle.add(key, canTag);
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
            } else {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            }
        }
        return CompletableFuture.completedFuture(new DefaultPlcSubscriptionResponse(request, answers));
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest request) {
        // CAN is push-only — there's nothing to unsubscribe on the bus side.
        // Consumer cleanup happens via onUnregisterConsumer.
        return CompletableFuture.completedFuture(new DefaultPlcUnsubscriptionResponse(request));
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer,
                                                        Collection<PlcSubscriptionHandle> handles) {
        DefaultPlcConsumerRegistration registration = new DefaultPlcConsumerRegistration(this, consumer,
            handles.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(registration, consumer);
        return registration;
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        if (registration instanceof DefaultPlcConsumerRegistration r) {
            consumers.remove(r);
        }
    }

}
