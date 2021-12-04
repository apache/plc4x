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
package org.apache.plc4x.protocol.ads;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.ads.readwrite.AmsTCPPacket;
import org.apache.plc4x.java.ads.readwrite.io.AmsTCPPacketIO;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;

public class ManualParserSerializerTest {

    public static void main(String[] args) throws Exception {
        byte[] data = Hex.decodeHex("0000f1000000c0a817c80101feffc0a817140101530309000500d1000000000000002f00000000000000c90000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012a00a07d0f7e1c8e31489f8cf1006b3604b888defc88f64eb6fbb217d11fc5cf5f148b0abf05407c941a93b8263301db0f4940d66875727a0000000000d2040000ab6459032bbf7e03b888defcb8a5b249044de82ee03a2ab8a577006f006c0066000000488bf0488bfa33c9");
        //byte[] data = Hex.decodeHex("0000f1000000c0a817c80101feffc0a817140101530309000500d1000000000000002f00000000000000c90000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012a00a07d0f7e1c8e31489f8cf1006b3604b888defc");
        ReadBuffer readBuffer = new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN);
        final AmsTCPPacket amsTCPPacket = AmsTCPPacketIO.staticParse(readBuffer);
        System.out.println(amsTCPPacket);
    }

}
