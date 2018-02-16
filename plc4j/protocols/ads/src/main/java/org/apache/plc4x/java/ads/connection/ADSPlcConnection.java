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

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.model.ADSAddress;
import org.apache.plc4x.java.ads.netty.ADSProtocol;
import org.apache.plc4x.java.ads.netty.Plc4XADSProtocol;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;

import java.net.InetAddress;
import java.util.concurrent.CompletableFuture;

public class ADSPlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    public static final int TCP_PORT = 48898;

    private final AMSNetId targetAmsNetId;

    private final AMSPort targetAmsPort;

    private final AMSNetId sourceAmsNetId;

    private final AMSPort sourceAmsPort;

    public ADSPlcConnection(InetAddress address, AMSNetId targetAmsNetId, AMSPort targetAmsPort) {
        this(address, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    public ADSPlcConnection(InetAddress address, Integer port, AMSNetId targetAmsNetId, AMSPort targetAmsPort) {
        this(address, port, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }


    public ADSPlcConnection(InetAddress address, AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        this(address, null, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    public ADSPlcConnection(InetAddress address, Integer port, AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        super(new TcpSocketChannelFactory(address, port));
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
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
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the ads protocol.
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new Plc4XADSProtocol(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort));
                pipeline.addLast(new ADSProtocol());
            }
        };
    }

    @Override
    public Address parseAddress(String addressString) throws PlcException {
        return ADSAddress.of(addressString);
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

    private static AMSNetId generateAMSNetId() {
        return AMSNetId.of("0.0.0.0.0.0");
    }

    private static AMSPort generateAMSPort() {
        return AMSPort.of(TCP_PORT);
    }

}
