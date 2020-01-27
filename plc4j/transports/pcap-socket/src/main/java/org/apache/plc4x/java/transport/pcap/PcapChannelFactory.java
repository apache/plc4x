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
package org.apache.plc4x.java.transport.pcap;

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.*;
import io.netty.channel.oio.OioEventLoopGroup;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.NettyChannelFactory;
import org.apache.plc4x.java.utils.pcapsockets.netty.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class PcapChannelFactory extends NettyChannelFactory implements HasConfiguration<PcapTransportConfiguration> {

    private static final Logger logger = LoggerFactory.getLogger(PcapChannelFactory.class);

    private PcapTransportConfiguration configuration;

    public PcapChannelFactory(PcapSocketAddress address) {
        super(address);
    }

    @Override
    public void setConfiguration(PcapTransportConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public Class<? extends Channel> getChannel() {
        return PcapSocketChannel.class;
    }

    @Override
    public void configureBootstrap(Bootstrap bootstrap) {
        if(configuration != null) {
            logger.info("Configuring Bootstrap with {}", configuration);
            bootstrap.option(PcapSocketChannelOption.PORT, configuration.getDefaultPort());
            bootstrap.option(PcapSocketChannelOption.PROTOCOL_ID, configuration.getProtocolId());
            bootstrap.option(PcapSocketChannelOption.SPEED_FACTOR, configuration.getReplaySpeedFactor());
            if(configuration.getPcapPacketHandler() != null) {
                bootstrap.option(PcapSocketChannelOption.PACKET_HANDLER, configuration.getPcapPacketHandler());
            }
        }
    }

    @Override
    public EventLoopGroup getEventLoopGroup() {
        return new OioEventLoopGroup();
    }

}
