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
package org.apache.plc4x.java.ads;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.configuration.AdsConfiguration;
import org.apache.plc4x.java.ads.discovery.AdsPlcDiscoverer;
import org.apache.plc4x.java.ads.field.AdsFieldHandler;
import org.apache.plc4x.java.ads.protocol.AdsProtocolLogic;
import org.apache.plc4x.java.ads.readwrite.AmsTCPPacket;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;

import java.util.function.ToIntFunction;

/**
 * Implementation of the ADS protocol, based on:
 * - ADS Protocol
 * - TCP
 * - Serial
 */
public class AdsPlcDriver extends GeneratedDriverBase<AmsTCPPacket> {

    @Override
    public String getProtocolCode() {
        return "ads";
    }

    @Override
    public String getProtocolName() {
        return "Beckhoff TwinCat ADS";
    }

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        return new DefaultPlcDiscoveryRequest.Builder(new AdsPlcDiscoverer());
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
    protected boolean canSubscribe() {
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
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    /**
     * This protocol doesn't have a disconnect procedure, so there is no need to wait for a login to finish.
     * @return false
     */
    @Override
    protected boolean awaitDisconnectComplete() {
        return false;
    }

    @Override
    protected ProtocolStackConfigurer<AmsTCPPacket> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(AmsTCPPacket.class, AmsTCPPacket::staticParse)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withProtocol(AdsProtocolLogic.class)
            .littleEndian()
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 6) {
                return (int) byteBuf.getUnsignedIntLE(byteBuf.readerIndex() + 2) + 6;
            }
            return -1;
        }
    }

}
