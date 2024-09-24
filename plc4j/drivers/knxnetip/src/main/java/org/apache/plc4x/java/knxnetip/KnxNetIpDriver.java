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
package org.apache.plc4x.java.knxnetip;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.knxnetip.discovery.KnxNetIpPlcDiscoverer;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpConfiguration;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpPcapReplayTransportConfiguration;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpRawSocketTransportConfiguration;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpUdpTransportConfiguration;
import org.apache.plc4x.java.knxnetip.context.KnxNetIpDriverContext;
import org.apache.plc4x.java.knxnetip.tag.KnxNetIpTag;
import org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpMessage;
import org.apache.plc4x.java.knxnetip.tag.KnxNetIpTagHandler;
import org.apache.plc4x.java.knxnetip.protocol.KnxNetIpProtocolLogic;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;

import java.util.Arrays;
import java.util.List;
import java.util.Optional;
import java.util.function.ToIntFunction;

public class KnxNetIpDriver extends GeneratedDriverBase<KnxNetIpMessage> {

    public static final int KNXNET_IP_PORT = 3671;

    @Override
    public String getProtocolCode() {
        return "knxnet-ip";
    }

    @Override
    public String getProtocolName() {
        return "KNXNet/IP";
    }

    @Override
    protected Class<? extends PlcConnectionConfiguration> getConfigurationClass() {
        return KnxNetIpConfiguration.class;
    }

    @Override
    protected Optional<Class<? extends PlcTransportConfiguration>> getTransportConfigurationClass(String transportCode) {
        switch (transportCode) {
            case "udp":
                return Optional.of(KnxNetIpUdpTransportConfiguration.class);
            case "pcap":
                return Optional.of(KnxNetIpPcapReplayTransportConfiguration.class);
            case "raw":
                return Optional.of(KnxNetIpRawSocketTransportConfiguration.class);
        }
        return Optional.empty();
    }

    @Override
    protected Optional<String> getDefaultTransportCode() {
        return Optional.of("udp");
    }

    @Override
    protected List<String> getSupportedTransportCodes() {
        return Arrays.asList("udp", "pcap", "raw");
    }

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        return new DefaultPlcDiscoveryRequest.Builder(new KnxNetIpPlcDiscoverer());
    }

    @Override
    protected boolean canPing() {
        return true;
    }

    @Override
    protected boolean canRead() {
        return false;
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
    protected boolean canDiscover() {
        return true;
    }

    @Override
    protected BaseOptimizer getOptimizer() {
        return new SingleTagOptimizer();
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new KnxNetIpTagHandler();
    }

    @Override
    protected boolean awaitDisconnectComplete() { return true; }

    @Override
    protected ProtocolStackConfigurer<KnxNetIpMessage> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(KnxNetIpMessage.class, KnxNetIpMessage::staticParse)
            .withProtocol(KnxNetIpProtocolLogic.class)
            .withDriverContext(KnxNetIpDriverContext.class)
            .withPacketSizeEstimator(PacketSizeEstimator.class)
            .build();
    }

    public static class PacketSizeEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 6) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 4);
            }
            return -1;
        }
    }

    @Override
    public KnxNetIpTag prepareTag(String tagAddress){
        return KnxNetIpTag.of(tagAddress);
    }

}
