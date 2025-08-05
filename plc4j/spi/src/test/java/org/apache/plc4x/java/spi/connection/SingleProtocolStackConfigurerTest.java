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

package org.apache.plc4x.java.spi.connection;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.embedded.EmbeddedChannel;
import java.util.function.ToIntFunction;
import org.apache.plc4x.java.api.EventPlcConnection;
import org.apache.plc4x.java.api.listener.ConnectionStateListener;
import org.apache.plc4x.java.api.listener.MessageExchangeListener;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
class SingleProtocolStackConfigurerTest {

    @Mock
    private MessageExchangeListener messageListener;
    @Mock
    private ConnectionStateListener connectionListener;

    @Test
    void testConnectionStateListener() throws Exception {
        TestChannelFactory channelFactory = new TestChannelFactory();

        SingleProtocolStackConfigurer<TestMessage> stackConfigurer = SingleProtocolStackConfigurer.builder(TestMessage.class, TestMessage::staticParse)
            .withProtocol(TestProtocol.class)
            .build();

        EventPlcConnection connection = new PlcConnectionFactory().create(channelFactory, stackConfigurer);
        connection.addEventListener(connectionListener);

        connection.connect();
        verify(connectionListener).connected();

        connection.close();
        verify(connectionListener).disconnected();
    }

    @Test
    void testConnectionStateListenerAppended() throws Exception {
        TestChannelFactory channelFactory = new TestChannelFactory();

        SingleProtocolStackConfigurer<TestMessage> stackConfigurer = SingleProtocolStackConfigurer.builder(TestMessage.class, TestMessage::staticParse)
            .withProtocol(TestProtocol.class)
            .build();

        EventPlcConnection connection = new PlcConnectionFactory().create(channelFactory, stackConfigurer);
        connection.addEventListener(connectionListener);

        connection.connect();
        verify(connectionListener).connected();

        // append listener after connection been made
        ConnectionStateListener dynamicListener = mock(ConnectionStateListener.class);
        connection.addEventListener(dynamicListener);

        connection.close();
        verify(connectionListener).disconnected();
        verify(dynamicListener).disconnected();
    }

    @Test
    void testMessageExchangeListener() throws Exception {
        TestChannelFactory channelFactory = new TestChannelFactory();

        SingleProtocolStackConfigurer<TestMessage> stackConfigurer = SingleProtocolStackConfigurer.builder(TestMessage.class, TestMessage::staticParse)
            .withProtocol(TestProtocol.class)
            .withPacketSizeEstimator(Estimator.class)
            .build();

        EventPlcConnection connection = new PlcConnectionFactory().doNotAwaitForDisconnect()
            .create(channelFactory, stackConfigurer);
        connection.addEventListener(messageListener);

        connection.connect();

        ByteBuf buffer = Unpooled.wrappedBuffer(new byte[] {0x00});
        EmbeddedChannel channel = channelFactory.getChannel();
        // send dummy message with 0x00
        channel.writeInbound(buffer.retain());
        assertTrue(channel.finish());

        // await completion of handshake
        connection.close();

        verify(messageListener).received(eq(new TestMessage(0)));
        verify(messageListener).sending(eq(new TestMessage(1)));
    }

    static class Estimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf value) {
            return 1;
        }
    }
}