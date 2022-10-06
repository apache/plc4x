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
package org.apache.plc4x.java.s7.readwrite;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.s7.readwrite.configuration.S7Configuration;
import org.apache.plc4x.java.s7.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.s7.readwrite.field.S7Field;
import org.apache.plc4x.java.s7.readwrite.optimizer.S7Optimizer;
import org.apache.plc4x.java.s7.readwrite.protocol.S7ProtocolLogic;
import org.apache.plc4x.java.s7.readwrite.field.S7PlcFieldHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;

import java.util.function.Consumer;
import java.util.function.ToIntFunction;

public class S7Driver extends GeneratedDriverBase<TPKTPacket> {

    public static final int ISO_ON_TCP_PORT = 102;

    @Override
    public String getProtocolCode() {
        return "s7";
    }

    @Override
    public String getProtocolName() {
        return "Siemens S7 (Basic)";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return S7Configuration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "tcp";
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
        return new S7Optimizer();
    }

    @Override
    protected S7PlcFieldHandler getFieldHandler() {
        return new S7PlcFieldHandler();
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
    protected ProtocolStackConfigurer<TPKTPacket> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(TPKTPacket.class, TPKTPacket::staticParse)
            .withProtocol(S7ProtocolLogic.class)
            .withDriverContext(S7DriverContext.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 4) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 2);
            }
            return -1;
        }
    }

    /** Consumes all Bytes till another Magic Byte is found */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {
        @Override
        public void accept(ByteBuf byteBuf) {
            while (byteBuf.getUnsignedByte(0) != TPKTPacket.PROTOCOLID) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

    @Override
    public S7Field prepareField(String query){
        return S7Field.of(query);
    }

}
