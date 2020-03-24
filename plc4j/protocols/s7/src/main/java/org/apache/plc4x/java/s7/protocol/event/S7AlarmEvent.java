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

package org.apache.plc4x.java.s7.protocol.event;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Duration;
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
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.netty.model.payloads.AlarmMessagePayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.AlarmMessageItem;
import org.apache.plc4x.java.s7.netty.model.payloads.items.AssociatedValueItem;
import org.apache.plc4x.java.s7.netty.model.payloads.items.MessageObjectItem;
import org.apache.plc4x.java.s7.netty.model.types.CpuServicesParameterSubFunctionGroup;

/**
 *
 * @author cgarcia
 */
public class S7AlarmEvent implements S7Event {
    
    public enum Fields{
        MAP,
        
        TYPE,
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
        
    }

    //private final Instant timeStamp;
    private final Map<String, Object> map; 

    public S7AlarmEvent(Map<String, Object> map) {
        this.map = map;                      
    }
    
    public static List<S7Event> getAlarmsEvents(CpuServicesParameterSubFunctionGroup subFunction, 
                                        AlarmMessagePayload payload){
        List<S7Event> alarmEvents = new ArrayList();
        Map<String, Object> map = new HashMap(); 
        Instant timeStamp;
        int index = 0;
        AssociatedValueItem itemValue;
        byte[] datos;

        AlarmMessageItem msgItem = payload.getMsg();
        List<Object> objects = msgItem.getMsgItems();
        for (Object object:objects){
            MessageObjectItem item = (MessageObjectItem) object;
            alarmEvents.add(createAlarmEvent(subFunction, msgItem, item));
        }

        return alarmEvents;
    }

    private static S7AlarmEvent createAlarmEvent(CpuServicesParameterSubFunctionGroup subFunction,
                                    AlarmMessageItem msgItem,
                                    MessageObjectItem message){
        Map<String, Object> map = new HashMap();
        AssociatedValueItem itemValue;
        byte[] datos;
        
        switch(subFunction){
            case ALARM_QUERY:{
                switch(message.getAlarmtype()){
                    case ALARM_S:
                        map.put(Fields.TYPE.toString(), CpuServicesParameterSubFunctionGroup.ALARM_S.getCode()); 
                        map.put(Fields.EVENT_ID.name(), message.getEventID());
                        map.put(Fields.EVENT_STATE.name(), message.getEventState());
                        map.put(Fields.ACKSTATE_GOING.name(), message.getAckStateGoing());
                        map.put(Fields.ACKSTATE_COMING.name(), message.getAckStateComming()); 
                        
                        map.put(Fields.TIMESTAMP_GOING.name(), message.getTimestampGoing().toInstant(ZoneOffset.UTC));
                        
                        AssociatedValueItem itemg = message.getComingValues().get(0);
                        datos = new byte[itemg.getLength()];
                        if (itemg.getLength() != 0) {
                            itemg.getData().getBytes(0, datos);
                        };
                        map.put(Fields.SIG_1_DATA_GOING.name(), datos);                        
                        
                        map.put(Fields.TIMESTAMP_COMING.name(), message.getTimestampComing().toInstant(ZoneOffset.UTC));
                        
                        AssociatedValueItem itemc = message.getComingValues().get(0);
                        datos = new byte[itemc.getLength()];                        
                        if (itemc.getLength() != 0) {
                            itemc.getData().getBytes(0, datos);                        
                        }
                        map.put(Fields.SIG_1_DATA_COMING.name(), datos);                        
                        break;
                        
                    case ALARM_8:                        
                        map.put(Fields.TYPE.toString(), CpuServicesParameterSubFunctionGroup.ALARM8);
                        break;
                    case SCAN:
                        map.put(Fields.TYPE.toString(), CpuServicesParameterSubFunctionGroup.SCAN);
                        break;
                    case NONE:
                        map.put(Fields.TYPE.toString(), CpuServicesParameterSubFunctionGroup.NONE);
                        break;
                    default:;
                }                
            }
            break;
            case ALARM_ACK_IND:;
                map.put(Fields.TYPE.toString(), CpuServicesParameterSubFunctionGroup.ALARM_S_IND.getCode()); 
                map.put(Fields.TIMESTAMP.name(), msgItem.getTimestamp().toInstant(ZoneOffset.UTC));
                map.put(Fields.EVENT_ID.name(), message.getEventID()); 
                map.put(Fields.EVENT_STATE.name(), message.getEventState());
                map.put(Fields.STATE.name(), message.getState());
                map.put(Fields.ACKSTATE_GOING.name(), message.getAckStateGoing());
                map.put(Fields.ACKSTATE_COMING.name(), message.getAckStateComming());                 
            break;
            case ALARM_SQ_IND:             
            case ALARM_S_IND:{
                map.put(Fields.TYPE.toString(), CpuServicesParameterSubFunctionGroup.ALARM_S_IND.getCode()); 
                map.put(Fields.TIMESTAMP.name(), msgItem.getTimestamp().toInstant(ZoneOffset.UTC));
                map.put(Fields.EVENT_ID.name(), message.getEventID());
                map.put(Fields.EVENT_STATE.name(), message.getEventState());
                map.put(Fields.STATE.name(), message.getState());
                map.put(Fields.ACKSTATE_GOING.name(), message.getAckStateGoing());
                map.put(Fields.ACKSTATE_COMING.name(), message.getAckStateComming()); 
                AssociatedValueItem itemc = message.getComingValues().get(0);
                datos = new byte[itemc.getLength()];
                if (itemc.getLength() != 0) {
                    itemc.getData().getBytes(0, datos);                        
                }
                map.put(Fields.SIG_1_DATA.name(), datos);                 
            };
            break;
            case NOTIFY8:{
                map.put(Fields.TYPE.toString(), CpuServicesParameterSubFunctionGroup.NOTIFY8.getCode());
                map.put(Fields.TIMESTAMP.name(), msgItem.getTimestamp().toInstant(ZoneOffset.UTC));
                map.put(Fields.EVENT_ID.name(), message.getEventID());
                map.put(Fields.EVENT_STATE.name(), message.getEventState());
                map.put(Fields.STATE.name(), message.getState());
                map.put(Fields.ACKSTATE_GOING.name(), message.getAckStateGoing());
                map.put(Fields.ACKSTATE_COMING.name(), message.getAckStateComming());
                map.put(Fields.EVENT_GOING.name(), message.getEventGoing());
                map.put(Fields.EVENT_COMING.name(), message.getEventComming());                
                map.put(Fields.EVENT_LAST_CHANGE.name(), message.getEventLastChange());  
                for(int i = 0; i < message.getNumberOfValues(); i++){                
                    itemValue = message.getComingValues().get(i);
                    datos = new byte[itemValue.getLength()];
                    map.put("SIG_"+ (i+1) +"_STATE", itemValue.getReturnCode().getCode());
                    itemValue.getData().getBytes(0, datos);
                    map.put("SIG_"+ (i+1) +"_DATA", datos);
                }                  
            }
            break;
        }
                
        return new S7AlarmEvent(map);
    }
    
    
    
    @Override
    public Instant getTimestamp() {
        return (Instant) map.get(Fields.TIMESTAMP.name());
    }

    @Override
    public PlcReadRequest getRequest() {
        throw new UnsupportedOperationException("Not supported yet.");
    }
        
    @Override
    public int getNumberOfValues(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Object getObject(String name) {
        switch(Fields.valueOf(name)){
            case MAP: return map;
            default:;
        }
        return null;
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
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Long getLong(String name, int index) {
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
    public boolean isValidByteArray(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidByteArray(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Byte[] getByteArray(String name) {
        return (Byte[]) map.get(name);
    }

    @Override
    public Byte[] getByteArray(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Byte[]> getAllByteArrays(String name) {
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
    public boolean isValidDuration(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public boolean isValidDuration(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Duration getDuration(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Duration getDuration(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Duration> getAllDuration(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    
    
    
    @Override
    public Map<String, Object> getMap() {
        return map;
    }
    
        
    
}
