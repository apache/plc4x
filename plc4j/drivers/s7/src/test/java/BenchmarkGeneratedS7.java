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

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;

import java.util.Arrays;

public class BenchmarkGeneratedS7 {

    public static void main(String[] args) throws Exception {
        byte[] rData = Hex.decodeHex("0300001611e00000000f00c2020100c1020311c0010a");
        long start = System.currentTimeMillis();
        int numRunsParse = 2000000;

        // Benchmark the parsing code
        TPKTPacket packet = null;
        for(int i = 0; i < numRunsParse; i++) {
            ReadBuffer rBuf = new ReadBufferByteBased(rData);
            packet = TPKTPacket.staticParse(rBuf);
        }
        long endParsing = System.currentTimeMillis();

        System.out.println("Parsed " + numRunsParse + " packets in " + (endParsing - start) + "ms");
        System.out.println("That's " + ((float) (endParsing - start) / numRunsParse) + "ms per packet");

        // Benchmark the serializing code
        int numRunsSerialize = 2000000;
        byte[] oData = null;
        for(int i = 0; i < numRunsSerialize; i++) {
            WriteBufferByteBased wBuf = new WriteBufferByteBased(packet.getLengthInBytes());
            packet.serialize(wBuf);
            oData = wBuf.getData();
        }
        long endSerializing = System.currentTimeMillis();

        System.out.println("Serialized " + numRunsSerialize + " packets in " + (endSerializing - endParsing) + "ms");
        System.out.println("That's " + ((float) (endSerializing - endParsing) / numRunsSerialize) + "ms per packet");
        if(!Arrays.equals(rData, oData)) {
            for(int i = 0; i < rData.length; i++) {
                if(rData[i] != oData[i]) {
                    System.out.println("Difference in byte " + i);
                }
            }
            System.out.println("Not equals");
        } else {
            System.out.println("Bytes equal");
        }
    }

}
