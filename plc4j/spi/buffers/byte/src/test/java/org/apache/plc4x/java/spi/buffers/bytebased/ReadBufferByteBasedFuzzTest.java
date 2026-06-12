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

import org.junit.jupiter.api.RepeatedTest;

import java.util.Random;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ReadBufferByteBasedFuzzTest {

    private static final Random RANDOM = new Random(1234);

    @RepeatedTest(50)
    void randomBitFuzzTest() throws Exception {
        byte[] randomData = new byte[64];
        RANDOM.nextBytes(randomData);

        ReadBufferByteBased buffer = new ReadBufferByteBased(randomData);

        int totalBits = randomData.length * 8;
        int bitsRead = 0;

        while (bitsRead < totalBits) {
            int bitsToRead = Math.min(RANDOM.nextInt(32) + 1, totalBits - bitsRead);

            int expected = extractBits(randomData, bitsRead, bitsToRead);
            byte[] actualBytes = buffer.readBits(bitsToRead);
            int actual = byteArrayToInt(actualBytes);

            assertEquals(expected, actual, "Mismatch at bitsRead=" + bitsRead + ", bitsToRead=" + bitsToRead);

            bitsRead += bitsToRead;
        }
    }

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

    private int byteArrayToInt(byte[] bytes) {
        if (bytes.length == 0) return 0;

        int result = 0;
        for (int i = 0; i < Math.min(4, bytes.length); i++) {
            result = (result << 8) | (bytes[i] & 0xFF);
        }
        return result;
    }
}
