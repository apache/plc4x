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

public class PlcDATE_AND_LTIME extends PlcIECValue<LocalDateTime> {

    public static PlcDATE_AND_LTIME of(Object value) {
        if (value instanceof PlcDATE_AND_LTIME) {
            return (PlcDATE_AND_LTIME) value;
        } else if (value instanceof LocalDateTime) {
            return new PlcDATE_AND_LTIME((LocalDateTime) value);
        } else if (value instanceof Byte) {
            return new PlcDATE_AND_LTIME((Byte) value);
        } else if (value instanceof Short) {
            return new PlcDATE_AND_LTIME((Short) value);
        } else if (value instanceof Integer) {
            return new PlcDATE_AND_LTIME((Integer) value);
        } else if (value instanceof Long) {
            return new PlcDATE_AND_LTIME((Long) value);
        } else if (value instanceof Float) {
            return new PlcDATE_AND_LTIME((Float) value);
        } else if (value instanceof Double) {
            return new PlcDATE_AND_LTIME((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcDATE_AND_LTIME((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcDATE_AND_LTIME((BigDecimal) value);
        } else {
            return new PlcDATE_AND_LTIME(LocalDateTime.parse(value.toString()));
        }
    }

    public PlcDATE_AND_LTIME(Byte nanosecondsSinceEpoch) {
        long secondsSinceEpoch = 0;
        long nannoSecondsOfSecond = nanosecondsSinceEpoch % 1000000000;
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, (int) nannoSecondsOfSecond,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_LTIME(Short nanosecondsSinceEpoch) {
        long secondsSinceEpoch = 0;
        long nannoSecondsOfSecond = nanosecondsSinceEpoch % 1000000000;
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, (int) nannoSecondsOfSecond,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_LTIME(Integer nanosecondsSinceEpoch) {
        long secondsSinceEpoch = nanosecondsSinceEpoch.longValue() / 1000000000;
        long nannoSecondsOfSecond = nanosecondsSinceEpoch % 1000000000;
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, (int) nannoSecondsOfSecond,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_LTIME(Long nanosecondsSinceEpoch) {
        long secondsSinceEpoch = nanosecondsSinceEpoch / 1000000000;
        long nannoSecondsOfSecond = nanosecondsSinceEpoch % 1000000000;
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, (int) nannoSecondsOfSecond,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_LTIME(Float nanosecondsSinceEpoch) {
        long secondsSinceEpoch = nanosecondsSinceEpoch.longValue() / 1000000000;
        long nannoSecondsOfSecond = nanosecondsSinceEpoch.longValue() % 1000000000;
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, (int) nannoSecondsOfSecond,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_LTIME(Double nanosecondsSinceEpoch) {
        long secondsSinceEpoch = nanosecondsSinceEpoch.longValue() / 1000000000;
        long nannoSecondsOfSecond = nanosecondsSinceEpoch.longValue() % 1000000000;
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, (int) nannoSecondsOfSecond,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_LTIME(BigInteger nanosecondsSinceEpoch) {
        long secondsSinceEpoch = nanosecondsSinceEpoch.longValue() / 1000000000;
        long nannoSecondsOfSecond = nanosecondsSinceEpoch.longValue() % 1000000000;
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, (int) nannoSecondsOfSecond,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_LTIME(BigDecimal nanosecondsSinceEpoch) {
        long secondsSinceEpoch = nanosecondsSinceEpoch.longValue() / 1000000000;
        long nannoSecondsOfSecond = nanosecondsSinceEpoch.longValue() % 1000000000;
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, (int) nannoSecondsOfSecond,
            ZoneOffset.UTC);
        this.isNullable = false;
    }

    public PlcDATE_AND_LTIME(LocalDateTime value) {
        this.value = value;
        this.isNullable = false;
    }


    public static PlcDATE_AND_LTIME ofNanosecondsSinceEpoch(BigInteger nanosecondsSinceEpoch) {
        return new PlcDATE_AND_LTIME(nanosecondsSinceEpoch.longValue());
    }

    public static PlcDATE_AND_LTIME ofSegments(int year, int month, int day, int hour, int minutes, int seconds, long nannosecondsOfSecond) {
        return new PlcDATE_AND_LTIME(LocalDateTime.of(year, month, day, hour, minutes, seconds, (int) nannosecondsOfSecond));
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.DATE_AND_LTIME;
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
