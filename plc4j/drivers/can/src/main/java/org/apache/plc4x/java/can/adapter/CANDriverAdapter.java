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
package org.apache.plc4x.java.can.adapter;

import java.util.concurrent.CompletableFuture;
import java.util.function.Function;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.can.adapter.conversation.ConversationContextWrapper;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.transport.can.CANTransport.FrameHandler;
import org.apache.plc4x.java.transport.can.FrameData;

public class CANDriverAdapter<C, T> extends Plc4xProtocolBase<C> {

    private final Plc4xCANProtocolBase<T> delegate;
    private final Class<C> wireType;
    private final Function<C, FrameData> adapter;
    private final FrameHandler<C, T> frameHandler;

    public CANDriverAdapter(Plc4xCANProtocolBase<T> delegate, Class<C> wireType, Function<C, FrameData> adapter, FrameHandler<C, T> frameHandler) {
        this.delegate = delegate;
        this.wireType = wireType;
        this.adapter = adapter;
        this.frameHandler = frameHandler;
    }

    @Override
    public void setDriverContext(DriverContext driverContext) {
        delegate.setDriverContext(driverContext);
    }

    @Override
    public DriverContext getDriverContext() {
        return delegate.getDriverContext();
    }

    @Override
    public void setContext(ConversationContext<C> context) {
        delegate.setContext(new ConversationContextWrapper<>(context, wireType, adapter, frameHandler));
    }

    @Override
    public void onConnect(ConversationContext<C> context) {
        delegate.onConnect(new ConversationContextWrapper<>(context, wireType, adapter, frameHandler));
    }

    @Override
    public void onDisconnect(ConversationContext<C> context) {
        delegate.onDisconnect(new ConversationContextWrapper<>(context, wireType, adapter, frameHandler));
    }

    @Override
    protected void decode(ConversationContext<C> context, C msg) throws Exception {
        delegate.decode(new ConversationContextWrapper<>(context, wireType, adapter, frameHandler), frameHandler.fromCAN(adapter.apply(msg)));
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        return delegate.read(readRequest);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        return delegate.write(writeRequest);
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        return delegate.subscribe(subscriptionRequest);
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return delegate.unsubscribe(unsubscriptionRequest);
    }

    @Override
    public void close(ConversationContext<C> context) {
        delegate.close(new ConversationContextWrapper<>(context, wireType, adapter, frameHandler));
    }

}
