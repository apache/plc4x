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
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.FieldItem;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.Map;

public class DefaultPlcReadResponse implements InternalPlcReadResponse {

    private final InternalPlcReadRequest request;
    private final Map<String, Pair<PlcResponseCode, FieldItem>> values;

    public DefaultPlcReadResponse(InternalPlcReadRequest request, Map<String, Pair<PlcResponseCode, FieldItem>> fields) {
        this.request = request;
        this.values = fields;
    }

    @Override
    public InternalPlcReadRequest getRequest() {
        return request;
    }

    @Override
    public int getNumberOfValues(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getNumberOfValues();
        }
        return 0;
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
            return null;
        }
        return values.get(name).getKey();
    }

    @Override
    public Map<String, Pair<PlcResponseCode, FieldItem>> getValues() {
        return values;
    }

    @Override
    public Object getObject(String name) {
        return getObject(name, 0);
    }

    @Override
    public Object getObject(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getObject(index);
        }
        return null;
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        // TODO: Implement this ...
        return null;
    }

    @Override
    public boolean isValidBoolean(String name) {
        return isValidBoolean(name, 0);
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidBoolean(index);
        }
        return false;
    }

    @Override
    public Boolean getBoolean(String name) {
        return getBoolean(name, 0);
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getBoolean(index);
        }
        return null;
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<Boolean> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getBoolean(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidByte(String name) {
        return isValidByte(name, 0);
    }

    @Override
    public boolean isValidByte(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidByte(index);
        }
        return false;
    }

    @Override
    public Byte getByte(String name) {
        return getByte(name, 0);
    }

    @Override
    public Byte getByte(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getByte(index);
        }
        return null;
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<Byte> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getByte(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidShort(String name) {
        return isValidShort(name, 0);
    }

    @Override
    public boolean isValidShort(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidShort(index);
        }
        return false;
    }

    @Override
    public Short getShort(String name) {
        return getShort(name, 0);
    }

    @Override
    public Short getShort(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getShort(index);
        }
        return null;
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<Short> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getShort(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidInteger(String name) {
        return isValidInteger(name, 0);
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidInteger(index);
        }
        return false;
    }

    @Override
    public Integer getInteger(String name) {
        return getInteger(name, 0);
    }

    @Override
    public Integer getInteger(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getInteger(index);
        }
        return null;
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<Integer> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getInteger(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidBigInteger(String name) {
        return isValidBigInteger(name, 0);
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidInteger(index);
        }
        return false;
    }

    @Override
    public BigInteger getBigInteger(String name) {
        return getBigInteger(name, 0);
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getBigInteger(index);
        }
        return null;
    }

    @Override
    public Collection<BigInteger> getAllBigIntegers(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<BigInteger> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getBigInteger(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidLong(String name) {
        return isValidLong(name, 0);
    }

    @Override
    public boolean isValidLong(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidLong(index);
        }
        return false;
    }

    @Override
    public Long getLong(String name) {
        return getLong(name, 0);
    }

    @Override
    public Long getLong(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getLong(index);
        }
        return null;
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<Long> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getLong(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidFloat(String name) {
        return isValidFloat(name, 0);
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidFloat(index);
        }
        return false;
    }

    @Override
    public Float getFloat(String name) {
        return getFloat(name, 0);
    }

    @Override
    public Float getFloat(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getFloat(index);
        }
        return null;
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<Float> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getFloat(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidDouble(String name) {
        return isValidDouble(name, 0);
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidDouble(index);
        }
        return false;
    }

    @Override
    public Double getDouble(String name) {
        return getDouble(name, 0);
    }

    @Override
    public Double getDouble(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getDouble(index);
        }
        return null;
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<Double> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getDouble(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidBigDecimal(String name) {
        return isValidBigDecimal(name, 0);
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidBigDecimal(index);
        }
        return false;
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        return getBigDecimal(name, 0);
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getBigDecimal(index);
        }
        return null;
    }

    @Override
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<BigDecimal> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getBigDecimal(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidString(String name) {
        return isValidString(name, 0);
    }

    @Override
    public boolean isValidString(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidString(index);
        }
        return false;
    }

    @Override
    public String getString(String name) {
        return getString(name, 0);
    }

    @Override
    public String getString(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getString(index);
        }
        return null;
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<String> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getString(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidTime(String name) {
        return isValidTime(name, 0);
    }

    @Override
    public boolean isValidTime(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidTime(index);
        }
        return false;
    }

    @Override
    public LocalTime getTime(String name) {
        return getTime(name, 0);
    }

    @Override
    public LocalTime getTime(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getTime(index);
        }
        return null;
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<LocalTime> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getTime(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidDate(String name) {
        return isValidDate(name, 0);
    }

    @Override
    public boolean isValidDate(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidDate(index);
        }
        return false;
    }

    @Override
    public LocalDate getDate(String name) {
        return getDate(name, 0);
    }

    @Override
    public LocalDate getDate(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getDate(index);
        }
        return null;
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<LocalDate> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getDate(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidDateTime(String name) {
        return isValidDateTime(name, 0);
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidDateTime(index);
        }
        return false;
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        return getDateTime(name, 0);
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getDateTime(index);
        }
        return null;
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<LocalDateTime> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getDateTime(i));
            }
            return values;
        }
        return null;
    }

    @Override
    public boolean isValidByteArray(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidByteArray(0);
        }
        return false;
    }

    @Override
    public boolean isValidByteArray(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.isValidDateTime(index);
        }
        return false;
    }

    @Override
    public Byte[] getByteArray(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getByteArray(0);
        }
        return null;
    }

    @Override
    public Byte[] getByteArray(String name, int index) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            return fieldInternal.getByteArray(index);
        }
        return null;
    }

    @Override
    public Collection<Byte[]> getAllByteArrays(String name) {
        FieldItem fieldInternal = getFieldInternal(name);
        if (fieldInternal != null) {
            int num = fieldInternal.getNumberOfValues();
            List<Byte[]> values = new ArrayList<>(num);
            for (int i = 0; i < num; i++) {
                values.add(fieldInternal.getByteArray(i));
            }
            return values;
        }
        return null;
    }

    private FieldItem getFieldInternal(String name) {
        // If this field doesn't exist, ignore it.
        if (values.get(name) == null) {
            return null;
        }
        if (values.get(name).getKey() != PlcResponseCode.OK) {
            return null;
        }
        return values.get(name).getValue();
    }

}
