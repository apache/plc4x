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
package org.apache.plc4x.java.examples.dummydriver.connection;

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelInitializer;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.PlcWriter;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;
import org.apache.plc4x.java.base.messages.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class DummyConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    @SuppressWarnings("unused")
    private static final Logger logger = LoggerFactory.getLogger(DummyConnection.class);

    @SuppressWarnings("unused")
    private boolean connected;

    public DummyConnection(InetAddress host) {
        super(new TcpSocketChannelFactory(host, 42));
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
            }
        };
    }

    @Override
    public Optional<PlcReadRequest.Builder> readRequestBuilder() {
        // TODO: Implement this ...
        return Optional.empty();
    }

    @Override
    public Optional<PlcWriteRequest.Builder> writeRequestBuilder() {
        // TODO: Implement this ...
        return Optional.empty();
    }

    @Override
    public Optional<PlcSubscriptionRequest.Builder> subscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public Optional<PlcUnsubscriptionRequest.Builder> unsubscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<InternalPlcReadResponse> readFuture = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcReadRequest, InternalPlcReadResponse> container =
            new PlcRequestContainer<>((InternalPlcReadRequest) readRequest, readFuture);
        channel.writeAndFlush(container);
        return readFuture
            .thenApply(PlcReadResponse.class::cast);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<InternalPlcWriteResponse> writeFuture = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcWriteRequest, InternalPlcWriteResponse> container =
            new PlcRequestContainer<>((InternalPlcWriteRequest) writeRequest, writeFuture);
        channel.writeAndFlush(container);
        return writeFuture
            .thenApply(PlcWriteResponse.class::cast);
    }

}
