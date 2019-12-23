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
package org.apache.plc4x.java.amsads.connection;

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.amsads.model.*;
import org.apache.plc4x.java.amsads.protocol.Ads2PayloadProtocol;
import org.apache.plc4x.java.amsads.protocol.Payload2TcpProtocol;
import org.apache.plc4x.java.amsads.protocol.Plc4x2AdsProtocol;
import org.apache.plc4x.java.amsads.protocol.util.LittleEndianDecoder;
import org.apache.plc4x.java.amsads.readwrite.*;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.InternalPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.InternalPlcSubscriptionHandle;
import org.apache.plc4x.java.spi.model.SubscriptionPlcField;
import org.apache.plc4x.java.spi.protocol.SingleItemToSingleRequestProtocol;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.time.Duration;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.stream.Collectors;

public class AdsTcpPlcConnection extends AdsAbstractPlcConnection implements PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsTcpPlcConnection.class);

    private static final int TCP_PORT = 48898;

    private static final long ADD_DEVICE_TIMEOUT = CONF.getLong("plc4x.adsconnection.add.device,timeout", 3000);
    private static final long DEL_DEVICE_TIMEOUT = CONF.getLong("plc4x.adsconnection.del.device,timeout", 3000);

    private static AtomicInteger localPorts = new AtomicInteger(30000);

    private Map<InternalPlcConsumerRegistration, Consumer<AdsDeviceNotificationRequest>> consumerRegistrations = new HashMap<>();

    private AdsTcpPlcConnection(InetAddress address, AmsNetId targetAmsNetId, int targetint) {
        this(address, targetAmsNetId, targetint, generateAmsNetId(), generateAmsPort());
    }

    private AdsTcpPlcConnection(InetAddress address, Integer port, AmsNetId targetAmsNetId, int targetint) {
        this(address, port, targetAmsNetId, targetint, generateAmsNetId(), generateAmsPort());
    }

    private AdsTcpPlcConnection(InetAddress address, AmsNetId targetAmsNetId, int targetint, AmsNetId sourceAmsNetId, int sourceint) {
        this(address, null, targetAmsNetId, targetint, sourceAmsNetId, sourceint);
    }

    private AdsTcpPlcConnection(InetAddress address, Integer port, AmsNetId targetAmsNetId, int targetint, AmsNetId sourceAmsNetId, int sourceint) {
        super(new TcpSocketChannelFactory(address, port != null ? port : TCP_PORT), targetAmsNetId, targetint, sourceAmsNetId, sourceint);
    }

    public static AdsTcpPlcConnection of(InetAddress address, AmsNetId targetAmsNetId, int targetint) {
        return new AdsTcpPlcConnection(address, targetAmsNetId, targetint);
    }

    public static AdsTcpPlcConnection of(InetAddress address, Integer port, AmsNetId targetAmsNetId, int targetint) {
        return new AdsTcpPlcConnection(address, port, targetAmsNetId, targetint);
    }

    public static AdsTcpPlcConnection of(InetAddress address, AmsNetId targetAmsNetId, int targetint, AmsNetId sourceAmsNetId, int sourceint) {
        return new AdsTcpPlcConnection(address, null, targetAmsNetId, targetint, sourceAmsNetId, sourceint);
    }

    public static AdsTcpPlcConnection of(InetAddress address, Integer port, AmsNetId targetAmsNetId, int targetint, AmsNetId sourceAmsNetId, int sourceint) {
        return new AdsTcpPlcConnection(address, port, targetAmsNetId, targetint, sourceAmsNetId, sourceint);
    }

    // TODO fix that
//    @Override
//    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
//        return new ChannelInitializer<Channel>() {
//            @Override
//            protected void initChannel(Channel channel) {
//                // Build the protocol stack for communicating with the ads protocol.
//                ChannelPipeline pipeline = channel.pipeline();
//                pipeline.addLast(new Payload2TcpProtocol());
//                pipeline.addLast(new Ads2PayloadProtocol());
//                pipeline.addLast(new Plc4x2AdsProtocol(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, fieldMapping));
//                pipeline.addLast(new SingleItemToSingleRequestProtocol(AdsTcpPlcConnection.this, AdsTcpPlcConnection.this, AdsTcpPlcConnection.this, timer, SingleItemToSingleRequestProtocol.SplitConfig.builder().dontSplitSubscribe().dontSplitUnsubscribe().build(), false));
//            }
//        };
//    }

    public InetAddress getRemoteAddress() {
        return ((TcpSocketChannelFactory) channelFactory).getAddress();
    }

    protected static AmsNetId generateAmsNetId() {
        try {
            String hostAddress = Inet4Address.getLocalHost().getHostAddress();
            String[] octets = hostAddress.split("\\.");
            return new AmsNetId(
                Short.parseShort(octets[3]),
                Short.parseShort(octets[2]),
                Short.parseShort(octets[1]),
                Short.parseShort(octets[0]),
                (short) 1,
                (short) 2
            );
        } catch (UnknownHostException e) {
            throw new PlcRuntimeException(e);
        }
    }

    protected static int generateAmsPort() {
        return localPorts.getAndIncrement();
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest plcSubscriptionRequest) {
        InternalPlcSubscriptionRequest internalPlcSubscriptionRequest = checkInternal(plcSubscriptionRequest, InternalPlcSubscriptionRequest.class);
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();

        Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> responseItems = internalPlcSubscriptionRequest.getSubscriptionPlcFieldMap().entrySet().stream()
            .map(subscriptionPlcFieldEntry -> {
                final String plcFieldName = subscriptionPlcFieldEntry.getKey();
                final SubscriptionPlcField subscriptionPlcField = subscriptionPlcFieldEntry.getValue();
                final PlcField field = Objects.requireNonNull(subscriptionPlcField.getPlcField());

                final long indexGroup;
                final long indexOffset;
                final AdsDataType adsDataType;
                final int numberOfElements;
                // If this is a symbolic field, it has to be resolved first.
                // TODO: This is blocking, should be changed to be async.
                if (field instanceof SymbolicAdsField) {
                    mapFields((SymbolicAdsField) field);
                    DirectAdsField directAdsField = fieldMapping.get(field);
                    if (directAdsField == null) {
                        throw new PlcRuntimeException("Unresolvable field " + field);
                    }
                    indexGroup = directAdsField.getIndexGroup();
                    indexOffset = directAdsField.getIndexOffset();
                    adsDataType = directAdsField.getAdsDataType();
                    numberOfElements = directAdsField.getNumberOfElements();
                }
                // If it's no symbolic field, we can continue immediately
                // without having to do any resolving.
                else if (field instanceof DirectAdsField) {
                    DirectAdsField directAdsField = (DirectAdsField) field;
                    indexGroup = directAdsField.getIndexGroup();
                    indexOffset = directAdsField.getIndexOffset();
                    adsDataType = directAdsField.getAdsDataType();
                    numberOfElements = directAdsField.getNumberOfElements();
                } else {
                    throw new IllegalArgumentException("Unsupported field type " + field.getClass());
                }

                final long transmissionMode;
                long cycleTime = 4000000;
                switch (subscriptionPlcField.getPlcSubscriptionType()) {
                    case CYCLIC:
                        transmissionMode = 3L;
                        cycleTime = subscriptionPlcField.getDuration().orElse(Duration.ofSeconds(1)).toMillis();
                        break;
                    case CHANGE_OF_STATE:
                        transmissionMode = 4L;
                        break;
                    default:
                        throw new PlcRuntimeException("Unmapped type " + subscriptionPlcField.getPlcSubscriptionType());
                }

                // Prepare the subscription request itself.
                AdsAddDeviceNotificationRequest adsAddDeviceNotificationRequest = new AdsAddDeviceNotificationRequest(
                    indexGroup,
                    indexOffset,
                    adsDataType.getTargetByteSize() * (long) numberOfElements,
                    transmissionMode,
                    cycleTime + 1,
                    cycleTime
                );

                // Send the request to the plc and wait for a response
                // TODO: This is blocking, should be changed to be async.
                CompletableFuture<InternalPlcProprietaryResponse<AdsAddDeviceNotificationResponse>> addDeviceFuture = new CompletableFuture<>();
                channel.writeAndFlush(new PlcRequestContainer<>(new DefaultPlcProprietaryRequest<>(adsAddDeviceNotificationRequest), addDeviceFuture));
                InternalPlcProprietaryResponse<AdsAddDeviceNotificationResponse> addDeviceResponse = getFromFuture(addDeviceFuture, ADD_DEVICE_TIMEOUT);
                AdsAddDeviceNotificationResponse response = addDeviceResponse.getResponse();

                // Abort if we got anything but a successful response.
                if (response.getResult() != 0L) {
                    throw new PlcRuntimeException("Error code received " + response.getResult());
                }
                PlcSubscriptionHandle adsSubscriptionHandle = new AdsSubscriptionHandle(this, plcFieldName, adsDataType, response.getNotificationHandle());
                return Pair.of(plcFieldName, Pair.of(PlcResponseCode.OK, adsSubscriptionHandle));
            })
            .collect(Collectors.toMap(Pair::getKey, Pair::getValue));

        future.complete(new DefaultPlcSubscriptionResponse(internalPlcSubscriptionRequest, responseItems));
        return future;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest plcUnsubscriptionRequest) {
        InternalPlcUnsubscriptionRequest internalPlcUnsubscriptionRequest = checkInternal(plcUnsubscriptionRequest, InternalPlcUnsubscriptionRequest.class);
        for (InternalPlcSubscriptionHandle internalPlcSubscriptionHandle : internalPlcUnsubscriptionRequest.getInternalPlcSubscriptionHandles()) {
            if (internalPlcSubscriptionHandle instanceof AdsSubscriptionHandle) {
                AdsSubscriptionHandle adsSubscriptionHandle = (AdsSubscriptionHandle) internalPlcSubscriptionHandle;
                AdsDeleteDeviceNotificationRequest adsDeleteDeviceNotificationRequest = new AdsDeleteDeviceNotificationRequest(
                    adsSubscriptionHandle.getNotificationHandle()
                );
                CompletableFuture<InternalPlcProprietaryResponse<AdsDeleteDeviceNotificationResponse>> deleteDeviceFuture =
                    new CompletableFuture<>();
                channel.writeAndFlush(new PlcRequestContainer<>(new DefaultPlcProprietaryRequest<>(adsDeleteDeviceNotificationRequest), deleteDeviceFuture));

                InternalPlcProprietaryResponse<AdsDeleteDeviceNotificationResponse> deleteDeviceResponse =
                    getFromFuture(deleteDeviceFuture, DEL_DEVICE_TIMEOUT);
                AdsDeleteDeviceNotificationResponse response = deleteDeviceResponse.getResponse();
                if (response.getResult() != 0L) {
                    throw new PlcRuntimeException("Non error code received " + response.getResult());
                }
            }
        }
        CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<>();
        future.complete(new DefaultPlcUnsubscriptionResponse(internalPlcUnsubscriptionRequest));
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
        Map<Long, AdsSubscriptionHandle> notificationHandleAdsSubscriptionHandleMap = Arrays.stream(internalPlcSubscriptionHandles)
            .map(subscriptionHandle -> checkInternal(subscriptionHandle, AdsSubscriptionHandle.class))
            .collect(Collectors.toConcurrentMap(AdsSubscriptionHandle::getNotificationHandle, Function.identity()));

        Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer =
            adsDeviceNotificationRequest -> Arrays.asList(adsDeviceNotificationRequest.getAdsStampHeaders()).forEach(adsStampHeader -> {
                BigInteger winTime = adsStampHeader.getTimestamp();
                BigInteger timeMillisSince16010101 = winTime.divide(BigInteger.valueOf(10_000));
                BigInteger EPOCH_DIFF_IN_MILLIS = BigInteger.valueOf((369L * 365L + 89L) * 86400L * 1000L);
                BigInteger subtract = timeMillisSince16010101.subtract(EPOCH_DIFF_IN_MILLIS);
                Instant timeStamp = new Date(subtract.longValue()).toInstant();

                Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> fields = new HashMap<>();
                Arrays.asList(adsStampHeader.getAdsNotificationSamples())
                    .forEach(adsNotificationSample -> {
                        Long notificationHandle = adsNotificationSample.getNotificationHandle();
                        byte[] data = adsNotificationSample.getData();
                        AdsSubscriptionHandle adsSubscriptionHandle = notificationHandleAdsSubscriptionHandleMap.get(notificationHandle);
                        if (adsSubscriptionHandle == null) {
                            // TODO: we might want to refactor this so that we don't subscribe to everything in the first place.
                            // TODO: rather than we add a Consumer with the handle as key
                            LOGGER.trace("We are not interested in this sample {} with handle {}", adsNotificationSample, notificationHandle);
                            return;
                        }
                        String plcFieldName = adsSubscriptionHandle.getPlcFieldName();
                        AdsDataType adsDataType = adsSubscriptionHandle.getAdsDataType();
                        try {
                            BaseDefaultFieldItem baseDefaultFieldItem = LittleEndianDecoder.decodeData(adsDataType, data);
                            fields.put(plcFieldName, Pair.of(PlcResponseCode.OK, baseDefaultFieldItem));
                        } catch (RuntimeException e) {
                            LOGGER.error("Can't decode {}", data, e);
                        }

                    });
                try {
                    PlcSubscriptionEvent subscriptionEventItem = new DefaultPlcSubscriptionEvent(timeStamp, fields);
                    consumer.accept(subscriptionEventItem);
                } catch (RuntimeException e) {
                    LOGGER.error("Can't dispatch adsStampHeader {}", adsStampHeader, e);
                }
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
    public boolean canSubscribe() {
        return true;
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        return new DefaultPlcSubscriptionRequest.Builder(this, new AdsPlcFieldHandler());
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
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
