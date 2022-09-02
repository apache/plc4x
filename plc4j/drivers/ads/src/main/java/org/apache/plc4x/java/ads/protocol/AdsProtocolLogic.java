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
package org.apache.plc4x.java.ads.protocol;

import org.apache.plc4x.java.ads.configuration.AdsConfiguration;
import org.apache.plc4x.java.ads.discovery.readwrite.*;
import org.apache.plc4x.java.ads.discovery.readwrite.AmsNetId;
import org.apache.plc4x.java.ads.field.*;
import org.apache.plc4x.java.ads.model.AdsSubscriptionHandle;
import org.apache.plc4x.java.ads.readwrite.*;
import org.apache.plc4x.java.ads.readwrite.DataItem;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.net.*;
import java.time.Duration;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Consumer;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

public class AdsProtocolLogic extends Plc4xProtocolBase<AmsTCPPacket> implements HasConfiguration<AdsConfiguration>, PlcSubscriber, PlcBrowser {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsProtocolLogic.class);

    private AdsConfiguration configuration;

    private final AtomicLong invokeIdGenerator = new AtomicLong(1);
    private final RequestTransactionManager tm;

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    //    private final ConcurrentHashMap<SymbolicAdsField, DirectAdsField> symbolicFieldMapping;
    private final ConcurrentHashMap<SymbolicAdsField, CompletableFuture<Void>> pendingResolutionRequests;

    private final Map<String, AdsSymbolTableEntry> symbolTable;
    private final Map<String, AdsDataTypeTableEntry> dataTypeTable;

    public AdsProtocolLogic() {
//        symbolicFieldMapping = new ConcurrentHashMap<>();
        pendingResolutionRequests = new ConcurrentHashMap<>();
        symbolTable = new HashMap<>();
        dataTypeTable = new HashMap<>();

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
    public void close(ConversationContext<AmsTCPPacket> context) {

    }

    @Override
    public void onConnect(ConversationContext<AmsTCPPacket> context) {
        final CompletableFuture<Void> future = new CompletableFuture<>();

        // If we have connection credentials available, try to set up the AMS routes.
        CompletableFuture<Void> setupAmsRouteFuture;
        if (context.getAuthentication() != null) {
            if (!(context.getAuthentication() instanceof PlcUsernamePasswordAuthentication)) {
                future.completeExceptionally(new PlcConnectionException(
                    "This type of connection only supports username-password authentication"));
                return;
            }
            PlcUsernamePasswordAuthentication usernamePasswordAuthentication =
                (PlcUsernamePasswordAuthentication) context.getAuthentication();

            setupAmsRouteFuture = setupAmsRoute(usernamePasswordAuthentication);
        } else {
            setupAmsRouteFuture = CompletableFuture.completedFuture(null);
        }

        // If the configuration asks us to load the symbol and data type tables, do so,
        // otherwise just mark the connection as completed instantly.
        setupAmsRouteFuture.whenComplete((unused, throwable) -> {
            if (configuration.isLoadSymbolAndDataTypeTables()) {
                LOGGER.debug("Fetching sizes of symbol and datatype table sizes.");
                CompletableFuture<Void> readSymbolTableFuture = readSymbolTableAndDatatypeTable(context);
                readSymbolTableFuture.whenComplete((unused2, throwable2) -> {
                    if (throwable2 != null) {
                        LOGGER.error("Error fetching symbol and datatype table sizes");
                    } else {
                        context.fireConnected();
                    }
                });
            } else {
                context.fireConnected();
            }
        });
    }

    protected CompletableFuture<Void> setupAmsRoute(PlcUsernamePasswordAuthentication authentication) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        new Thread(() -> {
            LOGGER.debug("Setting up remote AMS routes.");
            SocketAddress localSocketAddress = context.getChannel().localAddress();
            InetAddress localAddress = ((InetSocketAddress) localSocketAddress).getAddress();

            // Prepare the request message.
            AmsNetId sourceAmsNetId = new AmsNetId(
                configuration.getSourceAmsNetId().getOctet1(), configuration.getSourceAmsNetId().getOctet2(),
                configuration.getSourceAmsNetId().getOctet3(), configuration.getSourceAmsNetId().getOctet4(),
                configuration.getSourceAmsNetId().getOctet5(), configuration.getSourceAmsNetId().getOctet6());
            String routeName = String.format("PLC4X-%d.%d.%d.%d.%d.%d",
                sourceAmsNetId.getOctet1(), sourceAmsNetId.getOctet2(), sourceAmsNetId.getOctet3(),
                sourceAmsNetId.getOctet4(), sourceAmsNetId.getOctet5(), sourceAmsNetId.getOctet6());
            AdsDiscovery addOrUpdateRouteRequest = new AdsDiscovery(getInvokeId(), Operation.ADD_OR_UPDATE_ROUTE_REQUEST,
                sourceAmsNetId, AdsPortNumbers.SYSTEM_SERVICE,
                Arrays.asList(new AdsDiscoveryBlockRouteName(new AmsString(routeName)),
                    new AdsDiscoveryBlockAmsNetId(sourceAmsNetId),
                    new AdsDiscoveryBlockUserName(new AmsString(authentication.getUsername())),
                    new AdsDiscoveryBlockPassword(new AmsString(authentication.getPassword())),
                    new AdsDiscoveryBlockHostName(new AmsString(localAddress.getHostAddress()))));

            // Send the request to the PLC using a UDP socket.
            try (DatagramSocket adsDiscoverySocket = new DatagramSocket(AdsDiscoveryConstants.ADSDISCOVERYUDPDEFAULTPORT)) {
                // Serialize the message.
                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(
                    addOrUpdateRouteRequest.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
                addOrUpdateRouteRequest.serialize(writeBuffer);

                // Get the target IP from the connection
                SocketAddress remoteSocketAddress = context.getChannel().remoteAddress();
                InetAddress remoteAddress = ((InetSocketAddress) remoteSocketAddress).getAddress();

                // Create the UDP packet to the broadcast address.
                DatagramPacket discoveryRequestPacket = new DatagramPacket(
                    writeBuffer.getBytes(), writeBuffer.getBytes().length,
                    remoteAddress, AdsDiscoveryConstants.ADSDISCOVERYUDPDEFAULTPORT);
                adsDiscoverySocket.send(discoveryRequestPacket);

                // The actual length would be 32, but better be prepared for a more verbose response
                byte[] buf = new byte[100];
                DatagramPacket responsePacket = new DatagramPacket(buf, buf.length);
                adsDiscoverySocket.setSoTimeout(configuration.getTimeoutRequest());
                adsDiscoverySocket.receive(responsePacket);

                // Receive the response
                ReadBufferByteBased readBuffer = new ReadBufferByteBased(responsePacket.getData(), ByteOrder.LITTLE_ENDIAN);
                AdsDiscovery addOrUpdateRouteResponse = AdsDiscovery.staticParse(readBuffer);

                // Check if adding the route was successful
                if (addOrUpdateRouteResponse.getRequestId() == 1) {
                    for (AdsDiscoveryBlock block : addOrUpdateRouteResponse.getBlocks()) {
                        if (block.getBlockType() == AdsDiscoveryBlockType.STATUS) {
                            AdsDiscoveryBlockStatus statusBlock = (AdsDiscoveryBlockStatus) block;
                            if (statusBlock.getStatus() != Status.SUCCESS) {
                                future.completeExceptionally(new PlcException("Error adding AMS route"));
                                return;
                            }
                        }
                    }
                }

                future.complete(null);
            } catch (Exception e) {
                future.completeExceptionally(new PlcException("Error adding AMS route", e));
            }
        }).start();

        return future;
    }

    protected CompletableFuture<Void> readSymbolTableAndDatatypeTable(ConversationContext<AmsTCPPacket> context) {
        final CompletableFuture<Void> future = new CompletableFuture<>();

        // Initialize the request.
        AmsPacket amsPacket = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
            ReservedIndexGroups.ADSIGRP_SYMBOL_AND_DATA_TYPE_SIZES.getValue(), 0x00000000, 24);
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);
        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
            .unwrap(response -> (AdsReadResponse) response.getUserdata())
            .handle(responseAdsData -> {
                transaction.endRequest();
                if (responseAdsData.getResult() == ReturnCode.OK) {
                    ReadBuffer readBuffer = new ReadBufferByteBased(responseAdsData.getData());
                    try {
                        AdsTableSizes adsTableSizes = AdsTableSizes.staticParse(readBuffer);
                        LOGGER.debug("PLC contains {} symbols and {} data-types", adsTableSizes.getSymbolCount(), adsTableSizes.getDataTypeCount());

                        // Now we load the datatype definitions.
                        AmsPacket amsReadTablePacket = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
                            ReservedIndexGroups.ADSIGRP_DATA_TYPE_TABLE_UPLOAD.getValue(), 0x00000000, adsTableSizes.getDataTypeLength());
                        RequestTransactionManager.RequestTransaction transaction2 = tm.startRequest();
                        AmsTCPPacket amsReadTableTCPPacket = new AmsTCPPacket(amsReadTablePacket);
                        transaction2.submit(() -> context.sendRequest(amsReadTableTCPPacket)
                            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                            .onTimeout(future::completeExceptionally)
                            .onError((p, e) -> future.completeExceptionally(e))
                            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsReadTablePacket.getInvokeId())
                            .unwrap(response -> (AdsReadResponse) response.getUserdata())
                            .handle(responseAdsReadTableData -> {
                                transaction2.endRequest();
                                if (responseAdsData.getResult() == ReturnCode.OK) {
                                    // Parse the result.
                                    ReadBuffer rb = new ReadBufferByteBased(responseAdsReadTableData.getData());
                                    for (int i = 0; i < adsTableSizes.getDataTypeCount(); i++) {
                                        try {
                                            AdsDataTypeTableEntry adsDataTypeTableEntry = AdsDataTypeTableEntry.staticParse(rb);
                                            dataTypeTable.put(adsDataTypeTableEntry.getDataTypeName(), adsDataTypeTableEntry);
                                        } catch (ParseException e) {
                                            throw new RuntimeException(e);
                                        }
                                    }

                                    AmsPacket amsReadSymbolTablePacket = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                                        configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
                                        ReservedIndexGroups.ADSIGRP_SYM_UPLOAD.getValue(), 0x00000000, adsTableSizes.getSymbolLength());
                                    RequestTransactionManager.RequestTransaction transaction3 = tm.startRequest();
                                    AmsTCPPacket amsReadSymbolTableTCPPacket = new AmsTCPPacket(amsReadSymbolTablePacket);
                                    transaction3.submit(() -> context.sendRequest(amsReadSymbolTableTCPPacket)
                                        .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                                        .onTimeout(future::completeExceptionally)
                                        .onError((p, e) -> future.completeExceptionally(e))
                                        .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsReadSymbolTablePacket.getInvokeId())
                                        .unwrap(response -> (AdsReadResponse) response.getUserdata())
                                        .handle(responseAdsReadSymbolTableData -> {
                                            transaction3.endRequest();
                                            if (responseAdsData.getResult() == ReturnCode.OK) {
                                                ReadBuffer rb2 = new ReadBufferByteBased(responseAdsReadSymbolTableData.getData());
                                                for (int i = 0; i < adsTableSizes.getSymbolCount(); i++) {
                                                    try {
                                                        AdsSymbolTableEntry adsSymbolTableEntry = AdsSymbolTableEntry.staticParse(rb2);
                                                        symbolTable.put(adsSymbolTableEntry.getName(), adsSymbolTableEntry);
                                                    } catch (ParseException e) {
                                                        throw new RuntimeException(e);
                                                    }
                                                }
                                                future.complete(null);
                                            }
                                        }));
                                }
                            }));
                    } catch (ParseException e) {
                        future.completeExceptionally(new PlcException("Error loading the table sizes", e));
                    }
                } else {
                    // TODO: Implement this correctly.
                    future.completeExceptionally(new PlcException("Result is " + responseAdsData.getResult()));
                }
            }));
        return future;
    }

    @Override
    public void onDisconnect(ConversationContext<AmsTCPPacket> context) {
        super.onDisconnect(context);
        // TODO: Here we have to clean up all of the handles this connection acquired.
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        CompletableFuture<PlcBrowseResponse> future = new CompletableFuture<>();
        List<PlcBrowseItem> values = new ArrayList<>(symbolTable.size());
        for (AdsSymbolTableEntry symbol : symbolTable.values()) {
            // Get the datatype of this entry.
            AdsDataTypeTableEntry dataType = dataTypeTable.get(symbol.getDataTypeName());
            if (dataType == null) {
                System.out.printf("couldn't find datatype: %s%n", symbol.getDataTypeName());
                continue;
            }
            String itemName = (symbol.getComment() == null || symbol.getComment().isEmpty()) ? symbol.getName() : symbol.getComment();
            // Convert the plc value type from the ADS specific one to the PLC4X global one.
            org.apache.plc4x.java.api.types.PlcValueType plc4xPlcValueType = org.apache.plc4x.java.api.types.PlcValueType.valueOf(getPlcValueTypeForAdsDataType(dataType).toString());

            // If this type has children, add entries for its children.
            List<PlcBrowseItem> children = getBrowseItems(symbol.getName(), symbol.getGroup(), symbol.getOffset(), !symbol.getFlagReadOnly(), dataType);
            // Add the type itself.
            values.add(new DefaultPlcBrowseItem(symbol.getName(), itemName, plc4xPlcValueType, true, !symbol.getFlagReadOnly(), true, children));
        }
        DefaultPlcBrowseResponse response = new DefaultPlcBrowseResponse(browseRequest, PlcResponseCode.OK, values);
        future.complete(response);
        return future;
    }

    protected List<PlcBrowseItem> getBrowseItems(String basePath, long baseGroupId, long baseOffset, boolean parentWritable, AdsDataTypeTableEntry dataType) {
        if (dataType.getNumChildren() == 0) {
            return Collections.emptyList();
        }

        List<PlcBrowseItem> values = new ArrayList<>(dataType.getNumChildren());
        for (AdsDataTypeTableChildEntry child : dataType.getChildren()) {
            AdsDataTypeTableEntry childDataType = dataTypeTable.get(child.getDataTypeName());
            if (childDataType == null) {
                System.out.printf("couldn't find datatype: %s%n", child.getDataTypeName());
                continue;
            }
            String itemAddress = basePath + "." + child.getPropertyName();
            String itemName = (child.getComment() == null || child.getComment().isEmpty()) ? child.getPropertyName() : child.getComment();
            // Convert the plc value type from the ADS specific one to the PLC4X global one.
            org.apache.plc4x.java.api.types.PlcValueType plc4xPlcValueType = org.apache.plc4x.java.api.types.PlcValueType.valueOf(getPlcValueTypeForAdsDataType(childDataType).toString());
            // Recursively add all children of the current datatype.
            List<PlcBrowseItem> children = getBrowseItems(itemAddress, baseGroupId, baseOffset + child.getOffset(), parentWritable, childDataType);
            // Add the type itself.
            values.add(new DefaultPlcBrowseItem(basePath + "." + child.getPropertyName(), itemName, plc4xPlcValueType, true, parentWritable, true, children));
        }
        return values;
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        // Get all ADS addresses in their resolved state.
        final CompletableFuture<List<DirectAdsField>> directAdsFieldsFuture =
            getDirectAddresses(readRequest.getFields());

        // If all addresses were already resolved we can send the request immediately.
        if (directAdsFieldsFuture.isDone()) {
            final List<DirectAdsField> fields = directAdsFieldsFuture.getNow(null);
            if (fields != null) {
                return executeRead(readRequest, fields);
            } else {
                final CompletableFuture<PlcReadResponse> errorFuture = new CompletableFuture<>();
                errorFuture.completeExceptionally(new PlcException("Fields are null"));
                return errorFuture;
            }
        } else {
            // If there are still symbolic addresses that have to be resolved, send the
            // request as soon as the resolution is done.
            // In order to instantly be able to return a future, for the final result we have to
            // create a new one which is then completed later on. Unfortunately as soon as the
            // directAdsFieldsFuture is completed we still don't have the end result, but we can
            // now actually send the delayed read request ... as soon as that future completes
            // we can complete the initial one.
            CompletableFuture<PlcReadResponse> delayedRead = new CompletableFuture<>();
            directAdsFieldsFuture.handle((directAdsFields, throwable) -> {
                if (directAdsFields != null) {
                    final CompletableFuture<PlcReadResponse> delayedResponse =
                        executeRead(readRequest, directAdsFields);
                    delayedResponse.handle((plcReadResponse, throwable1) -> {
                        if (plcReadResponse != null) {
                            delayedRead.complete(plcReadResponse);
                        } else {
                            delayedRead.completeExceptionally(throwable1);
                        }
                        return this;
                    });
                } else {
                    delayedRead.completeExceptionally(throwable);
                }
                return this;
            });
            return delayedRead;
        }
    }

    protected CompletableFuture<PlcReadResponse> executeRead(PlcReadRequest readRequest,
                                                             List<DirectAdsField> directAdsFields) {
        // Depending on the number of fields, use a single item request or a sum-request
        if (directAdsFields.size() == 1) {
            // Do a normal (single item) ADS Read Request
            return singleRead(readRequest, directAdsFields.get(0));
        } else {
            // TODO: Check if the version of the remote station is at least TwinCAT v2.11 Build >= 1550 otherwise split up into single item requests.
            // Do a ADS-Sum Read Request.
            return multiRead(readRequest, directAdsFields);
        }
    }

    protected CompletableFuture<PlcReadResponse> singleRead(PlcReadRequest readRequest, DirectAdsField directAdsField) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        String dataTypeName = directAdsField.getAdsDataTypeName();
        AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(dataTypeName);
        long size = adsDataTypeTableEntry.getSize();

        AmsPacket amsPacket = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
            directAdsField.getIndexGroup(), directAdsField.getIndexOffset(), size * directAdsField.getNumberOfElements());
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
            .unwrap(amsResponsePacket -> (AdsReadResponse) amsResponsePacket.getUserdata())
            .handle(response -> {
                if (response.getResult() == ReturnCode.OK) {
                    final PlcReadResponse plcReadResponse = convertToPlc4xReadResponse(readRequest, response);
                    // Convert the response from the PLC into a PLC4X Response ...
                    future.complete(plcReadResponse);
                } else {
                    // TODO: Implement this correctly.
                    future.completeExceptionally(new PlcException("Result is " + response.getResult()));
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    protected CompletableFuture<PlcReadResponse> multiRead(PlcReadRequest readRequest, List<DirectAdsField> directAdsFields) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        // Calculate the size of all fields together.
        // Calculate the expected size of the response data.
        long expectedResponseDataSize = directAdsFields.stream().mapToLong(
            field -> {
                String dataTypeName = field.getAdsDataTypeName();
                AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(dataTypeName);
                long size = adsDataTypeTableEntry.getSize();
                // Status code + payload size
                return 4 + (size * field.getNumberOfElements());
            }).sum();

        // With multi-requests, the index-group is fixed and the index offset indicates the number of elements.
        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_MULTIPLE_READ.getValue(), directAdsFields.size(),
            expectedResponseDataSize, directAdsFields.stream().map(directAdsField -> {
            String dataTypeName = directAdsField.getAdsDataTypeName();
            AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(dataTypeName);
            long size = adsDataTypeTableEntry.getSize();
            return new AdsMultiRequestItemRead(
                directAdsField.getIndexGroup(), directAdsField.getIndexOffset(),
                (size * directAdsField.getNumberOfElements()));
        }).collect(Collectors.toList()), null);
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
            .unwrap(amsResponsePacket -> (AdsReadWriteResponse) amsResponsePacket.getUserdata())
            .handle(response -> {
                if (response.getResult() == ReturnCode.OK) {
                    final PlcReadResponse plcReadResponse = convertToPlc4xReadResponse(readRequest, response);
                    // Convert the response from the PLC into a PLC4X Response ...
                    future.complete(plcReadResponse);
                } else {
                    if (response.getResult() == ReturnCode.ADSERR_DEVICE_INVALIDSIZE) {
                        future.completeExceptionally(
                            new PlcException("The parameter size was not correct (Internal error)"));
                    } else {
                        future.completeExceptionally(new PlcException("Unexpected result " + response.getResult()));
                    }
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    protected PlcReadResponse convertToPlc4xReadResponse(PlcReadRequest readRequest, AmsPacket adsData) {
        ReadBuffer readBuffer = null;
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        if (adsData instanceof AdsReadResponse) {
            AdsReadResponse adsReadResponse = (AdsReadResponse) adsData;
            readBuffer = new ReadBufferByteBased(adsReadResponse.getData(), ByteOrder.LITTLE_ENDIAN);
            responseCodes.put(readRequest.getFieldNames().stream().findFirst().orElse(""),
                parsePlcResponseCode(adsReadResponse.getResult()));
        } else if (adsData instanceof AdsReadWriteResponse) {
            AdsReadWriteResponse adsReadWriteResponse = (AdsReadWriteResponse) adsData;
            readBuffer = new ReadBufferByteBased(adsReadWriteResponse.getData(), ByteOrder.LITTLE_ENDIAN);
            // When parsing a multi-item response, the error codes of each items come
            // in sequence and then come the values.
            for (String fieldName : readRequest.getFieldNames()) {
                try {
                    final ReturnCode result = ReturnCode.enumForValue(readBuffer.readUnsignedLong(32));
                    responseCodes.put(fieldName, parsePlcResponseCode(result));
                } catch (ParseException e) {
                    responseCodes.put(fieldName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }
        if (readBuffer != null) {
            Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
            for (String fieldName : readRequest.getFieldNames()) {
                DirectAdsField field;
                if (readRequest.getField(fieldName) instanceof DirectAdsField) {
                    field = (DirectAdsField) readRequest.getField(fieldName);
                } else {
                    SymbolicAdsField symbolicAdsField = (SymbolicAdsField) readRequest.getField(fieldName);
                    field = getDirectAdsFieldForSymbolicName(symbolicAdsField);
                }
                // If the response-code was anything but OK, we don't need to parse the payload.
                if (responseCodes.get(fieldName) != PlcResponseCode.OK) {
                    values.put(fieldName, new ResponseItem<>(responseCodes.get(fieldName), null));
                }
                // If the response-code was ok, parse the data returned.
                else {
                    values.put(fieldName, parseResponseItem(field, readBuffer));
                }
            }
            return new DefaultPlcReadResponse(readRequest, values);
        }
        return null;
    }

    private PlcResponseCode parsePlcResponseCode(ReturnCode adsResult) {
        if (adsResult == ReturnCode.OK) {
            return PlcResponseCode.OK;
        } else {
            // TODO: Implement this a little more ...
            return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    private ResponseItem<PlcValue> parseResponseItem(DirectAdsField field, ReadBuffer readBuffer) {
        try {
            String dataTypeName = field.getAdsDataTypeName();
            AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(dataTypeName);
            PlcValueType plcValueType = getPlcValueTypeForAdsDataType(adsDataTypeTableEntry);

            int strLen = 0;
            if ((plcValueType == PlcValueType.STRING) || (plcValueType == PlcValueType.WSTRING)) {
                strLen = (field instanceof AdsStringField) ? ((AdsStringField) field).getStringLength() : 256;
            }
            final int stringLength = strLen;
            if (field.getNumberOfElements() == 1) {
                return new ResponseItem<>(PlcResponseCode.OK, parsePlcValue(plcValueType, adsDataTypeTableEntry, stringLength, readBuffer));
            } else {
                // Fetch all
                final PlcValue[] resultItems = IntStream.range(0, field.getNumberOfElements()).mapToObj(i -> {
                    try {
                        return parsePlcValue(plcValueType, adsDataTypeTableEntry, stringLength, readBuffer);
                    } catch (ParseException e) {
                        LOGGER.warn("Error parsing field item of type: '{}' (at position {}})", field.getAdsDataTypeName(), i, e);
                    }
                    return null;
                }).toArray(PlcValue[]::new);
                return new ResponseItem<>(PlcResponseCode.OK, IEC61131ValueHandler.of(resultItems));
            }
        } catch (Exception e) {
            LOGGER.warn(String.format("Error parsing field item of type: '%s'", field.getAdsDataTypeName()), e);
            return new ResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null);
        }
    }

    private PlcValue parsePlcValue(PlcValueType plcValueType, AdsDataTypeTableEntry adsDataTypeTableEntry, int stringLength, ReadBuffer readBuffer) throws ParseException {
        switch (plcValueType) {
            case Struct:
                Map<String, PlcValue> properties = new HashMap<>();
                int startPos = readBuffer.getPos();
                int curPos = 0;
                for (AdsDataTypeTableChildEntry child : adsDataTypeTableEntry.getChildren()) {
                    if (child.getOffset() > curPos) {
                        long skipBytes = child.getOffset() - curPos;
                        for (long i = 0; i < skipBytes; i++) {
                            readBuffer.readByte();
                        }
                    }
                    String propertyName = child.getPropertyName();
                    AdsDataTypeTableEntry propertyDataTypeTableEntry = dataTypeTable.get(child.getDataTypeName());
                    PlcValueType propertyPlcValueType = getPlcValueTypeForAdsDataType(propertyDataTypeTableEntry);
                    PlcValue propertyValue = parsePlcValue(propertyPlcValueType, propertyDataTypeTableEntry, stringLength, readBuffer);
                    properties.put(propertyName, propertyValue);
                    curPos = readBuffer.getPos() - startPos;
                }
                return new PlcStruct(properties);
            case List:
                return parseArrayLevel(adsDataTypeTableEntry, adsDataTypeTableEntry.getArrayInfo(), readBuffer);
            default:
                return DataItem.staticParse(readBuffer, plcValueType, stringLength);
        }
    }

    private PlcValue parseArrayLevel(AdsDataTypeTableEntry adsDataTypeTableEntry, List<AdsDataTypeArrayInfo> arrayLayers, ReadBuffer readBuffer) throws ParseException {
        // If this is the last layer of the Array, parse the values themselves.
        if (arrayLayers.isEmpty()) {
            String dataTypeName = adsDataTypeTableEntry.getDataTypeName();
            dataTypeName = dataTypeName.substring(dataTypeName.lastIndexOf(" OF ") + 4);
            int stringLength = 0;
            if (dataTypeName.startsWith("STRING(")) {
                stringLength = Integer.parseInt(dataTypeName.substring(7, dataTypeName.length() - 1));
            } else if (dataTypeName.startsWith("WSTRING(")) {
                stringLength = Integer.parseInt(dataTypeName.substring(8, dataTypeName.length() - 1));
            }
            AdsDataTypeTableEntry elementDataTypeTableEntry = dataTypeTable.get(dataTypeName);
            PlcValueType plcValueType = getPlcValueTypeForAdsDataType(elementDataTypeTableEntry);
            return parsePlcValue(plcValueType, elementDataTypeTableEntry, stringLength, readBuffer);
        }

        List<PlcValue> elements = new ArrayList<>();
        List<AdsDataTypeArrayInfo> arrayInfo = adsDataTypeTableEntry.getArrayInfo();
        AdsDataTypeArrayInfo firstLayer = arrayInfo.get(0);
        for (int i = 0; i < firstLayer.getNumElements(); i++) {
            List<AdsDataTypeArrayInfo> remainingLayers = arrayInfo.subList(1, arrayInfo.size());
            elements.add(parseArrayLevel(adsDataTypeTableEntry, remainingLayers, readBuffer));
        }
        return new PlcList(elements);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        // Get all ADS addresses in their resolved state.
        final CompletableFuture<List<DirectAdsField>> directAdsFieldsFuture =
            getDirectAddresses(writeRequest.getFields());

        // If all addresses were already resolved we can send the request immediately.
        if (directAdsFieldsFuture.isDone()) {
            final List<DirectAdsField> fields = directAdsFieldsFuture.getNow(null);
            if (fields != null) {
                return executeWrite(writeRequest, fields);
            } else {
                final CompletableFuture<PlcWriteResponse> errorFuture = new CompletableFuture<>();
                errorFuture.completeExceptionally(new PlcException("Fields are null"));
                return errorFuture;
            }
        }
        // If there are still symbolic addresses that have to be resolved, send the
        // request as soon as the resolution is done.
        // In order to instantly be able to return a future, for the final result we have to
        // create a new one which is then completed later on. Unfortunately as soon as the
        // directAdsFieldsFuture is completed we still don't have the end result, but we can
        // now actually send the delayed read request ... as soon as that future completes
        // we can complete the initial one.
        else {
            CompletableFuture<PlcWriteResponse> delayedWrite = new CompletableFuture<>();
            directAdsFieldsFuture.handle((directAdsFields, throwable) -> {
                if (directAdsFields != null) {
                    final CompletableFuture<PlcWriteResponse> delayedResponse =
                        executeWrite(writeRequest, directAdsFields);
                    delayedResponse.handle((plcReadResponse, throwable1) -> {
                        if (plcReadResponse != null) {
                            delayedWrite.complete(plcReadResponse);
                        } else {
                            delayedWrite.completeExceptionally(throwable1);
                        }
                        return this;
                    });
                } else {
                    delayedWrite.completeExceptionally(throwable);
                }
                return this;
            });
            return delayedWrite;
        }
    }

    protected CompletableFuture<PlcWriteResponse> executeWrite(PlcWriteRequest writeRequest,
                                                               List<DirectAdsField> directAdsFields) {
        // Depending on the number of fields, use a single item request or a sum-request
        if (directAdsFields.size() == 1) {
            // Do a normal (single item) ADS Write Request
            return singleWrite(writeRequest, directAdsFields.get(0));
        } else {
            // TODO: Check if the version of the remote station is at least TwinCAT v2.11 Build >= 1550 otherwise split up into single item requests.
            // Do a ADS-Sum Read Request.
            return multiWrite(writeRequest, directAdsFields);
        }
    }

    protected CompletableFuture<PlcWriteResponse> singleWrite(PlcWriteRequest writeRequest, DirectAdsField directAdsField) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();

/*        final String fieldName = writeRequest.getFieldNames().iterator().next();
        final AdsField plcField = (AdsField) writeRequest.getField(fieldName);
        final PlcValue plcValue = writeRequest.getPlcValue(fieldName);
        final int stringLength;
        if (directAdsField.getAdsDataType() == AdsDataType.STRING) {
            stringLength = plcValue.getString().length() + 1;
        } else {
            if (directAdsField.getAdsDataType() == AdsDataType.WSTRING) {
                stringLength = (plcValue.getString().length() + 1) * 2;
            } else {
                stringLength = 0;
            }
        }
        try {
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(DataItem.getLengthInBytes(plcValue,
                plcField.getAdsDataType().getPlcValueType(), stringLength));
            DataItem.staticSerialize(writeBuffer, plcValue, plcField.getAdsDataType().getPlcValueType(), stringLength, ByteOrder.LITTLE_ENDIAN);
            AmsPacket amsPacket = new AdsWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
                0, getInvokeId(), directAdsField.getIndexGroup(), directAdsField.getIndexOffset(), writeBuffer.getData());
            AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

            // Start a new request-transaction (Is ended in the response-handler)
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            transaction.submit(() -> context.sendRequest(amsTCPPacket)
                .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
                .unwrap(amsResponsePacket -> (AdsWriteResponse) amsResponsePacket.getUserdata())
                .handle(response -> {
                    if (response.getResult() == ReturnCode.OK) {
                        final PlcWriteResponse plcWriteResponse = convertToPlc4xWriteResponse(writeRequest, response);
                        // Convert the response from the PLC into a PLC4X Response ...
                        future.complete(plcWriteResponse);
                    } else {
                        // TODO: Implement this correctly.
                        future.completeExceptionally(new PlcException("Unexpected return code " + response.getResult()));
                    }
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        } catch (Exception e) {
            future.completeExceptionally(new PlcException("Error"));
        }*/
        return future;
    }

    protected CompletableFuture<PlcWriteResponse> multiWrite(PlcWriteRequest writeRequest, List<DirectAdsField> directAdsFields) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();

        // Calculate the size of all fields together.
        // Calculate the expected size of the response data.
/*        int expectedRequestDataSize = directAdsFields.stream().mapToInt(
            field -> field.getAdsDataType().getNumBytes() * field.getNumberOfElements()).sum();
        byte[] writeBuffer = new byte[expectedRequestDataSize];
        int pos = 0;
        for (String fieldName : writeRequest.getFieldNames()) {
            final AdsField field = (AdsField) writeRequest.getField(fieldName);
            final PlcValue plcValue = writeRequest.getPlcValue(fieldName);
            final int stringLength;
            if (field.getAdsDataType() == AdsDataType.STRING) {
                stringLength = plcValue.getString().length() + 1;
            } else {
                if (field.getAdsDataType() == AdsDataType.WSTRING) {
                    stringLength = (plcValue.getString().length() + 1) * 2;
                } else {
                    stringLength = 0;
                }
            }
            try {
                WriteBufferByteBased itemWriteBuffer = new WriteBufferByteBased(DataItem.getLengthInBytes(plcValue,
                    field.getAdsDataType().getPlcValueType(), stringLength));
                DataItem.staticSerialize(itemWriteBuffer, plcValue,
                    field.getAdsDataType().getPlcValueType(), stringLength, ByteOrder.LITTLE_ENDIAN);
                int numBytes = itemWriteBuffer.getPos();
                System.arraycopy(itemWriteBuffer.getData(), 0, writeBuffer, pos, numBytes);
                pos += numBytes;
            } catch (Exception e) {
                throw new PlcRuntimeException("Error serializing data", e);
            }
        }

        // With multi-requests, the index-group is fixed and the index offset indicates the number of elements.
        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_MULTIPLE_WRITE.getValue(), directAdsFields.size(), (long) directAdsFields.size() * 4,
            directAdsFields.stream().map(directAdsField -> new AdsMultiRequestItemWrite(
                    directAdsField.getIndexGroup(), directAdsField.getIndexOffset(),
                    ((long) directAdsField.getAdsDataType().getNumBytes() * directAdsField.getNumberOfElements())))
                .collect(Collectors.toList()), writeBuffer);
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
            .unwrap(responseAmsPacket -> (AdsReadWriteResponse) responseAmsPacket.getUserdata())
            .handle(response -> {
                if (response.getResult() == ReturnCode.OK) {
                    final PlcWriteResponse plcWriteResponse = convertToPlc4xWriteResponse(writeRequest, response);
                    // Convert the response from the PLC into a PLC4X Response ...
                    future.complete(plcWriteResponse);
                } else {
                    // TODO: Implement this correctly.
                    future.completeExceptionally(new PlcException("Error"));
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));*/
        return future;
    }

    protected PlcWriteResponse convertToPlc4xWriteResponse(PlcWriteRequest writeRequest, AmsPacket adsData) {
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        if (adsData instanceof AdsWriteResponse) {
            AdsWriteResponse adsWriteResponse = (AdsWriteResponse) adsData;
            responseCodes.put(writeRequest.getFieldNames().stream().findFirst().orElse(""),
                parsePlcResponseCode(adsWriteResponse.getResult()));
        } else if (adsData instanceof AdsReadWriteResponse) {
            AdsReadWriteResponse adsReadWriteResponse = (AdsReadWriteResponse) adsData;
            ReadBuffer readBuffer = new ReadBufferByteBased(adsReadWriteResponse.getData(), ByteOrder.LITTLE_ENDIAN);
            // When parsing a multi-item response, the error codes of each items come
            // in sequence and then come the values.
            for (String fieldName : writeRequest.getFieldNames()) {
                try {
                    final ReturnCode result = ReturnCode.enumForValue(readBuffer.readUnsignedLong(32));
                    responseCodes.put(fieldName, parsePlcResponseCode(result));
                } catch (ParseException e) {
                    responseCodes.put(fieldName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }

        return new DefaultPlcWriteResponse(writeRequest, responseCodes);
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        // Get all ADS addresses in their resolved state.
        final CompletableFuture<List<DirectAdsField>> directAdsFieldsFuture =
            getDirectAddresses(subscriptionRequest.getFields()
                .stream()
                .map(field -> ((DefaultPlcSubscriptionField) field).getPlcField())
                .collect(Collectors.toList()));

        // If all addresses were already resolved we can send the request immediately.
        if (directAdsFieldsFuture.isDone()) {
            final List<DirectAdsField> fields = directAdsFieldsFuture.getNow(null);
            if (fields != null) {
                return executeSubscribe(subscriptionRequest);
            } else {
                final CompletableFuture<PlcSubscriptionResponse> errorFuture = new CompletableFuture<>();
                errorFuture.completeExceptionally(new PlcException("Fields are null"));
                return errorFuture;
            }
        }
        // If there are still symbolic addresses that have to be resolved, send the
        // request as soon as the resolution is done.
        // In order to instantly be able to return a future, for the final result we have to
        // create a new one which is then completed later on. Unfortunately as soon as the
        // directAdsFieldsFuture is completed we still don't have the end result, but we can
        // now actually send the delayed read request ... as soon as that future completes
        // we can complete the initial one.
        else {
            CompletableFuture<PlcSubscriptionResponse> delayedSubscribe = new CompletableFuture<>();
            directAdsFieldsFuture.handle((directAdsFields, throwable) -> {
                if (directAdsFields != null) {
                    final CompletableFuture<PlcSubscriptionResponse> delayedResponse =
                        executeSubscribe(subscriptionRequest);
                    delayedResponse.handle((plcSubscribeResponse, throwable1) -> {
                        if (plcSubscribeResponse != null) {
                            delayedSubscribe.complete(plcSubscribeResponse);
                        } else {
                            delayedSubscribe.completeExceptionally(throwable1);
                        }
                        return this;
                    });
                } else {
                    delayedSubscribe.completeExceptionally(throwable);
                }
                return this;
            });
            return delayedSubscribe;
        }
    }

    private CompletableFuture<PlcSubscriptionResponse> executeSubscribe(PlcSubscriptionRequest subscribeRequest) {
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();

        List<AmsTCPPacket> amsTCPPackets = subscribeRequest.getFields().stream()
            .map(field -> (DefaultPlcSubscriptionField) field)
            .map(field -> {
                AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(((DirectAdsField) field.getPlcField()).getAdsDataTypeName());
                DirectAdsField directAdsField = getDirectAdsFieldForSymbolicName(field);
                return new AmsTCPPacket(new AdsAddDeviceNotificationRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                    configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
                    0, getInvokeId(),
                    directAdsField.getIndexGroup(),
                    directAdsField.getIndexOffset(),
                    adsDataTypeTableEntry.getSize() * field.getNumberOfElements(),
                    field.getPlcSubscriptionType() == PlcSubscriptionType.CYCLIC ? 3 : 4, // if it's not cyclic, it's on change or event
                    0, // there is no api for that yet
                    field.getDuration().orElse(Duration.ZERO).toMillis()));
            })
            .collect(Collectors.toList());

        Map<String, ResponseItem<PlcSubscriptionHandle>> responses = new HashMap<>();

        // Start the first request-transaction (it is ended in the response-handler).
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(subscribeRecursively(
            subscribeRequest,
            subscribeRequest.getFieldNames().iterator(),
            responses,
            future,
            amsTCPPackets.iterator(),
            transaction));
        return future;
    }

    private Runnable subscribeRecursively(PlcSubscriptionRequest subscriptionRequest, Iterator<String> fieldNames,
                                          Map<String, ResponseItem<PlcSubscriptionHandle>> responses,
                                          CompletableFuture<PlcSubscriptionResponse> future,
                                          Iterator<AmsTCPPacket> amsTCPPackets,
                                          RequestTransactionManager.RequestTransaction transaction) {
        return () -> {
            AmsTCPPacket packet = amsTCPPackets.next();
            boolean hasMorePackets = amsTCPPackets.hasNext();
            String fieldName = fieldNames.next();
            context.sendRequest(packet)
                .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == packet.getUserdata().getInvokeId())
                .unwrap(responseAmsPacket -> (AdsAddDeviceNotificationResponse) responseAmsPacket.getUserdata())
                .handle(response -> {
                    if (response.getResult() == ReturnCode.OK) {
                        DefaultPlcSubscriptionField subscriptionField = (DefaultPlcSubscriptionField) subscriptionRequest.getField(fieldName);
                        AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(((DirectAdsField) subscriptionField.getPlcField()).getAdsDataTypeName());

                        // Collect notification handle from individual response.
                        responses.put(fieldName, new ResponseItem<>(
                            parsePlcResponseCode(response.getResult()),
                            new AdsSubscriptionHandle(this,
                                fieldName,
                                adsDataTypeTableEntry,
                                response.getNotificationHandle())));

                        // After receiving the last ADD_DEVICE_NOTIFICATION response, complete the PLC4X response.
                        if (!hasMorePackets) {
                            final PlcSubscriptionResponse plcSubscriptionResponse = new DefaultPlcSubscriptionResponse(subscriptionRequest, responses);
                            future.complete(plcSubscriptionResponse);
                        }
                    } else {
                        if (response.getResult() == ReturnCode.ADSERR_DEVICE_INVALIDSIZE) {
                            future.completeExceptionally(
                                new PlcException("The parameter size was not correct (Internal error)"));
                        } else {
                            future.completeExceptionally(new PlcException("Unexpected result " + response.getResult()));
                        }
                    }
                    // Finish the request-transaction.
                    transaction.endRequest();

                    // Submit the next transaction.
                    if (hasMorePackets) {
                        RequestTransactionManager.RequestTransaction nextTransaction = tm.startRequest();
                        nextTransaction.submit(subscribeRecursively(
                            subscriptionRequest, fieldNames, responses, future, amsTCPPackets, nextTransaction));
                    }
                });
        };
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<>();

        List<Long> notificationHandles = new ArrayList<>();
        unsubscriptionRequest.getSubscriptionHandles().stream()
            .filter(handle -> handle instanceof AdsSubscriptionHandle)
            .map(handle -> (AdsSubscriptionHandle) handle)
            .forEach(adsSubscriptionHandle -> {
                // Notification handle used for delete notification messages.
                notificationHandles.add(adsSubscriptionHandle.getNotificationHandle());
                // Remove consumers
                consumers.keySet().stream().filter(consumerRegistration ->
                        consumerRegistration.getSubscriptionHandles().contains(adsSubscriptionHandle))
                    .forEach(DefaultPlcConsumerRegistration::unregister);
            });

        List<AmsTCPPacket> amsTCPPackets = notificationHandles.stream().map(data -> new AmsTCPPacket(
            new AdsDeleteDeviceNotificationRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
                0, getInvokeId(), data))).collect(Collectors.toList());

        // Start the first request-transaction (it is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(unsubscribeRecursively(unsubscriptionRequest, future, amsTCPPackets.iterator(), transaction));
        return future;
    }

    private Runnable unsubscribeRecursively(PlcUnsubscriptionRequest unsubscriptionRequest,
                                            CompletableFuture<PlcUnsubscriptionResponse> future,
                                            Iterator<AmsTCPPacket> amsTCPPackets,
                                            RequestTransactionManager.RequestTransaction transaction) {
        return () -> {
            AmsTCPPacket packet = amsTCPPackets.next();
            boolean hasMorePackets = amsTCPPackets.hasNext();
            context.sendRequest(packet)
                .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == packet.getUserdata().getInvokeId())
                .unwrap(responseAmsPacket -> (AdsDeleteDeviceNotificationResponse) responseAmsPacket.getUserdata())
                .handle(response -> {
                    if (response.getResult() == ReturnCode.OK) {
                        // After receiving the last DELETE_DEVICE_NOTIFICATION response, complete the PLC4X response.
                        if (!hasMorePackets) {
                            final PlcUnsubscriptionResponse plcUnsubscriptionResponse = new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest);
                            future.complete(plcUnsubscriptionResponse);
                        }
                    } else {
                        // TODO: this is more guesswork than knowing it could actually occur
                        if (response.getResult() == ReturnCode.ADSERR_DEVICE_NOTIFYHNDINVALID) {
                            future.completeExceptionally(
                                new PlcException("The notification handle is invalid (Internal error)"));
                        } else {
                            future.completeExceptionally(new PlcException("Unexpected result " + response.getResult()));
                        }
                    }
                    // Finish the request-transaction.
                    transaction.endRequest();

                    // Submit the next transaction.
                    if (hasMorePackets) {
                        RequestTransactionManager.RequestTransaction nextTransaction = tm.startRequest();
                        nextTransaction.submit(unsubscribeRecursively(unsubscriptionRequest, future, amsTCPPackets, nextTransaction));
                    }
                });
        };
    }

    @Override
    protected void decode(ConversationContext<AmsTCPPacket> context, AmsTCPPacket msg) throws Exception {
        if (msg.getUserdata() instanceof AdsDeviceNotificationRequest) {
            AdsDeviceNotificationRequest notificationData = (AdsDeviceNotificationRequest) msg.getUserdata();
            List<AdsStampHeader> stamps = notificationData.getAdsStampHeaders();
            for (AdsStampHeader stamp : stamps) {
                // convert Windows FILETIME format to unix epoch
                long unixEpochTimestamp = stamp.getTimestamp().divide(BigInteger.valueOf(10000L)).longValue() - 11644473600000L;
                List<AdsNotificationSample> samples = stamp.getAdsNotificationSamples();
                for (AdsNotificationSample sample : samples) {
                    long handle = sample.getNotificationHandle();
                    for (DefaultPlcConsumerRegistration registration : consumers.keySet()) {
                        for (PlcSubscriptionHandle subscriptionHandle : registration.getSubscriptionHandles()) {
                            if (subscriptionHandle instanceof AdsSubscriptionHandle) {
                                AdsSubscriptionHandle adsHandle = (AdsSubscriptionHandle) subscriptionHandle;
                                if (adsHandle.getNotificationHandle() == handle)
                                    consumers.get(registration).accept(
                                        new DefaultPlcSubscriptionEvent(Instant.ofEpochMilli(unixEpochTimestamp),
                                            convertSampleToPlc4XResult(adsHandle, sample.getData())));
                            }
                        }
                    }
                }
            }
        }
    }

    private Map<String, ResponseItem<PlcValue>> convertSampleToPlc4XResult(AdsSubscriptionHandle subscriptionHandle, byte[] data) throws
        ParseException {
        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN);
        values.put(subscriptionHandle.getPlcFieldName(), new ResponseItem<>(PlcResponseCode.OK,
            DataItem.staticParse(readBuffer, getPlcValueTypeForAdsDataType(subscriptionHandle.getAdsDataType()), data.length)));
        return values;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        final DefaultPlcConsumerRegistration consumerRegistration =
            new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(consumerRegistration, consumer);
        return consumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        DefaultPlcConsumerRegistration consumerRegistration = (DefaultPlcConsumerRegistration) registration;
        consumers.remove(consumerRegistration);
    }

    protected CompletableFuture<List<DirectAdsField>> getDirectAddresses(List<PlcField> fields) {
        CompletableFuture<List<DirectAdsField>> future = new CompletableFuture<>();

        // Get all symbolic fields from the current request.
        // These potentially need to be resolved to direct addresses, if this has not been done before.
        final List<SymbolicAdsField> referencedSymbolicFields = fields.stream()
            .filter(SymbolicAdsField.class::isInstance)
            .map(SymbolicAdsField.class::cast)
            .collect(Collectors.toList());

        // Find out for which of these symbolic addresses no resolution has been initiated.
        final List<SymbolicAdsField> symbolicFieldsNeedingResolution = referencedSymbolicFields.stream()
            .filter(symbolicAdsField -> getDirectAdsFieldForSymbolicName(symbolicAdsField) == null)
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
                    resolutionFuture = resolveSingleSymbolicAddress(requiredResolutionFields.get(0));
                    pendingResolutionRequests.put(symbolicAdsField, resolutionFuture);
                } else {
                    resolutionFuture = resolveMultipleSymbolicAddresses(requiredResolutionFields);
                    for (SymbolicAdsField symbolicAdsField : requiredResolutionFields) {
                        pendingResolutionRequests.put(symbolicAdsField, resolutionFuture);
                    }
                }
            }

            // Create a global future which is completed as soon as all sub-futures for this request are completed.
            final CompletableFuture<Void> resolutionComplete =
                CompletableFuture.allOf(symbolicFieldsNeedingResolution.stream()
                    .map(pendingResolutionRequests::get)
                    .toArray(CompletableFuture[]::new));

            // Complete the future asynchronously as soon as all fields are resolved.
            resolutionComplete.handleAsync((unused, throwable) -> {
                if (throwable != null) {
                    return future.completeExceptionally(throwable.getCause());
                } else {
                    List<DirectAdsField> directAdsFields = new ArrayList<>(fields.size());
                    for (PlcField field : fields) {
                        if (field instanceof SymbolicAdsField) {
                            directAdsFields.add(getDirectAdsFieldForSymbolicName(field));
                        } else {
                            directAdsFields.add((DirectAdsField) field);
                        }
                    }
                    return future.complete(directAdsFields);
                }
            });
        } else {
            // If all fields were resolved, we can continue instantly.
            future.complete(fields.stream().map(plcField -> {
                if (plcField instanceof SymbolicAdsField) {
                    return getDirectAdsFieldForSymbolicName(plcField);
                } else {
                    return (DirectAdsField) plcField;
                }
            }).collect(Collectors.toList()));
        }
        return future;
    }

    protected CompletableFuture<Void> resolveSingleSymbolicAddress(SymbolicAdsField symbolicAdsField) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_SYM_HNDBYNAME.getValue(), 0,
            4, null,
            getNullByteTerminatedArray(symbolicAdsField.getSymbolicAddress()));
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
            .unwrap(AmsTCPPacket::getUserdata)
            .check(AdsReadWriteResponse.class::isInstance)
            .unwrap(AdsReadWriteResponse.class::cast)
            .handle(response -> {
                if (response.getResult() != ReturnCode.OK) {
                    future.completeExceptionally(new PlcException("Couldn't retrieve handle for symbolic field " +
                        symbolicAdsField.getSymbolicAddress() + " got return code " + response.getResult().name()));
                } else {
                    ReadBuffer readBuffer = new ReadBufferByteBased(response.getData(), ByteOrder.LITTLE_ENDIAN);
                    try {
                        // Read the handle.
                        long handle = readBuffer.readUnsignedLong(32);

/*                        DirectAdsField directAdsField = new DirectAdsField(
                            ReservedIndexGroups.ADSIGRP_SYM_VALBYHND.getValue(), handle,
                            symbolicAdsField.getAdsDataTypeName(), symbolicAdsField.getNumberOfElements());*/

                        // TODO: Find out how to read the datatype for the given symbolic field
                        //symbolicFieldMapping.put(symbolicAdsField, directAdsField);
                        future.complete(null);
                    } catch (ParseException e) {
                        future.completeExceptionally(e);
                    }
                }
                transaction.endRequest();
            }));
        return future;
    }

    protected CompletableFuture<Void> resolveMultipleSymbolicAddresses(List<SymbolicAdsField> symbolicAdsFields) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        // The expected response for every symbolic address is 12 bytes (8 bytes header and 4 bytes for the handle)
        long expectedResponseDataSize = (long) (symbolicAdsFields.size()) * 12;
        // Concatenate the string part of each symbolic address into one concatenated string and get the bytes.
        byte[] addressData = symbolicAdsFields.stream().map(
            SymbolicAdsField::getSymbolicAddress).collect(Collectors.joining("")).getBytes();
        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_MULTIPLE_READ_WRITE.getValue(),
            symbolicAdsFields.size(), expectedResponseDataSize, symbolicAdsFields.stream().map(symbolicAdsField ->
            new AdsMultiRequestItemReadWrite(ReservedIndexGroups.ADSIGRP_SYM_HNDBYNAME.getValue(), 0,
                4, symbolicAdsField.getSymbolicAddress().length())).collect(Collectors.toList()), null);
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> context.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
            .unwrap(AmsTCPPacket::getUserdata)
            .check(AdsReadWriteResponse.class::isInstance)
            .unwrap(AdsReadWriteResponse.class::cast)
            .handle(response -> {
                ReadBuffer readBuffer = new ReadBufferByteBased(response.getData(), ByteOrder.LITTLE_ENDIAN);
                Map<SymbolicAdsField, Long> returnCodes = new HashMap<>();
                // In the response first come the return codes and the data-lengths for each item.
                symbolicAdsFields.forEach(symbolicAdsField -> {
                    try {
                        // This should be 0 in the success case.
                        long returnCode = readBuffer.readUnsignedLong(32);
                        // This is always 4
                        long itemLength = readBuffer.readUnsignedLong(32);
                        assert itemLength == 4;

                        returnCodes.put(symbolicAdsField, returnCode);
                    } catch (ParseException e) {
                        throw new PlcRuntimeException(e);
                    }
                });
                // After reading the header-information, comes the data itself.
                symbolicAdsFields.forEach(symbolicAdsField -> {
                    try {
                        if (returnCodes.get(symbolicAdsField) == 0) {
                            // Read the handle.
                            long handle = readBuffer.readUnsignedLong(32);

                            /*DirectAdsField directAdsField = new DirectAdsField(
                                ReservedIndexGroups.ADSIGRP_SYM_VALBYHND.getValue(), handle,
                                symbolicAdsField.getAdsDataTypeName(), symbolicAdsField.getNumberOfElements());*/
                            // TODO: Find out how to read the datatype for the given symbolic field
                            //symbolicFieldMapping.put(symbolicAdsField, directAdsField);
                        } else {
                            // TODO: Handle the case of unsuccessful resolution ..
                        }
                    } catch (ParseException e) {
                        throw new PlcRuntimeException(e);
                    }
                });
                future.complete(null);
                transaction.endRequest();
            }));
        return future;
    }

    protected long getInvokeId() {
        long invokeId = invokeIdGenerator.getAndIncrement();
        // If we've reached the max value for a 16 bit transaction identifier, reset back to 1
        if (invokeIdGenerator.get() == 0xFFFFFFFF) {
            invokeIdGenerator.set(1);
        }
        return invokeId;
    }

    protected DirectAdsField getDirectAdsFieldForSymbolicName(PlcField field) {
        if (field instanceof DirectAdsField) {
            return (DirectAdsField) field;
        }

        SymbolicAdsField symbolicAdsField = (SymbolicAdsField) field;
        String symbolicAddress = symbolicAdsField.getSymbolicAddress();
        String[] addressParts = symbolicAddress.split("\\.");

        // If the number of parts are less than 2, we can find the entry in the symbol table directly.
        if (addressParts.length < 2) {
            // We can't find it, so we need to resolve it.
            if (!symbolTable.containsKey(symbolicAddress)) {
                return null;
            }

            AdsSymbolTableEntry adsSymbolTableEntry = symbolTable.get(symbolicAddress);
            AdsDataTypeTableEntry dataTypeTableEntry = dataTypeTable.get(adsSymbolTableEntry.getDataTypeName());
            return new DirectAdsField(adsSymbolTableEntry.getGroup(), adsSymbolTableEntry.getOffset(),
                dataTypeTableEntry.getDataTypeName(), dataTypeTableEntry.getArrayDimensions());
        }
        // Otherwise we'll have to crawl through the dataType definitions.
        else {
            String symbolName = addressParts[0] + "." + addressParts[1];
            AdsSymbolTableEntry adsSymbolTableEntry = symbolTable.get(symbolName);
            AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(adsSymbolTableEntry.getDataTypeName());
            return resolveDirectAdsFieldForSymbolicNameFromDataType(
                Arrays.asList(addressParts).subList(2, addressParts.length),
                adsSymbolTableEntry.getGroup(), adsSymbolTableEntry.getOffset(), adsDataTypeTableEntry);
        }
    }

    protected DirectAdsField resolveDirectAdsFieldForSymbolicNameFromDataType(List<String> remainingAddressParts, long currentGroup, long currentOffset, AdsDataTypeTableEntry adsDataTypeTableEntry) {
        if (remainingAddressParts.isEmpty()) {
            // TODO: Implement the Array support
            return new DirectAdsField(currentGroup, currentOffset, adsDataTypeTableEntry.getDataTypeName(), null);
        }

        // Go through all children looking for a matching one.
        for (AdsDataTypeTableChildEntry child : adsDataTypeTableEntry.getChildren()) {
            if (child.getPropertyName().equals(remainingAddressParts.get(0))) {
                AdsDataTypeTableEntry childAdsDataTypeTableEntry = dataTypeTable.get(child.getDataTypeName());
                return resolveDirectAdsFieldForSymbolicNameFromDataType(
                    remainingAddressParts.subList(1, remainingAddressParts.size()),
                    currentGroup, currentOffset + child.getOffset(), childAdsDataTypeTableEntry);
            }
        }

        throw new PlcRuntimeException(String.format("Couldn't find child with name '%s' for type '%s'",
            remainingAddressParts.get(0), adsDataTypeTableEntry.getDataTypeName()));
    }

    protected PlcValueType getPlcValueTypeForAdsDataType(AdsDataTypeTableEntry dataTypeTableEntry) {
        String dataTypeName = dataTypeTableEntry.getDataTypeName();
        if (dataTypeName.startsWith("STRING(")) {
            dataTypeName = "STRING";
        } else if (dataTypeName.startsWith("WSTRING(")) {
            dataTypeName = "WSTRING";
        }
        // First check, if this is a primitive type.
        try {
            return PlcValueType.valueOf(dataTypeName);
        } catch (IllegalArgumentException e) {
            // Then check if this is an array.
            if (dataTypeTableEntry.getArrayDimensions() > 0) {
                return PlcValueType.List;
            }
            // There seem to be some data types, that have odd names, but no children
            // So we'll check if their "simpleTypeName" matches instead.
            if(dataTypeTableEntry.getChildren().isEmpty()) {
                try {
                    dataTypeName = dataTypeTableEntry.getSimpleTypeName();
                    if (dataTypeName.startsWith("STRING(")) {
                        dataTypeName = "STRING";
                    } else if (dataTypeName.startsWith("WSTRING(")) {
                        dataTypeName = "WSTRING";
                    }

                    return PlcValueType.valueOf(dataTypeName);
                } catch (IllegalArgumentException e2) {
                    // In this case it's something we can't handle.
                    return PlcValueType.NULL;
                }
            }
            return PlcValueType.Struct;
        }
    }

    protected byte[] getNullByteTerminatedArray(String value) {
        byte[] valueBytes = value.getBytes();
        byte[] nullTerminatedBytes = new byte[valueBytes.length + 1];
        System.arraycopy(valueBytes, 0, nullTerminatedBytes, 0, valueBytes.length);
        return nullTerminatedBytes;
    }

}
