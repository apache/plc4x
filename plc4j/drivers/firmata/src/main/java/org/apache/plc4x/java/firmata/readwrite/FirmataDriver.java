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
package org.apache.plc4x.java.firmata.readwrite;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.firmata.readwrite.configuration.FirmataConfiguration;
import org.apache.plc4x.java.firmata.readwrite.context.FirmataDriverContext;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataField;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataFieldHandler;
import org.apache.plc4x.java.firmata.readwrite.protocol.FirmataProtocolLogic;
import org.apache.plc4x.java.spi.configuration.BaseConfiguration;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;

import java.util.function.Consumer;
import java.util.function.ToIntFunction;

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
    protected Class<? extends BaseConfiguration> getConfigurationType() {
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
    protected ProtocolStackConfigurer<FirmataMessage> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(FirmataMessage.class, FirmataMessage::staticParse)
            .withProtocol(FirmataProtocolLogic.class)
            .withDriverContext(FirmataDriverContext.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            // Every incoming message is to be treated as a response.
            .withParserArgs(true)
            .build();
    }

    /**
     * Estimate the Length of a Packet
     */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 1) {
                int type = byteBuf.getByte(byteBuf.readerIndex()) & 0xF0;
                switch (type) {
                    case 0xE0:
                    case 0x90:
                        return 3;
                    case 0xC0:
                    case 0xD0:
                        return 2;
                    case 0xF0: {
                        int commandType = byteBuf.getByte(byteBuf.readerIndex()) & 0x0F;
                        switch (commandType) {
                            case 0x00: {
                                try {
                                    int curPos = 1;
                                    // As long as there are more bytes available and we haven't found the terminating char, continue ...
                                    while ((byteBuf.readableBytes() > curPos + 1) && (byteBuf.getByte(byteBuf.readerIndex() + curPos) != (byte) 0xF7)) {
                                        curPos++;
                                    }
                                    if (byteBuf.getByte(byteBuf.readerIndex() + curPos) == (byte) 0xF7) {
                                        return curPos + 1;
                                    } else {
                                        return -1;
                                    }
                                } catch (Exception e) {
                                    throw new PlcRuntimeException("Invalid packet content", e);
                                }
                            }
                            case 0x04:
                            case 0x05:
                            case 0x09:
                                return 3;
                            case 0x0F:
                                return 1;
                            default:
                                throw new PlcRuntimeException("Invalid command type");
                        }
                    }
                    default:
                        throw new PlcRuntimeException("Invalid packet type");
                }
            }
            return -1;
        }
    }

    /**
     * Consumes all Bytes till one of the potential message type indicators
     */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {

        @Override
        public void accept(ByteBuf byteBuf) {
            while (!isPotentialStart(byteBuf.getUnsignedByte(byteBuf.readerIndex()))) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }

        private boolean isPotentialStart(short value) {
            switch (value & 0xF0) {
                case 0xE0:
                case 0x90:
                case 0xC0:
                case 0xD0:
                case 0xF0:
                    return true;
                default:
                    return false;
            }
        }
    }

    @Override
    public FirmataField prepareField(String query) {
        return FirmataField.of(query);
    }

}
