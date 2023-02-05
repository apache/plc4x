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

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;

public class PlcLINT extends PlcIECValue<Long> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %d - %d for a %s Value";
    static Long minValue = Long.MIN_VALUE;
    static Long maxValue = Long.MAX_VALUE;

    public static PlcLINT of(Object value) {
        if (value instanceof Boolean) {
            return new PlcLINT((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcLINT((Byte) value);
        } else if (value instanceof Short) {
            return new PlcLINT((Short) value);
        } else if (value instanceof Integer) {
            return new PlcLINT((Integer) value);
        } else if (value instanceof Long) {
            return new PlcLINT((Long) value);
        } else if (value instanceof Float) {
            return new PlcLINT((Float) value);
        } else if (value instanceof Double) {
            return new PlcLINT((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcLINT((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcLINT((BigDecimal) value);
        } else {
            return new PlcLINT((String) value);
        }
    }

    public PlcLINT(Boolean value) {
        this.value = value ? (long) 1 : (long) 0;
        this.isNullable = false;
    }

    public PlcLINT(Byte value) {
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcLINT(Short value) {
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcLINT(Integer value) {
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcLINT(Long value) {
        this.value = value;
        this.isNullable = false;
    }

    public PlcLINT(Float value) {
        if ((value < minValue) || (value > maxValue) || (value % 1 != 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcLINT(Double value) {
        if ((value < minValue) || (value > maxValue) || (value % 1 != 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcLINT(BigInteger value) {
        if ((value.compareTo(BigInteger.valueOf(minValue)) < 0) || (value.compareTo(BigInteger.valueOf(maxValue)) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = true;
    }

    public PlcLINT(BigDecimal value) {
        if ((value.compareTo(BigDecimal.valueOf(minValue)) < 0) || (value.compareTo(BigDecimal.valueOf(maxValue)) > 0) || (value.scale() > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = true;
    }

    public PlcLINT(String value) {
        try {
            this.value = Long.parseLong(value.trim());
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcLINT(long value) {
        this.value = value;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.LINT;
    }

    @Override
    public boolean isBoolean() {
        return true;
    }

    @Override
    public boolean getBoolean() {
        return (value != null) && !value.equals(0L);
    }

    @Override
    public boolean isByte() {
        return (value != null) && (value <= Byte.MAX_VALUE) && (value >= Byte.MIN_VALUE);
    }

    @Override
    public byte getByte() {
        return value.byteValue();
    }

    @Override
    public boolean isShort() {
        return (value != null) && (value <= Short.MAX_VALUE) && (value >= Short.MIN_VALUE);
    }

    @Override
    public short getShort() {
        return value.shortValue();
    }

    @Override
    public boolean isInteger() {
        return (value != null) && (value <= Integer.MAX_VALUE) && (value >= Integer.MIN_VALUE);
    }

    @Override
    public int getInteger() {
        return value.intValue();
    }

    @Override
    public boolean isLong() {
        return true;
    }

    @Override
    public long getLong() {
        return value;
    }

    @Override
    public boolean isBigInteger() {
        return true;
    }

    @Override
    public BigInteger getBigInteger() {
        return BigInteger.valueOf(value);
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
        return BigDecimal.valueOf(getDouble());
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
        return Long.toString(value);
    }

    public byte[] getBytes() {
        return new byte[]{
            (byte) ((value >> 56) & 0xff),
            (byte) ((value >> 48) & 0xff),
            (byte) ((value >> 40) & 0xff),
            (byte) ((value >> 32) & 0xff),
            (byte) ((value >> 24) & 0xff),
            (byte) ((value >> 16) & 0xff),
            (byte) ((value >> 8) & 0xff),
            (byte) (value & 0xff)
        };
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.writeLong(getClass().getSimpleName(), 64, value);
    }

}
