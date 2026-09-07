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
package org.apache.plc4x.java.canopen;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
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
import org.apache.plc4x.java.canopen.configuration.CANOpenConfiguration;
import org.apache.plc4x.java.canopen.conversation.CANConversation;
import org.apache.plc4x.java.canopen.conversation.SDODownloadConversation;
import org.apache.plc4x.java.canopen.conversation.SDOUploadConversation;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.CANOpenHeartbeatPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenNetworkPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPDO;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPDOPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.canopen.readwrite.DataItem;
import org.apache.plc4x.java.canopen.readwrite.IndexAddress;
import org.apache.plc4x.java.canopen.readwrite.NMTState;
import org.apache.plc4x.java.canopen.readwrite.NMTStateRequest;
import org.apache.plc4x.java.canopen.tag.CANOpenHeartbeatTag;
import org.apache.plc4x.java.canopen.tag.CANOpenNMTTag;
import org.apache.plc4x.java.canopen.tag.CANOpenPDOTag;
import org.apache.plc4x.java.canopen.tag.CANOpenSDOTag;
import org.apache.plc4x.java.canopen.tag.CANOpenTag;
import org.apache.plc4x.java.canopen.tag.CANOpenTagHandler;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
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
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.spi.values.PlcUSINT;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.Instant;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Map.Entry;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.Consumer;
import java.util.function.Predicate;

public class CANOpenConnection extends ConnectionBase<CANOpenConfiguration> implements CANConversation {

    private static final Logger LOGGER = LoggerFactory.getLogger(CANOpenConnection.class);

    public static final WithOption[] LITTLE_ENDIAN_OPTIONS = {
        WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    private CANOpenMessageCodec messageCodec;
    private volatile boolean connected = false;

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    /**
     * One-shot frame listeners registered by {@link #expect(Predicate, Duration)}.
     * Each is "the next frame that matches this predicate". The dispatcher walks
     * the list and gives the first match to one listener; non-matching listeners
     * stay registered for the next frame.
     */
    private final CopyOnWriteArrayList<FrameListener> listeners = new CopyOnWriteArrayList<>();

    private final ScheduledExecutorService scheduler =
        Executors.newSingleThreadScheduledExecutor(r -> {
            Thread t = new Thread(r, "canopen-scheduler");
            t.setDaemon(true);
            return t;
        });
    private ScheduledFuture<?> heartbeatTask;

    public CANOpenConnection(CANOpenConfiguration configuration,
                             TransportInstance<?> transportInstance,
                             AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new CANOpenTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        // Allow {@code PlcList} values to be passed through unwrapped, matching
        // the old driver's behavior for array writes.
        return new DefaultPlcValueHandler() {
            @Override
            public PlcValue newPlcValue(PlcTag tag, Object[] values) {
                if (values.length > 0 && values[0] instanceof PlcList) {
                    return (PlcList) values[0];
                }
                return super.newPlcValue(tag, values);
            }
        };
    }

    @Override
    public boolean isConnected() {
        return connected && messageCodec != null && messageCodec.isOpen();
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new CANOpenMessageCodec(transportInstance, this::handleIncomingFrame);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming CANopen data", e);
            }
        });
        connected = true;
        // CANopen heartbeat producer: announce BOOTED_UP once, then OPERATIONAL every 10s.
        if (configuration.isHeartbeat()) {
            try {
                messageCodec.sendFrame(heartbeatFrame(NMTState.BOOTED_UP));
            } catch (MessageCodecException e) {
                throw new PlcConnectionException("Failed to send initial CANopen heartbeat", e);
            }
            heartbeatTask = scheduler.scheduleAtFixedRate(() -> {
                try {
                    messageCodec.sendFrame(heartbeatFrame(NMTState.OPERATIONAL));
                } catch (MessageCodecException e) {
                    LOGGER.warn("Failed to send CANopen heartbeat", e);
                }
            }, 10, 10, TimeUnit.SECONDS);
        }
    }

    @Override
    public void close() throws Exception {
        connected = false;
        if (heartbeatTask != null) {
            heartbeatTask.cancel(false);
            heartbeatTask = null;
        }
        scheduler.shutdownNow();
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        consumers.clear();
        listeners.clear();
        super.close();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Conversation (used by SDOUploadConversation / SDODownloadConversation)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void sendToWire(CANOpenFrame frame) {
        try {
            messageCodec.sendFrame(frame);
        } catch (MessageCodecException e) {
            throw new RuntimeException("Failed to send CANopen frame", e);
        }
    }

    @Override
    public CompletableFuture<CANOpenFrame> expect(Predicate<CANOpenFrame> predicate, Duration timeout) {
        CompletableFuture<CANOpenFrame> future = new CompletableFuture<>();
        FrameListener listener = new FrameListener(predicate, future);
        listeners.add(listener);
        ScheduledFuture<?> timeoutHandle = scheduler.schedule(() -> {
            if (listeners.remove(listener)) {
                future.completeExceptionally(new TimeoutException("Timed out waiting for matching CANopen frame"));
            }
        }, timeout.toMillis(), TimeUnit.MILLISECONDS);
        future.whenComplete((f, t) -> {
            timeoutHandle.cancel(false);
            listeners.remove(listener);
        });
        return future;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming dispatch
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingFrame(CANOpenFrame frame) {
        // 1. Conversation listeners (SDO request/response) get first dibs.
        for (FrameListener listener : listeners) {
            if (listener.predicate.test(frame)) {
                if (listener.future.complete(frame)) {
                    return;
                }
            }
        }
        // 2. Frames we sent ourselves never need to be dispatched (would echo back from VirtualCAN-style transports).
        if (frame.getNodeId() == configuration.getNodeId()) {
            return;
        }
        // 3. Subscription dispatch.
        CANOpenService service = frame.getService();
        if (service == null) {
            return;
        }
        CANOpenPayload payload = frame.getPayload();
        if (service.getPdo() && payload instanceof CANOpenPDOPayload) {
            publishEvent(service, frame.getNodeId(), payload);
        } else if (service == CANOpenService.HEARTBEAT && payload instanceof CANOpenHeartbeatPayload) {
            publishEvent(service, frame.getNodeId(), payload);
        } else if (service == CANOpenService.NMT && payload instanceof CANOpenNetworkPayload) {
            publishEvent(service, frame.getNodeId(), payload);
        } else {
            LOGGER.debug("Decoded CANopen {} from {}, message {}", service, frame.getNodeId(), payload);
        }
    }

    private void publishEvent(CANOpenService service, int nodeId, CANOpenPayload payload) {
        for (Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            DefaultPlcConsumerRegistration registration = entry.getKey();
            Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            for (PlcSubscriptionHandle subscriptionHandle : registration.getSubscriptionHandles()) {
                if (!(subscriptionHandle instanceof CANOpenSubscriptionHandle handle)) {
                    continue;
                }
                if (!handle.matches(service, nodeId)) {
                    continue;
                }
                PlcSubscriptionEvent event = buildEvent(handle, nodeId, payload);
                if (event != null) {
                    consumer.accept(event);
                }
            }
        }
    }

    private PlcSubscriptionEvent buildEvent(CANOpenSubscriptionHandle handle, int nodeId, CANOpenPayload payload) {
        if (payload instanceof CANOpenPDOPayload pdoPayload) {
            CANOpenPDOTag tag = (CANOpenPDOTag) handle.getTag();
            byte[] data = pdoPayload.getPdo().getData();
            try {
                PlcValue value = DataItem.staticParse(
                    new org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased(data, LITTLE_ENDIAN_OPTIONS),
                    tag.getCanOpenDataType(), data.length);
                return new DefaultPlcSubscriptionEvent(Instant.now(), Collections.singletonMap(
                    handle.getName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK, value)));
            } catch (BufferException e) {
                LOGGER.warn("Could not parse PDO payload as {}: {}", tag.getCanOpenDataType(), e.getMessage());
                return new DefaultPlcSubscriptionEvent(Instant.now(), Collections.singletonMap(
                    handle.getName(), new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_DATA, new PlcNull())));
            }
        }
        if (payload instanceof CANOpenHeartbeatPayload hb) {
            NMTState state = hb.getState();
            Map<String, PlcValue> fields = new LinkedHashMap<>();
            fields.put("state", new PlcUSINT(state.getValue()));
            fields.put("node", new PlcUSINT(nodeId));
            return new DefaultPlcSubscriptionEvent(Instant.now(), Collections.singletonMap(
                handle.getName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcStruct(fields))));
        }
        if (payload instanceof CANOpenNetworkPayload net) {
            NMTStateRequest req = net.getRequest();
            Map<String, PlcValue> fields = new LinkedHashMap<>();
            fields.put("state", new PlcUSINT(req.getValue()));
            fields.put("node", new PlcUSINT(nodeId));
            return new DefaultPlcSubscriptionEvent(Instant.now(), Collections.singletonMap(
                handle.getName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcStruct(fields))));
        }
        return null;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read (SDO upload)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        if (readRequest.getTagNames().size() != 1) {
            CompletableFuture<PlcReadResponse> failed = new CompletableFuture<>();
            failed.completeExceptionally(new IllegalArgumentException("SDO requires single tag to be read"));
            return failed;
        }
        String tagName = readRequest.getTagNames().iterator().next();
        PlcTag tag = readRequest.getTags().get(0);
        if (!(tag instanceof CANOpenSDOTag sdoTag)) {
            CompletableFuture<PlcReadResponse> failed = new CompletableFuture<>();
            failed.completeExceptionally(new IllegalArgumentException("Only CANOpenSDOTag instances are supported"));
            return failed;
        }

        Duration timeout = Duration.ofMillis(Math.max(1, configuration.getRequestTimeout()));
        CompletableFuture<PlcValue> valueFuture = new CompletableFuture<>();
        SDOUploadConversation upload = new SDOUploadConversation(this,
            sdoTag.getNodeId(),
            sdoTag.getAnswerNodeId() == 0 ? sdoTag.getNodeId() : sdoTag.getAnswerNodeId(),
            new IndexAddress((int) sdoTag.getIndex(), sdoTag.getSubIndex()),
            sdoTag.getCanOpenDataType(),
            timeout);
        upload.execute(valueFuture);

        return valueFuture.handle((value, error) -> {
            Map<String, PlcResponseItem<PlcValue>> tags = new HashMap<>();
            if (error != null) {
                tags.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, new PlcNull()));
            } else {
                tags.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, value));
            }
            return new DefaultPlcReadResponse(readRequest, tags);
        });
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write (SDO download or PDO push)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        if (writeRequest.getTagNames().size() != 1) {
            CompletableFuture<PlcWriteResponse> failed = new CompletableFuture<>();
            failed.completeExceptionally(new IllegalArgumentException("You can write only one tag at the time"));
            return failed;
        }

        String tagName = writeRequest.getTagNames().iterator().next();
        PlcTag tag = writeRequest.getTags().get(0);
        if (!(tag instanceof CANOpenTag)) {
            CompletableFuture<PlcWriteResponse> failed = new CompletableFuture<>();
            failed.completeExceptionally(new IllegalArgumentException("Only CANOpenTag instances are supported"));
            return failed;
        }

        if (tag instanceof CANOpenSDOTag sdoTag) {
            return writeSdo(writeRequest, tagName, sdoTag);
        }
        if (tag instanceof CANOpenPDOTag pdoTag) {
            return writePdo(writeRequest, tagName, pdoTag);
        }

        CompletableFuture<PlcWriteResponse> failed = new CompletableFuture<>();
        failed.completeExceptionally(new IllegalArgumentException(
            "Only CANOpenSDOTag / CANOpenPDOTag instances are writable"));
        return failed;
    }

    private CompletableFuture<PlcWriteResponse> writeSdo(PlcWriteRequest writeRequest, String tagName, CANOpenSDOTag tag) {
        PlcValue value = writeRequest.getPlcValue(tagName);
        Duration timeout = Duration.ofMillis(Math.max(1, configuration.getRequestTimeout()));
        CompletableFuture<PlcResponseCode> codeFuture = new CompletableFuture<>();
        SDODownloadConversation download = new SDODownloadConversation(this,
            tag.getNodeId(),
            tag.getAnswerNodeId() == 0 ? tag.getNodeId() : tag.getAnswerNodeId(),
            new IndexAddress((int) tag.getIndex(), tag.getSubIndex()),
            value, tag.getCanOpenDataType(),
            timeout);
        download.execute(codeFuture);

        return codeFuture.handle((code, error) -> {
            PlcResponseCode result;
            if (error != null) {
                result = (error.getCause() instanceof org.apache.plc4x.java.canopen.transport.CANOpenAbortException
                    || error instanceof org.apache.plc4x.java.canopen.transport.CANOpenAbortException)
                    ? PlcResponseCode.REMOTE_ERROR : PlcResponseCode.INTERNAL_ERROR;
            } else {
                result = code;
            }
            return new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(tagName, result));
        });
    }

    private CompletableFuture<PlcWriteResponse> writePdo(PlcWriteRequest writeRequest, String tagName, CANOpenPDOTag tag) {
        PlcValue value = writeRequest.getPlcValue(tagName);
        try {
            int byteLen = DataItem.getLengthInBytes(value, tag.getCanOpenDataType(), value.getLength());
            WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[byteLen], LITTLE_ENDIAN_OPTIONS);
            DataItem.staticSerialize(buffer, value, tag.getCanOpenDataType(), value.getLength());
            CANOpenPDOPayload payload = new CANOpenPDOPayload(new CANOpenPDO(buffer.getBytes()));
            messageCodec.sendFrame(new CANOpenFrame((short) tag.getNodeId(), tag.getService(), payload));
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest,
                Collections.singletonMap(tagName, PlcResponseCode.OK)));
        } catch (BufferException | MessageCodecException e) {
            CompletableFuture<PlcWriteResponse> failed = new CompletableFuture<>();
            failed.completeExceptionally(e);
            return failed;
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe / consumer registry
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest request) {
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> answers = new LinkedHashMap<>();
        for (String key : request.getTagNames()) {
            DefaultPlcSubscriptionTag subscription = (DefaultPlcSubscriptionTag) request.getTag(key);
            if (subscription.getPlcSubscriptionType() != PlcSubscriptionType.EVENT) {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.UNSUPPORTED, null));
            } else if (subscription.getTag() instanceof CANOpenPDOTag pdoTag) {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, pdoTag)));
            } else if (subscription.getTag() instanceof CANOpenNMTTag nmtTag) {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, nmtTag)));
            } else if (subscription.getTag() instanceof CANOpenHeartbeatTag hbTag) {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.OK,
                    new CANOpenSubscriptionHandle(this, key, hbTag)));
            } else {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            }
        }
        return CompletableFuture.completedFuture(new DefaultPlcSubscriptionResponse(request, answers));
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest request) {
        // Consumer registrations are torn down via onUnregisterConsumer; nothing
        // to do on the bus side because CAN is push-only.
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

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private CANOpenFrame heartbeatFrame(NMTState state) {
        return new CANOpenFrame((short) configuration.getNodeId(),
            CANOpenService.HEARTBEAT, new CANOpenHeartbeatPayload(state));
    }

    private record FrameListener(Predicate<CANOpenFrame> predicate, CompletableFuture<CANOpenFrame> future) {
    }

}
