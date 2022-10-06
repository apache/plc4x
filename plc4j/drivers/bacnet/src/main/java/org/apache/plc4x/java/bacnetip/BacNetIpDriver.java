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
package org.apache.plc4x.java.bacnetip;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.bacnetip.configuration.BacNetIpConfiguration;
import org.apache.plc4x.java.bacnetip.field.BacNetIpFieldHandler;
import org.apache.plc4x.java.bacnetip.protocol.BacNetIpProtocolLogic;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;

import java.util.function.Consumer;
import java.util.function.ToIntFunction;

public class BacNetIpDriver extends GeneratedDriverBase<BVLC> {

    public static final int BACNET_IP_PORT = 47808;

    @Override
    public String getProtocolCode() {
        return "bacnet-ip";
    }

    @Override
    public String getProtocolName() {
        return "BACnet/IP";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return BacNetIpConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "udp";
    }

    @Override
    protected boolean canRead() {
        return false;
    }

    @Override
    protected boolean canWrite() {
        return false;
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected BacNetIpFieldHandler getFieldHandler() {
        return new BacNetIpFieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<BVLC> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(BVLC.class, BVLC::staticParse)
            .withProtocol(BacNetIpProtocolLogic.class)
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
            while (byteBuf.getUnsignedByte(0) != BVLC.BACNETTYPE) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

}
