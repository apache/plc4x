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

public class PlcUDINT extends PlcIECValue<Long> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %d - %d for a %s Value";
    static Long minValue = (long) 0;
    static Long maxValue = (long) Integer.MAX_VALUE * 2 + 1;

    public static PlcUDINT of(Object value) {
        if (value instanceof Boolean) {
            return new PlcUDINT((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcUDINT((Byte) value);
        } else if (value instanceof Short) {
            return new PlcUDINT((Short) value);
        } else if (value instanceof Integer) {
            return new PlcUDINT((Integer) value);
        } else if (value instanceof Long) {
            return new PlcUDINT((Long) value);
        } else if (value instanceof Float) {
            return new PlcUDINT((Float) value);
        } else if (value instanceof Double) {
            return new PlcUDINT((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcUDINT((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcUDINT((BigDecimal) value);
        } else {
            return new PlcUDINT((String) value);
        }
    }

    public PlcUDINT(Boolean value) {
        this.value = value ? (long) 1 : (long) 0;
        this.isNullable = false;
    }

    public PlcUDINT(Byte value) {
        if (value < minValue || value > maxValue) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcUDINT(Short value) {
        if (value < minValue || value > maxValue) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcUDINT(Integer value) {
        if (value < minValue || value > maxValue) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcUDINT(Long value) {
        if (value < minValue || value > maxValue) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value;
        this.isNullable = false;
    }

    public PlcUDINT(Float value) {
        if ((value < minValue) || (value > maxValue) || (value % 1 != 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcUDINT(Double value) {
        if ((value < minValue) || (value > maxValue) || (value % 1 != 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = false;
    }

    public PlcUDINT(BigInteger value) {
        if ((value.compareTo(BigInteger.valueOf(minValue)) < 0) || (value.compareTo(BigInteger.valueOf(maxValue)) > 0)) {
            throw new PlcInvalidTagException("Value of type " + value +
                " is out of range " + minValue + " - " + maxValue + " for a " +
                this.getClass().getSimpleName() + " Value");
        }
        this.value = value.longValue();
        this.isNullable = true;
    }

    public PlcUDINT(BigDecimal value) {
        if ((value.compareTo(BigDecimal.valueOf(minValue)) < 0) || (value.compareTo(BigDecimal.valueOf(maxValue)) > 0) || (value.scale() > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.longValue();
        this.isNullable = true;
    }

    public PlcUDINT(String value) {
        try {
            long val = Long.parseLong(value.trim());
            if (val >= minValue && val <= maxValue) {
                this.value = val;
                this.isNullable = false;
            } else {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcUDINT(long value) {
        if (value < minValue || value > maxValue) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.UDINT;
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
        byte[] bytes = new byte[4];
        bytes[0] = (byte) ((value >> 24) & 0xff);
        bytes[1] = (byte) ((value >> 16) & 0xff);
        bytes[2] = (byte) ((value >> 8) & 0xff);
        bytes[3] = (byte) (value & 0xff);
        return bytes;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.writeLong(getClass().getSimpleName(), 32, value);
    }

}
