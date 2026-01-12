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
import java.time.LocalTime;

public class PlcLTIME_OF_DAY extends PlcIECValue<LocalTime> {

    public static PlcLTIME_OF_DAY of(Object value) {
        if (value instanceof PlcLTIME_OF_DAY plcLTIMEOfDay) {
            return plcLTIMEOfDay;
        }
        if (value instanceof LocalTime localTime) {
            return new PlcLTIME_OF_DAY(localTime);
        }
        if (value instanceof Byte b) {
            return new PlcLTIME_OF_DAY(b);
        }
        if (value instanceof Short i) {
            return new PlcLTIME_OF_DAY(i);
        }
        if (value instanceof Integer i) {
            return new PlcLTIME_OF_DAY(i);
        }
        if (value instanceof Long l) {
            return new PlcLTIME_OF_DAY(l);
        }
        if (value instanceof Float v) {
            return new PlcLTIME_OF_DAY(v);
        }
        if (value instanceof Double v) {
            return new PlcLTIME_OF_DAY(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcLTIME_OF_DAY(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcLTIME_OF_DAY(bigDecimal);
        }
        return new PlcLTIME_OF_DAY(LocalTime.parse(value.toString()));
    }

    public static PlcLTIME_OF_DAY ofNanosecondsSinceMidnight(BigInteger nanosecondsSinceMidnight) {
        // TODO: Not 100% correct, we're loosing precision here
        return new PlcLTIME_OF_DAY(LocalTime.ofNanoOfDay(nanosecondsSinceMidnight.longValue()));
    }

    public PlcLTIME_OF_DAY(Byte millisecondsSinceMidnight) {
        this.value = LocalTime.ofNanoOfDay((long) millisecondsSinceMidnight * 1000000);
        this.isNullable = false;
    }

    public PlcLTIME_OF_DAY(Short millisecondsSinceMidnight) {
        this.value = LocalTime.ofNanoOfDay((long) millisecondsSinceMidnight * 1000000);
        this.isNullable = false;
    }

    public PlcLTIME_OF_DAY(Integer millisecondsSinceMidnight) {
        this.value = LocalTime.ofNanoOfDay((long) millisecondsSinceMidnight * 1000000);
        this.isNullable = false;
    }

    public PlcLTIME_OF_DAY(Long millisecondsSinceMidnight) {
        this.value = LocalTime.ofNanoOfDay(millisecondsSinceMidnight * 1000000);
        this.isNullable = false;
    }

    public PlcLTIME_OF_DAY(Float millisecondsSinceMidnight) {
        this.value = LocalTime.ofNanoOfDay(millisecondsSinceMidnight.longValue() * 1000000);
        this.isNullable = false;
    }

    public PlcLTIME_OF_DAY(Double millisecondsSinceMidnight) {
        this.value = LocalTime.ofNanoOfDay(millisecondsSinceMidnight.longValue() * 1000000);
        this.isNullable = false;
    }

    public PlcLTIME_OF_DAY(BigInteger millisecondsSinceMidnight) {
        this.value = LocalTime.ofNanoOfDay(millisecondsSinceMidnight.longValue() * 1000000);
        this.isNullable = false;
    }

    public PlcLTIME_OF_DAY(BigDecimal millisecondsSinceMidnight) {
        this.value = LocalTime.ofNanoOfDay(millisecondsSinceMidnight.longValue() * 1000000);
        this.isNullable = false;
    }

    public PlcLTIME_OF_DAY(LocalTime value) {
        this.value = value;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.LTIME_OF_DAY;
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
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        String valueString = value.toString();
        writeBuffer.writeString(valueString.getBytes(StandardCharsets.UTF_8).length * 8,
            valueString, WithOption.WithName(getClass().getSimpleName()), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
