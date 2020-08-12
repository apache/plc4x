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
package org.apache.plc4x.java.amsads;

import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPipeline;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.amsads.configuration.AdsConfiguration;
import org.apache.plc4x.java.amsads.field.AdsFieldHandler;
import org.apache.plc4x.java.amsads.protocol.AdsProtocolLogic;
import org.apache.plc4x.java.amsads.readwrite.AmsPacket;
import org.apache.plc4x.java.amsads.readwrite.io.AmsPacketIO;
import org.apache.plc4x.java.spi.GeneratedDriverByteToMessageCodec;
import org.apache.plc4x.java.spi.Plc4xNettyWrapper;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;

import java.util.Iterator;
import java.util.List;
import java.util.Map;

/**
 * Implementation of the ADS protocol, based on:
 * - ADS Protocol
 * - TCP
 * - Serial
 */
public class AMSADSPlcDriver extends GeneratedDriverBase<AmsPacket> {

    public static final int TCP_PORT = 48898;

    @Override
    public String getProtocolCode() {
        return "ads";
    }

    @Override
    public String getProtocolName() {
        return "Beckhoff TwinCat ADS";
    }

    @Override
    protected boolean canRead() {
        return true;
    }

    @Override
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return AdsConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "tcp";
    }

    @Override
    protected AdsFieldHandler getFieldHandler() {
        return new AdsFieldHandler();
    }

    @Override
    protected ProtocolStackConfigurer<AmsPacket> getStackConfigurer() {
        return new AdsSwitchingStackConfigurer(
            SingleProtocolStackConfigurer.builder(AmsPacket.class, AmsPacketIO.class)
                .withProtocol(AdsProtocolLogic.class)
                .littleEndian()
                .build());
    }

    /**
     * Custom Configurer for Switch TCP / Serial.
     */
    static class AdsSwitchingStackConfigurer implements ProtocolStackConfigurer<AmsPacket> {

        private final ProtocolStackConfigurer<AmsPacket> delegate;

        public AdsSwitchingStackConfigurer(ProtocolStackConfigurer<AmsPacket> delegate) {
            this.delegate = delegate;
        }

        @Override
        public Plc4xProtocolBase<AmsPacket> configurePipeline(Configuration configuration, ChannelPipeline pipeline, boolean passive) {
            final Plc4xProtocolBase<AmsPacket> protocolBase = delegate.configurePipeline(configuration, pipeline, passive);
            Iterator<Map.Entry<String, ChannelHandler>> iterator = pipeline.iterator();
            while (iterator.hasNext()) {
                final Map.Entry<String, ChannelHandler> entry = iterator.next();
                final Class<? extends ChannelHandler> aClass = entry.getValue().getClass();

                System.out.println(aClass);

                if (entry.getValue() instanceof Plc4xNettyWrapper) {
                    // Found handler
                    pipeline.addBefore(entry.getKey(), "idempotent-layer", new MessageToMessageCodec() {
                        @Override
                        protected void decode(ChannelHandlerContext channelHandlerContext, Object o, List list) throws Exception {
                            list.add(o);
                        }

                        @Override
                        protected void encode(ChannelHandlerContext channelHandlerContext, Object o, List list) throws Exception {
                            list.add(o);
                        }
                    });
                }

            }

            // After
            System.out.println("After");
            iterator = pipeline.iterator();
            while (iterator.hasNext()) {
                final Map.Entry<String, ChannelHandler> entry = iterator.next();
                final Class<? extends ChannelHandler> aClass = entry.getValue().getClass();

                System.out.println(aClass);
            }
            return protocolBase;
        }
    }

}
