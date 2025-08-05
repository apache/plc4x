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

public class PlcLDATE_AND_TIME extends PlcIECValue<LocalDateTime> {

    public static PlcLDATE_AND_TIME of(Object value) {
        if (value instanceof PlcLDATE_AND_TIME) {
            return (PlcLDATE_AND_TIME) value;
        } else if (value instanceof LocalDateTime) {
            return new PlcLDATE_AND_TIME((LocalDateTime) value);
        } else if (value instanceof Byte) {
            return new PlcLDATE_AND_TIME((Byte) value);
        } else if (value instanceof Short) {
            return new PlcLDATE_AND_TIME((Short) value);
        } else if (value instanceof Integer) {
            return new PlcLDATE_AND_TIME((Integer) value);
        } else if (value instanceof Long) {
            return new PlcLDATE_AND_TIME((Long) value);
        } else if (value instanceof Float) {
            return new PlcLDATE_AND_TIME((Float) value);
        } else if (value instanceof Double) {
            return new PlcLDATE_AND_TIME((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcLDATE_AND_TIME((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcLDATE_AND_TIME((BigDecimal) value);
        } else {
            return new PlcLDATE_AND_TIME(LocalDateTime.parse(value.toString()));
        }
    }

    public static PlcLDATE_AND_TIME ofNanosecondsSinceEpoch(BigInteger nanosecondsSinceEpoch) {
        BigInteger epochSecond = nanosecondsSinceEpoch.divide(BigInteger.valueOf(1000_000));
        BigInteger nanoOfSecond = nanosecondsSinceEpoch.mod(BigInteger.valueOf(1000_000));
        return new PlcLDATE_AND_TIME(LocalDateTime.ofEpochSecond(epochSecond.longValue(), nanoOfSecond.intValue(),
            ZoneOffset.UTC));
    }

    public PlcLDATE_AND_TIME(Byte millisecondsSinceEpoch) {
        long daysSinceEpoch = 0;
        long millisecondsOfDay = millisecondsSinceEpoch.longValue();
        this.value = LocalDateTime.of(LocalDate.ofEpochDay(daysSinceEpoch), LocalTime.ofNanoOfDay(millisecondsOfDay * 1000000));
        this.isNullable = false;
    }

    public PlcLDATE_AND_TIME(Short millisecondsSinceEpoch) {
        long daysSinceEpoch = 0;
        long millisecondsOfDay = millisecondsSinceEpoch.longValue();
        this.value = LocalDateTime.of(LocalDate.ofEpochDay(daysSinceEpoch), LocalTime.ofNanoOfDay(millisecondsOfDay * 1000000));
        this.isNullable = false;
    }

    public PlcLDATE_AND_TIME(Integer millisecondsSinceEpoch) {
        long daysSinceEpoch = millisecondsSinceEpoch.longValue() / 86400000;
        long millisecondsOfDay = millisecondsSinceEpoch.longValue() % 86400000;
        this.value = LocalDateTime.of(LocalDate.ofEpochDay(daysSinceEpoch), LocalTime.ofNanoOfDay(millisecondsOfDay * 1000000));
        this.isNullable = false;
    }

    public PlcLDATE_AND_TIME(Long millisecondsSinceEpoch) {
        long daysSinceEpoch = millisecondsSinceEpoch / 86400000;
        long millisecondsOfDay = millisecondsSinceEpoch % 86400000;
        this.value = LocalDateTime.of(LocalDate.ofEpochDay(daysSinceEpoch), LocalTime.ofNanoOfDay(millisecondsOfDay * 1000000));
        this.isNullable = false;
    }

    public PlcLDATE_AND_TIME(Float millisecondsSinceEpoch) {
        long daysSinceEpoch = millisecondsSinceEpoch.longValue() / 86400000;
        long millisecondsOfDay = millisecondsSinceEpoch.longValue() % 86400000;
        this.value = LocalDateTime.of(LocalDate.ofEpochDay(daysSinceEpoch), LocalTime.ofNanoOfDay(millisecondsOfDay * 1000000));
        this.isNullable = false;
    }

    public PlcLDATE_AND_TIME(Double millisecondsSinceEpoch) {
        long daysSinceEpoch = millisecondsSinceEpoch.longValue() / 86400000;
        long millisecondsOfDay = millisecondsSinceEpoch.longValue() % 86400000;
        this.value = LocalDateTime.of(LocalDate.ofEpochDay(daysSinceEpoch), LocalTime.ofNanoOfDay(millisecondsOfDay * 1000000));
        this.isNullable = false;
    }

    public PlcLDATE_AND_TIME(BigInteger millisecondsSinceEpoch) {
        long daysSinceEpoch = millisecondsSinceEpoch.longValue() / 86400000;
        long millisecondsOfDay = millisecondsSinceEpoch.longValue() % 86400000;
        this.value = LocalDateTime.of(LocalDate.ofEpochDay(daysSinceEpoch), LocalTime.ofNanoOfDay(millisecondsOfDay * 1000000));
        this.isNullable = false;
    }

    public PlcLDATE_AND_TIME(BigDecimal millisecondsSinceEpoch) {
        long daysSinceEpoch = millisecondsSinceEpoch.longValue() / 86400000;
        long millisecondsOfDay = millisecondsSinceEpoch.longValue() % 86400000;
        this.value = LocalDateTime.of(LocalDate.ofEpochDay(daysSinceEpoch), LocalTime.ofNanoOfDay(millisecondsOfDay * 1000000));
        this.isNullable = false;
    }

    public PlcLDATE_AND_TIME(LocalDateTime value) {
        this.value = value;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.DATE_AND_TIME;
    }

    public BigInteger getNanosecondsSinceEpoch() {
        Instant instant = getDateTime().toInstant(ZoneOffset.UTC);
        return BigInteger.valueOf(instant.getEpochSecond()).multiply(BigInteger.valueOf(1000_000_000)).add(BigInteger.valueOf(instant.getNano()));
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
