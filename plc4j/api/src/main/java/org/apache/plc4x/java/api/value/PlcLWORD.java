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

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;

import java.math.BigDecimal;
import java.math.BigInteger;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcLWORD extends PlcIECValue<BigInteger> {

    static BigInteger minValue = BigInteger.valueOf(0);
    static BigInteger maxValue = BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.valueOf(2)).add(BigInteger.valueOf(1));

    public PlcLWORD(Boolean value) {
        super();
        this.value = value ? BigInteger.valueOf(1) : BigInteger.valueOf(0);
        this.isNullable = false;
    }

    public PlcLWORD(Byte value) {
        super();
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) >= 0) && (val.compareTo(maxValue) <= 0)) {
            this.value = val;
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLWORD(Short value) {
        super();
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) >= 0) && (val.compareTo(maxValue) <= 0)) {
            this.value = val;
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLWORD(Integer value) {
        super();
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) >= 0) && (val.compareTo(maxValue) <= 0)) {
            this.value = val;
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLWORD(Long value) {
        super();
        BigInteger val = BigInteger.valueOf(value);
        if ((val.compareTo(minValue) >= 0) && (val.compareTo(maxValue) <= 0)) {
            this.value = val;
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLWORD(Float value) {
        super();
        try {
            BigInteger val = BigDecimal.valueOf(value).toBigInteger();
            if ((val.compareTo(minValue) >= 0) && (val.compareTo(maxValue) <= 0)) {
                this.value = val;
                this.isNullable = false;
            } else {
                throw new PlcInvalidFieldException("Value of type " + value +
                  " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
                  this.getClass().getSimpleName() + " Value");
            }
        } catch(Exception e) {
            throw new IllegalArgumentException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLWORD(Double value) {
        super();
        try {
            BigInteger val = BigDecimal.valueOf(value).toBigInteger();
            if ((val.compareTo(minValue) >= 0) && (val.compareTo(maxValue) <= 0)) {
                this.value = val;
                this.isNullable = false;
            } else {
                throw new PlcInvalidFieldException("Value of type " + value +
                  " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
                  this.getClass().getSimpleName() + " Value");
            }
        } catch(Exception e) {
            throw new IllegalArgumentException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLWORD(BigInteger value) {
        super();
        if ((value.compareTo(minValue) >= 0) && (value.compareTo(maxValue) <= 0)) {
            this.value = value;
            this.isNullable = false;
        } else {
            throw new PlcInvalidFieldException("Value of type " + value +
              " is out of range " + minValue + " - " + maxValue + " for a " +
              this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLWORD(BigDecimal value) {
        super();
        try {
            BigInteger val = value.toBigInteger();
            if ((val.compareTo(minValue) >= 0) && (val.compareTo(maxValue) <= 0)) {
                this.value = val;
                this.isNullable = false;
            } else {
                throw new PlcInvalidFieldException("Value of type " + value +
                  " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
                  this.getClass().getSimpleName() + " Value");
            }
        } catch(Exception e) {
            throw new IllegalArgumentException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    public PlcLWORD(String value) {
        super();
        try {
            BigInteger val = new BigInteger(value);
            if ((val.compareTo(minValue) >= 0) && (val.compareTo(maxValue) <= 0)) {
                this.value = val;
                this.isNullable = false;
            } else {
                throw new PlcInvalidFieldException("Value of type " + value +
                  " is out of range " + minValue + " - " + maxValue + " or has decimal places for a " +
                  this.getClass().getSimpleName() + " Value");
            }
        } catch(Exception e) {
            throw new IllegalArgumentException("Value of type " + value +
            " is out of range " + minValue + " - " + maxValue + " for a " +
            this.getClass().getSimpleName() + " Value");
        }
    }

    @Override
    @JsonIgnore
    public boolean isBigInteger() {
        return true;
    }

    @Override
    @JsonIgnore
    public BigInteger getBigInteger() {
        return value;
    }

    @JsonIgnore
    public BigInteger getLWORD() {
        return value;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return value.toString();
    }

    @JsonIgnore
    public byte[] getBytes() {
        byte[] tmp = value.toByteArray();
        byte[] bytes = new byte[8];
        for (int i = 0; i < bytes.length; i++) {
            if (i >= (bytes.length - tmp.length)) {
                bytes[i] = tmp[i - (bytes.length - tmp.length)];
            } else {
                bytes[i] = 0x00;
            }
        }
        return bytes;
    }

}
