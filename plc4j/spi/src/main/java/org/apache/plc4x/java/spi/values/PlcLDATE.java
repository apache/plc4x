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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.*;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcLDATE extends PlcSimpleValue<LocalDate> {

    public static PlcLDATE of(Object value) {
        if (value instanceof LocalDate) {
            return new PlcLDATE((LocalDate) value);
        } else if (value instanceof Long) {
            return new PlcLDATE(LocalDateTime.ofInstant(
                Instant.ofEpochSecond((long) value), ZoneId.systemDefault()).toLocalDate());
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    public static PlcLDATE ofNanosecondsSinceEpoch(BigInteger nanosecondsSinceEpoch) {
        BigInteger epochSecond = nanosecondsSinceEpoch.divide(BigInteger.valueOf(1000_000));
        BigInteger nanoOfSecond = nanosecondsSinceEpoch.mod(BigInteger.valueOf(1000_000));
        return new PlcLDATE(LocalDateTime.ofEpochSecond(epochSecond.longValue(), nanoOfSecond.intValue(),
            ZoneOffset.of(ZoneOffset.systemDefault().getId())).toLocalDate());
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcLDATE(@JsonProperty("value") LocalDate value) {
        super(value, true);
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcLDATE(@JsonProperty("value") BigInteger nanosecondsSinceEpoch) {
        super(LocalDateTime.ofEpochSecond(nanosecondsSinceEpoch.longValue() / 1000000,
            (int) (nanosecondsSinceEpoch.longValue() % 1000000),
            ZoneOffset.of(ZoneOffset.systemDefault().getId())).toLocalDate(), true);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.DATE;
    }

    public BigInteger getNanosecondsSinceEpoch() {
        Instant instant = getDateTime().toInstant(ZoneOffset.of(ZoneOffset.systemDefault().getId()));
        return BigInteger.valueOf(instant.getEpochSecond()).multiply(BigInteger.valueOf(1000_000_000)).add(BigInteger.valueOf(instant.getNano()));
    }

    @Override
    public boolean isLong() {
        return true;
    }

    @Override
    public long getLong() {
        Instant instant = value.atStartOfDay(ZoneId.systemDefault()).toInstant();
        return (instant.toEpochMilli() / 1000);
    }

    @Override
    @JsonIgnore
    public boolean isString() {
        return true;
    }

    @Override
    @JsonIgnore
    public String getString() {
        return value.toString();
    }

    @Override
    @JsonIgnore
    public boolean isDate() {
        return true;
    }

    @Override
    @JsonIgnore
    public LocalDate getDate() {
        return value;
    }

    @Override
    @JsonIgnore
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
