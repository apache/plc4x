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
package org.apache.plc4x.java.fins.readwrite;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.fins.readwrite.configuration.FinsTcpConfiguration;
import org.apache.plc4x.java.fins.readwrite.context.FinsTcpDriverContext;
import org.apache.plc4x.java.fins.readwrite.optimizer.FinsTcpOptimizer;
import org.apache.plc4x.java.fins.readwrite.protocol.FinsTcpProtocolLogic;
import org.apache.plc4x.java.fins.readwrite.protocol.FinsTcpProtocolStackConfigurer;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.function.Consumer;
import java.util.function.ToIntFunction;

public class FinsTcpDriver extends GeneratedDriverBase<FinsTcpMessage> {

    public static final int ISO_ON_TCP_PORT = 9600;
    private static final Logger log = LoggerFactory.getLogger(FinsTcpDriver.class);

    @Override
    public String getProtocolCode() {
        return "finstcp";
    }

    @Override
    public String getProtocolName() {
        return "omron finsTcp (Basic)";
    }

    @Override
    protected Class<? extends PlcConnectionConfiguration> getConfigurationClass() {
        return FinsTcpConfiguration.class;
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
    protected boolean canDiscover() {
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
    // TODO: Actually this is not quite true ... this is only true for some S7 devices
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected BaseOptimizer getOptimizer() {
        return new FinsTcpOptimizer();
    }

    @Override
    protected ProtocolStackConfigurer getStackConfigurer() {
        return FinsTcpProtocolStackConfigurer.builder(FinsTcpMessage.class, FinsTcpMessage::staticParse)
            .withProtocol(FinsTcpProtocolLogic.class)
            .withDriverContext(FinsTcpDriverContext.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .build();
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


    /**
     * Estimate the Length of a Packet
     */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 4) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 2);
            }
            return -1;
        }
    }

    /**
     * Consumes all Bytes till another Magic Byte is found
     */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {
        @Override
        public void accept(ByteBuf byteBuf) {
            while (byteBuf.getUnsignedByte(0) != FinsTcpMessage.PROTOCOLID) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

}
