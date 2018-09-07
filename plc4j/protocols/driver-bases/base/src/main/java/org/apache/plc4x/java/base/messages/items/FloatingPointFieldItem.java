/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.base.messages.items;

public class FloatingPointFieldItem extends FieldItem<Double> {

    public FloatingPointFieldItem(Double... values) {
        super(values);
    }

    @Override
    public Object getObject(int index) {
        return getDouble(index);
    }

    @Override
    public boolean isValidBoolean(int index) {
        return (getValue(index) != null);
    }

    @Override
    public Boolean getBoolean(int index) {
        if (isValidBoolean(index)) {
            return getValue(index) != 0L;
        }
        return null;
    }

    @Override
    public boolean isValidByte(int index) {
        Double value = getValue(index);
        return (value != null) && (value <= Byte.MAX_VALUE) && (value >= Byte.MIN_VALUE);
    }

    @Override
    public Byte getByte(int index) {
        if (isValidByte(index)) {
            return getValue(index).byteValue();
        }
        return null;
    }

    @Override
    public boolean isValidShort(int index) {
        Double value = getValue(index);
        return (value != null) && (value <= Short.MAX_VALUE) && (value >= Short.MIN_VALUE);
    }

    @Override
    public Short getShort(int index) {
        if (isValidShort(index)) {
            return getValue(index).shortValue();
        }
        return null;
    }

    @Override
    public boolean isValidInteger(int index) {
        Double value = getValue(index);
        return (value != null) && (value <= Integer.MAX_VALUE) && (value >= Integer.MIN_VALUE);
    }

    @Override
    public Integer getInteger(int index) {
        if (isValidInteger(index)) {
            return getValue(index).intValue();
        }
        return null;
    }

    @Override
    public boolean isValidLong(int index) {
        Double value = getValue(index);
        return (value != null) && (value <= Long.MAX_VALUE) && (value >= Long.MIN_VALUE);
    }

    @Override
    public Long getLong(int index) {
        if (isValidLong(index)) {
            return getValue(index).longValue();
        }
        return null;
    }

    @Override
    public boolean isValidFloat(int index) {
        Double value = getValue(index);
        return (value != null) && (value <= Float.MAX_VALUE) && (value >= Float.MIN_VALUE);
    }

    @Override
    public Float getFloat(int index) {
        if (isValidFloat(index)) {
            return getValue(index).floatValue();
        }
        return null;
    }

    @Override
    public boolean isValidDouble(int index) {
        return (getValue(index) != null);
    }

    @Override
    public Double getDouble(int index) {
        if (isValidDouble(index)) {
            return getValue(index);
        }
        return null;
    }

}


