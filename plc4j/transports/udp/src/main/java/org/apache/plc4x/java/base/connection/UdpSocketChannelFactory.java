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
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioDatagramChannel;
import org.apache.plc4x.java.spi.connection.NettyChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.SocketAddress;

public class UdpSocketChannelFactory extends NettyChannelFactory {

    private static final Logger logger = LoggerFactory.getLogger(UdpSocketChannelFactory.class);

    /**
     * @deprecated the next-gen drivers should use the {@link #UdpSocketChannelFactory(SocketAddress)}
     * constructor.
     */
    @Deprecated
    public UdpSocketChannelFactory(InetAddress address, int port) {
        super(new InetSocketAddress(address, port));
    }

    @Deprecated
    public UdpSocketChannelFactory(SocketAddress address) {
        super(address);
    }

    public UdpSocketChannelFactory() {
        // Default to use
    }

    @Override public Class<? extends Channel> getChannel() {
        return NioDatagramChannel.class;
    }

    @Override public void configureBootstrap(Bootstrap bootstrap) {
        // Do Nothing here
    }

}
