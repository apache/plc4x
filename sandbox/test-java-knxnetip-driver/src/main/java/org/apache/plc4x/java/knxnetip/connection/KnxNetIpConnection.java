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
package org.apache.plc4x.java.knxnetip.connection;

import io.netty.channel.*;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.UdpSocketChannelFactory;
import org.apache.plc4x.java.base.connection.protocol.DatagramUnpackingHandler;
import org.apache.plc4x.java.knxnetip.readwrite.KNXNetIPMessage;
import org.apache.plc4x.java.spi.Plc4xNettyWrapper;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.NettyPlcConnection;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.spi.events.DisconnectEvent;
import org.apache.plc4x.java.knxnetip.model.KnxNetIpField;
import org.apache.plc4x.java.knxnetip.protocol.KnxNetIpProtocolLogic;
import org.apache.plc4x.java.knxnetip.protocol.KnxNetIpProtocolMessage;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.parser.ConnectionParser;

import java.net.InetAddress;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

public class KnxNetIpConnection extends NettyPlcConnection implements PlcReader {

    public static final int KNXNET_IP_PORT = 3671;

    private final KnxNetIpConfiguration configuration;

    public KnxNetIpConnection(InetAddress address, String params) {
        this(new UdpSocketChannelFactory(address, KNXNET_IP_PORT), params);
    }

    public KnxNetIpConnection(ChannelFactory channelFactory, String params) {
        super(channelFactory, true);
        configuration = ConnectionParser.parse("a://1.1.1.1?" + params, KnxNetIpConfiguration.class);
    }

    @Override
    public PlcField prepareField(String fieldQuery) {
        return KnxNetIpField.of(fieldQuery);
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer<Channel>() {
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
                // Unpack the ByteBuf included in the DatagramPackage.
                pipeline.addLast(new DatagramUnpackingHandler());
                pipeline.addLast(new KnxNetIpProtocolMessage());

                Plc4xProtocolBase<KNXNetIPMessage> knxNetIpProtocolLogic = new KnxNetIpProtocolLogic(configuration);
                setProtocol(knxNetIpProtocolLogic);
                Plc4xNettyWrapper<KNXNetIPMessage> context =
                    new Plc4xNettyWrapper<>(pipeline, knxNetIpProtocolLogic, KNXNetIPMessage.class);
                pipeline.addLast(context);
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
    public void close() throws PlcConnectionException {
        if(channel == null) {
            super.close();
            return;
        }

        CompletableFuture<Void> disconnectFuture = new CompletableFuture<>();
        channel.pipeline().fireUserEventTriggered(new DisconnectEvent(disconnectFuture));
        try {
            // Wait for the connection to be disconnected.
            disconnectFuture.get(500, TimeUnit.MILLISECONDS);
            super.close();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException("Error closing connection");
        } catch (ExecutionException | TimeoutException e) {
            throw new PlcConnectionException("Error closing connection");
        }
    }

    @Override
    protected void sendChannelCreatedEvent() {
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

}
