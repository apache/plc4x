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
package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.util.regex.Pattern;
import java.util.stream.Stream;

/**
 * The AMSNetId consists of 6 bytes and addresses the transmitter or receiver. One possible AMSNetId would be e.g. "172.16.17.10.1.1". The storage arrangement in this example is as follows:
 * <p>
 * _____0     1     2     3     4     5
 * __+-----------------------------------+
 * 0 | 127 |  16 |  17 |  10 |   1 |   1 |
 * __+-----------------------------------+
 * <p>
 * <p>
 * The AMSNetId is purely logical and has usually no relation to the IP address. The AMSNetId is configured at the target system. At the PC for this the TwinCAT System Control is used. If you use other hardware, see the considering documentation for notes about settings of the AMS NetId.
 */
public class AMSNetId extends ByteValue {

    public static final Pattern AMS_NET_ID_PATTERN =
        Pattern.compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}");

    public static final int NUM_BYTES = 6;

    private AMSNetId(byte... values) {
        super(values);
        assertLength(NUM_BYTES);
    }

    public static AMSNetId of(byte... values) {
        return new AMSNetId(values);
    }

    public static AMSNetId of(int octed1, int octed2, int octed3, int octed4, int octed5, int octed6) {
        return new AMSNetId((byte) octed1, (byte) octed2, (byte) octed3, (byte) octed4, (byte) octed5, (byte) octed6);
    }

    public static AMSNetId of(String address) {
        if (!AMS_NET_ID_PATTERN.matcher(address).matches()) {
            throw new IllegalArgumentException(address + " must match " + AMS_NET_ID_PATTERN);
        }
        String[] split = address.split("\\.");
        byte[] bytes = ArrayUtils.toPrimitive(Stream.of(split).map(Integer::parseInt).map(Integer::byteValue).toArray(Byte[]::new));
        return new AMSNetId(bytes);
    }

    public static AMSNetId of(ByteBuf byteBuf) {
        byte[] values = new byte[NUM_BYTES];
        byteBuf.readBytes(values);
        return of(values);
    }

    @Override
    public String toString() {
        byte[] bytes = getBytes();
        return bytes[0] + "." + bytes[1] + "." + bytes[2] + "." + bytes[3] + "." + bytes[4] + "." + bytes[5];
    }
}
