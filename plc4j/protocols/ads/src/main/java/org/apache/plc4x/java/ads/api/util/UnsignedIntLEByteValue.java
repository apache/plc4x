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

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

import static java.lang.Long.toHexString;
import static org.apache.commons.lang3.StringUtils.leftPad;

public abstract class UnsignedIntLEByteValue extends ByteValue {

    public static final int UNSIGNED_INT_LE_NUM_BYTES = 4;

    private final long longValue;

    protected UnsignedIntLEByteValue(byte... value) {
        super(value);
        assertLength(UNSIGNED_INT_LE_NUM_BYTES);
        longValue = ByteBuffer.wrap(value)
            .order(ByteOrder.LITTLE_ENDIAN)
            .getInt();
    }

    protected UnsignedIntLEByteValue(long value) {
        super(ofLong(value));
        longValue = value;
    }

    protected UnsignedIntLEByteValue(String value) {
        this(Long.parseLong(value));
    }

    protected UnsignedIntLEByteValue(ByteBuf byteBuf) {
        this(byteBuf.readUnsignedIntLE());
    }

    private static byte[] ofLong(long value) {
        checkUnsignedBounds(value, UNSIGNED_INT_LE_NUM_BYTES);
        return ByteBuffer.allocate(UNSIGNED_INT_LE_NUM_BYTES)
            .order(ByteOrder.LITTLE_ENDIAN)
            .putInt((int) (value & 0xffff_ffff))
            .array();
    }

    public long getAsLong() {
        return longValue;
    }

    @Override
    public long getCalculatedLength() {
        return UNSIGNED_INT_LE_NUM_BYTES;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof UnsignedIntLEByteValue)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }

        UnsignedIntLEByteValue that = (UnsignedIntLEByteValue) o;

        return longValue == that.longValue;
    }

    @Override
    public int hashCode() {
        int result = super.hashCode();
        result = 31 * result + (int) (longValue ^ (longValue >>> 32));
        return result;
    }

    @Override
    public String toString() {
        return super.toString() + "{" +
            "longValue=" + getAsLong() +
            ",hexValue=0x" + leftPad(toHexString(getAsLong()), UNSIGNED_INT_LE_NUM_BYTES * 2, "0") +
            "} ";
    }
}
