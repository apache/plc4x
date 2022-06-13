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
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.simulated.field.SimulatedField;
import org.apache.plc4x.java.simulated.field.SimulatedFieldHandler;
import org.apache.plc4x.java.spi.connection.AbstractPlcConnection;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.PlcWriter;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

/**
 * Connection to a test device.
 * This class is not thread-safe.
 */
public class SimulatedConnection extends AbstractPlcConnection implements PlcReader, PlcWriter, PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(SimulatedConnection.class);

    private final SimulatedDevice device;

    private boolean connected = false;

    private final Map<PlcSubscriptionHandle, PlcConsumerRegistration> registrations = new ConcurrentHashMap<>();

    private final Map<Integer, Consumer<PlcSubscriptionEvent>> consumerIdMap = new ConcurrentHashMap<>();

    public SimulatedConnection(SimulatedDevice device) {
        super(true, true, true, new SimulatedFieldHandler(), new IEC61131ValueHandler(), null);
        this.device = device;
    }

    @Override
    public void connect() {
        connected = true;
    }

    @Override
    public boolean isConnected() {
        return connected;
    }

    @Override
    public void close() {
        connected = false;
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        Map<String, ResponseItem<PlcValue>> fields = new HashMap<>();
        for (String fieldName : readRequest.getFieldNames()) {
            SimulatedField field = (SimulatedField) readRequest.getField(fieldName);
            Optional<PlcValue> valueOptional = device.get(field);
            ResponseItem<PlcValue> fieldPair;
            boolean present = valueOptional.isPresent();
            fieldPair = present
                ? new ResponseItem<>(PlcResponseCode.OK, valueOptional.get())
                : new ResponseItem<>(PlcResponseCode.NOT_FOUND, null);
            fields.put(fieldName, fieldPair);
        }
        PlcReadResponse response = new DefaultPlcReadResponse(readRequest, fields);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        Map<String, PlcResponseCode> fields = new HashMap<>();
        for (String fieldName : writeRequest.getFieldNames()) {
            SimulatedField field = (SimulatedField) writeRequest.getField(fieldName);
            PlcValue value = writeRequest.getPlcValue(fieldName);
            device.set(field, value);
            fields.put(fieldName, PlcResponseCode.OK);
        }
        PlcWriteResponse response = new DefaultPlcWriteResponse(writeRequest, fields);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public String toString() {
        return String.format("simulated:%s", device);
    }

    /**
     * Blocking subscribe call
     *
     * @param subscriptionRequest subscription request containing at least one subscription request item.
     * @return the {@code PlcSubscriptionResponse}
     */
    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        LOGGER.info("subscribing {}", subscriptionRequest);
        Map<String, ResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
        subscriptionRequest.getFieldNames().forEach(name -> {
            LOGGER.info("creating handle for field name {}", name);
            PlcSubscriptionHandle handle = new DefaultPlcSubscriptionHandle(this);
            final PlcSubscriptionField subscriptionPlcField = subscriptionRequest.getField(name);
            switch (subscriptionPlcField.getPlcSubscriptionType()) {
                case CYCLIC:
                    LOGGER.info("Adding cyclic subscription for field name {}", name);
                    device.addCyclicSubscription(dispatchSubscriptionEvent(name, handle), handle, subscriptionPlcField, subscriptionPlcField.getDuration().orElseThrow(RuntimeException::new));
                    break;
                case CHANGE_OF_STATE:
                    LOGGER.info("Adding change of state subscription for field name {}", name);
                    device.addChangeOfStateSubscription(dispatchSubscriptionEvent(name, handle), handle, subscriptionPlcField);
                    break;
                case EVENT:
                    LOGGER.info("Adding event subscription for field name {}", name);
                    device.addEventSubscription(dispatchSubscriptionEvent(name, handle), handle, subscriptionPlcField);
                    break;
            }
            values.put(name, new ResponseItem<>(PlcResponseCode.OK, handle));
        });

        PlcSubscriptionResponse response = new DefaultPlcSubscriptionResponse(subscriptionRequest, values);
        return CompletableFuture.completedFuture(response);
    }

    private Consumer<PlcValue> dispatchSubscriptionEvent(String name, PlcSubscriptionHandle handle) {
        return plcValue -> {
            LOGGER.info("handling plc value {}", plcValue);
            PlcConsumerRegistration plcConsumerRegistration = registrations.get(handle);
            if (plcConsumerRegistration == null) {
                LOGGER.warn("no registration for handle {}", handle);
                return;
            }
            int consumerId = plcConsumerRegistration.getConsumerId();
            Consumer<PlcSubscriptionEvent> consumer = consumerIdMap.get(consumerId);
            if (consumer == null) {
                LOGGER.warn("no consumer for id {}", consumerId);
                return;
            }
            consumer.accept(
                new DefaultPlcSubscriptionEvent(
                    Instant.now(),
                    Collections.singletonMap(name, new ResponseItem<>(PlcResponseCode.OK, plcValue))
                )
            );
        };
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        LOGGER.info("unsubscribing {}", unsubscriptionRequest);
        device.removeHandles(unsubscriptionRequest.getSubscriptionHandles());

        PlcUnsubscriptionResponse response = new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        LOGGER.info("Registering consumer {} with handles {}", consumer, handles);
        PlcConsumerRegistration plcConsumerRegistration = new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        handles.stream()
            .map(PlcSubscriptionHandle.class::cast)
            .forEach(handle -> registrations.put(handle, plcConsumerRegistration));
        consumerIdMap.put(plcConsumerRegistration.getConsumerId(), consumer);
        return plcConsumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        LOGGER.info("Unregistering {}", registration);
        Iterator<Map.Entry<PlcSubscriptionHandle, PlcConsumerRegistration>> entryIterator = registrations.entrySet().iterator();
        while (entryIterator.hasNext()) {
            Map.Entry<PlcSubscriptionHandle, PlcConsumerRegistration> entry = entryIterator.next();
            if (!entry.getValue().equals(registration)) {
                LOGGER.debug("not the value we looking for {}. We are looking for {}", entry.getValue(), registration);
                continue;
            }
            PlcConsumerRegistration consumerRegistration = entry.getValue();
            int consumerId = consumerRegistration.getConsumerId();
            LOGGER.info("Removing consumer {}", consumerId);
            consumerIdMap.remove(consumerId);
            LOGGER.info("Removing handles {}", consumerRegistration.getSubscriptionHandles());
            device.removeHandles(consumerRegistration.getSubscriptionHandles());
            entryIterator.remove();
        }
    }
}
