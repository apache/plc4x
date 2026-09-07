/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
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
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * ADS carries an index group, an index offset and a length as four bytes each. A tag naming
 * something wider addresses whatever survives the narrowing, so it is refused - and refused as
 * an invalid tag, not as an unchecked number format error escaping the parse.
 */
public class DirectAdsTagBoundTest {

    @Test
    void aCountTooWideToBeANumberIsAnInvalidTagNotANumberFormatError() {
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsTag.of("1/2[0..99999999998]:INT"));
    }

    @Test
    void anIndexGroupTooWideToBeANumberIsAnInvalidTag() {
        assertThrows(PlcInvalidTagException.class,
            () -> DirectAdsTag.of("99999999999999999999999/2:INT"));
    }

    @Test
    void aHexIndexGroupTooWideToBeANumberIsAnInvalidTag() {
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsTag.of("0xFFFFFFFFFFFFFFFFFF/2:INT"));
    }

    @Test
    void anIndexGroupPastFourBytesIsRejected() {
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsTag.of("4294967296/2:INT"));
    }

    @Test
    void theLastIndexGroupFourBytesCanHoldIsStillAccepted() {
        assertEquals(0xFFFFFFFFL, DirectAdsTag.of("4294967295/2:INT").getIndexGroup());
    }

    @Test
    void aStringTagCountTooWideToBeANumberIsAlsoAnInvalidTag() {
        assertThrows(PlcInvalidTagException.class,
            () -> DirectAdsStringTag.of("1/2[0..99999999998]:STRING(80)"));
    }

    @Test
    void aStringTagIndexGroupPastFourBytesIsAlsoRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> DirectAdsStringTag.of("4294967296/2:STRING(80)"));
    }

    @Test
    /**
     * A count of zero used to be rejected. The notation cannot express it: [0] names the element
     * at offset 0, which is one element, and an empty bracket names nothing at all.
     */
    void anEmptySelectionIsRefused() {
        assertEquals(1, DirectAdsTag.of("1/2[0]:INT").getNumberOfElements());
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsTag.of("1/2[]:INT"));
    }

    @Test
    void aNegativeIndexGroupHandedInDirectlyIsAlsoRefused() {
        // The address parser cannot produce this, but the factory taking numbers can.
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsTag.of(-1L, 2L, "INT", 1));
    }

    @Test
    void anIndexOffsetPastFourBytesHandedInDirectlyIsAlsoRefused() {
        assertThrows(PlcInvalidTagException.class,
            () -> DirectAdsTag.of(1L, 0x100000000L, "INT", 1));
    }

    @Test
    void hexAddressesAreReadAsHex() {
        DirectAdsTag tag = DirectAdsTag.of("0x10/0xFF[0..1]:INT");
        assertEquals(0x10, tag.getIndexGroup());
        assertEquals(0xFF, tag.getIndexOffset());
        assertEquals(2, tag.getNumberOfElements());
    }

    @Test
    void aStringTagSharesTheSameChecks() {
        DirectAdsStringTag tag = DirectAdsStringTag.of("0x10/2[0..2]:STRING(80)");
        assertEquals(0x10, tag.getIndexGroup());
        assertEquals(2, tag.getIndexOffset());
        assertEquals(3, tag.getNumberOfElements());
        assertEquals(80, tag.getStringLength());
    }

    @Test
    void aStringTagEmptySelectionIsAlsoRefused() {
        assertEquals(1, DirectAdsStringTag.of("1/2[0]:STRING(80)").getNumberOfElements());
        assertThrows(PlcInvalidTagException.class,
            () -> DirectAdsStringTag.of("1/2[]:STRING(80)"));
    }

    @Test
    void aCountThatWouldNotFitAnIntIsRefused() {
        // Ten digits match the pattern but do not fit the int the count is kept in.
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsTag.of("1/2[0..2999999999]:INT"));
    }

    @Test
    void aPlausibleTagStillParses() {
        DirectAdsTag tag = DirectAdsTag.of("1/2[0..3]:INT");
        assertEquals(4, tag.getNumberOfElements());
        assertEquals(1, tag.getIndexGroup());
        assertEquals(2, tag.getIndexOffset());
    }
}
