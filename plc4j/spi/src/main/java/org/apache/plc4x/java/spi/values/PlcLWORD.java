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

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.util.BitSet;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcLWORD extends PlcIECValue<BigInteger> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %d - %d for a %s Value";
    static BigInteger minValue = BigInteger.valueOf(0);
    static BigInteger maxValue = BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.valueOf(2)).add(BigInteger.valueOf(1));

    public static PlcLWORD of(Object value) {
        if (value instanceof Boolean) {
            return new PlcLWORD((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcLWORD((Byte) value);
        } else if (value instanceof Short) {
            return new PlcLWORD((Short) value);
        } else if (value instanceof Integer) {
            return new PlcLWORD((Integer) value);
        } else if (value instanceof Long) {
            return new PlcLWORD((Long) value);
        } else if (value instanceof Float) {
            return new PlcLWORD((Float) value);
        } else if (value instanceof Double) {
            return new PlcLWORD((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcLWORD((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcLWORD((BigDecimal) value);
        } else {
            return new PlcLWORD((String) value);
        }
    }

    public PlcLWORD(Boolean value) {
        this.value = value ? BigInteger.valueOf(1) : BigInteger.valueOf(0);
        this.isNullable = false;
    }

    public PlcLWORD(Byte value) {
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcLWORD(Short value) {
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcLWORD(Integer value) {
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcLWORD(Long value) {
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = val;
        this.isNullable = false;
    }

    public PlcLWORD(Float value) {
        try {
            BigInteger val = BigDecimal.valueOf(value).toBigInteger();
            if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
                throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
    }

    public PlcLWORD(Double value) {
        try {
            BigInteger val = BigDecimal.valueOf(value).toBigInteger();
            if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
                throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
    }

    public PlcLWORD(BigInteger value) {
        if ((value.compareTo(minValue) < 0) || (value.compareTo(maxValue) > 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value;
        this.isNullable = false;
    }

    public PlcLWORD(BigDecimal value) {
        try {
            BigInteger val = value.toBigInteger();
            if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
                throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
    }

    public PlcLWORD(String value) {
        try {
            BigInteger val = new BigInteger(value.trim());
            if ((val.compareTo(minValue) < 0) || (val.compareTo(maxValue) > 0)) {
                throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
            this.value = val;
            this.isNullable = false;
        } catch (Exception e) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
        }
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

    @JsonIgnore
    public boolean[] getBooleanArray() {
        boolean[] booleanValues = new boolean[64];
        BitSet bitSet = BitSet.valueOf(this.value.toByteArray());
        for (int i = 0; i < 64; i++) {
            booleanValues[i] = bitSet.get(i);
        }
        return booleanValues;
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
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.writeBigInteger(getClass().getSimpleName(), 64, value);
    }

}
