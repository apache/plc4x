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
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.eip.base.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.base.tag.EipTagHandler;
import org.apache.plc4x.java.eip.readwrite.*;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.utils.subscriptionemulation.PollingSubscriptionConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.*;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.concurrent.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * EtherNet/IP (CIP encapsulation) TCP connection — direct port of the legacy
 * EipProtocolLogic to the SPI3 ConnectionBase model. Request/response correlation
 * uses the 8-byte senderContext echoed back by the device.
 */
public class EipTcpConnection extends PollingSubscriptionConnectionBase<EIPConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(EipTcpConnection.class);

    protected static final byte[] DEFAULT_SENDER_CONTEXT = "PLC4X   ".getBytes(StandardCharsets.US_ASCII);
    protected static final long EMPTY_SESSION_HANDLE = 0L;
    protected static final long EMPTY_INTERFACE_OPTIONS = 0L;
    protected static final long EMPTY_INTERFACE_HANDLE = 0L;

    protected final boolean bigEndian;

    private EipTcpMessageCodec messageCodec;
    // EIP has no transaction-id field. The original protocol logic relied on a
    // RequestTransactionManager(1) — at most one request in flight, responses
    // arrive in send order. With getMaxConcurrentRequests() == 1 we can use the
    // same FIFO scheme: enqueue the response future on send, complete the head
    // on receive. (Earlier this class encoded a counter into the 8-byte
    // senderContext for correlation, but that mutated the shared static
    // DEFAULT_SENDER_CONTEXT and made every outgoing packet's senderContext
    // unpredictable — which breaks scripted DriverTestsuite tests that expect
    // the literal "PLC4X   " bytes.)
    private final BlockingQueue<CompletableFuture<EipPacket>> pendingRequests = new LinkedBlockingDeque<>();

    protected long sessionHandle = EMPTY_SESSION_HANDLE;
    protected long connectionId = 0L;
    protected int sequenceCount = 1;
    protected boolean useConnectionManager = false;
    protected boolean useMessageRouter = false;
    protected boolean cipEncapsulationAvailable = false;

    private final NullAddressItem nullAddressItem = new NullAddressItem();
    private final List<PathSegment> routingAddress = new ArrayList<>();
    private short connectionPathSize = 0;
    private final int connectionSerialNumber = ThreadLocalRandom.current().nextInt(1, 0x10000);

    public EipTcpConnection(EIPConfiguration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        this(configuration, transportInstance, auditLog, configuration.isBigEndian());
    }

    public EipTcpConnection(EIPConfiguration configuration, TransportInstance<?> transportInstance, AuditLog auditLog, boolean bigEndian) {
        super(configuration, transportInstance, auditLog);
        this.bigEndian = bigEndian;
        initRoutingAddress();
    }

    private void initRoutingAddress() {
        if (configuration.getCommunicationPath() != null) {
            String[] splitConnectionPath = configuration.getCommunicationPath().split(",");
            if (splitConnectionPath.length % 2 == 0) {
                for (int i = 0; (i + 1) < splitConnectionPath.length; i += 2) {
                    switch (splitConnectionPath[i]) {
                        case "1":
                            int backplanePortId = Integer.parseInt(splitConnectionPath[i]);
                            int slot = Integer.parseInt(splitConnectionPath[i + 1]);
                            routingAddress.add(new PortSegment(new PortSegmentNormal((byte) backplanePortId, (short) slot)));
                            break;
                        case "2":
                            int ethernetPortId = Integer.parseInt(splitConnectionPath[i]);
                            String ipAddress = splitConnectionPath[i + 1];
                            int lengthString = ipAddress.length();
                            if ((ipAddress.length() % 2) != 0) {
                                ipAddress += "\0";
                            }
                            routingAddress.add(new PortSegment(new PortSegmentExtended((byte) ethernetPortId, (short) lengthString, ipAddress)));
                            break;
                        default:
                            LOGGER.error("Only backplane or Ethernet module routing is supported");
                    }
                }
            }
        } else {
            routingAddress.add(new PortSegment(new PortSegmentNormal((byte) 1, (short) configuration.getSlot())));
        }
        routingAddress.add(new LogicalSegment(new ClassID((byte) 0, (short) 2)));
        routingAddress.add(new LogicalSegment(new InstanceID((byte) 0, (short) 1)));

        int totalBytes = 0;
        for (PathSegment segment : routingAddress) {
            totalBytes += segment.getLengthInBytes();
        }
        if ((totalBytes % 2) != 0) {
            totalBytes += 1;
        }
        this.connectionPathSize = (short) (totalBytes / 2);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new EipTagHandler();
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
        return messageCodec != null && messageCodec.isOpen();
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new EipTcpMessageCodec(transportInstance, this::handleIncomingMessage, bigEndian);

        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming EIP data", e);
            }
        });

        try {
            doConnectHandshake().get(getConfiguration().getRequestTimeout() * 5L, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("Error during EIP connect handshake", e);
        }

        LOGGER.info("EIP TCP connection established");
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT, "EIP TCP connection established");
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    private CompletableFuture<Void> doConnectHandshake() {
        // 1) ListServices to confirm CIP encapsulation
        ListServicesRequest listServicesRequest = new ListServicesRequest(
            EMPTY_SESSION_HANDLE, CIPStatus.Success.getValue(), DEFAULT_SENDER_CONTEXT, 0L);
        return sendRequest(listServicesRequest).thenCompose(listResponse -> {
            if (listResponse.getStatus() == CIPStatus.Success.getValue() && listResponse instanceof ListServicesResponse lsr) {
                if (!lsr.getTypeIds().isEmpty() && lsr.getTypeIds().getFirst() instanceof ServicesResponse sr) {
                    this.cipEncapsulationAvailable = sr.getSupportsCIPEncapsulation();
                }
            }
            // 2) RegisterSession
            EipConnectionRequest connectionRequest = new EipConnectionRequest(
                EMPTY_SESSION_HANDLE, CIPStatus.Success.getValue(), DEFAULT_SENDER_CONTEXT, 0L);
            return sendRequest(connectionRequest);
        }).thenCompose(sessionResponse -> {
            if (!(sessionResponse instanceof EipConnectionResponse)) {
                // Some devices skip ahead to connection-manager — proceed without handshake
                return CompletableFuture.completedFuture(null);
            }
            if (sessionResponse.getStatus() != CIPStatus.Success.getValue()) {
                return CompletableFuture.failedFuture(new PlcRuntimeException(
                    "Got status code while registering session [" + sessionResponse.getStatus() + "]"));
            }
            this.sessionHandle = sessionResponse.getSessionHandle();
            LOGGER.debug("Got assigned with Session handle {}", sessionHandle);

            // 3) Probe connection-manager and message-router capabilities.
            //    Skipped when the user has already committed to the unconnected
            //    path — the probe is only useful for choosing between the
            //    connection-manager and message-router code paths, and some
            //    simulators (e.g. cpppo) return a payload that doesn't match
            //    the Logix layout, which would otherwise abort the handshake.
            if (getConfiguration().isForceUnconnectedOperation()) {
                return CompletableFuture.completedFuture(null);
            }
            return probeClassObjectSupport().thenCompose(v -> {
                if (useConnectionManager) {
                    return openConnectionManager();
                }
                return CompletableFuture.completedFuture(null);
            });
        });
    }

    // Using this method to continue having the GetAttributeAll that can be useful in the future
    private CompletableFuture<Void> checkClassObjectAttributes(CIPClassID classId) {
        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) classId.getValue()));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));
        UnConnectedDataItem exchange = new UnConnectedDataItem(
            new GetAttributeAllRequest(classSegment, instanceSegment));
        List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);
        CipRRData eipWrapper = new CipRRData(sessionHandle, CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT, 0L, EMPTY_INTERFACE_HANDLE, 0, typeIds);

        return sendRequest(eipWrapper).thenAccept(response -> {

            if (extractCipService(response) instanceof GetAttributeAllResponse gar &&
                gar.getStatus() == CIPStatus.Success.getValue() &&
                gar.getAttributes() != null) {
                LOGGER.debug("Identity getNumberActive {}", gar.getAttributes().getNumberActive());
            }
        });
    }

    private CipService extractCipService(EipPacket response) {
        if (response instanceof CipRRData rr
            && rr.getStatus() == CIPStatus.Success.getValue()
            && rr.getTypeIds().size() > 1
            && rr.getTypeIds().get(1) instanceof UnConnectedDataItem di) {
            return di.getService();
        }
        return null;
    }

    private CompletableFuture<Void> probeClassObjectSupport() {

        return checkClassObjectSupport(CIPClassID.ConnectionManager).thenCompose(hasSupport -> {
            useConnectionManager = hasSupport;
            return checkClassObjectSupport(CIPClassID.MessageRouter);
        }).thenCompose(hasSupport -> {
            useMessageRouter = hasSupport;
            return checkClassObjectAttributes(CIPClassID.Identity);
        }).exceptionally(e -> {
            // Treat any probe failure (timeout, parse error, malformed response,
            // ServiceNotSupported) as the state that was achieved. This keeps
            // the handshake working against non-Logix devices and simulators.
            LOGGER.debug("probeClassObjectSupport probe failed, keeping the state achieved", e);
            return null;
        });
    }

    private CompletableFuture<Boolean> checkClassObjectSupport(CIPClassID classId) {

        UnConnectedDataItem exchange = new UnConnectedDataItem(new GetAttributeSingleRequest(
            new LogicalSegment(new ClassID((byte) 0, (short) classId.getValue())),
            new LogicalSegment(new InstanceID((byte) 0, (short) 0)), // Class level discovery
            new LogicalSegment(new AttributeID((byte) 0, (short) 1))) // Attribute ID 1: Revision
        );
        List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);
        CipRRData eipWrapper = new CipRRData(sessionHandle, CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT, 0L, EMPTY_INTERFACE_HANDLE, 0, typeIds);

        return sendRequest(eipWrapper).thenCompose(response -> {

            if (extractCipService(response) instanceof GetAttributeSingleResponse gsr) {
                boolean hasSupport = gsr.getStatus() == CIPStatus.Success.getValue();
                LOGGER.debug("ClassId: {} status: {} hasSupport: {}", classId, gsr.getStatus(), hasSupport);
                return CompletableFuture.completedFuture(hasSupport);
            }
            return CompletableFuture.completedFuture(false);
        });
    }

    private CompletableFuture<Void> openConnectionManager() {
        UnConnectedDataItem exchange = new UnConnectedDataItem(
            new CipConnectionManagerRequest(
                new LogicalSegment(new ClassID((byte) 0, (short) CIPClassID.ConnectionManager.getValue())),
                new LogicalSegment(new InstanceID((byte) 0, (short) 1)),
                (byte) 0, (byte) 10, (short) 14, 536870914L, 33944L,
                this.connectionSerialNumber, 4919, 42L, (short) 3, 2101812L,
                new NetworkConnectionParameters(4002, false, (byte) 2, (byte) 0, true),
                2113537L,
                new NetworkConnectionParameters(4002, false, (byte) 2, (byte) 0, true),
                new TransportType(true, (byte) 2, (byte) 3),
                this.connectionPathSize, this.routingAddress));
        List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);
        CipRRData eipWrapper = new CipRRData(sessionHandle, CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT, 0L, EMPTY_INTERFACE_HANDLE, 0, typeIds);

        return sendRequest(eipWrapper).thenAccept(response -> {
            if (!(response instanceof CipRRData rr) || rr.getStatus() != 0L) {
                throw new PlcRuntimeException("Got status code while opening Connection Manager [" + response.getStatus() + "]");
            }
            UnConnectedDataItem dataItem = (UnConnectedDataItem) rr.getTypeIds().get(1);
            if (dataItem.getService() instanceof CipConnectionManagerResponse cmr) {
                this.connectionId = cmr.getOtConnectionId();
                LOGGER.debug("Got assigned with Connection Id {}", this.connectionId);
            }
        });
    }

    @Override
    public void close() throws Exception {
        try {
            disconnectSession();
        } catch (Exception e) {
            LOGGER.debug("Error during EIP disconnect", e);
        }
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        pendingRequests.forEach(future ->
            future.completeExceptionally(new PlcRuntimeException("Connection closed")));
        pendingRequests.clear();
        super.close();
        LOGGER.info("EIP TCP connection closed");
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    private void disconnectSession() {
        if (!isConnected() || sessionHandle == EMPTY_SESSION_HANDLE) {
            return;
        }
        if (connectionId != 0L) {
            UnConnectedDataItem exchange = new UnConnectedDataItem(
                new CipConnectionManagerCloseRequest(
                    (short) 2,
                    new LogicalSegment(new ClassID((byte) 0, (short) 6)),
                    new LogicalSegment(new InstanceID((byte) 0, (short) 1)),
                    (byte) 0, (byte) 10, (short) 14,
                    this.connectionSerialNumber, 4919, 42L,
                    this.connectionPathSize, this.routingAddress));
            List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);
            CipRRData closePkt = new CipRRData(sessionHandle, 0L, DEFAULT_SENDER_CONTEXT,
                0L, EMPTY_INTERFACE_HANDLE, 0, typeIds);
            try {
                sendRequest(closePkt).get(1, TimeUnit.SECONDS);
            } catch (Exception ignored) {
            }
        }
        EipDisconnectRequest disconnectRequest = new EipDisconnectRequest(
            sessionHandle, 0L, DEFAULT_SENDER_CONTEXT, 0L);
        try {
            // Many devices close the socket immediately after the disconnect — ignore
            sendRequest(disconnectRequest).get(50, TimeUnit.MILLISECONDS);
        } catch (Exception ignored) {
        }
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");
        PlcRuntimeException exception = new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause);
        pendingRequests.forEach(future -> future.completeExceptionally(exception));
        pendingRequests.clear();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Request/response correlation (FIFO, max-concurrent == 1)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private CompletableFuture<EipPacket> sendRequest(EipPacket request) {
        CompletableFuture<EipPacket> responseFuture = new CompletableFuture<>();
        pendingRequests.offer(responseFuture);

        try {
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.OUTGOING_MESSAGE, "Sending EIP request");
            }
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingRequests.remove(responseFuture);
            responseFuture.completeExceptionally(new PlcRuntimeException("Failed to send EIP request", e));
            return responseFuture;
        }

        long timeoutMs = getConfiguration().getRequestTimeout();
        responseFuture.orTimeout(timeoutMs, TimeUnit.MILLISECONDS)
            .whenComplete((r, e) -> {
                if (e instanceof TimeoutException) {
                    pendingRequests.remove(responseFuture);
                }
            });

        return responseFuture;
    }

    private void handleIncomingMessage(EipPacket packet) {
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.INCOMING_MESSAGE, "Received EIP response");
        }
        CompletableFuture<EipPacket> future = pendingRequests.poll();
        if (future != null) {
            future.complete(packet);
        } else {
            LOGGER.warn("Received EIP response with no pending request");
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        if (getConfiguration().isForceUnconnectedOperation() || (!useMessageRouter && !useConnectionManager)) {
            return readWithoutMessageRouter(readRequest);
        } else if (useMessageRouter && !useConnectionManager) {
            return readWithoutConnectionManager(readRequest);
        } else {
            return readWithConnectionManager(readRequest);
        }
    }

    private CompletableFuture<PlcReadResponse> readWithoutMessageRouter(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        Map<String, PlcResponseItem<PlcValue>> values = new ConcurrentHashMap<>();
        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

        // Chain per-tag CipRRData wraps so executeThrottled doesn't block the caller.
        List<CompletableFuture<Void>> tagFutures = new ArrayList<>();
        CompletableFuture<Void> chain = CompletableFuture.completedFuture(null);
        for (String tagName : request.getTagNames()) {
            // A tag whose address the builder couldn't parse is kept in the request with an
            // error code and a null tag - report that code instead of sending a request built
            // from it.
            PlcResponseCode requestCode = request.getTagResponseCode(tagName);
            if (requestCode != PlcResponseCode.OK) {
                values.put(tagName, new DefaultPlcResponseItem<>(requestCode, null));
                continue;
            }
            EipTag eipTag = (EipTag) request.getTag(tagName);
            CompletableFuture<Void> tagFuture = chain.thenComposeAsync(v -> executeThrottled(() -> {
                try {
                    CipReadRequest req = new CipReadRequest(toAnsi(eipTag.getTag()), elementCount(eipTag));
                    CipUnconnectedRequest requestItem = new CipUnconnectedRequest(
                        classSegment, instanceSegment, req,
                        (byte) getConfiguration().getBackplane(),
                        (byte) getConfiguration().getSlot());
                    List<TypeId> typeIds = Arrays.asList(nullAddressItem,
                        new UnConnectedDataItem(requestItem));
                    CipRRData rrdata = new CipRRData(sessionHandle, CIPStatus.Success.getValue(),
                        DEFAULT_SENDER_CONTEXT, 0L, EMPTY_INTERFACE_HANDLE, 0, typeIds);
                    return sendRequest(rrdata).thenApply(response -> {
                        if (!(response instanceof CipRRData rr)) {
                            values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                            return null;
                        }
                        UnConnectedDataItem dataItem = (UnConnectedDataItem) rr.getTypeIds().get(1);
                        if (dataItem.getService() instanceof CipConnectedResponse ccr && ccr.getStatus() == 0x03) {
                            values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
                        } else {
                            values.putAll(decodeSingleReadResponse(dataItem.getService(), tagName, eipTag));
                        }
                        return null;
                    });
                } catch (BufferException e) {
                    values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                    return CompletableFuture.<Void>completedFuture(null);
                }
            }));
            tagFutures.add(tagFuture);
            chain = tagFuture.handle((r, e) -> null);
        }

        return CompletableFuture.allOf(tagFutures.toArray(new CompletableFuture[0]))
            .thenApply(v -> {
                Map<String, PlcResponseItem<PlcValue>> ordered = new LinkedHashMap<>();
                for (String tn : request.getTagNames()) {
                    ordered.put(tn, values.getOrDefault(tn,
                        new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null)));
                }
                return (PlcReadResponse) new DefaultPlcReadResponse(request, ordered);
            });
    }

    private CompletableFuture<PlcReadResponse> readWithoutConnectionManager(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

        List<CipService> requests = new ArrayList<>(request.getNumberOfTags());
        for (String tagName : sendableTagNames(request)) {
            EipTag eipTag = (EipTag) request.getTag(tagName);
            try {
                requests.add(new CipReadRequest(toAnsi(eipTag.getTag()), elementCount(eipTag)));
            } catch (BufferException e) {
                return CompletableFuture.failedFuture(new PlcRuntimeException("Failed to read field", e));
            }
        }

        List<TypeId> typeIds = new ArrayList<>(2);
        typeIds.add(nullAddressItem);
        CipService outerService;
        if (requests.size() == 1) {
            outerService = requests.getFirst();
        } else {
            List<Integer> offsets = new ArrayList<>(requests.size());
            int offset = 2 + 2 * requests.size();
            for (CipService cipRequest : requests) {
                offsets.add(offset);
                offset += cipRequest.getLengthInBytes();
            }
            outerService = new MultipleServiceRequest(new Services(offsets, requests));
        }
        typeIds.add(new UnConnectedDataItem(new CipUnconnectedRequest(
            classSegment, instanceSegment, outerService,
            (byte) getConfiguration().getBackplane(),
            (byte) getConfiguration().getSlot())));

        CipRRData pkt = new CipRRData(sessionHandle, CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT, 0L, 0L, 0, typeIds);

        return executeThrottled(() -> sendRequest(pkt).thenApply(response -> {
            if (!(response instanceof CipRRData rr)) {
                return (PlcReadResponse) new DefaultPlcReadResponse(request, errorMap(request));
            }
            UnConnectedDataItem dataItem = (UnConnectedDataItem) rr.getTypeIds().get(1);
            return decodeReadResponse(dataItem.getService(), request);
        }));
    }

    private CompletableFuture<PlcReadResponse> readWithConnectionManager(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<CipService> requests = new ArrayList<>(request.getNumberOfTags());
        for (String tagName : sendableTagNames(request)) {
            EipTag eipTag = (EipTag) request.getTag(tagName);
            try {
                requests.add(new CipReadRequest(toAnsi(eipTag.getTag()), elementCount(eipTag)));
            } catch (BufferException e) {
                return CompletableFuture.failedFuture(new PlcRuntimeException("Failed to read field", e));
            }
        }

        ConnectedAddressItem addressItem = new ConnectedAddressItem(this.connectionId);
        List<TypeId> typeIds = new ArrayList<>(2);
        typeIds.add(addressItem);
        if (requests.size() == 1) {
            typeIds.add(new ConnectedDataItem(this.sequenceCount, requests.getFirst()));
        } else {
            List<Integer> offsets = new ArrayList<>(requests.size());
            int offset = 2 + 2 * requests.size();
            for (CipService cipRequest : requests) {
                offsets.add(offset);
                offset += cipRequest.getLengthInBytes();
            }
            typeIds.add(new ConnectedDataItem(this.sequenceCount,
                new MultipleServiceRequest(new Services(offsets, requests))));
        }
        SendUnitData pkt = new SendUnitData(sessionHandle, CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT, 0L, 0, typeIds);
        this.sequenceCount += 1;

        return executeThrottled(() -> sendRequest(pkt).thenApply(response -> {
            if (!(response instanceof SendUnitData sud)) {
                return (PlcReadResponse) new DefaultPlcReadResponse(request, errorMap(request));
            }
            ConnectedDataItem dataItem = (ConnectedDataItem) sud.getTypeIds().get(1);
            return decodeReadResponse(dataItem.getService(), request);
        }));
    }

    /**
     * The tag names that can actually be put on the wire. A tag whose address the builder
     * couldn't parse stays in the request with an error code and a {@code null} tag; it must
     * neither be sent nor occupy a slot in the response, which is mapped back to tags by
     * position. Request building and response decoding both filter through this method so the
     * two stay aligned.
     */
    private static List<String> sendableTagNames(PlcTagRequest request) {
        List<String> names = new ArrayList<>(request.getNumberOfTags());
        for (String tagName : request.getTagNames()) {
            if (request.getTagResponseCode(tagName) == PlcResponseCode.OK) {
                names.add(tagName);
            }
        }
        return names;
    }

    /** Response items for the tags the builder rejected, keyed by name. */
    private static Map<String, PlcResponseItem<PlcValue>> rejectedReadTags(PlcTagRequest request) {
        Map<String, PlcResponseItem<PlcValue>> rejected = new LinkedHashMap<>();
        for (String tagName : request.getTagNames()) {
            PlcResponseCode code = request.getTagResponseCode(tagName);
            if (code != PlcResponseCode.OK) {
                rejected.put(tagName, new DefaultPlcResponseItem<>(code, null));
            }
        }
        return rejected;
    }

    /** Response codes for the tags the builder rejected, keyed by name. */
    private static Map<String, PlcResponseCode> rejectedWriteTags(PlcTagRequest request) {
        Map<String, PlcResponseCode> rejected = new LinkedHashMap<>();
        for (String tagName : request.getTagNames()) {
            PlcResponseCode code = request.getTagResponseCode(tagName);
            if (code != PlcResponseCode.OK) {
                rejected.put(tagName, code);
            }
        }
        return rejected;
    }

    private Map<String, PlcResponseItem<PlcValue>> errorMap(DefaultPlcReadRequest request) {
        Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>(rejectedReadTags(request));
        for (String tn : sendableTagNames(request)) {
            values.put(tn, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
        }
        return values;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        if (getConfiguration().isForceUnconnectedOperation() || (!useMessageRouter && !useConnectionManager)) {
            return writeWithoutMessageRouter(writeRequest);
        } else if (useMessageRouter && !useConnectionManager) {
            return writeWithoutConnectionManager(writeRequest);
        } else {
            return writeWithConnectionManager(writeRequest);
        }
    }

    private CompletableFuture<PlcWriteResponse> writeWithoutMessageRouter(PlcWriteRequest writeRequest) {
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));
        Map<String, PlcResponseCode> values = new ConcurrentHashMap<>();

        List<CompletableFuture<Void>> tagFutures = new ArrayList<>();
        CompletableFuture<Void> chain = CompletableFuture.completedFuture(null);
        values.putAll(rejectedWriteTags(request));
        for (String fieldName : sendableTagNames(request)) {
            EipTag field = (EipTag) request.getTag(fieldName);
            PlcValue value = request.getPlcValue(fieldName);
            int elements = Math.max(field.getElementNb(), 1);
            CompletableFuture<Void> tagFuture = chain.thenComposeAsync(v -> executeThrottled(() -> {
                try {
                    byte[] data = encodeValue(value, field.getType());
                    CipWriteRequest writeReq = new CipWriteRequest(
                        toAnsi(field.getTag()), field.getType(), elements, data);
                    CipUnconnectedRequest requestItem = new CipUnconnectedRequest(
                        classSegment, instanceSegment, writeReq,
                        (byte) getConfiguration().getBackplane(),
                        (byte) getConfiguration().getSlot());
                    List<TypeId> typeIds = Arrays.asList(nullAddressItem,
                        new UnConnectedDataItem(requestItem));
                    CipRRData rrdata = new CipRRData(sessionHandle, 0L,
                        DEFAULT_SENDER_CONTEXT, EMPTY_INTERFACE_OPTIONS, EMPTY_INTERFACE_HANDLE, 0, typeIds);
                    return sendRequest(rrdata).thenApply(response -> {
                        if (response instanceof CipRRData rr) {
                            UnConnectedDataItem di = (UnConnectedDataItem) rr.getTypeIds().get(1);
                            if (di.getService() instanceof CipWriteResponse cwr) {
                                values.put(fieldName, decodeResponseCode(cwr.getStatus()));
                                return null;
                            }
                        }
                        values.put(fieldName, PlcResponseCode.INTERNAL_ERROR);
                        return null;
                    });
                } catch (BufferException e) {
                    values.put(fieldName, PlcResponseCode.INTERNAL_ERROR);
                    return CompletableFuture.<Void>completedFuture(null);
                }
            }));
            tagFutures.add(tagFuture);
            chain = tagFuture.handle((r, e) -> null);
        }

        return CompletableFuture.allOf(tagFutures.toArray(new CompletableFuture[0]))
            .thenApply(v -> {
                Map<String, PlcResponseCode> ordered = new LinkedHashMap<>();
                for (String tn : request.getTagNames()) {
                    ordered.put(tn, values.getOrDefault(tn, PlcResponseCode.INTERNAL_ERROR));
                }
                return (PlcWriteResponse) new DefaultPlcWriteResponse(request, ordered);
            });
    }

    private CompletableFuture<PlcWriteResponse> writeWithoutConnectionManager(PlcWriteRequest writeRequest) {
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<CipWriteRequest> items = new ArrayList<>(writeRequest.getNumberOfTags());
        for (String fieldName : sendableTagNames(request)) {
            EipTag field = (EipTag) request.getTag(fieldName);
            PlcValue value = request.getPlcValue(fieldName);
            int elements = Math.max(field.getElementNb(), 1);
            byte[] data = encodeValue(value, field.getType());
            try {
                items.add(new CipWriteRequest(toAnsi(field.getTag()), field.getType(), elements, data));
            } catch (BufferException e) {
                return CompletableFuture.failedFuture(new PlcRuntimeException("Failed to write field", e));
            }
        }

        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));
        CipService outer;
        if (items.size() == 1) {
            outer = items.getFirst();
        } else {
            short nb = (short) items.size();
            List<Integer> offsets = new ArrayList<>(nb);
            int offset = 2 + nb * 2;
            for (CipWriteRequest item : items) {
                offsets.add(offset);
                offset += item.getLengthInBytes();
            }
            List<CipService> services = new ArrayList<>(items);
            outer = new MultipleServiceRequest(new Services(offsets, services));
        }

        UnConnectedDataItem exchange = new UnConnectedDataItem(new CipUnconnectedRequest(
            classSegment, instanceSegment, outer,
            (byte) getConfiguration().getBackplane(),
            (byte) getConfiguration().getSlot()));
        List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);
        CipRRData pkt = new CipRRData(sessionHandle, 0L, DEFAULT_SENDER_CONTEXT,
            0L, EMPTY_INTERFACE_HANDLE, 0, typeIds);

        return executeThrottled(() -> sendRequest(pkt).thenApply(response -> {
            if (!(response instanceof CipRRData rr)) {
                return (PlcWriteResponse) new DefaultPlcWriteResponse(request, writeErrorMap(request));
            }
            UnConnectedDataItem di = (UnConnectedDataItem) rr.getTypeIds().get(1);
            return decodeWriteResponse(di.getService(), request);
        }));
    }

    private CompletableFuture<PlcWriteResponse> writeWithConnectionManager(PlcWriteRequest writeRequest) {
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<CipWriteRequest> items = new ArrayList<>(writeRequest.getNumberOfTags());
        for (String fieldName : sendableTagNames(request)) {
            EipTag field = (EipTag) request.getTag(fieldName);
            PlcValue value = request.getPlcValue(fieldName);
            int elements = Math.max(field.getElementNb(), 1);
            byte[] data = encodeValue(value, field.getType());
            try {
                items.add(new CipWriteRequest(toAnsi(field.getTag()), field.getType(), elements, data));
            } catch (BufferException e) {
                return CompletableFuture.failedFuture(new PlcRuntimeException("Failed to write field", e));
            }
        }
        ConnectedAddressItem addressItem = new ConnectedAddressItem(this.connectionId);
        CipService outer;
        if (items.size() == 1) {
            outer = items.getFirst();
        } else {
            short nb = (short) items.size();
            List<Integer> offsets = new ArrayList<>(nb);
            int offset = 2 + nb * 2;
            for (CipWriteRequest item : items) {
                offsets.add(offset);
                offset += item.getLengthInBytes();
            }
            List<CipService> services = new ArrayList<>(items);
            outer = new MultipleServiceRequest(new Services(offsets, services));
        }
        ConnectedDataItem exchange = new ConnectedDataItem(this.sequenceCount, outer);
        this.sequenceCount += 1;
        List<TypeId> typeIds = Arrays.asList(addressItem, exchange);
        SendUnitData pkt = new SendUnitData(sessionHandle, CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT, 0L, 0, typeIds);

        return executeThrottled(() -> sendRequest(pkt).thenApply(response -> {
            if (!(response instanceof SendUnitData sud)) {
                return (PlcWriteResponse) new DefaultPlcWriteResponse(request, writeErrorMap(request));
            }
            ConnectedDataItem di = (ConnectedDataItem) sud.getTypeIds().get(1);
            return decodeWriteResponse(di.getService(), request);
        }));
    }

    private Map<String, PlcResponseCode> writeErrorMap(DefaultPlcWriteRequest request) {
        Map<String, PlcResponseCode> values = new LinkedHashMap<>(rejectedWriteTags(request));
        for (String tn : sendableTagNames(request)) {
            values.put(tn, PlcResponseCode.INTERNAL_ERROR);
        }
        return values;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoders / decoders
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private byte[] toAnsi(String tag) throws BufferException {
        Pattern resourcePattern = Pattern.compile("([.\\[\\]])*([A-Za-z_0-9]+)");
        Matcher matcher = resourcePattern.matcher(tag);
        List<PathSegment> segments = new LinkedList<>();
        int lengthBytes = 0;
        while (matcher.find()) {
            String identifier = matcher.group(2);
            String qualifier = matcher.group(1);
            PathSegment newSegment;
            if ("[".equals(qualifier)) {
                newSegment = new LogicalSegment(new MemberID((byte) 0x00, Short.parseShort(identifier)));
            } else {
                newSegment = new DataSegment(new AnsiExtendedSymbolSegment(identifier,
                    (identifier.length() % 2 == 0) ? null : (short) 0));
            }
            segments.add(newSegment);
            lengthBytes += newSegment.getLengthInBytes();
        }
        WriteBufferByteBased buffer = messageCodec.createWriteBuffer(lengthBytes);
        for (PathSegment segment : segments) {
            segment.serialize(buffer);
        }
        return buffer.getBytes();
    }

    private PlcReadResponse decodeReadResponse(CipService p, PlcReadRequest readRequest) {
        Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>(rejectedReadTags(readRequest));
        List<String> sendable = sendableTagNames(readRequest);
        if (p instanceof CipReadResponse resp) {
            if (sendable.isEmpty()) {
                return new DefaultPlcReadResponse((DefaultPlcReadRequest) readRequest, values);
            }
            String tagName = sendable.getFirst();
            EipTag tag = (EipTag) readRequest.getTag(tagName);
            PlcResponseCode code = decodeResponseCode(resp.getStatus());
            PlcValue plcValue = null;
            if (code == PlcResponseCode.OK) {
                plcValue = parsePlcValue(tag, resp.getData().getData(), resp.getData().getDataType());
                if (plcValue == null) {
                    code = PlcResponseCode.INTERNAL_ERROR;
                }
            }
            values.put(tagName, new DefaultPlcResponseItem<>(code, plcValue));
        } else if (p instanceof MultipleServiceResponse responses) {
            int nb = responses.getServiceNb();
            List<CipService> arr = new ArrayList<>(nb);
            try {
                byte[] servicesData = responses.getServicesData();
                ReadBufferByteBased read = messageCodec.createReadBuffer(servicesData);
                int total = servicesData.length;
                for (int i = 0; i < nb; i++) {
                    int offset = responses.getOffsets().get(i) - responses.getOffsets().getFirst();
                    int length;
                    if (i == nb - 1) {
                        length = total - offset;
                    } else {
                        length = responses.getOffsets().get(i + 1) - offset - responses.getOffsets().getFirst();
                    }
                    arr.add(CipService.staticParse(read, false, length));
                }
            } catch (BufferException e) {
                throw new PlcRuntimeException(e);
            }
            Iterator<String> it = sendable.iterator();
            for (int i = 0; i < nb && it.hasNext(); i++) {
                String tagName = it.next();
                EipTag tag = (EipTag) readRequest.getTag(tagName);
                if (arr.get(i) instanceof CipReadResponse rr) {
                    PlcResponseCode code = rr.getStatus() == 0 ? PlcResponseCode.OK : PlcResponseCode.INTERNAL_ERROR;
                    PlcValue plcValue = null;
                    if (code == PlcResponseCode.OK) {
                        plcValue = parsePlcValue(tag, rr.getData().getData(), rr.getData().getDataType());
                        if (plcValue == null) {
                            code = PlcResponseCode.INTERNAL_ERROR;
                        }
                    }
                    values.put(tagName, new DefaultPlcResponseItem<>(code, plcValue));
                }
            }
        }
        return new DefaultPlcReadResponse(readRequest, values);
    }

    private Map<String, PlcResponseItem<PlcValue>> decodeSingleReadResponse(CipService p, String tagName, PlcTag tag) {
        Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
        if (!(p instanceof CipReadResponse resp)) {
            values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
            return values;
        }
        PlcResponseCode code = decodeResponseCode(resp.getStatus());
        PlcValue plcValue = null;
        if (code == PlcResponseCode.OK) {
            plcValue = parsePlcValue((EipTag) tag, resp.getData().getData(), resp.getData().getDataType());
            if (plcValue == null) {
                // Undecodable payload (unsupported type or a reply shorter than the requested
                // number of elements) - don't report that as a successful read.
                code = PlcResponseCode.INTERNAL_ERROR;
            }
        }
        values.put(tagName, new DefaultPlcResponseItem<>(code, plcValue));
        return values;
    }

    /**
     * Number of elements to ask the device for. An array tag ({@code %tag[0]:DINT:8}) has to
     * request all of its elements - see GH-1008 - otherwise the device returns a single element
     * and the decoder below has nothing to read the remaining ones from.
     */
    private static int elementCount(EipTag tag) {
        return Math.max(tag.getElementNb(), 1);
    }

    /** Whether the type is stored as a flat sequence of equally sized elements. */
    private static boolean isFixedSize(CIPDataTypeCode type) {
        return switch (type) {
            case SINT, INT, DINT, LINT, REAL, LREAL, BOOL -> true;
            default -> false;
        };
    }

    // Package-private and static so the decoding can be tested without a connection.
    static PlcValue parsePlcValue(EipTag tag, byte[] rawData, CIPDataTypeCode type) {
        final int STRING_LEN_OFFSET = 2;
        final int STRING_DATA_OFFSET = 6;
        ByteBuffer data = ByteBuffer.wrap(rawData).order(ByteOrder.LITTLE_ENDIAN);
        int nb = elementCount(tag);
        if (nb > 1) {
            // Never read past what the device actually sent - a short reply must not blow up
            // the whole response with an IndexOutOfBoundsException. Only the fixed-size types
            // below are laid out element by element; STRING/STRUCTURED carry their own length.
            int elementSize = isFixedSize(type) ? type.getSize() : 0;
            if (elementSize > 0 && rawData.length < nb * elementSize) {
                LOGGER.warn("Device returned {} bytes for tag '{}', expected {} for {} elements of {}.",
                    rawData.length, tag.getTag(), nb * elementSize, nb, type);
                return null;
            }
            List<PlcValue> list = new ArrayList<>();
            int index = 0;
            for (int i = 0; i < nb; i++) {
                switch (type) {
                    case DINT:
                        list.add(new PlcDINT(data.getInt(index)));
                        index += type.getSize();
                        break;
                    case INT:
                        list.add(new PlcINT(data.getShort(index)));
                        index += type.getSize();
                        break;
                    case SINT:
                        list.add(new PlcSINT(data.get(index)));
                        index += type.getSize();
                        break;
                    case REAL:
                        list.add(new PlcREAL(data.getFloat(index)));
                        index += type.getSize();
                        break;
                    case LREAL:
                        list.add(new PlcLREAL(data.getDouble(index)));
                        index += type.getSize();
                        break;
                    case LINT:
                        list.add(new PlcLINT(data.getLong(index)));
                        index += type.getSize();
                        break;
                    case BOOL:
                        list.add(new PlcBOOL(data.get(index) != 0));
                        index += type.getSize();
                        break;
                    case STRING:
                    case STRUCTURED:
                        short structuredType = data.getShort(0);
                        short structuredLen = data.getShort(STRING_LEN_OFFSET);
                        if (structuredType == CIPStructTypeCode.STRING.getValue()) {
                            list.add(new PlcSTRING(new String(rawData, STRING_DATA_OFFSET, structuredLen, StandardCharsets.UTF_8)));
                        }
                        return list.isEmpty() ? null : new PlcList(list);
                    default:
                        return null;
                }
            }
            return new PlcList(list);
        }
        switch (type) {
            case SINT:
                return new PlcSINT(data.get(0));
            case INT:
                return new PlcINT(data.getShort(0));
            case DINT:
                return new PlcDINT(data.getInt(0));
            case LINT:
                return new PlcLINT(data.getLong(0));
            case REAL:
                return new PlcREAL(data.getFloat(0));
            case LREAL:
                return new PlcLREAL(data.getDouble(0));
            case BOOL:
                return new PlcBOOL(data.get(0) != 0);
            case STRING:
            case STRUCTURED:
                short structuredType = data.getShort(0);
                short structuredLen = data.getShort(STRING_LEN_OFFSET);
                if (structuredType == CIPStructTypeCode.STRING.getValue()) {
                    return new PlcSTRING(new String(rawData, STRING_DATA_OFFSET, structuredLen, StandardCharsets.UTF_8));
                }
                return null;
            default:
                return null;
        }
    }

    private PlcWriteResponse decodeWriteResponse(CipService p, PlcWriteRequest writeRequest) {
        Map<String, PlcResponseCode> responses = new LinkedHashMap<>(rejectedWriteTags(writeRequest));
        List<String> sendable = sendableTagNames(writeRequest);
        if (p instanceof CipWriteResponse resp) {
            if (sendable.isEmpty()) {
                return new DefaultPlcWriteResponse((DefaultPlcWriteRequest) writeRequest, responses);
            }
            String fieldName = sendable.getFirst();
            responses.put(fieldName, decodeResponseCode(resp.getStatus()));
            return new DefaultPlcWriteResponse(writeRequest, responses);
        } else if (p instanceof MultipleServiceResponse resp) {
            int nb = resp.getServiceNb();
            List<CipService> arr = new ArrayList<>(nb);
            try {
                byte[] servicesData = resp.getServicesData();
                ReadBufferByteBased read = messageCodec.createReadBuffer(servicesData);
                int total = servicesData.length;
                for (int i = 0; i < nb; i++) {
                    int offset = resp.getOffsets().get(i);
                    int length;
                    if (offset == nb - 1) {
                        length = total - offset;
                    } else {
                        length = resp.getOffsets().get(i + 1) - offset;
                    }
                    arr.add(CipService.staticParse(read, false, length));
                }
            } catch (BufferException e) {
                throw new PlcRuntimeException(e);
            }
            Iterator<String> it = sendable.iterator();
            for (int i = 0; i < nb && it.hasNext(); i++) {
                String fieldName = it.next();
                if (arr.get(i) instanceof CipWriteResponse writeResponse) {
                    responses.put(fieldName, decodeResponseCode(writeResponse.getStatus()));
                }
            }
            return new DefaultPlcWriteResponse(writeRequest, responses);
        }
        for (String tn : sendable) {
            responses.put(tn, PlcResponseCode.INTERNAL_ERROR);
        }
        return new DefaultPlcWriteResponse(writeRequest, responses);
    }

    private byte[] encodeValue(PlcValue value, CIPDataTypeCode type) {
        ByteBuffer buffer = ByteBuffer.allocate(type.getSize()).order(ByteOrder.LITTLE_ENDIAN);
        switch (type) {
            case BOOL:
            case SINT:
                buffer.put(value.getByte());
                break;
            case INT:
                buffer.putShort(value.getShort());
                break;
            case DINT:
                buffer.putInt(value.getInteger());
                break;
            case LINT:
                buffer.putLong(value.getLong());
                break;
            case REAL:
                buffer.putFloat(value.getFloat());
                break;
            case LREAL:
                buffer.putDouble(value.getDouble());
                break;
            case STRING:
            case STRUCTURED:
                buffer.putInt(value.getString().length());
                buffer.put(value.getString().getBytes(StandardCharsets.UTF_8), 0, value.getString().length());
                break;
            default:
                break;
        }
        return buffer.array();
    }

    private PlcResponseCode decodeResponseCode(int status) {
        return status == 0 ? PlcResponseCode.OK : PlcResponseCode.INTERNAL_ERROR;
    }

}
