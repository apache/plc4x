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
package org.apache.plc4x.java.spi.buffers.bytebased;

import java.util.Random;

class WriteBufferByteBasedFuzzTest {

    private static final Random RANDOM = new Random(5678);

    /*@RepeatedTest(50)
    void randomBitFuzzTest() {
        byte[] expected = new byte[64];
        RANDOM.nextBytes(expected);

        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[64]);

        int bitsWritten = 0;
        while (bitsWritten < expected.length * 8) {
            int bitsToWrite = Math.min(RANDOM.nextInt(32) + 1, expected.length * 8 - bitsWritten);
            int value = extractBits(expected, bitsWritten, bitsToWrite);
            buffer.writeBits(value, bitsToWrite);
            bitsWritten += bitsToWrite;
        }

        byte[] actual = buffer.getBytes();
        assertArrayEquals(expected, actual);
    }*/

    private int extractBits(byte[] data, int bitOffset, int numBits) {
        int value = 0;
        for (int i = 0; i < numBits; i++) {
            int absoluteBit = bitOffset + i;
            int byteIndex = absoluteBit / 8;
            int bitIndex = absoluteBit % 8;
            boolean bit = (data[byteIndex] & (0x80 >> bitIndex)) != 0;
            value = (value << 1) | (bit ? 1 : 0);
        }
        return value;
    }

}
