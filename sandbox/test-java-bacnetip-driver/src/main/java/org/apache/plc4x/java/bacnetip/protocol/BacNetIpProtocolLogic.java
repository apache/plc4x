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
package org.apache.plc4x.java.bacnetip.protocol;

import java.time.Instant;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Consumer;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcInteger;
import org.apache.plc4x.java.api.value.PlcLong;
import org.apache.plc4x.java.api.value.PlcString;
import org.apache.plc4x.java.api.value.PlcStruct;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.bacnetip.configuration.BacNetIpConfiguration;
import org.apache.plc4x.java.bacnetip.ede.model.Datapoint;
import org.apache.plc4x.java.bacnetip.ede.model.EdeModel;
import org.apache.plc4x.java.bacnetip.field.BacNetIpField;
import org.apache.plc4x.java.bacnetip.readwrite.APDU;
import org.apache.plc4x.java.bacnetip.readwrite.APDUComplexAck;
import org.apache.plc4x.java.bacnetip.readwrite.APDUConfirmedRequest;
import org.apache.plc4x.java.bacnetip.readwrite.APDUError;
import org.apache.plc4x.java.bacnetip.readwrite.APDUSimpleAck;
import org.apache.plc4x.java.bacnetip.readwrite.APDUUnconfirmedRequest;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequest;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequestConfirmedCOVNotification;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequestReadProperty;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequestSubscribeCOV;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequestWriteProperty;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetTag;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetTagWithContent;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequest;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestIAm;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestWhoHas;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestWhoIs;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.bacnetip.readwrite.BVLCForwardedNPDU;
import org.apache.plc4x.java.bacnetip.readwrite.BVLCOriginalBroadcastNPDU;
import org.apache.plc4x.java.bacnetip.readwrite.BVLCOriginalUnicastNPDU;
import org.apache.plc4x.java.bacnetip.readwrite.NPDU;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.InternalPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.apache.plc4x.java.spi.model.InternalPlcSubscriptionHandle;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class BacNetIpProtocolLogic extends Plc4xProtocolBase<BVLC> implements HasConfiguration<BacNetIpConfiguration>, PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(BacNetIpProtocolLogic.class);

    private AtomicInteger requestId = new AtomicInteger();

    private Map<Integer, Consumer<PlcSubscriptionEvent>> consumerIdMap = new ConcurrentHashMap<>();

    @Override
    public void setConfiguration(BacNetIpConfiguration configuration) {

    }

    @Override
    public void onConnect(ConversationContext<BVLC> context) {
        context.fireConnected();

        APDUUnconfirmedRequest apdu = new APDUUnconfirmedRequest(
            new BACnetUnconfirmedServiceRequestWhoIs((short) 0x08, null, null)
        );

        BVLCOriginalBroadcastNPDU broadcastNPDU = new BVLCOriginalBroadcastNPDU(
            new NPDU(
                (short) 0,
                true,
                false,
                false,
                false,
                (byte) 0x0,
                0,
                (short) 0,
                new short[] {},
                0,
                (short) 0,
                new short[] {},
                (short) 0,
                null,
                apdu
            )
        );
        context.sendToWire(broadcastNPDU);
    }

    @Override
    public void close(ConversationContext<BVLC> context) {
        // Nothing to do here ...
    }

    @Override
    protected void decode(ConversationContext<BVLC> context, BVLC msg) throws Exception {
        NPDU npdu = null;
        if(msg instanceof BVLCOriginalUnicastNPDU) {
            BVLCOriginalUnicastNPDU bvlcOriginalUnicastNPDU = (BVLCOriginalUnicastNPDU) msg;
            npdu = bvlcOriginalUnicastNPDU.getNpdu();
        } else if (msg instanceof BVLCForwardedNPDU) {
            BVLCForwardedNPDU bvlcForwardedNPDU = (BVLCForwardedNPDU) msg;
            npdu = bvlcForwardedNPDU.getNpdu();
        } else if (msg instanceof BVLCOriginalBroadcastNPDU) {
            BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) msg;
            npdu = bvlcOriginalBroadcastNPDU.getNpdu();
        }

        if(npdu != null) {
            if(npdu.getApdu() instanceof APDUConfirmedRequest) {
                APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                final BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                // A value change subscription event.
                if(serviceRequest instanceof BACnetConfirmedServiceRequestConfirmedCOVNotification) {
                    BACnetConfirmedServiceRequestConfirmedCOVNotification valueChange =
                        (BACnetConfirmedServiceRequestConfirmedCOVNotification) serviceRequest;

                }
                // Someone read a value.
                else if(serviceRequest instanceof BACnetConfirmedServiceRequestReadProperty) {
                    // Ignore this ...
                }
                // Someone wrote a value.
                else if(serviceRequest instanceof BACnetConfirmedServiceRequestWriteProperty) {
                    // Ignore this ...
                } else if(serviceRequest instanceof BACnetConfirmedServiceRequestSubscribeCOV) {
                    // Ignore this ...
                } else {
                    LOGGER.debug(String.format("Unexpected ConfirmedServiceRequest type: %s", serviceRequest.getClass().getName()));
                }
            } else if(npdu.getApdu() instanceof APDUUnconfirmedRequest) {
                APDUUnconfirmedRequest unconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                final BACnetUnconfirmedServiceRequest serviceRequest = unconfirmedRequest.getServiceRequest();
                if(serviceRequest instanceof BACnetUnconfirmedServiceRequestWhoHas) {
                    // Ignore this ...
                } else if(serviceRequest instanceof BACnetUnconfirmedServiceRequestWhoIs){
                    // Ignore this ...
                } else if(serviceRequest instanceof BACnetUnconfirmedServiceRequestIAm){
                    // Ignore this ...
                } else if(serviceRequest instanceof BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer){
                    // Ignore this ...
                } else {
                    LOGGER.debug(String.format("Unexpected UnconfirmedServiceRequest type: %s", serviceRequest.getClass().getName()));
                }
            } else if(npdu.getApdu() instanceof APDUError) {
                APDUError apduError = (APDUError) npdu.getApdu();
            } else if(npdu.getApdu() instanceof APDUSimpleAck) {
                // Ignore this ...
            } else if(npdu.getApdu() instanceof APDUComplexAck) {
                // Ignore this ...
            } else if((npdu.getApdu() == null) && (npdu.getNlm() != null)){
                // "Who is router?" & "I am router" messages.
                // Ignore this ...
            } else {
                LOGGER.debug(String.format("Unexpected NPDU type: %s", npdu.getClass().getName()));
            }
        }
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        Map<String, ResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
        for (String fieldName : subscriptionRequest.getFieldNames()) {
            values.put(fieldName, new ResponseItem<>(PlcResponseCode.OK, new DefaultPlcSubscriptionHandle(this)));
        }

        return CompletableFuture.completedFuture(
            new DefaultPlcSubscriptionResponse((InternalPlcSubscriptionRequest) subscriptionRequest, values));
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> collection) {
        final DefaultPlcConsumerRegistration consumerRegistration =
            new DefaultPlcConsumerRegistration(this, consumer, collection.toArray(new InternalPlcSubscriptionHandle[0]));
        consumerIdMap.put(consumerRegistration.getConsumerHash(), consumer);
        return consumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration plcConsumerRegistration) {
        DefaultPlcConsumerRegistration consumerRegistration = (DefaultPlcConsumerRegistration) plcConsumerRegistration;
        consumerIdMap.remove(consumerRegistration.getConsumerHash());
    }

    protected void publishEvent(BacNetIpField field, PlcValue plcValue) {
        // Create a subscription event from the input.
        final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
            Collections.singletonMap("event", new ResponseItem(PlcResponseCode.OK, plcValue)));

        // Send the subscription event to all listeners.
        for (Consumer<PlcSubscriptionEvent> consumer : consumerIdMap.values()) {
            // TODO: Check if the subscription matches the current field ..
            consumer.accept(event);
        }
    }

    private String toString(BacNetIpField field) {
        return field.getDeviceIdentifier() + "/" + field.getObjectType() + "/" + field.getObjectInstance();
    }

}
