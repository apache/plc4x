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
package org.apache.plc4x.java.amsads.protocol;

import org.apache.plc4x.java.amsads.configuration.AdsConfiguration;
import org.apache.plc4x.java.amsads.field.DirectAdsField;
import org.apache.plc4x.java.amsads.field.SymbolicAdsField;
import org.apache.plc4x.java.amsads.readwrite.*;
import org.apache.plc4x.java.amsads.readwrite.types.CommandId;
import org.apache.plc4x.java.amsads.readwrite.types.ReservedIndexGroups;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;

public class AdsProtocolLogic extends Plc4xProtocolBase<AmsPacket> implements HasConfiguration<AdsConfiguration> {

    private AdsConfiguration configuration;
    public static final State DEFAULT_COMMAND_STATE = new State(
        false, false, false, false, false, false, true, false, false);

    private ConversationContext<AmsPacket> adsDriverContext;
    private static final AtomicLong invokeIdGenerator = new AtomicLong(0);
    private RequestTransactionManager tm;

    private ConcurrentHashMap<SymbolicAdsField, DirectAdsField> symbolicFieldMapping;
    private ConcurrentHashMap<SymbolicAdsField, CompletableFuture<DirectAdsField>> pendingResolutionRequests;

    public AdsProtocolLogic() {
        symbolicFieldMapping = new ConcurrentHashMap<>();
        pendingResolutionRequests = new ConcurrentHashMap<>();

        // Initialize Transaction Manager.
        // Until the number of concurrent requests is successfully negotiated we set it to a
        // maximum of only one request being able to be sent at a time. During the login process
        // No concurrent requests can be sent anyway. It will be updated when receiving the
        // S7ParameterSetupCommunication response.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void setConfiguration(AdsConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public void setContext(ConversationContext<AmsPacket> context) {
        super.setContext(context);
        adsDriverContext = context;
    }

    @Override
    public void close(ConversationContext<AmsPacket> context) {

    }

    @Override
    public void onConnect(ConversationContext<AmsPacket> context) {
        // AMS/ADS doesn't know a concept of a connect.
        context.fireConnected();
    }

    @Override
    public void onDisconnect(ConversationContext<AmsPacket> context) {
        super.onDisconnect(context);
        // TODO: Here we have to clean up all of the handles this connection acquired.
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;

        // Get all ADS addresses in their resolved state.
        final List<DirectAdsField> directAdsFields = getDirectAddresses(readRequest.getFields());

        // Depending on the number of fields, use a single item request or a sum-request
        if(directAdsFields.size() == 1) {
            // Do a normal (single item) ADS Read Request
            return singleRead(readRequest, directAdsFields.get(0));
        } else {
            // Do a ADS-Sum Read Request.
            return multiRead(readRequest, directAdsFields);
        }
    }

    protected CompletableFuture<PlcReadResponse> singleRead(PlcReadRequest readRequest, DirectAdsField directAdsField) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        int size = directAdsField.getAdsDataType().getTargetByteSize() * directAdsField.getNumberOfElements();
        AdsData adsData = new AdsReadRequest(directAdsField.getIndexGroup(), directAdsField.getIndexOffset(), size);
        AmsPacket amsPacket = new AmsPacket(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            CommandId.ADS_READ, DEFAULT_COMMAND_STATE, 0, invokeIdGenerator.getAndIncrement(), adsData);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsPacket)
            .expectResponse(AmsPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getInvokeId() == amsPacket.getInvokeId())
            .unwrap(AmsPacket::getData)
            .handle(responseAdsData -> {
                final PlcReadResponse plcReadResponse = convertToPlc4xReadResponse(readRequest, responseAdsData);
                // Convert the response from the PLC into a PLC4X Response ...
                future.complete(plcReadResponse);
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    protected CompletableFuture<PlcReadResponse> multiRead(PlcReadRequest readRequest, List<DirectAdsField> directAdsFields) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        // Calculate the size of all fields together.
        int size = 0;

        // TODO: Add the items ...
        List<AdsReadRequest> items = new ArrayList<>(directAdsFields.size());

        // With multi-requests, the index-group is fixed and the index offset indicates the number of elements.
        AdsData adsData = new AdsReadWriteRequest(
            ReservedIndexGroups.ADSIGRP_MULTIPLE_READ.getValue(), directAdsFields.size(), size,
            items.toArray(new AdsReadRequest[0]), new byte[0]);

        AmsPacket amsPacket = new AmsPacket(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            CommandId.ADS_READ, DEFAULT_COMMAND_STATE, 0, invokeIdGenerator.getAndIncrement(), adsData);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsPacket)
            .expectResponse(AmsPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getInvokeId() == amsPacket.getInvokeId())
            .unwrap(AmsPacket::getData)
            .handle(responseAdsData -> {
                final PlcReadResponse plcReadResponse = convertToPlc4xReadResponse(readRequest, responseAdsData);
                // Convert the response from the PLC into a PLC4X Response ...
                future.complete(plcReadResponse);
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    protected PlcReadResponse convertToPlc4xReadResponse(PlcReadRequest readRequest, AdsData adsData) {
        return null;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        return super.write(writeRequest);
    }

    @Override
    protected void decode(ConversationContext<AmsPacket> context, AmsPacket msg) throws Exception {
        super.decode(context, msg);
    }

    protected List<DirectAdsField> getDirectAddresses(List<PlcField> fields) {
        // Get all symbolic fields from the current request.
        // These potentially need to be resolved to direct addresses, if this has not been done before.
        final List<SymbolicAdsField> referencedSymbolicFields = fields.stream()
            .filter(plcField -> plcField instanceof SymbolicAdsField)
            .map(plcField -> (SymbolicAdsField) plcField).collect(Collectors.toList());

        // Find out for which of these symbolic addresses no resolution has been initiated.
        final List<SymbolicAdsField> symbolicFieldsNeedingResolution = referencedSymbolicFields.stream()
            .filter(symbolicAdsField -> symbolicFieldMapping.containsKey(symbolicAdsField))
            .collect(Collectors.toList());

        // If there are unresolved symbolic addresses, initiate the resolution
        if(!symbolicFieldsNeedingResolution.isEmpty()) {
            // If a previous request initiated a resolution request, join that resolutions future instead.
            // If not, initiate a new resolution request.
            final CompletableFuture<Void> resolutionComplete =
                CompletableFuture.allOf(symbolicFieldsNeedingResolution.stream().map(symbolicAdsField -> {
                    if (pendingResolutionRequests.containsKey(symbolicAdsField)) {
                        return pendingResolutionRequests.get(symbolicAdsField);
                    } else {
                        // Initiate a new resolution-request and add that to the pending resolution requests.
                        CompletableFuture<DirectAdsField> internalResolutionFuture =
                            resolveSymbolicAddress(symbolicAdsField);
                        // Create a second future which will be completed as soon as the resolution result has
                        // been added to the map.
                        CompletableFuture<DirectAdsField> resolutionFuture = new CompletableFuture<>();
                        // Make sure the resolved address is added to the mapping.
                        internalResolutionFuture.thenAccept(directAdsField -> {
                            symbolicFieldMapping.put(symbolicAdsField, directAdsField);
                            // Now we can tell the other waiting processes about the result.
                            resolutionFuture.complete(directAdsField);
                        });
                        return resolutionFuture;
                    }
                }).toArray(CompletableFuture[]::new));
            // Wait for the resolution to finish.
            try {
                resolutionComplete.get(configuration.getTimeoutSymbolicAddressResolution(), TimeUnit.MILLISECONDS);
            } catch (TimeoutException e) {
                // TODO: Return an error indicating a timeout during resolution.
            } catch (InterruptedException | ExecutionException e) {
                // TODO: Return an error indicating a timeout an internal server error.
            }
        }

        // So here all fields should be resolved so we can continue normally.
        return fields.stream().map(plcField -> {
            if(plcField instanceof SymbolicAdsField) {
                return symbolicFieldMapping.get(plcField);
            } else {
                return (DirectAdsField) plcField;
            }
        }).collect(Collectors.toList());
    }

    protected CompletableFuture<DirectAdsField> resolveSymbolicAddress(SymbolicAdsField symbolicAdsField) {
        return null;
    }

}
