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
package org.apache.plc4x.java.ads.connection;

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.*;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioSocketChannel;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.model.ADSAddress;
import org.apache.plc4x.java.ads.netty.ADSProtocol;
import org.apache.plc4x.java.ads.netty.Plc4XADSProtocol;
import org.apache.plc4x.java.api.connection.AbstractPlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.Address;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.concurrent.CompletableFuture;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ADSPlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    public static final int TCP_PORT = 48898;

    private final String hostName;

    private final Integer suppliedPort;

    private final AMSNetId targetAmsNetId;

    private final AMSPort targetAmsPort;

    private final AMSNetId sourceAmsNetId;

    private final AMSPort sourceAmsPort;

    private EventLoopGroup workerGroup;
    private Channel channel;

    public ADSPlcConnection(String hostName, AMSNetId targetAmsNetId, AMSPort targetAmsPort) {
        this(hostName, targetAmsNetId, targetAmsPort, null, null);
    }

    public ADSPlcConnection(String hostName, Integer port, AMSNetId targetAmsNetId, AMSPort targetAmsPort) {
        this(hostName, port, targetAmsNetId, targetAmsPort, null, null);
    }


    public ADSPlcConnection(String hostName, AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        this(hostName, null, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    public ADSPlcConnection(String hostName, Integer port, AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        this.hostName = hostName;
        this.suppliedPort = port;
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
    }

    public String getHostName() {
        return hostName;
    }

    public AMSNetId getTargetAmsNetId() {
        return targetAmsNetId;
    }

    public AMSPort getTargetAmsPort() {
        return targetAmsPort;
    }

    public AMSNetId getSourceAmsNetId() {
        return sourceAmsNetId;
    }

    public AMSPort getSourceAmsPort() {
        return sourceAmsPort;
    }

    @Override
    public void connect() throws PlcConnectionException {
        workerGroup = new NioEventLoopGroup();

        try {
            InetAddress serverInetAddress = InetAddress.getByName(hostName);

            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(workerGroup);
            bootstrap.channel(NioSocketChannel.class);
            bootstrap.option(ChannelOption.SO_KEEPALIVE, true);
            bootstrap.option(ChannelOption.TCP_NODELAY, true);
            bootstrap.handler(new ChannelInitializer() {
                @Override
                protected void initChannel(Channel channel) throws Exception {
                    // Build the protocol stack for communicating with the ads protocol.
                    ChannelPipeline pipeline = channel.pipeline();
                    pipeline.addLast(new Plc4XADSProtocol());
                    pipeline.addLast(new ADSProtocol());
                }
            });
            // Start the client.
            ChannelFuture f = bootstrap.connect(serverInetAddress, suppliedPort != null ? suppliedPort : TCP_PORT).sync();
            f.awaitUninterruptibly();
            // Wait till the session is finished initializing.
            channel = f.channel();
        } catch (UnknownHostException e) {
            throw new PlcConnectionException("Unknown Host " + hostName, e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException(e);
        }
    }

    @Override
    public void close() throws Exception {
        workerGroup.shutdownGracefully();
    }

    @Override
    public Address parseAddress(String addressString) throws PlcException {
        Matcher matcher = Pattern.compile("(?<targetAmsNetId>" + AMSNetId.AMS_NET_ID_PATTERN + "):(?<targetAmsPort>" + AMSPort.AMS_PORT_PATTERN + ")").matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Address string doesn't match the format '{targetAmsNetId}:{targetAmsPort}'");
        }
        AMSNetId targetAmsNetId = AMSNetId.of(matcher.group("targetAmsNetId"));
        AMSPort targetAmsPort = AMSPort.of(matcher.group("targetAmsPort"));
        return new ADSAddress(targetAmsNetId, targetAmsPort);
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
