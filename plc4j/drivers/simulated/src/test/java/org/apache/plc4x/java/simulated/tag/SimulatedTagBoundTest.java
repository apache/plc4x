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
package org.apache.plc4x.java.simulated.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;



/**
 * The count in a simulated tag's address decides the size of the array the device fills with
 * made-up data. That size was worked out in an int, so a large enough count did not produce a
 * large array - it produced a negative one.
 */
public class SimulatedTagBoundTest {

    /** How many elements the tag reports, read the way the device reads it. */
    private static int elementsOf(String address) {
        SimulatedTag tag = SimulatedTag.of(address);
        return tag.getArrayInfo().isEmpty() ? 1 : tag.getArrayInfo().get(0).getSize();
    }

    @Test
    void aCountWhoseSizeWouldNotFitAnIntIsRefused() {
        // 400000000 LREALs is 3.2e9 bytes, which as an int is negative.
        assertThrows(PlcInvalidTagException.class, () -> SimulatedTag.of("RANDOM/foo:LREAL[400000000]"));
    }

    @Test
    void aCountThatWouldBeMerelyEnormousIsAlsoRefused() {
        assertThrows(PlcInvalidTagException.class, () -> SimulatedTag.of("RANDOM/foo:LREAL[300000000]"));
    }

    @Test
    void aCountTooWideToBeANumberIsAnInvalidTagNotANumberFormatError() {
        assertThrows(PlcInvalidTagException.class, () -> SimulatedTag.of("RANDOM/foo:LREAL[99999999999]"));
    }

    @Test
    void aWiderElementLeavesRoomForFewerOfThem() {
        // The budget is in bytes, so the same count passes as a byte and fails as an eight-byte
        // double: 4194304 of them is 32MiB, past the budget, while as SINT it is 4MiB.
        assertEquals(4194304, elementsOf("RANDOM/foo:SINT[4194304]"));
        assertThrows(PlcInvalidTagException.class, () -> SimulatedTag.of("RANDOM/foo:LREAL[4194304]"));
    }

    @Test
    void aPlausibleCountStillParses() {
        assertEquals(16, elementsOf("RANDOM/foo:LREAL[16]"));
    }

    @Test
    void aTagWithNoCountIsStillOneElement() {
        assertEquals(1, elementsOf("RANDOM/foo:LREAL"));
    }
}
