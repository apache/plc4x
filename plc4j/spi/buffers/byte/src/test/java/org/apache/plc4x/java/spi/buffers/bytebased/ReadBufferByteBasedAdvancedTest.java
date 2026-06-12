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

import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderBigEndian;
import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderLittleEndian;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingTwosComplement;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ReadBufferByteBasedAdvancedTest {

    /**
     * Uses a BIG endian buffer and does both a BIG and LITTLE read from it.
     */
    @Test
    void testMixEndiannessReadings() throws Exception {
        byte[] data = {
            (byte) 0x12, (byte) 0x34, // for short
            (byte) 0x56, (byte) 0x78, (byte) 0x9A, (byte) 0xBC // for int
        };
        ReadBufferByteBased buffer = new ReadBufferByteBased(data, EncodingTwosComplement.optionEncodingTwosComplement());

        // Read a short (big endian)
        int shortValue = buffer.readSignedShort(16, ByteOrderBigEndian.optionByteOrderBigEndian());
        assertEquals(0x1234, shortValue);

        // Read an int (little endian)
        int intValue = buffer.readSignedInt(32, ByteOrderLittleEndian.optionByteOrderLittleEndian());
        assertEquals(0xBC9A7856, intValue);
    }

    /**
     * Doing unaligned readings (reading a field starting not at the beginning of a full byte and crossing over
     * to the next.
     */
    @Test
    void testCrossingByteBoundary() throws Exception {
        byte[] data = {(byte) 0b1011_0110, (byte) 0b0011_1101};
        ReadBufferByteBased buffer = new ReadBufferByteBased(data, EncodingTwosComplement.optionEncodingTwosComplement());

        // Read 5 bits: 10110
        byte[] firstBytes = buffer.readBits(5);
        assertEquals(0b10110, firstBytes[0]);

        // Read the next 7 bits: 1100011
        byte[] secondBytes = buffer.readBits(7);
        assertEquals(0b1100011, secondBytes[0]);
    }

    /**
     * Creating a buffer with one endianness and a sub-buffer with another.
     */
    @Test
    void testSubBufferEndianness() throws Exception {
        byte[] data = {
            (byte) 0xAB, (byte) 0xCD, (byte) 0xEF, (byte) 0x01
        };
        ReadBufferByteBased buffer = new ReadBufferByteBased(data, EncodingTwosComplement.optionEncodingTwosComplement(), ByteOrderLittleEndian.optionByteOrderLittleEndian());

        ReadBufferByteBased subBuffer = buffer.createSubBuffer(16);

        int value = subBuffer.readSignedShort(16, ByteOrderBigEndian.optionByteOrderBigEndian());
        assertEquals(-21555, value);

        int remaining = buffer.readSignedShort(16); // using LITTLE (default)
        assertEquals(0x01EF, remaining);
    }

    /**
     * Generally unaligned reads.
     */
    @Test
    void testWeirdBitWidths() throws Exception {
        byte[] data = {(byte) 0b1111_0000, (byte) 0b1100_1100, (byte) 0b1010_1010};
        ReadBufferByteBased buffer = new ReadBufferByteBased(data, EncodingTwosComplement.optionEncodingTwosComplement());

        byte[] v1Bytes = buffer.readBits(4); // 1111
        assertEquals(0b1111, v1Bytes[0]);

        byte[] v2Bytes = buffer.readBits(5); // 00001
        assertEquals(0b00001, v2Bytes[0]);

        byte[] v3Bytes = buffer.readBits(7); // 1001100
        assertEquals(0b1001100, v3Bytes[0]);

        byte[] v4Bytes = buffer.readBits(8); // 10101010
        assertEquals(0b10101010, v4Bytes[0] & 0xFF);
    }

    /**
     * Creating sub buffers from sub buffers.
     */
    @Test
    void testDeepNestedSubBuffers() throws Exception {
        byte[] data = {(byte) 0xDE, (byte) 0xAD, (byte) 0xBE, (byte) 0xEF};
        ReadBufferByteBased buffer = new ReadBufferByteBased(data, EncodingTwosComplement.optionEncodingTwosComplement());

        ReadBufferByteBased sub1 = buffer.createSubBuffer(16);
        ReadBufferByteBased sub2 = sub1.createSubBuffer(8);

        assertEquals(-34, sub2.readSignedByte(8));

        assertEquals(-83, sub1.readSignedByte(8));
        assertEquals(-66, buffer.readSignedByte(8));
        assertEquals(-17, buffer.readSignedByte(8));
    }
}
