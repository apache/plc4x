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

import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;

public class PlcBREAL extends PlcIECValue<BigDecimal> {

    public static PlcBREAL of(Object value) {
        if (value instanceof Boolean) {
            return new PlcBREAL((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcBREAL((Byte) value);
        } else if (value instanceof Short) {
            return new PlcBREAL((Short) value);
        } else if (value instanceof Integer) {
            return new PlcBREAL((Integer) value);
        } else if (value instanceof Long) {
            return new PlcBREAL((Long) value);
        } else if (value instanceof Float) {
            return new PlcBREAL((Float) value);
        } else if (value instanceof Double) {
            return new PlcBREAL((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcBREAL((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcBREAL((BigDecimal) value);
        } else {
            return new PlcBREAL((String) value);
        }
    }

    public PlcBREAL(Boolean value) {
        this.value = value ? BigDecimal.valueOf(1) : BigDecimal.valueOf(0);
        this.isNullable = false;
    }

    public PlcBREAL(Byte value) {
        this.value = BigDecimal.valueOf(value);;
        this.isNullable = false;
    }

    public PlcBREAL(Short value) {
        this.value = BigDecimal.valueOf(value);;
        this.isNullable = false;
    }

    public PlcBREAL(Integer value) {
        this.value = BigDecimal.valueOf(value);;
        this.isNullable = false;
    }

    public PlcBREAL(Long value) {
        this.value = BigDecimal.valueOf(value);;
        this.isNullable = false;
    }

    public PlcBREAL(Float value) {
        this.value = BigDecimal.valueOf(value);
        this.isNullable = false;
    }

    public PlcBREAL(Double value) {
        this.value = BigDecimal.valueOf(value);
        this.isNullable = false;
    }

    public PlcBREAL(BigInteger value) {
        this.value = new BigDecimal(value);
        this.isNullable = false;
    }

    public PlcBREAL(BigDecimal value) {
        this.value = value;
        this.isNullable = false;
    }

    public PlcBREAL(String value) {
        this.value = new BigDecimal(value.trim());;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.ULINT;
    }

    @Override
    public boolean isBoolean() {
        return true;
    }

    @Override
    public boolean getBoolean() {
        return (value != null) && !value.equals(BigDecimal.ZERO);
    }

    @Override
    public boolean isByte() {
        return (value != null) && (value.compareTo(BigDecimal.valueOf(Byte.MAX_VALUE)) <= 0) && (value.compareTo(BigDecimal.valueOf(Byte.MIN_VALUE)) >= 0);
    }

    @Override
    public byte getByte() {
        return value.byteValue();
    }

    @Override
    public boolean isShort() {
        return (value != null) && (value.compareTo(BigDecimal.valueOf(Short.MAX_VALUE)) <= 0) && (value.compareTo(BigDecimal.valueOf(Short.MIN_VALUE)) >= 0);
    }

    @Override
    public short getShort() {
        return value.shortValue();
    }

    @Override
    public boolean isInteger() {
        return (value != null) && (value.compareTo(BigDecimal.valueOf(Integer.MAX_VALUE)) <= 0) && (value.compareTo(BigDecimal.valueOf(Integer.MIN_VALUE)) >= 0);
    }

    @Override
    public int getInteger() {
        return value.intValue();
    }

    @Override
    public boolean isLong() {
        return (value != null) && (value.compareTo(BigDecimal.valueOf(Long.MAX_VALUE)) <= 0) && (value.compareTo(BigDecimal.valueOf(Long.MIN_VALUE)) >= 0);
    }

    @Override
    public long getLong() {
        return value.longValue();
    }

    @Override
    public boolean isBigInteger() {
        return true;
    }

    @Override
    public BigInteger getBigInteger() {
        return value.toBigInteger();
    }

    @Override
    public boolean isFloat() {
        return true;
    }

    @Override
    public float getFloat() {
        return value.floatValue();
    }

    @Override
    public boolean isDouble() {
        return true;
    }

    @Override
    public double getDouble() {
        return value.doubleValue();
    }

    @Override
    public boolean isBigDecimal() {
        return true;
    }

    @Override
    public BigDecimal getBigDecimal() {
        return value;
    }

    @Override
    public boolean isString() {
        return true;
    }

    @Override
    public String getString() {
        return toString();
    }

    @Override
    public String toString() {
        return value.toString();
    }

    public byte[] getBytes() {
        // TODO: Not sure if this is correct ...
        byte[] tmp = value.unscaledValue().toByteArray();
        byte[] bytes = new byte[8];
        for (int i = 0; i < bytes.length; i++) {
            if (i >= (bytes.length - tmp.length)) {
                bytes[i] = tmp[i - (bytes.length - tmp.length)];
            } else {
                bytes[i] = 0x00;
            }
        }
        return bytes;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.writeBigDecimal(getClass().getSimpleName(), 64, value);
    }

}
