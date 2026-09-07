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
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.math.BigDecimal;
import java.math.BigInteger;

public class PlcLREAL extends PlcIECValue<Double> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %f - %f for a %s Value";
    public static final Double MIN_VALUE = -Double.MAX_VALUE;
    public static final Double MAX_VALUE = Double.MAX_VALUE;

    public static PlcLREAL of(Object value) {
        if (value instanceof PlcLREAL plcLREAL) {
            return plcLREAL;
        }
        if (value instanceof Boolean b) {
            return new PlcLREAL(b);
        }
        if (value instanceof Byte b) {
            return new PlcLREAL(b);
        }
        if (value instanceof Short i) {
            return new PlcLREAL(i);
        }
        if (value instanceof Integer i) {
            return new PlcLREAL(i);
        }
        if (value instanceof Long l) {
            return new PlcLREAL(l);
        }
        if (value instanceof Float v) {
            return new PlcLREAL(v);
        }
        if (value instanceof Double v) {
            return new PlcLREAL(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcLREAL(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcLREAL(bigDecimal);
        }
        return new PlcLREAL(value.toString());
    }

    public PlcLREAL(Boolean value) {
        this.value = value ? (Double) 1.0 : (Double) 0.0;
        this.isNullable = false;
    }

    public PlcLREAL(Byte value) {
        this.value = value.doubleValue();
        this.isNullable = false;
    }

    public PlcLREAL(Short value) {
        this.value = value.doubleValue();
        this.isNullable = false;
    }

    public PlcLREAL(Integer value) {
        this.value = value.doubleValue();
        this.isNullable = false;
    }

    public PlcLREAL(Long value) {
        this.value = value.doubleValue();
        this.isNullable = false;
    }

    public PlcLREAL(Float value) {
        this.value = value.doubleValue();
        this.isNullable = false;
    }

    public PlcLREAL(Double value) {
        this.value = value;
        this.isNullable = false;
    }

    public PlcLREAL(BigInteger value) {
        BigDecimal val = new BigDecimal(value);
        if ((val.compareTo(BigDecimal.valueOf(MIN_VALUE)) < 0) || (val.compareTo(BigDecimal.valueOf(MAX_VALUE)) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = val.doubleValue();
        this.isNullable = true;
    }

    public PlcLREAL(BigDecimal value) {
        if ((value.compareTo(BigDecimal.valueOf(MIN_VALUE)) < 0) || (value.compareTo(BigDecimal.valueOf(MAX_VALUE)) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value.doubleValue();
        this.isNullable = true;
    }

    public PlcLREAL(String value) {
        try {
            this.value = Double.parseDouble(value.trim());
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
    }

    public PlcLREAL(double value) {
        this.value = value;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.LREAL;
    }

    @Override
    public boolean isBoolean() {
        return true;
    }

    @Override
    public boolean getBoolean() {
        return (value != null) && !value.equals((double) 0);
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
        return (value != null) && (value <= Long.MAX_VALUE) && (value >= Long.MIN_VALUE);
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
        return BigInteger.valueOf(getLong());
    }

    @Override
    public boolean isFloat() {
        return (value != null) && (value <= Float.MAX_VALUE) && (value >= -Float.MAX_VALUE);
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
        return value;
    }

    @Override
    public boolean isBigDecimal() {
        return true;
    }

    @Override
    public BigDecimal getBigDecimal() {
        return new BigDecimal(value);
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
        return Double.toString(value);
    }

    @Override
    public byte[] getRaw() {
        return getBytes();
    }

    public byte[] getBytes() {
        long longBits = Double.doubleToRawLongBits(value);
        return new byte[]{
            (byte) ((longBits >> 56) & 0xff),
            (byte) ((longBits >> 48) & 0xff),
            (byte) ((longBits >> 40) & 0xff),
            (byte) ((longBits >> 32) & 0xff),
            (byte) ((longBits >> 24) & 0xff),
            (byte) ((longBits >> 16) & 0xff),
            (byte) ((longBits >> 8) & 0xff),
            (byte) (longBits & 0xff)
        };
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.writeDouble(64, value, WithOption.WithName(getClass().getSimpleName()));
    }

}
