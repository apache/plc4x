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

}
