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

package org.apache.plc4x.protocol.df1;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.df1.readwrite.DF1Symbol;
import org.apache.plc4x.java.df1.readwrite.io.DF1SymbolIO;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.junit.jupiter.api.Test;

public class IOTest {

    @Test
    public void testXml() throws Exception {
        byte[] rData = Hex.decodeHex("10020900010001001100021003546F");
        ObjectMapper mapper = new XmlMapper().enableDefaultTyping();
        ReadBuffer rBuf = new ReadBufferByteBased(rData, false);
        DF1Symbol symbol = new DF1SymbolIO().parse(rBuf);
        String xml = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(symbol);
        System.out.println(xml);
        DF1Symbol symbol2 = mapper.readValue(xml, DF1Symbol.class);
        System.out.println(symbol2);
    }

    @Test
    public void testJson() throws Exception {
        byte[] rData = Hex.decodeHex("10020A0941000100FFFF1003DFB9");
        ObjectMapper mapper = new ObjectMapper().enableDefaultTyping();
        ReadBuffer rBuf = new ReadBufferByteBased(rData, false);
        DF1Symbol symbol = new DF1SymbolIO().parse(rBuf);
        String json = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(symbol);
        System.out.println(json);
        DF1Symbol symbol2 = mapper.readValue(json, DF1Symbol.class);
        System.out.println(symbol2);
    }

    /*@Test
    public void testParser() throws Exception {
        byte[] rData = Hex.decodeHex("0610020500180801c0a82a46c4090801c0a82a46c40a0203");
        long start = System.currentTimeMillis();
        int numRunsParse = 20000;

        KNXNetIPMessageIO knxNetIPMessageIO = new KNXNetIPMessageIO();

        // Benchmark the parsing code
        KNXNetIPMessage packet = null;
        for(int i = 0; i < numRunsParse; i++) {
            ReadBuffer rBuf = new ReadBuffer(rData);
            packet = knxNetIPMessageIO.parse(rBuf);
        }
        long endParsing = System.currentTimeMillis();

        System.out.println("Parsed " + numRunsParse + " packets in " + (endParsing - start) + "ms");
        System.out.println("That's " + ((float) (endParsing - start) / numRunsParse) + "ms per packet");

        // Benchmark the serializing code
        int numRunsSerialize = 20000;
        byte[] oData = null;
        for(int i = 0; i < numRunsSerialize; i++) {
            WriteBuffer wBuf = new WriteBuffer(packet.getLengthInBytes());
            knxNetIPMessageIO.serialize(wBuf, packet);
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
    }*/

}
