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

public class PlcREAL extends PlcIECValue<Float> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %f - %f for a %s Value";
    static final Float minValue = -Float.MAX_VALUE;
    static final Float maxValue = Float.MAX_VALUE;

    public static PlcREAL of(Object value) {
        if (value instanceof Boolean) {
            return new PlcREAL((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcREAL((Byte) value);
        } else if (value instanceof Short) {
            return new PlcREAL((Short) value);
        } else if (value instanceof Integer) {
            return new PlcREAL((Integer) value);
        } else if (value instanceof Long) {
            return new PlcREAL((Long) value);
        } else if (value instanceof Float) {
            return new PlcREAL((Float) value);
        } else if (value instanceof Double) {
            return new PlcREAL((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcREAL((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcREAL((BigDecimal) value);
        } else {
            return new PlcREAL((String) value);
        }
    }

    public PlcREAL(Boolean value) {
        this.value = value ? (Float) 1.0f : (Float) 0.0f;
        this.isNullable = false;
    }

    public PlcREAL(Byte value) {
        this.value = value.floatValue();
        this.isNullable = false;
    }

    public PlcREAL(Short value) {
        this.value = value.floatValue();
        this.isNullable = false;
    }

    public PlcREAL(Integer value) {
        this.value = value.floatValue();
        this.isNullable = false;
    }

    public PlcREAL(Float value) {
        this.value = value;
        this.isNullable = false;
    }

    public PlcREAL(Double value) {
        if ((value < minValue) || (value > maxValue)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.floatValue();
        this.isNullable = false;
    }

    public PlcREAL(BigInteger value) {
        BigDecimal val = new BigDecimal(value);
        if ((val.compareTo(BigDecimal.valueOf(minValue)) < 0) || (val.compareTo(BigDecimal.valueOf(maxValue)) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val.floatValue();
        this.isNullable = true;
    }

    public PlcREAL(BigDecimal value) {
        if ((value.compareTo(BigDecimal.valueOf(minValue)) < 0) || (value.compareTo(BigDecimal.valueOf(maxValue)) > 0) || (value.scale() > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.floatValue();
        this.isNullable = true;
    }

    public PlcREAL(String value) {
        try {
            this.value = Float.parseFloat(value.trim());
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
    }

    public PlcREAL(float value) {
        this.value = value;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.REAL;
    }

    @Override
    public boolean isBoolean() {
        return true;
    }

    @Override
    public boolean getBoolean() {
        return (value != null) && !value.equals(0.0f);
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
        return true;
    }

    @Override
    public float getFloat() {
        return value;
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
        return Float.toString(value);
    }

    public byte[] getBytes() {
        int intBits = Float.floatToIntBits(value);
        return new byte[]{(byte) (intBits >> 24), (byte) (intBits >> 16), (byte) (intBits >> 8), (byte) (intBits)};
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.writeDouble(getClass().getSimpleName(), 32, value);
    }

}
