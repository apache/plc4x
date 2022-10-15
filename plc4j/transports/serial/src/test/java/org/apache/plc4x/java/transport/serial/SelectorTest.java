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
package org.apache.plc4x.java.transport.serial;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.read.ListAppender;
import io.netty.bootstrap.Bootstrap;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.channel.nio.NioEventLoop;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.Instant;
import java.util.Arrays;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

public class SelectorTest {

    /*
     * WARNING: This test is fragile as it tied to logging output behavior and format.
     * However, I could not find my reliable hook points, and it was better than not having testing at all
     */
    @Disabled("This test is best executed manually as necessary due to fragility and to avoid increasing build times")
    @Test
    public void idleDeviceTest() throws Exception {
        final var maxTestRunTime = Duration.ofSeconds(30);

        /*
         * Configure loggers so test's logic can determine if we have done enough loops
         * and whether netty issues warnings that should not happen. E.g.
         */
        final var listAppender = new ListAppender<ILoggingEvent>();
        {
            final var loopLogger = LoggerFactory.getLogger(NioEventLoop.class);
            final var selectorLogger = LoggerFactory.getLogger(SerialPollingSelector.class);

            listAppender.start();

            /*
             * Enable debug mode of the loggers we need and grab copy of the output for test logic
             */
            Arrays.asList(loopLogger, selectorLogger)
                .stream()
                .map(Logger.class::cast)
                .forEach(l -> {
                    l.setLevel(Level.DEBUG);
                    l.addAppender(listAppender);
                });
            ;
        }


        /*
         * Setup IO loop that will check for read events to come in
         */
        final EventLoopGroup evenLoopGroup;
        {
            final var bootstrap = new Bootstrap();
            final var commPortAddress = new SerialSocketAddress("test", new IdleSerialChannelHandler());
            final var channelFactory = new SerialChannelFactory(commPortAddress);
            evenLoopGroup = channelFactory.getEventLoopGroup();

            bootstrap
                .group(evenLoopGroup)
                .channel(SerialChannel.class)
                .handler(new SimpleChannelInboundHandler<ByteBuf>() {
                    @Override
                    protected void channelRead0(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf) throws Exception {
                        throw new UnsupportedOperationException();
                    }
                })
                .remoteAddress(commPortAddress);

            /*
             * Start the processing and wait for connection to be established
             */
            bootstrap.connect().sync();
        }

        /*
         * Wait until 2k or more select loops happen or maximum time passes
         */
        final var start = Instant.now();

        /*
         * With current netty select has to return more than 512 time early to trigger warning so entry number was
         * chosen with a little safety buffer.
         */
        while (listAppender.list.size() < 600
            && Duration.between(start, Instant.now()).compareTo(maxTestRunTime) <= 0)
        {
            Thread.sleep(Duration.ofSeconds(1).toMillis());
        }

        /*
         * Stop netty
         */
        evenLoopGroup.shutdownGracefully(0, 4, TimeUnit.SECONDS).sync();

        /*
         * Check results.
         * Line: ....SerialPollingSelector - returning from select with 0 events"
         * means one completed call to select ideally
         *
         * Line ....NioEventLoop - Selector.select() returned prematurely 512 times in a row; rebuilding Selector ...
         * means the netty warning got tripped
         */

        final var selectorWarnings =
            listAppender.list.stream()
                .filter(log -> log.getFormattedMessage().contains("Selector.select() returned prematurely"))
                .collect(Collectors.toList());

        Assertions.assertEquals(0, selectorWarnings.size(), "Detected selector warnings");

        /*
         * Second check in case above breaks, given fragility of depending on logs
         */
        final var selectorExecutions =
            listAppender.list.stream()
                .filter(log -> log.getFormattedMessage().contains("returning from select"))
                .collect(Collectors.toList());

        final int maxSelectorCalls = 1;

        Assertions.assertTrue(selectorExecutions.size() <= maxSelectorCalls,
            String.format("Selector called less than %d times", maxSelectorCalls));
    }

}
