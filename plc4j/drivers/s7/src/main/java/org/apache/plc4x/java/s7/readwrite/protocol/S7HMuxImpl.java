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
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandler.Sharable;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.embedded.EmbeddedChannel;
import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.handler.timeout.ReadTimeoutHandler;
import io.netty.util.AttributeKey;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.s7.readwrite.configuration.S7Configuration;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.spi.events.DisconnectEvent;
import org.apache.plc4x.java.spi.events.DisconnectedEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.LocalTime;
import java.util.List;
import java.util.logging.Level;

/**
 * Implementation of a multiplexing channel, from an embedded channel to two
 * possible TCP connections, primary and secondary.
 * The objective is to allow connections to individual systems
 * with a two CP (PN CPUs, CP343-1, CP443-1 or similar), or H-type systems
 * (S7-400H or S7-1500H).
 * The user App must be in charge of restoring the requests or
 * subscriptions that it is requesting.
 */
@Sharable
public class S7HMuxImpl extends MessageToMessageCodec<ByteBuf, ByteBuf> implements S7HMux {

    private static final Logger logger = LoggerFactory.getLogger(S7HMuxImpl.class);

    /*
     * This attribute indicates to the other handlers that the channel is connected
     * or disconnected because a switch is being made between TCP channels or
     * both TCP channels are disconnected.
     * Default value: false
     */
    final static AttributeKey<Boolean> IS_CONNECTED = AttributeKey.valueOf("IS_CONNECTED");

    /*
     * This attribute indicates to the other handlers that the channel is connected
     * or disconnected because a switch is being made between TCP channels or
     * both TCP channels are disconnected.
     * Default value: false
     */
    final static AttributeKey<Boolean> WAS_CONNECTED = AttributeKey.valueOf("WAS_CONNECTED");


    /*
     * This attribute indicates to the other handlers which channel is being used,
     * this allows the request to be properly prepared.
     * For example, in the case of a CPU with two CPs, you should change
     * the "slot", in the case of H systems, you should change the "rack",
     * the correct values will be defined in the connection URL.
     * Default value: true
     */
    final static AttributeKey<Boolean> IS_PRIMARY = AttributeKey.valueOf("IS_PRIMARY");

    /*
     * This is the maximum waiting time for reading on the TCP channel.
     * As there is no traffic, it must be assumed that the connection with the
     * interlocutor was lost and it must be restarted.
     * When the channel is closed, the "fail over" is carried out
     * in case of having the secondary channel, or it is expected that it
     * will be restored automatically, which is done every 4 seconds.
     * Default value: 8 sec.
     */
    final static AttributeKey<Integer> READ_TIME_OUT = AttributeKey.valueOf("READ_TIME_OUT");

    /*
     * If your application requires sampling times greater than the
     * set "watchdog" time, it is important that the PING option is activated,
     * this will prevent the TCP channel from being closed unnecessarily.
     * Default value: false
     */
    final static AttributeKey<Boolean> IS_PING_ACTIVE = AttributeKey.valueOf("IS_PIN_ACTIVE");

    /*
     * Time value in seconds at which the execution of the PING will be scheduled.
     * Generally set by developer experience, but generally should be the same
     * as READ_TIME_OUT / 2.
     * Default value: -1
     */
    final static AttributeKey<Integer> PING_TIME = AttributeKey.valueOf("PING_TIME");

    /*
     * Time for supervision of TCP channels. If the channel is not active,
     * a safe stop of the EventLoop must be performed, to ensure that
     * no additional tasks are created.
     * Default value: 4
     */
    final static AttributeKey<Integer> RETRY_TIME = AttributeKey.valueOf("RETRY_TIME");

    ChannelHandlerContext embedCtx = null;
    protected Channel embededChannel = null;
    protected Channel tcpChannel = null;
    protected Channel primaryChannel = null;
    protected Channel secondaryChannel = null;

    /*
     * From S7ProtocolLogic
     * TODO: Evaluate if the "embed_ctx" is really required since we set
     *  the Embeded channel when we created it.
     */
    @Override
    protected void encode(ChannelHandlerContext ctx, ByteBuf outBB, List<Object> list) {
        if ((embedCtx == null) && (ctx.channel() instanceof EmbeddedChannel)) embedCtx = ctx;
        if ((tcpChannel != null) && (embedCtx == ctx)) {
            tcpChannel.writeAndFlush(outBB.copy());
        }
        list.add(outBB.copy());
    }

    /*
     * To S7ProtocolLogic
     * The information received here from the channel "tcp_channel" is sent to
     * the pipeline of the channel "embeded_channel"
     */
    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf inBB, List<Object> list) throws Exception {
        embedCtx.fireChannelRead(inBB.copy());
    }

    @Override
    public void channelRegistered(ChannelHandlerContext ctx) throws Exception {
        super.channelRegistered(ctx);
        logger.debug("channelRegistered: " + ctx.name());
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        super.exceptionCaught(ctx, cause);
        logger.debug("exceptionCaught: " + ctx.name());
    }

    @Override
    public void channelWritabilityChanged(ChannelHandlerContext ctx) throws Exception {
        super.channelWritabilityChanged(ctx);
        logger.debug("channelWritabilityChanged: " + ctx.name());
    }

    /*
     * The events detected here flow from the S7ProtocolLogic object.
     * Upon receiving the "ConnectedEvent" event, we must safely add the watchdog
     * to the pipeline of the "tcp_channel" connection.
     * The supervision time can be defined in the connection URL,
     * the default value being 0 secs disabling this function.
     * This value being defined experimentally, a typical value is 8 seconds.
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        logger.info(LocalTime.now() + " userEventTriggered: " + ctx.name() + " Event: " + evt);

        if (evt instanceof ConnectedEvent) {
            try {
                ChannelHandler watchdog = tcpChannel.pipeline().get("watchdog");
                if (watchdog != null) {
                    tcpChannel.pipeline().remove(watchdog);
                }

            } catch (Exception ex) {
                logger.info(ex.toString());
            }
            try {
                if ((embededChannel.attr(READ_TIME_OUT).get() > 0) &&
                    embededChannel.attr(IS_PING_ACTIVE).get())
                    tcpChannel.pipeline().addFirst("watchdog",
                        new ReadTimeoutHandler(embededChannel.attr(READ_TIME_OUT).get()));
                if (tcpChannel.isActive()) {
                    embededChannel.attr(IS_CONNECTED).set(true);
                } else {
                    embededChannel.attr(IS_CONNECTED).set(false);
                }
            } catch (Exception ex) {
                logger.info(ex.toString());
            }
        }

        if (evt instanceof DisconnectEvent) {
//            logger.info("userEventTriggered -> DisconnectEvent");
        }
        
        // trigger other event handlers after IS_CONNECTED was set
        super.userEventTriggered(ctx, evt);
    }

    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) throws Exception {
        super.channelReadComplete(ctx);
        logger.debug(LocalTime.now() + " channelReadComplete: " + ctx.name());
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        super.channelInactive(ctx);
        logger.debug("channelInactive: " + ctx.name());
    }

    @Override
    public void channelActive(ChannelHandlerContext ctx) throws Exception {
        super.channelActive(ctx);
        logger.debug("channelActive: " + ctx.name());
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        super.channelUnregistered(ctx);
        logger.debug("{} channelUnregistered: {}", LocalTime.now(), ctx.name());
        String strCanal = (tcpChannel == primaryChannel) ? "PRIMARY" : "SECONDARY";
        logger.info("Unregistered of channel: {}", strCanal);
        //TODO: If embedded channel is closed, we need close all channels
        if (ctx == embedCtx) return;

        if (tcpChannel == ctx.channel()) {
            embededChannel.attr(IS_CONNECTED).set(false);
            embededChannel.attr(WAS_CONNECTED).set(true);
            embededChannel.pipeline().fireUserEventTriggered(new DisconnectedEvent());
        }

        if (embedCtx != null)
            logger.info(embedCtx.executor().toString());

        if ((tcpChannel == primaryChannel) &&
            (primaryChannel == ctx.channel()))
            if ((!primaryChannel.isActive()) &&
                (secondaryChannel != null))
                if (secondaryChannel.isActive()) {
                    synchronized (tcpChannel) {
                        logger.info("Using secondary TCP channel.");
                        tcpChannel = secondaryChannel;
                        embededChannel.attr(IS_PRIMARY).set(false);
                        embededChannel.pipeline().fireUserEventTriggered(new ConnectEvent());
                    }
                }


        if ((tcpChannel == secondaryChannel) &&
            (secondaryChannel == ctx.channel()))
            if ((!secondaryChannel.isActive()) &&
                (primaryChannel != null))
                if (primaryChannel.isActive()) {
                    synchronized (tcpChannel) {
                        logger.info("Using primary TCP channel.");
                        tcpChannel = primaryChannel;
                        embededChannel.attr(IS_PRIMARY).set(true);
                        embededChannel.pipeline().fireUserEventTriggered(new ConnectEvent());
                    }
                }
    }

    @Override
    public void setEmbededhannel(Channel embeded_channel, PlcConnectionConfiguration configuration) {
        final S7Configuration conf = (S7Configuration) configuration;
        this.embededChannel = embeded_channel;
        this.embededChannel.attr(IS_CONNECTED).set(false);
        this.embededChannel.attr(WAS_CONNECTED).set(false);
        this.embededChannel.attr(IS_PRIMARY).set(true);

        //From the URL
        this.embededChannel.attr(READ_TIME_OUT).set(conf.getReadTimeout());
        this.embededChannel.attr(IS_PING_ACTIVE).set(conf.getPing());
        this.embededChannel.attr(PING_TIME).set(conf.getPingTime());
        this.embededChannel.attr(RETRY_TIME).set(conf.getRetryTime());
    }

    @Override
    public void setPrimaryChannel(Channel primary_channel) {
        if ((this.primaryChannel == null) && (tcpChannel == null)) {
            if (primary_channel != null) {
                this.primaryChannel = primary_channel;
                tcpChannel = primary_channel;
                embededChannel.attr(IS_PRIMARY).set(true);
            }
        } else if ((this.primaryChannel == null) ||
            ((tcpChannel == secondaryChannel)) && (tcpChannel.isActive())) {
            this.primaryChannel = primary_channel;

        } else if ((!this.primaryChannel.isActive()) && (tcpChannel == secondaryChannel)) {
            this.primaryChannel = primary_channel;

        } else if (((!this.primaryChannel.isActive()) && (tcpChannel == this.primaryChannel)) ||
            (primary_channel.isActive())) {
            synchronized (tcpChannel) {                
                tcpChannel.close();
                this.primaryChannel = primary_channel;
                tcpChannel = primary_channel;
                embededChannel.attr(IS_PRIMARY).set(true);

                if (tcpChannel.isActive()) {
                    logger.info("Reassigns the inactive primary channel and send ConnectEvent..");
                    embedCtx.fireUserEventTriggered(new ConnectEvent());
                }
            }
        } else if (primary_channel.isActive()) {
            synchronized (tcpChannel) {
                tcpChannel.close();
                this.primaryChannel = primary_channel;
                tcpChannel = primary_channel;
                embededChannel.attr(IS_PRIMARY).set(true);
                logger.info("Reassigns the primary channel and send ConnectEvent.");
                if (tcpChannel.isActive()) {
                    embedCtx.fireUserEventTriggered(new ConnectEvent());
                }
            }            
        }
    }

    @Override
    public void setSecondaryChannel(Channel secondary_channel) {
        if ((this.primaryChannel == null) && (tcpChannel == null)) {
            if (secondary_channel != null) {
                this.secondaryChannel = secondary_channel;
                tcpChannel = secondary_channel;
                embededChannel.attr(IS_PRIMARY).set(false);
            }

        } else if ((this.secondaryChannel == null) ||
            ((tcpChannel == primaryChannel)) && (tcpChannel.isActive())) {
            this.secondaryChannel = secondary_channel;

        } else if ((!this.secondaryChannel.isActive()) && (tcpChannel == primaryChannel)) {
            this.secondaryChannel = secondary_channel;

        } else if (((!this.secondaryChannel.isActive()) && (tcpChannel == this.secondaryChannel)) ||
            (secondary_channel.isActive())) {
            synchronized (tcpChannel) {
                tcpChannel.close();
                this.secondaryChannel = secondary_channel;
                tcpChannel = secondary_channel;
                embededChannel.attr(IS_PRIMARY).set(false);
            }
            if (tcpChannel.isActive()) {
                embedCtx.fireUserEventTriggered(new ConnectEvent());
            }
        }
    }

    @Override
    public Channel getTCPChannel() {
        return tcpChannel;
    }

}
