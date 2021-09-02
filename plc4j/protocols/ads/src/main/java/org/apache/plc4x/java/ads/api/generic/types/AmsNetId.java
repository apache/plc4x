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
package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.util.regex.Pattern;
import java.util.stream.Stream;

/**
 * It is not only possible to exchange data between TwinCAT modules on one PC, it is even possible to do so by ADS
 * methods between multiple TwinCAT PC's on the network.
 * <p>
 * Every PC on the network can be uniquely identified by a TCP/IP address, such as "172.1.2.16". The AdsAmsNetId is an
 * extension of the TCP/IP address and identifies a TwinCAT message router, e.g. "172.1.2.16.1.1". TwinCAT message
 * routers exist on every TwinCAT PC, and on every Beckhoff BCxxxx bus controller (e.g. BC3100, BC8100, BC9000, ...).
 * <p>
 * The AmsNetId consists of 6 bytes and addresses the transmitter or receiver. One possible AmsNetId would be e.g.
 * "172.16.17.10.1.1". The storage arrangement in this example is as follows:
 * <p>
 * _____0     1     2     3     4     5
 * __+-----------------------------------+
 * 0 | 127 |  16 |  17 |  10 |   1 |   1 |
 * __+-----------------------------------+
 * <p>
 * The AmsNetId is purely logical and has usually no relation to the IP address. The AmsNetId is configured at the
 * target system. At the PC for this the TwinCAT System Control is used. If you use other hardware, see the considering
 * documentation for notes about settings of the AMS NetId.
 *
 * @see <a href="https://infosys.beckhoff.com/content/1033/tcadscommon/html/tcadscommon_identadsdevice.htm?id=3991659524769593444">ADS device identification</a>
 */
public class AmsNetId extends ByteValue {

    public static final Pattern AMS_NET_ID_PATTERN =
        Pattern.compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}");

    public static final int NUM_BYTES = 6;

    private AmsNetId(byte... values) {
        super(values);
        assertLength(NUM_BYTES);
    }

    public static AmsNetId of(byte... values) {
        return new AmsNetId(values);
    }

    public static AmsNetId of(int octed1, int octed2, int octed3, int octed4, int octed5, int octed6) {
        return new AmsNetId((byte) octed1, (byte) octed2, (byte) octed3, (byte) octed4, (byte) octed5, (byte) octed6);
    }

    public static AmsNetId of(String address) {
        if (!AMS_NET_ID_PATTERN.matcher(address).matches()) {
            throw new IllegalArgumentException(address + " must match " + AMS_NET_ID_PATTERN);
        }
        String[] split = address.split("\\.");
        byte[] bytes = ArrayUtils.toPrimitive(Stream.of(split).map(Integer::parseInt).map(Integer::byteValue).toArray(Byte[]::new));
        return new AmsNetId(bytes);
    }

    public static AmsNetId of(ByteBuf byteBuf) {
        byte[] values = new byte[NUM_BYTES];
        byteBuf.readBytes(values);
        return of(values);
    }

    @Override
    public String toString() {
        byte[] bytes = getBytes();
        return (bytes[0] & 0xff) + "." + (bytes[1] & 0xff) + "." + (bytes[2] & 0xff) + "." + (bytes[3] & 0xff) + "." + (bytes[4] & 0xff) + "." + (bytes[5] & 0xff);
    }
}
