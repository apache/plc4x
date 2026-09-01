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

import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * The same selection means the same thing here as on every other driver (SC-001).
 *
 * <p>Each driver asserts this against its own address syntax; the assertions are deliberately
 * identical, because that sameness is the success criterion.
 */
class S7ArrayParityTest {

    @Test
    void aRangeOfEightStartingAtZero() {
        List<ArrayInfo> dimensions = S7Tag.of("%DB42:28.0[0..7]:BYTE").getArrayInfo();

        assertEquals(1, dimensions.size(), "one dimension");
        assertEquals(8, dimensions.get(0).getSize(), "eight elements");
        assertEquals(0, dimensions.get(0).getLowerBound() - dimensions.get(0).getBase(),
            "starting at offset zero");
        assertTrue(dimensions.get(0).isRange(), "written as a range");
    }

    @Test
    void aBareIndexIsAScalar() {
        assertTrue(S7Tag.of("%DB42:28.0[4]:BYTE").getArrayInfo().isEmpty());
    }

    /**
     * The case the notation exists to distinguish: a range spanning one element is an array of
     * one, while a bare index is a scalar. No element count can tell them apart, so a driver that
     * derives its shape from the count alone silently collapses them - which is how the SLMP tag
     * reported a one-element range as a scalar while plc4go's reported a list.
     */
    @Test
    void aOneElementRangeIsAnArrayOfOne() {
        List<ArrayInfo> dimensions = S7Tag.of("%DB42:28.0[4..4]:BYTE").getArrayInfo();

        assertEquals(1, dimensions.size(), "one dimension");
        assertEquals(1, dimensions.get(0).getSize(), "one element");
        assertTrue(dimensions.get(0).isRange(), "written as a range");

        assertTrue(S7Tag.of("%DB42:28.0[4]:BYTE").getArrayInfo().isEmpty(),
            "while the bare index of the same element stays a scalar");
    }

    @Test
    void anOmittedSelectionIsAScalar() {
        assertTrue(S7Tag.of("%DB42:28.0:BYTE").getArrayInfo().isEmpty());
    }
}
