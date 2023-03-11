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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.time.LocalTime;

public class PlcTIME_OF_DAY extends PlcSimpleValue<LocalTime> {

    public static PlcTIME_OF_DAY of(Object value) {
        if (value instanceof LocalTime) {
            return new PlcTIME_OF_DAY((LocalTime) value);
        } else if(value instanceof Long) {
            return new PlcTIME_OF_DAY((Long) value);
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    public static PlcTIME_OF_DAY ofMillisecondsSinceMidnight(long millisecondsSinceMidnight) {
        return new PlcTIME_OF_DAY(LocalTime.ofNanoOfDay(millisecondsSinceMidnight * 1000_000));
    }

    public PlcTIME_OF_DAY(LocalTime value) {
        super(value, true);
    }

    public PlcTIME_OF_DAY(long millisecondsSinceMidnight) {
        super(LocalTime.ofNanoOfDay(millisecondsSinceMidnight * 1000_000), true);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.TIME_OF_DAY;
    }

    public long getMillisecondsSinceMidnight() {
        return ((long) value.toSecondOfDay() * 1000) + (value.getNano() / 1000_000);
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
