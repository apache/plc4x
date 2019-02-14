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

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelOption;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioSocketChannel;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

import java.net.InetAddress;

public class TcpSocketChannelFactory implements ChannelFactory {

    private final InetAddress address;
    private final int port;

    public TcpSocketChannelFactory(InetAddress address, int port) {
        this.address = address;
        this.port = port;
    }

    @Override
    public Channel createChannel(ChannelHandler channelHandler)
        throws PlcConnectionException {
        try {
            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(new NioEventLoopGroup());
            bootstrap.channel(NioSocketChannel.class);
            bootstrap.option(ChannelOption.SO_KEEPALIVE, true);
            bootstrap.option(ChannelOption.TCP_NODELAY, true);
            bootstrap.handler(channelHandler);
            // Start the client.
            ChannelFuture f = bootstrap.connect(address, port).sync();
            f.awaitUninterruptibly();
            // Wait till the session is finished initializing.
            return f.channel();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException("Error creating channel.", e);
        } catch (Exception e) {
            throw new PlcConnectionException("Error creating channel.", e);
        }
    }

    public InetAddress getAddress() {
        return address;
    }

    public int getPort() {
        return port;
    }
}
