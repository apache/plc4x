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
package org.apache.plc4x.java.can.adapter.conversation;

import io.netty.channel.Channel;
import java.time.Duration;
import java.util.function.Function;

import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.transport.can.CANTransport.FrameHandler;
import org.apache.plc4x.java.transport.can.FrameData;

public class ConversationContextWrapper<C, T> implements ConversationContext<T> {

    private final ConversationContext<C> delegate;
    private final Class<C> wireType;
    private final Function<C, FrameData> adapter;
    private final FrameHandler<C, T> frameHandler;

    private final PlcAuthentication authentication;

    public ConversationContextWrapper(ConversationContext<C> delegate, Class<C> wireType, Function<C, FrameData> adapter, FrameHandler<C, T> frameHandler, PlcAuthentication authentication) {
        this.delegate = delegate;
        this.wireType = wireType;
        this.adapter = adapter;
        this.frameHandler = frameHandler;
        this.authentication = authentication;
    }

    @Override
    public PlcAuthentication getAuthentication() {
        return authentication;
    }

    @Override
    public Channel getChannel() {
        return delegate.getChannel();
    }

    @Override
    public boolean isPassive() {
        return delegate.isPassive();
    }

    @Override
    public void sendToWire(T msg) {
        delegate.sendToWire(frameHandler.toCAN(msg));
    }

    @Override
    public void fireConnected() {
        delegate.fireConnected();
    }

    @Override
    public void fireDisconnected() {
        delegate.fireDisconnected();
    }

    @Override
    public void fireDiscovered(Configuration c) {
        delegate.fireDiscovered(c);
    }

    @Override
    public SendRequestContext<T> sendRequest(T packet) {
        return new SendRequestContextWrapper<>(delegate.sendRequest(frameHandler.toCAN(packet)), wireType, adapter, frameHandler);
    }

    @Override
    public ExpectRequestContext<T> expectRequest(Class<T> clazz, Duration timeout) {
        throw new UnsupportedOperationException();
    }

}
