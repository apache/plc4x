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
package org.apache.plc4x.merlot.modbus.dev.api;

import io.netty.buffer.ByteBuf;

public interface ModbusDevice {

    /**
     * Enable / Disables device operation
     * @param Enabled
     */
    public void setEnabled(boolean Enabled);
    
    /**
     * Returns the current status of the device.
     * @return true if the device is enabled, false otherwise.
     */
    public boolean getEnabled();    	
	
    /**
     * Set de the device identification. The firts byte from the Modbus ADU
     * @param ui
     */
    public void setUnitIdentifier(byte UID);
    
    /**
     * Get de the device identification.
     * @return the UID of the device.
     */
    public byte getUnitIdentifier();
    
    /**
     * Set de the device description. 
     * @param description Description of the device.
     */
    public void setUnitDescription(String description);
    
    /**
     * Get de the device Description.
     * @return The description of the device.
     */
    public String getUnitDescription();    
    
    /**
     * Get the "discrete inputs" from the device. Complete reference to memory area.
     * 
     */
    public ByteBuf getDiscreteInputs();
    
    /**
     * Set a particular "discrete input" in the device. Only for internal use.
     * @param register
     * @param state
     */
    public void setDiscreteInput(int register,  boolean state);
    
    /**
     * Get a particular "discrete input" in the device. 
     * @param register
     * @param state
     */
    public boolean getDiscreteInput(int register);    
    
    /**
     * Get the "coil" inputs from the device. Complete reference to memory area.
     * 
     */
    public ByteBuf getCoils(); 
    
    /**
     * Set a particular "coil" in the device. Only for internal & external use.
     * @param register
     * @param state
     */
    public void setCoil(int register,  boolean state);
    
    /**
     * Get a particular "discrete input" in the device. 
     * @param register
     */
    public boolean getCoil(int register);  
    
    /**
     * Get the "input registers" inputs from the device. Complete reference to memory area.
     * 
     */
    public ByteBuf getInputRegisters(); 
    
    /**
     * Set a particular "input register" in the device. Only for internal use.
     * @param register
     * @param value
     */
    public void setInputRegister(int register,  short value);
    
    /**
     * Get a particular "input register" in the device. 
     * @param register
     */
    public short getInputRegister(int register);  
    
    /**
     * Get the "input registers" inputs from the device. Complete reference to memory area.
     * 
     */
    public ByteBuf getHoldingRegisters(); 
    
    /**
     * Set a particular "input register" in the device. Only for internal & external use.
     * @param register
     * @param value
     */
    public void setHoldingRegister(int register,  short value);
    
    /**
     * Get a particular "input register" in the device. 
     * @param register
     */
    public short getHoldingRegister(int register); 
    
    /**
     * Perform a dump of the memory contents in hexadecimal / Ascii.
     * @param index Start register address.
     * @param length number of records to be transferred from registers.
     */
    public String HexDump(int index, int length);    
    
    
}
