/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.simulator.server.bacnet;

import io.netty.bootstrap.Bootstrap;
import io.netty.bootstrap.ServerBootstrap;
import io.netty.buffer.ByteBuf;
import io.netty.channel.*;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioDatagramChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.DatagramPacketDecoder;
import io.netty.handler.codec.DatagramPacketEncoder;
import io.netty.handler.codec.MessageToMessageDecoder;
import io.netty.handler.codec.MessageToMessageEncoder;
import org.apache.plc4x.java.bacnetip.BacNetIpDriver;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.bacnetip.readwrite.BacnetConstants;
import org.apache.plc4x.java.cbus.CBusDriver;
import org.apache.plc4x.java.cbus.readwrite.CBusConstants;
import org.apache.plc4x.java.cbus.readwrite.CBusMessage;
import org.apache.plc4x.java.cbus.readwrite.CBusOptions;
import org.apache.plc4x.java.cbus.readwrite.RequestContext;
import org.apache.plc4x.java.spi.connection.GeneratedProtocolMessageCodec;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.apache.plc4x.simulator.PlcSimulatorConfig;
import org.apache.plc4x.simulator.exceptions.SimulatorException;
import org.apache.plc4x.simulator.model.Context;
import org.apache.plc4x.simulator.server.ServerModule;
import org.apache.plc4x.simulator.server.bacnet.protocol.BacnetServerAdapter;

import java.util.List;

public class BacnetServerModule implements ServerModule {

    private EventLoopGroup loopGroup;
    private EventLoopGroup workerGroup;
    private Context context;
    private PlcSimulatorConfig config;

    @Override
    public String getName() {
        return "Bacnet";
    }

    @Override
    public void setConfig(PlcSimulatorConfig config) {
        this.config = config;
    }

    @Override
    public void setContext(Context context) {
        this.context = context;
    }

    @Override
    public void start() throws SimulatorException {
        if (loopGroup != null) {
            return;
        }

        try {
            loopGroup = new NioEventLoopGroup();

            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(loopGroup)
                .channel(NioDatagramChannel.class)
                .option(ChannelOption.SO_BROADCAST, true)
                .handler(new ChannelInitializer<NioDatagramChannel>() {
                    @Override
                    public void initChannel(NioDatagramChannel channel) {
                        ChannelPipeline pipeline = channel.pipeline();
                        pipeline.addLast(new DatagramPacketDecoder(new MessageToMessageDecoder<ByteBuf>() {
                            @Override
                            protected void decode(ChannelHandlerContext ctx, ByteBuf msg, List<Object> out) throws Exception {
                                byte[] bytes = new byte[msg.readableBytes()];
                                msg.readBytes(bytes);
                                out.add(BVLC.staticParse(new ReadBufferByteBased(bytes)));
                            }
                        }));
                        pipeline.addLast(new DatagramPacketEncoder<>(new MessageToMessageEncoder<BVLC>() {
                            @Override
                            protected void encode(ChannelHandlerContext ctx, BVLC msg, List<Object> out) throws Exception {
                                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(msg.getLengthInBytes());
                                msg.serialize(writeBuffer);
                                out.add(writeBuffer.getBytes());
                            }
                        }));
                        pipeline.addLast(new BacnetServerAdapter(context));
                    }
                });

            int port = BacnetConstants.BACNETUDPDEFAULTPORT;
            if (config.getCBusPort() != null) {
                port = Integer.parseInt(config.getCBusPort());
            }
            String host = config.getHost();
            if (host != null) {
                bootstrap.bind(host, port).sync();
            } else {
                bootstrap.bind(port).sync();
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new SimulatorException(e);
        }
    }

    @Override
    public void stop() {
        if (workerGroup == null) {
            return;
        }

        workerGroup.shutdownGracefully();
        loopGroup.shutdownGracefully();
    }

}
