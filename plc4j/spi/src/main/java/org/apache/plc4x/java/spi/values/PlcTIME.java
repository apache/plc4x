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

import java.math.BigDecimal;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.Period;
import java.time.format.DateTimeParseException;
import java.time.temporal.ChronoUnit;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcTIME extends PlcSimpleValue<Duration> {

    public static PlcTIME of(Object value) {
        if (value instanceof Duration) {
            return new PlcTIME((Duration) value);
        } else if (value instanceof Integer) {
            return new PlcTIME(Duration.of((long) value, ChronoUnit.MILLIS));
        } else if (value instanceof Long) {
            return new PlcTIME(Duration.of((long) value, ChronoUnit.MILLIS));
        } else if (value instanceof Number) {
            return new PlcTIME(((Number) value).longValue());
        } else if (value instanceof String) {
            try {
                return new PlcTIME(Duration.parse((String) value));
            } catch (DateTimeParseException e) {
                try {
                    return new PlcTIME(Period.parse((String) value).get(ChronoUnit.MILLIS));
                }catch (DateTimeParseException e1)
                {
                    return new PlcTIME(new BigDecimal((String) value).longValue());
                }
            }
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    public static PlcTIME ofMilliseconds(long milliseconds) {
        return new PlcTIME(Duration.ofMillis(milliseconds));
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcTIME(@JsonProperty("value") Duration value) {
        super(value, true);
    }

    public PlcTIME(@JsonProperty("value") long milliseconds) {
        super(Duration.ofMillis(milliseconds), true);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.TIME;
    }

    public long getMilliseconds() {
        return value.toMillis();
    }

    @Override
    public boolean isInteger() {
        return true;
    }

    @Override
    public boolean isLong() {
        return true;
    }

    @Override
    @JsonIgnore
    public boolean isString() {
        return true;
    }

    @Override
    public boolean isDuration() {
        return true;
    }

    @Override
    public int getInteger() {
        return (int) (value.get(ChronoUnit.NANOS) / 1000000);
    }

    @Override
    public long getLong() {
        return value.get(ChronoUnit.NANOS) / 1000000;
    }

    @Override
    public Duration getDuration() {
        return value;
    }

    @Override
    @JsonIgnore
    public String getString() {
        return value.toString();
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
            valueString.getBytes(StandardCharsets.UTF_8).length*8,
            valueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
