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

package org.apache.plc4x.java.utils;

import java.math.BigDecimal;
import java.math.BigInteger;

public class IoBuffer {

    public int getPos() {
        return 0;
    }

    public boolean readBit() {
        return false;
    }

    public byte readUnsignedByte(int bitLength) {
        return 0;
    }

    public short readUnsignedShort(int bitLength) {
        return 0;
    }

    public int readUnsignedInt(int bitLength) {
        return 0;
    }

    public long readUnsignedLong(int bitLength) {
        return 0;
    }

    public BigInteger readUnsignedBigInteger(int bitLength) {
        return BigInteger.ZERO;
    }

    public byte readByte(int bitLength) {
        return 0;
    }

    public short readShort(int bitLength) {
        return 0;
    }

    public int readInt(int bitLength) {
        return 0;
    }

    public long readLong(int bitLength) {
        return 0;
    }

    public BigInteger readBigInteger(int bitLength) {
        return BigInteger.ZERO;
    }

    public float readFloat(int bitLength) {
        return 0.0f;
    }

    public double readDouble(int bitLength) {
        return 0.0;
    }

    public BigDecimal readBigDecimal(int bitLength) {
        return BigDecimal.ZERO;
    }

}
