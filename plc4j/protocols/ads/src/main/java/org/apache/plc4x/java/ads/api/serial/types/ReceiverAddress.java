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
package org.apache.plc4x.java.ads.api.serial.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import static java.lang.Integer.toHexString;
import static org.apache.commons.lang3.StringUtils.leftPad;

public class ReceiverAddress extends ByteValue {

    public static final int NUM_BYTES = 1;

    public static final ReceiverAddress RS232_COMM_ADDRESS = ReceiverAddress.of((byte) 0);

    private ReceiverAddress(byte value) {
        super(value);
    }

    private ReceiverAddress(byte[] value) {
        super(value);
    }

    private ReceiverAddress(ByteBuf byteBuf) {
        this(byteBuf.readByte());
    }

    public static ReceiverAddress of(byte value) {
        return new ReceiverAddress(value);
    }

    public static ReceiverAddress of(byte... value) {
        return new ReceiverAddress(value);
    }

    public static ReceiverAddress of(String value) {
        return new ReceiverAddress(Byte.valueOf(value));
    }

    public static ReceiverAddress of(ByteBuf byteBuf) {
        return new ReceiverAddress(byteBuf);
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