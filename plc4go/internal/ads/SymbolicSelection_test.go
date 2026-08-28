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

package ads

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

// A device declaring one symbol, MAIN.arr, as ARRAY [0..9] OF DINT at group 0x4020, offset 100.
// Ten four-byte elements, so every offset in these tests is checkable by hand.
func connectionWithArraySymbol() *Connection {
	dint := driverModel.NewAdsDataTypeTableEntryBuilder().
		WithSize(4).
		WithMainName("DINT").
		WithSecondaryName("DINT").
		MustBuild()
	array := driverModel.NewAdsDataTypeTableEntryBuilder().
		WithSize(40).
		WithMainName("ARRAY [0..9] OF DINT").
		WithSecondaryName("ARRAY [0..9] OF DINT").
		WithArrayDimensions(1).
		WithArrayInfo(driverModel.NewAdsDataTypeArrayInfo(0, 10)).
		MustBuild()
	symbol := driverModel.NewAdsSymbolTableEntryBuilder().
		WithGroup(0x4020).
		WithOffset(100).
		WithSize(40).
		WithName("MAIN.arr").
		WithDataTypeName("ARRAY [0..9] OF DINT").
		MustBuild()

	return &Connection{
		log: zerolog.Nop(),
		driverContext: &DriverContext{
			dataTypeTable: map[string]driverModel.AdsDataTypeTableEntry{
				"DINT": dint, "ARRAY [0..9] OF DINT": array,
			},
			symbolTable: map[string]driverModel.AdsSymbolTableEntry{"MAIN.arr": symbol},
		},
	}
}

func resolve(t *testing.T, address string, selection []apiModel.ArrayInfo) (*model.DirectPlcTag, error) {
	t.Helper()
	return connectionWithArraySymbol().resolveSymbolicTag(context.Background(),
		model.SymbolicPlcTag{SymbolicAddress: address, PlcTag: model.PlcTag{ArrayInfo: selection}})
}

// Without a selection the whole array is read from where the symbol table puts it - the behaviour
// every existing address relies on.
func TestSymbolicResolution_withoutASelectionReadsTheWholeArray(t *testing.T) {
	tag, err := resolve(t, "MAIN.arr", nil)

	require.NoError(t, err)
	assert.Equal(t, uint32(100), tag.IndexOffset)
	assert.Equal(t, uint32(40), tag.TransferSizeInBytes(), "all ten elements")
	require.Len(t, tag.DecodeArrayInfo(), 1)
	assert.Equal(t, uint32(10), tag.DecodeArrayInfo()[0].GetNumElements())
}

// A range moves the start to the first selected element and transfers only what it spans. Before
// this was applied the selection was parsed and then dropped: the read covered the whole array
// from offset 100, and the caller got ten values where four were asked for.
func TestSymbolicResolution_aRangeNarrowsTheReadToItsElements(t *testing.T) {
	tag, err := resolve(t, "MAIN.arr", []apiModel.ArrayInfo{
		&spiModel.DefaultArrayInfo{LowerBound: 1, UpperBound: 4, Range: true},
	})

	require.NoError(t, err)
	assert.Equal(t, uint32(104), tag.IndexOffset, "one four-byte element past the start")
	assert.Equal(t, uint32(16), tag.TransferSizeInBytes(), "four elements")
	require.Len(t, tag.DecodeArrayInfo(), 1)
	assert.Equal(t, uint32(4), tag.DecodeArrayInfo()[0].GetNumElements())
	assert.Equal(t, "DINT", tag.DataType.GetMainName(), "an element, not the array")
}

// A bare index is a scalar, and a range of one is a list of one. The two select the same element
// and differ only in shape, which is the distinction the notation exists to make.
func TestSymbolicResolution_aBareIndexIsAScalarAndAOneElementRangeIsNot(t *testing.T) {
	scalar, err := resolve(t, "MAIN.arr", []apiModel.ArrayInfo{
		&spiModel.DefaultArrayInfo{LowerBound: 3, UpperBound: 3},
	})
	require.NoError(t, err)

	listOfOne, err := resolve(t, "MAIN.arr", []apiModel.ArrayInfo{
		&spiModel.DefaultArrayInfo{LowerBound: 3, UpperBound: 3, Range: true},
	})
	require.NoError(t, err)

	assert.Equal(t, uint32(112), scalar.IndexOffset, "three elements past the start")
	assert.Equal(t, scalar.IndexOffset, listOfOne.IndexOffset, "the same element")
	assert.Equal(t, uint32(4), scalar.TransferSizeInBytes())
	assert.Equal(t, uint32(4), listOfOne.TransferSizeInBytes())

	assert.Empty(t, scalar.DecodeArrayInfo(), "decoded as a scalar")
	assert.Len(t, listOfOne.DecodeArrayInfo(), 1, "decoded as a list of one")
}

// A selection the device cannot satisfy is refused. Reading the wrong elements cannot be told
// apart from reading the right ones once the values come back, so this must not be approximated.
func TestSymbolicResolution_anOutOfBoundsSelectionIsRefused(t *testing.T) {
	_, err := resolve(t, "MAIN.arr", []apiModel.ArrayInfo{
		&spiModel.DefaultArrayInfo{LowerBound: 8, UpperBound: 12, Range: true},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "[0..9]", "the error names what the PLC declares")
}

// The direct counterpart: an address that names a memory location selects out of that location,
// so its selection is the whole of its shape. The request size was already multiplied by the
// count while the decoder was handed the scalar type's own (empty) shape, so three of these four
// elements were transferred and dropped.
func TestDirectResolution_aSelectionBecomesTheDecodedShape(t *testing.T) {
	connection := connectionWithArraySymbol()
	parsed, err := NewTagHandler().ParseTag("0x4020/100[0..3]:DINT")
	require.NoError(t, err)

	tag, err := connection.directTagFor(context.Background(), parsed)

	require.NoError(t, err)
	assert.Equal(t, uint32(16), tag.TransferSizeInBytes(), "four four-byte elements")
	require.Len(t, tag.DecodeArrayInfo(), 1, "decoded as a list")
	assert.Equal(t, uint32(4), tag.DecodeArrayInfo()[0].GetNumElements())
}

func TestDirectResolution_aScalarKeepsItsScalarShape(t *testing.T) {
	connection := connectionWithArraySymbol()
	parsed, err := NewTagHandler().ParseTag("0x4020/100:DINT")
	require.NoError(t, err)

	tag, err := connection.directTagFor(context.Background(), parsed)

	require.NoError(t, err)
	assert.Equal(t, uint32(4), tag.TransferSizeInBytes())
	assert.Empty(t, tag.DecodeArrayInfo())
}

// A second symbol, MAIN.grid, declared as ARRAY [0..9,0..4] OF DINT at group 0x4020, offset 500:
// ten rows of five four-byte elements, laid out row-major, so a row is 20 bytes.
func connectionWithGridSymbol() *Connection {
	connection := connectionWithArraySymbol()
	grid := driverModel.NewAdsDataTypeTableEntryBuilder().
		WithSize(200).
		WithMainName("ARRAY [0..9,0..4] OF DINT").
		WithSecondaryName("ARRAY [0..9,0..4] OF DINT").
		WithArrayDimensions(2).
		WithArrayInfo(
			driverModel.NewAdsDataTypeArrayInfo(0, 10),
			driverModel.NewAdsDataTypeArrayInfo(0, 5),
		).
		MustBuild()
	connection.driverContext.dataTypeTable["ARRAY [0..9,0..4] OF DINT"] = grid
	connection.driverContext.symbolTable["MAIN.grid"] = driverModel.NewAdsSymbolTableEntryBuilder().
		WithGroup(0x4020).
		WithOffset(500).
		WithSize(200).
		WithName("MAIN.grid").
		WithDataTypeName("ARRAY [0..9,0..4] OF DINT").
		MustBuild()
	return connection
}

func resolveGrid(t *testing.T, selection ...apiModel.ArrayInfo) (*model.DirectPlcTag, error) {
	t.Helper()
	return connectionWithGridSymbol().resolveSymbolicTag(context.Background(),
		model.SymbolicPlcTag{SymbolicAddress: "MAIN.grid", PlcTag: model.PlcTag{ArrayInfo: selection}})
}

// One element of a two-dimensional array: row 3, column 2, which row-major puts at
// 500 + 3*20 + 2*4.
func TestGridSelection_oneElement(t *testing.T) {
	tag, err := resolveGrid(t,
		&spiModel.DefaultArrayInfo{LowerBound: 3, UpperBound: 3},
		&spiModel.DefaultArrayInfo{LowerBound: 2, UpperBound: 2})

	require.NoError(t, err)
	assert.Equal(t, uint32(568), tag.IndexOffset)
	assert.Equal(t, uint32(4), tag.TransferSizeInBytes())
	assert.Empty(t, tag.DecodeArrayInfo(), "two bare indices name one element, which is a scalar")
}

// Part of one row is contiguous: row 3, columns 1..3.
func TestGridSelection_partOfOneRow(t *testing.T) {
	tag, err := resolveGrid(t,
		&spiModel.DefaultArrayInfo{LowerBound: 3, UpperBound: 3},
		&spiModel.DefaultArrayInfo{LowerBound: 1, UpperBound: 3, Range: true})

	require.NoError(t, err)
	assert.Equal(t, uint32(564), tag.IndexOffset, "row 3, column 1")
	assert.Equal(t, uint32(12), tag.TransferSizeInBytes(), "three elements")
	require.Len(t, tag.DecodeArrayInfo(), 1, "the bare row index collapses; the range remains")
	assert.Equal(t, uint32(3), tag.DecodeArrayInfo()[0].GetNumElements())
}

// Whole rows are contiguous too, and stay two-dimensional: rows 1..2, every column. No address
// can ask for this - the parser allows a range only in the last dimension - so this pins the
// resolver itself, which is what a tag built in code reaches.
func TestGridSelection_wholeRows(t *testing.T) {
	tag, err := resolveGrid(t,
		&spiModel.DefaultArrayInfo{LowerBound: 1, UpperBound: 2, Range: true},
		&spiModel.DefaultArrayInfo{LowerBound: 0, UpperBound: 4, Range: true})

	require.NoError(t, err)
	assert.Equal(t, uint32(520), tag.IndexOffset, "the start of row 1")
	assert.Equal(t, uint32(40), tag.TransferSizeInBytes(), "two rows of five")
	require.Len(t, tag.DecodeArrayInfo(), 2, "two rows of five, not a flat ten")
	assert.Equal(t, uint32(2), tag.DecodeArrayInfo()[0].GetNumElements())
	assert.Equal(t, uint32(5), tag.DecodeArrayInfo()[1].GetNumElements())
}

// The case the contiguity rule exists for: part of every row is ten separate runs, and the
// contiguous block a single read returns is not what was asked for. The address parser refuses a
// range before the last dimension for the same reason, so this is the resolver standing behind
// that guarantee for a tag built in code.
func TestGridSelection_partOfEveryRowIsRefused(t *testing.T) {
	_, err := resolveGrid(t,
		&spiModel.DefaultArrayInfo{LowerBound: 0, UpperBound: 9, Range: true},
		&spiModel.DefaultArrayInfo{LowerBound: 1, UpperBound: 3, Range: true})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "contiguous")
}

// Naming some of the dimensions is refused: what the unnamed ones select would be a guess.
func TestGridSelection_aPartialAddressIsRefused(t *testing.T) {
	_, err := resolveGrid(t, &spiModel.DefaultArrayInfo{LowerBound: 3, UpperBound: 3})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "every dimension")
}

// Bounds are held per dimension, not just on the first.
func TestGridSelection_anOutOfBoundsColumnIsRefused(t *testing.T) {
	_, err := resolveGrid(t,
		&spiModel.DefaultArrayInfo{LowerBound: 3, UpperBound: 3},
		&spiModel.DefaultArrayInfo{LowerBound: 3, UpperBound: 7, Range: true})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "[0..4]")
}
