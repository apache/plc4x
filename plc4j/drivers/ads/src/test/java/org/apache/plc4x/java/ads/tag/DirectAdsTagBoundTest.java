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
        assertThrows(PlcInvalidTagException.class, () -> DirectAdsTag.of("1/2:INT[99999999999]"));
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
            () -> DirectAdsStringTag.of("1/2:STRING(80)[99999999999]"));
    }

    @Test
    void aStringTagIndexGroupPastFourBytesIsAlsoRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> DirectAdsStringTag.of("4294967296/2:STRING(80)"));
    }

    @Test
    void aPlausibleTagStillParses() {
        DirectAdsTag tag = DirectAdsTag.of("1/2:INT[4]");
        assertEquals(4, tag.getNumberOfElements());
        assertEquals(1, tag.getIndexGroup());
        assertEquals(2, tag.getIndexOffset());
    }
}
