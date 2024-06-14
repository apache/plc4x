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

package org.apache.plc4x.java.modbus.utils;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ByteOrderUtilsTest {

    static final int INT_INPUT =   0x01020304;
    static final long LONG_INPUT = 0x0102030405060708L;

    @Test
    void testIntToBigEndian() {
        assertEquals(0x01020304, ByteOrderUtils.toBigEndian(INT_INPUT));
        assertEquals(INT_INPUT, ByteOrderUtils.toBigEndian(ByteOrderUtils.toBigEndian(INT_INPUT)));
    }

    @Test
    void testIntToLittleEndian() {
        assertEquals(0x04030201, ByteOrderUtils.toLittleEndian(INT_INPUT));
        assertEquals(INT_INPUT, ByteOrderUtils.toLittleEndian(ByteOrderUtils.toLittleEndian(INT_INPUT)));
    }

    @Test
    void testIntToBigEndianWordSwap() {
        assertEquals(0x03040102, ByteOrderUtils.toBigEndianWordSwap(INT_INPUT));
        assertEquals(INT_INPUT, ByteOrderUtils.toBigEndianWordSwap(ByteOrderUtils.toBigEndianWordSwap(INT_INPUT)));
    }

    @Test
    void testIntToLittleEndianWordSwap() {
        assertEquals(0x02010403, ByteOrderUtils.toLittleEndianWordSwap(INT_INPUT));
        assertEquals(INT_INPUT, ByteOrderUtils.toLittleEndianWordSwap(ByteOrderUtils.toLittleEndianWordSwap(INT_INPUT)));
    }

    @Test
    void testLongToBigEndian() {
        assertEquals(0x0102030405060708L, ByteOrderUtils.toBigEndian(LONG_INPUT));
        assertEquals(LONG_INPUT, ByteOrderUtils.toBigEndian(ByteOrderUtils.toBigEndian(LONG_INPUT)));
    }

    @Test
    void testLongToLittleEndian() {
        assertEquals(0x0807060504030201L, ByteOrderUtils.toLittleEndian(LONG_INPUT));
        assertEquals(LONG_INPUT, ByteOrderUtils.toLittleEndian(ByteOrderUtils.toLittleEndian(LONG_INPUT)));
    }

    @Test
    void testLongToBigEndianWordSwap() {
        assertEquals(0x0201040306050807L, ByteOrderUtils.toBigEndianWordSwap(LONG_INPUT));
        assertEquals(LONG_INPUT, ByteOrderUtils.toBigEndianWordSwap(ByteOrderUtils.toBigEndianWordSwap(LONG_INPUT)));
    }

    @Test
    void testLongToLittleEndianWordSwap() {
        assertEquals(0x0708050603040102L, ByteOrderUtils.toLittleEndianWordSwap(LONG_INPUT));
        assertEquals(LONG_INPUT, ByteOrderUtils.toLittleEndianWordSwap(ByteOrderUtils.toLittleEndianWordSwap(LONG_INPUT)));
    }

}