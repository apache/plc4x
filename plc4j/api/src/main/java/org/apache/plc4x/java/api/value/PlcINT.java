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

import org.apache.plc4x.java.api.value.*;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.LinkedList;
import java.util.List;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcINT extends PlcIECValue<Short> {

    static Short minValue = Short.MIN_VALUE;
    static Short maxValue = Short.MAX_VALUE;

    public PlcINT(Boolean value) {
        super();
        this.value = value ? new Short((short) 1) : new Short((short) 0);
        this.isNullable = false;
    }

    public PlcINT(Byte value) {
        super();
        this.value = (Short) value.shortValue();
        this.isNullable = false;
    }

    public PlcINT(Short value) {
        super();
        this.value = value;
        this.isNullable = false;
    }

    public PlcINT(Integer value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcINT(Long value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcINT(Float value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcINT(Double value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcINT(BigInteger value) {
        super();
        if ((value.compareTo(BigInteger.valueOf(minValue)) >= 0) && (value.compareTo(BigInteger.valueOf(maxValue)) <= 0)) {
            this.value = (Short) value.shortValue();
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcINT(BigDecimal value) {
        super();
        if ((value.compareTo(BigDecimal.valueOf(minValue)) >= 0) && (value.compareTo(BigDecimal.valueOf(maxValue)) <= 0) && (value.scale() <= 0)) {
            this.value = (Short) value.shortValue();
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcINT(String value) {
        super();
        try {
            Short val = Short.valueOf(value).shortValue();
            this.value = val;
            this.isNullable = false;
        }
        catch(Exception e) {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a PLCINT Value");
        }
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcINT(@JsonProperty("value") short value) {
        super();
        this.value = value;
        this.isNullable = false;
    }

    @Override
    @JsonIgnore
    public boolean isShort() {
        return true;
    }

    @Override
    @JsonIgnore
    public short getShort() {
        return value;
    }

    public short getINT() {
        return value;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return Short.toString(value);
    }

    public byte[] getBytes() {
        byte[] bytes = new byte[2];
        bytes[0] = (byte)((value >> 8) & 0xff);
        bytes[1] = (byte)(value & 0xff);
        return bytes;
    }

}
