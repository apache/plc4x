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
package org.apache.plc4x.java.knxnetip;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
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
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpConfiguration;
import org.apache.plc4x.java.knxnetip.ets.EtsParser;
import org.apache.plc4x.java.knxnetip.ets.model.EtsModel;
import org.apache.plc4x.java.knxnetip.ets.model.GroupAddress;
import org.apache.plc4x.java.knxnetip.model.KnxNetIpSubscriptionHandle;
import org.apache.plc4x.java.knxnetip.readwrite.Apdu;
import org.apache.plc4x.java.knxnetip.readwrite.ApduData;
import org.apache.plc4x.java.knxnetip.readwrite.ApduDataContainer;
import org.apache.plc4x.java.knxnetip.readwrite.ApduDataGroupValueRead;
import org.apache.plc4x.java.knxnetip.readwrite.ApduDataGroupValueResponse;
import org.apache.plc4x.java.knxnetip.readwrite.ApduDataGroupValueWrite;
import org.apache.plc4x.java.knxnetip.readwrite.KnxDatapointMainType;
import org.apache.plc4x.java.knxnetip.readwrite.KnxDatapointType;
import org.apache.plc4x.java.knxnetip.readwrite.CEMIPriority;
import org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequest;
import org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequestInformationTunnelConnection;
import org.apache.plc4x.java.knxnetip.readwrite.ConnectionResponse;
import org.apache.plc4x.java.knxnetip.readwrite.ConnectionResponseDataBlockTunnelConnection;
import org.apache.plc4x.java.knxnetip.readwrite.ConnectionStateRequest;
import org.apache.plc4x.java.knxnetip.readwrite.ConnectionStateResponse;
import org.apache.plc4x.java.knxnetip.readwrite.DisconnectRequest;
import org.apache.plc4x.java.knxnetip.readwrite.HPAIControlEndpoint;
import org.apache.plc4x.java.knxnetip.readwrite.HPAIDataEndpoint;
import org.apache.plc4x.java.knxnetip.readwrite.HPAIDiscoveryEndpoint;
import org.apache.plc4x.java.knxnetip.readwrite.HostProtocolCode;
import org.apache.plc4x.java.knxnetip.readwrite.IPAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KnxAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KnxDatapoint;
import org.apache.plc4x.java.knxnetip.readwrite.KnxGroupAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KnxGroupAddress2Level;
import org.apache.plc4x.java.knxnetip.readwrite.KnxGroupAddress3Level;
import org.apache.plc4x.java.knxnetip.readwrite.KnxGroupAddressFreeLevel;
import org.apache.plc4x.java.knxnetip.readwrite.KnxLayer;
import org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpMessage;
import org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpTunneling;
import org.apache.plc4x.java.knxnetip.readwrite.LDataExtended;
import org.apache.plc4x.java.knxnetip.readwrite.LDataFrame;
import org.apache.plc4x.java.knxnetip.readwrite.LDataInd;
import org.apache.plc4x.java.knxnetip.readwrite.LDataReq;
import org.apache.plc4x.java.knxnetip.readwrite.SearchRequest;
import org.apache.plc4x.java.knxnetip.readwrite.SearchResponse;
import org.apache.plc4x.java.knxnetip.readwrite.ServiceId;
import org.apache.plc4x.java.knxnetip.readwrite.Status;
import org.apache.plc4x.java.knxnetip.readwrite.TunnelingRequest;
import org.apache.plc4x.java.knxnetip.readwrite.TunnelingRequestDataBlock;
import org.apache.plc4x.java.knxnetip.readwrite.TunnelingResponse;
import org.apache.plc4x.java.knxnetip.readwrite.TunnelingResponseDataBlock;
import org.apache.plc4x.java.knxnetip.tag.KnxNetIpTag;
import org.apache.plc4x.java.knxnetip.tag.KnxNetIpTagHandler;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcPingResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.transport.udp.UdpTransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.values.PlcValueAdapter;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetSocketAddress;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.Comparator;
import java.util.EnumSet;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeMap;
import java.util.concurrent.CompletableFuture;
import java.util.regex.Pattern;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

/**
 * KNXNet/IP UDP connection. Direct port of the legacy {@code KnxNetIpProtocolLogic}
 * and {@code KnxNetIpDriverContext} to the SPI3 {@link ConnectionBase} model.
 *
 * <p>The protocol is request-driven on the control plane and push-driven on the
 * data plane. {@code onConnect()} runs a synchronous {@code SearchRequest →
 * ConnectionRequest} handshake against the gateway, then schedules a 60-second
 * {@code ConnectionStateRequest} heartbeat to keep the tunnel open. Once
 * connected the gateway forwards every group write on the KNX bus as an
 * {@code LDataInd}; this class decodes those (optionally through an ETS
 * project model for typed values) and dispatches matching events to any
 * registered subscription consumers.</p>
 *
 * <p>Writes are wrapped in a {@code TunnelingRequest} and matched to the
 * gateway's {@code TunnelingResponse} ack by sequence counter.</p>
 */
public class KnxNetIpConnection extends ConnectionBase<KnxNetIpConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(KnxNetIpConnection.class);

    /** How often to refresh the tunnel with a ConnectionStateRequest. */
    private static final long HEARTBEAT_INTERVAL_MS = 60_000L;

    private KnxNetIpMessageCodec messageCodec;
    private final ScheduledExecutorService heartbeatExecutor =
        Executors.newSingleThreadScheduledExecutor(r -> {
            Thread t = new Thread(r, "KnxNetIp-Heartbeat");
            t.setDaemon(true);
            return t;
        });
    private ScheduledFuture<?> heartbeatTask;

    /** Session state — populated during the handshake. */
    private volatile boolean handshakeComplete = false;
    private IPAddress localIPAddress;
    private int localPort;
    private KnxAddress gatewayAddress;
    private String gatewayName;
    private short communicationChannelId;
    private KnxAddress clientKnxAddress;
    private byte groupAddressType;
    private KnxLayer tunnelConnectionType;
    private EtsModel etsModel;

    /** Pending search-response future during the handshake. */
    private volatile CompletableFuture<SearchResponse> pendingSearch;
    /** Pending connection-response future during the handshake. */
    private volatile CompletableFuture<ConnectionResponse> pendingConnection;
    /** Pending connection-state-response future (heartbeat + ping). */
    private final Map<Short, CompletableFuture<ConnectionStateResponse>> pendingConnectionState =
        new ConcurrentHashMap<>();
    /** Pending tunneling-response futures keyed by sequence counter. */
    private final Map<Short, CompletableFuture<TunnelingResponse>> pendingTunneling =
        new ConcurrentHashMap<>();
    /**
     * Pending {@code GroupValueRead} requests keyed by destination group
     * address (as the textual {@code main/middle/sub} form). KNX reads work by
     * sending a read request to the bus and waiting for some other device to
     * answer with a {@code GroupValueResponse} indication routed to that same
     * GA. We complete the future when the matching response indication shows
     * up in {@link #handleTunnelingRequest}.
     */
    private final Map<String, CompletableFuture<byte[]>> pendingReads = new ConcurrentHashMap<>();

    /** Outbound sequence counter for TunnelingRequest. */
    private final AtomicInteger sequenceCounter = new AtomicInteger(0);

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    public KnxNetIpConnection(KnxNetIpConfiguration configuration,
                              TransportInstance<?> transportInstance,
                              AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
        loadEtsModel(configuration);
        tunnelConnectionType = KnxLayer.valueOf("TUNNEL_" + configuration.getConnectionType());
    }

    private void loadEtsModel(KnxNetIpConfiguration configuration) {
        if (configuration.knxprojFile == null) {
            this.groupAddressType = (byte) configuration.groupAddressNumLevels;
            return;
        }
        if (!configuration.knxprojFile.exists() || !configuration.knxprojFile.isFile()) {
            throw new PlcRuntimeException(String.format(
                "File specified with 'knxproj-file-path' does not exist or is not a file: '%s'",
                configuration.knxprojFile));
        }
        this.etsModel = new EtsParser().parse(configuration.knxprojFile, configuration.knxprojPassword);
        this.groupAddressType = etsModel.getGroupAddressType();
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new KnxNetIpTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    protected int getMaxConcurrentRequests() {
        return 1;
    }

    @Override
    public boolean isConnected() {
        return handshakeComplete && messageCodec != null && messageCodec.isOpen();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Connect / close
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new KnxNetIpMessageCodec(transportInstance, this::handleIncomingMessage);

        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming KNXnet/IP data", e);
            }
        });

        // The UDP transport exposes its bound local socket address via the
        // TransportInstance — that's the endpoint the gateway will use as
        // the data endpoint for replies, so we need to give it back to the
        // gateway in the HPAI fields.
        InetSocketAddress localAddress = resolveLocalAddress();
        this.localIPAddress = new IPAddress(localAddress.getAddress().getAddress());
        this.localPort = localAddress.getPort();

        try {
            doHandshake().get(getConfiguration().getRequestTimeout() * 3L, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("Error during KNXnet/IP connect handshake", e);
        }

        handshakeComplete = true;
        startHeartbeat();

        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT,
                String.format("KNXnet/IP connection established to %s (channel id %d)",
                    gatewayName, communicationChannelId));
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    private InetSocketAddress resolveLocalAddress() throws PlcConnectionException {
        if (!(transportInstance instanceof UdpTransportInstance udp)) {
            throw new PlcConnectionException(
                "KNXnet/IP requires a UDP transport, got " + transportInstance.getClass().getName());
        }
        return udp.getLocalAddress().orElseThrow(() -> new PlcConnectionException(
            "UDP transport has no bound local address (socket not connected?)"));
    }

    private CompletableFuture<Void> doHandshake() {
        pendingSearch = new CompletableFuture<>();
        SearchRequest searchRequest = new SearchRequest(
            new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));

        return sendMessage(searchRequest, pendingSearch)
            .thenCompose(this::sendConnectionRequest);
    }

    private CompletableFuture<Void> sendConnectionRequest(SearchResponse searchResponse) {
        // Make sure the gateway actually supports tunneling.
        ServiceId tunnelingService = searchResponse.getDibSuppSvcFamilies().getServiceIds().stream()
            .filter(serviceId -> serviceId instanceof KnxNetIpTunneling)
            .findFirst()
            .orElse(null);
        if (tunnelingService == null) {
            return CompletableFuture.failedFuture(new PlcRuntimeException(
                "KNXnet/IP gateway does not support tunneling — refusing to connect."));
        }

        this.gatewayAddress = searchResponse.getDibDeviceInfo().getKnxAddress();
        this.gatewayName = new String(searchResponse.getDibDeviceInfo().getDeviceFriendlyName()).trim();
        LOGGER.info("Found KNXnet/IP Gateway '{}' with KNX address '{}.{}.{}'",
            gatewayName, gatewayAddress.getMainGroup(),
            gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup());

        pendingConnection = new CompletableFuture<>();
        ConnectionRequest connectionRequest = new ConnectionRequest(
            new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort),
            new HPAIDataEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort),
            new ConnectionRequestInformationTunnelConnection(tunnelConnectionType));

        return sendMessage(connectionRequest, pendingConnection)
            .thenAccept(this::applyConnectionResponse);
    }

    private void applyConnectionResponse(ConnectionResponse connectionResponse) {
        if (connectionResponse.getStatus() != Status.NO_ERROR) {
            throw new PlcRuntimeException(String.format(
                "Connection to KNXnet/IP gateway '%s' refused with status %s",
                gatewayName, connectionResponse.getStatus()));
        }
        this.communicationChannelId = connectionResponse.getCommunicationChannelId();
        ConnectionResponseDataBlockTunnelConnection tunnelBlock =
            (ConnectionResponseDataBlockTunnelConnection) connectionResponse.getConnectionResponseDataBlock();
        this.clientKnxAddress = tunnelBlock.getKnxAddress();
        LOGGER.info("Connected to KNXnet/IP Gateway '{}', assigned client KNX address '{}.{}.{}' (channel {})",
            gatewayName, clientKnxAddress.getMainGroup(), clientKnxAddress.getMiddleGroup(),
            clientKnxAddress.getSubGroup(), communicationChannelId);
    }

    private <T> CompletableFuture<T> sendMessage(KnxNetIpMessage message, CompletableFuture<T> waiting) {
        try {
            messageCodec.send(message);
        } catch (MessageCodecException e) {
            waiting.completeExceptionally(new PlcRuntimeException(
                "Failed to send " + message.getClass().getSimpleName(), e));
            return waiting;
        }
        return waiting.orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS);
    }

    private void startHeartbeat() {
        heartbeatTask = heartbeatExecutor.scheduleAtFixedRate(
            this::sendHeartbeat, HEARTBEAT_INTERVAL_MS, HEARTBEAT_INTERVAL_MS, TimeUnit.MILLISECONDS);
    }

    private void sendHeartbeat() {
        try {
            sendConnectionStateRequest()
                .orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
                .whenComplete((response, error) -> {
                    if (error != null) {
                        LOGGER.error("KNX heartbeat failed", error);
                        return;
                    }
                    if (response.getStatus() != Status.NO_ERROR) {
                        LOGGER.error("KNX heartbeat got non-OK status: {}", response.getStatus());
                    }
                });
        } catch (Exception e) {
            LOGGER.error("Failed to send KNX heartbeat", e);
        }
    }

    private CompletableFuture<ConnectionStateResponse> sendConnectionStateRequest() {
        CompletableFuture<ConnectionStateResponse> future = new CompletableFuture<>();
        pendingConnectionState.put(communicationChannelId, future);
        ConnectionStateRequest request = new ConnectionStateRequest(
            communicationChannelId,
            new HPAIControlEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
        try {
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingConnectionState.remove(communicationChannelId, future);
            future.completeExceptionally(new PlcRuntimeException("Failed to send ConnectionStateRequest", e));
        }
        return future;
    }

    @Override
    public void close() throws Exception {
        handshakeComplete = false;
        if (heartbeatTask != null) {
            heartbeatTask.cancel(false);
        }
        heartbeatExecutor.shutdownNow();

        // Try to be polite and tell the gateway we're going away. The gateway
        // doesn't always reply, so we just fire and forget.
        if (messageCodec != null && messageCodec.isOpen()) {
            try {
                DisconnectRequest disconnect = new DisconnectRequest(communicationChannelId,
                    new HPAIControlEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
                messageCodec.send(disconnect);
            } catch (Exception e) {
                LOGGER.debug("Failed to send DisconnectRequest (gateway may already be gone): {}", e.getMessage());
            }
        }

        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        failAllPending(new PlcRuntimeException("Connection closed"));
        consumers.clear();
        super.close();
        LOGGER.info("KNXnet/IP connection to '{}' closed", gatewayName);
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        handshakeComplete = false;
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");
        failAllPending(new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause));
    }

    private void failAllPending(Throwable cause) {
        if (pendingSearch != null && !pendingSearch.isDone()) {
            pendingSearch.completeExceptionally(cause);
        }
        if (pendingConnection != null && !pendingConnection.isDone()) {
            pendingConnection.completeExceptionally(cause);
        }
        pendingConnectionState.values().forEach(f -> f.completeExceptionally(cause));
        pendingConnectionState.clear();
        pendingTunneling.values().forEach(f -> f.completeExceptionally(cause));
        pendingTunneling.clear();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming messages
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingMessage(KnxNetIpMessage msg) {
        try {
            if (msg instanceof SearchResponse sr) {
                if (pendingSearch != null && !pendingSearch.isDone()) {
                    pendingSearch.complete(sr);
                }
            } else if (msg instanceof ConnectionResponse cr) {
                if (pendingConnection != null && !pendingConnection.isDone()) {
                    pendingConnection.complete(cr);
                }
            } else if (msg instanceof ConnectionStateResponse cs) {
                CompletableFuture<ConnectionStateResponse> f =
                    pendingConnectionState.remove(cs.getCommunicationChannelId());
                if (f != null) {
                    f.complete(cs);
                }
            } else if (msg instanceof TunnelingResponse tr) {
                short seq = tr.getTunnelingResponseDataBlock().getSequenceCounter();
                CompletableFuture<TunnelingResponse> f = pendingTunneling.remove(seq);
                if (f != null) {
                    f.complete(tr);
                }
            } else if (msg instanceof TunnelingRequest tunnelingRequest) {
                handleTunnelingRequest(tunnelingRequest);
            }
        } catch (Exception e) {
            LOGGER.error("Error handling incoming KNXnet/IP message {}", msg, e);
        }
    }

    private void handleTunnelingRequest(TunnelingRequest tunnelingRequest) throws Exception {
        short channelId = tunnelingRequest.getTunnelingRequestDataBlock().getCommunicationChannelId();
        if (channelId != communicationChannelId) {
            LOGGER.debug("Ignoring TunnelingRequest for channel {} (ours is {})", channelId, communicationChannelId);
            return;
        }

        if (tunnelingRequest.getCemi() instanceof LDataInd dataInd) {
            LDataFrame frame = dataInd.getDataFrame();
            if (frame instanceof LDataExtended ext) {
                Apdu apdu = ext.getApdu();
                if (apdu instanceof ApduDataContainer container) {
                    ApduData dataApdu = container.getDataApdu();
                    if (dataApdu instanceof ApduDataGroupValueWrite groupWrite) {
                        processCemiData(ext.getSourceAddress(), ext.getDestinationAddress(),
                            groupWrite.getDataFirstByte(), groupWrite.getData());
                    } else if (dataApdu instanceof ApduDataGroupValueResponse groupResponse) {
                        // A device answered our GroupValueRead. Complete the
                        // pending read future keyed by destination GA.
                        completePendingRead(ext.getDestinationAddress(),
                            groupResponse.getDataFirstByte(), groupResponse.getData());
                    }
                }
            }
        }

        // Confirm receipt regardless of whether we cared about the payload —
        // a missed ack will make the gateway resend the same indication.
        short sequenceCounter = tunnelingRequest.getTunnelingRequestDataBlock().getSequenceCounter();
        TunnelingResponse ack = new TunnelingResponse(
            new TunnelingResponseDataBlock(communicationChannelId, sequenceCounter, Status.NO_ERROR));
        messageCodec.send(ack);
    }

    private void processCemiData(KnxAddress sourceAddress, byte[] destinationGroupAddress,
                                 byte firstByte, byte[] restBytes) throws BufferException {
        byte[] payload = new byte[1 + restBytes.length];
        payload[0] = firstByte;
        System.arraycopy(restBytes, 0, payload, 1, restBytes.length);

        ReadBufferByteBased addressBuffer = new ReadBufferByteBased(destinationGroupAddress,
            KnxNetIpMessageCodec.BUFFER_OPTIONS);
        KnxGroupAddress knxGroupAddress = KnxGroupAddress.staticParse(addressBuffer, groupAddressType);
        String destinationAddress = toGroupAddressString(knxGroupAddress);

        if (etsModel == null) {
            LOGGER.info("Raw Message: '{}' to: '{}' payload: '{}'",
                toString(sourceAddress), destinationAddress, Hex.encodeHexString(payload));
            return;
        }

        GroupAddress groupAddress = etsModel.getGroupAddresses().get(destinationAddress);
        if (groupAddress == null || groupAddress.getType() == null) {
            LOGGER.warn("Message from '{}' to unknown group address '{}' payload '{}'",
                toString(sourceAddress), destinationAddress, Hex.encodeHexString(payload));
            return;
        }

        ReadBufferByteBased rawDataReader = new ReadBufferByteBased(payload,
            KnxNetIpMessageCodec.BUFFER_OPTIONS);
        PlcValue value = KnxDatapoint.staticParse(rawDataReader, groupAddress.getType());

        String areaName = etsModel.getTopologyName(
            destinationAddress.substring(0, destinationAddress.indexOf('/')));
        String lineName = etsModel.getTopologyName(
            destinationAddress.substring(0, destinationAddress.indexOf('/', destinationAddress.indexOf('/') + 1)));

        // Attach ETS metadata directly to the value so consumers can read it
        // via PlcValue#getMetaData(...).
        if (value instanceof PlcValueAdapter adapter) {
            adapter.addMetaData("groupAddress", new PlcSTRING(groupAddress.getGroupAddress()));
            adapter.addMetaData("sourceAddress", new PlcSTRING(toString(sourceAddress)));
            if (groupAddress.getName() != null) {
                adapter.addMetaData("name", new PlcSTRING(groupAddress.getName()));
            }
            String description = buildDescription(groupAddress, areaName, lineName);
            if (description != null) {
                adapter.addMetaData("description", new PlcSTRING(description));
            }
            if (groupAddress.getFunction() != null) {
                adapter.addMetaData("location", new PlcSTRING(groupAddress.getFunction().getSpaceName()));
                adapter.addMetaData("function", new PlcSTRING(groupAddress.getFunction().getName()));
            }
            if (areaName != null) {
                adapter.addMetaData("area", new PlcSTRING(areaName));
            }
            if (lineName != null) {
                adapter.addMetaData("line", new PlcSTRING(lineName));
            }
            adapter.addMetaData("dataPointType", new PlcSTRING(groupAddress.getType().getName()));
        }

        publishEvent(groupAddress, value);
    }

    /**
     * Produces a "location function" blurb for the {@code description} metadata
     * field — e.g. {@code "Flur Energieversorgung"}. Falls back to the raw GA
     * description text if no function/space info is available.
     */
    private static String buildDescription(GroupAddress groupAddress, String areaName, String lineName) {
        if (groupAddress.getFunction() != null) {
            String space = groupAddress.getFunction().getSpaceName();
            String func = groupAddress.getFunction().getName();
            if (space != null && func != null) {
                return space + " " + func;
            }
            if (func != null) {
                return func;
            }
        }
        if (lineName != null) {
            return lineName;
        }
        return areaName;
    }

    /**
     * Decodes the destination GA from a {@code GroupValueResponse} indication
     * back to its textual form and hands the raw payload off to whichever
     * read future is waiting on it. If no one is waiting, the indication is
     * dropped — it's almost certainly an answer to someone else on the bus.
     */
    private void completePendingRead(byte[] destinationAddress, byte firstByte, byte[] restBytes) {
        String dest;
        try {
            ReadBufferByteBased addressBuffer = new ReadBufferByteBased(destinationAddress,
                KnxNetIpMessageCodec.BUFFER_OPTIONS);
            KnxGroupAddress addr = KnxGroupAddress.staticParse(addressBuffer, groupAddressType);
            dest = toGroupAddressString(addr);
        } catch (BufferException e) {
            LOGGER.debug("Failed to decode destination address on GroupValueResponse", e);
            return;
        }
        CompletableFuture<byte[]> future = pendingReads.remove(dest);
        if (future == null) {
            return;
        }
        byte[] payload = new byte[1 + restBytes.length];
        payload[0] = firstByte;
        System.arraycopy(restBytes, 0, payload, 1, restBytes.length);
        future.complete(payload);
    }

    private void publishEvent(GroupAddress groupAddress, PlcValue plcValue) {
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            DefaultPlcConsumerRegistration registration = entry.getKey();
            Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
                if (handle instanceof KnxNetIpSubscriptionHandle sub && sub.matches(groupAddress)) {
                    PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
                        Collections.singletonMap(sub.getName(),
                            new DefaultPlcResponseItem<>(PlcResponseCode.OK, plcValue)));
                    consumer.accept(event);
                }
            }
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Ping
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        CompletableFuture<PlcPingResponse> future = new CompletableFuture<>();
        sendConnectionStateRequest()
            .orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
            .whenComplete((response, error) -> {
                if (error != null) {
                    future.complete(new DefaultPlcPingResponse(pingRequest, PlcResponseCode.REMOTE_ERROR));
                    return;
                }
                future.complete(new DefaultPlcPingResponse(pingRequest,
                    response.getStatus() == Status.NO_ERROR
                        ? PlcResponseCode.OK
                        : PlcResponseCode.REMOTE_ERROR));
            });
        return future;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        if (writeRequest.getTagNames().isEmpty()) {
            future.complete(new DefaultPlcWriteResponse(writeRequest, Collections.emptyMap()));
            return future;
        }

        // KNX driver is single-tag-per-request — same shape as the legacy
        // implementation, driven by the SingleTagOptimizer wiring that is
        // applied by the platform when this driver was registered.
        String tagName = writeRequest.getTagNames().iterator().next();
        PlcTag tag = writeRequest.getTag(tagName);
        if (tag == null) {
            // Tag parsing failed in the request builder — surface that as an
            // INVALID_ADDRESS response code rather than NPE'ing on tag.getClass().
            future.complete(new DefaultPlcWriteResponse(writeRequest,
                Collections.singletonMap(tagName, PlcResponseCode.INVALID_ADDRESS)));
            return future;
        }
        if (!(tag instanceof KnxNetIpTag knxTag)) {
            future.completeExceptionally(new PlcRuntimeException(
                "Unsupported tag type " + tag.getClass().getSimpleName()));
            return future;
        }

        byte[] destinationAddress;
        try {
            destinationAddress = toKnxAddressData(knxTag);
        } catch (Exception e) {
            future.completeExceptionally(e);
            return future;
        }

        if (sequenceCounter.compareAndSet(Short.MAX_VALUE, 0)) {
            // Wrapped; nothing else to do.
        }
        short seq = (short) sequenceCounter.getAndIncrement();

        PlcValue value = writeRequest.getPlcValue(tagName);
        byte dataFirstByte;
        byte[] data;
        try {
            byte[] payload = serializePayload(value, destinationAddress);
            dataFirstByte = payload[0];
            data = new byte[payload.length - 1];
            System.arraycopy(payload, 1, data, 0, payload.length - 1);
        } catch (Exception e) {
            future.completeExceptionally(e);
            return future;
        }

        // Frame layout, which is known to work on real KNX gateways.
        // The legacy plc4x driver had frameType=false and numbered=true
        // here — those values produce a frame the gateway silently drops
        // without raising an error.
        TunnelingRequest knxRequest = new TunnelingRequest(
            new TunnelingRequestDataBlock(communicationChannelId, seq),
            new LDataReq((short) 0, new ArrayList<>(0),
                new LDataExtended(true, false, CEMIPriority.LOW, false, false,
                    true, (byte) 6, (byte) 0, clientKnxAddress, destinationAddress,
                    new ApduDataContainer(false, (byte) 0,
                        new ApduDataGroupValueWrite(dataFirstByte, data)))));

        CompletableFuture<TunnelingResponse> ackFuture = new CompletableFuture<>();
        pendingTunneling.put(seq, ackFuture);
        try {
            messageCodec.send(knxRequest);
        } catch (MessageCodecException e) {
            pendingTunneling.remove(seq);
            future.completeExceptionally(new PlcRuntimeException("Failed to send TunnelingRequest", e));
            return future;
        }

        ackFuture.orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
            .whenComplete((response, error) -> {
                if (error instanceof TimeoutException) {
                    pendingTunneling.remove(seq);
                }
                if (error != null) {
                    future.completeExceptionally(error);
                    return;
                }
                PlcResponseCode code = response.getTunnelingResponseDataBlock().getStatus() == Status.NO_ERROR
                    ? PlcResponseCode.OK
                    : PlcResponseCode.INTERNAL_ERROR;
                future.complete(new DefaultPlcWriteResponse(writeRequest,
                    Collections.singletonMap(tagName, code)));
            });
        return future;
    }

    private byte[] serializePayload(PlcValue value, byte[] destinationAddress) throws BufferException {
        if (etsModel != null) {
            String destinationAddressString = etsModel.parseGroupAddress(destinationAddress);
            GroupAddress groupAddress = etsModel.getGroupAddresses().get(destinationAddressString);
            if (groupAddress == null || groupAddress.getType() == null) {
                throw new PlcRuntimeException(
                    "ETS model didn't specify group address '" + destinationAddressString +
                        "' or didn't define a type for it.");
            }
            int length = KnxDatapoint.getLengthInBytes(value, groupAddress.getType());
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[length],
                KnxNetIpMessageCodec.BUFFER_OPTIONS);
            KnxDatapoint.staticSerialize(writeBuffer, value, groupAddress.getType());
            return writeBuffer.getBytes();
        }

        // No ETS model: only single-byte or list-of-bytes writes work, with the
        // first byte capped at 6 bits (raw KNX group write semantics).
        if (value.isByte()) {
            checkSixBitByte(value.getByte());
            return new byte[]{value.getByte()};
        }
        if (value.isList()) {
            List<? extends PlcValue> list = value.getList();
            if (list.isEmpty()) {
                throw new PlcRuntimeException(
                    "If no ETS model is provided, a write payload must be at least one byte.");
            }
            byte[] out = new byte[list.size()];
            for (int i = 0; i < list.size(); i++) {
                PlcValue element = list.get(i);
                if (!element.isByte()) {
                    throw new PlcRuntimeException(
                        "If no ETS model is provided, all write payload elements must be bytes.");
                }
                if (i == 0) {
                    checkSixBitByte(element.getByte());
                }
                out[i] = element.getByte();
            }
            return out;
        }
        throw new PlcRuntimeException(
            "If no ETS model is provided, the only supported types for writing are single byte or list of bytes.");
    }

    private static void checkSixBitByte(byte b) {
        if (b > 63 || b < 0) {
            throw new PlcRuntimeException(
                "If no ETS model is provided, value of the first byte must be between 0 and 63.");
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * KNX reads are asynchronous: send a {@code GroupValueRead} to a group
     * address and wait for any device on the bus to answer with a matching
     * {@code GroupValueResponse} indication. Not every GA has a device that
     * responds — some only emit unsolicited writes on change. Reads against
     * those will time out.
     *
     * <p>Tags must carry their own {@code :DPT…} suffix or be defined in the
     * loaded ETS project; otherwise we can't decode the response payload.</p>
     */
    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        return readSequentially(readRequest,
            new ArrayList<>(readRequest.getTagNames()), new LinkedHashMap<>());
    }

    private CompletableFuture<PlcReadResponse> readSequentially(PlcReadRequest readRequest,
                                                                List<String> remaining,
                                                                Map<String, PlcResponseItem<PlcValue>> results) {
        if (remaining.isEmpty()) {
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse(readRequest, results));
        }
        String tagName = remaining.get(0);
        List<String> next = remaining.subList(1, remaining.size());
        PlcTag tag = readRequest.getTag(tagName);
        if (!(tag instanceof KnxNetIpTag knxTag)) {
            results.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            return readSequentially(readRequest, next, results);
        }
        return readSingleTag(knxTag)
            .handle((value, error) -> {
                if (error != null) {
                    LOGGER.warn("Read of {} failed", tagName, error);
                    PlcResponseCode code = error instanceof TimeoutException
                        ? PlcResponseCode.NOT_FOUND   // no device on the bus answered our read
                        : PlcResponseCode.INTERNAL_ERROR;
                    results.put(tagName, new DefaultPlcResponseItem<>(code, null));
                } else {
                    results.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, value));
                }
                return null;
            })
            .thenCompose(v -> readSequentially(readRequest, next, results));
    }

    private CompletableFuture<PlcValue> readSingleTag(KnxNetIpTag tag) {
        byte[] destinationAddress;
        String dest;
        try {
            destinationAddress = toKnxAddressData(tag);
            ReadBufferByteBased addressBuffer = new ReadBufferByteBased(destinationAddress,
                KnxNetIpMessageCodec.BUFFER_OPTIONS);
            KnxGroupAddress addr = KnxGroupAddress.staticParse(addressBuffer, groupAddressType);
            dest = toGroupAddressString(addr);
        } catch (Exception e) {
            return CompletableFuture.failedFuture(e);
        }
        KnxDatapointType dpt = resolveDatapointType(tag, dest);
        if (dpt == null) {
            return CompletableFuture.failedFuture(new PlcRuntimeException(
                "Cannot decode read response for '" + tag.getAddressString() +
                    "': supply a :DPT… suffix on the tag or load a knxproj file."));
        }

        CompletableFuture<byte[]> responseFuture = new CompletableFuture<>();
        pendingReads.put(dest, responseFuture);

        short seq = (short) sequenceCounter.getAndIncrement();
        TunnelingRequest readRequestMessage = new TunnelingRequest(
            new TunnelingRequestDataBlock(communicationChannelId, seq),
            new LDataReq((short) 0, new ArrayList<>(0),
                new LDataExtended(true, false, CEMIPriority.LOW, false, false,
                    true, (byte) 6, (byte) 0, clientKnxAddress, destinationAddress,
                    new ApduDataContainer(false, (byte) 0, new ApduDataGroupValueRead()))));

        CompletableFuture<TunnelingResponse> ackFuture = new CompletableFuture<>();
        pendingTunneling.put(seq, ackFuture);
        try {
            messageCodec.send(readRequestMessage);
        } catch (MessageCodecException e) {
            pendingTunneling.remove(seq);
            pendingReads.remove(dest);
            return CompletableFuture.failedFuture(new PlcRuntimeException("Failed to send GroupValueRead", e));
        }

        // Wait for both the tunneling ack and the actual GroupValueResponse
        // indication. Honor the configured request timeout on the wait for
        // the indication — the gateway will ack quickly, but the bus device
        // that holds the value may not answer at all.
        return ackFuture
            .orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
            .thenCompose(ack -> {
                if (ack.getTunnelingResponseDataBlock().getStatus() != Status.NO_ERROR) {
                    pendingReads.remove(dest);
                    return CompletableFuture.failedFuture(new PlcRuntimeException(
                        "Tunneling ack for read failed: " +
                            ack.getTunnelingResponseDataBlock().getStatus()));
                }
                return responseFuture
                    .orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
                    .whenComplete((d, err) -> {
                        if (err != null) {
                            pendingReads.remove(dest);
                        }
                    });
            })
            .thenApply(payload -> {
                try {
                    ReadBufferByteBased buf = new ReadBufferByteBased(payload,
                        KnxNetIpMessageCodec.BUFFER_OPTIONS);
                    return KnxDatapoint.staticParse(buf, dpt);
                } catch (BufferException e) {
                    throw new PlcRuntimeException("Failed to decode read response", e);
                }
            });
    }

    /**
     * Resolves the {@link KnxDatapointType} for a tag, preferring an explicit
     * {@code :DPT…} suffix on the tag and falling back to the ETS model.
     * Returns {@code null} if neither source has a usable type — the caller
     * should surface that as an error since the response cannot be decoded.
     */
    private KnxDatapointType resolveDatapointType(KnxNetIpTag tag, String destinationAddress) {
        if (tag.getDptId() != null) {
            KnxDatapointType byTag = resolveDatapointTypeFromHint(tag.getDptId());
            if (byTag != null) {
                return byTag;
            }
        }
        if (etsModel == null) {
            return null;
        }
        GroupAddress ga = etsModel.getGroupAddresses().get(destinationAddress);
        return ga != null ? ga.getType() : null;
    }

    /**
     * Maps a textual {@code DPT<n>[.<sub>]} hint to the generated
     * {@link KnxDatapointType} constant. With a sub-number, we look for an
     * exact match; without one, we fall back to the first DPT in the matching
     * main type — enough for the simple cases (DPT1, DPT5, DPT9, …).
     */
    static KnxDatapointType resolveDatapointTypeFromHint(String dptId) {
        if (dptId == null) {
            return null;
        }
        String body = dptId.toUpperCase().startsWith("DPT") ? dptId.substring(3) : dptId;
        int dot = body.indexOf('.');
        int mainNum;
        Integer subNum = null;
        try {
            if (dot >= 0) {
                mainNum = Integer.parseInt(body.substring(0, dot));
                subNum = Integer.parseInt(body.substring(dot + 1));
            } else {
                mainNum = Integer.parseInt(body);
            }
        } catch (NumberFormatException e) {
            return null;
        }
        KnxDatapointMainType mainType = KnxDatapointMainType.firstEnumForFieldNumber(mainNum);
        if (mainType == null) {
            return null;
        }
        for (KnxDatapointType type : KnxDatapointType.values()) {
            if (type.getDatapointMainType() != mainType) {
                continue;
            }
            if (subNum == null || type.getNumber() == subNum) {
                return type;
            }
        }
        return null;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Browse — backed by the loaded ETS project model
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * KNX is event-driven on the bus, so every browseable GA is delivered as
     * a push subscription only. We never advertise cyclic / change-of-state.
     */
    private static final Set<PlcSubscriptionType> KNX_SUBSCRIPTION_TYPES =
        Collections.unmodifiableSet(EnumSet.of(PlcSubscriptionType.EVENT));

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        Map<String, List<PlcBrowseItem>> collected = new HashMap<>();
        PlcBrowseRequestInterceptor collector = (queryName, query, item) -> {
            collected.computeIfAbsent(queryName, k -> new ArrayList<>()).add(item);
            return true;
        };
        return onBrowseWithInterceptor(browseRequest, collector)
            .thenApply(emptyResp -> {
                Map<String, PlcResponseCode> codes = new HashMap<>();
                Map<String, List<PlcBrowseItem>> values = new TreeMap<>();
                for (String name : browseRequest.getQueryNames()) {
                    codes.put(name, emptyResp.getResponseCode(name));
                    List<PlcBrowseItem> items = collected.getOrDefault(name, Collections.emptyList());
                    items.sort(Comparator.comparing(i -> i.getTag().getAddressString()));
                    values.put(name, items);
                }
                return new DefaultPlcBrowseResponse(browseRequest, codes, values);
            });
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowseWithInterceptor(PlcBrowseRequest browseRequest,
                                                                          PlcBrowseRequestInterceptor interceptor) {
        if (etsModel == null) {
            return CompletableFuture.failedFuture(new PlcRuntimeException(
                "Browse is not available — no knxproj-file-path configured."));
        }
        Map<String, PlcResponseCode> codes = new HashMap<>();
        Map<String, List<PlcBrowseItem>> emptyValues = new TreeMap<>();
        for (String queryName : browseRequest.getQueryNames()) {
            PlcQuery query = browseRequest.getQuery(queryName);
            String queryString = query != null ? query.getQueryString() : "*";
            Pattern pattern = compileQueryPattern(queryString);
            for (GroupAddress info : etsModel.getGroupAddresses().values()) {
                String addressString = info.getGroupAddress();
                if (pattern != null && !pattern.matcher(addressString).matches()) {
                    continue;
                }
                interceptor.intercept(queryName, query, createBrowseItem(info));
            }
            codes.put(queryName, PlcResponseCode.OK);
            emptyValues.put(queryName, Collections.emptyList());
        }
        return CompletableFuture.completedFuture(
            new DefaultPlcBrowseResponse(browseRequest, codes, emptyValues));
    }

    /**
     * Compiles a browse query string ({@code *}, {@code 1/*}, {@code 1/2/*},
     * {@code 1/2/3}) into a regex. Returns {@code null} for "match all" so the
     * caller can skip pattern matching entirely on the hot path.
     */
    private static Pattern compileQueryPattern(String queryString) {
        if (queryString == null || queryString.isEmpty()
            || "*".equals(queryString) || "**".equals(queryString) || "*/*/*".equals(queryString)) {
            return null;
        }
        String[] parts = queryString.split("/");
        String regex;
        if (parts.length == 2 && "*".equals(parts[1])) {
            regex = parts[0] + "/\\d+/\\d+";
        } else if (parts.length == 3 && "*".equals(parts[2])) {
            regex = parts[0] + "/" + parts[1] + "/\\d+";
        } else {
            regex = queryString.replace("*", "\\d+");
        }
        return Pattern.compile("^" + regex + "$");
    }

    private PlcBrowseItem createBrowseItem(GroupAddress info) {
        KnxDatapointType dpt = info.getType();
        String tagAddress = info.getGroupAddress();
        if (dpt != null && dpt != KnxDatapointType.DPT_UNKNOWN && dpt.getDatapointMainType() != null) {
            tagAddress += ":DPT" + dpt.getDatapointMainType().getNumber();
        }
        KnxNetIpTag tag = KnxNetIpTag.of(tagAddress);

        Map<String, PlcValue> options = new LinkedHashMap<>();
        if (dpt != null && dpt != KnxDatapointType.DPT_UNKNOWN) {
            options.put("dataPointType", new PlcSTRING(dpt.name()));
            if (dpt.getDatapointMainType() != null) {
                options.put("dataPointMainType", new PlcSTRING(dpt.getDatapointMainType().name()));
            }
            options.put("unitOfMeasurement", new PlcSTRING(dpt.getName()));
        }
        if (info.getFunction() != null) {
            options.put("function", new PlcSTRING(info.getFunction().getName()));
            options.put("location", new PlcSTRING(info.getFunction().getSpaceName()));
        }

        return new DefaultPlcBrowseItem(
            tag,
            info.getName() != null ? info.getName() : info.getGroupAddress(),
            true,                       // readable
            true,                       // writable
            KNX_SUBSCRIPTION_TYPES,
            false,                      // not publishable
            null,                       // no array info
            null,                       // no children
            options);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe / consumer registry
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();
        for (String tagName : subscriptionRequest.getTagNames()) {
            DefaultPlcSubscriptionTag subscriptionTag =
                (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(tagName);
            if (!(subscriptionTag.getTag() instanceof KnxNetIpTag knxTag)) {
                values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            } else {
                values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK,
                    new KnxNetIpSubscriptionHandle(this, tagName, knxTag)));
            }
        }

        // If the request was built with setConsumer(...) or addEventTagAddress(name, addr, consumer),
        // auto-register that consumer here. The SPI plumbs the consumer through the request but
        // doesn't fire it anywhere — that wiring is per-driver.
        if (subscriptionRequest instanceof DefaultPlcSubscriptionRequest req) {
            Consumer<PlcSubscriptionEvent> requestConsumer = req.getConsumer();
            for (Map.Entry<String, PlcResponseItem<PlcSubscriptionHandle>> entry : values.entrySet()) {
                if (entry.getValue().getResponseCode() != PlcResponseCode.OK) {
                    continue;
                }
                PlcSubscriptionHandle handle = entry.getValue().getValue();
                Consumer<PlcSubscriptionEvent> perTag = req.getTagConsumer(entry.getKey());
                if (perTag != null) {
                    handle.register(perTag);
                }
                if (requestConsumer != null) {
                    handle.register(requestConsumer);
                }
            }
        }

        return CompletableFuture.completedFuture(
            new DefaultPlcSubscriptionResponse(subscriptionRequest, values));
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        // KNX has no per-tag wire unsubscribe — the gateway streams every
        // LDataInd whether we want it or not, we just stop matching consumers.
        return CompletableFuture.completedFuture(new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest));
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer,
                                                        Collection<PlcSubscriptionHandle> handles) {
        DefaultPlcConsumerRegistration registration =
            new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(registration, consumer);
        return registration;
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        if (registration instanceof DefaultPlcConsumerRegistration r) {
            consumers.remove(r);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private byte[] toKnxAddressData(KnxNetIpTag tag) {
        WriteBufferByteBased address = new WriteBufferByteBased(new byte[2], KnxNetIpMessageCodec.BUFFER_OPTIONS);
        try {
            switch (groupAddressType) {
                case 3:
                    address.writeUnsignedShort(5, Short.parseShort(tag.getMainGroup()));
                    address.writeUnsignedByte(3, Byte.parseByte(tag.getMiddleGroup()));
                    address.writeUnsignedShort(8, Short.parseShort(tag.getSubGroup()));
                    break;
                case 2:
                    address.writeUnsignedShort(5, Short.parseShort(tag.getMainGroup()));
                    address.writeUnsignedShort(11, Short.parseShort(tag.getSubGroup()));
                    break;
                case 1:
                    address.writeUnsignedShort(16, Short.parseShort(tag.getSubGroup()));
                    break;
                default:
                    throw new PlcRuntimeException("Unsupported group address type " + groupAddressType);
            }
        } catch (Exception e) {
            throw new PlcRuntimeException("Error converting tag into knx address data.", e);
        }
        return address.getBytes();
    }

    private static String toString(KnxAddress knxAddress) {
        return knxAddress.getMainGroup() + "." + knxAddress.getMiddleGroup() + "." + knxAddress.getSubGroup();
    }

    private static String toGroupAddressString(KnxGroupAddress groupAddress) {
        if (groupAddress instanceof KnxGroupAddress3Level level3) {
            return level3.getMainGroup() + "/" + level3.getMiddleGroup() + "/" + level3.getSubGroup();
        }
        if (groupAddress instanceof KnxGroupAddress2Level level2) {
            return level2.getMainGroup() + "/" + level2.getSubGroup();
        }
        if (groupAddress instanceof KnxGroupAddressFreeLevel free) {
            return Integer.toString(free.getSubGroup());
        }
        throw new PlcRuntimeException("Unsupported Group Address Type " + groupAddress.getClass().getName());
    }

}
