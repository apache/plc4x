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

import java.util.Map.Entry;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.spi.metadata.DefaultMetadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.PlcStruct;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithRenderAsList;

public class DefaultPlcReadResponse implements PlcReadResponse, Serializable {

    private final PlcReadRequest request;
    private final Map<String, PlcResponseItem<PlcValue>> values;
    private final Map<String, Metadata> metadata;

    public DefaultPlcReadResponse(PlcReadRequest request,
                                  Map<String, PlcResponseItem<PlcValue>> values) {
        this(request, values, Collections.emptyMap());
    }

    public DefaultPlcReadResponse(PlcReadRequest request,
                                  Map<String, PlcResponseItem<PlcValue>> values,
                                  Map<String, Metadata> metadata) {
        this.request = request;
        this.values = values;
        this.metadata = Collections.unmodifiableMap(metadata);
    }

    @Override
    public PlcReadRequest getRequest() {
        return request;
    }

    @Override
    public Metadata getTagMetadata(String tag) {
        return metadata.getOrDefault(tag, DefaultMetadata.EMPTY);
    }

    @Override
    public PlcValue getAsPlcValue() {
        Map<String, PlcValue> structMap = new HashMap<>();
        for (String tagName : request.getTagNames()) {
            PlcValue plcValue = getPlcValue(tagName);
            structMap.put(tagName, plcValue);
        }
        return new PlcStruct(structMap);
    }

    @Override
    public PlcValue getPlcValue(String name) {
        return values.getOrDefault(name, new DefaultPlcResponseItem<>(null, null)).getValue();
    }

    public PlcResponseItem<PlcValue> getPlcResponseItem(String name) {
        return values.getOrDefault(name, new DefaultPlcResponseItem<>(null, null));
    }

    @Override
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
    public Collection<String> getTagNames() {
        return request.getTagNames();
    }

    @Override
    public PlcTag getTag(String name) {
        return request.getTag(name);
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        if (values.get(name) == null) {
            throw new PlcInvalidTagException(name);
        }
        return values.get(name).getResponseCode();
    }

    public Map<String, PlcResponseItem<PlcValue>> getValues() {
        return values;
    }

    @Override
    public Object getObject(String name) {
        if(getTagInternal(name) != null) {
            return getTagInternal(name).getObject();
        }
        return null;
    }

    @Override
    public Object getObject(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        if(tagInternal != null) {
            return tagInternal.getObject();
        }
        return null;
    }

    @Override
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
    public boolean isValidBoolean(String name) {
        return isValidBoolean(name, 0);
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        PlcValue tagInternal = getTagInternal(name);
        return tagInternal.isBoolean();
    }

    @Override
    public Boolean getBoolean(String name) {
        return getBoolean(name, 0);
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getBoolean();
    }

    @Override
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
    public boolean isValidByte(String name) {
        return isValidByte(name, 0);
    }

    @Override
    public boolean isValidByte(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isByte();
    }

    @Override
    public Byte getByte(String name) {
        return getByte(name, 0);
    }

    @Override
    public Byte getByte(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getByte();
    }

    @Override
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
    public boolean isValidShort(String name) {
        return isValidShort(name, 0);
    }

    @Override
    public boolean isValidShort(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isShort();
    }

    @Override
    public Short getShort(String name) {
        return getShort(name, 0);
    }

    @Override
    public Short getShort(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getShort();
    }

    @Override
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
    public boolean isValidInteger(String name) {
        return isValidInteger(name, 0);
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isInteger();
    }

    @Override
    public Integer getInteger(String name) {
        return getInteger(name, 0);
    }

    @Override
    public Integer getInteger(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getInteger();
    }

    @Override
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
    public boolean isValidBigInteger(String name) {
        return isValidBigInteger(name, 0);
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isBigInteger();
    }

    @Override
    public BigInteger getBigInteger(String name) {
        return getBigInteger(name, 0);
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getBigInteger();
    }

    @Override
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
    public boolean isValidLong(String name) {
        return isValidLong(name, 0);
    }

    @Override
    public boolean isValidLong(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isLong();
    }

    @Override
    public Long getLong(String name) {
        return getLong(name, 0);
    }

    @Override
    public Long getLong(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getLong();
    }

    @Override
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
    public boolean isValidFloat(String name) {
        return isValidFloat(name, 0);
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isFloat();
    }

    @Override
    public Float getFloat(String name) {
        return getFloat(name, 0);
    }

    @Override
    public Float getFloat(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getFloat();
    }

    @Override
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
    public boolean isValidDouble(String name) {
        return isValidDouble(name, 0);
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isDouble();
    }

    @Override
    public Double getDouble(String name) {
        return getDouble(name, 0);
    }

    @Override
    public Double getDouble(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getDouble();
    }

    @Override
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
    public boolean isValidBigDecimal(String name) {
        return isValidBigDecimal(name, 0);
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isBigDecimal();
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        return getBigDecimal(name, 0);
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getBigDecimal();
    }

    @Override
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
    public boolean isValidString(String name) {
        return isValidString(name, 0);
    }

    @Override
    public boolean isValidString(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isString();
    }

    @Override
    public String getString(String name) {
        return getString(name, 0);
    }

    @Override
    public String getString(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getString();
    }

    @Override
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
    public boolean isValidTime(String name) {
        return isValidTime(name, 0);
    }

    @Override
    public boolean isValidTime(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isTime();
    }

    @Override
    public LocalTime getTime(String name) {
        return getTime(name, 0);
    }

    @Override
    public LocalTime getTime(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getTime();
    }

    @Override
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
    public boolean isValidDate(String name) {
        return isValidDate(name, 0);
    }

    @Override
    public boolean isValidDate(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isDate();
    }

    @Override
    public LocalDate getDate(String name) {
        return getDate(name, 0);
    }

    @Override
    public LocalDate getDate(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getDate();
    }

    @Override
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
    public boolean isValidDateTime(String name) {
        return isValidDateTime(name, 0);
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.isDateTime();
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        return getDateTime(name, 0);
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        PlcValue tagInternal = getTagIndexInternal(name, index);
        return tagInternal.getDateTime();
    }

    @Override
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

    public void add(String key, PlcResponseItem<PlcValue> value) {
        values.put(key, value);
    }

    public Map<String, PlcResponseItem<PlcValue>> getMap() {
        return values;
    }

    protected PlcValue getTagInternal(String name) {
        Objects.requireNonNull(name, "Name argument required");
        // If this tag doesn't exist, ignore it.
        if (values.get(name) == null) {
            throw new PlcInvalidTagException(name);
        }
        if (values.get(name).getResponseCode() != PlcResponseCode.OK) {
            throw new PlcRuntimeException(
                "Tag '" + name + "' could not be fetched, response was " + values.get(name).getResponseCode());
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

        writeBuffer.pushContext("request");
        if(request instanceof Serializable) {
            ((Serializable) request).serialize(writeBuffer);
        }
        writeBuffer.popContext("request");

        writeBuffer.pushContext("values", WithRenderAsList(true));
        for (Map.Entry<String, PlcResponseItem<PlcValue>> valueEntry : values.entrySet()) {
            String tagName = valueEntry.getKey();
            writeBuffer.pushContext(tagName);
            PlcResponseItem<PlcValue> valueResponse = valueEntry.getValue();
            if (!(valueResponse instanceof Serializable)) {
                throw new RuntimeException("Error serializing. PlcResponseItem doesn't implement Serializable");
            }
            ((Serializable) valueResponse).serialize(writeBuffer);
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("values");

        if (metadata != null && !metadata.isEmpty()) {
            writeBuffer.pushContext("metadata", WithRenderAsList(true));

            for (Entry<String, Metadata> entry : metadata.entrySet()) {
                if (entry.getValue() instanceof Serializable) {
                    writeBuffer.pushContext(entry.getKey());
                    ((Serializable) entry.getValue()).serialize(writeBuffer);
                    writeBuffer.popContext(entry.getKey());
                }
            }

            writeBuffer.popContext("metadata");
        }

        writeBuffer.popContext("PlcReadResponse");
    }

}
