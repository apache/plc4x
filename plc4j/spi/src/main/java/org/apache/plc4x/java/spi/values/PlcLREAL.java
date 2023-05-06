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

public class PlcLREAL extends PlcIECValue<Double> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %d - %d for a %s Value";
    static final Double minValue = -Double.MAX_VALUE;
    static final Double maxValue = Double.MAX_VALUE;

    public static PlcLREAL of(Object value) {
        if (value instanceof Boolean) {
            return new PlcLREAL((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcLREAL((Byte) value);
        } else if (value instanceof Short) {
            return new PlcLREAL((Short) value);
        } else if (value instanceof Integer) {
            return new PlcLREAL((Integer) value);
        } else if (value instanceof Long) {
            return new PlcLREAL((Long) value);
        } else if (value instanceof Float) {
            return new PlcLREAL((Float) value);
        } else if (value instanceof Double) {
            return new PlcLREAL((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcLREAL((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcLREAL((BigDecimal) value);
        } else {
            return new PlcLREAL((String) value);
        }
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
        if ((val.compareTo(BigDecimal.valueOf(minValue)) < 0) || (val.compareTo(BigDecimal.valueOf(maxValue)) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val.doubleValue();
        this.isNullable = true;
    }

    public PlcLREAL(BigDecimal value) {
        if ((value.compareTo(BigDecimal.valueOf(minValue)) < 0) || (value.compareTo(BigDecimal.valueOf(maxValue)) > 0) || (value.scale() > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.doubleValue();
        this.isNullable = true;
    }

    public PlcLREAL(String value) {
        try {
            this.value = Double.parseDouble(value.trim());
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
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
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.writeDouble(getClass().getSimpleName(), 64, value);
    }

}
