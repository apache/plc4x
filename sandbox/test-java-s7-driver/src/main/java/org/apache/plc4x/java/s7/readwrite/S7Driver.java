/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.s7.readwrite;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.s7.readwrite.connection.S7Configuration;
import org.apache.plc4x.java.s7.readwrite.protocol.S7ProtocolLogic;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcFieldHandler;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.connection.NettyChannelFactory;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.osgi.service.component.annotations.Component;

import java.util.function.Consumer;
import java.util.function.Function;

@Component(service = PlcDriver.class, immediate = true)
public class S7Driver extends GeneratedDriverBase<TPKTPacket> {

    private static final int ISO_ON_TCP_PORT = 102;

    @Override
    public String getProtocolCode() {
        return "s7ng";
    }

    @Override
    public String getProtocolName() {
        return "Siemens S7 (Basic)";
    }

    @Override protected int getDefaultPortIPv4() {
        return ISO_ON_TCP_PORT;
    }

    @Override protected S7PlcFieldHandler getFieldHandler() {
        return new S7PlcFieldHandler();
    }

    @Override protected Class<? extends NettyChannelFactory> getTransportChannelFactory() {
        return TcpSocketChannelFactory.class;
    }

    @Override protected ProtocolStackConfigurer<TPKTPacket> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(TPKTPacket.class)
            .withProtocol(S7ProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements Function<ByteBuf, Integer> {

        @Override public Integer apply(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 4) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 2);
            }
            return -1;
        }
    }

    /** Consumes all Bytes till another Magic Byte is found */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {

        @Override public void accept(ByteBuf byteBuf) {
            while (byteBuf.getUnsignedByte(0) != TPKTPacket.PROTOCOLID) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }
}
