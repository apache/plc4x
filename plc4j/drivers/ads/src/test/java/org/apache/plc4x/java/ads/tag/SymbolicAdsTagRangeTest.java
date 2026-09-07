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
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Where a range may appear in a symbolic address.
 *
 * <p>A range has to be contiguous to be a single read. Only the last dimension of the trailing
 * selection qualifies: everything before it names one element of one structure, and a range
 * there would ask for several disjoint pieces at once - which no single request can fetch.
 */
class SymbolicAdsTagRangeTest {

    @Test
    void aRangeOnTheLastSegmentIsAccepted() {
        assertEquals(8, SymbolicAdsTag.of("MAIN.g_arr[1..8]").getSelection().get(0).getSize());
        assertEquals(4, SymbolicAdsTag.of("MAIN.g_arr[1].member[2..5]").getSelection().get(0).getSize());
    }

    @Test
    void bareIndicesOnInteriorSegmentsAreAccepted() {
        assertNotNull(SymbolicAdsTag.of("MAIN.g_arr[1].member[2]"));
        assertNotNull(SymbolicAdsTag.of("MAIN.g_arr[1][2].member"));
    }

    /** Member b of elements 1 to 3 is three separate reads, so the address is refused. */
    @Test
    void aRangeOnAnInteriorSegmentIsRefused() {
        assertThrows(PlcInvalidTagException.class, () -> SymbolicAdsTag.of("MAIN.g_arr[1..3].member"));
    }

    /** A strided slice of a multi-dimensional array is not contiguous either. */
    @Test
    void aRangeOnAnythingButTheLastDimensionIsRefused() {
        assertThrows(PlcInvalidTagException.class, () -> SymbolicAdsTag.of("MAIN.g_arr[1..3][2]"));
        assertNotNull(SymbolicAdsTag.of("MAIN.g_arr[1][2..5]"), "the last dimension may span");
    }

    @Test
    void anAddressWithNoSelectionStatesNone() {
        assertTrue(SymbolicAdsTag.of("MAIN.g_arr").getSelection().isEmpty());
    }
}
