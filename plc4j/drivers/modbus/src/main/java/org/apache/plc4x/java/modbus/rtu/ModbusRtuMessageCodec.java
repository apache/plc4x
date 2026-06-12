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
package org.apache.plc4x.java.modbus.rtu;

import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.readwrite.ModbusADU;
import org.apache.plc4x.java.modbus.readwrite.ModbusRtuADU;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * MessageCodec for Modbus RTU protocol.
 * Handles the encoding and decoding of Modbus RTU ADU (Application Data Unit) messages.
 *
 * RTU framing: 1 byte address + PDU + 2 bytes CRC.
 * Since RTU has no explicit length field, we attempt to parse the available data
 * and wait for more data if the parse fails.
 */
public class ModbusRtuMessageCodec extends MessageCodecBase<ModbusRtuADU> {

    // Minimum RTU message: address (1) + function code (1) + CRC (2) = 4 bytes
    private static final int MODBUS_RTU_MIN_SIZE = 4;

    public ModbusRtuMessageCodec(TransportInstance<?> transportInstance, Consumer<ModbusRtuADU> messageHandler) {
        super("Modbus RTU", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return MODBUS_RTU_MIN_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        // RTU does not have an explicit length field in the header.
        // We return the total available bytes and let parseMessage handle validation.
        // If the parse fails due to incomplete data, the base class will wait for more data.
        return availableBytes;
    }

    @Override
    protected ModbusRtuADU parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        ModbusADU modbusParsed = ModbusADU.staticParse(readBuffer, DriverType.MODBUS_RTU, true);
        if (!(modbusParsed instanceof ModbusRtuADU modbusADU)) {
            throw new BufferException("Parsed message is not a ModbusRtuADU");
        }
        return modbusADU;
    }

}
