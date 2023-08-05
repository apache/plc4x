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
package org.apache.plc4x.merlot.modbus.dev.impl;


import org.apache.plc4x.merlot.modbus.dev.api.ModbusDevice;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.osgi.service.device.Device;


public class ModbusDeviceImpl implements Device,ModbusDevice{
	

        static final String MODBUS_DEVICE_CATEGORY = "mb1";
    
	private boolean Enabled = false;
	private byte 	UID = 0;
	private ByteBuf DiscreteInputs = null;
	private ByteBuf Coils = null;
	private ByteBuf InputRegisters = null;
	private ByteBuf HoldingRegisters = null;
	private String  UnitDescription = null;
	
	public ModbusDeviceImpl(int nHoldingRegisters) {
            HoldingRegisters = Unpooled.buffer(nHoldingRegisters * 2);
            DiscreteInputs = HoldingRegisters;
            Coils = HoldingRegisters;
            InputRegisters = HoldingRegisters;
	}
	
	public ModbusDeviceImpl(int nDiscreteInputs, int nCoils, int nInputRegisters, int nHoldingRegisters) {
            DiscreteInputs = Unpooled.buffer((nDiscreteInputs / 8));
            Coils = Unpooled.buffer((nCoils / 8));
            InputRegisters = Unpooled.buffer(nInputRegisters * 2);
            HoldingRegisters = Unpooled.buffer(nHoldingRegisters * 2);		
	}

	public void setEnabled(boolean Enabled) {
            this.Enabled = Enabled;
	}

	public boolean getEnabled() {
            return this.Enabled;
	}

	public void setUnitIdentifier(byte UID) {
            this.UID = UID;
	}

	public byte getUnitIdentifier() {
            return this.UID;
	}

	public void setUnitDescription(String description) {	
            this.UnitDescription = description;
		
	}

	public String getUnitDescription() {
            return this.UnitDescription;
	}

	public ByteBuf getDiscreteInputs() {
            return this.DiscreteInputs;
	}

	public void setDiscreteInput(int register, boolean state) {
            int intByte = (register / 8);
            int index = (register % 8);

            int temp = DiscreteInputs.getByte(intByte);
            if (state) {
                temp = temp | 1 << index; // sets 1 at given index
            } else {
                temp = temp & (~(1 << index)); // sets 0 at given index
            }
            DiscreteInputs.setByte(intByte, temp);
	}

	public boolean getDiscreteInput(int register) {
            int intByte = (register / 8);
            int index = (register % 8);
            return (DiscreteInputs.getByte(intByte) & (1 << index)) != 0;	
	}

	public ByteBuf getCoils() {
            return this.Coils;
	}

	public void setCoil(int register, boolean state) {

            int intByte = (register / 8);
            int index = (register % 8);

            int temp = Coils.getByte(intByte);
            if (state) {
                temp = temp | 1 << index; // sets 1 at given index
            } else {
                temp = temp & (~(1 << index)); // sets 0 at given index
            }
            Coils.setByte(intByte, temp);	
		
	}

	public boolean getCoil(int register) {
            int intByte = (register / 8);
            int index = (register % 8);
            return (Coils.getByte(intByte) & (1 << index)) != 0;	
	}

	public ByteBuf getInputRegisters() {
            return this.InputRegisters;
	}

	public void setInputRegister(int register, short value) {
            InputRegisters.setShort(register*2, value);		
	}

	public short getInputRegister(int register) {
            return InputRegisters.getShort(register*2);
	}

	public ByteBuf getHoldingRegisters() {
            return this.HoldingRegisters;
	}

	public void setHoldingRegister(int register, short value) {
            HoldingRegisters.setShort(register*2, value);
	}

	public short getHoldingRegister(int register) {
            return HoldingRegisters.getShort(register*2);
	}

	public String HexDump(int index, int length) {    
            return null;
	}

	public void noDriverFound() {
            this.Enabled = false;
	}
	
    @Override
	public String toString() {
		return "ModbusDeviceImpl [Enabled=" + Enabled + ", UID=" + UID + ", UnitDescription=" + UnitDescription 
				+ ", DiscreteInputs=" + DiscreteInputs.capacity() + ", Coils=" + Coils.capacity() 
				+ ", InputRegisters=" + InputRegisters.capacity() 
				+ ", HoldingRegisters=" + HoldingRegisters.capacity() + "]";
	}	
}
