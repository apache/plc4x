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

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
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

    @BeforeEach
    void setUp() throws Exception {
        wrapper = new Plc4xNettyWrapper<>(channelPipeline, false, protocol, Date.class);

        ArgumentCaptor<ConversationContext<Date>> captor = ArgumentCaptor.forClass(ConversationContext.class);
        doNothing().when(protocol).onConnect(captor.capture());

        when(channelHandlerContext.channel()).thenReturn(channel);

        wrapper.userEventTriggered(channelHandlerContext, new ConnectEvent());
        conversationContext = captor.getValue();
    }

    @Test // see PLC4X-207 / PLC4X-257
    void conversationTimeoutTest() throws Exception {
        AtomicBoolean timeout = new AtomicBoolean(false);
        AtomicBoolean handled = new AtomicBoolean(false);
        AtomicBoolean error = new AtomicBoolean(false);

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
        assertFalse(timeout.get(), "timeout");
        assertFalse(handled.get(), "handled");
        assertFalse(error.get(), "error");

        wrapper.decode(channelHandlerContext, new Date(), new ArrayList<>());

        assertTrue(timeout.get());

        assertTrue(timeout.get(), "timeout");
        assertFalse(handled.get(), "handled");
        assertFalse(error.get(), "error");

    }
}