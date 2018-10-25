/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.protocol.util;

import io.netty.channel.ChannelDuplexHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPromise;
import io.netty.util.concurrent.ScheduledFuture;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.SocketAddress;
import java.util.ArrayDeque;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * Inspired by {@code ChannelTrafficShapingHandler} this limiter ensures only one message is sent at a time.
 */
public class SingleMessageRateLimiter extends ChannelDuplexHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(SingleMessageRateLimiter.class);

    private final ArrayDeque<ToSend> messagesQueue = new ArrayDeque<>();

    private AtomicBoolean messageOnTheWay = new AtomicBoolean(false);

    private ScheduledFuture<?> sender;

    @Override
    public void bind(ChannelHandlerContext ctx, SocketAddress localAddress, ChannelPromise promise) throws Exception {
        LOGGER.debug("bind({}, {}, {})", ctx, localAddress, promise);
        super.bind(ctx, localAddress, promise);
    }

    @Override
    public void deregister(ChannelHandlerContext ctx, ChannelPromise promise) throws Exception {
        LOGGER.debug("deregister({}, {})", ctx, promise);
        super.deregister(ctx, promise);
    }

    @Override
    public void connect(ChannelHandlerContext ctx, SocketAddress remoteAddress, SocketAddress localAddress, ChannelPromise promise) throws Exception {
        LOGGER.debug("connect({}, {}, {}, {})", ctx, remoteAddress, localAddress, promise);
        sender = ctx.executor().scheduleAtFixedRate(() -> {
            LOGGER.trace("Woke up and doing work messageOnTheWay:{}, messageQueue:{}", messageOnTheWay, messagesQueue);
            if (!messagesQueue.isEmpty() && messageOnTheWay.compareAndSet(false, true)) {
                ToSend pop = messagesQueue.pop();
                LOGGER.debug("Sending {}", pop);
                pop.channelHandlerContext.writeAndFlush(pop.objectToSend, pop.promise);
                LOGGER.debug("Send {}", pop);
            }
        }, 100, 10, TimeUnit.MILLISECONDS);
        super.connect(ctx, remoteAddress, localAddress, promise);
    }

    @Override
    public void disconnect(ChannelHandlerContext ctx, ChannelPromise promise) throws Exception {
        LOGGER.debug("disconnect({}, {}, {}, {})", ctx, promise);
        sender.cancel(false);
        super.disconnect(ctx, promise);
    }

    @Override
    public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) {
        LOGGER.trace("(<--OUT): {}, {}, {}", ctx, msg, promise);
        if (messageOnTheWay.compareAndSet(false, true)) {
            ctx.write(msg, promise);
        } else {
            messagesQueue.add(new ToSend(ctx, msg, promise));
        }
    }

    @Override
    public void read(ChannelHandlerContext ctx) throws Exception {
        LOGGER.trace("(-->In): {}", ctx);
        messageOnTheWay.set(false);
        super.read(ctx);
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        LOGGER.trace("(-->ERR): {}", ctx, cause);
        messageOnTheWay.set(false);
        super.exceptionCaught(ctx, cause);
    }

    private static final class ToSend {
        final ChannelHandlerContext channelHandlerContext;
        final Object objectToSend;
        final ChannelPromise promise;

        private ToSend(ChannelHandlerContext channelHandlerContext, Object objectToSend, ChannelPromise promise) {
            this.channelHandlerContext = channelHandlerContext;
            this.objectToSend = objectToSend;
            this.promise = promise;
        }

        @Override
        public String toString() {
            return "ToSend{" +
                "channelHandlerContext=" + channelHandlerContext +
                ", objectToSend=" + objectToSend +
                ", promise=" + promise +
                '}';
        }
    }
}
