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
package org.apache.plc4x.java.spi.messages;

import com.fasterxml.jackson.annotation.*;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcReadResponse implements PlcReadResponse, Serializable {

    private final PlcReadRequest request;
    private final Map<String, ResponseItem<PlcValue>> values;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcReadResponse(@JsonProperty("request") PlcReadRequest request,
                                  @JsonProperty("values") Map<String, ResponseItem<PlcValue>> values) {
        this.request = request;
        this.values = values;
    }

    @Override
    public PlcReadRequest getRequest() {
        return request;
    }

    @Override
    @JsonIgnore
    public PlcValue getAsPlcValue() {
        Map<String, PlcValue> structMap = new HashMap<>();
        for (String tagName : request.getTagNames()) {
            PlcValue plcValue = getPlcValue(tagName);
            structMap.put(tagName, plcValue);
        }
        return new PlcStruct(structMap);
    }

    @Override
    @JsonIgnore
    public PlcValue getPlcValue(String name) {
        return values.getOrDefault(name, new ResponseItem<>(null, null)).getValue();
    }

    @Override
    @JsonIgnore
    public int getNumberOfValues(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            return plcList.getLength();
        } else {
            return 1;
        }
    }

    @Override
    @JsonIgnore
    public Collection<String> getTagNames() {
        return request.getTagNames();
    }

    @Override
    @JsonIgnore
    public PlcTag getTag(String name) {
        return request.getTag(name);
    }

    @Override
    @JsonIgnore
    public PlcResponseCode getResponseCode(String name) {
        if (values.get(name) == null) {
            throw new PlcInvalidTagException(name);
        }
        return values.get(name).getCode();
    }

    @JsonIgnore
    public Map<String, ResponseItem<PlcValue>> getValues() {
        return values;
    }

    @Override
    @JsonIgnore
    public Object getObject(String name) {
        if(getTagInternal(name) != null) {
            return getTagInternal(name).getObject();
        }
        return null;
    }

    @Override
    @JsonIgnore
    public Object getObject(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        if(tagInternal != null) {
            return tagInternal.getObject();
        }
        return null;
    }

    @Override
    @JsonIgnore
    public Collection<Object> getAllObjects(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<Object> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getObject());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getObject());
    }

    @Override
    @JsonIgnore
    public boolean isValidBoolean(String name) {
        return isValidBoolean(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidBoolean(String name, int index) {
        PlcValue tagInternal = getTagInternal(name);
        return tagInternal.isBoolean();
    }

    @Override
    @JsonIgnore
    public Boolean getBoolean(String name) {
        return getBoolean(name, 0);
    }

    @Override
    @JsonIgnore
    public Boolean getBoolean(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getBoolean();
    }

    @Override
    @JsonIgnore
    public Collection<Boolean> getAllBooleans(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<Boolean> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getBoolean());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getBoolean());
    }

    @Override
    @JsonIgnore
    public boolean isValidByte(String name) {
        return isValidByte(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidByte(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isByte();
    }

    @Override
    @JsonIgnore
    public Byte getByte(String name) {
        return getByte(name, 0);
    }

    @Override
    @JsonIgnore
    public Byte getByte(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getByte();
    }

    @Override
    @JsonIgnore
    public Collection<Byte> getAllBytes(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<Byte> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getByte());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getByte());
    }

    @Override
    @JsonIgnore
    public boolean isValidShort(String name) {
        return isValidShort(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidShort(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isShort();
    }

    @Override
    @JsonIgnore
    public Short getShort(String name) {
        return getShort(name, 0);
    }

    @Override
    @JsonIgnore
    public Short getShort(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getShort();
    }

    @Override
    @JsonIgnore
    public Collection<Short> getAllShorts(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<Short> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getShort());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getShort());
    }

    @Override
    @JsonIgnore
    public boolean isValidInteger(String name) {
        return isValidInteger(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidInteger(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isInteger();
    }

    @Override
    @JsonIgnore
    public Integer getInteger(String name) {
        return getInteger(name, 0);
    }

    @Override
    @JsonIgnore
    public Integer getInteger(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getInteger();
    }

    @Override
    @JsonIgnore
    public Collection<Integer> getAllIntegers(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<Integer> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getInteger());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getInteger());
    }

    @Override
    @JsonIgnore
    public boolean isValidBigInteger(String name) {
        return isValidBigInteger(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidBigInteger(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isBigInteger();
    }

    @Override
    @JsonIgnore
    public BigInteger getBigInteger(String name) {
        return getBigInteger(name, 0);
    }

    @Override
    @JsonIgnore
    public BigInteger getBigInteger(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getBigInteger();
    }

    @Override
    @JsonIgnore
    public Collection<BigInteger> getAllBigIntegers(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<BigInteger> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getBigInteger());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getBigInteger());
    }

    @Override
    @JsonIgnore
    public boolean isValidLong(String name) {
        return isValidLong(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidLong(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isLong();
    }

    @Override
    @JsonIgnore
    public Long getLong(String name) {
        return getLong(name, 0);
    }

    @Override
    @JsonIgnore
    public Long getLong(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getLong();
    }

    @Override
    @JsonIgnore
    public Collection<Long> getAllLongs(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<Long> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getLong());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getLong());
    }

    @Override
    @JsonIgnore
    public boolean isValidFloat(String name) {
        return isValidFloat(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidFloat(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isFloat();
    }

    @Override
    @JsonIgnore
    public Float getFloat(String name) {
        return getFloat(name, 0);
    }

    @Override
    @JsonIgnore
    public Float getFloat(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getFloat();
    }

    @Override
    @JsonIgnore
    public Collection<Float> getAllFloats(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<Float> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getFloat());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getFloat());
    }

    @Override
    @JsonIgnore
    public boolean isValidDouble(String name) {
        return isValidDouble(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidDouble(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isDouble();
    }

    @Override
    @JsonIgnore
    public Double getDouble(String name) {
        return getDouble(name, 0);
    }

    @Override
    @JsonIgnore
    public Double getDouble(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getDouble();
    }

    @Override
    @JsonIgnore
    public Collection<Double> getAllDoubles(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<Double> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getDouble());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getDouble());
    }

    @Override
    @JsonIgnore
    public boolean isValidBigDecimal(String name) {
        return isValidBigDecimal(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidBigDecimal(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isBigDecimal();
    }

    @Override
    @JsonIgnore
    public BigDecimal getBigDecimal(String name) {
        return getBigDecimal(name, 0);
    }

    @Override
    @JsonIgnore
    public BigDecimal getBigDecimal(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getBigDecimal();
    }

    @Override
    @JsonIgnore
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<BigDecimal> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getBigDecimal());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getBigDecimal());
    }

    @Override
    @JsonIgnore
    public boolean isValidString(String name) {
        return isValidString(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidString(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isString();
    }

    @Override
    @JsonIgnore
    public String getString(String name) {
        return getString(name, 0);
    }

    @Override
    @JsonIgnore
    public String getString(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getString();
    }

    @Override
    @JsonIgnore
    public Collection<String> getAllStrings(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<String> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getString());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getString());
    }

    @Override
    @JsonIgnore
    public boolean isValidTime(String name) {
        return isValidTime(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidTime(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isTime();
    }

    @Override
    @JsonIgnore
    public LocalTime getTime(String name) {
        return getTime(name, 0);
    }

    @Override
    @JsonIgnore
    public LocalTime getTime(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getTime();
    }

    @Override
    @JsonIgnore
    public Collection<LocalTime> getAllTimes(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<LocalTime> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getTime());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getTime());
    }

    @Override
    @JsonIgnore
    public boolean isValidDate(String name) {
        return isValidDate(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidDate(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isDate();
    }

    @Override
    @JsonIgnore
    public LocalDate getDate(String name) {
        return getDate(name, 0);
    }

    @Override
    @JsonIgnore
    public LocalDate getDate(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getDate();
    }

    @Override
    @JsonIgnore
    public Collection<LocalDate> getAllDates(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<LocalDate> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getDate());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getDate());
    }

    @Override
    @JsonIgnore
    public boolean isValidDateTime(String name) {
        return isValidDateTime(name, 0);
    }

    @Override
    @JsonIgnore
    public boolean isValidDateTime(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isDateTime();
    }

    @Override
    @JsonIgnore
    public LocalDateTime getDateTime(String name) {
        return getDateTime(name, 0);
    }

    @Override
    @JsonIgnore
    public LocalDateTime getDateTime(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getDateTime();
    }

    @Override
    @JsonIgnore
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        PlcValue tagInternal = getTagInternal(name);
        if(tagInternal instanceof PlcList) {
            PlcList plcList = (PlcList) tagInternal;
            List<LocalDateTime> items = new ArrayList<>(plcList.getLength());
            for (PlcValue plcValue : plcList.getList()) {
                items.add(plcValue.getDateTime());
            }
            return items;
        }
        return Collections.singletonList(tagInternal.getDateTime());
    }

    @JsonAnySetter
    public void add(String key, ResponseItem<PlcValue> value) {
        values.put(key, value);
    }

    @JsonAnyGetter
    public Map<String, ResponseItem<PlcValue>> getMap() {
        return values;
    }

    protected PlcValue getTagInternal(String name) {
        Objects.requireNonNull(name, "Name argument required");
        // If this tag doesn't exist, ignore it.
        if (values.get(name) == null) {
            throw new PlcInvalidTagException(name);
        }
        if (values.get(name).getCode() != PlcResponseCode.OK) {
            throw new PlcRuntimeException(
                "Tag '" + name + "' could not be fetched, response was " + values.get(name).getCode());
        }
        // No need to check for "null" as this is already captured by the constructors.
        return values.get(name).getValue();
    }

    protected PlcValue getTagIndexInternal(String name, int index) {
        final PlcValue values = getTagInternal(name);
        if(values instanceof PlcList) {
            PlcList plcList = (PlcList) values;
            if(index > (plcList.getLength() - 1)) {
                return null;
            }
            return plcList.getIndex(index);
        }
        if(index != 0) {
            return null;
        }
        return values;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcReadResponse");

        if(request instanceof Serializable) {
            ((Serializable) request).serialize(writeBuffer);
        }
        writeBuffer.pushContext("values");
        for (Map.Entry<String, ResponseItem<PlcValue>> valueEntry : values.entrySet()) {
            String tagName = valueEntry.getKey();
            writeBuffer.pushContext(tagName);
            ResponseItem<PlcValue> valueResponse = valueEntry.getValue();
            valueResponse.serialize(writeBuffer);
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("values");

        writeBuffer.popContext("PlcReadResponse");
    }

}
