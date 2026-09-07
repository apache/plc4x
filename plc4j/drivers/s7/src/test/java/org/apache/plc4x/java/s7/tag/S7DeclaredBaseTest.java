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
package org.apache.plc4x.java.s7.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Addressing an array that does not start at zero.
 *
 * <p>A TIA declaration of {@code ARRAY[1..10] OF BYTE} at {@code %DB42:28.0} is addressed with
 * the indices the PLC program shows, and the declared lower bound tells the driver where the
 * data really begins. The written indices are kept as written; the offset is derived.
 */
class S7DeclaredBaseTest {

    /**
     * S7 addresses a memory offset, so the declared base is consumed while parsing: the byte
     * offset moves and the reported selection is zero-based. The base has no meaning once it has
     * been applied - see FR-014.
     */
    @Test
    void theBaseIsConsumedAndTheReportedSelectionIsZeroBased() {
        S7Tag tag = S7Tag.of("%DB42:28.0[4..7;1]:BYTE");

        ArrayInfo dimension = tag.getArrayInfo().get(0);
        assertEquals(0, dimension.getLowerBound(), "normalised once the base was applied");
        assertEquals(3, dimension.getUpperBound());
        assertEquals(4, dimension.getSize(), "still four elements");
        assertEquals(31, tag.getByteOffset(), "the base moved the offset instead");
    }

    /** The base moves where the read starts: element 4 of an array declared from 1 is the fourth. */
    @Test
    void theDeclaredBaseMovesTheByteOffset() {
        assertEquals(28, S7Tag.of("%DB42:28.0[1..4;1]:BYTE").getByteOffset(), "starts at the array");
        assertEquals(31, S7Tag.of("%DB42:28.0[4..7;1]:BYTE").getByteOffset(), "three BYTEs in");
        assertEquals(28, S7Tag.of("%DB42:28.0[1;1]:BYTE").getByteOffset(), "the first element");
    }

    /** Element size scales the offset - a WORD is two bytes, so the same index moves twice as far. */
    @Test
    void theOffsetScalesWithTheElementSize() {
        assertEquals(28 + 3 * 2, S7Tag.of("%DB42:28.0[4..7;1]:WORD").getByteOffset());
    }

    @Test
    void anIndexBelowTheDeclaredBaseIsRejected() {
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB42:28.0[0;1]:BYTE"));
    }

    /** Without a declared base the indices are zero-based, so nothing moves. */
    @Test
    void withoutADeclaredBaseTheIndicesAreZeroBased() {
        assertEquals(28, S7Tag.of("%DB42:28.0[0..3]:BYTE").getByteOffset());
        assertEquals(0, S7Tag.of("%DB42:28.0[0..3]:BYTE").getArrayInfo().get(0).getBase());
    }
}
