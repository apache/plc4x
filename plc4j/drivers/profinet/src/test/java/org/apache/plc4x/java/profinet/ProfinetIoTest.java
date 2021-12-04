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
package org.apache.plc4x.java.profinet;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.profinet.readwrite.DceRpc_Packet;
import org.apache.plc4x.java.profinet.readwrite.io.DceRpc_PacketIO;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;

public class ProfinetIoTest {

    public static void main(String[] args) throws Exception {
        long val1 = -559939584;
        long val2 = 3735027712L;
        System.out.println(Long.toHexString(val1));
        System.out.println(Long.toHexString(val2));
        //  -559939584 == ffffffffdea00000
        // 3735027712l == 00000000dea00000

        final byte[] rawResponse = Hex.decodeHex("04020a00100000000000a0de976cd111827100010904002a0100a0de976cd111827100a02442df7d29f2443745f2e246b8035e0ab26b8293a400000001000000000000000000ffffffff5a000000000000000000460000009401000000000000460000008101001e01000001654519352df3b6428f874371217c2b510002883f990006ef889281020008010000010001bbf08102000801000002000280008103000801000001000000c8");
        ReadBuffer readBuffer = new ReadBufferByteBased(rawResponse);
        final DceRpc_Packet dceRpc_packet = DceRpc_PacketIO.staticParse(readBuffer);
        System.out.println(dceRpc_packet);
    }


}
