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
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Duration;
import java.time.temporal.ChronoUnit;

public class PlcTIME extends PlcIECValue<Duration> {

    public static PlcTIME of(Object value) {
        if (value instanceof PlcTIME) {
            return (PlcTIME) value;
        } else if (value instanceof Duration) {
            return new PlcTIME((Duration) value);
        } else if (value instanceof Byte) {
            return new PlcTIME((Byte) value);
        } else if (value instanceof Short) {
            return new PlcTIME((Short) value);
        } else if (value instanceof Integer) {
            return new PlcTIME((Integer) value);
        } else if (value instanceof Long) {
            return new PlcTIME((Long) value);
        } else if (value instanceof Float) {
            return new PlcTIME((Float) value);
        } else if (value instanceof Double) {
            return new PlcTIME((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcTIME((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcTIME((BigDecimal) value);
        } else {
            return new PlcTIME(Duration.parse(value.toString()));
        }
        //throw new PlcRuntimeException("Invalid value type");
    }

    public static PlcTIME ofMilliseconds(long milliseconds) {
        return new PlcTIME(Duration.ofMillis(milliseconds));
    }

    public PlcTIME(Byte milliseconds) {
        this.value = Duration.ofMillis(milliseconds);
        this.isNullable = false;
    }

    public PlcTIME(Short milliseconds) {
        this.value = Duration.ofMillis(milliseconds);
        this.isNullable = false;
    }

    public PlcTIME(Integer milliseconds) {
        this.value = Duration.ofMillis(milliseconds);
        this.isNullable = false;
    }

    public PlcTIME(Long milliseconds) {
        this.value = Duration.ofMillis(milliseconds);
        this.isNullable = false;
    }

    public PlcTIME(Float milliseconds) {
        this.value = Duration.ofMillis(milliseconds.longValue());
        this.isNullable = false;
    }

    public PlcTIME(Double milliseconds) {
        this.value = Duration.ofMillis(milliseconds.longValue());
        this.isNullable = false;
    }

    public PlcTIME(BigInteger milliseconds) {
        this.value = Duration.ofMillis(milliseconds.longValue());
        this.isNullable = false;
    }

    public PlcTIME(BigDecimal milliseconds) {
        this.value = Duration.ofMillis(milliseconds.longValue());
        this.isNullable = false;
    }

    public PlcTIME(Duration value) {
        this.value = value;
        this.isNullable = false;
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
        return value.toMillis();
    }

    @Override
    public Duration getDuration() {
        return value;
    }

    @Override
    public String getString() {
        return value.toString();
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
//        String valueString = value.toString();
//        writeBuffer.writeString(getClass().getSimpleName(),
//            valueString.getBytes(StandardCharsets.UTF_8).length*8,
//            valueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        writeBuffer.writeLong(getClass().getSimpleName(), 32, value.toMillis());        
    }

}
