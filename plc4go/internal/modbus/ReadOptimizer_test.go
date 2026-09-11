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

package modbus

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
)

// The cases here are ported from plc4j's ModbusReadOptimizerTest, including the ones that pin down
// shipped bug fixes - the merged request losing its unit-id (GH-2686) and a coil array answering
// with only its first element (GH-2060).

// tagAt builds a tag straight from its wire address, the way plc4j's test constructs one. Going
// through the address parser instead would drag the one-based logical addressing in and make every
// expectation here one off from the block addresses it is about.
func tagAt(tagType TagType, address uint16, quantity uint16, datatype readWriteModel.ModbusDataType, config tagConfig) modbusTag {
	return newTagFromWireAddress(tagType, address, quantity, datatype, 1, config, quantity > 1)
}

func coilAt(address uint16, quantity uint16) modbusTag {
	return tagAt(Coil, address, quantity, readWriteModel.ModbusDataType_BOOL, tagConfig{})
}

func holdingRegisterAt(address uint16, quantity uint16) modbusTag {
	return tagAt(HoldingRegister, address, quantity, readWriteModel.ModbusDataType_INT, tagConfig{})
}

func onUnit(unitId uint8) tagConfig {
	return tagConfig{unitId: &unitId}
}

func inByteOrder(byteOrder ByteOrder) tagConfig {
	return tagConfig{byteOrder: &byteOrder}
}

// named numbers the tags the way plc4j's test does, so that the expectations can name them.
func named(tags ...modbusTag) []namedTag {
	namedTags := make([]namedTag, 0, len(tags))
	for i, tag := range tags {
		namedTags = append(namedTags, namedTag{name: fmt.Sprintf("tag%d", i), tag: tag})
	}
	return namedTags
}

// expectedBlock is what one block is checked against. The tag names say which tags ended up in it,
// which is the half of the merging decision the addresses alone don't show.
type expectedBlock struct {
	tagType  TagType
	unitId   *uint8
	address  uint16
	quantity uint16
	tags     []string
}

func unitId(value uint8) *uint8 {
	return &value
}

func assertBlocks(t *testing.T, expected []expectedBlock, actual []readBlock) {
	t.Helper()
	require.Len(t, actual, len(expected))
	for i, want := range expected {
		got := actual[i]
		assert.Equal(t, want.tagType, got.tagType, "block %d: area", i)
		assert.Equal(t, want.unitId, got.unitId, "block %d: unit-id", i)
		assert.Equal(t, want.address, got.address, "block %d: address", i)
		assert.Equal(t, want.quantity, got.quantity, "block %d: quantity", i)
		names := make([]string, 0, len(got.tags))
		for _, member := range got.tags {
			names = append(names, member.name)
			assert.Equal(t, member.tag.Address-got.address, member.offset, "block %d: offset of %q", i, member.name)
		}
		assert.Equal(t, want.tags, names, "block %d: members", i)
	}
}

// A block read covers one contiguous run of coils, and the run is cut to whatever fits into
// maxCoilQuantity addresses measured from the first tag of the block.
func TestOptimizeReadsGroupsCoils(t *testing.T) {
	for _, test := range []struct {
		name   string
		tags   []namedTag
		blocks []expectedBlock
	}{
		{
			name:   "a single coil is read on its own",
			tags:   named(coilAt(0, 1)),
			blocks: []expectedBlock{{tagType: Coil, address: 0, quantity: 1, tags: []string{"tag0"}}},
		},
		{
			name:   "two adjacent coils become one two-coil read",
			tags:   named(coilAt(0, 1), coilAt(1, 1)),
			blocks: []expectedBlock{{tagType: Coil, address: 0, quantity: 2, tags: []string{"tag0", "tag1"}}},
		},
		{
			// The window width is the gap tolerance, so a gap of 99 coils is merged straight over.
			name:   "two coils with a large gap still become one read",
			tags:   named(coilAt(0, 1), coilAt(100, 1)),
			blocks: []expectedBlock{{tagType: Coil, address: 0, quantity: 101, tags: []string{"tag0", "tag1"}}},
		},
		{
			name: "two coils further apart than a request reaches are read separately",
			tags: named(coilAt(0, 1), coilAt(2100, 1)),
			blocks: []expectedBlock{
				{tagType: Coil, address: 0, quantity: 1, tags: []string{"tag0"}},
				{tagType: Coil, address: 2100, quantity: 1, tags: []string{"tag1"}},
			},
		},
		{
			// The window is anchored at the first tag, not slid along, so a third tag is measured
			// against the start of its block rather than against its predecessor.
			name: "the window is measured from the first tag of a block",
			tags: named(coilAt(0, 1), coilAt(1999, 1), coilAt(2001, 1)),
			blocks: []expectedBlock{
				{tagType: Coil, address: 0, quantity: 2000, tags: []string{"tag0", "tag1"}},
				{tagType: Coil, address: 2001, quantity: 1, tags: []string{"tag2"}},
			},
		},
		{
			name:   "an array occupies one coil per element",
			tags:   named(coilAt(0, 1), coilAt(6, 4)),
			blocks: []expectedBlock{{tagType: Coil, address: 0, quantity: 10, tags: []string{"tag0", "tag1"}}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertBlocks(t, test.blocks, optimizeReadsWithSpecLimits(test.tags))
		})
	}
}

// The register areas are cut the same way, only counted in registers rather than coils.
func TestOptimizeReadsGroupsRegisters(t *testing.T) {
	for _, test := range []struct {
		name   string
		tags   []namedTag
		blocks []expectedBlock
	}{
		{
			name:   "a single register is read on its own",
			tags:   named(holdingRegisterAt(0, 1)),
			blocks: []expectedBlock{{tagType: HoldingRegister, address: 0, quantity: 1, tags: []string{"tag0"}}},
		},
		{
			name:   "two adjacent registers become one two-register read",
			tags:   named(holdingRegisterAt(0, 1), holdingRegisterAt(1, 1)),
			blocks: []expectedBlock{{tagType: HoldingRegister, address: 0, quantity: 2, tags: []string{"tag0", "tag1"}}},
		},
		{
			name:   "two registers with a large gap still become one read",
			tags:   named(holdingRegisterAt(0, 1), holdingRegisterAt(100, 1)),
			blocks: []expectedBlock{{tagType: HoldingRegister, address: 0, quantity: 101, tags: []string{"tag0", "tag1"}}},
		},
		{
			name: "two registers further apart than a request reaches are read separately",
			tags: named(holdingRegisterAt(0, 1), holdingRegisterAt(2100, 1)),
			blocks: []expectedBlock{
				{tagType: HoldingRegister, address: 0, quantity: 1, tags: []string{"tag0"}},
				{tagType: HoldingRegister, address: 2100, quantity: 1, tags: []string{"tag1"}},
			},
		},
		{
			// A DINT is two registers wide, so a single one of them spans registers 1 and 2 -
			// the block has to cover both or the tag cannot be decoded out of it.
			name: "a wide data type occupies the registers it needs",
			tags: named(tagAt(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_DINT, tagConfig{})),
			blocks: []expectedBlock{
				{tagType: HoldingRegister, address: 1, quantity: 2, tags: []string{"tag0"}},
			},
		},
		{
			// The tags arrive unsorted; the block is still anchored at the lowest address.
			name: "tags are sorted by address before they are cut into blocks",
			tags: named(holdingRegisterAt(5, 1), holdingRegisterAt(1, 1)),
			blocks: []expectedBlock{
				{tagType: HoldingRegister, address: 1, quantity: 5, tags: []string{"tag1", "tag0"}},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertBlocks(t, test.blocks, optimizeReadsWithSpecLimits(test.tags))
		})
	}
}

// Five areas, addressed independently and read with different function codes, so address 0 means
// five unrelated things and a block may never span two of them.
func TestOptimizeReadsNeverMergesAcrossAreas(t *testing.T) {
	blocks := optimizeReadsWithSpecLimits(named(
		coilAt(0, 1),
		tagAt(HoldingRegister, 0, 1, readWriteModel.ModbusDataType_INT, tagConfig{}),
		tagAt(InputRegister, 0, 1, readWriteModel.ModbusDataType_INT, tagConfig{}),
		tagAt(ExtendedRegister, 0, 1, readWriteModel.ModbusDataType_INT, tagConfig{}),
		tagAt(DiscreteInput, 0, 1, readWriteModel.ModbusDataType_BOOL, tagConfig{}),
	))

	assertBlocks(t, []expectedBlock{
		{tagType: Coil, address: 0, quantity: 1, tags: []string{"tag0"}},
		{tagType: HoldingRegister, address: 0, quantity: 1, tags: []string{"tag1"}},
		{tagType: InputRegister, address: 0, quantity: 1, tags: []string{"tag2"}},
		{tagType: ExtendedRegister, address: 0, quantity: 1, tags: []string{"tag3"}},
		{tagType: DiscreteInput, address: 0, quantity: 1, tags: []string{"tag4"}},
	}, blocks)
}

// A block read reaches exactly one unit. Merging tags that address different units would send the
// request to whichever one won, and dropping the unit-id from the merged request would send it to
// the connection's default - both of which read the wrong device (plc4j GH-2686).
func TestOptimizeReadsKeepsUnitsApart(t *testing.T) {
	for _, test := range []struct {
		name   string
		tags   []namedTag
		blocks []expectedBlock
	}{
		{
			name: "registers on different units are not merged",
			tags: named(
				tagAt(HoldingRegister, 0, 1, readWriteModel.ModbusDataType_INT, onUnit(2)),
				tagAt(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_INT, onUnit(3)),
			),
			blocks: []expectedBlock{
				{tagType: HoldingRegister, unitId: unitId(2), address: 0, quantity: 1, tags: []string{"tag0"}},
				{tagType: HoldingRegister, unitId: unitId(3), address: 1, quantity: 1, tags: []string{"tag1"}},
			},
		},
		{
			name: "coils on different units are not merged",
			tags: named(
				tagAt(Coil, 0, 1, readWriteModel.ModbusDataType_BOOL, onUnit(2)),
				tagAt(Coil, 1, 1, readWriteModel.ModbusDataType_BOOL, onUnit(3)),
			),
			blocks: []expectedBlock{
				{tagType: Coil, unitId: unitId(2), address: 0, quantity: 1, tags: []string{"tag0"}},
				{tagType: Coil, unitId: unitId(3), address: 1, quantity: 1, tags: []string{"tag1"}},
			},
		},
		{
			name: "tags on the same unit are merged and the merged block keeps the unit-id",
			tags: named(
				tagAt(HoldingRegister, 0, 1, readWriteModel.ModbusDataType_INT, onUnit(7)),
				tagAt(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_INT, onUnit(7)),
			),
			blocks: []expectedBlock{
				{tagType: HoldingRegister, unitId: unitId(7), address: 0, quantity: 2, tags: []string{"tag0", "tag1"}},
			},
		},
		{
			// "whatever the connection uses" is not the same statement as a number, so a tag that
			// named no unit forms its own group even if the connection's default is that number.
			name: "a tag without a unit-id is not merged with one that has it",
			tags: named(
				tagAt(HoldingRegister, 0, 1, readWriteModel.ModbusDataType_INT, tagConfig{}),
				tagAt(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_INT, onUnit(2)),
			),
			blocks: []expectedBlock{
				{tagType: HoldingRegister, unitId: nil, address: 0, quantity: 1, tags: []string{"tag0"}},
				{tagType: HoldingRegister, unitId: unitId(2), address: 1, quantity: 1, tags: []string{"tag1"}},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertBlocks(t, test.blocks, optimizeReadsWithSpecLimits(test.tags))
		})
	}
}

// The optimizer only ever merges requests, it never takes one apart: a tag wider than a single
// request can carry passes through whole rather than being split into halves of a value. Such a
// tag is rejected when its address is parsed, so it only reaches here from a caller that built it
// by hand - but silently halving it would be worse than passing it on.
func TestOptimizeReadsPassesAnOversizedTagThrough(t *testing.T) {
	blocks := optimizeReadsWithSpecLimits(named(
		tagAt(HoldingRegister, 0, 200, readWriteModel.ModbusDataType_INT, tagConfig{}),
		tagAt(HoldingRegister, 300, 1, readWriteModel.ModbusDataType_INT, tagConfig{}),
	))

	assertBlocks(t, []expectedBlock{
		{tagType: HoldingRegister, address: 0, quantity: 200, tags: []string{"tag0"}},
		{tagType: HoldingRegister, address: 300, quantity: 1, tags: []string{"tag1"}},
	}, blocks)
}

// The ceilings are parameters, so a gateway that cannot take a full-width request can be given a
// narrower one without the grouping changing shape.
func TestOptimizeReadsHonoursTheGivenLimits(t *testing.T) {
	tags := named(holdingRegisterAt(0, 1), holdingRegisterAt(10, 1))

	assertBlocks(t, []expectedBlock{
		{tagType: HoldingRegister, address: 0, quantity: 11, tags: []string{"tag0", "tag1"}},
	}, optimizeReads(tags, maxCoilQuantity, 20))

	assertBlocks(t, []expectedBlock{
		{tagType: HoldingRegister, address: 0, quantity: 1, tags: []string{"tag0"}},
		{tagType: HoldingRegister, address: 10, quantity: 1, tags: []string{"tag1"}},
	}, optimizeReads(tags, maxCoilQuantity, 5))
}

func TestOptimizeReadsOnNoTags(t *testing.T) {
	assert.Empty(t, optimizeReadsWithSpecLimits(nil))
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Splitting the response
////////////////////////////////////////////////////////////////////////////////////////////////////

// boolsOf unpacks a list of BOOLs. PlcList.GetBoolArray only answers for a list of one, so the
// elements are read individually, the way plc4j's test walks getIndex(i).
func boolsOf(value apiValues.PlcValue) []bool {
	bools := make([]bool, 0, value.GetLength())
	for _, element := range value.GetList() {
		bools = append(bools, element.GetBool())
	}
	return bools
}

// splitSingle runs one tag through the optimizer and splits the given payload for the single block
// it is expected to produce.
func splitSingle(t *testing.T, tag modbusTag, blockData []byte) tagReadResult {
	t.Helper()
	results := splitAll(t, named(tag), blockData, BigEndianOrder)
	require.Len(t, results, 1)
	return results[0]
}

// splitAll runs the given tags through the optimizer and splits the given payload for the single
// block they are expected to produce.
func splitAll(t *testing.T, tags []namedTag, blockData []byte, defaultByteOrder ByteOrder) []tagReadResult {
	t.Helper()
	blocks := optimizeReadsWithSpecLimits(tags)
	require.Len(t, blocks, 1)
	return splitBlockResponse(context.Background(), blocks[0], apiModel.PlcResponseCode_OK, blockData, defaultByteOrder)
}

// A coil array occupies n coils and the block read covers all of them, so all n values have to
// come back - answering with the first one used to drop the rest on the floor (plc4j GH-2060).
func TestSplitBlockResponseReturnsEveryElementOfACoilArray(t *testing.T) {
	// Eight coils from address 0, packed least significant bit first: 1,0,0,0,1,1,0,1.
	result := splitSingle(t, coilAt(0, 8), []byte{0b10110001})

	require.Equal(t, apiModel.PlcResponseCode_OK, result.responseCode)
	require.True(t, result.value.IsList())
	assert.Equal(t, []bool{true, false, false, false, true, true, false, true}, boolsOf(result.value))
}

// An array sitting inside the block has to be read from the right bit and keep reading into the
// following byte when it runs past bit 7.
func TestSplitBlockResponseReadsACoilArrayAcrossAByteBoundary(t *testing.T) {
	// The scalar coil at address 0 anchors the block there, so the array at address 6 really does
	// start at bit 6 of the first byte. Bits 6,7 of byte 0 are 0,1 and bits 0,1 of byte 1 are 0,1.
	results := splitAll(t, []namedTag{
		{name: "scalar", tag: coilAt(0, 1)},
		{name: "array", tag: coilAt(6, 4)},
	}, []byte{0b10000000, 0b00000010}, BigEndianOrder)

	require.Len(t, results, 2)
	assert.Equal(t, "array", results[1].name)
	assert.Equal(t, []bool{false, true, false, true}, boolsOf(results[1].value))
}

// A single coil keeps reading as a plain BOOL rather than as a list of one, so merging a tag into
// a block never changes the type the caller sees.
func TestSplitBlockResponseReturnsAScalarForASingleCoil(t *testing.T) {
	result := splitSingle(t, coilAt(0, 1), []byte{0b00000001})

	require.Equal(t, apiModel.PlcResponseCode_OK, result.responseCode)
	assert.False(t, result.value.IsList())
	assert.True(t, result.value.GetBool())
}

// A coil holds a single bit, so anything wider has to be reported rather than silently answered
// with the first bit.
func TestSplitBlockResponseReportsANonBoolCoilAsUnsupported(t *testing.T) {
	result := splitSingle(t, tagAt(Coil, 0, 1, readWriteModel.ModbusDataType_INT, tagConfig{}), []byte{0b00000001, 0b00000000})

	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, result.responseCode)
	assert.Nil(t, result.value)
	assert.ErrorContains(t, result.err, "only BOOL is")
}

// A device answering with fewer coils than it was asked for must fail that tag, not panic.
func TestSplitBlockResponseOnAShortCoilResponse(t *testing.T) {
	result := splitSingle(t, coilAt(0, 16), []byte{0b00000001})

	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.responseCode)
	assert.Nil(t, result.value)
}

// The registers of the block are sliced at the member's own offset, so every tag gets the value it
// asked for and not the one at the start of the block.
func TestSplitBlockResponseSlicesRegistersAtTheTagsOffset(t *testing.T) {
	results := splitAll(t, []namedTag{
		{name: "first", tag: holdingRegisterAt(0, 1)},
		{name: "third", tag: holdingRegisterAt(2, 1)},
	}, []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x03}, BigEndianOrder)

	require.Len(t, results, 2)
	assert.Equal(t, apiModel.PlcResponseCode_OK, results[0].responseCode)
	assert.Equal(t, int16(1), results[0].value.GetInt16())
	assert.Equal(t, apiModel.PlcResponseCode_OK, results[1].responseCode)
	assert.Equal(t, int16(3), results[1].value.GetInt16())
}

// A tag wider than one register is decoded across the registers it occupies.
func TestSplitBlockResponseDecodesAWideTagFromTheBlock(t *testing.T) {
	results := splitAll(t, []namedTag{
		{name: "word", tag: holdingRegisterAt(0, 1)},
		{name: "double", tag: tagAt(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_DINT, tagConfig{})},
	}, []byte{0x00, 0x01, 0x00, 0x01, 0x00, 0x00}, BigEndianOrder)

	require.Len(t, results, 2)
	assert.Equal(t, int32(0x00010000), results[1].value.GetInt32())
}

// A block may mix byte orders: the order is a property of the tag, and only a tag that declared
// none falls back to the connection's default.
func TestSplitBlockResponseAppliesEachTagsOwnByteOrder(t *testing.T) {
	results := splitAll(t, []namedTag{
		{name: "default", tag: holdingRegisterAt(0, 1)},
		{name: "little", tag: tagAt(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_INT, inByteOrder(LittleEndianOrder))},
	}, []byte{0x00, 0x2A, 0x2A, 0x00}, BigEndianOrder)

	require.Len(t, results, 2)
	assert.Equal(t, int16(42), results[0].value.GetInt16(), "the connection default applies to a tag without one of its own")
	assert.Equal(t, int16(42), results[1].value.GetInt16(), "the tag's own order applies to it")
}

// A device answering with fewer registers than it was asked for must fail the tags that fall off
// the end rather than reading past the payload.
func TestSplitBlockResponseOnAShortRegisterResponse(t *testing.T) {
	results := splitAll(t, []namedTag{
		{name: "present", tag: holdingRegisterAt(0, 1)},
		{name: "missing", tag: holdingRegisterAt(1, 1)},
	}, []byte{0x00, 0x01}, BigEndianOrder)

	require.Len(t, results, 2)
	assert.Equal(t, apiModel.PlcResponseCode_OK, results[0].responseCode)
	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, results[1].responseCode)
	assert.Nil(t, results[1].value)
}

// Merging reads means a device that rejects one address takes the whole block down with it, so
// every member is answered with the block's own code.
func TestSplitBlockResponseFailsEveryTagOfAFailedBlock(t *testing.T) {
	blocks := optimizeReadsWithSpecLimits(named(holdingRegisterAt(0, 1), holdingRegisterAt(1, 1)))
	require.Len(t, blocks, 1)

	results := splitBlockResponse(context.Background(), blocks[0], apiModel.PlcResponseCode_REMOTE_ERROR, nil, BigEndianOrder)

	require.Len(t, results, 2)
	for _, result := range results {
		assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, result.responseCode, result.name)
		assert.Nil(t, result.value, result.name)
	}
}
