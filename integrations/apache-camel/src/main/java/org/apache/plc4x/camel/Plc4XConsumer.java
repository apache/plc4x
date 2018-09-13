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

import org.apache.camel.*;
import org.apache.camel.spi.ExceptionHandler;
import org.apache.camel.support.LoggingExceptionHandler;
import org.apache.camel.support.ServiceSupport;
import org.apache.camel.util.AsyncProcessorConverterHelper;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcSubscriber;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.temporal.ChronoUnit;
import java.util.Collection;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

public class Plc4XConsumer extends ServiceSupport implements Consumer, java.util.function.Consumer<PlcSubscriptionEvent> {
    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XConsumer.class);

    private Plc4XEndpoint endpoint;
    private AsyncProcessor processor;
    private ExceptionHandler exceptionHandler;
    private PlcConnection plcConnection;
    private String fieldQuery;
    private Class<?> dataType;
    private PlcSubscriptionResponse subscriptionResponse;


    public Plc4XConsumer(Plc4XEndpoint endpoint, Processor processor) throws PlcException {
        this.endpoint = endpoint;
        this.dataType = endpoint.getDataType();
        this.processor = AsyncProcessorConverterHelper.convert(processor);
        this.exceptionHandler = new LoggingExceptionHandler(endpoint.getCamelContext(), getClass());
        String plc4xURI = endpoint.getEndpointUri().replaceFirst("plc4x:/?/?", "");
        this.plcConnection = endpoint.getPlcDriverManager().getConnection(plc4xURI);
        this.fieldQuery = endpoint.getAddress();
    }

    @Override
    public String toString() {
        return "Plc4XConsumer[" + endpoint + "]";
    }

    @Override
    public Endpoint getEndpoint() {
        return endpoint;
    }

    public ExceptionHandler getExceptionHandler() {
        return exceptionHandler;
    }

    public void setExceptionHandler(ExceptionHandler exceptionHandler) {
        this.exceptionHandler = exceptionHandler;
    }

    @Override
    protected void doStart() throws InterruptedException, ExecutionException, PlcException {
        PlcSubscriber plcSubscriber = plcConnection.getSubscriber().orElseThrow(
            () -> new PlcException("Connection doesn't support subscriptions."));
        // TODO: Is it correct to only support one field?
        PlcSubscriptionRequest request = plcSubscriber.subscriptionRequestBuilder()
            .addCyclicField("default", fieldQuery, Duration.of(3, ChronoUnit.SECONDS)).build();
        plcSubscriber.register(request, this);
    }

    @Override
    protected void doStop() throws InterruptedException, ExecutionException, TimeoutException, PlcException {
        PlcSubscriber plcSubscriber = plcConnection.getSubscriber().orElseThrow(
            () -> new PlcException("Connection doesn't support subscriptions."));
        CompletableFuture<PlcUnsubscriptionResponse> unsubscriptionFuture = plcSubscriber.unsubscribe(builder -> builder.addHandles(subscriptionResponse.getSubscriptionHandles()));
        PlcUnsubscriptionResponse unsubscriptionResponse = unsubscriptionFuture.get(5, TimeUnit.SECONDS);
        // TODO: Handle the response ...
        try {
            plcConnection.close();
        } catch (Exception e) {
            LOGGER.error("Error closing connection", e);
        }
    }

    private PlcSubscriber getSubscriber() {
        return plcConnection.getSubscriber().orElseThrow(() -> new RuntimeException("No subscriber available"));
    }

    @Override
    public void accept(PlcSubscriptionEvent plcSubscriptionEvent) {
        LOGGER.debug("Received {}", plcSubscriptionEvent);
        try {
            Exchange exchange = endpoint.createExchange();
            exchange.getIn().setBody(unwrapIfSingle(plcSubscriptionEvent.getAllObjects("default")));
            processor.process(exchange);
        } catch (Exception e) {
            exceptionHandler.handleException(e);
        }
    }

    private Object unwrapIfSingle(Collection collection) {
        if (collection.isEmpty()) {
            return null;
        }
        if (collection.size() == 1) {
            return collection.iterator().next();
        }
        return collection;
    }

}
