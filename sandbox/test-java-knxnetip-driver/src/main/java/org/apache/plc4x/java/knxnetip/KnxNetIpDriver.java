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

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.base.connection.UdpSocketChannelFactory;
import org.apache.plc4x.java.knxnetip.connection.KnxNetIpConfiguration;
import org.apache.plc4x.java.knxnetip.connection.KnxNetIpFieldHandler;
import org.apache.plc4x.java.knxnetip.protocol.KnxNetIpProtocolLogic;
import org.apache.plc4x.java.knxnetip.readwrite.KNXNetIPMessage;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.NettyChannelFactory;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;

import java.util.function.Function;

public class KnxNetIpDriver extends GeneratedDriverBase<KNXNetIPMessage> {

    public static final int KNXNET_IP_PORT = 3671;

    @Override
    public String getProtocolCode() {
        return "knxnet-ip";
    }

    @Override
    public String getProtocolName() {
        return "KNXNet/IP";
    }

    @Override protected int getDefaultPortIPv4() {
        return KNXNET_IP_PORT;
    }

    @Override protected PlcFieldHandler getFieldHandler() {
        return new KnxNetIpFieldHandler();
    }

    @Override protected Class<? extends NettyChannelFactory> getTransportChannelFactory() {
        return UdpSocketChannelFactory.class;
    }

    @Override protected ProtocolStackConfigurer<KNXNetIPMessage> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(KNXNetIPMessage.class)
            .withProtocol(KnxNetIpProtocolLogic.class)
            .withPacketSizeEstimator(PacketSizeEstimator.class)
            .build();
    }

    public static class PacketSizeEstimator implements Function<ByteBuf, Integer> {

        @Override public Integer apply(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 6) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 4);
            }
            return -1;
        }

    }
}
