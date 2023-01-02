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
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.time.*;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcDATE_AND_TIME extends PlcSimpleValue<LocalDateTime> {

    public static PlcDATE_AND_TIME of(Object value) {
        if (value instanceof LocalDateTime) {
            return new PlcDATE_AND_TIME((LocalDateTime) value);
        } else if (value instanceof Long) {
            return new PlcDATE_AND_TIME(LocalDateTime.ofInstant(
                Instant.ofEpochSecond((long) value), ZoneId.systemDefault()));
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    public static PlcDATE_AND_TIME ofSecondsSinceEpoch(long secondsSinceEpoch) {
        return new PlcDATE_AND_TIME(LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0,
            ZoneOffset.UTC));
    }

    public static PlcDATE_AND_TIME ofSegments(int year, int month, int day, int hour, int minutes, int seconds, int nanoseconds) {
        return new PlcDATE_AND_TIME(LocalDateTime.of(year, month, day, hour, minutes, seconds, nanoseconds));
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcDATE_AND_TIME(@JsonProperty("value") LocalDateTime value) {
        super(value, true);
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcDATE_AND_TIME(@JsonProperty("value") long secondsSinceEpoch) {
        super(LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0,
            ZoneOffset.UTC), true);
    }

    public PlcDATE_AND_TIME(@JsonProperty("year") int year, @JsonProperty("month") int month, @JsonProperty("day") int day, @JsonProperty("hour") int hour, @JsonProperty("minutes") int minutes, @JsonProperty("seconds") int seconds, @JsonProperty("nanoseconds") int nanoseconds) {
        super(LocalDateTime.of(year, month, day, hour, minutes, seconds, nanoseconds), true);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.DATE_AND_TIME;
    }

    public long getSecondsSinceEpoch() {
        Instant instant = getDateTime().toInstant(ZoneOffset.of(ZoneOffset.UTC.getId()));
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
        Instant instant = value.atZone(ZoneId.systemDefault()).toInstant();
        return instant.getEpochSecond();
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
    public boolean isTime() {
        return true;
    }

    @Override
    @JsonIgnore
    public LocalTime getTime() {
        return value.toLocalTime();
    }

    @Override
    @JsonIgnore
    public boolean isDate() {
        return true;
    }

    @Override
    @JsonIgnore
    public LocalDate getDate() {
        return value.toLocalDate();
    }

    @Override
    @JsonIgnore
    public boolean isDateTime() {
        return true;
    }

    @Override
    @JsonIgnore
    public LocalDateTime getDateTime() {
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
        writeBuffer.writeString(getClass().getSimpleName(), valueString.getBytes(StandardCharsets.UTF_8).length*8,StandardCharsets.UTF_8.name(),valueString);
    }

}
