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
package org.apache.plc4x.java.modbus.ascii;

import io.netty.buffer.ByteBuf;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.modbus.ascii.config.ModbusAsciiConfiguration;
import org.apache.plc4x.java.modbus.ascii.protocol.ModbusAsciiProtocolLogic;
import org.apache.plc4x.java.modbus.base.field.ModbusField;
import org.apache.plc4x.java.modbus.base.field.ModbusFieldHandler;
import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.readwrite.ModbusAsciiADU;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleFieldOptimizer;
import org.apache.plc4x.java.spi.values.PlcValueHandler;

import java.nio.charset.StandardCharsets;
import java.util.function.ToIntFunction;

public class ModbusAsciiDriver extends GeneratedDriverBase<ModbusAsciiADU> {

    @Override
    public String getProtocolCode() {
        return "modbus-ascii";
    }

    @Override
    public String getProtocolName() {
        return "Modbus ASCII";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return ModbusAsciiConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "serial";
    }

    /**
     * Modbus doesn't have a login procedure, so there is no need to wait for a login to finish.
     * @return false
     */
    @Override
    protected boolean awaitSetupComplete() {
        return false;
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
    protected BaseOptimizer getOptimizer() {
        return new SingleFieldOptimizer();
    }

    @Override
    protected ModbusFieldHandler getFieldHandler() {
        return new ModbusFieldHandler();
    }

    @Override
    protected org.apache.plc4x.java.api.value.PlcValueHandler getValueHandler() {
        return new PlcValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<ModbusAsciiADU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(ModbusAsciiADU.class,
                new ModbusAsciiInput(), new ModbusAsciiOutput())
            .withProtocol(ModbusAsciiProtocolLogic.class)
            .withPacketSizeEstimator(ModbusAsciiDriver.ByteLengthEstimator.class)
            // Every incoming message is to be treated as a response.
            .withParserArgs(DriverType.MODBUS_ASCII, true)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 1) {
                return byteBuf.readableBytes();
            }
            return -1;
        }
    }

    @Override
    public ModbusField prepareField(String query){
        return ModbusField.of(query);
    }

    public static class ModbusAsciiInput implements MessageInput<ModbusAsciiADU> {
        @Override
        public ModbusAsciiADU parse(ReadBuffer io, Object... args) throws ParseException {
            final short startChar = io.readShort(8);
            // Check if the message starts with the ":" char.
            if(startChar != 0x3A) {
                throw new ParseException(String.format("Expected starting ':' character but got %c", startChar));
            }
            // Read in all the bytes in the message.
            final ReadBufferByteBased bufferByteBased = (ReadBufferByteBased) io;
            // Read in all bytes except the last two ones, which contain a line-break and carriage-return.
            final byte[] bytes = bufferByteBased.getBytes(bufferByteBased.getPos(), bufferByteBased.getTotalBytes() - 2);
            // Convert the bytes into a string (Which is the hex-encoded message)
            final String inputString = new String(bytes, StandardCharsets.UTF_8);
            // Decode the encoded string back into a byte-array.
            try {
                final byte[] decodedBytes = Hex.decodeHex(inputString);
                // Parse the now decoded bytes as normal Message.
                final ReadBufferByteBased readBuffer = new ReadBufferByteBased(decodedBytes);
                return (ModbusAsciiADU) ModbusAsciiADU.staticParse(readBuffer, DriverType.MODBUS_ASCII, true);
            } catch (DecoderException e) {
                throw new ParseException(String.format("Error parsing incoming message: %s", inputString), e);
            }
        }
    }

    public static class ModbusAsciiOutput implements MessageOutput<ModbusAsciiADU> {
        @Override
        public WriteBufferByteBased serialize(ModbusAsciiADU value, Object... args) throws SerializationException {
            // First serialize the packet the normal way.
            WriteBufferByteBased writeBufferByteBased = new WriteBufferByteBased(value.getLengthInBytes());
            value.serialize(writeBufferByteBased);
            // Get the bytes.
            final byte[] decodedBytes = writeBufferByteBased.getBytes();
            // Now encode each byte as two hex values.
            final String hexString = Hex.encodeHexString(decodedBytes).toUpperCase();
            // Create a new WriteBuffer with the encoded data.
            WriteBufferByteBased encodedWriteBuffer = new WriteBufferByteBased(hexString.length() + 3);
            // Write the leading ":"
            encodedWriteBuffer.writeShort(8, (short) 0x3a);
            encodedWriteBuffer.writeByteArray(hexString.getBytes(StandardCharsets.UTF_8));
            // Write the ending line-break, carriage return.
            encodedWriteBuffer.writeShort(8, (short) 0x0d);
            encodedWriteBuffer.writeShort(8, (short) 0x0a);
            return encodedWriteBuffer;
        }
    }

}