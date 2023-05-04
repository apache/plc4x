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
package org.apache.plc4x.java.openprotocol;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.openprotocol.config.OpenProtocolConfiguration;
import org.apache.plc4x.java.openprotocol.protocol.OpenProtocolProtocolLogic;
import org.apache.plc4x.java.openprotocol.readwrite.OpenProtocolMessage;
import org.apache.plc4x.java.openprotocol.tag.OpenProtocolTag;
import org.apache.plc4x.java.openprotocol.tag.OpenProtocolTagHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;
import org.apache.plc4x.java.spi.values.PlcValueHandler;

import java.util.function.ToIntFunction;

public class OpenProtocolDriver extends GeneratedDriverBase<OpenProtocolMessage> {

    @Override
    public String getProtocolCode() {
        return "open-protocol";
    }

    @Override
    public String getProtocolName() {
        return "Open-Protocol";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return OpenProtocolConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "tcp";
    }

    @Override
    public PlcDriverMetadata getMetadata() {
        return () -> false;
    }

    @Override
    protected boolean awaitSetupComplete() {
        return false;
    }

    @Override
    protected boolean awaitDisconnectComplete() {
        return false;
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
    protected BaseOptimizer getOptimizer() {
        return new SingleTagOptimizer();
    }

    @Override
    protected OpenProtocolTagHandler getTagHandler() {
        return new OpenProtocolTagHandler();
    }

    @Override
    protected org.apache.plc4x.java.api.value.PlcValueHandler getValueHandler() {
        return new PlcValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<OpenProtocolMessage> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(OpenProtocolMessage.class,
                OpenProtocolMessage::staticParse)
            .withProtocol(OpenProtocolProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .build();
    }

    /** Estimate the Length of a Packet */
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
    public OpenProtocolTag prepareTag(String tagAddress){
        return OpenProtocolTag.of(tagAddress);
    }

}
