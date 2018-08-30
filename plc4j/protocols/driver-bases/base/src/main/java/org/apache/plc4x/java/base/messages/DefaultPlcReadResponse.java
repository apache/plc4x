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
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collection;
import java.util.Map;

public class DefaultPlcReadResponse implements PlcReadResponse {

    private final PlcReadRequest request;
    private final Map<String, Pair<PlcResponseCode, byte[][]>> values;

    public DefaultPlcReadResponse(PlcReadRequest request, Map<String, Pair<PlcResponseCode, byte[][]>> values) {
        this.request = request;
        this.values = values;
    }

    @Override
    public PlcReadRequest getRequest() {
        return request;
    }

    @Override
    public int getNumValues(String name) {
        if(values.get(name) == null) {
            return 0;
        }
        if(values.get(name).getKey() != PlcResponseCode.OK) {
            return 0;
        }
        return values.get(name).getValue().length;
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
        if(values.get(name) == null) {
            return null;
        }
        return values.get(name).getKey();
    }

    @Override
    public boolean isRaw(String name) {
        return false;
    }

    @Override
    public byte[] getRaw(String name) {
        return getRaw(name, 0);
    }

    @Override
    public byte[] getRaw(String name, int index) {
        return new byte[0];
    }

    @Override
    public Collection<byte[]> getAllRaws(String name) {
        return null;
    }

    @Override
    public boolean isObject(String name) {
        return false;
    }

    @Override
    public Object getObject(String name) {
        return getObject(name, 0);
    }

    @Override
    public Object getObject(String name, int index) {
        return null;
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        return null;
    }

    @Override
    public boolean isBoolean(String name) {
        return false;
    }

    @Override
    public Boolean getBoolean(String name) {
        return getBoolean(name, 0);
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        return null;
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        return null;
    }

    @Override
    public boolean isByte(String name) {
        return false;
    }

    @Override
    public Byte getByte(String name) {
        return getByte(name, 0);
    }

    @Override
    public Byte getByte(String name, int index) {
        return null;
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        return null;
    }

    @Override
    public boolean isShort(String name) {
        return false;
    }

    @Override
    public Short getShort(String name) {
        return getShort(name, 0);
    }

    @Override
    public Short getShort(String name, int index) {
        return null;
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        return null;
    }

    @Override
    public boolean isInteger(String name) {
        return false;
    }

    @Override
    public Integer getInteger(String name) {
        return getInteger(name, 0);
    }

    @Override
    public Integer getInteger(String name, int index) {
        return null;
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        return null;
    }

    @Override
    public boolean isLong(String name) {
        return false;
    }

    @Override
    public Long getLong(String name) {
        return getLong(name, 0);
    }

    @Override
    public Long getLong(String name, int index) {
        return null;
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        return null;
    }

    @Override
    public boolean isFloat(String name) {
        return false;
    }

    @Override
    public Float getFloat(String name) {
        return getFloat(name, 0);
    }

    @Override
    public Float getFloat(String name, int index) {
        return null;
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        return null;
    }

    @Override
    public boolean isDouble(String name) {
        return false;
    }

    @Override
    public Double getDouble(String name) {
        return getDouble(name, 0);
    }

    @Override
    public Double getDouble(String name, int index) {
        return null;
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        return null;
    }

    @Override
    public boolean isString(String name) {
        return false;
    }

    @Override
    public String getString(String name) {
        return getString(name, 0);
    }

    @Override
    public String getString(String name, int index) {
        return null;
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        return null;
    }

    @Override
    public boolean isTime(String name) {
        return false;
    }

    @Override
    public LocalTime getTime(String name) {
        return getTime(name, 0);
    }

    @Override
    public LocalTime getTime(String name, int index) {
        return null;
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        return null;
    }

    @Override
    public boolean isDate(String name) {
        return false;
    }

    @Override
    public LocalDate getDate(String name) {
        return getDate(name, 0);
    }

    @Override
    public LocalDate getDate(String name, int index) {
        return null;
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        return null;
    }

    @Override
    public boolean isDateTime(String name) {
        return false;
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        return getDateTime(name, 0);
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        return null;
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        return null;
    }

}
