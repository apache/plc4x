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
package org.apache.plc4x.java.spi.drivers.model;

import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Rendering a selection and parsing it back yields the same selection (SC-006).
 *
 * <p>The spelling may normalise - a base of 0 is dropped, the comma form becomes one bracket per
 * dimension - but never the meaning, so a second pass is a fixed point.
 */
class ArrayNotationRoundTripTest {

    @ParameterizedTest
    @ValueSource(strings = {
        "[4]", "[0..7]", "[4;1]", "[4..7;1]", "[0]", "[7..7]",
        "[1..2][0..5]", "[4..7;1][7..10;2]", "[0][1][2]",
        "[0,1]", "[1..2,3..4]", "[0,1][2]", "[1..2;1,3..4;1]",
    })
    void aSelectionSurvivesBeingRenderedAndParsedAgain(String written) {
        var parsed = ArrayNotationParser.parse(written, "tag" + written);
        var reparsed = ArrayNotationParser.parse(
            ArrayNotationParser.render(parsed), "tag");

        assertEquals(parsed, reparsed, written);
    }

    /** Rendering is a fixed point: the canonical form renders to itself. */
    @ParameterizedTest
    @ValueSource(strings = {"[4]", "[0..7]", "[4..7;1]", "[1..2][0..5]", "[0,1]", "[1..2,3..4]"})
    void renderingTheCanonicalFormChangesNothingFurther(String written) {
        String once = ArrayNotationParser.render(ArrayNotationParser.parse(written, "tag"));
        String twice = ArrayNotationParser.render(ArrayNotationParser.parse(once, "tag"));

        assertEquals(once, twice, written);
    }

    /** The scalar/array distinction survives the trip, which is what FR-024 turns on. */
    @ParameterizedTest
    @ValueSource(strings = {"[4]", "[4..4]", "[0..7]", "[0,1]", "[0..0,1..1]"})
    void whetherEachDimensionWasARangeSurvives(String written) {
        var parsed = ArrayNotationParser.parse(written, "tag");
        var reparsed = ArrayNotationParser.parse(ArrayNotationParser.render(parsed), "tag");

        for (int i = 0; i < parsed.size(); i++) {
            assertEquals(parsed.get(i).isRange(), reparsed.get(i).isRange(),
                written + " dimension " + i);
        }
    }
}
