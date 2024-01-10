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
package org.apache.plc4x.merlot.api.impl;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.Set;
import static java.util.stream.Collectors.toList;
import static java.util.stream.Collectors.toMap;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.osgi.service.dal.FunctionData;
import static org.osgi.service.dal.FunctionData.DESCRIPTION;
import static org.osgi.service.dal.FunctionData.FIELD_TIMESTAMP;

/*
*
*/
public class PlcValueImpl extends FunctionData  implements PlcValue {

    public static final String	PV	= "pv";    
       
    private PlcValueType value_type;
    
    public PlcValueImpl(PlcValueBuilder builder){
        super(builder.timestamp, builder.metadata);
        this.value_type = builder.value_type;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return value_type; 
    }

    @Override
    public Object getObject() {
        return getMetadata().get(PV);
    }

    @Override
    public boolean isSimple() {
        return !getMetadata().get(PV).getClass().isArray();
    }

    @Override
    public boolean isNullable() {
        return true;
    }

    @Override
    public boolean isNull() {
        return (null == getMetadata().get(PV));
    }

    @Override
    public boolean is(Class<?> clazz) {
        return getMetadata().get(PV).getClass().isAssignableFrom(clazz);
    }

    @Override
    public boolean isConvertibleTo(Class<?> clazz) {
        return getMetadata().get(PV).getClass().isAssignableFrom(clazz);
    }

    @Override
    public <T> T get(Class<T> clazz) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isBoolean() {
        return (getMetadata().get(PV) instanceof Boolean);
    }

    @Override
    public boolean getBoolean() {
        return (Boolean) getMetadata().get(PV);
    }

    @Override
    public boolean isByte() {
        return (getMetadata().get(PV) instanceof Byte);
    }

    @Override
    public byte getByte() {
        return (Byte) getMetadata().get(PV);
    }

    @Override
    public boolean isShort() {
        return (getMetadata().get(PV) instanceof Short);
    }

    @Override
    public short getShort() {
        return (Short) getMetadata().get(PV);
    }

    @Override
    public boolean isInteger() {
        return (getMetadata().get(PV) instanceof Integer);
    }

    @Override
    public int getInteger() {
        return (Integer) getMetadata().get(PV);
    }

    @Override
    public int getInt() {
        return (int) getMetadata().get(PV);
    }

    @Override
    public boolean isLong() {
        return (getMetadata().get(PV) instanceof Long);
    }

    @Override
    public long getLong() {
        return (Long) getMetadata().get(PV);
    }

    @Override
    public boolean isBigInteger() {
        return (getMetadata().get(PV) instanceof BigInteger);
    }

    @Override
    public BigInteger getBigInteger() {
        return (BigInteger) getMetadata().get(PV);
    }

    @Override
    public boolean isFloat() {
        return (getMetadata().get(PV) instanceof Float);
    }

    @Override
    public float getFloat() {
        return (Float) getMetadata().get(PV);
    }

    @Override
    public boolean isDouble() {
        return (getMetadata().get(PV) instanceof Double);
    }

    @Override
    public double getDouble() {
        return (Double) getMetadata().get(PV);
    }

    @Override
    public boolean isBigDecimal() {
        return (getMetadata().get(PV) instanceof BigDecimal); 
    }

    @Override
    public BigDecimal getBigDecimal() {
        return (BigDecimal) getMetadata().get(PV); 
    }

    @Override
    public boolean isString() {
        return (getMetadata().get(PV) instanceof String);   
    }

    @Override
    public String getString() {
        return (String) getMetadata().get(PV); 
    }

    @Override
    public boolean isDuration() {
        return (getMetadata().get(PV) instanceof Duration);  
    }

    @Override
    public Duration getDuration() {
        return (Duration) getMetadata().get(PV);  
    }

    @Override
    public boolean isTime() {
        return (getMetadata().get(PV) instanceof LocalTime);  
    }

    @Override
    public LocalTime getTime() {
        return (LocalTime) getMetadata().get(PV); 
    }

    @Override
    public boolean isDate() {
        return (getMetadata().get(PV) instanceof LocalDate); 
    }

    @Override
    public LocalDate getDate() {
        return (LocalDate) getMetadata().get(PV);  
    }

    @Override
    public boolean isDateTime() {
        return (getMetadata().get(PV) instanceof LocalDateTime); 
    }

    @Override
    public LocalDateTime getDateTime() {
        return (LocalDateTime) getMetadata().get(PV);  
    }

    @Override
    public byte[] getRaw() {
        throw new UnsupportedOperationException("Not supported yet.");  
    }

    @Override
    public boolean isList() {
         return getMetadata().get(PV).getClass().isArray();
    }

    @Override
    public int getLength() {
        return getMetadata().size();
    }

    @Override
    public PlcValue getIndex(int i) {
        throw new UnsupportedOperationException("Not supported yet.");  
    }

    @Override
    public List<? extends PlcValue> getList() {
        return (List<? extends PlcValue>) getMetadata().values().stream().
                filter(o-> (o instanceof PlcValue)).
                collect(toList());
    }

    @Override
    public boolean isStruct() {
        return  (getMetadata().size() > 3);
    }

    @Override
    public Set<String> getKeys() {
        return getMetadata().keySet();
    }

    @Override
    public boolean hasKey(String key) {
        return getMetadata().containsKey(key);
    }

    @Override
    public PlcValue getValue(String key) {
        return new PlcValueImpl.
                PlcValueBuilder().
                putDescription(key).
                putTimeStamp((long) getMetadata().get(FIELD_TIMESTAMP)).
                putValueType(value_type).
                putMetaData(key, getMetadata().get(key)).
                build();
    }

    @Override
    public Map<String, ? extends PlcValue> getStruct() {
        return (Map<String,? extends PlcValue>) (List<? extends PlcValue>) getMetadata().
                entrySet().stream().
                filter(e -> (e.getValue() instanceof PlcValue)).
                collect(toMap(e -> e.getKey(), e -> e.getValue()));
    }

    @Override
    public Set<String> getMetaDataNames() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public boolean hasMetaData(String key) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public PlcValue getMetaData(String key) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }
        

    public static class PlcValueBuilder {
                 
        private final Map<String, Object > metadata;   
    	private long timestamp;
        private PlcValueType value_type;        
        
        public PlcValueBuilder() {          
            metadata = new Hashtable<>();
        }
        
        public PlcValueBuilder putDescription(String desc) {
            metadata.put(DESCRIPTION, desc);
            return this;
        }              
        
        public PlcValueBuilder putTimeStamp(long timestamp) {
            metadata.put(FIELD_TIMESTAMP, timestamp);
            return this;
        } 
        
        public PlcValueBuilder putValueType(PlcValueType value_type) {
            this.value_type = value_type;
            return this;
        }           
        
        public PlcValueBuilder putMetaData(String k, Object o) {
            if (!FIELD_TIMESTAMP.equalsIgnoreCase(k) &&
                !DESCRIPTION.equalsIgnoreCase(k)) {
                this.metadata.put(k, o);
            }
            return this;
        }
        
        public PlcValueImpl build() {
            PlcValueImpl plcvalue = new PlcValueImpl(this);
            validatePlcGroupObject(plcvalue);
            return plcvalue;
        }        
        
        private void validatePlcGroupObject(PlcValue plcvalue) {
            //
        }           
        
    }    
    
}
