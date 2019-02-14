/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.simulated.connection;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.model.*;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

/**
 * Connection to a test device.
 * This class is not thread-safe.
 */
public class SimulatedPlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter, PlcSubscriber {

    private final TestDevice device;

    private boolean connected = false;

    private Map<InternalPlcSubscriptionHandle, InternalPlcConsumerRegistration> registrations = new ConcurrentHashMap<>();

    private Map<Integer, Consumer<PlcSubscriptionEvent>> consumerIdMap = new ConcurrentHashMap<>();

    public SimulatedPlcConnection(TestDevice device) {
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
    public boolean canRead() {
        return true;
    }

    @Override
    public boolean canWrite() {
        return true;
    }

    @Override
    public boolean canSubscribe() {
        return true;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new TestFieldHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return new DefaultPlcWriteRequest.Builder(this, new TestFieldHandler());
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        return new DefaultPlcSubscriptionRequest.Builder(this, new TestFieldHandler());
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        InternalPlcReadRequest request = checkInternal(readRequest, InternalPlcReadRequest.class);
        Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> fields = new HashMap<>();
        for (String fieldName : request.getFieldNames()) {
            TestField field = (TestField) request.getField(fieldName);
            Optional<BaseDefaultFieldItem> fieldItemOptional = device.get(field);
            ImmutablePair<PlcResponseCode, BaseDefaultFieldItem> fieldPair;
            boolean present = fieldItemOptional.isPresent();
            fieldPair = present
                ? new ImmutablePair<>(PlcResponseCode.OK, fieldItemOptional.get())
                : new ImmutablePair<>(PlcResponseCode.NOT_FOUND, null);
            fields.put(fieldName, fieldPair);
        }
        PlcReadResponse response = new DefaultPlcReadResponse(request, fields);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        InternalPlcWriteRequest request = checkInternal(writeRequest, InternalPlcWriteRequest.class);
        Map<String, PlcResponseCode> fields = new HashMap<>();
        for (String fieldName : request.getFieldNames()) {
            TestField field = (TestField) request.getField(fieldName);
            BaseDefaultFieldItem fieldItem = request.getFieldItem(fieldName);
            device.set(field, fieldItem);
            fields.put(fieldName, PlcResponseCode.OK);
        }
        PlcWriteResponse response = new DefaultPlcWriteResponse(request, fields);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public String toString() {
        return String.format("test:%s", device);
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        InternalPlcSubscriptionRequest request = checkInternal(subscriptionRequest, InternalPlcSubscriptionRequest.class);
        LinkedHashMap<String, SubscriptionPlcField> subscriptionPlcFieldMap = request.getSubscriptionPlcFieldMap();

        Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> values = new HashMap<>();
        subscriptionPlcFieldMap.forEach((name, subscriptionPlcField) -> {
            InternalPlcSubscriptionHandle handle = new DefaultPlcSubscriptionHandle(this);
            switch (subscriptionPlcField.getPlcSubscriptionType()) {
                case CYCLIC:
                    device.addCyclicSubscription(dispatchSubscriptionEvent(name, handle), handle, (TestField) subscriptionPlcField.getPlcField(), subscriptionPlcField.getDuration().orElseThrow(RuntimeException::new));
                    break;
                case CHANGE_OF_STATE:
                    device.addChangeOfStateSubscription(dispatchSubscriptionEvent(name, handle), handle, (TestField) subscriptionPlcField.getPlcField());
                    break;
                case EVENT:
                    device.addEventSubscription(dispatchSubscriptionEvent(name, handle), handle, (TestField) subscriptionPlcField.getPlcField());
                    break;
            }
            values.put(name, Pair.of(PlcResponseCode.OK, handle));
        });

        PlcSubscriptionResponse response = new DefaultPlcSubscriptionResponse(request, values);
        return CompletableFuture.completedFuture(response);
    }

    private Consumer<BaseDefaultFieldItem> dispatchSubscriptionEvent(String name, InternalPlcSubscriptionHandle handle) {
        return fieldItem -> {
            InternalPlcConsumerRegistration plcConsumerRegistration = registrations.get(handle);
            if (plcConsumerRegistration == null) {
                return;
            }
            int consumerHash = plcConsumerRegistration.getConsumerHash();
            Consumer<PlcSubscriptionEvent> consumer = consumerIdMap.get(consumerHash);
            if (consumer == null) {
                return;
            }
            consumer.accept(new DefaultPlcSubscriptionEvent(Instant.now(), Collections.singletonMap(name, Pair.of(PlcResponseCode.OK, fieldItem))));
        };
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        InternalPlcUnsubscriptionRequest request = checkInternal(unsubscriptionRequest, InternalPlcUnsubscriptionRequest.class);

        device.removeHandles(request.getInternalPlcSubscriptionHandles());

        PlcUnsubscriptionResponse response = new DefaultPlcUnsubscriptionResponse(request);
        return CompletableFuture.completedFuture(response);
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        InternalPlcConsumerRegistration plcConsumerRegistration = new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new InternalPlcSubscriptionHandle[0]));
        handles.stream()
            .map(InternalPlcSubscriptionHandle.class::cast)
            .forEach(handle -> registrations.put(handle, plcConsumerRegistration));
        consumerIdMap.put(plcConsumerRegistration.getConsumerHash(), consumer);
        return plcConsumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        Iterator<Map.Entry<InternalPlcSubscriptionHandle, InternalPlcConsumerRegistration>> entryIterator = registrations.entrySet().iterator();
        while (entryIterator.hasNext()) {
            Map.Entry<InternalPlcSubscriptionHandle, InternalPlcConsumerRegistration> entry = entryIterator.next();
            if (!entry.getValue().equals(registration)) {
                continue;
            }
            InternalPlcConsumerRegistration value = entry.getValue();
            int consumerHash = value.getConsumerHash();
            consumerIdMap.remove(consumerHash);
            device.removeHandles(value.getAssociatedHandles());
            entryIterator.remove();
        }
    }
}
