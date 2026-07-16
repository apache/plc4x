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

package org.apache.plc4x.java.tools.eventpump;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

/**
 * Wrapper for PlcReadResponse that applies value transformations.
 * <p>
 * This class delegates most calls to the underlying response, but returns
 * transformed values for tags that have transformation expressions configured.
 */
class TransformedPlcReadResponse implements PlcReadResponse {

    private final PlcReadResponse delegate;
    private final Map<String, PlcValue> transformedValues;

    /**
     * Create a transformed response.
     *
     * @param delegate The original response
     * @param transformedValues Map of tag names to transformed values
     */
    public TransformedPlcReadResponse(PlcReadResponse delegate, Map<String, PlcValue> transformedValues) {
        this.delegate = delegate;
        this.transformedValues = new HashMap<>(transformedValues);
    }

    private PlcValue getValueOrDelegate(String tagName) {
        return transformedValues.containsKey(tagName) ? transformedValues.get(tagName) : delegate.getPlcValue(tagName);
    }

    @Override
    public PlcReadRequest getRequest() {
        return delegate.getRequest();
    }

    @Override
    public Collection<String> getTagNames() {
        return delegate.getTagNames();
    }

    @Override
    public PlcTag getTag(String name) {
        return delegate.getTag(name);
    }

    @Override
    public PlcResponseCode getResponseCode(String tagName) {
        return delegate.getResponseCode(tagName);
    }

    @Override
    public Metadata getTagMetadata(String name) {
        return delegate.getTagMetadata(name);
    }

    @Override
    public PlcValue getPlcValue(String tagName) {
        return getValueOrDelegate(tagName);
    }

    @Override
    public PlcValue getAsPlcValue() {
        return delegate.getAsPlcValue();
    }

    @Override
    public int getNumberOfValues(String name) {
        return delegate.getNumberOfValues(name);
    }

    @Override
    public Object getObject(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getObject();
        }
        return delegate.getObject(name);
    }

    @Override
    public Object getObject(String name, int index) {
        return delegate.getObject(name, index);
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        return delegate.getAllObjects(name);
    }

    // Boolean methods
    @Override
    public boolean isValidBoolean(String name) {
        return delegate.isValidBoolean(name);
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        return delegate.isValidBoolean(name, index);
    }

    @Override
    public Boolean getBoolean(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getBoolean();
        }
        return delegate.getBoolean(name);
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        return delegate.getBoolean(name, index);
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        return delegate.getAllBooleans(name);
    }

    // Byte methods
    @Override
    public boolean isValidByte(String name) {
        return delegate.isValidByte(name);
    }

    @Override
    public boolean isValidByte(String name, int index) {
        return delegate.isValidByte(name, index);
    }

    @Override
    public Byte getByte(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getByte();
        }
        return delegate.getByte(name);
    }

    @Override
    public Byte getByte(String name, int index) {
        return delegate.getByte(name, index);
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        return delegate.getAllBytes(name);
    }

    // Short methods
    @Override
    public boolean isValidShort(String name) {
        return delegate.isValidShort(name);
    }

    @Override
    public boolean isValidShort(String name, int index) {
        return delegate.isValidShort(name, index);
    }

    @Override
    public Short getShort(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getShort();
        }
        return delegate.getShort(name);
    }

    @Override
    public Short getShort(String name, int index) {
        return delegate.getShort(name, index);
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        return delegate.getAllShorts(name);
    }

    // Integer methods
    @Override
    public boolean isValidInteger(String name) {
        return delegate.isValidInteger(name);
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        return delegate.isValidInteger(name, index);
    }

    @Override
    public Integer getInteger(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getInteger();
        }
        return delegate.getInteger(name);
    }

    @Override
    public Integer getInteger(String name, int index) {
        return delegate.getInteger(name, index);
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        return delegate.getAllIntegers(name);
    }

    // BigInteger methods
    @Override
    public boolean isValidBigInteger(String name) {
        return delegate.isValidBigInteger(name);
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        return delegate.isValidBigInteger(name, index);
    }

    @Override
    public BigInteger getBigInteger(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getBigInteger();
        }
        return delegate.getBigInteger(name);
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        return delegate.getBigInteger(name, index);
    }

    @Override
    public Collection<BigInteger> getAllBigIntegers(String name) {
        return delegate.getAllBigIntegers(name);
    }

    // Long methods
    @Override
    public boolean isValidLong(String name) {
        return delegate.isValidLong(name);
    }

    @Override
    public boolean isValidLong(String name, int index) {
        return delegate.isValidLong(name, index);
    }

    @Override
    public Long getLong(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getLong();
        }
        return delegate.getLong(name);
    }

    @Override
    public Long getLong(String name, int index) {
        return delegate.getLong(name, index);
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        return delegate.getAllLongs(name);
    }

    // Float methods
    @Override
    public boolean isValidFloat(String name) {
        return delegate.isValidFloat(name);
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        return delegate.isValidFloat(name, index);
    }

    @Override
    public Float getFloat(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getFloat();
        }
        return delegate.getFloat(name);
    }

    @Override
    public Float getFloat(String name, int index) {
        return delegate.getFloat(name, index);
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        return delegate.getAllFloats(name);
    }

    // Double methods
    @Override
    public boolean isValidDouble(String name) {
        return delegate.isValidDouble(name);
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        return delegate.isValidDouble(name, index);
    }

    @Override
    public Double getDouble(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getDouble();
        }
        return delegate.getDouble(name);
    }

    @Override
    public Double getDouble(String name, int index) {
        return delegate.getDouble(name, index);
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        return delegate.getAllDoubles(name);
    }

    // BigDecimal methods
    @Override
    public boolean isValidBigDecimal(String name) {
        return delegate.isValidBigDecimal(name);
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        return delegate.isValidBigDecimal(name, index);
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getBigDecimal();
        }
        return delegate.getBigDecimal(name);
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        return delegate.getBigDecimal(name, index);
    }

    @Override
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        return delegate.getAllBigDecimals(name);
    }

    // String methods
    @Override
    public boolean isValidString(String name) {
        return delegate.isValidString(name);
    }

    @Override
    public boolean isValidString(String name, int index) {
        return delegate.isValidString(name, index);
    }

    @Override
    public String getString(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getString();
        }
        return delegate.getString(name);
    }

    @Override
    public String getString(String name, int index) {
        return delegate.getString(name, index);
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        return delegate.getAllStrings(name);
    }

    // Time methods
    @Override
    public boolean isValidTime(String name) {
        return delegate.isValidTime(name);
    }

    @Override
    public boolean isValidTime(String name, int index) {
        return delegate.isValidTime(name, index);
    }

    @Override
    public LocalTime getTime(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getTime();
        }
        return delegate.getTime(name);
    }

    @Override
    public LocalTime getTime(String name, int index) {
        return delegate.getTime(name, index);
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        return delegate.getAllTimes(name);
    }

    // Date methods
    @Override
    public boolean isValidDate(String name) {
        return delegate.isValidDate(name);
    }

    @Override
    public boolean isValidDate(String name, int index) {
        return delegate.isValidDate(name, index);
    }

    @Override
    public LocalDate getDate(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getDate();
        }
        return delegate.getDate(name);
    }

    @Override
    public LocalDate getDate(String name, int index) {
        return delegate.getDate(name, index);
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        return delegate.getAllDates(name);
    }

    // DateTime methods
    @Override
    public boolean isValidDateTime(String name) {
        return delegate.isValidDateTime(name);
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        return delegate.isValidDateTime(name, index);
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        if (transformedValues.containsKey(name)) {
            return transformedValues.get(name).getDateTime();
        }
        return delegate.getDateTime(name);
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        return delegate.getDateTime(name, index);
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        return delegate.getAllDateTimes(name);
    }
}
