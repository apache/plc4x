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
package org.apache.plc4x.java.slmp;

import org.apache.plc4x.java.slmp.readwrite.SlmpReadRequest;
import org.apache.plc4x.java.slmp.readwrite.SlmpRequestFrame3E;
import org.apache.plc4x.java.slmp.readwrite.SlmpWriteRequest;
import org.apache.plc4x.java.slmp.tag.SlmpTag;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

class SlmpRequestBuildTest {

    private static String toHex(byte[] b) {
        StringBuilder sb = new StringBuilder();
        for (byte x : b) {
            sb.append(String.format("%02x", x));
        }
        return sb.toString();
    }

    @Test
    void buildsBatchReadFrameMatchingSh080008Example() throws Exception {
        // SH-080008 section 8.2 (Batch Read and Write): read D350, 2 words.
        SlmpTag tag = SlmpTag.of("D350[0..1]:WORD");
        SlmpReadRequest data = new SlmpReadRequest(
            tag.getDeviceNumber(), tag.getDeviceCode(), tag.getNumberOfPoints());
        SlmpRequestFrame3E frame = new SlmpRequestFrame3E(0x0000, 0x0401, 0x0000, data);

        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[frame.getLengthInBytes()]);
        frame.serialize(buffer);

        assertEquals("500000ffff03000c000000010400005e0100a80200", toHex(buffer.getBytes()));
    }

    @Test
    void buildsBatchWriteFrameForSingleWord() throws Exception {
        // Batch Write D350 = one WORD 0x1234 (little-endian payload 34 12).
        SlmpTag tag = SlmpTag.of("D350");
        SlmpWriteRequest data = new SlmpWriteRequest(
            tag.getDeviceNumber(), tag.getDeviceCode(), tag.getNumberOfPoints(),
            new byte[]{0x34, 0x12});
        SlmpRequestFrame3E frame = new SlmpRequestFrame3E(0x0000, 0x1401, 0x0000, data);

        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[frame.getLengthInBytes()]);
        frame.serialize(buffer);

        // requestDataLength = 6 + (3+1+2+2) = 14 (0x000E); command 0x1401 -> 01 14 LE.
        assertEquals("500000ffff03000e000000011400005e0100a801003412", toHex(buffer.getBytes()));
    }
}
