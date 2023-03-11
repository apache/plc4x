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
package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.types.PlcValueType;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * Base Type of all Types.
 */
public interface PlcValue {

    PlcValueType getPlcValueType();
    
    Object getObject();

    // Simple Types

    boolean isSimple();

    boolean isNullable();

    boolean isNull();

    // Generic (\o/ Sebastian)

    boolean is(Class<?> clazz);

    boolean isConvertibleTo(Class<?> clazz);

    <T> T get(Class<T> clazz);

    // Boolean

    boolean isBoolean();

    boolean getBoolean();

    // Integer

    boolean isByte();

    byte getByte();

    boolean isShort();

    short getShort();

    boolean isInteger();

    int getInteger();

    int getInt();

    boolean isLong();

    long getLong();

    boolean isBigInteger();

    BigInteger getBigInteger();

    // Floating Point

    boolean isFloat();

    float getFloat();

    boolean isDouble();

    double getDouble();

    boolean isBigDecimal();

    BigDecimal getBigDecimal();

    // String

    boolean isString();

    String getString();

    // Time

    boolean isDuration();

    Duration getDuration();

    boolean isTime();

    LocalTime getTime();

    boolean isDate();

    LocalDate getDate();

    boolean isDateTime();

    LocalDateTime getDateTime();

    // Raw Access

    byte[] getRaw();

    // List Methods

    boolean isList();

    int getLength();

    PlcValue getIndex(int i);

    List<? extends PlcValue> getList();

    // Struct Methods

    boolean isStruct();

    Set<String> getKeys();

    boolean hasKey(String key);

    PlcValue getValue(String key);

    Map<String, ? extends PlcValue> getStruct();

}
