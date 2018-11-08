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

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collection;

// TODO: FIXME: implement me
public class DefaultPlcSubscriptionEvent implements InternalPlcSubscriptionEvent {

    public final Instant timestamp;

    public final byte[] bytes;

    public DefaultPlcSubscriptionEvent(Instant timestamp, byte[] bytes) {
        this.timestamp = timestamp;
        this.bytes = bytes;
    }

    @Override
    public int getNumberOfValues(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Object getObject(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Object getObject(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public boolean isValidBoolean(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Boolean getBoolean(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Byte getByte(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Short getShort(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Integer getInteger(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public boolean isValidBigInteger(String name) {
        return false;
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        return false;
    }

    @Override
    public BigInteger getBigInteger(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<BigInteger> getAllBigIntegers(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Long getLong(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Float getFloat(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Double getDouble(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public boolean isValidBigDecimal(String name) {
        return false;
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        return false;
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public String getString(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public LocalTime getTime(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public LocalDate getDate(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        throw new NotImplementedException("TODO: implement me");
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
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public boolean isValidByteArray(String name) {
        return false;
    }

    @Override
    public boolean isValidByteArray(String name, int index) {
        return false;
    }

    @Override
    public Byte[] getByteArray(String name) {
        return new Byte[0];
    }

    @Override
    public Byte[] getByteArray(String name, int index) {
        return new Byte[0];
    }

    @Override
    public Collection<Byte[]> getAllByteArrays(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Collection<String> getFieldNames() {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public PlcField getField(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public PlcReadRequest getRequest() {
        throw new NotImplementedException("TODO: implement me");
    }

    @Override
    public Instant getTimestamp() {
        return timestamp;
    }

    @Override
    public byte[] getBytes() {
        return bytes;
    }
}
