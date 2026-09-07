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
package org.apache.plc4x.java.spi.buffers.bytebased.encoding;

import java.math.BigDecimal;
import java.math.BigInteger;

public abstract class BaseEncodingDefault implements EncodingDefault {

    @Override
    public byte[] encodeByte(int numBits, byte value) {
        throw new UnsupportedOperationException("encoding byte is not supported by this encoding");
    }

    @Override
    public byte decodeByte(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding byte is not supported by this encoding");
    }

    @Override
    public byte[] encodeShort(int numBits, short value) {
        throw new UnsupportedOperationException("encoding short is not supported by this encoding");
    }

    @Override
    public short decodeShort(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding short is not supported by this encoding");
    }

    @Override
    public byte[] encodeInt(int numBits, int value) {
        throw new UnsupportedOperationException("encoding int is not supported by this encoding");
    }

    @Override
    public int decodeInt(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding int is not supported by this encoding");
    }

    @Override
    public byte[] encodeLong(int numBits, long value) {
        throw new UnsupportedOperationException("encoding long is not supported by this encoding");
    }

    @Override
    public long decodeLong(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding long is not supported by this encoding");
    }

    @Override
    public byte[] encodeBigInteger(int numBits, BigInteger value) {
        throw new UnsupportedOperationException("encoding BigInteger is not supported by this encoding");
    }

    @Override
    public BigInteger decodeBigInteger(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding BigInteger is not supported by this encoding");
    }

    @Override
    public byte[] encodeFloat(int numBits, float value) {
        throw new UnsupportedOperationException("encoding float is not supported by this encoding");
    }

    @Override
    public float decodeFloat(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding float is not supported by this encoding");
    }

    @Override
    public byte[] encodeDouble(int numBits, double value) {
        throw new UnsupportedOperationException("encoding double is not supported by this encoding");
    }

    @Override
    public double decodeDouble(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding double is not supported by this encoding");
    }

    @Override
    public byte[] encodeBigDecimal(int numBits, BigDecimal value) {
        throw new UnsupportedOperationException("encoding string is not supported by this encoding");
    }

    @Override
    public BigDecimal decodeBigDecimal(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding string is not supported by this encoding");
    }

    @Override
    public byte[] encodeString(int numBits, String value) {
        throw new UnsupportedOperationException("encoding string is not supported by this encoding");
    }

    @Override
    public String decodeString(int numBits, byte[] bytes) {
        throw new UnsupportedOperationException("decoding string is not supported by this encoding");
    }
}
