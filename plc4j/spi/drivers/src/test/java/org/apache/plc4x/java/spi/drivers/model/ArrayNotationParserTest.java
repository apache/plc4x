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
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.junit.jupiter.params.provider.ValueSource;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * The one definition of the array notation, per
 * {@code specs/002-unified-array-notation/contracts/array-notation-grammar.md}.
 *
 * <p>Every row of that contract's semantics table and every row of its rejection table has a case
 * here. The grammar is shared by every driver, so a gap here is a gap everywhere.
 */
class ArrayNotationParserTest {

    // --- the semantics table ---

    @ParameterizedTest(name = "{0} -> {1} element(s) at offset {2}..{3}")
    @CsvSource({
        // expression, size, firstOffset, lastOffset
        "'[4]',        1, 4, 4",
        "'[0..7]',     8, 0, 7",
        "'[4;1]',      1, 3, 3",
        "'[4..7;1]',   4, 3, 6",
        "'[0]',        1, 0, 0",
        "'[7..7]',     1, 7, 7",
    })
    void singleDimensionResolvesToTheDocumentedOffsets(String expression, int size, int first, int last) {
        List<ArrayInfo> dimensions = ArrayNotationParser.parse(expression, "tag" + expression);

        assertEquals(1, dimensions.size());
        ArrayInfo only = dimensions.get(0);
        assertEquals(size, only.getSize(), "size");
        assertEquals(first, only.getLowerBound() - only.getBase(), "first offset");
        assertEquals(last, only.getUpperBound() - only.getBase(), "last offset");
    }

    /**
     * A bare index and a one-element range cover the same element but are not the same selection:
     * the first yields a scalar and the second an array of one. Equal bounds cannot tell them
     * apart, so the written form is carried on the dimension.
     */
    @Test
    void aBareIndexIsNotTheSameAsAOneElementRange() {
        List<ArrayInfo> index = ArrayNotationParser.parse("[4]", "tag[4]");
        List<ArrayInfo> range = ArrayNotationParser.parse("[4..4]", "tag[4..4]");

        assertNotEquals(range, index);
        assertFalse(index.get(0).isRange());
        assertTrue(range.get(0).isRange());
        assertEquals(index.get(0).getLowerBound(), range.get(0).getLowerBound());
        assertEquals(1, index.get(0).getSize());
        assertEquals(1, range.get(0).getSize());
    }

    @Test
    void writtenBoundsArePreservedNotResolved() {
        ArrayInfo dimension = ArrayNotationParser.parse("[4..7;1]", "tag[4..7;1]").get(0);

        assertEquals(4, dimension.getLowerBound(), "lower bound is as written");
        assertEquals(7, dimension.getUpperBound(), "upper bound is as written");
        assertEquals(1, dimension.getBase(), "declared base");
        assertEquals(4, dimension.getSize());
    }

    @Test
    void baseDefaultsToZero() {
        assertEquals(0, ArrayNotationParser.parse("[4]", "tag[4]").get(0).getBase());
    }

    @Test
    void multipleDimensionsKeepTheirWrittenOrder() {
        List<ArrayInfo> dimensions = ArrayNotationParser.parse("[1..2][0..5]", "tag[1..2][0..5]");

        assertEquals(2, dimensions.size());
        assertEquals(1, dimensions.get(0).getLowerBound());
        assertEquals(2, dimensions.get(0).getUpperBound());
        assertEquals(0, dimensions.get(1).getLowerBound());
        assertEquals(5, dimensions.get(1).getUpperBound());
    }

    @Test
    void eachDimensionCarriesItsOwnBase() {
        List<ArrayInfo> dimensions =
            ArrayNotationParser.parse("[4..7;1][7..10;2]", "tag[4..7;1][7..10;2]");

        assertEquals(2, dimensions.size());
        assertEquals(3, dimensions.get(0).getLowerBound() - dimensions.get(0).getBase());
        assertEquals(6, dimensions.get(0).getUpperBound() - dimensions.get(0).getBase());
        assertEquals(5, dimensions.get(1).getLowerBound() - dimensions.get(1).getBase());
        assertEquals(8, dimensions.get(1).getUpperBound() - dimensions.get(1).getBase());
    }

    // --- the rejection table ---

    @ParameterizedTest
    @ValueSource(strings = {
        "[]",            // no index
        "[7..4]",        // upper below lower
        "[0;1]",         // resolved offset negative
        "[-1]",          // negative component
        "[1..-2]",       // negative component
        "[a]",           // non-numeric
        "[1..x]",        // non-numeric
        "[1..2;]",       // empty base
        "[1..]",         // missing upper bound
        "[..2]",         // missing lower bound
    })
    void malformedExpressionsAreRejected(String expression) {
        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse(expression, "tag" + expression));
        assertTrue(thrown.getMessage().contains("tag" + expression),
            () -> "message should name the address: " + thrown.getMessage());
    }

    // --- driver constraints ---

    @Test
    void indexBeyondTheProtocolMaximumIsRejected() {
        AddressConstraints eip = AddressConstraints.SINGLE_DIMENSION.withMaxIndex(255);

        assertDoesNotThrow(() -> ArrayNotationParser.parse("[255]", "tag[255]", eip));
        // The bound is on where the selection starts, not where it ends: a CIP request carries a
        // start index and an element count, so a long run from an encodable start is fine.
        assertDoesNotThrow(() -> ArrayNotationParser.parse("[0..300]", "tag[0..300]", eip));
        assertDoesNotThrow(() -> ArrayNotationParser.parse("[256;1]", "tag[256;1]", eip));

        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse("[256]", "tag[256]", eip));
        assertTrue(thrown.getMessage().contains("255"),
            () -> "message should name the real bound: " + thrown.getMessage());
    }

    @Test
    void moreDimensionsThanTheProtocolCarriesIsRejected() {
        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse("[1][2]", "tag[1][2]", AddressConstraints.SINGLE_DIMENSION));
        assertTrue(thrown.getMessage().contains("1"), thrown::getMessage);
    }

    @Test
    void interiorRangeIsRejectedWhereOnlyTheTrailingDimensionMaySpan() {
        AddressConstraints trailingOnly = AddressConstraints.UNCONSTRAINED
            .withOnlyTrailingDimensionMayBeRange(true);

        assertDoesNotThrow(() -> ArrayNotationParser.parse("[1][0..3]", "tag[1][0..3]", trailingOnly));
        assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse("[0..3][1]", "tag[0..3][1]", trailingOnly));

        // A one-element range is still a range. Judging this by the span let "[1..1][2]" through,
        // handing the driver a leading range it has no element count for.
        assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse("[1..1][2]", "tag[1..1][2]", trailingOnly));

        // A single index in the same position stays legal, which is what the constraint is for.
        assertDoesNotThrow(() -> ArrayNotationParser.parse("[1][2]", "tag[1][2]", trailingOnly));
    }

    // --- splitting an address ---

    @ParameterizedTest(name = "{0} -> address {1}, expression {2}")
    @CsvSource({
        "'myTag[0..7]',   'myTag',   '[0..7]'",
        "'myTag',         'myTag',   ''",
        "'a.b[2]',        'a.b',     '[2]'",
        "'40001[0..3]',   '40001',   '[0..3]'",
        "'t[1..2][0..5]', 't',       '[1..2][0..5]'",
    })
    void trailingExpressionIsSplitFromTheAddress(String input, String address, String expression) {
        assertEquals(address, ArrayNotationParser.addressPart(input));
        assertEquals(expression, ArrayNotationParser.expressionPart(input));
    }

    /**
     * Only a strictly numeric trailing run counts. OPC UA string identifiers may legitimately
     * contain brackets, and those must be left on the address.
     */
    @Test
    void nonNumericBracketsAreNotAnArrayExpression() {
        assertEquals("Some[Node]Name", ArrayNotationParser.addressPart("Some[Node]Name"));
        assertEquals("", ArrayNotationParser.expressionPart("Some[Node]Name"));
    }

    // --- rendering back ---

    @ParameterizedTest
    @ValueSource(strings = {"[4]", "[0..7]", "[4;1]", "[4..7;1]", "[1..2][0..5]", "[4..7;1][7..10;2]"})
    void renderingReproducesTheCanonicalForm(String expression) {
        assertEquals(expression, ArrayNotationParser.render(ArrayNotationParser.parse(expression, "tag")));
    }

    /**
     * Canonical form omits what is defaulted - a base of 0 - but never the range form, because
     * dropping that would turn an array of one into a scalar.
     */
    @Test
    void canonicalFormOmitsDefaultsButKeepsTheRangeForm() {
        assertEquals("[4..4]", ArrayNotationParser.render(ArrayNotationParser.parse("[4..4;0]", "tag")));
        assertEquals("[4]", ArrayNotationParser.render(ArrayNotationParser.parse("[4;0]", "tag")));
        assertEquals("[0..7]", ArrayNotationParser.render(ArrayNotationParser.parse("[0..7;0]", "tag")));
    }

    // --- guidance for addresses written before the migration ---

    @ParameterizedTest(name = "{0} -> {1}")
    @CsvSource({
        "'holding-register:1:INT[4]',   'holding-register:1[0..3]:INT'",
        "'%DB42:28.0:BYTE[8]',          '%DB42:28.0[0..7]:BYTE'",
        "'%DB1:0:STRING(40)[3]',        '%DB1:0[0..2]:STRING(40)'",
        "'D100:WORD[2]',                'D100[0..1]:WORD'",
        "'0x4020/0:DINT[4]',            '0x4020/0[0..3]:DINT'",
        "'myTag:DINT:8',                'myTag[0..7]:DINT'",
        "'foo:INT[1]',                  'foo[0]:INT'",
    })
    void anAddressWrittenBeforeTheMigrationIsRewritten(String legacy, String current) {
        assertEquals(current, ArrayNotationParser.currentFormOf(legacy).orElse(null), legacy);
    }

    @ParameterizedTest
    @ValueSource(strings = {"myTag", "myTag[0..3]:DINT", "holding-register:1[0..3]:INT", "nonsense"})
    void anAddressThatIsNotInTheOldShapeGetsNoRewrite(String address) {
        assertTrue(ArrayNotationParser.currentFormOf(address).isEmpty(), address);
    }

    @Test
    void theErrorNamesBothTheExpectedFormAndTheRewrite() {
        PlcInvalidTagException thrown =
            ArrayNotationParser.invalidAddress("holding-register:1:INT[4]", "{address}[range]:{TYPE}");

        assertTrue(thrown.getMessage().contains("holding-register:1:INT[4]"), thrown::getMessage);
        assertTrue(thrown.getMessage().contains("{address}[range]:{TYPE}"), thrown::getMessage);
        assertTrue(thrown.getMessage().contains("holding-register:1[0..3]:INT"), thrown::getMessage);
    }

    @Test
    void anErrorForSomethingElseJustNamesTheExpectedForm() {
        PlcInvalidTagException thrown = ArrayNotationParser.invalidAddress("nonsense", "{address}");
        assertFalse(thrown.getMessage().contains("now written"), thrown::getMessage);
    }

    // --- the comma spelling ---

    /**
     * Allen-Bradley and others write the dimensions of one array inside a single bracket. It is
     * the same selection, so it parses the same - and renders back in the one canonical form.
     */
    @ParameterizedTest(name = "{0} == {1}")
    @CsvSource({
        "'[0,1]',            '[0][1]'",
        "'[1..2,3..4]',      '[1..2][3..4]'",
        "'[1..2;1,3..4;1]',  '[1..2;1][3..4;1]'",
        "'[0,1,2]',          '[0][1][2]'",
        "'[1..2,3]',         '[1..2][3]'",
    })
    void theCommaSpellingIsTheSameSelection(String comma, String brackets) {
        assertEquals(
            ArrayNotationParser.parse(brackets, "tag" + brackets),
            ArrayNotationParser.parse(comma, "tag" + comma));
    }

    @ParameterizedTest(name = "{0} -> {1}")
    @CsvSource({
        "'[0,1]',           '[0][1]'",
        "'[1..2,3..4]',     '[1..2][3..4]'",
        "'[0][1]',          '[0][1]'",
    })
    void renderingAlwaysProducesOneBracketPerDimension(String written, String canonical) {
        assertEquals(canonical,
            ArrayNotationParser.render(ArrayNotationParser.parse(written, "tag")));
    }

    @Test
    void aMixtureOfBothSpellingsIsAccepted() {
        assertEquals("[0][1][2]",
            ArrayNotationParser.render(ArrayNotationParser.parse("[0,1][2]", "tag")));
    }

    @ParameterizedTest
    @ValueSource(strings = {"[0,]", "[,1]", "[0,,1]", "[0, 1]"})
    void aMalformedCommaListIsRejected(String expression) {
        assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse(expression, "tag" + expression));
    }

    /** The address split has to see the comma spelling as an array expression too. */
    @Test
    void theCommaSpellingIsSplitFromTheAddress() {
        assertEquals("myTag", ArrayNotationParser.addressPart("myTag[1..2,3..4]"));
        assertEquals("[1..2,3..4]", ArrayNotationParser.expressionPart("myTag[1..2,3..4]"));
    }

    // --- what the caller receives ---

    /**
     * getArrayInfo() describes the value the caller gets, so a consumer can decide from it alone
     * whether to render a scalar or a list. A bare index is a scalar; a range is an array even
     * when it spans a single element.
     */
    @ParameterizedTest
    @CsvSource({
        "'[1]',           true",
        "'[4]',           true",
        "'[4;1]',         true",
        "'[1][2]',        true",
        "'[1..1]',        false",
        "'[0..7]',        false",
        "'[4..7;1]',      false",
        "'[1][0..5]',     false",
        "'',              false",
    })
    void aBareIndexSelectsAScalarButARangeDoesNot(String expression, boolean scalar) {
        assertEquals(scalar, ArrayNotationParser.selectsSingleElement(expression), expression);
    }

    @Test
    void anAbsentExpressionRendersAsNothing() {
        assertEquals("", ArrayNotationParser.render(List.of()));
        assertTrue(ArrayNotationParser.parse("", "tag").isEmpty());
    }

    @Test
    void refusesARangeThatCannotBeCounted() {
        // getSize() is inclusive and computed in an int, so [0..2147483647] would wrap to a
        // negative element count from a selection the syntax accepted.
        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> ArrayNotationParser.parse("[0..2147483647]", "%test[0..2147483647]",
                AddressConstraints.UNCONSTRAINED));
        assertTrue(thrown.getMessage().contains("more than can be counted"), thrown.getMessage());

        // One below the limit still parses, and counts what it says.
        assertEquals(Integer.MAX_VALUE, ArrayNotationParser
            .parse("[0..2147483646]", "%test[0..2147483646]", AddressConstraints.UNCONSTRAINED)
            .get(0).getSize());
    }
}
