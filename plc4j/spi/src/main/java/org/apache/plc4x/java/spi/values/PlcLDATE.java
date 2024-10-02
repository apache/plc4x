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
import java.time.ZoneOffset;

public class PlcLDATE extends PlcIECValue<LocalDate> {

    public static PlcLDATE of(Object value) {
        if (value instanceof PlcLDATE) {
            return (PlcLDATE) value;
        } else if (value instanceof LocalDate) {
            return new PlcLDATE((LocalDate) value);
        } else if (value instanceof Byte) {
            return new PlcLDATE((Byte) value);
        } else if (value instanceof Short) {
            return new PlcLDATE((Short) value);
        } else if (value instanceof Integer) {
            return new PlcLDATE((Integer) value);
        } else if (value instanceof Long) {
            return new PlcLDATE((Long) value);
        } else if (value instanceof Float) {
            return new PlcLDATE((Float) value);
        } else if (value instanceof Double) {
            return new PlcLDATE((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcLDATE((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcLDATE((BigDecimal) value);
        } else {
            return new PlcLDATE(LocalDate.parse(value.toString()));
        }
    }

    public static PlcLDATE ofNanosecondsSinceEpoch(BigInteger nanosecondsSinceEpoch) {
        BigInteger epochSecond = nanosecondsSinceEpoch.divide(BigInteger.valueOf(1000_000));
        BigInteger nanoOfSecond = nanosecondsSinceEpoch.mod(BigInteger.valueOf(1000_000));
        return new PlcLDATE(LocalDateTime.ofEpochSecond(epochSecond.longValue(), nanoOfSecond.intValue(),
            ZoneOffset.UTC).toLocalDate());
    }

    public PlcLDATE(LocalDate value) {
        this.value = value;
        this.isNullable = false;
    }

    public PlcLDATE(Byte secondsSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofEpochSecond((long) secondsSinceEpoch, 0,
            ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcLDATE(Short secondsSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofEpochSecond((long) secondsSinceEpoch, 0,
            ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcLDATE(Integer secondsSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofEpochSecond((long) secondsSinceEpoch, 0,
            ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcLDATE(Long secondsSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofEpochSecond((long) secondsSinceEpoch, 0,
            ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcLDATE(Float secondsSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch.longValue(), 0,
            ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcLDATE(Double secondsSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch.longValue(), 0,
            ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcLDATE(BigInteger secondsSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch.longValue(), 0,
            ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    public PlcLDATE(BigDecimal secondsSinceEpoch) {
        // REMARK: Yes, I'm using LocalDataTime.ofInstant as LocalDate.ofInstant is marked "JDK 1.9"
        this.value = LocalDateTime.ofEpochSecond(secondsSinceEpoch.longValue(), 0,
            ZoneOffset.UTC).toLocalDate();
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.DATE;
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
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        String valueString = value.toString();
        writeBuffer.writeString(getClass().getSimpleName(),
            valueString.getBytes(StandardCharsets.UTF_8).length * 8,
            valueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
