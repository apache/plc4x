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
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;
import java.util.Properties;

/**
 * Adapter with sensible defaults for a Netty Based Channel Factory.
 * <p>
 * By Default Nettys {@link NioEventLoopGroup} is used.
 * Transports which have to use a different EventLoopGroup have to override {@link #getEventLoopGroup()}.
 */
public abstract class NettyChannelFactory implements ChannelFactory {

    private static final Logger logger = LoggerFactory.getLogger(NettyChannelFactory.class);
    private static final int PING_TIMEOUT_MS = 1_000;

    /**
     * TODO should be removed together with the Construcotr.
     */
    private SocketAddress address;
    private Properties properties;

    /**
     * @Deprecated Only there for Retrofit
     */
    @Deprecated
    public NettyChannelFactory(SocketAddress address) {
        this.address = address;
    }

    public NettyChannelFactory() {
        // Default Constructor to Use
    }

    /**
     * Channel to Use, e.g. NiO, EiO
     */
    public abstract Class<? extends Channel> getChannel();

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
     * Transports which have to use a different EventLoopGroup have to override {@link #getEventLoopGroup()}.
     */
    public EventLoopGroup getEventLoopGroup() {
        return new NioEventLoopGroup();
    }

    /**
     * @Deprecated use {@link #createChannel(SocketAddress, ChannelHandler)} instead.
     */
    @Deprecated
    @Override
    public Channel createChannel(ChannelHandler channelHandler) throws PlcConnectionException {
        if (this.address == null) {
            throw new IllegalStateException("This Method should only be used with the constructor which takes an Address");
        }
        return this.createChannel(address, channelHandler);
    }

    @Override
    public Channel createChannel(SocketAddress socketAddress, ChannelHandler channelHandler) throws PlcConnectionException {
        if (this.address == null) {
            this.address = socketAddress;
        }
        try {
            final EventLoopGroup workerGroup = getEventLoopGroup();

            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(workerGroup);
            bootstrap.channel(getChannel());
            // Callback to allow subclasses to modify the Bootstrap
            configureBootstrap(bootstrap);
            bootstrap.handler(channelHandler);
            // Start the client.
            final ChannelFuture f = bootstrap.connect(socketAddress);
            f.addListener(future -> {
                if (!future.isSuccess()) {
                    logger.info("Unable to connect, shutting down worker thread.");
                    workerGroup.shutdownGracefully();
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

    @Deprecated
    public InetAddress getAddress() {
        return ((InetSocketAddress) this.address).getAddress();
    }

    @Deprecated
    public int getPort() {
        return ((InetSocketAddress) this.address).getPort();
    }

    // TODO do we want to keep this like that?
    @Override
    public void ping() throws PlcException {
        // TODO: Replace this check with a more accurate one ...
        InetSocketAddress address = new InetSocketAddress(getAddress(), getPort());
        try (Socket s = new Socket()) {
            s.connect(address, PING_TIMEOUT_MS);
            // TODO keep the address for a (timely) next request???
            s.setReuseAddress(true);
        } catch (Exception e) {
            throw new PlcConnectionException("Unable to ping remote host");
        }
    }

    public Properties getProperties() {
        // Null Safety for older implementations
        if (properties == null) {
            return new Properties();
        }
        return properties;
    }

    public void setProperties(Properties properties) {
        this.properties = properties;
    }

    protected String getProperty(String key) {
        return ((String) getProperties().get(key));
    }

    protected boolean hasProperty(String key) {
        return getProperties().contains(key);
    }

    protected String getPropertyOrDefault(String key, String defaultValue) {
        return getProperties().getProperty(key, defaultValue);
    }

}
