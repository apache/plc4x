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

import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.readwrite.ModbusADU;
import org.apache.plc4x.java.modbus.readwrite.ModbusAsciiADU;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;

import java.nio.charset.StandardCharsets;
import java.util.function.Consumer;

/**
 * MessageCodec for Modbus ASCII protocol.
 * Handles the encoding and decoding of Modbus ASCII ADU (Application Data Unit) messages.
 *
 * <p>ASCII framing: each frame starts with ':' (0x3A), followed by hex-encoded binary data
 * (address + PDU + LRC), and ends with CR+LF. The codec handles the hex encoding/decoding
 * layer between the binary ADU representation and the wire format.</p>
 *
 * <p>The minimum frame is: ':' + address (2 hex chars) + function code (2 hex chars) + LRC (2 hex chars) + CR + LF = 9 bytes.</p>
 */
public class ModbusAsciiMessageCodec extends MessageCodecBase<ModbusAsciiADU> {

    // Minimum ASCII message: ':' (1) + address (2) + function (2) + LRC (2) + CR (1) + LF (1) = 9 bytes
    private static final int MODBUS_ASCII_MIN_SIZE = 9;

    private static final char[] HEX_CHARS = "0123456789ABCDEF".toCharArray();

    public ModbusAsciiMessageCodec(TransportInstance<?> transportInstance, Consumer<ModbusAsciiADU> messageHandler) {
        super("Modbus ASCII", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return MODBUS_ASCII_MIN_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        // ASCII framing: scan for CR+LF terminator to determine message boundary.
        try {
            byte[] data = transportInstance.peekReadableBytes(availableBytes);
            for (int i = 0; i < data.length - 1; i++) {
                if (data[i] == '\r' && data[i + 1] == '\n') {
                    return i + 2; // include CR+LF in the message
                }
            }
            // CR+LF not found yet — wait for more data
            return -1;
        } catch (Exception e) {
            return -1;
        }
    }

    @Override
    public void send(ModbusAsciiADU message) throws MessageCodecException {
        try {
            // Serialize the ADU to binary bytes (address + PDU + LRC)
            int messageSize = message.getLengthInBytes();
            WriteBufferByteBased writeBuffer = createWriteBuffer(messageSize);
            message.serialize(writeBuffer);
            byte[] binaryBytes = writeBuffer.getBytes();

            // Encode as ASCII frame: ':' + hex(binary) + CR + LF
            StringBuilder sb = new StringBuilder(1 + binaryBytes.length * 2 + 2);
            sb.append(':');
            for (byte b : binaryBytes) {
                sb.append(HEX_CHARS[(b >> 4) & 0x0F]);
                sb.append(HEX_CHARS[b & 0x0F]);
            }
            sb.append('\r');
            sb.append('\n');

            byte[] asciiBytes = sb.toString().getBytes(StandardCharsets.US_ASCII);

            if (logger.isTraceEnabled()) {
                logger.trace("Sending Modbus ASCII frame: {}", sb.substring(0, sb.length() - 2));
            }

            transportInstance.write(asciiBytes);
        } catch (BufferException e) {
            throw new MessageCodecException("Failed to serialize Modbus ASCII message", e);
        } catch (TransportException e) {
            throw new MessageCodecException("Failed to send Modbus ASCII message", e);
        }
    }

    @Override
    protected ModbusAsciiADU parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        // The readBuffer contains the raw ASCII frame: ':' + hex data + CR + LF
        byte[] asciiFrame = readBuffer.getBytes();

        // Validate frame structure
        if (asciiFrame[0] != ':') {
            throw new BufferException("Modbus ASCII frame must start with ':'");
        }

        // Extract hex data between ':' and CR+LF
        // Frame: ':' (1 byte) + hex data (N bytes) + CR (1 byte) + LF (1 byte)
        int hexStart = 1;
        int hexEnd = asciiFrame.length - 2; // exclude CR+LF
        int hexLen = hexEnd - hexStart;
        if (hexLen < 6 || hexLen % 2 != 0) {
            throw new BufferException("Invalid Modbus ASCII frame: hex data length=" + hexLen);
        }

        // Decode hex to binary
        byte[] binaryData = new byte[hexLen / 2];
        for (int i = 0; i < binaryData.length; i++) {
            int hi = Character.digit(asciiFrame[hexStart + i * 2], 16);
            int lo = Character.digit(asciiFrame[hexStart + i * 2 + 1], 16);
            if (hi < 0 || lo < 0) {
                throw new BufferException("Invalid hex character in Modbus ASCII frame at position " + (hexStart + i * 2));
            }
            binaryData[i] = (byte) ((hi << 4) | lo);
        }

        // Parse the binary ADU
        ReadBufferByteBased binaryBuffer = new ReadBufferByteBased(binaryData);
        ModbusADU modbusParsed = ModbusADU.staticParse(binaryBuffer, DriverType.MODBUS_ASCII, true);
        if (!(modbusParsed instanceof ModbusAsciiADU modbusADU)) {
            throw new BufferException("Parsed message is not a ModbusAsciiADU");
        }
        return modbusADU;
    }

}
