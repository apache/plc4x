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
package org.apache.plc4x.java.firmata;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.ConnectionStateChangeType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.firmata.configuration.FirmataConfiguration;
import org.apache.plc4x.java.firmata.model.FirmataSubscriptionHandle;
import org.apache.plc4x.java.firmata.tag.FirmataTag;
import org.apache.plc4x.java.firmata.tag.FirmataTagAnalog;
import org.apache.plc4x.java.firmata.tag.FirmataTagDigital;
import org.apache.plc4x.java.firmata.tag.FirmataTagHandler;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.*;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;
import org.apache.plc4x.java.firmata.readwrite.FirmataCommandSetDigitalPinValue;
import org.apache.plc4x.java.firmata.readwrite.FirmataCommandSetPinMode;
import org.apache.plc4x.java.firmata.readwrite.FirmataCommandSysex;
import org.apache.plc4x.java.firmata.readwrite.FirmataCommandSystemReset;
import org.apache.plc4x.java.firmata.readwrite.FirmataMessage;
import org.apache.plc4x.java.firmata.readwrite.FirmataMessageAnalogIO;
import org.apache.plc4x.java.firmata.readwrite.FirmataMessageCommand;
import org.apache.plc4x.java.firmata.readwrite.FirmataMessageDigitalIO;
import org.apache.plc4x.java.firmata.readwrite.FirmataMessageSubscribeAnalogPinValue;
import org.apache.plc4x.java.firmata.readwrite.FirmataMessageSubscribeDigitalPinValue;
import org.apache.plc4x.java.firmata.readwrite.PinMode;
import org.apache.plc4x.java.firmata.readwrite.SysexCommandReportFirmwareResponse;

/**
 * Firmata connection. Direct port of the legacy {@code FirmataProtocolLogic}
 * and {@code FirmataDriverContext} to the SPI3 {@link ConnectionBase} model.
 *
 * <p>Transport-agnostic: works over serial (the classic Firmata UART path),
 * TCP (Firmata-over-WiFi/Ethernet boards), or the in-process test transport.
 * The byte stream is identical — what differs is only how those bytes get
 * to the simulated/real device.</p>
 *
 * <p>Firmata is push-only: after a one-time {@code SystemReset} handshake the
 * board emits unsolicited {@code AnalogIO} / {@code DigitalIO} messages whenever
 * a configured pin changes. The connection caches the latest value per pin and
 * forwards changes to any registered subscription consumers whose subscribed
 * range intersects the changed pin.
 */
public class FirmataConnection extends ConnectionBase<FirmataConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(FirmataConnection.class);

    private FirmataMessageCodec messageCodec;

    private final AtomicBoolean connected = new AtomicBoolean(false);
    private volatile CompletableFuture<SysexCommandReportFirmwareResponse> pendingFirmwareReport;

    /** Latest sampled value for every analog pin we've seen. */
    private final Map<Integer, AtomicInteger> analogValues = new ConcurrentHashMap<>();
    /** Latest digital-pin state (one bit per pin). */
    private final BitSet digitalValues = new BitSet();

    /** Pin-mode tracking (formerly held by {@code FirmataDriverContext}). */
    private final Map<Integer, PinMode> digitalPins = new HashMap<>();
    private final Map<Integer, PinMode> analogPins = new HashMap<>();

    /** Active subscription consumers keyed by their registration. */
    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    public FirmataConnection(FirmataConfiguration configuration,
                             TransportInstance<?> transportInstance,
                             AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new FirmataTagHandler();
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
        return connected.get() && messageCodec != null && messageCodec.isOpen();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Connect / close
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new FirmataMessageCodec(transportInstance, this::handleIncomingMessage);

        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming Firmata data", e);
            }
        });

        // Send a SystemReset and wait for the SysexCommandReportFirmwareResponse the
        // board sends after a reset (the version-only ProtocolVersion message that
        // precedes it is discarded by handleIncomingMessage while !connected).
        pendingFirmwareReport = new CompletableFuture<>();
        try {
            messageCodec.send(new FirmataMessageCommand(new FirmataCommandSystemReset()));
            SysexCommandReportFirmwareResponse fw = pendingFirmwareReport
                .orTimeout(getConfiguration().getRequestTimeout(), TimeUnit.MILLISECONDS)
                .get();
            String name = new String(fw.getFileName(), StandardCharsets.UTF_8);
            LOGGER.info("Connected to Firmata host running version {}.{} with name {}",
                fw.getMajorVersion(), fw.getMinorVersion(), name);
        } catch (Exception e) {
            try {
                close();
            } catch (Exception ignored) {
                // ignore
            }
            throw new PlcConnectionException("Error during Firmata connect handshake", e);
        }

        connected.set(true);
        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONNECT, "Firmata serial connection established");
        }
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTED, null);
    }

    @Override
    public void close() throws Exception {
        connected.set(false);
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        if (pendingFirmwareReport != null && !pendingFirmwareReport.isDone()) {
            pendingFirmwareReport.completeExceptionally(new PlcRuntimeException("Connection closed"));
        }
        consumers.clear();
        super.close();
        LOGGER.info("Firmata serial connection closed");
        fireConnectionStateChanged(ConnectionStateChangeType.DISCONNECTED, null);
    }

    @Override
    protected void onTransportDisconnected(Throwable cause) {
        super.onTransportDisconnected(cause);
        connected.set(false);
        fireConnectionStateChanged(ConnectionStateChangeType.CONNECTION_LOST,
            cause != null ? cause.getMessage() : "Connection closed by remote");
        if (pendingFirmwareReport != null && !pendingFirmwareReport.isDone()) {
            pendingFirmwareReport.completeExceptionally(new PlcRuntimeException(
                cause != null ? "Connection lost: " + cause.getMessage() : "Connection closed by remote", cause));
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Incoming messages
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingMessage(FirmataMessage msg) {
        // Until the handshake completes, only react to the ReportFirmware response.
        // Anything else (notably the ProtocolVersion preamble) is intentionally
        // dropped so it doesn't pollute the per-pin caches before we know what
        // pins the user is interested in.
        if (!connected.get()) {
            if (msg instanceof FirmataMessageCommand mc
                && mc.getCommand() instanceof FirmataCommandSysex sysex
                && sysex.getCommand() instanceof SysexCommandReportFirmwareResponse fw
                && pendingFirmwareReport != null) {
                pendingFirmwareReport.complete(fw);
            }
            return;
        }

        if (msg instanceof FirmataMessageAnalogIO analogIO) {
            int pin = analogIO.getPin();
            int analogValue = decodeAnalogValue(analogIO.getData());
            AtomicInteger current = analogValues.get(pin);
            if (current == null || analogValue != current.intValue()) {
                analogValues.put(pin, new AtomicInteger(analogValue));
                publishAnalogEvents(pin);
            }
        } else if (msg instanceof FirmataMessageDigitalIO digitalIO) {
            BitSet incoming = decodeDigitalValues(digitalIO.getPinBlock(), digitalIO.getData());
            BitSet changed = new BitSet();
            for (int i = 0; i < 8; i++) {
                int bitPos = i + (8 * digitalIO.getPinBlock());
                if (digitalValues.get(bitPos) != incoming.get(bitPos)) {
                    changed.set(bitPos, true);
                    digitalValues.set(bitPos, incoming.get(bitPos));
                }
            }
            publishDigitalEvents(changed);
        } else {
            LOGGER.debug("Unexpected Firmata message {}", msg);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Write
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        try {
            List<FirmataMessage> messages = buildWriteMessages(writeRequest);
            for (FirmataMessage message : messages) {
                messageCodec.send(message);
            }
            // Firmata has no acknowledgement for writes — once we've shipped the
            // bytes there is no further confirmation to wait for.
            Map<String, PlcResponseCode> result = new LinkedHashMap<>();
            for (String tagName : writeRequest.getTagNames()) {
                result.put(tagName, PlcResponseCode.OK);
            }
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, result));
        } catch (Exception e) {
            return CompletableFuture.failedFuture(e);
        }
    }

    private List<FirmataMessage> buildWriteMessages(PlcWriteRequest writeRequest) {
        List<FirmataMessage> messages = new ArrayList<>();
        for (String tagName : writeRequest.getTagNames()) {
            PlcTag tag = writeRequest.getTag(tagName);
            if (!(tag instanceof FirmataTagDigital digitalTag)) {
                throw new PlcRuntimeException("Writing only supported for digital pins");
            }
            PlcValue plcValue = writeRequest.getPlcValue(tagName);
            if (digitalTag.getNumberOfElements() > 1 && plcValue.isList()) {
                PlcList list = (PlcList) plcValue;
                if (list.getList().size() != digitalTag.getNumberOfElements()) {
                    throw new PlcRuntimeException(
                        "Required " + digitalTag.getNumberOfElements() + " values but got " + list.getList().size());
                }
            }
            for (int i = 0; i < digitalTag.getNumberOfElements(); i++) {
                int pin = digitalTag.getAddress() + i;
                PinMode existing = digitalPins.get(pin);
                if (existing == null) {
                    digitalPins.put(pin, PinMode.PinModeOutput);
                    messages.add(new FirmataMessageCommand(
                        new FirmataCommandSetPinMode((short) pin, PinMode.PinModeOutput)));
                } else if (existing != PinMode.PinModeOutput) {
                    throw new PlcRuntimeException(
                        "Pin " + pin + " already configured as " + existing.name());
                }
                messages.add(new FirmataMessageCommand(
                    new FirmataCommandSetDigitalPinValue((short) pin, plcValue.getIndex(i).getBoolean())));
            }
        }
        return messages;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe / unsubscribe
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        try {
            List<FirmataMessage> messages = buildSubscribeMessages(subscriptionRequest);
            for (FirmataMessage message : messages) {
                messageCodec.send(message);
            }
            Map<String, PlcResponseItem<PlcSubscriptionHandle>> result = new LinkedHashMap<>();
            for (String tagName : subscriptionRequest.getTagNames()) {
                DefaultPlcSubscriptionTag subscriptionTag =
                    (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(tagName);
                FirmataTag tag = (FirmataTag) subscriptionTag.getTag();
                result.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK,
                    new FirmataSubscriptionHandle(this, tagName, tag)));
            }
            return CompletableFuture.completedFuture(new DefaultPlcSubscriptionResponse(subscriptionRequest, result));
        } catch (Exception e) {
            return CompletableFuture.failedFuture(e);
        }
    }

    private List<FirmataMessage> buildSubscribeMessages(PlcSubscriptionRequest subscriptionRequest) {
        Map<Integer, PinMode> requestedDigital = new LinkedHashMap<>();
        Map<Integer, PinMode> requestedAnalog = new LinkedHashMap<>();
        for (String tagName : subscriptionRequest.getTagNames()) {
            PlcTag rawTag = subscriptionRequest.getTag(tagName);
            DefaultPlcSubscriptionTag subscriptionTag = (DefaultPlcSubscriptionTag) rawTag;
            PlcTag tag = subscriptionTag.getTag();
            if (tag instanceof FirmataTagDigital digital) {
                PinMode mode = digital.getPinMode() != null ? digital.getPinMode() : PinMode.PinModeInput;
                if (mode != PinMode.PinModeInput && mode != PinMode.PinModePullup) {
                    throw new PlcInvalidTagException("Subscription tag must be of type 'INPUT' (default) or 'PULLUP'");
                }
                for (int pin = digital.getAddress(); pin < digital.getAddress() + digital.getNumberOfElements(); pin++) {
                    requestedDigital.put(pin, mode);
                }
            } else if (tag instanceof FirmataTagAnalog analog) {
                for (int pin = analog.getAddress(); pin < analog.getAddress() + analog.getNumberOfElements(); pin++) {
                    requestedAnalog.put(pin, PinMode.PinModeInput);
                }
            } else {
                throw new PlcRuntimeException("Unsupported tag type " + tag.getClass().getSimpleName());
            }
        }

        // Drop any pins that are already subscribed in the same mode; mismatched
        // mode requests for an already-claimed pin are an error.
        requestedDigital.entrySet().removeIf(entry -> {
            PinMode existing = digitalPins.get(entry.getKey());
            if (existing == null) {
                return false;
            }
            if (existing != entry.getValue()) {
                throw new PlcInvalidTagException(String.format(
                    "Error setting digital pin %d to mode %s, pin is already set to mode %s",
                    entry.getKey(), entry.getValue(), existing));
            }
            return true;
        });
        requestedAnalog.entrySet().removeIf(entry -> analogPins.containsKey(entry.getKey()));

        List<FirmataMessage> messages = new ArrayList<>();
        Set<Integer> digitalPortsToReport = new HashSet<>();
        for (Map.Entry<Integer, PinMode> entry : requestedDigital.entrySet()) {
            int pin = entry.getKey();
            PinMode mode = entry.getValue();
            digitalPins.put(pin, mode);
            messages.add(new FirmataMessageCommand(new FirmataCommandSetPinMode((short) pin, mode)));
            // SubscribeDigitalPinValue (0xD0 | port) is per-port (8 pins), not
            // per-pin. Send one REPORT_DIGITAL per port the subscription
            // touches.
            digitalPortsToReport.add(pin / 8);
        }
        for (int port : digitalPortsToReport) {
            messages.add(new FirmataMessageSubscribeDigitalPinValue((byte) port, true));
        }
        for (Map.Entry<Integer, PinMode> entry : requestedAnalog.entrySet()) {
            int pin = entry.getKey();
            analogPins.put(pin, entry.getValue());
            messages.add(new FirmataMessageSubscribeAnalogPinValue((byte) pin, true));
        }
        return messages;
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        // Firmata doesn't expose a clean per-tag unsubscribe — once a pin is in
        // reporting mode the board keeps streaming updates. The legacy driver
        // returned null here; we return a completed response so the API contract
        // is honoured.
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
    // Event publication
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void publishAnalogEvents(int pin) {
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            DefaultPlcConsumerRegistration registration = entry.getKey();
            Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
                if (handle instanceof FirmataSubscriptionHandle sub
                    && sub.getTag() instanceof FirmataTagAnalog analog
                    && analog.getAddress() <= pin
                    && analog.getAddress() + analog.getNumberOfElements() >= pin) {
                    List<PlcValue> values = new ArrayList<>(analog.getNumberOfElements());
                    for (int i = analog.getAddress(); i < analog.getAddress() + analog.getNumberOfElements(); i++) {
                        AtomicInteger v = analogValues.get(i);
                        values.add(new PlcDINT(v != null ? v.intValue() : -1));
                    }
                    deliverEvent(consumer, sub.getName(), values);
                }
            }
        }
    }

    private void publishDigitalEvents(BitSet changedBits) {
        if (changedBits.cardinality() == 0) {
            return;
        }
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            DefaultPlcConsumerRegistration registration = entry.getKey();
            Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
                if (handle instanceof FirmataSubscriptionHandle sub
                    && sub.getTag() instanceof FirmataTagDigital digital
                    && digital.getBitSet().intersects(changedBits)) {
                    // Emit one value per pin in the subscription's pin range,
                    // not zero..bitSet.length() which would include all the
                    // unset bits below the highest subscribed pin too.
                    int from = digital.getAddress();
                    int to = from + digital.getNumberOfElements();
                    List<PlcValue> values = new ArrayList<>(digital.getNumberOfElements());
                    for (int i = from; i < to; i++) {
                        values.add(new PlcBOOL(digitalValues.get(i)));
                    }
                    deliverEvent(consumer, sub.getName(), values);
                }
            }
        }
    }

    private void deliverEvent(Consumer<PlcSubscriptionEvent> consumer, String tagName, List<PlcValue> values) {
        PlcValue value = values.size() == 1 ? values.get(0) : new PlcList(values);
        PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
            Collections.singletonMap(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, value)));
        consumer.accept(event);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Wire-format helpers (unchanged from the legacy protocol logic)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private static int decodeAnalogValue(List<Byte> data) {
        // Firmata analog values are encoded as a 14-bit integer split across two
        // bytes: bits 0-6 of byte0 are the LSBs, bits 0-6 of byte1 are the MSBs.
        return ((data.get(0) & 0xFF) | (data.get(1) << 7)) & 0xFFFF;
    }

    private static BitSet decodeDigitalValues(int byteBlock, List<Byte> data) {
        int singleByte = mergeDigitalBytes(data);
        if (byteBlock > 0) {
            singleByte = singleByte * (256 * byteBlock);
        }
        byte[] bitSetData = BigInteger.valueOf(singleByte).toByteArray();
        ArrayUtils.reverse(bitSetData);
        return BitSet.valueOf(bitSetData);
    }

    private static int mergeDigitalBytes(List<Byte> data) {
        byte result = data.get(0);
        result = (byte) (result | (((data.get(1) & 0x01) == 0x01) ? 0x80 : 0x00));
        return result & 0xFF;
    }

}
