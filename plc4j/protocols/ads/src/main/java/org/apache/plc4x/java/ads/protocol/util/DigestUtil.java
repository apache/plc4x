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
package org.apache.plc4x.java.ads.protocol.util;

import com.github.snksoft.crc.CRC;
import org.apache.plc4x.java.ads.api.util.ByteReadable;

public class DigestUtil {

    public static final CRC.Parameters CRC16 = CRC.Parameters.CRC16;

    public static final CRC.Parameters CRC16_ADS = new CRC.Parameters(
        CRC16.getWidth(),
        CRC16.getPolynomial(),
        0xFFFF,
        CRC16.isReflectIn(),
        CRC16.isReflectOut(),
        CRC16.getFinalXor());

    private static CRC crc = new CRC(CRC16_ADS);

    private DigestUtil() {
        // Utility class
    }

    public static int calculateCrc16(ByteReadable... byteReadables) {
        if (byteReadables.length == 1) {
            return calculateCrc16(byteReadables[0].getBytes());
        }
        long currentCrcValue = crc.init();
        for (ByteReadable byteReadable : byteReadables) {
            currentCrcValue = crc.update(currentCrcValue, byteReadable.getBytes());
        }
        short finalCrc = crc.finalCRC16(currentCrcValue);
        return Short.toUnsignedInt(Short.reverseBytes(finalCrc));
    }

    public static int calculateCrc16(byte[] bytes) {
        short finalCrc = (short) crc.calculateCRC(bytes);
        return Short.toUnsignedInt(Short.reverseBytes(finalCrc));
    }

}
