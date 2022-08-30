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
package org.apache.plc4x.java.spi.connection;

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.SocketAddress;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Adapter with sensible defaults for a Netty Based Channel Factory.
 * <p>
 * By Default Nettys {@link NioEventLoopGroup} is used.
 * Transports which have to use a different EventLoopGroup have to override {@link #getEventLoopGroup()}.
 */
public abstract class NettyChannelFactory implements ChannelFactory {

    private static final Logger logger = LoggerFactory.getLogger(NettyChannelFactory.class);

    private final Map<Channel, EventLoopGroup> eventLoops = new ConcurrentHashMap<>();

    /**
     * TODO should be removed together with the Constructor.
     */
    private SocketAddress address;

    protected NettyChannelFactory(SocketAddress address) {
        this.address = address;
    }

    /**
     * Channel to Use, e.g. NiO, EiO
     */
    public abstract Class<? extends Channel> getChannel();

    /**
     * We need to be able to override this in the TestChanelFactory
     *
     * @return the Bootstrap instance we will be using to initialize the channel.
     */
    protected Bootstrap createBootstrap() {
        return new Bootstrap();
    }

    /**
     * This Method is used to modify the Bootstrap Element of Netty, if one wishes to do so.
     * E.g. for Protocol Specific extension.
     * For TCP e.g.
     * <code>
     * bootstrap.option(ChannelOption.SO_KEEPALIVE, true);
     * bootstrap.option(ChannelOption.TCP_NODELAY, true);
     * </code>
     */
    public abstract void configureBootstrap(Bootstrap bootstrap);

    /**
     * Event Loop Group to use.
     * Has to be in accordance with {@link #getChannel()}
     * otherwise a Runtime Exception will be produced by Netty
     * <p>
     * By Default Nettys {@link NioEventLoopGroup} is used.
     * Transports which have to use a different EventLoopGroup have to override {#getEventLoopGroup()}.
     */
    public EventLoopGroup getEventLoopGroup() {
        return new NioEventLoopGroup();
    }

    @Override
    public Channel createChannel(ChannelHandler channelHandler) throws PlcConnectionException {
        try {
            Bootstrap bootstrap = createBootstrap();

            EventLoopGroup workerGroup = getEventLoopGroup();
            if (workerGroup != null) {
                bootstrap.group(workerGroup);
            }

            bootstrap.channel(getChannel());
            // Callback to allow subclasses to modify the Bootstrap

            configureBootstrap(bootstrap);
            bootstrap.handler(channelHandler);
            // Start the client.
            final ChannelFuture f = bootstrap.connect(address);
            f.addListener(future -> {
                if (!future.isSuccess()) {
                    logger.info("Unable to connect, shutting down worker thread.");
                    if (workerGroup != null) {
                        workerGroup.shutdownGracefully();
                    }
                }
            });

            final Channel channel = f.channel();

            if (workerGroup != null) {
                // Shut down the workerGroup when channel closing to avoid open too many files
                channel.closeFuture().addListener(future -> workerGroup.shutdownGracefully());
                // Add to event-loop group
                eventLoops.put(channel, workerGroup);
            }

            // It seems the embedded channel operates differently.
            // Intentionally using the class name as we don't want to require a
            // hard dependency on the test-channel.
            if (!"Plc4xEmbeddedChannel".equals(channel.getClass().getSimpleName())) {
                // Wait for sync
                f.sync();
                // Wait till the session is finished initializing.
                f.awaitUninterruptibly(); // jf: unsure if we need that
            }

            return channel;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException("Error creating channel.", e);
        } catch (Throwable t) {
            throw new PlcConnectionException("Error creating channel.", t);
        }
    }

    @Override
    public void closeEventLoopForChannel(Channel channel) {
        if (eventLoops.containsKey(channel)) {
            logger.info("Channel is closed, closing worker Group also");
            EventLoopGroup eventExecutors = eventLoops.get(channel);
            eventLoops.remove(channel);
            eventExecutors.shutdownGracefully().awaitUninterruptibly();
            logger.info("Worker Group was closed successfully!");
        } else {
            logger.warn("Trying to remove EventLoop for Channel {} but have none stored", channel);
        }
    }

}
