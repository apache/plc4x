/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.values;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.ZoneOffset;

public class PlcDATE extends PlcIECValue<LocalDate> {

    public static PlcDATE of(Object value) {
        if (value instanceof PlcDATE plcDATE) {
            return plcDATE;
        }
        if (value instanceof LocalDate localDate) {
            return new PlcDATE(localDate);
        }
        if (value instanceof Byte b) {
            return new PlcDATE(b);
        }
        if (value instanceof Short i) {
            return new PlcDATE(i);
        }
        if (value instanceof Integer i) {
            return new PlcDATE(i);
        }
        if (value instanceof Long l) {
            return new PlcDATE(l);
        }
        if (value instanceof Float v) {
            return new PlcDATE(v);
        }
        if (value instanceof Double v) {
            return new PlcDATE(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcDATE(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcDATE(bigDecimal);
        }
        return new PlcDATE(LocalDate.parse(value.toString()));
    }

    public static PlcDATE ofSecondsSinceEpoch(long secondsSinceEpoch) {
        return new PlcDATE(LocalDateTime.ofInstant(
            Instant.ofEpochSecond(secondsSinceEpoch), ZoneOffset.UTC).toLocalDate());
    }

    public static PlcDATE ofDaysSinceEpoch(int daysSinceEpoch) {
        // 86400 = 24 hours x 60 Minutes x 60 Seconds
        return new PlcDATE(LocalDateTime.ofInstant(
            Instant.ofEpochSecond(((long) daysSinceEpoch) * 86400), ZoneOffset.UTC).toLocalDate());
    }

    /**
     * Creates a PlcDATE from explicit date components.
     * Used by protocols that transmit date as individual BCD-encoded fields
     * (e.g., Schneider UMAS).
     *
     * @param year  the year (e.g., 2026)
     * @param month the month (1-12)
     * @param day   the day of month (1-31)
     * @return the PlcDATE value
     */
    public static PlcDATE ofSegments(int year, int month, int day) {
        return new PlcDATE(LocalDate.of(year, month, day));
    }

    public PlcDATE(LocalDate value) {
        this.value = value;
        this.isNullable = false;
    }

    public PlcDATE(Byte daysSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofInstant(
            Instant.ofEpochSecond(((long) daysSinceEpoch) * 86400), ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcDATE(Short daysSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofInstant(
            Instant.ofEpochSecond(((long) daysSinceEpoch) * 86400), ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcDATE(Integer daysSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofInstant(
            Instant.ofEpochSecond(((long) daysSinceEpoch) * 86400), ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcDATE(Long daysSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofInstant(
            Instant.ofEpochSecond((daysSinceEpoch) * 86400), ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcDATE(Float daysSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofInstant(
            Instant.ofEpochSecond(daysSinceEpoch.longValue() * 86400), ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcDATE(Double daysSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofInstant(
            Instant.ofEpochSecond(daysSinceEpoch.longValue() * 86400), ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcDATE(BigInteger daysSinceEpoch) {
        this.value = LocalDateTime.ofInstant(
            Instant.ofEpochSecond(daysSinceEpoch.longValue() * 86400), ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcDATE(BigDecimal daysSinceEpoch) {
        this.value = LocalDateTime.ofInstant(
            Instant.ofEpochSecond(daysSinceEpoch.longValue() * 86400), ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.DATE;
    }

    public long getSecondsSinceEpoch() {
        return value.atStartOfDay(ZoneOffset.UTC).toEpochSecond();
    }

    public int getDaysSinceEpoch() {
        return (int) (value.atStartOfDay(ZoneOffset.UTC).toEpochSecond() / 86400);
    }

    @Override
    public boolean isLong() {
        return true;
    }

    @Override
    public long getLong() {
        Instant instant = value.atStartOfDay(ZoneOffset.UTC).toInstant();
        return (instant.toEpochMilli() / 1000);
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
    public boolean isDate() {
        return true;
    }

    @Override
    public LocalDate getDate() {
        return value;
    }

    @Override
    public LocalDateTime getDateTime() {
        return value.atStartOfDay();
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        String valueString = value.toString();
        writeBuffer.writeString(valueString.getBytes(StandardCharsets.UTF_8).length * 8,
            valueString, WithOption.WithName(getClass().getSimpleName()), WithOption.WithEncoding("UTF8"));
    }

}
