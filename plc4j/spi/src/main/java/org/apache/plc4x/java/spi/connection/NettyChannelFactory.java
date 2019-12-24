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

package org.apache.plc4x.java.spi.connection;

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoop;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.GenericFutureListener;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.SocketAddress;

/**
 * Adapter with sensible defaults for a Netty Based Channel Factory.
 */
public abstract class NettyChannelFactory implements ChannelFactory {

    private static final Logger logger = LoggerFactory.getLogger(NettyChannelFactory.class);

    /**
     * Channel to Use, e.g. NiO, EiO
     */
    public abstract Class<? extends Channel> getChannel();

    /**
     * Event Loop Group to use.
     * Has to be in accordance with {@link #getChannel()}
     * otherwise a Runtime Exception will be produced by Netty
     */
    public abstract EventLoopGroup getEventLoopGroup();

    @Override public Channel createChannel(SocketAddress socketAddress, ChannelHandler channelHandler) throws PlcConnectionException {
        try {
            final EventLoopGroup workerGroup = getEventLoopGroup();

            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(workerGroup);
            bootstrap.channel(getChannel());
            bootstrap.option(ChannelOption.SO_KEEPALIVE, true);
            bootstrap.option(ChannelOption.TCP_NODELAY, true);
            // TODO we should use an explicit (configurable?) timeout here
            // bootstrap.option(ChannelOption.CONNECT_TIMEOUT_MILLIS, 1000);
            bootstrap.handler(channelHandler);
            // Start the client.
            final ChannelFuture f = bootstrap.connect(socketAddress);
            f.addListener(new GenericFutureListener<Future<? super Void>>() {
                @Override public void operationComplete(Future<? super Void> future) throws Exception {
                    if (!future.isSuccess()) {
                        logger.info("Unable to connect, shutting down worker thread.");
                        workerGroup.shutdownGracefully();
                    }
                }
            });
            // Wait for sync
            f.sync();
            f.awaitUninterruptibly(); // jf: unsure if we need that
            // Wait till the session is finished initializing.
            return f.channel();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException("Error creating channel.", e);
        } catch (Exception e) {
            throw new PlcConnectionException("Error creating channel.", e);
        }
    }

    @Override public void ping() throws PlcException {

    }
}
