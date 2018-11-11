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

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.BitSet;

public class DefaultBooleanFieldItem extends BaseDefaultFieldItem<Boolean> {

    public DefaultBooleanFieldItem(Boolean... values) {
        super(values);
    }

    @Override
    public Object getObject(int index) {
        return getValue(index);
    }

    @Override
    public boolean isValidBoolean(int index) {
        return getValue(index) != null;
    }

    @Override
    public Boolean getBoolean(int index) {
        if (!isValidBoolean(index)) {
            throw new PlcIncompatibleDatatypeException(Boolean.class, index);
        }
        return getValue(index);
    }

    @Override
    public boolean isValidByte(int index) {
        return isValidBoolean(index);
    }

    @Override
    public Byte getByte(int index) {
        if (!isValidByte(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        return getValue(index) ? (byte) 1 : (byte) 0;
    }

    @Override
    public boolean isValidShort(int index) {
        return isValidBoolean(index);
    }

    @Override
    public Short getShort(int index) {
        if (!isValidShort(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        return getValue(index) ? (short) 1 : (short) 0;
    }

    @Override
    public boolean isValidInteger(int index) {
        return isValidBoolean(index);
    }

    @Override
    public Integer getInteger(int index) {
        if (!isValidInteger(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        return getValue(index) ? 1 : 0;
    }

    @Override
    public boolean isValidLong(int index) {
        return isValidBoolean(index);
    }

    @Override
    public Long getLong(int index) {
        if (!isValidLong(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        return getValue(index) ? 1L : 0L;
    }

    @Override
    public boolean isValidBigInteger(int index) {
        return isValidBoolean(index);
    }

    @Override
    public BigInteger getBigInteger(int index) {
        if (!isValidBigInteger(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        return getValue(index) ? BigInteger.ONE : BigInteger.ZERO;
    }

    @Override
    public boolean isValidFloat(int index) {
        return isValidBoolean(index);
    }

    @Override
    public Float getFloat(int index) {
        if (!isValidFloat(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        return getValue(index) ? 1.0F : 0.0F;
    }

    @Override
    public boolean isValidDouble(int index) {
        return isValidBoolean(index);
    }

    @Override
    public Double getDouble(int index) {
        if (!isValidDouble(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        return getValue(index) ? 1.0D : 0.0D;
    }

    @Override
    public boolean isValidBigDecimal(int index) {
        return isValidBoolean(index);
    }

    @Override
    public BigDecimal getBigDecimal(int index) {
        if (!isValidBigDecimal(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        return getValue(index) ? BigDecimal.ONE : BigDecimal.ZERO;
    }

    public Byte getCompleteByte(int index) {
        if (!isValidByte(index)) {
            throw new PlcIncompatibleDatatypeException(Byte.class, index);
        }
        BitSet bitSet = new BitSet();
        int i = 0;
        if(getValues() != null) {
            for (Boolean value : getValues()) {
                bitSet.set(i, value);
                i++;
            }
        }
        // TODO: In this case the real max index is smaller than the numValues the object reports.
        // Calculate the real number of bytes.
        int numBytes = (i / 8) + ((i % 8 == 0) ? 0 : 1);
        byte[] bytes = bitSet.toByteArray();
        if (numBytes - 1 < index) {
            return null;
        }
        // If the highest level bits are all false, the number of bytes is smaller than it should be,
        // So we have to fix that and return a 0-vale byte instead.
        if((bytes.length < numBytes) && (bytes.length <= index)) {
            return 0x00;
        }
        return bytes[index];
    }

    // TODO: implement other methods according to getByte
}

