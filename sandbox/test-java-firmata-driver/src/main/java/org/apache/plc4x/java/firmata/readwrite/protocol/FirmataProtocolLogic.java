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
package org.apache.plc4x.java.firmata.readwrite.protocol;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcBoolean;
import org.apache.plc4x.java.api.value.PlcList;
import org.apache.plc4x.java.firmata.readwrite.*;
import org.apache.plc4x.java.firmata.readwrite.context.FirmataDriverContext;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataField;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldDigital;
import org.apache.plc4x.java.firmata.readwrite.model.FirmataSubscriptionHandle;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.InternalPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.InternalPlcSubscriptionHandle;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class FirmataProtocolLogic extends Plc4xProtocolBase<FirmataMessage> implements PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(FirmataProtocolLogic.class);

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

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
            .only(SysexCommandReportFirmware.class)
            .handle(sysexCommandReportFirmware -> {
                String name = new String(sysexCommandReportFirmware.getFileName(), StandardCharsets.UTF_8);
                LOGGER.info(String.format("Connected to Firmata host running version %s.%s with name %s",
                    sysexCommandReportFirmware.getMajorVersion(), sysexCommandReportFirmware.getMinorVersion(),
                    name));
                context.fireConnected();
            });
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        return null;
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
            Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> result = new HashMap<>();
            for (String fieldName : subscriptionRequest.getFieldNames()) {
                FirmataField field = (FirmataField) subscriptionRequest.getField(fieldName);
                result.put(fieldName, Pair.of(PlcResponseCode.OK, new FirmataSubscriptionHandle(this, fieldName, field)));
            }
            future.complete(new DefaultPlcSubscriptionResponse((InternalPlcSubscriptionRequest) subscriptionRequest, result));
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
    protected void decode(ConversationContext<FirmataMessage> context, FirmataMessage msg) throws Exception {
        if(msg instanceof FirmataMessageCommand) {
            // Ignore ... for now ...
        } else {
            if(msg instanceof FirmataMessageAnalogIO) {
                // Analog values are single value messages (Value for one port only)
                FirmataMessageAnalogIO analogIO = (FirmataMessageAnalogIO) msg;
                int pin = analogIO.getPin();
                int analogValue = getAnalogValue(analogIO.getData());
                // If this is the first value, just add it.
                if(analogValues.get(pin) == null) {
                    analogValues.put(pin, new AtomicInteger(analogValue));
                    // TODO: Send an changed event ...
                }
                // If there was a value before and this is different to the current one, update it.
                else if(analogValue != analogValues.get(pin).intValue()) {
                    analogValues.get(pin).set(analogValue);
                    // TODO: Send an changed event ...
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
    }

    @Override
    public void close(ConversationContext<FirmataMessage> context) {

    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> collection) {
        final DefaultPlcConsumerRegistration consumerRegistration =
            new DefaultPlcConsumerRegistration(this, consumer, collection.toArray(new InternalPlcSubscriptionHandle[0]));
        consumers.put(consumerRegistration, consumer);
        return consumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration plcConsumerRegistration) {
        DefaultPlcConsumerRegistration consumerRegistration = (DefaultPlcConsumerRegistration) plcConsumerRegistration;
        consumers.remove(consumerRegistration);
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
            for (InternalPlcSubscriptionHandle handle : registration.getAssociatedHandles()) {
                if(handle instanceof FirmataSubscriptionHandle) {
                    FirmataSubscriptionHandle subscriptionHandle = (FirmataSubscriptionHandle) handle;
                    // Check if the subscription matches this current event
                    // (The bit subscribed to in this field actually changed).
                    if(subscriptionHandle.getField() instanceof FirmataFieldDigital) {
                        FirmataFieldDigital digitalField = (FirmataFieldDigital) subscriptionHandle.getField();
                        // If at least one bit of the current subscription changed it's value,
                        // send out an update event with all of its current values.
                        if(digitalField.getBitSet().intersects(changedBits)) {
                            List<PlcBoolean> values = new ArrayList<>(digitalField.getBitSet().cardinality());
                            for(int i = 0; i < digitalField.getBitSet().length(); i++) {
                                values.add(new PlcBoolean(bitValues.get(i)));
                            }
                            // If it's just one element, return this as a direct PlcValue
                            if(values.size() == 1) {
                                final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
                                    Collections.singletonMap(subscriptionHandle.getName(),
                                        Pair.of(PlcResponseCode.OK, values.get(0))));
                                consumer.accept(event);
                            }
                            // If it's more, return a PlcList instead.
                            else {
                                final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
                                    Collections.singletonMap(subscriptionHandle.getName(),
                                        Pair.of(PlcResponseCode.OK, new PlcList(values))));
                                consumer.accept(event);
                            }
                        }
                    }
                }
            }
        }
    }

    protected int getAnalogValue(byte[] data) {
        return 0;
    }

    protected int convertToSingleByteRepresentation(byte[] data) {
        byte result = data[0];
        result = (byte) (result | (((data[1] & 0x01) == 0x01) ? 0x80 : 0x00));
        return result & 0xFF;
    }

    protected BitSet getDigitalValues(int byteBlock, byte[] data) {
        int singleByte = convertToSingleByteRepresentation(data);
        if(byteBlock > 0) {
            singleByte = singleByte * (256 * byteBlock);
        }
        byte[] bitSetData = BigInteger.valueOf(singleByte).toByteArray();
        ArrayUtils.reverse(bitSetData);
        BitSet bitSet = BitSet.valueOf(bitSetData);
        return bitSet;
    }

}
