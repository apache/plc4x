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
package org.apache.plc4x.java.ethernetip.connection;

import io.netty.channel.*;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.ethernetip.netty.EnipCodec;
import org.apache.plc4x.java.ethernetip.netty.Plc4XEtherNetIpProtocol;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class EtherNetIpTcpPlcConnection extends BaseEtherNetIpPlcConnection {

    // Port 44818
    private static final int ETHERNET_IP_TCP_PORT = 0xAF12;

    private static final Logger logger = LoggerFactory.getLogger(EtherNetIpTcpPlcConnection.class);

    public EtherNetIpTcpPlcConnection(InetAddress address, String params) {
        this(new TcpSocketChannelFactory(address, ETHERNET_IP_TCP_PORT), params);
        logger.info("Configured EtherNetIpTcpPlcConnection with: host-name {}", address.getHostAddress());
    }

    public EtherNetIpTcpPlcConnection(InetAddress address, int port, String params) {
        this(new TcpSocketChannelFactory(address, port), params);
        logger.info("Configured EtherNetIpTcpPlcConnection with: host-name {}", address.getHostAddress());
    }

    public EtherNetIpTcpPlcConnection(ChannelFactory channelFactory, String params) {
        super(channelFactory, params);
    }

    @Override
    public Optional<PlcSubscriptionRequest.Builder> subscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public Optional<PlcUnsubscriptionRequest.Builder> unsubscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new ChannelInboundHandlerAdapter() {
                    @Override
                    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
                        if (evt instanceof ConnectedEvent) {
                            sessionSetupCompleteFuture.complete(null);
                        } else {
                            super.userEventTriggered(ctx, evt);
                        }
                    }
                });
                pipeline.addLast(new EnipCodec());
                pipeline.addLast(new Plc4XEtherNetIpProtocol());
            }
        };
    }

    @Override
    protected void sendChannelCreatedEvent() {
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

}
