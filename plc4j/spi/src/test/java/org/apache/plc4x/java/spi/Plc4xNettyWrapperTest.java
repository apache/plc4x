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
package org.apache.plc4x.java.spi;

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPipeline;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Date;
import java.util.concurrent.atomic.AtomicBoolean;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class Plc4xNettyWrapperTest {

    @Mock
    Plc4xProtocolBase<Date> protocol;
    @Mock
    ChannelPipeline channelPipeline;
    @Mock
    ChannelHandlerContext channelHandlerContext;
    @Mock
    Channel channel;

    Plc4xNettyWrapper<Date> wrapper;

    ConversationContext<Date> conversationContext;
    private AtomicBoolean timeout;
    private AtomicBoolean handled;
    private AtomicBoolean error;

    @BeforeEach
    void setUp() throws Exception {
        wrapper = new Plc4xNettyWrapper<>(channelPipeline, false, protocol, null, Date.class);

        ArgumentCaptor<ConversationContext<Date>> captor = ArgumentCaptor.forClass(ConversationContext.class);
        doNothing().when(protocol).onConnect(captor.capture());

        when(channelHandlerContext.channel()).thenReturn(channel);

        wrapper.userEventTriggered(channelHandlerContext, new ConnectEvent());
        conversationContext = captor.getValue();

        timeout = new AtomicBoolean(false);
        handled = new AtomicBoolean(false);
        error = new AtomicBoolean(false);
    }

    @Test // see PLC4X-207 / PLC4X-257
    void conversationTimeoutTest() throws Exception {
        ConversationContext.ContextHandler handler = conversationContext.sendRequest(new Date())
            .expectResponse(Date.class, Duration.ofMillis(500))
            .onTimeout(e -> {
                timeout.set(true);
            })
            .onError((value, throwable) -> {
                error.set(true);
            })
            .handle((answer) -> {
                handled.set(true);
            });

        Thread.sleep(750);

        verify(false, false, false);
        wrapper.decode(channelHandlerContext, new Date(), new ArrayList<>());
        verify(true, false, false);

    }

    @Test // see PLC4X-207 / PLC4X-257
    void conversationWithNoTimeoutTest() throws Exception {
        ConversationContext.ContextHandler handler = conversationContext.sendRequest(new Date())
            .expectResponse(Date.class, Duration.ofMillis(500))
            .onTimeout(e -> {
                timeout.set(true);
            })
            .onError((value, throwable) -> {
                error.set(true);
            })
            .handle((answer) -> {
                handled.set(true);
            });

        verify(false, false, false);
        wrapper.decode(channelHandlerContext, new Date(), new ArrayList<>());
        verify(false, false, true);
    }

    void verify(boolean isTimeout, boolean isError, boolean isHandled) {
        assertThat(timeout.get()).describedAs("Expected timeout state %b", isTimeout)
            .isEqualTo(isTimeout);
        assertThat(error.get()).describedAs("Expected error state %b", isError)
            .isEqualTo(isError);
        assertThat(handled.get()).describedAs("Expected handled state %b", isHandled)
            .isEqualTo(isHandled);
    }
}