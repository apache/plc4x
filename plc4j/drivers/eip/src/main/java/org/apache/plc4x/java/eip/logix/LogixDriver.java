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
package org.apache.plc4x.java.eip.logix;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.base.protocol.EipProtocolLogic;
import org.apache.plc4x.java.eip.logix.configuration.LogixConfiguration;
import org.apache.plc4x.java.eip.readwrite.EipPacket;
import org.apache.plc4x.java.eip.readwrite.IntegerEncoding;
import org.apache.plc4x.java.eip.base.tag.EipTagHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.*;
import org.apache.plc4x.java.spi.transport.Transport;

import java.util.ServiceLoader;
import java.util.function.Consumer;
import java.util.function.ToIntFunction;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.configure;

public class LogixDriver extends GeneratedDriverBase<EipPacket> {
    public static final int PORT = 44818;

    @Override
    public String getProtocolCode() {
        return "logix";
    }

    @Override
    public String getProtocolName() {
        return "Logix CIP";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return LogixConfiguration.class;
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new EipTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new org.apache.plc4x.java.spi.values.PlcValueHandler();
    }

    @Override
    protected boolean awaitDisconnectComplete() {
        return true;
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
    protected ProtocolStackConfigurer<EipPacket> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(EipPacket.class, EipPacket::staticParse)
            .withProtocol(EipProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withParserArgs(IntegerEncoding.LITTLE_ENDIAN, true)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .littleEndian()
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 4) {
                //Second word for the size and then add the header size 24
                int size = byteBuf.getUnsignedShortLE(byteBuf.readerIndex()+2)+24;
                return size;
            }
            return -1;
        }
    }

     /**Consumes all Bytes till another Magic Byte is found */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {
        @Override
        public void accept(ByteBuf byteBuf) {
            while (byteBuf.getUnsignedByte(0) != 0x00) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

    @Override
    public EipTag prepareTag(String query){
        return EipTag.of(query);
    }

}
