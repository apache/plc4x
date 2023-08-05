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
package org.apache.plc4x.merlot.modbus.svr.impl;

import org.apache.plc4x.merlot.modbus.svr.api.ModbusServer;
import org.apache.plc4x.merlot.modbus.svr.core.ModbusServerADUDecoder;
import org.apache.plc4x.merlot.modbus.svr.core.ModbusServerADUEncoder;
import org.apache.plc4x.merlot.modbus.svr.core.ModbusServerADUHandler;
import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.LengthFieldBasedFrameDecoder;
import io.netty.handler.logging.LoggingHandler;
import io.netty.handler.traffic.GlobalChannelTrafficShapingHandler;
import java.net.SocketAddress;
import java.time.Duration;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import org.osgi.framework.BundleContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ModbusServerImpl implements ModbusServer {

    private final Logger LOGGER = LoggerFactory.getLogger(ModbusServer.class.getName());

    private BundleContext bc;
    private EventLoopGroup bossGroup;
    private EventLoopGroup workerGroup;
    private ServerBootstrap b;
    private GlobalChannelTrafficShapingHandler traffic;

    private int port = 502;
    private String host = "0.0.0.0";
    private SocketAddress[] saddress = null;
    private boolean started = false;
    private long checkInterval = 5000;

    private Date date = new Date();
    private Instant StartTime;
    private Instant EndTime;
    private Boolean Running = false;

    public void init() throws Exception {
        StartServer();
    }

    public void destroy() throws Exception {
        ShutDownServer();
    }

    public void start() {
        if (!started) {
            StartServer();
        };
    }

    public void stop() {
        if (started){
            ShutDownServer();
        }
    }

    public void setBundleContext(BundleContext bundleContext) {
        this.bc = bundleContext;
    }

    public void setPort(int port) {
        this.port = port;
    }

    public int getPort() {
        return this.port;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public String getHost() {
        return this.host;
    }

    @Override
    public boolean isStarted() {
        return started;
    }
        
    @Override
    public void setSocketAddress(SocketAddress[] saddress) {
        this.saddress = saddress;
    }

    @Override
    public SocketAddress[] getSocketAddress() {
        return this.saddress;
    }

    public void StartServer() {
        bossGroup = new NioEventLoopGroup(); // (1)
        workerGroup = new NioEventLoopGroup();
        StartTime = Instant.now();
        traffic = new GlobalChannelTrafficShapingHandler(workerGroup, checkInterval);
        try {
            b = new ServerBootstrap(); // (2)
            b.group(bossGroup, workerGroup)
                    .channel(NioServerSocketChannel.class) // (3)
                    .childHandler(new ChannelInitializer<SocketChannel>() { // (4)
                        @Override
                        public void initChannel(SocketChannel ch) throws Exception {
                            ch.pipeline().addLast(
                                new LengthFieldBasedFrameDecoder(1024, 4, 2),
                                new LoggingHandler(),
                                traffic,                                
                                new ModbusServerADUDecoder(),
                                new ModbusServerADUEncoder(),
                                new ModbusServerADUHandler(bc));
                        }
                    })
                    .option(ChannelOption.SO_BACKLOG, 128) // (5)
                    .childOption(ChannelOption.SO_KEEPALIVE, true); // (6)

            // Bind and start to accept incoming connections.
            List<ChannelFuture> futures = new ArrayList<>();
            if (saddress == null) {
                b.bind(host, port).sync(); // (7)
            } else {
                for (SocketAddress sa:saddress){
                    futures.add(b.bind(sa)); // (8)
                }
                for (ChannelFuture f: futures) {
                    f.sync(); // (9)
                }
            }
            Running = true;
            // Wait until the server socket is closed.
            // In this example, this does not happen, but you can do that to gracefully
            // shut down your server.

            //f.channel().closeFuture().sync();
        } catch (Exception ex) {
            LOGGER.info("Failed to create the MODBUS server: " + ex.getMessage());
            ex.printStackTrace();
        } finally {
            ;
            // workerGroup.shutdownGracefully();
            //bossGroup.shutdownGracefully();	
        }
    }

    public void ShutDownServer() {
        LOGGER.info("Shutdown MODBUS server.");
        workerGroup.shutdownGracefully();
        bossGroup.shutdownGracefully();
        Running = false;
    }

    public String getDate() {
        return this.date.toString();
    }

    public String getElapseTime() {
        if (Running) {
            EndTime = Instant.now();
        } else {
            EndTime = StartTime;
        };
        return Duration.between(StartTime, EndTime).toString();
    }

    @Override
    public long getCheckInterval() {
        return traffic.trafficCounter().checkInterval();
    }

    @Override
    public void setCheckInterval(long newCheckInterval) {
        traffic.trafficCounter().configure(newCheckInterval);
    }

    @Override
    public long getLastReadThroughput() {
        return traffic.trafficCounter().lastReadThroughput();
    }

    @Override
    public long getLastWriteThroughput() {
        return traffic.trafficCounter().lastWriteThroughput();
    }

    @Override
    public long getCumulativeReadBytes() {
        return traffic.trafficCounter().cumulativeReadBytes();
    }

    @Override
    public long getCumulativeWrittenBytes() {
        return traffic.trafficCounter().cumulativeWrittenBytes();
    }

    @Override
    public long getCurrentReadBytes() {
        return traffic.trafficCounter().currentReadBytes();
    }

    @Override
    public long getCurrentWrittenBytes() {
        return traffic.trafficCounter().currentWrittenBytes();
    }

    
    
}
