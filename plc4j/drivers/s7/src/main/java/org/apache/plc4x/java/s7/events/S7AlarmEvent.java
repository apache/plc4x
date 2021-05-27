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

package org.apache.plc4x.java.s7.events;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageObjectPushType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessagePushType;
import org.apache.plc4x.java.s7.readwrite.AssociatedValueType;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarm8;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmS;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmSC;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmSQ;
import org.apache.plc4x.java.s7.readwrite.S7PayloadNotify;
import org.apache.plc4x.java.s7.readwrite.S7PayloadNotify8;

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
  
    private final Instant timeStamp;
    private Map<String, Object> map = new HashMap();   

    public S7AlarmEvent(Object obj) {
             
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
        
        
        this.timeStamp = null;
        if (msg == null) return;
        
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
            for (AssociatedValueType value:values) {
                map.put("SIG_"+i+"_DATA_STATUS", value.getReturnCode().getValue());
                map.put("SIG_"+i+"_DATA_SIZE", value.getTransportSize().getValue());
                map.put("SIG_"+i+"_DATA_LENGTH", value.getValueLength());
                map.put("SIG_"+i+"_DATA", value.getData());
                i=+1;
            }
            
        }

    }
    
    
    
    @Override
    public Map<String, Object> getMap() {
        return map;
    }

    @Override
    public Instant getTimestamp() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Byte getByte(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Short getShort(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
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
