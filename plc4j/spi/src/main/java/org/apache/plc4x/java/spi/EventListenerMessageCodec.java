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

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.api.listener.EventListener;
import org.apache.plc4x.java.api.listener.MessageExchangeListener;
import org.apache.plc4x.java.spi.generation.Message;

import java.util.List;

/**
 * Codec which propagate received or sent messages to connection event listeners.
 */
public class EventListenerMessageCodec extends MessageToMessageCodec<Message, Message> {

    private final List<EventListener> listeners;

    public EventListenerMessageCodec(List<EventListener> listeners) {
        this.listeners = listeners;
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, Message msg, List<Object> out) throws Exception {
        if (msg != null) {
            for (EventListener listener : listeners) {
                if (listener instanceof MessageExchangeListener) {
                    ((MessageExchangeListener) listener).sending(msg);
                }
            }
        }
        out.add(msg);
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, Message msg, List<Object> out) throws Exception {
        if (msg != null) {
            for (EventListener listener : listeners) {
                if (listener instanceof MessageExchangeListener) {
                    ((MessageExchangeListener) listener).received(msg);
                }
            }
        }
        out.add(msg);
    }

}
