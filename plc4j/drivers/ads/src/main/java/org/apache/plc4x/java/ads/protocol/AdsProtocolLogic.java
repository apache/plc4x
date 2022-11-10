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
import org.apache.plc4x.java.ads.model.AdsSubscriptionHandle;
import org.apache.plc4x.java.ads.readwrite.*;
import org.apache.plc4x.java.ads.readwrite.DataItem;
import org.apache.plc4x.java.ads.tag.AdsTag;
import org.apache.plc4x.java.ads.tag.DirectAdsTag;
import org.apache.plc4x.java.ads.tag.SymbolicAdsTag;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.net.*;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.Consumer;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

public class AdsProtocolLogic extends Plc4xProtocolBase<AmsTCPPacket> implements HasConfiguration<AdsConfiguration>, PlcSubscriber, PlcBrowser {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsProtocolLogic.class);

    private AdsConfiguration configuration;

    private String adsVersion;
    private String deviceName;

    private final AtomicLong invokeIdGenerator = new AtomicLong(1);
    private final RequestTransactionManager tm;

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    private final ConcurrentHashMap<SymbolicAdsTag, CompletableFuture<Void>> pendingResolutionRequests;

    private int symbolVersion;
    private long onlineVersion;
    private final Map<String, AdsSymbolTableEntry> symbolTable;
    private final Map<String, AdsDataTypeTableEntry> dataTypeTable;
    private final ReentrantLock invalidationLock;

    public AdsProtocolLogic() {
//        symbolicTagMapping = new ConcurrentHashMap<>();
        pendingResolutionRequests = new ConcurrentHashMap<>();
        symbolTable = new HashMap<>();
        dataTypeTable = new HashMap<>();
        invalidationLock = new ReentrantLock();

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
            if (!configuration.isLoadSymbolAndDataTypeTables()) {
                future.completeExceptionally(new PlcConnectionException(
                    "Lazy loading is generally planned, but not implemented yet. " +
                        "If you are in need for this feature, please reach out to the community."));
            }
            //if (configuration.isLoadSymbolAndDataTypeTables()) {
                // Execute a ReadDeviceInfo command
                AmsPacket readDeviceInfoRequest = new AdsReadDeviceInfoRequest(
                    configuration.getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
                    configuration.getSourceAmsNetId(), 800, 0, getInvokeId());
                RequestTransactionManager.RequestTransaction readDeviceInfoTx = tm.startRequest();
                readDeviceInfoTx.submit(() -> context.sendRequest(new AmsTCPPacket(readDeviceInfoRequest))
                    .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                    .onTimeout(future::completeExceptionally)
                    .onError((p, e) -> future.completeExceptionally(e))
                    .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == readDeviceInfoRequest.getInvokeId())
                    .unwrap(response -> (AdsReadDeviceInfoResponse) response.getUserdata())
                    .handle(readDeviceInfoResponse -> {
                        readDeviceInfoTx.endRequest();
                        if (readDeviceInfoResponse.getResult() != ReturnCode.OK) {
                            // TODO: Handle this
                            future.completeExceptionally(new PlcException("Result is " + readDeviceInfoResponse.getResult()));
                            return;
                        }

                        // Get the twin-cat version and PLC name.
                        adsVersion = String.format("%d.%d.%d", readDeviceInfoResponse.getMajorVersion(),
                            readDeviceInfoResponse.getMinorVersion(), readDeviceInfoResponse.getVersion());
                        deviceName = new String(readDeviceInfoResponse.getDevice()).trim();

                        // Read the online version number (Address; GroupID: 0xF004 (read symbol by name),Offset: 0, Read length: 4, ... Payload: "TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt")
                        AmsPacket readOnlineVersionNumberRequest = new AdsReadWriteRequest(
                            configuration.getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
                            configuration.getSourceAmsNetId(), 800, 0, getInvokeId(),
                            ReservedIndexGroups.ADSIGRP_SYM_VALBYNAME.getValue(), 0, 4, null,
                            "TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt".getBytes(StandardCharsets.UTF_8));
                        RequestTransactionManager.RequestTransaction readOnlineVersionNumberTx = tm.startRequest();
                        readOnlineVersionNumberTx.submit(() -> context.sendRequest(new AmsTCPPacket(readOnlineVersionNumberRequest))
                            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                            .onTimeout(future::completeExceptionally)
                            .onError((p, e) -> future.completeExceptionally(e))
                            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == readOnlineVersionNumberRequest.getInvokeId())
                            .unwrap(response -> (AdsReadWriteResponse) response.getUserdata())
                            .handle(readOnlineVersionNumberResponse -> {
                                readOnlineVersionNumberTx.endRequest();
                                if (readOnlineVersionNumberResponse.getResult() != ReturnCode.OK) {
                                    // TODO: Handle this
                                    future.completeExceptionally(new PlcException("Result is " + readOnlineVersionNumberResponse.getResult()));
                                    return;
                                }
                                try {
                                    ReadBuffer rb = new ReadBufferByteBased(readOnlineVersionNumberResponse.getData());
                                    onlineVersion = rb.readUnsignedLong(32);

                                    // Read the offline version number (Address: GroupID: 0xF008, Offset: 0, Read length: 1)
                                    AmsPacket readSymbolVersionNumberRequest = new AdsReadRequest(
                                        configuration.getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
                                        configuration.getSourceAmsNetId(), 800, 0, getInvokeId(),
                                        ReservedIndexGroups.ADSIGRP_SYM_VERSION.getValue(), 0, 1);
                                    RequestTransactionManager.RequestTransaction readSymbolVersionNumberTx = tm.startRequest();
                                    readSymbolVersionNumberTx.submit(() -> context.sendRequest(new AmsTCPPacket(readSymbolVersionNumberRequest))
                                        .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                                        .onTimeout(future::completeExceptionally)
                                        .onError((p, e) -> future.completeExceptionally(e))
                                        .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == readSymbolVersionNumberRequest.getInvokeId())
                                        .unwrap(response -> (AdsReadResponse) response.getUserdata())
                                        .handle(readSymbolVersionNumberResponse -> {
                                            readSymbolVersionNumberTx.endRequest();
                                            if (readSymbolVersionNumberResponse.getResult() != ReturnCode.OK) {
                                                // TODO: Handle this
                                                future.completeExceptionally(new PlcException("Result is " + readSymbolVersionNumberResponse.getResult()));
                                                return;
                                            }
                                            try {
                                                ReadBuffer rb2 = new ReadBufferByteBased(readSymbolVersionNumberResponse.getData());
                                                symbolVersion = rb2.readUnsignedInt(8);

                                                LOGGER.debug("Fetching sizes of symbol and datatype table sizes.");
                                                CompletableFuture<Void> readSymbolTableFuture = readSymbolTableAndDatatypeTable(context);
                                                readSymbolTableFuture.whenComplete((unused2, throwable2) -> {
                                                    if (throwable2 != null) {
                                                        LOGGER.error("Error fetching symbol and datatype table sizes");
                                                    } else {
                                                        context.fireConnected();
                                                    }
                                                });
                                            } catch (ParseException e) {
                                                future.completeExceptionally(new PlcConnectionException("Error reading the symbol version of data type and symbol data.", e));
                                            }
                                        }));
                                } catch (ParseException e) {
                                    future.completeExceptionally(new PlcConnectionException("Error reading the online version of data type and symbol data.", e));
                                }
                            }));
                    }));
            /*} else {
                context.fireConnected();
            }*/
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

        // Read the data-type and symbol table sizes
        AmsPacket readDataAndSymbolTableSizesRequest = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
            ReservedIndexGroups.ADSIGRP_SYMBOL_AND_DATA_TYPE_SIZES.getValue(), 0x00000000, 24);
        RequestTransactionManager.RequestTransaction readDataAndSymbolTableSizesTx = tm.startRequest();
        readDataAndSymbolTableSizesTx.submit(() -> context.sendRequest(new AmsTCPPacket(readDataAndSymbolTableSizesRequest))
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == readDataAndSymbolTableSizesRequest.getInvokeId())
            .unwrap(response -> (AdsReadResponse) response.getUserdata())
            .handle(readDataAndSymbolTableSizesResponse -> {
                readDataAndSymbolTableSizesTx.endRequest();
                if (readDataAndSymbolTableSizesResponse.getResult() != ReturnCode.OK) {
                    // TODO: Handle this
                    future.completeExceptionally(new PlcException("Reading data type and symbol table sizes failed: " + readDataAndSymbolTableSizesResponse.getResult()));
                    return;
                }
                try {
                    ReadBuffer readBuffer = new ReadBufferByteBased(readDataAndSymbolTableSizesResponse.getData());
                    AdsTableSizes adsTableSizes = AdsTableSizes.staticParse(readBuffer);
                    LOGGER.debug("PLC contains {} symbols and {} data-types", adsTableSizes.getSymbolCount(), adsTableSizes.getDataTypeCount());

                    // Now we load the datatype definitions.
                    AmsPacket readDataTypeTableRequest = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                        configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
                        ReservedIndexGroups.ADSIGRP_DATA_TYPE_TABLE_UPLOAD.getValue(), 0x00000000, adsTableSizes.getDataTypeLength());
                    RequestTransactionManager.RequestTransaction readDataTypeTableTx = tm.startRequest();
                    AmsTCPPacket amsReadTableTCPPacket = new AmsTCPPacket(readDataTypeTableRequest);
                    readDataTypeTableTx.submit(() -> context.sendRequest(amsReadTableTCPPacket)
                        .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                        .onTimeout(future::completeExceptionally)
                        .onError((p, e) -> future.completeExceptionally(e))
                        .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == readDataTypeTableRequest.getInvokeId())
                        .unwrap(response -> (AdsReadResponse) response.getUserdata())
                        .handle(readDataTypeTableResponse -> {
                            readDataTypeTableTx.endRequest();
                            if (readDataTypeTableResponse.getResult() != ReturnCode.OK) {
                                // TODO: Handle this
                                future.completeExceptionally(new PlcException("Reading data type table failed: " + readDataTypeTableResponse.getResult()));
                                return;
                            }
                            // Parse the result.
                            ReadBuffer rb = new ReadBufferByteBased(readDataTypeTableResponse.getData());
                            for (int i = 0; i < adsTableSizes.getDataTypeCount(); i++) {
                                try {
                                    AdsDataTypeTableEntry adsDataTypeTableEntry = AdsDataTypeTableEntry.staticParse(rb);
                                    dataTypeTable.put(adsDataTypeTableEntry.getDataTypeName(), adsDataTypeTableEntry);
                                } catch (ParseException e) {
                                    throw new RuntimeException(e);
                                }
                            }

                            // Now we load the symbol definitions
                            AmsPacket readSymbolTableRequest = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                                configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
                                ReservedIndexGroups.ADSIGRP_SYM_UPLOAD.getValue(), 0x00000000, adsTableSizes.getSymbolLength());
                            RequestTransactionManager.RequestTransaction readSymbolTableTx = tm.startRequest();
                            AmsTCPPacket amsReadSymbolTableTCPPacket = new AmsTCPPacket(readSymbolTableRequest);
                            readSymbolTableTx.submit(() -> context.sendRequest(amsReadSymbolTableTCPPacket)
                                .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                                .onTimeout(future::completeExceptionally)
                                .onError((p, e) -> future.completeExceptionally(e))
                                .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == readSymbolTableRequest.getInvokeId())
                                .unwrap(response -> (AdsReadResponse) response.getUserdata())
                                .handle(readSymbolTableResponse -> {
                                    readSymbolTableTx.endRequest();
                                    if (readSymbolTableResponse.getResult() != ReturnCode.OK) {
                                        // TODO: Handle this
                                        future.completeExceptionally(new PlcException("Reading symbol table failed: " + readSymbolTableResponse.getResult()));
                                        return;
                                    }
                                    ReadBuffer rb2 = new ReadBufferByteBased(readSymbolTableResponse.getData());
                                    for (int i = 0; i < adsTableSizes.getSymbolCount(); i++) {
                                        try {
                                            AdsSymbolTableEntry adsSymbolTableEntry = AdsSymbolTableEntry.staticParse(rb2);
                                            symbolTable.put(adsSymbolTableEntry.getName(), adsSymbolTableEntry);
                                        } catch (ParseException e) {
                                            throw new RuntimeException(e);
                                        }
                                    }

                                    LinkedHashMap<String, PlcSubscriptionTag> subscriptionTags = new LinkedHashMap<>();
                                    // Subscribe to online-version changes (get the address from the collected data for symbol: "TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt")
                                    subscriptionTags.put("onlineVersion", new DefaultPlcSubscriptionTag(
                                        PlcSubscriptionType.CHANGE_OF_STATE,
                                        new SymbolicAdsTag("TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt", org.apache.plc4x.java.api.types.PlcValueType.UDINT, Collections.emptyList()),
                                        Duration.ofMillis(1000)));
                                    // Subscribe to symbol-version changes (Address: GroupID: 0xF008, Offset: 0, Read length: 1)
                                    subscriptionTags.put("symbolVersion", new DefaultPlcSubscriptionTag(
                                        PlcSubscriptionType.CHANGE_OF_STATE,
                                        new DirectAdsTag(0xF008, 0x0000, "USINT", 1),
                                        Duration.ofMillis(1000)));
                                    LinkedHashMap<String, List<Consumer<PlcSubscriptionEvent>>> consumer = new LinkedHashMap<>();
                                    consumer.put("onlineVersion", Collections.singletonList(plcSubscriptionEvent -> {
                                        long oldVersion = onlineVersion;
                                        long newVersion = plcSubscriptionEvent.getPlcValue("onlineVersion").getLong();
                                        if(oldVersion != newVersion) {
                                            if(invalidationLock.tryLock()) {
                                                LOGGER.info("Detected change of the 'online-version', invalidating data type and symbol information.");
                                                CompletableFuture<Void> reloadingFuture = readSymbolTableAndDatatypeTable(context);
                                                reloadingFuture.whenComplete((unused, throwable) -> {
                                                    if(throwable != null) {
                                                        LOGGER.error("Error reloading data type and symbol data", throwable);
                                                    }
                                                    invalidationLock.unlock();
                                                });
                                            }
                                        }
                                    }));
                                    consumer.put("symbolVersion", Collections.singletonList(plcSubscriptionEvent -> {
                                        int oldVersion = symbolVersion;
                                        int newVersion = plcSubscriptionEvent.getPlcValue("symbolVersion").getInteger();
                                        if(oldVersion != newVersion) {
                                            if(invalidationLock.tryLock()) {
                                                LOGGER.info("Detected change of the 'symbol-version', invalidating data type and symbol information.");
                                                CompletableFuture<Void> reloadingFuture = readSymbolTableAndDatatypeTable(context);
                                                reloadingFuture.whenComplete((unused, throwable) -> {
                                                    if(throwable != null) {
                                                        LOGGER.error("Error reloading data type and symbol data", throwable);
                                                    }
                                                    invalidationLock.unlock();
                                                });
                                            }
                                        }
                                    }));
                                    PlcSubscriptionRequest subscriptionRequest = new DefaultPlcSubscriptionRequest(this, subscriptionTags, consumer);
                                    CompletableFuture<PlcSubscriptionResponse> subscriptionResponseCompletableFuture = subscribe(subscriptionRequest);

                                    // Wait for the subscription to be finished
                                    subscriptionResponseCompletableFuture.whenComplete((plcSubscriptionResponse, throwable) -> {
                                        if(throwable == null) {
                                            future.complete(null);
                                        }
                                    });
                                }));
                        }));
                    } catch (ParseException e) {
                        future.completeExceptionally(new PlcException("Error loading the table sizes", e));
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
        return browseWithInterceptor(browseRequest, item -> true);
    }

    public CompletableFuture<PlcBrowseResponse> browseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        CompletableFuture<PlcBrowseResponse> future = new CompletableFuture<>();
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        Map<String, List<PlcBrowseItem>> values = new HashMap<>();
        for (String queryName : browseRequest.getQueryNames()) {
            PlcQuery query = browseRequest.getQuery(queryName);
            List<PlcBrowseItem> resultsForQuery = new ArrayList<>();
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
                Map<String, PlcBrowseItem> childMap = new HashMap<>();
                for (PlcBrowseItem child : children) {
                    childMap.put(child.getName(), child);
                }

                // Populate a map of protocol-dependent options.
                Map<String, PlcValue> options = new HashMap<>();
                options.put("comment", new PlcSTRING(symbol.getComment()));
                options.put("group-id", new PlcUDINT(symbol.getGroup()));
                options.put("offset", new PlcUDINT(symbol.getOffset()));
                options.put("size-in-bytes", new PlcUDINT(symbol.getSize()));

                if(plc4xPlcValueType == org.apache.plc4x.java.api.types.PlcValueType.List) {
                    List<ArrayInfo> arrayInfo = new ArrayList<>();
                    for (AdsDataTypeArrayInfo adsDataTypeArrayInfo : dataType.getArrayInfo()) {
                        arrayInfo.add(new DefaultArrayInfo(
                            (int) adsDataTypeArrayInfo.getLowerBound(), (int) adsDataTypeArrayInfo.getUpperBound()));
                    }
                    DefaultListPlcBrowseItem item = new DefaultListPlcBrowseItem( new SymbolicAdsTag(symbol.getName(), plc4xPlcValueType, arrayInfo), itemName,
                        true, !symbol.getFlagReadOnly(), true, childMap, options);

                    // Check if this item should be added to the result
                    if(interceptor.intercept(item)) {
                        // Add the type itself.
                        resultsForQuery.add(item);
                    }
                } else {
                    DefaultPlcBrowseItem item = new DefaultPlcBrowseItem(new SymbolicAdsTag(symbol.getName(), plc4xPlcValueType, Collections.emptyList()), itemName, true,
                        !symbol.getFlagReadOnly(), true, childMap, options);

                    // Check if this item should be added to the result
                    if(interceptor.intercept(item)) {
                        // Add the type itself.
                        resultsForQuery.add(item);
                    }
                }
            }
            responseCodes.put(queryName, PlcResponseCode.OK);
            values.put(queryName, resultsForQuery);
            throw new RuntimeException("Gotta implement this ... Currently ignoring all queries.");
        }
        DefaultPlcBrowseResponse response = new DefaultPlcBrowseResponse(browseRequest, responseCodes, values);
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
            Map<String, PlcBrowseItem> childMap = new HashMap<>();
            for (PlcBrowseItem ch : children) {
                childMap.put(ch.getName(), ch);
            }

            // Populate a map of protocol-dependent options.
            Map<String, PlcValue> options = new HashMap<>();
            options.put("comment", new PlcSTRING(child.getComment()));
            options.put("group-id", new PlcUDINT(baseGroupId));
            options.put("offset", new PlcUDINT(baseOffset + child.getOffset()));
            options.put("size-in-bytes", new PlcUDINT(childDataType.getSize()));

            if(plc4xPlcValueType == org.apache.plc4x.java.api.types.PlcValueType.List) {
                List<ArrayInfo> arrayInfo = new ArrayList<>();
                for (AdsDataTypeArrayInfo adsDataTypeArrayInfo : childDataType.getArrayInfo()) {
                    arrayInfo.add(new DefaultArrayInfo(
                        (int) adsDataTypeArrayInfo.getLowerBound(), (int) adsDataTypeArrayInfo.getUpperBound()));
                }
                // Add the type itself.
                values.add(new DefaultListPlcBrowseItem(new SymbolicAdsTag(basePath + "." + child.getPropertyName(), plc4xPlcValueType, arrayInfo), itemName,
                    true, parentWritable, true, childMap, options));
            } else {
                // Add the type itself.
                values.add(new DefaultPlcBrowseItem(new SymbolicAdsTag(basePath + "." + child.getPropertyName(), plc4xPlcValueType, Collections.emptyList()), itemName,
                    true, parentWritable, true, childMap, options));
            }
        }
        return values;
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        // Get all ADS addresses in their resolved state.
        final CompletableFuture<Map<AdsTag, DirectAdsTag>> directAdsTagsFuture =
            getDirectAddresses(readRequest.getTags());

        // If all addresses were already resolved we can send the request immediately.
        if (directAdsTagsFuture.isDone()) {
            final Map<AdsTag, DirectAdsTag> resolvedTags = directAdsTagsFuture.getNow(null);
            if (resolvedTags != null) {
                return executeRead(readRequest, resolvedTags);
            } else {
                final CompletableFuture<PlcReadResponse> errorFuture = new CompletableFuture<>();
                errorFuture.completeExceptionally(new PlcException("Tags are null"));
                return errorFuture;
            }
        } else {
            // If there are still symbolic addresses that have to be resolved, send the
            // request as soon as the resolution is done.
            // In order to instantly be able to return a future, for the final result we have to
            // create a new one which is then completed later on. Unfortunately as soon as the
            // directAdsTagsFuture is completed we still don't have the end result, but we can
            // now actually send the delayed read request ... as soon as that future completes
            // we can complete the initial one.
            CompletableFuture<PlcReadResponse> delayedRead = new CompletableFuture<>();
            directAdsTagsFuture.handle((directAdsTags, throwable) -> {
                if (directAdsTags != null) {
                    final CompletableFuture<PlcReadResponse> delayedResponse =
                        executeRead(readRequest, directAdsTags);
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
                                                             Map<AdsTag, DirectAdsTag> resolvedTags) {
        // Depending on the number of tags, use a single item request or a sum-request
        if (resolvedTags.size() == 1) {
            // Do a normal (single item) ADS Read Request
            return singleRead(readRequest, resolvedTags.values().stream().findFirst().get());
        } else {
            // TODO: Check if the version of the remote station is at least TwinCAT v2.11 Build >= 1550 otherwise split up into single item requests.
            // Do a ADS-Sum Read Request.
            return multiRead(readRequest, resolvedTags);
        }
    }

    protected CompletableFuture<PlcReadResponse> singleRead(PlcReadRequest readRequest, DirectAdsTag directAdsTag) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        String dataTypeName = directAdsTag.getPlcDataType();
        AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(dataTypeName);
        long size = adsDataTypeTableEntry.getSize();

        AmsPacket amsPacket = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
            directAdsTag.getIndexGroup(), directAdsTag.getIndexOffset(), size * directAdsTag.getNumberOfElements());
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

    protected CompletableFuture<PlcReadResponse> multiRead(PlcReadRequest readRequest, Map<AdsTag, DirectAdsTag> resolvedTags) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        // Calculate the size of all tags together.
        // Calculate the expected size of the response data.
        long expectedResponseDataSize = resolvedTags.values().stream().mapToLong(
            tag -> {
                String dataTypeName = tag.getPlcDataType();
                AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(dataTypeName);
                long size = adsDataTypeTableEntry.getSize();
                // Status code + payload size
                return 4 + (size * tag.getNumberOfElements());
            }).sum();

        // With multi-requests, the index-group is fixed and the index offset indicates the number of elements.
        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_MULTIPLE_READ.getValue(), resolvedTags.size(),
            expectedResponseDataSize, readRequest.getTagNames().stream().map(tagName -> {
                AdsTag tag = (AdsTag) readRequest.getTag(tagName);
                DirectAdsTag directAdsTag = resolvedTags.get(tag);
                String dataTypeName = directAdsTag.getPlcDataType();
                AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(dataTypeName);
                long size = adsDataTypeTableEntry.getSize();
                return new AdsMultiRequestItemRead(
                    directAdsTag.getIndexGroup(), directAdsTag.getIndexOffset(),
                    (size * directAdsTag.getNumberOfElements()));
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
            responseCodes.put(readRequest.getTagNames().stream().findFirst().orElse(""),
                parsePlcResponseCode(adsReadResponse.getResult()));
        } else if (adsData instanceof AdsReadWriteResponse) {
            AdsReadWriteResponse adsReadWriteResponse = (AdsReadWriteResponse) adsData;
            readBuffer = new ReadBufferByteBased(adsReadWriteResponse.getData(), ByteOrder.LITTLE_ENDIAN);
            // When parsing a multi-item response, the error codes of each item comes
            // in sequence and then come the values.
            for (String tagName : readRequest.getTagNames()) {
                try {
                    final ReturnCode result = ReturnCode.enumForValue(readBuffer.readUnsignedLong(32));
                    responseCodes.put(tagName, parsePlcResponseCode(result));
                } catch (ParseException e) {
                    responseCodes.put(tagName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }
        if (readBuffer != null) {
            Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
            for (String tagName : readRequest.getTagNames()) {
                DirectAdsTag tag;
                if (readRequest.getTag(tagName) instanceof DirectAdsTag) {
                    tag = (DirectAdsTag) readRequest.getTag(tagName);
                } else {
                    SymbolicAdsTag symbolicAdsTag = (SymbolicAdsTag) readRequest.getTag(tagName);
                    tag = getDirectAdsTagForSymbolicName(symbolicAdsTag);
                }
                // If the response-code was anything but OK, we don't need to parse the payload.
                if (responseCodes.get(tagName) != PlcResponseCode.OK) {
                    values.put(tagName, new ResponseItem<>(responseCodes.get(tagName), null));
                }
                // If the response-code was ok, parse the data returned.
                else {
                    values.put(tagName, parseResponseItem(tag, readBuffer));
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

    private ResponseItem<PlcValue> parseResponseItem(DirectAdsTag tag, ReadBuffer readBuffer) {
        try {
            String dataTypeName = tag.getPlcDataType();
            AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(dataTypeName);
            PlcValueType plcValueType = getPlcValueTypeForAdsDataType(adsDataTypeTableEntry);

            int strLen = 0;
            if ((plcValueType == PlcValueType.STRING) || (plcValueType == PlcValueType.WSTRING)) {
                // Extract the string length from the data type name.
                strLen = Integer.parseInt(dataTypeName.substring(dataTypeName.indexOf("(") + 1, dataTypeName.indexOf(")")));
            }
            final int stringLength = strLen;
            if (tag.getNumberOfElements() == 1) {
                return new ResponseItem<>(PlcResponseCode.OK, parsePlcValue(plcValueType, adsDataTypeTableEntry, stringLength, readBuffer));
            } else {
                // Fetch all
                final PlcValue[] resultItems = IntStream.range(0, tag.getNumberOfElements()).mapToObj(i -> {
                    try {
                        return parsePlcValue(plcValueType, adsDataTypeTableEntry, stringLength, readBuffer);
                    } catch (ParseException e) {
                        LOGGER.warn("Error parsing tag item of type: '{}' (at position {}})", tag.getPlcDataType(), i, e);
                    }
                    return null;
                }).toArray(PlcValue[]::new);
                return new ResponseItem<>(PlcResponseCode.OK, PlcValueHandler.of(resultItems));
            }
        } catch (Exception e) {
            LOGGER.warn(String.format("Error parsing tag item of type: '%s'", tag.getPlcDataType()), e);
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
                    int strLen = 0;
                    if ((propertyPlcValueType == PlcValueType.STRING) || (propertyPlcValueType == PlcValueType.WSTRING)) {
                        String dataTypeName = propertyDataTypeTableEntry.getDataTypeName();
                        // Extract the string length from the data type name.
                        strLen = Integer.parseInt(dataTypeName.substring(dataTypeName.indexOf("(") + 1, dataTypeName.indexOf(")")));
                    }
                    PlcValue propertyValue = parsePlcValue(propertyPlcValueType, propertyDataTypeTableEntry, strLen, readBuffer);
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
        final CompletableFuture<Map<AdsTag, DirectAdsTag>> directAdsTagsFuture =
            getDirectAddresses(writeRequest.getTags());

        // If all addresses were already resolved we can send the request immediately.
        if (directAdsTagsFuture.isDone()) {
            final Map<AdsTag, DirectAdsTag> resolvedTags = directAdsTagsFuture.getNow(null);
            if (resolvedTags != null) {
                return executeWrite(writeRequest, resolvedTags);
            } else {
                final CompletableFuture<PlcWriteResponse> errorFuture = new CompletableFuture<>();
                errorFuture.completeExceptionally(new PlcException("Tags are null"));
                return errorFuture;
            }
        }
        // If there are still symbolic addresses that have to be resolved, send the
        // request as soon as the resolution is done.
        // In order to instantly be able to return a future, for the final result we have to
        // create a new one which is then completed later on. Unfortunately as soon as the
        // directAdsTagsFuture is completed we still don't have the end result, but we can
        // now actually send the delayed read request ... as soon as that future completes
        // we can complete the initial one.
        else {
            CompletableFuture<PlcWriteResponse> delayedWrite = new CompletableFuture<>();
            directAdsTagsFuture.handle((directAdsTags, throwable) -> {
                if (directAdsTags != null) {
                    final CompletableFuture<PlcWriteResponse> delayedResponse =
                        executeWrite(writeRequest, directAdsTags);
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
                                                               Map<AdsTag, DirectAdsTag> resolvedTags) {
        // Depending on the number of tags, use a single item request or a sum-request
        if (resolvedTags.size() == 1) {
            // Do a normal (single item) ADS Write Request
            return singleWrite(writeRequest, resolvedTags.values().stream().findFirst().get());
        } else {
            // TODO: Check if the version of the remote station is at least TwinCAT v2.11 Build >= 1550 otherwise split up into single item requests.
            // Do a ADS-Sum Read Request.
            return multiWrite(writeRequest, resolvedTags);
        }
    }

    protected CompletableFuture<PlcWriteResponse> singleWrite(PlcWriteRequest writeRequest, DirectAdsTag directAdsTag) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();

        final String tagName = writeRequest.getTagNames().iterator().next();
        final PlcValue plcValue = writeRequest.getPlcValue(tagName);

        try {
            byte[] serializedValue = serializePlcValue(plcValue, directAdsTag.getPlcDataType());
            AmsPacket amsPacket = new AdsWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
                0, getInvokeId(), directAdsTag.getIndexGroup(), directAdsTag.getIndexOffset(), serializedValue);
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
            future.completeExceptionally(new PlcException("Error", e));
        }
        return future;
    }

    protected CompletableFuture<PlcWriteResponse> multiWrite(PlcWriteRequest writeRequest, Map<AdsTag, DirectAdsTag> resolvedTags) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();

        int numTags = writeRequest.getTags().size();
        // Serialize all tags.
        List<byte[]> serializedTags = new ArrayList<>(numTags);
        Map<DirectAdsTag, AdsDataTypeTableEntry> directAdsTags = new LinkedHashMap<>(numTags);
        for (String tagName : writeRequest.getTagNames()) {
            final AdsTag tag = (AdsTag) writeRequest.getTag(tagName);
            final DirectAdsTag directAdsTag = resolvedTags.get(tag);
            final PlcValue plcValue = writeRequest.getPlcValue(tagName);
            final AdsDataTypeTableEntry dataType = dataTypeTable.get(directAdsTag.getPlcDataType());
            try {
                byte[] serializedValue = serializePlcValue(plcValue, directAdsTag.getPlcDataType());
                serializedTags.add(serializedValue);
                directAdsTags.put(directAdsTag, dataType);
            } catch (Exception e) {
                future.completeExceptionally(new PlcException("Error serializing data", e));
                return future;
            }
        }

        // Calculate the size of all serialized tags together.
        int serializedSize = serializedTags.stream().mapToInt(
            serializedTag -> serializedTag.length).sum();

        // Copy all serialized tags into one buffer.
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(serializedSize);
        for (byte[] serializedTag : serializedTags) {
            try {
                writeBuffer.writeByteArray("", serializedTag);
            } catch (SerializationException e) {
                future.completeExceptionally(new PlcException("Error serializing data", e));
                return future;
            }
        }

        // With multi-requests, the index-group is fixed and the index offset indicates the number of elements.
        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_MULTIPLE_WRITE.getValue(), serializedSize,
            (long) numTags * 4,
            directAdsTags.entrySet().stream().map(entry -> new AdsMultiRequestItemWrite(
                    entry.getKey().getIndexGroup(), entry.getKey().getIndexOffset(),
                    entry.getValue().getEntryLength()))
                .collect(Collectors.toList()), writeBuffer.getBytes());
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
            }));
        return future;
    }

    protected byte[] serializePlcValue(PlcValue plcValue, String datatypeName) throws SerializationException {
        // First check, if we have type information available.
        if (!dataTypeTable.containsKey(datatypeName)) {
            throw new SerializationException("Could not find data type: " + datatypeName);
        }

        // Get the data type, allocate enough memory and serialize the value based on the
        // structure defined by the data type.
        AdsDataTypeTableEntry dataType = dataTypeTable.get(datatypeName);
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased((int) dataType.getSize());
        List<AdsDataTypeArrayInfo> arrayInfo = dataType.getArrayInfo();
        serializeInternal(plcValue, dataType, arrayInfo, writeBuffer);
        return writeBuffer.getBytes();
    }

    protected void serializeInternal(PlcValue contextValue,
                                     AdsDataTypeTableEntry dataType,
                                     List<AdsDataTypeArrayInfo> arrayInfo,
                                     WriteBufferByteBased writeBuffer) throws SerializationException {

        // An array type: Recursively iterate over the elements
        if (arrayInfo.size() > 0) {
            if (!contextValue.isList()) {
                throw new SerializationException("Expected a PlcList, but got a " + contextValue.getPlcValueType().name());
            }
            AdsDataTypeArrayInfo curArrayLevel = arrayInfo.get(0);
            List<? extends PlcValue> list = contextValue.getList();
            if(curArrayLevel.getNumElements() != list.size()) {
                throw new SerializationException(String.format(
                    "Expected a PlcList of size %d, but got one of size %d", curArrayLevel.getNumElements(), list.size()));
            }
            for (PlcValue plcValue : list) {
                serializeInternal(plcValue, dataType, arrayInfo.subList(1, arrayInfo.size()), writeBuffer);
            }
        }

        // A complex type
        else if (dataType.getChildren().size() > 0) {
            if (!contextValue.isStruct()) {
                throw new SerializationException("Expected a PlcStruct, but got a " + contextValue.getPlcValueType().name());
            }
            PlcStruct plcStruct = (PlcStruct) contextValue;
            for (AdsDataTypeTableChildEntry child : dataType.getChildren()) {
                AdsDataTypeTableEntry childDataType = dataTypeTable.get(child.getDataTypeName());
                if (!plcStruct.hasKey(child.getPropertyName())) {
                    throw new SerializationException("PlcStruct is missing a child with the name " + child.getPropertyName());
                }
                PlcValue childValue = plcStruct.getValue(child.getPropertyName());
                serializeInternal(childValue, childDataType, childDataType.getArrayInfo(), writeBuffer);
            }
        }

        // A simple type
        else {
            PlcValueType plcValueType = getPlcValueTypeForAdsDataType(dataType);
            if (plcValueType == null) {
                throw new SerializationException("Unsupported simple type: " + dataType.getDataTypeName());
            }
            int stringLength = 0;
            if ((plcValueType == PlcValueType.STRING) || (plcValueType == PlcValueType.WSTRING)) {
                String stringTypeName = dataType.getDataTypeName();
                stringLength = Integer.parseInt(
                    stringTypeName.substring(stringTypeName.indexOf("(") + 1, stringTypeName.indexOf(")")));
            }
            DataItem.staticSerialize(writeBuffer, contextValue, plcValueType, stringLength);
        }
    }

    protected PlcWriteResponse convertToPlc4xWriteResponse(PlcWriteRequest writeRequest, AmsPacket adsData) {
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        if (adsData instanceof AdsWriteResponse) {
            AdsWriteResponse adsWriteResponse = (AdsWriteResponse) adsData;
            responseCodes.put(writeRequest.getTagNames().stream().findFirst().orElse(""),
                parsePlcResponseCode(adsWriteResponse.getResult()));
        } else if (adsData instanceof AdsReadWriteResponse) {
            AdsReadWriteResponse adsReadWriteResponse = (AdsReadWriteResponse) adsData;
            ReadBuffer readBuffer = new ReadBufferByteBased(adsReadWriteResponse.getData(), ByteOrder.LITTLE_ENDIAN);
            // When parsing a multi-item response, the error codes of each items come
            // in sequence and then come the values.
            for (String tagName : writeRequest.getTagNames()) {
                try {
                    final ReturnCode result = ReturnCode.enumForValue(readBuffer.readUnsignedLong(32));
                    responseCodes.put(tagName, parsePlcResponseCode(result));
                } catch (ParseException e) {
                    responseCodes.put(tagName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }

        return new DefaultPlcWriteResponse(writeRequest, responseCodes);
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        // Get all ADS addresses in their resolved state.
        final CompletableFuture<Map<AdsTag, DirectAdsTag>> directAdsTagFutures =
            getDirectAddresses(subscriptionRequest.getTags()
                .stream()
                .map(tag -> ((DefaultPlcSubscriptionTag) tag).getTag())
                .collect(Collectors.toList()));

        // If all addresses were already resolved we can send the request immediately.
        if (directAdsTagFutures.isDone()) {
            final Map<AdsTag, DirectAdsTag> resolvedTags = directAdsTagFutures.getNow(null);
            if (resolvedTags != null) {
                return executeSubscribe(subscriptionRequest, resolvedTags);
            } else {
                final CompletableFuture<PlcSubscriptionResponse> errorFuture = new CompletableFuture<>();
                errorFuture.completeExceptionally(new PlcException("Tags are null"));
                return errorFuture;
            }
        }
        // If there are still symbolic addresses that have to be resolved, send the
        // request as soon as the resolution is done.
        // In order to instantly be able to return a future, for the final result we have to
        // create a new one which is then completed later on. Unfortunately as soon as the
        // directAdsTagsFuture is completed we still don't have the end result, but we can
        // now actually send the delayed read request ... as soon as that future completes
        // we can complete the initial one.
        else {
            CompletableFuture<PlcSubscriptionResponse> delayedSubscribe = new CompletableFuture<>();
            directAdsTagFutures.handle((tagMapping, throwable) -> {
                if (tagMapping != null) {
                    final CompletableFuture<PlcSubscriptionResponse> delayedResponse =
                        executeSubscribe(subscriptionRequest, tagMapping);
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

    private CompletableFuture<PlcSubscriptionResponse> executeSubscribe(PlcSubscriptionRequest subscribeRequest, Map<AdsTag, DirectAdsTag> resolvedTags) {
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();

        List<AmsTCPPacket> amsTCPPackets = subscribeRequest.getTags().stream()
            .map(tag -> (DefaultPlcSubscriptionTag) tag)
            .map(tag -> {
                AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(resolvedTags.get((AdsTag) tag.getTag()).getPlcDataType());
                DirectAdsTag directAdsTag = getDirectAdsTagForSymbolicName(tag.getTag());
                // TODO: We should implement multi-dimensional arrays here ...
                int numberOfElements = (tag.getArrayInfo().size() == 0) ? 1 : tag.getArrayInfo().get(0).getSize();
                return new AmsTCPPacket(new AdsAddDeviceNotificationRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                    configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
                    0, getInvokeId(),
                    directAdsTag.getIndexGroup(),
                    directAdsTag.getIndexOffset(),
                    adsDataTypeTableEntry.getSize() * numberOfElements,
                    tag.getPlcSubscriptionType() == PlcSubscriptionType.CYCLIC ? AdsTransMode.CYCLIC : AdsTransMode.ON_CHANGE, // if it's not cyclic, it's on change or event
                    0, // there is no api for that yet
                    tag.getDuration().orElse(Duration.ZERO).toMillis()));
            })
            .collect(Collectors.toList());

        Map<String, ResponseItem<PlcSubscriptionHandle>> responses = new HashMap<>();

        // Start the first request-transaction (it is ended in the response-handler).
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(subscribeRecursively(
            subscribeRequest,
            subscribeRequest.getTagNames().iterator(),
            resolvedTags,
            responses,
            future,
            amsTCPPackets.iterator(),
            transaction));
        return future;
    }

    private Runnable subscribeRecursively(PlcSubscriptionRequest subscriptionRequest,
                                          Iterator<String> tagNames,
                                          Map<AdsTag, DirectAdsTag> resolvedTags,
                                          Map<String, ResponseItem<PlcSubscriptionHandle>> responses,
                                          CompletableFuture<PlcSubscriptionResponse> future,
                                          Iterator<AmsTCPPacket> amsTCPPackets,
                                          RequestTransactionManager.RequestTransaction transaction) {
        return () -> {
            AmsTCPPacket packet = amsTCPPackets.next();
            boolean hasMorePackets = amsTCPPackets.hasNext();
            String tagName = tagNames.next();
            context.sendRequest(packet)
                .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == packet.getUserdata().getInvokeId())
                .unwrap(responseAmsPacket -> (AdsAddDeviceNotificationResponse) responseAmsPacket.getUserdata())
                .handle(response -> {
                    if (response.getResult() == ReturnCode.OK) {
                        DefaultPlcSubscriptionTag subscriptionTag = (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(tagName);
                        AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get((resolvedTags.get((AdsTag) subscriptionTag.getTag())).getPlcDataType());

                        // Collect notification handle from individual response.
                        responses.put(tagName, new ResponseItem<>(
                            parsePlcResponseCode(response.getResult()),
                            new AdsSubscriptionHandle(this,
                                tagName,
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
                            subscriptionRequest, tagNames, resolvedTags, responses, future, amsTCPPackets, nextTransaction));
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
        values.put(subscriptionHandle.getTagName(), new ResponseItem<>(PlcResponseCode.OK,
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

    protected CompletableFuture<Map<AdsTag, DirectAdsTag>> getDirectAddresses(List<PlcTag> tags) {
        CompletableFuture<Map<AdsTag, DirectAdsTag>> future = new CompletableFuture<>();

        // Get all symbolic tags from the current request.
        // These potentially need to be resolved to direct addresses, if this has not been done before.
        final List<SymbolicAdsTag> referencedSymbolicTags = tags.stream()
            .filter(SymbolicAdsTag.class::isInstance)
            .map(SymbolicAdsTag.class::cast)
            .collect(Collectors.toList());

        // Find out for which of these symbolic addresses no resolution has been initiated.
        final List<SymbolicAdsTag> symbolicTagsNeedingResolution = referencedSymbolicTags.stream()
            .filter(symbolicAdsTag -> getDirectAdsTagForSymbolicName(symbolicAdsTag) == null)
            .collect(Collectors.toList());

        // If there are unresolved symbolic addresses, initiate the resolution
        if (!symbolicTagsNeedingResolution.isEmpty()) {
            // Get a list of symbolic addresses for which no resolution request has been sent yet
            // (A parallel request initiated a bit earlier might have already initiated a resolution
            // which has not yet been completed)
            final List<SymbolicAdsTag> requiredResolutionTags =
                symbolicTagsNeedingResolution.stream().filter(symbolicAdsTags ->
                    !pendingResolutionRequests.containsKey(symbolicAdsTags)).collect(Collectors.toList());
            // If there are tags for which no resolution request has been sent yet,
            // send a request.
            if (!requiredResolutionTags.isEmpty()) {
                CompletableFuture<Void> resolutionFuture;
                // Create a future which will be completed as soon as the
                // resolution result has been added to the map.
                if (requiredResolutionTags.size() == 1) {
                    SymbolicAdsTag symbolicAdsTag = requiredResolutionTags.get(0);
                    resolutionFuture = resolveSingleSymbolicAddress(requiredResolutionTags.get(0));
                    pendingResolutionRequests.put(symbolicAdsTag, resolutionFuture);
                } else {
                    resolutionFuture = resolveMultipleSymbolicAddresses(requiredResolutionTags);
                    for (SymbolicAdsTag symbolicAdsTag : requiredResolutionTags) {
                        pendingResolutionRequests.put(symbolicAdsTag, resolutionFuture);
                    }
                }
            }

            // Create a global future which is completed as soon as all sub-futures for this request are completed.
            final CompletableFuture<Void> resolutionComplete =
                CompletableFuture.allOf(symbolicTagsNeedingResolution.stream()
                    .map(pendingResolutionRequests::get)
                    .toArray(CompletableFuture[]::new));

            // Complete the future asynchronously as soon as all tags are resolved.
            resolutionComplete.handleAsync((unused, throwable) -> {
                if (throwable != null) {
                    return future.completeExceptionally(throwable.getCause());
                } else {
                    Map<AdsTag, DirectAdsTag> directAdsTagMapping = new HashMap<>(tags.size());
                    for (PlcTag tag : tags) {
                        if (tag instanceof SymbolicAdsTag) {
                            directAdsTagMapping.put((AdsTag) tag, getDirectAdsTagForSymbolicName(tag));
                        } else {
                            directAdsTagMapping.put((AdsTag) tag, (DirectAdsTag) tag);
                        }
                    }
                    return future.complete(directAdsTagMapping);
                }
            });
        } else {
            // If all tags were resolved, we can continue instantly.
            Map<AdsTag, DirectAdsTag> directAdsTagMapping = new HashMap<>(tags.size());
            for (PlcTag tag : tags) {
                if (tag instanceof SymbolicAdsTag) {
                    directAdsTagMapping.put((AdsTag) tag, getDirectAdsTagForSymbolicName(tag));
                } else {
                    directAdsTagMapping.put((AdsTag) tag, (DirectAdsTag) tag);
                }
            }
            future.complete(directAdsTagMapping);
        }
        return future;
    }

    protected CompletableFuture<Void> resolveSingleSymbolicAddress(SymbolicAdsTag symbolicAdsTag) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_SYM_HNDBYNAME.getValue(), 0,
            4, null,
            getNullByteTerminatedArray(symbolicAdsTag.getSymbolicAddress()));
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
                    future.completeExceptionally(new PlcException("Couldn't retrieve handle for symbolic tag " +
                        symbolicAdsTag.getSymbolicAddress() + " got return code " + response.getResult().name()));
                } else {
                    ReadBuffer readBuffer = new ReadBufferByteBased(response.getData(), ByteOrder.LITTLE_ENDIAN);
                    try {
                        // Read the handle.
                        long handle = readBuffer.readUnsignedLong(32);

/*                        DirectAdsTag directAdsTag = new DirectAdsTag(
                            ReservedIndexGroups.ADSIGRP_SYM_VALBYHND.getValue(), handle,
                            symbolicAdsTag.getAdsDataTypeName(), symbolicAdsTag.getNumberOfElements());*/

                        // TODO: Find out how to read the datatype for the given symbolic tag
                        //symbolicTagMapping.put(symbolicAdsTag, directAdsTag);
                        future.complete(null);
                    } catch (ParseException e) {
                        future.completeExceptionally(e);
                    }
                }
                transaction.endRequest();
            }));
        return future;
    }

    protected CompletableFuture<Void> resolveMultipleSymbolicAddresses(List<SymbolicAdsTag> symbolicAdsTags) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        // The expected response for every symbolic address is 12 bytes (8 bytes header and 4 bytes for the handle)
        long expectedResponseDataSize = (long) (symbolicAdsTags.size()) * 12;
        // Concatenate the string part of each symbolic address into one concatenated string and get the bytes.
        byte[] addressData = symbolicAdsTags.stream().map(
            SymbolicAdsTag::getSymbolicAddress).collect(Collectors.joining("")).getBytes();
        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_MULTIPLE_READ_WRITE.getValue(),
            symbolicAdsTags.size(), expectedResponseDataSize, symbolicAdsTags.stream().map(symbolicAdsTag ->
            new AdsMultiRequestItemReadWrite(ReservedIndexGroups.ADSIGRP_SYM_HNDBYNAME.getValue(), 0,
                4, symbolicAdsTag.getSymbolicAddress().length())).collect(Collectors.toList()), null);
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
                Map<SymbolicAdsTag, Long> returnCodes = new HashMap<>();
                // In the response first come the return codes and the data-lengths for each item.
                symbolicAdsTags.forEach(symbolicAdsTag -> {
                    try {
                        // This should be 0 in the success case.
                        long returnCode = readBuffer.readUnsignedLong(32);
                        // This is always 4
                        long itemLength = readBuffer.readUnsignedLong(32);
                        assert itemLength == 4;

                        returnCodes.put(symbolicAdsTag, returnCode);
                    } catch (ParseException e) {
                        throw new PlcRuntimeException(e);
                    }
                });
                // After reading the header-information, comes the data itself.
                symbolicAdsTags.forEach(symbolicAdsTag -> {
                    try {
                        if (returnCodes.get(symbolicAdsTag) == 0) {
                            // Read the handle.
                            long handle = readBuffer.readUnsignedLong(32);

                            /*DirectAdsTag directAdsTag = new DirectAdsTag(
                                ReservedIndexGroups.ADSIGRP_SYM_VALBYHND.getValue(), handle,
                                symbolicAdsTag.getAdsDataTypeName(), symbolicAdsTag.getNumberOfElements());*/
                            // TODO: Find out how to read the datatype for the given symbolic Tag
                            //symbolicTagMapping.put(symbolicAdsTag, directAdsTag);
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

    protected DirectAdsTag getDirectAdsTagForSymbolicName(PlcTag tag) {
        if (tag instanceof DirectAdsTag) {
            return (DirectAdsTag) tag;
        }

        SymbolicAdsTag symbolicAdsTag = (SymbolicAdsTag) tag;
        String symbolicAddress = symbolicAdsTag.getSymbolicAddress();
        String[] addressParts = symbolicAddress.split("\\.");

        // If the number of parts are less than 2, we can find the entry in the symbol table directly.
        if (addressParts.length < 2) {
            // We can't find it, so we need to resolve it.
            if (!symbolTable.containsKey(symbolicAddress)) {
                return null;
            }
            AdsSymbolTableEntry adsSymbolTableEntry = symbolTable.get(symbolicAddress);
            if(adsSymbolTableEntry == null) {
                throw new PlcInvalidTagException("Couldn't resolve symbolic address: " + symbolicAddress);
            }
            AdsDataTypeTableEntry dataTypeTableEntry = dataTypeTable.get(adsSymbolTableEntry.getDataTypeName());
            if(dataTypeTableEntry == null) {
                throw new PlcInvalidTagException(
                    "Couldn't resolve datatype: '" + adsSymbolTableEntry.getDataTypeName() +
                        "' for address: '" + ((SymbolicAdsTag) tag).getSymbolicAddress() + "'");
            }
            return new DirectAdsTag(adsSymbolTableEntry.getGroup(), adsSymbolTableEntry.getOffset(),
                dataTypeTableEntry.getDataTypeName(), dataTypeTableEntry.getArrayDimensions());
        }
        // Otherwise we'll have to crawl through the dataType definitions.
        else {
            String symbolName = addressParts[0] + "." + addressParts[1];
            AdsSymbolTableEntry adsSymbolTableEntry = symbolTable.get(symbolName);
            if(adsSymbolTableEntry == null) {
                throw new PlcInvalidTagException("Couldn't resolve symbolic address: " + symbolName);
            }
            AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTable.get(adsSymbolTableEntry.getDataTypeName());
            if(adsDataTypeTableEntry == null) {
                throw new PlcInvalidTagException(
                    "Couldn't resolve datatype: '" + adsSymbolTableEntry.getDataTypeName() +
                        "' for address: '" + ((SymbolicAdsTag) tag).getSymbolicAddress() + "'");
            }
            return resolveDirectAdsTagForSymbolicNameFromDataType(
                Arrays.asList(addressParts).subList(2, addressParts.length),
                adsSymbolTableEntry.getGroup(), adsSymbolTableEntry.getOffset(), adsDataTypeTableEntry);
        }
    }

    protected DirectAdsTag resolveDirectAdsTagForSymbolicNameFromDataType(List<String> remainingAddressParts, long currentGroup, long currentOffset, AdsDataTypeTableEntry adsDataTypeTableEntry) {
        if (remainingAddressParts.isEmpty()) {
            // TODO: Implement the Array support
            return new DirectAdsTag(currentGroup, currentOffset, adsDataTypeTableEntry.getDataTypeName(), 1);
        }

        // Go through all children looking for a matching one.
        for (AdsDataTypeTableChildEntry child : adsDataTypeTableEntry.getChildren()) {
            if (child.getPropertyName().equals(remainingAddressParts.get(0))) {
                AdsDataTypeTableEntry childAdsDataTypeTableEntry = dataTypeTable.get(child.getDataTypeName());
                return resolveDirectAdsTagForSymbolicNameFromDataType(
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
