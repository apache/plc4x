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
package org.apache.plc4x.java.ads.api.util;

import io.netty.buffer.ByteBuf;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

import static java.lang.Integer.toHexString;
import static org.apache.commons.lang3.StringUtils.leftPad;

public abstract class UnsignedShortLEByteValue extends ByteValue {

    public static final int UNSIGNED_SHORT_LE_NUM_BYTES = 2;

    private final int intValue;

    protected UnsignedShortLEByteValue(byte... value) {
        super(value);
        assertLength(UNSIGNED_SHORT_LE_NUM_BYTES);
        intValue = ByteBuffer.wrap(value)
            .order(ByteOrder.LITTLE_ENDIAN)
            .getShort();
    }

    protected UnsignedShortLEByteValue(int value) {
        super(ofInt(value));
        intValue = value;
    }

    protected UnsignedShortLEByteValue(String value) {
        this(Integer.parseInt(value));
    }

    protected UnsignedShortLEByteValue(ByteBuf byteBuf) {
        this(byteBuf.readUnsignedShortLE());
    }

    private static byte[] ofInt(int value) {
        checkUnsignedBounds(value, UNSIGNED_SHORT_LE_NUM_BYTES);
        return ByteBuffer.allocate(UNSIGNED_SHORT_LE_NUM_BYTES)
            .order(ByteOrder.LITTLE_ENDIAN)
            .putShort((short) (value & 0xffff))
            .array();
    }

    public int getAsInt() {
        return intValue;
    }

    @Override
    public long getCalculatedLength() {
        return UNSIGNED_SHORT_LE_NUM_BYTES;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof UnsignedShortLEByteValue)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }

        UnsignedShortLEByteValue that = (UnsignedShortLEByteValue) o;

        return intValue == that.intValue;
    }

    @Override
    public int hashCode() {
        int result = super.hashCode();
        result = 31 * result + intValue;
        return result;
    }

    @Override
    public String toString() {
        return super.toString() + "{" +
            "intValue=" + getAsInt() +
            ",hexValue=0x" + leftPad(toHexString(getAsInt()), UNSIGNED_SHORT_LE_NUM_BYTES * 2, "0") +
            "}";
    }
}
