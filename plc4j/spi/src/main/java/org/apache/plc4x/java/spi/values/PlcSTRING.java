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

package org.apache.plc4x.java.spi.values;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.w3c.dom.Element;

import java.math.BigDecimal;
import java.math.BigInteger;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcSTRING extends PlcSimpleValue<String> {

    static int maxLength = 254;

    public static PlcSTRING of(Object value) {
        if (value instanceof String) {
            return new PlcSTRING((String) value);
        }
        return new PlcSTRING(String.valueOf(value));
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcSTRING(@JsonProperty("value") String value) {
        super(value, true);
        if (value.length() > maxLength) {
            throw new IllegalArgumentException(
                String.format("String length %d exceeds allowed maximum for type String (max %d)", value.length(), maxLength));
        }
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
    @SuppressWarnings("all")
    public boolean isBoolean() {
        try {
            Boolean.parseBoolean(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    @JsonIgnore
    public boolean getBoolean() {
        return Boolean.parseBoolean(value);
    }

    @Override
    @JsonIgnore
    @SuppressWarnings("all")
    public boolean isByte() {
        try {
            Byte.parseByte(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    @JsonIgnore
    public byte getByte() {
        return Byte.parseByte(value);
    }

    @Override
    @JsonIgnore
    @SuppressWarnings("all")
    public boolean isShort() {
        try {
            Short.parseShort(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    @JsonIgnore
    public short getShort() {
        return Short.parseShort(value);
    }

    @Override
    @JsonIgnore
    public boolean isInteger() {
        try {
            Integer.parseInt(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    @JsonIgnore
    public int getInteger() {
        return Integer.parseInt(value);
    }

    @Override
    @JsonIgnore
    public boolean isLong() {
        try {
            Long.parseLong(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    @JsonIgnore
    public long getLong() {
        return Long.parseLong(value);
    }

    @Override
    @JsonIgnore
    public boolean isBigInteger() {
        try {
            new BigInteger(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    @JsonIgnore
    public BigInteger getBigInteger() {
        return new BigInteger(value);
    }

    @Override
    @JsonIgnore
    public boolean isFloat() {
        try {
            Float.parseFloat(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    @JsonIgnore
    public float getFloat() {
        return Float.parseFloat(value);
    }

    @Override
    @JsonIgnore
    public boolean isDouble() {
        try {
            Double.parseDouble(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    @JsonIgnore
    public double getDouble() {
        return Double.parseDouble(value);
    }

    @Override
    @JsonIgnore
    public boolean isBigDecimal() {
        try {
            new BigDecimal(value);
            return true;
        } catch (Exception e) {
            return false;
        }
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

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        // TODO: Implement
    }

    @Override
    public void xmlSerialize(Element parent) {
        // TODO: Implement
    }

}
