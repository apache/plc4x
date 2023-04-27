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

import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
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

public class PlcNull implements PlcValue, Serializable {

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.NULL;
    }

    @Override
    public Object getObject() {
        return null;
    }

    @Override
    public boolean isSimple() {
        return true;
    }

    @Override
    public boolean isNullable() {
        return true;
    }

    @Override
    public boolean isNull() {
        return true;
    }

    @Override
    public boolean is(Class<?> clazz) {
        return false;
    }

    @Override
    public boolean isConvertibleTo(Class<?> clazz) {
        return false;
    }

    @Override
    public <T> T get(Class<T> clazz) {
        return null;
    }

    @Override
    public boolean isBoolean() {
        return false;
    }

    @Override
    public boolean getBoolean() {
        return false;
    }

    @Override
    public boolean isByte() {
        return false;
    }

    @Override
    public byte getByte() {
        return 0;
    }

    @Override
    public boolean isShort() {
        return false;
    }

    @Override
    public short getShort() {
        return 0;
    }

    @Override
    public boolean isInteger() {
        return false;
    }

    @Override
    public int getInteger() {
        return 0;
    }

    @Override
    public int getInt() {
        return 0;
    }

    @Override
    public boolean isLong() {
        return false;
    }

    @Override
    public long getLong() {
        return 0;
    }

    @Override
    public boolean isBigInteger() {
        return false;
    }

    @Override
    public BigInteger getBigInteger() {
        return null;
    }

    @Override
    public boolean isFloat() {
        return false;
    }

    @Override
    public float getFloat() {
        return 0;
    }

    @Override
    public boolean isDouble() {
        return false;
    }

    @Override
    public double getDouble() {
        return 0;
    }

    @Override
    public boolean isBigDecimal() {
        return false;
    }

    @Override
    public BigDecimal getBigDecimal() {
        return null;
    }

    @Override
    public boolean isString() {
        return false;
    }

    @Override
    public String getString() {
        return null;
    }

    @Override
    public boolean isDuration() {
        return false;
    }

    @Override
    public Duration getDuration() {
        return null;
    }

    @Override
    public boolean isTime() {
        return false;
    }

    @Override
    public LocalTime getTime() {
        return null;
    }

    @Override
    public boolean isDate() {
        return false;
    }

    @Override
    public LocalDate getDate() {
        return null;
    }

    @Override
    public boolean isDateTime() {
        return false;
    }

    @Override
    public LocalDateTime getDateTime() {
        return null;
    }

    @Override
    public byte[] getRaw() {
        return new byte[0];
    }

    @Override
    public boolean isList() {
        return false;
    }

    @Override
    public int getLength() {
        return 0;
    }

    @Override
    public PlcValue getIndex(int i) {
        return null;
    }

    @Override
    public List<? extends PlcValue> getList() {
        return null;
    }

    @Override
    public boolean isStruct() {
        return false;
    }

    @Override
    public Set<String> getKeys() {
        return null;
    }

    @Override
    public boolean hasKey(String key) {
        return false;
    }

    @Override
    public PlcValue getValue(String key) {
        return null;
    }

    @Override
    public Map<String, ? extends PlcValue> getStruct() {
        return null;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcNull");
        writeBuffer.popContext("PlcNull");
    }
}
