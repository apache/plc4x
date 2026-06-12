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
package org.apache.plc4x.java.bacnetip;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.bacnetip.configuration.BacNetIpConfiguration;
import org.apache.plc4x.java.bacnetip.ede.EdeParser;
import org.apache.plc4x.java.bacnetip.ede.model.Datapoint;
import org.apache.plc4x.java.bacnetip.ede.model.EdeModel;
import org.apache.plc4x.java.bacnetip.readwrite.APDUConfirmedRequest;
import org.apache.plc4x.java.bacnetip.readwrite.APDUError;
import org.apache.plc4x.java.bacnetip.readwrite.APDUUnconfirmedRequest;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetApplicationTag;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequest;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequestConfirmedCOVNotification;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConstructedDataElement;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetConstructedDataUnspecified;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetPropertyIdentifier;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetPropertyValue;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequest;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.bacnetip.readwrite.BVLCForwardedNPDU;
import org.apache.plc4x.java.bacnetip.readwrite.BVLCOriginalBroadcastNPDU;
import org.apache.plc4x.java.bacnetip.readwrite.BVLCOriginalUnicastNPDU;
import org.apache.plc4x.java.bacnetip.readwrite.NPDU;
import org.apache.plc4x.java.bacnetip.tag.BacNetIpTag;
import org.apache.plc4x.java.bacnetip.tag.BacNetIpTagHandler;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcUnsubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.spi.values.PlcUDINT;
import org.apache.plc4x.java.spi.values.PlcUSINT;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.time.Instant;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Map.Entry;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

/**
 * BACnet/IP is a passive listener — there's no handshake, no read or write, and
 * subscriptions don't filter on a specific device (no working address parser yet).
 * The connection just receives BVLC frames, extracts COV notifications, and
 * broadcasts an enriched {@link PlcStruct} event to every registered consumer.
 */
public class BacNetIpConnection extends ConnectionBase<BacNetIpConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(BacNetIpConnection.class);

    private BacNetIpMessageCodec messageCodec;
    private volatile boolean connected = false;
    private EdeModel edeModel;

    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    public BacNetIpConnection(BacNetIpConfiguration configuration,
                              TransportInstance<?> transportInstance,
                              AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new BacNetIpTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    public boolean isConnected() {
        return connected && messageCodec != null && messageCodec.isOpen();
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        if (configuration.getEdeFile() != null) {
            File edeFile = configuration.getEdeFile();
            if (!edeFile.exists() || !edeFile.isFile()) {
                throw new PlcConnectionException("File specified with 'ede-file-path' does not exist or is not a file: '"
                    + edeFile + "'");
            }
            edeModel = new EdeParser().parseFile(edeFile);
        } else if (configuration.getEdeDirectory() != null) {
            File edeDirectory = configuration.getEdeDirectory();
            if (!edeDirectory.exists() || !edeDirectory.isDirectory()) {
                throw new PlcConnectionException("File specified with 'ede-directory-path' does not exist or is not a directory: '"
                    + edeDirectory + "'");
            }
            edeModel = new EdeParser().parseDirectory(edeDirectory);
        }

        messageCodec = new BacNetIpMessageCodec(transportInstance, this::handleIncomingMessage);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming BACnet/IP data", e);
            }
        });
        connected = true;
    }

    @Override
    public void close() throws Exception {
        connected = false;
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        consumers.clear();
        super.close();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding pipeline
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private void handleIncomingMessage(BVLC msg) {
        NPDU npdu = null;
        if (msg instanceof BVLCOriginalUnicastNPDU bvlcUnicast) {
            npdu = bvlcUnicast.getNpdu();
        } else if (msg instanceof BVLCForwardedNPDU bvlcForwarded) {
            npdu = bvlcForwarded.getNpdu();
        } else if (msg instanceof BVLCOriginalBroadcastNPDU bvlcBroadcast) {
            npdu = bvlcBroadcast.getNpdu();
        }

        if (npdu == null) {
            LOGGER.warn("Unmapped BVLC {}", msg);
            return;
        }

        if (npdu.getApdu() instanceof APDUConfirmedRequest apduConfirmedRequest) {
            decodeConfirmedRequest(apduConfirmedRequest);
        } else if (npdu.getApdu() instanceof APDUUnconfirmedRequest apduUnconfirmedRequest) {
            decodeUnconfirmedRequest(apduUnconfirmedRequest);
        } else if (npdu.getApdu() instanceof APDUError) {
            // Ignore.
        } else if (npdu.getApdu() == null && npdu.getNlm() != null) {
            // "Who is router?" / "I am router" messages — ignore.
        } else {
            LOGGER.debug("Unexpected NPDU APDU type: {}", npdu.getApdu() == null ? "null" : npdu.getApdu().getClass().getName());
        }
    }

    private void decodeConfirmedRequest(APDUConfirmedRequest apduConfirmedRequest) {
        BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
        if (!(serviceRequest instanceof BACnetConfirmedServiceRequestConfirmedCOVNotification valueChange)) {
            return;
        }
        long deviceIdentifier = valueChange.getMonitoredObjectIdentifier().getPayload().getInstanceNumber();
        int objectType = valueChange.getMonitoredObjectIdentifier().getPayload().getObjectType().getValue();
        long objectInstance = valueChange.getMonitoredObjectIdentifier().getPayload().getInstanceNumber();
        BacNetIpTag curTag = new BacNetIpTag(deviceIdentifier, objectType, objectInstance);

        for (BACnetPropertyValue baCnetPropertyValue : valueChange.getListOfValues().getData()) {
            if (baCnetPropertyValue.getPropertyIdentifier().getValue() != BACnetPropertyIdentifier.PRESENT_VALUE) {
                continue;
            }
            BACnetConstructedDataElement propertyValue = baCnetPropertyValue.getPropertyValue();
            Map<String, PlcValue> enriched = baseFields(curTag, deviceIdentifier, objectType, objectInstance);
            enriched.put("tagNumber", new PlcUSINT(propertyValue.getPeekedTagHeader().getActualTagNumber()));
            enriched.put("lengthValueType", new PlcUSINT(propertyValue.getPeekedTagHeader().getActualLength()));
            applyEdeEnrichment(curTag, enriched);
            publishEvent(curTag, new PlcStruct(enriched));
        }
    }

    private void decodeUnconfirmedRequest(APDUUnconfirmedRequest unconfirmedRequest) {
        BACnetUnconfirmedServiceRequest serviceRequest = unconfirmedRequest.getServiceRequest();
        if (!(serviceRequest instanceof BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification valueChange)) {
            return;
        }
        long deviceIdentifier = valueChange.getMonitoredObjectIdentifier().getPayload().getInstanceNumber();
        int objectType = valueChange.getMonitoredObjectIdentifier().getPayload().getObjectType().getValue();
        long objectInstance = valueChange.getMonitoredObjectIdentifier().getPayload().getInstanceNumber();
        BacNetIpTag curTag = new BacNetIpTag(deviceIdentifier, objectType, objectInstance);

        for (BACnetPropertyValue baCnetPropertyValue : valueChange.getListOfValues().getData()) {
            if (baCnetPropertyValue.getPropertyIdentifier().getValue() != BACnetPropertyIdentifier.PRESENT_VALUE) {
                continue;
            }
            BACnetApplicationTag baCnetTag = ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData())
                .getData().get(0).getApplicationTag();
            Map<String, PlcValue> enriched = baseFields(curTag, deviceIdentifier, objectType, objectInstance);
            enriched.put("tagNumber", new PlcUSINT(baCnetTag.getActualTagNumber()));
            enriched.put("lengthValueType", new PlcUSINT(baCnetTag.getActualLength()));
            applyEdeEnrichment(curTag, enriched);
            publishEvent(curTag, new PlcStruct(enriched));
        }
    }

    private Map<String, PlcValue> baseFields(BacNetIpTag tag, long deviceIdentifier, int objectType, long objectInstance) {
        Map<String, PlcValue> enriched = new LinkedHashMap<>();
        enriched.put("deviceIdentifier", new PlcUDINT(deviceIdentifier));
        enriched.put("objectType", new PlcDINT(objectType));
        enriched.put("objectInstance", new PlcUDINT(objectInstance));
        enriched.put("address", new PlcSTRING(tag.getDeviceIdentifier() + "/" + tag.getObjectType() + "/" + tag.getObjectInstance()));
        return enriched;
    }

    private void applyEdeEnrichment(BacNetIpTag tag, Map<String, PlcValue> enriched) {
        if (edeModel == null) {
            return;
        }
        Datapoint datapoint = edeModel.getDatapoint(tag);
        if (datapoint != null) {
            enriched.putAll(datapoint.toPlcValues());
        }
    }

    private void publishEvent(BacNetIpTag tag, PlcValue plcValue) {
        DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
            Collections.singletonMap("event", new DefaultPlcResponseItem<>(PlcResponseCode.OK, plcValue)));
        for (Consumer<PlcSubscriptionEvent> consumer : consumers.values()) {
            // No tag-level filtering yet — the address pattern is still a TODO.
            consumer.accept(event);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe / consumer registry
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest request) {
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
        for (String tagName : request.getTagNames()) {
            values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, new BacNetIpSubscriptionHandle(this)));
        }
        return CompletableFuture.completedFuture(new DefaultPlcSubscriptionResponse(request, values));
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest request) {
        return CompletableFuture.completedFuture(new DefaultPlcUnsubscriptionResponse(request));
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer,
                                                        Collection<PlcSubscriptionHandle> handles) {
        DefaultPlcConsumerRegistration registration = new DefaultPlcConsumerRegistration(this, consumer,
            handles.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(registration, consumer);
        return registration;
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        if (registration instanceof DefaultPlcConsumerRegistration r) {
            consumers.remove(r);
        }
    }

}
