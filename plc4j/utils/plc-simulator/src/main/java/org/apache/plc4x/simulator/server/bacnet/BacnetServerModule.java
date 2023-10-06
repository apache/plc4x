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
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.DatagramPacket;
import io.netty.channel.socket.nio.NioDatagramChannel;
import io.netty.handler.codec.MessageToMessageDecoder;
import io.netty.handler.codec.MessageToMessageEncoder;
import io.netty.util.concurrent.FastThreadLocal;
import org.apache.plc4x.java.bacnetip.BacNetIpDriver;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.bacnetip.readwrite.BacnetConstants;
import org.apache.plc4x.java.spi.connection.GeneratedProtocolMessageCodec;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.simulator.PlcSimulatorConfig;
import org.apache.plc4x.simulator.exceptions.SimulatorException;
import org.apache.plc4x.simulator.model.Context;
import org.apache.plc4x.simulator.server.ServerModule;
import org.apache.plc4x.simulator.server.bacnet.protocol.BacnetServerAdapter;

import java.net.InetSocketAddress;
import java.util.List;

public class BacnetServerModule implements ServerModule {

    private EventLoopGroup loopGroup;
    private EventLoopGroup workerGroup;
    private Context context;
    private PlcSimulatorConfig config;

    private final FastThreadLocal<String> returnToSender = new FastThreadLocal<>();

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
                        channel.pipeline()
                            .addLast(new MessageToMessageDecoder<DatagramPacket>() {
                                @Override
                                protected void decode(ChannelHandlerContext ctx, DatagramPacket msg, List<Object> out) throws Exception {
                                    final ByteBuf content = msg.content();
                                    out.add(content.retain());
                                    String value = msg.sender().getHostString();
                                    returnToSender.set(value);
                                }
                            })
                            .addLast(new MessageToMessageEncoder<ByteBuf>() {
                                @Override
                                protected void encode(ChannelHandlerContext ctx, ByteBuf msg, List<Object> out) throws Exception {
                                    msg.retain();
                                    out.add(new DatagramPacket(msg, new InetSocketAddress(returnToSender.get(), 47808)));
                                }
                            })
                            .addLast(new GeneratedProtocolMessageCodec<>(BVLC.class,
                                BVLC::staticParse, ByteOrder.BIG_ENDIAN,
                                null,
                                new BacNetIpDriver.ByteLengthEstimator(),
                                new BacNetIpDriver.CorruptPackageCleaner()))
                            .addLast(new BacnetServerAdapter(context));
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
