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
import java.time.LocalTime;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcTIME_OF_DAY extends PlcSimpleValue<LocalTime> {

    public static PlcTIME_OF_DAY of(Object value) {
        if (value instanceof LocalTime) {
            return new PlcTIME_OF_DAY((LocalTime) value);
        } else if (value instanceof Long) {
            return new PlcTIME_OF_DAY(LocalTime.ofSecondOfDay(((long) value) / 1000));
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcTIME_OF_DAY(@JsonProperty("value") LocalTime value) {
        super(value, true);
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcTIME_OF_DAY(@JsonProperty("value") Long value) {
        super(LocalTime.ofNanoOfDay(value * 1000000), true);
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
