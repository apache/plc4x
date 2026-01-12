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
import java.time.*;

public class PlcDATE_AND_LTIME extends PlcIECValue<LocalDateTime> {

    public static PlcDATE_AND_LTIME of(Object value) {
        if (value instanceof PlcDATE_AND_LTIME plcDATEAndLtime) {
            return plcDATEAndLtime;
        }
        if (value instanceof LocalDateTime localDateTime) {
            return new PlcDATE_AND_LTIME(localDateTime);
        }
        if (value instanceof Byte b) {
            return new PlcDATE_AND_LTIME(b);
        }
        if (value instanceof Short i) {
            return new PlcDATE_AND_LTIME(i);
        }
        if (value instanceof Integer i) {
            return new PlcDATE_AND_LTIME(i);
        }
        if (value instanceof Long l) {
            return new PlcDATE_AND_LTIME(l);
        }
        if (value instanceof Float v) {
            return new PlcDATE_AND_LTIME(v);
        }
        if (value instanceof Double v) {
            return new PlcDATE_AND_LTIME(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcDATE_AND_LTIME(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcDATE_AND_LTIME(bigDecimal);
        }
        return new PlcDATE_AND_LTIME(LocalDateTime.parse(value.toString()));
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
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        String valueString = value.toString();
        writeBuffer.writeString(valueString.getBytes(StandardCharsets.UTF_8).length * 8,
            valueString, WithOption.WithName(getClass().getSimpleName()), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
