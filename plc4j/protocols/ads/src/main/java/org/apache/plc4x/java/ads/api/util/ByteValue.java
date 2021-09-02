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
package org.apache.plc4x.java.ads.api.util;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;

import java.math.BigInteger;
import java.util.Arrays;

import static java.util.Objects.requireNonNull;

public class ByteValue implements ByteReadable {

    protected final byte[] value;

    protected ByteValue(byte... value) {
        this.value = requireNonNull(value);
    }

    protected ByteValue(ByteBuf byteBuf) {
        requireNonNull(byteBuf);
        value = new byte[byteBuf.readableBytes()];
        byteBuf.readBytes(value);
    }

    protected void assertLength(int length) {
        if (value.length != length) {
            throw new IllegalArgumentException("Expected length " + length + " got " + value.length);
        }
    }

    public static void checkUnsignedBounds(long value, int numberOfBytes) {
        double upperBound = Math.pow((double) 2, (double) (8 * numberOfBytes));
        if (value < 0 || value >= upperBound) {
            throw new IllegalArgumentException("Value must between 0 and " + upperBound + ". Was " + value);
        }
    }

    public static void checkUnsignedBounds(BigInteger value, int numberOfBytes) {
        BigInteger upperBound = BigInteger.valueOf(2).pow(8 * numberOfBytes);
        if (value.compareTo(BigInteger.ZERO) < 0 || value.compareTo(upperBound) >= 0) {
            throw new IllegalArgumentException("Value must between 0 and " + upperBound + ". Was " + value);
        }
    }

    @Override
    public byte[] getBytes() {
        return value;
    }

    @Override
    public ByteBuf getByteBuf() {
        return Unpooled.wrappedBuffer(value);
    }

    @Override
    public long getCalculatedLength() {
        return value.length;
    }

    public static ByteValue of(byte... values) {
        return new ByteValue(values);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof ByteValue))
            return false;

        ByteValue byteValue = (ByteValue) o;

        return Arrays.equals(value, byteValue.value);
    }

    @Override
    public int hashCode() {
        return Arrays.hashCode(value);
    }

    @Override
    public String toString() {
        // TODO: maybe we could find a way to implement this to string
        return getClass().getSimpleName() + "@" + Integer.toHexString(hashCode()) + "{bytes=" + value.length + "}";
    }
}
