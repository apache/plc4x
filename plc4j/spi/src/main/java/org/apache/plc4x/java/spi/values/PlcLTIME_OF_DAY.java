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

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.LocalTime;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcLTIME_OF_DAY extends PlcSimpleValue<LocalTime> {

    public static PlcLTIME_OF_DAY of(Object value) {
        if (value instanceof LocalTime) {
            return new PlcLTIME_OF_DAY((LocalTime) value);
        } else if (value instanceof Long) {
            return new PlcLTIME_OF_DAY(LocalTime.ofSecondOfDay(((Long) value) / 1000));
        } else if (value instanceof BigInteger) {
            // TODO: Not 100% correct, we're loosing precision here
            return new PlcLTIME_OF_DAY(LocalTime.ofSecondOfDay(((BigInteger) value).longValue() / 1000));
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    public static PlcLTIME_OF_DAY ofNanosecondsSinceMidnight(BigInteger nanosecondsSinceMidnight) {
        // TODO: Not 100% correct, we're loosing precision here
        return new PlcLTIME_OF_DAY(LocalTime.ofNanoOfDay(nanosecondsSinceMidnight.longValue()));
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcLTIME_OF_DAY(@JsonProperty("value") LocalTime value) {
        super(value, true);
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcLTIME_OF_DAY(@JsonProperty("value") Long nanosecondsSinceMidnight) {
        super(LocalTime.ofNanoOfDay(nanosecondsSinceMidnight), true);
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcLTIME_OF_DAY(@JsonProperty("value") BigInteger nanosecondsSinceMidnight) {
        // TODO: Not 100% correct, we're loosing precision here
        super(LocalTime.ofNanoOfDay(nanosecondsSinceMidnight.longValue()), true);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.TIME_OF_DAY;
    }

    public BigInteger getNanosecondsSinceMidnight() {
        return BigInteger.valueOf(value.toSecondOfDay()).multiply(BigInteger.valueOf(1000_000_000)).add(BigInteger.valueOf(value.getNano()));
    }

    @Override
    public boolean isLong() {
        return true;
    }

    @Override
    public long getLong() {
        return ((long) value.toSecondOfDay()) * 1000;
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
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        String valueString = value.toString();
        writeBuffer.writeString(getClass().getSimpleName(), valueString.getBytes(StandardCharsets.UTF_8).length*8,StandardCharsets.UTF_8.name(),valueString);
    }

}
