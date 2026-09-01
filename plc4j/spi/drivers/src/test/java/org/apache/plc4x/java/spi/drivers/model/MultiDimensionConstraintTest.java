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

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Multi-dimensional selections, and what happens on a protocol that cannot express them.
 *
 * <p>A driver addressing linear memory has no second dimension to carry, so it says so rather
 * than quietly reading the first one.
 */
class MultiDimensionConstraintTest {

    @Test
    void dimensionsComeBackInTheOrderTheyWereWritten() {
        List<ArrayInfo> dimensions = ArrayNotationParser.parse("[1..2][0..5]", "myTag[1..2][0..5]");

        assertEquals(2, dimensions.size());
        assertEquals(1, dimensions.get(0).getLowerBound());
        assertEquals(2, dimensions.get(0).getUpperBound());
        assertEquals(0, dimensions.get(1).getLowerBound());
        assertEquals(5, dimensions.get(1).getUpperBound());
    }

    @Test
    void theCommaSpellingYieldsTheSameOrder() {
        assertEquals(
            ArrayNotationParser.parse("[1..2][0..5]", "t"),
            ArrayNotationParser.parse("[1..2,0..5]", "t"));
    }

    /** The limit is named, so the user learns it is the protocol and not the syntax. */
    @Test
    void aSecondDimensionIsRefusedWhereTheProtocolCarriesOne() {
        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse("[1..2][0..5]", "myTag[1..2][0..5]",
                AddressConstraints.SINGLE_DIMENSION));

        assertTrue(thrown.getMessage().contains("2 dimensions"), thrown::getMessage);
        assertTrue(thrown.getMessage().contains("at most 1"), thrown::getMessage);
    }

    @Test
    void theCommaSpellingIsRefusedTheSameWay() {
        assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse("[1..2,0..5]", "myTag[1..2,0..5]",
                AddressConstraints.SINGLE_DIMENSION));
    }

    /** One dimension is fine on such a driver, however it is spelled. */
    @Test
    void oneDimensionIsStillAccepted() {
        assertDoesNotThrow(() -> ArrayNotationParser.parse("[0..5]", "myTag[0..5]",
            AddressConstraints.SINGLE_DIMENSION));
    }

    /** An omitted selection is no dimensions at all, which every protocol can carry. */
    @Test
    void anOmittedSelectionIsNotADimension() {
        assertTrue(ArrayNotationParser.parse("", "myTag", AddressConstraints.SINGLE_DIMENSION).isEmpty());
    }
}
