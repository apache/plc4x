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
import org.apache.plc4x.java.modbus.ascii.context.ModbusAsciiContext;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.modbus.ascii.config.ModbusAsciiConfiguration;
import org.apache.plc4x.java.modbus.ascii.protocol.ModbusAsciiProtocolLogic;
import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.readwrite.ModbusAsciiADU;
import org.apache.plc4x.java.modbus.tcp.config.ModbusTcpTransportConfiguration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;
import java.util.function.Consumer;
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
    protected Class<? extends PlcConnectionConfiguration> getConfigurationClass() {
        return ModbusAsciiConfiguration.class;
    }

    @Override
    protected Optional<Class<? extends PlcTransportConfiguration>> getTransportConfigurationClass(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return Optional.of(ModbusTcpTransportConfiguration.class);
        }
        return Optional.empty();
    }

    @Override
    protected Optional<String> getDefaultTransportCode() {
        return Optional.of("serial");
    }

    @Override
    protected List<String> getSupportedTransportCodes() {
        return Arrays.asList("tcp", "serial");
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
    protected boolean canPing() {
        return true;
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
        return new SingleTagOptimizer();
    }

    @Override
    protected ProtocolStackConfigurer<ModbusAsciiADU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(ModbusAsciiADU.class, new ModbusAsciiInput(), new ModbusAsciiOutput())
            .withProtocol(ModbusAsciiProtocolLogic.class)
            .withDriverContext(ModbusAsciiContext.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            // A Modbus ASCII packet has the absolute minimum size of 9 (if it has absolutely no payload)
            // That's one starting character ":" two chars at the end CR + LF
            // The message itself is encoded as two ascii digits of a hex-encoded byte and the minimum of
            // a Modbus ASCII PDU is 3 bytes, which is 6 bytes in hex-encoded ascii string.
            if (byteBuf.readableBytes() >= 9) {
                // Fetch what's currently in the buffer
                byte[] buf = new byte[byteBuf.readableBytes()];
                byteBuf.getBytes(byteBuf.readerIndex(), buf);
                ReadBufferByteBased reader = new ReadBufferByteBased(buf);

                // Try to parse the buffer content.
                try {
                    ModbusAsciiInput input = new ModbusAsciiInput();
                    ModbusAsciiADU modbusADU = input.parse(reader);

                    // Make sure we only read one message.
                    return (modbusADU.getLengthInBytes() * 2) + 3;
                } catch (ParseException e) {
                    return -1;
                }
                // If we're getting this error, manually compact the buffer.
                // Hopefully now there will be enough space for another attempt.
                catch (ArrayIndexOutOfBoundsException e) {
                    byteBuf.discardReadBytes();
                    return -1;
                }
            }
            return -1;
        }
    }

    /**
     * Consumes all Bytes until the starting character ":" is found
     */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {
        @Override
        public void accept(ByteBuf byteBuf) {
            while (byteBuf.getUnsignedByte(0) != 0x3A) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

    @Override
    public ModbusTag prepareTag(String tagAddress){
        return ModbusTag.of(tagAddress);
    }

    public static class ModbusAsciiInput implements MessageInput<ModbusAsciiADU> {
        @Override
        public ModbusAsciiADU parse(ReadBuffer io) throws ParseException {
            // A Modbus ASCII message starts with an ASCII character ":" and is ended by two characters CRLF
            // (Carriage-Return + Line-Feed)
            // The actual payload is that each byte of the message is encoded by a string representation of it's
            // two hex-digits.
            final short startChar = io.readShort(8);
            // Check if the message starts with the ":" char.
            if(startChar != 0x3A) {
                throw new ParseException(String.format("Expected starting ':' character but got %c", startChar));
            }
            // Read in all the bytes in the message.
            final ReadBufferByteBased bufferByteBased = (ReadBufferByteBased) io;
            // Read in all bytes except the last two ones, which contain a carriage-return and line-break.
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
        public WriteBufferByteBased serialize(ModbusAsciiADU value) throws SerializationException {
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