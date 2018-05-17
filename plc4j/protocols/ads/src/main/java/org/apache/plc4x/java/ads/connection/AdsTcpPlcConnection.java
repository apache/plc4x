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
import org.apache.plc4x.java.ads.api.commands.AdsDeviceNotificationRequest;
import org.apache.plc4x.java.ads.api.commands.types.AdsNotificationSample;
import org.apache.plc4x.java.ads.api.commands.types.AdsStampHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.util.ByteValue;
import org.apache.plc4x.java.ads.protocol.Ads2PayloadProtocol;
import org.apache.plc4x.java.ads.protocol.Payload2TcpProtocol;
import org.apache.plc4x.java.ads.protocol.Plc4x2AdsProtocol;
import org.apache.plc4x.java.api.connection.PlcSubscriber;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcNotification;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.stream.Collectors;

public class AdsTcpPlcConnection extends AdsAbstractPlcConnection implements PlcSubscriber {

    private static final int TCP_PORT = 48898;

    private final Map<Consumer<PlcNotification>, Consumer<AdsDeviceNotificationRequest>> subscriberMap = new HashMap<>();

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
        return AmsPort.of(TCP_PORT);
    }

    @Override
    public void subscribe(Consumer<PlcNotification> consumer, Address address) {
        Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer = adsDeviceNotificationRequest -> {
            for (AdsStampHeader adsStampHeader : adsDeviceNotificationRequest.getAdsStampHeaders()) {
                Date timeStamp = adsStampHeader.getTimeStamp().getAsDate();
                // TODO: where do we implement the mapping. Better move it into the ...
                List<Object> values = adsStampHeader.getAdsNotificationSamples()
                    .stream()
                    .map(AdsNotificationSample::getData)
                    .map(ByteValue::getBytes)
                    .map(data -> (Object) data)
                    .collect(Collectors.toList());
                consumer.accept(new PlcNotification(timeStamp, values));
            }
        };
        subscriberMap.put(consumer, adsDeviceNotificationRequestConsumer);
        getChannel().pipeline().get(Plc4x2AdsProtocol.class).addConsumer(adsDeviceNotificationRequestConsumer);
    }

    @Override
    public void unsubscribe(Consumer<PlcNotification> consumer) {
        Consumer<AdsDeviceNotificationRequest> adsDeviceNotificationRequestConsumer = subscriberMap.remove(consumer);
        if (adsDeviceNotificationRequestConsumer != null) {
            getChannel().pipeline().get(Plc4x2AdsProtocol.class).removeConsumer(adsDeviceNotificationRequestConsumer);
        }
    }
}
