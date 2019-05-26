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

package org.apache.plc4x.java.base.connection;

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.embedded.EmbeddedChannel;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

import java.util.concurrent.CompletableFuture;


public class NettyPlcConnectionTest implements WithAssertions {

    private final ChannelFactory channelFactory = new ChannelFactory() {
        @Override
        public Channel createChannel(ChannelHandler channelHandler) throws PlcConnectionException {
            return new EmbeddedChannel();
        }

        @Override
        public void ping() {
            // Ignore ...
        }

    };

    NettyPlcConnection SUT = new NettyPlcConnection(channelFactory, true) {
        @Override
        protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
            sessionSetupCompleteFuture.complete(null);

            return new ChannelHandler() {
                @Override
                public void handlerAdded(ChannelHandlerContext ctx) {
                }

                @Override
                public void handlerRemoved(ChannelHandlerContext ctx) {
                }

                @Override
                public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
                }
            };
        }

        @Override
        protected void sendChannelCreatedEvent() {
            channel.pipeline().fireUserEventTriggered(new ConnectEvent());
        }
    };

    @Test
    public void connect() throws Exception {
        SUT.connect();
        Channel channel = SUT.getChannel();
        assertThat(channel).isNotNull();
    }

    @Test
    public void close() throws Exception {
        SUT.close();
    }

    @Test
    public void isConnected() {
        SUT.isConnected();
    }

}
