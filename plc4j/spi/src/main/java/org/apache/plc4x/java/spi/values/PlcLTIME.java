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

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.temporal.ChronoUnit;

public class PlcLTIME extends PlcSimpleValue<Duration> {

    public static PlcLTIME of(Object value) {
        if (value instanceof Duration) {
            return new PlcLTIME((Duration) value);
        } else if(value instanceof Integer) {
            return new PlcLTIME(Duration.of(((Integer) value).longValue(), ChronoUnit.MILLIS));
        } else if(value instanceof Long) {
            return new PlcLTIME(Duration.of((long) value, ChronoUnit.NANOS));
        } else if(value instanceof BigInteger) {
            // TODO: Not 100% correct, we're loosing precision here
            return new PlcLTIME(Duration.of(((BigInteger) value).longValue(), ChronoUnit.NANOS));
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    public static PlcLTIME ofNanoseconds(long nanoseconds) {
        return new PlcLTIME(Duration.ofNanos(nanoseconds));
    }

    public static PlcLTIME ofNanoseconds(BigInteger nanoseconds) {
        // TODO: Not 100% correct, we're loosing precision here
        return new PlcLTIME(Duration.ofNanos(nanoseconds.longValue()));
    }

    public PlcLTIME(Duration value) {
        super(value, true);
    }

    public PlcLTIME(long nanoseconds) {
        super(Duration.ofNanos(nanoseconds), true);
    }

    public PlcLTIME(BigInteger nanoseconds) {
        // TODO: Not 100% correct, we're loosing precision here
        super(Duration.ofNanos(nanoseconds.longValue()), true);
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
        return (int) (value.get(ChronoUnit.NANOS) / 1000000);
    }

    @Override
    public long getLong() {
        return value.get(ChronoUnit.NANOS);
    }

    @Override
    public BigInteger getBigInteger() {
        return BigInteger.valueOf(value.get(ChronoUnit.NANOS));
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
        String valueString = value.toString();
        writeBuffer.writeString(getClass().getSimpleName(),
            valueString.getBytes(StandardCharsets.UTF_8).length*8,
            valueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
