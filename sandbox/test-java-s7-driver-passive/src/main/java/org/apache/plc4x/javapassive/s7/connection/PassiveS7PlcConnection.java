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
package org.apache.plc4x.javapassive.s7.connection;

import io.netty.channel.*;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.DefaultPlcFieldHandler;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.base.connection.RawSocketChannelFactory;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketIpAddress;
import org.apache.plc4x.java.utils.rawsockets.netty.TcpIpPacketHandler;
import org.apache.plc4x.javapassive.s7.protocol.HelloWorldProtocol;
import org.apache.plc4x.javapassive.s7.protocol.PassiveS7Protocol;

import java.util.concurrent.CompletableFuture;

public class PassiveS7PlcConnection extends NettyPlcConnection implements PlcReader {

    public PassiveS7PlcConnection(RawSocketIpAddress address, String params) {
        this(new RawSocketChannelFactory(address.getDeviceName(), null,
            address.getPort(), -1, new TcpIpPacketHandler()), params);
    }

    public PassiveS7PlcConnection(ChannelFactory channelFactory, String param) {
        super(channelFactory, false);
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
                pipeline.addLast(new PassiveS7Protocol());
                pipeline.addLast(new HelloWorldProtocol());
            }
        };
    }

    @Override
    public boolean canRead() {
        return true;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new DefaultPlcFieldHandler() {
            @Override
            public PlcField createField(String fieldQuery) throws PlcInvalidFieldException {
                return new PlcField() {
                    @Override
                    public Class<?> getDefaultJavaType() {
                        return String.class;
                    }
                };
            }
        });
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        return new CompletableFuture<>();
    }

}
