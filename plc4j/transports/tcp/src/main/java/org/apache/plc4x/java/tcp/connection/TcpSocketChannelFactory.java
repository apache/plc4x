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
package org.apache.plc4x.java.tcp.connection;

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.GenericFutureListener;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.NettyChannelFactory;
import org.apache.plc4x.java.spi.connection.NettyPlcConnection;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import sun.reflect.generics.reflectiveObjects.NotImplementedException;

import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;

public class TcpSocketChannelFactory extends NettyChannelFactory {

    private static final Logger logger = LoggerFactory.getLogger(TcpSocketChannelFactory.class);

    public static final String SO_KEEPALIVE = "SO_KEEPALIVE";
    public static final String CONNECT_TIMEOUT_MILLIS = "CONNECT_TIMEOUT_MILLIS";
    private static final String TCP_NODELAY = "TCP_NODELAY";

    public TcpSocketChannelFactory() {
        // Default Constructor
    }

    /**
     * Only there for retrofit
     */
    @Deprecated
    public TcpSocketChannelFactory(SocketAddress address) {
        super(address);
    }

    /**
     * Only there for retrofit
     */
    @Deprecated
    public TcpSocketChannelFactory(InetAddress address, int port) {
        this(new InetSocketAddress(address, port));
    }


    @Override public Class<? extends Channel> getChannel() {
        return NioSocketChannel.class;
    }

    @Override public void configureBootstrap(Bootstrap bootstrap) {
        boolean keepalive = Boolean.parseBoolean(getPropertyOrDefault(SO_KEEPALIVE, "true"));
        boolean nodelay = Boolean.parseBoolean(getPropertyOrDefault(TCP_NODELAY, "true"));
        int connectTimeout = Integer.parseInt(getPropertyOrDefault(CONNECT_TIMEOUT_MILLIS, "1000"));
        logger.info("Configuring Bootstrap with properties\n\t{} {}\n\t{} {}\n\t{} {}", SO_KEEPALIVE, keepalive, TCP_NODELAY, nodelay, CONNECT_TIMEOUT_MILLIS, connectTimeout);
        bootstrap.option(ChannelOption.SO_KEEPALIVE, keepalive);
        bootstrap.option(ChannelOption.TCP_NODELAY, nodelay);
        bootstrap.option(ChannelOption.CONNECT_TIMEOUT_MILLIS, connectTimeout);
    }

}
