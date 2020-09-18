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
public class PlcCHAR extends PlcIECValue<Short> {

    static Short minValue = 0;
    static Short maxValue = (short) Byte.MAX_VALUE * 2 + 1;

    public PlcCHAR(Boolean value) {
        super();
        this.value = value ? new Short((short) 1) : new Short((short) 0);
        this.isNullable = false;
    }

    public PlcCHAR(Character value) {
        super();
        Integer val = (int) value;
        if ((val >= minValue) && (val <= maxValue)) {
            this.value = (Short) val.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Character '" + Character.toString(value) + "', Value " + val +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(Byte value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(Short value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(Integer value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(Long value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(Float value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value " + value +
              " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(Double value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = (Short) value.shortValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value " + value +
              " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(BigInteger value) {
        super();
        if ((value.compareTo(BigInteger.valueOf(minValue)) >= 0) && (value.compareTo(BigInteger.valueOf(maxValue)) <= 0)) {
            this.value = (Short) value.shortValue();
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(BigDecimal value) {
        super();
        if ((value.compareTo(BigDecimal.valueOf(minValue)) >= 0) && (value.compareTo(BigDecimal.valueOf(maxValue)) <= 0) && (value.scale() <= 0)) {
            this.value = (Short) value.shortValue();
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcCHAR(String value) {
        super();
        try {
            Short val = (short) value.charAt(0);
            if ((val >= minValue) && (val <= maxValue)) {
                this.value = (short) val;
                this.isNullable = false;
            } else {
                throw new PlcInvalidFieldException("Value " + value +
                  " is out of range " + minValue + " - " + maxValue + " for a " +
                  this.getClass().getSimpleName() + " Value");
            }
        }
        catch(Exception e) {
            throw new PlcInvalidFieldException("Value " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcCHAR(@JsonProperty("value") short value) {
        super();
        if ((value >= minValue) && (value <= maxValue)) {
            this.value = value;
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
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

    public short getCHAR() {
        return value;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return Character.toString(Character.valueOf((char) ((short) value)));
    }

    public byte[] getBytes() {
        byte[] bytes = new byte[1];
        bytes[0] = (byte)(value & 0xff);
        return bytes;
    }

}
