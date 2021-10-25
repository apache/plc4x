/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.events;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.Charset;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneOffset;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesPush;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesSubscribeResponse;

/**
 *
 * @author cgarcia
 */
public class S7CyclicEvent implements S7Event {

    public enum Fields{
        TYPE,
        JOBID,
        TIMESTAMP,        
        ITEMSCOUNT,
        REQUEST,
        MAP,
        RETURNCODE_,
        TRANSPORTSIZE_,
        DATA_
    }    
    
    private final PlcSubscriptionRequest request;
    
    private final Instant timeStamp;
    private final Map<String, Object> map;    
    
    private int j;
    
    public S7CyclicEvent(PlcSubscriptionRequest request, short jobid, S7PayloadUserDataItemCyclicServicesPush event) {
        this.map = new HashMap();
        this.timeStamp = Instant.now(); 
        this.request = request;
        map.put(Fields.TYPE.name(), "CYCEVENT");         
        map.put(Fields.TIMESTAMP.name(),this.timeStamp);
        map.put(Fields.JOBID.name(), jobid);
        map.put(Fields.ITEMSCOUNT.name(), event.getItemsCount());        
        for (int i=0; i<event.getItemsCount(); i++){
            map.put(Fields.RETURNCODE_.name()+i, event.getItems()[i].getReturnCode().getValue());
            map.put(Fields.TRANSPORTSIZE_.name()+i, event.getItems()[i].getTransportSize().getValue());
            byte[] buffer = new byte[event.getItems()[i].getData().length];
            j = 0;
            for(short s:event.getItems()[i].getData()){
                buffer[j] = (byte) s;
                j ++;
            }
            map.put(Fields.DATA_.name()+i, buffer);  
        }
    }
    
    public S7CyclicEvent(PlcSubscriptionRequest request, short jobid, S7PayloadUserDataItemCyclicServicesSubscribeResponse event) {
        this.map = new HashMap();
        this.timeStamp = Instant.now(); 
        this.request = request;
        map.put(Fields.TYPE.name(), "CYCEVENT");         
        map.put(Fields.TIMESTAMP.name(),this.timeStamp);
        map.put(Fields.JOBID.name(), jobid);
        map.put(Fields.ITEMSCOUNT.name(), event.getItemsCount());
        for (int i=0; i<event.getItemsCount(); i++){
            map.put(Fields.RETURNCODE_.name()+i, event.getItems()[i].getReturnCode().getValue());
            map.put(Fields.TRANSPORTSIZE_.name()+i, event.getItems()[i].getTransportSize().getValue());
            byte[] buffer = new byte[event.getItems()[i].getData().length];
            j = 0;
            for(short s:event.getItems()[i].getData()){
                buffer[j] = (byte) s;
                j ++;
            }
            map.put(Fields.DATA_.name()+i, buffer); 
        }            
    }
        
    @Override
    public Map<String, Object> getMap() {
        return this.map;
    }

    @Override
    public Instant getTimestamp() {
        return this.timeStamp;
    }

    @Override
    public PlcReadRequest getRequest() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcValue getAsPlcValue() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcValue getPlcValue(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public int getNumberOfValues(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Object getObject(String name) {
        if ("REQUEST".equals(name)) return request;
        return null;
    }

    @Override
    public Object getObject(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidBoolean(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Boolean getBoolean(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidByte(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidByte(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Byte getByte(String name) {
       if (!(map.get(name) instanceof Byte)) 
            throw new UnsupportedOperationException("Field is not a Byte. Required Byte type.");
        return (byte) map.get(name);
    }

    @Override
    public Byte getByte(String name, int index) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");     
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index*Byte.BYTES;
        return byteBuf.getByte(pos);
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        byte[] array = (byte[]) map.get(name);
        
        List<Byte> list = IntStream.range(0, array.length).
                mapToObj(i -> array[i]).collect(Collectors.toList());
        
        return list;
    }

    @Override
    public boolean isValidShort(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidShort(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Short getShort(String name) {
        if (!(map.get(name) instanceof Short)) return null;
        return (short) map.get(name);
    }

    @Override
    public Short getShort(String name, int index) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");    
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index*Short.BYTES;
        return byteBuf.getShort(index);
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Short> list = new ArrayList();
        while(byteBuf.isReadable(Short.BYTES)) list.add(byteBuf.readShort());        
        return list;
    }

    @Override
    public boolean isValidInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Integer getInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Integer getInteger(String name, int index) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");     
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index*Integer.BYTES;
        return byteBuf.getInt(pos);
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");        
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Integer> list = new ArrayList();
        while(byteBuf.isReadable(Integer.BYTES)) list.add(byteBuf.readInt());        
        return list;
    }

    @Override
    public boolean isValidBigInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public BigInteger getBigInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<BigInteger> getAllBigIntegers(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidLong(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidLong(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Long getLong(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Long getLong(String name, int index) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");     
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index*Long.BYTES;
        return byteBuf.getLong(pos);
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");     
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Long> list = new ArrayList();
        while(byteBuf.isReadable(Long.BYTES)) list.add(byteBuf.readLong());        
        return list;
    }

    @Override
    public boolean isValidFloat(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Float getFloat(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Float getFloat(String name, int index) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");     
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index*Float.BYTES;
        return byteBuf.getFloat(pos);
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");       
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Float> list = new ArrayList();
        while(byteBuf.isReadable(Float.BYTES)) list.add(byteBuf.readFloat());        
        return list;
    }

    @Override
    public boolean isValidDouble(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Double getDouble(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Double getDouble(String name, int index) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");    
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index*Double.BYTES;
        return byteBuf.getDouble(pos);
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
       if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");      
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Double> list = new ArrayList();
        while(byteBuf.isReadable(Double.BYTES)) list.add(byteBuf.readDouble());        
        return list;
    }

    @Override
    public boolean isValidBigDecimal(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidString(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidString(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public String getString(String name) {
        if (!(map.get(name) instanceof byte[])) 
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        return byteBuf.toString(Charset.defaultCharset());
    }

    @Override
    public String getString(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidTime(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public LocalTime getTime(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public LocalTime getTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidDate(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidDate(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public LocalDate getDate(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public LocalDate getDate(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidDateTime(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<String> getFieldNames() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcField getField(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }
    
}
