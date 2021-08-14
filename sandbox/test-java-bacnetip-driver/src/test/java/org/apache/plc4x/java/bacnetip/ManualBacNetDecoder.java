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
package org.apache.plc4x.java.bacnetip;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.bacnetip.readwrite.io.BVLCIO;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;

public class ManualBacNetDecoder {

    public static void main(String[] args) throws Exception {
        final byte[] bytes = Hex.decodeHex("810a002b01040205790109011c020000142c000002f93a06b24e09552e44434a00002f096f2e8204002f4f");
        ReadBuffer readBuffer = new ReadBufferByteBased(bytes);
        final BVLC packet = BVLCIO.staticParse(readBuffer);
        System.out.println(packet);
    }

}
