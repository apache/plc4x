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
package org.apache.plc4x.java.ads.api.commands.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.math.BigInteger;
import java.nio.ByteBuffer;

public class TimeStamp extends ByteValue {

    public static final int NUM_BYTES = 8;

    TimeStamp(byte... values) {
        super(values);
        assertLength(NUM_BYTES);
    }

    public static TimeStamp of(long value) {
        return of(BigInteger.valueOf(value));
    }

    public static TimeStamp of(BigInteger value) {
        checkUnsignedBounds(value, NUM_BYTES);
        byte[] valueBytes = value.toByteArray();
        int length = valueBytes.length;
        return new TimeStamp(ByteBuffer.allocate(NUM_BYTES)
            // LE
            .put(length > 0 ? valueBytes[0] : 0)
            .put(length > 1 ? valueBytes[1] : 0)
            .put(length > 2 ? valueBytes[2] : 0)
            .put(length > 3 ? valueBytes[3] : 0)

            .put(length > 4 ? valueBytes[4] : 0)
            .put(length > 5 ? valueBytes[5] : 0)
            .put(length > 6 ? valueBytes[6] : 0)
            .put(length > 7 ? valueBytes[7] : 0)
            .array());
    }

    public static TimeStamp of(byte... values) {
        return new TimeStamp(values);
    }

    public static TimeStamp of(ByteBuf byteBuf) {
        byte[] values = new byte[NUM_BYTES];
        byteBuf.readBytes(values);
        return of(values);
    }
}
