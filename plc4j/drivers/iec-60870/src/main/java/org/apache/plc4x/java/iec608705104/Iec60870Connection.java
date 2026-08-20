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
package org.apache.plc4x.java.iec608705104;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.iec608705104.configuration.Iec608705014Configuration;
import org.apache.plc4x.java.iec608705104.messages.Iec608705104PlcSubscriptionEvent;
import org.apache.plc4x.java.iec608705104.model.Iec608705104SubscriptionHandle;
import org.apache.plc4x.java.iec608705104.protocol.Iec608705104TagParser;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104TagHandler;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.*;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.*;
import java.util.concurrent.*;
import java.util.function.Consumer;
import org.apache.plc4x.java.iec608705104.readwrite.APDU;
import org.apache.plc4x.java.iec608705104.readwrite.APDUIFormat;
import org.apache.plc4x.java.iec608705104.readwrite.APDUSFormat;
import org.apache.plc4x.java.iec608705104.readwrite.APDUUFormatStartDataTransferActivation;
import org.apache.plc4x.java.iec608705104.readwrite.APDUUFormatStartDataTransferConfirmation;
import org.apache.plc4x.java.iec608705104.readwrite.APDUUFormatTestFrameActivation;
import org.apache.plc4x.java.iec608705104.readwrite.APDUUFormatTestFrameConfirmation;
import org.apache.plc4x.java.iec608705104.readwrite.ASDU;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObject;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithSevenByteTime;
import org.apache.plc4x.java.iec608705104.readwrite.InformationObjectWithTreeByteTime;
import org.apache.plc4x.java.iec608705104.readwrite.SevenOctetBinaryTime;
import org.apache.plc4x.java.iec608705104.readwrite.ThreeOctetBinaryTime;

/**
 * IEC 60870-5-104 TCP connection. Direct port of the legacy
 * {@code Iec608705104Protocol} to the SPI3 {@link ConnectionBase} model.
 *
 * <p>The protocol is push-driven once a session is established:
 * <ol>
 *   <li>The driver opens a TCP connection and exchanges a {@code TestFrame
 *       Activation/Confirmation} handshake to verify the link.</li>
 *   <li>It then sends a {@code StartDataTransfer Activation} and waits for
 *       the {@code Confirmation}, after which the controlled station is
 *       free to emit unsolicited I-format ASDUs.</li>
 *   <li>Each received ASDU is decoded into one or more {@link PlcValue}s
 *       and dispatched to all registered subscription consumers whose
 *       handle matches.</li>
 * </ol>
 *
 * <p>No reads or writes are exposed — only subscriptions.</p>
 */
public class Iec60870Connection extends ConnectionBase<Iec608705014Configuration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Iec60870Connection.class);

    /** Send an S-format acknowledgement every N unconfirmed I-format frames. */
    private static final int ACK_EVERY_N_FRAMES = 8;

    private Iec60870MessageCodec messageCodec;

    /** Single-flight pending response future for the handshake. */
    private final BlockingQueue<CompletableFuture<APDU>> pendingHandshakeResponses = new LinkedBlockingDeque<>();

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    private volatile boolean handshakeComplete = false;
    private int unconfirmedPackets;

    public Iec60870Connection(Iec608705014Configuration configuration,
                              TransportInstance<?> transportInstance,
                              AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new Iec608705104TagHandler();
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
        messageCodec = new Iec60870MessageCodec(transportInstance, this::handleIncomingMessage);

        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming IEC 60870-5-104 data", e);
            }
        });

        try {
            doHandshake().get(getConfiguration().getRequestTimeout() * 5L, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("Error during IEC 60870-5-104 connect handshake", e);
        }

        handshakeComplete = true;
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT, "IEC 60870-5-104 connection established");
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    private CompletableFuture<Void> doHandshake() {
        // 1) TestFrame Activation → Confirmation
        APDUUFormatTestFrameActivation testFrameActivation = new APDUUFormatTestFrameActivation(0x43);
        return sendExpectingReply(testFrameActivation)
            .thenCompose(testFrameResponse -> {
                if (!(testFrameResponse instanceof APDUUFormatTestFrameConfirmation)) {
                    return CompletableFuture.failedFuture(new PlcRuntimeException(
                        "Expected APDUUFormatTestFrameConfirmation, got " + testFrameResponse.getClass().getSimpleName()));
                }
                // 2) StartDataTransfer Activation → Confirmation
                APDUUFormatStartDataTransferActivation startDataTransferActivation =
                    new APDUUFormatStartDataTransferActivation(0x07);
                return sendExpectingReply(startDataTransferActivation);
            })
            .thenAccept(startDataTransferResponse -> {
                if (!(startDataTransferResponse instanceof APDUUFormatStartDataTransferConfirmation)) {
                    throw new PlcRuntimeException(
                        "Expected APDUUFormatStartDataTransferConfirmation, got " + startDataTransferResponse.getClass().getSimpleName());
                }
            });
    }

    private CompletableFuture<APDU> sendExpectingReply(APDU request) {
        CompletableFuture<APDU> responseFuture = new CompletableFuture<>();
        pendingHandshakeResponses.offer(responseFuture);
        try {
            if (auditLog.isEnabled()) {
                auditLog.write(AuditLogEventType.OUTGOING_MESSAGE,
                    "Sending " + request.getClass().getSimpleName());
            }
            messageCodec.send(request);
        } catch (MessageCodecException e) {
            pendingHandshakeResponses.remove(responseFuture);
            responseFuture.completeExceptionally(new PlcRuntimeException("Failed to send IEC 60870-5-104 request", e));
            return responseFuture;
        }
        responseFuture.orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
            .whenComplete((r, e) -> {
                if (e instanceof TimeoutException) {
                    pendingHandshakeResponses.remove(responseFuture);
                }
            });
        return responseFuture;
    }

    @Override
    public void close() throws Exception {
        handshakeComplete = false;
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        pendingHandshakeResponses.forEach(f ->
            f.completeExceptionally(new PlcRuntimeException("Connection closed")));
        pendingHandshakeResponses.clear();
        consumers.clear();
        super.close();
        LOGGER.info("IEC 60870-5-104 connection closed");
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        handshakeComplete = false;
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");
        PlcRuntimeException exception = new PlcRuntimeException(
            cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause);
        pendingHandshakeResponses.forEach(f -> f.completeExceptionally(exception));
        pendingHandshakeResponses.clear();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming messages
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingMessage(APDU apdu) {
        // While the handshake is in progress, every reply is a U-format
        // control frame answering our outstanding request — dispatch to the
        // pending future.
        if (!handshakeComplete) {
            CompletableFuture<APDU> future = pendingHandshakeResponses.poll();
            if (future != null) {
                future.complete(apdu);
                return;
            }
            LOGGER.warn("Received pre-handshake APDU with no pending request: {}", apdu.getClass().getSimpleName());
            return;
        }

        // Post-handshake the controlled station may probe us with a test
        // frame at any time. Echo the confirmation.
        if (apdu instanceof APDUUFormatTestFrameActivation) {
            try {
                messageCodec.send(new APDUUFormatTestFrameConfirmation(0x83));
            } catch (MessageCodecException e) {
                LOGGER.error("Failed to send test-frame confirmation", e);
            }
            return;
        }

        // Real telemetry: I-format ASDUs.
        if (apdu instanceof APDUIFormat apduiFormat) {
            unconfirmedPackets++;
            if (unconfirmedPackets >= ACK_EVERY_N_FRAMES) {
                try {
                    messageCodec.send(new APDUSFormat(0x01, apduiFormat.getReceiveSequenceNo() + 1));
                } catch (MessageCodecException e) {
                    LOGGER.error("Failed to send S-format acknowledgement", e);
                }
                unconfirmedPackets = 0;
            }
            processAsdu(apduiFormat.getAsdu());
        }
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
            if (!(subscriptionTag.getTag() instanceof Iec608705104Tag)) {
                values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            } else {
                values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK,
                    new Iec608705104SubscriptionHandle(this, (Iec608705104Tag) subscriptionTag.getTag())));
            }
        }
        return CompletableFuture.completedFuture(new DefaultPlcSubscriptionResponse(subscriptionRequest, values));
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        // No per-tag wire unsubscribe in IEC-104; reporting is per-session.
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
    // ASDU → PlcSubscriptionEvent dispatch
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void processAsdu(ASDU asdu) {
        int asduAddress = asdu.getAsduAddressField();
        for (InformationObject informationObject : asdu.getInformationObjects()) {
            int objectAddress = informationObject.getAddress();
            Iec608705104Tag tag = new Iec608705104Tag(asduAddress, objectAddress);
            PlcValue plcValue = Iec608705104TagParser.parseTag(informationObject, asdu.getTypeIdentification());
            LocalDateTime eventTime = extractEventTime(informationObject);
            publishEvent(eventTime, tag, plcValue, asduAddress, objectAddress);
        }
    }

    private static LocalDateTime extractEventTime(InformationObject informationObject) {
        if (informationObject instanceof InformationObjectWithTreeByteTime threeByte) {
            return convertCp24Time2aToCalendar(threeByte.getCp24Time2a());
        }
        if (informationObject instanceof InformationObjectWithSevenByteTime sevenByte) {
            return convertCp56Time2aToCalendar(sevenByte.getCp56Time2a());
        }
        return LocalDateTime.now();
    }

    /**
     * The CP24Time2a only carries minutes + milliseconds; we anchor it to the
     * current host clock for the year/month/day/hour fields.
     */
    static LocalDateTime convertCp24Time2aToCalendar(ThreeOctetBinaryTime cp24Time2) {
        LocalDateTime now = LocalDateTime.now();
        return LocalDateTime.of(
            now.getYear(), now.getMonthValue(), now.getDayOfMonth(), now.getHour(),
            cp24Time2.getMinutes(),
            cp24Time2.getMilliseconds() / 1000,
            (cp24Time2.getMilliseconds() % 1000) * 1_000_000);
    }

    /**
     * Devices typically send CP56Time2a in UTC. Convert to local time using
     * the JVM's default timezone (with DST honoured per the bit in the frame).
     */
    static LocalDateTime convertCp56Time2aToCalendar(SevenOctetBinaryTime cp56Time2) {
        TimeZone localTimeZone = TimeZone.getDefault();
        Duration localTimeZoneOffsetFromUtc = Duration.ofMillis(localTimeZone.getRawOffset());
        if (cp56Time2.getDaylightSaving()) {
            localTimeZoneOffsetFromUtc = localTimeZoneOffsetFromUtc.plus(Duration.ofMillis(localTimeZone.getDSTSavings()));
        }
        return LocalDateTime.of(
                2000 + cp56Time2.getYear(),
                cp56Time2.getMonth(),
                cp56Time2.getDay(),
                cp56Time2.getHour(),
                cp56Time2.getMinutes(),
                cp56Time2.getMilliseconds() / 1000,
                (cp56Time2.getMilliseconds() % 1000) * 1_000_000)
            .minus(localTimeZoneOffsetFromUtc);
    }

    /**
     * Hands one information object to every subscription that covers it.
     *
     * @param tag          the address of the information object, as it goes into the event.
     * @param asduAddress  the ASDU address of the information object, matched against the
     *                     subscriptions. Passed separately rather than read back off
     *                     {@code tag}, whose getters report {@link Iec608705104Tag#WILDCARD_ADDRESS}
     *                     instead of a value for a wildcarded address.
     * @param objectAddress the information object address, matched the same way.
     */
    private void publishEvent(LocalDateTime eventTime, Iec608705104Tag tag, PlcValue plcValue,
                              int asduAddress, int objectAddress) {
        Instant timestamp = eventTime.atZone(ZoneId.systemDefault()).toInstant();
        String tagKey = tag.getAddressString();
        Iec608705104PlcSubscriptionEvent event = new Iec608705104PlcSubscriptionEvent(
            timestamp,
            Collections.singletonMap(tagKey, tag),
            Collections.singletonMap(tagKey, new DefaultPlcResponseItem<>(PlcResponseCode.OK, plcValue)));

        // IEC-104 has no per-tag filtering on the wire: the station simply
        // reports everything it has. So the filtering happens here, by
        // matching the address of the incoming information object against the
        // (possibly wildcarded) tag of every registered subscription handle.
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            if (matchesAnyHandle(entry.getKey(), asduAddress, objectAddress)) {
                entry.getValue().accept(event);
            }
        }
    }

    /**
     * @return {@code true} if any subscription handle of {@code registration}
     * covers the given information object. Registrations without a single
     * IEC 60870-5-104 handle never match — there is no tag to match against.
     */
    static boolean matchesAnyHandle(PlcConsumerRegistration registration, int adsuAddress, int objectAddress) {
        for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
            if ((handle instanceof Iec608705104SubscriptionHandle iecHandle)
                && iecHandle.getTag().matches(adsuAddress, objectAddress)) {
                return true;
            }
        }
        return false;
    }

}
