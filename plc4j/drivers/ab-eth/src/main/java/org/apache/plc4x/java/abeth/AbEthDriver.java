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
package org.apache.plc4x.java.abeth;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.abeth.configuration.AbEthConfiguration;
import org.apache.plc4x.java.abeth.configuration.AbEthTcpTransportConfiguration;
import org.apache.plc4x.java.abeth.tag.AbEthTag;
import org.apache.plc4x.java.abeth.protocol.AbEthProtocolLogic;
import org.apache.plc4x.java.abeth.readwrite.CIPEncapsulationPacket;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.function.ToIntFunction;

public class AbEthDriver extends GeneratedDriverBase<CIPEncapsulationPacket> {

    public static final int AB_ETH_PORT = 2222;

    @Override
    public String getProtocolCode() {
        return "ab-eth";
    }

    @Override
    public String getProtocolName() {
        return "Allen Bradley ETH";
    }


    @Override
    protected Class<? extends PlcConnectionConfiguration> getConfigurationClass() {
        return AbEthConfiguration.class;
    }

    @Override
    protected Optional<Class<? extends PlcTransportConfiguration>> getTransportConfigurationClass(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return Optional.of(AbEthTcpTransportConfiguration.class);
        }
        return Optional.empty();
    }

    @Override
    protected Optional<String> getDefaultTransportCode() {
        return Optional.of("raw");
    }

    @Override
    protected List<String> getSupportedTransportCodes() {
        return Collections.singletonList("tcp");
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
    protected boolean canRead() {
        return true;
    }

    @Override
    protected ProtocolStackConfigurer<CIPEncapsulationPacket> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(CIPEncapsulationPacket.class, CIPEncapsulationPacket::staticParse)
            .withProtocol(AbEthProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .build();
    }

    /**
     * Estimate the Length of a Packet
     */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 4) {
                // In the mspec we subtract 28 from the full size ... so here we gotta add it back.
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 2) + 28;
            }
            return -1;
        }
    }

    @Override
    public PlcTag prepareTag(String tagAddress) {
        return AbEthTag.of(tagAddress);
    }

}
