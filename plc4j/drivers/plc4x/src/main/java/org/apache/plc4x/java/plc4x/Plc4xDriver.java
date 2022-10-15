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
package org.apache.plc4x.java.plc4x;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.plc4x.config.Plc4xConfiguration;
import org.apache.plc4x.java.plc4x.field.Plc4xFieldHandler;
import org.apache.plc4x.java.plc4x.protocol.Plc4xProtocolLogic;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xMessage;
import org.apache.plc4x.java.spi.configuration.BaseConfiguration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;

import java.util.function.ToIntFunction;

public class Plc4xDriver extends GeneratedDriverBase<Plc4xMessage> {

    @Override
    public String getProtocolCode() {
        return "plc4x";
    }

    @Override
    public String getProtocolName() {
        return "PLC4X (Proxy-Protocol)";
    }

    @Override
    protected Class<? extends BaseConfiguration> getConfigurationType() {
        return Plc4xConfiguration.class;
    }

    @Override
    protected PlcFieldHandler getFieldHandler() {
        return new Plc4xFieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    @Override
    protected String getDefaultTransport() {
        // TODO: This should be TLS (which we currently don't have yet).
        return "tcp";
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
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected ProtocolStackConfigurer<Plc4xMessage> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(Plc4xMessage.class, Plc4xMessage::staticParse)
            .withProtocol(Plc4xProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 3) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 1);
            }
            return -1;
        }
    }

}
