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
import org.apache.plc4x.java.s7.netty.model.payloads.items.CpuDiagnosticMessageItem;

/**
 *
 * @author cgarcia
 */
public class S7SysEvent implements S7Event{
 
    public static final String EVENT_ID = "EVENT_ID";
    public static final String PRIORITY_CLASS = "PRIORITY_CLASS";
    public static final String OB_NUMBER = "OB_NUMBER";
    public static final String DAT_ID = "DAT_ID";
    public static final String INFO1 = "INFO1";
    public static final String INFO2 = "INFO2";
    
    private final Instant timeStamp;
    private Map<String, Object> map = new HashMap();

    public S7SysEvent(CpuDiagnosticMessageItem item) {
        map.put(EVENT_ID, item.getEventID());
        map.put(PRIORITY_CLASS, item.getPriorityClass());
        map.put(OB_NUMBER, item.getObNumber());
        map.put(DAT_ID, item.getDatID());
        map.put(INFO1, item.getInfo1());
        map.put(INFO2, item.getInfo2());
        this.timeStamp = item.getTimeStamp().toInstant(ZoneOffset.UTC);
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
    public int getNumberOfValues(String name) {
        throw new UnsupportedOperationException("Not supported yet."); 
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
        Object object = map.get(name);
        if (object instanceof Byte){
            return (Byte) object;
        }
        return null;
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
        Object object = map.get(name);
        if (object instanceof Short){
            return (Short) object;
        }
        return null;
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
        Object object = map.get(name);
        if (object instanceof Integer){
            return (Integer) object;
        }
        return null; 
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
        throw new UnsupportedOperationException("Not supported yet."); 
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
        return map.keySet();
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
        return "S7SysEvent{" + "timeStamp=" + timeStamp + ", map=" + map + '}';
    }
    
    
    
}
