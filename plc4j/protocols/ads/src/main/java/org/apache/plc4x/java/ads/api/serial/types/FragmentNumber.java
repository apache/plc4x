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
package org.apache.plc4x.java.ads.api.serial.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import static java.lang.Integer.toHexString;
import static org.apache.commons.lang3.StringUtils.leftPad;

public class FragmentNumber extends ByteValue {

    public static final int NUM_BYTES = 1;

    private FragmentNumber(byte value) {
        super(value);
    }

    private FragmentNumber(byte[] value) {
        super(value);
    }

    private FragmentNumber(ByteBuf byteBuf) {
        this(byteBuf.readByte());
    }

    public static FragmentNumber of(byte value) {
        return new FragmentNumber(value);
    }

    public static FragmentNumber of(byte... value) {
        return new FragmentNumber(value);
    }

    public static FragmentNumber of(String value) {
        return new FragmentNumber(Byte.valueOf(value));
    }

    public static FragmentNumber of(ByteBuf byteBuf) {
        return new FragmentNumber(byteBuf);
    }

    public byte getAsByte() {
        return value[0];
    }

    @Override
    public String toString() {
        return super.toString() + "{" +
            "byteValue=" + (getAsByte() & 0xFF) +
            ",hexValue=0x" + leftPad(toHexString(getAsByte() & 0xFF), NUM_BYTES * 2, "0") +
            "}";
    }
}