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

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collection;

public class DefaultPlcSubscriptionEvent implements PlcSubscriptionEvent {

    @Override
    public int getNumberOfValues(String name) {
        return 0;
    }

    @Override
    public boolean isValidObject(String name) {
        return false;
    }

    @Override
    public boolean isValidObject(String name, int index) {
        return false;
    }

    @Override
    public Object getObject(String name) {
        return null;
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
    public boolean isValidBoolean(String name) {
        return false;
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        return false;
    }

    @Override
    public Boolean getBoolean(String name) {
        return null;
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
    public boolean isValidByte(String name) {
        return false;
    }

    @Override
    public boolean isValidByte(String name, int index) {
        return false;
    }

    @Override
    public Byte getByte(String name) {
        return null;
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
    public byte[] getByteArray(String name) {
        return new byte[0];
    }

    @Override
    public boolean isValidShort(String name) {
        return false;
    }

    @Override
    public boolean isValidShort(String name, int index) {
        return false;
    }

    @Override
    public Short getShort(String name) {
        return null;
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
    public boolean isValidInteger(String name) {
        return false;
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        return false;
    }

    @Override
    public Integer getInteger(String name) {
        return null;
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
    public boolean isValidLong(String name) {
        return false;
    }

    @Override
    public boolean isValidLong(String name, int index) {
        return false;
    }

    @Override
    public Long getLong(String name) {
        return null;
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
    public boolean isValidFloat(String name) {
        return false;
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        return false;
    }

    @Override
    public Float getFloat(String name) {
        return null;
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
    public boolean isValidDouble(String name) {
        return false;
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        return false;
    }

    @Override
    public Double getDouble(String name) {
        return null;
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
    public boolean isValidString(String name) {
        return false;
    }

    @Override
    public boolean isValidString(String name, int index) {
        return false;
    }

    @Override
    public String getString(String name) {
        return null;
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
    public boolean isValidTime(String name) {
        return false;
    }

    @Override
    public boolean isValidTime(String name, int index) {
        return false;
    }

    @Override
    public LocalTime getTime(String name) {
        return null;
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
    public boolean isValidDate(String name) {
        return false;
    }

    @Override
    public boolean isValidDate(String name, int index) {
        return false;
    }

    @Override
    public LocalDate getDate(String name) {
        return null;
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
    public boolean isValidDateTime(String name) {
        return false;
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        return false;
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        return null;
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        return null;
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        return null;
    }

    @Override
    public Collection<String> getFieldNames() {
        return null;
    }

    @Override
    public PlcField getField(String name) {
        return null;
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        return null;
    }

    @Override
    public PlcReadRequest getRequest() {
        return null;
    }

}
