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
package org.apache.plc4x.java.ads.api.util;

import io.netty.buffer.ByteBuf;

import java.nio.ByteBuffer;

public abstract class UnsignedShortLEByteValue extends ByteValue {

    public static final int NUM_BYTES = 2;

    private final int intValue;

    protected UnsignedShortLEByteValue(byte... value) {
        super(value);
        assertLength(NUM_BYTES);
        intValue = (getBytes()[1] << 8 | getBytes()[0]) & 0xff;
    }

    protected UnsignedShortLEByteValue(int value) {
        super(ofInt(value));
        checkUnsignedBounds(value, NUM_BYTES);
        intValue = value;
    }

    protected UnsignedShortLEByteValue(String value) {
        this(Integer.parseInt(value));
    }

    protected UnsignedShortLEByteValue(ByteBuf byteBuf) {
        this(byteBuf.readUnsignedShortLE());
    }

    private static byte[] ofInt(long value) {
        return ByteBuffer.allocate(NUM_BYTES)
            // LE
            .put((byte) (value & 0xff))
            .put((byte) (value >> 8 & 0xff))
            .array();
    }

    public int getAsInt() {
        return intValue;
    }

    @Override
    public String toString() {
        return getClass().getSimpleName() + "{" +
            "intValue=" + getAsInt() +
            "} " + super.toString();
    }
}
