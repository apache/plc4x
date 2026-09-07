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

import org.apache.plc4x.java.spi.values.PlcDATE;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcREAL;
import org.apache.plc4x.java.spi.values.PlcUDINT;
import org.apache.plc4x.java.spi.values.PlcUINT;
import org.apache.plc4x.java.spi.values.PlcWORD;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.Test;

import java.time.LocalDate;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class SlmpDataTypeEncodeTest {

    private static String toHex(byte[] b) {
        StringBuilder sb = new StringBuilder();
        for (byte x : b) {
            sb.append(String.format("%02x", x));
        }
        return sb.toString();
    }

    @Test
    void encodeSingleWordLittleEndian() {
        assertEquals("3412", toHex(SlmpDataType.WORD.encode(new PlcWORD(0x1234), 1)));
    }

    @Test
    void encodeWordHighBitUnsigned() {
        assertEquals("ffff", toHex(SlmpDataType.WORD.encode(new PlcWORD(65535), 1)));
    }

    @Test
    void encodeSignedIntTwosComplement() {
        assertEquals("ffff", toHex(SlmpDataType.INT.encode(new PlcINT((short) -1), 1)));
    }

    @Test
    void encodeDintLowWordFirst() {
        // 0x00010002 -> word0 0x0002 (02 00), word1 0x0001 (01 00)
        assertEquals("02000100", toHex(SlmpDataType.DINT.encode(new PlcDINT(0x00010002), 1)));
    }

    @Test
    void encodeRealIeee754() {
        // 1.0f = 0x3F800000, little-endian 00 00 80 3f
        assertEquals("0000803f", toHex(SlmpDataType.REAL.encode(new PlcREAL(1.0f), 1)));
    }

    @Test
    void encodeListWhenQuantityGreaterThanOne() {
        PlcValue list = new PlcList(List.of(new PlcWORD(0x56AB), new PlcWORD(0x170F)));
        assertEquals("ab560f17", toHex(SlmpDataType.WORD.encode(list, 2)));
    }

    @Test
    void encodeReturnsNullOnArityMismatch() {
        assertNull(SlmpDataType.WORD.encode(new PlcWORD(0x1234), 2));
        assertNull(SlmpDataType.WORD.encode(new PlcList(List.of(new PlcWORD(1))), 2));
    }

    @Test
    void encodeReturnsNullOnListInputWhenQuantityIsOne() {
        assertNull(SlmpDataType.WORD.encode(new PlcList(List.of(new PlcWORD(1))), 1));
    }

    @Test
    void encodeReturnsNullOnIncompatibleValue() {
        // PlcDATE.getInt() throws PlcIncompatibleDatatypeException (no numeric coercion);
        // the narrowed catch must map this to null rather than propagating.
        assertNull(SlmpDataType.WORD.encode(new PlcDATE(LocalDate.of(2020, 1, 1)), 1));
    }

    @Test
    void encodeDecodeRoundTrips() {
        byte[] bytes = SlmpDataType.UDINT.encode(new PlcUDINT(0xDEADBEEFL), 1);
        assertEquals(0xDEADBEEFL, SlmpDataType.UDINT.decode(bytes, 1).getLong());
    }

    @Test
    void encodeUintRoundTrips() {
        byte[] bytes = SlmpDataType.UINT.encode(new PlcUINT(40000), 1);
        assertEquals(40000, SlmpDataType.UINT.decode(bytes, 1).getInt());
    }
}
