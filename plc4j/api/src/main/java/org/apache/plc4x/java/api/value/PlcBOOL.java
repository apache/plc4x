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

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcBOOL extends PlcIECValue<Boolean> {

    static int minValue = 0;
    static int maxValue = 1;

    public PlcBOOL(Boolean value) {
        super();
        this.value = value;
        this.isNullable = true;
    }

    public PlcBOOL(Byte value) {
        super();
        if ((value != null) && (value >= minValue && value <= maxValue)) {
            this.value = value >= 1;
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcBOOL(Short value) {
        super();
        if ((value != null) && (value >= minValue && value <= maxValue)) {
            this.value = value >= 1;
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcBOOL(Integer value) {
        super();
        if ((value != null) && (value >= minValue && value <= maxValue)) {
            this.value = value >= 1;
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcBOOL(Long value) {
        super();
        if ((value != null) && (value >= minValue && value <= maxValue)) {
            this.value = value >= 1;
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcBOOL(Float value) {
        super();
        if ((value != null) && (value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = value >= 1;
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcBOOL(Double value) {
        super();
        if ((value != null) && (value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = value >= 1;
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcBOOL(BigInteger value) {
        super();
        if ((value != null) && (value.compareTo(BigInteger.valueOf(minValue)) >= 0) && (value.compareTo(BigInteger.valueOf(maxValue)) <= 0)) {
            this.value = value.compareTo(BigInteger.valueOf(maxValue)) >= 0;
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcBOOL(BigDecimal value) {
        super();
        if ((value != null) && (value.compareTo(BigDecimal.valueOf(minValue)) >= 0) && (value.compareTo(BigDecimal.valueOf(maxValue)) <= 0) && (value.scale() <= 0)) {
            this.value = value.compareTo(BigDecimal.valueOf(maxValue)) >= 0;
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcBOOL(String value) {
        super();
        try {
            try {
                this.value = Integer.parseInt(value.trim()) > 0;
            } catch(NumberFormatException e) {
                //parseBoolean expects a string "true" or "false"
                this.value = Boolean.parseBoolean(value.trim());
            }
            this.isNullable = false;
        }
        catch(Exception e) {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcBOOL(@JsonProperty("value") boolean value) {
        super();
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
        return ((value != null) && value) ? new byte[] {0x01} : new byte[] {0x00};
    }

    @Override
    @JsonIgnore
    public String toString() {
        return Boolean.toString(value);
    }

}
