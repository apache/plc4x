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

public class PlcLDATE extends PlcIECValue<LocalDate> {

    public static PlcLDATE of(Object value) {
        if (value instanceof PlcLDATE plcLDATE) {
            return plcLDATE;
        }
        if (value instanceof LocalDate localDate) {
            return new PlcLDATE(localDate);
        }
        if (value instanceof Byte b) {
            return new PlcLDATE(b);
        }
        if (value instanceof Short i) {
            return new PlcLDATE(i);
        }
        if (value instanceof Integer i) {
            return new PlcLDATE(i);
        }
        if (value instanceof Long l) {
            return new PlcLDATE(l);
        }
        if (value instanceof Float v) {
            return new PlcLDATE(v);
        }
        if (value instanceof Double v) {
            return new PlcLDATE(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcLDATE(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcLDATE(bigDecimal);
        }
        return new PlcLDATE(LocalDate.parse(value.toString()));
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
        return PlcValueType.LDATE;
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
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        String valueString = value.toString();
        writeBuffer.writeString(valueString.getBytes(StandardCharsets.UTF_8).length * 8,
            valueString, WithOption.WithName(getClass().getSimpleName()), WithOption.WithEncoding("UTF8"));
    }

}
