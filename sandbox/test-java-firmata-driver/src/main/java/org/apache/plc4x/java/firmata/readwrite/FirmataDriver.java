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
package org.apache.plc4x.java.firmata.readwrite;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.firmata.readwrite.configuration.FirmataConfiguration;
import org.apache.plc4x.java.firmata.readwrite.context.FirmataDriverContext;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldHandler;
import org.apache.plc4x.java.firmata.readwrite.io.FirmataMessageIO;
import org.apache.plc4x.java.firmata.readwrite.protocol.FirmataProtocolLogic;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.osgi.service.component.annotations.Component;

import java.util.HashSet;
import java.util.Set;
import java.util.function.Consumer;
import java.util.function.ToIntFunction;

@Component(service = PlcDriver.class, immediate = true)
public class FirmataDriver extends GeneratedDriverBase<FirmataMessage> {

    @Override
    public String getProtocolCode() {
        return "firmata";
    }

    @Override
    public String getProtocolName() {
        return "Firmata";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return FirmataConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "serial";
    }

    @Override
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected PlcFieldHandler getFieldHandler() {
        return new FirmataFieldHandler();
    }

    @Override
    protected ProtocolStackConfigurer<FirmataMessage> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(FirmataMessage.class, FirmataMessageIO.class)
            .withProtocol(FirmataProtocolLogic.class)
            .withDriverContext(FirmataDriverContext.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            ByteBuf tmp = Unpooled.buffer(1024);
            if (byteBuf.readableBytes() >= 1) {
                int type = byteBuf.getByte(byteBuf.readerIndex()) & 0xF0;
                tmp.writeByte(byteBuf.getByte(byteBuf.readerIndex()));
                switch (type) {
                    case 0xE0:
                    case 0x90: return 3;
                    case 0xC0:
                    case 0xD0: return 2;
                    case 0xF0: {
                        int commandType = byteBuf.getByte(byteBuf.readerIndex()) & 0x0F;
                        switch (commandType) {
                            case 0x00: {
                                try {
                                    int curPos = 1;
                                    // As long as there are more bytes available and we haven't found the terminating char, continue ...
                                    while ((byteBuf.readableBytes() > curPos + 1) && (byteBuf.getByte(byteBuf.readerIndex() + curPos) != (byte) 0xF7)) {
                                        tmp.writeByte(byteBuf.getByte(byteBuf.readerIndex() + curPos));
                                        curPos++;
                                    }
                                    if (byteBuf.getByte(byteBuf.readerIndex() + curPos) == (byte) 0xF7) {
                                        tmp.writeByte(byteBuf.getByte(byteBuf.readerIndex() + curPos));
                                        return curPos + 1;
                                    } else {
                                        return -1;
                                    }
                                } catch(Exception e) {
                                    byte[] data = new byte[tmp.readableBytes()];
                                    tmp.readBytes(data);
                                    System.out.println("Error processing: " + Hex.encodeHexString(data));
                                }
                            }
                            case 0x04:
                            case 0x05:
                            case 0x09:
                                return 3;
                            case 0x0F:
                                return 1;
                        }
                    }
                    default: return 128;
                }
            }
            return -1;
        }
    }

    /** Consumes all Bytes till one of the potential message type indicators */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {

        static Set<Byte> commands = new HashSet<>();
        {
            commands.add((byte) 0xE0);
            commands.add((byte) 0x90);
            commands.add((byte) 0xC0);
            commands.add((byte) 0xD0);
            commands.add((byte) 0xF0);
        }

        @Override
        public void accept(ByteBuf byteBuf) {
            while (!commands.contains((byte) (byteBuf.getUnsignedByte(0) & 0xF0))) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

}
