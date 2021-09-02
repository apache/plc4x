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

import io.netty.channel.*;
import io.netty.util.HashedWheelTimer;
import io.netty.util.Timer;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcIoException;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.events.*;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.apache.plc4x.java.api.value.PlcValueHandler;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

public class DefaultNettyPlcConnection extends AbstractPlcConnection implements ChannelExposingConnection {

    /**
     * a {@link HashedWheelTimer} shall be only instantiated once.
     */
    // TODO: maybe find a way to make this configurable per jvm
    protected final static Timer timer = new HashedWheelTimer();
    protected final static long DEFAULT_DISCONNECT_WAIT_TIME = 10000L;
    private static final Logger logger = LoggerFactory.getLogger(DefaultNettyPlcConnection.class);

    protected Configuration configuration;
    protected final ChannelFactory channelFactory;
    protected final boolean awaitSessionSetupComplete;
    protected final boolean awaitSessionDisconnectComplete;
    protected final boolean awaitSessionDiscoverComplete;
    protected final ProtocolStackConfigurer stackConfigurer;
    protected final CompletableFuture<Void> sessionDisconnectCompleteFuture = new CompletableFuture<>();

    protected Channel channel;
    protected boolean connected;

    public DefaultNettyPlcConnection(boolean canRead, boolean canWrite, boolean canSubscribe,
                                     PlcFieldHandler fieldHandler, PlcValueHandler valueHandler, Configuration configuration,
                                     ChannelFactory channelFactory, boolean awaitSessionSetupComplete,
                                     boolean awaitSessionDisconnectComplete, boolean awaitSessionDiscoverComplete, ProtocolStackConfigurer stackConfigurer, BaseOptimizer optimizer) {
        super(canRead, canWrite, canSubscribe, fieldHandler, valueHandler, optimizer);
        this.configuration = configuration;
        this.channelFactory = channelFactory;
        this.awaitSessionSetupComplete = awaitSessionSetupComplete;
        //Used to signal that a disconnect has completed while closing a connection.
        this.awaitSessionDisconnectComplete = awaitSessionDisconnectComplete;
        //Used to signal that discovery has been completed
        this.awaitSessionDiscoverComplete = awaitSessionDiscoverComplete;
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
            CompletableFuture<Configuration> sessionDiscoveredCompleteFuture = new CompletableFuture<>();

            if(channelFactory == null) {
                throw new PlcConnectionException("No channel factory provided");
            }

            // Inject the configuration
            ConfigurationFactory.configure(configuration, channelFactory);

            // Have the channel factory create a new channel instance.
            if (awaitSessionDiscoverComplete) {
                channel = channelFactory.createChannel(getChannelHandler(sessionSetupCompleteFuture, sessionDisconnectCompleteFuture, sessionDiscoveredCompleteFuture));
                channel.closeFuture().addListener(future -> {
                    if (!sessionDiscoveredCompleteFuture.isDone()) {
                        //Do Nothing
                        try {
                            sessionDiscoveredCompleteFuture.complete(null);
                        } catch (Exception e) {
                            //Do Nothing
                        }

                    }
                });
                channel.pipeline().fireUserEventTriggered(new DiscoverEvent());

                // Wait till the connection is established.
                sessionDiscoveredCompleteFuture.get();
            }

            channel = channelFactory.createChannel(getChannelHandler(sessionSetupCompleteFuture, sessionDisconnectCompleteFuture, sessionDiscoveredCompleteFuture));
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

    /**
     * Close the connection by firstly calling disconnect and waiting for a DisconnectedEvent occurs and then calling
     * Close() to finish up any other clean up.
     * @throws PlcConnectionException
     */
    @Override
    public void close() throws PlcConnectionException {
        logger.debug("Closing connection to PLC, await for disconnect = {}", awaitSessionDisconnectComplete);
        channel.pipeline().fireUserEventTriggered(new DisconnectEvent());
        try {
            if (awaitSessionDisconnectComplete) {
                sessionDisconnectCompleteFuture.get(DEFAULT_DISCONNECT_WAIT_TIME, TimeUnit.MILLISECONDS);
            }
        } catch (Exception e) {
            logger.error("Timeout while trying to close connection");
        }
        channel.pipeline().fireUserEventTriggered(new CloseConnectionEvent());
        channel.close().awaitUninterruptibly();

        if (!sessionDisconnectCompleteFuture.isDone()) {
            sessionDisconnectCompleteFuture.complete(null);
        }

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

    @Override
    public Channel getChannel() {
        return channel;
    }

    public ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture, CompletableFuture<Void> sessionDisconnectCompleteFuture, CompletableFuture<Configuration> sessionDiscoverCompleteFuture) {
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
                        } else if (evt instanceof DisconnectedEvent) {
                            sessionDisconnectCompleteFuture.complete(null);
                        } else if (evt instanceof DiscoveredEvent) {
                            sessionDiscoverCompleteFuture.complete(((DiscoveredEvent) evt).getConfiguration());
                        } else {
                            super.userEventTriggered(ctx, evt);
                        }
                    }
                });
                // Initialize via Transport Layer
                channelFactory.initializePipeline(pipeline);
                // Initialize Protocol Layer
                setProtocol(stackConfigurer.configurePipeline(configuration, pipeline, channelFactory.isPassive()));
            }
        };
    }

    protected void sendChannelCreatedEvent() {
        logger.trace("Channel was created, firing ChannelCreated Event");
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

}
