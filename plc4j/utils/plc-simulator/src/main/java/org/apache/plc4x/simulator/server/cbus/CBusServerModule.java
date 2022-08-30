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
package org.apache.plc4x.simulator.server.cbus;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import org.apache.plc4x.java.cbus.CBusDriver;
import org.apache.plc4x.java.cbus.readwrite.CBusConstants;
import org.apache.plc4x.java.cbus.readwrite.CBusMessage;
import org.apache.plc4x.java.cbus.readwrite.CBusOptions;
import org.apache.plc4x.java.cbus.readwrite.RequestContext;
import org.apache.plc4x.java.spi.connection.GeneratedProtocolMessageCodec;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.simulator.PlcSimulatorConfig;
import org.apache.plc4x.simulator.exceptions.SimulatorException;
import org.apache.plc4x.simulator.model.Context;
import org.apache.plc4x.simulator.server.ServerModule;
import org.apache.plc4x.simulator.server.cbus.protocol.CBusServerAdapter;

public class CBusServerModule implements ServerModule {

    private EventLoopGroup loopGroup;
    private EventLoopGroup workerGroup;
    private Context context;
    private PlcSimulatorConfig config;

    @Override
    public String getName() {
        return "C-BUS";
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
            workerGroup = new NioEventLoopGroup();

            ServerBootstrap bootstrap = new ServerBootstrap();
            bootstrap.group(loopGroup, workerGroup)
                .channel(NioServerSocketChannel.class)
                .childHandler(new ChannelInitializer<SocketChannel>() {
                    @Override
                    public void initChannel(SocketChannel channel) {
                        ChannelPipeline pipeline = channel.pipeline();
                        pipeline.addLast(new GeneratedProtocolMessageCodec<>(CBusMessage.class,
                            CBusMessage::staticParse, ByteOrder.BIG_ENDIAN,
                            new Object[]{false, new RequestContext(false), new CBusOptions(false, false, false, false, false, false, false, false, false)},
                            new CBusDriver.ByteLengthEstimator(),
                            new CBusDriver.CorruptPackageCleaner()));
                        pipeline.addLast(new CBusServerAdapter(context));
                    }
                }).option(ChannelOption.SO_BACKLOG, 128)
                .childOption(ChannelOption.SO_KEEPALIVE, true);

            bootstrap.bind(config.getHost(), CBusConstants.CBUSTCPDEFAULTPORT).sync();
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
