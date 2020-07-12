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

import java.math.BigDecimal;
import java.math.BigInteger;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcString extends PlcSimpleValue<String> {

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcString(@JsonProperty("value") String value) {
        super(value, true);
    }

    @Override
    @JsonIgnore
    public boolean isString() {
        return true;
    }

    @Override
    @JsonIgnore
    public String getString() {
        return value;
    }

    @Override
    @JsonIgnore
    public boolean isBoolean() {
        // TODO: Implement ...
        return super.isBoolean();
    }

    @Override
    @JsonIgnore
    public boolean getBoolean() {
        return Boolean.parseBoolean(value);
    }

    @Override
    @JsonIgnore
    public boolean isByte() {
        // TODO: Implement ...
        return super.isByte();
    }

    @Override
    @JsonIgnore
    public byte getByte() {
        return Byte.parseByte(value);
    }

    @Override
    @JsonIgnore
    public boolean isShort() {
        // TODO: Implement ...
        return super.isShort();
    }

    @Override
    @JsonIgnore
    public short getShort() {
        return Short.parseShort(value);
    }

    @Override
    @JsonIgnore
    public boolean isInteger() {
        // TODO: Implement ...
        return super.isInteger();
    }

    @Override
    @JsonIgnore
    public int getInteger() {
        return Integer.parseInt(value);
    }

    @Override
    @JsonIgnore
    public boolean isLong() {
        // TODO: Implement ...
        return super.isLong();
    }

    @Override
    @JsonIgnore
    public long getLong() {
        return Long.parseLong(value);
    }

    @Override
    @JsonIgnore
    public boolean isBigInteger() {
        // TODO: Implement ...
        return super.isBigInteger();
    }

    @Override
    @JsonIgnore
    public BigInteger getBigInteger() {
        return new BigInteger(value);
    }

    @Override
    @JsonIgnore
    public boolean isFloat() {
        // TODO: Implement ...
        return super.isFloat();
    }

    @Override
    @JsonIgnore
    public float getFloat() {
        return Float.parseFloat(value);
    }

    @Override
    @JsonIgnore
    public boolean isDouble() {
        // TODO: Implement ...
        return super.isDouble();
    }

    @Override
    @JsonIgnore
    public double getDouble() {
        return Double.parseDouble(value);
    }

    @Override
    @JsonIgnore
    public boolean isBigDecimal() {
        // TODO: Implement ...
        return super.isBigDecimal();
    }

    @Override
    @JsonIgnore
    public BigDecimal getBigDecimal() {
        return new BigDecimal(value);
    }

    @Override
    @JsonIgnore
    public int getLength() {
        return value.length();
    }

    @Override
    @JsonIgnore
    public String toString() {
        return "\"" + value + "\"";
    }

}
