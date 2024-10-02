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

package org.apache.plc4x.java.iec608705104.readwrite.protocol;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.iec608705104.readwrite.*;
import org.apache.plc4x.java.iec608705104.readwrite.configuration.Iec608705014Configuration;
import org.apache.plc4x.java.iec608705104.readwrite.messages.Iec608705104PlcSubscriptionEvent;
import org.apache.plc4x.java.iec608705104.readwrite.model.Iec608705104SubscriptionHandle;
import org.apache.plc4x.java.iec608705104.readwrite.tag.Iec608705104Tag;
import org.apache.plc4x.java.iec608705104.readwrite.tag.Iec608705104TagHandler;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.PlcBrowser;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;

import java.time.Duration;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

public class Iec608705104Protocol extends Plc4xProtocolBase<APDU> implements HasConfiguration<Iec608705014Configuration>, PlcSubscriber, PlcBrowser {

    private Iec608705014Configuration configuration;
    private final RequestTransactionManager tm;

    private int unconfirmedPackets;

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    public Iec608705104Protocol() {
        // We're starting with allowing only one message in-flight.
        this.tm = new RequestTransactionManager(1);
        unconfirmedPackets = 0;
    }

    @Override
    public void setConfiguration(Iec608705014Configuration configuration) {
        this.configuration = configuration;
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new Iec608705104TagHandler();
    }

    @Override
    public void close(ConversationContext<APDU> context) {
        tm.shutdown();
    }

    @Override
    public void onConnect(ConversationContext<APDU> context) {
        // First we exchange a test-frame
        APDUUFormatTestFrameActivation testFrameActivation = new APDUUFormatTestFrameActivation(0x43);
        RequestTransactionManager.RequestTransaction testFrameTx = tm.startRequest();
        testFrameTx.submit(() -> context.sendRequest(testFrameActivation)
            .expectResponse(APDU.class, Duration.ofMillis(configuration.getRequestTimeout()))
            .onTimeout(e -> context.getChannel().pipeline().fireExceptionCaught(e))
            .onError((p, e) -> context.getChannel().pipeline().fireExceptionCaught(e))
            .only(APDUUFormatTestFrameConfirmation.class)
            .handle(testFrameResponse -> {
                testFrameTx.endRequest();

                // Next send the start-data-transfer packet.
                APDUUFormatStartDataTransferActivation startDataTransferActivation = new APDUUFormatStartDataTransferActivation(0x07);
                RequestTransactionManager.RequestTransaction startDataTransferTx = tm.startRequest();
                startDataTransferTx.submit(() -> context.sendRequest(startDataTransferActivation)
                    .expectResponse(APDU.class, Duration.ofMillis(configuration.getRequestTimeout()))
                    .onTimeout(e -> context.getChannel().pipeline().fireExceptionCaught(e))
                    .onError((p, e) -> context.getChannel().pipeline().fireExceptionCaught(e))
                    .only(APDUUFormatStartDataTransferConfirmation.class)
                    .handle(startDataTransferResponse -> {
                        startDataTransferTx.endRequest();
                        context.fireConnected();
                    }));
            }));
    }

    @Override
    protected void decode(ConversationContext<APDU> context, APDU msg) throws Exception {
        // When receiving a test-frame, send the expected response.
        if (msg instanceof APDUUFormatTestFrameActivation) {
            APDUUFormatTestFrameConfirmation testFrameConfirmation = new APDUUFormatTestFrameConfirmation(0x83);
            context.sendToWire(testFrameConfirmation);
        }
        // When receiving incoming data, process that.
        else if (msg instanceof APDUIFormat) {
            APDUIFormat apduiFormat = (APDUIFormat) msg;

            // Make sure we send an acknowledgement packet every few packets.
            unconfirmedPackets++;
            if (unconfirmedPackets >= 8) {
                // Confirm the reception of the packet.
                APDUSFormat confirmPacket = new APDUSFormat(0x01, apduiFormat.getReceiveSequenceNo() + 1);
                context.sendToWire(confirmPacket);
                unconfirmedPackets = 0;
            }

            // Handle the incoming messages.
            processData(apduiFormat.getAsdu());
        }
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
        for (String tagName : subscriptionRequest.getTagNames()) {
            final DefaultPlcSubscriptionTag tag = (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(tagName);
            if (!(tag.getTag() instanceof Iec608705104Tag)) {
                values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            } else {
                values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK,
                    new Iec608705104SubscriptionHandle(this, (Iec608705104Tag) tag.getTag())));
            }
        }
        return CompletableFuture.completedFuture(
            new DefaultPlcSubscriptionResponse(subscriptionRequest, values));
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

    protected void processData(ASDU asdu) {
        int asduAddress = asdu.getAsduAddressField();
        for (InformationObject informationObject : asdu.getInformationObjects()) {
            int objectAddress = informationObject.getAddress();
            Iec608705104Tag tag = new Iec608705104Tag(asduAddress, objectAddress);
            PlcValue plcValue = Iec608705104TagParser.parseTag(informationObject, asdu.getTypeIdentification());

            // If this is a datatype that comes with time-information, parse this
            // and use this instead of the generated timestamp.
            LocalDateTime eventTime;
            if (informationObject instanceof InformationObjectWithTreeByteTime) {
                InformationObjectWithTreeByteTime informationObjectWithTreeByteTime = (InformationObjectWithTreeByteTime) informationObject;
                ThreeOctetBinaryTime time = informationObjectWithTreeByteTime.getCp24Time2a();
                eventTime = convertCp24Time2aToCalendar(time);
            } else if (informationObject instanceof InformationObjectWithSevenByteTime) {
                InformationObjectWithSevenByteTime informationObjectWithSevenByteTime = (InformationObjectWithSevenByteTime) informationObject;
                SevenOctetBinaryTime time = informationObjectWithSevenByteTime.getCp56Time2a();
                eventTime = convertCp56Time2aToCalendar(time);
            } else {
                eventTime = LocalDateTime.now();
            }

            // Send the event out to all subscribed listeners.
            publishEvent(eventTime, tag, plcValue);
        }
    }

    protected LocalDateTime convertCp24Time2aToCalendar(ThreeOctetBinaryTime cp24Time2) {
        LocalDateTime now = LocalDateTime.now();
        return LocalDateTime.of(now.getYear(), now.getMonthValue(), now.getDayOfMonth(), now.getHour(), cp24Time2.getMinutes(), cp24Time2.getMilliseconds() / 1000, (cp24Time2.getMilliseconds() % 1000) * 1000000);
    }

    protected LocalDateTime convertCp56Time2aToCalendar(SevenOctetBinaryTime cp56Time2) {
        // It seems that the time is sent in UTC, so we need to convert that into our local timezone.
        TimeZone localTimeZone = TimeZone.getDefault();
        Duration localTimeZoneOffsetFromUTC = Duration.ofMillis(localTimeZone.getRawOffset());
        if (cp56Time2.getDaylightSaving()) {
            Duration daylightSavingOffset = Duration.ofMillis(localTimeZone.getDSTSavings());
            localTimeZoneOffsetFromUTC = localTimeZoneOffsetFromUTC.plus(daylightSavingOffset);
        }
        return LocalDateTime.of(2000 + cp56Time2.getYear(), cp56Time2.getMonth(), cp56Time2.getDay(), cp56Time2.getHour(), cp56Time2.getMinutes(), cp56Time2.getMilliseconds() / 1000, (cp56Time2.getMilliseconds() % 1000) * 1000000)
            .minus(localTimeZoneOffsetFromUTC);
    }

    protected void publishEvent(LocalDateTime timeStamp, Iec608705104Tag tag, PlcValue plcValue) {
        // Create a subscription event from the input.
        final PlcSubscriptionEvent event = new Iec608705104PlcSubscriptionEvent(
            timeStamp.atZone(ZoneId.systemDefault()).toInstant(),
            Collections.singletonMap(tag.toString(), tag),
            Collections.singletonMap(tag.toString(),
                new DefaultPlcResponseItem<>(PlcResponseCode.OK, plcValue)));

        // Try sending the subscription event to all listeners.
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            final DefaultPlcConsumerRegistration registration = entry.getKey();
            final Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            // Only if the current data point matches the subscription, publish the event to it.
            for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
                if (handle instanceof Iec608705104SubscriptionHandle) {
                    Iec608705104SubscriptionHandle subscriptionHandle = (Iec608705104SubscriptionHandle) handle;
                    // Check if the subscription matches this current event.
                    if (/*subscriptionHandle.getTag().matchesGroupAddress(groupAddress)*/true) {
                        consumer.accept(event);
                    }
                }
            }
        }
    }

}
