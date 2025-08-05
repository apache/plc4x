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

package org.apache.plc4x.java.iec608705104;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.iec608705104.readwrite.APDU;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;

public class ManualParserSerializerTest {

    public static void main(String[] args) throws Exception {
        byte[] data = Hex.decodeHex("0b00687010000000091101000200020000e60100010000d00b00020000e60100010000d00b00020000e60100010000d00b00020000e60100010000d00b00020000e60100010000d00b00020000e60100010000d00b00020000e60100010000d00b00020000e60100010000d00b00020000e6");
        ReadBuffer readBuffer = new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN);
        final APDU apdu = APDU.staticParse(readBuffer);
        System.out.println(apdu);
    }

}
