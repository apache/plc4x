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
import io.netty.channel.EventLoopGroup;
import io.netty.channel.oio.OioEventLoopGroup;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.GenericFutureListener;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketAddress;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketChannel;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketChannelOption;
import org.apache.plc4x.java.utils.pcapsockets.netty.TcpIpPacketHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;

public class PcapChannelFactory implements ChannelFactory {

    private static final Logger logger = LoggerFactory.getLogger(PcapChannelFactory.class);

    private final File pcapFile;

    public PcapChannelFactory(File pcapFile) {
        this.pcapFile = pcapFile;
    }

    @Override
    public Channel createChannel(ChannelHandler channelHandler)
        throws PlcConnectionException {
        try {
            final EventLoopGroup workerGroup = new OioEventLoopGroup();

            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(workerGroup);
            bootstrap.channel(PcapSocketChannel.class);
            bootstrap.option(PcapSocketChannelOption.PACKET_HANDLER, new TcpIpPacketHandler());
            bootstrap.option(PcapSocketChannelOption.SPEED_FACTOR, 1.000f);
            bootstrap.handler(channelHandler);

            // Start the client.
            ChannelFuture f = bootstrap.connect(new PcapSocketAddress(pcapFile)).sync();
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

    public File getPcapFile() {
        return pcapFile;
    }

}
