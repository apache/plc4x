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
import java.time.Duration;

public class PlcLTIME extends PlcIECValue<Duration> {

    public static PlcLTIME of(Object value) {
        if (value instanceof PlcLTIME plcLTIME) {
            return plcLTIME;
        }
        if (value instanceof Duration duration) {
            return new PlcLTIME(duration);
        }
        if (value instanceof Byte b) {
            return new PlcLTIME(b);
        }
        if (value instanceof Short i) {
            return new PlcLTIME(i);
        }
        if (value instanceof Integer i) {
            return new PlcLTIME(i);
        }
        if (value instanceof Long l) {
            return new PlcLTIME(l);
        }
        if (value instanceof Float v) {
            return new PlcLTIME(v);
        }
        if (value instanceof Double v) {
            return new PlcLTIME(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcLTIME(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcLTIME(bigDecimal);
        }
        return new PlcLTIME(Duration.parse(value.toString()));
        //throw new PlcRuntimeException("Invalid value type");
    }

    public static PlcLTIME ofNanoseconds(long nanoseconds) {
        return new PlcLTIME(Duration.ofNanos(nanoseconds));
    }

    public static PlcLTIME ofNanoseconds(BigInteger nanoseconds) {
        // TODO: Not 100% correct, we're loosing precision here
        return new PlcLTIME(Duration.ofNanos(nanoseconds.longValue()));
    }

    public PlcLTIME(Byte nanoseconds) {
        this.value = Duration.ofNanos(nanoseconds);
        this.isNullable = false;
    }

    public PlcLTIME(Short nanoseconds) {
        this.value = Duration.ofNanos(nanoseconds);
        this.isNullable = false;
    }

    public PlcLTIME(Integer nanoseconds) {
        this.value = Duration.ofNanos(nanoseconds);
        this.isNullable = false;
    }

    public PlcLTIME(Long nanoseconds) {
        this.value = Duration.ofNanos(nanoseconds);
        this.isNullable = false;
    }

    public PlcLTIME(Float nanoseconds) {
        this.value = Duration.ofNanos(nanoseconds.longValue());
        this.isNullable = false;
    }

    public PlcLTIME(Double nanoseconds) {
        this.value = Duration.ofNanos(nanoseconds.longValue());
        this.isNullable = false;
    }

    public PlcLTIME(BigInteger nanoseconds) {
        this.value = Duration.ofNanos(nanoseconds.longValue());
        this.isNullable = false;
    }

    public PlcLTIME(BigDecimal nanoseconds) {
        this.value = Duration.ofNanos(nanoseconds.longValue());
        this.isNullable = false;
    }

    public PlcLTIME(Duration value) {
        this.value = value;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.LTIME;
    }

    public long getNanoseconds() {
        return value.toNanos();
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
    public boolean isBigInteger() {
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
        return (int) (value.toNanos() / 1000000);
    }

    @Override
    public long getLong() {
        return value.toNanos();
    }

    @Override
    public BigInteger getBigInteger() {
        return BigInteger.valueOf(value.toNanos());
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
        String valueString = value.toString();
        writeBuffer.writeString(valueString.getBytes(StandardCharsets.UTF_8).length * 8,
            valueString, WithOption.WithName(getClass().getSimpleName()), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
