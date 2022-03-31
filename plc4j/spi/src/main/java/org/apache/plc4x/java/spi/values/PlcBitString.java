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

import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcBitString extends PlcList {

    public PlcBitString(boolean[] values) {
        super(toBitString(values));
    }

    public PlcBitString(short byteBitString) {
        super(toBitString(BigInteger.valueOf(byteBitString), 8));
    }

    public PlcBitString(int wordBitString) {
        super(toBitString(BigInteger.valueOf(wordBitString), 16));
    }

    public PlcBitString(long dwordBitString) {
        super(toBitString(BigInteger.valueOf(dwordBitString), 32));
    }

    public PlcBitString(BigInteger lwordBitString) {
        super(toBitString(lwordBitString, 64));
    }

    @Override
    public boolean isShort() {
        return getList().size() == 8;
    }

    @Override
    public short getShort() {
        if(!isShort()) {
            throw new PlcRuntimeException("getShort requires 8 boolean values");
        }
        return fromBitString(getList(), 8).shortValue();
    }

    @Override
    public boolean isInteger() {
        return getList().size() == 16;
    }

    @Override
    public int getInteger() {
        if(!isInteger()) {
            throw new PlcRuntimeException("getInteger requires 16 boolean values");
        }
        return fromBitString(getList(), 16).intValue();
    }

    @Override
    public boolean isLong() {
        return getList().size() == 32;
    }

    @Override
    public long getLong() {
        if(!isLong()) {
            throw new PlcRuntimeException("getLong requires 32 boolean values");
        }
        return fromBitString(getList(), 32).longValue();
    }

    @Override
    public boolean isBigInteger() {
        return getList().size() == 64;
    }

    @Override
    public BigInteger getBigInteger() {
        if(!isBigInteger()) {
            throw new PlcRuntimeException("getBigInteger requires 64 boolean values");
        }
        return fromBitString(getList(), 64);
    }

    @Override
    public Object getObject() {
        boolean[] result = new boolean[getLength()];
        for (int i = 0; i < getLength(); i++) {
            result[i] = getList().get(i).getBoolean();
        }
        return result;
    }

    public static List<PlcValue> toBitString(boolean[] booleanValues) {
        if((booleanValues.length != 8) && (booleanValues.length != 16) &&
            (booleanValues.length != 32) && (booleanValues.length != 64)) {
            throw new PlcRuntimeException("invalid number of values");
        }
        List<PlcValue> values = new ArrayList<>(booleanValues.length);
        for (boolean booleanValue : booleanValues) {
            values.add(new PlcBOOL(booleanValue));
        }
        return values;
    }

    public static List<PlcValue> toBitString(BigInteger bigInteger, int numBits) {
        if(bigInteger.bitCount() > numBits) {
            throw new PlcRuntimeException("value too big");
        }
        // Convert the numeric value into an array of bits.
        List<PlcValue> values = new ArrayList<>(numBits);
        for (int i = numBits - 1; i >= 0; i--) {
            values.add(new PlcBOOL(bigInteger.testBit(i)));
        }
        return values;
    }

    public static BigInteger fromBitString(List<PlcValue> bitString, int numBits) {
        BigInteger bigInteger = new BigInteger(String.valueOf(0L));
        for (int i = 0; i < numBits; i++) {
            final PlcValue plcValue = bitString.get(i);
            if(plcValue.getBoolean()) {
                bigInteger = bigInteger.setBit((numBits - 1) - i);
            }
        }
        return bigInteger;
    }

}
