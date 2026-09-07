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

public class PlcULINT extends PlcIECValue<BigInteger> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %s - %s for a %s Value";
    public static final BigInteger MIN_VALUE = BigInteger.ZERO;
    public static final BigInteger MAX_VALUE = BigInteger.ONE.shiftLeft(64).subtract(BigInteger.ONE); // 2^64 - 1

    public static PlcULINT of(Object value) {
        if (value instanceof PlcULINT plcULINT) {
            return plcULINT;
        }
        if (value instanceof Boolean b) {
            return new PlcULINT(b);
        }
        if (value instanceof Byte b) {
            return new PlcULINT(b);
        }
        if (value instanceof Short i) {
            return new PlcULINT(i);
        }
        if (value instanceof Integer i) {
            return new PlcULINT(i);
        }
        if (value instanceof Long l) {
            return new PlcULINT(l);
        }
        if (value instanceof Float v) {
            return new PlcULINT(v);
        }
        if (value instanceof Double v) {
            return new PlcULINT(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcULINT(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcULINT(bigDecimal);
        }
        return new PlcULINT(value.toString());
    }

    public PlcULINT(Boolean value) {
        this.value = value ? BigInteger.ONE : BigInteger.ZERO;
        this.isNullable = false;
    }

    public PlcULINT(Byte value) {
        BigInteger val = BigInteger.valueOf(value);
        if (val.compareTo(MIN_VALUE) < 0) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcULINT(Short value) {
        BigInteger val = BigInteger.valueOf(value);
        if (val.compareTo(MIN_VALUE) < 0) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcULINT(Integer value) {
        BigInteger val = BigInteger.valueOf(value);
        if (val.compareTo(MIN_VALUE) < 0) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcULINT(Long value) {
        BigInteger val = BigInteger.valueOf(value);
        if (val.compareTo(MIN_VALUE) < 0) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcULINT(Float value) {
        try {
            BigInteger val = BigDecimal.valueOf(value).toBigInteger();
            if (val.compareTo(MIN_VALUE) < 0 || val.compareTo(MAX_VALUE) > 0) {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()), e);
        }
    }

    public PlcULINT(Double value) {
        try {
            BigInteger val = BigDecimal.valueOf(value).toBigInteger();
            if (val.compareTo(MIN_VALUE) < 0 || val.compareTo(MAX_VALUE) > 0) {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()), e);
        }
    }

    public PlcULINT(BigInteger value) {
        if (value.compareTo(MIN_VALUE) < 0 || value.compareTo(MAX_VALUE) > 0) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value;
        this.isNullable = false;
    }

    public PlcULINT(BigDecimal value) {
        try {
            BigInteger val = value.toBigInteger();
            if (val.compareTo(MIN_VALUE) < 0 || val.compareTo(MAX_VALUE) > 0) {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()), e);
        }
    }

    public PlcULINT(String value) {
        try {
            BigInteger val = new BigInteger(value.trim());
            if (val.compareTo(MIN_VALUE) < 0 || val.compareTo(MAX_VALUE) > 0) {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()), e);
        }
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
        return (value != null) && !value.equals(BigInteger.ZERO);
    }

    @Override
    public boolean isByte() {
        return (value != null) && (value.compareTo(BigInteger.valueOf(Byte.MAX_VALUE)) <= 0) && (value.compareTo(BigInteger.valueOf(Byte.MIN_VALUE)) >= 0);
    }

    @Override
    public byte getByte() {
        return value.byteValue();
    }

    @Override
    public boolean isShort() {
        return (value != null) && (value.compareTo(BigInteger.valueOf(Short.MAX_VALUE)) <= 0) && (value.compareTo(BigInteger.valueOf(Short.MIN_VALUE)) >= 0);
    }

    @Override
    public short getShort() {
        return value.shortValue();
    }

    @Override
    public boolean isInteger() {
        return (value != null) && (value.compareTo(BigInteger.valueOf(Integer.MAX_VALUE)) <= 0) && (value.compareTo(BigInteger.valueOf(Integer.MIN_VALUE)) >= 0);
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
        return value.longValue();
    }

    @Override
    public boolean isBigInteger() {
        return true;
    }

    @Override
    public BigInteger getBigInteger() {
        return value;
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
        return value.toString();
    }

    @Override
    public byte[] getRaw() {
        return getBytes();
    }

    public byte[] getBytes() {
        byte[] bytes = new byte[8];
        bytes[0] = (byte) (value.shiftRight(56).and(BigInteger.valueOf(0xff)).intValue());
        bytes[1] = (byte) (value.shiftRight(48).and(BigInteger.valueOf(0xff)).intValue());
        bytes[2] = (byte) (value.shiftRight(40).and(BigInteger.valueOf(0xff)).intValue());
        bytes[3] = (byte) (value.shiftRight(32).and(BigInteger.valueOf(0xff)).intValue());
        bytes[4] = (byte) (value.shiftRight(24).and(BigInteger.valueOf(0xff)).intValue());
        bytes[5] = (byte) (value.shiftRight(16).and(BigInteger.valueOf(0xff)).intValue());
        bytes[6] = (byte) (value.shiftRight(8).and(BigInteger.valueOf(0xff)).intValue());
        bytes[7] = (byte) (value.and(BigInteger.valueOf(0xff)).intValue());
        return bytes;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.writeUnsignedBigInteger(64, value, WithOption.WithName(getClass().getSimpleName()));
    }

}
