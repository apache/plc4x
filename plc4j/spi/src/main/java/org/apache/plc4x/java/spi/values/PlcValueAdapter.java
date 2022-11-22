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
package org.apache.plc4x.java.spi.values;

import com.fasterxml.jackson.annotation.JsonIgnore;
import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;
import java.util.Map;
import java.util.Set;

public abstract class PlcValueAdapter implements PlcValue, Serializable {

    @Override
    public Object getObject() {
        return null;
    }

    @Override
    @JsonIgnore
    public boolean isSimple() {
        return false;
    }

    @Override
    @JsonIgnore
    public boolean isNullable() {
        return false;
    }

    @Override
    @JsonIgnore
    public boolean isNull() {
        return false;
    }

    @Override
    @JsonIgnore
    public boolean is(Class<?> clazz) {
        return false;
    }

    @Override
    @JsonIgnore
    public boolean isConvertibleTo(Class<?> clazz) {
        return false;
    }

    @Override
    @JsonIgnore
    public <T> T get(Class<T> clazz) {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isBoolean() {
        return false;
    }

    @Override
    @JsonIgnore
    public boolean getBoolean() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isByte() {
        return false;
    }

    @Override
    @JsonIgnore
    public byte getByte() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isShort() {
        return false;
    }

    @Override
    @JsonIgnore
    public short getShort() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isInteger() {
        return false;
    }

    @Override
    @JsonIgnore
    public int getInteger() {
        throw new PlcIncompatibleDatatypeException("");
    }

    /**
     * int vs. Integer is the only primitive type where the complex type doesn't have just a capitalized first letter.
     * @return well ... an Integer ... ahem ... int ...
     */
    @Override
    @JsonIgnore
    public int getInt() {
        return getInteger();
    }

    @Override
    @JsonIgnore
    public boolean isLong() {
        return false;
    }

    @Override
    @JsonIgnore
    public long getLong() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isBigInteger() {
        return false;
    }

    @Override
    @JsonIgnore
    public BigInteger getBigInteger() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isFloat() {
        return false;
    }

    @Override
    @JsonIgnore
    public float getFloat() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isDouble() {
        return false;
    }

    @Override
    @JsonIgnore
    public double getDouble() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isBigDecimal() {
        return false;
    }

    @Override
    @JsonIgnore
    public BigDecimal getBigDecimal() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isString() {
        return false;
    }

    @Override
    @JsonIgnore
    public String getString() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isDuration() {
        return false;
    }

    @Override
    @JsonIgnore
    public Duration getDuration() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isTime() {
        return false;
    }

    @Override
    @JsonIgnore
    public LocalTime getTime() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isDate() {
        return false;
    }

    @Override
    @JsonIgnore
    public LocalDate getDate() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isDateTime() {
        return false;
    }

    @Override
    @JsonIgnore
    public LocalDateTime getDateTime() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public byte[] getRaw() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isList() {
        return false;
    }

    @Override
    @JsonIgnore
    public int getLength() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public PlcValue getIndex(int i) {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public List<? extends PlcValue> getList() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean isStruct() {
        return false;
    }

    @Override
    @JsonIgnore
    public Set<String> getKeys() {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public boolean hasKey(String key) {
        return false;
    }

    @Override
    @JsonIgnore
    public PlcValue getValue(String key) {
        throw new PlcIncompatibleDatatypeException("");
    }

    @Override
    @JsonIgnore
    public Map<String, ? extends PlcValue> getStruct() {
        throw new PlcIncompatibleDatatypeException("");
    }

}
