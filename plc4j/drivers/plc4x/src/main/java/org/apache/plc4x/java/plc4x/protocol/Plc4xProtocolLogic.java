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
package org.apache.plc4x.java.plc4x.protocol;

import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.plc4x.config.Plc4xConfiguration;
import org.apache.plc4x.java.plc4x.tag.Plc4xTag;
import org.apache.plc4x.java.plc4x.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4xProtocolLogic extends Plc4xProtocolBase<Plc4xMessage> implements HasConfiguration<Plc4xConfiguration> {

    private final Logger logger = LoggerFactory.getLogger(Plc4xProtocolLogic.class);

    private final AtomicInteger txIdGenerator = new AtomicInteger(1);
    private String remoteConnectionString;
    private Duration requestTimeout;
    private RequestTransactionManager tm;
    private int connectionId;

    @Override
    public void setConfiguration(Plc4xConfiguration configuration) {
        this.tm = new RequestTransactionManager(1);
        this.remoteConnectionString = configuration.getRemoteConnectionString();
        this.requestTimeout = Duration.ofMillis(configuration.getRequestTimeout());
        this.connectionId = 0;
    }

    @Override
    public void onConnect(ConversationContext<Plc4xMessage> context) {
        final int requestId = txIdGenerator.getAndIncrement();

        Plc4xConnectRequest connectRequest = new Plc4xConnectRequest(requestId, remoteConnectionString);

        context.sendRequest(connectRequest)
            .onTimeout(e -> {
                logger.warn("Timeout during Connection establishing, closing channel...");
                context.getChannel().close();
            })
            .expectResponse(Plc4xMessage.class, requestTimeout)
            .check(p -> p.getRequestId() == requestId)
            .unwrap(plc4xMessage -> (Plc4xConnectResponse) plc4xMessage)
            .handle(connectResponse -> {
                // Save the connection id.
                connectionId = connectResponse.getConnectionId();
                logger.debug("Got Plc4x Connection Response");
                context.fireConnected();
            });
    }

    @Override
    public void onDisconnect(ConversationContext<Plc4xMessage> context) {
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest apiReadRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        // Prepare the request.
        List<Plc4xTagRequest> plc4xTags = new ArrayList<>(apiReadRequest.getNumberOfTags());
        for (String tagName : apiReadRequest.getTagNames()) {
            final Plc4xTag plc4xTag =
                (Plc4xTag) apiReadRequest.getTag(tagName);
            Plc4xTagRequest plc4XTagRequest = new Plc4xTagRequest(
                new org.apache.plc4x.java.plc4x.readwrite.Plc4xTag(tagName, plc4xTag.getAddressString() + ":" + plc4xTag.getPlcValueType().name()));
            plc4xTags.add(plc4XTagRequest);
        }
        final int requestId = txIdGenerator.getAndIncrement();
        Plc4xReadRequest plc4xReadRequest = new Plc4xReadRequest(requestId, connectionId, plc4xTags);

        // Send the request and await a response.
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        context.sendRequest(plc4xReadRequest)
            .expectResponse(Plc4xMessage.class, requestTimeout)
            .onTimeout(future::completeExceptionally)
            .check(plc4xMessage -> plc4xMessage.getRequestId() == requestId)
            .unwrap(plc4xMessage -> (Plc4xReadResponse) plc4xMessage)
            .check(plc4xReadResponse -> plc4xReadResponse.getConnectionId() == connectionId)
            .handle(plc4xReadResponse -> {
                Map<String, ResponseItem<PlcValue>> apiResponses = new HashMap<>();
                // Create the API response from the incoming message.
                for (Plc4xTagValueResponse plc4xTag : plc4xReadResponse.getTags()) {
                    final Plc4xResponseCode plc4xResponseCode = plc4xTag.getResponseCode();
                    final PlcResponseCode apiResponseCode = PlcResponseCode.valueOf(plc4xResponseCode.name());
                    apiResponses.put(plc4xTag.getTag().getName(),
                        new ResponseItem<>(apiResponseCode, plc4xTag.getValue()));
                }

                // Send it back to the calling process.
                future.complete(new DefaultPlcReadResponse(apiReadRequest, apiResponses));

                // Finish the request-transaction.
                transaction.endRequest();
            });
        return future;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();

        // Prepare the request.
        List<Plc4xTagValueRequest> tags = new ArrayList<>(writeRequest.getNumberOfTags());
        for (String tagName : writeRequest.getTagNames()) {
            final Plc4xTag plc4xTag = (Plc4xTag) writeRequest.getTag(tagName);
            final Plc4xValueType plc4xValueType = Plc4xValueType.valueOf(plc4xTag.getPlcValueType().name());
            final PlcValue plcValue = writeRequest.getPlcValue(tagName);
            Plc4xTagValueRequest tagRequest = new Plc4xTagValueRequest(
                new org.apache.plc4x.java.plc4x.readwrite.Plc4xTag(tagName, plc4xTag.getAddressString() + ":" + plc4xTag.getPlcValueType().name()), plc4xValueType, plcValue);
            tags.add(tagRequest);
        }
        final int requestId = txIdGenerator.getAndIncrement();
        Plc4xWriteRequest write = new Plc4xWriteRequest(requestId, connectionId, tags);

        // Send the request and await a response.
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        context.sendRequest(write)
            .expectResponse(Plc4xMessage.class, requestTimeout)
            .onTimeout(future::completeExceptionally)
            .check(p -> p.getRequestId() == requestId)
            .unwrap(plc4xMessage -> (Plc4xWriteResponse) plc4xMessage)
            .check(plc4xReadResponse -> plc4xReadResponse.getConnectionId() == connectionId)
            .handle(plc4xWriteResponse -> {
                Map<String, PlcResponseCode> apiResponses = new HashMap<>();
                // Create the API response from the incoming message.
                for (Plc4xTagResponse plc4xTag : plc4xWriteResponse.getTags()) {
                    final Plc4xResponseCode plc4xResponseCode = plc4xTag.getResponseCode();
                    final PlcResponseCode apiResponseCode = PlcResponseCode.valueOf(plc4xResponseCode.name());
                    apiResponses.put(plc4xTag.getTag().getName(), apiResponseCode);
                }

                // Send it back to the calling process.
                future.complete(new DefaultPlcWriteResponse(writeRequest, apiResponses));

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

    @Override
    protected void decode(ConversationContext<Plc4xMessage> context, Plc4xMessage msg) throws Exception {
        super.decode(context, msg);
    }

}
