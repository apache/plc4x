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
public class PlcLINT extends PlcIECValue<Long> {

    static Long minValue = (long) 0;
    static Long maxValue = (long) Long.MAX_VALUE;

    public PlcLINT(Boolean value) {
        super();
        this.value = value ? (long) 1 : (long) 0;
        this.isNullable = false;
    }

    public PlcLINT(Byte value) {
        super();
        this.value = (Long) value.longValue();
        this.isNullable = false;
    }

    public PlcLINT(Short value) {
        super();
        this.value = (Long) value.longValue();
        this.isNullable = false;
    }

    public PlcLINT(Integer value) {
        super();
        this.value = (Long) value.longValue();
        this.isNullable = false;
    }

    public PlcLINT(Long value) {
        super();
        this.value = value;
        this.isNullable = false;
    }

    public PlcLINT(Float value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = (Long) value.longValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLINT(Double value) {
        super();
        if ((value >= minValue) && (value <= maxValue) && (value % 1 == 0)) {
            this.value = (Long) value.longValue();
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLINT(BigInteger value) {
        super();
        if ((value.compareTo(BigInteger.valueOf(minValue)) >= 0) && (value.compareTo(BigInteger.valueOf(maxValue)) <= 0)) {
            this.value = (Long) value.longValue();
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLINT(BigDecimal value) {
        super();
        if ((value.compareTo(BigDecimal.valueOf(minValue)) >= 0) && (value.compareTo(BigDecimal.valueOf(maxValue)) <= 0) && (value.scale() <= 0)) {
            this.value = (Long) value.longValue();
            this.isNullable = true;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLINT(String value) {
        super();
        try {
            Long val = Long.parseLong(value);
            this.value = val;
            this.isNullable = false;
        }
        catch(Exception e) {
            throw new IllegalArgumentException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcLINT(@JsonProperty("value") long value) {
        super();
        this.value = value;
        this.isNullable = false;
    }

    @Override
    @JsonIgnore
    public boolean isLong() {
        return true;
    }

    @Override
    @JsonIgnore
    public long getLong() {
        return value;
    }

    public long getLINT() {
        return value;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return Long.toString(value);
    }

    public byte[] getBytes() {
        byte[] bytes = new byte[8];
        bytes[0] = (byte)((value >> 56) & 0xff);
        bytes[1] = (byte)((value >> 48) & 0xff);
        bytes[2] = (byte)((value >> 40) & 0xff);
        bytes[3] = (byte)((value >> 32) & 0xff);
        bytes[4] = (byte)((value >> 24) & 0xff);
        bytes[5] = (byte)((value >> 16) & 0xff);
        bytes[6] = (byte)((value >> 8) & 0xff);
        bytes[7] = (byte)(value & 0xff);
        return bytes;
    }

}
