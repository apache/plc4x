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
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.model.AdsAddress;
import org.apache.plc4x.java.ads.model.AdsSubscriptionHandle;
import org.apache.plc4x.java.ads.model.SymbolicAdsAddress;
import org.apache.plc4x.java.ads.protocol.Ads2PayloadProtocol;
import org.apache.plc4x.java.ads.protocol.Payload2TcpProtocol;
import org.apache.plc4x.java.ads.protocol.Plc4x2AdsProtocol;
import org.apache.plc4x.java.ads.protocol.util.LittleEndianDecoder;
import org.apache.plc4x.java.api.connection.PlcSubscriber;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.SubscriptionEventItem;
import org.apache.plc4x.java.api.messages.items.SubscriptionRequestItem;
import org.apache.plc4x.java.api.messages.items.SubscriptionResponseItem;
import org.apache.plc4x.java.api.messages.items.UnsubscriptionRequestItem;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.api.types.ResponseCode;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class AdsTcpPlcConnection extends AdsAbstractPlcConnection implements PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsTcpPlcConnection.class);

    private static final int TCP_PORT = 48898;

    private static final long ADD_DEVICE_TIMEOUT = CONF.getLong("plc4x.adsconnection.add.device,timeout", 3000);
    private static final long DEL_DEVICE_TIMEOUT = CONF.getLong("plc4x.adsconnection.del.device,timeout", 3000);

    private static AtomicInteger localPorts = new AtomicInteger(30000);

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
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        // TODO: Make this multi-value
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        if(subscriptionRequest.getNumberOfItems() == 1) {
            SubscriptionRequestItem<?> subscriptionRequestItem = subscriptionRequest.getRequestItem().orElse(null);

            Objects.requireNonNull(subscriptionRequestItem);
            Objects.requireNonNull(subscriptionRequestItem.getConsumer());
            Objects.requireNonNull(subscriptionRequestItem.getAddress());
            Objects.requireNonNull(subscriptionRequestItem.getDatatype());

            Address address = subscriptionRequestItem.getAddress();
            Class<?> datatype = subscriptionRequestItem.getDatatype();

            IndexGroup indexGroup;
            IndexOffset indexOffset;
            // If this is a symbolic address, it has to be resolved first.
            // TODO: This is blocking, should be changed to be async.
            if (address instanceof SymbolicAdsAddress) {
                mapAddress((SymbolicAdsAddress) address);
                AdsAddress adsAddress = addressMapping.get(address);
                if (adsAddress == null) {
                    throw new PlcRuntimeException("Unresolvable address" + address);
                }
                indexGroup = IndexGroup.of(adsAddress.getIndexGroup());
                indexOffset = IndexOffset.of(adsAddress.getIndexOffset());
            }
            // If it's no symbolic address, we can continue immediately
            // without having to do any resolving.
            else if (address instanceof AdsAddress) {
                AdsAddress adsAddress = (AdsAddress) address;
                indexGroup = IndexGroup.of(adsAddress.getIndexGroup());
                indexOffset = IndexOffset.of(adsAddress.getIndexOffset());
            } else {
                throw new IllegalArgumentException("Unsupported address type " + address.getClass());
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
                LittleEndianDecoder.getLengthFor(datatype, 1),
                TransmissionMode.DefinedValues.ADSTRANS_SERVERCYCLE,
                MaxDelay.of(0),
                CycleTime.of(4000000)
            );

            // Send the request to the plc and wait for a response
            // TODO: This is blocking, should be changed to be async.
            CompletableFuture<PlcProprietaryResponse<AdsAddDeviceNotificationResponse>> addDeviceFuture = new CompletableFuture<>();
            channel.writeAndFlush(new PlcRequestContainer<>(new PlcProprietaryRequest<>(adsAddDeviceNotificationRequest), addDeviceFuture));
            PlcProprietaryResponse<AdsAddDeviceNotificationResponse> addDeviceResponse = getFromFuture(addDeviceFuture, ADD_DEVICE_TIMEOUT);
            AdsAddDeviceNotificationResponse response = addDeviceResponse.getResponse();

            // Abort if we got anything but a successful response.
            if (response.getResult().toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
                throw new PlcRuntimeException("Error code received " + response.getResult());
            }
            AdsSubscriptionHandle adsSubscriptionHandle = new AdsSubscriptionHandle(response.getNotificationHandle());
            future.complete(new PlcSubscriptionResponse(subscriptionRequest, Collections.singletonList(
                new SubscriptionResponseItem<>(subscriptionRequestItem, adsSubscriptionHandle, ResponseCode.OK))));

            Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer =
                adsDeviceNotificationRequest -> adsDeviceNotificationRequest.getAdsStampHeaders().forEach(adsStampHeader -> {
                    Calendar timeStamp = Calendar.getInstance();
                    timeStamp.setTime(adsStampHeader.getTimeStamp().getAsDate());

                    adsStampHeader.getAdsNotificationSamples()
                        .forEach(adsNotificationSample -> {
                            Data data = adsNotificationSample.getData();
                            try {
                                @SuppressWarnings("unchecked")
                                List<?> decodeData = LittleEndianDecoder.decodeData(datatype, data.getBytes());
                                SubscriptionEventItem subscriptionEventItem =
                                    new SubscriptionEventItem(subscriptionRequestItem, timeStamp, decodeData);
                                subscriptionRequestItem.getConsumer().accept(subscriptionEventItem);
                            } catch (PlcProtocolException | RuntimeException e) {
                                LOGGER.error("Can't decode {}", data, e);
                            }
                        });
                });
            // TODO: What's this for? Is this still needed if we use the consumers in the subscriptions?
            getChannel().pipeline().get(Plc4x2AdsProtocol.class).addConsumer(adsDeviceNotificationRequestConsumer);
        }
        return future;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        for (UnsubscriptionRequestItem unsubscriptionRequestItem : unsubscriptionRequest.getRequestItems()) {
            Objects.requireNonNull(unsubscriptionRequestItem);
            if(unsubscriptionRequestItem.getSubscriptionHandle() instanceof AdsSubscriptionHandle) {
                AdsSubscriptionHandle adsSubscriptionHandle =
                    (AdsSubscriptionHandle) unsubscriptionRequestItem.getSubscriptionHandle();
                AdsDeleteDeviceNotificationRequest adsDeleteDeviceNotificationRequest =
                    AdsDeleteDeviceNotificationRequest.of(
                        targetAmsNetId,
                        targetAmsPort,
                        sourceAmsNetId,
                        sourceAmsPort,
                        Invoke.NONE,
                        adsSubscriptionHandle.getNotificationHandle()
                    );
                CompletableFuture<PlcProprietaryResponse<AdsDeleteDeviceNotificationResponse>> deleteDeviceFuture =
                    new CompletableFuture<>();
                channel.writeAndFlush(new PlcRequestContainer<>(new PlcProprietaryRequest<>(
                    adsDeleteDeviceNotificationRequest), deleteDeviceFuture));

                PlcProprietaryResponse<AdsDeleteDeviceNotificationResponse> deleteDeviceResponse =
                    getFromFuture(deleteDeviceFuture, DEL_DEVICE_TIMEOUT);
                AdsDeleteDeviceNotificationResponse response = deleteDeviceResponse.getResponse();
                if (response.getResult().toAdsReturnCode() != AdsReturnCode.ADS_CODE_0) {
                    throw new PlcRuntimeException("Non error code received " + response.getResult());
                }
            }
        }
        CompletableFuture<PlcUnsubscriptionResponse> future = new CompletableFuture<>();
        future.complete(new PlcUnsubscriptionResponse());
        return future;
    }
}
