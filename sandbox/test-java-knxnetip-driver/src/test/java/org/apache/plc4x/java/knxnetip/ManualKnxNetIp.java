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
package org.apache.plc4x.java.knxnetip;

import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.base.connection.PcapChannelFactory;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.knxnetip.connection.PassiveKnxNetIpPlcConnection;
import org.apache.plc4x.java.knxnetip.readwrite.KNXNetIPMessage;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketAddress;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketChannelConfig;
import org.apache.plc4x.java.utils.pcapsockets.netty.UdpIpPacketHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.util.List;

public class ManualKnxNetIp {

    private static final Logger logger = LoggerFactory.getLogger(ManualKnxNetIp.class);

    private static NettyPlcConnection connection;

    public static void main(String[] args) throws Exception {
        try {
            connection = new PassiveKnxNetIpPlcConnection(new PcapChannelFactory(
                new File(args[0]), null,
                PassiveKnxNetIpDriver.KNXNET_IP_PORT, PcapSocketAddress.ALL_PROTOCOLS,
                PcapSocketChannelConfig.SPEED_REALTIME, new UdpIpPacketHandler()), "",
                new PlcMessageToMessageCodec<KNXNetIPMessage, PlcRequestContainer>() {

                    @Override
                    protected void decode(ChannelHandlerContext channelHandlerContext, KNXNetIPMessage packet, List<Object> list) throws Exception {
                        System.out.println(packet);
                    }

                    @Override
                    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
                        // Ignore this as we don't send anything.
                    }
                });
            connection.connect();
        } catch (PlcConnectionException e) {
            logger.error("An error occurred starting the BACnet/IP driver", e);
            throw new Exception(e);
        }
    }

}
