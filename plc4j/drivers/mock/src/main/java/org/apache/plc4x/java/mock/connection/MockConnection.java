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
package org.apache.plc4x.java.mock.connection;

import org.apache.commons.lang3.Validate;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseRequestInterceptor;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
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
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.mock.configuration.MockConfiguration;
import org.apache.plc4x.java.mock.tag.MockTagHandler;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcPingResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collection;
import java.util.Collections;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.stream.Collectors;

/**
 * In-process mock connection. All read/write/subscribe operations delegate to
 * a user-supplied {@link MockDevice} so tests can assert on the calls.
 *
 * <p>No transport — passes {@code null} to {@link ConnectionBase}, which is
 * safe because the mock connection never calls {@code startReceiving} /
 * {@code isAsyncTransport}.</p>
 */
public class MockConnection extends ConnectionBase<MockConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(MockConnection.class);

    private final PlcAuthentication authentication;
    private final MockTagHandler tagHandler = new MockTagHandler();
    private volatile MockDevice device;

    public MockConnection(PlcAuthentication authentication) {
        this(authentication, new MockConfiguration(),
            AuditLog.builder().withSource("mock").build());
    }

    public MockConnection(PlcAuthentication authentication,
                          MockConfiguration configuration,
                          AuditLog auditLog) {
        super(configuration, (TransportInstance<?>) null, auditLog);
        this.authentication = authentication;
    }

    public MockDevice getDevice() {
        return device;
    }

    public void setDevice(MockDevice device) {
        LOGGER.info("Set Mock Device on Mock Connection {} with device {}", this, device);
        this.device = device;
    }

    public PlcAuthentication getAuthentication() {
        return authentication;
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return tagHandler;
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    protected void onConnect() {
        // No-op — the mock connection is "connected" as soon as a device is set.
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
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        return CompletableFuture.completedFuture(new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK));
    }

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending read request to MockDevice");
            Map<String, PlcResponseItem<PlcValue>> response = readRequest.getTagNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    name -> device.read(readRequest.getTag(name).getAddressString())));
            return new DefaultPlcReadResponse(readRequest, response);
        });
    }

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending write request to MockDevice");
            Map<String, PlcResponseCode> response = writeRequest.getTagNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    name -> device.write(writeRequest.getTag(name).getAddressString(),
                        writeRequest.getPlcValue(name))));
            return new DefaultPlcWriteResponse(writeRequest, response);
        });
    }

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending subscribe request to MockDevice");
            Map<String, PlcResponseItem<PlcSubscriptionHandle>> response =
                subscriptionRequest.getTagNames().stream()
                    .collect(Collectors.toMap(
                        Function.identity(),
                        name -> device.subscribe(subscriptionRequest.getTag(name).getAddressString())));
            return new DefaultPlcSubscriptionResponse(subscriptionRequest, response);
        });
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending unsubscribe request to MockDevice");
            device.unsubscribe();
            return new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest);
        });
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer,
                                                        Collection<PlcSubscriptionHandle> handles) {
        Validate.notNull(device, "No device is set in the mock connection!");
        return device.register(consumer, handles);
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        Validate.notNull(device, "No device is set in the mock connection!");
        device.unregister(registration);
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Validate.notNull(device, "No device is set in the mock connection!");
            LOGGER.debug("Sending browse request to MockDevice");
            return new DefaultPlcBrowseResponse(browseRequest, Collections.emptyMap(), Collections.emptyMap());
        });
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowseWithInterceptor(PlcBrowseRequest browseRequest,
                                                                          PlcBrowseRequestInterceptor interceptor) {
        return onBrowse(browseRequest);
    }

}
