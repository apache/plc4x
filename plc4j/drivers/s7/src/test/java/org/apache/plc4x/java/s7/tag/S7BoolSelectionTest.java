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
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * BOOLs are packed eight to a byte, so the start of a selection on a BOOL array counts bits from
 * the written bit offset. It used to count bytes, as for every other type, which put
 * {@code %DB1:16.0[4]:BOOL} at 20.0 instead of 16.4.
 */
class S7BoolSelectionTest {

    @Test
    void aStartIndexMovesTheBitOffset() {
        S7Tag tag = S7Tag.of("%DB1:16.0[4..7]:BOOL");

        assertEquals(16, tag.getByteOffset());
        assertEquals(4, tag.getBitOffset());
        assertEquals(4, tag.getNumberOfElements());
    }

    @Test
    void aStartIndexCarriesIntoTheNextByte() {
        S7Tag tag = S7Tag.of("%DB1.DBX16.6[3]:BOOL");

        assertEquals(17, tag.getByteOffset());
        assertEquals(1, tag.getBitOffset());
        assertEquals(1, tag.getNumberOfElements());
    }

    @Test
    void aDeclaredBaseIsHonoured() {
        // Elements 9 to 16 of an array declared from 1 start at offset 8, one byte in.
        S7Tag tag = S7Tag.of("%M10.0[9..16;1]:BOOL");

        assertEquals(11, tag.getByteOffset());
        assertEquals(0, tag.getBitOffset());
        assertEquals(8, tag.getNumberOfElements());
    }

    @Test
    void otherTypesStillStepByTheirSize() {
        assertEquals(28, S7Tag.of("%DB1:20[4..7]:WORD").getByteOffset());
    }

    @Test
    void aSelectionPastTheAddressableAreaIsRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> S7Tag.of("%DB1:16.0[" + Integer.MAX_VALUE + "]:BOOL"));
    }
}
