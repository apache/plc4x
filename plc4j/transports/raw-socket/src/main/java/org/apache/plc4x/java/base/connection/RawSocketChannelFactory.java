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
import io.netty.channel.*;
import io.netty.channel.oio.OioEventLoopGroup;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.GenericFutureListener;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.utils.rawsockets.netty.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;

public class RawSocketChannelFactory implements ChannelFactory {

    private static final Logger logger = LoggerFactory.getLogger(RawSocketChannelFactory.class);

    private static final int PING_TIMEOUT_MS = 1_000;

    private final String deviceName;
    private final InetAddress address;
    private final int port;
    private final int protocolId;
    private final PacketHandler packetHandler;

    public RawSocketChannelFactory(String deviceName, InetAddress address, int port, int protocolId, PacketHandler packetHandler) {
        this.deviceName = deviceName;
        this.address = address;
        this.port = port;
        this.protocolId = protocolId;
        this.packetHandler = packetHandler;
    }

    @Override
    public Channel createChannel(ChannelHandler channelHandler)
        throws PlcConnectionException {
        try {
            final EventLoopGroup workerGroup = new OioEventLoopGroup();

            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(workerGroup);
            bootstrap.channel(RawSocketChannel.class);
            bootstrap.option(RawSocketChannelOption.PACKET_HANDLER, packetHandler);
            // TODO we should use an explicit (configurable?) timeout here
            // bootstrap.option(ChannelOption.CONNECT_TIMEOUT_MILLIS, 1000);
            bootstrap.handler(channelHandler);

            // Start the client.
            ChannelFuture f = bootstrap.connect(new RawSocketIpAddress(deviceName, protocolId, address, port)).sync();
            f.addListener(new GenericFutureListener<Future<? super Void>>() {
                @Override public void operationComplete(Future<? super Void> future) throws Exception {
                    if (!future.isSuccess()) {
                        logger.info("Unable to connect, shutting down worker thread.");
                        workerGroup.shutdownGracefully();
                    }
                }
            });
            // Wait for sync
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
    public void ping() throws PlcException {
        // Raw-Sockets are absolutely passive, so there is nothing to do for a ping.
    }

    public String getDeviceName() {
        return deviceName;
    }

    public InetAddress getAddress() {
        return address;
    }

    public int getPort() {
        return port;
    }

    public int getProtocolId() {
        return protocolId;
    }

}
