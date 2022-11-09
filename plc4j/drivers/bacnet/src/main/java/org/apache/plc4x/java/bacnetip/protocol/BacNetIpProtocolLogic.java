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
package org.apache.plc4x.java.bacnetip.protocol;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.bacnetip.configuration.BacNetIpConfiguration;
import org.apache.plc4x.java.bacnetip.ede.EdeParser;
import org.apache.plc4x.java.bacnetip.ede.model.Datapoint;
import org.apache.plc4x.java.bacnetip.ede.model.EdeModel;
import org.apache.plc4x.java.bacnetip.tag.BacNetIpTag;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.time.Instant;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

public class BacNetIpProtocolLogic extends Plc4xProtocolBase<BVLC> implements HasConfiguration<BacNetIpConfiguration>, PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(BacNetIpProtocolLogic.class);

    private EdeModel edeModel;

    private final Map<Integer, Consumer<PlcSubscriptionEvent>> consumerIdMap = new ConcurrentHashMap<>();

    @Override
    public void setConfiguration(BacNetIpConfiguration configuration) {
        if (configuration.getEdeFilePath() != null) {
            File edeFile = new File(configuration.getEdeFilePath());
            if (!edeFile.exists() || !edeFile.isFile()) {
                throw new PlcRuntimeException(String.format(
                    "File specified with 'ede-file-path' does not exist or is not a file: '%s'",
                    configuration.getEdeFilePath()));
            }
            edeModel = new EdeParser().parseFile(edeFile);
        } else if (configuration.getEdeDirectoryPath() != null) {
            File edeDirectory = new File(configuration.getEdeDirectoryPath());
            if (!edeDirectory.exists() || !edeDirectory.isDirectory()) {
                throw new PlcRuntimeException(String.format(
                    "File specified with 'ede-directory-path' does not exist or is not a directory: '%s'",
                    configuration.getEdeDirectoryPath()));
            }
            edeModel = new EdeParser().parseDirectory(edeDirectory);
        }
    }

    @Override
    public void onConnect(ConversationContext<BVLC> context) {
        if (context.isPassive()) {
            context.fireConnected();
        } else {
            throw new PlcRuntimeException("Active connections not yet supported");
        }
    }

    @Override
    public void close(ConversationContext<BVLC> context) {
        // Nothing to do here ...
    }

    @Override
    protected void decode(ConversationContext<BVLC> context, BVLC msg) throws Exception {
        NPDU npdu = null;
        if (msg instanceof BVLCOriginalUnicastNPDU) {
            BVLCOriginalUnicastNPDU bvlcOriginalUnicastNPDU = (BVLCOriginalUnicastNPDU) msg;
            npdu = bvlcOriginalUnicastNPDU.getNpdu();
        } else if (msg instanceof BVLCForwardedNPDU) {
            BVLCForwardedNPDU bvlcForwardedNPDU = (BVLCForwardedNPDU) msg;
            npdu = bvlcForwardedNPDU.getNpdu();
        } else if (msg instanceof BVLCOriginalBroadcastNPDU) {
            BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) msg;
            npdu = bvlcOriginalBroadcastNPDU.getNpdu();
        }

        if (npdu == null) {
            LOGGER.warn("Ummapped BVLC {}", msg);
            return;
        }

        if (npdu.getApdu() instanceof APDUConfirmedRequest) {
            APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
            decodeConfirmedRequest(apduConfirmedRequest);
        } else if (npdu.getApdu() instanceof APDUUnconfirmedRequest) {
            APDUUnconfirmedRequest unconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
            decodeUnconfirmedRequest(unconfirmedRequest);
        } else if (npdu.getApdu() instanceof APDUError) {
            APDUError apduError = (APDUError) npdu.getApdu();
        } else if (npdu.getApdu() instanceof APDUSimpleAck) {
            // Ignore this ...
        } else if (npdu.getApdu() instanceof APDUComplexAck) {
            // Ignore this ...
        } else if ((npdu.getApdu() == null) && (npdu.getNlm() != null)) {
            // "Who is router?" & "I am router" messages.
            // Ignore this ...
        } else {
            LOGGER.debug(String.format("Unexpected NPDU type: %s", npdu.getClass().getName()));
        }
    }

    private void decodeConfirmedRequest(APDUConfirmedRequest apduConfirmedRequest) {
        final BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
        if (serviceRequest instanceof BACnetConfirmedServiceRequestConfirmedCOVNotification) { // A value change subscription event.
            BACnetConfirmedServiceRequestConfirmedCOVNotification valueChange =
                (BACnetConfirmedServiceRequestConfirmedCOVNotification) serviceRequest;

            long deviceIdentifier = valueChange.getMonitoredObjectIdentifier().getPayload().getInstanceNumber();
            int objectType = valueChange.getMonitoredObjectIdentifier().getPayload().getObjectType().getValue();
            long objectInstance = valueChange.getMonitoredObjectIdentifier().getPayload().getInstanceNumber();
            BacNetIpTag curTag = new BacNetIpTag(deviceIdentifier, objectType, objectInstance);

            // The actual value change is in the notifications ... iterate through them to get it.
            for (BACnetPropertyValue baCnetPropertyValue : valueChange.getListOfValues().getData()) {
                // These are value change notifications. Ignore the rest.
                if (baCnetPropertyValue.getPropertyIdentifier().getValue() == BACnetPropertyIdentifier.PRESENT_VALUE) {
                    BACnetConstructedDataElement propertyValue = baCnetPropertyValue.getPropertyValue();

                    // Initialize an enriched version of the PlcStruct.
                    final Map<String, PlcValue> enrichedPlcValue = new HashMap<>();
                    enrichedPlcValue.put("deviceIdentifier", new PlcUDINT(deviceIdentifier));
                    enrichedPlcValue.put("objectType", new PlcDINT(objectType));
                    enrichedPlcValue.put("objectInstance", new PlcUDINT(objectInstance));
                    enrichedPlcValue.put("address", new PlcSTRING(toString(curTag)));

                    // From the original BACNet tag
                    enrichedPlcValue.put("tagNumber", PlcValueHandler.of(propertyValue.getPeekedTagHeader().getActualTagNumber()));
                    enrichedPlcValue.put("lengthValueType", PlcValueHandler.of(propertyValue.getPeekedTagHeader().getActualLength()));

                    // Use the information in the edeModel to enrich the information.
                    if (edeModel != null) {
                        final Datapoint datapoint = edeModel.getDatapoint(curTag);
                        if (datapoint != null) {
                            // Add all the attributes from the ede file.
                            enrichedPlcValue.putAll(datapoint.toPlcValues());
                        }
                    }
                    // Send out the enriched event.
                    publishEvent(curTag, new PlcStruct(enrichedPlcValue));
                }
            }
        } else if (serviceRequest instanceof BACnetConfirmedServiceRequestReadProperty) { // Someone read a value.
            // Ignore this ...
        } else if (serviceRequest instanceof BACnetConfirmedServiceRequestWriteProperty) { // Someone wrote a value.
            // Ignore this ...
        } else if (serviceRequest instanceof BACnetConfirmedServiceRequestSubscribeCOV) {
            // Ignore this ...
        } else {
            LOGGER.debug(String.format("Unexpected ConfirmedServiceRequest type: %s", serviceRequest.getClass().getName()));
        }
    }

    private void decodeUnconfirmedRequest(APDUUnconfirmedRequest unconfirmedRequest) {
        final BACnetUnconfirmedServiceRequest serviceRequest = unconfirmedRequest.getServiceRequest();
        if (serviceRequest instanceof BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) { // A value change subscription event.
            BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification valueChange =
                (BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) serviceRequest;

            long deviceIdentifier = valueChange.getMonitoredObjectIdentifier().getPayload().getInstanceNumber();
            int objectType = valueChange.getMonitoredObjectIdentifier().getPayload().getObjectType().getValue();
            long objectInstance = valueChange.getMonitoredObjectIdentifier().getPayload().getInstanceNumber();
            BacNetIpTag curTag = new BacNetIpTag(deviceIdentifier, objectType, objectInstance);

            // The actual value change is in the notifications ... iterate through them to get it.
            for (BACnetPropertyValue baCnetPropertyValue : valueChange.getListOfValues().getData()) {
                // These are value change notifications. Ignore the rest.
                if (baCnetPropertyValue.getPropertyIdentifier().getValue() == BACnetPropertyIdentifier.PRESENT_VALUE) {
                    BACnetApplicationTag baCnetTag = ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();

                    // Initialize an enriched version of the PlcStruct.
                    final Map<String, PlcValue> enrichedPlcValue = new HashMap<>();
                    enrichedPlcValue.put("deviceIdentifier", new PlcUDINT(deviceIdentifier));
                    enrichedPlcValue.put("objectType", new PlcDINT(objectType));
                    enrichedPlcValue.put("objectInstance", new PlcUDINT(objectInstance));
                    enrichedPlcValue.put("address", new PlcSTRING(toString(curTag)));

                    // From the original BACNet tag
                    enrichedPlcValue.put("tagNumber", PlcValueHandler.of(baCnetTag.getActualTagNumber()));
                    enrichedPlcValue.put("lengthValueType", PlcValueHandler.of(baCnetTag.getActualLength()));

                    // Use the information in the edeModel to enrich the information.
                    if (edeModel != null) {
                        final Datapoint datapoint = edeModel.getDatapoint(curTag);
                        if (datapoint != null) {
                            // Add all the attributes from the ede file.
                            enrichedPlcValue.putAll(datapoint.toPlcValues());
                        }
                    }
                    // Send out the enriched event.
                    publishEvent(curTag, new PlcStruct(enrichedPlcValue));
                }
            }
        } else if (serviceRequest instanceof BACnetUnconfirmedServiceRequestWhoHas) {
            // Ignore this ...
        } else if (serviceRequest instanceof BACnetUnconfirmedServiceRequestWhoIs) {
            // Ignore this ...
        } else if (serviceRequest instanceof BACnetUnconfirmedServiceRequestIAm) {
            // Ignore this ...
        } else if (serviceRequest instanceof BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer) {
            // Ignore this ...
        } else {
            LOGGER.debug(String.format("Unexpected UnconfirmedServiceRequest type: %s", serviceRequest.getClass().getName()));
        }
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        Map<String, ResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
        for (String tagName : subscriptionRequest.getTagNames()) {
            values.put(tagName, new ResponseItem<>(PlcResponseCode.OK, new DefaultPlcSubscriptionHandle(this)));
        }
        return CompletableFuture.completedFuture(new DefaultPlcSubscriptionResponse(subscriptionRequest, values));
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> collection) {
        final DefaultPlcConsumerRegistration consumerRegistration =
            new DefaultPlcConsumerRegistration(this, consumer, collection.toArray(new PlcSubscriptionHandle[0]));
        consumerIdMap.put(consumerRegistration.getConsumerId(), consumer);
        return consumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration plcConsumerRegistration) {
        DefaultPlcConsumerRegistration consumerRegistration = (DefaultPlcConsumerRegistration) plcConsumerRegistration;
        consumerIdMap.remove(consumerRegistration.getConsumerId());
    }

    protected void publishEvent(BacNetIpTag tag, PlcValue plcValue) {
        // Create a subscription event from the input.
        final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
            Collections.singletonMap("event", new ResponseItem<>(PlcResponseCode.OK, plcValue)));

        // Send the subscription event to all listeners.
        for (Consumer<PlcSubscriptionEvent> consumer : consumerIdMap.values()) {
            // TODO: Check if the subscription matches the current tag ..
            consumer.accept(event);
        }
    }

    private String toString(BacNetIpTag tag) {
        return tag.getDeviceIdentifier() + "/" + tag.getObjectType() + "/" + tag.getObjectInstance();
    }

}
