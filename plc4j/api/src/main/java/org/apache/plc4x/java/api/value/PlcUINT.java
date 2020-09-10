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
import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;

import org.apache.plc4x.java.api.value.*;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.LinkedList;
import java.util.List;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcUINT extends PlcIECValue<Integer> {

    BigInteger minValue = BigInteger.valueOf(0);
    BigInteger maxValue = BigInteger.valueOf(((long) Short.MAX_VALUE) * 2 + 1);

    public PlcUINT(Integer value) {
        super();
        if ((BigInteger.valueOf(value).compareTo(minValue) >= 0) && (BigInteger.valueOf(value).compareTo(maxValue) <= 0)) {
          this.value = value;
          this.isNullable = false;
        } else {
          throw new IllegalArgumentException("Value of type " + value + " is out of range " + minValue + " - " + maxValue + " for a PLCUINT Value");
        }
    }

    public PlcUINT(Short value) {
        super();
        if ((BigInteger.valueOf(value).compareTo(minValue) >= 0) && (BigInteger.valueOf(value).compareTo(maxValue) <= 0)) {
          this.value = (Integer) value.intValue();
          this.isNullable = false;
        } else {
          throw new IllegalArgumentException("Value of type " + value + " is out of range " + minValue + " - " + maxValue + " for a PLCUINT Value");
        }
    }

    public PlcUINT(String value) {
        super();
        try {
          Integer val = Integer.parseInt(value);
          if ((BigInteger.valueOf(val).compareTo(minValue) >= 0) && (BigInteger.valueOf(val).compareTo(maxValue) <= 0)) {
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
    public PlcUINT(@JsonProperty("value") int value) {
      super();
      if ((BigInteger.valueOf(value).compareTo(minValue) >= 0) && (BigInteger.valueOf(value).compareTo(maxValue) <= 0)) {
        this.value = value;
        this.isNullable = false;
      } else {
        throw new IllegalArgumentException("Value of type " + value + " is out of range " + minValue + " - " + maxValue + " for a PLCUINT Value");
      }
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

    public int getUINT() {
        return value;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return Integer.toString(value);
    }

    public byte[] getBytes() {
        byte[] bytes = new byte[2];
        bytes[0] = (byte)((value >> 8) & 0xff);
        bytes[1] = (byte)(value & 0xff);
        return bytes;
    }

}
