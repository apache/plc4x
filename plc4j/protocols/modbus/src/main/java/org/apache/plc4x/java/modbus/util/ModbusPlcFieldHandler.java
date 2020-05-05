/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at
 
     http://www.apache.org/licenses/LICENSE-2.0
 
 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.modbus.util;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneOffset;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.DefaultPlcFieldHandler;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultBooleanFieldItem;
import org.apache.plc4x.java.modbus.messages.items.DefaultModbusByteArrayFieldItem;
import org.apache.plc4x.java.modbus.model.*;

public class ModbusPlcFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) throws PlcInvalidFieldException {
        if (MaskWriteRegisterModbusField.ADDRESS_PATTERN.matcher(fieldQuery).matches()) {
            return MaskWriteRegisterModbusField.of(fieldQuery);
        } else if (ReadDiscreteInputsModbusField.ADDRESS_PATTERN.matcher(fieldQuery).matches()) {
            return ReadDiscreteInputsModbusField.of(fieldQuery);
        } else if (ReadHoldingRegistersModbusField.ADDRESS_PATTERN.matcher(fieldQuery).matches()) {
            return ReadHoldingRegistersModbusField.of(fieldQuery);
        } else if (ReadInputRegistersModbusField.ADDRESS_PATTERN.matcher(fieldQuery).matches()) {
            return ReadInputRegistersModbusField.of(fieldQuery);
        } else if (CoilModbusField.ADDRESS_PATTERN.matcher(fieldQuery).matches()) {
            return CoilModbusField.of(fieldQuery);
        } else if (RegisterModbusField.ADDRESS_PATTERN.matcher(fieldQuery).matches()) {
            return RegisterModbusField.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public BaseDefaultFieldItem encodeBoolean(PlcField field, Object[] values) {
        ModbusField modbusField = (ModbusField) field;
        List<Boolean> booleanValues = new LinkedList<>();
        for (Object value : values) {
            if (value instanceof Boolean) {
                Boolean booleanValue = (Boolean) value;
                booleanValues.add(booleanValue);
            } else if (value instanceof Byte) {
                Byte byteValue = (Byte) value;
                BitSet bitSet = BitSet.valueOf(new byte[]{byteValue});
                for (int i = 0; i < 8; i++) {
                    booleanValues.add(bitSet.get(i));
                }
            } else if (value instanceof Short) {
                Short shortValue = (Short) value;
                BitSet bitSet = BitSet.valueOf(new long[]{shortValue});
                for (int i = 0; i < 16; i++) {
                    booleanValues.add(bitSet.get(i));
                }
            } else if (value instanceof Integer) {
                Integer integerValue = (Integer) value;
                BitSet bitSet = BitSet.valueOf(new long[]{integerValue});
                for (int i = 0; i < 32; i++) {
                    booleanValues.add(bitSet.get(i));
                }
            } else if (value instanceof Long) {
                long longValue = (Long) value;
                BitSet bitSet = BitSet.valueOf(new long[]{longValue});
                for (int i = 0; i < 64; i++) {
                    booleanValues.add(bitSet.get(i));
                }
            } else {
                throw new IllegalArgumentException(
                    "Value of type " + value.getClass().getName() +
                        " is not assignable to " + modbusField + " fields.");
            }
        }
        return new DefaultBooleanFieldItem(booleanValues.toArray(new Boolean[0]));
    }

    @Override
    public BaseDefaultFieldItem encodeByteArray(PlcField field, Object[] values) {
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();
        for (Object value : values) {
            if (value instanceof byte[]) {
                byte[] byteArray = (byte[]) value;
                byteArrays.add(ArrayUtils.toObject(byteArray));
            } else if (value instanceof Byte[]) {
                Byte[] byteArray = (Byte[]) value;
                byteArrays.add(byteArray);
            } else {
                throw new IllegalArgumentException(
                    "Value of type " + value.getClass().getName() +
                        " is not assignable to " + modbusField + " fields.");
            }
        }
        return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
    }

    @Override
    public BaseDefaultFieldItem encodeShort(PlcField field, Object[] values) {
        int size = Short.BYTES;   
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {  
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof Short){
                    byteBuf.writeShort((Short) values[i]);
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }

    @Override
    public BaseDefaultFieldItem encodeLong(PlcField field, Object[] values) {
        int size = Long.BYTES;          
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {  
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof Long){
                    byteBuf.writeLong((Long) values[i]);
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }

    @Override
    public BaseDefaultFieldItem encodeDouble(PlcField field, Object[] values) {
        int size = Double.BYTES;          
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {  
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof Double){
                    byteBuf.writeDouble((Double) values[i]);
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }

    @Override
    public BaseDefaultFieldItem encodeInteger(PlcField field, Object[] values) {
        int size = Integer.BYTES;         
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {  
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof Integer){
                    byteBuf.writeInt((Integer) values[i]);
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }
    
    @Override
    public BaseDefaultFieldItem encodeFloat(PlcField field, Object[] values) {
        int size = Float.BYTES;        
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {  
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof Float){
                    byteBuf.writeFloat((Float) values[i]);
                }
            }
            
            BytesToRegisters(backend, byteArrays);
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }

    @Override
    public BaseDefaultFieldItem encodeDuration(PlcField field, Object[] values) {
        int size = Long.BYTES;
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {  
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof Duration){
                    byteBuf.writeLong(((Duration) values[i]).toNanos());
                }
            }
            
            BytesToRegisters(backend, byteArrays);
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }

    @Override
    public BaseDefaultFieldItem encodeTime(PlcField field, Object[] values) {
        int size = Long.BYTES;        
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {  
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof LocalTime){
                    byteBuf.writeLong(((LocalTime) values[i]).toNanoOfDay());
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }
    
    @Override
    public BaseDefaultFieldItem encodeDate(PlcField field, Object[] values) {
        int size = Long.BYTES;        
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {   
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof LocalDate){
                    byteBuf.writeLong(((LocalDate) values[i]).toEpochDay());
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }

    @Override
    public BaseDefaultFieldItem encodeDateTime(PlcField field, Object[] values) {
        int size = Long.BYTES;        
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {   
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof LocalDateTime){
                    byteBuf.writeLong(((LocalDateTime) values[i]).toEpochSecond(ZoneOffset.UTC));
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }

    @Override
    public BaseDefaultFieldItem encodeBigDecimal(PlcField field, Object[] values) {
        int size = 4; //BigInteger
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {   
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof BigDecimal){
                    byteBuf.writeBytes(((BigDecimal) values[i]).toBigInteger().toByteArray());
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }

    @Override
    public BaseDefaultFieldItem encodeBigInteger(PlcField field, Object[] values) {
        int size = 4; //BigInteger
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();        
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();
        if ((modbusField.getQuantity()*2/size) == values.length) {   
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof BigInteger){
                    byteBuf.writeBytes(((BigInteger) values[i]).toByteArray());
                }
            }
            
            BytesToRegisters(backend, byteArrays);     
            
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }
    
    @Override
    public BaseDefaultFieldItem encodeString(PlcField field, Object[] values) {
        int size = 0; //BigInteger
        ModbusField modbusField = (ModbusField) field;
        List<Byte[]> byteArrays = new LinkedList<>();  
        for (Object value:values) size+=((String) value).getBytes().length;
        size+=((size % 2) == 0)?0:1;
        byte[] backend = new byte[values.length * size];
        ByteBuf byteBuf = Unpooled.wrappedBuffer(backend);
        byteBuf.clear();

        if ((modbusField.getQuantity()*2) == values.length) {   
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof String){
                    byteBuf.writeBytes(((String) values[i]).getBytes());
                }
            }
            if(byteBuf.readableBytes()<size) byteBuf.writeByte(0);
            BytesToRegisters(backend, byteArrays);
            return new DefaultModbusByteArrayFieldItem(byteArrays.toArray(new Byte[0][0]));
        } else {
            throw new IllegalArgumentException(
                "Value of type " + values.getClass().getName() +
                    " is not assignable to " + modbusField + " fields.");                    
        }
    }
    
    
    private static void BytesToRegisters(byte[] backend, List<Byte[]> registers){
        Byte[] tempBytes = null;        
        int quantity = (backend.length) / 2;
        for (int j=0; j < quantity; j++){
            tempBytes = new Byte[2];
            tempBytes[0] = backend[(j*2)];
            tempBytes[1] = backend[((j*2+1))];
            registers.add(tempBytes);
        }        
    }
   
    
    
    
}
