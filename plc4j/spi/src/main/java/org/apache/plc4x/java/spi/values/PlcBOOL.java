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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcBOOL extends PlcIECValue<Boolean> {

    private static final String VALUE_OUT_OF_RANGE = "Value of type %s is out of range %d - %d for a %s Value";
    static int minValue = 0;
    static int maxValue = 1;

    public static PlcBOOL of(Object value) {
        if (value instanceof Boolean) {
            return new PlcBOOL((Boolean) value);
        } else if (value instanceof Byte) {
            return new PlcBOOL((Byte) value);
        } else if (value instanceof Short) {
            return new PlcBOOL((Short) value);
        } else if (value instanceof Integer) {
            return new PlcBOOL((Integer) value);
        } else if (value instanceof Long) {
            return new PlcBOOL((Long) value);
        } else if (value instanceof Float) {
            return new PlcBOOL((Float) value);
        } else if (value instanceof Double) {
            return new PlcBOOL((Double) value);
        } else if (value instanceof BigInteger) {
            return new PlcBOOL((BigInteger) value);
        } else if (value instanceof BigDecimal) {
            return new PlcBOOL((BigDecimal) value);
        } else {
            return new PlcBOOL((String) value);
        }
    }

    public PlcBOOL(Boolean value) {
        this.value = value;
        this.isNullable = true;
    }

    public PlcBOOL(Byte value) {
        if ((value == null) || (value < minValue || value > maxValue)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value >= 1;
        this.isNullable = true;
    }

    public PlcBOOL(Short value) {
        if ((value == null) || (value < minValue || value > maxValue)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value >= 1;
        this.isNullable = true;
    }

    public PlcBOOL(Integer value) {
        if ((value == null) || (value < minValue || value > maxValue)) {
            throw new PlcInvalidFieldException(String.format("Value of type %d is out of range %d - %d for a %s Value", value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value >= 1;
        this.isNullable = true;
    }

    public PlcBOOL(Long value) {
        if ((value == null) || (value < minValue || value > maxValue)) {
            throw new PlcInvalidFieldException(String.format("Value of type %d is out of range %d - %d for a %s Value", value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value >= 1;
        this.isNullable = true;
    }

    public PlcBOOL(Float value) {
        if ((value == null) || (value < minValue) || (value > maxValue) || (value % 1 != 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value >= 1;
        this.isNullable = true;
    }

    public PlcBOOL(Double value) {
        if ((value == null) || (value < minValue) || (value > maxValue) || (value % 1 != 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value >= 1;
        this.isNullable = true;
    }

    public PlcBOOL(BigInteger value) {
        if ((value == null) || (value.compareTo(BigInteger.valueOf(minValue)) < 0) || (value.compareTo(BigInteger.valueOf(maxValue)) > 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.compareTo(BigInteger.valueOf(maxValue)) >= 0;
        this.isNullable = true;
    }

    public PlcBOOL(BigDecimal value) {
        if ((value == null) || (value.compareTo(BigDecimal.valueOf(minValue)) < 0) || (value.compareTo(BigDecimal.valueOf(maxValue)) > 0) || (value.scale() > 0)) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()));
        }
        this.value = value.compareTo(BigDecimal.valueOf(maxValue)) >= 0;
        this.isNullable = true;
    }

    public PlcBOOL(String value) {
        try {
            this.value = parseValue(value);
            this.isNullable = false;
        } catch (RuntimeException e) {
            throw new PlcInvalidFieldException(String.format(VALUE_OUT_OF_RANGE, value, minValue, maxValue, this.getClass().getSimpleName()), e);
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

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcBOOL(@JsonProperty("value") boolean value) {
        this.value = value;
        this.isNullable = true;
    }

    @Override
    @JsonIgnore
    public boolean isBoolean() {
        return true;
    }

    @Override
    @JsonIgnore
    public boolean getBoolean() {
        return (value != null) && value;
    }

    @Override
    @JsonIgnore
    public boolean isByte() {
        return true;
    }

    @Override
    @JsonIgnore
    public byte getByte() {
        return (byte) (((value != null) && value) ? 1 : 0);
    }

    @Override
    @JsonIgnore
    public boolean isShort() {
        return true;
    }

    @Override
    @JsonIgnore
    public short getShort() {
        return (short) (((value != null) && value) ? 1 : 0);
    }

    @Override
    @JsonIgnore
    public boolean isInteger() {
        return true;
    }

    @Override
    @JsonIgnore
    public int getInteger() {
        return ((value != null) && value) ? 1 : 0;
    }

    @Override
    @JsonIgnore
    public boolean isLong() {
        return true;
    }

    @Override
    @JsonIgnore
    public long getLong() {
        return ((value != null) && value) ? 1 : 0;
    }

    @Override
    @JsonIgnore
    public boolean isBigInteger() {
        return true;
    }

    @Override
    @JsonIgnore
    public BigInteger getBigInteger() {
        return value ? BigInteger.ONE : BigInteger.ZERO;
    }

    @Override
    @JsonIgnore
    public boolean isFloat() {
        return true;
    }

    @Override
    @JsonIgnore
    public float getFloat() {
        return ((value != null) && value) ? 1.0f : 0.0f;
    }

    @Override
    @JsonIgnore
    public boolean isDouble() {
        return true;
    }

    @Override
    @JsonIgnore
    public double getDouble() {
        return ((value != null) && value) ? 1.0 : 0.0;
    }

    @Override
    @JsonIgnore
    public boolean isBigDecimal() {
        return true;
    }

    @Override
    @JsonIgnore
    public BigDecimal getBigDecimal() {
        return value ? BigDecimal.ONE : BigDecimal.ZERO;
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

    @JsonIgnore
    public byte[] getBytes() {
        return ((value != null) && value) ? new byte[]{0x01} : new byte[]{0x00};
    }

    @Override
    @JsonIgnore
    public String toString() {
        return Boolean.toString(value);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.writeBit(getClass().getSimpleName(), value);
    }

}
