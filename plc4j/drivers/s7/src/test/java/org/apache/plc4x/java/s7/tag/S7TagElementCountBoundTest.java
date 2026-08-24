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
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1:0:INT[99999999999]"));
    }

    @Test
    void aCountSpanningMoreThanTheAddressableAreaIsRejected() {
        // Two million LREALs are sixteen megabytes, and no S7 area is that large.
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1:0:LREAL[2000000]"));
    }

    @Test
    void theSameCountOfSingleBytesIsWithinTheAreaAndStillParses() {
        assertEquals(2000000, S7Tag.of("%DB1:0:BYTE[2000000]").getNumberOfElements());
    }

    @Test
    void aFixedLengthStringCountIsBoundedByWhatOneStringCosts() {
        // 9999 strings of 254 characters plus their two length bytes is past the area; the
        // optimizer would have multiplied that out in an int before anybody looked at it.
        assertThrows(PlcInvalidTagException.class,
            () -> S7StringFixedLengthTag.of("%DB1:0:STRING(254)[9999]"));
    }

    @Test
    void aFixedLengthStringCountThatFitsStillParses() {
        assertEquals(8000, S7StringFixedLengthTag.of("%DB1:0:STRING(254)[8000]").getNumberOfElements());
    }

    @Test
    void aVarLengthStringCountIsBoundedByTheLengthTheDriverAssumes() {
        assertThrows(PlcInvalidTagException.class,
            () -> S7StringVarLengthTag.of("%DB1:0:STRING[9999]"));
    }

    @Test
    void aWideStringCostsTwiceAsMuchPerElement() {
        // The same count that fits as STRING does not fit as WSTRING.
        assertEquals(8000, S7StringFixedLengthTag.of("%DB1:0:STRING(254)[8000]").getNumberOfElements());
        assertThrows(PlcInvalidTagException.class,
            () -> S7StringFixedLengthTag.of("%DB1:0:WSTRING(254)[8000]"));
    }

    @Test
    void aPlausibleCountStillParses() {
        assertEquals(64, S7Tag.of("%DB1:0:INT[64]").getNumberOfElements());
    }
}
