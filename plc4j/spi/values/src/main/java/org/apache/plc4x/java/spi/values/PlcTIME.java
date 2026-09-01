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
import java.time.Duration;

public class PlcTIME extends PlcIECValue<Duration> {

    public static PlcTIME of(Object value) {
        if (value instanceof PlcTIME plcTIME) {
            return plcTIME;
        }
        if (value instanceof Duration duration) {
            return new PlcTIME(duration);
        }
        if (value instanceof Byte b) {
            return new PlcTIME(b);
        }
        if (value instanceof Short i) {
            return new PlcTIME(i);
        }
        if (value instanceof Integer i) {
            return new PlcTIME(i);
        }
        if (value instanceof Long l) {
            return new PlcTIME(l);
        }
        if (value instanceof Float v) {
            return new PlcTIME(v);
        }
        if (value instanceof Double v) {
            return new PlcTIME(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcTIME(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcTIME(bigDecimal);
        }
        return new PlcTIME(Duration.parse(value.toString()));
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
        return (int) value.toMillis();
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
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        // Durations render as ISO-8601 (like PlcLTIME), not as a raw millisecond count.
        String valueString = value.toString();
        writeBuffer.writeString(valueString.getBytes(java.nio.charset.StandardCharsets.UTF_8).length * 8,
            valueString, WithOption.WithName(getClass().getSimpleName()), WithOption.WithEncoding("UTF8"));
    }

}
