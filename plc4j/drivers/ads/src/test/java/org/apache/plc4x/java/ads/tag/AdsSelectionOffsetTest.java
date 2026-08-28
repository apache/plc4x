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
package org.apache.plc4x.java.ads.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * An ADS index offset is a <em>byte</em> offset, while a selection offset counts elements. They
 * coincide only for a one-byte type, so an unscaled offset reads from inside an earlier element.
 */
class AdsSelectionOffsetTest {

    @Test
    @DisplayName("a one-byte type advances one byte per element")
    void oneBytePerElement() {
        assertEquals(3, DirectAdsTag.of("0x4020/0[3]:BYTE").getIndexOffset());
    }

    @Test
    @DisplayName("a four-byte type advances four bytes per element")
    void fourBytesPerElement() {
        // The fourth DINT begins twelve bytes along, not three.
        assertEquals(12, DirectAdsTag.of("0x4020/0[3]:DINT").getIndexOffset());
    }

    @Test
    @DisplayName("an eight-byte type advances eight bytes per element")
    void eightBytesPerElement() {
        assertEquals(16, DirectAdsTag.of("0x4020/0[2]:LREAL").getIndexOffset());
    }

    @Test
    @DisplayName("a range starts where its first element starts")
    void aRangeStartsAtItsFirstElement() {
        assertEquals(8, DirectAdsTag.of("0x4020/0[2..5]:DINT").getIndexOffset());
        assertEquals(4, DirectAdsTag.of("0x4020/0[2..5]:DINT").getNumberOfElements());
    }

    @Test
    @DisplayName("a string advances by its declared length plus its terminator")
    void stringsAdvanceByTheirStorageSize() {
        // A STRING(10) occupies 11 bytes, so the third begins 22 bytes along.
        assertEquals(22, DirectAdsStringTag.of("0x4020/0[2]:STRING(10)").getIndexOffset());
        // A WSTRING(10) occupies 22, so the third begins 44 bytes along.
        assertEquals(44, DirectAdsStringTag.of("0x4020/0[2]:WSTRING(10)").getIndexOffset());
    }

    @Test
    @DisplayName("a selection on a type whose width is unknown here is refused, not guessed")
    void anUnknownWidthIsRefused() {
        // The device's data-type table is not available while parsing an address, so there is no
        // size to place the element with. Reading from the wrong offset would be worse than saying so.
        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> DirectAdsTag.of("0x4020/0[3]:MY_UDT"));
        assertTrue(thrown.getMessage().contains("MY_UDT"), thrown.getMessage());
    }

    @Test
    @DisplayName("without a selection an unknown type is still fine")
    void anUnknownTypeWithoutASelectionIsUnaffected() {
        assertEquals(0, DirectAdsTag.of("0x4020/0:MY_UDT").getIndexOffset());
    }

    @Test
    @DisplayName("a bare index is a scalar; a one-element range is a list of one")
    void aRangeIsAnArrayEvenWhenItSpansOneElement() {
        assertTrue(DirectAdsTag.of("0x4020/0[4]:DINT").getArrayInfo().isEmpty(),
            "a bare index selects one element, which is a scalar");
        assertEquals(1, DirectAdsTag.of("0x4020/0[4..4]:DINT").getArrayInfo().size(),
            "a range is an array even when it spans one element");
        assertTrue(DirectAdsTag.of("0x4020/0[4..4]:DINT").getArrayInfo().get(0).isRange());
    }

    @Test
    @DisplayName("a string tag renders an address that parses back")
    void theStringTagRoundTrips() {
        // It used to render "0x16416/0:STRING(10)[4]" - decimal digits behind a hex prefix, and
        // the suffix form the pattern no longer accepts. Neither survives a re-parse.
        DirectAdsStringTag tag = DirectAdsStringTag.of("0x4020/0[0..3]:STRING(10)");
        assertEquals("0x4020/0[0..3]:STRING(10)", tag.getAddressString());
        assertEquals(tag, DirectAdsStringTag.of(tag.getAddressString()));
    }

    @Test
    @DisplayName("a direct tag renders an address that parses back to the same tag")
    void theDirectTagRoundTrips() {
        // 16416 is 0x4020. Rendered as "0x16416" it re-parses as 91158 - the driver would read a
        // different memory location than the one the tag names.
        DirectAdsTag tag = DirectAdsTag.of("0x4020/0[0..3]:DINT");
        assertEquals("0x4020/0[0..3]:DINT", tag.getAddressString());
        assertEquals(tag, DirectAdsTag.of(tag.getAddressString()));
        assertEquals(0x4020, DirectAdsTag.of(tag.getAddressString()).getIndexGroup());
    }
}
