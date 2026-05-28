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
package org.apache.plc4x.java.umas.readwrite;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;
import org.apache.plc4x.java.umas.readwrite.configuration.UmasConfiguration;
import org.apache.plc4x.java.umas.readwrite.configuration.UmasTcpTransportConfiguration;
import org.apache.plc4x.java.umas.readwrite.context.UmasDriverContext;
import org.apache.plc4x.java.umas.readwrite.protocol.UmasProtocolLogic;
import org.apache.plc4x.java.umas.readwrite.tag.SymbolicUmasTag;
import org.apache.plc4x.java.umas.readwrite.tag.UmasTag;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.function.ToIntFunction;

/**
 * PLC4J driver for the UMAS protocol (Schneider Electric Modicon PLCs).
 * UMAS is tunneled inside Modbus/TCP using function code 0x5A.
 * Connection URL format: {@code umas:tcp://host:port?unit-identifier=0}
 */
public class UmasDriver extends GeneratedDriverBase<ModbusTcpADU> {

    @Override
    public String getProtocolCode() {
        return "umas";
    }

    @Override
    public String getProtocolName() {
        return "UMAS (Schneider Electric)";
    }

    @Override
    protected Class<? extends PlcConnectionConfiguration> getConfigurationClass() {
        return UmasConfiguration.class;
    }

    @Override
    protected Optional<Class<? extends PlcTransportConfiguration>> getTransportConfigurationClass(String transportCode) {
        if ("tcp".equals(transportCode)) {
            return Optional.of(UmasTcpTransportConfiguration.class);
        }
        return Optional.empty();
    }

    @Override
    protected Optional<String> getDefaultTransportCode() {
        return Optional.of("tcp");
    }

    @Override
    protected List<String> getSupportedTransportCodes() {
        return Collections.singletonList("tcp");
    }

    @Override
    protected boolean awaitSetupComplete() {
        return true;
    }

    @Override
    protected boolean awaitDisconnectComplete() {
        return true;
    }

    @Override
    protected boolean canPing() {
        return true;
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
    protected boolean canBrowse() {
        return true;
    }

    @Override
    protected BaseOptimizer getOptimizer() {
        return new SingleTagOptimizer();
    }

    @Override
    protected ProtocolStackConfigurer<ModbusTcpADU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(
                ModbusTcpADU.class,
                (io) -> {
                    // UMAS responses use function key 0xFE and need the original request's
                    // function key for type discrimination. Peek at the transaction ID from
                    // the MBAP header (first 2 bytes, big-endian) without advancing the
                    // read position, then look up the tracked function key.
                    byte[] header = ((ReadBufferByteBased) io).getBytes(0, 2);
                    int transactionId = ((header[0] & 0xFF) << 8) | (header[1] & 0xFF);
                    short fk = UmasFunctionKeyTracker.consumeFunctionKey(transactionId);
                    return (ModbusTcpADU) ModbusTcpADU.staticParse(io, fk);
                })
            .withProtocol(UmasProtocolLogic.class)
            .withDriverContext(UmasDriverContext.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .build();
    }

    @Override
    public PlcTag prepareTag(String tagAddress) {
        return SymbolicUmasTag.of(tagAddress);
    }

    /**
     * Estimates packet length from the Modbus/TCP MBAP header.
     * Header layout: transactionId(2) + protocolId(2) + length(2) + unitId(1) = 7 bytes.
     * Total size = 6 + length field value.
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

}
