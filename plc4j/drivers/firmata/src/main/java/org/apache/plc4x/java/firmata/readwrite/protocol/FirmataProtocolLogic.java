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
package org.apache.plc4x.java.firmata.readwrite.protocol;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.firmata.readwrite.*;
import org.apache.plc4x.java.firmata.readwrite.context.FirmataDriverContext;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataField;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldAnalog;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldDigital;
import org.apache.plc4x.java.firmata.readwrite.model.FirmataSubscriptionHandle;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcList;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class FirmataProtocolLogic extends Plc4xProtocolBase<FirmataMessage> implements PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(FirmataProtocolLogic.class);

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private AtomicBoolean connected = new AtomicBoolean(false);
    private Map<Integer, AtomicInteger> analogValues = new HashMap<>();
    private BitSet digitalValues = new BitSet();

    private Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    @Override
    public void onConnect(ConversationContext<FirmataMessage> context) {
        LOGGER.debug("Sending Firmata Reset Command");
        FirmataMessageCommand resetCommandMessage = new FirmataMessageCommand(new FirmataCommandSystemReset());
        context.sendRequest(resetCommandMessage)
            .expectResponse(FirmataMessage.class, REQUEST_TIMEOUT)
            .only(FirmataMessageCommand.class)
            .unwrap(FirmataMessageCommand::getCommand)
            .only(FirmataCommandSysex.class)
            .unwrap(FirmataCommandSysex::getCommand)
            .only(SysexCommandReportFirmwareResponse.class)
            .handle(sysexCommandReportFirmware -> {
                String name = new String(sysexCommandReportFirmware.getFileName(), StandardCharsets.UTF_8);
                LOGGER.info(String.format("Connected to Firmata host running version %s.%s with name %s",
                    sysexCommandReportFirmware.getMajorVersion(), sysexCommandReportFirmware.getMinorVersion(),
                    name));
                connected.set(true);
                context.fireConnected();
            });
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        try {
            final List<FirmataMessage> firmataMessages =
                ((FirmataDriverContext) getDriverContext()).processWriteRequest(writeRequest);
            for (FirmataMessage firmataMessage : firmataMessages) {
                context.sendToWire(firmataMessage);
            }
            // There's unfortunately no ack response :-(
            Map<String, PlcResponseCode> result = new HashMap<>();
            for (String fieldName : writeRequest.getFieldNames()) {
                result.put(fieldName, PlcResponseCode.OK);
            }
            future.complete(new DefaultPlcWriteResponse(writeRequest, result));
        } catch (PlcRuntimeException e) {
            future.completeExceptionally(e);
        }
        return future;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        try {
            final List<FirmataMessage> firmataMessages =
                ((FirmataDriverContext) getDriverContext()).processSubscriptionRequest(subscriptionRequest);
            for (FirmataMessage firmataMessage : firmataMessages) {
                context.sendToWire(firmataMessage);
            }
            Map<String, ResponseItem<PlcSubscriptionHandle>> result = new HashMap<>();
            for (String fieldName : subscriptionRequest.getFieldNames()) {
                DefaultPlcSubscriptionField subscriptionField =
                    (DefaultPlcSubscriptionField) subscriptionRequest.getField(fieldName);
                FirmataField field = (FirmataField) subscriptionField.getPlcField();
                result.put(fieldName, new ResponseItem<>(PlcResponseCode.OK,
                    new FirmataSubscriptionHandle(this, fieldName, field)));
            }
            future.complete(new DefaultPlcSubscriptionResponse(subscriptionRequest, result));
        } catch (PlcRuntimeException e) {
            future.completeExceptionally(e);
        }
        return future;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return null;
    }

    @Override
    protected void decode(ConversationContext<FirmataMessage> context, FirmataMessage msg) {
        // Especially when we restart there might already be data incoming before we actually are finished
        // setting up. Just ignore all incoming data until we're officially connected.
        if(!connected.get()) {
            return;
        }

        if(msg instanceof FirmataMessageAnalogIO) {
            // Analog values are single value messages (Value for one port only)
            FirmataMessageAnalogIO analogIO = (FirmataMessageAnalogIO) msg;
            int pin = analogIO.getPin();
            int analogValue = getAnalogValue(analogIO.getData());
            // If this is the first value, or the value changed, send update events..
            if((analogValues.get(pin) == null) || (analogValue != analogValues.get(pin).intValue())) {
                analogValues.put(pin, new AtomicInteger(analogValue));
                publishAnalogEvents(pin, analogValue);
            }
        } else if(msg instanceof FirmataMessageDigitalIO) {
            // Digital values come 8 pins together (ignoring the pin value, which is always 0).
            FirmataMessageDigitalIO digitalIO = (FirmataMessageDigitalIO) msg;
            BitSet newDigitalValues = getDigitalValues(digitalIO.getPinBlock(), digitalIO.getData());

            // Compare the currently set bits with the ones from the last time to see what's changed.
            BitSet changedBits = new BitSet();
            for (int i = 0; i < 8; i++) {
                int bitPos = i + (8 * digitalIO.getPinBlock());
                if (digitalValues.get(bitPos) != newDigitalValues.get(bitPos)) {
                    changedBits.set(bitPos, true);
                    digitalValues.set(bitPos, newDigitalValues.get(bitPos));
                }
            }

            // Send out update events.
            publishDigitalEvents(changedBits, digitalValues);
        } else {
            LOGGER.debug(String.format("Unexpected message %s", msg.toString()));
        }
    }

    @Override
    public void close(ConversationContext<FirmataMessage> context) {
        connected.set(false);
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> collection) {
        final DefaultPlcConsumerRegistration consumerRegistration =
            new DefaultPlcConsumerRegistration(this, consumer, collection.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(consumerRegistration, consumer);
        return consumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration plcConsumerRegistration) {
        DefaultPlcConsumerRegistration consumerRegistration = (DefaultPlcConsumerRegistration) plcConsumerRegistration;
        consumers.remove(consumerRegistration);
    }

    protected void publishAnalogEvents(int pin, int value) {
        // Try sending the subscription event to all listeners.
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            final DefaultPlcConsumerRegistration registration = entry.getKey();
            final Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            // Only if the current data point matches the subscription, publish the event to it.
            for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
                if (handle instanceof FirmataSubscriptionHandle) {
                    FirmataSubscriptionHandle subscriptionHandle = (FirmataSubscriptionHandle) handle;
                    // Check if the subscription matches this current event
                    // (The bit subscribed to in this field actually changed).
                    if (subscriptionHandle.getField() instanceof FirmataFieldAnalog) {
                        FirmataFieldAnalog analogField = (FirmataFieldAnalog) subscriptionHandle.getField();
                        // Check if this field would include the current pin.
                        if((analogField.getAddress() <= pin) &&
                            (analogField.getAddress() + analogField.getNumberOfElements() >= pin)) {
                            // Build an update event containing the current values for all subscribed fields.
                            List<PlcValue> values = new ArrayList<>(analogField.getNumberOfElements());
                            for(int i = analogField.getAddress(); i < analogField.getAddress() + analogField.getNumberOfElements(); i++) {
                                if(analogValues.containsKey(i)) {
                                    values.add(new PlcDINT(analogValues.get(i).intValue()));
                                }
                                // This could be the case if only some of the requested array values are available
                                else {
                                    values.add(new PlcDINT(-1));
                                }
                            }
                            sendUpdateEvents(consumer, subscriptionHandle.getName(), values);
                        }
                    }
                }
            }
        }
    }

    protected void publishDigitalEvents(BitSet changedBits, BitSet bitValues) {
        // If nothing changed, no need to do anything.
        if(changedBits.cardinality() == 0) {
            return;
        }
        // Try sending the subscription event to all listeners.
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            final DefaultPlcConsumerRegistration registration = entry.getKey();
            final Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            // Only if the current data point matches the subscription, publish the event to it.
            for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
                if(handle instanceof FirmataSubscriptionHandle) {
                    FirmataSubscriptionHandle subscriptionHandle = (FirmataSubscriptionHandle) handle;
                    // Check if the subscription matches this current event
                    // (The bit subscribed to in this field actually changed).
                    if(subscriptionHandle.getField() instanceof FirmataFieldDigital) {
                        FirmataFieldDigital digitalField = (FirmataFieldDigital) subscriptionHandle.getField();
                        // If at least one bit of the current subscription changed it's value,
                        // send out an update event with all of its current values.
                        if(digitalField.getBitSet().intersects(changedBits)) {
                            List<PlcValue> values = new ArrayList<>(digitalField.getBitSet().cardinality());
                            for(int i = 0; i < digitalField.getBitSet().length(); i++) {
                                values.add(new PlcBOOL(bitValues.get(i)));
                            }
                            sendUpdateEvents(consumer, subscriptionHandle.getName(), values);
                        }
                    }
                }
            }
        }
    }

    protected void sendUpdateEvents(Consumer<PlcSubscriptionEvent> consumer, String fieldName, List<PlcValue> values) {
        // If it's just one element, return this as a direct PlcValue
        if(values.size() == 1) {
            final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
                Collections.singletonMap(fieldName, new ResponseItem<>(PlcResponseCode.OK, values.get(0))));
            consumer.accept(event);
        }
        // If it's more, return a PlcList instead.
        else {
            final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
                Collections.singletonMap(fieldName, new ResponseItem<>(PlcResponseCode.OK, new PlcList(values))));
            consumer.accept(event);
        }
    }

    protected int getAnalogValue(List<Byte> data) {
        // In Firmata analog values are encoded as a 14bit integer with the least significant bits being located in
        // the bits 0-6 of the birst byte and the second half as the 0-6 bits of the second byte.
        return ((data.get(0) & 0xFF)| (data.get(1) << 7)) & 0xFFFF;
    }

    protected int convertToSingleByteRepresentation(List<Byte> data) {
        byte result = data.get(0);
        result = (byte) (result | (((data.get(1) & 0x01) == 0x01) ? 0x80 : 0x00));
        return result & 0xFF;
    }

    protected BitSet getDigitalValues(int byteBlock, List<Byte> data) {
        int singleByte = convertToSingleByteRepresentation(data);
        if(byteBlock > 0) {
            singleByte = singleByte * (256 * byteBlock);
        }
        byte[] bitSetData = BigInteger.valueOf(singleByte).toByteArray();
        ArrayUtils.reverse(bitSetData);
        return BitSet.valueOf(bitSetData);
    }

}
