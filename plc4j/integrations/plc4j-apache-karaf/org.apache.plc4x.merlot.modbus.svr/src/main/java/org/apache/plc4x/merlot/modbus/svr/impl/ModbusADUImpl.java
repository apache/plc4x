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
package org.apache.plc4x.merlot.modbus.svr.impl;

import org.apache.plc4x.merlot.modbus.svr.api.ModbusADU;
import io.netty.buffer.ByteBuf;
import static io.netty.buffer.Unpooled.directBuffer;

/**
 *
 * @author cgarcia
 *
 * Modbus TCP/IP Application Data Unit (ADU)
 * +-----------------------------------------------------------------------+ |
 * Transaction Identifier | 2 bytes | For synchronization between | | | |
 * messages of server& client |
 * +-----------------------------------------------------------------------+ |
 * Protocol Identifier | 2 bytes | Zero for MODBUS/TCP |
 * +-----------------------------------------------------------------------+ |
 * Length Field | 2 bytes | Number of remaining bytes in this | | | | frame |
 * +-----------------------------------------------------------------------+ |
 * Unit Identifier | 1 byte | Slave Address (255 if not used) |
 * +-----------------------------------------------------------------------+ |
 * Function code | 1 byte | Function codes as in other variants|
 * +-----------------------------------------------------------------------+ |
 * Data bytes | n bytes | Data as response or commands |
 * +-----------------------------------------------------------------------+
 */
public class ModbusADUImpl implements ModbusADU {

    public ModbusADUImpl() {
        super();
    }

    short TransactionID = 0;
    short ProtocolID = 0;
    short LengthField = 0;
    byte UnitID = 0;
    byte FunctionCode = 0;
    ByteBuf Data = null;

    public void setTransactionID(short TransactionID) {
        this.TransactionID = TransactionID;
    }

    public short getTransactionID() {
        return this.TransactionID;
    }

    public void setProtocolID(short ProtocolID) {
        this.ProtocolID = ProtocolID;
    }

    public short getProtocolID() {
        return this.ProtocolID;
    }

    public void setLengthField(short LengthField) {
        this.LengthField = LengthField;
    }

    public short getLengthField() {
        return this.LengthField;
    }

    public void setUnitID(byte UnitID) {
        this.UnitID = UnitID;
    }

    public byte getUnitID() {
        return this.UnitID;
    }

    public void setFunctionCode(byte FunctionCode) {
        this.FunctionCode = FunctionCode;
    }

    public byte getFunctionCode() {
        return this.FunctionCode;
    }

    public void setData(ByteBuf Data) {
        this.Data = Data;
    }

    public ByteBuf getData() {
        return this.Data;
    }

    public ByteBuf getSerialized() {
        ByteBuf buffer = directBuffer(this.getLengthField() + 6);
        buffer.writeShort(this.getTransactionID());
        buffer.writeShort(this.getProtocolID());
        buffer.writeShort(this.getLengthField());
        buffer.writeByte(this.getUnitID());
        buffer.writeByte(this.getFunctionCode());
        this.Data.markReaderIndex();
        buffer.writeBytes(this.getData());
        this.Data.resetReaderIndex();
        return buffer;
    }

}
