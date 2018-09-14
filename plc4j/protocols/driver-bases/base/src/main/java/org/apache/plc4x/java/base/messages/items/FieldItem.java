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

import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;

public abstract class FieldItem<T> {

    private T[] values;

    protected FieldItem(T[] values) {
        this.values = values;
    }

    public int getNumberOfValues() {
        return values.length;
    }

    public abstract Object getObject(int index);

    public boolean isValidBoolean(int index) {
        return false;
    }

    public Boolean getBoolean(int index) {
        return null;
    }

    public boolean isValidByte(int index) {
        return false;
    }

    public Byte getByte(int index) {
        return null;
    }

    public boolean isValidShort(int index) {
        return false;
    }

    public Short getShort(int index) {
        return null;
    }

    public boolean isValidInteger(int index) {
        return false;
    }

    public Integer getInteger(int index) {
        return null;
    }

    public boolean isValidBigInteger(int index) {
        return false;
    }

    public BigInteger getBigInteger(int index) {
        return null;
    }

    public boolean isValidLong(int index) {
        return false;
    }

    public Long getLong(int index) {
        return null;
    }

    public boolean isValidFloat(int index) {
        return false;
    }

    public Float getFloat(int index) {
        return null;
    }

    public boolean isValidDouble(int index) {
        return false;
    }

    public Double getDouble(int index) {
        return null;
    }

    public boolean isValidString(int index) {
        return false;
    }

    public String getString(int index) {
        return null;
    }

    public boolean isValidTime(int index) {
        return false;
    }

    public LocalTime getTime(int index) {
        return null;
    }

    public boolean isValidDate(int index) {
        return false;
    }

    public LocalDate getDate(int index) {
        return null;
    }

    public boolean isValidDateTime(int index) {
        return false;
    }

    public LocalDateTime getDateTime(int index) {
        return null;
    }

    public T[] getValues() {
        return values;
    }

    protected T getValue(int index) {
        if (index < 0 || (index >= values.length)) {
            return null;
        }
        return values[index];
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof FieldItem)) {
            return false;
        }
        FieldItem<?> fieldItem = (FieldItem<?>) o;
        return Arrays.equals(values, fieldItem.values);
    }

    @Override
    public int hashCode() {
        return Arrays.hashCode(values);
    }

    @Override
    public String toString() {
        return "FieldItem{" +
            "values=" + Arrays.toString(values) +
            '}';
    }
}
