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
package org.apache.plc4x.java.profinet;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetChannel;
import org.apache.plc4x.java.profinet.config.ProfinetDevices;
import org.apache.plc4x.java.profinet.discovery.ProfinetPlcDiscoverer;
import org.apache.plc4x.java.profinet.protocol.ProfinetProtocolLogic;
import org.apache.plc4x.java.profinet.readwrite.Ethernet_Frame;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.profinet.tag.ProfinetTagHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.pcap4j.core.*;

import java.util.HashMap;
import java.util.function.ToIntFunction;

public class ProfinetDriver extends GeneratedDriverBase<Ethernet_Frame> {

    public static final String DRIVER_CODE = "profinet";

    @Override
    public String getProtocolCode() {
        return DRIVER_CODE;
    }

    @Override
    public String getProtocolName() {
        return "Profinet";
    }

    @Override
    public PlcDriverMetadata getMetadata() {
        return new PlcDriverMetadata() {
            @Override
            public boolean canDiscover() {
                return true;
            }
        };
    }

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        try {
            ProfinetChannel channel = new ProfinetChannel(Pcaps.findAllDevs(), new HashMap<>());
            ProfinetPlcDiscoverer discoverer = new ProfinetPlcDiscoverer(channel);
            channel.setDiscoverer(discoverer);
            return new DefaultPlcDiscoveryRequest.Builder(discoverer);
        } catch (PcapNativeException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return ProfinetConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "raw";
    }

    @Override
    protected boolean awaitSetupComplete() {
        return true;
    }

    /**
     * This protocol doesn't have a disconnect procedure, so there is no need to wait for a login to finish.
     *
     * @return false
     */
    @Override
    protected boolean awaitDisconnectComplete() {
        return false;
    }

    @Override
    protected boolean canRead() {
        return false;
    }

    @Override
    protected boolean canWrite() {
        return false;
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected boolean canBrowse() {
        return true;
    }

    @Override
    protected BaseOptimizer getOptimizer() {
        return new SingleTagOptimizer();
    }

    @Override
    protected ProfinetTagHandler getTagHandler() {
        return new ProfinetTagHandler();
    }

    @Override
    protected  org.apache.plc4x.java.api.value.PlcValueHandler getValueHandler() {
        return new org.apache.plc4x.java.spi.values.PlcValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<Ethernet_Frame> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(Ethernet_Frame.class, Ethernet_Frame::staticParse)
            .withProtocol(ProfinetProtocolLogic.class)
            .withDriverContext(ProfinetDriverContext.class)
            // Every incoming message is to be treated as a response.
            .withParserArgs(true)
            .build();
    }

    /**
     * Estimate the Length of a Packet
     */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 6) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 4) + 6;
            }
            return -1;
        }

    }

    @Override
    public ProfinetTag prepareTag(String query) {
        return ProfinetTag.of(query);
    }

}
