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

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.*;
import org.apache.plc4x.java.api.connection.AbstractPlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.examples.dummydriver.model.DummyAddress;
import org.apache.plc4x.java.examples.dummydriver.netty.DummyProtocol;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketAddress;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

public class DummyConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    private static final Logger logger = LoggerFactory.getLogger(DummyConnection.class);

    private final String hostName;

    private EventLoopGroup workerGroup;
    private Channel channel;

    public DummyConnection(String hostName) {
        this.hostName = hostName;
    }

    public String getHostName() {
        return hostName;
    }

    @Override
    public void connect() throws PlcConnectionException {
        workerGroup = new DefaultEventLoopGroup() {
        };

        try {
            // As we don't just want to wait till the connection is established,
            // define a future we can use to signal back that the s7 session is
            // finished initializing.
            CompletableFuture<Void> sessionSetupCompleteFuture = new CompletableFuture<>();

            RawSocketAddress serverSocketAddress = new RawSocketAddress(hostName);

            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(workerGroup);
            bootstrap.channel(RawSocketChannel.class);
            bootstrap.handler(new ChannelInitializer() {
                @Override
                protected void initChannel(Channel channel) throws Exception {
                    ChannelPipeline pipeline = channel.pipeline();
                    pipeline.addLast(new DummyProtocol());
                }
            });
            // Start the client.
            ChannelFuture f = bootstrap.connect(serverSocketAddress).sync();
            f.awaitUninterruptibly();
            // Wait till the session is finished initializing.
            channel = f.channel();

            sessionSetupCompleteFuture.get();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException(e);
        }
        catch (ExecutionException e) {
            throw new PlcConnectionException(e);
        }
    }

    @Override
    public void close() throws Exception {
        if((channel != null) && channel.isOpen()) {
            channel.closeFuture().await();
        }

        if (workerGroup != null) {
            workerGroup.shutdownGracefully();
        }
    }

    @Override
    public Address parseAddress(String addressString) {
        return new DummyAddress(Integer.parseInt(addressString));
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> readFuture = new CompletableFuture<>();
        PlcRequestContainer<PlcReadRequest, PlcReadResponse> container =
            new PlcRequestContainer<>(readRequest, readFuture);
        channel.writeAndFlush(container);
        return readFuture;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> writeFuture = new CompletableFuture<>();
        PlcRequestContainer<PlcWriteRequest, PlcWriteResponse> container =
            new PlcRequestContainer<>(writeRequest, writeFuture);
        channel.writeAndFlush(container);
        return writeFuture;
    }

}
