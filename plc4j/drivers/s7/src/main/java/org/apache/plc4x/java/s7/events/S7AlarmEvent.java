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
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.Charset;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneOffset;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageAckObjectPushType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageAckPushType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageObjectPushType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageObjectQueryType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessagePushType;
import org.apache.plc4x.java.s7.readwrite.AssociatedQueryValueType;
import org.apache.plc4x.java.s7.readwrite.AssociatedValueType;
import org.apache.plc4x.java.s7.readwrite.DateAndTime;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarm8;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmAckInd;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmS;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmSQ;
import org.apache.plc4x.java.s7.readwrite.S7PayloadNotify;
import org.apache.plc4x.java.s7.readwrite.S7PayloadNotify8;
import org.apache.plc4x.java.s7.readwrite.types.AlarmType;

/**
 *
 * @author cgarcia
 */
public class S7AlarmEvent implements S7Event {

    
    public enum Fields{
        MAP,
        
        TYPE,
        ASSOCIATED_VALUES,
        TIMESTAMP,
        TIMESTAMP_GOING,
        TIMESTAMP_COMING,        
        
        EVENT_ID,
        EVENT_STATE,
        STATE,
        ACKSTATE_GOING,
        ACKSTATE_COMING,
        EVENT_GOING,
        EVENT_COMING,
        EVENT_LAST_CHANGE,
    
        SIG,
        SIG_1,
        SIG_2,
        SIG_3,
        SIG_4,    
        SIG_5,
        SIG_6,
        SIG_7,
        SIG_8,

        SIG_STATE,
        SIG_1_STATE,
        SIG_2_STATE,
        SIG_3_STATE,
        SIG_4_STATE,    
        SIG_5_STATE,
        SIG_6_STATE,
        SIG_7_STATE,
        SIG_8_STATE,    

        SIG_DATA,
        SIG_DATA_GOING,        
        SIG_DATA_GOING_STATUS,    
        SIG_DATA_GOING_SIZE,  
        SIG_DATA_GOING_LENGTH,         
        
        SIG_DATA_COMING,
        SIG_DATA_COMING_STATUS,    
        SIG_DATA_COMING_SIZE,  
        SIG_DATA_COMING_LENGTH,          
        
        SIG_1_DATA,
        SIG_2_DATA,
        SIG_3_DATA,
        SIG_4_DATA,    
        SIG_5_DATA,
        SIG_6_DATA,
        SIG_7_DATA,
        SIG_8_DATA,
        
        SIG_1_DATA_GOING,
        SIG_2_DATA_GOING,
        SIG_3_DATA_GOING,
        SIG_4_DATA_GOING,    
        SIG_5_DATA_GOING,
        SIG_6_DATA_GOING,
        SIG_7_DATA_GOING,
        SIG_8_DATA_GOING,
        SIG_1_DATA_COMING,
        SIG_2_DATA_COMING,
        SIG_3_DATA_COMING,
        SIG_4_DATA_COMING,    
        SIG_5_DATA_COMING,
        SIG_6_DATA_COMING,
        SIG_7_DATA_COMING,
        SIG_8_DATA_COMING,
        
        SIG_1_DATA_STATUS,
        SIG_2_DATA_STATUS,
        SIG_3_DATA_STATUS,
        SIG_4_DATA_STATUS,    
        SIG_5_DATA_STATUS,
        SIG_6_DATA_STATUS,
        SIG_7_DATA_STATUS,
        SIG_8_DATA_STATUS,   
        
        SIG_1_DATA_SIZE,
        SIG_2_DATA_SIZE,
        SIG_3_DATA_SIZE,
        SIG_4_DATA_SIZE,    
        SIG_5_DATA_SIZE,
        SIG_6_DATA_SIZE,
        SIG_7_DATA_SIZE,
        SIG_8_DATA_SIZE,    
        
        SIG_1_DATA_LENGTH,
        SIG_2_DATA_LENGTH,
        SIG_3_DATA_LENGTH,
        SIG_4_DATA_LENGTH,    
        SIG_5_DATA_LENGTH,
        SIG_6_DATA_LENGTH,
        SIG_7_DATA_LENGTH,
        SIG_8_DATA_LENGTH,            
        
    };
  
    private Instant timeStamp;
    private final Map<String, Object> map;   

    public S7AlarmEvent(Object obj) {
        this.map = new HashMap();
        
        if ((obj instanceof AlarmMessageObjectQueryType)) {
            AlarmMessageObjectQueryType event = (AlarmMessageObjectQueryType) obj;
            map.put(Fields.EVENT_ID.name(), event.getEventId());
            
            if (event.getAlarmType() == AlarmType.ALARM_S) {
                map.put(Fields.TYPE.name(), "ALARMS_QUERY");    
                map.put(Fields.ASSOCIATED_VALUES.name(), 1); 
                
                DateAndTime tcoming = event.getTimeComing();
                int year = (tcoming.getYear()>=90)?tcoming.getYear()+1900:tcoming.getYear()+2000;        
                 LocalDateTime ldt = LocalDateTime.of(year,
                        tcoming.getMonth(),
                        tcoming.getDay(),
                        tcoming.getHour(),
                        tcoming.getMinutes(), 
                        tcoming.getSeconds(), 
                        tcoming.getMsec()*1000000);
                this.timeStamp = ldt.toInstant(ZoneOffset.UTC);
                map.put(S7AlarmEvent.Fields.TIMESTAMP_COMING.name(),this.timeStamp);  
                
                DateAndTime tgoing = event.getTimeGoing();
                /*
                year = (tgoing.getYear()>=90)?tgoing.getYear()+1900:tgoing.getYear()+2000;        
                ldt = LocalDateTime.of(year,
                        tgoing.getMonth(),
                        tgoing.getDay(),
                        tgoing.getHour(),
                        tgoing.getMinutes(), 
                        tgoing.getSeconds(), 
                        tgoing.getMsec()*1000000);
                Instant timegoing = ldt.toInstant(ZoneOffset.UTC);
                */
                map.put(S7AlarmEvent.Fields.TIMESTAMP_GOING.name(),Instant.now());  

                AssociatedQueryValueType datacoming = event.getValueComing();
                AssociatedQueryValueType datagoing = event.getValueGoing();
                
                byte[] buffercoming = new byte[datacoming.getData().length];
                byte[] buffergoing  = new byte[datagoing.getData().length];
                                
                int j = 0;
                for (short s:datacoming.getData()) {
                    buffercoming[j] = (byte) s;
                    j++;
                }
                map.put(Fields.SIG_DATA_COMING.name(),  buffercoming);
                map.put(Fields.SIG_DATA_COMING_STATUS.name(), datacoming.getReturnCode().getValue());
                map.put(Fields.SIG_DATA_COMING_SIZE.name(), datacoming.getTransportSize().getValue());
                map.put(Fields.SIG_DATA_COMING_LENGTH.name(), (short) datacoming.getValueLength());                
                
                j = 0;
                for (short s:datagoing.getData()) {
                    buffergoing[j] = (byte) s;
                    j++;
                }
                    
                map.put(Fields.SIG_DATA_GOING.name(),  buffergoing);                   
                map.put(Fields.SIG_DATA_GOING_STATUS.name(), datagoing.getReturnCode().getValue());
                map.put(Fields.SIG_DATA_GOING_SIZE.name(), datagoing.getTransportSize().getValue());
                map.put(Fields.SIG_DATA_GOING_LENGTH.name(), (short) datagoing.getValueLength());              
                                                           
            } else if (event.getAlarmType() == AlarmType.ALARM_8) {
                map.put(Fields.TYPE.name(), "ALARM8_QUERY");
            };

            map.put(Fields.SIG_1_DATA_GOING.name(), event.getAckStateGoing().getSIG_1());
            map.put(Fields.SIG_2_DATA_GOING.name(), event.getAckStateGoing().getSIG_2());
            map.put(Fields.SIG_3_DATA_GOING.name(), event.getAckStateGoing().getSIG_3());
            map.put(Fields.SIG_4_DATA_GOING.name(), event.getAckStateGoing().getSIG_4());
            map.put(Fields.SIG_5_DATA_GOING.name(), event.getAckStateGoing().getSIG_5());
            map.put(Fields.SIG_6_DATA_GOING.name(), event.getAckStateGoing().getSIG_6());
            map.put(Fields.SIG_7_DATA_GOING.name(), event.getAckStateGoing().getSIG_7());
            map.put(Fields.SIG_8_DATA_GOING.name(), event.getAckStateGoing().getSIG_8()); 

            map.put(Fields.SIG_1_DATA_COMING.name(), event.getAckStateComing().getSIG_1());
            map.put(Fields.SIG_2_DATA_COMING.name(), event.getAckStateComing().getSIG_2());
            map.put(Fields.SIG_3_DATA_COMING.name(), event.getAckStateComing().getSIG_3());
            map.put(Fields.SIG_4_DATA_COMING.name(), event.getAckStateComing().getSIG_4());            
            map.put(Fields.SIG_5_DATA_COMING.name(), event.getAckStateComing().getSIG_5());
            map.put(Fields.SIG_6_DATA_COMING.name(), event.getAckStateComing().getSIG_6());
            map.put(Fields.SIG_7_DATA_COMING.name(), event.getAckStateComing().getSIG_7());
            map.put(Fields.SIG_8_DATA_COMING.name(), event.getAckStateComing().getSIG_8());                   
                
            
        } else if (obj instanceof S7PayloadAlarmAckInd) {
            AlarmMessageAckPushType msg = ((S7PayloadAlarmAckInd) obj).getAlarmMessage(); 
            DateAndTime dt = msg.getTimeStamp();
            int year = (dt.getYear()>=90)?dt.getYear()+1900:dt.getYear()+2000;        
             LocalDateTime ldt = LocalDateTime.of(year,
                    dt.getMonth(),
                    dt.getDay(),
                    dt.getHour(),
                    dt.getMinutes(), 
                    dt.getSeconds(), 
                    dt.getMsec()*1000000);
            this.timeStamp = ldt.toInstant(ZoneOffset.UTC);
            map.put(S7SysEvent.Fields.TIMESTAMP.name(),this.timeStamp);            

            AlarmMessageAckObjectPushType[] items = msg.getMessageObjects();
            for (AlarmMessageAckObjectPushType item:items){
                map.put(Fields.EVENT_ID.name(), item.getEventId()); 
                map.put(Fields.TYPE.name(), "ALARMACK_IND");  
                map.put(Fields.ASSOCIATED_VALUES.name(), item.getNumberOfValues()); 
                
                map.put(Fields.SIG_1_DATA_GOING.name(), item.getAckStateGoing().getSIG_1());
                map.put(Fields.SIG_2_DATA_GOING.name(), item.getAckStateGoing().getSIG_2());
                map.put(Fields.SIG_3_DATA_GOING.name(), item.getAckStateGoing().getSIG_3());
                map.put(Fields.SIG_4_DATA_GOING.name(), item.getAckStateGoing().getSIG_4());
                map.put(Fields.SIG_5_DATA_GOING.name(), item.getAckStateGoing().getSIG_5());
                map.put(Fields.SIG_6_DATA_GOING.name(), item.getAckStateGoing().getSIG_6());
                map.put(Fields.SIG_7_DATA_GOING.name(), item.getAckStateGoing().getSIG_7());
                map.put(Fields.SIG_8_DATA_GOING.name(), item.getAckStateGoing().getSIG_8()); 

                map.put(Fields.SIG_1_DATA_COMING.name(), item.getAckStateComing().getSIG_1());
                map.put(Fields.SIG_2_DATA_COMING.name(), item.getAckStateComing().getSIG_2());
                map.put(Fields.SIG_3_DATA_COMING.name(), item.getAckStateComing().getSIG_3());
                map.put(Fields.SIG_4_DATA_COMING.name(), item.getAckStateComing().getSIG_4());            
                map.put(Fields.SIG_5_DATA_COMING.name(), item.getAckStateComing().getSIG_5());
                map.put(Fields.SIG_6_DATA_COMING.name(), item.getAckStateComing().getSIG_6());
                map.put(Fields.SIG_7_DATA_COMING.name(), item.getAckStateComing().getSIG_7());
                map.put(Fields.SIG_8_DATA_COMING.name(), item.getAckStateComing().getSIG_8());                   
            }            
            
        } else {       

            AlarmMessagePushType msg = null;

            if (obj instanceof S7PayloadAlarm8)
                msg = ((S7PayloadAlarm8) obj).getAlarmMessage();
            if (obj instanceof S7PayloadNotify)
                msg = ((S7PayloadNotify) obj).getAlarmMessage(); 
            if (obj instanceof S7PayloadAlarmSQ)
                msg = ((S7PayloadAlarmSQ) obj).getAlarmMessage(); 
            if (obj instanceof S7PayloadAlarmS)
                msg = ((S7PayloadAlarmS) obj).getAlarmMessage();                             
            if (obj instanceof S7PayloadNotify8)
                msg = ((S7PayloadNotify8) obj).getAlarmMessage();        

            DateAndTime dt = msg.getTimeStamp();
            int year = (dt.getYear()>=90)?dt.getYear()+1900:dt.getYear()+2000;        
             LocalDateTime ldt = LocalDateTime.of(year,
                    dt.getMonth(),
                    dt.getDay(),
                    dt.getHour(),
                    dt.getMinutes(), 
                    dt.getSeconds(), 
                    dt.getMsec()*1000000);
            this.timeStamp = ldt.toInstant(ZoneOffset.UTC);
            map.put(S7SysEvent.Fields.TIMESTAMP.name(),this.timeStamp);

            AlarmMessageObjectPushType[] items = msg.getMessageObjects();
            for (AlarmMessageObjectPushType item:items){
                map.put(Fields.EVENT_ID.name(), item.getEventId());

                if (obj instanceof S7PayloadAlarm8)
                    map.put(Fields.TYPE.name(), "ALARM8");
                if (obj instanceof S7PayloadNotify)
                    map.put(Fields.TYPE.name(), "NOTIFY");
                if (obj instanceof S7PayloadAlarmSQ)
                    map.put(Fields.TYPE.name(), "ALARMSQ");
                if (obj instanceof S7PayloadAlarmS)
                    map.put(Fields.TYPE.name(), "ALARMS");                             
                if (obj instanceof S7PayloadNotify8)
                    map.put(Fields.TYPE.name(), "NOTIFY8");
          

                map.put(Fields.ASSOCIATED_VALUES.name(), item.getNumberOfValues());


                map.put(Fields.SIG_1.name(), item.getEventState().getSIG_1());
                map.put(Fields.SIG_2.name(), item.getEventState().getSIG_2());
                map.put(Fields.SIG_3.name(), item.getEventState().getSIG_3());
                map.put(Fields.SIG_4.name(), item.getEventState().getSIG_4());
                map.put(Fields.SIG_5.name(), item.getEventState().getSIG_5());
                map.put(Fields.SIG_6.name(), item.getEventState().getSIG_6());
                map.put(Fields.SIG_7.name(), item.getEventState().getSIG_7());
                map.put(Fields.SIG_8.name(), item.getEventState().getSIG_8());


                map.put(Fields.SIG_1_STATE.name(), item.getLocalState().getSIG_1());
                map.put(Fields.SIG_2_STATE.name(), item.getLocalState().getSIG_2());
                map.put(Fields.SIG_3_STATE.name(), item.getLocalState().getSIG_3());
                map.put(Fields.SIG_4_STATE.name(), item.getLocalState().getSIG_4()); 
                map.put(Fields.SIG_5_STATE.name(), item.getLocalState().getSIG_5());
                map.put(Fields.SIG_6_STATE.name(), item.getLocalState().getSIG_6());
                map.put(Fields.SIG_7_STATE.name(), item.getLocalState().getSIG_7());
                map.put(Fields.SIG_8_STATE.name(), item.getLocalState().getSIG_8());  

                map.put(Fields.SIG_1_DATA_GOING.name(), item.getAckStateGoing().getSIG_1());
                map.put(Fields.SIG_2_DATA_GOING.name(), item.getAckStateGoing().getSIG_2());
                map.put(Fields.SIG_3_DATA_GOING.name(), item.getAckStateGoing().getSIG_3());
                map.put(Fields.SIG_4_DATA_GOING.name(), item.getAckStateGoing().getSIG_4());
                map.put(Fields.SIG_5_DATA_GOING.name(), item.getAckStateGoing().getSIG_5());
                map.put(Fields.SIG_6_DATA_GOING.name(), item.getAckStateGoing().getSIG_6());
                map.put(Fields.SIG_7_DATA_GOING.name(), item.getAckStateGoing().getSIG_7());
                map.put(Fields.SIG_8_DATA_GOING.name(), item.getAckStateGoing().getSIG_8()); 

                map.put(Fields.SIG_1_DATA_COMING.name(), item.getAckStateComing().getSIG_1());
                map.put(Fields.SIG_2_DATA_COMING.name(), item.getAckStateComing().getSIG_2());
                map.put(Fields.SIG_3_DATA_COMING.name(), item.getAckStateComing().getSIG_3());
                map.put(Fields.SIG_4_DATA_COMING.name(), item.getAckStateComing().getSIG_4());            
                map.put(Fields.SIG_5_DATA_COMING.name(), item.getAckStateComing().getSIG_5());
                map.put(Fields.SIG_6_DATA_COMING.name(), item.getAckStateComing().getSIG_6());
                map.put(Fields.SIG_7_DATA_COMING.name(), item.getAckStateComing().getSIG_7());
                map.put(Fields.SIG_8_DATA_COMING.name(), item.getAckStateComing().getSIG_8());             

                AssociatedValueType[] values = item.getAssociatedValues();
                int i=1;
                int j = 0;
                for (AssociatedValueType value:values) {
                    map.put("SIG_"+i+"_DATA_STATUS", value.getReturnCode().getValue());
                    map.put("SIG_"+i+"_DATA_SIZE", value.getTransportSize().getValue());
                    map.put("SIG_"+i+"_DATA_LENGTH", (short) value.getValueLength());
                    byte[] data = new byte[value.getData().length];
                    j = 0;
                    for (short s:value.getData()) {
                        data[j] = (byte) s;
                        j++;
                    }
                    map.put("SIG_"+i+"_DATA",  data);
                    i++;
                }
            
            }
            
        } 
    };
    
    
    
    @Override
    public Map<String, Object> getMap() {
        return map;
    }

    @Override
    public Instant getTimestamp() {
        return timeStamp;
    }

    @Override
    public PlcReadRequest getRequest() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PlcValue getAsPlcValue() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PlcValue getPlcValue(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public int getNumberOfValues(String name) {
        return (int) map.get(Fields.ASSOCIATED_VALUES.name());
    }

    @Override
    public Object getObject(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Object getObject(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBoolean(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Boolean getBoolean(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidByte(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidByte(String name, int index) {  
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Byte getByte(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Byte getByte(String name, int index) {
        if (map.get(name) instanceof byte[]) {
            ByteBuf bytebuf = Unpooled.wrappedBuffer((byte[])map.get(name));
            if (bytebuf.readableBytes() >= index * Byte.BYTES) {
                return bytebuf.getByte(index);
            }
        }; 
        throw new UnsupportedOperationException("Not supported yet.");        
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidShort(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidShort(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Short getShort(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Short getShort(String name, int index) {
        if (map.get(name) instanceof byte[]) {
            ByteBuf bytebuf = Unpooled.wrappedBuffer((byte[])map.get(name));
            if (bytebuf.readableBytes() >= index * Short.BYTES) {
                return bytebuf.getShort(index);
            }
        };       
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Integer getInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Integer getInteger(String name, int index) {
        if (map.get(name) instanceof byte[]) {
            ByteBuf bytebuf = Unpooled.wrappedBuffer((byte[])map.get(name));
            if (bytebuf.readableBytes() >= index * Integer.BYTES) {
                return bytebuf.getInt(index);
            }
        }; 
        throw new UnsupportedOperationException("Not supported yet.");        
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBigInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public BigInteger getBigInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<BigInteger> getAllBigIntegers(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidLong(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidLong(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Long getLong(String name) {
        return (long) map.get(name);
    }

    @Override
    public Long getLong(String name, int index) {
        if (map.get(name) instanceof byte[]) {
            ByteBuf bytebuf = Unpooled.wrappedBuffer((byte[])map.get(name));
            if (bytebuf.readableBytes() >= index * Long.BYTES) {
                return bytebuf.getLong(index);
            }
        }; 
        throw new UnsupportedOperationException("Not supported yet.");   
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidFloat(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Float getFloat(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Float getFloat(String name, int index) {
        if (map.get(name) instanceof byte[]) {
            ByteBuf bytebuf = Unpooled.wrappedBuffer((byte[])map.get(name));
            if (bytebuf.readableBytes() >= index * Long.BYTES) {
                return bytebuf.getFloat(index);
            }
        }; 
        throw new UnsupportedOperationException("Not supported yet.");   
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDouble(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Double getDouble(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Double getDouble(String name, int index) {
        if (map.get(name) instanceof byte[]) {
            ByteBuf bytebuf = Unpooled.wrappedBuffer((byte[])map.get(name));
            if (bytebuf.readableBytes() >= index * Double.BYTES) {
                return bytebuf.getDouble(index);
            }
        }; 
        throw new UnsupportedOperationException("Not supported yet."); 
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBigDecimal(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidString(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidString(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public String getString(String name) {
        if (map.get(name) instanceof byte[]) {
            ByteBuf bytebuf = Unpooled.wrappedBuffer((byte[])map.get(name));
             return bytebuf.readCharSequence(bytebuf.readableBytes(), Charset.forName("utf-8")).toString();            
        }; 
        throw new UnsupportedOperationException("Not supported yet."); 
    }

    @Override
    public String getString(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidTime(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalTime getTime(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalTime getTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDate(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDate(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalDate getDate(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalDate getDate(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDateTime(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<String> getFieldNames() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PlcField getField(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public String toString() {
        ByteBuf buffer;
        StringBuilder str = new StringBuilder();
        str.append(Fields.EVENT_ID.name() + ": " + Long.toHexString((long) map.get(Fields.EVENT_ID.name())));
        str.append("\r\n");        
        str.append(Fields.TYPE.name() + ": " + ((String) map.get(Fields.TYPE.name())));
        str.append("\r\n");
        short nvalues = (short) map.get(Fields.ASSOCIATED_VALUES.name());
        str.append(Fields.ASSOCIATED_VALUES.name() + ": " + nvalues); 
        str.append("\r\n");
        str.append(Fields.TIMESTAMP.name() + ": " + ((Instant) map.get(Fields.TIMESTAMP.name()))); 
        str.append("\r\n");
        
        str.append("SIG_X: ");
        for (int i=8; i>=1; i--) {
            str.append(map.get("SIG_"+i) + " ");
        }        
        str.append("\r\n"); 
        str.append("SIG_X_DATA_GOING: ");
        for (int i=8; i>=1; i--) {
            str.append(map.get("SIG_"+i+"_DATA_GOING") + " ");
        }  
        str.append("\r\n"); 
        str.append("SIG_X_DATA_COMING: ");
        for (int i=8; i>=1; i--) {
            str.append(map.get("SIG_"+i+"_DATA_COMING") + " ");
        }   
        str.append("\r\n");        
        for (int i=1; i<= nvalues; i++){
            str.append("SIG_"+i+"_DATA: " + ((short) map.get("SIG_"+i+"_DATA_STATUS")));
            str.append(" " + ((short) map.get("SIG_"+i+"_DATA_SIZE")));
            str.append(" " + ((short) map.get("SIG_"+i+"_DATA_LENGTH")));
            buffer = Unpooled.wrappedBuffer((byte[]) map.get("SIG_"+i+"_DATA"));
            str.append(" "  + ByteBufUtil.hexDump(buffer));   
            str.append("\r\n");            
        }
        
        
        return str.toString();

    }
    


    
}
