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

public abstract class UnsignedIntLEByteValue extends ByteValue {

    public static final int NUM_BYTES = 4;

    protected final long longValue;

    public UnsignedIntLEByteValue(byte... value) {
        super(value);
        assertLength(NUM_BYTES);
        longValue = getBytes()[3] << 24 | getBytes()[2] << 16 | getBytes()[1] << 8 | getBytes()[0];
    }

    public UnsignedIntLEByteValue(long value) {
        super(ofLong(value));
        checkUnsignedBounds(value, NUM_BYTES);
        longValue = value;
    }

    public UnsignedIntLEByteValue(ByteBuf byteBuf) {
        this(byteBuf.readUnsignedIntLE());
    }

    protected static byte[] ofLong(long value) {
        return ByteBuffer.allocate(NUM_BYTES)
            // LE
            .put((byte) (value & 0xff))
            .put((byte) (value >> 8 & 0xff))
            .put((byte) (value >> 16 & 0xff))
            .put((byte) (value >> 24 & 0xff))
            .array();
    }

    public long getAsLong() {
        return longValue;
    }

    @Override
    public String toString() {
        return String.valueOf(getAsLong());
    }
}
