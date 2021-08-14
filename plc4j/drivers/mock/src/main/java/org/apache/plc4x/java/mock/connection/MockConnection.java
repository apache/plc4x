/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.mock.connection;

import org.apache.commons.lang3.Validate;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.mock.field.MockField;
import org.apache.plc4x.java.mock.field.MockFieldHandler;
import org.apache.plc4x.java.mock.field.MockValueHandler;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcUnsubscriptionRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.PlcWriter;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collection;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.stream.Collectors;

public class MockConnection implements PlcConnection, PlcReader, PlcWriter, PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(MockConnection.class);

    private final PlcAuthentication authentication;

    private MockDevice device;

    public MockConnection(PlcAuthentication authentication) {
        this.authentication = authentication;
    }

    public MockDevice getDevice() {
        return device;
    }

    public void setDevice(MockDevice device) {
        LOGGER.info("Set Mock Device on Mock Connection {} with device {}", this, device);
        this.device = device;
    }

    @Override
    public void connect() {
        // do nothing
    }

    @Override
    public CompletableFuture<Void> ping() {
        CompletableFuture<Void> future = new CompletableFuture<>();
        future.completeExceptionally(new UnsupportedOperationException());
        return future;
    }

    @Override
    public boolean isConnected() {
        return device != null;
    }

    @Override
    public void close() {
        LOGGER.info("Closing MockConnection with device {}", device);
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        return new PlcConnectionMetadata() {
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
        };
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new MockFieldHandler());
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending read request to MockDevice");
            Map<String, ResponseItem<PlcValue>> response = readRequest.getFieldNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    name -> device.read(((MockField) readRequest.getField(name)).getAddress())
                    )
                );
            return new DefaultPlcReadResponse((DefaultPlcReadRequest) readRequest, response);
        });
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending write request to MockDevice");
            Map<String, PlcResponseCode> response = writeRequest.getFieldNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    name -> device.write(((MockField) writeRequest.getField(name)).getAddress(),
                        ((MockField) writeRequest.getField(name)).getPlcValue())
                    )
                );
            return new DefaultPlcWriteResponse((DefaultPlcWriteRequest) writeRequest, response);
        });
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending subsribe request to MockDevice");
            Map<String, ResponseItem<PlcSubscriptionHandle>> response = subscriptionRequest.getFieldNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    name -> device.subscribe(((MockField) subscriptionRequest.getField(name)).getAddress())
                    )
                );
            return new DefaultPlcSubscriptionResponse((DefaultPlcSubscriptionRequest) subscriptionRequest, response);
        });
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending subsribe request to MockDevice");
            device.unsubscribe();
            return new DefaultPlcUnsubscriptionResponse((DefaultPlcUnsubscriptionRequest) unsubscriptionRequest);
        });
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        return device.register(consumer, handles);
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        device.unregister(registration);
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return new DefaultPlcWriteRequest.Builder(this, new MockFieldHandler(), new MockValueHandler());
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        return new DefaultPlcSubscriptionRequest.Builder(this, new MockFieldHandler());
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
    }

    public PlcAuthentication getAuthentication() {
        return authentication;
    }

}
