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
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneOffset;

public class PlcDATE_AND_TIME extends PlcIECValue<LocalDateTime> {

    public static PlcDATE_AND_TIME of(Object value) {
        if (value instanceof PlcDATE_AND_TIME) {
            return (PlcDATE_AND_TIME) value;
        } else if (value instanceof LocalDateTime) {
            return new PlcDATE_AND_TIME((LocalDateTime) value);
        } else if (value instanceof Byte) {
            return new PlcDATE_AND_TIME((Byte) value);
        } else if (value instanceof Short) {
            return new PlcDATE_AND_TIME((Short) value);
        } else if (value instanceof Integer) {
            return new PlcDATE_AND_TIME((Integer) value);
        } else if (value instanceof Long) {
            return new PlcDATE_AND_TIME((Long) value);
        } else if (value instanceof Float) {
            return new PlcDATE_AND_TIME((Float) value);
        } else if (value instanceof Double) {
            return new PlcDATE_AND_TIME((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcDATE_AND_TIME((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcDATE_AND_TIME((BigDecimal) value);
        } else {
            return new PlcDATE_AND_TIME(LocalDateTime.parse(value.toString()));
        }
    }

    public static PlcDATE_AND_TIME ofSecondsSinceEpoch(long secondsSinceEpoch) {
        return new PlcDATE_AND_TIME(LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0,
            ZoneOffset.UTC));
    }

    public static PlcDATE_AND_TIME ofSegments(int year, int month, int day, int hour, int minutes, int seconds, int nanoseconds) {
        return new PlcDATE_AND_TIME(LocalDateTime.of(year, month, day, hour, minutes, seconds, nanoseconds));
    }

    public PlcDATE_AND_TIME(Byte secondsSinceEpoch) {
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(Short secondsSinceEpoch) {
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(Integer secondsSinceEpoch) {
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(Long secondsSinceEpoch) {
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(Float secondsSinceEpoch) {
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch.longValue(), 0,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(Double secondsSinceEpoch) {
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch.longValue(), 0,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(BigInteger secondsSinceEpoch) {
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch.longValue(), 0,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(BigDecimal secondsSinceEpoch) {
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch.longValue(), 0,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(LocalDateTime value) {
        this.value = value;
        this.isNullable = false;
    }

    public PlcDATE_AND_TIME(int year, int month, int day, int hour, int minutes, int seconds, int nanoseconds) {
        this.value = LocalDateTime.of(year, month, day, hour, minutes, seconds, nanoseconds);
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.DATE_AND_TIME;
    }

    public long getSecondsSinceEpoch() {
        Instant instant = getDateTime().toInstant(ZoneOffset.UTC);
        return instant.getEpochSecond();
    }

    public int getYear() {
        return value.getYear();
    }

    public int getMonth() {
        return value.getMonthValue();
    }

    public int getDay() {
        return value.getDayOfMonth();
    }

    public int getDayOfWeek() {
        return value.getDayOfWeek().getValue();
    }

    public int getHour() {
        return value.getHour();
    }

    public int getMinutes() {
        return value.getMinute();
    }

    public int getSeconds() {
        return value.getSecond();
    }

    public int getNanoseconds() {
        return value.getNano();
    }

    @Override
    public boolean isLong() {
        return true;
    }

    @Override
    public long getLong() {
        Instant instant = value.atZone(ZoneOffset.UTC).toInstant();
        return instant.getEpochSecond();
    }

    @Override
    public boolean isString() {
        return true;
    }

    @Override
    public String getString() {
        return value.toString();
    }

    @Override
    public boolean isTime() {
        return true;
    }

    @Override
    public LocalTime getTime() {
        return value.toLocalTime();
    }

    @Override
    public boolean isDate() {
        return true;
    }

    @Override
    public LocalDate getDate() {
        return value.toLocalDate();
    }

    @Override
    public boolean isDateTime() {
        return true;
    }

    @Override
    public LocalDateTime getDateTime() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        String valueString = value.toString();
        writeBuffer.writeString(getClass().getSimpleName(),
            valueString.getBytes(StandardCharsets.UTF_8).length*8,
            valueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
