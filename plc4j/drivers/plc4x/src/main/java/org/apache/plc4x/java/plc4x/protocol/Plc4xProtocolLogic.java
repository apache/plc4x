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
package org.apache.plc4x.java.plc4x.protocol;

import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.plc4x.config.Plc4xConfiguration;
import org.apache.plc4x.java.plc4x.context.Plc4xDriverContext;
import org.apache.plc4x.java.plc4x.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4xProtocolLogic extends Plc4xProtocolBase<Plc4xMessage> implements HasConfiguration<Plc4xConfiguration> {

    private final Logger logger = LoggerFactory.getLogger(Plc4xProtocolLogic.class);

    private final AtomicInteger txIdGenerator = new AtomicInteger(1);
    private String remoteConnectionString;
    private Duration requestTimeout;
    private RequestTransactionManager tm;
    private Plc4xDriverContext plc4xDriverContext;

    @Override
    public void setConfiguration(Plc4xConfiguration configuration) {
        this.remoteConnectionString = configuration.getRemoteConnectionString();
        this.requestTimeout = Duration.ofMillis(configuration.getRequestTimeout());
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);
        plc4xDriverContext = (Plc4xDriverContext) driverContext;
    }

    @Override
    public void onConnect(ConversationContext<Plc4xMessage> context) {
        final int txId = txIdGenerator.getAndIncrement();

        Plc4xConnectRequest connectRequest = new Plc4xConnectRequest(txId, remoteConnectionString);

        context.sendRequest(connectRequest)
            .onTimeout(e -> {
                logger.warn("Timeout during Connection establishing, closing channel...");
                context.getChannel().close();
            })
            .expectResponse(Plc4xMessage.class, requestTimeout)
            .onTimeout(e -> {
                logger.warn("Timeout during Connection establishing, closing channel...");
                context.getChannel().close();
            })
            .check(p -> p.getRequestId() == txId)
            .unwrap(plc4xMessage -> (Plc4xConnectResponse) plc4xMessage)
            .handle(connectResponse -> {
                // Save the connection id.
                plc4xDriverContext.setConnectionId(connectResponse.getConnectionId());
                logger.debug("Got Plc4x Connection Response");
                context.fireConnected();
            });
    }

    @Override
    public void onDisconnect(ConversationContext<Plc4xMessage> context) {
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        List<Plc4xFieldRequest> fields = new ArrayList<>(readRequest.getNumberOfFields());
        for (String fieldName : readRequest.getFieldNames()) {
            final org.apache.plc4x.java.plc4x.field.Plc4xField field = (org.apache.plc4x.java.plc4x.field.Plc4xField) readRequest.getField(fieldName);
            Plc4xFieldRequest fieldRequest = new Plc4xFieldRequest(new Plc4xField(fieldName, field.getAddress()));
            fields.add(fieldRequest);
        }
        final int txId = txIdGenerator.getAndIncrement();
        Plc4xReadRequest read = new Plc4xReadRequest(plc4xDriverContext.getConnectionId(), txId, fields);

        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        context.sendRequest(read)
            .expectResponse(Plc4xMessage.class, requestTimeout)
            .onTimeout(future::completeExceptionally)
            .check(p -> p.getRequestId() == txId)
            .unwrap(plc4xMessage -> (Plc4xReadResponse) plc4xMessage)
            .handle(plc4xReadResponse -> {
                // TODO: Do something with this response.
                future.complete(new DefaultPlcReadResponse(readRequest, null));//Map.of("test", null)
                // Finish the request-transaction.
                transaction.endRequest();
            });
        return future;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();

        List<Plc4xFieldValueRequest> fields = new ArrayList<>(writeRequest.getNumberOfFields());
        for (String fieldName : writeRequest.getFieldNames()) {
            final org.apache.plc4x.java.plc4x.field.Plc4xField field = (org.apache.plc4x.java.plc4x.field.Plc4xField) writeRequest.getField(fieldName);
            final Plc4xValueType plc4xValueType = Plc4xValueType.BOOL;
            final PlcValue plcValue = writeRequest.getPlcValue(fieldName);
            Plc4xFieldValueRequest fieldRequest = new Plc4xFieldValueRequest(
                new Plc4xField(fieldName, field.getAddress()), plc4xValueType, plcValue);
            fields.add(fieldRequest);
        }
        final int txId = txIdGenerator.getAndIncrement();
        Plc4xWriteRequest write = new Plc4xWriteRequest(plc4xDriverContext.getConnectionId(), txId, fields);

        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        context.sendRequest(write)
            .expectResponse(Plc4xMessage.class, requestTimeout)
            .onTimeout(future::completeExceptionally)
            .check(p -> p.getRequestId() == txId)
            .unwrap(plc4xMessage -> (Plc4xWriteResponse) plc4xMessage)
            .handle(plc4xWriteResponse -> {
                // TODO: Do something with this response.
                future.complete(new DefaultPlcWriteResponse(writeRequest, null)); // Map.of("test", PlcResponseCode.OK)
                // Finish the request-transaction.
                transaction.endRequest();
            });
        return future;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        return super.subscribe(subscriptionRequest);
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return super.unsubscribe(unsubscriptionRequest);
    }

    @Override
    public void close(ConversationContext<Plc4xMessage> context) {
        // Nothing to do here ...
    }

}
