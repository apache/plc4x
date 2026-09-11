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
	"cmp"
	"context"
	"slices"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/values"
)

// This is the read optimizer, ported from plc4j's ModbusReadOptimizer.
//
// A request asks for one tag, and without this every tag costs a round trip. That is affordable
// over plain TCP and ruinous behind a TCP<->RTU gateway, where every transaction is serialized
// onto a slow serial line: a few hundred tags turn into a few hundred round trips and the gateway
// falls over. So tags that can be answered by one request are merged into a block read here, and
// the single response is split apart again afterwards.
//
// Nothing in this file talks to a connection. optimizeReads decides what to ask for and
// splitBlockResponse decodes what came back; issuing the requests is the caller's business.
//
// NOTE: nothing calls this yet. The driver still issues one request per tag (see Reader.Read,
// which rejects multi-tag requests outright, and the SingleItemRequestInterceptor the connection
// installs). Wiring it in changes per-tag failure semantics - a failed block fails every tag in
// it - and so is a separate, separately reviewable step. Until then this file changes no
// behaviour.

// namedTag pairs a tag with the name the read request knows it by. Tags travel by value here, the
// way they do everywhere else in this package - castToModbusTagFromPlcTag hands out a value too.
type namedTag struct {
	name string
	tag  modbusTag
}

// blockTag is one member of a block: the tag itself plus where its data sits inside the block's
// response.
type blockTag struct {
	name string
	tag  modbusTag
	// offset is how far into the block this tag's data starts, counted in whatever the area
	// addresses - bits for the coil and discrete-input areas, registers for the three register
	// areas. It is derivable from the two addresses; it is spelled out so that splitting a
	// response never has to redo the subtraction, and so that a caller can see the layout.
	offset uint16
}

// readBlock is one request the optimizer wants sent: a run of addresses in one area, on one unit,
// covering every tag in tags.
type readBlock struct {
	tagType TagType
	// unitId is the unit every tag in this block declared, or nil when they all left it to the
	// connection's default. It has to be carried through to the request: a block read reaches
	// exactly one unit, and dropping the unit-id here would silently send the merged request to
	// the connection's default unit instead (plc4j GitHub issue #2686).
	unitId *uint8
	// address is the wire address the block starts at.
	address uint16
	// quantity is how many addresses the block covers - coils for the two bit areas, registers
	// for the three register areas. This is the count the request carries.
	quantity uint16
	tags     []blockTag
}

// tagReadResult is what one tag of a block ends up with once the response has been split.
type tagReadResult struct {
	name         string
	responseCode apiModel.PlcResponseCode
	// value is set only for a result that is OK.
	value apiValues.PlcValue
	// err is why responseCode is not OK, kept for logging. The response itself carries only the
	// code, so without this the reason a tag failed would be lost at the point it is known.
	err error
}

// areaOrder is the order the blocks of the individual areas are emitted in. Five separate areas,
// never merged with one another: they are addressed independently and read with different function
// codes, so address 7 means five unrelated things. The order itself is the one plc4j's
// optimizeReads walks, kept so both implementations produce the same request sequence.
var areaOrder = []TagType{Coil, HoldingRegister, InputRegister, ExtendedRegister, DiscreteInput}

// unitGroupKey identifies one group of tags that may be merged with one another: same area, same
// unit. A tag that declared no unit-id forms its own group rather than joining the group of the
// unit that happens to be the connection's default - the connection may be re-configured, and the
// tag asked for "whatever the connection uses", which is not the same statement as a number.
type unitGroupKey struct {
	tagType   TagType
	hasUnitId bool
	unitId    uint8
}

// optimizeReadsWithSpecLimits merges tags using the largest request the modbus specification
// allows, which is also what plc4j defaults to.
func optimizeReadsWithSpecLimits(tags []namedTag) []readBlock {
	return optimizeReads(tags, maxCoilQuantity, maxRegisterQuantity)
}

// optimizeReads groups the given tags by area and unit and cuts each group into the blocks that
// have to be read, in a deterministic order.
//
// maxCoilsPerRequest and maxRegistersPerRequest bound how wide one block may get; see
// maxCoilQuantity and maxRegisterQuantity for the ceilings the protocol itself imposes.
func optimizeReads(tags []namedTag, maxCoilsPerRequest uint16, maxRegistersPerRequest uint16) []readBlock {
	groups := make(map[unitGroupKey][]namedTag, len(tags))
	for _, tag := range tags {
		key := unitGroupKey{tagType: tag.tag.TagType}
		if tag.tag.UnitId != nil {
			key.hasUnitId, key.unitId = true, *tag.tag.UnitId
		}
		groups[key] = append(groups[key], tag)
	}

	// Map iteration order is random, and a driver that sends its requests in a different order on
	// every read is untestable and miserable to follow in a packet capture.
	keys := make([]unitGroupKey, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	slices.SortFunc(keys, compareUnitGroups)

	var blocks []readBlock
	for _, key := range keys {
		group := groups[key]
		// Sorted by address, so that the window below only ever has to look forwards. The name
		// breaks ties, for the same reason the keys are sorted.
		slices.SortFunc(group, func(a, b namedTag) int {
			return cmp.Or(cmp.Compare(a.tag.Address, b.tag.Address), cmp.Compare(a.name, b.name))
		})
		blocks = append(blocks, blocksOf(key, group, maxPerRequestFor(key.tagType, maxCoilsPerRequest, maxRegistersPerRequest))...)
	}
	return blocks
}

// compareUnitGroups orders the groups: areas in the order plc4j emits them, and within an area the
// tags that named no unit before those that did, then by unit-id.
func compareUnitGroups(a, b unitGroupKey) int {
	return cmp.Or(
		cmp.Compare(slices.Index(areaOrder, a.tagType), slices.Index(areaOrder, b.tagType)),
		cmp.Compare(boolAsInt(a.hasUnitId), boolAsInt(b.hasUnitId)),
		cmp.Compare(a.unitId, b.unitId),
	)
}

func boolAsInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// maxPerRequestFor is the ceiling that applies to one block of the given area. The two bit areas
// are counted in coils, everything else in registers.
// maxExtendedRegisterQuantity is the widest extended-register block that still answers within the
// 253-byte Modbus PDU limit. FC 0x14 spends a byte count plus two bytes of per-item framing that
// FC 0x03/0x04 do not, and a block crossing a 10000-register file boundary pays that framing twice.
const maxExtendedRegisterQuantity = 123

func maxPerRequestFor(tagType TagType, maxCoilsPerRequest uint16, maxRegistersPerRequest uint16) uint16 {
	switch tagType {
	case Coil, DiscreteInput:
		return maxCoilsPerRequest
	case ExtendedRegister:
		// Extended registers travel over FC 0x14 (Read File Record), whose response carries a
		// byte count and then per-item framing (data length + reference type) that the plain
		// register function codes do not. A full 125-register block therefore answers with a PDU
		// over the 253-byte Modbus limit, and a block straddling a 10000-register file boundary
		// becomes two items and overflows by more. 123 is the widest that fits both shapes.
		if maxRegistersPerRequest > maxExtendedRegisterQuantity {
			return maxExtendedRegisterQuantity
		}
		return maxRegistersPerRequest
	default:
		return maxRegistersPerRequest
	}
}

// blocksOf cuts one area/unit group, which arrives sorted by address, into blocks.
//
// The window is anchored at the first tag of a block and is exactly maxPerRequest addresses wide.
// Every following tag whose end fits inside that window joins the block, whether or not anything
// was addressed in between; the first one that does not fit anchors the next window.
//
// There is deliberately no contiguity requirement and no separate gap tolerance: the width of the
// window *is* the gap tolerance (plc4j ModbusReadOptimizer.optimizeRegisters). Tags at register 0
// and register 100 therefore become one 101 register read - reading 99 registers nobody asked for
// costs a handful of bytes, while the second round trip it saves costs a whole transaction on a
// serial line. Tags at 0 and 2100 stay two reads.
func blocksOf(key unitGroupKey, group []namedTag, maxPerRequest uint16) []readBlock {
	var blocks []readBlock
	var current []namedTag
	// Held as int rather than uint16: the end of the window is the first tag's address plus the
	// ceiling, which runs past the end of the address space for a tag near the top of it.
	var first, last, window int

	for _, entry := range group {
		address := int(entry.tag.Address)
		end := address + int(tagSpan(entry.tag))
		switch {
		case len(current) == 0:
			// The first tag of a block sets the block's end with no cap check of its own, so a
			// single tag wider than maxPerRequest passes through unsplit. That is the plc4j
			// contract and it is the right one: this optimizer only ever merges requests, it
			// never takes one apart, and quietly halving a tag's read would hand the caller
			// half a value. A tag that wide is rejected when it is parsed (see
			// validateAddressAndQuantity), so it can only get here from a caller that built it
			// by hand.
			first, last, window = address, end, address+int(maxPerRequest)
		case end > window:
			blocks = append(blocks, newReadBlock(key, first, last, current))
			current = nil
			first, last, window = address, end, address+int(maxPerRequest)
		default:
			// Not max(last, end) for tidiness: the group is sorted by address, but a wide tag
			// followed by a narrow one inside it would otherwise shrink the block.
			last = max(last, end)
		}
		current = append(current, entry)
	}
	if len(current) > 0 {
		blocks = append(blocks, newReadBlock(key, first, last, current))
	}
	return blocks
}

// newReadBlock assembles the block covering [first, last) out of the tags that fell into it.
func newReadBlock(key unitGroupKey, first int, last int, tags []namedTag) readBlock {
	// first/last are ints because an oversized pass-through tag can push the span past what a
	// uint16 holds. Narrowing that unchecked would wrap: a 65536-register span becomes a request
	// for nothing, and a block anchored near the top of the address space would ask the device to
	// read past 65535. Such a tag cannot come from the address parser, which rejects it, only from
	// a caller that built one by hand. Clamping the SPAN (not the end address) is what matters:
	// a 65536-register span narrows to quantity 0, a request for nothing, which a device answers
	// happily and which therefore reads as success. A clamped-but-oversized quantity is rejected
	// by the device instead, which is the honest outcome for a request that should not exist.
	span := last - first
	if span > maxWireAddress {
		span = maxWireAddress
	}
	block := readBlock{
		tagType:  key.tagType,
		address:  uint16(first),
		quantity: uint16(span),
		tags:     make([]blockTag, 0, len(tags)),
	}
	if key.hasUnitId {
		unitId := key.unitId
		block.unitId = &unitId
	}
	for _, entry := range tags {
		block.tags = append(block.tags, blockTag{
			name:   entry.name,
			tag:    entry.tag,
			offset: entry.tag.Address - block.address,
		})
	}
	return block
}

// tagSpan is how many addresses a tag occupies in its area: coils in a bit area, where one element
// is one address, and registers everywhere else.
//
// This is the same count the unoptimized read puts into its request (see readRequestPdu), which is
// what makes a block that covers a tag's span enough to answer it.
func tagSpan(tag modbusTag) uint64 {
	return registerCountOf(tag.TagType, tag.Datatype, uint64(tag.Quantity), tag.StringLength)
}

// splitBlockResponse takes the payload one block read came back with and answers every tag of the
// block out of it.
//
// blockData is the raw payload of the response: the packed coil bytes for the two bit areas, the
// registers for the other three. defaultByteOrder is the connection's, used for every tag that did
// not declare one of its own.
func splitBlockResponse(ctx context.Context, block readBlock, blockResponseCode apiModel.PlcResponseCode, blockData []byte, defaultByteOrder ByteOrder) []tagReadResult {
	results := make([]tagReadResult, 0, len(block.tags))
	if blockResponseCode != apiModel.PlcResponseCode_OK {
		// One request answered for all of them, so its failure is every member's failure. Merging
		// reads means a device that rejects one address takes the whole block down with it - the
		// price of the optimization, and the reason the block's own code is reported rather than
		// a code invented here.
		for _, member := range block.tags {
			results = append(results, tagReadResult{name: member.name, responseCode: blockResponseCode})
		}
		return results
	}
	for _, member := range block.tags {
		switch block.tagType {
		case Coil, DiscreteInput:
			results = append(results, splitBits(member, blockData))
		default:
			results = append(results, splitRegisters(ctx, member, blockData, defaultByteOrder))
		}
	}
	return results
}

// splitBits answers one tag of a coil or discrete-input block out of the packed bits that came
// back.
func splitBits(member blockTag, blockData []byte) tagReadResult {
	if member.tag.Datatype != readWriteModel.ModbusDataType_BOOL {
		// A coil carries a single bit, and assembling coils into wider types is not implemented.
		// Answering with the first bit would look like a successful read of something else
		// entirely, so this is reported instead (plc4j ModbusReadOptimizer.splitResponse).
		return tagReadResult{
			name:         member.name,
			responseCode: apiModel.PlcResponseCode_UNSUPPORTED,
			err:          errors.Errorf("reading coils/discrete inputs as %s is not supported, only BOOL is", member.tag.Datatype),
		}
	}
	// An array tag occupies as many coils as it has elements, and the block covers all of them -
	// so every element is extracted, not just the first (plc4j GitHub issue #2060).
	firstBit := int(member.offset)
	numberOfElements := int(member.tag.Quantity)
	if lastByte := (firstBit + numberOfElements - 1) / 8; lastByte >= len(blockData) {
		return tagReadResult{
			name:         member.name,
			responseCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
			err:          errors.Errorf("the device answered with %d bytes of coils, which doesn't reach coil %d", len(blockData), firstBit+numberOfElements-1),
		}
	}
	bits := make([]apiValues.PlcValue, 0, numberOfElements)
	for i := range numberOfElements {
		// Coils are packed least significant bit first, so a run of them crosses into the next
		// byte at bit 8 rather than restarting there.
		bit := firstBit + i
		bits = append(bits, values.NewPlcBOOL(blockData[bit/8]&(1<<(bit%8)) != 0))
	}
	if numberOfElements == 1 {
		// A single coil reads as a plain BOOL, not as a list of one - the same shape
		// ParseRegisters gives a single value, so merging a tag into a block never changes the
		// type the caller sees.
		return tagReadResult{name: member.name, responseCode: apiModel.PlcResponseCode_OK, value: bits[0]}
	}
	return tagReadResult{name: member.name, responseCode: apiModel.PlcResponseCode_OK, value: values.NewPlcList(bits)}
}

// splitRegisters answers one tag of a register block out of the registers that came back.
func splitRegisters(ctx context.Context, member blockTag, blockData []byte, defaultByteOrder ByteOrder) tagReadResult {
	// The block is a flat run of registers, so a member starts a whole number of registers into
	// it. That even offset is what lets ParseRegisters be reused unchanged: it decodes positionally
	// from bit 0 of the slice it is handed and never learns an address, so a slice that starts
	// mid-register would silently decode garbage.
	byteOffset := int(member.offset) * 2
	// Whole REGISTERS, not the tag's packed byte length. The two differ for any type whose
	// payload is an odd number of bytes - a scalar CHAR, a STRING of odd declared length - and
	// the difference is not cosmetic: ParseRegisters byte-swaps in place for the two *_BYTE_SWAP
	// orders and can only swap whole pairs, so an odd-length slice leaves its last byte unswapped
	// while the unoptimized read, which is handed the device's whole registers, swaps it. Slicing
	// to the payload length would therefore decode a DIFFERENT value than the same tag read on its
	// own, silently and with an OK response code. ParseRegisters ignores bytes it does not need,
	// which is exactly what the unoptimized path already relies on.
	byteLength := int(tagSpan(member.tag)) * 2
	if byteOffset+byteLength > len(blockData) {
		return tagReadResult{
			name:         member.name,
			responseCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
			err:          errors.Errorf("the device answered with %d bytes, which doesn't cover the %d bytes at offset %d", len(blockData), byteLength, byteOffset),
		}
	}
	// The byte order is the tag's own if it declared one - a block may well mix them, and the
	// swapping modes are handled inside ParseRegisters rather than being redone here.
	value, err := ParseRegisters(ctx, blockData[byteOffset:byteOffset+byteLength], member.tag.Datatype,
		member.tag.Quantity, member.tag.resolveByteOrder(defaultByteOrder), member.tag.StringLength)
	if err != nil {
		return tagReadResult{
			name:         member.name,
			responseCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
			err:          errors.Wrap(err, "error parsing data item"),
		}
	}
	return tagReadResult{name: member.name, responseCode: apiModel.PlcResponseCode_OK, value: value}
}
