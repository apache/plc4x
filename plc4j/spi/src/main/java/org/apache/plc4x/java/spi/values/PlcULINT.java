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

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcULINT extends PlcIECValue<BigInteger> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %d - %d for a %s Value";
    static BigInteger minValue = BigInteger.valueOf(0);
    static BigInteger maxValue = BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.valueOf(2)).add(BigInteger.valueOf(1));

    public static PlcULINT of(Object value) {
        if (value instanceof Boolean) {
            return new PlcULINT((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcULINT((Byte) value);
        } else if (value instanceof Short) {
            return new PlcULINT((Short) value);
        } else if (value instanceof Integer) {
            return new PlcULINT((Integer) value);
        } else if (value instanceof Long) {
            return new PlcULINT((Long) value);
        } else if (value instanceof Float) {
            return new PlcULINT((Float) value);
        } else if (value instanceof Double) {
            return new PlcULINT((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcULINT((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcULINT((BigDecimal) value);
        } else {
            return new PlcULINT((String) value);
        }
    }

    public PlcULINT(Boolean value) {
        this.value = value ? BigInteger.valueOf(1) : BigInteger.valueOf(0);
        this.isNullable = false;
    }

    public PlcULINT(Byte value) {
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcULINT(Short value) {
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcULINT(Integer value) {
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcULINT(Long value) {
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcULINT(Float value) {
        try {
            BigInteger val = BigDecimal.valueOf(value).toBigInteger();
            if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
    }

    public PlcULINT(Double value) {
        try {
            BigInteger val = BigDecimal.valueOf(value).toBigInteger();
            if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
    }

    public PlcULINT(BigInteger value) {
        if ((value.compareTo(minValue) < 0) || (value.compareTo(maxValue) > 0)) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value;
        this.isNullable = false;
    }

    public PlcULINT(BigDecimal value) {
        try {
            BigInteger val = value.toBigInteger();
            if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
    }

    public PlcULINT(String value) {
        try {
            BigInteger val = new BigInteger(value.trim());
            if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.ULINT;
    }

    @Override
    @JsonIgnore
    public boolean isBoolean() {
        return true;
    }

    @Override
    @JsonIgnore
    public boolean getBoolean() {
        return (value != null) && !value.equals(BigInteger.ZERO);
    }

    @Override
    @JsonIgnore
    public boolean isByte() {
        return (value != null) && (value.compareTo(BigInteger.valueOf(Byte.MAX_VALUE)) <= 0) && (value.compareTo(BigInteger.valueOf(Byte.MIN_VALUE)) >= 0);
    }

    @Override
    @JsonIgnore
    public byte getByte() {
        return value.byteValue();
    }

    @Override
    @JsonIgnore
    public boolean isShort() {
        return (value != null) && (value.compareTo(BigInteger.valueOf(Short.MAX_VALUE)) <= 0) && (value.compareTo(BigInteger.valueOf(Short.MIN_VALUE)) >= 0);
    }

    @Override
    @JsonIgnore
    public short getShort() {
        return value.shortValue();
    }

    @Override
    @JsonIgnore
    public boolean isInteger() {
        return (value != null) && (value.compareTo(BigInteger.valueOf(Integer.MAX_VALUE)) <= 0) && (value.compareTo(BigInteger.valueOf(Integer.MIN_VALUE)) >= 0);
    }

    @Override
    @JsonIgnore
    public int getInteger() {
        return value.intValue();
    }

    @Override
    @JsonIgnore
    public boolean isLong() {
        return (value != null) && (value.compareTo(BigInteger.valueOf(Long.MAX_VALUE)) <= 0) && (value.compareTo(BigInteger.valueOf(Long.MIN_VALUE)) >= 0);
    }

    @Override
    @JsonIgnore
    public long getLong() {
        return value.longValue();
    }

    @Override
    @JsonIgnore
    public boolean isBigInteger() {
        return true;
    }

    @Override
    @JsonIgnore
    public BigInteger getBigInteger() {
        return value;
    }

    @Override
    @JsonIgnore
    public boolean isFloat() {
        return true;
    }

    @Override
    @JsonIgnore
    public float getFloat() {
        return value.floatValue();
    }

    @Override
    @JsonIgnore
    public boolean isDouble() {
        return true;
    }

    @Override
    @JsonIgnore
    public double getDouble() {
        return value.doubleValue();
    }

    @Override
    @JsonIgnore
    public boolean isBigDecimal() {
        return true;
    }

    @Override
    @JsonIgnore
    public BigDecimal getBigDecimal() {
        return new BigDecimal(value);
    }

    @Override
    @JsonIgnore
    public boolean isString() {
        return true;
    }

    @Override
    @JsonIgnore
    public String getString() {
        return toString();
    }

    @Override
    @JsonIgnore
    public String toString() {
        return value.toString();
    }

    @JsonIgnore
    public byte[] getBytes() {
        byte[] tmp = value.toByteArray();
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
        writeBuffer.writeBigInteger(getClass().getSimpleName(), 64, value);
    }

}
