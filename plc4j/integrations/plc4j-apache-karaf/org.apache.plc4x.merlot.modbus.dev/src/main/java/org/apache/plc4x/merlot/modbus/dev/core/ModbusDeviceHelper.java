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
package org.apache.plc4x.merlot.modbus.dev.core;

import org.apache.plc4x.merlot.modbus.dev.api.ModbusDevice;
import org.epics.pvdata.pv.ScalarType;

public class ModbusDeviceHelper {
    
    /*
    *
    */
    public static Double getValue(ModbusDevice mbdev, ScalarType type, int regtype, int index, boolean blnLE){
        Double value = null;
        int byteIndex = (index-1)*2;
        switch(type){
            case pvBoolean:                
                if (regtype == 0){
                    value = mbdev.getCoil(index)?1.0:0.0;
                } else if (regtype == 1){
                    value = mbdev.getDiscreteInput(index)?1.0:0.0;
                }
            break;
            case pvByte:
            case pvUByte:   
                //TODO: 
                byteIndex = index;
                byte temp = 0;
                if (regtype == 3){
                    temp = mbdev.getInputRegisters().getByte(index);    
                } else if (regtype == 4) {
                    temp = mbdev.getHoldingRegisters().getByte(index); 
                }
                value = (Double)(temp * 1.0);
            break; 
            case pvDouble:
                if (regtype == 3){
                    if (blnLE) {                    
                        value = mbdev.getInputRegisters().getDoubleLE(byteIndex);    
                    } else {
                        value = mbdev.getInputRegisters().getDouble(byteIndex);   
                    }
                } else if (regtype == 4) {
                    if (blnLE) {                     
                        value = mbdev.getHoldingRegisters().getDoubleLE(byteIndex); 
                    } else {
                        value = mbdev.getHoldingRegisters().getDouble(byteIndex);                         
                    }
                }
            break;  
            case pvFloat:
                if (regtype == 3){
                    if (blnLE) {                                  
                        value = mbdev.getInputRegisters().getFloatLE(byteIndex) * 1.0;    
                    } else {
                        value = mbdev.getInputRegisters().getFloat(byteIndex) * 1.0;                         
                    }
                } else if (regtype == 4) {
                    if (blnLE) {                        
                        value = mbdev.getHoldingRegisters().getFloatLE(byteIndex) * 1.0; 
                    } else {
                        value = mbdev.getHoldingRegisters().getFloat(byteIndex) * 1.0;                         
                    }
                }                
            break; 
            case pvInt:
            case pvUInt:  
                if (regtype == 3){
                    if (blnLE) {                          
                        value = mbdev.getInputRegisters().getIntLE(byteIndex) * 1.0;    
                    } else {
                        value = mbdev.getInputRegisters().getInt(byteIndex) * 1.0;                          
                    }
                } else if (regtype == 4) {
                    if (blnLE) {                    
                        value = mbdev.getHoldingRegisters().getIntLE(byteIndex) * 1.0; 
                    } else {
                        value = mbdev.getHoldingRegisters().getInt(byteIndex) * 1.0;                         
                    }
                }    
            break;             
            case pvLong:
            case pvULong: 
                if (regtype == 3){
                    if (blnLE) {                        
                        value = mbdev.getInputRegisters().getLongLE(byteIndex) * 1.0;    
                    } else {
                        value = mbdev.getInputRegisters().getLong(byteIndex) * 1.0;                         
                    }
                } else if (regtype == 4) {
                    if (blnLE) {   
                        value = mbdev.getHoldingRegisters().getLongLE(byteIndex) * 1.0; 
                    } else {
                        value = mbdev.getHoldingRegisters().getLong(byteIndex) * 1.0;                         
                    }
                }                  
            break;    
            case pvShort:
            case pvUShort: 
                if (regtype == 3){
                    if (blnLE) {                       
                        value = mbdev.getInputRegisters().getShortLE(byteIndex) * 1.0;    
                    } else {
                        value = mbdev.getInputRegisters().getShort(byteIndex) * 1.0;                          
                    }
                } else if (regtype == 4) {
                    if (blnLE) {                        
                        value = mbdev.getHoldingRegisters().getShortLE(byteIndex) * 1.0; 
                    } else {
                        value = mbdev.getHoldingRegisters().getShort(byteIndex) * 1.0;                         
                    }
                } 
            break;   
            case pvString:
            break;  
           
        }
        return value;        
    }
    
    /*
    *
    */    
    public static void putValue(Double value, ModbusDevice mbdev, ScalarType type, int regtype, int index, boolean blnLE){
        int byteIndex = (index-1)*2;        
        switch(type){
            case pvBoolean:  
                if (regtype == 0){
                    mbdev.setCoil(index, (value>0)?true:false);
                } else if (regtype == 1){
                    mbdev.setDiscreteInput(index, (value>0)?true:false);
                }                
            break;
            case pvByte:
            case pvUByte: { 
                short newValue = value.byteValue();
                byteIndex = index * Byte.BYTES;
                if (regtype == 3){
                    mbdev.getInputRegisters().setByte(byteIndex, newValue);    
                } else if (regtype == 4) {
                    mbdev.getHoldingRegisters().setByte(byteIndex, newValue); 
                }                 
            }
            break; 
            case pvDouble: {              
                if (regtype == 3){
                    if (blnLE) {
                        mbdev.getInputRegisters().setDoubleLE(byteIndex, value);
                    } else {
                        mbdev.getInputRegisters().setDouble(byteIndex, value);
                    }
                } else if (regtype == 4) {
                    if (blnLE) {
                        mbdev.getHoldingRegisters().setDoubleLE(byteIndex, value); 
                    } else {
                        mbdev.getHoldingRegisters().setDouble(byteIndex, value); 
                    }
                }                    
            }
            break;  
            case pvFloat: {
                float newValue = value.floatValue();              
                if (regtype == 3){
                    if (blnLE) {
                        mbdev.getInputRegisters().setFloatLE(byteIndex, newValue);    
                    } else {
                        mbdev.getInputRegisters().setFloat(byteIndex, newValue);
                    }
                } else if (regtype == 4) {
                    if (blnLE) {                    
                        mbdev.getHoldingRegisters().setFloatLE(byteIndex, newValue);
                    } else {
                        mbdev.getHoldingRegisters().setFloat(byteIndex, newValue);                        
                    }
                }   
            }
            break; 
            case pvInt:
            case pvUInt: { 
                int newValue = value.intValue();                 
                if (regtype == 3){
                    if (blnLE) {
                        mbdev.getInputRegisters().setIntLE(byteIndex, newValue);     
                    } else {
                        mbdev.getInputRegisters().setInt(byteIndex, newValue); 
                    }
                } else if (regtype == 4) {
                    if (blnLE) {                    
                        mbdev.getHoldingRegisters().setIntLE(byteIndex, newValue); 
                    } else {
                        mbdev.getHoldingRegisters().setInt(byteIndex, newValue);                         
                    }
                }   
            }
            break;             
            case pvLong:
            case pvULong: {
                long newValue = value.longValue();               
                if (regtype == 3){
                    if (blnLE) {                    
                        mbdev.getInputRegisters().setLongLE(byteIndex, newValue);    
                    } else {
                        mbdev.getInputRegisters().setLong(byteIndex, newValue);                        
                    }
                } else if (regtype == 4) {
                    if (blnLE) {                       
                        mbdev.getHoldingRegisters().setLongLE(byteIndex, newValue);
                    } else {
                        mbdev.getHoldingRegisters().setLong(byteIndex, newValue);                        
                    }
                }   
            }
            break;    
            case pvShort:
            case pvUShort: {   
                short newValue = value.shortValue();
                if (regtype == 3){
                    if (blnLE) {                       
                        mbdev.getInputRegisters().setShortLE(byteIndex, newValue);
                    } else {
                        mbdev.getInputRegisters().setShort(byteIndex, newValue);                        
                    }
                } else if (regtype == 4) {
                    if (blnLE) {                      
                        mbdev.getHoldingRegisters().setShortLE(byteIndex, newValue); 
                    } else {
                        mbdev.getHoldingRegisters().setShort(byteIndex, newValue);                         
                    }
                }   
            }
            break;   
            case pvString:
            break;  
           
        }       
    }
    
    
    
}
