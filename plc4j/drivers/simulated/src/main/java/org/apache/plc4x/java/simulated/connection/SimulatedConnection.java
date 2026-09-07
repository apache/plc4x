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
package org.apache.plc4x.java.simulated.connection;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcPingRequest;
import org.apache.plc4x.java.api.messages.PlcPingResponse;
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
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.simulated.configuration.SimulatedConfiguration;
import org.apache.plc4x.java.simulated.tag.SimulatedTag;
import org.apache.plc4x.java.simulated.tag.SimulatedTagHandler;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcPingResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionRequest;
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
import java.util.Objects;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

/**
 * Connection to a {@link SimulatedDevice}. Lives entirely in-process — there
 * is no real transport — so this class skips {@code startReceiving} and most
 * of the {@link ConnectionBase} polling machinery. The transport instance
 * passed to the superclass is {@code null}, which is safe because the only
 * code paths that touch it ({@code startReceiving} / {@code stopReceiving} /
 * {@code isAsyncTransport}) are never invoked on this connection.
 */
public class SimulatedConnection extends ConnectionBase<SimulatedConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(SimulatedConnection.class);

    private final SimulatedDevice device;

    private boolean connected = false;

    /** Consumers registered against subscription handles, keyed by handle. */
    private final Map<PlcSubscriptionHandle, PlcConsumerRegistration> registrations = new ConcurrentHashMap<>();

    /** Lookup from registration id → consumer, used by the dispatch path. */
    private final Map<Integer, Consumer<PlcSubscriptionEvent>> consumerIdMap = new ConcurrentHashMap<>();

    public SimulatedConnection(SimulatedDevice device) {
        this(device, new SimulatedConfiguration(), AuditLog.builder().build());
    }

    public SimulatedConnection(SimulatedDevice device,
                               SimulatedConfiguration configuration,
                               AuditLog auditLog) {
        super(configuration, (TransportInstance<?>) null, auditLog);
        this.device = Objects.requireNonNull(device, "device");
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new SimulatedTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    protected void onConnect() {
        connected = true;
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    @Override
    public boolean isConnected() {
        return connected;
    }

    @Override
    public void close() {
        connected = false;
        device.shutdown();
        registrations.clear();
        consumerIdMap.clear();
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    @Override
    public String toString() {
        return String.format("simulated:%s", device);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Ping
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        return CompletableFuture.completedFuture(new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read / write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        Map<String, PlcResponseItem<PlcValue>> tags = new HashMap<>();
        for (String tagName : readRequest.getTagNames()) {
            // A tag the builder couldn't parse is kept in the request with its error code and a
            // null tag, so echo that code instead of dereferencing the tag (as onWrite does).
            PlcResponseCode requestCode = readRequest.getTagResponseCode(tagName);
            if (requestCode != PlcResponseCode.OK) {
                tags.put(tagName, new DefaultPlcResponseItem<>(requestCode, null));
                continue;
            }
            SimulatedTag tag = (SimulatedTag) readRequest.getTag(tagName);
            Optional<PlcValue> value = device.get(tag);
            tags.put(tagName, value
                .map(v -> new DefaultPlcResponseItem<>(PlcResponseCode.OK, v))
                .orElseGet(() -> new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null)));
        }
        return CompletableFuture.completedFuture(new DefaultPlcReadResponse(readRequest, tags));
    }

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        Map<String, PlcResponseCode> tags = new HashMap<>();
        for (String tagName : writeRequest.getTagNames()) {
            PlcResponseCode requestCode = writeRequest.getTagResponseCode(tagName);
            if (requestCode == PlcResponseCode.OK) {
                SimulatedTag tag = (SimulatedTag) writeRequest.getTag(tagName);
                PlcValue value = writeRequest.getPlcValue(tagName);
                device.set(tag, value);
                tags.put(tagName, PlcResponseCode.OK);
            } else {
                tags.put(tagName, requestCode);
            }
        }
        return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, tags));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe / unsubscribe
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();
        for (String name : subscriptionRequest.getTagNames()) {
            // No invalid-address guard needed here: unlike the read/write builders, the
            // subscription builder rejects an unparseable address by throwing, so a request
            // that reaches this point only holds tags that parsed.
            SimulatedSubscriptionHandle handle = new SimulatedSubscriptionHandle(this, name);
            DefaultPlcSubscriptionTag subscriptionTag =
                (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(name);
            switch (subscriptionTag.getPlcSubscriptionType()) {
                case CYCLIC:
                    device.addCyclicSubscription(dispatchSubscriptionEvent(name, handle), handle, subscriptionTag,
                        subscriptionTag.getDuration().orElseThrow(
                            () -> new PlcRuntimeException("Cyclic subscription needs a polling duration")));
                    break;
                case CHANGE_OF_STATE:
                    device.addChangeOfStateSubscription(dispatchSubscriptionEvent(name, handle), handle, subscriptionTag);
                    break;
                case EVENT:
                    device.addEventSubscription(dispatchSubscriptionEvent(name, handle), handle, subscriptionTag);
                    break;
            }
            values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
        }

        // Honor setConsumer(...) / per-tag consumers on the request: SPI stores
        // them but nothing on the framework side calls register(...) for us.
        if (subscriptionRequest instanceof DefaultPlcSubscriptionRequest req) {
            Consumer<PlcSubscriptionEvent> requestConsumer = req.getConsumer();
            for (Map.Entry<String, PlcResponseItem<PlcSubscriptionHandle>> entry : values.entrySet()) {
                PlcSubscriptionHandle handle = entry.getValue().getValue();
                Consumer<PlcSubscriptionEvent> perTag = req.getTagConsumer(entry.getKey());
                if (perTag != null) {
                    handle.register(perTag);
                }
                if (requestConsumer != null) {
                    handle.register(requestConsumer);
                }
            }
        }

        return CompletableFuture.completedFuture(new DefaultPlcSubscriptionResponse(subscriptionRequest, values));
    }

    private Consumer<PlcValue> dispatchSubscriptionEvent(String name, PlcSubscriptionHandle handle) {
        return plcValue -> {
            PlcConsumerRegistration registration = registrations.get(handle);
            if (registration == null) {
                return;
            }
            Consumer<PlcSubscriptionEvent> consumer = consumerIdMap.get(registration.getConsumerId());
            if (consumer == null) {
                return;
            }
            consumer.accept(new DefaultPlcSubscriptionEvent(
                Instant.now(),
                Collections.singletonMap(name,
                    new DefaultPlcResponseItem<>(PlcResponseCode.OK, plcValue))));
        };
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        device.removeHandles(unsubscriptionRequest.getSubscriptionHandles());
        return CompletableFuture.completedFuture(new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest));
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer,
                                                        Collection<PlcSubscriptionHandle> handles) {
        PlcConsumerRegistration registration = new DefaultPlcConsumerRegistration(this, consumer,
            handles.toArray(new PlcSubscriptionHandle[0]));
        handles.forEach(handle -> registrations.put(handle, registration));
        consumerIdMap.put(registration.getConsumerId(), consumer);
        return registration;
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        var iterator = registrations.entrySet().iterator();
        while (iterator.hasNext()) {
            var entry = iterator.next();
            if (!entry.getValue().equals(registration)) {
                continue;
            }
            consumerIdMap.remove(entry.getValue().getConsumerId());
            device.removeHandles(entry.getValue().getSubscriptionHandles());
            iterator.remove();
        }
    }

}
