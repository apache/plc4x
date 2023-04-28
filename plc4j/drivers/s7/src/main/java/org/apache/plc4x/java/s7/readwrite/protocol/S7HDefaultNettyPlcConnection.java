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
package org.apache.plc4x.java.s7.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.channel.Channel;
import io.netty.channel.embedded.EmbeddedChannel;
import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.DefaultNettyPlcConnection;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.events.CloseConnectionEvent;
import org.apache.plc4x.java.spi.events.DisconnectEvent;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

/**
 * @author cgarcia
 */
public class S7HDefaultNettyPlcConnection extends DefaultNettyPlcConnection implements Runnable {

    private static final Logger logger = LoggerFactory.getLogger(S7HDefaultNettyPlcConnection.class);

    private static final String MULTIPLEXOR = "MULTIPLEXOR";
    private Boolean closed = false;

    private ScheduledFuture<?> scf = null;

    protected final ChannelFactory secondaryChannelFactory;
    protected Channel primary_channel = null;
    protected Channel secondary_channel = null;
    protected final MessageToMessageCodec<ByteBuf, ByteBuf> s7hmux;

    protected int slice_ping = 0;
    protected int slice_retry_time = 0;

    public S7HDefaultNettyPlcConnection(boolean canRead,
                                        boolean canWrite,
                                        boolean canSubscribe,
                                        boolean canBrowse,
                                        PlcTagHandler tagHandler,
                                        PlcValueHandler valueHandler,
                                        Configuration configuration,
                                        ChannelFactory channelFactory,
                                        ChannelFactory secondaryChannelFactory,
                                        boolean awaitSessionSetupComplete,
                                        boolean awaitSessionDisconnectComplete,
                                        boolean awaitSessionDiscoverComplete,
                                        ProtocolStackConfigurer<TPKTPacket> stackConfigurer,
                                        BaseOptimizer optimizer,
                                        PlcAuthentication authentication) {
        super(canRead,
            canWrite,
            canSubscribe,
            canBrowse,
            tagHandler,
            valueHandler,
            configuration,
            channelFactory,
            awaitSessionSetupComplete,
            awaitSessionDisconnectComplete,
            awaitSessionDiscoverComplete,
            stackConfigurer,
            optimizer,
            authentication);
        this.secondaryChannelFactory = secondaryChannelFactory;
        this.s7hmux = new S7HMuxImpl();
    }

    @Override
    public void connect() throws PlcConnectionException {
        try {
            // As we don't just want to wait till the connection is established,
            // define a future we can use to signal back that the s7 session is
            // finished initializing.
            CompletableFuture<Void> sessionSetupCompleteFuture = new CompletableFuture<>();
            CompletableFuture<Configuration> sessionDiscoveredCompleteFuture = new CompletableFuture<>();

            if (channelFactory == null) {
                throw new PlcConnectionException("No primary channel factory provided");
            }

            // Inject the configuration
            ConfigurationFactory.configure(configuration, channelFactory);

            if (secondaryChannelFactory != null)
                ConfigurationFactory.configure(configuration, secondaryChannelFactory);

            channel = new EmbeddedChannel(getChannelHandler(sessionSetupCompleteFuture, sessionDisconnectCompleteFuture, sessionDiscoveredCompleteFuture));
            channel.pipeline().addFirst(s7hmux);
            ((S7HMux) s7hmux).setEmbededhannel(channel);
            //channel.pipeline().addFirst((new LoggingHandler(LogLevel.INFO))); 
            /*
            channel.closeFuture().addListener(future -> {
                if (!sessionSetupCompleteFuture.isDone()) {
                    sessionSetupCompleteFuture.completeExceptionally(
                        new PlcIoException("Connection terminated by remote"));
                }
            });
            */
            doPrimaryTcpConnections();

            if (secondaryChannelFactory != null)
                doSecondaryTcpConnections();

            //If it is not possible to generate a TCP connection.
            //Safety shutdown all executors in the channels.
            if (primary_channel == null)
                if (secondary_channel == null) {
                    sendChannelDisconectEvent();
                    throw new PlcConnectionException("Connection is not possible.");
                }

            scf = channel.eventLoop().scheduleWithFixedDelay(this, 1, 1, TimeUnit.SECONDS);
            
            /*            
            primary_channel.closeFuture().addListener(future -> {
                if (!sessionDiscoveredCompleteFuture.isDone()) {
                    //Do Nothing
                    try {
                        sessionDiscoveredCompleteFuture.complete(null);
                    } catch (Exception e) {
                        //Do Nothing
                    }

                }
            });            
            */


            // Send an event to the pipeline telling the Protocol filters what's going on.
            sendChannelCreatedEvent();

            // Wait till the connection is established.
            if (awaitSessionSetupComplete) {
                sessionSetupCompleteFuture.get();
            }
            //channel.pipeline().write(new ConnectedEvent());
            // Set the connection to "connected"
            connected = true;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException(e);
        } catch (ExecutionException e) {
            throw new PlcConnectionException(e);
        }
    }

    @Override
    public void close() throws PlcConnectionException {
        if (closed) return;
        try {
            scf.cancel(true);
        } catch (Exception ex) {
            logger.info(ex.toString());
        }
        if (primary_channel != null)
            if (primary_channel.isActive()) {
                try {
                    primary_channel.pipeline().remove(MULTIPLEXOR);
                    primary_channel.pipeline().fireUserEventTriggered(new CloseConnectionEvent());
                    primary_channel.eventLoop().shutdownGracefully();

                } catch (Exception ex) {
                    logger.info(ex.toString());
                }
            }

        if (secondary_channel != null)
            if (secondary_channel.isActive()) {
                secondary_channel.pipeline().remove(MULTIPLEXOR);
                secondary_channel.pipeline().fireUserEventTriggered(new CloseConnectionEvent());
                secondary_channel.eventLoop().shutdownGracefully();
            }

        channel.pipeline().fireUserEventTriggered(new DisconnectEvent());
        closed = true;
    }


    public void doPrimaryTcpConnections() {
        try {
            primary_channel = channelFactory.createChannel(new LoggingHandler(LogLevel.TRACE));
        } catch (Exception ex) {
            logger.info(ex.toString());
        }
        if (primary_channel != null)
            if (primary_channel.isActive())
                primary_channel.pipeline().addFirst(MULTIPLEXOR, s7hmux);
        ((S7HMux) s7hmux).setPrimaryChannel(primary_channel);
    }

    public void doSecondaryTcpConnections() {
        try {
            secondary_channel = secondaryChannelFactory.createChannel(new LoggingHandler(LogLevel.TRACE));
        } catch (Exception ex) {
            logger.info(ex.toString());
        }
        if (secondary_channel != null)
            if (secondary_channel.isActive())
                secondary_channel.pipeline().addFirst(MULTIPLEXOR, s7hmux);
        ((S7HMux) s7hmux).setSecondaryChannel(secondary_channel);
    }

    /*
     * All handlers on the channel are notified that a disconnection has been
     * generated, generally during the first connection.
     * In this way, a controlled shutdown of the execution services is achieved.
     * The user application must take the measures to make the connection again.
     */
    protected void sendChannelDisconectEvent() {
        logger.trace("Channels was not created, firing DisconnectEvent Event");
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new DisconnectEvent());
    }

    /*
     * To avoid creating new tasks associated with the supervision of the
     * driver, I execute these tasks in a few slices of time. This allows me
     * to keep a simple control of the state machines.
     */
    @Override
    public void run() {
        if (slice_retry_time >= channel.attr(S7HMuxImpl.RETRY_TIME).get()) {
            if (primary_channel != null) {
                if (!primary_channel.isActive()) {
                    logger.info("Creating primary connection.");
                    primary_channel.eventLoop().shutdownGracefully();
                    doPrimaryTcpConnections();
                }
            } else {
                logger.info("Creating first primary connection.");
                doPrimaryTcpConnections();
            }

            if (secondary_channel != null) {
                if (!secondary_channel.isActive()) {
                    logger.info("Creating secondary connection.");
                    secondary_channel.eventLoop().shutdownGracefully();
                    doSecondaryTcpConnections();
                }
            } else {
                if (secondaryChannelFactory != null) {
                    logger.info("Creating first secondary connection.");
                    doSecondaryTcpConnections();
                }
            }
            slice_retry_time = 0;
        }

        if ((channel.attr(S7HMuxImpl.IS_PING_ACTIVE).get())) {
            if (slice_ping >= channel.attr(S7HMuxImpl.PING_TIME).get()) {
                ping();
                slice_ping = 0;
            }
        } else slice_ping = 0;

        slice_retry_time++;
        slice_ping++;
    }

    @Override
    public CompletableFuture<Void> ping() {
        if (channel.attr(S7HMuxImpl.IS_CONNECTED).get()) {
            channel.eventLoop().execute(() -> {
                PlcReadRequest.Builder builder = readRequestBuilder();

                builder.addTagAddress("value", "%MX1.0:BOOL");
                PlcReadRequest readRequest = builder.build();
                try {
                    PlcReadResponse readResponse = readRequest.execute().get(2, TimeUnit.SECONDS);
                    logger.info("PING: " + readResponse.getResponseCode("value"));
                } catch (Exception e) {
                    logger.info("PING: {}", e.getMessage(), e);
                }
            });
        }
        return null;
    }

}
