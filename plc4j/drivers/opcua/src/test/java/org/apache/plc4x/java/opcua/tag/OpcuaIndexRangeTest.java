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
package org.apache.plc4x.java.opcua.tag;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * The OPC UA IndexRange the shared notation produces.
 *
 * <p>OPC-UA was the reference implementation the grammar was extracted from, so its wire output
 * must be exactly what it was before: 0-based, inclusive, one entry per dimension, comma
 * separated - and absent altogether when the address selects nothing, which is how OPC UA asks
 * for a whole node.
 */
class OpcuaIndexRangeTest {

    @Test
    void aSingleIndexBecomesOneEntry() {
        assertEquals("8", OpcuaTag.of("ns=2;i=MyInt[8];DINT").getIndexRange());
    }

    @Test
    void aRangeBecomesLowColonHigh() {
        assertEquals("3:8", OpcuaTag.of("ns=2;i=MyInt[3..8];DINT").getIndexRange());
    }

    /** The declared base is applied, so the wire form is always 0-based. */
    @Test
    void theDeclaredBaseIsResolvedAway() {
        assertEquals("2:7", OpcuaTag.of("ns=2;i=MyInt[3..8;1];DINT").getIndexRange());
    }

    @Test
    void dimensionsAreCommaJoined() {
        assertEquals("1:2,0:5", OpcuaTag.of("ns=2;i=MyInt[1..2][0..5];DINT").getIndexRange());
    }

    /** The comma spelling of the same selection produces the same IndexRange. */
    @Test
    void theCommaSpellingProducesTheSameIndexRange() {
        assertEquals(
            OpcuaTag.of("ns=2;i=MyInt[1..2][0..5];DINT").getIndexRange(),
            OpcuaTag.of("ns=2;i=MyInt[1..2,0..5];DINT").getIndexRange());
    }

    /**
     * No selection means no IndexRange, which is how OPC UA asks for the whole node - the value
     * of FR-022 for a driver that can determine the extent.
     */
    @Test
    void noSelectionMeansNoIndexRange() {
        assertNull(OpcuaTag.of("ns=2;i=MyInt;DINT").getIndexRange());
        assertTrue(OpcuaTag.of("ns=2;i=MyInt;DINT").getArrayInfo().isEmpty());
    }

    /** A bare index yields a scalar, so it reports no array info even though it has a range. */
    @Test
    void aBareIndexIsAScalarToTheCaller() {
        assertEquals("8", OpcuaTag.of("ns=2;i=MyInt[8];DINT").getIndexRange());
        assertTrue(OpcuaTag.of("ns=2;i=MyInt[8];DINT").getArrayInfo().isEmpty());
        assertFalse(OpcuaTag.of("ns=2;i=MyInt[8..8];DINT").getArrayInfo().isEmpty());
    }
}
