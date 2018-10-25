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
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.protocol.Ads2PayloadProtocol;
import org.apache.plc4x.java.ads.protocol.Payload2SerialProtocol;
import org.apache.plc4x.java.ads.protocol.Plc4x2AdsProtocol;
import org.apache.plc4x.java.ads.protocol.util.SingleMessageRateLimiter;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.base.connection.SerialChannelFactory;
import org.apache.plc4x.java.base.protocol.SingleItemToSingleRequestProtocol;

import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class AdsSerialPlcConnection extends AdsAbstractPlcConnection {

    private AdsSerialPlcConnection(String serialPort, AmsNetId targetAmsNetId, AmsPort targetAmsPort) {
        this(serialPort, targetAmsNetId, targetAmsPort, generateAMSNetId(), generateAMSPort());
    }

    private AdsSerialPlcConnection(String serialPort, AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort) {
        super(new SerialChannelFactory(serialPort), targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    public static AdsSerialPlcConnection of(String serialPort, AmsNetId targetAmsNetId, AmsPort targetAmsPort) {
        return new AdsSerialPlcConnection(serialPort, targetAmsNetId, targetAmsPort);
    }

    public static AdsSerialPlcConnection of(String serialPort, AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort) {
        return new AdsSerialPlcConnection(serialPort, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the ads protocol.
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new Payload2SerialProtocol());
                pipeline.addLast(new SingleMessageRateLimiter());
                pipeline.addLast(new Ads2PayloadProtocol());
                pipeline.addLast(new Plc4x2AdsProtocol(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, fieldMapping));
                pipeline.addLast(new SingleItemToSingleRequestProtocol(AdsSerialPlcConnection.this, AdsSerialPlcConnection.this, timer));
            }
        };
    }

}
