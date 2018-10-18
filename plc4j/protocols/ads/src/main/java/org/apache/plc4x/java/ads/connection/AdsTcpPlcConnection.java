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
package org.apache.plc4x.java.ads.connection;

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.model.*;
import org.apache.plc4x.java.ads.protocol.Ads2PayloadProtocol;
import org.apache.plc4x.java.ads.protocol.Payload2TcpProtocol;
import org.apache.plc4x.java.ads.protocol.Plc4x2AdsProtocol;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.base.model.InternalPlcConsumerRegistration;
import org.apache.plc4x.java.base.model.InternalPlcSubscriptionHandle;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
import org.apache.plc4x.java.base.protocol.SingleItemToSingleRequestProtocol;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.time.temporal.ChronoUnit;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;
import java.util.stream.Collectors;

public class AdsTcpPlcConnection extends AdsAbstractPlcConnection implements PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsTcpPlcConnection.class);

    private static final int TCP_PORT = 48898;

    private static final long ADD_DEVICE_TIMEOUT = CONF.getLong("plc4x.adsconnection.add.device,timeout", 3000);
    private static final long DEL_DEVICE_TIMEOUT = CONF.getLong("plc4x.adsconnection.del.device,timeout", 3000);

    private static AtomicInteger localPorts = new AtomicInteger(30000);

    private Map<InternalPlcConsumerRegistration, Consumer<AdsDeviceNotificationRequest>> consumerRegistrations = new HashMap<>();

    private AdsTcpPlcConnection(InetAddress address, AmsNetId targetAmsNetId, AmsPort targetAmsPort) {
        this(address, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    private AdsTcpPlcConnection(InetAddress address, Integer port, AmsNetId targetAmsNetId, AmsPort targetAmsPort) {
        this(address, port, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    private AdsTcpPlcConnection(InetAddress address, AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort) {
        this(address, null, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    private AdsTcpPlcConnection(InetAddress address, Integer port, AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort) {
        super(new TcpSocketChannelFactory(address, port != null ? port : TCP_PORT), targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    public static AdsTcpPlcConnection of(InetAddress address, AmsNetId targetAmsNetId, AmsPort targetAmsPort) {
        return new AdsTcpPlcConnection(address, targetAmsNetId, targetAmsPort);
    }

    public static AdsTcpPlcConnection of(InetAddress address, Integer port, AmsNetId targetAmsNetId, AmsPort targetAmsPort) {
        return new AdsTcpPlcConnection(address, port, targetAmsNetId, targetAmsPort);
    }

    public static AdsTcpPlcConnection of(InetAddress address, AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort) {
        return new AdsTcpPlcConnection(address, null, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    public static AdsTcpPlcConnection of(InetAddress address, Integer port, AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort) {
        return new AdsTcpPlcConnection(address, port, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the ads protocol.
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new Payload2TcpProtocol());
                pipeline.addLast(new Ads2PayloadProtocol());
                pipeline.addLast(new Plc4x2AdsProtocol(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, fieldMapping));
                pipeline.addLast(new SingleItemToSingleRequestProtocol(AdsTcpPlcConnection.this, AdsTcpPlcConnection.this, timer));
            }
        };
    }

    public InetAddress getRemoteAddress() {
        return ((TcpSocketChannelFactory) channelFactory).getAddress();
    }

    protected static AmsNetId generateAMSNetId() {
        try {
            return AmsNetId.of(Inet4Address.getLocalHost().getHostAddress() + ".1.1");
        } catch (UnknownHostException e) {
            throw new PlcRuntimeException(e);
        }
    }

    protected static AmsPort generateAMSPort() {
        return AmsPort.of(localPorts.getAndIncrement());
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest plcSubscriptionRequest) {
        InternalPlcSubscriptionRequest internalPlcSubscriptionRequest = checkInternal(plcSubscriptionRequest, InternalPlcSubscriptionRequest.class);
        // TODO: Make this multi-value
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        if (internalPlcSubscriptionRequest.getNumberOfFields() != 1) {
            throw new PlcNotImplementedException("Multirequest on subscribe not implemented yet");
        }

        SubscriptionPlcField subscriptionPlcField = internalPlcSubscriptionRequest.getSubscriptionFields().get(0);
        PlcField field = subscriptionPlcField.getPlcField();

        IndexGroup indexGroup;
        IndexOffset indexOffset;
        AdsDataType adsDataType;
        int numberOfElements;
        // If this is a symbolic field, it has to be resolved first.
        // TODO: This is blocking, should be changed to be async.
        if (field instanceof SymbolicAdsField) {
            mapFields((SymbolicAdsField) field);
            DirectAdsField directAdsField = fieldMapping.get(field);
            if (directAdsField == null) {
                throw new PlcRuntimeException("Unresolvable field " + field);
            }
            indexGroup = IndexGroup.of(directAdsField.getIndexGroup());
            indexOffset = IndexOffset.of(directAdsField.getIndexOffset());
            adsDataType = directAdsField.getAdsDataType();
            numberOfElements = directAdsField.getNumberOfElements();
        }
        // If it's no symbolic field, we can continue immediately
        // without having to do any resolving.
        else if (field instanceof DirectAdsField) {
            DirectAdsField directAdsField = (DirectAdsField) field;
            indexGroup = IndexGroup.of(directAdsField.getIndexGroup());
            indexOffset = IndexOffset.of(directAdsField.getIndexOffset());
            adsDataType = directAdsField.getAdsDataType();
            numberOfElements = directAdsField.getNumberOfElements();
        } else {
            throw new IllegalArgumentException("Unsupported field type " + field.getClass());
        }

        final TransmissionMode transmissionMode;
        long cycleTime = 4000000;
        switch (subscriptionPlcField.getPlcSubscriptionType()) {
            case CYCLIC:
                transmissionMode = TransmissionMode.DefinedValues.ADSTRANS_SERVERCYCLE;
                cycleTime = subscriptionPlcField.getDuration().orElseThrow(IllegalStateException::new).get(ChronoUnit.MILLIS);
                break;
            case CHANGE_OF_STATE:
                transmissionMode = TransmissionMode.DefinedValues.ADSTRANS_SERVERONCHA;
                break;
            default:
                throw new PlcRuntimeException("Unmapped type " + subscriptionPlcField.getPlcSubscriptionType());
        }

        // Prepare the subscription request itself.
        AdsAddDeviceNotificationRequest adsAddDeviceNotificationRequest = AdsAddDeviceNotificationRequest.of(
            targetAmsNetId,
            targetAmsPort,
            sourceAmsNetId,
            sourceAmsPort,
            Invoke.NONE,
            indexGroup,
            indexOffset,
            Length.of(adsDataType.getTargetByteSize() * (long) numberOfElements),
            transmissionMode,
            MaxDelay.of(0),
            CycleTime.of(cycleTime)
        );

        // Send the request to the plc and wait for a response
        // TODO: This is blocking, should be changed to be async.
        CompletableFuture<InternalPlcProprietaryResponse<AdsAddDeviceNotificationResponse>> addDeviceFuture = new CompletableFuture<>();
        channel.writeAndFlush(new PlcRequestContainer<>(new DefaultPlcProprietaryRequest<>(adsAddDeviceNotificationRequest), addDeviceFuture));
        InternalPlcProprietaryResponse<AdsAddDeviceNotificationResponse> addDeviceResponse = getFromFuture(addDeviceFuture, ADD_DEVICE_TIMEOUT);
        AdsAddDeviceNotificationResponse response = addDeviceResponse.getResponse();

        // Abort if we got anything but a successful response.
        if (response.getResult().toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            throw new PlcRuntimeException("Error code received " + response.getResult());
        }
        AdsSubscriptionHandle adsSubscriptionHandle = new AdsSubscriptionHandle(this, response.getNotificationHandle());

        Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> responseItems = internalPlcSubscriptionRequest.getFieldNames()
            .stream()
            .collect(Collectors.toMap(
                fieldName -> fieldName,
                ignored -> Pair.of(PlcResponseCode.OK, adsSubscriptionHandle)
            ));

        future.complete(new DefaultPlcSubscriptionResponse(internalPlcSubscriptionRequest, responseItems));
        return future;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest plcUnsubscriptionRequest) {
        InternalPlcUnsubscriptionRequest internalPlcUnsubscriptionRequest = checkInternal(plcUnsubscriptionRequest, InternalPlcUnsubscriptionRequest.class);
        for (InternalPlcSubscriptionHandle internalPlcSubscriptionHandle : internalPlcUnsubscriptionRequest.getInternalPlcSubscriptionHandles()) {
            if (internalPlcSubscriptionHandle instanceof AdsSubscriptionHandle) {
                AdsSubscriptionHandle adsSubscriptionHandle = (AdsSubscriptionHandle) internalPlcSubscriptionHandle;
                AdsDeleteDeviceNotificationRequest adsDeleteDeviceNotificationRequest =
                    AdsDeleteDeviceNotificationRequest.of(
                        targetAmsNetId,
                        targetAmsPort,
                        sourceAmsNetId,
                        sourceAmsPort,
                        Invoke.NONE,
                        adsSubscriptionHandle.getNotificationHandle()
                    );
                CompletableFuture<InternalPlcProprietaryResponse<AdsDeleteDeviceNotificationResponse>> deleteDeviceFuture =
                    new CompletableFuture<>();
                channel.writeAndFlush(new PlcRequestContainer<>(new DefaultPlcProprietaryRequest<>(adsDeleteDeviceNotificationRequest), deleteDeviceFuture));

                InternalPlcProprietaryResponse<AdsDeleteDeviceNotificationResponse> deleteDeviceResponse =
                    getFromFuture(deleteDeviceFuture, DEL_DEVICE_TIMEOUT);
                AdsDeleteDeviceNotificationResponse response = deleteDeviceResponse.getResponse();
                if (response.getResult().toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
                    throw new PlcRuntimeException("Non error code received " + response.getResult());
                }
            }
        }
        CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<>();
        future.complete(new DefaultPlcUnsubscriptionResponse());
        return future;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        return register(consumer, handles.toArray(new PlcSubscriptionHandle[0]));
    }

    public InternalPlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, PlcSubscriptionHandle... handles) {
        Objects.requireNonNull(consumer);
        Objects.requireNonNull(handles);
        InternalPlcSubscriptionHandle[] internalPlcSubscriptionHandles = new InternalPlcSubscriptionHandle[handles.length];
        for (int i = 0; i < handles.length; i++) {
            internalPlcSubscriptionHandles[i] = checkInternal(handles[i], InternalPlcSubscriptionHandle.class);
        }

        InternalPlcConsumerRegistration internalPlcConsumerRegistration = new DefaultPlcConsumerRegistration(this, consumer, internalPlcSubscriptionHandles);

        Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer =
            adsDeviceNotificationRequest -> adsDeviceNotificationRequest.getAdsStampHeaders().forEach(adsStampHeader -> {
                Calendar timeStamp = Calendar.getInstance();
                timeStamp.setTime(adsStampHeader.getTimeStamp().getAsDate());

                adsStampHeader.getAdsNotificationSamples()
                    .forEach(adsNotificationSample -> {
                        Data data = adsNotificationSample.getData();
                        try {
                            PlcSubscriptionEvent subscriptionEventItem = new DefaultPlcSubscriptionEvent(timeStamp, data.getBytes());
                            consumer.accept(subscriptionEventItem);
                        } catch (RuntimeException e) {
                            LOGGER.error("Can't decode {}", data, e);
                        }
                    });
            });

        // Store the reference for so it can be uses for later
        consumerRegistrations.put(internalPlcConsumerRegistration, adsDeviceNotificationRequestConsumer);
        // register the actual consumer.
        getChannel().pipeline().get(Plc4x2AdsProtocol.class).addConsumer(adsDeviceNotificationRequestConsumer);

        return internalPlcConsumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration plcConsumerRegistration) {
        InternalPlcConsumerRegistration internalPlcConsumerRegistration = checkInternal(plcConsumerRegistration, InternalPlcConsumerRegistration.class);
        Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer = consumerRegistrations.remove(internalPlcConsumerRegistration);
        if (adsDeviceNotificationRequestConsumer == null) {
            return;
        }
        getChannel().pipeline().get(Plc4x2AdsProtocol.class).removeConsumer(adsDeviceNotificationRequestConsumer);
    }

    @Override
    public Optional<PlcSubscriptionRequest.Builder> subscriptionRequestBuilder() {
        return Optional.of(new DefaultPlcSubscriptionRequest.Builder(this, new AdsPlcFieldHandler()));
    }

    @Override
    public Optional<PlcUnsubscriptionRequest.Builder> unsubscriptionRequestBuilder() {
        return Optional.of(new DefaultPlcUnsubscriptionRequest.Builder(this));
    }

    @Override
    public void close() throws PlcConnectionException {
        try {
            consumerRegistrations.values().forEach(getChannel().pipeline().get(Plc4x2AdsProtocol.class)::removeConsumer);
            List<PlcSubscriptionHandle> collect = consumerRegistrations.keySet().stream()
                .map(InternalPlcConsumerRegistration::getAssociatedHandles)
                .flatMap(Collection::stream)
                .map(PlcSubscriptionHandle.class::cast)
                .collect(Collectors.toList());

            PlcUnsubscriptionRequest plcUnsubscriptionRequest = new DefaultPlcUnsubscriptionRequest.Builder(this).addHandles(collect).build();
            unsubscribe(plcUnsubscriptionRequest).get(5, TimeUnit.SECONDS);

            consumerRegistrations.clear();
        } catch (InterruptedException e) {
            LOGGER.warn("Exception while closing", e);
            Thread.currentThread().interrupt();
        } catch (RuntimeException | ExecutionException | TimeoutException e) {
            LOGGER.warn("Exception while closing", e);
        }
        super.close();
    }
}
