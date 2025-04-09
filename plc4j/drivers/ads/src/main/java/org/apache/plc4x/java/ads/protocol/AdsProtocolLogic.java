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
import org.apache.plc4x.java.ads.discovery.readwrite.AmsNetId;
import org.apache.plc4x.java.ads.discovery.readwrite.*;
import org.apache.plc4x.java.ads.model.AdsSubscriptionHandle;
import org.apache.plc4x.java.ads.readwrite.*;
import org.apache.plc4x.java.ads.tag.AdsTag;
import org.apache.plc4x.java.ads.tag.AdsTagHandler;
import org.apache.plc4x.java.ads.tag.DirectAdsStringTag;
import org.apache.plc4x.java.ads.tag.DirectAdsTag;
import org.apache.plc4x.java.ads.tag.SymbolicAdsTag;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.spi.metadata.DefaultMetadata;
import org.apache.plc4x.java.api.metadata.time.TimeSource;
import org.apache.plc4x.java.api.model.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcTagItem;
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
    public void close(ConversationContext<AmsTCPPacket> context) {
        tm.shutdown();
    }

    @Override
    public void setConfiguration(AdsConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new AdsTagHandler();
    }

    @Override
    public void onConnect(ConversationContext<AmsTCPPacket> context) {

        // If we have connection credentials available, try to set up the AMS routes.
        CompletableFuture<Void> setupAmsRouteFuture;
        if (context.getAuthentication() != null) {
            if (!(context.getAuthentication() instanceof PlcUsernamePasswordAuthentication)) {
                context.getChannel().pipeline().fireExceptionCaught(new PlcConnectionException(
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
                context.fireConnected();
                // Instead of aborting here, we allow connecting, as the user might be using raw-addresses
                // This is particularly important for using the driver together with ADS over EtherCAT for accessing
                // direct EtherCAT devices.
//                future.completeExceptionally(new PlcConnectionException(
//                    "Lazy loading is generally planned, but not implemented yet. " +
//                        "If you are in need for this feature, please reach out to the community."));
            } else {
                // Execute a ReadDeviceInfo command
                AmsPacket readDeviceInfoRequest = new AdsReadDeviceInfoRequest(
                    configuration.getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
                    configuration.getSourceAmsNetId(), 800, 0, getInvokeId());
                RequestTransactionManager.RequestTransaction readDeviceInfoTx = tm.startRequest();
                readDeviceInfoTx.submit(() -> context.sendRequest(new AmsTCPPacket(readDeviceInfoRequest))
                    .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                    .onTimeout(e -> context.getChannel().pipeline().fireExceptionCaught(e))
                    .onError((p, e) -> context.getChannel().pipeline().fireExceptionCaught(e))
                    .unwrap(AmsTCPPacket::getUserdata)
                    .check(userdata -> userdata.getInvokeId() == readDeviceInfoRequest.getInvokeId())
                    .only(AdsReadDeviceInfoResponse.class)
                    .handle(readDeviceInfoResponse -> {
                        readDeviceInfoTx.endRequest();
                        if (readDeviceInfoResponse.getResult() != ReturnCode.OK) {
                            context.getChannel().pipeline().fireExceptionCaught(new PlcConnectionException(
                                "Error reading device info. Got: " + readDeviceInfoResponse.getResult()));
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
                            .onTimeout(e -> context.getChannel().pipeline().fireExceptionCaught(e))
                            .onError((p, e) -> context.getChannel().pipeline().fireExceptionCaught(e))
                            .unwrap(AmsTCPPacket::getUserdata)
                            .check(userdata -> userdata.getInvokeId() == readOnlineVersionNumberRequest.getInvokeId())
                            .only(AdsReadWriteResponse.class)
                            .handle(readOnlineVersionNumberResponse -> {
                                readOnlineVersionNumberTx.endRequest();
                                if (readOnlineVersionNumberResponse.getResult() != ReturnCode.OK) {
                                    context.getChannel().pipeline().fireExceptionCaught(new PlcConnectionException(
                                        "Error reading online version number. Got: " + readOnlineVersionNumberResponse.getResult()));
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
                                        .onTimeout(e -> context.getChannel().pipeline().fireExceptionCaught(e))
                                        .onError((p, e) -> context.getChannel().pipeline().fireExceptionCaught(e))
                                        .unwrap(AmsTCPPacket::getUserdata)
                                        .check(userdata -> userdata.getInvokeId() == readSymbolVersionNumberRequest.getInvokeId())
                                        .only(AdsReadResponse.class)
                                        .handle(readSymbolVersionNumberResponse -> {
                                            readSymbolVersionNumberTx.endRequest();
                                            if (readSymbolVersionNumberResponse.getResult() != ReturnCode.OK) {
                                                context.getChannel().pipeline().fireExceptionCaught(new PlcConnectionException(
                                                    "Error reading symbol version number. Got: " + readDeviceInfoResponse.getResult()));
                                                return;
                                            }
                                            try {
                                                ReadBuffer rb2 = new ReadBufferByteBased(readSymbolVersionNumberResponse.getData());
                                                symbolVersion = rb2.readUnsignedInt(8);

                                                LOGGER.debug("Fetching sizes of symbol and datatype table sizes.");
                                                CompletableFuture<Void> readSymbolTableFuture = readSymbolTableAndDatatypeTable(context);
                                                readSymbolTableFuture.whenComplete((unused2, throwable2) -> {
                                                    if (throwable2 != null) {
                                                        context.getChannel().pipeline().fireExceptionCaught(new PlcConnectionException(
                                                            "Error reading symbol and datatype table sizes", throwable2));
                                                    } else {
                                                        context.fireConnected();
                                                    }
                                                });
                                            } catch (ParseException e) {
                                                context.getChannel().pipeline().fireExceptionCaught(new PlcConnectionException(
                                                    "Error parsing reading symbol and datatype table sizes response.", e));
                                            }
                                        }));
                                } catch (ParseException e) {
                                    context.getChannel().pipeline().fireExceptionCaught(new PlcConnectionException(
                                        "Error parsing online version number response.", e));
                                }
                            }));
                    }));
            }
        });
    }

    protected CompletableFuture<Void> setupAmsRoute(PlcUsernamePasswordAuthentication authentication) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        new Thread(() -> {
            LOGGER.debug("Setting up remote AMS routes.");
            SocketAddress localSocketAddress = conversationContext.getChannel().localAddress();
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
                SocketAddress remoteSocketAddress = conversationContext.getChannel().remoteAddress();
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
            .unwrap(AmsTCPPacket::getUserdata)
            .check(userdata -> userdata.getInvokeId() == readDataAndSymbolTableSizesRequest.getInvokeId())
            .only(AdsReadResponse.class)
            .handle(readDataAndSymbolTableSizesResponse -> {
                readDataAndSymbolTableSizesTx.endRequest();
                if (readDataAndSymbolTableSizesResponse.getResult() != ReturnCode.OK) {
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
                        .unwrap(AmsTCPPacket::getUserdata)
                        .check(userdata -> userdata.getInvokeId() == readDataTypeTableRequest.getInvokeId())
                        .only(AdsReadResponse.class)
                        .handle(readDataTypeTableResponse -> {
                            readDataTypeTableTx.endRequest();
                            if (readDataTypeTableResponse.getResult() != ReturnCode.OK) {
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
                                .unwrap(AmsTCPPacket::getUserdata)
                                .check(userdata -> userdata.getInvokeId() == readSymbolTableRequest.getInvokeId())
                                .only(AdsReadResponse.class)
                                .handle(readSymbolTableResponse -> {
                                    readSymbolTableTx.endRequest();
                                    if (readSymbolTableResponse.getResult() != ReturnCode.OK) {
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

                                    LinkedHashMap<String, PlcTagItem<PlcSubscriptionTag>> subscriptionTags = new LinkedHashMap<>();
                                    // Subscribe to online-version changes (get the address from the collected data for symbol: "TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt")
                                    subscriptionTags.put("onlineVersion", new DefaultPlcTagItem<>(new DefaultPlcSubscriptionTag(
                                        PlcSubscriptionType.CHANGE_OF_STATE,
                                        new SymbolicAdsTag("TwinCAT_SystemInfoVarList._AppInfo.OnlineChangeCnt", org.apache.plc4x.java.api.types.PlcValueType.UDINT, Collections.emptyList()),
                                        Duration.ofMillis(1000))));
                                    // Subscribe to symbol-version changes (Address: GroupID: 0xF008, Offset: 0, Read length: 1)
                                    subscriptionTags.put("symbolVersion", new DefaultPlcTagItem<>(new DefaultPlcSubscriptionTag(
                                        PlcSubscriptionType.CHANGE_OF_STATE,
                                        new DirectAdsTag(0xF008, 0x0000, "USINT", 1),
                                        Duration.ofMillis(1000))));
                                    Consumer<PlcSubscriptionEvent> consumer = plcSubscriptionEvent -> {
                                        for (String tagName : plcSubscriptionEvent.getTagNames()) {
                                            switch (tagName) {
                                                case "onlineVersion": {
                                                    long oldVersion = onlineVersion;
                                                    long newVersion = plcSubscriptionEvent.getPlcValue("onlineVersion").getLong();
                                                    if (oldVersion != newVersion) {
                                                        if (invalidationLock.tryLock()) {
                                                            LOGGER.info("Detected change of the 'online-version', invalidating data type and symbol information.");
                                                            CompletableFuture<Void> reloadingFuture = readSymbolTableAndDatatypeTable(context);
                                                            reloadingFuture.whenComplete((unused, throwable) -> {
                                                                if (throwable != null) {
                                                                    LOGGER.error("Error reloading data type and symbol data", throwable);
                                                                }
                                                                invalidationLock.unlock();
                                                            });
                                                        }
                                                    }
                                                    break;
                                                }
                                                case "symbolVersion": {
                                                    int oldVersion = symbolVersion;
                                                    int newVersion = plcSubscriptionEvent.getPlcValue("symbolVersion").getInteger();
                                                    if (oldVersion != newVersion) {
                                                        if (invalidationLock.tryLock()) {
                                                            LOGGER.info("Detected change of the 'symbol-version', invalidating data type and symbol information.");
                                                            CompletableFuture<Void> reloadingFuture = readSymbolTableAndDatatypeTable(context);
                                                            reloadingFuture.whenComplete((unused, throwable) -> {
                                                                if (throwable != null) {
                                                                    LOGGER.error("Error reloading data type and symbol data", throwable);
                                                                }
                                                                invalidationLock.unlock();
                                                            });
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    };

                                    PlcSubscriptionRequest subscriptionRequest = new DefaultPlcSubscriptionRequest(this, subscriptionTags, consumer, Collections.emptyMap());
                                    CompletableFuture<PlcSubscriptionResponse> subscriptionResponseCompletableFuture = subscribe(subscriptionRequest);

                                    // Wait for the subscription to be finished
                                    subscriptionResponseCompletableFuture.whenComplete((plcSubscriptionResponse, throwable) -> {
                                        if (throwable == null) {
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
            List<PlcBrowseItem> resultsForQuery = new ArrayList<>();
            for (AdsSymbolTableEntry symbol : symbolTable.values()) {
                // Get the datatype of this entry.
                Optional<AdsDataTypeTableEntry> dataTypeTableEntry = getDataTypeTableEntry(symbol.getDataTypeName());
                if (dataTypeTableEntry.isEmpty()) {
                    System.out.printf("couldn't find datatype: %s%n", symbol.getDataTypeName());
                    continue;
                }
                AdsDataTypeTableEntry dataType = dataTypeTableEntry.get();

                // Convert the plc value type from the ADS specific one to the PLC4X global one.
                PlcValueType plcValueType = getPlcValueTypeForAdsDataTypeForBrowse(dataType);

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

                List<ArrayInfo> arrayInfo = new ArrayList<>(dataType.getArrayInfo().size());
                List<ArrayInfo> itemArrayInfo = new ArrayList<>(dataType.getArrayInfo().size());
                for (AdsDataTypeArrayInfo adsDataTypeArrayInfo : dataType.getArrayInfo()) {
                    arrayInfo.add(new DefaultArrayInfo(
                        (int) adsDataTypeArrayInfo.getLowerBound(), (int) adsDataTypeArrayInfo.getUpperBound()));
                    itemArrayInfo.add(new DefaultArrayInfo(
                        (int) adsDataTypeArrayInfo.getLowerBound(), (int) adsDataTypeArrayInfo.getUpperBound()));
                }
                DefaultPlcBrowseItem item = new DefaultPlcBrowseItem(new SymbolicAdsTag(symbol.getName(), plcValueType, arrayInfo), symbol.getName(),
                    true, !symbol.getFlagReadOnly(), true, false, itemArrayInfo, childMap, options);

                // Check if this item should be added to the result
                if (interceptor.intercept(item)) {
                    // Add the type itself.
                    resultsForQuery.add(item);
                }
            }
            responseCodes.put(queryName, PlcResponseCode.OK);
            values.put(queryName, resultsForQuery);
        }
        DefaultPlcBrowseResponse response = new DefaultPlcBrowseResponse(browseRequest, responseCodes, values);
        future.complete(response);
        return future;
    }

    protected List<PlcBrowseItem> getBrowseItems(String basePath, long baseGroupId, long baseOffset, boolean parentWritable, AdsDataTypeTableEntry dataType) {
        // If this is an array type, then we need to lookup it's elemental type and use that instead
        if(dataType.getArrayDimensions() > 0) {
            Optional<AdsDataTypeTableEntry> dataTypeTableEntry = getDataTypeTableEntry(dataType.getSimpleTypeName());
            if(dataTypeTableEntry.isEmpty()) {
                LOGGER.warn("couldn't find datatype: {}", dataType.getSimpleTypeName());
                return Collections.emptyList();
            }
            dataType = dataTypeTableEntry.get();
        }

        if (dataType.getNumChildren() == 0) {
            return Collections.emptyList();
        }

        List<PlcBrowseItem> values = new ArrayList<>(dataType.getNumChildren());
        for (AdsDataTypeTableChildEntry child : dataType.getChildren()) {
            Optional<AdsDataTypeTableEntry> dataTypeTableEntry = getDataTypeTableEntry(child.getDataTypeName());
            if(dataTypeTableEntry.isEmpty()) {
                LOGGER.warn("couldn't find datatype: {} for child {}", dataType.getSimpleTypeName(), child.getPropertyName());
                continue;
            }
            AdsDataTypeTableEntry childDataType = dataTypeTableEntry.get();
            String itemAddress = basePath + "." + child.getPropertyName();

            // Convert the plc value type from the ADS specific one to the PLC4X global one.
            org.apache.plc4x.java.api.types.PlcValueType plc4xPlcValueType =
                org.apache.plc4x.java.api.types.PlcValueType.valueOf(getPlcValueTypeForAdsDataTypeForBrowse(childDataType).toString());

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

            List<ArrayInfo> arrayInfo = new ArrayList<>(childDataType.getArrayInfo().size());
            List<ArrayInfo> itemArrayInfo = new ArrayList<>(childDataType.getArrayInfo().size());
            for (AdsDataTypeArrayInfo adsDataTypeArrayInfo : childDataType.getArrayInfo()) {
                arrayInfo.add(new DefaultArrayInfo(
                    (int) adsDataTypeArrayInfo.getLowerBound(), (int) adsDataTypeArrayInfo.getUpperBound()));
                itemArrayInfo.add(new DefaultArrayInfo(
                    (int) adsDataTypeArrayInfo.getLowerBound(), (int) adsDataTypeArrayInfo.getUpperBound()));
            }
            // Add the type itself.
            values.add(new DefaultPlcBrowseItem(new SymbolicAdsTag(
                basePath + "." + child.getPropertyName(), plc4xPlcValueType, arrayInfo), child.getPropertyName(),
                true, parentWritable, true, false, itemArrayInfo, childMap, options));
        }
        return values;
    }

    @Override
    public CompletableFuture<PlcPingResponse> ping(PlcPingRequest pingRequest) {
        CompletableFuture<PlcPingResponse> future = new CompletableFuture<>();

        AmsPacket readDeviceInfoRequest = new AdsReadDeviceInfoRequest(
            configuration.getTargetAmsNetId(), DefaultAmsPorts.RUNTIME_SYSTEM_01.getValue(),
            configuration.getSourceAmsNetId(), 800, 0, getInvokeId());

        RequestTransactionManager.RequestTransaction readDeviceInfoTx = tm.startRequest();
        readDeviceInfoTx.submit(() -> conversationContext.sendRequest(new AmsTCPPacket(readDeviceInfoRequest))
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(e -> conversationContext.getChannel().pipeline().fireExceptionCaught(e))
            .onError((p, e) -> conversationContext.getChannel().pipeline().fireExceptionCaught(e))
            .unwrap(AmsTCPPacket::getUserdata)
            .check(userdata -> userdata.getInvokeId() == readDeviceInfoRequest.getInvokeId())
            .only(AdsReadDeviceInfoResponse.class)
            .handle(readDeviceInfoResponse -> {
                readDeviceInfoTx.endRequest();
                future.complete(new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK));
            }));
        return future;
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
            AdsTag adsTag = (AdsTag) readRequest.getTags().stream().findFirst().orElseThrow();
            DirectAdsTag directAdsTag = resolvedTags.get(adsTag);
            // Do a normal (single item) ADS Read Request
            return singleRead(readRequest, directAdsTag);
        } else {
            // TODO: Check if the version of the remote station is at least TwinCAT v2.11 Build >= 1550 otherwise split up into single item requests.
            // Do a ADS-Sum Read Request.
            return multiRead(readRequest, resolvedTags);
        }
    }

    protected CompletableFuture<PlcReadResponse> singleRead(PlcReadRequest readRequest, DirectAdsTag directAdsTag) {
        // This can happen, if the resolution of a symbolic tag did not work.
        if(directAdsTag == null) {
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse(readRequest, Collections.singletonMap(
                readRequest.getTagNames().stream().findFirst().orElseThrow(),
                new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null))));
        }

        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        String dataTypeName = directAdsTag.getPlcDataType();
        Optional<AdsDataTypeTableEntry> dataTypeTableEntry = getDataTypeTableEntry(dataTypeName);
        if(dataTypeTableEntry.isEmpty()) {
            return CompletableFuture.failedFuture(new PlcRuntimeException("couldn't find datatype: " + dataTypeName));
        }
        long size = dataTypeTableEntry.get().getSize();

        AmsPacket amsPacket = new AdsReadRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(), 0, getInvokeId(),
            directAdsTag.getIndexGroup(), directAdsTag.getIndexOffset(), size * directAdsTag.getNumberOfElements());
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .unwrap(AmsTCPPacket::getUserdata)
            .check(userdata -> userdata.getInvokeId() == amsPacket.getInvokeId())
            .only(AdsReadResponse.class)
            .handle(response -> {
                // result metadata
                Metadata metadata = new DefaultMetadata.Builder()
                    .put(PlcMetadataKeys.RECEIVE_TIMESTAMP, System.currentTimeMillis())
                    .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.ASSUMPTION)
                    .build();
                if (response.getResult() == ReturnCode.OK) {
                    final PlcReadResponse plcReadResponse = convertToPlc4xReadResponse(readRequest, Map.of((AdsTag) readRequest.getTags().get(0), directAdsTag), response, metadata);
                    // Convert the response from the PLC into a PLC4X Response ...
                    future.complete(plcReadResponse);
                } else {
                    // TODO: Implement this correctly.
                    future.completeExceptionally(new PlcException("Unexpected return code " + response.getResult()));
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    protected CompletableFuture<PlcReadResponse> multiRead(PlcReadRequest readRequest, Map<AdsTag, DirectAdsTag> resolvedTags) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();

        // Create a list of all successfully resolved tags.
        List<AdsTag> successfullyResolvedTags = readRequest.getTagNames().stream()
            .map(tagName -> (AdsTag) readRequest.getTag(tagName))
            .filter(adsTag -> resolvedTags.get(adsTag) != null)
            .collect(Collectors.toList());

        // Calculate the expected size of the response data.
        long expectedResponseDataSize = successfullyResolvedTags.stream().mapToLong(
            adsTag -> {
                DirectAdsTag directAdsTag = resolvedTags.get(adsTag);
                String dataTypeName = directAdsTag.getPlcDataType();
                Optional<AdsDataTypeTableEntry> dataTypeTableEntry = getDataTypeTableEntry(dataTypeName);
                if(dataTypeTableEntry.isEmpty()) {
                    LOGGER.warn("couldn't find datatype: {}", dataTypeName);
                    return 0;
                }

                long size = dataTypeTableEntry.get().getSize();
                // Status code + payload size
                return 4 + (size * directAdsTag.getNumberOfElements());
            }).sum();

        // With multi-requests, the index-group is fixed and the index offset indicates the number of elements.
        AmsPacket amsPacket = new AdsReadWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_MULTIPLE_READ.getValue(), successfullyResolvedTags.size(),
            expectedResponseDataSize,
            successfullyResolvedTags.stream().map(tag -> {
                DirectAdsTag directAdsTag = resolvedTags.get(tag);
                String dataTypeName = directAdsTag.getPlcDataType();
                Optional<AdsDataTypeTableEntry> dataTypeTableEntry = getDataTypeTableEntry(dataTypeName);
                long size = dataTypeTableEntry.map(AdsDataTypeTableEntry::getSize).orElse(0L);
                return new AdsMultiRequestItemRead(
                    directAdsTag.getIndexGroup(), directAdsTag.getIndexOffset(),
                    (size * directAdsTag.getNumberOfElements()));
            }).collect(Collectors.toList()),
            null);
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .unwrap(AmsTCPPacket::getUserdata)
            .check(userdata -> userdata.getInvokeId() == amsPacket.getInvokeId())
            .only(AdsReadWriteResponse.class)
            .handle(response -> {
                Metadata metadata = new DefaultMetadata.Builder()
                    .put(PlcMetadataKeys.RECEIVE_TIMESTAMP, System.currentTimeMillis())
                    .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.ASSUMPTION)
                    .build();
                if (response.getResult() == ReturnCode.OK) {
                    final PlcReadResponse plcReadResponse = convertToPlc4xReadResponse(readRequest, resolvedTags, response, metadata);
                    // Convert the response from the PLC into a PLC4X Response ...
                    future.complete(plcReadResponse);
                } else if (response.getResult() == ReturnCode.ADSERR_DEVICE_INVALIDSIZE) {
                    future.completeExceptionally(
                        new PlcException("The parameter size was not correct (Internal error)"));
                } else {
                    future.completeExceptionally(new PlcException("Unexpected result " + response.getResult()));
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    protected PlcReadResponse convertToPlc4xReadResponse(PlcReadRequest readRequest, Map<AdsTag, DirectAdsTag> resolvedTags, AmsPacket adsData, Metadata responseMetadata) {
        ReadBuffer readBuffer = null;
        Map<String, Metadata> metadata = new HashMap<>();
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();

        // Read the response codes first
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
                    PlcTag tag = readRequest.getTag(tagName);
                    if(resolvedTags.get((AdsTag) tag) != null) {
                        ReturnCode result = ReturnCode.enumForValue(readBuffer.readUnsignedLong(32));
                        responseCodes.put(tagName, parsePlcResponseCode(result));
                    } else {
                        responseCodes.put(tagName, PlcResponseCode.INVALID_ADDRESS);
                    }
                } catch (ParseException e) {
                    responseCodes.put(tagName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }

        // Read the values next
        if (readBuffer != null) {
            Map<String, PlcResponseItem<PlcValue>> values = new HashMap<>();
            for (String tagName : readRequest.getTagNames()) {
                metadata.put(tagName, new DefaultMetadata.Builder(responseMetadata).build());
                // If the response-code was anything but OK, we don't need to parse the payload.
                if(responseCodes.get(tagName) != PlcResponseCode.OK) {
                    values.put(tagName, new DefaultPlcResponseItem<>(responseCodes.get(tagName), null));
                }
                // If the response-code was ok, parse the data returned.
                else {
                    DirectAdsTag directAdsTag = resolvedTags.get((AdsTag) readRequest.getTag(tagName));
                    values.put(tagName, parseResponseItem(directAdsTag, readBuffer));
                }
            }
            return new DefaultPlcReadResponse(readRequest, values, metadata);
        }
        return null;
    }

    private PlcResponseCode parsePlcResponseCode(ReturnCode adsResult) {
        if (adsResult == ReturnCode.OK) {
            return PlcResponseCode.OK;
        } else if(adsResult == ReturnCode.ADSERR_DEVICE_SYMBOLNOTFOUND) {
            return PlcResponseCode.INVALID_ADDRESS;
        } else {
            return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    private PlcResponseItem<PlcValue> parseResponseItem(DirectAdsTag tag, ReadBuffer readBuffer) {
        try {
            String dataTypeName = tag.getPlcDataType();
            Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(dataTypeName);
            if(dataTypeTableEntryOptional.isEmpty()) {
                return new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null);
            }
            AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTableEntryOptional.get();
            PlcValueType plcValueType = getPlcValueTypeForAdsDataType(adsDataTypeTableEntry);

            int strLen = 0;
            if (tag instanceof DirectAdsStringTag) {
                strLen = ((DirectAdsStringTag) tag).getStringLength();
            }
            final int stringLength = strLen;
            if (tag.getNumberOfElements() == 1) {
                ReadBufferByteBased readBufferByteBased = ((ReadBufferByteBased) readBuffer);
                // Sometimes the ADS device just sends shorter strings than we asked for.
                int remainingBytes = readBufferByteBased.getTotalBytes() - readBufferByteBased.getPos();
                final int singleStringLength = Math.min(remainingBytes - 1, stringLength);
                return new DefaultPlcResponseItem<>(PlcResponseCode.OK, parsePlcValue(plcValueType, adsDataTypeTableEntry, singleStringLength, readBuffer));
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
                return new DefaultPlcResponseItem<>(PlcResponseCode.OK, DefaultPlcValueHandler.of(tag, resultItems));
            }
        } catch (Exception e) {
            LOGGER.warn(String.format("Error parsing tag item of type: '%s'", tag.getPlcDataType()), e);
            return new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null);
        }
    }

    private PlcValue parsePlcValue(PlcValueType plcValueType, AdsDataTypeTableEntry adsDataTypeTableEntry, int stringLength, ReadBuffer readBuffer) throws ParseException {
        switch (plcValueType) {
            case Struct:
                Map<String, PlcValue> properties = new HashMap<>();
                int startPos = readBuffer.getPos();
                int curPos = 0;
                for (AdsDataTypeTableChildEntry child : adsDataTypeTableEntry.getChildren()) {
                    // In some cases the starting position of the data is not where we are expecting it.
                    // So we need to skip some bytes.
                    if (child.getOffset() > curPos) {
                        long skipBytes = child.getOffset() - curPos;
                        for (long i = 0; i < skipBytes; i++) {
                            readBuffer.readByte();
                        }
                    }

                    String propertyName = child.getPropertyName();
                    Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(child.getDataTypeName());
                    if(dataTypeTableEntryOptional.isEmpty()) {
                        throw new ParseException(String.format("couldn't find datatype: %s", child.getDataTypeName()));
                    }
                    AdsDataTypeTableEntry propertyDataTypeTableEntry = dataTypeTableEntryOptional.get();
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
            Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(dataTypeName);
            if(dataTypeTableEntryOptional.isEmpty()) {
                throw new ParseException(String.format("couldn't find datatype: %s", dataTypeName));
            }
            AdsDataTypeTableEntry elementDataTypeTableEntry = dataTypeTableEntryOptional.get();
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
            AdsTag adsTag = (AdsTag) writeRequest.getTags().stream().findFirst().orElseThrow();
            DirectAdsTag directAdsTag = resolvedTags.get(adsTag);
            // Do a normal (single item) ADS Write Request
            return singleWrite(writeRequest, directAdsTag);
        } else {
            // TODO: Check if the version of the remote station is at least TwinCAT v2.11 Build >= 1550 otherwise split up into single item requests.
            // Do a ADS-Sum Read Request.
            return multiWrite(writeRequest, resolvedTags);
        }
    }

    protected CompletableFuture<PlcWriteResponse> singleWrite(PlcWriteRequest writeRequest, DirectAdsTag directAdsTag) {
        // This can happen, if the resolution of a symbolic tag did not work.
        if(directAdsTag == null) {
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, Collections.singletonMap(
                writeRequest.getTagNames().stream().findFirst().orElseThrow(),
                PlcResponseCode.INVALID_ADDRESS)));
        }

        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();

        final String tagName = writeRequest.getTagNames().iterator().next();
        final PlcValue plcValue = writeRequest.getPlcValue(tagName);

        // Serialize the value
        byte[] serializedValue;
        try {
            serializedValue = serializePlcValue(plcValue, directAdsTag.getPlcDataType());
        } catch (Exception e) {
            future.completeExceptionally(new PlcException("Error serializing data tag value for tag '" + tagName + "'", e));
            return future;
        }

        AmsPacket amsPacket = new AdsWriteRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
            configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
            0, getInvokeId(), directAdsTag.getIndexGroup(), directAdsTag.getIndexOffset(), serializedValue);
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .unwrap(AmsTCPPacket::getUserdata)
            .check(userdata -> userdata.getInvokeId() == amsPacket.getInvokeId())
            .only(AdsWriteResponse.class)
            .handle(response -> {
                // result metadata
                Metadata eventMetadata = new DefaultMetadata.Builder()
                    .put(PlcMetadataKeys.RECEIVE_TIMESTAMP, System.currentTimeMillis())
                    .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.ASSUMPTION)
                    .build();
                if (response.getResult() == ReturnCode.OK) {
                    final PlcWriteResponse plcWriteResponse = convertToPlc4xWriteResponse(writeRequest, Collections.singletonMap((AdsTag) writeRequest.getTag(tagName), directAdsTag), response, eventMetadata);
                    // Convert the response from the PLC into a PLC4X Response ...
                    future.complete(plcWriteResponse);
                } else {
                    // TODO: Implement this correctly.
                    future.completeExceptionally(new PlcException("Unexpected return code " + response.getResult()));
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    protected CompletableFuture<PlcWriteResponse> multiWrite(PlcWriteRequest writeRequest, Map<AdsTag, DirectAdsTag> resolvedTags) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();

        int numTags = writeRequest.getTags().size();

        // Serialize all tags.
        List<byte[]> serializedTags = new ArrayList<>(numTags);
        Map<DirectAdsTag, AdsDataTypeTableEntry> tagDatatypes = new LinkedHashMap<>(numTags);
        for (String tagName : writeRequest.getTagNames()) {
            final AdsTag adsTag = (AdsTag) writeRequest.getTag(tagName);
            // Skip invalid addresses.
            if(resolvedTags.get(adsTag) == null) {
                continue;
            }
            final DirectAdsTag directAdsTag = resolvedTags.get(adsTag);
            final PlcValue plcValue = writeRequest.getPlcValue(tagName);
            Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(directAdsTag.getPlcDataType());
            if(dataTypeTableEntryOptional.isPresent()) {
                AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTableEntryOptional.get();
                try {
                    byte[] serializedValue = serializePlcValue(plcValue, directAdsTag.getPlcDataType());
                    serializedTags.add(serializedValue);
                    tagDatatypes.put(directAdsTag, adsDataTypeTableEntry);
                } catch (Exception e) {
                    future.completeExceptionally(new PlcException("Error serializing data", e));
                    return future;
                }
            } else {
                future.completeExceptionally(new PlcException(String.format(
                    "couldn't find datatype: %s", directAdsTag.getPlcDataType())));
                return future;
            }
        }

        // Calculate the size of all serialized tags together.
        int serializedSize = serializedTags.stream().mapToInt(
            serializedTag -> serializedTag.length).sum();

        // Copy all serialized tags into one buffer.
        // This is intentionally not "LittleEndian" as we're just concatenating the serialized values.
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
            0, getInvokeId(), ReservedIndexGroups.ADSIGRP_MULTIPLE_WRITE.getValue(), serializedTags.size(),
            (long) numTags * 4,
            tagDatatypes.entrySet().stream().map(entry -> new AdsMultiRequestItemWrite(
                    entry.getKey().getIndexGroup(), entry.getKey().getIndexOffset(),
                    entry.getValue().getSize()))
                .collect(Collectors.toList()), writeBuffer.getBytes());
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .unwrap(AmsTCPPacket::getUserdata)
            .check(userdata -> userdata.getInvokeId() == amsPacket.getInvokeId())
            .only(AdsReadWriteResponse.class)
            .handle(response -> {
                // result metadata
                Metadata eventMetadata = new DefaultMetadata.Builder()
                    .put(PlcMetadataKeys.RECEIVE_TIMESTAMP, System.currentTimeMillis())
                    .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.ASSUMPTION)
                    .build();

                if (response.getResult() == ReturnCode.OK) {
                    final PlcWriteResponse plcWriteResponse = convertToPlc4xWriteResponse(writeRequest, resolvedTags, response, eventMetadata);
                    // Convert the response from the PLC into a PLC4X Response ...
                    future.complete(plcWriteResponse);
                } else {
                    future.completeExceptionally(new PlcException("Unexpected result " + response.getResult()));
                }
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;
    }

    protected byte[] serializePlcValue(PlcValue plcValue, String datatypeName) throws SerializationException {
        // Get the data type, allocate enough memory and serialize the value based on the
        // structure defined by the data type.
        Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(datatypeName);
        if(dataTypeTableEntryOptional.isEmpty()) {
            throw new SerializationException("Could not find data type: " + datatypeName);
        }
        AdsDataTypeTableEntry dataType = dataTypeTableEntryOptional.get();
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased((int) dataType.getSize(), ByteOrder.LITTLE_ENDIAN);
        List<AdsDataTypeArrayInfo> arrayInfo = dataType.getArrayInfo();
        serializeInternal(plcValue, dataType, arrayInfo, writeBuffer);
        return writeBuffer.getBytes();
    }

    protected void serializeInternal(PlcValue contextValue,
                                     AdsDataTypeTableEntry dataType,
                                     List<AdsDataTypeArrayInfo> arrayInfo,
                                     WriteBufferByteBased writeBuffer) throws SerializationException {

        // An array type: Recursively iterate over the elements
        if (!arrayInfo.isEmpty()) {
            if (!contextValue.isList()) {
                throw new SerializationException("Expected a PlcList, but got a " + contextValue.getPlcValueType().name());
            }
            AdsDataTypeArrayInfo curArrayLevel = arrayInfo.get(0);
            List<? extends PlcValue> list = contextValue.getList();
            if (curArrayLevel.getNumElements() != list.size()) {
                throw new SerializationException(String.format(
                    "Expected a PlcList of size %d, but got one of size %d", curArrayLevel.getNumElements(), list.size()));
            }
            Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(dataType.getSimpleTypeName());
            if(dataTypeTableEntryOptional.isEmpty()) {
                throw new SerializationException("Could not find data type: " + dataType.getSimpleTypeName());
            }
            AdsDataTypeTableEntry childDataType = dataTypeTableEntryOptional.get();
            for (PlcValue plcValue : list) {
                serializeInternal(plcValue, childDataType, arrayInfo.subList(1, arrayInfo.size()), writeBuffer);
            }
        }

        // A complex type
        else if (!dataType.getChildren().isEmpty()) {
            if (!contextValue.isStruct()) {
                throw new SerializationException("Expected a PlcStruct, but got a " + contextValue.getPlcValueType().name());
            }
            PlcStruct plcStruct = (PlcStruct) contextValue;
            int startPos = writeBuffer.getPos();
            int curPos = 0;
            for (AdsDataTypeTableChildEntry child : dataType.getChildren()) {
                Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(child.getDataTypeName());
                if(dataTypeTableEntryOptional.isEmpty()) {
                    throw new SerializationException("Could not find data type: " + child.getDataTypeName());
                }
                AdsDataTypeTableEntry childDataType = dataTypeTableEntryOptional.get();
                if (!plcStruct.hasKey(child.getPropertyName())) {
                    throw new SerializationException("PlcStruct is missing a child with the name " + child.getPropertyName());
                }
                // In some cases the starting position of the data is not where we are expecting it.
                // So we need to add some fill-bytes.
                if (child.getOffset() > curPos) {
                    long fillBytes = child.getOffset() - curPos;
                    for(long i = 0; i < fillBytes; i++) {
                        writeBuffer.writeByte("fillByte", (byte) 0x00);
                    }
                }

                PlcValue childValue = plcStruct.getValue(child.getPropertyName());
                serializeInternal(childValue, childDataType, childDataType.getArrayInfo(), writeBuffer);

                curPos = writeBuffer.getPos() - startPos;
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

    protected PlcWriteResponse convertToPlc4xWriteResponse(PlcWriteRequest writeRequest, Map<AdsTag, DirectAdsTag> resolvedTags, AmsPacket adsData, Metadata eventMtadata) {
        Map<String, PlcResponseCode> responseCodes = new HashMap<>();
        Map<String, Metadata> metadata = new HashMap<>();
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
                // result metadata
                metadata.put(tagName, eventMtadata);

                AdsTag adsTag = (AdsTag) writeRequest.getTag(tagName);
                // Skip invalid addresses.
                if(resolvedTags.get(adsTag) == null) {
                    responseCodes.put(tagName, PlcResponseCode.INVALID_ADDRESS);
                    continue;
                }
                try {
                    final ReturnCode result = ReturnCode.enumForValue(readBuffer.readUnsignedLong(32));
                    responseCodes.put(tagName, parsePlcResponseCode(result));
                } catch (ParseException e) {
                    responseCodes.put(tagName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }

        return new DefaultPlcWriteResponse(writeRequest, responseCodes, metadata);
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
                String dataTypeName = resolvedTags.get((AdsTag) tag.getTag()).getPlcDataType();
                Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(dataTypeName);
                if(dataTypeTableEntryOptional.isEmpty()) {
                    return null;
                }
                AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTableEntryOptional.get();
                DirectAdsTag directAdsTag = getDirectAdsTagForSymbolicName(tag.getTag());
                // TODO: We should implement multi-dimensional arrays here ...
                int numberOfElements = (tag.getArrayInfo().isEmpty()) ? 1 : tag.getArrayInfo().get(0).getSize();
                return new AmsTCPPacket(new AdsAddDeviceNotificationRequest(configuration.getTargetAmsNetId(), configuration.getTargetAmsPort(),
                    configuration.getSourceAmsNetId(), configuration.getSourceAmsPort(),
                    0, getInvokeId(),
                    directAdsTag.getIndexGroup(),
                    directAdsTag.getIndexOffset(),
                    adsDataTypeTableEntry.getSize() * numberOfElements,
                    tag.getPlcSubscriptionType() == PlcSubscriptionType.CYCLIC ? AdsTransMode.CYCLIC : AdsTransMode.ON_CHANGE, // if it's not cyclic, it's on change or event
                    0, // there is no api for that yet
                    tag.getDuration().orElse(Duration.ZERO).toMillis()));
            }).filter(Objects::nonNull)
            .collect(Collectors.toList());

        Map<String, PlcResponseItem<PlcSubscriptionHandle>> responses = new HashMap<>();

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
                                          Map<String, PlcResponseItem<PlcSubscriptionHandle>> responses,
                                          CompletableFuture<PlcSubscriptionResponse> future,
                                          Iterator<AmsTCPPacket> amsTCPPackets,
                                          RequestTransactionManager.RequestTransaction transaction) {
        return () -> {
            AmsTCPPacket packet = amsTCPPackets.next();
            boolean hasMorePackets = amsTCPPackets.hasNext();
            String tagName = tagNames.next();
            conversationContext.sendRequest(packet)
                .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .unwrap(AmsTCPPacket::getUserdata)
                .check(userdata -> userdata.getInvokeId() == packet.getUserdata().getInvokeId())
                .only(AdsAddDeviceNotificationResponse.class)
                .handle(response -> {
                    if (response.getResult() == ReturnCode.OK) {
                        DefaultPlcSubscriptionTag subscriptionTag = (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(tagName);

                        String dataTypeName = resolvedTags.get((AdsTag) subscriptionTag.getTag()).getPlcDataType();
                        Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(dataTypeName);
                        if(dataTypeTableEntryOptional.isEmpty()) {
                            future.completeExceptionally(new PlcRuntimeException("Could not find data type: " + dataTypeName));
                            return;
                        }
                        AdsDataTypeTableEntry adsDataTypeTableEntry = dataTypeTableEntryOptional.get();
                        // Collect notification handle from individual response.
                        responses.put(tagName, new DefaultPlcResponseItem<>(
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
            conversationContext.sendRequest(packet)
                .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .unwrap(AmsTCPPacket::getUserdata)
                .check(userdata -> userdata.getInvokeId() == packet.getUserdata().getInvokeId())
                .only(AdsDeleteDeviceNotificationResponse.class)
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
            long receiveTs = System.currentTimeMillis();
            for (AdsStampHeader stamp : stamps) {
                // convert Windows FILETIME format to unix epoch
                long unixEpochTimestamp = stamp.getTimestamp().divide(BigInteger.valueOf(10000L)).longValue() - 11644473600000L;
                // result metadata
                Metadata eventMetadata = new DefaultMetadata.Builder()
                    .put(PlcMetadataKeys.RECEIVE_TIMESTAMP, receiveTs)
                    .put(PlcMetadataKeys.TIMESTAMP, unixEpochTimestamp)
                    .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.SOFTWARE)
                    .build();
                List<AdsNotificationSample> samples = stamp.getAdsNotificationSamples();
                for (AdsNotificationSample sample : samples) {
                    long handle = sample.getNotificationHandle();
                    for (DefaultPlcConsumerRegistration registration : consumers.keySet()) {
                        for (PlcSubscriptionHandle subscriptionHandle : registration.getSubscriptionHandles()) {
                            if (subscriptionHandle instanceof AdsSubscriptionHandle) {
                                AdsSubscriptionHandle adsHandle = (AdsSubscriptionHandle) subscriptionHandle;
                                if (adsHandle.getNotificationHandle() == handle) {
                                    Map<String, Metadata> metadata = new HashMap<>();
                                    Instant timestamp = Instant.ofEpochMilli(unixEpochTimestamp);
                                    DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(timestamp, convertSampleToPlc4XResult(adsHandle, sample.getData(), metadata, eventMetadata));
                                    consumers.get(registration).accept(event);
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    private Map<String, PlcResponseItem<PlcValue>> convertSampleToPlc4XResult(AdsSubscriptionHandle subscriptionHandle, byte[] data, Map<String, Metadata> tagMetadata, Metadata metadata) throws
        ParseException {
        Map<String, PlcResponseItem<PlcValue>> values = new HashMap<>();
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN);
        values.put(subscriptionHandle.getTagName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK,
            DataItem.staticParse(readBuffer, getPlcValueTypeForAdsDataType(subscriptionHandle.getAdsDataType()), data.length)));
        tagMetadata.put(subscriptionHandle.getTagName(), new DefaultMetadata.Builder(metadata).build());
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
                // We will ignore the direct output of the global future as we're more interested in the
                // individual sub-future results (Otherwise one failing resolution will fail all resolutions)
                Map<AdsTag, DirectAdsTag> directAdsTagMapping = new HashMap<>(tags.size());
                for (PlcTag tag : tags) {
                    if (tag instanceof SymbolicAdsTag) {
                        // If we needed to actively resolve the address ...
                        if(pendingResolutionRequests.containsKey(tag) && pendingResolutionRequests.get(tag).isCompletedExceptionally()) {
                            directAdsTagMapping.put((AdsTag) tag, null);
                        } else {
                            directAdsTagMapping.put((AdsTag) tag, getDirectAdsTagForSymbolicName(tag));
                        }
                    } else {
                        directAdsTagMapping.put((AdsTag) tag, (DirectAdsTag) tag);
                    }
                }
                return future.complete(directAdsTagMapping);
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
        transaction.submit(() -> conversationContext.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
            .unwrap(AmsTCPPacket::getUserdata)
            .only(AdsReadWriteResponse.class)
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
                4, symbolicAdsTag.getSymbolicAddress().length())).collect(Collectors.toList()), addressData);
        AmsTCPPacket amsTCPPacket = new AmsTCPPacket(amsPacket);

        // Start a new request-transaction (Is ended in the response-handler)
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        transaction.submit(() -> conversationContext.sendRequest(amsTCPPacket)
            .expectResponse(AmsTCPPacket.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(responseAmsPacket -> responseAmsPacket.getUserdata().getInvokeId() == amsPacket.getInvokeId())
            .unwrap(AmsTCPPacket::getUserdata)
            .only(AdsReadWriteResponse.class)
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
                        ReturnCode returnCode = ReturnCode.enumForValue(returnCodes.get(symbolicAdsTag));
                        // 0 indicates a successful result.
                        if (returnCode == ReturnCode.OK) {
                            // Read the handle.
                            long handle = readBuffer.readUnsignedLong(32);

                            // TODO: Finish the parsing of the response information and possibly the reading of
                            //  datatype information for the current tag.
                            /*DirectAdsTag directAdsTag = new DirectAdsTag(
                                ReservedIndexGroups.ADSIGRP_SYM_VALBYHND.getValue(), handle,
                                symbolicAdsTag.getAdsDataTypeName(), symbolicAdsTag.getNumberOfElements());
                            symbolicTagMapping.put(symbolicAdsTag, directAdsTag);*/
                        }

                        // The user provided an invalid address
                        else if(returnCode == ReturnCode.ADSERR_DEVICE_SYMBOLNOTFOUND) {
                            pendingResolutionRequests.put(symbolicAdsTag, CompletableFuture.failedFuture(
                                new PlcInvalidTagException("Could not resolve tag " + symbolicAdsTag.getSymbolicAddress())));
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
            Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(adsSymbolTableEntry.getDataTypeName());
            if(dataTypeTableEntryOptional.isEmpty()) {
                return null;
            }
            AdsDataTypeTableEntry dataTypeTableEntry = dataTypeTableEntryOptional.get();
            return new DirectAdsTag(adsSymbolTableEntry.getGroup(), adsSymbolTableEntry.getOffset(),
                dataTypeTableEntry.getDataTypeName(), dataTypeTableEntry.getArrayDimensions());
        }
        // Otherwise we'll have to crawl through the dataType definitions.
        else {
            String symbolName = addressParts[0] + "." + addressParts[1];
            // We can't find it, so we need to resolve it.
            if (!symbolTable.containsKey(symbolName)) {
                return null;
            }
            AdsSymbolTableEntry adsSymbolTableEntry = symbolTable.get(symbolName);
            Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(adsSymbolTableEntry.getDataTypeName());
            if(dataTypeTableEntryOptional.isEmpty()) {
                return null;
            }
            AdsDataTypeTableEntry dataTypeTableEntry = dataTypeTableEntryOptional.get();
            return resolveDirectAdsTagForSymbolicNameFromDataType(
                Arrays.asList(addressParts).subList(2, addressParts.length),
                adsSymbolTableEntry.getGroup(), adsSymbolTableEntry.getOffset(), dataTypeTableEntry);
        }
    }

    protected DirectAdsTag resolveDirectAdsTagForSymbolicNameFromDataType(List<String> remainingAddressParts, long currentGroup, long currentOffset, AdsDataTypeTableEntry adsDataTypeTableEntry) {
        if (remainingAddressParts.isEmpty()) {
            // TODO: Implement the Array support
            if (adsDataTypeTableEntry.getDataType() == AdsDataType.CHAR.getValue()) {
                int stringLength = (int) adsDataTypeTableEntry.getSize() - 1;
                return new DirectAdsStringTag(currentGroup, currentOffset, adsDataTypeTableEntry.getDataTypeName(), stringLength, 1);
            } else if (adsDataTypeTableEntry.getDataType() == AdsDataType.WCHAR.getValue()) {
                int stringLength = (int) (adsDataTypeTableEntry.getSize() - 2) / 2;
                return new DirectAdsStringTag(currentGroup, currentOffset, adsDataTypeTableEntry.getDataTypeName(), stringLength, 1);
            } else {
                return new DirectAdsTag(currentGroup, currentOffset, adsDataTypeTableEntry.getDataTypeName(), 1);
            }
        }

        // Go through all children looking for a matching one.
        for (AdsDataTypeTableChildEntry child : adsDataTypeTableEntry.getChildren()) {
            if (child.getPropertyName().equals(remainingAddressParts.get(0))) {
                Optional<AdsDataTypeTableEntry> dataTypeTableEntryOptional = getDataTypeTableEntry(child.getDataTypeName());
                if(dataTypeTableEntryOptional.isEmpty()) {
                    throw new PlcRuntimeException("Could not resolve data type " + child.getDataTypeName());
                }
                AdsDataTypeTableEntry childAdsDataTypeTableEntry = dataTypeTableEntryOptional.get();
                return resolveDirectAdsTagForSymbolicNameFromDataType(
                    remainingAddressParts.subList(1, remainingAddressParts.size()),
                    currentGroup, currentOffset + child.getOffset(), childAdsDataTypeTableEntry);
            }
        }

        throw new PlcRuntimeException(String.format("Couldn't find child with name '%s' for type '%s'",
            remainingAddressParts.get(0), adsDataTypeTableEntry.getDataTypeName()));
    }

    protected PlcValueType getPlcValueTypeForAdsDataTypeForBrowse(AdsDataTypeTableEntry dataTypeTableEntry) {
        String dataTypeName = (!dataTypeTableEntry.getSimpleTypeName().isEmpty() && !dataTypeTableEntry.getDataTypeName().equals("BOOL")) ?
            dataTypeTableEntry.getSimpleTypeName() : dataTypeTableEntry.getDataTypeName();
        if (dataTypeName.startsWith("STRING(")) {
            dataTypeName = "STRING";
        } else if (dataTypeName.startsWith("WSTRING(")) {
            dataTypeName = "WSTRING";
        }
        // First check, if this is a primitive type.
        try {
            return PlcValueType.valueOf(dataTypeName);
        } catch (IllegalArgumentException e) {
            return PlcValueType.Struct;
        }
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
            if (dataTypeTableEntry.getChildren().isEmpty()) {
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

    protected Optional<AdsDataTypeTableEntry> getDataTypeTableEntry(String name) {
        if(dataTypeTable.containsKey(name)) {
            return Optional.of(dataTypeTable.get(name));
        }
        try {
            AdsDataType adsDataType;
            int numBytes;

            if (name.startsWith("STRING(")) {
                adsDataType = AdsDataType.valueOf("CHAR");
                numBytes = Integer.parseInt(name.substring(7, name.length() - 1)) + 1;
            }
            else if (name.startsWith("WSTRING(")) {
                adsDataType = AdsDataType.valueOf("WCHAR");
                numBytes = Integer.parseInt(name.substring(8, name.length() - 1)) * 2 + 2;
            }
            else {
                adsDataType = AdsDataType.valueOf(name);
                numBytes = adsDataType.getNumBytes();
            }

            // !It seems that the dataType value differs from system to system,
            // !However, we never really seem to use that value, so I would say it doesn't matter.
            return Optional.of(new AdsDataTypeTableEntry(
                128, 1, 0, 0, numBytes, 0,
                adsDataType.getValue(), 0, 0, 0,
                name, "", "",
                Collections.emptyList(), Collections.emptyList(), new byte[0]));
        } catch (IllegalArgumentException e) {
            return Optional.empty();
        }
    }

    protected byte[] getNullByteTerminatedArray(String value) {
        byte[] valueBytes = value.getBytes();
        byte[] nullTerminatedBytes = new byte[valueBytes.length + 1];
        System.arraycopy(valueBytes, 0, nullTerminatedBytes, 0, valueBytes.length);
        return nullTerminatedBytes;
    }

}
