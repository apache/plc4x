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
package org.apache.plc4x.java.abeth.connection;

import io.netty.channel.*;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.abeth.model.AbEthField;
import org.apache.plc4x.java.abeth.protocol.AbEthProtocol;
import org.apache.plc4x.java.abeth.protocol.Plc4xAbEthProtocol;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.concurrent.CompletableFuture;

public class AbEthPlcConnection extends NettyPlcConnection implements PlcReader {

    private static final int AB_ETH_PORT = 2222;
    private static final Logger logger = LoggerFactory.getLogger(AbEthPlcConnection.class);

    private final int station;


    public AbEthPlcConnection(InetAddress address, int station, String params) {
        this(new TcpSocketChannelFactory(address, AB_ETH_PORT), station, params);
        logger.info("Setting up AB-ETH Connection with: host-name {}", address.getHostAddress());
    }

    public AbEthPlcConnection(ChannelFactory channelFactory, int station, String params) {
        super(channelFactory, true);
        this.station = station;

        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                } else {
                    logger.debug("Unknown no-value parameter {}", paramName);
                }
            }
        }
    }

    @Override
    protected void sendChannelCreatedEvent() {
        logger.trace("Channel was created, firing ChannelCreated Event");
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

    @Override
    public PlcField prepareField(String fieldQuery) throws PlcInvalidFieldException {
        return AbEthField.of(fieldQuery);
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
                pipeline.addLast(new AbEthProtocol());
                pipeline.addLast(new Plc4xAbEthProtocol(station));
            }
        };
    }

    @Override
    public boolean canRead() {
        return true;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new AbEthFieldHandler());
    }

    @Override
    public void close() {
        logger.debug("Closing PlcConnection...");
        // Close the channel gracefully
        if ((channel != null) && channel.isOpen()) {
            logger.debug("Channel is still connected, Closing channel...");
            // Close the channel
            channel.close();

            // Do some additional cleanup operations ...
            // In normal operation, the channels event loop has a parent, however when running with
            // the embedded channel for unit tests, parent is null.
            if (channel.eventLoop().parent() != null) {
                logger.trace("Shutting down EventLoop gracefully...");
                channel.eventLoop().parent().shutdownGracefully();
                logger.trace("Eventloop is shutted down");
            }
        }
        // Dereference
        channel = null;
        connected = false;
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