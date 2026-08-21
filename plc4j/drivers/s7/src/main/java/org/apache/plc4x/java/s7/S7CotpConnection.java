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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.s7.readwrite.*;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.s7.context.S7DriverContext;
import org.apache.plc4x.java.s7.optimizer.S7BlockReadOptimizer;
import org.apache.plc4x.java.s7.optimizer.S7ReadChunk;
import org.apache.plc4x.java.s7.optimizer.S7WriteChunk;
import org.apache.plc4x.java.s7.tag.S7StringFixedLengthTag;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.s7.tag.S7PlcTagHandler;
import org.apache.plc4x.java.s7.tag.S7AlarmTag;
import org.apache.plc4x.java.s7.userdata.S7AlarmQueryHandle;
import org.apache.plc4x.java.s7.userdata.S7AlarmSubscriptionHandle;
import org.apache.plc4x.java.s7.userdata.S7AlarmSubscriptionService;
import org.apache.plc4x.java.s7.userdata.S7BrowseService;
import org.apache.plc4x.java.s7.userdata.S7CyclicSubscriptionHandle;
import org.apache.plc4x.java.s7.userdata.S7CyclicSubscriptionService;
import org.apache.plc4x.java.s7.userdata.S7SzlService;
import org.apache.plc4x.java.s7.userdata.S7SzlService.S7DeviceIdentification;
import org.apache.plc4x.java.s7.userdata.S7UserDataPushDispatcher;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.messages.PlcBrowseRequestInterceptor;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.SzlId;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcPingResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.stream.IntStream;

/**
 * S7 connection over the COTP transport. The COTP transport handles TPKT framing and the
 * COTP CR/CC handshake; this connection only deals with raw S7 PDUs (the inner payload of
 * COTP-DT packets, starting with magic 0x32).
 *
 * <p>The S7 SetupCommunication exchange happens in {@link #onConnect()} once the transport is
 * already open. After that, read/write requests use a tpdu-reference-keyed pendingRequests map
 * for response correlation.
 */
public class S7CotpConnection extends ConnectionBase<S7Configuration> {

    /** Safety net so a device that always claims "more data" cannot loop forever. */
    private static final int MAX_SZL_CONTINUATIONS = 16;

    private static final Logger LOGGER = LoggerFactory.getLogger(S7CotpConnection.class);

    private final AtomicInteger tpduGenerator = new AtomicInteger(10);
    private final Map<Integer, CompletableFuture<S7Message>> pendingRequests = new ConcurrentHashMap<>();
    private final S7DriverContext driverContext;
    private final S7BlockReadOptimizer optimizer = new S7BlockReadOptimizer();
    private final S7UserDataPushDispatcher pushDispatcher = new S7UserDataPushDispatcher();
    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();
    private final Set<S7AlarmSubscriptionHandle> alarmHandles = ConcurrentHashMap.newKeySet();
    private final List<java.io.Closeable> alarmDispatcherHooks = new ArrayList<>();
    private final Map<Short, List<S7CyclicSubscriptionHandle>> cyclicHandlesByJob = new ConcurrentHashMap<>();
    private final List<java.io.Closeable> cyclicDispatcherHooks = new ArrayList<>();

    private S7CotpMessageCodec messageCodec;

    public S7CotpConnection(S7Configuration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
        this.driverContext = new S7DriverContext();
        this.driverContext.applyConfiguration(configuration);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new S7PlcTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    /** Public accessor needed by {@link S7HCotpConnection} when wrapping two inner connections. */
    @Override
    public S7Configuration getConfiguration() {
        return super.getConfiguration();
    }

    @Override
    protected int getMaxConcurrentRequests() {
        // ConnectionBase invokes this from its constructor — before our fields are
        // initialized — so we read straight from the configuration. The negotiated
        // max-amq-callee from SetupCommunication doesn't change throttle behavior since
        // we serialize at the codec level via tpduRef anyway.
        return Math.max(1, getConfiguration().getMaxAmqCallee());
    }

    @Override
    public boolean isConnected() {
        return messageCodec != null && messageCodec.isOpen();
    }

    /**
     * Per-device capability flags. Read/Write var are always supported. Browse and
     * Subscribe both ride S7Comm UserData services, so they're enabled iff the SZL
     * probe at connect time succeeded.
     */
    @Override
    public PlcConnectionMetadata getMetadata() {
        boolean userData = driverContext.isUserDataServicesSupported();
        return new PlcConnectionMetadata() {
            @Override public boolean isReadSupported()      { return true; }
            @Override public boolean isWriteSupported()     { return true; }
            @Override public boolean isSubscribeSupported() { return userData; }
            @Override public boolean isBrowseSupported()    { return userData; }
        };
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new S7CotpMessageCodec(transportInstance, this::handleIncomingMessage);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming S7 data", e);
            }
        });

        try {
            doSetupCommunication().get(driverContext.getReadTimeout() * 2L, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("Error during S7 SetupCommunication", e);
        }

        // SZL capability probe. A successful response gives us the article number → controller
        // type and tells us the device speaks S7Comm UserData services (browse, alarms, cyclic
        // subscriptions). On failure (LOGO et al.), keep the connection — just mark UserData
        // unavailable so getMetadata() reflects reality and the higher layers don't try.
        probeUserDataCapability();

        LOGGER.info("S7 connection established (pdu-size={}, max-amq-caller={}, max-amq-callee={}, "
                + "controllerType={}, userdata-services={})",
            driverContext.getPduSize(), driverContext.getMaxAmqCaller(), driverContext.getMaxAmqCallee(),
            driverContext.getControllerType(), driverContext.isUserDataServicesSupported());
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT, "S7 connection established");
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    private void probeUserDataCapability() {
        // If the user pinned the controller type, trust it and skip the SZL probe entirely
        // (matches the pre-SPI3 driver, which only read SZL to auto-detect the type). The
        // UserData-service capability (browse/alarms/cyclic) is then derived from the known type.
        if (driverContext.getControllerType() != ControllerType.ANY) {
            boolean supported = supportsUserDataServices(driverContext.getControllerType());
            driverContext.setUserDataServicesSupported(supported);
            LOGGER.info("Controller type pinned to {}; skipping SZL probe (userdata-services={})",
                driverContext.getControllerType(), supported);
            return;
        }
        // Try the modern COMPONENT_IDENTIFICATION SZL first — works on S7-1200/1500 and is
        // tolerated by S7-300/400. Fall back to the legacy MODULE_IDENTIFICATION SZL if the
        // first attempt errors out or doesn't yield a recognisable article number. Either
        // success counts as "device speaks UserData services".
        S7SzlService.ProbeResult result = trySzlProbe(S7SzlService.COMPONENT_IDENTIFICATION, 0x0001);
        if (result == null) {
            result = trySzlProbe(S7SzlService.MODULE_IDENTIFICATION, 0x0000);
        }
        if (result != null) {
            driverContext.setArticleNumber(result.articleNumber());
            driverContext.setControllerType(result.controllerType());
            driverContext.setUserDataServicesSupported(true);
            LOGGER.info("SZL probe ok: article='{}', controllerType={}",
                result.articleNumber(), driverContext.getControllerType());
        } else {
            // Common on LOGO and other non-SZL devices. Connection stays usable for plain
            // Read/Write var; advanced services just stay disabled.
            driverContext.setUserDataServicesSupported(false);
            LOGGER.info("SZL probe yielded no usable identification; S7Comm UserData services "
                + "disabled for this device");
        }
    }

    /**
     * @return whether the given (explicitly configured) controller type speaks the S7Comm
     * UserData services that back browse, alarm and cyclic subscriptions. Used when the SZL
     * probe is skipped because the user pinned the controller type.
     */
    static boolean supportsUserDataServices(ControllerType type) {
        return switch (type) {
            case S7_300, S7_400, S7_1200, S7_1500 -> true;
            default -> false; // S7_200, LOGO, ANY
        };
    }

    private S7SzlService.ProbeResult trySzlProbe(SzlId szlId, int szlIndex) {
        try {
            // Draining matters even for the probe: a chained response the driver walks away
            // from leaves the sequence open, and the CPU then refuses every later SZL read on
            // this connection with 0xD406.
            byte[] szlData = readSzlData(szlId, szlIndex, this::sendInternal)
                .get(driverContext.getReadTimeout() * (MAX_SZL_CONTINUATIONS + 1L),
                    TimeUnit.MILLISECONDS);
            return S7SzlService.parseProbeResponse(szlData);
        } catch (Exception e) {
            LOGGER.debug("SZL {}/{} probe failed: {}", szlId.getSublistList(), szlIndex, e.getMessage());
            return null;
        }
    }

    /**
     * Read one SZL list, following the chain if the PLC splits it across several PDUs.
     *
     * <p>A response with {@code lastDataUnit} set means the CPU is holding the rest of the
     * list. Continuation chunks carry no header of their own - they simply continue the byte
     * stream, possibly starting mid-record - so the chunks are concatenated and parsed as one.
     *
     * @param sender how to dispatch each request; the connect-time probe sends directly,
     *               while normal operation goes through the request throttle.
     */
    private CompletableFuture<byte[]> readSzlData(SzlId szlId, int szlIndex,
            Function<S7Message, CompletableFuture<S7Message>> sender) {
        S7Message request = S7SzlService.buildRequest(getTpduId(), szlId, szlIndex);
        return sender.apply(request).thenCompose(response ->
            appendSzlChunks(szlId, szlIndex, response, S7SzlService.szlDataOf(response), 0, sender));
    }

    private CompletableFuture<byte[]> appendSzlChunks(SzlId szlId, int szlIndex, S7Message response,
            byte[] soFar, int continuations,
            Function<S7Message, CompletableFuture<S7Message>> sender) {
        if (!S7SzlService.moreDataFollows(response)) {
            return CompletableFuture.completedFuture(soFar);
        }
        if (continuations >= MAX_SZL_CONTINUATIONS) {
            // Give up rather than loop forever. The sequence stays open, so later SZL reads on
            // this connection will fail - worth a warning rather than silent truncation.
            LOGGER.warn("SZL {}/{} still reports more data after {} continuations; giving up "
                + "with {} bytes. Further SZL reads on this connection may be refused.",
                szlId.getSublistList(), szlIndex, continuations, soFar.length);
            return CompletableFuture.completedFuture(soFar);
        }
        S7Message next = S7SzlService.buildContinuationRequest(getTpduId(), szlId, szlIndex,
            S7SzlService.sequenceNumberOf(response));
        return sender.apply(next).thenCompose(chunk ->
            appendSzlChunks(szlId, szlIndex, chunk, concat(soFar, S7SzlService.szlDataOf(chunk)),
                continuations + 1, sender));
    }

    private static byte[] concat(byte[] first, byte[] second) {
        byte[] combined = new byte[first.length + second.length];
        System.arraycopy(first, 0, combined, 0, first.length);
        System.arraycopy(second, 0, combined, first.length, second.length);
        return combined;
    }

    /**
     * Read everything the CPU is willing to tell us about itself: order code, hardware and
     * firmware version, the configured names, the serial number and — where the device
     * implements it — the protection level currently in force.
     *
     * <p>The connect-time probe deliberately keeps only what the driver itself needs; this is
     * the on-demand version for callers building a device inventory. Three SZL lists are read
     * one after the other and merged. Each one is best-effort: devices implement different
     * subsets (S7-1200/1500 typically reject the protection-status list), so a rejected or
     * unanswered list leaves its fields {@code null} instead of failing the whole call.
     *
     * @return the merged identification, never {@code null}; {@link S7DeviceIdentification#EMPTY}
     * if the device answered none of the lists.
     */
    public CompletableFuture<S7DeviceIdentification> readDeviceIdentification() {
        return readSzl(S7SzlService.MODULE_IDENTIFICATION, 0x0000,
            S7SzlService::parseModuleIdentification)
            .thenCompose(module -> readSzl(S7SzlService.COMPONENT_IDENTIFICATION, 0x0000,
                S7SzlService::parseComponentIdentification)
                .thenApply(module::mergedWith))
            .thenCompose(identification -> readSzl(S7SzlService.PROTECTION_STATUS, 0x0004,
                S7SzlService::parseProtectionStatus)
                .thenApply(identification::mergedWith));
    }

    /**
     * Issue one SZL read and parse it, collapsing every failure mode — transport error, PLC
     * rejection, unparsable payload, no answer at all — into an empty identification. The
     * timeout matters: a device that simply ignores an SZL it doesn't implement would
     * otherwise leave the caller waiting forever.
     */
    private CompletableFuture<S7DeviceIdentification> readSzl(
            SzlId szlId, int szlIndex, Function<byte[], S7DeviceIdentification> parser) {
        return readSzlData(szlId, szlIndex, request -> executeThrottled(() -> sendInternal(request))
                .orTimeout(driverContext.getReadTimeout(), TimeUnit.MILLISECONDS))
            .handle((szlData, error) -> {
                if (error != null) {
                    LOGGER.debug("SZL {}/{} read failed: {}",
                        szlId.getSublistList(), szlIndex, error.getMessage());
                    return S7DeviceIdentification.EMPTY;
                }
                try {
                    return parser.apply(szlData);
                } catch (Exception e) {
                    LOGGER.debug("SZL {}/{} response could not be parsed: {}",
                        szlId.getSublistList(), szlIndex, e.getMessage());
                    return S7DeviceIdentification.EMPTY;
                }
            });
    }

    private CompletableFuture<Void> doSetupCommunication() {
        S7ParameterSetupCommunication setup = new S7ParameterSetupCommunication(
            driverContext.getMaxAmqCaller(), driverContext.getMaxAmqCallee(), driverContext.getPduSize());
        S7Message request = new S7MessageRequest(getTpduId(), setup, null);
        return sendInternal(request).thenAccept(response -> {
            if (response.getParameter() instanceof S7ParameterSetupCommunication setupResp) {
                driverContext.setMaxAmqCaller(setupResp.getMaxAmqCaller());
                driverContext.setMaxAmqCallee(setupResp.getMaxAmqCallee());
                driverContext.setPduSize(setupResp.getPduLength());
                setMaxConcurrentRequests(driverContext.getMaxAmqCallee());
            } else {
                throw new PlcRuntimeException("Unexpected response to S7 SetupCommunication: " +
                    (response.getParameter() == null ? "null" : response.getParameter().getClass().getSimpleName()));
            }
        });
    }

    @Override
    public void close() throws Exception {
        // Drain active subscriptions cleanly before tearing the codec down: tell the PLC to
        // stop pushing alarms / cyclic data so we don't leak server-side state and so the
        // cancel requests actually make the wire (without this drain they race the TCP FIN
        // and get dropped by the codec shutdown).
        drainActiveSubscriptions();
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        pendingRequests.values().forEach(f ->
            f.completeExceptionally(new PlcRuntimeException("Connection closed")));
        pendingRequests.clear();
        super.close();
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    private void drainActiveSubscriptions() {
        List<CompletableFuture<Void>> drains = new ArrayList<>();
        if (!alarmHandles.isEmpty()) {
            alarmHandles.clear();
            drains.add(deactivateAlarmSubscription());
        }
        if (!cyclicHandlesByJob.isEmpty()) {
            Set<Short> jobs = new java.util.HashSet<>(cyclicHandlesByJob.keySet());
            cyclicHandlesByJob.clear();
            drains.add(cancelCyclicJobs(jobs));
        }
        if (drains.isEmpty()) {
            return;
        }
        try {
            // Bounded wait — if a cancel times out, we still close cleanly afterwards.
            CompletableFuture.allOf(drains.toArray(CompletableFuture[]::new))
                .get(driverContext.getReadTimeout(), TimeUnit.MILLISECONDS);
        } catch (Exception ignored) {
            // Best-effort drain; continue with close regardless.
        }
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");
        PlcRuntimeException ex = new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause);
        pendingRequests.values().forEach(f -> f.completeExceptionally(ex));
        pendingRequests.clear();
    }

    private void handleIncomingMessage(S7Message message) {
        int tpduRef = message.getTpduReference();
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.INCOMING_MESSAGE, "Received S7 response, tpduRef=" + tpduRef);
        }
        CompletableFuture<S7Message> future = pendingRequests.remove(tpduRef);
        if (future != null) {
            future.complete(message);
            return;
        }
        // No pending request — could be an unsolicited UserData push (cyclic data, alarms,
        // mode-change). Route through the dispatcher; only complain if nobody owned it.
        if (pushDispatcher.dispatch(message)) {
            return;
        }
        LOGGER.warn("Received S7 response for unknown tpduRef {}", tpduRef);
    }

    /** Visible for tests / future service wiring. */
    public S7UserDataPushDispatcher getPushDispatcher() {
        return pushDispatcher;
    }

    private CompletableFuture<S7Message> sendInternal(S7Message request) {
        int tpduRef = request.getTpduReference();
        CompletableFuture<S7Message> future = new CompletableFuture<>();
        pendingRequests.put(tpduRef, future);
        try {
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.OUTGOING_MESSAGE, "Sending S7 request, tpduRef=" + tpduRef);
            }
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingRequests.remove(tpduRef);
            future.completeExceptionally(new PlcRuntimeException("Failed to send S7 request", e));
            return future;
        }
        long timeout = Math.max(1000, driverContext.getReadTimeout());
        future.orTimeout(timeout, TimeUnit.MILLISECONDS).whenComplete((r, e) -> {
            if (e instanceof TimeoutException) {
                pendingRequests.remove(tpduRef);
            }
        });
        return future;
    }

    private int getTpduId() {
        int id = tpduGenerator.getAndIncrement();
        if (id >= 0xFFFE) {
            tpduGenerator.set(10);
        }
        return id & 0xFFFF;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Ping
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        // No dedicated ping in S7. Issue a tiny zero-tag SetupCommunication-equivalent: a
        // bare ReadVar with one M0 byte. If we can build a request and get any response we're alive.
        S7Address addr = new S7AddressAny(TransportSize.BYTE, 1, 0, MemoryArea.FLAGS_MARKERS, 0, (byte) 0);
        S7VarRequestParameterItemAddress item = new S7VarRequestParameterItemAddress(addr);
        S7Message request = new S7MessageRequest(getTpduId(),
            new S7ParameterReadVarRequest(Collections.singletonList(item)), null);
        return executeThrottled(() -> sendInternal(request))
            .handle((resp, err) -> new DefaultPlcPingResponse(pingRequest,
                err != null ? PlcResponseCode.INTERNAL_ERROR : PlcResponseCode.OK));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Read
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        // A tag whose address the builder couldn't parse stays in the request with an error code
        // and a null tag. It is reported per tag rather than failing the whole request, so one
        // typo doesn't stop the request's other tags from being read.
        Map<String, PlcResponseItem<PlcValue>> rejectedTags = new LinkedHashMap<>();
        for (String tagName : readRequest.getTagNames()) {
            PlcResponseCode requestCode = readRequest.getTagResponseCode(tagName);
            if (requestCode != PlcResponseCode.OK) {
                rejectedTags.put(tagName, new DefaultPlcResponseItem<>(requestCode, null));
                continue;
            }
            org.apache.plc4x.java.api.model.PlcTag tag = readRequest.getTag(tagName);
            if (!(tag instanceof S7Tag)) {
                // Not a parse failure but a tag object of the wrong type handed to the API.
                return CompletableFuture.failedFuture(
                    new PlcProtocolException("Unsupported tag type " + tag.getClass().getName()));
            }
        }
        if (rejectedTags.size() == readRequest.getTagNames().size()) {
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse(readRequest, rejectedTags));
        }

        List<S7ReadChunk> chunks;
        try {
            chunks = optimizer.splitReadRequest(readRequest, driverContext);
        } catch (PlcRuntimeException e) {
            // A tag the optimizer can't fit at all is reported per-tag rather than aborting the whole request.
            Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
            for (String t : readRequest.getTagNames()) {
                values.put(t, rejectedTags.getOrDefault(t,
                    new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_DATA, null)));
            }
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse(readRequest, values));
        }
        if (chunks.isEmpty()) {
            return CompletableFuture.completedFuture(
                new DefaultPlcReadResponse(readRequest, new LinkedHashMap<>()));
        }

        // For each chunk, send one S7 read message; combine the per-chunk decoders into the final response.
        Map<String, PlcResponseItem<PlcValue>> finalValues = new LinkedHashMap<>();
        // Pre-seed in original request order so output order is deterministic. Tags the builder
        // rejected already have their code and are never overwritten (only nulls are).
        for (String t : readRequest.getTagNames()) {
            finalValues.put(t, rejectedTags.get(t));
        }
        CompletableFuture<Void> chain = CompletableFuture.completedFuture(null);
        for (S7ReadChunk chunk : chunks) {
            chain = chain.thenCompose(v -> sendReadChunk(chunk, finalValues));
        }
        return chain.thenApply(v -> {
            // Replace any still-null entries (shouldn't happen) with INTERNAL_ERROR.
            for (Map.Entry<String, PlcResponseItem<PlcValue>> e : finalValues.entrySet()) {
                if (e.getValue() == null) {
                    e.setValue(new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
            }
            return new DefaultPlcReadResponse(readRequest, finalValues);
        });
    }

    private CompletableFuture<Void> sendReadChunk(S7ReadChunk chunk, Map<String, PlcResponseItem<PlcValue>> out) {
        List<S7VarRequestParameterItem> items = new ArrayList<>(chunk.slots().size());
        for (S7ReadChunk.Slot slot : chunk.slots()) {
            items.add(slot.requestItem());
        }
        S7Message request = new S7MessageRequest(getTpduId(),
            new S7ParameterReadVarRequest(items), null);
        return executeThrottled(() -> sendInternal(request))
            .handle((response, error) -> {
                if (error != null) {
                    for (S7ReadChunk.Slot slot : chunk.slots()) {
                        for (S7ReadChunk.Binding b : slot.bindings()) {
                            mergeBindingFailure(out, b, PlcResponseCode.INTERNAL_ERROR);
                        }
                    }
                    return null;
                }
                applyChunkResponse(chunk, response, out);
                return null;
            });
    }

    private void applyChunkResponse(S7ReadChunk chunk, S7Message responseMessage,
                                    Map<String, PlcResponseItem<PlcValue>> out) {
        short errorClass = 0, errorCode = 0;
        if (responseMessage instanceof S7MessageResponseData md) {
            errorClass = md.getErrorClass();
            errorCode = md.getErrorCode();
        } else if (responseMessage instanceof S7MessageResponse mr) {
            errorClass = mr.getErrorClass();
            errorCode = mr.getErrorCode();
        }
        if (errorClass != 0 || errorCode != 0) {
            PlcResponseCode code = mapPlcErrorCode(errorClass, errorCode);
            for (S7ReadChunk.Slot slot : chunk.slots()) {
                for (S7ReadChunk.Binding b : slot.bindings()) {
                    mergeBindingFailure(out, b, code);
                }
            }
            return;
        }
        if (!(responseMessage.getPayload() instanceof S7PayloadReadVarResponse payload)) {
            for (S7ReadChunk.Slot slot : chunk.slots()) {
                for (S7ReadChunk.Binding b : slot.bindings()) {
                    mergeBindingFailure(out, b, PlcResponseCode.INTERNAL_ERROR);
                }
            }
            return;
        }
        List<S7VarPayloadDataItem> payloadItems = payload.getItems();
        for (int i = 0; i < chunk.slots().size(); i++) {
            S7ReadChunk.Slot slot = chunk.slots().get(i);
            if (i >= payloadItems.size()) {
                for (S7ReadChunk.Binding b : slot.bindings()) {
                    mergeBindingFailure(out, b, PlcResponseCode.INTERNAL_ERROR);
                }
                continue;
            }
            S7VarPayloadDataItem item = payloadItems.get(i);
            PlcResponseCode code = mapDataTransportError(item.getReturnCode());
            byte[] data = item.getData();
            for (S7ReadChunk.Binding b : slot.bindings()) {
                if (code != PlcResponseCode.OK) {
                    mergeBindingFailure(out, b, code);
                    continue;
                }
                try {
                    decodeBindingInto(out, b, data);
                } catch (Exception e) {
                    LOGGER.warn("Error decoding tag {}", b.tagName(), e);
                    mergeBindingFailure(out, b, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }
    }

    private void decodeBindingInto(Map<String, PlcResponseItem<PlcValue>> out, S7ReadChunk.Binding b, byte[] data)
        throws BufferException {
        S7Tag originalTag = b.originalTag();
        if (b.isSplitFragment()) {
            // Split-array reassembly: we expect a BYTE array in PlcRawByteArray for now.
            int elementSize = originalTag.getDataType().getSizeInBytes();
            int byteOffsetInOriginal = b.elementOffset() * elementSize;
            PlcResponseItem<PlcValue> existing = out.get(b.tagName());
            byte[] full;
            if (existing != null && existing.getValue() instanceof PlcRawByteArray prba) {
                full = prba.getRaw();
            } else {
                full = new byte[originalTag.getNumberOfElements() * elementSize];
                out.put(b.tagName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcRawByteArray(full)));
            }
            int copyLen = Math.min(data.length, full.length - byteOffsetInOriginal);
            System.arraycopy(data, 0, full, byteOffsetInOriginal, copyLen);
            return;
        }
        // Slice block-merged payload.
        int sliceLen = sliceLengthFor(originalTag);
        byte[] slice;
        if (b.payloadByteOffset() == 0 && data.length == sliceLen) {
            slice = data;
        } else {
            slice = new byte[sliceLen];
            System.arraycopy(data, b.payloadByteOffset(), slice, 0, Math.min(sliceLen, data.length - b.payloadByteOffset()));
        }
        PlcValue value = parsePlcValue(originalTag, slice);
        out.put(b.tagName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK, value));
    }

    private static int sliceLengthFor(S7Tag tag) {
        if (tag.getDataType() == TransportSize.BOOL) {
            return Math.max(1, (tag.getNumberOfElements() + 7) / 8);
        }
        if (tag instanceof S7StringFixedLengthTag fixed) {
            int bytesPerChar = fixed.getDataType() == TransportSize.WSTRING ? 2 : 1;
            return tag.getNumberOfElements() * (fixed.getStringLength() + 2) * bytesPerChar;
        }
        return tag.getNumberOfElements() * tag.getDataType().getSizeInBytes();
    }

    private static void mergeBindingFailure(Map<String, PlcResponseItem<PlcValue>> out,
                                            S7ReadChunk.Binding b, PlcResponseCode code) {
        // For split fragments a single failure poisons the whole tag; that matches the
        // legacy behaviour and avoids returning a half-filled byte array as success.
        out.put(b.tagName(), new DefaultPlcResponseItem<>(code, null));
    }


    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        // Same as for reads: a tag the builder rejected is reported per tag, not by failing the
        // whole request.
        Map<String, PlcResponseCode> rejectedTags = new LinkedHashMap<>();
        for (String tagName : writeRequest.getTagNames()) {
            PlcResponseCode requestCode = writeRequest.getTagResponseCode(tagName);
            if (requestCode != PlcResponseCode.OK) {
                rejectedTags.put(tagName, requestCode);
                continue;
            }
            org.apache.plc4x.java.api.model.PlcTag tag = writeRequest.getTag(tagName);
            if (!(tag instanceof S7Tag)) {
                return CompletableFuture.failedFuture(new PlcProtocolException(
                    "Unsupported tag type " + tag.getClass().getName()));
            }
        }
        if (rejectedTags.size() == writeRequest.getTagNames().size()) {
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, rejectedTags));
        }

        List<S7WriteChunk> chunks;
        try {
            chunks = optimizer.splitWriteRequest(writeRequest, driverContext);
        } catch (PlcRuntimeException e) {
            Map<String, PlcResponseCode> codes = new LinkedHashMap<>();
            for (String t : writeRequest.getTagNames()) {
                codes.put(t, rejectedTags.getOrDefault(t, PlcResponseCode.INVALID_DATA));
            }
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, codes));
        }
        if (chunks.isEmpty()) {
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, new LinkedHashMap<>()));
        }

        Map<String, PlcResponseCode> finalCodes = new LinkedHashMap<>();
        for (String t : writeRequest.getTagNames()) {
            finalCodes.put(t, rejectedTags.get(t));
        }
        CompletableFuture<Void> chain = CompletableFuture.completedFuture(null);
        for (S7WriteChunk chunk : chunks) {
            chain = chain.thenCompose(v -> sendWriteChunk(chunk, writeRequest, finalCodes));
        }
        return chain.thenApply(v -> {
            for (Map.Entry<String, PlcResponseCode> e : finalCodes.entrySet()) {
                if (e.getValue() == null) e.setValue(PlcResponseCode.INTERNAL_ERROR);
            }
            return new DefaultPlcWriteResponse(writeRequest, finalCodes);
        });
    }

    private CompletableFuture<Void> sendWriteChunk(S7WriteChunk chunk, PlcWriteRequest writeRequest,
                                                   Map<String, PlcResponseCode> out) {
        List<S7VarRequestParameterItem> parameterItems = new ArrayList<>(chunk.tagNames().size());
        List<S7VarPayloadDataItem> payloadItems = new ArrayList<>(chunk.tagNames().size());
        for (String tagName : chunk.tagNames()) {
            S7Tag s7Tag = (S7Tag) writeRequest.getTag(tagName);
            parameterItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(s7Tag)));
            S7VarPayloadDataItem item = serializePlcValue(s7Tag, writeRequest.getPlcValue(tagName));
            if (item == null) {
                out.put(tagName, PlcResponseCode.INVALID_DATA);
                // Skip this tag in the actual request — would otherwise corrupt the layout.
                continue;
            }
            payloadItems.add(item);
        }
        if (parameterItems.size() != payloadItems.size()) {
            // At least one serialization failed; remaining tags can still go through.
            // Re-build the parallel lists trimmed to successful ones.
            parameterItems = new ArrayList<>();
            payloadItems = new ArrayList<>();
            for (String tagName : chunk.tagNames()) {
                if (out.get(tagName) == PlcResponseCode.INVALID_DATA) continue;
                S7Tag s7Tag = (S7Tag) writeRequest.getTag(tagName);
                parameterItems.add(new S7VarRequestParameterItemAddress(encodeS7Address(s7Tag)));
                payloadItems.add(serializePlcValue(s7Tag, writeRequest.getPlcValue(tagName)));
            }
        }
        if (parameterItems.isEmpty()) {
            return CompletableFuture.completedFuture(null);
        }
        S7Message request = new S7MessageRequest(getTpduId(),
            new S7ParameterWriteVarRequest(parameterItems),
            new S7PayloadWriteVarRequest(payloadItems));
        final List<String> sentTagNames = new ArrayList<>();
        for (String tagName : chunk.tagNames()) {
            if (out.get(tagName) != PlcResponseCode.INVALID_DATA) {
                sentTagNames.add(tagName);
            }
        }
        return executeThrottled(() -> sendInternal(request))
            .handle((response, error) -> {
                if (error != null) {
                    for (String n : sentTagNames) out.put(n, PlcResponseCode.INTERNAL_ERROR);
                    return null;
                }
                short errorClass = 0, errorCode = 0;
                if (response instanceof S7MessageResponseData md) {
                    errorClass = md.getErrorClass();
                    errorCode = md.getErrorCode();
                } else if (response instanceof S7MessageResponse mr) {
                    errorClass = mr.getErrorClass();
                    errorCode = mr.getErrorCode();
                }
                if (errorClass != 0 || errorCode != 0) {
                    PlcResponseCode code = mapPlcErrorCode(errorClass, errorCode);
                    for (String n : sentTagNames) out.put(n, code);
                    return null;
                }
                if (!(response.getPayload() instanceof S7PayloadWriteVarResponse payload)) {
                    for (String n : sentTagNames) out.put(n, PlcResponseCode.INTERNAL_ERROR);
                    return null;
                }
                List<S7VarPayloadStatusItem> items = payload.getItems();
                for (int i = 0; i < sentTagNames.size(); i++) {
                    String n = sentTagNames.get(i);
                    if (i >= items.size()) {
                        out.put(n, PlcResponseCode.INTERNAL_ERROR);
                    } else {
                        out.put(n, mapDataTransportError(items.get(i).getReturnCode()));
                    }
                }
                return null;
            });
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding/decoding helpers (ported from s7-light S7ProtocolLogic)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private S7Address encodeS7Address(S7Tag s7Tag) {
        TransportSize transportSize = s7Tag.getDataType();
        int numElements = s7Tag.getNumberOfElements();
        if (transportSize == TransportSize.STRING) {
            transportSize = TransportSize.CHAR;
            int stringLength = (s7Tag instanceof S7StringFixedLengthTag f) ? f.getStringLength() : 254;
            numElements = numElements * (stringLength + 2);
        } else if (transportSize == TransportSize.WSTRING) {
            transportSize = TransportSize.CHAR;
            int stringLength = (s7Tag instanceof S7StringFixedLengthTag f) ? f.getStringLength() : 254;
            numElements = numElements * (stringLength + 2) * 2;
        } else if ((transportSize == TransportSize.BOOL) && (s7Tag.getNumberOfElements() > 1)) {
            numElements = (s7Tag.getNumberOfElements() + 7) / 8;
            transportSize = TransportSize.BYTE;
        }
        if (transportSize.getCode() == 0x00) {
            numElements = numElements * transportSize.getSizeInBytes();
            transportSize = TransportSize.BYTE;
        }
        return new S7AddressAny(transportSize, numElements, s7Tag.getBlockNumber(),
            s7Tag.getMemoryArea(), s7Tag.getByteOffset(), s7Tag.getBitOffset());
    }

    private S7VarPayloadDataItem serializePlcValue(S7Tag tag, PlcValue plcValue) {
        try {
            DataTransportSize transportSize = tag.getDataType().getDataTransportSize();
            int stringLength = (tag instanceof S7StringFixedLengthTag f) ? f.getStringLength() : 254;
            ByteBuffer byteBuffer = null;
            if (tag.getDataType() == TransportSize.BYTE && tag.getNumberOfElements() > 1) {
                byteBuffer = ByteBuffer.allocate(tag.getNumberOfElements());
                byteBuffer.put(plcValue.getRaw());
            } else if (tag.getDataType() == TransportSize.BOOL && tag.getNumberOfElements() > 1) {
                if (!(plcValue instanceof PlcList plcList)) {
                    throw new PlcRuntimeException("Expected PlcList for multi-element BOOL");
                }
                int numBytes = (tag.getNumberOfElements() + 7) / 8;
                byteBuffer = ByteBuffer.allocate(numBytes);
                for (int i = 0; i < tag.getNumberOfElements(); i++) {
                    if (plcList.getIndex(i) instanceof PlcBOOL b && b.getBoolean()) {
                        int curByte = i / 8;
                        int curBit = i % 8;
                        byteBuffer.put(curByte, (byte) (1 << curBit | byteBuffer.get(curByte)));
                    }
                }
                transportSize = DataTransportSize.BYTE_WORD_DWORD;
            } else {
                for (int i = 0; i < tag.getNumberOfElements(); i++) {
                    int lengthInBits = DataItem.getLengthInBits(plcValue.getIndex(i),
                        tag.getDataType().getDataProtocolId(), driverContext.getControllerType(), stringLength);
                    if (tag.getDataType() == TransportSize.STRING) {
                        lengthInBits = Math.min(lengthInBits, (stringLength * 8) + 16);
                    } else if (tag.getDataType() == TransportSize.WSTRING) {
                        lengthInBits = Math.min(lengthInBits, (stringLength * 16) + 32);
                    } else if (tag.getDataType() == TransportSize.S5TIME) {
                        lengthInBits = lengthInBits * 8;
                    }
                    WriteBufferByteBased writeBuffer = new WriteBufferByteBased(
                        new byte[(int) Math.ceil(((float) lengthInBits) / 8.0f)],
                        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                        WithOption.WithSignedIntegerEncoding("twos-complement"),
                        WithOption.WithFloatEncoding("IEEE754"),
                        WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
                    DataItem.staticSerialize(writeBuffer, plcValue.getIndex(i),
                        tag.getDataType().getDataProtocolId(), driverContext.getControllerType(), stringLength);
                    if (byteBuffer == null) {
                        byteBuffer = ByteBuffer.allocate(writeBuffer.getBytes().length * tag.getNumberOfElements());
                    }
                    byteBuffer.put(writeBuffer.getBytes());
                }
            }
            if (byteBuffer != null) {
                return new S7VarPayloadDataItem(DataTransportErrorCode.OK, transportSize, byteBuffer.array());
            }
        } catch (BufferException e) {
            LOGGER.warn("Error serializing tag {}", tag.getDataType().name(), e);
        }
        return null;
    }

    private PlcValue parsePlcValue(S7Tag tag, byte[] data) throws BufferException {
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
        int stringLength = (tag instanceof S7StringFixedLengthTag f) ? f.getStringLength() : 254;
        if (tag.getNumberOfElements() == 1) {
            return DataItem.staticParse(readBuffer, tag.getDataType().getDataProtocolId(),
                driverContext.getControllerType(), stringLength);
        }
        if (tag.getDataType() == TransportSize.BYTE) {
            return new PlcRawByteArray(data);
        }
        PlcValue[] resultItems = IntStream.range(0, tag.getNumberOfElements()).mapToObj(i -> {
            try {
                return DataItem.staticParse(readBuffer, tag.getDataType().getDataProtocolId(),
                    driverContext.getControllerType(), stringLength);
            } catch (BufferException e) {
                LOGGER.warn("Error parsing element {} of tag {}", i, tag.getDataType().name(), e);
                return null;
            }
        }).toArray(PlcValue[]::new);
        return DefaultPlcValueHandler.of(tag, resultItems);
    }

    private static PlcResponseCode mapDataTransportError(DataTransportErrorCode code) {
        if (code == null) return PlcResponseCode.INTERNAL_ERROR;
        return switch (code) {
            case OK -> PlcResponseCode.OK;
            case NOT_FOUND, INVALID_ADDRESS -> PlcResponseCode.INVALID_ADDRESS;
            case DATA_TYPE_NOT_SUPPORTED -> PlcResponseCode.INVALID_DATATYPE;
            case ACCESS_DENIED -> PlcResponseCode.ACCESS_DENIED;
            default -> PlcResponseCode.INTERNAL_ERROR;
        };
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Browse
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        return onBrowseWithInterceptor(browseRequest, (queryName, query, item) -> true);
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowseWithInterceptor(PlcBrowseRequest browseRequest,
                                                                           PlcBrowseRequestInterceptor interceptor) {
        // Browse rides UserData services. If the SZL probe at connect time said the device
        // doesn't speak UserData, return UNSUPPORTED per query rather than throwing.
        if (!driverContext.isUserDataServicesSupported()) {
            Map<String, PlcResponseCode> codes = new LinkedHashMap<>();
            Map<String, List<PlcBrowseItem>> values = new LinkedHashMap<>();
            for (String queryName : browseRequest.getQueryNames()) {
                codes.put(queryName, PlcResponseCode.UNSUPPORTED);
                values.put(queryName, Collections.emptyList());
            }
            return CompletableFuture.completedFuture(new DefaultPlcBrowseResponse(browseRequest, codes, values));
        }

        boolean subscribable = getMetadata().isSubscribeSupported();
        List<PlcBrowseItem> items = new ArrayList<>(S7BrowseService.staticAreas(subscribable));

        // Issue a "list of blocks of type DB" request. Treat any failure (timeout, PLC
        // doesn't honor block functions) as "no DBs" — static areas still go in the response.
        S7Message request = S7BrowseService.buildListBlocksOfTypeRequest(
            getTpduId(), S7BrowseService.BlockType.DB);
        return sendInternal(request)
            .orTimeout(driverContext.getReadTimeout(), TimeUnit.MILLISECONDS)
            .handle((response, error) -> {
                if (error == null) {
                    List<Integer> blockNumbers = S7BrowseService.parseListBlocksOfTypeResponse(response);
                    LOGGER.debug("Block-list browse found {} DB(s)", blockNumbers.size());
                    items.addAll(S7BrowseService.dataBlocksFromNumbers(blockNumbers, subscribable));
                } else {
                    LOGGER.debug("Block-list browse failed; returning static areas only: {}",
                        error.getMessage());
                }
                Map<String, PlcResponseCode> codes = new LinkedHashMap<>();
                Map<String, List<PlcBrowseItem>> values = new LinkedHashMap<>();
                for (String queryName : browseRequest.getQueryNames()) {
                    PlcQuery query = browseRequest.getQuery(queryName);
                    List<PlcBrowseItem> filtered = new ArrayList<>(items.size());
                    for (PlcBrowseItem item : items) {
                        if (interceptor.intercept(queryName, query, item)) {
                            filtered.add(item);
                        }
                    }
                    codes.put(queryName, PlcResponseCode.OK);
                    values.put(queryName, filtered);
                }
                return new DefaultPlcBrowseResponse(browseRequest, codes, values);
            });
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Alarm subscription (S7Comm UserData MsgSubscription, group=0x04, subfn=0x02)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        if (!driverContext.isUserDataServicesSupported()) {
            return CompletableFuture.completedFuture(
                buildSubscriptionResponse(subscriptionRequest, PlcResponseCode.UNSUPPORTED, null));
        }

        // Three flavours of subscription are supported:
        //   - S7AlarmTag (PUSH)  + EVENT  → alarm push (S7Comm UserData group=0x04, subfn 0x02)
        //   - S7AlarmTag (QUERY) + EVENT  → one-shot AlarmQuery (group=0x04, subfn 0x13)
        //   - S7Tag              + CYCLIC → cyclic value push (group=0x02)
        // Any other combination is rejected as INVALID_ADDRESS.
        List<String> alarmTagNames = new ArrayList<>();
        List<String> queryTagNames = new ArrayList<>();
        List<String> cyclicTagNames = new ArrayList<>();
        for (String name : subscriptionRequest.getTagNames()) {
            PlcSubscriptionTag subTag = subscriptionRequest.getTag(name);
            PlcTag inner = unwrap(subTag);
            PlcSubscriptionType type = subTag.getPlcSubscriptionType();
            if (inner instanceof S7AlarmTag at
                && type == org.apache.plc4x.java.api.types.PlcSubscriptionType.EVENT) {
                if (at.getKind() == S7AlarmTag.Kind.QUERY) {
                    queryTagNames.add(name);
                } else {
                    alarmTagNames.add(name);
                }
            } else if (inner instanceof S7Tag && (type == PlcSubscriptionType.CYCLIC)) {
                cyclicTagNames.add(name);
            } else {
                return CompletableFuture.completedFuture(
                    buildSubscriptionResponse(subscriptionRequest, PlcResponseCode.INVALID_ADDRESS, null));
            }
        }

        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();

        CompletableFuture<Void> alarmStage = alarmTagNames.isEmpty()
            ? CompletableFuture.completedFuture(null)
            : ensureAlarmSubscriptionActive().thenAccept(active -> {
                for (String name : alarmTagNames) {
                    if (active) {
                        S7AlarmSubscriptionHandle handle = new S7AlarmSubscriptionHandle(this, name);
                        alarmHandles.add(handle);
                        values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
                    } else {
                        values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null));
                    }
                }
            });

        CompletableFuture<Void> cyclicStage = cyclicTagNames.isEmpty()
            ? CompletableFuture.completedFuture(null)
            : subscribeCyclic(subscriptionRequest, cyclicTagNames, values);

        CompletableFuture<Void> queryStage = queryTagNames.isEmpty()
            ? CompletableFuture.completedFuture(null)
            : runAlarmQueries(subscriptionRequest, queryTagNames, values);

        return CompletableFuture.allOf(alarmStage, cyclicStage, queryStage)
            .thenApply(v -> new DefaultPlcSubscriptionResponse(subscriptionRequest, values));
    }

    /**
     * Run a one-shot AlarmQuery for each requested tag, in parallel. Each tag carries its
     * own QueryType (ALARM_S vs ALARM_8), so we issue one wire request per tag rather than
     * trying to bundle them.
     */
    private CompletableFuture<Void> runAlarmQueries(PlcSubscriptionRequest req, List<String> tagNames,
                                                    Map<String, PlcResponseItem<PlcSubscriptionHandle>> values) {
        List<CompletableFuture<Void>> stages = new ArrayList<>(tagNames.size());
        for (String name : tagNames) {
            S7AlarmTag tag = (S7AlarmTag) unwrap(req.getTag(name));
            S7Message request = S7AlarmSubscriptionService.buildAlarmQueryRequest(
                getTpduId(), tag.getQueryType());
            stages.add(sendInternal(request)
                .orTimeout(driverContext.getReadTimeout(), TimeUnit.MILLISECONDS)
                .handle((response, error) -> {
                    if (error != null || response == null) {
                        LOGGER.debug("Alarm query failed for {}: {}", name,
                            error == null ? "no response" : error.getMessage());
                        values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null));
                        return null;
                    }
                    byte[] payload = S7AlarmSubscriptionService.parseAlarmQueryResponse(response);
                    if (payload == null) {
                        LOGGER.debug("Alarm query rejected by PLC for {}", name);
                        values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null));
                        return null;
                    }
                    S7AlarmQueryHandle handle = new S7AlarmQueryHandle(this, name, payload);
                    values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
                    return null;
                }));
        }
        return CompletableFuture.allOf(stages.toArray(CompletableFuture[]::new));
    }

    /**
     * Issue cyclic-subscribe requests grouped by interval (one wire request per distinct
     * interval) and populate {@code values} with the per-tag response items. Tags that share
     * an interval get bundled into a single subscription with the PLC, sharing one jobId.
     */
    private CompletableFuture<Void> subscribeCyclic(PlcSubscriptionRequest req, List<String> tagNames,
                                                    Map<String, PlcResponseItem<PlcSubscriptionHandle>> values) {
        // Group tags by their (TimeBase, factor) pair so callers requesting the same cadence
        // ride one PLC subscription rather than one-per-tag.
        Map<S7CyclicSubscriptionService.IntervalSpec, List<String>> byInterval = new LinkedHashMap<>();
        for (String name : tagNames) {
            PlcSubscriptionTag subTag = req.getTag(name);
            S7CyclicSubscriptionService.IntervalSpec interval =
                S7CyclicSubscriptionService.pickInterval(subTag.getDuration().orElse(null));
            byInterval.computeIfAbsent(interval, k -> new ArrayList<>()).add(name);
        }

        List<CompletableFuture<Void>> stages = new ArrayList<>(byInterval.size());
        for (Map.Entry<S7CyclicSubscriptionService.IntervalSpec, List<String>> entry : byInterval.entrySet()) {
            stages.add(subscribeCyclicGroup(req, entry.getKey(), entry.getValue(), values));
        }
        return CompletableFuture.allOf(stages.toArray(CompletableFuture[]::new));
    }

    private CompletableFuture<Void> subscribeCyclicGroup(PlcSubscriptionRequest req,
                                                         S7CyclicSubscriptionService.IntervalSpec interval,
                                                         List<String> tagNames,
                                                         Map<String, PlcResponseItem<PlcSubscriptionHandle>> values) {
        List<S7Tag> tags = new ArrayList<>(tagNames.size());
        for (String name : tagNames) {
            tags.add((S7Tag) unwrap(req.getTag(name)));
        }
        S7Message request = S7CyclicSubscriptionService.buildSubscribeRequest(getTpduId(), tags, interval);
        return sendInternal(request)
            .orTimeout(driverContext.getReadTimeout(), TimeUnit.MILLISECONDS)
            .handle((response, error) -> {
                if (error != null || response == null) {
                    LOGGER.debug("Cyclic subscribe failed: {}",
                        error == null ? "no response" : error.getMessage());
                    for (String name : tagNames) {
                        values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null));
                    }
                    return null;
                }
                Short jobId = S7CyclicSubscriptionService.parseSubscribeResponse(response);
                if (jobId == null) {
                    LOGGER.debug("Cyclic subscribe rejected by PLC");
                    for (String name : tagNames) {
                        values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null));
                    }
                    return null;
                }
                installCyclicIndicationHandlers();
                List<S7CyclicSubscriptionHandle> handles = new ArrayList<>(tagNames.size());
                for (int i = 0; i < tagNames.size(); i++) {
                    String name = tagNames.get(i);
                    S7CyclicSubscriptionHandle handle =
                        new S7CyclicSubscriptionHandle(this, name, tags.get(i), jobId, i);
                    handles.add(handle);
                    values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
                }
                cyclicHandlesByJob.put(jobId, handles);
                LOGGER.debug("Cyclic subscribe active: jobId={}, tags={}, intervalMs={}",
                    jobId, tagNames, interval.toMillis());
                return null;
            });
    }

    private void installCyclicIndicationHandlers() {
        synchronized (cyclicDispatcherHooks) {
            if (!cyclicDispatcherHooks.isEmpty()) {
                return;
            }
            cyclicDispatcherHooks.add(pushDispatcher.register(
                S7CyclicSubscriptionService.CPU_FUNCTION_GROUP_CYCLIC,
                S7CyclicSubscriptionService.CPU_FUNCTION_TYPE_PUSH,
                S7CyclicSubscriptionService.CYCLIC_PUSH_SUBFUNCTION,
                this::dispatchCyclicPush));
        }
    }

    private void dispatchCyclicPush(S7Message message) {
        Short jobId = S7CyclicSubscriptionService.getJobId(message);
        if (jobId == null) return;
        List<S7CyclicSubscriptionHandle> handles = cyclicHandlesByJob.get(jobId);
        if (handles == null || handles.isEmpty()) return;
        List<byte[]> items = S7CyclicSubscriptionService.parsePushItems(message);
        if (items == null) return;

        // Decode each item using its corresponding tag's data type and emit one event per
        // consumer registration whose handle set covers this jobId.
        Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
        for (S7CyclicSubscriptionHandle handle : handles) {
            int idx = handle.getItemIndex();
            if (idx >= items.size()) {
                values.put(handle.getTagName(), new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null));
                continue;
            }
            try {
                PlcValue value = parsePlcValue(handle.getTag(), items.get(idx));
                values.put(handle.getTagName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK, value));
            } catch (Exception e) {
                LOGGER.debug("Cyclic push decode failed for tag {}: {}", handle.getTagName(), e.getMessage());
                values.put(handle.getTagName(), new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
            }
        }
        DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(java.time.Instant.now(), values);
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            for (PlcSubscriptionHandle h : entry.getKey().getSubscriptionHandles()) {
                if (h instanceof S7CyclicSubscriptionHandle ch && ch.getJobId() == jobId) {
                    entry.getValue().accept(event);
                    break;    // one event per registration even if multiple handles match
                }
            }
        }
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        Set<Short> cyclicJobsToCancel = new java.util.HashSet<>();
        for (PlcSubscriptionHandle h : unsubscriptionRequest.getSubscriptionHandles()) {
            if (h instanceof S7AlarmSubscriptionHandle ah) {
                alarmHandles.remove(ah);
            } else if (h instanceof S7CyclicSubscriptionHandle ch) {
                short jobId = ch.getJobId();
                List<S7CyclicSubscriptionHandle> remaining = cyclicHandlesByJob.get(jobId);
                if (remaining != null) {
                    remaining.remove(ch);
                    if (remaining.isEmpty()) {
                        cyclicHandlesByJob.remove(jobId);
                        cyclicJobsToCancel.add(jobId);
                    }
                }
            }
        }
        // Wait for both deactivations (when triggered) before completing the response so the
        // caller's subsequent connection close doesn't race ahead of the cancel requests.
        CompletableFuture<Void> alarmStage = alarmHandles.isEmpty()
            ? deactivateAlarmSubscription()
            : CompletableFuture.completedFuture(null);
        CompletableFuture<Void> cyclicStage = cancelCyclicJobs(cyclicJobsToCancel);
        return CompletableFuture.allOf(alarmStage, cyclicStage)
            .handle((v, t) -> new DefaultPlcUnsubscriptionResponse(unsubscriptionRequest));
    }

    private CompletableFuture<Void> cancelCyclicJobs(Set<Short> jobIds) {
        if (jobIds.isEmpty()) {
            return CompletableFuture.completedFuture(null);
        }
        List<CompletableFuture<Void>> stages = new ArrayList<>(jobIds.size());
        for (Short jobId : jobIds) {
            try {
                stages.add(sendInternal(S7CyclicSubscriptionService.buildUnsubscribeRequest(getTpduId(), jobId))
                    .orTimeout(driverContext.getReadTimeout(), TimeUnit.MILLISECONDS)
                    .handle((response, error) -> {
                        if (error != null) {
                            LOGGER.debug("Cyclic cancel send failed for jobId {} (ignoring): {}",
                                jobId, error.getMessage());
                        }
                        return null;
                    }));
            } catch (Exception e) {
                // Codec already shut down — nothing to send.
                stages.add(CompletableFuture.completedFuture(null));
            }
        }
        // If no other cyclic subscription is active, also tear down the push handler.
        if (cyclicHandlesByJob.isEmpty()) {
            synchronized (cyclicDispatcherHooks) {
                for (java.io.Closeable hook : cyclicDispatcherHooks) {
                    try { hook.close(); } catch (Exception ignored) { /* local-only impl */ }
                }
                cyclicDispatcherHooks.clear();
            }
        }
        return CompletableFuture.allOf(stages.toArray(CompletableFuture[]::new));
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer,
                                                         Collection<PlcSubscriptionHandle> handles) {
        DefaultPlcConsumerRegistration registration = new DefaultPlcConsumerRegistration(
            this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(registration, consumer);
        // Query handles carry their result inline (it was fetched at subscribe time). Deliver
        // it as a single PlcSubscriptionEvent now and the handle goes inert thereafter.
        Map<String, PlcResponseItem<PlcValue>> queryValues = null;
        for (PlcSubscriptionHandle handle : handles) {
            if (handle instanceof S7AlarmQueryHandle qh) {
                if (queryValues == null) {
                    queryValues = new LinkedHashMap<>();
                }
                queryValues.put(qh.getTagName(),
                    new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcRawByteArray(qh.getPayload())));
            }
        }
        if (queryValues != null) {
            consumer.accept(new DefaultPlcSubscriptionEvent(java.time.Instant.now(), queryValues));
        }
        return registration;
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        if (registration instanceof DefaultPlcConsumerRegistration r) {
            consumers.remove(r);
        }
    }

    private CompletableFuture<Boolean> ensureAlarmSubscriptionActive() {
        synchronized (alarmDispatcherHooks) {
            if (!alarmDispatcherHooks.isEmpty()) {
                return CompletableFuture.completedFuture(Boolean.TRUE);
            }
        }
        S7Message request = S7AlarmSubscriptionService.buildSubscribeRequest(
            getTpduId(), driverContext.getControllerType());
        return sendInternal(request)
            .orTimeout(driverContext.getReadTimeout(), TimeUnit.MILLISECONDS)
            .handle((response, error) -> {
                if (error != null) {
                    LOGGER.debug("Alarm subscription failed: {}", error.getMessage(), error);
                    return Boolean.FALSE;
                }
                if (!S7AlarmSubscriptionService.parseSubscribeResponse(response)) {
                    LOGGER.debug("Alarm subscription rejected by PLC");
                    return Boolean.FALSE;
                }
                installAlarmIndicationHandlers();
                LOGGER.debug("Alarm subscription active");
                return Boolean.TRUE;
            });
    }

    private void installAlarmIndicationHandlers() {
        synchronized (alarmDispatcherHooks) {
            if (!alarmDispatcherHooks.isEmpty()) {
                return;
            }
            for (Integer subfn : S7AlarmSubscriptionService.ALARM_INDICATION_SUBFUNCTIONS) {
                alarmDispatcherHooks.add(pushDispatcher.register(
                    S7AlarmSubscriptionService.CPU_FUNCTION_GROUP_CPU,
                    S7AlarmSubscriptionService.CPU_FUNCTION_TYPE_PUSH,
                    subfn,
                    this::dispatchAlarmIndication));
            }
        }
    }

    /**
     * Tear down alarm subscription state and tell the PLC to stop pushing. Returns a future
     * that completes once the ALARM_S_ABORT request has either been answered, timed out, or
     * failed — so {@code onUnsubscribe} can await it before returning to the caller. Without
     * waiting, the caller's subsequent connection close races ahead of the request and the
     * cancel never makes it onto the wire.
     */
    private CompletableFuture<Void> deactivateAlarmSubscription() {
        synchronized (alarmDispatcherHooks) {
            for (java.io.Closeable hook : alarmDispatcherHooks) {
                try {
                    hook.close();
                } catch (Exception ignored) {
                    // hook.close() is local-only and never throws in our impl.
                }
            }
            alarmDispatcherHooks.clear();
        }
        try {
            return sendInternal(S7AlarmSubscriptionService.buildUnsubscribeRequest(
                    getTpduId(), driverContext.getControllerType()))
                .orTimeout(driverContext.getReadTimeout(), TimeUnit.MILLISECONDS)
                .handle((response, error) -> {
                    if (error != null) {
                        LOGGER.debug("Alarm cancel send failed (ignoring): {}", error.getMessage());
                    }
                    return null;
                });
        } catch (Exception e) {
            // Codec already shut down — nothing to send. Treat as already-deactivated.
            return CompletableFuture.completedFuture(null);
        }
    }

    private void dispatchAlarmIndication(S7Message message) {
        PlcValue payload = S7AlarmSubscriptionService.parseAlarmIndication(message);
        if (payload == null) {
            return;
        }
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            for (PlcSubscriptionHandle handle : entry.getKey().getSubscriptionHandles()) {
                if (handle instanceof S7AlarmSubscriptionHandle ah) {
                    Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
                    values.put(ah.getTagName(), new DefaultPlcResponseItem<>(PlcResponseCode.OK, payload));
                    entry.getValue().accept(new DefaultPlcSubscriptionEvent(java.time.Instant.now(), values));
                }
            }
        }
    }

    private static PlcTag unwrap(PlcSubscriptionTag tag) {
        return tag == null ? null : tag.getTag();
    }

    private DefaultPlcSubscriptionResponse buildSubscriptionResponse(PlcSubscriptionRequest request,
                                                                     PlcResponseCode code,
                                                                     PlcSubscriptionHandle handle) {
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();
        for (String name : request.getTagNames()) {
            values.put(name, new DefaultPlcResponseItem<>(code, handle));
        }
        return new DefaultPlcSubscriptionResponse(request, values);
    }

    private static PlcResponseCode mapPlcErrorCode(short errorClass, short errorCode) {
        // Ported error mapping from s7-light: 129/4 means PUT/GET disabled.
        if (errorClass == 0x81 && errorCode == 4) return PlcResponseCode.ACCESS_DENIED;
        // An S7-300 reports the same refusal as 0x83/0x04 rather than 0x81/0x04 (GH-599). The
        // generic reading of class 0x83 is "no resources available", so only this exact pairing
        // is treated as a refusal - anything else in that class keeps falling through.
        if (errorClass == 0x83 && errorCode == 4) return PlcResponseCode.ACCESS_DENIED;
        if (errorClass == 0x85) return PlcResponseCode.ACCESS_DENIED;
        return PlcResponseCode.INTERNAL_ERROR;
    }

}
