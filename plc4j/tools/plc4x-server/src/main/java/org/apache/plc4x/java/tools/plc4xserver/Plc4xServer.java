/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.tools.plc4xserver;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xConstants;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xMessage;
import org.apache.plc4x.java.spi.connection.GeneratedProtocolMessageCodec;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.tools.plc4xserver.protocol.Plc4xServerAdapter;

import java.util.function.ToIntFunction;

public class Plc4xServer {

    private EventLoopGroup loopGroup;
    private EventLoopGroup workerGroup;

    public Plc4xServer() {
    }

    public void start() throws PlcException {
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
                        pipeline.addLast(new GeneratedProtocolMessageCodec<>(
                            Plc4xMessage.class, Plc4xMessage::staticParse, ByteOrder.BIG_ENDIAN, null,
                            new ByteLengthEstimator(), null));
                        pipeline.addLast(new Plc4xServerAdapter());
                    }
                }).option(ChannelOption.SO_BACKLOG, 128)
                .childOption(ChannelOption.SO_KEEPALIVE, true);

            bootstrap.bind(Plc4xConstants.PLC4XTCPDEFAULTPORT).sync();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcException(e);
        }
    }

    public void stop() {
        if(workerGroup == null) {
            return;
        }

        workerGroup.shutdownGracefully();
        loopGroup.shutdownGracefully();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 3) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 1);
            }
            return -1;
        }
    }

    public static void main(String[] args) throws Exception {
        new Plc4xServer().start();
    }

}
