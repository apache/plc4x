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
package org.apache.plc4x.simulator.server.s7;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import org.apache.plc4x.java.s7.readwrite.S7Driver;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.s7.readwrite.io.TPKTPacketIO;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedProtocolMessageCodec;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.simulator.exceptions.SimulatorExcepiton;
import org.apache.plc4x.simulator.model.Context;
import org.apache.plc4x.simulator.server.ServerModule;
import org.apache.plc4x.simulator.server.s7.protocol.S7Step7ServerAdapter;

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.configure;

public class S7ServerModule implements ServerModule {

    private static final int ISO_ON_TCP_PORT = 102;

    private EventLoopGroup loopGroup;
    private EventLoopGroup workerGroup;
    // private Context context;
    private S7PlcHandler handler;

    @Override
    public String getName() {
        return "S7-STEP7";
    }

    @Override
    public void setContext(Context context) {
        // Do Nothing...
    }

    public void setHandler(S7PlcHandler handler) {
        this.handler = handler;
    }

    @Override
    public void start() throws SimulatorExcepiton {
        if (this.handler == null) {
            throw new IllegalStateException("A handler has to be set to start the PLC Server");
        }
        if(loopGroup != null) {
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
                        pipeline.addLast(new GeneratedProtocolMessageCodec<>(TPKTPacket.class, new TPKTPacketIO(), true, null,
                            new S7Driver.ByteLengthEstimator(),
                            new S7Driver.CorruptPackageCleaner()));
                        pipeline.addLast(new S7Step7ServerAdapter(handler));
                    }
                }).option(ChannelOption.SO_BACKLOG, 128)
                .childOption(ChannelOption.SO_KEEPALIVE, true);

            bootstrap.bind(ISO_ON_TCP_PORT).sync();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new SimulatorExcepiton(e);
        }
    }

    @Override
    public void stop() {
        if(workerGroup == null) {
            return;
        }

        workerGroup.shutdownGracefully();
        loopGroup.shutdownGracefully();
    }

}
