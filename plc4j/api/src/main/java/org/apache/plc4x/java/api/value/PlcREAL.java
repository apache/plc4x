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
public class PlcREAL extends PlcIECValue<Float> {

    Float minValue = (Float) Float.MIN_VALUE;
    Float maxValue = (Float) Float.MAX_VALUE;

    public PlcREAL(Float value) {
        super();
        if ((value > minValue) && (value < maxValue)) {
          this.value = value;
          this.isNullable = false;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcREAL(Integer value) {
        super();
        if ((value > minValue) && (value < maxValue)) {
          this.value = (Float) value.floatValue();
          this.isNullable = false;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcREAL(Short value) {
        super();
        if ((value > minValue) && (value < maxValue)) {
          this.value = (Float) value.floatValue();
          this.isNullable = false;
        } else {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcREAL(String value) {
        super();
        try {
          Float val = Float.parseFloat(value);
          if ((val > minValue) && (val < maxValue)) {
            this.value = val;
            this.isNullable = false;
          } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
          }
        }
        catch(Exception e) {
          throw new PlcInvalidFieldException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcREAL(@JsonProperty("value") float value) {
      super();
      if ((value > minValue) && (value < maxValue)) {
        this.value = new Float(value);
        this.isNullable = false;
      } else {
        throw new PlcInvalidFieldException("Value of type " + value +
          " is out of range " + minValue + " - " + maxValue + " for a " +
          this.getClass().getSimpleName() + " Value");
      }
    }

    @Override
    @JsonIgnore
    public boolean isFloat() {
        return true;
    }

    @Override
    @JsonIgnore
    public float getFloat() {
        return value;
    }

    public float getREAL() {
        return value;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return Float.toString(value);
    }

    public byte[] getBytes() {
        int intBits =  Float.floatToIntBits(value);
	      return new byte[] { (byte) (intBits >> 24), (byte) (intBits >> 16), (byte) (intBits >> 8), (byte) (intBits) };
    }

}
