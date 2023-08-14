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
package org.apache.plc4x.merlot.modbus.sim.impl;

import org.apache.plc4x.merlot.modbus.dev.api.ModbusDevice;
import org.apache.plc4x.merlot.modbus.sim.api.ModbusSim;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import io.netty.buffer.ByteBuf;
import java.util.Random;
import org.epics.pvdata.pv.ScalarType;
import parser.MathExpression;


public class ModbusSimImpl implements ModbusSim {

    protected ModbusDevice mbdev;
    protected String strFunction = "";
    protected String[] strVariables = null;
    protected MathExpression eval;
    protected TagPoint tagF = null;
    protected TagPoint tagX = null;
    protected TagPoint tagY = null;
    protected TagPoint tagZ = null;
    protected Random random = null;
      
    protected boolean init = false;
    protected boolean started = false;
    
    protected double f = 0.0;
    protected double x = 0.0;
    protected double y = 0.0;
    protected double z = 0.0;
    

    public ModbusSimImpl(ModbusDevice mbdev) {
        this.mbdev = mbdev;
    }        
        
    @Override
    public void init() {
        if (strVariables.length>0)
        if (strVariables[0] != null){
            tagF = getTag(mbdev, strVariables[0]);
        }
        if (strVariables.length>1)
        if (strVariables[1] != null){
            tagX = getTag(mbdev, strVariables[1]);           
        }
        if (strVariables.length>2)
        if (strVariables[2] != null){          
            tagY = getTag(mbdev, strVariables[2]);        
        }
        if (strVariables.length>3)
        if (strVariables[3] != null){       
            tagZ = getTag(mbdev, strVariables[3]);           
        }                
    }

    @Override
    public void destroy() {
        if (tagF != null) tagF.bb.release();        
        if (tagX != null) tagX.bb.release();
        if (tagY != null) tagY.bb.release();
        if (tagZ != null) tagZ.bb.release();                
    }

    @Override
    public void start() {
        started = true;
    }

    @Override
    public void stop() {
        started = false;
    }
    
    @Override
    public void setFunction(String strFunction) {
    }

    @Override
    public String getFunction() {
        return strFunction;
    }    

    @Override
    public void setVariables(String... variables) {
        this.strVariables = variables;
    }

    @Override
    public String[] getVariables() {
        return this.strVariables;
    }
            
    @Override
    public void execute(JobContext context) {
        throw new UnsupportedOperationException("Not supported yet.");
    }
    
    public double getF(){
        return f;
    };

    @Override
    public void setF(double F) {
       if (tagF != null){
           putValue(F, mbdev, tagF, 0);
       }
    }
    
    public double getX(){
        return x;
    };    
    
    @Override
    public void setX(double X) {
       if (tagX != null){
           putValue(X, mbdev, tagX, 0);
       }
    }    
    
    public double getY(){
        return y;
    };    
    
    @Override
    public void setY(double Y) {
       if (tagY != null){
           putValue(Y, mbdev, tagY, 0);
       }
    }        
    
    public double getZ(){
        return z;
    };     
    
    @Override
    public void setZ(double Z) {
       if (tagZ != null){
           putValue(Z, mbdev, tagZ, 0);
       }
    }    
    
    
    public static Double getValue(ModbusDevice mbdev, TagPoint tag, int index){
        Double value = null;
        int byteIndex = 0;
        switch(tag.stType){
            case pvBoolean:   
                if ("coil".equalsIgnoreCase(tag.mbType)){
                    value = mbdev.getCoil(index)?1.0:0.0;
                } else if ("discrete-input".equalsIgnoreCase(tag.mbType)){
                    value = mbdev.getDiscreteInput(index)?1.0:0.0;
                }
            break;
            case pvByte:
            case pvUByte:   
                //TODO: 
                byteIndex = index;
                value = (Double)((tag.bb.getShort(byteIndex) & 0x00ff) * 1.0);
            break; 
            case pvDouble:
                byteIndex = index * Double.BYTES;
                value = tag.bb.getDouble(0);
            break;  
            case pvFloat:
                byteIndex = index * Float.BYTES;
                value = (Double) (tag.bb.getFloat(byteIndex) * 1.0);
            break; 
            case pvInt:
            case pvUInt:  
                byteIndex = index * Integer.BYTES;
                value = (Double) (tag.bb.getInt(byteIndex) * 1.0);
            break;             
            case pvLong:
            case pvULong: 
                byteIndex = index * Long.BYTES;
                value = (Double) (tag.bb.getLong(byteIndex) * 1.0);                
            break;    
            case pvShort:
            case pvUShort: 
                byteIndex = index * Short.BYTES;
                value = (Double)(tag.bb.getShort(byteIndex) * 1.0);  
            break;   
            case pvString:
            break;  
           
        }
        return value;        
    }
    
    public static void putValue(Double value, ModbusDevice mbdev, TagPoint tag, int index){
        int byteIndex = 0;        
        switch(tag.stType){
            case pvBoolean:  
                if ("coil".equalsIgnoreCase(tag.mbType)){
                    mbdev.setCoil(index, (value>0)?true:false);
                } else if ("discrete-input".equalsIgnoreCase(tag.mbType)){
                    mbdev.setDiscreteInput(index, (value>0)?true:false);
                }                
            break;
            case pvByte:
            case pvUByte: { 
                short newValue = value.byteValue();
                byteIndex = index * Byte.BYTES;
                tag.bb.setShort(byteIndex, newValue);
            }
            break; 
            case pvDouble: {
                byteIndex = index * Double.BYTES;                
                tag.bb.setDouble(byteIndex, value);
            }
            break;  
            case pvFloat: {
                float newValue = value.floatValue();
                byteIndex = index * Float.BYTES;                 
                tag.bb.setFloat(byteIndex, newValue);
            }
            break; 
            case pvInt:
            case pvUInt: { 
                int newValue = value.intValue();
                byteIndex = index * Integer.BYTES;                   
                tag.bb.setInt(byteIndex, newValue);
            }
            break;             
            case pvLong:
            case pvULong: {
                long newValue = value.longValue();
                byteIndex = index * Long.BYTES;                 
                tag.bb.setLong(byteIndex, newValue);
            }
            break;    
            case pvShort:
            case pvUShort: {   
                short newValue = value.shortValue();
                byteIndex = index * Short.BYTES;
                tag.bb.setShort(byteIndex, newValue);
            }
            break;   
            case pvString:
            break;  
           
        }       
    }
    
    public  TagPoint getTag(ModbusDevice mbdev, String address){
        String[] fields = address.split(":");
        String scalar = "boolean";
        int mbAddress = -1;        
        int range = -1;
        int start = -1;
        int end = -1;  
        
        if (fields.length == 2) {
            start = fields[1].indexOf("[");
            end = fields[1].indexOf("]"); 
            if ((start == -1) && (end == -1)) {
                mbAddress = Integer.parseInt(fields[1]);
                range = 1;
            } else if ((start > 0) && (end > start)) {
                mbAddress = Integer.parseInt(fields[1].substring(0, start));
                start += 1;
                String strRange = fields[1].substring(start,end);
                range = Integer.parseInt(strRange);
            }            
        } else if (fields.length > 2) {
            start = fields[2].indexOf("[");
            end = fields[2].indexOf("]"); 
            
            mbAddress = Integer.parseInt(fields[1]);
            
            if ((start == -1) && (end == -1)) {
                scalar = fields[2];
                range = 1;
            } else if ((start > 0) && (end > start)) {
                scalar = fields[2].substring(0, start);
                range = Integer.parseInt(fields[2].substring(start+1,end));
            }            
        }
        
        TagPoint tag = new TagPoint();
        
        if ("coil".equalsIgnoreCase(fields[0])){
            tag.mbType = "coil";
        }

        if ("discrete-input".equalsIgnoreCase(fields[0])){
            tag.mbType = "discrete-input";
        }

        if ("input-register".equalsIgnoreCase(fields[0])){
            tag.mbType = "input-register";
        }

        if ("holding-register".equalsIgnoreCase(fields[0])){
            tag.mbType = "holding-register";
        }
        
        tag.mbAddress = mbAddress;
        tag.stType = ScalarType.getScalarType(scalar);        
        tag.range = range;
        
        int size_in_bytes = 0;
        switch(tag.stType){
            case pvBoolean:                   
            break;
            case pvByte:
            case pvUByte: { 
                start = tag.mbAddress*2;
                end = start + tag.range * Byte.BYTES;
                if (tag.mbType.equalsIgnoreCase("input-register")){
                    tag.bb = mbdev.getInputRegisters().retainedSlice(start, end);
                } else {
                    tag.bb = mbdev.getHoldingRegisters().retainedSlice(start, end);                    
                }
            }
            break; 
            case pvDouble: { 
                start = tag.mbAddress*2;
                end = start + tag.range * Double.BYTES;
                if (tag.mbType.equalsIgnoreCase("input-register")){
                    tag.bb = mbdev.getInputRegisters().retainedSlice(start, end);
                } else {
                    tag.bb = mbdev.getHoldingRegisters().retainedSlice(start, end);                    
                }
            }
            break;  
            case pvFloat: { 
                start = (tag.mbAddress-1)*2;
                end = start + tag.range * Float.BYTES;
                if (tag.mbType.equalsIgnoreCase("input-register")){
                    tag.bb = mbdev.getInputRegisters().retainedSlice(start, end);
                } else {
                    System.out.println(start + " : " + end);
                    tag.bb = mbdev.getHoldingRegisters().retainedSlice(start, end);                    
                }
            }
            break; 
            case pvInt:
            case pvUInt:  { 
                start = (tag.mbAddress-1)*2;
                end = start + tag.range * Integer.BYTES;
                if (tag.mbType.equalsIgnoreCase("input-register")){
                    tag.bb = mbdev.getInputRegisters().retainedSlice(start, end);
                } else {
                    tag.bb = mbdev.getHoldingRegisters().retainedSlice(start, end);                    
                }
            }
            break;             
            case pvLong:
            case pvULong: { 
                start = (tag.mbAddress-1)*2;
                end = start + tag.range * Long.BYTES;
                if (tag.mbType.equalsIgnoreCase("input-register")){
                    tag.bb = mbdev.getInputRegisters().retainedSlice(start, end);
                } else {
                    tag.bb = mbdev.getHoldingRegisters().retainedSlice(start, end);                    
                }
            }
            break;    
            case pvShort:
            case pvUShort: { 
                start = (tag.mbAddress-1)*2;
                end = start + tag.range * Short.BYTES;              
                if (tag.mbType.equalsIgnoreCase("input-register")){
                    tag.bb = mbdev.getInputRegisters().retainedSlice(start, end);
                } else {
                    tag.bb = mbdev.getHoldingRegisters().retainedSlice(start, end);                    
                }
            }
            break;   
            case pvString:
            break;
  
        }            
        
        return tag;
    }
   
    class TagPoint {
        public ByteBuf bb;
        public String mbType;
        public ScalarType stType;
        public int mbAddress;  
        public int range;
    }
    
}
