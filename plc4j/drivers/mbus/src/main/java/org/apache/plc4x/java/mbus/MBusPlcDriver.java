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
package org.apache.plc4x.java.mbus;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.mbus.field.MBusFieldHandler;
import org.apache.plc4x.java.mbus.protocol.MBusProtocolLogic;
import org.apache.plc4x.java.mbus.readwrite.MBusFrame;
import org.apache.plc4x.java.mbus.readwrite.io.MBusFrameIO;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;

import java.util.function.ToIntFunction;

public class MBusPlcDriver extends GeneratedDriverBase<MBusFrame> {

    @Override
    public String getProtocolCode() {
        return "mbus";
    }

    @Override
    public String getProtocolName() {
        return "Meter-Bus";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return null;
    }

    @Override
    protected PlcFieldHandler getFieldHandler() {
        return new MBusFieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return null;
    }

    @Override
    protected String getDefaultTransport() {
        return null;
    }

    @Override
    protected ProtocolStackConfigurer<MBusFrame> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(MBusFrame.class, MBusFrameIO.class)
            .withPacketSizeEstimator(MBusEstimator.class)
            .withProtocol(MBusProtocolLogic.class)
            .littleEndian()
            .build();
    }

    static class MBusEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 1) {
                int type = byteBuf.getByte(byteBuf.readerIndex()) & 0xFF;
                switch (type) {
                    case 0x68:
                        // length is sent as second filed right after frame type
                        return byteBuf.getByte(byteBuf.readerIndex() + 1) & 0xFF;
                    case 0xE5:
                        return 1;
                    case 0x10: // short frame
                        return 5;
                }
            }

            return -1;
        }
    }
}
