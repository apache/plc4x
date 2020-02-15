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

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.util.HashedWheelTimer;
import io.netty.util.Timer;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcIoException;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.events.CloseConnectionEvent;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

public class DefaultNettyPlcConnection extends AbstractPlcConnection {

    /**
     * a {@link HashedWheelTimer} shall be only instantiated once.
     */
    // TODO: maybe find a way to make this configurable per jvm
    protected final static Timer timer = new HashedWheelTimer();
    private static final Logger logger = LoggerFactory.getLogger(DefaultNettyPlcConnection.class);

    protected final Configuration configuration;
    protected final ChannelFactory channelFactory;
    protected final boolean awaitSessionSetupComplete;
    protected final ProtocolStackConfigurer stackConfigurer;

    protected Channel channel;
    protected boolean connected;

    public DefaultNettyPlcConnection(boolean canRead, boolean canWrite, boolean canSubscribe,
                                     PlcFieldHandler fieldHandler, Configuration configuration,
                                     ChannelFactory channelFactory, boolean awaitSessionSetupComplete,
                                     ProtocolStackConfigurer stackConfigurer, BaseOptimizer optimizer) {
        super(canRead, canWrite, canSubscribe, fieldHandler, optimizer);
        this.configuration = configuration;
        this.channelFactory = channelFactory;
        this.awaitSessionSetupComplete = awaitSessionSetupComplete;
        this.stackConfigurer = stackConfigurer;

        this.connected = false;
    }

    @Override
    public void connect() throws PlcConnectionException {
        try {
            // As we don't just want to wait till the connection is established,
            // define a future we can use to signal back that the s7 session is
            // finished initializing.
            CompletableFuture<Void> sessionSetupCompleteFuture = new CompletableFuture<>();

            if(channelFactory == null) {
                throw new PlcConnectionException("No channel factory provided");
            }

            // Inject the configuration
            ConfigurationFactory.configure(configuration, channelFactory);

            // Have the channel factory create a new channel instance.
            channel = channelFactory.createChannel(getChannelHandler(sessionSetupCompleteFuture));
            channel.closeFuture().addListener(future -> {
                if (!sessionSetupCompleteFuture.isDone()) {
                    sessionSetupCompleteFuture.completeExceptionally(
                        new PlcIoException("Connection terminated by remote"));
                }
            });
            // Send an event to the pipeline telling the Protocol filters what's going on.
            sendChannelCreatedEvent();

            // Wait till the connection is established.
            if (awaitSessionSetupComplete) {
                sessionSetupCompleteFuture.get();
            }

            // Set the connection to "connected"
            connected = true;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException(e);
        } catch (ExecutionException e) {
            throw new PlcConnectionException(e);
        }
    }

    /*@Override
    public CompletableFuture<Void> ping() {
        CompletableFuture<Void> future = new CompletableFuture<>();
        try {
            // Relay the actual pinging to the channel factory ...
            channelFactory.ping();
            // If we got here, the ping was successful.
            future.complete(null);
        } catch (PlcException e) {
            // If we got here, something went wrong.
            future.completeExceptionally(e);
        }
        return future;
    }*/

    @Override
    public void close() throws PlcConnectionException {
        // TODO call protocols close method
        channel.pipeline().fireUserEventTriggered(new CloseConnectionEvent());
        // Close channel
        channel.close().awaitUninterruptibly();
        channel = null;
        connected = false;
    }

    /**
     * Check if the communication channel is active (channel.isActive()) and the driver for a given protocol
     * has finished establishing the connection.
     */
    @Override
    public boolean isConnected() {
        return connected && channel.isActive();
    }

    public Channel getChannel() {
        return channel;
    }

    public ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        if (stackConfigurer == null) {
            throw new IllegalStateException("No Protocol Stack Configurer is given!");
        }
        /*if (factory == null) {
            throw new IllegalStateException("No Instance Factory is Present!");
        }*/
        return new ChannelInitializer<Channel>() {
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
                // Initialize via Transport Layer
                channelFactory.initializePipeline(pipeline);
                // Initialize Protocol Layer
                setProtocol(stackConfigurer.configurePipeline(configuration, pipeline));
            }
        };
    }

    protected void sendChannelCreatedEvent() {
        logger.trace("Channel was created, firing ChannelCreated Event");
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

}
