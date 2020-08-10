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
import org.apache.plc4x.java.amsads.readwrite.io.DataItemIO;
import org.apache.plc4x.java.amsads.readwrite.types.CommandId;
import org.apache.plc4x.java.amsads.readwrite.types.ReservedIndexGroups;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

public class AdsProtocolLogic extends Plc4xProtocolBase<AmsPacket> implements HasConfiguration<AdsConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsProtocolLogic.class);

    private AdsConfiguration configuration;
    public static final State DEFAULT_COMMAND_STATE = new State(
        false, false, false, false, false, true, false, false, false);

    private ConversationContext<AmsPacket> adsDriverContext;
    private final AtomicLong invokeIdGenerator = new AtomicLong(1);
    private RequestTransactionManager tm;

    private ConcurrentHashMap<SymbolicAdsField, DirectAdsField> symbolicFieldMapping;
    private ConcurrentHashMap<SymbolicAdsField, CompletableFuture<Void>> pendingResolutionRequests;

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
        if (directAdsFields.size() == 1) {
            // Do a normal (single item) ADS Read Request
            return singleRead(readRequest, directAdsFields.get(0));
        } else {
            // Do a ADS-Sum Read Request.
            return multiRead(readRequest, directAdsFields);
        }
    }

    protected CompletableFuture<PlcReadResponse> singleRead(PlcReadRequest readRequest, DirectAdsField directAdsField) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        int size = directAdsField.getAdsDataType().getNumBytes() * directAdsField.getNumberOfElements();
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

        // With multi-requests, the index-group is fixed and the index offset indicates the number of elements.
        AdsData adsData = new AdsReadWriteRequest(
            ReservedIndexGroups.ADSIGRP_MULTIPLE_READ.getValue(), directAdsFields.size(), size,
            directAdsFields.stream().map(directAdsField -> new AdsReadWriteRequest(
                directAdsField.getIndexGroup(), directAdsField.getIndexOffset(), directAdsField.getNumberOfElements(),
                new AdsReadWriteRequest[0], new byte[0])).toArray(AdsReadWriteRequest[]::new), new byte[0]);

        AmsPacket amsPacket = new AmsPacket(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            CommandId.ADS_READ_WRITE, DEFAULT_COMMAND_STATE, 0, invokeIdGenerator.getAndIncrement(), adsData);

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
        ReadBuffer readBuffer = null;
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        if (adsData instanceof AdsReadResponse) {
            AdsReadResponse adsReadResponse = (AdsReadResponse) adsData;
            readBuffer = new ReadBuffer(adsReadResponse.getData());
            responseCodes.put(readRequest.getFieldNames().stream().findFirst().orElse(""),
                parsePlcResponseCode(adsReadResponse.getResult()));
        } else if (adsData instanceof AdsReadWriteResponse) {
            AdsReadWriteResponse adsReadWriteResponse = (AdsReadWriteResponse) adsData;
            readBuffer = new ReadBuffer(adsReadWriteResponse.getData());
            // When parsing a multi-item response, the error codes of each items come
            // in sequence and then come the values.
            for (String fieldName : readRequest.getFieldNames()) {
                try {
                    final long result = readBuffer.readUnsignedLong(32);
                    responseCodes.put(fieldName, parsePlcResponseCode(result));
                } catch (ParseException e) {
                    responseCodes.put(fieldName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }
        if(readBuffer != null) {
            Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
            for (String fieldName : readRequest.getFieldNames()) {
                DirectAdsField directAdsField = (DirectAdsField) readRequest.getField(fieldName);
                // If the response-code was anything but OK, we don't need to parse the payload.
                if(responseCodes.get(fieldName) != PlcResponseCode.OK) {
                    values.put(fieldName, new ResponseItem<>(responseCodes.get(fieldName), null));
                }
                // If the response-code was ok, parse the data returned.
                else {
                    values.put(fieldName, parsePlcValue(directAdsField, readBuffer));
                }
            }
            return new DefaultPlcReadResponse((InternalPlcReadRequest) readRequest, values);
        }
        return null;
    }

    private PlcResponseCode parsePlcResponseCode(long adsResult) {
            if (adsResult == 0L) {
                return PlcResponseCode.OK;
            } else {
                // TODO: Implement this a little more ...
                return PlcResponseCode.INTERNAL_ERROR;
            }
    }

    private ResponseItem<PlcValue> parsePlcValue(DirectAdsField field, ReadBuffer readBuffer) {
        try {
            if (field.getNumberOfElements() == 1) {
                return new ResponseItem<>(PlcResponseCode.OK,
                    DataItemIO.staticParse(readBuffer, field.getAdsDataType()));
            } else {
                // Fetch all
                final PlcValue[] resultItems = IntStream.range(0, field.getNumberOfElements()).mapToObj(i -> {
                    try {
                        return DataItemIO.staticParse(readBuffer, field.getAdsDataType());
                    } catch (ParseException e) {
                        LOGGER.warn("Error parsing field item of type: '{}' (at position {}})", field.getAdsDataType(), i, e);
                    }
                    return null;
                }).toArray(PlcValue[]::new);
                return new ResponseItem<>(PlcResponseCode.OK, PlcValues.of(resultItems));
            }
        } catch (ParseException e) {
            LOGGER.warn(String.format("Error parsing field item of type: '%s'", field.getAdsDataType()), e);
            return new ResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null);
        }
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
            .filter(symbolicAdsField -> !symbolicFieldMapping.containsKey(symbolicAdsField))
            .collect(Collectors.toList());

        // If there are unresolved symbolic addresses, initiate the resolution
        if (!symbolicFieldsNeedingResolution.isEmpty()) {
            // Get a list of symbolic addresses for which no resolution request has been sent yet
            // (A parallel request initiated a bit earlier might have already initiated a resolution
            // which has not yet been completed)
            final List<SymbolicAdsField> requiredResolutionFields =
                symbolicFieldsNeedingResolution.stream().filter(symbolicAdsField ->
                    !pendingResolutionRequests.containsKey(symbolicAdsField)).collect(Collectors.toList());
            // If there are fields for which no resolution request has been sent yet,
            // send a request.
            if (!requiredResolutionFields.isEmpty()) {
                CompletableFuture<Void> resolutionFuture;
                // Create a future which will be completed as soon as the
                // resolution result has been added to the map.
                if (requiredResolutionFields.size() == 1) {
                    SymbolicAdsField symbolicAdsField = requiredResolutionFields.get(0);
                    resolutionFuture = resolveSymbolicAddress(requiredResolutionFields.get(0));
                    pendingResolutionRequests.put(symbolicAdsField, resolutionFuture);
                } else {
                    resolutionFuture = resolveSymbolicAddresses(requiredResolutionFields);
                    for (SymbolicAdsField symbolicAdsField : requiredResolutionFields) {
                        pendingResolutionRequests.put(symbolicAdsField, resolutionFuture);
                    }
                }
            }

            // Create a global future which is completed as soon as all sub-futures for this request are completed.
            final CompletableFuture<Void> resolutionComplete =
                CompletableFuture.allOf(symbolicFieldsNeedingResolution.stream()
                    .map(symbolicAdsField -> pendingResolutionRequests.get(symbolicAdsField))
                    .toArray(CompletableFuture[]::new));

            // BLOCKING: Wait for the resolution to finish.
            // TODO: Make this asynchronous ...
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
            if (plcField instanceof SymbolicAdsField) {
                return symbolicFieldMapping.get(plcField);
            } else {
                return (DirectAdsField) plcField;
            }
        }).collect(Collectors.toList());
    }

    protected CompletableFuture<Void> resolveSymbolicAddress(SymbolicAdsField symbolicAdsField) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        // TODO: Instead of using 4 we need the size of the expected response
        AdsData adsData = new AdsReadWriteRequest(ReservedIndexGroups.ADSIGRP_SYM_HNDBYNAME.getValue(), 0, 4, null,
            getNullByteTerminatedArray(symbolicAdsField.getSymbolicField()));
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
            .check(adsDataResponse -> adsDataResponse instanceof AdsReadWriteResponse)
            .unwrap(adsDataResponse -> (AdsReadWriteResponse) adsDataResponse)
            .handle(responseAdsData -> {
                ReadBuffer readBuffer = new ReadBuffer(responseAdsData.getData());
                try {
                    // This should be 0 in the success case.
                    long returnCode = readBuffer.readLong(32);
                    // This is always 4
                    long itemLength = readBuffer.readLong(32);
                    // Get the handle from the response.
                    long handle = readBuffer.readLong(32);
                    if (returnCode == 0) {
                        DirectAdsField directAdsField = new DirectAdsField(ReservedIndexGroups.ADSIGRP_SYM_VALBYHND.getValue(),
                            handle, symbolicAdsField.getAdsDataType(), symbolicAdsField.getNumberOfElements());
                        symbolicFieldMapping.put(symbolicAdsField, directAdsField);
                        future.complete(null);
                    } else {
                        // TODO: Handle the case of unsuccessful resolution ..
                    }
                } catch (ParseException e) {
                    e.printStackTrace();
                }
            }));
        return future;
    }

    protected CompletableFuture<Void> resolveSymbolicAddresses(List<SymbolicAdsField> symbolicAdsFields) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        // TODO: Instead of using 4 we need the size of the expected response
        AdsData adsData = new AdsReadWriteRequest(ReservedIndexGroups.ADSIGRP_MULTIPLE_GET_HANDLE.getValue(),
            symbolicAdsFields.size(), 4, symbolicAdsFields.stream().map(symbolicAdsField ->
            new AdsReadWriteRequest(ReservedIndexGroups.ADSIGRP_SYM_HNDBYNAME.getValue(), 0, 4, null,
                getNullByteTerminatedArray(symbolicAdsField.getSymbolicField()))).toArray(AdsReadWriteRequest[]::new), null);
        AmsPacket amsPacket = new AmsPacket(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            CommandId.ADS_READ_WRITE, DEFAULT_COMMAND_STATE, 0, invokeIdGenerator.getAndIncrement(), adsData);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsPacket)
            .expectResponse(AmsPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getInvokeId() == amsPacket.getInvokeId())
            .unwrap(AmsPacket::getData)
            .check(adsDataResponse -> adsDataResponse instanceof AdsReadWriteResponse)
            .unwrap(adsDataResponse -> (AdsReadWriteResponse) adsDataResponse)
            .handle(responseAdsData -> {
                ReadBuffer readBuffer = new ReadBuffer(responseAdsData.getData());
                Map<SymbolicAdsField, Long> returnCodes = new HashMap<>();
                symbolicAdsFields.forEach(symbolicAdsField -> {
                    try {
                        // This should be 0 in the success case.
                        long returnCode = readBuffer.readLong(32);
                        // This is always 4
                        long itemLength = readBuffer.readLong(32);

                        returnCodes.put(symbolicAdsField, returnCode);
                    } catch (ParseException e) {
                        e.printStackTrace();
                    }
                });
                symbolicAdsFields.forEach(symbolicAdsField -> {
                    try {
                        if (returnCodes.get(symbolicAdsField) == 0) {
                            // Read the handle.
                            long handle = readBuffer.readLong(32);

                            DirectAdsField directAdsField = new DirectAdsField(
                                ReservedIndexGroups.ADSIGRP_SYM_VALBYHND.getValue(), handle,
                                symbolicAdsField.getAdsDataType(), symbolicAdsField.getNumberOfElements());
                            symbolicFieldMapping.put(symbolicAdsField, directAdsField);
                        } else {
                            // TODO: Handle the case of unsuccessful resolution ..
                        }
                    } catch (ParseException e) {
                        e.printStackTrace();
                    }
                });
                future.complete(null);
            }));
        return future;
    }

    protected byte[] getNullByteTerminatedArray(String value) {
        byte[] valueBytes = value.getBytes();
        byte[] nullTerminatedBytes = new byte[valueBytes.length + 1];
        System.arraycopy(valueBytes, 0, nullTerminatedBytes, 0, valueBytes.length);
        return nullTerminatedBytes;
    }

}
