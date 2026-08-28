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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The one definition of the array notation in plc4go, per the grammar contract shared with
// plc4j (specs/002-unified-array-notation/contracts/array-notation-grammar.md).
//
// The cases below are the plc4j suite's cases, deliberately the same ones rather than merely
// similar: the two bindings share no code, so the only evidence that an address means the same
// thing in both is that both satisfy the same specification against the same inputs.

// --- the semantics table ---

func TestArrayNotationParser_SingleDimensionResolvesToTheDocumentedOffsets(t *testing.T) {
	for _, tt := range []struct {
		expression string
		size       uint32
		first      uint32
		last       uint32
	}{
		{"[4]", 1, 4, 4},
		{"[0..7]", 8, 0, 7},
		{"[4;1]", 1, 3, 3},
		{"[4..7;1]", 4, 3, 6},
		{"[0]", 1, 0, 0},
		{"[7..7]", 1, 7, 7},
	} {
		t.Run(tt.expression, func(t *testing.T) {
			dimensions, err := ParseArrayExpression(tt.expression, "tag"+tt.expression, Unconstrained)
			require.NoError(t, err)
			require.Len(t, dimensions, 1)

			only := dimensions[0]
			assert.Equal(t, tt.size, only.GetSize(), "size")
			assert.Equal(t, tt.first, only.GetLowerBound()-only.GetBase(), "first offset")
			assert.Equal(t, tt.last, only.GetUpperBound()-only.GetBase(), "last offset")
		})
	}
}

// A bare index and a one-element range cover the same element but are not the same selection:
// the first yields a scalar and the second an array of one.
func TestArrayNotationParser_ABareIndexIsNotTheSameAsAOneElementRange(t *testing.T) {
	index, err := ParseArrayExpression("[4]", "tag[4]", Unconstrained)
	require.NoError(t, err)
	arrayRange, err := ParseArrayExpression("[4..4]", "tag[4..4]", Unconstrained)
	require.NoError(t, err)

	assert.NotEqual(t, arrayRange, index)
	assert.False(t, index[0].IsRange())
	assert.True(t, arrayRange[0].IsRange())
	assert.Equal(t, index[0].GetLowerBound(), arrayRange[0].GetLowerBound())
	assert.Equal(t, uint32(1), index[0].GetSize())
	assert.Equal(t, uint32(1), arrayRange[0].GetSize())
}

func TestArrayNotationParser_WrittenBoundsArePreservedNotResolved(t *testing.T) {
	dimensions, err := ParseArrayExpression("[4..7;1]", "tag[4..7;1]", Unconstrained)
	require.NoError(t, err)

	assert.Equal(t, uint32(4), dimensions[0].GetLowerBound(), "lower bound is as written")
	assert.Equal(t, uint32(7), dimensions[0].GetUpperBound(), "upper bound is as written")
	assert.Equal(t, uint32(1), dimensions[0].GetBase(), "declared base")
	assert.Equal(t, uint32(4), dimensions[0].GetSize())
}

func TestArrayNotationParser_BaseDefaultsToZero(t *testing.T) {
	dimensions, err := ParseArrayExpression("[4]", "tag[4]", Unconstrained)
	require.NoError(t, err)
	assert.Equal(t, uint32(0), dimensions[0].GetBase())
}

func TestArrayNotationParser_MultipleDimensionsKeepTheirWrittenOrder(t *testing.T) {
	dimensions, err := ParseArrayExpression("[1..2][0..5]", "tag[1..2][0..5]", Unconstrained)
	require.NoError(t, err)

	require.Len(t, dimensions, 2)
	assert.Equal(t, uint32(1), dimensions[0].GetLowerBound())
	assert.Equal(t, uint32(2), dimensions[0].GetUpperBound())
	assert.Equal(t, uint32(0), dimensions[1].GetLowerBound())
	assert.Equal(t, uint32(5), dimensions[1].GetUpperBound())
}

func TestArrayNotationParser_EachDimensionCarriesItsOwnBase(t *testing.T) {
	dimensions, err := ParseArrayExpression("[4..7;1][7..10;2]", "tag[4..7;1][7..10;2]", Unconstrained)
	require.NoError(t, err)

	require.Len(t, dimensions, 2)
	assert.Equal(t, uint32(3), dimensions[0].GetLowerBound()-dimensions[0].GetBase())
	assert.Equal(t, uint32(6), dimensions[0].GetUpperBound()-dimensions[0].GetBase())
	assert.Equal(t, uint32(5), dimensions[1].GetLowerBound()-dimensions[1].GetBase())
	assert.Equal(t, uint32(8), dimensions[1].GetUpperBound()-dimensions[1].GetBase())
}

// --- the rejection table ---

func TestArrayNotationParser_MalformedExpressionsAreRejected(t *testing.T) {
	for _, expression := range []string{
		"[]",      // no index
		"[7..4]",  // upper below lower
		"[0;1]",   // resolved offset is negative
		"[-1]",    // negative component
		"[1..-2]", // negative component
		"[a]",     // non-numeric
		"[1..x]",  // non-numeric
		"[1..2;]", // empty base
		"[1..]",   // missing upper bound
		"[..2]",   // missing lower bound
		"[0,]",    // trailing comma
		"[,1]",    // leading comma
		"[0,,1]",  // empty dimension
		"[0, 1]",  // space in the list
	} {
		t.Run(expression, func(t *testing.T) {
			_, err := ParseArrayExpression(expression, "tag"+expression, Unconstrained)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "tag"+expression, "the message must name the address")
		})
	}
}

// --- driver constraints ---

func TestArrayNotationParser_IndexBeyondTheProtocolMaximumIsRejected(t *testing.T) {
	eip := SingleDimension.WithMaxIndex(255)

	_, err := ParseArrayExpression("[255]", "tag[255]", eip)
	require.NoError(t, err)
	// The bound is on where the selection starts, not where it ends: a CIP request carries a
	// start index and an element count, so a long run from an encodable start is fine.
	_, err = ParseArrayExpression("[0..300]", "tag[0..300]", eip)
	require.NoError(t, err)
	_, err = ParseArrayExpression("[256;1]", "tag[256;1]", eip)
	require.NoError(t, err)

	_, err = ParseArrayExpression("[256]", "tag[256]", eip)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "255", "the message must name the real bound")
}

func TestArrayNotationParser_MoreDimensionsThanTheProtocolCarriesIsRejected(t *testing.T) {
	_, err := ParseArrayExpression("[1][2]", "tag[1][2]", SingleDimension)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1")
}

func TestArrayNotationParser_InteriorRangeIsRejectedWhereOnlyTheTrailingDimensionMaySpan(t *testing.T) {
	trailingOnly := Unconstrained.WithOnlyTrailingDimensionMayBeRange(true)

	_, err := ParseArrayExpression("[1][0..3]", "tag[1][0..3]", trailingOnly)
	require.NoError(t, err)

	_, err = ParseArrayExpression("[0..3][1]", "tag[0..3][1]", trailingOnly)
	require.Error(t, err)
}

// --- splitting an address ---

func TestArrayNotationParser_TrailingExpressionIsSplitFromTheAddress(t *testing.T) {
	for _, tt := range []struct{ input, address, expression string }{
		{"myTag[0..7]", "myTag", "[0..7]"},
		{"myTag", "myTag", ""},
		{"a.b[2]", "a.b", "[2]"},
		{"40001[0..3]", "40001", "[0..3]"},
		{"t[1..2][0..5]", "t", "[1..2][0..5]"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.address, AddressPart(tt.input))
			assert.Equal(t, tt.expression, ExpressionPart(tt.input))
		})
	}
}

// Only a strictly numeric trailing run counts. An identifier that happens to contain brackets is
// left on the address.
func TestArrayNotationParser_NonNumericBracketsAreNotAnArrayExpression(t *testing.T) {
	assert.Equal(t, "Some[Node]Name", AddressPart("Some[Node]Name"))
	assert.Equal(t, "", ExpressionPart("Some[Node]Name"))
}

// --- rendering back ---

func TestArrayNotationParser_RenderingReproducesTheCanonicalForm(t *testing.T) {
	for _, expression := range []string{
		"[4]", "[0..7]", "[4;1]", "[4..7;1]", "[1..2][0..5]", "[4..7;1][7..10;2]",
	} {
		t.Run(expression, func(t *testing.T) {
			dimensions, err := ParseArrayExpression(expression, "tag", Unconstrained)
			require.NoError(t, err)
			assert.Equal(t, expression, RenderArrayExpression(dimensions))
		})
	}
}

// Canonical form omits what is defaulted - a base of 0 - but never the range form, because
// dropping that would turn an array of one into a scalar.
func TestArrayNotationParser_CanonicalFormOmitsDefaultsButKeepsTheRangeForm(t *testing.T) {
	for _, tt := range []struct{ written, canonical string }{
		{"[4..4;0]", "[4..4]"},
		{"[4;0]", "[4]"},
		{"[0..7;0]", "[0..7]"},
	} {
		t.Run(tt.written, func(t *testing.T) {
			dimensions, err := ParseArrayExpression(tt.written, "tag", Unconstrained)
			require.NoError(t, err)
			assert.Equal(t, tt.canonical, RenderArrayExpression(dimensions))
		})
	}
}

func TestArrayNotationParser_AnAbsentExpressionRendersAsNothing(t *testing.T) {
	assert.Equal(t, "", RenderArrayExpression(nil))
	dimensions, err := ParseArrayExpression("", "tag", Unconstrained)
	require.NoError(t, err)
	assert.Empty(t, dimensions)
}

// --- the comma spelling ---

// Allen-Bradley and others write the dimensions of one array inside a single bracket. It is the
// same selection, so it parses the same - and renders back in the one canonical form.
func TestArrayNotationParser_TheCommaSpellingIsTheSameSelection(t *testing.T) {
	for _, tt := range []struct{ comma, brackets string }{
		{"[0,1]", "[0][1]"},
		{"[1..2,3..4]", "[1..2][3..4]"},
		{"[1..2;1,3..4;1]", "[1..2;1][3..4;1]"},
		{"[0,1,2]", "[0][1][2]"},
		{"[1..2,3]", "[1..2][3]"},
	} {
		t.Run(tt.comma, func(t *testing.T) {
			viaComma, err := ParseArrayExpression(tt.comma, "tag"+tt.comma, Unconstrained)
			require.NoError(t, err)
			viaBrackets, err := ParseArrayExpression(tt.brackets, "tag"+tt.brackets, Unconstrained)
			require.NoError(t, err)
			assert.Equal(t, viaBrackets, viaComma)
		})
	}
}

func TestArrayNotationParser_RenderingAlwaysProducesOneBracketPerDimension(t *testing.T) {
	for _, tt := range []struct{ written, canonical string }{
		{"[0,1]", "[0][1]"},
		{"[1..2,3..4]", "[1..2][3..4]"},
		{"[0][1]", "[0][1]"},
		{"[0,1][2]", "[0][1][2]"},
	} {
		t.Run(tt.written, func(t *testing.T) {
			dimensions, err := ParseArrayExpression(tt.written, "tag", Unconstrained)
			require.NoError(t, err)
			assert.Equal(t, tt.canonical, RenderArrayExpression(dimensions))
		})
	}
}

func TestArrayNotationParser_TheCommaSpellingIsSplitFromTheAddress(t *testing.T) {
	assert.Equal(t, "myTag", AddressPart("myTag[1..2,3..4]"))
	assert.Equal(t, "[1..2,3..4]", ExpressionPart("myTag[1..2,3..4]"))
}

// --- what the caller receives ---

// GetArrayInfo describes the value the caller gets, so a consumer can decide from it alone
// whether to render a scalar or a list. A bare index is a scalar; a range is an array even when
// it spans a single element.
func TestArrayNotationParser_ABareIndexSelectsAScalarButARangeDoesNot(t *testing.T) {
	for _, tt := range []struct {
		expression string
		scalar     bool
	}{
		{"[1]", true},
		{"[4]", true},
		{"[4;1]", true},
		{"[1][2]", true},
		{"[1..1]", false},
		{"[0..7]", false},
		{"[4..7;1]", false},
		{"[1][0..5]", false},
		{"", false},
	} {
		t.Run(tt.expression, func(t *testing.T) {
			assert.Equal(t, tt.scalar, SelectsSingleElement(tt.expression))
		})
	}
}

// --- guidance for addresses written before the migration ---

func TestArrayNotationParser_AnAddressWrittenBeforeTheMigrationIsRewritten(t *testing.T) {
	for _, tt := range []struct{ legacy, current string }{
		{"holding-register:1:INT[4]", "holding-register:1[0..3]:INT"},
		{"%DB42:28.0:BYTE[8]", "%DB42:28.0[0..7]:BYTE"},
		{"%DB1:0:STRING(40)[3]", "%DB1:0[0..2]:STRING(40)"},
		{"D100:WORD[2]", "D100[0..1]:WORD"},
		{"0x4020/0:DINT[4]", "0x4020/0[0..3]:DINT"},
		{"myTag:DINT:8", "myTag[0..7]:DINT"},
		{"foo:INT[1]", "foo[0]:INT"},
	} {
		t.Run(tt.legacy, func(t *testing.T) {
			current, ok := CurrentFormOf(tt.legacy)
			require.True(t, ok, tt.legacy)
			assert.Equal(t, tt.current, current)
		})
	}
}

func TestArrayNotationParser_AnAddressThatIsNotInTheOldShapeGetsNoRewrite(t *testing.T) {
	for _, address := range []string{"myTag", "myTag[0..3]:DINT", "holding-register:1[0..3]:INT", "nonsense"} {
		t.Run(address, func(t *testing.T) {
			_, ok := CurrentFormOf(address)
			assert.False(t, ok, address)
		})
	}
}

// --- round trip ---

func TestArrayNotationParser_ASelectionSurvivesBeingRenderedAndParsedAgain(t *testing.T) {
	for _, written := range []string{
		"[4]", "[0..7]", "[4;1]", "[4..7;1]", "[0]", "[7..7]",
		"[1..2][0..5]", "[4..7;1][7..10;2]", "[0][1][2]",
		"[0,1]", "[1..2,3..4]", "[0,1][2]", "[1..2;1,3..4;1]",
	} {
		t.Run(written, func(t *testing.T) {
			parsed, err := ParseArrayExpression(written, "tag"+written, Unconstrained)
			require.NoError(t, err)
			reparsed, err := ParseArrayExpression(RenderArrayExpression(parsed), "tag", Unconstrained)
			require.NoError(t, err)
			assert.Equal(t, parsed, reparsed)
		})
	}
}

// A range the syntax accepts but no count can hold must be refused at the parse, not silently
// wrapped: [0..4294967295] spans 2^32 elements, which is zero in a uint32.
func TestParseArrayExpression_refusesARangeThatCannotBeCounted(t *testing.T) {
	_, err := ParseArrayExpression("[0..4294967295]", "%test[0..4294967295]", Unconstrained)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "more than can be counted")

	// One below the limit still parses, and counts what it says.
	dimensions, err := ParseArrayExpression("[0..4294967294]", "%test[0..4294967294]", Unconstrained)
	require.NoError(t, err)
	assert.Equal(t, uint32(4294967295), dimensions[0].GetSize())
}
