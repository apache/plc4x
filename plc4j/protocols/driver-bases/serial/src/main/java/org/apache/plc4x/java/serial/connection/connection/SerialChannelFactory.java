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
package org.apache.plc4x.java.serial.connection.connection;

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.GenericFutureListener;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.base.connection.ChannelFactory;

import java.net.SocketAddress;
import java.util.concurrent.Executor;

public class SerialChannelFactory implements ChannelFactory {

    private final String serialPort;

    public SerialChannelFactory(String serialPort) {
        this.serialPort = serialPort;
    }

    @Override
    public Channel createChannel(ChannelHandler channelHandler)
        throws PlcConnectionException {
        SocketAddress address = new SerialSocketAddress(serialPort);

        try {
            Bootstrap bootstrap = new Bootstrap();
            final NioEventLoopGroup eventLoop = new NioEventLoopGroup(0, (Executor) null, new SerialSelectorProvider());
            bootstrap.group(eventLoop);
            bootstrap.channel(SerialChannel.class);
            bootstrap.handler(channelHandler);
            // Start the client.
            ChannelFuture f = bootstrap.connect(address);
            f.addListener(new GenericFutureListener<Future<? super Void>>() {
                    @Override public void operationComplete(Future<? super Void> future) throws Exception {
                        if (future.isSuccess()) {
                            System.out.println("Connection sucesfull!");
                        } else {
                            System.out.println("Connection not sucessfull: " + future.cause().getMessage());
                            eventLoop.shutdownGracefully();
                        }
                    }
                });
            f.sync();
            f.awaitUninterruptibly();
            // Wait till the session is finished initializing.
            return f.channel();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException("Error creating channel.", e);
        }
    }

    @Override
    public void ping() {
        // TODO: Do some sort of check as soon as we know how ...
    }

    public String getSerialPort() {
        return serialPort;
    }
}
