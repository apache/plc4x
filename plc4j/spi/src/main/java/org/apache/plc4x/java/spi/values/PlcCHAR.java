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
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;

public class PlcCHAR extends PlcIECValue<Short> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %d - %d for a %s Value";
    static Short minValue = 0;
    static Short maxValue = (short) Byte.MAX_VALUE * 2 + 1;

    public static PlcCHAR of(Object value) {
        if (value instanceof Boolean) {
            return new PlcCHAR((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcCHAR((Byte) value);
        } else if (value instanceof Short) {
            return new PlcCHAR((Short) value);
        } else if (value instanceof Integer) {
            return new PlcCHAR((Integer) value);
        } else if (value instanceof Long) {
            return new PlcCHAR((Long) value);
        } else if (value instanceof Float) {
            return new PlcCHAR((Float) value);
        } else if (value instanceof Double) {
            return new PlcCHAR((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcCHAR((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcCHAR((BigDecimal) value);
        } else {
            return new PlcCHAR((String) value);
        }
    }

    public PlcCHAR(Boolean value) {
        super();
        this.value = value ? Short.valueOf((short) 1) : Short.valueOf((short) 0);
        this.isNullable = false;
    }

    public PlcCHAR(Character value) {
        super();
        Integer val = (int) value;
        if ((val >= minValue) && (val <= maxValue)) {
            this.value = val.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(Byte value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(Short value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value;
            this.isNullable = false;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(Integer value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(Long value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(Float value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(Double value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(BigInteger value) {
        super();
        if ((value.compareTo(BigInteger.valueOf(minValue)) >= 0) && (value.compareTo(BigInteger.valueOf(maxValue)) <= 0)) {
            this.value = value.shortValue();
            this.isNullable = true;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(BigDecimal value) {
        super();
        if ((value.compareTo(BigDecimal.valueOf(minValue)) >= 0) && (value.compareTo(BigDecimal.valueOf(maxValue)) <= 0) && (value.scale() <= 0)) {
            this.value = value.shortValue();
            this.isNullable = true;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(String value) {
        super();
        try {
            //If there is a extra space around the character trim it, unless you are actually sending a space
            String s = value.trim();
            if (s.length() == 0) {
                s = " ";
            }
            short val = (short) s.charAt(0);
            if ((val >= minValue) && (val <= maxValue)) {
                this.value = val;
                this.isNullable = false;
            } else {
                throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
            }
        } catch (Exception e) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    public PlcCHAR(short value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value;
            this.isNullable = false;
        } else {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.CHAR;
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
        return (value != null) && (value <= Byte.MAX_VALUE) && (value >= Byte.MIN_VALUE);
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
    public Object getObject() {
        return Character.toString((char) ((short) value));
    }

    @Override
    public String toString() {
        return Character.toString((char) ((short) value));
    }

    public byte[] getBytes() {
        byte[] bytes = new byte[1];
        bytes[0] = (byte) (value & 0xff);
        return bytes;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        String valueString = value.toString();
        writeBuffer.writeString(getClass().getSimpleName(),
            16,
            valueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
