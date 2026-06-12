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
package org.apache.plc4x.java.utils.testutils.pcap;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FramingSpecTest {

    @Test
    void testRecordAccessors() {
        FramingSpec spec = new FramingSpec(2, 2, true, 6);

        assertEquals(2, spec.lengthFieldOffset());
        assertEquals(2, spec.lengthFieldSize());
        assertTrue(spec.bigEndian());
        assertEquals(6, spec.lengthAdjustment());
    }

    @Test
    void testNegativeOffsetThrows() {
        assertThrows(IllegalArgumentException.class, () ->
            new FramingSpec(-1, 2, true, 0));
    }

    @Test
    void testInvalidLengthFieldSizeThrows() {
        assertThrows(IllegalArgumentException.class, () ->
            new FramingSpec(0, 3, true, 0));
        assertThrows(IllegalArgumentException.class, () ->
            new FramingSpec(0, 0, true, 0));
        assertThrows(IllegalArgumentException.class, () ->
            new FramingSpec(0, 5, true, 0));
    }

    @Test
    void testValidLengthFieldSizes() {
        assertDoesNotThrow(() -> new FramingSpec(0, 1, true, 0));
        assertDoesNotThrow(() -> new FramingSpec(0, 2, true, 0));
        assertDoesNotThrow(() -> new FramingSpec(0, 4, true, 0));
    }

    @Test
    void testNegativeAdjustmentIsAllowed() {
        // Some protocols might use negative adjustment (total < length field value)
        assertDoesNotThrow(() -> new FramingSpec(0, 2, true, -2));
    }

    @Test
    void testZeroOffsetIsValid() {
        FramingSpec spec = new FramingSpec(0, 2, false, 10);
        assertEquals(0, spec.lengthFieldOffset());
        assertFalse(spec.bigEndian());
    }
}
