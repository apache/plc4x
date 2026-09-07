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
package org.apache.plc4x.java.simulated.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.simulated.tag.SimulatedTag;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * The same selection means the same thing here as on every other driver (SC-001).
 *
 * <p>Each driver asserts this against its own address syntax; the assertions are deliberately
 * identical, because that sameness is the success criterion.
 */
class SimulatedArrayParityTest {

    @Test
    void aRangeOfEightStartingAtZero() {
        List<ArrayInfo> dimensions = SimulatedTag.of("RANDOM/foo[0..7]:INT").getArrayInfo();

        assertEquals(1, dimensions.size(), "one dimension");
        assertEquals(8, dimensions.get(0).getSize(), "eight elements");
        assertEquals(0, dimensions.get(0).getLowerBound() - dimensions.get(0).getBase(),
            "starting at offset zero");
        assertTrue(dimensions.get(0).isRange(), "written as a range");
    }

    /**
     * This driver names a variable and knows no element size, so it can only select from the
     * first element. A selection that starts elsewhere is refused rather than silently read from
     * the start - see FR-034.
     */
    @Test
    void aSelectionMustStartAtTheFirstElement() {
        assertThrows(PlcInvalidTagException.class, () -> SimulatedTag.of("RANDOM/foo[4]:INT"));
        assertThrows(PlcInvalidTagException.class, () -> SimulatedTag.of("RANDOM/foo[4..7]:INT"));

        assertTrue(SimulatedTag.of("RANDOM/foo[0]:INT").getArrayInfo().isEmpty(), "[0] is still a scalar");
    }

    /**
     * The case the notation exists to distinguish: a range spanning one element is an array of
     * one, while a bare index is a scalar. No element count can tell them apart, so a driver that
     * derives its shape from the count alone silently collapses them - which is how the SLMP tag
     * reported a one-element range as a scalar while plc4go's reported a list.
     */
    @Test
    void aOneElementRangeIsAnArrayOfOne() {
        List<ArrayInfo> dimensions = SimulatedTag.of("RANDOM/foo[0..0]:INT").getArrayInfo();

        assertEquals(1, dimensions.size(), "one dimension");
        assertEquals(1, dimensions.get(0).getSize(), "one element");
        assertTrue(dimensions.get(0).isRange(), "written as a range");

        assertTrue(SimulatedTag.of("RANDOM/foo[0]:INT").getArrayInfo().isEmpty(),
            "while the bare index of the same element stays a scalar");
    }

    @Test
    void anOmittedSelectionIsAScalar() {
        assertTrue(SimulatedTag.of("RANDOM/foo:INT").getArrayInfo().isEmpty());
    }
}
