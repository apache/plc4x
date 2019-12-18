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
package org.apache.plc4x.java.passive.knxnetip.connection;

import io.netty.channel.*;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.RawSocketChannelFactory;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.NettyPlcConnection;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.passive.knxnetip.model.KnxNetIpField;
import org.apache.plc4x.java.passive.knxnetip.protocol.KnxNetIpProtocol;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketAddress;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketIpAddress;
import org.apache.plc4x.java.utils.rawsockets.netty.UdpIpPacketHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;

public class PassiveKnxNetIpPlcConnection extends NettyPlcConnection implements PlcReader {

    private static final Logger logger = LoggerFactory.getLogger(PassiveKnxNetIpPlcConnection.class);

    private final ChannelHandler handler;

    public PassiveKnxNetIpPlcConnection(RawSocketIpAddress address, String params, ChannelHandler handler) {
        this(new RawSocketChannelFactory(address.getDeviceName(), null,
            address.getPort(), RawSocketAddress.ALL_PROTOCOLS, new UdpIpPacketHandler()), params, handler);
    }

    public PassiveKnxNetIpPlcConnection(ChannelFactory channelFactory, String params, ChannelHandler handler) {
        super(channelFactory, true);
        this.handler = handler;
    }

    @Override
    protected void sendChannelCreatedEvent() {
        // As this type of protocol doesn't require any form of connection, we just send the connected event.
        channel.pipeline().fireUserEventTriggered(new ConnectedEvent());
    }

    @Override
    public PlcField prepareField(String fieldQuery) throws PlcInvalidFieldException {
        return KnxNetIpField.of(fieldQuery);
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the s7 protocol.
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
                pipeline.addLast(new KnxNetIpProtocol());
                pipeline.addLast(handler);
            }
        };
    }

    @Override
    public boolean canRead() {
        return true;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new KnxNetIpFieldHandler());
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        InternalPlcReadRequest internalReadRequest = checkInternal(readRequest, InternalPlcReadRequest.class);
        CompletableFuture<InternalPlcReadResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcReadRequest, InternalPlcReadResponse> container =
            new PlcRequestContainer<>(internalReadRequest, future);
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });
        return future
            .thenApply(PlcReadResponse.class::cast);
    }

}
