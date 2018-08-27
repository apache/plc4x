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
package org.apache.plc4x.camel;

import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcSubscriber;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.items.SubscriptionEventItem;
import org.apache.plc4x.java.api.messages.items.SubscriptionResponseItem;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Calendar;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.function.Consumer;
import java.util.stream.Collectors;

import static org.mockito.Mockito.*;

public class MockDriver implements PlcDriver {

    public static final Logger LOGGER = LoggerFactory.getLogger(MockDriver.class);

    ExecutorService executorService = Executors.newFixedThreadPool(10);

    @Override
    public String getProtocolCode() {
        return "mock";
    }

    @Override
    public String getProtocolName() {
        return "Mock Protocol Implementation";
    }

    @Override
    public PlcConnection connect(String url) {
        // Mock a connection.
        PlcConnection plcConnectionMock = mock(PlcConnection.class, RETURNS_DEEP_STUBS);
        try {
            when(plcConnectionMock.prepareField(anyString())).thenReturn(mock(PlcField.class));
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
        when(plcConnectionMock.getWriter()).thenReturn(Optional.of(mock(PlcWriter.class, RETURNS_DEEP_STUBS)));

        // Mock a typical subscriber.
        PlcSubscriber plcSubscriber = mock(PlcSubscriber.class, RETURNS_DEEP_STUBS);
        when(plcSubscriber.subscribe(any())).thenAnswer(invocation -> {
            LOGGER.info("Received {}", invocation);
            PlcSubscriptionRequest subscriptionRequest = invocation.getArgument(0);
            List<SubscriptionResponseItem<?>> responseItems =
                subscriptionRequest.getRequestItems().stream().map(subscriptionRequestItem -> {
                    Consumer consumer = subscriptionRequestItem.getConsumer();
                    executorService.submit(() -> {
                        while (!Thread.currentThread().isInterrupted()) {
                            consumer.accept(new SubscriptionEventItem<>(null, Calendar.getInstance(), Collections.singletonList("HelloWorld")));
                            try {
                                TimeUnit.MILLISECONDS.sleep(100);
                            } catch (InterruptedException e) {
                                Thread.currentThread().interrupt();
                                throw new RuntimeException(e);
                            }
                        }
                    });
                    return new SubscriptionResponseItem<>(subscriptionRequestItem,
                        mock(PlcSubscriptionHandle.class, RETURNS_DEEP_STUBS), PlcResponseCode.OK);
                }).collect(Collectors.toList());
            PlcSubscriptionResponse response = new PlcSubscriptionResponse(subscriptionRequest, responseItems);
            CompletableFuture<PlcSubscriptionResponse> responseFuture = new CompletableFuture<>();
            responseFuture.complete(response);
            return responseFuture;
        });
        when(plcConnectionMock.getSubscriber()).thenReturn(Optional.of(plcSubscriber));
        return plcConnectionMock;
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) {
        return connect(null);
    }

}
