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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcStruct;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;

public class DefaultPlcSubscriptionEvent implements PlcSubscriptionEvent {

    private final Instant timestamp;
    private final Map<String, PlcResponseItem<PlcValue>> values;

    public DefaultPlcSubscriptionEvent(Instant timestamp, Map<String, PlcResponseItem<PlcValue>> values) {
        this.timestamp = timestamp;
        this.values = values;
    }

    @Override
    public Instant getTimestamp() {
        return timestamp;
    }

    @Override
    public PlcReadRequest getRequest() {
        return null;
    }

    @Override
    public Metadata getTagMetadata(String tag) {
        return Metadata.EMPTY;
    }

    @Override
    public PlcValue getAsPlcValue() {
        Map<String, PlcValue> structMap = new HashMap<>();
        for (String tagName : getTagNames()) {
            structMap.put(tagName, getPlcValue(tagName));
        }
        return new PlcStruct(structMap);
    }

    @Override
    public PlcValue getPlcValue(String name) {
        PlcResponseItem<PlcValue> item = values.get(name);
        return item != null ? item.getValue() : null;
    }

    @Override
    public int getNumberOfValues(String name) {
        PlcValue v = getTagInternal(name);
        if (v instanceof PlcList plcList) {
            return plcList.getLength();
        }
        return 1;
    }

    @Override
    public Collection<String> getTagNames() {
        return values.keySet();
    }

    @Override
    public PlcTag getTag(String name) {
        return null;
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        if (values.get(name) == null) {
            throw new PlcInvalidTagException(name);
        }
        return values.get(name).getResponseCode();
    }

    @Override
    public Object getObject(String name) {
        PlcValue v = getTagInternal(name);
        return v != null ? v.getObject() : null;
    }

    @Override
    public Object getObject(String name, int index) {
        PlcValue v = getTagIndexInternal(name, index);
        return v != null ? v.getObject() : null;
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        PlcValue v = getTagInternal(name);
        if (v instanceof PlcList plcList) {
            List<Object> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getObject());
            }
            return items;
        }
        return Collections.singletonList(v.getObject());
    }

    @Override public boolean isValidBoolean(String name) { return getTagInternal(name).isBoolean(); }
    @Override public boolean isValidBoolean(String name, int index) { return getTagInternal(name).isBoolean(); }
    @Override public Boolean getBoolean(String name) { return getBoolean(name, 0); }
    @Override public Boolean getBoolean(String name, int index) { return getTagIndexInternal(name, index).getBoolean(); }
    @Override public Collection<Boolean> getAllBooleans(String name) { return getAllValues(name, PlcValue::getBoolean); }

    @Override public boolean isValidByte(String name) { return getTagInternal(name).isByte(); }
    @Override public boolean isValidByte(String name, int index) { return getTagIndexInternal(name, index).isByte(); }
    @Override public Byte getByte(String name) { return getByte(name, 0); }
    @Override public Byte getByte(String name, int index) { return getTagIndexInternal(name, index).getByte(); }
    @Override public Collection<Byte> getAllBytes(String name) { return getAllValues(name, PlcValue::getByte); }

    @Override public boolean isValidShort(String name) { return getTagInternal(name).isShort(); }
    @Override public boolean isValidShort(String name, int index) { return getTagIndexInternal(name, index).isShort(); }
    @Override public Short getShort(String name) { return getShort(name, 0); }
    @Override public Short getShort(String name, int index) { return getTagIndexInternal(name, index).getShort(); }
    @Override public Collection<Short> getAllShorts(String name) { return getAllValues(name, PlcValue::getShort); }

    @Override public boolean isValidInteger(String name) { return getTagInternal(name).isInteger(); }
    @Override public boolean isValidInteger(String name, int index) { return getTagIndexInternal(name, index).isInteger(); }
    @Override public Integer getInteger(String name) { return getInteger(name, 0); }
    @Override public Integer getInteger(String name, int index) { return getTagIndexInternal(name, index).getInteger(); }
    @Override public Collection<Integer> getAllIntegers(String name) { return getAllValues(name, PlcValue::getInteger); }

    @Override public boolean isValidBigInteger(String name) { return getTagInternal(name).isBigInteger(); }
    @Override public boolean isValidBigInteger(String name, int index) { return getTagIndexInternal(name, index).isBigInteger(); }
    @Override public BigInteger getBigInteger(String name) { return getBigInteger(name, 0); }
    @Override public BigInteger getBigInteger(String name, int index) { return getTagIndexInternal(name, index).getBigInteger(); }
    @Override public Collection<BigInteger> getAllBigIntegers(String name) { return getAllValues(name, PlcValue::getBigInteger); }

    @Override public boolean isValidLong(String name) { return getTagInternal(name).isLong(); }
    @Override public boolean isValidLong(String name, int index) { return getTagIndexInternal(name, index).isLong(); }
    @Override public Long getLong(String name) { return getLong(name, 0); }
    @Override public Long getLong(String name, int index) { return getTagIndexInternal(name, index).getLong(); }
    @Override public Collection<Long> getAllLongs(String name) { return getAllValues(name, PlcValue::getLong); }

    @Override public boolean isValidFloat(String name) { return getTagInternal(name).isFloat(); }
    @Override public boolean isValidFloat(String name, int index) { return getTagIndexInternal(name, index).isFloat(); }
    @Override public Float getFloat(String name) { return getFloat(name, 0); }
    @Override public Float getFloat(String name, int index) { return getTagIndexInternal(name, index).getFloat(); }
    @Override public Collection<Float> getAllFloats(String name) { return getAllValues(name, PlcValue::getFloat); }

    @Override public boolean isValidDouble(String name) { return getTagInternal(name).isDouble(); }
    @Override public boolean isValidDouble(String name, int index) { return getTagIndexInternal(name, index).isDouble(); }
    @Override public Double getDouble(String name) { return getDouble(name, 0); }
    @Override public Double getDouble(String name, int index) { return getTagIndexInternal(name, index).getDouble(); }
    @Override public Collection<Double> getAllDoubles(String name) { return getAllValues(name, PlcValue::getDouble); }

    @Override public boolean isValidBigDecimal(String name) { return getTagInternal(name).isBigDecimal(); }
    @Override public boolean isValidBigDecimal(String name, int index) { return getTagIndexInternal(name, index).isBigDecimal(); }
    @Override public BigDecimal getBigDecimal(String name) { return getBigDecimal(name, 0); }
    @Override public BigDecimal getBigDecimal(String name, int index) { return getTagIndexInternal(name, index).getBigDecimal(); }
    @Override public Collection<BigDecimal> getAllBigDecimals(String name) { return getAllValues(name, PlcValue::getBigDecimal); }

    @Override public boolean isValidString(String name) { return getTagInternal(name).isString(); }
    @Override public boolean isValidString(String name, int index) { return getTagIndexInternal(name, index).isString(); }
    @Override public String getString(String name) { return getString(name, 0); }
    @Override public String getString(String name, int index) { return getTagIndexInternal(name, index).getString(); }
    @Override public Collection<String> getAllStrings(String name) { return getAllValues(name, PlcValue::getString); }

    @Override public boolean isValidTime(String name) { return getTagInternal(name).isTime(); }
    @Override public boolean isValidTime(String name, int index) { return getTagIndexInternal(name, index).isTime(); }
    @Override public LocalTime getTime(String name) { return getTime(name, 0); }
    @Override public LocalTime getTime(String name, int index) { return getTagIndexInternal(name, index).getTime(); }
    @Override public Collection<LocalTime> getAllTimes(String name) { return getAllValues(name, PlcValue::getTime); }

    @Override public boolean isValidDate(String name) { return getTagInternal(name).isDate(); }
    @Override public boolean isValidDate(String name, int index) { return getTagIndexInternal(name, index).isDate(); }
    @Override public LocalDate getDate(String name) { return getDate(name, 0); }
    @Override public LocalDate getDate(String name, int index) { return getTagIndexInternal(name, index).getDate(); }
    @Override public Collection<LocalDate> getAllDates(String name) { return getAllValues(name, PlcValue::getDate); }

    @Override public boolean isValidDateTime(String name) { return getTagInternal(name).isDateTime(); }
    @Override public boolean isValidDateTime(String name, int index) { return getTagIndexInternal(name, index).isDateTime(); }
    @Override public LocalDateTime getDateTime(String name) { return getDateTime(name, 0); }
    @Override public LocalDateTime getDateTime(String name, int index) { return getTagIndexInternal(name, index).getDateTime(); }
    @Override public Collection<LocalDateTime> getAllDateTimes(String name) { return getAllValues(name, PlcValue::getDateTime); }

    private <T> Collection<T> getAllValues(String name, java.util.function.Function<PlcValue, T> extractor) {
        PlcValue v = getTagInternal(name);
        if (v instanceof PlcList plcList) {
            List<T> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(extractor.apply(plcValue));
            }
            return items;
        }
        return Collections.singletonList(extractor.apply(v));
    }

    private PlcValue getTagInternal(String name) {
        Objects.requireNonNull(name, "Name argument required");
        if (values.get(name) == null) {
            throw new PlcInvalidTagException(name);
        }
        if (values.get(name).getResponseCode() != PlcResponseCode.OK) {
            throw new PlcRuntimeException(
                "Tag '" + name + "' could not be fetched, response was " + values.get(name).getResponseCode());
        }
        return values.get(name).getValue();
    }

    private PlcValue getTagIndexInternal(String name, int index) {
        PlcValue v = getTagInternal(name);
        if (v instanceof PlcList plcList) {
            if (index > (plcList.getLength() - 1)) {
                return null;
            }
            return plcList.getIndex(index);
        }
        if (index != 0) {
            return null;
        }
        return v;
    }

}
