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
package org.apache.plc4x.java.s7.readwrite.connection;

import io.netty.buffer.ByteBuf;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler.Sharable;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.embedded.EmbeddedChannel;
import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.handler.timeout.ReadTimeoutHandler;
import io.netty.util.AttributeKey;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.DisconnectEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.LocalTime;
import java.util.List;

/**
 * Implementation of a multiplexing channel, from an embedded channel to two
 * possible TCP connections, primary and secondary.
 * The objective is to allow connections to individual systems
 * with a two CP (PN CPUs, CP343-1, CP443-1 or similar), or H-type systems
 * (S7-400H or S7-1500H).
 * <p>
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
    public final static AttributeKey<Boolean> IS_CONNECTED = AttributeKey.valueOf("IS_CONNECTED");

    /*
     * This attribute indicates to the other handlers which channel is being used,
     * this allows the request to be properly prepared.
     * For example, in the case of a CPU with two CPs, you should change
     * the "slot", in the case of H systems, you should change the "rack",
     * the correct values will be defined in the connection URL.
     * Default value: true
     */
    public final static AttributeKey<Boolean> IS_PRIMARY = AttributeKey.valueOf("IS_PRIMARY");

    /*
     * This is the maximum waiting time for reading on the TCP channel.
     * As there is no traffic, it must be assumed that the connection with the
     * interlocutor was lost and it must be restarted.
     * When the channel is closed, the "fail over" is carried out
     * in case of having the secondary channel, or it is expected that it
     * will be restored automatically, which is done every 4 seconds.
     * Default value: 8 sec.
     */
    public final static AttributeKey<Integer> READ_TIME_OUT = AttributeKey.valueOf("READ_TIME_OUT");

    /*
     * If your application requires sampling times greater than the
     * set "watchdog" time, it is important that the PING option is activated,
     * this will prevent the TCP channel from being closed unnecessarily.
     * Default value: false
     */
    public final static AttributeKey<Boolean> IS_PING_ACTIVE = AttributeKey.valueOf("IS_PIN_ACTIVE");

    /*
     * Time value in seconds at which the execution of the PING will be scheduled.
     * Generally set by developer experience, but generally should be the same
     * as READ_TIME_OUT / 2.
     * Default value: -1
     */
    public final static AttributeKey<Integer> PING_TIME = AttributeKey.valueOf("PING_TIME");

    /*
     * Time for supervision of TCP channels. If the channel is not active,
     * a safe stop of the EventLoop must be performed, to ensure that
     * no additional tasks are created.
     * Default value: 4
     */
    public final static AttributeKey<Integer> RETRY_TIME = AttributeKey.valueOf("RETRY_TIME");

    ChannelHandlerContext embedCtx = null;
    protected Channel embededChannel = null;
    protected Channel tcpChannel = null;
    protected Channel primaryChannel = null;
    protected Channel secondaryChannel = null;

    /*
     * From S7ProtocolLogic
     * TODO: Evaluate if the "embed_ctx" is really required since we set
     * the Embedded channel when we created it.
     */
    @Override
    protected void encode(ChannelHandlerContext ctx, ByteBuf outBB, List<Object> list) {
        logger.debug("ENCODE: {}", outBB.toString());
        if ((embedCtx == null) && (ctx.channel() instanceof EmbeddedChannel)) {
            embedCtx = ctx;
        }
        if ((tcpChannel != null) && (embedCtx == ctx)) {
            tcpChannel.writeAndFlush(outBB.copy());
        } else {
            list.add(outBB.copy());
        }
    }

    /*
     * To S7ProtocolLogic
     * The information received here from the channel "tcp_channel" is sent to
     * the pipeline of the channel "embedded_channel"
     */
    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf inbb, List<Object> list) throws Exception {
        embedCtx.fireChannelRead(inbb.copy());
    }

    @Override
    public void channelRegistered(ChannelHandlerContext ctx) throws Exception {
        super.channelRegistered(ctx);
        logger.debug("channelRegistered: {}", ctx.name());
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        super.exceptionCaught(ctx, cause);
        logger.debug("exceptionCaught: {}", ctx.name());
    }

    @Override
    public void channelWritabilityChanged(ChannelHandlerContext ctx) throws Exception {
        super.channelWritabilityChanged(ctx);
        logger.debug("channelWritabilityChanged: {}", ctx.name());
    }

    /*
     * The events detected here flow from the S7ProtocolLogic object.
     * Upon receiving the "ConnectEvent" event, we must safely add the watchdog
     * to the pipeline of the "tcp_channel" connection.
     * The supervision time can be defined in the connection URL,
     * the default value being 8 secs, this value being defined experimentally.
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        super.userEventTriggered(ctx, evt);
        logger.info("{}  userEventTriggered: {} Event: {}", LocalTime.now(), ctx.name(), evt);
        if (evt instanceof ConnectEvent) {
            try {
                tcpChannel.pipeline().remove("watchdog");
            } catch (Exception ex) {
                logger.info(ex.toString());
            }
            try {
                tcpChannel.pipeline().addFirst("watchdog", new ReadTimeoutHandler(30));
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
            logger.info("DisconnectEvent");
        }

    }

    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) throws Exception {
        super.channelReadComplete(ctx);
        logger.debug("{} channelReadComplete: {}", LocalTime.now(), ctx.name());
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        super.channelInactive(ctx);
        logger.debug("channelInactive: {}", ctx.name());
    }

    @Override
    public void channelActive(ChannelHandlerContext ctx) throws Exception {
        super.channelActive(ctx);
        logger.debug("channelActive: {}", ctx.name());
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        super.channelUnregistered(ctx);
        logger.debug("{} channelUnregistered: {}", LocalTime.now(), ctx.name());
        String strCanal = (tcpChannel == primaryChannel) ? "PRIMARY" : "SECONDARY";
        logger.debug("Unregistered of channel: {}", strCanal);
        //TODO: If embedded channel is closed, we need close all channels
        if (ctx == embedCtx) return;

        if (tcpChannel == ctx.channel())
            embededChannel.attr(IS_CONNECTED).set(false);

        logger.info(embedCtx.executor().toString());

        if ((tcpChannel == primaryChannel) &&
            (primaryChannel == ctx.channel()))
            if ((!primaryChannel.isActive()) &&
                (secondaryChannel != null)) {
                if (secondaryChannel.isActive())
                    synchronized (tcpChannel) {
                        logger.info("Using secondary TCP channel.");
                        tcpChannel = secondaryChannel;
                        embededChannel.attr(IS_PRIMARY).set(false);
                        embededChannel.pipeline().fireUserEventTriggered(new ConnectEvent());
                    }
            }


        if ((tcpChannel == secondaryChannel) &&
            (secondaryChannel == ctx.channel()))
            if ((!secondaryChannel.isActive() &&
                (primaryChannel.isActive()))) {
                synchronized (tcpChannel) {
                    logger.info("Using primary TCP channel.");
                    tcpChannel = primaryChannel;
                    embededChannel.attr(IS_PRIMARY).set(true);
                    embededChannel.pipeline().fireUserEventTriggered(new ConnectEvent());
                }
            }
    }


    @Override
    public void setEmbeddedChannel(Channel embeded_channel) {
        this.embededChannel = embeded_channel;
        this.embededChannel.attr(IS_CONNECTED).set(false);
        this.embededChannel.attr(IS_PRIMARY).set(true);
        this.embededChannel.attr(READ_TIME_OUT).set(8);
        this.embededChannel.attr(IS_PING_ACTIVE).set(false);
        this.embededChannel.attr(PING_TIME).set(-1);
        this.embededChannel.attr(RETRY_TIME).set(8);
    }

    public void setPrimaryChannel(Channel primary_channel) {
        if ((this.primaryChannel == null) && (tcpChannel == null)) {
            if (primary_channel != null) {
                this.primaryChannel = primary_channel;
                tcpChannel = primary_channel;
                embededChannel.attr(IS_PRIMARY).set(true);
            }
        } else if ((!this.primaryChannel.isActive()) && (tcpChannel == secondaryChannel)) {
            this.primaryChannel = primary_channel;
        } else if ((!this.primaryChannel.isActive()) && (tcpChannel == this.primaryChannel)) {
            synchronized (tcpChannel) {
                tcpChannel.close();
                this.primaryChannel = primary_channel;
                tcpChannel = primary_channel;
                embededChannel.attr(IS_PRIMARY).set(true);
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
        } else if ((this.secondaryChannel == null) || (tcpChannel == primaryChannel)) {
            this.secondaryChannel = secondary_channel;
        } else if ((!this.secondaryChannel.isActive()) && (tcpChannel == primaryChannel)) {
            this.secondaryChannel = secondary_channel;
        } else if ((!this.secondaryChannel.isActive()) && (tcpChannel == this.secondaryChannel)) {
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


}
