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

public class PlcUSINT extends PlcIECValue<Short> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %d - %d for a %s Value";
    public static final Short MIN_VALUE = 0;
    public static final Short MAX_VALUE = 255; // Unsigned 8-bit: 0 to 2^8-1

    public static PlcUSINT of(Object value) {
        if (value instanceof PlcUSINT plcUSINT) {
            return plcUSINT;
        }
        if (value instanceof Boolean b) {
            return new PlcUSINT(b);
        }
        if (value instanceof Byte b) {
            return new PlcUSINT(b);
        }
        if (value instanceof Short i) {
            return new PlcUSINT(i);
        }
        if (value instanceof Integer i) {
            return new PlcUSINT(i);
        }
        if (value instanceof Long l) {
            return new PlcUSINT(l);
        }
        if (value instanceof Float v) {
            return new PlcUSINT(v);
        }
        if (value instanceof Double v) {
            return new PlcUSINT(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcUSINT(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcUSINT(bigDecimal);
        }
        return new PlcUSINT(value.toString());
    }

    public PlcUSINT(Boolean value) {
        this.value = value ? (short) 1 : (short) 0;
        this.isNullable = false;
    }

    public PlcUSINT(Byte value) {
        if ((value < MIN_VALUE) || (value > MAX_VALUE)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value.shortValue();
        this.isNullable = false;
    }

    public PlcUSINT(Short value) {
        if ((value < MIN_VALUE) || (value > MAX_VALUE)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value;
        this.isNullable = false;
    }

    public PlcUSINT(Integer value) {
        if ((value < MIN_VALUE) || (value > MAX_VALUE)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value.shortValue();
        this.isNullable = false;
    }

    public PlcUSINT(Long value) {
        if ((value < MIN_VALUE) || (value > MAX_VALUE)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value.shortValue();
        this.isNullable = false;
    }

    public PlcUSINT(Float value) {
        if ((value < MIN_VALUE) || (value > MAX_VALUE)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value.shortValue();
        this.isNullable = false;
    }

    public PlcUSINT(Double value) {
        if ((value < MIN_VALUE) || (value > MAX_VALUE)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value.shortValue();
        this.isNullable = false;
    }

    public PlcUSINT(BigInteger value) {
        if ((value.compareTo(BigInteger.valueOf(MIN_VALUE)) < 0) || (value.compareTo(BigInteger.valueOf(MAX_VALUE)) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value.shortValue();
        this.isNullable = false;
    }

    public PlcUSINT(BigDecimal value) {
        if ((value.compareTo(BigDecimal.valueOf(MIN_VALUE)) < 0) || (value.compareTo(BigDecimal.valueOf(MAX_VALUE)) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, MIN_VALUE, MAX_VALUE, this.getClass().getSimpleName()));
        }
        this.value = value.shortValue();
        this.isNullable = false;
    }

    public PlcUSINT(String value) {
        try {
            short val = Short.parseShort(value.trim());
            if ((val < MIN_VALUE) || (val > MAX_VALUE)) {
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
        return PlcValueType.USINT;
    }

    @Override
    public boolean isBoolean() {
        return true;
    }

    @Override
    public boolean getBoolean() {
        return (value != null) && !value.equals((short) 0);
    }

    @Override
    public boolean isByte() {
        return value != null;
    }

    @Override
    public byte getByte() {
        return value.byteValue();
    }

    @Override
    public boolean isShort() {
        return true;
    }

    @Override
    public short getShort() {
        return value;
    }

    @Override
    public boolean isInteger() {
        return true;
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
        return BigInteger.valueOf(getLong());
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
        return BigDecimal.valueOf(getFloat());
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
        return Short.toString(value);
    }

    @Override
    public byte[] getRaw() {
        return getBytes();
    }

    public byte[] getBytes() {
        byte[] bytes = new byte[1];
        bytes[0] = (byte) (value & 0xff);
        return bytes;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.writeUnsignedShort(8, value, WithOption.WithName(getClass().getSimpleName()));
    }

}
