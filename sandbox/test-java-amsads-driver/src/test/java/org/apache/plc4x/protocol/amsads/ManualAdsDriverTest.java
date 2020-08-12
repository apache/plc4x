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
package org.apache.plc4x.protocol.amsads;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.amsads.readwrite.AmsTCPPacket;
import org.apache.plc4x.java.amsads.readwrite.io.AmsTCPPacketIO;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.spi.generation.ReadBuffer;

public class ManualAdsDriverTest {

    public static void main(String[] args) throws Exception {
        // Working request:
        // 000050000000c0a8171401015303c0a817c80101feff09000400300000000000000002000000 81f0000002000000080000002000000006f00000000000000400000006f0000000000000040000000300801b3702001b
        // Working response:
        // 000030000000c0a817cd0101feffc0a81714010153030900050010000000000000000100000000000000080000000000000000000000

        // PLC4X request:
        // 000050000000c0a8171401015303c0a817c80101feff09000400300000000000000002000000 80f0000002000000020000002000000005f000000000801d010000000000000005f000003702001d0100000000000000
        // PLC4X response:
        // 000028000000c0a817c80101feffc0a8171401015303090005000800000000000000020000000507000000000000



        byte[] data = Hex.decodeHex("000028000000c0a817c80101feffc0a8171401015303090005000800000000000000010000000507000000000000");
        ReadBuffer readBuffer = new ReadBuffer(data, true);
        final AmsTCPPacket amsTCPPacket = AmsTCPPacketIO.staticParse(readBuffer);
        System.out.println(amsTCPPacket);
    }

}
