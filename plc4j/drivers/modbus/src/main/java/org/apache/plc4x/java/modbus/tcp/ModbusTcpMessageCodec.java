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
package org.apache.plc4x.java.modbus.tcp;

import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.readwrite.ModbusADU;
import org.apache.plc4x.java.modbus.readwrite.ModbusTcpADU;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * MessageCodec for Modbus TCP protocol.
 * Handles the encoding and decoding of Modbus TCP ADU (Application Data Unit) messages.
 */
public class ModbusTcpMessageCodec extends MessageCodecBase<ModbusTcpADU> {

    private static final int MODBUS_TCP_HEADER_SIZE = 6; // Transaction ID (2) + Protocol ID (2) + Length (2)

    public ModbusTcpMessageCodec(TransportInstance<?> transportInstance, Consumer<ModbusTcpADU> messageHandler) {
        super("Modbus TCP", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return MODBUS_TCP_HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        // Extract the length field (bytes 4-5, big-endian)
        // Length includes unit ID (1 byte) + PDU length
        int length = ((header[4] & 0xFF) << 8) | (header[5] & 0xFF);
        return MODBUS_TCP_HEADER_SIZE + length;
    }

    @Override
    protected ModbusTcpADU parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        ModbusADU modbusParsed = ModbusADU.staticParse(readBuffer, DriverType.MODBUS_TCP, true);
        if (!(modbusParsed instanceof ModbusTcpADU modbusADU)) {
            throw new BufferException("Parsed message is not a ModbusTcpADU");
        }
        return modbusADU;
    }

}
