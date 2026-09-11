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
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
)

// ---------------------------------------------------------------------------
// invariants every optimizeReads result must satisfy
// ---------------------------------------------------------------------------

func checkInvariants(t *testing.T, in []namedTag, blocks []readBlock, maxCoils, maxRegs uint16) {
	t.Helper()

	seen := map[string]int{}
	for _, block := range blocks {
		require.NotEmpty(t, block.tags, "a block with no tags was emitted")
		limit := int(maxPerRequestFor(block.tagType, maxCoils, maxRegs))
		for _, member := range block.tags {
			seen[member.name]++
			// offset must be the stated subtraction
			assert.Equal(t, member.tag.Address-block.address, member.offset,
				"offset of %q in block %+v", member.name, block)
			// coverage: the block must actually contain the whole tag
			span := int(tagSpan(member.tag))
			tagStart := int(member.tag.Address)
			tagEnd := tagStart + span
			blockStart := int(block.address)
			blockEnd := blockStart + int(block.quantity)
			assert.GreaterOrEqual(t, tagStart, blockStart,
				"tag %q starts before its block (block %+v)", member.name, block)
			assert.LessOrEqual(t, tagEnd, blockEnd,
				"block %+v does NOT cover tag %q (tag %d..%d, block %d..%d)",
				block, member.name, tagStart, tagEnd, blockStart, blockEnd)
		}
		// width: either within the ceiling, or a single pass-through tag
		if int(block.quantity) > limit && len(block.tags) != 1 {
			t.Errorf("block %+v is %d wide, over the %d ceiling, with %d tags",
				block, block.quantity, limit, len(block.tags))
		}
		// Must not run off the end of the 16 bit address space -- unless a member tag already did.
		//
		// The optimizer may never MANUFACTURE an out-of-range request by merging in-range tags;
		// that is the bug this guards. But a single tag that is itself out of range passes through
		// unsplit (the documented contract), and emitting it verbatim is the better of the bad
		// options: the device answers an illegal-data-address exception, which is loud. Clamping
		// it would return fewer coils than asked for with an OK code, which reads as success.
		if int(block.address)+int(block.quantity) > 65536 {
			memberAlreadyOutOfRange := false
			for _, member := range block.tags {
				if int(member.tag.Address)+int(tagSpan(member.tag)) > 65536 {
					memberAlreadyOutOfRange = true
					break
				}
			}
			if !memberAlreadyOutOfRange {
				t.Errorf("block %+v was merged past the 65535 wire-address ceiling (ends at %d) from in-range tags",
					block, int(block.address)+int(block.quantity))
			}
		}
	}
	for _, tag := range in {
		if seen[tag.name] != 1 {
			t.Errorf("tag %q appears %d times in the output, want exactly 1", tag.name, seen[tag.name])
		}
	}
	assert.Len(t, seen, len(in), "output names %v vs input count %d", seen, len(in))
}

// ---------------------------------------------------------------------------
// 1. degenerate inputs
// ---------------------------------------------------------------------------

func TestAdvEmptyAndSingle(t *testing.T) {
	assert.Empty(t, optimizeReadsWithSpecLimits(nil))
	assert.Empty(t, optimizeReadsWithSpecLimits([]namedTag{}))

	single := named(coilAt(7, 1))
	blocks := optimizeReadsWithSpecLimits(single)
	require.Len(t, blocks, 1)
	checkInvariants(t, single, blocks, maxCoilQuantity, maxRegisterQuantity)
}

// Two tags that carry the same name. A PlcReadRequest keys tags by name so this cannot arrive
// from the API, but the optimizer accepts a slice and has no guard.
func TestAdvDuplicateNames(t *testing.T) {
	in := []namedTag{
		{name: "same", tag: holdingRegisterAt(0, 1)},
		{name: "same", tag: holdingRegisterAt(5, 1)},
	}
	blocks := optimizeReadsWithSpecLimits(in)
	total := 0
	for _, b := range blocks {
		total += len(b.tags)
	}
	assert.Equal(t, 2, total, "both entries must survive into blocks")

	results := splitBlockResponse(context.Background(), blocks[0], apiModel.PlcResponseCode_OK,
		make([]byte, int(blocks[0].quantity)*2), BigEndianOrder)
	byName := map[string]int{}
	for _, r := range results {
		byName[r.name]++
	}
	t.Logf("splitBlockResponse produced %d results, distinct names %d: %v", len(results), len(byName), byName)
	assert.Equal(t, 1, len(byName), "DEFECT if >1: two results share one name, a map keyed by name loses one")
}

// Identical addresses, different names.
func TestAdvDuplicateAddresses(t *testing.T) {
	in := named(holdingRegisterAt(10, 1), holdingRegisterAt(10, 1), holdingRegisterAt(10, 1))
	blocks := optimizeReadsWithSpecLimits(in)
	checkInvariants(t, in, blocks, maxCoilQuantity, maxRegisterQuantity)
	require.Len(t, blocks, 1)
	assert.Equal(t, uint16(1), blocks[0].quantity)
	assert.Len(t, blocks[0].tags, 3)
}

// ---------------------------------------------------------------------------
// 2. the top of the address space
// ---------------------------------------------------------------------------

func TestAdvTagAt65535(t *testing.T) {
	in := named(holdingRegisterAt(65535, 1))
	blocks := optimizeReadsWithSpecLimits(in)
	require.Len(t, blocks, 1)
	t.Logf("block: addr=%d qty=%d (ends at %d)", blocks[0].address, blocks[0].quantity,
		int(blocks[0].address)+int(blocks[0].quantity))
	checkInvariants(t, in, blocks, maxCoilQuantity, maxRegisterQuantity)
}

// Merging tags near the ceiling: legally-parseable tags only.
func TestAdvMergeNearCeiling(t *testing.T) {
	in := named(holdingRegisterAt(65420, 1), holdingRegisterAt(65534, 1))
	blocks := optimizeReadsWithSpecLimits(in)
	for _, b := range blocks {
		t.Logf("block: addr=%d qty=%d end=%d tags=%d", b.address, b.quantity,
			int(b.address)+int(b.quantity), len(b.tags))
	}
	checkInvariants(t, in, blocks, maxCoilQuantity, maxRegisterQuantity)
}

// A coil block anchored at the very top with an oversized pass-through tag: the block's
// quantity is fine but address+quantity wraps past the wire ceiling.
func TestAdvOversizedTagAtCeiling(t *testing.T) {
	in := named(coilAt(65535, 2000))
	blocks := optimizeReadsWithSpecLimits(in)
	require.Len(t, blocks, 1)
	t.Logf("block addr=%d qty=%d -> reads coils %d..%d",
		blocks[0].address, blocks[0].quantity,
		blocks[0].address, int(blocks[0].address)+int(blocks[0].quantity)-1)
	checkInvariants(t, in, blocks, maxCoilQuantity, maxRegisterQuantity)
}

// ---------------------------------------------------------------------------
// 3. uint16 truncation of the block quantity
// ---------------------------------------------------------------------------

// A hand-built tag whose span does not fit into uint16. newReadBlock stores
// quantity: uint16(last - first) with no check.
func TestAdvQuantityTruncation(t *testing.T) {
	// STRING(65535) x 3 in a holding register -> 98303 registers
	wide := newTagFromWireAddress(HoldingRegister, 0, 3,
		readWriteModel.ModbusDataType_STRING, 65535, tagConfig{}, true)
	span := tagSpan(wide)
	t.Logf("tagSpan = %d registers (does not fit uint16: %v)", span, span > 65535)

	in := []namedTag{{name: "wide", tag: wide}}
	blocks := optimizeReadsWithSpecLimits(in)
	require.Len(t, blocks, 1)
	t.Logf("block addr=%d quantity=%d  (span %d)", blocks[0].address, blocks[0].quantity, span)

	// The span cannot be asked for: a Modbus quantity is 16 bits, so 98303 is not expressible at
	// all. The requirement is therefore not "preserve the span" but "never wrap": the block must
	// come out clamped to the widest representable request, never as quantity 0, which is a
	// request for nothing that a device answers happily and which would read as success.
	assert.NotZero(t, blocks[0].quantity, "an oversized span must never wrap to a request for nothing")
	assert.Equal(t, uint16(maxWireAddress), blocks[0].quantity,
		"an oversized span must be clamped to the widest representable quantity")
}

// The same truncation, but landing on quantity == 0, i.e. a request that asks for nothing.
func TestAdvQuantityTruncatesToZero(t *testing.T) {
	// STRING(32768) x 2 -> 65536 bytes -> 32768 registers. Need exactly 65536 registers:
	// WSTRING(32768) x 2 -> widthBits 32768*16 = 524288 -> 2 values = 1048576 bits
	// = 131072 bytes -> 65536 registers.
	wide := newTagFromWireAddress(HoldingRegister, 0, 2,
		readWriteModel.ModbusDataType_WSTRING, 32768, tagConfig{}, true)
	span := tagSpan(wide)
	t.Logf("tagSpan = %d registers", span)
	require.Equal(t, uint64(65536), span, "test setup: want a span of exactly 65536")

	blocks := optimizeReadsWithSpecLimits([]namedTag{{name: "wide", tag: wide}})
	require.Len(t, blocks, 1)
	t.Logf("block addr=%d quantity=%d", blocks[0].address, blocks[0].quantity)
	assert.NotEqual(t, uint16(0), blocks[0].quantity,
		"DEFECT: a block that asks for 65536 registers is emitted as quantity 0")
}

// ---------------------------------------------------------------------------
// 4. degenerate ceilings from config
// ---------------------------------------------------------------------------

func TestAdvZeroCeiling(t *testing.T) {
	done := make(chan []readBlock, 1)
	in := named(coilAt(0, 1), coilAt(1, 1), coilAt(2, 1), holdingRegisterAt(0, 1))
	go func() { done <- optimizeReads(in, 0, 0) }()
	select {
	case blocks := <-done:
		for _, b := range blocks {
			t.Logf("block %v addr=%d qty=%d tags=%d", b.tagType, b.address, b.quantity, len(b.tags))
		}
		checkInvariants(t, in, blocks, 0, 0)
	case <-t.Context().Done():
		t.Fatal("optimizeReads hung with a zero ceiling")
	}
}

func TestAdvCeilingOfOne(t *testing.T) {
	in := named(holdingRegisterAt(0, 1), holdingRegisterAt(1, 1), holdingRegisterAt(2, 1))
	blocks := optimizeReads(in, 1, 1)
	checkInvariants(t, in, blocks, 1, 1)
	assert.Len(t, blocks, 3, "a ceiling of one register cannot merge anything")
}

// A ceiling of 65535 with a tag near the top: window = address + 65535 overflows uint16 but is
// held as int, so this only proves the int widening actually works.
func TestAdvMaxCeiling(t *testing.T) {
	in := named(holdingRegisterAt(65000, 1), holdingRegisterAt(65535, 1))
	blocks := optimizeReads(in, 65535, 65535)
	for _, b := range blocks {
		t.Logf("block addr=%d qty=%d end=%d", b.address, b.quantity, int(b.address)+int(b.quantity))
	}
	checkInvariants(t, in, blocks, 65535, 65535)
}

// ---------------------------------------------------------------------------
// 5. extended registers and the 10000-register file boundary
// ---------------------------------------------------------------------------

func extendedAt(address uint16, quantity uint16) modbusTag {
	return tagAt(ExtendedRegister, address, quantity, readWriteModel.ModbusDataType_INT, tagConfig{})
}

// A block the optimizer produces has to be expressible as FC 0x14 items. The response item's
// dataLength is a uint8 carrying COUNT(data)+1, and the response's own byteCount is a uint8
// carrying the total size of all items.
func TestAdvExtendedBlockFitsTheFileRecordFraming(t *testing.T) {
	for _, test := range []struct {
		name string
		tags []namedTag
	}{
		{"straddles the boundary", named(extendedAt(9995, 1), extendedAt(10100, 1))},
		{"inside one file", named(extendedAt(100, 1), extendedAt(220, 1))},
		{"at the very top", named(extendedAt(65410, 1), extendedAt(65534, 1))},
	} {
		t.Run(test.name, func(t *testing.T) {
			blocks := optimizeReadsWithSpecLimits(test.tags)
			checkInvariants(t, test.tags, blocks, maxCoilQuantity, maxRegisterQuantity)
			for _, b := range blocks {
				groups := splitExtendedRegister(b.address, b.quantity)
				byteCount := 0
				for _, g := range groups {
					itemData := int(g.lengthWords) * 2
					// dataLength is uint8 = COUNT(data)+1
					if itemData+1 > 255 {
						t.Errorf("item %+v needs dataLength %d, which does not fit a uint8", g, itemData+1)
					}
					byteCount += 2 + itemData // dataLength + referenceType + data
				}
				pdu := 2 + byteCount // functionCode + byteCount + items
				t.Logf("block addr=%d qty=%d -> %d file group(s), response byteCount=%d, PDU=%d bytes",
					b.address, b.quantity, len(groups), byteCount, pdu)
				if byteCount > 255 {
					t.Errorf("DEFECT: response byteCount %d does not fit the uint8 field (block %+v)", byteCount, b)
				}
				if pdu > 253 {
					t.Errorf("DEFECT: FC 0x14 response PDU would be %d bytes, over the 253 byte modbus limit (block %+v)", pdu, b)
				}
			}
		})
	}
}

// The worst case the ceiling allows: a full 125-register extended block straddling a file
// boundary.
func TestAdvExtendedFullWidthBlock(t *testing.T) {
	// 9990..10114 inclusive is 125 registers - the plain-register ceiling. Extended registers
	// cannot be that wide, so the optimizer must NOT merge these into one block. Asserting a
	// single block here is what the defective version did; assert the invariant instead.
	in := named(extendedAt(9990, 1), extendedAt(10114, 1))
	blocks := optimizeReadsWithSpecLimits(in)
	require.NotEmpty(t, blocks)
	for _, b := range blocks {
		t.Logf("block addr=%d qty=%d", b.address, b.quantity)
		groups := splitExtendedRegister(b.address, b.quantity)
		byteCount := 0
		for _, g := range groups {
			t.Logf("  file=%d record=%d words=%d", g.fileNumber, g.recordNumber, g.lengthWords)
			byteCount += 2 + int(g.lengthWords)*2
		}
		pdu := 2 + byteCount
		t.Logf("response byteCount=%d PDU=%d", byteCount, pdu)
		assert.LessOrEqual(t, pdu, 253,
			"every extended block the optimizer emits must answer within one modbus PDU")
	}
}

// splitExtendedRegister on a block whose address is 0. Extended registers are addressed from
// zero on the wire, so this is the legitimate bottom of the area.
func TestAdvExtendedAtZeroAndTop(t *testing.T) {
	require.NotPanics(t, func() {
		g := splitExtendedRegister(0, 1)
		t.Logf("addr 0: %+v", g)
	})
	require.NotPanics(t, func() {
		g := splitExtendedRegister(65535, 125)
		t.Logf("addr 65535 len 125: %+v", g)
	})
	require.NotPanics(t, func() {
		g := splitExtendedRegister(0, 0)
		t.Logf("len 0: %+v (nil groups -> a request with no items)", g)
	})
	// A request whose length wraps the file numbering.
	require.NotPanics(t, func() {
		g := splitExtendedRegister(65535, 65535)
		t.Logf("addr 65535 len 65535: %d groups, last=%+v", len(g), g[len(g)-1])
	})
}

// ---------------------------------------------------------------------------
// 6. randomised property test
// ---------------------------------------------------------------------------

func TestAdvPropertyEveryTagAppearsExactlyOnce(t *testing.T) {
	types := []TagType{Coil, DiscreteInput, HoldingRegister, InputRegister, ExtendedRegister}
	dtypes := []readWriteModel.ModbusDataType{
		readWriteModel.ModbusDataType_BOOL,
		readWriteModel.ModbusDataType_INT,
		readWriteModel.ModbusDataType_DINT,
		readWriteModel.ModbusDataType_LREAL,
	}
	rng := rand.New(rand.NewSource(1))
	for iteration := range 3000 {
		n := rng.Intn(30)
		in := make([]namedTag, 0, n)
		for i := range n {
			tagType := types[rng.Intn(len(types))]
			dt := dtypes[rng.Intn(len(dtypes))]
			if tagType == Coil || tagType == DiscreteInput {
				dt = readWriteModel.ModbusDataType_BOOL
			}
			addr := uint16(rng.Intn(65536))
			quantity := uint16(1 + rng.Intn(8))
			// keep the tag inside the address space, the way the parser would
			span := registerCountOf(tagType, dt, uint64(quantity), 1)
			if uint64(addr)+span > 65535 {
				addr = uint16(65535 - span)
			}
			cfg := tagConfig{}
			if rng.Intn(3) == 0 {
				u := uint8(rng.Intn(4))
				cfg.unitId = &u
			}
			in = append(in, namedTag{
				name: fmt.Sprintf("t%02d", i),
				tag:  newTagFromWireAddress(tagType, addr, quantity, dt, 1, cfg, quantity > 1),
			})
		}
		maxCoils := uint16(1 + rng.Intn(2500))
		maxRegs := uint16(1 + rng.Intn(200))
		blocks := optimizeReads(in, maxCoils, maxRegs)
		t.Run(fmt.Sprintf("iteration%d", iteration), func(t *testing.T) {
			if t.Failed() {
				return
			}
			checkInvariants(t, in, blocks, maxCoils, maxRegs)
		})
		if t.Failed() {
			t.Logf("failing input (maxCoils=%d maxRegs=%d):", maxCoils, maxRegs)
			for _, tag := range in {
				t.Logf("  %s: type=%v addr=%d qty=%d dt=%v unit=%v",
					tag.name, tag.tag.TagType, tag.tag.Address, tag.tag.Quantity, tag.tag.Datatype, tag.tag.UnitId)
			}
			t.FailNow()
		}
	}
}

// Shuffling the input must not change the output: a driver that reorders its requests between
// reads is untestable.
func TestAdvDeterministicUnderShuffling(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	base := named(
		coilAt(0, 1), coilAt(2500, 1), holdingRegisterAt(0, 1), holdingRegisterAt(300, 1),
		tagAt(InputRegister, 5, 1, readWriteModel.ModbusDataType_INT, onUnit(3)),
		tagAt(InputRegister, 6, 1, readWriteModel.ModbusDataType_INT, onUnit(1)),
		extendedAt(9995, 1), extendedAt(10005, 1),
		tagAt(DiscreteInput, 9, 1, readWriteModel.ModbusDataType_BOOL, tagConfig{}),
	)
	want := render(optimizeReadsWithSpecLimits(base))
	for range 200 {
		shuffled := append([]namedTag(nil), base...)
		rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		got := render(optimizeReadsWithSpecLimits(shuffled))
		require.Equal(t, want, got, "block layout changed when the input order changed")
	}
}

func render(blocks []readBlock) string {
	out := ""
	for _, b := range blocks {
		unit := "-"
		if b.unitId != nil {
			unit = fmt.Sprint(*b.unitId)
		}
		names := make([]string, 0, len(b.tags))
		for _, m := range b.tags {
			names = append(names, fmt.Sprintf("%s@%d", m.name, m.offset))
		}
		out += fmt.Sprintf("[%v u=%s a=%d q=%d %v]", b.tagType, unit, b.address, b.quantity, names)
	}
	return out
}

// ---------------------------------------------------------------------------
// 7. splitBlockResponse against hand-built tags
// ---------------------------------------------------------------------------

func TestAdvSplitBitsWithZeroQuantity(t *testing.T) {
	member := blockTag{name: "zero", tag: newTagFromWireAddress(Coil, 0, 0,
		readWriteModel.ModbusDataType_BOOL, 1, tagConfig{}, false), offset: 0}
	require.NotPanics(t, func() {
		r := splitBits(member, []byte{0xFF})
		t.Logf("quantity 0 -> code=%v value=%v", r.responseCode, r.value)
	})
	// offset 0, quantity 0, empty block data: lastByte = (0+0-1)/8 = 0 in Go, so the guard
	// compares 0 >= 0 and reports an error rather than indexing.
	require.NotPanics(t, func() {
		r := splitBits(member, nil)
		t.Logf("quantity 0, no data -> code=%v err=%v", r.responseCode, r.err)
	})
}

func TestAdvSplitBitsHugeOffset(t *testing.T) {
	member := blockTag{name: "far", tag: coilAt(0, 8), offset: 65535}
	require.NotPanics(t, func() {
		r := splitBits(member, make([]byte, 4))
		assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, r.responseCode)
	})
}

func TestAdvSplitRegistersOffsetOverflow(t *testing.T) {
	member := blockTag{name: "far", tag: holdingRegisterAt(0, 1), offset: 65535}
	require.NotPanics(t, func() {
		r := splitRegisters(context.Background(), member, make([]byte, 8), BigEndianOrder)
		assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, r.responseCode)
	})
}

// Round trip: build a fake device image, let the optimizer merge, split the response and
// compare every tag against what an unoptimized single read of that tag would have produced.
func TestAdvRoundTripMatchesUnoptimizedRead(t *testing.T) {
	image := make([]byte, 2*4096)
	for i := range image {
		image[i] = byte(i*7 + 3)
	}
	ctx := context.Background()
	rng := rand.New(rand.NewSource(42))
	dtypes := []readWriteModel.ModbusDataType{
		readWriteModel.ModbusDataType_INT,
		readWriteModel.ModbusDataType_DINT,
		readWriteModel.ModbusDataType_REAL,
		readWriteModel.ModbusDataType_LREAL,
	}
	for iteration := range 500 {
		n := 1 + rng.Intn(8)
		in := make([]namedTag, 0, n)
		for i := range n {
			dt := dtypes[rng.Intn(len(dtypes))]
			quantity := uint16(1 + rng.Intn(4))
			addr := uint16(rng.Intn(800))
			in = append(in, namedTag{
				name: fmt.Sprintf("t%d", i),
				tag:  newTagFromWireAddress(HoldingRegister, addr, quantity, dt, 1, tagConfig{}, quantity > 1),
			})
		}
		blocks := optimizeReadsWithSpecLimits(in)
		got := map[string]string{}
		for _, b := range blocks {
			data := image[int(b.address)*2 : int(b.address)*2+int(b.quantity)*2]
			for _, r := range splitBlockResponse(ctx, b, apiModel.PlcResponseCode_OK, data, BigEndianOrder) {
				require.Equal(t, apiModel.PlcResponseCode_OK, r.responseCode,
					"iteration %d: %s failed: %v", iteration, r.name, r.err)
				got[r.name] = fmt.Sprint(r.value)
			}
		}
		for _, tag := range in {
			words, err := tag.tag.lengthWords()
			require.NoError(t, err)
			raw := image[int(tag.tag.Address)*2 : int(tag.tag.Address)*2+int(words)*2]
			want, err := ParseRegisters(ctx, raw, tag.tag.Datatype, tag.tag.Quantity, BigEndianOrder, tag.tag.StringLength)
			require.NoError(t, err)
			require.Equal(t, fmt.Sprint(want), got[tag.name],
				"iteration %d: tag %s (addr=%d qty=%d dt=%v) decoded differently through the optimizer",
				iteration, tag.name, tag.tag.Address, tag.tag.Quantity, tag.tag.Datatype)
		}
		require.Len(t, got, len(in), "iteration %d: a tag was dropped", iteration)
	}
}

// The bit areas, same round trip.
func TestAdvRoundTripCoils(t *testing.T) {
	ctx := context.Background()
	rng := rand.New(rand.NewSource(99))
	bits := make([]byte, 512)
	for i := range bits {
		bits[i] = byte(i*31 + 17)
	}
	for iteration := range 500 {
		n := 1 + rng.Intn(8)
		in := make([]namedTag, 0, n)
		for i := range n {
			quantity := uint16(1 + rng.Intn(20))
			addr := uint16(rng.Intn(2000))
			in = append(in, namedTag{
				name: fmt.Sprintf("c%d", i),
				tag:  newTagFromWireAddress(Coil, addr, quantity, readWriteModel.ModbusDataType_BOOL, 1, tagConfig{}, quantity > 1),
			})
		}
		blocks := optimizeReadsWithSpecLimits(in)
		got := map[string][]bool{}
		for _, b := range blocks {
			nbytes := (int(b.quantity) + 7) / 8
			data := make([]byte, nbytes)
			for i := range nbytes {
				// pack the device's coils starting at the block's address
				for bit := range 8 {
					globalBit := int(b.address) + i*8 + bit
					if bits[globalBit/8]&(1<<(globalBit%8)) != 0 {
						data[i] |= 1 << bit
					}
				}
			}
			for _, r := range splitBlockResponse(ctx, b, apiModel.PlcResponseCode_OK, data, BigEndianOrder) {
				require.Equal(t, apiModel.PlcResponseCode_OK, r.responseCode,
					"iteration %d: %s failed: %v", iteration, r.name, r.err)
				got[r.name] = flatten(r.value)
			}
		}
		for _, tag := range in {
			want := make([]bool, 0, tag.tag.Quantity)
			for i := range int(tag.tag.Quantity) {
				b := int(tag.tag.Address) + i
				want = append(want, bits[b/8]&(1<<(b%8)) != 0)
			}
			require.Equal(t, want, got[tag.name],
				"iteration %d: coil tag %s (addr=%d qty=%d) decoded differently through the optimizer",
				iteration, tag.name, tag.tag.Address, tag.tag.Quantity)
		}
	}
}

func flatten(v apiValues.PlcValue) []bool {
	if v == nil {
		return nil
	}
	if !v.IsList() {
		return []bool{v.GetBool()}
	}
	out := make([]bool, 0, len(v.GetList()))
	for _, e := range v.GetList() {
		out = append(out, e.GetBool())
	}
	return out
}

// ---------------------------------------------------------------------------
// 8. odd byte lengths + the byte-swapping orders
// ---------------------------------------------------------------------------

// splitRegisters slices blockData to lengthInBytes, which is ODD for a scalar CHAR and for a
// STRING of odd declared length. ParseRegisters byte-swaps in place for the two *_BYTE_SWAP
// orders and can only swap whole pairs, so the trailing byte of an odd slice never moves - while
// the unoptimized read, which is handed whole registers (2*lengthWords bytes), does move it.
// Same tag, same device, different value.
func TestAdvOddByteLengthUnderByteSwap(t *testing.T) {
	ctx := context.Background()
	// two registers of a device: 0x0102 0x0304
	image := []byte{0x01, 0x02, 0x03, 0x04}

	for _, order := range []ByteOrder{BigEndianOrder, LittleEndianOrder, BigEndianByteSwapOrder, LittleEndianByteSwapOrder} {
		for _, tc := range []struct {
			name         string
			datatype     readWriteModel.ModbusDataType
			stringLength uint16
		}{
			{"CHAR", readWriteModel.ModbusDataType_CHAR, 1},
			{"STRING(3)", readWriteModel.ModbusDataType_STRING, 3},
		} {
			t.Run(order.String()+"/"+tc.name, func(t *testing.T) {
				tag := newTagFromWireAddress(HoldingRegister, 0, 1, tc.datatype, tc.stringLength, tagConfig{}, false)
				words, err := tag.lengthWords()
				require.NoError(t, err)
				t.Logf("lengthInBytes=%d lengthWords=%d",
					lengthInBytes(tag.Datatype, tag.Quantity, tag.StringLength), words)

				// what the unoptimized Reader does: hand ParseRegisters whole registers
				want, err := ParseRegisters(ctx, image[:int(words)*2], tag.Datatype, tag.Quantity, order, tag.StringLength)
				require.NoError(t, err)

				// what the optimizer does: slice the block to lengthInBytes
				got := splitRegisters(ctx, blockTag{name: "t", tag: tag, offset: 0}, image, order)
				require.Equal(t, apiModel.PlcResponseCode_OK, got.responseCode, "%v", got.err)

				assert.Equal(t, fmt.Sprintf("%v", want), fmt.Sprintf("%v", got.value),
					"DEFECT: merging this tag into a block changes the value it decodes to")
			})
		}
	}
}

// The same defect reached the way a user would: two CHAR tags a few registers apart on a
// connection whose default byte order is BIG_ENDIAN_BYTE_SWAP.
func TestAdvCharTagsMergedIntoABlockDecodeWrong(t *testing.T) {
	ctx := context.Background()
	in := named(
		newTagFromWireAddress(HoldingRegister, 0, 1, readWriteModel.ModbusDataType_CHAR, 1, tagConfig{}, false),
		newTagFromWireAddress(HoldingRegister, 3, 1, readWriteModel.ModbusDataType_CHAR, 1, tagConfig{}, false),
	)
	blocks := optimizeReadsWithSpecLimits(in)
	require.Len(t, blocks, 1)
	// registers 0..3 of the device
	image := []byte{'A', 'a', 'B', 'b', 'C', 'c', 'D', 'd'}
	for _, r := range splitBlockResponse(ctx, blocks[0], apiModel.PlcResponseCode_OK, image, BigEndianByteSwapOrder) {
		t.Logf("%s -> %v", r.name, r.value)
	}
	// The unoptimized read of tag1 (register 3) byte-swaps 'D','d' to 'd','D' and reads 'd'.
	results := splitBlockResponse(ctx, blocks[0], apiModel.PlcResponseCode_OK, image, BigEndianByteSwapOrder)
	// Compare the decoded character, not the PlcValue's display form ("CHAR(8bit):d").
	assert.Equal(t, "d", results[1].value.GetString(),
		"a merged CHAR must decode the same byte the unoptimized read would")
}
