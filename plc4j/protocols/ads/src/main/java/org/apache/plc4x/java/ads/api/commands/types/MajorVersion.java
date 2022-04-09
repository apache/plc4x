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
package org.apache.plc4x.java.ads.api.commands.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.util.Arrays;

public class MajorVersion extends ByteValue {

    public static final int NUM_BYTES = 1;

    private MajorVersion(byte... values) {
        super(values);
        assertLength(NUM_BYTES);
    }

    public static MajorVersion of(byte... values) {
        return new MajorVersion(values);
    }

    public static MajorVersion of(int value) {
        return new MajorVersion((byte) value);
    }

    public static MajorVersion of(String value) {
        return of(Integer.parseInt(value));
    }

    public static MajorVersion of(ByteBuf byteBuf) {
        byte[] values = new byte[NUM_BYTES];
        byteBuf.readBytes(values);
        return of(values);
    }

    @Override
    public long getCalculatedLength() {
        return NUM_BYTES;
    }

    public byte getAsByte() {
        return getBytes()[0];
    }

    @Override
    public String toString() {
        return "MajorVersion{" +
            "value=" + Arrays.toString(value) +
            "} " + super.toString();
    }
}
