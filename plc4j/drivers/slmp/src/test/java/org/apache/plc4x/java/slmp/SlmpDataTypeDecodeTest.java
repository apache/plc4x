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

import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class SlmpDataTypeDecodeTest {

    // Wire bytes are little-endian, words LSB-first. D350=0x56AB -> bytes ab 56.
    private static byte[] hex(String s) {
        byte[] out = new byte[s.length() / 2];
        for (int i = 0; i < out.length; i++) {
            out[i] = (byte) Integer.parseInt(s.substring(i * 2, i * 2 + 2), 16);
        }
        return out;
    }

    @Test
    void wordsPerElement() {
        assertEquals(1, SlmpDataType.WORD.getWordsPerElement());
        assertEquals(1, SlmpDataType.INT.getWordsPerElement());
        assertEquals(1, SlmpDataType.UINT.getWordsPerElement());
        assertEquals(2, SlmpDataType.DINT.getWordsPerElement());
        assertEquals(2, SlmpDataType.UDINT.getWordsPerElement());
        assertEquals(2, SlmpDataType.REAL.getWordsPerElement());
    }

    @Test
    void decodeSingleWordUnsigned() {
        PlcValue v = SlmpDataType.WORD.decode(hex("ab56"), 1);
        assertEquals(0x56AB, v.getInt());
    }

    @Test
    void decodeWordHighBitStaysUnsigned() {
        // 0xFFFF little-endian must decode to 65535 (unsigned), NOT -1.
        PlcValue v = SlmpDataType.WORD.decode(hex("ffff"), 1);
        assertEquals(65535, v.getInt());
    }

    @Test
    void decodeSingleIntSigned() {
        // 0xFFFF little-endian -> -1 as signed 16
        PlcValue v = SlmpDataType.INT.decode(hex("ffff"), 1);
        assertEquals(-1, v.getInt());
    }

    @Test
    void decodeDintTwoWordsLowWordFirst() {
        // value 0x00010002 stored low-word-first: word0=0x0002 (bytes 02 00), word1=0x0001 (bytes 01 00)
        PlcValue v = SlmpDataType.DINT.decode(hex("02000100"), 1);
        assertEquals(0x00010002, v.getInt());
    }

    @Test
    void decodeRealTwoWords() {
        // 1.0f = 0x3F800000; little-endian bytes 00 00 80 3f
        PlcValue v = SlmpDataType.REAL.decode(hex("0000803f"), 1);
        assertEquals(1.0f, v.getFloat(), 0.0f);
    }

    @Test
    void decodeListWhenQuantityGreaterThanOne() {
        // D350=0x56AB (ab 56), D351=0x170F (0f 17)
        PlcValue v = SlmpDataType.WORD.decode(hex("ab560f17"), 2);
        assertTrue(v.isList());
        assertEquals(0x56AB, v.getList().get(0).getInt());
        assertEquals(0x170F, v.getList().get(1).getInt());
    }

    @Test
    void decodeReturnsNullOnShortResponse() {
        // asking for 2 words but only 1 word of data
        assertNull(SlmpDataType.WORD.decode(hex("ab56"), 2));
    }
}
