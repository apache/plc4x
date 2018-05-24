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
import org.apache.plc4x.java.ads.model.AdsAddress;
import org.apache.plc4x.java.ads.model.SymbolicAdsAddress;
import org.apache.plc4x.java.ads.protocol.Ads2PayloadProtocol;
import org.apache.plc4x.java.ads.protocol.Payload2TcpProtocol;
import org.apache.plc4x.java.ads.protocol.Plc4x2AdsProtocol;
import org.apache.plc4x.java.ads.protocol.util.LittleEndianDecoder;
import org.apache.plc4x.java.api.connection.PlcSubscriber;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.RequestItem;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class AdsTcpPlcConnection extends AdsAbstractPlcConnection implements PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsTcpPlcConnection.class);

    private static final int TCP_PORT = 48898;

    private static final long ADD_DEVICE_TIMEOUT = CONF.getLong("plc4x.adsconnection.add.device,timeout", 3000);
    private static final long DEL_DEVICE_TIMEOUT = CONF.getLong("plc4x.adsconnection.del.device,timeout", 3000);

    private static AtomicInteger localPorts = new AtomicInteger(30000);

    private final Map<Pair<Consumer<? extends PlcNotification>, Address>, Pair<Consumer<AdsDeviceNotificationRequest>, NotificationHandle>> subscriberMap = new HashMap<>();

    private final Map<NotificationHandle, Consumer<? extends PlcNotification>> handleConsumerMap = new HashMap<>();

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
                pipeline.addLast(new Plc4x2AdsProtocol(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, addressMapping));
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
    public <T extends R, R> void subscribe(Consumer<PlcNotification<R>> consumer, Address address, Class<T> dataType) {
        Objects.requireNonNull(consumer);
        Objects.requireNonNull(address);
        IndexGroup indexGroup;
        IndexOffset indexOffset;
        if (address instanceof SymbolicAdsAddress) {
            mapAddresses(new PlcRequest() {
                {
                    requestItems.add(new RequestItem(Void.class, address) {
                        // Not relevant
                    });
                }
            });
            AdsAddress adsAddress = addressMapping.get(address);
            if (adsAddress == null) {
                throw new PlcRuntimeException("Unresolvable address" + address);
            }
            indexGroup = IndexGroup.of(adsAddress.getIndexGroup());
            indexOffset = IndexOffset.of(adsAddress.getIndexOffset());
        } else if (address instanceof AdsAddress) {
            AdsAddress adsAddress = (AdsAddress) address;
            indexGroup = IndexGroup.of(adsAddress.getIndexGroup());
            indexOffset = IndexOffset.of(adsAddress.getIndexOffset());
        } else {
            throw new IllegalArgumentException("Unssuported address type " + address.getClass());
        }
        AdsAddDeviceNotificationRequest adsAddDeviceNotificationRequest = AdsAddDeviceNotificationRequest.of(
            targetAmsNetId,
            targetAmsPort,
            sourceAmsNetId,
            sourceAmsPort,
            Invoke.NONE,
            indexGroup,
            indexOffset,
            Length.of(1),
            TransmissionMode.of(3),
            MaxDelay.of(0),
            CycleTime.of(4000000)
        );

        CompletableFuture<PlcProprietaryResponse<AdsAddDeviceNotificationResponse>> addDeviceFuture = new CompletableFuture<>();
        channel.writeAndFlush(new PlcRequestContainer<>(new PlcProprietaryRequest<>(adsAddDeviceNotificationRequest), addDeviceFuture));
        PlcProprietaryResponse<AdsAddDeviceNotificationResponse> addDeviceResponse;
        try {
            addDeviceResponse = addDeviceFuture.get(ADD_DEVICE_TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException e) {
            LOGGER.warn("Interrupted!", e);
            Thread.currentThread().interrupt();
            throw new PlcRuntimeException(e);
        } catch (ExecutionException | TimeoutException e) {
            throw new PlcRuntimeException(e);
        }
        AdsAddDeviceNotificationResponse response = addDeviceResponse.getResponse();
        if (response.getResult().toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
            throw new PlcRuntimeException("Non error code received " + response.getResult());
        }
        NotificationHandle notificationHandle = response.getNotificationHandle();
        handleConsumerMap.put(notificationHandle, consumer);

        Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer =
            adsDeviceNotificationRequest -> adsDeviceNotificationRequest.getAdsStampHeaders().forEach(adsStampHeader -> {
                Date timeStamp = adsStampHeader.getTimeStamp().getAsDate();

                adsStampHeader.getAdsNotificationSamples()
                    .forEach(adsNotificationSample -> {
                        Consumer<? extends PlcNotification> plcNotificationConsumer = handleConsumerMap.get(adsNotificationSample.getNotificationHandle());
                        if (plcNotificationConsumer == null) {
                            LOGGER.warn("Unmapped notification received ", adsNotificationSample.getNotificationHandle());
                            return;
                        }
                        Data data = adsNotificationSample.getData();
                        try {
                            consumer.accept(new PlcNotification(timeStamp, address, LittleEndianDecoder.decodeData(dataType, data.getBytes())));
                        } catch (PlcProtocolException | RuntimeException e) {
                            LOGGER.error("Can't decode {}", data, e);
                        }
                    });
            });
        subscriberMap.put(Pair.of(consumer, address), Pair.of(adsDeviceNotificationRequestConsumer, notificationHandle));
        getChannel().pipeline().get(Plc4x2AdsProtocol.class).addConsumer(adsDeviceNotificationRequestConsumer);
    }

    @Override
    public <R> void unsubscribe(Consumer<PlcNotification<R>> consumer, Address address) {
        Pair<Consumer<AdsDeviceNotificationRequest>, NotificationHandle> handlePair = subscriberMap.remove(Pair.of(consumer, address));
        if (handlePair != null) {
            NotificationHandle notificationHandle = handlePair.getRight();
            handleConsumerMap.remove(notificationHandle);
            AdsDeleteDeviceNotificationRequest adsDeleteDeviceNotificationRequest = AdsDeleteDeviceNotificationRequest.of(
                targetAmsNetId,
                targetAmsPort,
                sourceAmsNetId,
                sourceAmsPort,
                Invoke.NONE,
                notificationHandle
            );
            CompletableFuture<PlcProprietaryResponse<AdsDeleteDeviceNotificationResponse>> deleteDeviceFuture = new CompletableFuture<>();
            channel.writeAndFlush(new PlcRequestContainer<>(new PlcProprietaryRequest<>(adsDeleteDeviceNotificationRequest), deleteDeviceFuture));

            PlcProprietaryResponse<AdsDeleteDeviceNotificationResponse> deleteDeviceResponse;
            try {
                deleteDeviceResponse = deleteDeviceFuture.get(DEL_DEVICE_TIMEOUT, TimeUnit.MILLISECONDS);
            } catch (InterruptedException e) {
                LOGGER.warn("Interrupted!", e);
                Thread.currentThread().interrupt();
                throw new PlcRuntimeException(e);
            } catch (ExecutionException | TimeoutException e) {
                throw new PlcRuntimeException(e);
            }
            AdsDeleteDeviceNotificationResponse response = deleteDeviceResponse.getResponse();
            if (response.getResult().toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
                throw new PlcRuntimeException("Non error code received " + response.getResult());
            }

            getChannel().pipeline().get(Plc4x2AdsProtocol.class).removeConsumer(handlePair.getLeft());
        }
    }
}
