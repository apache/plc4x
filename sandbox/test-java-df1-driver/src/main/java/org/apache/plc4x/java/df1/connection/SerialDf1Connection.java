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
package org.apache.plc4x.java.df1.connection;

import io.netty.channel.*;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.serial.connection.connection.SerialChannelFactory;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.df1.protocol.Df1Protocol;
import org.apache.plc4x.java.df1.protocol.Plc4XDf1Protocol;
import org.apache.plc4x.java.df1.util.Df1FieldHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;

public class SerialDf1Connection extends BaseDf1Connection implements PlcReader, PlcWriter {

    private static final Logger logger = LoggerFactory.getLogger(SerialDf1Connection.class);

    private short localAddr;
    private short remoteAddr;

    public SerialDf1Connection(String comPortName, String params) {
        this(new SerialChannelFactory(comPortName), params);
    }

    public SerialDf1Connection(ChannelFactory channelFactory, String params) {
        super(channelFactory, false);
        this.localAddr = (short) 0x00;
        this.remoteAddr = (short) 0x09;

        // Override some of the settings, if they are asked for.
        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        case "local-addr":
                            this.localAddr = Short.parseShort(paramValue);
                            break;
                        case "remote-addr":
                            this.remoteAddr = Short.parseShort(paramValue);
                            break;
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                }
            }
        }
    }

    @Override
    public boolean canRead() {
        return true;
    }

    @Override
    public boolean canWrite() {
        return true;
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
                pipeline.addLast(new Df1Protocol(localAddr, remoteAddr));
                pipeline.addLast(new Plc4XDf1Protocol());
            }
        };
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new Df1FieldHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return new DefaultPlcWriteRequest.Builder(this, new Df1FieldHandler());
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<InternalPlcReadResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcReadRequest, InternalPlcReadResponse> container =
            new PlcRequestContainer<>((InternalPlcReadRequest) readRequest, future);
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });
        return future
            .thenApply(PlcReadResponse.class::cast);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<InternalPlcWriteResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcWriteRequest, InternalPlcWriteResponse> container =
            new PlcRequestContainer<>((InternalPlcWriteRequest) writeRequest, future);
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });
        return future
            .thenApply(PlcWriteResponse.class::cast);
    }

}
