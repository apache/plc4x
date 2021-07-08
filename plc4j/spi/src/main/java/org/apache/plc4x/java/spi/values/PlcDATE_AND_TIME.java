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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.w3c.dom.Element;

import java.nio.charset.StandardCharsets;
import java.time.*;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcDATE_AND_TIME extends PlcSimpleValue<LocalDateTime> {

    public static PlcDATE_AND_TIME of(Object value) {
        if (value instanceof LocalDateTime) {
            return new PlcDATE_AND_TIME((LocalDateTime) value);
        } else if (value instanceof Long) {
            return new PlcDATE_AND_TIME(LocalDateTime.ofInstant(
                Instant.ofEpochSecond((long) value), ZoneId.of("UTC")));
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcDATE_AND_TIME(@JsonProperty("value") LocalDateTime value) {
        super(value, true);
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcDATE_AND_TIME(@JsonProperty("value") Long value) {
        super(LocalDateTime.ofInstant(
            Instant.ofEpochSecond(value), ZoneId.of("UTC")), true);
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
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        String valueString = value.toString();
        writeBuffer.writeString(getClass().getSimpleName(), valueString.getBytes(StandardCharsets.UTF_8).length*8,StandardCharsets.UTF_8.name(),valueString);
    }

}
