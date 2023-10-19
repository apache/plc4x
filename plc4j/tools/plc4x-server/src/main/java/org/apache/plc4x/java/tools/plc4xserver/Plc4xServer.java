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

package org.apache.plc4x.java.tools.plc4xserver;

import static java.lang.Runtime.getRuntime;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import java.io.IOException;
import java.net.ServerSocket;
import java.util.Arrays;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Future;
import java.util.function.ToIntFunction;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xConstants;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xMessage;
import org.apache.plc4x.java.spi.connection.GeneratedProtocolMessageCodec;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.tools.plc4xserver.protocol.Plc4xServerAdapter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Plc4xServer {

    public static final String SERVER_PORT_PROPERTY = "plc4x.server.port";
    public static final String SERVER_PORT_ENVIRONMENT_VARIABLE = "PLC4X_SERVER_PORT";
    public static int DEFAULT_PORT = Plc4xConstants.PLC4XTCPDEFAULTPORT;

    private static final Logger LOG = LoggerFactory.getLogger(Plc4xServerAdapter.class);

    private EventLoopGroup loopGroup;
    private EventLoopGroup workerGroup;
    private ChannelFuture channelFuture;
    private Integer port;

    public static void main(String[] args) throws Exception {
        final Plc4xServer server = new Plc4xServer();

        Future<Void> serverFuture = server.start(
                Arrays.stream(args).findFirst() // port number given as first command line argument
                        .or(() -> Optional.ofNullable(System.getProperty(SERVER_PORT_PROPERTY)))
                        .or(() -> Optional.ofNullable(System.getenv(SERVER_PORT_ENVIRONMENT_VARIABLE)))
                        .map(Integer::parseInt)
                        .orElse(DEFAULT_PORT)
        );
        CompletableFuture<Void> serverRunning = new CompletableFuture<>();
        getRuntime().addShutdownHook(new Thread(() -> serverRunning.complete(null)));

        try {
            LOG.info("Server is configured to listen on TCP port {}", server.getPort());
            serverFuture.get();
            LOG.info("Server is ready.");
            serverRunning.get();
        } catch (InterruptedException e) {
            throw new PlcRuntimeException(e);
        } finally {
            LOG.info("Server is shutting down...");
            server.stop();
        }
    }

    public Integer getPort() {
        return port;
    }

    public Future<Void> start() {
        return start(0);
    }

    public Future<Void> start(int port) {
        if (port == 0) {
            this.port = findRandomFreePort();
        } else {
            this.port = port;
        }

        if (loopGroup != null) {
            return CompletableFuture.completedFuture(null);
        }

        loopGroup = new NioEventLoopGroup();
        workerGroup = new NioEventLoopGroup();

        channelFuture = new ServerBootstrap()
                .group(loopGroup, workerGroup)
                .channel(NioServerSocketChannel.class)
                .childHandler(new SocketChannelChannelInitializer())
                .option(ChannelOption.SO_BACKLOG, 128)
                .childOption(ChannelOption.SO_KEEPALIVE, true)
                .bind(this.port);

        return channelFuture;
    }

    public void stop() {
        if (workerGroup == null) {
            return;
        }

        channelFuture.cancel(true);

        workerGroup.shutdownGracefully();
        loopGroup.shutdownGracefully();

        workerGroup = null;
        loopGroup = null;
    }

    private static class SocketChannelChannelInitializer extends ChannelInitializer<SocketChannel> {

        @Override
        public void initChannel(SocketChannel channel) {
            ChannelPipeline pipeline = channel.pipeline();
            pipeline.addLast(
                    new GeneratedProtocolMessageCodec<>(
                            Plc4xMessage.class,
                            Plc4xMessage::staticParse,
                            ByteOrder.BIG_ENDIAN,
                            null,
                            new ByteLengthEstimator(),
                            null
                    )
            );
            pipeline.addLast(new Plc4xServerAdapter());
        }
    }

    private static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {

        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 3) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 1);
            }
            return -1;
        }
    }

    private static int findRandomFreePort() {
        final int port;
        try (ServerSocket socket = new ServerSocket(0)) {
            port = socket.getLocalPort();
        } catch (IOException e) {
            throw new RuntimeException("Couldn't find any free port.", e);
        }
        return port;
    }
}
