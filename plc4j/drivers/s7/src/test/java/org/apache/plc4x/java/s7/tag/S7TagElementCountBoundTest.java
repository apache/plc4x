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
package org.apache.plc4x.java.s7.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * An element count in a tag address is a request to allocate, so it has to be answered before
 * anything is allocated for it, and it has to be answered as an invalid tag rather than as
 * whatever the number parser happened to raise.
 */
public class S7TagElementCountBoundTest {

    @Test
    void aCountTooWideToBeANumberIsAnInvalidTagNotANumberFormatError() {
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1:0[0..99999999998]:INT"));
    }

    @Test
    void aCountSpanningMoreThanTheAddressableAreaIsRejected() {
        // Two million LREALs are sixteen megabytes, and no S7 area is that large.
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1:0[0..1999999]:LREAL"));
    }

    @Test
    void theSameCountOfSingleBytesIsWithinTheAreaAndStillParses() {
        assertEquals(2000000, S7Tag.of("%DB1:0[0..1999999]:BYTE").getNumberOfElements());
    }

    /**
     * A start index is scaled by what one element costs, and the parser accepts indices up to
     * Integer.MAX_VALUE. Scaled in int, a large one wrapped into a small offset that looked
     * perfectly valid, and the tag then read a different location without complaint.
     */
    @Test
    void aSelectionOffsetThatOverflowsIsRejectedRatherThanWrapped() {
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1:0[268435456]:LREAL"));
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1.DBW0[536870912]:LREAL"));
    }

    /**
     * The same scaling without an overflow still has to land inside the addressable area - it was
     * not checked at all after the selection was applied.
     */
    @Test
    void aSelectionOffsetPastTheAddressableAreaIsRejected() {
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1:0[1000000]:LREAL"));
    }

    @Test
    void aSelectionOffsetInsideTheAddressableAreaStillParses() {
        assertEquals(800, S7Tag.of("%DB1:0[100]:LREAL").getByteOffset());
    }

    @Test
    void aFixedLengthStringCountIsBoundedByWhatOneStringCosts() {
        // 9999 strings of 254 characters plus their two length bytes is past the area; the
        // optimizer would have multiplied that out in an int before anybody looked at it.
        assertThrows(PlcInvalidTagException.class,
            () -> S7StringFixedLengthTag.of("%DB1:0[0..9998]:STRING(254)"));
    }

    @Test
    void aFixedLengthStringCountThatFitsStillParses() {
        assertEquals(8000, S7StringFixedLengthTag.of("%DB1:0[0..7999]:STRING(254)").getNumberOfElements());
    }

    @Test
    void aVarLengthStringCountIsBoundedByTheLengthTheDriverAssumes() {
        assertThrows(PlcInvalidTagException.class,
            () -> S7StringVarLengthTag.of("%DB1:0[0..9998]:STRING"));
    }

    @Test
    void aWideStringCostsTwiceAsMuchPerElement() {
        // The same count that fits as STRING does not fit as WSTRING.
        assertEquals(8000, S7StringFixedLengthTag.of("%DB1:0[0..7999]:STRING(254)").getNumberOfElements());
        assertThrows(PlcInvalidTagException.class,
            () -> S7StringFixedLengthTag.of("%DB1:0[0..7999]:WSTRING(254)"));
    }

    @Test
    void aPlausibleCountStillParses() {
        assertEquals(64, S7Tag.of("%DB1:0[0..63]:INT").getNumberOfElements());
    }
}
