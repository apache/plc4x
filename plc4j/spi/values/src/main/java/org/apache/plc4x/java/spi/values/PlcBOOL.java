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

public class PlcBOOL extends PlcIECValue<Boolean> {

    private static final String VALUE_OUT_OF_RANGE = "Value of is required for a %s Value";

    public static PlcBOOL of(Object value) {
        if (value instanceof PlcBOOL plcBOOL) {
            return plcBOOL;
        }
        if (value instanceof Boolean b) {
            return new PlcBOOL(b);
        }
        if (value instanceof Byte b) {
            return new PlcBOOL(b);
        }
        if (value instanceof Short i) {
            return new PlcBOOL(i);
        }
        if (value instanceof Integer i) {
            return new PlcBOOL(i);
        }
        if (value instanceof Long l) {
            return new PlcBOOL(l);
        }
        if (value instanceof Float v) {
            return new PlcBOOL(v);
        }
        if (value instanceof Double v) {
            return new PlcBOOL(v);
        }
        if (value instanceof BigInteger bigInteger) {
            return new PlcBOOL(bigInteger);
        }
        if (value instanceof BigDecimal bigDecimal) {
            return new PlcBOOL(bigDecimal);
        }
        return new PlcBOOL(value.toString());
    }

    public PlcBOOL(Boolean value) {
        this.value = value;
        this.isNullable = true;
    }

    public PlcBOOL(Byte value) {
        if (value == null) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, this.getClass().getSimpleName()));
        }
        this.value = value != 0;
        this.isNullable = true;
    }

    public PlcBOOL(Short value) {
        if (value == null) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, this.getClass().getSimpleName()));
        }
        this.value = value != 0;
        this.isNullable = true;
    }

    public PlcBOOL(Integer value) {
        if (value == null) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, this.getClass().getSimpleName()));
        }
        this.value = value != 0;
        this.isNullable = true;
    }

    public PlcBOOL(Long value) {
        if (value == null) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, this.getClass().getSimpleName()));
        }
        this.value = value != 0;
        this.isNullable = true;
    }

    public PlcBOOL(Float value) {
        if (value == null) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, this.getClass().getSimpleName()));
        }
        this.value = value != 0.0;
        this.isNullable = true;
    }

    public PlcBOOL(Double value) {
        if (value == null) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, this.getClass().getSimpleName()));
        }
        this.value = value != 0.0d;
        this.isNullable = true;
    }

    public PlcBOOL(BigInteger value) {
        if (value == null) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, this.getClass().getSimpleName()));
        }
        this.value = value.compareTo(BigInteger.ZERO) != 0;
        this.isNullable = true;
    }

    public PlcBOOL(BigDecimal value) {
        if (value == null) {
            throw new PlcInvalidTagException(String.format(VALUE_OUT_OF_RANGE, this.getClass().getSimpleName()));
        }
        this.value = value.compareTo(BigDecimal.ZERO) != 0;
        this.isNullable = true;
    }

    public PlcBOOL(String value) {
        try {
            this.value = parseValue(value);
            this.isNullable = false;
        } catch (RuntimeException e) {
            throw new PlcInvalidTagException(String.format("Value %s could not be parsed to %s Value", value, this.getClass().getSimpleName()), e);
        }
    }

    private boolean parseValue(String value) {
        try {
            return Integer.parseInt(value.trim()) > 0;
        } catch (NumberFormatException ignore) {
            //parseBoolean expects a string "true" or "false"
            return Boolean.parseBoolean(value.trim());
        }
    }

    public PlcBOOL(boolean value) {
        this.value = value;
        this.isNullable = true;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.BOOL;
    }

    @Override
    public boolean isBoolean() {
        return true;
    }

    @Override
    public boolean getBoolean() {
        return (value != null) && value;
    }

    @Override
    public boolean isByte() {
        return true;
    }

    @Override
    public byte getByte() {
        return (byte) (((value != null) && value) ? 1 : 0);
    }

    @Override
    public boolean isShort() {
        return true;
    }

    @Override
    public short getShort() {
        return (short) (((value != null) && value) ? 1 : 0);
    }

    @Override
    public boolean isInteger() {
        return true;
    }

    @Override
    public int getInteger() {
        return ((value != null) && value) ? 1 : 0;
    }

    @Override
    public boolean isLong() {
        return true;
    }

    @Override
    public long getLong() {
        return ((value != null) && value) ? 1 : 0;
    }

    @Override
    public boolean isBigInteger() {
        return true;
    }

    @Override
    public BigInteger getBigInteger() {
        return value ? BigInteger.ONE : BigInteger.ZERO;
    }

    @Override
    public boolean isFloat() {
        return true;
    }

    @Override
    public float getFloat() {
        return ((value != null) && value) ? 1.0f : 0.0f;
    }

    @Override
    public boolean isDouble() {
        return true;
    }

    @Override
    public double getDouble() {
        return ((value != null) && value) ? 1.0 : 0.0;
    }

    @Override
    public boolean isBigDecimal() {
        return true;
    }

    @Override
    public BigDecimal getBigDecimal() {
        return value ? BigDecimal.ONE : BigDecimal.ZERO;
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
    public byte[] getRaw() {
        return getBytes();
    }

    public byte[] getBytes() {
        return ((value != null) && value) ? new byte[]{0x01} : new byte[]{0x00};
    }

    @Override
    public String toString() {
        return Boolean.toString(value);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.writeBit(value, WithOption.WithName(getClass().getSimpleName()));
    }

}
