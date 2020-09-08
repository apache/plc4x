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

package org.apache.plc4x.java.modbus.readwrite;

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
public class PlcINT extends PlcValueAdapter {

    BigInteger minValue = BigInteger.valueOf(Short.MIN_VALUE);
    BigInteger maxValue = BigInteger.valueOf(Short.MAX_VALUE);

    final Short value;
    final boolean isNullable;

    public PlcINT(Short value) {
      if ((BigInteger.valueOf(value).compareTo(minValue) > 0) && (BigInteger.valueOf(value).compareTo(maxValue) < 0)) {
        this.value = value;
        this.isNullable = false;
      } else {
        throw new PlcIncompatibleDatatypeException("");
      }
    }

    public PlcINT(Integer value) {
      if ((BigInteger.valueOf(value).compareTo(minValue) > 0) && (BigInteger.valueOf(value).compareTo(maxValue) < 0)) {
        this.value = (Short) value.shortValue();
        this.isNullable = false;
      } else {
        throw new PlcIncompatibleDatatypeException("");
      }
    }

    public PlcINT(String value) {
      Short val;
      try {
        val = Short.valueOf((String) value).shortValue();
      }
      catch(Exception e) {
        throw new PlcIncompatibleDatatypeException("");
      }
      if ((BigInteger.valueOf(val).compareTo(minValue) > 0) && (BigInteger.valueOf(val).compareTo(maxValue) < 0)) {
        this.value = val;
        this.isNullable = false;
      } else {
        throw new PlcIncompatibleDatatypeException("");
      }
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcINT(@JsonProperty("value") short value) {
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

    @Override
    public Object getObject() {
        return value;
    }

    @Override
    @JsonIgnore
    public int getLength() {
        return 2;
    }

    @Override
    @JsonIgnore
    public boolean isSimple() {
        return true;
    }

    @Override
    @JsonIgnore
    public boolean isNullable() {
        return isNullable;
    }

    @Override
    @JsonIgnore
    public boolean isNull() {
        return isNullable && value == null;
    }

    public byte[] getBytes() {
        byte[] bytes = new byte[2];
        bytes[0] = (byte)((value >> 8) & 0xff);
        bytes[1] = (byte)(value & 0xff);
        return bytes;
    }

    //public static int getList(String[] values) {
    //    List<String> plcIntList = new LinkedList<>();
    //    for (int i = 0; i < values.length; i++) {
    //      plcIntList.add(new PlcINT((String) values[i]));
    //    }
    //    return new PlcList(plcIntList);
    //}

    /**
     * Convenience method to simplify accessing items with a list syntax.
     * @param i item number
     * @return if i == 0 returns itself, otherwise throws an exception.
     */
    @Override
    @JsonIgnore
    public PlcValue getIndex(int i) {
        if(i == 0) {
            return this;
        }
        return super.getIndex(i);
    }

}
