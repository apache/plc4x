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
package org.apache.plc4x.java.iec608705104.readwrite;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.iec608705104.readwrite.configuration.Iec608705014Configuration;
import org.apache.plc4x.java.iec608705104.readwrite.configuration.Iec608705014TcpTransportConfiguration;
import org.apache.plc4x.java.iec608705104.readwrite.protocol.Iec608705104Protocol;
import org.apache.plc4x.java.iec608705104.readwrite.tag.Iec608705104TagHandler;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationTypeProvider;
import org.apache.plc4x.java.spi.values.PlcValueHandler;

import java.util.function.Consumer;
import java.util.function.ToIntFunction;

/**
 * Implementation of the ADS protocol, based on:
 * - ADS Protocol
 * - TCP
 * - Serial
 */
public class Iec60870514PlcDriver extends GeneratedDriverBase<APDU> implements TransportConfigurationTypeProvider {

    @Override
    public String getProtocolCode() {
        return "iec-60870-5-104";
    }

    @Override
    public String getProtocolName() {
        return "IEC 60870-5-104";
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    public Class<? extends PlcConnectionConfiguration> getConfigurationType() {
        return Iec608705014Configuration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "tcp";
    }

    @Override
    protected Iec608705104TagHandler getTagHandler() {
        return new Iec608705104TagHandler();
    }

    @Override
    protected org.apache.plc4x.java.api.value.PlcValueHandler getValueHandler() {
        return new PlcValueHandler();
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
    protected ProtocolStackConfigurer<APDU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(APDU.class, APDU::staticParse)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .withProtocol(Iec608705104Protocol.class)
            .littleEndian()
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 2) {
                return byteBuf.getUnsignedByte( byteBuf.readerIndex() + 1) + 2;
            }
            return -1;
        }
    }

    /** Consumes all Bytes till another Magic Byte is found */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {
        @Override
        public void accept(ByteBuf byteBuf) {
            while (byteBuf.getUnsignedByte(0) != APDU.STARTBYTE) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

    @Override
    public Class<? extends TransportConfiguration> getTransportConfigurationType(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return Iec608705014TcpTransportConfiguration.class;
        }
        return null;
    }

}
