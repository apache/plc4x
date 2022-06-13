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
package org.apache.plc4x.java.cbus;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.cbus.configuration.CBusConfiguration;
import org.apache.plc4x.java.cbus.context.CBusDriverContext;
import org.apache.plc4x.java.cbus.protocol.CBusProtocolLogic;
import org.apache.plc4x.java.cbus.readwrite.CBusCommand;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;

import java.util.function.Consumer;
import java.util.function.ToIntFunction;

public class CBusDriver extends GeneratedDriverBase<CBusCommand> {

    @Override
    public String getProtocolCode() {
        return "c-bus";
    }

    @Override
    public String getProtocolName() {
        return "Clipsal C-Bus";
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
    protected Class<? extends Configuration> getConfigurationType() {
        return CBusConfiguration.class;
    }

    @Override
    protected PlcFieldHandler getFieldHandler() {
        return null;
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return null;
    }

    @Override
    protected ProtocolStackConfigurer<CBusCommand> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(CBusCommand.class, CBusCommand::staticParse)
            .withProtocol(CBusProtocolLogic.class)
            .withDriverContext(CBusDriverContext.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .build();
    }

    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            for(int i = 0; i < byteBuf.readableBytes() - 1; i++) {
                if((byteBuf.getUnsignedByte(i) == (short) 0x0D) && (byteBuf.getUnsignedByte(i + 1) == (short) 0x0A)) {
                    return i + 1;
                }
            }
            return -1;
        }
    }

    /** Consumes all Bytes till a backslash is found */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {
        @Override
        public void accept(ByteBuf byteBuf) {
            // Consume every byte until the next byte would be a backslash.
            while (byteBuf.getUnsignedByte(0) != '\\') {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

}
