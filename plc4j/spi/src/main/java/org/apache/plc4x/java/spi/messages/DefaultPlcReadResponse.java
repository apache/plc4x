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
package org.apache.plc4x.java.spi.messages;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValues;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.function.Function;
import java.util.stream.Collectors;

public class DefaultPlcReadResponse implements InternalPlcReadResponse, PlcReadResponse {

    private final InternalPlcReadRequest request;
    private final Map<String, Pair<PlcResponseCode, PlcValue>> values;

    public DefaultPlcReadResponse(InternalPlcReadRequest request, Map<String, Pair<PlcResponseCode, PlcValue>> fields) {
        this.request = request;
        this.values = fields;
    }

    @Override
    public InternalPlcReadRequest getRequest() {
        return request;
    }

    @Override public PlcValue getAsPlcValue() {
        return PlcValues.of(request.getFieldNames().stream()
            .collect(Collectors.toMap(Function.identity(), name -> PlcValues.of(getObject(name)))));
    }

    @Override
    public int getNumberOfValues(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            return plcList.getNumberOfValues();
        } else {
            return 1;
        }
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
    public Map<String, Pair<PlcResponseCode, PlcValue>> getValues() {
        return values;
    }

    @Override
    public Object getObject(String name) {
        return getFieldInternal(name).getObject();
    }

    @Override
    public Object getObject(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getObject();
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<Object> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getObject());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getObject());
    }

    @Override
    public boolean isValidBoolean(String name) {
        return isValidBoolean(name, 0);
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        PlcValue fieldInternal = getFieldInternal(name);
        return fieldInternal.isBoolean();
    }

    @Override
    public Boolean getBoolean(String name) {
        return getBoolean(name, 0);
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getBoolean();
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<Boolean> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getBoolean());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getBoolean());
    }

    @Override
    public boolean isValidByte(String name) {
        return isValidByte(name, 0);
    }

    @Override
    public boolean isValidByte(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isByte();
    }

    @Override
    public Byte getByte(String name) {
        return getByte(name, 0);
    }

    @Override
    public Byte getByte(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getByte();
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<Byte> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getByte());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getByte());
    }

    @Override
    public boolean isValidShort(String name) {
        return isValidShort(name, 0);
    }

    @Override
    public boolean isValidShort(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isShort();
    }

    @Override
    public Short getShort(String name) {
        return getShort(name, 0);
    }

    @Override
    public Short getShort(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getShort();
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<Short> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getShort());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getShort());
    }

    @Override
    public boolean isValidInteger(String name) {
        return isValidInteger(name, 0);
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isInteger();
    }

    @Override
    public Integer getInteger(String name) {
        return getInteger(name, 0);
    }

    @Override
    public Integer getInteger(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getInteger();
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<Integer> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getInteger());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getInteger());
    }

    @Override
    public boolean isValidBigInteger(String name) {
        return isValidBigInteger(name, 0);
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isBigInteger();
    }

    @Override
    public BigInteger getBigInteger(String name) {
        return getBigInteger(name, 0);
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getBigInteger();
    }

    @Override
    public Collection<BigInteger> getAllBigIntegers(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<BigInteger> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getBigInteger());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getBigInteger());
    }

    @Override
    public boolean isValidLong(String name) {
        return isValidLong(name, 0);
    }

    @Override
    public boolean isValidLong(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isLong();
    }

    @Override
    public Long getLong(String name) {
        return getLong(name, 0);
    }

    @Override
    public Long getLong(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getLong();
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<Long> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getLong());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getLong());
    }

    @Override
    public boolean isValidFloat(String name) {
        return isValidFloat(name, 0);
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isFloat();
    }

    @Override
    public Float getFloat(String name) {
        return getFloat(name, 0);
    }

    @Override
    public Float getFloat(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getFloat();
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<Float> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getFloat());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getFloat());
    }

    @Override
    public boolean isValidDouble(String name) {
        return isValidDouble(name, 0);
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isDouble();
    }

    @Override
    public Double getDouble(String name) {
        return getDouble(name, 0);
    }

    @Override
    public Double getDouble(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getDouble();
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<Double> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getDouble());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getDouble());
    }

    @Override
    public boolean isValidBigDecimal(String name) {
        return isValidBigDecimal(name, 0);
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isBigDecimal();
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        return getBigDecimal(name, 0);
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getBigDecimal();
    }

    @Override
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<BigDecimal> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getBigDecimal());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getBigDecimal());
    }

    @Override
    public boolean isValidString(String name) {
        return isValidString(name, 0);
    }

    @Override
    public boolean isValidString(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isString();
    }

    @Override
    public String getString(String name) {
        return getString(name, 0);
    }

    @Override
    public String getString(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getString();
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<String> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getString());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getString());
    }

    @Override
    public boolean isValidTime(String name) {
        return isValidTime(name, 0);
    }

    @Override
    public boolean isValidTime(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isTime();
    }

    @Override
    public LocalTime getTime(String name) {
        return getTime(name, 0);
    }

    @Override
    public LocalTime getTime(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getTime();
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<LocalTime> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getTime());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getTime());
    }

    @Override
    public boolean isValidDate(String name) {
        return isValidDate(name, 0);
    }

    @Override
    public boolean isValidDate(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isDate();
    }

    @Override
    public LocalDate getDate(String name) {
        return getDate(name, 0);
    }

    @Override
    public LocalDate getDate(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getDate();
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<LocalDate> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getDate());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getDate());
    }

    @Override
    public boolean isValidDateTime(String name) {
        return isValidDateTime(name, 0);
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.isDateTime();
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        return getDateTime(name, 0);
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        PlcValue fieldInternal = getFieldIndexInternal(name, index);
        return fieldInternal.getDateTime();
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        PlcValue fieldInternal = getFieldInternal(name);
        if(fieldInternal instanceof PlcList) {
            PlcList plcList = (PlcList) fieldInternal;
            List<LocalDateTime> items = new ArrayList<>(plcList.getNumberOfValues());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getDateTime());
            }
            return items;
        }
        return Collections.singletonList(fieldInternal.getDateTime());
    }

    protected PlcValue getFieldInternal(String name) {
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

    protected PlcValue getFieldIndexInternal(String name, int index) {
        final PlcValue field = getFieldInternal(name);
        if(field instanceof PlcList) {
            PlcList plcList = (PlcList) field;
            if(index > (plcList.getNumberOfValues() - 1)) {
                return null;
            }
            return plcList.getIndex(index);
        }
        if(index != 0) {
            return null;
        }
        return field;
    }

}
