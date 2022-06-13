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
package org.apache.plc4x.java.transport.test;

import io.netty.bootstrap.Bootstrap;
import io.netty.bootstrap.EmbeddedBootstrap;
import io.netty.channel.Channel;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.NettyChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.SocketAddress;

public class TestChannelFactory extends NettyChannelFactory implements HasConfiguration<TestTransportConfiguration> {

    private static final Logger logger = LoggerFactory.getLogger(TestChannelFactory.class);

    private TestTransportConfiguration configuration;

    public TestChannelFactory(SocketAddress address) {
        super(address);
    }

    @Override
    public void setConfiguration(TestTransportConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public Class<? extends Channel> getChannel() {
        return Plc4xEmbeddedChannel.class;
    }

    @Override
    public boolean isPassive() {
        return false;
    }

    @Override
    protected Bootstrap createBootstrap() {
        return new EmbeddedBootstrap();
    }

    @Override
    public EventLoopGroup getEventLoopGroup() {
        return null;
    }

    @Override
    public void configureBootstrap(Bootstrap bootstrap) {
        bootstrap.localAddress(new TestSocketAddress("lalala"));
        if(configuration != null) {
            logger.info("Configuring Bootstrap with {}", configuration);
        }
    }

}
