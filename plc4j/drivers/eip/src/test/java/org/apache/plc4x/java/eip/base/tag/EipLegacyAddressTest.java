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
package org.apache.plc4x.java.eip.base.tag;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Addresses written before the array notation was unified must not parse.
 *
 * <p>EIP's brackets already meant an index, so what changed here is the separate ':elementNb'
 * suffix that used to carry the count. An address still using it does not parse, rather than
 * quietly reading one element where it used to read several.
 */
class EipLegacyAddressTest {

    @Test
    void theElementCountSuffixNoLongerParses() {
        assertNull(EipTag.of("myArray[0]:DINT:8"));
        assertNull(EipTag.of("myTag:DINT:4"));
        assertNull(EipTag.of("myTag:4"));
    }

    @Test
    void theReplacementFormParses() {
        assertEquals(8, EipTag.of("myArray[0..7]:DINT").getElementNb());
        assertEquals(4, EipTag.of("myTag[0..3]:DINT").getElementNb());
    }

    /** matches() and of() agree, so a tag handler sees the same answer either way. */
    @Test
    void matchesAgreesWithOf() {
        assertFalse(EipTag.matches("myArray[0]:DINT:8"));
        assertTrue(EipTag.matches("myArray[0..7]:DINT"));
    }
}
