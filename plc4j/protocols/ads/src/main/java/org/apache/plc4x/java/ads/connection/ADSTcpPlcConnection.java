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
import org.apache.plc4x.java.ads.protocol.ADS2TcpProtocol;
import org.apache.plc4x.java.ads.protocol.Plc4X2ADSProtocol;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;

import java.net.InetAddress;
import java.util.concurrent.CompletableFuture;

public class ADSTcpPlcConnection extends ADSAbstractPlcConnection {

    private static final int TCP_PORT = 48898;

    public ADSTcpPlcConnection(InetAddress address, AMSNetId targetAmsNetId, AMSPort targetAmsPort) {
        this(address, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    public ADSTcpPlcConnection(InetAddress address, Integer port, AMSNetId targetAmsNetId, AMSPort targetAmsPort) {
        this(address, port, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    public ADSTcpPlcConnection(InetAddress address, AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        this(address, null, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    public ADSTcpPlcConnection(InetAddress address, Integer port, AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort) {
        super(new TcpSocketChannelFactory(address, port != null ? port : TCP_PORT), targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the ads protocol.
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new Plc4X2ADSProtocol(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort));
                pipeline.addLast(new ADS2TcpProtocol());
            }
        };
    }

    protected static AMSPort generateAMSPort() {
        return AMSPort.of(TCP_PORT);
    }

}
