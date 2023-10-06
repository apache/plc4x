/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.values;

import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;

public class PlcSTRING extends PlcSimpleValue<String> {

    public static PlcSTRING of(Object value) {
        if (value instanceof String) {
            return new PlcSTRING((String) value);
        }
        return new PlcSTRING(String.valueOf(value));
    }

    public PlcSTRING(String value) {
        super(value, true);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.STRING;
    }

    @Override
    public boolean isString() {
        return true;
    }

    @Override
    public String getString() {
        return value;
    }

    @Override
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
    public boolean getBoolean() {
        return Boolean.parseBoolean(value);
    }

    @Override
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
    public byte getByte() {
        return Byte.parseByte(value);
    }

    @Override
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
    public short getShort() {
        return Short.parseShort(value);
    }

    @Override
    public boolean isInteger() {
        try {
            Integer.parseInt(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    public int getInteger() {
        return Integer.parseInt(value);
    }

    @Override
    public boolean isLong() {
        try {
            Long.parseLong(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    public long getLong() {
        return Long.parseLong(value);
    }

    @Override
    public boolean isBigInteger() {
        try {
            new BigInteger(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    public BigInteger getBigInteger() {
        return new BigInteger(value);
    }

    @Override
    public boolean isFloat() {
        try {
            Float.parseFloat(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    public float getFloat() {
        return Float.parseFloat(value);
    }

    @Override
    public boolean isDouble() {
        try {
            Double.parseDouble(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    public double getDouble() {
        return Double.parseDouble(value);
    }

    @Override
    public boolean isBigDecimal() {
        try {
            new BigDecimal(value);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    public BigDecimal getBigDecimal() {
        return new BigDecimal(value);
    }

    @Override
    public int getLength() {
        return value.length();
    }

    @Override
    public String toString() {
        return "\"" + value + "\"";
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        String valueString = value;
        writeBuffer.writeString(getClass().getSimpleName(),
            valueString.getBytes(StandardCharsets.UTF_8).length*8,
            valueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
    }

}
