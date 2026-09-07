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
package org.apache.plc4x.java.umas.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * What a UMAS symbolic address states about array selection.
 *
 * <p>The driver cannot execute a selection yet - its read reference is built from a symbol's own
 * block and offset, with no per-element arithmetic - so the connection refuses one rather than
 * returning the whole variable (FR-033). The address itself still parses, which is what these
 * assertions cover.
 */
class SymbolicUmasTagSelectionTest {

    @Test
    void aBareAddressStatesNoSelection() {
        assertTrue(SymbolicUmasTag.of("MyVar").getSelection().isEmpty());
        assertTrue(SymbolicUmasTag.of("MyVar").getArrayInfo().isEmpty());
    }

    /** A bare index is one element, so it reports no array info to the caller (FR-023). */
    @Test
    void aBareIndexIsAScalarToTheCaller() {
        assertEquals(1, SymbolicUmasTag.of("MyVar[1]").getSelection().size());
        assertTrue(SymbolicUmasTag.of("MyVar[1]").getArrayInfo().isEmpty());
    }

    @Test
    void aRangeStatesItsDimensions() {
        assertEquals(8, SymbolicUmasTag.of("MyVar[1..8]").getSelection().get(0).getSize());
        assertFalse(SymbolicUmasTag.of("MyVar[1..8]").getArrayInfo().isEmpty());
    }

    /** Only the last dimension may span, as everywhere else (FR-030). */
    @Test
    void aRangeBeforeTheLastDimensionIsRefused() {
        assertThrows(PlcInvalidTagException.class, () -> SymbolicUmasTag.of("MyVar[1..3].member"));
        assertThrows(PlcInvalidTagException.class, () -> SymbolicUmasTag.of("MyVar[1..3][2]"));
        assertNotNull(SymbolicUmasTag.of("MyVar[1][2..5]"));
    }

    /** Interior indices stay part of the symbolic path. */
    @Test
    void interiorIndicesAreAcceptedAsPath() {
        assertNotNull(SymbolicUmasTag.of("MyVar[1].member[2]"));
    }

    @Test
    void theDeclaredBaseIsCarriedForVerification() {
        assertEquals(1, SymbolicUmasTag.of("MyVar[1..8;1]").getDeclaredBase());
        assertNull(SymbolicUmasTag.of("MyVar[1..8]").getDeclaredBase());
    }
}
