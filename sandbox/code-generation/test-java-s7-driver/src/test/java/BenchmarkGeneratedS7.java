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

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.s7.TPKTPacket;
import org.apache.plc4x.java.s7.io.TPKTPacketIO;
import org.apache.plc4x.java.utils.ReadBuffer;

public class BenchmarkGeneratedS7 {

    public static void main(String[] args) throws Exception {
        byte[] data = Hex.decodeHex("0300006702f080320100000001005600000407120a10060001032b84000160120a10020001032b840001a0120a10010001032b840001a9120a10050001032b84000150120a10020001032b84000198120a10040001032b84000140120a10020001032b84000190");
        long start = System.currentTimeMillis();
        int numRuns = 20000;
        for(int i = 0; i < numRuns; i++) {
            ReadBuffer buf = new ReadBuffer(data);
            TPKTPacket packet = TPKTPacketIO.parse(buf);
        }
        long end = System.currentTimeMillis();
        System.out.println("Parsed " + numRuns + " packets in " + (end - start) + "ms");
        System.out.println("That's " + ((float) (end - start) / numRuns) + "ms per packet");
    }

}
