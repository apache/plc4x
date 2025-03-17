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
import java.util.List;
import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcPingResponse;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.DefaultNettyPlcConnection;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.events.CloseConnectionEvent;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.spi.events.DisconnectEvent;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.*;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcSubscriptionRequest;

/**
 * This object generates the main connection and includes the management
 * of multiple connections to the PLC, in the case of S7H there
 * are two connections.
 * Here the reference to the multiplexer object (S7HMuxImpl) is maintained,
 * which is in charge of managing the physical TCP connections.
 * <p>
 * TODO: It should be able to run in the "channel" executor.
 * You should be able to remove the "executor".
 */
public class S7HPlcConnection extends DefaultNettyPlcConnection implements Runnable {

    private static final Logger logger = LoggerFactory.getLogger(S7HPlcConnection.class);

    private static final String MULTIPLEXER = "MULTIPLEXER";
    private Boolean closed = false;

    private ScheduledFuture<?> scf = null;
    private final ScheduledExecutorService executor = Executors.newSingleThreadScheduledExecutor(new BasicThreadFactory.Builder()
        .namingPattern("plc4x-s7ha-thread-%d")
        .daemon(true)
        .priority(Thread.MAX_PRIORITY)
        .build());

    protected final ChannelFactory secondaryChannelFactory;
    protected Channel primaryChannel = null;
    protected Channel secondaryChannel = null;
    protected final MessageToMessageCodec<ByteBuf, ByteBuf> s7hmux;

    protected int slicePing = 0;
    protected int sliceRetryTime = 0;

    public S7HPlcConnection(
        boolean canPing,
        boolean canRead,
        boolean canWrite,
        boolean canSubscribe,
        boolean canBrowse,
        PlcValueHandler valueHandler,
        PlcConnectionConfiguration configuration,
        ChannelFactory channelFactory,
        ChannelFactory secondaryChannelFactory,
        boolean fireDiscoverEvent,
        boolean awaitSessionSetupComplete,
        boolean awaitSessionDisconnectComplete,
        boolean awaitSessionDiscoverComplete,
        ProtocolStackConfigurer<TPKTPacket> stackConfigurer,
        BaseOptimizer optimizer,
        PlcAuthentication authentication) {
        super(canPing,
            canRead,
            canWrite,
            canSubscribe,
            canBrowse,
            valueHandler,
            configuration,
            channelFactory,
            fireDiscoverEvent,
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
            CompletableFuture<PlcConnectionConfiguration> sessionDiscoveredCompleteFuture = new CompletableFuture<>();

            if (channelFactory == null) {
                throw new PlcConnectionException("No primary channel factory provided");
            }

            // Inject the configuration
            ConfigurationFactory.configure(configuration, channelFactory);

            if (secondaryChannelFactory != null) {
                ConfigurationFactory.configure(configuration, secondaryChannelFactory);
            }

            if (null == channel) {
                channel = new EmbeddedChannel(
                    getChannelHandler(sessionSetupCompleteFuture,
                        sessionDisconnectCompleteFuture,
                        sessionDiscoveredCompleteFuture));

                channel.pipeline().addFirst("Multiplexor", s7hmux);

            }

            ((S7HMux) s7hmux).setEmbededhannel(channel, configuration);
//            channel.pipeline().addFirst((new LoggingHandler("CEOS"))); 
            /*
            channel.closeFuture().addListener(future -> {
                if (!sessionSetupCompleteFuture.isDone()) {
                    sessionSetupCompleteFuture.completeExceptionally(
                        new PlcIoException("Connection terminated by remote"));
                }
            });
            */
            doPrimaryTcpConnections();

            if (secondaryChannelFactory != null) {
                doSecondaryTcpConnections();
            }

            //If it is not possible to generate a TCP connection.
            //Safety shutdown all executors in the channels.
            if (primaryChannel == null) {
                if (secondaryChannel == null) {
                    sendChannelDisconectEvent();
                    throw new PlcConnectionException("Connection is not possible.");
                }
            }

            scf = executor.scheduleAtFixedRate(this, 1, 1, TimeUnit.SECONDS); 
            
            /*            
            primary_channel.closeFuture().addListener(future -> {/watch?v=TmENMZFUU_0&list=RDlBlx1JffMQ4&index=27
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
            //((EmbeddedChannel) channel).runPendingTasks();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException(e);
        } catch (ExecutionException e) {
            throw new PlcConnectionException(e);
        }
    }

    @Override
    public void close() throws PlcConnectionException {
        if (closed) {
            return;
        }
        try {
            scf.cancel(true);
        } catch (Exception ex) {
            logger.info(ex.toString());
        }
        if (primaryChannel != null) {
            if (primaryChannel.isActive()) {
                try {
                    primaryChannel.pipeline().remove(MULTIPLEXER);
                    primaryChannel.pipeline().fireUserEventTriggered(new CloseConnectionEvent());
                    primaryChannel.eventLoop().shutdownGracefully();
                } catch (Exception ex) {
                    logger.info(ex.toString());
                }
            }
        }

        if (secondaryChannel != null) {
            if (secondaryChannel.isActive()) {
                secondaryChannel.pipeline().remove(MULTIPLEXER);
                secondaryChannel.pipeline().fireUserEventTriggered(new CloseConnectionEvent());
                secondaryChannel.eventLoop().shutdownGracefully();
            }
        }

        channel.pipeline().fireUserEventTriggered(new DisconnectEvent());
        scf.cancel(true);
        executor.shutdown();
        closed = true;
    }

    @Override
    public boolean isConnected() {
        return channel.attr(S7HMuxImpl.IS_CONNECTED).get();
    }

    /**
     * Subscriptions are only supported in a small subset of the S7 devices.
     *
     * @return true, if the device supports subscriptions.
     */
    @Override
    public boolean isSubscribeSupported() {
        Plc4xProtocolBase<?> protocol = getProtocol();
        S7DriverContext s7driverContext = (S7DriverContext) protocol.getDriverContext();
        return (s7driverContext.getControllerType() == ControllerType.S7_300) || (s7driverContext.getControllerType() == ControllerType.S7_400);
    }

    public void doPrimaryTcpConnections() {
        try {
            primaryChannel = channelFactory.createChannel(new LoggingHandler(LogLevel.TRACE));
        } catch (Exception ex) {
            primaryChannel = null;
            logger.info(ex.toString());
        }
        if (primaryChannel != null) {
            if (primaryChannel.isActive()) {
                primaryChannel.pipeline().addFirst(MULTIPLEXER, s7hmux);
                ((S7HMux) s7hmux).setPrimaryChannel(primaryChannel);
            }
        }
    }

    public void doSecondaryTcpConnections() {
        try {
            secondaryChannel = secondaryChannelFactory.createChannel(new LoggingHandler(LogLevel.TRACE));
        } catch (Exception ex) {
            secondaryChannel = null;
            logger.info(ex.toString());
        }
        if (secondaryChannel != null) {
            if (secondaryChannel.isActive()) {
                secondaryChannel.pipeline().addFirst(MULTIPLEXER, s7hmux);
                ((S7HMux) s7hmux).setSecondaryChannel(secondaryChannel);
            }
        }
    }

    /*
     * All handlers on the channel are notified that a disconnection has been
     * generated, generally during the first connection.
     * In this way, a controlled shutdown of the execution services is achieved.
     * The user application must take the measures to make the connection again.
     */
    protected void sendChannelDisconectEvent() {
        logger.trace("Channel was not created, firing DisconnectEvent Event");
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new DisconnectEvent());
    }

    /*
     * To avoid creating new tasks associated with the supervision of the
     * driver, I execute these tasks in a few slices of time. This allows me
     * to keep a simple control of the state machines.
     * Verifies disconnection of the channel, which notifies the client
     * that they are listening.
     */
    @Override
    public void run() {

        /*
         * Here a driver reconnection is detected. If so, a connection event
         * is generated, so the user application can generate its requests again.
         */
        if (channel.attr(S7HMuxImpl.WAS_CONNECTED).get() &&
            channel.attr(S7HMuxImpl.IS_CONNECTED).get()) {
            channel.attr(S7HMuxImpl.WAS_CONNECTED).set(false);
            channel.pipeline().fireUserEventTriggered(new ConnectedEvent());
        }

        /*
         * Execute the ping towards the PLC. Its purpose is to keep
         * the TCP channel active, for very long sampling periods.
         */
        if (channel.attr(S7HMuxImpl.IS_PING_ACTIVE).get()) {
            if (slicePing >= channel.attr(S7HMuxImpl.PING_TIME).get()) {
                ping();
                slicePing = 0;
            }
            slicePing++;
        } else {
            slicePing = 0;
        }

        /*
         * Here it is verified if the channels are active. If they are not,
         * the connection is created again.
         * For H type systems, at least one of the connections must be active,
         * and if the other is down, the connection must be lifted.
         * S7HMuxImpl, includes the switching logic between the TCP links.
         */
        if (sliceRetryTime >= channel.attr(S7HMuxImpl.RETRY_TIME).get()) {

            if (primaryChannel != null) {
                if (!primaryChannel.isActive()) {
                    logger.info("Creating primary connection.");
                    primaryChannel.eventLoop().shutdownGracefully();
                    doPrimaryTcpConnections();
                } else if (null == secondaryChannel) {
                    if (channel.attr(S7HMuxImpl.WAS_CONNECTED).get() &&
                        !channel.attr(S7HMuxImpl.IS_CONNECTED).get()) {
                        logger.info("Reconnecting primary channel.");
                        if (null != ((S7HMux) s7hmux).getTCPChannel()) {
                            ((S7HMux) s7hmux).setPrimaryChannel(primaryChannel);
                        }
                    }
                }
            } else {
                logger.info("Creating first primary connection.");
                doPrimaryTcpConnections();
            }

            if (secondaryChannel != null) {
                if (!secondaryChannel.isActive()) {
                    logger.info("Creating secondary connection.");
                    secondaryChannel.eventLoop().shutdownGracefully();
                    doSecondaryTcpConnections();
                } else if (null == primaryChannel) {
                    if ((channel.attr(S7HMuxImpl.WAS_CONNECTED).get()) &&
                        (!channel.attr(S7HMuxImpl.IS_CONNECTED).get())) {
                        logger.info("Reconnecting secondary channel.");
                        if (null != ((S7HMux) s7hmux).getTCPChannel()) {
                            ((S7HMux) s7hmux).setSecondaryChannel(secondaryChannel);
                        }
                    }
                }
            } else {
                if (secondaryChannelFactory != null) {
                    logger.info("Creating first secondary connection.");
                    doSecondaryTcpConnections();
                }
            }
            sliceRetryTime = 0;
        }

        if (channel.attr(S7HMuxImpl.RETRY_TIME).get() > 0) {
            sliceRetryTime++;
        }

        connected = channel.attr(S7HMuxImpl.IS_CONNECTED).get();
        logger.trace("*************************************************\r\n"
            + "INSTAMCIA PRIMARIO      : " + ((null == primaryChannel) ? "NULL" : primaryChannel.toString()) + "\r\n"
            + "ACTIVO PRIMARIO         : " + ((null == primaryChannel) ? "NULL" : primaryChannel.isActive()) + "\r\n"
            + "INSTAMCIA SECUNDARIO    : " + ((null == secondaryChannel) ? "NULL" : secondaryChannel.toString()) + "\r\n"
            + "ACTIVO SECUNDARIO       : " + ((null == secondaryChannel) ? "NULL" : secondaryChannel.isActive()) + "\r\n"
            + "CANAL CONECTADO?        : " + channel.attr(S7HMuxImpl.IS_CONNECTED).get() + "\r\n"
            + "CANAL ESTUVO CONECTADO? : " + channel.attr(S7HMuxImpl.WAS_CONNECTED).get() + "\r\n"
            + "CONTADORES              : " + slicePing + " : " + sliceRetryTime + "\r\n"
            + "*************************************************");
    }

    /*
     * PING to PLC.
     */
    @Override
    public CompletableFuture<? extends PlcPingResponse> ping() {
        if (channel.attr(S7HMuxImpl.IS_CONNECTED).get()) {

            //channel.eventLoop().execute(()->{
            executor.execute(() -> {
                PlcReadRequest.Builder builder = readRequestBuilder();

                builder.addTagAddress("value", "%MX1.0:BOOL");
                PlcReadRequest readRequest = builder.build();
                try {
                    PlcReadResponse readResponse = readRequest.execute().get(2, TimeUnit.SECONDS);
                    logger.debug("PING: " + readResponse.getResponseCode("value"));
                } catch (Exception ex) {
                    logger.info("PING: " + ex);
                }
            });
        }
        return null;
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        if (!isSubscribeSupported()) {
            throw new PlcUnsupportedOperationException("The connection does not support subscription");
        }
        return new S7PlcSubscriptionRequest.Builder(this, getPlcTagHandler());        
    }        

}
