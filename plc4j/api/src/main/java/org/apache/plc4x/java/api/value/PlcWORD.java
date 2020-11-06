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

package org.apache.plc4x.java.api.value;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.BitSet;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcWORD extends PlcIECValue<Integer> {

    static Integer minValue = 0;
    static Integer maxValue = Short.MAX_VALUE * 2 + 1;

    public PlcWORD(Boolean value) {
        super();
        this.value = value ? (Integer) 1 : (Integer) 0;
        this.isNullable = false;
    }

    public PlcWORD(Byte value) {
        super();
        this.value = value.intValue();
        this.isNullable = false;
    }

    public PlcWORD(Short value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value.intValue();
            this.isNullable = false;
        } else {
            throw new IllegalArgumentException("Value of type " + value + " is out of range " + minValue + " - " + maxValue + " for a PLCUINT Value");
        }
    }

    public PlcWORD(Integer value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value;
            this.isNullable = false;
        } else {
            throw new IllegalArgumentException("Value of type " + value + " is out of range " + minValue + " - " + maxValue + " for a PLCUINT Value");
        }
    }

    public PlcWORD(Long value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value.intValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcWORD(Float value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = value.intValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcWORD(Double value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = value.intValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcWORD(BigInteger value) {
        super();
        if ((value.compareTo(BigInteger.valueOf(minValue)) >= 0) && (value.compareTo(BigInteger.valueOf(maxValue)) <= 0)) {
            this.value = value.intValue();
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcWORD(BigDecimal value) {
        super();
        if ((value.compareTo(BigDecimal.valueOf(minValue)) >= 0) && (value.compareTo(BigDecimal.valueOf(maxValue)) <= 0) && (value.scale() <= 0)) {
            this.value = value.intValue();
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcWORD(String value) {
        super();
        try {
            int val = Integer.parseInt(value.trim());
            if ((val >= minValue) && (val <= maxValue)) {
                this.value = val;
                this.isNullable = false;
            } else {
                throw new IllegalArgumentException("Value of type " + value + " is out of range " + minValue + " - " + maxValue + " for a PLCUINT Value");
            }
        }
        catch(Exception e) {
            throw new IllegalArgumentException("Value of type " + value + " is out of range " + minValue + " - " + maxValue + " for a PLCUINT Value");
        }
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcWORD(@JsonProperty("value") int value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value;
            this.isNullable = false;
        } else {
            throw new IllegalArgumentException("Value of type " + value + " is out of range " + minValue + " - " + maxValue + " for a PLCUINT Value");
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
        return (value != null) && !value.equals(0);
    }

    @JsonIgnore
    public boolean[] getBooleanArray() {
        boolean[] booleanValues = new boolean[16];
        BitSet bitSet = BitSet.valueOf(new long[]{this.value});
        for (int i = 0; i < 16; i++) {
            booleanValues[i] = bitSet.get(i);
        }
        return booleanValues;
    }

    @Override
    @JsonIgnore
    public boolean isByte() {
        return (value != null) && (value <= Byte.MAX_VALUE) && (value >= Byte.MIN_VALUE);
    }

    @Override
    @JsonIgnore
    public byte getByte() {
        return value.byteValue();
    }

    @Override
    @JsonIgnore
    public boolean isShort() {
        return (value != null) && (value <= Short.MAX_VALUE) && (value >= Short.MIN_VALUE);
    }

    @Override
    @JsonIgnore
    public short getShort() {
        return value.shortValue();
    }

    @Override
    @JsonIgnore
    public boolean isInteger() {
        return true;
    }

    @Override
    @JsonIgnore
    public int getInteger() {
        return value;
    }

    @Override
    @JsonIgnore
    public boolean isLong() {
        return true;
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
        return BigInteger.valueOf(getLong());
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
        return BigDecimal.valueOf(getFloat());
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
        return Integer.toString(value);
    }

    @JsonIgnore
    public byte[] getBytes() {
        byte[] bytes = new byte[2];
        bytes[0] = (byte)((value >> 8) & 0xff);
        bytes[1] = (byte)(value & 0xff);
        return bytes;
    }

}
