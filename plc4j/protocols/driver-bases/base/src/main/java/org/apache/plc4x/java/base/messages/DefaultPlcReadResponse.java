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
package org.apache.plc4x.java.base.messages;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;

public class DefaultPlcReadResponse implements InternalPlcReadResponse {

    private final InternalPlcReadRequest request;
    private final Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> values;

    public DefaultPlcReadResponse(InternalPlcReadRequest request, Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> fields) {
        this.request = request;
        this.values = fields;
    }

    @Override
    public InternalPlcReadRequest getRequest() {
        return request;
    }

    @Override
    public int getNumberOfValues(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getNumberOfValues();
    }

    @Override
    public Collection<String> getFieldNames() {
        return request.getFieldNames();
    }

    @Override
    public PlcField getField(String name) {
        return request.getField(name);
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        if (values.get(name) == null) {
            throw new PlcInvalidFieldException(name);
        }
        return values.get(name).getKey();
    }

    @Override
    public Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> getValues() {
        return values;
    }

    @Override
    public Object getObject(String name) {
        if(getFieldInternal(name).getNumberOfValues()>1) {
            return getAllObjects(name);
        }
        else{
            return getObject(name, 0);
        }
    }

    @Override
    public Object getObject(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getObject(index);
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Object> objectList = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            objectList.add(fieldInternal.getObject(i));
        }
        return objectList;
    }

    @Override
    public boolean isValidBoolean(String name) {
        return isValidBoolean(name, 0);
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidBoolean(index);
    }

    @Override
    public Boolean getBoolean(String name) {
        return getBoolean(name, 0);
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getBoolean(index);
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Boolean> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getBoolean(i));
        }
        return values;
    }

    @Override
    public boolean isValidByte(String name) {
        return isValidByte(name, 0);
    }

    @Override
    public boolean isValidByte(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidByte(index);
    }

    @Override
    public Byte getByte(String name) {
        return getByte(name, 0);
    }

    @Override
    public Byte getByte(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getByte(index);
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Byte> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getByte(i));
        }
        return values;
    }

    @Override
    public boolean isValidShort(String name) {
        return isValidShort(name, 0);
    }

    @Override
    public boolean isValidShort(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidShort(index);
    }

    @Override
    public Short getShort(String name) {
        return getShort(name, 0);
    }

    @Override
    public Short getShort(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getShort(index);
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Short> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getShort(i));
        }
        return values;
    }

    @Override
    public boolean isValidInteger(String name) {
        return isValidInteger(name, 0);
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidInteger(index);
    }

    @Override
    public Integer getInteger(String name) {
        return getInteger(name, 0);
    }

    @Override
    public Integer getInteger(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getInteger(index);
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Integer> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getInteger(i));
        }
        return values;
    }

    @Override
    public boolean isValidBigInteger(String name) {
        return isValidBigInteger(name, 0);
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidInteger(index);
    }

    @Override
    public BigInteger getBigInteger(String name) {
        return getBigInteger(name, 0);
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getBigInteger(index);
    }

    @Override
    public Collection<BigInteger> getAllBigIntegers(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<BigInteger> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getBigInteger(i));
        }
        return values;
    }

    @Override
    public boolean isValidLong(String name) {
        return isValidLong(name, 0);
    }

    @Override
    public boolean isValidLong(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidLong(index);
    }

    @Override
    public Long getLong(String name) {
        return getLong(name, 0);
    }

    @Override
    public Long getLong(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getLong(index);
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Long> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getLong(i));
        }
        return values;
    }

    @Override
    public boolean isValidFloat(String name) {
        return isValidFloat(name, 0);
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidFloat(index);
    }

    @Override
    public Float getFloat(String name) {
        return getFloat(name, 0);
    }

    @Override
    public Float getFloat(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getFloat(index);
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Float> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getFloat(i));
        }
        return values;
    }

    @Override
    public boolean isValidDouble(String name) {
        return isValidDouble(name, 0);
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidDouble(index);
    }

    @Override
    public Double getDouble(String name) {
        return getDouble(name, 0);
    }

    @Override
    public Double getDouble(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getDouble(index);
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Double> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getDouble(i));
        }
        return values;
    }

    @Override
    public boolean isValidBigDecimal(String name) {
        return isValidBigDecimal(name, 0);
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidBigDecimal(index);
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        return getBigDecimal(name, 0);
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getBigDecimal(index);
    }

    @Override
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<BigDecimal> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getBigDecimal(i));
        }
        return values;
    }

    @Override
    public boolean isValidString(String name) {
        return isValidString(name, 0);
    }

    @Override
    public boolean isValidString(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidString(index);
    }

    @Override
    public String getString(String name) {
        return getString(name, 0);
    }

    @Override
    public String getString(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getString(index);
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<String> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getString(i));
        }
        return values;
    }

    @Override
    public boolean isValidTime(String name) {
        return isValidTime(name, 0);
    }

    @Override
    public boolean isValidTime(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidTime(index);
    }

    @Override
    public LocalTime getTime(String name) {
        return getTime(name, 0);
    }

    @Override
    public LocalTime getTime(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getTime(index);
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<LocalTime> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getTime(i));
        }
        return values;
    }

    @Override
    public boolean isValidDate(String name) {
        return isValidDate(name, 0);
    }

    @Override
    public boolean isValidDate(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidDate(index);
    }

    @Override
    public LocalDate getDate(String name) {
        return getDate(name, 0);
    }

    @Override
    public LocalDate getDate(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getDate(index);
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<LocalDate> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getDate(i));
        }
        return values;
    }

    @Override
    public boolean isValidDateTime(String name) {
        return isValidDateTime(name, 0);
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidDateTime(index);
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        return getDateTime(name, 0);
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getDateTime(index);
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<LocalDateTime> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getDateTime(i));
        }
        return values;
    }

    @Override
    public boolean isValidByteArray(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidByteArray(0);
    }

    @Override
    public boolean isValidByteArray(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.isValidByteArray(index);
    }

    @Override
    public Byte[] getByteArray(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getByteArray(0);
    }

    @Override
    public Byte[] getByteArray(String name, int index) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        return fieldInternal.getByteArray(index);
    }

    @Override
    public Collection<Byte[]> getAllByteArrays(String name) {
        BaseDefaultFieldItem fieldInternal = getFieldInternal(name);
        int num = fieldInternal.getNumberOfValues();
        List<Byte[]> values = new ArrayList<>(num);
        for (int i = 0; i < num; i++) {
            values.add(fieldInternal.getByteArray(i));
        }
        return values;
    }

    protected BaseDefaultFieldItem getFieldInternal(String name) {
        Objects.requireNonNull(name, "Name argument required");
        // If this field doesn't exist, ignore it.
        if (values.get(name) == null) {
            throw new PlcInvalidFieldException(name);
        }
        if (values.get(name).getKey() != PlcResponseCode.OK) {
            throw new PlcRuntimeException("Field '" + name + "' could not be fetched, response was " + values.get(name).getKey());
        }
        // No need to check for "null" as this is already captured by the constructors.
        return values.get(name).getValue();
    }

}
