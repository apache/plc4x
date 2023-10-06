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

package org.apache.plc4x.java.spi.internal;

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandlerContext;
import java.time.Duration;
import java.util.function.Consumer;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.spi.events.DisconnectedEvent;
import org.apache.plc4x.java.spi.events.DiscoveredEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class DefaultConversationContext<T1> implements ConversationContext<T1> {
    private final Logger logger = LoggerFactory.getLogger(DefaultConversationContext.class);

    private final Consumer<HandlerRegistration> handlerRegistrar;

    private final ChannelHandlerContext channelHandlerContext;
    private final PlcAuthentication authentication;
    private final boolean passive;

    public DefaultConversationContext(Consumer<HandlerRegistration> handlerRegistrar,
                                      ChannelHandlerContext channelHandlerContext,
                                      PlcAuthentication authentication,
                                      boolean passive) {
        this.handlerRegistrar = handlerRegistrar;
        this.channelHandlerContext = channelHandlerContext;
        this.authentication = authentication;
        this.passive = passive;
    }
    @Override
    public Channel getChannel() {
        return channelHandlerContext.channel();
    }

    public PlcAuthentication getAuthentication() {
        return authentication;
    }

    @Override
    public boolean isPassive() {
        return passive;
    }

    @Override
    public void sendToWire(T1 msg) {
        logger.trace("Sending to wire {}", msg);
        channelHandlerContext.channel().writeAndFlush(msg);
    }

    @Override
    public void fireConnected() {
        logger.trace("Firing Connected!");
        channelHandlerContext.pipeline().fireUserEventTriggered(new ConnectedEvent());
    }

    @Override
    public void fireDisconnected() {
        logger.trace("Firing Disconnected!");
        channelHandlerContext.pipeline().fireUserEventTriggered(new DisconnectedEvent());
    }

    @Override
    public void fireDiscovered(Configuration c) {
        logger.trace("Firing Discovered!");
        channelHandlerContext.pipeline().fireUserEventTriggered(new DiscoveredEvent(c));
    }

    @Override
    public SendRequestContext<T1> sendRequest(T1 packet) {
        return new DefaultSendRequestContext<>(handler -> {
            logger.trace("Adding Response Handler ...");
            handlerRegistrar.accept(handler);
        }, packet, this);
    }

    @Override
    public ExpectRequestContext<T1> expectRequest(Class<T1> clazz, Duration timeout) {
        return new DefaultExpectRequestContext<>(handler -> {
            logger.trace("Adding Request Handler ...");
            handlerRegistrar.accept(handler);
        }, clazz, timeout, this);
    }
}